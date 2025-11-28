package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParamRouter_MatchAndExtractParams(t *testing.T) {
	router := NewParamRouter()
	called := false

	router.HandleFunc("GET", "/api/databases/{dbname}/tables/{tname}", func(w http.ResponseWriter, r *http.Request, params map[string]string) {
		called = true
		if params["dbname"] != "testdb" {
			t.Errorf("dbname param falsch: %v", params["dbname"])
		}
		if params["tname"] != "users" {
			t.Errorf("tname param falsch: %v", params["tname"])
		}
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	req := httptest.NewRequest("GET", "/api/databases/testdb/tables/users", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != 200 {
		t.Errorf("StatusCode = %d, erwartet 200", resp.StatusCode)
	}
	if !called {
		t.Error("Handler wurde nicht aufgerufen")
	}
}

func TestParamRouter_NotFound(t *testing.T) {
	router := NewParamRouter()
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/api/unknown/path", nil)
	router.ServeHTTP(w, req)
	resp := w.Result()
	if resp.StatusCode != 404 {
		t.Errorf("StatusCode = %d, erwartet 404", resp.StatusCode)
	}
}

