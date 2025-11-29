package query

import (
	"fmt"
	"github.com/a-digi/coco-db/src/table/fields"
)

type ValidationError struct {
	Field   string
	Code    string
	Message string
}

// ValidateQuery prüft die Query-Struktur und die Filter gegen das TableMeta-Schema
func ValidateQuery(qr *Query, meta *fields.TableMeta) []ValidationError {
	errors := []ValidationError{}
	if qr == nil {
		errors = append(errors, ValidationError{"_query", "ERR_QUERY_NIL", "Query is nil"})
		return errors
	}
	if qr.Filter == nil {
		errors = append(errors, ValidationError{"filter", "ERR_FILTER_MISSING", "Missing filter object"})
	}
	if qr.Limit < 0 {
		errors = append(errors, ValidationError{"limit", "ERR_LIMIT_INVALID", "Limit must be >= 0"})
	}
	if qr.Offset < 0 {
		errors = append(errors, ValidationError{"offset", "ERR_OFFSET_INVALID", "Offset must be >= 0"})
	}
	if meta == nil {
		errors = append(errors, ValidationError{"_meta", "ERR_META_NIL", "TableMeta is nil"})
		return errors
	}
	fieldsByName := map[string]fields.FieldMeta{}
	for _, f := range meta.Fields {
		fieldsByName[f.Name] = f
	}
	allowAdditional := false
	if meta.AllowAdditionalFields != nil && *meta.AllowAdditionalFields {
		allowAdditional = true
	}
	// Filter-Felder prüfen
	for k, v := range qr.Filter {
		f, ok := fieldsByName[k]
		if !ok {
			if !allowAdditional {
				errors = append(errors, ValidationError{k, "ERR_FIELD_NOT_ALLOWED", "Feld ist nicht im Schema definiert"})
			}
			continue
		}
		// Typprüfung und Operatorprüfung
		switch val := v.(type) {
		case map[string]interface{}:
			for op, opVal := range val {
				if !isAllowedOperator(op, f.Type) {
					errors = append(errors, ValidationError{k, "ERR_OPERATOR_NOT_ALLOWED", fmt.Sprintf("Operator '%s' ist für Typ '%s' nicht erlaubt", op, f.Type)})
				}
				if !isValidValueForType(opVal, f.Type) {
					errors = append(errors, ValidationError{k, "ERR_TYPE_MISMATCH", fmt.Sprintf("Wert für Operator '%s' passt nicht zu Typ '%s'", op, f.Type)})
				}
				// Optional: Werteformat prüfen (z.B. Datum, Enum, Pattern)
			}
		default:
			// Direkter Vergleichswert (implizit eq)
			if !isValidValueForType(val, f.Type) {
				errors = append(errors, ValidationError{k, "ERR_TYPE_MISMATCH", fmt.Sprintf("Wert passt nicht zu Typ '%s'", f.Type)})
			}
		}
	}
	return errors
}

// isAllowedOperator prüft, ob der Operator für den Feldtyp zulässig ist
func isAllowedOperator(op string, typ string) bool {
	// Basis-Operatoren pro Typ
	allowed := map[string][]string{
		"string":  {"eq", "neq", "like", "fulltext", "partial"},
		"int":     {"eq", "neq", "gt", "gte", "lt", "lte"},
		"integer": {"eq", "neq", "gt", "gte", "lt", "lte"},
		"float":   {"eq", "neq", "gt", "gte", "lt", "lte"},
		"bool":    {"eq", "neq"},
		"boolean": {"eq", "neq"},
		"date":    {"eq", "neq", "gt", "gte", "lt", "lte", "between"},
		"json":    {"like", "fulltext", "partial"},
	}
	for _, a := range allowed[typ] {
		if op == a {
			return true
		}
	}
	return false
}

// isValidValueForType prüft, ob der Wert zum Feldtyp passt
func isValidValueForType(val interface{}, typ string) bool {
	switch typ {
	case "string":
		_, ok := val.(string)
		return ok
	case "int", "integer":
		// JSON-Zahlen werden als float64 geparst, aber auch int akzeptieren
		switch val.(type) {
		case float64, int, int64, int32, float32:
			return true
		default:
			return false
		}
	case "float":
		switch val.(type) {
		case float64, float32, int, int64, int32:
			return true
		default:
			return false
		}
	case "bool", "boolean":
		_, ok := val.(bool)
		return ok
	case "date":
		_, ok := val.(string)
		return ok // Optional: Formatprüfung
	case "json":
		_, ok := val.(string)
		return ok // Optional: JSON-String prüfen
	default:
		return false
	}
}
