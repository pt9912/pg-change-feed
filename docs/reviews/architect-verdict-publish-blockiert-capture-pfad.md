# Architect-Verdikt: `Publish` hält den Capture-Pfad an — Entkopplung im `Broadcaster`

**Rolle:** Architect (Modul 8)

**Anlass:** Finding F-1 (MEDIUM, `SPEC-020`/Port-Godoc pinnen nur die
Hälfte der Zustellsemantik), F-5 (INFO, Godoc-Satz ohne Deckung) und das im
Verdikt-Abschnitt ausdrücklich als **echt** bestätigte Broadcaster-Risiko aus
[`review-slice-069.md`](review-slice-069.md). Der Reviewer hat das Risiko
nicht als Finding gegen den Code geführt, sondern als **Entscheidung**
gereicht, die in `slice-070` <!-- d-check:status-provenance --> fällig und dort
sichtbar zu machen ist — damit läuft der Rollenwechsel Reviewer → Architect →
Implementer (Modul 8 §Konflikt-Pfad, Verdikt 2: Folge-ADR statt stiller
Korrektur).

**Rolleninhaber:** pt9912 (Architect-Zug, anderer Kontext als der
Implementer-Lauf von `slice-069` <!-- d-check:status-provenance --> und als
der Reviewer-Lauf)

**Datum:** 2026-09-14

**Bezug:** [`LH-FA-SST-008`](../../spec/lastenheft.md) (Happy Path, Boundary,
Negative — insbesondere das Boundary-Kriterium), [`LH-FA-REA-001`](../../spec/lastenheft.md),
[`LH-FA-CON-003`](../../spec/lastenheft.md)/[`LH-FA-CON-005`](../../spec/lastenheft.md),
[`LH-QA-REL-001`](../../spec/lastenheft.md) (keine stillen **Daten**verluste),
[`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md) (Teilfrage 2
Anknüpfung/Design, Teilfrage 3 Zustellsemantik — in der Puffer-Klausel durch
die Folge-ADR korrigiert), [`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md)
(zweiter Zustellweg über denselben `Broadcaster`),
[`ADR-0055`](../plan/adr/0055-nats-change-notification-wecksignal.md) Punkt 1
(Nachvollziehbarkeit beim Store, ein verlorenes Signal ist folgenlos) und
Punkt 4 (Capture-kritischer Pfad unantastbar),
[`ADR-0011`](../plan/adr/0011-persist-before-ack.md),
[`ADR-0027`](../plan/adr/0027-capture-application-service.md),
[`SPEC-020`](../../spec/pflichtenheft.md) (Zustellsemantik-Zeile),
`internal/adapters/driven/grpcstream/broadcaster.go`,
`internal/adapters/driving/grpc/server.go` (`StreamChanges`),
`internal/application/port/outbound/changestream.go`,
`internal/application/usecase/capture/service.go`,
`docs/plan/planning/in-progress/slice-069-grpc-streaming-adapter-grundgeruest.md` <!-- d-check:status-provenance -->
(noch offen, Fixrunde), `docs/plan/planning/open/slice-070-grpc-capture-integration.md` <!-- d-check:status-provenance -->
(der Aufrufer), `docs/plan/planning/welle-19.md`

**Erzeugte Artefakte dieses Zugs:**

