# Review-Report: slice-code-kommentare-kennungen — 2026-09-26

**Review-Art:** Code — der Diff liefert die geschärfte Regel „Herkunft im Go-Kommentar ist ein Feld“
(`AGENTS.md` §3.7, Schritt 20 des Implementer-Ablaufs, ein neuer MEDIUM-Unterpunkt im Reviewer-Skill),
das Werkzeug `make kommentar-kennungen` (Go-Programm, Aufrufer, Vertrag, zwei Make-Ziele, zwei
README-Zeilen) und den Erstbeleg (drei Kommentarblöcke in `changestream.go`); geprüft gegen Plan,
[ADR-0083](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md), `AGENTS.md` §3.7/§3.9/§3.12/§3.13,
Baseline-Regelwerk und das Architect-Verdikt `architect-verdict-slice-chronik-in-code-kommentar`
(Modul 10 §Drei Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Slice `slice-code-kommentare-kennungen` (ohne Welle, Harness-Querschnitt), Diff-Range
`1021f6fe..71cf9537` (7 Commits, 11 Dateien, +983/−36; Baum sauber, nicht gepusht).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ in der Form des Diff-Stands
`71cf9537` (mit dem neuen Unterpunkt „Herkunft als mehrere Felder, Kette, ‚ff.‘ oder
Spec-Wiederholung“, den dieser Lauf zugleich prüft).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-26.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-code-kommentare-kennungen` (§1 Ziel und Zielform, §2 Liefer-Punkte als Bezug, §3
  Plan, Nachzug-Tabelle, Basis-Messung, Stichprobe und Suchlauf-Feld, §4 Rückführungen, §6 Risiken);
  Architect-Verdikt `architect-verdict-slice-chronik-in-code-kommentar`
- [ADR-0083](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von Aussagen, Grenze
  „kein Sensor auf Prosa“), [ADR-0060](../plan/adr/0060-grpc-streaming-mechanismus.md) (Anker des
  Erstbelegs), [ADR-0066](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md) (Broadcaster,
  Gegenstand der Zusagen-Prüfung)
- `AGENTS.md` (Hard Rules §3.1, §3.2, §3.3, §3.6, §3.7, §3.9, §3.12, §3.13), `harness/conventions.md`
  (`MR-000`/`MR-001`)
- Baseline `v6.9.0` · `regelwerk/grundlagen-harness-dateien.md` §Was ein Kommentar trägt
  (Hard Rule „Herkunft als **ein** auflösbares Feld“), `regelwerk/grundlagen-traceability.md`
  §Herkunfts-Anker; Kennungs-Muster des Repos aus `.d-check.yml` (`ids.patterns`)
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-harness-suchlauf-nachmessen.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem Implementer-Bericht übernommen;
Exit-Codes ungepiped in Log-Dateien gesichert, gedruckte Zeilen zitiert; Scratchpad der Sitzung):

- **Läufe am Stand `71cf9537`** (einer zugleich; `free -m` vorher 13,5 GB verfügbar; dangling Volumes
  `docker volume ls -qf dangling=true | wc -l` 34 vor und 34 nach allen Läufen; kein `prune`;
  `git status --short` danach leer):
  `make test-kommentar-kennungen` Exit 0, gedruckt „ok … tools/harness/kommentar-kennungen“;
  `make test` (Race) Exit 0, 45 Pakete `ok`, darunter `tools/harness/kommentar-kennungen`;
  `make a-check` Exit 0, gedruckt „gesamt: 0 Befund(e)“;
  `make coverage-gate` Exit 0, gedruckt „coverage-gate: OK — Coverage 85.00% erfüllt Schwelle 80%“;
  `make gates` (einmal) Exit 0, gedruckt „d-check: 1255 Datei(en) geprüft, 0 Befund(e)“ und
  „gesamt: 0 Befund(e)“; `make commit-traceability RANGE=1021f6fe..71cf9537` Exit 0, „7 Commit(s)“;
  `gofmt -l` im gepinnten Toolchain-Image über `tools/harness/kommentar-kennungen` und
  `changestream.go`: keine Ausgabe, Exit 0.
- **Basis-Messung nachgemessen** (Werkzeug am Stand `71cf9537`): `make kommentar-kennungen COUNT=1`
  → **597**, `… TESTS=exclude` → **400**, `… TESTS=only` → **197**. Der Nenner des Plans: eine Kopie des
  Programms mit der Schwelle „mindestens eine Kennung“, gefahren gegen den Baum, meldet **1474** Blöcke
  (782 Nicht-Test, 692 Test); 597 von 1474 sind 40,5 % (abgeleitet), unter der Hälfte — die
  Rückführung (a) aus §4 tritt nicht ein. Die Plan-Zahlen 600/403/197 (Stand `0d333120`) passen zu
  597/400/197 nach dem Erstbeleg (drei Blöcke, alle in Nicht-Test-Dateien).
- **Suchlauf-Feld:** `make suchlauf-nachmessen PLAN=<Plan-Datei>` Exit 0, gedruckt
  „suchlauf-nachmessen: 16 Zeilen stimmen“. Von Hand mit `git grep` nachgefahren (ohne das Werkzeug,
  Ist gleich Soll): Zeile 1 (Parent `7b70b34a`, Kennungs-Zeilen) 2233; Zeile 2 (`-l`) 243 Dateien;
  Zeile 3 (`-l`, `:!*_test.go`) 125; Zeile 5 (`ff.`, Parent) 8; Zeile 6 (`auflösbares Feld`, Parent) 2;
  Zeile 7 (`§3\.7`, `diff`) 12; Zeile 9 (`diff`, Kennungs-Zeilen) 2236; Zeile 10 (`diff`, `-l`) 244.
  Die Zusammensetzung der Diff-Zahlen (`ff.` 8, drei Kennungen 33, 2236) siehe F-6; Zeile 8
  (`Zeile 266` am Stand `1021f6fe`) 1 (`slice-harness-fmt-check`). Der Zeilenversatz der LOW-Zeile des Reviewer-Skills ist real:
  `gofmt -l` steht am Stand `1021f6fe` auf Zeile 266, jetzt auf Zeile 286 — die Meldung des Implementers an
  die fremde Plan-Datei trifft zu.
- **Eingabeseiten-Mutationen am Programm** (je eine Kopie in einem Wegwerf-Verzeichnis, Test im gepinnten
  Toolchain-Image, Exit ungepiped; das Original blieb unverändert, `git status` sauber):

  | Nr. | Mutation an `main.go` | Ergebnis |
  |---|---|---|
  | M1 | Schwelle „mindestens zwei“ auf „mindestens drei“ | rot (`TestCandidate`, `TestRunModes`, `TestRunDiffMode`) |
  | M2 | Kompaktform-Schleife nach der ersten Nummer abgebrochen | rot (`TestScanText`) |
  | M3 | Ziffernbreiten-Prüfung der Fortsetzung entfernt | rot (`TestScanText`) |
  | M4 | Überlappungsprüfung des Diff-Modus invertiert | rot (`TestOverlaps`, `TestRunDiffMode`) |
  | M5 | Zeichenketten-Literale zusätzlich gezählt | rot (`TestBlocksOf`) |
  | M6 | Ausschluss der Wurzel `gen` entfernt | rot (`TestRunModes`) |
  | M7 | `-tests`-Filter vertauscht | rot (`TestRunModes`) |
  | M8 | Direktiven-Ausschluss (`//go:`) entfernt | rot (`TestBlocksOf`) |
  | M9 | Ende des hinzugefügten Bereichs um eins verlängert | rot (`TestParseDiff`) |
  | M10 | Diff-Filter im Lauf durch eine Bedingung ersetzt, die den Bereich nicht liest | rot (`TestRunDiffMode`) |
  | M11 | „ff.“-Erkennung entfernt | rot (`TestScanText`, `TestBlocksOf`) |
  | M12 | Exit 1 bei Kandidaten entfernt | rot (`TestRunModes`, `TestRunDiffMode`) |
  | M13 | Hunk ohne hinzugefügte Zeile (`n > 0` zu `n >= 0`) trägt einen Bereich | rot (`TestParseDiff`) |
  | M14 | Direktiven-Ausschluss nur für `//go:build` | rot (`TestBlocksOf`) |
  | M15 | Ausschluss der Wurzel `.harness` entfernt | rot (`TestRunModes`) |
  | M16 | Ausschluss auf einem Datei-Pfad-Argument entfernt | rot (`TestRunModes`) |
  | M17 | Ausschluss der Wurzel `.git` entfernt | **grün** (siehe F-8) |

- **Aufrufer `tools/harness/kommentar-kennungen.sh` von Hand** (kein Tabellentest vorhanden):
  `DIFF=<unbekannter Name>` Exit 2 mit Meldung; `TESTS=alle` Exit 2; ein fehlender Pfad Exit 2;
  fehlende `TOOLCHAIN_IMAGE`-Variable Exit 2; `PATHS='a; touch <Datei>'` und `DIFF='--output=<Datei>'`
  legen keine Datei an (Exit 2, jedes Wort ein eigenes Argument, kein `eval`); `PATHS='internal/*'` wird
  nicht expandiert; die Temp-Datei des Diffs bleibt nach Läufen nicht zurück. Der Diff-Lauf gegen `HEAD~40`
  meldet 39 Blöcke; mit `diff.mnemonicPrefix=true` in der Git-Konfiguration meldet derselbe Lauf **0**
  (F-1); mit `diff.noprefix=true` weiterhin 39.
- **Erstbeleg gegen den Code:** `git diff 1021f6fe..71cf9537 -- internal/application/port/outbound/changestream.go`
  ohne Zeile, die nicht mit `//` beginnt: **0** (Befehl: die geänderten Zeilen des Diffs ohne die
  `+++`/`---`-Köpfe, gefiltert auf Nicht-`//`). `make kommentar-kennungen PATHS=…/changestream.go` Exit 0,
  Ausgabe leer; drei Zeilen mit Kennung, alle `ADR-0060` (Parent: elf). Zusagen von `Publish` gegen
  `internal/adapters/driven/grpcstream/broadcaster.go:96-117` gelesen: nicht blockierend
  (`select … default`), begrenzte Warteschlange je Abonnent (`queueCapacity`), Überlauf verworfen für den
  eintreffenden Change, Fehler nur bei `nil`-Change (`ErrChangeStream`) oder beendetem `ctx`, an die zum
  Aufrufzeitpunkt registrierten Abonnenten; die Aufrufstelle `capture/service.go:150` fängt den Fehler
  mit einer Warnung ab und lässt ihn nicht in den Rückgabewert; Persistierung und Bestätigung liegen vor
  dem Aufruf. `ChangeStorePort` und die View `cdc.changes` existieren (`git grep`); die einzige
  `ChangeStreamPort`-Implementierung ist der `Broadcaster` (`var _ outbound.ChangeStreamPort`).
- **Stichprobe der Liste** (acht Zeilen aus `make kommentar-kennungen`, jede zwanzigste ab der zehnten
  nach `awk 'NR%20==10'`, selbst gelesen): vier Verstöße im Sinn der Regel (Reihung gleichrangiger Anker
  bzw. Kette): `systemclock.go:1-6`, `receive.go:451-453`, `table_schema.go:7-10`, `http/server.go:54-57`
  (mit „ff.“); vier **Grenzfälle**, in denen zwei Anker je eine eigene Aussage der Stelle tragen:
  `domain/errors/errors.go:35-37`, `port/inbound/consumer.go:30-32`, `queries/queries.go:241-244`,
  `port/inbound/verwaltung.go:83-89` (siehe F-7).
- **Commits:** alle sieben Betreffs nennen `ADR-0083`, keiner trägt `SPEC-*`/`ARC-*`; die beiden
  Lifecycle-Moves (`ec6cb07e`, `0d333120`) sind reine Renames (0 Zeilen), Inhalt in eigenen Commits
  (`AGENTS.md` §3.3). Keine Spur von `sed -i`/`perl -pi` in hinzugefügten Zeilen.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | Der `DIFF`-Pfad zerlegt die `+++`-Zeilen des `git diff` mit festem Präfix `b/`; `diff.mnemonicPrefix=true` in der Git-Konfiguration des Aufrufers lässt `git diff` `w/…` statt `b/…` ausgeben (ein abweichendes `diff.dstPrefix` wirkt gleich, nicht gemessen). Gemessen: derselbe Lauf gegen `HEAD~40` meldet 39 Kandidaten, mit `diff.mnemonicPrefix=true` **0** bei Exit 0 — das Werkzeug meldet „kein Kandidat“, ohne dass ein Block gelesen wurde. Der Aufrufer pinnt `--no-color --no-ext-diff`, nicht die Präfix-Optionen; weder Programm noch Aufrufer melden einen Diff, in dem kein Pfad einer Datei des Arbeitsverzeichnisses entspricht. | Maintainability; `regel-weiter-als-ihr-sensor`-Restrisiko „Werkzeug grün, Muster trifft nicht“ (Plan §6); Vertrag §Grenze 4 (grün ist kein Beleg) | `tools/harness/kommentar-kennungen.sh:48`; `tools/harness/kommentar-kennungen/main.go:166-171` | ja — `GIT_CONFIG_COUNT=1 GIT_CONFIG_KEY_0=diff.mnemonicPrefix GIT_CONFIG_VALUE_0=true` vor `bash tools/harness/kommentar-kennungen.sh "" 1 "" HEAD~40` (mit gesetzter `TOOLCHAIN_IMAGE`), Ausgabe 0 statt 39; der Tabellentest hat keinen Fall dafür (der Diff-Strom kommt dort als Literal) | Werkzeug liest Nutzer-Konfiguration ohne Pin — grün ohne Aussage |
| F-2 | LOW | Der Aufrufer trägt keinen Tabellentest; der Vertrag nennt das („Nicht gebunden ist der Aufrufer“). Seine Zweige (unbekannte Diff-Basis, `TESTS`-Wert, Temp-Datei, Argument-Zerlegung, Exit-Weitergabe) habe ich von Hand gefahren und in Ordnung gefunden; die eine Naht, die der Go-Test nicht deckt — der Diff-Strom aus `git` —, ist die von F-1. Die Klasse „fehlende Negativtests bei neuem öffentlichem Vertrag“ wäre MEDIUM; LOW, weil die Lücke im Vertrag benannt ist und der Plan den Go-Test als Form der Bindung setzt. | Maintainability; fehlende Negativtests bei neuem öffentlichem Vertrag (Reviewer-Skill) | `tools/harness/kommentar-kennungen.sh:33-70`; `harness/sensors/kommentar-kennungen.md` §Test | ja — `DIFF=<unbekannt>` und `TESTS=alle` real Exit 2 (gefahren), aber kein Test hält das fest | Vertrags-Exit-Zweige des Aufrufers ohne Testbindung |
| F-3 | LOW | Der Godoc von `ErrChangeStream` sagt im ersten Satz, der Fehler „trägt die Fehlerklasse“ des Stream-Publish, und im dritten, er trage „keine der Fehlerklassen (`transient` und die übrigen der Fehlerklassifikation)“. Die Zusage ist wahr (die Sentinel-Variable trägt keine Klasse; `ErrNotify` trägt `transient`), der Wortlaut liest sich widersprüchlich; der Rang-Zeiger auf die Fehlerklassifikation, der die „übrigen“ auflösbar machte, ist mit der Bereinigung auf einen Anker entfallen. | `AGENTS.md` §3.7 (Zusage, Rang-Zeiger) | `internal/application/port/outbound/changestream.go:10-17` | nein | Kommentar zweideutig, Rang-Zeiger der Abgrenzung verloren |
| F-4 | LOW | Zwei Einfügungen trennen ein bestehendes Bezugspaar: In `AGENTS.md` §3.7 steht der neue Absatz zwischen dem Falsch/Richtig-Paar der Kommentar-Regel und „**Zustandsfelder ebenso:**“, dessen „ebenso“ die Kommentar-Regel meint, nicht die Herkunfts-Regel. In `implement-slice.md` Schritt 20 steht der neue Block vor „**Grenze dieser Selbstprüfung**“, dessen „Dieser Schritt“ und „bislang 4/4“ den Chronik-Lauf meinen; der neue Block trägt daneben eine zweite „**Grenze:**“ mit derselben Aussage („erste, nicht tragende Linie, die tragende ist der Reviewer“). | Maintainability; Nachzug widerspricht dem Nachbarn im selben Träger (Reviewer-Skill; hier ohne Widerspruch der Aussagen, nur Bezug unklar) | `AGENTS.md:182-211` und Folgeabsatz „Zustandsfelder ebenso“; `.claude/commands/implement-slice.md:243-263` | ja — Lesen des Kontexts (`git diff -U20`) | Einfügung trennt Bezugspaar |
| F-5 | LOW | Die 31 in `AGENTS.md` eingefügten Zeilen verschieben §3.13 von Zeile 449 auf 480; der Register-Beleg `evidence/welle-d-check-verkoerperung.md` (Eintrag `zitat-nennt-die-falsche-stelle`) zitiert „`AGENTS.md:451-455` (korrigiert)“ und meint den Anfang von §3.13 — die Zeilen 451-455 stehen jetzt am Ende von §3.12. Das Suchlauf-Feld des Plans sucht `§3\.7`, nicht das Muster `AGENTS.md:<Zahl>`; die Nachzug-Tabelle nennt den Träger nicht. Träger außerhalb des Diffs, deshalb LOW und Meldung an den Planner statt stiller Änderung (`AGENTS.md` §3.13). | `AGENTS.md` §3.13 (Träger-Nachzug; Zeilen-Lokatoren fängt nur der Reviewer) | `docs/plan/planning/observations/BEO-PGC/zitat-nennt-die-falsche-stelle/evidence/welle-d-check-verkoerperung.md:28` | ja — `git grep -n -E 'AGENTS\.md:[0-9]+' -- . ':!docs/reviews'` (ein Treffer) und `grep -n '^### 3.13' AGENTS.md` (Zeile 480 statt 449 am Stand `1021f6fe`) | Zeilen-Lokator in fremdem Träger driftet durch Einfügung |
| F-6 | INFO | Zwei `diff`-Zeilen des Suchlauf-Felds stimmen in der Zahl, nicht in der Zusammensetzung, die der Plan nennt: „`ff.`“ 8 → 8 ist 7 Bestandszeilen (der Erstbeleg nimmt eine) plus eine Testfixture (`main_test.go:88`, ein Quelltext-Literal mit `// … ff.`); „drei Kennungen“ 31 → 33 sind 31 plus zwei Fixture-Zeilen (`main_test.go:70`, `:79`); 2233 → 2236 sind 2233 − 8 (`changestream.go`: elf Zeilen mit Kennung → drei) + 11 (`main_test.go`). Das Feld druckt die Zahlen richtig; dass die Wirkung des Erstbelegs auf „`ff.`“ in der Zahl nicht sichtbar ist, steht nirgends. | `AGENTS.md` §3.12 (Zahl im Träger, Ursprung) | `docs/plan/planning/in-progress/slice-code-kommentare-kennungen.md` §3 Suchlauf-Feld, Zeilen `diff 8`/`diff 33`/`diff 2236` | ja — `git grep -n -E '//.* ff\.' -- '*.go'` listet die acht Fundstellen | Zahl gleich, Zusammensetzung anders |
| F-7 | INFO | Die Stichprobe von 30 Kandidaten mit Klasse und Urteil, die Risiko 1 des Plans (§6) belegen soll, steht „im Bericht des Implementers“, in keinem committeten Träger (der Plan nennt Anzahl und Anteil, nicht die Liste). Meine acht Zeilen derselben Auswahlregel: vier Verstöße, vier Grenzfälle mit zwei Ankern und je einer eigenen Aussage (z. B. `errors.go:35-37`: die Vorwärts-Regel des ACK und die Idempotenz derselben Position); die Regel „höchstens eine Kennung“ nimmt dort einem der beiden Sätze seinen Anker. 4 von 8 liegen über dem Anteil der Implementer-Stichprobe (5 von 30); acht Zeilen sind keine Stichprobe, die ihn widerlegt. Die Regel ist die Zielform des Plans (§1) und strenger als der Baseline-Wortlaut „ein auflösbares Feld“; die Rückführung (a) misst die Gesamtquote (40,5 %), nicht den Grenzfall-Anteil. | Plan §6 Risiko 1 und 4; Baseline `v6.9.0` · `regelwerk/grundlagen-harness-dateien.md` §Was ein Kommentar trägt | Plan §3 „Stichprobe“; `internal/domain/errors/errors.go:35-37`, `internal/application/port/inbound/consumer.go:30-32`, `internal/adapters/driven/postgresstorage/queries/queries.go:241-244`, `internal/application/port/inbound/verwaltung.go:83-89` | nein | Stichproben-Beleg ohne committeten Träger; Regel enger als der Baseline-Wortlaut |
| F-8 | INFO | Der Ausschluss der Wurzel `.git` ist im Vertrag genannt und im Programm gesetzt, aber nicht an den Test gebunden (M17 grün); `gen`, `sdks`, `.harness` sind gebunden. Ein `.go` unter `.git` gibt es im Normalfall nicht. | Maintainability | `tools/harness/kommentar-kennungen/main.go:45`; `main_test.go` `TestRunModes` | ja — M17 an einer Kopie | Ausnahme-Eintrag ohne Testbindung |
| F-9 | INFO | Die Basis-Messung des Plans (600/403/197) steht am Stand `0d333120`, an dem das Werkzeug nur im Arbeitsbaum lag; sie ist mit keinem Repo-Befehl an diesem Stand wiederholbar, ihre Nachmessung am Diff-Stand (597/400/197, Nenner 1474) stimmt. Der neue Reviewer-Unterpunkt nennt als Herkunft `slice-chronik-in-code-kommentar`, der Plan (§7) erwartet als Register-Eintrag `kommentar-herkunft-als-kette`; beides sind Closure-Nachzüge des Planners. | `AGENTS.md` §3.12 | Plan §3 „Basis-Messung“; `.harness/skills/reviewer.md:282-283` | nein | Zahl an einem nicht wiederholbaren Stand |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `tools/harness/kommentar-kennungen/main.go` (Zählregel, Blockgrenzen, Kennungsalphabet, Diff-Parser, Modi, Exit-Codes) | geprüft, ohne Befund über F-1, F-8 hinaus: Kennungsmuster stimmen mit den `ids.patterns` und `commits.id-patterns` in `.d-check.yml` überein (kein Kennungs-Fund im Baum mit anderer Ziffernbreite); Blockgrenzen sind die des Go-Parsers (Endkommentar eigener Block, Direktive zählt nicht, Literal ist kein Kommentar, Blockkommentar); dieselbe Kennung zweimal zählt einmal; Kompaktform und Bereich je Nummer; Löschung, gelöschte Datei (`/dev/null`) und andere Datei tragen im Diff-Modus keinen Bereich; ein nicht lesbarer Quelltext oder Hunk-Kopf gibt Exit 2; die Zeichenform „ff.“ direkt an der Kennung ohne Leerzeichen erkennt das Programm nicht (Grenze nach Vertrag: Formen außerhalb der genannten bleiben unerkannt). 16 von 17 Mutationen rot |
| `tools/harness/kommentar-kennungen/main_test.go` | geprüft, ohne Befund über F-8 hinaus: Fälle je Eingabeseite (Zählregel, Blöcke, Diff-Bereiche, Überlappung, Modi, Ausschlüsse, Exit-Codes) und unterscheidbare Erwartungswerte (M1, M12 rot: eine Schwellen- oder Exit-Mutation fällt an mehreren Fällen) |
| `tools/harness/kommentar-kennungen.sh` (Docker-Aufruf, Temp-Datei, Argument-Zerlegung, Exit-Weitergabe) | geprüft, ohne Befund über F-1, F-2 hinaus: `--network none`, gepinntes Toolchain-Image, Repo lesend gemountet, `go build` + `exec` (Exit 1 und 2 bleiben unterscheidbar), Diff in einer Temp-Datei mit geprüftem `git`-Exit (`AGENTS.md` §3.9), `trap` löscht die Datei, `PATHS`/`DIFF`/`TESTS` erreichen das Programm nur als einzelne Argumente (kein `eval`, keine Expansion, keine Option, die ein Kommando startet), Diff-Basis wird vor dem Aufruf als Commit geprüft |
| `harness/sensors/kommentar-kennungen.md` | geprüft, ohne Befund: Vertrag nennt Form-statt-Wahrheit im ersten Absatz, kein Gate, keine Ausnahmeliste, Abgrenzung zum Chronik-Lauf (kein Satz-Subjekt, keine Slice-/Wellen-Nummer), Host-Werkzeuge, Exit-Tabelle (über `make` Exit 2, `Fehler <n>` — gefahren) und die fünf Grenzen; die Beschreibung von `sdks/` und `gen/` als ausgenommen stimmt mit Programm und Plan überein; das Werkzeug wertet weder ein Satz-Subjekt noch eine Slice-Nummer aus — das Verdikt (kein Textmuster-Sensor über Satz-Subjekte) wird nicht umgangen |
| `Makefile` (zwei Ziele), `harness/README.md` (zwei Zeilen) | geprüft, ohne Befund: `.PHONY`, Hilfe-Zeilen, kein Eintrag in einem Gate-Bündel, `make test` deckt das Paket (45 Pakete `ok`), `make test-kommentar-kennungen` netzlos Docker-only; README-Zeilen nennen Aufruf, Grenze und Bindung im Indikativ |
| `AGENTS.md` §3.7 (Wortlaut) | geprüft, ohne Befund über F-4 hinaus: konkret (eine Kennung, keine Kette, keine Kompaktform, kein „ff.“, keine Spec-Wiedergabe, Kopplung nennt die Stelle) und prüfbar; widerspricht der Baseline nicht (`v6.9.0` · `regelwerk/grundlagen-harness-dateien.md` §Was ein Kommentar trägt: „ein auflösbares Feld … nie als Absatz“ — der Wortlaut konkretisiert es, F-7 nennt die Strenge); das „Richtig“-Beispiel trägt genau eine Kennung, das „Falsch“-Beispiel steht im Codeblock; keine Zeilen-Lokatoren auf §3.7; der eine Lokator auf `AGENTS.md` außerhalb der Records trifft §3.13 und driftet (F-5) |
| `.claude/commands/implement-slice.md` Schritt 20 | geprüft, ohne Befund über F-4 hinaus: Aufruf `make kommentar-kennungen DIFF=<Basis>` mit dem Hinweis auf `git add` neuer Dateien, Chronik-Lauf bleibt, Grenze genannt; die Nummern der Schritte 20/21 und alle Verweise auf sie (`grep` „Schritt 2[01]“) unverändert |
| `.harness/skills/reviewer.md` (neuer MEDIUM-Unterpunkt) | geprüft, ohne Befund über F-9 hinaus: Einstufung MEDIUM für die **Form** und Eskalation auf den HIGH-Punkt „Kommentar trägt keine der Kommentar-Klassen“ bei falscher Wiedergabe oder nicht getragenem Verhalten ist begründet und doppelt den HIGH-Punkt nicht (dieser deckt Konjunktiv, abwesenden Text, Zusage ohne Träger — nicht die Kennungszahl); Lauf ist „Probe, kein Beleg“; der Skopus „im Diff neu geschriebener oder geänderter Kommentar“ passt zu Vertrag §Grenze 5 (eine berührte Zeile bringt den Block in den Lauf); Zeilenversatz der LOW-Zeile (266 → 286) gemeldet und zutreffend |
| `internal/application/port/outbound/changestream.go` | geprüft, ohne Befund über F-3 hinaus: nur Kommentarzeilen geändert (Nicht-`//`-Zeilen im Diff: 0), `gofmt -l` leer, ein Anker je Block, die Zusagen von `Publish` (verteilt an die zum Aufrufzeitpunkt registrierten Abonnenten, blockiert nie, begrenzte Warteschlange je Abonnent, Überlauf verworfen als Drop-Newest, keine Zustellgarantie, Fehler nur bei ungültigem Aufruf oder beendetem `ctx`, keine Rückwirkung auf Persistierung und Bestätigung) stehen vollständig und sind gegen `Broadcaster.Publish` und die Aufrufstelle wahr; Verweis auf den Lesezugriffsweg (`ChangeStorePort`, View `cdc.changes`) ohne Spec-Wiederholung; nichts zugesagt, was der Adapter nicht tut; die Zusagen zum optionalen Port und zum `CaptureService` stimmen mit `capture/service.go` |
| Plan `slice-code-kommentare-kennungen` §3 (Nachzug, Basis-Messung, Träger-Meldungen) | geprüft, ohne Befund über F-5, F-6, F-9 hinaus: die Abweichungen vom Plan-Wortlaut (Temp-Datei statt Pipe, `go build` + `exec`, Kompaktform mit `…`) stehen begründet in der Nachzug-Tabelle; Meldungen an `slice-harness-fmt-check` (Zeile 266) und an die fremden Träger (Bereinigungs-Plan, Register) tragen Frist und Adresse; `.a-check.yml` unverändert und richtig (Gruppe `tooling`, nur Standardbibliothek; `make a-check` Exit 0); Coverage-Fläche unverändert (`Dockerfile`-Stufe `coverage` zählt `./internal/...`, `./cmd/...`, `./gen/...`; `tools/` nicht Teil; Schwelle 80 %) |
| Hard Rules und Commit-Struktur | geprüft, ohne Befund: keine Suppression (kein `//nolint`), kein host-lokaler absoluter Pfad in den neuen Dokumenten, keine Toolchain-Installation auf dem Host, `git mv` rein und Inhalt getrennt, Betreffs mit `ADR-0083` ohne `SPEC-`/`ARC-`; das Benutzerhandbuch bleibt unberührt und muss es (keine Betreiber-Oberfläche, keine neue `CDC_*`-Variable, keine `cdc.*`-Funktion); keine Spec-/Lastenheft-Änderung |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 4 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Werkzeug liest Nutzer-Konfiguration ohne Pin — grün ohne Aussage ·
Vertrags-Exit-Zweige des Aufrufers ohne Testbindung · Kommentar zweideutig, Rang-Zeiger der Abgrenzung
verloren · Einfügung trennt Bezugspaar · Zeilen-Lokator in fremdem Träger driftet · Zahl gleich, Zusammensetzung anders · Stichproben-Beleg ohne
committeten Träger · Ausnahme-Eintrag ohne Testbindung · Zahl an einem nicht wiederholbaren Stand

## Verdikt

**Merge-blockierend:** ja — wegen F-1 (MEDIUM): der Diff-Modus, der dem Implementer und dem Reviewer als
Probe dient, meldet unter einer gängigen Git-Konfiguration „kein Kandidat“, ohne zu lesen; das ist die
Ausprägung des Restrisikos „Werkzeug grün, Muster trifft nicht“ aus Plan §6, die das Werkzeug ehrlich
lässt oder nicht. Kein HIGH: das Werkzeug ist advisory (kein Gate), der Erstbeleg trägt seine Zusagen
gegen den Code, Regel, Träger und Vertrag sind konsistent, alle Läufe (`make test`, `make a-check`,
`make coverage-gate`, `make gates`) grün. Die LOW-Findings (F-2 bis F-5) und die INFO-Findings gehen ohne
eigene Fixrunde an den Planner bzw. den Implementer der Fixrunde von F-1 mit.

**Übergabe:** F-1 (und, in derselben Runde mitgenommen, F-2 bis F-4) gehen an den Implementer; F-5 geht als Träger-Meldung an den Planner (Frist: die Closure dieses Slice); die
**Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Die DoD-Zeile
„Review durchgeführt“ im Slice-Plan bleibt offen (eine Fixrunde folgt; sie wird bei Schritt 21 des
Implementer-Ablaufs nachgezogen). Dieser Report selbst ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser
Skill, dieses Modell, dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und muss es nicht.
Der Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat (Modul 11;
anderes Prüf-Artefakt, anderer Eingabe-Kontext).
