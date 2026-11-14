#!/usr/bin/env python3
"""
Capture images from a V4L2 camera (Logitech C615 on Jetson) for YOLO dataset labeling.
"""

import argparse
import os
import sys
import time

import cv2


def main() -> int:
    parser = argparse.ArgumentParser(description="Capture dataset frames from V4L2 camera")
    parser.add_argument("--camera", type=int, default=0, help="V4L2 camera index (default: 0)")
    parser.add_argument("--out", default="dataset/images", help="output directory for JPEGs")
    parser.add_argument("--interval", type=float, default=0.5, help="seconds between saved frames")
    parser.add_argument("--count", type=int, default=200, help="number of frames to capture (0 = until Ctrl-C)")
    args = parser.parse_args()

    os.makedirs(args.out, exist_ok=True)
    cap = cv2.VideoCapture(args.camera, cv2.CAP_V4L2)
    if not cap.isOpened():
        print(f"Failed to open /dev/video{args.camera}", file=sys.stderr)
        return 1
    cap.set(cv2.CAP_PROP_FOURCC, cv2.VideoWriter_fourcc(*"MJPG"))
    cap.set(cv2.CAP_PROP_FRAME_WIDTH, 1280)
    cap.set(cv2.CAP_PROP_FRAME_HEIGHT, 720)
    cap.set(cv2.CAP_PROP_FPS, 30)

    saved = 0
    next_save = time.monotonic()
    try:
        while args.count == 0 or saved < args.count:
            ret, frame = cap.read()
            if not ret or frame is None:
                continue
            now = time.monotonic()
            if now >= next_save:
                saved += 1
                filename = os.path.join(args.out, f"x23w_{int(time.time())}_{saved:04d}.jpg")
                cv2.imwrite(filename, frame)
                print(f"saved {saved}: {filename}")
                next_save = now + max(0.05, args.interval)
    except KeyboardInterrupt:
        pass
    finally:
        cap.release()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
