# Beleg: slice-089

Vorgang: `slice-089` — die Register-Ausgänge verkörpern.

Fund: `harness/sensors/coverage-gate.md` §Zählbasis erklärte einen **Mechanismus
eines Werkzeugs**, ohne das Werkzeug zu belegen: „…`tools/coverage-gate.sh`
formatiert `%.2f` — daraus werden `71.9%` und `71.94%`." **Gemessen** (Reviewer,
am Skript und live): das Skript parst die **gedruckte** Zeile (`:35`) und
formatiert **sie** (`:47`): `71.9` → `71.90%`. `71.94 %` ist die
**deduplizierte Rechnung** und steht vier Zeilen darüber korrekt als „`71,94 %`,
gedruckt `71.90%`" — die Zeile warf zwei verschiedene Dinge in einen Topf und
nannte damit **zwei verschiedene gedruckte Zeilen für denselben Stand**.

Gefunden als `review-slice-089` **F-1 (HIGH)**. **Das Werkzeug lag vor** — die
Erklärung stand in einem Dokument über genau dieses Skript, und ein Blick in
`tools/coverage-gate.sh:35/:47` hätte sie widerlegt.

**Besonderheit:** der Beleg stammt aus dem Vorgang, der die **Herkunfts-Regel**
(`AGENTS.md` §3.12) einführt — die Regel gegen den ungeprüften Mechanismus-Satz
ist damit **an ihrer eigenen Einführung** fällig geworden.

Quelle: `docs/reviews/review-slice-089.md` (F-1) · `tools/coverage-gate.sh`
(`:35` parst, `:47` formatiert) · `harness/sensors/coverage-gate.md` §Zählbasis.
