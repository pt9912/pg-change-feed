# Slice transformationen-e2e-wirkung: E2E Wirkung — Regel per SQL am laufenden Feed-Container, alle Zustellwege, Neustart, Ausschluss und Regel

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-transformationen](../welle-transformationen.md).

**Bezug:** [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) (Happy Path,
Boundary), [`LH-QA-SEC-004`](../../../../spec/lastenheft.md) (Wert nirgends im
Image), [`LH-FA-CFG-005`](../../../../spec/lastenheft.md) (Ausschluss gilt
zuerst), [`LH-FA-REA-005`](../../../../spec/lastenheft.md) (dieselbe Change
beim erneuten Lesen), [`LH-FA-SST-002`](../../../../spec/lastenheft.md)
(SQL-Lesezugriff), [`LH-FA-SST-006`](../../../../spec/lastenheft.md)
(HTTP-API), [`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Live-Streaming),
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Folgepflicht 5 (E2E-Belege) und §Fitness Function (realer Rundlauf),
[`ADR-0030`](../../adr/0030-testpyramide.md) (E2E-Tier).

**Berührte Spec-Stellen:** [`SPEC-019`](../../../../spec/pflichtenheft.md),
[`SPEC-020`](../../../../spec/pflichtenheft.md),
[`SPEC-021`](../../../../spec/pflichtenheft.md),
[`SPEC-022`](../../../../spec/pflichtenheft.md),
[`SPEC-024`](../../../../spec/pflichtenheft.md) — gelesen, nicht geändert (das
Nachrichtenschema bleibt: Live-Nachrichten
([`SPEC-020`](../../../../spec/pflichtenheft.md)/[`SPEC-021`](../../../../spec/pflichtenheft.md)/[`SPEC-024`](../../../../spec/pflichtenheft.md))
tragen zehn Felder, die HTTP-Antwort von
[`SPEC-022`](../../../../spec/pflichtenheft.md) dreizehn, `origin`
inbegriffen; gemessen an `readChangeResponse`,
`internal/adapters/driving/http/readchanges.go`).

**Verantwortlich:** Implementer-Agent, 2026-09-26.

**Autor:** Planner-Agent, Welle-Eröffnung
[welle-transformationen](../welle-transformationen.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Wirkung der Transformationen ist am **laufenden** Feed-Container
real belegt: eine Regel, per SQL beantragt und `applied`, prägt jede danach
erfasste Change auf allen Wegen, die sie liefern, in derselben Form; die
Konfliktfreiheit greift real (`failed` mit Text, Regelstand unverändert); der
Regelstand überlebt einen realen Neustart; der Ausschluss gilt vor der Regel.
Die Belege stehen in `make test-integration` und machen
[`LH-FA-CFG-007`](../../../../spec/lastenheft.md) im RTM-Lauf (`make
doc-trace`) zur Nicht-Waise (Träger: die vom Runner geschriebene Zeile in
[`docs/user/e2e-abdeckung.md`](../../../user/e2e-abdeckung.md)).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Nichtanwendbarkeit und Abhilfe** — `e2e-abhilfe`: der Erfassungspfad endet
  dort, das Szenario hat eine eigene Container-Ende-Grenze und braucht die
  deterministische Startreihenfolge aus `start-reihenfolge`.
- **Der Backfill-Beleg** — `backfill-pfad` trägt die Backfill-Phase mit Regel.
- **SDK-Realserver-E2E** — trägt `slice-sdk-regel-realserver-e2e` (Start nach
  diesem Slice, vor `betriebsdoku`): die SDK-Tiers (`make
  test-sdk-*-integration`) fahren dort je vier Phasen gegen eine aktive Regel;
  `make test-integration` trägt keinen SDK-Client.
- **Regeln je Consumer oder Zustellweg** —
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 6 Option D: alle Wege sehen dieselbe Form; ein Unterschied je Weg
  wäre ein Re-Evaluierungs-Trigger.
- **Rückwirkung auf gespeicherte Changes** — die zuvor erfasste Change behält
  ihre Rohform; das ist Teil des Belegs (Punkt 1), keine Fähigkeit, die dieser
  Slice liefert.

## 2. Definition of Done

- [ ] Happy Path und Boundary am laufenden Container: für eine per
      `cdc.enable_table` aktivierte, nicht in `CDC_TABLES` gelistete Tabelle
      beantragt `cdc.set_transformation` (`rename_column` und `map_value`), der
      Antrag wird `applied` (Poll auf `status`); die danach erfasste Change
      trägt über `cdc.changes`, `GET /changes` (`tools/harness/httpclient`),
      den gRPC-Stream (`grpcclient`), den SSE-Stream (`sseclient`) und den
      NATS-Vollinhalts-Stream (`natsstreamsub`) dieselbe transformierte Form
      (Zielname, abgebildeter Wert, Quellschlüssel fehlt); die zuvor erfasste
      Change bleibt in Rohform lesbar; je eine K1–K4-Verletzung endet real
      `failed` mit dem Text der Spec, und die danach erfasste Change trägt
      weiter die gültige Form. *Zu belegen durch:* ein realer, grüner `make
      test-integration`-Lauf; jeder Negativfall an seine Eingabe gebunden (die
      Gegenprobe mit gültigem Antrag endet `applied`).
- [ ] Neustart-Festigkeit und Ausschluss+Regel: nach einem **realen** `docker
      restart` leitet der Prozessstart den Regelstand aus den `applied`-Zeilen
      ab — die danach erfasste Change trägt die Form;
      `cdc.remove_transformation` stellt die Rohform für künftige Changes
      wieder her; `cdc.exclude_column` auf einer Spalte mit Regel: die danach
      erfasste Change trägt weder Quellnamen noch Zielnamen noch Wert im Image
      ([`LH-QA-SEC-004`](../../../../spec/lastenheft.md)), vor und nach dem
      Neustart, und die nicht ausgeschlossene Spalte bleibt darin. *Zu belegen
      durch:* derselbe `make test-integration`-Lauf; die Phase steht vor der
      Container-Ende-Grenze und vor dem Upgrade-Tausch des Runners (Lesen des
      Runners am Start).
- [ ] Die E2E-Abdeckung ist getragen: jede neue `func TestE2E*` und jede neue
      Runner-Phase trägt die Kennung
      [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) im Deklarations-Anker,
      wird von mindestens einem `-run`-Muster des Runners erfasst (Befehl im
      Bericht; `BEO-PGC/test-runner-stiller-ausschluss`, offen, 2×), und die
      Zeile steht im Runner-Erzeugnis
      [`docs/user/e2e-abdeckung.md`](../../../user/e2e-abdeckung.md); `make
      doc-trace` führt [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) nicht
      mehr unter den Waisen (am Stand nachgemessen, Zahl mit Ursprung,
      [`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A;
      [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) bleibt als Waise
      sichtbar). *Zu belegen durch:* `make doc-trace` (Ausgabe im Bericht) und
      `make docs-check`.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: `harness/README.md` §Sensors — die Zeile `make
      test-integration` nennt die neuen Belege, die Zeile `make doc-trace`
      trägt die nachgemessene Waisen-Aussage; das Benutzerhandbuch bleibt bis
      `betriebsdoku` unberührt (§1, Welle §4 Abweichung 3).
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
      der Closure der Welle
      [welle-transformationen](../welle-transformationen.md) (die Roadmap führt
      sie unter *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `test/integration/integration_test.go` | update | `TestE2E…`-Funktionen mit [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) (Go-Hälfte der Abdeckung: Happy, Boundary, Neustart, Ausschluss+Regel). |
| `tools/harness/run-integration-tests.sh` | update | Phasen für Happy Path (SQL, alle Wege), Boundary (K-Verletzung), Neustart, Ausschluss+Regel; `-run`-Muster; `abdeckung_declare`-Anker. |
| `tools/harness/httpclient`, `grpcclient`, `sseclient`, `natsstreamsub` | prüfen | Wegwerf-Clients geben Schlüssel und Werte der Row Images aus; trägt einer nur festen Feldzugriff, ist das ein Plan-Nachzug. |
| `docs/user/e2e-abdeckung.md` | Erzeugnis | kommt aus dem Runner, wird nicht von Hand geschrieben. |
| `harness/README.md` §Sensors (`make test-integration`, `make doc-trace`) | update | Aufzählung der Belege; Waisen-Messung nachgemessen. |
| `compose.yaml` | prüfen | keine Änderung erwartet: die Tabelle wird über `cdc.enable_table` aktiviert, `CDC_TABLES` bleibt unberührt. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „der Inhalt von
`make test-integration`“, „die Waisen im RTM-Lauf“, „die Zeilen der
E2E-Abdeckungstabelle“; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Beschreibung der Belege von `make test-integration` | `grep -rn 'Spaltenausschluss-Rundlauf' harness docs README.md` | *(Implementer trägt ein)* | die Zeile in `harness/README.md` §Sensors trägt die neuen Belege; weitere Aufzähler nachziehen |
| Waisen-Aussage im RTM-Träger | `grep -rn 'Waisen' harness docs` | *(Implementer trägt ein)* | Zahl und Kennungen am Stand nachmessen (`make doc-trace`), Ursprung nennen |
| Zeilenzahlen der Abdeckungstabelle | `git diff --stat` auf `tools/harness/run-integration-tests.sh` und `docs/user/e2e-abdeckung.md` | *(Implementer trägt ein)* | die Tabelle ist Erzeugnis; ein Runner-Eingriff verschiebt ihre Zeilennummern — neu erzeugen |
| Jede neue `TestE2E*`-Funktion ist von einem `-run`-Muster erfasst | `grep -n 'func TestE2E' test/integration/integration_test.go` gegen `grep -n -- '-run' tools/harness/run-integration-tests.sh` | *(Implementer trägt ein)* | Abgleich im Bericht; eine unerfasste Funktion ist ein Befund |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-transformationen-map-value` in
`done/` liegt (beide Regeltypen bestehen, `rename_column` und `map_value` sind
Teil des Happy Path), `slice-harness-fmt-check` in `done/` liegt (Kante:
`test/integration/integration_test.go` ist damit formatiert, `make fmt-check`
endet am Baum mit Exit 0; dieser Slice erweitert die Datei und lässt das Ziel
nach seinem Diff grün, Schritt 18 des Implementer-Ablaufs) und kein anderer Slice in `in-progress/` liegt
(WIP-Limit 1). `slice-backfill-e2e` liegt zu diesem Zeitpunkt in `done/`
(Voraussetzung von `backfill-pfad`); der Runner trägt die Backfill-Phasen
bereits, dieser Slice fügt seine eigenen an.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Happy Path,
  Neustart und Ausschluss+Regel nicht in einem Review tragen — der abtrennbare
  Teil ist der zweite Liefer-Punkt (Neustart, Ausschluss+Regel) als eigener
  Slice mit Start nach diesem.
- `in-progress` → `open` (blockiert): falls ein Wegwerf-Client die Schlüssel
  eines Images nicht ausgibt oder ein Zustellweg die Form abweicht (dann liegt
  ein Befund vor, der an `kern-rename` zurückgeht — ein Re-Evaluierungs-Trigger
  von
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 6 wird geprüft, nicht angenommen).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + ein realer, grüner `make
test-integration`-Lauf + Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Ein Zustellweg trägt eine abweichende Form.**
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  begründet Option D damit, dass alle Wege dieselbe Change sehen; belegt ist
  das erst am realen Lauf. *Erwartet, zu belegen durch:* der Happy Path über
  alle fünf Wege. **Ausgang:** *(bei Closure)*
- **Laufzeit von `make test-integration`** wächst
  (`BEO-PGC/test-integration-retention-timing-flake`, verkörpert, 3×); die
  Pipeline `e2e.yml` fährt das Testpaket unverändert (keine Workflow-Änderung,
  [`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht). *Erwartet, zu
  belegen durch:* ein realer Lauf, die Zahl trägt ihren Lauf. **Ausgang:**
  *(bei Closure)*
