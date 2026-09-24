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

- [ ] Bench: `tools/bench-backfill.sh` (Arbeitsname) läuft in `make bench` mit
      und druckt die Kopierdauer je Tabellengröße und die Abweichung der
      Schätzung von der tatsächlichen Zeilenzahl; jede gedruckte Zahl trägt
      Größe und Lauf. *Zu belegen durch:* ein realer `make bench`-Lauf (Exit
      ungefiltert gesichert), Ausgabe im Bericht; ein gemessener Lauf, nicht
      abgeleitet.
- [ ] Toleranz, Richtgröße und Warnungen: die Richtgröße folgt aus der Messung
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
- [ ] Benennung: jede Stelle, die die Zeilenzahl nennt, trägt das Wort „geschätzt";
      die Richtgröße heißt „Richtgröße" oder „Orientierung", nie „Grenze", „Limit"
      oder „maximal"; die Toleranz heißt „Startwert" mit dem Zusatz „Setzung ohne
      Messung" ([`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 3, Punkt 6). *Zu belegen durch:* Review des Diffs
      (kein Gate) und der Suchlauf in §3.
- [ ] Träger: die `make bench`-Beschreibung ist auf vier Skripte gezogen
      (Makefile-Kommentar und -Hilfetext, `harness/README.md` §Sensors,
      ggf. `docs/user/bench-abdeckung.md`); das Handbuch trägt die Änderungshistorie.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: Handbuch „Grenzwerte“ nennt die Toleranz als „Startwert, Setzung ohne Messung" und die Richtgröße mit ihrem Ursprung (gemessen oder abgeleitet), Host und Lauf; Änderungshistorie eine Zeile.
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
| `Makefile` (Target `bench`) | update | ruft das vierte Skript; Hilfetext und Kommentar zählen neu. |
| `internal/application/usecase/backfill/` (+ Tests, Arbeitsname `warn.go`) | update | Toleranz- und Richtgrößen-Konstante, Auswertung beider Warnungen an einer Stelle; Fake-Uhr-Tests der Grenzfälle. |
| Run-Zustands-Port, Annahme-Port und ihre Adapter | prüfen | tragen das Warn-Ergebnis in die zwei Warn-Spalten aus `run-store` (Warnung 1 über `Admit`, Warnung 2 über das Fortschritts-Update); ein Bedarf über die dort angelegten Spalten hinaus wäre ein Plan-Nachzug. |
| `internal/bootstrap/backfill.go` (`diagnoseBackfillStatus`, `formatBackfillRun`; von `Diagnose` in `wiring.go` gerufen) | prüfen | liest die zwei Warn-Spalten aus der View (`sql-administration`, Abfrage `WHERE source_id = $1`); enthält keine Auswertung und keine Konstante. |
| `internal/bootstrap/diagnose_test.go` | prüfen | zeigt die Warnung aus der View an; keine Grenzfälle der Auswertung (die liegen im Use-Case-Test). |
| `docs/user/benutzerhandbuch.md` | update | Richtgröße im Abschnitt „Grenzwerte“ mit Ursprung; Diagnose-Beispiel; Änderungshistorie. |
| `harness/README.md` §Sensors | update | Zeile `make bench` (drei → vier Skripte) — aus dem realen Lauf geschrieben. |
| `tools/schema/schema.yaml` | prüfen (nicht ändern) | die View `backfill_status` trägt die zwei Warn-Spalten seit `slice-backfill-sql-administration`; ein Bedarf an einer Signaturänderung wäre ein Plan-Nachzug und ein Rückführungsgrund (§4), keine stille Änderung. |
| `docs/user/bench-abdeckung.md` | prüfen | die Datei bindet `LH-QA-PER-001`…`003` an durchgesetzte Schwellen; eine Zeile ohne Schwelle für `LH-FA-CAP-009` passt nicht zu dieser Form — Entscheidung im Slice, im Bericht begründet. |

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
| Zeile `make bench` in `harness/README.md` („drei eigenständige Bench-Skripte") | `grep -rn 'drei' harness/README.md docs/user/bench-abdeckung.md Makefile` (Fundstellen nach Bench-Bezug lesen) | *(Implementer trägt ein)* | nachziehen |
| Hilfetext und Kommentar des Targets `bench` | Lesen von `Makefile` | *(Implementer trägt ein)* | nachziehen |
| Bench-Abdeckungs-Träger | Lesen von `docs/user/bench-abdeckung.md` und dem Generator in `tools/bench-batch-vs-single.sh` | *(Implementer trägt ein)* | Entscheidung siehe §3, Zeile `bench-abdeckung.md` |
| Handbuch-Grenzwerte | Lesen von §Grenzwerte in `docs/user/benutzerhandbuch.md` | *(Implementer trägt ein)* | Richtgröße ergänzen; die Replication-Verbindungs-Zahl ist Sache von `sql-administration` |
| Toleranz-Konstante an genau einer Stelle; Wortwahl der Richtgröße | `grep -rn 'Toleranz\|Richtgröße\|Orientierung' internal docs harness --include=*.go --include=*.md` (Konstante, Handbuch, Kommentare) | *(Implementer trägt ein)* | eine Definition der Toleranz; „geschätzt" an jeder Nennung der Zeilenzahl; kein „Grenze"/„Limit"/„maximal" am Begriff |
| Build-Kontext | Lesen von `.dockerignore` und `Dockerfile` | *(Implementer trägt ein)* | Bench-Skripte laufen auf dem Host; keine Docker-Stufe kopiert `tools/bench-*.sh` — sonst Freigabe nach [`ADR-0085`](../../adr/0085-build-kontext-ausnahme-test-only-zweck.md) |

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

DoD vollständig + `make gates` grün + ein realer `make bench`-Lauf + `make test`
grün + Closure-Notiz mit Lerneintrag geschrieben.

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
