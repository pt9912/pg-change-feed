# BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die
CI/CD-Infrastruktur, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Workflow, der bewusst **nicht** als
Required-Status-Check hinterlegt ist (Nutzerentscheidung 2026-09-13:
„Bei jedem PR, aber als eigener, nicht-blockierender Workflow" —
`slice-056`), kann dauerhaft
rot laufen, ohne dass ein Merge dadurch verhindert wird und ohne dass
irgendein Mechanismus erzwingt, dass jemand hinsieht. Dasselbe
Alarmmüdigkeits-Risiko benennt [`ADR-0051`](../../../../adr/0051-cicd-pipeline-github-actions.md)
bereits für die advisory-Nachtläufe (CVE-Scan, Pin-Freshness) — hier
tritt es erstmals auf einen PR-Trigger-Workflow übertragen auf, dessen
Rot-Zustand deutlich häufiger sichtbar würde als ein nächtlicher.

Deklaration: `slice-056-e2e-workflow-in-ci.md` §6, vom Verifier
(Verifikationsbericht zu `slice-056`) empfohlen.
