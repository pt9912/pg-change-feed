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

- [x] Happy Path und Boundary am laufenden Container: für eine per
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
- [x] Neustart-Festigkeit und Ausschluss+Regel: nach einem **realen** `docker
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
- [x] Die E2E-Abdeckung ist getragen: jede neue `func TestE2E*` und jede neue
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
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: `harness/README.md` §Sensors — die Zeile `make
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
| `test/integration/integration_test.go` | update → **nicht realisiert** | geplant waren die `TestE2E…`-Funktionen in dieser Datei. Sie stehen in einer neuen Datei (nächste Zeile) nach dem Muster von `backfill_e2e_test.go`; die Datei trägt Hilfen und den Abdeckungs-Erzeuger und bleibt unverändert, der Erzeuger liest jede Go-Datei des Verzeichnisses (`abdeckungsZeilen`), die Zeilen der neuen Funktionen erscheinen in der Tabelle. |
| `test/integration/transformation_e2e_test.go` | neu | Go-Hälfte der Abdeckung, zwei Funktionen mit [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) im Doc-Kommentar: `TestE2ETransformationRulesShapeBothImages` (Happy Path — `rename_column` und `map_value` prägen Alt- und Neu-Bild von INSERT, UPDATE und DELETE an einer Tabelle mit voller Replica-Identität; Rohform vor den Regeln; nicht abgebildeter Wert, `NULL`; zwei Lesungen gleich) und `TestE2ETransformationConflictsFailWithSpecText` (Boundary — K1, K2, K3 in beiden Formen, K4 in beiden Formen, je `failed` mit dem Klartext und der Adresse der Spec; Regelstand danach unverändert; Gegenprobe `applied`). Der Neustart braucht `docker restart` und steht im Runner. Die Boundary steht in Go statt im Runner, weil sie nur SQL-Funktionen, `cdc.administration_request` und `cdc.changes` über pgx nutzt und keinen Container-Zugriff braucht, den Klartext von [`SPEC-019`](../../../../spec/pflichtenheft.md) wörtlich vergleicht und im ersten `go test`-Aufruf des Runners läuft; die Datei importiert keine interne Anwendungslogik (E2E-Tier, [`ADR-0030`](../../adr/0030-testpyramide.md)). Die Hilfen (`newBackfillEnv`, `enableTable`, `awaitRequestApplied`) stammen aus `backfill_e2e_test.go`. |
| `tools/harness/run-integration-tests.sh` | update | zwei Phasen nach der Leerlauf-Bestätigung und vor dem Upgrade-Tausch: Happy Path über alle fünf Zustellwege (mit `rename_column` und `map_value`, Rohform davor und danach) und Neustart mit Ausschluss (zwei Neustarts, jeder mit Startzeit-Vergleich über `tf_restart_feed`; Rücknahme der Regeln; die zweite Phase liest die Regelformen `TF_RULE_*` aus dem Block der ersten, im Kommentar benannt); die Regeln der Phasen werden am Ende zurückgenommen; die zwei neuen Testfunktionen stehen im `-run`-Muster des ersten `go test`-Aufrufs; die K1–K4-Boundary steht in der Go-Datei, nicht im Runner; Kopfkommentar nennt die Phasen; je Phase ein `abdeckung_declare`-Anker mit [`LH-FA-CFG-007`](../../../../spec/lastenheft.md). |
| `tools/harness/httpclient`, `grpcclient`, `sseclient`, `natsstreamsub` | geprüft, unverändert | alle vier geben das vollständige Row Image aus (`new_image=<JSON>`; gemessen `git grep -n 'new_image' -- tools/harness/httpclient tools/harness/grpcclient tools/harness/sseclient tools/harness/natsstreamsub`, Zeile im Block unten): kein fester Feldzugriff, kein Plan-Nachzug. Grenze der Aussage „alle Wege dieselbe Form“: die Clients der drei Stream-Wege und der HTTP-Client drucken das Neu-Bild; die Phase belegt die Form auf den fünf Wegen daher für eine eingefügte Zeile (INSERT). Die Bilder von UPDATE und DELETE, das Alt-Bild eingeschlossen, belegt allein `TestE2ETransformationRulesShapeBothImages` über `cdc.changes`; der Deklarations-Anker der Phase im Runner nennt beide Hälften. |
| `docs/user/e2e-abdeckung.md` | Erzeugnis | kommt aus dem Runner, wird nicht von Hand geschrieben. Die `Ort`-Angaben sind Zeilennummern des Quelltexts am Stand des schreibenden Laufs: jede Änderung von Go- oder Runner-Zeilen oberhalb einer Deklaration verschiebt sie, und nur ein Lauf des Runners schreibt sie neu. Die Zwischenstände vor dem letzten Lauf tragen abweichende Nummern (Kommentar-Kürzung in `66f60c8b`, damit die Nummern des Endstands stimmen); am Endstand ist die Datei die Ausgabe des Runners (`abdeckung_schreiben` schreibt sie nach dem letzten Phasen-Anker; ein weiterer Lauf ohne Quelländerung meldet `E2E-Abdeckungstabelle unverändert`). |
| `harness/README.md` §Sensors (`make test-integration`, `make doc-trace`) | update | Aufzählung der Belege (vier neue Belege in der Zeile `make test-integration`); Waisen-Messung nachgemessen (Zeile `make doc-trace`). |
| `compose.yaml` | geprüft, unverändert | die Tabellen werden über `cdc.enable_table` aktiviert, `CDC_TABLES` bleibt unberührt. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „der Inhalt von
`make test-integration`“, „die Waisen im RTM-Lauf“, „die Zeilen der
E2E-Abdeckungstabelle“; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Beschreibung der Belege von `make test-integration` | Zeilen 1–2 im Block unten (Symbolname einer bestehenden Phase, ganzer Baum ohne `docs/reviews`, `done/` und die Baseline) | **Gefunden.** `Spaltenausschluss-Rundlauf`: Parent 12, Diff 12, verteilt auf drei Dateien ohne diesen Plan: `tools/harness/run-integration-tests.sh` (Deklaration, Kommentare, Ausgabezeile der Phase), `harness/README.md` (die Zeile `make test-integration`, ein Aufzähler) und `docs/user/e2e-abdeckung.md` (Erzeugnis). **Nichtgefunden:** kein weiterer Träger unter `docs/user`, `spec` oder in einem Handbuch zählt die Belege von `make test-integration` auf (das Benutzerhandbuch nennt den Lauf nur als Beleg-Ort einzelner Backfill-Aussagen, Zeilen 579 und 612, und in einer Zeile der Änderungshistorie). | `harness/README.md` (vier Belege in der Zeile `make test-integration`), Kopfkommentar des Runners (die Transformations-Rundläufe); die Abdeckungstabelle ist Erzeugnis |
| Waisen-Aussage im RTM-Träger | Zeilen 3–6 im Block unten | **Gefunden.** Zählwort und Zahl „2 Waisen“: Parent 2, Diff 1 — die Zeile `make doc-trace` in `harness/README.md` (Parent: „**2 Waisen** — `LH-FA-CFG-007`/`LH-FA-CFG-008`“, Diff nachgezogen) und die Beleg-Datei `slice-sdk-csharp-reale2e` im Beobachtungs-Register (Record des 2026-09-23, unverändert). **Nachgemessen:** `make doc-trace` am Parent-Stand `38c7b3bc`: 80 Anforderungen, **1 Waise**, `LH-FA-CFG-008`; `LH-FA-CFG-007` trägt dort schon `E2E` über die Zeile der Phase „Backfill-Regelstand“ (Commit `fc0b8d38`). Am Diff-Stand: dieselbe Zahl, `LH-FA-CFG-007` trägt fünf Zeilen (`git grep -n 'LH-FA-CFG-007' -- docs/user/e2e-abdeckung.md`: Parent 1, Diff 5). **Nichtgefunden:** kein weiterer Träger mit einer Waisen-Zahl außer den Sätzen zum Ziel `LH-FA-CFG-007 verlässt die Waisen` in `docs/plan/planning/welle-transformationen.md` (Zeilen 53, 124–127, 168; Parent 2 Treffer des Musters, Diff 2). | Die Zeile `make doc-trace` trägt 80 Anforderungen und 1 Waise mit Datum und Ursprung. Die Sätze der Welle-Datei sind ein fremder Träger (Suche: `git grep -n -E 'verlässt die Waisen|nicht mehr (unter den )?Waise' -- docs/plan/planning/welle-transformationen.md`, zwei Treffer wie im Suchlauf-Block) und werden gemeldet, nicht mitgeändert — Adresse: Planner; Frist: die Closure dieses Slice, bis dahin zieht der Planner die Sätze nach oder benennt den Träger mit Adresse: das Ziel ist bereits vor diesem Slice erreicht, dieser Slice fügt vier Träger hinzu |
| Zeilenzahlen der Abdeckungstabelle | Zeilen 13–14 im Block unten (`LH-FA-CFG-007` in der Tabelle) und `git diff --stat 38c7b3bc -- tools/harness/run-integration-tests.sh docs/user/e2e-abdeckung.md test/integration harness/README.md` (gemessen am Diff-Stand nach dem Lauf von `make test-integration`) | **Gefunden.** Die Tabelle trägt 54 Zeilen (16 Go-Zeilen und 38 Bash-Zeilen statt 14 und 36, nachmessbar als `^func TestE2E`- und `^abdeckung_declare`-Zeilen: Zeilen 17–20 im Block unten; der Runner druckt sie beim Schreiben als „E2E-Abdeckungstabelle aus 16 Go-Zeilen und 38 Bash-Zeilen“) statt 50; Stat: `docs/user/e2e-abdeckung.md` 40 Einfügungen und 36 Löschungen, `harness/README.md` 2/2, `test/integration/transformation_e2e_test.go` 324 neu, `tools/harness/run-integration-tests.sh` 301 Einfügungen und 2 Löschungen (Kopfkommentar, `-run`-Muster und der Block der zwei Phasen; gemessen mit `git diff --numstat 38c7b3bc` am Stand `4800d75b`, Summe der vier Dateien 667 und 40, abgeleitet). **Nichtgefunden:** keine Zeilennummer der Tabelle ist von Hand geschrieben; jede Verschiebung der Runner-Zeilen ist im Erzeugnis vom Runner nachgezogen. | die Tabelle ist Erzeugnis; sie wurde vom Runner neu geschrieben („E2E-Abdeckungstabelle geschrieben“) |
| Jede neue `TestE2E*`-Funktion ist von einem `-run`-Muster erfasst | Zeilen 7–12 im Block unten | **Gefunden.** `func TestE2E` unter `test/integration`: Parent 19, Diff 21 (die zwei neuen Funktionen); beide Namen stehen im `-run`-Muster des ersten `go test`-Aufrufs (Parent 0 Treffer je Name, Diff 1 für `TestE2ETransformationConflictsFailWithSpecText`, Diff 2 für `TestE2ETransformationRulesShapeBothImages`: Muster und Deklarations-Anker der Phase Happy Path, der die Grenze der Zustellwege nennt). Dass beide laufen und bestehen (`--- PASS`), steht in keinem committeten Träger: der Beleg ist ein Lauf von `make test-integration` und **übernommen** aus dem Bericht des Implementers; die Wiederholung trägt die Verifikation. **Nichtgefunden:** keine `func TestE2E*` ohne Muster: alle 16 Go-Zeilen der Tabelle sind von den vier `-run`-Aufrufen des Runners erfasst (die zwei neuen im ersten Aufruf; Abgleich: die Namen aus `git grep -hoE '^func TestE2E[A-Za-z0-9_]+' -- test/integration` gegen die Zeilen von `git grep -n -F -e "-run '" -- tools/harness/run-integration-tests.sh`, fünf Treffer: die vier E2E-Aufrufe und der Aufruf des Erzeugers `TestAbdeckungstabelleZeilen`). | Abgleich mit den zwei Befehlen der Zelle |
| Die vier Wegwerf-Clients (Plan-Zeile „prüfen“) | Zeilen 15–16 im Block unten | **Gefunden.** `new_image` in den vier Client-Verzeichnissen: Parent 7, Diff 7 — jeder Client gibt das vollständige Row Image aus. **Nichtgefunden:** kein fester Feldzugriff auf einzelne Schlüssel in den Clients (sie reichen `json.RawMessage` bzw. `[]byte` durch). | keiner, unverändert |

