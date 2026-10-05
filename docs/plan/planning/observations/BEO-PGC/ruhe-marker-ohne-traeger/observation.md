# BEO-PGC/ruhe-marker-ohne-traeger

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Lifecycle-Disziplin
der Roadmap unter `docs/plan/planning/`, keine eigene Sub-Area im Sinn der
Modus-Deklaration).

Die Beobachtung: Der Ruhe-Marker unter *Offene Wellen* der Roadmap steht genau dann,
wenn `in-progress/` keinen Slice trägt (Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur, Bullet *Offene Wellen*: „in **beide** Richtungen“). Im Repo hat
das Setzen und Entfernen keinen Träger: kein Schritt unter `.claude/commands/` nennt
es, und kein Sensor hält den Marker gegen das Verzeichnis. Gemessen im Lauf von
`slice-baseline-6-14-0-dokumente-nachziehen` (Befund 4): der Marker stand
unverändert seit `d38cf9e7` (2026-09-23, `git blame`), über alle seither
beanspruchten Slices hinweg; bis `f4112fe9` wurde er je Slice per Commit gesetzt.
Der Slice entfernte ihn und setzte ihn bei seiner Closure wieder ein — als eigene
DoD-Pflicht, weil der Review (F-3, LOW) die Übergabe nur in Belege-Prosa fand.

**Warum das zählt:** Der Marker ist die deklarierte Redundanz der Roadmap zum
Verzeichnis. Ohne Träger driftet er still, und die Antwort auf „was läuft gerade“
aus der Roadmap ist falsch, während `ls in-progress/` richtig antwortet.

**Abgrenzung.** `BEO-PGC/roadmap-kontext-verwaist-nach-zeilen-entfernung` betrifft
Text, der nach einer Umplanung stehen bleibt; hier fehlt der Träger eines
Zustandsübergangs. `BEO-PGC/regel-weiter-als-ihr-sensor` betrifft eine benannte
Lücke zwischen Regel und Sensor; hier war die Lücke nicht benannt.
