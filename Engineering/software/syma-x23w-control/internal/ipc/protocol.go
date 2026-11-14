package ipc

import (
	"time"

	"github.com/nick-fedchik/drone-do/Engineering/software/syma-x23w-control/internal/control"
)

const ProtocolVersion = 1
const recentCommandLimit = 6

type Request struct {
	Version  int            `json:"version"`
	Type     string         `json:"type"`
	Sequence uint64         `json:"sequence,omitempty"`
	Frame    *control.Frame `json:"frame,omitempty"`
}

type Snapshot struct {
	Running          bool            `json:"running"`
	CommandsReceived uint64          `json:"commands_received"`
	LastSequence     uint64          `json:"last_sequence,omitempty"`
	LastFrame        *control.Frame  `json:"last_frame,omitempty"`
	UpdatedAt        *time.Time      `json:"updated_at,omitempty"`
	RecentCommands   []CommandRecord `json:"recent_commands,omitempty"`
}

type CommandRecord struct {
	Number uint64        `json:"number"`
	Frame  control.Frame `json:"frame"`
}

type Response struct {
	Version  int      `json:"version"`
	Type     string   `json:"type"`
	Sequence uint64   `json:"sequence,omitempty"`
	Snapshot Snapshot `json:"snapshot"`
	Error    string   `json:"error,omitempty"`
}
