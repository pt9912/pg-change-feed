# Slice slice-033: Typ-Auswertung und Fehlerklasse `schema` für inkompatible Typänderungen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-10`](welle-10.md) — der Nachweis, dass die
`ADR-0015`-Folgepflicht real eingelöst ist, ist `welle-10`s Closure-Trigger
(§3); dieser Slice ist der letzte der drei geplanten Slices und schließt
`LH-FA-SCH-004`s Negative-Fall real.

**Bezug:** [`LH-FA-SCH-004`](../../../../spec/lastenheft.md), `ADR-0015`
(nur gelesen — keine aktive ADR wird geändert, `ADR-0015` bleibt
`Accepted`).

**Berührte Spec-Stellen:** [`SPEC-008`](../../../../spec/pflichtenheft.md)
(Fehlerklasse `schema`), `LH-FA-SCH-004.a`.

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-12.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `Assembler.observeRelation`
(`internal/adapters/driving/replication/mapper/mapper.go`) meldet den
Fall `relationOther` (Spalte entfernt, Typ einer bestehenden Spalte
geändert, Spalte umbenannt — jede nicht sicher als reine
Obermengen-Erweiterung erkennbare Relation-Änderung) ab sofort als
**sichtbaren Fehler der Fehlerklasse `schema`**
([`SPEC-008`](../../../../spec/pflichtenheft.md)), statt ihn wie bisher
(`slice-032`) konservativ stillschweigend zu ignorieren. Das schließt
`LH-FA-SCH-004`s Negative-Fall real: „Given eine inkompatible
Typänderung, when sie eintritt, then ist sie erkennbar gemeldet; die
Daten werden nicht still fehlinterpretiert." `TestMVPSchemaChangeIncompatibleTypeChange`
(`slice-030`) läuft danach mit angepasster Erwartung an seinen zweiten
Testfall (PostgreSQL-seitig zugelassene Typänderung) grün: **statt**
stillschweigender Übernahme erwartet er jetzt einen sichtbaren
`schema`-Fehler. `TestMVPSchemaChangeAddColumn` (kompatible Erweiterung,
bereits von `slice-032` real geschlossen) bleibt unverändert im
Verhalten — dieser Slice rührt nur den `relationOther`-Zweig an.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Entfernte Spalten differenziert behandeln** (`LH-FA-SCH-003`) — Bestand
  bleibt bewusst stehen: eine entfernte Spalte fällt unter denselben
  `relationOther`-Zweig wie jede andere nicht sicher interpretierbare
  Änderung; eine feinere, `LH-FA-SCH-003`-spezifische Behandlung (z. B.
  definiertes Verhalten für historische Werte) ist ein eigener, hier
  ausgeschlossener Vorgang (bereits in `slice-030` §1 als Out-of-Scope
  benannt).
- **Wiederherstellung/Recovery nach einem `schema`-Fehler** (z. B.
  automatischer Neustart mit neuer Baseline-Version) — anderer Vorgang:
  dieser Slice meldet den Fehler sichtbar, definiert aber nicht, wie ein
  Operator danach fortsetzt (das ist bereits heute das etablierte
  Verhalten anderer `schema`-Klassen-Fehler in diesem Repo, z. B.
  `ErrTruncateUnsupported` — konsistentes Verhalten, keine neue
  Recovery-Fähigkeit).
- **Änderung der Klassifikationslogik selbst
  (`classifyRelationColumns`)** — Bestand bleibt bewusst stehen: die
  Drei-Fälle-Unterscheidung aus `slice-032` ist korrekt und bereits
  real verifiziert; dieser Slice ändert nur, **was im dritten Fall
  passiert**, nicht die Klassifikation selbst.

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

- [x] `observeRelation` meldet `relationOther` als sichtbaren Fehler der
      Fehlerklasse `schema` (Muster: `Fehlerklasse schema: …`, analog zu
      `ErrTruncateUnsupported`/`ErrChangeWithoutBegin` in derselben
      Datei) statt `nil` zurückzugeben. Unit-getestet (Assembler-Tests).
      Beleg: `internal/adapters/driving/replication/mapper/mapper.go`
      (`ErrIncompatibleSchemaChange`, `observeRelation`), Test
      `TestConsumeRelationOtherChangeReportsSchemaError` (`mapper_test.go`,
      grün über `make test`).
