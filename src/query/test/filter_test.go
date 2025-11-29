package query_test

import (
	"testing"

	"github.com/a-digi/coco-db/src/query"
	"github.com/a-digi/coco-db/src/table/fields"
)

func TestFilterEntries_Basic(t *testing.T) {
	meta := &fields.TableMeta{
		TableName: "users",
		Fields: []fields.FieldMeta{
			{Name: "email", Type: "string"},
			{Name: "age", Type: "int"},
			{Name: "isActive", Type: "bool"},
		},
	}
	entries := []map[string]interface{}{
		{"email": "foo@bar.de", "age": 20.0, "isActive": true},
		{"email": "bar@foo.de", "age": 17.0, "isActive": false},
		{"email": "baz@foo.de", "age": 30.0, "isActive": true},
	}

	t.Run("eq filter", func(t *testing.T) {
		q := &query.Query{Filter: map[string]interface{}{"email": "foo@bar.de"}}
		res := query.FilterEntries(entries, q, meta)
		if len(res) != 1 || res[0]["email"] != "foo@bar.de" {
			t.Errorf("expected 1 result, got %+v", res)
		}
	})

	t.Run("gt filter", func(t *testing.T) {
		q := &query.Query{Filter: map[string]interface{}{"age": map[string]interface{}{"gt": 18}}}
		res := query.FilterEntries(entries, q, meta)
		if len(res) != 2 {
			t.Errorf("expected 2 results, got %+v", res)
		}
	})

	t.Run("and logic", func(t *testing.T) {
		q := &query.Query{Filter: map[string]interface{}{"isActive": true, "age": map[string]interface{}{"gte": 18}}}
		res := query.FilterEntries(entries, q, meta)
		if len(res) != 2 {
			t.Errorf("expected 2 results, got %+v", res)
		}
	})

	t.Run("like filter", func(t *testing.T) {
		q := &query.Query{Filter: map[string]interface{}{"email": map[string]interface{}{"like": "foo*@bar.de"}}}
		res := query.FilterEntries(entries, q, meta)
		if len(res) != 1 || res[0]["email"] != "foo@bar.de" {
			t.Errorf("expected 1 result for LIKE, got %+v", res)
		}
	})

	t.Run("partial filter", func(t *testing.T) {
		q := &query.Query{Filter: map[string]interface{}{"email": map[string]interface{}{"partial": "foo"}}}
		res := query.FilterEntries(entries, q, meta)
		if len(res) != 2 {
			t.Errorf("expected 2 results for PARTIAL, got %+v", res)
		}
	})

	t.Run("fulltext filter", func(t *testing.T) {
		q := &query.Query{Filter: map[string]interface{}{"email": map[string]interface{}{"fulltext": "foo bar"}}}
		res := query.FilterEntries(entries, q, meta)
		if len(res) != 1 || res[0]["email"] != "foo@bar.de" {
			t.Errorf("expected 1 result for FULLTEXT, got %+v", res)
		}
	})
}
