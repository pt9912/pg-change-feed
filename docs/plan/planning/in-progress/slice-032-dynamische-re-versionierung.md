# Slice slice-032: Dynamische Re-Versionierung im Consume-Pfad

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-10`](../welle-10.md) — der Nachweis, dass die
`ADR-0015`-Folgepflicht real eingelöst ist, ist `welle-10`s Closure-Trigger
(§3), kein Einzel-Slice-DoD; dieser Slice schließt `LH-FA-SCH-005`s
Boundary real, `LH-FA-SCH-004`s Negative-Fall bleibt `slice-033`.

**Bezug:** [`LH-FA-SCH-005`](../../../../spec/lastenheft.md),
[`LH-FA-SCH-001`](../../../../spec/lastenheft.md)/`002` (bereits erfüllt,
hier real geschlossen statt nur teilweise), `ADR-0015` (nur gelesen —
keine aktive ADR wird geändert, `ADR-0015` bleibt `Accepted`).

**Berührte Spec-Stellen:** [`SPEC-004`](../../../../spec/pflichtenheft.md),
`LH-FA-SCH-004.a` (der bislang unbelegte „Metadata-Pfad").

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

**Ziel:** `Assembler.Consume`
(`internal/adapters/driving/replication/mapper/mapper.go`) behandelt
`*decode.Relation`-Ereignisse real, statt sie im `default`-Zweig zu
verwerfen: bei unveränderter Spaltenform (Name **und** Typ-OID identisch
zur zuletzt bekannten `TableSchema`) geschieht nichts; bei einer
**kompatiblen Erweiterung** (mindestens eine neue Spalte, alle
bestehenden Spalten unverändert) wird über den `SchemaStorePort`
(`slice-031`) eine neue `SchemaVersionID` registriert und die
`TableBinding` des Assemblers auf diese Version aktualisiert — künftig
erfasste Changes referenzieren die neue Version, ältere bleiben
unverändert über die alte lesbar (`LH-FA-SCH-005`s Boundary, real belegt
durch `TestMVPSchemaChangeAddColumn`/`TestMVPSchemaChangeIncompatibleTypeChange`
aus `slice-030`, die dieser Slice ohne Konzeptänderung grün machen soll).
Dazu muss der Decoder
(`internal/adapters/driving/replication/decode/decode.go`) die
Spalten-Typ-OID einer `pgoutput`-Relation-Message mitführen (`Column`
trägt bislang nur `Name`/`Key`) — ohne diese Information lässt sich
„unverändert" von „kompatible Erweiterung" nicht zuverlässig von
„inkompatible Änderung" unterscheiden.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Erkennung/Meldung inkompatibler Änderungen als Fehlerklasse `schema`**
  (`LH-FA-SCH-004`s Negative-Fall) — Folge-Slice `slice-033` übernimmt
  das explizit. Dieser Slice trifft **nur** die Entscheidung „unverändert
  vs. kompatible Erweiterung"; jede Relation-Änderung, die **keine**
  reine Obermengen-Erweiterung ist (Spalte entfernt, Spaltentyp einer
  bestehenden Spalte geändert, Spalte umbenannt), bleibt in diesem Slice
  bei der bisherigen, konservativen Behandlung (keine neue Version
  registriert, `TableBinding` unverändert) — kein sichtbarer Fehler, kein
  stiller Fortschritt zu falschen Daten; das ist bewusst ein Zwischenstand,
  den `slice-033` real schließt.
- **Änderung der bestehenden `EnableTable`-Erstaktivierung** — Bestand
  bleibt bewusst stehen: `EnableTableService`/`InsertSchemaVersion`
  bleiben unverändert (Version 1 bei Erstaktivierung bleibt korrekt).
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

- [x] Decoder trägt die Spalten-Typ-OID einer `pgoutput`-Relation-Message
      (`decode.Column` erweitert); unit-getestet. Beleg:
      `internal/adapters/driving/replication/decode/decode.go` (`Column.TypeOID`,
      `relationColumns`), Test `TestDecodeRelationColumnTypeOID`
      (`decode_test.go`).
- [x] `Assembler.Consume` behandelt `*decode.Relation` real: unverändert
      → no-op; kompatible Erweiterung → neue `SchemaVersionID` über
      `SchemaStorePort` registriert, `TableBinding` aktualisiert. Unit-
      getestet (Assembler-Tests mit einem Fake/Stub `SchemaStorePort`).
      Beleg: `internal/adapters/driving/replication/mapper/mapper.go`
      (`observeRelation`, `classifyRelationColumns`), Tests
      `TestConsumeRelationUnchanged`, `TestConsumeRelationCompatibleExtension`,
      `TestConsumeRelationOtherChangeStaysConservative`,
      `TestConsumeRelationBackfillsMissingTableSchema`,
      `TestConsumeRelationNilSchemaStoreStaysNoop` (`mapper_test.go`).
- [x] `LH-FA-SCH-005` real geschlossen: `TestMVPSchemaChangeAddColumn`
      (`slice-030`) läuft ohne Konzeptänderung grün und belegt
      unterscheidbare `schema_version`-Werte vor/nach `ALTER TABLE ADD
      COLUMN`. `make gates` und `make test-integration` dreimal in Folge
      grün. Beleg: `test/integration/integration_test.go` (Assertion
      dreht auf „unterscheidbar"), drei aufeinanderfolgende grüne
      `make test-integration`-Läufe.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update, falls ein öffentlicher Vertrag berührt wird —
      Implementer entscheidet und begründet im Plan-Nachzug. Kein
      öffentlicher Vertrag geändert: `CDC_TABLES`-Format, Port-Signaturen
      (`SchemaStorePort` selbst) und `ADR-0015` bleiben unverändert; die
      neue `Config.SchemaStore`/`Consume(ctx, …)`-Signatur ist interne
      Adapter-/Assembler-Verdrahtung ohne öffentlichen Vertrag.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Entfällt: Repo ist Greenfield (`harness/conventions.md` §Modus-Deklaration), `../reconciliation.md` existiert nicht.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/adapters/driving/replication/decode/decode.go` | update | `Column` um Typ-OID erweitert |
