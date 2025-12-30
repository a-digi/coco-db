package query

import (
	"fmt"
	"net/http"
	"time"

	"github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/response"
	"github.com/a-digi/coco-db/src/table/fields"
	"github.com/a-digi/coco-db/src/index"
)

// QueryHandler encapsulates DataDir and Logger for query endpoints
// and enables an object-oriented handler architecture
// Example: handler := &QueryHandler{DataDir: "./data", Logger: &logger.NoopLogger{}}
type QueryHandler struct {
	DataDir string
	Logger  logger.Logger
}

// TableQueryHandler processes a table query and returns an APIResponse
func (h *QueryHandler) TableQueryHandler(dbName, tableName string, r *http.Request) *response.APIResponse {
	defer func() {
		if rec := recover(); rec != nil {
			h.Logger.Error(fmt.Sprintf("panic in TableQueryHandler: %v", rec))
		}
	}()
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
	filterResult, err := FilterEngine(h.DataDir, dbName, tableName, queryObj, meta)
	if err != nil {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusInternalServerError, "ERR_QUERY_EXEC", err.Error(), execTime)
	}
	execTime := time.Since(start).String()
	return response.WriteSuccess(filterResult, execTime)
}

// joinDepth indicates the current join nesting depth (for limitation)
func (h *QueryHandler) queryWithJoins(dbName, tableName string, query *Query, meta *fields.TableMeta, joinDefs []JoinDef, joinDepth, maxJoinDepth int, joinPath map[string]struct{}) (*FilterResult, error) {
	if joinDepth > maxJoinDepth {
		return nil, fmt.Errorf("Maximum join depth (%d) exceeded", maxJoinDepth)
	}
	if joinPath == nil {
		joinPath = make(map[string]struct{})
	}
	pathKey := dbName + "." + tableName
	if _, exists := joinPath[pathKey]; exists {
		return nil, fmt.Errorf("Cyclic join detected: %s", pathKey)
	}
	joinPath[pathKey] = struct{}{}

	startTable := time.Now()
	fmt.Printf("[JOIN-TRACE] Processing table: %s.%s | joinDepth=%d\n", dbName, tableName, joinDepth)
	// Find the results from the main table first
	filterResult, err := FilterEngine(h.DataDir, dbName, tableName, query, meta)
	if err != nil {
		return nil, err
	}

	fmt.Printf("[JOIN-TRACE] Table: %s.%s | Entries: %d | FileOpens: %d | RAMHits: %d | Time: %s\n", dbName, tableName, len(filterResult.Entries), filterResult.FileOpens, filterResult.RAMHits, time.Since(startTable))

	entries := filterResult.Entries
	fileOpens := filterResult.FileOpens
	ramHits := filterResult.RAMHits

    // Process each join definition
	for _, join := range joinDefs {
		joinStart := time.Now()
		fmt.Printf("[JOIN-TRACE] → JOIN: %s ON %+v | Parent entries: %d\n", join.Table, join.On, len(entries))
		// Load metadata for the join table
		joinMeta, err := fields.LoadTableMeta(h.DataDir, dbName, join.Table)

		if err != nil {
			fmt.Printf("[JOIN-TRACE]   [ERROR] Could not load meta for join table %s: %v\n", join.Table, err)
			continue
		}
        // Collect join IDs from the parent entries
		joinIDs := filterJoinIDs(entries, join.On)

		fmt.Printf("[JOIN-TRACE]   Join IDs collected: %d\n", len(joinIDs))
		// Build join query
		joinQuery := buildJoinFilter(join.On, joinIDs, join.Filter, join.Join)

        // Execute the join query recursively
		joinResult, err := h.queryWithJoins(dbName, join.Table, joinQuery, joinMeta, join.Join, joinDepth+1, maxJoinDepth, copyJoinPath(joinPath))
		if err != nil {
			fmt.Printf("[JOIN-TRACE]   [ERROR] Join query failed for %s: %v\n", join.Table, err)
			joinResult = &FilterResult{Entries: []map[string]interface{}{}}
		}
		if joinResult == nil {
			joinResult = &FilterResult{Entries: []map[string]interface{}{}}
		}
		fmt.Printf("[JOIN-TRACE]   Join table: %s | Join results: %d | FileOpens: %d | RAMHits: %d | Time: %s\n", join.Table, len(joinResult.Entries), joinResult.FileOpens, joinResult.RAMHits, time.Since(joinStart))
		if join.Fields != nil && len(join.Fields) > 0 {
			for j := range joinResult.Entries {
				for k := range joinResult.Entries[j] {
					found := false
					for _, f := range join.Fields {
						if k == f {
							found = true
							break
						}
					}
					if !found {
						delete(joinResult.Entries[j], k)
					}
				}
			}
		}
		reg := index.GetRegistry()
		var joinIndexField string
		if len(join.Fields) > 0 {
			joinIndexField = join.Fields[0]
		} else {
			for dst := range join.On {
				joinIndexField = dst
				break
			}
		}
		idxKey := dbName + "." + join.Table + "." + joinIndexField
		idxObj, idxOk := reg.Get(idxKey)
		for i := range entries {
			var matchList []map[string]interface{}
			for dst, src := range join.On {
				parentVal, ok := entries[i][src]
				if !ok {
					continue
				}
				parentValStr, okStr := parentVal.(string)
				if idxOk && okStr {
					if ids, found := idxObj[parentValStr]; found {
						if idList, ok := ids.([]string); ok {
							fmt.Printf("[JOIN-TRACE]   [RAM-INDEX] ParentID: %s → JoinIDs: %v\n", parentValStr, idList)
							for _, id := range idList {
								for _, e := range joinResult.Entries {
									if e[dst] == id {
										matchList = append(matchList, e)
									}
								}
							}
						}
					}
				} else {
					fmt.Printf("[JOIN-TRACE]   [NO-RAM-INDEX] ParentID: %v\n", parentVal)
					for _, e := range joinResult.Entries {
						if entries[i][src] == e[dst] {
							matchList = append(matchList, e)
						}
					}
				}
			}
			entries[i][join.Table] = matchList
		}
		fileOpens += joinResult.FileOpens
		ramHits += joinResult.RAMHits
	}
	return &FilterResult{
		Entries:   entries,
		FileOpens: fileOpens,
		RAMHits:   ramHits,
	}, nil
}

