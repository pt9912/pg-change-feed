**Vorgang:** slice-transformationen-kern-rename (Review-Fund F-2 MEDIUM, Mutation M15; Fixrunde `TestConsumeEveryRuleIsCheckedForApplicability`)

**Fund:** Die Anwendbarkeits-Prüfung des `Assembler` (`checkTransformations`, prüft jede Regel der Bindung gegen die Spalten der Relation vor der Serialisierung) war nur für die **erste** Regel einer Bindung an ihre Eingabeseite gebunden. Der Reviewer beschränkte die Schleife per Mutation auf `rules[:1]`: beide betroffenen Pakete blieben grün, weil jeder der Nichtanwendbarkeits-Tests genau eine Regel trug. Der Fall „zweite Regel nicht anwendbar, erste anwendbar“ hatte keinen Test. Die Fixrunde band die Position der Regel mit vier Fällen (zweite von zwei, erste von zwei, letzte von drei, mittlere von drei); der Verifier setzte dieselbe Mutation erneut und sah sie rot.

**Form (Ausprägung):** die Eingabeseite einer **Schleife über eine Liste** ist die **Position** des betroffenen Elements; ein Test mit einem Element bindet nur den Fall „das eine Element“. Neben der Schleife verlangt derselbe Vorgang die Bindung an die Verzweigungen der Regelmenge (Prüfung gegen die Spalten der Relation, Prüfung gegen die umbenannten Schlüssel im Bild, jeweils mit der Mutation, die genau diese Prüfung entfernt).

Quelle: `docs/reviews/review-slice-transformationen-kern-rename.md` (F-2, Mutation M15) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-transformationen-kern-rename.md` (§4 Zeile M3, §5 Zeile F-2). <!-- d-check:status-provenance -->
