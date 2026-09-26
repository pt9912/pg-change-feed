# Review-Report: slice-transformationen-kern-rename — 2026-09-26

**Review-Art:** Code — der Diff liefert den Regeltyp `rename_column` in der Domäne (Konstruktor mit
Invarianten, `CheckApplicable`, reine Auswertung), die um den Regelsatz erweiterte gemeinsame
Row-Image-Funktion `BuildRowImage`, den Regelstand der Bindung im `Assembler` (`SetTransformation`,
`RemoveTransformation`, Erhalt bei `AddBinding`-Merge und `setSchemaVersion`), die
Anwendbarkeits-Prüfung vor der Serialisierung samt `ErrTransformationNotApplicable`, die vier neuen
Domänen-Sentinels und die Abbildung in `classifyRunError`; geprüft gegen Plan, ADRs, Spec-Stellen und
`AGENTS.md` Hard Rules (Modul 10 §Drei Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe
(Modul 11).

**Gegenstand:** Slice `slice-transformationen-kern-rename` (Welle `welle-transformationen`), Diff-Range
`7310dbd1..770fc754` (7 Commits, 14 Dateien, +1474/−97; Baum sauber, nicht gepusht).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere
HIGH-Klassen ergänzt, u. a. Kommentar-Zusage, Zahl-im-Träger mit `suchlauf`-Probe, Beleg-Satz,
Zusage-ohne-Eingabeseite, Nachzug-Nachbar).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-26.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-transformationen-kern-rename` (§1 Ziel und Abgrenzung, §2 DoD-Wortlaut als Bezug,
  §3 Plan mit Festlegungen, Suchlauf-Feld und Belegen des Laufs, §6 Risiken)
- [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) (Teilfrage
  2/3/4/5/6, Folgepflicht 2 und 7),
  [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) (Teilfrage 2, die eine
  Konstruktionsstelle), [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) (Muster
  `ExcludeColumn`), [`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md) (sieben Fehlerklassen)
- [`SPEC-030`](../../spec/pflichtenheft.md) (Regelform, Bezeichner, Position, Anwendbarkeit),
  [`SPEC-002`](../../spec/pflichtenheft.md), [`SPEC-008`](../../spec/pflichtenheft.md),
  [`SPEC-019`](../../spec/pflichtenheft.md) (K1–K4 und Fehlertexte),
  [`ARC-001`](../../spec/architecture.md), [`ARC-005`](../../spec/architecture.md),
  [`ARC-007`](../../spec/architecture.md)
- [`LH-FA-CFG-007`](../../spec/lastenheft.md), [`LH-FA-CFG-005`](../../spec/lastenheft.md),
  [`LH-QA-SEC-004`](../../spec/lastenheft.md), [`LH-FA-DAT-005`](../../spec/lastenheft.md),
  [`LH-FA-CAP-008`](../../spec/lastenheft.md), [`LH-FA-ADM-003`](../../spec/lastenheft.md),
  [`LH-FA-SCH-004`](../../spec/lastenheft.md)
- `AGENTS.md` (Hard Rules §3.1, §3.3, §3.7, §3.9, §3.12, §3.13), `harness/conventions.md`
  (`MR-000`/`MR-001`/`MR-002`)
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-harness-suchlauf-nachmessen.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem Implementer-Bericht übernommen;
Exit-Codes ungepiped in Log-Dateien gesichert, gedruckte Zeilen zitiert):

- **Sensoren am Stand `770fc754`:** `make test` (Race-Detector) Exit 0, 42 Pakete `ok`, keine
  `FAIL`-/`panic`-Zeile; `make a-check` Exit 0, gedruckt „gesamt: 0 Befund(e)“; `make coverage-gate`
  Exit 0, gedruckt „coverage-gate: OK — Coverage 83.80% erfüllt Schwelle 80%“; `make gates` (einmal)
  Exit 0, gedruckt „coverage-gate: OK — Coverage 83.80% erfüllt Schwelle 80%“, „d-check: 1208 Datei(en)
  geprüft, 0 Befund(e)“, „generated-sync: OK“, „gesamt: 0 Befund(e)“, `git status --short` danach leer;
  `make commit-traceability RANGE=7310dbd1..770fc754` Exit 0, gedruckt „OK — 7 Commit(s) in
  "7310dbd1..770fc754", Betreffs ohne Struktur-ID“.
