#!/usr/bin/env python3
"""
Lightweight HTTP MJPEG preview server on Jetson for viewing Logitech C615
(with optional TensorRT YOLO + ArUco overlay) in VLC or a browser on AORUS5.
"""

import argparse
from http.server import BaseHTTPRequestHandler, ThreadingHTTPServer
import sys
import threading
import time

import cv2

from aruco_lander import detect_aruco


class FrameBuffer:
    def __init__(self) -> None:
        self.lock = threading.Lock()
        self.jpeg: bytes = b""

    def set(self, data: bytes) -> None:
        with self.lock:
            self.jpeg = data

    def get(self) -> bytes:
        with self.lock:
            return self.jpeg


def make_handler(buf: FrameBuffer):
    class Handler(BaseHTTPRequestHandler):
        def do_GET(self) -> None:
            if self.path == "/" or self.path == "/index.html":
                body = b'<html><body style="margin:0;background:#111"><img src="/stream.mjpg" style="max-width:100vw;max-height:100vh"/></body></html>'
                self.send_response(200)
                self.send_header("Content-Type", "text/html")
                self.send_header("Content-Length", str(len(body)))
                self.end_headers()
                self.wfile.write(body)
                return
            if self.path != "/stream.mjpg":
                self.send_error(404)
                return

            self.send_response(200)
            self.send_header("Content-Type", "multipart/x-mixed-replace; boundary=frame")
            self.end_headers()
            try:
                while True:
                    frame = buf.get()
                    if frame:
                        self.wfile.write(b"--frame\r\nContent-Type: image/jpeg\r\nContent-Length: " +
                                         str(len(frame)).encode() + b"\r\n\r\n" + frame + b"\r\n")
                    time.sleep(0.03)
            except (BrokenPipeError, ConnectionResetError):
                pass

        def log_message(self, format: str, *args) -> None:
            return

    return Handler


def main() -> int:
    parser = argparse.ArgumentParser(description="HTTP MJPEG preview server for VLC/browser")
    parser.add_argument("--camera", type=int, default=0, help="V4L2 camera index")
    parser.add_argument("--port", type=int, default=8081, help="HTTP listen port")
    parser.add_argument("--engine", default="", help="optional TensorRT .engine for live overlay")
    parser.add_argument("--conf", type=float, default=0.55, help="YOLO confidence threshold")
    parser.add_argument("--aruco", action="store_true", help="overlay ArUco DICT_4X4_50 markers")
    args = parser.parse_args()

    detector = None
    if args.engine:
        from trt_detector import TrtYoloDetector
        detector = TrtYoloDetector(args.engine)

    cap = cv2.VideoCapture(args.camera, cv2.CAP_V4L2)
    if not cap.isOpened():
        if detector is not None:
            detector.close()
        print(f"Failed to open /dev/video{args.camera}", file=sys.stderr)
        return 1
    cap.set(cv2.CAP_PROP_FOURCC, cv2.VideoWriter_fourcc(*"MJPG"))
    cap.set(cv2.CAP_PROP_FRAME_WIDTH, 1280)
    cap.set(cv2.CAP_PROP_FRAME_HEIGHT, 720)
    cap.set(cv2.CAP_PROP_FPS, 30)

    buf = FrameBuffer()
    stop_event = threading.Event()

    def capture_loop() -> None:
        while not stop_event.is_set():
            ret, frame = cap.read()
            if not ret or frame is None:
                continue
            if args.aruco:
                det = detect_aruco(frame)
                if det is not None:
                    mid, corners = det
                    pts = corners.astype(int).reshape((-1, 1, 2))
                    cv2.polylines(frame, [pts], True, (255, 200, 0), 2)
                    cv2.putText(frame, f"ArUco {mid}", tuple(pts[0][0]), cv2.FONT_HERSHEY_SIMPLEX, 0.7, (255, 200, 0), 2)
            if detector is not None:
                for d in detector.detect(frame, conf_thresh=args.conf):
                    x, y, w, h = [int(v) for v in d["bbox_xywh"]]
                    cv2.rectangle(frame, (x, y), (x + w, y + h), (0, 255, 0), 2)
                    cv2.putText(frame, f"{d['label']} {d['confidence']:.2f}", (x, max(20, y - 8)),
                                cv2.FONT_HERSHEY_SIMPLEX, 0.7, (0, 255, 0), 2)
            ok, enc = cv2.imencode(".jpg", frame, [int(cv2.IMWRITE_JPEG_QUALITY), 80])
            if ok:
                buf.set(enc.tobytes())

    worker = threading.Thread(target=capture_loop, daemon=True)
    worker.start()

    server = ThreadingHTTPServer(("0.0.0.0", args.port), make_handler(buf))
    print(f"MJPEG stream ready at http://jetson:{args.port}/stream.mjpg (or / in browser)", flush=True)
    try:
        server.serve_forever()
    except KeyboardInterrupt:
        pass
    finally:
        stop_event.set()
        server.server_close()
        cap.release()
        if detector is not None:
            detector.close()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
