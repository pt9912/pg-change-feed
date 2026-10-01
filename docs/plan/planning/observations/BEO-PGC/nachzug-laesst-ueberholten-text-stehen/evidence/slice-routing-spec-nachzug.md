**Vorgang:** slice-routing-spec-nachzug (Review F-1, MEDIUM; sechzehnte Datei des Eintrags)

**Fund:** Der Spec-Nachzug fügte `LH-FA-CAP-009.a` den Absatz „Ziel der Backfill-Changes“ hinzu
(Backfill-Changes tragen `route_target` „des Regelstands zum Run“). Der unveränderte
Nachbarabsatz „Fail-closed vor dem Commit“ im selben Abschnitt führte als zu prüfenden Stand
weiter nur den Ausschlussstand und den Regelstand der Transformationen; der Begriff „Regelstand
zum Run“ war undefiniert, und der Folge-Plan `slice-routing-backfill-pfad` setzte eine
Erweiterung der Prüfung (Klasse `configuration`) voraus, die die Spec nicht trug. Die Lücke
lag zugleich in der Entscheidung (Architect-Frage A-1 des Reviews, geschlossen durch
`ADR-0139` vor dem Spec-Text); die Fixrunde zog den Nachbarabsatz nach.

**Form (Ausprägung):** dieselbe Klasse am Träger **Pflichtenheft**: ein hinzugefügter Absatz
macht den im Diff unberührten Nachbarabsatz desselben Abschnitts unvollständig. Gefunden mit
der Lese-Probe des Reviewers (Kontext um die hinzugefügten Zeilen); vor dem Merge gefunden,
in der Fixrunde behoben, Ausgang unverändert **verkörpert**.

Quelle: `docs/reviews/review-slice-routing-spec-nachzug.md` (F-1, A-1) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-routing-spec-nachzug.md` (§4). <!-- d-check:status-provenance -->
