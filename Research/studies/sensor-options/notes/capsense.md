# Infineon CapSense for an indoor drone docking box

> **WP type:** Analysis Report — supporting research note<br>
> **Status:** Imported draft; technical claims and conclusions have not been revalidated.<br>
> **Origin:** [`copter-capsense/research_notes/Технології Infineon для дока/capsense.md`](https://github.com/nick-fedchik/drafts/blob/32d179e513072ba7c08e4f9664fdb3c9f3e99c1b/copter-capsense/research_notes/%D0%A2%D0%B5%D1%85%D0%BD%D0%BE%D0%BB%D0%BE%D0%B3%D1%96%D1%97%20Infineon%20%D0%B4%D0%BB%D1%8F%20%D0%B4%D0%BE%D0%BA%D0%B0/capsense.md) at `32d179e513072ba7c08e4f9664fdb3c9f3e99c1b`.<br>
> **Migration:** Content preserved; internal links remapped where possible.


Research snapshot: 2026-11-11. Scope: Infineon official product pages and official Infineon documentation/source repositories. “CapSense” below means the PSoC-integrated CAPSENSE™ system; Infineon’s separately marketed capacitive-sensor IC page is not assumed to describe PSoC CapSense.

## What hardware and guidance fit

### Takeaway
The most directly relevant current evaluation platforms found are **CY8CKIT-041S-MAX (PSoC 4100S Max)** and **CY8CKIT-040T (PSoC 4000T)**: Infineon’s own pages list a capacitive proximity sensor on each. They demonstrate a proximity sensor, not a prebuilt multi-zone docking-box array. For zones, a PSoC CapSense design with separate, individually scanned electrodes is the plausible direction; exact device channel/pin budget, layout, and scan timing still need validation.

### Cited Findings
- Infineon lists CY8CKIT-041S-MAX as a PSoC 4100S Max kit, describes a capacitive proximity sensor on its expansion board, and identifies 5th-generation CapSense. The page lists a proximity sensor alongside touchpad, slider, and buttons. [CY8CKIT-041S-MAX product page](https://www.infineon.com/evaluation-board/CY8CKIT-041S-MAX)
- Infineon lists CY8CKIT-040T as a PSoC 4000T kit with a capacitive proximity sensor, button, and touchpad; the page identifies it as an active/preferred product at research time. [CY8CKIT-040T product page](https://www.infineon.com/evaluation-board/CY8CKIT-040T)
- The official CAPSENSE Middleware README identifies supported sensing methods (self-capacitance CSD, mutual-capacitance CSX, inductive ISX) and widgets including buttons, matrix buttons, sliders, touchpads, and proximity sensors. It links the current ModusToolbox CapSense Configurator and Tuner guides. [Infineon CAPSENSE Middleware README](https://github.com/Infineon/capsense/blob/master/README.md)
- The middleware README links **CAPSENSE™ Design Guide** (AN85951, PSoC 4/PSoC 6) and PSoC 4 documentation. These are the appropriate official design references for electrode/layout and device-specific configuration decisions; use the Configurator/Tuner guides for setting up and inspecting a concrete design. [AN85951 CAPSENSE Design Guide](https://www.infineon.com/row/public/documents/30/42/infineon-an85951-psoc-4-psoc-6-capsense-design-guide-applicationnotes-en.pdf) · [ModusToolbox CAPSENSE Configurator User Guide](https://www.infineon.com/ModusToolboxCapSenseConfig) · [ModusToolbox CAPSENSE Tuner Tool Guide](https://www.infineon.com/ModusToolboxCapSenseTuner) · [PSoC 4 documentation](https://documentation.infineon.com/psoc4/)
- CY8CKIT-145-40XX is an active/preferred PSoC 4000S prototyping kit with three CapSense buttons and a five-segment slider; its official page does not list a proximity electrode. CY8CKIT-149 has buttons and a six-segment slider, likewise not a listed proximity sensor. They are relevant for basic touch/slider evaluation, not the closest out-of-box proximity demonstration. [CY8CKIT-145-40XX product page](https://www.infineon.com/evaluation-board/CY8CKIT-145-40XX) · [CY8CKIT-149 product page](https://www.infineon.com/evaluation-board/CY8CKIT-149)
- Infineon’s current general “Capacitive sensors” category describes high-sensitivity capacitive ICs, I/Q demodulation, active shielding, and multi-channel sensing for automotive-facing use cases. This is a separate product framing from the PSoC CapSense MCU platform; the page does not establish a suitable docking-box development kit. [Infineon Capacitive sensors](https://www.infineon.com/products/sensor/capacitive-sensors)
- I did not establish from a current official product page a current MBR3 development kit with multi-zone proximity capability. Do not assume MBR3 is equivalent to a programmable PSoC multi-electrode implementation or select it without verifying a specific current device, documentation, and evaluation hardware.

### Inferences
- **Evaluation approach, not a purchase commitment:** begin with one proximity-capable PSoC kit (CY8CKIT-041S-MAX or CY8CKIT-040T) to validate the enclosure/overlay and target-distance signal. The 041S-MAX is a stronger candidate if the later design will need more custom firmware/peripheral flexibility; that is an engineering preference, not an Infineon claim that its included electrode is a multi-zone array.
- A segmented grid should be prototyped as separate self-capacitance electrodes (for example, a few pads around/beneath the landing area), each assigned/scanned independently if the selected device’s pin/channel resources permit. Treat CapSense widgets/slider segments as configurable signal inputs, not evidence that a specific kit provides spatially independent proximity zones.
- The dedicated CY8CKIT-040T/041S-MAX proximity electrode is the supported starting point; moving to custom electrodes may need a custom PCB/flex and suitable connections. Confirm these connections and supported sensor counts from the selected kit user guide, device datasheet, and Configurator before committing.

### Gaps
- Official material reviewed did not give a guaranteed detection distance for the proposed box, target, or overlay, nor prove that either kit’s stock proximity pad can be repurposed as an array.
- MBR3’s current product/kit status and exact proximity configuration were not verified from current official documentation in this research; it is deliberately not recommended.

## What controls proximity feasibility

### Takeaway
CapSense measures changes in the electrode’s capacitive coupling, not distance as an independent geometric quantity. Feasibility depends on electrode and target geometry, the dielectric overlay, parasitic capacitance and nearby ground/shielding, plus scan/tuning choices. There is no source-backed universal “drone detection range” to quote for this enclosure.

### Cited Findings
- Infineon’s CAPSENSE Middleware describes proximity as a supported widget and provides self-capacitance (CSD) and mutual-capacitance (CSX) sensing methods. [CAPSENSE Middleware README](https://github.com/Infineon/capsense/blob/master/README.md)
- Infineon points designers to AN85951, **CAPSENSE™ Design Guide**, for sensor design, and the Configurator/Tuner guides for configuring and tuning designs. [AN85951 CAPSENSE Design Guide](https://www.infineon.com/row/public/documents/30/42/infineon-an85951-psoc-4-psoc-6-capsense-design-guide-applicationnotes-en.pdf) · [Configurator User Guide](https://www.infineon.com/ModusToolboxCapSenseConfig) · [Tuner Tool Guide](https://www.infineon.com/ModusToolboxCapSenseTuner)
- Infineon’s current capacitive-sensor category lists active shielding and wide measurement range as features of its sensor IC offerings. This is not a quantified promise for a PSoC CapSense electrode or this box. [Infineon Capacitive sensors](https://www.infineon.com/products/sensor/capacitive-sensors)

### Inferences
- **Electrode/target:** larger or differently shaped electrodes and greater conductive target area can change coupling and useful sensing distance; an airframe’s conductive parts, battery/PCB, frame material, orientation, and grounding can therefore change the response. A carbon-fiber frame, metals, and conductive coatings may behave differently from a plastic mock-up. These effects require measurements in the actual stack-up.
- **Overlay and spacing:** a nonconductive box wall or grid coating can be sensed through, but its thickness/material and the spacing between electrode and target alter coupling. A grounded metal grid or close ground plane may shield or reshape the field; nearby fan metalwork and wiring are part of the electrode environment, not negligible background.
- **Shielding and ground:** grounded copper behind/around an electrode can screen the field; a driven shield may be used in suitable designs to shape/extend the field, but its implementation is device/layout-specific. Avoid promising that a shield always improves range; consult AN85951 and validate with the actual ground/grid/fan configuration.
- **Sampling:** scan rate/latency depends on the selected sensor configuration, sensor count, resolution, filters/averaging, and firmware schedule. More zones and aggressive smoothing can affect response time. Choose an application response target, then measure scan period and false-trigger behavior on the intended MCU rather than inferring it from a kit’s marketing page.
- For a practical first experiment, fix the electrode/overlay geometry, log tuned/raw sensor values and thresholds while stepping conductive targets through known locations/distances, then repeat with fan wiring, grounded structures, and the final box assembled. Test each zone separately before combining logic.

### Gaps
- The source material cited here does not provide dimensions, an electrode-to-target distance, or scan-rate values for the proposed landing grid. Those require a specified device, electrode stack-up, configurator settings, and measured prototype.

## Which states it can distinguish (and what requires rotor tests)

### Takeaway
CapSense can plausibly report **near-field occupancy/proximity changes by zone**. By itself it cannot identify “drone” versus “hand” semantically or infer whether airflow came from drone propellers or an external blower. A manually moved inert mock-up is useful for geometry/occupancy and false-trigger tests, but not for validating real rotor airflow, vibration, or motor/ESC interference.

### Cited Findings
- Infineon documents CapSense proximity widgets and promotes capacitive sensing for touch/proximity-type interactions; the documented measurement is capacitive sensing, not an airflow or semantic object-classification modality. [CAPSENSE Middleware README](https://github.com/Infineon/capsense/blob/master/README.md) · [Infineon Capacitive sensors](https://www.infineon.com/products/sensor/capacitive-sensors)
- CY8CKIT-041S-MAX and CY8CKIT-040T product pages document proximity sensors, while their listed kit features do not claim drone/hand classification or airflow-source discrimination. [CY8CKIT-041S-MAX](https://www.infineon.com/evaluation-board/CY8CKIT-041S-MAX) · [CY8CKIT-040T](https://www.infineon.com/evaluation-board/CY8CKIT-040T)

### Inferences
- **(a) Drone present with propeller airflow:** CapSense may detect the drone body/nearby conductive structures as occupancy, whether rotors are stopped or running. It will not directly measure the propeller airflow. If body occupancy is the intended trigger, test sensor response with actual drone geometry (or a representative conductive mock-up) and the drone’s position/orientation.
- **(b) External airflow, e.g. hair dryer:** airflow itself does not change capacitance in a way that identifies its source. Heating, humidity/condensation, EMI, or movement of a target could affect readings or create artifacts, but these are confounds, not a reliable “external airflow” signature.
- **(c) Hand insertion:** a hand’s presence can produce a proximity signal without mechanical contact. That supports a non-contact intrusion/occupancy trigger, not a guaranteed identity label. A hand and a drone body may produce overlapping signals depending on size, material, grounding, distance, and position.
- **Claims to avoid:** “CapSense detects a drone,” “distinguishes a human hand from a drone,” “detects propeller airflow,” or “determines the airflow source.” A defensible claim is “CapSense detects a tuned capacitive-proximity event in the tested region.”
- **Mock-up PoC is suited to:** verify basic hand/no-touch response, sensor zone coverage, approach/departure thresholds, enclosure/overlay effects, blind spots, and whether airflow alone causes false events with fans off/on. A manually moved inert mock-up can map body-shaped conductive targets and position sensitivity, but cannot simulate a real flight controller, rotor wash, vibration, or motor electronics.
- **Needs real rotor testing:** distinguish/characterize operational effects of real propellers, motors, ESCs, vibration, fan/rotor airflow interaction, and electrical noise. Use a protected, controlled rig for this phase. If source classification is essential, combine proximity with an independent modality such as the reported camera/Jetson, fan/rotor telemetry, or an airflow/pressure sensor; the fusion algorithm also needs labelled real-test data.
- The reported Jetson J501 64G and Logitech UB webcam could support a separate vision/ground-truth or sensor-fusion study; neither changes CapSense’s sensing modality. The later OpenWRT Wi‑Fi 6 router/local-navigation work is system integration, not evidence of proximity feasibility.

### Gaps
- No official Infineon source found claims object classification, airflow sensing, or drone-propeller test validation for CapSense. The proposed hand/drone overlap and mock-up limits are engineering inferences and must be empirically tested.
- The user’s exact drone dimensions/material, box/grid construction, electrode placement, and desired warning distance are unspecified; classification performance cannot be estimated without them.
