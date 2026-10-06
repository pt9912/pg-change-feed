**Vorgang:** slice-zitat-vergleich-werkzeug (Review F-4, MEDIUM; F-5, F-6, LOW; Re-Review zur Fixrunde F-1, MEDIUM; F-3, F-4, LOW; Re-Review zu Fixrunde 2 F-1, HIGH; F-2, F-3, LOW; Re-Review zu Fixrunde 3 F-1, LOW)

**Fund:** `tools/harness/zitat-vergleich.sh` bildet die Markdown-Einheiten von `ADR-0158`/`ADR-0159` nach. Jeder der vier Review-Läufe fand mindestens eine Eingabeform, an der das Werkzeug still `cmp 0` lieferte und die weder Kommentar noch Vertrag als Grenze nannte: `id` in eingerücktem Code, Listen-Fence und HTML-Kommentar (Review F-4); ein Lokator `L7` oder `L5-3`, für den der Vertrag Exit 2 zusagte (F-5); eine schließende `#`-Folge (F-6); Backtick-Parität vor `<!--` (Re-Review F-1), Einzug aus Leerzeichen und Tab (F-3), Kommentar in der `id`-Zeile (F-4); `id` in einem Code-Span mit Text davor, eine Form mit Fundstelle im Vertrag selbst (Re-Review 2 F-1, HIGH), `<!--` in eingerücktem Code nach einer `id`-Zeile (F-2), eingerückter Code im Blockzitat (F-3); ein escapter Backtick vor einem Code-Span (Re-Review 3 F-1). Ab der zweiten Fixrunde gilt im Werkzeug „Mehrdeutigkeit endet mit Exit 2“; die konstruierten Formen ohne Fundstelle stehen als benannte Grenzen im Vertrag, die letzte mit der Closure.

**Form (Ausprägung):** Erstauftreten — nachgebildetes Format, unbenannte Eingabeform, falsches Grün. Ein Vorgang, vier Review-Läufe.

Quelle: `docs/reviews/review-slice-zitat-vergleich-werkzeug.md` (F-4 bis F-6) <!-- d-check:status-provenance -->
· `docs/reviews/review-slice-zitat-vergleich-werkzeug-fixrunde.md` (F-1, F-3, F-4) <!-- d-check:status-provenance -->
· `docs/reviews/review-slice-zitat-vergleich-werkzeug-fixrunde-2.md` (F-1 bis F-3) <!-- d-check:status-provenance -->
· `docs/reviews/review-slice-zitat-vergleich-werkzeug-fixrunde-3.md` (F-1). <!-- d-check:status-provenance -->
