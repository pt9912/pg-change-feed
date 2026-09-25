# Slice backfill-bench-richtgroesse: Bench der Kopierdauer je Tabellengröße — Warn-Richtgröße aus der Messung, Auswertung der beiden Warnungen im Use Case des Runs

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-backfill-bestand](../welle-backfill-bestand.md).

**Bezug:** [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Out-of-Scope: Parallelisierung/Durchsatz über sehr
große Tabellen ist Ausbaustufe — hier die Messung, ab wann eine Tabelle dafür
in Frage kommt), [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Festlegung 3 (große Tabellen: Warnung, keine
Ablehnung; Richtgröße zu messen), Teilfrage 4 (Ein-Transaktions-Form),
[`ADR-0054`](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md) §(b) (Benchmark-Infrastruktur), [`ADR-0104`](../../adr/0104-benchmark-schwellen-per-001-002-003.md) (Benchmark-Schwellen),
[`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 3 (Warn-Kriterium: Toleranz als Startwert, zwei Warnungen, Richtgröße
aus der Messung).

**Berührte Spec-Stellen:** [`SPEC-025`](../../../../spec/pflichtenheft.md) (Bench-Schwellen — gelesen; **kein** neuer
Wert in §3 ohne schärfende ADR, siehe §6), [`SPEC-029`](../../../../spec/pflichtenheft.md) (Run-Zustand mit
`estimated_rows`) — gelesen, nicht geändert.

**Verantwortlich:** Implementer-Agent, 2026-09-25.

**Autor:** Planner-Agent, Welle-Eröffnung [welle-backfill-bestand](../welle-backfill-bestand.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Aus einer Messung entsteht die **Warn-Richtgröße** für große
Tabellen — und aus ihr eine Warnung, nie eine Ablehnung.

- **Bench.** Ein viertes eigenständiges Bench-Skript im Muster von
  `tools/bench-source-impact.sh`, `tools/bench-scaling.sh` und
  `tools/bench-batch-vs-single.sh` (Arbeitsname `tools/bench-backfill.sh`, dieselbe
  Umgebung über `tools/bench-lib.sh`): Kopierdauer eines Runs je Tabellengröße
  (`finished_at − started_at` des Runs; die Wartezeit in `queued` zählt nicht), an einer Stufenfolge von Tabellengrößen (die Stufen
  legt der Slice fest und nennt sie im Bericht); gemessen wird auch, wie nah die
  **geschätzte** Zeilenzahl (`pg_class.reltuples`) der tatsächlichen kommt — in
  einer frisch befüllten, in einer analysierten und in einer nie analysierten
  Tabelle. `make bench` fährt das Skript mit; es ist eine **Messung** mit
  gedruckter Zahl, kein Pass/Fail gegen eine Schwelle (die Richtgröße ist ihr
  Ergebnis, nicht ihre Vorbedingung).
- **Toleranz, Richtgröße und Warnungen** ([`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 3). Die **Toleranz**
  ist eine benannte Konstante an genau **einer** Stelle im Use-Case-Paket des Runs,
  mit dem Startwert 10 Minuten Kopierdauer — „Startwert, Setzung ohne Messung";
  die Uhr ist der `ClockPort`. Die **Richtgröße in Zeilen** folgt aus der Messung
  dieses Slice (die Zeilenzahl, die der Bench in der Toleranz kopiert, auf eine
  runde Zahl abgerundet, die Rundungsregel nennt der Bericht; kopiert der Bench
  nicht die volle Toleranzdauer, ist die Zahl aus der gemessenen Rate
  hochgerechnet und als **abgeleitet** gekennzeichnet) und steht als
  benannte Konstante im Code (Anker: der gemessene Lauf) und im Handbuch
  (Abschnitt „Grenzwerte“, mit Ursprung, Host und Lauf). Die **Auswertung beider
  Warnungen** liegt an **einer** Stelle im Use Case des Runs — nicht in
  `bootstrap.Diagnose`: Warnung (1) beim Antrag, wenn die **geschätzte**
  Zeilenzahl über der Richtgröße liegt; Warnung (2) zur Laufzeit, wenn ein Run
  länger als die Toleranz läuft, geprüft bei jedem Fortschritts-Update (je Block)
  und beim Abschluss. Das Ergebnis steht in den zwei Warn-Spalten der Run-Zeile —
  Warnung (1) schreibt `Admit`, Warnung (2) das Fortschritts-Update des Run-Zustands-Adapters (`UPDATE`) —,
  `cdc.backfill_status` reicht sie durch — die View trägt die zwei Spalten seit
  `slice-backfill-sql-administration` in ihrer Signatur, dieser Slice **ändert die
  Signatur nicht** und füllt nur Werte —, `diagnose` liest sie aus der View.
  Eine unbekannte Schätzung (`NULL`) warnt nicht (Warnung (1) entfällt), sagt
  aber „unbekannt"; Warnung (2) greift unabhängig von der Schätzung. **Keine
  Ablehnung, kein Abbruch, keine Statusänderung** ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Festlegung 3).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Pass/Fail-Schwellenwert im Bench** — [`ADR-0104`](../../adr/0104-benchmark-schwellen-per-001-002-003.md) regelt die
  Schwellen der drei bestehenden Bench-Kennungen; die Richtgröße ist eine
  gemessene Orientierung, keine Anforderung. Ein Schwellenwert im Bench wäre eine
  neue Gate-Aussage ([`AGENTS.md`](../../../../AGENTS.md) §3.6).
