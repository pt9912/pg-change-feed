**Vorgang:** slice-066
**Fund:** Neue Objektklasse derselben Werkzeuggrenze. Bislang belegt: die
Erstanlage-Konvergenz von Views (aufgelöst seit `slice-015`/`slice-016`)
und von SQL-Funktionen (`slice-036`, offen). `slice-066` erweitert
`cdc.administration_request` um eine Spalte **und** die
`request_kind`-CHECK-Klausel an einer **bereits bestehenden** Tabelle:
real gemessen (Reviewer-Rolle, PostgreSQL 18 / d-migrate 1.3.1,
`docs/reviews/review-slice-066.md`) konvergiert die **Spalte** (Exit 0),
die **CHECK-Klausel** jedoch nicht — Exit 5 (`POST_EXECUTE_DRIFT`), und
beim Änderungsversuch entfällt die bestehende Klausel sogar, ohne dass die
neue entsteht (auch unter umbenanntem Constraint, also unabhängig vom
Namen). Aufgelöst über die etablierte Ausweichform
`tools/schema/nacharbeit-administration.sql` (idempotent, dreifach
nachgemessen), die damit einen zweiten Grund trägt — neben der
Funktionsklasse nun auch die CHECK-Klausel an einer Bestandstabelle.
