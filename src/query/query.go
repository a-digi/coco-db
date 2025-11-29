package query

import (
	"encoding/json"
	"net/http"
	"fmt"
)

// QueryRequest repräsentiert die Query-API-Struktur
// (Filter, Limit, Offset, Sortierung, Joins)
type QueryRequest struct {
	Filter map[string]interface{} `json:"filter"`
	Limit  int                    `json:"limit"`
	Offset int                    `json:"offset"`
	Sort   []string               `json:"sort"`
	Join   []map[string]interface{} `json:"join"` // optional, für globale Query
}

// ParseQueryRequest liest und parst den JSON-Body in ein QueryRequest-Struct
func ParseQueryRequest(r *http.Request) (*QueryRequest, error) {
	var qr QueryRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&qr); err != nil {
		return nil, fmt.Errorf("Malformed JSON: %w", err)
	}
	return &qr, nil
}

// (Die Funktion ValidateQueryRequest wurde in src/query/validate.go ausgelagert)
