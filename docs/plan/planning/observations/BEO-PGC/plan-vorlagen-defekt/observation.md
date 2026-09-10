# BEO-PGC/plan-vorlagen-defekt

**Sub-Area:** Planning-Harness (Slice-Pläne beim Füllen; Sub-Area-Kürzel
`PGC` aus der Modus-Deklaration)

Die Beobachtung: Beim Füllen mehrerer Slice-Pläne in einem Zug (Skript
oder wiederholtes Editieren, typisch bei einer Welle-Eröffnung) trägt
§2 (Definition of Done) zwei Defekte, die kein Gate erkennt: eine
doppelte DoD-Zeile (`make gates` grün zweimal) und ein unaufgelöster
Vorlagen-Platzhalter (`<Schnittstelle X>`). Beide Defekte entstehen aus
demselben Fill-Fehler und treten korreliert auf.

Deklaration: `.claude/commands/plan-welle.md` (§2-Form-Prüfung nach dem
Füllen, seit slice-010); kein d-check-Modul erkennt Zeilen-Duplikate
oder Platzhalter-Reste in einer bestimmten Spalte — die Klasse ist
maschinell nicht belegt.
