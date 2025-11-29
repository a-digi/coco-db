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
    - [x] POST /api/v1/db/{dbName}/tables/{tableName}/entries
    - [x] Erwartet ein JSON-Objekt im Request-Body

2. **Router-Handler implementieren**
    - [x] Lese dbName und tableName aus der URL
    - [x] Lese und parse den Request-Body in ein map[string]interface{} (JSON-Decoding findet im Routing statt)
    - [x] Rufe InsertEntry(dbName, tableName, entry) auf

3. **Schema laden**
    - [x] Lade das TableMeta-Schema für dbName und tableName

4. **Validierung**
    - [x] Rufe ValidateEntry(entry, schema) auf
    - [x] Falls Fehler: Gib Fehlerdetails zurück (z. B. als Fehlerobjekt)

5. **Persistierung**
    - [ ] Wenn keine Fehler: Speichere den Eintrag (z. B. in Datei, DB, etc.)
    - [ ] Gib Erfolg zurück (z. B. true, ID, etc.)

6. **Fehlerbehandlung**
    - [ ] Bei internen Fehlern: Gib Fehlerobjekt zurück

#### Zusätzliche Validierung: Verbotenes ID-Feld
- Beim Anlegen eines Eintrags darf das Feld `ID` (Groß- oder Kleinschreibung, also `ID` oder `id`) **nicht** im Request-Body enthalten sein.
- Wird das Feld `ID` oder `id` im Eintrag gefunden, wird der Request mit einem Validierungsfehler (z. B. Fehlercode `ERR_FORBIDDEN_ID_FIELD`) abgelehnt.
- Die ID wird ausschließlich vom System generiert und dem Eintrag zugewiesen.

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

---

## Persistierung und Versionierung von Einträgen

- **Ablagepfad für Einträge:**
  - Jeder Eintrag wird unter folgendem Pfad gespeichert:
    `<DataDir>/<dbName>/<tableName>/entries/<entryId>/{entryId}.json`
  - `<entryId>` ist ein eindeutiger Identifier für den Eintrag (z. B. UUID oder generierter Wert).

- **Versionierung:**
  - Für jeden Eintrag gibt es im Ordner `<DataDir>/<dbName>/<tableName>/entries/<entryId>/` eine Datei `version.json`.
  - Die Datei `version.json` enthält ein Array von Versionseinträgen.
  - Ein Versionseintrag besteht aus:
    - `entryId`: ID des Eintrags
    - `created_at`: Zeitstempel der Erstellung
    - `versioned_at`: Zeitstempel, wann diese Version durch eine neue ersetzt wurde (für die aktuelle Version leer)
    - `version_number`: fortlaufende Nummer (beginnend bei 1, wird bei jeder neuen Version erhöht)
  - **Wird ein aktiver Eintrag versioniert, wird die bisherige Datei `{entryId}.json` nach `<DataDir>/<dbName>/<tableName>/entries/<entryId>/version/{versionNumber}.json` verschoben.**

- **Ablauf beim Hinzufügen eines Eintrags:**
  1. Existiert der Ordner `<entryId>` nicht, wird er angelegt und eine neue `version.json` mit dem ersten Versionseintrag erstellt.
  2. Existiert bereits ein aktiver Eintrag (d. h. ein Eintrag ohne `versioned_at`), wird dieser vor dem Hinzufügen des neuen Eintrags versioniert:
     - Das Feld `versioned_at` des bisherigen Eintrags wird mit dem aktuellen Zeitstempel gefüllt.
     - Der neue Eintrag erhält die nächste `version_number` und ein leeres `versioned_at`.
  3. Die Versionseinträge in `version.json` sind chronologisch sortiert und die `version_number` ist immer eindeutig und aufsteigend.

- **Ablauf beim Editieren eines Eintrags:**
  1. Lade den aktuellen Eintrag (`{entryId}.json`) und die zugehörige `version.json`.
  2. Validierung des neuen Eintrags gegen das aktuelle Schema.
  3. Versioniere den bisherigen Stand:
     - Setze das Feld `versioned_at` des bisherigen Eintrags auf den aktuellen Zeitstempel.
     - Verschiebe die bisherige Datei `{entryId}.json` nach `version/{versionNumber}.json`.
     - Erhöhe die `version_number`.
  4. Speichere den neuen Eintrag als `{entryId}.json` und trage ihn als neue Version in `version.json` ein (mit leerem `versioned_at`).

- **Ablauf beim Löschen eines Eintrags:**
  1. Lade die `version.json` und prüfe, ob ein aktiver Eintrag existiert.
  2. Setze das Feld `versioned_at` des aktiven Eintrags auf den aktuellen Zeitstempel.
  3. Verschiebe die Datei `{entryId}.json` nach `version/{versionNumber}.json` (optional: setze ein Lösch-Flag oder entferne die Datei nicht endgültig, um Historie zu bewahren).
  4. Es existiert kein aktiver Eintrag mehr für diese `entryId`.

- **Ablauf beim Lesen eines Eintrags:**
  1. Lade die Datei `{entryId}.json` für den aktuellen Stand.
  2. Optional: Lade aus dem Unterordner `version/` eine bestimmte Version (`{versionNumber}.json`) für historische Abfragen.
  3. Die `version.json` gibt Auskunft über alle Versionen und deren Zeitstempel.

- **Beispiel für version.json:**

```json
[
  {
    "entryId": "abc123",
    "created_at": "2025-11-29T12:00:00Z",
    "versioned_at": "2025-11-29T13:00:00Z",
    "version_number": 1
  },
  {
    "entryId": "abc123",
    "created_at": "2025-11-29T13:00:00Z",
    "versioned_at": "",
    "version_number": 2
  }
]
```

- **Hinweis:**
  - Die Versionierung wird in einer eigenen Datei `src/table/entries/version.go` implementiert.
  - Die Versionierung ist Pflicht, bevor ein neuer Eintrag als aktiv gespeichert wird.
  - Editieren, Löschen und Lesen greifen auf die gleiche Versionierungslogik und -struktur zurück.