- [`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
  (Accepted, `Supersedes ADR-0060` **nur** die Puffer-Klausel aus Teilfrage 3
  Option A und die sie wiederholende Design-Zeile in Teilfrage 2), samt
  Index-Zeile und Supersedes-Vermerk in
  [`docs/plan/adr/README.md`](../plan/adr/README.md)
- Präzisierung
  [`SPEC-020`](../../spec/pflichtenheft.md) (Zeile *Zustellgarantie* + neue
  Zeile *Erzeuger-Blockade*)
- Plan-Präzisierung `slice-069` <!-- d-check:status-provenance --> (§3
  Implementer-Entscheidung zur `Publish`-Semantik korrigiert; §6-Risiko trägt
  seinen Ausgang) und `slice-070` <!-- d-check:status-provenance --> (§1/§3 auf
  den nicht-blockierenden `Publish` nachgezogen, Kennungs-Fehler `SPEC-019`
  → `SPEC-020` berichtigt)
- diese Verdikt-Datei

**Nicht erzeugt:** ein **neuer Folge-Slice**. Die Umsetzung ist ein Fix in
`slice-069` <!-- d-check:status-provenance --> — der `Broadcaster` ist dort
gerade entstanden, der Slice ist noch **nicht** in `done/` und trägt bereits
eine offene Rückkante (F-1, F-2 → Fixrunde). Ein eigener Folge-Slice würde
einen bereits nicht geschlossenen Liefer-Punkt doppeln.

---

## Verdikt

**Das Broadcaster-Risiko ist echt, und die Entkopplung gehört in den
`Broadcaster`.** Der ungepufferte, blockierende Send in `Publish` legt die
Blockade dorthin, wo der Kanal liegt; kein Aufrufer kann sie auflösen, ohne
selbst zu blockieren oder unbeschränkt zu wachsen. `Publish` übergibt künftig
**nicht-blockierend** an eine **begrenzte Empfangs-Warteschlange je Abonnent**
und **verwirft bei Überlauf** den eintreffenden Change für diesen Abonnenten.
Der Erzeuger hält nie an; der Aufrufer (`CaptureService`, `slice-070` <!-- d-check:status-provenance -->)
bleibt unverändert bei synchronem `Publish` und isoliert nur den Fehler.

**Die Zustellsemantik wird vollständig gepinnt.** Ab jetzt gilt: *Ein nicht
verbundener **oder langsamer lesender** Empfänger verliert die betroffenen
Nachrichten ersatzlos; der Erzeuger hält nie an.* Ein verlorenes Stream-Signal
ist **kein Datenverlust** — die Nachvollziehbarkeit liegt beim
`ChangeStorePort` ([`ADR-0055`](../plan/adr/0055-nats-change-notification-wecksignal.md)
Punkt 1), verpasste Changes bleiben über den bestehenden Lesezugriffsweg und
die bestätigte Consumer-Position nachholbar
([`LH-FA-SST-008`](../../spec/lastenheft.md) Boundary). Damit ist F-1
aufgelöst: die Semantik-Zeile und der Port-Godoc tragen beide Hälften — die
verwerfende **und** die entkoppelnde.

**Die Lösung ist nicht die verworfene Sache unter neuem Namen.** Der
Unterschied ist der, den `ADR-0060`s Verwerfung trägt: **Verwerfen statt
Nachliefern**. Die Warteschlange liefert nie nach, überlebt keine Subskription
und trägt keinen Positionsbegriff; Replay- (Option B) und Ring-Puffer
(Option C) existierten beide, um einem **getrennten/wiederkehrenden** Consumer
etwas **nachzuliefern**.

**Träger:** Folge-ADR — [`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md),
`Supersedes ADR-0060` **nur** für die betroffene Klausel (Muster
`ADR-0048`/`ADR-0062`/`ADR-0063`/`ADR-0065`). Eine bloße Präzisierung in
`slice-070`s Plan plus `SPEC-020` genügt **nicht**: `ADR-0060` Teilfrage 3
Option A und die Design-Zeile aus Teilfrage 2 sind `Accepted` und sagen „kein
Puffer" ausdrücklich zu; eine Accepted-ADR wird nicht nachgebessert
(`AGENTS.md` §3.5). **Umsetzung:** Fix in `slice-069` <!-- d-check:status-provenance -->
(Fixrunde), kein eigener Folge-Slice.

---

## Befundlage (eigener Code-Gang)

Drei Eigenschaften des ausgelieferten Standes, am Code nachvollzogen — sie
sind die Grundlage des Verdikts, nicht die Wiederholung des Reviewer-Texts:

1. **Der Kanal ist ungepuffert, der Send blockierend.**
   `internal/adapters/driven/grpcstream/broadcaster.go`,
   `Subscribe` legt `make(chan *model.Change)` ohne Kapazität an;
   `Publish` bedient die Abonnenten **nacheinander** mit
   `select { case sub.changes <- change: case <-sub.done: case <-ctx.Done(): }`.
   Ohne bereiten Leser wartet dieser Send; die einzigen Ausgänge sind
   `sub.done` (Abmeldung) oder `ctx.Done()`.
