package registry

import (
	"fmt"
	"strings"
)

type Registry[T any] map[string]func(map[string]interface{}) (T, error)

// var registry [T]Registry = make(map[string]func(map[string]interface{}) (T, error))

func New[T any]() Registry[T] {
	reg := make(map[string]func(map[string]interface{}) (T, error))
	return reg
}

func (r Registry[T]) Add(name string, factory func(map[string]interface{}) (T, error)) {
	if _, exists := r[name]; exists {
		panic(fmt.Sprintf("cannot register duplicate charger type: %s", name))
	}
	r[name] = factory
}

func (r Registry[T]) Get(name string) (func(map[string]interface{}) (T, error), error) {
	factory, exists := r[name]
	if !exists {
		return nil, fmt.Errorf("charger type not registered: %s", name)
	}
	return factory, nil
}


// NewFromConfig creates charger from configuration
func (r Registry[T]) NewFromConfig(typ string, other map[string]interface{}) (v T, err error) {
	factory, err := r.Get(strings.ToLower(typ))
	if err == nil {
		if v, err = factory(other); err != nil {
			err = fmt.Errorf("cannot create type '%s': %w", typ, err)
		}
	} else {
		err = fmt.Errorf("invalid charger type: %s", typ)
	}

	return
}
