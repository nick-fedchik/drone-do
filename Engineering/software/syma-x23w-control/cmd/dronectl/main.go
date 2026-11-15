package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/control"
	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/ipc"
	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/sequence"
	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/tui"
	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/x23wcmd"
	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/x23wprobe"
)

func main() {
	os.Exit(run(os.Args[1:], os.Stdout, os.Stderr))
}

func run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		printUsage(stderr)
		return 2
	}

	switch args[0] {
	case "protocol":
		if len(args) != 2 || args[1] != "status" {
			printUsage(stderr)
			return 2
		}
		fmt.Fprintln(stdout, "X23W protocol: UNVERIFIED")
		fmt.Fprintln(stdout, "Live radio transport: DISABLED")
		return 0
	case "simulate":
		return simulate(args[1:], stdout, stderr)
	case "status":
		return status(args[1:], stdout, stderr)
	case "send":
		return send(args[1:], stdout, stderr)
	case "tui":
		return controllerTUI(args[1:], stdout, stderr)
	case "probe":
		return probeWiFi(args[1:], stdout, stderr)
	case "video":
		return captureVideo(args[1:], stdout, stderr)
	case "bridge":
		return runBridge(args[1:], stdout, stderr)
	default:
		printUsage(stderr)
		return 2
	}
}

func status(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("status", flag.ContinueOnError)
	flags.SetOutput(stderr)
	socketPath := flags.String("socket", "", "path to local emulator Unix socket")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "status accepts only --socket")
		return 2
	}
	path, err := resolveSocketPath(*socketPath)
	if err != nil {
		fmt.Fprintf(stderr, "resolve socket path: %v\n", err)
		return 1
	}
	client, err := ipc.Dial(path)
	if err != nil {
		fmt.Fprintf(stderr, "connect to emulator: %v\n", err)
		return 1
	}
	defer client.Close()

	snapshot, err := client.Status()
	if err != nil {
		fmt.Fprintf(stderr, "request emulator status: %v\n", err)
		return 1
	}
	if err := json.NewEncoder(stdout).Encode(snapshot); err != nil {
		fmt.Fprintf(stderr, "write status: %v\n", err)
		return 1
	}
	return 0
}

func send(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("send", flag.ContinueOnError)
	flags.SetOutput(stderr)
	socketPath := flags.String("socket", "", "path to local emulator Unix socket")
	throttle := flags.Float64("throttle", 0, "normalized throttle in [0,1]")
	roll := flags.Float64("roll", 0, "normalized roll in [-1,1]")
	pitch := flags.Float64("pitch", 0, "normalized pitch in [-1,1]")
	yaw := flags.Float64("yaw", 0, "normalized yaw in [-1,1]")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "send accepts only --socket and control values")
		return 2
	}
	path, err := resolveSocketPath(*socketPath)
	if err != nil {
		fmt.Fprintf(stderr, "resolve socket path: %v\n", err)
		return 1
	}
	client, err := ipc.Dial(path)
	if err != nil {
		fmt.Fprintf(stderr, "connect to emulator: %v\n", err)
		return 1
	}
	defer client.Close()

	frame := control.Frame{
		At:       time.Now().UTC(),
		Throttle: *throttle,
		Roll:     *roll,
		Pitch:    *pitch,
		Yaw:      *yaw,
	}
	snapshot, err := client.Send(1, frame)
	if err != nil {
		fmt.Fprintf(stderr, "send to emulator: %v\n", err)
		return 1
	}
	if err := json.NewEncoder(stdout).Encode(snapshot); err != nil {
		fmt.Fprintf(stderr, "write acknowledgement: %v\n", err)
		return 1
	}
	return 0
}

func controllerTUI(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("tui", flag.ContinueOnError)
	flags.SetOutput(stderr)
	socketPath := flags.String("socket", "", "path to local emulator Unix socket")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "tui accepts only --socket")
		return 2
	}
	path, err := resolveSocketPath(*socketPath)
	if err != nil {
		fmt.Fprintf(stderr, "resolve socket path: %v\n", err)
		return 1
	}
	client, err := ipc.Dial(path)
	if err != nil {
		fmt.Fprintf(stderr, "connect to emulator: %v\n", err)
		return 1
	}
	defer client.Close()
	if err := tui.RunController(client, os.Stdin, stdout); err != nil {
		fmt.Fprintf(stderr, "controller TUI: %v\n", err)
		return 1
	}
	return 0
}

func resolveSocketPath(path string) (string, error) {
	if path != "" {
		return path, nil
	}
	return ipc.DefaultSocketPath()
}

func probeWiFi(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("probe", flag.ContinueOnError)
	flags.SetOutput(stderr)
	target := flags.String("target", x23wprobe.DefaultTargetIP, "drone Wi-Fi IP address")
	duration := flags.Duration("duration", 2*time.Second, "video/probe listen duration")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "probe accepts only --target and --duration")
		return 2
	}
	report, err := x23wprobe.Probe(*target, *duration)
	if err != nil {
		fmt.Fprintf(stderr, "probe failed: %v\n", err)
		return 1
	}
	if err := json.NewEncoder(stdout).Encode(report); err != nil {
		fmt.Fprintf(stderr, "write probe report: %v\n", err)
		return 1
	}
	return 0
}

