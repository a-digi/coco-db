// src/table/fields/validate.go
// Zentrale Validierungsfunktion für Einträge gegen TableMeta/FieldMeta
// Siehe Checkliste Punkt 2 in project/felder_und_unterstuetzte_datentypen.md

package fields

import (
	"fmt"
	"github.com/a-digi/coco-db/src/table/fields/types"
	"github.com/a-digi/coco-db/src/table/fields/validate"
)

// ValidateEntry prüft ein Dokument gegen das TableMeta-Schema
// Gibt nil zurück, wenn alles gültig ist, sonst eine Liste von ValidationError
func ValidateEntry(entry map[string]interface{}, meta types.TableMeta) []types.ValidationError {
	errors := []types.ValidationError{}
	fieldsByName := map[string]types.FieldMeta{}
	for _, f := range meta.Fields {
		fieldsByName[f.Name] = f
	}

	// 1. Pflichtfelder prüfen
	for _, f := range meta.Fields {
		if f.Required {
			if _, ok := entry[f.Name]; !ok {
				errors = append(errors, types.ValidationError{
					Field:   f.Name,
					Code:    "ERR_FIELD_MISSING",
					Message: fmt.Sprintf("Feld '%s' ist erforderlich", f.Name),
				})
			}
		}
	}

	// 2. Felder prüfen
	for k, v := range entry {
		f, ok := fieldsByName[k]
		if !ok {
			if meta.AllowAdditionalFields != nil && *meta.AllowAdditionalFields {
				continue // erlaubt
			}
			errors = append(errors, types.ValidationError{
				Field:   k,
				Code:    "ERR_FIELD_NOT_ALLOWED",
				Message: fmt.Sprintf("Feld '%s' ist nicht im Schema definiert", k),
			})
			continue
		}
		// Typprüfung und Constraints
		if err := validateFieldValue(v, f); err != nil {
			errors = append(errors, *err)
		}
	}

	return errors
}

// Exportiert für rekursive Validatoren
func ValidateFieldValue(value interface{}, meta types.FieldMeta) *types.ValidationError {
	return validateFieldValue(value, meta)
}

// validateFieldValue prüft Typ und Constraints für ein Feld
// Gibt nil zurück, wenn alles gültig ist, sonst ValidationError
func validateFieldValue(value interface{}, meta types.FieldMeta) *types.ValidationError {
	// 1. Nullable prüfen
	if value == nil {
		if meta.Nullable != nil && *meta.Nullable {
			return nil
		}
		return &types.ValidationError{
			Field:   meta.Name,
			Code:    "ERR_NULL_NOT_ALLOWED",
			Message: "Feld darf nicht null sein",
		}
	}

	switch meta.Type {
	case "string":
		return validate.ValidateString(value, meta)
	case "integer":
		return validate.ValidateInteger(value, meta)
	case "boolean":
		return validate.ValidateBoolean(value, meta)
	case "json":
		return validate.ValidateJSON(value, meta)
	case "date":
		return validate.ValidateDate(value, meta)
	default:
		return &types.ValidationError{
			Field:   meta.Name,
			Code:    "ERR_TYPE_UNKNOWN",
			Message: "Unbekannter Feldtyp",
		}
	}
}
