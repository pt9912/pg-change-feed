**Vorgang:** Auftrag des Auftraggebers zur Kennungsdichte in Code-Kommentaren (2026-09-26),
Beobachtung am Godoc von `Publish`

**Fund:** Der Godoc von `Publish` in `internal/application/port/outbound/changestream.go` trug am
Stand `7b70b34a` in einem Absatz `ADR-0060` Teilfrage 3, `ADR-0066`, `LH-FA-REA-001` mit „ff.“,
`LH-FA-CON-003`/`005` (zwei Entscheidungen, drei Anforderungen, eine Kompaktform) und gab den
Inhalt der Spec in eigenen Worten wieder. Die Datei trug elf Kommentarzeilen mit einer Kennung
(`git show 7b70b34a:internal/application/port/outbound/changestream.go`, gezählt mit
`grep -c` über das Muster der vier Kennungsarten aus dem Suchlauf-Feld des Plans; gemessen
2026-09-26). Der Fund kam vom Auftraggeber, nicht
aus einem Review: kein Reviewer-Punkt nannte die Zählform.

Quelle: Plan `slice-code-kommentare-kennungen` §1 (Ausgangslage).
