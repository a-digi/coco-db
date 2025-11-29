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
- Die Metadaten werden in `{datadir}/{db}/tables.json` aktualisiert.
- Fehler werden als JSON mit success: false und error-Objekt zurückgegeben.
- Die Felder `httpCode` und `executionTime` sind immer enthalten.
