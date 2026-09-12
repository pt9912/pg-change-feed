Zustand: offen — Ausgang: **weiter offen** → Ursachenbehebung (unskopierte
`DELETE`/`DROP SCHEMA CASCADE` in mehreren `postgresstorage`-Tests durch
gezielte, konsumenten-spezifische Bereinigung ersetzen) berührt mehrere
bestehende Testdateien zugleich; kein Slice dafür existiert.
Zähler (abgeleitet): 1× (evidence/slice-022.md). Nachrichtlich: ein
früheres, andersartiges Workaround für dasselbe Grundmuster
(`sqlviews_test.go`s Dateinamen-Sortierung, `slice-010`-Zeitraum) entstand
vor der Registrierung dieser Beobachtung und trägt keinen eigenen
Zähler-Beitrag (ein Vorgang zählt einmal, Modul 6).
