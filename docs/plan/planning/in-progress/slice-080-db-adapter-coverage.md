# Slice slice-080: DB-Adapter-Coverage — eigene, subjekt-qualifizierte Messung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** ohne Welle — die Closure-Bedingung ist die eigene DoD (eine **zweite,
getrennte** Messung mit eigener Schwelle), kein repo-weites *Mehr*. Ausdrücklich
**nicht** Teil von `welle-20`: deren §6 schließt sie als eigenen Vorgang aus.

**Bezug:** [`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 3 (die Messung, ihr Träger und ihre Schwelle, unverändert nach dem
bootstrap-aware-Muster aus [`ADR-0054`](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
§(a)); [`ADR-0030`](../../adr/0030-testpyramide.md) (die Test-Tiers);
`harness/README.md` §Werkzeuge (die beiden Träger-Läufe).

**Berührte Spec-Stellen:** — (kein Spec-Stratum: eine Messung an bestehenden
Tests, ohne Vertragsänderung).

**Verantwortlich:** pt9912.

**Autor:** pt9912 (Planner). **Datum:** 2026-09-15.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die drei Pakete, die aus dem Unit-Messgegenstand ausgeschlossen sind
(`postgresstorage` **ohne** das Unterpaket `mapper`, `postgresack`,
`replication/receive`), bekommen eine **eigene Coverage-Zahl** — gemessen dort,
wo ihre Tests ohnehin gegen einen echten PostgreSQL laufen (`make test-store`,
`make test-replication`), mit eigenem `-coverprofile` und eigener, nach
demselben bootstrap-aware-Muster kalibrierter Schwelle. Die Zahl trägt **immer
ihr Subjekt** („DB-Adapter-Coverage") und heißt nie „die Coverage" des Repos
([`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 4).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein Container in `make gates`** — der von `ADR-0071` verworfene Preis
  (Option C): das Gate läuft bei jedem Commit.
- **Die Unit-Zahl** — sie bleibt die Gate-getragene; diese Messung ist eine
  **zweite, getrennte**, und ihre Schwelle gilt nur für ihr Subjekt.
- **Die Tests selbst** — der Slice ändert keinen Test; er **misst** die
  bestehenden. Fehlt Abdeckung, ist das ein Befund dieses Slice, kein Auftrag.
- **Eine Zusammenführung beider Zahlen** — zwei Zahlen, zwei Gegenstände, ein
  Name je Gegenstand; eine gemeinsame Zahl wäre die Doppelquelle, die
  `ADR-0071` Punkt 4 ausschließt.
- **Produktionscode unter `internal/**`** — Schicht-Abgrenzung: Messung, Runner,
  Workflow.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

**Drei Liefer-Punkte** — die Kriterien darunter sind ihre Prüf-Form, kein
vierter Punkt:

**Liefer-Punkt 1 — die Messung existiert.**

- [x] Die beiden Träger-Läufe erzeugen **je ein** `-coverprofile` für ihren
      Testbestand; die Profile werden **gemergt** und als **eine** Zahl
      ausgegeben — das Merge-Verfahren ist benannt.
- [x] Die **Herkunft der Zahl** steht dabei (welcher Lauf, welches Profil,
      welche Deduplizierung): die Zählbasis-Regel, die `slice-079` für die
      Unit-Zahl verkörpert hat, gilt hier analog.

**Liefer-Punkt 2 — die Schwelle ist kalibriert.**

- [x] Der erste reale Wert wird **gemessen** und die Stufe nach dem
      bootstrap-aware-Muster gesetzt (Ist-Stand, abgerundet auf die volle
      5-%-Stufe; `ADR-0054` §(a)) — keine erfundene Zahl, keine Senkung.
- [x] Ihr **Träger** ist der nicht-blockierende Workflow
      (`.github/workflows/e2e.yml`), **nicht** `make gates`; die Zahl trägt ihr
      Subjekt im Namen und heißt nie „die Coverage".