- **Suchlauf-Feld des Plans:** `make suchlauf-nachmessen PLAN=<Plan-Datei>` Exit 0, gedruckt
  „suchlauf-nachmessen: 14 Zeilen stimmen“. **Von Hand** mit `git grep` nachgefahren, ohne das Werkzeug:
  `BuildRowImage` über `internal/*.go` am Parent `1852b2f9` 32 Trefferzeilen, am Diff 49, ohne Tests am
  Diff 10; `TableBinding{` am Diff 52, ohne Tests 5 Stellen (`config_file.go:250`, `wiring.go:370`, `:381`,
  `:425`, `:1404`); `ErrIncompatibleSchemaChange` über `internal docs spec harness` ohne Records 35 Zeilen
  roh, davon 3 in der Plan-Datei, also 32; `BuildRowImage(columns` über `docs spec harness` ohne Records
  6 Zeilen roh, davon 3 in der Plan-Datei, also 3 (`0115` Zeilen 129 und 196, `backfill-pfad` Zeile 172).
  Träger, die der Bericht als „nicht gefunden“ führt, selbst gesucht: `classifyRunError` in
  `docs spec harness` (nur `0049`, `0112`, offene Pläne mit eigener Adresse), `ErrTruncateUnsupported`
  (0 Treffer außerhalb des Codes), `TableBinding` in `docs spec harness` (nur `Accepted`-ADRs, Register
  und offene Pläne). Kein weiterer beschreibender Träger gefunden.
- **Byte-Gleichheit ohne Regel:** `git diff -w 7310dbd1..770fc754 -- internal/domain/model/rowimage_test.go`
  gelesen: die Erwartungs-Zeichenketten der Byte-Tabelle stehen unverändert, die Tabelle ist in die
  Paketvariable `rowImageByteCases` gezogen und läuft mit `nil` und mit `[]Transformation{}` durch
  dieselbe Erwartung.
- **Kollisions-Probe an `BuildRowImage`** (Wegwerf-Test in einer Kopie des Baums): Spalten `id`, `name`,
  `customer_name`, Regel `name` nach `customer_name` liefert `{"id":"1","customer_name":"Ada","customer_name":"Kunde"}`
  ohne Fehler; zwei Regeln mit gleichem Ziel `z` auf `a` und `b` liefern `{"id":"1","z":"x","z":"y"}`;
  die Kette `a` nach `b`, `b` nach `c` auf den Spalten `id`, `a`, `b` liefert `{"id":"1","b":"x","c":"y"}`
  (einmalige Auswertung gegen die Original-Spaltennamen, keine sequentielle).
