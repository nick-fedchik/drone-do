package x23wcmd

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/control"
	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/ipc"
)

type BridgeConfig struct {
	TargetIP   string
	ArmLive    bool
	Watchdog   time.Duration
	MaxTilt    float64
	SendRateHz int
	SocketPath string
	StopMode   string
}

// RunBridge exposes the local Unix IPC server and translates incoming frames
// into 20 Hz 0xA5 0x5A packets for 192.168.29.1:25000 when ArmLive is true.
func RunBridge(ctx context.Context, server *ipc.Server, cfg BridgeConfig) error {
	if cfg.TargetIP == "" {
		cfg.TargetIP = "192.168.29.1"
	}
	if cfg.Watchdog <= 0 {
		cfg.Watchdog = 500 * time.Millisecond
	}
	if cfg.MaxTilt <= 0 || cfg.MaxTilt > 1 {
		cfg.MaxTilt = 0.35
	}
	if cfg.SendRateHz <= 0 {
		cfg.SendRateHz = 20
	}

	var ctrlConn, cmdConn *net.UDPConn
	var ctrlAddr, cmdAddr *net.UDPAddr
	if cfg.ArmLive {
		ip := net.ParseIP(cfg.TargetIP)
		if ip == nil {
			return fmt.Errorf("invalid target IP %q", cfg.TargetIP)
		}
		var err error
		ctrlConn, err = net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 20000})
		if err != nil {
			return fmt.Errorf("bind UDP 20000: %w", err)
		}
		defer ctrlConn.Close()

		cmdConn, err = net.ListenUDP("udp4", &net.UDPAddr{IP: net.IPv4zero, Port: 0})
		if err != nil {
			return fmt.Errorf("bind UDP command socket: %w", err)
		}
		defer cmdConn.Close()

		ctrlAddr = &net.UDPAddr{IP: ip, Port: 20000}
		cmdAddr = &net.UDPAddr{IP: ip, Port: FlightUDPPort}

		_, _ = ctrlConn.WriteToUDP([]byte{'J', 'H', 'C', 'M', 'D', 0x10, 0x00}, ctrlAddr)
		time.Sleep(25 * time.Millisecond)
		_, _ = ctrlConn.WriteToUDP([]byte{'J', 'H', 'C', 'M', 'D', 0x20, 0x00, 0x00, 0x00, 0x00, 0x00}, ctrlAddr)
		time.Sleep(15 * time.Millisecond)
		_, _ = ctrlConn.WriteToUDP([]byte{'J', 'H', 'C', 'M', 'D', 0xd0, 0x01}, ctrlAddr)
	}

	tickPeriod := time.Second / time.Duration(cfg.SendRateHz)
	ticker := time.NewTicker(tickPeriod)
	defer ticker.Stop()

	keepalive := time.NewTicker(500 * time.Millisecond)
	defer keepalive.Stop()

	for {
		select {
		case <-ctx.Done():
			if cfg.ArmLive && cmdConn != nil && ctrlConn != nil {
				stopPkt, err := EncodeFlightFrame(control.Frame{At: time.Now().UTC(), Throttle: 0.0}, stopFlags(cfg.StopMode))
				if err == nil {
					for range 5 {
						_, _ = cmdConn.WriteToUDP(stopPkt, cmdAddr)
						time.Sleep(20 * time.Millisecond)
					}
				}
				_, _ = ctrlConn.WriteToUDP([]byte{'J', 'H', 'C', 'M', 'D', 0xd0, 0x02}, ctrlAddr)
			}
			return nil
		case <-keepalive.C:
			if cfg.ArmLive && ctrlConn != nil {
				_, _ = ctrlConn.WriteToUDP([]byte{'J', 'H', 'C', 'M', 'D', 0xd0, 0x01}, ctrlAddr)
			}
		case <-ticker.C:
			pkt, err := PacketFromSnapshot(server.Snapshot(), time.Now().UTC(), cfg.Watchdog, cfg.MaxTilt, cfg.StopMode)
			if err != nil {
				continue
			}
			if cfg.ArmLive && cmdConn != nil {
				_, _ = cmdConn.WriteToUDP(pkt, cmdAddr)
			}
		}
	}
}

func stopFlags(mode string) uint8 {
	switch mode {
	case "land":
		return FlagLand
	case "zero":
		return 0
	default:
		return FlagStop
	}
}

// PacketFromSnapshot applies watchdog timeout and tilt clamping to an IPC snapshot.
func PacketFromSnapshot(snap ipc.Snapshot, now time.Time, watchdog time.Duration, maxTilt float64, stopMode string) ([]byte, error) {
	if snap.LastFrame == nil || snap.UpdatedAt == nil || now.Sub(*snap.UpdatedAt) > watchdog {
		return EncodeFlightFrame(control.Frame{At: now, Throttle: 0.0}, stopFlags(stopMode))
	}
	f := *snap.LastFrame
	f.At = now
	f.Roll = clampAxis(f.Roll, maxTilt)
	f.Pitch = clampAxis(f.Pitch, maxTilt)
	f.Yaw = clampAxis(f.Yaw, maxTilt)
	if f.Throttle <= 0.001 {
		return EncodeFlightFrame(control.Frame{At: now, Throttle: 0.0}, stopFlags(stopMode))
	}
	return EncodeFlightFrame(f, 0)
}

func clampAxis(v, limit float64) float64 {
	if v > limit {
		return limit
	}
	if v < -limit {
		return -limit
	}
	return v
}
