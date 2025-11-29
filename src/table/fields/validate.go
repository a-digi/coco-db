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

// ValidatorHandler ist eine Funktion, die prüft, ob sie für den Typ zuständig ist und ggf. validiert
// Gibt nil zurück, wenn sie nicht zuständig ist
// Gibt *types.ValidationError zurück, wenn sie zuständig ist (Fehler oder nil bei Erfolg)
type ValidatorHandler func(value interface{}, meta types.FieldMeta) *types.ValidationError

func stringHandler(value interface{}, meta types.FieldMeta) *types.ValidationError {
	if meta.Type == "string" {
		return validate.ValidateString(value, meta)
	}
	return nil
}

func integerHandler(value interface{}, meta types.FieldMeta) *types.ValidationError {
	if meta.Type == "integer" {
		return validate.ValidateInteger(value, meta)
	}
	return nil
}

func booleanHandler(value interface{}, meta types.FieldMeta) *types.ValidationError {
	if meta.Type == "boolean" {
		return validate.ValidateBoolean(value, meta)
	}
	return nil
}

func jsonHandler(value interface{}, meta types.FieldMeta) *types.ValidationError {
	if meta.Type == "json" {
		return validate.ValidateJSON(value, meta)
	}
	return nil
}

func dateHandler(value interface{}, meta types.FieldMeta) *types.ValidationError {
	if meta.Type == "date" {
		return validate.ValidateDate(value, meta)
	}
	return nil
}

var validatorChain = []ValidatorHandler{
	stringHandler,
	integerHandler,
	booleanHandler,
	jsonHandler,
	dateHandler,
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

	for _, handler := range validatorChain {
		err := handler(value, meta)
		// Handler ist zuständig, gibt Fehler oder nil zurück
		if meta.Type == "string" || meta.Type == "integer" || meta.Type == "boolean" || meta.Type == "json" || meta.Type == "date" {
			return err
		}
		if err != nil {
			return err
		}
	}
	return &types.ValidationError{
		Field:   meta.Name,
		Code:    "ERR_TYPE_UNKNOWN",
		Message: "Unbekannter Feldtyp",
	}
}
