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
      §3).
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
      ausgeschlossene Spalte) erfasst `map_value` über die Domänen-Menge ohne
      manuelle Ergänzung — das Image trägt weder Quellschlüssel noch Zielname
      noch Quellwert noch abgebildeten Wert; der Paritätstest des
      Backfill-Pfads erfasst `map_value` ohne Strukturänderung. *Zu belegen
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
| `internal/domain/model/transformation.go` (aus `kern-rename`; + Test) | update | zweiter Regeltyp: Konstruktor-Invarianten, Auswertung, Parser-Zweig. Gelieferte Form: die Konstante `TransformationMapValue`, `NewMapValue`, das Feld `values` (die Zuordnung in kanonischer Kodierung: Paare in aufsteigender Ordnung des Schlüssels, je Feld ein Vier-Byte-Längenpräfix), `Values()`, die Auswertung in `applyTransformations` und die Anwendbarkeit ohne Zielname-Prüfung für `map_value` in `CheckApplicable`. |
| `internal/domain/model/transformationspec.go` (Parser-Zweig) | update | (Nachzug, aus der Zeile oben herausgelöst) `ParseTransformationSpec` liest `values` strikt (nicht leeres Objekt, Zeichenketten als Werte), `allowedRuleKeys` nennt `kind`/`column`/`values`, `Build` baut die Regel. |
| `internal/domain/errors/errors.go` (Doc-Kommentare) | update | (Nachzug) die Kommentare von `ErrInvalidTransformation` und `ErrInvalidRuleSpec` nennen `values`; der Block von `ErrInvalidTransformation` trug zwei Kennungen und trägt eine (`AGENTS.md` §3.7). |
| `internal/domain/model/transformation_mapvalue_test.go` (neu), `transformation_test.go`, `transformationspec_test.go` | update | (Nachzug) Tests der Domäne für `map_value`. Zwei bestehende Erwartungen ändern sich zwangsläufig: `TestTransformationKindsIsAClosedSet` (die Menge trägt zwei Typen) und `TestTransformationSpecBuild` (der Platzhalter für einen unbekannten Regeltyp war `map_value`, jetzt `nope`). Die Kodierung der Zuordnung trägt Felder bis 70 000 Byte (`TestMapValueEncodingCarriesFieldsUpTo70000Bytes`) und oberhalb von 16 777 216 Byte (`TestMapValueEncodingCarriesFieldsBeyondSixteenMiB`, die Stufe des vierten Längen-Bytes; gemessen unter `-race` rund 0,1 s, Speicher-Spitze zwischen 200 und 300 MiB: bei 200 MiB Cgroup-Grenze bricht der Testlauf ab, bei 300 MiB läuft er durch); das Kürzen des Längen-Präfixes auf drei Byte färbt den zweiten Test rot (gesehen). |
| `internal/adapters/driving/replication/mapper/` (Eigenschaftstest) | update | Regeltyp-Menge aus der Domäne — `map_value` wird ohne manuelle Ergänzung erfasst. Gelieferte Form: `transformation_test.go` (Fixture-Schalter `ruleFor` trägt den Fall `map_value`, die Erwartung je Regeltyp steht als Schlüssel-Wert-Paar am Fall) und die neue Datei `transformation_mapvalue_test.go` (beide Images, Abwesenheit, Metadaten, Anwendbarkeit, Determinismus, Live-Reload, `-race`). |
| Paritätstest des Backfill-Pfads (Ort aus `backfill-pfad`) | prüfen → update | tabellengetrieben: die Schleife über `model.TransformationKinds()` erfasst `map_value` ohne Strukturänderung; der Fixture-Schalter `parityRule` in `internal/bootstrap/backfill_image_parity_test.go` bekommt den Fall `map_value` (Übergabe aus `backfill-pfad`, siehe unten), die Prüfung der sichtbaren Wirkung je Regeltyp steht am Fall. Die Fixrunde nach dem Review bindet den Ausschluss an der Eingabeseite: ist die Spalte der Regel ausgeschlossen, trägt das Bild weder ihren Schlüssel noch ihren Quellwert noch die Wirkung der Regel; die Mutation, die in `BuildRowImage` den Ausschluss für eine Spalte mit `map_value`-Regel aufhebt, färbt den Paritätstest rot (gesehen, drei Fälle `map_value/Regel an …/ausgeschlossen`). |
| `internal/adapters/driving/replication/mapper/mapper.go` (Doc-Kommentar von `ErrTransformationNotApplicable`) | update | (Nachzug, Fixrunde nach dem Review) nur Kommentar, kein Verhaltensdiff: der Kommentar nannte als Nichtanwendbarkeit den Zielnamen, den nicht jede Regel trägt; er bindet den Fall jetzt an eine Regel mit Zielname (ohne einen Regeltyp beim Namen zu nennen, das Suchlauf-Feld zählt Produktivcode außerhalb der Domäne ohne Regeltypnamen) und trägt eine Kennung statt drei (`AGENTS.md` §3.7). Die einzige Produktivdatei außerhalb der Domäne im Diff; DoD Punkt 2 nennt sie als Ausnahme. |
| `internal/application/usecase/backfill/transformation_test.go` | update | (Nachzug, Übergabe aus `backfill-pfad`) Fixture-Schalter `ruleFor` bekommt den Fall `map_value`; die Erwartungen der Bild-Tests folgen dem Regeltyp; neue Fälle: Zustandswechsel des Regelstands mit `map_value` und Nichtanwendbarkeit im Run. |
| `internal/adapters/driven/postgresstorage/administrationrequest_test.go` | update | (Nachzug, über den Plan hinaus, nur Test) `TestTableActivationTransformationRulesReadMapValueThroughJsonb`: die Regelform von `map_value` läuft durch die reale `jsonb`-Spalte (Schlüssel-Umordnung, Normalisierung) und wird zur selben Regel gefaltet wie die gebaute; Tier `make test-store`. Grund: der Slice ändert die Faltung des Regelstands, und der `jsonb`-Weg ist der einzige Ort, an dem die Ordnung der Zuordnung nicht vom Aufrufer kommt. |
| `internal/application/usecase/settransformation/service_test.go`, `internal/bootstrap/administration_internal_test.go` | update | (Nachzug, DoD Punkt 2) der Use-Case-Test eines `map_value`-Antrags am unveränderten Use-Case-Code und sein Pendant über die Verdrahtung `applyAdministrationRequest`. |

