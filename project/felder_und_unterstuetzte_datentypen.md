# 5. Felder und unterstützte Datentypen

## Ziel
Dieser Plan beschreibt die Umsetzung von Punkt 5 der Roadmap: Zentrale Validierung von Einträgen (Dokumenten) anhand der in `meta.json` definierten Felder, Datentypen und Constraints. Er dient als technische Spezifikation und To-Do-Liste für die Implementierung.

---

## Implementierungsplan für src/table/fields

### Aufgabenübersicht mit Checklisten

1. [x] **1. Definition der Datenstrukturen**
    - [x] Erweiterung/Definition von `FieldMeta` und ggf. `TableMeta` für alle Typen und Constraints
    - [x] Dokumentation der unterstützten Typen und Constraints

2. [x] **2. Zentrale Validierungsfunktion**
    - [x] Funktion `ValidateEntry(entry map[string]interface{}, meta TableMeta) error`
    - [x] Prüfung auf Pflichtfelder (`required`)
    - [x] Typprüfung (`string`, `integer`, `boolean`, `json`, `date`)
    - [x] Prüfung aller Constraints (`minLength`, `maxLength`, `pattern`, `enum`, `nullable`, `default`, `description`)
    - [x] Prüfung auf zusätzliche Felder (`allowAdditionalFields`)
    - [x] Fehlercodes/-nachrichten

3. [x] **3. Einzelne Typ- und Constraint-Validatoren**
    - [x] String-Validator (inkl. minLength, maxLength, pattern, enum)
    - [x] Integer-Validator (inkl. minLength, maxLength, enum)
    - [x] Boolean-Validator (inkl. enum)
    - [x] JSON-Validator (inkl. maxLength, minLength)
    - [x] Date-Validator (ISO 8601)
    - [x] Nullable- und Default-Handling

4. [ ] **4. Fehlerbehandlung**
    - [x] Präzise Fehlercodes/-nachrichten

6. [ ] **6. Tests**
    - [x] Unit-Tests für alle Typen und Constraints (inkl. Grenzfälle) (nur teilweise, String ist abgedeckt, andere Typen sollten ergänzt werden)
    - [ ] Integrationstests für Einfügen/Aktualisieren/Validieren
    - [x] Tests für Fehlerfälle, Defaultwerte, optionale Felder

7. [x] **7. Dokumentation und Beispiele**
    - [x] Dokumentation der Validierungslogik und Beispiele für meta.json

8. [ ] **8. (Optional) Schema-Migration und Versionierung**
    - [ ] Unterstützung für schemaVersion und Migration bestehender Einträge

---

## Hinweise
- Die Checklisten werden nach erfolgreichem Testen gemeinsam abgehakt.
- Die Validierungslogik soll modular und erweiterbar sein.
