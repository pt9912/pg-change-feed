Zustand: offen — Ausgang noch nicht zugewiesen. Erstes Auftreten, unter der
3×-Schwelle für eine reguläre Skill-Verkörperung — aber mit **real
nachgewiesener Sensor-Wirkung**, nicht nur hypothetisch: Die reine
Backtick-Paritäts-Korrektur an
`docs/reviews/review-slice-sdk-csharp-pack-werkzeug.md` <!-- d-check:status-provenance -->
(Commit `69d2dc00`) legte einen zweiten, bis dahin vom Backtick-Defekt
verdeckten `id-unlinked`-Befund (nackte `ADR-0106`-Erwähnung, Zeile 132)
offen, der über den gesamten bisherigen Merge-Zeitraum trotz wiederholt
grünem `make gates`/`make docs-check` real unentdeckt blieb. Beide Funde
sind in eigenen, separat committeten Fixes behoben.

Kandidat für eine künftige Schärfung — Vorschlag an einen künftigen
Architect-Zug, unabhängig von der 3×-Schwelle, weil der Wirkmechanismus
bereits an einem realen Dokument bestätigt ist (kein reiner Verdacht mehr):
eine zusätzliche Lese-Pflicht im Reviewer-Skill
(`.harness/skills/reviewer.md`), bei jedem Fund der Klasse „scheinbar
nackte Kennung" zusätzlich die Backtick-Gesamtparität des ganzen Dokuments
zu prüfen (Gesamtzahl gerade), bevor die Kennung selbst als alleinige
Fehlstelle behandelt wird — sonst bleibt ein zweiter, durch denselben
Defekt verdeckter Fund unentdeckt.

Zähler (abgeleitet): 1× (evidence/slice-sdk-python-pack-werkzeug.md) — ein
Vorgang, zwei zusammenhängende Funde (Paritäts-Defekt + dadurch verdeckter
`id-unlinked`-Befund).
