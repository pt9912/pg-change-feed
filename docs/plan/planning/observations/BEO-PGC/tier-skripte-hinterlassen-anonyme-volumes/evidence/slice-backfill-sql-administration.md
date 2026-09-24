**Vorgang:** slice-backfill-sql-administration (Verifikation V-5, Implementer-Bericht der Fixrunde)

**Fund:** Die Tier-Läufe (`make test-store`, `make test-replication`) hinterlassen je Lauf anonyme Volumes des PostgreSQL-Testcontainers (`docker rm -f` ohne `-v`; das Image deklariert `VOLUME /var/lib/postgresql`). Der Verifier zählte sechs anonyme Volumes, die seine Tier-Läufe erzeugten (Erzeugungszeit im Lauf-Fenster), und entfernte sie je einzeln per `docker volume rm` (dangling vor dem Lauf 34, nach der Entfernung 34; gemessen). Der Implementer der Fixrunde zählte 34 → 41 während seiner Läufe (**übernommen**, Bericht im Repo nicht abgelegt). Kein Slice-Befund, Eigenschaft der Tier-Pfade.

Quelle: `docs/reviews/verifikation-slice-backfill-sql-administration.md` (§1 Absatz „Tier-Hygiene", V-5). <!-- d-check:status-provenance -->
