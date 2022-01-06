package types

import "fmt"

// Array represents BSON array.
type Array []any

// Get returns a value at the given index.
func (a Array) Get(index int) (any, error) {
	if l := len(a); index < 0 || index >= l {
		return nil, fmt.Errorf("types.Array.Get: index %d is out of bounds [0-%d)", index, l)
	}

	return a[index], nil
}

// Set sets the value at the given index.
func (a Array) Set(index int, value any) error {
	if l := len(a); index < 0 || index >= l {
		return fmt.Errorf("types.Array.Set: index %d is out of bounds [0-%d)", index, l)
	}

	if err := validateValue(value); err != nil {
		return fmt.Errorf("types.Array.Set: %w", err)
	}

	a[index] = value
	return nil
}
