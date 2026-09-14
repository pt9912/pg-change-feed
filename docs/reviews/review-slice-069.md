# Review-Report: slice-069 — 2026-09-14

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten);
DoD-/Spec-Konformität ist Verifier-Aufgabe und **nicht** Gegenstand dieses
Reports.

**Gegenstand:** `slice-069`, Diff `067bf3d..af742cb` — vier Commits
(`7c02d15` Adapter-Grundgerüst, `9f4817d` Lebenszyklus-Tests, `7066916`
Image-Digest, `af742cb` DoD-Häkchen), 20 Dateien (+1499/−26).
Der im Auftrag genannte Elter-Commit `4085f18` ist der reine
`next→in-progress`-Move, **trägt aber zwei fremde Planner-Commits nach sich**
(`a31597d`, `067bf3d`) — der Slice-Diff beginnt deshalb erst bei `067bf3d`
(siehe F-8).

**Skill:** `.harness/skills/reviewer.md` @ `f07ba3d` (letzte Schärfung
2026-09-13, unverändert seitdem).
**Modell:** deepseek-v4.1-flash:cloud[1m] · **Datum:** 2026-09-14.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-069`
  (vollständig, inkl. der vier Plan-Nachzüge, §6 und §8)
- [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) (vollständig,
  alle sechs Teilfragen — insbesondere Teilfrage 2 Option A *verworfen*,
  Teilfrage 3 Zustellsemantik, Teilfrage 4 Auth, Teilfrage 5 Platzierung,
  Teilfrage 6 Bootstrap —, §Konsequenzen, §Fitness Function,
  §Slice-Schnitt-Empfehlung)
- [`SPEC-020`](../../spec/pflichtenheft.md) (neu, `spec/pflichtenheft.md:339-359`),
  [`SPEC-002`](../../spec/pflichtenheft.md), `SPEC-018`
- [`LH-FA-SST-008`](../../spec/lastenheft.md) (Happy Path / Boundary /
  Negative wörtlich), `LH-FA-REA-001`, `LH-FA-CON-003`/`005`
- [`ADR-0057`](../plan/adr/0057-http-grpc-api.md) (Token-Klassen),
  [`ADR-0055`](../plan/adr/0055-nats-change-notification-wecksignal.md)
  (Wecksignal-Abgrenzung), [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)
  (Digest-Semantik), [`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
- `AGENTS.md` §3.1/§3.2/§3.7/§3.8/§3.9/§4/§5/§6;
  `harness/conventions.md` (`MR-000` ID-Schema)
- `docs/plan/planning/observations/BEO-PGC/` — `rollen-verdrahtung` (4×,
  eingetreten), `adapter-fehler-ausgang` (2×, weiter offen),
  `a-check-null-abdeckung` (3×, verkörpert),
  `slice-chronik-in-code-kommentar` (verkörpert), `commit-traceability-kein-vorab-hook`
  (3×, Ausgang geplant), `report-nackte-id-ohne-link`

