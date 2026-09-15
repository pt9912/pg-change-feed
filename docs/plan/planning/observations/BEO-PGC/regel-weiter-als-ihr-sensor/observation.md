# BEO-PGC/regel-weiter-als-ihr-sensor

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft das Verhältnis
zwischen der Reichweite einer Zusage und dem Umfang dessen, was ihr Sensor
prüfen kann, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Eine Regel oder Zusage ist **weiter gefasst als ihr Sensor**.
Die Lücke ist meist gewollt und benannt — Beispiele im Fenced-Block, eine
Disziplin-Regel, ein Werkzeug, das nur `.md` liest —, aber sie wird nur gedeckt,
solange die benannten Stellen mitgepflegt werden, und in **keinem** Gate fällt
sie auf. Belegt an `slice-078`: `AGENTS.md` §3.11 deckt seit `ADR-0075` die
ganze Markdown-Fläche einschließlich der Fenced-Blöcke, während das Modul
`hostpaths` Fences per Design frei lässt und Nicht-`.md`-Dateien nicht liest —
gemessen: ein Host-Pfad im Fence lässt das Modul grün, die Regel verbietet ihn.
