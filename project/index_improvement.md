# Index-Verbesserungsvorschläge für hohe Query-Performance

## Empfehlungen

1 **Indexe beim Serverstart in den RAM laden**
  - [x] Lade alle relevanten Indexdateien (z.B. BTree, JSON) beim Start des Servers in einen zentralen In-Memory-Cache.
    - [x] Definiere eine zentrale Index-Registry/Cache-Struktur (z.B. als Go-Map oder Singleton).
    - [x] Implementiere eine Funktion `LoadAllIndexes(dataDir string)`, die rekursiv alle Indexdateien (z.B. `index_*.json`) in allen Tabellenverzeichnissen findet und lädt.
    - [x] Für jeden gefundenen Index: Lese die Datei, parse sie (z.B. als Map oder BTree-Objekt) und speichere sie in der Registry unter dem Schlüssel (Datenbank, Tabelle, Indexname).
    - [x] Stelle sicher, dass die Registry threadsicher ist (z.B. mit sync.RWMutex).
    - [x] Rufe diese Funktion beim Serverstart auf und logge Anzahl, Ladezeit und Speicherverbrauch der geladenen Indexe.
    - [x] Schreibe Unit-Tests für das Laden und die Registry.
    - [x] Dokumentiere die Registry-API für Query-Handler und Filter-Engine.
    - [ ] (Optional) Implementiere ein Interface für verschiedene Index-Typen (BTree, Hash, ...).
    - [ ] (Optional) Füge einen Health-Check hinzu, der prüft, ob alle Indexe geladen wurden.
  - [ ] Halte die Indexe während der Laufzeit aktuell (z.B. mit Event Sourcing/Index-Worker).

2 **Zentralen Index-Cache/Registry implementieren**
  - [x] Entwickle eine zentrale Komponente (z.B. Singleton oder globales Objekt), die alle geladenen Indexe verwaltet.
  - [x] Alle Query-Handler und Filter-Engines greifen auf diesen Cache zu, statt Indexdateien von der Festplatte zu lesen.

3 **Indexe für schnelle Lookups bei Filter- und Join-Operationen nutzen**
  - [ ] Nutze die In-Memory-Indexe, um Filter und Joins effizient (O(1) oder O(log n)) auszuführen.
  - [ ] Reduziere sequentielle Dateizugriffe auf ein Minimum.

4  **Index-Worker/Consumer für Aktualisierung**
  - [ ] Bei jeder Änderung an den Einträgen (Insert, Update, Delete) wird ein Event ins Event-Log geschrieben.
  - [ ] Ein Index-Worker verarbeitet diese Events und aktualisiert die In-Memory-Indexe sowie die persistierten Indexdateien.

5 **Startup-Check und Rebuild**
  - [ ] Beim Serverstart prüft der Index-Worker, ob das Event-Log vollständig verarbeitet ist und holt ggf. fehlende Index-Updates nach.
  - [ ] Optional: Biete einen Rebuild-Befehl an, um alle Indexe aus dem Event-Log neu zu erstellen.

6 **Monitoring und Fehlerbehandlung**
  - [ ] Überwache die Aktualität und Konsistenz der Indexe.
  - [ ] Logge Fehler und Rückstände im Index-Worker.

## 1.3 Architektur- und Schnittstellendesign für den zentralen Index-Cache (Implementierung)

Die zentrale Komponente ist `IndexRegistry` (Singleton, threadsicher), die alle geladenen Indexe im RAM hält. Sie bietet folgende Schnittstellen:

- `GetRegistry() *IndexRegistry` – liefert die Singleton-Instanz
- `Set(key string, data IndexData)` – speichert/aktualisiert einen Index
- `Get(key string) (IndexData, bool)` – liest einen Index
- `LoadAllIndexes(dataDir string) error` – lädt alle Indexdateien rekursiv in die Registry

**Implementierung:**
- Siehe `src/index/registry.go` und `src/index/registry_test.go` für vollständigen Go-Code und Unit-Tests.
- Die Registry ist generisch und kann für verschiedene Index-Typen erweitert werden.
- Die Nutzung erfolgt über den Key `db.table.index` (z.B. `poseidon.users.email`).
- Die Registry kann beim Serverstart initialisiert und für Query- und Filter-Operationen verwendet werden.

**Beispiel für die Nutzung im Serverstart:**

```go
import "src/index"

func main() {
    err := index.LoadAllIndexes("/data")
    if err != nil {
        log.Fatalf("Index-Laden fehlgeschlagen: %v", err)
    }
    // ... Restlicher Serverstart
}
```

## 1.4 Prototyp für das Laden und Aktualisieren der Indexe im RAM

Die IndexRegistry kann jetzt nicht nur beim Serverstart geladen, sondern auch zur Laufzeit aktualisiert werden:

**Neue Methode:**
```go
// Prototyp: Index-Update im RAM nach Datenänderung
// action: "insert", "update", "delete"
func (r *IndexRegistry) UpdateIndexInMemory(db, table, index, key string, value interface{}, action string)
```
- Fügt einen Eintrag hinzu, aktualisiert oder entfernt ihn im RAM-Index.
- Wird nach Insert/Update/Delete aufgerufen, um die Indexe aktuell zu halten.

**Beispiel für die Nutzung nach einem Insert:**
```go
reg := index.GetRegistry()
reg.UpdateIndexInMemory("poseidon", "users", "email", "foo@bar.de", entryID, "insert")
```

**Empfohlene Integration:**
- Nach jeder Datenänderung (Insert/Update/Delete) in der Datenbank wird diese Methode aufgerufen.
- Ein Event-Worker kann das Event-Log konsumieren und so die Indexe im RAM synchron halten.

## Vorteile
- [ ] Deutlich schnellere Filter- und Join-Operationen (keine Full Table Scans)
- [ ] Geringere Latenz und höherer Durchsatz bei komplexen Queries
- [ ] Bessere Skalierbarkeit und Robustheit

---

**Empfohlene nächste Schritte:**
- [ ] Prototyp für das Laden und Aktualisieren der Indexe im RAM
- [ ] Integration in Query- und Filter-Engine