**Eigene Sensor-Läufe** (Exit-Code jeweils ungepiped in einem eigenen Schritt
ermittelt, `AGENTS.md` §3.9):

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make gates` | **0** | baseline-verify, docs-check (540 Dateien, 0 Befunde), commit-traceability (5 Commits, Betreffs ohne Struktur-ID), a-check (0 Befunde; Hinweis: 2 Dateien in keiner Schicht, unverändert), coverage-gate (48,30 % ≥ 35 %) |
| `make test` | **0** | `internal/adapters/driven/grpcstream` und `internal/adapters/driving/grpc` grün, kein `FAIL` im Lauf |
| `make a-check` | **0** | 0 Befunde gegen den vollen Baum |
| `make proto-generate` | **0** | Determinismus-Nachprüfung, s. Negativbefund 6 — zweiter Lauf byte-identisch |

---

## Findings

### F-1 — `SPEC-020` pinnt nur die halbe Zustellsemantik: die Blockade-Hälfte fehlt

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md)
  §Folgepflicht („neuer `SPEC-*`-Eintrag für … die Stream-Semantik") ·
  §Teilfrage 3 („ohne Empfänger wird die Nachricht verworfen (kein Puffer)")
- `pfad`: `spec/pflichtenheft.md:356` (Zeile *Zustellgarantie*), gegen
  `internal/adapters/driven/grpcstream/broadcaster.go:88-107` und
  `internal/application/port/outbound/changestream.go:31-39`
- `befund`: Die Zeile beschreibt den Fall „Consumer empfängt nicht" als
  „verpasst die Nachricht ersatzlos"; der implementierte `Publish` verwirft in
  diesem Fall aber nicht, sondern hält die Übergabe an einen **registrierten,
  gerade nicht lesenden** Empfänger an, bis er liest, sich abmeldet oder `ctx`
  endet (im Plan-Nachzug `:164-170` und im Code benannt, in `SPEC-020` nicht).
  Der Godoc des Ports sagt ebenfalls zu, ein „zum Verteilungszeitpunkt nicht
  empfangsbereiter" Abonnent erhalte den Change „nicht nachgeliefert" — die
  Implementierung wartet auf genau diese Empfangsbereitschaft. Die vom
  [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) an `SPEC-020`
  delegierte Semantik-Aussage ist damit unvollständig gegenüber dem, was der
  Aufrufer in `slice-070` vorfindet.
- `verifizierbar`: nein — Prosa-/Vertrags-Abgleich, kein Gate-Gegenstand
- `klasse`: „Zustellsemantik nur halb gepinnt"

### F-2 — Fehlender Negativtest der Token-Konfigurationsgrenze im neuen Adapter

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md)
  §Teilfrage 4 · `.harness/skills/reviewer.md` MEDIUM „fehlende Negativtests
  bei neuem öffentlichem Vertrag"
- `pfad`: `internal/adapters/driving/grpc/server_test.go:49-73`
  (`startTestServer`), gegen
  `internal/adapters/driving/http/server_test.go:194-211`
- `befund`: Jeder Testserver des neuen Adapters trägt `TokenReader`/`TokenAdmin`
  gesetzt; die Grenze „ein **leer konfiguriertes** Token trifft nie ein
  Aufruf-Token" — der real aufgetretene `slice-059`-Fehlerpfad, den die
  Sichtung in Plan §8 (`:294-297`, `BEO-PGC/rollen-verdrahtung`) ausdrücklich
  als Erinnerung führt — hat hier keinen Testfall, obwohl `classifyToken` im
  `grpc`-Paket neu vorliegt. Der HTTP-Adapter deckt genau diese Grenze mit
  `TestClassifyTokenLeereKonfiguration` ab. Am **Code** ist die Grenze
  eingehalten (Negativbefund 3) — ungeprüft ist sie trotzdem.
- `verifizierbar`: ja — `make test`; der Fall läuft erst mit einem Testfall,
  kein bestehendes Gate fordert ihn
- `klasse`: „fehlender Negativtest der Token-Konfigurationsgrenze"

### F-3 — Token-Klassifikation steht zweimal in zwei Driving-Adaptern, ohne deklarierten Gewinner

- `kategorie`: LOW
- `quelle`: Maintainability · `.a-check.yml` (keine `adapters→adapters`-Kante)
- `pfad`: `internal/adapters/driving/grpc/interceptor.go:31-55`, gegen
  `internal/adapters/driving/http/middleware.go:14-38`
- `befund`: `role`, die drei Konstanten und `classifyToken` sind im neuen
  Adapter eine zweite, bis auf Kommentare gleichlautende Fassung derselben
  sicherheitsrelevanten Zuordnung; die Layer-Regel verbietet den direkten
  Import der einen Adapter-Hälfte durch die andere, ein gemeinsamer Ort ist
  also nicht trivial — die Doppelung ist aber an keiner der beiden Stellen als
  solche benannt, und eine Änderung an einer Hälfte erreicht die andere nicht.
- `verifizierbar`: nein
- `klasse`: „sicherheitsrelevante Klassifikationslogik in zwei Adaptern dupliziert"

### F-4 — `SPEC-020` zitiert für die Nachrichtenfelder `SPEC-002`, das diese Namen nicht führt

- `kategorie`: LOW
- `quelle`: `SPEC-002` (`spec/pflichtenheft.md:167-181`) · Maintainability
- `pfad`: `spec/pflichtenheft.md:354`
- `befund`: Die Zeile nennt „dieselben Felder wie der Domain-Change
  (`SPEC-002`)" und listet dann `old_image`/`new_image`; `SPEC-002` führt
  `old_data`/`new_data` und beschreibt das Speicherschema (`cdc.change`),
  nicht den Domain-Typ. Gemeint ist offenkundig `internal/domain/model.Change`
  (`OldImage`/`NewImage`, `change.go:46-56`) — die zitierte Stelle trägt den
  Satz nicht.
- `verifizierbar`: nein — Kennungs-Abgleich in der Doku, kein Gate-Gegenstand
- `klasse`: „Spec-Zitat verfehlt die benannte Stelle"

### F-5 — `Subscribe`-Vertrag verbietet das Lesen nach `cancel`; der paketeigene Test liest danach

- `kategorie`: INFO
- `quelle`: Maintainability (`AGENTS.md` §3.7 — ein Kommentar beschreibt, was da ist)
- `pfad`: `internal/adapters/driven/grpcstream/broadcaster.go:52-57`, gegen
  `internal/adapters/driven/grpcstream/broadcaster_test.go:129-134`
- `befund`: Der Godoc sagt zu, „nach der Abmeldung darf der Aufrufer den
  gelieferten Kanal nicht mehr lesen — er wird nicht geschlossen";
  `TestAbmeldenIstIdempotent` liest den Kanal nach der Abmeldung
  (`:130-134`). Folgenlos, weil der Kanal nach dem Entfernen des Empfängers
  nicht mehr bedient wird — der Satz trägt seine Begründung („nicht
  geschlossen") aber nicht selbst.
- `verifizierbar`: nein
- `klasse`: „Vertrags-Satz ohne Deckung im Verhalten"

### F-6 — Kein Sensor gegen Drift zwischen `.proto`-Quelle und committetem Erzeugnis

- `kategorie`: INFO
- `quelle`: [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md)
  §Konsequenzen (Docker-only-Codegenerierung als Folgepflicht)
- `pfad`: `Makefile:101-110` (`proto-generate`), Erzeugnis in
  `internal/adapters/driving/grpc/streamv1/`
- `befund`: Der erzeugte Go-Code ist committet und wird von `make test`/`make
  image` mitkompiliert; ein Edit am `.proto` ohne anschließendes
  `make proto-generate` fällt keinem Gate auf — der Generator läuft nur, wenn
  ihn jemand aufruft. Das Repo kennt committete Erzeugnisse ohne Sync-Gate
  bereits (`tools/schema/plan.yaml`, `tools/schema/down.sql` aus
  `make schema-rollout`), deshalb kein Defekt dieses Slice — die Bindung
  bleibt aber offen.
- `verifizierbar`: nein — genau das ist der Befund
- `klasse`: „committetes Generator-Erzeugnis ohne Sync-Sensor"

### F-7 — Der generierte Code zählt in die Coverage-Messung

- `kategorie`: INFO
- `quelle`: [`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
  (Coverage-Gate, `-coverpkg ./internal/...`)
- `pfad`: `internal/adapters/driving/grpc/streamv1/changestream.pb.go`
- `befund`: Das Paket `streamv1` liegt unter `internal/...` und wird von
  `-coverpkg` mitgemessen; im Gate-Lauf dieses Reviews steht es mit
  `0.0% of statements` in der Ausgabe (die Accessor-Funktionen erscheinen
  einzeln darunter). Die Gesamtschwelle bleibt unberührt (48,30 % ≥ 35 %).
  Ob Generator-Erzeugnisse in die Messung gehören, ist eine Frage der
  Reifestufen-Bindung, nicht dieses Diff.
