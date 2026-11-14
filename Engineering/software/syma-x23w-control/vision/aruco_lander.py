#!/usr/bin/env python3
"""
ArUco landing-pad tracker for Jetson + Logitech C615 -> droneemu local Unix IPC.

Reads frames from V4L2 camera (/dev/video0) or a test image, detects an ArUco
marker (DICT_4X4_50 by default), computes normalized control corrections, and
streams IPC command frames to droneemu. No real aircraft commands are sent.
"""

import argparse
import datetime as dt
import json
import math
import os
import socket
import sys
import time
from typing import Optional, Tuple

import cv2
import numpy as np


def default_socket_path() -> str:
    cache_home = os.environ.get("XDG_CACHE_HOME") or os.path.join(os.path.expanduser("~"), ".cache")
    return os.path.join(cache_home, "drone-do", "syma-x23w.sock")


def clamp(value: float, low: float, high: float) -> float:
    return max(low, min(high, value))


def detect_aruco(
    frame: np.ndarray,
    dict_id: int = cv2.aruco.DICT_4X4_50,
    target_id: Optional[int] = None,
) -> Optional[Tuple[int, np.ndarray]]:
    dictionary = cv2.aruco.getPredefinedDictionary(dict_id)
    detector = cv2.aruco.ArucoDetector(dictionary, cv2.aruco.DetectorParameters())
    corners, ids, _ = detector.detectMarkers(frame)
    if ids is None or len(ids) == 0:
        return None
    for i, marker_id in enumerate(ids.flatten()):
        if target_id is None or int(marker_id) == target_id:
            return int(marker_id), corners[i][0]
    return None


def compute_control(
    marker_corners: Optional[np.ndarray],
    frame_width: int,
    frame_height: int,
    kp_xy: float = 0.8,
    kp_yaw: float = 0.5,
    hover_throttle: float = 0.45,
    target_area_ratio: float = 0.12,
) -> dict:
    """Maps ArUco marker offset/orientation in image space to normalized [-1,1] / [0,1] controls."""
    if marker_corners is None or frame_width <= 0 or frame_height <= 0:
        return {"throttle": 0.0, "roll": 0.0, "pitch": 0.0, "yaw": 0.0, "aligned": False}

    center = marker_corners.mean(axis=0)
    err_x = (float(center[0]) - frame_width / 2.0) / (frame_width / 2.0)
    err_y = (frame_height / 2.0 - float(center[1])) / (frame_height / 2.0)

    top_edge = marker_corners[1] - marker_corners[0]
    yaw_rad = math.atan2(float(top_edge[1]), float(top_edge[0]))
    err_yaw = clamp(yaw_rad / math.pi, -1.0, 1.0)

    area = float(cv2.contourArea(marker_corners.astype(np.float32)))
    area_ratio = area / float(frame_width * frame_height)

    roll = clamp(kp_xy * err_x, -1.0, 1.0)
    pitch = clamp(kp_xy * err_y, -1.0, 1.0)
    yaw = clamp(kp_yaw * err_yaw, -1.0, 1.0)

    aligned = abs(err_x) < 0.15 and abs(err_y) < 0.15 and abs(err_yaw) < 0.15
    if aligned and area_ratio >= target_area_ratio:
        throttle = 0.0
    elif aligned:
        throttle = clamp(hover_throttle * (1.0 - min(1.0, area_ratio / target_area_ratio) * 0.6), 0.15, 1.0)
    else:
        throttle = hover_throttle

    return {
        "throttle": round(throttle, 3),
        "roll": round(roll, 3),
        "pitch": round(pitch, 3),
        "yaw": round(yaw, 3),
        "aligned": aligned,
        "area_ratio": round(area_ratio, 4),
    }


