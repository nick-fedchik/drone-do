#!/usr/bin/env python3
"""
Combined landing controller: external Logitech C615 tracks both the ArUco
landing pad and the drone (TensorRT YOLOv8n), then streams relative alignment
commands to droneemu over local Unix IPC.
"""

import argparse
import json
import math
import os
import socket
import sys
import threading
import time
from typing import Optional, Tuple

import cv2
import numpy as np

from aruco_lander import clamp, default_socket_path, detect_aruco, send_ipc_command


def bbox_center_and_area(bbox_xywh: list, frame_w: int, frame_h: int) -> Tuple[float, float, float]:
    x, y, w, h = bbox_xywh
    cx = x + w / 2.0
    cy = y + h / 2.0
    area_ratio = (w * h) / float(max(1, frame_w * frame_h))
    return cx, cy, area_ratio


def compute_relative_control(
    pad_corners: Optional[np.ndarray],
    drone_bbox_xywh: Optional[list],
    frame_w: int,
    frame_h: int,
    kp_xy: float = 1.0,
    kp_yaw: float = 0.5,
    hover_throttle: float = 0.5,
    descent_throttle: float = 0.35,
    max_tilt: float = 0.35,
    smoothed_err: Optional[Tuple[float, float]] = None,
    alpha: float = 0.4,
    invert_roll: bool = False,
    invert_pitch: bool = False,
    invert_yaw: bool = False,
) -> Tuple[dict, Optional[Tuple[float, float]]]:
    """Computes clamped, smoothed controls to steer the detected drone toward the ArUco pad."""
    if pad_corners is None or drone_bbox_xywh is None or frame_w <= 0 or frame_h <= 0:
        return {"throttle": 0.0, "roll": 0.0, "pitch": 0.0, "yaw": 0.0, "aligned": False}, None

    pad_center = pad_corners.mean(axis=0)
    drone_cx, drone_cy, _ = bbox_center_and_area(drone_bbox_xywh, frame_w, frame_h)

    raw_x = (float(pad_center[0]) - drone_cx) / (frame_w / 2.0)
    raw_y = (drone_cy - float(pad_center[1])) / (frame_h / 2.0)
    if smoothed_err is not None:
        err_x = alpha * raw_x + (1.0 - alpha) * smoothed_err[0]
        err_y = alpha * raw_y + (1.0 - alpha) * smoothed_err[1]
    else:
        err_x, err_y = raw_x, raw_y

    top_edge = pad_corners[1] - pad_corners[0]
    yaw_rad = math.atan2(float(top_edge[1]), float(top_edge[0]))
    err_yaw = clamp(yaw_rad / math.pi, -1.0, 1.0)

    tilt_lim = clamp(max_tilt, 0.0, 1.0)
    roll = clamp((-1.0 if invert_roll else 1.0) * kp_xy * err_x, -tilt_lim, tilt_lim)
    pitch = clamp((-1.0 if invert_pitch else 1.0) * kp_xy * err_y, -tilt_lim, tilt_lim)
    yaw = clamp((-1.0 if invert_yaw else 1.0) * kp_yaw * err_yaw, -tilt_lim, tilt_lim)

    aligned = abs(err_x) < 0.10 and abs(err_y) < 0.10
    throttle = descent_throttle if aligned else hover_throttle

    return {
        "throttle": round(throttle, 3),
        "roll": round(roll, 3),
        "pitch": round(pitch, 3),
        "yaw": round(yaw, 3),
        "aligned": aligned,
    }, (err_x, err_y)


