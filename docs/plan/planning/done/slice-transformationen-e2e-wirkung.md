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
Die Belege stehen in `make test-integration` und fügen dem RTM-Träger, der vom
Runner geschriebenen Tabelle
[`docs/user/e2e-abdeckung.md`](../../../user/e2e-abdeckung.md), vier Zeilen für
[`LH-FA-CFG-007`](../../../../spec/lastenheft.md) hinzu. Die Anforderung ist im
RTM-Lauf (`make doc-trace`) schon vor diesem Slice gedeckt: die Phase
„Backfill-Regelstand“ (`fc0b8d38`) trägt die erste Zeile, gemessen am Parent-Stand
`38c7b3bc` (§3, Zeile „Waisen-Aussage im RTM-Träger“).

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
      Change bleibt in Rohform lesbar; jede der sechs Zeilen K1 bis K4 der
      Fehlertexte ([`SPEC-019`](../../../../spec/pflichtenheft.md): K1 und K2 je
      eine Zeile, K3 und K4 je zwei) endet real `failed` mit dem Text der
      Spec, und die danach erfasste Change trägt weiter die gültige Form. *Zu belegen durch:* ein realer, grüner `make
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
      sichtbar); das Kriterium prüft den Stand des Diffs und galt am
      Parent-Stand `38c7b3bc` bereits (§3, gemessen). *Zu belegen durch:* `make
      doc-trace` (Ausgabe im Bericht) und `make docs-check`.
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
      `betriebsdoku` unberührt (§6 letzte Zeile, Welle §4 Abweichung 3).
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
      der Closure der Welle
      [welle-transformationen](../welle-transformationen.md) (die Roadmap führt
      sie unter *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `test/integration/integration_test.go` | unverändert (die Änderungs-Art „update“ des Plans ist nicht realisiert) | die `TestE2E…`-Funktionen stehen in einer neuen Datei (nächste Zeile) nach dem Muster von `backfill_e2e_test.go`; diese Datei trägt Hilfen und den Abdeckungs-Erzeuger, der Erzeuger liest jede Go-Datei des Verzeichnisses (`abdeckungsZeilen`), die Zeilen der neuen Funktionen erscheinen in der Tabelle. |
| `test/integration/transformation_e2e_test.go` | neu | Go-Hälfte der Abdeckung, zwei Funktionen mit [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) im Doc-Kommentar: `TestE2ETransformationRulesShapeBothImages` (Happy Path — `rename_column` und `map_value` prägen Alt- und Neu-Bild von INSERT, UPDATE und DELETE an einer Tabelle mit voller Replica-Identität; Rohform vor den Regeln; nicht abgebildeter Wert, `NULL`; zwei Lesungen gleich) und `TestE2ETransformationConflictsFailWithSpecText` (Boundary — die sechs Zeilen K1, K2, K3 (zwei) und K4 (zwei) von [`SPEC-019`](../../../../spec/pflichtenheft.md) (Zeilen 753 bis 758), je `failed` mit dem Klartext und der Adresse der Spec; Regelstand danach unverändert; Gegenprobe `applied`). Der Neustart braucht `docker restart` und steht im Runner. Die Boundary steht in Go statt im Runner, weil sie nur SQL-Funktionen, `cdc.administration_request` und `cdc.changes` über pgx nutzt und keinen Container-Zugriff braucht, den Klartext von [`SPEC-019`](../../../../spec/pflichtenheft.md) wörtlich vergleicht und im ersten `go test`-Aufruf des Runners läuft; die Datei importiert keine interne Anwendungslogik (E2E-Tier, [`ADR-0030`](../../adr/0030-testpyramide.md)). Die Hilfen (`newBackfillEnv`, `enableTable`, `awaitRequestApplied`) stammen aus `backfill_e2e_test.go`. |
| `tools/harness/run-integration-tests.sh` | update | zwei Phasen nach der Leerlauf-Bestätigung und vor dem Upgrade-Tausch: Happy Path über alle fünf Zustellwege (mit `rename_column` und `map_value`, Rohform davor und danach) und Neustart mit Ausschluss (zwei Neustarts, jeder mit Startzeit-Vergleich über `tf_restart_feed`; Rücknahme der Regeln; die zweite Phase liest die Regelformen `TF_RULE_*` aus dem Block der ersten, im Kommentar benannt); die Regeln der Phasen werden am Ende zurückgenommen; die zwei neuen Testfunktionen stehen im `-run`-Muster des ersten `go test`-Aufrufs; die K1–K4-Boundary steht in der Go-Datei, nicht im Runner; Kopfkommentar nennt die Phasen; je Phase ein `abdeckung_declare`-Anker mit [`LH-FA-CFG-007`](../../../../spec/lastenheft.md). |
| `tools/harness/httpclient`, `grpcclient`, `sseclient`, `natsstreamsub` | geprüft, unverändert | alle vier geben das vollständige Row Image aus (`new_image=<JSON>`; gemessen `git grep -n 'new_image' -- tools/harness/httpclient tools/harness/grpcclient tools/harness/sseclient tools/harness/natsstreamsub`, Zeile im Block unten): kein fester Feldzugriff, kein Plan-Nachzug. Grenze der Aussage „alle Wege dieselbe Form“: die Clients der drei Stream-Wege und der HTTP-Client drucken das Neu-Bild; die Phase belegt die Form auf den fünf Wegen daher für eine eingefügte Zeile (INSERT). Die Bilder von UPDATE und DELETE, das Alt-Bild eingeschlossen, belegt allein `TestE2ETransformationRulesShapeBothImages` über `cdc.changes`; der Deklarations-Anker der Phase im Runner nennt beide Hälften. **Restfläche mit Adresse:** UPDATE, DELETE und das Alt-Bild auf den vier weiteren Wegen sind aus der Architektur hergeleitet (alle Wege lesen dieselbe Change, [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 6 Option D), nicht erprobt. Kein offener Slice trägt sie: `slice-sdk-regel-realserver-e2e` führt `old_image` einer UPDATE-/DELETE-Change ausdrücklich als nicht Gegenstand (§6 dort; gelesen bei der Closure), und die Pläne von `e2e-abhilfe` und `betriebsdoku` nennen sie nicht (Suche bei der Closure). Adresse: das Closure-Kriterium „Restfläche der Zustellwege“ in [welle-transformationen](../welle-transformationen.md) §3, ein Ereignis, das eintreten kann (die Roadmap führt die Welle unter *Offene Wellen*). |
| `docs/user/e2e-abdeckung.md` | Erzeugnis | kommt aus dem Runner, wird nicht von Hand geschrieben. Die `Ort`-Angaben sind Zeilennummern des Quelltexts am Stand des schreibenden Laufs: jede Änderung von Go- oder Runner-Zeilen oberhalb einer Deklaration verschiebt sie, und nur ein Lauf des Runners schreibt sie neu. Die Zwischenstände vor dem letzten Lauf tragen abweichende Nummern; am Endstand ist die Datei die Ausgabe des Runners (`abdeckung_schreiben` schreibt sie nach dem letzten Phasen-Anker; ein weiterer Lauf ohne Quelländerung meldet `E2E-Abdeckungstabelle unverändert`, gedruckt in der Verifikation und in beiden Legs des Laufs von `e2e.yml` nach dem Push, §7). |
| `harness/README.md` §Sensors (`make test-integration`, `make doc-trace`) | update | Aufzählung der Belege (vier neue Belege in der Zeile `make test-integration`); Waisen-Messung nachgemessen (Zeile `make doc-trace`). |
| `compose.yaml` | geprüft, unverändert | die Tabellen werden über `cdc.enable_table` aktiviert, `CDC_TABLES` bleibt unberührt. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „der Inhalt von
`make test-integration`“, „die Waisen im RTM-Lauf“, „die Zeilen der
E2E-Abdeckungstabelle“; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Beschreibung der Belege von `make test-integration` | Zeilen 1–2 im Block unten (Symbolname einer bestehenden Phase, ganzer Baum ohne `docs/reviews`, `done/` und die Baseline) | **Gefunden.** `Spaltenausschluss-Rundlauf`: Parent 12, Diff 12, verteilt auf drei Dateien ohne diesen Plan: `tools/harness/run-integration-tests.sh` (Deklaration, Kommentare, Ausgabezeile der Phase), `harness/README.md` (die Zeile `make test-integration`, ein Aufzähler) und `docs/user/e2e-abdeckung.md` (Erzeugnis). **Nichtgefunden:** kein weiterer Träger unter `docs/user`, `spec` oder in einem Handbuch zählt die Belege von `make test-integration` auf (das Benutzerhandbuch nennt den Lauf nur als Beleg-Ort einzelner Backfill-Aussagen, Zeilen 579 und 612, und in einer Zeile der Änderungshistorie). | `harness/README.md` (vier Belege in der Zeile `make test-integration`), Kopfkommentar des Runners (die Transformations-Rundläufe); die Abdeckungstabelle ist Erzeugnis |
| Waisen-Aussage im RTM-Träger | Zeilen 3–6 im Block unten | **Gefunden.** Zählwort und Zahl „2 Waisen“: Parent 2, Diff 1 — die Zeile `make doc-trace` in `harness/README.md` (Parent: „**2 Waisen** — `LH-FA-CFG-007`/`LH-FA-CFG-008`“, Diff nachgezogen) und die Beleg-Datei `slice-sdk-csharp-reale2e` im Beobachtungs-Register (Record des 2026-09-23, unverändert). **Nachgemessen:** `make doc-trace` am Parent-Stand `38c7b3bc`: 80 Anforderungen, **1 Waise**, `LH-FA-CFG-008`; `LH-FA-CFG-007` trägt dort schon `E2E` über die Zeile der Phase „Backfill-Regelstand“ (Commit `fc0b8d38`). Am Diff-Stand: dieselbe Zahl, `LH-FA-CFG-007` trägt fünf Zeilen (`git grep -n 'LH-FA-CFG-007' -- docs/user/e2e-abdeckung.md`: Parent 1, Diff 5). **Nichtgefunden:** kein weiterer Träger mit einer Waisen-Zahl außer den Sätzen zum Ziel `LH-FA-CFG-007 verlässt die Waisen` in `docs/plan/planning/welle-transformationen.md` (Zeilen 53, 124–127, 168; Parent 2 Treffer des Musters; nach dem Nachzug der Closure Diff 1, die Zeile 53 beschreibt den Zustand des Bündels und ist wahr). | Die Zeile `make doc-trace` trägt 80 Anforderungen und 1 Waise mit Datum und Ursprung. Die Welle-Datei war ein fremder Träger und ist in der Closure dieses Slice nachgezogen (Frist der Meldung eingehalten): die Tabellenzeile des Slice nennt vier weitere E2E-Träger für die Anforderung, und der Closure-Trigger der Welle nennt den Beginn der Deckung (Phase „Backfill-Regelstand“) und die Messung an beiden Ständen |
| Zeilenzahlen der Abdeckungstabelle | Zeilen 13–14 im Block unten (`LH-FA-CFG-007` in der Tabelle) und `git diff --stat 38c7b3bc -- tools/harness/run-integration-tests.sh docs/user/e2e-abdeckung.md test/integration harness/README.md` (gemessen am Diff-Stand nach dem Lauf von `make test-integration`) | **Gefunden.** Die Tabelle trägt 54 Zeilen (16 Go-Zeilen und 38 Bash-Zeilen statt 14 und 36, nachmessbar als `^func TestE2E`- und `^abdeckung_declare`-Zeilen: Zeilen 17–20 im Block unten; der Runner druckt sie beim Schreiben als „E2E-Abdeckungstabelle aus 16 Go-Zeilen und 38 Bash-Zeilen“) statt 50; Stat: `docs/user/e2e-abdeckung.md` 40 Einfügungen und 36 Löschungen, `harness/README.md` 2/2, `test/integration/transformation_e2e_test.go` 324 neu, `tools/harness/run-integration-tests.sh` 301 Einfügungen und 2 Löschungen (Kopfkommentar, `-run`-Muster und der Block der zwei Phasen; gemessen mit `git diff --numstat 38c7b3bc` am Stand `4800d75b`, Summe der vier Dateien 667 und 40, abgeleitet). **Nichtgefunden:** keine Zeilennummer der Tabelle ist von Hand geschrieben; jede Verschiebung der Runner-Zeilen ist im Erzeugnis vom Runner nachgezogen. | die Tabelle ist Erzeugnis; sie wurde vom Runner neu geschrieben („E2E-Abdeckungstabelle geschrieben“) |
| Jede neue `TestE2E*`-Funktion ist von einem `-run`-Muster erfasst | Zeilen 7–12 im Block unten | **Gefunden.** `func TestE2E` unter `test/integration`: Parent 19, Diff 21 (die zwei neuen Funktionen); beide Namen stehen im `-run`-Muster des ersten `go test`-Aufrufs (Parent 0 Treffer je Name, Diff 1 für `TestE2ETransformationConflictsFailWithSpecText`, Diff 2 für `TestE2ETransformationRulesShapeBothImages`: Muster und Deklarations-Anker der Phase Happy Path, der die Grenze der Zustellwege nennt). Dass beide laufen und bestehen (`--- PASS`), steht in der Verifikation (Wiederholung von `make test-integration`, Exit 0) und in den Job-Logs des Laufs von `e2e.yml` nach dem Push (§7); der Bericht des Implementers (**übernommen**) ist kein Träger. **Nichtgefunden:** keine `func TestE2E*` ohne Muster: alle 16 Go-Zeilen der Tabelle sind von den vier `-run`-Aufrufen des Runners erfasst (die zwei neuen im ersten Aufruf; Abgleich: die Namen aus `git grep -hoE '^func TestE2E[A-Za-z0-9_]+' -- test/integration` gegen die Zeilen von `git grep -n -F -e "-run '" -- tools/harness/run-integration-tests.sh`, fünf Treffer: die vier E2E-Aufrufe und der Aufruf des Erzeugers `TestAbdeckungstabelleZeilen`). | Abgleich mit den zwei Befehlen der Zelle |
| Die vier Wegwerf-Clients (Plan-Zeile „prüfen“) | Zeilen 15–16 im Block unten | **Gefunden.** `new_image` in den vier Client-Verzeichnissen: Parent 7, Diff 7 — jeder Client gibt das vollständige Row Image aus. **Nichtgefunden:** kein fester Feldzugriff auf einzelne Schlüssel in den Clients (sie reichen `json.RawMessage` bzw. `[]byte` durch). | keiner, unverändert |
| Träger mit dem Lifecycle-Ort dieses Slice (Closure: der Move nach `done/`) | Zeilen 21–24 im Block unten (Parent der Closure `e5d6b6bf`, ganzer Baum ohne `docs/reviews`, `done/` und die Baseline; Zeilen 21–22 das Verzeichnis `open/`, Zeilen 23–24 `next/` und `in-progress/`) | **Gefunden.** `open/slice-transformationen-e2e-wirkung`: Parent 1 — `docs/plan/planning/open/slice-sdk-regel-realserver-e2e.md`, Träger-Tabelle, Inline-Code-Pfad mit dem Verzeichnis `open/`; Diff 0 nach dem Nachzug (die Kennung ohne Verzeichnis). **Nichtgefunden:** kein Träger nennt `next/` oder `in-progress/` als Ort dieses Slice (Zeilen 23–24: Parent 0, Diff 0); die Start-Trigger der drei Folge-Slices (`start-reihenfolge`, `e2e-abhilfe`, `sdk-regel-realserver-e2e`) nennen „liegt in `done/`“ als Bedingung und bleiben mit dem Move erfüllt. | Kennung statt Pfad (Fremddatei, in der Closure gezogen) |

```suchlauf
38c7b3bc 12 -n -F 'Spaltenausschluss-Rundlauf' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 12 -n -F 'Spaltenausschluss-Rundlauf' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
38c7b3bc 2 -n -F '2 Waisen' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 1 -n -F '2 Waisen' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
38c7b3bc 2 -n -E 'verlässt die Waisen|nicht mehr (unter den )?Waise' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 1 -n -E 'verlässt die Waisen|nicht mehr (unter den )?Waise' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
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
e5d6b6bf 1 -n -F 'open/slice-transformationen-e2e-wirkung' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 0 -n -F 'open/slice-transformationen-e2e-wirkung' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
e5d6b6bf 0 -n -E '(next|in-progress)/slice-transformationen-e2e-wirkung' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 0 -n -E '(next|in-progress)/slice-transformationen-e2e-wirkung' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
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
  alle fünf Wege. **Ausgang: entfallen** für das Neu-Bild einer eingefügten
  Zeile — gRPC, SSE und NATS empfangen dieselbe Zeile (gedruckt:
  `new_image={"id":"11","customer_name":"TfNeu","status":"open"}`, gleiche
  `change_id`), `cdc.changes` und `GET /changes` tragen dieselbe Form, gemessen in
  der Verifikation (Lauf am HEAD `c3348898`) und im Lauf von `e2e.yml` nach dem
  Push (Lauf 36268704311, beide PostgreSQL-Legs; die drei Zeilen gelesen im
  Job-Log beider Legs). **Restfläche (benannte Grenze, kein Risiko-Ausgang):**
  UPDATE, DELETE und das Alt-Bild auf den Stream-Wegen sind hergeleitet, nicht
  erprobt; Adresse in §3 (Zeile der Wegwerf-Clients): das Closure-Kriterium in
  [welle-transformationen](../welle-transformationen.md) §3.
