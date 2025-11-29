# Tabellen auflisten (List Tables)

## Endpunkt

```
GET /api/databases/{dbname}/tables
```

## Beispiel-Request

```bash
curl http://localhost:2022/api/databases/testdb/tables
```

## Beispiel-Response (Erfolg)

```json
{
  "success": true,
  "data": {
    "tables": ["users", "products"]
  },
  "httpCode": 200,
  "executionTime": "..."
}
```

## Beispiel-Response (leere Liste)

```json
{
  "success": true,
  "data": {
    "tables": []
  },
  "httpCode": 200,
  "executionTime": "..."
}
```

## Hinweise
- Die Tabellennamen werden aus `{datadir}/{db}/tables.json` gelesen.
- Existiert die Datei nicht, wird sie als leeres Array angelegt.
- Fehler werden als JSON mit success: false und error-Objekt zurückgegeben.

