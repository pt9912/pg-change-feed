# Verifikationsbericht: slice-081 — 2026-09-15

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag (`slice-081` §2) und die im Slice referenzierten Entscheidungen
([`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 5 — die Naht ist **design**-begründet —,
[`ADR-0077`](../plan/adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md)
— die Neu-Bemessung beim Subjekt-Transfer —,
[`ADR-0078`](../plan/adr/0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
— Teil-Supersede, Transfer-Nachweis statt Summen-Konstanz —,
[`ADR-0041`](../plan/adr/0041-a-check-maschinenform-architekturpruefung.md)
— `.a-check.yml` unverändert —; `AGENTS.md` §3.9, §3.10).
**Nicht** gegen den Diff als solchen (Reviewer-Aufgabe; die zwei Reports
`review-slice-081.md` und `review-slice-081-fixrunde.md` wurden als Kontext
gelesen, **nicht** als Beleg übernommen) und **nicht** gegen realen Bedarf
(Validator, hier nicht ausgelöst).

**Frischer Kontext:** Dieser Lauf hat den vollständigen Slice-Plan in der
Fassung von `HEAD` gelesen, dazu die drei ADRs und die berührten Träger.
**Alle** Zahlen dieses Berichts stammen aus eigenen, in dieser Sitzung
gefahrenen Läufen bzw. aus einer eigenen Zählung; kein Beleg des
Implementers oder des Reviewers wurde übernommen. Exit-Codes je in eigenem,
ungepiptem Schritt gelesen (`AGENTS.md` §3.9); Gate-Lauf und Folgehandlung
getrennt beauftragt. Die Baumänderung, die die schema-rollenden Testläufe an
`tools/schema/plan.yaml` hinterlassen, wurde nach **jedem** Lauf real
zurückgenommen; der Arbeitsbaum ist am Ende dieses Laufs unverändert
(`git status --porcelain` leer). Ein für die Vor-Zahlen angelegter Worktree
wurde real abgeräumt (`git worktree list` = nur der Hauptbaum). Eigene
Testcontainer und Docker-Netze wurden durch die Skripte selbst abgeräumt.

**Gegenstand:** `slice-081`
(`docs/plan/planning/in-progress/slice-081-executor-naht.md`), geprüfter
Stand `HEAD` = `f5ba274`. Der Slice liegt weiterhin in `in-progress/`; der
`git mv` nach `done/` ist **nicht** erfolgt.

---

## 1. Eigene Messungen dieses Laufs

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `make gates` (ungepiped, Ausgabe in Logdatei, Exit direkt gelesen) | **0** | baseline-verify `v6.5.0` OK (54 Dateien) · d-check **681** Dateien / **0** Befunde · commit-traceability OK (5 Commits, Betreffe ohne Struktur-ID) · a-check **0** Befunde · coverage-gate OK — `Coverage 71.30% erfüllt Schwelle 70%` |
| 2 | `make test-store` | **0** | `ok …/postgresstorage 4.158s` · `DB-Adapter-Coverage: 73.38% (gedeckt 477 von 650 Statements; Profile gemergt: store,replication)` · `db-coverage: OK … erfuellt Schwelle 70%` |
| 3 | `make test-replication` | **0** | `ok …/postgresack 0.023s` · `ok …/replication/receive 4.152s` (Messphase) bzw. `4.348s` (Tierphase) · dieselbe Zahl `73.38%` · `db-coverage: OK … Schwelle 70%` |
| 4 | Eigene Zählung `merged.coverprofile` (Dedup über Block-Position) | — | `postgresstorage=472` · `postgresack=23` · `replication/receive=155` · total **650**, gedeckt **477** → **73,38 %** |
| 5 | Eigene Zählung `store.coverprofile`, je Datei (Dedup) | — | `tableactivation 132 · store 93 · consumerstate 87 · administrationrequest 58 · schemastore 52 · heartbeat 27 · schema 17 · options 6` = **472** — zeichengenau die Aufteilung aus §3(b) |
| 6 | `go test -coverpkg=…/sqlexec/... -coverprofile` im gepinnten Toolchain-Container | **0** | **152** Statements: `translate.go` 147 · `statement.go` 3 · `errors.go` 2 — die Aufteilung aus §3(b)/`ADR-0078` exakt |
| 7 | Eigene Unit-Messung (Package-Liste + `-coverpkg` der `coverage`-Stufe nachgebaut), Toolchain-Container | **0** | Gegenstand **1831** Statements, gedeckt **1308** → **71,44 %** (die dokumentierte Lauf-Varianz ±2; §3 nennt 1306) |
| 8 | **Vor-Zahlen** am Vorgängerbaum: eigener Worktree auf `4c7ea8a~1` (= `69d12f4`), dieselben zwei Messungen | **0** | Unit-Gegenstand **1679**, gedeckt **1171** → **69,74 %** · `postgresstorage` **610** |
| 9 | `DB_COVERAGE_THRESHOLD=80 bash tools/harness/db-coverage.sh` | **1** | `db-coverage: FAIL — DB-Adapter-Coverage 73.38% unter Schwelle 80%` |
| 10 | `DB_COVERAGE_THRESHOLD=74` / `=73` / Default | **1** / **0** / **0** | Schwelle am gemessenen Wert scharf: 74 rot, 73 grün, 70 grün |
| 11 | `make coverage-gate THRESHOLD=75` | **2** (make) | `coverage-gate: FAIL — Coverage 71.30% unter Schwelle 75%` — eigene Gegenprobe der Unit-Rampe über dem Ist |
| 12 | **Mutation A** — `DB`-Schnittstelle um eine Methode erweitert, die der Pool nicht trägt | **1** | `seam.go:43:12: … *pgxpool.Pool does not implement DB (missing method SeamProbeNonexistent)` — die Zusicherung ist tragend, nicht dekorativ |
| 13 | **Mutation B** — `ReadExcludedColumns`: letzter Einschluss streicht den Eintrag nicht mehr | **1** | `FAIL: TestReadExcludedColumnsLastInclusionRemovesEntry` (`map[public.feed:[]]`) — die Fake-Seite prüft die Übersetzung wirklich |
| 14 | Struktur-Proben (`git diff`, `git grep`) | — | siehe §3 und §5 |
| 15 | `git status --porcelain` am Ende; `git worktree list` | — | leer; nur der Hauptbaum — kein Lauf-Artefakt stehengeblieben |

Beide Mutationen wurden unmittelbar danach real zurückgenommen (Datei-Kopie
bzw. `git checkout`); die Läufe 12/13 sind die Exit-Codes der mutierten
Bäume, nicht des Ist-Stands.

---

## 2. DoD-Konformität je Kriterium

Gelesen im **heutigen** Text (`HEAD`); der Plan wurde mehrfach geändert.

### Liefer-Punkt 1 — die Naht existiert

- **„Adapter hängen an adapter-eigener Schnittstelle (`Query`/`Exec`/`QueryRow`)
  statt an `*pgxpool.Pool`; der reale Pool erfüllt sie (`var _ sqlexec.DB =
  (*pgxpool.Pool)(nil)`, `…/sqlexec/seam.go`)."** — **bestätigt.**
  `seam.go` deklariert `Executor`/`DB`; die Zusicherung steht in Zeile 42.
  **Mutation A** zeigt, dass sie bricht, sobald die Schnittstelle divergiert
  (Exit 1) — der Kompilier-Beleg ist scharf.
- **„Die Zeilen-Schnittstelle bleibt `pgx.Rows`"** — **bestätigt.**
  `Executor.Query` liefert `pgx.Rows`; `Executor.QueryRow` liefert `pgx.Row`
  (`seam.go:24-26`). Eine `sqlexec.Rows`-Deklaration existiert **nicht** mehr
  (`git grep` leer) — der Plan behauptet sie auch nirgends mehr (§3/§1
  berichtigt, `f5ba274`).
- **„Kein Verhaltens-Change: die reale Verdrahtung geht unverändert durch
  `make test-store` und `make test-replication` (Exit 0)."** — **bestätigt.**
  Beide Läufe Exit **0** (Messungen 2/3), die Adapter-Tests liefen real
  (`postgresack`, `postgresstorage`, `receive` je `ok`). Die **bestehenden**
  DB-Testdateien sind **byte-identisch** (`git diff --name-status` über
  `'*_test.go'` im Slice-Bereich zeigt genau **eine neue** Datei). Der
  Produktions-Diff der sieben Adapterdateien ist mechanisch
  (`pool`→`db`, `fmt.Errorf("%w: %w", …)`→`sqlexec.Classify`,
  `errors.Is(err, pgx.ErrNoRows)`→`sqlexec.IsAbsent`, entfernte Inline-Funktion
  `collectRecords` lebt unverändert in `sqlexec.ReadChanges`) — siehe §5.

### Liefer-Punkt 2 — die reine Logik ist prüfbar

- **„Fehlerklassifikation und Zeilen-Übersetzung sind reine Funktionen und
  netzlos geprüft (Fälle: Erfolg, Fehlerklasse, leeres Ergebnis)."** —
  **bestätigt.** `Classify`/`IsAbsent`/`Statement.fail` sind zustandsfrei;
  `translate_test.go` fährt netzlos (Lauf 6/12/13 im Container `--network none`)
  und deckt Erfolg (`TestReadChangesIssuesQueryAndTranslatesRows`),
  Fehlerklasse (`…ClassifiesQueryFailure`/`…ScanFailure`/`…IterationFailure`)
  und Leerfall (`…EmptyResultIsNoError`) ab. **Mutation B** (anderer Schnitt
  als die vier Mutationen des Reviews) wird gefangen → die Tests prüfen die
  Übersetzung wirklich, nicht nur ihre Existenz.
- **„Die Fake-Seite fährt die Verklebung (Query/Exec-Aufruf, Scan-Schleife,
  Fehlerpfad) — und ist ausdrücklich kein Ersatz der realen DB-Tests."** —
  **bestätigt.** Der Fake zeichnet `Query`/`QueryRow`/`Exec`-Aufrufe auf,
  fährt die Scan-Schleife und den `rows.Err()`-Pfad; die zwei realen
  DB-gestützten Läufe bleiben grün und unverändert (Messungen 2/3).

### Liefer-Punkt 3 — die Wirkung ist benannt, nicht angestrebt

- **„Der Coverage-Effekt wird als Folge dokumentiert — nicht als Zweck:
  `69,70 % → 71,30 %` (`make coverage-gate`, beide Läufe Exit 0; Gegenstand
  `1679 → 1831` Statements — §3)."** — **bestätigt (mit benannter Grenze).**
  Der Nach-Lauf ist real gemessen (Messung 1: `71.30%`, Exit 0; Gegenstand
  1831, Messung 7). Der Vor-Lauf ist **nicht** als Gate-Lauf nachfahrbar (der
  Baum ist vergangen, die damalige Stufe war 65 %); die Vor-Zahl ist aber
  **unabhängig reproduziert** (Messung 8: 1679 Statements, 1171 gedeckt =
  69,74 %; die gedruckte Stufe zeigt `69,7 %` → `69,70 %`). Die Design-Priorität
  ist in §1 ausdrücklich voranstellt („Die Rechtfertigung ist Design, nicht die
  Zahl"); das Gate wurde **nicht** angefasst, um die Zahl zu bewegen.
- **„Die Neu-Bemessung aus `ADR-0077` in der Fassung von `ADR-0078` ist
  umgesetzt: `DB_COVERAGE_THRESHOLD` 75 → 70 (`tools/harness/db-coverage.sh`),
  `THRESHOLD` 65 → 70 (`harness/mk/coverage.mk`), und die Träger-Doku
  (`db-adapter-coverage.md`, `coverage-gate.md`, `harness/README.md` §Sensors,
  `AGENTS.md` §4) nennt die neuen Stufen."** — **bestätigt.**
  `db-coverage.sh:54` = `70`, `coverage.mk:15` = `70`; Endstufen 80 in beiden
  Sensor-Docs; `harness/README.md` §Sensors und `AGENTS.md` §4 nennen die
  Rampe `Einstieg 70 % → Endstufe 80 %`. Kein lebender Träger nennt noch eine
  geltende alte Stufe (Scan über `*.md`/`*.mk`/`*.sh` außerhalb von
  Zitat-/Lauf-Beleg-Zeilen: leer).
- **„Der Transfer-Nachweis liegt vollständig bei — alle drei Belege
  (`ADR-0078` §Entscheidung)."** — **bestätigt**, unabhängig nachgeprüft,
  siehe §3.
- **„`make gates` grün (Exit direkt, ungepiped)."** — **bestätigt**
  (Messung 1, Exit 0).

### Review-Zeile

- **„Review durchgeführt, Report unter `docs/reviews/` liegt vor … 0 HIGH in
  der Fixrunde."** — **bestätigt** (beide Dateien vorhanden; die Fixrunde
  weist 0 HIGH / 0 MEDIUM aus). Der Auftrag schließt eine Nachprüfung des
  Review-Stils aus; nur die **Existenz** und die DoD-Häkchenlage sind geprüft.

### Die fünf Closure-Punkte (noch `[ ]`)

`Closure-Notiz`, `Reconciliation-Register`, `Beobachtungs-Register`,
`Risiko-Ausgänge §6`, `Drei Paarungen` — alle **offen**. Das ist der
**normale** Stand eines Slice in `in-progress/`: diese fünf Punkte sind
Closure-Pflichten, die **nach** der Verifikation und teils **nach** dem
`git mv` liegen (die Paarungen suchen in `done/`). Kein DoD-Mangel — sie
sind der Gegenstand der Planner-Closure, nicht dieses Berichts. §6 trägt
weiter `<bei Closure>` in allen drei Risiken, §7 die Vorlagen-Platzhalter;
auch das ist vor der Closure erwartet.

**Ergebnis: jedes Liefer-Kriterium bestätigt; kein Kriterium nicht
bestätigt und keines, das der DoD-Text behauptet, aber der Code nicht
hergibt.**

---

## 3. Der Transfer-Nachweis (`ADR-0078` §Entscheidung) — unabhängig nachgeprüft

**(a) Ankunft.** Eigene Messung am Vorgängerbaum (Worktree auf `4c7ea8a~1`)
und am Ist-Stand:

| Größe | vor | nach | Quelle |
|---|---|---|---|
| Unit-Gegenstand (Statements) | **1679** | **1831** | Messungen 7/8 |
| `postgresstorage` (Statements) | **610** | **472** | Messungen 5/8 |

`k_ab = 610 − 472 = **138**` · `k_auf = 1831 − 1679 = **152**` ·
`k_auf (152) ≥ k_ab (138)` ✓ · Differenz `**14**`. Die 152 sind **neuer Code**
des neuen Pakets `sqlexec` (`errors.go` 2 · `statement.go` 3 · `translate.go`
147 — Messung 6), nicht verlagerte Menge. `postgresack` (23) und
`replication/receive` (155) sind in **beiden** Ständen gleich (Messungen 4/8).

**(b) Paket-Granularitäts-Diff — nicht die Summe.** Eigene Reproduktion von
`git diff --name-status 4c7ea8a~1..8e9fe4f -- internal/`: ausschließlich
Dateien unter `internal/adapters/driven/postgresstorage/` plus die **fünf
neuen** Dateien unter `postgresstorage/sqlexec/` (`errors.go`, `seam.go`,
`statement.go`, `translate.go`, `translate_test.go`); **kein** anderes Paket
ist berührt (`postgresack`, `replication/receive` unverändert). Der verlagerte
Code ist im abfließenden Gegenstand **vollständig abgegangen**: kein Rest von
`collectRecords`/`pgx.Rows` in den DB-Gegenstands-Paketen außerhalb von
`sqlexec/` (eigener `grep`: leer). Die **aggregierte Summe** wird ausdrücklich
**nicht** als Träger geführt — §3(a) nennt sie nur als **Lesehilfe** und sagt
in §3(b) ausdrücklich, dass sie den Nachweis **nicht** trägt. Die Anforderung
aus `ADR-0078` ist damit eingehalten.

**(c) Kein Verhalten verloren.** `git diff --name-status … -- '*_test.go'`
zeigt genau **eine neue** Datei (`sqlexec/translate_test.go`), **keine**
Löschung (`--diff-filter=D … -- '*_test.go'` leer). Die realen,
dienst-gestützten Läufe des abfließenden Gegenstands sind grün (Messungen 2/3,
Exit 0 — die Träger von `ADR-0078` Beleg (c)). Die vier Adapter-Testdateien
selbst sind unverändert.

**Ergebnis: (a), (b) und (c) tragen; die Summen-Konstanz wird korrekt nicht
als Träger verwendet.**

---

## 4. Ist die neue Schwelle noch eine Schwelle? — eigene Gegenproben

- **DB-Rampe:** `DB_COVERAGE_THRESHOLD=80` → Exit **1**; `=74` → Exit **1**;
  `=73` → Exit **0**; Default `70` → Exit **0**. Die Schwelle ist am
  gemessenen Ist (73,38 %) **scharf** — eine Senkung, die nichts mehr prüft,
  wäre hier nicht eingetreten.
- **Unit-Rampe:** `make coverage-gate THRESHOLD=75` gegen den Ist 71,30 % →
  **make Exit 2**, Skript `FAIL — Coverage 71.30% unter Schwelle 75%`. Die
  Anhebung 65 → 70 ist keine Entwertung; sie liegt am gemessenen Ist (71,3 %)
  und ist ihrerseits scharf (`71,3 ≥ 70`, Abstand 1,3 pp).

---

## 5. ADR-Konformität

- **`ADR-0071` Punkt 5** — die Naht ist zulässig und **design**-begründet,
  „dass der Wert danach steigt, ist Folge, nicht Zweck". Der Slice trägt das
  (§1: „Die Rechtfertigung ist Design, nicht die Zahl"), und die Umsetzung
  trifft die Design-Begründung (schmale Abhängigkeit: `db sqlexec.DB` statt
  `pool *pgxpool.Pool`; Fehlerklassifikation als reine Funktion: `Classify`).
  **Konform.** Der von Punkt 5 aufgerufene Re-Evaluierungs-Trigger (b) ist
  über `ADR-0077`/`ADR-0078` eingelöst.
- **`ADR-0077`** — die zwei Einstiege (DB 75 → 70, Unit 65 → 70) und die
  unveränderten Endstufen 80 sind in den zwei Werkzeug-Trägern und der Doku
  umgesetzt. `AGENTS.md` §3.6 ist über den `Accepted`-ADR-Träger erfüllt.
  **Konform.**
- **`ADR-0078` (Teil-Supersede)** — der **dreiteilige** Transfer-Nachweis ist
  beigebracht (§3), die **Summe** ist als Allein-Beleg **nicht** verwendet,
  der Regressions-Riegel ist in §3(b) wörtlich benannt. Der ADR-Index führt
  den Nachfolge-Zeiger (`ADR-0077 → ADR-0078, teilweise`). **Konform.**
- **`ADR-0041`** — `.a-check.yml` ist im gesamten Slice-Bereich **unverändert**
  (`git diff --name-status` leer); `a-check` meldet **0** Befunde. Die Naht
  liegt innerhalb des driven Adapters. **Konform.**

---

## 6. Plan-vs-Code-Diff (`§3` gegen den Ist-Zustand)

Die §3-Träger-Tabelle ist gegen den committeten Stand geprüft; **keine Zeile
behauptet eine abwesende Struktur**:

| §3-Zeile | Ist-Zustand |
|---|---|
| `sqlexec/**` (`seam.go`, `statement.go`, `errors.go`, `translate.go`) — neu | alle vier Dateien vorhanden; **keine** `Rows`-Deklaration (die entfernte ist nicht mehr zugesagt) |
| die sechs pool-tragenden Adapter-Dateien · `schema.go` — refactor | genau **7** modifizierte Dateien (`store`, `consumerstate`, `heartbeat`, `tableactivation`, `schemastore`, `administrationrequest`, `schema`) — deckungsgleich |
| `translate_test.go` — neu | vorhanden; Fakes (`pgx.Rows`/`pgx.Row`/`Executor`) — deckungsgleich |
| `db-coverage.sh` · `coverage.mk` — update | 70 / 70 gesetzt — deckungsgleich |
| die vier Doku-Träger — update | alle vier geändert und mit den neuen Stufen — deckungsgleich |
| „Nicht angefasst: `queries`, `.a-check.yml`, `spec/**`" | `queries` nicht im Diff; `.a-check.yml` und `spec/**` nicht im Diff (beide Proben leer) — deckungsgleich |
| §3 „Abweichung": `postgresack` **6** / `receive` **18** Symbole, beide unberührt | eigene Zählung: **6** / **18**; beide Pakete nicht im Diff — deckungsgleich |

Die „Berichtigung des Vertrags" in §1 (die Verengung auf vier Methoden
verträgt sich nicht mit der Pool-Zusicherung; benannte Grenze: **10** Methoden
im Fake) beschreibt den Ist-Zustand: `fakeRows` implementiert **zehn** Methoden
(`Next`/`Scan`/`Err`/`Close` + sechs `nil`-Attrappen
`CommandTag`/`FieldDescriptions`/`Values`/`RawValues`/`Conn`/`TypeMap`).

---

## 7. Folge-Slice-Adressen und `AGENTS.md` §3.10

- **`slice-084` (`open/slice-084-postgresack-naht.md`)** und **`slice-085`
  (`open/slice-085-receive-naht.md`)** existieren im Planning-Lifecycle und
  **nehmen die Sendung an**: `slice-084` §1 führt als Gegenstand
  `postgresack` und schließt `receive` (→ `slice-085`) aus; `slice-085` §1
  führt als Gegenstand `receive` und schließt `postgresack` (→ `slice-084`)
  aus. Die Symbolzahlen (6 / 18) stimmen mit §3 überein. Die Folge-Slice-Paarung
  des §3 löst damit auf. **Bestätigt.**
- **`AGENTS.md` §3.10** greift **nicht**: `.github/workflows/**` ist im
  gesamten Slice-Bereich `4c7ea8a~1..HEAD` **unverändert** (`git diff
  --name-status … -- .github/` liefert keine Zeile). Der Slice berührt keinen
  Workflow; die Regel für einen Post-Push-Lauf ist nicht ausgelöst.
  **Bestätigt.**

---

## 8. Was dieser Lauf strukturell nicht prüfen konnte

- **`make test-integration`** (Compose-Rundlauf) — nicht gefahren; braucht den
  vollen Compose-Stack und gehört zum E2E-Tier. Der Slice ändert keine
  Compose-Datei und keine Sicht, aber die reale Container-Verdrahtung der
  Adapter ist so nicht gegengeprüft; ihr Wächter ist der nicht-blockierende
  Workflow.
- **`make image`** (OCI-Build) und **`.github/workflows/e2e.yml`** — nicht
  gefahren; der Post-Push-Lauf ist repo-extern. Da `.github/workflows/**`
  **nicht** geändert ist (§7), ist hier auch keine `§3.10`-Bringschuld offen.
- **Der Vor-Gate-Lauf `make coverage-gate` (69,70 %, Exit 0)** — nicht
  nachfahrbar (Baum und Stufe vergangen); die Vor-Zahl ist über eine eigene
  Statement-/Deckungs-Messung am Vorgängerbaum gestützt (Messungen 7/8), nicht
  über einen reproduzierten Gate-Lauf.

---

## 9. Beobachtungen außerhalb des DoD (nicht blockierend)

- **Zwei Genauigkeits-Formen für dieselbe Zahl.** Die DoD-Zeile (Liefer-Punkt 3)
  nennt die **gedruckten** Stufen-Werte `69,70 % → 71,30 %`; §3 nennt die
  **gerechneten** Werte `69,74 %` (1171/1679) und `71,33 %` (1306/1831). Beide
  sind je für sich korrekt (eigene Messung 8: 1171/1679 = 69,74 %); die
  Differenz ist die Ausgabepräzision (`go tool cover` druckt eine
  Nachkommastelle). Kein Fakt-Fehler — nur zwei Konventionen im selben
  Dokument.
- **`welle-20` §1 nennt weiter `1679` Statements** als Größe der netzlos
  prüfbaren Fläche; sie ist seit dieser Naht `1831`. Die Zeile liegt
  **außerhalb** der von der DoD benannten vier Doku-Träger und wurde vom
  Review (`F-6`, INFO) bereits als benannter, ausgelagerter Punkt geführt. Sie
  ist **kein** DoD-Gegenstand dieses Slice, wartet aber auf eine Adresse
  (Supersede-Zeiger oder Nachzug durch den Planner-/Architect-Zug).
- **`ADR-0071` §Entscheidung Punkt 2** ist arithmetisch mit sich nicht
  konsistent: sie nennt „Dieselbe gedeckte Menge (**1215** Statements) über
  den kleineren Nenner ergibt **69,74 %**"; `1215 / 1679 = 72,37 %`, gemessen
  sind es **1171 / 1679 = 69,74 %** (Messung 8). `ADR-0071` ist `Accepted` und
  unberührt (`AGENTS.md` §3.5); `ADR-0078` berichtigt bereits eine andere Zahl
  derselben ADR. **Nicht** Gegenstand dieses Slice, aber ein Kandidat für einen
  Nachfolge-Zeiger des Architect-Zugs.

---

## Verdikt

**DoD-Konformität: bestätigt.** Jedes Liefer-Kriterium aus §2 ist am heutigen
Text und am committeten Code erfüllt. Die naht-tragende Zusicherung, die
reinen Funktionen, die zwei realen DB-Läufe (Exit 0), die Rampen-Neu-Bemessung
mit ihren Trägern, der **dreiteilige Transfer-Nachweis** (`ADR-0078`) und
`make gates` (Exit 0) sind durch eigene Läufe und eigene Zählungen gedeckt —
kein Beleg ist aus einem fremden Bericht übernommen.

**Entscheidungs-Konformität: bestätigt** für `ADR-0071` Punkt 5, `ADR-0077`,
`ADR-0078` und `ADR-0041` (§5). Die neue Schwelle ist **keine** stille
Entwertung: beide Rampen werden gegen den Ist-Stand **rot** (DB 74 → Exit 1,
Unit 75 → make Exit 2).

**Aus meiner Sicht closure-fähig.** Es fehlt **nichts**, was die Verifikation
tragen müsste; die fünf offenen §2-Kästchen sind die Closure-Pflichten der
Planner-Rolle (Closure-Notiz, Register, Risiko-Ausgänge, drei Paarungen) und
liegen **nach** diesem Schritt. Unabhängig prüfbar geblieben sind nur die
repo-externen bzw. vergangenen Läufe (§8) — keiner davon ist eine offene
DoD-Zusage. Die drei Beobachtungen in §9 sind **nicht** blockierend und
außerhalb des Slice-Gegenstands.

Dieser Bericht ist ein **Lauf-Beleg** (dieser Stand, diese Läufe, dieses
Ergebnis); er ersetzt keinen Review (Modul 10, schon erfolgt) und keine
Validation (Modul 12, hier nicht ausgelöst).
