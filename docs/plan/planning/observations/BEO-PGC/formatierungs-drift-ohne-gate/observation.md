# BEO-PGC/formatierungs-drift-ohne-gate

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Go-Formatierung der
Quelltexte, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Das Repo hat kein `gofmt`-Gate (kein Linter, `AGENTS.md` §3.2). Ein Diff
kann Dateien hinterlassen, die `gofmt -l` im gepinnten Toolchain-Image meldet; kein Lauf
von `make gates` färbt sich rot. Belegt: am Parent `2d47d8a7` meldete `gofmt -l internal
tools cmd` zwei Dateien (`seam_test.go`, `log_test.go`, beide Bestand), am Diff-Stand des
Slice sechs — die vier zusätzlichen sind Dateien des Diffs (Einrückung einer
Feldausrichtung in `wiring.go`, drei einzeilige `t.Cleanup(func() { … })` über der
Zeilenlänge, die `gofmt` umbricht). Die zwei Bestandsdateien stehen weiter in der Liste.

**Warum das zählt:** Ohne Gate hält allein der Reviewer die Formatierung; jeder Diff, der
sie nicht prüft, bringt neue Abweichungen und der Bestand wächst.

Deklaration: `slice-backfill-sql-administration`, Review F-8 (LOW).
