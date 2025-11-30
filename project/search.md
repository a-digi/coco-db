# Implementierung einer Search-Funktionalität (API-Query)

Diese Anleitung beschreibt Schritt für Schritt, wie eine flexible Search-Funktionalität für die Datenbank implementiert wird, die komplexe Filter, Sortierung, Pagination, Joins und Feldselektion unterstützt.

## Ziel
Eine API, die es ermöglicht, mit einer Query wie unten beschrieben, Datenbankeinträge abzufragen und strukturierte Ergebnisse zurückzugeben.

Beispiel-Query:

```
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

---

## Schritt-für-Schritt-Anleitung

- [ ] **1. Query-Parser erweitern/erstellen**
    - Query-Objekt muss Felder wie `filter`, `limit`, `offset`, `sort`, `join`, `fields` unterstützen
    - Validierung der Query-Struktur und Datentypen

- [ ] **2. Filter-Engine nutzen/erweitern**
    - Filter auf die Haupttabelle anwenden (z.B. users)
    - Unterstützung für Operatoren wie `gte`, `lte`, `like`, etc.

- [ ] **3. Sortierung, Limit und Offset anwenden**
    - Ergebnisse nach den gewünschten Feldern sortieren
    - Pagination (limit/offset) anwenden

- [ ] **4. Joins implementieren/aufrufen**
    - Für jedes Join-Objekt: 
        - Join-Tabelle laden
        - Join-Bedingung(en) anwenden
        - Join-Filter anwenden
        - Nur gewünschte Felder übernehmen
        - Rekursive Joins unterstützen (bis max. Tiefe)

- [ ] **5. Feldselektion anwenden**
    - Nur die in `fields` angegebenen Felder in das Ergebnis aufnehmen
    - Verschachtelte Felder (z.B. `orders`) korrekt abbilden

- [ ] **6. Antwortstruktur aufbauen**
    - Ergebnisse als Array von Objekten zurückgeben
    - Fehlerbehandlung und Statuscodes implementieren

- [ ] **7. API-Endpunkt definieren**
    - POST `/api/{dbname}/search` oder `/api/{dbname}/query`
    - Request-Body: Query-Objekt (JSON)
    - Response: Gefilterte, sortierte, paginierte und ggf. gejointe Ergebnisse

- [ ] **8. Tests schreiben**
    - Unit- und Integrationstests für alle Komponenten (Parser, Filter, Join, Sortierung, Pagination, Feldselektion)
    - Beispiel-Queries und erwartete Ergebnisse abdecken

- [ ] **9. Dokumentation und Beispiele ergänzen**
    - API-Dokumentation aktualisieren
    - Beispiel-Requests und -Responses dokumentieren

---

## Hinweise
- Bestehende Komponenten wie `FilterEngine`, `QueryHandler`, etc. wiederverwenden, wo möglich
- Auf Performance und Speicherverbrauch achten (Indexnutzung, Streaming, etc.)
- Fehler und Grenzfälle (z.B. ungültige Filter, zu tiefe Joins) sauber behandeln

---

## Weiterführende Links
- [docs/query/query.md](../docs/query/query.md)
- [src/query/](../src/query/)
- [src/query/filter_engine.go](../src/query/filter_engine.go)
- [src/query/query_handlers.go](../src/query/query_handlers.go)
