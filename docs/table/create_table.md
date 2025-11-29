# Tabelle anlegen (Create Table)

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
