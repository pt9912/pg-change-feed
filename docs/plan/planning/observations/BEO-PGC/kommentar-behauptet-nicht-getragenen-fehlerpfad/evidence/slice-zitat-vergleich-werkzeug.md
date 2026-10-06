**Vorgang:** slice-zitat-vergleich-werkzeug (Review F-1, HIGH; Re-Review zu Fixrunde 2 F-1, HIGH; elfte Datei des Eintrags)

**Fund:** Der Kopfkommentar von `tools/harness/zitat-vergleich.sh` sagte zu, die Datei definiere mit `source` geladen nur die Funktionen; sie setzte beim Laden `set -uo pipefail` in der Shell des Aufrufers (Review F-1). Der Kommentar von `idform` sagte zu, eine `id` in Inline-Code nicht zu lesen; der Code prüfte nur einen Backtick unmittelbar vor `<a`, und eine `id` in einem Code-Span mit Text davor galt als Anker, mit Fundstelle im Vertrag des Werkzeugs (Re-Review 2 F-1). Beide durch Nachfahren des zugesagten Pfads gefunden, je in einer Fixrunde getragen.

**Form (Ausprägung):** **Zusage ohne Code** (Ausgang bzw. Erkennung) in einem Harness-Skript. Beide vor dem Merge vom Reviewer gefunden; HIGH, daher Datei trotz Deckel. Ausgang unverändert **verkörpert**.

Quelle: `docs/reviews/review-slice-zitat-vergleich-werkzeug.md` (F-1) <!-- d-check:status-provenance -->
· `docs/reviews/review-slice-zitat-vergleich-werkzeug-fixrunde-2.md` (F-1). <!-- d-check:status-provenance -->