- **Laufzeit von `make test-integration`** wächst
  (`BEO-PGC/test-integration-retention-timing-flake`, verkörpert, 3×); die
  Pipeline `e2e.yml` fährt das Testpaket unverändert (keine Workflow-Änderung,
  [`AGENTS.md`](../../../../AGENTS.md) §3.10 greift nicht). *Erwartet, zu
  belegen durch:* ein realer Lauf, die Zahl trägt ihren Lauf. **Ausgang:
  entfallen** — die Laufzeit wächst um die zwei Phasen und bleibt weit unter der
  Grenze `timeout-minutes: 60` von `e2e.yml`. Gemessen (Befehl:
  `gh api repos/pt9912/pg-change-feed/actions/runs/<Lauf>/jobs`, Schritt
  „Compose-Integrationstest (Black-Box-E2E)“, `started_at` bis `completed_at`):
  im Lauf 36268704311 nach dem Push 579 s (PostgreSQL 18, 20:13:05 bis 20:22:44)
  und 491 s (PostgreSQL 17, 20:13:00 bis 20:21:11); im Lauf 36263188406 am
  Parent-Stand `c246ba4f` 501 s (18) und 493 s (17); Unterschied +78 s und −2 s
  (abgeleitet), die Streuung des Runners übersteigt den Beitrag der zwei Phasen.
  Die Phasen selbst: 70,3 s (18) und 56,3 s (17) von der Meldung der
  Leerlauf-Bestätigung bis zur Meldung der Neustart-Phase (abgeleitet aus den
  Zeitstempeln des Job-Logs). Lokal: 346 s, gemessen in der Verifikation
  (**übernommen**, Host mit 20 Kernen).
