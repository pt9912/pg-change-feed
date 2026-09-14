# Slice 069: gRPC-Streaming-Adapter-Grundgerüst

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-19.

**Bezug:** [LH-FA-SST-008](../../../../spec/lastenheft.md) (Haupt-Bezug),
[LH-FA-REA-001](../../../../spec/lastenheft.md),
[LH-FA-CON-003](../../../../spec/lastenheft.md),
[LH-FA-CON-005](../../../../spec/lastenheft.md) (Fortsetzungs-Referenz,
Boundary-Kriterium),
[ADR-0060](../../adr/0060-grpc-streaming-mechanismus.md),
[ADR-0066](../../adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
(Zustellsemantik des `Broadcaster` — superseded `ADR-0060`s Puffer-Klausel).

**Berührte Spec-Stellen:** `SPEC-020` (neu anzulegen durch diesen Slice —
Protobuf-Nachrichtenschema, RPC-Methodenname, Stream-Semantik; von
[ADR-0060](../../adr/0060-grpc-streaming-mechanismus.md) als Folgepflicht
benannt). **Plan-Korrektur 2026-09-14:** Der Slice-Plan nannte dafür
zunächst `SPEC-019`; diese Kennung hat inzwischen `slice-066` belegt
(Feldform von `cdc.administration_request`), deshalb rückt dieser Slice auf
die nächste freie Kennung `SPEC-020`.

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Rolleninhaber: Planner-Lauf, 2026-09-14). **Datum:** 2026-09-14.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Das gRPC-Streaming-Adapter-Grundgerüst aus
[ADR-0060](../../adr/0060-grpc-streaming-mechanismus.md) steht: neuer
Outbound Port `ChangeStreamPort`, der In-Prozess-`Broadcaster`
(`internal/adapters/driven/grpcstream/`), der gRPC-Driving-Adapter samt
Auth-Interceptor (`internal/adapters/driving/grpc/`), die Docker-only
Protobuf-/buf-Toolchain-Anbindung und die additive
`CDC_GRPC_ADDR`-Bootstrap-Verdrahtung — jeweils isoliert getestet, aber noch
nicht an den Capture-Pfad angeschlossen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Capture-Integration** (`CaptureService.WithChangeStream`, die
  tatsächliche Publish-Kette `… → ACK Source → Notify → Stream-Publish`) —
  Folge-Slice `slice-070` übernimmt es; dieser Slice liefert nur den Port
  und den Broadcaster, der Aufrufer entsteht erst dort
  ([ADR-0060](../../adr/0060-grpc-streaming-mechanismus.md)
  §Slice-Schnitt-Empfehlung).
- **gRPC-Beispiel-Client und E2E-Rundlauf gegen den laufenden
  Feed-Container** — Folge-Slice `slice-071` übernimmt es; ohne
  `slice-070`s Publish-Kette gäbe es hier noch kein reales Ereignis zu
  belegen.
- **HTTP/SSE-Endpunkt** — anderer Vorgang: `slice-072`
  ([ADR-0061](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md)) nutzt den
  hier gebauten `Broadcaster` als zweiten Abonnenten, ist aber ein eigener
  Liefer-Punkt in einem anderen Adapter-Paket
  (`internal/adapters/driving/http/`).
- **Stream-internes Replay und tabellen-granulare Filterung** — Bestand
  bleibt bewusst stehen: Beide sind laut
  [ADR-0060](../../adr/0060-grpc-streaming-mechanismus.md)
  §Re-Evaluierungs-Trigger bewusst vertagt und bräuchten je eine eigene
  Folge-ADR, falls ein konkreter Bedarf entsteht — kein Gegenstand dieser
  Welle (siehe `welle-19.md` §6).
- **Änderungen an `ChangeNotificationPort`/`natsnotify`** — Schicht-
  Abgrenzung: `ChangeNotificationPort` bleibt unverändert
  ([ADR-0060](../../adr/0060-grpc-streaming-mechanismus.md) Teilfrage 2,
  Option A explizit verworfen); dieser Slice fasst keine Datei dieses
  Pakets an.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [LH-FA-SST-008](../../../../spec/lastenheft.md) (Teilaspekt
      Authn/Authz-Boundary des Server-Grundgerüsts) erfüllt, Test
      referenziert: `internal/adapters/driving/grpc` — ein
      Stream-Öffnungsversuch ohne oder mit unbekanntem
      `authorization`-Metadata-Wert endet mit gRPC-Status `Unauthenticated`;
      ein gültiges `reader`- oder `admin`-Token öffnet den Stream
      erfolgreich (Fitness Function aus
      [ADR-0060](../../adr/0060-grpc-streaming-mechanismus.md)).
