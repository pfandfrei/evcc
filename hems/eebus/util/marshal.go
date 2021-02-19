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
	var e []map[string]interface{}

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
	var e []map[string]json.RawMessage

	if data[0] == byte('[') {
		// fmt.Printf("unmarshal: slice %s\n", string(data))
		if err := json.Unmarshal(data, &e); err != nil {
			return err
		}
	} else {
		// fmt.Printf("unmarshal: map %s\n", string(data))
		e = append(e, map[string]json.RawMessage{})
		if err := json.Unmarshal(data, &e[0]); err != nil {
			return err
		}
	}

	for _, m := range e {
		if len(m) > 1 {
			return fmt.Errorf("unmarshal: invalid map %v", m)
		}

		// extract map
		var key string
		var val json.RawMessage
		for k, v := range m {
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

		iface := field.Value()
		switch field.Kind() {
		case reflect.Struct, reflect.Slice:
			iface = reflect.New(reflect.TypeOf(field.Value())).Interface() // struct to interface
		}

		// fmt.Printf("iface: %T %s\n", iface, reflect.TypeOf(field.Value()).Kind())
		if err := json.Unmarshal(val, &iface); err != nil {
			return err
		}

		switch field.Kind() {
		case reflect.Int:
			iface = int(iface.(float64))

		case reflect.Struct, reflect.Slice:
			elem := reflect.ValueOf(iface).Elem() // de-reference
			switch elem.Kind() {
			case reflect.Struct, reflect.Slice:
				iface = elem.Interface()
			case reflect.String:
				iface = elem.String()
			}
		}

		if err := field.Set(iface); err != nil {
			// fmt.Printf("set: %s=%+v %v\n", field.Name(), iface, err)
			return err
		}
	}

	return nil
}
