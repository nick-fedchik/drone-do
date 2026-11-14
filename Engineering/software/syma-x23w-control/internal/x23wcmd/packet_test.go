package x23wcmd

import (
	"testing"
	"time"

	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/control"
)

func TestEncodeFlightFrameNeutralChecksums(t *testing.T) {
	pkt, err := EncodeFlightFrame(control.Frame{At: time.Now().UTC()})
	if err != nil {
		t.Fatal(err)
	}
	if len(pkt) != packetLen {
		t.Fatalf("len = %d, want %d", len(pkt), packetLen)
	}
	if pkt[0] != 0xA5 || pkt[1] != 0x5A || pkt[2] != 12 {
		t.Fatalf("unexpected header: %x", pkt[:3])
	}
	// Neutral throttle=0, pitch/yaw/roll=0x80 (128)
	if pkt[3] != 0 || pkt[4] != 0x80 || pkt[5] != 0x80 || pkt[6] != 0x80 {
		t.Fatalf("unexpected neutral axes: %x", pkt[3:7])
	}
	// Verify outer 16-bit LE sum
	var sum uint16
	for i := 0; i < 13; i++ {
		sum += uint16(pkt[i])
	}
	got := uint16(pkt[13]) | (uint16(pkt[14]) << 8)
	if got != sum {
		t.Fatalf("outer checksum = 0x%04x, want 0x%04x", got, sum)
	}
}

func TestEncodeFlightFrameRejectsInvalidFrame(t *testing.T) {
	if _, err := EncodeFlightFrame(control.Frame{At: time.Now().UTC(), Throttle: 1.2}); err == nil {
		t.Fatal("expected error for out-of-range throttle")
	}
}
