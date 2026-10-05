# `make kommentar-kennungen` — listet Kommentarblöcke, die ihre Herkunft nicht als ein Feld tragen

## Vertrag

Das Werkzeug prüft die **Form**, nicht die **Wahrheit**: es zählt, wie viele
verschiedene Kennungen ein Kommentarblock trägt. Ein Lauf ohne Kandidat sagt
nicht „die Kommentare sind konform“ — eine Spec- oder ADR-Aussage in eigenen
Worten hinter einer einzigen Kennung erkennt es nicht, ebenso wenig einen
Kommentar, der nicht zutrifft oder keine der fünf Kommentar-Klassen trägt. Das
sind Lese-Handlungen des Reviewers
([`AGENTS.md`](../../AGENTS.md) §3.7, §3.12; Baseline-Regelwerk
`grundlagen-harness-dateien.md` §Was ein Kommentar trägt).

Gegenstand ist die Regel „Herkunft im Kommentar ist ein auflösbares Feld“
([`AGENTS.md`](../../AGENTS.md) §3.7): höchstens **eine** Kennung je Kommentar,
keine Kette, keine Kompaktform, kein „ff.“. `make kommentar-kennungen` liest die
Kommentarblöcke der `.go`-Dateien und der Nicht-Go-Zeilenkommentar-Formen
(`.sh`, `.mk`, `.yml`, `.yaml`, `.sql`, Makefile, Dockerfile) und meldet jeden,
der sie verletzt (**Kandidat**). Die Herkunft der Regel und die Grenze „kein
Sensor über Prosa“ führt
[`ADR-0083`](../../docs/plan/adr/0083-herkunft-von-aussagen-in-traegern.md).

**Kein Gate.** Das Ziel steht in keinem Gate-Bündel (`make gates`) und hat
**keinen Ausnahme-Pfad** (keine Marker im Quelltext, keine Ausnahmeliste,
[`AGENTS.md`](../../AGENTS.md) §3.2): ein Lauf über den Bestand liefert
Kandidaten, bis der Bestand bereinigt ist, und die Aufnahme in `make gates`
brauchte eine ADR ([`AGENTS.md`](../../AGENTS.md) §3.6).

**Abgrenzung zum Chronik-Kandidatenlauf.** Das Werkzeug zählt verschiedene
Kennungen je Kommentarblock. Es liest weder ein Satz-Subjekt noch eine
Slice-/Wellen-Nummer und unterscheidet nicht „Testfall-Provenienz“ von
„Produktionsverhalten-Chronik“; das ist ein Urteil über das Satz-Subjekt, für
das der Architect-Verdikt
[`architect-verdict-slice-chronik-in-code-kommentar`](../../docs/reviews/architect-verdict-slice-chronik-in-code-kommentar.md)
einen Textmuster-Sensor ausschließt. Der Chronik-Kandidatenlauf in
`.claude/commands/implement-slice.md` Schritt 20 bleibt daneben unverändert.

## Kandidat

Ein **Kommentarblock** ist je Form definiert:

- **Go (`.go`):** eine Kommentargruppe des Go-Parsers (`go/parser`,
  `ParseComments`): eine Folge von Kommentaren ohne Leerzeile und ohne Code
  dazwischen; ein Endkommentar hinter Code bildet einen eigenen Block.
  Direktiven (`//go:…`) zählen nicht mit. Eine Kennung in einem
  Zeichenketten-Literal ist kein Kommentar und wird nicht gelesen.
- **Zeilenkommentare (`#` in `.sh`/`.mk`/`.yml`/`.yaml`/Makefile/Dockerfile,
  `--` in `.sql`):** eine Folge aufeinanderfolgender vollzeiliger
  Kommentarzeilen (erster Nicht-Leerraum-Text beginnt mit dem Marker); jede
  andere Zeile — auch eine Leerzeile — beendet den Block, und eine nur aus dem
  Marker bestehende Zeile (`#`, `--`) ist Grenz-Marker. Nachgestellte
  Kommentaranteile hinter Code liest die Messung nicht (die Trennstelle ist
  in diesen Formen syntaktisch mehrdeutig).

Nicht im Messraum: Markdown (Prosa-Träger ohne Kommentarform; die
Kennungs-Linkpflicht in Prosa trägt d-check `ids`), `gen/**` und `sdks/**`
(excludedRoots) sowie `.harness/**` (vendored Baseline, SHA-gepinnt).

Eine **Kennung** hat eine von vier Arten: `ADR-NNNN`, `LH-FA-XXX-NNN` oder
`LH-QA-XXX-NNN`, `SPEC-NNN`, `ARC-NNN`. Ein Block ist Kandidat, wenn er

