#!/usr/bin/env python3
"""
Fine-tune YOLOv8n on AORUS5 (RTX 3070) and export static 640x640 ONNX for Jetson TensorRT.

Run with: ~/.venvs/drone-yolo/bin/python vision/train_yolo.py --data path/to/data.yaml
"""

import argparse
import os
import shutil

from ultralytics import YOLO


def main() -> int:
    parser = argparse.ArgumentParser(description="Fine-tune YOLOv8n and export ONNX for TensorRT")
    parser.add_argument("--data", required=True, help="path to YOLO data.yaml")
    parser.add_argument("--weights", default="yolov8n.pt", help="base weights (default: yolov8n.pt)")
    parser.add_argument("--epochs", type=int, default=60, help="training epochs")
    parser.add_argument("--imgsz", type=int, default=640, help="input image size")
    parser.add_argument("--batch", type=int, default=16, help="batch size")
    parser.add_argument("--device", default="0", help="CUDA device or cpu")
    parser.add_argument("--out-onnx", default="syma_x23w.onnx", help="output ONNX path")
    args = parser.parse_args()

    model = YOLO(args.weights)
    results = model.train(
        data=args.data,
        epochs=args.epochs,
        imgsz=args.imgsz,
        batch=args.batch,
        device=args.device,
    )
    best_weights = os.path.join(str(results.save_dir), "weights", "best.pt")
    trained = YOLO(best_weights)
    exported = trained.export(format="onnx", imgsz=args.imgsz, opset=17, simplify=False)
    shutil.copy2(str(exported), args.out_onnx)
    print(f"Best weights: {best_weights}")
    print(f"Exported ONNX: {args.out_onnx}")
    print(
        f"Next on Jetson: /usr/src/tensorrt/bin/trtexec --onnx={os.path.basename(args.out_onnx)} "
        "--saveEngine=syma_x23w.engine --fp16"
    )
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
