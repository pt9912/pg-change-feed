**Vorgang:** slice-transformationen-map-value (Review F-1, MEDIUM, daher Datei trotz Deckel)

**Fund:** Der Doc-Kommentar von `TestBackfillAndWALImagesAreByteEqualWithRules`
(`internal/bootstrap/backfill_image_parity_test.go`) sagte zu, die sichtbare Wirkung
der Regel — bei `map_value` der abgebildete Wert statt des Quellwerts — stehe im Bild,
wo die Spalte einen Wert und keinen Ausschluss trägt, „und fehlt sonst“. Der Test
trug für `map_value` kein `silent` (`nil`); sein Zweig „NULL oder ausgeschlossen“ prüfte
nichts und endete mit `continue`. Die Eingabeseite der Zusage ist die Ausschlussmenge.
Der Reviewer mutierte sie: `containsName(excluded, column) && !mapsColumn(rules, column)`
in `BuildRowImage` hebt den Ausschluss nur für eine Spalte mit `map_value`-Regel auf. Die
Mutation ließ den Paritätstest grün und färbte drei andere Tests rot
(`TestExcludedColumnIsUnreachableForEveryRuleKind`, `TestExecuteRulesNeverLeakExcludedColumns`,
`TestBuildRowImageMapValue`). Die Sicherheitsaussage war damit an drei Stellen gebunden, die
Zusage dieses Kommentars für `map_value` an keiner.

**Form (Ausprägung):** die Zusage steht im **Kommentar eines Tests**, dessen Assertion
sie für einen Regeltyp nicht bindet, und sie ist an der Eingabeseite (die Ausschlussmenge)
nur durch Mutation sichtbar. Der Reviewer stufte MEDIUM statt HIGH, weil die Eigenschaft
selbst an drei Stellen rot färbbar ist; die Begründung steht im Report. Die Fixrunde band
den Ausschluss im Paritätstest (drei Fälle `map_value/Regel an …/ausgeschlossen`); die
Verifikation bestätigte, dass dieselbe Mutation den Paritätstest jetzt rot färbt (M7, rot in
vier Paketen).

Quelle: `docs/reviews/review-slice-transformationen-map-value.md` (F-1) <!-- d-check:status-provenance -->
· `docs/reviews/verify-slice-transformationen-map-value.md` (§5 Zeile F-1, §4 M7). <!-- d-check:status-provenance -->
