**Vorgang:** slice-capture-leerlauf-quellbelege (Review F-7, INFO; Verifikation §5 Zeile F-7)

**Fund:** Zwei Slice-Pläne verlinkten einander mit festem Lifecycle-Verzeichnis: der Plan `slice-capture-leerlauf-quellbelege` (in `in-progress/`) den Träger-Slice `slice-wal-fehlerschwelle-ausgangsklasse` über `../open/slice-wal-fehlerschwelle-ausgangsklasse.md` (drei Stellen), dessen Plan (in `open/`) den Quellbeleg-Slice über `../in-progress/slice-capture-leerlauf-quellbelege.md` (eine Stelle). Der jeweils nächste Move (`open` → `next`, `in-progress` → `done`) bricht den Link des anderen Dokuments (`docs-check` Modul `links`, `target-missing`). Der Reviewer nannte es als Träger außerhalb des Code-Diffs (Planner-Commits), der Verifier bestätigte es; die Planner-Closure stellte alle vier Stellen auf die Kennung um, vor dem ersten Move.

**Form (Ausprägung):** dieselbe Klasse mit einem **neuen Träger-Typ**: ein Slice-Plan als Quelle des Links (die zwei `structure`-Regeln in `.d-check.yml` decken `docs/reviews/**` und `observation.md`, nicht Slice-Pläne); gefunden vor dem Move, vom Reviewer, das Gate hätte den Fehler am Move gefangen. Schwere INFO, Ausgang unverändert **verkörpert**.

Quelle: `docs/reviews/review-slice-capture-leerlauf-quellbelege.md` (F-7) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-capture-leerlauf-quellbelege.md` (§5 Zeile F-7). <!-- d-check:status-provenance -->
