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
voraus. Er ändert in der Fixrunde allein die Last-Aussage der Zeile
Zusatz-Subjekt von `SPEC-024` („nicht gemessen“ auf „ohne Abonnent und ohne
Schwelle gemessen, keine Last-Zusage“) und die Geschichte-Tabelle (Review F-4).

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

- [x] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) und
      [`LH-FA-SST-008`](../../../../spec/lastenheft.md) (A): eine Change mit Ziel
      erzeugt genau zwei Veröffentlichungen — `cdc.stream.<source_id>.<schema>.<table>`
      (unverändert) und `cdc.route.<source_id>.<ziel>` mit byte-gleichem Payload —, eine
      Change ohne Ziel genau eine; der Fehlschlag der zweiten Veröffentlichung
      beendet weder die Schleife noch verändert er die erste (Test mit einer
      Verbindung, die die zweite verweigert). *Zu belegen durch:* Adapter-Tabellentest
      (`make test`), Zählung der Veröffentlichungen je Change am Test. *Beleg
      (Verifier):* Verifikations-Report §2 Zeile 1 (`TestPublishCountsPublicationsPerChange`,
      `TestRouteFailureStaysLocal`, `TestRouteSubjectSurvivesSkippedTableSubject`; Mutationen
      M1, M2, M3 rot).
- [x] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (B): am realen NATS-Server
      empfängt ein Abonnent auf `cdc.route.<source_id>.<ziel>` nur die Changes dieses
      Ziels (mindestens zwei Ziele und eine Change ohne Ziel), ein Zielname mit `-` und
      `_` wird angenommen, ein Wildcard-Abonnent auf `cdc.route.<source_id>.*` empfängt
      beide Ziele, ein Abonnent der bestehenden `cdc.stream…`-Form und der des
      Wecksignals empfangen unverändert. *Zu belegen durch:* `make test-notify`
      (erweiterter Lauf mit gedruckter Zeile je Fall). *Beleg (Verifier):*
      Verifikations-Report §1 und §2 Zeile 2 (`make test-notify` Exit 0, gedruckte Zeilen
      je Fall; Mutationen N1 und N2 rot). **Grenze:** der Wecksignal-Abonnent ist nur
      negativ belegt (null Nachrichten, im Test läuft kein Wecksignal-Sender); „das
      Wecksignal funktioniert unverändert“ trägt der unberührte `natsnotify`-Adapter mit
      seinem Test im selben Lauf.
- [x] **Kosten der zweiten Veröffentlichung** (C): gemessen — Veröffentlichungen je
      Sekunde (oder Zeit je 10 000) mit und ohne Ziel am Testcontainer-NATS, Median von
      fünf Läufen, die gedruckte Zeile steht im Bericht, ein abgeleiteter Wert ist als
      abgeleitet gekennzeichnet ([`AGENTS.md`](../../../../AGENTS.md) §3.12) — **oder**
      ausdrücklich als nicht gemessen geführt mit dem Re-Evaluierungs-Trigger der ADR
      („Messung der zweiten NATS-Veröffentlichung zeigt Druck am Publisher"). Das
      Ergebnis geht in die Handbuch-Adresse (`slice-routing-betriebsdoku`). *Beleg
      (Verifier):* Verifikations-Report §1 (gedruckte Kostenzeile) und §2 Zeile 3:
      gemessen, Verhältnis mit/ohne über fünf Läufe 0,92 bis 1,30 (abgeleitet), kein
      Aufschlag auflösbar, keine Schwelle.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9). *Beleg
      (Verifier):* Verifikations-Report §1 (Exit 0, `coverage-gate` 82,00 %, `a-check`
      0 Befunde).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8). *Beleg:* Review
      [`review-slice-routing-nats-subjekt`](../../../reviews/review-slice-routing-nats-subjekt.md)
      (0 HIGH, 2 MEDIUM F-1 und F-2, 2 LOW F-3 und F-4, 2 INFO F-5 und F-6) und die
      **Gegenprüfung der Fixrunde `d8371372` durch den Verifier**
      ([`verifikation-slice-routing-nats-subjekt`](../../../reviews/verifikation-slice-routing-nats-subjekt.md)
      §5 und §7: F-1 bis F-4 geschlossen). Es gibt **kein separates Re-Review**.
      Begründung: die Fixrunde enthält keine geänderte Produktionslogik (der Diff von
      `publisher.go` ändert nur Kommentare), sondern Godoc, Tests, das Skript und die
      Abschwächung einer Spec-Zeile; das Skript hat der Verifier in drei Zuständen
      ausgeführt (grün, `FAIL`, `SKIP`), Mutation M3 färbt den neuen Test rot. Kein
      offenes HIGH/MEDIUM.
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-routing-nats-subjekt.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13). *Beleg
      (Verifier):* Verifikations-Report §1 und §2 Zeile 6 (11 Zeilen stimmen).
