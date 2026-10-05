**Vorgang:** slice-sdk-kompat-version-parameter (Review F-1, HIGH; Verifikation V-1, LOW; zehnte Datei des Eintrags)

**Fund:** Der Kommentar von `version_lesen` in `tools/harness/run-sdk-kompat-tests.sh`, die Meldung und der Vertrag `harness/targets/sdk-kompat.md` (§Versionen, Exit-Code-Tabelle) sagten Exit 2 für jede Version zu, die nicht der Form `X.Y.Z` folgt; der Code ließ jede Version mit Suffix durch (`0.7.0-rc.1` baute). Der Reviewer fuhr den zugesagten Ausgang mit einer Version mit Suffix nach (F-1); die Fixrunde trug die Zusage im Code. Im selben Vorgang fand der Verifier V-1: Runner-Kopf und Vertrag sagten zu, ein Eingabefehler beende den Lauf „vor jedem Bau“; die Prüfung läuft je Sprache in der Schleife, die vorderen Sprachen sind dann gebaut und gelaufen. Die Mutationen waren alle mit einer Sprache gefahren. Die Closure berichtigte den Wortlaut, der Ausgang Exit 2 bleibt.

**Form (Ausprägung):** **Zusage ohne Code** (F-1, Ausgang) und **Allaussage über einen Fall** (V-1, an einer Sprache gemessen, für den Lauf über alle Sprachen gesagt). Beide vor dem Merge von Lesern gefunden; Ausgang unverändert **verkörpert**.

Quelle: `docs/reviews/review-slice-sdk-kompat-version-parameter.md` (F-1) <!-- d-check:status-provenance -->
· `docs/reviews/verify-slice-sdk-kompat-version-parameter.md` (V-1). <!-- d-check:status-provenance -->
