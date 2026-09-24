**Vorgang:** slice-backfill-run-usecase (Review-Fund F-9)

**Fund:** Der Reviewer prüfte die Klassen-Abbildung des Runs: `internal` als Rückfall ist zulässig (die Menge von `SPEC-008` hat sieben Klassen), `ErrSchemaVersionUnknown` → `configuration` ist eine Setzung ohne ADR-Text, die zur Klasse passt; `schema` vergibt der Run nicht. Ob der Backfill-Pfad der Transformationen sie für „Regel nicht anwendbar" braucht, ist eine offene Frage an den Architect; `slice-transformationen-backfill-pfad` nennt die Nichtanwendbarkeit ohne Klasse und trägt das Kurzverdikt als Start-Trigger.

Quelle: `docs/reviews/review-slice-backfill-run-usecase.md` (F-9) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-backfill-run-usecase.md` (§7, Zeile `ADR-0111` Teilfrage 5). <!-- d-check:status-provenance -->
