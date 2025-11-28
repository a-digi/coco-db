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
  "message": "Datenbank erfolgreich gelöscht"
}
```

## Beispiel-Response (Fehler: Nicht gefunden)

```json
{
  "success": false,
  "error": {
    "code": "ERR_DB_NOT_FOUND",
    "message": "Datenbank nicht gefunden"
  }
}
```

## Hinweise
- Der Name der zu löschenden Datenbank wird in der URL angegeben.
- Fehler werden als JSON mit success: false und error-Objekt zurückgegeben.
