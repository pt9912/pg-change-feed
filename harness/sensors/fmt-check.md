# `make fmt-check` — meldet Go-Dateien, die `gofmt` nicht als formatiert führt

## Vertrag

`make fmt-check` führt `gofmt -l` über **alle** Go-Dateien unter der Repo-Wurzel
im gepinnten Toolchain-Image aus (`TOOLCHAIN_IMAGE`, `docker run --network none`,
das Verzeichnis lesend gemountet; `tools/harness/fmt-check.sh`), druckt die
abweichenden Dateien (ein Pfad je Zeile, relativ zur Wurzel) und endet mit dem
Exit-Code der Tabelle unten. Es schreibt nichts: es gibt keinen `gofmt -w`-Pfad
und keinen `--user`-Workaround. Eine gemeldete Datei wird nach der Ausgabe von
`gofmt -d` (im selben Image, lesend) mit dem Edit-Werkzeug des Laufs korrigiert,
nie mit einem in-place schreibenden Textwerkzeug
([`AGENTS.md`](../../AGENTS.md) §3.1).

`gofmt -l` endet auch bei Abweichung mit Exit 0. Das Werkzeug wertet deshalb
die **Ausgabe** aus, nicht den Exit-Code des Formatierers; der Exit-Code des
Formatierers zählt nur, wenn er einen Fehler meldet (Syntaxfehler).

**Kein Gate.** Das Ziel steht in keinem Gate-Bündel (`make gates`). Die Aufnahme
als Gate braucht eine ADR ([`AGENTS.md`](../../AGENTS.md) §3.6, §4). Der Grund für
„Werkzeug ohne Gate“ steht im Register-Eintrag
`BEO-PGC/formatierungs-drift-ohne-gate` (`state.md`): der Fund kam in allen drei
Vorgängen von einem Leser vor dem Merge, Schwere LOW, je eine kleine Zahl Dateien.
**Trigger der Gate-Aufnahme (Kenntnis):** ein weiteres Auftreten, das der Reviewer
trotz gelaufenem Schritt 18 des Implementer-Ablaufs findet; dann Architect-Frage
mit ADR-Vorschlag.

## Aufruf

```text
make fmt-check
tools/harness/fmt-check.sh [Verzeichnis]     # Verzeichnis-Argument für den Test
```

Das Skript ruft `docker run --rm --network none -v <Verzeichnis>:/src:ro -w /src
<TOOLCHAIN_IMAGE> sh -c …` auf. Im Container zählt es zuerst die Go-Dateien
(reguläre Dateien mit Endung `.go`, deren Name nicht mit einem Punkt beginnt —
dieselben, die `gofmt` liest), dann läuft `gofmt -l .` über den ganzen Baum. Auf
stdout stehen die abweichenden Pfade; die Zählung steht als eine Zeile
`fmt-check: <n> Go-Dateien geprüft, …` (bei Abweichung auf stderr, sonst auf
stdout).

Host-Werkzeuge: `bash`, `git` (Repo-Wurzel) und `realpath` sowie `docker`
([`AGENTS.md`](../../AGENTS.md) §3.1, Klasse „Host-Werkzeug ohne Installation“);
`gofmt` läuft im Container, nicht auf dem Host.

## Exit-Codes

| Exit | Bedeutung |
|---|---|
| 0 | jede Go-Datei ist formatiert |
| 1 | mindestens eine Datei weicht ab (die Pfade stehen auf stdout) |
| 2 | Eingabe- oder Formatierer-Fehler: mehr als ein Argument, kein Verzeichnis, `TOOLCHAIN_IMAGE` fehlt, Syntaxfehler in einer Go-Datei, Docker-Fehler, oder **keine Go-Datei** unter dem Verzeichnis |

**Leer ist nicht bestanden:** ein Verzeichnis ohne Go-Datei endet mit Exit 2, nicht
mit Exit 0 — ein falsch gemounteter Pfad meldete sonst grün. Über `make` kommt
jeder Exit ≠ 0 als der Make-eigene Exit `2` an; die Unterscheidung von 1 und 2
trägt die Ausgabe (Pfade bzw. Meldung) oder der direkte Aufruf des Skripts.

