package validate

import (
	"fmt"
	"github.com/a-digi/coco-db/src/table/fields/types"
)

// ValidateInteger prüft einen Wert vom Typ integer und die zugehörigen Constraints
func ValidateInteger(value interface{}, meta types.FieldMeta) *types.ValidationError {
	var ival int64
	switch v := value.(type) {
	case int:
		ival = int64(v)
	case int32:
		ival = int64(v)
	case int64:
		ival = v
	case float64:
		ival = int64(v)
	default:
		return &types.ValidationError{
			Field:   meta.Name,
			Code:    "ERR_TYPE_MISMATCH",
			Message: "Feldtyp muss integer sein",
		}
	}
	if meta.MinLength != nil && ival < int64(*meta.MinLength) {
		return &types.ValidationError{
			Field:   meta.Name,
			Code:    "ERR_MIN_LENGTH",
			Message: fmt.Sprintf("Feldwert muss mindestens %d sein", *meta.MinLength),
		}
	}
	if meta.MaxLength != nil && ival > int64(*meta.MaxLength) {
		return &types.ValidationError{
			Field:   meta.Name,
			Code:    "ERR_MAX_LENGTH",
			Message: fmt.Sprintf("Feldwert darf maximal %d sein", *meta.MaxLength),
		}
	}
	if len(meta.Enum) > 0 {
		found := false
		for _, ev := range meta.Enum {
			if evi, ok := ev.(float64); ok && int64(evi) == ival {
				found = true
				break
			}
		}
		if !found {
			return &types.ValidationError{
				Field:   meta.Name,
				Code:    "ERR_ENUM",
				Message: "Feldwert ist nicht in der erlaubten Werteliste",
			}
		}
	}
	return nil
}
