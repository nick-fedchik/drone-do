# Report: Jetson J501 Wi-Fi, Syma X23W protocol, and TensorRT bench verification

* **WP type:** Report / Experiment
* **Version:** 0.1.0
* **Date:** 2026-11-14
* **Status:** Bench-verified for Jetson Wi-Fi, X23W video/query, and TensorRT inference; flight commands unverified
* **Related WI:** [#1](https://github.com/nick-fedchik/drone-do/issues/1), [#3](https://github.com/nick-fedchik/drone-do/issues/3)
* **Source / baseline:** `Engineering/software/syma-x23w-control/` (`dronectl`), `Seeed-Studio/Linux_for_Tegra` (`r36.4.3`), `aircrack-ng/rtl8188eus`, `aivenlau/SyMaDemo_android`
* **Author / reviewer:** Session bench verification (AORUS5 + `mfedchyk@jetson`); unreviewed

## Question and scope

Verify the practical hardware, Wi-Fi link, camera protocol, and on-device vision pipeline on the Jetson compute node with the live Syma X23W aircraft on the bench:

- Confirm the physical Jetson carrier/image identity and add a working 2.4 GHz USB Wi-Fi interface (`TP-Link TL-WN722N v2/v3`).
- Verify the Syma X23W Wi-Fi AP, Joyhonest `IC_GPRTPB` handshake, read-only firmware query, and MJPEG video stream from Go (`dronectl`).
- Build and benchmark an FP16 TensorRT object-detection engine (`yolov8n.engine`) on Jetson Orin using captured X23W frames.
- Assess whether `kuzmin-machine` (KPOL / PYNQ-Z1) or `Pico SC0915` is needed for the flight/CV loop.

Excluded: live motor/flight-command testing on UDP `25000` (propellers remained installed on the aircraft).

## Method and evidence

### 1. Jetson J501 and USB Wi-Fi driver (`8188eu`)

- Physical label confirmed by user: `reServer Industrial J501-1`, `s/n: 102991854250100123` (Seeed Studio SKU `102991854`).
- OS/kernel state verified via SSH: Ubuntu `22.04.5 LTS` (`aarch64`), running kernel `5.15.148-tegra` (`#1 SMP PREEMPT Wed May 6 16:01:13 UTC 2026`), custom Seeed image `mfi_reserver-agx-orin-64g-j501-gmsl-6.2-36.4.3-2026-05-06.tar.gz` (branch `r36.4.3`), DTB `/boot/dtb/kernel_tegra234-j501x-0000+p3701-0005-reserver-gmsl.dtb`. All `nvidia-l4t-kernel*` packages are on `apt` `hold` and must not be overwritten with vanilla NVIDIA packages.
- Because `/lib/modules/5.15.148-tegra/build` pointed to an unavailable CI build path, a module build tree was prepared at `/usr/src/kernel-jammy-src-r36.4.3` from NVIDIA R36.4.3 `public_sources.tbz2` (`kernel_src.tbz2`), overlaid with `Seeed-Studio/Linux_for_Tegra` (`r36.4.3`), `/proc/config.gz` (`CONFIG_LOCALVERSION="-tegra"`), and the installed `Module.symvers`.
- `TP-Link TL-WN722N v2/v3` (`lsusb`: `2357:010c`, Realtek `RTL8188EUS`) driver `8188eu.ko` was compiled from `aircrack-ng/rtl8188eus` (`CONFIG_PLATFORM_ARM64_RPI = y`), installed, and loaded (`wlx58044f20d142`, `phy#0`, managed mode, `txpower 13.00 dBm`).

### 2. Syma X23W Wi-Fi and Go (`dronectl`) verification

- Network scan on both Jetson (`wlx58044f20d142`) and AORUS5 detected open AP `FPV_WIFI__140e` (`FC:58:4A:5F:14:0E`, channel 1).
- Jetson connected with `ipv4.never-default yes` (preserving default route on `eno1`): Jetson IP `192.168.29.100/24`, drone IP `192.168.29.1` (`ping` ~2.4–4.9 ms, 0% loss, no open common TCP ports).
- Cross-compiled `dronectl` (`GOOS=linux GOARCH=arm64`) was executed on Jetson against `192.168.29.1`:
  - `dronectl probe` sent `JHCMD` handshake frames to UDP `20000` and the read-only `FDWN` version query (`Command.java`) to UDP `20001`.
  - `dronectl video` captured UDP `10900` payloads and reassembled JPEG frames (`image/jpeg` `DecodeConfig`).

### 3. TensorRT 10.3 FP16 inference on Jetson Orin

- Installed JetPack `6.2.1` components verified on Jetson: CUDA `12.6`, cuDNN `9.3`, TensorRT `10.3.0` (`trtexec`), OpenCV `4.8.0`, VPI `3.2.4`; added `cuda-python` (`12.6.x`) via `pip --user`.
- Exported `yolov8n.onnx` (`1x3x640x640`, opset 17, 12.2 MB) on AORUS5 and compiled `~/models/yolov8n.engine` (`8.8 MB`, FP16) on Jetson using `trtexec`.
- Ran `trtexec` benchmark and Python `tensorrt` + `cuda-python` + `cv2` inference over the 48 captured X23W frames.

### 4. Auxiliary hardware and `kuzmin-machine` assessment

- Connected `Raspberry Pi Pico SC0915` on AORUS5 enumerates as `1209:3020 Generic LogicAnalyzer` (`/dev/ttyACM0`); `sigrok-cli` and `pulseview` are installed on AORUS5.
- Reviewed `~/GitHub/nick-fedchik/kuzmin-machine` (`S2`, `Engineering/README.md`, `pynq-z1/README.md`).

## Results

| Check | Measured observation | Notes / caveats |
|---|---|---|
| `8188eu` USB Wi-Fi on Jetson | Interface `wlx58044f20d142` active; scans 2.4 GHz APs up to `-31 dBm` | Built out-of-tree against `r36.4.3` overlay |
| `dronectl probe` (`UDP 20000`) | `105`-byte `JHCMD` reply (`4a48434d442000...`) | Confirms `IC_GPRTPB` camera control port |
| `dronectl probe` (`UDP 20001`) | `31`-byte `FDWN` reply: `F.WD1B06-20220325-005` | Read-only query from `Command.java` verified |
| `dronectl video` (`UDP 10900`) | `645` UDP packets in `3 s`; `48` frames (`16 FPS`, `640x480`, `12–19 KB`/frame) | `jpeg.DecodeConfig` accepted 48 headers; `cv2.imread` reported `premature end of data segment` on 7 frames (full `jpeg.Decode` validation needed) |
| `trtexec --fp16` (`yolov8n.engine`) | `339.5 qps`, mean GPU compute `2.94 ms`, host latency `3.35 ms` | Engine size `8.8 MB` on AGX Orin 64GB |
| End-to-end Python TRT + `cv2` | `48` frames in `0.861 s` (`55.7 FPS` incl. JPEG decode + NMS) | `0` COCO detections > `0.35` on bench-facing frames |

## Conclusion and limitations

1. **Jetson + Go + TensorRT pipeline is viable:** Jetson AGX Orin connects directly to `FPV_WIFI__140e` (`192.168.29.1`) over `TL-WN722N`, decodes live `640x480 @ 16 FPS` MJPEG in Go, and runs FP16 YOLOv8n inference at `55.7 FPS` end-to-end (>3x headroom over camera frame rate).
2. **Flight control remains unverified:** `jh_wifi.cpp` indicates `IC_GPRTPB` flight commands use `0xA5 0x5A` framing on UDP `25000`, but this has **not** been tested on the aircraft (propellers remain installed).
3. **JPEG validation gap:** `dronectl video` should validate full frames with `jpeg.Decode` rather than `jpeg.DecodeConfig` so truncated UDP-reassembled frames are dropped before inference.
4. **Role of `Pico` and `kuzmin-machine` (PYNQ-Z1):**
   - `Pico SC0915` (`1209:3020`) is ready on AORUS5 for sniffing the physical RC transmitter or internal UART/SPI buses without spinning motors.
   - `kuzmin-machine` (KPOL / BRAM / CME on PYNQ-Z1, ~100 ns hardware pipeline) is **not** needed for `16 FPS` Wi-Fi flight/CV; keep it as a candidate only for later deterministic docking/charging safety interlocks.

## History

| Version | Date | Change |
|---|---|---|
| 0.1.0 | 2026-11-14 | Record Jetson J501 Wi-Fi driver build, live X23W `IC_GPRTPB` probe/video capture, and TensorRT YOLOv8n benchmark. |
