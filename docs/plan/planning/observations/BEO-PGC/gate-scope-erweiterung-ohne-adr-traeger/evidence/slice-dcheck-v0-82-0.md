**Vorgang:** slice-dcheck-v0-82-0 (Review F-3, MEDIUM)

**Fund:** Der Pin-Commit `7f796ef0` hob `d-check.mk` auf d-check v0.82.0. Ab v0.80.0 meldet das
Modul `hostpaths` auch Pfade unter der Tilde; `make docs-check` (in `make gates`) blockiert damit
eine Pfad-Klasse, die weder `AGENTS.md` §3.11 noch `ADR-0075` Entscheidung 2 verbot. Der Plan
behandelte das als Meldung an den Planner; die Reichweite der Regel ist aber eine Frage an eine
`Accepted`-ADR. Aufgelöst über `ADR-0160` Entscheidung 3 und 4 (Teil-Supersede von `ADR-0075`
Entscheidung 2, §3.11 in Fassung 3).

**Form (Ausprägung):** dieselbe Klasse wie die beiden ersten Belege (Gate-Scope wächst ohne den
Träger, den die zuständige ADR benennt), neuer Weg: nicht die Konfigurationsdatei des Gates
(`.a-check.yml`, `.d-check.yml`), sondern der **Werkzeug-Pin** trägt die Erweiterung — ein
Pin-Update sieht aus wie Pflege (`ADR-0051` Entscheidung 7) und ändert doch, was das Gate
blockiert. Vor dem Merge vom Reviewer gefunden. Drittes Auftreten.

Quelle: `docs/reviews/review-slice-dcheck-v0-82-0.md` (F-3). <!-- d-check:status-provenance -->
