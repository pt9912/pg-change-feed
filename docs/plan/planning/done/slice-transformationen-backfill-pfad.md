# Slice transformationen-backfill-pfad: Backfill-Pfad — die Regelauswertung im Run, Fail-closed um den Regelstand, Beleg am laufenden System

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-transformationen](../welle-transformationen.md).

**Bezug:** [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) (die ausgelieferte
Change trägt die transformierte Form, nicht die Rohform),
[`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Backfill des Bestands),
[`LH-QA-SEC-004`](../../../../spec/lastenheft.md) (Wert nirgends im Image),
[`LH-FA-CAP-008`](../../../../spec/lastenheft.md),
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Folgepflicht 7 (Bindung künftiger Erzeugungspfade),
[`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 2
(eine Funktion für WAL- und Backfill-Pfad) und Teilfrage 4 (Fail-closed vor dem
Commit),
[`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md) (Klasse
`schema` für eine im Run nicht anwendbare Regel, Folgepflichten 2 und 3).

**Berührte Spec-Stellen:** [`SPEC-008`](../../../../spec/pflichtenheft.md)
(Fehlerklassen), [`SPEC-029`](../../../../spec/pflichtenheft.md) (Run-Zustand,
durch `slice-backfill-spec-nachzug`) — gelesen; die Zeile `schema` von
[`SPEC-008`](../../../../spec/pflichtenheft.md) und der Satz zum Run in
[`LH-FA-CAP-009.a`](../../../../spec/pflichtenheft.md) sind Gegenstand von
`slice-transformationen-spec-nachzug` (Folgepflicht 1 von
[`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md)), nicht dieses
Plans.

**Verantwortlich:** Implementer-Agent, 2026-09-26.

**Autor:** Planner-Agent, Welle-Eröffnung
[welle-transformationen](../welle-transformationen.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Backfill-Changes tragen dieselbe transformierte Form wie WAL-Changes:
der Run liest den Regelstand (Port aus `antragsweg-usecase`) je Block neu und
baut das Bild über die gemeinsame Row-Image-Funktion **mit** dem Regelsatz; die
Fail-closed-Prüfung vor dem Commit prüft zusätzlich, dass der Regelstand
demselben Stand entspricht, mit dem die Blöcke gebaut wurden; eine auf den
Bestand nicht anwendbare Regel endet den Run sichtbar; ein E2E-Beleg zeigt die
transformierte Form eines Backfill-Bestands am laufenden System. Das ist
Kopplung K2 der Welle [welle-backfill-bestand](../done/welle-backfill-bestand.md)
§5.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Änderung am WAL-Pfad, am `Assembler` oder am Antragsweg** — `kern-rename`,
  `antragsweg-*`; dieser Slice ruft die dort gelieferten Bausteine und ergänzt
  den Run.
- **Der zweite Regeltyp** — `map-value`; der Paritäts- und der Eigenschaftstest
  dieses Slice zählen die Regeltypen aus der Domänen-Menge auf und erfassen
  `map_value`, sobald er dort steht.
- **Ein zweiter Bild-Bau-Weg im Run** — der Run hat **eine** Stelle, an der
  Ausschluss- und Regelstand in das Bild eingehen (Welle
  [welle-backfill-bestand](../done/welle-backfill-bestand.md) §5 K2); ein zweiter
  Weg wäre ein zweiter Träger derselben Aussage.
- **Checkpoint, Wiederaufnahme, Live-Zustellung von Backfill-Changes** —
  Out-of-Scope der Welle
  [welle-backfill-bestand](../done/welle-backfill-bestand.md).
- **Eine Regelform, die nur der Backfill kennt** — Regeln gelten je Tabelle für
  beide Erzeugungspfade; eine Backfill-spezifische Regel ist keine Fähigkeit
  von
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md).

## 2. Definition of Done

- [x] Regelauswertung im Run: jeder Block liest den Regelstand über den Port
      neu (neben dem Ausschlussstand) und baut das Bild über die gemeinsame
      Funktion mit dem Regelsatz; ein Backfill-Change trägt bei gleicher Zeile
      und gleicher Regelmenge ein byte-gleiches Bild wie die WAL-Change am
      Ausgang der gemeinsamen Bild-Konstruktion (Byte-Gleichheit ist die Aussage von
      [`ADR-0115`](../../adr/0115-backfill-spaltenwerte-text-ergebnisformat.md)
      Festlegung 4 an dieser Stelle; die Spec sagt über den `jsonb`-Lesepfad
      Schlüsselmenge und Werte zu:
      [`LH-FA-CAP-009.a`](../../../../spec/pflichtenheft.md) „Markierung“;
      Paritätstest, tabellengetrieben über die Regeltypen der Domäne); ein
      Backfill-Aufrufer der gemeinsamen Funktion mit leerer Regelmenge, den
      `kern-rename` hinterlassen hat, ist ersetzt. *Zu belegen durch:* `make
      test` (Race-Detector) und — wenn der Paritätstest der Backfill-Welle im
      Replication-Tier liegt (am Start gelesen) — `make test-replication`.
