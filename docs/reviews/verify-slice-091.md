# Verifikationsbericht: slice-091 — 2026-09-16

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den DoD-Vertrag
(`slice-091` §2, LP1–LP3), die im Slice referenzierte Entscheidung
[`ADR-0082`](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
(§Konsequenzen „test-only" und „Proto-Runtime-Interna" · §Kontext (2) die
`±2`-Schwankung · die Cluster-Tabelle) sowie die Hard Rules `AGENTS.md` §3.1,
§3.6, §3.7, §3.9, §3.11, §3.12. **Nicht** gegen den Diff als solchen
(Reviewer-Aufgabe, mit [`review-slice-091.md`](review-slice-091.md)
abgeschlossen) und **nicht** gegen realen Bedarf (Validator, nicht ausgelöst).

**Frischer Kontext:** Diese Sitzung hat den Slice-Plan am Stand `HEAD` gelesen,
dazu `ADR-0082`, den Review-Report und die drei Commits des Vorgangs. Review
und Implementer-Bericht wurden als **Kontext** gelesen, ihre Zahlen **nicht**
übernommen: jede Aussage unten stammt aus einem hier selbst gefahrenen Lauf —
inklusive der Zahlen, die der Review als „nachgemessen, sie halten" führt
(§5) und der Mutationsangaben der Fixrunde (§6, V-Negativbefund unten).
Exit-Codes sind je **ungepiped** und in eigenem Schritt gelesen (`AGENTS.md`
§3.9); Gate-Lauf und Folgehandlung waren getrennt beauftragt. Logs liegen
**außerhalb** des Baums.

**Gegenstand:** `slice-091`, Stand `HEAD` = `f9cd5e4`, Zweig `main`, Baum
sauber (`git status --porcelain` leer — vor **und** nach jedem Lauf). Die drei
Commits des Vorgangs: `f90c3f4` (Tests), `db0c200` (Review-Report), `f9cd5e4`
(Fixrunde). Der Slice liegt in `in-progress/`; der `git mv` nach `done/` ist
**nicht** erfolgt.

**Zwei Zahlen zur Schreibweise.** (a) Das Gate-Skript endet im Rot-Fall mit
**Exit 1**, `make` kapselt den Rezept-Fehlschlag zu **Exit 2**; alle
Rot-Belege dieses Berichts, die nicht über `make` laufen, stehen als **1**
(direkter `go test`-Aufruf). (b) Alle Prozentzahlen dieses Berichts sind die
**gedruckten** Zeilen eines konkreten Laufs; die Zahl ist eine **Lauf-Größe**,
kein Zustand (§3.12 Instanz A, Befund aus §5).

---

## 1. Eigene Messungen dieses Laufs

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| 1 | `make gates` (Log in Datei, Exit **danach** aus eigener Datei) | **0** | sechs Checks: `baseline-verify v6.5.0 OK — 54 Dateien` · `coverage-gate: OK — Coverage 77.20% erfüllt Schwelle 70%` · `d-check: 748 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s)` · `generated-sync: OK` (beide `.pb.go` byte-gleich) · `a-check: gesamt: 0 Befund(e)`; danach `git status --porcelain` leer |
| 2 | `make test` (netzlos, `--network none`, `-race ./...`) | **0** | **32 × `ok`, 0 × `FAIL`**; die vier Pakete des Clusters C darunter: `natsnotify` 1,014 s · `grpc` 1,034 s · `grpc/streamv1` 1,016 s · `http` 1,253 s |
| 3 | Profil der netzlos prüfbaren Fläche am Stand `f9cd5e4` (eigener Lauf, Zählbasis der `coverage`-Stufe) | **0** | `total: (statements) 77.3 %`; dedupliziert über die Block-Position: **432 ungedeckt / 1471 gedeckt / Nenner 1903** |
| 4 | dasselbe am Parent `22cebec` (Arbeitsbaum-Kopie außerhalb des Repos) | **0** | `total: (statements) 74.8 %`; **479 ungedeckt / 1424 gedeckt / 1903** |
| 5 | Cluster C paketweise, beide Stände (#3/Gegenstand #4) | **0** | `http` **29 → 3** · `grpc` **3 → 0** · `streamv1` **25 → 10** · `natsnotify` **5 → 2** — vorher **62** (exakt die Zeile aus `ADR-0082` Kontext (4b)), nachher **15**, Zuwachs **+47**, Nenner unverändert |
| 6 | ungedeckte **Blöcke** des Clusters C am Stand `f9cd5e4`, mit Statement-Zahl | **0** | `http/readchanges.go:181.3,182.1` (1) · `http/retention.go:47.4,49.1` (2) · `natsnotify/notify.go:136.2,137.12` (2) · `streamv1/changestream.pb.go` 12 Blöcke (59, 62, 73, 78, 168, 171, 182, 187, 217, 218, 220, 241: die Summe der **Statements** ist **10**, **zwei** dieser Blöcke tragen **0** (`:62`, `:171`)) · `streamv1/changestream_grpc.pb.go:97.84,97.84` (**0**) |
| 7 | `internal/bootstrap` isoliert, **2 × 8 Läufe** | **0 / 0** | `runWALRetentionCheck` (`internal/bootstrap/wiring.go:981`) **87,5 % in 13 Läufen**, **100,0 % in 3 Läufen**; der instabile Block ist `internal/bootstrap/wiring.go:991.5,992.13` und trägt **2 Statements** |
| 8 | `internal/adapters/driven/grpcstream` isoliert, **8 Läufe** | **0** (8 ×) | `broadcaster.go:96 Publish` **100,0 % in allen acht Läufen** — der Flap liegt **nicht** dort (siehe V-2) |
| 9 | volle Fläche, **47 Läufe** (Serien 5 + 10 + 14 + 16, plus #3 und #1) | **0** (je) | ungedeckt **434** (250 ungedeckte Blöcke, `runWALRetentionCheck` 87,5 %) oder **432** (249 Blöcke, 100,0 %) — **in genau einem Lauf 435** (`77.1 %`, `runWALRetentionCheck` 87,5 %); in den **40** Läufen mit Blockmitschnitt war der `wiring.go`-Block der **einzige** instabile |
| 10 | Umfang des Vorgangs: `git diff --name-only f90c3f4^..f9cd5e4` | **0** | **10 Pfade**: neun `*_test.go` + `docs/reviews/review-slice-091.md`; **kein** Produktcode; `harness/mk/coverage.mk` unberührt (`THRESHOLD ?= 70`) |
| 11 | `make doc-commits RANGE=f90c3f4^..f9cd5e4` | **0** | `d-check: 748 Datei(en) geprüft, 0 Befund(e)` |
| 12 | `make doc-immutable RANGE=f90c3f4^..f9cd5e4` | **0** | 0 Befunde. Derselbe Aufruf **ohne** `RANGE` endet **Exit 2** (`flag needs an argument: --range`) — die MR-Immutabilitäts-Prüfung braucht die Range, sie hat hier kein Objekt (der Vorgang berührt keine `MR-*`-Datei) |
| 13 | `+`-Zeilen der Go-Dateien des Vorgangs: `slice-[0-9]\|welle-[0-9]` · `t.Parallel` · host-lokale Pfadform | **0 / 0 / 0** | 0 · 0 · 0 (ein grobes Pfad-Muster trifft 15 ×, das sind **HTTP-URL-Pfade** in Testaufrufen wie `/consumers/acknowledge` — keine Dateisystem-Pfade) |
| 14 | Fristen: `3 * time.Second` vor dem Slice / nach `f90c3f4` / nach der Fixrunde | **0** | `sse_test.go` 6 → 9 → 9 (die vier **neuen** Fristen sind jetzt `haengeFrist` = 30 s); `grpc/server_test.go` 3 → 4 → 4 (die eine neue Frist ebenso). Die **neun vorbestehenden** Fristen bleiben unangetastet |
| 15 | DoD-Häkchen des Slice-Plans · §7 · §6-Ausgänge | **0** | **0 von 12** Häkchen gesetzt; §7 trägt in **6 von 6** Inhaltszeilen Platzhalter; **4 von 4** §6-Risiken stehen auf `— **Ausgang:** <…>` |
| 16 | `ls docs/plan/planning/reconciliation.md` | **1** | Datei existiert nicht (Greenfield) → das §2-Item „entfällt" ist nachweislich richtig |
| 17 | Beobachtungs-Register: Dateien im Vorgang · Belege der zwei berührten Einträge | **0** | im Vorgang **keine** Registerdatei geändert; `negativtest-ohne-bindung-an-seine-eingabe/evidence/` = **4** (`slice-083/086/087/088`), `test-integration-retention-timing-flake/evidence/` = **2** (`slice-057/090`) |
| 18 | Mutationsproben (Details §6/Negativbefunde) | je **1** | **13 von 13** Proben an **in-diff**-Negativtests färben rot; **1** Probe an einem **nicht** im Diff liegenden Test bleibt grün (V-7) |

---

## 2. DoD-Konformität, Kriterium für Kriterium

**Liefer-Punkt 1 — die Tests existieren und sind netzlos grün.**

| Kriterium (§2) | Befund |
|---|---|
| für die vier Pakete des Clusters C liegen Tests vor | **erfüllt** — neun `*_test.go` über die vier Pakete (#10); `streamv1` erhält mit `changestream_test.go` seine erste Testdatei |
| der Gate-Lauf **fährt sie wirklich** (`make test`, netzlos) | **erfüllt** — `make test` läuft `go test -race ./...` unter `--network none`, **Exit 0**, 32 × `ok`, darunter die vier Pakete (#2). Zusätzlich fährt die `coverage`-Stufe selbst `go test` über die `go list`-Liste (31 Pakete, die vier enthalten) und bricht bei rotem Testlauf ab — `make gates` **Exit 0** mit `77.20 %` (#1) ist damit auch ihr Beleg |
| `make gates` ist grün | **erfüllt** — Exit **0** aus separater Datei gelesen (#1) |
| **der Zuwachs wird als Zahl mit ihrem Lauf genannt** | **materiell erfüllt, Träger offen** — gemessen: **62 → 15** (`http` 29→3, `grpc` 3→0, `streamv1` 25→10, `natsnotify` 5→2), **+47**, Nenner 1903 unverändert (#5). Diese Zahl steht heute **nur** in den Commit-Messages (git-Historie — laut `ADR-0083` §Geltungsbereich ausdrücklich **kein** Doku-Träger) und im Review-Report (Lauf-Beleg); §7 ist leer (#15). §7 muss sie mit ihrem Lauf nennen — die Zahl dafür steht in §5 dieses Berichts |

**Liefer-Punkt 2 — die Negativtests binden ihre Ablehnung an die Eingabe.**

| Kriterium (§2) | Befund |
|---|---|
| geprüft an **mindestens drei** Tests durch Mutation der **Produktionsseite** | **erfüllt** — **13 eigene Proben** an 13 Tests (Negativbefunde unten). Jede Probe schaltet die Ablehnung an der Produktionsseite ab (Decode-Zweig, `401`-Zweig, Interceptor-Prüfung, Sendefehler, `Flusher`-Zweig, RecvMsg-Fehler, Empty-Source-Ablehnung) |
| die Ablehnung hängt am **Eingabewert**, nicht am Fake | **erfüllt** — alle fünf `…UngueltigesJSONEndetMit400` des Diffs, der gRPC-Sendefehler-Test, der SSE-`Flusher`-Fall, der Interceptor-Negativfall, die beiden `401`-Fälle und der `streamv1`-Empfangsfehler färben **rot**, sobald die Produktionsseite nicht mehr ablehnt (Details unten) |
| die **Gegenrichtung** ist mitgemessen | **erfüllt** — Probe „**jeder** Body endet mit 400" (`consumer.go:117`) färbt `TestRemoveConsumerUngueltigesJSONEndetMit400` **rot**: die Gegenprobe des Tests bindet real |
| ein Test, der grün bleibt, bindet nicht | **eine Stelle außerhalb des Diffs** fällt darunter → **V-7** (INFO, benannt, nicht gezählt) |

**Liefer-Punkt 3 — die unerreichbaren Statements sind benannt.**

| Kriterium (§2) | Befund |
|---|---|
| jedes der vier Pakete nennt seine verbleibenden ungedeckten Statements **einzeln** | **erfüllt in der Substanz, Träger teils offen** — `grpc` hat **0** ungedeckte Statements (nichts zu benennen, #5). Die anderen drei sind einzeln erhoben und ihre Gründe halten (§5, V-3): `http` = 3, `natsnotify` = 2, `streamv1` = 10 (#6) |
| „für keine Eingabe erreichbar" (Klasse 1) trägt | **erfüllt, selbst nachgeprüft** — `readchanges.go:181`: `model.NewSourcePosition` scheitert nur bei leerer Quelle oder `offset == 0`; die Quelle ist oberhalb Pflichtfeld (`:139`), `offset >= 1` ist `:176` erzwungen → der Zweig ist für **keine** Eingabe erreichbar. `retention.go:47–49`: `NewRetentionPolicy` hat genau die Bedingung `Nanos < 0`, und `model.NewDuration` (`:40`) hat sie Zeilen vorher ausgeschlossen → ebenso unerreichbar. Kein „Rest nicht erreichbar" ohne Namen |
| „lebender Dienst" (Klasse 2) trägt | **erfüllt** — `natsnotify/notify.go:136–137` ist der Publish-**Erfolgspfad** (`a.log.Debug` + `return nil`); er braucht einen laufenden NATS-Server und wird von `make test-notify` gefahren (`CDC_NATS_TEST_URL`), das den Test real grün sieht |
| die **dritte Klasse** (Proto-Runtime-Interna) trägt, und `ADR-0082` schließt sie wirklich aus | **erfüllt, selbst nachgeprüft** — alle **zehn** Statements einzeln gegen den erzeugten Quelltext gehalten: `Change.String` (`59`), `ProtoReflect` (`73`), `Descriptor` (`78`), `StreamChangesRequest.String` (`168`), `ProtoReflect` (`182`), `Descriptor` (`187`), `rawDescGZIP` (`217`, `218`, `220`), Wiedereintritt von `init` (`241`). `ADR-0082` §Konsequenzen sagt wörtlich: die 25 Statements des erzeugten Pakets seien „teils Proto-Runtime-Interna (`String`, `ProtoMessage`, `Descriptor`, `rawDescGZIP`, `init`) — ein Test, der *sie* abdeckt, prüft erzeugten Code gegen sich selbst", und der Cluster-C-Slice „soll die Getter und die Client-/Handler-Stubs prüfen, **nicht** die Runtime-Interna". Die Klasse ist also **entschieden**, keine Ausrede. Zwei Genauigkeits-Notizen dazu in **V-8** (INFO) |
| die Differenz `62 → ≈52` ist die **benannte** Lücke | **erfüllt, und die Rechnung stimmt fast** — real erreichbar waren **47** (nicht ≈52), benannt sind **15** (#5). Die Zahl steht in der Commit-Message; §7 muss sie tragen (LP1) |

**Die Closure-Pflichten aus §2.**

| Kriterium (§2) | Befund |
|---|---|
| `make gates` grün | **erfüllt** (#1) |
| Review durchgeführt, Report unter `docs/reviews/`, kein Self-Review | **formell erfüllt, Substanz offen** — `review-slice-091.md` liegt vor (0 HIGH, 1 MEDIUM, 1 LOW, 2 INFO) und prüft `f90c3f4`; die **Fixrunde `f9cd5e4` hat keine Review** → **V-4** |
| Verifikation durchgeführt, Report unter `docs/reviews/verify-slice-091.md` | **erfüllt** mit diesem Bericht; das Häkchen ist offen (#15) |
| Closure-Notiz mit Steering-Loop-Lerneintrag | **nicht erfüllt** — §7 trägt in allen sechs Inhaltszeilen Platzhalter (#15) |
| Reconciliation-Register fortgeschrieben *(entfällt …)* | **entfällt nachweislich** (#16) |
| Beobachtungs-Register fortgeschrieben | **nicht erfüllt** — im Vorgang keine Registerdatei (#17). Meine Läufe liefern **kein** drittes Auftreten von `test-integration-retention-timing-flake` (alle Läufe grün, kein Ergebniswechsel bei unverändertem Stand) und **kein** weiteres der Klasse `negativtest-…` aus dem Diff — der Fund V-7 gehört **benannt**, nicht gezählt |
| Jedes Risiko aus §6 trägt einen Ausgang | **nicht erfüllt** — 4 von 4 stehen auf `<…>` (#15) |
| Die drei Paarungen sind getragen | **nicht erfüllt** — fällt laut Zeile der `welle-20`-Closure zu |

**Ergebnis:** Die **drei Liefer-Punkte** und die Gate-Zeile **tragen** — gemessen,
nicht gelesen. Offen sind die regulären Closure-Pflichten (Häkchen, §7,
Register-Zeile, §6-Ausgänge, Paarungen) sowie die vier Punkte, die nur diese
Rolle sieht: **V-1** (die Zahl), **V-3** (der Träger der LP3-Benennung),
**V-4** (die ungeprüfte Fixrunde), **V-7** (die sechste `400`-Stelle).

---

## 3. Entscheidungs-Konformität — hält der Diff, was Plan §1 und `ADR-0082` zusagen?

| Zusage | Befund |
|---|---|
| **kein Produkt-Code** (Plan §1: „Was dieser Slice liefert: Tests … Er ändert **keinen** Produkt-Code, außer einer Naht, die ein Test wirklich braucht") | **eingehalten** — zehn Pfade im Vorgang, neun davon `*_test.go`; **kein** `internal/**`-/`cmd/**`-Go außerhalb von Tests (#10). Es wurde **keine** Naht gezogen |
| **`THRESHOLD` unangetastet** (Plan §1: das Anheben ist Wellen-Closure-Arbeit; „ein Slice, der seine eigene Messlatte mitzieht, könnte sein Grün nicht mehr belegen") | **eingehalten** — `harness/mk/coverage.mk:15` steht auf `THRESHOLD ?= 70`, und die Datei ist im Vorgang **nicht** berührt (#10) |
| **Cluster D und A nicht berührt** (Plan §1: eigene Slices) | **eingehalten** — kein Pfad aus `internal/application/usecase/**`, `internal/domain/model/**`, `internal/adapters/driven/telemetry`, `internal/bootstrap` oder `cmd/pg-change-feed` im Vorgang (#10). Auch Cluster B (`postgresstorage`, `replication`) ist unberührt |
| `ADR-0082` §Konsequenzen **„test-only"** | **eingehalten** — dieselbe Messung; keine Zeile `internal/**`/`cmd/**` außerhalb von Tests ändert sich „um Coverage zu gewinnen" (wörtlich die Folgepflicht der ADR) |
| `ADR-0082` §Konsequenzen: „der Slice … soll die **Getter und die Client-/Handler-Stubs** prüfen, nicht die Runtime-Interna" | **eingehalten** — der `streamv1`-Test prüft Getter (auf dem Nullwert und gegen die Gegenprobe), die drei Fehlerzweige des erzeugten Clients, den `Unimplemented`-Stub und die Handler-Verklebung; die zehn Runtime-Interna bleiben ungedeckt **und** sind als solche benannt (#6, §2/LP3) |
| Plan §3: `welle-20.md` §4 = „**nicht**" (die erreichte Zahl gehört in die Closure-Notiz, nicht in die Welle) | **eingehalten** — `welle-20.md` ist im Vorgang nicht berührt (#10) |
| Plan §3: `harness/sensors/coverage-gate.md` „update, **nur falls** eine Zahl dort gegen die Messung driftet" | **Entscheidung offen** — die Datei ist nicht angefasst (zulässig, ihre Zahlen tragen durchweg ihren Lauf). Zwei Stellen verdienen die Entscheidung trotzdem: die `±2`-Ursache ist **contra-gemessen** (**V-2**) und das `streamv1`-Bild ist gegen den heutigen Stand alt (**V-6**, INFO) |
| Plan §1: die DB-Adapter-Coverage, der Test-Runner und die Cluster D/A bleiben draußen | **eingehalten** — kein Pfad der drei Gegenstände im Vorgang (#10) |

---

## 4. Plan-vs-Code-Diff

Verglichen wurde gegen die §3-Liste des Plans am Stand `HEAD` (der Plan ist in
diesem Vorgang **nicht** nachgezogen worden — der Diff enthält keine
Plan-Änderung, #10; ein Vergleich gegen einen Vorstand ist damit nicht nötig).

| §3-Zeile | Geliefert | Urteil |
|---|---|---|
| `internal/adapters/driving/http/**` — Test neu | `consumer_test.go`, `retention_test.go`, `server_test.go`, `sse_test.go`, `verwaltung_test.go` | **Plan eingehalten** |
| `internal/adapters/driving/grpc/**` — Test neu | `interceptor_test.go`, `server_test.go` | **Plan eingehalten** |
| `internal/adapters/driven/natsnotify/**` — Test neu | `notify_test.go` | **Plan eingehalten** |
| `internal/adapters/driving/grpc/streamv1/**` — Test neu, **falls** dort prüfbarer Code liegt; „der generierte Stub ist **nicht** Gegenstand" | `changestream_test.go` | **Plan eingehalten** — der Test prüft Verhalten, nicht den erzeugten Text; `make generated-sync` bleibt grün (#1), die neue Testdatei ist nicht als Erzeugnis gekennzeichnet |
| `harness/sensors/coverage-gate.md` — „update, nur falls …" | nicht angefasst | **Plan eingehalten** mit benannter Folge-Entscheidung (§3, V-2, V-6) |
| — | `docs/reviews/review-slice-091.md` (neu) | **kein Plan-Bruch**: der Report ist das Übergabe-Artefakt der Reviewer-Rolle (Modul 8), kein Liefer-Punkt |

**Was der Diff nicht enthält, obwohl der Plan es nennt:** nichts. Der Vorgang
ist vollständig durch die vier Testzeilen gedeckt.

---

## 5. Die Zahlen, die in §7 landen werden — selbst gemessen (`AGENTS.md` §3.12)

| Zahl | Wo sie heute steht | Mein Lauf | Befund |
|---|---|---|---|
| Cluster C **62 → 15**, Zuwachs **+47** | Commit `f90c3f4` (Message), Review-Report | **62 → 15**, paketweise **29/25/3/5 → 3/10/0/2** (#5) | **hält, exakt** — und ist **flap-frei**, weil die dokumentierte Schwankung in einem anderen Paket liegt: `internal/bootstrap` trägt in **beiden** Vergleichsläufen **334** ungedeckte Statements (#3, #4) |
| Nenner **1903** | `ADR-0082`, Sensor-Dokument | **1903** in beiden Ständen (#3, #4) | **hält** |
| Gate-Total **77,20 %** | Commit `f90c3f4`: „Gate druckt 77,20 %" | `make gates` druckt **77.20 %** (#1) | **hält** — als **Lauf-Größe** |
| Gate-Total **77,2 % / 77,3 %** | Commit-Message, Review-Report | `77.2 %` (434 ungedeckt), `77.3 %` (432), **einmal 77.1 %** (435) (#9) | **hält, aber mit drei beobachteten Enden** → **V-1** |
| **„+49, nicht +47"** (Gate-Zähler) | Commit `f90c3f4` (Message), Review-Report | **erklärt, aber nicht die einzige zulässige Zahl** — die beobachteten Enden des Diff-Stands sind **432** (Schwankung „gedeckt") und **434** (Schwankung „ungedeckt"); der Parent ist von mir **einmal** gemessen — **479** (Schwankung „gedeckt", #4) —, sein anderes Ende **481** ist nach **derselben** Schwankung **abgeleitet**, nicht gemessen. Die Paarung **481 → 432** ergibt **+49**, die **gemessene** Paarung **479 → 432** ergibt **+47**, die Paarung **479 → 434** **+45** (#3, #4, #9) | **Die `±2`-Schwankung erklärt die Differenz +49 vs. +47 vollständig** — sie verschiebt **denselben** Quelltext um 2 Statements, und **beide Enden sind am Diff-Stand in meinen Läufen real aufgetreten** (`432`/`434`, #7, #9); am Parent ist das zweite Ende aus derselben Schwankung abgeleitet. Nur: „+49" ist damit **eine** zulässige Zahl aus einem Band 45–49, nicht **der** Zuwachs |
| `±2`-Quelle: `runWALRetentionCheck` | `ADR-0082` §Kontext (2) | **bestätigt** — 13 × 87,5 %, 3 × 100,0 % in 16 isolierten Läufen; instabiler Block `internal/bootstrap/wiring.go:991.5,992.13`, **2 Statements** (#7) | **hält** — und ist der **einzige** instabile Block in 40 Läufen mit Blockmitschnitt (#9) |
| `±2`-Quelle: `grpcstream/broadcaster.go` `Publish` | `harness/sensors/coverage-gate.md` §Zählbasis | **100,0 % in 8 von 8** isolierten Läufen und in **47 von 47** Läufen der vollen Fläche (#8, #9) | **hält nicht** → **V-2** |
| `streamv1` **86 Statements, 61 gedeckt (70,9 %, Lauf slice-089)** | `harness/sensors/coverage-gate.md` §Grenze 1 | 86 Statements unverändert, **76 gedeckt** (**88,4 %**) am heutigen Stand (#5) | **als Lauf-Beleg zulässig** (die Zeile nennt ihren Lauf), als Bild des Pakets nicht mehr aktuell → **V-6** |
| Register-Stand **4×** / **2×** | §8 des Plans | `negativtest-…` **4** Belege (`slice-083/086/087/088`), `test-integration-retention-timing-flake` **2** (`slice-057/090`) (#17) | **hält** |

**Was in §7 gehört.** Die Aussage über den **Gegenstand** ist die
Statement-Zahl (Nenner 1903, Cluster C 15, Zuwachs **+47**) — sie hängt am
Code-Stand und schwankt nicht. Die Aussage über den **Lauf** ist die gedruckte
Gate-Zeile — sie gehört mit ihrem Lauf genannt (der Lauf dieses Berichts:
`make gates` am Stand `f9cd5e4`, **77.20 %**) und mit dem Hinweis, dass
dieselbe Messung am selben Stand in meinen 47 Läufen **77,1 / 77,2 / 77,3 %**
ergab. Wird der Gate-**Zähler**-Zuwachs genannt, trägt er den Hinweis auf die
`±2`-Schwankung und deren gemessene Quelle (`runWALRetentionCheck`, **nicht**
`broadcaster.go` — V-2).

---

## 6. Findings

### V-1 — Die Gate-Prozentzahl hat mindestens **drei** beobachtete Enden; die dokumentierte `±2`-Schwankung erklärt zwei davon

- `kategorie`: MEDIUM
- `pfad`: `harness/mk/coverage.mk:15` (`THRESHOLD`) · gedruckte Zeile der
  `coverage`-Stufe · `internal/bootstrap/wiring.go:981`, dort der Block
  `:991.5,992.13` · die Zahlen in `f90c3f4` (Commit-Message) und im
  Review-Report
- `befund`: Ich habe die netzlos prüfbare Fläche **47 ×** über denselben
  Quelltext gemessen (Zählbasis der `coverage`-Stufe: eigene Paketliste aus
  `go list`, `-coverpkg` über den Gegenstand, `-covermode=atomic`,
  dedupliziert über die Block-Position). Ergebnis: **432** ungedeckt
  (249 ungedeckte Blöcke, `runWALRetentionCheck` 100,0 %, gedruckt
  `77.3 %`), **434** (250 Blöcke, 87,5 %, `77.2 %`) — und in **einem** Lauf
  **435** (gedruckt `77.1 %`) bei ebenfalls 87,5 % an derselben Funktion. Die
  beiden ersten Enden sind genau die dokumentierten; das dritte liegt **eins
  unter** dem Band, das `±2` aufspannt. In den **40** Läufen mit
  Blockmitschnitt war `internal/bootstrap/wiring.go:991.5,992.13` der
  **einzige** instabile Block — die Ursache des dritten Endes habe ich damit
  **nicht** lokalisiert (eine zweite, seltene Flap-Quelle im Umfang von 1
  Statement, oder ein Messartefakt; die Beobachtung hat sich in 40 weiteren
  Läufen nicht wiederholt).
- `verifizierbar`: ja — dieselbe Stufe nachfahren; die Einzelläufe über
  `go test -count=1 -coverpkg="$(go list ./internal/... ./cmd/... | grep -vE
  '(^|/)(postgresstorage|postgresack|replication/receive)$' | tr '\n' ',')"
  -coverprofile=… -covermode=atomic` und die letzten beiden Zeilen von
  `go tool cover -func` bzw. die Block-Auswertung des Profils
- `urteil`: **vor der Closure zu adressieren, nicht zu beheben.** Die
  **Lieferung** ist davon unberührt: der Zuwachs **+47** ist flap-frei, weil
  beide Vergleichsläufe denselben `internal/bootstrap`-Stand tragen (#5). Die
  Folge betrifft **§7** (Zahl mit Lauf nennen, bevorzugt die Statement-Zahl;
  siehe §5) und die `welle-20`-Closure: die Rampe nach 80 % hat laut
  `ADR-0082` Kontext (5) **13 Statements** Abstand — dort wird eine
  schwankende Zahl zur schwankenden Gate-Entscheidung. Die Beobachtung ist
  für Tests unsichtbar (alle 47 Läufe grün) und für ein diff-skopiertes
  Review nur über eine Wiederholungsmessung sichtbar: die Klasse gehört
  **benannt** in das Beobachtungs-Register, mit ihren zwei Adressen
  (`zahl-in-traeger-driftet-gegen-die-messung`, 6×, verkörpert in
  `AGENTS.md` §3.12 — hier die Nuance „gemessen, aber mit drittem Ende";
  `test-integration-retention-timing-flake`, 2×, derselbe Test-Gegenstand).
  Ob das ein Auftreten oder eine zweite Beobachtung ist, entscheidet das
  Register-Urteil (Modul 6: „Mensch urteilt, Maschine prüft Deckung").

### V-2 — Der Zahlen-Träger nennt als Ursache der `±2`-Schwankung eine Funktion, die nicht schwankt

- `kategorie`: LOW
- `pfad`: `harness/sensors/coverage-gate.md` §Zählbasis (der Klammerzusatz zur
  Schwankung) gegen `internal/adapters/driven/grpcstream/broadcaster.go:96`
  (`Publish`) und `internal/bootstrap/wiring.go:981`
- `befund`: Der Träger führt die Schwankung auf „ein `ctx`-abhängiger Pfad in
  `internal/adapters/driven/grpcstream/broadcaster.go` — `Publish` mit bereits
  beendetem `ctx`" zurück. Gemessen: `Publish` steht in **8 von 8** isolierten
  Läufen des Pakets und in **47 von 47** Läufen der vollen Fläche auf
  **100,0 %**; dagegen schwankt `runWALRetentionCheck`
  (`internal/bootstrap/wiring.go:981`) real (13 × 87,5 % / 3 × 100,0 %) — und
  `ADR-0082` §Kontext (2) hat an **zwei** Läufen über denselben Quelltext
  gemessen, dass sich die Läufe in **genau einer** Funktion unterscheiden,
  eben dieser. Die Größe `±2` stimmt, die **Begründung** trägt ihren Satz
  nicht: sie nennt keinen Beleg-Anker, und der gemessene Träger ist ein
  anderer (§3.12 Instanz B).
- `verifizierbar`: ja — `go test -count=1 -coverpkg=… -coverprofile=… …` je
  Paket, achtmal, dann `go tool cover -func` und die `Publish`- bzw.
  `runWALRetentionCheck`-Zeile vergleichen
- `urteil`: **kein Reparatur-Ort in diesem Diff** (die Datei ist im Vorgang
  zu Recht unberührt, Plan §3 „nur falls eine Zahl driftet"; die Zahl selbst
  driftet nicht). Der Punkt ist ein **Adress-Fund für den Planner**: wer §7
  schreibt oder die Rampe hochschaltet, soll die Schwankung der **gemessenen**
  Funktion zuordnen — sonst wandert die falsche Ursache in die nächste
  Fassung. Ein Ein-Satz-Nachzug im Sensor-Dokument (oder eine Zitat-Korrektur
  über den Weg, den `AGENTS.md` §3.5 für `Accepted`-Dokumente freigibt) genügt.

### V-3 — Die LP3-Benennung für zwei der vier Pakete lebt nur in Lauf-Belegen

- `kategorie`: LOW
- `pfad`: `internal/adapters/driving/http/readchanges.go:181`,
  `internal/adapters/driving/http/retention.go:47–49`,
  `internal/adapters/driven/natsnotify/notify.go:136–137` gegen §7 des
  Slice-Plans (Platzhalter) und `docs/reviews/review-slice-091.md`
- `befund`: LP3 verlangt, dass jedes Paket seine verbleibenden ungedeckten
  Statements **einzeln mit Grund** nennt. Für `streamv1` steht das **im
  Träger** (Kopfkommentar von `changestream_test.go`, mit `ADR-0082`
  §Konsequenzen als Grund). Für `driving/http` (drei Statements) und
  `driven/natsnotify` (zwei) steht die Benennung heute **nur** im
  Review-Report und in der Commit-Message — beides Lauf-Belege, von denen der
  Report „über Läufe hinweg nicht gelesen" wird (Modul 10) und die
  Commit-Historie nach `ADR-0083` §Geltungsbereich **kein** Doku-Träger ist.
  `driving/grpc` hat nichts zu benennen (0 ungedeckt).
- `verifizierbar`: ja — Volltextsuche nach den drei Stellen über `docs/` und
  `internal/`; sie treffen außerhalb des Review-Reports nichts
- `urteil`: **§7 nachziehen.** Die Substanz hält (die zwei Klassen sind real,
  §2/LP3); es fehlt der Träger, nicht der Inhalt. Die Form: je Paket die
  Stelle mit Datei **und** Zeile und den Grund — für `http` „für keine Eingabe
  erreichbar (oberhalb bereits erzwungen)", für `natsnotify` „lebender Dienst
  (`make test-notify`)".

### V-4 — Die Fixrunde `f9cd5e4` ist von **keiner** Review gedeckt

- `kategorie`: MEDIUM
- `pfad`: `docs/reviews/` (nur `review-slice-091.md`) gegen
  `internal/adapters/driven/natsnotify/notify_test.go`,
  `internal/adapters/driving/grpc/interceptor_test.go`,
  `internal/adapters/driving/grpc/server_test.go`,
  `internal/adapters/driving/http/sse_test.go`
- `befund`: Der Review-Report beurteilt ausdrücklich den Stand `f90c3f4` und
  schließt mit „**Merge-blockierend: ja** — 1 MEDIUM (F-1)" und
  „**Rückgabe-Pfeil Reviewer → Implementer: nötig** (F-1, F-2)". Danach hat
  die Fixrunde **vier** Testdateien geändert (64 Zeilen hinzu, 23 entfernt),
  darunter die **Semantik** eines Tests (F-1: der Test prüft jetzt ein
  zweites Kettenglied; das ist eine neue Zusage, keine Textkorrektur) und die
  Form von vier Fristen (F-3: 3 s → 30 s). Für diesen Stand existiert **kein**
  Review-Artefakt: `ls docs/reviews/` kennt nur `review-slice-091.md`. Der
  Haus-Präzedenzfall `slice-089` hatte für seine Fixrunde einen eigenen
  Delta-Report (`review-slice-089-delta.md`). Die **Substanz** der drei
  Fixrunden-Angaben habe ich unabhängig bestätigt (§6/Negativbefunde: eine
  Probe rot, zwei grün, die zwei natsnotify-Proben rot) — die Klassen
  §3.7-Kommentarform und die **neue** Zusage in `natsnotify` hat aber kein
  Reviewer gesehen (Modul 8: „wer geschrieben hat, reviewt nicht").
- `verifizierbar`: ja — `ls docs/reviews/ | grep 091`;
  `git diff --name-status f9cd5e4^..f9cd5e4`
- `urteil`: **eine Entscheidung ist vor der Closure zu treffen.** Entweder ein
  Delta-Review auf `f9cd5e4` (Präzedenz `slice-089`) oder eine in §7
  ausdrücklich festgehaltene Entscheidung, dass die Fixrunde ohne Delta
  schließt und warum. Beides ist Planner-/Reviewer-Sache; still bleiben darf
  es nicht, weil der Report selbst die Rückgabe als nötig bezeichnet hat.

### V-5 — Die DoD-Häkchen, die §7-Notiz und die vier §6-Ausgänge sind offen

- `kategorie`: LOW
- `pfad`: `docs/plan/planning/in-progress/slice-091-coverage-cluster-c.md`
  §2 (12 × `[ ]`), §6 (4 × `<…>`), §7 (6 × `<…>`)
- `befund`: **0 von 12** Häkchen gesetzt, obwohl LP1 (Test-only-Umfang, Gate
  grün, vier Pakete `ok`), LP2 (13 rote Mutationsproben) und `make gates`
  (Exit 0) materiell belegt sind — die Belege stehen in diesem Bericht. Die
  verkörperte Regel `BEO-PGC/dod-checkbox-nachzug` verlangt, dass eine
  Fixrunde, die einen offenen DoD-Punkt auflöst, ihr Häkchen **im
  Fixrunden-Commit** setzt; `f9cd5e4` hat das nicht getan, `f90c3f4` für die
  drei Liefer-Punkte ebenso wenig. Dazu: **vier** §6-Risiken ohne Ausgang und
  §7 vollständig Platzhalter.
- `verifizierbar`: ja — `grep -n '^ *- \[' …/slice-091-coverage-cluster-c.md`;
  `grep -n 'Ausgang' …/slice-091-coverage-cluster-c.md`
- `urteil`: **Planner-Closure-Arbeit.** Kein Liefer-Defekt; die Häkchen
  tragen die Belege dieses Berichts. Für die §6-Ausgänge liefert dieser Lauf
  Material: Risiko 1 (die „≈52") ist **eingetreten** (real 47) und benannt,
  Risiko 2 (Zeitabhängigkeit) hat mit V-1 seinen Beleg und seine Grenze,
  Risiko 3 (`negativtest`-Klasse) ist im Diff **nicht** eingetreten (13 Proben
  rot) — mit dem benannten Fund außerhalb des Diffs (V-7), Risiko 4
  (Coverage-Theater) trägt mit den Mutationsproben seinen Gegenbeleg.

### V-6 — Das `streamv1`-Bild im Zahlen-Träger ist gegen den heutigen Stand alt

- `kategorie`: INFO
- `pfad`: `harness/sensors/coverage-gate.md` §Grenze 1 (die Aufzählung der
  fünf Pakete ohne Testdatei)
- `befund`: Dort steht für `internal/adapters/driving/grpc/streamv1`
  „**86 Statements, 61 gedeckt** (**70,9 %**, abgeleitet aus 61/86; Lauf
  `slice-089`)". Gemessen am Stand `f9cd5e4`: Nenner unverändert **86**,
  gedeckt **76** (**88,4 %**), ungedeckt **10** (#5, #6). Die Zeile nennt
  ihren Lauf und ist damit **kein** Drift im Sinn der Klasse (die Zahl ist ein
  Lauf-Beleg) — sie ist als Beschreibung des Pakets nur nicht mehr aktuell,
  und sie steht in der Datei, die der Plan in §3 als „update, nur falls eine
  Zahl dort gegen die Messung driftet" führt.
- `verifizierbar`: ja — Block-Auswertung des eigenen Profils für das Paket;
  die Zeile im Träger
- `urteil`: **Entscheidung bei der Closure, kein Blocker.** Auffrischen (mit
  neuem Lauf-Beleg) oder ausdrücklich stehen lassen — beides ist vertretbar;
  unentschieden stehen lassen wäre die dritte, unbenannte Option.

### V-7 — Der sechste `…UngueltigesJSONEndetMit400` bindet seine Ablehnung **nicht** — er liegt außerhalb dieses Diffs

- `kategorie`: INFO
- `pfad`: `internal/adapters/driving/http/server_test.go:161`
  (`TestRegisterConsumerUngueltigesJSONEndetMit400`) gegen
  `internal/adapters/driving/http/registerconsumer.go:38`
- `befund`: Es gibt **sechs** Tests dieses Namens im Paket, nicht fünf. Fünf
  liegen im Diff dieses Vorgangs (`consumer.go:35`, `:117`,
  `verwaltung.go:45`, `:96`, `retention.go:36`) und binden (#18/Negativbefunde).
  Der sechste ist **nicht** Teil des Diffs (er stammt aus dem Adapter-Grundgerüst
  `7b6b253`, lange vor diesem Slice) — und er bindet **nicht**: schaltet man
  den `Decode`-Fehlerzweig in `registerconsumer.go:38` ab, bleibt er **grün**
  (Exit 0 gemessen), weil der leere Body danach an der Domänen-Invariante
  (`ErrEmptyIdentifier`) als `400` endet: er prüft die Formgrenze nicht, er
  sieht nur irgendwo ein `400`. Die Zahl „fünf" im Review ist damit richtig
  **als Zahl der Diff-Tests**; sie ist als Aussage über das Paket
  missverständlich, und der Befund ist ein neues Vorkommen derselben Klasse
  mit einem Vorgang, der **lange geschlossen** ist und **keine** Evidence-Datei
  hat (`negativtest-ohne-bindung-an-seine-eingabe/evidence/` führt
  `slice-083/086/087/088`, #17).
- `verifizierbar`: ja — `sed` auf `registerconsumer.go:38`
  (`err != nil` → `false && err != nil`), dann
  `go test -count=1 -run TestRegisterConsumerUngueltigesJSONEndetMit400
  ./internal/adapters/driving/http/` (gemessen Exit **0**); dieselbe Probe an
  den fünf Diff-Tests färbt **rot**
- `urteil`: **benannt, nicht gezählt** (Modul 6: „Ein Vorkommen **ohne**
  abgeschlossenen Vorgang bekommt keinen Beleg und bewegt den Zähler nicht; es
  gehört trotzdem in den Eintrag"). Kein Repair-Ort dieses Slice (er ändert
  keinen Produktionscode und keinen Bestandstest), aber ein Vorgang für
  `open/` — und ein Beleg dafür, dass die Zählung des Registers an dieser
  Stelle eine Lücke hat, die der Planner beim Lese-Schritt sehen sollte.

### V-8 — Zwei Genauigkeits-Notizen zur dritten Klasse (LP3)

- `kategorie`: INFO
- `pfad`: `internal/adapters/driving/grpc/streamv1/changestream_test.go:8–12`
  gegen `internal/adapters/driving/grpc/streamv1/changestream_grpc.pb.go:97`
  und `ADR-0082` §Konsequenzen
- `befund`: Zwei Punkte, beide **ohne** Wirkung auf LP3, beide für einen
  späteren Leser wert: (a) der Testkopf nennt sechs Klassen (`String`,
  `ProtoMessage`, `ProtoReflect`, `Descriptor`, Deskriptor-Kompression,
  `init`), `ADR-0082` §Konsequenzen nennt fünf — `ProtoReflect` fehlt dort. Es
  ist dieselbe Klasse (der ADR sagt „teils … (`…`)" und nennt Beispiele), aber
  die Aufzählung ist nicht deckungsgleich. (b) Von den **13** ungedeckten
  Blöcken des Pakets tragen **drei 0 Statements** (`ProtoMessage` ×2 und
  `mustEmbedUnimplementedChangeStreamServer`); die Aussage „die zehn sind
  Proto-Runtime-Interna" ist über die **Statements** exakt wahr (#6) — der
  dritte 0-Statement-Block ist **keine** Runtime-Interna und im Testkopf nicht
  genannt. Wer „ungedeckte Statements" und „ungedeckte Blöcke" verwechselt,
  liest die Benennung als vollständiger, als sie ist.
- `verifizierbar`: ja — Block-Auswertung des eigenen Profils für `streamv1`
  mit Statement-Zahl je Block; die Zeilen 59–241 in `changestream.pb.go`
- `urteil`: **keine Änderung nötig.** LP3s Buchstabe spricht von Statements,
  und über die Statements ist die Klasse vollständig gedeckt. Für eine
  spätere Auffrischung der Benennung ist es der Hinweis, dass der
  Zeichenvorrat der Benennung (5 bzw. 6 Namen) und der der Messung (15 Blöcke)
  zwei verschiedene Dinge sind.

---

## 7. Negativbefunde

- **geprüft, ohne Befund: „test-only" — die §1-Zusage und die Folgepflicht aus
  `ADR-0082` §Konsequenzen.** Zehn Pfade im Vorgang, **neun** davon
  `*_test.go`, der zehnte der Review-Report; **kein** `internal/**`/`cmd/**`
  außerhalb von Tests, kein `git mv`, keine ADR-Textänderung, kein
  `proto/`, keine `.pb.go`, kein `Makefile`-Ziel (#10). `THRESHOLD` steht
  unverändert auf **70**; die Cluster A, B und D sind pfadmäßig unberührt.
- **geprüft, ohne Befund: die Tests laufen wirklich im Gate.** Die
  `coverage`-Stufe fährt `go test` über die aus `go list` gebildete Liste (31
  Pakete, die vier des Clusters C darunter) und bricht bei rotem Testlauf ab;
  `make gates` ist **Exit 0** mit `77.20 %` (#1). `make test` ist **Exit 0**
  mit **32 × `ok`, 0 × `FAIL`** unter `--network none`, die vier Pakete
  ausgewiesen (#2). Die Tests laufen also **nicht** „lokal daneben": sie sind
  Teil beider Pfade.
- **geprüft, ohne Befund: die LP2-Bindung — 13 eigene Mutationsproben an der
  Produktionsseite, 13 × rot.** `consumer.go:35` · `verwaltung.go:45` ·
  `verwaltung.go:96` · `retention.go:36` · `consumer.go:117` (Gegenrichtung:
  **jeder** Body `400`) → die fünf `…UngueltigesJSONEndetMit400` **rot**;
  `grpc/server.go:134` Sendefehler verworfen → `TestStreamChangesSendefehlerWirdWeitergereicht`
  **rot**; `sse.go:94` `Flusher`-Zweig aus → `TestStreamOhneFlusherEndetMit500`
  **rot** (Panic im Handler); `interceptor.go:93` Prüfung aus →
  Interceptor-Negativfall **rot**; `middleware.go:81` `401`-Zweig aus →
  `TestStreamOhneTokenEndetMit401` und `TestRegisterConsumerOhneTokenEndetMit401`
  **rot**; `notify.go:123` Empty-Source-Ablehnung aus →
  `TestNotifyRejectsEmptySourceID` **rot**; `changestream_grpc.pb.go:121`
  `RecvMsg`-Fehler → `return nil` → `TestStreamChangesHandlerReichtEmpfangsfehlerWeiter`
  **rot**. Jede Probe mit Patch-Nachweis (Datei musste sich ändern) und eigener
  Kontrolle am unmutierten Stand (**Exit 0**). Die Klasse
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` tritt im Diff damit
  **nicht** ein fünftes Mal auf; der Fund V-7 liegt außerhalb.
- **geprüft, ohne Befund: die Mutationsangaben der Fixrunde halten alle drei.**
  `notify.go:81` (`log: outbound.NoopLog`) → `TestNewWithLogReichtDenLogPortDurch`
  **rot**; die zweite im Kommentar genannte rote Probe (`newOptions` ohne
  Options-Schleife) → **rot**; `interceptor.go:71` mit **gültigem** Token
  (`"reader-token"`) → Interceptor-Negativtest **rot**; mit **unbekanntem**
  Wert (`"probe"`) → **grün**; Entfernen des `!ok`-Zweigs (mit `md, _ :=` statt
  `md, ok :=`, damit es kompiliert) → **grün**. Genau die Verteilung, die die
  beiden Kommentare behaupten — die F-1-Korrektur hat den Test **gebunden**,
  nicht den Kommentar zurückgenommen.
- **geprüft, ohne Befund: die dritte Klasse ist entschieden, nicht erfunden.**
  Alle zehn ungedeckten Statements des `streamv1`-Pakets einzeln gegen den
  erzeugten Quelltext gehalten: `String` (2), `ProtoReflect` (2), `Descriptor`
  (2), Deskriptor-Kompression (3: `Do`-Aufruf, `CompressGZIP`-Zuweisung,
  Rückgabe), Wiedereintritt von `init` (1). `ADR-0082` §Konsequenzen schließt
  genau diese Klasse mit Begründung aus („prüft erzeugten Code gegen sich
  selbst"). Die zwei Notizen dazu stehen als V-8.
- **geprüft, ohne Befund: die zwei „für keine Eingabe erreichbar"-Stellen sind
  wirklich unerreichbar.** `readchanges.go:181` — `model.NewSourcePosition`
  scheitert nur bei leerer Quelle oder `offset == 0`; die Quelle ist oberhalb
  Pflichtfeld (`:139`), `offset >= 1` ist `:176` erzwungen. `retention.go:47–49`
  — `model.NewRetentionPolicy` (`internal/domain/model/retention.go`) hat genau
  die Bedingung `Nanos < 0`, `model.NewDuration` (`timepoint.go`) hat sie
  Zeilen vorher ausgeschlossen. Beide sind **tote Doppelprüfungen**, nicht
  bloß untestiert — die Benennung ist ehrlich, keine Kaschierung.
- **geprüft, ohne Befund: das Gate-Grün und die Zahl dazu sind reproduziert.**
  `make gates` **Exit 0**: `baseline-verify v6.5.0 OK — 54 Dateien`,
  `coverage-gate OK — Coverage 77.20% erfüllt Schwelle 70%`,
  `d-check: 748 Datei(en), 0 Befund(e)`, `commit-traceability OK — 5 Commit(s)`,
  `generated-sync OK` (beide `.pb.go` byte-gleich), `a-check: gesamt: 0 Befund(e)`;
  danach Baum sauber (#1). Mein eigenes Profil derselben Fläche kommt auf
  **77,2–77,3 %** je nach Schwankungsende (#3, #9) — die gedruckte Zeile des
  Gates und die eigene Auswertung sind **dieselbe** Messung in anderer
  Ausgabepräzision (jede Zuordnung `1903 − ungedeckt` trifft die gedruckte
  Prozentzeile exakt).
- **geprüft, ohne Befund: die zwei Doc-Gates über den Vorgang.**
  `make doc-commits RANGE=f90c3f4^..f9cd5e4` **Exit 0**, 748 Dateien,
  0 Befunde; `make doc-immutable RANGE=f90c3f4^..f9cd5e4` **Exit 0**, 0
  Befunde. Der Aufruf **ohne** `RANGE` endet Exit 2 — der Vorgang berührt
  keine `MR-*`-Datei, die Prüfung hat hier kein Objekt (#11, #12).
- **geprüft, ohne Befund: §3.7/§3.11 in den `+`-Zeilen.** Keine
  Slice-/Wellen-Nummer als Begründung (`slice-[0-9]`/`welle-[0-9]`: **0**
  Treffer), kein `t.Parallel` (**0**), keine host-lokale absolute Pfadform
  (**0**; die 15 Treffer eines groben Musters sind HTTP-URL-Pfade in
  Testaufrufen) — alle drei Greps über die `+`-Zeilen der **Go-Dateien**
  (#13). `docs-check` ist im Gate-Lauf grün (#1).
- **geprüft, ohne Befund: die Fixrunde hat keine vorbestehende Frist
  angetastet.** Vor dem Slice trug `sse_test.go` **6** und
  `grpc/server_test.go` **3** Fristen von `3 * time.Second`; der Diff fügte
  vier hinzu (9 und 4) und hob im Fixrunden-Commit **genau diese vier** auf
  30 s — die neun vorbestehenden stehen unverändert (#14). Die
  Zurückhaltung ist belegt, nicht behauptet.
- **geprüft, ohne Befund: das §2-Item „Reconciliation-Register … (entfällt)" ist
  nachweislich richtig.** `docs/plan/planning/reconciliation.md` existiert
  nicht (#16).

---

## 8. Was ich nicht prüfen konnte

- **Die Läufe des Implementers und des Reviewers.** Die Commit-Message nennt
  „+49, nicht +47" und der Review „M-14…M-18: 5 × rot" — beider Läufe sind
  nicht mehr einsehbar. Ich habe die **Substanz** unabhängig gemessen
  (die `±2`-Schwankung erklärt das Band 45…49, die fünf Diff-Tests färben
  rot); die **konkreten** Läufe sind nicht reproduzierbar und damit als Belege
  nicht überprüfbar.
- **Die Ursache des dritten Endes aus V-1.** Ein Lauf von 47 lag bei 435
  ungedeckten Statements; in 40 weiteren Läufen trat er nicht wieder auf, und
  in den 40 Läufen mit Blockmitschnitt war nur der `wiring.go`-Block instabil.
  Ein zweiter Flapper im Umfang von 1 Statement bleibt damit **unlokalisiert** —
  die Aussage „die `±2` sind vollständig erklärt" gilt für das beobachtete
  **Band**, nicht für jedes mögliche Ende.
- **Der reale Post-Push-Lauf.** Dieser Vorgang ändert **keinen** Workflow —
  `AGENTS.md` §3.10 greift dem Buchstaben nach nicht. Die Wirkung der neuen
  Tests auf die CI-Laufzeit bleibt bis zum nächsten echten Lauf unbelegt; die
  CI fährt `make test` bereits (nicht-blockierend für den E2E-Workflow).
- **Die Stabilität über die Rampe.** Ob die Zahl bei `THRESHOLD=80` (der
  `welle-20`-Closure) stabil genug ist, ist mit diesem Stand nicht
  entscheidbar: dafür braucht es die Messung **am** Hochschalt-Punkt, nicht die
  Extrapolation (13 Statements Abstand laut `ADR-0082` Kontext (5)).

---

## 9. Verdikt

**Die drei Liefer-Punkte aus §2 tragen — gemessen, nicht gelesen.** LP1: die
neun Testdateien über die vier Pakete sind real geliefert, laufen im Gate
**und** in `make test` netzlos (**Exit 0**, 32 × `ok`), `make gates` ist
**Exit 0**; der Zuwachs ist **62 → 15**, **+47**, Nenner unverändert 1903 — und
diese Zahl ist **flap-frei**, weil beide Vergleichsläufe denselben
`internal/bootstrap`-Stand tragen. LP2: **13 eigene Mutationsproben** an der
Produktionsseite färben **13 × rot**, darunter die fünf `…UngueltigesJSONEndetMit400`
des Diffs, die Gegenrichtung und die drei Klassen der DoD (`401`,
`Unauthenticated`, verweigerte Publikation). LP3: die verbleibenden
ungedeckten Statements sind für **alle** vier Pakete einzeln erhoben; die
Klassen „für keine Eingabe erreichbar" (drei HTTP-Statements, selbst gegen die
Domänen-Konstruktoren nachgeprüft), „lebender Dienst" (zwei in `natsnotify`)
und „Proto-Runtime-Interna" (zehn in `streamv1`, jeder gegen den erzeugten
Quelltext gehalten) **tragen** — die dritte ist in `ADR-0082` §Konsequenzen
entschieden, keine Ausrede.

**Entscheidungs-Konformität: hält.** Kein Produktcode, keine Naht,
`THRESHOLD` unverändert 70, Cluster A, B und D pfadmäßig unberührt, die
Plan-Zeile `welle-20.md` = „nicht" eingehalten, die `streamv1`-Zeile im Sinn
der ADR-Konsequenzen ausgefüllt. Der Plan-vs-Code-Diff ergibt **keine
unbegründete Abweichung**: alle vier Testzeilen sind geliefert, keine Zeile
fehlt.

**`done/`-fähig ist der Slice damit noch nicht — vier Punkte fehlen:**

1. **§7** — die Closure-Notiz mit Lerneintrag, und darin **die Zahl mit ihrem
   Lauf** (§5: die Statement-Zahl ist der Zustand, die gedruckte Gate-Zeile
   die Lauf-Größe) und **die LP3-Benennung für `http` und `natsnotify`**
   (**V-3**).
2. **V-4** — die Fixrunde `f9cd5e4` hat **keine** Review; ein Delta-Review
   (Präzedenz `slice-089`) oder eine ausdrückliche §7-Entscheidung fehlt.
3. **V-5** — die Häkchen (LP1–LP3, `make gates`, Review, Verifikation), die
   vier **§6-Risiko-Ausgänge** und die Register-Zeile; für die Ausgänge liefert
   §6/V-5 das Material.
4. **V-1/V-2** — die Zahl des Fortschritts ist als **beweglich** zu behandeln:
   sie trägt ihren Lauf, und wer die Schwankung begründet, nennt die
   **gemessene** Funktion (`runWALRetentionCheck`), nicht die im Sensor-Dokument
   genannte; die Register-Frage (Auftreten oder zweite Beobachtung) entscheidet
   der Lese-Schritt.

Nach 1–4 ist der Slice `done/`-fähig. **Kein Liefer-Defekt, kein rotes Gate:**
die offenen Punkte sind Übergebenes und Träger-Fragen, keine gebrochene Zusage.
Die zwei LOW/INFO-Funde außerhalb des Diffs (V-2, V-6, V-7) sind Adressen, keine
Reparatur-Orte dieses Slice.

---

**Beleg-Lage dieses Reports:** Jede Zahl stammt aus einem der Läufe in §1, je
in eigener Werkzeug-Beauftragung gefahren; der Gate-Lauf (#1) und seine
Auswertung waren **zwei** Schritte, sein Exit-Code wurde aus einer separaten
Datei gelesen, nie durch eine Pipe (§3.9). Die Mutationsproben liefen auf
**Arbeitsbaum-Kopien außerhalb des Repos**; der Baum ist nach diesem Bericht
sauber, **kein** Commit, keine Änderung an Artefakten des Slice.

## Nachtrag — Abschluss-Prüfung nach Closure und vierter Runde · 2026-09-16

**Auftrag:** enger Nachtrag zu **einem** Durchgang — die vierte Runde, das
D-2-Band, die Closure-Inhalte (vier §6-Ausgänge, Register-Zähler, §7-Zahlen)
und ein letztes Verdikt. §1–§9 dieses Berichts bleiben der Stand `f9cd5e4`.

**Stand:** `HEAD` = `58ff54e`, Baum sauber (`git status --porcelain` leer — vor
und nach jedem Lauf). Neue Commits seit §9: `a7d7f7b` (dritte Runde: W-1
Flap-Ursache, W-2 `streamv1`-Bild, W-3 sechster Test), `7e21e1b`
(**Delta-Review** `review-slice-091-delta.md` über `f9cd5e4` + `a7d7f7b`),
`af3ea9f` (vierte Runde: D-1), `58ff54e` (Closure: §6, §7, Häkchen, Register).
`git diff --name-only f9cd5e4..58ff54e` = **14 Pfade**, davon **13** Markdown
und **eine** Go-Datei: `internal/adapters/driving/http/server_test.go` (W-3).
Produktcode ist weiterhin **keine** Zeile berührt.

### N-1 Eigene Messungen dieses Nachtrags

| # | Lauf | Exit | Ergebnis |
|---|---|---|---|
| N-1.1 | `make gates` am Stand `58ff54e` (Log in Datei, Exit **danach** aus eigener Datei) | **0** | `baseline-verify v6.5.0 OK — 54 Dateien` · `coverage-gate: OK — Coverage 77.20% erfüllt Schwelle 70%` · `d-check: 756 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s), Betreffs ohne Struktur-ID` · `generated-sync: OK` · `a-check: gesamt: 0 Befund(e)`; danach Status leer |
| N-1.2 | `go list -f '{{.ImportPath}} Test={{len .TestGoFiles}} XTest={{len .XTestGoFiles}}'` über den Gegenstand (Arbeitsbaum-Kopie von `58ff54e`) | **0** | **31** Pakete; `Test=0 XTest=0`: **4** (`postgresstorage/queries`, `application/port/inbound`, `domain/errors`, `cmd/pg-change-feed`); `TestGoFiles=0`: **23**; davon `XTest>0`: **19** — 23 = 4 + 19 |
| N-1.3 | `streamv1` **eigener** Lauf (`-coverpkg` auf das Paket bzw. ohne), am Stand `58ff54e` | **0 / 0** | beide **`total: 52.3 %`** = **45 von 86** |
| N-1.4 | dieselbe Fläche im **Gegenstand** (volle Paketliste, `-coverpkg` über den Gegenstand) | **0** | `streamv1` **76 von 86** gedeckt (10 ungedeckt), **88,4 %** (abgeleitet aus 76/86); Gesamt in diesem Lauf `77.3 %` |
| N-1.5 | die Lauf-Ausgabe der Gegenstands-Pakete | **0** | `[no test files]` genau für die drei Pakete `queries`, `port/inbound`, `domain/errors`; `coverage: 0.0% of statements` **genau einmal** — `cmd/pg-change-feed`; `streamv1` steht als `ok … coverage: 2.4%` (= 45/1903) |
| N-1.6 | `go test -race -count=20` über die vier Cluster-C-Pakete | **0** | 4 × `ok`, **keine** Frist gefeuert (`http` 5,585 s · `grpc` 1,413 s · `streamv1` 1,013 s · `natsnotify` 1,052 s für alle zwanzig Iterationen) |
| N-1.7 | Register-Zähler gegen `ls evidence/ \| wc -l` (sieben Einträge) | **0** | `test-integration-retention-timing-flake` **3** (`slice-057/090/091`) · `beleg-befehl-traegt-seinen-satz-nicht` **3** (`slice-084/085/091`) · `negativtest-ohne-bindung-an-seine-eingabe` **5** (`slice-083/086/087/088/091`) · `arbeit-ueberholt-stehenden-traeger` **1** (neu, `slice-091`) · `zahl-in-traeger-driftet-gegen-die-messung` **6** (unverändert) |
| N-1.8 | DoD-Häkchen am Stand `58ff54e` / Zahl der DoD-Zeilen am verifizierten Stand `f9cd5e4` | **0** | **11 von 11** gesetzt; die Zeilen-Zahl war **11**, nicht 12 (siehe **N-6**) |
| N-1.9 | `grep -l af3ea9f docs/reviews/*.md` · `grep -n 'Fünf\|fünf' harness/sensors/coverage-gate.md` · Chronik-Wörter in den `+`-Zeilen von `af3ea9f` | **1 / 1 / 1** | kein Review nennt `af3ea9f`; **kein** Rest „Fünf/fünf" im Sensor-Dokument; **keine** Chronik-Wendung in seinen `+`-Zeilen → **N-1** |

### N-2 Die vierte Runde — die Zahlen tragen, und die neue Rolle ist substanziiert

**Alle Zahlen der neuen Fassung sind nachgemessen und halten** (N-1.2…N-1.5):
**4** Pakete mit `Test=0 XTest=0` — genau die vier genannten, in den zwei
Rollen (drei ohne ausführbare Statements, `cmd/pg-change-feed` vollständig
ungedeckt) —, **23** mit `TestGoFiles=0`, **19** davon ausschließlich mit
externem Testpaket. Die Gruppe ist damit mechanisch, und die Begründung nennt
die Zahl, die sie erzeugt.

**Trägt die neue Rolle von `streamv1`, oder ist sie umbenannt?** Sie **trägt**,
und die Unterscheidung ist gemessen, nicht behauptet: das Paket führt ein
**eigenes** externes Testpaket (`XTestGoFiles = 1`, `TestGoFiles = 0`) und hat
damit einen **eigenen Testlauf**, in dem es **45 von 86** Statements trägt
(`52,3 %`) — im **Gegenstand** sind es zugleich **76 von 86** (**88,4 %**, abgeleitet aus 76/86), weil
andere Testbinaries mit `-coverpkg` über die Paketgrenze messen. Die zwei
Zahlen sind verschieden **und** beide gemessen; „weder ohne Testdatei noch
allein über fremde Testpakete gedeckt" ist damit eine Eigenschaft, keine
Etikettierung. Zwei Nebenbelege: die Lauf-Ausgabe weist `streamv1` als
`ok … coverage: 2.4 %` aus (= 45/1903, dieselbe 45), und `coverage: 0.0% of
statements` trägt im Lauf **genau eine** Zeile — `cmd/pg-change-feed`
(N-1.5). Die zwei Nachbarsätze, die D-1 als mit diesem Slice falsch geworden
benannt hat, sind also beide berichtigt: die Überschrift führt **vier** Pakete
(kein „fünf" mehr im Dokument, N-1.9), und der `0.0 %`-Satz ist auf `cmd`
begrenzt.

### N-3 D-2 — die zwei Bänder: unterscheidbar, und die Ein-Wort-Frage

**Sie sind unterscheidbar** — beide Bänder tragen ihren Lauf-Marker
(„Lauf `slice-084`; im Lauf `slice-085` erneut beobachtet" bzw. „…, Lauf
`slice-091`"), und sie liegen ~100 Statements auseinander, weil sie **zwei
verschiedene Code-Stände** beschreiben. **Der schwache Punkt ist ein anderer:**
beide werden mit „**desselben** Stands" eingeleitet — einmal „desselben Stands"
(1369/1371), einmal „desselben **Produktionsstands**" (1468–1471) —, und beide
Male ist ein **anderer** Referent gemeint. Dass die Zahlen um ~100 steigen, hat
einen Grund (die Zwischenstände haben die Quote gehoben), und der Grund steht
nicht da.

**Ist die angebotene Ein-Wort-Umformulierung nötig?** **Für `done/`-fähig:
nein. Gerechtfertigt: ja.** Der Leser, der sich verwechselt, nimmt 1369
(71,9 %) für den heutigen Deckungsstand — 5 pp daneben, und genau diese Größe
entscheidet über den Abstand zur 80-%-Rampe. Die Form ist deshalb richtig (das
**Band** statt eines Einzelwerts), ihre Kennzeichnung ist vollständig, und die
Verwendung in §Grenze Punkt 4 ist in beiden Lesarten die konservative Richtung
(mit 1469 statt 1369 werden alle drei Rückrechnungs-Quotienten **höher**, die
Aussage „zwei bleiben grün" also eher bestätigt). Wenn etwas geändert wird,
dann **nicht mehr als das eine Wort** — den zweiten Referenten markieren
(etwa „Über acht Läufe am Stand **dieses Vorgangs**"). Eine Umformulierung
mehr wäre die Überholung an einer Stelle, die nicht driftet.

### N-4 Die Closure-Inhalte — §6, Register, §7

**Die vier §6-Ausgänge — Substanz geprüft:**

| Risiko | Ausgang (§6) | Mein Befund |
|---|---|---|
| **R1** die ≈52 sind eine Über-Schätzung | *eingetreten — und begrenzt:* real **47**, Puffer 13 → **8** | **trägt.** 47 ist meine eigene Messung (`62 → 15`); die **fünf** Differenz-Statements sind mit **denselben** Block-Positionen benannt, die ich in §1 Nr. 6 erhoben habe (`readchanges.go:181.3,182.1`, `retention.go:47.4,49.1`, `notify.go:136.2,137.12`); 47/52 = 90 % ist richtig gerechnet und trägt die Aussage „nicht *weit* unter ≈52". Die **8** ist eine **abgeleitete** Differenz (13 − 5): die 13 ist zitiert (`ADR-0082`), die 5 gemessen, die Subtraktion korrekt — als abgeleitet ist sie nicht gekennzeichnet, und sie unterstellt B (liegt in `done/`) und D (noch nicht geschnitten) auf ihren ADR-Zahlen → **N-4-Notiz** unten |
| **R2** ein Test wird zeitabhängig | *entfallen für die Tests — und bestätigt für den Bestand* | **trägt, und es ist die ehrliche Form.** `-race -count=20` über die vier Pakete: **Exit 0, 4 × ok, keine Frist gefeuert** (N-1.6) → für *die Tests* ist das Risiko entfallen. Dass die **Zahl** trotzdem schwankt (77,20 gegen 77,30 %; 2/1903 = **0,105 pp** — nachgerechnet), ist dem **Bestand** zugeordnet und nicht diesem Slice: derselbe Gegenstand, zwei Träger, sauber getrennt. Die Klasse trifft damit **nicht** in diesem Slice ein, und das wird nicht verschwiegen |
| **R3** ein Negativtest bindet an den Fake | *eingetreten — und behoben, zweimal* | **Substanz trägt, die Klassenzuordnung des ersten Falls nicht** → **N-2** |
| **R4** Coverage-Theater | *entfallen — gemessen:* Review **24** Proben (20 rot), Verifier **13** (13 rot) | **trägt.** Meine **13** in-diff-Proben waren sämtlich rot (§6/Negativbefunde), und die Zerlegung ist intern konsistent (24 + 13 = **37** Proben, 20 + 13 = **33** rot — genau die Zahlen der Commit-Message). Hinweis zur Lesart: die 37 sind die Proben **zweier** Berichte, nicht aller vier Runden (die Delta-Runden prüften weitere acht); die Untergrenze ist damit **konservativ**, aber „37" sollte niemand als Gesamtzahl lesen |

**Die Register-Zähler stimmen mit `ls evidence/` überein** (N-1.7): 3 · 3 · 5 ·
1 (neu) — und `zahl-in-traeger-driftet-gegen-die-messung` bleibt bei **6**,
was §7 auch nicht anders behauptet. Der neue Eintrag
`arbeit-ueberholt-stehenden-traeger` liegt vollständig vor
(`observation.md` ✓ mit Sub-Area-Angabe · `state.md` ✓ · `evidence/slice-091.md`
✓), und sein Beleg führt den Parent-Stand **als Messung** (`git archive
f90c3f4^` + `go list`) — die Aussage „die Liste *war* mechanisch" ist damit
datiert und nicht behauptet. Der Beleg zu `negativtest-…/evidence/slice-091.md`
nennt als Vorkommen korrekt **V-7** und trennt „Ursprung (`7b6b253`, `slice-061`)
≠ Vorkommen" — genau die Zuordnung, die §6 R3 offenlässt (**N-2**).
Die Paarungen-Vorprüfung in §7 hält: alle vier genannten Register-Adressen
existieren und führen ein nicht leeres `evidence/`.

**Nennt §7 eine Zahl ohne ihren Lauf oder ohne ihr Band?** **Das Band steht**
(„die Prozentzahl ist ein **Band** (77,1–77,3 %)") — das ist die richtige Form
für die bewegliche Zahl, und sie ist gegen meine 47 Läufe gedeckt. **Der Lauf
fehlt**, und zwar in einem Satz, der das Gegenteil behauptet: §7 schließt mit
„… und **beide Zahlen tragen ihren Lauf**", während im **ganzen Plan** kein
Lauf-Marker steht (gemessen: `grep` über den Plan → kein Treffer für
`Lauf slice-091`). Die Zahlen selbst sind richtig (62 → 15 ist meine Messung und
**stabil**; das Band deckt meine 47 Läufe) — falsch ist die Aussage **über** die
Notiz → **N-3**.

### N-5 Findings dieses Nachtrags

**N-1 — `af3ea9f` ist von keinem Review gedeckt (V-4 rezidiviert, schwächer).**
*Kategorie:* LOW. *Pfad:* `docs/reviews/` (drei Berichte, keiner nennt
`af3ea9f`) gegen `harness/sensors/coverage-gate.md` und §2 des Slice-Plans.
*Befund:* Der Delta-Review `7e21e1b` erklärt `f9cd5e4` **und** `a7d7f7b` zu
seinem Gegenstand und verlangt für D-1 einen Rückgabe-Pfeil („Ein-Stelle-
Korrektur …: die Gruppe auf ihren mechanischen Bestand bringen … und den
`0.0%`-Schlusssatz auf `cmd/pg-change-feed` begrenzen"); die vierte Runde
`af3ea9f` führt **genau diese** zwei Teile aus (+36/−26, zwei Träger) — geprüft
in N-2 —, und dieselbe Runde setzt in §2 die Häkchen **LP1, LP2, LP3,
`make gates`, Review, Verifikation, Reconciliation**. Für ihren Stand existiert
**kein** Review-Artefakt, obwohl der Delta-Review sein eigenes Verdikt mit
„DoD-Häkchen … **bleibt offen** — der Slice braucht eine (kurze) Fixrunde"
geschlossen hat. *Warum LOW und nicht MEDIUM wie `verify-slice-090`s
gleichnamige V-2:* dort änderte die Fixrunde den **Kern des Lieferwerts** und
legte einen neuen Träger an; hier führt sie eine vom Reviewer **wörtlich
vorgeschriebene** Korrektur in einem Zahlen-Träger aus, deren Zahlen dieser
Nachtrag vollständig nachgemessen hat (4 / 23 / 19 · 45/86 · 76/86 · die
Lauf-Zeilen), und deren neue Sätze ich als letzter Leser geprüft habe (kein
„fünf" mehr, keine Chronik-Wendung, Lauf-Marker vorhanden, N-1.9).
*Verifizierbar:* ja — `grep -l af3ea9f docs/reviews/*.md`;
`git diff --name-status f90c3f4^..58ff54e`. *Urteil:* **eine Zeile fehlt, nicht
eine Runde.** Entweder ein kurzes Delta auf `af3ea9f` (Präzedenz `slice-089`)
oder die ausdrückliche §7-Zeile „die vierte Runde führt die vorgeschriebene
Ein-Stelle-Korrektur aus; ihre Zahlen sind im Verifikations-Nachtrag geprüft,
ein weiteres Delta entfällt deshalb". Still bleiben darf es nicht.

**N-2 — §6 R3 zählt „zweimal" und nimmt F-1 in eine Klasse, die die erste Review
ausdrücklich außerhalb ihres Buchstabens verortet hat.**
*Kategorie:* LOW. *Pfad:* §6 R3 gegen `review-slice-091.md` §Antwort (2) und
`observations/BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe/evidence/slice-091.md`.
*Befund:* R3 sagt „eingetreten — und behoben, **zweimal**" und führt als ersten
Fall F-1 (`notify_test.go`: die `LogPort`-Weitergabe war nur bis zum ersten
Kettenglied gebunden). Die erste Review schreibt genau dazu: „der einzige Fund
dieser Achse (F-1) liegt **außerhalb des Buchstabens der Klasse**: die
Eingabeseiten-Mutation färbt dort rot — es fehlt das **mittlere** Kettenglied,
nicht die Eingabeseite" — und führt F-1 unter `beleg-befehl-traegt-seinen-satz-nicht`.
Der Register-Beleg dieses Slice sagt dasselbe („dass sie im selben Zug an einem
neuen Test (**F-1, dort als `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`
geführt**) und an diesem vorbestehenden auftrat"): die **Evidenz-Datei ist
präziser als der §6-Ausgang**. Die Substanz von R3 ist wahr (zwei Defekte
gefunden, beide gebunden, beide mutationsbelegt), die **Klassenzuordnung** des
ersten Falls ist es nicht. *Wirkung auf den Zähler:* **keine** — das Register
zählt Vorgänge, und `negativtest-…` steht zu Recht bei **5×** (das Vorkommen ist
V-7, N-1.7). *Verifizierbar:* ja — die beiden Zitatstellen; §6 R3. *Urteil:*
Klassenzuordnung der ersten Hälfte auf `beleg-befehl-…` umstellen (F-1 dort
belassen, wo der Register-Beleg es schon führt) oder „zweimal" fallenlassen.

**N-3 — §7 behauptet von zwei Zahlen, sie trügen ihren Lauf; im Plan steht
kein Lauf-Marker.** *Kategorie:* LOW. *Pfad:* §7
(`docs/plan/planning/in-progress/slice-091-coverage-cluster-c.md`, die
Lerneintrag-Zeile „Was hat funktioniert"). *Befund:* Der Satz nennt
**62 → 15** (ungedeckt) und das Band **77,1–77,3 %** und schließt mit „und
**beide Zahlen tragen ihren Lauf**". Gemessen enthält der **ganze** Plan keinen
Lauf-Marker (`grep 'Lauf slice-091'` → kein Treffer); die Zahlen sind ohne ihn
nicht auf einen Lauf zurückführbar — der Leser kann den Zeitpunkt nicht raten
(§3.12 Instanz A). Die **Substanz** ist in Ordnung: die Band-Form ist die
richtige Antwort auf die bewegliche Zahl, der Träger
`harness/sensors/coverage-gate.md` nennt für sein Band `Lauf slice-091`, und
`62 → 15` ist meine eigene, in 47 Läufen **stabile** Messung.
*Verifizierbar:* ja — `grep -n 'Lauf' …/slice-091-coverage-cluster-c.md`.
*Urteil:* einen Lauf-Marker in den Satz (etwa „… tragen ihren Lauf
(`slice-091`)") **oder** die Behauptung streichen und die Beobachtung als das
formulieren, was sie ist (die Trennung von Zustand und Lauf hat getragen).

**N-4 — der Puffer „13 → 8" ist eine abgeleitete Differenz ohne Kennzeichnung.**
*Kategorie:* INFO. *Pfad:* §6 R1, letzter Satz. *Befund:* Die 13 ist aus
`ADR-0082` zitiert, die 5 ist die gemessene Differenz (≈52 − 47), die Subtraktion
ist korrekt — die **8** wird aber als Tatsache vorgetragen, nicht als
**abgeleiteter** Wert (§3.12 Instanz A nennt Differenzen ausdrücklich), und sie
steht unter der Annahme, dass B (in `done/`) und D (noch nicht geschnitten) ihre
ADR-Zahlen liefern. *Urteil:* ein Wort („rechnerisch") genügt; kein Blocker.

**N-5 — „37 Mutationsproben" ist die Summe zweier Berichte, nicht aller vier
Runden.** *Kategorie:* INFO. *Pfad:* §6 R4 und §7 (erstes Element). *Befund:*
§6 legt die Zusammensetzung offen (Review **24**, Verifier **13**); die
Delta-Runden haben darüber hinaus geprüft (F-1: drei, F-2: drei, W-3: zwei).
Die Summe ist damit eine **Untergrenze** — konservativ, aber wer sie später als
Gesamtzahl des Vorgangs liest, zählt zu niedrig. *Urteil:* keine Änderung nötig;
die Zerlegung steht ja daneben.

**N-6 — Korrektur an meinem eigenen Bericht: es waren 11 DoD-Zeilen, nicht 12.**
*Kategorie:* INFO (Selbstkorrektur, `AGENTS.md` §3.12). *Pfad:* §1 Nr. 15, §2,
V-5 und §9 dieses Berichts („0 von **12** Häkchen"). *Befund:* Gemessen am
verifizierten Stand `f9cd5e4` trägt der Plan **11** DoD-Zeilen
(`grep -c '^ *- \['` → **11**), von denen **keine** gesetzt war; gesetzt sind
jetzt **11 von 11** (N-1.8). Die Zahl in diesem Bericht stammt nicht aus einer
Zählung, sondern aus dem Nachbarbericht — genau die Übernahme, gegen die §3.12
geschrieben ist. *Urteil:* **Lauf-Belege werden nicht rückdatiert** — die
Korrektur steht hier, der Bericht bleibt, wie er war. Der Slice-090-Vergleich
in V-5 bleibt in der Sache richtig.

**N-7 — die untere Kante des Bands ist selten, und das steht nicht dabei.**
*Kategorie:* INFO. *Pfad:* `harness/sensors/coverage-gate.md` §Zählbasis (der
Acht-Läufe-Absatz). *Befund:* Die Form ist richtig (Band, nicht Einzelwert) und
sie trägt ihren Lauf; die untere Kante **1468** trat in meinen 47 Läufen der
vollen Fläche **einmal** auf, in den 54 Läufen des Delta-Reviews **keinmal** —
die Formulierung „über **acht** Läufe … zwischen 1468 und 1471" ist als
Lauf-Beleg zulässig (sie ist ihrer), lässt aber offen, dass ihr unteres Ende ein
~1-%-Fall ist. *Urteil:* optional ein Halbsatz; **kein** Drift (die Zahl ist
nicht als Dauerwert ausgegeben). D-2s Substanz ist damit erledigt.

### N-8 Verdikt des Nachtrags

**Die vierte Runde trägt, und die neue Rolle von `streamv1` ist substanziiert.**
Alle Zahlen sind an **einem** Stand (`58ff54e`) unabhängig nachgemessen:
**4** Pakete `Test=0 XTest=0` (genau die vier genannten), **23** mit
`TestGoFiles=0`, **19** davon mit externem Testpaket; `streamv1` **45 von 86**
(`52,3 %`) im **eigenen** Lauf gegen **76 von 86** (`88,4 %`) im **Gegenstand**;
die Lauf-Ausgabe stützt beides (`[no test files]` genau für die drei
statement-losen Pakete, `coverage: 0.0% of statements` genau einmal für `cmd`).
Die zwei von D-1 benannten Nachbarsätze sind berichtigt, kein „fünf" mehr im
Dokument, kein Chronik-Ton in den `+`-Zeilen. **D-1 ist damit erledigt.**

**Die Closure-Inhalte tragen** — die drei Liefer-Punkte unverändert (§2), die
vier §6-Risiken mit je **einem** Ausgang, die Register-Zähler deckungsgleich mit
`ls evidence/` (3 · 3 · 5 · 1 neu, `zahl-in-traeger-…` unverändert 6), das
Häkchen-Bild vollständig (11/11), und `make gates` am Closure-Stand **Exit 0**
(`77.20 %`, d-check 756 Dateien, 0 Befunde).

**Offen sind drei Ein-Zeilen-Nachträge** — `af3ea9f` ohne Review-Artefakt
(**N-1**, LOW), die Klassenzuordnung in §6 R3 (**N-2**, LOW) und der fehlende
Lauf-Marker in §7 (**N-3**, LOW) —, dazu zwei INFO-Nuancen (N-4, N-7) und eine
Selbstkorrektur (N-6). Keiner ist ein Liefer-Defekt, keiner rührt an Zahlen,
die tragen, und keiner berührt Produktcode oder Schwelle.

**`done/`-fähig: ja** — mit der Auflage, dass die drei LOW-Nachträge **als
Text** geschrieben werden, bevor der Slice schließt: N-2 und N-3 gehören in die
Closure-Notiz (sie wandert mit nach `done/`), N-1 braucht eine Entscheidung —
kurzes Delta auf `af3ea9f` **oder** die ausdrückliche Zeile in §7. Alles drei
ist Ein-Zeilen-Arbeit an vorhandenen Trägern; keine neue Runde, kein neues
Artefakt, keine Messung.

**Offen über diesen Slice hinaus** (unverändert aus §8): der reale Post-Push-Lauf
und die Stabilität der Zahl bei `THRESHOLD=80` — die `welle-20`-Closure misst
dort, nicht hier.
