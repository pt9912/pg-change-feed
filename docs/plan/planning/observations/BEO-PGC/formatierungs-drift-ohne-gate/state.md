Zustand: offen — **Schwelle erreicht** (3×; der Lese-Schritt der Closure von
`welle-transformationen` liest den Eintrag, Verkörperung 3b: Planner → Architect → Planner).
Ausgang-Vorschlag: der Fund ist in allen drei Vorgängen von einem Leser (Reviewer, Verifier)
vor dem Merge gemacht worden; der Diff enthält jeweils eine kleine Zahl Dateien mit Abweichung.
Kandidaten, vom Kleinen zum Großen: (1) ein Schritt im Implementer-Ablauf vor dem Handoff
(`.claude/commands/implement-slice.md`, vor Schritt 20): `gofmt -l` im gepinnten Toolchain-Image
über die Go-Dateien des eigenen Diffs (`git diff --name-only <Parent> -- '*.go'`), Ausgabe im
Bericht; (2) ein diff-skopiertes make-Ziel mit derselben Aufgabe (Werkzeug, kein Gate — der
Bestand hat sechs Abweichungen, gemessen mit `gofmt -l internal tools cmd test` im
gepinnten Toolchain-Image am Stand dieser Closure: `queries/queries.go`,
`mapper/transformation_test.go`, `receive/seam_test.go`, `outbound/log_test.go`,
`retention/service_test.go`, `test/integration/integration_test.go`); (3) ein `gofmt`-Gate über den
ganzen Baum (Bestandsbereinigung, Entscheidung mit ADR, `AGENTS.md` §3.6). Der Planner nennt
(1) als kleinste wirksame Form; die Wahl ist ein Verdikt des Architects. Ein Textwerkzeug am
Repo (`sed -i`) ist die Ursache des dritten Belegs und bleibt eine Prozessregel des
Auftraggebers, kein Repo-Träger.
Zähler (abgeleitet): **3×** (evidence/slice-backfill-sql-administration.md,
evidence/slice-backfill-e2e.md, evidence/slice-transformationen-antragsweg-usecase.md).
