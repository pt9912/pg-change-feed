Zustand: **verkörpert** — Ausgang: **verkörpert** → neue Hard Rule
`AGENTS.md` §3.9 („Exit-Code eines Gate-Laufs wird direkt geprüft, nie
durch eine Pipe/einen Wrapper hindurch"), plus Kurzverweis darauf in
`AGENTS.md` §6 Schritt 6 und `harness/README.md` §Minimal agent workflow
Schritt 6 — verkörpert `seit welle-15`. Architect-Verdikt:
`docs/reviews/architect-verdict-pipe-maskiert-make-exit-code.md`
(Diagnose: strukturell anders als die Chronik-/Handbuch-Präzedenzfälle —
der Fehler entsteht in der Shell-Ausführung selbst, hinterlässt kein
zweites, unabhängig einsehbares Artefakt, das eine zweite Rolle prüfen
könnte; daher eine zentrale Hard Rule statt einer zweiten Ebene in
`.harness/skills/reviewer.md` oder den Rollen-Command-Dateien; kein
mechanischer Sensor möglich, da kein Objekt im Working Tree). Wellenloser
Lese-Schritt-Träger entfällt hier: Der Lese-Schritt lief als Teil der
laufenden `welle-15`-Closure (Modul 8 §Rollen-Sequenz für eine Welle,
Schritt 3b).

Zähler (abgeleitet): **4×** (evidence/slice-054.md, evidence/slice-055.md,
evidence/slice-058.md, evidence/slice-078.md) — Schwelle erreicht, Ausgang im
Lese-Schritt der `welle-15`-Closure zugewiesen. Der vierte Beleg liegt **nach**
der Verkörperung und trifft die **geschärfte** Hälfte: nicht die maskierende
Pipe, sondern die fehlende Konditionierung von Gate-Lauf und Folgehandlung über
zwei Schritte. Die Regel steht; ihr Wächter bleibt Disziplin — kein Sensor fängt
sie.
