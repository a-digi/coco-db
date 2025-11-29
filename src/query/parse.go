package query

import (
	"encoding/json"
	"net/http"
	"fmt"
)

// ParseQuery liest und parst den JSON-Body in ein Query-Objekt
func ParseQuery(r *http.Request) (*Query, error) {
	var qr Query
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&qr); err != nil {
		return nil, fmt.Errorf("Malformed JSON: %w", err)
	}
	return &qr, nil
}
