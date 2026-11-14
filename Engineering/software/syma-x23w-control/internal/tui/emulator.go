package tui

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/ipc"
)

type emulatorModel struct {
	ctx      context.Context
	server   *ipc.Server
	snapshot ipc.Snapshot
	width    int
}

type refreshMessage time.Time

func RunEmulator(ctx context.Context, server *ipc.Server, input *os.File, output io.Writer) error {
	model := emulatorModel{
		ctx:      ctx,
		server:   server,
		snapshot: server.Snapshot(),
		width:    100,
	}
	program := tea.NewProgram(model,
		tea.WithAltScreen(),
		tea.WithInput(input),
		tea.WithOutput(output))
	if _, err := program.Run(); err != nil {
		return fmt.Errorf("run emulator TUI: %w", err)
	}
	return nil
}

func (m emulatorModel) Init() tea.Cmd {
	return scheduleRefresh()
}

func (m emulatorModel) Update(message tea.Msg) (tea.Model, tea.Cmd) {
	switch message := message.(type) {
	case tea.WindowSizeMsg:
		m.width = message.Width
	case tea.KeyMsg:
		if message.Type == tea.KeyCtrlC || message.Type == tea.KeyEsc || message.String() == "q" {
			return m, tea.Quit
		}
	case refreshMessage:
		if m.ctx.Err() != nil {
			return m, tea.Quit
		}
		m.snapshot = m.server.Snapshot()
		return m, scheduleRefresh()
	}
	return m, nil
}

func (m emulatorModel) View() string {
	return renderEmulatorDashboard(m.snapshot, m.width)
}

func scheduleRefresh() tea.Cmd {
	return tea.Tick(refreshTime(), func(now time.Time) tea.Msg {
		return refreshMessage(now)
	})
}