- **Eine neue Testfunktion fällt still aus dem Runner**
  (`BEO-PGC/test-runner-stiller-ausschluss`, offen, 2×). *Erwartet, zu belegen
  durch:* der Abgleich aus DoD Punkt 3. **Ausgang: entfallen** — 16 Funktionen
  `func TestE2E*`, jede im `-run`-Muster eines der vier E2E-Aufrufe des Runners
  (Abgleich bei der Closure gefahren: 16 Namen, 0 ohne Muster; die beiden neuen
  laufen im Log beider Legs von `e2e.yml` mit `--- PASS`).
- **Geteilter Zustand zwischen Phasen**
  (`BEO-PGC/test-isolation-geteilter-zustand`, offen, 1×;
  `BEO-PGC/e2e-metrik-boundary-nur-reviewer-belegt`, offen, 1×): Regeln einer
  Phase wirken auf eine spätere Phase derselben Tabelle. *Erwartet, zu belegen
  durch:* eine eigene Tabelle je Phase und die Rücknahme der Regel im
  Aufräumen. **Ausgang: entfallen** — jede Phase führt eine eigene Tabelle
  (`feed_e2e_transform`, `feed_e2e_transform_restart`, die Go-Tests eigene
  Tabellennamen), beide Phasen nehmen ihre Regeln am Ende zurück; die
  nachfolgenden Phasen (Upgrade-Rundlauf, Schema-Tests) laufen in beiden Legs
  von `e2e.yml` grün. Die Kopplung der zweiten Phase an die Regelformen
  `TF_RULE_*` der ersten steht im Kommentar des Runners (Review F-8).
