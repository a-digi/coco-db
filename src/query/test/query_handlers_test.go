package query_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"fmt"
	"path/filepath"

	"github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/query"
	"github.com/a-digi/coco-db/src/table/fields"
)

// Entferne setupTestTable aus dieser Datei, da sie jetzt in test_helpers.go liegt

func TestTableQueryHandler_Success(t *testing.T) {
	h := &query.QueryHandler{DataDir: "./testdata", Logger: &logger.NoopLogger{}}
	// Nutze die importierte Funktion zum Setup der Testtabelle
	setupTestTable("./testdata", "testdb", "users", t)
	setupTestRegistryForUsers() // RAM-Index für Users-Tabelle setzen
	t.Cleanup(func() { os.RemoveAll("./testdata") })
	queryBody := map[string]interface{}{
		"filter": map[string]interface{}{"id": "1"},
	}
	body, _ := json.Marshal(queryBody)
	r := httptest.NewRequest("POST", "/api/databases/testdb/tables/users/query", bytes.NewReader(body))
	resp := h.TableQueryHandler("testdb", "users", r)
	if !resp.Success || resp.HttpCode != http.StatusOK {
		t.Errorf("Erwartet Success und Status 200, bekommen: %+v", resp)
	}
	data, ok := resp.Data.([]interface{})
	if !ok || len(data) != 1 {
		t.Errorf("Erwartet ein Eintrag-Array, bekommen: %+v", data)
		return
	}
	entry, ok := data[0].(map[string]interface{})
	if !ok {
		t.Errorf("Erwartet map[string]interface{} für entry, bekommen: %+v", data[0])
		return
	}
	t.Logf("Entry: %+v", entry)
	if entry["id"] != "1" || entry["name"] != "Test" || fmt.Sprintf("%v", entry["age"]) != "42" {
		t.Errorf("Erwartet einen Eintrag mit id=1, name=Test, age=42, bekommen: %+v", entry)
	}
}

func TestTableQueryHandler_InvalidJSON(t *testing.T) {
	h := &query.QueryHandler{DataDir: "./testdata", Logger: &logger.NoopLogger{}}
	r := httptest.NewRequest("POST", "/api/databases/testdb/tables/users/query", bytes.NewReader([]byte("{")))
	resp := h.TableQueryHandler("testdb", "users", r)
	if resp.Success || resp.HttpCode != http.StatusBadRequest {
		t.Errorf("Erwartet Fehler bei ungültigem JSON, bekommen: %+v", resp)
	}
}

func TestTableQueryHandler_NotFound(t *testing.T) {
	h := &query.QueryHandler{DataDir: "./testdata", Logger: &logger.NoopLogger{}}
	// Kein Setup, Tabelle existiert nicht
	r := httptest.NewRequest("POST", "/api/databases/testdb/tables/users/query", bytes.NewReader([]byte(`{"filter":{"id":"1"}}`)))
	resp := h.TableQueryHandler("testdb", "users", r)
	if resp.Success || resp.HttpCode != http.StatusNotFound {
		t.Errorf("Erwartet Fehler bei fehlender Tabelle, bekommen: %+v", resp)
	}
}

func TestQueryHandler_Success(t *testing.T) {
	h := &query.QueryHandler{DataDir: "./testdata", Logger: &logger.NoopLogger{}}
	// Nutze die importierte Funktion zum Setup der Testtabelle
	setupTestTable("./testdata", "testdb", "users", t)
	setupTestRegistryForUsers() // RAM-Index für Users-Tabelle setzen
	t.Cleanup(func() { os.RemoveAll("./testdata") })
	queryBody := map[string]interface{}{
		"filter": map[string]interface{}{"id": "1"},
	}
	body, _ := json.Marshal(queryBody)
	r := httptest.NewRequest("POST", "/api/query", bytes.NewReader(body))
	resp := h.QueryHandler("testdb", r)
	if !resp.Success || resp.HttpCode != http.StatusOK {
		t.Errorf("Erwartet Success und Status 200, bekommen: %+v", resp)
	}
	results, ok := resp.Data.([]interface{})
	if !ok || len(results) == 0 {
		t.Errorf("Erwartet mindestens ein Tabellenergebnis, bekommen: %+v", resp.Data)
		return
	}
	found := false
	var foundEntry map[string]interface{}
	for _, tbl := range results {
		tblMap, ok := tbl.(map[string]interface{})
		if !ok {
			continue
		}
		if tblMap["table"] == "users" {
			entries, ok := tblMap["entries"].([]interface{})
			if !ok {
				continue
			}
			for _, e := range entries {
				entry, ok := e.(map[string]interface{})
				if ok && entry["id"] == "1" {
					found = true
					foundEntry = entry
					break
				}
			}
		}
	}
	if !found {
		t.Errorf("Erwartet einen Eintrag mit id=1 in Tabelle users, bekommen: %+v", resp.Data)
		for i, tbl := range results {
			tblMap, ok := tbl.(map[string]interface{})
			if !ok {
				continue
			}
			t.Logf("Tabelle %d: %v", i, tblMap)
			entries, ok := tblMap["entries"].([]interface{})
			if ok {
				for j, e := range entries {
					entry, ok := e.(map[string]interface{})
					if ok {
						t.Logf("  Entry %d: %v", j, entry)
					}
				}
			}
		}
	} else {
		// Erwarte, dass der Eintrag alle Felder enthält
		if foundEntry["name"] != "Test" || fmt.Sprintf("%v", foundEntry["age"]) != "42" {
			t.Errorf("Gefundener Eintrag stimmt nicht: %+v", foundEntry)
		}
	}
}

