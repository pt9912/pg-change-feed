**Vorgang:** slice-zitat-korrektur-vergleichseinheit (Review F-9, INFO; gemessen in `ADR-0159` Entscheidung 2)

**Fund:** `harness/conventions/MR-001-technik-dokument-heisst-pflichtenheft.md` zitiert als ersetzte Baseline-Regel „grundlagen-source-precedence.md §Spec-Straten“ mit dem Anker `#spec-straten-mehr-als-ein-spec-dokument`. Das Heading „Spec-Straten: mehr als ein Spec-Dokument“ steht in `grundlagen-referenz-richtung.md`; in `grundlagen-source-precedence.md` löst der Anker an keinem Stand auf (`vergleich` an `5d8855d9~1` und `5d8855d9`: „keine Einheit“, Exit 2). `make docs-check` meldete es nicht.

**Form (Ausprägung):** bekannte Form — ein Verweis auf eine Stelle eines anderen Dokuments nennt die falsche Datei; die gemeinte Stelle ist eindeutig (der Anker ist im Regelwerk nur dort ein Heading, und `grundlagen-source-precedence.md` verweist selbst für §Spec-Straten dorthin). Gefunden lange nach dem Merge des zitierenden Eintrags, daher eine Datei. Behoben bei der Closure dieses Slice als Zitat-Korrektur am MR-Eintrag (`ADR-0073`, `ADR-0159` Entscheidung 2, Beleg „nicht messbar“ mit Grund). Ausgang unverändert **verkörpert**.

Quelle: `docs/reviews/review-slice-zitat-korrektur-vergleichseinheit.md` (F-9). <!-- d-check:status-provenance -->