2. **`ctx` bindet nicht, es lebt so lange wie der Prozess.** Der einzige
   `Publish`-Aufrufer ist die Capture-Schleife; ihr `ctx` endet mit dem
   Prozess. Der Halt ist damit real unbeschränkt — die „Schranke", die der
   Reviewer im `ctx` benennt, ist keine.
3. **Der gRPC-Handler liest und sendet in einer Goroutine.**
   `internal/adapters/driving/grpc/server.go`, `StreamChanges` liest
   `changes` und ruft `stream.Send` **in derselben** Goroutine. Sobald das
   Flow-Control-Fenster des Clients voll ist — der Client also nicht liest —
   blockiert `stream.Send`, und die Goroutine ist **nicht** auf Empfang
   geparkt. Genau in diesem Zustand trifft der ungepufferte Send keinen Leser.

Die Kette ist damit: nicht lesender Client → `stream.Send` blockiert → Handler
liest `changes` nicht → `Publish` blockiert am ungepufferten Kanal →
`Receive` jeder weiteren Transaktion wartet → `cdc_capture_lag` wächst bis zum
Stillstand. Zusätzlich Head-of-Line: `Publish` bedient die Abonnenten
nacheinander, ein gestörter Erster hält alle Übrigen.

**Bestätigt, nicht angenommen:** Der Fix braucht **keine** Änderung an den
Driving-Adaptern. `changeSubscriber` (`server.go:29-31`) und
`Subscribe() (<-chan *model.Change, func())` bleiben in Signatur und
Semantik identisch; eine begrenzte Kanalkapazität ist für den Abonnenten
unsichtbar. Damit erbt auch der SSE-Weg aus
[`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md) dieselbe
Entkopplung ohne eigenen Code, und `.a-check.yml` bleibt unberührt.

---

## Frage 1 — Wirkort: `Broadcaster`, Aufrufer oder beides?

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (ungepuffert/blockierend, synchroner `Publish` in `slice-070` <!-- d-check:status-provenance -->) | kein Code-Eingriff, `ADR-0060` wörtlich erfüllt | der nicht lesende Client hält über den ungepufferten Kanal den `Receive` jeder weiteren Transaktion an — genau die Verzögerung des Capture-kritischen Pfads, die [`ADR-0011`](../plan/adr/0011-persist-before-ack.md)/[`ADR-0027`](../plan/adr/0027-capture-application-service.md) und [`ADR-0055`](../plan/adr/0055-nats-change-notification-wecksignal.md) Punkt 4 ausschließen; dazu Head-of-Line über alle Empfänger |
| B — nur der Aufrufer: `slice-070` <!-- d-check:status-provenance --> ruft nicht-blockierend (Goroutine je Change oder Deadline-`ctx`) | kein Eingriff am `Broadcaster`, dessen Textkörper bleibt unberührt | „Goroutine je Change" wächst mit dem Rückstand (eine je Transaktion, die auf den gestörten Empfänger wartet); eine Deadline-`ctx` verlagert den Halt nur und legt **jedem** Publish eine Wartezeit auf — und trifft die falsche Aussage, weil `Publish` die Abonnenten **nacheinander** bedient: ein blockierender Erster hält die Übrigen bis zum Ablauf. Der Aufrufer kann die serielle Schleife im `Broadcaster` nicht überspringen; die Invariante läge in jedem Aufrufer statt an einer Stelle |
| **C — im `Broadcaster`: begrenzte Empfangs-Warteschlange je Abonnent + nicht-blockierender Send mit Drop-Newest (gewählt)** | eine Stelle für alle Abonnenten und beide Zustellwege; beide Driving-Adapter unverändert; der Erzeuger kann strukturell nicht warten; kein Replay-Zustand | superseded eine Klausel einer `Accepted`-ADR; begrenzter Speicher je Abonnent; ein langsamer Consumer verliert Stream-Nachrichten still (per Festlegung) |
| D — asynchroner Versand je Abonnent mit eigener Goroutine und begrenzter Warteschlange | trennt Sender-/Empfänger-Fortschritt sichtbar | liefert dieselbe Garantie nur bei nicht-blockierendem Erzeuger→Warteschlange-Handoff, ist dann C plus eine zusätzliche Goroutine mit eigener Lebensdauer und `cancel`-Aufräumen; keine Zusatz-Zusage |
| E — Entkopplung im Driving-Adapter (Broadcaster bleibt ungepuffert, der Handler schiebt nicht-blockierend in eine eigene Warteschlange) | `ADR-0060`s „kein Puffer" bliebe wörtlich erfüllt | hängt an einer **Scheduling**-Eigenschaft (der Drainer ist *fast immer*, aber nicht *immer* auf Empfang geparkt) statt an einer strukturellen; das Warteschlangen-Muster würde in jedem nachgelagerten Adapter dupliziert (gRPC jetzt, SSE in `slice-072` <!-- d-check:status-provenance -->) — dieselbe Duplikations-Sorge wie F-3 |

**Fazit Frage 1:** **C**. Der Kanalbesitzer ist der `Broadcaster`; wer die
Blockade hält, muss sie auflösen. Der Aufrufer wird **nicht** angefasst — sobald
`Publish` nicht blockiert, ist die Best-Effort-Kette nach `ACK Source` in Zeit
und Fehler isoliert.

## Frage 2 — Warum das nicht die verworfene Option B/C aus `ADR-0060` ist

Der Unterschied ist nicht der Name, sondern die Richtung: **Verwerfen statt
Nachliefern.**

- **Option B (Replay ab bestätigter Consumer-Position)** und **Option C
  (In-Memory-Ring-Puffer begrenzter Tiefe)** stehen in `ADR-0060` Teilfrage 3,
  um einem **getrennten** Consumer etwas **nachzuliefern** — B mit Zusage
  (Position, Konsistenzgrenze, Live-Übergang), C als kurze
  Verbindungslücken-Überbrückung. Beide brauchen Zustand, der einen **Connect
  überlebt**.
- Die gewählte Warteschlange liefert **nie** nach: sie wird bei `cancel` mit
  dem Abonnenten verworfen, steht **keiner** neuen Subskription zur Verfügung
  und trägt keinen Positionsbegriff. Sie wirkt ausschließlich **innerhalb**
  einer bestehenden, **verbundenen** Subskription, um deren momentane
  Nicht-Lesebereitschaft zu absorbieren — der Grund ist ein blockierender
  `stream.Send`, nicht eine getrennte Leitung. Ein Client, der sich trennt und
  neu verbindet, erhält aus ihr **kein** Byte.

Die `ADR-0060`-Zusage „ohne Empfänger wird die Nachricht verworfen" bleibt in
Kraft; präzisiert wird nur, dass dieses Verwerfen **nicht-blockierend** ist und
eine **begrenzte** Warteschlange davor liegt (was `ADR-0060` als „kein Puffer"
verneinte). Dieselbe Trennung zieht
[`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md) für `Last-Event-ID`
— protokoll-natives Replay bleibt ungenutzt; diese ADR ändert daran nichts und
bestätigt es.

