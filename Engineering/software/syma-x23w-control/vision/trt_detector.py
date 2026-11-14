#!/usr/bin/env python3
"""
TensorRT 10.3 YOLOv8n FP16 detector for Jetson (C615 camera or static image).
"""

import argparse
import json
import sys
import time
from typing import List, Optional

import cv2
import numpy as np
import tensorrt as trt
from cuda import cudart

COCO_CLASSES = [
    "person", "bicycle", "car", "motorcycle", "airplane", "bus", "train", "truck", "boat",
    "traffic light", "fire hydrant", "stop sign", "parking meter", "bench", "bird", "cat",
    "dog", "horse", "sheep", "cow", "elephant", "bear", "zebra", "giraffe", "backpack",
    "umbrella", "handbag", "tie", "suitcase", "frisbee", "skis", "snowboard", "sports ball",
    "kite", "baseball bat", "baseball glove", "skateboard", "surfboard", "tennis racket",
    "bottle", "wine glass", "cup", "fork", "knife", "spoon", "bowl", "banana", "apple",
    "sandwich", "orange", "broccoli", "carrot", "hot dog", "pizza", "donut", "cake",
    "chair", "couch", "potted plant", "bed", "dining table", "toilet", "tv", "laptop",
    "mouse", "remote", "keyboard", "cell phone", "microwave", "oven", "toaster", "sink",
    "refrigerator", "book", "clock", "vase", "scissors", "teddy bear", "hair drier", "toothbrush",
]


class TrtYoloDetector:
    def __init__(self, engine_path: str, class_names: Optional[List[str]] = None) -> None:
        logger = trt.Logger(trt.Logger.WARNING)
        with open(engine_path, "rb") as f:
            runtime = trt.Runtime(logger)
            self.engine = runtime.deserialize_cuda_engine(f.read())
        if self.engine is None:
            raise RuntimeError(f"failed to load TensorRT engine: {engine_path}")

        self.context = self.engine.create_execution_context()
        self.in_name = self.engine.get_tensor_name(0)
        self.out_name = self.engine.get_tensor_name(1)
        self.in_shape = tuple(self.engine.get_tensor_shape(self.in_name))
        self.out_shape = tuple(self.engine.get_tensor_shape(self.out_name))

        num_classes = self.out_shape[1] - 4
        if class_names:
            self.class_names = class_names
        elif num_classes == 1:
            self.class_names = ["syma_x23w"]
        elif num_classes == len(COCO_CLASSES):
            self.class_names = COCO_CLASSES
        else:
            self.class_names = [f"class_{i}" for i in range(num_classes)]

        self.in_bytes = int(np.prod(self.in_shape) * np.dtype(np.float32).itemsize)
        self.out_bytes = int(np.prod(self.out_shape) * np.dtype(np.float32).itemsize)
        err, self.d_in = cudart.cudaMalloc(self.in_bytes)
        if err != 0:
            raise RuntimeError(f"cudaMalloc input failed: {err}")
        err, self.d_out = cudart.cudaMalloc(self.out_bytes)
        if err != 0:
            cudart.cudaFree(self.d_in)
            raise RuntimeError(f"cudaMalloc output failed: {err}")

        self.context.set_tensor_address(self.in_name, self.d_in)
        self.context.set_tensor_address(self.out_name, self.d_out)
        err, self.stream = cudart.cudaStreamCreate()
        if err != 0:
            self.close()
            raise RuntimeError(f"cudaStreamCreate failed: {err}")
        self.h_out = np.empty(self.out_shape, dtype=np.float32)

    def detect(self, img: np.ndarray, conf_thresh: float = 0.35, iou_thresh: float = 0.45) -> list:
        img_h, img_w = img.shape[:2]
        in_h, in_w = self.in_shape[2], self.in_shape[3]
        scale = min(in_w / img_w, in_h / img_h)
        nw, nh = int(img_w * scale), int(img_h * scale)
        resized = cv2.resize(img, (nw, nh))
        canvas = np.full((in_h, in_w, 3), 114, dtype=np.uint8)
        dx, dy = (in_w - nw) // 2, (in_h - nh) // 2
        canvas[dy : dy + nh, dx : dx + nw] = resized

        blob = np.ascontiguousarray(canvas[:, :, ::-1].transpose(2, 0, 1), dtype=np.float32) / 255.0
        cudart.cudaMemcpyAsync(
            self.d_in, blob.ctypes.data, self.in_bytes, cudart.cudaMemcpyKind.cudaMemcpyHostToDevice, self.stream
        )
        self.context.execute_async_v3(self.stream)
        cudart.cudaMemcpyAsync(
            self.h_out.ctypes.data, self.d_out, self.out_bytes, cudart.cudaMemcpyKind.cudaMemcpyDeviceToHost, self.stream
        )
        cudart.cudaStreamSynchronize(self.stream)

        preds = self.h_out[0].T
        scores = preds[:, 4:].max(axis=1)
        keep = scores >= conf_thresh
        preds, scores = preds[keep], scores[keep]
        if len(preds) == 0:
            return []

        cls_ids = preds[:, 4:].argmax(axis=1)
        boxes = []
        for cx, cy, bw, bh in preds[:, :4]:
            boxes.append([
                float((cx - bw / 2 - dx) / scale),
                float((cy - bh / 2 - dy) / scale),
                float(bw / scale),
                float(bh / scale),
            ])

        indices = cv2.dnn.NMSBoxes(boxes, scores.tolist(), conf_thresh, iou_thresh)
        detections = []
        if len(indices) > 0:
            for idx in indices.flatten():
                cid = int(cls_ids[idx])
                label = self.class_names[cid] if cid < len(self.class_names) else f"class_{cid}"
                detections.append({
                    "label": label,
                    "confidence": round(float(scores[idx]), 4),
                    "bbox_xywh": [round(v, 1) for v in boxes[idx]],
                })
        return detections

    def close(self) -> None:
        if getattr(self, "stream", None) is not None:
            cudart.cudaStreamDestroy(self.stream)
            self.stream = None
        if getattr(self, "d_in", None) is not None:
            cudart.cudaFree(self.d_in)
            self.d_in = None
        if getattr(self, "d_out", None) is not None:
            cudart.cudaFree(self.d_out)
            self.d_out = None


