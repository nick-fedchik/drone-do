#!/usr/bin/env python3
"""
Minimal single-class YOLO box labeler.

Drag a rectangle around the drone, then:
  Enter/Space  save the box and go to the next image
  n            no drone in this image (writes an empty label file = negative)
  r            redraw the box
  b            go back one image
  q / Esc      quit (progress is kept; re-run to continue)

Run with: ~/.venvs/drone-yolo/bin/python vision/label_boxes.py --images DIR
Labels are written next to the images as <name>.txt (class 0).
"""

import argparse
import glob
import os

import cv2


def main() -> int:
    parser = argparse.ArgumentParser(description="Single-class YOLO box labeler")
    parser.add_argument("--images", required=True, help="directory with .jpg images")
    parser.add_argument("--max-width", type=int, default=1280, help="display width")
    args = parser.parse_args()

    paths = sorted(glob.glob(os.path.join(args.images, "*.jpg")))
    if not paths:
        print("no images found")
        return 1

    state = {"start": None, "end": None, "drag": False}

    def on_mouse(event, x, y, flags, param):
        if event == cv2.EVENT_LBUTTONDOWN:
            state.update(start=(x, y), end=(x, y), drag=True)
        elif event == cv2.EVENT_MOUSEMOVE and state["drag"]:
            state["end"] = (x, y)
        elif event == cv2.EVENT_LBUTTONUP:
            state.update(end=(x, y), drag=False)

    win = "label: drag box, Enter=save, n=no drone, r=redraw, b=back, q=quit"
    cv2.namedWindow(win)
    cv2.setMouseCallback(win, on_mouse)

    i = 0
    while i < len(paths):
        label_path = os.path.splitext(paths[i])[0] + ".txt"
        if os.path.exists(label_path) and not state.get("revisit"):
            i += 1
            continue
        state.update(start=None, end=None, drag=False, revisit=False)
        img = cv2.imread(paths[i])
        h, w = img.shape[:2]
        scale = min(1.0, args.max_width / w)
        shown = cv2.resize(img, (int(w * scale), int(h * scale))) if scale < 1.0 else img

        while True:
            canvas = shown.copy()
            if state["start"] and state["end"]:
                cv2.rectangle(canvas, state["start"], state["end"], (0, 255, 0), 2)
            cv2.putText(canvas, f"{i + 1}/{len(paths)}", (10, 30), cv2.FONT_HERSHEY_SIMPLEX, 1, (0, 0, 255), 2)
            cv2.imshow(win, canvas)
            key = cv2.waitKey(20) & 0xFF
            if key in (ord("q"), 27):
                return 0
            if key == ord("r"):
                state.update(start=None, end=None)
            elif key == ord("b") and i > 0:
                i -= 1
                state["revisit"] = True
                break
            elif key == ord("n"):
                open(label_path, "w").close()
                i += 1
                break
            elif key in (13, 32) and state["start"] and state["end"]:
                x1, y1 = state["start"]
                x2, y2 = state["end"]
                x1, x2 = sorted((x1, x2))
                y1, y2 = sorted((y1, y2))
                if x2 - x1 < 5 or y2 - y1 < 5:
                    continue
                dw, dh = shown.shape[1], shown.shape[0]
                cx, cy = (x1 + x2) / 2 / dw, (y1 + y2) / 2 / dh
                bw, bh = (x2 - x1) / dw, (y2 - y1) / dh
                with open(label_path, "w") as f:
                    f.write(f"0 {cx:.6f} {cy:.6f} {bw:.6f} {bh:.6f}\n")
                i += 1
                break
    cv2.destroyAllWindows()
    print("done")
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
