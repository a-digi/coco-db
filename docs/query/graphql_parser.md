# GraphQL-Style Query-Parser Dokumentation

## Ziel

Der Parser `ParseGraphQLStyleQuery` unterstützt jetzt explizite Feldselektion ("fields") auf Root-Ebene und in Joins. Damit kann ein echtes GraphQL-ähnliches Query wie folgt verarbeitet werden:

## Beispiel-Query (GraphQL-Stil)

```graphql
query {
  users(
    filter: {
      created_at: { gte: "2025-11-01T00:00:00Z", lte: "2025-11-30T23:59:59Z" }
      email: { like: "*@gmail.com" }
      age: { gte: 18 }
    }
    limit: 10
    offset: 0
    sort: ["created_at", "age"]
    join: [
      {
        table: "orders"
        on: { user_id: "id" }
        filter: { status: "open" }
        fields: ["id", "amount", "status"]
      }
    ]
    fields: ["id", "name", "email", "created_at", "age", "orders"]
  ) {
    id
    name
    email
    created_at
    age
    orders {
      id
      amount
      status
    }
  }
}
```

## JSON-Format nach dem Parsen

```json
{
  "filter": {
    "created_at": { "gte": "2025-11-01T00:00:00Z", "lte": "2025-11-30T23:59:59Z" },
    "email": { "like": "*@gmail.com" },
    "age": { "gte": 18 }
  },
  "limit": 10,
  "offset": 0,
  "sort": ["created_at", "age"],
  "fields": ["id", "name", "email", "created_at", "age", "orders"],
  "join": [
    {
      "table": "orders",
      "on": { "user_id": "id" },
      "filter": { "status": "open" },
      "fields": ["id", "amount", "status"]
    }
  ]
}
```

## Hinweise
- Die Feldauswahl ("fields") ist jetzt auf Root-Ebene und in Joins möglich.
- Die Response kann so gestaltet werden, dass sie exakt die gewünschten Felder enthält.
- Die Query ist weiterhin JSON-basiert, aber die Feldlogik entspricht echtem GraphQL.

---

**Letzte Änderung: 30.11.2025**