def main() -> int:
    parser = argparse.ArgumentParser(description="Run YOLOv8n TensorRT detection on Jetson")
    parser.add_argument("--engine", default="/home/mfedchyk/models/yolov8n.engine", help="TensorRT .engine path")
    parser.add_argument("--camera", type=int, default=0, help="V4L2 camera index (default: 0)")
    parser.add_argument("--image", default="", help="optional image path instead of live camera")
    parser.add_argument("--frames", type=int, default=10, help="frames to process from camera (0 = loop)")
    parser.add_argument("--conf", type=float, default=0.35, help="confidence threshold")
    args = parser.parse_args()

    detector = TrtYoloDetector(args.engine)
    try:
        if args.image:
            img = cv2.imread(args.image)
            if img is None:
                print(f"Failed to read image: {args.image}", file=sys.stderr)
                return 1
            t0 = time.monotonic()
            dets = detector.detect(img, conf_thresh=args.conf)
            ms = (time.monotonic() - t0) * 1000.0
            print(json.dumps({"latency_ms": round(ms, 2), "detections": dets}))
            return 0

        cap = cv2.VideoCapture(args.camera, cv2.CAP_V4L2)
        if not cap.isOpened():
            print(f"Failed to open /dev/video{args.camera}", file=sys.stderr)
            return 1
        cap.set(cv2.CAP_PROP_FOURCC, cv2.VideoWriter_fourcc(*"MJPG"))
        cap.set(cv2.CAP_PROP_FRAME_WIDTH, 1280)
        cap.set(cv2.CAP_PROP_FRAME_HEIGHT, 720)
        cap.set(cv2.CAP_PROP_FPS, 30)
        try:
            count = 0
            while args.frames == 0 or count < args.frames:
                ret, frame = cap.read()
                if not ret or frame is None:
                    continue
                count += 1
                t0 = time.monotonic()
                dets = detector.detect(frame, conf_thresh=args.conf)
                ms = (time.monotonic() - t0) * 1000.0
                print(json.dumps({"frame": count, "latency_ms": round(ms, 2), "detections": dets}))
        finally:
            cap.release()
    finally:
        detector.close()
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
