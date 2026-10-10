**Vorgang:** slice-spec-festlegungen-doku-gates (Review F-7, MEDIUM)

**Fund:** `ADR-0163` nimmt zwei Formen an, an denen `make docs-check` grün endet, obwohl die Regel
verletzt sein kann: der Kommentar `d-check:status-provenance` hebt einen Token-Befund von `matrix`
auch in `spec/` auf (auch in Inline-Code, M33), und `matrix.exempt-paths` nimmt die ADRs 0039 und
0041 auch von den Regeln aus; dazu ist ein Token in einem Fence für `matrix` kein Treffer (R13).
Der Vertrag `harness/sensors/docs-check.md` §Grenze („was das Grün nicht abdeckt“) nannte keine der
drei Formen. Die Fixrunde `7f139836` trägt sie als Grenze 12, das Architect-Verdikt zu
`matrix-inactive` ergänzte dort den Status nur am Link.

**Form (Ausprägung):** zweites Auftreten, anderer Träger-Typ als das erste: kein nachgebildetes
Format eines eigenen Werkzeugs, sondern ein fremdes Werkzeug (d-check) mit angenommenen Lücken, die
die Entscheidung kennt und der Vertrag nicht nennt. Gemeinsam ist das falsche Grün an einer Form,
die der Vertrag nicht als Grenze führt.

Quelle: `docs/reviews/review-slice-spec-festlegungen-doku-gates.md` (F-7). <!-- d-check:status-provenance -->
