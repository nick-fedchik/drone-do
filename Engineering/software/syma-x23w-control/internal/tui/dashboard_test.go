package tui

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/control"
	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/ipc"
)

func TestEmulatorDashboardFitsTerminalAndLabelsMockData(t *testing.T) {
	snapshot := ipc.Snapshot{
		Running:          true,
		CommandsReceived: 2,
		LastFrame:        &control.Frame{Throttle: 0.3, Roll: 0.1},
		RecentCommands: []ipc.CommandRecord{
			{Number: 1, Frame: control.Frame{}},
			{Number: 2, Frame: control.Frame{Throttle: 0.3, Roll: 0.1}},
		},
	}
	for _, width := range []int{52, 72, 90, 100, 120} {
		view := renderEmulatorDashboard(snapshot, width)
		for lineNumber, line := range strings.Split(view, "\n") {
			if got := lipgloss.Width(line); got > width {
				t.Errorf("width %d line %d is %d cells: %q", width, lineNumber+1, got, line)
			}
		}
		if !strings.Contains(view, "MOCK") || !strings.Contains(view, "PROTOCOL UNVERIFIED") ||
			!strings.Contains(view, "RECENT COMMANDS") || !strings.Contains(view, "COMBINED INPUT") {
			t.Errorf("dashboard at width %d is missing expected status panels", width)
		}
	}
}

func TestManeuverLabelDescribesInputNotMotion(t *testing.T) {
	if got := maneuverLabel(&control.Frame{}); got != "NEUTRAL INPUT" {
		t.Fatalf("neutral label = %q", got)
	}
	if got := maneuverLabel(&control.Frame{Pitch: -0.5}); got != "PITCH -0.5" {
		t.Fatalf("pitch label = %q", got)
	}
	if got := maneuverLabel(&control.Frame{Throttle: 0.2, Yaw: 0.1}); got != "COMBINED INPUT" {
		t.Fatalf("combined label = %q", got)
	}
}
