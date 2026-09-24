**Vorgang:** slice-backfill-change-origin (Review-Fund F-5, Verifikation V-3, Planner-Closure)

**Fund:** `ADR-0114` Entscheidung 3 sagt, die Rechte der neu angelegten View setze `nacharbeit-roles.sql` „im selben Lauf"; Entscheidung 7 belegt das für die Rolle `cdc_reader`. `tools/schema/nacharbeit-roles.sql` (Zeile 105) vergibt auf `cdc.changes` nur dieses `SELECT`. Der Reviewer fand, dass Handbuch und Makefile-Kommentar nur `cdc_reader` nennen und den Verlust anderer Rechte verschweigen (F-5, LOW); die Fixrunde 2 zog Handbuch (Version 1.46, Punkt „Eigene Rechte") und Makefile-Kommentar nach und meldete die Lücke der ADR an den Architect, ohne die ADR zu ändern. Der Verifier ordnete sie als ADR-Lücke ein, nicht als Slice-Befund: der Slice setzt den Wortlaut um (V-3, INFO). Nach dem Doku-/Kommentar-Zug `c54f873a` trägt das Target-Dokument `harness/targets/schema-rollout.md` die Grenze (§Grenze Punkt 3), der Makefile-Kommentar verweist dorthin.

Quelle: `docs/reviews/review-slice-backfill-change-origin.md` (F-5) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-backfill-change-origin.md` (V-3, §4). <!-- d-check:status-provenance -->
