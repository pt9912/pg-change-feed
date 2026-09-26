**Vorgang:** slice-transformationen-backfill-pfad (Review F-8, INFO; Verifikation §7 „Kenntnis (Prozess)“)

**Fund:** Alle drei Rollen des Vorgangs machten denselben Fehlgriff, jeweils ohne Wirkung auf das
Repo: der Implementer meldete ein `sed -i` auf `/dev/null`, der Reviewer beim Aufbau seiner
Mutations-Skripte ein `sed -i` mit der Fehlermeldung „kann nicht bearbeitet werden“, der Verifier
eines auf einer Scratch-Skriptdatei. Die Mutationen liefen danach über Python-Ersetzungen auf
einer Kopie im Scratchpad; der Diff trägt keine Spur eines Textwerkzeugs.

**Form (Ausprägung):** das Auftreten in **drei Rollen desselben Vorgangs** zeigt, dass die Regel
nicht am Wissen einer Rolle hängt, sondern am fehlenden committeten Träger; der Fehlgriff blieb
diesmal folgenlos, weil das Ziel keine Repo-Datei war.

Quelle: `docs/reviews/review-slice-transformationen-backfill-pfad.md` (F-8) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-transformationen-backfill-pfad.md` (§5 Zeile F-8, §7). <!-- d-check:status-provenance -->
