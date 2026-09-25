**Vorgang:** slice-backfill-bench-richtgroesse (Architect-Verdikt Befund 2, Verifikation V-1)

**Fund:** Der DoD-Punkt „Bench“ verlangte ursprünglich ein „reales `make bench`-Lauf mit Exit 0“, der Closure-Trigger dann ein Exit 2 von `make bench`; beides ist eine Eigenschaft des Messhosts, nicht des Slice. Auf dem Messhost endet `tools/bench-source-impact.sh` oberhalb der 35-%-Schwelle (Verdikt: Latenz des Festschreibens, `fdatasync` 2.956 µs je Operation, Host-Last widerlegt; Verifikation §1: 92,0 %), das neue Skript läuft einzeln mit Exit 0. Der Implementer meldete das Kriterium als nicht erreichbar und formulierte den Beleg um; das Verdikt wies die Ratifizierung dem Planner zu, die Verifikation (V-1) fand die Umformulierung in der Fixrunde ohne erkennbare Rolle und den Closure-Trigger weiter mit „Exit 2“. Die Closure fasst beides hostunabhängig (Beleg des Gegenstands einzeln; Ganz-Target-Ergebnis zur Kenntnis mit Host).

Quelle: `docs/reviews/architect-verdict-backfill-wal-rueckstand-und-bench-rot.md` (Befund 2) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-backfill-bench-richtgroesse.md` (V-1, §1). <!-- d-check:status-provenance -->
