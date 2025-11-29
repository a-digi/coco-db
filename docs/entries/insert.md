# Eintrag anlegen (Insert Entry)

## Endpunkt

```
POST /api/databases/{dbname}/tables/{tablename}/entries
```

## Beschreibung

Legt einen neuen Eintrag in der angegebenen Tabelle einer Datenbank an. Der Eintrag wird validiert, versioniert und persistent gespeichert.

## Request

- **Pfadparameter:**
  - `dbname` (string): Name der Datenbank
  - `tablename` (string): Name der Tabelle
- **Body:**
  - JSON-Objekt mit den Feldern des neuen Eintrags, z. B.:

```json
{
  "foo": "bar",
  "num": 42
}
```

## Antwort

- **Erfolg:**
  - HTTP-Status: `201 Created`
  - Body:
    ```json
    {
      "data": {
        "entryId": "<uuid>"
      },
      "httpCode": 201
    }
    ```
- **Fehler:**
  - HTTP-Status: `400 Bad Request` (bei Validierungsfehlern oder ungültigem JSON)
  - Body:
    ```json
    {
      "error": "ERR_INVALID_JSON",
      "message": "Ungültiges JSON: ..."
    }
    ```
  - HTTP-Status: `400 Bad Request` (bei Fehlern beim Einfügen)
    ```json
    {
      "error": "ERR_INSERT_ENTRY",
      "message": "...Fehlerbeschreibung..."
    }
    ```

## Beispielaufruf

```
curl -X POST \
  http://localhost:2022/api/databases/testdb/tables/users/entries \
  -H 'Content-Type: application/json' \
  -d '{"username": "alice", "age": 30}'
```

## Hinweise
- Der Eintrag wird automatisch versioniert und mit einer eindeutigen ID gespeichert.
- Die Validierung erfolgt anhand des Tabellenschemas (meta.json).
- Bei Fehlern wird eine strukturierte Fehlermeldung zurückgegeben.

