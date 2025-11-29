# 5. Felder und unterstützte Datentypen

## Ziel
Dieser Plan beschreibt die Umsetzung von Punkt 5 der Roadmap: Zentrale Validierung von Einträgen (Dokumenten) anhand der in `meta.json` definierten Felder, Datentypen und Constraints. Er dient als technische Spezifikation und To-Do-Liste für die Implementierung.

---

## Implementierungsplan für src/table/fields

### Aufgabenübersicht mit Checklisten

1. [x] **1. Definition der Datenstrukturen**
    - [x] Erweiterung/Definition von `FieldMeta` und ggf. `TableMeta` für alle Typen und Constraints
    - [x] Dokumentation der unterstützten Typen und Constraints

2. [x] **2. Zentrale Validierungsfunktion**
    - [x] Funktion `ValidateEntry(entry map[string]interface{}, meta TableMeta) error`
    - [x] Prüfung auf Pflichtfelder (`required`)
    - [x] Typprüfung (`string`, `integer`, `boolean`, `json`, `date`)
    - [x] Prüfung aller Constraints (`minLength`, `maxLength`, `pattern`, `enum`, `nullable`, `default`, `description`)
    - [x] Prüfung auf zusätzliche Felder (`allowAdditionalFields`)
    - [x] Fehlercodes/-nachrichten

3. [x] **3. Einzelne Typ- und Constraint-Validatoren**
    - [x] String-Validator (inkl. minLength, maxLength, pattern, enum)
    - [x] Integer-Validator (inkl. minLength, maxLength, enum)
    - [x] Boolean-Validator (inkl. enum)
    - [x] JSON-Validator (inkl. maxLength, minLength)
    - [x] Date-Validator (ISO 8601)
    - [x] Nullable- und Default-Handling

4. [ ] **4. Fehlerbehandlung**
    - [x] Präzise Fehlercodes/-nachrichten

6. [ ] **6. Tests**
    - [x] Unit-Tests für alle Typen und Constraints (inkl. Grenzfälle) (nur teilweise, String ist abgedeckt, andere Typen sollten ergänzt werden)
    - [ ] Integrationstests für Einfügen/Aktualisieren/Validieren
    - [x] Tests für Fehlerfälle, Defaultwerte, optionale Felder

7. [x] **7. Dokumentation und Beispiele**
    - [x] Dokumentation der Validierungslogik und Beispiele für meta.json

---

## Plan für die konkrete Einfüge-Logik und REST-API-Anbindung

### Ziel
Definiere, wie neue Einträge über eine REST-API entgegengenommen, validiert und gespeichert werden.

### Schritt-für-Schritt-Checkliste

1. **API-Endpunkt definieren**
    - [ ] POST /api/v1/db/{dbName}/tables/{tableName}/entries
    - [ ] Erwartet ein JSON-Objekt im Request-Body

2. **Router-Handler implementieren**
    - [ ] Lese dbName und tableName aus der URL
    - [ ] Lese und parse den Request-Body in ein map[string]interface{} (JSON-Decoding findet im Routing statt)
    - [ ] Rufe InsertEntry(dbName, tableName, entry) auf

3. **Schema laden**
    - [ ] Lade das TableMeta-Schema für dbName und tableName

4. **Validierung**
    - [ ] Rufe ValidateEntry(entry, schema) auf
    - [ ] Falls Fehler: Gib Fehlerdetails zurück (z. B. als Fehlerobjekt)

5. **Persistierung**
    - [ ] Wenn keine Fehler: Speichere den Eintrag (z. B. in Datei, DB, etc.)
    - [ ] Gib Erfolg zurück (z. B. true, ID, etc.)

6. **Fehlerbehandlung**
    - [ ] Bei internen Fehlern: Gib Fehlerobjekt zurück

### Beispiel: Pseudocode für die Einfügefunktion

```go
func InsertEntry(dbName string, tableName string, entry map[string]interface{}) error {
    schema := LoadTableMeta(dbName, tableName)
    errors := ValidateEntry(entry, schema)
    if len(errors) > 0 {
        return errors // oder Fehlerstruktur
    }
    // entry enthält jetzt Defaultwerte
    err := PersistEntry(dbName, tableName, entry)
    if err != nil {
        // Logging und Fehlerbehandlung
        return err
    }
    return nil
}
```

### Hinweise
- Das JSON-Decoding findet im Routing statt, nicht in InsertEntry.
- Die Validierungslogik bleibt unverändert und wird vor der Persistierung aufgerufen.
- Defaultwerte werden durch ValidateEntry gesetzt.
- Fehlerhafte Einträge werden abgelehnt und mit Fehlerdetails beantwortet.
- Die Persistierung ist austauschbar (Datei, DB, etc.).
