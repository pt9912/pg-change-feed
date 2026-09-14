# ADR-0067: Publish-Einbindung an der Capture-Aufrufstelle — Fitness-Function-Zeile korrigiert

**Status:** Accepted — Supersedes [`ADR-0066`](0066-broadcaster-begrenzte-empfangswarteschlange.md)
(nur deren **dritte Fitness-Function-Zeile**: die Aussage „Ein
`ChangeStreamPort`, dessen `Publish` nicht zurückkehrt (oder dessen Subscriber
nicht liest), darf `Capture()` nicht anhalten; der Fehler-/Zeitisolations-
Regressionstest aus `slice-070` <!-- d-check:status-provenance --> deckt **beide Hälften**". Die übrigen
Festlegungen dieser ADR — §Entscheidung (synchroner Aufruf, keine
caller-seitige Goroutine), die Festlegungen 1–4, die verworfene Option B, die
beiden ersten Fitness-Function-Zeilen und die Re-Evaluierungs-Trigger — sind
**hiermit bestätigt** und bleiben unverändert bestehen; sie werden hier nicht
wiederholt)

**Datum:** 2026-09-14

**Autor:** pt9912 (Architect-Rolle, unabhängiger Architect-Zug — anderer
Kontext als der Implementer-Lauf von `slice-070` <!-- d-check:status-provenance -->,
der die Reduktion selbst benannte, und als der Reviewer-Lauf, der F-1 als
Entscheidung weiterreichte; Modul 8 §Konflikt-Pfad, Verdikt 1: die
Fitness-Function-Zeile ist der Defekt, die Entscheidung gilt)

**Bezug:** [`ADR-0066`](0066-broadcaster-begrenzte-empfangswarteschlange.md)
(§Entscheidung, §Verglichene Alternativen Option B, §Fitness Function —
superseded nur in der dritten Zeile), [`ADR-0060`](0060-grpc-streaming-mechanismus.md)
Teilfrage 2 (die Aufrufkette `Receive → … → ACK Source → Notify → Stream-Publish`),
[`ADR-0055`](0055-nats-change-notification-wecksignal.md) Punkt 4 (der
Capture-kritische Pfad ist unantastbar), [`ADR-0011`](0011-persist-before-ack.md),
[`ADR-0027`](0027-capture-application-service.md) (Persist-before-ACK-Kette),
[`LH-FA-SST-008`](../../../spec/lastenheft.md) (Happy Path, Boundary),
[`SPEC-020`](../../../spec/pflichtenheft.md) (Zeilen *Zustellgarantie*,
*Erzeuger-Blockade*, *Fehler bei Publish-Fehlschlag* — unberührt),
`docs/reviews/review-slice-070.md` <!-- d-check:status-provenance --> (Anlass:
F-1 und F-7), `docs/plan/planning/in-progress/slice-070-grpc-capture-integration.md` <!-- d-check:status-provenance -->
(§2-Fußnote, §3 Implementer-Abweichung), `.a-check.yml` (Schicht-Edges —
`app → adapters` fehlt), `internal/adapters/driven/grpcstream/broadcaster_test.go`
(`TestNichtLesenderAbonnentHaeltPublishNichtAn`), `internal/application/usecase/capture/service_test.go`
(`TestCapturePublishOhneRueckkehrHaeltKritischeKetteNichtAn`,
`TestCaptureKehrtOhneUndMitNichtLesendemStreamEmpfaengerZurueck`)

