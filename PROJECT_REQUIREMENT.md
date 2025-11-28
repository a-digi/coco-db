# Projektanforderungen: JSON-Dokumentdatenbank

## 1. Datenbanken
- Es können beliebig viele Datenbanken existieren. Jede Datenbank ist eine eigenständige logische Einheit und wird durch einen eindeutigen Namen oder eine ID identifiziert.
- Jede Datenbank erhält einen eigenen Ordner: `/data/{datenbankname}/`
- Die Verwaltung von Datenbanken (Anlegen, Löschen, Auflisten) erfolgt über die Server-API.

## 2. Tabellen
- Innerhalb jeder Datenbank können beliebig viele Tabellen existieren. Tabellen dienen der logischen Gruppierung von Dokumenten.
- Für jede Tabelle gibt es einen eigenen Unterordner: `/data/{datenbankname}/{tabellenname}/`
- Jede Tabelle besitzt einen eigenen Index und einen eigenen Speicherort für Dokumente.
- Die Dokumente einer Tabelle werden als einzelne JSON-Dateien gespeichert: `/data/{datenbankname}/{tabellenname}/entries/{dokumenten_id}.json`
- Jede Tabelle besitzt eine eigene Metadaten-Datei im JSON-Format: `/data/{datenbankname}/{tabellenname}/meta.json`

## 3. Felder und unterstützte Datentypen
- Unterstützte Datentypen für Felder:
  - string (Text)
  - number (Ganzzahlen und Fließkommazahlen)
  - boolean (true/false)
  - json (JSON-Objekte)
  - date (ISO 8601, als string gespeichert)
- Felder werden in der Metadaten-Datei (`meta.json`) als Array von Objekten mit Namen, Typ und optionalen Eigenschaften (z. B. required, default, description) beschrieben.
- Optional können weitere Eigenschaften oder Einschränkungen für Werte definiert werden, z. B.:
  - `minLength`, `maxLength` (für string, number und json)
  - `nullable` (für alle Typen, erlaubt null-Werte)
  - `pattern` (für string, Regex)
  - `enum` (für string, number, boolean)

**Beispiel für Felder-Definitionen mit Datentypen und optionalen Eigenschaften:**

- string:
  ```json
  { "name": "username", "type": "string", "required": true, "minLength": 3, "maxLength": 20, "nullable": false }
  ```
- number:
  ```json
  { "name": "age", "type": "number", "minLength": 1, "maxLength": 3, "nullable": true }
  ```
- boolean:
  ```json
  { "name": "isActive", "type": "boolean", "default": false, "nullable": false }
  ```
- json:
  ```json
  { "name": "address", "type": "json", "maxLength": 5, "nullable": true }
  ```
- date:
  ```json
  { "name": "created_at", "type": "date", "nullable": false }
  ```

## 4. Index Engine (Tabellenbasiert)
Die Indexsuche erfolgt im Speicher (In-Memory) für schnelle Zugriffe. Zusätzlich wird der Index jeder Tabelle regelmäßig in einer Datei persistiert, um nach einem Neustart wiederhergestellt werden zu können.

**Index-Dateiformat und Anzahl der Indexdateien:**
- Pro Tabelle existiert mindestens eine Indexdatei für den Primärschlüssel: `/data/{datenbankname}/{tabellenname}/index.jsonl`.
- Für jeden definierten Sekundärindex (z. B. auf ein bestimmtes Feld) wird eine zusätzliche Indexdatei angelegt, z. B.: `/data/{datenbankname}/{tabellenname}/index_{feldname}.jsonl`.
- Die Anzahl der Indexdateien pro Tabelle entspricht also: 1 (Primärindex) + Anzahl der Sekundärindizes.
- Jede Indexdatei ist im `.jsonl`-Format (JSON Lines), wobei jede Zeile einen Indexeintrag als JSON-Objekt repräsentiert.
- Dieses Format ist leicht zu parsen, kompatibel mit vielen Tools und ermöglicht effizientes sequentielles Lesen und Schreiben.
- Der Index ist immer tabellenbasiert, d.h. jede Tabelle verwaltet ihre eigenen Indexdateien unabhängig von anderen Tabellen.

## 5. Filter & Query

**Abfragearten:**
- Unterstützung von einfachen Filtern (z. B. Suche nach Feldwerten wie Gleichheit, Bereich, Teilstring).
- Unterstützung von komplexen Filtern mit logischen Operatoren (AND, OR, NOT).
- Unterstützung von Joins über mehrere Tabellen hinweg.
- Vergleichsoperatoren: =, !=, <, <=, >, >=, IN, NOT IN, LIKE/Pattern.

