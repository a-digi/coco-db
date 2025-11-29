package validate

import (
	"fmt"
	"regexp"
	"github.com/a-digi/coco-db/src/table/fields/types"
)

// ValidateString prüft einen Wert vom Typ string und die zugehörigen Constraints
func ValidateString(value interface{}, meta types.FieldMeta) *types.ValidationError {
	str, ok := value.(string)
	if !ok {
		return &types.ValidationError{
			Field:   meta.Name,
			Code:    "ERR_TYPE_MISMATCH",
			Message: "Feldtyp muss string sein",
		}
	}
	if meta.MinLength != nil && len(str) < *meta.MinLength {
		return &types.ValidationError{
			Field:   meta.Name,
			Code:    "ERR_MIN_LENGTH",
			Message: fmt.Sprintf("Feld muss mindestens %d Zeichen lang sein", *meta.MinLength),
		}
	}
	if meta.MaxLength != nil && len(str) > *meta.MaxLength {
		return &types.ValidationError{
			Field:   meta.Name,
			Code:    "ERR_MAX_LENGTH",
			Message: fmt.Sprintf("Feld darf maximal %d Zeichen lang sein", *meta.MaxLength),
		}
	}
	if meta.Pattern != "" {
		re, err := regexp.Compile(meta.Pattern)
		if err == nil && !re.MatchString(str) {
			return &types.ValidationError{
				Field:   meta.Name,
				Code:    "ERR_PATTERN",
				Message: "Feld entspricht nicht dem geforderten Muster",
			}
		}
	}
	if len(meta.Enum) > 0 {
		found := false
		for _, ev := range meta.Enum {
			if s, ok := ev.(string); ok && s == str {
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
