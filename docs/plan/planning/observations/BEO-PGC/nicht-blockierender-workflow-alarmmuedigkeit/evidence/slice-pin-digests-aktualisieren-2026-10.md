**Vorgang:** slice-pin-digests-aktualisieren-2026-10
**Fund:** Zweites Auftreten, anderer Workflow derselben Klasse: der
advisory-Nachtlauf `upstream-drift.yml` (kein Required-Status-Check, kein
Merge-Blocker) endete in allen 14 Läufen seit dem frühesten gelisteten
(2026-09-20) mit `failure`, kein einziger `success`
(gemessen: `gh run list --workflow upstream-drift.yml --limit 100 --json
conclusion,createdAt`, 14 Läufe, `[{"failure":14}]`, 2026-10-03). Der Lauf
37112891025 druckte unter anderem `DRIFT astral-sh/setup-uv@v10.1.0
Tag-Frische: gepinnt v10.1.0, neuester Release v10.2.0` und fünf
Digest-Drifts; niemand reagierte, bis der Slice den Anlass aufgriff. Ob die
Rot-Ursache über die 14 Läufe gleich blieb, ist nicht gemessen. Die Beobachtung
trifft hier den Nachtlauf, den
[`ADR-0051`](../../../../../adr/0051-cicd-pipeline-github-actions.md) selbst als
Alarmmüdigkeits-Risiko nennt.