func copyJoinPath(orig map[string]struct{}) map[string]struct{} {
	newMap := make(map[string]struct{}, len(orig))

	for k, v := range orig {
		newMap[k] = v
	}

	return newMap
}

// QueryHandler processes a global query and returns an APIResponse
func (h *QueryHandler) QueryHandler(dbName string, r *http.Request) *response.APIResponse {
	defer func() {
		if rec := recover(); rec != nil {
			h.Logger.Error(fmt.Sprintf("panic in QueryHandler: %v", rec))
		}
	}()
	start := time.Now()
	// Use ParseSearchQuery for GraphQL-like queries
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

	filterResult, err := h.queryWithJoins(dbName, tableName, queryObj, meta, queryObj.Join, 1, 8, nil)

	if err != nil {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusInternalServerError, "ERR_QUERY_EXEC", err.Error(), execTime)
	}

	execTime := time.Since(start).String()

	return response.WriteSuccess(filterResult, execTime)
}

// Export the function so it can be used in the router
func (h *QueryHandler) QueryWithJoins(dbName, tableName string, query *Query, meta *fields.TableMeta, joinDefs []JoinDef, joinDepth, maxJoinDepth int, joinPath map[string]struct{}) ([]map[string]interface{}, error) {
	result, err := h.queryWithJoins(dbName, tableName, query, meta, joinDefs, joinDepth, maxJoinDepth, joinPath)

	if err != nil {
		return nil, err
	}

	return result.Entries, nil
}

// filterJoinIDs extrahiert die Join-IDs aus den Parent-Entries anhand der Join-On Felder
func filterJoinIDs(entries []map[string]interface{}, joinOn map[string]string) map[interface{}]struct{} {
	joinIDs := make(map[interface{}]struct{})
	for i := range entries {
		for _, src := range joinOn {
			val, ok := entries[i][src]
			if ok {
				joinIDs[val] = struct{}{}
			}
		}
	}

	return joinIDs
}

// buildJoinFilter erstellt den Filter für den Join anhand der Join-IDs und mergen mit bestehenden Filtern
// und gibt direkt ein Query-Objekt zurück
func buildJoinFilter(joinOn map[string]string, joinIDs map[interface{}]struct{}, joinFilters map[string]interface{}, joinJoin []JoinDef) *Query {
	joinFilter := make(map[string]interface{})
	for dst := range joinOn {
		var idList []interface{}
		for id := range joinIDs {
			idList = append(idList, id)
		}
		if len(idList) == 1 {
			joinFilter[dst] = idList[0]
		} else if len(idList) > 1 {
			joinFilter[dst] = map[string]interface{}{"in": idList}
		}
	}
	// Merge mit bestehenden Filtern
	for k, v := range joinFilters {
		joinFilter[k] = v
	}
	return &Query{
		Filter: joinFilter,
		Limit:  0,
		Offset: 0,
		Sort:   nil,
		Join:   joinJoin,
	}
}
