package query

import (
	"net/http"
	"time"
	"fmt"
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
	start := time.Now()
	queryObj, err := ParseQuery(r, 100, 1000)
	if err != nil {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusBadRequest, "ERR_INVALID_QUERY", err.Error(), execTime)
	}
	meta, err := fields.LoadTableMeta(h.DataDir, dbName, tableName)
	if err != nil {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusNotFound, "ERR_META_NOT_FOUND", err.Error(), execTime)
	}
	entries, err := FilterEngine(h.DataDir, dbName, tableName, queryObj, meta)
	if err != nil {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusInternalServerError, "ERR_QUERY_EXEC", err.Error(), execTime)
	}
	execTime := time.Since(start).String()
	return response.WriteSuccess(entries, execTime)
}

// joinDepth gibt die aktuelle Verschachtelungstiefe an (für Begrenzung)
func (h *QueryHandler) queryWithJoins(dbName, tableName string, query *Query, meta *fields.TableMeta, joinDefs []JoinDef, joinDepth, maxJoinDepth int, joinPath map[string]struct{}) ([]map[string]interface{}, error) {
	if joinDepth > maxJoinDepth {
		return nil, fmt.Errorf("Maximale Join-Tiefe (%d) überschritten", maxJoinDepth)
	}
	if joinPath == nil {
		joinPath = make(map[string]struct{})
	}
	pathKey := dbName + "." + tableName
	if _, exists := joinPath[pathKey]; exists {
		return nil, fmt.Errorf("Zyklischer Join erkannt: %s", pathKey)
	}
	joinPath[pathKey] = struct{}{}
	entries, err := FilterEngine(h.DataDir, dbName, tableName, query, meta)
	if err != nil {
		return nil, err
	}
	for i := range entries {
		for _, join := range joinDefs {
			joinMeta, err := fields.LoadTableMeta(h.DataDir, dbName, join.Table)
			if err != nil {
				continue
			}
			joinQuery := &Query{
				Filter: join.Filter,
				Limit:  0,
				Offset: 0,
				Sort:   nil,
				Join:   join.Join,
			}
			// Join-Bedingung: Mapping Ziel-Feld (in Join-Tabelle) → Quell-Feld (im Haupteintrag)
			joinFilter := map[string]interface{}{}
			for dst, src := range join.On {
				val, ok := entries[i][src]
				if !ok {
					if h.Logger != nil {
						h.Logger.Warning(fmt.Sprintf("Join-Mapping: Quellfeld '%s' fehlt im Eintrag: %+v", src, entries[i]))
					}
					continue
				}
				joinFilter[dst] = val
			}
			// Filter kombinieren
			for k, v := range joinQuery.Filter {
				joinFilter[k] = v
			}
			joinQuery.Filter = joinFilter
			// Rekursiver Join mit aktualisiertem joinPath
			joinResults, err := h.queryWithJoins(dbName, join.Table, joinQuery, joinMeta, join.Join, joinDepth+1, maxJoinDepth, copyJoinPath(joinPath))
			if err != nil {
				joinResults = []map[string]interface{}{} // Fehler: leeres Array statt null
			}
			if joinResults == nil {
				joinResults = []map[string]interface{}{}
			}
			if join.Fields != nil && len(join.Fields) > 0 {
				// Nur gewünschte Felder übernehmen
				for j := range joinResults {
					for k := range joinResults[j] {
						found := false
						for _, f := range join.Fields {
							if k == f {
								found = true
								break
							}
						}
						if !found {
							delete(joinResults[j], k)
						}
					}
				}
			}
			// Debug: Logge Join-Filter und Ergebnis-Anzahl
			if h.Logger != nil {
				h.Logger.Info(fmt.Sprintf("Join: %s, Filter: %+v, Treffer: %d", join.Table, joinQuery.Filter, len(joinResults)))
			}
			entries[i][join.Table] = joinResults
		}
	}
	return entries, nil
}

func copyJoinPath(orig map[string]struct{}) map[string]struct{} {
	newMap := make(map[string]struct{}, len(orig))
	for k, v := range orig {
		newMap[k] = v
	}
	return newMap
}

// QueryHandler verarbeitet eine globale Query und gibt eine APIResponse zurück
func (h *QueryHandler) QueryHandler(dbName string, r *http.Request) *response.APIResponse {
	start := time.Now()
	// Nutze ParseSearchQuery für GraphQL-ähnliche Queries
	queryObj, tableName, err := ParseSearchQuery(r)
	if err != nil {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusBadRequest, "ERR_INVALID_QUERY", err.Error(), execTime)
	}

	meta, err := fields.LoadTableMeta(h.DataDir, dbName, tableName)
	if err != nil {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusNotFound, "ERR_META_NOT_FOUND", err.Error(), execTime)
	}
	entries, err := h.queryWithJoins(dbName, tableName, queryObj, meta, queryObj.Join, 1, 8, nil)
	if err != nil {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusInternalServerError, "ERR_QUERY_EXEC", err.Error(), execTime)
	}
	execTime := time.Since(start).String()
	return response.WriteSuccess(entries, execTime)
}

// Exportiere die Funktion, damit sie im Router verwendet werden kann
func (h *QueryHandler) QueryWithJoins(dbName, tableName string, query *Query, meta *fields.TableMeta, joinDefs []JoinDef, joinDepth, maxJoinDepth int, joinPath map[string]struct{}) ([]map[string]interface{}, error) {
	return h.queryWithJoins(dbName, tableName, query, meta, joinDefs, joinDepth, maxJoinDepth, joinPath)
}
