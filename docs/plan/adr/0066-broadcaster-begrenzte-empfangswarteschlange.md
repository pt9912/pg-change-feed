# ADR-0066: Broadcaster-Zustellung — begrenzte Empfangs-Warteschlange je Abonnent

**Status:** Accepted — Supersedes [`ADR-0060`](0060-grpc-streaming-mechanismus.md)
(nur die Puffer-Klausel: der Pro-Satz in Teilfrage 3 Option A
(„minimale Zustandshaltung im `Broadcaster` (kein Puffer …)") und die ihn
wiederholende Design-Zeile in Teilfrage 2 („ohne Empfänger wird die Nachricht
verworfen (kein Puffer, siehe Teilfrage 3)"); die übrigen Festlegungen dieser
ADR — Mechanismus/Wahl gRPC-Server-Streaming und Verwerfung von SSE/„NATS
erweitern" (Teilfrage 1), Anknüpfungspunkt am Capture-Pfad mit getrenntem
`ChangeStreamPort` (Teilfrage 2 abgesehen von der genannten Design-Zeile),
die Replay-Verwerfung selbst (Teilfrage 3 Option B/C bleiben verworfen),
Authn/Authz (Teilfrage 4), Adapter-Platzierung (Teilfrage 5) und
Bootstrap-Aktivierung (Teilfrage 6) — bleiben unverändert bestehen und werden
hier nicht wiederholt)

**Datum:** 2026-09-14

**Autor:** pt9912 (Architect-Rolle, unabhängiger Architect-Zug — anderer
Kontext als der Implementer-Lauf von `slice-069` <!-- d-check:status-provenance -->
und als der Reviewer-Lauf, der das Risiko als Entscheidung weiterreichte;
Modul 8 §Konflikt-Pfad, Verdikt 2: Folge-ADR statt stiller Korrektur)

**Bezug:** [`LH-FA-SST-008`](../../../spec/lastenheft.md) (Haupt-Bezug —
Live-Streaming, insbesondere das Boundary-Kriterium), [`LH-FA-REA-001`](../../../spec/lastenheft.md)
(bestehender Lesezugriffsweg als Nachhol-Pfad), [`LH-FA-CON-003`](../../../spec/lastenheft.md),
[`LH-FA-CON-005`](../../../spec/lastenheft.md) (bestätigte Consumer-Position),
[`LH-QA-REL-001`](../../../spec/lastenheft.md) (keine stillen **Daten**verluste),
[`ADR-0060`](0060-grpc-streaming-mechanismus.md) (Mechanismus und Port-/Adapter-
Design — superseded nur in der Puffer-Klausel), [`ADR-0061`](0061-http-sse-zusaetzlich-zu-grpc.md)
(zweiter Zustellweg über denselben `Broadcaster`; Teilfrage 3 dieser ADR
bleibt gültig, ihre Wiedergabe der Puffer-Klausel liest jetzt diese ADR),
[`ADR-0055`](0055-nats-change-notification-wecksignal.md) Punkt 1 (die
Nachvollziehbarkeit liegt beim Store, ein verlorenes Signal ist folgenlos) und
Punkt 4 (der Capture-kritische Pfad ist unantastbar),
[`ADR-0011`](0011-persist-before-ack.md),
[`ADR-0027`](0027-capture-application-service.md) (Persist-before-ACK-Kette
`Receive → Decode → Persist → COMMIT Store → ACK Source`),
[`ADR-0034`](0034-ports-nach-faehigkeiten.md) (Port-Zuschnitt nach Fähigkeit),
Review zu `slice-069` <!-- d-check:status-provenance --> (Anlass: F-1, F-5 und das
Broadcaster-Risiko),
`docs/plan/planning/in-progress/slice-069-grpc-streaming-adapter-grundgeruest.md` <!-- d-check:status-provenance -->
(Umsetzung: Fixrunde des noch offenen Slice), `docs/plan/planning/open/slice-070-grpc-capture-integration.md` <!-- d-check:status-provenance -->
(der Aufrufer, der die Zustellsemantik vorfindet)

**Schärft:** [`SPEC-020`](../../../spec/pflichtenheft.md) (die Zeile
*Zustellgarantie* — sie pinnt die blockierende Hälfte nach),
[`ARC-005`](../../../spec/architecture.md) (macht die gRPC-/SSE-Driving-Adapter
über den gemeinsamen In-Prozess-`Broadcaster` zustellseitig konkret)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0060`](0060-grpc-streaming-mechanismus.md) entschied Teilfrage 3
„Fire-and-Forget ohne Replay" und legte in Teilfrage 2 fest, dass der
Driven-Adapter `Broadcaster` (`internal/adapters/driven/grpcstream/`) ein
„In-Prozess-Fan-out **ohne Puffer**" ist. `slice-069` <!-- d-check:status-provenance -->
hat das wörtlich umgesetzt: `Subscribe()` liefert einen **ungepufferten**
Kanal, `Publish` übergibt jeden Change direkt per blockierendem Send
(`select { case sub.changes <- change: case <-sub.done: case <-ctx.Done(): }`).

`ADR-0060` prüfte nur den Fall **ohne** Empfänger („ein `Publish` ohne aktiven
Subscriber blockiert nicht und liefert keinen Fehler" — Fitness Function).
Der Reviewer-Lauf zu `slice-069` <!-- d-check:status-provenance --> hat den
zweiten Fall am Code nachvollzogen und als echt bestätigt: ein **registrierter,
gerade nicht lesender** Empfänger hält `Publish` an. Gebunden ist dieser Halt
allein durch `ctx` (Lebensdauer der Capture-Schleife, also des Prozesses) oder
die Abmeldung.

**Warum das den Capture-kritischen Pfad real erreicht.** `slice-070` <!-- d-check:status-provenance -->
bindet `Publish` nach der in [`ADR-0060`](0060-grpc-streaming-mechanismus.md)
Teilfrage 2 gezeichneten Kette **synchron** ein
(`Receive → Decode → Persist → COMMIT Store → ACK Source → Notify (best
effort) → Stream-Publish (best effort)`) und isoliert nur den **Fehler**, nicht
die **Zeit**. Der Grund für den Halt ist kein Defekt, sondern die normale
Betriebslage eines Server-Streamings: der gRPC-Handler liest den Kanal und ruft
`stream.Send` **in derselben Goroutine** (`internal/adapters/driving/grpc/server.go`,
`StreamChanges`); `stream.Send` blockiert, sobald das Flow-Control-Fenster des
Clients voll ist — also genau dann, wenn der Client nicht liest. In diesem
Zustand ist die Handler-Goroutine **nicht** auf Empfang geparkt, der
ungepufferte Send trifft keinen bereiten Leser, und `Publish` wartet. Über die
gemeinsamen Kanäle staut das zusätzlich **alle** Empfänger (Head-of-Line), weil
`Publish` die Abonnenten nacheinander bedient.

Damit berührt der ausgelieferte Zustand zwei Zusagen, die `ADR-0060` selbst als
Konstraint führt:

- **Der Capture-kritische Pfad ist unantastbar** ([`ADR-0011`](0011-persist-before-ack.md)/
  [`ADR-0027`](0027-capture-application-service.md); [`ADR-0055`](0055-nats-change-notification-wecksignal.md)
  Punkt 4 zieht dieselbe Grenze für das NATS-Wecksignal): Ein nicht lesender
  Client hält über `Publish` den `Receive` **jeder weiteren** Transaktion an —
  nicht den ACK dieser, aber den nächsten Empfang; der `cdc_capture_lag` wächst
  bis zum Stillstand. Das ist genau die Verzögerung, die der best-effort-Schritt
  nie auslösen darf.
- **`LH-FA-SST-008`s Boundary-Kriterium verlangt keinen Nachlieferungs-Zustand**
  („bleiben über den bestehenden Lesezugriffsweg abrufbar und über die
  bestätigte Consumer-Position fortsetzbar") — der ausgelieferte Code wartet
  aber gerade auf die Empfangsbereitschaft eines **verbundenen** Empfängers,
  was weder ein Akzeptanzkriterium verlangt noch semantisch zugesagt ist.

Der Reviewer hat korrekt festgehalten: `ADR-0060`s „kein Puffer" steht einer
Lösung nicht im Weg, solange die Lösung nicht dieselbe verworfene Sache unter
neuem Namen ist (Teilfrage 3 Option B/C — Verwerfen vs. **Nachliefern**); der
`ctx`-Parameter trägt die Schranke bereits. Diese ADR zieht die Konsequenz.

**Warum jetzt und nicht in `slice-070`.** Der `Broadcaster` ist in `slice-069` <!-- d-check:status-provenance -->
entstanden und noch nicht in `done/` — der Slice trägt bereits eine offene
Rückkante (Review-Findings F-1, F-2 → Fixrunde). Eine Änderung am Broadcaster
ist damit eine Änderung an einem **noch nicht geschlossenen** Slice, kein
Eingriff in Bestand; `slice-070` <!-- d-check:status-provenance --> findet dann eine Zustellsemantik vor, die den
Capture-Pfad strukturell nicht anhalten kann.

## Entscheidung

Wir wählen: **Die Entkopplung gehört in den `Broadcaster`.** Jeder Abonnent
erhält eine **begrenzte Empfangs-Warteschlange**; `Publish` übergibt jeden
Change mit einem **nicht-blockierenden Send** und einer expliziten
**Verwerfungsregel**. Der Erzeuger (`CaptureService`) hält nie auf einen
Abonnenten an. Der **Aufrufer bleibt unverändert** — `slice-070` <!-- d-check:status-provenance -->
ruft `Publish` weiter synchron in der Best-Effort-Kette und isoliert den
Fehler; eine caller-seitige Goroutine entfällt.

Vier Festlegungen:

1. **Senke darf die Quelle nie aufhalten — strukturell, nicht per Schranke.**
   `Publish` bedient jeden Abonnenten mit einem nicht-blockierenden Send
   (`select { case sub.changes <- change: default: … }`). Der Aufruf kehrt
   unabhängig davon zurück, ob ein Abonnent liest, abgemeldet ist oder sein
   Client das Flow-Control-Fenster geschlossen hat. Kein `done`-Signal und kein
   `ctx.Done()` im Handoff ist mehr nötig — es gibt nichts mehr zu warten; ein
   bereits beendetes `ctx` wird vorab geprüft. Der `Publish`-Fehler entsteht
   weiter nur aus einem ungültigen Aufruf (`ErrChangeStream`).
2. **Begrenzte Warteschlange je Abonnent statt ungepuffertem Kanal.**
   `Subscribe()` liefert einen Kanal mit begrenzter Kapazität
   (benannte Paket-Konstante, Vorgabe `64`); die Signatur
   `Subscribe() (<-chan *model.Change, func())` bleibt unverändert, die
   Kapazität ist für den Abonnenten unsichtbar. Die Warteschlange absorbiert
   einen Erzeuger-Burst, während der Abonnent zwischen zwei Reads steht (z. B.
   in einem blockierenden `stream.Send`).
3. **Explizite Verwerfungsregel: Drop-Newest.** Ist die Warteschlange voll,
   wird der **eintreffende** Change für diesen Abonnenten verworfen; bereits
   eingereihte Changes werden nicht verdrängt (keine Umsortierung). Die
   Warteschlange macht **keine** Zustellzusage — sie ist ein
   Entkopplungs-Puffer, kein Replay-Zustand: sie trägt keine Position, keine
   Abonnenten-Identität, wird bei `cancel` mit dem Abonnenten verworfen und
   steht **keiner** neuen Subskription zur Verfügung. Der Kanal wird weiterhin
   **nie geschlossen** — ein nicht-blockierender Send auf einen geschlossenen
   Kanal panikte ebenso wie ein blockierender, und ein `cancel`, das schließt,
   träfe einen gleichzeitigen `Publish`-Send.
4. **Semantik-Festlegung (Wecksignal-Charakter).** Die Zustellsemantik lautet:
   *Ein nicht verbundener **oder langsamer lesender** Empfänger verliert die
   betroffenen Nachrichten ersatzlos; der Erzeuger hält nie an.* Ein
   verlorenes Stream-Signal ist **kein Datenverlust**: die Nachvollziehbarkeit
   liegt beim `ChangeStorePort` ([`ADR-0055`](0055-nats-change-notification-wecksignal.md)
   Punkt 1), verpasste Changes bleiben über den bestehenden Lesezugriffsweg
   ([`LH-FA-REA-001`](../../../spec/lastenheft.md) ff.) und die bestätigte
   Consumer-Position ([`LH-FA-CON-003`](../../../spec/lastenheft.md)/005)
   nachholbar ([`LH-FA-SST-008`](../../../spec/lastenheft.md) Boundary,
   [`LH-QA-REL-001`](../../../spec/lastenheft.md) unberührt).

### Warum das nicht die verworfene Option B/C aus `ADR-0060` ist

Der Unterschied ist genau der, den die Verwerfung trägt: **Verwerfen statt
Nachliefern.**

- **Option B (Replay ab bestätigter Consumer-Position)** und **Option C
  (Ring-Puffer begrenzter Tiefe)** existieren in
  [`ADR-0060`](0060-grpc-streaming-mechanismus.md) Teilfrage 3, um einem
  **getrennten oder wiederkehrenden** Consumer etwas **nachzuliefern** —
  Option B garantiert es, Option C overbrückt kurze Verbindungslücken. Beide
  fügen eine Zustellzusage hinzu und brauchen Zustand, der einen **Connect**
  überlebt (Position, Ring-Inhalt).
- Die hier gewählte Warteschlange liefert **nie** nach: sie wird bei `cancel`
  verworfen, steht keiner neuen Subskription zur Verfügung und trägt keinen
  Positionsbegriff. Sie wirkt ausschließlich **innerhalb** einer bestehenden,
  **verbundenen** Subskription, um deren momentane Nicht-Lesebereitschaft zu
  absorbieren. Ein Client, der sich trennt und neu verbindet, erhält aus ihr
  **kein** Byte — er holt über den Lesezugriffsweg nach, genau wie bisher.

Die `ADR-0060`-Fitness-Function-Aussage „ohne Empfänger wird die Nachricht
verworfen" bleibt damit in Kraft; präzisiert wird nur, dass dieses Verwerfen
**nicht-blockierend** geschieht und eine begrenzte Warteschlange vor dem
Verwerfen liegt (was `ADR-0060` als „kein Puffer" verneinte).

### Verhältnis zu `ADR-0061`

[`ADR-0061`](0061-http-sse-zusaetzlich-zu-grpc.md) nutzt denselben
`Broadcaster` als zweiten Abonnenten und wiederholt in Teilfrage 3 dessen
„kein Puffer"-Zustand. Die **Entscheidung** jener ADR — Fire-and-Forget ohne
Replay, `Last-Event-ID` bewusst ungenutzt — bleibt unverändert gültig und wird
von dieser ADR **bestätigt** (kein Replay). Die von ihr zitierte
Puffer-Klausel stammt aus [`ADR-0060`](0060-grpc-streaming-mechanismus.md) und
liest sich mit dieser ADR als das, was sie meint: *kein Replay-Puffer*. Der
SSE-Weg erbt die begrenzte Empfangs-Warteschlange des gemeinsamen
`Broadcaster` unverändert — beide Zustellwege bleiben für denselben Consumer
verhaltensgleich.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun: der `Broadcaster` bleibt ungepuffert/blockierend, `slice-070` <!-- d-check:status-provenance --> ruft `Publish` synchron auf | kein Code-Eingriff; `ADR-0060`s „kein Puffer" bleibt wörtlich erfüllt | ein verbundener, nicht lesender Client hält über den ungepufferten Kanal `Publish` und damit den `Receive` jeder weiteren Transaktion an (`cdc_capture_lag` wächst bis zum Stillstand) — genau die Verzögerung des Capture-kritischen Pfads, die [`ADR-0011`](0011-persist-before-ack.md)/[`ADR-0027`](0027-capture-application-service.md) und [`ADR-0055`](0055-nats-change-notification-wecksignal.md) Punkt 4 ausschließen; zusätzlich Head-of-Line-Blocking über alle Empfänger |
| B — nur der Aufrufer: `slice-070` <!-- d-check:status-provenance --> ruft `Publish` nicht-blockierend (eigene Goroutine je Change oder `Publish` mit kurz-Deadline-`ctx`) | kein Eingriff am `Broadcaster`; `ADR-0060`s Textkörper bleibt unberührt | „Goroutine je Change” ist unbeschränkt (jede Transaktion startet eine, die auf den gestörten Empfänger wartet) und wächst mit dem Rückstand; eine Deadline-`ctx` verlagert den Halt nur in ein Zeitfenster und legt **jedem** Publish eine Wartezeit auf, auch wenn kein Client gestört ist — und sie trifft die **falsche** Aussage: `Publish` bedient die Abonnenten nacheinander, ein blockierender Erster hält die Übrigen bis zum Ablauf; der Aufrufer kann die serielle Schleife im `Broadcaster` nicht überspringen, ohne selbst zu blockieren. Die Invariante läge zudem in jedem Aufrufer statt an einer Stelle |
| **C — im `Broadcaster`: begrenzte Empfangs-Warteschlange je Abonnent + nicht-blockierender Send mit Drop-Newest (gewählt)** | eine Stelle für alle Abonnenten und beide Zustellwege; **beide Driving-Adapter bleiben unverändert** (`changeSubscriber`/`Subscribe`-Signatur identisch); der Erzeuger kann strukturell nicht warten; kein Replay-Zustand, keine zweite Wahrheit | superseded eine Klausel einer `Accepted`-ADR (Teilfrage 3 Option A von [`ADR-0060`](0060-grpc-streaming-mechanismus.md)); ein begrenzter Speicher je Abonnent; ein langsamer Consumer verliert Stream-Nachrichten **still** (per Festlegung, kein Log je verworfenem Change) |
| D — asynchroner Versand je Abonnent mit eigener Goroutine und begrenzter Warteschlange | trennt Sender- und Empfänger-Fortschritt sichtbar | liefert dieselbe Garantie nur, wenn der Erzeuger→Warteschlange-Handoff selbst nicht-blockierend ist — dann ist es C plus eine zusätzliche Goroutine je Abonnent, die Lebensdauer und `cancel`-Aufräumen tragen muss; keine zusätzliche Zusage, mehr bewegliche Teile |
| E — Entkopplung im Driving-Adapter: der gRPC-/SSE-Handler zieht weiter aus einem **ungepufferten** Kanal, schiebt aber nicht-blockierend in eine eigene begrenzte Warteschlange und sendet daraus (`Broadcaster` bleibt wörtlich pufferlos) | `ADR-0060`s „kein Puffer" bliebe sogar am Wortlaut erfüllt | die Nicht-Blockade hängt dann an einer **Scheduling**-Eigenschaft (die Drainer-Goroutine ist *fast immer* auf Empfang geparkt, aber während sie in ihre eigene Warteschlange schiebt oder unter GC-Pause **nicht**) statt an einer strukturellen (Send ohne Warten); das Warteschlangen-Muster würde in **jedem** nachgelagerten Adapter dupliziert (gRPC jetzt, SSE in `slice-072` <!-- d-check:status-provenance -->, später ein Dritter) — dieselbe Duplikations-Sorge, die das Review bereits für die Token-Klassifikation (F-3) benennt |

**Zur Alternative A im Verhältnis zu dieser Entscheidung:** A ist nicht bloß
schwächer, sondern unzutreffend — sie benennt die Blockade nicht, sie lässt sie
stehen. Der Reviewer hat sie als **echt** bestätigt und nur ihre Auflösung der
Entscheidung offengelassen; diese ADR entscheidet sie.

## Konsequenzen

- Positiv: Der Capture-kritische Pfad ist strukturell nicht anhaltbar. Ein
  verbundener, nicht lesender Client — gleich ob gRPC-`Send` oder
  SSE-Socket-Schreibvorgang — kann `Publish` nicht mehr aufhalten; die
  Best-Effort-Kette nach `ACK Source` bleibt best-effort in **Zeit** und
  **Fehler**.
- Positiv: **eine** Stelle. Beide Driving-Adapter bleiben unverändert; das
  `changeSubscriber`-Interface und die `Subscribe`-Signatur aus
  [`ADR-0060`](0060-grpc-streaming-mechanismus.md) Teilfrage 2 gelten fort, kein
  Adapter-Paket importiert ein anderes (Richtungskonvention aus
  `spec/architecture.md` §1 unberührt, `.a-check.yml` unverändert).
- Positiv: kein Schritt in Richtung der verworfenen Optionen. Die
  Warteschlange fügt **keine** Zustellzusage hinzu und trägt keinen Zustand über
  eine Subskription hinaus — die Nachvollziehbarkeit bleibt ausschließlich beim
  `ChangeStorePort`.
- Negativ: ein Consumer, der langsamer liest als Changes eintreffen, verliert
  Stream-Nachrichten **still**. Das ist die zugesagte Semantik (Wecksignal-
  Charakter), kein Fehlerpfad — ein Log je verworfenem Change wäre bei einem
  dauerhaft nicht lesenden Client ein Log-Flut, deshalb **kein** solches Log.
  Wer den Verlust sichtbar zählen will, braucht eine eigene Metrik (nicht
  Gegenstand dieser ADR).
- Negativ: [`ADR-0060`](0060-grpc-streaming-mechanismus.md) und diese ADR
  müssen zusammengelesen werden, um den `Broadcaster`-Vertrag vollständig zu
  verstehen (Muster der übrigen engen Klausel-Korrekturen des Repos:
  `ADR-0048`, `ADR-0062`, `ADR-0063`, `ADR-0065`).
- Negativ: die Kapazität (`64`) ist eine gewählte Einstellgröße. Anders als
  [`ADR-0060`](0060-grpc-streaming-mechanismus.md) Teilfrage 3 Option C ist sie
  **keine** unvorhersagbare Zusicherung: sie verspricht nichts, sie weitet nur
  die Burst-Toleranz eines verbundenen Abonnenten. Eine echte Zusicherung
  (Backpressure mit begrenztem Lag **statt** Verwerfen) bräuchte eine eigene
  Folge-ADR (Re-Evaluierungs-Trigger).
- Folgepflicht: `slice-069` <!-- d-check:status-provenance --> (Fixrunde des
  noch offenen Slice, `internal/adapters/driven/grpcstream/`) setzt Festlegungen
  1–3 um: begrenzter Kanal in `Subscribe`, nicht-blockierender Send mit
  Drop-Newest in `Publish`, Paket-Kommentar/Godoc auf die neue Semantik, ein
  Regressionstest „ein registrierter, nicht lesender Abonnent hält `Publish`
  nicht an". Der Kanal wird weiterhin nie geschlossen.
- Folgepflicht: `internal/application/port/outbound/changestream.go` — der
  `Publish`-Godoc wird auf die Sätze aus §Entscheidung 4/§Semantik-Festlegung
  nachgezogen (er sagt heute zu, ein „nicht empfangender" Abonnent erhalte
  nichts nachgeliefert, und lässt die Erzeuger-Blockade aus).
- Folgepflicht: `SPEC-020` ([`SPEC-020`](../../../spec/pflichtenheft.md))
  pinnt die blockierende Hälfte der Zustellsemantik (Zeile *Zustellgarantie*
  plus eine Zeile *Erzeuger-Blockade*). Die Sätze stehen im vorausgehenden
  Architect-Verdikt zu Publish blockiert den Capture-Pfad.
- Folgepflicht: `slice-070` <!-- d-check:status-provenance --> bleibt bei
  synchronem `Publish` ohne caller-seitige Goroutine; sein Plan-Vermerk trägt
  den Verweis auf diese ADR.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Unit-Test (`internal/adapters/driven/grpcstream`, Whitebox) | Ein **registrierter, nicht lesender** Abonnent hält `Publish` nicht an: `Publish` kehrt zurück, obwohl niemand liest; liest der Abonnent danach, erhält er die zwischengespeicherten Changes in Reihenfolge; über die Kapazität hinausgehende Changes sind für ihn verloren und stauen die übrigen Abonnenten nicht | `make test` |
| Go-Unit-Test (`internal/adapters/driven/grpcstream`, Whitebox) | Ein `Publish` ohne aktiven Subscriber blockiert nicht und liefert keinen Fehler (bestehende Fire-and-Forget-Regression, aus [`ADR-0060`](0060-grpc-streaming-mechanismus.md), bleibt unverändert) | `make test` |
| Go-Unit-Test (`internal/application/usecase/capture`, Whitebox) | Ein `ChangeStreamPort`, dessen `Publish` nicht zurückkehrt (oder dessen Subscriber nicht liest), darf `Capture()` nicht anhalten; der Fehler-/Zeitisolations-Regressionstest aus `slice-070` <!-- d-check:status-provenance --> deckt beide Hälften | `make test` |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Zwei Trigger, beide benannt statt „permanent"; die Festlegungen 1 und 4 sind
von keinem der beiden betroffen:

1. **Ein Bedarf an Stream-internem Replay wird konkret benannt.** Dann greift
   der Re-Evaluierungs-Trigger von [`ADR-0060`](0060-grpc-streaming-mechanismus.md)
   unverändert: eine eigene Folge-ADR bewertet Option B/C neu. Diese ADR
   ändert daran nichts.
2. **Eine Zustellgarantie über Fire-and-Forget hinaus wird verlangt** (z. B.
   Backpressure mit begrenztem Capture-Lag **statt** Verwerfen, oder eine
   Metrik/Konfiguration der Warteschlangentiefe als öffentliche Zusage). Dann
   entscheidet eine Folge-ADR zwischen einer wachsenden Warteschlange mit
   Rückstau-Semantik und der heute gewählten Verwerfungsregel.

Sonst permanent — die Entkopplung im `Broadcaster` (Erzeuger hält nie an,
Empfänger verwirft bei Überlauf) gilt unabhängig von der Anzahl der
Zustellwege und vom Nachrichtenvolumen.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-14 | Accepted — Anlass: `slice-069` <!-- d-check:status-provenance -->-Review (F-1, F-5 und das Broadcaster-Risiko) bestätigt real, dass ein registrierter, nicht lesender Empfänger `Publish` und damit den `Receive` jeder weiteren Transaktion anhält; unabhängiger Architect-Zug entscheidet die Entkopplung im `Broadcaster` und superseded `ADR-0060`s Puffer-Klausel (Teilfrage 3 Option A, Design-Zeile Teilfrage 2) | Review zu `slice-069`, der Architect-Verdikt zu Publish blockiert den Capture-Pfad |
| 2026-09-18 | Zitat-Korrektur — `docs/reviews/**`-Pfade durch Kennung ersetzt (`ADR-0073`) | `c2bc868` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0066` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