## Wer es aufruft

- **Implementer**, Schritt 18 (`.claude/commands/implement-slice.md`, Absatz
  „Format“): vor der „fertig“-Meldung und nach jeder Fixrunde; die gedruckte Zeile
  steht im Bericht.
- **Reviewer** (`.harness/skills/reviewer.md`, Punkt LOW): derselbe Lauf als Probe.

## Grenze

1. **Formatierung, nicht Semantik.** Das Werkzeug sagt, ob `gofmt` die Datei
   umschriebe; ob der Code richtig ist, sagt es nicht. Kein Lint, kein Import-
   Sortieren jenseits dessen, was `gofmt` selbst tut ([`AGENTS.md`](../../AGENTS.md)
   §3.2: das Repo führt keinen Linter).
2. **`gofmt` formt Doc-Kommentare um.** Ein Paar `''` oder ein Backtick-Paar in
   einem Doc-Kommentar wird zu einem typografischen Anführungszeichen; wer den
   Kommentar so meint, schreibt ihn anders, der Inhalt ändert sich sonst mit
   der Formatierung.
3. **Nur Go.** Markdown, YAML, Shell, C#, Kotlin und Python haben in diesem Repo
   keinen Formatierer im Werkzeugsatz; erzeugter Code (`gen/`, `sdks/`) liegt im
   Suchraum, sofern er Go ist, und muss `gofmt`-fest sein.
4. **Grün ist kein Beleg für den Diff.** Der Lauf misst den Baum, nicht den
   eigenen Diff: eine Datei außerhalb des Diffs, die abweicht, färbt ihn rot.

## Test

`make test-fmt-check` (`tools/harness/run-fmt-check-tests.sh`, netzlos) fährt
echte Docker-Läufe gegen Wegwerf-Verzeichnisse im Temp-Verzeichnis. Je Zusage die
Mutation ihrer Eingabeseite, die den Fall rot färbt:

| Zusage | Mutation am Werkzeug | Fall, der rot wird |
|---|---|---|
| Abweichung heißt Exit 1, die Ausgabe wertet | nur der Exit-Code des Formatierers gilt (Auswertung der Ausgabe entfernt) | unformatierte Datei · gemischt · Unterverzeichnis · Eingabe bleibt · Pfad mit Leerzeichen |
| der ganze Baum wird gelesen | `gofmt -l ./*.go` statt `gofmt -l .` | Unterverzeichnis |
| leer ist nicht bestanden | Zählung der Go-Dateien entfernt | leeres Verzeichnis · Verzeichnis ohne Go-Datei |
| dieselben Dateien wie `gofmt` | die Zählung nimmt auch Dateien mit Punkt am Namensanfang | Verzeichnis ohne Go-Datei |
| Syntaxfehler heißt Exit 2 | Exit-Code des Formatierers ignoriert | Syntaxfehler |
| Docker-Fehler heißt Exit 2 | der Exit des Docker-Aufrufs wird durchgereicht (125) | Docker-Fehler |
| nur die abweichende Datei wird genannt | die Ausgabe trägt zusätzlich jede Datei des Verzeichnisses | gemischt · formatierte Datei · Syntaxfehler |
| lesender Mount, kein Schreibpfad | `gofmt -l -w` ohne `:ro` | Eingabe bleibt byte-gleich |
| ein Verzeichnis-Argument, ein Mount | `"$dir"` ohne Anführungszeichen | Pfad mit Leerzeichen |
| mehr als ein Argument ist ein Eingabefehler | die Argumentzahl-Prüfung entfernt | zu viele Argumente |
| ein fehlendes Verzeichnis ist ein Eingabefehler | die Verzeichnis-Prüfung entfernt | kein Verzeichnis |
| eine fehlende `TOOLCHAIN_IMAGE` ist ein Eingabefehler | die Prüfung der Variable entfernt | `TOOLCHAIN_IMAGE` fehlt |
