# Projekt-Roadmap: JSON-Dokumentdatenbank

Diese Roadmap orientiert sich an der Struktur und den Begrifflichkeiten der PROJECT_REQUIREMENT.md, um maximale Konsistenz und Nachvollziehbarkeit zu gewährleisten.

## 1. Initialisierung & Installation
- [x] Implementierung des Initialisierungskommandos (`./coco-db init --data-dir=...`) zur Anlage der Verzeichnisstruktur und Prüfung der Schreibbarkeit
- [x] Sicherstellen, dass das Datenverzeichnis beim Serverstart als Pflichtparameter übergeben und mit dem letzten Pfad abgeglichen wird
- [x] Dokumentation und Validierung aller Startparameter

## 2. Server
- [x] Anforderungen und Schnittstellen klären (REST/HTTP-API, JSON, Authentifizierung/Autorisierung)
- [x] Architektur und Struktur (Trennung von Server-Logik, Routing, Request-Handlern, Geschäftslogik)
- [x] API-Design (Definition der Endpunkte, Request-/Response-Formate)
- [x] Infrastruktur (Initialisierung, Konfiguration, Logging, Monitoring)
- [x] Testbarkeit und Dokumentation (API-Dokumentation, Unit-/Integrationstests)
- [x] Erweiterbarkeit (Vorbereitung für spätere Features wie Auth, Event Sourcing, Indexe)
- [x] Sicherheit (Input-Validierung, Fehlervermeidung, optionale Authentifizierung/Autorisierung)
- [x] Entwicklung der API (REST/HTTP) für alle Kernoperationen (CRUD, Index, Events, Transaktionen)
- [x] Sicherstellung, dass alle Ein- und Ausgaben im JSON-Format erfolgen
- [x] Implementierung von Authentifizierung und Autorisierung (optional)

## 3. Datenbanken
- [x] Implementierung der Verwaltung von Datenbanken (Anlegen, Löschen, Auflisten)
- [x] Anlegen der Verzeichnisstruktur `/data/{datenbankname}/`

## 4. Tabellen
- [x] Implementierung der Verwaltung von Tabellen innerhalb einer Datenbank (Anlegen, Löschen, Auflisten)
- [x] Anlegen der Verzeichnisstruktur `/data/{datenbankname}/{tabellenname}/`
- [x] Speichern und Laden der Metadaten-Datei `meta.json`

## 5. Felder und unterstützte Datentypen
- [x] Implementiere eine zentrale Validierungsfunktion, die alle Einträge (Dokumente) gegen das meta.json-Schema prüft
- [x] Prüfe, ob alle Pflichtfelder vorhanden sind und keine unbekannten Felder enthalten sind (sofern nicht erlaubt)
- [x] Prüfe, ob die Typen der Felder mit den Vorgaben in meta.json übereinstimmen (string, int, bool, date, json)
- [x] Prüfe, ob Defaultwerte korrekt gesetzt werden, falls ein Feld fehlt und ein Default definiert ist
- [x] minLength, maxLength (für Strings)
- [ ] min, max (für numerische Werte)
- [x] nullable (Feld darf null sein)
- [x] pattern (Regulärer Ausdruck für Strings)
- [x] enum (Werte müssen aus einer vorgegebenen Liste stammen)
- [ ] unique (optional, für spätere Index-Validierung)
- [x] required (Pflichtfeld)
- [x] allowAdditionalFields (ob zusätzliche Felder erlaubt sind)
- [x] Fehlerhafte Einträge werden mit präzisen Fehlercodes und -nachrichten abgelehnt (z. B. ERR_FIELD_MISSING, ERR_TYPE_MISMATCH, ERR_CONSTRAINT_FAILED)
- [x] Fehler werden geloggt und in der API-Response zurückgegeben
- [x] Unit-Tests für alle Feldtypen und Constraints (inkl. Grenzfälle) (für String, Integer, Boolean, Date, JSON)
- [x] Integrationstests für das Einfügen, Aktualisieren und Validieren von Dokumenten
- [x] Tests für fehlerhafte und gültige Einträge, Defaultwerte, optionale Felder, etc.

## 6. Index Engine (Tabellenbasiert) – Schritt-für-Schritt-Plan

- [x] **1. Indexdatenstruktur entwerfen**
    - [x] Definition der Index-Metadatenstruktur (z.B. IndexMeta, IndexType, Felder, unique/sparse)
    - [x] Definition der In-Memory-Indexstruktur (z.B. Map, Tree, etc.)