class FailSafeTracker:
    """Tracks alignment phase (SEARCHING/ALIGNING/DESCENDING) and cuts throttle on timeout."""

    def __init__(
        self,
        loss_timeout: float = 0.25,
        hover_throttle: float = 0.5,
        descent_throttle: float = 0.35,
        align_frames: int = 3,
    ) -> None:
        self.loss_timeout = loss_timeout
        self.hover_throttle = hover_throttle
        self.descent_throttle = descent_throttle
        self.align_frames = max(1, align_frames)
        self.align_streak = 0
        self.phase = "SEARCHING"
        self.last_valid_at: Optional[float] = None
        self.smoothed_err: Optional[Tuple[float, float]] = None

    def step(
        self,
        now: float,
        pad_corners: Optional[np.ndarray],
        drone_bbox_xywh: Optional[list],
        frame_w: int,
        frame_h: int,
        max_tilt: float = 0.35,
        invert_roll: bool = False,
        invert_pitch: bool = False,
        invert_yaw: bool = False,
    ) -> dict:
        if pad_corners is not None and drone_bbox_xywh is not None:
            ctrl, self.smoothed_err = compute_relative_control(
                pad_corners,
                drone_bbox_xywh,
                frame_w,
                frame_h,
                hover_throttle=self.hover_throttle,
                descent_throttle=self.descent_throttle,
                max_tilt=max_tilt,
                smoothed_err=self.smoothed_err,
                invert_roll=invert_roll,
                invert_pitch=invert_pitch,
                invert_yaw=invert_yaw,
            )
            self.last_valid_at = now
            if ctrl["aligned"]:
                self.align_streak += 1
            else:
                self.align_streak = 0
            if self.align_streak >= self.align_frames:
                self.phase = "DESCENDING"
                ctrl["throttle"] = round(self.descent_throttle, 3)
            else:
                self.phase = "ALIGNING"
                ctrl["throttle"] = round(self.hover_throttle, 3)
            ctrl["phase"] = self.phase
            return ctrl

        self.smoothed_err = None
        self.align_streak = 0
        if self.last_valid_at is not None and (now - self.last_valid_at) <= self.loss_timeout:
            self.phase = "HOLD_LOSS"
            return {"throttle": self.hover_throttle, "roll": 0.0, "pitch": 0.0, "yaw": 0.0, "aligned": False, "phase": self.phase}
        self.phase = "SEARCHING"
        return {"throttle": 0.0, "roll": 0.0, "pitch": 0.0, "yaw": 0.0, "aligned": False, "phase": self.phase}


