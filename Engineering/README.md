# Engineering

* **WP type:** Guide
* **Version:** 0.6.0
* **Date:** 2026-11-13
* **Status:** Initial structure; imported specifications remain drafts

## Specifications

- [System architecture](specifications/system-architecture.md)
- [Sensing and docking](specifications/sensing-and-docking.md)
- [Network and aircraft control](specifications/network-and-control.md)

The [specification index](specifications/README.md) describes the imported
documents and evidence labels. The
[verification plan](validation/verification-plan.md) is a proposed sequence,
not approval to conduct a flight or charging test. Bench results for Jetson
J501 Wi-Fi, X23W `IC_GPRTPB` video/query, and TensorRT are recorded in
[Jetson/X23W bench report](validation/jetson-x23w-bench-report.md); see
[pre-flight checklist](validation/pre-flight-checklist.md) before blades-on tests.

## Decisions and software

Technical decision records belong in `Engineering/decisions/`; create this
directory when the first decision is recorded.

When software is introduced, place each independently buildable module or
component in its own folder under `Engineering/software/`, regardless of
language. Keep that module's source, tests, build files, and documentation
together; do not make `Engineering/software/` itself a single-language
module root.

For Go, put `go.mod` at the individual module's root and use conventional
`cmd/` and `internal/` packages where they fit. Other languages keep their
own manifests and toolchain files in their module folder. Keep TinyGo
firmware separate when it has a different target or toolchain. Do not create
modules or empty folders before there is code to place there.

Hardware descriptions and designs belong in `Engineering/hardware/` once
those Work Products exist.

## Software

- [Syma X23W control utility](software/syma-x23w-control/README.md) —
  Bubble Tea controller/pseudo-drone IPC PoC; X23W protocol remains
  unverified and live aircraft transport is disabled. See its
  [tester guide](software/syma-x23w-control/TESTING.md).

## Version history

| Version | Date | Change |
|---|---|---|
| 0.1.0 | 2026-11-12 | Initial Engineering index. |
| 0.2.0 | 2026-11-12 | Define technical decision and future software placement. |
| 0.3.0 | 2026-11-12 | Clarify support for multiple language-independent software modules. |
| 0.4.0 | 2026-11-12 | Index the offline Syma X23W control utility module. |
| 0.5.0 | 2026-11-12 | Update the Syma utility index for its local IPC emulator PoC. |
| 0.6.0 | 2026-11-13 | Note Bubble Tea-based controller and emulator dashboards. |