- mindestens **zwei verschiedene** Kennungen trägt (dieselbe Kennung zweimal
  zählt einmal; eine Unterkennung wie `LH-FA-CAP-006.a` zählt als ihre
  Kennung), oder
- „ff.“ hinter einer Kennung trägt (auch über einen Zeilenumbruch hinweg).

Eine **Kompaktform** zählt je Nummer als eine Kennung: `LH-FA-CFG-001/002/003`
sind drei, `LH-QA-SEC-001…003` sind zwei (die Endpunkte des Bereichs, die
Zwischennummer nennt der Text nicht). Die Fortsetzung hat die Ziffernbreite der
Kennung davor; `ADR-0060/005` setzt keine Kennung fort. Formen, die Kennungen
anders verbinden (`bis`, Bindestrich-Bereich), erkennt das Werkzeug nicht.

## Aufruf

```text
make kommentar-kennungen [PATHS=<Pfade>] [COUNT=1] [TESTS=exclude|only] [DIFF=<Basis>]
```

| Variable | Wirkung |
|---|---|
| `PATHS` | Leerzeichen-getrennte Pfade (Dateien oder Verzeichnisse) relativ zur Repo-Wurzel; Standard: der Baum. Ausgenommen sind `gen/`, `sdks/`, `.git`, `.harness` (erzeugter Code; die SDK-Bäume prüft `make sdk-public-doc-check` strenger — dort ist jede Kennung verboten; Repository-Inneres; vendored Baseline). |
| `COUNT=1` | druckt nur die Zahl der Kandidaten, Exit 0 |
| `TESTS=exclude` / `TESTS=only` | nur Nicht-Test-Dateien bzw. nur `*_test.go`; ohne Wert beide |
| `DIFF=<Basis>` | nur Blöcke, die eine seit `<Basis>` **hinzugefügte** Zeile überlappen (Eingabe: `git diff -U0 <Basis> -- '*.go' '*.sh' '*.mk' '*.yml' '*.yaml' '*.sql' 'Makefile' 'Dockerfile' '**/Dockerfile'` — dieselbe Pfadspec wie die `tools/harness/kommentar-kennungen.sh` setzt); ein Block, den der Diff nur berührt, ohne eine Zeile zu ändern, und eine reine Löschung zählen nicht. Bestandskandidaten färben den Lauf eines Implementers nicht. Neue Dateien vorher `git add` (der Diff sieht nur getrackte Dateien) |

Ausgabe je Kandidat: `Datei:Zeile-Zeile  Kennungen` (bei „ff.“ steht es am
Ende). Die Kandidaten stehen nach Dateipfad und Zeile sortiert.

Das Programm `tools/harness/kommentar-kennungen/` (nur Standardbibliothek)
läuft im gepinnten Toolchain-Image (`TOOLCHAIN_IMAGE`, `--network none`, Repo
lesend gemountet); `tools/harness/kommentar-kennungen.sh` ruft es auf. Host-
Werkzeuge: `bash`, `git` (für `DIFF`), `docker` — dieselbe Klasse wie bei
`make suchlauf-nachmessen`; installiert wird nichts. Der Diff-Strom entsteht auf
dem Host in einer Temp-Datei, nicht in einer Pipe: `pipefail` meldet den Exit
des rechten Glieds und ließe ein fehlgeschlagenes `git` hinter Exit 1 des
Programms verschwinden ([`AGENTS.md`](../../AGENTS.md) §3.9).

Der Aufrufer pinnt die Form des Stroms gegen die Git-Konfiguration des Nutzers:
`git -c core.quotePath=false diff -U0 --no-color --no-ext-diff --no-textconv
--src-prefix=a/ --dst-prefix=b/` (ohne sie ändern `diff.mnemonicPrefix`,
`diff.noprefix`, `color.diff` oder `diff.external` den Strom). Das Programm liest
die Zieldatei-Zeilen mit dem Präfix `b/` und endet bei einer `+++`-Zeile mit
anderem Präfix mit Exit 2, statt „kein Kandidat“ zu melden; Zeilen innerhalb
eines Hunks (nach der Zahl im Hunk-Kopf) sind Inhalt, auch wenn sie mit `+++ `
beginnen.

## Exit-Codes

