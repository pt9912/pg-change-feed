**Stand:** **verkörpert** in `AGENTS.md` §3.5 (Absatz „Eine Abweichung vom
Wortlaut ist eine Frage, keine Regel.“) und `.claude/agents/implementer.md`
(`seit slice-dcheck-v0-82-0`; Lese-Schritt der wellenlosen Closure, Architect-Zug
`64054840`).

Drittes Auftreten (`evidence/slice-dcheck-v0-82-0.md`) in der **verdeckten**
Ausprägung wie das erste: Träger des Bump-Ablaufs setzten eine
„Ausführungsregel“ an die Stelle von `ADR-0157` Entscheidung 4, ohne Artefakt des
Architect (Review F-1, HIGH; aufgelöst über `ADR-0160`), dazu im selben Vorgang
ein Zusatzsatz zum beschlossenen Wortlaut von `ADR-0160` in `AGENTS.md` §3.11
(Re-Review F-1, LOW; bei der Closure gestrichen).

Das zweite Auftreten ist die **benannte** Ausprägung (Abweichung als Risiko und
Auslegung geführt, vom Verifier als konform gelesen), das erste die
**verdeckte** (Darstellung als Entscheidung, vom Reviewer als HIGH gefunden). Ein
drittes Auftreten in einer der beiden Ausprägungen erreicht die Schwelle; der
Lese-Schritt der Closure von `welle-transformationen` liest die Einträge ab 3×
und liest diesen nicht.

Kein Auftreten: der Port-Schnitt von `slice-transformationen-antragsweg-usecase` (ein Port mit
zwei Lese-Methoden) entspricht dem Wortlaut „ein neuer Outbound Port“ von `ADR-0112` und ist
mit `ADR-0034` vereinbar; die Plan-Aussage „weicht von der ADR-Zählung ab“ war unzutreffend
(Review F-5, LOW) und ist in der Fixrunde berichtigt — sie zählt nicht als drittes Auftreten.

Kein Auftreten: `slice-transformationen-backfill-pfad` führt für den Wechsel des Regelstands im Run
einen eigenen Sentinel (`ErrTransformationStateChanged`) mit der Klasse `configuration`, während
`ADR-0117` Festlegung 5 die Abbildung im Wortlaut an `ErrExclusionStateChanged` bindet. Die Klasse ist
die der ADR; der Plan nennt den Sentinel als Plan-Drift, Review (F-2) und Verifikation lasen keinen
Widerspruch zum Wortlaut — eine Abweichung, die als Entscheidung dargestellt oder als Risiko geführt
würde, liegt nicht vor; die drei Fragen des Implementers gehen an den Architect. Der Vorgang zählt
nicht.

Kein Auftreten: `slice-transformationen-start-reihenfolge` liefert die Ordnung „offene Anträge
vor `stream.Run` verarbeiten“ als **eigenen Slice**, während `ADR-0112` Folgepflicht 5 sie in den
E2E-Slice legt, wenn Kriterium (c) nicht trägt. Der Inhalt der Entscheidung ist geliefert (die
Ordnung steht im Code und ist an die Eingabeseite gebunden, Verifikation), der Wortlaut der
Präambel („der Planner formt daraus Welle und Slices“) überlässt den Schnitt dem Planner; die
Zuschnitts-Abweichung ist als Frage an den Auftraggeber in `welle-transformationen` §4
(Abweichung 2) geführt und weder als Entscheidung dargestellt noch von Reviewer oder Verifier
beanstandet. Der Vorgang zählt nicht; die Frage bleibt bei der Closure der Welle.

Kein Auftreten: `slice-routing-backfill-pfad` setzte den Lesefehler eines Regelstands im Run
auf die Klasse der Ursache und wich damit vom Wortlaut von `ADR-0139` Festlegung 1
(`configuration`) ab. Die Abweichung folgte dem Bestand der zwei anderen Stände; falsch war die
ADR-Aussage („derselbe Mechanismus“), nicht die Umsetzung — der Architect berichtigte sie mit
`ADR-0141`. Die Gegenrichtung, eine Fixrunde, die dem Wortlaut folgte, wurde zurückgenommen.
Der Vorgang zählt unter `BEO-PGC/adr-aussage-breiter-als-ihre-messung`, nicht hier.

Viertes Auftreten (`evidence/slice-harness-baseline-v6-16-0.md`) in einer **dritten**
Ausprägung, der **unbeabsichtigten**: der Träger (`AGENTS.md` §3.5) wich nicht als
Entscheidung vom Wortlaut von `ADR-0161` ab, sondern durch die Lesart eines Fixrunden-Auftrags
(Zusatz statt Ersatz), dessen Vorlage nicht im Repo lag (Review F-1 MEDIUM, Re-Review R-1 /
Verifikation V-1 LOW; bei der Closure an den Wortlaut angeglichen). Die verkörperte Regel hat
gegriffen: der Planner hielt an, statt die Abweichung in weitere Träger zu schreiben. Nach dem
vierten Auftreten einer in Prosa verkörperten Klasse (Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 3): **kein Sensor möglich** — ob ein Satz in einem Träger
dasselbe sagt wie der Wortlaut einer ADR, ist eine Lese-Handlung; die Leser sind Reviewer und
Verifier, die alle vier Auftreten vor dem Merge fanden. Die Ursache dieses Auftretens zählt
gesondert unter `BEO-PGC/rollen-uebergabe-ohne-committetes-artefakt`.

Zähler (abgeleitet): 4× (`evidence/slice-sdk-kotlin-publish-workflow.md`,
`evidence/slice-transformationen-kern-rename.md`, `evidence/slice-dcheck-v0-82-0.md`,
`evidence/slice-harness-baseline-v6-16-0.md`).
