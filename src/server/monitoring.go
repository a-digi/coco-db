// Monitoring-Stub für Health-Checks und Metriken
// Hier können später Health-Checks, Prometheus-Integration o.ä. ergänzt werden

package server

import (
	"net/http"
	"github.com/a-digi/coco-db/src/response"
)

// HealthHandler liefert einen einfachen Health-Status
func HealthHandler(w http.ResponseWriter, r *http.Request) *response.APIResponse {
	return response.WriteSuccess(w, map[string]string{"status": "ok"}, "")
}
