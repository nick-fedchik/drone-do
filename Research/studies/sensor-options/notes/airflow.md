# Infineon sensors and electronics for a drone docking box

> **WP type:** Analysis Report — supporting research note<br>
> **Status:** Imported draft; technical claims and conclusions have not been revalidated.<br>
> **Origin:** [`copter-capsense/research_notes/Технології Infineon для дока/airflow.md`](https://github.com/nick-fedchik/drafts/blob/32d179e513072ba7c08e4f9664fdb3c9f3e99c1b/copter-capsense/research_notes/%D0%A2%D0%B5%D1%85%D0%BD%D0%BE%D0%BB%D0%BE%D0%B3%D1%96%D1%97%20Infineon%20%D0%B4%D0%BB%D1%8F%20%D0%B4%D0%BE%D0%BA%D0%B0/airflow.md) at `32d179e513072ba7c08e4f9664fdb3c9f3e99c1b`.<br>
> **Migration:** Content preserved; internal links remapped where possible.


Evidence checked for the requested status date: 2026-11-11.

## Which Infineon sensor families fit, and are they direct airflow or differential-pressure sensors?

### Takeaway
Infineon has useful absolute-pressure, temperature, magnetic speed/position, capacitive-proximity, and mmWave radar technologies, but the reviewed portfolio does not establish a directly suitable airflow-velocity or differential-pressure sensor. Pressure, temperature, magnetic, and capacitive signals are indirect or complementary measurements; none identifies the source of airflow alone.

### Cited Findings
- Infineon’s pressure portfolio describes XENSIV™ pressure sensors as pressure transducers and characterizes its air-pressure IC portfolio as absolute-pressure sensors, including barometric (BAP) and manifold (MAP) categories. The product page does not describe a differential-pressure or direct air-velocity sensor. — [Pressure sensors (Infineon)](https://www.infineon.com/products/sensor/pressure-sensors)
- XENSIV™ DPS368 is explicitly a digital **barometric air pressure (BAP)** sensor; its product page gives a 300–1200 hPa range, I²C/SPI interfaces, and says it measures pressure and temperature. Its named document is **“DPS368 Datasheet.”** — [DPS368 product page and datasheet](https://www.infineon.com/part/DPS368)
- Infineon’s magnetic speed portfolio describes non-contact linear/rotational speed measurement using Hall, GMR, and TMR technologies and targets magnetic encoders or ferromagnetic gear-tooth wheels. Listed families include TLE494x Hall, TLE504x GMR, and TLE554x TMR wheel-speed sensors. Its angle-sensor family reports magnetic-field orientation/rotational position. These need a suitable magnetic target; they do not sense air movement. — [Magnetic speed sensors](https://www.infineon.com/products/sensor/magnetic-speed-sensors); [Angle sensors](https://www.infineon.com/products/sensor/magnetic-position-sensors/angle-sensors)
- XENSIV™ 60 GHz radar sensors can provide motion, distance/range, speed, and (depending on configuration) angle/object information, and can detect objects through non-conductive materials. **BGT60TR13C** is a 60 GHz radar with antennas in package; its product page lists presence detection, range, motion, and speed. Named document: **“BGT60TR13C Datasheet.”** — [Radar sensors](https://www.infineon.com/products/sensor/radar-sensors); [BGT60TR13C product page and datasheet](https://www.infineon.com/part/BGT60TR13C)
- Infineon’s capacitive-sensor page describes detection of human touch/proximity and hands-on detection. This is a plausible short-range hand-presence channel, but it does not measure airflow or determine its source. — [Capacitive sensors](https://www.infineon.com/products/sensor/capacitive-sensors)

### Inferences
- A DPS368 or another absolute-pressure part could record pressure fluctuations at a location, but it is not a differential-pressure flow transducer. Inferring velocity would require a calibrated mechanical flow element/pressure arrangement and suitable signal processing; using two absolute sensors as a pair would be an experimental workaround, not a documented direct differential sensor.
- Radar is the strongest reviewed Infineon candidate for detecting a hand or drone entering without contact. Validate performance with the dock geometry and stationary as well as moving targets; radar presence alone does not establish that rotor wash is occurring.
- Magnetic speed/angle ICs are only relevant if a compatible magnet/encoder can be attached to a fan or, later, to a custom drone rotor. They are not general-purpose object-presence sensors.

### Gaps
- No directly suitable Infineon airflow-velocity/anemometer product or dedicated differential-pressure part for this dock was identified in the reviewed official portfolio pages. A third-party airflow sensor (e.g. hot-wire/thermal anemometer or a differential-pressure element with a characterized flow fixture) is likely needed for ground-truth flow.
- The DPS368 product description confirms pressure and temperature readings, but this review did not establish that its temperature channel is a calibrated measurement of incoming air temperature rather than a sensor/device temperature channel; do not treat it as an airflow thermometer without checking the datasheet and validating placement/thermal response.

## What Infineon temperature and presence sensors can contribute; can temperature identify a hair dryer?

### Takeaway
Temperature can flag a sufficiently hot incoming stream, but cannot distinguish a hair dryer from rotor wash as a source, and it fails when the external air is unheated or a dryer is on cool mode. Use it as an auxiliary feature, not a source classifier.

### Cited Findings
- Infineon offers a family of factory-calibrated digital temperature sensors with I²C output, low power, and compact packages. The category page is titled **“Temperature sensors.”** — [Temperature sensors (Infineon)](https://www.infineon.com/products/sensor/temperature-sensors)
- DPS368 is described as measuring pressure and temperature as well as being an absolute BAP sensor. — [DPS368 product page and datasheet](https://www.infineon.com/part/DPS368)
- The BGT60TR13C and radar portfolio support presence/motion and distance/range sensing; the capacitive portfolio supports touch/proximity sensing. — [BGT60TR13C product page and datasheet](https://www.infineon.com/part/BGT60TR13C); [Capacitive sensors](https://www.infineon.com/products/sensor/capacitive-sensors)

### Inferences
- A separate digital temperature sensor placed in the airflow could distinguish a hot stream from ambient-temperature air if its response time, thermal coupling, and location are suitable. However, the measured temperature does not uniquely identify a hair dryer: other warm sources can produce the same reading, and a cool/unheated external stream can resemble rotor wash thermally.
- Use a sensor array or multiple points (e.g., above and below the grid) and record transients as well as steady state. A single temperature sensor can lag, be warmed by the enclosure/electronics, or miss a narrow jet.
- Neither CapSense/capacitive proximity, pressure, nor temperature alone should be described as identifying the airflow source. Radar/capacitance help answer “is an object/hand present?”; flow instrumentation and controlled state context answer “is air moving/how?”

### Gaps
- The official temperature-family page describes generic product-family capabilities but does not, in the retrieved content, identify a particular Infineon temperature part/board validated for fast free-air measurements in this enclosure. Check a specific device datasheet and response-time/packaging details before selection.
- No source-specific signature (hair dryer versus drone rotor wash) is specified by the reviewed Infineon material.

## Which Infineon motor-control/power products could monitor box fans, and what does current not tell us?

### Takeaway
Infineon motor-control devices are relevant only when the box-fan motor type and electrical interface match, especially for a controllable low-voltage BLDC fan. Motor current can be useful telemetry, but current alone does not prove fan RPM, airflow, or the presence/source of external air.

### Cited Findings
- Infineon’s motor-driver portfolio covers low-voltage MOTIX™ solutions (12–60 V) and iMOTION™ integrated control for high-voltage three-phase motors (below 600 V), including gate drivers, integrated drivers, and power modules. — [Motor drivers (Infineon)](https://www.infineon.com/products/power/motor-drivers)
- **MOTIX™ 6EDL7141** is a three-phase gate driver for battery-supplied BLDC/PMSM control; its product data lists PWM control and three current-sense amplifiers (CSA). Its named document is **“6EDL7141 Datasheet.”** — [6EDL7141 product page and datasheet](https://www.infineon.com/part/6EDL7141)
- **MOTIX™ TLE9879QXA40** is a single-chip three-phase motor-control device integrating an Arm Cortex-M3, analog/mixed-signal peripherals, LIN, and driver stages; the product page says it fits BLDC applications such as pumps and fans. Named document: **“TLE9879QXA40 Datasheet.”** — [TLE9879QXA40 product page and datasheet](https://www.infineon.com/part/TLE9879QXA40)

### Inferences
- If the below-grid fans are low-voltage BLDC and the controller is accessible, a suitable MOTIX controller/driver could control them and expose motor electrical telemetry; a dedicated tach output or magnetic speed target is a stronger RPM measurement. If the installed fans are mains AC or an enclosed proprietary assembly, these low-voltage BLDC parts are not a drop-in monitoring solution; identify motor topology, voltage, and interface first.
- Current trends can indicate a motor is energized or that its electrical load changed, but load, supply voltage, commutation/control mode, and fan mechanics all affect current. A current change cannot reliably establish a specific RPM, airflow speed/direction, blocked/failed airflow, drone presence, or external hair-dryer flow.
- Knowing the built-in fan command/RPM helps subtract its baseline from airflow observations, but it does not distinguish an external flow source by itself.

### Gaps
- The box-fan motor voltage/type, available tach/feedback line, control access, and current waveform are unknown. Product fit cannot be finalized until those are measured/identified.
- The listed products are motor-control components, not airflow or source-classification sensors.

## What sensor combination and experiment can separate flow-only, object-only, and object-plus-rotor-flow?

### Takeaway
Use independent presence and flow channels, with fan-state telemetry as context. Infineon radar (optionally capacitive proximity) can detect the non-contact object; a direct third-party flow sensor should establish airflow. Classify source only from controlled, measured combinations—not from one Infineon sensor.

### Cited Findings
- Infineon radar supports non-contact object/motion/range detection, while its capacitive sensors support touch/proximity; neither page describes an airflow-source classifier. — [Radar sensors](https://www.infineon.com/products/sensor/radar-sensors); [Capacitive sensors](https://www.infineon.com/products/sensor/capacitive-sensors)
- XENSIV™ pressure parts are presented as absolute pressure transducers; magnetic speed sensors require an appropriate magnetic/ferromagnetic target. — [Pressure sensors](https://www.infineon.com/products/sensor/pressure-sensors); [Magnetic speed sensors](https://www.infineon.com/products/sensor/magnetic-speed-sensors)
- The present PoC uses a manually moved inert mock-up and has no drone. — User-provided setup context (not a product claim).

### Inferences
- Recommended prototype channels: (1) BGT60TR13C radar aimed into the dock for non-contact presence/range/motion; optionally a capacitive electrode for close hand-proximity redundancy; (2) a calibrated, spatially placed airflow instrument from a supplier with a direct suitable flow sensor; (3) temperature as an auxiliary hot-stream feature; (4) built-in fan tach/command and, where useful, motor current/voltage telemetry.
- Label and repeat at least these conditions with built-in fans both off and on: (a) no object/no external flow baseline; (b) external airflow only, including a hair dryer on hot and cool and an unheated fan at varied angles/distances; (c) stationary hand/mock-up intrusion with no applied airflow; (d) object intrusion plus controlled flow. Record synchronized radar, airflow, temperature, and fan state. Establish thresholds/classification from held-out repeats rather than assuming a universal signature.
- The manually moved inert mock-up can validate object-presence detection and flow-only versus object-only plumbing, but cannot validate propeller-generated downwash. The “object + rotor flow” class must wait for a later integrated drone test; use a safely constrained/guarded setup and compare with propellers stopped, motors running, and controlled rotor speeds while built-in fans are separately characterized.
- A practical rule is: flow high + radar absent = flow-only candidate; radar present + flow near fan-only baseline = object-only candidate; radar present + additional reproducible, spatial/temporal flow signature correlated with drone motor/rotor state = object-plus-rotor-flow candidate. The last condition remains an empirical inference and must be tested against an unheated external-flow counterexample.

### Gaps
- No evaluated direct airflow sensor, thresholds, or test data are available yet. In particular, the PoC’s inert mock-up cannot answer whether rotor downwash is separable from an unheated external jet.
- A single temperature or absolute-pressure sensor, CapSense channel, motor-current reading, or radar-presence event is insufficient to identify airflow origin reliably.