- **Eine neue Testfunktion fällt still aus dem Runner**
  (`BEO-PGC/test-runner-stiller-ausschluss`, offen, 2×). *Erwartet, zu belegen
  durch:* der Abgleich aus DoD Punkt 3. **Ausgang:** *(bei Closure)*
- **Geteilter Zustand zwischen Phasen**
  (`BEO-PGC/test-isolation-geteilter-zustand`, offen, 1×;
  `BEO-PGC/e2e-metrik-boundary-nur-reviewer-belegt`, offen, 1×): Regeln einer
  Phase wirken auf eine spätere Phase derselben Tabelle. *Erwartet, zu belegen
  durch:* eine eigene Tabelle je Phase und die Rücknahme der Regel im
  Aufräumen. **Ausgang:** *(bei Closure)*
- **Neustart-Beleg misst Timing statt Zustand.** Der Poll auf `applied` und die
  Wartezeit nach `docker restart` sind Zeitannahmen. *Erwartet, zu belegen
  durch:* Poll auf Zustand mit Frist, nicht feste Wartezeit. **Ausgang:** *(bei
  Closure)*
- **Das Handbuch entsteht erst danach** (Aufschub mit Adresse
  `slice-transformationen-betriebsdoku`, dessen §2 den aufgeschobenen
  Gegenstand nennt). **Ausgang:** *(bei Closure)*

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
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-transformationen](../welle-transformationen.md) (offen) — die Prüfung
  läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Test-Runner und Wegwerf-Clients sind keine eigenen
Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×, einschlägig — DoD Punkt
3), `BEO-PGC/test-isolation-geteilter-zustand` (offen, 1×) und
`BEO-PGC/e2e-metrik-boundary-nur-reviewer-belegt` (offen, 1×, Risiko §6),
`BEO-PGC/test-integration-retention-timing-flake` (verkörpert, 3×),
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 6×, Gegenprobe
in DoD Punkt 1), `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (verkörpert,
8×, jeder genannte Beleg-Befehl wird gefahren),
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert, 13×,
Waisen-Zahl), `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×,
Suchlauf §3).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
