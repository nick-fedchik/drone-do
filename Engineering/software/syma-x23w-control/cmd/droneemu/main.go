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

	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/ipc"
	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/tui"
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
	case "serve":
		return serve(args[1:], stdout, stderr)
	case "tui":
		return runTUI(args[1:], stdout, stderr)
	default:
		printUsage(stderr)
		return 2
	}
}

func serve(args []string, stdout, stderr io.Writer) int {
	path, err := parseSocketPath("serve", args, stderr)
	if err != nil {
		return 2
	}
	server, err := ipc.Listen(path)
	if err != nil {
		fmt.Fprintf(stderr, "start emulator: %v\n", err)
		return 1
	}
	defer server.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	serveDone := make(chan error, 1)
	go func() { serveDone <- server.Serve(ctx) }()
	fmt.Fprintf(stdout, "Pseudo-drone listening on %s (local Unix socket; no flight physics)\n", path)
	for {
		select {
		case snapshot := <-server.Updates():
			if err := json.NewEncoder(stdout).Encode(snapshot); err != nil {
				fmt.Fprintf(stderr, "write emulator event: %v\n", err)
				return 1
			}
		case err := <-serveDone:
			if err != nil {
				fmt.Fprintf(stderr, "emulator server: %v\n", err)
				return 1
			}
			return 0
		case <-ctx.Done():
			if err := <-serveDone; err != nil {
				fmt.Fprintf(stderr, "emulator server: %v\n", err)
				return 1
			}
			return 0
		}
	}
}

func runTUI(args []string, stdout, stderr io.Writer) int {
	path, err := parseSocketPath("tui", args, stderr)
	if err != nil {
		return 2
	}
	server, err := ipc.Listen(path)
	if err != nil {
		fmt.Fprintf(stderr, "start emulator: %v\n", err)
		return 1
	}
	defer server.Close()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	serveDone := make(chan error, 1)
	go func() { serveDone <- server.Serve(ctx) }()
	tuiErr := tui.RunEmulator(ctx, server, os.Stdin, stdout)
	stop()
	serveErr := <-serveDone
	if tuiErr != nil {
		fmt.Fprintf(stderr, "emulator TUI: %v\n", tuiErr)
		return 1
	}
	if serveErr != nil {
		fmt.Fprintf(stderr, "emulator server: %v\n", serveErr)
		return 1
	}
	return 0
}

func parseSocketPath(command string, args []string, stderr io.Writer) (string, error) {
	flags := flag.NewFlagSet(command, flag.ContinueOnError)
	flags.SetOutput(stderr)
	path := flags.String("socket", "", "path to local Unix socket")
	if err := flags.Parse(args); err != nil {
		return "", err
	}
	if flags.NArg() != 0 {
		fmt.Fprintf(stderr, "%s accepts only --socket\n", command)
		return "", fmt.Errorf("unexpected positional argument")
	}
	if *path != "" {
		return *path, nil
	}
	resolved, err := ipc.DefaultSocketPath()
	if err != nil {
		fmt.Fprintf(stderr, "resolve socket path: %v\n", err)
		return "", err
	}
	return resolved, nil
}

func printUsage(w io.Writer) {
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, "  droneemu serve [--socket PATH]")
	fmt.Fprintln(w, "  droneemu tui [--socket PATH]")
	fmt.Fprintln(w, "The emulator is local-only and does not model flight physics.")
}
