package table

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/a-digi/coco-db/src/logger"
	"github.com/a-digi/coco-db/src/response"
)

// TableCreator kapselt die Abhängigkeiten für das Anlegen von Tabellen
// und macht die Konfiguration (dataDir, Logger) explizit
//
type TableCreator struct {
	DataDir string
	Logger  logger.Logger
}

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
func (tc *TableCreator) HandleCreateTable(dbname string, meta TableMeta) *response.APIResponse {
	start := time.Now()
	if dbname == "" {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusBadRequest, "ERR_DB_NAME_MISSING", "Datenbankname fehlt", execTime)
	}
	if !TableNamePattern.MatchString(dbname) {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusBadRequest, "ERR_DB_INVALID_NAME", "Ungültiger Datenbankname", execTime)
	}

	if !TableNamePattern.MatchString(meta.TableName) {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusBadRequest, "ERR_TABLE_INVALID_NAME", "Ungültiger Tabellenname", execTime)
	}
	if strings.ToLower(meta.TableName) == "meta" || strings.ToLower(meta.TableName) == "entries" || strings.ToLower(meta.TableName) == "indexes" {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusBadRequest, "ERR_TABLE_RESERVED_NAME", "Tabellenname ist reserviert", execTime)
	}

	// Felder-Validierung: Mindestens ein Feld, alle Pflichtfelder müssen Namen und Typ haben
	if len(meta.Fields) == 0 {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusBadRequest, "ERR_FIELDS_MISSING", "Mindestens ein Feld muss definiert sein", execTime)
	}
	for _, f := range meta.Fields {
		if f.Name == "" || f.Type == "" {
			execTime := time.Since(start).String()
			return response.WriteErrorInternal(http.StatusBadRequest, "ERR_FIELD_INVALID", "Jedes Feld muss einen Namen und Typ haben", execTime)
		}
	}

	dbDir := filepath.Join(tc.DataDir, dbname)
	if _, err := os.Stat(dbDir); os.IsNotExist(err) {
		if err := os.MkdirAll(dbDir, 0755); err != nil {
			execTime := time.Since(start).String()
			tc.Logger.Error(fmt.Sprintf("[TABLE_CREATE] Fehler beim Anlegen DB-Verzeichnis: %v", err))
			return response.WriteErrorInternal(http.StatusInternalServerError, "ERR_IO", "Fehler beim Anlegen des Datenbankverzeichnisses: "+err.Error(), execTime)
		}
	}
	tableDir := filepath.Join(dbDir, meta.TableName)
	if _, err := os.Stat(tableDir); err == nil {
		execTime := time.Since(start).String()
		tc.Logger.Warning(fmt.Sprintf("[TABLE_CREATE] Tabelle %s/%s existiert bereits", dbname, meta.TableName))
		return response.WriteErrorInternal(http.StatusConflict, "ERR_TABLE_EXISTS", "Tabelle existiert bereits", execTime)
	}
	if err := os.MkdirAll(filepath.Join(tableDir, "entries"), 0755); err != nil {
		execTime := time.Since(start).String()
		tc.Logger.Error(fmt.Sprintf("[TABLE_CREATE] Fehler beim Anlegen entries: %v", err))
		return response.WriteErrorInternal(http.StatusInternalServerError, "ERR_IO", "Fehler beim Anlegen des entries-Verzeichnisses: "+err.Error(), execTime)
	}
	if err := os.MkdirAll(filepath.Join(tableDir, "indexes"), 0755); err != nil {
		execTime := time.Since(start).String()
		tc.Logger.Error(fmt.Sprintf("[TABLE_CREATE] Fehler beim Anlegen indexes: %v", err))
		return response.WriteErrorInternal(http.StatusInternalServerError, "ERR_IO", "Fehler beim Anlegen des indexes-Verzeichnisses: "+err.Error(), execTime)
	}

	metaPath := filepath.Join(tableDir, "meta.json")
	if meta.SchemaVersion == 0 {
		meta.SchemaVersion = 1
	}
	f, err := os.Create(metaPath)
	if err != nil {
		execTime := time.Since(start).String()
		tc.Logger.Error(fmt.Sprintf("[TABLE_CREATE] Fehler beim Anlegen meta.json: %v", err))
		return response.WriteErrorInternal(http.StatusInternalServerError, "ERR_IO", "Fehler beim Anlegen der meta.json: "+err.Error(), execTime)
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(meta); err != nil {
		execTime := time.Since(start).String()
		tc.Logger.Error(fmt.Sprintf("[TABLE_CREATE] Fehler beim Schreiben meta.json: %v", err))
		return response.WriteErrorInternal(http.StatusInternalServerError, "ERR_IO", "Fehler beim Schreiben der meta.json: "+err.Error(), execTime)
	}

	tc.Logger.Info(fmt.Sprintf("[TABLE_CREATE] Tabelle %s/%s erfolgreich angelegt", dbname, meta.TableName))
	if metaBytes, err := json.MarshalIndent(meta, "", "  "); err == nil {
		tc.Logger.Info("[TABLE_CREATE] meta.json: " + string(metaBytes))
	}

	execTime := time.Since(start).String()
	return response.WriteSuccess(meta, execTime)
}
