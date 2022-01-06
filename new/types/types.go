package types

import (
	"fmt"
)

// validateValue validates value.
func validateValue(value any) error {
	switch value := value.(type) {
	case Document:
		return value.validate()
	case Array:
		return nil
	default:
		return fmt.Errorf("types.validateValue: unsupported type: %T", value)
	}
}
