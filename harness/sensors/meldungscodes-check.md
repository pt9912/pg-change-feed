# `make meldungscodes-check` — Meldungscodes in Quelltext, Code-Tabelle und Handbuch-Katalog gleich

## Vertrag

Wird dieses Target rot, weichen die Meldungscodes (`PCF-<S><NNNN>`) von Quelltext,
Code-Tabelle und Handbuch-Katalog voneinander ab
(`tools/harness/meldungscodes-check.sh`). Die Tabelle
`internal/domain/messagecode/codes.go` ist die Quelle der Wahrheit; der Katalog im
Benutzerhandbuch führt dieselbe Menge; der Quelltext verwendet nur Codes der
Tabelle ([`ADR-0144`](../../docs/plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md)
Festlegung 5, [`LH-QA-OPS-001`](../../spec/lastenheft.md)).

**Geprüft** (ERE, kein `grep -P`; die Token-Suche liest `PCF-[A-Za-z0-9]*`, die Form
prüft `^PCF-[EWI][0-9]{4}$`):

| Nr. | Prüfung | Befund |
|---|---|---|
| 1 | Jeder Token `PCF-…` in Quelltext, Tabelle und Handbuch hat die Form `PCF-[EWI][0-9]{4}` | `Code in falscher Form: <datei>:<zeile>:<token> (erwartet PCF-[EWI][0-9]{4})` |
| 2 | Jeder Code im Quelltext und im Handbuch steht in der Tabelle | `Code ohne Eintrag in der Tabelle <tabelle>: <datei>:<zeile>:<code>` |
| 3a | Jeder Tabellen-Code hat eine Katalog-Zeile | `Tabellen-Code ohne Katalog-Zeile in <handbuch>: <code>` |
| 3b | Jede Katalog-Zeile nennt einen Tabellen-Code | `Katalog-Zeile ohne Eintrag in der Tabelle <tabelle>: <code>` |

Eine **Katalog-Zeile** ist eine Markdown-Tabellenzeile des Handbuchs, deren erste
Zelle (optional in Backticks) mit einem Code beginnt; je Zeile gilt der erste Code.
Ein Code in Prosa oder in einer späteren Zelle ist keine Katalog-Zeile (Prüfung 2
gilt für ihn trotzdem). Ein **zurückgezogener** Code bleibt in Tabelle und Katalog und
wird wie jeder andere Code verglichen; den Status liest das Gate nicht.

**Gegenstand:**

| Menge | Quelle |
|---|---|
| Tabelle | `internal/domain/messagecode/codes.go` |
| Katalog | `docs/user/benutzerhandbuch.md` |
| Quelltext | Produktions-Go (`*.go` ohne `*_test.go` und `*.pb.go`) unter `internal/`, `cmd/` und `tools/schema/`, ohne das Paket `internal/domain/messagecode/`; Skripte `*.sh` unter `tools/schema/` und `examples/` |

