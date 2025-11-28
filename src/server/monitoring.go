// Monitoring-Stub für Health-Checks und Metriken
// Hier können später Health-Checks, Prometheus-Integration o.ä. ergänzt werden

package server

import (
	"net/http"
	"encoding/json"
	"github.com/a-digi/coco-db/src/response"
)

// HealthHandler liefert einen einfachen Health-Status
func HealthHandler(w http.ResponseWriter, r *http.Request) *response.APIResponse {
	resp := response.WriteSuccess(map[string]string{"status": "ok"}, "")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(resp)
	return resp
}
