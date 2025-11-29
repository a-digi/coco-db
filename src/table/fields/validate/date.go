package validate

import (
	"time"
	"github.com/a-digi/coco-db/src/table/fields/types"
)

// Nur ISO 8601-Prüfung für date
func ValidateDate(value interface{}, meta types.FieldMeta) *types.ValidationError {
	if s, ok := value.(string); ok {
		_, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return &types.ValidationError{
				Field:   meta.Name,
				Code:    "ERR_TYPE_MISMATCH",
				Message: "Feldtyp muss ISO 8601 Datum/Zeit sein",
			}
		}
	} else {
		return &types.ValidationError{
			Field:   meta.Name,
			Code:    "ERR_TYPE_MISMATCH",
			Message: "Feldtyp muss string (ISO 8601) sein",
		}
	}
	return nil
}
