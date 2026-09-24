**Vorgang:** slice-backfill-change-origin (Review-Fund F-6, Verifikation V-6, Planner-Closure)

**Fund:** Der Reviewer ließ den Guard-Test `tools/harness/run-schema-rollout-guard-test.sh` einmal ungekürzt laufen; danach zeigte `git status` `tools/schema/plan.yaml` und `tools/schema/down.sql` modifiziert (F-6, INFO, gemessen). Fixrunde 2 ergänzte im Skript die Sicherung beider Dateien beim Start und die Wiederherstellung in `cleanup`; der Verifier fuhr den Guard-Test danach ungekürzt und fand `git status --short` leer. Derselbe Verifier maß, dass `make test-store` und `make test-replication` die `target`-Zeile von `plan.yaml` überschreiben (V-6, INFO) und nahm die Datei nach jedem Lauf per `git checkout` zurück; der Slice-Commit `5f126971` hatte die Compose-Bezeichnung nachträglich wiederherstellen müssen. Beide Funde gehören zu einem Vorgang.

Quelle: `docs/reviews/review-slice-backfill-change-origin.md` (F-6) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-backfill-change-origin.md` (§1 Nebenbeobachtung, V-6). <!-- d-check:status-provenance -->
