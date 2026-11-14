package control

import (
	"fmt"
	"math"
	"time"
)

// Frame is an abstract, normalized control input for offline tooling.
// Its axes are not mapped to the unverified X23W wire protocol.
type Frame struct {
	At       time.Time `json:"at"`
	Throttle float64   `json:"throttle"`
	Roll     float64   `json:"roll"`
	Pitch    float64   `json:"pitch"`
	Yaw      float64   `json:"yaw"`
}

func (f Frame) Validate() error {
	if f.At.IsZero() {
		return fmt.Errorf("at must be a non-zero timestamp")
	}
	if !inRange(f.Throttle, 0, 1) {
		return fmt.Errorf("throttle must be finite and in [0,1]")
	}
	if !inRange(f.Roll, -1, 1) {
		return fmt.Errorf("roll must be finite and in [-1,1]")
	}
	if !inRange(f.Pitch, -1, 1) {
		return fmt.Errorf("pitch must be finite and in [-1,1]")
	}
	if !inRange(f.Yaw, -1, 1) {
		return fmt.Errorf("yaw must be finite and in [-1,1]")
	}
	return nil
}

func inRange(value, min, max float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0) && value >= min && value <= max
}
