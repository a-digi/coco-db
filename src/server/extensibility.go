// Vorbereitungen für Erweiterbarkeit: Schnittstellen und Strukturen für zukünftige Features
// Hier werden Platzhalter und Abstraktionen für Auth, Event Sourcing, Indexe etc. definiert

package server

// AuthProvider definiert die Schnittstelle für Authentifizierungsmechanismen
// Kann später durch verschiedene Implementierungen (z.B. JWT, Basic Auth) ersetzt werden
type AuthProvider interface {
	Authenticate(token string) (userID string, err error)
}

// EventStore definiert die Schnittstelle für Event Sourcing
// Ermöglicht verschiedene Backends (Datei, DB, etc.)
type EventStore interface {
	AppendEvent(event interface{}) error
	GetEvents(since int64) ([]interface{}, error)
}

// IndexEngine definiert die Schnittstelle für Indexverwaltung
// Ermöglicht verschiedene Index-Backends
// (z.B. In-Memory, Datei, später DB)
type IndexEngine interface {
	CreateIndex(field string, unique bool) error
	DropIndex(field string) error
	ListIndexes() ([]string, error)
}

// Erweiterbarkeit: Weitere Schnittstellen und Platzhalter können hier ergänzt werden