- **Eingabeseiten-Mutationen** (in einer Kopie des Baums per Zeichenketten-Ersetzung gesetzt, je einzeln
  gefahren, Kopie danach aus `git show HEAD:` zurückgesetzt; das Arbeitsverzeichnis blieb unberührt):

  | Nr. | Mutation | Ergebnis |
  |---|---|---|
  | M1 | umbenannter Schlüssel ans Ende des Bilds gesetzt | rot: `TestBuildRowImageRenameColumn`, `TestConsumeRuleTargetComparisonIsCaseSensitive`, `TestConsumeRuleSurvivesSchemaBump`, `TestSetAndRemoveTransformationOnLiveBinding` |
  | M2 | Kollisions-Prüfung in `CheckApplicable` entfernt | rot: `TestCheckApplicable`, `TestConsumeRuleTargetCollidesWithRelationColumnIsNotApplicable`, `TestConsumeRuleTargetCollisionAfterCompatibleExtension` |
  | M3a | Ausschluss-Prüfung in `BuildRowImage` entfernt | rot: mindestens zehn Tests in beiden Paketen, darunter `TestBuildRowImageExcludedValueNowhere` |
  | M3b | Ausschluss erst gegen den Zielschlüssel (Regel vor Ausschluss) | rot: `TestBuildRowImageRenameColumn`, `TestExcludedColumnIsUnreachableForEveryRuleKind` |
  | M4 | `AddBinding` erhält den Regelstand nicht | rot: `TestAddBindingKeepsRuleState` |
  | M5 | Grenze `to` 63 auf 64 | rot: `TestNewRenameColumnInvariants` |
  | M6 | Prüfung erst nach dem Vorrücken der Sequenz | rot: `TestConsumeRuleRemedyRestoresCapture` |
  | M7a/b | `withTransformation`/`withoutTransformation` bauen in der übergebenen Liste um | rot: `TestTransformationListsAreReplacedNotMutated` |
  | M8 | `SetTransformation` ohne Sperre, `-race` | rot: `WARNING: DATA RACE`, `TestAssemblerTransformationsAreRaceFree` |
  | M9 | `ErrTransformationNotApplicable` aus `classifyRunError` entfernt | rot: `TestClassifyRunErrorMapsKnownSentinelsToADR0023Classes` |
  | M10 | `TransformationKinds` liefert einen zweiten Typ | rot: `TestTransformationKindsIsAClosedSet`, `TestExcludedColumnIsUnreachableForEveryRuleKind` |
  | M11 | Zielname ohne `encoding/json`-Maskierung | rot: `TestBuildRowImageBytes`, `TestBuildRowImageRenameColumn`, `TestBuildRowImageEmptyRuleSetKeepsBytes` |
  | M12 | `setSchemaVersion` verliert den Regelstand | rot: `TestConsumeRuleSurvivesSchemaBump`, `TestConsumeRuleTargetCollisionAfterCompatibleExtension` |
  | M13 | `RemoveTransformation` ignoriert den Namen | rot: zwei Tests |
  | M14b | letzte statt erster treffender Regel entscheidet | rot: `TestBuildRowImageRenameColumn` |
  | M16 | Alt-Image ohne Regelstand | rot: `TestConsumeRenameColumnAppliesToBothImages`, `TestExcludedColumnIsUnreachableForEveryRuleKind` |
  | M17 | Spalten-Prüfung in `CheckApplicable` entfernt | rot: drei Tests (`TestCheckApplicable`, zwei im `mapper`-Paket) |
  | M18 | Wert verändert | rot: mindestens neun Tests |
  | M20 | Länge in Runes statt Byte | rot: `TestNewRenameColumnInvariants` |
  | M14a | Auswertung als Kette (Spalte nach der Regel neu belegt) | **grün** — bei angewendeter Vorbedingung `CheckApplicable` nicht erreichbar (Kette verletzt die Kollisionsprüfung), damit kein Befund |
  | M15 | `checkTransformations` prüft nur die erste Regel | **grün** (F-2) |
  | M19 | Prüfreihenfolge Kollision vor fehlender Spalte | **grün** (F-6) |

- **Benchmark (§6-Risiko „Kosten der Prüfung“) nachgemessen:** `BenchmarkAssemblerChange` am Parent
  `1852b2f9` (Benchmark-Datei ohne die Regel-Variante in die Kopie gelegt) und `…WithRules` am Diff, je
  zwölf Läufe (`-count 6 -benchtime 2s`, zwei Durchgänge, ohne `-race`, im Toolchain-Container von
  `make test` mit anderem Kommando): Parent 4640–5140 ns/op, Diff ohne Regel 4846–5252 ns/op, Diff mit
  zwei Regeln 5043–5456 ns/op, überall 2444 B/op und 86 allocs/op. Mediane (abgeleitet): 4734, 4979,
  5149 ns/op, also +245 ns und +415 ns. Die Spanne des Parents allein (500 ns) ist so breit wie der
  Aufschlag; die Aussage „0,2–0,6 µs, nicht gegen Rauschen abgesichert“ trägt die Messung.
- **Kommentar-Probe nach `AGENTS.md` §3.7:** `git diff 7310dbd1..770fc754 --name-only -- '*.go' | xargs -r grep -nE 'slice-[0-9]+|welle-[0-9]+|vor diesem|nach diesem|seit diesem|weiterhin|früher|bisher|neu '`
  liefert nur „neu aufgebaut“/„neu registriert“ im Sinn der Sache und Vorbestand; keine
  Slice-/Wellen-Chronik. Die namentliche Adresse an `service.go:486-488` ist ein Rang-Zeiger.