**Schärft:** — (Prozess-ADR ohne Spec-Stratum, wie `ADR-0063`/`ADR-0064`;
trifft eine Fitness-Function-Korrektur, ändert keine Lastenheft-/
Pflichtenheft-Zusage)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0066`](0066-broadcaster-begrenzte-empfangswarteschlange.md) (Accepted,
2026-09-14) verlegt die Entkopplung zwischen Erzeuger und Empfänger in den
`Broadcaster`: begrenzte Empfangs-Warteschlange je Abonnent, nicht-blockierender
Send, Drop-Newest. Ihre §Entscheidung zieht daraus die Konsequenz für die
Aufrufstelle ausdrücklich und eng: *„Der **Aufrufer bleibt unverändert** —
`slice-070` <!-- d-check:status-provenance --> ruft `Publish` weiter synchron in
der Best-Effort-Kette und isoliert den Fehler; eine caller-seitige Goroutine
entfällt."* Genau die caller-seitige Goroutine (oder ein `Publish`-Aufruf mit
kurz-Deadline-`ctx`) führt dieselbe ADR in §Verglichene Alternativen als
**Option B (verworfen)**.

Ihre **dritte Fitness-Function-Zeile** verlangt dagegen einen Testfall für den
Zustand, den dieselbe ADR an der Aufrufstelle ausschließt:

> Ein `ChangeStreamPort`, dessen `Publish` nicht zurückkehrt (oder dessen
> Subscriber nicht liest), darf `Capture()` nicht anhalten; der
> Fehler-/Zeitisolations-Regressionstest aus `slice-070` <!-- d-check:status-provenance --> deckt **beide
> Hälften**.

Drei Sätze, nicht gleichzeitig erfüllbar:

1. §Entscheidung: synchroner Aufruf, keine caller-seitige Goroutine, keine
   Deadline.
2. §Verglichene Alternativen Option B: verworfen.
3. §Fitness Function, dritte Zeile: verlangt für den Fall „`Publish` kehrt
   nicht zurück" genau die Zeit-Isolation der verworfenen Option B.

Mit einem synchronen `Publish`-Aufruf im Goroutinen-Kontext von `Capture()`
gibt es an dieser Aufrufstelle **keinen** Mechanismus, der `Capture()`
zurückkehren ließe, während `Publish` hängt. Die einzigen Kandidaten — eigene
Goroutine je Change, Deadline-`ctx` — sind exakt Option B. Die Zeile wäre nur
mit der Maßnahme grün, die dieselbe ADR verwirft.

**Die Klammer-Hälfte ist an dieser Schicht zusätzlich strukturell
unerreichbar.** „oder dessen Subscriber nicht liest" ließe sich so lesen, dass
der Capture-Test den **realen** `Broadcaster` anbinden soll. Das verbietet die
Schichtregel: `.a-check.yml` führt keine `app → adapters`-Kante, `a-check`
scannt `**/*.go` einschließlich `_test.go`; eine Probe-Datei mit
`grpcstream`-Import im Paket `capture_test` endet real mit `app-impurity`
(Exit 2, `make a-check` — im Review ausgeführt, F-7). Der Capture-Test kann
den realen Port an dieser Stelle nicht importieren. Wo die Nicht-Blockade real
belegt wird, ist das Paket `internal/adapters/driven/grpcstream`
(`TestNichtLesenderAbonnentHaeltPublishNichtAn`).

**Was an der Aufrufstelle real getragen ist.** Der ausgelieferte
`TestCapturePublishOhneRueckkehrHaeltKritischeKetteNichtAn`
(`internal/application/usecase/capture/service_test.go`) parkt `Capture()` über
einen vertragsbrüchigen Doppel **in** `Publish`, liest in der Parkposition
Store- und ACK-Stand ab und gibt erst danach frei. Er zeigt damit exakt die
Gegenrichtung der Fitness-Function-Zeile: `Capture()` hält an, solange
`Publish` nicht zurückkehrt — und die Capture-kritische Kette `Persist → ACK`
ist trotzdem vollständig. Das ist die Eigenschaft, die diese Schicht
**tatsächlich** trägt: die kritische Kette steht, **bevor** der erste Change
das Haus verlässt, nicht „ein vertragsbrüchiger Port hält `Capture()` nicht
an".

**Warum das ein Architect-Träger ist.** Der Defekt liegt im `Accepted`-Text
von [`ADR-0066`](0066-broadcaster-begrenzte-empfangswarteschlange.md), nicht am
Code: die Implementierung folgt §Entscheidung und ist gegenüber ihr
vollständig (Review F-1, §„Zum ADR-Widerspruch": legitime Abweichung). `AGENTS.md`
§3.5 schließt die In-place-Korrektur aus; Modul 8 §Konflikt-Pfad (Verdikt 1)
führt die Korrektur als Folge-ADR mit `Supersedes`. Der Slice darf nicht still
abschließen, solange die Zeile unrevidiert dasteht.

## Entscheidung

Wir wählen: **Die dritte Fitness-Function-Zeile wird ersetzt; die Entscheidung
der [`ADR-0066`](0066-broadcaster-begrenzte-empfangswarteschlange.md) —
synchroner `Publish`-Aufruf ohne caller-seitige Goroutine — wird bestätigt.**

Drei Festlegungen:

1. **Die Entscheidung gilt unverändert.** Der Aufrufer ruft `Publish` synchron
   in der Best-Effort-Kette, unmittelbar nach `ACK Source` und dem
   Notify-Schritt, und isoliert nur den **Fehler**. Eine caller-seitige
   Goroutine, ein Deadline-`ctx` oder jede andere Zeit-Isolation an dieser
   Stelle bleibt ausgeschlossen. Sie wird hiermit gegen drei Prüfungen
   bestätigt:
   - **gegen [`ADR-0055`](0055-nats-change-notification-wecksignal.md) Punkt 4
     und [`ADR-0011`](0011-persist-before-ack.md)/[`ADR-0027`](0027-capture-application-service.md):**
     die Capture-kritische Kette `Receive → Decode → Persist → COMMIT Store →
     ACK Source` darf durch keinen nachgelagerten Schritt verzögert werden. Die
     Entkopplung trägt der Port **strukturell** (nicht-blockierender Send,
     kein Warten auf Empfangsbereitschaft); die Aufrufstelle braucht deshalb
     keine eigene Maßnahme.
   - **gegen das, was `slice-069` <!-- d-check:status-provenance --> am Port
     belegt hat:** `TestNichtLesenderAbonnentHaeltPublishNichtAn` zeigt real,
     dass ein registrierter, nie lesender Empfänger `Publish` — auch über die
     Kapazität hinaus — nicht anhält. Die Zusage „der Erzeuger hält nie an"
     ist damit an der Stelle belegt, an der sie wirkt: am Port.
   - **gegen die Contras der verworfenen Option B:** eine Goroutine je Change
     ist unbeschränkt und wächst mit dem Rückstand; eine Deadline-`ctx` legt
     **jedem** `Publish` eine Wartezeit auf, auch ohne gestörten Client, und
     verlagert den Halt nur in ein Zeitfenster. Eine caller-seitige Maßnahme
     wäre überdies eine zweite Stelle, an der die Invariante liegen müsste.
2. **Die zutreffende Zusage ist geschichtet — Port und Aufrufstelle.** Die
   Nicht-Blockade ist eine **Vertragspflicht des Ports** und wird am Port
   belegt (Zeile 1 der Fitness Function aus
   [`ADR-0066`](0066-broadcaster-begrenzte-empfangswarteschlange.md), Paket
   `internal/adapters/driven/grpcstream`). Der **Anteil der Capture-Schicht**
   ist eine andere, dort prüfbare Eigenschaft: *die Aufrufstelle fügt kein
   eigenes Warten hinzu — die Capture-kritische Kette `Persist → ACK` ist
   abgeschlossen, **bevor** `Publish` betreten wird, und `Capture()` kehrt
   gegen einen vertragstreuen Port ohne eigene Zeit-Isolation zurück.*
   **Nicht** zugesagt ist: „ein `ChangeStreamPort`, dessen `Publish` seinen
   Vertrag verletzt und nicht zurückkehrt, hält `Capture()` nicht an" — bei
   synchronem Aufruf gibt es dafür konstruktiv keinen Mechanismus, und ein
   Test darauf wäre keine Deckung, sondern eine Behauptung.
3. **Der Ersatztext der dritten Fitness-Function-Zeile** lautet:

   > Die Aufrufstelle fügt kein eigenes Warten hinzu: (a) `Capture()` betritt
   > `Publish` erst, nachdem `Persist` und `Source-ACK` geschrieben sind (die
   > Capture-kritische Kette `Persist → ACK` ist vollständig, bevor der erste
   > Change das Haus verlässt), und ruft `Publish` synchron in der
   > Best-Effort-Kette; (b) `Capture()` kehrt gegen einen **vertragstreuen**
   > `ChangeStreamPort` ohne eigene Zeit-Isolation zurück — bei getrenntem
   > Client wie bei einem registrierten, nicht lesenden Empfänger (gegen einen
   > Port-Doppel, weil diese Schicht keinen Adapter importieren darf,
   > `.a-check.yml`). **Nicht** Gegenstand dieser Zeile: ein `ChangeStreamPort`,
   > dessen `Publish` nicht zurückkehrt — bei synchronem Aufruf ohne
   > caller-seitige Goroutine (§Entscheidung) gibt es an der Aufrufstelle
   > keinen Mechanismus, der `Capture()` zurückkehren ließe; die Nicht-Blockade
   > trägt der Port (Zeile 1) und wird in seinem Paket belegt.

   Die beiden ersten Fitness-Function-Zeilen der
   [`ADR-0066`](0066-broadcaster-begrenzte-empfangswarteschlange.md) bleiben
   unverändert in Kraft.

### Was diese ADR nicht ändert

- **`SPEC-020` bleibt unberührt.** Ihre Zeilen *Zustellgarantie*,
  *Erzeuger-Blockade* und *Fehler bei Publish-Fehlschlag* pinnen bereits die
  zutreffende Semantik („`Publish` blockiert nie auf einen Abonnenten … der
  Fehlschlag geht nicht in den Rückgabewert ein"). Der Widerspruch lag allein
  in der Fitness-Function-Zeile der ADR, nicht in der Spec.
- **Kein Code-Eingriff.** Der ausgelieferte Zustand ist gegenüber
  §Entscheidung vollständig; diese ADR ist keine Implementer-Arbeit.
- **`slice-070` <!-- d-check:status-provenance -->s §6-Risiken** (Publish vor
  `ACK`; Test-Duplikation) sind am Code nicht eingetreten und werden regulär
  bei Closure ausgegangen.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun: die Zeile bleibt, die Reduktion bleibt Implementer-/Plan-Notiz | kein Eingriff an einer `Accepted`-ADR | die `Accepted`-ADR trägt weiter eine maschinell prüfbare Aussage, die nichts erfüllt; der nächste Verifier-/Planner-Lauf muss den Widerspruch erneut herleiten — genau die stille Korrektur, die `AGENTS.md` §3.5 und Modul 8 §Konflikt-Pfad ausschließen (Review-Verdikt 3: „Zeile bleibt, Reduktion bleibt" ist kein Verdikt) |
| B — Fitness-Function-Zeile geschichtet korrigieren (Port = Nicht-Blockade, Aufrufstelle = kritische Kette vollständig vor `Publish` + kein eigenes Warten), §Entscheidung bestätigt **(gewählt)** | behebt den Defekt an seiner Stelle (der ADR-Text), keine Code-Änderung, kein Implementer-Rückweg; die beiden Hälften liegen je an der Schicht, die sie tragen kann; die Entscheidung wird ausdrücklich bestätigt, statt als offen stehen zu bleiben | superseded eine Zeile einer `Accepted`-ADR; die Formulierung muss die Klammer-Hälfte explizit ausgrenzen, sonst wiederholt sie den Fehler |
| C — Zeit-Isolation gehört doch an den Aufrufer: §Entscheidung und Option B dieser ADR werden `supersedes`d | die Zeile bliebe wörtlich erfüllbar; eine Frist-Isolation wäre an einer Stelle sichtbar | die Contras von Option B gelten unverändert (unbeschränkte Goroutine je Change, Wartezeit auf jedem `Publish`, serielle Schleife im `Broadcaster` nicht überspringbar) und die Invariante läge in jedem Aufrufer statt am Port; sie widerspräche zudem der strukturellen Entkopplung (`ADR-0066` Festlegungen 1–3), die der Port bereits trägt — eine echte Verhaltensänderung mit Implementer-Runde für ein Problem, das der Port schon löst |
| D — die Zeile ersatzlos streichen (die Eigenschaft trägt bereits Zeile 1 plus der Happy-Path-/Ordnungstest) | kleinster Eingriff; keine neue Formulierung nötig | die **Aufrufstelle** hätte keine benannte Fitness Function mehr für ihren eigenen Anteil — „kein eigenes Warten hinzufügen" und „kritische Kette steht vor `Publish`" wären ungeprüfte Zusagen; die Ordnung wird heute anwesenheits-, nicht grenz-getragen (der Parkpositionstest trägt sie) |
| E — die Zeile aus der Fitness-Function-Tabelle in Prosa überführen („nicht maschinell prüfbar") | ehrlich über die Grenze der Schicht | verschenkt den maschinell prüfbaren Anteil, den diese Schicht **hat** (kritische Kette vor `Publish`, kein eigenes Warten); eine Zusage ohne Fitness Function, obwohl eine existiert, ist schwächer als die korrigierte Zeile |

**Zur Alternative C im Verhältnis zu dieser Entscheidung:** Sie ist die einzige
echte Alternative mit Verhaltensfolge, und sie ist zulässig — aber sie wäre
eine **größere** Folge-ADR (sie superseded §Entscheidung und Option B der
[`ADR-0066`](0066-broadcaster-begrenzte-empfangswarteschlange.md), nicht nur
eine Zeile). Ihr Träger wäre ein neuer Implementer-Lauf. Diese ADR wählt sie
nicht: nichts am ausgelieferten Zustand noch an der Spec verlangt sie, und die
Entkopplung, die sie herstellen würde, trägt der Port bereits strukturell.

## Konsequenzen

- Positiv: Der Widerspruch innerhalb [`ADR-0066`](0066-broadcaster-begrenzte-empfangswarteschlange.md)
  ist aufgelöst — die eine Hälfte (Entscheidung) gilt ausdrücklich bestätigt,
  die andere (dritte FF-Zeile) ist korrigiert. Kein Lauf muss ihn erneut
  herleiten.
- Positiv: **kein Code-Eingriff, kein Implementer-Rückweg.** Der ausgelieferte
  Diff von `slice-070` <!-- d-check:status-provenance --> folgt §Entscheidung; die Korrektur ist reiner Text an der
  Entscheidungslage.
- Positiv: Die Zusage liegt je an der Schicht, die sie tragen kann — die
  Nicht-Blockade am Port (`internal/adapters/driven/grpcstream`), die
  kritische-Kette-vor-`Publish`-Eigenschaft an der Aufrufstelle
  (`internal/application/usecase/capture`, gegen einen vertragstreuen Doppel).
  Die Formulierung grenzt die an dieser Schicht **unerreichbare** Behauptung
  ausdrücklich aus.
- Negativ: [`ADR-0066`](0066-broadcaster-begrenzte-empfangswarteschlange.md)
  und diese ADR müssen zusammengelesen werden, um die Fitness Function des
  `ChangeStreamPort`-Aufrufs vollständig zu verstehen (Muster der übrigen
  engen Klausel-Korrekturen des Repos: `ADR-0048`, `ADR-0062`, `ADR-0063`,
  `ADR-0065`).
- Negativ: die Prüfung der Aufrufstelle bleibt gegen einen **Doppel**
  (vertragstreu), nicht gegen den realen `Broadcaster` — das ist die
  Schichtgrenze (`.a-check.yml`), kein Mangel dieser Zeile; der reale Beleg
  liegt im Port-Paket.
- Folgepflicht: `slice-070` <!-- d-check:status-provenance --> (Planner-Zug) —
  §2-Fußnote und §3-Implementer-Abweichung verweisen heute auf die
  unrevidierte dritte Zeile der [`ADR-0066`](0066-broadcaster-begrenzte-empfangswarteschlange.md)
  und formulieren ihren Inhalt als **Reduktion** („ist nicht als Test gebaut").
  Mit dieser ADR ist die zitierte Zeile ersetzt; der Plan soll auf
  [`ADR-0067`](0067-capture-publish-einbindung-fitness-function-korrektur.md)
  verweisen und die beiden gebauten Belege **als die zutreffende Zeile** nennen,
  nicht als Abweichung. Der Slice-Plan ist Planner-/Implementer-Arbeit; diese
  ADR benennt den Nachzug, ändert ihn nicht.
- Folgepflicht: die Finding-Klasse „Fitness-Function-Zeile gegen die eigene
  Entscheidung" (Review-Summary zu `slice-070` <!-- d-check:status-provenance -->)
  geht bei der Slice-Closure in §7 und von dort in den Zähler
  (Beobachtungs-Register) — Träger ist die Closure, nicht diese ADR.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Unit-Test (`internal/adapters/driven/grpcstream`, Whitebox) | Die Nicht-Blockade ist Vertragspflicht des **Ports**: ein registrierter, gerade nicht lesender Empfänger hält `Publish` nicht an — auch über die Warteschlangen-Kapazität hinaus (`TestNichtLesenderAbonnentHaeltPublishNichtAn`; unverändert die erste Zeile der [`ADR-0066`](0066-broadcaster-begrenzte-empfangswarteschlange.md)-Fitness-Function) | `make test` |
| Go-Unit-Test (`internal/application/usecase/capture`, Whitebox) | Die **Aufrufstelle** fügt kein eigenes Warten hinzu: (a) `Capture()` betritt `Publish` erst, nachdem `Persist` und `Source-ACK` geschrieben sind (`TestCapturePublishOhneRueckkehrHaeltKritischeKetteNichtAn`); (b) `Capture()` kehrt gegen einen vertragstreuen `ChangeStreamPort` ohne eigene Zeit-Isolation zurück, bei getrenntem Client wie bei einem registrierten, nicht lesenden Empfänger (`TestCaptureKehrtOhneUndMitNichtLesendemStreamEmpfaengerZurueck`). **Nicht** Gegenstand: ein vertragsbrüchiger `Publish`, der nicht zurückkehrt — an dieser Aufrufstelle konstruktiv nicht abfangbar (§Entscheidung) | `make test` |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Ein Trigger, benannt statt „permanent":

1. **Eine caller-seitige Zeit-Isolation wird doch verlangt** — etwa weil ein
   konkreter Bedarf an einer Frist an der Aufrufstelle benannt wird oder weil
   der Port seine Nicht-Blockade nicht mehr strukturell trägt. Dann entscheidet
   eine eigene Folge-ADR über §Entscheidung und Option B der
   [`ADR-0066`](0066-broadcaster-begrenzte-empfangswarteschlange.md) (Alternative
   C oben); die dort notierten Contras gelten dann unverändert.

Sonst permanent — die geschichtete Zusage (Nicht-Blockade am Port, kritische
Kette vor `Publish` an der Aufrufstelle) gilt unabhängig vom
Nachrichtenvolumen und von der Zahl der Zustellwege.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-14 | Accepted — Anlass: `slice-070` <!-- d-check:status-provenance -->-Review F-1 (HIGH) belegt den Widerspruch zwischen §Entscheidung, Option B und der dritten Fitness-Function-Zeile der [`ADR-0066`](0066-broadcaster-begrenzte-empfangswarteschlange.md); F-7 belegt real, dass die Klammer-Hälfte an der Capture-Schicht unerreichbar ist (`make a-check` `app-impurity`, Exit 2). Unabhängiger Architect-Zug bestätigt die Entscheidung (synchroner Aufruf, keine caller-seitige Goroutine) und ersetzt die dritte Fitness-Function-Zeile (Verdikt 1, Modul 8 §Konflikt-Pfad) | `docs/reviews/review-slice-070.md` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0067` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
