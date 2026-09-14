# Slice 070: gRPC-Capture-Integration

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-19.

**Bezug:** [LH-FA-SST-008](../../../../spec/lastenheft.md) (Haupt-Bezug),
[ADR-0060](../../adr/0060-grpc-streaming-mechanismus.md),
[ADR-0066](../../adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
(Zustellsemantik des `Broadcaster` — `Publish` blockiert nie auf einen
Abonnenten).

**Berührte Spec-Stellen:** [SPEC-020](../../../../spec/pflichtenheft.md)
(Referenz — Publish-Aufrufkontrakt und Zustellsemantik bereits durch
[ADR-0060](../../adr/0060-grpc-streaming-mechanismus.md),
[ADR-0066](../../adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
und `slice-069` fixiert, keine inhaltliche Änderung erwartet).

> **Plan-Korrektur 2026-09-14 (Architect-Verdikt, `ADR-0066`):** Die
> Kennung der berührten Spec-Stelle ist `SPEC-020` (Stream-Semantik), nicht
> `SPEC-019` (Feldform von `cdc.administration_request`). Der Stream-Publish
> ruft `Publish` **synchron** in der Best-Effort-Kette und isoliert nur den
> **Fehler** — die **Zeit** braucht er nicht zu isolieren: `ADR-0066` macht
> `Publish` nicht-blockierend (begrenzte Empfangs-Warteschlange je Abonnent,
> Drop-Newest), der Erzeuger hält nie auf einen Abonnenten an. Eine
> caller-seitige Goroutine oder Deadline ist ausdrücklich **nicht** Teil
> dieses Slice.

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

**Ziel:** `CaptureService` erhält die neue, optionale Konstruktions-Option
`WithChangeStream` und die tatsächliche Publish-Kette
(`Receive → Decode → Persist → COMMIT Store → ACK Source → Notify (best
effort) → Stream-Publish (best effort)`) — der in `slice-069` gebaute
`ChangeStreamPort`/`Broadcaster` bekommt damit erstmals einen realen
Aufrufer, isoliert und fehlerisoliert getestet. Der Publish-Aufruf ist
nicht-blockierend ([ADR-0066](../../adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)):
die Fehlerisolation an der Aufrufstelle genügt, eine Zeit-Isolation (eigene
Goroutine, Deadline) ist weder nötig noch Gegenstand dieses Slice.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Port/Broadcaster/Server-Grundgerüst selbst** — bereits geliefert durch
  `slice-069` (Bestand); dieser Slice verdrahtet nur den bestehenden
  Vertrag, ändert seine Signatur nicht.
- **gRPC-Beispiel-Client und E2E-Rundlauf** — Folge-Slice `slice-071`
  übernimmt es; dieser Slice liefert die Publish-Kette, nicht ihren
  externen Beleg.
- **HTTP/SSE-Endpunkt** — anderer Vorgang: `slice-072`
  ([ADR-0061](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md)) nutzt den
  hier erstmals befüllten `Broadcaster` als zweiten Abonnenten, ist aber
  ein eigener Liefer-Punkt in einem anderen Adapter-Paket.
- **Änderungen an `ChangeNotificationPort`/`natsnotify` oder am
  bestehenden `Receive → Decode → Persist → COMMIT Store → ACK
  Source`-Kernpfad** — Schicht-Abgrenzung: Der Stream-Publish-Schritt
  reiht sich **nach** dem bestehenden Notify-Schritt ein und darf ihn
  nicht verändern ([ADR-0060](../../adr/0060-grpc-streaming-mechanismus.md)
  Teilfrage 2, dieselbe Grenze wie
  [ADR-0011](../../adr/0011-persist-before-ack.md)/
  [ADR-0027](../../adr/0027-capture-application-service.md)); dieser Slice
  fasst keine Zeile des Kernpfads vor dem Stream-Publish-Aufruf an.
- **Stream-internes Replay und tabellen-granulare Filterung** — Bestand
  bleibt bewusst stehen, dieselbe Vertagung wie in `slice-069` §1
  begründet; kein Gegenstand dieser Welle.

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] [LH-FA-SST-008](../../../../spec/lastenheft.md) Happy-Path erfüllt,
      Test referenziert:
      `internal/application/usecase/capture` — `CaptureService` ruft nach
      `ACK Source` und nach dem bestehenden Notify-Schritt genau einmal
      `ChangeStreamPort.Publish` je `Change` der committed Transaktion auf
      (`WithChangeStream`-Option gesetzt).
- [x] [ADR-0060](../../adr/0060-grpc-streaming-mechanismus.md) Teilfrage 2
      vollständig umgesetzt: Fehlerisolations-Regressionstest — ein
      fehlschlagender `ChangeStreamPort` darf `Capture()`s Rückgabewert
      nicht beeinflussen, wenn `store`/`ack` erfolgreich waren, analog zum
      bestehenden `ChangeNotificationPort`-Test
      ([ADR-0055](../../adr/0055-nats-change-notification-wecksignal.md));
      Fire-and-Forget-Regressionstest bei getrenntem Client (kein aktiver
      Subscriber blockiert `Capture()` nicht) **und** bei einem
      **registrierten, nicht lesenden** Empfänger — der Aufruf kehrt ohne
      Zeit-Isolation zurück ([ADR-0066](../../adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)).
      *(Die dritte Fitness-Function-Zeile des
      [ADR-0066](../../adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
      war gegen dessen eigene Entscheidung formuliert und wurde durch
      [ADR-0067](../../adr/0067-capture-publish-einbindung-fitness-function-korrektur.md)
      korrigiert — die zwei gebauten Tests sind die **zutreffende** Zeile,
      keine Reduktion; Begründung: §3.)*
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8);
      Report: [`review-slice-070`](../../../reviews/review-slice-070.md),
      Verdikt 2 HIGH / 0 MEDIUM / 2 LOW / 4 INFO, F-1 durch
      [ADR-0067](../../adr/0067-capture-publish-einbindung-fitness-function-korrektur.md)
      entschieden, F-2/F-3/F-4/F-6 in der Fixrunde behoben, damit ohne weitere
      Fixrunde geschlossen (Nachzug nach `.harness/skills/reviewer.md`
      §DoD-Checkbox-Nachzug ohne Fixrunde).
- [x] Doku-Update `internal/bootstrap/wiring.go`-Kommentar zur
      `CDC_GRPC_ADDR`-Verdrahtung, falls sich der Aktivierungspfad seit
      `slice-069` sichtbar ändert. *(Eingetreten: der `Broadcaster` wird an
      der Capture-Verdrahtung konstruiert und trägt den
      `ChangeStreamPort`; beide `CDC_GRPC_ADDR`-Kommentarblöcke
      nachgezogen.)*
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
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
| `internal/application/usecase/capture/service.go` | update | neue Option `WithChangeStream(stream outbound.ChangeStreamPort) Option`; Publish-Aufruf nach `ACK Source`/Notify, synchron (nicht-blockierend, [ADR-0066](../../adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)) |
| `internal/application/usecase/capture/service_test.go` | update | Fehlerisolations-Regressionstest, Fire-and-Forget-Regressionstest |
| `internal/bootstrap/wiring.go` | update | `CaptureService` erhält den in `slice-069` konstruierten `*grpcstream.Broadcaster` über `WithChangeStream`, wenn `CDC_GRPC_ADDR` gesetzt ist |

**Implementer-Entscheidungen und -Abweichungen (Plan-Nachzug im selben Lauf):**

- **Plan-Nachzug (nach `ADR-0067`) — Zeit-Isolations-Test:** Die
  Fitness-Function-Zeile aus
  [`ADR-0066`](../../adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
  („ein `ChangeStreamPort`, dessen `Publish` nicht zurückkehrt … darf
  `Capture()` nicht anhalten") widersprach der Entscheidung derselben ADR
  (synchroner Aufruf, Option B ausdrücklich verworfen) und ist deshalb
  **nicht als Test gebaut** — der Reviewer hat das als HIGH bestätigt und an
  den Architect gereicht;
  [`ADR-0067`](../../adr/0067-capture-publish-einbindung-fitness-function-korrektur.md)
  superseded genau diese Zeile. Die zutreffende Zusage ist geschichtet: das
  Nicht-Blockieren ist **Vertragspflicht des Ports** und wird **am Port**
  belegt (Paket `internal/adapters/driven/grpcstream`,
  `slice-069`-Fixrunde); der Anteil der Capture-Schicht ist die folgende
  Eigenschaft. Die Anschlussstelle testet die zwei Hälften, die ihr zufallen:
  (a) die Capture-kritische Kette `Persist → ACK` steht vollständig und
  unberührt, wenn der Port bei `Publish` eintritt
  (`TestCapturePublishOhneRueckkehrHaeltKritischeKetteNichtAn`),
  (b) `Capture()` kehrt ohne eigene Zeit-Isolation zurück — bei getrenntem
  Client und bei einem registrierten, nicht lesenden Empfänger
  (`TestCaptureKehrtOhneUndMitNichtLesendemStreamEmpfaengerZurueck`).
- **Keine caller-seitige Goroutine, keine Deadline** — bestätigt §1 und
  [`ADR-0066`](../../adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
  §Entscheidung; `internal/application/usecase/capture` importiert
  `internal/adapters/driven/grpcstream` nicht (`.a-check.yml` kennt keine
  `app → adapters`-Kante), der Capture-Test trägt deshalb einen
  vertragstreuen Port-Doppel.
- **Keine Zeile des Kernpfads vor dem Stream-Publish-Aufruf angefasst**
  (§1-Schicht-Abgrenzung): der bestehende `Receive → … → ACK Source`-Pfad
  und der Notify-Block bleiben unverändert; der neue Block steht danach.
- **Fixrunde nach Review** (`../../../reviews/review-slice-070.md`): vier
  Findings am Code behoben — F-2 (der `wiring.go`-Kommentar behauptete, ein
  gRPC-Server-Startfehler erreiche den Rückgabewert von `Run`; tatsächlich
  läuft der Server in eigener Goroutine und der Fehler wird dort über
  `log.Error` gemeldet, wie beim HTTP-Adapter), F-3 (NATS-Kommentar direkt
  an sein `if cfg.NatsURL != ""` gerückt), F-4
  (`erwarteteZustellung` → `erwarteteErfolgreicheAufrufe`: die Größe zählt
  zurückgekehrte `Publish`-Aufrufe, nicht zugestellte Changes), F-6 (die
  Begründung des verworfenen `tx.Changes()`-Fehlers steht jetzt einmal im
  Helfer `changesOfCommittedTransaction`, den `Capture()` und
  `distinctTables` gemeinsam nutzen). F-1 ist durch
  [ADR-0067](../../adr/0067-capture-publish-einbindung-fitness-function-korrektur.md)
  entschieden, F-5/F-7/F-8 sind INFO ohne Code-Änderung.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-069` liegt in `done/` — der neue
Outbound Port `ChangeStreamPort` und der `Broadcaster` müssen existieren,
bevor `CaptureService` sie verdrahten kann; WIP-Limit frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Wenn sich die
  Fehlerisolations- und Fire-and-Forget-Tests als eigenständiger,
  aufwändigerer Liefer-Punkt herausstellen als hier geschätzt (etwa weil
  `CaptureService`s bestehende Test-Fixtures für den neuen Port
  substanziell umgebaut werden müssen).
- `in-progress` → `open` (blockiert — Carveout?): Wenn `slice-069`s
  lokal deklariertes Interface (`changeSubscriber`) beim Anschluss nicht
  zu `CaptureService`s Aufrufmuster passt und eine Design-Korrektur an
  `slice-069` selbst nötig wird (Rücksprache mit dem Architect-Verdikt).

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

- Ein zusätzlicher Schritt im `Capture()`-Aufrufpfad (Stream-Publish nach
  Notify) könnte versehentlich vor statt nach `ACK Source` einsortiert
  werden und damit den unantastbaren Kernpfad
  ([ADR-0011](../../adr/0011-persist-before-ack.md)) berühren. —
  **Ausgang:** wird bei Closure zugewiesen.
- Der Fehlerisolations-Test könnte einen bereits bestehenden,
  strukturell ähnlichen `ChangeNotificationPort`-Test unbeabsichtigt
  duplizieren statt ihn zu spiegeln, was Testpflege-Aufwand doppelt. —
  **Ausgang:** wird bei Closure zugewiesen.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Der berührte Pfad
(`internal/application/usecase/capture/`, `internal/bootstrap/`)
qualifiziert keine eigene Sub-Area gegenüber dem in
`harness/conventions.md` §Modus-Deklaration deklarierten Default `*`
(Kürzel `PGC`) — bleibt unter dem Default.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`docs/plan/planning/observations/`) durchgegangen, Treffer für `PGC`:

- `BEO-PGC/rollen-verdrahtung` — 4×, **verkörpert** seit `slice-023`.
  Betrifft dieses Slice nicht direkt (keine DSN-/Rollen-Verdrahtung),
  genannt zur Vollständigkeit der Sichtung.
- `BEO-PGC/adapter-fehler-ausgang` — 2×, **weiter offen**. Betrifft dieses
  Slice unmittelbar: Der Fehlerisolations-Test ist die konkrete
  Umsetzung der in dieser Beobachtung verlangten Disziplin für den neuen
  `ChangeStreamPort`-Aufruf.
- `BEO-PGC/slice-chronik-in-code-kommentar` — 5×, **verkörpert**.
  Implementer-Warnung: keine Slice-/Wellen-Verweise in
  `service.go`-Produktionscode-Kommentaren.
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
