# Slice capture-transient-wiederholung: Klasse `transient` — Wiederholung mit begrenztem Backoff im Capture-Pfad

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — der Slice trägt keine Closure-Bedingung, die von
seiner DoD verschieden wäre. Er hat keine Kante zu einer offenen Welle; die
Roadmap entscheidet über seine Priorität.

**Bezug:** [`LH-QA-REL-002`](../../../../spec/lastenheft.md) (kontrollierter
Neustart), [`LH-QA-REL-001`](../../../../spec/lastenheft.md) (kein
Datenverlust, Persist-before-ACK),
[`LH-FA-ADM-003`](../../../../spec/lastenheft.md) (sichtbarer Fehlerzustand),
[`ADR-0023`](../../adr/0023-fehlerklassifikation.md) (Fehlerklassen),
[`ADR-0012`](../../adr/0012-at-least-once.md) (at-least-once),
[`ADR-0049`](../../adr/0049-replication-fehlerklassen-schwellen.md)
(Replikations-Fehlerklassen und Schwellen), Architect-Verdikt
[`architect-verdict-welle-backfill-bestand-lese-schritt`](../../../reviews/architect-verdict-welle-backfill-bestand-lese-schritt.md)
§8 (Slice C).

**Berührte Spec-Stellen:** [`SPEC-008`](../../../../spec/pflichtenheft.md)
(Zeile `transient`: „Erneut versuchen mit begrenztem Backoff“) — gelesen; ein
Nachzug der Zeile gehört zur Umsetzung, wenn die ADR eine Bedingung
präzisiert.

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Closure der Welle
[welle-backfill-bestand](../done/welle-backfill-bestand.md). **Datum:** 2026-09-25.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein vorübergehend nicht verfügbarer Quell- oder Speicher-Dienst
beendet den Container nicht beim ersten Fehler: der Capture-Pfad wiederholt
mit begrenztem Backoff, wie die Aktion der Klasse `transient` in
[`SPEC-008`](../../../../spec/pflichtenheft.md) sie verlangt, und endet erst,
wenn die Grenze erschöpft ist. Der Erfassungspfad endet auf jeden
Adapter-Fehler mit Prozess-Ausgang 1 (Ausgangslage des Slice, gemessen am
Kommentar an `Run` in `internal/bootstrap/wiring.go`); die Fortsetzung trägt der
Neustart durch den Aufrufer (Kommentar an `Run` in `internal/bootstrap/wiring.go`,
Handbuch-Abschnitt „Neustart nach einem Fehler“, `restart: "no"` in
`compose.yaml`).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Entscheidung über die Form der Wiederholung.** Wo sie liegt (Adapter,
  Capture Service, Composition Root), welche Grenzen sie trägt, welche Fehler
  `transient` sind und wie sie sichtbar ist, entscheidet der Architect in einer
  ADR vor dem Start (§4); der Slice setzt sie um. Ohne die ADR wäre die Form
  eine Auslegung des Implementers.
- **Ein Restart-Supervisor im Container oder eine Änderung von `restart:
  "no"`.** Der Neustart liegt beim Aufrufer
  ([`LH-QA-REL-002`](../../../../spec/lastenheft.md)); ob die Wiederholung im
  Prozess ihn ergänzt, klärt die ADR.
- **Wiederholung für die Klassen `permission`, `configuration`, `schema` und
  `storage`.** [`SPEC-008`](../../../../spec/pflichtenheft.md) verlangt dort
  einen sichtbaren Fehler ohne stillen Retry bzw. kein Source-ACK.
- **Der Backfill-Run.** Ein `transient`-Fehler des Runs endet ihn `failed`
  ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  Teilfrage 4 und 5, run-lokal); eine Wiederholung dort ist ein anderer
  Vorgang.

## 2. Definition of Done

