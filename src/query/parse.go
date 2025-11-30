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

	// Robuste Extraktion des Parameterblocks (Klammerzählung)
	paramsStart := tableStart + tableEnd + 1
	parenCount := 1
	paramsEnd := paramsStart
	for paramsEnd < len(queryStr) && parenCount > 0 {
		if queryStr[paramsEnd] == '(' {
			parenCount++
		} else if queryStr[paramsEnd] == ')' {
			parenCount--
		}
		paramsEnd++
	}
	if parenCount != 0 {
		return nil, "", errors.New("Parameterblock nicht gefunden (Klammern nicht ausgeglichen)")
	}
	paramsBlock := queryStr[paramsStart : paramsEnd-1]

	// Logging für Debugging
	fmt.Println("[DEBUG] paramsBlock:", paramsBlock)

	// Extrahiere Felder (zwischen erstem '{' nach ')' und passender '}')
	fieldsStart := strings.Index(queryStr[paramsEnd:], "{")
	fieldsEnd := strings.LastIndex(queryStr, "}")
	if fieldsStart == -1 || fieldsEnd == -1 {
		return nil, "", errors.New("Feldblock nicht gefunden")
	}
	fieldsBlock := queryStr[paramsEnd+fieldsStart+1 : fieldsEnd]

	// 3. Parameter-Block in JSON-ähnliches Format umwandeln (robuster)
	paramsBlock = strings.ReplaceAll(paramsBlock, "\n", " ")
	paramsBlock = strings.ReplaceAll(paramsBlock, "\t", " ")
	paramsBlock = strings.ReplaceAll(paramsBlock, "'", "\"")
	paramsBlock = strings.ReplaceAll(paramsBlock, ",", ", ")
	paramsBlock = strings.ReplaceAll(paramsBlock, "  ", " ")
	paramsBlock = strings.ReplaceAll(paramsBlock, "=", ": ")
	paramsBlock = strings.ReplaceAll(paramsBlock, "True", "true")
	paramsBlock = strings.ReplaceAll(paramsBlock, "False", "false")
	paramsBlock = strings.ReplaceAll(paramsBlock, "None", "null")
	paramsBlock = strings.ReplaceAll(paramsBlock, "\r", " ")
	paramsBlock = strings.TrimSpace(paramsBlock)
	paramsBlock = replaceColonsAndCommasOutsideStrings(paramsBlock)
	paramsBlock = quoteKeys(paramsBlock)

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
		return nil, "", fmt.Errorf("Parameterblock konnte nicht geparst werden: %w\nparamsBlock: %s\nparamsJSON: %s", err, paramsBlock, paramsJSON)
	}
	qr.Fields = fields
	qr.IsSearchQuery = true

	return &qr, tableName, nil
}

// Verbesserte Hilfsfunktion: Ersetze Doppelpunkte und Kommas nur außerhalb von Strings und Arrays
func replaceColonsAndCommasOutsideStrings(s string) string {
	var result strings.Builder
	inString := false
	bracketDepth := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '"' {
			inString = !inString
		}
		if !inString {
			if c == '{' || c == '[' {
				bracketDepth++
			}
			if c == '}' || c == ']' {
				bracketDepth--
			}
			if c == ':' {
				result.WriteString(": ")
				continue
			}
			if c == ',' {
				result.WriteString(", ")
				continue
			}
		}
		result.WriteByte(c)
	}
	return result.String()
}

// Robuste Hilfsfunktion: Keys in JSON-ähnlichem String quotieren (auch nach Kommas, {, [ und mit Whitespace)
func quoteKeys(s string) string {
	var result strings.Builder
	inString := false
	nextIsKey := true // Nach {, [, oder , erwarten wir einen Key
	for i := 0; i < len(s); {
		if s[i] == '"' {
			inString = !inString
			result.WriteByte(s[i])
			i++
			nextIsKey = false
			continue
		}
		if !inString && (s[i] == '{' || s[i] == '[' || s[i] == ',') {
			result.WriteByte(s[i])
			nextIsKey = true
			i++
			continue
		}
		if !inString && (s[i] == ' ' || s[i] == '\n' || s[i] == '\t' || s[i] == '\r') {
			result.WriteByte(s[i])
			i++
			continue
		}
		if !inString && nextIsKey && ((s[i] >= 'a' && s[i] <= 'z') || (s[i] >= 'A' && s[i] <= 'Z') || s[i] == '_') {
			start := i
			for i < len(s) && (s[i] == '_' || (s[i] >= 'a' && s[i] <= 'z') || (s[i] >= 'A' && s[i] <= 'Z') || (s[i] >= '0' && s[i] <= '9')) {
				i++
			}
			// Whitespace zwischen Key und : überspringen
			for i < len(s) && (s[i] == ' ' || s[i] == '\n' || s[i] == '\t' || s[i] == '\r') {
				result.WriteByte(s[i])
				i++
			}
			if i < len(s) && s[i] == ':' {
				result.WriteByte('"')
				result.WriteString(s[start:i])
				result.WriteByte('"')
				nextIsKey = false
				continue
			} else {
				result.WriteString(s[start:i])
				nextIsKey = false
				continue
			}
		}
		result.WriteByte(s[i])
		nextIsKey = false
		i++
	}
	return result.String()
}
