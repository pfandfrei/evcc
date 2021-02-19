package spine

import (
	"encoding/json"
	"fmt"

	"github.com/andig/evcc/hems/eebus/ship/message"
)

// {"datagram":[
// 	{"header":[
// 		{"specificationVersion":"1.2.0"},
// 		{"addressSource":[
// 			{"device":"d:_i:3210_ESystemsMtg-CEM"},{"entity":[0]},{"feature":0}
// 		]},
// 		{"addressDestination":[
// 			{"entity":[0]},{"feature":0}
// 		]},
// 		{"msgCounter":5876},
// 		{"cmdClassifier":"read"}
// 	]},
// 	{"payload":[
// 		{"cmd":[
// 			[{"nodeManagementDetailedDiscoveryData":[]}]
// 		]}
// 	]}
// ]}

func init() {
	d := message.CmiData{}
	b, err := json.Marshal(d)
	fmt.Println(string(b), err)
}
