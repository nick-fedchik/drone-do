# Syma X23W control utility

* **WP type:** Source Code
* **Version:** 0.4.0
* **Date:** 2026-11-14
* **Status:** Local IPC controller/emulator + read-only X23W Wi-Fi probe/video capture; flight commands unverified
* **Related WI:** [#3](https://github.com/nick-fedchik/drone-do/issues/3)
* **Related WI:** [#4](https://github.com/nick-fedchik/drone-do/issues/4)

This Go module contains `dronectl` (the controller) and `droneemu` (the
pseudo-drone). They exchange versioned JSON messages over a local Unix domain
socket, and `dronectl` also includes read-only `probe` and `video` commands
for the Joyhonest `IC_GPRTPB` camera/Wi-Fi board (`192.168.29.1`).

## Safety and protocol boundary

- On the user's X23W (`FPV_WIFI__140e`, `192.168.29.1`), read-only `FDWN`
  firmware query (`UDP 20001`, firmware `F.WD1B06-20220325-005`), `JHCMD`
  handshake (`UDP 20000`), and `640x480 @ 16 FPS` MJPEG streaming (`UDP 10900`)
  have been verified on the bench.
- Flight-control packet mapping (`0xA5 0x5A` framing on `UDP 25000`) and
  loss-of-link behavior remain **unverified**. Live flight transport is
  disabled.
- The IPC transport is a local Unix domain socket only. No TCP, UDP, or radio
  transport is implemented. `protocol status` reports that live radio
  transport is disabled.
- `dronectl status`, `send`, and `tui` communicate only with `droneemu`.
  Their abstract control values are not aircraft commands.
- The socket is created with mode `0600` under the user's cache directory by
  default. Existing socket paths are never removed automatically.
- `simulate` validates an abstract, normalized input sequence only. It does
  not emulate aircraft physics, validate control response, or establish
  flight safety.
- Do not add a live transport until an X23W-specific capture, protocol
  analysis, and separately reviewed bench-test procedure exist.

## Commands

From this directory:

```sh
go test ./...
go vet ./...
go run ./cmd/dronectl protocol status
go run ./cmd/dronectl simulate --input examples/neutral.jsonl
go run ./cmd/dronectl probe --target 192.168.29.1 --duration 2s
go run ./cmd/dronectl video --target 192.168.29.1 --duration 3s --out /tmp/x23w-frames
go run ./cmd/droneemu serve
```

For the interactive emulator and controller walkthrough, see
[TESTING.md](TESTING.md).

Start `droneemu serve` in one terminal, then use these commands in another:

```sh
go run ./cmd/dronectl status
go run ./cmd/dronectl send --throttle 0.25 --roll -0.5 --pitch 0.75 --yaw 0.1
go run ./cmd/dronectl status
go run ./cmd/dronectl tui
```

`droneemu tui` can be used instead of `droneemu serve` to display the current
emulator state interactively. `dronectl tui` sends a neutral frame on startup
and exit. Its keys adjust normalized values in the pseudo-emulator only.
Both TUIs use Bubble Tea. The emulator dashboard uses a compact, panel-based
layout inspired by BashTOP's information density, with airframe status,
control-input meters, battery placeholder, and recent-command history. The
battery is a fixed `100%` mock value with no discharge model; airframe
specifications and aircraft telemetry remain unknown.

`simulate` accepts one JSON object per line. `at` is an RFC 3339 timestamp;
`throttle` is normalized to `[0,1]`, and `roll`, `pitch`, and `yaw` to
`[-1,1]`. These are abstract ranges, not verified X23W wire-field mappings.
Frames must be strictly time-ordered. Unknown fields and invalid values are
rejected.

## Next protocol work

Capture the official app's traffic for the exact X23W board/firmware in a
controlled local setup after the aircraft is available and charged. Preserve
the capture and document the evidence before implementing any packet encoder
or live transport.

## Version history

| Version | Date | Change |
|---|---|---|
| 0.1.0 | 2026-11-12 | Offline-only CLI scaffold and input validation. |
| 0.2.0 | 2026-11-12 | Add local Unix-socket IPC, controller and pseudo-drone CLIs/TUIs, tests, and tester guide. |
| 0.3.0 | 2026-11-13 | Replace hand-built TUIs with Bubble Tea and add a telemetry-style emulator dashboard. |
| 0.4.0 | 2026-11-14 | Add read-only `probe` and `video` commands for `IC_GPRTPB` (`192.168.29.1`). |
