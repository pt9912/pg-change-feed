# Verifikationsbericht: slice-sdk-sse-client-schema-table-filter — 2026-10-02

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den DoD-Vertrag
des Slice-Plans
([`slice-sdk-sse-client-schema-table-filter.md`](../plan/planning/in-progress/slice-sdk-sse-client-schema-table-filter.md)
§2) und die §6-Risiko-Ausgänge, in frischem Kontext. Geprüft werden die
**Belege**, nicht die Behauptung ([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B).
Nicht der Diff als solcher (Reviewer, abgeschlossen mit
[`review-slice-sdk-sse-client-schema-table-filter.md`](review-slice-sdk-sse-client-schema-table-filter.md))
und nicht der reale Bedarf (Validator).

**Gegenstand:** `git diff e97fd62d HEAD` — sieben Commits: `ed5db509` (SDK-Packages),
`d9990c1f` (Beispiele), `07b10726` (Doku, Plan, Suchlauf), `0adfc502` (Review),
`81ad7643` (Fixrunde: Flag-Verdrahtung der Beispiele), `410e22ac` (Fixrunde:
Eingabetabelle), `9b1023c7` (Fixrunde: Plantext, Risiko). 29 Dateien, `718 insertions`,
`81 deletions` (`git diff --stat`, gelesen).

**Bezug:** [`LH-FA-SST-008`](../../spec/lastenheft.md),
[`LH-FA-SST-009`](../../spec/lastenheft.md),
[`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) Teilfrage 4,
[`ADR-0134`](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md),
[`ADR-0090`](../plan/adr/0090-beispiel-clients-volle-matrix.md).

**Frischer Kontext:** Plan (vollständig), Review-Report (vollständig), die
Produktions-Diffs der drei Package-Clients und der drei Beispiele, die Testdiffs
beider Fixrunden, `internal/adapters/driving/http/sse.go`
(`parseStreamChangesFilter`) gelesen; Mutationen und Läufe unten selbst gefahren,
keine Behauptung des Implementers oder Reviewers übernommen.

---

## 1. Eigene Sensor-Belege (ungepiped, [`AGENTS.md`](../../AGENTS.md) §3.9)

| Sensor | Ausgang | gedruckte Zeile aus meinem Lauf |
|---|---|---|
| `make gates` | **Exit 0** (Log in eine Datei, `ec=$?` danach gelesen) | `coverage-gate: OK — Coverage 82.00% erfüllt Schwelle 80%`; `d-check: 1544 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID`; `generated-sync: OK — das committete Erzeugnis ist byte-gleich …`; `sdk-public-doc-check: keine interne Kennung unter sdks`; `gesamt: 0 Befund(e)` (a-check); `baseline-verify: v6.13.0 OK — 54 Dateien` |
| `make test` | Exit 0 | alle Pakete `ok`, kein `FAIL` |
| `make docs-check` | Exit 0 | `d-check: 1544 Datei(en) geprüft, 0 Befund(e)` |
| `make sdk-public-doc-check` | Exit 0 | `keine interne Kennung unter sdks` |
| `make fmt-check` | Exit 0 | `323 Go-Dateien geprüft, alle formatiert` |
| `make suchlauf-nachmessen PLAN=…` | Exit 0 | `suchlauf-nachmessen: 10 Zeilen stimmen` |
| `make kommentar-kennungen DIFF=e97fd62d` | Exit 0 | keine Kandidaten |
| `make doc-commits RANGE=e97fd62d..HEAD` | Exit 0 | `--enable commits`: `0 Befund(e)` |
| `make doc-immutable RANGE=e97fd62d..HEAD` | Exit 0 | `--enable vcs`: `0 Befund(e)` |
| `make sdk-pack-python` / `-csharp` / `-kotlin`, `make examples-csharp` / `-kotlin` | je Exit 0 | **aber alle Bau-Schritte `CACHED`, keine Testzeile gedruckt** — der Cache-Lauf belegt keine Ausführung (siehe §3) |

Der Cache-Befund ist die Bestätigung der Plan-Warnung: Auf einem Arbeitsbaum, dessen
Eingaben schon einmal gebaut wurden, ist der grüne Exit der Pack- und Examples-Ziele
**kein** Testbeleg. Die Testzahlen holte ich mit `docker build --no-cache --target build`
je Wurzel nach (§3).

## 2. Mutationen (selbst gefahren)

Methode: `git archive HEAD` in einen Scratchpad-Ordner, Mutation als
`sed … > Kopie`, Rücknahme per `cp`; kein `sed -i`; Testlauf im Toolchain-Container
(Go) bzw. `docker build --target build` mit `BUILDKIT_PROGRESS=plain`. Die Bau-Läufe der
Mutationen setzten keine Tags (keine Mutations-Images); `git status --short` im
Echtrepo war nach den Mutationen und nach allen Läufen **leer**.

| # | Instanz | Mutation (Kopie) | gesehene Farbe |
|---|---|---|---|
| G1 | Go-Beispiel `examples/sse-client` | `fs.StringVar(&cfg.table, "table", …)` gelöscht | **rot** (`FAIL … examples/sse-client`) |
| G2 | Go-Beispiel | `config.streamURL`: `c.table` → `""` | **rot** |
| G3 | Go-Beispiel | `stream.go`: `"schema": schema` aus der Map entfernt | **rot** (`TestStreamURLCarriesSchemaAndTable`, Fall „schema eu") |
| C1 | C#-Beispiel | `StreamUrl(Config)`: `cfg.Table` → `""` | **rot**: `Failed: 2, Passed: 23, Total: 25` (`CliTests.FlagsWireIntoTheStreamUrl`) |
| C2 | C#-Beispiel | `case "--table"` aus `Cli.Parse` entfernt | **rot**: `Failed: 3, Passed: 22, Total: 25` |
| K1 | Kotlin-Beispiel | `streamUrl(Config)`: `cfg.table` → `""` | **rot**: `13 tests completed, 1 failed` (`CliTest.flagsWireIntoTheStreamUrl`) |
| K2 | Kotlin-Beispiel | `"--table" -> …` aus `Cli.parse` entfernt | **rot**: `13 tests completed, 2 failed` |
| P-C# | C#-Package | `schema=` im Query-Aufbau zu `xschema=` | **rot**: `Failed: 9, Passed: 150, Total: 159` |
| P-Py | Python-Package | `("schema", schema)` aus dem Tupel entfernt | **rot**: `3 failed, 147 passed` |
| P-Kt | Kotlin-Package | `"schema" to schema` entfernt | **rot**: `116 tests completed, 3 failed` |

Damit sind M5/M6 des Reviews (Flag gelöscht, `cfg.table` → `""`) im Go-Beispiel
geschlossen (G1, G2), und die hergeleitete Lücke bei C# und Kotlin ist **gemessen**
geschlossen (C1, C2, K1, K2).

**Rest-Lücke, hergeleitet, nicht gemutet:** Der Aufruf an der Einstiegsstelle selbst —
Go `main()` ruft `cfg.streamURL()`, C# `Program.cs:45` `SseStream.StreamUrl(cfg)`, Kotlin
`Main.kt:47` `SseStream.streamUrl(cfg)` — ist von keinem Test gebunden; ein Ersetzen
dieses einen Aufrufs durch die Form ohne `schema`/`table` ließe die Tests grün. Das ist
eine Zeile je Beispiel, die Test-Grenze liegt jetzt an der `main`-Funktion statt an der
Parser/URL-Naht (Einstufung LOW, kein Blocker; das vorhandene `-target` trug sie schon
vor dem Slice).

## 3. DoD-Zeilen einzeln

| # | DoD-Zeile | Beleg und Befund | Verdikt |
|---|---|---|---|
| 1 | (A) Packages: `schema`/`table` optional, nur gesetzt auf dem Draht; Tests für gesetzt, `null`, `""`, Konjunktion mit `target`; dieselbe Eingabetabelle `null`/`""`/`eu`/`a&b=c`/`a+b`/`100%`/`ü`/`a b`; Belegbefehl: Bau der Test-Stufe ohne Cache oder rote Mutation | Quelltext gelesen: C# `Uri.EscapeDataString`, Kotlin `percentEncode`, Python `httpx params`; alle drei bauen die Reihenfolge `schema`, `table`, `target`, `null` fehlt, `""` erscheint. **Testzahlen ohne Cache, selbst gefahren:** Python `150 passed`; C# `Passed: 159, Failed: 0, Total: 159`; Kotlin `Task :test`, `BUILD SUCCESSFUL` (kein Test-Zähler in Gradles Normalausgabe, die rote Variante druckt `116 tests completed`). Rote Mutationen P-C#, P-Py, P-Kt. Eingabetabelle vollständig in allen drei Sprachen (Zeilen `a+b`, `100%`, `ü`, `a b` je Parameter `schema`, `table`, `target`) | **erfüllt** |
| 2 | (B) Beispiele Go/C#/Kotlin: Flags; Tests Flag → Anfrage; Belegbefehle `make test`, `make examples-csharp`, `make examples-kotlin` | `make test` Exit 0. Ohne Cache: C# SSE `Passed: 25`, Kotlin `:sse-client:test` ausgeführt, `BUILD SUCCESSFUL`. Verdrahtungs-Tests (`TestFlagsWireIntoTheStreamURL`, `FlagsWireIntoTheStreamUrl`, `flagsWireIntoTheStreamUrl`) fahren echte Flag-Argumente bis zur URL; Mutationen G1, G2, C1, C2, K1, K2 rot. Rest-Lücke: Einstiegsaufruf (siehe §2) | **erfüllt** (mit LOW-Rest) |
| 3 | (C) `make sdk-public-doc-check` Exit 0; READMEs, `examples/README.md`, Handbuch (Version, Historie) nennen die Parameter; „`target` ist der einzige Filter" ersetzt; Zählwort erst nach (A) | Exit 0 (eigener Lauf, §1). Handbuch `Version: 1.87`, Historienzeile 1.87 im selben Diff, `Stand: 2026-10-02`. Eigener Suchlauf über den ganzen Baum (Muster `target only`, `only filter`, `einzige…Filter`, `filtered by (delivery )?target`, `does not set`, `not by table`, `setzen sie am SSE`; Ausnahmen `docs/reviews`, `done/`, `.harness/baseline`): die verbleibenden Treffer sind die Historienzeilen 1.85/1.86 (bewusst), die Plan-Prosa, eine Zeile in [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md) (anderer Gegenstand), der NATS-Satz der Python-README („NATS stream … by delivery target only", wahr: das NATS-Subjekt filtert nicht über Query-Parameter) und der Handbuch-Satz zu `GET /changes` („wählt nur den Filter `target`", anderer Weg). Keine stehengebliebene Aussage zum SSE-Client | **erfüllt** |
| 4 | `make gates` grün, Exit ungefiltert | Exit 0, §1 | **erfüllt** |
| 5 | Review durchgeführt, kein offenes HIGH/MEDIUM | Report liegt vor; **F-1 (MEDIUM) ist im Report noch als „merge-blockierend: ja" geführt**, die Fixrunde `81ad7643` behebt ihn nach meinen Mutationen (G1, G2, C1, C2, K1, K2 rot). Ein Re-Review fand nicht statt; der Report selbst trägt keinen Schließungsvermerk (Bedingung B1) | **bedingt erfüllt** |
| 6 | §3.13-Suchlauf: Feld trägt Gefundenes und Nichtgefundenes, beide Stände; `make suchlauf-nachmessen` Exit 0 | Exit 0, `10 Zeilen stimmen` (u. a. `diff 88` für die Beispiel-Flag-Zeile, `diff 0` für interne Kennungen unter `sdks/`); Feld in §3 des Plans trägt je Träger Befund und „Nicht gefunden". Eigener Suchlauf (Punkt 3) fand nichts darüber hinaus | **erfüllt** |
| 7 | Doku-Update; kein SPEC-/ARC-Eintrag; Package-Versionen bleiben | `git diff e97fd62d HEAD --stat` auf `docs/user/version.md`, die `.csproj`-, `pyproject.toml`-, `build.gradle.kts`-Version, `internal`, `cmd`, `spec`: **leer**. Kein Release, kein Tag | **erfüllt** |
| 8 | Closure-Notiz mit Steering-Loop-Eintrag | Plan §7 trägt „Wird bei der Closure gefüllt" — **offen** (Planner) | offen, erwartet |
| 9 | Reconciliation-Register | entfällt (Greenfield) | erfüllt |
| 10 | Beobachtungs-Register fortgeschrieben | im Diff **kein** Eintrag unter `observations/` — offen (Planner, §6) | offen, erwartet |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | alle vier Risiken tragen „bei der Closure einzutragen" — offen (Planner) | offen, erwartet |
| 12 | Drei Paarungen (Anker · Folge-Slice · Register) | Slice ohne Welle; Folge-Slice-Bedarf siehe §6 | offen, erwartet |

Keine DoD-Häkchen von mir gesetzt.

## 4. Entscheidungs-Konformität

- **[`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) Teilfrage 4** (SSE: dieselben Query-Parameter `schema`/`table`): die Clients senden nur die drei dem Server bekannten Namen (`streamChangesParams`), ein anderer Name endet am Server mit `400`. Gelesen: `parseStreamChangesFilter(r.URL.Query())` (`sse.go`).
- **Draht-Bytes und Server-Dekodierung (selbst gemessen):** `url.ParseQuery` liefert `schema=a+b` → `"a b"`, `schema=a%20b` → `"a b"`, `schema=a%2Bb` → `"a+b"`, `schema=%C3%BC` → `"ü"`, `schema=100%25` → `"100%"`, `schema=a%26b%3Dc` → `"a&b=c"`. Die erwarteten Draht-Bytes sind in den drei Package-Sprachen für `a+b` (`%2B`), `100%` (`%25`), `ü` (`%C3%BC`) und `a&b=c` (`%26`/`%3D`) **gleich**; ein Leerzeichen sendet C#/Kotlin als `%20`, Python als `+` — als einzige Abweichung in jedem Test benannt (Kommentar „the server decodes both alike") und im Python-Docstring dokumentiert; der Server dekodiert beide Formen zum selben Wert (gemessen).
- **[`ADR-0134`](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md):** `make sdk-public-doc-check` Exit 0; Suchlauf-Zeile `diff 0` für `ADR-|LH-FA|LH-QA|SPEC-|ARC-` unter `sdks/` stimmt.
- **[`ADR-0090`](../plan/adr/0090-beispiel-clients-volle-matrix.md):** Beispiele bleiben Doku mit Bau-Bindung; die drei Beispiel-Wurzeln bauen und testen (§3 Zeile 2).
- **Hard Rules:** §3.1 kein Host-Werkzeug im Diff; §3.2 kein `//nolint`/`# noqa`; §3.7 `make kommentar-kennungen DIFF=e97fd62d` Exit 0; §3.9 Exit-Codes aus Dateien/ungepiped gelesen; §3.11 dieser Report nennt keinen Host-Pfad.
- **Commit-Betreffs:** `make doc-commits` und der Standing-Gate-Anteil in `make gates` Exit 0; kein `SPEC-*`/`ARC-*` im Betreff.

## 5. Plan-vs-Code-Diff

- **Package-Clients (C#, Kotlin, Python):** wie §3 Punkt „Ansatz" geplant — neue Parameter hinter `target`, kein Verschieben bestehender Parameter, Wire-Reihenfolge `schema`, `table`, `target`. Die Parameter-Reihenfolge des SSE-Clients (`ct, target, schema, table`) weicht von der des gRPC-Clients (`schema, table, ct, target`) ab; der Plan legt sie so fest („nach `target` angefügt"), die READMEs sagen „by name" (Review F-4, INFO).
- **Kein Server-/Produktivcode im Diff:** `git diff e97fd62d HEAD --stat` auf `internal`, `cmd`, `spec` leer. Alle Änderungen liegen in `sdks/`, `examples/`, `docs/`.
- **Refactor der Beispiele (Umfangsfrage des Auftrags):** Go: `parseFlags()` → `parseConfig(args, getenv)` mit eigenem `flag.NewFlagSet(…, flag.ContinueOnError)` sowie `config.streamURL()`; C# und Kotlin: eine Überladung `StreamUrl(Config)`/`streamUrl(Config)`. Der Plan führt diese Zeile in §3 ausdrücklich („Fixrunde nach dem Review … `parseConfig`, `config.streamURL`, Überladung"), und sie ist die kleinste Form, die die vom Review geforderte Bindung ermöglicht — der Umfang wird **nicht gesprengt**. Eine Verhaltensänderung im Go-Beispiel ist zu benennen: `flag.CommandLine` mit `ExitOnError` ist durch ein `FlagSet` mit `ContinueOnError` ersetzt, `main` meldet den Fehler selbst (`sse-client: …`) und beendet mit Exit 2; ein `-h` endete zuvor mit Exit 0 und Usage, jetzt mit der Fehlerzeile und Exit 2. Das betrifft ein Beispiel-Programm, keine Zusage des Plans; Einstufung INFO.

## 6. Register und Nachzug für die Closure (Planner)

Im Diff steht **kein** Register-Eintrag. Vorschläge, kein Beschluss meinerseits:

- [`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`](../plan/planning/observations/BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall/state.md) steht auf `gestrichen` (3×). Dieser Slice ist **kein Divergenz-Fund am Randwert**: der Plan nannte die Eingabetabelle als Pflicht je Sprache, die Fixrunde erweiterte sie auf `+`, `%`, `ü`, Leerzeichen, und die Draht-Bytes sind in drei Sprachen gleich (außer dem benannten Leerzeichen). Die Gegenmaßnahme „derselbe Eingabesatz je Sprache" hat **erneut gewirkt** — ein Vermerk „Gegenmaßnahme beobachtet, kein neuer Anfall" wäre belegt; der Reviewer-Fund F-2 war eine Lücke in der Tabelle, nicht in einer Sprache.
- [`BEO-PGC/docker-cache-ueberspringt-tests-still`](../plan/planning/observations/BEO-PGC/docker-cache-ueberspringt-tests-still/state.md) steht auf `offen` (1×, Schwelle 3×). In diesem Slice trat die Falle **real** auf: alle fünf Ziele (`sdk-pack-python`, `-csharp`, `-kotlin`, `examples-csharp`, `examples-kotlin`) liefen Exit 0 aus dem Cache **ohne eine Testzeile**; die Testzahl entstand erst im `--no-cache`-Bau. Das ist ein zweites Auftreten der Klasse (vom Plan vorab benannt und von Reviewer und Verifier durch Mutation bzw. `--no-cache` gelöst). Ob es als `evidence/`-Datei (2×) oder als Vermerk steht, entscheidet der Planner.
- Kandidat für einen neuen Eintrag oder einen Vermerk bei [`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`](../plan/planning/observations/BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe/state.md) (Deckel 14×): Review F-1 „Verdrahtung Flag → Anfrage ungebunden", gefunden vom Reviewer vor dem Merge, mit Mutation belegt, in der Fixrunde behoben.
- **Risiken aus §6 des Plans** (Ausgänge einzutragen): Randfall-Kopie → entfallen (gemessen: Eingabetabelle in drei Sprachen, Mutationen rot); Cache → eingetreten, gelöst mit `--no-cache`; Binärkompatibilität → weiter offen bis zum Release (quellkompatibel, kein Release in diesem Slice); interne Kennung → entfallen (`make sdk-public-doc-check` Exit 0).
- **Folge-Slices:** keiner nötig für den Gegenstand. Realserver-Belege der SDK-Tiers für `schema`/`table` am SSE-Client sind im Plan ausdrücklich ausgeschlossen (eigener Zug im Muster `make test-sdk-*-integration`); der Planner kann sie als Folge-Slice benennen. Die Dopplungs-Aussage des Handbuchs („Beispiel-Clients und Packages nehmen alle drei Parameter entgegen") ist auf Unit-Ebene belegt, nicht am Realserver.

## 7. Verdikt

**Bestanden, unter einer Bedingung.**

- **B1 (formal, Reviewer/Planner):** Der Review-Report führt F-1 weiter als „Merge-blockierend: ja"; die DoD-Zeile „Review durchgeführt, kein offenes HIGH/MEDIUM" ist damit auf dem Papier offen. **Ob ein weiterer Reviewer-Durchgang nötig ist:** für den Code **nein** — die Fixrunde änderte nur Beispiel-Verdrahtung (testbare Funktion), Tests und Plantext; Produktionscode der Packages blieb unverändert, und ich habe die sechs Verdrahtungs-Mutationen und die drei Package-Mutationen selbst rot gesehen. Es genügt ein Nachtrag im Review-Report (Schließung F-1 bis F-3 mit Verweis auf die Fixrunde und diesen Bericht) oder eine Planner-Notiz in der Closure; ein vollständiger zweiter Review ist nicht erforderlich. Wird die Einstufung trotzdem gewünscht (HIGH-Liste führt die Klasse „Zusage ohne Bindung"), ist ein kurzer Re-Review der beiden Fixrunden-Commits `81ad7643` und `410e22ac` ausreichend.

**Abweichungen (alle nicht blockierend):**

1. LOW — Einstiegsaufruf (`main`/`Program.cs`/`Main.kt`) an die URL-Funktion ungebunden (§2, hergeleitet).
2. INFO — Go-Beispiel: `-h` endet nun mit Exit 2 statt 0 (§5).
3. INFO — Parameter-Reihenfolge SSE-Client gegen gRPC-Client verschieden (Review F-4).
4. INFO — `make sdk-pack-*`/`make examples-*` liefern bei unverändertem Eingang keinen Testbeleg (Cache); Belegbefehl bleibt `--no-cache` auf die Test-Stufe oder eine rote Mutation.

**Offen für den Planner:** Closure-Notiz (§7 des Plans), Risiko-Ausgänge (§6), Register-Fortschreibung (§6 dieses Berichts), DoD-Häkchen, `git mv` nach `done/`, Zählwort-Fortschreibung. Kein Push, kein Release, keine Versionsänderung.
