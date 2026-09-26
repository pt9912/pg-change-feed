# `make suchlauf-nachmessen` — misst das Suchlauf-Feld eines Slice-Plans nach

## Vertrag

`make suchlauf-nachmessen PLAN=<Plan-Datei>` liest die Codeblöcke mit dem Etikett
`suchlauf` aus einem Slice-Plan, führt jede Zeile als `git grep` aus, zählt die
Trefferzeilen und vergleicht sie mit der im Plan genannten Zahl
(`tools/harness/suchlauf-nachmessen.sh`; netzlos, nur `bash` und `git`, kein
Docker). Es wiederholt eine **vom Plan deklarierte Messung** — Befehl, Stand,
Zahl —, es entscheidet nicht, ob die Messung die richtige ist
([`AGENTS.md`](../../AGENTS.md) §3.13 §Suchform,
[`ADR-0083`](../../docs/plan/adr/0083-herkunft-von-aussagen-in-traegern.md)).

**Kein Gate.** Der Stand `diff` bewegt sich mit jedem Commit: eine Zahl, die
heute stimmt, ist morgen mit Recht eine andere. Die Messung gehört in den Lauf
dessen, der sie braucht — Implementer nach jeder Fixrunde, Reviewer,
Verifier, Planner. Das Ziel steht in keinem Gate-Bündel (`make gates`).

## Form der Zeile

Ein Codeblock je Plan-Feld mit dem Etikett `suchlauf`, eine Zeile je Messung:

```text
<Stand> <Soll> <Argumente von git grep>
```

- **`<Stand>`** ist eine Commit-Kennung (7 bis 40 Hexziffern; der Parent, nie
  `HEAD` und kein Branch- oder Tag-Name) oder das Wort `diff` — der Arbeitsbaum
  beim Nachmessen (tracked Dateien, wie `git grep` sie sieht; neue Dateien
  vorher `git add`).
- **`<Soll>`** ist die im Plan genannte Zahl der Trefferzeilen.
- **`<Argumente von git grep>`** sind Optionen und Suchmuster; ein alleinstehendes
  `--` trennt sie vom Pathspec. Ohne Pathspec ist der Suchraum der ganze Baum.
  `'…'` und `"…"` gruppieren Wörter (ohne Expansion, ohne Backslash-Escape),
  das Muster steht also so im Block, wie es in der Shell stünde.

Das Werkzeug ruft `git grep -n <Optionen> [<Stand>] -- <Pathspec> <Ausschluss>`
auf. Der **Ausschluss** ist die Plan-Datei selbst unter ihrem Dateinamen in
jedem Verzeichnis (`:(exclude,glob)**/<Dateiname>`): der Plan trägt Muster und
Beschreibung der Suche, und er liegt am Parent-Stand oft unter einem anderen
Lifecycle-Verzeichnis. Gezählt werden die gedruckten Trefferzeilen; Optionen, die
die Ausgabeform ändern (`-c`, `-l`, `-A`), ändern damit auch, was die Zahl
zählt.

Eine Tabellenzelle trägt die Zeile nicht: das Escape `\|` aus einer Zelle liefert
als `-E`-Muster keinen Treffer. Bestehende Pläne tragen keinen solchen Block; ein
Plan trägt ihn, sobald sein Implementer das Feld füllt.

## Ausgabe und Ausgänge

Je Zeile eine Ausgabezeile `OK` oder `ABWEICHUNG` mit Soll, Ist und der
Plan-Zeile:

```text
ABWEICHUNG  soll=3 ist=2  a32a1931 3 -E 'zeichengenau' -- spec
```

| Exit | Bedeutung |
|---|---|
| 0 | jede Zeile stimmt |
| 1 | mindestens eine Abweichung zwischen Soll und Ist |
| 2 | Eingabefehler, bevor ein Befehl läuft: kein `suchlauf`-Block oder ein leerer, ein nicht geschlossener Block, ein Stand außerhalb der Form (`HEAD`, Name, unbekannte Kennung), ein Soll ohne Zahl, eine Zeile ohne Muster; oder `git grep` selbst endet mit einem Fehler (Exit über 1) |

**Leer ist nicht bestanden:** ein Plan ohne `suchlauf`-Zeile endet mit Exit 2,
nicht mit Exit 0. Über `make` kommt jeder Exit ≠ 0 als der Make-eigene Exit `2`
an.

## Grenze

Das Werkzeug prüft **Zahlen und Stände**, nicht die **Vollständigkeit** von
Suchraum und Suchmuster. Ob das Muster den Symbolnamen, das Zählwort und die
Beschreibung der bewegten Eigenschaft trägt und ob der Suchraum den ganzen Baum
deckt, bleibt eine Lese-Handlung des Reviewers ([`AGENTS.md`](../../AGENTS.md)
§3.13). Ein grüner Lauf sagt: die Zahlen im Plan sind die Zahlen dieser Befehle an
diesen Ständen — nicht, dass das Feld vollständig ist. Ein Sensor „jede Zahl trägt
ihren Ursprung“ bleibt ausgeschlossen
([`ADR-0083`](../../docs/plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
§Entscheidung 4).

Die Zahl für `diff` gilt für den Arbeitsbaum des Laufs; der Aufruf-Beleg im
Bericht nennt Kommando und Zeitpunkt.

## Test

`make test-suchlauf-nachmessen` (`tools/harness/run-suchlauf-nachmessen-tests.sh`,
netzlos) fährt fünf Fälle gegen ein Wegwerf-Repo: stimmt · weicht ab · Selbstverweis
ausgeschlossen (der Plan trägt das Suchwort selbst, am Parent-Stand unter einem
anderen Verzeichnis) · `HEAD` abgelehnt · kein Block (Tabelle, `text`-Block,
leerer Block).
