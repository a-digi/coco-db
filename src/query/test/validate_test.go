package query_test

import (
	"testing"
	"github.com/a-digi/coco-db/src/query"
	"github.com/a-digi/coco-db/src/table/fields"
)

func TestValidateQueryRequest_Basics(t *testing.T) {
	meta := &fields.TableMeta{
		TableName: "users",
		Fields: []fields.FieldMeta{
			{Name: "email", Type: "string", Required: true},
			{Name: "age", Type: "int"},
			{Name: "isActive", Type: "bool"},
		},
	}

	t.Run("valid simple query", func(t *testing.T) {
		qr := &query.QueryRequest{
			Filter: map[string]interface{}{"email": "foo@bar.de", "age": map[string]interface{}{"gte": 18}},
			Limit:  10,
			Offset: 0,
		}
		errs := query.ValidateQueryRequest(qr, meta)
		if len(errs) != 0 {
			t.Errorf("expected no errors, got: %+v", errs)
		}
	})

	t.Run("missing filter", func(t *testing.T) {
		qr := &query.QueryRequest{Limit: 5, Offset: 0}
		errs := query.ValidateQueryRequest(qr, meta)
		if len(errs) == 0 {
			t.Error("expected error for missing filter")
		}
	})

	t.Run("unknown field", func(t *testing.T) {
		qr := &query.QueryRequest{
			Filter: map[string]interface{}{"foo": "bar"},
			Limit:  1,
		}
		errs := query.ValidateQueryRequest(qr, meta)
		found := false
		for _, e := range errs {
			if e.Code == "ERR_FIELD_NOT_ALLOWED" {
				found = true
			}
		}
		if !found {
			t.Error("expected ERR_FIELD_NOT_ALLOWED error")
		}
	})

	t.Run("type mismatch", func(t *testing.T) {
		qr := &query.QueryRequest{
			Filter: map[string]interface{}{"age": "notanumber"},
		}
		errs := query.ValidateQueryRequest(qr, meta)
		found := false
		for _, e := range errs {
			if e.Code == "ERR_TYPE_MISMATCH" {
				found = true
			}
		}
		if !found {
			t.Error("expected ERR_TYPE_MISMATCH error")
		}
	})

	t.Run("invalid operator", func(t *testing.T) {
		qr := &query.QueryRequest{
			Filter: map[string]interface{}{"email": map[string]interface{}{"gt": "foo@bar.de"}},
		}
		errs := query.ValidateQueryRequest(qr, meta)
		found := false
		for _, e := range errs {
			if e.Code == "ERR_OPERATOR_NOT_ALLOWED" {
				found = true
			}
		}
		if !found {
			t.Error("expected ERR_OPERATOR_NOT_ALLOWED error")
		}
	})
}

