# Index-API: Endpunkte und Beispiele

## Indexabfrage

**GET** `/api/databases/{dbname}/tables/{tablename}/index/{indexname}/query?value=...`

- Liefert alle EntryIDs, die im Index für den Wert gespeichert sind.
- Antwort:
```json
{
  "success": true,
  "data": { "entryIds": ["id1", "id2"] },
  "httpCode": 200
}
```
- Fehlerfall (Index nicht gefunden):
```json
{
  "success": true,
  "data": { "entryIds": [] },
  "httpCode": 200
}
```

## Fehlerfälle
- **Indexverletzung (unique):**
  - Beim Insert/Update wird ein Fehler zurückgegeben, wenn der Wert bereits existiert und unique=true ist.
- **Inkonsistenz:**
  - Kann per Konsistenzprüfung erkannt werden (siehe engine.md).
- **Fehlende Indexdatei:**
  - Wird automatisch neu erstellt oder als leer behandelt.

## Beispiel für Indexdefinition in meta.json
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

## Hinweise
- Die API ist erweiterbar für Range-Queries und Multi-Field-Indizes.
- Siehe auch: [engine.md](./engine.md) für technische Details.
