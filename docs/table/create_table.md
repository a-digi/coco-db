# Tabelle anlegen (Create Table)

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

## Endpunkt

```
POST /api/databases/{dbname}/tables
```

## Beispiel-Request

```bash
curl -X POST http://localhost:2022/api/databases/testdb/tables \
  -H "Content-Type: application/json" \
  -d '{
    "tableName": "users",
    "fields": [
      { "name": "id", "type": "string", "required": true },
      { "name": "email", "type": "string", "required": true }
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
      { "name": "email", "type": "string", "required": true }
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

## Beispiel-Response (Fehler: Tabelle existiert bereits)

```json
{
  "success": false,
  "error": {
    "code": "ERR_TABLE_EXISTS",
    "message": "Tabelle existiert bereits"
  },
  "httpCode": 409,
  "executionTime": "..."
}
```

## Hinweise
- Die Metadaten der Tabelle werden in `{datadir}/{db}/tables.json` gespeichert.
- Die Felder `httpCode` und `executionTime` sind immer enthalten.
- Fehler werden als JSON mit success: false und error-Objekt zurückgegeben.
