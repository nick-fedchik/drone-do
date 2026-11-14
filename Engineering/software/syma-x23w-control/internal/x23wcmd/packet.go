package x23wcmd

import (
	"fmt"
	"math"

	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/control"
)

const (
	// FlightUDPPort is the Joyhonest IC_GPRTPB flight-command UDP port from jh_wifi.cpp.
	FlightUDPPort = 25000
	innerLen      = 10
	packetLen     = 15
)

// EncodeFlightFrame converts a normalized control.Frame into a 15-byte
// Joyhonest IC_GPRTPB 0xA5 0x5A packet (for future bench verification; not sent here).
func EncodeFlightFrame(frame control.Frame) ([]byte, error) {
	if err := frame.Validate(); err != nil {
		return nil, fmt.Errorf("invalid frame: %w", err)
	}

	inner := make([]byte, innerLen)
	inner[0] = byte(math.Round(frame.Throttle * 255.0))
	inner[1] = signedAxisByte(frame.Pitch)
	inner[2] = signedAxisByte(frame.Yaw)
	inner[3] = signedAxisByte(frame.Roll)
	inner[4] = 0x20
	inner[5] = 0x20
	inner[6] = 0x20
	inner[7] = 0x20
	inner[8] = 0x00

	var xor uint8
	for i := 0; i < 9; i++ {
		xor ^= inner[i]
	}
	inner[9] = xor + 0x55

	pkt := make([]byte, packetLen)
	pkt[0] = 0xA5
	pkt[1] = 0x5A
	pkt[2] = innerLen + 2
	copy(pkt[3:3+innerLen], inner)

	var sum uint16
	for i := 0; i < 3+innerLen; i++ {
		sum += uint16(pkt[i])
	}
	pkt[13] = byte(sum & 0xFF)
	pkt[14] = byte(sum >> 8)
	return pkt, nil
}

func signedAxisByte(v float64) byte {
	// Center is 0x80 (128); range [-1, 1] maps to [0x01, 0xFF].
	scaled := 128.0 + v*127.0
	if scaled < 0 {
		scaled = 0
	} else if scaled > 255 {
		scaled = 255
	}
	return byte(math.Round(scaled))
}
