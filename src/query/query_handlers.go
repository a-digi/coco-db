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
	filterResult, err := FilterEngine(h.DataDir, dbName, tableName, query, meta)

	if err != nil {
		return nil, err
	}

	entries := filterResult.Entries
	fileOpens := filterResult.FileOpens
	ramHits := filterResult.RAMHits

	for _, join := range joinDefs {
		joinMeta, err := fields.LoadTableMeta(h.DataDir, dbName, join.Table)
		if err != nil {
			continue
		}
		// Collect all relevant join IDs from parent entries
		joinIDs := make(map[interface{}]struct{})
		for i := range entries {
			for _, src := range join.On {
				val, ok := entries[i][src]
				if ok {
					joinIDs[val] = struct{}{}
				}
			}
		}
		// Build in-filter for join field
		joinFilter := map[string]interface{}{}

		for dst := range join.On {
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

		// Combine filters
		for k, v := range join.Filter {
			joinFilter[k] = v
		}

		joinQuery := &Query{
			Filter: joinFilter,
			Limit:  0,
			Offset: 0,
			Sort:   nil,
			Join:   join.Join,
		}

		// Recursive join with updated joinPath
		joinResult, err := h.queryWithJoins(dbName, join.Table, joinQuery, joinMeta, join.Join, joinDepth+1, maxJoinDepth, copyJoinPath(joinPath))

		if err != nil {
			joinResult = &FilterResult{Entries: []map[string]interface{}{}}
		}

		if joinResult == nil {
			joinResult = &FilterResult{Entries: []map[string]interface{}{}}
		}

		if join.Fields != nil && len(join.Fields) > 0 {
			// Only keep desired fields
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

		// Debug: Log join filter and result count
		if h.Logger != nil {
			h.Logger.Info(fmt.Sprintf("Join: %s, Filter: %+v, Hits: %d", join.Table, joinQuery.Filter, len(joinResult.Entries)))
		}

		// Map join results to parent entries
		for i := range entries {
			var matchList []map[string]interface{}
			for _, e := range joinResult.Entries {
				for dst, src := range join.On {
					if entries[i][src] == e[dst] {
						matchList = append(matchList, e)
					}
				}
			}
			entries[i][join.Table] = matchList
		}

		fileOpens += joinResult.FileOpens
		ramHits += joinResult.RAMHits
		// Nested loop join with RAM index
		reg := index.GetRegistry()
		// Dynamically determine the join index field (instead of assuming: join.Fields[0])
		var joinIndexField string

		if len(join.Fields) > 0 {
			joinIndexField = join.Fields[0]
		} else {
			// Fallback: first field from On-mapping
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
							//println("[NestedLoopJoin] ParentID:", parentValStr, "→ JoinIDs:", idList, "(RAM-Index used: true)")
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
					//println("[NestedLoopJoin] ParentID:", parentVal, "(RAM-Index used: false)")
					for _, e := range joinResult.Entries {
						if entries[i][src] == e[dst] {
							matchList = append(matchList, e)
						}
					}
				}
			}
			entries[i][join.Table] = matchList
		}
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
