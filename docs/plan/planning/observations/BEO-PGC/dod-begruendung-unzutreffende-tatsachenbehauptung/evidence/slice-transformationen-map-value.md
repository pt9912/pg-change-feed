**Vorgang:** slice-transformationen-map-value (Verifikation V-1, LOW; zehnte Datei des Eintrags)

**Fund:** DoD Punkt 1 nannte nach der Fixrunde als Ausnahme „bis auf die zwei zwangsläufig
geänderten“ Erwartungen (`TestTransformationKindsIsAClosedSet`, `TestTransformationSpecBuild`);
DoD Punkt 3 sagte, der Eigenschaftstest erfasse `map_value` „ohne manuelle Ergänzung“. Der
Verifier fuhr die Testdateien des Parent-Stands gegen den Produktivcode von `HEAD`: sechs Tests
färbten sich rot, nicht zwei — die zwei genannten und vier Abbrüche `Regeltyp "map_value" ohne
Fall in diesem Test` in den drei Fixture-Schaltern der Eigenschafts- und Paritätstests. Die vier
Abbrüche sind Absicht (ein neuer Typ braucht seinen Fall; der Plan §3 nennt die Schalter und die
Übergabe aus `slice-transformationen-backfill-pfad`), standen aber nicht im DoD-Wortlaut; der Fall
im Schalter ist die manuelle Ergänzung, ohne die die Schleife über die Domänen-Menge abbricht. Die
Closure fuhr den Gegenlauf nach (Kopie von `HEAD`, acht Testdateien auf `e5a11979`, `go test -race`
über fünf Pakete: Exit 1, sechs Zeilen `--- FAIL`) und zog den Wortlaut beider Punkte nach.

**Form (Ausprägung):** Instanz B von `AGENTS.md` §3.12 an einem **Zählwort im DoD-Wortlaut**:
„zwei“ nannte eine Menge, die nicht am Gegenstand gemessen war; der Diff war richtig, der Plan
nannte die vier Schalter an anderer Stelle. Vorstufe im selben Träger: Review F-2 (DoD-Wortlaut
gegen die Auslegung in §3), als Nachzug-Deckel-Fall in der Closure-Notiz geführt.

Quelle: `docs/reviews/verify-slice-transformationen-map-value.md` (V-1, §4 Gegenlauf) <!-- d-check:status-provenance -->
· `docs/reviews/review-slice-transformationen-map-value.md` (F-2). <!-- d-check:status-provenance -->
