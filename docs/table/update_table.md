# Tabelle aktualisieren (Update Table)

## Hinweise zu unterstützten Datentypen für Felder

- Die folgenden Datentypen sind laut Projektanforderung erlaubt:
  - **string** (Text)
  - **number** (Ganzzahlen und Fließkommazahlen)
  - **boolean** (true/false)
  - **json** (JSON-Objekte)
  - **date** (ISO 8601, als string gespeichert)

### Beispiel für Felder mit allen erlaubten Datentypen

```json
[
  { "name": "username", "type": "string", "required": true, "minLength": 3, "maxLength": 20, "nullable": false },
  { "name": "age", "type": "number", "minLength": 1, "maxLength": 3, "nullable": true },
  { "name": "isActive", "type": "boolean", "default": false, "nullable": false },
  { "name": "address", "type": "json", "maxLength": 5, "nullable": true },
  { "name": "created_at", "type": "date", "default": "now", "nullable": false }
]
```

# Tabelle aktualisieren (Update Table)

## Endpunkt

```
PUT /api/databases/{dbname}/tables/{tablename}
```

## Beispiel-Request

```bash
curl -X PUT http://localhost:2022/api/databases/testdb/tables/users \
  -H "Content-Type: application/json" \
  -d '{
    "tableName": "users",
    "fields": [
      { "name": "id", "type": "string", "required": true },
      { "name": "email", "type": "string", "required": true },
      { "name": "active", "type": "bool", "required": false }
    ],
    "indexes": [
      { "name": "primary", "type": "primary", "fields": ["id"], "unique": true }
    ],
    "options": { "documentLimit": 1000 }
  }'
```

## Beispiel-Response (Erfolg)

```json
{
  "success": true,
  "data": {
    "tableName": "users",
    "schemaVersion": 1,
    "fields": [
      { "name": "id", "type": "string", "required": true },
      { "name": "email", "type": "string", "required": true },
      { "name": "active", "type": "bool", "required": false }
    ],
    "indexes": [
      { "name": "primary", "type": "primary", "fields": ["id"], "unique": true }
    ],
    "options": { "documentLimit": 1000 }
  },
  "httpCode": 200,
  "executionTime": "..."
}
```

## Beispiel-Response (Fehler: Tabelle nicht gefunden)

```json
{
  "success": false,
  "error": {
    "code": "ERR_TABLE_NOT_FOUND",
    "message": "Tabelle nicht gefunden"
  },
  "httpCode": 404,
  "executionTime": "..."
}
```

## Hinweise
- Der Datenbankname und Tabellenname werden aus der URL extrahiert.
- Die neuen Metadaten werden im JSON-Body übergeben und ersetzen die bestehende meta.json der Tabelle.
- Fehler werden als JSON mit success: false und error-Objekt zurückgegeben.
- Die Felder `httpCode` und `executionTime` sind immer enthalten.

# Tabelle aktualisieren (Update Table)

## Indexe hinzufügen und entfernen

- Über das Feld `indexes` im Request-Body können Indexe hinzugefügt oder entfernt werden.
- Die übergebenen Indexe ersetzen die bisherigen Indexe vollständig.
- Um einen Index zu entfernen, lasse ihn einfach im neuen Array weg.
- Um einen Index hinzuzufügen, ergänze ihn im Array.
- Die Änderung wird in der meta.json gespeichert und die Indexdateien werden entsprechend angelegt oder gelöscht.

**Beispiel: Index hinzufügen**

```json
"indexes": [
  { "name": "primary", "type": "primary", "fields": ["id"], "unique": true },
  { "name": "email", "type": "secondary", "fields": ["email"], "unique": true }
]
```

**Beispiel: Index entfernen**

```json
"indexes": [
  { "name": "primary", "type": "primary", "fields": ["id"], "unique": true }
]
```

## Hinweise zu Default-Werten für Datumsfelder (date)

- Für Felder vom Typ `date` können folgende Default-Werte verwendet werden:
  - Fester ISO-8601-String, z.B. `"2025-11-30T00:00:00Z"`
  - Platzhalter:
    - `"now"`   → aktuelles Datum und Uhrzeit (UTC, RFC3339)
    - `"today"` → aktuelles Tagesdatum (UTC, 00:00:00 Uhr)
- Beispiel für ein Feld mit Default-Wert:

```json
{
  "name": "updated_at",
  "type": "date",
  "default": "now"
}
```
