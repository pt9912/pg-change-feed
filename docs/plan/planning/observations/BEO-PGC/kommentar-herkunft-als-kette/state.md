Zustand: **verkörpert** — Ausgang: **verkörpert** → `AGENTS.md` §3.7 §Herkunft im Go-Kommentar
(höchstens **eine** Kennung je Kommentar als Rang-Zeiger auf die Norm; keine Kette, keine
Kompaktform, kein „ff.“; keine Wiedergabe von Spec- oder ADR-Inhalt in eigenen Worten; eine
Kopplung nennt die mitzuändernde Stelle), `make kommentar-kennungen` (Vertrag:
`harness/sensors/kommentar-kennungen.md`; Kandidat = Kommentarblock mit mindestens zwei
verschiedenen Kennungen oder „ff.“ hinter einer Kennung; kein Gate, keine Ausnahmeliste),
`.claude/commands/implement-slice.md` Schritt 20 (diff-skopierter Lauf) und
`.harness/skills/reviewer.md` (MEDIUM-Unterpunkt „Herkunft als mehrere Felder, Kette, „ff.“ oder
Spec-Wiederholung“, der Lauf ist Probe, kein Beleg) · seit slice-code-kommentare-kennungen.

Zähler (abgeleitet): 2× (evidence/changestream-publish-godoc.md,
evidence/slice-code-kommentare-kennungen.md) — unter der Schwelle von 3×; der Ausgang ist
zugewiesen, weil die Regel mit dem Slice landet, der den Eintrag anlegt (dieselbe Arbeit trägt
Beobachtung und Verkörperung). Beide Belegdateien haben denselben Ursprung: das erste Auftreten ist
die Stelle, das zweite die Messung des Bestands.

Bestand: 597 Kandidaten (400 Nicht-Test, 197 Test) am Stand `d13ab81e`, `make kommentar-kennungen
COUNT=1` (gemessen 2026-09-26). Die Bereinigung trägt `slice-code-kommentare-bereinigung`
(Datei in `open/`); jede ihrer Tranchen misst vor und nach.

**Benannte Grenze.** Das Werkzeug prüft die Form, nicht die Wahrheit: eine Spec-Aussage in eigenen
Worten hinter **einer** Kennung ist kein Kandidat, ebenso ein Kommentar, der eine Zusage macht,
die der Code nicht trägt (Lese-Handlung des Reviewers; `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`).
Ein Sensor über Prosa bleibt ausgeschlossen (`ADR-0083`).

**Offene Frage — Grenzfälle mit zwei Ankern.** 5 von 30 Kandidaten der Stichprobe (17 %,
abgeleitet) und 4 von 8 der Review-Stichprobe tragen zwei Anker mit je einer eigenen Aussage der
Stelle; die Regel verlangt dort einen Anker und die Stelle. Adresse: `slice-code-kommentare-bereinigung`
(DoD Liefer-Punkt 3: ein Kandidat mit zwei Ankern, den der Implementer als konform begründet, geht
als Meldung an den Planner — Konkretisierung von `AGENTS.md` §3.7, kein Ausnahme-Eintrag; §4
Rückführung (a)).

**Trigger für die Gate-Frage** (Architect-Frage mit ADR-Vorschlag, `AGENTS.md` §3.6 und §4): die
Bereinigung endet mit Kandidatenzahl 0, und ein Kommentar mit Kennungs-Kette erreicht trotz
gelaufenem Werkzeug den Review. Vorher färbte ein Gate den Bestand rot.
