package index

import (
	"encoding/json"
	"os"
	"path/filepath"

	fields "github.com/a-digi/coco-db/src/table/fields"
)

// BuildIndexesFromEntries baut alle Indizes einer Tabelle aus den vorhandenen Einträgen neu auf.
// tableDir: Verzeichnis der Tabelle (enthält meta.json, entries/)
// meta: Metadaten der Tabelle (inkl. Indexdefinitionen)
func BuildIndexesFromEntries(tableDir string, meta fields.TableMeta) error {
	entriesDir := filepath.Join(tableDir, "entries")
	files, err := os.ReadDir(entriesDir)
	if err != nil {
		return err
	}
	for _, idxMeta := range meta.Indexes {
		idx := NewBTreeIndex(IndexMeta{
			Name:   idxMeta.Name,
			Fields: idxMeta.Fields,
			Type:   IndexTypeBTree,
			Unique: idxMeta.Unique,
			Sparse: idxMeta.Sparse,
		})
		for _, entryFolder := range files {
			if !entryFolder.IsDir() {
				continue
			}
			entryId := entryFolder.Name()
			entryFile := filepath.Join(entriesDir, entryId, entryId+".json")
			data, err := os.ReadFile(entryFile)
			if err != nil {
				continue // Fehlerhafte Einträge überspringen
			}
			var entry map[string]interface{}
			if err := json.Unmarshal(data, &entry); err != nil {
				continue
			}
			// Extrahiere Indexwert (nur Single-Field-Index für MVP)
			if len(idxMeta.Fields) == 1 {
				key := entry[idxMeta.Fields[0]]
				if key != nil {
					idx.Insert(key, entryId)
				}
			}
			// TODO: Multi-Field-Index (Composite Key) später
		}
		// Index speichern
		_ = idx.SaveToFile(tableDir)
	}
	return nil
}

