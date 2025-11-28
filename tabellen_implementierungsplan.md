**Detaillierter Implementierungsplan für Punkt 4: Tabellen (basierend auf PROJECT_REQUIREMENT und Roadmap)**

---

### 1. API-Design & Endpunkte
- [x] Definiere REST-API-Endpunkte für Tabellenoperationen:
    - [x] POST   `/api/databases/{dbname}/tables`         → Tabelle anlegen
    - [ ] DELETE `/api/databases/{dbname}/tables/{tname}` → Tabelle löschen
    - [ ] GET    `/api/databases/{dbname}/tables`         → Tabellen auflisten
    - [ ] GET    `/api/databases/{dbname}/tables/{tname}` → Metadaten einer Tabelle abrufen
    - [ ] PUT    `/api/databases/{dbname}/tables/{tname}` → Metadaten einer Tabelle bearbeiten (Edit)

### 2. Tabellen-Anlage (Create)
- [x] Validierung des Tabellennamens (Pattern, Länge, keine Duplikate)
- [x] Validierung und ggf. Anlegen des Datenbankverzeichnisses
- [x] Anlegen des Tabellenverzeichnisses und meta.json
- [x] Fehlerbehandlung: Existiert bereits, ungültige Felder, etc.
- [x] Logging aller Operationen

### 3. Verzeichnisstruktur & Dateianlage
- [x] Lege für jede Tabelle ein Verzeichnis an: `/data/{datenbankname}/{tabellenname}/`
- [x] Lege in jedem Tabellenverzeichnis eine `meta.json` an (Tabellen-Metadaten, Felder, Typen, Constraints, schemaVersion, etc.)
- [x] Lege Unterverzeichnisse für Einträge und Indizes an:
    - `/data/{datenbankname}/{tabellenname}/entries/` (für Dokumente)
    - `/data/{datenbankname}/{tabellenname}/indexes/` (für Indexdateien)

### 4. Tabellen-Metadaten (meta.json)
- [x] Definiere das Schema für meta.json:
    - Tabellenname, schemaVersion, Felder (Name, Typ, Constraints), Indexdefinitionen, Optionen (z. B. Dokumenten-Limit)
- [x] Validierung und Defaultwerte für meta.json implementieren
- [x] Versionierung und optionale Migration vorbereiten

### 5. Tabellen-Liste (List)
- [ ] Auflisten aller Tabellen einer Datenbank (Verzeichnisse unterhalb von /data/{datenbankname}/)
- [ ] Optional: Filterung nach gültigen Tabellen (nur mit meta.json)
- [ ] Rückgabe als JSON-Array

### 6. Tabellen-Löschen (Delete)
- [ ] Validierung des Tabellennamens
- [ ] Löschen des Tabellenverzeichnisses (rekursiv)
- [ ] Fehlerbehandlung: Nicht gefunden, gesperrt, etc.
- [ ] Logging

### 7. Tabellen-Metadaten abrufen (Get)
- [ ] meta.json einer Tabelle lesen und als JSON zurückgeben
- [ ] Fehlerbehandlung: Nicht gefunden, ungültig, etc.

### 8. Tests & Validierung
- [x] Unit- und Integrationstests für alle Endpunkte und Fehlerfälle
- [x] Testfälle für Namensvalidierung, Duplikate, fehlerhafte Metadaten, etc.

### 9. Dokumentation
- [ ] API-Dokumentation der Endpunkte und Beispiele
- [ ] Beschreibung des meta.json-Schemas
- [ ] Hinweise zu Fehlercodes und Logging

---

**Jeder Schritt wird nach Umsetzung mit Tests und Review überprüft.**

Wenn du Anpassungen oder Ergänzungen wünschst, gib Bescheid!
