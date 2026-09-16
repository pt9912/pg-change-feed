# Review-Report: slice-094 — **Delta** (Fixrunde `8292766`) · 2026-09-16

**Review-Art:** Code (Delta) — geprüft gegen **Plan und Entscheidungen**. DoD-/Spec-Konformität
(Verifier), §6-Ausgänge, Register und die drei Paarungen sind **nicht** Gegenstand.

**Gegenstand:** `8292766` (Parent `3b7f8c7`) — **3** Dateien, +85/−29:
`cmd/pg-change-feed/main_test.go`, `internal/bootstrap/run_test.go`,
`harness/sensors/coverage-gate.md`. Der **Code**-Anteil gegen den Erst-Review-Stand ist auf zwei
`*_test.go` begrenzt (`git diff --stat 32b8b9d..8292766 -- cmd internal` → 2 Dateien, +77/−23).
Kein Produktionscode, kein `THRESHOLD`, keine Naht, keine ADR. `f64794b` (Planner-Nachzug in
`welle-20.md` §1) ist **nicht** Gegenstand — Stellungnahme unter §Schwerpunkte (5).

**Skill:** `.harness/skills/reviewer.md` @ `8292766` · **Datum:** 2026-09-16.

---

## Findings

### D-1 — Die Ordnungszahl „den **vierten** Dispatch-Zweig jeder Modus-Verzweigung" löst für drei der vier Modi nicht auf

- `kategorie`: **LOW** · `quelle`: §3.12 Instanz B (Beleg-Anker)
- `pfad`: `cmd/pg-change-feed/main_test.go:175-176`
- `befund`: Der neue Testkopf sagt, der Test trage „den **vierten** Dispatch-Zweig **jeder**
  Modus-Verzweigung". Gezählt an `main.go` gibt es **vier** Exit-Punkte nur im
  `acknowledge-consumer`-Block (`:74`, `:79`, `:84`, `:86`); die drei anderen tragen zwei, drei
  und zwei (`--healthcheck` `:42`/`:44`, `register-consumer` `:55`/`:60`/`:62`, `diagnose`
  `:99`/`:101`). Der Satz nach dem Doppelpunkt ist richtig und gemessen — **die Ordnungszahl**
  trägt ihn nicht.
- `verifizierbar`: **ja** — Exit-Punkte je Modus-Block in `main.go`
- `klasse`: „Ordnungszahl im Kommentar, die der Code nicht auflöst"

### D-2 — „nicht auf den Modus-Namen" verliert den einschränkenden Partikel, den der Satz braucht

- `kategorie`: INFO · `quelle`: §3.12 Instanz B
- `pfad`: `cmd/pg-change-feed/main_test.go:184-185` (gegen die Assertion-Tabelle `:194-197`)
- `befund`: „Die Zeile wird zusätzlich auf die moduseigene Präfix-Form geprüft, **nicht auf den
  Modus-Namen**" — jede geprüfte Zeichenkette **beginnt** mit dem Modus-Namen; gemeint ist „nicht
  auf den **blossen** Modus-Namen", wie es der Schwester-Kommentar derselben Datei präzise sagt
  (`:225-226`). Die Assertion selbst ist richtig und gebunden. `beleg-befehl-traegt-seinen-satz-nicht`
  ist **nicht** getroffen — der Satz ist enger gemeint, als er dasteht.
- `klasse`: „verlorener Qualifikator in einem neuen Kommentar-Satz"

---

## Delta-Bilanz — was die Fixrunde schließt, was offen bleibt

| Finding | Stand nach `8292766` | Beleg |
|---|---|---|
| **F-1** (HIGH) | **geschlossen** | Sensor-Satz Klausel für Klausel nachgemessen: `cmd` **49/49**, Gesamt **1581/1903 = 83,08 %**, gedruckt `83.1%`; vier Aufruf-Blöcke tragen `count > 0` |
| **F-2** (LOW) | **geschlossen** | Probe: moduseigene Meldung durch den Fallback ersetzt → **rot, EC 1**; derselbe Mutationsstand mit dem Test-Stand `32b8b9d` → **grün, EC 0** (der Übergang ist real) |
| **F-3** (LOW) | **geschlossen** | `Run` verschluckt den ersten Konstruktor-Fehler → **cmd grün, `bootstrap` rot** |
| **F-4** (LOW) | **geschlossen** | `PATH` benannt; die zweite Ausnahme (`GOCOVERDIR`) steht im Folgesatz |
| **F-6** (INFO) | **geschlossen** | Haltepunkt „eine Stufe später am Schema-Store" gemessen (`ErrSchemaStoreStorage`) |
| **F-5** (INFO) | **vom Planner nachgezogen** (`f64794b`) — Enthaltung richtig, s. (5) |
| **F-7** (INFO) | **unverändert offen** — `coverage-gate.md:166` trägt weiter „in derselben Aufrufform" ohne Antezedens; geht an die Closure |

## Negativbefunde

