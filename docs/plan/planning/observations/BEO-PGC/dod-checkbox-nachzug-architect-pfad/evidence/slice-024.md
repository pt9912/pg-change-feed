# Beleg: slice-024 (ADR — Fehlerklassen-Trennung und Schwellen-Präzisierung)

Vorgang: slice-024.

Fund: Der Verifier (`docs/reviews/verify-slice-024.md`, V-1) stellte fest,
dass die DoD-Zeile „Review durchgeführt" trotz real abgeschlossenem, sauberem
Review (`docs/reviews/review-slice-024.md`, 0 HIGH/MEDIUM/LOW, 1 INFO) auf
`[ ]` stehen blieb. Dieselbe Symptom-Klasse wie
`BEO-PGC/dod-checkbox-nachzug`, aber über einen anderen Pfad: Dieser Slice
lief über die Architect-Rolle (ADR + Spec-Präzisierung), nicht über den
Implementer-Workflow `.claude/commands/implement-slice.md`, in dem die
bereits verkörperte Regel (Pre-completion-Checkliste Schritt 18/21) lebt.
Die verkörperte Regel griff deshalb nicht, weil ihr Träger rollen-gebunden
ist.

**Ausgang: weiter offen.** Ob eine zweite, rollen-übergreifende Verankerung
nötig ist (z. B. eine Checkliste im Architect-Skill/-Workflow, oder ein
Sensor, der DoD-Checkboxen gegen vorhandene Review-/Verify-Reports
gegenprüft, unabhängig von der liefernden Rolle), ist eine Entscheidung, kein
Ein-Zeilen-Fix.

Quelle: `docs/reviews/verify-slice-024.md` V-1.
