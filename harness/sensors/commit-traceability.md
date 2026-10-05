# `make commit-traceability` — Traceability-Kennung je Commit-Message

Ausführliche Fassung der Index-Zeile aus [`harness/README.md` §Sensors](../README.md#sensors-feedback-gates); die Zeile dort trägt einen Satz und verlinkt hierher.

## Vertrag

Commit-Message-Traceability: je Message der Range ≥ 1 `LH-*`-/`ADR-*`-Kennung (d-check Modul `commits`, Befund `commit-untraceable`) und keine `SPEC-*`/`ARC-*`-Kennung im Betreff (`tools/harness/commit-traceability.sh`); Standing-Gate über die letzten 5 Commits (`RANGE=base..head` überschreibt)

## Bindung

[`ADR-0045`](../../docs/plan/adr/0045-commit-traceability-standing-gate.md) · seit slice-006
