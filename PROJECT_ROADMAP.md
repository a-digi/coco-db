# Projekt-Roadmap: JSON-Dokumentdatenbank

Diese Roadmap orientiert sich an der Struktur und den Begrifflichkeiten der PROJECT_REQUIREMENT.md, um maximale Konsistenz und Nachvollziehbarkeit zu gewährleisten.

## 1. Datenbanken
- [ ] Implementierung der Verwaltung von Datenbanken (Anlegen, Löschen, Auflisten)
- [ ] Anlegen der Verzeichnisstruktur `/data/{datenbankname}/`

## 2. Tabellen
- [ ] Implementierung der Verwaltung von Tabellen innerhalb einer Datenbank (Anlegen, Löschen, Auflisten)
- [ ] Anlegen der Verzeichnisstruktur `/data/{datenbankname}/{tabellenname}/`
- [ ] Speichern und Laden der Metadaten-Datei `meta.json`

## 3. Felder und unterstützte Datentypen
- [ ] Validierung von Einträgen anhand der in `meta.json` definierten Felder und Typen
- [ ] Durchsetzung aller Constraints (minLength, maxLength, nullable, pattern, enum, etc.)

## 4. Index Engine (Tabellenbasiert)
- [ ] Aufbau und Verwaltung der In-Memory-Indizes pro Tabelle
- [ ] Persistenz der Indexdateien (`index.jsonl`, `index_{feldname}.jsonl`)
- [ ] Unterstützung von Primär- und Sekundärindizes inkl. Unique/Sparse

## 5. Filter & Query
- [ ] Implementierung von Abfrage- und Filtermechanismen auf Basis der Indizes und/oder vollständiger Iteration

## 6. Speicherstruktur
- [ ] Speicherung aller Einträge als einzelne JSON-Dateien unter `/data/{datenbankname}/{tabellenname}/entries/{dokumenten_id}.json`
- [ ] Sicherstellung der ACID-Prinzipien für alle Operationen
- [ ] Validierung und Fehlerbehandlung gemäß Spezifikation
- [ ] Umsetzung von Backup- und Recovery-Strategien
- [ ] Logging und Monitoring

## 7. Event Sourcing
- [ ] Implementierung des Event Sourcing: Jede Änderung erzeugt ein Event, das chronologisch und unveränderlich gespeichert wird
- [ ] Mechanismen zum Wiederherstellen des Systemzustands aus Events
- [ ] Integration der Events in Backup und Recovery

## 8. Sekundärindizes
- [ ] Erweiterung der Index Engine um zusätzliche, benutzerdefinierte Sekundärindizes

## 9. Server
- [ ] Entwicklung der API (REST/HTTP) für alle Kernoperationen (CRUD, Index, Events, Transaktionen)
- [ ] Sicherstellung, dass alle Ein- und Ausgaben im JSON-Format erfolgen
- [ ] Implementierung von Authentifizierung und Autorisierung (optional)

## 10. Erweiterte Features & Qualitätssicherung
- [ ] Soft Deletes, Audit Trail, Schema-Versionierung, Migration
- [ ] Dokumenten-Limitierung, Health-Checks, Performance-Optimierungen
- [ ] Umfassende Unit- und Integrationstests
- [ ] Pflege der technischen Dokumentation und API-Referenz

---

**Empfohlene Reihenfolge:**
Starte mit den Punkten 1–3 (Basisdatenstrukturen, Validierung), dann Index Engine und Speicherstruktur, gefolgt von Event Sourcing und Server/API. Jeder Schritt sollte durch Tests und Dokumentation begleitet werden.
