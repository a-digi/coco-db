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
	"github.com/a-digi/coco-db/src/table/fields"
	"github.com/a-digi/coco-db/src/table/index"
)

// TableCreator kapselt die Abhängigkeiten für das Anlegen von Tabellen
// und macht die Konfiguration (dataDir, Logger) explizit
//
type TableCreator struct {
	DataDir string
	Logger  logger.Logger
}

// --- Hilfsfunktionen für Validierung (aus edit.go) ---
var TableNamePattern = regexp.MustCompile(`^[a-zA-Z0-9_-]{3,32}$`)

func validateTableName(name string) error {
	if !TableNamePattern.MatchString(name) {
		return fmt.Errorf("Ungültiger Tabellenname")
	}
	reserved := map[string]struct{}{"meta": {}, "entries": {}, "indexes": {}}
	if _, found := reserved[strings.ToLower(name)]; found {
		return fmt.Errorf("Tabellenname ist reserviert")
	}
	return nil
}

func validateFields(fields []fields.FieldMeta) error {
	if len(fields) == 0 {
		return fmt.Errorf("Mindestens ein Feld muss definiert sein")
	}
	fieldNames := make(map[string]struct{})
	for _, f := range fields {
		if f.Name == "" || f.Type == "" {
			return fmt.Errorf("Jedes Feld muss einen Namen und Typ haben")
		}
		if _, exists := fieldNames[f.Name]; exists {
			return fmt.Errorf("Feldname kommt mehrfach vor: %s", f.Name)
		}
		fieldNames[f.Name] = struct{}{}
	}
	return nil
}

// Fehlercode-Wrapper für Validierung
func validateTableNameWithCode(name string) (string, error) {
	if !TableNamePattern.MatchString(name) {
		return "ERR_TABLE_INVALID_NAME", fmt.Errorf("Ungültiger Tabellenname")
	}
	reserved := map[string]struct{}{"meta": {}, "entries": {}, "indexes": {}}
	if _, found := reserved[strings.ToLower(name)]; found {
		return "ERR_TABLE_RESERVED_NAME", fmt.Errorf("Tabellenname ist reserviert")
	}
	return "", nil
}

func validateFieldsWithCode(fields []fields.FieldMeta) (string, error) {
	if len(fields) == 0 {
		return "ERR_FIELDS_MISSING", fmt.Errorf("Mindestens ein Feld muss definiert sein")
	}
	fieldNames := make(map[string]struct{})
	for _, f := range fields {
		if f.Name == "" || f.Type == "" {
			return "ERR_FIELD_INVALID", fmt.Errorf("Jedes Feld muss einen Namen und Typ haben")
		}
		if _, exists := fieldNames[f.Name]; exists {
			return "ERR_FIELD_DUPLICATE", fmt.Errorf("Feldname kommt mehrfach vor: %s", f.Name)
		}
		fieldNames[f.Name] = struct{}{}
	}
	return "", nil
}

// HandleCreateTable verarbeitet das Anlegen einer neuen Tabelle (POST /api/databases/{dbname}/tables)
func (tc *TableCreator) HandleCreateTable(dbname string, meta fields.TableMeta) *response.APIResponse {
	start := time.Now()
	if dbname == "" {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusBadRequest, "ERR_DB_NAME_MISSING", "Datenbankname fehlt", execTime)
	}
	if !TableNamePattern.MatchString(dbname) {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusBadRequest, "ERR_DB_INVALID_NAME", "Ungültiger Datenbankname", execTime)
	}

	if code, err := validateTableNameWithCode(meta.TableName); err != nil {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusBadRequest, code, err.Error(), execTime)
	}
	if code, err := validateFieldsWithCode(meta.Fields); err != nil {
		execTime := time.Since(start).String()
		return response.WriteErrorInternal(http.StatusBadRequest, code, err.Error(), execTime)
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

	// meta.json schreiben
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

	// Indexdefinitionen validieren
	if indexErrors := fields.ValidateIndexes(&meta); len(indexErrors) > 0 {
		execTime := time.Since(start).String()
		errMsg := "Indexdefinitionen ungültig: " + strings.Join(indexErrors, "; ")
		tc.Logger.Error("[TABLE_CREATE] " + errMsg)
		return response.WriteErrorInternal(http.StatusBadRequest, "ERR_INDEX_DEFINITION", errMsg, execTime)
	}

	// Indexdateien initial anlegen (leere Indexstruktur persistieren)
	for _, idxMeta := range meta.Indexes {
		idx := index.NewBTreeIndex(index.IndexMeta{
			Name:   idxMeta.Name,
			Fields: idxMeta.Fields,
			Type:   index.IndexTypeBTree,
			Unique: idxMeta.Unique,
			Sparse: idxMeta.Sparse,
		})
		if err := idx.SaveToFile(tableDir); err != nil {
			tc.Logger.Error("[TABLE_CREATE] Fehler beim Anlegen Indexdatei: " + err.Error())
		}
	}

	// tables.json der Datenbank aktualisieren
	tablesJsonPath := filepath.Join(dbDir, "tables.json")
	var tablesMeta []fields.TableMeta
	if _, err := os.Stat(tablesJsonPath); err == nil {
		content, err := os.ReadFile(tablesJsonPath)
		if err == nil {
			_ = json.Unmarshal(content, &tablesMeta)
		}
	}
	// Füge die neue Tabelle hinzu (ersetzt ggf. vorhandene mit gleichem Namen)
	tableReplaced := false
	for i, t := range tablesMeta {
		if t.TableName == meta.TableName {
			tablesMeta[i] = meta
			tableReplaced = true
			break
		}
	}
	if !tableReplaced {
		tablesMeta = append(tablesMeta, meta)
	}
	f2, err := os.Create(tablesJsonPath)
	if err == nil {
		enc2 := json.NewEncoder(f2)
		enc2.SetIndent("", "  ")
		_ = enc2.Encode(tablesMeta)
		f2.Close()
	}

	tc.Logger.Info(fmt.Sprintf("[TABLE_CREATE] Tabelle %s/%s erfolgreich angelegt", dbname, meta.TableName))
	if metaBytes, err := json.MarshalIndent(meta, "", "  "); err == nil {
		tc.Logger.Info("[TABLE_CREATE] meta.json: " + string(metaBytes))
	}

	execTime := time.Since(start).String()
	return response.WriteSuccess(meta, execTime)
}
