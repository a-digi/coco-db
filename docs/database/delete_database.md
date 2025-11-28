# Datenbank löschen (Delete Database)

## Endpunkt

```
DELETE /api/databases/{dbname}
```

## Beispiel-Request

```bash
curl -X DELETE http://localhost:2022/api/databases/testdb
```

## Beispiel-Response (Erfolg)

```json
{
  "success": true,
  "data": {
    "name": "testdb"
  },
  "httpCode": 200,
  "executionTime": "..."
}
```

## Beispiel-Response (Fehler: Nicht gefunden)

```json
{
  "success": false,
  "error": {
    "code": "ERR_DB_NOT_FOUND",
    "message": "Datenbank nicht gefunden"
  },
  "httpCode": 404,
  "executionTime": "..."
}
```

## Hinweise
- Der Name der zu löschenden Datenbank wird in der URL angegeben.
- Fehler werden als JSON mit success: false und error-Objekt zurückgegeben.
- Die Felder `httpCode` und `executionTime` sind immer enthalten.
