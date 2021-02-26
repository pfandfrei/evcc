package message

import (
	"time"

	"github.com/andig/evcc/hems/eebus/util"
)

// init constants
const (
	CmiTypeInit byte = 0

	CmiTimeout                  = 60 * time.Second
	CmiHelloProlongationTimeout = 30 * time.Second

	CmiHelloPhasePending = "pending"
	CmiHelloPhaseReady   = "ready"
	CmiHelloPhaseAborted = "aborted"
)

type CmiHelloMsg struct {
	ConnectionHello ConnectionHello `json:"connectionHello"`
}

type ConnectionHello struct {
	Phase               string `json:"phase"`
	Waiting             int    `json:"waiting,omitempty"`
	ProlongationRequest bool   `json:"prolongationRequest,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m ConnectionHello) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *ConnectionHello) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}
