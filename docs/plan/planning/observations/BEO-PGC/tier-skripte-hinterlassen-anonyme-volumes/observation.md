# BEO-PGC/tier-skripte-hinterlassen-anonyme-volumes

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Test-/Tier-Skripte
unter `tools/harness/`, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Die Test- und Tier-Skripte starten PostgreSQL-Testcontainer mit
`docker run -d --name …` und räumen sie mit `docker rm -f <container>` ab, ohne `-v`
(`tools/harness/run-store-tests.sh`, `run-replication-tests.sh`; dasselbe Muster
steht in `run-schema-rollout-guard-test.sh`, dessen Volume-Ertrag nicht gemessen ist).
Das gepinnte Image
`postgres:18-alpine@sha256:63bdc97d…` deklariert ein `VOLUME` (`/var/lib/postgresql`,
gemessen mit `docker image inspect` beim Schreiben dieses Eintrags), das der Container
als anonymes Volume anlegt; `docker rm -f` ohne `-v` entfernt das Volume nicht. Jeder
Tier-Lauf hinterlässt deshalb ein oder mehrere dangling-Volumes mit einem PostgreSQL-Datenverzeichnis.
Belegt: der Verifier zählte sechs anonyme Volumes, die seine Tier-Läufe erzeugten, und
entfernte sie je einzeln (dangling vor dem Lauf 34, nach der Entfernung 34); der Implementer
der Fixrunde zählte 34 → 41 dangling während seiner Läufe (Zahl **übernommen** aus dem
Bericht des Implementers, im Repo nicht abgelegt). Beim Schreiben dieses Eintrags druckt
`docker volume ls -qf dangling=true | wc -l` 34.

**Warum das zählt:** Ein Entwicklerrechner, der Tier-Läufe wiederholt (Review, Verifikation,
Closure), füllt die Platte mit Datenverzeichnissen; die Bereinigung ist Handarbeit
(`docker volume rm` je Volume — ein `prune` trifft fremde Volumes).

Deklaration: `slice-backfill-sql-administration`, Verifikation V-5 (INFO).
