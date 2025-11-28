# Datenbank anlegen (Create Database)

## Endpunkt

```
POST /api/databases
```

## Beispiel-Request

```bash
curl -X POST http://localhost:2022/api/databases \
  -H "Content-Type: application/json" \
  -d '{
    "name": "testdb"
  }'
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

## Beispiel-Response (Fehler: Datenbank existiert)

```json
{
  "success": false,
  "error": {
    "code": "ERR_DB_EXISTS",
    "message": "Datenbank existiert bereits"
  },
  "httpCode": 409,
  "executionTime": "..."
}
```

## Hinweise
- Der Name der Datenbank wird im JSON-Body übergeben.
- Fehler werden als JSON mit success: false und error-Objekt zurückgegeben.
- Die Felder `httpCode` und `executionTime` sind immer enthalten.