def main() -> int:
    from trt_detector import TrtYoloDetector

    parser = argparse.ArgumentParser(description="Combined pad + drone TensorRT lander for droneemu")
    parser.add_argument("--engine", default="/home/mfedchyk/models/yolov8n.engine", help="TensorRT .engine path")
    parser.add_argument("--socket", default=default_socket_path(), help="droneemu Unix socket path")
    parser.add_argument("--camera", type=int, default=0, help="V4L2 camera index (default: 0)")
    parser.add_argument("--marker-id", type=int, default=0, help="ArUco DICT_4X4_50 pad ID")
    parser.add_argument("--target-label", default="", help="optional YOLO class filter (e.g. syma_x23w)")
    parser.add_argument("--conf", type=float, default=0.55, help="YOLO confidence threshold")
    parser.add_argument("--max-tilt", type=float, default=0.35, help="max roll/pitch/yaw amplitude [0,1]")
    parser.add_argument("--hover-throttle", type=float, default=0.5, help="neutral hover throttle [0,1] (0.5 -> 0x80)")
    parser.add_argument("--descent-throttle", type=float, default=0.35, help="gentle descent throttle when aligned")
    parser.add_argument("--loss-timeout", type=float, default=0.25, help="seconds before fail-safe cuts throttle")
    parser.add_argument("--invert-roll", action="store_true", help="invert roll sign")
    parser.add_argument("--invert-pitch", action="store_true", help="invert pitch sign")
    parser.add_argument("--invert-yaw", action="store_true", help="invert yaw sign")
    parser.add_argument("--frames", type=int, default=0, help="stop after N frames (0 = until Ctrl-C)")
    parser.add_argument("--hz", type=float, default=10.0, help="control loop rate in Hz")
    parser.add_argument("--status-file", default="", help="optional JSON status file updated every frame")
    parser.add_argument("--preview-port", type=int, default=0, help="optional HTTP MJPEG preview port (e.g. 8081)")
    args = parser.parse_args()

    detector = TrtYoloDetector(args.engine)
    cap = cv2.VideoCapture(args.camera, cv2.CAP_V4L2)
    if not cap.isOpened():
        detector.close()
        print(f"Failed to open /dev/video{args.camera}", file=sys.stderr)
        return 1
    cap.set(cv2.CAP_PROP_FOURCC, cv2.VideoWriter_fourcc(*"MJPG"))
    cap.set(cv2.CAP_PROP_FRAME_WIDTH, 1280)
    cap.set(cv2.CAP_PROP_FRAME_HEIGHT, 720)
    cap.set(cv2.CAP_PROP_FPS, 30)

    sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
    sock.settimeout(3.0)
    try:
        sock.connect(args.socket)
    except OSError as exc:
        cap.release()
        detector.close()
        print(f"Connect to emulator socket {args.socket}: {exc}", file=sys.stderr)
        return 1

    preview_buf = None
    preview_server = None
    if args.preview_port > 0:
        from http.server import ThreadingHTTPServer
        from stream_preview import FrameBuffer, make_handler
        preview_buf = FrameBuffer()
        preview_server = ThreadingHTTPServer(("0.0.0.0", args.preview_port), make_handler(preview_buf))
        threading.Thread(target=preview_server.serve_forever, daemon=True).start()
        print(f"MJPEG preview ready at http://jetson:{args.preview_port}/stream.mjpg", flush=True)

    seq = 0
    period = 1.0 / max(1.0, args.hz)
    tracker = FailSafeTracker(
        loss_timeout=args.loss_timeout,
        hover_throttle=args.hover_throttle,
        descent_throttle=args.descent_throttle,
    )
    try:
        while args.frames == 0 or seq < args.frames:
            t0 = time.monotonic()
            ret, frame = cap.read()
            if not ret or frame is None:
                continue

            h, w = frame.shape[:2]
            pad_det = detect_aruco(frame, target_id=args.marker_id)
            pad_corners = pad_det[1] if pad_det is not None else None

            dets = detector.detect(frame, conf_thresh=args.conf)
            if args.target_label:
                dets = [d for d in dets if d["label"] == args.target_label]
            drone_bbox = dets[0]["bbox_xywh"] if dets else None
            drone_conf = dets[0]["confidence"] if dets else 0.0

            ctrl = tracker.step(
                t0,
                pad_corners,
                drone_bbox,
                w,
                h,
                max_tilt=args.max_tilt,
                invert_roll=args.invert_roll,
                invert_pitch=args.invert_pitch,
                invert_yaw=args.invert_yaw,
            )
            seq += 1
            resp = send_ipc_command(sock, seq, ctrl)
            if args.status_file:
                status_payload = {
                    "seq": seq,
                    "pad": pad_det is not None,
                    "drone": drone_bbox is not None,
                    "drone_conf": drone_conf,
                    "ctrl": ctrl,
                    "emu_cmds": resp["snapshot"]["commands_received"],
                }
                tmp_path = args.status_file + ".tmp"
                with open(tmp_path, "w", encoding="utf-8") as sf:
                    json.dump(status_payload, sf)
                os.replace(tmp_path, args.status_file)
            if preview_buf is not None:
                if pad_corners is not None:
                    pts = pad_corners.astype(int).reshape((-1, 1, 2))
                    cv2.polylines(frame, [pts], True, (255, 200, 0), 2)
                    cv2.putText(frame, f"Pad {pad_det[0]}", tuple(pts[0][0]), cv2.FONT_HERSHEY_SIMPLEX, 0.7, (255, 200, 0), 2)
                if drone_bbox is not None:
                    x, y, bw, bh = [int(v) for v in drone_bbox]
                    cv2.rectangle(frame, (x, y), (x + bw, y + bh), (0, 255, 0), 2)
                    cv2.putText(frame, f"syma_x23w {drone_conf:.2f}", (x, max(20, y - 8)),
                                cv2.FONT_HERSHEY_SIMPLEX, 0.7, (0, 255, 0), 2)
                ok, enc = cv2.imencode(".jpg", frame, [int(cv2.IMWRITE_JPEG_QUALITY), 80])
                if ok:
                    preview_buf.set(enc.tobytes())
            print(
                f"seq={seq} pad={pad_det is not None} drone={drone_bbox is not None} "
                f"thr={ctrl['throttle']:.2f} roll={ctrl['roll']:+.2f} "
                f"pitch={ctrl['pitch']:+.2f} emu_cmds={resp['snapshot']['commands_received']}"
            )
            dt_loop = time.monotonic() - t0
            if dt_loop < period:
                time.sleep(period - dt_loop)
    except KeyboardInterrupt:
        pass
    finally:
        if preview_server is not None:
            preview_server.server_close()
        cap.release()
        detector.close()
        sock.close()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