- [x] **2. Indexdefinition in meta.json ermöglichen**
    - [x] Erweiterung von meta.json um Indexdefinitionen (z.B. Primär-/Sekundärindex, unique, sparse)
    - [x] Validierung der Indexdefinitionen beim Anlegen/Ändern einer Tabelle

- [x] **3. Indexdateien persistieren**
    - [x] Format und Speicherort für Indexdateien festlegen (`index.jsonl`, `index_{feldname}.jsonl`)
    - [x] Routinen zum Laden und Speichern der Indexdateien implementieren

- [ ] **4. Indexaufbau und -aktualisierung**
    - [x] Indexaufbau beim Start (Initialisierung aus bestehenden Einträgen)
    - [x] Indexaktualisierung bei Insert, Update, Delete von Einträgen
    - [x] Konsistenzprüfung zwischen Index und Daten

- [ ] **5. Indexabfragen**
    - [x] API/Methoden für schnelle Suche nach Einträgen über Index (z.B. GetByField, RangeQuery)
    - [x] Unterstützung für Primär- und Sekundärindizes
    - [x] Unterstützung für unique/sparse-Index

- [ ] **6. Fehlerbehandlung und Tests**
    - [x] Fehlerfälle (z.B. Indexverletzung, Inkonsistenz) behandeln
    - [x] Unit- und Integrationstests für Indexaufbau, -aktualisierung und -abfrage

- [ ] **7. Dokumentation**
    - [x] Dokumentation der Index-Engine, Indexdefinitionen und API

---

Jeder Schritt sollte einzeln umgesetzt und getestet werden. Die Reihenfolge ist empfohlen, kann aber je nach Architektur angepasst werden.

## 7. Filter & Query
- [ ] Implementierung von Abfrage- und Filtermechanismen auf Basis der Indizes und/oder vollständiger Iteration
- [ ] Unterstützung rekursiver Joins (maxJoinDepth = 64)
- [ ] Substring/Pattern-Suche (LIKE), optionale Vorbereitung für Tokenizer/Volltextsuche
- [ ] Fehlerbehandlung und Validierung aller Query-Parameter

## 8. Speicherstruktur & ACID
- [ ] Speicherung aller Einträge als einzelne JSON-Dateien unter `/data/{datenbankname}/{tabellenname}/entries/{dokumenten_id}.json`
- [ ] Sicherstellung der ACID-Prinzipien für alle Operationen
- [ ] Validierung und Fehlerbehandlung gemäß Spezifikation
- [ ] Umsetzung von Backup- und Recovery-Strategien
- [ ] Logging und Monitoring

## 9. Event Sourcing
- [ ] Implementierung des Event Sourcing: Jede Änderung erzeugt ein Event, das chronologisch und unveränderlich gespeichert wird
- [ ] Mechanismen zum Wiederherstellen des Systemzustands aus Events
- [ ] Integration der Events in Backup und Recovery
- [ ] Audit Trail, Undo/Redo, Zeitreisen

## 10. Sekundärindizes
- [ ] Erweiterung der Index Engine um zusätzliche, benutzerdefinierte Sekundärindizes

## 11. Erweiterte Features & Qualitätssicherung
- [ ] Soft Deletes, Audit Trail, Schema-Versionierung, Migration
- [ ] Dokumenten-Limitierung, Health-Checks, Performance-Optimierungen
- [ ] Umfassende Unit- und Integrationstests
- [ ] Pflege der technischen Dokumentation und API-Referenz
- [ ] Optionale Erweiterung: Volltextsuche mit Tokenizer und Inverted Index

---

**Empfohlene Reihenfolge:**
Starte mit Initialisierung & Installation, dann Server, Datenbanken, Tabellen, Felder, Index Engine, Speicherstruktur & ACID, gefolgt von Filter & Query, Event Sourcing und Server/API. Jeder Schritt sollte durch Tests und Dokumentation begleitet werden.

---

## ToDo: Schema-Versionierung und Migration

- [ ] meta.json enthält ein Feld `schemaVersion` (Pflichtfeld für jede Tabelle)
- [ ] Bei Änderung von `schemaVersion`: Migration bestehender Einträge auf das neue Schema ermöglichen
- [ ] Optional: Migrationsskripte oder automatische Anpassung der Einträge implementieren

Diese Punkte sind für die Zukunft vorgesehen und noch nicht umgesetzt. Sie sind essenziell für Wartbarkeit und Weiterentwicklung des Systems.
