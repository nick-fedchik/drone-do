package tui

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/control"
	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/ipc"
)

const mockBatteryPercent = 100

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("45"))
	panelTitle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("81"))
	healthy    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("42"))
	warning    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("214"))
	dim        = lipgloss.NewStyle().Foreground(lipgloss.Color("244"))
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("39")).
			Padding(0, 1)
)

func renderEmulatorDashboard(snapshot ipc.Snapshot, width int) string {
	width = dashboardWidth(width)
	contentWidth := width - 2
	frame := control.Frame{}
	if snapshot.LastFrame != nil {
		frame = *snapshot.LastFrame
	}

	headerText := "X23W / EMULATOR"
	if width >= 80 {
		headerText = "SYMA X23W  /  LOCAL PSEUDO-EMULATOR"
	}
	header := titleStyle.Render(headerText) + "   " + statusBadge(snapshot.Running)
	airframe := renderPanel("AIRFRAME", airframeRows(snapshot), minInt(28, contentWidth))
	battery := renderPanel("BATTERY", batteryRows(), minInt(28, contentWidth))
	commands := renderPanel("RECENT COMMANDS", commandRows(snapshot, contentWidth), contentWidth)

	var layout string
	switch {
	case width >= 100:
		leftWidth, centerWidth := 27, 35
		rightWidth := contentWidth - leftWidth - centerWidth - 2
		left := lipgloss.JoinVertical(lipgloss.Left,
			renderPanel("AIRFRAME", airframeRows(snapshot), leftWidth),
			renderPanel("BATTERY", batteryRows(), leftWidth))
		layout = lipgloss.JoinHorizontal(lipgloss.Top, left, " ",
			renderPanel("CONTROL INPUT", controlRows(frame, centerWidth), centerWidth), " ",
			renderPanel("RECENT COMMANDS", commandRows(snapshot, rightWidth), rightWidth))
	case width >= 72:
		leftWidth := contentWidth / 2
		rightWidth := contentWidth - leftWidth - 1
		left := lipgloss.JoinVertical(lipgloss.Left, airframe, battery)
		top := lipgloss.JoinHorizontal(lipgloss.Top, left, " ",
			renderPanel("CONTROL INPUT", controlRows(frame, rightWidth), rightWidth))
		layout = lipgloss.JoinVertical(lipgloss.Left, top, commands)
	default:
		layout = lipgloss.JoinVertical(lipgloss.Left,
			renderPanel("AIRFRAME", airframeRows(snapshot), contentWidth),
			renderPanel("BATTERY", batteryRows(), contentWidth),
			renderPanel("CONTROL INPUT", controlRows(frame, contentWidth), contentWidth),
			renderPanel("RECENT COMMANDS", commandRows(snapshot, contentWidth), contentWidth))
	}

	footerText := "q quit | refresh 250 ms | no flight model or aircraft telemetry"
	if width < 72 {
		footerText = "q quit | mock telemetry only"
	}
	footer := dim.Render(fitText(footerText, contentWidth))
	return lipgloss.JoinVertical(lipgloss.Left, header, "", layout, "", footer)
}

func renderControllerDashboard(frame control.Frame, snapshot ipc.Snapshot, message string, width int) string {
	width = dashboardWidth(width)
	contentWidth := width - 2
	header := titleStyle.Render("X23W CONTROL  /  LOCAL EMULATOR ONLY") + "   " + statusBadge(snapshot.Running)
	controls := renderPanel("CONTROL INPUT", controlRows(frame, contentWidth), contentWidth)
	statusRows := []string{
		"IPC      " + statusText(snapshot.Running),
		fmt.Sprintf("COMMANDS %d", snapshot.CommandsReceived),
		"PROFILE  " + maneuverLabel(&frame),
		"RADIO    DISABLED",
		"PROTOCOL UNVERIFIED",
	}
	if snapshot.UpdatedAt != nil {
		statusRows = append(statusRows, "UPDATED  "+snapshot.UpdatedAt.Format("15:04:05 UTC"))
	}
	status := renderPanel("EMULATOR STATUS", statusRows, contentWidth)
	help := "r/f thr | a/d roll | w/s pitch\nq/e yaw | 0 neutral | x exit"
	if message != "" {
		help += "\n" + warning.Render(message)
	}
	helpPanel := renderPanel("KEYS", []string{help}, contentWidth)
	if width >= 82 {
		leftWidth := contentWidth / 2
		rightWidth := contentWidth - leftWidth - 1
		return lipgloss.JoinVertical(lipgloss.Left, header, "",
			lipgloss.JoinHorizontal(lipgloss.Top,
				renderPanel("EMULATOR STATUS", statusRows, leftWidth), " ",
				renderPanel("CONTROL INPUT", controlRows(frame, rightWidth), rightWidth)),
			helpPanel)
	}
	return lipgloss.JoinVertical(lipgloss.Left, header, "", status, controls, helpPanel)
}

