**Stand:** offen

Zweites Auftreten (2×), unter der 3×-Schwelle für eine reguläre
Skill-/Regel-Verkörperung. Kein Ausgang zugewiesen — kein Ausgang unter der
Schwelle ist der Normalzustand, kein Rückstand.

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

Zähler (abgeleitet): 2× (`evidence/slice-sdk-kotlin-publish-workflow.md`,
`evidence/slice-transformationen-kern-rename.md`).
