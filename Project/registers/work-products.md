# Work Product register

* **WP type:** Register
* **Version:** 0.5.0
* **Date:** 2026-11-13
* **Status:** Initial inventory; imported material is not revalidated.

## Research

| WP | Type | State / boundary |
|---|---|---|
| [Drone docking research concept](../../Research/problems/drone-docking.md) | Problem | Imported CapSense-focused draft; no longer represents CapSense as the project's primary technology. |
| [Original cabinet concept memo](../../Research/sources/author/original-cabinet-concept.md) | Source | Imported source copy; see migration register for origin. |
| [Infineon tools and technology notes](../../Research/sources/author/infineon-tools.md) | Source | Imported author notes; unverified technology candidates and limitations. |
| [Docking concept demonstration v5](../../Research/demonstrations/docking-v5/README.md) | Demonstration | Source scripts, Blender scene, and illustrative render; not hardware evidence. |
| [Infineon technology report](../../Research/studies/sensor-options/report.md) | Report / Analysis | Candidate mapping and caveats; not a procurement decision. |
| [CapSense and HMI report](../../Research/studies/capsense-hmi/report.md) | Report / Analysis | Background study; CapSense is an optional sensor/HMI line. |
| [Infineon technology map](../../Research/studies/sensor-options/infineon-technology-map.md) | Analysis Report | Imported engineering-oriented map retained as research input. |
| Supporting research notes | Report / Analysis support | Preserved under each study's `notes/` directory; not independently revalidated. |

## Engineering

| WP | Type | State / boundary |
|---|---|---|
| [Specifications index](../../Engineering/specifications/README.md) | Guide | Index for imported working specifications. |
| [System architecture](../../Engineering/specifications/system-architecture.md) | Specification | Draft; reflects source context and assumptions. |
| [Sensing and docking](../../Engineering/specifications/sensing-and-docking.md) | Specification | Draft; CapSense and other modalities remain hypotheses. |
| [Network and control](../../Engineering/specifications/network-and-control.md) | Specification | Draft; X23W protocol is explicitly unverified. |
| [Verification plan](../../Engineering/validation/verification-plan.md) | Plan | Proposed sequence, not an approved test authorization. |
| [Jetson/X23W bench report](../../Engineering/validation/jetson-x23w-bench-report.md) | Report / Experiment | Verified Jetson J501 `8188eu` Wi-Fi, X23W `IC_GPRTPB` video/motor bridge, retrained `syma_x23w.engine` (`0/60` FP), flight trim drift tests, and engineering analysis of platform limitations for fine closed-loop control. |
| [Pre-flight checklist](../../Engineering/validation/pre-flight-checklist.md) | Plan | Safety and axis-polarity checks before installing propellers for hover/landing tests. |
| [Syma X23W control utility](../../Engineering/software/syma-x23w-control/README.md) | Source Code | Go CLI and Bubble Tea controller/pseudo-drone plus read-only X23W `probe`/`video`; flight commands remain disabled. See its [tester guide](../../Engineering/software/syma-x23w-control/TESTING.md). |

## Project

| WP | Type | State / boundary |
|---|---|---|
| [Hardware and procurement register](hardware-procurement.md) | Register | Candidate list only; no purchase approved. |
| Project policies, guides, and templates | Policy / Guide / Template | Initial project conventions; review as the repository evolves. |
| [Open questions](open-questions.md) | Register | Current project facts and unresolved decisions. |
| [Source migration register](source-migration.md) | Register | Provenance and migration disposition. |

The inventory is a navigation aid, not evidence that listed work is complete
or accepted. Update it when adding or retiring a primary Work Product.

## Version history

| Version | Date | Change |
|---|---|---|
| 0.1.0 | 2026-11-12 | Initial Work Product inventory. |
| 0.2.0 | 2026-11-12 | Add imported technology notes and concept demonstration. |
| 0.3.0 | 2026-11-12 | Add the offline Syma X23W control utility scaffold. |
| 0.4.0 | 2026-11-12 | Update the Syma utility entry for the local controller/emulator PoC. |
| 0.5.0 | 2026-11-13 | Note the Bubble Tea TUI and emulator dashboard. |
