# Zeitkonstanz einer Prüfung ist Erwartung ohne Messung

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Token-Prüfung der Driving-Adapter, `internal/application/port/apiauth`).

Eine Spec-Zusage verlangt eine zeitkonstante Prüfung des API-Tokens
([`SPEC-035`](../../../../../../spec/pflichtenheft.md)); kein Test misst Laufzeit, und ein Laufzeit-Test
wäre auf einem geteilten Runner nicht belastbar. Die Zusage trägt allein die Lesung des Codes (kein früher
Abbruch, Vergleich über Digests gleicher Länge); sie bleibt eine Erwartung, keine gemessene Eigenschaft.

## Benannt, nicht gezählt

Der Reviewer nennt in F-2 (`docs/reviews/review-slice-api-token-mehrfach-konfiguration.md`), <!-- d-check:status-provenance -->
dass der Kommentar der Funktion `Classify` die Grenze etwas weiter formuliert, als der Code sie trägt (leeres
Aufruf-Token kehrt vor dem Hash zurück, die Hash-Zeit hängt an der Länge des Aufruf-Tokens).
