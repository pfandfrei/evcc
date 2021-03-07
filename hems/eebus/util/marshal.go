package util

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/fatih/structs"
)

// Marshal is the SHIP serialization
func Marshal(v interface{}) ([]byte, error) {
	e := make([]map[string]interface{}, 0)

	for _, f := range structs.Fields(v) {
		if !f.IsExported() {
			continue
		}

		jsonTag := f.Tag("json")
		if f.IsZero() && strings.HasSuffix(jsonTag, ",omitempty") {
			continue
		}

		key := f.Name()
		if jsonTag != "" {
			key = strings.TrimSuffix(jsonTag, ",omitempty")
		}

		m := map[string]interface{}{key: f.Value()}
		e = append(e, m)
	}

	return json.Marshal(e)
}

// Unmarshal is the SHIP de-serialization
func Unmarshal(data []byte, v interface{}) error {
	var ar []map[string]json.RawMessage

	// convert input to json array
	if data[0] != byte('[') {
		data = append([]byte{'['}, append(data, ']')...)
	}
	if err := json.Unmarshal(data, &ar); err != nil {
		return err
	}

	// convert array elements to struct members
	for _, ae := range ar {
		if len(ae) > 1 {
			return fmt.Errorf("unmarshal: invalid map %v", ae)
		}

		// extract 1-element map
		var key string
		var val json.RawMessage
		for k, v := range ae {
			key = k
			val = v
		}

		// find field
		var field *structs.Field
		for _, f := range structs.Fields(v) {
			name := f.Name()
			if jsonTag := f.Tag("json"); jsonTag != "" {
				name = strings.TrimSuffix(jsonTag, ",omitempty")
			}

			if name == key {
				field = f
				break
			}
		}

		if field == nil {
			return fmt.Errorf("unmarshal: field not found: %s", key)
		}

		// iface := field.Value()
		iface := reflect.New(
			reflect.TypeOf(field.Value()),
		).Interface()
		// iface := field.Value()
		fmt.Printf("iface 1: %T %s\n", iface, field.Kind())

		// switch field.Kind() {
		// case reflect.Struct, reflect.Slice:
		// 	iface = reflect.New(reflect.TypeOf(field.Value())).Interface() // struct to interface
		// }

		fmt.Println(string(val))
		fmt.Printf("iface 2: %T %s %s\n", iface, reflect.TypeOf(field.Value()).String(), reflect.TypeOf(field.Value()).Kind())
		err := json.Unmarshal(val, iface)
		if err != nil {
			// panic(err)
			return err
		}

		// switch field.Kind() {
		// case reflect.Int:
		// 	iface = int(iface.(float64))

		// case reflect.Struct, reflect.Slice:
		// 	elem := reflect.ValueOf(iface).Elem() // de-reference
		// 	switch elem.Kind() {
		// 	case reflect.Struct, reflect.Slice:
		// 		iface = elem.Interface()
		// 	case reflect.String:
		// 		iface = elem.String()
		// 	}
		// }

		fmt.Printf("set: %s=%+v (%T)\n", field.Name(), iface, iface)
		if err := field.Set(iface); err != nil {
			// fmt.Printf("set: %s=%+v %v\n", field.Name(), iface, err)
			return err
		}
	}

	return nil
}
