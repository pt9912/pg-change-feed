# harness/mk/doc-gate.mk — Doc-Gate-Fragment, emittiert von ai-harness-init (slice-034).
# Bindet das tool-generierte d-check.mk ein (Befund-Gate docs-check) und haengt
# docs-check an GATE_CHECKS an; der Root-Aggregator faehrt es via make gates.
include d-check.mk
GATE_CHECKS += docs-check

# Commit-Traceability als Standing-Gate (ADR-0045): die Traceability-Regel
# (harness/README.md §Traceability rules) ist mechanisch getragen — positive
# Haelfte d-check Modul commits, Grenz-Haelfte Shell-Sensor, beide im Target
# commit-traceability (d-check.mk) über die letzten fünf Commits. Die
# Bindung inkrafttet erst, nachdem die Verkörperungs-Commits das
# 5-Commit-Fenster kennungstragend gemacht haben.
GATE_CHECKS += commit-traceability