def send_ipc_command(sock: socket.socket, sequence: int, ctrl: dict) -> dict:
    now = dt.datetime.now(dt.timezone.utc).isoformat().replace("+00:00", "Z")
    payload = {
        "version": 1,
        "type": "command",
        "sequence": sequence,
        "frame": {
            "at": now,
            "throttle": ctrl["throttle"],
            "roll": ctrl["roll"],
            "pitch": ctrl["pitch"],
            "yaw": ctrl["yaw"],
        },
    }
    sock.sendall((json.dumps(payload) + "\n").encode("utf-8"))
    buf = b""
    while b"\n" not in buf:
        chunk = sock.recv(4096)
        if not chunk:
            raise RuntimeError("IPC connection closed by emulator")
        buf += chunk
    line = buf.split(b"\n", 1)[0]
    resp = json.loads(line.decode("utf-8"))
    if resp.get("type") == "error":
        raise RuntimeError(f"emulator rejected command: {resp.get('error')}")
    return resp


def generate_marker_png(path: str, marker_id: int = 0, size: int = 400) -> None:
    dictionary = cv2.aruco.getPredefinedDictionary(cv2.aruco.DICT_4X4_50)
    img = cv2.aruco.generateImageMarker(dictionary, marker_id, size)
    bordered = cv2.copyMakeBorder(img, 40, 40, 40, 40, cv2.BORDER_CONSTANT, value=255)
    cv2.imwrite(path, bordered)


def main() -> int:
    parser = argparse.ArgumentParser(description="ArUco visual landing tracker for droneemu IPC")
    parser.add_argument("--socket", default=default_socket_path(), help="droneemu Unix socket path")
    parser.add_argument("--camera", type=int, default=0, help="V4L2 camera index (default: 0)")
    parser.add_argument("--image", default="", help="optional static image instead of live camera")
    parser.add_argument("--marker-id", type=int, default=0, help="ArUco DICT_4X4_50 marker ID (-1 for any)")
    parser.add_argument("--frames", type=int, default=0, help="stop after N frames (0 = run until Ctrl-C)")
    parser.add_argument("--hz", type=float, default=10.0, help="control loop rate in Hz")
    parser.add_argument("--generate-marker", default="", help="save printable ArUco marker PNG and exit")
    args = parser.parse_args()

    if args.generate_marker:
        marker_id = 0 if args.marker_id < 0 else args.marker_id
        generate_marker_png(args.generate_marker, marker_id=marker_id)
        print(f"Saved ArUco DICT_4X4_50 id={marker_id} to {args.generate_marker}")
        return 0

    target_id = None if args.marker_id < 0 else args.marker_id
    cap = None
    static_img = None
    if args.image:
        static_img = cv2.imread(args.image)
        if static_img is None:
            print(f"Failed to read image: {args.image}", file=sys.stderr)
            return 1
    else:
        cap = cv2.VideoCapture(args.camera, cv2.CAP_V4L2)
        if not cap.isOpened():
            print(f"Failed to open camera /dev/video{args.camera}", file=sys.stderr)
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
        print(f"Connect to emulator socket {args.socket}: {exc}", file=sys.stderr)
        if cap is not None:
            cap.release()
        return 1

    seq = 0
    period = 1.0 / max(1.0, args.hz)
    try:
        while args.frames == 0 or seq < args.frames:
            t0 = time.monotonic()
            if static_img is not None:
                frame = static_img
            else:
                ret, frame = cap.read()
                if not ret or frame is None:
                    continue

            h, w = frame.shape[:2]
            det = detect_aruco(frame, target_id=target_id)
            corners = det[1] if det is not None else None
            ctrl = compute_control(corners, w, h)
            seq += 1
            resp = send_ipc_command(sock, seq, ctrl)
            status = f"marker={det[0]}" if det is not None else "marker=NONE"
            print(
                f"seq={seq} {status} thr={ctrl['throttle']:.2f} "
                f"roll={ctrl['roll']:+.2f} pitch={ctrl['pitch']:+.2f} yaw={ctrl['yaw']:+.2f} "
                f"emu_cmds={resp['snapshot']['commands_received']}"
            )
            dt_loop = time.monotonic() - t0
            if dt_loop < period:
                time.sleep(period - dt_loop)
    except KeyboardInterrupt:
        pass
    finally:
        if cap is not None:
            cap.release()
        sock.close()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
