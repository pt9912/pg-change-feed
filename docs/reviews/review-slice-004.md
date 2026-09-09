# Review-Report: slice-004 Implementer-Diff — 2026-09-09

**Review-Art:** Diff-Review (Implementer-Range `c9a068f..193716e`, 3 Commits) —
*wogegen*: Slice-Plan §1/§3 (Plan-Treue, Bewertung der drei gemeldeten
Entscheidungen + Plan-Defekt-Meldung), ADR-Bezüge ([`ADR-0011`](../plan/adr/README.md),
0009, 0010/0032/0038-Fortgeltung, 0023, 0029-Regel 1/3/6/7, 0042, 0005), Hard
Rules (`AGENTS.md` §3.1 Docker-only, §3.7 Kommentar-Klassen, §5
Dokumentations-Regeln), Traceability (LH-*/ADR-* je Commit, keine
Struktur-IDs, keine superseded-Referenzen als tragende Anker), neue
Angriffsfläche (erster realer Treiber-Code: SQL-Übersetzung, Fehler-/
Credential-Grenze, Testcontainer-Hygiene, Digest-Pins), Implementer-Risiken
(a)–(c). Keine DoD-Prüfung — das ist der Verifier (Modul 11).

**Gegenstand:** `cefa18c` (CDC-Schema-DDL) · `37405a6`
(PostgresChangeStoreAdapter + Port-Erweiterung) · `193716e` (Adapter-Tests
gegen reale PostgreSQL, Targets `make test`/`make test-store`, Plan-§3-Nachzug)
— Basis `c9a068f` (slice-004 `next` → `in-progress`).

**Skill:** `.harness/skills/reviewer.md` (Stand 2026-09-09, geschärft: vier
repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen) · Gerüst:
`docs/reviews/review-report.template.md` (Form wie `review-slice-003.md`).

**Modell:** Claude Code (glm-5.3-flash) · **Datum:** 2026-09-09

**Eingangs-Kontext:**

- Diff `git diff c9a068f..193716e` (17 Dateien, +1318/−8: `postgresstorage/`
  inkl. `queries/`+`mapper/`, Port-Erweiterung, Tests, DDL, `go.mod`/`go.sum`,
  Makefile, `tools/harness/run-store-tests.sh`, Dockerfile-Kommentar,
  harness-README-Zeilen)
- `docs/plan/planning/in-progress/slice-004-postgres-store-adapter.md` (§1–§8,
  am Range-Start `c9a068f` und am Range-Head)
- `docs/plan/adr/README.md` · [`ADR-0005`](../plan/adr/README.md) · [`ADR-0009`](../plan/adr/README.md) · [`ADR-0010`](../plan/adr/README.md) ·
  [`ADR-0011`](../plan/adr/README.md) · [`ADR-0023`](../plan/adr/README.md) · [`ADR-0029`](../plan/adr/README.md) · [`ADR-0030`](../plan/adr/README.md) · [`ADR-0032`](../plan/adr/README.md) · [`ADR-0038`](../plan/adr/README.md) · [`ADR-0039`](../plan/adr/README.md) ·
  [`ADR-0042`](../plan/adr/README.md) — Status: Superseded sind [`ADR-0038`](../plan/adr/README.md) (→ 0039) und
  [`ADR-0039`](../plan/adr/README.md) (→ 0042, Rest-Fortgeltung ausdrücklich)
