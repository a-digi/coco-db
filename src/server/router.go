// Routing und HTTP-Handler für die Kernoperationen (CRUD, Index, Events, Transaktionen)
// Hier werden die Endpunkte gemäß API-Design angebunden

package server

import (
	"net/http"
)

// SetupRouter initialisiert das Routing für alle Kernoperationen
func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Beispiel-Handler für Health-Check
	mux.HandleFunc("/health", HealthHandler)

	// TODO: Weitere Endpunkte für CRUD, Index, Events, Transaktionen
	// z.B. mux.HandleFunc("/api/databases", DatabaseHandler)
	//      mux.HandleFunc("/api/databases/{dbname}/tables", TableHandler)
	//      mux.HandleFunc("/api/databases/{dbname}/tables/{tablename}/documents", DocumentHandler)
	//      mux.HandleFunc("/api/databases/{dbname}/tables/{tablename}/indexes", IndexHandler)
	//      mux.HandleFunc("/api/events", EventHandler)
	//      mux.HandleFunc("/api/transactions", TransactionHandler)

	return mux
}

