Zustand: offen — Ausgang noch nicht zugewiesen. Erstes Auftreten, weit unter
der 3×-Schwelle. Real bestätigt und behoben an einer Stelle
(`docs/reviews/review-slice-sdk-csharp-pack-werkzeug.md`), aber ohne
nachgewiesene Sensor-Wirkung in diesem Fall (`make gates`/`make docs-check`
liefen über den gesamten Zeitraum grün, siehe `observation.md`) — deshalb
kein Carveout, keine Gate-Änderung, nur der Beleg.

Kandidat für eine künftige Schärfung, falls ein zweites Auftreten mit
tatsächlicher Sensor-Wirkung (ein `id-unlinked`-Fehlalarm oder ein
übersehener echter Fund wegen verschobener Parität) hinzukommt: eine
zusätzliche Lese-Pflicht im Reviewer-Skill (`.harness/skills/reviewer.md`),
bei einem Fund der Klasse „scheinbar nackte Kennung" zusätzlich die
Backtick-Gesamtparität des Dokuments zu prüfen (z. B. Gesamtzahl gerade),
bevor die Kennung selbst als Fehlstelle behandelt wird.

Zähler (abgeleitet): 1× (evidence/slice-sdk-python-pack-werkzeug.md).
