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

# Kennungsfreie Nutzerdokumentation (ADR-0143): ein grep-Wächter über die
# benannte Liste der Nutzerdokumente unter docs/user/ (Handbuch, Standard,
# version.md) auf interne Kennungen und Links nach docs/plan/ bzw.
# docs/reviews/; eine unklassifizierte docs/user/*.md oder eine genannte,
# fehlende Datei endet mit Exit 2. Netzlos
# (tools/harness/handbuch-public-doc-check.sh). Der Tabellentest zur
# Wächter-Logik bleibt Werkzeug, kein Gate (ADR-0143 Festlegung 5).
.PHONY: handbuch-public-doc-check
handbuch-public-doc-check: ## Gate: keine interne Kennung in den Nutzerdokumenten unter docs/user/ (netzlos, grep; ADR-0143)
	@bash tools/harness/handbuch-public-doc-check.sh

GATE_CHECKS += handbuch-public-doc-check

.PHONY: test-handbuch-public-doc-check
test-handbuch-public-doc-check: ## Tabellentest gegen tools/harness/handbuch-public-doc-check.sh (netzlos)
	@bash tools/harness/run-handbuch-public-doc-check-tests.sh

# Kennungsfreie Programm-Ausgaben (ADR-0144 Festlegung 10): ein grep-Wächter
# über die Ausgabe-Literale der Produktions-Go-Dateien (internal/, cmd/,
# tools/schema/) und die echo-/printf-Zeilen der Skripte unter tools/schema/
# und examples/ auf interne Kennungen; Lesefehler und ein leerer Gegenstand
# enden mit Exit 2. Netzlos (tools/harness/ausgabe-kennungen-check.sh). Der
# Tabellentest zur Wächter-Logik bleibt Werkzeug, kein Gate.
.PHONY: ausgabe-kennungen-check
ausgabe-kennungen-check: ## Gate: keine interne Kennung in Ausgabe-Literalen des Servers und seiner Betriebs-Skripte (netzlos, grep; ADR-0144)
	@bash tools/harness/ausgabe-kennungen-check.sh

GATE_CHECKS += ausgabe-kennungen-check

.PHONY: test-ausgabe-kennungen-check
test-ausgabe-kennungen-check: ## Tabellentest gegen tools/harness/ausgabe-kennungen-check.sh (netzlos)
	@bash tools/harness/run-ausgabe-kennungen-check-tests.sh

# Meldungscodes (ADR-0144 Festlegung 5): ein grep-Wächter gleicht die Codes
# des Quelltexts (Produktions-Go, Skripte unter tools/schema/ und examples/),
# der Code-Tabelle (internal/domain/messagecode/codes.go) und des Katalogs im
# Benutzerhandbuch ab; Lesefehler, ein leerer Gegenstand und eine Tabelle ohne
# Code enden mit Exit 2. Netzlos (tools/harness/meldungscodes-check.sh). Der
# Tabellentest zur Wächter-Logik bleibt Werkzeug, kein Gate.
.PHONY: meldungscodes-check
meldungscodes-check: ## Gate: Meldungscodes in Quelltext, Code-Tabelle und Handbuch-Katalog gleich (netzlos, grep; ADR-0144)
	@bash tools/harness/meldungscodes-check.sh

GATE_CHECKS += meldungscodes-check

.PHONY: test-meldungscodes-check
test-meldungscodes-check: ## Tabellentest gegen tools/harness/meldungscodes-check.sh (netzlos)
	@bash tools/harness/run-meldungscodes-check-tests.sh
