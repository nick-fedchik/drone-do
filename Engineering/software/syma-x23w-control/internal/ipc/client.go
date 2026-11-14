package ipc

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/control"
)

const requestTimeout = 3 * time.Second

type Client struct {
	conn    *net.UnixConn
	encoder *json.Encoder
	decoder *json.Decoder
}

func DefaultSocketPath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("resolve user cache directory: %w", err)
	}
	return filepath.Join(cacheDir, "drone-do", "syma-x23w.sock"), nil
}

func Dial(path string) (*Client, error) {
	address := &net.UnixAddr{Name: path, Net: "unix"}
	conn, err := net.DialUnix("unix", nil, address)
	if err != nil {
		return nil, fmt.Errorf("connect to emulator socket: %w", err)
	}
	return &Client{
		conn:    conn,
		encoder: json.NewEncoder(conn),
		decoder: json.NewDecoder(conn),
	}, nil
}

func (c *Client) Status() (Snapshot, error) {
	response, err := c.roundTrip(Request{Version: ProtocolVersion, Type: "status"})
	if err != nil {
		return Snapshot{}, err
	}
	if response.Type != "status" {
		return Snapshot{}, fmt.Errorf("unexpected response type %q", response.Type)
	}
	return response.Snapshot, nil
}

func (c *Client) Send(sequence uint64, frame control.Frame) (Snapshot, error) {
	if sequence == 0 {
		return Snapshot{}, fmt.Errorf("sequence must be greater than zero")
	}
	if err := frame.Validate(); err != nil {
		return Snapshot{}, fmt.Errorf("invalid control frame: %w", err)
	}
	response, err := c.roundTrip(Request{
		Version:  ProtocolVersion,
		Type:     "command",
		Sequence: sequence,
		Frame:    &frame,
	})
	if err != nil {
		return Snapshot{}, err
	}
	if response.Type != "ack" {
		return Snapshot{}, fmt.Errorf("unexpected response type %q", response.Type)
	}
	return response.Snapshot, nil
}

func (c *Client) roundTrip(request Request) (Response, error) {
	if err := c.conn.SetDeadline(time.Now().Add(requestTimeout)); err != nil {
		return Response{}, fmt.Errorf("set IPC deadline: %w", err)
	}
	if err := c.encoder.Encode(request); err != nil {
		return Response{}, fmt.Errorf("write IPC request: %w", err)
	}

	var response Response
	if err := c.decoder.Decode(&response); err != nil {
		return Response{}, fmt.Errorf("read IPC response: %w", err)
	}
	if response.Version != ProtocolVersion {
		return Response{}, fmt.Errorf("unsupported IPC response version %d", response.Version)
	}
	if response.Type == "error" {
		return Response{}, fmt.Errorf("emulator rejected request: %s", response.Error)
	}
	return response, nil
}

func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	if err := c.conn.Close(); err != nil {
		return fmt.Errorf("close IPC connection: %w", err)
	}
	return nil
}