- [x] Die Wiederholung im Capture-Pfad steht nach der ADR des Architects: ein
      `transient`-Fehler (Beispiel: SQLSTATE 55006, Slot noch aktiv, den
      `START_REPLICATION` nicht wiederholt) führt zu Wiederholungen mit
      begrenztem Backoff (Anfangswert, Obergrenze, Erschöpfung laut ADR); jede
      Wiederholung setzt an der zuletzt bestätigten Position fort
      (`confirmed_flush_lsn`,
      [`ADR-0012`](../../adr/0012-at-least-once.md)), es entsteht kein ACK vor
      der Persistierung
      ([`LH-QA-REL-001`](../../../../spec/lastenheft.md)); ist die Grenze
      erschöpft, endet der Prozess mit der Klasse `transient` und dem
      Fehlerzustand im Heartbeat
      ([`LH-FA-ADM-003`](../../../../spec/lastenheft.md)). *Zu belegen durch:*
      `make test` (Race-Detector) gegen Fakes mit deterministischer Uhr — je
      Grenze (Wiederholungszahl, Backoff-Folge, Erschöpfung) ein Test, dessen
      Mutation der Eingabeseite (Grenze verschoben) rot färbt; dazu
      `make test-replication` mit dem Fall „Slot noch aktiv“
      (`TestStreamRestartsOnExistingSlot`-Muster ohne Wartezeit auf die
      Freigabe im Test).
- [x] Die Klassen bleiben getrennt: ein `permission`-, `configuration`-,
      `schema`- oder `storage`-Fehler wird nicht wiederholt und endet den Prozess
      (Rückfall-Verhalten von `classifyRunError`). *Zu belegen
      durch:* `make test` — je Klasse ein Negativtest an seine Eingabe gebunden
      (Mutation: die Klasse als wiederholbar behandeln färbt rot).
- [x] Die Träger folgen: der Kommentar an `Run`, der die `transient`-Aktion als
      nicht getragen nennt, die Container-Vertrags-Zeile in `compose.yaml` und der
      Handbuch-Abschnitt „Neustart nach einem Fehler“ nennen die Wiederholung
      und ihre Grenze; die Zeile `transient` von
      [`SPEC-008`](../../../../spec/pflichtenheft.md) trägt die Bedingung der ADR,
      wenn sie sie präzisiert; die Handbuch-Version und die Änderungshistorie
      tragen eine Zeile. *Zu belegen durch:* Lesen der Träger und `make
      docs-check`.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: siehe dritter Liefer-Punkt (Handbuch, `harness/README.md`
      falls ein Lauf-Beleg entsteht).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure der nächsten Welle (die Roadmap führt
      [welle-transformationen](../done/welle-transformationen.md) unter *Offene
      Wellen*, das Ereignis kann eintreten; ein Slice ohne Welle wird von ihr
      mitgeprüft).

