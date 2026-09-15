# Beleg: slice-078

Vorgang: `slice-078` — die Aktivierung des `hostpaths`-Moduls und die
Festlegung seiner Reichweite.

Fund: `AGENTS.md` §3.11 deckt seit `ADR-0075` die **ganze Markdown-Fläche
einschließlich der Fenced-Blöcke**, während das Modul `hostpaths` Fences per
Design frei lässt und Nicht-`.md`-Dateien gar nicht liest. Gemessen: ein
Host-Pfad in einem Fence lässt das Modul **grün**, während die Regel ihn
verbietet; der repo-weite Grep findet ihn. Die Lücke ist an beiden Trägern
benannt und der Wächter dort ist das Review — sie ist damit kein stiller
Zustand, aber auch **kein Gate**.

Quelle: `docs/reviews/verify-slice-078.md` (Mutation M-1) ·
`docs/reviews/review-slice-078-delta.md` (Mutation A/B) ·
`docs/plan/adr/0075-hostpaths-reichweite-und-wortlaut.md` ·
`AGENTS.md` §3.11 · `harness/sensors/docs-check.md` §Grenze Punkt 8.
