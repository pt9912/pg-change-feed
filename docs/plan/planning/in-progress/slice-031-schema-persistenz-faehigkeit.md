# Slice slice-031: Schema-Persistenz-Fähigkeit — `TableSchema`-Modell, `SchemaStorePort`, Adapter (ohne Live-Verdrahtung)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-10`](../welle-10.md) — der Nachweis, dass die
`ADR-0015`-Folgepflicht real eingelöst ist, ist `welle-10`s Closure-Trigger
(§3), kein Einzel-Slice-DoD; dieser Slice liefert nur die
Persistenz-Grundlage, ohne die sich `LH-FA-SCH-004`/`005` selbst noch nicht
schließen (das leisten `slice-032`/`slice-033`).

**Bezug:** [`SPEC-004`](../../../../spec/pflichtenheft.md) (technologie-
unabhängiges `TableSchema`/`SchemaVersion`-Modell), `ADR-0015` (nur
gelesen — keine aktive ADR wird geändert, `ADR-0015` bleibt `Accepted`).

**Berührte Spec-Stellen:** [`SPEC-004`](../../../../spec/pflichtenheft.md)
— die Kennung; die Architektur-Sicht zeigt `SchemaStorePort (ARC-004)`
bereits im Sequenzdiagramm zu `LH-FA-CFG-001.a` als vorgesehenen Port.

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

**Ziel:** Die Persistenz-Grundlage für `ADR-0015`s Option C bauen, ohne den
laufenden Erfassungspfad zu berühren: ein `TableSchema`-Domänenmodell
(Spaltenmenge inkl. Typ-/OID-Information je Version), ein neuer Outbound
Port `SchemaStorePort` (aktuelle Schema-Version einer Tabelle lesen, neue
Version registrieren, `TableSchema` zu einer Version lesen) und ein
Postgres-Adapter, der ihn über eine neue Tabelle `cdc.table_schema`
(ausgerollt über d-migrate, `ADR-0043`) real persistiert. Architect-Skizze:
[`docs/reviews/architect-verdict-slice-030-adr-0015.md`](../../../reviews/architect-verdict-slice-030-adr-0015.md)
§Umsetzungsskizze, Schritte 1–3.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Verdrahtung in `Assembler.Consume`** (dynamische Re-Versionierung bei
  echten `*decode.Relation`-Ereignissen) — Folge-Slice `slice-032`
  übernimmt das; dieser Slice liefert nur die Fähigkeit, keine
  Verhaltensänderung am laufenden Erfassungspfad.
- **Spalten-Typ-/OID-Auswertung im Decoder und Fehlerklasse `schema` für
  inkompatible Typänderungen** — Folge-Slice `slice-033` übernimmt das;
  `TableSchema` hält hier nur die Datenstruktur, keine
  Vergleichs-/Kompatibilitätslogik.
- **Änderung der bestehenden `EnableTable`-Erstaktivierung** — Bestand
  bleibt bewusst stehen: `EnableTableService`/`InsertSchemaVersion` bleiben
  unverändert (Version 1 bei Erstaktivierung bleibt korrekt, siehe
  Architect-Skizze Schritt 6); der neue Port wird nur additiv verdrahtet,
  nicht in den bestehenden Aktivierungspfad eingebaut.
- **CLI-/Verwaltungs-Exposition der Schema-Historie** — kein Lastenheft-
  Kriterium verlangt das in diesem Umfang; anderer Vorgang, falls je
  gebraucht.

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

- [x] `TableSchema`-Domänenmodell (`internal/domain/model/`) trägt eine
      Spaltenmenge mit Typ-/OID-Information je `SchemaVersionID`; Unit-
      getestet (`make test`). Beleg: `internal/domain/model/table_schema.go`,
      `internal/domain/model/table_schema_test.go` — `make test` grün.
