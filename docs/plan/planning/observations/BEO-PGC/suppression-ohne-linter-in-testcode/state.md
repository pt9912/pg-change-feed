Zustand: offen (**1×**) — unter der Schwelle, kein Ausgang zugewiesen. Behoben im Vorgang:
die Fixrunde `03a1a20a` fängt `Exception` statt `BaseException` und trägt einen Zusage-Kommentar;
`git grep -n -E "noqa|nolint|type: ignore|pylint:|# pragma" -- sdks tools` trifft keine Zeile
(Verifikations-Report, gemessen). Die Entscheidung, `AGENTS.md` §3.2 sinngemäß auch für das
Python-`# noqa` zu lesen, traf der Hauptlauf; ein Textnachzug in §3.2 (Python-Form nennen) ist
nicht beschlossen — der Anlass dafür ist ein zweites Auftreten.

Zähler (abgeleitet): **1×** (evidence/slice-routing-sdk-realserver-e2e.md).