- **Neustart-Beleg misst Timing statt Zustand.** Der Poll auf `applied` und die
  Wartezeit nach `docker restart` sind Zeitannahmen. *Erwartet, zu belegen
  durch:* Poll auf Zustand mit Frist, nicht feste Wartezeit. **Ausgang:
  entfallen** — der Beleg des Neustarts ist die Startzeit des Containers vor und
  nach jedem `docker restart` (`tf_restart_feed`), die Wartezeit ein Poll auf
  den Zustand mit Frist; beide Neustart-Aufrufe sind einzeln mutiert und rot
  (Verifikation R1, R2), auf dem Runner stehen beide Startzeit-Paare im
  Job-Log (PostgreSQL 18: 20:20:39,92 → 20:21:42,70 und 20:21:42,70 →
  20:21:49,51; PostgreSQL 17: 20:19:29,35 → 20:20:18,37 und 20:20:18,37 →
  20:20:24,95).
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
  Workflow selbst ist unverändert). **Ausgang: entfallen** — auf dem Runner
  von `e2e.yml` (Lauf 36268704311, `headSha` `e5d6b6bf`, beide Legs
  `success`) endet die Phase Happy Path in beiden Legs mit „1 Stream-Zeile(n)“:
  der erste Einfüge-Versuch genügte (Zeile mit `id` 11), alle drei Clients
  melden `RECEIVED` mit derselben `change_id`, das Job-Log enthält keine Zeile
  `TIMEOUT` (`grep -c TIMEOUT` je Job-Log: 0; beide Job-Logs über `gh api
  repos/pt9912/pg-change-feed/actions/jobs/<Job>/logs` gelesen, Ausschnitte,
  keine Geheimnisse). Die Messung gilt für zwei Läufe an einem Commit, nicht
  für die Streuung des Runners.
