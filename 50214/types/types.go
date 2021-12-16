package types

import (
	"fmt"
)

// validateValue validates value.
func validateValue(value any) error {
	switch value := value.(type) {
	case float64:
		return nil
	case string:
		return nil
	case Document:
		return value.validate()
	case Array:
		return nil
	case bool:
		return nil
	case nil:
		return nil
	case int32:
		return nil
	case int64:
		return nil
	default:
		return fmt.Errorf("types.validateValue: unsupported type: %T", value)
	}
}
