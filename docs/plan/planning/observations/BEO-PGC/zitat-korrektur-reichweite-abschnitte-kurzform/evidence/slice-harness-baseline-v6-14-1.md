# Beleg: slice-harness-baseline-v6-14-1

Vorgang: `slice-harness-baseline-v6-14-1` — drittes Auftreten der Klasse,
befundet als F-2 (MEDIUM) im `review-slice-harness-baseline-v6-14-1.md`.

Fund: Die Zitat-Korrektur an
[`ADR-0095`](../../../../../adr/0095-review-klasse-exempt-status-check.md)
(`eadf3054`) änderte erneut das Versions-Pfadsegment in der Options-Tabelle von
§Verglichene Alternativen — einem Abschnitt, den die Abschnitte-Liste von
[`ADR-0073`](../../../../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
§Entscheidung 1 und `AGENTS.md` §3.5 (Stand vor dem Slice) als unberührbar
führten. Referent gemessen gleich: `templates/.d-check.yml` zwischen v6.14.0
und v6.14.1 roh `cmp` Exit 0 (Slice-Plan, Abschnitt „Verweise, Pins und
Records“; Verifikation `verify-slice-harness-baseline-v6-14-1.md` §3).

Zweite Form im selben Vorgang (zählt nicht ein zweites Mal): Review F-3
(MEDIUM) — die Pin-Umstellung in `MR-001` bis `MR-004` (`5d8855d9`) ist eine
Zitat-Korrektur an immutablen Einträgen, die weder Plan noch Commit-Message als
solche führten; `make doc-immutable` über die volle Range endet mit vier
`core-drift-vcs`.

Behandlung: Konflikt-Pfad über den Architect (Modul 8). Ausgang
[`ADR-0157`](../../../../../adr/0157-zitat-korrektur-reichweite-nach-aussage-und-mr-pins.md)
(`f1f6ae70`, Teil-Supersede von `ADR-0073`): Reichweite nach Aussage statt nach
Abschnitt (Entscheidung 1), Wortlaut von `AGENTS.md` §3.5 (Entscheidung 2,
`01399b76`), MR-Einträge in der Klasse mit Record-Belegform und eigener
Pin-Commit beim Bump (Entscheidungen 3 und 4).
