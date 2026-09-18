# Slice slice-030: Black-Box-E2E für Schema-Änderungen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-9`](welle-9.md) — der Nachweis, dass
Schema-Änderungen real gegen den Compose-Stack erfasst und gelesen werden,
ist welle-9s Closure-Trigger (§3), kein Einzel-Slice-DoD.

**Bezug:** [`LH-FA-SCH-001`](../../../../spec/lastenheft.md),
[`LH-FA-SCH-002`](../../../../spec/lastenheft.md),
[`LH-FA-SCH-004`](../../../../spec/lastenheft.md),
[`LH-FA-SCH-005`](../../../../spec/lastenheft.md) — dieser Slice erweitert
die Testabdeckung für bestehende Verträge, ändert sie nicht. Kein aktives
ADR wird geändert.

**Berührte Spec-Stellen:** — (reine Testinfrastruktur, keine
Verhaltensänderung an einer Spec-Stelle).
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

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

**Ziel:** Der Compose-Integrationstest bekommt neue Testfälle für
Schema-Änderungen an einer aktivierten Tabelle — ausschließlich über
externe Schnittstellen (SQL-DDL auf die Quelltabelle, Lesen über
`cdc.changes`): (a) `ALTER TABLE ADD COLUMN` wird erfasst, die
`schema_version` erhöht sich, ältere Changes bleiben unverändert lesbar
(`LH-FA-SCH-001`/`002`/`005`); (b) eine inkompatible Typänderung wird
sichtbar als Fehler gemeldet, nicht still fehlinterpretiert
(`LH-FA-SCH-004`, Negative-Fall).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Entfernte Spalten** (`LH-FA-SCH-003`) — anderer Vorgang: eigener
  Testfall mit eigener Akzeptanzkriterien-Prüfung, kein Bestandteil dieses
  ≤3-Liefer-Punkte-Slices; Kandidat für einen Folge-Slice, falls dieser
  hier zu groß wird.
- **Änderung des Schema-Änderungs-Erkennungsmechanismus selbst** — Bestand
  bleibt bewusst stehen: `internal/domain/model/schema_version.go` (oder
  Äquivalent) wird nur *geprüft*, nicht geändert.
- **Lesepfad-Vertragstest (SQL-View vs. Go-Adapter)** — Folge-Slice
  `slice-029` liefert das bereits; dieser Slice nutzt denselben externen
  Lesezugriffsweg, erfindet aber kein zweites Vergleichsmuster.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

**Planner-Korrektur bei Closure (analog `slice-028`s F-1-Präzedenz):** Das
ursprüngliche Ziel oben unterstellte, `schema_version` erhöhe sich real und
eine inkompatible Typänderung werde real als Fehlerklasse `schema`
gemeldet — beides unterstellte eine bereits vorhandene Fähigkeit, die
lediglich noch keine Testabdeckung hatte. Real geliefert (und das ist der
eigentliche Wert dieses black-box-testenden Slices) ist etwas anderes und
ebenso Vollständiges: **der empirische Nachweis am realen Compose-Stack,
was das System tatsächlich tut** — inklusive des Fundes, dass zwei
Lastenheft-Akzeptanzkriterien (`LH-FA-SCH-004` Negative-Fall,
`LH-FA-SCH-005` Boundary) strukturell nicht erfüllt sind, weil eine
`ADR-0015`-Folgepflicht nie eingelöst wurde (Reviewer-Fund F-1 HIGH,
Architect-Verdikt in
[`docs/reviews/architect-verdict-slice-030-adr-0015.md`](../../../reviews/architect-verdict-slice-030-adr-0015.md)).
Die tatsächliche Fähigkeits-Lieferung (`SchemaStorePort`, dynamische
Re-Versionierung, Typ-Fehlerklasse) ist **nicht** Gegenstand dieses
Testinfrastruktur-Slices — sie ist als eigene Feature-Welle
„Schema-Evolution-Nachlieferung (`ADR-0015`)" in der Roadmap vorgemerkt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] `LH-FA-SCH-001`/`002` real belegt, `LH-FA-SCH-005` als realer Fund
      dokumentiert (korrigierte Fassung, siehe §1 Planner-Korrektur — das
      ursprüngliche Kriterium unterstellte eine bereits vorhandene
      dynamische Schema-Versionierung).
      `LH-FA-SCH-001`/`002` real belegt
      (`TestMVPSchemaChangeAddColumn`, `test/integration/integration_test.go`)
      — neue Spalte im Row Image danach erfasster Changes erkennbar, ältere
      Change unverändert/ohne die neue Spalte lesbar. `LH-FA-SCH-005`s
      Unterscheidbarkeits-Boundary **strukturell nicht erfüllbar** — real
      nachgewiesener Fund (statische Schema-Version-Bindung, kein
      Metadata-Pfad implementiert), siehe §3 Plan-Nachzug, Reviewer F-1
      (HIGH, `docs/reviews/review-slice-030.md`) und Architect-Verdikt
      (`docs/reviews/architect-verdict-slice-030-adr-0015.md`). Ausgang:
      `BEO-PGC/schema-evolution-nicht-dynamisch` registriert, Feature-Welle
      „Schema-Evolution-Nachlieferung (`ADR-0015`)" in Roadmap vorgemerkt.
