import unittest

import numpy as np

from pad_drone_lander import FailSafeTracker, compute_relative_control


class TestPadDroneLander(unittest.TestCase):
    def test_steers_drone_toward_pad_and_clamps_tilt(self):
        pad = np.array([[900, 300], [1000, 300], [1000, 400], [900, 400]], dtype=np.float32)
        drone_bbox = [200.0, 300.0, 100.0, 100.0]
        ctrl, err = compute_relative_control(pad, drone_bbox, 1280, 720, max_tilt=0.35)
        self.assertIsNotNone(err)
        self.assertEqual(ctrl["roll"], 0.35)
        self.assertFalse(ctrl["aligned"])
        self.assertGreater(ctrl["throttle"], 0.0)

    def test_missing_drone_or_pad_returns_neutral(self):
        pad = np.array([[900, 300], [1000, 300], [1000, 400], [900, 400]], dtype=np.float32)
        ctrl, err = compute_relative_control(pad, None, 1280, 720)
        self.assertIsNone(err)
        self.assertEqual(ctrl["throttle"], 0.0)
        self.assertEqual(ctrl["roll"], 0.0)

    def test_failsafe_holds_brief_dropout_then_cuts_throttle(self):
        pad = np.array([[600, 320], [680, 320], [680, 400], [600, 400]], dtype=np.float32)
        drone_bbox = [590.0, 310.0, 100.0, 100.0]
        tracker = FailSafeTracker(loss_timeout=0.5, hover_throttle=0.45)

        c1 = tracker.step(1.0, pad, drone_bbox, 1280, 720)
        self.assertGreater(c1["throttle"], 0.0)

        c2 = tracker.step(1.2, pad, None, 1280, 720)
        self.assertEqual(c2["throttle"], 0.45)
        self.assertEqual(c2["roll"], 0.0)

        c3 = tracker.step(1.8, pad, None, 1280, 720)
        self.assertEqual(c3["throttle"], 0.0)

    def test_axis_inversion_and_descent_after_align_streak(self):
        pad = np.array([[600, 320], [680, 320], [680, 400], [600, 400]], dtype=np.float32)
        aligned_drone = [590.0, 310.0, 100.0, 100.0]
        tracker = FailSafeTracker(loss_timeout=0.25, hover_throttle=0.5, descent_throttle=0.35, align_frames=3)

        c1 = tracker.step(1.0, pad, aligned_drone, 1280, 720)
        c2 = tracker.step(1.1, pad, aligned_drone, 1280, 720)
        self.assertEqual(c1["phase"], "ALIGNING")
        self.assertEqual(c2["throttle"], 0.5)

        c3 = tracker.step(1.2, pad, aligned_drone, 1280, 720)
        self.assertEqual(c3["phase"], "DESCENDING")
        self.assertEqual(c3["throttle"], 0.35)

        right_pad = np.array([[900, 300], [1000, 300], [1000, 400], [900, 400]], dtype=np.float32)
        left_drone = [200.0, 300.0, 100.0, 100.0]
        inv, _ = compute_relative_control(right_pad, left_drone, 1280, 720, max_tilt=0.35, invert_roll=True)
        self.assertEqual(inv["roll"], -0.35)


if __name__ == "__main__":
    unittest.main()
