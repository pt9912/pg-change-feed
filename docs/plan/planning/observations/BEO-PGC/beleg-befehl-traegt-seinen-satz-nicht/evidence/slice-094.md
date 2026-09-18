# Beleg: slice-094

Vorgang: `slice-094` — Coverage Cluster A (Prozess-Rand).

Fund: In vier Argument-Fehler-Fällen des neuen Prozess-Rand-Tests band die
**Ausgabe-Hälfte** nicht — der Test prüfte `Contains(stderr, <Modus-Name>)`,
und der Modus-Name steht **auch** in der **generischen** Fallback-Zeile
(`cmd/pg-change-feed/main.go:104`, die alle Modi aufzählt). Der genannte Beleg
trug seinen Satz also nur zur Hälfte: der **Exit-Code** war gebunden, die
**Zeile** nicht.

Der Reviewer hat es mit genau der Mutation gezeigt, die den Träger entfernt: die
moduseigene Meldung (`main.go:73`) durch den Fallback-Text ersetzt, Ausgang
bleibt 2 — die Suite war **grün**. Nach der Fixrunde färbt dieselbe Mutation
**rot** (EC 1), und der Delta-Review hat den **Übergang** unabhängig
nachgefahren: derselbe Mutationsstand mit dem **alten** Testfile ist grün, mit
dem neuen rot.

**Was diesen Beleg von den vier früheren unterscheidet:** dort war der Beleg ein
**Befehl** (`slice-084`, `-085`), ein **Testkommentar** (`slice-091`) oder eine
**Adresse** (`slice-093`); hier ist es eine **Assertion**. Vier Formen, ein
Kern: der Leser kann nicht nachschlagen, was der Satz ihm zu prüfen gibt — er
prüft etwas, das auch woanders herkommt.

**Und eine zweite Grenze desselben Slice gehört hierher, obwohl sie nicht
dieselbe Klasse ist:** eine Mutation, die die `GOCOVERDIR`-Weitergabe an die
Kindprozesse entfernt, lässt **alle** Tests grün und `cmd` auf **0 von 49**
zurückfallen (Gesamt 80,50 %). Die Coverage-Zusage dieses Slice ist damit
**nicht testgewahrt** — sie steht als **§Grenze 7** im Sensor-Dokument, mit
„Wächter: keiner". Benannt, nicht still.

Quelle: Review zu `slice-094` (F-2) ·
Delta-Review zu `slice-094` (Negativbefunde, Messungen 4 und 5) ·
`cmd/pg-change-feed/main_test.go` (berichtigt in `8292766`) ·
`harness/sensors/coverage-gate.md` §Grenze Punkt 7.
