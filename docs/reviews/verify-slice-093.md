# Verifikationsbericht: slice-093 — 2026-09-16

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den DoD-Vertrag
(`slice-093` §2, LP1–LP3), die im Slice referzierten Entscheidungen
[ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
(Schnittmaß, die fünf dienstgebundenen Funktionen, §Konsequenzen „test-only"),
[ADR-0085](../plan/adr/0085-build-kontext-ausnahme-test-only-zweck.md)
(der Träger der `.dockerignore`-Ausnahme und ihre vier Merkmale),
[ADR-0071](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(Messgegenstand), [ADR-0044](../plan/adr/0044-image-beleg-semantik.md)
(Image-Beleg) sowie die Hard Rules `AGENTS.md` §3.2, §3.5, §3.7, §3.9, §3.11,
§3.12. **Nicht** gegen den Diff als solchen (Reviewer-Aufgabe) und **nicht**
gegen realen Bedarf (Validator, nicht ausgelöst).

**Frischer Kontext.** Der Slice-Plan wurde am Stand `HEAD` vollständig gelesen
(§1–§8), dazu [ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md),
[ADR-0085](../plan/adr/0085-build-kontext-ausnahme-test-only-zweck.md), die zwei
Review-Reports, die sechs Commits, der Zieltests des LP2 und die zwei berührten
Register-Einträge. Implementer-Bericht und Reviews waren **Kontext**, ihre
Zahlen **nicht** übernommen: jede Zahl dieses Berichts stammt aus einem hier
selbst gefahrenen Lauf — einschließlich der Review-Findings, die ich
nachgemessen habe (F-1 an Merkmal (i)–(iv), F-2 an acht eigenen Mutationen,
F-3 an zwei eigenen Mutationen, D-2 an genau den zwei Mutationen, die der
Delta-Review als **grün** gemeldet hatte). Exit-Codes sind je **ungepiped** und
in **eigenem** Schritt gelesen (`AGENTS.md` §3.9); Gate-Lauf und Auswertung
waren getrennt beauftragt. Arbeitsbaum-Kopien, Profile und Logs liegen
**außerhalb** des Baums (Scratch-Verzeichnis); der Baum ist vor und nach jedem
Lauf sauber (`git status --porcelain` leer).

**Gegenstand.** `HEAD` = `e140363`, Zweig `main`, Baum sauber. Der Vorgang sind
**sechs** Commits über `3ed5311`: `014f29c` (Runde 1), `66eb313` (Review),
`4298c4c` (Fixrunde 1), `583a2bf` + `03fc53a` ([ADR-0085](../plan/adr/0085-build-kontext-ausnahme-test-only-zweck.md) und ihr Träger,
Architect-Zug zu F-1), `6b7a5a6` (Delta-Review), `e140363` (Runde 4). Der Slice
liegt in `in-progress/`; der `git mv` nach `done/` ist **nicht** erfolgt.

**Zwei Zahlen zur Schreibweise.** (a) Das Gate-Skript endet im Rot-Fall mit
**Exit 1**; `make` kapselt den Rezept-Fehlschlag zu **Exit 2** — beide stehen
unten getrennt, wo sie anfallen. (b) Prozentwerte sind die **gedruckten**
Zeilen eines konkreten Laufs; wo ich einen Wert **zurückrechne**, steht das
dabei (`AGENTS.md` §3.12 Instanz A).

---

## 1. Eigene Messungen dieses Laufs

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `make gates` (Log in Datei, Exit **danach** aus eigener Datei) | **0** | sechs Checks: `baseline-verify: v6.5.0 OK — 54 Dateien` · `coverage-gate: OK — Coverage 80.00% erfüllt Schwelle 70%` (gedruckte `total:`-Zeile derselben Stufe: `80.0%`) · `d-check: 767 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s)` · `generated-sync: OK` (beide `.pb.go` byte-gleich) · `a-check: gesamt: 0 Befund(e)`; Baum danach leer |
| 2 | `make coverage-gate THRESHOLD=80` (Wellen-Beleg grün) | **0** | gedruckt `total: (statements) 80.0%`; `coverage-gate: OK — Coverage 80.00% erfüllt Schwelle 80%` |
| 3 | `make coverage-gate THRESHOLD=85` (Wellen-Beleg rot) | **2** (make) | Stufe bricht: `coverage-gate: FAIL — Coverage 80.00% unter Schwelle 85%`, Skript-EC **1** (`did not complete successfully: exit code: 1`), `make: *** [harness/mk/coverage.mk:24: coverage-gate] Fehler 1` |
| 4 | `make doc-commits RANGE=3ed5311..e140363` | **0** | 767 Dateien, **0** Befunde — jeder der sechs Commits trägt seine Kennung |
| 5 | `make doc-immutable RANGE=3ed5311..e140363` | **0** | 767 Dateien, **0** Befunde |
| 6 | `make test` (netzlos, `--network none`, `-race ./...`) | **0** | **32 × `ok`, 0 × `FAIL`, 0 × `DATA RACE`** |
| 7 | `go test -race -count=20 ./internal/bootstrap/ ./internal/adapters/driven/telemetry/` (netzlos, gepinnte Race-Image) | **0** | 2 × `ok`, **kein** Flake, **kein** Race in **20** Wiederholungen der neuen Tests |
| 8 | Eigenes Profil der netzlos prüfbaren Fläche am Diff-Stand (Zählbasis der `coverage`-Stufe: `go list`-Liste, `-coverpkg` über den Gegenstand, `-covermode=atomic`, **dedupliziert über die Block-Position**), **6 ×** | 6 × **0** | gedeckt ∈ {**1523** (5 ×), **1522** (1 ×)}, ungedeckt ∈ {**380**, **381**}; **Nenner 1903** (6 ×); gedruckt 6 × `80.0%` |
| 9 | dasselbe am **Parent** `3ed5311` (Arbeitsbaum-Kopie über `git archive`, außerhalb des Repos) | **0** | gedeckt **1493**, ungedeckt **410**, Nenner **1903**; gedruckt `78.5%` |
| 10 | paketweise, dedupliziert (Diff-Stand / Parent) | **0** | `bootstrap` **598/290/308** (Parent **598/262/336**) · `telemetry` **6/6/0** (Parent **6/4/2**); unverändert: `cmd/pg-change-feed` **49** · `streamv1` **10** · `replication/mapper` **7** · `http` **3** · `natsnotify` **2** · `replication/decode` **1** |
| 11 | LP3 — Zuordnung **jeder** ungedeckten Bootstrap-Block-Position zu einer Funktion (Funktionsgrenzen aus `grep -n "^func "`) | **0** | **139** ungedeckte Blöcke, **sämtlich** in `wiring.go`; `Run` **198** · `Diagnose` **73/81** · `Healthcheck` **14/22** · `RegisterConsumer` **13/17** · `AcknowledgeConsumer` **10/14** → **308**; **außerhalb der fünf: 0 Blöcke, 0 Statements** |
| 12 | **13** neue Testfunktionen der vier neuen Dateien im netzlosen `-v`-Lauf | **0** | jede **1 ×** `=== RUN`; **110** × `--- PASS`, **21** × `--- SKIP` (vorbestehende DSN-gestützte Tests), **0** × `--- FAIL` |
| 13 | **7 Mutationsproben** an `tools/schema/nacharbeit-roles.sql` (Arbeitsbaum-Kopie) + Baseline-Kontrolle, `go test -run TestRolloutDatei ./internal/bootstrap/`, netzlos | Kontrolle **0**, Proben je **1** | Kontrolle grün; **7 × rot**, jede mit ihrer benannten Regel (§4) |
| 14 | End-to-End: `docker build --target coverage` (Rezept der `coverage`-Stufe, `COVERAGE_THRESHOLD=70`) über einen Kontext mit der Mutation `GRANT CREATE ON SCHEMA cdc TO cdc_reader` | **1** | `--- FAIL: TestRolloutDateiTraegtDieRechteDerVerdrahtung`, Stufen-Bau bricht ab — die Mutation färbt **das Gate**, nicht nur den Test |
| 15 | **2 Telemetrie-Mutationen** + Kontrolle (`Warn` → `ErrorContext` / → `InfoContext`) | Kontrolle **0**, Proben je **1** | rot mit `level-Feld: ERROR, wollen WARN` bzw. `Log-Zeile ist kein gültiges JSON: unexpected end of JSON input` — die zwei im Kommentar zitierten Meldungen, wörtlich |
| 16 | `docker buildx build` des Runtime-Images + Digest- und Binär-Vergleich (`ADR-0044`) | **0** | Digest **`sha256:84bdca56…17020`** = `harness/image-hash.txt`; Binär-sha256 frisch **`bbbf7135…452c`** = Binär-sha256 des geführten `:dev`-Images (**byte-identisch**) |
| 17 | Umfang: `git diff --name-status 3ed5311..e140363`; `-- 'internal/**' 'cmd/**'` ohne `_test.go`; `-- harness/mk Makefile Dockerfile`; `-- tools/` | **0** | **12** Pfade; die zweite, dritte und vierte Abfrage **leer** → kein Produktcode, kein `THRESHOLD`, keine Datei unter `tools/` |
| 18 | `THRESHOLD` · Kontext-Ableitung · Register · Reviews · §6/§7 | **0** | `harness/mk/coverage.mk:15` = `THRESHOLD ?= 70`; **166** Dateien unter `cmd/`+`internal/` + `go.mod`/`go.sum` + **zwei** `tools/`-Negationen = **170** Kontexteinträge (**6** Negationen insgesamt); `docs/reviews/` trägt **zwei** 093-Reports; §6: **4 von 4** Risiken auf `<…>`; §7: **6 von 6** Zeilen Platzhalter |
| 19 | `make gates` **mit diesem Bericht im Baum**, erneut (Log in Datei, Exit danach aus eigener Datei) | **0** | sechs Checks; `d-check: 768 Datei(en) geprüft, 0 Befund(e)` (eine Datei mehr = dieser Bericht) · `coverage-gate: OK — Coverage 80.00% erfüllt Schwelle 70%` · `commit-traceability: OK — 5 Commit(s)`. Der **vorige** Lauf derselben Anordnung endete **EC 2** mit **einem** Befund — `verify-slice-093.md:32 ADR-0085 id-unlinked`; nach dem Verlinken der Kennung **0** Befunde. Der Befund war dieser Bericht selbst, nicht der Gegenstand |

---

## 2. DoD-Konformität, Kriterium für Kriterium

### Liefer-Punkt 1 — die Tests existieren und sind netzlos grün

| Kriterium (§2) | Befund |
|---|---|
| für den `bootstrap`-Rest und `telemetry` liegen Tests vor | **erfüllt** — **4** neue Dateien (`wiring_rest_internal_test.go`, `config_file_rest_internal_test.go`, `roles_rollout_file_internal_test.go`, `telemetry/slog_levels_internal_test.go`) mit **13** Testfunktionen, dazu **2** erweiterte Bestandsdateien (#17) |
| der Gate-Lauf **fährt sie wirklich** | **erfüllt** — die `coverage`-Stufe bildet ihre Paketliste aus `go list ./internal/... ./cmd/...` und bricht bei rotem Testlauf ab; mein Nachbau derselben Liste, netzlos: `TESTEXIT=0`, **27** × `ok` (#8); alle **13** neuen Funktionen laufen real (#12); zusätzlich `make test` netzlos **Exit 0, 32 × `ok`** (#6) und **20** Wiederholungen ohne Flake (#7) |
| `make gates` ist grün | **erfüllt** — **Exit 0** aus separater Datei gelesen, sechs Checks grün, Baum danach leer (#1) |
| **der Zuwachs wird als Zahl mit ihrem Lauf genannt** | **materiell erfüllt, Träger offen** — gemessen: `bootstrap` **336 → 308** ugedeckt (gedeckt **262 → 290**, **+28**), `telemetry` **2 → 0** (**+2**), zusammen **+30** Statements bei **unverändertem Nenner 1903** (gedeckt **1493 → 1523**, ungedeckt **410 → 380**; #8, #9, #10). Die Zahl steht heute nur in der Commit-Message (git-Historie — laut `ADR-0083` §Geltungsbereich **kein** Doku-Träger) und in den Lauf-Belegen; §7 ist leer (#18). Die Zahl samt Lauf für §7 steht in **§5** dieses Berichts |
| **die erreichte Quote mit ihrem Band** | **materiell erfüllt, Träger offen** — gedruckte Zeile der Stufe `80.0%`; das Band desselben Stands ist **1 Statement** breit: gedeckt **1522–1523** (6 eigene Läufe: 1 × 1522, 5 × 1523; #8), gedruckt an **beiden** Enden `80.0%` (gemessen, nicht zurückgerechnet) |

### Liefer-Punkt 2 — ein Test misst den verdrahteten Gegenstand, nicht seinen eigenen Aufbau

| Kriterium (§2) | Befund |
|---|---|
| an das **reale Artefakt** gebunden, nicht an etwas, das der Test selbst herstellt oder entzieht | **erfüllt** — der Test liest `tools/schema/nacharbeit-roles.sql` über `filepath.Join(filepath.Dir(runtime.Caller(0)), …)` (arbeitsverzeichnis-unabhängig), entfernt Kommentare und parst `GRANT`-/`CREATE ROLE`-Anweisungen; er enthält **kein** `REVOKE`, **kein** ausgeführtes `GRANT` und **keinen** `Exec(`-Aufruf |
| **durch Mutation belegt, nicht durch Lesen** | **erfüllt** — **7** eigene Proben am Artefakt (#13), jede rot: `USAGE` entfernt → Regel (0) · `SELECT` aus dem Heartbeat-Grant für `cdc_admin` entfernt → (1) · `DELETE`-Grant entfernt → (2) · `cdc.retention_blockers` aus dem Reader-Grant gestrichen → (5) · `ALL TABLES IN SCHEMA` angehängt → (6) · `GRANT CREATE ON SCHEMA cdc` angehängt → (6a) · **neuntes Objekt** (`cdc.administration_request` an beide Rollen) → (7); Baseline-Kontrolle grün. **Die zwei Mutationen, die der Delta-Review als grün gemeldet hatte** (`GRANT CREATE`/`GRANT ALL ON SCHEMA cdc TO cdc_reader`), sind **rot** — der Code ist der engeren Ausnahme nachgezogen, wie `e140363` es zusagt |
| die Bindung endet nicht am Test, sondern am **Gate** | **erfüllt** — die Mutation `GRANT CREATE ON SCHEMA cdc TO cdc_reader` lässt den **Bau der `coverage`-Stufe** real umfallen (#14, Bau-EC 1); LP2s Kernaussage hält **end-to-end** |
| „vorher war sie es nirgends" ist **nicht** mehr die Formulierung des Trägers | **erfüllt** — der Testkopf sagt „bliebe **dort** unsichtbar" (die Schwester-Tests) und benennt die reale Ende-zu-Ende-Seite; die **Commit-Message** von `014f29c` trug die weitere Form — sie ist git-Historie, kein Träger (`ADR-0083` §Geltungsbereich) |

### Liefer-Punkt 3 — die unerreichbaren Statements sind benannt

| Kriterium (§2) | Befund |
|---|---|
| die fünf dienstgebundenen Funktionen **einzeln mit Grund** | **erfüllt** — gemessen (#11): `Run` **198/198** · `Diagnose` **73/81** · `Healthcheck` **14/22** · `RegisterConsumer` **13/17** · `AcknowledgeConsumer` **10/14** = **308**; die Zahl und ihre Notation stimmen mit [ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) §Kontext (4a) **Stück für Stück** überein. Der Grund ist der gemessene `Ping`-Riegel (`pgxpool.New` + `pool.Ping`, `receive.NewStream`, `nats.Connect`) |
| **jede weitere unerreichbare Stelle** | **erfüllt (leer erfüllbar)** — im `bootstrap`-Rest gibt es **keine** Stelle außerhalb der fünf: **0** ungedeckte Blöcke außerhalb ihrer Funktionsgrenzen (#11); `telemetry` steht bei **0/6** (#10). Die Aussage ist damit nicht zu weit formuliert; „Rest nicht erreichbar" ohne Namen trifft niemanden |

### Die Closure-Pflichten aus §2

| Kriterium (§2) | Befund |
|---|---|
| `make gates` grün | **erfüllt** (#1) |
| Review durchgeführt, Report unter `docs/reviews/`, **Delta-Review bei Fixrunde** | **nicht erfüllt** — `review-slice-093.md` (`014f29c`) und `review-slice-093-delta.md` (`4298c4c`) liegen vor; der Delta-Review weist mit D-1 (HIGH) und D-2 (MEDIUM) eine **weitere** Fixrunde aus („Rückgabe-Pfeil … nötig"), `e140363` führt sie aus — für diesen Stand existiert **kein** Review-Artefakt → **V-1** |
| Verifikation durchgeführt, Report unter `docs/reviews/verify-slice-093.md` | **erfüllt** mit diesem Bericht; das Häkchen ist offen (#18) |
| Closure-Notiz mit Steering-Loop-Lerneintrag | **nicht erfüllt** — §7 trägt in **6 von 6** Inhaltszeilen Platzhalter (#18); Material liefert §6 dieses Berichts |
| Reconciliation-Register fortgeschrieben *(entfällt …)* | **entfällt nachweislich** — `docs/plan/planning/reconciliation.md` existiert nicht (Greenfield-Bootstrap) |
| Beobachtungs-Register fortgeschrieben — **kein Zähler wird gesetzt** | **nicht erfüllt** — im Vorgang ist **keine** Registerdatei geändert (#17). Die **Substanz** liefert §6 dieses Berichts: `rollen-test-abdeckungsluecken` (Punkt 1 geschlossen, Punkt 2 offen) und `arbeit-ueberholt-stehenden-traeger` (zweite Gelegenheit, Träger berichtigt) |
| Jedes Risiko aus §6 trägt einen Ausgang | **nicht erfüllt** — 4 von 4 stehen auf `<…>` (#18). Material für alle vier liefert §5/§6 dieses Berichts (R1 nicht eingetreten · R2 nicht eingetreten, mit mutierten Gegenproben · R3 mit Gegenbeleg · R4 eingetreten und bedient) |
| Die drei Paarungen sind getragen | **nicht erfüllt** — fällt laut Zeile der `welle-20`-Closure zu |

**Ergebnis:** Die **drei Liefer-Punkte** und die Gate-Zeile **tragen** — gemessen,
nicht gelesen. Offen sind die regulären Closure-Pflichten (Häkchen, §7 samt
Zahlen und Lauf, Register-Entscheidung, vier §6-Ausgänge, Paarungen) und **ein**
Punkt, der nur diese Rolle sieht: **V-1** (die vierte Runde ohne Review).

---

## 3. Entscheidungs-Konformität — hält der Vorgang, was Plan §1/§3 und die ADRs zusagen?

| Zusage | Befund |
|---|---|
| **Plan §1:** „Was dieser Slice liefert: Tests … Er ändert **keinen** Produkt-Code" | **eingehalten** — **12** Pfade, **8** davon `*_test.go`; die Abfrage `git diff … -- 'internal/**' 'cmd/**'` ohne `_test.go` ist **leer** (#17). Keine Naht, kein `git mv` |
| **`ADR-0082` §Konsequenzen „test-only"** — abgelöst durch [ADR-0085](../plan/adr/0085-build-kontext-ausnahme-test-only-zweck.md) Festlegung 2 | **eingehalten in der geschärften Fassung** — „kein Griff in den Messgegenstand" (keine Datei unter `internal/**`/`cmd/**` außer `_test.go`) ist wörtlich erfüllt (#17); die drei Messgrößen (Paketfilter/`-coverpkg`/`-covermode`, Zählung in `tools/coverage-gate.sh`) und `THRESHOLD` sind unberührt (#17, #18) |
| **`ADR-0085` Festlegung 3, Merkmal (i) „gelesen, nicht gebaut"** | **eingehalten** — `tools/schema/nacharbeit-roles.sql` ist Eingabe eines Tests; im Produktionscode kommt sie nur in **Kommentaren** vor (`wiring.go:490`, `:708`), es gibt keine Code-Referenz (#14 zeigt die Wirkung) |
| **`ADR-0085` Merkmal (ii) „genau eine Datei"** | **eingehalten** — der Diff fügt **eine** Negation hinzu (`!tools/schema/nacharbeit-roles.sql`), kein Verzeichnis, kein Muster, keine Wildcard; Kontext-Ableitung **166 + 4 = 170** gegen **169** ohne sie (#18) |
| **`ADR-0085` Merkmal (iii) „der Leser steht dabei"** | **eingehalten** — der `.dockerignore`-Kopf nennt je Eintrag seinen Leser (`tools/coverage-gate.sh` ← die Stufe selbst, `nacharbeit-roles.sql` ← `roles_rollout_file_internal_test.go`) und die Klasse → **mit einem falschen Richtungssatz**, siehe **V-3** |
| **`ADR-0085` Merkmal (iv) „Image unberührt — mit Beleg"** | **eingehalten** — eigener Bau: Digest **`sha256:84bdca56…`** identisch mit dem geführten Beleg `harness/image-hash.txt`; inhalts-entscheidend (`ADR-0044`) der Binär-sha256: **`bbbf7135…452c`** in **beiden** Images, byte-identisch (#16) |
| **`ADR-0085` Festlegung 5 „kein Sensor, keine Messfläche, kein `THRESHOLD`"** | **eingehalten** — `harness/mk/` ist im Vorgang **nicht** berührt, `THRESHOLD ?= 70` steht unverändert (#17, #18); die Sensor-Doku ist Doku, kein Gate-Werkzeug |
| **`ADR-0085` Folgepflicht (Implementer-Zug): `harness/sensors/coverage-gate.md` §Grenze bekommt den Punkt zur Klasse** | **eingehalten** — §Grenze Punkt **6** benennt die vier Merkmale und führt auf `ADR-0085` (`03fc53a`) |
| **`ADR-0085` Folgepflicht (Planner-Zug): die geschärfte Klausel gehört in die Slice-Pläne** | **offen (kein Bestandsträger verletzt)** — die Probe über `*.md` ohne `.harness/` und `docs/reviews/` findet die alte Klausel nur in `ADR-0082` selbst (die zwei abgelösten Zeilen) und in einem Register-`state.md`; **kein** Slice-Plan führt sie. Der Auftrag ist damit zukunftsgerichtet, nicht verletzt |
| **Plan §1:** Anheben von `THRESHOLD` ist Wellen-Closure-Arbeit | **eingehalten** — `THRESHOLD ?= 70`, Datei nicht berührt (#17) |
| **Plan §1:** Cluster **D2** und Cluster **A** = andere Slices | **eingehalten** — `cmd/pg-change-feed` unverändert **49** ungedeckt; `Run` bleibt Teil der fünf *nicht* bewegten Funktionen (#10, #11) |
| **Plan §1:** die Rollen-Rollout-Datei wird **nicht** geändert | **eingehalten** — `git diff … -- tools/` ist **leer** (#17); alle Mutationen liefen auf Arbeitsbaum-Kopien **außerhalb** des Repos |
| **Plan §3:** `welle-20.md` §4 = **„nicht"** | **eingehalten** — die Datei ist im Vorgang nicht enthalten (#17) |
| **Plan §3:** `harness/sensors/coverage-gate.md` — update, **nur falls** eine Zahl driftet **oder** eine Stelle bei Niederschrift eine bewegliche Größe als Ist-Stand führt | **eingehalten — die Bedingung greift jetzt** (anders als in `slice-092`, V-1): die zwei Änderungen sind (a) die durch diese Arbeit **wahr gewordene** Berichtigung des Schwankungs-Satzes (`014f29c`, gedeckt durch §6-Risiko 4 und `ADR-0085` Festlegung 2 — die „überholte Aussage") und (b) §Grenze 6 als **Folgepflicht** derselben ADR (`03fc53a`). Beide sind **benannt**; keine Zahl driftete still |
| **`AGENTS.md` §3.7** in den neuen Kommentaren | **eingehalten** — kein Zustandsfeld, keine Chronik; die Abwesenheits-Formulierung des Erstbefunds (D-1) ist ersetzt. Der **Zähl**-Defekt der neuen Fassung steht als **V-2** |
| **`AGENTS.md` §3.2 / §3.11** | **eingehalten** — kein `//nolint` in den neuen Tests, kein host-lokaler Pfad; `make docs-check` grün über den ganzen Vorgang (#1, #4, #5) |
| **`AGENTS.md` §3.5** (Accepted-ADRs immutable) | **eingehalten** — `ADR-0082` ist **nicht** in-place geändert; die zwei Klauseln sind per Folge-ADR abgelöst, der ADR-Index ist fortgeschrieben (`docs/plan/adr/README.md` Zeile 100) |

---

## 4. Plan-vs-Code-Diff

Verglichen gegen die §3-Liste des Plans am Stand `HEAD`.

| §3-Zeile | Geliefert | Urteil |
|---|---|---|
| `internal/bootstrap/**_test.go` — Test neu/update | **5** Dateien: `wiring_rest_internal_test.go` (neu), `config_file_rest_internal_test.go` (neu), `roles_rollout_file_internal_test.go` (neu), `administration_internal_test.go` (erweitert), `walretention_internal_test.go` (erweitert) | **Plan eingehalten** |
| `internal/adapters/driven/telemetry/**_test.go` — Test neu/update | `slog_levels_internal_test.go` (neu) | **Plan eingehalten** |
| `harness/sensors/coverage-gate.md` — update, nur unter der Bedingung | zweimal angefasst (`014f29c`, `03fc53a`) — beide Anlässe benannt (§3) | **Plan eingehalten** (mit **V-3** als Wortlaut-Rest in der `.dockerignore`, nicht in dieser Datei) |
| `docs/plan/planning/welle-20.md` §4 — **nicht** | nicht angefasst | **Plan eingehalten** |
| — | `.dockerignore` (eine Negation + Kopfkommentar) | **kein Plan-Bruch**: durch **`ADR-0085`** nachträglich getragen (Festlegung 3, vier Merkmale) und in der ADR-Geschichte als ihr Anlass geführt; die Abweichung ist damit nicht mehr undokumentiert |
| — | `docs/plan/adr/0085-*.md` (neu), `docs/plan/adr/README.md`, `docs/reviews/review-slice-093*.md` (neu) | **kein Plan-Bruch**: Übergabe-Artefakte der Architect- und Reviewer-Rolle (Modul 8), keine Liefer-Punkte |

**Was der Diff nicht enthält, obwohl der Plan es nennt:** nichts. Die zwei
Test-Zeilen sind vollständig geliefert; Produktcode, `THRESHOLD`, `tools/schema/`
und die Wellen-Datei sind pfadmäßig unberührt.

---

## 5. Die Zahlen — selbst gemessen (`AGENTS.md` §3.12), und was in §7 gehört

| Zahl | Wo sie heute steht | Mein Lauf | Befund |
|---|---|---|---|
| Zuwachs **+30** Statements | Commit-Message `014f29c` | gedeckt **1493 → 1523**, ungedeckt **410 → 380**, Nenner **1903** in **beiden** Ständen (#8, #9) | **hält, exakt** |
| `bootstrap` **336 → 308** | Commit-Message | paketweise, dedupliziert: **262 → 290** gedeckt, **336 → 308** ungedeckt bei **598** gesamt (#10) | **hält, exakt** |
| `telemetry` **2 → 0** | Commit-Message | **6/6/0** gegen Parent **6/4/2** (#10) | **hält, exakt** |
| Nenner **1903** | Plan §1, `ADR-0082`, Sensor-Dokument | **1903** an beiden Ständen, **6 ×** wiederholt (#8, #9) | **hält** — Zustandsgröße, unverändert |
| **1523/1903 = 80,00 %** | Commit-Message | gedruckt `total: (statements) 80.0%` in **jedem** der drei Gate-Läufe (#1, #2, #8); exakt **80,0315 %** (**zurückgerechnet**, keine gedruckte Zeile); gedeckt **1523** in 5 von 6 Läufen, **1522** in einem (#8) | **hält als gedruckte Zeile**; „80,00 %" ist die `%.2f`-Ausgabe des Gate-Skripts über derselben Zahl, das Band ist **1 Statement** breit |
| **Band** der Quote | Plan §2/LP1 | gedeckt **1522–1523** (1 × 1522, 5 × 1523; #8); gedruckt an **beiden** Enden `80.0%` — am unteren Ende **gemessen**, nicht geschlossen | **hält**, und die Schwankung ist genau der eine Block `wiring.go:1091.4,1092.1` des `runAdministration`-Zweigs |
| **vier Wellen-Belege** (grün bei 80, rot bei 85) | `ADR-0082` §Konsequenzen | `THRESHOLD=80` → **EC 0**, `OK — Coverage 80.00% erfüllt Schwelle 80%` (#2); `THRESHOLD=85` → **FAIL**, Skript-EC 1, make-EC **2** (#3) | **hält, beide Hälften selbst gefahren** — für §7 gehören **beide** hinein, mit ihrer Schwelle |
| **die zwei Träger der Schwankung** | Sensor-Dokument §Zählbasis | 6 eigene Läufe: `wiring.go:991.5,992.13` **count > 0 in 6 von 6**; `wiring.go:1091.4,1092.1` in **1 von 6** auf `count = 0` (#8) | **hält** — „nur noch **einer**" trägt; der Berichtigungs-Satz des Slice ist **gemessen**, nicht geglaubt |

**Was in §7 gehört** (LP1 verlangt es, `AGENTS.md` §3.12 wendet es an):

1. **die Zustands-Zahl des Gegenstands** — **+28** Statements in `bootstrap`
   (**336 → 308**) und **+2** in `telemetry` (**2 → 0**), zusammen **+30**;
   **Nenner 1903 unverändert**. Diese Zahlen hängen am Code-Stand und schwanken
   nicht.
2. **die Lauf-Größe mit ihrem Lauf** — die gedruckte Zeile **`80.0%`** (= `80.00%`
   in der Ausgabe des Gate-Skripts) des Laufs dieses Vorgangs, **mit ihrem
   Band**: gedeckt **1522–1523** von 1903 über **6** eigene Läufe; die
   Schwankung ist der eine Block `internal/bootstrap/wiring.go:1091.4,1092.1`.
3. **die zwei Wellen-Belege mit ihrer Schwelle** — grün bei `THRESHOLD=80`
   (EC 0), rot bei `THRESHOLD=85` (Skript-EC 1 / make-EC 2).
4. **die LP3-Aussage in ihrer gemessenen Form** — „**308** ungedeckte Statements
   liegen **sämtlich** in den fünf dienstgebundenen Funktionen
   (`Run` 198/198, `Diagnose` 73/81, `Healthcheck` 14/22, `RegisterConsumer`
   13/17, `AcknowledgeConsumer` 10/14); außerhalb ihrer: **0**".
5. **die Adresse des Testkopfes einlösen** — `roles_rollout_file_internal_test.go`
   verweist für die Exit-Codes der sieben Mutationen auf **§7**. Bleibt diese
   Angabe dort aus, zeigt der Kommentar auf eine leere Stelle → **V-4**.

### Die Frage der Welle: kippt `79,9 %`?

**Heute nicht — aber der Puffer ist genau ein Statement.** Gemessen: bei
**1522** gedeckt druckt `go tool cover` **`80.0%`** (eigener Lauf, #8) und das
Gate besteht bei `THRESHOLD=80`. Der gedruckte Wert kippt erst bei **≤ 1521**
(**zurückgerechnet**: 79,9264 % → `79.9%`; 1522/1903 = 79,9790 % → `80.0%`),
und das Gate liest genau diese gedruckte Zeile. Das beobachtete Band reicht bis
**1522** — es liegt damit **am Rand**, nicht darüber: ein weiteres verlorenes
Statement (statt der Schwankung) druckt `79.9%` und färbt die 80-%-Schwelle rot.

---

## 6. Die zwei Register-Einträge dieses Slice

| Eintrag | Befund |
|---|---|
| `BEO-PGC/rollen-test-abdeckungsluecken` (**2×**) — **Punkt 1** | **geschlossen, gemessen**: der neue Test liest den **Datei-Inhalt** (kein `REVOKE`, kein `GRANT`, kein `Exec`), und die Mutation am Artefakt färbt ihn **und das Gate** rot (#13, #14). Der alte `REVOKE`/`GRANT`-Test bleibt daneben stehen — sein Gegenstand ist die laufende Instanz; die **Blindheit gegenüber der Rollout-Datei** ist weg |
| derselbe Eintrag — **Punkt 2** | **offen, unberührt**: „Replication-Stream- und ACK-Adapter sind gegen Rollen-Vertauschung nicht testgesichert, weil `tools/harness/run-replication-tests.sh` keine rollenbeschränkten Login-Test-Identitäten bereitstellt". Der Vorgang enthält **keine** Datei aus `tools/harness/` (#17); der Punkt bleibt Register-Sache. **Präzisierung aus F-4:** der Integrations-Tier liest den Heartbeat-Grant real — die geschlossene Lücke ist die **im Gate-Bündel**, nicht „nirgends" |
| `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (**1×**) — **Richtung 1: wird der Satz wirklich falsch?** | **ja, gemessen** — der Satz „Der Träger der Schwankung … ist in **zwei** Blöcken gemessen" war bei Niederschrift wahr (Parent `3ed5311`: der `991.5,992.13`-Block **count = 0**, #9) und ist durch diese Arbeit falsch geworden: der neue Test fährt ihn deterministisch, der Block trug in **6 von 6** eigenen Läufen `count > 0` (#8). Die Berichtigung ist der **zweite** Vorgang dieser Klasse und damit die Schwelle in Sicht — Entscheidung des Lese-Schritts |
| derselbe Eintrag — **Richtung 2: macht die Korrektur Neues falsch?** | **nicht in der Sache** — die vier berichtigten Aussagen halten (`count > 0` in jedem Lauf · „nur noch einer" · die `slice-091`-Herkunft als Marker · die vier `slice-093`-Läufe); der **Nenner/Ableitungs**-Apparat ist zusätzlich selbst nachgemessen (§5). Zwei **neue** Aussage-Defekte stehen daneben und sind nicht aus der Berichtigung entstanden, sondern aus ihrer Nachbarschaft: **V-2** (im Testkommentar) und **V-3** (im `.dockerignore`-Kopf) |

**Beiläufige eigene Messung am Stand `f90c3f4` (slice-091, *außerhalb* dieses
Diffs, 8 Läufe, netzlos):** gedeckt **1469** (7 ×) / **1471** (1 ×); der Block
`991.5,992.13` trug in **genau einem** dieser acht Läufe `count > 0` — die
Hälfte „der Takt-Zweig trug in **einem** dieser Läufe `count > 0`" des
berichtigten Absatzes ist damit **reproduziert**. Das untere Ende **1468** habe
ich in acht Läufen **nicht** gesehen (es braucht beide Blöcke ungedeckt; belegt
ist es in `verify-slice-091.md` #9 mit **einem** Lauf von 47) — der Satz ist
damit nicht widerlegt, sein unteres Ende aber auch nicht von mir bestätigt. Die
Probe gehört nicht in die Zählung dieses Slice.

---

## 7. Findings

### V-1 — Die vierte Runde `e140363` ist von **keiner** Review gedeckt; §2 verlangt für sie ein Delta-Review

- `kategorie`: **MEDIUM**
- `quelle`: §2 des Slice-Plans, Zeile „Review durchgeführt, Report unter
  `docs/reviews/` liegt vor … **Weist der Review eine Fixrunde aus, deckt ein
  Delta-Review sie ab** — `slice-092` hat das dreimal gekostet; die Zeile steht
  hier, damit sie greift, bevor sie jemand herleiten muss" · Präzedenz
  `verify-slice-092.md` V-2 (MEDIUM), `verify-slice-091.md` V-4
- `pfad`: `docs/reviews/` (nur `review-slice-093.md` und
  `review-slice-093-delta.md`) gegen `e140363`
  (`internal/bootstrap/roles_rollout_file_internal_test.go`, +41/−16)
- `befund`: Der Delta-Review beurteilt **`4298c4c`** und weist mit **D-1
  (HIGH)** und **D-2 (MEDIUM)** eine **weitere** Fixrunde aus („**Merge-blockierend: ja**
  … **Rückgabe-Pfeil Reviewer → Implementer: nötig, kurz** — vier Adressen").
  `e140363` führt sie aus — der Commit-Text nennt sie selbst D-1…D-5. Für
  **diesen** Stand existiert **kein** Review-Artefakt: repo-weit nennt **kein**
  Report die Kennung `e140363` (`grep -rl e140363 docs/reviews/` → leer), und
  `git log 6b7a5a6..HEAD` führt **genau einen** Commit. Die Hausform dieses
  Repos für „Fixrunde nach Review" ist ein **eigener Delta-Report** — die
  Präzedenz ist zweifach belegt (`review-slice-089-delta.md`,
  `review-slice-091-delta.md`, `review-slice-092-delta.md`).
- `verifizierbar`: ja — `ls docs/reviews/ | grep 093`; `grep -rln e140363
  docs/reviews/`; `git log --oneline 6b7a5a6..HEAD`; `git show e140363`
- `urteil`: **Substanz geprüft, Artefakt fehlt.** Ich habe die Runde unabhängig
  nachgemessen — sie ist **code-tragend**, nicht bloß Prosa: sie fügt die
  Prüfung **(6a)** hinzu (neue Assertion), und die zwei Mutationen, die der
  Delta-Review als grün gemeldet hatte, sind dadurch **rot** (#13); **kein**
  Produktcode, `THRESHOLD` unberührt. Die Severity bleibt trotzdem **MEDIUM**:
  das Kriterium steht im DoD und ist nicht erfüllt — und der Unterschied zu
  `slice-092` (dort doku-only) macht die Lücke nicht kleiner, sondern größer.
  Zwei Ausgänge sind zulässig: ein **kurzes Delta-Review** auf `e140363` oder
  die **ausdrückliche §7-Zeile**, dass die Runde die vorgeschriebene Korrektur
  wörtlich ausführt, ihre Wirkung in diesem Bericht nachgemessen ist und ein
  Delta deshalb entfällt. Still bleiben darf es nicht.

### V-2 — „Vier Prüfungen" trägt nicht: der Doppelpunkt zählt **drei** auf

- `kategorie`: **LOW**
- `quelle`: `AGENTS.md` §3.12 Instanz A (eine Zahl im Träger trägt ihren
  Ursprung) · §3.7 (der Kommentar schreibt an den, der die Stelle ändert)
- `pfad`: `internal/bootstrap/roles_rollout_file_internal_test.go:176-179`
- `befund`: Der Satz der vierten Runde lautet:

  ```text
  176 // Vier Prüfungen sind **Regeln über den geparsten Grant-Bestand**, keine
  177 // Namenslisten: (6) kein Rollen-Grant über eine ganze Objektklasse,
  178 // (6a) auf einem Schema-Objekt nur `USAGE` und (7) der Leser teilt kein
  179 // Objekt mit einer schreibenden Rolle.
  ```

  Aufgezählt sind **drei** Regeln — (6), (6a), (7) —, angekündigt sind
  **vier**. Die Vorfassung sagte „**Zwei** Prüfungen … (6) und (7)" und war mit
  ihrer Aufzählung deckungsgleich; die Runde hat die Regel (6a) ergänzt und die
  Zahl übers Ziel gezogen. Welche die gemeinte vierte ist, ist aus dem Satz
  **nicht** auflösbar: die Prüfung (0) (Schema-`USAGE` für alle drei Rollen) ist
  nicht genannt, und die Prüfungen (1)–(5) sind Namenslisten, also gerade nicht
  gemeint. Wer die Regeln ändert, zählt vier und findet drei.
- `verifizierbar`: ja — `sed -n '176,184p' internal/bootstrap/roles_rollout_file_internal_test.go`;
  `git show e140363 -- internal/bootstrap/roles_rollout_file_internal_test.go`
- `urteil`: **ein Wort** — „Vier" → „Drei", oder die vierte Prüfung wird
  genannt. Kein Liefer-Defekt, keine Zahl, die eine Messung behauptet; die
  drei genannten Regeln sind jede einzeln mutiert und rot (#13).

### V-3 — Der `.dockerignore`-Kopf sagt die **Richtung** der Ausnahme verkehrt

- `kategorie`: **LOW**
- `quelle`: `ADR-0085` Kontext (3) („Die Größe der Ausnahme — eine Datei, nicht
  ein Baum": Kontext **169 → 170**) · `AGENTS.md` §3.7 (Ist-Zustand) ·
  §3.12 Instanz B
- `pfad`: `.dockerignore:3`
- `befund`: Der von `03fc53a` neu geschriebene Kopf sagt:

  ```text
  3 # wandern). Die zwei Ausnahmen verkleinern den KONTEXT, nicht eine einzelne
  4 # Stage — jede `COPY . .`-Stufe sieht sie, und jede Ausnahme nennt ihren Leser:
  ```

  Die beiden Negationen **vergroßern** den Kontext: sie nehmen zwei Dateien in
  die Menge auf, die `*` ausschließt. Gemessen/abgeleitet: **166** Dateien
  unter `cmd/`+`internal/` + `go.mod` + `go.sum` + **zwei** `tools/`-Negationen
  = **170** Kontexteinträge, ohne die neue Negation **169** (#18); die
  ADR-Ableitung sagt dasselbe. Die tragenden Aussagen des Kopfes
  (kontextweit statt stage-gebunden · jeder Eintrag nennt seinen Leser · das
  Runtime-Image bleibt unberührt) sind wahr — der **Richtungssatz** ist es
  nicht. Die Formulierung ist genau die Stelle, die **`ADR-0085`** als
  Folgepflicht „auf den Ist-Zustand gezogen" haben wollte.
- `verifizierbar`: ja — `sed -n '1,12p' .dockerignore`;
  `git show 03fc53a -- .dockerignore`;
  `git ls-tree -r HEAD --name-only -- cmd internal | wc -l` (**166**);
  `grep -c '^!' .dockerignore` (**6**)
- `urteil`: **eine Wendung** — „wirken auf den KONTEXT" statt „verkleinern den
  KONTEXT". Kein Liefer-Defekt: Merkmal (ii) und (iv) sind davon unberührt und
  beide gemessen (#16, #18).

### V-4 — Die Adresse „§7 der Closure-Notiz" löst erst mit der Closure auf — und ist dort einzulösen

- `kategorie`: **INFO**
- `quelle`: `AGENTS.md` §3.12 Instanz B („ein Kriterium, das erst nach der
  Arbeit belegt werden kann, ist eine **Zusage**") · Delta-Review **D-3**, das
  genau diesen Adressaten verlangt hat
- `pfad`: `internal/bootstrap/roles_rollout_file_internal_test.go:186-187`
  gegen §7 des Slice-Plans
- `befund`: Der Kommentar sagt jetzt „Rot färbende Mutationen (real gefahren;
  die Exit-Codes führt **§7 der Closure-Notiz von `slice-093`**)" und führt
  **sieben** Mutationen auf. §7 trägt heute **6 von 6** Platzhalter (#18) —
  die Adresse zeigt also auf eine noch leere Stelle. Das ist **kein Defekt
  jetzt**, sondern eine **Pflicht der Closure**: erst mit den Exit-Codes (und
  den Zahlen aus §5 dieses Berichts) löst die Adresse auf. Unterbleibt es,
  steht dort dieselbe Form, die D-3 zu Recht als „Beleg-Adressat ohne
  auflösbares Artefakt" geführt hat — nur mit einem gültigen Dateinamen.
- `verifizierbar`: ja — `sed -n '184,194p' internal/bootstrap/roles_rollout_file_internal_test.go`;
  `sed -n '/## 7. Closure-Notiz/,/## 8./p' docs/plan/planning/in-progress/slice-093-coverage-cluster-d2.md`
- `urteil`: **keine Reparatur am Code** — die Einlösung ist §7 selbst. Der
  Bericht führt die sieben Proben samt ihrem Rot-Nachweis deshalb in §2/LP2
  und §1 #13 mit, damit die Closure sie abschreiben kann.

---

## 8. Negativbefunde

- **geprüft, ohne Befund: LP2 bindet wirklich an das reale Artefakt — 7 eigene
  Mutationen, 7 × rot, Baseline grün** (#13). Die Proben sind am **Artefakt**
  gesetzt (nicht am Test), jede nennt ihre Regel in der Fatal-Zeile, und die
  Fälle, die der Delta-Review als **grün** gemeldet hatte (`GRANT CREATE`/
  `GRANT ALL ON SCHEMA cdc TO cdc_reader`), sind jetzt **rot**: `e140363`
  Prüfung (6a). Zeile **`neuntes Objekt`** (`cdc.administration_request` an
  beide Rollen) fällt an Regel (7).
- **geprüft, ohne Befund: die Bindung endet am Gate.** `docker build --target
  coverage` über einen Kontext mit der Mutation `GRANT CREATE ON SCHEMA cdc TO
  cdc_reader` endet **EC 1** mit `--- FAIL: TestRolloutDateiTraegtDieRechteDerVerdrahtung`
  (#14) — die Zusage „im Gate sichtbar" ist damit **end-to-end** gemessen, nicht
  am Testlauf allein.
- **geprüft, ohne Befund: LP3 hält exakt.** **139** ungedeckte Blöcke, **308**
  Statements, **sämtlich** in den fünf Funktionen; die paketweise Summe
  außerhalb ihrer ist **0**. Die fünf Einzelzahlen stimmen **Stück für Stück**
  mit [ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
  §Kontext (4a) überein (#11).
- **geprüft, ohne Befund: alle Zahlen des Implementers halten.** `+30`,
  `336 → 308`, `2 → 0`, Nenner `1903`, `1523` — jede einzeln nachgemessen, jede
  exakt (#8, #9, #10); die Quote **`80.0%`** ist die **gedruckte** Zeile dreier
  eigener Gate-Läufe (#1, #2) und beider Band-Enden (#8).
- **geprüft, ohne Befund: die zwei Wellen-Belege sind beide real.** Grün bei
  `THRESHOLD=80` (EC **0**), rot bei `THRESHOLD=85` (Skript-EC **1**, make-EC
  **2**) — selbst gefahren, je mit eigenem Exit (#2, #3). Für §7 gehören
  **beide** hinein, mit ihrer Schwelle.
- **geprüft, ohne Befund: `make gates` trägt den ganzen Vorgang.** EC **0** aus
  separater Datei gelesen, sechs Checks grün (#1); `make doc-commits` und
  `make doc-immutable` über die volle Slice-Range je **EC 0**, 767 Dateien,
  0 Befunde (#4, #5) — die zwei Dokument-Gates, die dieser Bericht selbst
  mitfahren muss.
- **geprüft, ohne Befund: die neuen Tests sind netzlos und laufen im Gate.**
  **13** neue Funktionen, jede **1 ×** `=== RUN` im netzlosen Lauf, **110**
  × `PASS`, **0 × FAIL** (#12); `make test` mit Race-Detector **Exit 0, 32 ×
  ok** (#6); **20** Wiederholungen der zwei Pakete unter `-race` ohne Flake und
  ohne Race (#7). Die zeitabhängige Stelle (`runWALRetentionCheck`-Test mit
  2-s-Zusage) hat in **20** Läufen nicht gefeuert — das Risiko
  „test-integration-retention-timing-flake" ist für die **neuen** Tests nicht
  eingetreten.
- **geprüft, ohne Befund: `ADR-0085`s Träger hält, soweit er den Gegenstand
  betrifft.** Merkmal (i): die Datei ist Test-Eingabe, im Produktionscode nur
  als Kommentarpfad. (ii): genau **eine** Negation, keine Muster-/Verzeichnis-
  Ausnahme. (iii): jeder Eintrag nennt seinen Leser (Richtungssatz → **V-3**).
  (iv): **eigener** Bau — Digest gleich `harness/image-hash.txt`, Binär-sha256
  über **beide** Images byte-identisch (#16). Der Weg ohne Build-Konfiguration
  (Bind-Mount) ist in `ADR-0085` §Verglichene Alternativen Option D mit
  Trigger (a)/(d) vertagt und **nicht** still gegangen.
- **geprüft, ohne Befund: F-1, F-2 und F-3 des ersten Reviews sind geschlossen.**
  F-1: `ADR-0085` (Accepted, Index-Zeile 100) trägt die Abweichung, die
  `.dockerignore` nennt Klasse und Leser (#18, §3). F-2: fünf grüne Regressionen
  sind rot geworden (#13, #14). F-3: **beide** Telemetrie-Mutationen färben rot
  mit **genau** den zwei zitierten Meldungen (#15) — der Kommentar trägt seinen
  Satz jetzt.
- **geprüft, ohne Befund: D-1, D-3, D-4 und D-5 des Delta-Reviews sind
  geschlossen.** D-1: die Abwesenheits-Wendung ist weg
  (`grep -rnE "ersetzt eine frühere|sie war handverlesen" --include=*.go .` →
  kein Treffer), der Satz steht im Präsens über das, was die Regeln **tun**.
  D-3: „Lauf-Bericht" ist ersetzt; repo-weit kein zweiter Treffer. D-4: die
  Begründung von (6) hat das Subjekt auf das gezogen, was die Messung zeigt
  (die übrigen Zeilen bleiben stehen, aufgehoben ist die **Trennung**), und die
  Nebenaussage „nimmt jede künftige Tabelle mit" ist entfernt — mit `ALL TABLES
  IN SCHEMA` färbt die Probe (6) rot (#13). D-5: die Nicht-Anker-Form ist ohne
  erfundene Slice-Nummer entfernt, `tools/schema/nacharbeit-roles.sql:101`
  bleibt unangetastet (§1-Ausschluss, #17).
- **geprüft, ohne Befund: kein Produktcode, `THRESHOLD` unberührt, die
  ausgeschlossene Datei unberührt.** 12 Pfade; die drei Abfragen
  `internal/**`+`cmd/**` ohne `_test.go`, `harness/mk Makefile Dockerfile` und
  `tools/` sind **leer** (#17); `THRESHOLD ?= 70` steht (#18). Alle Mutationen
  liefen auf Arbeitsbaum-Kopien **außerhalb** des Repos; der Baum ist vor und
  nach jedem Lauf sauber.
- **geprüft, ohne Befund: `BEO-PGC/rollen-test-abdeckungsluecken` Punkt 1 ist
  geschlossen — Punkt 2 nicht.** §6 dieser Prüfung; die Datei liegt unverändert
  bei `evidence/slice-023.md`, `evidence/slice-028.md` (`2×`), und im Vorgang
  ist **keine** Registerdatei geändert (#17) — der Zähler folgt den Dateien,
  die Dateiung der zwei Vorkommen ist Closure-Entscheidung.
- **geprüft, ohne Befund: `AGENTS.md` §3.2, §3.11, §3.5.** Kein `//nolint` in
  den neuen Tests, kein host-lokaler Pfad (dieser Bericht ebenso), `ADR-0082`
  nicht in-place geändert; `make docs-check` grün über den ganzen Vorgang.

---

## 9. Was ich nicht prüfen konnte

- **Die Läufe des Implementers und der Reviews.** Die Commit-Messages nennen
  „Vier Laeufe … konstant bei 1523", der Review „20 Mutationsproben", der
  Delta-Review „18 LP2-Läufe". Beider Läufe sind nicht mehr einsehbar; ich habe
  die **Substanz** unabhängig gemessen (6 eigene Profile, 8 eigene
  Mutationsproben am Artefakt, 2 Telemetrie-Proben, 1 Gate-Mutationsbau) — ein
  **Gegenbeispiel** hat sich nicht gezeigt. Die konkreten Läufe sind als Belege
  nicht überprüfbar.
- **Das untere Ende `1468` des `slice-091`-Bands.** Meine 8 Läufe am Stand
  `f90c3f4` erreichen **1469/1471**; der Zustand `1468` braucht beide Blöcke
  ungedeckt und ist in `verify-slice-091.md` #9 mit **einem** Lauf von **47**
  belegt. Der Satz ist damit nicht widerlegt — aber von mir auch nicht
  bestätigt. **Außerhalb** dieses Diffs.
- **`make test-store`, `make test-replication`, `make test-integration`,
  `make test-notify`.** Kein Gate, kein Bestandteil der Lieferung; die neue
  LP2-Bindung ist netzlos geführt (das ist ihr Punkt). Der Integrations-Tier,
  auf den F-4 verweist, wurde **nicht** gefahren — F-4s Aussage ist damit eine
  Lese-Kette am Skript, keine eigene Messung.
- **Der reale Post-Push-Lauf.** Dieser Vorgang ändert **keinen** Workflow —
  `AGENTS.md` §3.10 greift dem Buchstaben nach nicht. Die Wirkung der neuen
  Tests auf die CI-Laufzeit bleibt bis zum nächsten echten Lauf unbelegt.
- **Die Stabilität über die Rampe.** Ob die Quote bei `THRESHOLD=80` (der
  `welle-20`-Closure) stabil genug ist, ist mit einem Stand nicht entscheidbar —
  mein Befund dazu (§5) ist: **1 Statement** Puffer bis `79,9%`, und das
  beobachtete Band füllt ihn aus.
- **Der Zähler-Stand nach diesem Vorgang.** Ob
  `arbeit-ueberholt-stehenden-traeger` auf **2×** geht, ist eine
  Register-Entscheidung (Modul 6: „Mensch urteilt, Maschine prüft Deckung") —
  dieser Bericht liefert nur die Messung, die ihr zugrunde liegt.

---

## 10. Verdikt

**Die drei Liefer-Punkte aus §2 tragen — gemessen, nicht gelesen.**

- **LP1:** **13** neue Testfunktionen laufen real im netzlosen Lauf (#12,
  `TESTEXIT=0`, **110** × PASS) **und** in der `coverage`-Stufe des Gates; der
  Zuwachs ist **+30** Statements (`bootstrap` **336 → 308**, `telemetry`
  **2 → 0**) bei **unverändertem Nenner 1903**, die Quote **`80.0%`** gedruckt
  in drei eigenen Gate-Läufen, das Band **1522–1523** (**1** Statement) an
  **beiden** Enden `80.0%` — die Zahl mit ihrem Lauf steht in §5.
- **LP2:** der Test liest das **reale Artefakt** und färbt bei **7** eigenen
  Mutationen rot, jede mit ihrer Regel; die zwei Mutationen, die der
  Delta-Review als **grün** gemeldet hatte, sind **rot**, und eine davon lässt
  den **Bau der Gate-Stufe** real umfallen. Das ist der Kern dieses Slice und er
  hält — end-to-end.
- **LP3:** **308** ungedeckte Statements liegen **sämtlich** in den fünf
  dienstgebundenen Funktionen (`Run` 198/198 · `Diagnose` 73/81 · `Healthcheck`
  14/22 · `RegisterConsumer` 13/17 · `AcknowledgeConsumer` 10/14); **außerhalb
  ihrer: 0**.

**Entscheidungs-Konformität: hält.** Kein Produktcode (**12** Pfade, **8**
Testdateien), `THRESHOLD ?= 70` unberührt, die drei Messgrößen unverändert,
`tools/schema/` unberührt, `welle-20.md` unberührt, `ADR-0082` nicht in-place
geändert. Die `.dockerignore`-Ausnahme ist von `ADR-0085` getragen und erfüllt
alle vier Merkmale — Merkmal (iv) mit **eigenem** Bau: Digest gleich dem
geführten Beleg, Binär-sha256 über beide Images **byte-identisch**. Die zwei
Wortlaut-Reste stehen als **V-2** und **V-3** (beide LOW, beide nicht
blockierend).

**Der Plan-vs-Code-Diff ergibt keine unbegründete Abweichung.** Die zwei
Test-Zeilen sind vollständig geliefert; die zwei weiteren Pfade ohne
Test-Endung sind durch `ADR-0085` gedeckt bzw. Rollen-Übergabe-Artefakte. Die
`BEO-PGC/rollen-test-abdeckungsluecken`-Punkt-1-Zusage ist geschlossen, Punkt 2
bleibt offen; der Berichtigungs-Anlass ist als zweiter Vorgang von
`arbeit-ueberholt-stehenden-traeger` **gemessen** wahr.

**`done/`-fähig ist der Slice damit noch nicht — fünf Punkte fehlen:**

1. **V-1** — die vierte Runde `e140363` hat **kein** Review-Artefakt; §2
   verlangt für sie ein **Delta-Review** (Präzedenz dreifach). **Das ist der
   einzige Punkt dieser Liste, der ein DoD-Kriterium direkt verletzt.**
2. **§7** — die Closure-Notiz mit Lerneintrag, darin **die Zahl mit ihrem Lauf**
   und **beide** Wellen-Belege mit ihrer Schwelle (§5) sowie die **sieben**
   Mutations-Exit-Codes, auf die der Testkopf adressiert (**V-4**).
3. **V-2 und V-3** — zwei Wörter: „Vier" → „Drei" (Testkommentar), „verkleinern"
   → „wirken auf" (`.dockerignore`-Kopf). Beide sitzen in Dateien **dieses**
   Vorgangs und gehören **vor** den `git mv`.
4. **Register-Entscheidung** — ob `evidence/slice-093.md` in beide Einträge
   kommt (`rollen-test-abdeckungsluecken`: Punkt 1 belegt, Punkt 2 offen;
   `arbeit-ueberholt-stehenden-traeger`: zweiter Vorgang, Schwelle in Sicht).
5. **Die Häkchen, die vier §6-Ausgänge** — Material für alle vier liefert §5/§6
   dieses Berichts (R1 nicht eingetreten · R2 nicht eingetreten · R3 mit
   Gegenbeleg · R4 eingetreten und bedient); die **Paarungen** gehören der
   `welle-20`-Closure.

Nach 1–5 ist der Slice `done/`-fähig. **Kein Liefer-Defekt, kein rotes Gate:**
die offenen Punkte sind ein fehlendes Übergabe-Artefakt (V-1), zwei Wörter
(V-2, V-3), die Erfüllung einer angelegten Adresse (V-4) und die regulären
Closure-Pflichten — keine gebrochene Zusage.

**Ein Befund für die Welle, über diesen Slice hinaus:** die 80-%-Schwelle steht
mit diesem Stand **genau** auf dem gedruckten Rand — **1523** ist der Bedarf,
das Band reicht bis **1522**, und **erst 1521** druckt `79.9%`. Der Puffer der
`welle-20`-Closure ist damit **ein Statement**, nicht ein Prozentpunkt: jede
weitere Bewegung des `runAdministration`-Zweigs oder des Wächter-Blocks färbt
das Hochschalten bei `THRESHOLD=80` rot. Die `welle-20`-Closure sollte den
Beleg deshalb **mit** seinem Band führen und den Block
`internal/bootstrap/wiring.go:1091.4,1092.1` als seinen Träger benennen.

---

**Beleg-Lage dieses Berichts:** jede Zahl stammt aus einem der Läufe in §1, je
in eigener Werkzeug-Beauftragung gefahren; die drei Gate-Läufe (#1, #2, #3) und
ihre Auswertung waren **zwei** Schritte, ihre Exit-Codes wurden aus separaten
Dateien gelesen, nie durch eine Pipe (`AGENTS.md` §3.9). Die Mutationsproben
liefen auf **Arbeitsbaum-Kopien außerhalb des Repos** (`git archive`), Logs und
Profile liegen außerhalb des Baums; die Repo-Datei `tools/schema/nacharbeit-roles.sql`
wurde **nicht** angefasst (nach jeder Probe gegen das Original geprüft), das
temporäre Image-Tag ist entfernt. Der Baum ist nach diesem Bericht sauber,
**kein** Commit, keine Änderung an Artefakten des Slice, an `THRESHOLD`, an
Produktcode oder an einem Träger außerhalb dieses Berichts.

---

## Nachtrag — Abschluss-Prüfung nach Delta-2, V-2/V-3-Fix und Closure · 2026-09-16

**Rolle:** Verifier (Modul 11), **ein** Durchgang, **frischer** Kontext für diesen
Nachtrag. Auftrag: die zwei nach dem Delta-2 korrigierten Sätze, die Auflösung
von **V-4**, ein letztes Verdikt. **Nicht** gegen realen Bedarf (Validator, nicht
ausgelöst).

**Stand:** `HEAD` = `fc762ed`, Zweig `main`, Baum sauber. Neue Commits seit §1 des
Hauptberichts: `9f9bffc` (**V-2/V-3**), `a41a9de` (**Closure**: §6-Ausgänge, §7,
elf Häkchen, Register, dieser Bericht), `5662f78` (Delta-2-Report), `fc762ed`
(zwei Weiten aus dem Delta-2). Der Slice liegt weiter in `in-progress/`; der
`git mv` nach `done/` ist **nicht** erfolgt. Exit-Codes sind je **ungepiped** und
in eigenem Schritt gelesen; die Mutationsproben liefen auf Arbeitsbaum-Kopien
**außerhalb** des Repos.

### N-1 Eigene Messungen des Nachtrags

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| N-1.1 | `make gates` am Stand `fc762ed` (Log in Datei, Exit **danach** aus eigener Datei) | **0** | sechs Checks: `baseline-verify v6.5.0 OK — 54 Dateien` · `coverage-gate: OK — Coverage 80.00% erfüllt Schwelle 70%` (gedruckte `total:`-Zeile: `80.0%`) · `d-check: 772 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s)` · `generated-sync: OK`; `a-check: gesamt: 0 Befund(e)`; Baum danach leer |
| N-1.2 | die **sieben** im Testkopf und in §7 genannten Mutationen, je einzeln (Containerkopie, `--network none`), plus Baseline | Baseline **0**, Proben je **1** | **7 × rot, jede mit ihrer Regel**: `SELECT` aus dem Heartbeat-Grant (`:215`, (1)) · `DELETE`-Grant für `cdc_admin` (`:224`, (2)) · `DELETE` an `cdc_capture` (`:239`, (3)) · Schema-USAGE entfernt (`:204`, (0)) · `retention_blockers` gestrichen (`:259`, (5)) · `ALL TABLES IN SCHEMA` angehängt (`:273`, (6)) · `GRANT CREATE ON SCHEMA cdc` angehängt (`:289`, (6a)) |
| N-1.3 | **zwei** Schema-`USAGE`-Grants: `… ON SCHEMA public TO cdc_reader` angehängt (die Mutation, die der Delta-2 als grün meldet) | **0** | **grün** — die Ausnahme von (7) überspringt **zwei** Elemente; reproduziert, ohne neuen Fatal |
| N-1.4 | dasselbe Paar, aber mit `CREATE`: `GRANT USAGE, CREATE ON SCHEMA public TO cdc_reader;` angehängt | **1** | rot an `:289`, Regel **(6a)**, Fatal nennt `schema public` — „(6a) hält **jedes** von ihnen auf `USAGE`" trägt auch für das **zweite** Schema-Objekt |
| N-1.5 | `GRANT USAGE ON SCHEMA public TO cdc_capture, cdc_reader;` (Schema **beider** Seiten) | **0** | grün — die Ausnahme greift klassenweit; (6a) hält das Objekt bei `[usage]` |
| N-1.6 | **Regel (6) deaktiviert** (`if false && istKlassenGrant(objekt)`), dazu `GRANT SELECT ON ALL TABLES IN SCHEMA cdc TO cdc_reader;` | **0** | **grün, kein Fatal** — der rote Fall kommt **allein aus (6)**; (7) erreicht ihn nicht (Bestätigung von Delta-2 **D-3**, INFO) |
| N-1.7 | `grep -c '^!' .dockerignore` · Kontext-Ableitung | **0** | **6** Negationen; die zwei Negationen der Ausnahme-Klasse sind benannt (`tools/coverage-gate.sh`, `tools/schema/nacharbeit-roles.sql`); 166 + `go.mod`/`go.sum` + 2 = **170** (ohne die neue Negation **169**) |
| N-1.8 | V-4: Adresse des Testkopfes gegen **§7** (committed in `a41a9de`) | **0** | §7 `:313-322` führt die **sieben** Mutationen in **derselben Reihenfolge** wie der Testkopf und sagt „alle real gefahren, alle **EC = 1**, Kontrolle **EC = 0**" — **deckungsgleich** mit N-1.2, Zeile für Zeile geprüft |
| N-1.9 | Register-Paarung (c): jedes Verzeichnis unter `observations/BEO-PGC/` gegen sein `evidence/` | **0** | **eines** von 63 ohne Beleg-Datei — `architect-verdikt-ablageort-uneinheitlich`, und dessen `state.md` trägt `gestrichen` **mit Begründung** (Modul 6: gestrichen heißt nicht gelöscht) → **kein Defekt**; die vier in §7 genannten Einträge tragen **2 / 4 / 2 / 2** Belege; `open/` ist leer (§7s Folge-Slice-Aussage) |

### N-2 Die zwei neuen Sätze

**`.dockerignore:3-7` (`fc762ed`) — trägt.** Der Kopf benennt das Paar, das er
meint, **namentlich** („Die beiden `!`-Zeilen für `tools/coverage-gate.sh` und
`tools/schema/nacharbeit-roles.sql`"), die Richtung stimmt (**VERGRÖSSERN**;
abgeleitet **169 → 170**, eigene Zählung der Negationen **6**, ohne Vorwärts-Deixis),
und die Singular-Form („**eine** Negation hebt den Ausschluss für **genau diese
Datei** auf") trifft die Klasse, deren Merkmal (ii) **genau eine** Datei je
Ausnahme verlangt. Die tragenden Aussagen (kontextweit statt stage-gebunden ·
jeder Eintrag nennt seinen Leser · Runtime-Image unberührt) sind unverändert und
waren schon in der Vorprüfung gemessen.

**`roles_rollout_file_internal_test.go:299-300` (`fc762ed`) — trägt, mit einer
benannten Rest-Weite.** Der Satz lautet jetzt:

```text
299 	// Ausgenommen ist ein Schema-Objekt: (6a) hält jedes von ihnen auf
300 	// `USAGE` fest, und die Vorbedingung beider Seiten ist kein Objektzugriff.
```

- **Nicht mehr (keine Über-Behauptung):** die falsche Kardinalitäts-Begründung
  („einelementig, weil (6a) …") ist **weg**; der Satz behauptet **keine** Anzahl.
- **Beide Hälften messen sich:** „(6a) hält **jedes** von ihnen auf `USAGE` fest"
  hält auch für das **zweite** Schema-Objekt (N-1.4 rot an (6a), Fatal nennt
  `schema public`) — die Aussage ist also nicht auf `schema cdc` beschränkt; und
  „die Vorbedingung beider Seiten ist kein Objektzugriff" ist für den
  `cdc`-Grant wahr (drei Rollen halten `USAGE`, kein Objektzugriff).
- **Nicht weniger — die Weite:** die Ausnahme im **Code** ist klassenweit (jedes
  Objekt mit Präfix `schema `), und der **Satz** benennt genau diese Klasse
  („ein Schema-Objekt") — die benannte Ausnahme ist damit so weit wie die im
  Code, anders als in den zwei Vorfassungen. Was der Satz **nicht** ausspricht:
  dass ein Schema-Objekt **ohne** die Eigenschaft „Vorbedingung **beider**
  Seiten" mit übersprungen wird — gemessen (N-1.3, N-1.5) bleibt
  `GRANT USAGE ON SCHEMA public TO cdc_reader` (und dieselbe Zeile an beide
  Rollen) **grün**. Das ist **keine** least-privilege-Lücke ((6a) hält jedes
  Schema-Objekt bei `[usage]`; ein Schema-`USAGE` ist kein Objektzugriff), aber
  es ist eine **unbenannte Grenze der Begründung** — festgehalten, damit sie
  nicht still bleibt.

### N-3 V-4 löst auf — die Adresse trägt die Codes, und sie stimmen

§7 führt seit `a41a9de` die sieben Mutationen **mit ihren Exit-Codes** („alle
real gefahren, alle **EC = 1**, Kontrolle **EC = 0**") und die **zwei
Wellen-Belege mit ihrer Schwelle** (`THRESHOLD=80` → EC 0 · `THRESHOLD=85` →
Skript-EC 1 / `make`-EC 2) sowie die Zahlen mit Lauf und Band (**+30**,
`bootstrap` **336 → 308**, `telemetry` **2 → 0**, Nenner **1903**, Quote
**1522–1523 von 1903**, beide Enden gedruckt `80.0%`) — das sind genau die
Angaben, die der Hauptbericht §5 für §7 verlangt hatte. Gegen meine eigenen
Läufe gehalten: **sieben von sieben** Zeilen stimmen (N-1.2 gegen N-1.8), und
die Listen in Testkopf und §7 sind **dieselben sieben in derselben
Reihenfolge**. **V-4 ist damit erfüllt**; die Adresse zeigt auf eine
**committete**, ausgefüllte Sektion.

### N-4 Findings des Nachtrags

#### N-A — Drei Zustands-Zeilen im Beobachtungs-Register tragen eine Zahl, die der abgeleitete Zähler derselben Datei widerlegt

- `kategorie`: **LOW**
- `quelle`: `AGENTS.md` §3.12 Instanz A · Modul 6 („Der Zähler wird **abgeleitet**,
  nicht geführt … Ein gespeicherter Zähler neben einer Belegliste sind zwei
  Quellen für denselben Zustand") · Klasse
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (**7×**)
- `pfad`: `docs/plan/planning/observations/BEO-PGC/arbeit-ueberholt-stehenden-traeger/state.md:1`
  · `…/beleg-befehl-traegt-seinen-satz-nicht/state.md:1` ·
  `…/rollen-test-abdeckungsluecken/state.md:1-6` gegen die `Zähler`-Zeile
  derselben Dateien
- `befund`: drei Stellen, alle in **Registern**, die die Closure `a41a9de`
  angefasst hat:
  - `arbeit-ueberholt-stehenden-traeger` — Kopfzeile „Zustand: offen (**1×**) —
    unter der Schwelle"; die Zähler-Zeile derselben Datei sagt **2×**
    (`evidence/slice-091.md`, `evidence/slice-093.md`). Die Kopfzeile war bei
    ihrer Niederschrift richtig; die Closure hat den Zähler erhöht und die
    Kopfzeile stehen gelassen.
  - `beleg-befehl-traegt-seinen-satz-nicht` — Kopfzeile „Zustand: offen
    (**2×**) — unter der Schwelle, **kein Ausgang zugewiesen**"; die
    Zähler-Zeile sagt **4×** und „**Schwelle erreicht**", der Absatz darunter
    weist den Ausgang dem Lese-Schritt der `welle-20`-Closure zu. Die Kopfzeile
    widerspricht damit **beiden** Hälften des Ist-Zustands (Zahl **und**
    Schwellen-Lage) und ist schon vor diesem Vorgang stale gewesen.
  - `rollen-test-abdeckungsluecken` — die Kopfzeile („Zustand: offen — Ausgang:
    **weiter offen** → ein Test, der den tatsächlichen
    `nacharbeit-roles.sql`-Inhalt liest (**Punkt 1**), sowie … (**Punkt 2**)")
    führt **Punkt 1** weiter als offene Arbeit auf, während derselbe
    `state.md`-Absatz zwei Zeilen später schreibt: „**Punkt 1 ist mit
    `slice-093` geschlossen**". Der *Zustand* des Eintrags (offen, weiter offen)
    bleibt richtig — Punkt 2 trägt ihn —, die **Aufzählung** ist überholt.
- `verifizierbar`: ja — `head -12` je `state.md`;
  `ls …/evidence/ | wc -l` je Eintrag (**2 · 4 · 2**);
  `git show a41a9de -- …/state.md`
- `urteil`: **eine Zeile je Datei**, keine Messung nötig — Zahl bzw. Aufzählung
  auf den abgeleiteten Stand ziehen. **Nicht blockierend:** der Zähler folgt in
  allen drei Fällen den Dateien und ist dort richtig; die veraltete Kopie steht
  daneben. Es ist aber **dieselbe Klasse**, die der Slice in vier Runden
  bekämpft hat, und zwei der drei Stellen hat die Closure selbst erzeugt bzw.
  verschärft.

#### N-B — „Sechs Vorkommen in vier Runden" (§7) ist eine Zahl ohne genannte Basis — und die Reports dieses Slice dokumentieren neun

- `kategorie`: **LOW**
- `quelle`: `AGENTS.md` §3.12 Instanz A (jede Zahl eines Doku-Trägers trägt
  ihren Ursprung, und eine abgeleitete trägt ihre Rechnung) · dieselbe Klasse
  wie **V-2** des Hauptberichts
- `pfad`: `docs/plan/planning/in-progress/slice-093-coverage-cluster-d2.md:284`
- `befund`: Der Satz „… ein Beleg-Adressat löste nicht auf. **Sechs Vorkommen in
  vier Runden**, alle an derselben Art von Stelle." steht am Ende einer
  Aufzählung, die **vier** Beispiele nennt (abwesender Text · Ausnahme im Code
  weiter als im Satz · Begründung mit unzutreffendem Subjekt · Adressat ohne
  Artefakt). Welche Menge die **sechs** sind, sagt der Satz nicht; die Reports
  dieses Vorgangs dokumentieren **neun** Sätze dieser Art — `review-slice-093`
  F-3 · `review-slice-093-delta` D-1, D-2, D-3, D-4 · dieser Bericht V-2, V-3 ·
  `review-slice-093-delta-2` D-1, D-2 (die INFO-Fälle D-5/Δ1 und D-3/D-4/Δ2
  nicht mitgezählt). Eine Lesart, unter der die Sechs aufgeht (Δ1 D-1…D-4 +
  V-2 + V-3), ist möglich — sie ist aber **nicht benannt**, und mit den zwei
  Sätzen des Delta-2 ist sie auf **acht** gewachsen.
- `verifizierbar`: ja — `sed -n '278,286p' <Slice-Plan>`;
  `grep -nE "^### (F|D|V)-[0-9]" docs/reviews/review-slice-093*.md docs/reviews/verify-slice-093.md`
- `urteil`: **ein Zusatz** — die Basis nennen („…, die beiden Delta-Reports und
  dieser Bericht zusammen") oder die Zahl auf den Ist-Stand ziehen. Nicht
  blockierend; die Aussage **hinter** der Zahl („alle an derselben Art von
  Stelle") trägt und ist der eigentliche Lerneintrag.

#### N-C — Die zwei INFO-Zeilen des Delta-2 bleiben unverändert im Baum — und **D-3** ist nachgemessen

- `kategorie`: **INFO**
- `quelle`: `review-slice-093-delta-2.md` D-3/D-4 (beide INFO, „kein
  Rückgabe-Pfeil") · `AGENTS.md` §3.12 Instanz B
- `pfad`: `internal/bootstrap/roles_rollout_file_internal_test.go:267` (D-3) und
  `:264-266` (D-4)
- `befund`: `fc762ed` hat **D-1 und D-2** des Delta-2 aufgenommen (die zwei
  Sätze aus **N-2**); **D-3** („aufgehoben ist die **Trennung, die (7) prüft**")
  und **D-4** („**jedes** zu diesem Zeitpunkt existierende Objekt des Schemas"
  — `ALL TABLES IN SCHEMA` deckt tabellenartige Objekte, nicht Sequenzen)
  stehen unverändert. **D-3 habe ich selbst nachgemessen** (N-1.6): mit
  deaktivierter Regel (6) und angehängtem `GRANT SELECT ON ALL TABLES IN SCHEMA
  cdc TO cdc_reader` bleibt der Test **grün**, kein Fatal — der rote Fall kommt
  allein aus (6), die Attribution an (7) ist damit **widerlegt**.
- `verifizierbar`: ja — `sed -n '263,270p' <Testdatei>`; N-1.6
- `urteil`: **keine Reparatur verlangt** — der Delta-2 führt beide ausdrücklich
  als **INFO** ohne Rückgabe-Pfeil, und beide sind Weiten von Begründungen ohne
  Wirkung auf die Bindung (jede der drei Regeln färbt in allen Proben rot). Sie
  stehen hier, damit sie vor der Closure **nicht still** bleiben; will die
  Planner-Runde sie mitnehmen, ist es je eine Wendung.

### N-5 Verdikt des Nachtrags

**Der Slice ist `done/`-fähig — die Lieferung trägt, die Closure ist
vollständig, und die vier Punkte des Hauptberichts sind erledigt.**

- **V-1 erfüllt:** der zweite Delta-Review deckt `e140363` **und** `9f9bffc`
  (sein Gegenstand sind beide Commits), Verdikt „die Runde trägt", **0 HIGH,
  0 MEDIUM**, kein Rückgabe-Pfeil.
- **V-2 und V-3 behoben und nachgeprüft:** „Drei" statt „Vier" (`:176`) ist mit
  der Aufzählung (6), (6a), (7) deckungsgleich; der `.dockerignore`-Kopf nennt
  das Paar und die Richtung (N-2).
- **V-4 eingelöst:** §7 trägt die sieben Exit-Codes, und sie stimmen mit meinen
  eigenen sieben Läufen **Zeile für Zeile** (N-3).
- **DoD:** **11 von 11** Häkchen; §6 trägt vier Ausgänge (R1 *entfallen*, R2
  *eingetreten und behoben*, R3 *entfallen*, R4 *eingetreten und behoben*); §7
  trägt Zahlen mit Lauf, das Band, beide Wellen-Belege und die Codes; das
  Register ist fortgeschrieben (drei Belege, `rollen-test-abdeckungsluecken`
  **ohne** neuen — der Slice **schließt** dessen Punkt 1); `make gates` am
  Stand `fc762ed` **EC 0** (N-1.1).
- **Entscheidungs-Konformität unverändert:** kein Produktcode, `THRESHOLD ?= 70`,
  `tools/schema/` unberührt; die zwei neuen Commits ändern einen Kommentar, eine
  Konfigurationszeile und Doku.
- **Rest, nicht blockierend:** **N-A** (drei Zustands-Zeilen im Register),
  **N-B** (die Zahl „sechs" in §7), **N-C** (die zwei INFO-Weiten des Delta-2).
  Alle drei sind **Ein-Zeilen-Nachträge** ohne neue Messung; **N-A** und **N-B**
  empfehle ich **vor** dem `git mv` bzw. im selben Zug — die Registerdateien
  wandern nicht mit, der Slice-Plan **schon**.
- **Offen über diesen Slice hinaus** (unverändert): die `welle-20`-Closure misst
  die Rampenstufe bei `THRESHOLD=80` — der Puffer ist **ein Statement** — und
  führt die drei Paarungen; **Cluster A** (`cmd/pg-change-feed`) ist der letzte
  Schnitt der Welle; `open/` ist leer (N-1.9).

**Nicht gefahren:** `make test-store`/`-replication`/`-integration`/`-notify`
(kein Gate, kein Gegenstand dieses Nachtrags), `make image` (kein
Image-Eingriff in den zwei neuen Commits — `Dockerfile` und der Bau-Kontext
sind unberührt) und die Paarungen (der `welle-20`-Closure zugewiesen).

**Beleg-Lage dieses Nachtrags:** jede Zahl stammt aus einem der Läufe N-1.1…N-1.9,
je in eigener Werkzeug-Beauftragung; der Gate-Lauf (N-1.1) und seine Auswertung
waren **zwei** Schritte, sein Exit-Code wurde aus einer separaten Datei gelesen,
nie durch eine Pipe (`AGENTS.md` §3.9). Die Mutationsproben liefen auf
Arbeitsbaum-Kopien **außerhalb** des Repos; `tools/schema/nacharbeit-roles.sql`
und die Testdatei sind im Repo unverändert (nach jeder Probe gegen das Original
geprüft). `make docs-check` über den Stand **mit** diesem Nachtrag: **EC 0**,
**772** Dateien, **0** Befunde. Der Baum trägt nach diesem Nachtrag allein die
Änderung an diesem Bericht — **kein** Commit.
