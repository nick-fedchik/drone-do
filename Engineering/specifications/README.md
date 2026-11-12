# Copter CapSense: specifications

> **WP type:** Guide<br>
> **Status:** Imported draft; technical claims and conclusions have not been revalidated.<br>
> **Origin:** [`copter-capsense/specs/README.md`](https://github.com/nick-fedchik/drafts/blob/32d179e513072ba7c08e4f9664fdb3c9f3e99c1b/copter-capsense/specs/README.md) at `32d179e513072ba7c08e4f9664fdb3c9f3e99c1b`.<br>
> **Migration:** Content preserved; internal links remapped where possible.


**Status:** working draft, v0.1 — 2026-11-11
**Purpose:** turn the research concept and recent hardware clarifications into an engineering specification that can be reviewed and refined with the CapSense team.

## Documents

- [System architecture](system-architecture.md) — objective, boundaries, reported equipment, and system-level requirements.
- [Sensing and docking](sensing-and-docking.md) — CapSense, radar/ToF/camera, airflow, fans, HMI, and charging hypotheses.
- [Network and control](network-and-control.md) — J501, home network, candidate OpenWrt link, and provisional Syma X23W control protocol.
- [Infineon technology mapping](../../Research/studies/sensor-options/infineon-technology-map.md) — candidate Infineon families, fit, limitations, and PoC priority.
- [Verification plan](../validation/verification-plan.md) — staged experiments, evidence to collect, and acceptance criteria to set.

## Evidence labels

- **Confirmed by user** — explicitly stated as part of the current setup or intent.
- **Candidate** — a proposed component or architecture, not yet selected or tested.
- **Unverified** — a claim that needs confirmation from the actual hardware, firmware, or a measurement.
- **TBD** — a requirement or parameter intentionally left open until the team defines it.

These are research and PoC specifications, not a production design, safety case, procurement authorization, or claim of flight readiness. The root-level research concept, procurement register, and reports remain the source documents for their respective topics.
