**Vorgang:** slice-schema-rollout-zentrale-idempotenz-wache
**Fund:** Zwei getrennte Instanzen. (1) Der Implementer trug in einer
selbst gefundenen Fixrunde („skip → allow-destructive") Vorher/Nachher-
Prosa über die verworfene Vorgänger-Implementierung in `Makefile` und
`tools/schema/rolloutguard/guard.go` ein („statt beim bloßen Überspringen
von `--execute` verlustig zu gehen", „ein reines Überspringen … ließ …
nicht zurückkommen"). Der unabhängige Reviewer fing beide Stellen vor
Merge (F-1, HIGH; Fixrunde real geprüft, auf Ist-Zustand mit
Regressionstest-Anker umgeschrieben) — siebtes/achtes Auftreten,
dieselbe Diagnose wie bisher: die Verkörperung trägt, kein Hard-Rule-
Verstoß hat `main` erreicht.

(2) **Neue Grenze der bestehenden Verkörperung:** Eine dritte,
wortidentische Instanz derselben Formulierung stand in
`harness/README.md` (Sensors-Tabelle, `make schema-rollout`-Zeile) — vom
Reviewer explizit gegen §3.13 geprüft (bestanden), aber **nicht** gegen
den Reviewer-Skill-HIGH-Punkt „Slice-/Wellen-Chronik in
**Produktionscode**-Kommentar", weil `harness/README.md` Doku-Prosa ist,
kein Code-Kommentar — außerhalb des wörtlichen Skopus des bisherigen
HIGH-Punkts. Gefunden wurde sie erst vom **Verifier** (nicht vom
Reviewer), der den vollen Diff-unabhängigen Kontext las. Erstes real
belegtes Beispiel dafür, dass dieselbe Chronik-Klasse außerhalb von
Produktionscode auftritt und die heutige Verteidigungslinie (Reviewer-
HIGH-Punkt, Skopus „Produktionscode-Kommentar") sie nicht zwingend
abdeckt — dem nächsten Lese-Schritt zur Prüfung vorgelegt, ob der
HIGH-Punkt um Doku-Träger erweitert werden sollte.