**Umfang:** M — Schätzung, nicht gemessen; die Größe hängt an der ADR: liegt die
Wiederholung im Adapter `receive`, bleibt der Zug klein; reicht sie in die
Composition Root, ist die Rückführung in §4 zu prüfen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Ort der Wiederholung nach der ADR (Adapter `internal/adapters/driving/replication/receive` oder `internal/bootstrap/wiring.go`) und Tests | update | Wiederholung mit Backoff, Erschöpfung, Klasse `transient`. |
| `internal/bootstrap/wiring.go` (`classifyRunError`, Kommentar an `Run`) | update | Abbildung der wiederholbaren Fehler auf `transient`; der Kommentar nennt die Aktion als getragen. |
| `compose.yaml`, `docs/user/benutzerhandbuch.md`, `spec/pflichtenheft.md` (Zeile `transient`) | update | Container-Vertrags-Kommentar, Handbuch-Abschnitt, Bedingung der Klasse. |
| `internal/adapters/driving/replication/receive/stream_test.go`, Tests der Composition Root | update | Store-Tier-Fall „Slot noch aktiv“, Fake-Tests mit Uhr. |
| `internal/adapters/driving/replication/receive/receive.go`, `walretention.go`, neu `serverfault_test.go` | update / neu | Fixrunde: `serverFault` klassifiziert Server-Fehler nach SQLSTATE (`ErrPermission` für 42501/Klasse 28, `ErrRejected` für nicht transiente Abweisungen); `Stream.Close` schließt eine nicht gestartete Verbindung. |
| `internal/bootstrap/wiring.go` (`runStreamWithRetry`, `retryableStreamError`, `classifyRunError`, Zyklus-Closure) | update | Fixrunde: Rücksetzung der Episode nach einem Zyklus von mindestens 30 s, WARN mit Versuchszähler, INFO bei Fortsetzung, `permission`/`ErrRejected` nicht wiederholt, `Stream.Close` bei Fehlern nach `NewStream`. |
| `internal/bootstrap/stream_retry_internal_test.go` | update | Fixrunde: Tests an die Werte der ADR (2 s, 30 s, 5 min, Faktor 2) und an Rücksetzung, Log-Inhalt, `permission` gebunden. |
| `internal/adapters/driving/replication/receive/receive.go` (`Config.OnStreaming`, `serverFault` mit 25006) | update | Fixrunde 2 (`ADR-0136`): der Stream meldet nach der Bestätigung von `START_REPLICATION` den Streaming-Beginn; 25006 gehört zur wiederholten Positivliste. |
| `internal/bootstrap/wiring.go` (`runStreamWithRetry`, neu `runStreamCycle`, `cycleStream`, `cycleService`) | update | Fixrunde 2: Stabilität ab dem Streaming-Signal statt ab Zyklus-Beginn (ohne Signal nie Rücksetzung), INFO „fortgesetzt“ aus dem Signal, Aufbau unter einer Frist von 30 s, `Run` am Prozess-Kontext; die Zyklus-Closure ist als `runStreamCycle` testbar herausgelöst (die Schließ-Pfade sind Gegenstand von Tests). |
| `internal/bootstrap/stream_cycle_internal_test.go`, `internal/adapters/driving/replication/receive/streaming_signal_test.go` | neu | Fixrunde 2: Tests der Schließ-Pfade und der Aufbau-Frist (Loopback-Listener ohne Antwort, netzlos), des Streaming-Signals und von `Stream.Close`. |
| `internal/bootstrap/stream_retry_internal_test.go`, `internal/bootstrap/replication_stream_retry_internal_test.go`, `internal/bootstrap/administration_startorder_internal_test.go`, `internal/adapters/driving/replication/receive/serverfault_test.go` | update | Fixrunde 2: Skript-Zyklus mit Aufbau-Dauer, Streaming-Dauer und Signal-Ausfall; Quelltext-Test der Weitergabe von `OnStreaming` und der Stellung von `stream.Run` in `runStreamCycle`; `serverFault`-Tabelle je SQLSTATE-Klasse und Grenzfall. |
| `compose.yaml` (Kommentar), `docs/user/benutzerhandbuch.md` (1.83), `spec/pflichtenheft.md` (`SPEC-008` Zeile `transient`) | update | Fixrunde 2: Träger der SQLSTATE-Auswahl, des Stabilitätsmaßes und der Aufbau-Frist. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „der Capture-Pfad endet
auf jeden Adapter-Fehler mit Ausgang 1“; beide Stände gemessen; die Befehle
stehen im Codeblock, der Implementer trägt Stand und Trefferzahl ein):**

```suchlauf
cd3a1c60 156 -n -i -E 'transient|Backoff|erneut versuchen|restart: "no"|Neustart nach einem Fehler|kontrollierte Fortsetzung|Ausgang 1' -- internal spec docs/user harness compose.yaml
diff 185 -n -i -E 'transient|Backoff|erneut versuchen|restart: "no"|Neustart nach einem Fehler|kontrollierte Fortsetzung|Ausgang 1' -- internal spec docs/user harness compose.yaml
cd3a1c60 2 -n -i -E 'endet auf jeden|jeden Adapter-Fehler|nicht wiederholt' -- internal docs/user harness
diff 1 -n -i -E 'endet auf jeden|jeden Adapter-Fehler|nicht wiederholt' -- internal docs/user harness
```

**Fixrunde 2 (`ADR-0136`) — bewegte Eigenschaften: „was eine Rücksetzung der
Episode trägt“ (Zählwort `30 s`, Beschreibung `stabil`/`gestreamt`, Symbol
`streamRetryStableAfter`) und „welche Server-Fehler wiederholt werden“
(Symbole `ErrRejected`/`SQLSTATE`, Beschreibung `Server-Abweisung`); Parent ist
der Stand vor der Fixrunde:**

