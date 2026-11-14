package ipc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type Server struct {
	listener *net.UnixListener
	path     string

	mu       sync.RWMutex
	snapshot Snapshot

	connMu   sync.Mutex
	conns    map[*net.UnixConn]struct{}
	handlers sync.WaitGroup

	updates chan Snapshot
	close   sync.Once
}

func Listen(path string) (*Server, error) {
	if path == "" {
		return nil, fmt.Errorf("socket path is required")
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return nil, fmt.Errorf("create socket directory: %w", err)
	}
	if info, err := os.Lstat(path); err == nil {
		if info.Mode()&os.ModeSocket == 0 {
			return nil, fmt.Errorf("socket path exists and is not a Unix socket")
		}
		return nil, fmt.Errorf("socket path already exists; verify no emulator is running before removing it")
	} else if !errors.Is(err, os.ErrNotExist) {
		return nil, fmt.Errorf("inspect socket path: %w", err)
	}

	listener, err := net.ListenUnix("unix", &net.UnixAddr{Name: path, Net: "unix"})
	if err != nil {
		return nil, fmt.Errorf("listen on Unix socket: %w", err)
	}
	listener.SetUnlinkOnClose(true)
	if err := os.Chmod(path, 0o600); err != nil {
		_ = listener.Close()
		return nil, fmt.Errorf("restrict socket permissions: %w", err)
	}

	return &Server{
		listener: listener,
		path:     path,
		snapshot: Snapshot{Running: true},
		conns:    make(map[*net.UnixConn]struct{}),
		updates:  make(chan Snapshot, 1),
	}, nil
}

func (s *Server) Path() string {
	return s.path
}

func (s *Server) Updates() <-chan Snapshot {
	return s.updates
}

func (s *Server) Snapshot() Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return cloneSnapshot(s.snapshot)
}

func (s *Server) Serve(ctx context.Context) error {
	defer s.handlers.Wait()
	defer s.closeConnections()
	for {
		if ctx.Err() != nil {
			return nil
		}
		if err := s.listener.SetDeadline(time.Now().Add(250 * time.Millisecond)); err != nil {
			if ctx.Err() != nil || errors.Is(err, net.ErrClosed) {
				return nil
			}
			return fmt.Errorf("set listener deadline: %w", err)
		}
		conn, err := s.listener.AcceptUnix()
		if err != nil {
			if ctx.Err() != nil || errors.Is(err, net.ErrClosed) {
				return nil
			}
			var networkError net.Error
			if errors.As(err, &networkError) && networkError.Timeout() {
				continue
			}
			return fmt.Errorf("accept IPC connection: %w", err)
		}

		s.connMu.Lock()
		s.conns[conn] = struct{}{}
		s.handlers.Add(1)
		s.connMu.Unlock()
		go s.handle(conn)
	}
}

func (s *Server) Close() error {
	var closeErr error
	s.close.Do(func() {
		if err := s.listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
			closeErr = fmt.Errorf("close Unix listener: %w", err)
		}
		s.closeConnections()
	})
	return closeErr
}

func (s *Server) handle(conn *net.UnixConn) {
	defer s.handlers.Done()
	defer func() {
		_ = conn.Close()
		s.connMu.Lock()
		delete(s.conns, conn)
		s.connMu.Unlock()
	}()

	decoder := json.NewDecoder(conn)
	decoder.DisallowUnknownFields()
	encoder := json.NewEncoder(conn)
	for {
		var request Request
		if err := decoder.Decode(&request); err != nil {
			return
		}
		response := s.process(request)
		if err := encoder.Encode(response); err != nil {
			return
		}
	}
}

func (s *Server) process(request Request) Response {
	response := Response{
		Version:  ProtocolVersion,
		Sequence: request.Sequence,
	}
	if request.Version != ProtocolVersion {
		response.Type = "error"
		response.Error = fmt.Sprintf("unsupported IPC version %d", request.Version)
		response.Snapshot = s.Snapshot()
		return response
	}

	switch request.Type {
	case "status":
		if request.Frame != nil {
			response.Type = "error"
			response.Error = "status request must not include a frame"
			break
		}
		response.Type = "status"
	case "command":
		if request.Sequence == 0 {
			response.Type = "error"
			response.Error = "command sequence must be greater than zero"
			break
		}
		if request.Frame == nil {
			response.Type = "error"
			response.Error = "command frame is required"
			break
		}
		if err := request.Frame.Validate(); err != nil {
			response.Type = "error"
			response.Error = err.Error()
			break
		}
		s.mu.Lock()
		if s.snapshot.LastFrame != nil && !request.Frame.At.After(s.snapshot.LastFrame.At) {
			s.mu.Unlock()
			response.Type = "error"
			response.Error = "command timestamps must be strictly increasing"
			break
		}
		frame := *request.Frame
		s.snapshot.CommandsReceived++
		s.snapshot.LastSequence = request.Sequence
		s.snapshot.LastFrame = &frame
		updatedAt := time.Now().UTC()
		s.snapshot.UpdatedAt = &updatedAt
		s.snapshot.RecentCommands = append(s.snapshot.RecentCommands, CommandRecord{
			Number: s.snapshot.CommandsReceived,
			Frame:  frame,
		})
		if len(s.snapshot.RecentCommands) > recentCommandLimit {
			s.snapshot.RecentCommands = append([]CommandRecord(nil), s.snapshot.RecentCommands[len(s.snapshot.RecentCommands)-recentCommandLimit:]...)
		}
		snapshot := cloneSnapshot(s.snapshot)
		s.mu.Unlock()
		response.Type = "ack"
		response.Snapshot = snapshot
		s.publish(snapshot)
		return response
	default:
		response.Type = "error"
		response.Error = fmt.Sprintf("unsupported request type %q", request.Type)
	}
	response.Snapshot = s.Snapshot()
	return response
}

func (s *Server) publish(snapshot Snapshot) {
	select {
	case s.updates <- snapshot:
	default:
		select {
		case <-s.updates:
		default:
		}
		select {
		case s.updates <- snapshot:
		default:
		}
	}
}

func (s *Server) closeConnections() {
	s.connMu.Lock()
	defer s.connMu.Unlock()
	for conn := range s.conns {
		_ = conn.Close()
	}
}

func cloneSnapshot(snapshot Snapshot) Snapshot {
	if snapshot.LastFrame != nil {
		frame := *snapshot.LastFrame
		snapshot.LastFrame = &frame
	}
	if snapshot.UpdatedAt != nil {
		updatedAt := *snapshot.UpdatedAt
		snapshot.UpdatedAt = &updatedAt
	}
	snapshot.RecentCommands = append([]CommandRecord(nil), snapshot.RecentCommands...)
	return snapshot
}
