// Verwaltung von Datenbanken: Anlegen, Löschen, Auflisten
// Entspricht Punkt 3 der Roadmap und den Anforderungen aus PROJECT_REQUIREMENT.md

package database

import (
	"regexp"
)

var DbNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,32}$`)

// DatabaseRequest repräsentiert das Request-Objekt für das Anlegen einer Datenbank
type DatabaseRequest struct {
	Name string `json:"name"`
}

// Hilfsfunktion für Prefix-Matching
func startsWith(s, prefix string) bool {
	return len(s) >= len(prefix) && s[:len(prefix)] == prefix
}
