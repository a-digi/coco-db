// Monitoring-Stub für Health-Checks und Metriken
// Hier können später Health-Checks, Prometheus-Integration o.ä. ergänzt werden

package server

import "net/http"

// HealthHandler liefert einen einfachen Health-Status
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

