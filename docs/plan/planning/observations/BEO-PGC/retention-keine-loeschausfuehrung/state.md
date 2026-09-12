Zustand: offen — Ausgang: **weiter offen** → tatsächliche Löschausführung
(Use-Case/CLI/Job, der `RetentionPolicy.AllowsDeletion` real aufruft und
`cdc.change`-Zeilen löscht) und die Metrik `cdc_storage_bytes` bauen, oder
`LH-FA-RET-002`…`006` bewusst als noch nicht umgesetzt markieren, falls das
Ziel aktuell nicht verfolgt wird; kein Slice dafür existiert.
Zähler (abgeleitet): 0× — noch kein abgeschlossener Vorgang trägt einen
Beleg (siehe observation.md, „Benannt, nicht gezählt").
