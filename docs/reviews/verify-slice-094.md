# Verifikationsbericht: slice-094 — 2026-09-16

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den DoD-Vertrag
(`slice-094` §2, LP1–LP3), die im Slice referzierten Entscheidungen
[ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
(Schnittmaß, Prozess-Rand als Puffer, Re-Exec-Harness „ohne Produktionsänderung",
§Konsequenzen *test-only*) ·
[ADR-0071](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(Messgegenstand) ·
[ADR-0077](../plan/adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md)
(die Rampe) · [ADR-0047](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md)
(Rollen-DSN) ·
[ADR-0083](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
(Herkunft von Aussagen in Trägern) sowie die Hard Rules `AGENTS.md` §3.2, §3.5,
§3.6, §3.7, §3.9, §3.11, §3.12. **Nicht** gegen den Diff als solchen
(Reviewer-Aufgabe) und **nicht** gegen realen Bedarf (Validator, nicht
ausgelöst).

**Frischer Kontext.** Der Slice-Plan wurde am Stand `HEAD` vollständig gelesen
(§1–§8), dazu [ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md),
die zwei vorliegenden Review-Reports, der Sensor
[`coverage-gate.md`](../../harness/sensors/coverage-gate.md), die zwei neuen
Testdateien, `cmd/pg-change-feed/main.go`, `internal/bootstrap/wiring.go` und der
berührte Register-Eintrag `BEO-PGC/arbeit-ueberholt-stehenden-traeger`.
Implementer-Bericht und Reviews waren **Kontext**, ihre Zahlen **nicht**
übernommen: jede Zahl dieses Berichts stammt aus einem hier selbst gefahrenen
Lauf — einschließlich der Review-Findings, die ich nachgemessen habe (F-2 an
**acht** eigenen Mutationsläufen über **zwei** Teststände, die drei reparierten
Sätze an `go list` und am Profil, §Grenze 7 an **fünf** eigenen Läufen).

**Gegenstand.** `HEAD` = `d839975`, Zweig `main`, Baum sauber. Der Vorgang sind
**sieben** Commits über `4ba09ff`: `32b8b9d` (Arbeit, Cluster A), `9e9df55`
(Review), `3b7f8c7` (Planner-Nachzug `welle-20` §1), `8292766` (Fixrunde auf
F-1/F-2/F-3/F-4/F-6), `f64794b` (Planner-Nachzug `welle-20` §1), `3bfaddb`
(Delta-Review), `d839975` (D-1, D-2, F-7, zwei Herkünfte). Der Slice liegt in
`in-progress/`; der `git mv` nach `done/` ist **nicht** erfolgt.

**Zwei Zahlen zur Schreibweise.** (a) Das Gate-Skript endet im Rot-Fall mit
**Exit 1**; `make` kapselt den Rezept-Fehlschlag zu **Exit 2** — beide stehen
unten getrennt, wo sie anfallen. (b) Prozentwerte sind die **gedruckten** Zeilen
eines konkreten Laufs; wo ich einen Wert **zurückrechne**, steht das dabei
(`AGENTS.md` §3.12 Instanz A).

**Beleg-Lage.** Jeder Exit-Code ist **ungepiped** ermittelt und in einem
**eigenen**, abgeschlossenen Schritt aus einer separaten Datei gelesen
(`AGENTS.md` §3.9); Gate-Lauf und Auswertung waren getrennt beauftragt. Die
Mutationsproben liefen auf Arbeitsbaum-Kopien **außerhalb** des Baums
(`git archive`), netzlos (`--network none`); die Repo-Dateien wurden **nicht**
angefasst (nach jeder Probe gegen das Original geprüft). Der Baum ist vor und
nach jedem Lauf sauber (`git status --porcelain` leer).

---

## 1. Eigene Messungen dieses Laufs

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `make gates` (Log in Datei, Exit **danach** aus eigener Datei) | **0** | `baseline-verify: v6.5.0 OK — 54 Dateien` · `coverage-gate: OK — Coverage 83.10% erfüllt Schwelle 70%` (gedruckte `total:`-Zeile derselben Stufe: `83.1%`) · `d-check: 775 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` · `generated-sync: OK` (beide `.pb.go` byte-gleich) · `a-check: gesamt: 0 Befund(e)`; Baum danach leer |
| 2 | `make coverage-gate THRESHOLD=80` (Rampen-Beleg, obere Stufe) | **0** | gedruckt `total: (statements) 83.0%`; `coverage-gate: OK — Coverage 83.00% erfüllt Schwelle 80%` |
| 3 | `make coverage-gate THRESHOLD=85` (Rampen-Beleg, Rot-Hälfte) | **2** (make) | Stufe bricht: gedruckt `83.1%`; `coverage-gate: FAIL — Coverage 83.10% unter Schwelle 85%`; Skript-EC **1** (`did not complete successfully: exit code: 1`); `make: *** [harness/mk/coverage.mk:24: coverage-gate] Fehler 1` |
| 4 | Eigenes Profil am Diff-Stand — Nachbau der Stufe: `go list ./internal/... ./cmd/...` ohne die drei ausgenommenen Pakete, `-coverpkg`, `-covermode=atomic`, `--network none`, **dedupliziert über die Block-Position**, **7 ×** | 7 × **0** | gedruckte Zeile 7 × `83.1%`; gedeckt **1581** (7 ×), ungedeckt **322**, Nenner **1903** (7 ×); `wiring.go` **222/521** (7 ×) |
| 5 | dasselbe am **Parent** `4ba09ff` (Arbeitsbaum-Kopie außerhalb des Baums, `git archive`) | **0** | gedeckt **1523**, ungedeckt **380**, Nenner **1903**, gedruckt `80.0%`; `cmd/pg-change-feed` **0/49**; `wiring.go` **213/521** |
| 6 | paketweise, dedupliziert (Diff-Stand) | 0 | `cmd/pg-change-feed` **49/49** · `internal/bootstrap` **299/598**; die acht Dateien mit ungedeckten Statements summieren zu **322** (= 1903 − 1581) |
| 7 | LP3 — Zuordnung **jeder** ungedeckten Block-Position der `wiring.go` zu einer Funktion (Funktionsgrenzen aus `grep -n '^func '`) | 0 | **134** ungedeckte Blöcke, **sämtlich** in `wiring.go`: `Run` **189/198** · `Diagnose` **73/81** · `Healthcheck` **14/22** · `RegisterConsumer` **13/17** · `AcknowledgeConsumer` **10/14** → **299**; **außerhalb der fünf: 0 Blöcke, 0 Statements** |
| 8 | die vier `cmd`-Aufruf-Blöcke, blockgenau | 0 | `main.go:44.3,44.82` · `:62.3,62.77` · `:86.3,86.88` · `:101.3,101.79` — je **max-count = 1**; die vier Zeilen rufen `Healthcheck`/`RegisterConsumer`/`AcknowledgeConsumer`/`Diagnose` (Zeilennamen gelesen) |
| 9 | `go list -f '{{.ImportPath}} Test={{len .TestGoFiles}} XTest={{len .XTestGoFiles}}'` über den Gegenstand (Diff-Stand) | 0 | **31** Pakete; **genau drei** `Test=0 XTest=0` — `postgresstorage/queries`, `application/port/inbound`, `domain/errors`; `TestGoFiles=0` trifft **22**, davon **19** mit `XTest>0`; `cmd/pg-change-feed` **Test=1 XTest=0** |
| 10 | dasselbe am **Parent** `4ba09ff` | 0 | **vier** Pakete `Test=0 XTest=0`, darunter `cmd/pg-change-feed` |
| 11 | die drei Pakete **ohne ausführbare Statements** im Profil | 0 | **0** Profilzeilen je Paket — sie kommen im Profil nicht vor |
| 12 | `-v`-Lauf der zwei Pakete, netzlos | **0** | **6** neue Testfunktionen, jede **1 ×** `=== RUN`; 145 × `=== RUN`, 67 × `--- PASS`, 21 × `--- SKIP`, **0 × `--- FAIL`** |
| 13 | `go test -count=20` der zwei Pakete, netzlos | **0** | 20 × grün, **kein** Flake; **0** Treffer für `time.`/`Sleep`/`After(`/`Timeout`/`Deadline` in den zwei neuen Dateien |
| 14 | `make test` (netzlos, `-race`) | **0** | **33 × `ok`, 0 × `FAIL`, 0 × `DATA RACE`**; `cmd/pg-change-feed` und `internal/bootstrap` je `ok` |
| 15 | **8 Mutationsläufe** am Diff-Stand (Arbeitsbaum-Kopie, netzlos, je eigener Aufruf) + Kontrolle | Kontrolle **0**, Proben je **1** | Kontrolle grün; alle acht rot, je mit benanntem Subtest/Test (§2/LP2) |
| 16 | **F-2-Probe in beiden Ständen**: moduseigene Meldung `main.go:73` durch den **Fallback-Text** von `main.go:104` ersetzt (Ausgang bleibt 2) | Testfile-Stand `32b8b9d` **0**, Diff-Stand **1** | **vorher grün, jetzt rot** — selbst gefahren, nicht übernommen; am Diff-Stand genau die zwei `acknowledge-consumer`-Subtests |
| 17 | §Grenze-7-Probe: `kindUmgebung` ohne `GOCOVERDIR`-Weitergabe (Arbeitsbaum-Kopie), **5 ×** | 5 × **0** | **alle** Tests grün (28 × `ok`, 0 × `FAIL`); `cmd` **0/49**; gedeckt **1532** (4 ×) / **1531** (1 ×); gedruckt 5 × `80.5%` |
| 18 | `make doc-commits RANGE=4ba09ff..d839975` | **0** | 775 Dateien, **0** Befunde — jeder der sieben Commits trägt seine Kennung |
| 19 | `make doc-immutable RANGE=4ba09ff..d839975` | **0** | 775 Dateien, **0** Befunde |
| 20 | Umfang: `git diff --name-status 4ba09ff..d839975`; `-- 'internal/**' 'cmd/**'` ohne `_test.go`; `-- harness/mk Makefile Dockerfile`; `-- tools/`; `-- docs/plan/planning/observations/` | 0 | **7** Pfade (2 × `A` `*_test.go`, 2 × Review-Report, Sensor-Doku, Slice-Plan, `welle-20.md`); die **vier** übrigen Abfragen sind **leer** |
| 21 | Zustands-Prüfungen | 0 | `harness/mk/coverage.mk:15` = `THRESHOLD ?= 70`; §6: **4 von 4** Risiken auf `<…>`; §7: **7 von 7** Inhaltszeilen Platzhalter; `docs/plan/planning/reconciliation.md` existiert **nicht**; `docs/reviews/` trägt **zwei** 094-Reports (kein dritter) |
| 22 | Doppel-Überschriften im Slice-Plan; wer den Anker verlinkt | 0 | `## 6. Risiken und offene Punkte## 6. Risiken und offene Punkte` und `## 7. Closure-Notiz## 7. Closure-Notiz` — **kein** Dokument verlinkt den Anker (`grep -rn 'slice-094-coverage-cluster-a.md#'` → leer), `git show e160e26:` zeigt beide schon bei der Anlage |
| 23 | Träger-grep nach der bewegten Eigenschaft (Auftrag aus `BEO-PGC/arbeit-ueberholt-stehenden-traeger`) | 0 | `cmd/pg-change-feed` + `49` / Paket-gruppen-Sätze / `TestGoFiles` / `coverage: 0.0%` über `--include='*.md'` ohne `docs/reviews/`: **kein** weiterer Treffer außerhalb des Diffs |
| 24 | `make gates` **mit diesem Bericht im Baum**, erneut (Log in Datei, Exit **danach** aus eigener Datei) | **0** | sechs Checks: `baseline-verify v6.5.0 OK — 54 Dateien` · `coverage-gate: OK — Coverage 83.10% erfüllt Schwelle 70%` · `d-check: 776 Datei(en) geprüft, 0 Befund(e)` (eine Datei mehr = dieser Bericht) · `commit-traceability: OK — 5 Commit(s)` · `generated-sync: OK` · `a-check: gesamt: 0 Befund(e)`; der Baum trägt danach allein diesen Bericht |

---

## 2. DoD-Konformität, Kriterium für Kriterium

### Liefer-Punkt 1 — die Tests existieren und sind netzlos grün

| Kriterium (§2) | Befund |
|---|---|
| die Tests existieren | **erfüllt** — zwei neue Dateien (`cmd/pg-change-feed/main_test.go`, `internal/bootstrap/run_test.go`) mit **6** Testfunktionen; die Re-Exec-Rolle trägt `TestMain` über den Marker `PGC_FEED_MAIN_REEXEC` (#12) |
| der Gate-Lauf **fährt sie wirklich** | **erfüllt** — die `coverage`-Stufe bildet ihre Paketliste aus `go list ./internal/... ./cmd/...` (ohne die drei ausgenommenen) und bricht bei rotem Testlauf ab; `cmd/pg-change-feed` steht in dieser Liste und sein Testlauf endet `ok` (#4, #12); **49 von 49** Statements des Pakets sind nur erreichbar, weil die Zähler der Kindprozesse über `GOCOVERDIR` ins Profil mergen (#4, #17) |
| `make gates` ist grün | **erfüllt** — **Exit 0** aus separater Datei gelesen, sechs Checks grün, Baum danach leer; erneut grün **mit** diesem Bericht im Baum (#1, #24) |
| **der Zuwachs wird als Zahl mit ihrem Lauf genannt** | **materiell erfüllt, Träger offen** — gemessen: `cmd/pg-change-feed` **0 → 49** gedeckt, `wiring.go` **213 → 222**, `internal/bootstrap` ungedeckt **336 → 299**, Gesamt gedeckt **1523 → 1581** bei **unverändertem Nenner 1903** (#4, #5, #6). Die Zahl steht heute nur in **Commit-Messages** (git-Historie — laut [ADR-0083](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) §Geltungsbereich **kein** Doku-Träger), in den zwei Review-Reports (Lauf-Belege, über Läufe hinweg nicht gelesen) und im Sensor-Dokument als `49 von 49` **mit** Lauf-Anker; der **Zuwachs** selbst und sein **Band** sind in **keinem** Träger. §7 ist leer (#21). Material für §7 steht in **§5** dieses Berichts |
| **die erreichte Quote mit ihrem Band** | **materiell erfüllt, Träger offen** — die gedruckte Zeile lautet `83.1%` in **zehn** Läufen (#1, #3, 7 × #4, #24) und `83.0%` in **einem** (#2); das Band desselben Stands ist **1** Statement breit: gedeckt **1580–1581** von 1903. Die **Größe** des Bands trägt das Sensor-Dokument §Zählbasis („höchstens **1** Statement", dort als **abgeleitet** gekennzeichnet) — die **Lage** dieses Laufs trägt kein Träger (#21) |
| **die zwei Rampen-Belege** | **erfüllt** — beide Hälften selbst gefahren, je mit eigenem Exit: `THRESHOLD=80` → **EC 0**, `OK — Coverage 83.00% erfüllt Schwelle 80%` (#2); `THRESHOLD=85` → **FAIL**, Skript-EC **1**, make-EC **2** (#3) |

### Liefer-Punkt 2 — die Prozess-Verträge am Exit-Code, nicht am Text

| Kriterium (§2) | Befund |
|---|---|
| der Prozess endet mit dem **vereinbarten** Exit-Code | **erfüllt** — die Tests binden **1** (Sondermodi am fehlenden Dienst), **2** (Argument-Fehler und argumentloser Lauf ohne Vorbedingung) und **0** (`--version`); die Verträge selbst stehen im Benutzerhandbuch (Rollen-Verdrahtung der vier Sondermodi) und in `main.go`. **Eigene Exit-Code-Mutation:** `os.Exit(2)` → `os.Exit(3)` im Zweig „unbekanntes Argument" färbt `TestArgumentFehlerEndenMitAusgang2/unbekanntes_Argument` **rot** (#15, M2) |
| die Ausgabe ist an ihre **Eingabeseite** gebunden | **erfüllt** — **acht** eigene Mutationen, **acht** rot: `--healthcheck` → `cfg.AdminDSN` (#15, M1) · `diagnose` → `cfg.AdminDSN` (M3) · `RegisterConsumer` → `cfg.ReaderDSN` (M4) · `ConfigFromEnv` nennt für den fehlenden `CDC_READER_DSN` den Namen `CDC_ADMIN_DSN` (M7) · dazu M2, M5, M6 und die F-2-Probe (#16) |
| **nicht** „irgendein Text irgendwo" | **erfüllt** — die entscheidende Probe ist die F-2-Mutation: mit dem **Fallback-Text** von `main.go:104` an der Stelle der moduseigenen Meldung bleibt die Suite am Teststand `32b8b9d` **grün (EC 0)** und wird am Diff-Stand **rot (EC 1)** (#16). Der Übergang ist damit **von mir** gemessen, nicht aus dem Delta-Review übernommen; am Diff-Stand fallen genau die zwei `acknowledge-consumer`-Subtests |
| der Auslöser ist bedient | **erfüllt in der Sache** — `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (**4×**) beschreibt den Träger, der seinen Satz nicht trägt; der Sondermodus-Test prüft jetzt die **konkrete Verstoß-Formulierung** je Fall, und die generische Zeile trägt sie nicht mehr (#15, #16). **Kein neuer Register-Beleg** — Entscheidung des Lese-Schritts (§2-Zeile „Beobachtungs-Register fortgeschrieben" ist offen, #21) |

### Liefer-Punkt 3 — die unerreichbaren Statements sind benannt

| Kriterium (§2) | Befund |
|---|---|
| die **189** hinter dem `Ping`-Riegel, namentlich als Block | **erfüllt** — gemessen (#7): `Run` **189/198**, `Diagnose` **73/81**, `Healthcheck` **14/22**, `RegisterConsumer` **13/17**, `AcknowledgeConsumer` **10/14** = **299**; der Grund ist der gemessene `Ping`-Riegel (`pgxpool.New` + `pool.Ping`, `receive.NewStream`, `nats.Connect`). Die fünf Funktionen mit Grund und Zahl führt [ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) §Kontext (4a); §1 des Plans nennt die **189** |
| **jede weitere** unerreichbare Stelle einzeln | **erfüllt** — im `bootstrap`-Rest liegt **keine**: **0** ungedeckte Blöcke außerhalb der fünf Funktionsgrenzen (#7). Für `cmd/pg-change-feed` ist die Aussage gegenstandslos geworden: **49 von 49** gedeckt, die vier früher offenen Aufruf-Blöcke sind es jetzt (#8) — das Sensor-Dokument nennt sie namentlich mit Zeile und Funktion, und es sagt den Grund für den Rest: dienstgebunden ist der **Rumpf**, nicht der Aufruf |
| die vier `cmd`-Aufrufe, die jetzt gedeckt sind | **erfüllt** — `main.go:44`/`:62`/`:86`/`:101` tragen `count > 0`; die Zuordnung Zeile → Funktion ist gelesen und stimmt (#8) |

### Die Closure-Pflichten aus §2

| Kriterium (§2) | Befund |
|---|---|
| `make gates` grün | **erfüllt** (#1, #24) |
| Review durchgeführt, Report unter `docs/reviews/`, **Delta-Review deckt eine Fixrunde ab** | **nicht erfüllt** — das Review zu `slice-094` (`32b8b9d`) und der Delta-Review zu `slice-094` (`8292766`) liegen vor; der Delta-Review weist mit D-1 und D-2 und dem **offenen F-7** eine weitere Runde aus („Fixrunde: **ja**"), `d839975` führt sie aus — für **diesen** Stand existiert **kein** Review-Artefakt → **V-1** |
| Verifikation durchgeführt, Report unter `docs/reviews/verify-slice-094.md` | **erfüllt** mit diesem Bericht; das Häkchen ist offen (#21) |
| Closure-Notiz mit Steering-Loop-Lerneintrag | **nicht erfüllt** — §7 trägt in **7 von 7** Inhaltszeilen Platzhalter (#21); Material für die Zahlen liefert §5, für den Register-Teil §6 |
| Reconciliation-Register fortgeschrieben *(entfällt …)* | **entfällt nachweislich** — `docs/plan/planning/reconciliation.md` existiert nicht (#21) |
| Beobachtungs-Register fortgeschrieben — **kein Zähler wird gesetzt** | **nicht erfüllt** — im Vorgang ist **keine** Registerdatei geändert (#20). Die **Substanz** liefert §6: der dritte Vorgang von `arbeit-ueberholt-stehenden-traeger` (die drei reparierten Sätze) und der Auslöser von `beleg-befehl-traegt-seinen-satz-nicht` |
| Jedes Risiko aus §6 trägt einen Ausgang | **nicht erfüllt** — **4 von 4** stehen auf `<…>` (#21). Material für alle vier liefert §5/§8 dieses Berichts (R1 *entfallen* · R2 *entfallen* · R3 *entfallen, mit Gegenbeleg* · R4 *eingetreten und bedient*) |
| Die drei Paarungen sind getragen | **nicht erfüllt** — fällt laut Zeile der `welle-20`-Closure zu |

**Ergebnis:** Die **drei Liefer-Punkte** und die Gate-Zeile **tragen** — gemessen,
nicht gelesen. Offen sind die regulären Closure-Pflichten (Häkchen, §7 samt
Zahlen und Band, Register-Entscheidung, vier §6-Ausgänge, Paarungen) und **ein**
Punkt, den nur diese Rolle sieht: **V-1** (die Fixrunde ohne Review-Artefakt).

---

## 3. Entscheidungs-Konformität — hält der Vorgang, was Plan §1/§3 und die ADRs zusagen?

| Zusage | Befund |
|---|---|
| **Plan §1:** „Was dieser Slice liefert: Tests. Er ändert **keinen** Produkt-Code" | **eingehalten** — **7** Pfade, **2** davon `*_test.go`; die Abfrage `git diff … -- 'internal/**' 'cmd/**'` ohne `_test.go` ist **leer** (#20); `cmd/pg-change-feed/main.go` und `internal/bootstrap/wiring.go` sind **nicht** im Diff |
| **Plan §1:** der Weg ist der Re-Exec-Harness und er ändert keinen Produktionscode | **eingehalten** — die **49/49** sind am **unveränderten** `main.go` gemessen (#4 auf dem Baum, #15 auf einer Kopie desselben Stands); kein Umbau, keine Naht, kein neues Interface |
| **Plan §1:** Anheben von `THRESHOLD` ist Wellen-Closure-Arbeit | **eingehalten** — `THRESHOLD ?= 70` unverändert (#21); `harness/mk Makefile Dockerfile` im Vorgang **nicht** berührt (#20) |
| **Plan §1:** kein eigenes Gate für die Prozess-Verträge | **eingehalten** — kein neues Make-Target; der Harness läuft in `make test` und in der `coverage`-Stufe (#1, #14) |
| **Plan §1:** die DB-Adapter-Coverage ist ein anderer Messgegenstand | **eingehalten** — `tools/` im Vorgang **nicht** berührt (#20); die drei ausgenommenen Pakete bleiben ausgenommen (#4) |
| **Plan §1/§4:** Rückführungs-Bedingungen (Sondermodus ohne Naht nicht prüfbar · netzloser Anteil weit unter ≈54) | **nicht eingetreten** — alle vier Sondermodi sind netzlos prüfbar (#15), der Ertrag ist **58** Statements (`+49` `cmd`, `+9` `wiring.go`; **abgeleitet** 1581 − 1523) und liegt **über** der ADR-Schätzung ≈54 |
| **Plan §3:** `welle-20.md` §4 = **„nicht"** | **eingehalten** — die Datei hat **einen** Hunk, `@@ -34,2 +34,7 @@`, und er liegt in der **§1**-Prosa; die §4-Cluster-Tabelle ist unberührt. Die §1-Änderung selbst ist ein **Planner**-Zug nach Review F-5 (`3b7f8c7`, `f64794b`), kein stiller Eingriff |
| **Plan §3:** `harness/sensors/coverage-gate.md` — update, **nur falls** eine Zahl driftet **oder** eine Stelle bei Niederschrift eine bewegliche Größe als Ist-Stand führt | **eingehalten — die Bedingung greift** — zwei Anlässe, beide benannt: die drei durch diese Arbeit **wahr gewordenen** Berichtigungen in §Grenze 1 (siehe §6) und §Grenze 7 als **Selbst-Fund** (die Deckung von `cmd` hängt an der `GOCOVERDIR`-Weitergabe, #17). Keine Zahl driftete still |
| **[ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) §Konsequenzen, Folgepflicht (Implementer-Zug): „test-only … keine Zeile `internal/**` oder `cmd/**` außerhalb von Tests"** | **eingehalten** — wörtlich: die Abfrage über `internal/**`+`cmd/**` ohne `_test.go` ist leer (#20). Die vier übrigen Pfade liegen außerhalb dieser zwei Bäume (Sensor-Doku, Slice-Plan, `welle-20.md`, Review-Reports) und sind durch Plan §3 bzw. Modul 8 (Übergabe-Artefakte) gedeckt |
| **[ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) §Konsequenzen, Folgepflicht (Welle-Closure): „grün bei `THRESHOLD=80`, rot bei `THRESHOLD=85`"** | **eingehalten, soweit dieser Slice sie trägt** — beide Hälften sind real gefahren (#2, #3); die **Hochschaltung** selbst ist Wellen-Closure-Arbeit und **nicht** erfolgt (#21) |
| **[ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) §Konsequenzen: „Kein Produktionscode, kein Messmechanismus, keine Schwelle wird berührt"** | **eingehalten** — `THRESHOLD` unberührt (#21), die drei Messgrößen (Paketfilter, `-coverpkg`, `-covermode`) unberührt, `tools/coverage-gate.sh` unberührt (#20) |
| **[ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) Folgepflicht (Planner-Zug): die §4-Zeile der Welle berichtigen** | **erledigt außerhalb des Implementer-Diffs** — `welle-20.md` §1 trägt seit `3b7f8c7` den Planstand-Hinweis; die §4-Zeile selbst führt **Soll**-Zahlen und bleibt nach Plan §3 stehen |
| **`AGENTS.md` §3.5** (Accepted-ADRs immutable) | **eingehalten** — [ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) ist **nicht** im Vorgang (#20), also nicht in-place geändert |
| **`AGENTS.md` §3.6** (Schwellen nur per ADR) | **eingehalten** — keine Schwellen-Änderung (#21); die zwei Rampen-Läufe sind Overrides auf der Kommandozeile |
| **`AGENTS.md` §3.7** (Ist-Zustand, keine Chronik) in den neuen Sätzen | **eingehalten** — die neuen Kommentare und Sensor-Sätze beschreiben den geltenden Zustand („49 von 49", „Dienstgebunden ist der **Rumpf**"); keine Abwesenheits- oder Vorher/Nachher-Wendung |
| **`AGENTS.md` §3.2** (Suppression-Verbot) / **§3.11** (kein host-lokaler Pfad) | **eingehalten** — **0** Treffer für `nolint` in den zwei neuen Testdateien; `make docs-check` grün über den ganzen Vorgang (#1, #18, #19, #24) |
| **`AGENTS.md` §3.12** in den geänderten Trägern | **eingehalten** — die neuen Zahlen des Sensor-Dokuments tragen ihren Lauf („Lauf `slice-094`"); die Rückrechnungen sind als solche benannt; die `welle-20` §1-Passage trennt **Planstand** (327, 49), **Zitat** (336, ADR-Stand) und **eigenen Messwert** (299, 0) und gibt jedem seinen Ursprung — die Ungleichartigkeit, die der Delta-Review als zwei Herkünfte in **einer** Klammer gerügt hatte, ist damit behoben |

---

## 4. Plan-vs-Code-Diff

Verglichen gegen die §3-Liste des Plans am Stand `HEAD`.

| §3-Zeile | Geliefert | Urteil |
|---|---|---|
| `cmd/pg-change-feed/**_test.go` — Test neu, Argument-Dispatch über einen Re-Exec-Harness | `cmd/pg-change-feed/main_test.go` (neu, 282 Zeilen, 5 Testfunktionen) | **Plan eingehalten** |
| `internal/bootstrap/**_test.go` — Test neu/update, netzlos erreichbarer Teil des `Run`-Fehlerpfads | `internal/bootstrap/run_test.go` (neu, 62 Zeilen, 1 Testfunktion) | **Plan eingehalten** |
| `harness/sensors/coverage-gate.md` — update, nur unter der Bedingung | zweimal angefasst (`32b8b9d`, `d839975`) — beide Anlässe benannt (§3) | **Plan eingehalten** |
| `docs/plan/planning/welle-20.md` §4 — **nicht** | §4 nicht angefasst; **ein** Hunk in §1 (`3b7f8c7`, `f64794b`) | **Plan eingehalten** (§1 ist nicht die ausgeschlossene Stelle; die Änderung ist ein Planner-Zug nach Review F-5) |
| — | Slice-Plan `§2` (fünf Häkchen, `32b8b9d`) | **kein Plan-Bruch**: Pflege des eigenen Plans |
| — | das Review zu `slice-094` und dessen Delta-Review (neu) | **kein Plan-Bruch**: Übergabe-Artefakte der Reviewer-Rolle (Modul 8), keine Liefer-Punkte |

**Was der Diff nicht enthält, obwohl der Plan es nennt:** nichts. Die zwei
Test-Zeilen sind vollständig geliefert; Produktcode, `THRESHOLD`, `tools/` und
die §4-Tabelle sind pfadmäßig unberührt (#20).

---

## 5. Die Zahlen — selbst gemessen (`AGENTS.md` §3.12), und was in §7 gehört

| Zahl | Wo sie heute steht | Mein Lauf | Befund |
|---|---|---|---|
| `cmd/pg-change-feed` **49/49** | Commit-Message `8292766`, Sensor-Dokument §Grenze 1 (mit Lauf-Anker) | **49/49** gedeckt (#6), in **7** Profillläufen identisch; Parent **0/49** (#5) | **hält, exakt** — und sie ist der **stabile** Teil: das Paket kennt kein Band (`GOCOVERDIR`-abhängig, #17) |
| `wiring.go` **222/521** | Commit-Message `32b8b9d`/`8292766` | **222/521** (#4, 7 ×); Parent **213/521** (#5) | **hält**; das Band dieses Teilgegenstands reicht bis **221** (`zurückgerechnet` aus der gedruckten `83.0%` von #2, kein eigener Profillauf mit Block-Auswertung) |
| Gesamt **1581/1903** | Commit-Message `8292766` | **1581** (7 ×, #4); Parent **1523** (#5); die gedruckte Zeile `83.1%` in #1, #3, allen sieben Profilen und #24, `83.0%` in #2 | **hält als oberes Band-Ende**; das Band reicht bis **1580** (beobachtet, #2) |
| **83,08 %** | Commit-Message `8292766` (Review-Reports nennen dieselbe Zahl) | die **gedruckte** Zeile ist `83.1%` (zehn Läufe) bzw. `83.0%` (einer); `1581/1903 = 83,0794 %` (**zurückgerechnet**, keine gedruckte Zeile) | **hält als Rückrechnung** — im Träger gehört die **gedruckte** Zeile hin, mit ihrem Band |
| **+58** | **nirgends** — weder Commit-Message noch Träger (Aufgabe dieser Prüfung) | **abgeleitet**: 1581 − 1523 = **58**; das untere Band-Ende ergibt **57** | **hält als abgeleitete Zahl** — sie braucht beide Herkünfte (beide Läufe) und den Vermerk *abgeleitet* |
| Nenner **1903** | Sensor-Dokument §Zählbasis, Plan §1, [ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) | **1903** an **beiden** Ständen (8 Läufe, #4, #5) | **hält** — Zustandsgröße, unverändert |
| LP3-Verteilung **189/73/14/13/10** | Commit-Message `32b8b9d`; §1 des Plans nennt nur die **189** | **189 · 73 · 14 · 13 · 10 = 299**, außerhalb der fünf **0** (#7) | **hält, exakt** — die Differenz zur ADR-Tabelle (198 → 189) ist der in diesem Slice gewonnene Fehlerpfad des ersten Konstruktors |
| **vier Wellen-Belege** | [ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) §Konsequenzen | #2 und #3, je mit eigenem Exit | **hält, beide Hälften selbst gefahren** |
| `1532/1903 = 80,50 %` (Grenze 7) | Sensor-Dokument §Grenze 7 (mit Lauf-Anker) | **1532** (4 ×) / **1531** (1 ×), gedruckt 5 × `80.5%`, `cmd` 0/49, alles grün (#17) | **hält**; das Band ist **1** Statement breit (1531–1532) und druckt an **beiden** Enden `80.5%` |

**Was in §7 gehört** (LP1 verlangt es, `AGENTS.md` §3.12 wendet es an):

1. **Zustands-Größen** (hängen am Code-Stand, schwanken nicht): Nenner **1903**;
   `cmd/pg-change-feed` **49** Statements gesamt; `wiring.go` **521** gesamt.
2. **die Zuwachs-Zahl mit ihrem Lauf** — `cmd` **0 → 49**, `wiring.go`
   **213 → 222**, `internal/bootstrap` ungedeckt **336 → 299**, Gesamt gedeckt
   **1523 → 1581**; die Summe **+58** als **abgeleitet** kennzeichnen, mit beiden
   Läufen als Herkunft (Parent-Lauf und eigener Lauf), oder ganz weglassen.
3. **die erreichte Quote mit ihrem Lauf und ihrem Band** — gedruckt **`83.1%`**
   (oberes Ende) bzw. **`83.0%`** (unteres), gedeckt **1580–1581 von 1903**,
   Band **1** Statement; der Schwankungs-Träger steht im Sensor-Dokument
   §Zählbasis (`internal/bootstrap/wiring.go`, der `runAdministration`-Zweig).
4. **die zwei Rampen-Belege mit ihrer Schwelle** — grün bei `THRESHOLD=80`
   (EC **0**), rot bei `THRESHOLD=85` (Skript-EC **1** / make-EC **2**).
5. **die LP3-Aussage in ihrer gemessenen Form** — „**299** ungedeckte Statements
   liegen **sämtlich** in den fünf dienstgebundenen Funktionen (`Run` 189/198 ·
   `Diagnose` 73/81 · `Healthcheck` 14/22 · `RegisterConsumer` 13/17 ·
   `AcknowledgeConsumer` 10/14); außerhalb ihrer: **0**"; für `cmd`:
   **49 von 49**, die vier Aufrufe `main.go:44`/`:62`/`:86`/`:101` netzlos
   erreichbar.

### Die Frage der Welle: trägt der Puffer?

**Ja — und er ist jetzt belastbar, nicht mehr am Rand.** Nach `slice-093` stand
der Puffer bei `THRESHOLD=80` auf **einem** Statement. Gemessen heute: **1580 bis
1581** gedeckt, der grüne Beleg bei **80 %** läuft real durch (#2), und der
Abstand zur Schwelle ist **58** Statements (**abgeleitet**). Der Rot-Beleg bei
**85 %** bleibt bestehen (#3) — die Stufe prüft real, und die Decke aus
[ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
§Kontext (5) ist damit erreicht statt angenähert.

---

## 6. Die drei reparierten Sätze — `BEO-PGC/arbeit-ueberholt-stehenden-traeger`

Der Eintrag steht im Register bei **2×** (`evidence/slice-091.md`,
`evidence/slice-093.md`). Die Frage dieser Prüfung: **waren die drei Sätze bei
ihrer Niederschrift wahr und werden sie durch *diese* Arbeit falsch** — und macht
die Reparatur nichts Neues falsch.

| # | Satz (`§Grenze 1`) | Bei Niederschrift | Nach der Arbeit | Mein Beleg |
|---|---|---|---|---|
| 1 | „**Vier** Pakete des Gegenstands führen keine Testdatei" (und „genau **vier** Pakete `Test=0 XTest=0`") | **wahr** — am Parent `4ba09ff` liefert `go list` **vier** Pakete `Test=0 XTest=0`, darunter `cmd/pg-change-feed` (#10) | **falsch geworden** — der Re-Exec-Harness gibt `cmd` ein eigenes Testpaket; am Diff-Stand sind es **genau drei** (#9) | gemessen, beide Stände |
| 2 | „`TestGoFiles` allein trifft **23** der 31 Pakete" (und „die **23** zerfallen in diese vier und **19**") | **wahr für den Parent** (23 = 4 + 19) | **falsch geworden** — am Diff-Stand trifft `TestGoFiles=0` **22** Pakete, sie zerfallen in **3** + **19** (#9) | gemessen |
| 3 | „`cmd/pg-change-feed` trägt **49 Statements, alle mit `count = 0`** … das **einzige** Paket des Gegenstands ohne ein einziges gedecktes Statement … die Zeile `coverage: 0.0% of statements`" | **wahr** — Parent: `cmd` **0/49** gedeckt (#5), die Zeile `coverage: 0.0%` war die seine | **falsch geworden** — `cmd` ist **49/49** (#6); **kein** Paket druckt mehr `coverage: 0.0% of statements` (0 Treffer im Stufenlauf, #4) und **kein** Paket mit ausführbaren Statements trägt null gedeckte Statements (per Paket dedupliziert: **0** solche Pakete, #6) | gemessen |

**Macht die Reparatur nichts Neues falsch?** Die neuen und umgeschriebenen Sätze
der §Grenze 1 sind Satz für Satz nachgemessen:

- „Das Paket führt ein eigenes Testpaket (`TestGoFiles` = **1**,
  `XTestGoFiles` = **0**)" — **hält** (#9).
- „Das Paket trägt **49 von 49** Statements gedeckt (Lauf `slice-094`)" —
  **hält** (#6, 7 Läufe).
- „die Aufrufe `main.go:44` (`Healthcheck`), `:62` (`RegisterConsumer`),
  `:86` (`AcknowledgeConsumer`) und `:101` (`Diagnose`) sind netzlos erreichbar"
  — **hält**, je `count > 0` und die Zuordnung Zeile → Funktion stimmt (#8).
- „Dienstgebunden ist der **Rumpf** dieser vier Funktionen, nicht ihr Aufruf" —
  **hält**: der Aufruf-Block liegt in `cmd` (49/49), der Rumpf in den fünf
  Funktionen (#7).
- „genau drei Pakete `Test=0 XTest=0` … `TestGoFiles` allein trifft **22** … die
  **22** zerfallen in diese drei und **19**" — **hält**, Zahl für Zahl (#9).
- „die drei tragen **keine ausführbaren Statements** — im Profil kommen sie nicht
  vor" — **hält**: **0** Profilzeilen je Paket, der Lauf weist sie als
  `[no test files]` aus (3 Zeilen, #11, #4).
- „In der Aufrufform der Stufe trägt **kein** Paket die Zeile
  `coverage: 0.0% of statements`" — **hält** (0 Treffer im Stufenlauf, #4).

**Befund:** Alle drei Sätze waren bei Niederschrift **wahr** und sind durch
diese Arbeit **falsch geworden** — das ist der **dritte** Vorgang der Klasse
(Erstauftreten `slice-091`, zweiter `slice-093`), und die Reparatur macht
**nichts Neues falsch**. Ob daraus `evidence/slice-094.md` und damit der
Schwellen-Übertritt auf **3×** wird, ist die Register-Entscheidung des
Lese-Schritts (Modul 6: „Mensch urteilt, Maschine prüft Deckung"); dieser Bericht
liefert die Messung, die ihr zugrunde liegt. Der Träger-grep des Auftrags ist
**leer** ausgegangen: **kein** weiteres Dokument beschreibt die bewegte
Eigenschaft (#23) — außer dem §Kontext (4a) von
[ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md),
der den **Parent-Stand** mit Lauf-Label führt und als `Accepted`-ADR nach
`AGENTS.md` §3.5 **nicht** in-place zu ändern ist.

---

## 7. Findings

### V-1 — Die Fixrunde `d839975` ist von **keiner** Review gedeckt; §2 verlangt für sie ein Delta-Review

- `kategorie`: **MEDIUM**
- `quelle`: §2 des Slice-Plans, Zeile „Review durchgeführt, Report unter
  `docs/reviews/` liegt vor … **Weist der Review eine Fixrunde aus, deckt ein
  Delta-Review sie ab**" · Präzedenz im Verifikationsbericht zu `slice-092` V-2 (MEDIUM),
  Verifikationsbericht zu `slice-093` V-1 (MEDIUM)
- `pfad`: `docs/reviews/` (nur das Review zu `slice-094` und
  dessen Delta-Review) gegen `d839975`
  (`cmd/pg-change-feed/main_test.go` 4 Zeilen, `welle-20.md` 9 Zeilen,
  `harness/sensors/coverage-gate.md` 5 Zeilen)
- `befund`: Der Delta-Review beurteilt **`8292766`** und weist mit **D-1**
  („die Ordnungszahl löst für drei der vier Modi nicht auf"), **D-2** („verliert
  den einschränkenden Partikel") und dem ausdrücklich **offen** gelassenen
  **F-7** („`in derselben Aufrufform` ohne Antezedens") eine **weitere** Runde
  aus — wörtlich: „**Fixrunde:** **ja, aber nur ein Kommentar-Nachzug** — D-1
  (ein Wort) und D-2 (ein Partikel); er kann im selben Zug laufen wie der offene
  F-7-Satz." `d839975` führt genau diese drei aus (Commit-Text: „D-1 … D-2 …
  F-7 … Und mein Wellen-Nachzug … zwei Herkünfte"). Für **diesen** Stand
  existiert **kein** Review-Artefakt: `ls docs/reviews/ | grep 094` liefert
  **zwei** Dateien (#21), und die Hausform dieses Repos für „Fixrunde nach
  Review" ist ein **eigener** Nachtrag (Präzedenz bei den Delta-Reviews zu `slice-089`,
  `slice-091`, `slice-092`,
  und der zweiten Delta-Review-Runde zu `slice-093`).
- `verifizierbar`: ja — `ls docs/reviews/ | grep 094`;
  `git log --oneline 4ba09ff..HEAD`; `git show d839975`; das Verdikt des
  Delta-Reviews (`## Verdikt`, Zeile „Fixrunde")
- `urteil`: **Substanz geprüft, Artefakt fehlt.** Ich habe die Runde unabhängig
  nachgemessen — sie ist **tragend, nicht Prosa**: D-1 und D-2 sind im
  Testkommentar nachweislich umgesetzt (`:176` „Zweig, der die vier Sondermodi
  durchlässt statt sie abzuweisen"; `:185` „nicht auf den **blossen**
  Modus-Namen"), F-7 ist im Sensor-Dokument eingelöst (`:166` „In der
  Aufrufform der Stufe (`-coverpkg` über den ganzen Gegenstand)" — das
  Antezedens steht jetzt im Satz), und die zwei Herkünfte in der
  `welle-20`-Klammer sind getrennt benannt (§3). Die Severity bleibt trotzdem
  **MEDIUM**: das Kriterium steht im DoD und ist **nicht** erfüllt. Zwei Ausgänge
  sind zulässig — ein kurzer Delta-Nachtrag auf `d839975` oder die ausdrückliche
  §7-Zeile, dass die Runde die verlangte Korrektur wörtlich ausführt und ihre
  Wirkung in diesem Bericht nachgemessen ist. **Still bleiben darf es nicht.**

### V-2 — Zwei Überschriften des Slice-Plans sind im Text verdoppelt

- `kategorie`: **LOW**
- `quelle`: `AGENTS.md` §3.7 (ein Kommentar/Träger beschreibt, was da ist) ·
  Struktur-Erwartung der Vorlage
  (`.harness/baseline/v6.5.0/templates/docs/plan/planning/slice.template.md`)
- `pfad`: `docs/plan/planning/in-progress/slice-094-coverage-cluster-a.md:182`,
  `:205`
- `befund`: Zwei Sektionsüberschriften tragen ihren Text doppelt, ohne Trenner:

  ```text
  182 ## 6. Risiken und offene Punkte## 6. Risiken und offene Punkte
  205 ## 7. Closure-Notiz## 7. Closure-Notiz
  ```

  Gerendert lautet die Überschrift damit „6. Risiken und offene Punkte## 6.
  Risiken und offene Punkte"; der Anker ist nicht der, den die Vorlage und die
  Nachbardateien (`slice-093` u. a.: `## 6. …`, `## 7. …`) erzeugen. Der Zustand
  ist **nicht** von dieser Arbeit erzeugt, aber Teil **dieses** Vorgangs: die
  Datei ist seit ihrer Anlage so (`git show e160e26:…` zeigt beide Zeilen) und
  wurde in `32b8b9d` angefasst. **Kein Sensor sieht es** —
  `make docs-check` ist grün (#1), weil **niemand** den Anker verlinkt
  (`grep -rn 'slice-094-coverage-cluster-a.md#'` → leer, #22).
- `verifizierbar`: ja — `sed -n '182p;205p' <Slice-Plan>`;
  `grep -n '^## ' <Slice-Plan>`; `git show e160e26:docs/plan/planning/open/slice-094-coverage-cluster-a.md | grep -n '^## '`
- `urteil`: **zwei Wörter löschen** — die Verdopplung hinter `## 6.`/`## 7.`
  entfernen. Kein Liefer-Defekt, kein Gate-Fall; die Datei wandert bei der
  Closure nach `done/`, und dort liest jeder Lauf ihre Überschriften als
  Gliederung. Erledigung **vor** dem `git mv`.

### V-3 — Der Re-Exec-Harness hat keinen Watchdog; ein Kindprozess, der nicht endet, bringt den Gate-Lauf zum **Hängen** statt zum **Scheitern**

- `kategorie`: **INFO**
- `quelle`: §2 LP1 („die Tests … sind netzlos **grün**" — der Gate-Lauf muss
  entscheiden, nicht stehenbleiben) · Plan §8 (die Sichtung führt
  `BEO-PGC/test-integration-retention-timing-flake` als **3×**, „Schwelle
  erreicht", und benennt für diesen Slice die Zeitabhängigkeit als
  nächstliegenden Fehler) · `AGENTS.md` §3.12 Instanz B
- `pfad`: `cmd/pg-change-feed/main_test.go:101-120` (`fahreProzess`) gegen
  `harness/sensors/coverage-gate.md` §Grenze
- `befund`: `fahreProzess` setzt weder `context`/`CommandContext` noch ein
  Zeitlimit; gewertet wird `kind.Run()` — das **Prozess-Ende**. Das ist für den
  Normalfall die richtige Form (die Tests sind end-gebunden, **0** Treffer für
  `Sleep`/`Timeout`/`Deadline`, #13), aber es gibt **keine** Obergrenze: endet ein
  Kindprozess nicht, steht der Lauf, statt rot zu werden. Die einzige Annahme,
  unter der das eintreten kann, ist ein **Listener auf Port 1** — in der
  netzlosen Umgebung der Stufe (`--network none`) ausgeschlossen, in einer
  Umgebung mit Netz aber nicht. Weder der Slice-Plan noch das Sensor-Dokument
  nennen die Grenze; §Grenze 7 benennt die andere, gemessene Grenze
  (`GOCOVERDIR`).
- `verifizierbar`: ja — `sed -n '101,120p' cmd/pg-change-feed/main_test.go`;
  `grep -nE 'Context|Timeout|Deadline|Sleep' cmd/pg-change-feed/main_test.go`;
  `grep -n '^7\.' harness/sensors/coverage-gate.md`; #13 (20 Läufe grün)
- `urteil`: **keine Reparatur verlangt** — die Grenze ist in der Gate-Umgebung
  nicht erreichbar, und der Reviewer hat sie als Restrisiko benannt. Sie steht
  hier, damit sie beim Closure-§7-Eintrag **nicht still** bleibt: sie ist die
  eine Annahme, auf der LP1s „netzlos grün" ruht. Will die Closure sie führen,
  ist es ein Halbsatz.

---

## 8. Negativbefunde

- **geprüft, ohne Befund: LP2 bindet wirklich an die Eingabeseite — acht eigene
  Mutationen, acht rot, Kontrolle grün** (#15, #16). `--healthcheck` →
  `cfg.AdminDSN` färbt genau den `--healthcheck`-Subtest; `diagnose` →
  `cfg.AdminDSN` den `diagnose`-Subtest; `RegisterConsumer` → `cfg.ReaderDSN` den
  `register-consumer`-Subtest; ein falsch benannter fehlender DSN den
  `--healthcheck`-Subtest **und** den bestehenden `TestConfigFromEnvOhneVorbedingung`;
  `os.Exit(2)` → `os.Exit(3)` den `unbekanntes_Argument`-Subtest.
- **geprüft, ohne Befund: der Übergang „grün → rot" der F-2-Probe ist real.**
  Dieselbe Mutation (moduseigene Meldung durch den Fallback-Text ersetzt) am
  Teststand `32b8b9d` **EC 0**, am Diff-Stand **EC 1** mit genau den zwei
  `acknowledge-consumer`-Subtests (#16) — die Fixrunde hat die Ausgabe-Hälfte
  wirklich gebunden, nicht nur den Kommentar geändert.
- **geprüft, ohne Befund: die Bindung von `run_test.go` trägt über den Sentinel,
  nicht über den DSN-Namen.** Eigene Mutation: `Run` erreicht zuerst den
  Schema-Store (`ErrSchemaStoreStorage`) → **bootstrap rot, cmd grün** (#15, M6).
  Das bestätigt den Satz des Testkopfs, dass ein verschluckender `Run` eine Stufe
  später endet.
- **geprüft, ohne Befund: LP3 hält exakt.** **134** ungedeckte Blöcke, **299**
  Statements, **sämtlich** in den fünf Funktionen; außerhalb ihrer **0**. Die
  Verteilung **189/73/14/13/10** stimmt Stück für Stück, und die Differenz zur
  ADR-Tabelle (`Run` 198 → 189) ist der in diesem Slice gewonnene Fehlerpfad
  (#7).
- **geprüft, ohne Befund: alle Zahlen des Vorgangs halten.** `49/49`,
  `222/521`, `1523`, `1581`, Nenner `1903`, `299` — jede einzeln nachgemessen,
  jede exakt (#4, #5, #6, #7); `83,08 %` ist die **zurückgerechnete** Form der
  gedruckten `83.1%`, `+58` ist **abgeleitet** und hält.
- **geprüft, ohne Befund: die zwei Rampen-Belege sind beide real.** Grün bei
  `THRESHOLD=80` (EC **0**), rot bei `THRESHOLD=85` (Skript-EC **1**, make-EC
  **2**), je selbst gefahren (#2, #3).
- **geprüft, ohne Befund: `make gates` trägt den ganzen Vorgang.** EC **0** aus
  separater Datei gelesen, sechs Checks grün (#1), erneut grün mit diesem
  Bericht im Baum (#24); `make doc-commits` und `make doc-immutable` über die
  **volle** Slice-Range je **EC 0**, 775 Dateien, 0 Befunde (#18, #19) — das ist
  nötig, weil das Standing-Gate nur die letzten **5** Commits liest
  (`commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"`) und die zwei ersten
  Commits des Vorgangs sonst ungeprüft blieben.
- **geprüft, ohne Befund: die neuen Tests sind netzlos, deterministisch und
  laufen im Gate.** **6** neue Funktionen, jede **1 ×** `=== RUN`, **0 × FAIL**
  (#12); **20** Wiederholungen ohne Flake (#13); `make test -race` **EC 0,
  33 × ok, 0 × DATA RACE** (#14). Die zeitabhängige Stelle
  (`runWALRetentionCheck`) hat in 60 Läufen nicht gefeuert; das Risiko
  `test-integration-retention-timing-flake` ist für die neuen Tests **nicht**
  eingetreten.
- **geprüft, ohne Befund: §Grenze 7 hält, beide Zahlen.** Ohne die
  `GOCOVERDIR`-Weitergabe bleiben **alle** Tests grün und `cmd` fällt auf
  **0/49**; gemessen **1532** (4 ×) / **1531** (1 ×), gedruckt 5 × `80.5%` (#17)
  — die Zahl des Dokuments ist das obere Band-Ende und druckt an beiden Enden
  gleich.
- **geprüft, ohne Befund: die drei reparierten Sätze der §Grenze 1 waren bei
  Niederschrift wahr und sind durch diese Arbeit falsch geworden** (§6) — der
  dritte Vorgang der Klasse, nachgemessen an beiden Ständen (#5, #9, #10).
  Weiterer Träger außerhalb des Diffs: **keiner** (#23); der §Kontext (4a) von
  [ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
  führt den Parent-Stand mit Lauf-Label und ist als `Accepted`-ADR nach
  `AGENTS.md` §3.5 unberührbar.
- **geprüft, ohne Befund: kein Produktcode, `THRESHOLD` unberührt, kein
  Messmechanismus berührt.** **7** Pfade; die vier Abfragen
  `internal/**`+`cmd/**` ohne `_test.go`, `harness/mk`+`Makefile`+`Dockerfile`,
  `tools/` und `docs/plan/planning/observations/` sind **leer** (#20);
  `THRESHOLD ?= 70` steht (#21). Alle Mutationen liefen auf Arbeitsbaum-Kopien
  außerhalb des Baums.
- **geprüft, ohne Befund: der Plan-vs-Code-Diff ergibt keine unbegründete
  Abweichung** (§4). Die zwei Test-Zeilen sind vollständig geliefert; die
  `welle-20.md`-Änderung liegt in §1 (nicht in der ausgeschlossenen §4-Tabelle)
  und ist ein Planner-Zug nach Review F-5; die zwei weiteren Pfade sind
  Reviewer-Übergabe-Artefakte.
- **geprüft, ohne Befund: `AGENTS.md` §3.2, §3.5, §3.6, §3.7, §3.11.** Kein
  `//nolint` in den neuen Tests, keine in-place-Änderung an einer
  `Accepted`-ADR, keine Schwellen-Änderung, kein host-lokaler Pfad (dieser
  Bericht ebenso); `make docs-check` grün über den ganzen Vorgang (#1, #18, #19,
  #24).

---

## 9. Was ich nicht prüfen konnte

- **Die Läufe des Implementers und der Reviews.** Die Commit-Messages nennen
  „14 Mutationen", der Review „12 Mutations-Proben", der Delta-Review „19 neue
  bzw. umgeschriebene Sätze". Beider Läufe sind nicht mehr einsehbar; ich habe
  die **Substanz** unabhängig gemessen (7 eigene Profile, 8 eigene
  Mutationsläufe über zwei Stände, 5 Läufe der §Grenze-7-Probe) — ein
  **Gegenbeispiel** hat sich nicht gezeigt.
- **Der 1580-Lauf im Detail.** Die untere Band-Grenze ist **beobachtet** (die
  gedruckte `83.0%` von #2), aber ich habe sie in **sieben** eigenen Profilen
  nicht reproduziert; die Zuordnung des fehlenden Statements zum
  `runAdministration`-Zweig ist **zurückgerechnet**, nicht von mir blockgenau
  gemessen (§5).
- **Die dritte Review-Runde.** `docs/reviews/` trägt **zwei** 094-Reports (#21);
  ein Nachtrag auf `d839975` ist zum Zeitpunkt dieser Messung **nicht** im Baum.
  **V-1** bleibt damit offen, bis er vorliegt oder §7 ihn ausdrücklich abdeckt.
- **`make test-store`, `make test-replication`, `make test-notify`,
  `make test-integration`, `make image`.** Kein Gate, kein Bestandteil der
  Lieferung; die neue Bindung ist netzlos geführt (das ist ihr Punkt). Der
  Integrations-Tier, der `Run` real fährt, wurde **nicht** gefahren.
- **Der reale Post-Push-Lauf.** Dieser Vorgang ändert **keinen** Workflow —
  `AGENTS.md` §3.10 greift dem Buchstaben nach nicht. Die Wirkung der neuen
  Tests auf die CI-Laufzeit bleibt bis zum nächsten echten Lauf unbelegt.
- **Die Stabilität über die Rampe.** Mein Befund „1580–1581, Puffer 58
  Statements" gilt für **diesen** Stand; ob die Stufe bei 80 % dauerhaft grün
  bleibt, entscheidet die `welle-20`-Closure mit ihrem eigenen Lauf.
- **Der Register-Zähler nach diesem Vorgang.** Ob
  `arbeit-ueberholt-stehenden-traeger` auf **3×** geht und welchen Ausgang der
  Lese-Schritt zuweist, ist eine Register-Entscheidung (Modul 6: „Mensch
  urteilt, Maschine prüft Deckung") — §6 liefert nur die Messung.
- **Die drei Paarungen.** Der `welle-20`-Closure zugewiesen (§2-Zeile).

---

## 10. Verdikt

**Die drei Liefer-Punkte aus §2 tragen — gemessen, nicht gelesen.**

- **LP1:** **6** neue Testfunktionen laufen real im netzlosen Lauf (#12, EC **0**,
  0 × `FAIL`) **und** in der `coverage`-Stufe des Gates (#4); der Zuwachs ist
  **+58** Statements (`cmd` **0 → 49**, `wiring.go` **213 → 222**, `bootstrap`
  ungedeckt **336 → 299**) bei **unverändertem Nenner 1903**, die Quote ist die
  **gedruckte** Zeile `83.1%` (unteres Band-Ende `83.0%`), das Band **1580–1581**
  = **1** Statement; **beide** Rampen-Belege sind selbst gefahren (EC **0** bei
  80, Skript-EC **1** / make-EC **2** bei 85). Der **Träger** der Zuwachs-Zahl
  und des Bands ist offen — §5 legt ihn vor.
- **LP2:** Die Exit-Codes sind gebunden (eigene Mutation `os.Exit(2)` →
  `os.Exit(3)` rot), und die Ausgabe ist an ihre **Eingabeseite** gebunden:
  **acht** eigene Mutationen, **acht** rot — darunter die verlangte
  `--healthcheck` → `cfg.AdminDSN` (#15) und der **selbst gefahrene** Übergang
  „grün am Stand `32b8b9d` → rot am Diff-Stand" (#16). Das ist der Kern dieses
  Slice, und er hält.
- **LP3:** **299** ungedeckte Statements liegen **sämtlich** in den fünf
  dienstgebundenen Funktionen (`Run` 189/198 · `Diagnose` 73/81 · `Healthcheck`
  14/22 · `RegisterConsumer` 13/17 · `AcknowledgeConsumer` 10/14); außerhalb
  ihrer **0**. Die vier `cmd`-Aufrufe `main.go:44`/`:62`/`:86`/`:101` sind
  gedeckt (#8).

**Entscheidungs-Konformität: hält.** Kein Produkt-Code (**7** Pfade, **2**
Testdateien; `internal/**` und `cmd/**` außerhalb der Tests leer), `THRESHOLD ?=
70` unberührt, kein Messmechanismus berührt, `tools/` unberührt,
`cmd/pg-change-feed/main.go` unberührt — der Re-Exec-Harness trägt die 49/49
**ohne** eine Zeile Produktionsänderung, genau wie
[ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
es verfügt. Die Berichtigungen am Sensor-Dokument sind durch die Messung gedeckt
und tragen ihren Lauf; die `welle-20` §1-Passage liegt außerhalb der
ausgeschlossenen §4-Tabelle.

**Der Plan-vs-Code-Diff ergibt keine unbegründete Abweichung.**

**`done/`-fähig ist der Slice damit noch nicht — vier Punkte fehlen:**

1. **V-1** — die Fixrunde `d839975` hat **kein** Review-Artefakt; §2 verlangt für
   sie einen Delta-Nachtrag. **Das ist der einzige Punkt dieser Liste, der ein
   DoD-Kriterium direkt verletzt.**
2. **§7** — die Closure-Notiz (7 von 7 Zeilen Platzhalter), darin **die Zahlen
   mit Lauf und Band** und **beide** Rampen-Belege mit ihrer Schwelle (§5) sowie
   die LP3-Aussage in ihrer gemessenen Form.
3. **V-2 und V-3** — zwei Wörter in §6/§7 des Slice-Plans (Doppel-Überschriften)
   und die benannte Watchdog-Grenze als Halbsatz, wenn die Closure sie führen
   will. Beide sitzen in Dateien **dieses** Vorgangs; **V-2** gehört **vor** den
   `git mv`.
4. **Die regulären Closure-Pflichten** — Register-Entscheidung
   (`arbeit-ueberholt-stehenden-traeger` dritter Vorgang,
   `beleg-befehl-traegt-seinen-satz-nicht` Auslöser; Substanz in §6), die **vier
   §6-Ausgänge** (Material in §5/§8: R1 *entfallen*, R2 *entfallen*, R3
   *entfallen mit Gegenbeleg*, R4 *eingetreten und bedient*), die Häkchen, die
   **drei Paarungen** (der `welle-20`-Closure).

Nach 1–4 ist der Slice `done/`-fähig. **Kein Liefer-Defekt, kein rotes Gate:**
die offenen Punkte sind ein fehlendes Übergabe-Artefakt (V-1), zwei Wörter
(V-2), eine benannte Grenze (V-3) und die regulären Closure-Pflichten — keine
gebrochene Zusage.

---

**Beleg-Lage dieses Berichts:** jede Zahl stammt aus einem der Läufe in §1, je in
eigener Werkzeug-Beauftragung gefahren; die vier Gate-Läufe (#1, #2, #3, #24) und
ihre Auswertung waren **zwei** Schritte, ihre Exit-Codes wurden aus separaten
Dateien gelesen, nie durch eine Pipe (`AGENTS.md` §3.9); #24 lief **mit** diesem
Bericht im Baum und ist grün. Die Mutationsproben liefen auf Arbeitsbaum-Kopien
**außerhalb** des Repos (`git archive`), netzlos; die Repo-Dateien wurden
**nicht** angefasst (nach jeder Probe gegen das Original geprüft). Der Baum
trägt nach diesem Bericht allein diesen Bericht — **kein** Commit, keine
Änderung an Artefakten des Slice, an `THRESHOLD`, an Produktcode oder an einem
Träger außerhalb dieses Berichts.
