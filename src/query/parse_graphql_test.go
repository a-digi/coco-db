package query

import (
	"testing"
	"reflect"
)

func TestParseGraphQLStyleQuery_Basic(t *testing.T) {
	input := map[string]interface{}{
		"filter": map[string]interface{}{"email": map[string]interface{}{"like": "*@gmail.com"}},
		"limit": float64(10),
		"offset": float64(2),
		"sort": []interface{}{ "email", "id" },
	}
	q, err := ParseGraphQLStyleQuery(input)
	if err != nil {
		t.Fatalf("ParseGraphQLStyleQuery failed: %v", err)
	}
	if q.Limit != 10 || q.Offset != 2 {
		t.Errorf("Limit/Offset falsch: got %d/%d", q.Limit, q.Offset)
	}
	if q.Filter["email"].(map[string]interface{})["like"] != "*@gmail.com" {
		t.Errorf("Filter falsch: %+v", q.Filter)
	}
	if !reflect.DeepEqual(q.Sort, []string{"email", "id"}) {
		t.Errorf("Sort falsch: %+v", q.Sort)
	}
}

func TestParseGraphQLStyleQuery_Join(t *testing.T) {
	input := map[string]interface{}{
		"filter": map[string]interface{}{"active": true},
		"join": []interface{}{
			map[string]interface{}{
				"table": "user_roles",
				"on": map[string]interface{}{"id": "user_id"},
				"fields": []interface{}{ "role_id" },
				"join": []interface{}{
					map[string]interface{}{
						"table": "roles",
						"on": map[string]interface{}{"role_id": "id"},
						"fields": []interface{}{ "name", "description" },
					},
				},
			},
		},
	}
	q, err := ParseGraphQLStyleQuery(input)
	if err != nil {
		t.Fatalf("ParseGraphQLStyleQuery failed: %v", err)
	}
	if len(q.Join) != 1 {
		t.Fatalf("Join nicht erkannt")
	}
	j := q.Join[0]
	if j.Table != "user_roles" || j.On["id"] != "user_id" {
		t.Errorf("Join falsch: %+v", j)
	}
	if len(j.Fields) != 1 || j.Fields[0] != "role_id" {
		t.Errorf("Fields falsch: %+v", j.Fields)
	}
	if len(j.Join) != 1 || j.Join[0].Table != "roles" {
		t.Errorf("Verschachtelter Join fehlt: %+v", j.Join)
	}
	if j.Join[0].Fields[0] != "name" || j.Join[0].Fields[1] != "description" {
		t.Errorf("Fields im verschachtelten Join falsch: %+v", j.Join[0].Fields)
	}
}

func TestParseGraphQLStyleQuery_Errors(t *testing.T) {
	_, err := ParseGraphQLStyleQuery(map[string]interface{}{ "limit": "notanumber" })
	if err == nil {
		t.Error("Fehler erwartet für ungültigen limit-Typ")
	}
	_, err = ParseGraphQLStyleQuery(map[string]interface{}{ "join": "notanarray" })
	if err == nil {
		t.Error("Fehler erwartet für ungültigen join-Typ")
	}
}

