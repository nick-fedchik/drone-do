package tui

import (
	"fmt"
	"io"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/control"
	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/ipc"
)

const axisStep = 0.1
const throttleStep = 0.05

type controllerModel struct {
	client   *ipc.Client
	frame    control.Frame
	snapshot ipc.Snapshot
	sequence uint64
	message  string
	width    int
}

func RunController(client *ipc.Client, input *os.File, output io.Writer) error {
	snapshot, err := client.Status()
	if err != nil {
		return fmt.Errorf("read emulator status: %w", err)
	}
	frame := control.Frame{At: time.Now().UTC()}
	snapshot, err = client.Send(1, frame)
	if err != nil {
		return fmt.Errorf("send initial neutral frame: %w", err)
	}

	model := controllerModel{
		client:   client,
		frame:    frame,
		snapshot: snapshot,
		sequence: 1,
		width:    100,
	}
	program := tea.NewProgram(model,
		tea.WithAltScreen(),
		tea.WithInput(input),
		tea.WithOutput(output))
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("run controller TUI: %w", err)
	}
	return nil
}

func (m controllerModel) Init() tea.Cmd {
	return nil
}

func (m controllerModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		m.width = message.Width
	case tea.KeyMsg:
		if message.Type == tea.KeyCtrlC || message.Type == tea.KeyEsc || keyRune(message) == 'x' {
			m.sendNeutral()
			return m, tea.Quit
		}
		key := keyRune(message)
		if key == 0 {
			return m, nil
		}
		changed, hint := applyControllerKey(&m.frame, key)
		if !changed {
			m.message = hint
			return m, nil
		}
		m.message = ""
		m.frame.At = time.Now().UTC()
		m.sequence++
		snapshot, err := m.client.Send(m.sequence, m.frame)
		if err != nil {
			m.message = fmt.Sprintf("IPC error: %v", err)
			return m, nil
		}
		m.snapshot = snapshot
	}
	return m, nil
}

func (m controllerModel) View() string {
	return renderControllerDashboard(m.frame, m.snapshot, m.message, m.width)
}

func (m *controllerModel) sendNeutral() {
	m.frame = control.Frame{At: time.Now().UTC()}
	m.sequence++
	snapshot, err := m.client.Send(m.sequence, m.frame)
	if err != nil {
		m.message = fmt.Sprintf("neutral command failed: %v", err)
		return
	}
	m.snapshot = snapshot
}

func keyRune(key tea.KeyMsg) byte {
	if key.Type == tea.KeySpace {
		return '0'
	}
	if key.Type != tea.KeyRunes || len(key.Runes) != 1 || key.Runes[0] > 127 {
		return 0
	}
	return byte(key.Runes[0])
}

func applyControllerKey(frame *control.Frame, key byte) (bool, string) {
	switch key {
	case 'w':
		frame.Pitch = clamp(frame.Pitch+axisStep, -1, 1)
	case 's':
		frame.Pitch = clamp(frame.Pitch-axisStep, -1, 1)
	case 'a':
		frame.Roll = clamp(frame.Roll-axisStep, -1, 1)
	case 'd':
		frame.Roll = clamp(frame.Roll+axisStep, -1, 1)
	case 'q':
		frame.Yaw = clamp(frame.Yaw-axisStep, -1, 1)
	case 'e':
		frame.Yaw = clamp(frame.Yaw+axisStep, -1, 1)
	case 'r':
		frame.Throttle = clamp(frame.Throttle+throttleStep, 0, 1)
	case 'f':
		frame.Throttle = clamp(frame.Throttle-throttleStep, 0, 1)
	case '0':
		frame.Throttle, frame.Roll, frame.Pitch, frame.Yaw = 0, 0, 0, 0
	default:
		return false, fmt.Sprintf("Unknown key %q", key)
	}
	return true, ""
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
