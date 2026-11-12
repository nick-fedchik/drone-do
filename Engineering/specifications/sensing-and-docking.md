# Sensing and docking specification

> **WP type:** Specification<br>
> **Status:** Imported draft; technical claims and conclusions have not been revalidated.<br>
> **Origin:** [`copter-capsense/specs/sensing-and-docking.md`](https://github.com/nick-fedchik/drafts/blob/32d179e513072ba7c08e4f9664fdb3c9f3e99c1b/copter-capsense/specs/sensing-and-docking.md) at `32d179e513072ba7c08e4f9664fdb3c9f3e99c1b`.<br>
> **Migration:** Content preserved; internal links remapped where possible.


**Status:** working draft, v0.1 — 2026-11-11
**Primary question:** what can CapSense contribute to detection and state estimation around an indoor docking bay?

## 1. Sensor roles

| Channel | Intended evidence | Explicit limitation |
|---|---|---|
| CapSense on PSoC | Near-field capacitive change; possible hand/object proximity and activation sequence across separate zones. | Does not identify a drone, measure airflow, provide guaranteed 3D position, or prove safe landing. Metal mesh and nearby ground can shield or distort the field. |
| XENSIV radar | Motion/range/angle evidence, depending on sensor and setup. | Not a documented drone-versus-hand or airflow-source classifier. Multipath and near-range behavior must be measured. |
| REAL3 ToF | Depth/geometry evidence for a visible, unobstructed scene. | A ready-to-use Infineon ToF evaluation module was not confirmed in the reviewed materials; optical limits and integration remain TBD. |
| RGB camera | Visual record and possible geometry/vision reference. | Requires line of sight and lighting; an ordinary camera does not directly measure depth or airflow. |
| Direct airflow sensor | Air speed/direction at one or more defined points in the bay. | Likely third-party; an absolute-pressure reading or motor current is not a substitute for calibrated flow measurement. |
| Fan feedback | Fan command, tachometer/RPM if available, and electrical state. | Fan state does not establish actual local airflow or identify external airflow. |
| CapSense HMI | Operator status/maintenance interaction on a panel. | Not the sole emergency stop, interlock, or flight-control input. |

## 2. Functional requirements

| ID | Requirement | Verification |
|---|---|---|
| SEN-01 | The CapSense experiment shall log time-stamped readings/events independently for each tested electrode zone. | Inspect synchronized logs and compare with marked mock-up positions. |
| SEN-02 | The baseline test shall include empty bay, approach, partial entry, full occupancy, departure, and a hand entering without mechanical contact. | Repeated video-labelled trials. |
| SEN-03 | The test shall record geometry, material, target orientation, gap, electrode construction, and environmental conditions for each run. | Test record completeness. |
| SEN-04 | Metal grid, fan wiring/state, and nearby power electronics shall be tested as potential sources of shielding, coupling, or electrical interference. | Repeat baseline under controlled configurations. |
| SEN-05 | Airflow-only trials shall use an independent direct airflow measurement; CapSense shall not be used to claim airflow detection. | Compare airflow sensor and CapSense logs. |
| SEN-06 | Any combined state estimator shall retain an `UNKNOWN/AMBIGUOUS` outcome when channels disagree or confidence is inadequate. | Inject missing/conflicting sensor data. |
| SEN-07 | CapSense-only events shall not trigger aircraft takeoff/landing, charging, magnet engagement, or safety-critical actuation. | Design review and test interlocks. |
| SEN-08 | A pass-through gate shall use separated, timestamped sensing zones; a single capacitance threshold crossing shall not be treated as a validated count or speed measurement. | Compare against labelled reference trials. |
| SEN-09 | Speed shall be estimated from the time difference between corresponding events at two known sensing planes, or from an independently validated motion sensor. | Compare with an optical/radar reference; report uncertainty. |
| SEN-10 | Gate performance shall be evaluated across the allowed flight path, object orientation, speed, motor state, and environmental conditions. | Controlled test matrix; numerical acceptance limits TBD. |

Numeric values for range, latency, repeatability, detection probability, and false-positive rate are **TBD**. The team must set them before testing rather than treating a demonstration as proof of performance.

### Pass-through sensing gate concept

CapSense could be used to build a thin, lightweight, non-contact sensing gate—circular or rectangular—through which a drone flies without stopping. The frame would carry electrically isolated electrode segments around the aperture, mounted on a nonconductive support. A change in one or more zones can indicate that an object entered the near-field sensing volume; a sequence of zone events may provide coarse crossing direction or location.

The gate should be treated as a **passage detector**, not as a physical barrier or a drone-identification system. One portal can plausibly count well-separated crossing events after calibration, but a single electrode's signal duration is not a reliable speed measurement: it also depends on the drone's size, path, orientation, and sensor range. For speed, use two sensing planes a known distance apart, timestamp their corresponding detections, and estimate `speed = distance / time difference`. The baseline and event-matching method must be selected experimentally. Alternatively, pair CapSense detection with a validated optical light curtain or radar speed measurement.

This concept is a hypothesis, not a demonstrated CapSense capability at a specified gate size or flight speed. The electric field is near-field and may be weak at the larger distances needed for a wide aperture. Signal strength can vary substantially with the drone's conductive parts, plastic body, battery, orientation, altitude, and lateral path. Humidity, nearby metal, people, fan wiring, and motor/ESC interference may shift the baseline or mask a crossing. Segmentation and differential measurements can help characterize these effects but do not eliminate the need for tests.

The first prototype should use a nonconductive frame and segmented electrodes, with no continuous metal mesh or close ground plane across the aperture. Test one and two sensing planes against a camera or other independent reference, including empty-gate passages, repeated passes, varying paths/orientations/speeds, motor-off versus guarded motor-on trials, and closely spaced crossings. Keep `CROSSING_DETECTED` separate from dock states such as `OCCUPIED`; a crossing event does not establish identity, safe passage, or landing.

## 3. Proposed state model

The following are observation labels, not flight-controller states:

- `CROSSING_DETECTED`: a validated sequence of gate-zone events indicates passage through the aperture; direction and speed are reported only if the corresponding configuration has been validated.
- `EMPTY`: no validated presence evidence in the dock volume.
- `APPROACH`: a target is changing position toward the bay; requires an independent spatial/motion cue.
- `PARTIAL_OCCUPANCY`: one or more local zones report a target, but the full docking region is not confirmed.
- `OCCUPIED`: validated evidence indicates a target within the defined dock volume.
- `FLOW_ONLY`: airflow is measured, but no target is independently detected.
- `AMBIGUOUS`: evidence is conflicting, outside validated bounds, or missing.

`OCCUPIED` does not mean “identified drone,” “landed,” “rotors stopped,” or “safe to charge.” These require separate evidence and interlocks.

## 4. Dock and airflow hypotheses

- The landing grid may be made denser and supported by four fans beneath it to create a controllable air cushion. This is an unvalidated mechanical/flow hypothesis.
- The experimental question is whether occupancy sensing can distinguish “object present with associated flow” from “external airflow with no object” and “hand intrusion without contact.”
- CapSense can contribute local proximity/occupancy evidence. It does not measure or classify the air source.
- A direct airflow sensor should measure flow at defined points above/below the grid. A hot/cold temperature signal may be auxiliary only; it cannot reliably distinguish a hair dryer from rotor wash.
- A manually moved inert mock-up can validate geometry and presence sensing, but cannot validate propeller wash, vibration, motor EMI, or flight behavior.

## 5. Charging, retention, and wire routing

These topics are recorded as future investigations, not PoC requirements:

1. **Contact charging:** detect and align electrical terminals, route leads, measure contact integrity, and interlock charging.
2. **Contactless charging:** determine whether an inductive/resonant architecture is practical at the required aircraft power and alignment tolerance.
3. **Electromagnetic retention:** quantify holding/release force, coil heating, power loss, and safe release on power failure.
4. **Cable routing:** establish a mechanical method that cannot snag rotors or interfere with docking.

Before designing these systems, specify battery chemistry/cell count, charging voltage/current/power, allowed temperature, alignment tolerance, connector/contact rating, retention load, and release behavior. No single CapSense event may authorize charging or retention.

## 6. HMI

- CapSense HMI is suitable for noncritical local status/maintenance actions, subject to testing with the real panel, gloves, moisture, contamination, and EMI.
- Sensor labels in the demo are illustrative. In the actual system, labels/indicators should correspond to measured and logged channel states.
- The HMI shall not be the only means of emergency stop or safe maintenance access. Any safety function needs a separate risk-based design.

## 7. Infineon candidates

- **PSoC CAPSENSE:** best first Infineon experiment for proximity and electrode tuning. Candidate evaluation boards include CY8CKIT-041S-MAX or CY8CKIT-040T; neither is a ready-made multi-zone landing grid.
- **XENSIV BGT60TR13C:** candidate 60 GHz FMCW radar experiment for motion/range/angle; the published minimum detection range is 0.2 m and must be considered in mounting.
- **REAL3 ToF:** potential depth channel, but current kit/module, optics, and software integration need confirmation before procurement.
- **MOTIX/iMOTION and OptiMOS:** possible fan/coil/power-stage components only after motor, voltage, current, and charging requirements are known.
- **DPS368:** absolute barometric pressure/temperature sensor, not a direct airflow-velocity or differential-pressure sensor.

See [Infineon technology mapping](../../Research/studies/sensor-options/infineon-technology-map.md) and the [research report](../../Research/studies/sensor-options/report.md).
