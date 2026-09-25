Zustand: offen — Ausgang noch nicht zugewiesen. Bei 3× wäre zu prüfen,
ob der Reviewer-Skill-HIGH-Punkt „Slice-/Wellen-Chronik in
Produktionscode-Kommentar" um Test-Harness-Skripte (`tools/harness/*.sh`)
erweitert werden sollte, oder ob der engere Skopus (nur `internal/**`)
bewusst richtig bleibt — beides ist eine Architect-Entscheidung, nicht
vorwegzunehmen bei 1×.

Der zweite Beleg (`slice-backfill-e2e`, F-5) trifft eine Ausprägung außerhalb des
Produktionscodes, aber mit derselben Skopus-Frage: zwei Kommentare (ein
Test-Godoc in `test/integration/`, ein Runner-Kommentar in
`tools/harness/run-integration-tests.sh`) tragen neben Zusage bzw. Kopplung eine
Nebenklausel über die verworfene oder nicht gewählte Alternative („statt ohne
Überlappung grün zu werden“; „würde diese verzögern“, Konjunktiv) — die Form, die
`AGENTS.md` §3.7 nennt. Der Reviewer fand sie, die Fixrunde formulierte beide im
Indikativ um. Keine Slice-/Wellen-Nummer, keine Chronik; ob der Skopus des
Reviewer-Skill-Punkts diese Klasse in Tests und Skripten trägt, ist dieselbe
Architect-Frage.

Der dritte Beleg (`slice-backfill-slot-leerlauf-bestaetigung`, F-6) trifft ein Test-Paket
(`internal/bootstrap`): der Kommentar der exklusiven Instanz trug neben der Kopplung eine
Konjunktiv-Nebenklausel über den nicht gewählten Aufbau; der Reviewer fand sie, die Fixrunde
setzte den Indikativ. **Schwelle 3× erreicht mit diesem Slice**; die Skopus-Frage des
Reviewer-Skill-Punkts (Tests und Skripte neben `internal/**`) liegt beim Architect und
wird im Lese-Schritt der Welle-Closure von `welle-backfill-bestand` gelesen, nicht
entschieden.

Zähler (abgeleitet): 3× (evidence/slice-e2e-drei-rtm-luecken.md,
evidence/slice-backfill-e2e.md, evidence/slice-backfill-slot-leerlauf-bestaetigung.md).
