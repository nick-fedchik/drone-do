package tui

import (
	"testing"

	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/control"
)

func TestApplyControllerKeyClampsAxesAndThrottle(t *testing.T) {
	frame := control.Frame{Throttle: 0.99, Roll: 0.95}
	for range 5 {
		changed, _ := applyControllerKey(&frame, 'r')
		if !changed {
			t.Fatal("throttle key was not applied")
		}
		changed, _ = applyControllerKey(&frame, 'd')
		if !changed {
			t.Fatal("roll key was not applied")
		}
	}
	if frame.Throttle != 1 || frame.Roll != 1 {
		t.Fatalf("controls exceeded expected bounds: %+v", frame)
	}

	changed, _ := applyControllerKey(&frame, '0')
	if !changed || frame.Throttle != 0 || frame.Roll != 0 || frame.Pitch != 0 || frame.Yaw != 0 {
		t.Fatalf("neutral key result = %+v, changed=%t", frame, changed)
	}
}

func TestApplyControllerKeyRejectsUnknownKey(t *testing.T) {
	frame := control.Frame{}
	changed, message := applyControllerKey(&frame, 'z')
	if changed || message == "" {
		t.Fatalf("unknown key result: changed=%t message=%q", changed, message)
	}
}
