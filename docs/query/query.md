# Query-API

## Übersicht
Die Query-API ermöglicht das Filtern, Suchen und Sortieren von Einträgen in Tabellen oder global über alle Tabellen hinweg. Sie unterstützt Filter, Sortierung, Paginierung und Joins.

**Basis-URL:**
```
http://localhost:2022
```

---

## 1. Tabellen-Query

### Endpunkt
```
POST /api/databases/{dbname}/tables/{tablename}/query
```

### Beispiel-Request
```bash
curl -X POST http://localhost:2022/api/databases/testdb/tables/users/query \
  -H "Content-Type: application/json" \
  -d '{
    "filter": { "id": "1" },
    "limit": 10,
    "offset": 0,
    "sort": ["age"]
  }'
```

### Beispiel-Response (Erfolg)
```json
{
  "success": true,
  "data": [
    { "id": "1", "name": "Max", "age": 22 }
  ],
  "httpCode": 200,
  "executionTime": "..."
}
```

### Beispiel-Response (Fehler)
```json
{
  "success": false,
  "error": {
    "code": "ERR_INVALID_QUERY",
    "message": "Missing required field: filter"
  },
  "httpCode": 400,
  "executionTime": "..."
}
```

---

## 2. Globale Query

### Endpunkt
```
POST /api/{dbname}/query
```

### Beispiel-Request
```bash
curl -X POST http://localhost:2022/api/testdb/query \
  -H "Content-Type: application/json" \
  -d '{
    "filter": { "email": { "like": "max" } },
    "limit": 5,
    "sort": ["email"]
  }'
```

### Beispiel-Response (Erfolg)
```json
{
  "success": true,
  "data": [
    { "id": "1", "email": "max@foo.de", "age": 22, "database": "testdb", "table": "users" }
  ],
  "httpCode": 200,
  "executionTime": "..."
}
```

---

## 2. Globale Query (GraphQL-Stil)

### Endpunkt
```
POST /api/{dbname}/query
```

### Beispiel-Request (GraphQL-Style)
```bash
curl -X POST http://localhost:2022/api/testdb/query \
  -H "Content-Type: application/json" \
  -d '{
    "users": {
      "filter": { "email": { "like": "*@gmail.com" } },
      "limit": 5,
      "sort": ["email"],
      "join": [
        {
          "table": "user_roles",
          "on": { "id": "user_id" },
          "fields": ["role_id"],
          "join": [
            {
              "table": "roles",
              "on": { "role_id": "id" },
              "fields": ["name", "description"]
            }
          ]
        }
      ]
    }
  }'
```

### Beispiel-Response (GraphQL-Style)
```json
{
  "success": true,
  "data": [
    {
      "id": "1",
      "name": "Max Mustermann",
      "email": "max@gmail.com",
      "user_roles": [
        {
          "role_id": "admin",
          "roles": [
            { "name": "Admin", "description": "Administrator mit allen Rechten" }
          ]
        }
      ]
    }
    // ...weitere Nutzer...
  ],
  "httpCode": 200,
  "executionTime": "..."
}
```

**Hinweis:**
- Die Response ist ein Array mit verschachtelten Objekten, die der Join-Struktur der Query entsprechen.
- Die Query kann beliebig tief verschachtelte Joins enthalten.
- Die Felder und Filter können an die tatsächlichen Daten angepasst werden.

---

## 3. Request-Body (Schema)

```json
{
  "filter": { "feld": "wert", "alter": { "gte": 18 } },
  "limit": 10,
  "offset": 0,
  "sort": ["feld1", "feld2"],
  "join": [ { "table": "andere_tabelle", "on": { "user_id": "id" } } ]
}
```
- **filter**: Filterbedingungen (siehe unten)
- **limit**: Maximale Anzahl Ergebnisse (optional)
- **offset**: Offset für Paginierung (optional)
- **sort**: Sortierfelder (optional)
- **join**: Joins mit anderen Tabellen (optional, nur global)

### Unterstützte Filter-Operatoren
- Gleichheit: `{ "feld": "wert" }`
- Bereich: `{ "feld": { "gte": 18, "lte": 99 } }`
- LIKE/Partial: `{ "feld": { "like": "max" } }`
- Fulltext: `{ "feld": { "fulltext": "foo bar" } }`

---

## 4. Hinweise
- Die Query-API gibt immer die Felder `success`, `data` (oder `error`), `httpCode` und `executionTime` zurück.
- Fehler werden als JSON mit success: false und error-Objekt zurückgegeben.
- Die globale Query durchsucht alle Tabellen/Datenbanken (aktuell Demo: testdb/users).
- Die Ausführungszeit (`executionTime`) ist immer enthalten.

---

## 5. Beispiele für Filter

```json
{
  "filter": {
    "created_at": { "gte": "2025-11-01T00:00:00Z", "lte": "2025-11-30T23:59:59Z" },
    "email": "foo@bar.de",
    "age": { "gte": 18 },
    "name": { "like": "Max" }
  },
  "limit": 10,
  "offset": 0,
  "sort": ["created_at", "age"]
}
```

**Erklärung:**
- Das Beispiel filtert alle Einträge, die im Zeitraum vom 1.11.2025 bis 30.11.2025 erstellt wurden (`created_at`), zusätzlich nach E-Mail, Mindestalter und Namens-Pattern.
- Die Filter-Operatoren `gte` (größer/gleich) und `lte` (kleiner/gleich) funktionieren auch für Felder vom Typ `date` (ISO 8601).

---

## 6. Fehlerbehandlung
- Fehlerhafte Requests liefern einen passenden httpCode und ein error-Objekt.
- Beispiel: Fehlendes Pflichtfeld `filter` → 400 Bad Request

---

## 7. Weitere Hinweise
- Die API ist für POST-Requests mit JSON-Body ausgelegt.
- Für produktive Nutzung: Siehe Hinweise zu Joins und globaler Query in der Hauptdokumentation.
