package x23wprobe

import (
	"bytes"
	"errors"
	"fmt"
	"image/jpeg"
	"net"
	"os"
	"path/filepath"
	"time"
)

const rtpHeaderSize = 8

type VideoReport struct {
	Target        string  `json:"target"`
	DurationSec   float64 `json:"duration_sec"`
	Packets       int     `json:"packets"`
	FramesDecoded int     `json:"frames_decoded"`
	FPS           float64 `json:"fps"`
	Width         int     `json:"width,omitempty"`
	Height        int     `json:"height,omitempty"`
	SavedDir      string  `json:"saved_dir,omitempty"`
	LastFrameFile string  `json:"last_frame_file,omitempty"`
}

// ReassembleJPEGs extracts complete JPEG images from Joyhonest 8-byte-prefixed UDP payloads.
func ReassembleJPEGs(packets [][]byte) [][]byte {
	var frames [][]byte
	var current bytes.Buffer
	inFrame := false

	for _, pkt := range packets {
		if len(pkt) <= rtpHeaderSize {
			continue
		}
		payload := pkt[rtpHeaderSize:]
		if len(payload) >= 2 && payload[0] == 0xFF && payload[1] == 0xD8 {
			current.Reset()
			inFrame = true
		}
		if !inFrame {
			continue
		}
		current.Write(payload)
		if idx := bytes.Index(payload, []byte{0xFF, 0xD9}); idx >= 0 {
			data := current.Bytes()
			if eoi := bytes.LastIndex(data, []byte{0xFF, 0xD9}); eoi >= 0 {
				frame := make([]byte, eoi+2)
				copy(frame, data[:eoi+2])
				frames = append(frames, frame)
			}
			current.Reset()
			inFrame = false
		}
	}
	return frames
}

func CaptureVideo(target string, duration time.Duration, outDir string) (VideoReport, error) {
	if target == "" {
		target = DefaultTargetIP
	}
	if duration <= 0 {
		duration = 3 * time.Second
	}

	report := VideoReport{
		Target:      target,
		DurationSec: duration.Seconds(),
		SavedDir:    outDir,
	}

	ctrlConn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: ControlPort})
	if err != nil {
		return report, fmt.Errorf("bind UDP %d: %w", ControlPort, err)
	}
	defer ctrlConn.Close()

	vidConn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: VideoPort})
	if err != nil {
		return report, fmt.Errorf("bind UDP %d: %w", VideoPort, err)
	}
	defer vidConn.Close()

	ip := net.ParseIP(target)
	if ip == nil {
		return report, fmt.Errorf("invalid target IP %q", target)
	}
	ctrlAddr := &net.UDPAddr{IP: ip, Port: ControlPort}

	if _, err := ctrlConn.WriteToUDP([]byte{'J', 'H', 'C', 'M', 'D', 0x10, 0x00}, ctrlAddr); err != nil {
		return report, fmt.Errorf("send JHCMD 0x10: %w", err)
	}
	time.Sleep(25 * time.Millisecond)
	if _, err := ctrlConn.WriteToUDP([]byte{'J', 'H', 'C', 'M', 'D', 0x20, 0x00, 0x00, 0x00, 0x00, 0x00}, ctrlAddr); err != nil {
		return report, fmt.Errorf("send JHCMD 0x20: %w", err)
	}
	time.Sleep(15 * time.Millisecond)
	if _, err := ctrlConn.WriteToUDP([]byte{'J', 'H', 'C', 'M', 'D', 0xd0, 0x01}, ctrlAddr); err != nil {
		return report, fmt.Errorf("send JHCMD 0xd0: %w", err)
	}

	var packets [][]byte
	buf := make([]byte, 4096)
	deadline := time.Now().Add(duration)
	nextKeepalive := time.Now().Add(500 * time.Millisecond)
	_ = vidConn.SetReadDeadline(deadline)

	for time.Now().Before(deadline) {
		if time.Now().After(nextKeepalive) {
			_, _ = ctrlConn.WriteToUDP([]byte{'J', 'H', 'C', 'M', 'D', 0xd0, 0x01}, ctrlAddr)
			nextKeepalive = time.Now().Add(500 * time.Millisecond)
		}
		n, _, err := vidConn.ReadFromUDP(buf)
		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				break
			}
			break
		}
		if n > 0 {
			pkt := make([]byte, n)
			copy(pkt, buf[:n])
			packets = append(packets, pkt)
		}
	}

	_, _ = ctrlConn.WriteToUDP([]byte{'J', 'H', 'C', 'M', 'D', 0xd0, 0x02}, ctrlAddr)
	report.Packets = len(packets)

	if outDir != "" {
		if err := os.MkdirAll(outDir, 0o755); err != nil {
			return report, fmt.Errorf("create output directory: %w", err)
		}
	}

	for _, frame := range ReassembleJPEGs(packets) {
		cfg, err := jpeg.DecodeConfig(bytes.NewReader(frame))
		if err != nil {
			continue
		}
		report.FramesDecoded++
		report.Width = cfg.Width
		report.Height = cfg.Height
		if outDir != "" {
			name := filepath.Join(outDir, fmt.Sprintf("frame_%04d.jpg", report.FramesDecoded))
			if err := os.WriteFile(name, frame, 0o644); err != nil {
				return report, fmt.Errorf("write frame: %w", err)
			}
			report.LastFrameFile = name
		}
	}
	if report.DurationSec > 0 {
		report.FPS = float64(report.FramesDecoded) / report.DurationSec
	}
	return report, nil
}