- **Umgebung:** ein schwerer Docker-Lauf zugleich (`free -m` vor dem `make gates`-Lauf 17,6 GB
  verfügbar); ein eigenes Volume `review-gocache` für den Go-Build-Cache der Mutationsläufe angelegt und
  nach dem Lauf entfernt; kein `prune`.
- **Commits:** alle sieben Betreffs nennen `LH-FA-CFG-007` und `ADR-0112`, keiner trägt `SPEC-*`/`ARC-*`;
  die beiden Lifecycle-Moves (`64b8eb9e`, `1852b2f9`) sind reine Renames (100 %, 0 Zeilen), Inhalt steht in
  eigenen Commits (`AGENTS.md` §3.3).
- **Nicht gefahren (Grenze):** `make test-store`, `make test-replication`, `make test-integration`
  (der Diff führt keinen Weg zur Datenbank, keinen SQL-Zugriff und keine Betreiber-Oberfläche ein).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Der Doc-Kommentar von `BuildRowImage` sagt „eine nicht anwendbare Regel wirkt hier nicht“; für den Fall „Zielname kollidiert“ trägt der Code das nicht: die Regel wirkt und das Bild trägt zwei gleichnamige Schlüssel ohne Fehler (Probe: `{"id":"1","customer_name":"Ada","customer_name":"Kunde"}`). Der zweite Halbsatz („ergäbe zwei gleichnamige Schlüssel“) widerspricht dem ersten im selben Satz; der Plan §3 („eine nicht anwendbare Regel wirkt in der Funktion nicht“) trägt dieselbe Zusage. `jsonb` behielte bei doppeltem Schlüssel den letzten Wert; der künftige Backfill-Aufrufer liest diesen Kommentar als Vertrag seiner Vorbedingung. | Kommentar-Klasse Zusage (Reviewer-Skill, HIGH-Liste); `AGENTS.md` §3.7; [`SPEC-030`](../../spec/pflichtenheft.md) Anwendbarkeit | `internal/domain/model/rowimage.go:29-31`; Plan §3 Festlegungen, dritter Punkt | ja — Wegwerf-Test `BuildRowImage` mit Spalten `id`, `name`, `customer_name` und Regel `name` nach `customer_name`; kein Test der Datei bindet das Verhalten bei verletzter Vorbedingung | Kommentar behauptet nicht getragenen Fehlerpfad |
