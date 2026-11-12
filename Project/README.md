# Project workspace

* **WP type:** Guide
* **Version:** 0.2.0
* **Status:** Working project structure

## Governance and navigation

| Work Product | Purpose |
|---|---|
| [Project manual](guides/project-manual.md) | Create and maintain Work Items and Work Products |
| [Layout policy](policies/layout.md) | Place material in Research, Engineering, or Project |
| [Work Item policy](policies/work-items.md) | Classify tracked work and research questions |
| [Work Product policy](policies/work-products.md) | Classify and format project outputs |
| [Template guide](templates/README.md) | Find Work Product and Work Item templates |
| [Work Product register](registers/work-products.md) | Index current project documents |
| [Migration register](registers/source-migration.md) | Trace imported CapSense material to its source |
| [Hardware register](registers/hardware-procurement.md) | Inventory and procurement candidates |
| [Open questions](registers/open-questions.md) | Record current facts, hypotheses, and unresolved choices |

Technical decisions belong in `Engineering/decisions/`; project-scope and
governance decisions belong in `Project/decisions/`. Create either directory
when its first decision record is added.

Templates are aids, not a requirement to create a separate file for every
question or small task. Work Item types and Work Product types describe
different things: work being tracked versus its resulting artifact.

## Version history

| Version | Date | Change |
|---|---|---|
| 0.1.0 | 2026-11-12 | Initial R&D workspace structure. |
| 0.2.0 | 2026-11-12 | Clarify decision-record ownership. |