- **Ein Eintrag von Toleranz oder Richtgröße in `spec/pflichtenheft.md` §3** — die
  Regel dieses Abschnitts verlangt, dass eine ADR den Wert in ihrem `Schärft:`-Feld
  deklariert; weder [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) noch [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) tut das (Festlegung 3, Punkt 1: die
  Toleranz steuert nur eine Warnung). Soll einer der Werte Vertrag werden, braucht
  er eine Folge-ADR.
- **Ein Konfigurationsschlüssel für Toleranz oder Richtgröße** — die Konstante ist
  billiger nachzuschärfen als eine neue Betreiber-Oberfläche; Konfigurierbarkeit ist
  ein Re-Evaluierungs-Trigger von [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) (Folge-ADR).
- **Checkpoint, Block-Transaktionen, Parallelisierung** — der
  Re-Evaluierungs-Trigger für sie ist „Kopierdauer über der Betriebs-Toleranz"
  ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md)); diese Messung liefert dafür Zahlen, entscheidet ihn nicht.
- **Aussagen über andere Hardware** — die Zahl gilt für den gemessenen Host und
  nennt ihn.

## 2. Definition of Done

- [x] Bench: `tools/bench-backfill.sh` (Arbeitsname) ist in `make bench`
      verdrahtet (als letztes Skript) und druckt die Kopierdauer je
      Tabellengröße und die Abweichung der Schätzung von der tatsächlichen
      Zeilenzahl; jede gedruckte Zahl trägt Größe und Lauf. *Erreichbarer
      Beleg — das Kriterium „ein realer `make bench`-Lauf mit Exit 0“ ist auf
      dem Messhost **nicht erreichbar** und gilt nicht als erfüllt:* (a) das
      Skript einzeln mit Exit 0, gedruckte Zeilen mit Größe und Lauf
      (`20260925T012459Z`, Implementer-Nachmessung der Fixrunde;
      `20260925T011036Z`, Review-Report; `20260925T000439Z`, übernommen aus dem
      Lauf-Bericht des Implementer-Laufs, im Repository nicht auflösbar); (b)
      `make bench` endet mit Exit 2 am vorhandenen Skript
      `tools/bench-source-impact.sh` (`LH-QA-PER-001`): gemessen in der Fixrunde
      Overhead 91,7 % gegen die 35-%-Schwelle (ohne CDC 16.094 ms, mit CDC
      30.848 ms, Median von je 5 Läufen); übernommen 87,5 % und 93,9 %
      (Implementer-Läufe, im Repository nicht auflösbar) und 95,8 % am Parent
      `933ab054`; das Architect-Verdikt (Ursache: Latenz des Festschreibens auf
      dem Datenträger des Hosts, `pg_test_fsync` `fdatasync` 2.956 µs je
      Operation; Host-Last widerlegt) führt sie als Eigenschaft dieses Hosts —
      keine Schwellen- und keine Skript-Änderung
      ([`AGENTS.md`](../../../../AGENTS.md) §3.6), `make bench` ist kein Gate.
- [x] Toleranz, Richtgröße und Warnungen: die Richtgröße folgt aus der Messung
      (Herleitung im Bericht, jede Zahl mit Ursprung: gemessen · übernommen ·
      abgeleitet; Host und Lauf genannt); die Toleranz-Konstante (Startwert 10
      Minuten) und die Richtgröße stehen je an **genau einer** Stelle im Code; die
      Auswertung beider Warnungen liegt im Use Case des Runs und nicht in
      `bootstrap.Diagnose`; eine Warnung ändert weder Status noch Ablauf, und keine
      Stelle lehnt einen Antrag wegen der Größe ab. *Zu belegen durch:* Unit-Tests
      mit Fake-Uhr und je einer Mutation (`make test`) — Warnung (2) genau an der
      Toleranz (darunter keine, darüber gesetzt), Warnung (1) an der Richtgröße
      (genau auf, eins darüber), unbekannte Schätzung → keine Warnung (1) und nie
      `0`; ein Store-Test, dass die View das Ergebnis der Run-Zeile durchreicht
      (`make test-store`); der Suchlauf in §3. Die Signatur von
      `cdc.backfill_status` bleibt unverändert (Spaltenzahl, -reihenfolge,
      -namen und -typen wie am Start dieses Slice; eine Signaturänderung einer
      bestehenden View kostet nach
      [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md) ein
      Lesefenster im Rollout) — *zu belegen durch:* `git diff` gegen den
      Start-Stand zeigt keine Änderung an der View-Definition in
      `tools/schema/schema.yaml`.
- [x] Benennung: jede Stelle, die die Zeilenzahl nennt, trägt das Wort „geschätzt";
      die Richtgröße heißt „Richtgröße" oder „Orientierung", nie „Grenze", „Limit"
      oder „maximal"; die Toleranz heißt „Startwert" mit dem Zusatz „Setzung ohne
      Messung" ([`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 3, Punkt 6). *Zu belegen durch:* Review des Diffs
      (kein Gate) und der Suchlauf in §3.
