# API-Design: JSON-Dokumentdatenbank

## 1. Datenbankverwaltung
### 1.1 Datenbank anlegen
- **POST** `/api/databases`
- Request: `{ "name": "string" }`
- Response: `{ "success": true, "database": { ... } }`

### 1.2 Datenbank löschen
- **DELETE** `/api/databases/{dbname}`
- Response: `{ "success": true }`

### 1.3 Datenbanken auflisten
- **GET** `/api/databases`
- Response: `{ "databases": [ ... ] }`

## 2. Tabellenverwaltung
### 2.1 Tabelle anlegen
- **POST** `/api/databases/{dbname}/tables`
- Request: `{ "name": "string", "schema": { ... } }`
- Response: `{ "success": true, "table": { ... } }`

### 2.2 Tabelle löschen
- **DELETE** `/api/databases/{dbname}/tables/{tablename}`
- Response: `{ "success": true }`

### 2.3 Tabellen auflisten
- **GET** `/api/databases/{dbname}/tables`
- Response: `{ "tables": [ ... ] }`

## 3. Dokumentenverwaltung (CRUD)
### 3.1 Dokument anlegen
- **POST** `/api/databases/{dbname}/tables/{tablename}/documents`
- Request: `{ ... }` (JSON-Dokument)
- Response: `{ "success": true, "id": "string" }`

### 3.2 Dokument abrufen
- **GET** `/api/databases/{dbname}/tables/{tablename}/documents/{id}`
- Response: `{ ... }` (JSON-Dokument)

### 3.3 Dokument aktualisieren
- **PUT** `/api/databases/{dbname}/tables/{tablename}/documents/{id}`
- Request: `{ ... }` (JSON-Dokument)
- Response: `{ "success": true }`

### 3.4 Dokument löschen
- **DELETE** `/api/databases/{dbname}/tables/{tablename}/documents/{id}`
- Response: `{ "success": true }`

### 3.5 Dokumente abfragen (Filter, Query)
- **POST** `/api/databases/{dbname}/tables/{tablename}/query`
- Request: `{ "filter": { ... }, "limit": 10, "offset": 0 }`
- Response: `{ "documents": [ ... ] }`

## 4. Indexverwaltung
### 4.1 Index anlegen
- **POST** `/api/databases/{dbname}/tables/{tablename}/indexes`
- Request: `{ "field": "string", "type": "primary|secondary", "unique": true|false }`
- Response: `{ "success": true }`

### 4.2 Index löschen
- **DELETE** `/api/databases/{dbname}/tables/{tablename}/indexes/{field}`
- Response: `{ "success": true }`

### 4.3 Indizes auflisten
- **GET** `/api/databases/{dbname}/tables/{tablename}/indexes`
- Response: `{ "indexes": [ ... ] }`

## 5. Event- und Transaktionsmanagement
### 5.1 Events abrufen
- **GET** `/api/events?since=timestamp`
- Response: `{ "events": [ ... ] }`

### 5.2 Transaktion starten
- **POST** `/api/transactions/start`
- Response: `{ "transaction_id": "string" }`

### 5.3 Transaktion commit
- **POST** `/api/transactions/{id}/commit`
- Response: `{ "success": true }`

### 5.4 Transaktion rollback
- **POST** `/api/transactions/{id}/rollback`
- Response: `{ "success": true }`

## 6. Authentifizierung (optional)
- **POST** `/api/auth/login`
- Request: `{ "username": "string", "password": "string" }`
- Response: `{ "token": "string" }`

## Fehlerbehandlung
- Fehler werden als JSON mit `error`-Feld und passendem HTTP-Statuscode zurückgegeben:
  - `{ "error": "Beschreibung des Fehlers" }`

## Hinweise
- Alle Ein- und Ausgaben erfolgen im JSON-Format.
- Erweiterbarkeit für zukünftige Features (z.B. Audit, Health-Check, Volltextsuche) ist vorgesehen.
- Authentifizierung/Autorisierung kann optional aktiviert werden.

