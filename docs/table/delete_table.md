# Tabelle löschen (Delete Table)
- Die Felder `httpCode` und `executionTime` sind immer enthalten.
- Fehler werden als JSON mit success: false und error-Objekt zurückgegeben.
- Die Tabelle wird aus `{datadir}/{db}/tables.json` entfernt und das Tabellenverzeichnis gelöscht.
## Hinweise

```
}
  "executionTime": "..."
  "httpCode": 404,
  },
    "message": "Tabelle nicht gefunden"
    "code": "ERR_TABLE_NOT_FOUND",
  "error": {
  "success": false,
{
```json

## Beispiel-Response (Fehler: Tabelle nicht gefunden)

```
}
  "executionTime": "..."
  "httpCode": 200,
  },
    "table": "users"
  "data": {
  "success": true,
{
```json

## Beispiel-Response (Erfolg)

```
curl -X DELETE http://localhost:2022/api/databases/testdb/tables/users
```bash

## Beispiel-Request

```
DELETE /api/databases/{dbname}/tables/{tablename}
```

## Endpunkt


