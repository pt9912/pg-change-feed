# BEO-PGC/arbeitsbaum-race-ohne-offenlegung

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft das
Verhältnis von Commit-Text zu committem Inhalt bei mehreren, zugleich in
einem Arbeitsbaum laufenden Slices).

Die Beobachtung: Laufen mehrere Slices zugleich in einem Arbeitsbaum,
hebt ein Commit geänderte Dateien mit hoch, die zum Gegenstand eines
**anderen** Slice gehören — inhaltlich korrekt, aber der Commit-Text
nennt nur den eigenen Slice. Der Commit-Betreff behauptet einen Umfang,
den der Diff (neben dem eigenen) noch einen fremden trägt; die
Herkunft des fremden Anteils bleibt nur über die Berichte der
beteiligten Slices rekonstruierbar, nicht über `git log`.

Konkret (Kotlin-Review F-4, INFO): der Diff von `b239d849` (Betreff
nennt ausschließlich den Kotlin-Slice) hob
`docs/user/benutzerhandbuch.md` von Version 1.79 auf 1.81 und trug
damit **beide** Changelog-Zeilen — 1.80 (Python-SDK, eigener Gegenstand
des parallel laufenden Python-Slices) und 1.81 (Kotlin). Der
nachfolgende Python-Feature-Commit `12403d9f` bestätigt die Übernahme
im eigenen Commit-Text; `git diff b239d849 12403d9f --
docs/user/benutzerhandbuch.md` ist leer. Der `b239d849`-Commit-Text
selbst benennt die Fremdübernahme nicht.

**Warum das zählt:** Kein Inhalt war falsch, und die Traceability-Kennung
(`LH-FA-SST-009`) ist für beide Slices identisch — aber ein späterer
Leser von `git show b239d849` ordnet die Python-Zeile dem Kotlin-Slice
zu, weil der Commit-Text sie nicht ausspart und nicht erklärt. Die
billige Gegenmaßnahme ist die Offenlegung im Commit-Text („hebt
ebenfalls mit hoch: <Datei> aus <Slice>, läuft parallel") oder — wenn
zeilenscharf möglich — die Trennung der Commits; dieselbe Race-Klasse
ist bereits zweimal in derselben Welle transparent berichtet worden
(Ein-Zeilen-Link-Fix im C#-Verifikations-Report, Handbuch-Verschnitt
des Python-Zugs), hier erstmals als Befund ohne Offenlegung im
Commit-Text.