- **Der neue Test fährt wirklich alle vier Sondermodi.** `go test -v -run
  TestSondermodiMitVollstaendigerUmgebungNennenIhreRolle` → **4 Subtests RUN + PASS**, EC 0.
- **Die Ausgabe ist an die Eingabeseite gebunden — fünf eigene Mutationen, jede rot.**
  `--healthcheck` → `cfg.AdminDSN`: rot. `diagnose` → `cfg.AdminDSN`: rot. `--healthcheck` ruft
  `Diagnose`: rot. `RegisterConsumer`/`AcknowledgeConsumer` → `cfg.ReaderDSN`: rot. Reale
  stderr-Zeilen: `database=reader-db` (healthcheck), `database=admin-db` (register),
  `database=admin-db` (acknowledge), `database=reader-db` (diagnose).
- **Der Harness bleibt netzlos und deterministisch.** Alle Läufe `--network none`;
  `-count=20` → EC 0; **fünf** volle Profil-Läufe → fünfmal identisch `cmd 49/49`, `1581/1903`,
  28 × ok, 0 × FAIL. Kein `time.`/`Sleep`/`After(`/`Timeout`/`Deadline` in den neuen Zeilen.
  `make test` (-race, netzlos) → **EC 0**, **33 × ok, 0 × FAIL, 0 × DATA RACE**.
- **§Grenze 7 hält unverändert — nachgemessen.** `kindUmgebung` ohne `GOCOVERDIR`-Weitergabe:
  alle Tests grün, `cmd` **0/49**, Gesamt **1532/1903 = 80,50 %**, gedruckt `80.5%`.
- **„Dienstgebunden ist der **Rumpf** dieser vier Funktionen, nicht ihr Aufruf" — die gemessene
  Gestalt.** Vor und nach der Fixrunde: `wiring.go` **222/521**, Funktions-Quoten
  `23,5 / 28,6 / 36,4 / 9,9 %` **unverändert**; bewegt haben sich genau die vier Aufruf-Blöcke
  `main.go:44/:62/:86/:101` (je `count=1`). Der Satz behauptet weder mehr noch weniger.
- **`run_test.go`s neuer Satz hält.** Mutiertes `Run`: Fehler trägt `database=capture-db`, der
  Test endet an `run_test.go:54` — rot über die **Klassen-Prüfung**, wie der Kommentar sagt.
- **`TestArgumentFehlerEndenMitAusgang2` — die vier neuen `nennt`-Werte sind byte-genau.**
- **Die Nachbar-Aussagen des Sensor-Dokuments halten:** `go list` → **31** Pakete, **3** ×
  `Test=0 XTest=0`, **22** × `TestGoFiles=0`; **3** `[no test files]`-Zeilen; **0** ×
  `coverage: 0.0% of statements`.
