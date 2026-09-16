# Review-Report: slice-091 — Coverage Cluster C (Zustell- und Betriebs-Rand) · 2026-09-16

**Review-Art:** Code — geprüft gegen **Plan und Entscheidungen**
(Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-10-review-harness.md`
§Drei Review-Arten). DoD-/Spec-Konformität (Verifier, Modul 11),
§6-Risiko-Ausgänge, Beobachtungs-Register und die drei Paarungen
(Planner-Closure) sind **nicht** Gegenstand dieses Reports.

**Gegenstand:** `f90c3f4` (Parent `22cebec`), ein Commit — neun Dateien, 774
hinzugefügte Zeilen, **ausschließlich `*_test.go`**. Gemessen:
`git diff --name-only f90c3f4^..f90c3f4` listet neun Pfade, davon **null**
ohne die Endung `_test.go`. Arbeitsbaum sauber (`git status --porcelain` leer).
Kein Produktcode, keine Naht, kein Gate, keine Schwelle berührt.

**Skill:** `.harness/skills/reviewer.md` @ `f90c3f4` — die vier
repo-spezifischen HIGH-Unterpunkte (u. a. „Zusage ohne Bindung an ihre
Eingabeseite“ und „Zahl im Träger“) gehören zum Prüfraster.
**Modell:** `deepseek-v4.1-flash:cloud[1m]` · **Datum:** 2026-09-16.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-091-coverage-cluster-c` vollständig (§1–§8), die offene
  Welle `welle-20`
- `ADR-0082` (das Schnittmaß; §Konsequenzen „test-only“, Cluster-Tabelle,
  §Re-Evaluierungs-Trigger), `ADR-0071` (Messgegenstand), `ADR-0077`
  (Rampe 70 → 80), `ADR-0057`/`ADR-0060`/`ADR-0055` (die vier Pakete),
  `ADR-0024` (`LogPort`, `LH-QA-OPS-004`), `ADR-0081`/`ADR-0056`/`ADR-0061`
  (die berührten Zusagen), `ADR-0080` (Nahtform — hier nicht angewandt)
- Berührte `LH-*`: `LH-FA-SST-006`, `LH-FA-SST-007`, `LH-FA-SST-008`,
  `LH-QA-OPS-004`, `SPEC-017`, `SPEC-020`
- `AGENTS.md` §3.1, §3.3, §3.6, §3.7, §3.9, §3.11, §3.12, §4, §6 ·
  `harness/conventions.md` (MR-000 ID-Schema)
- Letzte fünf Reviews am gleichen Modul: `review-slice-085`/`-086`/`-087`/
  `-088`/`-089` — keine Klasse doppelt zu diesem Diff
- Beobachtungs-Register `BEO-PGC/` (Zähler mit `ls evidence/ | wc -l`
  nachgezählt, nicht aus Prosa übernommen — §8 des Plans stimmt in **allen
  sechs** genannten Zahlen)

---

## Findings

### F-1 — Der Test bindet nur den Konstruktionsaufruf; das im Namen behauptete Durchreichen ist ungeprüft, und die zweite genannte Mutation färbt nicht rot

- `kategorie`: MEDIUM
- `quelle`: Maintainability · `ADR-0024`/`LH-QA-OPS-004` ·
  `AGENTS.md` §3.12 Instanz B · Register-Muster
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (dortige Mechanik:
  „eine Aussage, die nicht an ihre Eingabeseite gebunden ist, ist grün ohne
  Aussage“) · `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`
- `pfad`: `internal/adapters/driven/natsnotify/notify_test.go:190-197`
  (Testkörper `:198-209`, Produktionsstelle `notify.go:76-81`)
- `befund`: Der Testkopf sagt „der über `WithLog` übergebene `LogPort` ist
  der, über den der Adapter protokolliert“, der Name sagt
  „ReichtDenLogPortDurch“. Gemessen ist nur das **erste** Kettenglied: der
  Aufruf `o.log.Info(...)` in `New` (`notify.go:80`) läuft über den
  injizierten Port. Das **mittlere** Glied — die Weitergabe in das
  Adapterfeld (`notify.go:81`) — ist durch keine Zusage gedeckt: mit
  `log: outbound.NoopLog` an genau dieser Stelle bleibt der neue Test **und
  das ganze Paket** grün (`go test ./internal/adapters/driven/natsnotify/`,
  EC=0), obwohl der Adapter danach über keinen realen LogPort mehr
  protokolliert. Dieselbe Mutation, die der Testkopf als „rot färbend“
  nennt („oder in `New` den Log auf `outbound.NoopLog` festlegen“), leert
  die Aufzeichnung nicht — die behauptete Wirkung („bliebe die Aufzeichnung
  leer“) tritt bei der genannten Mutation nicht ein. Kein anderer Test des
  Pakets fängt die Stelle: `TestNotifyWrapsPublishFailureAsTransient`
  (`:116`) und `TestNotifyPublishesEmptyPayloadOnSubject` (`:135`) rufen
  `New` **ohne** Option, laufen also über `NoopLog`.
- `verifizierbar`: ja — Mutation `log: o.log` → `log: outbound.NoopLog`
  (`notify.go:81`), dann `go test ./internal/adapters/driven/natsnotify/`
  (gemessen EC=0). Die Eingabeseiten-Mutation des Kopfes (`newOptions`-Schleife
  fallenlassen) färbt dagegen real rot (gemessen EC=1) — die Kette ist damit
  **teilweise** gebunden, nicht ganz; das ist der Grund für MEDIUM und nicht
  für HIGH.
- `klasse`: „Durchreichen behauptet, nur das erste Kettenglied gebunden“

### F-2 — Die genannte Rot-Mutation des Interceptor-Tests trägt ihren Satz nicht

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12 Instanz B („eine Begründung, die eine Tatsache
  über den Gegenstand behauptet, nennt den Beleg-Anker … oder sie ist als
  erwartet formuliert“) · Register `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`
- `pfad`: `internal/adapters/driving/grpc/interceptor_test.go:56-57`
  (gegen `internal/adapters/driving/grpc/interceptor.go:68-80`)
- `befund`: Der Testkopf nennt als rot färbende Mutation „in `credentialToken`
  an der Stelle ohne eingehende Metadata einen Wert statt `""` zurückgeben —
  dann wird der Handler erreicht“. Gemessen: mit `return "probe"` an
  `interceptor.go:71` bleibt der Test **grün** (EC=0) — `classifyToken`
  ordnet einen unbekannten Wert weiterhin `roleNone` zu, der Handler wird
  nicht erreicht. Rot färbt erst ein **gültiger** Token-Wert an derselben
  Stelle (gemessen mit `return "reader-token"`: EC=1). Die Aussage ist damit
  als Rezept nicht reproduzierbar; wer sie wörtlich nachfährt, findet die
  Zusage schwächer, als sie ist. Der zweite Satz desselben Kommentars
  („Das bloße Entfernen des `!ok`-Zweigs färbt diesen Test **nicht** rot“)
  hält dagegen — gemessen EC=0.
- `verifizierbar`: ja — Mutation `interceptor.go:71` (`""` → `"probe"` bzw.
  → `"reader-token"`), dann `go test -run
  TestAuthStreamInterceptorOhneEingehendeMetadataEndetMitUnauthenticated
  ./internal/adapters/driving/grpc/`
- `klasse`: `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (dritter Vorgang,
  **eine** Evidence-Datei — F-1 und F-2 sind zwei Funde im **selben** Vorgang,
  der Zähler zählt Vorgänge, nicht Funde)

### F-3 — Die drei Sekunden sind ihrer Form nach eine Zusicherung, nicht nur ein Hänge-Schutz

- `kategorie`: INFO
- `quelle`: Maintainability · Register
  `BEO-PGC/test-integration-retention-timing-flake` (Mechanismus: das Ergebnis
  wechselt bei unverändertem Stand)
- `pfad`: `internal/adapters/driving/http/sse_test.go:400`, `:438`, `:475` ·
  `internal/adapters/driving/grpc/server_test.go:326`
- `befund`: In drei der vier neuen Fristen ist die Frist **zugleich** die
  Beobachtung: „der Handler endete“ ist nur über `<-fertig` innerhalb von
  `3 * time.Second` feststellbar; läuft die Goroutine nicht binnen 3 s an,
  färbt der Test rot, obwohl das Verhalten richtig ist. Der Nenner ist eine
  Uhr, kein Zustand — die Grenze ist also benannt, aber sie ist keine reine
  Hänge-Sicherung. Gemessen ist die Marge groß: 20 Wiederholungen je Paket
  unter `-race` (`-count=20`) laufen grün durch, die Frist feuert in keinem
  Lauf; die vier Pakete brauchen zusammen **9,18 s für alle 20 Iterationen**
  (gemessene Paketzeiten 5,757 + 1,345 + 1,025 + 1,050 s) — der ganze
  `http`-Testlauf mit allen seinen Handlern liegt also bei ≈0,29 s je
  Iteration, die Frist bei 3 s **je einzelnem Handler**. Ein
  Beleg für ein drittes Auftreten der Register-Beobachtung ist aus diesem
  Diff nicht ableitbar: kein Lauf wechselte bei unverändertem Stand das
  Ergebnis.
- `verifizierbar`: ja — `go test -race -count=20 <vier Pakete>`, gemessen
  EC=0 (4 × ok)
- `klasse`: „Frist als Zusicherung statt als Hänge-Schutz“

### F-4 — Die Aussage „alle nur `docs/`“ über den bewegten HEAD ist contra-gemessen (Träger außerhalb des Repos)

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 Instanz B · Register
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (Verwandtschaft:
  Aussage über den Gegenstand gegen die Messung)
- `pfad`: — (die Aussage steht im Übergabe-Bericht an den Reviewer, nicht in
  einem committeten Träger; kein Reparatur-Ort in diesem Diff)
- `befund`: Der Bericht nennt den Bewegungslauf „`3de7fb1` → `22cebec`, ~60
  Commits, **alle nur `docs/`**“. Gemessen: `git log --oneline
  3de7fb1..22cebec | wc -l` = **55**, und der Bereich berührt **zwölf** Pfade
  außerhalb von `docs/` — darunter die Cluster-B-Testdateien
  (`internal/adapters/driving/replication/decode/*`, `replication/mapper/*`,
  `postgresstorage/sqlexec/translate_test.go`, `postgresstorage/mapper/mapper_test.go`)
  und `slice-090`s Werkzeug-Zug (`tools/harness/generated-sync.sh`,
  `harness/mk/generated-sync.mk`, AGENTS.md, `.harness/skills/reviewer.md`).
  **Die Substanz der Frage ist davon unberührt:** keiner der **neun**
  Diff-Pfade kommt im Bereich vor (`git diff --name-only 3de7fb1..22cebec`
  enthält keine der neun Dateien), und die Zahlen des Slices sind am Parent
  reproduzierbar (Messung unten: 62 ungedeckt paketweise exakt). Die
  Fehlaussage ist damit ein Genauigkeits-Mangel des Berichts, kein
  Befund am Diff.
- `verifizierbar`: ja — `git log --oneline 3de7fb1..22cebec | wc -l`;
  `git diff --name-only 3de7fb1..22cebec | grep -v '^docs/'`;
  `git diff --name-only 3de7fb1..22cebec | grep -Ff <die neun Diff-Pfade>`
  (leer)
- `klasse`: „Aussage über den Gegenstand ohne ihren Beleg-Anker“

---

## Negativbefunde

- geprüft, ohne Befund: **die Zahlen des Commits — selbst nachgemessen.** Ein
  eigener Lauf über die Paketliste des Gates (`go list ./internal/... ./cmd/...`
  ohne `postgresstorage|postgresack|replication/receive`), Profil
  `-coverpkg=<alle> -covermode=atomic`, dedupliziert über die Block-Position
  (dieselbe Zählbasis wie `tools/harness/db-coverage.sh`), im netzlosen
  Toolchain-Container:
  | Stand | ungedeckt | Nenner | `go tool cover -func` total |
  |---|---|---|---|
  | `22cebec` (Parent) | **479** | **1903** | **74.8 %** |
  | `f90c3f4` (Diff) | **434** | **1903** | **77.2 %** |

  Paketweise Cluster C, beide Stände: `driving/http` **29 → 3**,
  `driving/grpc/streamv1` **25 → 10**, `driving/grpc` **3 → 0**,
  `driven/natsnotify` **5 → 2** — vorher **62** (exakt die Zeile aus
  `ADR-0082` Kontext (4b)), nachher **15**, Zuwachs **+47**, Nenner
  unverändert. Die im Commit behauptete Abweichung „+49, nicht +47“ ist die
  dokumentierte Schwankung und **ebenfalls nachgemessen**:
  `runWALRetentionCheck` (`internal/bootstrap/wiring.go:981`) trägt im
  Parent-Lauf **100.0 %**, im Diff-Lauf **87.5 %** (±2 Statements). Alles in
  der Commit-Message genannte Zahlenwerk hält.
- geprüft, ohne Befund: **die fünf benannten unerreichbaren Statements**
  (LP3-Substanz, hier nur auf Ehrlichkeit geprüft). Ungedeckt bleiben real:
  `internal/adapters/driving/http/readchanges.go:181` (1), `retention.go:47`
  (2), `internal/adapters/driven/natsnotify/notify.go:136` (2) plus die zehn
  Proto-Runtime-Interna in `changestream*.pb.go`. Die drei HTTP-Statements
  sind **tote Doppelprüfungen**, nicht bloß untestiert:
  `readchanges.go:176` erzwingt `offset >= 1` und `source` ist oberhalb als
  Pflichtfeld geprüft (`:140`), `model.NewSourcePosition` kann danach nicht
  mehr scheitern (`internal/domain/model/position.go:23-30`); `retention.go:45`
  ruft `NewRetentionPolicy`, dessen einzige Bedingung
  (`internal/domain/model/retention.go:23`) „`Nanos < 0`“ ist —
  `model.NewDuration` (`internal/domain/model/timepoint.go:54-59`) hat sie
  Zeilen vorher bereits ausgeschlossen. `notify.go:136` ist der
  Publish-Erfolgspfad und braucht einen lebenden NATS-Server (`make test-notify`).
- geprüft, ohne Befund: **die Prüf-Kraft der neuen Zusagen — 20 von 24 eigenen
  Mutationsproben färben rot.** (Vier Proben waren in meinem ersten Anlauf
  falsch gepatcht — der Patch traf nicht, die Probe war leer; sie sind in der
  Tabelle unten durch nachgeführte Proben ersetzt, nicht mitgezählt.)
  Reproduziert wurden die im Diff **genannten** Mutationen und, wo der
  Kommentar keine nennt, die Eingabeseiten-Mutation des Gegenstands:

  | Probe | Produktionsmutation | Test | Ergebnis |
  |---|---|---|---|
  | M-2 | `grpc/server.go:133` Sendefehler verworfen | `TestStreamChangesSendefehlerWirdWeitergereicht` | **rot** |
  | M-3 | `natsnotify/notify.go:46` Options-Schleife weg | `TestNewWithLogReichtDenLogPortDurch` | **rot** |
  | M-5 | `interceptor.go:71` gültiger Token statt `""` | `TestAuthStreamInterceptor…Unauthenticated` | **rot** |
  | M-6 | `interceptor.go:70` `!ok`-Zweig entfernt | derselbe | grün (**vom Kopf als Negativkontrolle benannt**) |
  | M-7 | `http/sse.go:114` Marshal-Fehler verworfen | `TestStreamNichtKodierbareChangeBeendetDenStream` | **rot** |
  | M-9 | `http/errors.go:25` 404 → 500 | `TestGetConsumerPositionFehlerWirdNachKlasseAbgebildet` | **rot** |
  | M-10 | `http/server.go:128` `ErrServerClosed` als Fehler | `TestServerShutdownVorStartIstRegulaererAusgang` | **rot** |
  | M-11 | `http/server.go:128` Bindfehler verworfen | `TestServerStartMeldetBindfehler` | **rot** |
  | M-12 | `http/sse.go:93` `!ok`-Zweig entfernt | `TestStreamOhneFlusherEndetMit500` | **rot** (Panic im Handler) |
  | M-14…M-18 | je `Decode`-Fehlerzweig entfernt (`consumer.go:35`, `:117`, `verwaltung.go:45`, `:96`, `retention.go:36`) | die fünf `…UngueltigesJSONEndetMit400` | **5 × rot** |
  | M-19 | `changestream_grpc.pb.go:57` `NewStream`-Fehler → `return stream, nil` (Zweig `:56`) | `TestStreamChangesClientReichtFehlerWeiter` | **rot** |
  | M-22 | `changestream.pb.go:102` `GetSequence` ohne Nullwert-Zweig | `TestChangeGetterTragenDenNullwertNilSicher` | **rot** |
  | M-23 | `grpc/server.go:94` Listener-Fehler verworfen | `TestServeMeldetListenerFehler` | **rot** |
  | M-24 | `changestream_grpc.pb.go:121` `RecvMsg`-Fehler → `nil` | `TestStreamChangesHandlerReichtEmpfangsfehlerWeiter` | **rot** |
  | M-25 | `changestream_grpc.pb.go:95` Statuscode geändert | `TestUnimplementedServerEndetMitUnimplemented` | **rot** |
  | M-27 | `http/sse.go:119` `Fprintf`-Fehler verworfen | `TestStreamSchreibfehlerBeendetDenStream` | **rot** |
  | M-28 | `consumer.go:117` **jeder** Body mit 400 abgelehnt | `TestRemoveConsumerUngueltigesJSONEndetMit400` | **rot** — über die **Gegenprobe** |
  | M-4 | `interceptor.go:71` `"probe"` statt `""` | `TestAuthStreamInterceptor…` | grün → **F-2** |
  | M-26/M-26b | `notify.go:81` `log: outbound.NoopLog`, **ganzes Paket** | `TestNewWithLogReichtDenLogPortDurch` | grün → **F-1** |

  Die fünf `…UngueltigesJSONEndetMit400`-Tests sind Muster-Beispiele für LP2:
  jeder paart den Negativfall (nicht dekodierbarer Body → 400) mit einer
  **Gegenprobe** (gültiger Body → 404 aus **demselben** Fake). Die Gegenprobe
  ist nicht dekorativ: M-28 (jeder Body abgelehnt) färbt real rot — die
  Gegenrichtung ist mitgemessen, nicht nur behauptet.
- geprüft, ohne Befund: **keine `t.Parallel`-Nutzung, kein geteilter Zustand.**
  `grep -c 't.Parallel'` über den Diff = **0**; alle neuen Doppel
  (`recordingLog`, `fakeServerStream`/`fakeChangeServerStream`,
  `fakeClientConn`/`fakeClientStream`, `responseWriterOhneFlusher`,
  `schreibfehlerWriter`, die fünf `…FailingUseCase`) sind pro Test
  instanziiert, keiner ist paketweit veränderlich; die fünf neuen Fakes folgen
  dem im Paket bereits vorhandenen Haus-Muster (ein kleiner Fehler-Fake je
  Use-Case-Schnittstelle, vgl. `consumer_test.go:93`, `retention_test.go:29`).
- geprüft, ohne Befund: **Netzlosigkeit.** Die vier Pakete laufen unter
  `--network none`: `go test -race ./internal/adapters/driving/http/
  ./internal/adapters/driving/grpc/ ./internal/adapters/driving/grpc/streamv1/
  ./internal/adapters/driven/natsnotify/` → 4 × `ok`, EC=0. Der
  `grpc`-Test nutzt `bufconn` (in-memory), der `streamv1`-Test gar keinen
  Transport (`fakeClientConn`/`fakeServerStream`), die `sse`-Tests
  `httptest.NewRequest`/`httptest.NewRecorder`, `natsnotify` die
  Loopback-Adresse `nats://127.0.0.1:1` (Port 1, kein Server).
- geprüft, ohne Befund: **die Zusicherung „braucht keinen Socket“**
  (`http/server_test.go`, `TestServerShutdownVorStartIstRegulaererAusgang`).
  Nachgemessen an der Standardbibliothek: `net/http.(*Server).ListenAndServe`
  prüft `s.shuttingDown()` **vor** `net.Listen` und kehrt mit
  `ErrServerClosed` zurück (`<GOROOT>/src/net/http/server.go:3450-3456`, das
  `GOROOT` des gepinnten Toolchain-Containers) — der Kommentar des Tests hält
  wörtlich.
- geprüft, ohne Befund: **die genannten Kontrollen existieren.** Der in
  `sse_test.go` zitierte Kontrolltest `TestStreamReaderTokenOeffnetTraegtChange`
  steht real in derselben Datei (`:178`); der in `grpc/server_test.go`
  zitierte reguläre Ausgang `grpc.ErrServerStopped` hat seinen Träger in
  `TestStartUndShutdown`. Die im Interceptor-Test behauptete
  Unerreichbarkeit über eine reale Verbindung ist mit dem Profil
  gegengedeckt: im Parent-Lauf waren genau `interceptor.go:71`,
  `grpc/server.go:95` und `grpc/server.go:134` ungedeckt — der `!ok`-Zweig
  war also unter den bestehenden `bufconn`-Tests real unerreicht.
- geprüft, ohne Befund: **§3.7/§3.12 in den Kommentaren des Diffs.** Keine
  Slice-/Wellen-Nummer in einem Kommentar (`grep -nE 'slice-[0-9]|welle-[0-9]'`
  über die `+`-Zeilen: leer), kein Vorher/Nachher-Vokabular, kein
  abwesender Text. Die Testköpfe tragen die Klassen Zusage/Kopplung/Abgrenzung
  und zitieren `ADR-*`/`LH-*`/`SPEC-*` sowie Register-Pfade als Herkunft. Die
  Zahlen, die sie nennen, sind Mutationsangaben — die zwei tragenden stehen
  als F-1/F-2, nicht als Formfehler.
- geprüft, ohne Befund: **§3.11.** Keine `+`-Zeile des Diffs nennt einen
  host-lokalen absoluten Pfad (Präfix-Grep über die `+`-Zeilen: leer); der
  Diff ist test-only, `make docs-check` ist grün.
- geprüft, ohne Befund: **die §3-Bedingung des Plans für
  `harness/sensors/coverage-gate.md`.** Die Datei wurde nicht angefasst — zu
  Recht: ihre Statement-Zahlen tragen durchweg ihren **Lauf** (`Lauf
  slice-081/084/085/089`) und sind in §Zählbasis ausdrücklich als bewegliche
  Zustands- bzw. Lauf-Größen deklariert; der Nenner `1903` gilt nach meiner
  Messung unverändert. Die Zeile „trägt **86 Statements, 61 gedeckt** (**70,9 %**, abgeleitet aus
  61/86; Lauf `slice-089`)“ ist nach diesem Diff **nicht mehr aktuell** (gemessen: 76 von 86 = 88,4 %) — sie ist als Lauf-Beleg zitiert,
  nicht als Ist-Stand, und damit kein Drift im Sinn der Klasse. Für die
  Closure ist es eine Auffrischungs-Option, kein Fund.
- geprüft, ohne Befund: **Plan-vs-Code.** §3 nennt vier Test-Pakete — alle
  vier real geliefert; die Zeilen `welle-20.md` §4 = „**nicht**“ und
  `docs/user/benutzerhandbuch.md` (nicht gelistet, weil keine
  Betreiber-Oberfläche entsteht) sind eingehalten. Die `streamv1`-Zeile („Test
  neu, **falls** dort prüfbarer Code liegt; der generierte Stub ist **nicht**
  Gegenstand“) ist mit `ADR-0082` §Konsequenzen im Einklang: der Diff prüft
  Getter, Client- und Handler-Verklebung und lässt die Proto-Runtime-Interna
  liegen — kein Test gegen den erzeugten Text, nur gegen sein Verhalten.
- geprüft, ohne Befund: **`ADR-0082` Folgepflicht „test-only“.** Gemessen:
  neun Dateien, alle `_test.go`; `internal/**.go`/`cmd/**.go` außerhalb von
  Tests unberührt. Kein `git mv`, kein ADR-Text, keine Schwelle (§3.3, §3.5,
  §3.6 unberührt), `THRESHOLD` unverändert 70.
- geprüft, ohne Befund: **Register-Zähler des Plans §8** —
  `ls <eintrag>/evidence/ | wc -l`: `negativtest-ohne-bindung-an-seine-eingabe`
  **4**, `zahl-in-traeger-driftet-gegen-die-messung` **6**,
  `dod-begruendung-unzutreffende-tatsachenbehauptung` **4**,
  `test-integration-retention-timing-flake` **2**, `a-check-null-abdeckung`
  **3**, `db-gegenstand-enthaelt-netzlos-geprueften-code` **2** — alle sechs
  wie im Plan genannt.
- geprüft, ohne Befund: **`make gates` real.** Log in eine Datei, Exit-Code
  **danach separat** gelesen (§3.9): **EC=0**. Sechs Checks gefahren
  (`baseline-verify` OK, `coverage-gate` **OK — Coverage 77.20 % erfüllt
  Schwelle 70 %**, `commit-traceability` OK, `generated-sync` OK — beide
  `.pb.go` byte-gleich, `docs-check` grün, `a-check` `gesamt: 0 Befund(e)`),
  danach leerer `git status --porcelain`. Die neue Testdatei im
  Erzeugnis-Verzeichnis stört `generated-sync` nicht (sie ist nicht als
  Erzeugnis gekennzeichnet).
- geprüft, ohne Befund: **`make test` real** (voller Unit-Lauf mit
  Race-Detector, netzlos) — EC=0, **32 × `ok`, 0 × `FAIL`**. Zusätzlich
  `-race -count=20` über die vier berührten Pakete: EC=0 (siehe F-3).
- geprüft, ohne Befund: **Traceability.** Betreff trägt `ADR-0082`, keine
  Struktur-ID (`SPEC-*`/`ARC-*`); `make commit-traceability` grün über die
  letzten fünf Commits. Keine neue Kennung, kein neues Präfix, damit MR-000
  unberührt. Keine neue Betreiber-Oberfläche (`CDC_*`, `cdc.*`, Endpunkt) —
  die beiden HIGH-Klassen Handbuch-Zug und Versionshistorie haben kein
  Objekt.

## Eigene Messungen (Exit-Codes direkt, ungepiped)

| Lauf | Exit | Ergebnis |
|---|---|---|
| `git diff --name-only f90c3f4^..f90c3f4` | **0** | 9 Pfade, **0** ohne `_test.go` |
| `git log --oneline 3de7fb1..22cebec \| wc -l` | **0** | **55**, nicht ~60 |
| `git diff --name-only 3de7fb1..22cebec \| grep -v '^docs/'` | **0** | 12 Pfade → **F-4** |
| Profil bei `22cebec` (Container, netzlos) | **0** | 479 ungedeckt / 1903, 74.8 %; Cluster C = **62** |
| Profil bei `f90c3f4` | **0** | 434 ungedeckt / 1903, 77.2 %; Cluster C = **15** |
| `runWALRetentionCheck` in beiden Läufen | **0** | 100.0 % → 87.5 % (±2) — erklärt +49 vs +47 |
| `make coverage-gate` (via `make gates`) | **0** | `OK — Coverage 77.20% erfüllt Schwelle 70%` |
| `make gates` (Log, EC separat gelesen) | **0** | sechs Checks grün; danach `git status --porcelain` leer |
| `make test` (voller Lauf, `-race`, `--network none`) | **0** | 32 × `ok`, 0 × `FAIL` |
| `go test -race -count=20` (vier Pakete) | **0** | 4 × `ok`, keine Frist gefeuert |
| 24 Mutationsproben (siehe Tabelle) | — | 20 rot, 4 grün (davon 1 benannte Negativkontrolle) |
| Mutationsprobe M-26b (ganzes `natsnotify`-Paket) | **0** | grün → **F-1** |
| `grep -n "ListenAndServe" …/net/http/server.go` (GOROOT) | **0** | `shuttingDown()` vor `net.Listen` (Z. 3450-3456) |
| `ls <eintrag>/evidence/ \| wc -l` (7 Einträge) | **0** | 4·6·4·2·3·2 — deckt §8 des Plans |

Der Gate-Lauf und seine Auswertung sind **zwei Schritte**: `make gates` lief
ungepiped in eine Log-Datei, der Exit-Code separat gelesen und erst danach in
einem eigenen Werkzeug-Aufruf ausgewertet; kein Folgekommando hing an einer
Pipe oder an einem Wrapper (§3.9).

## Antwort auf die Schwerpunkte

**(1) Prüf-Kraft statt Zeilen — trägt, mit einem Fund.** 24 Mutationsproben,
20 rot. Die fünf JSON-Formtests tragen ihre Ablehnung an der Eingabeseite und
sind gegen die Gegenrichtung mitgemessen (M-28). Kein Test in diesem Diff
prüft nur „dass Code läuft“: jede Probe mit einer Produktionsmutation färbt
rot — **außer** der einen in F-1. Zwei Zusagen decken Zweige, die nur ein
Test-Writer erreicht (Interceptor-`!ok`, SSE-`Flusher`) — sie sind als solche
im Testkopf benannt und in `ADR-0082` §Konsequenzen („Getter und
Client-/Handler-Stubs prüfen, nicht die Runtime-Interna“) gedeckt; das ist
benannte Doppelprüfung, kein Theater.

**(2) `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (4×) — erfüllt, an
neun Tests nachgeprüft.** Die Bindung hält bei allen fünf
`…UngueltigesJSONEndetMit400` (M-14…M-18: 5 × rot), beim
Interceptor-Negativfall (M-5 rot), beim SSE-`Flusher`-Fall (M-12 rot), beim
gRPC-Sendefehler (M-2 rot) und in `streamv1` (M-19/M-24 rot). Der Diff liefert
damit **kein** fünftes Auftreten; der Zähler bleibt bei **4×** und der
Lese-Schritt der `welle-20`-Closure bleibt aufgerufen. Der einzige Fund
dieser Achse (F-1) liegt **außerhalb** des Buchstabens der Klasse: die
Eingabeseiten-Mutation (Options-Schleife) färbt dort rot — es fehlt das
**mittlere** Kettenglied, nicht die Eingabeseite.

**(3) Keine zeitabhängigen Tests — gehalten, mit benannter Grenze.**
20 Läufe je Paket unter `-race` grün, keine Frist gefeuert; die Frist je Handler (3 s)
liegt um mehr als eine Größenordnung über der Laufzeit **aller** Tests
seines Pakets (gemessen ≈0,29 s je Iteration im größten der vier). In **drei** Fällen ist
die Frist ihrer Form nach dennoch eine Zusicherung (F-3, INFO). Ein dritter
Beleg für `test-integration-retention-timing-flake` ist aus diesem Diff
**nicht** ableitbar: kein Lauf wechselte bei unverändertem Stand das Ergebnis.

**(4) Die Zahlen sind Lauf-Belege — nachgemessen, sie halten.** `make gates`
druckt in meinem eigenen Lauf `77.20 %`; mein eigenes Profil druckt `77.2 %`
und liefert dedupliziert 434 ungedeckt von 1903 (77,19 %; die zweite
Nachkommastelle der gedruckten Zeile ist Go-Rundung). Die paketweise
Reproduktion der `62` aus `ADR-0082` (4b) gilt exakt (29/25/3/5), die
Endzahl `15` ebenfalls (3/10/0/2), der Zuwachs `+47`, der Nenner unverändert.
Keine Zahl im Diff oder im Plan ist ungemessen; die zwei Abweichungen
(+49 vs +47, 77,30 % vs 77,20 %) sind beide auf die dokumentierte
`runWALRetentionCheck`-Schwankung zurückgeführt — in meinen Läufen sichtbar
als 100.0 % → 87.5 % an derselben Funktion. **Offen (Closure, nicht Review):**
Der Diff selbst trägt keine Zahl; die Zahlen leben bisher nur in der
Commit-Message (git-Historie, kein Doku-Träger) — §7 muss sie mit ihrem Lauf
nennen (LP1 verlangt das ohnehin).

**(5) Der bewegte HEAD berührt den Diff nicht — die Aussage über ihn stimmt
nicht.** Keiner der neun Diff-Pfade kommt im Bereich `3de7fb1..22cebec` vor,
und die Zahlen sind am Parent unabhängig reproduziert. Der Bereich enthält
aber die Cluster-B-Testdateien und `slice-090`s Werkzeug-Zug — „alle nur
`docs/`“ ist contra-gemessen (F-4, INFO, kein Träger im Repo).

**(6) Die Mutationsmethode — reproduzierbar, mit einer Ausnahme.** 24 Proben
nachgefahren (Tabelle oben); 20 färben rot. Die Proben reproduzieren die im
Diff **genannten** Mutationsangaben und ergänzen sie um drei eigenständig
gewählte Eingabeseiten-Mutationen (M-5: gültiger Token statt eines beliebigen
Werts an derselben Stelle; M-9: eine Fehlerklasse auf `500` statt aller;
M-28: **jeder** Body abgelehnt — die Gegenrichtung der Gegenprobe). Die
genannten Mutationen halten bis auf zwei: F-2 (Interceptor, dort trägt die genannte
Mutation den Satz nicht) und F-1 (natsnotify, dort ist die genannte zweite
Mutation wirkungslos und die Stelle selbst ungeprüft). `M-6` (Entfernen des
`!ok`-Zweigs) ist wie im Kopf vorhergesagt grün — der Zweig ist
verhaltensgleich und damit zu Recht keine Zusage.

**Was nicht geprüft wurde:** DoD/LP1–LP3 (Verifier), die §6-Risiko-Ausgänge
und die drei Paarungen (Planner-Closure), die Zählung des
Beobachtungs-Registers durch die Closure, und ob die neuen Tests im kalten
CI-Lauf innerhalb des Workflow-Timeouts bleiben (dafür fehlt der reale
Post-Push-Lauf; `AGENTS.md` §3.10 gilt für den betroffenen Workflow, den
dieser Diff **nicht** ändert).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | **0** |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** „Durchreichen behauptet, nur das erste
Kettenglied gebunden“ (F-1) · `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`
(F-2 **und** die zweite Mutationsangabe in F-1 — **zwei Funde, ein Vorgang**:
das Register zählt Vorgänge, `slice-091` trägt **eine** Evidence-Datei und
hebt diesen Eintrag damit von 2× auf **3×**) · „Frist als Zusicherung statt
als Hänge-Schutz“ (F-3) · „Aussage über den Gegenstand ohne ihren
Beleg-Anker“ (F-4, Träger außerhalb des Repos)

## Verdikt

**Merge-blockierend:** ja — **1 MEDIUM** (F-1). HIGH ist **keiner** dabei;
F-1 ist eine Ein-Stellen-Korrektur in einer Testdatei dieses Diffs und keine
Verletzung einer Hard Rule.

**Rückgabe-Pfeil Reviewer → Implementer:** **nötig** (F-1, F-2) — beide
Fundstellen liegen im Diff und in der Verantwortung des Implementers
(`internal/adapters/driven/natsnotify/notify_test.go`,
`internal/adapters/driving/grpc/interceptor_test.go`). F-3 und F-4 liefern
keinen Rückgabe-Pfeil: F-3 ist die benannte Grenze einer Form, F-4 hat keinen
Träger im Repo. Die Finding-Klassen gehen zusätzlich in die Slice-Closure §7
und von dort in den Zähler.

**DoD-Häkchen „Review durchgeführt, Report unter `docs/reviews/` liegt vor“:**
bleibt **offen** — der Slice braucht eine Fixrunde (`.harness/skills/reviewer.md`
§DoD-Checkbox-Nachzug ohne Fixrunde greift nicht).

**Übergabe:** Dieser Report ist ein **Lauf-Beleg** (dieser Diff, dieser Skill,
dieses Modell, dieses Verdikt) und ersetzt keine Verifikation — DoD-/Spec-
Konformität prüft der Verifier separat (Modul 11).
