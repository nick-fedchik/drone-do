# Infineon Wi-Fi for OpenWrt and indoor drone navigation

> **WP type:** Analysis Report — supporting research note<br>
> **Status:** Imported draft; technical claims and conclusions have not been revalidated.<br>
> **Origin:** [`copter-capsense/research_notes/Технології Infineon для дока/wifi.md`](https://github.com/nick-fedchik/drafts/blob/32d179e513072ba7c08e4f9664fdb3c9f3e99c1b/copter-capsense/research_notes/%D0%A2%D0%B5%D1%85%D0%BD%D0%BE%D0%BB%D0%BE%D0%B3%D1%96%D1%97%20Infineon%20%D0%B4%D0%BB%D1%8F%20%D0%B4%D0%BE%D0%BA%D0%B0/wifi.md) at `32d179e513072ba7c08e4f9664fdb3c9f3e99c1b`.<br>
> **Migration:** Content preserved; internal links remapped where possible.


## Which current Infineon Wi-Fi 6 / AIROC products and boards are documented?

### Takeaway
Infineon documents AIROC Wi-Fi 6/6E combo parts and embedded evaluation boards, but those facts establish radio capabilities—not a drop-in Linux/OpenWrt router or Jetson interface. A useful compatibility distinction: the current CYW955513EVK-01 is for the CYW55513/55512/55511 family, whereas the public Wi-Fi Host Driver (WHD) Wi-Fi 6 table lists 55500, 55900, and 55572; the public docs do not establish that the EVK is supported by WHD or OpenWrt.

### Cited Findings
- Infineon's active CYW955513EVK-01 page identifies it as an evaluation board for AIROC CYW55513/2/1, a 1x1 Wi-Fi 6E (802.11ax) and Bluetooth 5.4 combo; it describes embedded audio development on the chip's 192-MHz Arm Cortex-M33. [Infineon CYW955513EVK-01](https://www.infineon.com/evaluation-board/CYW955513EVK-01)
- Infineon's WHD repository describes its software as an embedded Wi-Fi host driver; the Wi-Fi 6 table lists chip families 55500, 55900, and 55572, with SDIO bus support listed for 55500 and 55572 and OCI for 55900. It lists WPA2/WPA3 AP operation for all three. [Infineon WHD README](https://github.com/Infineon/wifi-host-driver/blob/master/README.md)
- WHD API documentation describes a porting layer requiring platform HAL, RTOS, BSP/bus definitions, and network/buffer abstractions; it supports station and AP modes, but this is not documentation of an OpenWrt package or Linux cfg80211 driver. [WHD API reference](https://infineon.github.io/wifi-host-driver/html/index.html)
- The cited public documentation does not supply a current Infineon Wi-Fi 6 router model, a Linux/Jetson module compatibility matrix, or a specific Infineon Wi-Fi 6 FTM-capable board. The 55513 EVK page is specifically oriented to embedded/Bluetooth audio and Cortex-M33 use. [EVK page](https://www.infineon.com/evaluation-board/CYW955513EVK-01); [WHD README](https://github.com/Infineon/wifi-host-driver/blob/master/README.md)

### Inferences
- The EVK could help explore the specific 55513-family radio in its intended development environment, but it should not be assumed to function as a USB/SDIO Linux Wi-Fi adapter, OpenWrt AP, or Jetson client without board-interface, driver, firmware, and host-support documentation.
- WHD AP mode does not imply OpenWrt compatibility: WHD is an embedded abstraction/API stack, while OpenWrt depends on support for the exact target radio and driver integration.

### Gaps
- Verify official product datasheet/driver availability and interface for any selected 55513/55572 board, plus documented Linux kernel and OpenWrt version support. The cited material does not close these gaps.

## Is OpenWrt support documented for specific Infineon chipsets/drivers?

### Takeaway
No official OpenWrt support for a specific AIROC Wi-Fi 6 chipset was established in the sources checked. OpenWrt support must be verified for the exact router/device and its drivers; the phrase “Wi-Fi 6” alone is not evidence that an OpenWrt image supports its radio.

### Cited Findings
- OpenWrt's Table of Hardware is a listing of supported devices, and its documentation points users to that device-specific inventory. [OpenWrt Table of Hardware](https://openwrt.org/toh/start)
- OpenWrt's wireless configuration docs state that the physical radio and driver determine available configuration; `mac80211` is the usual type for platforms other than the listed Broadcom `brcm47xx` case. This is general configuration guidance, not a chipset compatibility guarantee. [OpenWrt wireless configuration](https://openwrt.org/docs/guide-user/network/wifi/basic)
- Infineon's public WHD repository calls WHD an embedded host driver and names Mbed OS, FreeRTOS, and Azure RTOS as portability examples; its API reference requires a HAL/RTOS/BSP port. Neither source identifies OpenWrt as a supported target. [WHD README](https://github.com/Infineon/wifi-host-driver/blob/master/README.md); [WHD API reference](https://infineon.github.io/wifi-host-driver/html/index.html)
- The Linux wireless documentation explains that Linux 802.11 devices need cfg80211 directly or via mac80211 for a consistent nl80211 userspace interface. That Linux architecture does not itself establish that an Infineon chipset driver exists or is packaged in OpenWrt. [Linux cfg80211 documentation](https://www.kernel.org/doc/html/latest/driver-api/80211/cfg80211.html)

### Inferences
- For the user's existing Wi-Fi 6 OpenWrt router, first identify the exact manufacturer/model and radio chipset, then verify its OpenWrt device page, driver/firmware and supported mode. Do not select or buy an Infineon chip based solely on Wi-Fi 6 branding.
- For a drone client built around the reported Jetson J501 64G, require explicit support for the exact module, host bus, kernel, firmware, regulatory configuration, and AP/client mode. The reported computer does not by itself establish radio compatibility.

### Gaps
- The router model and radio chipset are unknown. No official OpenWrt page, package, board support, or driver source for a named AIROC Wi-Fi 6 chipset was located in the cited sources; this is a bounded search finding, not proof that no third-party or out-of-tree port exists.

## Is RTT/FTM available, and what should Wi-Fi do in sensor fusion?

### Takeaway
No cited Infineon product, evaluation-board, or WHD documentation confirms FTM/RTT support for a named AP/client pair. Treat ordinary Wi-Fi as communications, device identity/presence, and telemetry by default; treat ranging only as a separately validated experiment, not as navigation or docking truth.

### Cited Findings
- Google's Wi-Fi RTT documentation says ranging requires FTM-capable hardware on both the requesting device and the AP, plus a supported software/API stack; the documented Wi-Fi RTT implementation measures distance using FTM and describes typical 1–2 m positioning when distances to three or more APs are multilaterated. This is Android platform guidance and not a performance promise for Infineon hardware, OpenWrt, or a drone. [Android Wi-Fi RTT](https://developer.android.com/develop/connectivity/wifi/wifi-rtt)
- The WHD public feature/API pages identify STA/AP and security features, but the inspected material does not document FTM initiator/responder roles, RTT APIs, ranging calibration, or accuracy. [WHD README](https://github.com/Infineon/wifi-host-driver/blob/master/README.md); [WHD API reference](https://infineon.github.io/wifi-host-driver/html/index.html)
- OpenWrt device support and standard wireless configuration documentation do not establish FTM operation for an unspecified radio. [OpenWrt ToH](https://openwrt.org/toh/start); [OpenWrt wireless configuration](https://openwrt.org/docs/guide-user/network/wifi/basic)

### Inferences
- A viable Wi-Fi RTT test needs a specifically documented FTM-capable AP and client, both-end driver/firmware support, user-space/API access on each side, and empirical testing in the intended environment. A Wi-Fi 6 or 6E label, AIROC branding, or successful association does not establish any of these.
- Even vendor-stated meter-scale RTT positioning is much coarser than a docking target; multipath, antenna placement/orientation, and platform implementation would require evaluation before use. For autonomous indoor flight/docking, use Wi-Fi for command/status and telemetry, and validate navigation with suitable onboard sensing (for example visual/inertial sensing and a close-range docking cue) rather than assuming Wi-Fi gives position.
- The no-drone-yet/custom-platform context makes this a future integration question: choose a radio only after platform interfaces and the flight/docking sensor stack are defined.

### Gaps
- No reliable cited evidence was found that CYW55513/55512/55511, CYW55572, or the CYW955513EVK-01 implement Wi-Fi FTM/RTT, nor that any named OpenWrt router provides an FTM responder. No ranging accuracy for this user's future setup can be claimed.