- **Die Schwankungs-Behauptung der §Zählbasis („höchstens 1 Statement") hält:** zwei Gate-Läufe
  desselben Stands druckten `83.10%` und `83.00%` → Spanne **1** Statement.
- **§3.2 / §3.5 / §3.6 / §3.11 — geprüft, ohne Befund.** Kein `//nolint`, keine Accepted-ADR
  berührt, kein `THRESHOLD`, kein host-lokaler Pfad. Betreff mit `ADR-0082`, ohne Struktur-ID.
- **Rollen-/Scope-Grenze gehalten.** Kein Produktionscode, keine Datei außerhalb der §3-Liste.

## Eigene Messungen (Exit-Codes direkt, ungepiped; Mutationen in Kopien außerhalb des Repos)

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `make gates` (HEAD) | 0 | `d-check` **774/0** · coverage-gate OK `83.10%` @70 · übrige grün |
| 2 | `make docs-check` | 0 | 774 Dateien, 0 Befunde |
| 3 | **Mutation:** `Run` verschluckt den ersten Konstruktor-Fehler | cmd **0** / bootstrap **1** | die Aufteilung, die der neue Satz behauptet |
| 4 | **F-2-Probe:** moduseigene Meldung → Fallback-Text | **1** | rot, genau die zwei `acknowledge-consumer`-Subtests |
| 5 | **Dieselbe Mutation**, Testfile-Stand `32b8b9d` | **0** | **grün** — „vorher grün, jetzt rot" real |
| 6 | **Mutation:** `--healthcheck` erhält `cfg.AdminDSN` | **1** | rot (`database=admin-db`, eigene Rolle fehlt) |
| 7 | **Mutationen:** `diagnose`→`AdminDSN` · `--healthcheck`→`Diagnose` · register/ack→`ReaderDSN` | je **1** | je rot auf den betroffenen Subtests |
| 8 | **§Grenze-7-Probe:** ohne `GOCOVERDIR` | 0 | alle Tests grün; `cmd` **0/49**, **1532/1903 = 80,50 %** |
| 9 | `make coverage-gate THRESHOLD=80` | 0 | `OK — Coverage 83.10%` |
| 10 | `make coverage-gate THRESHOLD=85` | **2** | `FAIL — Coverage 83.00% unter Schwelle 85%` |
| 11 | `make test` (-race, `--network none`) | 0 | **33 × ok**, 0 × FAIL, 0 × DATA RACE |
| 12 | `go test -count=20 ./cmd/pg-change-feed/` | 0 | 20 × grün, kein Flake |
| 13 | 5 × Profil-Lauf (block-dedupliziert) | 0 | 5 × `cmd 49/49`, `1581/1903`; Stand `32b8b9d`: `45/49`, `1577/1903` |
| 14 | `go tool cover -func` (vor/nach) | — | Funktions-Quoten **unverändert**; `wiring.go` 222/521 |
| 15 | Reale Modus-Starts mit vollständigem ENV | 1 | vier stderr-Zeilen mit den Rollen-DB-Namen |
| 16 | `go list`-Gruppierung + Profil-Zeilen | 0 | 31 · 3 · 22 · 3 `[no test files]` · 0 × `coverage: 0.0%` |

## Antwort auf die Schwerpunkte

**(1) F-1 — die Kernkorrektur trägt, vierfach gemessen.** (a) alle vier Modi gefahren; (b) an die
**Eingabeseite** gebunden, fünf Mutationen rot; (c) die verlangte Mutation
`--healthcheck` → `cfg.AdminDSN` **färbt rot**; (d) der Harness bleibt netzlos und
deterministisch.

**(2) F-2 — richtig gebunden, nicht zurückgenommen.** Der Übergang „grün → rot" ist real
gemessen (Messungen 4 und 5). Die vier neuen `nennt`-Werte sind byte-genaue Teilstrings.

**(3) F-3/F-4/F-6 — die neuen Sätze stimmen und sind nicht zu weit.**

**(4) §Grenze 7 — selbst gefahren, beide Zahlen halten.**

**(5) Die Enthaltung des Implementers war richtig — und der Nachzug ist stabil, aber nur als
Record.** *Zur Enthaltung:* `welle-20.md` §1 liegt außerhalb der §3-Datei-Liste des Slice und
außerhalb des Reviewer→Implementer-Pfeils; F-5 war INFO und nannte die Zuständigkeit selbst; und
der Implementer hat **nicht geschwiegen**, sondern gemeldet — das ist das Übergabe-Artefakt, das
Modul 8 verlangt („kein Rollenwechsel ohne Artefakt"). Eine Fremd-Datei still mitzuändern wäre
genau der blinde Übergang. *Zum Nachzug:* Der Satz nennt den **Vorgang** (`slice-094`) und ist
damit eine Aussage über das Ergebnis eines Laufs statt ein freistehender Ist-Stand — die Form, die
§3.12 Instanz A verlangt. Umstoßen kann ihn nur ein Leser, nicht die nächste Runde:
(i) `slice-094` ist **noch nicht geschlossen** — der Satz nimmt vorweg, was die Closure-Notiz laut
Folgesatz trägt; beide Stellen müssen bei der Closure übereinstimmen. (ii) Legt ein späterer Slice
eine ungedeckte Zeile in `cmd/pg-change-feed/main.go` an, bleibt „hat auf 0 gebracht" wahr,
verlangt aber, ihn als Record zu lesen. **Eine Ungleichartigkeit bleibt:** die Parenthese stellt
`(336)` (Zitat aus §4, Herkunft dort benannt, ADR-Stand) neben die `0` (eigener Messwert dieses
Slice) — zwei **verschiedene Ursprünge** in einer Klammer. Gemessen: `internal/bootstrap` trägt am
heutigen Stand **299** ungedeckte Statements, nicht 336; wer die 336 als Ist-Wert liest, liest sie
falsch (§4s Herkunfts-Block sagt es; deshalb kein Finding).

**(6) Ja — zwei neue Sätze sind falsch geworden, beide am Rand.** Von **19** neuen bzw.
umgeschriebenen Sätzen tragen **17**; die zwei Defekte sind Wort-Ränder (D-1, D-2), **kein**
Mechanismus betroffen. Die **vier** Sätze, die eine Mutation behaupten, sind alle nachgefahren und
**halten**.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | **0** |
| MEDIUM | 0 |
| LOW | **1** |
| INFO | 1 |

## Verdikt

**Die Fixrunde trägt.** F-1 ist geschlossen und nicht nur behauptet; F-2/F-3/F-4/F-6 sind
geschlossen und je durch eine eigene Mutation belegt; F-7 bleibt als INFO offen (geht an die
Closure).

**Merge-blockierend:** **nein** — 0 HIGH, 0 MEDIUM.

**Fixrunde:** **ja, aber nur ein Kommentar-Nachzug** — D-1 (ein Wort) und D-2 (ein Partikel);
er kann im selben Zug laufen wie der offene F-7-Satz. Kein Mechanismus, keine Zahl, kein Gate
berührt.

**DoD-Häkchen „Review durchgeführt":** bleibt **offen**.

**Nicht gefahren:** `make test-store`, `make test-replication`, `make test-notify`,
`make test-integration`, `make image` — kein Build-Kontext, kein `THRESHOLD`, keine Naht und kein
Produktionscode berührt. `make test` **wurde** gefahren (EC 0).
