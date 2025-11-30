package index

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// BTreeIndex ist eine einfache In-Memory-Baumstruktur für Indexzwecke.
// Für Produktivbetrieb sollte eine bewährte B+Tree-Bibliothek verwendet werden.
type BTreeIndex struct {
	meta   IndexMeta
	root   *btreeNode
	stub   map[string][]string // Nur für Persistenz-Test, bis echte BTree-Logik steht
}

type btreeNode struct {
	keys     []interface{}
	entryIDs [][]string // mehrere EntryIDs pro Key möglich (für nicht-unique)
	children []*btreeNode
	leaf     bool
}

func NewBTreeIndex(meta IndexMeta) *BTreeIndex {
	return &BTreeIndex{
		meta: meta,
		root: &btreeNode{leaf: true},
		stub: map[string][]string{},
	}
}

func (b *BTreeIndex) Insert(key interface{}, entryId string) error {
	// Für den Test: Key als String in stub-Map eintragen
	k, ok := key.(string)
	if !ok {
		return nil // Nur String-Keys für MVP
	}
	if b.stub == nil {
		b.stub = map[string][]string{}
	}
	b.stub[k] = append(b.stub[k], entryId)
	return nil
}

func (b *BTreeIndex) Delete(key interface{}, entryId string) error {
	// TODO: B+Tree-Delete-Logik implementieren
	return nil
}

func (b *BTreeIndex) Update(oldKey, newKey interface{}, entryId string) error {
	// TODO: B+Tree-Update-Logik implementieren
	return nil
}

func (b *BTreeIndex) Find(key interface{}) ([]string, error) {
	// TODO: B+Tree-Find-Logik implementieren
	return nil, nil
}

func (b *BTreeIndex) Range(start, end interface{}) ([]string, error) {
	// TODO: B+Tree-Range-Logik implementieren
	return nil, nil
}

// SaveToFile speichert den Index als JSON-Objekt in eine Datei (z.B. index_{name}.json)
func (b *BTreeIndex) SaveToFile(dir string) error {
	indexPath := filepath.Join(dir, "indexes", "index_"+b.meta.Name+".json")
	data := b.toMap()
	f, err := os.Create(indexPath)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(data)
}

// LoadFromFile lädt den Index aus einer Datei (JSON-Objekt)
func (b *BTreeIndex) LoadFromFile(dir string) error {
	indexPath := filepath.Join(dir, "index_"+b.meta.Name+".json")
	f, err := os.Open(indexPath)
	if err != nil {
		return err
	}
	defer f.Close()
	var data map[string][]string
	dec := json.NewDecoder(f)
	if err := dec.Decode(&data); err != nil {
		return err
	}
	b.fromMap(data)
	return nil
}

// toMap serialisiert den Index als map[string][]string (Stub, da BTree-Logik fehlt)
func (b *BTreeIndex) toMap() map[string][]string {
	return b.stub
}

// fromMap lädt die Daten in den Index (Stub, da BTree-Logik fehlt)
func (b *BTreeIndex) fromMap(m map[string][]string) {
	b.stub = m
}
