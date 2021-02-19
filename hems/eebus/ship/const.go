package ship

import "crypto/tls"

// CipherSuites are the SHIP cipher suites
var CipherSuites = []uint16{
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
	tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA256,
}

// protocol constants
const (
	Scheme      = "wss://"
	SubProtocol = "ship"
)
