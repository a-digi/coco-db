# Datumsbasierte Filter in der Query-API

## Funktionsweise

Felder vom Typ `date` werden als ISO-8601-Strings gespeichert und können in der Query-API mit Bereichsoperatoren durchsucht werden. Die wichtigsten Operatoren sind:

- `eq`: Exakte Übereinstimmung
- `gte`: Größer oder gleich (ab diesem Datum)
- `lte`: Kleiner oder gleich (bis zu diesem Datum)
- `gt`:  Größer als (nach diesem Datum)
- `lt`:  Kleiner als (vor diesem Datum)

Die Werte müssen immer im ISO-8601-Format (z.B. `2025-11-30T00:00:00Z`) angegeben werden.

## Beispiele für Datums-Filter

### 1. Einträge ab einem bestimmten Datum
```json
{
  "filter": { "created_at": { "gte": "2025-11-01T00:00:00Z" } }
}
```
→ Findet alle Einträge, die am oder nach dem 1.11.2025 erstellt wurden.

### 2. Einträge bis zu einem bestimmten Datum
```json
{
  "filter": { "created_at": { "lte": "2025-11-30T23:59:59Z" } }
}
```
→ Findet alle Einträge, die am oder vor dem 30.11.2025 erstellt wurden.

### 3. Einträge in einem Zeitraum
```json
{
  "filter": { "created_at": { "gte": "2025-11-01T00:00:00Z", "lte": "2025-11-30T23:59:59Z" } }
}
```
→ Findet alle Einträge, die im November 2025 erstellt wurden.

### 4. Exaktes Datum
```json
{
  "filter": { "created_at": { "eq": "2025-11-30T00:00:00Z" } }
}
```
→ Findet alle Einträge, die exakt am 30.11.2025 um 00:00:00 UTC erstellt wurden.

### 5. Vor/nach einem Datum
```json
{
  "filter": { "created_at": { "lt": "2025-12-01T00:00:00Z" } }
}
```
→ Findet alle Einträge, die vor dem 1.12.2025 erstellt wurden.

## Hinweise
- Die Operatoren können kombiniert werden (z.B. `gte` und `lte` für einen Zeitraum).
- Das Datumsformat muss immer RFC3339/ISO-8601 sein (z.B. `2025-11-30T00:00:00Z`).
- Wenn `lte` kleiner als `gte` ist, gibt die API einen Fehler zurück.
- Die Filter sind immer inklusiv (einschließlich der Grenzwerte bei `gte`/`lte`).

## Zusammenfassung
- Für Datumsbereiche immer beide Operatoren angeben:
  - Beispiel: `{ "created_at": { "gte": "2025-11-01T00:00:00Z", "lte": "2025-11-30T23:59:59Z" } }`
- Für exakte Suche nur `eq` verwenden.
- Für offene Bereiche nur `gte` oder nur `lte` verwenden.

