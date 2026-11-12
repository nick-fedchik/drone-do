# Drone docking and remote-control R&D

This repository organizes research, experiments, and engineering work for
indoor drone control, passage sensing, and docking concepts. It is an R&D
workspace: proposals, simulations, and concept renders are not evidence of
validated hardware or flight capability.

CapSense is an optional sensing line, not the project's baseline or primary
detection technology. The current imported CapSense material is retained as
unvalidated research input. Sensor selection, drone command protocols, and
performance claims require experiments on the actual equipment.

## Project map

| Area | Purpose | Entry point |
|---|---|---|
| `Research/` | Research questions, source material, analyses, and reports | [Research index](Research/README.md) |
| `Engineering/` | Specifications, designs, software, and validation plans | [Engineering index](Engineering/README.md) |
| `Project/` | Work Item / Work Product rules, templates, registers, and decisions | [Project guide](Project/README.md) |
| `.github/` | GitHub Work Item forms and pull request template | [Issue templates](.github/ISSUE_TEMPLATE/) |

## Working principles

- Separate user-confirmed facts, hypotheses, candidates, and verified results.
- Keep original sources and raw evidence traceable; do not treat migration as
  technical revalidation.
- Define experiment methods and acceptance criteria before collecting results.
- Treat unknown or conflicting sensor evidence as inconclusive.
- Keep flight stabilization and loss-of-link behavior with the aircraft's
  flight controller; early control tests require propellers removed.

See [Project manual](Project/guides/project-manual.md) for how to add Work
Items and Work Products.