**Auflösung der Übergabe aus `backfill-pfad` (Entscheidung dieses Slice, vor dem Code):**
`Transformation` bleibt über `==` vergleichbar; die Zuordnung von `map_value` steht als kanonische
Zeichenkette im Feld `values` (erster der beiden benannten Wege). `sameSet[T comparable]` in
`internal/application/usecase/backfill/service.go` bleibt unverändert, der Diff nennt keine
Produktivdatei unter `internal/application/usecase/`. Die Kosten: die Suche einer Zuordnung ist linear in
der Zahl der Paare und dekodiert nicht (`lookupMappedValue`). **Auslegung von DoD Punkt 2:** die
Ausnahme „Testdateien, die die Regeltypen aufzählen“ gilt für jedes Verzeichnis der Aufzählung, nicht nur
für `internal/adapters/**` — `internal/application/usecase/backfill/transformation_test.go`,
`internal/bootstrap/backfill_image_parity_test.go` (Fixture-Schalter) sowie
`internal/application/usecase/settransformation/service_test.go` und
`internal/bootstrap/administration_internal_test.go` (der geforderte Use-Case-Test) sind Testdateien mit
Diff; kein Produktivcode dort. Der Store-Test `administrationrequest_test.go` (`internal/adapters/**`)
zählt keine Regeltypen auf und steht als Test über den Plan hinaus in der Tabelle. Der Wortlaut von DoD
Punkt 1 und 2 trägt diese Auslegung seit der Fixrunde nach dem Review selbst (Finding F-2 in
`docs/reviews/review-slice-transformationen-map-value.md`): Punkt 1 nennt die zwei zwangsläufig geänderten
Erwartungen, Punkt 2 den Ausschluss `':!*_test.go'` und die Kommentar-Korrektur in `mapper.go`. **Beleg von
DoD Punkt 2:** der Diff-Stat im Bericht, gemessen mit
`git diff --stat e5a11979 -- internal/application/usecase internal/bootstrap tools/schema internal/adapters spec ':!*_test.go'`;
er nennt genau eine Datei, `internal/adapters/driving/replication/mapper/mapper.go` (nur Kommentar), keine
Zeile in `internal/application/usecase/`, `internal/bootstrap/`, `tools/schema/` oder `spec/`.