| Exit | Bedeutung |
|---|---|
| 0 | kein Kandidat, oder `COUNT=1` (die Zahl steht auf stdout) |
| 1 | mindestens ein Kandidat |
| 2 | Eingabefehler: unbekannter Wert für `TESTS`, `DIFF` ist kein Commit, ein Pfad fehlt, eine `.go`-Datei ist nicht lesbar, der Diff-Strom nicht parsbar oder eine Zieldatei-Zeile trägt nicht das Präfix `b/` |

Über `make` kommt jeder Exit ≠ 0 als der Make-eigene Exit `2` an; die Unterscheidung
von 1 und 2 trägt die Make-Meldung `Fehler <n>` (der Exit des Skripts), die Ausgabe
(Kandidatenzeilen bzw. Meldung auf stderr) oder der direkte Aufruf von
`tools/harness/kommentar-kennungen.sh`.

## Wer es aufruft

- **Implementer**, Schritt 20 (`.claude/commands/implement-slice.md`): vor der
  Übergabe `make kommentar-kennungen DIFF=<Basis>` über die eigenen Änderungen,
  jeder Kandidat wird umformuliert.
- **Reviewer** (`.harness/skills/reviewer.md`): derselbe Lauf als **Probe**, nicht
  als Beleg — er liest die Kommentare des Diffs zusätzlich auf Spec-Wiedergabe.
- **Bereinigung des Bestands:** die Kandidatenzahl (`COUNT=1`) ist ihre Messung.

## Grenze

1. **Form, nicht Wahrheit** (erster Absatz).
2. **Nur Code- und Konfig-Kommentare, keine Prosa.** Markdown liest das Werkzeug
   nicht — Prosa ist selbst der Träger; die Kennungs-Linkpflicht dort trägt
   d-check `ids`. `sdks/` deckt `make sdk-public-doc-check`, `gen/**` und
   `.harness/**` sind ausgenommen. Zeilenformen lesen vollzeilige Kommentare
   nach Konvention, nicht Syntax: YAML-Blockskalare, Shell-Heredocs und
   SQL-Zeichenketten können eine Zeile vortäuschen, die als Kommentar gelesen
   wird — im Bestand ist kein solcher Fall bekannt; neue Fälle sind Befund.
3. **Zählregel, keine Bewertung.** Ein Kommentar mit zwei Ankern kann eine legitime
   Kopplung oder Abgrenzung sein; das Werkzeug meldet ihn trotzdem, die Regel
   verlangt dort die Stelle (Datei, Funktion) statt einer Kennungsreihe
   ([`AGENTS.md`](../../AGENTS.md) §3.7).
4. **Grün ist kein Beleg.** Ein Kommentar mit einer Kennung, die eine Spec-Aussage
   in eigenen Worten wiedergibt, ist kein Kandidat.
5. **Diff-Modus sieht Zeilen, nicht Absätze.** Wer eine Zeile eines Blocks ändert,
   den Block also berührt, bringt ihn in den Lauf; wer nur einen Nachbarblock
   ändert, nicht.

## Test

`make test-kommentar-kennungen` (auch unter `make test`) fährt den Tabellentest
`tools/harness/kommentar-kennungen/main_test.go`: die Zählregel (eine Kennung ·
zwei verschiedene · dieselbe zweimal · alle vier Arten · Kompaktform mit
Schrägstrich und mit Auslassungszeichen · Nummer falscher Breite · „ff.“ mit und
ohne Kennung davor), die Blockgrenzen (Leerzeile · Endkommentar hinter Code ·
Zeichenketten-Literal · Direktive · Blockkommentar · „ff.“ über einen
Zeilenumbruch), den Diff-Modus (überlappende und nicht überlappende Zeile ·
Löschung · andere Datei · leerer Diff), die Modi `-count`/`-tests`, die
ausgenommenen Wurzeln (`gen`, `sdks`, `.harness`, `.git`), das Präfix der
Zieldatei-Zeile und die Exit-Codes. Die Nicht-Go-Formen tragen eigene Fälle:
Form-Auswahl je Name/Endung, Blockgrenze der Zeilenformen (Leerzeile ·
Grenz-Marker · Einrückung · nachgestellter Kommentar nicht gelesen · SQL
`--` · Makefile-Zielzeile mit `##`), Lauf über den Baum mit allen Formen
(`-tests only` liest sie nicht, `.txt` bleibt ungelesen) und Diff-Modus für
eine Shell-Datei.

