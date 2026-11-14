package ipc

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/control"
)

func TestLocalIPCCommandAndStatus(t *testing.T) {
	path := filepath.Join(t.TempDir(), "emulator.sock")
	server, err := Listen(path)
	if err != nil {
		t.Fatal(err)
	}
	mode, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if mode.Mode().Perm() != 0o600 {
		t.Fatalf("socket permissions = %04o, want 0600", mode.Mode().Perm())
	}

	ctx, cancel := context.WithCancel(context.Background())
	serveDone := make(chan error, 1)
	go func() { serveDone <- server.Serve(ctx) }()
	defer func() {
		cancel()
		if err := server.Close(); err != nil {
			t.Errorf("close server: %v", err)
		}
		if err := <-serveDone; err != nil {
			t.Errorf("serve: %v", err)
		}
	}()

	client, err := dialWithRetry(path)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()

	initial, err := client.Status()
	if err != nil {
		t.Fatal(err)
	}
	if !initial.Running || initial.CommandsReceived != 0 || initial.UpdatedAt != nil {
		t.Fatalf("initial snapshot = %+v", initial)
	}

	frame := control.Frame{
		At:       time.Now().UTC(),
		Throttle: 0.25,
		Roll:     -0.5,
		Pitch:    0.75,
		Yaw:      0.1,
	}
	var ack Snapshot
	for sequence := uint64(1); sequence <= uint64(recentCommandLimit+2); sequence++ {
		frame.At = frame.At.Add(time.Second)
		ack, err = client.Send(sequence, frame)
		if err != nil {
			t.Fatal(err)
		}
	}
	if ack.CommandsReceived != uint64(recentCommandLimit+2) || ack.LastSequence != uint64(recentCommandLimit+2) || ack.LastFrame == nil {
		t.Fatalf("acknowledgement snapshot = %+v", ack)
	}
	if *ack.LastFrame != frame {
		t.Fatalf("last frame = %+v, want %+v", *ack.LastFrame, frame)
	}
	if len(ack.RecentCommands) != recentCommandLimit || ack.RecentCommands[0].Number != 3 {
		t.Fatalf("recent command history = %+v", ack.RecentCommands)
	}

	status, err := client.Status()
	if err != nil {
		t.Fatal(err)
	}
	if status.CommandsReceived != uint64(recentCommandLimit+2) || status.LastFrame == nil || *status.LastFrame != frame {
		t.Fatalf("status snapshot = %+v", status)
	}
}

func TestListenPreservesNonSocketPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "occupied")
	if err := os.WriteFile(path, []byte("keep"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := Listen(path); err == nil {
		t.Fatal("Listen succeeded for a regular file")
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "keep" {
		t.Fatalf("path content = %q, want preserved file", got)
	}
}

func TestServerRejectsUnsupportedProtocolVersion(t *testing.T) {
	path := filepath.Join(t.TempDir(), "emulator.sock")
	server, err := Listen(path)
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	serveDone := make(chan error, 1)
	go func() { serveDone <- server.Serve(ctx) }()
	defer func() {
		cancel()
		_ = server.Close()
		if err := <-serveDone; err != nil {
			t.Errorf("serve: %v", err)
		}
	}()

	conn, err := dialUnixWithRetry(path)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()
	if err := json.NewEncoder(conn).Encode(Request{Version: ProtocolVersion + 1, Type: "status"}); err != nil {
		t.Fatal(err)
	}
	var response Response
	if err := json.NewDecoder(conn).Decode(&response); err != nil {
		t.Fatal(err)
	}
	if response.Type != "error" {
		t.Fatalf("response type = %q, want error", response.Type)
	}
}

func dialWithRetry(path string) (*Client, error) {
	var lastErr error
	for range 50 {
		client, err := Dial(path)
		if err == nil {
			return client, nil
		}
		lastErr = err
		time.Sleep(10 * time.Millisecond)
	}
	return nil, lastErr
}

func dialUnixWithRetry(path string) (*net.UnixConn, error) {
	var lastErr error
	for range 50 {
		conn, err := net.DialUnix("unix", nil, &net.UnixAddr{Name: path, Net: "unix"})
		if err == nil {
			return conn, nil
		}
		lastErr = err
		time.Sleep(10 * time.Millisecond)
	}
	return nil, lastErr
}
