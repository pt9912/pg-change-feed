# Review-Report: slice-backfill-change-origin — 2026-09-24

**Review-Art:** Code — der Diff führt das Feld `origin` von der Domäne bis zu den
Lesewegen (`cdc.change`, View `cdc.changes`, `GET /changes`) und setzt in der
Fixrunde den automatischen View-Vorlauf des Schema-Rollouts um; geprüft gegen
Plan, ADRs und `AGENTS.md` Hard Rules (Modul 10 §Drei Review-Arten). Kein
DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Slice `slice-backfill-change-origin`, Diff-Range
`09386619..HEAD` (`c1df2ac6`) ohne `docs/plan/adr` und `docs/reviews`. Slice-Commits
`3b7f828a`, `f8ae1fc6`, `f4ba82ab` (Lifecycle, Verantwortlich), `98dab17f`
(Domäne), `c40e2ead` (Store, Schema), `e95937cf` (HTTP, Handbuch), `2608df37`,
`5f126971`, `988e8a40` (Plan, `plan.yaml`); Fixrunden-Commits zu
[`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md)
`e01d4ba6` (rolloutguard), `14f971b7` (Makefile-Vorlauf, Guard-Test),
`fa954778` (Handbuch, README, Demo-Skript), `4744c4a5`, `c1df2ac6` (Plan). Nicht
Gegenstand (Constraint bzw. Eingang): die ADR-Datei, das Architect-Verdikt
(`4ba4ef27`, `3acd2c8d`), der Planner-Träger-Commit `ae97247a`.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(seither um weitere HIGH-Klassen ergänzt, u. a. Zahl-im-Träger, Beleg-Satz,
Zusage-ohne-Eingabeseite, Kommentar-Chronik, Handbuch-Zug).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-24.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-NNN` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** (`v6.9.0` ·
> `regelwerk/<datei>.md` §<Abschnitt>). Ein `pfad`-Feld auf den **geprüften
> Gegenstand** zitiert den Stand des Laufs und darf ihn festhalten.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-backfill-change-origin` (§1 Ziel, §3 Plan, Suchlauf-Felder
  und Rollout-Messung, §6 Risiken) und Welle `welle-backfill-bestand`
- [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage
  2 und 8 (Feld `origin`, Reichweite in den Lesewegen),
  [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) (Schema-Rollout,
  Guard),
  [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md) (Accepted,
  Constraint der Fixrunde),
  [`ADR-0064`](../plan/adr/0064-lh-qa-ops-005-testansatz-korrektur.md)
  (Upgrade-Testansatz, Re-Evaluierungs-Trigger 1),
  [`ADR-0058`](../plan/adr/0058-testansatz-fuenf-luecken.md) Entscheidung 2,
  [`ADR-0081`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md)
- [`LH-FA-CAP-009`](../../spec/lastenheft.md),
  [`LH-FA-DAT-006`](../../spec/lastenheft.md),
  [`LH-FA-SST-002`](../../spec/lastenheft.md),
  [`LH-FA-SST-006`](../../spec/lastenheft.md),
  [`LH-QA-OPS-005`](../../spec/lastenheft.md);
  [`SPEC-002`](../../spec/pflichtenheft.md),
  [`SPEC-022`](../../spec/pflichtenheft.md)
- `AGENTS.md` (Hard Rules §3.1, §3.3, §3.6, §3.7, §3.9, §3.11, §3.12, §3.13),
  `harness/conventions.md` (MR-000/MR-001)
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-backfill-row-image-gemeinsam.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem
Implementer-Bericht übernommen):

- **Gates am Stand `c1df2ac6`, Exit-Codes ungepiped gesichert:** `make test`
  (Race-Detector) Exit 0; `make a-check` Exit 0 („gesamt: 0 Befund(e)");
  `make coverage-gate` Exit 0 („Coverage 82.90% erfüllt Schwelle 80%");
  `make test-store` (Basislauf) Exit 0 („DB-Adapter-Coverage: 77.13%", Schwelle
  70 %). Die Unit-Tests von `tools/schema/rolloutguard` (18 Tests, gezählt in
  `guard_test.go`) laufen in `make test`.
- **Guard-Test real, ungekürzt, einmal:** `bash
  tools/harness/run-schema-rollout-guard-test.sh` → Exit 0, Laufzeit 1 min 53 s;
  Lauf 5 druckt „Tag v0.1.2: Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, mit
  Vorlauf), Exit 0 (Arbeitsbaum, zweiter Lauf); Zeile alttag-ch über cdc.changes
  lesbar". Nach dem Lauf: kein Container, kein Netz, kein Temp-Verzeichnis
  `schema-rollout-alt-tag.*` übrig (`docker ps -a`, `docker network ls`, `ls`);
  `git status` zeigt `tools/schema/plan.yaml` und `tools/schema/down.sql` geändert
  (siehe F-6) — von mir per `git checkout` zurückgenommen.
- **Zahlen der Suchlauf-Felder** (Plan §3) mit `git grep` an beiden Ständen
  nachgemessen: Diff-Feld `f4ba82ab`/`e95937cf` (6 Treffer `committed_at`, 33
  `Felder`, 19 `zehn Felder`, 3 `SELECT *`, 8 „dieselben Felder", 0 strikte
  Dekoder, 14 Schema-Dateien) und Fixrunden-Feld `ae97247a`/`fa954778` (33→42
  Guard-Treffer, 9→11 Lauf-Zahl, 2→0 „nicht idempotent") stimmen. Der
  Suchraum selbst ist unvollständig: siehe F-4.
- **`readChangeResponse`:** 12 → 13 JSON-Felder (`git show`, Parent gegen Kopf), 
  `origin` letztes, Reihenfolge der zwölf übrigen unverändert.
- **Live-Wege:** `git grep` über `internal` nach `Origin` außerhalb von Domäne,
  Mapper, Store, `readchanges.go`: keine Verwendung in gRPC, SSE, NATS — die
  Live-Wege bleiben bei zehn Feldern (`SPEC-020`/`SPEC-021`/`SPEC-024`);
  `make generated-sync` ist nicht berührt (`gen/`, `proto/` ohne Diff).
- **`plan.yaml`/`down.sql`:** die Fixrunde lässt beide unverändert
  (`git diff 5f126971 HEAD --stat` leer); die Ziel-Zeile in `plan.yaml` trägt die
  Compose-Bezeichnung (`cdc-test-postgres`), `origin` und
  `COALESCE(c.origin, 'wal')` stehen im Report, die Operations-IDs von `plan.yaml`
  und `down.sql` stimmen überein.
- **Kommentare (§3.7):** die hinzugefügten Zeilen in `*.go`, `Makefile`, `*.sh`,
  `*.yaml` per Textsuche auf `slice-`/`welle-`/„früher"/„bisher"/„wäre"/„würde"/
  „nicht mehr"/„jetzt" durchsucht: keine Treffer in Produktionscode-Kommentaren.
  Kein `//nolint`, keine host-lokalen absoluten Pfade im Diff. Skripte rufen
  ausschließlich `docker`, `make`, `git`, `tar`, `mktemp` (wie bestehende Runner);
  keine Host-Toolchain (§3.1).
- **Commit-Struktur:** `git show --stat -M` über die Lifecycle-Commits `3b7f828a`
  und `f4ba82ab`: reine Renames ohne Zeilenänderung, Inhalt (`f8ae1fc6`) getrennt —
  §3.3 erfüllt. Alle Betreffs tragen `LH-*`/`ADR-*`, keine `SPEC-`/`ARC-`-Kennung.
- **Spec:** [`SPEC-002`](../../spec/pflichtenheft.md) und
  [`SPEC-022`](../../spec/pflichtenheft.md) tragen `origin` (Zeilen 338, 345, 617
  `spec/pflichtenheft.md`; letztes Feld, fehlender Wert liest `wal`, kein Filter) — der Diff setzt das um, ändert die Spec nicht.

### Mutationen (Eingabeseite, selbst ausgeführt)

Datei nach jeder Mutation per `git checkout` zurückgenommen; Endstand `git diff`
leer, `git status` sauber.

| # | Mutation | Ort | Ergebnis |
|---|---|---|---|
| M1 | Zeile erhält immer `wal` statt der gesetzten Herkunft | `mapper.go` (`NewChangeRows`) | rot: `TestChangeRowRoundTripsOrigin` |
| M2 | `ToChange` setzt `change.Origin` nicht | `mapper.go` | rot: `TestChangeRowRoundTripsOrigin`, `TestReadChangesCarriesOrigin/backfill` |
| M3 | `OrDefault()` in der HTTP-Antwort entfernt | `readchanges.go` | rot: `TestReadChangesTraegtOriginAlsLetztesFeld/fehlender_Wert_liest_als_wal` |
| M4 | leere Zeichenkette liest nicht mehr als `wal` | `change.go` (`NewChangeOrigin`) | rot: `TestNewChangeOriginClosedSet`, `TestToChangeCarriesSchemaAndTable` |
| M5 | View: `COALESCE(c.origin, 'wal')` → `c.origin` | `schema.yaml` | rot (`make test-store`): `TestChangesViewCarriesOriginLikeReadChanges` |
| M6 | `SelectChanges`: `COALESCE` entfernt | `queries.go` | rot (`make test-store`): `TestPersistAndReadCarryChangeOrigin`, `TestChangesViewCarriesOriginLikeReadChanges` |
| G1 | Diagnose-Prüfung `!signatureDiagnosed[id]` entfernt | `guard.go` | rot: drei `TestDecideRefuses…Diagnostic…`/`…WithoutSignatureDiagnostic` |
| G2 | View-Namen-Regex (`viewIdentifier`) entfernt | `guard.go` | rot: `TestDecideRefusesViewNameThatIsNotAnIdentifier` (u. a. `changes; DROP TABLE cdc.change`) |
| G3 | JSON-Tag `operationId` umbenannt | `report.go` | rot: `TestParseReportDecodesRealViewSignatureReport` |
| G4 | Prüfung `ReplaceView`/`VIEW` entfernt | `guard.go` | rot: `TestDecideRefusesSignatureDiagnosticOnNonViewOperation` |
| G5 | `knownForeignObjects`-Prüfung entfernt | `guard.go` | rot (Unit): `TestDecideRefusesUnknownDestructiveBlocker`, `TestDecideViewSignatureWithUnknownBlockerRefusesEverything`; **Guard-Test Exit 0 — Mutation überlebt**, siehe F-2 |
| G6 | `MANUAL_ACTION_REQUIRED` ohne Operationen erlaubt | `guard.go` | rot: `TestDecideRefusesManualActionWithoutOperations` |

Die Zusagen von Domäne, Mapper, Store, View und HTTP sind an ihrer Eingabeseite
gebunden. Die Guard-Entscheidung ist über die Report-Felder (JSON-Tags gegen einen
real gemessenen Report) und über die Klassen-Merkmale an ihrer Eingabeseite gebunden.

**Sonde ohne Mutation (temporäre Testdatei, entfernt):** ein Blocker mit Grund
`DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION` und **ohne** `operationIds` liefert
`allowDestructive: true` — allein wie neben einem View-Signatur-Blocker (F-1).

---

## Findings

<!-- Kein Fließtext, kein Lösungsvorschlag im Befund. -->

### F-1 — Guard: Blocker `DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION` ohne Operationen gilt als „bekannt" (Vakuum)

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md)
  Entscheidung 2 („jeder Blocker … eine der bekannten Fremdobjekt-Operationen"),
  `AGENTS.md` §3.6 (Gate-Lockerung nur über ADR)
- `pfad`: `tools/schema/rolloutguard/guard.go:96-106` (Zweig
  `destructiveConfirmationReason`, `d.allowDestructive = true` vor der Schleife
  über `b.OperationIDs`)
- `befund`: Ein Blocker dieser Klasse mit leerer `operationIds`-Liste setzt
  `allowDestructive`, ohne dass eine bekannte Fremdobjekt-Operation geprüft wurde
  (Sonde: `{Reason: DESTRUCTIVE…}` allein → `allowDestructive: true`; neben einer
  View-Signatur-Klasse → `allowDestructive: true`, `dropViews: [changes]`). Der
  Zweig `manualActionReason` derselben Funktion lehnt den leeren Fall ausdrücklich
  ab (`guard.go:107-109`, eigener Test); das Verhalten des Zweigs
  `destructiveConfirmationReason` ist gegenüber dem Parent-Stand unverändert, der
  ADR-Wortlaut (jeder Blocker ist eine bekannte Operation) ist damit für diesen
  Fall nicht umgesetzt.
- `verifizierbar`: ja — Unit-Test mit einem Blocker `DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION`
  ohne `operationIds` (`make test`); ein Test dafür existiert nicht
- `klasse`: Sicherheitsnaht mit Vakuum-Zweig (Guard bindet „alle" ohne leere Menge auszuschließen)

### F-2 — Guard-Test Lauf 6 belegt die Bekannt-Liste des Guards nicht (Mutation überlebt)

- `kategorie`: MEDIUM
- `quelle`: Skill-HIGH-Klasse „Beleg trägt seinen Satz nicht" (Assertion hält auch aus
  einem anderen Pfad), hier zurückgestuft, weil die Eigenschaft anderweitig
  gebunden ist (G5: Unit-Tests rot)
- `pfad`: `tools/harness/run-schema-rollout-guard-test.sh:37-39` (Kopf: „muss
  weiterhin mit Exit 8 abbrechen (der Beleg, dass die Wache nicht pauschal
  durchlässt)"), `:245-249` (Assertion `run6_exit -eq 0`), `harness/README.md:151`
  (Zeile `make schema-rollout`, „Negativ-Abbruch")
- `befund`: Selbst gemessen: mit entfernter `knownForeignObjects`-Prüfung meldet der
  Guard „alle Blocker sind bekannte Fremdobjekte" und das Target läuft mit
  `--allow-destructive`, `make` endet trotzdem mit Fehler 8 (d-migrate bricht den
  `DropColumn` selbst ab); der Guard-Test endet Exit 0. Die Assertion prüft nur
  „nicht 0" (das Skript druckt „Exit 2", der Kopf sagt „Exit 8"); der Satz „Beleg,
  dass die Wache nicht pauschal durchlässt" steht am Träger, die Grenze steht nur
  im Slice-Plan (§3, „Grenze"), nicht im Skript oder in der Sensors-Zeile.
- `verifizierbar`: ja — Mutation G5 mit anschließendem
  `bash tools/harness/run-schema-rollout-guard-test.sh` (Exit 0 trotz Mutation)
- `klasse`: Beleg trägt seinen Satz nicht (Negativlauf mit fremdem Abbruchpfad)

### F-3 — Makefile-Kommentarblock widerspricht sich zu `--allow-destructive`

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Kommentar beschreibt, was da ist)
- `pfad`: `Makefile:293-295` gegen `Makefile:320-325`
- `befund`: Zeile 293-295 sagt, `rolloutguard` entscheide, ob „ausschließlich die
  sechs bekannten Objekte blockieren", „nur dann" laufe `--execute` mit
  `--allow-destructive`; der Block ab Zeile 320 lässt `--allow-destructive` auch
  neben der Klasse „View-Signatur" zu (der real gemessene Upgrade-Fall,
  `TestDecideViewSignatureWithKnownForeignObjects`). Der frühere Absatz wurde beim
  Ergänzen nicht mitgezogen.
- `verifizierbar`: nein
- `klasse`: Kommentarblock ohne Nachzug bei erweiterter Zusage

### F-4 — Suchlauf-Feld (Plan §3): Zeile „Anzahl-Formulierungen" unvollständig — Kotlin-KDoc nennt „twelve fields"

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.13 (Suchlauf trägt Gefundenes **und** Nichtgefundenes),
  §3.12 Instanz A
- `pfad`: Slice-Plan `slice-backfill-change-origin` §3, Suchlauf-Zeile
  „Anzahl-Formulierungen“ (Stand `c1df2ac6`);
  `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/sse/model/Change.kt:19-20`,
  `sdks/csharp/PgChangeFeed.Client/Sse/Models/Change.cs:17`,
  `sdks/csharp/PgChangeFeed.Client/Nats/Models/Change.cs:21`
- `befund`: Die Zeile nennt außerhalb des Suchraums nur den Python-Kommentar; die
  Kotlin-KDoc („`SPEC-022`, twelve fields", am Parent wahr, jetzt 13) und die zwei
  C#-Kommentare („eleven fields", am Parent bereits ungleich zu den gezählten 12
  Feldern) fehlen — die Suche lief über `docs/user spec` und `*.go *.proto *.py`,
  nicht über `*.kt`/`*.cs`. Die Stellen gehören zum Folge-Slice
  `slice-backfill-sdk-origin`, dessen Plan keine Zahl nennt
  (`grep -in 'twelve\|eleven\|zwölf' docs/plan/planning/open/slice-backfill-sdk-origin.md` leer).
- `verifizierbar`: ja — `git grep -n -i -E 'twelve|eleven' HEAD -- '*.kt' '*.cs'`
- `klasse`: Zahl im Träger driftet gegen die Messung (Suchraum zu eng)

### F-5 — `DROP VIEW` verwirft Betreiber-Privilegien auf `cdc.changes` außer `cdc_reader`

- `kategorie`: LOW
- `quelle`: [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md)
  Entscheidung 3 (Rechte setzt `nacharbeit-roles.sql`), `AGENTS.md` §3.12 Instanz B
- `pfad`: `docs/user/benutzerhandbuch.md:621-623` („Die Rechte der Rolle
  `cdc_reader` setzt derselbe Lauf wieder"), `Makefile:311-315`
- `befund`: Das Vorlauf-Statement entfernt die View samt ihrer ACL;
  `nacharbeit-roles.sql:105` grantet nur an `cdc_reader`. Ein vom Betreiber
  gesetztes `GRANT SELECT ON cdc.changes` an eine andere Rolle geht im Lauf mit
  Signaturänderung verloren; Handbuch und Makefile-Kommentar nennen nur `cdc_reader`,
  nicht den Verlust anderer Privilegien.
- `verifizierbar`: nein — kein Gate; Lauf 4 des Guard-Tests prüft nur `cdc_reader`
- `klasse`: Zusage benennt nur die getestete Rolle

### F-6 — Guard-Test hinterlässt `tools/schema/plan.yaml` und `down.sql` im Arbeitsbaum

- `kategorie`: INFO
- `quelle`: Maintainability (Zeile im Slice-Plan §3: „stellen es per `git checkout`
  wieder her")
- `pfad`: `tools/harness/run-schema-rollout-guard-test.sh` (`run_rollout .`,
  `make schema-rollout` ohne `-C`)
- `befund`: Nach dem Lauf sind `plan.yaml` und `down.sql` modifiziert (`git status`,
  gemessen); das Skript stellt sie nicht wieder her, der Kopf benennt das nicht. Das
  Verhalten stammt aus den Läufen 1-3 des Parent-Stands; die neuen Läufe 4 und 5
  schreiben mit.
- `verifizierbar`: ja — `git status --short` nach dem Skript
- `klasse`: Runner verändert committete Erzeugnisse (kein Zähler)

### F-7 — `tools/harness/run-integration-tests.sh:2711` beschreibt den Ist-Stand der Wache nicht mehr

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.13 (Träger nach bewegter Eigenschaft),
  Slice-Plan §3 (bewusst gemeldet: Zeilenanker in `docs/user/e2e-abdeckung.md`)
- `pfad`: `tools/harness/run-integration-tests.sh:2710-2712` („Exit 8 auf vier
  Fremdobjekten")
- `befund`: Der Kommentar nennt vier Fremdobjekte und einen real blockierten zweiten
  Rollout; die Wache besteht seit `schema-rollout-zentrale-idempotenz-wache`, es
  sind sechs Objekte, und `ADR-0114` schließt für das Schema den Trigger 1 von
  `ADR-0064`. Die Zeile ist eine Begründung der Test-Wahl in einem Runner, dessen
  Zeilennummern `docs/user/e2e-abdeckung.md` trägt (Spalte `Datei:Zeile`,
  z. B. Zeilen 314, 583, 658); ein zeilentreues Umschreiben verschöbe sie nicht.
  Die Meldung im Plan-Feld ist eine korrekte Übergabe; der Träger bleibt im Bestand.
- `verifizierbar`: nein
- `klasse`: Träger einer bewegten Eigenschaft nicht nachgezogen (gemeldet)

### F-8 — Vorlauf adressiert das Schema `cdc` fest; der Report trägt nur den View-Namen

- `kategorie`: INFO
- `quelle`: Maintainability; [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md)
  Entscheidung 3 (`DROP VIEW cdc.<name>`)
- `pfad`: `Makefile:358` (`DROP VIEW cdc.$$v`), `tools/schema/rolloutguard/guard.go:119`
  (Pfad genau ein Segment)
- `befund`: `path` der `ReplaceView`-Operation enthält kein Schema (real gemessen im
  Testreport: `["changes"]`); das Ziel-Schema `cdc` gilt über die Rollout-Vorbedingung
  `search_path = cdc` (`tools/schema/apply-rollout.sh`, `examples/bootstrap.sh`). Ein
  Ziel ohne diese Vorbedingung träfe der Vorlauf im falschen Schema oder liefe mit
  dem PostgreSQL-Fehler ab (`ON_ERROR_STOP`, kein stilles Löschen); Makefile-
  Kommentar und Handbuch nennen die Vorbedingung im Vorlauf-Zusammenhang nicht.
- `verifizierbar`: nein
- `klasse`: implizite Umgebungs-Vorbedingung (kein Zähler)

### F-9 — Blocker-Grund `MANUAL_ACTION_REQUIRED` ist enger als der ADR-Wortlaut (bewertet, kein Widerspruch)

- `kategorie`: INFO
- `quelle`: [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md)
  Entscheidung 1 und 2
- `pfad`: `tools/schema/rolloutguard/guard.go:107-122`, Slice-Plan §3 (Fixrunde,
  Zeile `rolloutguard`, „Präzisierung gegenüber dem ADR-Wortlaut")
- `befund`: Entscheidung 1 definiert die Klasse über Operation und Diagnose, nicht
  über den Blocker-Grund; die Umsetzung verlangt zusätzlich `MANUAL_ACTION_REQUIRED`.
  Enger heißt: ein künftig unter anderem Grund gemeldetes `ReplaceView`/
  `VIEW_SIGNATURE_INCOMPATIBLE` bricht mit Exit 8 ab statt den Vorlauf zu fahren —
  die sichere Richtung, kein Widerspruch zu Entscheidung 2 („jeder andere Blocker
  lässt beides ausfallen"). Die Präzisierung steht im Plan und im Kommentar von
  `guard.go` (Konstanten `manualActionReason`/`viewSignatureCode`).
- `verifizierbar`: ja — `TestDecideRefusesOtherDiagnosticCode`,
  `TestDecideManualActionWithoutSignatureDiagnostic`
- `klasse`: Umsetzung enger als Entscheidungstext (zulässig)

### F-10 — Gemeldet, nicht Befund: `.proto`/`.pb.go` nennen weiter „dieselben Felder wie der Domain-Typ"

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.13; Slice-Plan §3 (Suchlauf-Zeile „Kommentare in den
  Live-Wegen")
- `pfad`: `proto/cdc/stream/v1/changestream.proto:13`,
  `gen/cdc/stream/v1/changestream.pb.go:31`
- `befund`: Die Aussage ist mit `Change.Origin` an der gRPC-Nachricht ungenau; eine
  Änderung verlangt `make proto-generate` und berührt `make generated-sync`, die der
  Plan als unberührt festlegt. Die Spec (`SPEC-020`, Zeile 566) sagt „`origin` gehört
  nicht zur Nachricht" und trägt die Grenze. Der Plan meldet die Stelle korrekt.
- `verifizierbar`: ja — `git grep -n 'Feldern wie der Domain-Typ' HEAD -- proto gen`
- `klasse`: Träger einer bewegten Eigenschaft nicht nachgezogen (gemeldet)

## Negativbefunde

- geprüft, ohne Befund: `internal/domain/model`, `internal/domain/errors` —
  geschlossene Menge (`wal` | `backfill`, leere Zeichenkette liest `wal`,
  Konstruktor-Default `wal`), Fehler-Sentinel, Doc-Kommentare im Indikativ; Tests
  binden Eingabe (`""`, `"WAL"`, `" backfill"`, `"snapshot"`).
- geprüft, ohne Befund: `internal/adapters/driven/postgresstorage` (Mapper,
  `queries`, `sqlexec/translate.go`, `store.go`, `schema.sql`, Tests) — explizite
  Spaltenlisten, `origin` letzte Spalte von `InsertChange` und `SelectChanges`,
  einheitlicher `COALESCE`, zweiter Schema-Träger `schema.sql` in derselben Form wie
  `schema.yaml` (nullable, ohne Default, ohne CHECK); Scan-Reihenfolge (14. Spalte)
  gleich der Projektion; Byte-/Verhaltens-Erhalt der übrigen Felder (`make test`,
  `make test-store` grün).
- geprüft, ohne Befund: `tools/schema/schema.yaml` — View `changes` mit `origin` als
  letzter Spalte und `columns:`-Signatur; `tools/schema/plan.yaml`/`down.sql`
  konsistent (Operations-IDs, Fingerprint, Ziel-Zeile).
- geprüft, ohne Befund: `internal/adapters/driving/http/readchanges.go` (+ Test) —
  `origin` letztes Feld, `OrDefault()`, Reihenfolge der zwölf übrigen unverändert;
  `internal/adapters/driving/http/sse.go`,
  `internal/adapters/driven/natsstream/publisher.go` (+ Test) — nur Kommentare, „ohne
  `Origin`" wahr; keine Live-Wege-Änderung.
- geprüft, ohne Befund: `tools/schema/rolloutguard` — Klasse exakt „ReplaceView +
  `MANUAL_ACTION_REQUIRED` + `VIEW_SIGNATURE_INCOMPATIBLE` zur selben `operationId`",
  alles oder nichts; View-Name durch `^[a-z_][a-z0-9_]*$` (Go-`$` ohne
  Multiline-Nachlauf) und genau ein Pfadsegment beschränkt; kein Weg, dass ein
  unbekannter Blocker den neuen Pfad zu `--allow-destructive` führt (außer dem
  Vakuum-Fall F-1, der derselbe wie im Parent-Stand ist); `--allow-destructive` wird
  für die Klasse allein nicht gesetzt (`TestDecideViewSignatureAlone`).
- geprüft, ohne Befund: `Makefile` Target `schema-rollout` — Namen kommen validiert
  aus der Maschinenzeile `drop-view <name>` (kein Shell-/SQL-Injektionsweg, die
  Zeile wird nach der Go-Validierung ausgelesen, `for v in` ohne Metazeichen);
  `DROP VIEW` ohne `CASCADE`, `-v ON_ERROR_STOP=1`, `|| exit 1`; Exit-Codes hängen an
  keiner Pipe (`guard_out=$(…) && guard_ok=1 || guard_ok=0`, nur `grep`/`sed` auf
  gesichertem String); nicht-atomarer Vorlauf im Kommentar und im Handbuch benannt
  (Abbruch nach dem `DROP` → View fehlt bis zum Wiederholungslauf); Richtwert „rund
  7 s" trägt Ursprung (`ADR-0114`, Architect-Messung Szenario 7, ein Lauf) und den
  Zusatz „Messung, nicht garantiert".
- geprüft, ohne Befund: `tools/harness/run-schema-rollout-guard-test.sh` — Lauf 4
  (Vorlauf gemeldet, Soll-Signatur, `SELECT`-Recht, Zeile lesbar, Folgelauf ohne
  Vorlauf; abhängiges Objekt in `public` bleibt bestehen, Wiederholungslauf heilt),
  Lauf 5 (`git archive` des jüngsten `v*`-Tags in ein `mktemp`-Verzeichnis unter
  `TMPDIR`, nie im Repo-Baum; Tag und Exit-Codes gedruckt — Belegform gegen
  `ADR-0114` Entscheidung 7), `trap cleanup EXIT` räumt Container, Netz und
  Temp-Verzeichnis auch beim Fehlschlag; Ergebnis real Exit 0, nichts übrig.
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — `Version:` 1.43 → 1.45 mit
  zwei Historienzeilen, beide benannten Stellen (SQL-Beispiel, Feldliste der
  `GET /changes`-Antwort) gezogen; Abschnitt „Schema aktualisieren" nennt Reihenfolge
  (Rollout vor Container-Tausch), Vorlauf-Meldung, Lesefenster mit Ursprung
  („einzelne Architect-Messung"), „abgeleitet, nicht gemessen" für Feed-Container und
  laufenden Alt-Container; „keine Daten gehen verloren" — die Aussage über den
  Feed-Container habe ich selbst nachgemessen: kein Go-Code außerhalb der Tests liest
  die View `cdc.changes` (`--healthcheck` liest `cdc.heartbeat`).
- geprüft, ohne Befund: `harness/README.md` Zeile `make schema-rollout`,
  `examples/bootstrap.sh` (Kommentar und Meldung, Verhalten unverändert),
  `AGENTS.md` §3.14 — Beschreibung entspricht dem Code (sechs Läufe, Vorlauf,
  `--allow-destructive`-Bedingung).
- geprüft, ohne Befund: Lifecycle/Commit-Struktur (§3.3), Traceability der
  Commit-Betreffs, Suppression-Verbot (§3.2), host-lokale Pfade (§3.11), Docker-only
  (§3.1).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 3 |
| INFO | 5 |

**Finding-Klassen dieses Laufs:** Sicherheitsnaht mit Vakuum-Zweig · Beleg trägt
seinen Satz nicht (Negativlauf mit fremdem Abbruchpfad) · Kommentarblock ohne
Nachzug bei erweiterter Zusage · Zahl im Träger driftet gegen die Messung
(Suchraum zu eng) · Zusage benennt nur die getestete Rolle

## Verdikt

**Merge-blockierend:** ja — F-1 und F-2 (MEDIUM) gehen an den Implementer. Beide
sind klein und ohne Rückkante auf den Plan: F-1 ist eine Zeile im Guard plus ein
Unit-Test (Vakuum-Fall, gleiche Form wie der bestehende Zweig
`manualActionReason`), F-2 verlangt, dass die Grenze des Lauf 6 am Träger (Skript-Kopf,
Sensors-Zeile) steht oder Lauf 6 die Guard-Meldung bindet. F-3 bis F-5 sind LOW
und können in derselben Fixrunde mitgehen. F-6 bis F-10 sind INFO ohne erwartete
Aktion; F-7 und F-10 sind im Plan bereits als gemeldete Träger geführt.

Kein HIGH: die Sicherheitsnaht des Guards hält an allen selbst gesetzten
Eingabeseiten-Mutationen (G1-G4, G6 rot); die einzige Lücke der Entscheidungslogik
(F-1) ist gegenüber dem Parent-Stand unverändert und bei realen d-migrate-Reports
nicht beobachtet (Blocker tragen `operationIds`: der Guard nennt im Lauf 6 dieses
Reviews die `operationId` des realen `DropColumn`-Blockers). Das
DoD-Kästchen „Review durchgeführt" bleibt offen; es wird nach der Fixrunde regulär
nachgezogen (`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug, Fixrunde nötig).

**Übergabe:** Findings gehen an den Implementer; die **Finding-Klassen** gehen
zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Dieser Report ist ein
**Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses Modell, dieses Verdikt).
Er ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11).
