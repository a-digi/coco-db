# Suche (Search) – Query-API

Dieses Dokument beschreibt die Nutzung der Search-API mit Beispielen für Filter, Joins, Sortierung, Paginierung und Feldselektion.

## 1. Einfache Suche mit Filter

**Beispiel:** Finde alle Nutzer mit einer Gmail-Adresse, die nach dem 1.11.2025 erstellt wurden und mindestens 18 Jahre alt sind.

```graphql
query {
  users(
    filter: {
      created_at: { gte: "2025-11-01T00:00:00Z", lte: "2025-11-30T23:59:59Z" },
      email: { like: "*@gmail.com" },
      age: { gte: 18 }
    },
    limit: 10,
    offset: 0,
    sort: ["created_at", "age"],
    fields: ["id", "name", "email", "created_at", "age"]
  ) {
    id
    name
    email
    created_at
    age
  }
}
```

## 2. Suche mit Join (user_roles, roles)

**Beispiel:** Hole Nutzer mit ihren Rollen (über user_roles und roles):

```graphql
query {
  users(
    filter: {
      email: { like: "*@gmail.com" }
    },
    join: [
      {
        table: "user_roles",
        on: { user_id: "id" },
        join: [
          {
            table: "roles",
            on: { id: "role_id" },
            fields: ["id", "name"]
          }
        ],
        fields: ["role_id", "user_id", "roles"]
      }
    ],
    fields: ["id", "name", "email", "user_roles"]
  ) {
    id
    name
    email
    user_roles {
      role_id
      user_id
      roles {
        id
        name
      }
    }
  }
}
```

## 3. Komplexe Suche mit Filter, Join, Sortierung und Paginierung

```graphql
query {
  users(
    filter: {
      created_at: { gte: "2025-11-01T00:00:00Z", lte: "2025-11-30T23:59:59Z" },
      email: { like: "*@gmail.com" },
      age: { gte: 18 }
    },
    limit: 10,
    offset: 0,
    sort: ["created_at", "age"],
    join: [
      {
        table: "user_roles",
        on: { user_id: "id" },
        join: [
          {
            table: "roles",
            on: { id: "role_id" },
            fields: ["id", "name"]
          }
        ],
        fields: ["role_id", "user_id", "roles"]
      }
    ],
    fields: ["id", "name", "email", "created_at", "age", "user_roles"]
  ) {
    id
    name
    email
    created_at
    age
    user_roles {
      role_id
      user_id
      roles {
        id
        name
      }
    }
  }
}
```

## 4. CURL-Beispiel

```sh
curl -v -X POST http://localhost:2022/api/databases/poseidon/search -d 'query { users( filter: { created_at: { gte: "2025-11-01T00:00:00Z", lte: "2025-11-30T23:59:59Z" }, email: { like: "*@gmail.com" }, age: { gte: 18 } }, join: [ { table: "user_roles", on: { user_id: "id" }, join: [ { table: "roles", on: { id: "role_id" }, fields: ["id", "name"] } ], fields: ["role_id", "user_id", "roles"] } ], limit: 10, offset: 0, fields: ["id", "name", "email", "user_roles"] ) { id name email user_roles { role_id user_id roles { id name } } } }'
```

## 5. Hinweise

- **Filter:** Unterstützt Operatoren wie `eq`, `neq`, `gte`, `lte`, `like` (z.B. für Strings mit Wildcards `*`).
- **Joins:** Verschachtelte Joins sind möglich (siehe user_roles → roles).
- **Felder:** Mit `fields` kann die Rückgabe auf bestimmte Felder eingeschränkt werden.
- **Sortierung:** Mit `sort` kann nach mehreren Feldern sortiert werden.
- **Paginierung:** Mit `limit` und `offset` können Ergebnisse paginiert werden.

---

Weitere Beispiele und Spezialfälle können gerne ergänzt werden.
