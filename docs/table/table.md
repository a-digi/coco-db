# Tabellen-API Übersicht

Die Tabellen-API ermöglicht das Anlegen, Bearbeiten, Löschen und Auflisten von Tabellen innerhalb einer Datenbank.

## Übersicht der Endpunkte

- **Tabelle anlegen:** `POST /api/databases/{dbname}/tables`
- **Tabellen auflisten:** `GET /api/databases/{dbname}/tables`
- **Tabellen-Metadaten abrufen:** `GET /api/databases/{dbname}/tables/{tablename}`
- **Tabelle aktualisieren:** `PUT /api/databases/{dbname}/tables/{tablename}`
- **Tabelle löschen:** `DELETE /api/databases/{dbname}/tables/{tablename}`

## Hinweise
- Die Metadaten aller Tabellen einer Datenbank werden in `{datadir}/{db}/tables.json` verwaltet.
- Jede Tabelle besitzt zusätzlich eine eigene `meta.json` im jeweiligen Tabellenverzeichnis.
- Fehler werden als JSON mit success: false und error-Objekt zurückgegeben.
- Die Felder `httpCode` und `executionTime` sind immer enthalten.

