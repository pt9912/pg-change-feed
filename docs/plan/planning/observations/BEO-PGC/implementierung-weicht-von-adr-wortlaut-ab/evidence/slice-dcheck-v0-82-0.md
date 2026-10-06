**Vorgang:** slice-dcheck-v0-82-0 (Review F-1, HIGH; Re-Review zur Fixrunde F-1, LOW, gleich Verifikation V-1)

**Fund:** Nach der Messung, dass d-check v0.82.0 eine leere Commit-Range mit Exit 2 abbricht,
setzten die Träger des Bump-Ablaufs (`.claude/agents/verifier.md`, `.claude/agents/implementer.md`,
`harness/targets/pin-stale.md`) an die Stelle des Laufs, den `ADR-0157` Entscheidung 4 für jeden
Range-Teil um den Pin-Commit verlangt („jede endet mit Exit 0“), einen Nicht-Lauf mit eigenem
Beleg. Der Plan führte das als „Ausführungsregel“ einer Orchestrator-Entscheidung, nicht als
offene Frage an den Architect; ein Artefakt des Architect gab es nicht. Der Reviewer fand den
Widerspruch zur `Accepted`-ADR (HIGH). Aufgelöst über die Folge-ADR `ADR-0160` (Teil-Supersede von
`ADR-0157` Entscheidung 4, erster Spiegelstrich) und eine Fixrunde.

Im selben Vorgang, zweite Ausprägung: die Fixrunde übernahm den Wortlaut von `ADR-0160`
Entscheidung 3 in `AGENTS.md` §3.11 und fügte einen Satz an, der nicht im beschlossenen Text steht
(„Die Hausform ist `d-check`s `Dockerfile` …“). Re-Review F-1 und Verifikation V-1 (LOW); bei der
Closure gestrichen. Zwei Funde, ein Vorgang, eine Datei.

**Form (Ausprägung):** verdeckt (Abweichung als Entscheidung dargestellt, vom Reviewer als HIGH
gefunden) wie der Erstbeleg; Träger-Typ Agenten-Briefing und Target-Vertrag statt Workflow. Vor dem
Merge gefunden.

Quelle: `docs/reviews/review-slice-dcheck-v0-82-0.md` (F-1), <!-- d-check:status-provenance -->
`docs/reviews/review-slice-dcheck-v0-82-0-fixrunde.md` (F-1), <!-- d-check:status-provenance -->
`docs/reviews/verify-slice-dcheck-v0-82-0.md` (V-1). <!-- d-check:status-provenance -->
