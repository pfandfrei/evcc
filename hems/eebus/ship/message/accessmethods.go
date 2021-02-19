package message

import "github.com/andig/evcc/hems/eebus/util"

type CmiAccessMethodsRequest struct {
	AccessMethodsRequest []AccessMethodsRequest `json:"accessMethodsRequest"`
}

type AccessMethodsRequest struct{}

type CmiAccessMethods struct {
	AccessMethods []AccessMethods `json:"accessMethods"`
}

type AccessMethods struct {
	ID string `json:"id"`
	// DnsSDmDns string `json:"dnsSd_mDns,omitempty"`
	// Dns       *struct {
	// 	URI string `json:"uri"`
	// } `json:"dns,omitempty"`
}

// MarshalJSON is the SHIP serialization marshaller
func (m AccessMethods) MarshalJSON() ([]byte, error) {
	return util.Marshal(m)
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *AccessMethods) UnmarshalJSON(data []byte) error {
	return util.Unmarshal(data, &m)
}
