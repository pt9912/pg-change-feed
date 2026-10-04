# Verifikations-Report: Slice `tls-http-grpc-server` ([`LH-FA-SST-011`](../../spec/lastenheft.md), [`ADR-0150`](../plan/adr/0150-tls-und-mehrfach-token.md)) — 2026-10-04

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD- und Entscheidungs-Konformität
plus Plan-vs-Code-Diff, frischer Kontext, nach Implementierung, Review und
Fixrunde. Kein Reparieren, kein Setzen von DoD-Häkchen.

**Gegenstand:** `git diff c45a3cca HEAD` — `8e3fa08b` (Code), `08a796f2`
(Runner, Handbuch, Plan), `eca155ad` (Review), `aaa3ed3c` (Fixrunde),
`0115485d` (Review-Haken). Plan:
[`slice-tls-http-grpc-server`](../plan/planning/done/slice-tls-http-grpc-server.md);
Review: [`review-slice-tls-http-grpc-server`](review-slice-tls-http-grpc-server.md)
(0 HIGH, F-1/F-2 MEDIUM in der Fixrunde behoben, F-3/F-4 LOW, F-5/F-6 INFO).
Bezug: [`SPEC-034`](../../spec/pflichtenheft.md),
[`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md),
[`ADR-0152`](../plan/adr/0152-zugangsdaten-klasse-elf-schluessel.md)
(Zugangsdaten-Klasse; die zwei TLS-Schlüssel gehören nicht dazu, §2).

## Verdikt

**Bestanden, mit Bedingungen für den Planner (§9).** Jede Liefer-Punkt-Zeile und
jede Gate-/Lauf-Zeile der DoD ist an einem eigenen Lauf bestätigt; keine
Abweichung zwischen Beleg des Implementers und eigenem Lauf. Die Fixrunde
(Runner, Handbuch, Kommentare, `wiring.go`-Kommentar) ist voll mitgeprüft. Eine
Zusage bleibt *hergeleitet* (Runner-Slot-Zähler, §3). Offen bleibt nur
Planner-/Closure-Arbeit.

## 1. Eigene Sensor-Belege (dieser Lauf, Exit direkt, nie durch eine Pipe — §3.9)

Lange Läufe liefen im Hintergrund mit gesichertem Exit (`make …; echo $? > Datei`).
Kein Docker-Hub-Abruflimit aufgetreten. Nach `make test-integration` waren
`docker ps -a` und `docker network ls` frei von `cdc-*`-Resten (nur das
unbeteiligte `gitea`), `git status --porcelain` leer: kein `.pem` im Arbeitsbaum.

| Sensor | Exit | Gedruckter Beleg |
|---|---|---|
| `make image` | 0 | Lauf ohne Fehler (zuerst gefahren) |
| `make test` (nach `make mod-download`, Exit 0) | 0 | 52 Pakete `ok`, keine Zeile `FAIL` |
| `make test-store` | 0 | `db-coverage: OK — DB-Adapter-Coverage 83.03% erfuellt Schwelle 80%` |
| `make test-integration` | 0 | `run-integration-tests: TLS der HTTP-, SSE- und gRPC-Schnittstelle (LH-FA-SST-011) belegt — fünf Starts … PCF-E2009 bzw. PCF-E2010 im Log, ohne Slot; das gültige Paar legte über die Konfigurationsdatei den Slot an; … --healthcheck gesund (Ausgang 0, Zustand healthy), HTTP, gRPC, SSE (change_id=2445-1) und gRPC-Stream (change_id=2448-1) wurden über TLS bedient …; der Klartext-Versuch … HTTP: CLEARTEXT status=400 body="Client sent an HTTP request to an HTTPS server."; gRPC: CLEARTEXT code=Unavailable …; nach der Wiederherstellung ohne Override bedienten beide Adressen wieder im Klartext`; 0 Treffer `--- FAIL`; Schluss `Lauf abgeschlossen — E2E-Abdeckungstabelle aus 21 Go-Zeilen und 56 Bash-Zeilen` |
| `make gates` | 0 | `baseline-verify … 54 Dateien`, `coverage-gate: OK — Coverage 83.10% erfüllt Schwelle 80%`, `d-check: 1663 Datei(en) geprüft, 0 Befund(e)`, `commit-traceability: OK — 5 Commit(s)`, `generated-sync: OK`, `meldungscodes-check: 100 Codes …`, a-check `gesamt: 0 Befund(e)` |
| `make docs-check` | 0 | `d-check: 1663 Datei(en) geprüft, 0 Befund(e)` (mit dem Report erneut, siehe Commit-Zeile am Ende) |
| `make fmt-check` | 0 | `342 Go-Dateien geprüft, alle formatiert` |
| `make a-check` | 0 | `gesamt: 0 Befund(e)`; `git diff c45a3cca HEAD -- .a-check.yml` leer |
| `make handbuch-public-doc-check` | 0 | `keine interne Kennung in 3 Nutzerdokumenten unter docs/user/` |
| `make ausgabe-kennungen-check` | 0 | `keine interne Kennung in Ausgabe-Literalen von 130 Go-Dateien und 4 Skripten` |
| `make meldungscodes-check` | 0 | `100 Codes in Tabelle und Katalog gleich, Quelltext (128 Go-Dateien, 4 Skripte) nur mit Codes der Tabelle` |
| `make doc-trace` | 0 | `83 Anforderung(en), 1 Waise(n).`; Zeile `LH-FA-SST-011 … E2E | ok`; Waise `LH-FA-SST-010` |
| `make doc-immutable RANGE=c45a3cca..HEAD` | 0 | `0 Befund(e)` |
| `make doc-commits RANGE=c45a3cca..HEAD` | 0 | `0 Befund(e)` |
| `make commit-traceability` | 0 | `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `make kommentar-kennungen DIFF=c45a3cca COUNT=1` | 0 | `0` (kein Kandidat) |
| `make suchlauf-nachmessen PLAN=…` | 0 | `12 Zeilen stimmen` (je Zeile soll = ist: 8/9/22/6/5/29, diff 12/11/30/30/52/30) |
| `make examples-csharp` / `make examples-kotlin` | 0 / 0 | beide Exit 0; die Läufe endeten binnen Sekunden, die Schichten kamen weitgehend aus dem Docker-Cache (gleicher Inhalt wie im Lauf der Fixrunde) — *der Bau-Beleg ist ein Cache-Treffer, kein frischer Compile* |

## 2. DoD-Zeile für Zeile

| DoD-Zeile | Befund |
|---|---|
| Liefer-Punkt 1 (a) Konfiguration: Variablen, Datei-Schlüssel, Paar-Prüfung aus **einer** Funktion | bestätigt: `validateTLSPair` in `internal/bootstrap/tls.go`, gerufen aus `ConfigFromEnv` und `mergeConfig` (Diff gelesen); Codes `PCF-E2009`/`PCF-E2010` Klasse `configuration` in `codes.go`; die zwei Schlüssel nicht in `forbiddenFileCredentialKeys` (Test `TestTLSDateiSchluesselGehoerenNichtZurZugangsdatenKlasse` im grünen `make test`; die Rot-Angabe der zugehörigen Plan-Mutation ist **übernommen**, §3) |
| (b) Laden vor jeder Verbindung und vor dem Slot | bestätigt am Code: `Run` ruft `newTLSConfig` (Zeile 631) vor `postgresstorage.New` (636) und `receive.NewStream` (1249); `--healthcheck` lädt nicht (`main.go` ruft nur `ConfigFromEnvAndFile`); Mutation M2 (§3) rot |
| (c) Betrieb: `MinVersion` 1.2 explizit, kein mTLS, Unit-Tests je Fähigkeit | bestätigt: `tls.Config{Certificates, MinVersion: tls.VersionTLS12}`, kein `ClientAuth`/`CipherSuites`/`InsecureSkipVerify` im Produktionscode (`git grep`: nur Treffer in `tls_internal_test.go:242`, dort als Prüfung `NoClientCert`); Klartext-Verhalten gemessen und benannt (`400` bzw. `Unavailable`) |
| (d) Mutationsprobe | an fremder Kopie nachgefahren, §3 |
| Liefer-Punkt 2 Realserver-Phase | bestätigt am eigenen Lauf (§1); Reihenfolge: erst TLS-Erfolg, dann Klartext-Versuch (Runner-Code gelesen); Healthcheck unter TLS Exit 0 und `healthy`; fünf Negativstarts je mit Ausgang, Klasse, Code, Text, ohne `PRIVATE KEY`, ohne Slot; Gegenprobe: Slot entsteht mit gültigem Paar; Wiederherstellung im Klartext |
| `doc-trace` ohne Waise für `LH-FA-SST-011` | bestätigt (§1) |
| Liefer-Punkt 3 Handbuch 1.96 | bestätigt, §6 |
| Gates: `make gates`, `docs-check`, `test`, `test-store`, `test-integration`, `image`, `a-check` | alle Exit 0 (§1) |
| Review, Report liegt vor | bestätigt (Datei vorhanden, F-1/F-2 MEDIUM in `aaa3ed3c`, Fixrunde hier voll geprüft) |
| Doku-Update `harness/README.md` | bestätigt, eine Zeile in der `test-integration`-Zeile; die Aussage „ohne geöffneten Port“ ist entfernt |
| Closure-Notiz mit Lerneintrag; Beobachtungs-Register; §6-Ausgänge; drei Paarungen | **offen**, Planner/Closure (nicht Teil dieser Verifikation) |

## 3. Mutationen, selbst nachgefahren

Kopie des `HEAD`-Baums im Scratchpad (`git archive`), Mutation als
`sed … Datei > Kopie` und `cp` (Block-Verschiebung M2 über `sed -n`-Zeilenbereiche
nach Kopie, kein `-i`), Rücknahme per `cp`, danach `cmp` gegen den Echtbaum gleich.
Instanz: `go test -race -count=1 <Paket>` im Toolchain-Container, `--network none`,
Kopie als `:ro`-Mount (wie das Makefile-Ziel `test`).

| Nr. | Stelle · mutierte Eingabe | Instanz | gesehene Farbe |
|---|---|---|---|
| M1 | `validateTLSPair`, Fall 1 `&&` zu `\|\|` | `internal/bootstrap` | rot (Exit 1): u. a. `TestTLSPaarWirdGelesen`, `TestNewTLSConfigBautEineKonfiguration`, `TestNewTLSConfigMindestversionAmDraht`, dazu mehrere `TestMergeConfig…`/`TestConfigFromEnv…`-Tests |
| M2 | `Run`: `newTLSConfig` hinter `postgresstorage.New` verschoben (Fixrunden-Mutation) | `internal/bootstrap`, `-run TestRunLaedtTLSPaarVorJederVerbindung` | rot: `tls_internal_test.go:421: Run mit nicht ladbarem Paar: Fehlerklasse storage [PCF-E5001] … erwartet ErrTLSPairUnusable ohne Datenbankfehler` |
| M3 | HTTP `serve`: Verzweigung `TLSConfig != nil` zu `false` | `internal/adapters/driving/http` | rot: `TestServerTLSBedientJedeFaehigkeit`, `TestServerTLSLehntKlartextClientAb` |
| M4 | gRPC `New`: `cfg.TLSConfig != nil` zu `false` | `internal/adapters/driving/grpc` | rot: `TestServerTLSBedientStreamUndVerwaltung`, `TestServerTLSLehntKlartextClientAb` |
| M5 | `LoadX509KeyPair(certFile, certFile)` | `internal/bootstrap` | rot: `TestNewTLSConfigBautEineKonfiguration`, `…MindestversionAmDraht`, `…LiestDieKetteDerZertifikatsdatei`, `…PruftAblaufUndNamenNicht`, `TestRunLaedtTLSPaarVorJederVerbindung` |
| M6 | `MinVersion` auf `tls.VersionTLS10` | `internal/bootstrap` | rot: `TestNewTLSConfigBautEineKonfiguration`, `TestNewTLSConfigMindestversionAmDraht` |
| M7 | `Run`: Ladeaufruf auf `newTLSConfig("", "")` (Paar verworfen) | `internal/bootstrap` | rot: `TestRunLaedtTLSPaarVorJederVerbindung` |

Jede gefahrene Mutation färbte sich rot, an der angegebenen Instanz; HTTP und
gRPC je am eigenen Server gefahren (keine Herleitung). **Nicht gefahren:** die
Plan-Mutationen „Fall 2 `case false`“, „`MinVersion` entfernt“, `mergeConfig`
ohne Datei-Wert und `forbiddenFileCredentialKeys` ergänzt — deren Rot-Angaben im
Plan sind **übernommen**, nicht nachgemessen. **Nicht gefahren, *hergeleitet*:**
die Runner-Mutation „Laden hinter `NewStream`“ (Slot-Zähler des Runners). Sie
verlangt einen kompletten `make test-integration`-Lauf (rund 15 Minuten) mit
`make image-mutation`; der Beleg des Slot-Zählers an *dieser* Mutation fehlt
weiter. Die Zusage „kein Slot bei Ladefehler“ ruht damit auf dem Code-Aufbau
(Laden vor `NewStream`, gelesen) und dem Unit-Test hinter M2.

## 4. Sicherheitsaussagen

- `internal/bootstrap/tls.go`: `MinVersion: tls.VersionTLS12` explizit, kein
  mTLS (`ClientAuth` ungesetzt, Default `NoClientCert`, ein Test prüft es), keine
  `CipherSuites`, kein `InsecureSkipVerify` im Produktionscode. Gelesen und per
  `git grep` bestätigt.
- Kein Schlüsselmaterial in Fehlern: der Ladefehler nennt Pfade und die
  Ursache von `tls.LoadX509KeyPair` (Fehlertext ohne Schlüsselbytes);
  `TestNewTLSConfigLadefehlerEndenMitConfiguration` prüft ohne `PRIVATE KEY`, der
  Runner prüft jedes Negativ-Log ohne `PRIVATE KEY` (grün im eigenen Lauf).
  `git grep 'PRIVATE KEY'` im Produktionscode: nur `certgen` (schreibt in das
  Temp-Verzeichnis des Runners) und Tests.
- Ablauf und Name des Zertifikats werden nicht geprüft: die Grenze steht im
  Handbuch (Aufzählungspunkt „prüft beim Start weder das Ablaufdatum noch den
  Namen“), im Plan (§1, §2) und im Test `TestNewTLSConfigPruftAblaufUndNamenNicht`.
- Kein `.pem` im Arbeitsbaum (Runner-Zähler 0, `git status --porcelain` leer
  nach dem Lauf).
- `chmod 0644` auf die Wegwerf-Schlüssel im Runner (F-4): im Kommentar mit Grund
  benannt (Container als `nonroot`); Wegwerf-Material im Temp-Verzeichnis.

## 5. Reihenfolge und Verhalten ohne TLS

`Run`: Zeile 631 `newTLSConfig`, 636 erster `postgresstorage.New`, 1249
`receive.NewStream`; die zwei Server binden danach. HTTP `Start` bindet jetzt
synchron per `net.Listen`; ohne TLS bleibt das Verhalten unverändert: alle
übrigen Runner-Phasen liefen im eigenen Lauf grün im Klartext, `make test` 52
Pakete `ok`. Benannte Grenze (3) des Plans (`Addr` ohne Wert bindet einen
Zufallsport) trifft die Verdrahtung nicht, weil sie `Start` nur mit gesetzter
Adresse ruft.

## 6. Handbuch 1.96

- Kennungsfrei (`make handbuch-public-doc-check` Exit 0).
- Clients-Absatz nach der Fixrunde: nur Gemessenes. Nachgemessen: `git grep -i -E
  'tls|ssl|certificate'` über die Optionen-Dateien von C#, Kotlin und Python:
  kein Treffer (Exit 1); die gRPC-Beispiele sind im Klartext fest verdrahtet
  (`insecure.NewCredentials()` in `examples/grpc-client/main.go:102`,
  `usePlaintext()` in `Main.kt:61`, C# laut Kommentar mit h2c-Switch,
  `Program.cs:51`). Der Schluss „verbinden deshalb nicht mit einem TLS-Server“
  ist *hergeleitet* (kein Beispiel gegen einen TLS-Server gefahren; gestützt
  durch den gemessenen Klartext-Versuch des Wegwerf-Clients: `Unavailable`).
  Für Client-Pakete und HTTP-/SSE-Clients steht „nichts zugesagt“: wahr, nichts
  behauptet.
- „ohne geöffneten Port“: im lebenden Baum (ohne Reviews, `done/`, Beobachtungen,
  Baseline, ADRs) kein Treffer; der einzige Treffer von `kein Port`
  (`benutzerhandbuch.md:1477`, „kein Server, kein Port“) gehört zur
  deaktivierten HTTP-API und ist älter als der Slice (`git blame`: `10a12219`).
- Exit-Codes: 2 bei `PCF-E2009` (`ConfigFromEnvAndFile` endet in `main.go` vor
  `Run` mit `os.Exit(2)`), 1 bei `PCF-E2010` (`Run` liefert den Fehler,
  `os.Exit(1)`): wahr, und im eigenen Runner-Lauf je Start über `bf_expect`
  gemessen.
- „beide Fälle beenden den Prozess vor dem Aufbau der Datenbank-Verbindungen“
  stimmt mit dem Code überein (Zeile 631 vor 636) und trägt der Unit-Test hinter
  M2; am Container ist nur „kein Slot“ gemessen, nicht „keine Verbindung“.

## 7. Spec-Treue

[`SPEC-034`](../../spec/pflichtenheft.md): Variablen, Datei-Schlüssel, Paar, PEM-Datei mit Zertifikat
und Kette. „Zertifikat und Kette“ ist nur an `crypto/tls`
(`TestNewTLSConfigLiestDieKetteDerZertifikatsdatei`, zwei Glieder) gemessen,
nicht an einem Client mit Zwischenzertifikat; der Plan benennt die Grenze (§3,
Benannte Grenzen 4), die Aussage ist begrenzt. Das Handbuch sagt „das
Zertifikat und dahinter seine Kette“ ohne Zusage zur Vertrauenskette des
Clients.

## 8. Schichten und Umfang

`git diff c45a3cca HEAD --stat -- sdks .a-check.yml compose.yaml` ist leer: kein
SDK-Zug, `.a-check.yml` unverändert. Die Adapter nehmen eine fertige
`*tls.Config` entgegen, gebaut einmal in `internal/bootstrap`. SDK-TLS bleibt
dem Folge-Slice `sdk-tls-optionen` überlassen.

## 9. Offene Punkte und Vorschläge für den Planner

**§6-Ausgänge (Vorschlag, je genau ein Ausgang):**

| Risiko | Vorschlag |
|---|---|
| Verhalten gegenüber Klartext-Client ungeprüft | entfallen: gemessen (`net/http` `400` mit Text, `grpc-go` `Unavailable`), Unit-Tests und eigener Runner-Lauf |
| Negativ-Beleg „kein Klartext“ trifft nicht bereiten Server | entfallen: der Klartext-Versuch folgt dem belegten TLS-Erfolg (Runner gelesen, eigener Lauf grün); die Prüfung „Port nach beendetem Container“ ist entfernt |
| Zertifikat im Image oder Repo | entfallen: Runner-Zähler 0 und `git status --porcelain` leer am eigenen Lauf |
| Prozess liest die Zertifikatsdateien nicht | entfallen: `chmod 0644`, Start mit gültigem Paar im Container gelaufen (`healthy`) |
| Spec-Lücke zu „ungültig“ (Ablauf, Name) | weiter offen als benannte Spec-Lücke, mit Register- oder Lerneintrag (Handbuch und Test benennen sie) |
| Zwei Konfigurationswege lesen verschieden | entfallen: eine Funktion `validateTLSPair`, Mutation M1 färbt beide Wege; Container-Start im Runner über Umgebung **und** Datei |
| SDK-/Beispiel-Aussage breiter als ihre Messung | entfallen: Handbuch auf Gemessenes zurückgenommen; der Schluss zu den Beispielen steht als hergeleitet (§6 oben) |
| Runner-Phase stört die Folge-Phasen | entfallen: alle Folge-Phasen im eigenen Lauf grün, Container ohne Override wiederhergestellt |

**Register- und Lerneintrag-Vorschläge (nur Vorschläge):**

- Beobachtung `ready-ist-nicht-verbunden` (zwei Evidenz-Dateien am
  Planungsstand): dieser Slice hält die Reihenfolge ein, die Klasse ist hier
  **vermieden**, nicht eingetreten; keine neue Evidenz-Datei nötig.
- Neue Klasse (Review F-2): „Beleg gegen beendeten Container“ — eine Prüfung
  („kein Klartext-Server“), die einen schon beendeten Container liest und an
  keiner Eingabe rot werden kann; in der Fixrunde entfernt. Sie ist verwandt mit
  `negativtest-ohne-bindung-an-seine-eingabe` (25 Evidenz-Dateien, gezählt mit
  `ls` am Arbeitsbaum); ob sie ein eigenes Verzeichnis oder eine weitere
  Evidenz-Datei dort wird, entscheidet der Planner; ein Zähler wird nicht von
  Hand gesetzt.
- `inplace-textwerkzeug-am-repo-trotz-nutzerregel`: in diesem Lauf nicht
  beobachtet; kein Eintrag.
- Lerneintrag-Kandidaten (Plan §5): das gemessene Klartext-Verhalten von
  `net/http`/`grpc-go`; die benannte Spec-Lücke „Ablauf des Zertifikats“;
  die Mutationslücke „Runner-Slot-Zähler“ (hergeleitet).

**Release-Folge:** keine (kein SDK-Zug).

**Adressen für Folgearbeit (Vorschlag):**

- Beispiel-Nachzug (TLS-fähige gRPC-Beispiele in Go, C#, Kotlin): eigener
  Slice in `open/`, Arbeitsname `examples-grpc-tls`; Träger der Adresse ist die
  Closure von `tls-http-grpc-server`.
- `sdk-tls-optionen`: wartet weiter auf die Folge-Anforderung zu
  [`ADR-0150`](../plan/adr/0150-tls-und-mehrfach-token.md) §Konsequenzen; Adresse ist
  die Closure-Notiz §7 „Folge-Slices“, wie der Plan sie schon nennt.

**Slice C `otlp-metrik-export`:** [`slice-otlp-metrik-export`](../plan/planning/open/slice-otlp-metrik-export.md)
liegt in `open/`; seine zwei Verweise auf
[`slice-tls-http-grpc-server`](../plan/planning/done/slice-tls-http-grpc-server.md)
lösen auf (`make docs-check` Exit 0). Beim `git mv` nach `done/` sind sie
nachzuziehen.

**Bedingungen:**

1. Closure-Notiz (§7), Register-Zeile, §6-Ausgänge und die drei Paarungen vom
   Planner schreiben; keine DoD-Häkchen von mir.
2. Der Satz zu den gRPC-Beispielen („verbinden deshalb nicht mit einem
   TLS-Server“) bleibt *hergeleitet*; kein Anspruch auf einen Aufruf.
3. Die Runner-Slot-Mutation (Laden hinter `NewStream`) bleibt *hergeleitet*,
   es sei denn der Planner will den Lauf (ca. 15 Minuten, `make image-mutation`
   und `make image-mutation-rm`).
4. Beim `git mv` nach `done/` die Links aus Slice C und aus diesem Report
   nachziehen (`make docs-check`).