- [x] `SchemaStorePort` (`internal/application/port/outbound/`) neu
      definiert (`CurrentVersion`, `RegisterVersion`, `TableSchema` lesen
      — exakte Signatur Implementer-Entscheidung im Plan-Nachzug); ein
      Postgres-Adapter implementiert ihn real gegen eine neue Tabelle
      `cdc.table_schema` (`tools/schema/schema.yaml`, ausgerollt über
      d-migrate). Real gegen PostgreSQL getestet (`make test-store`-Muster):
      Schreiben, Lesen, Round-Trip. Beleg:
      `internal/application/port/outbound/schemastore.go`,
      `internal/adapters/driven/postgresstorage/schemastore.go`,
      `internal/adapters/driven/postgresstorage/schemastore_test.go`,
      `tools/schema/schema.yaml` (neue Tabelle `table_schema`) — `make
      test-store` grün (5 neue Tests, real gegen PostgreSQL im
      Testcontainer, Schema ausgerollt über d-migrate/`make schema-rollout`).
- [x] `make gates` grün. Beleg: `baseline-verify` OK (54 Dateien),
      `docs-check` 268 Dateien/0 Befunde, `commit-traceability` OK (5
      Commits, Range `HEAD~5..HEAD`), `a-check` 0 Befunde.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Beleg: [`docs/reviews/review-slice-031.md`](../../../reviews/review-slice-031.md)
      (2 HIGH, 1 MEDIUM), Fixrunde behoben in Commit `100ff2b`, bestätigt
      in [`docs/reviews/review-slice-031-fixrunde.md`](../../../reviews/review-slice-031-fixrunde.md)
      (alle drei Findings behoben, keine Regression). Verifikation in
      [`docs/reviews/verify-slice-031.md`](../../../reviews/verify-slice-031.md)
      (DoD eigenständig nachgeprüft, real reproduziert).
