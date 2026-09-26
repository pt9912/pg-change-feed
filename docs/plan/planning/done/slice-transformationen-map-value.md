# Slice transformationen-map-value: Regeltyp `map_value` — Wertabbildung in der Domäne, ohne Änderung am Antragsweg und am Wirkort

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-transformationen](../welle-transformationen.md).

**Bezug:** [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) (wertbasierte
Ableitung über die reine Spaltenauswahl hinaus),
[`LH-QA-SEC-004`](../../../../spec/lastenheft.md),
[`LH-FA-DAT-005`](../../../../spec/lastenheft.md),
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Teilfrage 2 und Folgepflicht 4 (jeder weitere Regeltyp danach braucht eine
Folge-ADR).

**Berührte Spec-Stellen:** [`SPEC-019`](../../../../spec/pflichtenheft.md) und
die Regelform (durch `slice-transformationen-spec-nachzug`, dort stehen beide
Regeltypen) — gelesen, nicht geändert.

**Verantwortlich:** Implementer-Agent, 2026-09-26.

**Autor:** Planner-Agent, Welle-Eröffnung
[welle-transformationen](../welle-transformationen.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** `map_value` (`column`, `values` als Objekt `alt → neu` aus
Zeichenketten) wirkt als zweiter Regeltyp: ist der Wert von `column` als
Zeichenkette ein Schlüssel von `values`, steht der zugeordnete Wert im Image,
jeder andere Wert bleibt unverändert, ein abwesender Wert bleibt abwesend. Der
Slice ändert **nur die Domäne** (Regeltyp, Konstruktor-Invarianten,
Parser-Zweig) und die Tests, die die Regeltypen aus der Domänen-Menge
aufzählen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Antragsweg und Wirkort** — der Use Case, das Schema, der `Assembler` und
  die Verdrahtung bleiben unverändert; das ist die Aussage von
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Folgepflicht 4 und wird am Diff belegt (DoD Punkt 2). Muss dort etwas ändern,
  ist das ein Plan-Nachzug, kein stiller Zusatz.
- **Die Spec-Zeile** —
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Folgepflicht 4 nennt „Domäne, Spec-Zeile, Tests“; die Spec führt beide
  Regeltypen bereits seit `spec-nachzug` (Welle §4, Abweichung 4).
- **Ein dritter Regeltyp** (`set_constant` u. ä.) — ein weiterer Typ braucht
  eine Folge-ADR
  ([`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  §Re-Evaluierungs-Trigger), er wird nicht mitgenommen.
- **Umkehrbarkeit** — `map_value` verliert Information, sobald mehrere
  Quellwerte auf denselben Zielwert abgebildet werden; das ist eine benannte
  Konsequenz von
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md),
  keine Eigenschaft, die dieser Slice behebt.

## 2. Definition of Done

- [x] `map_value` wirkt auf beide Images: Zuordnung, nicht zugeordneter Wert
      unverändert, Schlüsselposition der Quellspalte am Ausgang des `Assembler`
      (Festlegung von
      [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
      Teilfrage 3; die Spec sagt am Lesepfad Schlüsselmenge und Werte zu, nicht
      die Reihenfolge, [`SPEC-030`](../../../../spec/pflichtenheft.md)),
      Abwesenheit bleibt (NULL,
      unverändertes TOAST, ausgeschlossene Spalte, fehlendes Bild),
      Determinismus (gleiche Regelmenge und Relation → byte-gleiches Image am
      Ausgang des `Assembler`),
      mehrere Quellwerte mit demselben Zielwert zulässig; die
      Konstruktor-Invarianten entsprechen der Spec (Randfälle: leeres `values`,
      Abbildung auf sich selbst); der Parser lehnt unbekannte Schlüssel ab. *Zu
      belegen durch:* `make test` (Race-Detector), bestehende Tests ohne
      geänderte Erwartung — bis auf die zwei zwangsläufig geänderten
      (`TestTransformationKindsIsAClosedSet`, `TestTransformationSpecBuild`,
      §3) und die drei Fixture-Schalter der Regeltypen, die für den neuen Typ
      einen Fall brauchen (vier Tests, §3). Gemessen (Lauf der Closure,
      Gegenlauf in §3): gegen den Produktivcode färben die Testdateien des
      Parent-Stands genau diese sechs Tests rot; jeder andere bestehende Test
      bleibt grün.
- [x] Ohne Änderung an Antragsweg und Wirkort: `git diff --stat` gegen den
      Parent-Stand mit dem Ausschluss `':!*_test.go'` nennt weder
      `internal/application/usecase/`, `internal/bootstrap/`, `tools/schema/`,
      `internal/adapters/**` noch die Spec — bis auf die Kommentar-Korrektur
      ohne Verhaltensänderung an `internal/adapters/driving/replication/mapper/mapper.go`
      (§3); der Use Case akzeptiert einen `map_value`-Antrag über den Parser der
      Domäne. Die Testdateien mit Diff (Regeltypen-Aufzählungen, der geforderte
      Use-Case-Test, der Store-Test) liegen in allen vier Verzeichnissen der
      Aufzählung (§3, Auslegung). *Zu belegen durch:* der Diff-Stat im Bericht
      und ein Use-Case-Test, der einen `map_value`-Antrag ohne Änderung des
      Use-Case-Codes annimmt.
- [x] Die Fitness Function bleibt vollständig: der Eigenschaftstest (Regeltyp ×
      ausgeschlossene Spalte) erfasst `map_value` über die Domänen-Menge — die
      Schleife über die Menge läuft ohne Ergänzung, der Fixture-Schalter trägt
      den Fall `map_value` (die manuelle Ergänzung, §3) — das Image trägt weder
      Quellschlüssel noch Zielname noch Quellwert noch abgebildeten Wert; der
      Paritätstest des Backfill-Pfads erfasst `map_value` ohne
      Strukturänderung der Schleife. *Zu belegen
      durch:* `make test` und die Mutation, die `map_value` aus der
      Domänen-Menge entfernt — gemessen (Fixrunde, `go test -race` über
      `internal/domain/...`, `mapper`, `backfill`, `bootstrap`): sie färbt allein
      `TestTransformationKindsIsAClosedSet` rot, die Schleifen über die Menge
      laufen dann über einen Typ und bleiben grün; die Erfassung von `map_value`
      trägt die Schleife über die Menge, die Sicherung der Menge der Test mit
      fester Liste (Gegenrichtung, **übernommen** aus dem Review: ein dritter Typ
      in der Menge färbt fünf Tests); `make a-check` grün.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: entfällt — keine Betreiber-Oberfläche; der Regeltyp steht im
      Handbuch-Abschnitt von `slice-transformationen-betriebsdoku`.
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
| `internal/domain/model/transformation.go` (aus `kern-rename`; + Test) | update | zweiter Regeltyp: Konstruktor-Invarianten, Auswertung, Parser-Zweig. Gelieferte Form: die Konstante `TransformationMapValue`, `NewMapValue`, das Feld `values` (die Zuordnung in kanonischer Kodierung: Paare in aufsteigender Ordnung des Schlüssels, je Feld ein Vier-Byte-Längenpräfix), `Values()`, die Auswertung in `applyTransformations` und die Anwendbarkeit ohne Zielname-Prüfung für `map_value` in `CheckApplicable`. |
| `internal/domain/model/transformationspec.go` (Parser-Zweig) | update | (Nachzug, aus der Zeile oben herausgelöst) `ParseTransformationSpec` liest `values` strikt (nicht leeres Objekt, Zeichenketten als Werte), `allowedRuleKeys` nennt `kind`/`column`/`values`, `Build` baut die Regel. |
| `internal/domain/errors/errors.go` (Doc-Kommentare) | update | (Nachzug) die Kommentare von `ErrInvalidTransformation` und `ErrInvalidRuleSpec` nennen `values`; der Block von `ErrInvalidTransformation` trägt eine Kennung (`AGENTS.md` §3.7). |
| `internal/domain/model/transformation_mapvalue_test.go` (neu), `transformation_test.go`, `transformationspec_test.go` | update | (Nachzug) Tests der Domäne für `map_value`. Zwei bestehende Erwartungen ändern sich zwangsläufig: `TestTransformationKindsIsAClosedSet` (die Menge trägt zwei Typen) und `TestTransformationSpecBuild` (der Platzhalter für einen unbekannten Regeltyp ist `nope`, kein Regeltyp). Die Kodierung der Zuordnung trägt Felder bis 70 000 Byte (`TestMapValueEncodingCarriesFieldsUpTo70000Bytes`) und oberhalb von 16 777 216 Byte (`TestMapValueEncodingCarriesFieldsBeyondSixteenMiB`, die Stufe des vierten Längen-Bytes; die Wandzeit der beiden Fälle beträgt unter `-race` im 250-MiB-Lauf unten 0,06 s und 0,05 s). Speicher des zweiten Tests, **gemessen** am Testbinary (`go test -race -c ./internal/domain/model`, dann `docker run --network none --memory=<X> --memory-swap=<X>` mit `-test.run TestMapValueEncodingCarriesFieldsBeyondSixteenMiB`, ein Lauf je Grenze, Closure 2026-09-26): 150 MiB und 200 MiB enden mit Exit 137, 250 MiB und 300 MiB mit Exit 0; die Spitze liegt zwischen 200 MiB und 250 MiB. `make test` läuft ohne Speichergrenze (`docker run --rm --network none …` ohne `--memory`, gelesen im Makefile). Das Kürzen des Längen-Präfixes auf drei Byte färbt den zweiten Test rot (gesehen). |
| `internal/adapters/driving/replication/mapper/` (Eigenschaftstest) | update | Regeltyp-Menge aus der Domäne — `map_value` wird ohne manuelle Ergänzung erfasst. Gelieferte Form: `transformation_test.go` (Fixture-Schalter `ruleFor` trägt den Fall `map_value`, die Erwartung je Regeltyp steht als Schlüssel-Wert-Paar am Fall) und die neue Datei `transformation_mapvalue_test.go` (beide Images, Abwesenheit, Metadaten, Anwendbarkeit, Determinismus, Live-Reload, `-race`). |
| Paritätstest des Backfill-Pfads (Ort aus `backfill-pfad`) | prüfen → update | tabellengetrieben: die Schleife über `model.TransformationKinds()` erfasst `map_value` ohne Strukturänderung; der Fixture-Schalter `parityRule` in `internal/bootstrap/backfill_image_parity_test.go` bekommt den Fall `map_value` (Übergabe aus `backfill-pfad`, siehe unten), die Prüfung der sichtbaren Wirkung je Regeltyp steht am Fall. Der Ausschluss ist an der Eingabeseite gebunden: ist die Spalte der Regel ausgeschlossen, trägt das Bild weder ihren Schlüssel noch ihren Quellwert noch die Wirkung der Regel; die Mutation, die in `BuildRowImage` den Ausschluss für eine Spalte mit `map_value`-Regel aufhebt, färbt den Paritätstest rot (gesehen, drei Fälle `map_value/Regel an …/ausgeschlossen`). |
| `internal/adapters/driving/replication/mapper/mapper.go` (Doc-Kommentar von `ErrTransformationNotApplicable`) | update | (Nachzug) nur Kommentar, kein Verhaltensdiff: der Kommentar bindet den Zielnamen als Grund der Nichtanwendbarkeit an eine Regel mit Zielname (ohne einen Regeltyp beim Namen zu nennen, das Suchlauf-Feld zählt Produktivcode außerhalb der Domäne ohne Regeltypnamen; nicht jede Regel trägt einen Zielnamen) und trägt eine Kennung (`AGENTS.md` §3.7). Die einzige Produktivdatei außerhalb der Domäne im Diff; DoD Punkt 2 nennt sie als Ausnahme. |
| `internal/application/usecase/backfill/transformation_test.go` | update | (Nachzug, Übergabe aus `backfill-pfad`) Fixture-Schalter `ruleFor` bekommt den Fall `map_value`; die Erwartungen der Bild-Tests folgen dem Regeltyp; neue Fälle: Zustandswechsel des Regelstands mit `map_value` und Nichtanwendbarkeit im Run. |
| `internal/adapters/driven/postgresstorage/administrationrequest_test.go` | update | (Nachzug, über den Plan hinaus, nur Test) `TestTableActivationTransformationRulesReadMapValueThroughJsonb`: die Regelform von `map_value` läuft durch die reale `jsonb`-Spalte (Schlüssel-Umordnung, Normalisierung) und wird zur selben Regel gefaltet wie die gebaute; Tier `make test-store`. Grund: der Slice ändert die Faltung des Regelstands, und der `jsonb`-Weg ist der einzige Ort, an dem die Ordnung der Zuordnung nicht vom Aufrufer kommt. |
| `internal/application/usecase/settransformation/service_test.go`, `internal/bootstrap/administration_internal_test.go` | update | (Nachzug, DoD Punkt 2) der Use-Case-Test eines `map_value`-Antrags am unveränderten Use-Case-Code und sein Pendant über die Verdrahtung `applyAdministrationRequest`. |

**Auflösung der Übergabe aus `backfill-pfad` (Entscheidung dieses Slice, vor dem Code):**
`Transformation` bleibt über `==` vergleichbar; die Zuordnung von `map_value` steht als kanonische
Zeichenkette im Feld `values` (erster der beiden benannten Wege). `sameSet[T comparable]` in
`internal/application/usecase/backfill/service.go` bleibt unverändert, der Diff nennt keine
Produktivdatei unter `internal/application/usecase/`. Die Kosten: die Suche einer Zuordnung ist linear in
der Zahl der Paare und dekodiert nicht (`lookupMappedValue`; hergeleitet aus dem
Quelltext). Größenordnung, **übernommen** aus dem Verifikations-Report
`docs/reviews/verify-slice-transformationen-map-value.md` (V-3; dort gemessen mit
einem Wegwerf-Benchmark in einer Kopie, Schlüssel am Ende der Zuordnung als
ungünstigster Fall, 200 Wiederholungen, ohne `-race`, i9-13900H, nicht
nachgemessen): 80 ns bei 10 Paaren, 6,8 µs bei 1000 Paaren und 0,72 ms bei
100 000 Paaren je Wert und Regel. Die Kosten fallen je Wert einer Spalte mit
`map_value`-Regel an, je Zeile eines Backfill-Blocks und je Change im
Erfassungspfad. Die Spec bindet die Zahl der Paare nicht nach oben; ob sie eine
Obergrenze braucht, ist eine Spec-Frage (Adresse: `welle-transformationen` §5,
„Fragen für den nächsten Architect-Zug“, Punkt (d)); die Betreiber-Aussage
trägt `slice-transformationen-betriebsdoku` (§2, Übergabe-Block). **Auslegung von DoD Punkt 2:** die
Ausnahme „Testdateien, die die Regeltypen aufzählen“ gilt für jedes Verzeichnis der Aufzählung, nicht nur
für `internal/adapters/**` — `internal/application/usecase/backfill/transformation_test.go`,
`internal/bootstrap/backfill_image_parity_test.go` (Fixture-Schalter) sowie
`internal/application/usecase/settransformation/service_test.go` und
`internal/bootstrap/administration_internal_test.go` (der geforderte Use-Case-Test) sind Testdateien mit
Diff; kein Produktivcode dort. Der Store-Test `administrationrequest_test.go` (`internal/adapters/**`)
zählt keine Regeltypen auf und steht als Test über den Plan hinaus in der Tabelle. Der Wortlaut von DoD
Punkt 1 und 2 trägt diese Auslegung selbst (Herkunft: Finding F-2 in
`docs/reviews/review-slice-transformationen-map-value.md`, V-1 in
`docs/reviews/verify-slice-transformationen-map-value.md`): Punkt 1 nennt die zwei zwangsläufig geänderten
Erwartungen und die drei Fixture-Schalter, Punkt 2 den Ausschluss `':!*_test.go'` und die Kommentar-Korrektur in `mapper.go`.
**Gegenlauf der bestehenden Tests (gemessen, Closure 2026-09-26):** eine Kopie von `HEAD`
(`git archive`), in der die acht geänderten Testdateien auf ihren Stand von `e5a11979` zurückgesetzt und die
zwei neuen Testdateien entfernt sind; `go test -race` im Toolchain-Image ohne Netz über `./internal/domain/...`,
`./internal/adapters/driving/replication/mapper/`, `./internal/application/usecase/backfill/`,
`./internal/application/usecase/settransformation/` und `./internal/bootstrap/`; Exit 1, sechs Zeilen
`--- FAIL`: `TestTransformationKindsIsAClosedSet`, `TestTransformationSpecBuild` und die vier Abbrüche
`Regeltyp "map_value" ohne Fall in diesem Test` (`TestExcludedColumnIsUnreachableForEveryRuleKind`,
`TestExecuteBuildsImagesWithTheRuleSet`, `TestExecuteRulesNeverLeakExcludedColumns`,
`TestBackfillAndWALImagesAreByteEqualWithRules`); keine Panik, `settransformation` grün. **Beleg von
DoD Punkt 2:** der Diff-Stat im Bericht, gemessen mit
`git diff --stat e5a11979 -- internal/application/usecase internal/bootstrap tools/schema internal/adapters spec ':!*_test.go'`;
er nennt genau eine Datei, `internal/adapters/driving/replication/mapper/mapper.go` (nur Kommentar), keine
Zeile in `internal/application/usecase/`, `internal/bootstrap/`, `tools/schema/` oder `spec/`.

**Nicht realisiert (Grund):** der Doc-Kommentar von `BuildRowImage` (`internal/domain/model/rowimage.go`)
trägt keinen Satz zur Wertabbildung: sein Block trägt bereits mehrere Kennungen (Kandidat von
`make kommentar-kennungen`: der Block in `rowimage.go`, Zeilen 11 bis 45, neun Kennungen, gemessen in der Closure mit
`make kommentar-kennungen PATHS=internal/domain/model/rowimage.go`), und die Bereinigung liegt bei
`slice-code-kommentare-bereinigung`; der Kommentar bleibt wahr (er beschreibt die Position umbenannter
Schlüssel, `applyTransformations` trägt die Beschreibung der Wertabbildung).

**Übergabe aus `slice-transformationen-antragsweg-usecase`** (gemeldet, kein zusätzlicher
Umfang; Herkunft: Plan des Slice §6 und Review-Report
`review-slice-transformationen-antragsweg-usecase` Finding F-8, gelesen am Stand `2c22334f`):

- **Zwischenzustand der Regeltyp-Menge endet mit diesem Slice.** Der Use Case prüft gegen die
  Regeltyp-Menge der Domäne; bis zu diesem Slice besteht sie aus `rename_column`, ein
  `map_value`-Antrag endet `failed` mit `unbekannter Regeltyp`, während
  [`SPEC-019`](../../../../spec/pflichtenheft.md) ihn als Regeltyp führt. Das Ende belegt der
  Use-Case-Test des zweiten DoD-Punkts (ein `map_value`-Antrag ohne Änderung des Use-Case-Codes).
- **Rückfall auf einen Binärstand vor diesem Slice.** Eine `applied`-Zeile, die die Faltung
  (`model.FoldTransformations`) nicht mehr in eine Regel führt, endet als Fehler der Klasse
  `internal` und hält Prozessstart und jeden Regel-Antrag der ganzen Quelle an (bewusst: der
  Stand wird nie um eine Zeile verkürzt). Erstmals erreichbar mit diesem Slice: eine vermerkte
  `map_value`-Zeile ist für einen Binärstand ohne den Regeltyp nicht lesbar. Der Slice nennt die
  Grenze im Bericht und übergibt sie mit dem Wortlaut „Rückfall auf einen älteren Binärstand
  nach vermerkten `map_value`-Regeln hält die Quelle an“ an
  `slice-transformationen-betriebsdoku` (dessen §2, Abschnitt zur Dauerhaftigkeit), oder
  belegt, dass sie an anderer Stelle getragen ist. Beleg des heutigen Verhaltens: Review-Report
  F-8 (hergeleitet aus dem Quelltext, nicht erprobt).

**Übergabe aus `slice-transformationen-backfill-pfad`** (gemeldet, kein zusätzlicher
Umfang; Herkunft: Review-Report Finding F-3 und Verifikations-Report V-2 dieses Slice, gelesen
am Stand `bdff5a54`; Frist der Meldung: der Start dieses Slice). Die Fundstellen sind gemessen
mit `git grep -n sameSet -- internal` (sechs Trefferzeilen in
`internal/application/usecase/backfill/service.go`: Kommentar, Definition und vier Aufrufe) und
`git grep -n 'ohne Fall in diesem Test' -- internal` (zwei Treffer):

- **Der Run vergleicht Regeln als Menge über `==`.** `sameSet[T comparable]` in
  `internal/application/usecase/backfill/service.go` vergleicht den Regelstand, den `copyBlocks`
  je Block liest, mit dem Stand zu Beginn des Runs; die Funktion wird mit `model.Transformation`
  instanziiert, dessen vier Felder Zeichenketten sind (`transformation.go`), und ihr Doc-Kommentar
  nennt den Typ „über alle seine Felder vergleichbar“. Ein Regeltyp mit dem Objekt `values`
  (Map oder Slice als Feld) macht `Transformation` unvergleichbar; der Bau bricht dann am
  Übersetzer — sichtbar, nicht still.
- **Der Konflikt mit DoD Punkt 2 und §4 ist real.** DoD Punkt 2 verlangt einen Diff ohne
  `internal/application/usecase/`, §4 führt eine Änderung dort als Rückführung `in-progress` →
  `open`; der Vergleich des Runs liegt in `internal/application/usecase/backfill/`. Der Slice
  löst den Konflikt in seinem eigenen Plan, bevor Code entsteht; die Lösung bleibt dem Slice
  überlassen. Zwei Wege, nicht entschieden: `Transformation` bleibt vergleichbar (etwa `values`
  als kanonische Zeichenkette im Feld — kein Diff im Use Case) oder der Run vergleicht über
  eine Kennung bzw. Kanonisierung, die die Domäne liefert (Diff im Use Case, dann mit
  Ausnahme in DoD Punkt 2 und Plan-Nachzug).
- **Die zwei Fixture-Schalter der Regeltypen liegen in den von DoD Punkt 2 ausgenommenen
  Verzeichnissen.** `parityRule` (`internal/bootstrap/backfill_image_parity_test.go`, der
  Paritätstest des Backfill-Pfads, Ort aus dieser Closure) und `ruleFor`
  (`internal/application/usecase/backfill/transformation_test.go`, der Eigenschaftstest im
  Run) zählen die Regeltypen über `model.TransformationKinds()` auf und brechen mit
  `Regeltyp … ohne Fall in diesem Test` ab, sobald ein Typ dort steht, für den der Schalter
  keinen Fall trägt. Beide brauchen einen `map_value`-Fall; die Zeile „Paritätstest des
  Backfill-Pfads … prüfen“ in der Tabelle oben ist damit ein Test-Diff in `internal/bootstrap/`
  und `internal/application/usecase/backfill/`, den DoD Punkt 2 (Klammer „außer Testdateien,
  die die Regeltypen aufzählen“ nur für `internal/adapters/**`) als Ausnahme nennen muss; die
  Aussage „ohne Strukturänderung“ des Paritätstests (DoD Punkt 3) gilt für die Schleife über die
  Domänen-Menge, nicht für den Fixture-Schalter.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Menge der
Regeltypen“ (ein Typ → zwei); beide Stände gemessen):**

Die Befehle stehen im Block unten (eine Zeile je Messung, `<Stand> <Soll> <git grep>`); Parent ist
`e5a11979` (Stand vor dem ersten Code-Commit dieses Slice), `diff` der Arbeitsbaum. Der Suchraum ist der
ganze Baum ohne `docs/reviews/**`, die Records unter `done/` und `.harness/baseline/**`; jede
Einschränkung steht im Muster.

| Träger | Messung (Zeilen im Block) | Befund | Behandlung |
|---|---|---|---|
| Aufzählungen der Regeltypen im Code | 1–6, 13–14 und 21–22 | **Gefunden.** `rename_column` im Code ohne Tests: Parent 10, Diff 11 (Domäne, dazu zwei Test-Runner-Skripte unter `tools/harness/` als Testdaten; plus eine Nennung im Doc-Kommentar von `CheckApplicable`); `TransformationKinds` im Code ohne Tests: 5 und 5 (Definition und Doc-Kommentare der Domäne, kein Aufrufer); `TransformationRenameColumn` im Code ohne Tests: Parent 9, Diff 10; `TransformationMapValue` im Code: Parent 0, Diff 9. **Nichtgefunden:** kein Produktivcode außerhalb der Domäne (`internal/application`, `internal/bootstrap`, `internal/adapters`, `tools/schema`) nennt einen Regeltyp beim Namen — Parent 0, Diff 0 (Zeilen 13–14; Risiko „Antragsweg nicht generisch“). | jeder der vier Schalter der Domäne über den Regeltyp (`applyTransformations`, `ParseTransformationSpec`, `allowedRuleKeys`, `TransformationSpec.Build`) trägt beide Zweige; kein Träger außerhalb der Domäne zieht nach |
| Aufzählungen der Regeltypen in Docs | 9–12 (Zählwort, Hedge, Rückfall-Grenze) und 15–16 (Handbuch) | **Gefunden.** Sätze mit „beide/zwei Regeltypen“, „ein Regeltyp“ oder „nur `rename_column`“ im Suchraum ohne diesen Plan: 9 und 9 (Spec, ADR, Welle und Pläne anderer Slices); jeder nennt beide Typen, eine künftige Folge oder einen hypothetischen dritten Typ. `docs/user`: eine Nennung von `rename_column` (Zeile der Abdeckungstabelle, Testdatum eines Backfill-Belegs, ein Erzeugnis des Runners), Parent 1, Diff 1. **Nichtgefunden:** kein Satz, der die Menge als einen Typ beschreibt; das Benutzerhandbuch nennt keinen Regeltyp (0). | Spec trägt beide Typen bereits; Handbuch-Träger an `betriebsdoku` gemeldet (Handbuch-Abschnitt, Beispiele `rename_column`/`map_value`); die Rückfall-Grenze (Übergabe aus `antragsweg-usecase`) steht bereits in `slice-transformationen-betriebsdoku` §2 (Zeilen 11–12: eine Fundstelle, beide Stände) und wird dort nicht doppelt geführt |
| Tests, die eine feste Typ-Liste führen | 7–8 und 17–20 | **Gefunden.** `TransformationKinds` in Tests: Parent 9, Diff 9; die drei Fixture-Schalter „Regeltyp %q …“ (mapper, backfill, bootstrap-Parität): 3 und 3, je mit neuem Fall `map_value`; ein Test, der `map_value` als Platzhalter für einen unbekannten Regeltyp führte (`transformationspec_test.go`): Parent 1, Diff 0. **Nichtgefunden:** kein weiterer Test führt eine feste Regeltyp-Liste. | die drei Schalter tragen den Fall; die Erwartung von `TestTransformationKindsIsAClosedSet` und der Platzhalter von `TestTransformationSpecBuild` sind angepasst (§3-Tabelle) |
| Träger, die den Zielnamen beschreiben (Beschreibung statt Symbolname) | 23–24 (Code außerhalb der Domäne) und 25–26 (Spec, Handbuch) | **Gefunden.** „Zielname“ in Nicht-Test-Code außerhalb `internal/domain`: Parent 8, Diff 9 (Parent: vier Zeilen im `mapper`-Paket, eine im Backfill-Use-Case, zwei im Antrags-Use-Case, eine in der Verdrahtung; Diff: die Definition trägt das Wort in zwei Zeilen statt einer). Davon ist genau eine Stelle die Definition der Nichtanwendbarkeit — der Doc-Kommentar von `ErrTransformationNotApplicable` in `mapper.go` —; am Parent führte sie den Zielnamen als Grund jeder Regel, am Diff bindet sie ihn an eine Regel mit Zielname (§3-Zeile oben). Die übrigen sieben nennen K3 oder die Kollision zweier Regeln mit gleichem Zielnamen und bleiben für `rename_column` wahr. `spec` und `docs/user`: Parent und Diff gleich (die Spec trägt beide Typen). **Nichtgefunden:** keine weitere Stelle, die den Zielnamen als Eigenschaft jeder Regel beschreibt. | `mapper.go`-Kommentar nachgezogen; die übrigen Zeilen unverändert |

```suchlauf
e5a11979 10 -n rename_column -- internal tools ':!*_test.go'
diff 11 -n rename_column -- internal tools ':!*_test.go'
e5a11979 5 -n TransformationKinds -- internal ':!*_test.go'
diff 5 -n TransformationKinds -- internal ':!*_test.go'
e5a11979 9 -n TransformationRenameColumn -- internal ':!*_test.go'
diff 10 -n TransformationRenameColumn -- internal ':!*_test.go'
e5a11979 9 -n TransformationKinds -- 'internal/*_test.go'
diff 9 -n TransformationKinds -- 'internal/*_test.go'
e5a11979 9 -n -i -E 'beide Regeltypen|zwei Regeltypen|ein Regeltyp|einziger Regeltyp|nur .rename_column' -- docs spec harness README.md AGENTS.md ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 9 -n -i -E 'beide Regeltypen|zwei Regeltypen|ein Regeltyp|einziger Regeltyp|nur .rename_column' -- docs spec harness README.md AGENTS.md ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
e5a11979 1 -n -F 'vor `map-value`' -- docs/plan/planning
diff 1 -n -F 'vor `map-value`' -- docs/plan/planning
e5a11979 0 -n -E 'TransformationRenameColumn|TransformationMapValue|rename_column|map_value' -- internal/application internal/bootstrap internal/adapters tools/schema ':!*_test.go'
diff 0 -n -E 'TransformationRenameColumn|TransformationMapValue|rename_column|map_value' -- internal/application internal/bootstrap internal/adapters tools/schema ':!*_test.go'
e5a11979 1 -n -i -E 'regeltyp|map_value|rename_column' -- docs/user
diff 1 -n -i -E 'regeltyp|map_value|rename_column' -- docs/user
e5a11979 3 -n 'Regeltyp %q' -- internal
diff 3 -n 'Regeltyp %q' -- internal
e5a11979 1 -n -F 'kind: "map_value"' -- internal
diff 0 -n -F 'kind: "map_value"' -- internal
e5a11979 0 -n TransformationMapValue -- internal ':!*_test.go'
diff 9 -n TransformationMapValue -- internal ':!*_test.go'
e5a11979 8 -n -i Zielname -- internal cmd ':!*_test.go' ':!internal/domain'
diff 9 -n -i Zielname -- internal cmd ':!*_test.go' ':!internal/domain'
e5a11979 14 -n -i Zielname -- spec docs/user
diff 14 -n -i Zielname -- spec docs/user
```

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-transformationen-backfill-pfad`
in `done/` liegt (die Bindung der Erzeugungspfade steht, bevor der zweite Typ
hinzukommt), `slice-code-kommentare-kennungen` in `done/` liegt (Kante: der
Implementer dieses und jedes folgenden Slices der Welle läuft Schritt 20 mit dem
Werkzeug `make kommentar-kennungen`, und die Regel in
[`AGENTS.md`](../../../../AGENTS.md) §3.7 steht) und kein anderer Slice in
`in-progress/` liegt (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht erwartet — ein
  zweiter Regeltyp in der Domäne mit Tests; wächst der Zug um Änderungen an
  Antragsweg oder Wirkort, ist das der Rückführungs-Fall des zweiten
  Liefer-Punkts.
- `in-progress` → `open` (blockiert): falls der Use Case oder der `Assembler`
  für `map_value` geändert werden müssten (dann Architect-Frage: die Aussage
  „ohne Änderung an Antragsweg oder Wirkort“ von
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Folgepflicht 4 trüge nicht) oder falls die Spec die Randfälle anders
  festgelegt hat, als die Domäne sie tragen kann.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test` (Race-Detector) grün +
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Der Antragsweg ist nicht generisch über die Regeltypen** — ein Zweig im Use
  Case, der `rename_column` beim Namen nennt, machte den zweiten Typ zu einer
  Antragsweg-Änderung. *Erwartet, zu belegen durch:* der Diff-Stat und der
  Use-Case-Test aus DoD Punkt 2. **Ausgang:** *entfallen* — der Stat
  `git diff --stat e5a11979 -- internal/application/usecase internal/bootstrap tools/schema internal/adapters spec ':!*_test.go'`
  nennt genau eine Datei, `internal/adapters/driving/replication/mapper/mapper.go`
  (9 Einfügungen, 8 Löschungen, nur Kommentarzeilen; gemessen in der Closure);
  die Suchlauf-Zeilen 13 und 14 zählen 0 und 0 (kein Regeltypname im Produktivcode
  außerhalb der Domäne, Parent und Diff); der Use-Case-Test
  `TestSetTransformationAcceptsAndChecksMapValueThroughTheDomain` nimmt einen
  `map_value`-Antrag ohne Änderung des Use-Case-Codes an und färbt rot, wenn
  `allowedRuleKeys` `map_value` einen fremden Schlüssel erlaubt oder K2 entfällt
  (Mutationen der Verifikation, **übernommen**: M13 und M3 im Report
  `docs/reviews/verify-slice-transformationen-map-value.md` §4).
- **K3 gilt nur für Umbenennungen** (Zielnamen); `map_value` trägt keinen
  Zielnamen — eine Prüfung, die ihn erwartet, scheitert am zweiten Typ.
  *Erwartet, zu belegen durch:* der Test eines `map_value`-Antrags gegen die
  Konfliktfreiheit (K2 greift: eine Quellspalte trägt höchstens eine
  Spaltenregel). **Ausgang:** *entfallen* — `TestCheckConflictsMapValue`
  (`internal/domain/model/transformation_mapvalue_test.go`) trägt die Fälle „kein
  K3 gegen eine Umbenennung“ und „kein K3 gegen eine zweite Wertabbildung auf
  anderer Spalte“ (gelesen in der Closure) und die K2-Fälle; die Wache
  `s.to != ""` in `CheckConflicts` und die Zielnamen-Prüfung für `map_value` in
  `CheckApplicable` färben je rot (Mutationen M4 und M5 der Verifikation,
  **übernommen**; M5 ist in der Produktion ein äquivalenter Mutant, weil er nur
  über eine Spalte ohne Namen rot wird, Verifikation V-4).
- **Ein Wert, der als Zeichenkette kein Schlüssel ist, wird still abgebildet**
  (Normalisierung, Groß-/Kleinschreibung). *Erwartet, zu belegen durch:* Test
  der exakten Zeichenketten-Gleichheit; die Festlegung steht in der Spec.
  **Ausgang:** *entfallen* — `TestBuildRowImageMapValue` trägt die Fälle
  „Groß-/Kleinschreibung zählt“ und „Präfix eines Schlüssels ist kein
  Schlüssel“ (gelesen in der Closure); die Suche per `HasPrefix` und per
  `EqualFold` in `lookupMappedValue` färbt je rot (M11 und M16 der
  Verifikation, **übernommen**); der Review las `é` und `e` mit U+0301 als
  getrennte Schlüssel (Negativbefund im Review, **übernommen**).
- **Informationsverlust ist nicht sichtbar** (mehrere Quellwerte, ein
  Zielwert): kein Fehler, keine Warnung. Das ist eine benannte Konsequenz von
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md);
  die Betreiber-Aussage steht im Handbuch-Abschnitt von `betriebsdoku`.
  **Ausgang:** *weiter offen* — die Konsequenz gilt (mehrere Quellwerte auf einen
  Zielwert sind zulässig und melden nichts; Tests der Domäne) und ist bis zum
  Handbuch nicht als Betreiber-Aussage getragen. Adresse: `slice-transformationen-betriebsdoku`
  (`open/`), dessen §2 den Informationsverlust bei `map_value` als Stichwort nennt
  und ihn samt der Größenordnung der linearen Suche als committeten
  Übergabe-Block führt (Träger-Nachzug dieser Closure). Kein
  Register-Verzeichnis: der Träger ist ein offener Slice derselben Welle, deren
  Closure die Paarung prüft (Vorbild: das Risiko „Kosten der Lesung je Block“ in
  `slice-transformationen-backfill-pfad`).

## 7. Closure-Notiz

- **Was hat funktioniert:** (1) Die Übergabe-Blöcke der Vorgänger wurden vor dem Code
  entschieden und trugen: `Transformation` blieb über `==` vergleichbar (die Zuordnung steht als
  kanonische Zeichenkette in `values`, §3), `sameSet` und der Use Case blieben unberührt. Der Stat nach
  DoD Punkt 2 nennt genau eine Produktivdatei außerhalb der Domäne, `mapper.go` (9 Einfügungen, 8
  Löschungen, nur Kommentar; gemessen in der Closure), die Suchlauf-Zeilen 13 und 14 zählen 0 und 0
  (Parent und Diff). (2) Die Leser-Kette lief mit einer Fixrunde ohne Rückführung: der Review nennt 0
  HIGH, 1 MEDIUM, 3 LOW und 3 INFO (übernommen aus dem Report) und fand F-1 durch Mutation der
  Eingabeseite; die Verifikation nennt 0 HIGH, 0 MEDIUM, 2 LOW und 5 INFO (übernommen) und fuhr 21
  Mutationen, alle rot (übernommen; darunter M12 und M18 zu DoD Punkt 3, M7 zum Ausschluss, M1 zur
  Kodierung). (3) Die Closure maß nach, was der Bericht nur nannte: den Gegenlauf der bestehenden
  Tests (sechs rote, §3), den Speicher des 16-MiB-Tests (150 und 200 MiB Exit 137, 250 und 300 MiB
  Exit 0, §3), den Stat nach DoD Punkt 2 (eine Datei) und den Kandidaten-Block in `rowimage.go`
  (`make kommentar-kennungen`: Zeilen 11 bis 45, neun Kennungen).
- **Was ging anders als geplant:** (1) **Der DoD-Wortlaut zählte falsch.** DoD Punkt 1 nannte „die zwei
  zwangsläufig geänderten“ Erwartungen; gegen den Produktivcode färben die Testdateien des Parent-Stands
  sechs Tests rot, davon vier Abbrüche der Fixture-Schalter (Verifikation V-1; in der Closure
  nachgemessen: Exit 1, sechs Zeilen `--- FAIL`). Die Übergabe aus `backfill-pfad` verlangte, DoD Punkt 2
  müsse die Ausnahme selbst nennen; der Implementer löste sie als Auslegung in §3, der Wortlaut folgte
  erst nach Review F-2 und Verifikation V-1. DoD Punkt 1 und 3 sind in der Closure im Ist-Ton
  nachgezogen. (2) **Ein Testkommentar sagte zu, was der Test nicht band** (Review F-1, MEDIUM): der
  Paritätstest band den Ausschluss für `map_value` an keiner Eingabe; die Fixrunde band ihn (Commit
  `077e7531`), die Verifikation sah die Mutation rot. (3) **Ein Test behauptete mehr als er trieb**
  (F-3): „jede Stufe“ der Längen-Bytes bei einer Übung bis 70 000 Byte; die Fixrunde ergänzte den
  16-MiB-Test (Commit `901274c9`). Er ist der speicherhungrigste Test des Pakets (Spitze zwischen 200 und
  250 MiB, gemessen); `make test` setzt keine Speichergrenze (Makefile gelesen), für CI ist der Wert nicht
  gemessen; das Ereignis, das die Kenntnis abholt, ist ein Exit 137 im Lauf von `internal/domain/model`.
  (4) **Die Speicher-Zahl hatte keinen Lauf** (V-2): sie steht seit der Closure mit Befehl und Lauf in §3.
  (5) **Umfang über den Plan hinaus, im Plan nachgetragen:** der Store-Test zur `jsonb`-Spalte und der
  zweite Kodierungstest. (6) **Umfang:** 15 Commits, 18 Dateien, +2345/−361 einschließlich der zwei
  Lifecycle-Moves, des Review- und des Verifikations-Reports; die Produktivdateien: vier Go-Dateien mit
  177 Einfügungen und 24 Löschungen, die Tests zehn Dateien mit 1275 und 67 (gemessen mit
  `git diff --shortstat e5a11979..8b9b9c5b` in der Closure). (7) **Die Coverage streut:** 85,50 % und 85,40 %
  in den zwei `make gates`-Läufen der Verifikation (übernommen), 85,30 % im Lauf der Closure vor dem
  Inhalts-Commit (gedruckt: „coverage-gate: OK — Coverage 85.30% erfüllt Schwelle 80%“); Schwelle 80 % —
  der Beleg je eines Laufs, kein Ist-Stand.
- **Steering-Loop-Eintrag (Lerneintrag):** *(a) Neuer Sensor:* keiner gebaut. Das Messverfahren, das den
  DoD-Zählfehler fand, der **Gegenlauf** (Testdateien des Parent-Stands gegen den Produktivcode von
  `HEAD`), steht als Befehl in §3; ein Sensor ist ausgeschlossen: ob ein Zählwort im DoD zu einem
  Testlauf passt, ist eine Messhandlung des Verifiers (`AGENTS.md` §3.12 Instanz B), keine
  Dateieigenschaft. *(b) Geschärfte Regel:* keine Schärfung durch diese Closure. Die verkörperten
  Regeln haben gewirkt: der Reviewer fand F-1 (Zusage ohne Bindung an ihre Eingabeseite), F-3 (Testname
  breiter als die Messung), F-4 (Träger-Suchmuster ohne die bewegte Beschreibung) und F-5 (Beleg trägt
  den Satz nicht), der Verifier V-1 und V-2 (Zählwort und Zahl ohne Lauf) — alle vor dem Merge, keines
  über einem Sensor. Die Anwendung für den Planner als Verfasser: eine DoD-Zeile, die eine **Zahl über
  bestehende Tests** nennt, steht als Erwartung, bis der Gegenlauf gefahren ist (Instanz B, keine neue
  Regel). *(c) Benannte Spec-Lücke:* [`SPEC-030`](../../../../spec/pflichtenheft.md) bindet die Zahl der
  Paare von `values` nicht nach oben, die Suche einer Zuordnung ist linear in dieser Zahl (Größenordnung
  in §3, übernommen); ob die Spec eine Obergrenze führt, ist nicht entschieden. Adresse: `welle-transformationen`
  §5, Fragen für den nächsten Architect-Zug, Punkt (d) (ein Architect-Zug zur Spec, kein stilles
  Festlegen); die Betreiber-Aussage trägt der Übergabe-Block in §2 von
  `slice-transformationen-betriebsdoku`. *(d) Benannte Regelgrenze für den Architect:* das Anhängen von
  Text per Umleitung (`cat >>`) ist weder ein Guard-Treffer noch in [`AGENTS.md`](../../../../AGENTS.md)
  §3.1 ausdrücklich geregelt (Review F-7, Verifikation V-5); Adresse: `welle-transformationen` §5,
  Punkt (c), und `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` (`state.md`). *(e) Gelernt, bei
  1× keine Regel:* die Fitness Function von [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  hat zwei Hälften. Die Erfassung eines Typs trägt die Schleife über die Domänen-Menge, die Sicherung der
  Menge ein Test mit fester Liste. Die Mutation „Typ aus der Menge streichen“ belegt deshalb nur die
  Sicherung (sie färbt allein `TestTransformationKindsIsAClosedSet`; Review F-5, Verifikation M12 —
  DoD Punkt 3 trägt den Wortlaut); „Typ hinzufügen“ färbt fünf Tests (Verifikation M18, **übernommen**).
  Ein Typ, der als Konstante existiert und nicht in die Menge eingetragen wird, ist von den Schleifen
  und vom Listen-Test nicht erfasst — **hergeleitet** aus der Form der Tests, nicht gefahren; Adresse: die Folge-ADR eines
  dritten Regeltyps ([`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Folgepflicht 4) nennt beide Hälften.
- **Beobachtungs-Register (`../observations/`):** je Anfall eine Datei
  `evidence/slice-transformationen-map-value.md`, Zähler = Zahl der Dateien (gemessen mit
  `ls evidence | wc -l` am Stand dieser Closure). *Neue Belege:*
  `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` **18×** (F-1 MEDIUM, daher Datei trotz Deckel;
  verkörpert); `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` **10×** (V-1 LOW; der
  Eintrag erreicht mit dieser Datei die Zahl, ab der der Deckel gilt, das `state.md` nennt „Deckel bei
  10×“; verkörpert); `BEO-PGC/test-name-behauptet-mehr-als-der-test-treibt` **2×** (F-3 LOW; offen,
  unter der Schwelle); `BEO-PGC/gemeldete-ungenauigkeit-ohne-traeger` **2×** (F-6/V-3 INFO, die Form
  „benannte Grenze ohne Adresse“; offen, unter der Schwelle, Adresse im `state.md`);
  `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel` **5×** (F-7/V-5/V-6 INFO, eine
  **Regelgrenze**: `cat >>`, ein geblockter `sed -i` des Verifiers, ein `cd <Repo> && python3
  --version` des Planners, das den Guard an der in `MR-003` benannten Grenze passierte; verkörpert).
  *Deckel-Fälle ohne Datei, Finding-Kennung hier* (verkörpert ab 10×, vor dem Merge von Reviewer bzw.
  Verifier gefunden, Schwere ≤ LOW, bekannter Träger-Typ): F-2 (LOW, `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`,
  Deckel bei 10×: DoD-Wortlaut gegen die Auslegung im selben Plan), F-4 (LOW,
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger`, Deckel bei 32×: der Doc-Kommentar von
  `ErrTransformationNotApplicable`, ein Träger, den das Suchlauf-Feld mit dem Symbolnamen nicht fand),
  F-5 (INFO, `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`, Deckel bei 14×: der Beleg „die Tests
  färben sich rot“ trug seinen Satz nicht, gemessen färbt ein Test), V-2 (LOW,
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`, Deckel bei 23×: die Speicher-Zahl ohne Lauf).
  *`state.md` fortgeschrieben:* die fünf Einträge oben mit neuem Zähler; beim Eintrag zur
  Regelgrenze trägt es zusätzlich den Ausgang der Neubewertung (die Frage an den Nutzer, unten).
  *Kein Register-Anfall:* V-4 (ein äquivalenter Mutant, im Test-Kommentar benannt), V-7 (übernommene
  Messungen, im Report benannt). *Lese-Schritt der Closure von `welle-transformationen`:* aus diesem
  Slice erreicht neu kein Eintrag die Schwelle ohne Ausgang (die zwei Einträge unter der Schwelle stehen
  bei 2×); kein Slice ist wegen einer Beobachtung dieses Slice fällig.
- **Folge-Slices:** keine angelegt. Nächster Slice der Welle nach der Tabellenreihenfolge:
  `slice-transformationen-e2e-wirkung` (`open/`); sein Start-Trigger (§4 dort) ist mit dem Move dieses
  Slice erfüllt (`slice-transformationen-map-value` liegt in `done/`, `slice-harness-fmt-check` liegt in
  `done/`, in `in-progress/` liegt nur die Roadmap). Übergaben mit Adresse (gemeldet, Frist: diese
  Closure, gezogen): `slice-transformationen-betriebsdoku` (`open/`) — der Übergabe-Block in §2
  (Informationsverlust, Vergleich, Größenordnung der Suche, Rückfall-Grenze mit dem erprobten ersten
  Schritt); `welle-transformationen` — §5, Fragen für den nächsten Architect-Zug, Punkte (c) und (d);
  `slice-code-kommentare-bereinigung` (`open/`) — Zähler und Umleitungs-Weg in §6 und §8.
- **Risiken aus §6:** je ein Ausgang, mit Beleg in §6. *Entfallen:* Antragsweg nicht generisch · K3 gilt
  nur für Umbenennungen · Wert still abgebildet. *Weiter offen:* Informationsverlust ist nicht sichtbar
  (Adresse: `slice-transformationen-betriebsdoku` §2).
- **Frage an den Nutzer — Neubewertung der Kopf-Liste `tools/harness/blocked/go` (die
  Entscheidung liegt beim Nutzer; der Planner legt sie schriftlich vor und entscheidet sie nicht).**
  *Auslöser*, wörtlich aus dem Auflösungs-Trigger von `MR-003`: „sobald der Register-Eintrag eine weitere
  Beleg-Datei mit einem Host-Interpreter-Aufruf ohne Repo-Pfad im Befehlsstring und Wirkung auf eine
  Repo-Datei trägt, oder mit der Closure des nächsten Slice, dessen Läufe unter diesem Guard liefen
  (dann ist die Wirkung des Guards gemessen).“ Das zweite Kriterium ist mit dieser Closure eingetreten
  (dieser Slice folgt unmittelbar auf `slice-harness-guard-inplace-textwerkzeug`: der Start-Move
  `ce7acc8c` liegt nach dessen `→ done`-Move `e5a11979`, gelesen mit `git log`); das erste Kriterium ist
  nicht eingetreten (keine der drei Stellen der Beleg-Datei ist ein Host-Interpreter-Aufruf ohne
  Repo-Pfad mit Wirkung auf eine Repo-Datei). *Die Frage:* Soll das Fragment `tools/harness/blocked/go`
  angelegt werden — die Wortliste `go gofmt python python3 node dotnet java gradle uv`, die der Guard
  (Fragment-Ladung liegt im Guard, das Verzeichnis existiert nicht) am Kopf eines Kommando-Segments
  unbedingt blockt —, und mit welchem Umfang der Liste? *Was die Liste ändert* (hergeleitet aus dem
  Guard und dem Plan von `slice-harness-guard-inplace-textwerkzeug` §6, nicht erprobt): sie schlösse
  Host-`python`/`python3`-Aufrufe unabhängig vom Befehlsstring, auch die zwei Ränder, die der Guard heute
  nicht liest (Pfad ohne Repo-Namen, `cd <Repo> && python3 …`); `docker run … go` und `make …` blieben
  frei, weil `docker` und `make` der Kopf sind; ein harmloser Aufruf (`python3 --version`) würde
  ebenfalls blockiert; `perl` steht nicht in der Liste; die drei Vorfälle des Auslösers nennen keinen Host-`go`-Aufruf
  (Plan des Guard-Slice, übernommen; für die übrigen Einträge der Liste ist der Bestand nicht gesucht).
  *Gemessene Wirkung des Guards in den Läufen dieses Slice* (Ursprung je Zeile): ein Guard-Treffer —
  der `sed -i` des Verifiers auf einer Scratch-Hilfsdatei, geblockt, ohne Wirkung auf eine Repo-Datei
  (Verifikation V-5, **übernommen**); zwei Wege am Guard vorbei, beide ohne Wirkung auf eine Repo-Datei
  außer der beabsichtigten — das Anhängen per `cat >>` im Lauf des Implementers (Review F-7,
  **übernommen**; der Guard liest Umleitungen nicht) und das `cd <Repo> && python3 --version` des
  Planners dieser Closure (**gemessen** in dieser Sitzung: der Aufruf lief; die Hook-Eingabe mit
  diesem Befehlsstring am Guard ergibt Exit 0 ohne Ausgabe, mit einem Repo-Pfad als Argument
  die Meldung der Klasse `interp`); kein Host-Interpreter-Aufruf ohne Repo-Pfad mit Wirkung auf eine
  Repo-Datei nach dem Kenntnisstand von Review und Verifikation (übernommen). *Was nicht gemessen ist:*
  die Zahl der Blockierungen in den Sitzungen von Implementer und Reviewer — der Guard führt kein
  Protokoll, die Zahlen kommen nur aus den Berichten der Rollen; ein Slice ist eine Stichprobe, keine
  Wirkungsmessung. *Wege* (ohne Empfehlung des Planners; die Empfehlung des Guard-Plans lautete
  „nicht in diesem Slice, zuerst die Wirkung abwarten“, übernommen): (1) das Fragment mit der vollen
  Liste anlegen; (2) unverändert lassen — der Rest bleibt Sache des Reviews, die Neubewertung
  folgt dem ersten Kriterium; (3) eine kleinere Liste (etwa nur `python python3`, die einzigen
  Interpreter mit Belegen). Bei „ja“ legt der Planner einen eigenen Plan an; an `MR-003` ändert eine
  Entscheidung nichts (Auflösungs-Trigger dort: ein Auftreten trotz Guard ist eine Beleg-Datei im
  Register-Eintrag). Der Stand steht auch im `state.md` des Register-Eintrags
  `BEO-PGC/inplace-textwerkzeug-am-repo-trotz-nutzerregel`.
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-transformationen](../welle-transformationen.md) (offen) — die Closure der Welle prüft sie mit;
  die Slice-Closure trägt sie zusätzlich jetzt: *Anker:* die Steering-Loop-Einträge tragen kein
  `liegt in`-Feld eines verkörperten Ziels (nichts wird mit diesem Slice verkörpert); die Zielorte der
  Adressen (`welle-transformationen` §5 Punkte (c) und (d), `slice-transformationen-betriebsdoku` §2)
  existieren und tragen den Text (gelesen). *Folge-Slice:* keiner angelegt; `slice-transformationen-e2e-wirkung`,
  `slice-transformationen-betriebsdoku` und `slice-code-kommentare-bereinigung` existieren als Dateien in
  `open/`. *Register:* jede genannte Kennung `BEO-PGC/<slug>` existiert als Verzeichnis mit nicht leerem
  `evidence/` (geprüft mit `ls docs/plan/planning/observations/BEO-PGC/<slug>/evidence`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); die Domäne ist keine eigene Sub-Area — kein Anlass zur
Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 6×),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf §3 —
die Menge der Regeltypen ist die bewegte Eigenschaft),
`BEO-PGC/slice-chronik-in-code-kommentar` (verkörpert, 9×),
`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung` (verkörpert, 6×,
DoD Punkt 2 nennt seinen Beleg-Anker); übrige Einträge gesichtet, kein Bezug.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
