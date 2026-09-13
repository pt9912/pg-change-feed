# BEO-PGC/commit-traceability-kein-vorab-hook

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Harness-Infrastruktur, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: `make commit-traceability` ([`ADR-0045`](../../../../adr/0045-commit-traceability-standing-gate.md))
prüft die Commit-Traceability-Regeln (mindestens eine `LH-*`-/`ADR-*`-
Kennung im Betreff, keine `SPEC-*`/`ARC-*`-Struktur-ID im Betreff) nur
**nachträglich** — als Teil von `make gates`, manuell aufgerufen, nach dem
Commit. Es existiert kein lokaler `commit-msg`-Git-Hook, der einen
Verstoß bereits beim `git commit`-Aufruf zurückweist. Ein Verstoß wird
deshalb erst sichtbar, nachdem der Commit bereits erstellt (und in dieser
Session mehrfach bereits gepusht) wurde — die einzige verbleibende
Korrektur ist ein `git commit --amend`/History-Rewrite, der bei bereits
gepushten Commits potenziell einen Force-Push und damit eine
Nutzer-Bestätigung braucht (in dieser Session zusätzlich vom
Auto-Mode-Klassifikator als „Git Destructive" blockiert, sobald mehr als
ein einzelner `--amend` auf `HEAD` nötig ist).

Deklaration: direkte Nutzerfrage nach einem zweiten Vorfall dieser Klasse,
2026-09-13 — kein Slice hat diese Beobachtung ausgelöst.
