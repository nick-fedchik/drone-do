package sequence

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/control"
)

const maxLineSize = 1 << 20

// Read parses JSONL control frames, validates each frame, and requires
// strictly increasing timestamps.
func Read(r io.Reader) ([]control.Frame, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 4096), maxLineSize)

	frames := make([]control.Frame, 0)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Bytes()
		if len(bytes.TrimSpace(line)) == 0 {
			continue
		}

		var frame control.Frame
		decoder := json.NewDecoder(bytes.NewReader(line))
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&frame); err != nil {
			return nil, fmt.Errorf("line %d: decode frame: %w", lineNumber, err)
		}
		var trailing any
		if err := decoder.Decode(&trailing); err != io.EOF {
			if err == nil {
				return nil, fmt.Errorf("line %d: expected one JSON object", lineNumber)
			}
			return nil, fmt.Errorf("line %d: trailing data: %w", lineNumber, err)
		}
		if err := frame.Validate(); err != nil {
			return nil, fmt.Errorf("line %d: %w", lineNumber, err)
		}
		if len(frames) > 0 && !frame.At.After(frames[len(frames)-1].At) {
			return nil, fmt.Errorf("line %d: timestamps must be strictly increasing", lineNumber)
		}
		frames = append(frames, frame)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read sequence: %w", err)
	}
	if len(frames) == 0 {
		return nil, fmt.Errorf("sequence contains no frames")
	}
	return frames, nil
}

// ReadString is a small helper for tests and in-memory examples.
func ReadString(input string) ([]control.Frame, error) {
	return Read(strings.NewReader(input))
}
