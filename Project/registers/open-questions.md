# Project facts and open questions

* **WP type:** Register
* **Version:** 0.1.0
* **Date:** 2026-11-12
* **Status:** Working register; not a design baseline

| Topic | Current record | Status / next evidence |
|---|---|---|
| Project focus | R&D for drone control, passage sensing, and docking. CapSense is not the primary or assumed sensing technology. | User direction; detailed scope and priorities remain to be baselined. |
| Initial drone-control setup | The AORUS4 laptop has Wi-Fi and is intended for same-room experiments with the Syma X23W. | Confirm actual Wi-Fi band and connection behavior on the laptop/drone. |
| Drone protocol | The exact X23W board/firmware and command protocol are unknown. Prior X22W examples are not proof for X23W. | Capture/verify behavior using the user's own aircraft before implementing control packets. |
| Preferred implementation language | Go is preferred for the control software; TinyGo is to be considered. | Compare actual target needs. Standard Go is the initial candidate for Linux laptop/OpenWrt; TinyGo may suit a later MCU target. |
| OpenWrt | A router is not currently available; adapting OpenWrt is a possible later stage. | Select hardware first, then verify radio, architecture, driver, package, and routing support. |
| Gate sensing | The target should not require a tag or modification to the drone. Counting and speed measurement are research goals. | Compare camera/stereo or depth tracking and radar. No method is selected or validated yet. |
| CapSense | Retained as an optional research modality, not the gate/control baseline. | Test only if it answers a defined question; do not infer drone identity or speed from capacitance alone. |
| Docking, charging, retention | These remain separate R&D topics. | Define measurable requirements and safety constraints before hardware selection. |

## Version history

| Version | Date | Change |
|---|---|---|
| 0.1.0 | 2026-11-12 | Initial project facts and open questions. |