- [x] `LH-FA-SCH-004`s Negative-Fall real geschlossen:
      `TestMVPSchemaChangeIncompatibleTypeChange` (`slice-030`) läuft mit
      angepasster Erwartung grün — der PostgreSQL-seitig zugelassene
      Typänderungs-Fall erwartet jetzt einen sichtbaren `schema`-Fehler
      statt stillschweigender Übernahme. `TestMVPSchemaChangeAddColumn`
      bleibt unverändert grün (Regressionsschutz für die dynamische
      Re-Versionierung). `make gates` und `make test-integration` dreimal
      in Folge grün. Beleg: `test/integration/integration_test.go`
      (`awaitHeartbeatErrorClass`, Fall 2), drei aufeinanderfolgende
      grüne `make test-integration`-Läufe (Runner-Skript-Anpassung s.
      Plan-Nachzug).
- [x] `make gates` grün. Beleg: `make gates`-Lauf,
      `d-check: 275 Datei(en) geprüft, 0 Befund(e)` /
      `a-check: gesamt: 0 Befund(e)`.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Beleg: [`docs/reviews/review-slice-033.md`](../../../reviews/review-slice-033.md)
      (1 HIGH, 1 MEDIUM), Fixrunde behoben in Commit `ef16e38`, bestätigt
      in [`docs/reviews/review-slice-033-fixrunde.md`](../../../reviews/review-slice-033-fixrunde.md).
      Verifikation in
      [`docs/reviews/verify-slice-033.md`](../../../reviews/verify-slice-033.md)
      (DoD eigenständig nachgeprüft, `LH-FA-SCH-004` dreifach real gegen
      PostgreSQL bestätigt).