func TestQueryHandler_NotFound(t *testing.T) {
	h := &query.QueryHandler{DataDir: "./testdata", Logger: &logger.NoopLogger{}}
	r := httptest.NewRequest("POST", "/api/query", bytes.NewReader([]byte(`{"filter":{"id":"1"}}`)))
	resp := h.QueryHandler("testdb", r)
	if resp.Success || resp.HttpCode != http.StatusNotFound {
		t.Errorf("Erwartet Fehler bei fehlender Tabelle, bekommen: %+v", resp)
	}
}

func TestQueryHandler_Join_Success(t *testing.T) {
	h := &query.QueryHandler{DataDir: "./testdata", Logger: &logger.NoopLogger{}}
	// Setup users-Tabelle
	setupTestTable("./testdata", "testdb", "users", t)
	// Setup orders-Tabelle with Join-Feld user_id
	meta := &fields.TableMeta{
		TableName: "orders",
		Fields: []fields.FieldMeta{
			{Name: "id", Type: "string"},
			{Name: "user_id", Type: "string"},
			{Name: "amount", Type: "int"},
		},
		Indexes: []fields.IndexMeta{
			{Name: "id_idx", Fields: []string{"id"}, Type: "primary", Unique: true},
		},
	}
	metaPath := filepath.Join("./testdata", "testdb", "orders", "meta.json")
	_ = os.MkdirAll(filepath.Dir(metaPath), 0755)
	metaBytes, _ := json.Marshal(meta)
	_ = os.WriteFile(metaPath, metaBytes, 0644)
	// Ein Order-Eintrag
	entriesDir := filepath.Join("./testdata", "testdb", "orders", "entries")
	_ = os.MkdirAll(entriesDir, 0755)
	entry := map[string]interface{}{"id": "1", "user_id": "1", "amount": 99}
	entryBytes, _ := json.Marshal(entry)
	_ = os.WriteFile(filepath.Join(entriesDir, "1.json"), entryBytes, 0644)
	// Indexdatei für orders als Objekt anlegen (id → [id])
	idxDir := filepath.Join("./testdata", "testdb", "orders", "indexes")
	_ = os.MkdirAll(idxDir, 0755)
	idxPath := filepath.Join(idxDir, "index_id_idx.json")
	idxObj := map[string][]string{"1": {"1"}}
	idxBytes, _ := json.Marshal(idxObj)
	_ = os.WriteFile(idxPath, idxBytes, 0644)

	// RAM-Index für users-Tabelle bereitstellen
	userIdxKey := "testdb.users.id_idx"
	userReg := query.GetTestRegistry()
	userReg.Set(userIdxKey, map[string]interface{}{
		"1": []string{"1"},
	})
	// RAM-Index für orders-Tabelle bereitstellen
	orderIdxKey := "testdb.orders.id_idx"
	orderReg := query.GetTestRegistry()
	orderReg.Set(orderIdxKey, map[string]interface{}{
		"1": []string{"1"},
	})

	t.Cleanup(func() { os.RemoveAll("./testdata") })

	queryBody := map[string]interface{}{
		"filter": map[string]interface{}{"id": "1"},
		"join": []interface{}{
			map[string]interface{}{
				"table": "orders",
				"on": map[string]interface{}{"id": "user_id"},
				"fields": []interface{}{ "id", "amount" },
			},
		},
	}
	body, _ := json.Marshal(queryBody)
	r := httptest.NewRequest("POST", "/api/query", bytes.NewReader(body))
	resp := h.QueryHandler("testdb", r)
	if !resp.Success || resp.HttpCode != http.StatusOK {
		t.Errorf("Erwartet Success und Status 200, bekommen: %+v", resp)
	}
	results, ok := resp.Data.([]interface{})
	if !ok || len(results) == 0 {
		t.Errorf("Erwartet mindestens ein Tabellenergebnis, bekommen: %+v", resp.Data)
		return
	}
	var userEntry map[string]interface{}
	for _, tbl := range results {
		tblMap, ok := tbl.(map[string]interface{})
		if !ok || tblMap["table"] != "users" {
			continue
		}
		entries, ok := tblMap["entries"].([]interface{})
		if !ok || len(entries) == 0 {
			t.Errorf("Erwartet mindestens einen User-Eintrag, bekommen: %+v", entries)
			return
		}
		userEntry, ok = entries[0].(map[string]interface{})
		if !ok {
			t.Errorf("User-Eintrag ist kein map[string]interface{}: %+v", entries[0])
			return
		}
		break
	}
	if userEntry == nil {
		t.Errorf("Kein User-Eintrag gefunden")
		return
	}
	orders, ok := userEntry["orders"].([]interface{})
	if !ok || len(orders) == 0 {
		t.Errorf("Erwartet mindestens einen Order-Eintrag im Join, bekommen: %+v", userEntry["orders"])
		return
	}
	order, ok := orders[0].(map[string]interface{})
	if !ok {
		t.Errorf("Order-Eintrag ist kein map[string]interface{}: %+v", orders[0])
		return
	}
	if order["id"] != "1" || fmt.Sprintf("%v", order["amount"]) != "99" {
		t.Errorf("Order-Eintrag stimmt nicht: %+v", order)
	}
}