- [x] `LH-FA-SCH-004` real geprüft, Negative-Fall als realer Fund
      dokumentiert (korrigierte Fassung, siehe §1 Planner-Korrektur).
      Der PostgreSQL-seitige Ablehnungsfall real belegt
      (`TestMVPSchemaChangeIncompatibleTypeChange`). Der CDC-seitige
      Negative-Fall (PostgreSQL lässt zu, CDC meldet Fehlerklasse `schema`)
      ist mit dem aktuellen System **strukturell nicht herstellbar** — real
      nachgewiesener Fund (Text-Pass-through ohne Typprüfung), derselbe
      `ADR-0015`-Zusammenhang wie oben.
- [x] `make gates` grün, `make test-integration` dreimal in Folge grün. —
      `make gates`: `d-check` 259 Dateien / 0 Befund(e), `a-check` 0
      Befund(e), `commit-traceability` OK, `baseline-verify` OK (Lauf vor
      dem Plan-Nachzug-Commit dieses Slice). `make test-integration`: 3×
      in Folge grün, inklusive der beiden neuen Testfunktionen (Commit
      folgt).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Beleg: [`docs/reviews/review-slice-030.md`](../../../reviews/review-slice-030.md)
      (1 HIGH F-1, 1 LOW F-2) plus Architect-Verdikt
      [`docs/reviews/architect-verdict-slice-030-adr-0015.md`](../../../reviews/architect-verdict-slice-030-adr-0015.md)
      (Modul 8 Konflikt-Pfad, Verdikt 1: `ADR-0015` gilt unverändert fort).
