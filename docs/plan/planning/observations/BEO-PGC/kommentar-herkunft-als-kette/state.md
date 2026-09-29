Zustand: **verkörpert** — Ausgang: **verkörpert** → `AGENTS.md` §3.7 §Herkunft im Go-Kommentar
(höchstens **eine** Kennung je Kommentar als Rang-Zeiger auf die Norm; keine Kette, keine
Kompaktform, kein „ff.“; keine Wiedergabe von Spec- oder ADR-Inhalt in eigenen Worten; eine
Kopplung nennt die mitzuändernde Stelle), `make kommentar-kennungen` (Vertrag:
`harness/sensors/kommentar-kennungen.md`; Kandidat = Kommentarblock mit mindestens zwei
verschiedenen Kennungen oder „ff.“ hinter einer Kennung; kein Gate, keine Ausnahmeliste),
`.claude/commands/implement-slice.md` Schritt 20 (diff-skopierter Lauf) und
`.harness/skills/reviewer.md` (MEDIUM-Unterpunkt „Herkunft als mehrere Felder, Kette, „ff.“ oder
Spec-Wiederholung“, der Lauf ist Probe, kein Beleg) · seit slice-code-kommentare-kennungen.

Zähler (abgeleitet): 3× (evidence/changestream-publish-godoc.md,
evidence/slice-code-kommentare-kennungen.md,
evidence/slice-code-kommentare-bereinigung.md); der Ausgang ist
zugewiesen, weil die Regel mit dem Slice landet, der den Eintrag anlegt (dieselbe Arbeit trägt
Beobachtung und Verkörperung). Der Ursprung je Belegdatei: das erste Auftreten ist
die Stelle, das zweite die Messung des Bestands, das dritte die Lese-Funde der Bereinigung
(zwei neue Finding-Klassen: „Kürzung lässt hängende Herkunfts-Referenz zurück“ 4×,
„ungenau umgeformte Spec-Wiedergabe“ 1×).

Bestand: 597 Kandidaten (400 Nicht-Test, 197 Test) am Stand `d13ab81e`, `make kommentar-kennungen
COUNT=1` (gemessen 2026-09-26). Die Bereinigung trägt `slice-code-kommentare-bereinigung`
(Datei in `open/`); jede ihrer Tranchen misst vor und nach.

**Benannte Grenze.** Das Werkzeug prüft die Form, nicht die Wahrheit: eine Spec-Aussage in eigenen
Worten hinter **einer** Kennung ist kein Kandidat, ebenso ein Kommentar, der eine Zusage macht,
die der Code nicht trägt (Lese-Handlung des Reviewers; `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`).
Ein Sensor über Prosa bleibt ausgeschlossen (`ADR-0083`).

**Offene Frage — Grenzfälle mit zwei Ankern.** Bejaht: `slice-code-kommentare-bereinigung`
lieferte mit den 14 gekennzeichneten `TestE2E*`-Godocs den Fall, dessen Kennungs-Menge ein
Generator maschinell liest (`abdeckungsAdressiert`); die Konkretisierung von
`AGENTS.md` §3.7 (Erzeugnis-Eingabe ist Abdeckung, nicht Herkunft; Grenz-Vermerk am Block)
ist im Steering-Loop des Slice verkörpert · seit slice-code-kommentare-bereinigung.

**Trigger für die Gate-Frage** (Architect-Frage mit ADR-Vorschlag, `AGENTS.md` §3.6 und §4): die
Bereinigung endet mit Kandidatenzahl 0, und ein Kommentar mit Kennungs-Kette erreicht trotz
gelaufenem Werkzeug den Review. Vorher färbte ein Gate den Bestand rot.

**Ausgang der Gate-Frage:** nicht aufgestellt (`slice-code-kommentare-bereinigung` §7) — die
Bereinigung endet bei der strukturellen Restmenge 14, ein Gate stünde ohne Ausnahmeliste dauerhaft
rot. Neubewertung, wenn sich die Restmenge auflöst oder ein zweiter maschineller Leser hinzutritt.
