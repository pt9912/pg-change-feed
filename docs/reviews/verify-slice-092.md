# Verifikationsbericht: slice-092 — 2026-09-16

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den DoD-Vertrag
(`slice-092` §2, LP1–LP3), die im Slice referzierten Entscheidungen
[`ADR-0082`](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
(Schnittmaß, §Konsequenzen „test-only", Re-Evaluierungs-Trigger),
[`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(Messgegenstand), [`ADR-0077`](../plan/adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md)
(Rampe) sowie die Hard Rules `AGENTS.md` §3.1, §3.6, §3.7, §3.9, §3.11,
§3.12. **Nicht** gegen den Diff als solchen (Reviewer-Aufgabe, mit
dem Review zu `slice-092` abgeschlossen) und **nicht**
gegen realen Bedarf (Validator, nicht ausgelöst).

**Frischer Kontext.** Der Slice-Plan wurde am Stand `HEAD` vollständig
gelesen (§1–§8), dazu `ADR-0082`, der Review-Report, die drei Commits und die
zwei berührten Register-Einträge. Implementer-Bericht und Review waren
**Kontext**, ihre Zahlen **nicht** übernommen: jede Zahl dieses Berichts
stammt aus einem hier selbst gefahrenen Lauf — einschließlich der vier
Review-Findings, die ich nicht übernommen, sondern nachgemessen habe (F-1 an
allen vier Zahlen, F-2 am Träger und an der Messung, F-3 an der
**Parent-Fassung**, siehe Probe P-7). Exit-Codes sind je **ungepiped** und in
eigenem Schritt gelesen (`AGENTS.md` §3.9); Gate-Lauf und Auswertung waren
getrennt beauftragt. Logs und Arbeitsbaum-Kopien liegen **außerhalb** des
Baums; der Baum ist vor und nach jedem Lauf sauber.

**Gegenstand:** `HEAD` = `ad95444`, Zweig `main`, Baum sauber
(`git status --porcelain` leer — vor und nach jedem Lauf). Die drei Commits
des Vorgangs: `3de9547` (Tests), `d7b50e1` (Review-Report), `ad95444`
(Korrektur F-1/F-2). Der Slice liegt in `in-progress/`; der `git mv` nach
`done/` ist **nicht** erfolgt.

**Zwei Zahlen zur Schreibweise.** (a) Das Gate-Skript endet im Rot-Fall mit
**Exit 1**, `make` kapselt den Rezept-Fehlschlag zu **Exit 2**; die
Mutationsproben dieses Berichts laufen als direkter `go test`-Aufruf und
stehen darum als **1**. (b) Alle Prozentzahlen sind die **gedruckten** Zeilen
eines konkreten Laufs; wo ich einen Prozentwert aus einer Statement-Zahl
**zurückrechne**, steht das dabei (§3.12 Instanz A).

---

## 1. Eigene Messungen dieses Laufs

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `make gates` (Log in Datei, Exit **danach** aus eigener Datei) | **0** | sechs Checks: `baseline-verify: v6.5.0 OK — 54 Dateien` · `coverage-gate: OK — Coverage 78.60% erfüllt Schwelle 70%` (gedruckte `total:`-Zeile derselben Stufe: `78.6%`) · `d-check: 758 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s)` · `generated-sync: OK` (beide `.pb.go` byte-gleich) · `a-check: gesamt: 0 Befund(e)`; danach Baum leer |
| 2 | `make test` (netzlos, `--network none`, `-race ./...`) | **0** | **32 × `ok`, 0 × `FAIL`**; darunter alle **13** Use-Case-Pakete und `domain/model` je `ok` |
| 3 | `make doc-commits RANGE=3de9547^..ad95444` | **0** | 758 Dateien, 0 Befunde |
| 4 | `make doc-immutable RANGE=3de9547^..ad95444` | **0** | 758 Dateien, 0 Befunde |
| 5 | Profil der netzlos prüfbaren Fläche am Diff-Stand (Zählbasis der `coverage`-Stufe: eigene `go list`-Liste, `-coverpkg` über den Gegenstand, `-covermode=atomic`, dedupliziert über die Block-Position), netzlos | **0** | `TESTEXIT=0`, **31** Pakete im Gegenstand, gedruckt `78.5%`; dedupliziert: **Nenner 1903 · gedeckt 1493 · ungedeckt 410**, **226** ungedeckte Blöcke |
| 6 | dasselbe, **6 ×** wiederholt (derselbe Quelltext) | **0** (6 ×) | gedeckt ∈ {**1493** (5 ×), **1495** (1 ×)}, ungedeckt ∈ {**410**, **408**}; gedruckt 5 × `78.5%`, 1 × `78.6%` |
| 7 | dasselbe am **Parent** `2b7c6ec` (Arbeitsbaum-Kopie `git archive`, außerhalb des Repos) | **0** | gedruckt `77.3%`; **Nenner 1903 · gedeckt 1471 · ungedeckt 432**, 249 ungedeckte Blöcke |
| 8 | Block-Positions-Vergleich Parent → Diff (Mengendifferenz der ungedeckten Positionen) | **0** | **24** Blöcke neu gedeckt, **je 1 Statement = +24** (22 in `usecase/*/service.go`, 2 in `domain/model`: `source.go:31.2,31.40`, `transaction.go:29.3,30.1`); **1** Block neu ungedeckt = **2** Statements (`internal/bootstrap/wiring.go:991.5,992.13`); **Netto +22** |
| 9 | paketweise am Diff-Stand (dedupliziert) | **0** | **13** Use-Case-Pakete + `domain/model` **sämtlich `uncovered = 0`** (Tabelle §2/LP3); die einzigen Pakete mit Rest: `bootstrap` **336**, `cmd/pg-change-feed` **49**, `streamv1` **10**, `replication/mapper` **7**, `http` **3**, `natsnotify` **2**, `telemetry` **2**, `replication/decode` **1** — Summe **410**, keines davon in D1 |
| 10 | paketweise am **Parent** `2b7c6ec` | **0** | **10** der 13 Use-Case-Pakete mit `uncovered > 0` (Summe über die Pakete: **22**), `domain/model` **2** → **24**; `readchanges`, `excludecolumn`, `includecolumn` je **0** |
| 11 | `go list ./internal/application/usecase/...` · Gegenstandsliste · Gruppierung | **0** | **13** Use-Case-Pakete; Gegenstand **31** Pakete; `Test=0 XTest=0` **4**; `TestGoFiles=0` **23** (→ **19** mit externem Testpaket); md5 der `go list -f`-Ausgabe Parent **und** Diff **identisch** (`adda7355…`) |
| 12 | Umfang: `git diff --name-status 3de9547^..ad95444` | **0** | **17** Pfade: **13** `*_test.go` (**12** Use-Case + `domain/model/validation_test.go`), **3** Doku (`slice-092`-Plan, `welle-20.md` §4, `harness/sensors/coverage-gate.md`), **1** Review-Report. Pfade **ohne** `_test.go` und **außerhalb** `docs/`: **genau einer** — `harness/sensors/coverage-gate.md`; **kein** Produktcode |
| 13 | `THRESHOLD` | **0** | `harness/mk/coverage.mk:15` steht auf `THRESHOLD ?= 70`; die Datei ist im Vorgang **nicht** berührt |
| 14 | **7 Mutationsproben + 7 Kontrollläufe** (§6, Tabelle unten) | je **1** (Probe) / **0** (Kontrolle) | 6 Proben an **neuen** Tests: **6 × rot**; 1 Probe an der **Parent-Fassung** der reparierten Spaltenbindung: **grün** (P-7); alle 7 Kontrollen am unmutierten Stand: **Exit 0** |
| 15 | Register-Stände: `ls evidence/ | wc -l` | **0** | `arbeit-ueberholt-stehenden-traeger` **1** · `negativtest-ohne-bindung-an-seine-eingabe` **5**; **keine** Registerdatei im Vorgang geändert (`git diff --name-only … | grep -c observations` → **0**) |
| 16 | `ls docs/plan/planning/reconciliation.md` | **1** | Datei existiert nicht (Greenfield) → das §2-Item „entfällt" ist nachweislich richtig |
| 17 | erreichte Zahlen in stehenden Trägern: `grep -rn "1493\|1495\|78,5\|78.5\|78,6\|78.6\|408\|410"` über `*.md` außerhalb `docs/reviews/` und der vendored Baseline | **0** | **kein Treffer** in einem Doku-Träger (der einzige Roh-Treffer ist ein Actions-SHA in `AGENTS.md`); ebenso `grep -rn "D1"`: nur `welle-20.md:100` (Soll-Zeile) und der Slice-Plan |
| 18 | DoD-Häkchen · §6-Ausgänge · §7-Platzhalter (§2 des Plans, am Stand `HEAD`) | **0** | **0 von 11** Häkchen gesetzt; **4 von 4** §6-Risiken auf `— **Ausgang:** <…>`; **6 von 6** §7-Zeilen Platzhalter |
| 19 | `make gates` **erneut, mit diesem Bericht im Baum** (Log in Datei, Exit danach aus eigener Datei gelesen) | **0** | sechs Checks grün; `d-check: **759** Datei(en) geprüft, **0** Befund(e)` (eine Datei mehr = dieser Bericht); `coverage-gate: OK — Coverage **78.50%** erfüllt Schwelle 70%` — der zweite Gate-Lauf desselben Stands liegt am **anderen** Schwankungsende als #1 |

---

## 2. DoD-Konformität, Kriterium für Kriterium

### Liefer-Punkt 1 — die Tests existieren und sind netzlos grün

| Kriterium (§2) | Befund |
|---|---|
| für die **12** geänderten Use-Case-Pakete und `domain/model` liegen Tests vor | **erfüllt** — der Vorgang trägt **13** `*_test.go`: die **12** Use-Case-Pakete `acknowledge`, `capture`, `disable`, `enable`, `excludecolumn`, `includecolumn`, `list`, `position`, `register`, `remove`, `retention`, `status` plus `domain/model/validation_test.go` (#12); `readchanges` ist zu Recht nicht dabei (am Parent bereits `0/10` ungedeckt, #10) |
| der Gate-Lauf **fährt sie wirklich** | **erfüllt** — die `coverage`-Stufe bildet ihre Paketliste aus `go list ./internal/... ./cmd/...` (31 Pakete, alle 13 + `domain/model` darunter) und bricht bei rotem Testlauf ab; mein Nachbau derselben Liste netzlos: `TESTEXIT=0` (#5). Zusätzlich `make test` netzlos mit Race-Detector: **Exit 0, 32 × `ok`** (#2) |
| `make gates` ist grün | **erfüllt** — **Exit 0** aus separater Datei gelesen, sechs Checks grün (#1) |
| **der Zuwachs wird als Zahl mit ihrem Lauf genannt** | **materiell erfüllt, Träger offen** — gemessen: **+24** Statements als **24** Block-Positionen (22 + 2), **Netto +22**, Nenner 1903 unverändert (#8). Die Zahl steht heute **nur** in der Commit-Message (git-Historie — laut `ADR-0083` §Geltungsbereich **kein** Doku-Träger) und im Review-Report (Lauf-Beleg). §7 ist leer (#18). Die Zahl samt Lauf für §7 steht in **§5** dieses Berichts |

### Liefer-Punkt 2 — die Negativtests binden ihre Ablehnung an die Eingabe

| Kriterium (§2) | Befund |
|---|---|
| geprüft durch **Mutation an der Produktionsseite** | **erfüllt** — **6** eigene Proben an **neuen** Tests des Diffs, jede **rot**; jede mit eigener **Kontrolle** am unmutierten Stand (**Exit 0**) (#14, Tabelle unten) |
| die Ablehnung hängt am **Eingabewert**, nicht am Fake | **erfüllt** (Stichprobe von 6) — die Proben greifen die drei in LP2 genannten Klassen ab: **fremde/fehlende Tabelle** (P-2: `TableExists` mit fester Adresse → `TestDisableTableExistsErrorFollowsAddress` rot), **unbekannte Spalte** (P-6: `ColumnExists` mit fester Spalte → `TestExcludeColumnRejectsMissingSourceColumn` rot), **ungültige Dauer** (P-8: `NewDuration`-Riegel aus → `TestNewRetentionPolicyRejectsNegativeDuration` rot), dazu der Port-Fehler-Pfad des Zustands-Schreibens (P-1), die Publication-Adresse des Listenpfads (P-3) und die Freigabe-Menge der Retention (P-4) |
| die **Gegenrichtung** jeder Probe ist mitgemessen | **erfüllt** — P-2/P-3/P-6 mutieren genau die **Weitergabe der Kommando-Werte** (Adresse, Publication, Spalte); der Test färbt rot, **weil** der Produktionscode aufhört, den Eingabewert zu benutzen — nicht, weil ein Fake unabhängig davon ablehnt. Die Kontrollläufe derselben Tests am unmutierten Stand sind grün (#14) |
| ein Test, der **nicht** bindet, wäre sichtbar | **bestätigt an der Parent-Fassung** — P-7 wendet **dieselbe** Mutation (feste Spalte) auf den **Parent-Stand** an: `TestExcludeColumnRejectsMissingSourceColumn` bleibt **Exit 0** (grün). Der Commit-Text „zwei **vorbestehende** Tests, die genau die Klasse waren" ist damit **nachgemessen**, nicht übernommen (Review-F-3) |

### Liefer-Punkt 3 — die unerreichbaren Statements sind benannt

| Kriterium (§2) | Befund |
|---|---|
| jedes Paket, das danach noch ungedeckte Statements hat, nennt sie einzeln mit Grund | **erfüllt (leer erfüllbar)** — nachgemessen: **alle 13** Use-Case-Pakete und `domain/model` stehen bei **`uncovered = 0`** (#9). Für den D1-Umfang gibt es **nichts zu benennen**; die Aussage des Commits („im D1-Umfang bleibt kein ungedecktes Statement") hält. „Rest nicht erreichbar" ohne Namen trifft niemanden |
| der Umfang der Aussage ist nicht zu weit | **erfüllt** — die paketweise Messung am Diff-Stand weist **8** Pakete **außerhalb** D1 mit Rest aus (allen voran `bootstrap` 336 und `cmd` 49, #9); sie gehören laut §1 zu D2/A und sind **nicht** Gegenstand dieses Slice — die Aussage „im D1-Umfang" ist damit wörtlich richtig und nicht als „im Gegenstand" zu lesen |

| D1-Paket | ungedeckt / gesamt (Diff-Stand, #9) |
|---|---|
| `usecase/acknowledge` · `capture` · `disable` · `enable` · `excludecolumn` · `includecolumn` · `list` · `position` · `readchanges` · `register` · `remove` · `retention` · `status` | **0 / 10 · 42 · 22 · 18 · 7 · 7 · 15 · 7 · 10 · 8 · 7 · 18 · 19** |
| `domain/model` | **0 / 110** |

### Die Closure-Pflichten aus §2

| Kriterium (§2) | Befund |
|---|---|
| `make gates` grün | **erfüllt** (#1) |
| Review durchgeführt, Report unter `docs/reviews/`, **Delta-Review bei Fixrunde** | **nicht erfüllt** — das Review zu `slice-092` liegt vor (0 HIGH / 1 MEDIUM / 2 INFO) und beurteilt `3de9547`; der Review weist mit F-1 eine **Fixrunde** aus („Merge-blockierend: ja", Rückgabe-Pfeil nötig, „die Fixrunde deckt nach §2 des Plans ein **Delta-Review** ab"), und `ad95444` führt sie aus — für diesen Stand existiert **kein** Review-Artefakt (`ls docs/reviews/ | grep 092` → nur das Review zu `slice-092`) → **V-2** |
| Verifikation durchgeführt, Report unter `docs/reviews/` | **erfüllt** mit diesem Bericht; das Häkchen ist offen (#18) |
| Closure-Notiz mit Steering-Loop-Lerneintrag | **nicht erfüllt** — §7 trägt in **6 von 6** Inhaltszeilen Platzhalter (#18) |
| Reconciliation-Register fortgeschrieben *(entfällt …)* | **entfällt nachweislich** (#16) |
| Beobachtungs-Register fortgeschrieben — **kein Zähler wird gesetzt** | **nicht erfüllt** — im Vorgang ist **keine** Registerdatei geändert (#15). Die **Substanz** liefert dieser Lauf: `arbeit-ueberholt-stehenden-traeger` hat mit diesem Vorgang seine **zweite** Gelegenheit (F-1: vier Träger trugen eine Fehlzählung; F-2: ein Deixis-Satz im Sensor-Träger), `negativtest-ohne-bindung-an-seine-eingabe` einen **Reparatur-Vorgang** (§2/LP2). Beides gehört als Beleg-Entscheidung in die Closure — der Zähler folgt den Dateien, §7 notiert ihn |
| Jedes Risiko aus §6 trägt einen Ausgang | **nicht erfüllt** — 4 von 4 stehen auf `<…>` (#18). Material für drei der vier liefert dieser Lauf: **R1** (die `22`/`2` seien eine Über-Schätzung) ist **nicht eingetreten** — beide Zahlen halten exakt (#8, #10); **R2** (Negativtest bindet an den Fake) ist an den **neuen** Tests **nicht eingetreten** (6 Proben rot) und an **zwei vorbestehenden** Stellen **eingetreten und behoben** (P-7); **R3** (Coverage-Theater) trägt seinen Gegenbeleg (#14); **R4** (überholter Träger) ist **eingetreten** und benannt (F-1/F-2) |
| Die drei Paarungen sind getragen | **nicht erfüllt** — fällt laut Zeile der `welle-20`-Closure zu |

**Ergebnis:** Die **drei Liefer-Punkte** und die Gate-Zeile **tragen** — gemessen,
nicht gelesen. Offen sind die regulären Closure-Pflichten (Häkchen, §7 samt
Zahl und Lauf, Register-Entscheidung, vier §6-Ausgänge, Paarungen) und **ein**
Punkt, der nur diese Rolle sieht: **V-2** (die Fixrunde ohne Delta-Review).

---

## 3. Entscheidungs-Konformität — hält der Vorgang, was Plan §1/§3 und `ADR-0082` zusagen?

| Zusage | Befund |
|---|---|
| **Plan §1:** „Was dieser Slice liefert: Tests … Er ändert **keinen** Produkt-Code" | **eingehalten** — **17** Pfade im Vorgang, **13** davon `*_test.go`; der einzige Pfad ohne `_test.go` **außerhalb** `docs/` ist `harness/sensors/coverage-gate.md` (Doku) (#12). **Keine** Zeile `internal/**`/`cmd/**` außerhalb von Tests, keine Naht, kein `git mv` |
| **`ADR-0082` §Konsequenzen „test-only"** | **eingehalten** — dieselbe Messung (#12); „keine Zeile `internal/**` oder `cmd/**` außerhalb von Tests ändert sich, um Coverage zu gewinnen" ist wörtlich erfüllt |
| **Plan §1:** Anheben von `THRESHOLD` ist Wellen-Closure-Arbeit | **eingehalten** — `THRESHOLD ?= 70` unverändert, Datei nicht berührt (#13) |
| **Plan §1:** Cluster **D2** (`bootstrap`-Rest, `telemetry`) und Cluster **A** (`cmd`, `bootstrap`-`Run`) = andere Slices | **eingehalten** — **kein** Pfad aus `internal/bootstrap`, `internal/adapters/driven/telemetry` oder `cmd/pg-change-feed` im Vorgang; die Rest-Ungedeckten liegen unverändert dort (336 · 49 · 2, #9) |
| **Plan §1:** DB-Adapter-Coverage = anderer Messgegenstand | **eingehalten** — die drei Ausgenommenen bleiben aus der Messung (die Stufe nimmt sie über den `go list`-Filter aus, #5) |
| **Plan §3:** `welle-20.md` §4 = **„nicht"** (die **erreichte** Zahl gehört in die Closure-Notiz, nicht in die Welle) | **eingehalten** — der Vorgang ändert dort **eine** Zeile: die **Soll**-Zeile `application/usecase/* (10 Pakete)` → `(13 Pakete, 10 davon mit ungedeckten Statements)`. **Keine** erreichte Zahl in der Welle; der Grep nach `1493/1495/78,5/78,6/408/410` über alle stehenden Träger ist leer (#17) |
| **Plan §3:** `harness/sensors/coverage-gate.md` — „update, **nur falls** eine Zahl dort gegen die Messung driftet" | **nicht erfüllt** — kein Wert dieser Datei driftet gegen die Messung (nachgemessen, §5/Negativbefunde), die Datei wurde in `ad95444` dennoch angefasst → **V-1** |
| **Plan §3 (zweite Hälfte):** „**prüfend**, nicht nur nachziehend" — ob ein **anderer** Träger die bewegte Eigenschaft beschreibt (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`) | **erfüllt und eingetreten** — die Prüfung hat **einen** Fehler in **vier** Trägerstellen gefunden („zehn Pakete" als Fehlzählung des Globs) und berichtigt; die vier Zahlen sind jetzt je Aussage richtig (nachgezählt, §5) |
| **`AGENTS.md` §3.7** in den hinzugefügten Träger-Zeilen | **eingehalten** — kein Zustandsfeld, keine Chronik-Wendung; die zwei Vergangenheits-Formen in §1 („trugen ungedeckte Statements", „war bereits vollständig gedeckt") sind die **Unterscheidung** der vier Größen (Planungs-Stand vs. Vorgangs-Delta), kein Vorher/Nachher-Statusfeld |
| **`AGENTS.md` §3.11/§3.12** | **eingehalten** — `make docs-check` grün über die geänderten Träger (#1); die Zahl-Bewertung führt §5. Der **Rest**-Punkt der §3.12-Prüfung steht als V-1 (ein neu **mitgetragener** abgeleiteter Wert ohne Herkunfts-Marker) |

---

## 4. Plan-vs-Code-Diff

Verglichen gegen die §3-Liste des Plans am Stand `HEAD` (der Plan **ist** in
diesem Vorgang nachgezogen worden — `ad95444`; verglichen wird darum gegen die
geltende Fassung, nicht gegen einen Vorstand).

| §3-Zeile | Geliefert | Urteil |
|---|---|---|
| `internal/application/usecase/*/**_test.go` (**12** der **13** Pakete) — Test neu/update | **12** Dateien, genau die genannten Pakete (#12) | **Plan eingehalten** |
| `internal/domain/model/**_test.go` — Test neu/update | `validation_test.go` (#12) | **Plan eingehalten** |
| `harness/sensors/coverage-gate.md` — update, **nur falls** eine Zahl driftet | angefasst in `ad95444` (2 Zeilen, Deixis), **keine** Zahl geändert | **Abweichung von der Bedingung**, benannt → **V-1** |
| `docs/plan/planning/welle-20.md` §4 — **nicht** | eine **Soll**-Zeile berichtigt, **keine** erreichte Zahl | **Plan eingehalten** (die Zeile trägt Soll-Zahlen; die erreichte Zahl fehlt dort, #17) |
| — | das Review zu `slice-092` (neu) | **kein Plan-Bruch**: Übergabe-Artefakt der Reviewer-Rolle (Modul 8), kein Liefer-Punkt |

**Was der Diff nicht enthält, obwohl der Plan es nennt:** nichts. Die drei
Test-Zeilen sind vollständig geliefert; Produktcode, `THRESHOLD` und die
Cluster A/B/D2 sind pfadmäßig unberührt.

---

## 5. Die Zahlen — selbst gemessen (`AGENTS.md` §3.12), und was in §7 gehört

| Zahl | Wo sie heute steht | Mein Lauf | Befund |
|---|---|---|---|
| Planungsmaß **22** (Use-Cases) + **2** (`domain/model`) = **24** | Plan §1; Commit-Message `3de9547` | Parent-Profil: die **10** Pakete mit Rest summieren auf **22**, `domain/model` **2** (#10) | **hält, exakt** — die Plan-Zahl ist reproduzierbar |
| Zuwachs **+24** Statements | Commit-Message `3de9547` | Block-Positions-Differenz Parent → Diff: **24** Blöcke à **1** Statement, jeder benannt (#8) | **hält, exakt** (die Frage „sind es wirklich 24?" ist damit positiv beantwortet) |
| Netto **+22** (Schwankungsblock) | Commit-Message `3de9547` | **1** Block **neu ungedeckt**: `internal/bootstrap/wiring.go:991.5,992.13`, **2** Statements (#8) | **hält, exakt** — die Differenz ist vollständig erklärt |
| Nenner **1903** | Plan §1, `ADR-0082`, Sensor-Dokument | **1903** am Parent **und** am Diff (#5, #7) | **hält** — beide Stände |
| gedeckt/ungedeckt **1493/1495** · **408/410** | Commit-Message, Review | 6 eigene Läufe: gedeckt ∈ {**1493** (5 ×), **1495** (1 ×)}, ungedeckt ∈ {**410**, **408**} (#6) | **hält** — als **Band** zweier Enden desselben Stands |
| gedruckt **78,5 % / 78,6 %** | Commit-Message, Review | 5 × `78.5%`, 1 × `78.6%` (#6); die zwei `make gates`-Läufe drucken **78.60 %** (#1) und **78.50 %** (#19) — je ein anderes Schwankungsende | **hält** — Lauf-Größen; die Zuordnung Zahl → gedruckte Zeile stimmt: `1493/1903 = 78,455 % → 78.5 %`, `1495/1903 = 78,560 % → 78.6 %` (**abgeleitet**, nicht gedruckt) |
| **13** Pakete · **10** mit Rest · **12** mit Tests · **11** bewegte | Plan §1 (nach `ad95444`), `welle-20.md:100` | **13** (Glob) · **10** (Parent-Profil) · **12** (Diff-Dateiliste) · **11** (10 + `domain/model`; `readchanges`/`excludecolumn`/`includecolumn` bewegen sich nicht: `0 → 0`, #9/#10) | **hält, jede einzelne** — die vier Größen sind jetzt je Aussage richtig gesetzt (§1/§2/§5/§6/§8 widersprechen sich nicht mehr; F-1 aufgezählt und nachgezählt) |
| **`~100 Statements`** als Abstand der zwei Bänder | `harness/sensors/coverage-gate.md` §Zählbasis | Enden-Paarung: `1468 − 1369 = 99`, `1471 − 1371 = 100` (**abgeleitet**) | **trägt als Abstand der *gedeckten* Zahlen**, trägt **nicht** als Abstand der *Code-Stände*: der Nenner ist an beiden Bändern **1903** (#5, #7) → **V-1**, Rest-Punkt |
| „**21 Mutationen**, alle Exit 1" | Commit-Message `3de9547` | nicht nachfahrbar (§8); **6** eigene Proben + Kontrollen → kein Gegenbeispiel | **plausibel, nicht nachgemessen** — die Substanz trägt die Stichprobe |

**Was in §7 gehört** (LP1 verlangt es ohnehin, `AGENTS.md` §3.12 wendet es an):

1. **die Zustands-Zahl des Gegenstands** — **24** Statements in Cluster D1
   sind neu gedeckt (**22** Use-Cases + **2** `domain/model`), als **24**
   benannte Block-Positionen; **Nenner 1903 unverändert**; **Netto +22**, weil
   genau **ein** Block (`internal/bootstrap/wiring.go:991.5,992.13`, 2
   Statements) neu ungedeckt ist. Diese Zahlen hängen am Code-Stand und
   schwanken nicht.
2. **die Lauf-Größe mit ihrem Lauf** — die gedruckte Zeile des Gate-Laufs
   **78,60 %** (`make gates` am Stand `ad95444`, dieser Bericht) und, als
   **Band** formuliert, die sechs eigenen Profil-Läufe desselben Stands:
   gedeckt **1493–1495** (ungedeckt **410–408**), gedruckt **78,5 %/78,6 %** —
   die Schwankung ist die **eine** benannte Block-Position, nicht D1.
3. **die LP3-Aussage in ihrer gemessenen Form** — „im D1-Umfang bleibt **kein**
   ungedecktes Statement (13 Use-Case-Pakete + `domain/model`, je `0`,
   dedupliziert gemessen)" — und **nicht** als „im Gegenstand".
4. **keine erreichte Zahl in einen stehenden Träger** — `welle-20.md` §4 trägt
   weiter nur die **Soll**-Zahl (`24`); die erreichten Zahlen stehen in **keinem**
   `*.md` außerhalb von `docs/reviews/` (#17). Die §3-Zeile „`welle-20.md` = nicht"
   bleibt damit auch nach der Closure eingehalten.

---

## 6. Findings

### V-1 — Die §3-Bedingung für den Sensor-Träger ist **nicht** erfüllt: kein Wert driftet, die Datei wurde dennoch angefasst — und der Rest der Deixis ist eine abgeleitete Zahl ohne Herkunft

- `kategorie`: LOW
- `quelle`: Plan §3 Zeile 3 („update, **nur falls** eine Zahl dort gegen die
  Messung driftet") · `AGENTS.md` §3.12 Instanz A · Klasse
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (1×, Erstauftreten `slice-091`)
- `pfad`: `harness/sensors/coverage-gate.md:68-74` (§Zählbasis, der
  Acht-Läufe-Absatz) gegen `ad95444` (die zwei geänderten Zeilen) und
  `docs/plan/planning/in-progress/slice-092-coverage-cluster-d1.md:149` (§3-Zeile)
- `befund`: **Gemessen driftet keine Zahl der Datei.** Nenner **1903** heute
  bestätigt (#5, #7); `postgresstorage/mapper` **20/20** (§Grenze 1, mein Lauf:
  20 Statements, 0 ungedeckt, #9); „genau **vier** Pakete `Test=0 XTest=0`" →
  **4** (#11); „`TestGoFiles` allein trifft **23** der **31**" → **23/31**,
  davon **19** mit externem Testpaket (#11); `cmd` „**49** Statements, alle
  `count = 0`" → **0/49** (#9); `streamv1` „**76 von 86** im Gegenstand" → **76/86**
  (#9); die drei Ausgenommenen (**30/32**, **31/472**, **112/187**) sind
  ausdrücklich Lauf-Werte `slice-085`, und die **Rückrechnung** in §Grenze 4
  behält mit dem heutigen Zähler alle drei Ausgänge: `64,17 %` (rot), `76,79 %`
  und `78,71 %` (grün) — **abgeleitet** aus `1493` statt `1369`, keine eigene
  Messung. Die §3-Bedingung ist damit **im engen Sinn nicht erfüllt**; `ad95444`
  hat die zwei Zeilen dennoch geändert (nur die Wendung „**desselben, hier
  gegenständlichen** Produktionsstands" → „**desselben** Produktionsstands",
  `git show ad95444 -- harness/sensors/coverage-gate.md`). Die Änderung ist
  nachvollziehbar — der Review hatte mit **F-2** genau diese Stelle als
  alternde Deixis benannt —, sie ist aber **nirgends als Ausweitung der §3-Zeile
  benannt**: weder §3 noch §7 noch der Commit-Text sagen, dass hier eine
  *Prosa*-Alterung, kein *Zahl*-Drift, der Anlass war. **Zweiter Teil:** der
  Satz trägt die Unterscheidung der zwei Bänder jetzt an einer **abgeleiteten
  Zahl** (`~100 Statements auseinander`; gemessen `99`/`100` je Enden-Paarung,
  §5) **ohne Herkunfts-Marker** und hängt sie an die „**Code-Stände**", obwohl
  der Nenner an beiden Bändern **1903** ist — die ~100 ist der Abstand der
  **gedeckten** Zahlen. Das ist die Form, die §3.12 Instanz A für Differenzen
  ausdrücklich verlangt („als **abgeleitet** gekennzeichnet"), und der
  Commit-Text erklärt genau die Klauseln für alterungsfest, die es nicht ist.
- `verifizierbar`: ja — `git show ad95444 -- harness/sensors/coverage-gate.md`;
  Profil beider Stände und Block-Positions-Auswertung (#5, #7, #9, #11);
  `grep -n "~100 Statements" harness/sensors/coverage-gate.md`
- `urteil`: **zwei Optionen, eine ist zu wählen — still bleiben darf es nicht.**
  (a) Die Closure nimmt die Änderung als **innerhalb** der §3-Zeile an und
  schreibt eine Zeile dazu, *warum* hier Prosa statt Zahl der Anlass war; oder
  (b) sie liest die §3-Zeile wörtlich und trägt die **Abweichung** in §7 nach.
  Für (b) spricht ein Befund dieses Vorgangs selbst: die Bedingung „nur falls
  eine Zahl driftet" ist **zu eng** für einen Träger, dessen **Satz** altert,
  ohne dass eine Zahl driftet — genau die Aussage von
  `arbe-ueberholt-stehenden-traeger`, die dieser Slice zum zweiten Mal berührt.
  Der zweite Teil (das `~100` ohne Marker) ist ein **Ein-Wort-Nachtrag**
  („abgeleitet"), falls die Stelle ohnehin angefasst wird; er blockiert
  `done/` nicht.

### V-2 — Die Fixrunde `ad95444` ist von **keiner** Review gedeckt; das §2-Kriterium verlangt für sie ein Delta-Review

- `kategorie`: MEDIUM
- `quelle`: §2 des Slice-Plans, Zeile „Review durchgeführt … **Weist der Review
  eine Fixrunde aus, deckt ein Delta-Review sie ab** — die Lehre aus
  `slice-090` (V-2) und `slice-091` (N-1)" · der Review selbst: „**Merge-blockierend: ja**
  — 1 MEDIUM (F-1) … die Fixrunde deckt nach §2 des Plans ein **Delta-Review** ab"
- `pfad`: `docs/reviews/` (nur das Review zu `slice-092`) gegen `ad95444`
  (`docs/plan/planning/in-progress/slice-092-coverage-cluster-d1.md:49-58`
  §1 neu geschrieben, `:106` LP1, `:193` §5, `:217` §6, `:297` §8;
  `harness/sensors/coverage-gate.md` zwei Zeilen)
- `befund`: Der Review beurteilt ausdrücklich `3de9547` und weist mit F-1 eine
  Fixrunde aus („Rückgabe-Pfeil Reviewer → Implementer: **nötig**"). Danach hat
  `ad95444` **sechs** Planstellen und **zwei** Träger-Zeilen geändert — darunter
  eine **neue, vorher nicht dagewesene Prosa-Passage** (§1 „Drei Zahlen, drei
  Dinge", 8 Zeilen), nicht bloß ein Zahl-Austausch. Für diesen Stand existiert
  **kein** Review-Artefakt; die Präzedenz dieses Repos ist ein eigener
  Delta-Report (den Delta-Reviews zu `slice-089`, `slice-091`). Der
  Review hat die Fixrunde damit selbst als nötig bezeichnet, und die
  Hausform für „Fixrunde nach Review" ist die **Delta**-Prüfung — wer
  geschrieben hat, reviewt nicht (Modul 8).
- `verifizierbar`: ja — `ls docs/reviews/ | grep 092` (nur der Bericht zu
  `3de9547`); `git diff --name-status 3de9547^..ad95444` (#12);
  `git show ad95444 -- docs/plan/planning/in-progress/slice-092-coverage-cluster-d1.md`
- `urteil`: **Substanz geprüft, Artefakt fehlt.** Ich habe den Inhalt der
  Fixrunde unabhängig nachgemessen — die vier Zahlen halten je Aussage (13/10/12/11,
  §5), die Deixis ist entfernt (`grep "gegenständlichen"` → kein Treffer), kein
  Produktcode, kein Test berührt. **Die Severity ist trotzdem MEDIUM, weil
  das Kriterium im DoD steht und nicht erfüllt ist** — der Präzedenzfall
  `slice-091` (N-1/V-4) hatte für dieselbe Form MEDIUM. Der Unterschied zu dort
  ist real und mindert nur den Aufwand, nicht den Status: die Fixrunde hier ist
  **doku-only** und führt wörtlich aus, was der Review vorgeschrieben hat. Zwei
  Ausgänge sind zulässig: ein **kurzes Delta-Review** auf `ad95444` oder die
  **ausdrückliche §7-Zeile**, dass die Fixrunde die vorgeschriebene Korrektur
  ausführt, ihre Zahlen und Sätze in diesem Bericht nachgemessen sind und ein
  Delta deshalb entfällt. Still bleiben darf es nicht.

### V-3 — Die Commit-Message nennt „13 Use-Case-Pakete", der Diff ändert **12**

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 Instanz A (gilt für **Doku-Träger**; die
  Commit-Historie ist laut `ADR-0083` §Geltungsbereich ausdrücklich **kein**
  Träger) · §1 des Plans (die Unterscheidung 13/10/12/11)
- `pfad`: Commit `3de9547`, Betreff-Umfeld: „Tests fuer die **13** Use-Case-Pakete
  und domain/model. Kein Produkt-Code." gegen
  `git diff --name-only 3de9547^..3de9547` (**12** Use-Case-`*_test.go` +
  `domain/model/validation_test.go`) und §1 des Plans am Stand `HEAD`
  („**12** haben Tests bekommen (`readchanges` war bereits vollständig gedeckt)")
- `befund`: Der Diff ändert **12** der **13** Use-Case-Testdateien; für
  `readchanges` gibt es **keine** Teständerung (am Parent bereits `0/10`
  ungedeckt, #10). Die Zahl **13** ist in der Commit-Message als „die
  Use-Case-Pakete" (die Fläche) lesbar und dann nicht falsch — als Aussage über
  die **geänderten** Pakete ist sie es, und genau diese Aussage trägt der Plan
  an drei Stellen (LP1, §1, §3) mit **12**. Die Commit-Message ist der einzige
  Ort, an dem die Zahl des Zuwachses heute steht (§7 ist leer, #18) — wer sie
  für §7 abschreibt, kann die 13 in die falsche Aussage übernehmen.
- `verifizierbar`: ja — `git diff --name-only 3de9547^..3de9547 | grep -c "usecase.*_test.go"`;
  `git log -1 --format=%B 3de9547`
- `urteil`: **keine Reparatur (git-Historie wird nicht rückdatiert), aber ein
  Hinweis für §7:** die Zahl der **geänderten** Pakete ist **12**, die Zahl der
  Use-Case-Pakete im Glob **13**, die Zahl mit **bewegter** Deckung **11** — die
  drei gehören getrennt benannt, wie §1 des Plans es tut.

---

## 7. Negativbefunde

- **geprüft, ohne Befund: „test-only" — die §1-Zusage und die Folgepflicht aus
  `ADR-0082` §Konsequenzen.** **17** Pfade im Vorgang, davon **13** `*_test.go`
  und **4** Doku (drei Träger, ein Review-Report); **kein** `internal/**`/`cmd/**`
  außerhalb von Tests, kein `git mv`, keine ADR-Textänderung, kein `proto/`,
  keine `.pb.go`, kein `Makefile`-Ziel, keine Naht (#12, #13). `THRESHOLD`
  steht unverändert auf **70**; die Cluster A, B und D2 sind pfadmäßig unberührt
  — ihre Rest-Ungedeckten sind unverändert messbar (`bootstrap` **336**,
  `cmd` **49**, `telemetry` **2**; #9).
- **geprüft, ohne Befund: LP3 hält — Cluster D1 steht real bei `uncovered = 0`.**
  Paketweise am Diff-Stand, dedupliziert über die Block-Position: **alle 13**
  Use-Case-Pakete `0` ungedeckt, `domain/model` **0/110** (#9). Die Aussage „im
  D1-Umfang bleibt kein ungedecktes Statement" ist damit gemessen und **nicht**
  zu weit formuliert (die acht Pakete mit Rest liegen außerhalb D1).
- **geprüft, ohne Befund: die Plan-Zahlen des Vorgangs halten alle vier.**
  **13** (Glob) · **10** mit ungedeckten Statements am Parent (Summe **22**) ·
  **12** mit Tests · **11** mit bewegter Deckung (10 + `domain/model`;
  `readchanges`/`excludecolumn`/`includecolumn` bleiben bei `0` und bewegen sich
  nicht) — jede Zahl einzeln nachgezählt (#9, #10, #11, #12). **F-1 ist damit
  erledigt:** die vier Größen stehen je Aussage an genau einer Bedeutung, und
  §1/§2/§5/§6/§8 widersprechen sich nicht mehr.
- **geprüft, ohne Befund: der Zuwachs ist **+24**, als **24** Block-Positionen
  benannt.** Die Mengendifferenz der ungedeckten Block-Positionen Parent → Diff
  ist **exakt** die des Commits: **24** Blöcke à **1** Statement neu gedeckt
  (**22** in `usecase/*/service.go`, **2** in `domain/model`:
  `source.go:31.2,31.40`, `transaction.go:29.3,30.1`), **1** Block **neu
  ungedeckt** (`internal/bootstrap/wiring.go:991.5,992.13`, **2** Statements) →
  **Netto +22** (#8). Die Plan-Überlegung „die `22`/`2` sind eine
  Über-Schätzung" (§6 R1) ist damit **nicht eingetreten**; die Zahlen halten exakt.
- **geprüft, ohne Befund: LP2 — 6 eigene Produktions-Mutationen, 6 × rot, jede
  mit Kontrolle.** Die Proben und ihre Kontrollen:

  | Probe | Produktionsmutation (Datei:Zeile) | benannter Test | Probe | Kontrolle |
  |---|---|---|---|---|
  | P-1 | `acknowledge/service.go:61` Port-Fehler verworfen (`if false && err != nil`) | `TestAcknowledgeStateErrorFollowsConsumer` | **rot** | **Exit 0** |
  | P-2 | `disable/service.go:53` `TableExists` mit **fester Adresse** `public.t1` | `TestDisableTableExistsErrorFollowsAddress` | **rot** | **Exit 0** |
  | P-3 | `list/service.go:54` `Published` mit **fester Publication** `pub-2` | `TestListTablesPublishedErrorFollowsPublication` | **rot** | **Exit 0** |
  | P-4 | `retention/service.go:72` Freigabe unabhängig von der Policy (`… \|\| true`) | `TestRunDistinguishesEligibleChangesFromMixedSet` | **rot** | **Exit 0** |
  | P-5 | `domain/model/source.go:31` **konstanter Name** | `TestNewSourceCarriesIdentifierAndName` | **rot** | **Exit 0** |
  | P-6 | `excludecolumn/service.go:40` `ColumnExists` mit **fester Spalte** | `TestExcludeColumnRejectsMissingSourceColumn` | **rot** | **Exit 0** |
  | P-8 | `domain/model/timepoint.go:55` Dauer-Riegel aus (`if false`) | `TestNewRetentionPolicyRejectsNegativeDuration` | **rot** | **Exit 0** |
  | **P-7** | dieselbe Mutation wie P-6, aber auf der **Parent-Fassung** | `TestExcludeColumnRejectsMissingSourceColumn` (Parent) | **grün** | — |

  Jede Probe ist belegt als **angewandt** (`diff -q` gegen die unmutierte
  Datei: „Dateien … sind verschieden") und lief netzlos im gepinnten
  Toolchain-Container, `-count=1`, mit gelesener Testzeile (nicht bloß
  Exit 0 — die Kontrollläufe zeigen `ok`, nicht `[no tests to run]`; ein
  falscher `-run`-Name mit vacuous Exit 0 ist mir bei P-8 einmal unterlaufen
  und wurde verworfen und neu gefahren). **P-7 ist der Kern der
  LP2-Feststellung:** dieselbe Mutation an der **Parent**-Fassung bleibt
  **grün** — die zwei vorbestehenden Tests banden wirklich nicht, der Commit-Text
  ist nachgemessen (Review-F-3).
- **geprüft, ohne Befund: das Gate-Grün und die Zahl dazu sind reproduziert.**
  `make gates` **Exit 0**: `baseline-verify v6.5.0 OK — 54 Dateien`,
  `coverage-gate: OK — Coverage 78.60% erfüllt Schwelle 70%`, `d-check: 758
  Datei(en), 0 Befund(e)`, `commit-traceability: OK — 5 Commit(s)`,
  `generated-sync: OK` (beide `.pb.go` byte-gleich), `a-check: gesamt: 0
  Befund(e)`; danach Baum sauber (#1). Mein eigenes Profil derselben Fläche
  kommt auf **78,5–78,6 %** je Schwankungsende (#5, #6) — die gedruckte Zeile
  des Gates und die eigene Auswertung sind **dieselbe** Messung in anderer
  Ausgabepräzision: jede Zuordnung `1903 − ungedeckt` trifft die gedruckte
  Prozentzeile (#6).
- **geprüft, ohne Befund: die zwei Doc-Gates über den Vorgang.**
  `make doc-commits RANGE=3de9547^..ad95444` **Exit 0**, 758 Dateien,
  0 Befunde; `make doc-immutable RANGE=3de9547^..ad95444` **Exit 0**, 758
  Dateien, 0 Befunde (#3, #4).
- **geprüft, ohne Befund: `BEO-PGC/arbeit-ueberholt-stehenden-traeger` — die
  Prüfung des Commits hält, und kein *stehender* Satz wird durch diesen Slice
  falsch.** Die mechanische Gruppierung ist vor und nach dem Slice
  **byte-identisch**: `go list -f '{{.ImportPath}} Test={{len .TestGoFiles}}
  XTest={{len .XTestGoFiles}}'` über den Gegenstand liefert am Parent **und** am
  Diff dieselbe Ausgabe (md5 `adda7355479fc95c8414a9e9ff3947bf`), **31** Pakete,
  **4 × `Test=0 XTest=0`**, **23 × `Test=0`**, `streamv1` unverändert `0/1`
  (#11). Und ich habe **jede** Zahl des Trägers gegen meinen eigenen Lauf
  gehalten (Nenner, Gruppierung, `cmd` 49, `streamv1` 76/86, `mapper` 20/20,
  die drei Ausgenommenen, die Rückrechnung in §Grenze 4) — **keine** wird durch
  diesen Vorgang falsch; die zwei Vergangenheits-Bänder sind Lauf-Belege mit
  Lauf-Marker. Der **eine** Rest dieser Stelle steht als V-1.
- **geprüft, ohne Befund: keine erreichte Zahl in einem stehenden Träger.**
  `grep` über alle `*.md` außerhalb `docs/reviews/` und der vendored Baseline
  nach `1493|1495|78,5|78.5|78,6|78.6|408|410` → **kein** Treffer in einem
  Doku-Träger (der Roh-Treffer ist ein Actions-SHA in `AGENTS.md`); `welle-20.md:100`
  trägt die **Soll**-Zeile `(13 Pakete, 10 davon mit ungedeckten Statements) | 24`
  (#17). Die §3-Zeile „`welle-20.md` §4 = **nicht**" ist damit eingehalten, und
  die Ablage-Pflicht liegt unverändert bei §7.
- **geprüft, ohne Befund: die vier §6-Risiken — Material für drei Ausgänge.**
  Gemessen: **R1** (Über-Schätzung) **nicht eingetreten** (22+2 exakt, #8/#10);
  **R2** (Negativtest am Fake) an den **neuen** Tests **nicht eingetreten**
  (6 Proben rot) und an **zwei vorbestehenden** Stellen eingetreten und
  **behoben** (P-7 grün am Parent); **R3** (Coverage-Theater) mit Gegenbeleg —
  jede Probe mutiert die **Produktionsseite** und färbt ihren benannten Test
  rot; **R4** (überholter Träger) **eingetreten** und benannt (F-1: vier
  Trägerstellen; F-2: eine Deixis). Die **Ausgänge** selbst stehen weiter auf
  `<…>` (#18) — Planner-Closure-Arbeit.
- **geprüft, ohne Befund: der Umfang trägt keine zweite Quelle für den
  Zähler.** Im Vorgang ist **keine** Registerdatei geändert (#15);
  `arbeit-ueberholt-stehenden-traeger` steht bei **1** Beleg
  (`evidence/slice-091.md`), `negativtest-ohne-bindung-an-seine-eingabe` bei
  **5** (`slice-083/086/087/088/091`) — die Zähler folgen den Dateien, und die
  Dateiung der zwei Vorkommen dieses Vorgangs ist eine Closure-Entscheidung,
  keine Zahl, die dieser Bericht setzt.

---

## 8. Was ich nicht prüfen konnte

- **Die Läufe des Implementers und des Reviewers.** Die Commit-Message nennt
  „**21 Mutationen**, alle Exit 1", der Review „**32** eigene Proben, 32 × rot".
  Beider Läufe sind nicht mehr einsehbar. Ich habe die **Substanz** unabhängig
  gemessen (**6** eigene Produktions-Mutationen an neuen Tests, alle rot, plus
  **P-7** an der Parent-Fassung) — ein **Gegenbeispiel** hat sich in dieser
  Stichprobe nicht gezeigt; die **konkreten** Läufe sind als Belege nicht
  überprüfbar, und die Bindung der **übrigen** neuen Tests ruht auf den zwei
  Berichten.
- **Die dritte Nachkommastelle der gedruckten Zeile.** Die Zuordnung
  `1493 → 78.5 %` und `1495 → 78.6 %` ist **abgeleitet** (Division), nicht
  gedruckt; gedruckt gesehen habe ich in meinen Läufen genau die zwei Zeilen
  (§5, #6) — die zwei Gate-Läufe dieses Berichts drucken **78.60 %** (#1) und
  **78.50 %** (#19), also dasselbe Band. Ob dazwischen ein drittes Ende
  existiert, ist mit 8 Läufen (6 eigene + 2 Gate-Läufe) nicht ausgeschlossen — der Kontext (2) dieser ADR kennt für den
  Vorgängerstand ebenfalls genau zwei Enden.
- **Der reale Post-Push-Lauf.** Dieser Vorgang ändert **keinen** Workflow —
  `AGENTS.md` §3.10 greift dem Buchstaben nach nicht. Die Wirkung der neuen
  Tests auf die CI-Laufzeit bleibt bis zum nächsten echten Lauf unbelegt.
- **Die Stabilität über die Rampe.** Ob die Zahl bei `THRESHOLD=80` (der
  `welle-20`-Closure) stabil genug ist, ist mit diesem Stand nicht entscheidbar
  — dafür braucht es die Messung **am** Hochschalt-Punkt, nicht die
  Extrapolation.
- **Der Zähler-Stand nach diesem Vorgang.** Ob `arbeit-ueberholt-stehenden-traeger`
  auf **2×** und `negativtest-ohne-bindung-an-seine-eingabe` auf **6×** geht,
  ist eine Register-Entscheidung (Modul 6: „Mensch urteilt, Maschine prüft
  Deckung") — dieser Bericht liefert nur die Messung, die ihr zugrunde liegt.

---

## 9. Verdikt

**Die drei Liefer-Punkte aus §2 tragen — gemessen, nicht gelesen.** **LP1:**
**13** Testdateien sind real geliefert (**12** Use-Case-Pakete +
`domain/model`), sie laufen im Gate (`coverage`-Stufe über die `go list`-Liste,
netzlos, `TESTEXIT=0`) **und** in `make test` (netzlos, Race-Detector, **Exit 0,
32 × `ok`**); `make gates` ist **Exit 0**; der Zuwachs ist **+24** Statements als
**24** benannte Block-Positionen (22 + 2), **Netto +22**, Nenner **1903**
unverändert — die Zahl mit ihrem Lauf steht in §5. **LP2:** **6 eigene
Produktions-Mutationen** an neuen Tests färben **6 × rot**, jede mit grüner
Kontrolle; **P-7** beweist an der **Parent**-Fassung, dass die zwei reparierten
Spaltenbindungen vorher wirklich nicht banden. **LP3:** **alle 13** Use-Case-Pakete
und `domain/model` stehen bei `uncovered = 0` — im D1-Umfang gibt es **nichts**
zu benennen, und die Aussage ist nicht zu weit formuliert.

**Entscheidungs-Konformität: hält — mit einer benannten Ausnahme.** Kein
Produktcode (**17** Pfade, **13** Testdateien), `THRESHOLD` unverändert **70**,
Cluster A, B und D2 pfadmäßig unberührt, die Plan-Zeile `welle-20.md` = „nicht"
eingehalten (nur eine **Soll**-Zeile berichtigt, keine erreichte Zahl in einem
stehenden Träger), `ADR-0082` „test-only" erfüllt. **Die Ausnahme ist V-1:** die
§3-Bedingung für `harness/sensors/coverage-gate.md` („nur falls eine Zahl dort
gegen die Messung driftet") ist **nicht erfüllt** — kein Wert driftet, die Datei
wurde dennoch angefasst (Deixis-Korrektur nach Review-F-2), und der Rest der
Stelle ist eine abgeleitete `~100` ohne Herkunfts-Marker.

**Der Plan-vs-Code-Diff ergibt keine unbegründete Abweichung:** die drei
Test-Zeilen sind vollständig geliefert, keine Zeile fehlt; die vierte Zeile
(`coverage-gate.md`) ist die benannte Ausnahme (V-1).

**`done/`-fähig ist der Slice damit noch nicht — vier Punkte fehlen:**

1. **V-2** — die Fixrunde `ad95444` hat **kein** Review-Artefakt; §2 verlangt
   für sie ein **Delta-Review** (Präzedenz `slice-089`/`slice-091`), und **das
   ist der einzige Punkt dieser Liste, der ein DoD-Kriterium direkt verletzt.**
2. **V-1** — eine Entscheidung zur §3-Zeile: annehmen mit einer Zeile, warum
   hier Prosa der Anlass war, **oder** die Abweichung in §7 nachtragen; dazu
   ggf. das eine Wort „abgeleitet" an der `~100`.
3. **§7** — die Closure-Notiz mit Lerneintrag, darin **die Zahl mit ihrem Lauf**
   (§5: Zustands-Zahl und Lauf-Größe getrennt) und die LP3-Aussage in ihrer
   gemessenen Form.
4. **Die Häkchen, die vier §6-Ausgänge und die Register-Entscheidung** —
   Material für drei der vier Ausgänge liefert dieser Bericht (§7/Negativbefunde);
   die Paarungen gehören der `welle-20`-Closure.

Nach 1–4 ist der Slice `done/`-fähig. **Kein Liefer-Defekt, kein rotes Gate:**
die offenen Punkte sind ein fehlendes Übergabe-Artefakt (V-2), eine
Entscheidung über die Reichweite einer Plan-Zeile (V-1) und die regulären
Closure-Pflichten — keine gebrochene Zusage.

---

**Beleg-Lage dieses Berichts:** Jede Zahl stammt aus einem der Läufe in §1, je
in eigener Werkzeug-Beauftragung gefahren; der Gate-Lauf (#1) und seine
Auswertung waren **zwei** Schritte, sein Exit-Code wurde aus einer separaten
Datei gelesen, nie durch eine Pipe (§3.9). Die Mutationsproben liefen auf
**Arbeitsbaum-Kopien außerhalb des Repos** (Parent-Stand über `git archive`);
Logs und Profile liegen **außerhalb** des Baums. Der Baum ist nach diesem
Bericht sauber, **kein** Commit, keine Änderung an Artefakten des Slice,
`THRESHOLD`, Produktcode oder an einem Träger außerhalb dieses Berichts.