- [x] Doku-Update, falls ein öffentlicher Vertrag berührt wird — hier
      voraussichtlich keiner; Implementer entscheidet und begründet im
      Plan-Nachzug. **Begründung:** kein öffentlicher Vertrag berührt — die
      dritte `CDC_TABLES`-Bindung ist test-lokal (`compose.yaml`,
      `run-integration-tests.sh`); `docs/user/benutzerhandbuch.md`
      dokumentiert das `CDC_TABLES`-Format generisch, ohne die konkreten
      Compose-Tabellennamen zu nennen, und bleibt unverändert korrekt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. Siehe §7.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Entfällt: Repo ist Greenfield (`harness/conventions.md` §Modus-Deklaration), `../reconciliation.md` existiert nicht.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
      Beleg: `BEO-PGC/schema-evolution-nicht-dynamisch/` neu angelegt,
      Beleg `evidence/slice-030.md` — Zähler 1×. Siehe §7.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen). Siehe §6.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Entfällt hier: Repo mit Wellen-Betrieb — Prüfung läuft bei der `welle-9`-Closure.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/run-integration-tests.sh` oder `test/integration/integration_test.go` | update | neue Testfälle: `ALTER TABLE ADD COLUMN` real erfasst, `schema_version`-Wechsel geprüft, ältere Changes weiter lesbar; inkompatible Typänderung sichtbar als Fehler |
| `compose.yaml` | update | *Plan-Nachzug:* dritte aktivierte Testtabelle `public.feed_mvp_schema=tbl-mvp-schema:sv-mvp-schema` in `CDC_TABLES` ergänzt — eigene, von den bestehenden MVP-Referenztabellen isolierte Tabelle für die Schema-Änderungs-Testfälle. |
| `tools/harness/run-integration-tests.sh` | update | *Plan-Nachzug:* `CREATE TABLE public.feed_mvp_schema (id int PRIMARY KEY, name text, amount text)` vor dem Feed-Container-Start ergänzt (Vorbedingung der Aktivierung, wie bei den bestehenden `feed_mvp_*`-Tabellen). |
| `test/integration/integration_test.go` | update | *Plan-Nachzug:* zwei Testfunktionen ergänzt — `TestMVPSchemaChangeAddColumn` (`LH-FA-SCH-001` Happy Path, `LH-FA-SCH-002` Boundary, real gegen den Compose-Stack: neue Spalte im Row Image danach erfasster Changes, ältere Change unverändert/ohne die neue Spalte lesbar) und `TestMVPSchemaChangeIncompatibleTypeChange` (`LH-FA-SCH-004`: PostgreSQL lehnt eine tatsächlich inkompatible Typänderung selbst per DDL ab — Boundary-Fall aus §6 real belegt; eine PostgreSQL-seitig zugelassene Typänderung mit konvertierbaren Bestandsdaten wird von CDC unverändert übernommen). **Fund (beide Testfälle, empirisch am realen Compose-Stack verifiziert, nicht nur aus Code-Lektüre geschlossen):** (1) `LH-FA-SCH-005`s Boundary — unterscheidbare Schema-Versionen vor/nach einer Schemaänderung — hält mit dem aktuellen System nicht: `mapper.TableBinding.SchemaVersion` (`internal/adapters/driving/replication/mapper/mapper.go`) ist eine bei der Aktivierung aus `CDC_TABLES` statisch gebundene Kennung, unverändert für die Laufzeit des Feed-Containers; `mapper.Assembler.Consume` verwirft `*decode.Relation`-Ereignisse ungenutzt (Default-Zweig). Die an vier Stellen referenzierte dynamische Re-Versionierung „über den Metadata-Pfad" (`internal/bootstrap/wiring.go:331`, `internal/adapters/driving/replication/mapper/mapper.go:53`, `internal/adapters/driving/replication/receive/receive.go:62`, `internal/application/port/inbound/verwaltung.go:26`) hat im Repo keine Implementierung. (2) `LH-FA-SCH-004`s Negative-Fall — PostgreSQL lässt eine Typänderung zu, CDC decodiert sie aber nicht verlustfrei und meldet das sichtbar (Fehlerklasse `schema`) — ist mit dem aktuellen System nicht real herstellbar: `decode.tupleValues`/`mapper.rowImage` interpretieren jeden Spaltenwert ausschließlich als Text ohne eigene Typprüfung; eine „nicht sicher interpretierbare Schemaänderung" (`LH-FA-SCH-004.a`) kann daraus für einen von PostgreSQL zugelassenen Typwechsel nicht entstehen. Beide Funde ändern den Erkennungsmechanismus nicht (Out-of-Scope §1) — die Tests belegen den realen Ist-Zustand inklusive dieser Lücken, statt sie zu verdecken. DoD-Punkte entsprechend nicht als erfüllt markiert; Einordnung (Carveout/Folge-Slice/Beobachtung) bleibt der Planner-Closure vorbehalten. |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Priorisiert, `Verantwortlich:` gesetzt,
WIP-Limit (1 je Implementer) frei — keine harte Abhängigkeit von einem
anderen Slice.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  entfernte Spalten (`LH-FA-SCH-003`) oder weitere Typänderungs-Varianten
  ebenfalls in diesen Slice gehören sollen, gehört das zurück zur
  Zerlegung — die ursprüngliche Abgrenzung (§1) bleibt sonst nur auf dem
  Papier.
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

- Was genau als „inkompatible Typänderung" gilt, könnte von PostgreSQLs
  eigenem `ALTER COLUMN … TYPE`-Verhalten abweichen (PostgreSQL lehnt
  manche inkompatible Änderungen bereits selbst per DDL-Fehler ab, bevor
  CDC sie überhaupt sieht) — der Testfall muss einen Fall wählen, den
  PostgreSQL zulässt, aber CDC nicht verlustfrei decodieren kann.
  **Ausgang: weiter offen** → das Risiko ist eingetreten, aber in
  verschärfter Form: Es gibt aktuell **keinen** Fall, den PostgreSQL
  zulässt und CDC nicht verlustfrei decodiert, weil CDC gar keine
  Typprüfung vornimmt — kein Testfall-Auswahlproblem, sondern eine
  strukturelle Lücke. Wandert ins Beobachtungs-Register
  (`BEO-PGC/schema-evolution-nicht-dynamisch`), adressiert durch die
  vorgemerkte Feature-Welle „Schema-Evolution-Nachlieferung (`ADR-0015`)".
- Der bestehende Erkennungsmechanismus für Schemaänderungen könnte
  bestimmte DDL-Formen (z. B. `ALTER TABLE … ALTER COLUMN … TYPE` versus
  `DROP`+`ADD`) unterschiedlich behandeln, was der Testfall zufällig nicht
  trifft. **Ausgang: weiter offen** → real gefunden wurde etwas
  Grundlegenderes als DDL-Form-Inkonsistenz: **jede** Relation-Metadata-
  Änderung wird im `Assembler.Consume`-Default-Zweig verworfen, unabhängig
  von der DDL-Form. Dieselbe Beobachtung
  (`BEO-PGC/schema-evolution-nicht-dynamisch`) deckt das ab; kein
  gesondertes Risiko nötig.

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

- **Was hat funktioniert:** Der schwarze-Kasten-Ansatz dieser Welle (real
  gegen den laufenden Compose-Stack testen, statt Code zu lesen und zu
  vertrauen) hat genau das geleistet, wofür er gebaut wurde: einen realen,
  seit `ADR-0015`s Verabschiedung (2026-09-09) unentdeckten Verstoß gegen
  eine `Accepted`/`permanent` ADR-Folgepflicht aufzudecken — nicht durch
  Code-Lektüre vermutet, sondern empirisch am echten System bewiesen.
- **Was ging anders als geplant:** Der Implementer stieß bei beiden
  DoD-Kriterien auf real nicht herstellbare Fälle und dokumentierte das
  ehrlich als offene Checkbox statt eine Scheinabdeckung zu bauen. Der
  Reviewer stufte den Fund als HIGH (ADR-Verstoß) ein und eskalierte an den
  Architect (Modul 8 Konflikt-Pfad) statt ihn selbst zu bewerten. Das
  Architect-Verdikt bestätigte Verdikt 1: `ADR-0015` gilt unverändert
  fort, die Folgepflicht (`SchemaStorePort`, dynamische Re-Versionierung)
  wurde nie eingeplant — kein früherer Slice hat sie fälschlich als
  geliefert behauptet. Konsequenz: neue Feature-Welle
  „Schema-Evolution-Nachlieferung (`ADR-0015`)" in der Roadmap
  vorgemerkt (Größe L, ≥3 Slices laut Architect-Skizze), unmittelbar nach
  `welle-9` eingereiht — vor der bereits geplanten „E2E-Abdeckung —
  Verwaltung & Observability".
- **Steering-Loop-Eintrag:** Kein Eintrag erreicht mit diesem Slice 3× —
  der Normalfall. `BEO-PGC/schema-evolution-nicht-dynamisch` neu angelegt,
  steht bei 1×.
- **Beobachtungs-Register (`../observations/`):** `BEO-PGC/schema-evolution-nicht-dynamisch/`
  neu angelegt, Beleg `evidence/slice-030.md` — Zähler steht bei 1×.
- **Folge-Slices:** keine konkrete Slice-Datei — die Fähigkeits-Lieferung
  ist als eigene Feature-Welle „Schema-Evolution-Nachlieferung (`ADR-0015`)"
  in der Roadmap (*Nächste Wellen*) vorgemerkt, ihre Slices werden bei
  deren Eröffnung geschnitten (Modul 5: nicht alle Slices vor der ersten
  Implementation planen).
- **Risiken aus §6:** beide *weiter offen* → `BEO-PGC/schema-evolution-nicht-dynamisch`
  — siehe §6 für Begründung.
- **Drei Paarungen:** Repo **mit** Wellen-Betrieb (`welle-9` offen) —
  Prüfung läuft bei der `welle-9`-Closure.

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
Keine Treffer für `PGC`, die spezifisch die Schema-Änderungs-Behandlung
betreffen.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
