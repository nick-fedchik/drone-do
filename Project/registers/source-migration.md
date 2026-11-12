# Source migration register

* **WP type:** Register
* **Version:** 0.3.0
* **Date:** 2026-11-12
* **Status:** Completed document import; content review pending

## Scope and provenance

Markdown documents were copied from
`nick-fedchyk/drafts/copter-capsense` at repository commit
`32d179e513072ba7c08e4f9664fdb3c9f3e99c1b`. The source tree was not modified.
Each imported document carries its WP type, origin, and a notice that its
technical content has not been revalidated. Internal Markdown links were
remapped when a matching imported document existed.

The source `README.md` was not imported. The destination repository's root
`README.md` was written as new project navigation; the source README and the
rest of the source tree were left unchanged.
The available v5 Blender scene, its two source scripts, and combined
concept-render video are also copied as an illustrative Demonstration Work
Product. Historical videos named in the source README but absent from the
source tree were not copied. Generated Python bytecode caches are excluded.
The source README, source tree, and other originals remain unchanged.

## Imported Work Products

| Source path | Target Work Product | Type |
|---|---|---|
| `copter-capsense/research-concept.md` | [Research problem](../../Research/problems/drone-docking.md) | Problem |
| `copter-capsense/drone-docking-cabinet-concept.md` | [Original concept memo](../../Research/sources/author/original-cabinet-concept.md) | Source |
| `copter-capsense/hardware-procurement.md` | [Hardware register](hardware-procurement.md) | Register |
| `copter-capsense/specs/README.md` | [Specifications index](../../Engineering/specifications/README.md) | Guide |
| `copter-capsense/specs/system-architecture.md` | [System architecture](../../Engineering/specifications/system-architecture.md) | Specification |
| `copter-capsense/specs/sensing-and-docking.md` | [Sensing and docking](../../Engineering/specifications/sensing-and-docking.md) | Specification |
| `copter-capsense/specs/network-and-control.md` | [Network and control](../../Engineering/specifications/network-and-control.md) | Specification |
| `copter-capsense/specs/infineon-options.md` | [Infineon technology map](../../Research/studies/sensor-options/infineon-technology-map.md) | Analysis Report |
| `copter-capsense/specs/verification-plan.md` | [Verification plan](../../Engineering/validation/verification-plan.md) | Plan |
| `copter-capsense/reports/Технології Infineon для дока.md` | [Infineon technology report](../../Research/studies/sensor-options/report.md) | Report / Analysis |
| `copter-capsense/reports/Застосування CapSense поза авто.md` | [CapSense and HMI report](../../Research/studies/capsense-hmi/report.md) | Report / Analysis |
| `copter-capsense/Infineon-tools.txt` | [Infineon tools and technology notes](../../Research/sources/author/infineon-tools.md) | Source |
| `copter-capsense/research_notes/Технології Infineon для дока/*.md` | `Research/studies/sensor-options/notes/` | Supporting research notes |
| `copter-capsense/research_notes/Застосування CapSense поза авто/*.md` | `Research/studies/capsense-hmi/notes/` | Supporting research notes |
| `copter-capsense/demo/generate_docking_demo.py` | `Research/demonstrations/docking-v5/generate_docking_demo.py` | Source Code |
| `copter-capsense/demo/generate_gate_intro.py` | `Research/demonstrations/docking-v5/generate_gate_intro.py` | Source Code |
| `copter-capsense/demo/drone-docking-demo-v5-gate-intro.blend` | `Research/demonstrations/docking-v5/drone-docking-demo-v5-gate-intro.blend` | Model |
| `copter-capsense/demo/drone-docking-demo-v5-gate-and-docking.mp4` | `Research/demonstrations/docking-v5/drone-docking-demo-v5-gate-and-docking.mp4` | Demonstration |

## Migration limits

- Import preserved the source statements and conclusions; link changes and
  WP metadata were added for repository navigation.
- No independent fact-check, scientific review, hardware test, or approval
  was performed by the migration.
- The CapSense-centered framing is historical input. The project's current
  direction and open decisions are recorded in
  [open questions](open-questions.md).
- The source repository and its original files remain authoritative for
  comparison with the imported copies.

## Version history

| Version | Date | Change |
|---|---|---|
| 0.1.0 | 2026-11-12 | First document migration inventory. |
| 0.2.0 | 2026-11-12 | Clarify that the source README was not imported or modified. |
| 0.3.0 | 2026-11-12 | Include the Infineon note and available v5 concept-demo assets. |
