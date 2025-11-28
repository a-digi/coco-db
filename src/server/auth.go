// Optionale Authentifizierungs- und Autorisierungslogik für den Server
// Hier wird ein einfacher Authentifizierungsmechanismus vorbereitet (z.B. Token-basiert)

package server

import (
	"net/http"
	"github.com/a-digi/coco-db/src/response"
	"encoding/json"
)

// SimpleTokenAuth ist eine Beispielimplementierung für AuthProvider
// In einer echten Anwendung sollte dies durch ein sicheres Verfahren ersetzt werden

type SimpleTokenAuth struct {
	ValidTokens map[string]string // token -> userID
}

func (a *SimpleTokenAuth) Authenticate(token string) (string, error) {
	userID, ok := a.ValidTokens[token]
	if !ok {
		return "", http.ErrNoCookie // als Beispiel für Fehler
	}
	return userID, nil
}

// AuthMiddleware prüft das Vorhandensein und die Gültigkeit eines Tokens
func AuthMiddleware(auth AuthProvider, next func(http.ResponseWriter, *http.Request) *response.APIResponse) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			resp := response.WriteErrorInternal(http.StatusUnauthorized, "unauthorized", "Kein Auth-Token übergeben", "")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.HttpCode)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		_, err := auth.Authenticate(token)
		if err != nil {
			resp := response.WriteErrorInternal(http.StatusUnauthorized, "unauthorized", "Authentifizierung fehlgeschlagen", "")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(resp.HttpCode)
			_ = json.NewEncoder(w).Encode(resp)
			return
		}
		next(w, r)
	}
}
