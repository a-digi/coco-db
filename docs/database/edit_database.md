# Datenbank bearbeiten (Edit Database)

## Endpunkt

```
PUT /api/databases/{dbname}
```

## Beispiel-Request

```bash
curl -X PUT http://localhost:2022/api/databases/testdb \
  -H "Content-Type: application/json" \
  -d '{
    "newName": "neuername"
  }'
```

## Beispiel-Response (Erfolg)

```json
{
  "success": true,
  "data": {
    "oldName": "testdb",
    "newName": "neuername"
  },
  "httpCode": 200,
  "executionTime": "..."
}
```

## Beispiel-Response (Fehler: Name existiert bereits)

```json
{
  "success": false,
  "error": {
    "code": "ERR_DB_EXISTS",
    "message": "Ziel-Datenbankname existiert bereits"
  },
  "httpCode": 409,
  "executionTime": "..."
}
```

## Hinweise
- Der neue Name wird im JSON-Body übergeben.
- Fehler werden als JSON mit success: false und error-Objekt zurückgegeben.
- Die Felder `httpCode` und `executionTime` sind immer enthalten.
