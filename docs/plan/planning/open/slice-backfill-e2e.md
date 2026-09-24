# Slice backfill-e2e: E2E — Happy Path, Boundary und Negative am laufenden Feed-Container; gemessene Startposition eines frisch registrierten Consumers

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-backfill-bestand](../welle-backfill-bestand.md).

**Bezug:** [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Happy Path, Boundary, Negative — alle drei am
laufenden System), [`LH-FA-CAP-004`](../../../../spec/lastenheft.md), [`LH-FA-REA-001`](../../../../spec/lastenheft.md), [`LH-FA-CON-005`](../../../../spec/lastenheft.md)
(Anfangsposition eines neu registrierten Consumers), [`LH-QA-POR-003`](../../../../spec/lastenheft.md) (E2E-Ebene),
[`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Festlegung 1 (Startposition messen und dokumentieren) und
Folgepflicht 2 (`S5`), [`ADR-0030`](../../adr/0030-testpyramide.md) (Testpyramide, E2E-Tier), [`ADR-0058`](../../adr/0058-testansatz-fuenf-luecken.md)
(Testansatz — additive Belege am realen Container), [`ADR-0012`](../../adr/0012-at-least-once.md)
(at-least-once, Consumer arbeiten idempotent), [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 2 (Aufnahme einer
`queued`-Zeile beim Prozessstart).

**Berührte Spec-Stellen:** [`SPEC-019`](../../../../spec/pflichtenheft.md), [`SPEC-022`](../../../../spec/pflichtenheft.md), [`SPEC-029`](../../../../spec/pflichtenheft.md) — gelesen
als Vertrag der Belege, nicht geändert.

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-backfill-bestand](../welle-backfill-bestand.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die drei Akzeptanzkriterien von [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) sind in
`make test-integration` **am laufenden Feed-Container** belegt, im Muster der
bestehenden Rundläufe (ausschließlich externe Wege: `docker exec`, SQL gegen
`cdc.changes`/`cdc.backfill_status`, HTTP — kein Import aus `internal/**`):

- **Happy Path.** Eine Tabelle mit Bestand wird aktiviert, `SELECT
  cdc.backfill_table(…)` beantragt den Backfill, ein Poll auf
  `cdc.backfill_status` sieht `completed`; jede zum Startzeitpunkt vorhandene
  Zeile ist über `cdc.changes` **und** `GET /changes` als Change lesbar
  (Tabelle, Operation `INSERT`, Row Image, `origin = 'backfill'`), unterscheidbar
  von einer danach WAL-erfassten Änderung derselben Zeile (`origin = 'wal'`); die
  `change_id`s sind unabhängig gegen `cdc.changes` gehalten.
- **Boundary.** Nebenläufige `INSERT`/`UPDATE`/`DELETE` an derselben Tabelle
  während des Runs; danach ergibt das Log ab Log-Anfang, angewandt in
  Lese-Ordnung (`INSERT`/`UPDATE` als Upsert des Row Images, `DELETE` als
  Löschen), je Schlüssel den Quellstand (**Replay-Invariante**); eine leere
  Tabelle endet `completed` mit 0 Zeilen und schreibt keine Transaktion; ein
  zweiter Antrag bei aktivem Run endet `failed` mit Grund.
- **Negative.** `docker kill` des Feed-Containers **mitten im Run**, Neustart:
  der Run steht `interrupted`, **keine** Change des Runs ist sichtbar, ein
  erneuter Antrag ergibt einen neuen Run, der `completed` erreicht — der Bestand
  ist danach **einmal** und vollständig lesbar (keine Dopplung durch den
  Abbruch). Eine zum Abbruchzeitpunkt `queued` wartende Zeile (ein zweiter Antrag
  gegen eine zweite Tabelle, angenommen hinter dem hängenden Run) überlebt den
  Neustart und wird beim Prozessstart aufgenommen und ausgeführt, ohne dass ein
  neuer Antrag nötig ist ([`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 2).
- **Startposition.** Von welcher Position ein frisch registrierter Consumer
  startet, ist in [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Festlegung 1 „erwartet, nicht geprüft": der Slice
  misst sie (Registrierung über den CLI-Weg, Position über den HTTP-Weg und
  `cdc.consumer_status`) und trägt das Ergebnis samt Lauf ins Handbuch ein — als
  Bezugsmuster, wie ein Consumer den Bestand erhält.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Fail-closed-Prüfung am Container** (`exclude_column` während des Runs) — die
  Negativ-Tests mit Fakes in `run-usecase` tragen sie; ein Zeitfenster-Test am
  Container wäre nicht deterministisch.
- **Retention-Beleg für Backfill-Changes** — die Regel gilt unverändert
  ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 7); ein Beleg über 24 Stunden `MinAge` wäre keiner
  Laufzeit tragbar.
- **Bench und Warn-Richtgröße** — `bench-richtgroesse`.
- **Die Live-Wege** — Backfill-Changes gehen nicht in gRPC, SSE oder
  NATS-Vollinhalt (Welle §6); ein Nicht-Beleg dort ist kein Kriterium.
- **Die SDK-Läufe** (`make test-sdk-*-integration`) — `sdk-origin`.

## 2. Definition of Done

- [ ] Happy Path und Boundary am laufenden Feed-Container: Bestand lesbar und als
      `backfill` erkennbar über beide Lesewege, unterscheidbar von einer
      WAL-Änderung derselben Zeile, Replay-Invariante nach nebenläufigen
      Schreibern, leere Tabelle, zweiter Antrag. *Zu belegen durch:* ein realer
      `make test-integration`-Lauf mit Exit 0; der Bericht nennt den Lauf mit
      den gedruckten `change_id`s. **Ort der Replay-Invariante:** die
      Fitness-Function-Zeile in [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) nennt den Tier `make test-replication`,
      Folgepflicht 2 den E2E-Beleg (`S5`); dieser Slice belegt sie im E2E, weil erst
      dort Snapshot, Schreiber und WAL-Pfad komponiert laufen — der Architect
      bestätigt die Ortswahl im Review oder verlangt zusätzlich den Tier-Beleg.
- [ ] Negative: `docker kill` im laufenden Run → nach dem Neustart `interrupted`,
      keine sichtbare Change des Runs, erneuter Antrag erreicht `completed`, Bestand
      einmal und vollständig; eine zum Abbruchzeitpunkt `queued` wartende Zeile
      überlebt den Neustart und wird ausgeführt (`completed`, Bestand ihrer Tabelle
      lesbar). *Zu belegen durch:* `make test-integration`; **jedes** dieser
      Kriterien trägt je eine Mutation im Bericht (die Prüfung gegen die
      Eingabe gelenkt, der Lauf färbt rot —
      `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`, verkörpert).
- [ ] Startposition gemessen und im Handbuch mit ihrem Lauf genannt; der
      Runner deklariert die Phase(n) über `abdeckung_declare` mit der Kennung
      [`LH-FA-CAP-009`](../../../../spec/lastenheft.md), und `docs/user/e2e-abdeckung.md` trägt nach dem Lauf eine
      Zeile dafür (Erzeugnis des Runners, kein Lauf-Beleg, mitcommittet);
      `make doc-trace` führt [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) nicht mehr unter den Waisen. Jede
      neue `func TestE2E*` steht in einem `-run`-Muster des Runners oder ist eine
      deklarierte Runner-Phase — *zu belegen durch:* der Lauf zeigt jede in der
      `-v`-Ausgabe (`BEO-PGC/test-runner-stiller-ausschluss`, offen, 2×).
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: Handbuch §4 (Abschnitt „Bestand als Backfill überführen“) trägt die gemessene Startposition mit Lauf-Ursprung; die Änderungshistorie eine Zeile.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel ·
      neuer Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in
      diesem Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      von der Closure der Welle [welle-backfill-bestand](../welle-backfill-bestand.md) (die Roadmap führt sie unter
      *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `test/integration/integration_test.go` | update | neue `TestE2E*`-Funktionen für Boundary/Negative-Teile, die Go-seitig lesen (Replay-Anwendung, Bestandsvergleich). |
| `tools/harness/run-integration-tests.sh` | update | Phasen: Auslösung per `docker exec`/`psql`, Poll, nebenläufige Schreiber, `docker kill` samt Neustart; `abdeckung_declare` mit `LH-FA-CAP-009`; `-run`-Muster um neue Testfunktionen. |
| Kommentar zum zweiten `make schema-rollout`-Lauf in `tools/harness/run-integration-tests.sh` (gemeldet von `slice-backfill-change-origin`, Review F-7) | update | der Kommentar („… real blockierten zweiten `make schema-rollout`-Lauf (… Exit 8 auf vier Fremdobjekten …)") beschreibt den Stand der Wache nicht mehr: sechs Fremdobjekte, der zweite Lauf endet mit Exit 0. Der Runner wird in diesem Slice ohnehin bearbeitet, und die Zeilen-Anker von `docs/user/e2e-abdeckung.md` verschieben sich hier ohnehin (Erzeugnis, mitcommittet) — der Slice, der den Kommentar ändert, ist der Slice, der die Anker regeneriert. |
| `tools/harness/httpclient/main.go` | prüfen/update | liest `GET /changes` für den Wegwerf-Beleg; trägt er `origin`, ist der Beleg der Feldform am Wire. |
| `compose.yaml` | prüfen | Vertrag des Feed-Containers (Tabellen, DSN); eine Abweichung wäre ein Plan-Nachzug. |
| `docs/user/e2e-abdeckung.md` | regeneriert | Erzeugnis des Runners, mitcommittet. |
| `docs/user/benutzerhandbuch.md` | update | gemessene Startposition (Ursprung: der Lauf), Änderungshistorie. |
| `harness/README.md` §Sensors | update | die Zeile `make test-integration` trägt den Backfill-Rundlauf (aus derselben Messung geschrieben). |

**Ansatz-Vorschlag, zu belegen (nicht bindend):** Ein „Abbruch mitten im Run" ist
nur mit einem festen Haltepunkt deterministisch. Der Slot-Anlage blockiert bis
eine beim Aufruf laufende Schreibtransaktion endet (M3 in [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md), gemessen
5,27 s bei rund 5 s Restlaufzeit): hält der Test eine offene Schreibtransaktion
auf der Quelle, steht der Run in `running`, ohne dass er weiterkommt — dann
`docker kill`, Transaktion beenden, Neustart. Das setzt voraus, dass der Run
vor der Slot-Anlage `running` trägt; der Implementer prüft es am
Use-Case-Vertrag. Für den `queued`-Beleg beantragt der Test, während der erste Run
hängt, einen zweiten Run gegen eine zweite Tabelle (derselben Tabelle würde der
zweite Antrag als `failed` enden): der Ein-Worker-Betrieb hält ihn `queued`; nach
dem Neustart beobachtet der Test, dass er ohne neuen Antrag `completed` erreicht.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „die Menge der E2E-belegten Kennungen und ihre Zeilen-Anker", „die Beschreibung von `make test-integration`"; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Abdeckungs-Tabelle | Lauf von `make test-integration`, `git diff docs/user/e2e-abdeckung.md` | *(Implementer trägt ein)* | mitcommitten; Zeilen-Anker früherer Zeilen verschieben sich |
| Zeile `make test-integration` in `harness/README.md` (lange Aufzählung der Rundläufe) | Lesen | *(Implementer trägt ein)* | Rundlauf ergänzen, Rest unverändert |
| RTM-Träger | `make doc-trace` | *(Implementer trägt ein)* | `LH-FA-CAP-009` nicht mehr Waise; sonst Deklarations-Anker prüfen |
| CI-Träger der Läufe | Lesen von `.github/workflows/e2e.yml` (Trigger, Matrix) | *(Implementer trägt ein)* | unverändert; die Laufzeit-Frage steht in §6 |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `sql-administration` in `done/` liegt,
`make image` real gelaufen ist (`compose.yaml` trägt keinen `build:`-Block,
[`ADR-0044`](../../adr/0044-image-beleg-semantik.md)) und kein anderer Slice in `in-progress/` liegt (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Happy Path/Boundary
  und Negative nicht in einem Review tragen — der abtrennbare Teil ist das
  Negative (dann als zweiter Slice `backfill-e2e-abbruch` mit Start nach
  diesem).
- `in-progress` → `open` (blockiert): falls ein Abbruch mitten im Run ohne
  Eingriff in die Produktion nicht deterministisch herstellbar ist (dann
  Architect-Frage nach einem Test-Haltepunkt im Run); falls die Replay-Invariante
  real **nicht** gilt (dann ein Befund gegen [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 3, keine
  Testanpassung).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test-integration` real grün (Exit
ungefiltert gesichert) + Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Die Replay-Invariante gilt real nicht** — [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) führt sie als
  „hergeleitet, im Slice als Eigenschaftstest zu belegen". Ein Rot wäre kein
  Testfehler, sondern ein Befund gegen die Entscheidung; die Rückführung §4
  benennt das. *Erwartet, zu belegen durch:* der Lauf. **Ausgang:** *(bei Closure)*
- **Nichtdeterministischer Abbruch** — siehe Ansatz-Vorschlag; ohne Haltepunkt
  bliebe der Test flackernd. *Erwartet, zu belegen durch:* mehrere
  Wiederholungen im Bericht. **Ausgang:** *(bei Closure)*
- **Randfall-Belege nur beim Reviewer** (`BEO-PGC/e2e-metrik-boundary-nur-reviewer-belegt`,
  offen, 1×): die leere Tabelle und der zweite Antrag sind committete
  Testfälle, nicht Reviewer-Scratch-Läufe. *Erwartet, zu belegen durch:* die
  Testfunktionen im Diff. **Ausgang:** *(bei Closure)*
- **Stiller Ausschluss aus dem Runner** (`BEO-PGC/test-runner-stiller-ausschluss`,
  offen, 2×): eine neue `TestE2E*`-Funktion ohne `-run`-Muster läuft nie.
  *Erwartet, zu belegen durch:* die `-v`-Ausgabe des Laufs. Ein weiterer
  Auftritt der Klasse erreicht 3×. **Ausgang:** *(bei Closure)*
- **Geteilter Zustand zwischen Rundläufen** (`BEO-PGC/test-isolation-geteilter-zustand`,
  offen, 1×): eigene Tabellennamen und eigene Quelle je Phase; die
  Bestands-Tabelle wird nach dem Lauf abgeräumt. *Erwartet, zu belegen durch:*
  Lesen der Phasen. **Ausgang:** *(bei Closure)*
- **Laufzeit des erweiterten Testpakets.** `e2e.yml` fährt `make test-integration`
  je PostgreSQL-Version der Matrix; der Zuwachs verlängert jeden Lauf. Der
  Workflow bleibt unverändert ([`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht). *Erwartet, zu
  belegen durch:* die lokale Laufzeit vor und nach dem Zug im Bericht (mit
  Lauf-Ursprung); eine Aussage über den GitHub-Runner ist erst nach dem ersten
  Push-Lauf möglich und bleibt **weiter offen**, falls ein Timeout auftritt.
  **Ausgang:** *(bei Closure)*
- **Die Startposition** ist ein gemessener Wert mit Ursprung (der Lauf); eine
  Übernahme aus [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) oder aus dem Gedächtnis ist nicht zulässig
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A). *Erwartet, zu belegen durch:* das Handbuch nennt
  den Lauf. **Ausgang:** *(bei Closure)*
- **Die Abdeckungs-Zeilen-Anker** verschieben sich mit jeder Einfügung oberhalb
  bestehender Phasen; die regenerierte Datei wird committet. **Ausgang:** *(bei
  Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen" als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice gehört zu [welle-backfill-bestand](../welle-backfill-bestand.md) (offen) — die
  Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); die E2E-Werkzeuge unter `tools/harness/` und `test/`
sind keine eigene Sub-Area — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×, einschlägig — DoD),
`BEO-PGC/e2e-metrik-boundary-nur-reviewer-belegt` (offen, 1×, einschlägig —
Risiko §6), `BEO-PGC/test-isolation-geteilter-zustand` (offen, 1×, Risiko §6),
`BEO-PGC/test-integration-retention-timing-flake` (verkörpert, 3× — Zeitfenster
in Polls tragen ihre Begründung), `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`
(verkörpert, 6×, DoD), `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
(verkörpert, 13×, Startposition mit Lauf),
`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (verkörpert, 8× — die
Zeile in `harness/README.md` beschreibt nur, was der Lauf gefahren hat),
`BEO-PGC/anforderung-ohne-erkennbaren-nachweis` (offen, 1×, einschlägig: dieser
Slice löst die Waise `LH-FA-CAP-009` auf), `BEO-PGC/limit-fortsetzung-innerhalb-einer-position`
(offen, 0×, einschlägig — die Belege lesen den Bestand ohne `Limit`),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf §3),
`BEO-PGC/github-actions-unverifizierbar-lokal` (verkörpert, 7× — nicht
einschlägig: kein Workflow-Zug).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
