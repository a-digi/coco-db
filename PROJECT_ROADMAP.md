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
- [ ] Erweiterbarkeit (Vorbereitung für spätere Features wie Auth, Event Sourcing, Indexe)
- [ ] Sicherheit (Input-Validierung, Fehlervermeidung, optionale Authentifizierung/Autorisierung)
- [ ] Entwicklung der API (REST/HTTP) für alle Kernoperationen (CRUD, Index, Events, Transaktionen)
- [ ] Sicherstellung, dass alle Ein- und Ausgaben im JSON-Format erfolgen
- [ ] Implementierung von Authentifizierung und Autorisierung (optional)

## 3. Datenbanken
- [ ] Implementierung der Verwaltung von Datenbanken (Anlegen, Löschen, Auflisten)
- [ ] Anlegen der Verzeichnisstruktur `/data/{datenbankname}/`

## 4. Tabellen
- [ ] Implementierung der Verwaltung von Tabellen innerhalb einer Datenbank (Anlegen, Löschen, Auflisten)
- [ ] Anlegen der Verzeichnisstruktur `/data/{datenbankname}/{tabellenname}/`
- [ ] Speichern und Laden der Metadaten-Datei `meta.json`

## 5. Felder und unterstützte Datentypen
- [ ] Validierung von Einträgen anhand der in `meta.json` definierten Felder und Typen
- [ ] Durchsetzung aller Constraints (minLength, maxLength, nullable, pattern, enum, etc.)
- [ ] Unterstützung von Schema-Versionierung und Migration

## 6. Index Engine (Tabellenbasiert)
- [ ] Aufbau und Verwaltung der In-Memory-Indizes pro Tabelle
- [ ] Persistenz der Indexdateien (`index.jsonl`, `index_{feldname}.jsonl`)
- [ ] Unterstützung von Primär- und Sekundärindizes inkl. Unique/Sparse

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
