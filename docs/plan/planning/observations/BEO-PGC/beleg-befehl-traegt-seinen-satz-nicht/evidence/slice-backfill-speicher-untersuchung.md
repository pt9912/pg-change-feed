**Vorgang:** slice-backfill-speicher-untersuchung (Verifikation V-1, MEDIUM)

**Fund:** Das Suchlauf-Feld der Fixrunde sagte „Nicht gefunden: kein weiterer Träger außerhalb der Suchwurzeln nennt den Höchstwert 1,57“. Der Satz gilt nur für die Suchwurzeln des Befehls (`docs/user`, `harness`, `tools`, die zwei Messberichte, `docs/plan/planning/open`); `git grep -n -E '1,03 bis 1,57'` über den ganzen Baum nennt zusätzlich `docs/plan/adr/0124-…` (Zeilen 68, 191, 199) und das Architect-Verdikt zur Retention (Zeilen 152, 155, 167). Die Form ist der **Befehl**: die Negativaussage reichte über den Raum hinaus, den der Befehl absucht. Gefunden hat es der Verifier durch den Lauf ohne Pathspec.

Quelle: `docs/reviews/verifikation-slice-backfill-speicher-untersuchung.md` (§6 letzte Zeile, §10 V-1). <!-- d-check:status-provenance -->
