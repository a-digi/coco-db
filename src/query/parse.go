package query

import (
	"encoding/json"
	"net/http"
	"fmt"
)

// ParseQuery liest und parst den JSON-Body in ein Query-Objekt und prüft Pflichtfelder/Defaults
func ParseQuery(r *http.Request, defaultLimit, maxLimit int) (*Query, error) {
	var qr Query
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&qr); err != nil {
		return nil, fmt.Errorf("Malformed JSON: %w", err)
	}
	// Pflichtfeld Filter prüfen
	if qr.Filter == nil {
		return nil, fmt.Errorf("Missing required field: filter")
	}
	// Limit prüfen und ggf. setzen
	if qr.Limit < 0 {
		return nil, fmt.Errorf("Limit must be >= 0")
	}
	if qr.Limit == 0 && defaultLimit > 0 {
		qr.Limit = defaultLimit
	}
	if maxLimit > 0 && qr.Limit > maxLimit {
		qr.Limit = maxLimit
	}
	// Offset prüfen
	if qr.Offset < 0 {
		return nil, fmt.Errorf("Offset must be >= 0")
	}
	return &qr, nil
}
