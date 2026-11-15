package x23wcmd

import (
	"testing"
	"time"

	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/control"
	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/ipc"
)

func TestPacketFromSnapshotWatchdogTriggersStop(t *testing.T) {
	stale := time.Now().UTC().Add(-2 * time.Second)
	snap := ipc.Snapshot{
		Running:   true,
		UpdatedAt: &stale,
		LastFrame: &control.Frame{At: stale, Throttle: 0.6, Roll: 0.9},
	}
	pkt, err := PacketFromSnapshot(snap, time.Now().UTC(), 500*time.Millisecond, 0.35)
	if err != nil {
		t.Fatal(err)
	}
	// FlagStop sets cmd[8] (pkt[11]) to 0x10
	if pkt[11] != 0x10 {
		t.Fatalf("stale snapshot flag byte = 0x%02x, want 0x10 (FlagStop)", pkt[11])
	}
}

func TestPacketFromSnapshotClampsTilt(t *testing.T) {
	now := time.Now().UTC()
	snap := ipc.Snapshot{
		Running:   true,
		UpdatedAt: &now,
		LastFrame: &control.Frame{At: now, Throttle: 0.5, Roll: 0.9, Pitch: -0.8},
	}
	pkt, err := PacketFromSnapshot(snap, now, 500*time.Millisecond, 0.35)
	if err != nil {
		t.Fatal(err)
	}
	// Roll=0.35 -> 128 + 0.35*127 = 172 (0xAC), Pitch=-0.35 -> 128 - 0.35*127 = 84 (0x54)
	if pkt[6] != 172 || pkt[4] != 84 || pkt[11] != 0 {
		t.Fatalf("clamped axes = roll:%d pitch:%d flags:0x%02x", pkt[6], pkt[4], pkt[11])
	}
}
