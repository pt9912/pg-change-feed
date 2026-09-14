# Slice 072: HTTP/SSE-Streaming-Endpunkt

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
[ADR-0061](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md).

**Berührte Spec-Stellen:** [SPEC-018](../../../../spec/pflichtenheft.md)
(Erweiterung — SSE-Nachrichtenschema, Event-Feld-Layout, `event:`-Typ, falls
verwendet; Folgepflicht aus
[ADR-0061](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md)).

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

**Ziel:** Ein zweiter, paralleler Zustellweg für
[LH-FA-SST-008](../../../../spec/lastenheft.md) steht: der Endpunkt
`GET /changes/stream` im bestehenden
`internal/adapters/driving/http/`-Paket, gespeist vom selben
`ChangeStreamPort`/`Broadcaster` aus `slice-069`/`slice-070` als zweiter
Abonnent, geschützt durch dieselbe `withToken`-Middleware, samt
Bootstrap-Entkopplung des `Broadcaster` von `CDC_GRPC_ADDR` und real
belegt gegen den laufenden Feed-Container.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Port/Broadcaster/Server-Grundgerüst und Capture-Integration** —
  bereits geliefert durch `slice-069`/`slice-070` (Bestand); dieser Slice
  fügt nur einen zweiten `Subscribe()`-Aufrufer hinzu
  ([ADR-0061](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md)
  Teilfrage 1).
- **gRPC-Beispiel-Client/E2E-Beleg** — kein Folge-Slice-Bezug: `slice-071`
  ist eine unabhängige Nachweisform für einen anderen Konsumenten
  desselben `Broadcaster`; dieser Slice hängt **nicht** von ihm ab
  ([ADR-0061](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md)
  §Slice-Schnitt-Empfehlung).
- **`Last-Event-ID`-gestütztes Replay und tabellen-granulare Filterung
  des SSE-Streams** — Bestand bleibt bewusst stehen: Beide sind laut
  [ADR-0061](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md)
  §Re-Evaluierungs-Trigger bewusst vertagt (der `Last-Event-ID`-Header
  wird von Server und Handler bewusst ignoriert); kein Gegenstand dieser
  Welle.
- **Änderungen an bestehenden Endpunkten, der Middleware oder der
  Fehler-Antwortform des HTTP-Adapters** — Schicht-Abgrenzung: Dieser
  Slice fügt genau eine neue Route zum bestehenden `mux` hinzu
  ([ADR-0061](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md)
  Teilfrage 2); kein bestehender Handler, keine bestehende Middleware-Zeile
  wird verändert.
- **Ein drittes, eigenes Adapter-Paket oder ein zweites Feature-Gate**
  (`CDC_HTTP_SSE_ENABLED` o. ä.) — von
  [ADR-0061](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md) Teilfrage 2/5
  explizit verworfen; dieser Slice führt keine neue Konfigurationsvariable
  ein.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [LH-FA-SST-008](../../../../spec/lastenheft.md) (Happy-Path/Negative)
      und [ADR-0061](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md)
      Teilfrage 1/2/3/5 umgesetzt, Test referenziert:
      `internal/adapters/driving/http` — neue Route
      `GET /changes/stream` (`http.Flusher`-Flushing je Event,
      `Last-Event-ID` bewusst ignoriert); ein Aufruf ohne oder mit
      unbekanntem Bearer-Token endet mit `401`, vor jedem geschriebenen
      SSE-Event, ein gültiges `reader`-/`admin`-Token öffnet den Stream
      erfolgreich; Bootstrap-Entkopplung der `Broadcaster`-Konstruktion von
      `CDC_GRPC_ADDR` (verdrahtet, sobald `CDC_GRPC_ADDR` **oder**
      `CDC_HTTP_ADDR` gesetzt ist — Unit-Test in `internal/bootstrap`),
      `503`-Pfad bei gesetztem `CDC_HTTP_ADDR` ohne verdrahteten
      `Broadcaster`; Fire-and-Forget-Regressionstest (ein `Publish` ohne
      verbundenen SSE-Client blockiert nicht) (Fitness Function aus
      [ADR-0061](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md)).
