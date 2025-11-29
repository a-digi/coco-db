# Eintrag löschen (Delete Entry)

## Endpunkt

```
DELETE /api/databases/{dbname}/tables/{tablename}/entries/{entryid}
```

## Beschreibung

Löscht einen bestehenden Eintrag aus der angegebenen Tabelle einer Datenbank unwiderruflich (Hard Delete). Dabei werden alle Versionen und Metadaten des Eintrags entfernt.

## Request

- **Pfadparameter:**
  - `dbname` (string): Name der Datenbank
  - `tablename` (string): Name der Tabelle
  - `entryid` (string): ID des zu löschenden Eintrags

## Antwort

- **Erfolg:**
  - HTTP-Status: `200 OK`
  - Body:
    ```json
    {
      "data": {
        "entryId": "<entryid>",
        "deleted_at": "<timestamp>"
      },
      "httpCode": 200
    }
    ```
- **Fehler:**
  - HTTP-Status: `404 Not Found` (wenn der Eintrag nicht existiert)
    ```json
    {
      "error": "ERR_ENTRY_NOT_FOUND",
      "message": "Eintrag nicht gefunden"
    }
    ```
  - HTTP-Status: `400 Bad Request` (bei anderen Fehlern beim Löschen)
    ```json
    {
      "error": "ERR_DELETE_ENTRY",
      "message": "...Fehlerbeschreibung..."
    }
    ```

## Beispielaufruf

```
curl -X DELETE \
  http://localhost:2022/api/databases/testdb/tables/users/entries/123e4567-e89b-12d3-a456-426614174000
```

## Hinweise
- Der Eintrag wird vollständig und unwiderruflich gelöscht (Hard Delete).
- Es werden alle Versionen und Metadaten des Eintrags entfernt.
- Bei Fehlern wird eine strukturierte Fehlermeldung zurückgegeben.

