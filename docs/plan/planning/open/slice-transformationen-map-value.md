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

**Verantwortlich:** — (noch nicht priorisiert).

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

- [ ] `map_value` wirkt auf beide Images: Zuordnung, nicht zugeordneter Wert
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
      geänderte Erwartung.
- [ ] Ohne Änderung an Antragsweg und Wirkort: `git diff --stat` gegen den
      Parent-Stand nennt weder `internal/application/usecase/`,
      `internal/bootstrap/`, `tools/schema/`, `internal/adapters/**` (außer
      Testdateien, die die Regeltypen aufzählen) noch die Spec; der Use Case
      akzeptiert einen `map_value`-Antrag über den Parser der Domäne. *Zu
      belegen durch:* der Diff-Stat im Bericht und ein Use-Case-Test, der einen
      `map_value`-Antrag ohne Änderung des Use-Case-Codes annimmt.
- [ ] Die Fitness Function bleibt vollständig: der Eigenschaftstest (Regeltyp ×
      ausgeschlossene Spalte) erfasst `map_value` über die Domänen-Menge ohne
      manuelle Ergänzung — das Image trägt weder Quellschlüssel noch Zielname
      noch Quellwert noch abgebildeten Wert; der Paritätstest des
      Backfill-Pfads erfasst `map_value` ohne Strukturänderung. *Zu belegen
      durch:* `make test` und die Mutation, die `map_value` aus der
      Domänen-Menge entfernt (die Tests färben sich rot); `make a-check` grün.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: entfällt — keine Betreiber-Oberfläche; der Regeltyp steht im
      Handbuch-Abschnitt von `slice-transformationen-betriebsdoku`.
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
| `internal/domain/model/transformation.go` (aus `kern-rename`; + Test) | update | zweiter Regeltyp: Konstruktor-Invarianten, Auswertung, Parser-Zweig. |
| `internal/adapters/driving/replication/mapper/` (Eigenschaftstest) | update | Regeltyp-Menge aus der Domäne — `map_value` wird ohne manuelle Ergänzung erfasst. |
| Paritätstest des Backfill-Pfads (Ort aus `backfill-pfad`) | prüfen | tabellengetrieben: erfasst `map_value` ohne Strukturänderung. |

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

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Aufzählungen der Regeltypen im Code | `grep -rn 'rename_column' internal tools --include=*.go --include=*.sql --include=*.yaml` | *(Implementer trägt ein)* | jede Stelle, die die Regeltypen einzeln aufzählt, ist eine zweite Liste: sie zieht auf die Domänen-Menge oder trägt den neuen Typ |
| Aufzählungen der Regeltypen in Docs | `grep -rn 'rename_column' docs spec harness` | *(Implementer trägt ein)* | Spec trägt beide Typen bereits; Handbuch-Träger an `betriebsdoku` melden |
| Tests, die eine feste Typ-Liste führen | `grep -rn 'Regeltyp\|kind' internal --include=*_test.go` | *(Implementer trägt ein)* | auf die Domänen-Menge umstellen oder den Typ ergänzen |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-transformationen-backfill-pfad`
in `done/` liegt (die Bindung der Erzeugungspfade steht, bevor der zweite Typ
hinzukommt) und kein anderer Slice in `in-progress/` liegt (WIP-Limit 1).

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
