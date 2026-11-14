package sequence

import (
	"strings"
	"testing"
)

func TestReadValidSequence(t *testing.T) {
	input := `{"at":"2026-11-12T12:00:00Z","throttle":0,"roll":0,"pitch":0,"yaw":0}
{"at":"2026-11-12T12:00:01Z","throttle":0.1,"roll":0.05,"pitch":0,"yaw":0}`

	frames, err := ReadString(input)
	if err != nil {
		t.Fatalf("ReadString() error = %v", err)
	}
	if got, want := len(frames), 2; got != want {
		t.Fatalf("ReadString() returned %d frames, want %d", got, want)
	}
}

func TestReadRejectsOutOfRangeControl(t *testing.T) {
	input := `{"at":"2026-11-12T12:00:00Z","throttle":1.1,"roll":0,"pitch":0,"yaw":0}`
	if _, err := ReadString(input); err == nil {
		t.Fatal("ReadString() accepted out-of-range throttle")
	}
}

func TestReadRejectsNonIncreasingTimestamps(t *testing.T) {
	input := `{"at":"2026-11-12T12:00:01Z","throttle":0,"roll":0,"pitch":0,"yaw":0}
{"at":"2026-11-12T12:00:01Z","throttle":0,"roll":0,"pitch":0,"yaw":0}`
	if _, err := ReadString(input); err == nil {
		t.Fatal("ReadString() accepted duplicate timestamps")
	}
}

func TestReadRejectsUnknownFields(t *testing.T) {
	input := `{"at":"2026-11-12T12:00:00Z","throttle":0,"roll":0,"pitch":0,"yaw":0,"motor":1}`
	if _, err := ReadString(input); err == nil {
		t.Fatal("ReadString() accepted an unknown field")
	}
}

func TestReadRejectsMultipleObjectsOnOneLine(t *testing.T) {
	input := `{"at":"2026-11-12T12:00:00Z","throttle":0,"roll":0,"pitch":0,"yaw":0} {"at":"2026-11-12T12:00:01Z","throttle":0,"roll":0,"pitch":0,"yaw":0}`
	if _, err := ReadString(input); err == nil || !strings.Contains(err.Error(), "expected one JSON object") {
		t.Fatalf("ReadString() error = %v, want multiple-object error", err)
	}
}
