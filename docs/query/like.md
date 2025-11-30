# LIKE-Filter in der Query-API

## Funktionsweise

Der LIKE-Operator ermöglicht das Filtern von Feldern nach Mustern mit Wildcards. Er unterstützt folgende Platzhalter:

- `*` steht für beliebig viele Zeichen (wie `.*` in regulären Ausdrücken)
- `?` steht für genau ein Zeichen (wie `.` in regulären Ausdrücken)

Ohne Wildcards wird eine exakte Übereinstimmung verlangt.

## Beispiele für LIKE-Filter

### 1. Teilstring am Anfang
```json
{
  "filter": { "name": { "like": "Max*" } }
}
```
→ Findet alle Namen, die mit "Max" beginnen (z.B. "Maximilian", "Maxine").

### 2. Teilstring am Ende
```json
{
  "filter": { "email": { "like": "*@gmail.com" } }
}
```
→ Findet alle E-Mails, die auf "@gmail.com" enden.

### 3. Teilstring irgendwo im Wert
```json
{
  "filter": { "name": { "like": "*ax*" } }
}
```
→ Findet alle Namen, die "ax" irgendwo enthalten (z.B. "Max", "Alexander").

### 4. Exakter Vergleich (ohne Wildcards)
```json
{
  "filter": { "name": { "like": "Max" } }
}
```
→ Findet nur Einträge, bei denen das Feld exakt "Max" ist.

### 5. Einzelnes Zeichen als Platzhalter
```json
{
  "filter": { "name": { "like": "M?x" } }
}
```
→ Findet z.B. "Max", "Mix", aber nicht "Maxim".

## Hinweise
- Für Teilstring-Suche IMMER `*` als Wildcard verwenden.
- Der LIKE-Operator ist case-sensitive (Groß-/Kleinschreibung wird beachtet).
- Für unscharfe oder Fulltext-Suche siehe auch die Operatoren `partial` und `fulltext`.

## Zusammenfassung
- Verwende `*` und `?` für flexible Muster.
- Ohne Wildcards ist der Vergleich exakt.
- Beispiel für eine Query mit LIKE:

```json
{
  "filter": { "email": { "like": "*example.com" } },
  "limit": 10
}
```

