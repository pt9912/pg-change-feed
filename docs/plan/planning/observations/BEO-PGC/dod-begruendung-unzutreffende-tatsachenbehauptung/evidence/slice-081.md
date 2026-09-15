# Beleg: slice-081

Vorgang: `slice-081` — die Executor-Naht und die Rampen-Neu-Bemessung, die sie
ausgelöst hat.

Fund: Ein **DoD-Kriterium** des eigenen Slice trug eine Zahl, die nicht gemessen
war. §2 (Liefer-Punkt 3) sagte „Gegenstand **1679 → 1817** Statements"; 1817 ist
die *abgeleitete* Summe `1679 + 138`, nicht der gemessene Wert — gemessen sind
**1831** (der zweite Implementer-Lauf). Der Fund kam über die Verifikation des
Folge-Gegenstands: die Zahl stand auch in `ADR-0077` §Kontext, wohin sie über den
Auftrag des Planners an den Architect gelangt war, und widersprach dort der
Prozentzeile derselben Tabelle (`1306/1817 = 71,9 %`, gedruckt `71,3 %`).

Quelle: `docs/reviews/review-slice-081.md` (Review) ·
`docs/plan/planning/in-progress/slice-081-executor-naht.md` §1/§2 (Berichtigung) ·
`docs/plan/adr/0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md`
(Entscheidung 4: die Berichtigung reitet in der Folge-ADR, weil `ADR-0073`s
Zitat-Korrektur-Klasse eine inhaltliche Zahl nicht deckt).