- **Das Handbuch entsteht erst danach** (Aufschub mit Adresse
  `slice-transformationen-betriebsdoku`, dessen §2 den aufgeschobenen
  Gegenstand nennt). **Ausgang: eingetreten** als geplanter Aufschub;
  Folge-Slice mit Kennung: `slice-transformationen-betriebsdoku` (`open/`),
  dessen §2 die zwei SQL-Funktionen, die Wirkung, die Abhilfe und die Zeile der
  Fehlerklasse als Gegenstand des Handbuch-Zugs nennt (gelesen bei der
  Closure); dieser Slice berührt `docs/user/benutzerhandbuch.md` nicht (Diff
  `c246ba4f..e5d6b6bf`: nicht enthalten).

## 7. Closure-Notiz

- **Was hat funktioniert:** (1) Die Belege binden an ihre Eingabeseite: die
  Verifikation fuhr 23 Mutationen und sah alle rot (12 im Runner, 9 in den
  Go-Tests, 2 im Produktivcode der Verdrahtung des Prozessstarts; **übernommen**
  aus dem Report, §4), der Review davor 15 mit einer grünen (R10, der zweite
  Neustart; **übernommen**, führte zu F-2). (2) Die Leser-Kette lief mit einer
  Fixrunde ohne Rückführung: der Review nennt 1 HIGH, 1 MEDIUM, 4 LOW und 4
  INFO, die Verifikation 0 HIGH, 0 MEDIUM, 2 LOW und 4 INFO (beide
  **übernommen** aus den Summary-Tabellen), DoD Punkt 1 bis 7 bestätigt, F-1 bis
  F-10 durch die Fixrunde (`4800d75b`, `c3348898`) bearbeitet. (3) Der Beleg von
  [`AGENTS.md`](../../../../AGENTS.md) §3.10 liegt vor: der Lauf von `e2e.yml`
  nach dem Push (36268704311, `headSha` `e5d6b6bf`, beide PostgreSQL-Legs
  `success`, gemessen mit `gh run view 36268704311 --json conclusion,jobs`), dazu
  `ci` (36268704449) und `examples` (36268704347), beide `success` (gemessen mit
  `gh run view <Lauf> --json conclusion`); Laufzeiten und Fenster in §6. (4) Die
  Closure maß nach, was Report und Plan nur nannten: `make doc-trace` am
  Parent-Stand `38c7b3bc` und am Arbeitsbaum („80 Anforderung(en), 1 Waise(n).“,
  je `LH-FA-CFG-007` mit Deckung `E2E`, gedruckt), den Abgleich der 16 Funktionen
  `func TestE2E*` gegen die `-run`-Muster (0 ohne Muster), den Kandidatenlauf von
  Schritt 20 (0 Zeilen am Stand des Implementers, §7 Register) und die Job-Logs
  beider Legs (die zwei Phasen-Meldungen, keine Zeile `TIMEOUT`).