```suchlauf
8ad34a01 55 -n -i -E '30 s|stabil|streamRetryStableAfter|gestreamt|setzt die Episode|erfolgreiche[rn]? Zyklus' -- internal spec docs/user harness compose.yaml
diff 61 -n -i -E '30 s|stabil|streamRetryStableAfter|gestreamt|setzt die Episode|erfolgreiche[rn]? Zyklus' -- internal spec docs/user harness compose.yaml
8ad34a01 99 -n -i -E 'Server-Abweisung|SQLSTATE|ErrRejected' -- internal spec docs/user harness compose.yaml
diff 112 -n -i -E 'Server-Abweisung|SQLSTATE|ErrRejected' -- internal spec docs/user harness compose.yaml
```

| Träger | Befund | Behandlung |
|---|---|---|
| Kommentar an `Run` (`wiring.go`), Container-Vertrags-Zeile (`compose.yaml`), Handbuch, `SPEC-008` | nachgezogen: alle vier nennen die Wiederholung und ihre Grenze (`ADR-0135`); Suchlauf 1: 156 (`cd3a1c60`) → 185 (`diff`, +29 Wiederholungs-Erwähnungen), gemessen mit `make suchlauf-nachmessen` | „jede Aussage ‚trägt dieser Pfad nicht' folgt der Wiederholung" ✓ (der `Run`-Kommentar, `compose.yaml` und das Handbuch nennen die Wiederholung) |
| Fixrunde 2: Träger der Rücksetzung der Episode (Handbuch „Neustart nach einem Fehler“, Kommentar an `streamRetryStableAfter` und `runStreamWithRetry`, Test-Kommentare der Schwelle) | Suchlauf 3: 55 (`8ad34a01`) → 61 (`diff`); Fundstellen der Aussage „Zyklus von mindestens 30 s“ bis „Fehler gestreamt“ nachgezogen auf „ab der Bestätigung von `START_REPLICATION`“; die Historienzeile 1.82 des Handbuchs bleibt als Chronik stehen | kein Träger der alten Messgröße außer der Historienzeile ✓ (`git grep` nach `gestreamt|erfolgreiche[rn]? Zyklus|Zyklus.*30 s` trifft nur Handbuch-Abschnitt, Historie, `wiring.go`, Test-Kommentare) |
| Fixrunde 2: Träger der wiederholten Fehlermenge (`SPEC-008` Zeile `transient`, Handbuch, `compose.yaml`-Kommentar, `Run`-Kommentar, `ErrRejected`/`serverFault`-Kommentare) | Suchlauf 4: 99 (`8ad34a01`) → 112 (`diff`); `SPEC-008` nennt `ErrRejected` und die SQLSTATE-Auswahl, `compose.yaml` und der `Run`-Kommentar nennen die Server-Abweisung außerhalb der Auswahl | jede Nennung der wiederholten Fehler folgt der Positivliste ✓; Suchlauf 2 (`diff` 1) trifft die neue Handbuch-Zeile „Nicht wiederholt werden Berechtigungsfehler …“, die das geltende Verhalten beschreibt, keinen Rest der alten Aussage |
| Test-Kommentare, die auf die Freigabe des Slots warten (`TestStreamRestartsOnExistingSlot`) | Suchlauf 2: 2 → 1 — beide Ursprungs-Treffer sind nachgezogen, der eine verbleibende Treffer ist die Handbuch-Zeile der Fixrunde 2 (oben); der `reportFault`-Kommentar (`wiring.go`) nennt jetzt „nicht wiederholbaren Adapter-Fehler und die Erschöpfung der Wiederholung" | Kommentar und Test folgen dem Verhalten ✓ |

## 4. Trigger

