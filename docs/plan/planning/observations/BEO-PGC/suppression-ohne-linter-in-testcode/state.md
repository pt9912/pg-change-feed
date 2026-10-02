Zustand: **verkörpert** — `AGENTS.md` §3.2 nennt die Python-Formen (`# noqa`, `# type: ignore`, `# pylint:`) ausdrücklich (Anker `seit welle-routing`, Commit `f7d42312`; Zeile gelesen); der Auftraggeber hat die Ergänzung am 2026-10-02 freigegeben. Behoben im Vorgang:
die Fixrunde `03a1a20a` fängt `Exception` statt `BaseException` und trägt einen Zusage-Kommentar;
`git grep -n -E "noqa|nolint|type: ignore|pylint:|# pragma" -- sdks tools` trifft keine Zeile
(Verifikations-Report, gemessen). Die Entscheidung, `AGENTS.md` §3.2 sinngemäß auch für das
Python-`# noqa` zu lesen, traf der Hauptlauf; der Textnachzug in §3.2 ist umgesetzt.

Zähler (abgeleitet): **1×** (evidence/slice-routing-sdk-realserver-e2e.md).