| F-2 | MEDIUM | Die Anwendbarkeits-Prüfung ist nur für die erste Regel einer Bindung an ihre Eingabeseite gebunden: eine Mutation, die `checkTransformations` auf `rules[:1]` beschränkt, lässt beide betroffenen Pakete grün. Alle Nichtanwendbarkeits-Tests tragen genau eine Regel; der Fall „zweite Regel nicht anwendbar, erste anwendbar“ hat keinen Test. | Maintainability; fehlende Negativtests bei neuem öffentlichem Vertrag (Reviewer-Skill, MEDIUM-Liste); [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 4 („prüft die Regeln der Bindung“) | `internal/adapters/driving/replication/mapper/mapper.go:580-588`; Tests `transformation_test.go:170-262` | ja — Mutation M15 an einer Kopie, `make test` bleibt Exit 0 | Prüfung ohne Testbindung an die Position der Regel |
| F-3 | LOW | Die Tabelle in Plan §3 trägt überholte Zeilen neben ihren Nachzugs-Zeilen: „Aufrufer der gemeinsamen Funktion im Backfill-Pfad (falls vorhanden)“ steht neben der konkreten Zeile zu `service.go`, `mapper_test.go` als „update / neu“ neben der Zeile, die die Mapper-Tests in eigene Dateien legt und `mapper_test.go` unverändert lässt (`git diff 7310dbd1..770fc754 --stat` bestätigt: `mapper_test.go` nicht im Diff), „Datei und Name am Start gemessen“ neben dem gemessenen Namen. Nur die letzte Zeile verweist auf die ersetzte. | Maintainability; Nachzug widerspricht dem Nachbarn im selben Träger (Reviewer-Skill, MEDIUM-Liste — hier LOW: die Zeilen widersprechen sich nicht in der Aussage, sie stehen doppelt) | `docs/plan/planning/in-progress/slice-transformationen-kern-rename.md` §3, Zeilen 177–182 gegen 183–187 | ja — Lesen der Tabelle | Plan-Tabelle trägt überholte Zeile neben ihrem Nachzug |
| F-4 | LOW | Die Belege des Laufs nennen für `make coverage-gate` „83.70%“ ohne Stand und für `make gates` am Stand `cfca864b` „83.80%“; gemessen am Stand `770fc754` (Code unverändert seit `06719d4f`) zweimal 83.80 %. Die erste Zahl trägt keinen Lauf-Stand, der ihre Abweichung erklärt. | `AGENTS.md` §3.12 Instanz A; Zahl-im-Träger (Reviewer-Skill) — LOW: beide Zahlen stehen mit ihrem gedruckten Wortlaut im selben Absatz, die aktuelle steht mit Stand daneben | Plan §3 „Belege des Laufs“, Punkt „Sensoren“ | ja — `make coverage-gate` | Zahl im Träger ohne Lauf-Stand |
| F-5 | INFO | Die Domäne prüft die Konfliktfreiheit zwischen Regeln nicht und der `Assembler` ebenso wenig: zwei Regeln mit gleichem Zielnamen auf verschiedenen Spalten sind je einzeln anwendbar und ergeben `{"id":"1","z":"x","z":"y"}` ohne Fehler (Probe). K3 „Zielnamen untereinander verschieden“ liegt laut Plan §1 beim Antragsweg (`slice-transformationen-antragsweg-usecase`); `SetTransformation` benennt das in seinem Doc-Kommentar. Zur Laufzeit im Erfassungspfad steht dahinter kein zweiter Wächter. | [`SPEC-030`](../../spec/pflichtenheft.md) Randfälle und [`SPEC-019`](../../spec/pflichtenheft.md) K3; [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 3 | `internal/adapters/driving/replication/mapper/mapper.go:509-521`; `internal/domain/model/transformation.go:86-95` | ja — Wegwerf-Test | Rollen-übergreifende Invariante ohne Laufzeit-Wächter (delegiert, benannt) |
| F-6 | INFO | Drei Kleinigkeiten am Rand, keine Aktion erwartet: (a) wenn beide Bedingungen verletzt sind (Spalte fehlt, Zielname kollidiert), bestimmt die Prüfreihenfolge in `CheckApplicable` den Sentinel; kein Test bindet sie (M19 grün); (b) der Konstruktor prüft Zielname und Spalte nicht auf UTF-8-Gültigkeit — `rule_spec` als `jsonb` liefert nur gültiges UTF-8, ein ungültiger Wert würde in `json.Marshal` zu U+FFFD; (c) `NewAssembler` und `AddBinding` teilen die übergebene Regelliste mit dem Aufrufer (flache Kopie, wie `ExcludedColumns`); der Vertrag „ab dem Schreiben unverändert“ steht im Doc-Kommentar von `TableBinding`. | Maintainability | `internal/domain/model/transformation.go:86-95`, `:52-63`; `mapper.go:82-99` | nein | Randfälle ohne Bindung (benannt) |
| F-7 | INFO | Der Kommentar an `blockBuilder.build` („Der Run übergibt keine Transformationsregeln … Adresse `slice-transformationen-backfill-pfad`“) ist ein Rang-Zeiger und wahr am Stand des Diffs; er ist zeitgebunden und muss mit dem Slice, den er nennt, umgeschrieben werden (der Plan trägt das als §6-Risiko „Zwischenzustand im Backfill-Pfad“). Die gemeldeten fremden Träger (`slice-transformationen-backfill-pfad` Zeile 172, [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md) Zeilen 129 und 196) sind vollständig gefunden und richtig behandelt: die `Accepted`-ADR bleibt unberührt, der offene Plan geht mit Frist „Closure dieses Slice“ an den Planner. | `AGENTS.md` §3.13 (Träger-Nachzug) | `internal/application/usecase/backfill/service.go:486-488` | nein | Träger-Nachzug fremder Datei (gemeldet) |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/domain/model/rowimage.go` (`BuildRowImage`: Reihenfolge Ausschluss, Abwesenheit, Regel in einer Schleife; Position der Quellspalte; nil und leere Regelmenge; Escaping des Zielnamens; ausgeschlossene Spalte bringt weder Quell- noch Zielschlüssel noch Wert ins Bild) | geprüft, ohne Befund über F-1 hinaus: Mutationen M1, M3a/b, M11, M16, M18 rot; Byte-Tabelle mit unveränderten Erwartungswerten für `nil` und `[]Transformation{}` (`git diff -w`); einmalige Auswertung gegen die Original-Spaltennamen (Kette a→b, b→c entsteht nur bei verletzter Kollisionsprüfung) |
| `internal/domain/model/transformation.go` (Konstruktor, `CheckApplicable`, `applyTransformations`, Kinds-Menge) | geprüft, ohne Befund über F-5, F-6 hinaus: Invarianten je Eingabe an ihre Eingabeseite gebunden (M5, M20 rot; Bytes statt Runes durch den `ä`-Fall, U+0000 in `column` und `to`, leere Felder, `to` gleich `column` mit eigenem Sentinel); Nullwert trägt keine Spalte; erste treffende Regel entscheidet (M14b rot); `TransformationKinds` liefert je Aufruf eine neue Liste |
| `internal/domain/errors/errors.go` (vier Sentinels) | geprüft, ohne Befund: Namen und Texte unabhängig von der späteren Abbildung auf die Formzeilen aus [`SPEC-019`](../../spec/pflichtenheft.md) korrekt (`ErrTransformationTargetIsColumn` trennt den K3-Fall von `ErrInvalidTransformation`, `ErrEmptyIdentifier` trägt nur den Regelnamen-Fall dieses Konstruktors); die Kommentare tragen Zusage und Abgrenzung im Indikativ |
| `internal/adapters/driving/replication/mapper/mapper.go` (Regelstand, `SetTransformation`/`RemoveTransformation`, Erhalt, `checkTransformations`, `change`) | geprüft, ohne Befund über F-2 hinaus: Listen werden in beiden Helfern neu aufgebaut, kein Aliasing über `append` (M7a/b rot); Erhalt bei `AddBinding` und `setSchemaVersion` gebunden (M4, M12 rot); Set/Remove auf nicht getragener Bindung ohne Wirkung und ohne Beleben; Prüfung vor Serialisierung und vor dem Vorrücken der Sequenz (M6 rot), ein Schnappschuss der Bindung trägt Prüfung und Bild; spalten-entfernende Fälle enden weiter an `relationOther` (Test mit `stderrors.Is` in beide Richtungen); `%w`-Wrapping trägt Sentinel, Grund und Regelname; Sperre gebunden (M8 rot) |
| `internal/bootstrap/wiring.go` (`classifyRunError`) und `heartbeat_internal_test.go` | geprüft, ohne Befund: Abbildung auf `model.ErrorClassSchema` vor der Klasse `replication` eingeordnet; M9 rot; der Pfad `receive.Stream.process` gibt den Fehler des `Assembler` unverändert weiter und erreicht `Capture` nicht — der Kommentar an `ErrTransformationNotApplicable` trägt nur, was der Code trägt |
| `internal/application/usecase/backfill/service.go`, `internal/adapters/driven/postgressnapshot/snapshot_test.go` | geprüft, ohne Befund über F-7 hinaus: Aufrufer für den vierten Parameter nachgezogen, `nil` als Regelmenge, Kommentar als Rang-Zeiger |
| Tests (`transformation_test.go` beider Pakete, `transformation_internal_test.go`, `mapper_bench_test.go`, `rowimage_test.go`) | geprüft, ohne Befund über F-2, F-6 hinaus: 20 der 23 gefahrenen Mutationen färben je mindestens einen Test rot (die drei grünen: M14a als bei erfüllter Vorbedingung unerreichbar, M15 in F-2, M19 in F-6); der Eigenschaftstest zählt die Regeltypen aus `TransformationKinds()` auf und bricht bei unbekanntem Typ ab (M10 rot); Determinismus- und Nebenläufigkeits-Test laufen unter `-race` in `make test` (Exit 0); der Whitebox-Test bindet die Schnappschuss-Zusage mit freier Kapazität der Ausgangsliste |
| Benchmark und §6-Risiko „Kosten der Prüfung“ | geprüft, ohne Befund: Aussage „0,2–0,6 µs“ und „nicht gegen Rauschen abgesichert“ mit eigener Messung bestätigt (Mediane 4734, 4979, 5149 ns/op; identische Allokationen) |
| Suchlauf-Feld des Plans §3 (14 Zeilen, beide Stände) | geprüft, ohne Befund: mit dem Werkzeug Exit 0 und sieben Zeilen von Hand nachgefahren; Nichtgefundenes je Träger steht im Feld; die im Bericht als „nicht gefunden“ geführten Träger sind auch bei eigener Suche nicht auffindbar |
| Kommentare nach `AGENTS.md` §3.7 in allen geänderten Go-Dateien | geprüft, ohne Befund über F-1 hinaus: keine Slice-/Wellen-Chronik in Produktionspfaden, keine verworfene Alternative im Konjunktiv außer dem Halbsatz aus F-1, Test-Kommentare nennen Test als Subjekt |
| Hard Rules und Commit-Struktur | geprüft, ohne Befund: `make a-check` „gesamt: 0 Befund(e)“ (Domäne importiert nichts aus anderen Schichten), Coverage 83.80 % gegen 80 %, Docker-only (keine Toolchain-Installation), keine Suppression, `git mv` rein und Inhalt getrennt, Betreffs mit `LH-FA-CFG-007`/`ADR-0112` und ohne `SPEC-`/`ARC-` |
| Handbuch-Kandidatenlauf (`docs/user/benutzerhandbuch.md`) | geprüft, ohne Befund: keine Änderung am Handbuch im Diff, keine neue `CDC_*`-Variable, keine `cdc.*`-Funktion, kein Endpunkt (die Regel ist nur über die Assembler-Methoden setzbar); die Handbuch-Nachzüge stehen als Adresse im Plan (`slice-transformationen-betriebsdoku`) |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Kommentar behauptet nicht getragenen Fehlerpfad · Prüfung ohne Testbindung an die Position der Regel · Plan-Tabelle trägt überholte Zeile neben ihrem Nachzug · Zahl im Träger ohne Lauf-Stand · Rollen-übergreifende Invariante ohne Laufzeit-Wächter (delegiert, benannt) · Randfälle ohne Bindung (benannt) · Träger-Nachzug fremder Datei (gemeldet)

Hinweis zum Zähler: die Klasse „Kommentar behauptet nicht getragenen Fehlerpfad“ trägt im
Beobachtungs-Register den Stand offen, 2× (`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`,
laut Plan §8); F-1 wäre das dritte Auftreten (der Kommentar an `ErrTransformationNotApplicable`, den die
DoD ausdrücklich an dieser Klasse ausrichtet, trägt dagegen fehlerfrei — der Fund liegt am Nachbarkommentar
von `BuildRowImage`).

## Verdikt

**Merge-blockierend:** ja — F-1 (HIGH, Kommentar-Zusage, die der Code für den Kollisionsfall nicht trägt)
und F-2 (MEDIUM, Prüfung der zweiten und jeder weiteren Regel ohne Test) gehören in eine Fixrunde; beide
sind klein (ein Satz im Doc-Kommentar samt Plan-Festlegung, ein Test mit zwei Regeln). Der Rest des Diffs
trägt ohne Befund über LOW/INFO hinaus: Reihenfolge Ausschluss vor Regel, Position, Byte-Gleichheit ohne
Regel, Konstruktor-Invarianten, Schnappschuss-Zusage, Erhalt des Regelstands und Prüfung vor der
Serialisierung sind gegen 20 rot gesehene Mutationen an der Eingabeseite gebunden.

**Übergabe:** Findings gehen an den Implementer (Fixrunde nötig, deshalb bleibt die DoD-Zeile „Review
durchgeführt“ im Plan offen und wird bei Schritt 21 des Implementer-Workflows nachgezogen); die
**Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. F-3 (LOW) und
F-7 (gemeldete fremde Träger) gehen mit der Frist „Closure dieses Slice“ an den Planner beziehungsweise
werden dort mitgezogen. Dieser Report selbst ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill,
dieses Modell, dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und muss es nicht. Der
Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat (Modul 11; anderes
Prüf-Artefakt, anderer Eingabe-Kontext).
