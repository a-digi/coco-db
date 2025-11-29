# Eintrag aktualisieren (Update Entry)

## Endpunkt

```
PUT /api/databases/{dbname}/tables/{tablename}/entries/{entryid}
```

## Beschreibung

Aktualisiert einen bestehenden Eintrag in der angegebenen Tabelle einer Datenbank. Die neue Version des Eintrags wird validiert, versioniert und persistent gespeichert. Die alte Version bleibt erhalten (Versionierung).

## Request

- **Pfadparameter:**
  - `dbname` (string): Name der Datenbank
  - `tablename` (string): Name der Tabelle
  - `entryid` (string): ID des zu aktualisierenden Eintrags
- **Body:**
  - JSON-Objekt mit den neuen Feldern des Eintrags, z. B.:

```json
{
  "foo": "bar",
  "num": 43
}
```

## Antwort

- **Erfolg:**
  - HTTP-Status: `200 OK`
  - Body:
    ```json
    {
      "data": {
        "entryId": "<entryid>",
        "version": <neueVersion>,
        "updated_at": "<timestamp>"
      },
      "httpCode": 200
    }
    ```
- **Fehler:**
  - HTTP-Status: `400 Bad Request` (bei Validierungsfehlern oder ungültigem JSON)
    ```json
    {
      "error": "ERR_INVALID_JSON",
      "message": "Ungültiges JSON: ..."
    }
    ```
  - HTTP-Status: `404 Not Found` (wenn der Eintrag nicht existiert)
    ```json
    {
      "error": "ERR_ENTRY_NOT_FOUND",
      "message": "Eintrag nicht gefunden"
    }
    ```
  - HTTP-Status: `400 Bad Request` (bei anderen Fehlern beim Aktualisieren)
    ```json
    {
      "error": "ERR_UPDATE_ENTRY",
      "message": "...Fehlerbeschreibung..."
    }
    ```

## Beispielaufruf

```
curl -X PUT \
  http://localhost:2022/api/databases/testdb/tables/users/entries/123e4567-e89b-12d3-a456-426614174000 \
  -H 'Content-Type: application/json' \
  -d '{"username": "alice", "age": 31}'
```

## Hinweise
- Die neue Version des Eintrags wird automatisch versioniert und gespeichert.
- Die Validierung erfolgt anhand des Tabellenschemas (meta.json).
- Die alte Version bleibt erhalten und ist über die Versionierung nachvollziehbar.
- Bei Fehlern wird eine strukturierte Fehlermeldung zurückgegeben.

