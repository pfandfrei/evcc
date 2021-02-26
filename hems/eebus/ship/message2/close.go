package message

import (
	"time"

	"github.com/andig/evcc/hems/eebus/util"
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
	ConnectionClose ConnectionClose `json:"connectionClose"`
}

// ConnectionClose message
type ConnectionClose struct {
	Phase   string `json:"phase"`
	MaxTime int    `json:"maxTime,omitempty"`
	Reason  string `json:"reason,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m ConnectionClose) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *ConnectionClose) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}
