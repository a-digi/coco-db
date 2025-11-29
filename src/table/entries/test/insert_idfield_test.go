package entries_test

import (
	"testing"

	"github.com/a-digi/coco-db/src/logger"
	entries "github.com/a-digi/coco-db/src/table/entries"
)

func TestInsertEntry_ForbiddenIDField(t *testing.T) {
	tmpDir := t.TempDir()
	dbName := "testdb"
	tableName := "testtable"

	cases := []struct {
		name  string
		entry map[string]interface{}
	}{
		{"ID field upper", map[string]interface{}{"ID": 1}},
		{"id field lower", map[string]interface{}{"id": 1}},
		{"both id fields", map[string]interface{}{"id": 1, "ID": 2}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			resp := entries.InsertEntry(dbName, tableName, c.entry, tmpDir, &logger.NoopLogger{})
			if resp == nil || resp.HttpCode != 400 {
				t.Errorf("expected HttpCode 400 for forbidden ID field, got: %+v", resp)
			}
			if resp != nil && resp.Data != nil {
				if _, ok := resp.Data.(map[string]interface{})["entryId"]; ok {
					t.Errorf("entryId should not be returned on forbidden ID field")
				}
			}
			if resp != nil && resp.Error != nil {
				if resp.Error.Code != "ERR_FORBIDDEN_ID_FIELD" {
					t.Errorf("expected error code ERR_FORBIDDEN_ID_FIELD, got: %v", resp.Error.Code)
				}
			}
		})
	}
}
