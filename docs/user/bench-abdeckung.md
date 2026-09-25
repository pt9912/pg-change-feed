# Bench-Abdeckung je Lastenheft-Kennung

Erzeugt von den drei `tools/bench-*.sh`-Skripten mit Pass/Fail-Schwelle
(`make bench`,
[`ADR-0104`](../plan/adr/0104-benchmark-schwellen-per-001-002-003.md));
`tools/bench-backfill.sh` misst ohne Schwelle und trägt keine Zeile.
Jede Zeile bindet eine Kennung an ihre real durchgesetzte Pass/Fail-
Schwelle. Diese Datei ist eine **stabile Abdeckungs-Deklaration**, kein
Lauf-Beleg — der zuletzt gemessene Wert steht in der stdout-Ausgabe des
jeweiligen Laufs, nicht hier.

| Lastenheft-Kennung | Kurzbeschreibung | Schwelle (SPEC) | Ort |
| --- | --- | --- | --- |
| [`LH-QA-PER-001`](../../spec/lastenheft.md) | Schreibdurchsatz/-latenz derselben Insert-Last mit/ohne aktivierter CDC im Vergleich | Overhead (Median von 5 Läufen) ≤ 35% (`SPEC-025`) | `tools/bench-source-impact.sh` |
| [`LH-QA-PER-002`](../../spec/lastenheft.md) | Skalierbarkeit über die drei [`SPEC-014`](../../spec/pflichtenheft.md)-Lastenstufen | cdc_capture_lag ≤ 60s bei jeder Stufe ([`SPEC-013`](../../spec/pflichtenheft.md), wiederverwendet) | `tools/bench-scaling.sh` |
| [`LH-QA-PER-003`](../../spec/lastenheft.md) | Lesevorgang über größere Change-Mengen im Batch gegen Einzelabruf | Batch-Vorteil ≥ 10× (`SPEC-025`) | `tools/bench-batch-vs-single.sh` |
