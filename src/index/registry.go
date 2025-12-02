package index

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type IndexData map[string]interface{}

type IndexRegistry struct {
	mu    sync.RWMutex
	cache map[string]IndexData // key: db.table.index
}

var registry *IndexRegistry

func GetRegistry() *IndexRegistry {
	if registry == nil {
		registry = &IndexRegistry{
			cache: make(map[string]IndexData),
		}
	}
	return registry
}

func (r *IndexRegistry) Set(key string, data IndexData) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cache[key] = data
}

func (r *IndexRegistry) Get(key string) (IndexData, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	data, ok := r.cache[key]
	return data, ok
}

// Keys gibt alle Keys der Registry zurück (thread-safe)
func (r *IndexRegistry) Keys() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	keys := make([]string, 0, len(r.cache))
	for k := range r.cache {
		keys = append(keys, k)
	}
	return keys
}

// Lädt alle index_*.json Dateien rekursiv aus dataDir und misst Zeit/Speicher
func LoadAllIndexes(dataDir string) error {
	reg := GetRegistry()
	count := 0
	err := filepath.Walk(dataDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".json" && len(info.Name()) > 6 && info.Name()[:6] == "index_" {
			content, err := ioutil.ReadFile(path)
			if err != nil {
				return err
			}
			var idx IndexData
			if err := json.Unmarshal(content, &idx); err != nil {
				return err
			}
			key := buildIndexKeyFromPath(path, dataDir)
			reg.Set(key, idx)
			count++
		}
		return nil
	})

	return err
}

func buildIndexKeyFromPath(path, dataDir string) string {
	rel, _ := filepath.Rel(dataDir, path)
	// Annahme: dataDir/db/table/indexes/index_xxx.json
	parts := strings.Split(filepath.ToSlash(rel), "/")

	if len(parts) == 4 && parts[2] == "indexes" && strings.HasPrefix(parts[3], "index_") && strings.HasSuffix(parts[3], ".json") {
		db := parts[0]
		table := parts[1]
		index := parts[3][6 : len(parts[3])-5] // index_xxx.json → xxx
		key := db + "." + table + "." + index

		return key
	}

	return path
}

// Prototyp: Index-Update im RAM nach Datenänderung

// UpdateIndexInMemory aktualisiert einen bestimmten Index-Eintrag im RAM.
// action: "insert", "update", "delete"
func (r *IndexRegistry) UpdateIndexInMemory(db, table, index, key string, value interface{}, action string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	idxKey := db + "." + table + "." + index
	idx, ok := r.cache[idxKey]
	if !ok {
		idx = make(IndexData)
		r.cache[idxKey] = idx
	}
	switch action {
	case "insert":
		// Wert als []string behandeln und anhängen
		if v, ok := idx[key].([]string); ok {
			if newVals, ok := value.([]string); ok {
				idx[key] = append(v, newVals...)
			}
		} else {
			if newVals, ok := value.([]string); ok {
				idx[key] = newVals
			}
		}
	case "update":
		// Wert als []string setzen
		if newVals, ok := value.([]string); ok {
			idx[key] = newVals
		}
	case "delete":
		delete(idx, key)
	}
}

// DebugPrintIndex gibt die ersten n Einträge eines Index im RAM aus
func (r *IndexRegistry) DebugPrintIndex(key string, n int) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if _, ok := r.cache[key]; !ok {
		return
	}
}

func (r *IndexRegistry) CountAllEntries() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	total := 0
	for _, idx := range r.cache {
		total += len(idx)
	}
	return total
}
