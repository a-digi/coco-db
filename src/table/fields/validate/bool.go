package validate

import "github.com/a-digi/coco-db/src/table/fields/types"

// ValidateBoolean prüft einen Wert vom Typ boolean und die zugehörigen Constraints
func ValidateBoolean(value interface{}, meta types.FieldMeta) *types.ValidationError {
	b, ok := value.(bool)
	if !ok {
		return &types.ValidationError{
			Field:   meta.Name,
			Code:    "ERR_TYPE_MISMATCH",
			Message: "Feldtyp muss bool sein",
		}
	}
	if len(meta.Enum) > 0 {
		found := false
		for _, ev := range meta.Enum {
			if eb, ok := ev.(bool); ok && eb == b {
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
