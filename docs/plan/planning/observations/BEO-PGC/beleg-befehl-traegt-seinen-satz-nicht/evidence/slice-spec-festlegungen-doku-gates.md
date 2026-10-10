**Vorgang:** slice-spec-festlegungen-doku-gates (Re-Review N-1, HIGH)

**Fund:** Die Fixrunde-Liste des Plans (§3, Punkt F-9) sagte, der Satz „keine Zeile unter `spec/`
trägt den Marker `d-check:status-provenance` und zugleich ein Token der Klassen `adr` oder `slice`“
werde zusätzlich mit `git grep -c status-provenance -- spec/` belegt. Der Befehl druckt
`spec/pflichtenheft.md:2`: er zählt die Zeilen mit dem Marker und sagt nichts darüber, ob eine davon
ein Token trägt. §6 desselben Plans nannte den tragenden Befehl
(`git grep -n status-provenance -- spec/ | grep -cE 'ADR-[0-9]{4}|slice-[0-9]{3}'`, druckt `0`).
Der Reviewer fand es durch Ausführen; die Fixrunde `3122ff63` nennt den Befehl aus §6, der Verifier
fuhr ihn nach (`0`).

**Form (Ausprägung):** Form **Befehl** — die Zusammenfassung eines Belegs in einer zweiten Stelle
desselben Plans nennt einen kürzeren Befehl als die Stelle, die den Satz trägt. Schwere HIGH (der
Skill führt „Beleg trägt seinen Satz nicht“ unter HIGH), daher eine Datei trotz Deckel; vor dem
Merge gefunden.

Quelle: `docs/reviews/review-slice-spec-festlegungen-doku-gates-re-review.md` (N-1) <!-- d-check:status-provenance -->
· `docs/reviews/verify-slice-spec-festlegungen-doku-gates.md` (§1, Zeile „Review durchgeführt“). <!-- d-check:status-provenance -->
