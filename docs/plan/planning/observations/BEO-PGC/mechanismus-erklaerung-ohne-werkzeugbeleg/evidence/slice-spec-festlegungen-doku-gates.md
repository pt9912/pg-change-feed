**Vorgang:** slice-spec-festlegungen-doku-gates (Review F-1 bis F-4, HIGH; Verifikation A-1 und
Nachprüfung A-6, je DoD-Abweichung von Liefer-Punkt 1).

**Fund:** Die Festlegung `SPEC-040` (`spec/pflichtenheft.md` §7, `make docs-check`) beschrieb an
sechs Stellen, wie das gepinnte Werkzeug d-check entscheidet, ohne es dort gemessen zu haben. Die
Sätze stammten aus der Gegenprobe der Quellen (Vertrag, Kommentare in `.d-check.yml`, ADRs), nicht
aus einem Lauf:

- F-1: die Closure-Notiz-Regel gelte „jedem geschlossenen Slice-Plan“; das Werkzeug prüft nur
  `done/slice-*.md` (Mutation R1).
- F-2: `ids` verlange den Link auf das Dokument, das das Muster zuordnet; ein Link auf ein
  beliebiges Ziel genügt (R3).
- F-3: `matrix-inactive` melde einen Verweis; das Werkzeug meldet nur einen Link zwischen Dateien
  einer Klasse (R5–R7; entschieden im Architect-Verdikt zu `matrix-inactive`).
- F-4: ein umbrochener Link werde nicht gemeldet; gilt nur für den Linktext. Der Satz war aus einer
  Messung von `slice-077` an einem älteren d-check übernommen, ohne Kennzeichnung (R8–R10).
- A-1: die Ausgabe-Form (Ströme, letzte Zeile, Felder bei `hostpaths`, mehrzeilige YAML-Meldung)
  wich vom Werkzeug ab.
- A-6: die Berichtigung von A-1 legte eine Reihenfolge von Summen- und Befundzeilen über `make`
  fest, gestützt auf fünf Läufe; in 20 Läufen des Verifiers stand die Summe 6-mal dahinter.

Alle sechs sind im Slice an das gemessene Werkzeug gezogen (`7f139836`, `2ac39c31`, `bf8d26c2`,
`928de0ff`); kein Satz verlangte eine Änderung am Werkzeug.

**Form (Ausprägung):** neuer Träger-Typ: die Erklärung steht in einer **Spec-Festlegung**, die
`ADR-0163` Entscheidung 1 über das Werkzeug stellt, nicht in einer Gate-Doku. Die Gegenprobe der
Quellen (`.claude/commands/plan-welle.md` Schritt 6) lief vollständig (76 Zeilen) und fand keinen
der sechs: sie hält die Festlegung gegen die Quellen, und die Quellen waren selbst breiter oder
enger als das Werkzeug. Gefunden haben sie die Mutationen des Reviewers und die getrennten Ströme
des Verifiers. A-6 ist die Spielart „gemessen, aber zu schmal“ im selben Vorgang.

Quelle: `docs/reviews/review-slice-spec-festlegungen-doku-gates.md` (F-1 bis F-4) <!-- d-check:status-provenance -->
· `docs/reviews/verify-slice-spec-festlegungen-doku-gates.md` (A-1) <!-- d-check:status-provenance -->
· `docs/reviews/verify-slice-spec-festlegungen-doku-gates-nachpruefung.md` (A-6). <!-- d-check:status-provenance -->
