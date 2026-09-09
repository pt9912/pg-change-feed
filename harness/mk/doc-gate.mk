# harness/mk/doc-gate.mk — Doc-Gate-Fragment, emittiert von ai-harness-init (slice-034).
# Bindet das tool-generierte d-check.mk ein (Befund-Gate docs-check) und haengt
# docs-check an GATE_CHECKS an; der Root-Aggregator faehrt es via make gates.
include d-check.mk
GATE_CHECKS += docs-check
