# Slice routing-nats-subjekt: NATS-Zusatz-Subjekt — jede geroutete Change erscheint zusätzlich auf `cdc.route.<source_id>.<ziel>`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-routing](../welle-routing.md).

**Bezug:** [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (Routing —
Haupt-Bezug), [`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Live-Streaming),
[`LH-FA-SST-006`](../../../../spec/lastenheft.md) (Gleichwertigkeit der Zugriffswege),
[`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
Folgepflicht 5 und Teilfrage 5 (zusätzliche Veröffentlichung),
[`ADR-0100`](../../adr/0100-nats-dritter-vollinhalts-zustellweg.md) (NATS-Vollinhalt,
fire-and-forget), [`ADR-0055`](../../adr/0055-nats-change-notification-wecksignal.md)
(Wecksignal — bleibt unberührt).

**Berührte Spec-Stellen:**
[`SPEC-024`](../../../../spec/pflichtenheft.md) (NATS-Vollinhalt, Zusatz-Subjekt);
gelesen, nicht geändert: [`SPEC-017`](../../../../spec/pflichtenheft.md)
(Wecksignal). Die Spec führt: der Slice setzt `slice-routing-spec-nachzug`
voraus und ändert sie nicht.

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-routing](../welle-routing.md).
**Datum:** 2026-10-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der NATS-Vollinhalts-Publisher veröffentlicht eine Change mit Ziel
zusätzlich auf `cdc.route.<source_id>.<ziel>` (derselbe Payload wie auf
`cdc.stream.<source_id>.<schema>.<table>`); eine Change ohne Ziel erscheint nur auf
dem bestehenden Subjekt; das Subjekt `cdc.stream…` und das Wecksignal bleiben
unverändert. Drei Liefer-Punkte:

- (A) **Publisher:** die zweite Veröffentlichung nach dem Muster des bestehenden
  Subjekt-Baus (reservierte Zeichen geprüft); ihr Fehlschlag bleibt lokal wie der
  der ersten (fire-and-forget, `ADR-0100` Teilfrage 3) und berührt weder den
  `CaptureService` noch die erste Veröffentlichung;
- (B) **Beleg am realen NATS-Server:** ein Abonnent auf dem Subjekt eines Ziels
  empfängt genau die Changes dieses Ziels; ein Zielname mit `-` und `_` ist ein
  einzelnes gültiges Subjekt-Token (die ADR führt das als *hergeleitet*);
  ein Abonnent mit Wildcard `cdc.route.<source_id>.*` empfängt jedes Ziel;
- (C) **Kosten der zweiten Veröffentlichung:** die ADR nennt sie *nicht gemessen*
  (Konsequenz und Re-Evaluierungs-Trigger); der Slice misst sie am Testcontainer-NATS
  oder kennzeichnet sie ausdrücklich als nicht gemessen mit dem benannten Trigger.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Änderung des Subjekts `cdc.stream…` oder seines Payloads** — Teilfrage 5 lässt es
  unverändert; ein Abonnent der bestehenden Form darf nichts davon merken.
- **Das Wecksignal** (`cdc.changes.<source_id>.<schema>.<table>`, leerer Payload) —
  bleibt ohne Ziel (Teilfrage 5, Entscheidung der ADR).
- **Authentifizierung oder Subjekt-Berechtigungen am NATS-Server** — der Zustellweg
  bleibt tokenbasiert (`ADR-0100`); eine Beschränkung je Ziel wäre eine eigene
  Entscheidung (das Ziel ist Auswahl, kein Zugriffsschutz).
- **Nachvollziehbarkeit über NATS** — fire-and-forget ohne Zustellgarantie; das
  persistierte Label ist über `cdc.changes` und `GET /changes` lesbar.
- **SDK- und Beispiel-Clients** — `slice-routing-sdk-beispiel-target`.
- **Der Beleg am laufenden System** — `slice-routing-e2e` (Rundlauf durch den
  Feed-Container); dieser Slice belegt am Adapter-Test und an `make test-notify`.

## 2. Definition of Done

- [ ] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) und
      [`LH-FA-SST-008`](../../../../spec/lastenheft.md) (A): eine Change mit Ziel
      erzeugt genau zwei Veröffentlichungen — `cdc.stream.<source_id>.<schema>.<table>`
      (unverändert) und `cdc.route.<source_id>.<ziel>` mit byte-gleichem Payload —, eine
      Change ohne Ziel genau eine; der Fehlschlag der zweiten Veröffentlichung
      beendet weder die Schleife noch verändert er die erste (Test mit einer
      Verbindung, die die zweite verweigert). *Zu belegen durch:* Adapter-Tabellentest
      (`make test`), Zählung der Veröffentlichungen je Change am Test.