- [x] [ADR-0060](../../adr/0060-grpc-streaming-mechanismus.md) Teilfrage
      1/2/4/5/6 umgesetzt: neuer Outbound Port
      `internal/application/port/outbound/changestream.go`
      (`ChangeStreamPort`), Driven-Adapter
      `internal/adapters/driven/grpcstream/` (`Broadcaster`:
      `Subscribe()`/`Publish()`, nebenläufigkeitssicher, kein Puffer,
      Fire-and-Forget-Regressionstest gegen `make test` — ein `Publish`
      ohne aktiven Subscriber blockiert nicht und liefert keinen Fehler),
      Driving-Adapter `internal/adapters/driving/grpc/`-Server-Grundgerüst
      samt Auth-Interceptor, additive Bootstrap-Verdrahtung
      `CDC_GRPC_ADDR` (No-Op bei fehlender Adresse), Docker-only
      Protobuf-/buf-Build-Stufe für die Code-Generierung.
- [x] Fixrunde (Review `slice-069` F-1/F-2, [ADR-0066](../../adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)):
      `Publish` nicht-blockierend mit begrenzter Empfangs-Warteschlange je
      Abonnent (Drop-Newest, Kanal nie geschlossen), Paket-/Funktions-Godoc
      und `internal/application/port/outbound/changestream.go`-Godoc auf die
      Semantik aus [ADR-0066](../../adr/0066-broadcaster-begrenzte-empfangswarteschlange.md);
      Regressionstest „ein registrierter, nicht lesender Abonnent hält
      `Publish` nicht an"; Negativtest der Token-Konfigurationsgrenze
      (`classifyToken` bei leer konfiguriertem Token, analog
      `TestClassifyTokenLeereKonfiguration`).
- [x] `spec/pflichtenheft.md` erhält den neuen Eintrag
      [`SPEC-020`](../../../../spec/pflichtenheft.md) (konkretes
      Protobuf-/Nachrichtenschema, RPC-Methodenname, Stream-Semantik) —
      Folgepflicht aus
      [ADR-0060](../../adr/0060-grpc-streaming-mechanismus.md).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8);
      Report: `docs/reviews/review-slice-069.md`, Verdikt 0 HIGH, 2 MEDIUM /
      2 LOW / 4 INFO, F-1…F-5 in der Fixrunde behoben, damit ohne weitere
      Fixrunde geschlossen (Nachzug nach `.harness/skills/reviewer.md`
      §DoD-Checkbox-Nachzug ohne Fixrunde).
