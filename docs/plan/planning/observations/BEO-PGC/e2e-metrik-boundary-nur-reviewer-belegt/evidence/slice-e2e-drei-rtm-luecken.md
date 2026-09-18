**Vorgang:** e2e-drei-rtm-luecken
**Fund:** Review-Fund F-1 (`docs/reviews/review-slice-e2e-drei-rtm-luecken.md` <!-- d-check:status-provenance -->):
die neue Phase „Metriken-Minimum-Beleg" (`LH-QA-OPS-003`) prüft nur den
positiven Fall (`cdc_changes_pending`>0, `cdc_errors_total`>0). Der
Reviewer bestätigte den Randfall (gesunde Quelle ohne Fehlerzeile,
vollständig bestätigter Consumer mit `cdc_changes_pending`=0) real per
eigener, nicht committeter Scratch-PostgreSQL — korrekt, aber kein
Beleg im ausgelieferten Testbestand. Der Verifier bestätigte
unabhängig dieselbe Lücke ohne sie als DoD-Verstoß zu werten (das
Lastenheft verlangt für `LH-QA-OPS-003` keinen expliziten
Negativ-/Boundary-Fall). Erster Beleg dieser Klasse.
