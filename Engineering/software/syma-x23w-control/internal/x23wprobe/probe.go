package x23wprobe

import (
	"encoding/hex"
	"errors"
	"fmt"
	"net"
	"time"
)

const (
	DefaultTargetIP = "192.168.29.1"
	ControlPort     = 20000
	CommandPort     = 20001
	VideoPort       = 10900
)

type Report struct {
	Target          string `json:"target"`
	HandshakeReply  string `json:"handshake_reply_hex,omitempty"`
	HandshakeBytes  int    `json:"handshake_bytes"`
	VideoPackets    int    `json:"video_packets"`
	VideoFirstBytes string `json:"video_first_hex,omitempty"`
	CommandReply    string `json:"command_reply_hex,omitempty"`
	CommandBytes    int    `json:"command_bytes"`
}

// FirmwareVersionQuery is the read-only FDWN version query from Command.java.
func FirmwareVersionQuery() []byte {
	return []byte{'F', 'D', 'W', 'N', 0x20, 0x00, 0x09, 0x00, 0x00, 0x00}
}

// Probe performs non-flight UDP queries against a Joyhonest IC_GPRTPB camera/Wi-Fi board.
func Probe(target string, duration time.Duration) (Report, error) {
	if target == "" {
		target = DefaultTargetIP
	}
	if duration <= 0 {
		duration = 2 * time.Second
	}

	report := Report{Target: target}

	ctrlConn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: ControlPort})
	if err != nil {
		return report, fmt.Errorf("bind UDP %d: %w", ControlPort, err)
	}
	defer ctrlConn.Close()

	cmdConn, err := net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: CommandPort})
	if err != nil {
		return report, fmt.Errorf("bind UDP %d: %w", CommandPort, err)
	}
	defer cmdConn.Close()

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
	cmdAddr := &net.UDPAddr{IP: ip, Port: CommandPort}

	// Stream/handshake packets (no flight commands).
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

	// Read-only Command.java query on UDP 20001.
	if _, err := cmdConn.WriteToUDP(FirmwareVersionQuery(), cmdAddr); err != nil {
		return report, fmt.Errorf("send FDWN version query: %w", err)
	}

	deadline := time.Now().Add(duration)
	buf := make([]byte, 4096)

	_ = ctrlConn.SetReadDeadline(deadline)
	if n, _, err := ctrlConn.ReadFromUDP(buf); err == nil && n > 0 {
		report.HandshakeBytes = n
		preview := min(n, 24)
		report.HandshakeReply = hex.EncodeToString(buf[:preview])
	}

	_ = cmdConn.SetReadDeadline(time.Now().Add(400 * time.Millisecond))
	if n, _, err := cmdConn.ReadFromUDP(buf); err == nil && n > 0 {
		report.CommandBytes = n
		preview := min(n, 32)
		report.CommandReply = hex.EncodeToString(buf[:preview])
	}

	_ = vidConn.SetReadDeadline(deadline)
	for time.Now().Before(deadline) {
		n, _, err := vidConn.ReadFromUDP(buf)
		if err != nil {
			var netErr net.Error
			if errors.As(err, &netErr) && netErr.Timeout() {
				break
			}
			break
		}
		if n > 0 {
			report.VideoPackets++
			if report.VideoFirstBytes == "" {
				preview := min(n, 24)
				report.VideoFirstBytes = hex.EncodeToString(buf[:preview])
			}
		}
	}

	// Stop RTP stream cleanly.
	_, _ = ctrlConn.WriteToUDP([]byte{'J', 'H', 'C', 'M', 'D', 0xd0, 0x02}, ctrlAddr)
	return report, nil
}
