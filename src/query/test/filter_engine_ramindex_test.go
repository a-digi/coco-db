package query_test

import (
	"testing"
	"github.com/a-digi/coco-db/src/index"
)

func TestFilterEngine_UsesRAMIndex(t *testing.T) {
	dbName := "testdb"
	tableName := "users"
	idxName := "email_idx"
	idxKey := dbName + "." + tableName + "." + idxName

	// Simuliere RAM-Index
	reg := index.GetRegistry()
	reg.Set(idxKey, index.IndexData{"foo@bar.de": []string{"id1", "id2"}})

	// Prüfe, ob der RAM-Index genutzt wird
	idx, ok := reg.Get(idxKey)
	if !ok {
		t.Fatalf("RAM-Index nicht gefunden")
	}
	idsRaw, found := idx["foo@bar.de"]
	if !found {
		t.Fatalf("Key im RAM-Index nicht gefunden")
	}
	ids, ok := idsRaw.([]string)
	if !ok || len(ids) != 2 || ids[0] != "id1" || ids[1] != "id2" {
		t.Errorf("RAM-Index liefert falsche Werte: %+v", idsRaw)
	}
}