- [x] Träger: die `make bench`-Beschreibung ist auf vier Skripte gezogen
      (Makefile-Kommentar und -Hilfetext, `harness/README.md` §Sensors,
      ggf. `docs/user/bench-abdeckung.md`); das Handbuch trägt die Änderungshistorie.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
      Report `review-slice-backfill-bench-richtgroesse` (0 HIGH, 2 MEDIUM, 4 LOW, 6 INFO);
      kein offenes HIGH/MEDIUM nach der Fixrunde (F-1/F-2 nachgezogen, F-3 bis F-6
      behoben, F-7/F-8 ohne Aktion, F-9 bis F-12 im Handbuch/Plan gekennzeichnet).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: Handbuch „Grenzwerte“ nennt die Toleranz als „Startwert, Setzung ohne Messung" und die Richtgröße mit ihrem Ursprung (gemessen oder abgeleitet), Host und Lauf; Änderungshistorie eine Zeile.
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
| `tools/bench-backfill.sh` (Arbeitsname) | neu | die Messung im Muster der drei bestehenden Skripte; nutzt `tools/bench-lib.sh`. |
| `Makefile` (Target `bench`) | update | ruft das vierte Skript **zuletzt**, hinter `tools/bench-batch-vs-single.sh` (Fixrunde, Review F-12): der Abbruch einer Messung ohne Schwelle verdeckt so weder eine Schwellen-Prüfung noch die Erzeugung von `docs/user/bench-abdeckung.md`; Hilfetext und Kommentar zählen neu. |
| `internal/application/usecase/backfill/` (+ Tests, Arbeitsname `warn.go`) | update | Toleranz- und Richtgrößen-Konstante, Auswertung beider Warnungen an einer Stelle; Fake-Uhr-Tests der Grenzfälle. |
| Run-Zustands-Port, Annahme-Port und ihre Adapter | prüfen | tragen das Warn-Ergebnis in die zwei Warn-Spalten aus `run-store` (Warnung 1 über `Admit`, Warnung 2 über das Fortschritts-Update); ein Bedarf über die dort angelegten Spalten hinaus wäre ein Plan-Nachzug. |
| `internal/bootstrap/backfill.go` (`diagnoseBackfillStatus`, `formatBackfillRun`; von `Diagnose` in `wiring.go` gerufen) | prüfen | liest die zwei Warn-Spalten aus der View (`sql-administration`, Abfrage `WHERE source_id = $1`); enthält keine Auswertung und keine Konstante. |
| `internal/bootstrap/diagnose_test.go` | prüfen | zeigt die Warnung aus der View an; keine Grenzfälle der Auswertung (die liegen im Use-Case-Test). |
| `docs/user/benutzerhandbuch.md` | update | Richtgröße im Abschnitt „Grenzwerte“ mit Ursprung; Diagnose-Beispiel; Änderungshistorie. |
| `harness/README.md` §Sensors | update | Zeile `make bench` (drei → vier Skripte) — aus dem realen Lauf geschrieben. |
| `tools/schema/schema.yaml` | prüfen (nicht ändern) | die View `backfill_status` trägt die zwei Warn-Spalten seit `slice-backfill-sql-administration`; ein Bedarf an einer Signaturänderung wäre ein Plan-Nachzug und ein Rückführungsgrund (§4), keine stille Änderung. |
| `docs/user/bench-abdeckung.md` | prüfen | die Datei bindet `LH-QA-PER-001`…`003` an durchgesetzte Schwellen; eine Zeile ohne Schwelle für `LH-FA-CAP-009` passt nicht zu dieser Form — Entscheidung im Slice, im Bericht begründet. **Ergebnis:** regeneriert, keine neue Zeile; nur der vom Generator erzeugte Kopftext nennt, dass `tools/bench-backfill.sh` ohne Schwelle misst und keine Zeile trägt. |
| `tools/bench-lib.sh`, `tools/bench-source-impact.sh` | update (Plan-Nachzug) | `bench::median_of` steht einmal in der Bibliothek statt im Skript (das vierte Skript nutzt sie mit); Kopf- und Kommentartexte und der erzeugte Kopf von `docs/user/bench-abdeckung.md` zählen vier Skripte. |
| `harness/targets/bench-backfill.md` | neu (Plan-Nachzug) | Vertrag, Phasen, Parameter und Grenzen des Skripts an einem Ort; der Makefile-Kommentar bleibt bei einem Verweis. |
| `internal/application/usecase/backfill/export_test.go`, `warn_test.go` | neu (Plan-Nachzug) | Prüf-Zugriff der Tests des Pakets `backfill_test` auf die beiden Konstanten und Whitebox-Tests der Grenzfälle von `warn.go`. |
| `internal/domain/model/backfillrun.go` | update (Kommentar, Plan-Nachzug) | der Doc-Kommentar der Warn-Kennzeichnungen nennt die Auswertung des Use Cases als ihren Setzer. |
| `internal/adapters/driven/postgresstorage/backfillstatusview_test.go` | update (Plan-Nachzug) | die View reicht `warn_estimated_size` durch: der Test prüft neben `warn_duration` auch diese Spalte; DoD-Beleg „View trägt das Ergebnis“ (`make test-store`). Fixrunde (Review F-6): die Assertionen zeigen den Wert der Schätzung statt des Zeigers; die Bindung an die View-Spalte bleibt (Mutation `false AS warn_estimated_size` rot, Meldung „vwst_e: Schätzung 7000000, Warnungen false/false“). |
| `tools/bench-backfill.sh` (Fixrunde) | update (Plan-Nachzug) | Review F-4: `BENCH_BACKFILL_RUN_TIMEOUT_S` zählt Sekunden der Shell-Uhr (`SECONDS`), nicht Abfragen; F-5: `range_of` bildet den Median über `bench::median_of`, eine Stelle statt zwei Formeln. |
| `harness/targets/bench-backfill.md`, `harness/README.md` (Fixrunde) | update (Plan-Nachzug) | Review F-3: Vertrag und `make bench`-Zeile nennen die WAL-Messung des Skripts (Rückstand, gehaltenes WAL, Live-Commit-Freigabe, abgeleitete WAL-Schwellen-Zeile) und den Speicher 20 s nach dem letzten Run; F-4: Einheit der Variable. |
| `docs/user/benutzerhandbuch.md` (Fixrunde) | update (Plan-Nachzug) | 1.54: F-1 (WAL-Rückstand nicht an den Backfill gebunden, Folge-Slice, Betriebs-Abhilfen), F-2 (Speicher: Spitze im Run und Probe 20 s danach), F-9 (drei Läufe zur Richtgröße), F-10/F-11 (Live-Wirkung als drei Einzelläufe; nicht auflösbare Läufe als übernommen gekennzeichnet). |

