# `make ausgabe-kennungen-check` — keine interne Kennung in Ausgabe-Literalen des Servers und seiner Betriebs-Skripte

## Vertrag

Wird dieses Target rot, trägt ein Ausgabe-Literal des Servers oder eine
`echo`-/`printf`-Zeile eines Betriebs-Skripts eine **interne Kennung**
(`LH-…`, `ADR-…`, `SPEC-…`, `ARC-…`; `tools/harness/ausgabe-kennungen-check.sh`). Die Programm-Ausgabe ist die
Betreiber-Sicht ([`LH-QA-OPS-001`](../../spec/lastenheft.md)); Kennungen der
Anforderungs-, Spezifikations- und Entscheidungsdokumente gehören in Plan, Spec
und Commit ([`ADR-0144`](../../docs/plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md)
Festlegung 10).

**Muster** (ERE, zeilenweise `grep -anE`, kein `grep -P`), zweistufig:

```text
Go     = ^[^"`]*("[^"]*"[^"`]*)*["`][^"`]*(LH|ADR|SPEC|ARC)-[A-Z0-9]     und die Zeile beginnt nicht mit //
Shell  = (echo|printf)[[:space:]].*(LH|ADR|SPEC|ARC)-[A-Z0-9]            und die Zeile beginnt nicht mit #
```

Das Go-Muster liest „eine Kennung hinter einem **ungeschlossenen**
Anführungszeichen oder Backtick": vor dem öffnenden Zeichen stehen beliebig viele
geschlossene `"…"`-Literale. Eine Kennung im Bezeichner, in einem Kommentar ohne
Anführungszeichen davor und in einem Struct-Tag (`` `json:"x" ADR-…` ``) trifft
nicht.

**Gegenstand** ([`ADR-0144`](../../docs/plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md)
Festlegung 10):

| Wurzel | Dateien | Zeilen |
|---|---|---|
| `internal/`, `cmd/`, `tools/schema/` | `*.go` ohne `*_test.go` und `*.pb.go` | die Zeilen des Go-Musters |
| `tools/schema/`, `examples/` | `*.sh` | die `echo`-/`printf`-Zeilen des Shell-Musters |

Nicht Gegenstand: Tests, die Läufer unter `tools/harness/` (Entwickler-Ausgabe,
nicht Betreiber-Ausgabe), Quellcode-Kommentare
([`AGENTS.md`](../../AGENTS.md) §3.7 regelt sie, `make kommentar-kennungen` liest
sie), erzeugter Code unter `gen/`, Spec und ADR-Texte.

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | keine Treffer; der Erfolgstext nennt die Zahl der gelesenen Dateien (`ausgabe-kennungen-check: keine interne Kennung in Ausgabe-Literalen von <n> Go-Dateien und <m> Skripten`) |
| 1 | mindestens ein Treffer; je Treffer eine Zeile `<datei>:<zeile>:<text>` auf stderr, darunter die Sammelzeile `ausgabe-kennungen-check: interne Kennung in einem Ausgabe-Literal (siehe oben)` |
| 2 | **Lesefehler** (fail-closed, nie „sauber"): eine nicht lesbare Datei (`Lesefehler: <datei> (<Meldung von grep>)`, `grep`-Exit ≥ 2), ein nicht lesbares Verzeichnis (`Lesefehler beim Durchsuchen von <wurzel> (<Meldung von find>)`), eine fehlende Wurzel (`Lesefehler: Verzeichnis fehlt: <wurzel>`) und ein leerer Gegenstand (`Lesefehler: kein Gegenstand (<n> Go-Dateien, <m> Skripte); leer ist nicht bestanden`). `grep`-Exit 1 (kein Treffer) bleibt Erfolg; eine Datei mit NUL-Byte wird als Text gelesen (`grep -a`), ihre Treffer bleiben sichtbar |

Über `make` kommt jeder Exit ≠ 0 des Skripts als Make-eigener Exit `2` an.

## Overrides

| Variable | Wirkung |
|---|---|
| erstes Argument | abweichende Wurzel (Default: die Repo-Wurzel) — Aufrufform für den Tabellentest `tools/harness/run-ausgabe-kennungen-check-tests.sh` und für Mutationsproben an einer Kopie; das Gate ruft das Skript ohne Argument |

## Grenze — was das Grün nicht abdeckt

1. **Mehrzeiliges Raw-String-Literal.** Das Skript liest Zeilen; eine Kennung in
   der Mitte eines Backtick-Literals über mehrere Zeilen steht in einer Zeile ohne
   öffnendes Zeichen und trifft nicht. Akzeptiertes Negativ: alle 14 Treffer des
   Bestands waren einzeilig ([`ADR-0144`](../../docs/plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md)
   Festlegung 10); ein Fund eines mehrzeiligen Literals mit Kennung macht den
   Wächter zum AST-Programm (Muster `tools/harness/kommentar-kennungen/`).
2. **Heredoc und mehrzeiliges `echo`.** Eine `cat <<EOF`-Ausgabe und die
   Fortsetzungszeile eines `echo` mit `\` tragen kein `echo`/`printf` in der
   Zeile mit der Kennung und treffen nicht.
3. **Nachgestellter Kommentar mit zitiertem Literal** (`x := 1 // Fehler "ADR-…"`)
   trifft: ein Falsch-Positiv in der lauten Richtung, im Tabellentest als Fall
   benannt. Ebenso ein nachgestellter `# …`-Kommentar hinter einem `echo` mit
   Kennung.
4. **Block-Kommentar `/* … */`** wird nicht als Kommentar erkannt: eine Zeile darin
   mit einem Anführungszeichen vor der Kennung trifft.
5. **Einzeiliger Raw-String mit inneren Anführungszeichen und Rune-Literal `'"'`.**
   Gemessen mit dem Skript an einer Wegwerf-Wurzel (Tabellentest, Exit 0 = Lücke):
   `` var s = `{"a": "ADR-0043"}` `` und `x := '"' + " ADR-0043"` enden beide mit
   Exit 0 — die gerade Anzahl `"` vor der Kennung lässt das Muster das Backtick-
   bzw. das zweite Literal nicht als offen erkennen. Beide stehen als „benannte
   Grenze“-Fälle im Tabellentest; eine künftige Schließung färbt sie rot.
6. **Andere Ausgabewege.** Die Gegenstand-Liste ist abschließend: SQL-Funktionen
   (`RAISE`), SDK-Quellen unter `sdks/` (eigenes Gate `make sdk-public-doc-check`),
   Beispiel-Programme in Go/C#/Kotlin unter `examples/` und Texte, die zur
   Laufzeit aus Daten entstehen, liest das Gate nicht.
6. **Form, nicht Sinn.** Eine Ausgabe, die ihre Aussage ohne Kennung unvollständig
   lässt (Klammer-Rest, Satz ohne Subjekt), bleibt grün; diese Hälfte liest der
   Reviewer.
7. **Die Wächter-Logik ist im Gate-Lauf nicht gegen ihren Tabellentest
   gesichert.** `make test-ausgabe-kennungen-check` bleibt Werkzeug; der Ort, an
   dem eine Änderung am Wächter gegen diesen Vertrag gelesen wird, ist der
   Review-Diff. Die Lesefehler-Fälle des Tabellentests (`chmod 000`) entfallen
   unter root; der Lauf druckt das dann.

## Sperren

- **Host-Werkzeuge:** `bash`, `git`, `grep`, `find`, `sort`, `sed`, `head`,
  `mktemp` (`AGENTS.md` §3.1, POSIX-/coreutils-Basis). Netzlos und schnell
  (reines `grep`), kein Docker; das
  Skript steigt über `git rev-parse --show-toplevel` auf die Wurzel und legt ein
  Temp-Verzeichnis an, das es beim Ende entfernt. Der Wächter liest den
  Arbeitsbaum, nicht den Index.

## Tabellentest

`make test-ausgabe-kennungen-check` fährt den Tabellentest gegen
`tools/harness/ausgabe-kennungen-check.sh`
(`tools/harness/run-ausgabe-kennungen-check-tests.sh`, netzlos): je Kennungsart
ein Treffer im Go-Literal und im `echo`/`printf`, Literal mit Schrägstrich und
mit URL davor, Raw-String in Backticks; Nicht-Treffer Kommentarzeile,
nachgestellter Kommentar ohne Anführungszeichen, Struct-Tag, Test-Datei,
erzeugte Datei, Läufer-Verzeichnis, Dateien außerhalb der Wurzeln; der benannte
Falsch-Positiv (Grenze 3) und die benannten Grenzen (Grenze 1 und 5); fehlende
Wurzel, leerer Gegenstand und Lesefehler mit Exit 2, je mit Meldungstext. Der
Tabellentest bleibt Werkzeug — der Gegenstand des Gates ist der Bestand der
Ausgabe-Literale, nicht die Wächter-Logik (Muster
[`ADR-0134`](../../docs/plan/adr/0134-sdk-public-doc-check-gate-make-gates.md)
Teilfrage 2; Grenze 7).

## Bindung

[`ADR-0144`](../../docs/plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md)
Festlegung 10 (Aufnahme in `GATE_CHECKS`, Muster, Reichweite) ·
[`LH-QA-OPS-001`](../../spec/lastenheft.md) (Betriebsfähigkeit) ·
`tools/harness/ausgabe-kennungen-check.sh` · `harness/mk/doc-gate.mk` · seit
slice-meldungscodes-kennungsfreie-ausgaben.
