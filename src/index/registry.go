package index

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
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

// Lädt alle index_*.json Dateien rekursiv aus dataDir und misst Zeit/Speicher
func LoadAllIndexes(dataDir string) error {
	reg := GetRegistry()
	start := time.Now()
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
	dur := time.Since(start)
	log.Printf("[IndexRegistry] %d Indexe in %.2fs geladen.", count, dur.Seconds())
	return err
}

func buildIndexKeyFromPath(path, dataDir string) string {
	rel, _ := filepath.Rel(dataDir, path)
	// Annahme: dataDir/db/table/index_*.json
	partsFS := filepath.ToSlash(rel)
	partsArr := make([]string, 0)
	for _, p := range strings.Split(partsFS, "/") {
		if p != "" {
			partsArr = append(partsArr, p)
		}
	}
	if len(partsArr) >= 3 {
		db := partsArr[0]
		table := partsArr[1]
		index := partsArr[2][6 : len(partsArr[2])-5] // index_xxx.json → xxx
		return db + "." + table + "." + index
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
	case "insert", "update":
		idx[key] = value
	case "delete":
		delete(idx, key)
	}
}