## Frage 3 — Semantik-Festlegung (die Sätze für `SPEC-020` und den Godoc)

Die folgenden Sätze sind bereits in `SPEC-020` gepinnt (Zeile
*Zustellgarantie* neu gefasst, Zeile *Erzeuger-Blockade* neu); der
Port-Godoc in `internal/application/port/outbound/changestream.go` wird von der
Fixrunde auf denselben Stand gezogen (er sagt heute zu, ein „nicht
empfangender" Abonnent erhalte nichts nachgeliefert, und lässt die
Erzeuger-Blockade aus — F-1).

**`SPEC-020`, Zeile *Zustellgarantie* (neu):**

> keine (Fire-and-Forget, verlustbehaftet): ein Consumer, der nicht verbunden
> ist **oder langsamer liest als Changes eintreffen**, verpasst die betroffenen
> Nachrichten ersatzlos — je Abonnent trägt der `Broadcaster` eine begrenzte
> Empfangs-Warteschlange, deren Überlauf verworfen wird. Ein Replay innerhalb
> des Streams gibt es nicht; verpasste Changes bleiben über den bestehenden
> Lesezugriffsweg ([`LH-FA-REA-001`](../../spec/lastenheft.md) ff.) und die
> bestätigte Consumer-Position ([`LH-FA-CON-003`](../../spec/lastenheft.md)/005)
> nachholbar ([`LH-FA-SST-008`](../../spec/lastenheft.md) Boundary).

