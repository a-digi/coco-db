package query

// Query repräsentiert die Query-API-Struktur
// (Filter, Limit, Offset, Sortierung, Joins)
type Query struct {
	Filter map[string]interface{} `json:"filter"`
	Limit  int                    `json:"limit"`
	Offset int                    `json:"offset"`
	Sort   []string               `json:"sort"`
	Join   []map[string]interface{} `json:"join"` // optional, für globale Query
}
