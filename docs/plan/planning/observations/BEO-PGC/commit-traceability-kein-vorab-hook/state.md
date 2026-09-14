Zustand: **geplant** — Ausgang: **geplant** → Folge-Slice `slice-073`
(`docs/plan/planning/open/slice-073-commit-msg-git-hook.md`), Architect-Verdikt
`docs/reviews/architect-verdict-commit-traceability-kein-vorab-hook.md`,
gegengeprüft und in der ADR-Frage korrigiert durch
`docs/reviews/architect-verdict-commit-traceability-kein-vorab-hook-gegengeprueft.md`
+ [`ADR-0062`](../../../../adr/0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md)
(Supersedes `ADR-0045`, nur die Klausel „Kein commit-msg-Hook").
Zähler (abgeleitet): 3× (evidence/slice-038.md, evidence/review-slice-041.md,
evidence/slice-059.md) — Schwelle erreicht, Ausgang im Lese-Schritt der
`welle-16`-Closure zugewiesen (Baseline-Regelwerk `modul-08-agentenrollen.md`
§Rollen-Sequenz für eine Welle, Planner → Architect → Planner-Zug). Der im
Register selbst benannte Lösungsvorschlag (lokaler `commit-msg`-Git-Hook)
trägt als Design-Vorgabe; die Umsetzung selbst ist kein Ein-Zug-Fix (neues,
lauffähiges Skript mit eigenen Fehlerfällen, Aktivierungs-Abhängigkeit und
Duplikations-Risiko gegenüber dem bestehenden Standing-Gate, `ADR-0045`) und
läuft über den Implementer→Reviewer→Verifier-Weg von `slice-073`, nicht als
Architect-Ein-Zug-Prosa-Ergänzung. Zusätzlich brauchte der Vorschlag eine
ADR-Korrektur, bevor `slice-073` implementiert werden darf: `ADR-0045`s
eigener Entscheidungstext schloss einen `commit-msg`-Hook ausdrücklich aus
(„Kein commit-msg-Hook") — das hatte der erste, im Planner-Kontext als
Rollenspiel geführte Architect-Zug übersehen. Ein unabhängiger,
zweiter Architect-Zug stellte den Widerspruch fest und schrieb `ADR-0062`
(Supersedes `ADR-0045`, nur diese eine Klausel) nach, statt ihn stillschweigend
zu übergehen (Modul 8 §Rollen-Regeln, `AGENTS.md` §3.5). Herkunfts-Anker bei
Umsetzung: `seit slice-073`.
