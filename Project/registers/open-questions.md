# Project facts and open questions

* **WP type:** Register
* **Version:** 0.2.0
* **Date:** 2026-11-12
* **Status:** Working register; not a design baseline

| Topic | Current record | Status / next evidence |
|---|---|---|
| Project focus | R&D for drone control, passage sensing, and docking. CapSense is not the primary or assumed sensing technology. | User direction; detailed scope and priorities remain to be baselined. |
| Initial drone-control setup | Local host identified as GIGABYTE AORUS 5 SE: Ubuntu 26.04.1 x86_64, Intel Core i7-12700H, 64 GB class memory, Intel Wi-Fi, RTX 3070 Laptop GPU. | AORUS5 Wi-Fi association/control path to X23W is not tested. |
| Jetson compute node | Фізичний лейбл (підтверджено користувачем 2026-11-14): `reServer Industrial J501-1`, `s/n: 102991854250100123` (SKU `102991854`). SSH to `mfedchyk@jetson` succeeded: `Ubuntu 22.04.5 aarch64`, `5.15.148-tegra`, кастомний образ Seeed `r36.4.3` (`kernel_tegra234-j501x-0000+p3701-0005-reserver-gmsl.dtb`); пакети `nvidia-l4t-kernel*` на `hold` (не знімати). 12 Cortex-A78AE CPUs, 64 GB class UMA, Ollama 0.40.7-rc0-6-ge61ba8b4 активна (`http://jetson:11434`, 14B/32B моделі). | Підключений `TP-Link TL-WN722N v2/v3` (`2357:010c`, RTL8188EUS) працює: драйвер `8188eu.ko` зібрано під підготовлене дерево `5.15.148-tegra` (`Seeed-Studio/Linux_for_Tegra` `r36.4.3`), інтерфейс `wlx58044f20d142` (`phy#0`) активний. Model quality, inference latency, and quoted TOPS claims are not independently validated. |
| Additional boards | PYNQ-Z1 and Raspberry Pi Pico SC0915 (RP2040) confirmed present by user. Pico is intended for a 24-channel logic-analyzer setup. | Identify board revisions and Pico analyzer firmware, sampling limits, and voltage compatibility before use. |
| Test aircraft status | Syma X23W is present and powers on (`FPV_WIFI__140e`, `FC:58:4A:5F:14:0E`, channel 1, open AP, IP `192.168.29.1`). Jetson `wlx58044f20d142` connects as `192.168.29.100/24`. | Keep propellers removed for bench protocol tests; verify battery charge before longer runs. |
| Drone protocol | Live test on Jetson (2026-11-14) with `dronectl probe` verified Joyhonest `IC_GPRTPB` (`192.168.29.1`): `JHCMD` over UDP `20000` (105-byte reply), live MJPEG (`JFIF`) stream on UDP `10900` (396 packets/2s), and `Command.java` `FDWN` query over UDP `20001` returning firmware version `F.WD1B06-20220325-005`. Flight command packet mapping (`UDP 25000` `0xA5 0x5A` framing) is not yet verified on this aircraft. | Verify flight-command packet format safely on the bench (propellers removed) before enabling live flight transport in Go. |
| Control utility | Offline-only Go CLI exists at `Engineering/software/syma-x23w-control/`; it validates and simulates normalized JSONL frames and transmits no commands. | Add live transport only after X23W-specific protocol evidence and a separate controlled test plan. |
| Preferred implementation language | Go is preferred for the control software; TinyGo is to be considered. | Compare actual target needs. Standard Go is the initial candidate for Linux laptop/OpenWrt; TinyGo may suit a later MCU target. |
| OpenWrt | A router is not currently available; adapting OpenWrt is a possible later stage. | Select hardware first, then verify radio, architecture, driver, package, and routing support. |
| Gate sensing | The target should not require a tag or modification to the drone. Counting and speed measurement are research goals. | Compare camera/stereo or depth tracking and radar. No method is selected or validated yet. |
| CapSense | Retained as an optional research modality, not the gate/control baseline. | Test only if it answers a defined question; do not infer drone identity or speed from capacitance alone. |
| Docking, charging, retention | These remain separate R&D topics. | Define measurable requirements and safety constraints before hardware selection. |

## Version history

| Version | Date | Change |
|---|---|---|
| 0.1.0 | 2026-11-12 | Initial project facts and open questions. |
| 0.2.0 | 2026-11-12 | Record locally verified host/Jetson and user-confirmed test hardware. |
| 0.3.0 | 2026-11-12 | Record the offline-first X23W utility boundary. |