**Nicht realisiert (Grund):** der Doc-Kommentar von `BuildRowImage` (`internal/domain/model/rowimage.go`)
trägt keinen Satz zur Wertabbildung: sein Block trägt bereits mehrere Kennungen, ein Zusatz machte den
Block zum Kandidaten von `make kommentar-kennungen`, und die Bereinigung liegt bei
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
| Träger, die den Zielnamen beschreiben (Beschreibung statt Symbolname; Fixrunde nach dem Review) | 23–24 (Code außerhalb der Domäne) und 25–26 (Spec, Handbuch) | **Gefunden.** „Zielname“ in Nicht-Test-Code außerhalb `internal/domain`: Parent 8, Diff 9 (Parent: vier Zeilen im `mapper`-Paket, eine im Backfill-Use-Case, zwei im Antrags-Use-Case, eine in der Verdrahtung; Diff: die korrigierte Definition trägt das Wort in zwei Zeilen statt einer). Davon ist genau eine Stelle die Definition der Nichtanwendbarkeit — der Doc-Kommentar von `ErrTransformationNotApplicable` in `mapper.go` —; sie führte den Zielnamen als Grund jeder Regel und ist korrigiert (Fixrunde, §3-Zeile oben). Die übrigen sieben nennen K3 oder die Kollision zweier Regeln mit gleichem Zielnamen und bleiben für `rename_column` wahr. `spec` und `docs/user`: Parent und Diff gleich (die Spec trägt beide Typen). **Nichtgefunden:** keine weitere Stelle, die den Zielnamen als Eigenschaft jeder Regel beschreibt. | `mapper.go`-Kommentar nachgezogen; die übrigen Zeilen unverändert |

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
  Use-Case-Test aus DoD Punkt 2. **Ausgang:** *(bei Closure)*
- **K3 gilt nur für Umbenennungen** (Zielnamen); `map_value` trägt keinen
  Zielnamen — eine Prüfung, die ihn erwartet, scheitert am zweiten Typ.
  *Erwartet, zu belegen durch:* der Test eines `map_value`-Antrags gegen die
  Konfliktfreiheit (K2 greift: eine Quellspalte trägt höchstens eine
  Spaltenregel). **Ausgang:** *(bei Closure)*
- **Ein Wert, der als Zeichenkette kein Schlüssel ist, wird still abgebildet**
  (Normalisierung, Groß-/Kleinschreibung). *Erwartet, zu belegen durch:* Test
  der exakten Zeichenketten-Gleichheit; die Festlegung steht in der Spec.
  **Ausgang:** *(bei Closure)*
- **Informationsverlust ist nicht sichtbar** (mehrere Quellwerte, ein
  Zielwert): kein Fehler, keine Warnung. Das ist eine benannte Konsequenz von
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md);
  die Betreiber-Aussage steht im Handbuch-Abschnitt von `betriebsdoku`.
  **Ausgang:** *(bei Closure)*

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
