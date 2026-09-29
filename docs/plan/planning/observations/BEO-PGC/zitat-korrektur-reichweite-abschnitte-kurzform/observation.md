# BEO-PGC/zitat-korrektur-reichweite-abschnitte-kurzform

**Sub-Area:** Planning-Harness (Zitat-Korrektur an `Accepted`-ADRs;
Sub-Area-Kürzel `PGC` aus der Modus-Deklaration, Repo-Default).

Die Beobachtung: [`ADR-0073`](../../../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
trägt die Reichweite der Zitat-Korrektur an einem `Accepted`-ADR in **zwei
Formen**, die nicht dieselbe Menge benennen: §Entscheidung 1 listet die
unberührbaren **Abschnitte** (§Entscheidung, §Konsequenzen, §Verglichene
Alternativen, §Status, `Supersedes`-Kette — „unberührbar"), während die
**Kurzform** („das Gerüst darf sich ändern, die Aussage nie; der Referent
bleibt derselbe") die Korrektur an der Aussage und am Referenten festmacht
und damit auch Stellen zulässt, die in der Abschnitte-Liste stehen. Ein
strenger Lauf, der die Abschnitte-Liste liest, wertet dieselbe Korrektur
als §3.5-Verstoß (`AGENTS.md` §3.5), den die Kurzform trägt — die Spanne
wird erst sichtbar, wenn eine Korrektur genau dort sitzt.

Konkret (Erstauftreten): die
[`ADR-0095`](../../../../adr/0095-review-klasse-exempt-status-check.md)-Zitat-Korrektur
(`00d96eb7`, `5bb4eabc`
— Versions-Pfadsegment in der Options-Tabelle von §Verglichene Alternativen,
Referent unverändert, Commit nennt `ADR-0073`, §Geschichte-Zeile gesetzt)
sitzt genau in einem Abschnitt der Liste; Review F-4 (LOW) hat die Spanne
benannt, der Architect hat die Kurzform-Lesart im Audit §4.2/§7 als
tragende gezogen (Referent **gemessen**: die `v6.13.0`-Vorlage
`templates/.d-check.yml` trägt `matrix.status` weiterhin nur
klassen-übergreifend) und die engere Lesart (Folge-ADR) für neue Fälle
festgelegt.

**Warum das zählt:** Die Träger der beiden Formen sind zwei Norm-Texte
(`AGENTS.md` §3.5 und `ADR-0073`), deren Lesart auseinanderläuft, und der
Konflikt sitzt nicht im geänderten Code, sondern in der Auslegung — kein
Sensor kann die Reichweite einer Ausnahme-Regel messen. Die verfügbare
Falsifikation ist der strenge Gegen-Lese-Lauf: bis zur Klärung liest
derselbe Commit je nach Träger „zulässig" oder „Verstoß".