- [x] Fail-closed und Nichtanwendbarkeit: die Prüfung vor dem Commit vergleicht
      zusätzlich den Regelstand (Mengengleichheit unabhängig von der
      Reihenfolge; eine Zwischenabweichung — Regel gesetzt, dann entfernt —
      wird erkannt); eine Regel, die auf den Bestand nicht anwendbar ist
      (Spalte fehlt in der Spaltenliste des Snapshots, Zielname kollidiert),
      endet den Run `failed` mit der Klasse `schema` — einmal je Run, nachdem
      der Snapshot seine Spalten liefert und bevor die Schreibtransaktion
      öffnet, mit derselben Prüffunktion der Domäne wie der Erfassungspfad
      ([`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md)
      Festlegung 1 und 2); `classifyError` bildet den Sentinel der Prüfung auf
      `schema` ab (der Kommentar „vergibt `schema` nicht“ entfällt), ein
      Wechsel des Regelstands zwischen den Lesungen endet mit `configuration`
      (Festlegung 5); der Fehler ist run-lokal, ohne Heartbeat-Fehlerzustand
      und ohne den Capture-Pfad zu stoppen
      ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md)
      Teilfrage 5, Festlegung 3 von `ADR-0117`); ein Eigenschaftstest im Run
      (Regeltyp × ausgeschlossene Spalte) belegt
      [`LH-QA-SEC-004`](../../../../spec/lastenheft.md) für den Backfill-Pfad.
      *Zu belegen durch:* `make test` gegen Fakes, je
      Negativfall an seine Eingabe gebunden (Mutation der Prüfung färbt den
      Test rot); dazu der Lesefehler des Regelstands je Block und unmittelbar
      vor dem Commit — der Fake scheitert ab dem n-ten Aufruf, je Aufrufstelle
      eine Mutation, die ihren Fehler verwirft
      (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`).
- [x] E2E-Beleg in `make test-integration`: für eine Tabelle mit
      `rename_column`-Regel trägt ein Backfill-Run den Bestand über
      `cdc.changes` und `GET /changes` mit umbenanntem Schlüssel und `origin =
      'backfill'`; mit zusätzlichem `exclude_column` auf der umbenannten Spalte
      trägt kein Backfill-Image Quellnamen, Zielnamen oder Wert
      ([`LH-QA-SEC-004`](../../../../spec/lastenheft.md)); eine im Run nicht
      anwendbare Regel endet den Run `failed`/`schema` ohne Change, der
      Erfassungspfad läuft weiter, und nach der Abhilfe (Regel entfernen, neuer
      Antrag) endet ein neuer Run `completed`
      ([`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md)
      Folgepflicht 3; Festlegung 4 führt „ohne Prozessneustart“ als *erwartet*).
      *Zu belegen durch:*
      ein realer, grüner `make test-integration`-Lauf am laufenden
      Feed-Container (Zeile im Runner-Erzeugnis
      [`docs/user/e2e-abdeckung.md`](../../../user/e2e-abdeckung.md)).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
      Report: `docs/reviews/review-slice-transformationen-backfill-pfad.md` (Reviewer, ohne Fixrunde).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: entfällt — keine neue Betreiber-Oberfläche; die Aussage,
      dass Regeln auch für einen Backfill gelten, steht mit den übrigen im
      Handbuch-Abschnitt von `slice-transformationen-betriebsdoku` (Adresse in
      dessen §2).
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
      der Closure dieses Slice (§7) und zusätzlich von der Closure der Welle
      [welle-transformationen](../welle-transformationen.md) (die Roadmap führt
      sie unter *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/application/usecase/backfill/service.go` | update | `Ports` trägt den Regelstand-Port (`Transformations`, neun Ports); `copyBlocks` liest den Regelstand einmal nach dem Öffnen des Snapshots (Stand des Runs), prüft ihn mit `checkRulesApplicable` (`Transformation.CheckApplicable` gegen `snapshot.Columns()`) vor dem ersten `NextBlock` und vor `Begin`, liest ihn je Block und unmittelbar vor dem Commit neu und vergleicht als Menge (`sameSet`, ersetzt `sameNames`/`nameSet`); `blockBuilder.build` übergibt den Regelsatz des Blocks an `BuildRowImage` (der Kommentar mit der Adresse dieses Slice ist entfernt); `classifyError` bildet `ErrTransformationColumnMissing`/`ErrTransformationTargetCollides` auf `schema` und `ErrTransformationStateChanged` auf `configuration` ab. Kosten der Lesung je Block: als benannte Grenze im Doc-Kommentar von `copyBlocks` (der Port liest die Zeilen aller Tabellen der Quelle). |
| `internal/domain/errors/errors.go` | update | neuer Sentinel `ErrTransformationStateChanged` (Klasse `configuration`): [`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md) Festlegung 5 verlangt die Klasse für den Wechsel des Regelstands; ein eigener Sentinel lässt den Fehlertext den Regelstand nennen statt den Ausschluss (Plan-Drift: dieser Plan nannte keinen neuen Sentinel). |
| `internal/domain/model/rowimage.go` | update (Kommentar) | der Satz „der Backfill-Lauf mit leerer Regelmenge“ beschreibt den Ist-Zustand (Übergabe aus `kern-rename`). |
| `internal/bootstrap/wiring.go` | update | `Transformations: activation` in den `backfill.Ports` (dieselbe Adapter-Instanz wie `Exclusion`, dieselbe Rolle auf `cdc.administration_request`). |
| `internal/bootstrap/administration_roles_internal_test.go`, `internal/bootstrap/backfill_endtoend_test.go` | update | die zwei Test-Verdrahtungen mit realem Adapter tragen den neuen Port. |
| `internal/application/usecase/backfill/service_test.go`, `internal/application/usecase/backfill/transformation_test.go` (neu) | update, neu | `fakeRules` (Fehler ab dem n-ten Aufruf wie `fakeExclusion`) im Rig; die Tests der Regelauswertung, des Eigenschaftstests (Regeltyp × Regel-Spalte × ausgeschlossene Spalte), der Nichtanwendbarkeit (`schema`), des Zustandswechsels (`configuration`), des Lesefehlers je Aufrufstelle und der nicht lesbaren `applied`-Zeile (`internal`). |
| `internal/bootstrap/backfill_image_parity_test.go` (neu) | neu | Paritätstest tabellengetrieben über `model.TransformationKinds()`: echter `mapper.Assembler` gegen echten `BackfillTableService` (Fakes an den Ports des Runs), Bytes der Bilder gleich. Der Ort am Start gelesen: `TestImageParityWalAndBackfill` (`internal/adapters/driven/postgressnapshot/snapshot_test.go`, Replication-Tier) trägt die Typ-Parität über reale Textwerte, die regelunabhängig ist (eine Regel ändert Schlüssel, nie den Wert) und unverändert bleibt; der Regel-Fall liegt in `make test`, weil beide Pfade nur die Composition Root zugleich importieren darf (`.a-check.yml`, `composition_root`). |
| `tools/harness/run-integration-tests.sh` | update | neue Phase „Backfill-Regelstand“ (nach „Backfill-Boundary“): `rename_column`-Bestand über `cdc.changes` und `GET /changes`, Schlüsselmengen-Vergleich mit einer WAL-Change, `exclude_column` auf der umbenannten Spalte, nicht anwendbare Regel (`schema`, run-lokal) und Abhilfe. Die Phase ist eine Shell-Phase (SQL, HTTP): `test/integration/integration_test.go` bleibt unverändert, es entsteht kein `func TestE2E*` und kein `-run`-Muster (Nicht-Realisierung der Plan-Zeile, Grund: kein Go-Test nötig; `BEO-PGC/test-runner-stiller-ausschluss` betrifft die `-run`-Muster und trifft die Phase nicht). |
| `harness/README.md` | update | §3.13-Träger: die Zeile `make test-integration` zählt die Backfill-Rundläufe (sieben → acht) und nennt die Phase. |
| `docs/user/e2e-abdeckung.md` | Erzeugnis | kommt aus dem Runner, wird nicht von Hand geschrieben. |

**Übergaben aus `slice-backfill-run-usecase`** (gemeldet, kein zusätzlicher Umfang; alle
Stellen in `internal/application/usecase/backfill/service.go`):

- **Stelle des Bild-Baus.** `blockBuilder.build` baut das Bild je Zeile über
  `model.BuildRowImage(columns, row, excluded, rules)` mit `rules []model.Transformation`
  (der Aufruf im Run übergibt `nil`); das ist die einzige Stelle des Runs, an der ein Regelsatz
  eingeht.
- **Stelle der Fail-closed-Prüfung.** In `copyBlocks` liest `excludedColumns` je Block den
  Ausschlussstand neu und `sameNames` vergleicht ihn mit dem Stand des ersten Blocks; vor dem
  Commit prüfen `stillBound` die Bindung und derselbe Vergleich den Stand; ein nicht lesbarer
  Stand endet den Run wie eine Abweichung. Die Grenze der Prüfung (ein zwischen zwei
  Lesungen gesetzter und zurückgenommener Stand ist unsichtbar, der Stand trägt keine
  Historie) steht im Doc-Kommentar von `copyBlocks` und gilt für den Regelstand ebenso.
- **Port und Fake des Regelstands.** `Ports` bündelt die Pflicht-Ports des Use Cases; der
  Regelstand-Port kommt dort hinzu. Die Fakes des Use-Case-Tests lassen einen Aufruf ab dem
  n-ten scheitern (`fakeExclusion.errCall`); der Fake des Regelstands folgt diesem Muster.
- **Klassifikation.** `classifyError` (Abbildung des Sentinels auf `schema`:
  [`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md) Festlegung 1).

**Übergaben aus `slice-transformationen-kern-rename`** (gemeldet, kein zusätzlicher Umfang):

- **Vorbedingung `CheckApplicable`.** `BuildRowImage` prüft die Anwendbarkeit einer Regel nicht;
  die Vorbedingung des Aufrufers ist `model.Transformation.CheckApplicable(columns)` je Regel,
  hier gegen die Spalten des Snapshots — das ist die Prüffunktion der Domäne, die
  [`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md) Festlegung 2 meint (einmal je
  Run, vor der Schreibtransaktion). Eine Regel, deren Spalte nicht in `columns` steht, wirkt in
  `BuildRowImage` nicht.
- **Fehlerverhalten aus dem Bild-Bau.** `BuildRowImage` liefert `ErrTransformationTargetCollides`
  und kein Bild, wenn ein Zielschlüssel schon vergeben ist (eine Spalte aus `columns` oder ein
  anderer im Bild umbenannter Schlüssel); der Fall hängt an der Zeile (beide Quellwerte tragen
  einen Wert) und tritt daher erst im Block auf, nach dem Öffnen der Schreibtransaktion. Der
  Erfassungspfad ordnet den Fehler in `mapper` als `ErrTransformationNotApplicable` ein, der Run
  kann den Mapper-Sentinel nicht importieren: `classifyError` bildet die Domänen-Sentinels
  `ErrTransformationColumnMissing` und `ErrTransformationTargetCollides` auf `schema` ab, und der
  Run endet in beiden Fällen atomar `failed` ([`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md)
  Festlegung 1).
- **Zeitgebundene Kommentare.** Der Kommentar an `blockBuilder.build` (`service.go`, nennt die
  Adresse dieses Slice) und ein Satz im Doc-Kommentar von `BuildRowImage` (`rowimage.go`, „der
  Backfill-Lauf mit leerer Regelmenge“) beschreiben den Stand ohne Regelstand; dieser Slice
  schreibt beide um.
- **Überholte Signatur in einer `Accepted`-ADR.** [`ADR-0115`](../../adr/0115-backfill-spaltenwerte-text-ergebnisformat.md)
  führt die Signatur von `BuildRowImage` mit drei Parametern und den Satz „bleibt
  unverändert“ (Konstraints, Festlegung 3); der Satz beschreibt den Stand dieser Entscheidung,
  die Erweiterung um den Regelsatz trägt
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Folgepflicht 7. Kein Nachzug in der ADR (`AGENTS.md` §3.5), kein Folge-ADR.

**Übergabe aus `slice-transformationen-antragsweg-usecase`** (gemeldet, kein zusätzlicher
Umfang; Herkunft: Review-Report `review-slice-transformationen-antragsweg-usecase` Finding F-11
und Plan des Slice §6, gelesen am Stand `2c22334f`):

- **Die Lesung des Regelstands je Block liest die ganze Quelle.** Der Regelstand-Port
  (`TransformationPort.TransformationRules(ctx, source)`, `internal/application/port/outbound/transformation.go`)
  liefert den Stand **aller** Tabellen einer Quelle: das Statement
  (`SelectAppliedTransformationRequests`, `queries.go`) liest je Aufruf alle `applied`-Zeilen der
  zwei Transformations-Antragsarten mit `WHERE source_id = $1` — ohne Tabellenfilter — und faltet
  sie je Tabelle in `model.FoldTransformations`. Ein Aufruf je Block (Ziel dieses Slice: „je Block
  neu“, §1) ist damit je Block eine Lesung über die Zeilen aller Tabellen der Quelle, nicht nur
  der des Runs; der Ausschlussstand (`ColumnExclusionPort.ExcludedColumns(ctx, source)`, gelesen
  je Block in `copyBlocks`) hat dieselbe Form. Der Slice misst die Kosten (Zeilenzahl der Queue
  einer Quelle mal Blockzahl eines Runs) oder trägt sie als benannte Grenze; ein
  tabellenbezogener Lesezugriff änderte den Port und gehört dann in seinen Plan (Kosten,
  Rückführung §4). Dieser Plan trägt die Kosten der Lesung nicht (§1 nennt „je Block neu“, §6
  führt kein Kostenrisiko).
- **Fehler der Faltung im Run.** Eine `applied`-Zeile, die die Faltung nicht mehr in eine Regel
  führt, endet als Fehler der Klasse `internal` (nicht `storage`) und hält jede Lesung der Quelle
  an (Prozessstart und Regel-Anträge aller Tabellen der Quelle; bewusst, „der Stand wird nie um
  eine Zeile verkürzt“); ein Run, der den Regelstand je Block liest, endet an derselben Stelle.
  Der Slice nennt die Klasse, in der `classifyError` diesen Fehler abbildet, und bindet den Fall
  mit einem Test (Eingabeseite: eine nicht lesbare `applied`-Zeile im Fake des Regelstands).

**Festlegungen des Implementers (Semantik der Regelstand-Lesung im Run):**

- **Der Stand des Runs ist der Stand nach dem Öffnen des Snapshots.** Die Lesung zu Beginn (einmal je Run, vor dem ersten `NextBlock`, vor `Begin`) ist Grundlage der Anwendbarkeitsprüfung und der Vergleichsstand; eine Regel, die zwischen Antrag und Start des Runs gesetzt wurde, gilt für den Run. Ein Antrag während des Runs endet ihn mit `configuration`, sobald eine der folgenden Lesungen (je Block, unmittelbar vor dem Commit) ihn sieht ([`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md) Festlegung 5).
- **Kosten der Lesung (benannte Grenze, gemessen nicht).** Jede Lesung liest die `applied`-Zeilen der zwei Transformations-Antragsarten aller Tabellen der Quelle (Statement `SelectAppliedTransformationRequests`, ohne Tabellenfilter); ein Run liest den Regelstand `Blockzahl + 2`-mal, den Ausschlussstand `Blockzahl + 1`-mal (`DefaultBlockSize` 1000, `internal/adapters/driven/postgressnapshot/snapshot.go`). Die Kosten eines Runs sind damit Zeilen der Queue der Quelle mal `Blockzahl + 2` — abgeleitet, nicht gemessen; ein tabellenbezogener Lesezugriff änderte den Port (`TransformationPort`, Adapter, Fakes) und gehört in einen eigenen Plan, falls die Queue einer Quelle in die Größenordnung der Blockzahl wächst. Die Zahl der Lesungen steht im Doc-Kommentar von `copyBlocks`, die Kennzeichnung „abgeleitet, nicht gemessen“ steht hier, nicht im Kommentar; die Formel gilt ab einem Block (eine leere Tabelle endet ohne Schreibtransaktion: der Regelstand wird einmal, der Ausschluss nie gelesen). Adresse der Grenze: §6, Risiko „Kosten der Lesung je Block“.
- **Klasse der Faltungsfehler im Run:** `internal` (Rückfall von `classifyError`); der Fall ist mit einem Test gebunden (`TestExecuteUnreadableAppliedRowEndsRunAsInternal`, Eingabeseite: der reale Fehler von `model.FoldTransformations` an einer nicht lesbaren Regelform, je Aufrufstelle der Lesung).
- **Leere Tabelle:** die Prüfung der Anwendbarkeit läuft auch dann (sie hängt an Regel- und Spaltenmenge, nicht an einer Zeile, [`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md) Festlegung 2): ein Run über eine leere Tabelle mit nicht anwendbarer Regel endet `failed`/`schema`.

**Mutationen (je Zusage eine Eingabeseiten-Mutation; Zusage · mutierte Eingabe · gesehenes Rot; die Menge der Stellen steht in der Zusage):**

| Zusage | mutierte Eingabe | gesehenes Rot |
|---|---|---|
| Der Run baut das Bild mit dem Regelsatz des Blocks (Regeltypen der Domäne: `rename_column`, drei Blöcke, drei Zeilen je Fall) | `nil` statt `rules` an `BuildRowImage` (`build`) | `TestExecuteBuildsImagesWithTheRuleSet`, `TestExecuteRulesNeverLeakExcludedColumns`, `TestExecuteRuleTargetCollisionInBlockEndsRunAsSchema`, `TestBackfillAndWALImagesAreByteEqualWithRules` |
| Der Run liest den Regelstand der Tabelle des Runs (Schlüssel `schema.table`) | Schlüssel einer anderen Tabelle (`public.other`) in `transformationRules` | `TestExecuteBuildsImagesWithTheRuleSet`, `…InapplicableRule…`, `…StateChange…`, `…NeverLeak…`, `…FailureIsReportedInTheRunOnly`, Parität (sieben Tests) |
| Der Run liest den Regelstand der Quelle des Runs | leere Quelle an `TransformationRules` | `TestExecuteBuildsImagesWithTheRuleSet` |
| Die Parität von WAL- und Backfill-Bild (Regeltypen × Regel-Spalte × Ausschluss, 12 Fälle) | `nil` statt `binding.Transformations` an `BuildRowImage` im `Assembler` (Gegenseite) | `TestBackfillAndWALImagesAreByteEqualWithRules` |
| Ausschluss vor Regel im Run (`LH-QA-SEC-004`; 9 Fälle Regeltyp × Regel-Spalte × ausgeschlossene Spalte) | Ausschluss erst am Zielschlüssel in `BuildRowImage` prüfen | `TestExecuteRulesNeverLeakExcludedColumns` (dazu die Tests von `model` und `mapper`), `TestBackfillAndWALImagesAreByteEqualWithRules` |
| Nichtanwendbarkeit endet den Run `schema` vor dem ersten Block (6 Fälle: Spalte fehlt, Ziel kollidiert, leere Tabelle, zweite Regel, fremde Tabelle, anwendbare Regel) | Prüfung entfällt; Prüfung gegen `nil` statt der Snapshot-Spalten | `TestExecuteInapplicableRuleEndsRunAsSchema` (beide) |
| Die Klasse `schema` für beide Sentinels der Prüfung | Abbildung von `ErrTransformationTargetCollides` entfernt; Abbildung von `ErrTransformationColumnMissing` entfernt (je eine) | `TestExecuteInapplicableRuleEndsRunAsSchema` (beide), `TestExecuteRuleTargetCollisionInBlockEndsRunAsSchema` (erste) |
| Zustandswechsel des Regelstands endet `configuration` (11 Fälle; drei Lesestellen: Block, vor dem Commit; Vergleich als Menge) | Vergleich je Block entfernt; Vergleich vor dem Commit entfernt; Mengenvergleich als Längenvergleich (je eine) | `TestExecuteRuleStateChangeEndsRunAsConfiguration` (drei, zwei, zwei Fälle) |
| Abbildung `ErrTransformationStateChanged` → `configuration` | Abbildung entfernt (Klasse `internal`) | `TestExecuteRuleStateChangeEndsRunAsConfiguration` (alle Fehlerfälle) |
| Lesefehler des Regelstands endet den Run (fünf Aufrufstellen: zu Beginn, drei Blöcke, vor dem Commit) | Fehler verworfen, je Stelle eine Mutation (zu Beginn, je Block, vor dem Commit) | `TestExecuteRuleReadFailure` (der Fall der Stelle; je Mutation weitere Fälle der Faltungs-Klasse), `TestExecuteUnreadableAppliedRowEndsRunAsInternal` |
| Eine nicht lesbare `applied`-Zeile endet als `internal` | Rückfall von `classifyError` auf `storage` | `TestExecuteUnreadableAppliedRowEndsRunAsInternal` (vier Stellen) |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „welche
Erzeugungspfade `model.Change`-Bilder bauen und ob sie den Regelstand tragen“;
beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Alle Bild-Erzeuger | `git grep -n 'BuildRowImage('` und `git grep -n 'json.Marshal'` über `internal` ohne Tests (Block unten, Zeilen 1–4) | **Gefunden.** Parent `3973390e`: `BuildRowImage(` 4 Treffer — die Definition (`rowimage.go`), zwei Aufrufe in `mapper.go` (mit `binding.Transformations`) und ein Aufruf in `service.go` (mit `nil`); `json.Marshal` 4 Treffer — zwei in `rowimage.go` (die Konstruktion selbst), je einer in `natsstream/publisher.go` und `http/sse.go` (Serialisierung eines fertigen Bildes für die Zustellwege, kein Bild-Bau). Diff: 4 und 4, der Aufruf in `service.go` trägt den Regelsatz des Blocks. **Nichtgefunden:** keine zweite Konstruktion eines Row Images, kein zweiter Bild-Bau-Weg im Run. | keine Änderung an den Fundstellen außer dem Aufruf im Run; die zwei Zustellweg-Treffer sind Serialisierung, kein Bild-Bau. |
| Aufrufer mit leerer Regelmenge | `git grep -n -E 'BuildRowImage\(.*nil\)'` über `internal` ohne Tests (Block unten, Zeilen 5–6) und mit Tests (`-- 'internal/*_test.go'`, gelesen) | **Gefunden.** Parent 1 Treffer (`service.go`), Diff 0. Mit Tests: 9 Treffer, alle mit `nil` als Regelsatz — `rowimage_test.go` (sieben, der Vergleich mit leerer Regelmenge) und `postgressnapshot/snapshot_test.go` (zwei, `TestImageParityWalAndBackfill`: die reale Typ-Parität der Textwerte, regelunabhängig, siehe §3-Zeile des Paritätstests). **Nichtgefunden:** kein Nicht-Test-Aufrufer mit leerer Regelmenge im Backfill-Pfad mehr. | die zwei Aufrufer in `snapshot_test.go` bleiben: sie beweisen die Gleichheit der Rohtexte, und eine Regel ändert nur Schlüssel. |
| Fail-closed-Aufzählungen (Bindung, Ausschlussstand) in Doc-Kommentaren und Docs | `git grep -n Ausschlussstand` über `internal docs spec harness` ohne `docs/reviews`, `done/`, Baseline (Block unten, Zeilen 7–8) und `git grep -n sameNames` über `internal` (Zeilen 23–24) | **Gefunden.** Parent 92 Treffer, Diff 93 (ein Treffer mehr: der Doc-Kommentar von `copyBlocks`). Die Träger, die die Prüfung vor dem Commit beschreiben: `service.go` (Doc-Kommentare von `Ports` und `copyBlocks`, nachgezogen), `spec/pflichtenheft.md` §Fail-closed (trägt „Ausschlussstand sowie der Regelstand“ seit `slice-transformationen-spec-nachzug`), `spec/architecture.md` (Sequenzdiagramm „Bindung, Ausschluss- und Regelstand erneut prüfen“). `sameNames` (Symbol): 4 Treffer im Parent, alle `service.go`, 0 im Diff (ersetzt durch `sameSet`). **Nichtgefunden:** keine Aufzählung in `docs/user` (`benutzerhandbuch.md` trägt den Ausschlussstand nur als Dauerhaftigkeit der Spalten-Anträge, ohne Backfill-Bezug). | Aufzählungen in `service.go` tragen den Regelstand; Spec und Sicht tragen ihn schon. |
| Beschreibung des Backfill-Bild-Baus in Doku | `git grep -n -i Backfill` über `docs/user harness spec` (Block unten, Zeilen 9–10) und die Zählwörter „sieben Backfill“/„acht Backfill“ (Zeilen 11–14) | **Gefunden.** 236 Treffer im Parent, 237 im Diff (einer mehr: die Zeile der neuen Phase im Erzeugnis `docs/user/e2e-abdeckung.md`); die Aussagen über die Bildform stehen in `spec/pflichtenheft.md` §Fail-closed und §Sichtbarkeit, in `spec/architecture.md` (Form der Row Images, Sequenz) und in `harness/README.md`; Zählwort: „sieben Backfill“ 1 Treffer im Parent (`harness/README.md`, Zeile `make test-integration`), 0 im Diff; „acht Backfill“ 0 im Parent, 1 im Diff. **Nichtgefunden:** keine Handbuch-Stelle, die die Bildform des Backfills beschreibt (das Handbuch trägt die Transformationen noch nicht: `slice-transformationen-betriebsdoku`). | die Zeile `make test-integration` in `harness/README.md` zählt acht Rundläufe und nennt die Regelstand-Phase; Meldung an `slice-transformationen-betriebsdoku`: die Aussage, dass Regeln auch für einen Backfill gelten. |
| Adresse dieses Slice im Code | `git grep -n 'slice-transformationen-backfill-pfad'` über `internal tools harness spec docs/user Makefile` (Block unten, Zeilen 15–16) | **Gefunden.** Parent 1 Treffer (Kommentar in `blockBuilder.build`), Diff 0. **Nichtgefunden:** keine weitere Adresse auf diesen Slice außerhalb der Pläne. | der Kommentar beschreibt den Ist-Zustand ohne Adresse. |
| Klassifikations-Kommentar „vergibt `schema` nicht“ | `git grep -n -F 'vergibt der Run nicht'` über `internal` (Block unten, Zeilen 17–18) | **Gefunden.** Parent 2 Treffer (`classifyError`, Doc-Kommentar von `TestExecuteClassifiesFailures`), Diff 0. | beide umformuliert (die Klasse `schema` und der Zustandswechsel stehen im Doc-Kommentar von `classifyError`). |
| Aufrufer des Regelstand-Ports und der Prüffunktion | `git grep -n TransformationRules` und `git grep -n -E 'CheckApplicable\('` über `internal` ohne Tests (Block unten, Zeilen 19–22) | **Gefunden.** `TransformationRules`: Parent 13, Diff 15 (zwei mehr: die Lesung im Run und die Nennung im Doc-Kommentar von `copyBlocks`); `CheckApplicable(`: Parent 2 (Definition und `mapper.go`), Diff 3 (der Run ruft dieselbe Prüffunktion der Domäne). **Nichtgefunden:** keine Kopie der Prüffunktion im Run (Risiko „Prüfung steht zweimal“). | der Run trägt `checkRulesApplicable` als Schleife über `Transformation.CheckApplicable`, keine eigene Regellogik. |

```suchlauf
3973390e 4 -n 'BuildRowImage(' -- internal ':!*_test.go'
diff 4 -n 'BuildRowImage(' -- internal ':!*_test.go'
3973390e 4 -n 'json.Marshal' -- internal ':!*_test.go'
diff 4 -n 'json.Marshal' -- internal ':!*_test.go'
3973390e 1 -n -E 'BuildRowImage\(.*nil\)' -- internal ':!*_test.go'
diff 0 -n -E 'BuildRowImage\(.*nil\)' -- internal ':!*_test.go'
3973390e 92 -n Ausschlussstand -- internal docs spec harness ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 93 -n Ausschlussstand -- internal docs spec harness ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
3973390e 236 -n -i Backfill -- docs/user harness spec
diff 237 -n -i Backfill -- docs/user harness spec
3973390e 1 -n 'sieben Backfill' -- harness tools docs/user spec
diff 0 -n 'sieben Backfill' -- harness tools docs/user spec
3973390e 0 -n 'acht Backfill' -- harness tools docs/user spec
diff 1 -n 'acht Backfill' -- harness tools docs/user spec
3973390e 1 -n 'slice-transformationen-backfill-pfad' -- internal tools harness spec docs/user Makefile
diff 0 -n 'slice-transformationen-backfill-pfad' -- internal tools harness spec docs/user Makefile
3973390e 2 -n -F 'vergibt der Run nicht' -- internal
diff 0 -n -F 'vergibt der Run nicht' -- internal
3973390e 13 -n TransformationRules -- internal ':!*_test.go'
diff 15 -n TransformationRules -- internal ':!*_test.go'
3973390e 2 -n -E 'CheckApplicable\(' -- internal ':!*_test.go'
diff 3 -n -E 'CheckApplicable\(' -- internal ':!*_test.go'
3973390e 4 -n sameNames -- internal
diff 0 -n sameNames -- internal
```

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-backfill-run-usecase` in
`done/` liegt (Kopplung K2 der Welle
[welle-backfill-bestand](../done/welle-backfill-bestand.md) §5), **zusätzlich**
`slice-backfill-e2e` (der E2E-Beleg braucht das lauffähige System und die
Backfill-Phase im Runner) und `slice-transformationen-antragsweg-usecase`
(Regelstand-Port und Wirkung) in `done/` liegen und kein anderer Slice in
`in-progress/` liegt (WIP-Limit 1). Das Architect-Kurzverdikt zur
Nichtanwendbarkeit einer Regel im Run liegt vor: das Verdikt
[`architect-verdict-backfill-schema-klasse-rollen`](../../../reviews/architect-verdict-backfill-schema-klasse-rollen.md)
und [`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md) (`Accepted`,
Klasse `schema` im Run, run-lokal, Folgepflicht 7 von
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
für den Run ausgefüllt). **Der Übergangs-Commit `next` → `in-progress` nennt
beide** (`BEO-PGC/start-trigger-ohne-uebergabe-artefakt`, offen, 1×). Die
Abbildung, an der `ADR-0117` ansetzt, steht in `classifyError` am Use Case des Runs
(`internal/application/usecase/backfill/service.go`): ein nicht erkannter Fehler endet als
`internal`, `schema` vergibt der Run bis zur Umsetzung dieses Slice nicht
(Register: `BEO-PGC/run-fehlerklasse-schema-im-transformations-backfill`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls
  Run-Erweiterung, Fail-closed und E2E-Beleg nicht in einem Review tragen — der
  abtrennbare Teil ist der E2E-Beleg (dritter Liefer-Punkt) als eigener Slice
  mit Start nach diesem.
- `in-progress` → `open` (blockiert): falls der Run-Use-Case der
  Backfill-Welle keine einzige Stelle für den Bild-Bau trägt (dann Plan-Nachzug
  an die Backfill-Welle statt einer Zweitkopie), oder falls die Prüffunktion der
  Regelanwendbarkeit in der Domäne keinen Ort trägt, den Erfassungspfad und Run
  gemeinsam rufen
  ([`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md)
  Festlegung 2; Architect-Frage zum Zuschnitt).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test` (Race-Detector) grün + ein
realer, grüner `make test-integration`-Lauf + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Die Prüfung der Regelanwendbarkeit steht zweimal** (eine Kopie für den Run
  neben der Domänen-Funktion des Erfassungspfads). *Erwartet, zu belegen
  durch:* der Suchlauf über die Aufrufer der Prüffunktion und Review
  ([`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md)
  Fitness Function, Zeile „Review-Prüfpflicht“). **Ausgang:** *entfallen* —
  der Prüfkern `Transformation.CheckApplicable` steht einmal in der Domäne und
  beide Pfade rufen ihn (Suchlauf `CheckApplicable\(` ohne Tests: Parent 2,
  Diff 3, ein Aufrufer je Pfad; Review-Finding F-5 nennt die Zeile
  „Review-Prüfpflicht“ erfüllt). Die Schleife über die Regeln samt Fehlertext
  steht je Pfad (`checkTransformations` im Mapper, `checkRulesApplicable` im
  Run): eine benannte Grenze ohne Register-Eintrag, weil beide Schleifen an
  ihrer Eingabeseite gebunden sind (Reviewer M4–M7, Verifier M-D rot).
- **Der Bild-Bau des Runs hat zwei Wege** (Ausschluss über Bild-Funktion, Regel
  über eine zweite Stelle). *Erwartet, zu belegen durch:* der Suchlauf (Zeile
  1) und Review. **Ausgang:** *entfallen* — `BuildRowImage(` ohne Tests: 4 am
  Parent und 4 im Diff, der eine Aufrufer im Run trägt den Regelsatz des
  Blocks; der Review-Negativbefund findet keinen zweiten Erzeuger eines Row
  Images; Ausschluss vor Regel im Run ist gebunden
  (`TestExecuteRulesNeverLeakExcludedColumns`, Reviewer M19 rot).
- **Fail-closed ist zu lasch für den Regelstand**: eine Regel, die zwischen
  zwei Blöcken gesetzt und wieder entfernt wird, hinterlässt Blöcke mit
  abweichender Form. *Erwartet, zu belegen durch:* der Negativtest der
  Zwischenabweichung mit Mutation. **Ausgang:** *entfallen* —
  `TestExecuteRuleStateChangeEndsRunAsConfiguration` bindet Vergleich je Block,
  Vergleich vor dem Commit und Mengenvergleich (Reviewer M8–M11, Verifier
  M-B und M-G rot). Grenze, im Doc-Kommentar von `copyBlocks` benannt: ein
  zwischen zwei Lesungen gesetzter und wieder entfernter Stand ist unsichtbar,
  weil der Stand keine Historie trägt; jeder Block entsteht mit dem Stand
  seiner eigenen Lesung, der dem Stand zu Beginn gleicht (Review, Bewertung
  der dritten Frage), kein Bild wird damit falsch gebaut.
- **Der Paritätstest sitzt in einem Tier, das dieser Slice nicht fährt**
  (DB-gestützt). *Erwartet, zu belegen durch:* Lesen des Ortes am Start; liegt
  er im Replication-Tier, gehört `make test-replication` zur Closure.
  **Ausgang:** *entfallen* — der Ort am Start gelesen: `TestImageParityWalAndBackfill`
  liegt im Replication-Tier und ist regelunabhängig (unverändert); der
  Regel-Fall `TestBackfillAndWALImagesAreByteEqualWithRules` liegt in
  `internal/bootstrap` in `make test`. `make test-replication` gehört nicht
  zur Closure und wurde nicht gefahren; der Diff berührt den Replication-Tier
  nicht.
- **Laufzeit von `make test-integration`** wächst mit der Backfill-Phase
  (`BEO-PGC/test-integration-retention-timing-flake`, verkörpert, 3×).
  *Erwartet, zu belegen durch:* ein realer Lauf; die Zahl trägt ihren Lauf
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A). **Ausgang:**
  *entfallen* — kein Zeitfehlschlag; gemessen 5 min 23 s (Reviewer, 11:11:19
  bis 11:16:42) und 5 min 14 s (Verifier, 11:33:38 bis 11:38:52, `make image`
  gegen den Layer-Cache eingerechnet), je ein Lauf am Stand `17d0048a` am
  2026-09-26 (übernommen aus den beiden Reports). Der Zuwachs durch die Phase
  ist nicht gemessen: kein Lauf des Parent-Stands.
- **Fenster zwischen Regel-Setzbarkeit und Backfill-Bindung** (Übergabe aus
  `slice-transformationen-antragsweg-usecase`, dessen Risiko mit seiner Closure hierher
  wandert). Seit der Closure von `antragsweg-usecase` sind Regeln setzbar; ein Backfill-Run,
  der vor der Closure dieses Slice läuft, liefert die Rohform, dieselbe Tabelle also zwei
  Formen ([welle-transformationen](../welle-transformationen.md) §5, Fenster ein Slice
  lang). *Erwartet, zu belegen durch:* der E2E-Beleg des dritten Liefer-Punkts (Backfill-Bestand
  einer Tabelle mit `rename_column`-Regel trägt den umbenannten Schlüssel) und die Benennung
  des Fensters im Bericht. **Ausgang:** *entfallen mit der Closure dieses Slice* — der
  E2E-Lauf trägt den Bestand mit umbenanntem Schlüssel `customer_name` (Werte
  `RegelAlpha,RegelBeta`, Schlüsselmenge gleich der WAL-Change, gedruckt im Lauf des
  Verifiers am Stand `17d0048a`); das Fenster bestand von der Closure von
  `slice-transformationen-antragsweg-usecase` bis zu dieser.
- **Kommentare zum Fehlerpfad des Runs behaupten mehr, als der Code trägt**
  (`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`, verkörpert, 5×).
  *Erwartet, zu belegen durch:* Review liest die Kommentare der neuen Zweige.
  **Ausgang:** *entfallen* — der Review fuhr die Zusagen der Kommentare am Code
  nach („auch bei einer leeren Tabelle“, „bevor eine Zeile gelesen und bevor
  die Schreibtransaktion geöffnet wird“, „Blockzahl plus zwei Lesungen“); kein
  Auftreten, der Zähler bleibt 5×. Der Fund V-1 der Verifikation (die
  Kennzeichnung „abgeleitet, nicht gemessen“ fehlt im Doc-Kommentar von
  `copyBlocks`) trifft keine Zusage des Kommentars: der Kommentar nennt die Zahl
  der Lesungen richtig; er ist das nächste Risiko.
- **Kosten der Lesung je Block** (Verifikation V-1, Review F-1; Übergabe aus
  `slice-transformationen-antragsweg-usecase`, dessen Risiko-Punkt „Kosten der
  Regelstand-Lesung je Block“ hierher wanderte). Jede Lesung des Regelstands und
  des Ausschlusses liest die `applied`-Zeilen der Transformations- und
  Spalten-Antragsarten **aller** Tabellen der Quelle (ohne Tabellenfilter,
  `SelectAppliedTransformationRequests`); ein Run liest den Regelstand
  `Blockzahl + 2`-mal, den Ausschluss `Blockzahl + 1`-mal (abgeleitet, am Code
  nachgezählt; der Test zählt fünf Lesungen des Regelstands bei drei Blöcken,
  vom Review nachgezählt). Nicht gemessen sind die Kosten je Lesung an einer
  Queue realer Größe; ein tabellenbezogener Lesezugriff änderte Port, Adapter und
  Fakes. *Erwartet, zu belegen durch:* ein Architect-Verdikt mit Messung an einer
  Queue realer Größe oder dem Folge-Slice des tabellenbezogenen Lesezugriffs.
  **Ausgang:** *weiter offen* — Adresse: das Closure-Kriterium „Die Lesekosten des
  Backfill-Runs sind entschieden“ in §3 der Welle
  [welle-transformationen](../welle-transformationen.md); der Auslöser der
  früheren Fassung („falls die Queue in die Größenordnung der Blockzahl wächst“)
  ist durch die Messung im Verdikt ersetzt, weil eine Größenordnung vorab nicht
  bestimmbar ist. Der Doc-Kommentar von `copyBlocks` trägt die Zahl der Lesungen
  im Indikativ und die Ausdehnung auf alle Tabellen der Quelle; die Kennzeichnung
  „abgeleitet“ steht in §3 dieses Plans (kein Code-Zug: die Verifikation sieht
  keinen Fixrunden-Zwang).

## 7. Closure-Notiz

- **Was hat funktioniert:** (1) Die Leser-Kette lief ohne Fixrunde: der Review nennt 0 HIGH, 0 MEDIUM, 0 LOW und 8 INFO, die Verifikation 0 HIGH, 0 MEDIUM, 3 LOW (V-1, V-2, V-4) und 4 INFO (übernommen aus den Reports). Der Review mutierte 21-mal an der Eingabeseite, 20 rot, die 21. — die Verdrahtungszeile `Transformations: activation` — grün; der Verifier wiederholte sieben Go-Mutationen (sechs rot, dieselbe Zeile grün) und fuhr die Zeile zusätzlich durch die E2E-Kette (`make image`, `make test-integration`: Exit 2 an der ersten Backfill-Phase, der Run bleibt `running`; übernommen aus dem Verifikations-Report, §4). (2) Die Mutationstabelle des Implementers (elf Zeilen, gezählt an §3) stimmte in ihren Zahlen mit den Läufen des Reviews überein (sieben Tests bei der Mutation des Tabellenschlüssels, neun Fälle bei der Abbildung auf `configuration`, vier Stellen bei der Klasse `internal`). (3) Das Suchlauf-Feld trug: `make suchlauf-nachmessen` meldete 24 stimmende Zeilen (Lauf des Verifiers, Exit 0), neun der zwölf Paare fuhr er von Hand nach, alle gleich. (4) Der reale Lauf trägt den Beleg des dritten Liefer-Punkts (gedruckt im Lauf des Verifiers am Stand `17d0048a`): ein Bestand von drei Zeilen mit `rename_column`-Regel liegt über `cdc.changes` und `GET /changes` mit dem Schlüssel `customer_name`; nach `exclude_column` auf `name` tragen vier Bilder weder `name` noch `customer_name` noch einen Wert von `name`; eine nicht anwendbare Regel endet den Run `failed` mit `schema` ohne Change, der Erfassungspfad läuft weiter; der neue Antrag nach dem Entfernen der Regel übernimmt zwei Zeilen in Rohform. Die Erwartung von [`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md) Festlegung 4 („ohne Prozessneustart“, dort *erwartet*) ist damit erprobt. (5) Die Übergabe-Blöcke der Vorgänger-Pläne (Stelle des Bild-Baus, Stelle der Fail-closed-Prüfung, Vorbedingung `CheckApplicable`, Klasse der Faltungsfehler) trafen am Code zu (Review, Negativbefund zum Plan).
- **Was ging anders als geplant:** (1) **Plan-Drift ohne Widerspruch:** der Plan nannte keinen neuen Sentinel; die Abbildung auf `configuration` trägt `ErrTransformationStateChanged`, während [`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md) Festlegung 5 den Wortlaut an `ErrExclusionStateChanged` bindet — die Klasse ist die der ADR, Review und Verifikation lasen keinen Widerspruch (`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab`, Zähler bleibt 2×, `state.md` nennt den Vorgang als kein Auftreten). (2) **Shell-Phase statt Go-Test** in `make test-integration` (die Plan-Zeile nennt die Nicht-Realisierung mit Grund); der Paritätstest liegt in `make test`, `make test-replication` lief nicht. (3) **Die Behauptung im Bericht des Implementers**, der Doc-Kommentar von `copyBlocks` trage die Kennzeichnung „abgeleitet, nicht gemessen“, traf nicht zu (Verifikation V-1: `git grep` nach `abgeleitet|gemessen` in der Datei ohne Treffer); die Kennzeichnung steht in §3 dieses Plans, der Kommentar nennt die Zahl der Lesungen und die Ausdehnung auf alle Tabellen der Quelle. Ein Lauf-Bericht ist kein Träger im Geltungsbereich von [`ADR-0083`](../../adr/0083-herkunft-von-aussagen-in-traegern.md) und der Kommentar sagt nichts Nicht-Getragenes: kein Register-Anfall. (4) **Die Register-Datei `state.md` des Eintrags zur Run-Klasse brach beim Move** (Link auf `open/`, Review F-7); der Implementer stellte sie auf die Kennung um. (5) **Umfang:** neun Commits bis zum Review-Report, 15 Dateien, +1330/−128 einschließlich der zwei Lifecycle-Moves (übernommen aus dem Verifikations-Report, Range `3973390e..17d0048a`). (6) **Die Coverage streut:** 85,10 % (Implementer und Verifier mit `make coverage-gate` allein) und 85,00 % (Reviewer in zwei Läufen, Verifier im Schritt von `make gates`), Schwelle 80 % — der Beleg je eines Laufs, übernommen aus den Reports, kein Ist-Stand. (7) **Der Handbuch-Träger nahm die Sendung nur dem Wort nach an** (Verifikation V-4): die Run-Abhilfe steht seit dieser Closure als committeter Text im Plan der Adresse. (8) **Die `diff`-Zeilen des Suchlaufs bewegten sich mit den Closure-Nachzügen nicht:** `make suchlauf-nachmessen` meldete am Arbeitsbaum der Closure 24 stimmende Zeilen (gemessen, Exit 0), obwohl die Closure Träger unter `docs/plan` ergänzt (die Sollwerte in §3 blieben unverändert).
- **Steering-Loop-Eintrag (Lerneintrag):** *(a) Neuer Sensor, gebaut und verkörpert.* Drei Träger binden, was vorher Zusage im Text war: `TestBackfillAndWALImagesAreByteEqualWithRules` (`internal/bootstrap`, zwölf Fälle Regeltyp × Regel-Spalte × Ausschluss, `make test`) bindet [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Folgepflicht 7 für den Backfill-Pfad — beide Erzeuger erhalten denselben Regelsatz; `TestExecuteRulesNeverLeakExcludedColumns` (Backfill-Paket, neun Fälle, `make test`) bindet [`LH-QA-SEC-004`](../../../../spec/lastenheft.md) im Run; die Phase „Backfill-Regelstand“ des Runners von `make test-integration` bindet die Form am laufenden System und trägt die Verdrahtung. Anker: `BEO-PGC/run-fehlerklasse-schema-im-transformations-backfill`, `state.md` (Beleg-Anker nachgezogen). Grenze: der Paritätstest prüft die Verdrahtung, nicht die Funktion — beide Seiten rechnet dieselbe `BuildRowImage`, eine gemeinsam falsche Funktion bliebe dort grün; sie tragen `rowimage_test.go` und der Eigenschaftstest im Run (Review, Negativbefund). *(b) Geschärfte Regel, Vorschlag an den Architect* (Verkörperung 3b, `v6.9.0` · `regelwerk/modul-08-agentenrollen.md` §Rollen-Sequenz für eine Welle, Schritt 3b; der Planner schärft `AGENTS.md` nicht selbst): die Nutzerregel „kein `sed -i`/`perl -pi`“ steht in keinem committeten Text (`git grep -n -i -E 'sed -i|perl -pi' -- .claude harness AGENTS.md .harness/skills`: kein Treffer) und wurde in diesem Vorgang von allen drei Rollen verletzt, jeweils ohne Wirkung auf das Repo (Ziele `/dev/null` und eine Scratch-Datei). `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` erreicht damit **3×** (`slice-backfill-speicher-untersuchung`, `slice-transformationen-antragsweg-usecase` und dieser Slice, je in einem Review-Report belegt); der Vorschlag steht in seiner `state.md`: `AGENTS.md` §3.1 nennt das in-place schreibende Host-Textwerkzeug als Verstoß und die Agenten-Prompts unter `.claude/agents/` verweisen darauf; ein Sensor ist ausgeschlossen (ein Werkzeugaufruf hinterlässt keine Signatur in der Datei). **Adresse: der Lese-Schritt der Closure von `welle-transformationen` und der nächste Architect-Zug, der `AGENTS.md` §3.1 berührt**; der Orchestrator beauftragt den Zug. Zusätzlich erreicht `BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` mit V-4 den vierten Beleg (die Adresse trug das Stichwort, nicht den Ablauf); sein Ausgangs-Vorschlag Punkt (2) hätte den Fall am Absender gefangen, Adresse unverändert der Lese-Schritt der Welle-Closure. *(c) Benannte Lücken, mit Adresse.* Erstens: die Zusage „ein zweiter Regeltyp ohne Änderung am Wirkort“ ([`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Folgepflicht 4, DoD Punkt 2 von `slice-transformationen-map-value`) rechnete den Run nicht mit: `sameSet[T comparable]` koppelt ihn an die Vergleichbarkeit von `model.Transformation`, und zwei Fixture-Schalter der Regeltypen liegen in Verzeichnissen, die DoD Punkt 2 ausnimmt; Adresse: der Block „Übergabe aus `slice-transformationen-backfill-pfad`“ in §3 von `slice-transformationen-map-value` (Frist der Meldung: dessen Start). Zweitens: die Kosten je Lesung des Regelstands je Block (§6, „weiter offen“); Adresse: das Closure-Kriterium in §3 der Welle. Drittens: der Zeitpunkt, ab dem ein Regelstand für einen Run gilt (ab dem Öffnen des Snapshots, Festlegung des Implementers, Review F-2), steht in keinem Spec-Satz und in keiner Entscheidung; Adresse: `slice-transformationen-betriebsdoku` §2 (das Handbuch führt das Ist-Verhalten als Zusage, nicht als Spec-Aussage). *(d) Gelernt, bei 1× keine Regel:* eine Zusage über die Kopplung an einen Folge-Slice („der Wirkort bleibt unverändert“) hängt an den Typen, über die der Wirkort **heute** vergleicht und iteriert; der Review fand sie durch Lesen der Instanziierung von `sameSet` — keine Mutation hätte sie gefunden, denn `make test` bleibt grün und der Fehler erscheint erst am Übersetzer des Folge-Slice.
- **Beobachtungs-Register (`../observations/`):** je Anfall eine Datei `evidence/slice-transformationen-backfill-pfad.md`, Zähler = Zahl der Dateien (gemessen mit `ls evidence | wc -l` am Stand dieser Closure). *Neue Belege:* `BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` **4×** (V-4 LOW; offen, Schwelle erreicht, Vorschlag im `state.md`), `BEO-PGC/slice-pfad-als-link-in-berichten` **4×** (F-7 INFO, V-6 INFO; verkörpert — die benannte Grenze (3) der `structure`-Regeln trat ein, das Gate fing die Form am Move, keine Schärfung, Begründung im `state.md`). *Neuer Eintrag:* `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` **3×** (F-8 INFO; drei Belegdateien aus drei Vorgängen; Schwelle erreicht, Vorschlag und Adresse im `state.md`). *`state.md` nachgezogen ohne neue Datei:* `BEO-PGC/run-fehlerklasse-schema-im-transformations-backfill` (Abbildung `schema`/`configuration` getragen, Beleg-Anker), `BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` (kein Auftreten, Zähler bleibt 2×). *Deckel-Fall ohne Datei, Finding-Kennung hier* (verkörpert ab 10×, vor dem Merge von Reviewer bzw. Verifier gefunden, Schwere ≤ LOW, bekannter Träger-Typ): F-3 und V-2 (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`, Deckel bei 32×: `sameSet` und der Plan `slice-transformationen-map-value`, Träger-Typ fremder Plan). *Benannte Grenze, kein Register-Anfall:* F-4 und V-3 — die Verdrahtungszeile `Transformations: activation` ist in `make test` nicht gebunden (Mutation grün), im E2E rot; die Zusage „belegt am laufenden System“ ist gebunden. Zählte man sie zu `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (Deckel bei 14×), wäre sie ein Deckel-Fall (INFO, vor dem Merge gefunden); der Träger-Typ „Verdrahtungszeile der Composition Root“ ist im Bestand der Klasse nicht belegt (gemessen: keine der 16 Belegdateien nennt `Verdrahtung` oder `wiring`), ein zweites Auftreten dieses Typs bekäme deshalb eine Datei. *Kein Register-Anfall:* V-1 (die Zahl der Lesungen im Doc-Kommentar ist eine Strukturaussage, an Code und Test gebunden; die Kennzeichnung trägt §3, das Risiko §6), F-1 und F-2 (Fragen an den Architect, Ausgang in §6), F-5 (konform), F-6 und V-5 (Coverage-Streuung als Lauf-Beleg), V-7 (übernommene Messungen, im Report benannt). *Kein Anfall:* `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` und `BEO-PGC/test-integration-retention-timing-flake` (§6), `BEO-PGC/test-runner-stiller-ausschluss` (die Phase ist eine Shell-Phase ohne `-run`-Muster). *Lese-Schritt der Closure von `welle-transformationen`:* aus diesem Slice erreicht neu `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` die Schwelle ohne Ausgang, `BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` trägt den vierten Beleg (Vorschlag im `state.md`); die übrigen Einträge mit neuer Datei haben einen Ausgang.
- **Folge-Slices:** keine angelegt. Übergaben mit Adresse (gemeldet, Frist: diese Closure, gezogen): `slice-transformationen-map-value` (`open/`) — der Übergabe-Block zu `sameSet`, dem Konflikt mit seinem DoD Punkt 2 und den zwei Fixture-Schaltern (V-2); `slice-transformationen-betriebsdoku` (`open/`) — §2 trägt die Run-Abhilfe nach [`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md) Folgepflicht 4 (Klasse `schema` bei Nichtanwendbarkeit, `configuration` beim Wechsel des Regelstands, neuer `cdc.backfill_table`-Antrag, kein Neustart) und den Bezugspunkt des Regelstands im Run (V-4); die Welle [welle-transformationen](../welle-transformationen.md) — §3 trägt das Closure-Kriterium zu den Lesekosten des Runs (V-1, Review F-1). Die **Start-Bedingung von `slice-transformationen-map-value`** ist mit diesem Move erfüllt (gelesen in dessen §4: dieser Slice liegt nach dem Move in `done/`, in `in-progress/` liegt nur die Roadmap); der Plan ist bis auf den Übergabe-Block nicht geändert.
- **Risiken aus §6:** je ein Ausgang, mit Beleg in §6. *Entfallen:* Prüfung steht zweimal (Kern einmal, Schleife je Pfad benannt) · zwei Wege des Bild-Baus · Fail-closed zu lasch · Tier des Paritätstests · Laufzeit von `make test-integration` (5 min 14 s bis 23 s, gemessen von Reviewer und Verifier, je ein Lauf) · Fenster zwischen Regel-Setzbarkeit und Backfill-Bindung (mit dieser Closure) · Fehlerpfad-Kommentare. *Weiter offen:* Kosten der Lesung je Block (Adresse: Closure-Kriterium der Welle).
- **Drei Paarungen:** dieser Slice gehört zu [welle-transformationen](../welle-transformationen.md) (offen) — die Closure der Welle prüft sie mit; die Slice-Closure trägt sie zusätzlich jetzt: *Anker:* der Lerneintrag (a) trägt den Anker `BEO-PGC/run-fehlerklasse-schema-im-transformations-backfill` (`state.md`, am Ort existent) und die genannten Tests, die am Ort stehen; (b) ist ein Vorschlag an den Architect, kein `liegt in`-Feld eines verkörperten Ziels, seine Zielorte (`AGENTS.md` §3.1, `.claude/agents/`) existieren; *Folge-Slice:* keiner angelegt; die Übergabe-Adressen `slice-transformationen-map-value` und `slice-transformationen-betriebsdoku` existieren als Dateien in `open/`, das Closure-Kriterium steht in `welle-transformationen.md`; *Register:* jede genannte Kennung `BEO-PGC/<slug>` existiert als Verzeichnis mit nicht leerem `evidence/` (geprüft mit `ls docs/plan/planning/observations/BEO-PGC/<slug>/evidence`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Use Case des Runs, Composition Root und Test-Runner
sind keine eigenen Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/start-trigger-ohne-uebergabe-artefakt` (offen, 1×, einschlägig —
Start-Trigger), `BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft` (offen, 2×,
einschlägig — das Verdikt ist Vorab-Bedingung im Start-Trigger, nicht
Rückführungs-Bedingung), `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`
(verkörpert, 6×), `BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×,
Plan-Zeile Runner), `BEO-PGC/test-integration-retention-timing-flake`
(verkörpert, 3×, Risiko §6),
`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (offen, 2×),
`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` (offen, 1×, für die
Run-Klasse `schema` nicht einschlägig: sie ist mit
[`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md)
(`Supersedes` für einen Satzteil von
[`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md)
Teilfrage 5) entschieden),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf §3).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