- `spec/lastenheft.md` ([`LH-FA-REA-001`](../../spec/lastenheft.md)…006, [`LH-FA-DAT-002`](../../spec/lastenheft.md)/004, [`LH-FA-CAP-006`](../../spec/lastenheft.md)/008,
  [`LH-FA-RET-001`](../../spec/lastenheft.md), [`LH-FA-SCH-005`](../../spec/lastenheft.md), [`LH-QA-REL-001`](../../spec/lastenheft.md)/002/003),
  `spec/pflichtenheft.md` ([`SPEC-001`](../../spec/pflichtenheft.md)/002/003/004/008, [`LH-FA-REA-004.a`](../../spec/pflichtenheft.md),
  [`LH-FA-CAP-006.a`](../../spec/pflichtenheft.md), §2 Tabellen-Menge inkl. `cdc.consumer`,
  `cdc.consumer_position`, `cdc.capture_state`, Zeile „Keine Credentials in
  Logs"), `spec/architecture.md` ([`ARC-006`](../../spec/architecture.md)/009)
- Commit-Messagen der Range; `AGENTS.md` §3/§5; `harness/conventions.md`
  (MR-000); `harness/README.md` §Sensors (`image`-Zeile: Beleg gilt am HEAD);
  Beobachtungs-Register (`BEO-PGC/a-check-null-abdeckung`)
- Stand-alone-Prüfung je Commit-Menge: `go build`/`go vet` am Range-Head im
  gepinnten Toolchain-Container (`golang:1.27-alpine@sha256:cf6fca…`,
  docker-only, `AGENTS.md` §3.1) — grün; treiberfreie `go test ./...` — grün
  (Store-Tests skippen ohne DSN); **`make test-store` am Range-Head — grün**
  (8 Adapter-Tests gegen die gepinnte `postgres:18-alpine`-Instanz); drei
  Gates am Range-Head (`baseline-verify` OK, `d-check` 89 Dateien/0 Befunde,
  `a-check` 0 Befunde); zwei Mutations-Proben gegen die Store-Suite
  (Worktree außerhalb des Arbeitsbaums, siehe unten); DSN-Error-Probe gegen
  den gepinnten Treiber (Credential-Redaktion, siehe Negativbefunde)

---

## Findings

### F-1 — Plan-Nachzug fehlt für zwei der vier gemeldeten Entscheidungen (F-2-Klasse, 4. Auftreten — laufende Sequenz)

- `kategorie`: MEDIUM
- `quelle`: Slice-Plan §1/§3/§7 · Modul 5 („Wer später mitnimmt …, hat den
  Plan **geändert**, nicht nur ergänzt"; „Was nicht ausdrücklich
  ausgeschlossen ist, wandert im Zweifel hinein") · laufende Konflikt-Sequenz
  (Modul 8 — drittes Auftreten in review-slice-003 F-2)
- `pfad`: `docs/plan/planning/in-progress/slice-004-postgres-store-adapter.md`
  (§1-Ausschlüsse unverändert, §3 um vier Zeilen ergänzt — in `193716e`, mit
  dem Code, nicht vor ihm; §7 ungefüllt — Closure steht aus) gegen
  `internal/adapters/driven/postgresstorage/schema.sql:7-17` (DDL-Umfang) und
  Commit-Messages der Range
- `befund`: Vier gemeldete Punkte, zwei davon ohne Plan-Träger:
  **(b) „Position-Verwaltung" enger gelesen** — *Plan-Ergänzung, nicht
  Plan-Änderung*: die enge Lesung (Commit-Positionen je Transaktion als
  Ordnungs-/Bereichs-Größe des Lesens; Source-ACK-/Consumer-Positionen bei
  slice-005/006) deckt sich mit den bestehenden §1-Ausschlüssen (Replication-
  Stream → slice-005; Consumer-Use-Cases → anderer Vorgang); kein §1-Ausschluss
  und kein DoD-Punkt ist berührt. Sie ist aber eine Interpretation der
  Ziel-Zeile, die nur im Handoff lebt — Nachzug-Zeile (§7 oder §1-Fußnote) vor
  der Closure fällig. **(d) DDL-Umfang** — *Plan-Änderung, nicht nur
  Ergänzung*: §3 nannte „Schema-DDL (`SPEC-001`/`SPEC-002`)"; `SPEC-001`
  führt `cdc.consumer`, `cdc.consumer_position` und `cdc.capture_state` als
  vorgesehene Tabellen, die DDL trägt sie nicht, und §1 schließt genau diese
  Tabellen **nicht** aus (§1 schließt Consumer-*Use-Cases* und Retention-*
  Adapter* aus, nicht ihre Speicher-Seite). Die Schmälerung ist sachlich
  tragfähig — ihre Ports (`ConsumerStatePort`, `ReplicationAckPort`) liegen
  außerhalb dieses Store-Adapters — aber sie ist eine Grenzziehung, die nur in
  Code-Kommentar und Commit-Message steht, nicht am Plan. **(c)** ist in §3
  nachgezogen (Port-Zeile in `193716e`, allerdings mit dem Code im selben
  Commit statt vor ihm); **(a)** ist als Vorbestand bestätigt (siehe
  Design-Entscheidungen). Die Klasse steht damit beim **vierten** Auftreten —
  die Konflikt-Sequenz (Modul 8) läuft bereits seit review-slice-003; die
  Plan-Nachzüge (b/d) gehen als Übergabe-Artefakt an den Planner, nicht als
  informelle Handoff-Notiz.
- `verifizierbar`: ja — `SPEC-001`-Tabellen-Menge gegen die DDL; §1-Ausschlüsse
  gegen die gelieferte DDL; §3-Zeilen gegen die gelieferten Dateien
- `klasse`: Plan-Erweiterung ohne Plan-Nachzug (4. Auftreten — Sequenz läuft)

### F-2 — `harness/image-hash.txt` ist am Range-Head veraltet (Beleg-Kette gebrochen)

- `kategorie`: MEDIUM
- `quelle`: `harness/README.md` §Sensors (Zeile `make image`: „der Beleg gilt
  am HEAD: ein Zug, der Build-Kontext-Dateien ändert, erneuert ihn vor seiner
  Closure (Re-Build + Commit)") · [`ADR-0039`](../plan/adr/README.md) (Reproduzierbarkeits-Anker)
- `pfad`: `harness/image-hash.txt` (letzte Erneuerung `423cc5a` — **vor** der
  Range) gegen `go.mod`/`go.sum` (in `37405a6` geändert, pgx/v5 v5.11.0 + vier
  indirekte Dependencies) und `Dockerfile:14` (`COPY go.mod go.sum` in den
  deps-Layer)
- `befund`: Die Range ändert zwei Build-Kontext-Dateien, die in den deps-Layer
  kopiert werden — der Layer-Inhalt ändert sich, damit der Digest des gebauten
  Images. Die Implementer-Meldung „image-Digest unverändert (Binary importiert
  Adapter nicht — Verdrahtung slice-006)" trägt die *Runtime*-Verdrahtung,
  aber nicht die *Build-Kontext*-Regel: die deps-Layer-Änderung allein hebt den
  Digest. Der Beleg ist am HEAD nicht mehr gültig; vor der Closure ist er per
  `make image` zu erneuern (bewusster Commit, Modul 14).
- `verifizierbar`: ja — `git log -n1 -- harness/image-hash.txt` gegen die
  Range; `make image` am HEAD liefert einen anderen Digest
- `klasse`: Beleg am HEAD veraltet (Beleg-Kette gebrochen)

### F-3 — Fehlerklassen-Übersetzung (`ADR-0023`/`SPEC-008`) ohne Träger am ersten realen Adapter

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0023`](../plan/adr/README.md) (Konsequenz: „jeder Adapter trägt eine
  Übersetzungsverantwortung") · [`SPEC-008`](../../spec/pflichtenheft.md) (sieben stabile Kategorien;
  `storage` → kein Source-ACK) · Maintainability (unklare Fehlerbehandlung am
  Rand des Spec-Bereichs)
- `pfad`: `internal/adapters/driven/postgresstorage/store.go:31-44`
  (`New` meldet den Ping-Fehler roh; Kommentar zitiert „(`SPEC-008`, Klasse
  `storage`)"), `store.go:105-124` (`PersistTransaction` reicht rohe
  pgx-Fehler durch), `store.go:129-137` (`ReadChanges` ebenso)
- `befund`: Der erste reale Treiber-Adapter meldet technische Fehler roh
  (pgx-Verbindungs-, Constraint- und Pool-Fehler) an den Port. [`ADR-0023`](../plan/adr/README.md)
  legt die Übersetzung in stabile Kategorien als Adapter-Pflicht fest; der
  Adapter-Kommentar behauptet die Klasse `storage` für den Ping-Fehler, ohne
  sie zu tragen. Weder ist die Übersetzung implementiert noch die Grenze im
  Plan §6 mit Ausgang terminisiert. Für den Capture-Pfad (der bei *jedem*
  Fehler ohne ACK bleibt) ist das folgenfrei; für die Verbräucher der Klassen
  (slice-006-Verdrahtung, Betrieb über `cdc_errors_total{class}`) ist die
  Grenze unbenannt.
- `verifizierbar`: ja — Signatur-/Fehlermenge des Adapters gegen [`ADR-0023`](../plan/adr/README.md)
  und `SPEC-008`; kein Klassifizierungs-Träger im Adapter, kein §6-Ausgang
- `klasse`: ADR-Verantwortung ohne Träger (unbenannte Grenze am Spec-Rand)

### F-4 — Referenzen auf superseded ADRs ohne Fortgeltungs-Anker

- `kategorie`: LOW
- `quelle`: AGENTS.md §5 / Traceability-Regeln („keine superseded-ADR-Referenzen"
  als Negativbefund-Prüfpflicht seit review-slice-003) · [`ADR-0042`](../plan/adr/README.md)
  (Rest-Fortgeltung des [`ADR-0039`](../plan/adr/README.md)-Bestands)
- `pfad`: Plan-§3-Zeilen (neu in `193716e`: „je [`ADR-0039`](../plan/adr/README.md)" für
  row-Typen/queries/mapper; „[`ADR-0038`](../plan/adr/README.md) bleibt fortgeltend"),
  `internal/adapters/driven/postgresstorage/queries/queries.go:1-3` und
  `mapper/mapper.go:1-4` (Paket-Kommentare mit [`ADR-0039`](../plan/adr/README.md)), Commit `37405a6`
  („[`ADR-0038`](../plan/adr)-Fortgeltung"), `Dockerfile:17` („([`ADR-0039`](../plan/adr), fortgeltend)" —
  Vorbestand, von der Range-Edition des Kommentarblocks berührt)
- `befund`: [`ADR-0038`](../plan/adr/README.md) ist Superseded by [`ADR-0039`](../plan/adr/README.md), [`ADR-0039`](../plan/adr/README.md) Superseded by
  [`ADR-0042`](../plan/adr/README.md) — die getragenen Regeln (Go, CGO-frei, nativer Treiber,
  Row-Typen im Driven-Adapter) laufen über die Fortgeltungskette 0038→0039→0042
  und sind inhaltlich unberührt. Die neuen Referenzen im Plan und in den
  Code-Kommentaren nennen aber die superseded ADRs **ohne** den tragenden
  Anker (`ADR-0042` bzw. dessen Rest-Fortgeltungs-Klausel) — ein Leser am
  Range-Head schlägt „Superseded" nach und findet die Regel nicht.
- `verifizierbar`: ja — Referenzen im Diff gegen die Status-Zeilen der
  ADR-Dateien
- `klasse`: superseded-ADR-Referenz ohne Fortgeltungs-Anker

### F-5 — `tools/harness/run-store-tests.sh`: kein Datei-Abschluss, Test-Netz bleibt stehen

- `kategorie`: LOW
- `quelle`: Maintainability (Datei-Abschluss, 2. Auftreten der Klasse — F-6
  review-slice-001) · Testcontainer-Hygiene
- `pfad`: `tools/harness/run-store-tests.sh:63` (Datei endet ohne
  Zeilenumbruch) · `run-store-tests.sh:16-26` (`cleanup()` räumt nur den
  Container, nicht das Netz `cdc-store-test`)
- `befund`: Die neue Datei endet ohne Zeilenumbruch (zweites Auftreten der
  Klasse in diesem Repo). Der Aufräumpfad entfernt den PG-Container in jedem
  Ausgang, lässt aber das Docker-Netz zurück — nach dem Lauf dieser Review
  bleibt `cdc-store-test` bestehen. Daten und Caches liegen korrekt im
  Container/Volume, nichts im Arbeitsbaum; der Rest ist Hygiene am Rand.
- `verifizierbar`: ja — `tail -c1 tools/harness/run-store-tests.sh`;
  `docker network ls` nach einem `make test-store`-Lauf
- `klasse`: Datei-Abschluss (2. Auftreten) · Netz-Rückbau fehlt

### F-6 — Unerreichbare Bereichs-Prüfung in `mapper.ToPosition`

- `kategorie`: LOW
- `quelle`: Maintainability (toter Zweig als Grenz-Träger)
- `pfad`: `internal/adapters/driven/postgresstorage/mapper/mapper.go:93-95`
- `befund`: `uint64(commitPosition) > math.MaxInt64` ist für einen `int64 ≥ 1`
  unerreichbar — die Spalten-`CHECK`-Kante hält den Wert positiv (Kommentar
  benennt genau das), damit ist der Vergleich immer falsch. Die Grenze trägt
  allein die `commitPosition < 1`-Zeile; der tote Zweig liest sich als zweite
  Grenze und ist keine.
- `verifizierbar`: ja — Typ-Range-Analyse; Test
  `TestTransactionRowRoundTripsPosition` deckt den Zweig nie
- `klasse`: Toter Grenz-Zweig

### F-7 — `make gates` kompiliert Go nicht (Implementer-Risiko (b) bestätigt)

- `kategorie`: INFO
- `quelle`: `harness/mk/{baseline.mk,doc-gate.mk}` (GATE_CHECKS =
  baseline-verify, docs-check) · `a-check.mk` (liest Pfad-/Import-Ebene,
  „Kein Runtime-Urteil", `harness/sensors/a-check.md` Grenze 4) ·
  Implementer-Risiko (b) des Handoffs
- `pfad`: `Makefile:7-20,56-57` · `docs/plan/planning/in-progress/
  slice-004-postgres-store-adapter.md` §2 (DoD: „`make gates` grün")
- `befund`: Kein Gate des `make gates`-Bündels kompiliert oder testet den
  Go-Baum — ein nicht kompilierender Baum geht grün durch die Gates (dieser
  Review lief den Compile über einen eigenen Container-Lauf, nicht über ein
  Gate). Die Lücke ist seit dem ersten Binary real und wächst mit jedem
  Go-Slice. Notiz mit Zuständigkeit: **Architect/Planner** —
  Kompilier-/Test-Check als Gate (oder bootstrap-aware Stufung) ist eine
  Entscheidung, kein Implementer-Schritt; Kandidat für den
  Steering-Loop-Lerneintrag der Closure §7.
- `verifizierbar`: ja — `grep -rn "GATE_CHECKS +=" harness/mk/*.mk`;
  absichtlicher Kompilierfehler durch `make gates` (bleibt grün)
- `klasse`: Gate-Lücke konkretisiert (kein Compile-/Test-Gate)

### F-8 — Leere committed Transaktion: Grenze am Adapter-Test getragen, Port-Kontrakt-Zeile nennt sie nicht

- `kategorie`: INFO
- `quelle`: review-slice-003 F-4 („unbenannte Grenze am öffentlichen Vertrag")
  · `LH-FA-CAP-006.a`
- `pfad`: `internal/adapters/driven/postgresstorage/store_test.go` (Test
  `TestPersistCarriesEmptyCommittedTransaction`) ·
  `internal/application/port/outbound/changestore.go:93-97`
  (`PersistTransaction`-Kontrakt: „mit interner ID, Commit-Position und
  Changes")
- `befund`: Die Grenze aus slice-003 F-4 ist am realen Pfad jetzt **belegt**
  (leere committed Transaktion persistiert ihre Position, liest keinen
  Change) — der Befund ist damit am Verhalten benannt. Die Kontrakt-Zeile am
  Port nennt den Fall weiterhin nicht; mit dem ersten realen Treiber, der sie
  ausübt, ist sie eine Notiz für den Planner (Kontrakt-Schärfung), keine
  Nacharbeit am Diff.
- `verifizierbar`: ja — Test gegen die reale Instanz (läuft grün in
  `make test-store`)
- `klasse`: Grenze am Test getragen, am Kontrakt unbenannt

---

## Design-Entscheidungen des Implementers — Bewertung (Prüf-Fragen)

- **(a) Plan-Defekt „`make gates` grün" doppelt in §2:** bestätigt —
  `git show c9a068f` trägt die Dopplung bereits; sie ist ein Vorlagen-Glitch
  vor dem Slice und im Diff **nicht** berührt. Bewertung: Plan-Defekt-Korrektur
  durch den **Planner** (eine Zeile streichen), keine Plan-Änderung des
  Implementers, kein Finding gegen den Diff.
- **(b) „Position-Verwaltung" enger gefasst:** *Plan-Ergänzung, nicht
  Plan-Änderung* — die enge Lesung trägt sich aus §1 (Replication-Stream →
  slice-005, Consumer-Use-Cases → anderer Vorgang); kein §1-Ausschluss und
  kein DoD-Punkt berührt. Aber: F-1 — die Ergänzung steht nur im Handoff und
  braucht den Plan-Nachzug; die Klasse zählt (4. Auftreten).
- **(c) `ErrRangeInverted`/`ErrNonPositiveLimit`/`ChangeQuery`/`ChangeRecord`
  am Port:** [`ADR-0042`](../plan/adr)-konform — „Kontrakt-Typen eines Outbound-Ports leben
  analog am Port"; die Sentinels tragen die Port-Vertragsgrenzen (REA-001/003
  Negative), die Domänen-Ordnung bleibt Domänenwert (`ADR-0005`/`ADR-0029`).
  Plan-Nachzug ist fällig und **erfolgt** — §3-Zeile in `193716e` (mit dem
  Code im selben Commit, nicht vor ihm — notiert, keine eigenständige Regel
  verletzt). Kein Architekt-Verdikt nötig: die Platzierung folgt wörtlich dem
  Accepted-[`ADR-0042`](../plan/adr).
- **(d) DDL ohne Consumer-/Capture-State-Tabellen:** sachlich tragfähig
  (die Ports liegen außerhalb des Store-Adapters; die drei Paarungen und §1
  werden nicht berührt), aber formal eine **Plan-Änderung** — siehe F-1.

## ADR-Deckung (Adapter, Port, DDL)

| ADR | Aussage | Träger im Diff |
|---|---|---|
| [`ADR-0011`](../plan/adr/README.md) | Persist-before-ACK; idempotente Persistenz ist Pflicht | getragen — Transaktion + Changes in EINEM Store-Commit; ON CONFLICT DO NOTHING auf beiden Primärschlüsseln; Idempotenz jetzt **auch am Port-Kontrakt** (`changestore.go:87-91` — F-6 aus review-slice-003 ist aufgelöst); Mutations-Probe belegt die Tests |
| [`ADR-0009`](../plan/adr/README.md) | Ein Fähigkeits-Port für Lesen und Schreiben (Option C) | getragen — `ReadChanges` am `ChangeStorePort`, keine zweite Persistenz-Grenze; Port-Vertrag erweitert, Fakes nachgezogen |
| [`ADR-0010`](../plan/adr/README.md)/[`ADR-0032`](../plan/adr/README.md) | PostgreSQL als Store, Adapterdetail | getragen — pgx/v5-Typen bleiben im Adapter (pgxpool, pgx.Rows); Core sieht keinen Treibertyp |
| [`ADR-0032`](../plan/adr/README.md)-Fortgeltung ([`ADR-0038`](../plan/adr/README.md)-Kette) | nativer Go-Stack, CGO-frei | getragen — pgx/v5 v5.11.0, rein Go; Dockerfile `CGO_ENABLED=0` unverändert; Referenz-Hygiene F-4 |
| [`ADR-0023`](../plan/adr/README.md)/[`SPEC-008`](../../spec/pflichtenheft.md) | Fehlerklassen | **F-3** — rohe pgx-Fehler ohne Klassifizierungs-Träger; der Kommentar zitiert die Klasse |
| [`ADR-0029`](../plan/adr/README.md) Regel 1/3/6/7 | ACK-Nur-nach-Persistenz · offene Transaktionen unkonsumierbar · eindeutige Sequenz · Schema-Version | getragen — Port-Kontrakt-Zeile (Regel 1), `ErrTransactionNotCommitted`-Guard + Test (Regel 3), `UNIQUE (transaction_id, sequence)` + `CHECK (sequence >= 1)` in der DDL (Regel 6), `schema_version NOT NULL REFERENCES` (Regel 7) |
| [`ADR-0042`](../plan/adr/README.md) | Kontrakt-Typen am Port | getragen — `ChangeQuery`/`ChangeRecord`/Sentinels am Outbound-Port; keine Neuedefinition im Use Case (dort existiert noch kein lese-seitiger Use Case — der Re-Evaluierungs-Trigger bleibt unberührt) |
| [`ADR-0005`](../plan/adr/README.md) | Positionen nur innerhalb ihrer Quelle | getragen — `Validate` verwirft fremd-quellige Grenzen über `ErrSourceMismatch` (+ Test) |

## Neue Angriffsfläche — erster realer Treiber-Code

- **SQL-Injection:** keine — alle vier Queries tragen Parameter über das
  Extended-Protokoll (`$1…$5`), keine String-Interpolation; Tabellen-/Spalten-
  namen stehen als Konstanten in `queries/`. Die `fmt.Sprintf`-Interpolation in
  `store_test.go seedReference` nutzt nur Konstanten (Test-Code, kein
  Eingabepfad).
- **Credentials in Logs:** keine — der Adapter loggt nichts; die DSN-Error-
  Probe (fehlerhafte DSN mit Passwort gegen pgx v5.11.0) zeigt, dass der
  Treiber das Passwort in Parse-Fehlern **redaktiert** (`cdc:xxxxx@…`); die
  Grenze „Keine Credentials in Logs" (`Pflichtenheft`) bleibt gewahrt.
- **Testcontainer-Hygiene:** Daten im Container (`DROP SCHEMA … CASCADE` je
  Test, kein Volume in den Arbeitsbaum); Test-Läufe mounten `/src` read-only;
  Cache in `/tmp` und Docker-Volume; Rest: Netz bleibt stehen (F-5). Kein
  Arbeitsbaum-Eintrag (Modul 14 Besitz-Regel).
- **Digest-Pins:** beide Images gepinnt — Toolchain-Digest identisch zum
  Dockerfile (`cf6fca…`), PG-Digest `63bdc97…` in Makefile und Skript
  konsistent; Herkunft („docker manifest inspect postgres:18-alpine, amd64")
  im Kopf benannt.

## Test-Qualität — reale Treiber-Läufe und Mutations-Proben

`make test-store` am Range-Head: 8 Tests gegen die reale
Testcontainer-PostgreSQL, grün; treiberfreie `go test ./...` grün (5
Test-Pakete). Eigene Mutations-Proben gegen die Store-Suite (Worktree
`/tmp/rev-s004`, außerhalb des Arbeitsbaums, gepinnter Treiber-Container):

| Mutation | Erwartung | Ergebnis |
|---|---|---|
| End-Exklusivität gelockert (`commit_position < $3` → `<=`) | `TestReadCarriesPositionsAndRanges` rot | **rot** (Bereichsgrenze verletzt) |
| Deduplizierung entfernt (beide `ON CONFLICT … DO NOTHING` gestrichen) | `TestPersistTransactionIsIdempotent` rot | **rot** (Idempotenz verletzt) |

Die kritischen Zusagen — Bereichs-Ordnung (`LH-FA-REA-001/002`) und
Idempotenz (`ADR-0011`) — tragen ihre Tests, nicht nur ihre Kommentare. Die
drei Gates am Range-Head: `baseline-verify` OK · `d-check` 89 Dateien/0
Befunde · `a-check` 0 Befunde (keine neue Kante; `adapters → ports`/
`ports → domain` erfüllt).

## Implementer-Risiken — Bewertung

- **(a) image-Digest unverändert:** **F-2** — die Begründung (Binary importiert
  den Adapter nicht) trägt die Runtime-Verdrahtung, aber nicht die
  Build-Kontext-Regel; `go.mod`/`go.sum` ändern den deps-Layer, der Beleg ist
  am HEAD veraltet. Erneuerung vor der Closure.
- **(b) `make gates` kompiliert Go nicht:** **F-7** — bestätigt; Zuständigkeit
  Architect/Planner, Lerneintrag-Kandidat für die Closure §7.
- **(c) FK-Voraussetzung (Metadaten vor Persistenz):** kein Befund — die
  Grenze ist dreifach benannt (`schema.sql` Kopf-Kommentar, `schema.go`
  `ApplySchema`-Kommentar, `seedReference`-Testkommentar), im Indikativ und
  klassenkonform (§3.7).

## Plan-Nachzug — die drei Meldungen in §3/§7 (Planner-Sache, notiert)

(a) ist Vorbestand und Planner-Korrektur (F-1-Kontext, keine Plan-Änderung im
Diff) · (b) und (d) brauchen den Nachzug vor der Closure (F-1, Sequenz läuft)
· (c) ist in §3 nachgezogen (`193716e`). Die §7-Füllung ist Closure-Sache und
steht hier nicht zur Prüfung an.

## Negativbefunde

- geprüft, ohne Befund: **HIGH-Klassen über den Diff** — kein ADR-Verstoß auf
  Layer-Ebene (`a-check` 0 Befunde, keine undeclarierte Import-Kante, kein
  `ports → app`-Pfad), kein Sicherheits-Anti-Pattern (parametrisierte Queries,
  keine Injection-Fläche, Credential-Redaktion des Treibers probe-belegt),
  kein Korrektheitsfehler im kritischen Pfad (EIN Store-Commit je Transaktion,
  Mutations-Proben rot auf beiden Zusagen, Lesen trägt nur SELECT und
  verändert keine Position — `TestReadLeavesPersistedStateUnchanged`), keine
  Gate-Suppression, keine Norm nur im Template-Kommentar, kein
  Chronik-tragendes Zustandsfeld im Diff, kein Docker-only-Verstoß (alle
  Läufe über gepinnte Container, Modul-Cache im Volume, Daten im Container)
- geprüft, ohne Befund: **Stand-alone-Build am Range-Head** — `go build ./...`
  und `go vet ./...` im gepinnten Toolchain-Container grün; `go test ./...`
  treiberfrei grün; `make test-store` grün (8 Store-Tests gegen reale
  PostgreSQL)
- geprüft, ohne Befund: **Traceability-Grundpflege der drei Commits** — jeder
  trägt mindestens eine `LH-*`-/`ADR-*`-Kennung; alle genannten IDs existieren
  ([`LH-FA-REA-001`](../../spec/lastenheft.md)…006, [`LH-FA-DAT-002`](../../spec/lastenheft.md)/004, [`LH-FA-CAP-006`](../../spec/lastenheft.md)/008, [`LH-FA-RET-001`](../../spec/lastenheft.md),
  [`LH-FA-SCH-005`](../../spec/lastenheft.md), [`LH-QA-REL-001`](../../spec/lastenheft.md)/002, [`ADR-0009`](../plan/adr)/0010/0011/0032/0038/0039/0042);
  **keine `SPEC-*`-/`ARC-*`-Kennung in einer der drei Commit-Messagen** — die
  F-3-Klasse aus review-slice-003 tritt hier **nicht** erneut auf (4.
  Zählstand bleibt stehen); keine undeklarierten Präfixe (die
  superseded-Referenzen sind F-4, kein Präfix-Verstoß)
- geprüft, ohne Befund: **Spec-Stratum** — der Diff berührt keine Spec-Datei;
  die `SPEC-*`-Nennungen in Code-Kommentaren und DDL sind Rang-Zeiger, keine
  Erweiterung des Technik-Stratums
- geprüft, ohne Befund: **Kommentar-Klassen im neuen Code (§3.7)** —
  `queries.go`, `schema.go`, `schema.sql`, `store.go`, `mapper.go`,
  `changestore.go` (inkl. der Idempotenz-Kontrakt-Zeile), beide Testdateien,
  `run-store-tests.sh`, Makefile-Block — Klassen
  Zusage/Kopplung/Abgrenzung/Rang-Zeiger/Grenze im Indikativ, kein
  Konjunktiv über verworfene Alternativen, kein abwesender Text, kein
  abgebrochener Satz
- geprüft, ohne Befund: **§3-Datei-Menge des Plans (am Range-Head)** — alle
  gelieferten Dateien sind durch die (in `193716e` nachgezogenen) §3-Zeilen
  gedeckt; keine unbudgetierte Datei; die Dockerfile-Edition ist
  Kommentar-Pflege des deps-Absatzes, kein neuer Liefer-Punkt
- geprüft, ohne Befund: **Zwei-Quellen-Drift im Beobachtungs-Register** —
  `BEO-PGC/a-check-null-abdeckung/state.md` trägt den Zähler abgeleitet
  („Zähler (abgeleitet): … 3×, Ausgang im Lese-Schritt der Welle-1-Closure
  zugewiesen"); der Defekt F-7 aus review-slice-003 ist bei der Closure
  berichtigt. Neue Beobachtung aus diesem Diff angefallen: keine (kein
  zweites Auftreten einer neuen Klasse unter der Schwelle)
- geprüft, ohne Befund: **slice-003-F-6 (Idempotenz-Pflicht ohne
  Kontrakt-Träger)** — aufgelöst; die Kontrakt-Zeile steht am Port
  (`changestore.go:87-91`) und der Adapter belegt sie real

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 3 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Plan-Erweiterung ohne Plan-Nachzug (4.
Auftreten — Sequenz läuft) · Beleg am HEAD veraltet (image-hash) ·
ADR-Verantwortung ohne Träger (Fehlerklassen-Übersetzung) · superseded-ADR-
Referenz ohne Fortgeltungs-Anker · Datei-Abschluss (2. Auftreten) · Toter
Grenz-Zweig · Gate-Lücke konkretisiert (kein Compile-/Test-Gate) · Grenze am
Test getragen, am Kontrakt unbenannt

## Verdikt

**Merge-blockierend:** nein — der Diff ist inhaltlich schlüssig, kompiliert
und vet-t sauber stand-alone, die kritischen Zusagen (Idempotenz, Bereichs-
Ordnung, Lesen-ohne-Positionsänderung) sind gegen den realen Treiber belegt
und mutations-geprüft, die Port-Erweiterung folgt wörtlich [`ADR-0042`](../plan/adr), und die
drei Gates laufen grün.

**Blockierend für Closure:** ja, in drei Punkten, bevor der Slice nach
`done/` geht:

1. **F-2:** `harness/image-hash.txt` per `make image` am HEAD erneuern
   (bewusster Commit, Modul 14) — der Beleg gilt am HEAD, und die Range hat
   den deps-Layer geändert.
2. **F-1 als Sequenz-Fortsetzung (Modul 8, 4. Auftreten):** Plan-Nachzug für
   (b) (Ergänzungs-Zeile) und (d) (§1-Ausschluss mit Begründung „ein
   Folge-Slice übernimmt es" oder §3-Schmälerung) über den **Planner**;
   Übergabe-Artefakt: dieser Report-Abschnitt + Handoff-Meldung. (a) geht als
   Vorlagen-Glitch-Korrektur in denselben Planner-Zug.
3. **F-3:** die Fehlerklassen-Grenze ist zu benennen — entweder trägt der
   Adapter die `SPEC-008`-Übersetzung, oder der Plan §6/§1 terminisiert den
   Ausgang (welcher Slice trägt die Übersetzung). Stillschweigend darf der
   Adapter sie weder behaupten (Kommentar) noch fehlen lassen.

F-4/F-5/F-6 sind Vor-Closure-Nacharbeit ohne Blockier-Charakter (Anker-Hebung,
Datei-Abschluss, toter Zweig). F-7 geht als Steering-Loop-Kandidat in die
Closure §7 (Lerneintrag: Kompilier-/Test-Check als Gate — Architect/Planner).
DoD- und Spec-Konformität prüft der Verifier separat (Modul 11) —
insbesondere `make gates` grün als beobachtbarer Beleg (der nach F-2 den
erneuerten Image-Beleg einschließen sollte).