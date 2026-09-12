# BEO-PGC/dod-checkbox-nachzug

**Sub-Area:** Planning-Harness (Slice-Lifecycle; Sub-Area-Kürzel `PGC` aus
der Modus-Deklaration)

Die Beobachtung: Der Implementer hakt DoD-Checkboxen in §2 des Slice-Plans
nicht selbst ab, obwohl die Arbeit real erledigt ist — der Verifier findet
regelmäßig materiell abgeschlossene Punkte, deren Checkbox weiterhin auf
`[ ]` steht. Kein Sensor prüft DoD-Checkboxen gegen den tatsächlichen
Zustand (verifizierbar: nein) — die Korrektur bleibt bislang jedes Mal
manuelle Planner-Nacharbeit vor dem `git mv` nach `done/`.

## Benannt, nicht gezählt

Rückblickend erkannt bei der Anlage dieses Eintrags: Das Muster trat
bereits bei slice-015 und slice-016 auf, wurde dort aber nur direkt
behoben, nicht als Beobachtung geführt — der Register-Pfad existierte
für diese Klasse noch nicht. Beide Vorkommen sind nachträglich als
`evidence/`-Dateien erfasst (Modul 6 erlaubt das: der Zähler ist die Zahl
der Dateien, unabhängig davon, wann sie geschrieben wurden).
