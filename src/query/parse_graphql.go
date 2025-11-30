package query

import (
	"fmt"
)

// ParseGraphQLStyleQuery konvertiert ein beliebiges map[string]interface{} (aus dem GraphQL-Stil-Request)
// rekursiv in das interne Query-Struct.
func ParseGraphQLStyleQuery(obj map[string]interface{}) (*Query, error) {
	q := &Query{}
	for k, v := range obj {
		switch k {
		case "filter":
			if m, ok := v.(map[string]interface{}); ok {
				q.Filter = m
			} else {
				return nil, fmt.Errorf("filter muss ein Objekt sein")
			}
		case "limit":
			if f, ok := v.(float64); ok {
				q.Limit = int(f)
			} else if i, ok := v.(int); ok {
				q.Limit = i
			} else {
				return nil, fmt.Errorf("limit muss eine Zahl sein")
			}
		case "offset":
			if f, ok := v.(float64); ok {
				q.Offset = int(f)
			} else if i, ok := v.(int); ok {
				q.Offset = i
			} else {
				return nil, fmt.Errorf("offset muss eine Zahl sein")
			}
		case "sort":
			if arr, ok := v.([]interface{}); ok {
				q.Sort = make([]string, len(arr))
				for i, s := range arr {
					if str, ok := s.(string); ok {
						q.Sort[i] = str
					} else {
						return nil, fmt.Errorf("sort[%d] muss ein String sein", i)
					}
				}
			} else {
				return nil, fmt.Errorf("sort muss ein Array von Strings sein")
			}
		case "join":
			if arr, ok := v.([]interface{}); ok {
				joins := make([]JoinDef, 0, len(arr))
				for i, j := range arr {
					joinMap, ok := j.(map[string]interface{})
					if !ok {
						return nil, fmt.Errorf("join[%d] muss ein Objekt sein", i)
					}
					joinDef, err := parseJoinDef(joinMap)
					if err != nil {
						return nil, fmt.Errorf("join[%d]: %v", i, err)
					}
					joins = append(joins, *joinDef)
				}
				q.Join = joins
			} else {
				return nil, fmt.Errorf("join muss ein Array sein")
			}
		}
	}
	return q, nil
}

func parseJoinDef(obj map[string]interface{}) (*JoinDef, error) {
	j := &JoinDef{}
	for k, v := range obj {
		switch k {
		case "table":
			if s, ok := v.(string); ok {
				j.Table = s
			} else {
				return nil, fmt.Errorf("table muss ein String sein")
			}
		case "on":
			if m, ok := v.(map[string]interface{}); ok {
				j.On = make(map[string]string)
				for k2, v2 := range m {
					if s2, ok := v2.(string); ok {
						j.On[k2] = s2
					} else {
						return nil, fmt.Errorf("on[%s] muss ein String sein", k2)
					}
				}
			} else {
				return nil, fmt.Errorf("on muss ein Objekt sein")
			}
		case "filter":
			if m, ok := v.(map[string]interface{}); ok {
				j.Filter = m
			} else {
				return nil, fmt.Errorf("filter muss ein Objekt sein")
			}
		case "fields":
			if arr, ok := v.([]interface{}); ok {
				j.Fields = make([]string, len(arr))
				for i, s := range arr {
					if str, ok := s.(string); ok {
						j.Fields[i] = str
					} else {
						return nil, fmt.Errorf("fields[%d] muss ein String sein", i)
					}
				}
			} else {
				return nil, fmt.Errorf("fields muss ein Array von Strings sein")
			}
		case "join":
			if arr, ok := v.([]interface{}); ok {
				joins := make([]JoinDef, 0, len(arr))
				for i, j2 := range arr {
					joinMap, ok := j2.(map[string]interface{})
					if !ok {
						return nil, fmt.Errorf("join[%d] muss ein Objekt sein", i)
					}
					joinDef, err := parseJoinDef(joinMap)
					if err != nil {
						return nil, fmt.Errorf("join[%d]: %v", i, err)
					}
					joins = append(joins, *joinDef)
				}
				j.Join = joins
			} else {
				return nil, fmt.Errorf("join muss ein Array sein")
			}
		}
	}
	return j, nil
}
