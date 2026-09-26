**Vorgang:** slice-transformationen-antragsweg-schema (Plan §3 Zeile `nacharbeit-administration.sql`, Plan §6, `ADR-0125`)

**Fund:** Eine Funktion der Nacharbeit-Datei mit einem `jsonb`-Parameter (`cdc.set_transformation(…, rule_spec jsonb)`) ließ den zweiten `make schema-rollout` mit d-migrate-Exit 5 enden (`function set_transformation(text, text, text, text, json) does not exist`, Wiederherstellung `FULL_ROLLBACK_CONFIRMED`): d-migrate 1.3.1 meldet jede Funktion mit `json`- oder `jsonb`-Parameter im `--plan-only`-Report als `in:json` und rendert ihren Abbau als `DROP FUNCTION … (…, json)`; die Signatur mit `json` existiert bei einem `jsonb`-Parameter nicht. Mit `json` als Parameterart endet der zweite Lauf mit Exit 0 (gemessen im Plan, vom Reviewer mit dem `jsonb`-Zweig nachgestellt: Exit 0 und Exit 2 mit d-migrate-Fehler 5; unmutiert Exit 0 und 0). Entschieden mit `ADR-0125`; die Guard-Schreibweise `in:json` bindet `guard_test.go` mit einem Test, der `in:jsonb` ablehnt.

**Form (Ausprägung):** Unterfall der Objektklasse „SQL-Funktionen/Prozeduren“: die Signatur einer Nacharbeit-Funktion wird von d-migrate für `json`- und `jsonb`-Parameter gleich gemeldet und für den Abbau falsch gerendert; die Ausweichform trägt eine Funktion mit `jsonb`-Parameter nicht.

Quelle: `docs/reviews/review-slice-transformationen-antragsweg-schema.md` (Eigenständig durchgeführte Prüfungen, Begründung von `ADR-0125` nachgestellt) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-transformationen-antragsweg-schema.md` (§7 `ADR-0125`). <!-- d-check:status-provenance -->