- [x] Doku-Update, falls ein öffentlicher Vertrag berührt wird —
      Implementer entscheidet und begründet im Plan-Nachzug. Kein
      Doku-Update nötig: `LH-FA-SCH-004.a`/`SPEC-008` beschreiben den
      Vertrag bereits generisch („sichtbarer Fehler statt stiller
      Fehlinterpretation"); dieser Slice schließt die Implementierungslücke,
      ändert aber keinen Vertragstext. `mapper.ErrIncompatibleSchemaChange`
      ist eine neue exportierte Kennung derselben, bereits dokumentierten
      Fehlerklasse `schema` — kein neuer Vertrag.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. Siehe §7.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Entfällt: Repo ist Greenfield (`harness/conventions.md` §Modus-Deklaration), `../reconciliation.md` existiert nicht.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
      Beleg: (1) `evidence/slice-033.md` in
      `BEO-PGC/schema-evolution-nicht-dynamisch/` ergänzt — Ausgang
      **verkörpert** (2×, direkt aufgelöst analog
      `BEO-PGC/spec008-replication-luecke`). (2) Neue Beobachtung
      `BEO-PGC/test-runner-stiller-ausschluss/` angelegt (1×, weiter
      offen) für den vom Reviewer gefundenen F-2-Risikoklasse. Siehe §7.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen). Siehe §6.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Entfällt hier: Repo mit Wellen-Betrieb — Prüfung läuft bei der `welle-10`-Closure.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driving/replication/mapper/mapper.go` | update | `observeRelation` meldet `relationOther` als sichtbaren `schema`-Fehler |
| `test/integration/integration_test.go` | update | `TestMVPSchemaChangeIncompatibleTypeChange`s zweiter Testfall erwartet jetzt einen sichtbaren Fehler |
| `internal/adapters/driving/replication/mapper/mapper_test.go` | update | Plan-Nachzug: `TestConsumeRelationOtherChangeStaysConservative` umbenannt/umformuliert auf `TestConsumeRelationOtherChangeReportsSchemaError` |
| `internal/bootstrap/wiring.go` | update | Plan-Nachzug: `classifyRunError` übersetzt `mapper.ErrIncompatibleSchemaChange` in `model.ErrorClassSchema` (Konsistenz-Anforderung §1) |
| `internal/bootstrap/heartbeat_internal_test.go` | update | Plan-Nachzug: Testfall für `classifyRunError(mapper.ErrIncompatibleSchemaChange)` ergänzt |
| `tools/harness/run-integration-tests.sh` | update | Plan-Nachzug: der sichtbare `schema`-Fehler beendet den Erfassungspfad des einzigen, geteilten Feed-Containers dauerhaft (`restart: "no"`) — `TestMVPSchemaChangeIncompatibleTypeChange` läuft deshalb als eigener, letzter `go test`-Aufruf, nach Lasttest-Beleg und Black-Box-CLI-Rundlauf, statt im ursprünglichen gemeinsamen `go test ./test/integration/...`-Aufruf |

### Plan-Nachzug (nach Implementierung)

- **Fehler-Sentinel:** `mapper.ErrIncompatibleSchemaChange` — Muster
  identisch zu `ErrTruncateUnsupported`/`ErrChangeWithoutBegin`
  (`errors.New("Fehlerklasse schema: …")`); `observeRelation` wrappt ihn
  mit dem qualifizierten Tabellennamen (`fmt.Errorf("%w: %s", …,
  relation.QualifiedName())`). `internal/bootstrap/wiring.go`s
  `classifyRunError` ordnet ihn `model.ErrorClassSchema` zu — dieselbe
  Übersetzungsstelle wie für `decode.ErrSchema`/`ErrTruncateUnsupported`
  (Konsistenz-Anforderung aus §1, mit eigenem Testfall in
  `heartbeat_internal_test.go` belegt).
- **Reale Konsequenz, vor dem ersten grünen Lauf rot gesehen:** Der
  gemeldete `schema`-Fehler beendet `receive.Stream.Run` (`Consume` →
  `process` → `Run` gibt den Fehler zurück) und darüber `bootstrap.Run` →
  `os.Exit(1)`. Der Compose-Feed-Container trägt `restart: "no"` (kein
  Neustart-Vertrag) und **eine einzige** Replication-Verbindung für
  **alle** aktivierten Tabellen (`CDC_TABLES` in `compose.yaml`) — ein
  ausgelöster `schema`-Fehler an `feed_mvp_schema` beendet damit den
  gesamten Feed-Container, nicht nur den Erfassungspfad dieser einen
  Tabelle. Der erste Lauf mit ungeändertem `tools/harness/run-integration-tests.sh`
  bestätigte das real: der anschließende Lasttest-Beleg
  (`cdc_capture_lag`, braucht einen laufenden Container) und der
  Black-Box-CLI-Rundlauf hätten denselben, bereits beendeten Container
  gebraucht. Behoben durch Umstellen der Aufruf-Reihenfolge — nicht durch
  eine neue Recovery-Fähigkeit des Produkts (§1-Ausschluss bleibt
  unberührt): `go test ./test/integration/...` lief bisher als ein
  einziger Aufruf vor Lasttest-Beleg und CLI-Rundlauf; er läuft jetzt als
  zwei `-run`-gefilterte Aufrufe — die sechs unveränderten Testfunktionen
  vorn (unveränderte Position), `TestMVPSchemaChangeIncompatibleTypeChange`
  separat und zuletzt, nach dem Black-Box-CLI-Rundlauf. Eine künftig
  ergänzte Testfunktion dieses Pakets muss in das vordere `-run`-Muster
  aufgenommen werden, sofern sie den Container nicht selbst beendet — im
  Runner-Skript kommentiert.
- **Bild-Rebuild-Falle (Steering-Loop-relevant, s. §7 bei Closure):**
  Der erste `make test-integration`-Lauf nach der Code-Änderung blieb
  rot, weil das lokal geladene Image (`ghcr.io/pt9912/pg-change-feed:dev`)
  noch den Stand vor dieser Änderung trug — `make image` war nicht erneut
  gelaufen. `harness/README.md`s Werkzeug-Hinweis („Ein Zug, der
  Build-Kontext-Dateien ändert, läuft `make image` vor seiner Closure")
  trägt das bereits; der Implementer-Lauf hat ihn hier zunächst
  übersehen.
- **Beobachtungspfad des sichtbaren Fehlers:** `cdc.heartbeat.error_class`
  (`tools/schema/nacharbeit-heartbeat.sql`, bereits bestehende, für
  `cdc_reader` gegrantete Projektion von `cdc.process_heartbeat`) — derselbe
  externe SQL-Lesezugriffsweg wie `cdc.changes`, kein Container-Log- oder
  Exit-Code-Parsing im Testcode. `awaitHeartbeatErrorClass`
  (`test/integration/integration_test.go`) pollt darauf, `coalesce(error_class, '')`
  trägt den Normalbetrieb (`NULL`) als leere Zeichenkette. Zusätzlich
  bestätigt der Test, dass `cdc.changes` die Zeile `id=11` nicht trägt —
  ihre Transaktion trägt die auslösende Relation-Nachricht vor dem
  eigenen Commit, der Erfassungspfad endet davor.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-032` liegt in `done/`,