| `internal/adapters/driving/replication/mapper/mapper.go` | update | `Assembler.Consume` behandelt `*decode.Relation` real, nutzt `SchemaStorePort` |
| `internal/bootstrap/wiring.go` | update | `SchemaStorePort`-Adapter (aus `slice-031`) real in die Assembler-Konstruktion verdrahten |
| `internal/adapters/driving/replication/decode/decode_test.go` | update | Plan-Nachzug: Unit-Test für die neue Spalten-Typ-OID (`TestDecodeRelationColumnTypeOID`), Aufrufstelle an `Consume(ctx, …)`/`NewAssembler(…, nil)` angepasst |
| `internal/adapters/driving/replication/mapper/mapper_test.go` | update | Plan-Nachzug: `fakeSchemaStore`-Stub und Klassifikationstests (unverändert / kompatible Erweiterung / andere Änderung / Backfill / `nil`-Store) |
| `internal/adapters/driving/replication/receive/receive.go` | update | Plan-Nachzug: `Config.SchemaStore`-Feld, `Consume(ctx, …)` statt `Consume(…)` — der `SchemaStorePort`-Aufruf im Assembler braucht einen Kontext, den `Stream.process` bereits trägt |
| `test/integration/integration_test.go` | update | Plan-Nachzug: `TestMVPSchemaChangeAddColumn`s Funktionskommentar und letzte Prüfung an den real geschlossenen `LH-FA-SCH-005`-Boundary-Fall angepasst (Erwartung dreht von „dieselbe Version" auf „unterscheidbare Versionen") |
| `internal/bootstrap/walretention_endtoend_test.go` | update | Plan-Nachzug: `cdc.table_schema`-DDL inline ergänzt — `ApplySchema`/`schema.sql` trägt sie bewusst nicht (`slice-031`-Abgrenzung, `schemastore_test.go`), aber dieser Lauf verdrahtet `bootstrap.Run` real und durchläuft damit erstmals den `SchemaStorePort` im laufenden Erfassungspfad |
| `tools/schema/nacharbeit-roles.sql` | update | Plan-Nachzug: `GRANT INSERT` auf `cdc.schema_version` und `GRANT SELECT, INSERT` auf `cdc.table_schema` an `cdc_capture` — beide fehlten bislang jeder Rolle (vor diesem Slice griff nichts im laufenden Erfassungspfad auf `cdc.table_schema` zu); vor dem ersten `make test-integration`-Lauf ergänzt, `RegisterVersion`/`TableSchema` laufen darüber real über `CDC_CAPTURE_DSN` (dreimal grün bestätigt) |

