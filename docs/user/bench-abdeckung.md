# Bench-Abdeckung je Lastenheft-Kennung

Erzeugt von den drei `tools/bench-*.sh`-Skripten mit Pass/Fail-Schwelle
(`make bench`; Schwellen für `LH-QA-PER-002`/`LH-QA-PER-003`:
[`ADR-0104`](../plan/adr/0104-benchmark-schwellen-per-001-002-003.md),
für `LH-QA-PER-001`:
[`ADR-0155`](../plan/adr/0155-per-001-zusatzlatenz-relativ-zur-festschreib-latenz.md));
`tools/bench-backfill.sh` misst ohne Schwelle und trägt keine Zeile.
Jede Zeile bindet eine Kennung an ihre real durchgesetzte Pass/Fail-
Schwelle. Diese Datei ist eine **stabile Abdeckungs-Deklaration**, kein
Lauf-Beleg — der zuletzt gemessene Wert steht in der stdout-Ausgabe des
jeweiligen Laufs, nicht hier.

| Lastenheft-Kennung | Kurzbeschreibung | Schwelle (SPEC) | Ort |
| --- | --- | --- | --- |
| [`LH-QA-PER-001`](../../spec/lastenheft.md) | Zusatzlatenz je Schreibtransaktion derselben Insert-Last mit/ohne aktivierter CDC, bezogen auf die gemessene Festschreib-Latenz | Zusatzlatenz je Commit (Median von 5 Läufen) ≤ max(0.10 ms; 1.5 × Festschreib-Latenz) (`SPEC-025`) | `tools/bench-source-impact.sh` |
| [`LH-QA-PER-002`](../../spec/lastenheft.md) | Skalierbarkeit über die drei [`SPEC-014`](../../spec/pflichtenheft.md)-Lastenstufen | cdc_capture_lag ≤ 60s bei jeder Stufe ([`SPEC-013`](../../spec/pflichtenheft.md), wiederverwendet) | `tools/bench-scaling.sh` |
| [`LH-QA-PER-003`](../../spec/lastenheft.md) | Lesevorgang über größere Change-Mengen im Batch gegen Einzelabruf | Batch-Vorteil ≥ 10× (`SPEC-025`) | `tools/bench-batch-vs-single.sh` |
