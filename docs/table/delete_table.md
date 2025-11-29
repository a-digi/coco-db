# Tabelle löschen (Delete Table)

## Endpunkt

```
DELETE /api/databases/{dbname}/tables/{tablename}
```

## Beispiel-Request

```bash
curl -X DELETE http://localhost:2022/api/databases/testdb/tables/users
```

## Beispiel-Response (Erfolg)

```json
{
  "success": true,
  "data": {
    "table": "users"
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
- Die Tabelle wird aus `{datadir}/{db}/tables.json` entfernt und das Tabellenverzeichnis gelöscht.
- Fehler werden als JSON mit success: false und error-Objekt zurückgegeben.
- Die Felder `httpCode` und `executionTime` sind immer enthalten.