```suchlauf
38c7b3bc 12 -n -F 'Spaltenausschluss-Rundlauf' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 12 -n -F 'Spaltenausschluss-Rundlauf' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
38c7b3bc 2 -n -F '2 Waisen' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 1 -n -F '2 Waisen' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
38c7b3bc 2 -n -E 'verlässt die Waisen|nicht mehr (unter den )?Waise' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 2 -n -E 'verlässt die Waisen|nicht mehr (unter den )?Waise' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
38c7b3bc 19 -n -F 'func TestE2E' -- test/integration
diff 21 -n -F 'func TestE2E' -- test/integration
38c7b3bc 0 -n -F 'TestE2ETransformationRulesShapeBothImages' -- tools/harness/run-integration-tests.sh
diff 2 -n -F 'TestE2ETransformationRulesShapeBothImages' -- tools/harness/run-integration-tests.sh
38c7b3bc 0 -n -F 'TestE2ETransformationConflictsFailWithSpecText' -- tools/harness/run-integration-tests.sh
diff 1 -n -F 'TestE2ETransformationConflictsFailWithSpecText' -- tools/harness/run-integration-tests.sh
38c7b3bc 1 -n -F 'LH-FA-CFG-007' -- docs/user/e2e-abdeckung.md
diff 5 -n -F 'LH-FA-CFG-007' -- docs/user/e2e-abdeckung.md
38c7b3bc 7 -n -F 'new_image' -- tools/harness/httpclient tools/harness/grpcclient tools/harness/sseclient tools/harness/natsstreamsub
diff 7 -n -F 'new_image' -- tools/harness/httpclient tools/harness/grpcclient tools/harness/sseclient tools/harness/natsstreamsub
38c7b3bc 14 -n -E '^func TestE2E' -- test/integration
diff 16 -n -E '^func TestE2E' -- test/integration
38c7b3bc 36 -n -E '^abdeckung_declare ' -- tools/harness/run-integration-tests.sh
diff 38 -n -E '^abdeckung_declare ' -- tools/harness/run-integration-tests.sh
```

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
- **Fenster READY → erste Zeile trägt die Timeouts der Stream-Clients.** Die
  drei Clients starten gleichzeitig, die Phase fügt die erste Zeile erst ein,
  nachdem alle drei `READY` gemeldet haben. Ihre Fristen laufen ab dem eigenen
  `READY` (NATS 30 s in `NextMsg`, gRPC und SSE je 60 s ab der Stream-Öffnung);
  liegt das `READY` eines Clients um mehr als seine Frist vor dem des
  langsamsten, endet er mit `TIMEOUT`, und die Phase endet laut mit der
  Meldung, welcher Client nichts empfing (kein stilles Grün). Gemessen wurde
  das Fenster am Host des Reviews (20 Kerne: `go build` kalt 4,85 s bis
  6,15 s je Client, der erste Einfüge-Versuch genügte;
  `docs/reviews/review-slice-transformationen-e2e-wirkung.md` F-4); auf dem
  Runner von `e2e.yml` ist es nicht gemessen. Eine weitere Absicherung
  im Runner gehört nicht zu diesem Slice: die Einfügung folgt bereits der
  `READY`-Bestätigung aller Clients, und die Fristen liegen in den
  Wegwerf-Clients, die dieser Slice nicht ändert. *Erwartet, zu
  belegen durch:* der erste reale Lauf von `e2e.yml` nach dem Push
  ([`AGENTS.md`](../../../../AGENTS.md) §3.10, gehört zur Closure; der
  Workflow selbst ist unverändert). **Ausgang:** *(bei Closure)*
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