`Verantwortlich:` gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich,
  dass eine differenzierte Fehlerbehandlung je Unterfall (entfernte
  Spalte vs. Typänderung vs. Umbenennung) statt einer einheitlichen
  `schema`-Fehlerklasse nötig ist, gehört das zurück zur Zerlegung.
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** `make test-integration`
dreimal in Folge grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Ein sichtbarer `schema`-Fehler bei jeder nicht als Obermenge erkennbaren
  Änderung (auch harmlose Umbenennungen ohne Datenverlust) könnte in der
  Praxis zu häufig auslösen und den Erfassungspfad unnötig hart stoppen —
  eine feinere Unterscheidung wäre eine größere, hier nicht geleistete
  Fähigkeit. **Ausgang: entfallen** — das ist keine neu eingetretene
  Nebenwirkung, sondern eine bereits in §1 bewusst dokumentierte
  Design-Grenze („Änderung der Klassifikationslogik selbst" bleibt
  bestehen): die konservative Binär-Entscheidung entspricht exakt der
  Architect-Skizze (`docs/reviews/architect-verdict-slice-030-adr-0015.md`)
  und ist keine unbeabsichtigte Überreaktion.
- Der bestehende Fehlerbehandlungspfad (wie ein `schema`-Fehler aus
  `Consume` den Erfassungspfad tatsächlich beendet/meldet) könnte
  Annahmen treffen, die für einen während des laufenden Streams
  auftretenden Fehler (statt eines Fehlers beim initialen Decode) nicht
  zutreffen. **Ausgang: entfallen** — Reviewer und Verifier haben den
  vollständigen Propagationspfad (`Consume` → `receive.Stream.Run` →
  `bootstrap.Run` → `os.Exit(1)`) real nachvollzogen und dreifach gegen
  den laufenden Compose-Stack reproduziert: er verhält sich konsistent
  zu bereits bestehenden `schema`-Klassen-Fehlern
  (`ErrTruncateUnsupported` durchläuft denselben Pfad), keine
  falsche Annahme gefunden.
- Der Split der `go test -run`-Filterung in `tools/harness/run-integration-tests.sh`
  (sechs benannte Testfunktionen im vorderen Aufruf, genau
  `TestMVPSchemaChangeIncompatibleTypeChange` im hinteren) kann eine
  künftig zu `test/integration/integration_test.go` hinzugefügte
  Testfunktion, die in keinem der beiden `-run`-Muster auftaucht,
  dauerhaft und stillschweigend von `make test-integration` ausschließen
  — `go test -run` meldet keinen Fehler, solange mindestens eine andere
  Funktion im selben Aufruf matcht. **Ausgang: weiter offen** → wandert
  ins Beobachtungs-Register (`BEO-PGC/test-runner-stiller-ausschluss`,
  neu angelegt).
