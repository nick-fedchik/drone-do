# Infineon technology mapping

> **WP type:** Analysis Report<br>
> **Status:** Imported draft; technical claims and conclusions have not been revalidated.<br>
> **Origin:** [`copter-capsense/specs/infineon-options.md`](https://github.com/nick-fedchik/drafts/blob/32d179e513072ba7c08e4f9664fdb3c9f3e99c1b/copter-capsense/specs/infineon-options.md) at `32d179e513072ba7c08e4f9664fdb3c9f3e99c1b`.<br>
> **Migration:** Content preserved; internal links remapped where possible.


**Status:** working draft, v0.1 — 2026-11-11
**Purpose:** map Infineon product families to the docking research without implying component selection or guaranteed system performance.

## 1. Priority candidates

| Priority | Infineon family / example | Possible project role | Main caveat |
|---|---|---|---|
| 1 | PSoC CAPSENSE; CY8CKIT-041S-MAX or CY8CKIT-040T evaluation kits | First proximity/electrode tuning experiment; possible service-panel HMI. | The kits are not a ready-made multi-zone landing array. A metal grid can shield/distort the field; measure the actual stack-up. |
| 2 | XENSIV 60 GHz FMCW radar; BGT60TR13C / DEMO BGT60TR13C | Independent motion/range/angle experiment near the entry. | Published minimum detection range is 0.2 m; not a drone/hand or airflow-source classifier. |
| 3 | REAL3 ToF image sensors; IRS2975C, IRS2877A | Depth-image option for the final approach/occupancy geometry. | Product components are not confirmed as a current ready-to-use camera kit; optics, emitter, driver, and software path need confirmation. |
| Conditional | MOTIX/iMOTION motor-control families | Control/monitor a compatible BLDC fan or actuator. | Motor type, voltage, controller access, and feedback interface are unknown. |
| Conditional | OptiMOS MOSFETs and associated gate drivers | Switching for low-voltage fan, electromagnet, or charging power stage. | Select only after voltage/current/thermal/transient requirements are known. Not a complete charger. |
| Optional | AIROC Wi-Fi/Bluetooth | Embedded connectivity research if an exact supported host/radio platform is selected. | Does not establish OpenWrt or Jetson driver compatibility; not a substitute for choosing a supported OpenWrt router. |

## 2. CapSense / PSoC

**Recommended first Infineon experiment.** Use a proximity-capable PSoC evaluation board to validate electrode, overlay, target-distance, and interference behavior before designing a multi-zone structure. For zonal sensing, independently instrumented electrodes and the chosen PSoC's pin/channel resources must be checked in the CapSense configurator and device documentation.

CapSense evidence should be described as a capacitive proximity event, not “drone detected.” It does not directly measure airflow, object identity, distance as a calibrated 3D coordinate, or safe docking. Test metal/ground planes, grid materials, fan wiring, charger switching, temperature/humidity, and drift.

Sources: [CY8CKIT-041S-MAX](https://www.infineon.com/evaluation-board/CY8CKIT-041S-MAX), [CY8CKIT-040T](https://www.infineon.com/evaluation-board/CY8CKIT-040T), [CAPSENSE Middleware](https://github.com/Infineon/capsense/blob/master/README.md), [AN85951 design guide](https://www.infineon.com/row/public/documents/30/42/infineon-an85951-psoc-4-psoc-6-capsense-design-guide-applicationnotes-en.pdf).

## 3. Radar and ToF

- **XENSIV BGT60TR13C:** useful candidate for a short-range FMCW experiment with raw data, range/angle, and motion. The board and enclosure geometry must be tested, particularly targets closer than the stated 0.2 m minimum. It may add independent motion evidence; it does not establish drone identity or airflow origin.
- **REAL3 ToF:** conceptually suited to depth geometry. IRS2975C is described as a bare-die 240×180 sensor; IRS2877A is automotive-oriented and lists upscaled system resolution. The reviewed research did not establish a current Infineon off-the-shelf ToF camera kit and complete software/optics chain. Do not make it the first purchase without supplier confirmation.

Sources: [BGT60TR13C](https://www.infineon.com/part/BGT60TR13C), [DEMO BGT60TR13C](https://www.infineon.com/evaluation-board/DEMO-BGT60TR13C), [IRS2975C](https://www.infineon.com/part/IRS2975C), [IRS2877A](https://www.infineon.com/part/IRS2877A).

## 4. Fans, magnet, and charging

- **MOTIX/iMOTION:** consider only if the selected fans/actuator have a compatible motor topology and voltage. A controller can report electrical or tachometer-related state, but motor current alone is not airflow measurement.
- **OptiMOS / gate drivers:** possible power switches for a low-voltage fan, magnet coil, or custom power converter. Coil drive needs protection, current/thermal limits, and a defined safe-release behavior.
- **Charging:** separate contact and contactless charging studies. Semiconductor product families do not by themselves define a complete, safe charger. Select a topology only after the battery chemistry, cell count, voltage, power, alignment tolerance, thermal limits, and safety requirements are specified.
- **Airflow:** no directly suitable Infineon airflow-velocity sensor was established in the reviewed product research. DPS368 is an absolute barometric pressure/temperature sensor, not an anemometer or differential-pressure flow sensor. Use a suitable third-party flow instrument for the PoC.

Sources: [Infineon motor drivers](https://www.infineon.com/products/power/motor-drivers), [6EDL7141](https://www.infineon.com/part/6EDL7141), [TLE9879QXA40](https://www.infineon.com/part/TLE9879QXA40), [OptiMOS/MOSFET portfolio](https://www.infineon.com/products/power/mosfet), [DPS368](https://www.infineon.com/part/DPS368).

## 5. PoC recommendation

1. Start with one PSoC CapSense proximity kit and a manually positioned mock-up.
2. Log CapSense signals and video; test the intended grid/overlay and fan-off/fan-on interference.
3. Add a direct airflow sensor to separate object-presence questions from air movement.
4. Add a BGT60TR13C evaluation board only if range/motion evidence is a useful independent channel for the bay geometry.
5. Defer REAL3, motor-drive redesign, magnetic retention, and charging hardware until the earlier experiments and mechanical/electrical requirements justify them.

For researched caveats and source-level detail, see [Infineon technology report](report.md) and [CapSense/HMI report](../capsense-hmi/report.md).
