# Ein neuer E2E-Metrik-Beleg deckt nur den positiven Fall, der Boundary-Fall bleibt Reviewer-/Scratch-DB-Sache

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
Testabdeckungs-Disziplin von `tools/harness/run-integration-tests.sh`,
keine eigene Sub-Area im Sinn der Modus-Deklaration).

Ein neu geschriebener E2E-Beleg für eine Metrik/View (long-format,
mehrere Zeilen je nach Zustand) prüft real nur den Fall, der eine
positive, von Null verschiedene Zeile erzeugt — den negativen/
Randfall (gesunder Zustand ohne Zeile, vollständig bestätigter
Consumer mit Wert 0) prüft niemand über einen committeten,
wiederholbaren Testlauf. Ein unabhängiger Reviewer oder Verifier kann
diesen Randfall zwar per eigener, nicht committeter Scratch-DB
bestätigen — das ist ein einmaliger Beleg dieses einen Laufs, kein
dauerhafter Wächter gegen eine spätere Regression an genau dieser
Stelle.

## Benannt, nicht gezählt

- **`e2e-drei-rtm-luecken`** (Review-Fund F-1, `docs/reviews/review-slice-e2e-drei-rtm-luecken.md` <!-- d-check:status-provenance -->):
  die neue Phase „Metriken-Minimum-Beleg" (`LH-QA-OPS-003`) prüft
  `cdc_changes_pending`>0 und `cdc_errors_total`>0 für eine erzeugte
  Fehlerklasse — nicht aber, dass eine gesunde Quelle keine
  `cdc_errors_total`-Zeile trägt oder ein vollständig bestätigter
  Consumer `cdc_changes_pending`=0 zeigt. Der Reviewer bestätigte den
  Randfall korrekt per eigener Scratch-DB, aber das ist keine
  Wiederholung im ausgelieferten Testbestand. Dasselbe Muster gilt
  bereits für ältere Zeilen derselben View (`cdc_storage_bytes`,
  `cdc_consumer_lag`) — keine Regression durch diesen Slice, aber auch
  keine Verbesserung.
