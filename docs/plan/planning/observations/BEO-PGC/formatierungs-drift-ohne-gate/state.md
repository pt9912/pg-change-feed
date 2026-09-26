Zustand: **verkörpert (Schritt und Werkzeug), kein Gate** (3×).

Ausgang: kein Gate. Der Fund kam in allen drei Vorgängen von einem Leser (Reviewer, Verifier)
vor dem Merge, Schwere LOW, je eine kleine Zahl Dateien; ein Gate in `make gates` verlangt
eine ADR (`AGENTS.md` §3.6/§4). Die kleinste tragende Form ist ein Schritt im Implementer-Ablauf
mit Reviewer-Probe und ein `make fmt-check` als Werkzeug ohne Gate-Bindung (Vorbild
`make image-stale`, kein Gate, keine ADR).

Verkörpert: `make fmt-check` (Docker-only im gepinnten `TOOLCHAIN_IMAGE`, `--network none`,
`gofmt -l` über alle Go-Dateien des Baums, Exit 0 formatiert · 1 Abweichung · 2 Formatierer-Fehler
oder ohne Go-Datei) und `make test-fmt-check` (Tabellentest) — liegen in
`harness/sensors/fmt-check.md`, `harness/README.md` §Sensors (Werkzeug-Zeilen) und im `Makefile`;
`.claude/commands/implement-slice.md` Schritt 18 (Absatz „Format“ ruft `make fmt-check`, Lauf im
Bericht) und `.harness/skills/reviewer.md` (LOW-Liste: eine Go-Datei, die `make fmt-check` meldet) ·
seit slice-harness-fmt-check. Beleg-Anker: `make fmt-check` endet am Baum mit Exit 0
(gedruckt: „fmt-check: 254 Go-Dateien geprüft, alle formatiert“, Lauf der Closure).

Trigger für die Aufnahme als Gate (dann Architect-Frage mit ADR-Vorschlag, Kenntnis, kein Umfang):
ein weiteres Auftreten, das der Reviewer trotz gelaufenem Schritt 18 findet. Am Stand der Closure
von `slice-harness-fmt-check` ist der Trigger nicht eingetreten: Review und Verifikation dieses
Slice nennen keinen Format-Befund.

Zähler (abgeleitet): **3×** (evidence/slice-backfill-sql-administration.md,
evidence/slice-backfill-e2e.md, evidence/slice-transformationen-antragsweg-usecase.md).
