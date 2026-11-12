# Verification and experiment plan

> **WP type:** Plan<br>
> **Status:** Imported draft; technical claims and conclusions have not been revalidated.<br>
> **Origin:** [`copter-capsense/specs/verification-plan.md`](https://github.com/nick-fedchik/drafts/blob/32d179e513072ba7c08e4f9664fdb3c9f3e99c1b/copter-capsense/specs/verification-plan.md) at `32d179e513072ba7c08e4f9664fdb3c9f3e99c1b`.<br>
> **Migration:** Content preserved; internal links remapped where possible.


**Status:** working draft, v0.1 — 2026-11-11
**Purpose:** define what evidence is needed before progressing from concept to integrated testing.

## 1. General rules

- Record hardware revision, firmware/software version, configuration, geometry, environment, and trial ID for every run.
- Synchronize sensor logs with video and actuator/fan state to the extent practical.
- Set numerical thresholds for detection rate, false alarms, latency, range, repeatability, and loss-of-link response with the research team **before** the relevant test.
- Preserve raw readings as well as derived states. Do not report only a success/failure label.
- Keep “unknown/ambiguous” as a valid outcome when channels disagree or the sample is outside the tested envelope.
- A visual demo, product datasheet, or one successful trial is not a performance validation.

## 2. Staged plan

### Stage A — inventory and interface confirmation

Record:

- J501 model/module, operating system, Ethernet NIC names and capabilities.
- Deco X60/X20-4G topology and the actual LTE failover behavior.
- Candidate OpenWrt model, revision, chipset, supported OpenWrt image, STA mode, and wired routing.
- Syma X23W revision, Wi-Fi SSID, DHCP address/gateway, control/video ports, and app version.
- CapSense kit/device, electrodes, wiring, grid, enclosure materials, and supply/grounding.
- Whether PYNQ-Z1 boards are actually available and what specific experiment they would support.

**Exit:** no unverified model/port/protocol assumptions remain hidden in the test setup.

### Stage B — CapSense geometry and baseline

Test an empty bay and a manually moved inert mock-up across marked positions. Include approach, partial entry, occupancy, withdrawal, and hand intrusion. Compare one electrode/threshold against multiple zones. Repeat with the grid and final overlay in place.

**Record:** raw/tuned counts, zone events, timestamps, mock-up position/orientation, distance/gap, repeat number, video reference, environmental conditions.

**Exit:** a documented operating envelope and initial repeatability/false-event results exist. Numeric acceptance limits are TBD until agreed.

### Stage C — airflow and actuator effects

Use a direct airflow instrument at defined points. Compare no flow, built-in fans off/on, external airflow only (including hot and cool sources), and mock-up plus airflow. Separately log fan command/tach/current if available.

An inert mock-up cannot validate rotor wash. The experiment may establish that airflow-only events can be measured and whether they perturb CapSense; it cannot prove separation of propeller wash from a dryer.

**Exit:** flow readings, spatial placement, sensor response time, and fan-only baseline are documented.

### Stage D — independent spatial sensor

If the sensing question warrants it, evaluate radar (e.g. BGT60TR13C) for movement/range and/or a working ToF module for depth geometry. Use the same trajectories and synchronized logs as Stage B.

**Exit:** sensor contribution over CapSense/camera baseline is measured; do not claim object-classification accuracy unless labelled trials support it.

### Stage E — X23W network/protocol validation

1. Confirm OpenWrt device/firmware, STA association, and the actual path from J501's dedicated Ethernet interface.
2. Capture the user's X23W Wi-Fi traffic using the official application in a controlled setup.
3. Verify command IP, ports, local source port, initialization, packet length, field offsets, checksums, and video direction for the exact aircraft revision.
4. Compare with the public X22W implementation only as a hypothesis.
5. Test the neutral/low-authority path with propellers removed; document actual behavior when packets stop.

**Exit:** the X23W packet format and link-loss response are evidenced for the user's unit. No X22W-only assumption is carried into control software.

### Stage F — controlled aircraft integration

This stage is blocked until the control interface, operator procedure, test enclosure, flight controller behavior, and safety review are defined. Use a guarded area and independent shutdown. Keep high-level navigation separate from stabilization and the flight controller's failsafe.

**Exit:** thresholds and acceptance criteria approved by the research team and measured across repeated, labelled trials.

### Stage G — charging and electromagnetic retention

After defining battery, electrical, thermal, mechanical, and release requirements, test contact charging, inductive charging, terminal alignment, wire routing, and magnet retention as separate subsystems. Include loss of power and fault cases.

**Exit:** separate validated subsystem specifications and safety review; no charging or magnet operation is triggered by CapSense alone.

## 3. Evidence and acceptance metrics

For each applicable stage, report:

- Detection probability and false-positive/false-negative counts by scenario.
- Measured detection distance and spatial coverage by position/orientation.
- Event latency and repeatability.
- Sensitivity to enclosure/grid/fan/charger states, EMI, temperature, humidity, and drift.
- Sensor disagreement and `UNKNOWN/AMBIGUOUS` frequency.
- For radio control: packet loss, command timing, video behavior, reconnect behavior, and measured loss-of-command response.

Numerical pass/fail thresholds are **TBD** and must be agreed before testing. The initial research result may be “not suitable” or “inconclusive.”

## 4. Out of scope for the initial PoC

- Autonomous flight or landing with a toy drone.
- Proving a CapSense-based drone identity classifier.
- Proving that an external fan/hair dryer can be distinguished from propeller wash without direct flow data and rotor trials.
- Charging, electromagnetic retention, or automatic wiring.
- Multi-bay scheduling and fleet operation.
- Safety certification or production readiness.
