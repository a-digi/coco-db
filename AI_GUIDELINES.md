# AI Guidelines

Diese Richtlinien gelten für alle KI-gestützten und automatisierten Code-Generierungs- und Review-Prozesse in diesem Projekt.

## Grundprinzipien

1. **SOLID-Prinzipien**
    - Jeder generierte oder geprüfte Code muss die SOLID-Prinzipien einhalten:
        - **S**ingle Responsibility Principle (Jede Klasse/Funktion hat genau eine Aufgabe)
        - **O**pen/Closed Principle (Offen für Erweiterung, geschlossen für Modifikation)
        - **L**iskov Substitution Principle (Austauschbarkeit von Subtypen)
        - **I**nterface Segregation Principle (Schnittstellen sind spezifisch und klein)
        - **D**ependency Inversion Principle (Abhängigkeiten zu Abstraktionen, nicht zu konkreten Implementierungen)

2. **DRY-Prinzip (Don't Repeat Yourself)**
    - Redundanzen im Code sind zu vermeiden.
    - Wiederverwendbare Logik wird in Funktionen, Methoden oder Modulen gekapselt.
    - Copy-Paste-Code ist zu vermeiden.

3. **Separation of Concerns**
    - Funktionale Verantwortlichkeiten werden klar getrennt.
    - Jede Komponente, Klasse oder Funktion ist für einen klar abgegrenzten Aspekt zuständig.
    - Keine Vermischung von UI, Geschäftslogik, Datenzugriff oder Infrastruktur-Code.

4. **Dokumentation und Nachvollziehbarkeit**
    - Jeder KI-generierte oder -veränderte Code muss ausreichend dokumentiert werden (Kommentare, Docstrings, Änderungsprotokoll).
    - Die Motivation für komplexe oder ungewöhnliche Lösungen ist zu erläutern.
    - KI-generierte Änderungen sind im Commit-Text oder in Pull-Requests klar zu kennzeichnen.

5. **Testbarkeit und Tests**
    - KI-generierter Code muss testbar sein (z.B. durch Unit- oder Integrationstests).
    - Wo möglich, sollen automatisierte Tests für neue oder geänderte Funktionen mitgeliefert werden.
    - Bestehende Tests dürfen durch KI-Änderungen nicht brechen.

6. **Fehler- und Ausnahmebehandlung**
    - Fehlerfälle und Ausnahmen sind explizit zu behandeln.
    - Es dürfen keine stillen Fehler oder unklare Fehlermeldungen entstehen.

7. **Sicherheit und Datenschutz**
    - KI-generierter Code darf keine sensiblen Daten ungeschützt verarbeiten, speichern oder übertragen.
    - Sicherheitsrelevante Aspekte (z.B. Input-Validierung, Authentifizierung) sind besonders zu beachten.

8. **Konsistenz und Stil**
    - Der generierte Code muss sich an die im Projekt geltenden Code-Style-Guides halten (Formatierung, Benennung, Struktur).
    - Bei Unsicherheiten ist der bestehende Code als Referenz zu nutzen.

9. **Transparenz bei Einschränkungen**
    - Wenn die KI eine gewünschte Änderung nicht direkt umsetzen kann (z.B. wegen fehlender Schreibrechte), muss sie immer den vollständigen, betroffenen Dateiinhalt bereitstellen – niemals nur Code-Snippets.
    - **Wichtig:** In solchen Fällen ist es ausdrücklich untersagt, lediglich Code-Snippets oder Ausschnitte zu liefern. Es muss immer der vollständige, relevante Dateiinhalt bereitgestellt werden.

## Anwendung
- Diese Prinzipien gelten für alle KI-generierten Vorschläge, automatisierte Refactorings und Reviews.
- Bei Zielkonflikten ist stets die bestmögliche Balance zwischen Lesbarkeit, Wartbarkeit und Erweiterbarkeit zu wählen.
- Verstöße gegen diese Prinzipien sind zu dokumentieren und zu begründen.

---

*Letzte Aktualisierung: 25.11.2025*
