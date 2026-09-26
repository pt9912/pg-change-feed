**Vorgang:** slice-harness-suchlauf-nachmessen (Review F-6, INFO; Verifikation V-3, INFO)

**Fund:** Der Plan führte in §6 das Risiko „Docker-only für `git`“ mit der Erwartung, die Skript-Köpfe zeigten Präzedenz. Der Reviewer fand, dass der Vertrag nur `bash` und `git` nannte, das Werkzeug aber `realpath` braucht (Vertrag nachgezogen); Reviewer und Verifier nannten die Vereinbarkeit mit `AGENTS.md` §3.1 als Aussage des Architects, nicht als Befund. Der Planner löste das Risiko in der Closure mit gemessener Präzedenz (42 Dateien unter `tools/` mit `git rev-parse --show-toplevel`), die Regel selbst führt keine Klasse dafür.

Quelle: `docs/reviews/review-slice-harness-suchlauf-nachmessen.md` (F-6) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-harness-suchlauf-nachmessen.md` (§8 V-3). <!-- d-check:status-provenance -->