- [ ] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (B): am realen NATS-Server
      empfängt ein Abonnent auf `cdc.route.<source_id>.<ziel>` nur die Changes dieses
      Ziels (mindestens zwei Ziele und eine Change ohne Ziel), ein Zielname mit `-` und
      `_` wird angenommen, ein Wildcard-Abonnent auf `cdc.route.<source_id>.*` empfängt
      beide Ziele, ein Abonnent der bestehenden `cdc.stream…`-Form und der des
      Wecksignals empfangen unverändert. *Zu belegen durch:* `make test-notify`
      (erweiterter Lauf mit gedruckter Zeile je Fall).
- [ ] **Kosten der zweiten Veröffentlichung** (C): gemessen — Veröffentlichungen je
      Sekunde (oder Zeit je 10 000) mit und ohne Ziel am Testcontainer-NATS, Median von
      fünf Läufen, die gedruckte Zeile steht im Bericht, ein abgeleiteter Wert ist als
      abgeleitet gekennzeichnet ([`AGENTS.md`](../../../../AGENTS.md) §3.12) — **oder**
      ausdrücklich als nicht gemessen geführt mit dem Re-Evaluierungs-Trigger der ADR
      („Messung der zweiten NATS-Veröffentlichung zeigt Druck am Publisher"). Das
      Ergebnis geht in die Handbuch-Adresse (`slice-routing-betriebsdoku`).
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-routing-nats-subjekt.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: `harness/README.md` (Beschreibung von `make test-notify`) nur,
      soweit der Suchlauf eine bewegte Beschreibung findet; das Benutzerhandbuch
      (Abschnitt „Zugriff über den NATS-Vollinhalts-Stream") bleibt unberührt —
      Adresse: `slice-routing-betriebsdoku` §2 (Zusatz-Subjekt, Kosten-Aussage mit
      Ursprung); die NATS-Clients der SDKs und Beispiele folgen
      `slice-routing-sdk-beispiel-target`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der
      Closure der Welle [welle-routing](../welle-routing.md) (die Roadmap führt sie
      unter *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driven/natsstream/publisher.go` | update | die zweite Veröffentlichung in `publish` (Subjekt-Wurzel `cdc.route.`, Prüfung der reservierten Zeichen wie beim bestehenden Subjekt); der Publisher hält die Quelle (`sourceID`), das Ziel kommt aus `model.Change.RouteTarget`. |
| `internal/adapters/driven/natsstream/publisher_test.go` | update | Happy/Boundary/Negative nach [`LH-FA-CFG-008`](../../../../spec/lastenheft.md): Veröffentlichungen je Change, byte-gleicher Payload, Fehlschlag der zweiten, Ziel mit reservierten Zeichen (Verteidigung in der Tiefe, das Alphabet schließt es aus). |
| `tools/harness/run-notify-tests.sh` und der Testpaket-Teil, den das Skript fährt | update | realer NATS-Server: Abonnenten auf Ziel-Subjekt, Wildcard, `cdc.stream…`, Wecksignal; Kosten-Messung (Punkt C). |
| `harness/README.md` (Zeile `make test-notify`) | update, soweit der Suchlauf es findet | die Beschreibung nennt, was der Lauf belegt. |
| `docs/user/benutzerhandbuch.md` | **nicht** | Adresse `slice-routing-betriebsdoku` (§2). |
| `internal/adapters/driven/natsstream/publisher_nats_test.go` | neu | die Real-Server-Tests (`CDC_NATS_TEST_URL`, ohne Server übersprungen) und die Kostenmessung; eigene Datei, damit `publisher_test.go` netzlos bleibt. |
| `internal/adapters/driven/natsstream/publisher.go` (Nahtstelle) | update | die Verbindung des Publishers ist die kleine Schnittstelle `subjectPublisher` (`Publish`), `*nats.Conn` erfüllt sie; der Adapter-Tabellentest zählt damit die Veröffentlichungen je Change und lässt die zweite scheitern, ohne Server. `New` nimmt weiter `*nats.Conn`. |
| `tools/harness/run-notify-tests.sh` | update | ein zweiter `go test -v`-Lauf für `natsstream`, damit die gedruckten Zeilen im Lauf stehen. Fixrunde: der Lauf umfasst alle Tests des Pakets (kein `-run`-Präfix) und endet mit Exit 1, sobald ein Test sich überspringt (Review F-3). |
| `internal/adapters/driven/natsstream/publisher_test.go` (Fixrunde) | update | `TestRouteSubjectSurvivesSkippedTableSubject`: die Gegenrichtung der Unabhängigkeit (Tabellen-Subjekt nicht bildbar, Ziel-Subjekt erscheint); der `publish`-Godoc ist im Indikativ neu gefasst (Review F-1, F-2). |
| `spec/pflichtenheft.md` (Fixrunde) | update (fremde Datei, minimal) | `SPEC-024` Zeile Zusatz-Subjekt: „Last nicht gemessen" auf den gemessenen Umfang gezogen, ohne Zahl; Historienzeile (Review F-4). |
| `Makefile` (Hilfetext `test-notify`) | update | der Text nennt `natsstream` neben `natsnotify`. |
| `docs/plan/planning/open/slice-routing-betriebsdoku.md` | update (fremde Datei, minimal) | Übergabe-Block: Zusatz-Subjekt, erprobte Token-Aussage, gemessene Kosten-Aussage mit Ursprung (Punkt C: **gemessen**). |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „der NATS-Vollinhalts-Weg
veröffentlicht jede Change auf genau einem Subjekt `cdc.stream…`"; Parent ist
`30fd6cb5`; der Implementer ergänzt die `diff`-Zeilen und trägt Gefundenes und
Nichtgefundenes ein):**

```suchlauf
30fd6cb5 5 -n -E 'cdc\.stream\.' -- internal ':!*_test.go'
30fd6cb5 12 -n -E 'cdc\.stream\.' -- internal tools/harness test
30fd6cb5 64 -n -E 'cdc\.stream\.' -- spec docs/user sdks examples harness README.md
30fd6cb5 0 -n -E 'cdc\.route' -- internal tools test spec docs/user sdks examples harness README.md
30fd6cb5 7 -n subjectPrefix -- internal
diff 5 -n -E 'cdc\.stream\.' -- internal ':!*_test.go'
diff 22 -n -E 'cdc\.stream\.' -- internal tools/harness test
diff 66 -n -E 'cdc\.stream\.' -- spec docs/user sdks examples harness README.md
diff 20 -n -E 'cdc\.route' -- internal tools test
diff 3 -n -E 'cdc\.route' -- spec docs/user sdks examples harness README.md
diff 7 -n subjectPrefix -- internal
```

| Träger | Messung am Parent (`30fd6cb5`, gemessen am 2026-10-01) | Behandlung und Befund am Diff |
|---|---|---|
| Quellen, die das Subjekt-Schema führen | Zeilen 1 und 2: 5 Nicht-Test-Zeilen in `internal`, 12 mit Tests und `tools/harness`/`test` | jede Stelle gelesen: `publisher.go` (Paketkommentar, `subjectPrefix`, Typkommentar, `publish`-Kommentar) beschreibt das Tabellen-Subjekt, bleibt wahr (die Zusatz-Veröffentlichung steht im Typ- und `publish`-Kommentar); `tools/harness/natsstreamsub` und `run-integration-tests.sh` lesen `cdc.stream…` und sind unberührt. Diff: Nicht-Test 5 (unverändert), mit Tests 22 (die 10 neuen Zeilen sind die Tests dieses Slice) |
| Beschreibungen des Subjekt-Schemas in Spec, Handbuch, SDKs, Beispielen, Harness | Zeile 3: 64 Zeilen | Diff 66 (65 am Stand `72a59d5a` durch Commits nach dem Parent — die Spec-Zeile des Zusatz-Subjekts —, +1 `harness/README.md`-Zeile `make test-notify`, hier nachgezogen). Handbuch (9 Zeilen) → `slice-routing-betriebsdoku` (Übergabe-Block dort ergänzt); SDKs und Beispiele → `slice-routing-sdk-beispiel-target` (keine Änderung hier; ihre Subjekt-Beschreibungen bleiben wahr, sie kennen das Zusatz-Subjekt nicht). Nicht gefunden: eine Beschreibung „genau ein Subjekt je Change" in einem Träger außerhalb der Handbuch-Adresse |
| Namensraum `cdc.route` | Zeile 4: 0 Zeilen im Quell- und Doku-Suchraum | Diff: Quellen und Tests 20 Zeilen (3 `publisher.go`, die übrigen Tests dieses Slice, 19 am Stand vor der Fixrunde), Doku-Suchraum 3 (Spec-Zeilen und `harness/README.md`, nicht aus diesem Slice außer der README-Zeile); keine Kollision mit einem Bestand |
| Verwendung der Subjekt-Wurzel | Zeile 5: 7 Zeilen `subjectPrefix` | Diff 7 — Konstante und Nutzung unverändert; die neue Wurzel heißt `routeSubjectPrefix` und ist nicht Teil dieses Treffers |

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-routing-kern-label` liegt in `done/`
(das Feld `RouteTarget` an der Change, die der `Broadcaster` verteilt) und kein anderer
Slice liegt in `in-progress/` (WIP-Limit 1). Der Slice steht in der Tabelle nach
`slice-routing-lesewege` und braucht ihn nicht (Welle §5: untereinander unabhängig).
Voraussetzung am Start: ein Docker-Daemon mit dem gepinnten NATS-Image für
`make test-notify`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Kosten-Messung (C)
  wächst zu einer Benchmark-Infrastruktur — dann trennt sie sich ab
  (`slice-routing-nats-kosten`), (A) und (B) liefern; die ADR nennt die Messung als
  Re-Evaluierungs-Trigger, nicht als Voraussetzung der Umsetzung.
- `in-progress` → `open` (blockiert): ein Zielname mit `-` oder `_` wird vom realen
  Server nicht als einzelnes Token angenommen (die *hergeleitete* Aussage der ADR
  fiele) — Architect-Zug zu Teilfrage 1 vor der Weiterarbeit.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code
ungefiltert), `make test-notify` real grün, Suchlauf-Block nachgemessen,
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Zusätzliche Last am Publisher.** Jede geroutete Change erzeugt eine zweite
  Veröffentlichung; die Kosten sind nach der ADR *nicht gemessen*. — **Ausgang:** bei
  der Closure einzutragen (Messung mit gedruckter Zeile oder ausdrücklicher Verzicht
  mit Trigger).
- **Fire-and-forget gegen Erwartung eines Abonnenten.** Wer das Ziel-Subjekt abonniert,
  bekommt keine Zustellgarantie und keine Wiederholung (`ADR-0100` Teilfrage 3); ein
  Abonnent, der nach einem Verbindungsabbruch zurückkehrt, liest die Lücke über
  `cdc.changes` mit `route_target`. — **Ausgang:** bei der Closure einzutragen
  (Handbuch-Adresse).
- **Subjekt-Syntax *hergeleitet*.** Die ADR leitet aus der Syntax ab, dass das Alphabet
  des Zielnamens ein einzelnes Token bildet; am Server nicht geprüft. — **Ausgang:**
  bei der Closure einzutragen (Test am realen Server, Punkt B).
- **Ein Subjekt je Ziel, eine Quelle je Publisher.** Der Publisher hält genau eine
  `source_id`; mehrere Quellen im Prozess sind nach dem Bestand nicht vorgesehen
  (*erwartet*, `internal/adapters/driven/natsstream/publisher.go` `New`). —
  **Ausgang:** bei der Closure einzutragen.
- **Testcontainer-Abhängigkeit.** `make test-notify` braucht den gepinnten NATS-Digest
  und Docker; ein Pin-Wechsel ist ein bewusster Commit. — **Ausgang:** bei der Closure
  einzutragen.

## 7. Closure-Notiz

Wird bei der Closure gefüllt (vor dem `git mv` nach `done/`).

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit dem Pfad
`internal/adapters/driven/natsstream/` und dem Testskript `tools/harness/` — eine
Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** das Register wurde durchgegangen
(Zähler = Zahl der `evidence/`-Dateien, ausgezählt am 2026-10-01); die Treffer
stehen in §6 und in der Eröffnungs-Sichtung der Welle
[welle-routing](../welle-routing.md) §6:
`BEO-PGC/geschaetzter-wert-als-grenze`, `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
(verkörpert, 28×) und `BEO-PGC/adr-aussage-breiter-als-ihre-messung` (verkörpert, 10×)
— die Kosten-Aussage (Punkt C) trägt ihren Ursprung. Kein offener Eintrag erreicht mit
diesem Slice 3×.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