### Plan-Nachzug (nach Implementierung)

- **`Consume`-Signatur:** `Assembler.Consume(ctx context.Context, event
  decode.Event)` statt ohne `ctx` — der `SchemaStorePort`-Aufruf in
  `observeRelation` braucht einen Kontext; `Stream.process` trägt ihn
  bereits (`receive.go`), die Änderung reicht ihn nur durch. Alle
  Aufrufstellen (Tests, `receive.go`) sind entsprechend angepasst.
- **`NewAssembler`-Signatur:** drittes Argument `schemaStore
  outbound.SchemaStorePort`, `nil` zulässig — dann bleibt die
  Relation-Behandlung wirkungslos (Tests, die diesen Pfad nicht prüfen,
  und jeder bestehende Aufrufer vor diesem Slice). Die Aktivierungs-Map
  wird beim Anlegen kopiert (`copiedTables`), damit `observeRelation`
  eine aktualisierte `TableBinding` nicht in die Aufrufer-Map zurückschreibt.
- **Vergleichslogik (`classifyRelationColumns`):** drei Fälle über einen
  Namens→OID-Abgleich der bekannten gegen die eingehenden Spalten — fehlt
  eine bekannte Spalte in der eingehenden Relation oder trägt sie dort
  eine andere Typ-OID, gilt das *immer* als `relationOther` (konservativ),
  unabhängig von der Spaltenzahl. Sind alle bekannten Spalten unverändert
  vorhanden, entscheidet die Spaltenzahl zwischen `relationUnchanged`
  (gleich) und `relationCompatibleExtension` (eingehend größer).
- **Neue-Versions-ID-Format:** `<SourceTableID>-v<Versionsnummer>`
  (`nextSchemaVersionID`, z. B. `tbl-mvp-schema-v2`) — eindeutig je
  Tabelle, weil die Versionsnummer je Tabelle monoton steigt
  (`SchemaStorePort.CurrentVersion` liest die höchste). Kollisionsfrei
  über Tabellen hinweg, weil `SourceTableID` je Tabelle eindeutig ist.
