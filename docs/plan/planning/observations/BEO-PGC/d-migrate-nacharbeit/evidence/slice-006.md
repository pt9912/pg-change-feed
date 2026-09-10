# Beleg: slice-006

Vorgang: slice-006 (Docker-Compose-Umgebung und MVP-Integrationstest,
in `done/` seit seiner Closure, welle-2).

Fund: die Erstlieferung `tools/schema/schema.yaml` (Überführung aus der
handgeschriebenen DDL) verliert den CHECK `chk_change_operation` an die
d-migrate-1.2.0-Konvergenz-Grenze (`raw-sql-text-drift`, Exit 7 Drift /
E012 Autorenform — beide Formen im Implementer-Lauf rot gesehen); der
Constraint lebt als berichtete Nacharbeit im Rollout-Target.

Quelle: docs/reviews/review-slice-006.md (F-4/F-5), Implementer-Bericht
(slice-006, Zentraler Befund am Erstversatz).
