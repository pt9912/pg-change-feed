# BEO-PGC/exemption-ohne-reifegrenze

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft das
Verhältnis zwischen einer Sensor-Ausnahme und dem Zeitpunkt, ab dem sie
gilt, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Eine Sensor-Ausnahme (`.d-check.yml` `versions.exempt-paths`
o. ä.) wird mit dem Vorbild einer **bestehenden** Ausnahme begründet, ist
aber tatsächlich weiter gefasst: Das Vorbild schützt erst **ab einem
Abschluss-Ereignis** (`git mv` nach `done/`), die neue Ausnahme schützt
**ab dem ersten Commit**, der die Datei anlegt. Eine Datei, die während
ihrer aktiven Phase einen inzwischen veralteten Pin zitiert — per Tippfehler
oder Copy-Paste —, wird dadurch nicht mehr rot markiert, obwohl sie es beim
engeren Vorbild noch würde. Belegt an `slice-105`: `docs/reviews/**` erhielt
eine `versions.exempt-paths`-Ausnahme „mit derselben Begründung wie
`harness/conventions/done/**`" — real reproduziert bleibt ein stale
`v6.5.0`-Pfad in `docs/reviews/**` grün, derselbe Pin außerhalb der Ausnahme
(`harness/README.md`) färbt korrekt rot.
