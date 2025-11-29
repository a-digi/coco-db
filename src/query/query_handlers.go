package query

import (
	"net/http"
	"github.com/a-digi/coco-db/src/response"
	"github.com/a-digi/coco-db/src/table/fields"
	"github.com/a-digi/coco-db/src/logger"
)

// QueryHandler kapselt DataDir und Logger für Query-Endpunkte
// und ermöglicht eine objektorientierte Handler-Architektur
// Beispiel: handler := &QueryHandler{DataDir: "./data", Logger: &logger.NoopLogger{}}
type QueryHandler struct {
	DataDir string
	Logger  logger.Logger
}

// TableQueryHandler verarbeitet eine Tabellen-Query und gibt eine APIResponse zurück
func (h *QueryHandler) TableQueryHandler(dbName, tableName string, r *http.Request) *response.APIResponse {
	queryObj, err := ParseQuery(r, 100, 1000)
	if err != nil {
		return response.WriteErrorInternal(http.StatusBadRequest, "ERR_INVALID_QUERY", err.Error(), "")
	}
	meta, err := fields.LoadTableMeta(h.DataDir, dbName, tableName)
	if err != nil {
		return response.WriteErrorInternal(http.StatusNotFound, "ERR_META_NOT_FOUND", err.Error(), "")
	}
	entries, err := FilterEngine(h.DataDir, dbName, tableName, queryObj, meta)
	if err != nil {
		return response.WriteErrorInternal(http.StatusInternalServerError, "ERR_QUERY_EXEC", err.Error(), "")
	}

	return response.WriteSuccess(entries, "")
}

// QueryHandler verarbeitet eine globale Query und gibt eine APIResponse zurück
func (h *QueryHandler) QueryHandler(dbName string, r *http.Request) *response.APIResponse {
	queryObj, err := ParseQuery(r, 100, 1000)
	if err != nil {
		return response.WriteErrorInternal(http.StatusBadRequest, "ERR_INVALID_QUERY", err.Error(), "")
	}
	// Dummy: Nur users-Tabelle, TODO: alle Tabellen iterieren
	meta, err := fields.LoadTableMeta(h.DataDir, dbName, "users")
	if err != nil {
		return response.WriteErrorInternal(http.StatusNotFound, "ERR_META_NOT_FOUND", err.Error(), "")
	}
	entries, err := FilterEngine(h.DataDir, dbName, "users", queryObj, meta)
	if err != nil {
		return response.WriteErrorInternal(http.StatusInternalServerError, "ERR_QUERY_EXEC", err.Error(), "")
	}

	return response.WriteSuccess(entries, "")
}
