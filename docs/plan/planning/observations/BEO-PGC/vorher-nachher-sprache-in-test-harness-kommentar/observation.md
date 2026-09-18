# Implizite Vorher/Nachher-Sprache in einem Test-Harness-Bash-Kommentar, außerhalb des Produktionscode-Skopus

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area).

`BEO-PGC/slice-chronik-in-code-kommentar` und der zugehörige
Reviewer-Skill-HIGH-Punkt „Slice-/Wellen-Chronik in
Produktionscode-Kommentar" sind bewusst auf **Produktionscode**-Pfade
skopiert (`internal/**`, nicht `Test*`-Godoc). Dieselbe rhetorische
Klasse — ein Kommentar, der implizit über einen Vorher/Nachher-Zustand
spricht statt über den reinen Ist-Zustand — kann aber auch in einem
Test-Harness-Skript (`tools/harness/*.sh`) auftreten, wo der HIGH-Punkt
nicht greift. Ob dieselbe Disziplin dort ebenfalls gelten sollte (§3.7
nennt „Code, Konfiguration und Skripte" explizit, ohne Test-Harness
auszunehmen) oder ob der engere Skopus bewusst richtig ist, ist noch
nicht entschieden — diese Beobachtung hält den ersten realen Fall fest.

## Benannt, nicht gezählt

- **`e2e-drei-rtm-luecken`** (Review-Fund F-2, `docs/reviews/review-slice-e2e-drei-rtm-luecken.md` <!-- d-check:status-provenance -->):
  `tools/harness/run-integration-tests.sh:2184-2189`, Kommentar zur
  neuen „Metriken-Minimum-Beleg"-Phase: „Dieser Beleg schließt die
  beiden zuvor fehlenden real" — implizite Vorher/Nachher-Sprache. Der
  Reviewer stufte den Fund als INFO statt HIGH ein, weil der
  HIGH-Skopus explizit auf Produktionscode-Pfade begrenzt ist und der
  Satz zusätzlich an `LH-QA-OPS-003` verankert ist (kein bloßer
  Slice-Bezug).
