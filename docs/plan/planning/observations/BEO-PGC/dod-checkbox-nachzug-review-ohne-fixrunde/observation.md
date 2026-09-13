# BEO-PGC/dod-checkbox-nachzug-review-ohne-fixrunde

**Sub-Area:** Planning-Harness (Slice-Lifecycle; Sub-Area-Kürzel `PGC` aus
der Modus-Deklaration)

Die Beobachtung: `BEO-PGC/dod-checkbox-nachzug` (verkörpert seit welle-5)
verankert den DoD-Checkbox-Nachzug an zwei Stellen im Implementer-Workflow
(`.claude/commands/implement-slice.md`) — Schritt 18 (im selben Lauf, vor
Handoff an den Reviewer) und Schritt 21 (nach Reviewer-Findings, im Rahmen
einer Fixrunde). Für die DoD-Zeile „Review durchgeführt" greift **keiner**
der beiden Träger strukturell: Schritt 18 kann sie nicht erfüllen, weil das
Review erst *nach* dem Implementer-Handoff stattfindet; Schritt 21 greift nur,
wenn der Reviewer Findings meldet und eine Fixrunde läuft — bei einem
**sauberen** Review (0 HIGH/MEDIUM/LOW, wie bei `slice-045`) gibt es keine
Fixrunde und damit keinen erneuten Implementer-Lauf, der die Checkbox
nachziehen könnte. Der Verifier findet die Lücke deshalb bei jedem sauberen
Review erneut, unabhängig vom Implementer-Pfad — eine andere Ursache als
`BEO-PGC/dod-checkbox-nachzug-architect-pfad` (dort läuft die Substanz über
eine rollenfremde Rolle; hier läuft sie über den regulären
Implementer-Workflow, aber ohne einen zweiten Implementer-Lauf, an den sich
der Nachzug hängen könnte).

## Benannt, nicht gezählt

Dieselbe Symptom-Klasse trat bereits bei `slice-039`, `slice-043` und
`slice-044` auf (`verify-slice-039.md`, `verify-slice-043.md`,
`verify-slice-044.md`, jeweils ein sauberes oder bereits fixrunden-
abgeschlossenes Review, Checkbox dennoch offen) — alle drei liegen vor der
Anlage dieses Eintrags und zählen laut Präzedenzfall `BEO-PGC/plan-nachzug`
nicht mit (historische Vorkommen vor Registrierung zählen nicht).
