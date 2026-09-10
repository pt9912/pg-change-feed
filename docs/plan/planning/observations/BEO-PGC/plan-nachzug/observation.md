# BEO-PGC/plan-nachzug

**Sub-Area:** Planning-Harness (Slice-Pläne; Sub-Area-Kürzel `PGC` aus der
Modus-Deklaration)

Die Beobachtung: Der Implementer-Lauf erweitert die §3-Datei-Liste des
Slice-Plans über den gemeldeten Umfang hinaus, und der Plan-Nachzug
kommt als eigener Planner-Commit nach — im schlimmsten Fall gar nicht
(Review-Summary-Zeile: 8. Auftreten der Klasse, slice-008 kam ohne jeden
Plan-Commit aus). Die §3-Liste ist die Pfad-Kandidaten-Quelle für §8
(Sub-Area-Prüfungen); ohne Nachzug sichtet §8 einen falschen Bestand.

Deklaration: `.claude/commands/implement-slice.md` (Implementer-
Workflow, Plan-Nachzug-Schritt);
[`ADR-0045`](../../../../../plan/adr/README.md) trägt die Kennungs-
Klasse, die die Nachzug-Commits trägt.
