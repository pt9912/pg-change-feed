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

- [ ] Die Wiederholung im Capture-Pfad steht nach der ADR des Architects: ein
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
- [ ] Die Klassen bleiben getrennt: ein `permission`-, `configuration`-,
      `schema`- oder `storage`-Fehler wird nicht wiederholt und endet den Prozess
      (Rückfall-Verhalten von `classifyRunError`). *Zu belegen
      durch:* `make test` — je Klasse ein Negativtest an seine Eingabe gebunden
      (Mutation: die Klasse als wiederholbar behandeln färbt rot).
- [ ] Die Träger folgen: der Kommentar an `Run`, der die `transient`-Aktion als
      nicht getragen nennt, die Container-Vertrags-Zeile in `compose.yaml` und der
      Handbuch-Abschnitt „Neustart nach einem Fehler“ nennen die Wiederholung
      und ihre Grenze; die Zeile `transient` von
      [`SPEC-008`](../../../../spec/pflichtenheft.md) trägt die Bedingung der ADR,
      wenn sie sie präzisiert; die Handbuch-Version und die Änderungshistorie
      tragen eine Zeile. *Zu belegen durch:* Lesen der Träger und `make
      docs-check`.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: siehe dritter Liefer-Punkt (Handbuch, `harness/README.md`
      falls ein Lauf-Beleg entsteht).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure der nächsten Welle (die Roadmap führt
      [welle-transformationen](../welle-transformationen.md) unter *Offene
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

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „der Capture-Pfad endet
auf jeden Adapter-Fehler mit Ausgang 1“; beide Stände gemessen; die Befehle
stehen im Codeblock, der Implementer trägt Stand und Trefferzahl ein):**

```text
git grep -n -i -E 'transient|Backoff|erneut versuchen|restart: "no"|Neustart nach einem Fehler|kontrollierte Fortsetzung|Ausgang 1' -- internal spec docs/user harness compose.yaml
git grep -n -i -E 'endet auf jeden|jeden Adapter-Fehler|nicht wiederholt' -- internal docs/user harness
```

| Träger | Befund | Behandlung |
|---|---|---|
| Kommentar an `Run` (`wiring.go`), Container-Vertrags-Zeile (`compose.yaml`), Handbuch, `SPEC-008` | *(Implementer trägt ein)* | jede Aussage „trägt dieser Pfad nicht“ folgt der Wiederholung |
| Test-Kommentare, die auf die Freigabe des Slots warten (`TestStreamRestartsOnExistingSlot`) | *(Implementer trägt ein)* | Kommentar und Test folgen dem Verhalten |

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
  der ADR. **Ausgang:** *(bei Closure)*
- **Ein wiederholter Start liest Änderungen doppelt oder überspringt eine.**
  *Erwartet, zu belegen durch:* der Fortsetzungs-Test an `confirmed_flush_lsn`
  gegen die reale Instanz (`make test-replication`), Persist-before-ACK bleibt
  bindend ([`ADR-0012`](../../adr/0012-at-least-once.md)). **Ausgang:** *(bei
  Closure)*
- **Die Grenzwerte sind Startwerte ohne Messung.** *Erwartet, zu belegen durch:*
  die ADR nennt sie als Startwerte mit Nachschärfe-Trigger, und der Bericht
  trennt Messung von Setzung (`BEO-PGC/backfill-adapter-startwerte-ohne-messung`,
  1×). **Ausgang:** *(bei Closure)*
- **Kommentare an `Run` behaupten mehr, als der Code trägt**
  (`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`, 3×, verkörpert:
  Zusage-Klausel im Reviewer-Skill). *Erwartet, zu belegen durch:* der Reviewer
  fährt den zugesagten Pfad im Code nach. **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen“ als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
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