**Suchverfahren und Textsuche:**
- Für strukturierte Filter (Gleichheit, Bereich, IN, etc.) werden direkte Vergleiche auf den Feldwerten durchgeführt, unterstützt durch Indizes, sofern vorhanden.
- Für LIKE/Pattern-Suchen wird standardmäßig eine Substring-Suche (Teilstring-Match) verwendet.
- Für fortgeschrittene Textsuche (z. B. Suche nach einzelnen Wörtern, Wortstämmen, Phrasen) ist perspektivisch die Einführung eines Tokenizers und eines Inverted Index vorgesehen (ähnlich wie bei PostgreSQL oder Elasticsearch). In der ersten Version erfolgt die Textsuche jedoch ohne Tokenizer.
- Volltextsuche, Stemming und Stopword-Filter sind als optionale Erweiterung geplant und werden in der Dokumentation als zukünftige Features ausgewiesen.
- Die unterstützten Operatoren orientieren sich an gängigen Standards aus SQL- und dokumentenbasierten Datenbanken (z. B. =, !=, <, <=, >, >=, IN, NOT IN, LIKE/Pattern).

**Index-Nutzung:**
- Filter auf indizierte Felder nutzen immer den entsprechenden Index für die Suche.
- Filter auf nicht indizierte Felder führen zu einem vollständigen Scan aller Einträge (Full Table Scan).

**Query-Syntax (API):**
- Abfragen werden als JSON-Objekt an die API übergeben, z. B.:
  ```json
  {
    "table": "kunden",
    "filter": {
      "age": { "gte": 18, "lte": 65 },
      "isActive": true
    },
    "joins": [
      {
        "table": "bestellungen",
        "on": { "kunden.id": "bestellungen.kunden_id" },
        "type": "inner",
        "filter": { "status": "offen" },
        "joins": [
          {
            "table": "produkte",
            "on": { "bestellungen.produkt_id": "produkte.id" },
            "type": "left",
            "filter": { "kategorie": "digital" }
          }
        ]
      }
    ],
    "sort": [{ "field": "created_at", "direction": "desc" }],
    "limit": 20,
    "offset": 0
  }
  ```
- Das Feld `joins` ist rekursiv: Jede Join-Definition kann selbst wieder ein `joins`-Array enthalten, um beliebig tiefe Join-Hierarchien zu ermöglichen.
- Die maximale Tiefe für verschachtelte Joins beträgt 64 Ebenen (`maxJoinDepth = 64`).
- Jede Join-Definition besteht aus:
  - `table`: Name der Zieltabelle
  - `on`: Join-Bedingung (Feld in Haupttabelle → Feld in Zieltabelle)
  - `type`: Join-Typ (z. B. inner, left; optional, Standard: inner)
  - `filter`: Optionaler Filter auf die gejointe Tabelle
  - `joins`: Optional, Array weiterer Joins auf dieser Ebene
- Unterstützung von Sortierung, Limitierung und Pagination.

**Beispiele:**
- Einfache Abfrage: Alle aktiven Kunden über 30 Jahre
  ```json
  {
    "table": "kunden",
    "filter": {
      "isActive": true,
      "age": { "gt": 30 }
    }
  }
  ```
- Komplexe Abfrage mit Join: Kunden mit offenen Bestellungen, sortiert nach Erstellungsdatum
  ```json
  {
    "table": "kunden",
    "joins": [
      {
        "table": "bestellungen",
        "on": { "kunden.id": "bestellungen.kunden_id" },
        "type": "inner",
        "filter": { "status": "offen" }
      }
    ],
    "sort": [{ "field": "created_at", "direction": "desc" }],
    "limit": 20
  }
  ```
- Komplexe Abfrage mit verschachtelten Joins:
  ```json
  {
    "table": "kunden",
    "joins": [
      {
        "table": "bestellungen",
        "on": { "kunden.id": "bestellungen.kunden_id" },
        "type": "inner",
        "filter": { "status": "offen" },
        "joins": [
          {
            "table": "produkte",
            "on": { "bestellungen.produkt_id": "produkte.id" },
            "type": "left",
            "filter": { "kategorie": "digital" }
          }
        ]
      }
    ],
    "sort": [{ "field": "created_at", "direction": "desc" }],
    "limit": 20
  }
  ```

## 6. Speicherstruktur
Die eigentlichen Dateninhalte (Dokumente) werden ausschließlich persistent als einzelne JSON-Dateien gespeichert und nicht im RAM gehalten. Nur die Indexdaten der jeweiligen Tabellen befinden sich im Speicher (In-Memory), um schnelle Suchen und Zugriffe zu ermöglichen.

