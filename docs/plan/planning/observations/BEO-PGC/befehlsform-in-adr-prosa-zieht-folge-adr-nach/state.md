Zustand: **offen** — unter der Schwelle. Der Ausweg, den `ADR-0159` selbst
vorsieht (Re-Evaluierungs-Trigger (a)), ist für die Referent-Messung geliefert:
`make zitat-vergleich` mit Tabellentest `make test-zitat-vergleich`, Vertrag
`harness/targets/zitat-vergleich.md` · seit slice-zitat-vergleich-werkzeug
(Beleg-Anker: `git grep -n 'seit slice-zitat-vergleich-werkzeug' -- harness/README.md`).
Der `bash`-Block in `ADR-0159` Entscheidung 4 ist die historische Fassung; die
Randfälle der Messung gehen ins Skript und in den Tabellentest, nicht in eine
Folge-ADR. Die Klasse bleibt offen, weil sie jede ADR mit ausführbarer
Befehlsform betrifft: die Datei-Schleife des MR-Datei-`cmp` aus `ADR-0157`
Entscheidung 4 steht weiter als Befehlsform in einer ADR (Vertrag
§Grenzen, „Was nicht im Werkzeug steht“).

Zähler (abgeleitet): **1×** (evidence/slice-zitat-korrektur-vergleichseinheit.md).