**Start** (`next` → `in-progress`): eine Architect-Entscheidung liegt vor —
eine ADR mit Status `Accepted` (Zeile im ADR-Index
[`docs/plan/adr/README.md`](../../adr/README.md)), die die Wiederholungsform
festlegt: Ort, Backoff-Grenzen, welche Fehler `transient` sind, Sichtbarkeit
während der Wiederholung, Ausgang bei Erschöpfung und Verhältnis zu `restart:
"no"`; **der Übergangs-Commit `next` → `in-progress` nennt sie**
(`BEO-PGC/start-trigger-ohne-uebergabe-artefakt`, offen, 1×). Kein weiterer
Slice in `in-progress/` (WIP-Limit 1). Keine Priorität gesetzt; die Roadmap
entscheidet.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls die Wiederholung
  in die Composition Root reicht und Adapter, Verdrahtung und Doku nicht in
  einem Review tragen — der abtrennbare Teil ist der Zug am Adapter `receive`
  (Fall „Slot noch aktiv“) als eigener Slice.
- `in-progress` → `open` (blockiert): falls die Wiederholung die
  Persist-before-ACK-Zusage berührt (ein wiederholter Start liest vor der
  Persistierung des Vorgängers) — Architect-Frage, keine Umsetzung auf Verdacht.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test` (Race-Detector) grün + ein
realer `make test-replication`-Lauf + Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Die Wiederholung verzögert die Sichtbarkeit eines dauerhaften Fehlers**
  (ein `configuration`-Fehler, der wie `transient` aussieht, wird
  wiederholt statt gemeldet). *Erwartet, zu belegen durch:* die
  Negativtests je Klasse (zweiter Liefer-Punkt) und die Klassifikationstabelle
  der ADR. **Ausgang: eingetreten und vor dem Merge geschlossen** — der Review fand
  (F-4) `permission`-Fehler des Servers pauschal als `ErrReplication` wiederholt;
  `serverFault` trennt sie seitdem (`ErrPermission`, `ErrRejected`), gebunden durch
  `TestRunStreamWithRetryKlassenEndenOhneWiederholung` und die Tabelle in
  `serverfault_test.go` (Verifikations-Report §2 Zeile 2: je Klasse `waits == 0`).
  Die wiederholte Menge ist eine Positivliste ([`ADR-0136`](../../adr/0136-capture-wiederholung-stabilitaetsmass-und-sqlstate-auswahl.md)).
- **Ein wiederholter Start liest Änderungen doppelt oder überspringt eine.**
  *Erwartet, zu belegen durch:* der Fortsetzungs-Test an `confirmed_flush_lsn`
  gegen die reale Instanz (`make test-replication`), Persist-before-ACK bleibt
  bindend ([`ADR-0012`](../../adr/0012-at-least-once.md)). **Ausgang: nicht eingetreten, Beleg schwach — weiter offen als
  `slice-capture-retry-realtest-belege-schaerfen`.** Persist-before-ACK ist
  unberührt (Capture Service und ACK-Adapter stehen nicht im Diff, der
  ACK-Port wird je Zyklus neu gesetzt; Verifikations-Report §2 Zeile 1). Der
  Realtest `TestRunStreamWithRetrySlotStillActive` lief real (`make
  test-replication`, Exit 0, `internal/bootstrap` 13,8 s gegenüber 1,6 s ohne
  DSN) und belegt die Lieferung der Folge-Change; er prüft weder SQLSTATE 55006
  als Ursache noch die Start-Position an `confirmed_flush_lsn`, und `>= 2`
  schließt Duplikate nicht aus (Verifikations-Report §5, F-8).
- **Die Grenzwerte sind Startwerte ohne Messung.** *Erwartet, zu belegen durch:*
  die ADR nennt sie als Startwerte mit Nachschärfe-Trigger, und der Bericht
  trennt Messung von Setzung (`BEO-PGC/backfill-adapter-startwerte-ohne-messung`,
  1×). **Ausgang: weiter offen → `BEO-PGC/backfill-adapter-startwerte-ohne-messung`.**
  Die Werte (2 s Anfang, Faktor 2, 30 s Obergrenze, 5 min Gesamtfenster, 30 s
  Stabilität und Aufbau-Frist) sind **gesetzt, nicht gemessen**; die ADRs nennen sie
  als Startwerte, der Nachschärfe-Trigger steht in
  [`ADR-0136`](../../adr/0136-capture-wiederholung-stabilitaetsmass-und-sqlstate-auswahl.md)
  §Re-Evaluierungs-Trigger (regelmäßig längerer Aufbau oder Flapping-Episoden,
  die das Fenster umgehen). Eine Messung im Betrieb liegt nicht vor.
- **Der Ort des Aufrufs `START_REPLICATION` wandert vor diesem Slice.**
  `slice-start-vorlauf-grenze` verlegt den
  Aufruf von `NewStream` nach `Stream.Run` (`ADR-0128` Festlegung 1); der Fehler
  „Slot noch aktiv“ (SQLSTATE 55006) entsteht danach in `Run`, nicht mehr im
  Verbindungsaufbau, und die Beispiele in §1, §2 und §3 (Adapter `receive`,
  `TestStreamRestartsOnExistingSlot`) sind am Stand nach jenem Slice zu lesen.
  *Erwartet, zu belegen durch:* der Suchlauf in §3 an beiden Ständen und die ADR
  zur Wiederholungsform, die den Ort am Start nennt. **Ausgang: entfallen** — der
  Slice `slice-start-vorlauf-grenze` liegt in `done/`; der Zyklus liest den Ort
  am Stand danach: `Stream.Run` sendet `START_REPLICATION` und meldet den
  Streaming-Beginn erst nach dessen erfolgreicher Rückkehr, bei Startfehler nie
  (`TestRunSignalsNothingWhenStartReplicationFails`; Review Fixrunde 2, N-1).
- **Kommentare an `Run` behaupten mehr, als der Code trägt**
  (`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`, 3×, verkörpert:
  Zusage-Klausel im Reviewer-Skill). *Erwartet, zu belegen durch:* der Reviewer
  fährt den zugesagten Pfad im Code nach. **Ausgang: eingetreten, vor dem
  Merge geschlossen** — der Reviewer fuhr den Pfad nach und fand (F-1, HIGH) die
  Zusage „ein erfolgreicher Zyklus setzt die Episode zurück" in Kommentar und
  Handbuch ohne Code; ein Scratchpad-Test zeigte die Erschöpfung nach einer
  Stunde Betrieb. Die Fixrunden trugen sie (Rücksetzung ab Streaming-Signal);
  die Probe hat gegriffen, Beleg: `docs/reviews/review-slice-capture-transient-wiederholung.md`
  Fixrunde 1 und 2.

## 7. Closure-Notiz

- **Was hat funktioniert:** Die Wiederholung liegt am Stream-Zyklus der
  Composition Root ([`ADR-0135`](../../adr/0135-capture-transient-wiederholung-stream-zyklus.md)),
  Adapter, Capture Service und ACK-Adapter blieben unberührt; damit blieb
  Persist-before-ACK ohne Rückfrage bindend. Die Mutationsprobe des Reviewers
  fand, was Lesen nicht fand: zwei Runden, in der zweiten elf Mutationen, zehn rot
  und eine grün (`m8`, N-5; Reviewer-Bericht, übernommen), dazu die eigene
  Verifier-Mutation, die N-5 unabhängig reproduzierte (Verifikations-Report §4).
  Der Verifier bestätigte die DoD-Zeilen 1–4 und 6 als getragen; der Realtest
  lief real (`make test-replication`, Exit 0).
- **Was ging anders als geplant:** Zwei Fixrunden statt keiner, und eine
  Folge-ADR: [`ADR-0135`](../../adr/0135-capture-transient-wiederholung-stream-zyklus.md)
  sagte „ein erfolgreicher Zyklus setzt die Episode zurück" zu, ohne die
  Messgröße zu nennen (F-1, HIGH: nicht implementiert). Die erste Fixrunde maß
  die Dauer des Zyklus; erst der Re-Review fand, dass ein hängender Verbindungsaufbau
  damit das Gesamtfenster aufhob (N-1) und dass die SQLSTATE-Auswahl in keiner ADR
  stand (N-2). [`ADR-0136`](../../adr/0136-capture-wiederholung-stabilitaetsmass-und-sqlstate-auswahl.md)
  legt beides fest (Streaming-Signal nach `START_REPLICATION`, Positivliste).
  Offen bleiben drei Belege-Schwächen ohne DoD-Bruch (N-5, N-6, F-8), übergeben
  an die Folge-Slices unten.
- **Steering-Loop-Eintrag (Lerneintrag):** benannte Lücke im Architect-Übergabe-Artefakt:
  eine ADR, die eine **Rücksetzung oder Schwelle** zusagt, nennt die **Messgröße**
  (was gemessen wird, ab welchem Signal), nicht nur den Wert; ein „erfolgreicher
  Zyklus" ohne Signal ist eine Zusage ohne Eingabeseite, die nur Code-Lesen und
  Mutation fanden, kein Sensor. Der Eintrag ist gezählt, nicht verkörpert:
  `BEO-PGC/schwelle-ohne-benannte-messgroesse` (1×, angelegt mit diesem Slice). Die
  Spec-Lücke ist geschlossen, nicht offen: Die Zeile `transient` von
  [`SPEC-008`](../../../../spec/pflichtenheft.md) trägt die Bedingung der ADR.
- **Beobachtungs-Register (`../observations/`):** neu angelegt
  `BEO-PGC/schwelle-ohne-benannte-messgroesse/` mit `evidence/slice-capture-transient-wiederholung.md`
  (N-1, N-2); ergänzt `evidence/slice-capture-transient-wiederholung.md` in
  `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad/` (F-1, HIGH — Zähler 8×),
  in `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe/` (F-3 MEDIUM, N-5 LOW — Zähler 21×)
  und in `BEO-PGC/backfill-adapter-startwerte-ohne-messung/` (Grenzwerte gesetzt — Zähler 2×);
  `BEO-PGC/adapter-fehler-ausgang/` trägt den Ausgang `umgesetzt` mit diesem Slice als
  Anker. Keine weitere Beobachtung: F-8 (Realtest erzwingt den Fehlschlag indirekt) ist eine
  Ausprägung von `negativtest-ohne-bindung-an-seine-eingabe` unter LOW und steht nach der
  Deckel-Regel in dieser Notiz, ohne eigene Datei.
- **Folge-Slices:** `slice-capture-retry-aufbau-frist-bindung` (N-5: Wert der
  Aufbau-Frist in `Run` an einen Test binden) und `slice-capture-retry-realtest-belege-schaerfen`
  (F-8: Ursache 55006, `count == 2`, Start-Position gegen `confirmed_flush_lsn`;
  N-6: realer `pgconn` nach beendetem Aufbau-Kontext) — beide Dateien in `open/`.
  N-7 (Namensüberdeckung `cycleStream` im Test, INFO, keine Wirkung) bekommt keinen Slice.
- **Risiken aus §6:** je ein Ausgang, siehe §6 (Klassen: eingetreten, geschlossen ·
  doppeltes Lesen: weiter offen, Folge-Slice · Grenzwerte: weiter offen, Register ·
  Ort von `START_REPLICATION`: entfallen · Kommentar-Zusage: eingetreten, geschlossen).
- **Drei Paarungen:** dieser Slice hat keine Welle; die Prüfung läuft
  regelkonform bei der Closure der nächsten Welle.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Replication-Adapter und Composition Root sind keine
eigenen Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(Zähler gemessen am 2026-09-25 mit `ls evidence | wc -l` je Eintrag) —
`BEO-PGC/adapter-fehler-ausgang` (3×, Ausgang: dieser Slice),
`BEO-PGC/start-trigger-ohne-uebergabe-artefakt` (offen, 1×, einschlägig —
Start-Trigger), `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (3×,
verkörpert), `BEO-PGC/backfill-adapter-startwerte-ohne-messung` (offen, 1×,
Risiko §6), `BEO-PGC/kein-admin-weg-schema-fehler-recovery` (offen, 1×, nicht
einschlägig: dort geht es um die Klasse `schema`, nicht `transient`).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
