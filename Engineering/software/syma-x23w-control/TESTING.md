# Tester guide: local drone controller and pseudo-emulator

* **WP type:** Test Instructions
* **Version:** 0.2.0
* **Date:** 2026-11-13
* **Status:** Local software test only; not an aircraft test

## Purpose and limits

This procedure verifies local IPC between `dronectl` and `droneemu`. The
pseudo-emulator stores the most recently accepted normalized control frame
and counts commands. It does not simulate flight, sensors, battery, radio,
or the Syma X23W protocol. No real aircraft or radio link is involved.

## Prerequisites

- Go 1.24 or later
- Two interactive terminals for the TUI test
- Run commands from `Engineering/software/syma-x23w-control/`

The default socket is
`<user-cache-directory>/drone-do/syma-x23w.sock`. Both programs must run as
the same user. An alternate path can be supplied with `--socket PATH`.

## Automated checks

Run:

```sh
go test ./...
go vet ./...
```

Expected: both commands exit with status 0.

Confirm the aircraft boundary:

```sh
go run ./cmd/dronectl protocol status
```

Expected output includes `X23W protocol: UNVERIFIED` and
`Live radio transport: DISABLED`.

## CLI IPC smoke test

1. In terminal A, start the pseudo-drone:

   ```sh
   go run ./cmd/droneemu serve
   ```

   Expected: it reports the local Unix socket path and remains running.

2. In terminal B, query its initial status:

   ```sh
   go run ./cmd/dronectl status
   ```

   Expected JSON has `"running": true`, `"commands_received": 0`, and no
   `"last_frame"`.

3. Send one abstract control frame:

   ```sh
   go run ./cmd/dronectl send --throttle 0.25 --roll -0.5 --pitch 0.75 --yaw 0.1
   ```

   Expected: the command succeeds and the acknowledgement reports one
   received command with the values above. Terminal A also prints an event.

4. Query status again:

   ```sh
   go run ./cmd/dronectl status
   ```

   Expected: `commands_received` is `1`, and `last_frame` contains the sent
   values.

5. Verify validation rejects an out-of-range value:

   ```sh
   go run ./cmd/dronectl send --throttle 1.5
   ```

   Expected: a nonzero exit status and an error explaining that throttle must
   be in `[0,1]`. The emulator's command count must remain `1`.

Stop terminal A with Ctrl-C.

## Interactive TUI test

1. In terminal A, run:

   ```sh
   go run ./cmd/droneemu tui
   ```

   Expected: a Bubble Tea dashboard shows emulator/IPC status, an explicitly
   mocked fixed 100% battery, normalized control meters, and recent commands.
   Aircraft protocol and physical specifications are marked unknown or
   unverified. The emulator TUI uses `q` to stop.

2. In terminal B, run:

   ```sh
   go run ./cmd/dronectl tui
   ```

   Expected: the Bubble Tea controller screen identifies itself as
   local-emulator-only and displays the current normalized controls.

3. Press `r` once. The throttle should increase by `0.05`; terminal A should
   show an incremented command count and the updated frame.

4. Press `a`, `w`, and `q` once each. Roll, pitch, and yaw should update by
   `0.1` in their indicated directions. Values must remain within their
   documented normalized ranges.

5. Press `0`. All four controls should become zero. Press `x` to exit; the
   controller sends a final neutral frame before closing.

6. Confirm terminal A shows the final neutral frame, then press `q` to stop
   the emulator.

## ArUco landing tracker (Jetson + Logitech C615 -> `droneemu`)

1. Run the unit test on Jetson (uses OpenCV 4.8 `cv2.aruco`):

   ```sh
   python3 -m unittest discover -s vision -p "test_*.py"
   ```

2. Generate a printable `DICT_4X4_50` marker (ID `0`):

   ```sh
   python3 vision/aruco_lander.py --generate-marker /tmp/aruco_0.png
   ```

3. Start `droneemu tui` in terminal A and the tracker in terminal B:

   ```sh
   python3 vision/aruco_lander.py --camera 0 --hz 10
   ```

   Pointing the camera at the printed/displayed marker updates `roll`, `pitch`,
   `yaw`, and `throttle` in the `droneemu` dashboard over the local Unix socket.

## YOLOv8n training (AORUS5) and TensorRT inference (Jetson)

1. On **Jetson**, capture dataset frames from Logitech C615:

   ```sh
   python3 vision/capture_dataset.py --camera 0 --interval 0.5 --count 200 --out ~/dataset/images
   ```

2. Label the `syma_x23w` class in YOLO format and fine-tune on **AORUS5** (RTX 3070):

   ```sh
   ~/.venvs/drone-yolo/bin/python vision/train_yolo.py --data /path/to/data.yaml --out-onnx syma_x23w.onnx
   ```

3. On **Jetson**, compile the FP16 TensorRT engine and run `trt_detector.py`:

   ```sh
   /usr/src/tensorrt/bin/trtexec --onnx=syma_x23w.onnx --saveEngine=$HOME/models/syma_x23w.engine --fp16
   python3 vision/trt_detector.py --engine $HOME/models/syma_x23w.engine --camera 0 --frames 10
   ```

## Troubleshooting

- If the client reports that it cannot connect, confirm the emulator is still
  running and both programs use the same `--socket` path and user.
- If startup reports that the socket path already exists, first verify no
  emulator is running at that path. The program intentionally does not delete
  existing files or stale sockets automatically.
- The TUI requires a real interactive terminal. Use `serve` and the CLI
  commands in non-interactive environments.
- Any result from this guide is evidence only for local software IPC; it is
  not evidence about aircraft protocol compatibility, safe flight, or radio
  behavior.
