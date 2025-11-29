# Tabellen-API Übersicht

Die Tabellen-API ermöglicht das Anlegen, Bearbeiten, Löschen und Auflisten von Tabellen innerhalb einer Datenbank.

## Übersicht der Endpunkte

- **Tabelle anlegen:** `POST /api/databases/{dbname}/tables`
- **Tabellen auflisten:** `GET /api/databases/{dbname}/tables`
- **Tabellen-Metadaten abrufen:** `GET /api/databases/{dbname}/tables/{tablename}`
- **Tabelle aktualisieren:** `PUT /api/databases/{dbname}/tables/{tablename}`
- **Tabelle löschen:** `DELETE /api/databases/{dbname}/tables/{tablename}`

## Einzelne Tabelle abrufen (Get Table)

### Endpunkt

```
GET /api/databases/{dbname}/tables/{tablename}
```

### Beispiel-Request

```bash
curl http://localhost:2022/api/databases/testdb/tables/users
```

### Beispiel-Response (Erfolg)

```json
{
  "success": true,
  "data": {
    "tableName": "users",
    "schemaVersion": 1,
    "fields": [
      { "name": "id", "type": "string", "required": true },
      { "name": "email", "type": "string", "required": true }
    ],
    "indexes": [
      { "name": "primary", "type": "primary", "fields": ["id"], "unique": true }
    ],
    "options": { "documentLimit": 1000 }
  },
  "httpCode": 200,
  "executionTime": "..."
}
```

### Beispiel-Response (Fehler: Tabelle nicht gefunden)

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
- Die Metadaten aller Tabellen einer Datenbank werden in `{datadir}/{db}/tables.json` verwaltet.
- Jede Tabelle besitzt zusätzlich eine eigene `meta.json` im jeweiligen Tabellenverzeichnis.
- Fehler werden als JSON mit success: false und error-Objekt zurückgegeben.
- Die Felder `httpCode` und `executionTime` sind immer enthalten.
