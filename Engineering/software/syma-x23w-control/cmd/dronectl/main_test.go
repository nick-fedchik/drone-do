package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestProtocolStatusDisablesLiveTransport(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"protocol", "status"}, &stdout, &stderr); code != 0 {
		t.Fatalf("run() code = %d, want 0; stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "X23W protocol: UNVERIFIED") {
		t.Fatalf("status output does not mark the protocol unverified: %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "Live radio transport: DISABLED") {
		t.Fatalf("status output does not disable live transport: %q", stdout.String())
	}
}

func TestSimulatePrintsOfflineFrames(t *testing.T) {
	input := `{"at":"2026-11-12T12:00:00Z","throttle":0,"roll":0,"pitch":0,"yaw":0}
{"at":"2026-11-12T12:00:01Z","throttle":0.1,"roll":0,"pitch":0,"yaw":0}`
	path := filepath.Join(t.TempDir(), "sequence.jsonl")
	if err := os.WriteFile(path, []byte(input), 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	var stdout, stderr bytes.Buffer
	if code := run([]string{"simulate", "--input", path}, &stdout, &stderr); code != 0 {
		t.Fatalf("run() code = %d, want 0; stderr=%q", code, stderr.String())
	}
	if got := strings.Count(stdout.String(), `"mode":"SIMULATION_ONLY"`); got != 2 {
		t.Fatalf("simulation output has %d frames, want 2: %q", got, stdout.String())
	}
}

func TestSimulateRejectsMissingInput(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"simulate"}, &stdout, &stderr); code == 0 {
		t.Fatal("run() accepted a missing input file")
	}
}
