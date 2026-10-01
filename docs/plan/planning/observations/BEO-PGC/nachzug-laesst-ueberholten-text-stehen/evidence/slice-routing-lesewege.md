**Vorgang:** slice-routing-lesewege (Review F-2, MEDIUM; Re-Review F-N2, LOW; achtzehnte Datei des Eintrags)

**Fund:** Der Use Case `ReadChanges` bekam die Zusage, ein `target` außerhalb des Alphabets antworte leer ohne Port-Aufruf. Der Godoc derselben Funktion sagte wenige Zeilen davor weiter, der Bereichs- und Limit-Kontrakt des Ports „kommt unverändert zurück — der Use Case normiert nichts und entscheidet nichts“; für ein ungültiges Ziel trug der Code die erste Zusage nicht (F-1: Alphabet-Prüfung vor dem Lese-Kontrakt). Die Fixrunde (`8c3d6e2e`) setzte den Lese-Kontrakt vor die Alphabet-Prüfung und zog den Godoc nach. Danach stand der Plan selbst widersprüchlich: Kopf und DoD-Zeile „Doku-Update“ sagten, die Spec bleibe unberührt, die Fixrunde änderte `SPEC-022` und `SPEC-031`; DoD-Zeile 3 und `welle-routing.md` führten „leer, keinen Fehler“ ohne den Qualifier „bei sonst gültiger Anfrage“ (Verifikation V-1, V-2; Re-Review F-N2). Die Closure zog Plan, Welle und die Adressliste von `slice-routing-betriebsdoku` nach.

**Form (Ausprägung):** dieselbe Klasse an zwei Träger-Typen: Go-Godoc (Satz neben einem hinzugefügten Satz) und Slice-Plan (Kopf und DoD neben der Fixrunde, die eine Norm änderte). Gefunden mit der Lese-Probe des Reviewers (Kontext um die hinzugefügten Zeilen) und des Verifiers (Plan gegen Diff); vor dem Merge gefunden, Ausgang unverändert **verkörpert**.

Quelle: `docs/reviews/review-slice-routing-lesewege.md` (F-2) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-routing-lesewege.md` (V-1, V-2) <!-- d-check:status-provenance -->
· `docs/reviews/review-slice-routing-lesewege-fixrunde-1.md` (F-N2). <!-- d-check:status-provenance -->