- **Was ging anders als geplant:** (1) **Das Ziel „`LH-FA-CFG-007` verlässt die
  Waisen“ war vor dem Slice erreicht.** Die Phase „Backfill-Regelstand“ des
  Vorgänger-Slice (`fc0b8d38`) trägt die erste Zeile; am Parent-Stand `38c7b3bc`
  führt `make doc-trace` die Anforderung mit Deckung `E2E` (gemessen). Der Slice
  fügt vier Träger hinzu; die Zahl in `harness/README.md` (zwei Waisen) war seit
  `fc0b8d38` veraltet und ist mit diesem Slice richtiggestellt; Plan §1 und §2,
  die Slice-Tabelle und der Closure-Trigger der Welle stehen im Ist-Ton. (2)
  **Die Boundary K1 bis K4 steht in Go, die Happy-Path-Phasen im Runner**, in
  einer neuen Datei statt in `integration_test.go` (§3, mit Grund). (3) **Ein
  Kommentar nannte den Zustand ohne die Zusage im Konjunktiv** (Review F-1,
  HIGH): der Kandidatenlauf von Schritt 20 druckte am Stand des Implementers 0
  Zeilen, weil „trüge“ nicht in seiner Wortliste steht (gemessen, Register). (4)
  **Die Zusage „nach dem zweiten Neustart“ war an den zweiten Neustart nicht
  gebunden** (F-2, MEDIUM); die Fixrunde band beide Neustarts durch den
  Startzeit-Vergleich, die Verifikation sah beide Aufrufe einzeln rot. (5) **Ein
  Doc-Kommentar sagte „jeder Negativfall in genau einem Feld“** bei einem Fall
  anderer Antragsart (F-3). (6) **Umfang:** 11 Commits, 7 Dateien, +1344/−58
  einschließlich der drei Lifecycle-Moves und der zwei Berichte; ohne
  `docs/reviews` 5 Dateien mit 727 Einfügungen und 58 Löschungen (gemessen mit
  `git diff --shortstat c246ba4f e5d6b6bf` in der Closure); kein Produktivcode im
  Diff. (7) **Die Coverage streut:** 85,50 % im Lauf der Verifikation
  (**übernommen**), 85,30 % im Bericht des Implementers (**übernommen**), 85.40 %
  im `make gates` der Closure vor dem Inhalts-Commit (gedruckt: „coverage-gate:
  OK — Coverage 85.40% erfüllt Schwelle 80%“); Schwelle 80 % — der Beleg je eines
  Laufs, kein Ist-Stand. (8) **Nicht gefahren:** `make test`, `make test-store`,
  `make test-replication`, `make bench` (kein Produktivcode im Diff; Aussage der
  Verifikation, **übernommen**).
