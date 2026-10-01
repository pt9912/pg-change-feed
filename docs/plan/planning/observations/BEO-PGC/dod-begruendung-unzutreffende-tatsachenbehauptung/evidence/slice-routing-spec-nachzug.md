**Vorgang:** slice-routing-spec-nachzug (Review F-2, MEDIUM; zwölfte Datei des Eintrags)

**Fund:** `ADR-0138` führt „ein `target` außerhalb des Alphabets liefert eine leere Liste, kein
Fehler“ nur für `ReadChanges` und ausdrücklich als *Erwartung*. Die Spec setzte den Satz
unbedingt für gRPC-Stream, `GET /changes`, `ReadChanges` und (über Verweis) den SSE-Stream,
ohne Erwartungs- oder Herleitungs-Kennzeichnung; für die SQL-gestützten Wege blieb offen, was
ein `target` mit U+0000 tut, obwohl „kein Fehler“ zugesagt war. Die Setzung ging über die
Entscheidung hinaus. Geschlossen durch `ADR-0139` (Festlegung 2) und die Fixrunde der Spec.

**Form (Ausprägung):** Instanz B von `AGENTS.md` §3.12 am Träger **Pflichtenheft**: eine
Aussage über eine Menge („alle Lesewege“), breiter als die Entscheidung, die sie trägt. Die
Fixrunde hatte zusätzlich den Rand, dass Plan und Historie die Änderung der SSE-Zeile
(`SPEC-021`) berichteten, der Diff sie aber nicht enthielt (Verifikation V-1, LOW; unter dem
Deckel, hier mitgeführt): ein Nachzug-Bericht, der „X ergänzt“ sagt, ist gegen den Diff zu
prüfen. Vor dem Merge von Reviewer und Verifier gefunden, Ausgang unverändert **verkörpert**.

Quelle: `docs/reviews/review-slice-routing-spec-nachzug.md` (F-2, A-2) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-routing-spec-nachzug.md` (V-1, §4). <!-- d-check:status-provenance -->
