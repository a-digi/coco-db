# Tabellen auflisten (List Tables)

## Endpunkt

```
GET /api/databases/{dbname}/tables
```

## Beispiel-Request

```bash
curl -X GET http://localhost:2022/api/databases/testdb/tables
```

## Beispiel-Response (Erfolg)

```json
{
  "success": true,
  "data": {
    "tables": [
      "users",
      "orders"
    ]
  },
  "httpCode": 200,
  "executionTime": "..."
}
```

## Beispiel-Response (Fehler: Datenbank nicht gefunden)

```json
{
  "success": false,
  "error": {
    "code": "ERR_DB_NOT_FOUND",
    "message": "Not Found: Datenbank nicht gefunden"
  },
  "httpCode": 404,
  "executionTime": "..."
}
```

## Hinweise
- Der Datenbankname wird aus der URL extrahiert.
- Die Antwort enthält ein Array aller Tabellennamen in der Datenbank.
- Fehler werden als JSON mit success: false und error-Objekt zurückgegeben.
- Die Felder `httpCode` und `executionTime` sind immer enthalten.

