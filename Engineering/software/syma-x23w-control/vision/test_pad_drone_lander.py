import unittest

import numpy as np

from pad_drone_lander import compute_relative_control


class TestPadDroneLander(unittest.TestCase):
    def test_steers_drone_toward_pad_on_right(self):
        pad = np.array([[900, 300], [1000, 300], [1000, 400], [900, 400]], dtype=np.float32)
        drone_bbox = [200.0, 300.0, 100.0, 100.0]  # center at (250, 350)
        ctrl = compute_relative_control(pad, drone_bbox, 1280, 720)
        self.assertGreater(ctrl["roll"], 0.3)
        self.assertFalse(ctrl["aligned"])
        self.assertGreater(ctrl["throttle"], 0.0)

    def test_missing_drone_or_pad_returns_neutral(self):
        pad = np.array([[900, 300], [1000, 300], [1000, 400], [900, 400]], dtype=np.float32)
        ctrl = compute_relative_control(pad, None, 1280, 720)
        self.assertEqual(ctrl["throttle"], 0.0)
        self.assertEqual(ctrl["roll"], 0.0)


if __name__ == "__main__":
    unittest.main()
