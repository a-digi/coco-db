// Optionale Authentifizierungs- und Autorisierungslogik für den Server
// Hier wird ein einfacher Authentifizierungsmechanismus vorbereitet (z.B. Token-basiert)

package server

import (
	"net/http"
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
func AuthMiddleware(auth AuthProvider, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token == "" {
			WriteErrorResponse(w, http.StatusUnauthorized, "Kein Auth-Token übergeben")
			return
		}
		_, err := auth.Authenticate(token)
		if err != nil {
			WriteErrorResponse(w, http.StatusUnauthorized, "Authentifizierung fehlgeschlagen")
			return
		}
		next.ServeHTTP(w, r)
	})
}