- **Backfill statt Erstregistrierung (Risiko 2, §6):** Trägt die aktuell
  bekannte Version noch keine `TableSchema` (`outbound.ErrSchemaVersionUnknown`
  — der Regelfall für jede Tabelle direkt nach der statischen
  Erstaktivierung, die nur die Versions-Zeile schreibt, keine
  Spaltenform), trägt die erste real eintreffende Relation-Nachricht die
  Spaltenform zur **bestehenden** Version nach (`RegisterVersion` mit
  unveränderter `SchemaVersionID`/Versionsnummer) — keine neue Version,
  keine Bindungs-Änderung. Das ist eine bewusste Entscheidung, keine
  unerwartete Nebenwirkung: `RegisterVersion`s Idempotenz (`slice-031`,
  „Backfill" bereits im Namen) trägt genau diesen Fall.
- **`internal/bootstrap/wiring.go`:** eigener Pool für den SchemaStore
  (`postgresstorage.NewSchemaStore(ctx, cfg.CaptureDSN, …)`), gebunden an
  dieselbe Rolle wie Store und Stream (`cdc_capture`, `ADR-0047`) — ein
  eigener `defer schemaStore.Close()`, kein geteilter Pool mit dem
  ChangeStore-Adapter.
- **Rollen-Lücke real gefunden, nicht nur vermutet:** Vor dem ersten
  `make test-integration`-Lauf zeigte die Lektüre von
  `tools/schema/nacharbeit-roles.sql`, dass `cdc.table_schema` keiner
  Rolle gegrantet war und `cdc_capture` für `cdc.schema_version` nur
  `SELECT` trug — beide Lücken wurden vor dem ersten Lauf geschlossen
  (siehe Tabellenzeile oben).
- **Zweite reale Lücke, diesmal rot gesehen:** `make test-replication`
  (Nicht-Gate-Sensor, engster nützlicher Sensor für die berührte Schicht)
  lief zunächst rot — `TestWALRetentionThresholdEndToEnd`
  (`internal/bootstrap`) verdrahtet `bootstrap.Run` real gegen eine per
  `postgresstorage.ApplySchema` (eingebettete `schema.sql`) aufgesetzte
  Instanz; diese DDL trägt `cdc.table_schema` bewusst nicht (`slice-031`s
  Abgrenzung, dieselbe wie bei den Consumer-State-Tabellen) — bis zu
  diesem Slice traf das nie zu, weil nichts im laufenden Erfassungspfad
  darauf zugriff. Behoben durch eine inline `CREATE TABLE
  cdc.table_schema` im Test-Setup (mirrort den bestehenden Stil dieses
  Tests für seine eigenen Test-Tabellen), **nicht** durch eine Änderung an
  `schema.sql` selbst — die Abgrenzung aus `slice-031` bleibt bestehen,
  nur dieser eine reale Erfassungspfad-Lauf braucht die Tabelle jetzt
  zusätzlich.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-031` liegt in `done/`
(`SchemaStorePort`/`TableSchema` real vorhanden), `Verantwortlich:`
gesetzt, WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich,
  dass Decoder-Erweiterung und Mapper-Verdrahtung zusammen mehr als drei
  Liefer-Punkte oder mehr als zwei Schichten in einer Review-Sitzung
  nicht mehr prüfbar machen, gehört das zurück zur Zerlegung (z. B.
  Decoder-Erweiterung als eigener Folge-Slice).
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

- Die Grenze zwischen „kompatible Erweiterung" und „nicht sicher
  interpretierbare Änderung" könnte in der Praxis unschärfer sein als
  eine reine Obermengen-Prüfung auf Spaltennamen — z. B. eine Spalte, die
  entfernt und mit gleichem Namen, anderem Typ wieder hinzugefügt wird,
  in einer einzigen `pgoutput`-Relation-Message. **Ausgang:** <bei
  Closure einzutragen>
- Bestehende, bereits aktivierte Tabellen ohne je registrierte
  `TableSchema` (Version 1 aus der statischen Erstaktivierung, siehe
  `slice-031`s Backfill-Fähigkeit) könnten beim ersten real eintreffenden
  `*decode.Relation`-Ereignis eine unerwartete Erstregistrierung
  auslösen, wenn der Assembler keine vorhandene `TableSchema` findet.
  **Ausgang:** <bei Closure einzutragen>

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

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <Guide oder Sensor> <geschärft/ergänzt>: <was genau>
  — liegt in `<AGENTS.md §X | Makefile:<target> | .harness/skills/…>`.
  Auslöser: `BEO-<NNN>` (<slice-NNN>, <slice-MMM>, <slice-KKK> — 3×).
  *(Wurde mit diesem Slice nichts verkörpert — der Normalfall —, entfällt die
  Teil-Zeile `— liegt in …` ersatzlos. Der Eintrag ist dann gezählt, nicht
  verkörpert.)*
- **Beobachtungs-Register (`../observations/`):** <`BEO-<KUERZEL>/<slug>/` neu angelegt, Beleg `evidence/slice-NNN.md` | `evidence/slice-NNN.md` in `BEO-<KUERZEL>/<slug>/` ergaenzt — Zaehler steht damit bei <N>x | keine Beobachtung angefallen>
- **Folge-Slices:** <slice-NNN (<Titel>) — ist eine Datei in `open/`>
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** <nur im Repo ohne Wellen-Betrieb — Anker · Folge-Slice · Register, Ergebnis>

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
offen — dieser Slice liefert den zweiten Baustein der Auflösung, ohne sie
selbst zu vollenden — `slice-033` fehlt noch für `LH-FA-SCH-004`). Keiner
der übrigen Treffer erreicht mit diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
