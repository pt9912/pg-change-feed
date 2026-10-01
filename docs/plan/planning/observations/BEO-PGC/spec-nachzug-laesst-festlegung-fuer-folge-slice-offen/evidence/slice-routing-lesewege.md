**Vorgang:** slice-routing-lesewege (Review F-1, MEDIUM; Architect-Frage A-1; Verifikation V-1, V-3)

**Fund:** `ADR-0139` Festlegung 2 und die Zeilen `SPEC-022`/`SPEC-031` legten fest, dass ein `target` außerhalb des Alphabets mit leerer Liste und ohne Fehler antwortet, der Use Case dabei den Store nicht anfragt. Sie legten nicht fest, ob diese Antwort einen Fehler des Lese-Kontrakts (leere Quelle, `limit < 1`, invertierter Bereich, Start-/End-Position einer anderen Quelle) verdeckt. Der umsetzende Slice brauchte die Antwort beim ersten Satz Code (die Alphabet-Prüfung steht vor dem Port-Aufruf, der allein diese Fehler erzeugt); er meldete die Frage als „nicht entschieden“. Der Reviewer fand die Lücke als F-1 (MEDIUM) und wies sie dem Architect zu; entschieden hat sie der Hauptlauf der Sitzung (der Lese-Kontrakt gewinnt, Qualifier „bei sonst gültiger Anfrage“ in beiden Spec-Zeilen). Eine ADR hält die Entscheidung nicht fest; das Re-Review stufte sie als von der Festlegung gedeckte Auslegung ein.

**Form (Ausprägung):** die Festlegung, die fehlt, ist eine **Fehlerrangfolge** zwischen zwei Eingabe-Prüfungen am Lesepfad (Stärke einer Zusage), nicht ein Bezeichner-Vergleich oder eine Schicht-Zuordnung wie beim ersten Auftreten. Der Nachzug (`slice-routing-spec-nachzug`) hatte geliefert; die Lücke zeigte sich erst dem Leser, der den Randfall „zwei Kontraktfehler zugleich“ fragte.

Quelle: `docs/reviews/review-slice-routing-lesewege.md` (F-1, A-1) <!-- d-check:status-provenance -->
· `docs/reviews/review-slice-routing-lesewege-fixrunde-1.md` (F-N1) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-routing-lesewege.md` (V-1, V-3). <!-- d-check:status-provenance -->
