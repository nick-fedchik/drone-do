# Pre-flight bench & hover checklist: Syma X23W + Jetson J501

* **WP type:** Plan
* **Version:** 0.1.0
* **Date:** 2026-11-15
* **Status:** Bench preparation; not an autonomous flight approval
* **Related WI:** [#3](https://github.com/nick-fedchik/drone-do/issues/3), [#4](https://github.com/nick-fedchik/drone-do/issues/4)

## 1. Before installing propellers (blades OFF)

1. Confirm the physical Syma RC transmitter is **powered OFF**.
2. Connect Jetson `wlx58044f20d142` to `FPV_WIFI__140e` (`192.168.29.1`, `never-default yes`).
3. Fix `Logitech C615` (`/dev/video0`) above/beside the open landing pad with `assets/aruco_0.png` clearly visible and no bright reflections covering the marker.
4. ~~Verify stop mode~~ Done 2026-11-16 (blades off, throttle 0.5 for 4 s, phone timer): motors spin only after a `0x10` arm pulse; idle `stop` mode spun at idle ~9 s; plain `zero` idle ramps down in 8 s, `--end-pulse land` 8 s, `--end-pulse stop` 4.5–5 s (pulse length 1 s or 4 s makes no difference). Default is now `--stop-mode zero` with an automatic 0.5 s arm pulse; use `--end-pulse stop`. Motors cannot be cut instantly over Wi-Fi: keep ~5 s ramp-down in mind for the landing logic.
5. Run `pad_drone_lander.py` with `--max-tilt 0.15 --loss-timeout 0.25 --preview-port 8081` and confirm in `lander_status.json` / VLC:
   - `pad=True`, `drone=True`, zero false boxes on the empty pad;
   - Moving the drone left/right/forward/back produces the expected sign of `roll` and `pitch` (add `--invert-roll` / `--invert-pitch` / `--invert-yaw` if reversed).

## 2. Controlled hover test (blades ON, hand/tether safety)

1. Clear at least a `1.5 m x 1.5 m` unobstructed zone around the pad.
2. Start with `--max-tilt 0.15` and `--hover-throttle 0.50` (`0x80`).
3. Keep one hand ready on `Ctrl-C` in the `dronectl bridge` terminal; `Ctrl-C` sends 5x stop frames and `JHCMD\xd0\x02`.
4. Test short 2–3 s hover alignment first; observe drift and axis sign before enabling descent (`--descent-throttle 0.35`).
