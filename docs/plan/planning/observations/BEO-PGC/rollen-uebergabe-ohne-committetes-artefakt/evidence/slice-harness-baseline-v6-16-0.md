**Vorgang:** slice-harness-baseline-v6-16-0 (Re-Review zur Fixrunde R-1, LOW, gleich Verifikation V-1; Halt des Planners bei der Closure)

**Fund:** Die Vorgabe des Architect zur Fixrunde (Klammer in `AGENTS.md` §3.5, Zeile 10 von
`harness/targets/zitat-vergleich.md`, Lese-Rolle des Verifiers) stand nur im Bericht an den
Orchestrator; der Fixrunden-Auftrag gab sie wieder („mit dem `formnorm`-`cmp` als Beleg“), der
Implementer las sie als Zusatz zur `zitat-vergleich`-Zeile. Das Re-Review nannte sie im
Eingangs-Kontext „vom Auftraggeber wörtlich übernommen, nicht als Artefakt im Repo“ und fand die
Rest-Unschärfe (R-1). Der Planner hielt bei der Closure an, weil der Nachzug nach R-1 aus dem
Repo nur gegen den Wortlaut von `ADR-0161` zu entscheiden war; der Orchestrator nannte die
Ursache (sein Auftrag) und die Lesart (Ersatz, nach `ADR-0161` §Konsequenzen).

**Form (Ausprägung):** Übergabe Architect → Orchestrator → Implementer ohne committetes
Artefakt; Folge ein Träger, der vom ADR-Wortlaut abweicht
(`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab`). Vor dem Merge gefunden.

Quelle: `docs/reviews/review-slice-harness-baseline-v6-16-0-fixrunde.md` (Eingangs-Kontext, R-1), <!-- d-check:status-provenance -->
`docs/reviews/verify-slice-harness-baseline-v6-16-0.md` (V-1). <!-- d-check:status-provenance -->