**Übergabe aus `slice-backfill-snapshot-reader`** (gemeldet, kein zusätzlicher Umfang):
der Snapshot-Leser liest in Blöcken von `DefaultBlockSize` = 1.000 Zeilen, ein Startwert
ohne Messung; `B` zählt Zeilen, nicht Bytes, der Speicherbedarf eines Blocks ist `B` mal die
Zeilenbreite (Port-Doku `NextBlock`, nicht gemessen). Der Bench nennt `B` und die
Zeilenbreite seiner Tabellen im Ursprung jeder Kopier-Zahl; eine Richtgröße in Zeilen gilt
für diese Breite, Messungen mit breiten Zeilen (`jsonb`, `bytea`) sind ohne eigenen Lauf
ungemessen.

**Übergabe aus `slice-backfill-run-usecase`** (gemeldet, kein zusätzlicher Umfang):
`Request` legt den Run mit der Schätzung an (`RowEstimate`: der Nullwert ist „unbekannt",
die Zahl 0 einer analysierten, leeren Tabelle bleibt bekannt) und ruft `Admit` als letzten
Schritt; die Stelle für Warnung (1) ist der Run-Wert zwischen `WithEstimatedRows` und
`Admit`. Das Fortschritts-Update ist `progress` in
`internal/application/usecase/backfill/service.go`: je Block und einmal mit der
Anfangsposition (Zähler 0) ruft es `BackfillRunPort.RecordProgress` mit dem **ganzen** Run;
Warnung (2) setzt der Slice dort am Run-Wert, `Finish` und `Commit` tragen den Run beim
Abschluss mit denselben Feldern (`WarnEstimatedSize`, `WarnDuration`, bis zur Auswertung
`false`). Der Test-Baustein für die Uhr ist im Use-Case-Test vorhanden (`fakeClock`
läuft je Aufruf um einen Schritt weiter).

