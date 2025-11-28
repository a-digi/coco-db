# Datenbank lesen (Read Database)

## Endpunkt

```
GET /api/databases
```

## Beispiel-Request

```bash
curl -X GET http://localhost:2022/api/databases
```

## Beispiel-Response (Erfolg)

```json
{
  "success": true,
  "data": {
    "databases": [
      "testdb",
      "andere_db"
    ]
  },
  "httpCode": 200,
  "executionTime": "..."
}
```

## Hinweise
- Gibt eine Liste aller vorhandenen Datenbanken zurück.
- Fehler werden als JSON mit success: false und error-Objekt zurückgegeben.
- Die Felder `httpCode` und `executionTime` sind immer enthalten.
