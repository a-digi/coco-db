package test

import (
	"testing"
	"github.com/a-digi/coco-db/src/table"
	"github.com/a-digi/coco-db/src/logger"
)

func TestHandleCreateTable_InvalidIndexDefinitions(t *testing.T) {
	creator := &table.TableCreator{Logger: &logger.NoopLogger{}}
	cases := []struct {
		name string
		meta table.TableMeta
		wantCode string
	}{
		{
			name: "Doppelter Indexname",
			meta: table.TableMeta{
				TableName: "users",
				Fields: []table.FieldMeta{{Name: "id", Type: "string"}},
				Indexes: []table.IndexMeta{
					{Name: "idx1", Type: "primary", Fields: []string{"id"}},
					{Name: "idx1", Type: "secondary", Fields: []string{"id"}},
				},
			},
			wantCode: "ERR_INDEX_DEFINITION",
		},
		{
			name: "Unbekanntes Feld",
			meta: table.TableMeta{
				TableName: "users",
				Fields: []table.FieldMeta{{Name: "id", Type: "string"}},
				Indexes: []table.IndexMeta{
					{Name: "idx2", Type: "primary", Fields: []string{"notfound"}},
				},
			},
			wantCode: "ERR_INDEX_DEFINITION",
		},
		{
			name: "Ungültiger Index-Typ",
			meta: table.TableMeta{
				TableName: "users",
				Fields: []table.FieldMeta{{Name: "id", Type: "string"}},
				Indexes: []table.IndexMeta{
					{Name: "idx3", Type: "foo", Fields: []string{"id"}},
				},
			},
			wantCode: "ERR_INDEX_DEFINITION",
		},
		{
			name: "Index ohne Felder",
			meta: table.TableMeta{
				TableName: "users",
				Fields: []table.FieldMeta{{Name: "id", Type: "string"}},
				Indexes: []table.IndexMeta{
					{Name: "idx4", Type: "primary", Fields: []string{}},
				},
			},
			wantCode: "ERR_INDEX_DEFINITION",
		},
		{
			name: "Unique und Sparse kombiniert",
			meta: table.TableMeta{
				TableName: "users",
				Fields: []table.FieldMeta{{Name: "id", Type: "string"}},
				Indexes: []table.IndexMeta{
					{Name: "idx5", Type: "primary", Fields: []string{"id"}, Unique: true, Sparse: true},
				},
			},
			wantCode: "ERR_INDEX_DEFINITION",
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resp := creator.HandleCreateTable("testdb", c.meta)
			if resp == nil || resp.Success || resp.Error == nil || resp.Error.Code != c.wantCode {
				t.Errorf("%s: Fehlerfall nicht erkannt: %+v", c.name, resp)
			}
		})
	}
}

