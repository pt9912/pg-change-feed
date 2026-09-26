**Vorgang:** slice-harness-guard-inplace-textwerkzeug (Review-Funde F-3 und F-5, je MEDIUM)

**Fund:** Zwei Tatsachenbehauptungen im Plan und in den Trägern, die am Gegenstand nicht geprüft waren.

- **F-3.** Die Definition of Done und `MR-003` sagten „Bestandsregeln … bleiben unverändert“, und der Tabellentest führte zwei Fälle als „Bestand“ (`ls | xargs -n1 pip`, `env -i pip install x`). Der Reviewer maß am Guard des Parents (`git show e98d419c:…`): beide Fälle blocken dort nicht, und `command -v pip` blockte am neuen Guard, am Parent nicht. Die Optionen hinter einem Wrapper-Präfix wurden für **alle** Klassen übersprungen; die Erweiterung der Bestandsklasse war im Text als Unverändertes ausgewiesen.
- **F-5.** Der Plan (Ausgangslage) nannte die Python-Ersetzung auf einer Kopie im Scratchpad „der zulässige Weg nach §3.1“, ohne Anker; `AGENTS.md` §3.1 führt `python` unter den Werkzeugen des Containers und nennt für die Mutation auf der Kopie nur Edit/Write und `sed … Datei > Kopie`. Die Blockmeldung der Klasse `interp` und `MR-003` (Begründung) übernahmen den Weg und lenkten in jeder Sitzung zu einem Host-Interpreter-Aufruf, den die Regel im selben Slice nicht deckt.

Beide fand der Reviewer vor dem Merge durch Messen am Parent bzw. Lesen der Regel gegen die Meldung. Die Fixrunde wies die Erweiterung als solche aus (`MR-003` Adaption, eigene Tabellengruppe „Kopf-Erkennung“, am Parent rot gemessen) und ersetzte den Weg in Meldung, `MR-003` und Plan durch die drei Wege der Regel (Edit/Write, Repo-Werkzeug hinter `make`, `sed … > Kopie`).

**Form (Ausprägung):** Instanz B von `AGENTS.md` §3.12 an einem **Wächter-Plan**: „unverändert“ ist eine Aussage über den Parent und trägt ihren Beleg-Anker nur, wenn der Parent gemessen ist; „zulässiger Weg nach §3.1“ ist eine Aussage über eine Regel und trägt ihn nur, wenn die Regel gelesen ist. Der Diff war in beiden Fällen der Ort der Wirkung (eine Meldung in jeder Sitzung), nicht nur des Textes.

Quelle: `docs/reviews/review-slice-harness-guard-inplace-textwerkzeug.md` (F-3, F-5) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-harness-guard-inplace-textwerkzeug.md` (§5 Zeilen F-3 und F-5). <!-- d-check:status-provenance -->
