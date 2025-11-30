// Zentrale Index-Registry für In-Memory-Indexe
// Diese Registry wird beim Serverstart initialisiert und hält alle geladenen Indexe im RAM.
// Zugriff erfolgt threadsicher über RWMutex.
//
// Nutzung:
//   - registry := index.NewRegistry()
//   - registry.LoadAllIndexes(dataDir)
//   - registry.GetIndex(db, table, indexName)

package index

import (
	"encoding/json"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// IndexObj ist ein Platzhalter für verschiedene Index-Typen (z.B. BTree, Map, ...)
type IndexObj interface{}

// Registry hält alle geladenen Indexe im RAM
// Key: db.table.indexName
// Value: IndexObj (z.B. BTree, Map, ...)
type Registry struct {
	mu     sync.RWMutex
	cache  map[string]IndexObj
}

// NewRegistry erzeugt eine neue leere Registry
func NewRegistry() *Registry {
	return &Registry{
		cache: make(map[string]IndexObj),
	}
}

// LoadAllIndexes lädt alle Indexdateien aus dem Datenverzeichnis in die Registry
func (r *Registry) LoadAllIndexes(dataDir string) error {
	start := time.Now()
	count := 0
	err := filepath.WalkDir(dataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() || !strings.HasPrefix(d.Name(), "index_") || !strings.HasSuffix(d.Name(), ".json") {
			return nil
		}
		// Extrahiere db, table, indexName aus dem Pfad
		parts := strings.Split(filepath.ToSlash(path), "/")
		if len(parts) < 5 {
			return nil // erwartet: .../data/<db>/<table>/indexes/index_*.json
		}
		db := parts[len(parts)-5]
		table := parts[len(parts)-4]
		indexFile := d.Name()
		indexName := strings.TrimSuffix(strings.TrimPrefix(indexFile, "index_"), ".json")
		// Lese Indexdatei
		b, err := os.ReadFile(path)
		if err != nil {
			log.Printf("[IndexRegistry] Fehler beim Lesen von %s: %v", path, err)
			return nil
		}
		var idxObj map[string][]string
		if err := json.Unmarshal(b, &idxObj); err != nil {
			log.Printf("[IndexRegistry] Fehler beim Parsen von %s: %v", path, err)
			return nil
		}
		key := db + "." + table + "." + indexName
		r.mu.Lock()
		r.cache[key] = idxObj
		r.mu.Unlock()
		count++
		return nil
	})
	if err != nil {
		return err
	}
	log.Printf("[IndexRegistry] %d Indexe in %.2fs geladen", count, time.Since(start).Seconds())
	return nil
}

// GetIndex liefert einen geladenen Index aus der Registry
func (r *Registry) GetIndex(db, table, indexName string) (IndexObj, bool) {
	key := db + "." + table + "." + indexName
	r.mu.RLock()
	defer r.mu.RUnlock()
	idx, ok := r.cache[key]
	return idx, ok
}