**`SPEC-020`, neue Zeile *Erzeuger-Blockade*:**

> keine: `Publish` blockiert nie auf einen Abonnenten — auch ein verbundener,
> gerade nicht lesender Consumer hält den Capture-Pfad nicht an. Der Aufruf
> liefert nur bei ungültigem Aufruf einen Fehler; ein Fehlschlag geht nicht in
> den Rückgabewert des Capture-Aufrufs ein (Zeile *Fehler bei Publish-Fehlschlag*).

**Port-Godoc (`Publish`), Ziel-Wortlaut für die Fixrunde:**

> `Publish` verteilt einen Change an alle zum Aufrufzeitpunkt registrierten
> Stream-Abonnenten (`ADR-0060` Teilfrage 3, `ADR-0066`). Der Aufruf blockiert
> **nie** auf einen Abonnenten: jeder Abonnent trägt eine begrenzte
> Empfangs-Warteschlange; liest er nicht schnell genug, werden die über sie
> hinausgehenden Changes für ihn verworfen (Drop-Newest, nicht nachgeliefert).
> Der Aufruf trägt keine Zustellgarantie: ein Abonnent, der nicht verbunden ist
> oder langsamer liest als Changes eintreffen, verpasst die betroffenen Changes
> ersatzlos — die Nachvollziehbarkeit bleibt ausschließlich beim
> Lesezugriffsweg (`LH-FA-REA-001` ff.) und der bestätigten Consumer-Position
> (`LH-FA-CON-003`/005). Sein Fehler wird an der Aufrufstelle abgefangen und darf
> die bereits erfolgte Persistierung oder Bestätigung nicht beeinflussen.

**Zum Godoc-Satz aus F-5.** Der Vertragssatz „nach der Abmeldung darf der
Aufrufer den gelieferten Kanal nicht mehr lesen — er wird nicht geschlossen"
bleibt in der Sache richtig und trägt mit dieser ADR seinen Grund ausdrücklich:
der Kanal wird **nie** geschlossen, weil ein Send auf einen geschlossenen Kanal
auch im nicht-blockierenden `select` panikte und ein schließendes `cancel`
einen gleichzeitigen `Publish`-Send träfe. Ziel-Wortlaut für `Subscribe`:

> Nach der Abmeldung werden dem Aufrufer keine weiteren Changes zugestellt; er
> liest den gelieferten Kanal nicht weiter. Der Kanal wird **nicht**
> geschlossen: ein Schließen träfe einen gleichzeitig laufenden Sende-Versuch in
> `Publish` und panikte, `select` schützt davor nicht. Etwaig noch in der
> Warteschlange liegende Changes werden mit dem Abonnenten verworfen.

## Frage 4 — Träger und Umsetzung