**JSON als primäres Ein- und Ausgabeformat:**
- JSON ist das zentrale Datenformat für alle Ein- und Ausgaben des Servers.
- Alle Anfragen an die API (z. B. zum Erstellen, Aktualisieren, Abfragen oder Löschen von Dokumenten) sowie alle Antworten des Servers erfolgen im JSON-Format.
- Die Validierung, Speicherung und Übertragung der Daten basiert vollständig auf JSON.

**ACID-Konformität:**
- Das System MUSS die ACID-Prinzipien (Atomicity, Consistency, Isolation, Durability) für alle Schreib- und Leseoperationen gewährleisten. Das bedeutet insbesondere:
  - Schreibvorgänge sind atomar und führen entweder vollständig oder gar nicht statt.
  - Die Datenintegrität bleibt bei allen Operationen erhalten (Konsistenz).
  - Gleichzeitige Zugriffe beeinträchtigen sich nicht gegenseitig (Isolation).
  - Nach Abschluss einer Operation sind die Daten dauerhaft und über Systemabstürze hinweg gespeichert (Durability).

**Validierung und Fehlerbehandlung:**
- Vor dem Speichern oder Aktualisieren eines Dokuments wird eine umfassende Validierung gegen das in `meta.json` definierte Schema durchgeführt.
  - Alle Pflichtfelder (`required: true`) müssen vorhanden und gültig sein.
  - Datentypen müssen exakt mit der Definition übereinstimmen (z. B. darf ein Feld vom Typ `string` keine Zahl enthalten).
  - Einschränkungen wie `minLength`, `maxLength`, `min`, `max`, `pattern`, `enum` und `nullable` werden strikt geprüft.
  - Felder, die nicht im Schema definiert sind, werden standardmäßig abgelehnt, sofern nicht explizit als erlaubt konfiguriert.
  - Standardwerte (`default`) werden gesetzt, falls ein Wert fehlt und das Feld nicht `required` ist.
- Bei Validierungsfehlern wird der Schreibvorgang abgebrochen und eine detaillierte Fehlermeldung mit Angabe aller Verstöße zurückgegeben.
- Fehlercodes und -nachrichten sind konsistent und maschinenlesbar aufgebaut (z. B. `ERR_VALIDATION_FAILED`, `ERR_TYPE_MISMATCH`, `ERR_MISSING_REQUIRED_FIELD`).
- Fehlerhafte oder unvollständige Dokumente werden niemals gespeichert oder überschreiben bestehende Daten nicht.
- Alle Fehlerfälle werden protokolliert, um Nachvollziehbarkeit und Debugging zu gewährleisten.

**Schema-Versionierung:**
- Jede Tabelle enthält in ihrer meta.json ein Feld `schemaVersion` (z. B. `"schemaVersion": 1`).
- Änderungen am Tabellenschema werden durch Erhöhung der Versionsnummer nachvollziehbar gemacht.
- Bei Schema-Änderungen ist eine Migrationsstrategie zu definieren (z. B. automatische Anpassung bestehender Dokumente oder explizite Migration durch den Nutzer).

**Umgang mit unbekannten Feldern:**
- Standardmäßig dürfen Dokumente keine Felder enthalten, die nicht im Schema definiert sind (`allowAdditionalFields: false`).
- Optional kann in meta.json explizit erlaubt werden, zusätzliche Felder zu speichern (`allowAdditionalFields: true`).

**Index-Optionen:**
- Indexdefinitionen in meta.json können optionale Eigenschaften wie `unique` (Eindeutigkeit) und `sparse` (nur für vorhandene Werte) enthalten.
- Beispiel: `{ "name": "email", "type": "secondary", "fields": ["email"], "unique": true, "sparse": false }`

**Soft Deletes und Historie:**
- Optional kann ein Feld wie `deleted: true` für Soft Deletes verwendet werden, sodass gelöschte Dokumente nicht physisch entfernt, sondern als gelöscht markiert werden.
- Für Änderungs-Historie kann eine Audit-Trail-Strategie definiert werden (z. B. Speicherung von Änderungen mit Zeitstempel und Benutzer-ID).

**Backup- und Recovery-Strategie:**
- Es wird empfohlen, regelmäßige Backups der Daten- und Indexdateien zu erstellen.
- Die Wiederherstellung erfolgt durch Rückspielen der Backup-Dateien in die entsprechende Datenbankstruktur.

**API-Fehlercodes und -antworten:**
- Die API gibt standardisierte Fehlercodes und -nachrichten zurück (z. B. 400 für Validierungsfehler, 404 für nicht gefundene Ressourcen, 409 für Konflikte).
- Fehlerantworten enthalten maschinenlesbare Codes und eine menschenlesbare Beschreibung.

