# Slice transformationen-kern-rename: Kern — Regeltyp `rename_column` in der Domäne, Regelstand und Nichtanwendbarkeits-Prüfung im Assembler

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
Negative), [`LH-FA-CFG-005`](../../../../spec/lastenheft.md) (Ausschluss gilt
zuerst), [`LH-QA-SEC-004`](../../../../spec/lastenheft.md) (Wert nirgends im
Image), [`LH-FA-DAT-005`](../../../../spec/lastenheft.md) (Abwesenheit
erkennbar), [`LH-FA-CAP-008`](../../../../spec/lastenheft.md) (fehlendes Bild
bleibt fehlend), [`LH-FA-ADM-003`](../../../../spec/lastenheft.md) (sichtbarer
Fehlerzustand), [`LH-FA-SCH-004`](../../../../spec/lastenheft.md) (Fehlerpfad
der Schemaänderung),
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Teilfrage 2/3/4/5/6 und Folgepflicht 2,
[`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 2
(die gemeinsame Row-Image-Funktion),
[`ADR-0059`](../../adr/0059-spaltenauswahl-mechanismus.md) (Muster
`ExcludeColumn`), [`ADR-0023`](../../adr/0023-fehlerklassifikation.md)
(Fehlerklassen).

**Berührte Spec-Stellen:** [`SPEC-002`](../../../../spec/pflichtenheft.md) (Row
Images), [`SPEC-008`](../../../../spec/pflichtenheft.md) (Fehlerklasse
`schema`), [`ARC-001`](../../../../spec/architecture.md) (Domain Core: reine
Funktionen), [`ARC-005`](../../../../spec/architecture.md) (der Assembler ist
Teil des Driving Adapters Replication Stream),
[`ARC-007`](../../../../spec/architecture.md) (`classifyRunError` in der
Composition Root) — gelesen, nicht geändert (die Spec trägt
`slice-transformationen-spec-nachzug`).

**Verantwortlich:** Implementer-Agent, 2026-09-26.

**Autor:** Planner-Agent, Welle-Eröffnung
[welle-transformationen](../welle-transformationen.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** `rename_column` wirkt im Erfassungspfad: die Domäne trägt den
Regeltyp mit Konstruktor-Invarianten und einer reinen Auswertung, die **in der
gemeinsamen Row-Image-Funktion** nach dem Ausschluss ansetzt; der `Assembler`
hält den Regelstand je Bindung als unveränderliche Liste, ersetzt sie unter
`tablesMu` und erhält sie bei `AddBinding`-Merge und `setSchemaVersion`; eine
auf eine Change nicht anwendbare Regel endet den Erfassungspfad sichtbar
(`ErrTransformationNotApplicable`, Klasse `schema`), bevor ein Wert
serialisiert wird. Die Regel ist in diesem Slice nur über die
Assembler-Methoden setzbar — kein SQL-Weg.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Antragsweg, Use Cases, Schema, Dauerhaftigkeit** — `antragsweg-schema` und
  `antragsweg-usecase`; hier gibt es keine Datenbank und keinen Antrag. Ein
  Regelstand entsteht nur über `SetTransformation`.
- **K1–K4 je Tabelle** — die Konfliktfreiheits-Invarianten sind Prüfungen beim
  Antrag (Use Case, ihr Wortlaut steht in
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 3); der `Assembler` prüft die **Anwendbarkeit** an einer konkreten
  Change, nicht die Konfliktfreiheit der Regelmenge. Die Domäne trägt nur die
  Invarianten einer einzelnen Regel: `column` und `to` nicht leer und ohne
  U+0000, `to` höchstens 63 Byte UTF-8
  ([`SPEC-030`](../../../../spec/pflichtenheft.md), Bezeichner). Der Use Case
  bildet eine Verletzung dieser Invarianten auf `rule_spec ist ungültig` ab
  ([`SPEC-019`](../../../../spec/pflichtenheft.md), fünfte Formzeile, Adresse
  Regelname). `to` gleich `column` gehört zu K3 (`Zielname kollidiert mit einer
  Spalte der Tabelle`, Adresse Zielname, geprüft nach den Formzeilen): trägt die
  Domäne dafür einen eigenen Sentinel, bildet der Use Case ihn auf diesen Text
  ab und nicht auf `rule_spec ist ungültig`. Das Alphabet des Regelnamens
  (`[a-z0-9_]{1,63}`) und der Text `Regelname ist ungültig` gehören nicht zu
  diesem Slice: `antragsweg-usecase` legt fest, an welcher Stelle ein leerer
  oder ungültiger Regelname geprüft wird; der Text der Spec ist bindend.
- **`map_value`** — `slice-transformationen-map-value`; der Regeltyp-Satz ist
  hier ein Satz aus einem Typ, aber als Menge geführt, damit ein zweiter Typ
  keine Umbau-Kante braucht. Bis `map-value` kennt die Domäne nur
  `rename_column`; die Folge für den Antragsweg steht in
  `antragsweg-usecase` §6 (Zwischenzustand der Regeltyp-Menge).
- **Der Backfill-Pfad** — `slice-transformationen-backfill-pfad`; ein bereits
  vorhandener Aufrufer der gemeinsamen Funktion im Backfill-Pfad übergibt bis
  dahin die leere Regelmenge (benannt in §6).
- **Die Startreihenfolge** — `slice-transformationen-start-reihenfolge`.

## 2. Definition of Done

- [x] `rename_column` wirkt auf beide Images: der Schlüssel `column` steht
      unter `to`, am Ausgang des `Assembler` an der Position seiner
      Quellspalte (Relation-Spaltenreihenfolge; Festlegung von
      [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
      Teilfrage 3 — die Spec sagt am Lesepfad nur Schlüsselmenge und Werte zu,
      `jsonb` bewahrt die Reihenfolge nicht,
      [`SPEC-030`](../../../../spec/pflichtenheft.md)), der Wert bleibt
      unverändert; Abwesenheit
      bleibt Abwesenheit (NULL, unverändertes TOAST, ausgeschlossene Spalte,
      fehlendes Bild); `change_id`, `transaction_id`, `source_table_id`,
      `sequence`, `operation`, `schema_version`, `schema` und `table` bleiben
      unverändert; eine leere Regelmenge liefert Bytes, die dem Stand vor
      diesem Slice gleich sind. *Zu belegen durch:* `make test` (Race-Detector)
      — bestehende Mapper- und Domänen-Tests ohne geänderte Erwartungswerte,
      neue Tests je Eigenschaft.
- [x] Der `Assembler` trägt den Regelstand und die Prüfung:
      `SetTransformation`/`RemoveTransformation` ersetzen die Liste unter
      `tablesMu`, `AddBinding`-Merge und `setSchemaVersion` erhalten sie, ein
      Leser hält seinen Schnappschuss ohne eigene Sperre; ist die `column`
      einer Regel in `event.Relation.Columns` nicht enthalten oder kollidiert
      ihr Zielname mit einer Spalte der Relation, meldet `change`
      `mapper.ErrTransformationNotApplicable` **vor** jeder Serialisierung, die
      Transaktion wird weder persistiert noch bestätigt, und `classifyRunError`
      bildet den Fehler auf `model.ErrorClassSchema` ab. *Zu belegen durch:*
      `make test` — je ein Negativtest an seine Eingabe gebunden (Mutation der
      Prüfung färbt den Test rot,
      `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`), darunter der Fall,
      den erst diese Prüfung fängt: Zielname kollidiert nach einer kompatiblen
      Spalten-Erweiterung; die spalten-entfernenden Fälle enden weiterhin
      vorher an `relationOther`.
- [x] Die Fitness Function von
      [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
      steht: ein Eigenschaftstest im `mapper`-Paket (Regeltyp × ausgeschlossene
      Spalte: das Image trägt weder den Quellschlüssel noch einen Zielnamen
      noch den Quellwert, [`LH-QA-SEC-004`](../../../../spec/lastenheft.md)),
      der die Regeltypen aus der Domänen-Menge aufzählt statt aus einer zweiten
      Liste; ein Determinismus-Test (gleiche Regelmenge und Relation →
      byte-gleiches Image am Ausgang des `Assembler`); ein Nebenläufigkeits-Test (die Regelliste eines
      Lesers bleibt unter gleichzeitigem `SetTransformation` stabil). *Zu
      belegen durch:* `make test` mit Race-Detector; `make a-check` grün (die
      Regeltypen liegen in `internal/domain/**` und importieren aus keiner
      anderen Schicht); `make coverage-gate` grün (Domäne und
      `replication/mapper` liegen in der netzlos gemessenen Fläche — die
      `coverage`-Stufe des Dockerfile schließt nur die Pakete
      `postgresstorage`, `postgresack`, `postgressnapshot` und
      `replication/receive` selbst aus; kein neues Paket).
- [x] Der Kommentar-Träger folgt: der Doc-Kommentar von `TableBinding` und
      `AddBinding` nennt den Regelstand neben `ExcludedColumns`;
      `ErrTransformationNotApplicable` trägt einen Kommentar, der nur zusagt,
      was der Code trägt (der Erfassungspfad endet, kein ACK) — keine
      Behauptung, die Abhilfe wirke im gescheiterten Prozess
      ([`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
      führt sie als erwartet;
      `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`, offen, 2×).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: entfällt — kein öffentlicher Vertrag berührt (das
      Nachrichtenschema bleibt, keine Betreiber-Oberfläche); die Spec trägt
      `slice-transformationen-spec-nachzug`, das Handbuch
      `slice-transformationen-betriebsdoku`.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
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
| `internal/domain/model/transformation.go` (Arbeitsname; + Test) | neu | Regeltyp `rename_column`, Konstruktor-Invarianten, reine Auswertungsfunktion, keine geteilten Zustände; Menge der Regeltypen als eine Quelle. |
| die gemeinsame Row-Image-Funktion in `internal/domain/model` (aus `slice-backfill-row-image-gemeinsam`; Datei und Name am Start gemessen) | update | nimmt den Regelsatz als weiteren Parameter und wertet nach dem Ausschluss aus, innerhalb derselben Schleife; alle Aufrufer nachgezogen (Suchlauf). |
| `internal/adapters/driving/replication/mapper/mapper.go` | update | `TableBinding.Transformations`, `Assembler.SetTransformation`/`RemoveTransformation`, Erhalt bei `AddBinding`-Merge und `setSchemaVersion`, Anwendbarkeits-Prüfung in `change`, `ErrTransformationNotApplicable`. |
| `internal/adapters/driving/replication/mapper/mapper_test.go` (+ Datei für den Eigenschaftstest) | update / neu | Regeltests, Eigenschaftstest [`LH-QA-SEC-004`](../../../../spec/lastenheft.md), Determinismus, Nebenläufigkeit — nach dem Muster der `ExcludeColumn`-Tests. |
| `internal/bootstrap/wiring.go` (+ Test) | update | `classifyRunError` bildet `ErrTransformationNotApplicable` auf `model.ErrorClassSchema` ab. |
| Aufrufer der gemeinsamen Funktion im Backfill-Pfad (falls vorhanden) | update | reicht bis `backfill-pfad` die leere Regelmenge, ausdrücklich benannt. |
| `internal/application/usecase/backfill/service.go` | update | der Aufrufer im Backfill-Pfad (`blockBuilder.build`) übergibt `nil` als Regelmenge; der Kommentar an der Aufrufstelle nennt die Grenze mit der Adresse `slice-transformationen-backfill-pfad`. |
| `internal/domain/errors/errors.go` | update | vier Sentinels: `ErrInvalidTransformation` (Form), `ErrTransformationTargetIsColumn` (K3-Fall `to` gleich `column`, eigener Sentinel), `ErrTransformationColumnMissing` und `ErrTransformationTargetCollides` (die beiden Gründe der Nichtanwendbarkeit). |
| `internal/domain/model/rowimage_test.go`, `internal/adapters/driven/postgressnapshot/snapshot_test.go` | update | Aufrufer der Signatur nachgezogen; die Byte-Tabelle steht als Paketvariable `rowImageByteCases`, damit derselbe Satz Referenz-Bytes den Vergleich mit einer leeren Regelmenge trägt — Erwartungswerte unverändert. |
| `internal/domain/model/transformation_test.go` | neu | Konstruktor-Invarianten je Eingabe, Regeltyp-Menge, `CheckApplicable`, Wirkung in `BuildRowImage` (Position, Abwesenheit, Ausschluss, Maskierung, leere Regelmenge, Reinheit unter `-race`). |
| `internal/adapters/driving/replication/mapper/transformation_test.go`, `transformation_internal_test.go`, `mapper_bench_test.go` | neu | die Mapper-Tests liegen in eigenen Dateien statt in `mapper_test.go` (dort unverändert): Wirkung auf beide Images, Metadaten, Nichtanwendbarkeit, Live-Reload, Erhalt, Fitness-Function-Tests; der Whitebox-Test hält die Schnappschuss-Zusage (Liste wird neu aufgebaut); der Benchmark misst die Kosten der Prüfung. |
| `internal/bootstrap/heartbeat_internal_test.go` | update | zwei Fälle für `ErrTransformationNotApplicable` in der Sentinel-Tabelle von `classifyRunError`. |

**Festlegungen dieses Slice** (Auslegungen, wo der Plan-Text offen lässt; keine Abweichung von
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)):

- **Der Regelname ist Teil der Domänen-Regel** (`NewRenameColumn(name, column, to)`, Invariante
  nur „nicht leer“): `RemoveTransformation` und das Ersetzen adressieren je Name
  (K1). Alphabet und Länge des Namens prüft weiter der Aufrufer, der ihn entgegennimmt
  (`antragsweg-usecase`, Text `Regelname ist ungültig`) — die Domäne nimmt diese Prüfung nicht vorweg.
- **Die Anwendbarkeit ist eine Domänen-Methode** (`Transformation.CheckApplicable`,
  `SPEC-030` Anwendbarkeit); `mapper.ErrTransformationNotApplicable` wrappt Regelname, Tabelle
  und den Domänen-Grund. Der Backfill-Pfad ruft dieselbe Methode gegen die Snapshot-Spalten
  (Übergabe an `backfill-pfad`).
- **`BuildRowImage` prüft die Anwendbarkeit nicht selbst.** Sie ist Vorbedingung des Aufrufers und im
  Doc-Kommentar benannt; eine nicht anwendbare Regel wirkt in der Funktion nicht. Der Mapper prüft
  je Änderung vor jeder Serialisierung.
- **`SetTransformation` ersetzt eine Regel unter demselben Namen an ihrer Stelle**, statt einen zweiten
  Eintrag anzulegen; K1 bleibt Prüfung des Use Case.
- **Auswertung je Spalte:** die erste Regel, deren Quellspalte die Spalte ist, entscheidet (K2 hält
  je Quellspalte höchstens eine).

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „Signatur und
Aufrufer der gemeinsamen Row-Image-Funktion“, „der Feldsatz von
`TableBinding`“, „die Menge der Sentinels, die `classifyRunError` auf `schema`
abbildet“; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Aufrufer der gemeinsamen Funktion | `git grep -n BuildRowImage` über `internal/*.go`, Parent `1852b2f9` und Diff (Block unten, Zeilen 1–4) | Parent: 32 Trefferzeilen, davon 10 außerhalb der Tests; drei Aufrufstellen im Produktivcode (`Assembler.change` zwei, `blockBuilder.build` eine), neun in Tests (`rowimage_test.go` sieben, `snapshot_test.go` zwei), der Rest Doc-Kommentare (`snapshot.go`, `tablesnapshot.go`, `service.go`, `mapper.go`) und die Definition. Diff: 49 (17 mehr: die neuen Tests), außerhalb der Tests weiter 10. Nichtgefunden: kein weiterer Aufrufer, kein zweiter Bild-Erzeuger im Produktivcode (`json.Marshal` an Row Images: nur `rowimage.go`; `natsstream/publisher.go` und `http/sse.go` marshalen fertige Changes, keine Bilder). | alle drei Aufrufer und alle neun Test-Aufrufer nachgezogen; der Backfill-Aufrufer übergibt `nil`, Kommentar an der Stelle nennt `slice-transformationen-backfill-pfad`. Die Doc-Kommentare beschreiben die Funktion ohne Signatur und bleiben wahr. |
| Beschreibungen der Signatur in Doku | `git grep -n 'BuildRowImage(columns'` über `docs spec harness` ohne `docs/reviews`, `done/`, Baseline (Block unten, Zeile 9) | Parent 3: `docs/plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md` (Zeilen 129, 196: die Signatur mit drei Parametern und „bleibt unverändert“) und `docs/plan/planning/open/slice-transformationen-backfill-pfad.md` (Zeile 172: der Aufruf mit drei Argumenten). Diff: 3, unverändert. | **gemeldet, nicht mitgeändert** (fremde Dateien): die ADR ist `Accepted` und unberührbar (`AGENTS.md` §3.5), ihre Aussage beschreibt den Stand der Entscheidung; der Plan `backfill-pfad` nennt die Signatur als Übergabe — der Planner zieht die Zeile bei der Closure dieses Slice nach (Signatur trägt jetzt den vierten Parameter `rules []Transformation`). |
| Anlage- und Merge-Stellen von Bindungen | `git grep -n 'TableBinding{'` über `internal/*.go` (Block unten, Zeilen 5–8) | Parent: 43 Trefferzeilen, davon 5 außerhalb der Tests: `wiring.go` (`parseTables` mit zwei Zeilen, `activatedTableBindings`, Aktivierungs-Zweig der Antrags-Verarbeitung) und `config_file.go` (`mergeTables`). Diff: 52 (9 mehr, nur Tests), außerhalb der Tests weiter 5. Nichtgefunden: kein Ort, der einen Regelstand aus der Datenbank ableitet (der Feldsatz kennt `Transformations` erst mit diesem Slice: `git grep -c Transformations 1852b2f9 -- internal` druckt keine Datei). | in diesem Slice tragen `AddBinding`-Merge und `setSchemaVersion` den Regelstand mit; die drei Anlagestellen in `wiring.go`/`config_file.go` legen Bindungen ohne Regelstand an (`Transformations` nil) — die Anlage aus der Datenbank trägt `antragsweg-usecase`, so im Plan festgelegt. |
| Feldsatz-Beschreibungen von `ExcludedColumns` | `git grep -n ExcludedColumns` über `internal docs spec harness` ohne `docs/reviews`, `done/`, Baseline (Block unten, Zeilen 10–11) | Parent: 79 Trefferzeilen in 22 Dateien; die Beschreibung des Feldsatzes von `TableBinding` steht nur in `mapper.go` (Doc-Kommentare von `TableBinding`, `AddBinding`, `setSchemaVersion`); die übrigen Fundstellen sind Port, Adapter, Verdrahtung und Tests des Ausschlussstands. Diff: 81 (`mapper.go` eine mehr, `transformation_test.go` eine). Nichtgefunden: keine Doku außerhalb von `mapper.go`, die die Felder der Bindung aufzählt (`docs/user` trägt den Namen nicht). | die drei Doc-Kommentare in `mapper.go` nennen den Regelstand neben `ExcludedColumns` (`TableBinding`, `AddBinding`, `setSchemaVersion`, `Assembler`); die Ausschluss-Ports und ihre Adapter beschreiben keinen Feldsatz und bleiben. |
| Aufzählungen der Sentinels der Klasse `schema` | `git grep -n ErrIncompatibleSchemaChange` über `internal docs spec harness` ohne `docs/reviews`, `done/`, Baseline (Block unten, Zeilen 12–13) | Parent: 29 Trefferzeilen; Aufzählungen der Sentinels der Klasse `schema` stehen an zwei Stellen im Code: `classifyRunError` in `wiring.go` und die Sentinel-Tabelle in `heartbeat_internal_test.go`. Die übrigen: `mapper.go`/`mapper_test.go` (Erzeuger und Tests), [`ADR-0059`](../../adr/0059-spaltenauswahl-mechanismus.md), [`ADR-0063`](../../adr/0063-lh-fa-sch-003-testform-korrektur.md) und [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) (`Accepted`, unberührbar), `harness/README.md` (Zeile `make test-integration`) und `docs/user/e2e-abdeckung.md` (Erzeugnis) beschreiben den Pfad der Spaltenentfernung, keine Aufzählung; `slice-transformationen-e2e-abhilfe.md` ist ein offener Plan mit eigener Adresse. Diff: 32 (drei mehr in `transformation_test.go`). Nichtgefunden: keine Aufzählung in Spec oder Handbuch, die den Sentinel namentlich trägt. | `classifyRunError` und die Sentinel-Tabelle tragen den neuen Sentinel; die Handbuch-Zeile `schema` (§6, `docs/user/benutzerhandbuch.md`) nennt keinen Sentinel und ist die Adresse von `slice-transformationen-betriebsdoku` (Meldung, kein Nachzug hier: keine Betreiber-Oberfläche). |

```suchlauf
1852b2f9 32 -n BuildRowImage -- internal/*.go
diff 49 -n BuildRowImage -- internal/*.go
1852b2f9 10 -n BuildRowImage -- internal/*.go :!internal/**/*_test.go
diff 10 -n BuildRowImage -- internal/*.go :!internal/**/*_test.go
1852b2f9 43 -n TableBinding{ -- internal/*.go
diff 52 -n TableBinding{ -- internal/*.go
1852b2f9 5 -n TableBinding{ -- internal/*.go :!internal/**/*_test.go
diff 5 -n TableBinding{ -- internal/*.go :!internal/**/*_test.go
1852b2f9 3 -n BuildRowImage(columns -- docs spec harness :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
diff 3 -n BuildRowImage(columns -- docs spec harness :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
1852b2f9 79 -n ExcludedColumns -- internal docs spec harness :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
diff 81 -n ExcludedColumns -- internal docs spec harness :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
1852b2f9 29 -n ErrIncompatibleSchemaChange -- internal docs spec harness :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
diff 32 -n ErrIncompatibleSchemaChange -- internal docs spec harness :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
```

**Belege des Laufs (Implementer):**

- Sensoren: `make test` (Race-Detector) Exit 0, 42 Pakete `ok`; `make a-check` Exit 0, `gesamt: 0 Befund(e)`; `make coverage-gate` Exit 0, `coverage-gate: OK — Coverage 83.70% erfüllt Schwelle 80%` (Domäne und `replication/mapper` liegen in der gemessenen Fläche, kein neues Paket); `make gates` ungefiltert in ein Log geschrieben, Exit-Code gesondert gelesen: Exit 0 am Stand `cfca864b` (`coverage-gate: OK — Coverage 83.80% erfüllt Schwelle 80%`, `d-check: 1208 Datei(en) geprüft, 0 Befund(e)`, `commit-traceability: OK — 5 Commit(s)`, `generated-sync: OK`, `gesamt: 0 Befund(e)`); ein früherer Lauf am Stand davor endete mit Exit 2 an `id-unlinked` (nackte Kennung in diesem Plan), die Kennung ist verlinkt. `make suchlauf-nachmessen PLAN=<dieser Plan>` Exit 0, 14 Zeilen stimmen.
- Ausgang von §6 „Byte-Abweichung ohne Regel“: die Byte-Tabelle `rowImageByteCases` (Referenz-Bytes der Backfill-Welle) läuft mit `nil` und mit einer leeren, nicht-nil Regelliste durch dieselbe Erwartung; `git diff -w 1852b2f9 -- internal/domain/model/rowimage_test.go` zeigt keine geänderte Erwartungs-Zeichenkette (nur die Verschiebung der Tabelle in eine Paketvariable und das vierte Argument `nil`).
- Kosten der Prüfung (§6): `BenchmarkAssemblerChange` (zwölf Spalten, `Consume` Begin/Change/Commit, `go test -bench -count 6 -benchtime 2s`, im Toolchain-Container von `make test` ohne `-race`, `make` trägt kein Ziel für einen Go-Benchmark, daher derselbe `docker run` mit anderem Kommando) am Parent-Stand 4777–4958 ns/op, 2444 B/op, 86 allocs/op; am Diff-Stand ohne Regel 5005–5213 ns/op, 2444 B/op, 86 allocs/op; mit zwei Regeln (`BenchmarkAssemblerChangeWithRules`) 5233–5439 ns/op, 2444 B/op, 86 allocs/op. Die Allokationen bleiben gleich; die Laufzeit-Spanne liegt in der Größenordnung von 0,2–0,6 µs je Änderung und ist am selben Host in getrennten Läufen gemessen, nicht gegen Rauschen abgesichert.
- Kandidatenläufe: `git diff --name-only 1852b2f9 -- docs/user/benutzerhandbuch.md` druckt nichts (Handbuch unberührt, keine Versionshistorie-Pflicht); `git diff --name-only 1852b2f9 -- internal/bootstrap/ tools/schema/ internal/adapters/driving/` trifft `mapper.go`, die Mapper-Tests, `wiring.go` (eine `case`-Zeile in `classifyRunError`) und `heartbeat_internal_test.go` — keine neue Umgebungsvariable, keine SQL-Funktion, kein Endpunkt, also keine neue Betreiber-Oberfläche; Slice-Link-Kandidatenlauf auf die berührten Dokumente: 0 Treffer. §3.7-Probe (`git diff --name-only 1852b2f9 -- '*.go' | xargs -r grep -nE 'slice-[0-9]+|welle-[0-9]+|…'`): 0 Treffer; die namentliche Adresse `slice-transformationen-backfill-pfad` an der Aufrufstelle in `service.go` steht als benannte Grenze, Subjekt ist der Run.
- Zusage · mutierte Eingabe · gesehenes Rot (alle Mutationen einzeln gefahren, danach per `cmp` gegen die Sicherungskopie zurückgenommen):
  - `to` höchstens 63 Byte · Grenze auf 64 · `TestNewRenameColumnInvariants` (beide 64-Byte-Fälle); `to` ohne U+0000 · Prüfung entfernt · Fall „Zielname mit U+0000“; `column` ohne U+0000 · Prüfung entfernt · „Spalte mit U+0000“; `column` nicht leer · Prüfung entfernt · „Spalte leer“; `to` nicht leer · Prüfung entfernt · „Zielname leer“; `to` gleich `column` · Prüfung entfernt · „Zielname gleicht der Quellspalte“ und `TestNewRenameColumnTargetIsColumnIsNotAFormViolation`; Regelname nicht leer · Prüfung entfernt · „Regelname leer“.
  - Anwendbarkeit, Spalte fehlt · Prüfung in `CheckApplicable` entfernt · `TestCheckApplicable` (drei Fälle) und `TestConsumeRuleColumnMissingInRelationIsNotApplicable`; Anwendbarkeit, Zielname kollidiert · Prüfung entfernt · `TestCheckApplicable` (Fall Kollision) und beide Fälle von `TestConsumeRuleTargetCollidesWithRelationColumnIsNotApplicable` samt `TestConsumeRuleTargetCollisionAfterCompatibleExtension`.
  - Ausschluss gilt zuerst · (a) Ausschluss-Prüfung in `BuildRowImage` entfernt, (b) Regeln vor dem Ausschluss ausgewertet (der Ausschluss prüft den Zielschlüssel) · (a) `TestBuildRowImageBytes`, `TestBuildRowImageExcludedValueNowhere`, `TestBuildRowImageRenameColumn` und im `mapper`-Paket `TestExcludedColumnIsUnreachableForEveryRuleKind`; (b) `TestBuildRowImageRenameColumn` (Fall „ausgeschlossene Spalte trägt weder …“) und `TestExcludedColumnIsUnreachableForEveryRuleKind`.
  - Wirkung von `rename_column` · Zielname nicht gesetzt / Wert verändert / umbenannter Schlüssel ans Ende gesetzt · `TestBuildRowImageRenameColumn` (alle Wirkungsfälle bzw. die beiden Positionsfälle „mittlere Spalte“ und „erste Spalte“); Abwesenheit · NULL-Wert einer umbenannten Spalte erzeugt einen Zielschlüssel · Fall „NULL bleibt Abwesenheit, kein Zielschlüssel“; leere Regelmenge liefert dieselben Bytes · Auswertung verändert den Wert bei leerer Menge · `TestBuildRowImageBytes` (jeder Fall mit Wert) und `TestBuildRowImageEmptyRuleSetKeepsBytes`.
  - Regeltypen aus der Domänen-Menge aufgezählt · die Menge trägt einen unbekannten Typ · `TestExcludedColumnIsUnreachableForEveryRuleKind` (bricht mit „nicht abgedeckt“ ab).
  - Regelstand am Assembler · `binding.Transformations` in `change` durch `nil` (beide Images) · `TestConsumeRenameColumnAppliesToBothImages` und sechs weitere; nur das Alt-Image ohne Regelstand · `TestConsumeRenameColumnAppliesToBothImages`, `TestExcludedColumnIsUnreachableForEveryRuleKind`; Prüfung nicht gerufen · vier `NotApplicable`-Tests; Prüfung hinter dem Vorrücken der Sequenz · `TestConsumeRuleRemedyRestoresCapture`; `AddBinding` erhält den Regelstand nicht · `TestAddBindingKeepsRuleState`; `setSchemaVersion` verliert den Regelstand · `TestConsumeRuleSurvivesSchemaBump` und `TestConsumeRuleTargetCollisionAfterCompatibleExtension`; `SetTransformation` ohne Sperre · `TestAssemblerTransformationsAreRaceFree` unter `-race` (`WARNING: DATA RACE`).
  - Schnappschuss ohne eigene Sperre · Ersetzen bzw. Anhängen bzw. Entfernen in der übergebenen Liste · `TestTransformationListsAreReplacedNotMutated` (je der zugehörige Fall); Set/Remove auf einer nicht getragenen Bindung · Bindung belebt (Set bzw. Remove) · `TestSetAndRemoveTransformationOnLiveBinding`; Entfernen ignoriert den Namen · `TestSetAndRemoveTransformationOnLiveBinding`; Ersetzen unter demselben Namen hängt an · `TestSetAndRemoveTransformationOnLiveBinding` und der Schnappschuss-Fall „Ersetzen“.
  - Fehlerklasse · `ErrTransformationNotApplicable` aus `classifyRunError` entfernt · beide neuen Fälle von `TestClassifyRunErrorMapsKnownSentinelsToADR0023Classes`; der Grund nicht als `%w` im Fehler · `TestConsumeRuleColumnMissingInRelationIsNotApplicable`, `…TargetCollides…` und `…AfterCompatibleExtension`.
  - Schichtkante · Domäne importiert `internal/application/port/outbound` · `make a-check` Exit 2, `wrong-direction: domain -> ports`.
  - Nicht mutiert (Grund): die Doc-Kommentare (kein Wächter, Leser: Reviewer, Verifier); die Metadaten-Gleichheit (`TestConsumeRenameColumnLeavesMetadataUnchanged`) beruht auf `reflect.DeepEqual` über den ganzen `model.Change`, eine Mutation an `change` würde ein weiteres Feld verändern müssen — nicht einzeln gefahren.

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-backfill-row-image-gemeinsam`
in `done/` liegt (Kopplung K1 der Welle
[welle-backfill-bestand](../done/welle-backfill-bestand.md) §5: die Regelauswertung
hängt an der **einen** gemeinsamen Funktion, nicht an zwei Bild-Erzeugern),
`slice-transformationen-spec-nachzug` und `slice-harness-suchlauf-nachmessen`
in `done/` liegen (der zweite liefert das Nachmess-Werkzeug für das
§3.13-Suchlauf-Feld dieses Plans) und kein anderer Slice in `in-progress/`
liegt (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Domäne,
  Assembler und Fitness-Tests nicht in einem Review tragen — der abtrennbare
  Teil ist der dritte Liefer-Punkt (Eigenschafts-, Determinismus- und
  Nebenläufigkeits-Tests) als eigener Slice mit Start nach diesem.
- `in-progress` → `open` (blockiert): falls die gemeinsame Funktion der
  Backfill-Welle keinen Erweiterungspunkt ohne Signaturbruch für ihre Aufrufer
  erlaubt (Architect-Frage zum Zuschnitt), oder falls `make a-check` die
  Regeltypen in der Domäne ablehnt.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test` (Race-Detector) grün +
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Byte-Abweichung ohne Regel.** Der Regel-Zweig verändert das Image auch
  dann, wenn keine Regel gesetzt ist (Escaping, Schlüsselreihenfolge, Leerraum)
  — jedes erfasste Bild wäre betroffen. *Erwartet, zu belegen durch:* die
  bestehenden Mapper- und Domänen-Tests mit den Referenz-Bytes der
  Backfill-Welle ohne geänderte Erwartung. **Ausgang:** *(bei Closure)*
- **Schlüsselposition nach der Umbenennung.**
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 3 legt fest, dass der umbenannte Schlüssel die Position seiner
  Quellspalte behält; ein Umbau der Schleife könnte ihn ans Ende setzen. Die
  Position ist eine Festlegung am Ausgang des `Assembler`; die Spec sagt am
  Lesepfad Schlüsselmenge und Werte zu, nicht die Reihenfolge
  ([`SPEC-030`](../../../../spec/pflichtenheft.md)).
  *Erwartet, zu belegen durch:* ein Test mit Regel auf einer mittleren Spalte.
  **Ausgang:** *(bei Closure)*
- **Auswertung sieht eine ausgeschlossene Spalte.** Die Struktur der Auswertung
  (gleiche Schleife, Ausschluss zuerst) ist die tragende Zusage von
  [`LH-QA-SEC-004`](../../../../spec/lastenheft.md); eine Umsortierung der
  Schritte bräche sie still. *Erwartet, zu belegen durch:* der Eigenschaftstest
  mit Mutation (Ausschluss-Check entfernen → rot). **Ausgang:** *(bei Closure)*
- **Data Race auf der Regelliste** zwischen Capture-Goroutine und
  Administrations-Goroutine. *Erwartet, zu belegen durch:* der
  Nebenläufigkeits-Test unter `-race`, Muster der `ExcludeColumn`-Tests.
  **Ausgang:** *(bei Closure)*
- **Kosten der Anwendbarkeits-Prüfung im heißen Pfad.** Sie läuft je Change und
  vergleicht Regeln gegen `Relation.Columns`; die Größenordnung (Regeln je
  Tabelle klein) ist angenommen, nicht gemessen. *Erwartet, zu belegen durch:*
  ein `go test -bench` gegen den Parent-Stand oder eine begründete
  Nicht-Messung im Bericht. **Ausgang:** *(bei Closure)*
- **Zwischenzustand im Backfill-Pfad.** Ein Aufrufer im Backfill-Pfad übergibt
  bis `backfill-pfad` die leere Regelmenge; ein Regelstand kann erst mit
  `antragsweg-usecase` überhaupt entstehen, `backfill-pfad` folgt diesem
  unmittelbar (Welle §5). *Erwartet, zu belegen durch:* der Kommentar an der
  Aufrufstelle und die Reihenfolge der Welle. **Ausgang:** *(bei Closure:
  entfallen mit der Closure von `slice-transformationen-backfill-pfad`, dessen
  §2 die Stelle nennt)*
- **Der Ort der Auswertung weicht vom ADR-Wortlaut ab.**
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  nennt „in `rowImage`“, das private `rowImage` des Mappers entfällt mit
  `slice-backfill-row-image-gemeinsam`; die Auswertung sitzt in der gemeinsamen
  Domänen-Funktion (`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab`,
  offen, 1×). Die Auslegung ist in der Kopplung K1 begründet; sie ist im Review
  prüfbar und wird nicht als Entscheidung dargestellt. **Ausgang:** *(bei
  Closure)*

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
`*`/`PGC` (Greenfield); Domäne, Replication-Mapper und Composition Root sind
keine eigenen Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf §3),
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 6×, DoD Punkt
2), `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (offen, 2×,
einschlägig — DoD Punkt 4), `BEO-PGC/slice-chronik-in-code-kommentar`
(verkörpert, 9×, die neuen Kommentare tragen keine Slice-/Wellen-Chronik),
`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` (offen, 1×, einschlägig —
Risiko §6, letzter Punkt), `BEO-PGC/a-check-null-abdeckung` (verkörpert — die
Layer-Globs sind besetzt, `make a-check` prüft die neue Domänen-Datei),
`BEO-PGC/endstufe-unter-eigenem-messgegenstand-unerreichbar` (verkörpert —
neuer Code liegt in der netzlos gemessenen Fläche, Beleg in DoD Punkt 3).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