- **Folge-ADR: ja** —
  [`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md).
  `ADR-0060` Teilfrage 3 Option A („kein Puffer") und die Design-Zeile aus
  Teilfrage 2 sind **Accepted**; eine inhaltliche Änderung ist per
  `AGENTS.md` §3.5/Modul 8 ausgeschlossen. Die Korrektur ist eine **neue**
  Entscheidung, die die alte Klausel ablöst — eng gefasst (Muster
  `ADR-0048`/`ADR-0062`/`ADR-0063`/`ADR-0065`), nicht die ganze ADR. Eine
  reine Plan-/`SPEC-020`-Präzisierung genügt nicht, weil sie den Widerspruch
  zur Accepted-Klausel stehen ließe.
- **Umsetzung: Fix in `slice-069`** <!-- d-check:status-provenance -->, **kein
  neuer Folge-Slice.** Der `Broadcaster` ist in diesem Slice entstanden, er ist
  **nicht** in `done/`, und der Slice trägt bereits eine offene Rückkante
  (F-1, F-2 → Fixrunde). Der Fix bleibt innerhalb des bestehenden Liefer-Punkts
  „Adapter-Grundgerüst" (kein Liefer-Punkt-Zuwachs, kein Re-Cut nötig).
  `slice-070` <!-- d-check:status-provenance --> findet danach die
  nicht-blockierende Semantik vor; sein Plan wird nur nachgezogen, nicht neu
  geschnitten.
- **Kein neuer Registereintrag in diesem Zug.** Die Finding-Klasse zu F-1
  („Zustellsemantik nur halb gepinnt") gehört in die **Slice-Closure** §7 von
  `slice-069` <!-- d-check:status-provenance -->, wie der Reviewer sie bereits
  dorthin reicht (Modul 6: Eintrag bei der Slice-Closure; Beleg
  `evidence/review-slice-069.md`). Ein Registereintrag aus dem Architect-Zug
  heraus wäre eine zweite Schreibstelle für denselben Zähler.

---

## Was dieser Zug geändert hat — und was nicht

**Geändert:**

- `docs/plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md` (neu,
  Accepted) und `docs/plan/adr/README.md` (Zeile `ADR-0066`; Zeile
  `ADR-0060` mit dem Supersedes-Vermerk)
- `spec/pflichtenheft.md` — `SPEC-020`: Zeile *Zustellgarantie* neu gefasst,
  Zeile *Erzeuger-Blockade* ergänzt
- `docs/plan/planning/in-progress/slice-069-grpc-streaming-adapter-grundgeruest.md` — <!-- d-check:status-provenance -->
  §3 Implementer-Entscheidung zur `Publish`-Semantik korrigiert, §6-Risiko mit
  Ausgang, §2 um die Fixrunden-Punkte (Godoc/`SPEC-020`-Nachzug,
  Broadcaster-Regressionstest, F-2-Negativtest) ergänzt
- `docs/plan/planning/open/slice-070-grpc-capture-integration.md` — <!-- d-check:status-provenance -->
  §1/§3 auf den nicht-blockierenden `Publish` nachgezogen (keine caller-seitige
  Goroutine), Kennungs-Fehler `SPEC-019` → `SPEC-020` berichtigt
- diese Verdikt-Datei

**Nicht geändert:** `internal/**` (kein Code-Eingriff — der Broadcaster-,
Port- und Capture-Code wird von der Fixrunde geändert, nicht von diesem
Architect-Zug), `docs/plan/planning/welle-19.md` (der Wellen-Closure-Trigger
ist von diesem Verdikt unberührt: `slice-069` war ohnehin noch nicht in
`done/`), `spec/architecture.md`.

**Offen für den Implementer-Zug (Fixrunde `slice-069`):**

1. `internal/adapters/driven/grpcstream/broadcaster.go` — begrenzter Kanal in
   `Subscribe` (Paket-Konstante, Vorgabe `64`), nicht-blockierender Send mit
   Drop-Newest in `Publish`, Paket-/Funktions-Godoc auf die neue Semantik, kein
   Schließen des Kanals; Regressionstest „ein registrierter, nicht lesender
   Abonnent hält `Publish` nicht an".
2. `internal/application/port/outbound/changestream.go` — `Publish`-Godoc auf
   den Wortlaut aus Frage 3.
3. F-2 aus dem Review: Negativtest der Token-Konfigurationsgrenze
   (`classifyToken` bei leer konfiguriertem Token) analog
   `TestClassifyTokenLeereKonfiguration` im HTTP-Adapter.
4. F-4/F-5-Doku-Punkte: `SPEC-020`-Feldnamen-Quelle (`SPEC-002` →
   `model.Change`) berichtigen; Godoc-Satz zu `Subscribe` gemäß Frage 3.

**Offen für den Planner-Zug:**

- `welle-19` §6 trägt keine neue Grenze zu benennen — der Fix landet vor der
  Closure von `slice-069`; kein Drift-Eintrag nötig.
- Risiko-Ausgang in `slice-069` <!-- d-check:status-provenance --> §6 ist in
  diesem Zug gesetzt (**entfallen** — die Fixrunde dieses Slice beseitigt den
  Auslöser vor der Closure; siehe Plan §6).