- [x] Doku-Update: `harness/README.md` (Beschreibung von `make test-notify`) nur,
      soweit der Suchlauf eine bewegte Beschreibung findet; das Benutzerhandbuch
      (Abschnitt „Zugriff über den NATS-Vollinhalts-Stream") bleibt unberührt —
      Adresse: `slice-routing-betriebsdoku` §2 (Zusatz-Subjekt, Kosten-Aussage mit
      Ursprung); die NATS-Clients der SDKs und Beispiele folgen
      `slice-routing-sdk-beispiel-target`. Die Fixrunde änderte zusätzlich die
      Last-Aussage in der `SPEC-024`-Zeile (Review F-4). *Beleg (Verifier):*
      Verifikations-Report §2 Zeile 7.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert (§7: drei neue `evidence/`-Dateien, ein
      `state.md` mit Vermerk).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
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
  Veröffentlichung; die Kosten sind nach der ADR *nicht gemessen*. — **Ausgang:**
  entfallen (gemessen). Die Zeit für `publish` samt Flush ist je 10 000 Changes mit
  Median von fünf Läufen gemessen, ohne Abonnent, ohne Schwelle, am Testcontainer-NATS;
  das Verhältnis mit/ohne (**abgeleitet**) liegt über fünf Läufe zwischen 0,92 und
  1,30 (Verifier-Zeile: „ohne Ziel 20,856 ms (479477 Changes/s), mit Ziel 19,272 ms
  (518898 Changes/s)“, Verifikations-Report §1; die übrigen vier Läufe stehen im
  Übergabe-Block von `slice-routing-betriebsdoku` §2, der Review-Lauf als
  **übernommen**). Ein Aufschlag der zweiten Veröffentlichung ist an diesem Messaufbau
  nicht auflösbar; eine Last-Zusage folgt nicht (Spec-Zeile `SPEC-024` entsprechend
  abgeschwächt). Der Re-Evaluierungs-Trigger der ADR („Messung zeigt Druck am
  Publisher“) bleibt.
- **Fire-and-forget gegen Erwartung eines Abonnenten.** Wer das Ziel-Subjekt abonniert,
  bekommt keine Zustellgarantie und keine Wiederholung (`ADR-0100` Teilfrage 3); ein
  Abonnent, der nach einem Verbindungsabbruch zurückkehrt, liest die Lücke über
  `cdc.changes` mit `route_target`. — **Ausgang:** weiter offen bis zur Handbuch-Stelle.
  Die Lokalität des Fehlschlags ist belegt (`TestRouteFailureStaysLocal`); die
  Aussage an den Anwender trägt der Übergabe-Block in `slice-routing-betriebsdoku` §2.
- **Subjekt-Syntax *hergeleitet*.** Die ADR leitet aus der Syntax ab, dass das Alphabet
  des Zielnamens ein einzelnes Token bildet; am Server nicht geprüft. — **Ausgang:**
  entfallen (erprobt). Abonnent auf `cdc.route.src-route.eu-west_1` empfängt zwei
  Nachrichten, `.*` und `.>` empfangen alle drei (`make test-notify` Exit 0,
  Verifikations-Report §1).
- **Ein Subjekt je Ziel, eine Quelle je Publisher.** Der Publisher hält genau eine
  `source_id`; mehrere Quellen im Prozess sind nach dem Bestand nicht vorgesehen
  (*erwartet*, `internal/adapters/driven/natsstream/publisher.go` `New`). —
  **Ausgang:** entfallen als Risiko dieses Slice: `New` ist textgleich zum Vor-Stand
  (Review, Negativbefund zu `publisher.go`); die Annahme bleibt *erwartet*, ohne Test.
- **Testcontainer-Abhängigkeit.** `make test-notify` braucht den gepinnten NATS-Digest
  und Docker; ein Pin-Wechsel ist ein bewusster Commit. — **Ausgang:** entfallen.
  `make test-notify` Exit 0 mit den Pins des Skripts (Verifikations-Report §1).
- **Überspringen aus anderem Grund färbt den Lauf rot.** Das Skript wertet jedes
  `--- SKIP` im Paket `natsstream` als Fehler; heute überspringt sich allein `realConn`
  ohne Server (Verifikations-Report §5 F-3). — **Ausgang:** entfallen als Gefahr (die
  Wirkung ist die sichere Richtung), als Grenze im Kommentar von
  `tools/harness/run-notify-tests.sh` benannt; der `natsnotify`-Lauf im selben Skript
  läuft ohne `-v` und ohne SKIP-Zweig (Verifikations-Report §8, vor dem Slice vorhanden).
- **Feste Wartezeiten der Real-Server-Tests.** `publisher_nats_test.go` wartet mit festen
  Zeiten (100 ms Wartezeit, 300 ms `NextMsg`-Frist je Drain); an einem langsamen
  Server ist das eine mögliche, **nicht beobachtete** Quelle von Intermittenz
  (*hergeleitet*, Review F-6). — **Ausgang:** weiter offen. Adresse:
  `internal/adapters/driven/natsstream/publisher_nats_test.go`; Trigger: die erste
  beobachtete Intermittenz von `make test-notify` — dann ein Folge-Slice, der auf die
  Nachricht wartet statt auf eine Zeit.

## 7. Closure-Notiz

- **Was hat funktioniert:** Die drei Liefer-Punkte tragen: der Publisher veröffentlicht
  eine Change mit Ziel zusätzlich auf `cdc.route.<source_id>.<ziel>` mit byte-gleichem
  Payload, ohne Ziel auf genau einem Subjekt; der Fehlschlag der zweiten
  Veröffentlichung bleibt lokal; die Nahtstelle `subjectPublisher` macht die Zählung der
  Veröffentlichungen je Change ohne Server testbar; am realen Server ist der Zielname mit
  `-` und `_` ein einzelnes Token (erprobt). Gemessen im Lauf des Verifiers: `make test`,
  `make test-notify`, `make gates`, `make docs-check`, `make fmt-check`, Suchlauf (11 Zeilen)
  alle Exit 0 (Verifikations-Report §1). Der Mutationsweg über beide Wege (Unit im Race-Image
  und `make test-notify`) färbte die Pflicht-Mutationen auf beiden Wegen rot.
- **Was ging anders als geplant:** Eine Fixrunde (`d8371372`). Der Review fand, dass die
  Unabhängigkeit der beiden Veröffentlichungen nur in einer Richtung gebunden war (F-1,
  MEDIUM) und dass der umgeschriebene `publish`-Godoc einen Konjunktiv-Satz über die
  verworfene Alternative mitführte (F-2, MEDIUM); das Skript hing am Namenspräfix der
  Tests (F-3) und die Spec-Zeile sagte „Last nicht gemessen“, obwohl der Slice sie misst
  (F-4). Die Fixrunde ergänzte `TestRouteSubjectSurvivesSkippedTableSubject` (die Eingabe
  der Gegenrichtung: Tabellen-Subjekt nicht bildbar), setzte den Godoc in den Indikativ,
  ließ das Skript das ganze Paket ohne `-run` fahren und bei jedem `--- SKIP` mit Exit 1
  enden, und zog die Spec-Zeile auf den gemessenen Umfang. **Kein separates Re-Review:** der
  Verifier hat die Fixrunde gegengeprüft (Diff von `publisher.go` nur Kommentare, das Skript
  in drei Zuständen ausgeführt, M3 rot); die Fixrunde enthält keine geänderte
  Produktionslogik, nur Godoc, Tests, Skript und die Abschwächung einer Spec-Zeile
  (Verifikations-Report §7). Das ist die Gegenprüfung durch einen zweiten Kontext, nicht
  ein Review. Nach der Verifikation (V-1 bis V-4): Plan-Kopf und DoD-Zeile an die
  Spec-Änderung angepasst, eine Grenzen-Zeile im Skript-Kommentar, die README-Zeile für die
  Wecksignal-Wurzel auf „empfängt nichts von der zweiten Veröffentlichung“ (negativ belegt),
  die Kosten-Spanne im Übergabe-Block auf fünf Läufe (0,92 bis 1,30, Ursprung je Lauf:
  zwei Implementer-Läufe gemessen, der Review-Lauf **übernommen**, ein Fixrunden-Lauf
  gemessen, der Verifier-Lauf gemessen im Verifikations-Report und hier **übernommen**).
  Die Mutationszahlen tragen ihren Ursprung (Instanz A von
  [`AGENTS.md`](../../../../AGENTS.md) §3.12): Implementer sieben, **übernommen** (Bericht,
  nicht nachgefahren); Reviewer sieben selbst gefahren, **gemessen** (sechs rot, eine grün —
  das ist F-1); Fixrunde zwei selbst gefahren, **gemessen**; Verifier sechs selbst gefahren,
  **gemessen** (alle rot; die sieben des Reviews übernahm er, soweit nicht wiederholt).
  Zahlen aus den Reports.
- **Steering-Loop-Eintrag:** geschärfte Regel, kein neuer Sensor. Eine
  Unabhängigkeits-Zusage zweier Seiteneffekte („das Überspringen der einen verändert die
  andere nicht“) ist zwei Zusagen und braucht beide Richtungen als Test, mit der Eingabe
  der Gegenrichtung als Eingabe (hier: leeres Tabellen-Subjekt bei gesetztem Ziel); ein
  Test, der nur eine Richtung bindet, lässt die Mutation der anderen grün (M6 des Reviews).
  Zweitens: ein Test, der nur mit einem Fremdsystem läuft, darf sich nicht still
  überspringen, und die Absicherung gehört an das Skript, das ihn fährt (Exit 1 bei
  `--- SKIP`, ganzes Paket ohne `-run`), nicht an eine Namenspräfix-Konvention, die niemand
  prüft: ein `go test -run` ohne Treffer endet mit Exit 0. Die Handlung vor dem
  Reviewer-Handoff: je Zusage mit „und umgekehrt“ im Wortlaut die Gegenrichtung mutieren
  (Schritt 19 in `.claude/commands/implement-slice.md` verlangt das schon: „je Zusage eine
  benannte Eingabeseiten-Mutation“). Träger: die Lese-Handlung des Reviewers und das Skript;
  keine neue Regel im Text von `AGENTS.md`.
- **Beobachtungs-Register (`../observations/`):** drei neue `evidence/`-Dateien, zwei
  Vermerke in `state.md`:
  - **`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`** (verkörpert, Deckel bei 14×):
    F-1 (MEDIUM, daher Datei trotz Deckel), die Gegenrichtung der Unabhängigkeit ohne
    Eingabe-Bindung; Zähler **22×** → **23×** (`ls evidence | wc -l`). Der Reviewer fand
    sie vor dem Merge, die Regel hat gewirkt.
  - **`BEO-PGC/vorher-nachher-sprache-in-test-harness-kommentar`** (verkörpert): F-2
    (MEDIUM, daher Datei), der Konjunktiv über die verworfene Alternative im
    umgeschriebenen `publish`-Godoc; Zähler **7×** → **8×**. Der Reviewer fand ihn vor dem
    Merge (Lese-Handlung); der Satz stand im Vor-Stand und wurde mit dem umgebrochenen Block
    mitgeführt. Ob der Kandidatenlauf von Schritt 20 ihn gemeldet hätte, ist nicht nachgemessen.
  - **`BEO-PGC/test-runner-stiller-ausschluss`** (offen): F-3 (LOW), der Lauf von
    `make test-notify` hing am `-run`-Namenspräfix; dieselbe Klasse wie die `-run`-Muster
    von `run-integration-tests.sh`; Zähler **2×** → **3×**, Schwelle erreicht, Ausgang beim
    Lese-Schritt der Closure von [welle-routing](../welle-routing.md). Für das Skript
    `run-notify-tests.sh` ist die Lücke geschlossen (ganzes Paket, Exit 1 bei Skip); für
    `run-integration-tests.sh` bleibt sie offen.
  - **`BEO-PGC/fixrunde-ohne-reviewer-lesung`** (offen, 3×): **vierter Gegenbeleg, keine
    Datei.** Die Fixrunde änderte eine Anweisung (das Skript), ein Re-Review blieb aus,
    weil der Verifier sie in drei Zuständen **ausgeführt** hat. Der Kandidat-Wortlaut
    „sobald die Fixrunde Anweisungen ändert“ ist damit zu weit; tragfähiger: „Produktionslogik
    oder Norm geändert, oder nach der Fixrunde hat kein anderer Kontext sie ausgeführt“.
    Im `state.md` vermerkt, die Entscheidung bleibt beim Lese-Schritt der Closure von
    [welle-routing](../welle-routing.md).
  - **Keine Beobachtung** (Begründung): F-4 (LOW, Spec-Zeile „Last nicht gemessen“) ist ein
    Träger-Nachzug nach §3.13 und im Slice gezogen, der Fall ist durch die geltende Regel
    gedeckt; F-5 und F-6 (INFO) stehen als Grenzen in §2 und §6.
- **Folge-Slices:** keine neuen. Übergaben: `slice-routing-betriebsdoku` (Handbuch-Stelle
  „Zugriff über den NATS-Vollinhalts-Stream“: Zusatz-Subjekt, Kosten-Aussage mit der Spanne
  0,92 bis 1,30, Wecksignal nur negativ belegt), `slice-routing-sdk-beispiel-target` (NATS-Clients),
  `slice-routing-e2e` (Beleg am laufenden System).
- **Risiken aus §6:** Last **entfallen** (gemessen, ohne Aufschlag auflösbar); Fire-and-forget
  **weiter offen** (Adresse: Handbuch-Stelle in `slice-routing-betriebsdoku` §2);
  Subjekt-Syntax **entfallen** (erprobt); eine Quelle je Publisher **entfallen** (unverändert);
  Testcontainer **entfallen**; Überspringen aus anderem Grund **entfallen** als Gefahr, als
  Grenze benannt; feste Wartezeiten **weiter offen** (Adresse `publisher_nats_test.go`,
  Trigger: erste beobachtete Intermittenz).
- **Drei Paarungen:** der Slice gehört zu [welle-routing](../welle-routing.md) (offen) — die
  Prüfung läuft bei deren Closure; die DoD-Zeile bleibt deshalb `[ ]`. (a) Anker: der
  Lerneintrag verkörpert nichts neu (Ausprägung unter Register und Reviewer-Skill); (b)
  Folge-Slice: keiner neu, die genannten Pläne liegen unter `open/`; (c) Register: die
  genannten Kennungen existieren als Verzeichnis, jede trägt ein nicht leeres `evidence/`.
- **Validator (Modul 8):** entfällt — das Zusatz-Subjekt ohne Handbuch, SDK-Clients und E2E ist
  für Betreiber noch nicht als Ganzes nutzbar; der Nutzer-Bedarf
  ([`LH-FA-CFG-008`](../../../../spec/lastenheft.md)) wird erst durch den Wellen-Beleg
  validierbar.

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

