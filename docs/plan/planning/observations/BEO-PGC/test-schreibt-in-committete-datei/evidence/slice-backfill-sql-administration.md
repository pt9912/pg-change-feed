**Vorgang:** slice-backfill-sql-administration (Verifikation V-5, Review-Mutationen)

**Fund:** `make test-store` überschreibt über `tools/schema/apply-rollout.sh` die `target`-Zeile der committeten Datei `tools/schema/plan.yaml` (Tier-Ziel `postgres://cdc:***@cdc-store-test-pg:5432/cdc_test?…` statt des committeten Compose-Werts `cdc-test-postgres`); der Verifier nahm sie nach jedem Lauf zurück, der Reviewer nach jeder Mutation und jedem Tier-Lauf (`git checkout`). Gemessen im zweiten Store-Lauf des Verifiers: `git diff` nach dem Lauf zeigt allein die `target`-Zeile; `down.sql` unverändert. Derselbe Pfad wie bei den Vorbelegen, hier in den Läufen zweier weiterer Rollen desselben Vorgangs; die Läufe der zwei DB-Tiers besitzen weiterhin keine Wiederherstellung.

Quelle: `docs/reviews/verifikation-slice-backfill-sql-administration.md` (§3.2, V-5) <!-- d-check:status-provenance -->
· `docs/reviews/review-slice-backfill-sql-administration.md` (Abschnitt „Mutationen", Kopfzeile: „`tools/schema/plan.yaml`/`down.sql` als Tier-Nebeneffekt jeweils zurückgenommen"). <!-- d-check:status-provenance -->
