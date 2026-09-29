# Beleg: slice-harness-baseline-v6-13-0

Vorgang: `slice-harness-baseline-v6-13-0` — Erstauftreten der Klasse,
befundet als F-4 (LOW) im `review-slice-harness-baseline-v6-13-0.md`,
als `BEO-PGC`-Eintrag von Review- und Verifikations-Übergabe beauftragt.

Fund: Die [`ADR-0095`](../../../../../adr/0095-review-klasse-exempt-status-check.md)-Zitat-Korrektur
(`00d96eb7`, `5bb4eabc`) änderte das Versions-Pfadsegment in der
Options-Tabelle von §Verglichene Alternativen
(`docs/plan/adr/0095-review-klasse-exempt-status-check.md:107`) — einem
Abschnitt, den [`ADR-0073`](../../../../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
§Entscheidung 1 und `AGENTS.md` §3.5 als unberührbar listen. Nach der
[`ADR-0073`](../../../../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)-Kurzform
(„das Gerüst darf sich ändern, die Aussage nie; der Referent bleibt
derselbe") ist die Änderung zulässig — Referent **gemessen** unverändert
(`grep status` auf der `v6.13.0`-Vorlage `templates/.d-check.yml`:
`matrix.status` weiterhin nur klassen-übergreifend), Commit nennt
[`ADR-0073`](../../../../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md),
§Geschichte-Zeile gesetzt — aber die Spanne zwischen Abschnitte-Liste und
Kurzform blieb offen: der nächste strenge Lauf liest denselben Commit als
§3.5-Verstoß.

Behandlung: der Architect hat die Kurzform-Lesart als tragende gezogen und
im Audit §4.2/§7 niedergelegt (`9f1eb320`); [`ADR-0095`](../../../../../adr/0095-review-klasse-exempt-status-check.md)
selbst bleibt unberührt. Die engere Lesart für neue Fälle (Angleichung der
Abschnitte-Liste an die Kurzform) ist eine Folge-ADR und mit diesem
Eintrag adressiert, nicht in diesem Zug umgesetzt.

Einordnung: erste Form der Klasse — die Korrektur selbst war
regelkonform ausgeführt (Commit-Kennung, §Geschichte, Referent-Check);
der Befund sitzt in der Auslegungs-Spanne zweier Norm-Texte, nicht in
einer Ausführung. Details und Ausgang-Vorschlag: `state.md` dieses
Eintrags.
