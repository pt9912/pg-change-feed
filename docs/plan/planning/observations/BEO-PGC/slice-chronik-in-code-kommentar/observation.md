# BEO-PGC/slice-chronik-in-code-kommentar

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft Kommentar-
Disziplin im gesamten Quellcode, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Go-Quellcode-Kommentare (und, historisch, SQL-/
`schema.yaml`-Kommentare) tragen wiederholt Slice-Chronik — entweder eine
explizite Slice-Kennung (`slice-<NNN>`) als Begründungsquelle einer
Code-Entscheidung, oder implizite Chronik-Sprache (Vorher/Nachher-
Vergleiche wie „zeigt identisches Verhalten wie vor diesem Slice",
„nur noch", „jetzt") — statt ausschließlich den Ist-Zustand zu
beschreiben. Das ist bereits als `AGENTS.md` §3.7 Hard Rule kodifiziert
(„Ein Kommentar beschreibt, was da ist"), wird aber nicht mechanisch
geprüft — kein Gate/Sensor greift dagegen, nur menschliches/Reviewer-Urteil.

## Historie vor der Registrierung — benannt, nicht gezählt

`slice-018` (2026-09-12, welle-5, Commit-Zeitstempel-Arbeit): Nutzer-
Korrektur außerhalb des Review-Prozesses — Kommentare in `queries.go`,
`mapper.go` (`postgresstorage`) und `tools/schema/schema.yaml` trugen
„seit slice-018"/„nur noch". Vor dieser Registrierung, deshalb nicht im
Zähler.

Deklaration: Review zu `slice-041` (F-1, F-4),
Review-Report zur Fixrunde von `slice-041` (drittes Vorkommen).
