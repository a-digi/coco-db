// src/table/fields/types.go
// Definition der zentralen Datenstrukturen für Felder, Tabellen und Indizes
// Siehe Dokumentation in project/felder_und_unterstuetzte_datentypen.md

package fields

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// FieldMeta beschreibt ein Feld in meta.json gemäß Projektanforderungen.
type FieldMeta struct {
	// Name des Feldes (Pflichtfeld)
	Name string `json:"name"`
	// Datentyp: string, int, float, bool, object, array, date (ISO 8601)
	Type string `json:"type"`
	// Pflichtfeld?
	Required bool `json:"required,omitempty"`
	// Standardwert, falls nicht gesetzt
	Default interface{} `json:"default,omitempty"`
	// Minimale/maximale Länge (nur für string/array)
	MinLength *int `json:"minLength,omitempty"`
	MaxLength *int `json:"maxLength,omitempty"`
	// Minimaler/maximaler Wert (nur für int/float)
	Min *float64 `json:"min,omitempty"`
	Max *float64 `json:"max,omitempty"`
	// Feld darf null sein?
	Nullable *bool `json:"nullable,omitempty"`
	// Regex-Pattern (nur für string)
	Pattern string `json:"pattern,omitempty"`
	// Werteliste (enum)
	Enum []interface{} `json:"enum,omitempty"`
	// Beschreibung (optional)
	Description string `json:"description,omitempty"`
	// Eindeutigkeit (optional, für spätere Index-Validierung)
	Unique *bool `json:"unique,omitempty"`
	// Für Typ object/array: optionale Definition von Subfeldern
	Fields []FieldMeta `json:"fields,omitempty"`
}

// IndexMeta beschreibt einen Index in meta.json
// Typ: "primary" oder "secondary"
type IndexMeta struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Fields []string `json:"fields"`
	Unique bool     `json:"unique,omitempty"`
	Sparse bool     `json:"sparse,omitempty"`
}

// TableMeta entspricht exakt dem meta.json-Schema laut Projektanforderungen
type TableMeta struct {
	TableName             string                 `json:"tableName"`
	SchemaVersion         int                    `json:"schemaVersion"`
	Fields                []FieldMeta            `json:"fields"`
	Indexes               []IndexMeta            `json:"indexes,omitempty"`
	Options               map[string]interface{} `json:"options,omitempty"`
	AllowAdditionalFields *bool                  `json:"allowAdditionalFields,omitempty"`
}

// LoadTableMeta lädt meta.json als types.TableMeta
func LoadTableMeta(dataDir, dbName, tableName string) (*TableMeta, error) {
	metaPath := filepath.Join(dataDir, dbName, tableName, "meta.json")
	f, err := os.Open(metaPath)
	if err != nil {
		return nil, fmt.Errorf("meta.json nicht gefunden: %w", err)
	}
	defer f.Close()
	var meta TableMeta
	if err := json.NewDecoder(f).Decode(&meta); err != nil {
		return nil, fmt.Errorf("meta.json ungültig: %w", err)
	}
	return &meta, nil
}
