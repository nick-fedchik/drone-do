package x23wcmd

import (
	"bytes"
	"encoding/hex"
	"testing"
	"time"

	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/control"
)

func TestEncodeFlightFrameMatchesVerifiedNeutralAndStop(t *testing.T) {
	pkt, err := EncodeFlightFrame(control.Frame{At: time.Now().UTC(), Throttle: 0.5}, 0)
	if err != nil {
		t.Fatal(err)
	}
	wantNeutral, _ := hex.DecodeString("a55a0c808080802040202000b56004")
	if !bytes.Equal(pkt, wantNeutral) {
		t.Fatalf("neutral = %x, want %x", pkt, wantNeutral)
	}

	stopPkt, err := EncodeFlightFrame(control.Frame{At: time.Now().UTC(), Throttle: 0.5}, FlagStop)
	if err != nil {
		t.Fatal(err)
	}
	if stopPkt[11] != 0x10 {
		t.Fatalf("stop flag byte = 0x%02x, want 0x10", stopPkt[11])
	}
}

func TestEncodeFlightFrameRejectsInvalidFrame(t *testing.T) {
	if _, err := EncodeFlightFrame(control.Frame{At: time.Now().UTC(), Throttle: 1.2}, 0); err == nil {
		t.Fatal("expected error for out-of-range throttle")
	}
}