**Sicherheit und Zugriffskontrolle:**
- Optional kann Authentifizierung (z. B. API-Keys) und Autorisierung (z. B. Rollen, Rechte auf DB/Table-Ebene) implementiert werden.
- Zugriffsbeschränkungen können in der Server-Konfiguration oder in meta.json definiert werden.

**Performance- und Skalierbarkeitsoptionen:**
- Das System sollte für große Datenmengen und viele gleichzeitige Zugriffe optimiert werden (z. B. durch effiziente Indizes, Caching, ggf. Sharding).
- Hinweise zur Skalierung und zu empfohlenen Hardware-Anforderungen können dokumentiert werden.

**Dokumenten-Limitierung:**
- Optional kann eine maximale Anzahl an Dokumenten pro Tabelle oder Datenbank festgelegt werden.
- Überschreitet eine Tabelle das Limit, werden weitere Schreibvorgänge abgelehnt und ein entsprechender Fehlercode zurückgegeben.

**Logging und Monitoring:**
- Das System protokolliert alle relevanten Ereignisse (z. B. Fehler, Zugriffe, Änderungen) in maschinenlesbaren Logdateien.
- Monitoring-Schnittstellen (z. B. Health-Checks, Metriken) können bereitgestellt werden.

**Transaktionsunterstützung:**
- Das System kann mehrere Operationen in einer Transaktion zusammenfassen, die entweder vollständig oder gar nicht ausgeführt werden (All-or-Nothing-Prinzip).
- Transaktionen werden explizit über die API gestartet und abgeschlossen.

## 7. Event Sourcing
- Das System MUSS Event Sourcing unterstützen, d. h. jede Änderung an einem Dokument (Erstellung, Aktualisierung, Löschung) wird als separates, unveränderliches Event gespeichert.
- Alle Events werden chronologisch und unveränderlich abgelegt, sodass der komplette Zustand einer Tabelle oder eines Dokuments zu jedem beliebigen Zeitpunkt wiederhergestellt werden kann.
- Vorteile von Event Sourcing:
  - Ermöglicht vollständige Nachvollziehbarkeit und Revisionssicherheit aller Änderungen (Audit Trail).
  - Erhöht die Ausfallsicherheit: Nach einem Systemfehler kann der Zustand durch das erneute Abspielen aller Events exakt rekonstruiert werden.
  - Unterstützt fortgeschrittene Features wie Zeitreisen (Abfragen des Systemzustands zu einem beliebigen Zeitpunkt), Undo/Redo und flexible Datenmigrationen.
  - Erleichtert die Integration mit externen Systemen durch Event-Streams.
- Die Events werden persistent gespeichert und sind Teil des Backups und der Recovery-Strategie.

## 8. Sekundärindizes
Für häufig abgefragte Felder innerhalb einer Tabelle werden zusätzliche Indizes gepflegt, um gezielte Suchen zu beschleunigen.

## 9. Server
Ein Server stellt eine API zum Anlegen, Löschen und Verwalten von Datenbanken und Tabellen sowie zum Speichern und Abrufen der JSON-Dokumente bereit.

## 10. Installation & Setup

**Initialisierung der Datenbank:**
- Vor dem ersten Start des Servers sollte das Datenverzeichnis initialisiert werden. Dies kann über ein Terminal-Kommando erfolgen:
  ```sh
  ./coco-db init --data-dir=/pfad/zum/datenverzeichnis
  ```
- Der Befehl legt die notwendige Verzeichnisstruktur an und prüft, ob das Verzeichnis beschreibbar ist.
- Optional kann beim Initialisieren auch direkt eine erste Datenbank angelegt werden (z. B. mit `--database=meinedb`).
- Die Initialisierung ist Voraussetzung, bevor der Server im Produktivbetrieb gestartet werden kann.

**Starten des Servers (Beispiel):**
  ```sh
  ./coco-db --data-dir=/pfad/zum/datenverzeichnis
  ```
  - Das Datenverzeichnis MUSS explizit angegeben werden (siehe Speicherstruktur).

**Hinweise:**
- Die Konfiguration kann auch über Umgebungsvariablen oder eine Konfigurationsdatei erfolgen (siehe Dokumentation).
- Für den Produktivbetrieb wird empfohlen, das Datenverzeichnis auf einem persistenten und gesicherten Laufwerk zu betreiben.
- Weitere Startparameter und Konfigurationsoptionen werden in der technischen Dokumentation beschrieben.

---

**Hinweis:** Es wird kein Code geschrieben, bevor der Projektplan und die Roadmap final abgestimmt sind.
