**Vorgang:** slice-zitat-korrektur-vergleichseinheit (Review F-3, MEDIUM; F-5, F-6, LOW; Re-Review F-1 bis F-3, LOW; F-5, F-6, INFO)

**Fund:** Die Befehlsform der Referent-Messung stand als `bash`-Block in `ADR-0158` Entscheidung 6. Das Review fand sie fail-open unter `nullglob`/`failglob` (F-3), ohne gedruckte Zeile unter `set -e` (F-5) und nicht roh an abschließenden Leerzeilen (F-6); `ADR-0159` Entscheidung 4 ersetzte sie. Das Re-Review fand an der neuen Form gestapelte `id`s (F-1), `id` mit Attribut oder selbstschließend (F-2) und verschachtelte Fences (F-3), dazu fremde Tag-Paare (F-5) und die Locale (F-6). Drei ADRs in Folge (`ADR-0157`, `ADR-0158`, `ADR-0159`) in zwei Vorgängen; nach der Entscheidung des Auftraggebers ist `ADR-0159` die letzte Folge-ADR zur Befehlsform.

**Form (Ausprägung):** Erstauftreten — ausführbare Befehlsform als Verfahrensregel im Text einer `Accepted`-ADR, ohne Tabellentest. Adresse: Folge-Slice `slice-zitat-vergleich-werkzeug` (`ADR-0159` Re-Evaluierungs-Trigger (a)).

Quelle: `docs/reviews/review-slice-zitat-korrektur-vergleichseinheit.md` (F-3, F-5, F-6) <!-- d-check:status-provenance -->
· `docs/reviews/review-slice-zitat-korrektur-vergleichseinheit-fixrunde.md` (F-1 bis F-3, F-5, F-6). <!-- d-check:status-provenance -->
