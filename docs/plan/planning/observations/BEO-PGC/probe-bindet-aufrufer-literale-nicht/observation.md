# BEO-PGC/probe-bindet-aufrufer-literale-nicht

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft eine Probe im
Docker-Bau gegen den Aufrufer außerhalb des Bau-Kontexts, keine eigene Sub-Area
im Sinn der Modus-Deklaration).

Die Beobachtung: Eine Probe im Bau bindet eine Zusage an **ein** Literal ihrer
Umgebung (hier: die Namen der Gradle-Aufgaben im Dockerfile), der **Aufrufer** der
Zusage steht dagegen außerhalb des Bau-Kontexts (hier: die zwei Aufruf-Literale im
Workflow unter `.github/`) und ist von der Probe nicht gebunden. Läuft das
Aufruf-Literal auseinander, färbt kein Sensor; der Fehler zeigt sich erst im
realen Lauf. Zusätzlich hängt die Probe an einem Ausgabeformat des geprüften
Werkzeugs (`--dry-run`-Zeile `:<Name> SKIPPED`) und eine ihrer Prüfungen an einem
`grep -F`, das eine Kommentarzeile mit demselben Literal erfüllen würde.

**Abgrenzung zu `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`:** Dort ist
die Kette *Eingabe → Weitergabe → Ablehnung* im Test unterbrochen und wird durch
Mutation der Eingabeseite sichtbar; hier färbt jede Mutation der Probe-Seite
(fünf gefahren) die Probe, und die Lücke liegt in dem, was die Probe **nicht**
erreicht — der Aufrufer außerhalb ihres Kontexts. Bleibt getrennt, bis ein zweiter
Beleg zeigt, ob beide Klassen zusammenfallen.

Deklaration: Planner-Agent (Closure-Rolle), 2026-09-26, auf Grundlage des Reviews
zu `slice-sdk-kotlin-cloudsmith` (F-5, INFO).
