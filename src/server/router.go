// Routing und HTTP-Handler für die Kernoperationen (CRUD, Index, Events, Transaktionen)
// Hier werden die Endpunkte gemäß API-Design angebunden

package server

import (
	"net/http"
	"github.com/a-digi/coco-db/src/response"
	"github.com/a-digi/coco-db/src/database"
)

// apiHandler ist ein Wrapper, der Handler mit APIResponse-Signatur in http.HandlerFunc umwandelt
func apiHandler(fn func(http.ResponseWriter, *http.Request) *response.APIResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resp := fn(w, r)
		if resp == nil {
			response.WriteError(w, http.StatusInternalServerError, "internal_error", "Leere APIResponse", "")
		}
	}
}

// SetupRouter initialisiert das Routing für alle Kernoperationen
func SetupRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Beispiel-Handler für Health-Check
	mux.HandleFunc("/health", apiHandler(HealthHandler))

	// Datenbank-Endpunkte
	mux.HandleFunc("/api/databases", apiHandler(database.Handler))        // POST, GET
	mux.HandleFunc("/api/databases/", apiHandler(database.Handler)) // DELETE (mit Name im Pfad)

	// TODO: Weitere Endpunkte für CRUD, Index, Events, Transaktionen
	// z.B. mux.HandleFunc("/api/databases/{dbname}/tables", apiHandler(TableHandler))
	//      mux.HandleFunc("/api/databases/{dbname}/tables/{tablename}/documents", apiHandler(DocumentHandler))
	//      mux.HandleFunc("/api/databases/{dbname}/tables/{tablename}/indexes", apiHandler(IndexHandler))
	//      mux.HandleFunc("/api/events", apiHandler(EventHandler))
	//      mux.HandleFunc("/api/transactions", apiHandler(TransactionHandler))

	return mux
}
