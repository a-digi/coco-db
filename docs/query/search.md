# Globale Suche mit Joins (JSON-API mit verschachtelten Joins)

Dieses Beispiel zeigt, wie du die globale Suchfunktion der Query-API verwendest, um tabellenübergreifend nach Einträgen zu suchen und dabei Joins ähnlich wie in GraphQL zu nutzen. **Achtung:** Die API ist JSON-basiert und nicht vollwertig wie GraphQL – Einstiegspunkt ist nicht eine einzelne Tabelle, sondern es werden alle Tabellen durchsucht.

## Endpunkt
```
POST /api/{dbname}/query
```

## Beispiel: Globale Suche mit Join (JSON-API mit Joins)

### Konzeptuelle Abfrage (angelehnt an GraphQL, aber als JSON)

Die API nimmt ein JSON-Objekt mit den Feldern `filter`, `limit`, `offset`, `sort`, `join` entgegen. Die Query wird für alle Tabellen der Datenbank ausgeführt. Das Ergebnis ist ein Array mit Treffern pro Tabelle, nicht ein verschachteltes Objekt wie bei GraphQL.

### Beispiel-curl-Request

```bash
curl -X POST http://localhost:2022/api/poseidon/query \
  -H "Content-Type: application/json" \
  -d '{
    "filter": {
      "created_at": { "gte": "2025-11-01T00:00:00Z", "lte": "2025-11-30T23:59:59Z" },
      "email": { "like": "*@gmail.com" },
      "age": { "gte": 18 }
    },
    "limit": 10,
    "offset": 0,
    "sort": ["created_at", "age"],
    "join": [
      {
        "table": "orders",
        "on": { "user_id": "id" },
        "filter": { "status": "open" },
        "fields": ["id", "amount", "status"]
      }
    ]
  }'
```

### Beispiel-Response (Achtung: Ergebnisse sind nach Tabellen gruppiert)

```json
{
  "success": true,
  "data": [
    {
      "table": "users",
      "entries": [
        {
          "id": "1",
          "name": "Max Mustermann",
          "email": "max@gmail.com",
          "created_at": "2025-11-10T12:00:00Z",
          "age": 22,
          "orders": [
            { "id": "101", "amount": 99.99, "status": "open" },
            { "id": "102", "amount": 49.99, "status": "open" }
          ]
        }
      ]
    },
    // ...andere Tabellen...
  ],
  "httpCode": 200,
  "meta": {
    "totalCount": 1,
    "limit": 10,
    "offset": 0
  }
}
```

## Hinweise
- Die globale Suche kann beliebig verschachtelte Joins enthalten (empfohlen: max. Tiefe 5–10).
- Einstiegspunkt ist nicht eine einzelne Tabelle, sondern alle Tabellen der Datenbank.
- Die Join-Definitionen erlauben Filter, Feld-Auswahl und Mapping der Join-Bedingung auf jeder Ebene.
- Die Response ist nach Tabellen gruppiert, nicht wie bei GraphQL nach Einstiegspunkt verschachtelt.
- Die Syntax ist JSON-basiert, nicht GraphQL.
- Unterstützte Operatoren, Filter und Sortierung siehe Hauptdokumentation.
