# Beleg: slice-083

Vorgang: `slice-083` — der NATS-Beispielclient.

Fund: `examples/nats-client/subject_test.go` trägt einen Kommentar, der eine
Eigenschaft **zusagt**, die die Funktion nicht hat: sie setze „kein eigenes
Schema auf einen bereits absoluten Wert". Gemessen liefert
`ChangesURL("https://feed:8080", …)` →
`http://https:%2F%2Ffeed:8080/changes?…` — sie setzt es **unbedingt** auf. Und
der Test **übt den Fall nicht aus** (einziger Input `localhost:9090`, schema-los),
also bindet nichts die Zusage an ihre Eingabeseite.

Gefunden hat es der Reviewer (Review `review-slice-083` F-2) und der Verifier
hat es unabhängig reproduziert. Der Fall berührt **kein** DoD-Kriterium.

**Dritter Gegenstand in drei Tagen:** `slice-086` ein **Negativtest** (die
Ablehnung hing an keinem Eingabewert), `slice-087` ein **E2E-Beleg** (die
Filterprüfung hing nicht am Filter), `slice-083` ein **Kommentar samt Test**
(eine zugesagte Eigenschaft, die kein Input ausübt). Derselbe Mechanismus:
*eine Aussage, die nicht an ihre Eingabeseite gebunden ist, ist grün ohne
Aussage.*

Quelle: `docs/reviews/review-slice-083.md` (F-2) ·
`docs/reviews/verify-slice-083.md` (unabhängig reproduziert) ·
`examples/nats-client/subject_test.go`, `subject.go`.