- [x] `spec/pflichtenheft.md` erhält die Erweiterung von
      [SPEC-018](../../../../spec/pflichtenheft.md) um das
      SSE-Nachrichtenschema (Event-Feld-Layout, `event:`-Typ falls
      verwendet).
- [x] `tools/harness/sseclient/` (Wegwerf-Beispiel-Client) und ein realer
      SSE-E2E-Rundlauf gegen den laufenden Feed-Container
      (`make test-integration`): eine committed Änderung erreicht einen
      verbundenen SSE-Client mit vollständigem Inhalt; ein
      Verbindungsversuch ohne gültiges Token wird mit `401` abgelehnt.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`docs/reviews/review-slice-072.md`, `.harness/skills/reviewer.md`) —
      Rollenwechsel nach Schritt 8 des Minimal Agent Workflow (`AGENTS.md` §6),
      kein Self-Review (Modul 8).
- [x] Doku-Update `harness/README.md` §Sensors (`make test-integration`
      Zeile: neuer SSE-Rundlauf-Satz).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driving/http/sse.go` (Pfad Implementer-Entscheidung) | neu | Handler für `GET /changes/stream`, `changeSubscriber`-Interface-Nutzung |
| `internal/adapters/driving/http/router.go` (o. ä.) | update | neue Route in den bestehenden `mux`, über `withToken` verdrahtet |
| `internal/adapters/driving/http/*_test.go` | update | Auth-Pfad-Test, Fire-and-Forget-Test |
| `internal/bootstrap/wiring.go` | update | `Broadcaster`-Konstruktion von `CDC_GRPC_ADDR` auf „`CDC_GRPC_ADDR` oder `CDC_HTTP_ADDR`" entkoppelt |
| `internal/bootstrap/*_test.go` | update | Unit-Test für die Oder-Bedingung |
| `spec/pflichtenheft.md` | update | Erweiterung [SPEC-018](../../../../spec/pflichtenheft.md) um SSE-Nachrichtenschema |
| `tools/harness/sseclient/` | neu | Wegwerf-Beispiel-Client, analog `tools/harness/grpcclient/`/`tools/harness/httpclient/` |
| `tools/harness/run-integration-tests.sh` | update | Rundlauf-Erweiterung: SSE-Client verbindet, empfängt Change, negativer Token-Test |
| `harness/README.md` | update | `make test-integration`-Sensor-Zeile um SSE-Rundlauf-Satz ergänzt |
| `internal/adapters/driving/http/sse_test.go` | neu | Whitebox-Tests: `401` (fehlend/unbekannt), Happy Path mit vollständigem Event, `null`-Row-Image, `503` ohne Broadcaster, Freigabe der Subskription bei Verbindungsende |
| `internal/bootstrap/changestream_internal_test.go` | neu | Unit-Test der Oder-Bedingung `changeStreamEnabled` (`CDC_GRPC_ADDR` oder `CDC_HTTP_ADDR`) |
| `docs/user/benutzerhandbuch.md` | unverändert (Aufschub) | neue Betreiber-Oberfläche (Endpunkt `GET /changes/stream`, kein neues ENV-Feld); die Handbuch-Dokumentation holt `slice-077` nach — benannte Aufschub-Adresse statt Mitnahme in diesem Diff |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-069` **und** `slice-070` liegen
in `done/` — `ChangeStreamPort` und der `Broadcaster` müssen existieren
**und** in `CaptureService` verdrahtet sein, sonst sieht der SSE-Endpunkt
nie ein Event ([ADR-0061](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md)
§Slice-Schnitt-Empfehlung). **Nicht** erforderlich: `slice-071` (gRPC-
Beispiel-Client/E2E) — unabhängige Nachweisform, kein Trigger für diesen
Slice. Zusätzlich: `slice-061` (HTTP-API-Beispiel-Client/E2E) liegt in
`done/` — derselbe Adapter-Pfad
(`internal/adapters/driving/http/`, `harness/README.md`,
`internal/bootstrap/wiring.go`) wird von `slice-061` parallel bearbeitet;
dieser Slice beginnt erst, wenn `slice-061` gemerged ist (Merge-Konflikt-
Vermeidung, kein inhaltlicher Trigger). WIP-Limit frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Wenn die
  Bootstrap-Entkopplung (`CDC_GRPC_ADDR` oder `CDC_HTTP_ADDR`) sich als
  eigenständiger Liefer-Punkt mit substanzieller Nacharbeit an bestehenden
  Bootstrap-Tests herausstellt.
- `in-progress` → `open` (blockiert — Carveout?): Wenn `withToken`s
  Verhalten bei lang laufenden Verbindungen entgegen der in
  [ADR-0061](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md) Teilfrage 2
  dokumentierten Prüfung doch stört (z. B.
  Response-Buffering durch eine dazwischenliegende Schicht) und eine
  Design-Korrektur an der Middleware selbst nötig wird.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig (§2) + PR gemerged + `make gates` grün + realer
`make test-integration`-SSE-Rundlauf grün + Closure-Notiz mit
Steering-Loop-Lerneintrag geschrieben (§7).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Direkte Datei-Überschneidung mit `slice-061`
  (`internal/adapters/driving/http/`, `compose.yaml`,
  `tools/harness/run-integration-tests.sh`, `harness/README.md`,
  `internal/bootstrap/wiring.go`, `tools/harness/httpclient/`) — auch bei
  eingehaltenem Start-Trigger (§4) bleibt das Risiko einer inhaltlichen
  Überschneidung (z. B. beide Slices ändern denselben Router-Abschnitt). —
  **Ausgang: entfallen** — `slice-061` war längst geschlossen; die
  Erweiterungen trafen keine parallele Änderung.
- Die Bootstrap-Entkopplung ändert eine bestehende Bedingung
  (`CDC_GRPC_ADDR` → „`CDC_GRPC_ADDR` oder `CDC_HTTP_ADDR`") in
  `internal/bootstrap/wiring.go` — Regressionsrisiko für `slice-069`s
  bereits bestehenden gRPC-only-Verdrahtungspfad, falls die Oder-Bedingung
  fehlerhaft umgesetzt wird. — **Ausgang: entfallen** — der Verifier hat die
  Oder-Bedingung am Code bestätigt (derselbe `Broadcaster` an beide Adapter;
  beide Adressen leer → unverändertes Bestandsverhalten) und die Mutation
  `||`→`&&` färbt genau die zwei Sub-Fälle rot.
- `withToken`s Eignung für lang laufende Verbindungen ist in
  [ADR-0061](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md) Teilfrage 2
  geprüft, aber nicht real getestet — ein bislang unentdecktes
  Response-Buffering (z. B. durch einen künftig eingeführten
  Reverse-Proxy) würde erst hier sichtbar. — **Ausgang: entfallen** — der
  reale `make test-integration`-Lauf zeigt den Client die Änderung
  **während** der offenen Verbindung empfangen (`READY` → `RECEIVED`), also
  kein Buffering des Adapters.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** Der zweite Zustellweg brauchte keinen neuen
  Server und keinen neuen Port — derselbe HTTP-Adapter, dieselbe
  Token-Middleware, derselbe `Broadcaster` als zweiter Abonnent. Die
  Bootstrap-Entkopplung ist eine reine Oder-Bedingung
  (`changeStreamEnabled`), und der reale Lauf zeigt beide Wege gegen
  **denselben** laufenden Prozess grün (`RECEIVED change_id=967-1` über SSE,
  `REJECTED code=401` ohne Token). Der Verifier hat drei Mutationen selbst
  rot gesehen.
- **Was ging anders als geplant:** (a) Der SSE-Block lag zunächst als
  Erweiterung in `SPEC-018`; `ADR-0061`s Folgepflicht verlangt aber einen
  **neuen** `SPEC-*`-Eintrag, und `SPEC-018`s Intro grenzt sich ausdrücklich
  auf die neun Port-gedeckten Fähigkeiten aus `LH-FA-SST-006` ab — der
  Block ist deshalb bei der Closure nach `SPEC-021` herausgelöst (Review
  F-6, Verifier bestätigt). (b) Der `503`-Pfad ist über den regulären
  Start unerreichbar (Review F-2) — benannt, kein Fix. (c) Eine Verzweigung
  in `rowImage` ist durch keinen Test gebunden (Review F-4, mit Probe
  bestätigt) — benannte Grenze.
- **Steering-Loop-Eintrag:** keiner neu verkörpert — die zwei Klassen aus
  den Reviews sind benannt statt gezählt: „Aufschub-Adresse nimmt die
  Sendung nicht an" (`BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an`, mit
  diesem Slice **2×**) und „Adapter-Unit-Test verdeckt eine
  Bootstrap-Lücke" (`BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke`,
  mit diesem Slice **2×**).
  — liegt in `<AGENTS.md §X | Makefile:<target> | .harness/skills/…>`.
  Auslöser: `BEO-<NNN>` (<slice-NNN>, <slice-MMM>, <slice-KKK> — 3×).
  *(Wurde mit diesem Slice nichts verkörpert — der Normalfall —, entfällt die
  Teil-Zeile `— liegt in …` ersatzlos. Der Eintrag ist dann gezählt, nicht
  verkörpert.)*
- **Beobachtungs-Register (`../observations/`):** neues Verzeichnis
  `BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an/` angelegt, Beleg
  `evidence/slice-072.md` (Zähler 2×); `evidence/slice-072.md` in
  `BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke/` ergänzt (Zähler 2×).
- **Folge-Slices:** keiner aus diesem Slice — der Handbuch-Aufschub zeigt
  auf `slice-077` (in `open/`), dessen DoD-Punkt bei dieser Closure um den
  SSE-Endpunkt erweitert wurde.
- **Risiken aus §6:** alle drei entfallen — siehe §6.
- **Drei Paarungen:** Repo **mit** Wellen-Betrieb (`welle-19` offen) —
  Prüfung läuft bei der `welle-19`-Closure.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Der berührte Pfad
(`internal/adapters/driving/http/`, `internal/bootstrap/`, `tools/harness/`,
`spec/pflichtenheft.md`) qualifiziert keine eigene Sub-Area gegenüber dem
in `harness/conventions.md` §Modus-Deklaration deklarierten Default `*`
(Kürzel `PGC`) — bleibt unter dem Default.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`docs/plan/planning/observations/`) durchgegangen, Treffer für `PGC`:

- `BEO-PGC/rollen-verdrahtung` — 4×, **verkörpert** seit `slice-023`.
  Betrifft dieses Slice als Erinnerung: Die Wiederverwendung der
  bestehenden `CDC_API_TOKEN_READER`/`ADMIN`-Klassen folgt demselben
  Least-Privilege-Muster.
- `BEO-PGC/adapter-fehler-ausgang` — 2×, **weiter offen**. Betrifft dieses
  Slice: Der `503`-Pfad bei fehlender `Broadcaster`-Verdrahtung ist die
  konkrete Umsetzung eines klaren Fehler-Ausgangs statt eines stillen
  Blockierens.
- `BEO-PGC/slice-chronik-in-code-kommentar` — 5×, **verkörpert**.
  Implementer-Warnung: keine Slice-/Wellen-Verweise in
  Produktionscode-Kommentaren des neuen Handlers.
- `BEO-PGC/commit-traceability-kein-vorab-hook` — 3×, Ausgang steht noch
  aus (`welle-16`-Closure zuständig). Implementer-Warnung:
  Commit-Betreffs vor `git commit` gegen `SPEC-\d+` prüfen.

Keine weiteren Treffer für die in diesem Slice berührten Pfade.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF — Default-Sub-Area `*` (Kürzel `PGC`), siehe
`harness/conventions.md` §Modus-Deklaration pro Sub-Area. Kein eigener
Sub-Area-Block nötig.
