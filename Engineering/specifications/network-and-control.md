# Network and aircraft-control specification

> **WP type:** Specification<br>
> **Status:** Imported draft; technical claims and conclusions have not been revalidated.<br>
> **Origin:** [`copter-capsense/specs/network-and-control.md`](https://github.com/nick-fedchik/drafts/blob/32d179e513072ba7c08e4f9664fdb3c9f3e99c1b/copter-capsense/specs/network-and-control.md) at `32d179e513072ba7c08e4f9664fdb3c9f3e99c1b`.<br>
> **Migration:** Content preserved; internal links remapped where possible.


**Status:** working draft, v0.1 — 2026-11-11
**Scope:** network path for experiments with the Syma X23W; not a validated flight-control interface.

## 1. Known network context

| Node/link | Status |
|---|---|
| Three Deco X60 units | User-confirmed home mesh on the primary PON service. |
| Deco X20-4G | User-confirmed spare LTE device for loss of the PON link; keep it outside the drone-control path. Failover behavior/wiring are not verified here. |
| J501 reServer | User-confirmed Ethernet-only compute node, without built-in Wi-Fi. The dedicated second Ethernet link is the preferred candidate for the drone radio node. |
| M.2 Key E 2230 | User reports the position is occupied by NVMe storage; no internal wireless module is planned. |
| OpenWrt router | Candidate only. Model, radio chipset, hardware revision, OpenWrt release, and STA/forwarding support remain TBD. |
| Syma X23W | Candidate Wi-Fi-controlled toy aircraft. Actual SSID, IP range, ports, firmware, and packet format must be verified on the user's unit. |

## 2. Candidate physical/logical topology

```text
J501 Ethernet 1 ── Deco X60 mesh ── PON / home services
J501 dedicated Ethernet ── point-to-point ── OpenWrt LAN
OpenWrt Wi-Fi STA ── Syma X23W access point
```

The point-to-point link is the preferred candidate because it avoids requiring the J501 itself to have Wi-Fi. It also isolates the radio-control network from ordinary home traffic. Exact interface names, link speed, IP addresses, routes, and host-firewall rules are not specified yet.

If the J501 sends commands through OpenWrt, OpenWrt must route the intended traffic between its Ethernet-side interface and Wi-Fi client interface. Return video traffic may require additional routing or a local relay. A Wi-Fi STA association does not automatically provide transparent Ethernet bridging; do not assume that disabling NAT alone makes the end-to-end route work.

Any firewall should default to deny and allow only the required peers and traffic after the protocol is verified. Do not use broad `input ACCEPT` or unrestricted bidirectional forwarding as a final configuration. Do not add a route to the Deco network until the actual Deco firmware's routing capabilities and the need for access from other hosts are confirmed.

### Radio-node selection priority

The user prioritizes **replaceable antennas and engineering flexibility** over Wi-Fi 6 or mesh capability. Linksys WRT3200ACM and WRT32X are therefore candidates to evaluate, not approved purchases:

- OpenWrt's device records list four detachable antennas, 512 MB RAM, 256 MB NAND, two USB ports (one USB 3.0 and one USB 2.0/eSATA), serial access, and current OpenWrt support for both models. Both are discontinued and use the Marvell 88W8964/mwlwifi radio platform; Wi-Fi is 802.11ac, not Wi-Fi 6. See the [WRT3200ACM record](https://openwrt.org/toh/hwdata/linksys/linksys_wrt3200acm) and [WRT32X record](https://openwrt.org/toh/hwdata/linksys/linksys_wrt32x_v1_venom).
- OpenWrt describes the WRT32X as **nearly** hardware-identical to the WRT3200ACM, with a different OEM interface and network-prioritization feature. Do not treat them as 100% identical or assume one model's image/firmware procedure applies to the other.
- The WRT AC series has dual firmware partitions, but the other partition is only a fallback snapshot and is not guaranteed to be healthy. Recovery behavior and exact installation steps must be checked against the exact model and unit.
- OpenWrt documents a serial console at J1 using a 3.3 V, 6-pin JST-PH connector with 2.0 mm pitch. This is not a generic 2.54 mm header; use the correct cable/pinout and avoid connecting an incompatible voltage.
- Four detachable antennas make RF experiments more convenient, but do not by themselves provide direction finding, per-antenna RSSI, or CSI. Verify the actual connector type, band/chain mapping, and radio/driver support before buying antennas or adapters. Coax length adds loss, and directional antennas improve link geometry only when correctly aimed; they do not measure bearing.
- A separate USB Wi-Fi adapter may be useful for a second radio or packet capture only if its exact chipset, OpenWrt driver, monitor-mode behavior, antenna connector, and required packages are verified. Do not assume an AR9271 provides CSI or calibrated per-antenna RSSI; those are chipset/driver-specific capabilities.

The WRT3200ACM is a plausible used-radio candidate if its detachable antennas, serial access, and OpenWrt support match the experiment. Its age, discontinued status, Marvell driver behavior, exact unit revision, and 2.4 GHz STA reliability with the user's X23W must be checked before selection. Use the J501 for heavy compute/build work; the router's role should remain radio/network transport unless measured needs justify more.

## 3. Syma X23W protocol status

### Verified evidence — Syma X22W only

The public [`rooterkyberian/syma` repository](https://github.com/rooterkyberian/syma) describes **Syma X22W**, not X23W. Its README and Python controller report this candidate protocol for X22W:

- Drone address `172.16.10.1`, UDP destination port `5555`.
- Client binds local UDP port `6666`, sends a one-byte `00` initialization datagram, then sends 27-byte command datagrams. The reference controller defaults to 50 Hz.
- The 27-byte example has a 16-byte command prefix, followed by 9 control/trim/special bytes and two checksum bytes. In packet offsets: axes `16–19`, trims `20–23`, special command `24`, checksums `25–26`.
- Example neutral packet: `436d640001001200010404000a000000808080802020202000550f`.
- The X22W README describes an H.264 video stream using UDP port `6666`.
- The cited implementation documents UDP command delivery; it does not substantiate the pasted claim of TCP initialization.

References: [X22W protocol notes](https://github.com/rooterkyberian/syma/blob/master/README.md), [packet builder/constants](https://github.com/rooterkyberian/syma/blob/master/syma/__init__.py), [controller loop](https://github.com/rooterkyberian/syma/blob/master/syma/controller.py).

### Unverified for the user's X23W

Do **not** assume the X22W packet format, port, local source port, SSID prefix, subnet, video format, command rate, or loss-of-link behavior applies to X23W. The earlier packet table that put the four axes at offsets 19–22 is inconsistent with the X22W reference implementation; its axes are at offsets 16–19.

Before implementing a sender, capture the user's own Syma GO traffic or verify the exact X23W board/firmware. Use a packet capture from a permitted local test setup and compare initialization, neutral commands, stick changes, special commands, and video destination. Do not send takeoff/flight commands based only on the X22W repository.

## 4. Control allocation and safety

- The J501 is the candidate for high-level vision, logging, experiment control, and later navigation workloads.
- OpenWrt may provide Wi-Fi connectivity and a transport/relay service. It is not a safety-rated flight controller.
- Stabilization and loss-of-command behavior must remain with the aircraft's flight controller or a suitable onboard real-time controller. The X23W's actual response to lost Wi-Fi packets is TBD and must be tested.
- PYNQ-Z1 boards are not included in the architecture baseline until their availability and intended role are confirmed.
- Initial protocol tests shall be performed with propellers removed. Any powered rotor or flight test requires a separate controlled and guarded procedure.

## 5. Network requirements and open items

| ID | Requirement / question | Status |
|---|---|---|
| NET-01 | The OpenWrt device shall be checked by exact model/hardware revision against the selected OpenWrt release and the required Wi-Fi STA mode. | TBD; no device selected. |
| NET-02 | The J501-to-radio-node link shall be separate from the J501's home-mesh link if the second Ethernet interface is available. | Preferred candidate. |
| NET-03 | Addressing and routing shall be documented after observing the actual X23W DHCP and route behavior. | TBD; proposed `10.10.10.0/30` is only an example. |
| NET-04 | Firewall rules shall restrict control traffic to the required hosts and ports. | Required; exact rules depend on verified protocol. |
| NET-05 | Video forwarding/relay shall be a separate feature and shall not be assumed necessary for control-only experiments. | Optional; protocol TBD. |
| NET-06 | The radio node should have user-replaceable antennas and accessible engineering interfaces; exact connector, chain mapping, driver behavior, and router revision shall be verified. | User priority; candidate models only. |
| CTRL-01 | The control packet format shall be verified against the user's X23W before a command sender is implemented. | Blocking. |
| CTRL-02 | A Wi-Fi/router failure shall not be treated as a safe flight mode unless the aircraft's behavior is independently verified. | Required safety investigation. |

## 6. Recommended integration sequence

1. Record the exact X23W model/board revision, SSID, DHCP lease, and IP/port behavior using the official app.
2. Verify OpenWrt STA association and stable Ethernet connectivity without controlling the aircraft.
3. Capture/compare command traffic; then test only neutral packets with propellers removed.
4. Test low-authority commands and the actual loss-of-link response in a controlled, propeller-removed or guarded setup.
5. Only after protocol and routing are stable, move high-level command generation to J501 and test the dedicated Ethernet path.
6. Keep visual navigation, sensor fusion, and flight-safety decisions separate from the Wi-Fi relay until the aircraft/control stack is known.
