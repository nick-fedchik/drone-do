import os
import tempfile
import unittest

import cv2
import numpy as np

from aruco_lander import compute_control, detect_aruco, generate_marker_png


class TestArucoLander(unittest.TestCase):
    def test_detects_generated_marker_and_centers_offset(self):
        with tempfile.TemporaryDirectory() as tmp:
            marker_path = os.path.join(tmp, "marker0.png")
            generate_marker_png(marker_path, marker_id=0, size=200)
            marker = cv2.imread(marker_path)
            self.assertIsNotNone(marker)

            # Place marker on right half of a 1280x720 canvas.
            canvas = np.full((720, 1280, 3), 255, dtype=np.uint8)
            mh, mw = marker.shape[:2]
            canvas[220 : 220 + mh, 820 : 820 + mw] = marker

            det = detect_aruco(canvas, target_id=0)
            self.assertIsNotNone(det)
            marker_id, corners = det
            self.assertEqual(marker_id, 0)

            ctrl = compute_control(corners, 1280, 720)
            self.assertGreater(ctrl["roll"], 0.1)
            self.assertGreater(ctrl["throttle"], 0.0)
            self.assertLessEqual(ctrl["throttle"], 1.0)

    def test_missing_marker_returns_neutral(self):
        ctrl = compute_control(None, 1280, 720)
        self.assertEqual(ctrl["throttle"], 0.0)
        self.assertEqual(ctrl["roll"], 0.0)
        self.assertEqual(ctrl["pitch"], 0.0)
        self.assertEqual(ctrl["yaw"], 0.0)


if __name__ == "__main__":
    unittest.main()