**Liefer-Punkt 3 — der Beleg.**

- [x] Ein **realer** Lauf zeigt die Zahl, Exit direkt gelesen und ungepiped
      (`AGENTS.md` §3.9).
- [x] `make gates` bleibt **unverändert** grün — kein Container dort, keine
      zweite Schwelle im Bündel.

- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      `docs/reviews/review-slice-080.md` · Fixrunde
      `docs/reviews/review-slice-080-fixrunde.md`.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. **Entfällt** — das Repo ist durchgehend Greenfield, die Datei existiert nicht.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im
      Repo **mit** Wellen-Betrieb: die Prüfung trägt die **nächste
      Welle-Closure** (`welle-20` läuft und sammelt auch Slices ohne
      Wellen-Zugehörigkeit ein). Der genannte Folge-Slice `slice-082` existiert
      als Datei in `open/`.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/run-store-tests.sh` | update | `-coverprofile` für den Store-Testbestand, Ausgabe der Teilzahl |
| `tools/harness/run-replication-tests.sh` | update | dito für den Replication-Testbestand; Merge beider Profile + Schwellen-Prüfung |
| `tools/harness/db-coverage.sh` | **neu** | der Merge der zwei Profile zur **einen** Zahl, ihre Zählbasis (Dedup über die Block-Position), die Träger-Subjektliste (`--coverpkg`) und `DB_COVERAGE_THRESHOLD` als **ein** beweglicher Ort |
| `harness/sensors/db-adapter-coverage.md` | **neu** | die Bindung dieser **zweiten** Zahl: Subjekt, Schwelle, Träger, Zählbasis, Grenzen |
| `harness/sensors/coverage-gate.md` §Grenze | update | die Unit-Sensor-Doku verweist auf die eigene Messung und hält fest, dass `mapper` **nicht** in beiden Gegenständen liegt |
| `harness/README.md` §Werkzeuge | update | die beiden Läufe nennen die Zahl, die sie erzeugen, samt Sensor-Verweis |
| `.github/workflows/e2e.yml` | update | der Träger: die Messung läuft dort, nicht in `make gates` |

**Nicht in dieser Liste:** `Makefile` (kein neues Gate-Target), `internal/**`
(kein Produktionscode), die Tests selbst (sie werden gemessen, nicht geändert).

**Der genaue Zuschnitt der Runner-Änderungen entsteht im ersten Implementer-Lauf**
— die Liste nennt die Träger, nicht jede Zeile.

**Nachzug aus dem ersten Implementer-Lauf:**

- **Der Messlauf ist ein eigener `go test`-Aufruf über den Gegenstand**, nicht
  der Tier-weite `./...`-Lauf: so trägt die Messung einen eigenen Exit und ist
  von dem vorbestehend roten Paket außerhalb des Gegenstands (s. u.) unabhängig.
  `-coverpkg` instrumentiert dabei nur die in den Testbinaries **verlinkten**
  Gegenstands-Pakete — der Store-Lauf trägt 610 Statements für
  `postgresstorage` und **keine Zeile** für `postgresack`/`receive`; im
  Replication-Lauf (zwei Testbinaries) erscheint jede Position **zweimal**.
  Der Store-Lauf zieht `postgresstorage` aus dem `others`-Sammelaufruf heraus und
  führt es als Messlauf **zuletzt** — seine Tests räumen das `cdc`-Schema ab und
  dürfen dem vorgezogenen `bootstrap`-Aufruf nicht das ausgerollte Schema
  entziehen.
- **Vorbestehender roter Tier-Lauf, nicht durch diesen Slice verursacht:**
  `make test-replication`s Tier-weiter `go test ./...` ist rot
  (`internal/bootstrap` · `TestWALRetentionThresholdEndToEnd`): dessen Fixture
  baut das `cdc`-Schema per `DROP SCHEMA cdc CASCADE` + `ApplySchema` neu auf
  und trägt die seit `ADR-0050` von `bootstrap.Run` gelesenen Tabellen
  (`cdc.administration_request`, `cdc.process_heartbeat`) nicht. Real belegt:
  derselbe Fehlschlag auf dem **unveränderten** Runner-Skript.
- **Der Träger-Schritt ist zweigeteilt** (Fixrunde F-3):
  `tools/harness/run-replication-tests.sh` trägt zwei per Argument wählbare
  Phasen — `measure` (Profil, Merge, Schwellen-Prüfung; ihr Exit ist das Verdikt
  der Messung) und `tier` (`go test ./...`). Der Workflow führt sie als **zwei
  Schritte**; im gemeinsamen Schritt verschluckte der rote Tier-Exit das Verdikt
  der Messung. `make test-replication` ohne Argument fährt beide Phasen für
  einen lokalen Einzelaufruf. Der Tier-Schritt bleibt rot, bis das Fixture
  nachgezogen ist.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): [`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
Punkt 3 ist entschieden, die beiden Träger-Läufe existieren
(`make test-store`, `make test-replication`), `.github/workflows/e2e.yml` läuft,
`Verantwortlich:` gesetzt, WIP-Limit frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): zeigt sich, dass die
  zwei Läufe **keine gemeinsame** Zahl tragen können (verschiedene Testbestände,
  verschiedene Nenner, kein sinnvolles Merge), ist der Schnitt zu grob — dann
  gehört er zurück zur Zerlegung: **eine Zahl je Lauf** statt einer gemeinsamen.
- `in-progress` → `open` (blockiert — Carveout?): erweist sich, dass die Messung
  den nicht-blockierenden Workflow über seine Laufzeitgrenze treibt (der
  Container-Aufbau steht schon, die Profile nicht), braucht sie eine Entscheidung
  (Carveout oder eigener Träger) — **nicht** stilles Weglassen.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** ein **realer** Lauf zeigt die kalibrierte Zahl **und**
`make gates` unverändert grün **und** die Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- **Die zwei Läufe tragen verschiedene Testbestände.** Ein gemeinsamer Nenner
  über beide kann eine Zahl erzeugen, die keinen der beiden Gegenstände
  beschreibt. Antwort: eine Zahl **je Lauf** oder ein ausdrücklich benannter
  gemeinsamer Nenner samt Verfahren. — **Ausgang:** *entfallen — gestrichen mit
  Begründung*: der gemeinsame Nenner **beschreibt** beide Gegenstände, weil die
  zwei Läufe einander **nicht** überschneiden — die Dateimengen der Profile
  schneiden sich leer, und `788 = 610 + 155 + 23` ist exakt die Größe, die
  [`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  für den Gegenstand gemessen hat. Das Verfahren steht im Skript und in der
  Sensor-Doku §Zählbasis.
- **Die erste Kalibrierung könnte sehr niedrig liegen.** Das ist zulässig
  (bootstrap-aware: die Stufe folgt der Messung) — aber der Wert gehört
  **benannt**, samt dem, was ihn drückt. — **Ausgang:** *entfallen — gestrichen
  mit Begründung*: der erste Wert liegt bei **75,25 %** und damit nicht niedrig;
  die Stufe (75) und das, was sie drückt (`tableactivation.go` 70,8 %,
  `receive.go` 69,5 %, `administrationrequest.go` 66,7 %), sind benannt.
  **Als benannte Grenze festgehalten** (nicht als Ausgang): die Stufe liegt nur
  **2 Statements** über ihrer Schwelle (591 wären 75,0 %) — sie folgt der
  Messung, und ein Testverlust kippt sie auf rot. Das ist die Mechanik, kein
  Defekt.
- **`postgresstorage/mapper` liegt in beiden Messungen** — es bleibt im
  Unit-Gegenstand. Die zwei Zahlen **überlappen** in diesem Paket; wer sie
  addiert oder vergleicht, muss das wissen. — **Ausgang:** *entfallen —
  gestrichen mit Begründung*: die **Prämisse des Risikos trifft nicht zu** —
  `postgresstorage/mapper` ist aus dem DB-Subjekt **ausgeschlossen** (das Subjekt
  führt `postgresstorage` **ohne** `mapper`), die zwei Zahlen überlappen also
  nicht. Der Implementer hat das an den Profilen nachgewiesen, der Reviewer hat
  es bestätigt. Ein Risiko, dessen Voraussetzung widerlegt ist, kann nicht
  eintreten.

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** die **Zählbasis** — der teuerste Teil von `slice-079`
  trug hier sofort. Die Regel „gedeckt = mindestens ein Vorkommen `count > 0`"
  ist bei diesem Slice **tragend**, nicht kosmetisch: ohne sie fiele
  `replication/receive` von 112/155 auf **0/155** und der Wert auf 10,11 % —
  unter jede Schwelle. Und der **Merge** ist keine Notlösung: die zwei Läufe
  partitionieren die Subjekttests, ihre Profile schneiden sich leer, und der
  gemeinsame Nenner (788) stimmt exakt mit dem überein, was
  [`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  für den Gegenstand gemessen hat.
- **Was ging anders als geplant:**
  1. **Ein Risiko des Plans war falsch angesetzt.** §6 Risiko 3 nahm eine
     Überlappung bei `mapper` an; `mapper` ist aus dem DB-Subjekt
     **ausgeschlossen**. Der Implementer hat es an den Profilen widerlegt und
     **gemeldet statt erfüllt** — der Ausgang ist damit *entfallen*, nicht
     „nicht eingetreten".
  2. **Ein vorbestehender roter Test wurde sichtbar** — nicht erzeugt, sondern
     **aufgedeckt**: `make test-replication`s tier-weites `go test ./...`
     scheitert seit `ADR-0050` an einem Fixture, das die seither gelesenen
     Tabellen nicht mitbringt. Das Target ist **weder Gate noch CI**, deshalb
     war es unsichtbar. Belegt am **unveränderten** Runner (Reviewer) und mit
     der Mechanik am Artefakt (`walretention_endtoend_test.go`;
     `administration_endtoend_test.go:37` dokumentiert die Lücke selbst).
     Eingetragen als `BEO-PGC/roter-test-ohne-leser`; die **Adresse** ist
     `slice-082`.
  3. **Die neue Verdrahtung machte ihn in CI sichtbar.** Seit diesem Slice läuft
     `e2e` rot (**nicht-blockierend**): rot ist der **Tier-Schritt**, grün der
     **Mess-Schritt** — die Zahl ist in CI also beobachtbar, der Lauf als Ganzes
     nicht. Das gehört zu 2. und wird von `slice-082` zurückgenommen, **ohne**
     zu maskieren. **Zum Post-Push-Beleg (`AGENTS.md` §3.10):** der reale Lauf
     liegt vor (`gh run view`, beide Matrix-Legs PostgreSQL 17/18), und die
     **neuen** Schritte sind darin **grün** — die Workflow-Änderung ist damit in
     ihrem eigenen Umfang bestätigt; offen ist allein der fremde Tier-Schritt.
  4. **Eine Mechanismus-Erklärung war falsch** (Review F-1): „`-coverpkg`
     instrumentiert in jeder Testbinary den ganzen Gegenstand" — real erscheint
     nur das **verlinkte** Paket. Dieselbe Klasse wie in `slice-079`.
- **Steering-Loop-Einträge:**
  - *verkörperte Regel, angewandt und präzisiert*: die Zählbasis (`slice-079`) —
    hier mit dem **richtigen** Grund für die Deduplizierung (Duplikat innerhalb
    **eines** Profils, nicht über den ganzen Gegenstand).
  - *zwei neue Register-Einträge*: `BEO-PGC/roter-test-ohne-leser` (1×) und
    `BEO-PGC/mechanismus-erklaerung-ohne-werkzeugbeleg` (2× über zwei Vorgänge).
  - *kein neuer Sensor* — beide Klassen haben ihren Wächter im Review.
- **Beobachtungs-Register (`../observations/`):** zwei **neue** Verzeichnisse
  (s. o.), je mit `evidence/`. `BEO-PGC/regel-weiter-als-ihr-sensor` bleibt bei
  **2×**: dieser Slice fügt ihm keinen Beleg hinzu — die DB-Zahl ist eine
  **dritte, eigene** Messung, keine Änderung der beiden dort geführten Zusagen.
- **Folge-Slices:** `slice-082` (Replication-Fixture nachziehen) — die Adresse
  für die Punkte 2 und 3; er liegt in `open/`.
- **Risiken aus §6:** alle drei *entfallen, gestrichen mit Begründung*; die
  **2-Statement-Marge** ist als benannte Grenze in §6 festgehalten, nicht als
  Ausgang.
- **Drei Paarungen:** hier **nicht** geprüft — die §2-Zeile verweist sie an die
  nächste Welle-Closure (`welle-20` läuft und sammelt auch Slices ohne
  Wellen-Zugehörigkeit ein).
- **Zwei benannte Grenzen:** (a) der Kopf-Satz des `e2e`-Workflows („Kein
  Workflow-Schritt enthält Inline-Shell-Logik, die eines dieser Ziele umgeht")
  stand im Widerspruch zu den zwei neuen Zeilen und ist **nachgezogen**;
  (b) **zwei Wege zum selben Skript** (`make test-replication` und die zwei
  CI-Aufrufe) sind eine Divergenz-Fläche — eine spätere Rezeptur-Änderung
  erreicht die CI-Schritte nicht. Benannt, nicht behoben.
- **Eine Buchhaltungs-Lücke, benannt:**
  [`ADR-0071`](../../adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  §Fitness Function führt „neues Target (geplant)"; angelegt ist **kein** Target
  — die Messung hängt an den bestehenden Träger-Läufen plus einem Hilfsskript.
  Operativ konform, wörtlich nicht.
- **Archivierung:** nicht ausgeführt — das Repo führt kein Archivierungswerkzeug.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
repo-weite Default-Sub-Area `*`/`PGC` — sie deckt `tools/harness/**`,
`harness/sensors/**`, `harness/README.md` und `.github/workflows/` in **einem**
Kürzel. Eine feinere Sub-Area ist nicht zu bilden: „Test-Infrastruktur" ist keine
deklarierte Sub-Area, und das Register führt den Runner unter demselben Kürzel.

**Vorgelagert — offene Beobachtungen sichten:** Register
(`../observations/BEO-PGC/`) durchgegangen, gemergter Stand; zwei Treffer:

- `BEO-PGC/regel-weiter-als-ihr-sensor` (2×, offen): **kein** Treffer — die
  beiden Belege betreffen die host-lokale Pfad-Regel und die **Unit**-Fitness-
  Function. Dieser Slice fügt eine **dritte, eigene** Messung hinzu und ändert
  keine der beiden Zusagen.
- `BEO-PGC/endstufe-unter-eigenem-messgegenstand-unerreichbar` (1×,
  `verkörpert` → `ADR-0071`): der **Vorläufer** — dieser Slice setzt seinen
  Punkt 3 um. **Kein neuer Beleg**: derselbe Gegenstand, derselbe Träger.

**Kein** Eintrag erreicht mit diesem Slice die 3×-Schwelle; es entsteht kein
neues Verzeichnis.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (nur `*`/`PGC`)
— kein Modus-Begründungsblock. Die vier Pflichtkriterien tragen dennoch:
**Konventionen-Dichte** hoch (die beiden Träger-Läufe existieren,
`harness/sensors/` ist die Bindungsform), **Phase-Reife** hoch (das Gate läuft
seit `slice-049`, die Träger-Läufe seit `slice-050`), **Evidenz-/Diskrepanz-Risiko**
niedrig — die Zahl wird **gemessen**, nicht gegen einen Bestand inventarisiert —,
**Reconciliation-Aufwand** null.
