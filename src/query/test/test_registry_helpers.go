package query_test

import "github.com/a-digi/coco-db/src/query"

func setupTestRegistryForUsers() {
	idxKey := "testdb.users.id_idx"
	reg := query.GetTestRegistry()
	reg.Set(idxKey, map[string]interface{}{
		"1": []string{"1"},
	})
}

