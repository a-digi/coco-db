# Index-Engine: Architektur und Funktionsweise

## Überblick
Die Index-Engine ermöglicht schnelle Suchen und Sortierungen in Tabellen durch persistente Indexdateien. Unterstützt werden Primär- und Sekundärindizes, unique/sparse-Optionen und BTree-basierte Strukturen.

## Index-Metadaten (meta.json)
Jede Tabelle kann in ihrer meta.json beliebig viele Indizes definieren:

```json
"indexes": [
  {
    "name": "email_idx",
    "type": "secondary",
    "fields": ["email"],
    "unique": true,
    "sparse": false
  }
]
```
- **name**: Name des Index (Dateiname: `index_<name>.json`)
- **type**: "primary" oder "secondary"
- **fields**: Felder, die indexiert werden (aktuell nur Single-Field)
- **unique**: true = Wert darf nur einmal vorkommen
- **sparse**: true = nur Einträge mit Wert werden indexiert

## Indexdateien
- Speicherort: Im Tabellenverzeichnis als `index_<name>.json`
- Format: JSON-Objekt, Key = Wert, Value = Array von EntryIDs

Beispiel:
```json
{
  "foo@bar.de": ["id1"],
  "bar@baz.de": ["id2"]
}
```

## Indexaufbau und -aktualisierung
- Beim Start: Index wird aus allen Einträgen neu aufgebaut (Konsistenzprüfung möglich)
- Insert/Update/Delete: Indexdateien werden automatisch aktualisiert
- unique/sparse werden geprüft und durchgesetzt

## Fehlerfälle
- unique-Verletzung: Insert/Update wird abgelehnt
- Inkonsistenz: Kann per Konsistenzprüfung erkannt und gemeldet werden
- Fehlende Indexdatei: Wird beim Zugriff automatisch neu erstellt

## Beispiel: Indexabfrage
API:
```
GET /api/databases/{dbname}/tables/{tablename}/index/{indexname}/query?value=foo@bar.de
```
Antwort:
```json
{
  "success": true,
  "data": { "entryIds": ["id1"] },
  "httpCode": 200
}
```

## Erweiterung
- Unterstützung für Multi-Field-Indizes, Range-Queries, weitere Index-Typen (z.B. Hash, LSMTree) ist vorbereitet.

---

Weitere Details siehe Quellcode und API-Dokumentation.