Der Aufrufer `tools/harness/kommentar-kennungen.sh` trägt seinen eigenen
Tabellentest, `tools/harness/run-kommentar-kennungen-tests.sh` (bash, ein
Wegwerf-Repo im Temp-Verzeichnis, `make test-kommentar-kennungen` fährt ihn nach
dem Go-Test): Argumentzahl, fehlende `TOOLCHAIN_IMAGE`, unbekannte Diff-Basis ·
Exit-Weitergabe (0, 1, 2) · jedes Wort von `PATHS`/`DIFF` ein eigenes Argument
(keine Marker-Datei eines eingeschleusten Kommandos, kein Glob) · der Diff-Strom
unter fremder Git-Konfiguration byte-gleich zu dem ohne Eingriff · ein
fehlschlagendes `git diff` · keine Temp-Datei danach — mit einem Stub für
`docker` — und vier Läufe mit echtem Docker (Kandidat im Diff-Modus mit und ohne
`diff.mnemonicPrefix`, ungültiger `TESTS`-Wert, fehlender Pfad).

## Belege der Definition

- **Zählregel gegen den Bestand** (gemessen, Stand `0d333120`): 30 Kandidaten
  (jede zwanzigste Zeile der Ausgabe ab der zehnten, `awk 'NR%20==10'`) — 25
  Verstöße (Reihung gleichrangiger Anker, Kette, Kompaktform, „ff.“), 5 Grenzfälle
  (zwei Anker tragen je eine eigene Aussage der Stelle). Eine zweite Stichprobe
  derselben Auswahlregel im Review (acht Zeilen, Stand `71cf9537`): vier
  Verstöße, vier Grenzfälle — `internal/domain/errors/errors.go` (Zeilen 35–37),
  `internal/application/port/inbound/consumer.go` (30–32),
  `internal/adapters/driven/postgresstorage/queries/queries.go` (241–244),
  `internal/application/port/inbound/verwaltung.go` (83–89).
- **Strenger als der Baseline-Wortlaut.** Die Baseline verlangt „ein auflösbares
  Feld“; die Zählregel („höchstens eine Kennung“) folgt der Zielform von
  [`AGENTS.md`](../../AGENTS.md) §3.7 und meldet auch einen Block mit zwei Ankern,
  die je eine eigene Aussage tragen (Grenze 3). Das ist die Zielform des Plans, kein
  Zufall der Implementierung.

## Fassung im Gate-Index

Ausführliche Fassung der Index-Zeile aus [`harness/README.md` §Sensors](../README.md#sensors-feedback-gates); die Zeile dort trägt einen Satz und verlinkt hierher.

### `make kommentar-kennungen`

listet die Kommentarblöcke der Go-Dateien und der Nicht-Go-Zeilenkommentar-Formen (`.sh`, `.mk`, `.yml`, `.yaml`, `.sql`, Makefile, Dockerfile), die ihre Herkunft nicht als ein auflösbares Feld tragen (Kandidat: mindestens zwei verschiedene Kennungen `ADR-`/`LH-FA-`/`LH-QA-`/`SPEC-`/`ARC-` oder „ff.“ hinter einer Kennung, eine Kompaktform zählt je Nummer): `make kommentar-kennungen [PATHS=<Pfade>] [COUNT=1] [TESTS=exclude|only] [DIFF=<Basis>]` — `COUNT=1` druckt nur die Zahl, `DIFF=<Basis>` meldet nur Blöcke, die eine seit `<Basis>` hinzugefügte Zeile überlappen. Docker-only (Go-Programm `tools/harness/kommentar-kennungen/`, gepinntes Toolchain-Image, `--network none`; Host-Werkzeuge `bash`, `git`, `docker`). Prüft die **Form, nicht die Wahrheit**: eine Spec-Wiedergabe in eigenen Worten hinter einer Kennung erkennt es nicht, ein Lauf ohne Kandidat sagt nicht „die Kommentare sind konform“. Keine Ausnahmeliste; Exit 1 bei mindestens einem Kandidaten (über `make` als Exit 2). Aufrufer: Schritt 20 des Implementer-Ablaufs, Reviewer (Probe)

### `make test-kommentar-kennungen`

Tabellentests des Programms `tools/harness/kommentar-kennungen/` (Zählregel, Blockgrenzen, Diff-Modus samt Präfix der Zieldatei-Zeile, Modi, ausgenommene Wurzeln, Exit-Codes; der Go-Test läuft auch unter `make test`) und seines Aufrufers `tools/harness/kommentar-kennungen.sh` (`tools/harness/run-kommentar-kennungen-tests.sh`: Eingabefehler, Exit-Weitergabe, Argument-Zerlegung, gepinnte Form des Diff-Stroms unter fremder Git-Konfiguration, Temp-Datei; Stub-`docker` plus vier Läufe mit echtem Docker); Docker-only, netzlos
