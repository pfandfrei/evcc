package main

import (
	"encoding/json"
	"fmt"
)

// Version of the protocol
type Version struct {
	Major int `json:"major"`
}

// UnmarshalJSON is the SHIP serialization unmarshaller
func (m *Version) UnmarshalJSON(data []byte) error {
	println("invoked")
	return nil
}

func unmarshal(val interface{}) {
	v := `{"major":0}`
	fmt.Println(json.Unmarshal([]byte(v), &val), val)
}

func main() {
	ver0 := Version{}
	unmarshal(&ver0)

	ver1 := &Version{}
	unmarshal(&ver1)

	var ver2 interface{}
	ver2 = Version{}
	unmarshal(&ver2)

	var ver3 interface{}
	ver3 = &Version{}
	unmarshal(&ver3)
}
