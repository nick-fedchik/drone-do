package x23wcmd

import (
	"fmt"
	"math"

	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/control"
)

const (
	// FlightUDPPort is the verified Joyhonest IC_GPRTPB flight-command UDP port.
	FlightUDPPort = 25000
	innerLen      = 10
	packetLen     = 15

	FlagTakeoff   uint8 = 0x01 // cmd[7] |= 0x40 (bUp)
	FlagLand      uint8 = 0x02 // cmd[8] |= 0x08 (bDn)
	FlagStop      uint8 = 0x04 // cmd[8] |= 0x10 (bStop)
	FlagCalibrate uint8 = 0x08 // cmd[8] |= 0x20 (bAdj)
)

// EncodeFlightFrame converts a normalized control.Frame into the verified 15-byte
// Syma X23W (IC_GPRTPB 0xA5 0x5A + 10-byte Syma payload) packet for UDP 25000.
func EncodeFlightFrame(frame control.Frame, flags uint8) ([]byte, error) {
	if err := frame.Validate(); err != nil {
		return nil, fmt.Errorf("invalid frame: %w", err)
	}

	inner := make([]byte, innerLen)
	inner[0] = byte(math.Round(frame.Throttle * 255.0))
	inner[1] = signedAxisByte(frame.Pitch)
	inner[2] = signedAxisByte(frame.Yaw)
	inner[3] = signedAxisByte(frame.Roll)
	inner[4] = 0x20
	inner[5] = 0x40
	inner[6] = 0x20
	inner[7] = 0x20
	if flags&FlagTakeoff != 0 {
		inner[7] |= 0x40
	}
	inner[8] = 0x00
	if flags&FlagLand != 0 {
		inner[8] |= 0x08
	}
	if flags&FlagStop != 0 {
		inner[8] |= 0x10
	}
	if flags&FlagCalibrate != 0 {
		inner[8] |= 0x20
	}

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