- Der Split der `go test -run`-Filterung in `tools/harness/run-integration-tests.sh`
  (sechs benannte Testfunktionen im vorderen Aufruf, genau
  `TestMVPSchemaChangeIncompatibleTypeChange` im hinteren) kann eine
  künftig zu `test/integration/integration_test.go` hinzugefügte
  Testfunktion, die in keinem der beiden `-run`-Muster auftaucht,
  dauerhaft und stillschweigend von `make test-integration` ausschließen
  — `go test -run` meldet keinen Fehler, solange mindestens eine andere
  Funktion im selben Aufruf matcht. **Ausgang:** <bei Closure einzutragen>

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

- **Was hat funktioniert:** Der Implementer fand die reale
  Container-Sterblichkeits-Konsequenz des neuen Fehlers selbst über
  einen echten roten Testlauf, nicht durch Vermutung, und behob sie durch
  eine saubere Testreihenfolge-Anpassung statt einer Ad-hoc-Umgehung. Der
  Reviewer verfolgte den vollständigen Fehler-Propagationspfad
  eigenständig bis `os.Exit(1)` und bestätigte die Lösung als korrekt,
  fand aber real ein wiederkehrendes Kommentar-Disziplin-Finding (vierte
  Gelegenheit für dieselbe Fehlerklasse nach `slice-018`/`slice-031`) —
  diesmal in einem Shell-Skript statt Go-Code, was zeigt, dass die Regel
  (`AGENTS.md` §3.7) tatsächlich medienübergreifend gilt und geprüft
  werden muss.
- **Was ging anders als geplant:** Der ursprüngliche Plan unterschätzte
  die Tragweite des neuen sichtbaren Fehlers (Container-weiter Absturz,
  nicht nur ein lokal isolierter Fehlerpfad) — das erzwang eine
  Runner-Skript-Anpassung, die im Plan-Nachzug, nicht im ursprünglichen
  §3-Plan stand. Der Reviewer fand zusätzlich eine strukturelle Lücke im
  Testrunner selbst (F-2), die dieser Slice nicht behebt, sondern als
  Beobachtung weiterreicht.
- **Steering-Loop-Eintrag:** `BEO-PGC/schema-evolution-nicht-dynamisch`
  erreicht mit diesem Slice 2× und geht direkt auf Ausgang *verkörpert*
  (analog `BEO-PGC/spec008-replication-luecke`s Auflösung bei 2×) — die
  `ADR-0015`-Folgepflicht ist real eingelöst. Verkörpert in
  `internal/adapters/driving/replication/mapper/mapper.go`
  (`Assembler.observeRelation`, `classifyRelationColumns`) — liegt in
  `internal/adapters/driving/replication/mapper/mapper.go`. Auslöser:
  `BEO-PGC/schema-evolution-nicht-dynamisch` (slice-030, slice-033 — 2×,
  direkt aufgelöst).
- **Beobachtungs-Register (`../observations/`):**
  `evidence/slice-033.md` in `BEO-PGC/schema-evolution-nicht-dynamisch/`
  ergänzt — Ausgang **verkörpert** (2×). Neu angelegt:
  `BEO-PGC/test-runner-stiller-ausschluss/`, Beleg `evidence/slice-033.md`
  — Zähler 1×, weiter offen.
- **Folge-Slices:** keine — `welle-10` schließt mit diesem Slice; ein
  optionaler schlanker „E2E-Abdeckung — Schema-Evolution"-Folgeschritt
  (Architect-Verdikt-Vorschlag) ist bereits durch die in `slice-030`
  geschriebenen und jetzt real grünen Black-Box-Tests eingelöst, braucht
  keinen eigenen Folge-Slice.
- **Risiken aus §6:** zwei *entfallen*, eines *weiter offen* → siehe §6
  für Begründung.
- **Drei Paarungen:** Repo **mit** Wellen-Betrieb (`welle-10` offen) —
  Prüfung läuft bei der `welle-10`-Closure.

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
Repo-weite Default-Sub-Area `*`/`PGC`.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen.
Treffer für `PGC`: `BEO-PGC/schema-evolution-nicht-dynamisch` (1×, weiter
offen — dieser Slice liefert den dritten und letzten Baustein; die
Auflösung selbst entscheidet der Planner bei `welle-10`s Closure). Keiner
der übrigen Treffer erreicht mit diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
