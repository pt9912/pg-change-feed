Zustand: offen (**2×**) — unter der Schwelle, kein Ausgang zugewiesen. Die vier Dateien
des Diffs von `slice-backfill-sql-administration` sind in der Fixrunde formatiert
(Verifikation: `gofmt -l internal tools cmd` im gepinnten Toolchain-Image listet nur
`receive/seam_test.go` und `outbound/log_test.go`, beide Bestand, beide nicht im Diff).
Der zweite Beleg (`slice-backfill-e2e`, F-10) meldet eine dritte Bestandsdatei außerhalb des
Diffs (`test/integration/integration_test.go`); alle Dateien des Diffs sind formatiert. Ein
Träger ist nicht vorgeschlagen; ein `gofmt`-Gate wäre eine Entscheidung mit ADR, die über
diesen Eintrag hinausgeht.

Zähler (abgeleitet): **2×** (evidence/slice-backfill-sql-administration.md,
evidence/slice-backfill-e2e.md).