func runBridge(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("bridge", flag.ContinueOnError)
	flags.SetOutput(stderr)
	socketPath := flags.String("socket", "", "path to local IPC Unix socket")
	target := flags.String("target", x23wprobe.DefaultTargetIP, "drone Wi-Fi IP address")
	armLive := flags.Bool("arm-live", false, "explicitly enable live UDP 25000 flight packets (default: dry-run)")
	watchdog := flags.Duration("watchdog", 500*time.Millisecond, "IPC watchdog timeout before sending FlagStop")
	maxTilt := flags.Float64("max-tilt", 0.35, "max roll/pitch/yaw amplitude [0,1]")
	stopMode := flags.String("stop-mode", "stop", "stop packet mode on timeout/zero throttle: stop|land|zero")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "bridge accepts only flags")
		return 2
	}
	path, err := resolveSocketPath(*socketPath)
	if err != nil {
		fmt.Fprintf(stderr, "resolve socket path: %v\n", err)
		return 1
	}
	server, err := ipc.Listen(path)
	if err != nil {
		fmt.Fprintf(stderr, "listen IPC socket: %v\n", err)
		return 1
	}
	defer server.Close()

	mode := "DRY-RUN (no UDP flight packets sent)"
	if *armLive {
		mode = fmt.Sprintf("ARMED LIVE -> %s:%d (20 Hz, watchdog=%s, max_tilt=%.2f)", *target, x23wcmd.FlightUDPPort, *watchdog, *maxTilt)
	}
	fmt.Fprintf(stdout, "dronectl bridge listening on %s | mode: %s\n", path, mode)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	serveDone := make(chan error, 1)
	go func() { serveDone <- server.Serve(ctx) }()

	bridgeErr := x23wcmd.RunBridge(ctx, server, x23wcmd.BridgeConfig{
		TargetIP:   *target,
		ArmLive:    *armLive,
		Watchdog:   *watchdog,
		MaxTilt:    *maxTilt,
		SocketPath: path,
		StopMode:   *stopMode,
	})
	stop()
	serveErr := <-serveDone
	if bridgeErr != nil {
		fmt.Fprintf(stderr, "bridge error: %v\n", bridgeErr)
		return 1
	}
	if serveErr != nil {
		fmt.Fprintf(stderr, "IPC server error: %v\n", serveErr)
		return 1
	}
	return 0
}

func captureVideo(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("video", flag.ContinueOnError)
	flags.SetOutput(stderr)
	target := flags.String("target", x23wprobe.DefaultTargetIP, "drone Wi-Fi IP address")
	duration := flags.Duration("duration", 3*time.Second, "capture duration")
	outDir := flags.String("out", "", "optional directory to save decoded JPEG frames")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "video accepts only --target, --duration, and --out")
		return 2
	}
	report, err := x23wprobe.CaptureVideo(*target, *duration, *outDir)
	if err != nil {
		fmt.Fprintf(stderr, "video capture failed: %v\n", err)
		return 1
	}
	if err := json.NewEncoder(stdout).Encode(report); err != nil {
		fmt.Fprintf(stderr, "write video report: %v\n", err)
		return 1
	}
	return 0
}

func simulate(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("simulate", flag.ContinueOnError)
	flags.SetOutput(stderr)
	inputPath := flags.String("input", "", "path to a local JSONL control sequence")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	if *inputPath == "" || flags.NArg() != 0 {
		fmt.Fprintln(stderr, "simulate requires --input FILE")
		return 2
	}

	file, err := os.Open(*inputPath)
	if err != nil {
		fmt.Fprintf(stderr, "open input: %v\n", err)
		return 1
	}
	defer file.Close()

	frames, err := sequence.Read(file)
	if err != nil {
		fmt.Fprintf(stderr, "read input: %v\n", err)
		return 1
	}

	encoder := json.NewEncoder(stdout)
	for i, frame := range frames {
		result := struct {
			Mode  string `json:"mode"`
			Index int    `json:"index"`
			Frame any    `json:"frame"`
		}{
			Mode:  "SIMULATION_ONLY",
			Index: i + 1,
			Frame: frame,
		}
		if err := encoder.Encode(result); err != nil {
			fmt.Fprintf(stderr, "write simulation output: %v\n", err)
			return 1
		}
	}
	return 0
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  dronectl protocol status")
	fmt.Fprintln(w, "  dronectl simulate --input FILE")
	fmt.Fprintln(w, "  dronectl status [--socket PATH]")
	fmt.Fprintln(w, "  dronectl send [--socket PATH] [--throttle N] [--roll N] [--pitch N] [--yaw N]")
	fmt.Fprintln(w, "  dronectl tui [--socket PATH]")
	fmt.Fprintln(w, "  dronectl probe [--target IP] [--duration 2s]")
	fmt.Fprintln(w, "  dronectl video [--target IP] [--duration 3s] [--out DIR]")
	fmt.Fprintln(w, "  dronectl bridge [--socket PATH] [--target IP] [--arm-live] [--watchdog 500ms] [--max-tilt 0.35] [--stop-mode stop|land|zero]")
	fmt.Fprintln(w, "status/send/tui use the local pseudo-emulator; bridge defaults to dry-run unless --arm-live is set.")
}
