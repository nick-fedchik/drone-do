# System architecture specification

> **WP type:** Specification<br>
> **Status:** Imported draft; technical claims and conclusions have not been revalidated.<br>
> **Origin:** [`copter-capsense/specs/system-architecture.md`](https://github.com/nick-fedchik/drafts/blob/32d179e513072ba7c08e4f9664fdb3c9f3e99c1b/copter-capsense/specs/system-architecture.md) at `32d179e513072ba7c08e4f9664fdb3c9f3e99c1b`.<br>
> **Migration:** Content preserved; internal links remapped where possible.


**Status:** working draft, v0.1 — 2026-11-11
**Scope:** indoor drone docking research with CapSense as the primary research focus.

## 1. Objective

Evaluate whether spatially segmented CapSense electrodes can provide useful, repeatable evidence of an object approaching, entering, occupying, and leaving a drone docking bay. Determine whether that evidence adds value when combined with other sensing modalities.

The first PoC is a stationary, single-bay test with an inert mock-up. It does not demonstrate autonomous flight, drone identification, safe landing, charging, or operational reliability.

## 2. System boundary

The concept has six separable subsystems:

1. **Dock structure:** an indoor bay with an entry, dense landing grid, and four proposed fans below the grid.
2. **Presence and geometry sensing:** CapSense electrodes; optional radar, ToF, and visible-light camera.
3. **Airflow measurement:** a direct airflow instrument, likely third-party, plus state/telemetry from the dock fans if available.
4. **Compute and networking:** J501 reServer for higher-level processing; a separate OpenWrt-capable Wi-Fi node is a candidate radio bridge.
5. **Aircraft control:** a flight controller and a Wi-Fi control interface for the Syma X23W toy copter as a possible early test target. Its exact protocol is not yet verified.
6. **Power, charging, and retention:** separate future investigations for contact charging, inductive charging, terminal alignment, wire routing, and electromagnetic retention.

The demo video is a visual concept only; it is not evidence that these subsystems have been built or tested.

## 3. Reported setup and status

| Item | Status | Specification note |
|---|---|---|
| Three TP-Link Deco X60 units | Confirmed by user | Existing home mesh connected to the primary PON service. |
| TP-Link Deco X20-4G | Confirmed by user | Spare LTE unit intended for Internet backup if the PON link fails. It is not assigned to drone control. The actual failover wiring and behavior remain to be verified. |
| J501 reServer | Confirmed by user | No built-in Wi-Fi; Ethernet-only operation. A dedicated second Ethernet interface is the preferred candidate for the drone radio link. Exact NIC names and capabilities are TBD. |
| J501 physical state | Confirmed by user | Board is open, without an enclosure. User reports the M.2 Key E 2230 position is occupied by an NVMe drive. |
| OpenWrt router | Candidate | Exact model, hardware revision, radio chipset, OpenWrt release, and STA/forwarding support are not selected. |
| Radio-node priority | Confirmed by user | Detachable antennas and maximum engineering flexibility are more important than Wi-Fi 6 or mesh features. Router model remains unselected. |
| Syma X23W | Candidate test aircraft | The user is interested in controlling it over its Wi-Fi AP from their own software. Exact aircraft/camera-board revision and command protocol are TBD. |
| Jetson module/RAM configuration | Unverified | J501 is confirmed; the exact compute-module SKU and memory configuration need confirmation before performance assumptions. |
| Two PYNQ-Z1 boards | Unverified | Appeared in a proposed architecture text, but are not confirmed as available project equipment. |

## 4. Proposed logical architecture

```text
CapSense electrodes ─┐
Radar / ToF / camera ├──> local sensing and logging ──> state estimate: occupied / empty / unknown
Airflow sensor ──────┤
Dock fan telemetry ──┘

J501 reServer ── dedicated Ethernet ── OpenWrt candidate ── Wi-Fi STA ── drone AP
      │
      └── separate Ethernet ── Deco X60 home network ── PON / Internet

Deco X20-4G: separate, existing LTE-backup role; not in the drone control path.
```

The diagram is a target topology, not a completed network configuration. In particular, return routing for video, firewall policy, address allocation, and any access from home-network devices must be designed and tested.

## 5. System requirements

| ID | Requirement | Status / verification |
|---|---|---|
| SYS-01 | The initial experiment shall use one stationary bay and an inert mock-up; multi-bay assignment is out of the first PoC scope. | In scope; verify test setup. |
| SYS-02 | The system shall distinguish measured sensor events from inferred aircraft identity or safety state. | Required design rule; review state logic. |
| SYS-03 | When sensing channels conflict, are unavailable, or fall outside validated conditions, the reported state shall be `UNKNOWN/AMBIGUOUS`, not a confident docked/clear decision. | Required design rule; test fault cases. |
| SYS-04 | All performance limits and thresholds shall be set before the corresponding experiment. No numerical detection, latency, or false-alarm target is approved yet. | TBD with research team. |
| SYS-05 | The four-fan landing-cushion idea shall remain a mechanical/airflow hypothesis until measured with a defined fan and bay geometry. | Candidate; airflow test needed. |
| SYS-06 | Charging and electromagnetic retention shall remain separate future workstreams until the battery, voltage, current, mechanical loads, and release behavior are specified. | Deferred. |
| SYS-07 | The Full HD Blender demo shall be labelled concept art/visualization and shall not be presented as a validated prototype. | `demo/drone-docking-demo-v3-FullHD.mp4`. |

## 6. Safety and architecture constraints

- CapSense, Wi-Fi RSSI, an OpenWrt process, and Jetson software are not safety-rated flight-control functions.
- The aircraft flight controller (or a suitable real-time onboard controller) must own stabilization and loss-of-command behavior. The exact Syma X23W behavior on Wi-Fi loss must be tested, not assumed.
- No sensor alone proves that the aircraft is identified, landed, stopped, mechanically retained, or safe to charge.
- Initial radio/control tests shall be conducted with propellers removed. Any later rotor or flight test requires a guarded, controlled test setup and a separate safety review.
- Charging, fan actuation, magnet actuation, and flight control must not be enabled merely because one CapSense channel changes state.

## 7. Open decisions

- Final bay dimensions, material, grid geometry, and fan placement.
- Whether a custom flight platform is needed after the toy-aircraft networking experiment.
- Exact J501 module, Ethernet interfaces, and operating system.
- Exact OpenWrt device and supported driver/modes.
- Whether PYNQ-Z1 is actually part of the available test bench.
- Required performance thresholds and evaluation criteria.
- Battery chemistry/cell count, charging power, contacts/coil design, and retention force.

## 8. Related material

- [Research concept](../../Research/problems/drone-docking.md)
- [Hardware and procurement register](../../Project/registers/hardware-procurement.md)
- [Original concept notes](../../Research/sources/author/original-cabinet-concept.md)
- [Infineon research report](../../Research/studies/sensor-options/report.md)
- [CapSense/HMI research report](../../Research/studies/capsense-hmi/report.md)
