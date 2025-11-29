# 5. Felder und unterstützte Datentypen

## Ziel
Dieser Plan beschreibt die Umsetzung von Punkt 5 der Roadmap: Zentrale Validierung von Einträgen (Dokumenten) anhand der in `meta.json` definierten Felder, Datentypen und Constraints. Er dient als technische Spezifikation und To-Do-Liste für die Implementierung.

---

## Implementierungsplan für src/table/fields

### Aufgabenübersicht mit Checklisten

1. [ ] **1. Definition der Datenstrukturen**
    - [ ] Erweiterung/Definition von `FieldMeta` und ggf. `TableMeta` für alle Typen und Constraints
    - [ ] Dokumentation der unterstützten Typen und Constraints

2. [ ] **2. Zentrale Validierungsfunktion**
    - [ ] Funktion `ValidateEntry(entry map[string]interface{}, meta TableMeta) error`
    - [ ] Prüfung auf Pflichtfelder (`required`)
    - [ ] Typprüfung (`string`, `int`, `float`, `bool`, `object`, `array`, `date`)
    - [ ] Prüfung aller Constraints (`minLength`, `maxLength`, `min`, `max`, `pattern`, `enum`, `nullable`, `default`)
    - [ ] Prüfung auf zusätzliche Felder (`allowAdditionalFields`)
    - [ ] Fehlercodes/-nachrichten

3. [ ] **3. Einzelne Typ- und Constraint-Validatoren**
    - [ ] String-Validator (inkl. minLength, maxLength, pattern, enum)
    - [ ] Int-/Float-Validator (inkl. min, max, enum)
    - [ ] Bool-Validator
    - [ ] Object-/Array-Validator (rekursiv, falls nötig)
    - [ ] Date-Validator (ISO 8601)
    - [ ] Nullable- und Default-Handling

4. [ ] **4. Fehlerbehandlung und Logging**
    - [ ] Präzise Fehlercodes/-nachrichten
    - [ ] Logging von Validierungsfehlern

5. [ ] **5. Integration in Einfüge-/Update-Logik**
    - [ ] Validierung wird beim Einfügen/Aktualisieren von Dokumenten aufgerufen
    - [ ] Fehlerhafte Einträge werden abgelehnt

6. [ ] **6. Tests**
    - [ ] Unit-Tests für alle Typen und Constraints (inkl. Grenzfälle)
    - [ ] Integrationstests für Einfügen/Aktualisieren/Validieren
    - [ ] Tests für Fehlerfälle, Defaultwerte, optionale Felder

7. [ ] **7. Dokumentation und Beispiele**
    - [ ] Dokumentation der Validierungslogik und Beispiele für meta.json

8. [ ] **8. (Optional) Schema-Migration und Versionierung**
    - [ ] Unterstützung für schemaVersion und Migration bestehender Einträge

---

## Hinweise
- Die Checklisten werden nach erfolgreichem Testen gemeinsam abgehakt.
- Die Validierungslogik soll modular und erweiterbar sein.
