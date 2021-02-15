package ship

import (
	"time"
)

// init constants
const (
	CmiTypeInit byte = 0

	CmiHelloInitTimeout         = 60 * time.Second
	CmiHelloProlongationTimeout = 30 * time.Second

	CmiHelloPhasePending = "pending"
	CmiHelloPhaseReady   = "ready"
	CmiHelloPhaseAborted = "aborted"
)

type CmiHelloMsg struct {
	ConnectionHello []ConnectionHello `json:"connectionHello"`
}

type ConnectionHello struct {
	Phase               string `json:"phase"`
	Waiting             int    `json:"waiting,omitempty"`
	ProlongationRequest bool   `json:"prolongationRequest,omitempty"`
}