- **Steering-Loop-Eintrag (Lerneintrag):** *(a) Neuer Sensor:* keiner gebaut. Der
  Kandidatenlauf von `implement-slice` Schritt 20 hat eine **gemessene Lücke** (0
  Zeilen bei einem Konjunktiv-Kommentar, weil die Wortliste endlich ist); kein
  Sensor hier gebaut: eine Erweiterung der Wortliste ändert Schritt 20 von
  `.claude/commands/implement-slice.md` (ein Architect-Zug, nicht Sache dieser
  Closure), und eine endliche Liste von Konjunktiv-Formen bleibt unvollständig
  (hergeleitet, nicht erprobt; Schritt 20 selbst führt die Unterscheidung als
  Satz-Subjekt-Urteil). Adresse der Lücke: der Lese-Schritt der Closure von
  [welle-transformationen](../welle-transformationen.md) (der Eintrag
  `BEO-PGC/vorher-nachher-sprache-in-test-harness-kommentar` steht bei 7×).
  *(b) Geschärfte Regel:* keine neue Regel; zwei **Anwendungen** der
  verkörperten Regeln, die vor dem Merge gewirkt haben: die Zusage über den
  Zustand **nach einer realen Aktion** bindet die Aktion selbst (Startzeit
  davor ≠ danach je Neustart; Eingabeseite eines `docker restart`, Review F-2,
  Verifikation R1 und R2) — Anwendung von Schritt 19; und die **Zielaussage einer
  Welle-Tabelle** wird beim Schnitt des Slice-Plans am Ist-Stand gemessen
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B, der Planner ist Verfasser):
  die Zuschreibung „verlässt die Waisen“ wurde beim Schnitt der Welle nicht
  gemessen und war seit `fc0b8d38` falsch; den Fund machte der Suchlauf des
  Implementers (§3.13 wirkt durch Ausführen, nicht durch einen Sensor).
  *(c) Benannte Lücke:* die **Restfläche der Zustellwege** — die Aussage „alle
  Wege dieselbe Form“ ist für das Neu-Bild einer eingefügten Zeile auf den fünf
  Wegen erprobt (Menge: ein INSERT, `rename_column` und `map_value`, fünf Wege,
  zwei PostgreSQL-Legs), für UPDATE, DELETE und das Alt-Bild auf den vier Wegen
  neben `cdc.changes` **hergeleitet** (alle Wege lesen dieselbe Change,
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 6); Adresse: das Closure-Kriterium „Restfläche der Zustellwege“ in
  [welle-transformationen](../welle-transformationen.md) §3 (Folge-Slice mit
  Kennung oder als hergeleitet geführte Aussage mit der gemessenen Menge).
