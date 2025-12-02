# Roadmap: Join- und Index-Optimierung für speicheroptimierte Query-Engine

## 1. Architektur: Index-Nutzung und Filter-Engine
- Alle Indexdateien werden beim Serverstart in die zentrale Registry (RAM) geladen.
- Die Registry verwendet zusammengesetzte Schlüssel: `db.table.index` (z.B. `poseidon.users.email`).
- Die FilterEngine nutzt für Filter und Joins die RAM-Indizes, sofern vorhanden. Nur die IDs, die im Index gefunden werden, werden geöffnet.
- Bereichsfilter, Like-Filter und "in"-Filter werden direkt auf die RAM-Indizes angewendet.

## 2. Join-Logik: Optimale Nutzung der Indizes
- Nach dem Filtern der Haupttabelle (z.B. users) werden nur die passenden User-IDs weitergegeben.
- Für Joins (z.B. user_roles, roles) werden die Join-IDs aus den Parent-Objekten extrahiert.
- Die Join-Filter werden als "in"-Filter an die FilterEngine übergeben.
- Die FilterEngine prüft, ob ein RAM-Index für das Join-Feld existiert und öffnet nur die passenden Einträge.
- Falls kein Index existiert, wird ein vollständiger Scan durchgeführt.

## 3. Aktuelle Probleme und Optimierungspotenzial
- Indizes werden zwar geladen, aber die Join-Logik nutzt sie nicht immer konsequent. Es werden zu viele Dateien geöffnet.
- Die FilterEngine muss für "in"-Filter explizit den RAM-Index nutzen und nur die passenden IDs öffnen.
- Die Join-Logik muss so angepasst werden, dass sie garantiert nur die Einträge öffnet, die im RAM-Index gefunden wurden.
- Lazy Loading: Indizes sollten nur geladen werden, wenn sie für die aktuelle Query benötigt werden.

## 4. Geplante Verbesserungen
- **Join-Optimierung:** Die Join-Logik wird so erweitert, dass für jeden Join die passenden IDs aus dem RAM-Index geholt werden. Nur diese Einträge werden geöffnet.
- **Index-Initialisierung:** Sicherstellen, dass alle relevanten Join-Felder als Index in der meta.json definiert sind und beim Start geladen werden.
- **Logging:** Log-Ausgaben für genutzte und ungenutzte Indizes, um die Performance zu überwachen.
- **Testabdeckung:** Tests für Queries mit komplexen Joins, um die korrekte Nutzung der Indizes zu garantieren.

## 5. Best Practices für Query-Performance
- Definiere alle Join-Felder als Index in der meta.json der jeweiligen Tabelle.
- Nutze Filter und Joins so, dass möglichst viele Einträge über den RAM-Index gefunden werden.
- Vermeide vollständige Scans, indem du die Indexstruktur regelmäßig prüfst und optimierst.

## 6. Beispiel-Query und Ablauf
```graphql
query {
  users(
    filter: {
      created_at: { gte: "2025-11-01T00:00:00Z", lte: "2025-11-30T23:59:59Z" },
      email: { like: "*@gmail.com" },
      age: { gte: 18 }
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
    limit: 10,
    offset: 0,
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
**Ablauf:**
1. Filter auf users: RAM-Index für created_at, email, age wird genutzt.
2. Join auf user_roles: RAM-Index für user_id wird genutzt, nur passende Einträge werden geöffnet.
3. Join auf roles: RAM-Index für id wird genutzt, nur passende Einträge werden geöffnet.
4. Ergebnis wird paginiert und zurückgegeben.

## 7. Ziel
- Die Query-Engine soll für alle Filter und Joins ausschließlich die RAM-Indizes nutzen und nur die relevanten Einträge öffnen.
- Die Performance und Skalierbarkeit werden dadurch massiv verbessert.
- Die Roadmap wird regelmäßig erweitert, um neue Optimierungen und Best Practices zu dokumentieren.

