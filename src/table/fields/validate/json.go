package validate

import (
	"github.com/a-digi/coco-db/src/table/fields/types"
	"encoding/json"
	"fmt"
)

// ValidateJSON prüft einen Wert vom Typ json und die zugehörigen Constraints
func ValidateJSON(value interface{}, meta types.FieldMeta) *types.ValidationError {
	var jsonObj map[string]interface{}
	switch v := value.(type) {
	case string:
		if err := json.Unmarshal([]byte(v), &jsonObj); err != nil {
			return &types.ValidationError{
				Field:   meta.Name,
				Code:    "ERR_TYPE_MISMATCH",
				Message: "Feldtyp muss gültiges JSON sein",
			}
		}
	case map[string]interface{}:
		jsonObj = v
	default:
		return &types.ValidationError{
			Field:   meta.Name,
			Code:    "ERR_TYPE_MISMATCH",
			Message: "Feldtyp muss JSON sein",
		}
	}
	if meta.MinLength != nil && len(jsonObj) < *meta.MinLength {
		return &types.ValidationError{
			Field:   meta.Name,
			Code:    "ERR_MIN_LENGTH",
			Message: fmt.Sprintf("JSON-Objekt muss mindestens %d Felder haben", *meta.MinLength),
		}
	}
	if meta.MaxLength != nil && len(jsonObj) > *meta.MaxLength {
		return &types.ValidationError{
			Field:   meta.Name,
			Code:    "ERR_MAX_LENGTH",
			Message: fmt.Sprintf("JSON-Objekt darf maximal %d Felder haben", *meta.MaxLength),
		}
	}
	return nil
}

