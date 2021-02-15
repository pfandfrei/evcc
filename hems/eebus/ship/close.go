package ship

import (
	"time"
)

// connection close
const (
	CmiTypeEnd      byte = 3
	CmiCloseTimeout      = 100 * time.Millisecond

	ConnectionCloseReasonUnspecific        = "unspecific"
	ConnectionCloseReasonRemovedConnection = "removedConnection"

	CmiClosePhaseAnnounce = "announce"
	CmiClosePhaseConfirm  = "confirm"
)

// CmiCloseMsg is the close message
type CmiCloseMsg struct {
	ConnectionClose []ConnectionClose `json:"connectionClose"`
}

// ConnectionClose message
type ConnectionClose struct {
	Phase   string `json:"phase"`
	MaxTime int    `json:"maxTime,omitempty"`
	Reason  string `json:"reason,omitempty"`
}