- **Beobachtungs-Register (`../observations/`):** je Anfall eine Datei
  `evidence/slice-transformationen-e2e-wirkung.md`, Zähler = Zahl der Dateien
  (gemessen mit `ls evidence | wc -l` am Stand dieser Closure). *Neue Belege:*
  `BEO-PGC/vorher-nachher-sprache-in-test-harness-kommentar` **7×** (F-1 HIGH;
  verkörpert; dazu die benannte Lücke der Wortliste);
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` **19×** (F-2 MEDIUM, daher
  Datei trotz Deckel; verkörpert);
  `BEO-PGC/test-name-behauptet-mehr-als-der-test-treibt` **3×** (F-3 LOW; offen,
  **Schwelle erreicht**); `BEO-PGC/gemeldete-ungenauigkeit-ohne-traeger` **3×**
  (V-2 LOW, F-7 INFO; die Form „benannte Grenze ohne Adresse“; offen,
  **Schwelle erreicht**); `BEO-PGC/plan-zusage-erfuellung-ohne-committeten-anker`
  **2×** (F-5 LOW; offen); `BEO-PGC/arbeit-ueberholt-stehenden-traeger` **33×**
  (der Deckungsstand von `LH-FA-CFG-007` überholte zwei Träger, gefunden vom
  Suchlauf des Implementers nach dem Merge des bewegenden Slice, daher Datei
  trotz Deckel; verkörpert). *Deckel-Fälle ohne Datei, Finding-Kennung hier*
  (vor dem Merge von Reviewer bzw. Verifier gefunden, Schwere ≤ LOW, bekannter
  Träger-Typ): F-6 (LOW, `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`, Deckel
  bei 10×: die Meldung des fremden Trägers nannte im committeten Feld keine
  Frist; derselbe Vorgang wie die Datei zum Deckungsstand, V-3). *`state.md`
  fortgeschrieben:* die sechs Einträge mit neuem Zähler; bei den zwei Einträgen
  auf 3× steht, dass kein Slice fällig ist (ein Träger ist eine Zeile in Skill
  oder Regelwerk, also ein Architect-Zug des Lese-Schritts). *Kein
  Register-Anfall:* F-4 (ein Risiko, mit dem Runner-Beleg entfallen, §6), F-8,
  F-9, F-10 und V-1 (INFO, einmalig, im Plan oder Runner behoben, keine Klasse
  des Registers trägt sie), V-4 bis V-6 (Übergaben und übernommene Messungen, im
  Report benannt). *Lese-Schritt der Closure von `welle-transformationen`:* aus
  diesem Slice erreichen `test-name-behauptet-mehr-als-der-test-treibt` und
  `gemeldete-ungenauigkeit-ohne-traeger` die Schwelle ohne Ausgang; sie gehören in
  den Lese-Schritt der Welle-Closure, kein Slice ist wegen einer Beobachtung
  dieses Slice fällig.
- **Folge-Slices:** keine angelegt. Nächster Slice der Welle nach der
  Tabellenreihenfolge (§4 der Welle): `slice-transformationen-start-reihenfolge`
  (`open/`); sein Start-Trigger (§4 dort) ist mit dem Move dieses Slice erfüllt:
  `slice-transformationen-e2e-wirkung`, `slice-backfill-sql-administration` und
  `slice-antragsqueue-lesefehler-failed` liegen in `done/`, in `in-progress/`
  liegt nur die Roadmap (gelesen mit `ls` in der Closure). Danach folgt
  `slice-transformationen-e2e-abhilfe` (Start nach `slice-capture-leerlauf-quellbelege`),
  zuletzt `slice-transformationen-betriebsdoku` (Start nach
  `slice-sdk-regel-realserver-e2e`, dessen Start-Trigger mit diesem Move ebenfalls
  erfüllt ist). Übergaben mit Adresse (gemeldet, Frist: diese Closure, gezogen):
  `welle-transformationen` §3 (Restfläche der Zustellwege; Tabellenzeile des Slice
  und Herkunft der Deckung); `slice-sdk-regel-realserver-e2e` (der Pfad mit
  Verzeichnis in der Träger-Tabelle als Kennung).
- **Risiken aus §6:** je ein Ausgang, mit Beleg in §6. *Entfallen:* ein Zustellweg
  trägt eine abweichende Form (für das INSERT-Neu-Bild, Restfläche benannt) ·
  Laufzeit (579 s und 491 s im Lauf von `e2e.yml`) · stiller Ausschluss aus dem
  Runner (16 von 16) · geteilter Zustand · Neustart-Beleg misst Timing · Fenster
  READY → erste Zeile (Runner-Beleg, beide Legs). *Eingetreten:* das Handbuch
  entsteht erst danach, Folge-Slice mit Kennung
  `slice-transformationen-betriebsdoku`. *Weiter offen:* keines.
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-transformationen](../welle-transformationen.md) (offen) — die Prüfung
  läuft regelkonform bei deren Closure. Vorab gelesen: *Anker* — dieser Slice
  trägt kein Feld „liegt in“; *Folge-Slice* — `slice-transformationen-start-reihenfolge`,
  `slice-transformationen-e2e-abhilfe`, `slice-transformationen-betriebsdoku`,
  `slice-sdk-regel-realserver-e2e` liegen als Dateien in `open/`; *Register* —
  die sechs genannten Kennungen `BEO-PGC/…` sind Verzeichnisse mit nicht leerem
  `evidence/` (gemessen in der Closure).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Test-Runner und Wegwerf-Clients sind keine eigenen
Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Zähler
zum Stand der Planung; die Zähler der Closure stehen in §7) —
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
