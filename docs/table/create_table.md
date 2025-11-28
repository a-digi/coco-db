# Tabellen anlegen (Create Table)

## Endpunkt

```
POST /api/databases/{dbname}/tables
```

## Beispiel-Request

```bash
curl -X POST http://localhost:8080/api/databases/testdb/tables \
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
  "message": "Tabelle erfolgreich angelegt"
}
```

## Beispiel-Response (Fehler: Tabelle existiert)

```json
{
  "success": false,
  "error": {
    "code": "ERR_TABLE_EXISTS",
    "message": "Tabelle existiert bereits"
  }
}
```

## Hinweise
- Der Datenbankname wird aus der URL extrahiert.
- Die vollständige meta.json wird im Erfolgsfall zurückgegeben.
- Fehler werden als JSON mit success: false und error-Objekt zurückgegeben.

