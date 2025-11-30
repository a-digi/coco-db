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
	// Sort prüfen (falls vorhanden)
	if qr.Sort != nil {
		for i, s := range qr.Sort {
			if s == "" {
				return nil, fmt.Errorf("Sort[%d] must be a non-empty string", i)
			}
		}
	}
	// Join prüfen (falls vorhanden)
	if qr.Join != nil {
		for i, j := range qr.Join {
			if j.Table == "" {
				return nil, fmt.Errorf("Join[%d] missing required field: table", i)
			}
			if len(j.On) == 0 {
				return nil, fmt.Errorf("Join[%d] missing required field: on", i)
			}
			// Rekursive Validierung für verschachtelte Joins
			if j.Join != nil {
				for k, sub := range j.Join {
					if sub.Table == "" {
						return nil, fmt.Errorf("Join[%d].Join[%d] missing required field: table", i, k)
					}
					if len(sub.On) == 0 {
						return nil, fmt.Errorf("Join[%d].Join[%d] missing required field: on", i, k)
					}
				}
			}
		}
	}
	return &qr, nil
}
