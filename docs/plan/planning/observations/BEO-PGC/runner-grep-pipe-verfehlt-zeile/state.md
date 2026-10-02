Zustand: **verkörpert** — Anker: `slice-harness-grep-pipe-sigpipe-unter-pipefail`, Träger die
sechs Skripte unter `tools/harness/` am Commit `f8aca691` (35 Fundstellen `docker logs … | grep -q…`
ersetzt durch `… | grep -F … >/dev/null`), Beleg der e2e-Lauf 37036006973 (beide PostgreSQL-Legs
`success`). Das Regelwerk (`AGENTS.md` §3.9) trägt keinen Satz dazu: der Vorschlag steht in der
Closure-Notiz des Slice, die Freigabe liegt beim Auftraggeber (offen).

Zähler (abgeleitet): **1×** (evidence/slice-sdk-sse-filter-phase-verbindung-haertung.md).
Reproduktionsläufe sind keine Vorkommen. Fehlerrate der alten Form in vier Messungen, je 1000
Aufrufe: 10, 9 (Implementer), 16 (Reviewer), 6 (Verifier); neue Form 0.
