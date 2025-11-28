// Zentrale Validierungs- und Sicherheitsmechanismen für den Server
// Hier werden Input-Validierung, Fehlervermeidung und optionale Authentifizierung vorbereitet

package server

import (
	"encoding/json"
	"net/http"
	"regexp"
	"errors"
)

// ValidateJSON prüft, ob der Request-Body valides JSON ist
func ValidateJSON(r *http.Request, v interface{}) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(v); err != nil {
		return errors.New("Ungültiges JSON: " + err.Error())
	}
	return nil
}

// ValidateString prüft, ob ein String ein bestimmtes Pattern erfüllt
func ValidateString(value, pattern string) error {
	re := regexp.MustCompile(pattern)
	if !re.MatchString(value) {
		return errors.New("String entspricht nicht dem geforderten Pattern")
	}
	return nil
}

// WriteErrorResponse gibt eine standardisierte Fehlermeldung im JSON-Format zurück
func WriteErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// (Optional) Authentifizierungs-Stub für spätere Erweiterung
func RequireAuth(auth AuthProvider, w http.ResponseWriter, r *http.Request) (userID string, ok bool) {
	token := r.Header.Get("Authorization")
	if token == "" {
		WriteErrorResponse(w, http.StatusUnauthorized, "Kein Auth-Token übergeben")
		return "", false
	}
	userID, err := auth.Authenticate(token)
	if err != nil {
		WriteErrorResponse(w, http.StatusUnauthorized, "Authentifizierung fehlgeschlagen")
		return "", false
	}
	return userID, true
}

