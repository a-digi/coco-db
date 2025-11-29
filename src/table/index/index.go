package index

// IndexType beschreibt den Typ des Index (z.B. BTree, Hash, LSMTree)
type IndexType string

const (
	IndexTypeBTree IndexType = "btree"
)

// IndexMeta beschreibt die Metadaten eines Index
// (z.B. für meta.json und interne Verwaltung)
type IndexMeta struct {
	Name      string     `json:"name"`           // Name des Index
	Fields    []string   `json:"fields"`         // Felder, die indexiert werden
	Type      IndexType  `json:"type"`           // Index-Typ (nur "btree" wird unterstützt)
	Unique    bool       `json:"unique"`         // true = unique-Index
	Sparse    bool       `json:"sparse"`         // true = nur Einträge mit Wert werden indexiert
	Primary   bool       `json:"primary"`        // true = Primärindex
}

// Index ist das Interface für alle Index-Typen (BTree, Hash, ...)
type Index interface {
	Insert(key interface{}, entryId string) error
	Delete(key interface{}, entryId string) error
	Update(oldKey, newKey interface{}, entryId string) error
	Find(key interface{}) ([]string, error) // Liefert alle EntryIDs zum Key
	Range(start, end interface{}) ([]string, error) // Für Bereichsanfragen (nur BTree)
}

// Hinweis: Aktuell wird nur BTree als Index-Typ unterstützt.
// Weitere Implementierungen (BTree) folgen in eigenen Dateien.