**Übergabe aus `slice-backfill-run-store`** (gemeldet, kein zusätzlicher Umfang): die
Postgres-Adapter des Runs tragen zwei Startwerte ohne Messung — die Fristen (30 s je
Zustands-Operation, 5 min je Block und Commit; `backfillStateTimeout` und
`backfillBlockTimeout` in `internal/adapters/driven/postgresstorage/backfill.go`) und die
zeilenweise Einfügung eines Blocks (`AppendBlock`: je Change eine Anweisung in der einen
Schreibtransaktion). Die Kopierdauer je Tabellengröße, die der Bench misst, schließt die
Einfügeform ein; der Bench nennt sie im Ursprung jeder Kopier-Zahl, und eine Blockdauer nahe
oder über der Frist steht als Befund im Bericht. Beide Setzungen bleiben bis dahin ungemessen
(Review F-7 zu `slice-backfill-run-store`, INFO).

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Zahl und der Inhalt der Bench-Skripte hinter `make bench`"; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Zeile `make bench` in `harness/README.md` („drei eigenständige Bench-Skripte") | `grep -rn 'drei' harness/README.md docs/user/bench-abdeckung.md Makefile` (Fundstellen nach Bench-Bezug lesen); ergänzt um die Zählwörter `git grep -n -i -E '(drei\|vier\|zwei\|beiden\|sechs\|sieben\|fünf)[^\|]{0,40}bench' <Stand>` | **Parent** `933ab054`: `harness/README.md` Zeile 150 „drei eigenständige Bench-Skripte“, `Makefile` Zeile 188 „drei Skripte“, `docs/user/bench-abdeckung.md` Zeile 3 „drei `tools/bench-*.sh`-Skripten“, `tools/bench-lib.sh` Zeilen 9 „drei Bench-Skripte“, 145 „der beiden anderen Skripte“, 160 „letzten der drei Bench-Skripte“, 171 (erzeugter Kopftext) „drei“; `tools/bench-scaling.sh` Zeilen 3/32/100/103 „drei“ meinen die drei `SPEC-014`-Lastenstufen (anderer Gegenstand). **Diff-Stand:** `harness/README.md` „vier … (drei mit Schwelle, eines ohne)“, `Makefile` „vier Skripte“ und „drei mit Schwelle“, `docs/user/bench-abdeckung.md` „drei … mit Pass/Fail-Schwelle“, `tools/bench-lib.sh` „vier“, „der anderen“, „drei Schwellen-Skripte“ und der erzeugte Kopftext. **Nicht gefunden:** kein Träger in `.github/`, `AGENTS.md`, `README.md`, `docs/user/benutzerhandbuch.md`, `spec/` nennt die Zahl der Bench-Skripte; `ADR-0054` und `ADR-0104` nennen „drei Skripte“ für die drei `LH-QA-PER`-Belege (ihr Gegenstand, `Accepted`, unverändert) | nachgezogen |
| Hilfetext und Kommentar des Targets `bench` | Lesen von `Makefile` an beiden Ständen | **Parent:** Kommentar (Zeilen 181–186) „Drei eigenständige Skripte (`LH-QA-PER-001`…`003`)“ und Hilfetext „drei Skripte“; die Rezeptzeilen rufen drei Skripte. **Diff-Stand:** Kommentar „Vier eigenständige Skripte“ mit der Aufteilung drei mit Schwelle / eines ohne und einem Verweis auf `harness/targets/bench-backfill.md`, Hilfetext „vier Skripte“, vier Rezeptzeilen. **Nicht gefunden:** kein weiteres Makefile-Ziel und keine `harness/mk/*.mk`-Datei ruft ein `tools/bench-*.sh` | nachgezogen |
| Bench-Abdeckungs-Träger | Lesen von `docs/user/bench-abdeckung.md` und dem Generator in `tools/bench-batch-vs-single.sh` und `tools/bench-lib.sh` (`bench::render_abdeckung`) | **Parent:** die Datei trägt drei Zeilen (`LH-QA-PER-001`…`003`) und den Kopftext „drei `tools/bench-*.sh`-Skripten“; der Generator ist `bench::render_abdeckung` in `tools/bench-lib.sh`, aufgerufen vom letzten Schwellen-Skript `tools/bench-batch-vs-single.sh`. **Diff-Stand:** dieselben drei Zeilen, der Kopftext nennt die drei Schwellen-Skripte und dass `tools/bench-backfill.sh` ohne Schwelle misst und keine Zeile trägt; die Datei ist aus einem realen Lauf von `tools/bench-batch-vs-single.sh` regeneriert. **Nicht gefunden:** kein zweiter Generator, kein Träger außerhalb `bench-lib.sh`, der den Kopftext spiegelt | Entscheidung siehe §3, Zeile `bench-abdeckung.md`: keine Zeile für `LH-FA-CAP-009` |
| Handbuch-Grenzwerte | Lesen von §Grenzwerte in `docs/user/benutzerhandbuch.md` an beiden Ständen; ergänzt um die Zählwörter in den Backfill-Absätzen desselben Handbuchs (`grep -n 'drei\|vier\|zwei\|beiden' docs/user/benutzerhandbuch.md` im Abschnitt „Bestand als Backfill überführen“) | **Parent:** §Grenzwerte trägt fünf Aufzählungspunkte (Lebenszeichen-Takt, WAL-Rückstand-Messtakt, WAL-Schwellen, eine Quelle je Container, zwei Replication-Verbindungen mit der Backfill-Ergänzung) und keine Backfill-Zahl; der Absatz „Für den Lauf selbst gelten drei Betriebs-Vorbedingungen“ zählt drei Punkte (`SELECT`-Recht, Reserve, Snapshot-Haltedauer). **Diff-Stand:** §Grenzwerte trägt zusätzlich sieben Backfill-Punkte (Toleranz, Richtgröße, gemessene Werte, Speicher, WAL-Rückstand, Live-Wirkung, Schätzung); der Absatz zählt „vier“ und trägt den Punkt WAL-Rückstand. **Nicht gefunden:** kein weiterer Träger im Handbuch nennt die Zahl der Backfill-Vorbedingungen; die Replication-Verbindungs-Zahl im Punkt „zwei gleichzeitige Replication-Protokoll-Verbindungen“ blieb unberührt | Richtgröße ergänzt; die Replication-Verbindungs-Zahl ist Sache von `sql-administration` |
| Toleranz-Konstante an genau einer Stelle; Wortwahl der Richtgröße | `grep -rn 'Toleranz\|Richtgröße\|Orientierung' internal docs harness --include=*.go --include=*.md` (Konstante, Handbuch, Kommentare), an beiden Ständen | **Parent:** kein Treffer für die Toleranz oder die Richtgröße als Konstante; `Orientierung, keine Grenze` steht an `internal/application/port/outbound/tablesnapshot.go` (Schätzung) und `internal/domain/model/backfillrun.go`; `spec/pflichtenheft.md` nennt beide Begriffe in `LH-FA-CAP-009.a` und `SPEC-029`. **Diff-Stand:** die Toleranz ist genau einmal definiert (`copyDurationToleranceMinutes` in `internal/application/usecase/backfill/warn.go`, `copyDurationToleranceNanos` ist ihre Umrechnung), die Richtgröße genau einmal (`estimatedRowsGuideline` ebenda); die Tests lesen beide über `export_test.go`, `tools/bench-backfill.sh` liest die Toleranz aus `warn.go` (keine zweite Zahl; Mutation: Konstante umbenannt → Exit 1); das Handbuch nennt die Toleranz „Startwert, Setzung ohne Messung“ und die Richtgröße „Orientierung, keine Grenze“ mit Ursprung. „Grenze“ steht am Begriff nur in der Verneinung („keine Grenze“); kein „Limit“, kein „maximal“ an Toleranz oder Richtgröße. **Nicht gefunden:** keine Zeilenzahl-Konstante außerhalb `warn.go` (`grep -rn -e '4_000_000' -e '4\.000\.000' internal docs harness spec tools Makefile` trifft nur `warn.go` und die eine Handbuch-Angabe; das Beispiel der Diagnose-Ausgabe trägt eigene, erfundene Zahlen); kein Eintrag in `spec/pflichtenheft.md` §3 | eine Definition der Toleranz; „geschätzt" an jeder Nennung der Zeilenzahl; kein „Grenze"/„Limit"/„maximal" am Begriff |
| Build-Kontext | Lesen von `.dockerignore` und `Dockerfile` an beiden Ständen | **Parent und Diff-Stand:** `.dockerignore` schließt mit `*` alles aus und öffnet nur `cmd/`, `internal/`, `gen/`, `proto/` sowie die zwei einzelnen Dateien der `ADR-0085`-Klasse; das `Dockerfile` kopiert per `COPY . .` in den Stufen `coverage` (Zeile 98) und `build` (Zeile 123) und per `COPY proto/` in `proto-export`. **Nicht gefunden:** kein `COPY`, das `tools/bench-*.sh`, `harness/` oder `docs/` in einen Kontext bringt; die neue Datei `warn.go` liegt unter `internal/` (Go-Bau, im Kontext); keine neue `!`-Zeile nötig | Bench-Skripte laufen auf dem Host; keine Docker-Stufe kopiert `tools/bench-*.sh` — sonst Freigabe nach [`ADR-0085`](../../adr/0085-build-kontext-ausnahme-test-only-zweck.md) |
| zweite bewegte Eigenschaft: „die zwei Warn-Kennzeichnungen sind ohne Auswertung `false`“ | `git grep -n -i 'keine Auswertung\|beide bleiben\|bleiben .false\|bis zur Auswertung\|solange keine' <Stand> -- spec harness internal docs/user Makefile tools examples sdks` und `grep -n 'warn_estimated_size\|warn_duration\|Warnung' <Träger>` | **Parent:** `docs/user/benutzerhandbuch.md` Zeile 524 „In dieser Version setzt keine Auswertung sie: beide bleiben `false`“ und `internal/domain/model/backfillrun.go` Zeilen 64–65 „bleiben `false`, solange keine Auswertung sie setzt“; `spec/pflichtenheft.md` Zeilen 704/705 sagen „bleibt `false`, solange … keine Richtgröße festgelegt ist“ bzw. „gesetzt beim Fortschritts-Update (je Block) oder beim Abschluss“ (Vertrag der Auswertung, keine Aussage über einen Ist-Stand); `tools/schema/schema.yaml` Zeilen 273/408 „false, solange keine Auswertung sie setzt“ (Metadaten-Text der Tabelle und der View). **Diff-Stand:** Handbuch und `backfillrun.go` nennen die Auswertung des Use Cases als Setzer; die Spezifikation bleibt unverändert wahr. **Nicht gefunden:** kein weiterer Träger in `harness/`, `internal/bootstrap/` (`diagnose` zeigt die View-Werte ohne Aussage über den Setzer), `examples/`, `sdks/` | nachgezogen; `tools/schema/schema.yaml` bleibt unverändert (die DoD bindet „keine Änderung an der View-Definition“, und der Text „false, solange keine Auswertung sie setzt“ bleibt als Bedingung wahr) |
| **Fixrunde** — bewegte Eigenschaft: „der WAL-Rückstand hängt am Run/Backfill“ (Verdikt: hängt an WAL ohne Inhalt für die Publication) | `git grep -c -i -E 'WAL-Rückstand\|Fehlerschwelle\|wal_retention_(error\|warn)' <Stand> -- docs/user harness tools Makefile examples sdks internal README.md AGENTS.md ':!docs/plan' ':!docs/reviews'` (Treffer nach der Aussage über den Backfill gelesen) | **Parent** `50aa8da8`: 11 Dateien, davon Träger mit einer Aussage über den Run: `docs/user/benutzerhandbuch.md` (23 Treffer; §4-Punkt „WAL-Rückstand des Capture-Slots“, §9-Punkt „Backfill, WAL-Rückstand des Capture-Slots“) und dieser Plan (§3 „Befunde der Messung“, DoD „Bench“). **Diff-Stand:** 13 Dateien; Handbuch (27 Treffer), Plan, dazu `harness/README.md` (1) und `harness/targets/bench-backfill.md` (3) mit der Beschreibung der WAL-Messung des Skripts. **Nicht gefunden:** die übrigen Treffer sind der Schwellen-Code und seine Tests (`internal/bootstrap/`, `internal/adapters/driving/replication/receive/`) und ein Lokator in `harness/sensors/coverage-gate.md` — keine Aussage über den Backfill; kein Treffer in `examples/`, `sdks/`, `README.md`, `AGENTS.md`; `docs/plan/adr/**` und `docs/reviews/**` (`Accepted`/Records) bewusst ausgenommen; der offene Slice `slice-backfill-slot-leerlauf-bestaetigung` trägt die Nachzüge der Umsetzung selbst | nachgezogen (Handbuch, Plan) |
| **Fixrunde** — bewegte Eigenschaft: „Einheit und Bedeutung von `BENCH_BACKFILL_RUN_TIMEOUT_S`“ | `git grep -c 'RUN_TIMEOUT' <Stand> -- . ':!docs/reviews' ':!docs/plan/adr'` | **Parent:** `tools/bench-backfill.sh` (3), `harness/targets/bench-backfill.md` (1). **Diff-Stand:** dieselben zwei Dateien mit gleicher Trefferzahl (dieser Plan nennt die Variable in seinen eigenen Zeilen); Zählung in Sekunden der Shell-Uhr, Vertrag „Sekunden“. **Nicht gefunden:** kein Träger in `Makefile`, `harness/README.md`, Handbuch, `.github/` | nachgezogen |
| **Fixrunde** — bewegte Eigenschaft: „Reihenfolge der Skripte in `make bench`“ | `git grep -c 'bench-batch-vs-single' <Stand> -- . ':!docs/reviews' ':!docs/plan/adr'` | **Parent:** `Makefile` (1, Rezeptzeile vor dem neuen Skript), `harness/README.md` (1), `tools/bench-lib.sh` (1), `docs/user/bench-abdeckung.md` (1), `tools/bench-batch-vs-single.sh` (10) und Records unter `docs/plan/planning/done/` und `observations/` (6), dazu dieser Plan. **Diff-Stand:** gleiche Verteilung ohne diesen Plan (er nennt das Skript in seinen eigenen Zeilen); das Rezept ruft das neue Skript zuletzt. `harness/README.md` nennt die vier Skripte bereits in der Reihenfolge Quelle, Skalierung, Batch, Backfill; der Kommentar in `tools/bench-lib.sh` („Schwellen-Skripte in Ausführungsreihenfolge“) bleibt wahr. **Nicht gefunden:** kein Träger nennt eine Reihenfolge, die der neuen widerspricht | Makefile-Kommentar nennt den Grund |
| **Fixrunde** — bewegte Eigenschaft: „Werte der Messung im Handbuch: Speicher-Spitze, Kopierrate, Live-Wirkung“ | `git grep -c -E '467 MiB' <Stand> -- docs/user harness docs/plan/planning/in-progress tools Makefile`, ebenso `'7\.693'` und `'kein Unterschied'` | **Parent:** `467 MiB` je 1 Treffer in Plan und Handbuch; `7\.693` 2 Treffer im Handbuch; `kein Unterschied` 1 Treffer im Handbuch („ist kein Unterschied gemessen“). **Diff-Stand:** `467 MiB` je 1 (jetzt mit den Werten dreier Läufe und der Probe 20 s danach); `7\.693` 3 im Handbuch und 1 im Plan; `kein Unterschied` 1 im Handbuch (Satz „kein Unterschied ableitbar“ aus drei Einzelläufen), 1 in der Versionshistorie-Zeile, 1 im Plan. **Nicht gefunden:** kein Träger in `harness/`, `tools/`, `Makefile`, der die Werte trägt | nachgezogen |

