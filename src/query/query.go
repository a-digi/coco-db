package query

// JoinDef beschreibt einen Join in der globalen Query
// table: Name der Zieltabelle
// on: Join-Bedingung (Mapping Quellfeld → Zielfeld)
// filter: optionaler Filter für die Join-Tabelle
// fields: optionale Felderliste für die Join-Tabelle
// join: optionale weitere verschachtelte Joins
//
type JoinDef struct {
	Table  string                 `json:"table"`
	On     map[string]string      `json:"on"`
	Filter map[string]interface{} `json:"filter,omitempty"`
	Fields []string               `json:"fields,omitempty"`
	Join   []JoinDef              `json:"join,omitempty"`
}

// Query repräsentiert die Query-API-Struktur
// (Filter, Limit, Offset, Sortierung, Joins)
type Query struct {
	Filter map[string]interface{} `json:"filter"`
	Limit  int                    `json:"limit"`
	Offset int                    `json:"offset"`
	Sort   []string               `json:"sort"`
	Join   []JoinDef              `json:"join"`
	Fields []string               `json:"fields"`
}
