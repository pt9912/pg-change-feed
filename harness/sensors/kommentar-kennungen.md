# `make kommentar-kennungen` — listet Go-Kommentarblöcke, die ihre Herkunft nicht als ein Feld tragen

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
Kommentarblöcke der `.go`-Dateien und meldet jeden, der sie verletzt
(**Kandidat**). Die Herkunft der Regel und die Grenze „kein Sensor über Prosa“
führt [`ADR-0083`](../../docs/plan/adr/0083-herkunft-von-aussagen-in-traegern.md).

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

Ein **Kommentarblock** ist eine Kommentargruppe des Go-Parsers (`go/parser`,
`ParseComments`): eine Folge von Kommentaren ohne Leerzeile und ohne Code
dazwischen; ein Endkommentar hinter Code bildet einen eigenen Block. Direktiven
(`//go:…`) zählen nicht mit. Eine Kennung in einem Zeichenketten-Literal ist kein
Kommentar und wird nicht gelesen.

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
| `DIFF=<Basis>` | nur Blöcke, die eine seit `<Basis>` **hinzugefügte** Zeile überlappen (Eingabe: `git diff -U0 <Basis> -- '*.go'`); ein Block, den der Diff nur berührt, ohne eine Zeile zu ändern, und eine reine Löschung zählen nicht. Bestandskandidaten färben den Lauf eines Implementers nicht. Neue Dateien vorher `git add` (der Diff sieht nur getrackte Dateien) |

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

## Exit-Codes

| Exit | Bedeutung |
|---|---|
| 0 | kein Kandidat, oder `COUNT=1` (die Zahl steht auf stdout) |
| 1 | mindestens ein Kandidat |
| 2 | Eingabefehler: unbekannter Wert für `TESTS`, `DIFF` ist kein Commit, ein Pfad fehlt, eine `.go`-Datei ist nicht lesbar oder der Diff-Strom nicht parsbar |

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
2. **Nur Go-Kommentare.** Skripte, `Makefile`, `.sql`, `.yml` und Dockerfiles liest
   das Werkzeug nicht; `sdks/` deckt `make sdk-public-doc-check`.
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
ausgenommenen Wurzeln und die Exit-Codes. Nicht gebunden ist der Aufrufer
`tools/harness/kommentar-kennungen.sh` (Docker-Aufruf, Temp-Datei des Diffs): er
hat keinen Tabellentest.
