// src/table/fields/types/types.go
// Enthält zentrale Typen für Felder, Tabellen und Validierung
package types

type FieldMeta struct {
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Required    bool          `json:"required,omitempty"`
	Default     interface{}   `json:"default,omitempty"`
	MinLength   *int          `json:"minLength,omitempty"`
	MaxLength   *int          `json:"maxLength,omitempty"`
	Min         *float64      `json:"min,omitempty"`
	Max         *float64      `json:"max,omitempty"`
	Nullable    *bool         `json:"nullable,omitempty"`
	Pattern     string        `json:"pattern,omitempty"`
	Enum        []interface{} `json:"enum,omitempty"`
	Description string        `json:"description,omitempty"`
	Unique      *bool         `json:"unique,omitempty"`
	Fields      []FieldMeta   `json:"fields,omitempty"`
}

type IndexMeta struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"`
	Fields []string `json:"fields"`
	Unique bool     `json:"unique,omitempty"`
	Sparse bool     `json:"sparse,omitempty"`
}

type TableMeta struct {
	TableName             string                 `json:"tableName"`
	SchemaVersion         int                    `json:"schemaVersion"`
	Fields                []FieldMeta            `json:"fields"`
	Indexes               []IndexMeta            `json:"indexes,omitempty"`
	Options               map[string]interface{} `json:"options,omitempty"`
	AllowAdditionalFields *bool                  `json:"allowAdditionalFields,omitempty"`
}

type ValidationError struct {
	Field   string
	Code    string
	Message string
}