func renderPanel(title string, rows []string, width int) string {
	innerWidth := maxInt(8, width-4)
	content := panelTitle.Render(title)
	if len(rows) > 0 {
		content += "\n" + strings.Join(rows, "\n")
	}
	return panelStyle.Copy().Width(innerWidth).Render(content)
}

func airframeRows(snapshot ipc.Snapshot) []string {
	return []string{
		"MODEL   X23W / EMULATED",
		"STATE   " + statusText(snapshot.Running),
		"IPC     LOCAL UNIX",
		"RADIO   DISABLED",
		"PROTOCOL UNVERIFIED",
		"AIRFRAME DATA UNKNOWN",
		"PHYSICS NOT MODELLED",
	}
}

func batteryRows() []string {
	return []string{
		fmt.Sprintf("MOCK  %3d%%", mockBatteryPercent),
		"############",
		"fixed demo value",
		"no drain model",
	}
}

func controlRows(frame control.Frame, width int) []string {
	meterWidth := maxInt(20, width-4)
	return []string{
		"PROFILE  " + maneuverLabel(&frame),
		meter("THR", frame.Throttle, 0, 1, meterWidth, false),
		meter("ROLL", frame.Roll, -1, 1, meterWidth, true),
		meter("PITCH", frame.Pitch, -1, 1, meterWidth, true),
		meter("YAW", frame.Yaw, -1, 1, meterWidth, true),
		dim.Render("input only; no motion model"),
	}
}

func commandRows(snapshot ipc.Snapshot, width int) []string {
	rows := []string{fmt.Sprintf("TOTAL  %d accepted", snapshot.CommandsReceived)}
	if len(snapshot.RecentCommands) == 0 {
		rows = append(rows, dim.Render("waiting for first command"))
		return rows
	}
	for i := len(snapshot.RecentCommands) - 1; i >= 0; i-- {
		command := snapshot.RecentCommands[i]
		line := fmt.Sprintf("#%04d  %-14s %s",
			command.Number,
			maneuverLabel(&command.Frame),
			command.Frame.At.Local().Format("15:04:05"))
		rows = append(rows, fitText(line, maxInt(12, width-6)))
	}
	return rows
}

func statusBadge(running bool) string {
	if running {
		return healthy.Render("[ ONLINE ]")
	}
	return warning.Render("[ OFFLINE ]")
}

func statusText(running bool) string {
	if running {
		return healthy.Render("ONLINE")
	}
	return warning.Render("OFFLINE")
}

func maneuverLabel(frame *control.Frame) string {
	if frame == nil {
		return "WAITING FOR INPUT"
	}
	active := 0
	var direction string
	for _, axis := range []struct {
		name  string
		value float64
	}{
		{"ROLL", frame.Roll},
		{"PITCH", frame.Pitch},
		{"YAW", frame.Yaw},
	} {
		if math.Abs(axis.value) > 0.001 {
			active++
			direction = fmt.Sprintf("%s %+.1f", axis.name, axis.value)
		}
	}
	if active == 0 {
		if frame.Throttle > 0.001 {
			return "THROTTLE ONLY"
		}
		return "NEUTRAL INPUT"
	}
	if active > 1 || frame.Throttle > 0.001 {
		return "COMBINED INPUT"
	}
	return direction
}

func meter(label string, value, min, max float64, width int, signed bool) string {
	valueText := fmt.Sprintf("%.2f", value)
	if signed {
		valueText = fmt.Sprintf("%+.2f", value)
	}
	barWidth := maxInt(4, width-len(label)-len(valueText)-5)
	ratio := (value - min) / (max - min)
	ratio = math.Max(0, math.Min(1, ratio))
	filled := int(math.Round(ratio * float64(barWidth)))
	bar := strings.Repeat("#", filled) + strings.Repeat("-", barWidth-filled)
	return fmt.Sprintf("%-5s [%s] %s", label, bar, valueText)
}

func dashboardWidth(width int) int {
	if width <= 0 {
		return 100
	}
	return maxInt(48, width)
}

func fitText(text string, width int) string {
	runes := []rune(text)
	if len(runes) <= width {
		return text
	}
	if width <= 1 {
		return string(runes[:width])
	}
	return string(runes[:width-1]) + "~"
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func refreshTime() time.Duration {
	return 250 * time.Millisecond
}