- [x] Doku-Update `harness/README.md` §Sensors/Werkzeuge und `AGENTS.md`
      §4, falls ein neues `make`-Ziel für die Protobuf-/buf-Codegenerierung
      entsteht.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. *(Entfällt: `docs/plan/planning/reconciliation.md` existiert in diesem GF-Repo nicht.)*
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/application/port/outbound/changestream.go` | neu | `ChangeStreamPort` — `Publish(ctx, change) error`, dazu der Sentinel `ErrChangeStream`; Fixrunde F-1: `Publish`-Godoc auf die nicht-blockierende Zustellsemantik nachgezogen |
| `internal/adapters/driven/grpcstream/broadcaster.go` | neu | `Broadcaster`: `Subscribe()`/`Publish()`, In-Prozess-Fan-out; Fixrunde: begrenzte Empfangs-Warteschlange je Abonnent (`queueCapacity`), nicht-blockierender Send mit Drop-Newest |
| `internal/adapters/driven/grpcstream/broadcaster_test.go` | neu | Fire-and-Forget-Regressionstest; Fixrunde: Regression „ein registrierter, nicht lesender Abonnent hält `Publish` nicht an“ plus Warteschlangen-Test (begrenzt, Drop-Newest, Reihenfolge) |
| `internal/adapters/driving/grpc/server.go` | neu | gRPC-Server-Grundgerüst, lokales `changeSubscriber`-Interface, Domain↔Protobuf-Übersetzung |
| `internal/adapters/driving/grpc/interceptor.go` | neu | Auth-Interceptor (Metadata-Token-Prüfung); Fixrunde F-3: zweite Fassung von `classifyToken`/`role` als benannte Duplikation dokumentiert |
| `internal/adapters/driving/grpc/server_test.go` | neu | Unauthenticated-/Erfolgs-Pfad-Tests gegen `bufconn` mit lokalem Fake-Subscriber, dazu `Start`/`Shutdown`-Lebenszyklus und Bind-Fehlerpfad |
| `internal/adapters/driving/grpc/interceptor_test.go` | neu (Fixrunde F-2) | Negativtest der Token-Konfigurationsgrenze: `classifyToken` bei leer konfiguriertem Token, fail-closed am laufenden Adapter |
| `internal/adapters/driving/http/middleware.go` | update (Fixrunde F-3) | Gegenverweis auf die zweite `classifyToken`-Fassung im `grpc`-Adapter; Chronik-Verweis im Kommentar durch die Invariante ersetzt (`AGENTS.md` §3.7) |
| `proto/cdc/stream/v1/changestream.proto` | neu | Protobuf-Schema für Change-Nachricht + RPC-Methode (konkretisiert den Plan-Platzhalter `proto/`) |
| `internal/adapters/driving/grpc/streamv1/changestream.pb.go`, `changestream_grpc.pb.go` | neu | erzeugter Go-Code (committet, damit `make test`/`make image` ohne Codegen laufen) |
| `Dockerfile` | update | neue Stufe `proto` für protoc + `protoc-gen-go`/`protoc-gen-go-grpc` (Docker-only) |
| `Makefile` | update | neues Ziel `proto-generate` (+ `PROTO_IMAGE`/`PROTO_RUN_USER`) |
| `internal/bootstrap/wiring.go` | update | additive `CDC_GRPC_ADDR`-Verdrahtung (Server-Start, noch ohne `CaptureService`-Anschluss — folgt in `slice-070`) |
| `internal/bootstrap/wiring_test.go` | update | ConfigFromEnv-Test: `CDC_GRPC_ADDR` bleibt optional |
| `spec/pflichtenheft.md` | update | neuer Eintrag `SPEC-020` (§2), externe-Verträge-Zeile (§6), Historie-Zeile (§7); Fixrunde F-1/F-4: Erzeuger-Blockade-Zeile um den Kontext-Fehlerausgang präzisiert, Feldnamen-Quelle in `SPEC-020` von `SPEC-002` auf den Domain-Typ `model.Change` berichtigt |
| `harness/README.md`, `AGENTS.md` | update | neues `make`-Ziel `proto-generate` in der Werkzeuge-/Gate-Tabelle |
| `go.mod`, `go.sum` | update | neue direkte Abhängigkeiten `google.golang.org/grpc`, `google.golang.org/protobuf` |
| `harness/image-hash.txt` | update | Digest-Beleg nach `make image` (Build-Kontext geändert, `harness/README.md` §Werkzeuge) |

**Implementer-Entscheidungen und -Abweichungen (Plan-Nachzug im selben Lauf):**

- **Protobuf-Pfad konkretisiert:** `proto/cdc/stream/v1/changestream.proto`;
  der erzeugte Go-Code liegt unter
  `internal/adapters/driving/grpc/streamv1/` (Paket `streamv1`) — der
  `.proto`-Platzhalter aus §3 war ausdrücklich Implementer-Entscheidung.
- **Metadata-Wertform:** der Interceptor liest den Wert des
  `authorization`-Metadata-Eintrags in der Form `Bearer <token>` —
  dieselbe Wertform wie der HTTP-Header (`SPEC-018`); `ADR-0060`
  Teilfrage 4 lässt die Wertform ausdrücklich der Spezifikation
  (`SPEC-020` pinnt sie).
- **`Publish`-Blockade-Semantik (korrigiert durch Architect-Verdikt,
  [ADR-0066](../../adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)):**
  Der erste Entwurf übergab ungepuffert und hielt an einen registrierten,
  gerade nicht lesenden Empfänger an (`ADR-0060` Teilfrage 3, „kein
  Puffer"). [ADR-0066](../../adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
  löst die Puffer-Klausel ab: `Publish` übergibt **nicht-blockierend** an
  eine begrenzte Empfangs-Warteschlange je Abonnent (Drop-Newest), der
  Erzeuger hält nie an; ohne registrierten Empfänger unterbleibt die
  Verteilung (Fitness Function bleibt). Die Fixrunde setzt das samt Godoc
  und Regressionstest um — siehe §2 und §6.
- **Kein Unary-Interceptor:** `ChangeStream` trägt ausschließlich
  Streaming-RPCs; ein Unary-Interceptor hätte keinen Aufruf zu schützen.
- **Driving-Test ohne Broadcaster-Import:** der Whitebox-Test im
  `grpc`-Paket benutzt einen lokalen Fake-Subscriber statt
  `*grpcstream.Broadcaster` — ein Adapter-Paket darf kein anderes
  importieren (`.a-check.yml`, `adapters`-Layer hat keine
  `adapters→adapters`-Kante). Der reale Broadcaster wird in seinem eigenen
  Paket getestet.
- **`docs/user/benutzerhandbuch.md` unberührt** — Schicht-Abgrenzung: die
  ENV-Variablen des HTTP-Adapters (`CDC_HTTP_ADDR`/`CDC_API_TOKEN_*`) sind
  dort ebenfalls nicht dokumentiert; die Operator-Dokumentation des
  Streamings gehört zur Operator-Erreichbarkeit (`slice-071`/`slice-072`),
  nicht zum Adapter-Grundgerüst.
- **Reduktion:** der Plan-Platzhalter `harness/mk/*.mk` bleibt unberührt —
  das Codegen-Ziel braucht kein Gate-Fragment, es hängt an keinem
  `GATE_CHECKS`.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-19.md` liegt flach (eröffnet),
kein anderer Slice desselben Rolleninhabers (Implementer) liegt in
`in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Wenn sich die
  Docker-only-Protobuf-/buf-Toolchain-Integration während der Umsetzung
  als eigenständiger, von Port/Broadcaster/Server-Grundgerüst unabhängiger
  Liefer-Punkt herausstellt (z. B. weil ein gepinntes Toolchain-Image
  eigene Nacharbeit an mehreren Stellen erzwingt) — dann trennt ein
  Re-Schnitt Toolchain-Anbindung von Adapter-Grundgerüst.
- `in-progress` → `open` (blockiert — Carveout?): Wenn kein geeignetes,
  digest-gepinntes Protobuf-/buf-Toolchain-Image verfügbar ist oder das
  Docker-only-Prinzip (`AGENTS.md` §3.1) mit der Codegen-Toolchain nicht
  einhaltbar ist.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig (§2) + PR gemerged + `make gates` grün + Closure-Notiz mit
Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Die neue Docker-Build-Stufe für Protobuf-/buf-Codegenerierung könnte den
  bestehenden `make image`-Build verlängern oder brechen (neue,
  bislang ungetestete Toolchain-Abhängigkeit). — **Ausgang:** wird bei
  Closure zugewiesen.
- Paralleler Lauf `slice-061` ändert `internal/bootstrap/wiring.go`
  gleichzeitig (anderer Feature-Zweig, additiver `CDC_HTTP_ADDR`-Pfad) —
  Merge-Konflikt-Risiko beim additiven `CDC_GRPC_ADDR`-Zweig in derselben
  Datei. — **Ausgang:** wird bei Closure zugewiesen.
- Das lokal im `grpc`-Paket deklarierte Interface (`changeSubscriber`)
  könnte im ersten Entwurf enger oder weiter gefasst werden, als
  `slice-070`s Anschluss an `CaptureService.WithChangeStream` es braucht,
  und einen kleinen Nacharbeits-Zyklus auslösen. — **Ausgang:** wird bei
  Closure zugewiesen.
- Der `Broadcaster` übergibt ungepuffert und hält die Übergabe an einen
  registrierten, gerade nicht lesenden Empfänger an (`kein Puffer`,
  `ADR-0060` Teilfrage 3); ein langsamer Stream-Client könnte darüber den
  Capture-Pfad stauen, sobald `slice-070` `Publish` an die
  Best-Effort-Kette nach `ACK Source` anschließt. **Eingetreten** (Review
  `slice-069`, bestätigt am Code: `stream.Send` blockiert in derselben
  Goroutine, der der Kanal liest). — **Ausgang: entfallen** —
  [ADR-0066](../../adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
  macht `Publish` nicht-blockierend; die Fixrunde **dieses** Slice
  beseitigt den Auslöser vor der Closure, kein Folge-Slice nötig.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <Guide oder Sensor> <geschärft/ergänzt>: <was genau>
  — liegt in `<AGENTS.md §X | Makefile:<target> | .harness/skills/…>`.
  Auslöser: `BEO-<NNN>` (<slice-NNN>, <slice-MMM>, <slice-KKK> — 3×).
  *(Wurde mit diesem Slice nichts verkörpert — der Normalfall —, entfällt die
  Teil-Zeile `— liegt in …` ersatzlos. Der Eintrag ist dann gezählt, nicht
  verkörpert.)*
- **Beobachtungs-Register (`../observations/`):** <`BEO-<KUERZEL>/<slug>/` neu angelegt, Beleg `evidence/slice-NNN.md` | `evidence/slice-NNN.md` in `BEO-<KUERZEL>/<slug>/` ergaenzt — Zaehler steht damit bei <N>x | keine Beobachtung angefallen>
- **Folge-Slices:** <slice-NNN (<Titel>) — ist eine Datei in `open/`>
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** <nur im Repo ohne Wellen-Betrieb — Anker · Folge-Slice · Register, Ergebnis>

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Die berührten Pfade
(`internal/adapters/driven/grpcstream/`, `internal/adapters/driving/grpc/`,
`internal/application/port/outbound/`, die neue Dockerfile-Build-Stufe)
qualifizieren keine eigene Sub-Area gegenüber dem in
`harness/conventions.md` §Modus-Deklaration deklarierten Default `*`
(Kürzel `PGC`) — keine eigene `MR`-Adaption absehbar, keine eigenständig
abgleichbare Inventur-Linie getrennt vom übrigen Adapter-Bestand. Bleibt
unter dem Default.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`docs/plan/planning/observations/`) durchgegangen, Treffer für `PGC`:

- `BEO-PGC/rollen-verdrahtung` — 4×, **verkörpert** seit `slice-023`.
  Betrifft dieses Slice als Erinnerung: Die neue Token-Prüfung im
  Auth-Interceptor folgt demselben Least-Privilege-Muster (`reader`/`admin`)
  wie die bestehende DSN-/API-Token-Verdrahtung.
- `BEO-PGC/adapter-fehler-ausgang` — 2×, **weiter offen**. Betrifft dieses
  Slice: Der neue `Broadcaster`/gRPC-Server-Fehlerpfad muss einen klaren
  Fehler-Ausgang tragen (kein stiller Fehlerschluck), dieselbe noch offene
  Disziplin wie bei bestehenden Adaptern.
- `BEO-PGC/a-check-null-abdeckung` — 3×, **verkörpert** seit `welle-1`.
  Bestätigt bereits durch beide ADRs: Die neuen Pakete liegen vollständig
  in den bestehenden `.a-check.yml`-Layer-Globs.
- `BEO-PGC/slice-chronik-in-code-kommentar` — 5×, **verkörpert**.
  Implementer-Warnung: **keine** Slice-/Wellen-/Paragraph-Verweise in
  Produktionscode-Kommentaren der neuen Pakete (nur in Test-Godoc als
  etablierte Testfall-Provenienz zulässig, `.harness/skills/reviewer.md`
  HIGH-Punkt „Slice-/Wellen-Chronik in Produktionscode-Kommentar").
- `BEO-PGC/commit-traceability-kein-vorab-hook` — 3×, Ausgang steht noch
  aus (`welle-16`-Closure zuständig). Implementer-Warnung: Commit-Betreffs
  vor `git commit` gegen `SPEC-\d+` prüfen — kein `SPEC-*` im Betreff
  (`tools/harness/commit-traceability.sh`).

Keine weiteren Treffer für die in diesem Slice berührten Pfade.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF — Default-Sub-Area `*` (Kürzel `PGC`), siehe
`harness/conventions.md` §Modus-Deklaration pro Sub-Area. Kein eigener
Sub-Area-Block nötig.