Nicht Gegenstand: Tests, die Läufer unter `tools/harness/`, erzeugter Code, das
Paket der Tabelle selbst (es trägt die Form und die Konstanten; sein Test
`internal/domain/messagecode/messagecode_test.go` läuft in `make test`).

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | die Mengen sind gleich; der Erfolgstext nennt die Zahl der Codes und der gelesenen Dateien (`meldungscodes-check: <n> Codes in Tabelle und Katalog gleich, Quelltext (<g> Go-Dateien, <s> Skripte) nur mit Codes der Tabelle`) |
| 1 | mindestens ein Befund; je Befund eine Zeile auf stderr, darunter die Sammelzeile `meldungscodes-check: Quelltext, Tabelle und Katalog sind nicht gleich (siehe oben)` |
| 2 | **Lesefehler** (fail-closed, nie „gleich"): fehlende Wurzel, fehlende Tabelle oder fehlender Katalog (`Lesefehler: Code-Tabelle fehlt: …` bzw. `Katalog fehlt: …`), eine nicht lesbare Datei (`Lesefehler: <datei> (<Meldung von grep>)`), ein nicht lesbares oder fehlendes Verzeichnis, eine Tabelle ohne Code und ein leerer Gegenstand (`Lesefehler: kein Gegenstand (<g> Go-Dateien, <s> Skripte); leer ist nicht bestanden`) |

Über `make` kommt jeder Exit ≠ 0 des Skripts als Make-eigener Exit `2` an.

## Overrides

| Variable | Wirkung |
|---|---|
| erstes Argument | abweichende Wurzel (Default: die Repo-Wurzel) — Aufrufform für den Tabellentest `tools/harness/run-meldungscodes-check-tests.sh` und für Mutationsproben an einer Kopie; das Gate ruft das Skript ohne Argument |
| `CHECK` (Umgebung des Tabellentests) | Pfad des Prüflings für `run-meldungscodes-check-tests.sh` (Default `tools/harness/meldungscodes-check.sh`) — Mutationsläufe des Wächters an einer Kopie |

## Grenze — was das Grün nicht abdeckt

1. **Mengen, nicht Sinn.** Ein falsch beschriebener Katalog-Eintrag (Bedeutung oder
   Maßnahme), ein Code mit falscher Klasse im Katalog und ein Code, den kein Fehlerwert
   verwendet, bleiben grün. Die Klasse prüft der Go-Test des Pakets (erste Ziffer gegen
   Klasse, Dopplung, Rückfall je Klasse); die Bedeutung liest der Reviewer.
2. **Tabelle gleich Menge der Token in `codes.go`.** Ein Code in einem Kommentar der
   Datei zählt als Tabellen-Code und braucht dann eine Katalog-Zeile; der Go-Test
   liest die Tabelle als Wert und fängt einen Kommentar-Code nicht.
3. **Paket der Tabelle ausgenommen.** Ein Code-Literal in `messagecode.go` (außerhalb
   von `codes.go`) liest das Gate nicht.
4. **Andere Ausgabewege.** SQL-Funktionen, SDK-Quellen unter `sdks/`, die Läufer unter
   `tools/harness/` und Beispiel-Programme in Go/C#/Kotlin unter `examples/` liest das
   Gate nicht; eine Katalog-Zeile außerhalb des Benutzerhandbuchs zählt nicht.
5. **Nicht-Token-Formen.** Ein Code, der zur Laufzeit aus Teilen entsteht
   (`"PCF-" + …`), ist ein Token `PCF-` und damit ein Befund der Form; eine
   Konstante ohne Literal liest das Gate nicht — die Konstante ist die Tabelle.
6. **Die Wächter-Logik ist im Gate-Lauf nicht gegen ihren Tabellentest gesichert.**
   `make test-meldungscodes-check` bleibt Werkzeug; der Ort, an dem eine Änderung am
   Wächter gegen diesen Vertrag gelesen wird, ist der Review-Diff. Die vier
   Lesefehler-Fälle mit `chmod 000` entfallen unter root; der Lauf druckt das dann.

## Sperren

- **Host-Werkzeuge:** `bash`, `git`, `grep`, `find`, `sed` (ohne `-i`), `sort`, `comm`,
  `head`, `wc`, `mktemp` (`AGENTS.md` §3.1, POSIX-/coreutils-Basis). Kein Docker, kein
  Netz; das Skript steigt über `git rev-parse --show-toplevel` auf die Wurzel und legt
  ein Temp-Verzeichnis an, das es beim Ende entfernt. Der Wächter liest den
  Arbeitsbaum, nicht den Index.

## Bindung

[`ADR-0144`](../../docs/plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md)
Festlegung 5 (Aufnahme in `GATE_CHECKS`, Prüfungen, Heimat `harness/mk/doc-gate.mk`) ·
[`LH-QA-OPS-001`](../../spec/lastenheft.md) (Betriebsfähigkeit) ·
`tools/harness/meldungscodes-check.sh` · `harness/mk/doc-gate.mk` · seit
slice-meldungscodes-registry-fehlerkopf.

## Fassung im Gate-Index

Ausführliche Fassung der Index-Zeile aus [`harness/README.md` §Sensors](../README.md#sensors-feedback-gates); die Zeile dort trägt einen Satz und verlinkt hierher.

### `make meldungscodes-check`

gleicht die Meldungscodes `PCF-<S><NNNN>` mit `grep` (ERE, kein `-P`) ab: jeder Code im Quelltext (Produktions-Go unter `internal/`, `cmd/` und `tools/schema/` ohne das Paket der Tabelle, Skripte unter `tools/schema/` und `examples/`) und im Handbuch steht in der Code-Tabelle `internal/domain/messagecode/codes.go`, und die Codes der Tabelle sind gleich den Codes der Katalog-Zeilen im Benutzerhandbuch (`tools/harness/meldungscodes-check.sh`); ein Token `PCF-…` in falscher Form ist ein Befund. Fail-closed: ein Lesefehler, eine fehlende Tabelle oder ein fehlender Katalog, eine Tabelle ohne Code und ein leerer Gegenstand enden mit Exit 2. Netzlos und schnell (reines `grep`, kein Docker). Grenzen im Vertrag: das Gate vergleicht Mengen, nicht Sinn — Klasse und Bedeutung eines Codes liest der Go-Test des Pakets bzw. der Reviewer

### `make test-meldungscodes-check`

fährt den Tabellentest gegen `tools/harness/meldungscodes-check.sh` (`tools/harness/run-meldungscodes-check-tests.sh`: je Zweig ein Fall mit Meldungstext — Code im Go-Quelltext, in `cmd/` und im Skript ohne Tabelle, falsche Form (Buchstabe, drei und fünf Ziffern, Präfix ohne Code), Tabellen-Code ohne Katalog-Zeile, Katalog-Zeile ohne Tabellen-Code, Codes nur in Prosa oder in einer späteren Zelle, Prosa-Code ohne Tabelle; Nicht-Befunde saubere Menge, zurückgezogener Code in beiden Mengen, Test-Datei, erzeugte Datei, Paket der Tabelle, Läufer-Verzeichnis; fehlende Wurzel, Tabelle, Katalog und Verzeichnis, Tabelle ohne Code, leerer Gegenstand und nicht lesbare Dateien mit Exit 2); der Tabellentest bleibt Werkzeug — der Gegenstand des Gates ist der Bestand der Codes, nicht die Wächter-Logik (Muster [`ADR-0134`](../../docs/plan/adr/0134-sdk-public-doc-check-gate-make-gates.md) Teilfrage 2); netzlos
