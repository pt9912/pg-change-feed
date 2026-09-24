**Vorgang:** slice-backfill-snapshot-reader (Verifikation V-6)

**Fund:** `make test-replication` (und `make test-store`) überschreiben über `tools/schema/apply-rollout.sh` die `target`-Zeile der committeten Datei `tools/schema/plan.yaml`; der Verifier nahm die Datei nach jedem Tier-Lauf per `git checkout` zurück, der Arbeitsbaum war sonst nach dem Lauf dirty (V-6, INFO). Derselbe Pfad wie beim Erstbeleg, jetzt an zwei Tier-Läufen je PostgreSQL-Version; die Läufe der zwei DB-Tiers besitzen weiterhin keine Wiederherstellung.

Quelle: `docs/reviews/verifikation-slice-backfill-snapshot-reader.md` (§1 Tabelle, V-6). <!-- d-check:status-provenance -->