- [x] Doku-Update für `harness/README.md` §Sensors/`AGENTS.md`, falls ein
      neuer Sensor/Vertrag entsteht — Implementer entscheidet und begründet
      im Plan-Nachzug. **Begründung:** kein neuer Sensor/Vertrag entstanden
      — `make test-store` existiert bereits als dokumentierter Sensor
      (`harness/README.md` §Sensors), die neuen Tests erweitern ihn nur;
      kein neues `make`-Target.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. Siehe §7.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Entfällt: Repo ist Greenfield (`harness/conventions.md` §Modus-Deklaration), `../reconciliation.md` existiert nicht.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
      Keine Beobachtung angefallen — dieser Slice liefert nur einen
      Baustein zur Auflösung von `BEO-PGC/schema-evolution-nicht-dynamisch`
      (1×, unverändert); die Auflösung selbst erfolgt erst mit `slice-032`/
      `slice-033`/`welle-10`-Closure, nicht schon hier.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen). Siehe §6.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Entfällt hier: Repo mit Wellen-Betrieb — Prüfung läuft bei der `welle-10`-Closure.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/domain/model/table_schema.go` | neu | `TableSchema`-Domänenmodell: Spaltenmenge inkl. Typ-/OID je Version |
| `internal/application/port/outbound/schemastore.go` | neu | `SchemaStorePort`-Vertrag |
| `internal/adapters/driven/postgresstorage/` | update | Postgres-Adapter für `SchemaStorePort` |
| `tools/schema/schema.yaml` | update | neue Tabelle `cdc.table_schema` |
| `internal/bootstrap/wiring.go` | update | `SchemaStorePort`-Adapter verdrahten (additiv, keine Verhaltensänderung am laufenden Pfad) |

### Plan-Nachzug (nach Implementierung)

- **`TableSchema`-Form:** `Column{ Name string; OID ColumnOID }` (`ColumnOID
  uint32`) und `TableSchema{ VersionID SchemaVersionID; Columns []Column }`
  (`internal/domain/model/table_schema.go`). Die OID bleibt roh (keine
  Übersetzung in einen technologieunabhängigen Typ-Enum) — dieselbe
  Abgrenzung wie beim Architect-Verdikt (Schritt 4 der Skizze trennt
  Decoder-OID-Mitführung von der Typ-Kompatibilitätsprüfung, `slice-033`).
  `NewTableSchema` verlangt eine nichtleere `VersionID`, mindestens eine
  Spalte (`ErrEmptyColumns`, neu in `internal/domain/errors`) und
  nichtleere Spaltennamen; die Spaltenmenge wird beim Konstruieren
  kopiert. Kein separates Ordnungsfeld auf `Column` — die
  Spalten-Reihenfolge trägt die Position im Slice und (persistiert) die
  neue Spalte `ordinal_position`.
- **`SchemaStorePort`-Signatur:** `CurrentVersion(ctx, table
  model.SourceTableID) (model.SchemaVersion, bool, error)` (höchste
  registrierte Version; `bool` false = keine über diesen Port
  registrierte Version — die statische Erstaktivierung schreibt ihre
  Version-1-Zeile weiterhin über `TableActivationPort.Register`, siehe
  Ausschluss §1), `RegisterVersion(ctx, version model.SchemaVersion,
  schema model.TableSchema) (bool, error)` (Schema-Version- und
  TableSchema-Zeilen in einem Store-Commit; `version.ID != schema.VersionID`
  endet vor dem ersten SQL-Aufruf über den neuen Sentinel
  `outbound.ErrSchemaVersionMismatch`) und `TableSchema(ctx, versionID
  model.SchemaVersionID) (model.TableSchema, error)` (Abwesenheit über
  `outbound.ErrSchemaVersionUnknown` sichtbar, keine stille
  Fehlinterpretation, `LH-FA-SCH-004`). Fehlerklasse `storage` über den
  neuen Sentinel `outbound.ErrSchemaStoreStorage` — derselbe Aufbau wie
  `ErrConsumerStateStorage` (eigener Sentinel je Port, keine geteilte
  Klasse-Aktion mit dem ChangeStore-Sentinel).
- **Adapter-Ort:** eigene Datei
  `internal/adapters/driven/postgresstorage/schemastore.go`
  (`PostgresSchemaStoreAdapter`, `NewSchemaStore`) statt Erweiterung von
  `tableactivation.go` — eigenständiger Port mit eigenem Lebenszyklus
  (`Close`), analog zu `consumerstate.go`/`ConsumerStatePort`, nicht
  Erweiterung eines bestehenden Adapters.
- **Tabellenform `cdc.table_schema`:** pro Spalte eine Zeile
  (`schema_version_id`, `ordinal_position`, `column_name`, `column_oid`),
  PK `(schema_version_id, ordinal_position)`, UNIQUE
  `(schema_version_id, column_name)`; `column_oid` als `biginteger` (eine
  PostgreSQL-OID ist ein vorzeichenloses 32-Bit-Feld und passt nicht
  verlustfrei in ein vorzeichenbehaftetes `integer`). Nur in
  `tools/schema/schema.yaml` (d-migrate-Rollout) — nicht in der
  handgeschriebenen `schema.sql`/`ApplySchema`, dieselbe Abgrenzung wie
  bei den Consumer-State-Tabellen (`ADR-0043`); Adapter-Test folgt dem
  `consumerstate_test.go`-Muster (Tabellen-Existenz prüfen, kein `DROP
  SCHEMA … CASCADE`) statt dem `tableactivation_test.go`-Muster, weil
  `cdc.table_schema` nur über den d-migrate-Rollout entsteht.
- **Wiring — entfällt bewusst:** `internal/bootstrap/wiring.go` bleibt in
  diesem Slice unverändert. Keine bestehende Verdrahtungsstelle in `Run`
  ruft `SchemaStorePort` auf (Verdrahtung in `Assembler.Consume` ist
  ausdrücklich `slice-032`); eine reine Objekt-Konstruktion ohne Aufrufer
  würde nur eine zusätzliche, ungenutzte DB-Verbindung im
  Produktionspfad öffnen und schließen — ein neuer Fehlermodus
  (Erreichbarkeits-Fehler eines Ports, den dieser Lauf nicht braucht)
  ohne Gegenwert. Die §3-Plan-Zeile zu `wiring.go` oben ist damit nicht
  eingelöst; das ist eine bewusste Abweichung vom Vorab-Plan, keine
  stillschweigende.

### Fixrunde nach Review (`docs/reviews/review-slice-031.md`, F-1/F-2/F-3)

- **F-1 (HIGH, Kommentar-Chronik):** Die Vorwärtsverweise auf `slice-032`/
  `slice-033` und der Review-Report-Pfad in den Doc-Kommentaren von
  `internal/domain/model/table_schema.go`, `internal/application/port/
  outbound/schemastore.go` und `internal/adapters/driven/postgresstorage/
  schemastore.go` sind entfernt; die Kommentare beschreiben jetzt nur noch
  den Ist-Zustand von Modell, Port und Adapter (`AGENTS.md` §3.7). Die
  Out-of-Scope-/Folgepflicht-Information bleibt vollständig in §1 dieses
  Plans stehen — das ist ihr Ort.
- **F-2 (HIGH, falscher Doc-Kommentar):** Der `CurrentVersion`-Kommentar in
  `postgresstorage/schemastore.go` behauptete, Version 1 der statischen
  Erstaktivierung werde „erst nach der ersten `RegisterVersion`" sichtbar.
  Eigenständig gegen eine reale PostgreSQL-Instanz nachvollzogen
  (Aktivierung über `TableActivationAdapter.Register`, danach
  `CurrentVersion` ohne `RegisterVersion`-Aufruf): meldet sofort `ok=true,
  version=1`, weil beide Adapter dieselbe `cdc.schema_version`-Tabelle
  beschreiben. Kommentar auf dieses tatsächliche Verhalten korrigiert.
- **F-3 (MEDIUM, Backfill-Lücke):** `RegisterVersion` übersprang die
  Spalten-Einfügung immer dann, wenn die `schema_version`-Zeile bereits
  bestand (`registered == false`) — für Version 1 nach jeder
  `EnableTable`-Aktivierung der Fall, da diese dieselbe Zeile über
  `queries.InsertSchemaVersion` schreibt. Ein späteres Nachtragen der
  Spaltenform für exakt diese Version blieb dadurch wirkungslos
  (`TableSchema` dauerhaft `ErrSchemaVersionUnknown`). Fix: Die
  Spalten-Einfügung hängt jetzt an der Abwesenheit einer bestehenden
  Spaltenform selbst (neue Abfrage `queries.CountTableSchemaColumns`),
  nicht mehr an der Neuheit der `schema_version`-Zeile — eine
  `SchemaVersionID` ohne Spaltenform bekommt sie geschrieben (Rückkehr
  `true`), unabhängig davon, ob ihre Versions-Zeile neu ist oder über
  einen anderen Schreibpfad (`TableActivationAdapter.Register`) bereits
  bestand; eine Version mit bestehender Spaltenform bleibt unverändert
  (Rückkehr `false`, keine doppelte Spaltenform). Neuer Testfall
  `TestSchemaStoreRegisterVersionBackfillsColumns`
  (`schemastore_test.go`): simuliert die bestehende Schema-Version-Zeile,
  registriert die Spaltenform nach, liest sie über `TableSchema` zurück,
  prüft die Idempotenz einer erneuten Registrierung — real gegen
  PostgreSQL, `make test-store` grün.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Priorisiert, `Verantwortlich:` gesetzt,
WIP-Limit (1 je Implementer) frei — keine harte Abhängigkeit von einem
anderen Slice.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich, dass
  Domänenmodell, Port und Adapter zusammen mehr als drei Liefer-Punkte
  oder mehr als zwei Schichten in einer Review-Sitzung nicht mehr prüfbar
  machen, gehört das zurück zur Zerlegung (z. B. Adapter als eigener
  Folge-Slice).
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** Adapter-Test gegen reale
PostgreSQL grün **und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Die genaue Form von `TableSchema` (welche Typ-Repräsentation trägt eine
  Spalte — PostgreSQL-OID roh oder ein bereits übersetzter,
  technologieunabhängiger Typ-Enum) könnte sich erst in `slice-033`
  (Typ-Auswertung) als unpassend erweisen und eine Nacharbeit an diesem
  Slice erzwingen. **Ausgang: entfallen** — der Reviewer hat die Form
  (`Column{Name, OID}`) eigenständig als Grundlage für `slice-033`
  geprüft (`docs/reviews/review-slice-031.md`, Negativbefund): eine
  PostgreSQL-OID identifiziert den exakten Datentyp eindeutig und ist
  ausreichendes Rohmaterial für eine spätere Typ-Kompatibilitätsprüfung.
- Eine neue Tabelle `cdc.table_schema` im neutralen Schema könnte mit
  d-migrates View-Signatur-Handling (`ADR-0043`, bekannte Historie bei
  `views:`) in Konflikt geraten, obwohl es sich um eine reguläre Tabelle
  handelt. **Ausgang: entfallen** — der reale Rollout
  (`make schema-rollout`, regenerierte `plan.yaml`/`down.sql`) verlief
  ohne Konflikt; das bekannte View-Signatur-Problem betrifft
  ausschließlich `views:`-Knoten, `table_schema` ist eine reguläre
  Tabelle und dadurch strukturell nicht betroffen — vom Reviewer
  eigenständig anhand der erzeugten DDL bestätigt.

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

- **Was hat funktioniert:** Der Reviewer fand real, nicht nur behauptet,
  zwei HIGH-Findings (Kommentar-Klassen-Verstöße, dieselbe Fehlerklasse
  wie bereits einmal in `slice-018` korrigiert) und ein MEDIUM (eine
  echte API-Lücke in `RegisterVersion`, die `slice-032` sonst blockiert
  hätte) — durch eigene, empirische Proben gegen reale PostgreSQL, nicht
  durch bloße Code-Lektüre. Die anschließende Fixrunde und ihre
  Bestätigung liefen sauber getrennt (Implementer → Reviewer-Bestätigung
  → Verifier), ohne Architect-Eskalation, da keine ADR-/Plan-Frage
  vorlag.
- **Was ging anders als geplant:** Zwei Fixrunden-Findings (F-1, F-2)
  betrafen reine Kommentar-Formulierung, nicht die Fähigkeit selbst — ein
  wiederkehrendes Muster (dieselbe Fehlerklasse wie `slice-018`s
  Kommentar-Chronik-Fund). F-3 deckte auf, dass `RegisterVersion`
  ursprünglich keine Spaltenform für eine bereits per `EnableTable`
  existierende `SchemaVersionID` nachtragen konnte — behoben, bevor
  `slice-032` darauf hätte aufbauen müssen.
- **Steering-Loop-Eintrag:** Kein Eintrag erreicht mit diesem Slice 3× —
  der Normalfall. `BEO-PGC/schema-evolution-nicht-dynamisch` bleibt
  unverändert bei 1× (dieser Slice liefert einen Baustein, keinen neuen
  Beleg).
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung
  angefallen — siehe §2-Begründung.
- **Folge-Slices:** keine — `slice-032` (dynamische Re-Versionierung) und
  `slice-033` (Typ-Auswertung/Fehlerklasse `schema`) stehen bereits in
  `welle-10` §4 als vorgesehene Slices dieser Welle, werden aber noch
  nicht als Dateien angelegt (Modul 5: nicht alle Slices vor der ersten
  Implementation planen) — `slice-032` wird als nächster Schritt dieser
  Welle neu geschnitten, nicht durch diesen Closure-Schritt ausgelöst.
- **Risiken aus §6:** beide *entfallen* — siehe §6 für Begründung.
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
offen — dieser Slice liefert den ersten Baustein der Auflösung, ohne sie
selbst zu vollenden). Keiner der übrigen Treffer erreicht mit diesem Slice
3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
