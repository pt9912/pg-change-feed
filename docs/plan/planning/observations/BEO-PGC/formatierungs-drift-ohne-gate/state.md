Zustand: **verkörpert (Schritt) — Werkzeug offen, Adresse `slice-harness-fmt-check`** (3×).

Ausgang: kein Gate. Der Fund kam in allen drei Vorgängen von einem Leser (Reviewer, Verifier)
vor dem Merge, Schwere LOW, je eine kleine Zahl Dateien; ein Gate in `make gates` verlangt
eine ADR (`AGENTS.md` §3.6/§4) und eine Bestandsbereinigung vor dem ersten grünen Lauf. Die
kleinste tragende Form ist ein Schritt im Implementer-Ablauf mit Reviewer-Probe, dann ein
netzloses `make fmt-check` als Werkzeug ohne Gate-Bindung (Vorbild `make image-stale`,
kein Gate, keine ADR).

Verkörpert: `.claude/commands/implement-slice.md` Schritt 18 (Absatz „Format“, `gofmt -l` im
gepinnten Toolchain-Image über die Go-Dateien des eigenen Diffs, Lauf im Bericht) und
`.harness/skills/reviewer.md` (LOW-Liste: eine Go-Datei des Diffs, die `gofmt -l` meldet) ·
seit welle-transformationen. Beleg-Anker: `git ls-files '*.go' | xargs -r docker run --rm
--network none -v "$PWD":/src:ro -w /src <TOOLCHAIN_IMAGE> gofmt -l` meldet am Stand des
Verdikts sechs Bestandsdateien (`queries/queries.go`, `mapper/transformation_test.go`,
`receive/seam_test.go`, `outbound/log_test.go`, `retention/service_test.go`,
`test/integration/integration_test.go`).

Folgearbeit (Adresse `slice-harness-fmt-check`, Werkzeug): `make fmt-check` (Docker-only im
`TOOLCHAIN_IMAGE`, `--network none`, `gofmt -l` über alle Go-Dateien des Baums, Exit ≠ 0 bei
Abweichung, nicht in `make gates`), Sensor-Vertrag `harness/sensors/fmt-check.md`, Zeile in
`harness/README.md` §Werkzeuge, Tabellentest gegen ein Wegwerf-Verzeichnis, Formatierung der
sechs Bestandsdateien (ein reiner Format-Commit, Ausgabe von `gofmt -d` als Vorlage, kein
Textwerkzeug), und Umschreiben des Absatzes „Format“ in Schritt 18 auf `make fmt-check`.
Trigger für die Aufnahme als Gate (dann ADR): ein weiteres Auftreten, das der Reviewer trotz
gelaufenem Schritt findet.

Zähler (abgeleitet): **3×** (evidence/slice-backfill-sql-administration.md,
evidence/slice-backfill-e2e.md, evidence/slice-transformationen-antragsweg-usecase.md).
