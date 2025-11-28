package table

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/response"
)

// FieldMeta beschreibt ein Feld in meta.json gemäß PROJECT_REQUIREMENT.md
// Unterstützt alle geforderten Typen und Constraints
type FieldMeta struct {
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Required    bool          `json:"required,omitempty"`
	Default     interface{}   `json:"default,omitempty"`
	MinLength   *int          `json:"minLength,omitempty"`
	MaxLength   *int          `json:"maxLength,omitempty"`
	Min         *float64      `json:"min,omitempty"`
	Max         *float64      `json:"max,omitempty"`
	Nullable    *bool         `json:"nullable,omitempty"`
	Pattern     string        `json:"pattern,omitempty"`
	Enum        []interface{} `json:"enum,omitempty"`
	Description string        `json:"description,omitempty"`
}

// IndexMeta beschreibt einen Index in meta.json
type IndexMeta struct {
	Name   string   `json:"name"`
	Type   string   `json:"type"` // "primary" oder "secondary"
	Fields []string `json:"fields"`
	Unique bool     `json:"unique,omitempty"`
	Sparse bool     `json:"sparse,omitempty"`
}

// TableMeta entspricht exakt dem meta.json-Schema laut PROJECT_REQUIREMENT.md
type TableMeta struct {
	TableName            string                 `json:"tableName"`
	SchemaVersion        int                    `json:"schemaVersion"`
	Fields               []FieldMeta            `json:"fields"`
	Indexes              []IndexMeta            `json:"indexes,omitempty"`
	Options              map[string]interface{} `json:"options,omitempty"`
	AllowAdditionalFields *bool                 `json:"allowAdditionalFields,omitempty"`
}

var TableNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,32}$`)

// HandleCreateTable verarbeitet das Anlegen einer neuen Tabelle (POST /api/databases/{dbname}/tables)
func HandleCreateTable(w http.ResponseWriter, r *http.Request) {
	dbname := r.Header.Get("X-DB-Name")
	if dbname == "" {
		response.WriteError(w, http.StatusBadRequest, "ERR_DB_NAME_MISSING", "Datenbankname fehlt", "")
		return
	}
	if !TableNamePattern.MatchString(dbname) {
		response.WriteError(w, http.StatusBadRequest, "ERR_DB_INVALID_NAME", "Ungültiger Datenbankname", "")
		return
	}

	var meta TableMeta
	if err := json.NewDecoder(r.Body).Decode(&meta); err != nil {
		response.WriteError(w, http.StatusBadRequest, "ERR_INVALID_JSON", "Ungültiges JSON: "+err.Error(), "")
		return
	}

	if !TableNamePattern.MatchString(meta.TableName) {
		response.WriteError(w, http.StatusBadRequest, "ERR_TABLE_INVALID_NAME", "Ungültiger Tabellenname", "")
		return
	}
	if strings.ToLower(meta.TableName) == "meta" || strings.ToLower(meta.TableName) == "entries" || strings.ToLower(meta.TableName) == "indexes" {
		response.WriteError(w, http.StatusBadRequest, "ERR_TABLE_RESERVED_NAME", "Tabellenname ist reserviert", "")
		return
	}

	// Felder-Validierung: Mindestens ein Feld, alle Pflichtfelder müssen Namen und Typ haben
	if len(meta.Fields) == 0 {
		response.WriteError(w, http.StatusBadRequest, "ERR_FIELDS_MISSING", "Mindestens ein Feld muss definiert sein", "")
		return
	}
	for _, f := range meta.Fields {
		if f.Name == "" || f.Type == "" {
			response.WriteError(w, http.StatusBadRequest, "ERR_FIELD_INVALID", "Jedes Feld muss einen Namen und Typ haben", "")
			return
		}
	}

	dataDir := "./data" // TODO: Aus Konfiguration laden
	dbDir := filepath.Join(dataDir, dbname)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			response.WriteError(w, http.StatusInternalServerError, "ERR_IO", "Fehler beim Anlegen des Datenbankverzeichnisses: "+err.Error(), "")
			logger.LogError(fmt.Sprintf("[TABLE_CREATE] Fehler beim Anlegen DB-Verzeichnis: %v", err))
			return
		}
	}
	tableDir := filepath.Join(dbDir, meta.TableName)
	if _, err := os.Stat(tableDir); err == nil {
		response.WriteError(w, http.StatusConflict, "ERR_TABLE_EXISTS", "Tabelle existiert bereits", "")
		logger.LogWarn(fmt.Sprintf("[TABLE_CREATE] Tabelle %s/%s existiert bereits", dbname, meta.TableName))
		return
	}
	if err := os.MkdirAll(filepath.Join(tableDir, "entries"), 0755); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "ERR_IO", "Fehler beim Anlegen des entries-Verzeichnisses: "+err.Error(), "")
		logger.LogError(fmt.Sprintf("[TABLE_CREATE] Fehler beim Anlegen entries: %v", err))
		return
	}
	if err := os.MkdirAll(filepath.Join(tableDir, "indexes"), 0755); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "ERR_IO", "Fehler beim Anlegen des indexes-Verzeichnisses: "+err.Error(), "")
		logger.LogError(fmt.Sprintf("[TABLE_CREATE] Fehler beim Anlegen indexes: %v", err))
		return
	}

	metaPath := filepath.Join(tableDir, "meta.json")
	if meta.SchemaVersion == 0 {
		meta.SchemaVersion = 1
	}
	f, err := os.Create(metaPath)
	if err != nil {
		response.WriteError(w, http.StatusInternalServerError, "ERR_IO", "Fehler beim Anlegen der meta.json: "+err.Error(), "")
		logger.LogError(fmt.Sprintf("[TABLE_CREATE] Fehler beim Anlegen meta.json: %v", err))
		return
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(meta); err != nil {
		response.WriteError(w, http.StatusInternalServerError, "ERR_IO", "Fehler beim Schreiben der meta.json: "+err.Error(), "")
		logger.LogError(fmt.Sprintf("[TABLE_CREATE] Fehler beim Schreiben meta.json: %v", err))
		return
	}

	logger.LogInfo(fmt.Sprintf("[TABLE_CREATE] Tabelle %s/%s erfolgreich angelegt", dbname, meta.TableName))
	// Logge die vollständige meta.json als JSON-String
	if metaBytes, err := json.MarshalIndent(meta, "", "  "); err == nil {
		logger.LogInfo("[TABLE_CREATE] meta.json: " + string(metaBytes))
	}

	// Erfolgsantwort: gesamte meta.json zurückgeben
	response.WriteSuccess(w, meta, "Tabelle erfolgreich angelegt")
}