**Befunde der Messung** (Übergabe an Planner und Architect, kein zusätzlicher Umfang; Ursprung: Läufe von `tools/bench-backfill.sh` `20260924T233628Z` und `20260925T000439Z`, Zahlen und Host im Handbuch §Grenzwerte):

- **Die WAL-Fehlerschwelle kann vor der Richtgröße greifen — der WAL-Rückstand ist nicht an den Backfill gebunden.** Ein dritter Run über 1.000.000 Zeilen im selben CDC-Speicher überschritt die Fehlerschwelle von 1 GiB (`SPEC-013`): der Feed-Container endete mit Exit 1 (`replication`), der Run `interrupted` (übernommen aus dem Lauf-Bericht, im Repository nicht auflösbar). Beantwortet durch das Architect-Verdikt (`architect-verdict-backfill-wal-rueckstand-und-bench-rot`) und [`ADR-0120`](../../adr/0120-capture-slot-leerlauf-bestaetigung.md): der Slot bestätigt nur Commits, die Änderungen veröffentlichter Tabellen tragen; jedes WAL ohne solchen Inhalt hebt den Rückstand — ein Schreiber mit 200.000 Zeilen auf eine **nicht aktivierte** Tabelle hob ihn ohne Backfill um 33,5 MiB (gemessen im Verdikt, Lauf `wal-verdikt[20260925T004731Z]`, gedruckt in `ADR-0120` §Gemessen), ein Commit auf einer aktivierten Tabelle senkte ihn innerhalb von 3 s auf 0 MiB. Die Richtgröße (Zeit-Orientierung, [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)) bleibt unverändert und wird nicht an die WAL-Schwelle gekoppelt. Abhilfe: der Folge-Slice `slice-backfill-slot-leerlauf-bestaetigung` (Umsetzung von `ADR-0120`, kein Umfang dieses Slice); Betreiber-Abhilfen stehen im Handbuch (§4 „Bestand als Backfill überführen“, Punkt „WAL-Rückstand des Capture-Slots“): ein Commit auf einer aktivierten Tabelle und das Datei-Feld `wal_retention_error_bytes` (`SPEC-013`-Override, kein Umgebungsvariablen-Gegenstück).
- **Die Speicher-Spitze des Feed-Containers liegt weit über `B` mal Zeilenbreite** (Stufe 200.000 Zeilen, Spitze im Run: 467 MiB übernommen, 401,7 MiB gemessen im Lauf `20260925T012459Z`, 416,5 MiB im Review-Report; die Probe 20 s nach dem letzten Run liegt höher: 437,0 und 641,7 MiB; 1.544 MiB bei 1.000.000 Zeilen, übernommen; gegen etwa 74 KB je Block); die Ursache ist nicht untersucht.
- **`make bench` bricht am ersten roten Skript ab**; `LH-QA-PER-001` ist auf diesem Host rot (91,7 % in der Fixrunde gemessen, übernommen 87,5 % bis 95,8 %, auch am Parent). Ursache laut Verdikt: die Latenz des Festschreibens auf dem Datenträger des Messhosts (`pg_test_fsync` `fdatasync` 2.956 µs je Operation, gemessen im Verdikt; Host-Last widerlegt), keine Regression und keine CPU-Last; die 35-%-Schwelle ist damit eine Aussage über die Flush-Latenz des Messhosts. Der Slice berichtet nur: keine Schwellen- und keine Skript-Änderung ([`AGENTS.md`](../../../../AGENTS.md) §3.6), `make bench` ist kein Gate. Das neue Skript läuft als letztes und wird von `make bench` auf diesem Host nicht erreicht; es läuft einzeln.
- **Die Kopierrate streut zwischen Läufen um etwa den Faktor 2**; die Richtgröße hängt am Lauf, aus dem sie abgeleitet ist: Median der Stufe 200.000 Zeilen 7.693 (Lauf `20260925T000439Z`, übernommen), 8.933 (Lauf `20260925T011036Z`, Review-Report) und 8.559 Zeilen/s (Lauf `20260925T012459Z`, gemessen) ergeben 4.000.000, 5.000.000 und 5.000.000 Zeilen nach der Rundung auf eine Stelle; die Konstante trägt den kleinsten Wert (konservativer Startwert, Nachschärfen lässt `ADR-0113` zu).
- **Die Live-Wirkung ist ein Einzellauf-Vergleich:** Mediane von `cdc_capture_lag` im Run gegen ohne Run 0,48 gegen 0,50 s (übernommen), 0,59 gegen 0,40 s (Review-Report), 0,456 gegen 0,413 s (Lauf `20260925T012459Z`, gemessen) — in beiden Richtungen; keine Aussage „kein Unterschied“ ohne diese Einschränkung.

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `e2e` in `done/` liegt, kein
anderer Slice in `in-progress/` liegt **und** [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) den Status `Accepted` trägt
(geprüft an der Status-Spalte im ADR-Index): sie legt das Warn-Kriterium fest —
die Toleranz als Startwert 10 Minuten Kopierdauer, zwei Warnungen, die
Richtgröße erst aus der Messung (Festlegung 3). [`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) nennt „die
Betriebs-Toleranz" ohne Wert; ohne die Festlegung stünde ein selbst gewählter Wert
als Grenze im Code (`BEO-PGC/geschaetzter-wert-als-grenze`). Das ist eine
**Vorab-Bedingung** und steht am Start, nicht als Rückführung
(`BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Bench und
  Warn-Auswertung nicht in einem Review tragen — der abtrennbare Teil ist die
  Warn-Auswertung im Use Case samt ihren Tests.
- `in-progress` → `next` (Plan-Nachzug): falls die Warn-Auswertung eine
  Signaturänderung von `cdc.backfill_status` verlangt (über die zwei Warn-Spalten
  aus `slice-backfill-sql-administration` hinaus) — die Änderung ist nach
  [`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md) ein
  Lesefenster im Rollout und wird vor der Umsetzung im Plan benannt.
- `in-progress` → `open` (blockiert): falls die Schätzung (`reltuples`) so weit
  von der tatsächlichen Zeilenzahl abweicht, dass eine Warnung an ihr nichts
  trägt (dann Architect-Frage: statt der Schätzung eine andere Größe, etwa ein
  Seitenzähler der Tabelle).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `tools/bench-backfill.sh` einzeln mit Exit 0
und `make bench` mit dem in §2 (Bench) genannten Ausgang (Exit 2 an
`LH-QA-PER-001`, Ursache Host) + `make test` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Die Toleranz ist eine Setzung ohne Messung** ([`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 3: Startwert
  10 Minuten). Sie wird im Code und im Handbuch als „Startwert, Setzung ohne
  Messung" geführt und nicht als Ergebnis dieses Slice ausgegeben. *Erwartet, zu
  belegen durch:* der Bericht nennt die gemessene Kopierdauer je Stufe neben dem
  Startwert; weicht sie ab, ist das Nachschärfen der Konstante ein
  Re-Evaluierungs-Fall von [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md), kein Bruch. **Ausgang:** *(bei Closure)*
- **Die Richtgröße wird zur Grenze** (`BEO-PGC/geschaetzter-wert-als-grenze`,
  offen, 1×; dieser Slice ist der Ort, an dem der Fehler passiert): eine
  gemessene Zahl eines Hosts steht im Code und im Handbuch als allgemeine
  Grenze. *Erwartet, zu belegen durch:* jede Nennung trägt „gemessen auf …,
  Lauf …" und das Wort „Richtgröße"; keine Stelle lehnt ab. **Ausgang:** *(bei
  Closure)*
- **Die Schätzung ist so alt wie das letzte `ANALYZE`.** Eine frisch geladene,
  nie analysierte Tabelle trägt „unbekannt" (erwartet), die Warnung schweigt dann
  gerade bei den Tabellen, für die sie gedacht ist. *Erwartet, zu belegen durch:*
  die Schätz-Abweichungs-Messung; das Handbuch nennt den Zusammenhang.
  **Ausgang:** *(bei Closure)*
- **Die Messung hängt am Host** (Docker-Umgebung, Plattengeschwindigkeit). Das
  Handbuch nennt die Zahl als Orientierung mit ihrem Lauf; eine zweite
  Umgebung ist nicht Teil. **Ausgang:** *(bei Closure)*
- **Laufzeit von `make bench`** wächst; der Default-Modus fährt verkürzte Stufen
  (Muster `bench-scaling.sh`), ein Volllauf ist Option. *Erwartet, zu belegen
  durch:* Laufzeit im Bericht. **Ausgang:** *(bei Closure)*
- **Ein Wert in `spec/pflichtenheft.md` §3 ohne schärfende ADR** wäre gegen die
  Regel des Abschnitts; der Slice führt keinen (weder Toleranz noch Richtgröße).
  **Ausgang:** *(bei Closure)*
- **Warnung (2) prüft je Block und beim Abschluss:** ein einzelner Block, der
  länger als die Toleranz braucht, warnt erst danach — die Warnung ist eine
  Orientierung, kein Alarm (akzeptiertes Negativ von [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)). Das Handbuch sagt es
  so. **Ausgang:** *(bei Closure)*

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
`*`/`PGC` (Greenfield); das Bench-Werkzeug unter `tools/` und die Diagnose im
Bootstrap sind keine eigene Sub-Area — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/geschaetzter-wert-als-grenze` (offen, 1×, einschlägig — Start-Trigger und
Risiko §6), `BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft` (offen, 2×,
Start-Trigger), `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
(verkörpert, 13×, jede Zahl mit Ursprung), `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`
(verkörpert, 8×, die `make bench`-Zeile beschreibt nur, was der Lauf getan hat),
`BEO-PGC/coverage-stage-dockerignore-blockiert-tooling` (offen, 2×) und
`BEO-PGC/dockerignore-default-deny-blockiert-neuen-pfad` (offen, 1×) —
Suchlauf-Zeile Build-Kontext, `BEO-PGC/arbeit-ueberholt-stehenden-traeger`
(verkörpert, 26×, Suchlauf §3), `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`
(verkörpert, 6×, Warn-Tests mit Mutation).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
