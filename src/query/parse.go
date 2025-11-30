package query

import (
	"encoding/json"
	"net/http"
	"fmt"
	"errors"
	"io"
	"strings"
)

// ParseQuery liest und parst den JSON-Body in ein Query-Objekt und prüft Pflichtfelder/Defaults
func ParseQuery(r *http.Request, defaultLimit, maxLimit int) (*Query, error) {
	var qr Query
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&qr); err != nil {
		return nil, fmt.Errorf("Malformed JSON: %w", err)
	}
	// Pflichtfeld Filter prüfen
	if qr.Filter == nil {
		return nil, fmt.Errorf("Missing required field: filter")
	}
	// Limit prüfen und ggf. setzen
	if qr.Limit < 0 {
		return nil, fmt.Errorf("Limit must be >= 0")
	}
	if qr.Limit == 0 && defaultLimit > 0 {
		qr.Limit = defaultLimit
	}
	if maxLimit > 0 && qr.Limit > maxLimit {
		qr.Limit = maxLimit
	}
	// Offset prüfen
	if qr.Offset < 0 {
		return nil, fmt.Errorf("Offset must be >= 0")
	}
	// Sort prüfen (falls vorhanden)
	if qr.Sort != nil {
		for i, s := range qr.Sort {
			if s == "" {
				return nil, fmt.Errorf("Sort[%d] must be a non-empty string", i)
			}
		}
	}
	// Join prüfen (falls vorhanden)
	if qr.Join != nil {
		for i, j := range qr.Join {
			if j.Table == "" {
				return nil, fmt.Errorf("Join[%d] missing required field: table", i)
			}
			if len(j.On) == 0 {
				return nil, fmt.Errorf("Join[%d] missing required field: on", i)
			}
			// Rekursive Validierung für verschachtelte Joins
			if j.Join != nil {
				for k, sub := range j.Join {
					if sub.Table == "" {
						return nil, fmt.Errorf("Join[%d].Join[%d] missing required field: table", i, k)
					}
					if len(sub.On) == 0 {
						return nil, fmt.Errorf("Join[%d].Join[%d] missing required field: on", i, k)
					}
				}
			}
		}
	}
	return &qr, nil
}

// ParseSearchQuery parst eine Query im GraphQL-ähnlichen Format wie:
// query { users( ... ) { ... } }
// und wandelt sie in ein Query-Objekt um.
func ParseSearchQuery(r *http.Request) (*Query, string, error) {
	// 1. Body als String einlesen
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = io.ReadAll(r.Body)
	}
	queryStr := string(bodyBytes)

	// 2. Query-String parsen (explizit für das Format: query { users( ... ) { ... } })
	queryStr = strings.TrimSpace(queryStr)
	if !strings.HasPrefix(queryStr, "query {") {
		return nil, "", errors.New("Query muss mit 'query {' beginnen")
	}
	// Extrahiere Tabellennamen
	tableStart := strings.Index(queryStr, "{") + 1
	tableEnd := strings.Index(queryStr[tableStart:], "(")
	if tableEnd == -1 {
		return nil, "", errors.New("Tabellenname und Parameter erwartet (users(...)")
	}
	tableName := strings.TrimSpace(queryStr[tableStart : tableStart+tableEnd])

	// Extrahiere Parameter (zwischen erstem '(' und erstem ')')
	paramsStart := tableStart + tableEnd + 1
	paramsEnd := strings.Index(queryStr[paramsStart:], ")")
	if paramsEnd == -1 {
		return nil, "", errors.New("Parameterblock nicht gefunden")
	}
	paramsBlock := queryStr[paramsStart : paramsStart+paramsEnd]

	// Extrahiere Felder (zwischen erstem '{' nach ')' und passender '}')
	fieldsStart := strings.Index(queryStr[paramsStart+paramsEnd:], "{")
	fieldsEnd := strings.LastIndex(queryStr, "}")
	if fieldsStart == -1 || fieldsEnd == -1 {
		return nil, "", errors.New("Feldblock nicht gefunden")
	}
	fieldsBlock := queryStr[paramsStart+paramsEnd+fieldsStart+1 : fieldsEnd]

	// 3. Parameter-Block in JSON-ähnliches Format umwandeln (vereinfachte Annahme)
	paramsBlock = strings.ReplaceAll(paramsBlock, "\n", " ")
	paramsBlock = strings.ReplaceAll(paramsBlock, "\t", " ")
	paramsBlock = strings.ReplaceAll(paramsBlock, "'", "\"")
	paramsBlock = strings.ReplaceAll(paramsBlock, ":", ": ")
	paramsBlock = strings.ReplaceAll(paramsBlock, ",", ", ")

	// 4. Felder extrahieren (durch Komma getrennt, ggf. mit Subfeldern)
	fields := []string{}
	for _, f := range strings.Split(fieldsBlock, "\n") {
		f = strings.TrimSpace(f)
		if f != "" && !strings.HasSuffix(f, "{") && !strings.HasSuffix(f, "}") {
			fields = append(fields, f)
		}
	}

	// 5. Parameter-Block als JSON parsen
	paramsJSON := "{" + paramsBlock + "}"
	var qr Query
	if err := json.Unmarshal([]byte(paramsJSON), &qr); err != nil {
		return nil, "", fmt.Errorf("Parameterblock konnte nicht geparst werden: %w", err)
	}
	qr.Fields = fields

	return &qr, tableName, nil
}