- `verifizierbar`: ja — `make gates`, Ausgabe des `coverage-gate`-Schritts
- `klasse`: „Generator-Erzeugnis in der Coverage-Messung"

### F-8 — Der im Auftrag genannte Range trägt fremde Commits; der Arbeitsbaum trägt eine fremde Änderung

- `kategorie`: INFO
- `quelle`: Maintainability (Rollen-Trennung, Baseline-Regelwerk `v6.5.0` ·
  `regelwerk/modul-08-agentenrollen.md` §Rollen-Sequenz für einen Slice)
- `pfad`: `docs/plan/planning/open/slice-076-coverage-gate-reifestufe-40.md:28`
- `befund`: `4085f18..HEAD` enthält `a31597d` und `067bf3d` (Coverage-Reifestufe
  [`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md),
  `slice-076`-Datei, Architect-Verdikt) — Planner-Arbeit, die nicht zu diesem
  Slice gehört; der Slice-Diff ist `067bf3d..HEAD`, und dort ist die
  `slice-076`-Datei in **keinem** der vier Commits (nachgeprüft per
  `git diff --name-only`). Im Arbeitsbaum dieses Laufs liegt dieselbe Datei
  jedoch uncommittet modifiziert (`Verantwortlich: —` → `pt9912`); sie ist nicht
  Gegenstand dieses Reports und wird von ihm nicht mitcommittet.
- `verifizierbar`: ja — `git diff --name-only 067bf3d..HEAD` bzw.
  `git status --short`
- `klasse`: „Fremdänderung im Arbeitsbaum während des Reviews"

## Negativbefunde

- **geprüft, ohne Befund: `ChangeNotificationPort`/`natsnotify` unverändert
  (Auftrags-Punkt 1).** `git diff --name-only 067bf3d..HEAD --`
  `internal/application/port/outbound/changenotification.go`
  `internal/adapters/driven/natsnotify` liefert keine Ausgabe — und über den
  weiteren Range `4085f18..HEAD` ebenso wenig. [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md)
  §Teilfrage 2 Option A („Payload-Parameter ergänzen") ist damit nicht nur
  verworfen, sondern im Diff nachweislich unausgeführt; der neue
  `ChangeStreamPort` ist ein getrenntes Interface (`changestream.go`), der
  Port-Zuschnitt folgt Option B.
- **geprüft, ohne Befund: `Broadcaster`-Verträge am Code (Auftrags-Punkt 3).**
  Beide Kanäle ungepuffert (`broadcaster.go:59`, `make(chan *model.Change)`
  ohne Kapazität), Abmelden idempotent über `sync.Once` (`:66-74`), `changes`
  wird an keiner Stelle geschlossen (auch nicht in `cancel`), die Abmeldung
  trägt `done`. **Die Begründung gegen das Schließen trägt:** Schließt `cancel`
  den Kanal, trifft das einen `Publish`, der `sub.changes <- change` in seinem
  `select` gewählt hat — ein Senden auf einen geschlossenen Kanal panikte, und
  `select` schützt davor nicht. `done` als Abmelde-Signal ist der Träger, der
  diesen Fall ohne Panik auflöst; die Tests üben ihn real aus
  (`broadcaster_test.go:92-108`). Die Paket-Coverage der drei Funktionen
  `New`/`Subscribe`/`Publish` liegt im Gate-Lauf bei 100 %.
- **geprüft, ohne Befund: leeres/nicht konfiguriertes Token kann den Stream
  nicht öffnen (Auftrags-Punkt 2).** `classifyToken`
  (`interceptor.go:44-55`) bricht bei `token == ""` ab und verlangt für beide
  Klassen einen **nicht leeren** konfigurierten Wert; `credentialToken`
  (`:61-74`) liefert für fehlenden Eintrag, `values[0]` ohne `Bearer `-Vorsprung
  und für `Bearer ` mit leerem Rest jeweils `""`. Ein ungesetztes
  `CDC_API_TOKEN_READER`/`_ADMIN` bei gesetztem `CDC_GRPC_ADDR` lässt damit
  **keinen** Wert passieren (fail-closed, symmetrisch zum HTTP-Adapter). Die
  Grenze ist ungetestet (F-2), aber eingehalten.
- **geprüft, ohne Befund: `Unauthenticated` für beide Negativ-Hälften am
  laufenden Adapter (Auftrags-Punkt 2).** Drei Testfälle treffen die Klasse
  am RPC (`server_test.go:117-148`: keine Metadata, unbekannter Wert, falsche
  Wertform) — der dritte belegt genau die von `SPEC-020` gepinnte
  `Bearer `-Form. Die Happy-Hälfte ist real, nicht nur als Fehlerpfad: ein Change
  erreicht den Client mit allen zehn Feldern inkl. Row Images
  (`server_test.go:150-200`), für `reader` **und** `admin`, wie
  `ADR-0060` §Teilfrage 4 es fordert.
- **geprüft, ohne Befund: Schichtung und Importe (Auftrags-Punkt 6).**
  Importlisten gelesen: `grpcstream` → `port/outbound`, `domain/model`;
  `grpc` → `grpc/streamv1`, `port/outbound`, `domain/model` plus die
  gRPC-Bibliotheken; die Tests ebenso. **Kein** Adapter-Paket importiert ein
  anderes, das lokale `changeSubscriber`-Interface (`server.go:29-31`) wird
  strukturell von `*grpcstream.Broadcaster` erfüllt — die Zuordnung prüft der
  Compiler an der Verdrahtungsstelle (`wiring.go:689`), nicht erst der Test.
  `.a-check.yml` ist **nicht** im Diff (nachgeprüft per
  `git diff --name-only -- .a-check.yml`, leer) und `make a-check` läuft mit
  0 Befunden — [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md)
  §Teilfrage 5 Punkt 2 und §Fitness Function („`.a-check` unverändert") sind
  eingehalten.
- **geprüft, ohne Befund: Docker-only-Codegen und Determinismus
  (Auftrags-Punkt 4).** Kein Host-`protoc`/`buf`: `protoc` und die beiden
  Plugins kommen ausschließlich aus der neuen Dockerfile-Stufe `proto`
  (`Dockerfile:33-36`), deren Basis `deps` der bereits digest-gepinnte
  `golang:1.27-alpine@sha256:cf6fca66…` ist; die Plugins sind auf
  `v1.36.12`/`v1.6.2` gepinnt, `protobuf-dev` auf `31.1-r1`. Der Generator läuft
  über `--user "$(PROTO_RUN_USER)"` (`id -u:id -g`, dasselbe Muster wie
  `D_MIGRATE_RUN_USER`) in den Bind-Mount, `--network none`, **nicht** als
  root. **Determinismus nachgestellt:** zweiter `make proto-generate`-Lauf
  liefert für beide Erzeugnisse byte-identische Dateien
  (`sha256:be49f3dd…`/`sha256:2419c153…` vor und nach dem Lauf,
  `git status --short` auf `streamv1/` leer) — die Behauptung des Implementers
  trägt. Die neue Stufe steht **vor** `coverage`/`build`/`runtime`, der
  Default-Build (`make image`, ohne `--target`) trifft weiterhin die letzte
  Stufe (`runtime`) und wird von der Stufe nicht angefasst.
- **geprüft, ohne Befund: neue Abhängigkeiten sind die MVS-Minima, keine
  unbegründeten Sprünge (Auftrags-Punkt 5).** Gegen die `go.mod`-Dateien der
  Module geprüft: `google.golang.org/grpc` `v1.83.2` verlangt `x/net v0.58.0`,
  `x/sync v0.22.0`, `x/sys v0.47.0`, `genproto/googleapis/rpc` `…3dc84a4a5aaa`;
  `x/net v0.58.0` verlangt `x/crypto v0.55.0` und `x/text v0.41.0` — genau die
  im Diff gehobenen Stände. `google.golang.org/protobuf v1.36.12` ist die
  direkte, bewusste Setzung und trifft die Version des gepinnten
  `protoc-gen-go`. Zwei neue **direkte** Abhängigkeiten, kein weiterer
  Fremdbaum — [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md)
  §Konsequenzen („zweite und dritte direkte Nicht-PostgreSQL-Abhängigkeit") ist
  eingehalten; `go 1.27` in `go.mod` blieb unangetastet.
- **geprüft, ohne Befund: `SPEC-020` gegen den Code (Auftrags-Punkt 7).**
  Dienst und RPC: `cdc.stream.v1.ChangeStream` /
  `StreamChanges` mit `returns (stream Change)`
  (`streamv1/changestream_grpc.pb.go:28`, `changestream.proto:37-43`) — deckt
  sich mit der Zeile *Protokoll/Dienst* und *RPC*
  (`spec/pflichtenheft.md:351-352`). Die zehn Nachrichtenfelder, der leere
  Request und die Granularität „eine Nachricht je Zeilen-Change"
  (`server.go:142-155`, `changestream.proto:17-33`) decken sich mit den Zeilen
  *Request*, *Nachricht* und *Granularität* (`:353-355`; zur Feldnamen-Quelle
  siehe F-4). Die Metadata-Wertform `Bearer <token>`
  (`interceptor.go:25`, `:70-73`) deckt sich mit der Zeile
  *Authentifizierung* (`:358`) und mit `SPEC-018`s HTTP-Header-Form. Die
  Aktivierung über `CDC_GRPC_ADDR` (`wiring.go:108-115`, `:678-700`) deckt
  sich mit der Zeile *Aktivierung* (`:359`) und
  [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) §Teilfrage 6
  Option B.
- **geprüft, ohne Befund: additive Verdrahtung und keine Capture-Anbindung.**
  `ConfigFromEnv` liest `CDC_GRPC_ADDR` als optionale Größe
  (`wiring.go:278`), der Adapter startet nur bei gesetzter Adresse in einer
  eigenen Goroutine mit dem Muster des HTTP-Adapters (Log-Fehler statt
  Abbruch, `Shutdown` + `Wait` beim Beenden, `:766-770`); der neue
  `ConfigFromEnv`-Test prüft beide Hälften (`wiring_test.go:151-171`). Ein
  `ChangeStreamPort` wird **nicht** an den `CaptureService` gereicht
  (`grpcstream.New()` steht ausschließlich im `apigrpc.Config`), die
  Capture-Integration ist damit wie in Plan §1 ausgeschlossen — kein
  stiller Scope-Creep nach vorn. `internal/application/usecase/**` ist im Diff
  nicht enthalten.
- **geprüft, ohne Befund: Kommentar-Disziplin per `grep` gegen die geänderten
  `.go`/`.proto`-Dateien (Auftrags-Punkt 10).** Muster
  `slice-[0-9]+|welle-[0-9]+` über den Produktionscode
  (`grpcstream/`, `driving/grpc/server.go`, `interceptor.go`,
  `port/outbound/changestream.go`, `proto/`, `wiring.go`) ergibt **keinen**
  Treffer; die zwei Treffer im Umfeld (`server.go:7` „`spec/architecture.md`
  §1", `wiring_test.go:125` „von `ADR-0057`/`slice-059`") sind ein
  Rang-Zeiger bzw. eine **vorbestehende**, nicht von diesem Diff geänderte
  Testfall-Provenienz. Kein Treffer der Klasse
  `BEO-PGC/slice-chronik-in-code-kommentar`; die neuen Godocs nennen Zusage,
  Kopplung, Abgrenzung oder Grenze und tragen `ADR-*`/`LH-*`/`SPEC-*`-Zeiger.
- **geprüft, ohne Befund: kein Suppression-Pfad.** `grep` auf
  `nolint|noqa|SuppressMessage` über die neuen Pakete und den Port: kein
  Treffer (`AGENTS.md` §3.2).
- **geprüft, ohne Befund: Traceability und ID-Schema.** Die vier Commits
  tragen `LH-FA-SST-008`/`ADR-0060` im Betreff, kein `SPEC-*`/`ARC-*`;
  `make commit-traceability` meldet 5 Commits ohne Struktur-ID im Betreff
  (Exit 0). Die im Diff neu vergebenen Kennungen (`SPEC-020`) folgen `MR-000`;
  die Umbenennung von der im Plan ursprünglich genannten Kennung auf `SPEC-020`
  ist als Plan-Korrektur `:20-23` dokumentiert und trifft die freie Kennung
  (`SPEC-019` ist von `slice-066` belegt).
- **geprüft, ohne Befund: DoD-Häkchen `af742cb` (Auftrags-Punkt 8).** Der
  Commit ändert **eine** Zeile — `- [ ] make gates grün.` → `- [x]`. Kein
  weiterer DoD-Punkt, insbesondere nicht die Review-Zeile, die Closure-Notiz,
  die Register- oder die Risiko-Zeilen; kein Übergreifen in eine andere Rolle.
  In `7c02d15` gesetzt sind die umsetzungs- und gate-tragenden Zeilen; die
  Reconciliation-Zeile trägt ihre Entfall-Begründung im selben Satz.
- **geprüft, ohne Befund: `harness/image-hash.txt` ist ein echter Lauf-Beleg.**
  Der committete Digest `sha256:20dc7366bc0a…` stimmt exakt mit der ID des
  lokal gebauten `ghcr.io/pt9912/pg-change-feed:dev` überein, dessen
  `Created` (`2026-09-14T18:39:42+02:00`) zwischen dem Test-Commit (`18:38:37`)
  und dem Digest-Commit (`18:40:08`) liegt — der Beleg stammt aus einem realen
  `make image` dieses Zugs, die Begründung im Commit-Betreff trifft zu
  (Build-Kontext geändert, [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)).
- **geprüft, ohne Befund: `harness/README.md`/`AGENTS.md`-Doku-Update.**
  Beide neuen Zeilen beschreiben, was das Ziel tut (Docker-only, gepinnte
  Stufe, Bind-Mount, Aufrufer-uid, committeter Erzeugnis-Code, kein Gate) und
  stehen in den Tabellen, in denen die übrigen Werkzeug-Ziele stehen; das
  Target existiert (`.PHONY`, `Makefile:102`), also keine halluzinierte Zeile
  (`AGENTS.md` §4).
- **geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` unberührt — und
  unberührt-lassen trägt.** Die Begründung des Plans (`:179-183`) ist
  faktisch gedeckt: das Handbuch führt `CDC_CAPTURE_DSN`, `CDC_ADMIN_DSN`,
  `CDC_READER_DSN`, `CDC_SOURCE_ID`, `CDC_PUBLICATION`, `CDC_SLOT`,
  `CDC_TABLES` — keine der HTTP-Adapter-Variablen (`CDC_HTTP_ADDR`,
  `CDC_API_TOKEN_*`). Kein HIGH der Klasse
  `BEO-PGC/handbuch-versionshistorie-uebersprungen` (der Diff ändert die Datei
  nicht), die Operator-Erreichbarkeit bleibt
  [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md)
  §Slice-Schnitt-Empfehlung Punkt 3 überlassen.
- **geprüft, ohne Befund: die §3-Nachzüge bleiben im Scope.** Die neuen Zeilen
  benennen die tatsächlich entstandenen Pfade (auch `streamv1/`, `go.mod`,
  `image-hash.txt`) und ergänzen nichts, was nicht im Diff liegt; die
  Festlegungen, die der Plan als „Implementer-Entscheidung" führt
  (Protobuf-Pfad, Wertform, `Publish`-Blockade, kein Unary-Interceptor, Test
  ohne Broadcaster-Import, Handbuch, `harness/mk`-Reduktion), sind im Diff
  einzeln auffindbar. Kein Liefer-Punkt-Zuwachs.
- **geprüft, ohne Befund: `internal/bootstrap/wiring.go`,
  `wiring_test.go`, `internal/application/port/outbound/changestream.go`,
  `proto/cdc/stream/v1/`, `Dockerfile`/`Makefile`,
  `harness/README.md`/`AGENTS.md`/`harness/image-hash.txt`,
  Slice-Plan `slice-069`** — über die oben benannten Punkte hinaus kein
  Befund.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 2 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Zustellsemantik nur halb gepinnt ·
fehlender Negativtest der Token-Konfigurationsgrenze · sicherheitsrelevante
Klassifikationslogik in zwei Adaptern dupliziert · Spec-Zitat verfehlt die
benannte Stelle · Vertrags-Satz ohne Deckung im Verhalten · committetes
Generator-Erzeugnis ohne Sync-Sensor · Generator-Erzeugnis in der
Coverage-Messung · Fremdänderung im Arbeitsbaum während des Reviews

## Verdikt

**Merge-blockierend:** ja — 0 HIGH, 2 MEDIUM (F-1, F-2). Beide brauchen eine
Rückkante Reviewer → Implementer: F-2 einen Testfall am neuen Adapter, F-1
eine Ergänzung der von
[`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) delegierten
Stream-Semantik (Pflichtenheft-Zeile `spec/pflichtenheft.md:356`, dazu der
Godoc des Ports). **Kein DoD-Nachzug dieses Reports:** Weil diese Rückkante
existiert, bleibt die Zeile „Review durchgeführt, Report unter `docs/reviews/`
liegt vor" im Slice-Plan **offen** und wird regulär bei Schritt 21 des
Implementer-Workflows nachgezogen (`.harness/skills/reviewer.md`
§DoD-Checkbox-Nachzug ohne Fixrunde greift genau dann, wenn **keine** Fixrunde
kommt).

**Zum Broadcaster-Risiko für `slice-070` (Auftrags-Punkt 3) — es ist echt und
durch die dortige Anbindung nicht von selbst entschärft.** Die Anbindung, die
[`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) §Teilfrage 2
zeichnet, ruft `Publish` **synchron** in der Best-Effort-Kette nach `ACK
Source` auf und isoliert nur den **Fehler** (Log-Warnung, kein Eintrag in den
Rückgabewert) — nicht die **Zeit**. Der ausgelieferte `Publish` hält die
Übergabe an einen registrierten, gerade nicht lesenden Empfänger an, und
gebunden ist dieser Halt allein durch `ctx` (Ende des Aufrufs) bzw. die
Abmeldung; das `ctx` der Capture-Schleife lebt so lange wie der Prozess. Die
Kette ist damit real unbeschränkt: ein Client, der den Stream geöffnet hat und
nicht liest (pausierter/flow-control-blockierter Leser, dessen `stream.Send`
seinerseits hängt), hält über den ungepufferten Kanal den `Publish`-Aufruf und
damit die Capture-Schleife an — nicht den `ACK` **dieser** Transaktion, aber
den `Receive` **jeder weiteren**, also Wachstum des Capture-Lag bis zum
Stillstand, und zusätzlich Head-of-Line-Blocking über **alle** Empfänger
hinweg. Entschärft wird das nur durch eine Schranke an der Aufrufstelle
(Endlichkeit von `ctx`) oder durch eine andere Zustellsemantik — beides ist
eine Entscheidung, die `slice-070` treffen und sichtbar machen muss; der
Parameter `ctx` trägt sie bereits, die ADR verlangt keinen Puffer und steht
ihr nicht entgegen. Der Plan trägt genau dieses Risiko in §6 (`:239-244`,
Ausgang „wird bei Closure zugewiesen") und der Paket-Kommentar benennt die
Semantik — **deshalb kein Finding gegen diesen Code**, aber der offene Punkt
gehört in `slice-070` geschlossen und in der Zwischenzeit über F-1 in der
Semantik-Zeile geführt.

**Was dieser Lauf unabhängig bestätigt.** Auftrags-Punkt 1 ist belegt, nicht
geglaubt: keine Datei von `ChangeNotificationPort`/`natsnotify` im Diff,
über keinen der beiden Ranges. Auftrags-Punkt 2 ist am Code gehalten (ein
leeres Token kann den Stream nicht öffnen) und am laufenden Adapter in drei
Negativ- und zwei Erfolgsfällen belegt. Auftrags-Punkt 4 trägt die
Determinismus-Behauptung — nachgestellt, zweiter Lauf byte-identisch.
Auftrags-Punkt 5 ist sauber: nur die MVS-Minima der beiden neuen direkten
Abhängigkeiten, keine Sprünge. Auftrags-Punkt 6 hält die Schichtung ohne
`adapters→adapters`-Import und ohne `.a-check.yml`-Änderung. Auftrags-Punkt 8
ist ein sauberer Ein-Zeilen-Nachzug. Auftrags-Punkt 9: die Fremddatei
`slice-076` ist in keinem Slice-Commit — sie liegt aber uncommittet modifiziert
im Arbeitsbaum (F-8) und ist von diesem Report unberührt geblieben.

**Übergabe:** Findings gehen an den Implementer (F-1 zusätzlich als
Doku-/Semantik-Punkt an die Planer-Seite, weil die Zeile im Pflichtenheft
liegt); die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und
von dort in den Zähler. Dieser Report selbst ist ein **Lauf-Beleg** (Audit:
dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) — er wird über Läufe
hinweg nicht wieder gelesen. Der Report ersetzt keine Verifikation — DoD-/Spec-
Konformität prüft der Verifier separat (Baseline-Regelwerk `v6.5.0` ·
`regelwerk/modul-11-*.md`; anderes Prüf-Artefakt, anderer Eingabe-Kontext).

---

## Fixrunde (2. Lauf) — 2026-09-14

**Gegenstand:** die zwei Fix-Commits `0867835` (Broadcaster nicht-blockierend
mit begrenzter Empfangswarteschlange) und `2f6882c` (`SPEC-020`-Präzisierung),
beide real gelesen, nicht aus den Commit-Messages übernommen. Zusätzliche
Prüfgrundlage dieses Laufs:
[`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
(`Supersedes` [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md)
**nur** die Puffer-Klausel; Index-Vermerk in `docs/plan/adr/README.md`,
`ADR-0060`s Textkörper unangetastet — das Muster von
`ADR-0059`/`ADR-0065` ist eingehalten) und das Architect-Verdikt
[`architect-verdict-publish-blockiert-capture-pfad.md`](architect-verdict-publish-blockiert-capture-pfad.md).
Alle fünf umsetzbaren Findings (F-1…F-5) sind damit über den
Rollenweg Reviewer → Architect → Implementer zurückgekommen; die zwei
INFO-Befunde F-6/F-7 bleiben begründet liegen (unten).

**Eigene Sensor-Läufe** (Exit-Code je eigener Schritt, ungepiped, mechanisch
konditioniert):

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make test` (Stand nach der Fixrunde) | **0** | 28 Pakete `ok`, kein `FAIL` |
| `make test` + Mutation A (Sende-Zweig blockierend) | **2** | rot genau an der Kernzusage — `TestNichtLesenderAbonnentHaeltPublishNichtAn` und `TestWarteschlangeVerwirftUeberlaufOhneAndereZuStauen`, je `2.00s`-Frist |
| `make test` + Mutation C (nur der Leer-Token-Frühausstieg entfernt) | **0** | grün — das ist Befund F-10 |
| `make test` + Mutation B (Frühausstieg **und** beide `!= ""`-Wächter entfernt) | **2** | rot in **beiden** Hälften des F-2-Tests: `TestClassifyTokenLeereKonfigurationTrifftKeinToken` (`interceptor_test.go:17`) und `TestStreamChangesLeereTokenKonfigurationEndetMitUnauthenticated` am realen Adapter |
| `make test` (nach Rücknahme aller Mutationen) | **0** | Arbeitsbaum leer |
| `make gates` (nach Rücknahme) | **0** | baseline-verify, docs-check, commit-traceability, a-check, coverage-gate |

**Rücknahme und Blatt-Identität (Mutation A und B):** nach `git checkout --`
stimmt `git hash-object` beider Dateien mit `git rev-parse HEAD:<pfad>` überein
(`broadcaster.go` `36672aa9…`, `interceptor.go` `49f98fee…`),
`git status --porcelain` ist leer, kein Rest im Baum.

### Verdikt je Finding

| Finding | Verdikt | Beleg |
|---|---|---|
| F-1 (MEDIUM) | **behoben** | `spec/pflichtenheft.md:356-357`, `changestream.go:31-45` |
| F-2 (MEDIUM) | **behoben** | `interceptor_test.go:15-38`, `server_test.go:49-64` |
| F-3 (LOW) | **behoben** | `interceptor.go:44-50`, `http/middleware.go:26-32` |
| F-4 (LOW) | **behoben** | `spec/pflichtenheft.md:354` |
| F-5 (INFO) | **behoben** | `broadcaster.go:56-65`, `broadcaster_test.go:167-190` |
| F-6 (INFO) | **Begründung geteilt** — Klasse bleibt offen, gehört ins Register | unten |
| F-7 (INFO) | **Begründung geteilt** — Klasse bleibt offen, gehört ins Register | unten |

**F-1 — behoben.** Die zwei Zeilen
(`spec/pflichtenheft.md:356-357`) und der Port-Godoc (`changestream.go:31-45`)
tragen jetzt beide Hälften und decken sich mit dem Code: nicht-blockierend
(`broadcaster.go:110-115`, `select` mit `default`), begrenzte
Empfangs-Warteschlange je Abonnent (`queueCapacity = 64`, `:36`, `:67`),
Drop-Newest (der **eintreffende** Change fällt weg, die eingereihten bleiben in
Reihenfolge — der Test `TestWarteschlangeVerwirftUeberlaufOhneAndereZuStauen`
prüft beides), keine Zustellgarantie (Abmeldung verwirft die Warteschlange,
kein Replay, keine Subskriptions-übergreifende Position). **Die
*Erzeuger-Blockade*-Zeile ist nach der Ergänzung korrekt:** `Publish` hat genau
zwei Fehlerausgänge — `ErrChangeStream` beim fehlenden Change (`:97-99`) und
den rohen Kontext-Fehler beim bereits beendeten `ctx` (`:100-102`); die Zeile
nennt beide und die Einordnung des zweiten als *ungültiger Aufruf* trifft den
Code, weil der Aufruf vor jeder Verteilung abbricht. Der Wortlaut des Verdikts
wurde dabei überholt — der Implementer hat den Kontext-Ausgang ergänzt, den
Frage 3 des Verdikts noch offen ließ; `2f6882c` zieht die Zeile nach, statt die
Lücke zu lassen.

**F-2 — behoben.** `interceptor_test.go:15-22` übt die reine Funktion, `:29-38`
zusätzlich den **realen** `bufconn`-Adapter mit leerer Token-Konfiguration
(`startTestServerMitTokenKonfiguration(t, newFakeSubscriber(), "", "")`,
`server_test.go:56-64`) über drei Metadata-Formen (fehlend, `Bearer ` mit leerem
Rest, `Bearer beliebig`) — alle drei enden mit `Unauthenticated`. Die
Sensitivität ist nicht behauptet, sondern belegt: Mutation B färbt genau diesen
Test rot (Tabelle oben).

**F-3 — behoben.** Beide Fassungen tragen denselben Gegenverweis
(`interceptor.go:45-50`, `middleware.go:27-32`) samt der tragenden Begründung
(„Das `.a-check.yml`-Schichtenmodell führt keine `adapters→adapters`-Kante")
und der Kopplungs-Anweisung („beide Fassungen sind zusammen zu ändern"). **Die
Begründung trägt:** `.a-check.yml` führt die Kanten `adapters→ports` und
`adapters→domain` und ist im Diff unberührt; ein gemeinsamer Ort außerhalb der
Adapter-Schichten existiert nicht. Der HTTP-Kommentar hat dabei seinen
`slice-059`-Verweis gegen die Invariante getauscht — kein Chronik-Verlust, der
Vorgang steht in `git` und im Register.

**F-4 — behoben.** `spec/pflichtenheft.md:354` nennt jetzt den Domain-Typ
(`model.Change`, `OldImage`/`NewImage`) und die Datei — die zitierte Stelle
trägt den Satz.

**F-5 — behoben.** Der Godoc (`broadcaster.go:56-65`) sagt jetzt, was gilt
(„Nach der Abmeldung werden dem Aufrufer keine weiteren Changes zugestellt; er
liest den gelieferten Kanal nicht weiter") und nennt den Grund des
Nicht-Schließens ausdrücklich („ein Schließen träfe einen gleichzeitig
laufenden Sende-Versuch in `Publish` und panikte, `select` schützt davor
nicht"); `TestAbmeldenIstIdempotent` (`:167-190`) liest den Kanal nach `cancel`
nicht mehr und prüft stattdessen den entfernten Empfänger. Godoc und Test sind
deckungsgleich — der Widerspruch ist weg, nicht umformuliert.

**F-6 — Begründung geteilt, Klasse bleibt offen.** Kein Sync-Gate für
`.proto`↔`.pb.go` ist die richtige Entscheidung **jetzt**: ein solches Gate
kostete in jedem `make gates`-Lauf einen Codegen-Durchlauf (Docker-Build +
`protoc`), und das Repo führt dieselbe Klasse bereits ohne Sensor
(`tools/schema/plan.yaml`/`down.sql`). Der richtige Moment ist der Eintritt
eines zweiten Konsumenten der Schnittstelle (`slice-070`/`slice-071`) — bis
dahin ist die Lücke zulässig, aber **nicht verschwiegen**: die Klasse
„committetes Generator-Erzeugnis ohne Sync-Sensor" gehört bei der
Slice-Closure in §7 und von dort ins Register, damit sie nicht als erledigt
gilt.

**F-7 — Begründung geteilt, Klasse bleibt offen.** Die Messdefinition zu
ändern (generierten Code aus `-coverpkg` nehmen) wäre selbst eine
Schwellen-/Sensor-Entscheidung und gehört nicht in einen Fixlauf; das
`coverage-gate` ist bootstrap-aware (`ADR-0054`) und misst als Gesamtzahl, die
Generat-Coverage trägt derzeit rund um die Schwelle keinen Ausschlag
(48,30 % ≥ 35 %). Scharf wird die Frage erst an der Endstufe, weil
Generator-Erzeugnisse die Gesamtzahl strukturell drücken, ohne testbar zu sein —
der Träger dafür ist die Reifestufen-Linie (`slice-076` in `next/`). Auch hier:
Klasse bei der Closure in §7 ins Register, nicht still schließen.

### Neue Findings dieses Laufs

### F-9 — `SPEC-020` nennt die Verwerfungsrichtung nicht

- `kategorie`: INFO
- `quelle`: [`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
  Festlegung 3 („Explizite Verwerfungsregel: Drop-Newest … bereits eingereihte
  Changes werden nicht verdrängt") · Maintainability
- `pfad`: `spec/pflichtenheft.md:356` gegen
  `internal/adapters/driven/grpcstream/broadcaster.go:88-90` und
  `internal/application/port/outbound/changestream.go:33-36`
- `befund`: Die Zeile *Zustellgarantie* sagt „je Abonnent trägt der
  `Broadcaster` eine begrenzte Empfangs-Warteschlange, deren Überlauf verworfen
  wird" — die Richtung (der eintreffende Change fällt weg, der eingereihte
  bleibt) nennt sie nicht, während Port-Godoc und Paket-Kommentar „Drop-Newest"
  ausdrücklich führen und die ADR sie als eigene Festlegung führt. Der Wortlaut
  stammt aus dem Architect-Verdikt (Frage 3), ist also nicht vom Implementer
  abgewichen, sondern dort schon offen.
- `verifizierbar`: nein — Prosa-Abgleich, kein Gate-Gegenstand
- `klasse`: „Verwerfungsrichtung im Spec-Text nicht benannt"

### F-10 — Die benannte Mutation „Leer-Token-Frühausstieg" färbt den Lauf nicht rot

- `kategorie`: INFO
- `quelle`: Maintainability (Mutations-Beleg als Nachweis der Test-Sensitivität)
- `pfad`: `internal/adapters/driving/grpc/interceptor.go:51-59`
- `befund`: Der Frühausstieg `if token == "" { return roleNone }` ist mit den
  `!= ""`-Wächtern der beiden konfigurierten Klassen **doppelt** gedeckt: jede
  Hälfte allein hält die Klasse „leeres konfiguriertes Token" ab — Entfernen
  des Frühausstiegs allein lässt die Suite grün (selbst gestellt, Exit 0),
  Entfernen **beider** Hälften färbt sie rot (selbst gestellt, Exit 2). Die
  Verteidigung im Code ist damit nicht schmal, aber ein Beleg, der nur den
  Frühausstieg entfernt, weist keine Sensitivität nach. Der *tragende* Fall ist
  das Paar — und das ist belegt (F-2).
- `verifizierbar`: ja — `make test` mit der jeweiligen Mutation; beide Läufe in
  der Tabelle oben mit Exit-Code
- `klasse`: „Mutations-Beleg greift die redundante Hälfte"

## Summary (Fixrunde)

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Verwerfungsrichtung im Spec-Text nicht
benannt · Mutations-Beleg greift die redundante Hälfte

## Verdikt (Fixrunde)

**Merge-blockierend:** nein — **keine Fixrunde mehr nötig.** F-1 bis F-5 sind
behoben (Verdikt-Tabelle oben, jeder Punkt am Code/Test nachgeprüft, nicht an
der Commit-Message), die beiden neuen Findings sind INFO ohne erwartete Aktion,
F-6/F-7 bleiben begründet liegen. Die Kette läuft damit über
Reviewer → Architect (Folge-ADR [`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
+ Verdikt) → Implementer (Fixrunde) → Reviewer und schließt hier; der nächste
Rollenwechsel ist Reviewer → Verifier (DoD-/Spec-Konformität, Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-11-*.md`).

**Zur Kernzusage, unabhängig nachgestellt:** Der Test
`TestNichtLesenderAbonnentHaeltPublishNichtAn` pulsiert `queueCapacity+1` Changes
(`broadcaster_test.go:122`) und ist genau deshalb sensitiv — mit Mutation A
(blockierender Send in einen **64er**-Puffer) fällt er nach der 65. Übergabe in
die 2-Sekunden-Frist (`broadcaster_test.go:45`), rot an der zugesagten Stelle.
Die Begründung des Implementers trifft zu: ein einzelner Change hätte auch in
einen blockierenden Send in einen 64er-Puffer gepasst und den Test grün
gelassen. Damit ist `ADR-0066` Festlegung 1 („Senke darf die Quelle nie
aufhalten — strukturell, nicht per Schranke") testgetragen, nicht nur
kommentiert.

**DoD-Nachzug:** Mit diesem Verdikt ist die Rückkante geschlossen — die Zeile
„Review durchgeführt, Report unter `docs/reviews/` liegt vor" in §2 des
Slice-Plans ist in demselben Commit wie dieser Vermerk auf `[x]` gezogen
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde). Nur diese
eine Zeile; Closure-Notiz, Beobachtungs-Register und die drei Paarungen bleiben
offen (Planner-Arbeit).
