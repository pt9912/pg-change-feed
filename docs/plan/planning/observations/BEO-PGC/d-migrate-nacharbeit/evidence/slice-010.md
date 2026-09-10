# Beleg: slice-010

Vorgang: slice-010 (Lesen-Vollabdeckung und SQL-Schnittstelle, Welle 3).

Fund: dieselbe `raw-sql-text-drift`-Grenze wie bei `chk_change_operation`
(slice-006) trifft die drei SQL-Views für SST-002 — PostgreSQL formatiert
eine gespeicherte View über `pg_get_viewdef` um (Alias-Präfixe fallen,
Einrückung ändert sich), der textliche Vergleich von d-migrate 1.2.0
liest das als Drift (`Post-execute compare detected drift`, Exit 5).
Reproduziert mit einer minimalen Einzel-View gegen einen frischen
Postgres-Container. Gelöst mit derselben Ausweichform: berichtete
manuelle Nacharbeit (`tools/schema/nacharbeit-views.sql`, zweiter
psql-Schritt im `schema-rollout`-Target).

Quelle: Implementer-Bericht (slice-010), `tools/schema/nacharbeit-views.sql`.
