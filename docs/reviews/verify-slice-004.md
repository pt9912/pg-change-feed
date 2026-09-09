# Verifier-Report: slice-004 — 2026-09-09

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of Done,
10 Items nach Entfernung des doppelten „make gates"-Items in `13abe88`), §3
(Plan-vs-Code, Range `c9a068f..46fbf33` inkl. Plan-Nachzug `13abe88`), §6
(Risiko-Ausgänge) und ADR-Konformität ([`ADR-0011`](../plan/adr/0011-persist-before-ack.md) ·
[`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md)/[`SPEC-008`](../../spec/pflichtenheft.md) ·
[`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md) ·
[`ADR-0010`](../plan/adr/0010-postgresql-cdc-store.md)/[`ADR-0032`](../plan/adr/0032-postgresql-adapterdetail.md) ·
[`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) ·
[`ADR-0029`](../plan/adr/0029-domain-invarianten.md) Regel 3/6/7). Nicht geprüft:
Diff gegen Plan/Hard Rules (Reviewer, `review-slice-004.md`, Verdikt dort),
realer Bedarf (Validator). Expliziter Auftrag: Schiedsspruch zum
Review-F-2-Widerspruch (image-hash-Beleg).

**Gegenstand:** Implementer-Handoff zu
`docs/plan/planning/in-progress/slice-004-postgres-store-adapter.md` · Range
`c9a068f..46fbf33` (4 Commits: `cefa18c`, `37405a6`, `13abe88`, `46fbf33`) ·
Fortsetzung: `59e3568` (d-migrate-Verdrahtung, [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)) · Arbeitsbaum-HEAD
**`96c47af`** (während dieses Laufs gelandet: Image-Digest-Semantik
präzisiert — Message zitiert „verify-slice-004", das ist dieser Bericht;
die Semantik-Antwort des Commits wird unten gegen meine Belege geprüft) ·
**Fix-Commit `46fbf33` trägt F-3/F-4/F-5/F-6; F-2 ist als Widerspruch belegt
und hier geschieden.**

**Grundsatz:** Es wurden **keine Behauptungen übernommen** — jeder Sensor
unten wurde in diesem Lauf selbst gefahren; Ausgaben sind Belege über stdout.
Der Hauptbaum blieb read-only (`git status` clean am Ende, `git worktree
list` trägt nur den Hauptbaum); Image-Builds und Mutations-Probe liefen in
`/tmp`-Worktrees (nach Abschluss entfernt), die reale Testcontainer-
PostgreSQL stand in einem eigenen Container/Netz (nach Abschluss abgeräumt).

**Modell:** Claude Code (glm-5.3-flash) · **Datum:** 2026-09-09

**Eingangs-Kontext:**

- Slice-Plan §1–§8 (am HEAD `96c47af`) · `harness/README.md` (Sensors-/
  Werkzeuge-Tabellen) · `Makefile`, `a-check.mk`, `harness/mk/*.mk`,
  `.a-check.yml`, `Dockerfile`
- [`ADR-0011`](../plan/adr/0011-persist-before-ack.md) ·
  [`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md) ·
  [`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md) (Accepted) ·
  [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) (Accepted,
  neu) · [`ADR-0010`](../plan/adr/0010-postgresql-cdc-store.md)/
  [`ADR-0032`](../plan/adr/0032-postgresql-adapterdetail.md)/
  [`ADR-0029`](../plan/adr/0029-domain-invarianten.md)
- `docs/reviews/review-slice-004.md` (F-1…F-8) · Fix-Commit `46fbf33` im
  Volltext · Plan-Nachzug `13abe88` · `59e3568` · `96c47af` im Volltext
- `spec/pflichtenheft.md` (§2 Tabellen-Menge inkl. `cdc.consumer`,
  `cdc.consumer_position`, `cdc.capture_state`; `SPEC-001`/002/004/008) ·
  `spec/lastenheft.md` ([`LH-FA-REA-001`](../../spec/lastenheft.md)…006,
  [`LH-FA-RET-001`](../../spec/lastenheft.md), [`LH-QA-REL-001`](../../spec/lastenheft.md)/003)
- Beobachtungs-Register (`BEO-PGC/a-check-null-abdeckung`, Stand und
  evidence/)

---

## F-2-Schiedsspruch (expliziter Auftrag)

**Streit:** Review-F-2 behauptet, `harness/image-hash.txt` sei am Range-Head
veraltet — die Range ändere `go.mod`/`go.sum` (im deps-Layer), damit den
Digest. Der Implementer meldet, der Digest bleibe unverändert (Runtime-Stage
erbt den deps-Layer nicht, Binary byte-identisch).

**Belege, in diesem Lauf selbst gefahren** (jeder Build im `/tmp`-Worktree,
Hauptbaum unberührt; derselbe Builder, derselbe Buildx-Aufruf wie
`make image` — `docker buildx build --load --metadata-file …`):

| Build | Commit | containerimage.digest | Binary (sha256) |
|---|---|---|---|
| A | HEAD `59e3568` (Range-Ende) | `sha256:01b46e36683818d27b2707c2e154de93213df49c044314a94da6933c49557d8f` | `43c3aec075e989db779845f0f97cb5a629718baf54edc991a9efecff40947b4c` |
| B | `423cc5a` (Vor-pgx-Stand, Aufzeichnungs-Commit des Belegs `447eab…`) | `sha256:01b46e36683818d27b2707c2e154de93213df49c044314a94da6933c49557d8f` | `43c3aec075e989db779845f0f97cb5a629718baf54edc991a9efecff40947b4c` (byte-identisch zu A) |
| C | HEAD `96c47af` (Dockerfile-Kommentar geändert — Build-Kontext erneut geändert) | `sha256:01b46e36683818d27b2707c2e154de93213df49c044314a94da6933c49557d8f` | — |
| `--provenance=false` (A) | `59e3568` | identisch zu A | — |

Eingetragener Beleg: `harness/image-hash.txt` = `sha256:447eab3690a67c2669698c5e470b7277f8657082c99a75d18f5c4aa733f3273f`
(letzte Erneuerung `423cc5a`).

**Verdikt:**

1. **Die Review-F-2-Fassung ist falsch.** Der deps-Layer hebt den Digest
   dieses Images nicht: Der Runtime-Stage erbt ihn nicht (`FROM deps AS
   build`; Runtime `FROM distroless` + `COPY --from=build /out/…`), und der
   kontrollierte Doppel-Build (A vs. B) liefert **identischen Digest und
   byte-identisches Binary** trotz geänderten Build-Kontexts. Auch
   Dockerfile-Kommentaränderungen (B→C) bewegen den Digest nicht.
2. **Aber auch die Review-„verifizierbar"-Prozedur trägt nicht:** Mein
   Re-Build am *Aufzeichnungs-Commit* `423cc5a` liefert `01b46e…`, nicht den
   eingetragenen `447eab…`. Der Digest ist **builder-gebunden** — der
   Vergleich „frischer Build gegen eingetragenen Digest" über Umgebungen
   hinweg kann weder Staleness belegen noch widerlegen. Der Streit ist am
   Digest allein **nicht entscheidbar**; entscheidbar ist er am Inhalt — und
   der trägt (byte-identisches Binary, unveränderte Base-Pins).
3. **Die Implementer-Position trägt in der Sache**; die Semantik-Frage
   („Target-Semantik, cache-freier Lauf" als Architect-Target-Änderung) war
   berechtigt und ist mit `96c47af` beantwortet: Der Beleg ist der Digest des
   **exportierten Images** (Runtime-Stage: Basis + Binary); deps-Layer-
   Änderungen ohne Binary-Änderung lassen ihn unverändert — „unverändert ist
   ein gültiger Befund, kein Staleness-Zeichen". Der Re-Build-Pflicht-Teil
   der Regel bleibt (Re-Build vor Closure), der Commit entfällt bei
   unverändertem Digest. Meine Builds bestätigen genau diese Semantik.
4. **Residuum, nicht blockierend:** Die Formulierung „ändert sich genau dann,
   wenn sich das exportierte Image ändert" gilt **innerhalb einer Builder-
   Umgebung**; über Umgebungen hinweg ist der Digest nicht stabil (B.2). Der
   Dockerfile-Kopf trägt die Grenze bereits korrekt („Beleg, kein
   Wiederholungs-Schlüssel") — die Kette lautet über den Inhalt, nicht über
   den Digest-Vergleich.

**Folge für die Closure:** Kein Erneuerungs-Zwang aus F-2 in der
Review-Fassung. Die Semantik-Änderung `96c47af` ist als Architektur-Entscheidung
dokumentiert (Dockerfile-Kopf + Werkzeuge-Zeile; `make image` bleibt kein
Gate, AGENTS §3.6 greift nicht) — Bindungs-Anker der Zeile ist jetzt
[`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md) (Fortgeltungs-Rest der
[`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md)-Kette), konsistent mit der
Anker-Hebung des F-4-Fixes.

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` (HEAD `59e3568`) | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 91 Datei(en) geprüft, 0 Befund(e)` · `a-check … gesamt: 0 Befund(e)` — a-check im Bündel; der `adapters`-Glob (`internal/adapters/**`) matcht den echten Adapter-Content | 0 |
| `make gates` (HEAD `96c47af`, erneuter Lauf nach dem Zwischen-Commit) | dieselbe Ausgabe, 0 Befunde | 0 |
| `go build ./...` + `go vet ./...` + `gofmt -l .` im gepinnten Toolchain-Container (`golang:1.27-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125`, go1.27.1, `/src:ro`, `--network none`) | `VET_CLEAN` · `GOFMT_CLEAN` | 0 |
| `make test` (netzlos, gepinnter Container) | `ok …` für 5 Test-Pakete (`postgresstorage` skippt ohne DSN, `mapper`, `port/outbound`, `usecase/capture`, `domain/model`) | 0 |
| `make test-store` | `ok … internal/adapters/driven/postgresstorage 1.422s` gegen die gepinnte `postgres:18-alpine@sha256:63bdc97d…` | 0 |
| Store-Paket verbose, eigene Instanz desselben gepinnten PG-Images | **10/10 PASS** — `TestPersistTransactionIsIdempotent`, `TestPersistRejectsOpenTransaction`, `TestPersistCarriesEmptyCommittedTransaction`, `TestReadIsDeterministicallySorted`, `TestReadCarriesPositionsAndRanges`, `TestReadCarriesLimit`, `TestReadCarriesTableFilter`, `TestReadLeavesPersistedStateUnchanged`, `TestPersistCarriesStorageClass`, `TestNewCarriesStorageClass` | 0 |
| Mutations-Probe (Sentinel der Klasse `storage` entfernt: `storageFailure` gibt die rohe Ursache zurück; `/tmp`-Worktree, revertiert und entfernt) | `--- FAIL: TestPersistCarriesStorageClass` (`Fehler = ERROR: … SQLSTATE 23503, wollen Klasse storage`) · `--- FAIL: TestNewCarriesStorageClass` — **erwartet rot, rot gesehen** | 1 |
| Image-Builds A/B/C (Tabelle oben) | Digest-Invarianz + Binary-Identität | 0 |
| `make doc-commits RANGE=c9a068f..96c47af` | `d-check: 91 Datei(en) geprüft, 0 Befund(e)` (Traceability je Commit, volle Range) | 0 |
| `make doc-immutable RANGE=c9a068f..96c47af` | `d-check: 91 Datei(en) geprüft, 0 Befund(e)` (Immutabilität) | 0 |
| `make schema-validate` (ohne `tools/schema/schema.yaml`) | `FEHLER: tools/schema/schema.yaml fehlt — das neutrale Schema-YAML ist die Erstlieferung des d-migrate-Einbaus (ADR-0043)…` — sauberer, deklarierter Ausgang | 2 |
| `git worktree list` / `git status --porcelain` | nur der Hauptbaum; Arbeitsbaum clean — keine Einträge dieses Laufs | — |

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt — 10 Items nach Entfernung des Duplikats in `13abe88`)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | `PostgresChangeStoreAdapter` persistiert committed Transaktionen deduplizierbar (Idempotenz-Kontrakt [`ADR-0011`](../plan/adr/0011-persist-before-ack.md)) | **bestätigt** | `TestPersistTransactionIsIdempotent` grün gegen die reale gepinnte PostgreSQL (eigener Lauf, verbose-Beleg oben); Träger: Primärschlüssel `cdc.transaction(transaction_id)` + `ON CONFLICT DO NOTHING` (`schema.sql:53-58`, `queries.go`); Kontrakt-Zeile am Port (`changestore.go:95-99`); Teil-Beleg zu [`LH-FA-RET-001`](../../spec/lastenheft.md) |
| 2 | Deterministisches Lesen über [`SPEC-001`](../../spec/pflichtenheft.md)-Tabellen (Bereich/Limit/Filter) | **bestätigt** | Fünf Lese-Tests grün am realen Treiber (deterministische Sortierung, Bereich inklusiv/exklusiv, Limit, Tabellenfilter, Lesen ohne Positionsänderung); die Ordnung trägt `ORDER BY t.commit_position, c.transaction_id, c.sequence` (`queries.go:54`); Teil-Beleg zu [`LH-FA-REA-001`](../../spec/lastenheft.md)…006 |
| 3 | `make gates` grün | **bestätigt** | Drei Gates grün am HEAD — gefahren an `59e3568` **und** erneut an `96c47af` (Tabelle oben) |
| 4 | Review-Report unter `docs/reviews/` | **bestätigt — mit Einschränkung** | `docs/reviews/review-slice-004.md` committet (in `13abe88`); Rollenwechsel nach Schritt 8 eingehalten. Einschränkung: die dokumentierte Review-Basis `193716e` ist unerreichbar — V-2 |
| 5 | Doku-Update bzw. begründete Aussage „kein öffentlicher Vertrag berührt" | **nicht bestätigt in der Begründung (V-3)** | Die Aussage „kein öffentlicher Vertrag berührt (nur `internal/**`)" deckt die eigene Range nicht: `harness/README.md` (Rang 9, öffentlicher Einstiegsvertrag) wurde in `13abe88` (Werkzeuge-Zeilen `test`/`test-store`) und `59e3568` (Werkzeuge-Zeilen `schema-*`) geändert. Die *Änderungen selbst* sind korrekt (Regel 7 des Minimal Agent Workflow: neue Targets deklarieren); falsch ist allein die DoD-Begründung — das Häkchen steht auf einer Behauptung, die die Range selbst widerlegt |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **offen bis Closure** | §7 trägt Platzhalter; Slice liegt korrekt in `in-progress/` — fällig vor dem `git mv` (inkl. Finding-Klassen des Reviews, F-7 als Lerneintrag-Kandidat, V-1…V-6 dieses Berichts) |
| 7 | Reconciliation-Register | **entfällt** | `docs/plan/planning/reconciliation.md` existiert nicht (Repo ohne Brownfield-Bootstrap) — vom DoD-Wortlaut ausdrücklich vorgesehen |
| 8 | Beobachtungs-Register fortgeschrieben | **kein neues Auftreten — notiert** | `BEO-PGC/a-check-null-abdeckung` trägt den Ausgang `verkörpert · seit welle-1`; `evidence/` trägt `slice-001/002/003.md`. Kein neues Auftreten dieser Klasse: der `adapters`-Glob matcht seit diesem Slice echten Content, a-check 0 Befunde. **Neuer Kandidat für die Closure:** die F-7-Klasse des Reviews („kein Compile-/Test-Gate") ist eine *neue* Beobachtungsklasse, nicht das Wiederauftreten von `BEO-PGC` — ob sie als neues Verzeichnis ins Register geht, entscheidet der Lese-/Sichtungs-Schritt der Closure; die Zählung beginnt bei 1 (ein Vorgang, `slice-004`) |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **Voreinstellung getragen — Endwert fällig bei Closure** | (a) „reale Idempotenz-Deduplizierung ist stärker als der Fake — weiter offen": trägt jetzt den vorgesehenen Beleg (`TestPersistTransactionIsIdempotent` am realen Treiber); der Ausgang bleibt „weiter offen" mit slice-003-Termin (Welle 2) — tragende Lesart. (b) Row-Images/REPLICA IDENTITY — unverändert, Bewertung slice-006. Beide Endausgänge (eingetreten/entfallen/weiter offen) sind vor dem `git mv` zu setzen |
| 10 | Drei Paarungen (Anker · Folge-Slice · Register) | **delegiert — korrekt** | Repo **mit** Wellen-Betrieb (`docs/plan/planning/welle-2.md` flach vorhanden, Slice trägt `Welle: welle-2`): der DoD-Wortlaut weist die Prüfung der Welle-2-Closure zu; hier nicht fällig |

**Zwischenstand: 4/10 jetzt erfüllt** (Items 1–4) · 1 nicht bestätigt in der
Begründung (Item 5, V-3 — die getane Arbeit ist korrekt, die DoD-Aussage
falsch) · 1 entfällt (Item 7) · 3 erfüllen sich erst bei Closure (Items 6, 8
als Notiz, 9 Endwert) · 1 delegiert (Item 10).

## Plan-vs-Code-Diff (Range `c9a068f..46fbf33`, Plan-Nachzug `13abe88`)

**Plan-Korrekturen (`13abe88`) gegen den Auftrag:** bestätigt — (i) das
doppelte DoD-Item „`make gates` grün" ist entfernt (10 Items bleiben), (ii)
das Doku-Update-Item trägt seine Aussage (Begründung falsch, V-3), (iii) der
Plan-Nachzug steht in §1: (b) „Position-Verwaltung" als Plan-**Ergänzung**
(Commit-Positionen je Transaktion als Ordnungs-/Bereichs-Größe des Lesens;
Source-ACK-/Consumer-Positionen bei `ReplicationAckPort`/`ConsumerStatePort`,
slice-005/006) und die Konsequenz (d) als eine Plan-**Änderung**: die DDL
trägt die Consumer-/Capture-State-Tabellen von [`SPEC-001`](../../spec/pflichtenheft.md)
bewusst nicht (§3-Schmälerung + §1-Zeile). Code-seitig konsistent:
`schema.sql:7-17` nennt den Umfang im Kopf (5 Tabellen: `cdc.source`,
`cdc.source_table`, `cdc.schema_version`, `cdc.transaction`, `cdc.change`);
[`SPEC-001`](../../spec/pflichtenheft.md) §2 listet zusätzlich `cdc.consumer`,
`cdc.consumer_position`, `cdc.capture_state` — genau die drei, die der
Plan-Nachzug (d) ausnimmt, mit Begründung (ihre Ports liegen außerhalb des
Store-Adapters). **Die Grenzziehung ist damit am Plan, nicht nur im Code.**

**Deckung §3:** alle gelieferten Dateien sind durch die (in `13abe88`
nachgezogenen) §3-Zeilen gedeckt: `postgresstorage/*.go` (store, schema,
`queries/`, `mapper/`), `*_test.go` (store_test, mapper_test), `schema.sql`,
`changestore.go` (Port), `go.mod`/`go.sum` (pgx/v5 v5.11.0 + 4 indirekte
Deps), `Makefile` + `tools/harness/run-store-tests.sh`. **Unbudgetiert, klein
und zulässig:** `internal/application/port/outbound/changestore_test.go`
(neu, 80 Zeilen — §3 nennt den Port nur „geändert", nicht dessen Test; der
Umfang zählt zu `ADR-0042`-Kontrakt-Belegung, keine Eigenentscheidung) und
`usecase/capture/service_test.go` +7 (Fake-Nachzug der Port-Erweiterung).
Beide sind Test-Deckung bereits budgetierter Zeilen, keine neuen
Liefer-Punkte — notiert, kein Befund.

**Über die Range hinaus, am HEAD:** `59e3568` verdrahtet d-migrate
(`schema-validate`/`schema-rollout`) und `96c47af` präzisiert die
Image-Digest-Semantik. **Beide sind in §3 des Slice-Plans nicht budgetiert**
— V-6: entweder als §3-Plan-Erweiterung nachziehen oder als eigenes
wellenloses Slice führen (die [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)-Erstlieferung `tools/schema/schema.yaml`
ist als Folgepflicht in der ADR benannt und noch offen).

**Nicht in §3, zulässig:** Review-Report und Plan-Nachzug selbst (Lauf-Belege
pro Slice, Modul 5) · `docs/plan/adr/0043-*.md` + ADR-Index-Update ·
`harness/README.md`-Werkzeuge-Zeilen (Regel-7-Pflicht der neuen Targets) ·
`Dockerfile`-Kommentar-Edition (deps-Absatz, Anker-Hebung F-4).

## ADR-Konformität

- **[`ADR-0011`](../plan/adr/0011-persist-before-ack.md)** (Persist-before-ACK, Idempotenz): **konform und am realen Adapter belegt.** Persist mit Deduplizierung über die interne Transaktions-ID (`ON CONFLICT DO NOTHING` an `cdc.transaction`-PK und `UNIQUE (transaction_id, sequence)`), `TestPersistTransactionIsIdempotent` grün gegen die reale gepinnte PostgreSQL (eigener Lauf); die Pflicht steht am Port-Kontrakt und wird vom Adapter getragen.
- **[`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md)/[`SPEC-008`](../../spec/pflichtenheft.md)** (Fehlerklassen): **konform mit Ausnahme V-1.** Der Adapter trägt die Übersetzung über `storageFailure` → `outbound.ErrStorage` am Port (Sentinel am Port, `errors.Is`-Lesbarkeit, technische Ursache hinter der Klasse); beide Klassen-Tests grün, die rote Probe ist repliziert (Sentinel entfernt → beide rot, Beleg oben). **V-1:** der Commit-Ausgang des Store-Commits (`store.go:116` `return tx.Commit(ctx)`) meldet roh — die Klassen-Zusage der Kontrakt-Zeile und der Adapter-Kommentare trägt an dieser Stelle nicht.
- **[`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md)** (Kontrakt-Typen am Port): **konform.** `ChangeQuery`/`ChangeRecord` und die Sentinels (`ErrRangeInverted`, `ErrNonPositiveLimit`, `ErrStorage`) leben am Outbound-Port (`changestore.go`); der Re-Evaluierungs-Trigger („ein Port trägt einen Transport-Typ, den sein Use Case nicht aliasiert") ist nicht erfüllt; Fitness-Zeile durch `make a-check` (0 Befunde, `adapters → ports`-Kante) belegt.
- **[`ADR-0010`](../plan/adr/0010-postgresql-cdc-store.md)/[`ADR-0032`](../plan/adr/0032-postgresql-adapterdetail.md)** (PostgreSQL bleibt Adapterdetail; nativer Go-Stack): **konform.** pgx-Typen (`pgxpool`, `pgx.Rows`) bleiben im Adapter; der Port-Kontrakt kennt keinen Treibertyp; `CGO_ENABLED=0` im Dockerfile unverändert; pgx/v5 v5.11.0 ist rein Go.
- **[`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)** (neuer Accepted-Anker): **konform und konsistent mit dem Adapter-Code.** Die benannte Grenze trägt: die handgeschriebene DDL + `ApplySchema` bleiben die Quelle und der Testloader (Dockerfile-/Makefile-Kommentar und `schema.sql`-Kopf benennen sie im Indikativ als Bestand mit benanntem Ersatz); kein `tools/schema/schema.yaml` angelegt (Erstlieferung der Überführung steht aus — Folgepflicht der ADR, nicht dieses Slices); `make schema-validate` meldet die fehlende Quelle sauber (Exit 2, deklarierter Text — eigener Lauf); `schema-rollout` mit Pflicht-Report und Rollback-Artefakt; **kein Target hängt an `GATE_CHECKS`** (`grep`-Beleg: nur `baseline-verify` und `docs-check` akkumulieren); d-migrate-Image per Digest gepinnt.
- **[`ADR-0029`](../plan/adr/0029-domain-invarianten.md) Regel 3/6/7:** **konform** — `ErrTransactionNotCommitted`-Guard am Adapter (Regel 3), `UNIQUE (transaction_id, sequence)` + `CHECK (sequence >= 1)` in der DDL (Regel 6), `schema_version NOT NULL REFERENCES` (Regel 7).

## Fix-Commit `46fbf33` gegen die Review-Findings

| Finding | getragen im Fix? | Beleg |
|---|---|---|
| F-3 (Fehlerklasse ohne Träger) | **ja — mit Rest V-1** | `storageFailure` wickelt Treiber-Fehler in `outbound.ErrStorage` (New/Ping, Begin, Exec, Query, Scan, rows.Err); Sentinel am Port mit Kontrakt-Zeile; beide Klassen-Tests grün am realen Treiber; **rote Probe repliziert** (Sentinel entfernt → beide rot). Rest: der Commit-Ausgang (`tx.Commit`) meldet roh — V-1 |
| F-4 (superseded-Referenzen ohne Anker) | **ja** | `mapper.go:2-5` und `queries.go:2-5` zitieren jetzt „Paketstruktur je `ADR-0042`, der die Struktur-Regeln des abgelösten `ADR-0039` als Rest fortgilt"; Dockerfile-Build-Kommentar ebenso. Rest-Bestand `source.go:3` zitiert `ADR-0039` für die Paketstruktur (Fortgeltungs-Rest) — kein Befund (wie V-2 in verify-slice-003) |
| F-5 (kein Datei-Abschluss, Netz bleibt) | **ja** | `run-store-tests.sh` endet mit Zeilenumbruch (`0a`-Beleg); `cleanup()` entfernt Container **und** Netz (`docker network rm`), `trap cleanup EXIT` |
| F-6 (toter Grenz-Zweig) | **ja** | `uint64(commitPosition) > math.MaxInt64` ist entfernt; die Bereichs-Grenze liegt allein in `NewTransactionRow` |
| F-1 (Plan-Nachzug (b)/(d)) | **ja, über `13abe88`** | §1-Nachzug-Zeile mit Ergänzungs-/Änderungs-Benennung; §3-Zeilen; V-2 unten benennt die Träger-Hygiene des Commits |
| F-2 (image-hash veraltet) | **abgewiesen — Schiedsspruch oben** | deps-Layer-Änderung hebt den Digest nicht (kontrollierter Doppel-Build, byte-identisches Binary); Semantik mit `96c47af` präzisiert |
| F-7 (kein Compile-/Test-Gate) | **offen, korrekt bei Planner** | `GATE_CHECKS` trägt weiter nur `baseline-verify` + `docs-check` (+ `a-check`); Lerneintrag-Kandidat für die Closure §7 (Architect/Planner-Entscheidung) — kein DoD-Bruch, die DoD-Zeile „`make gates` grün" sagt nichts über Compile-Abdeckung |
| F-8 (leere Transaktion am Kontrakt unbenannt) | **notiert — Planner** | Test trägt die Grenze am realen Pfad (grün); die Kontrakt-Zeile nennt sie nicht — Kontrakt-Schärfung als Planner-Notiz, keine Nacharbeit am Diff |

## Befunde

### V-1 — Der Commit-Ausgang des Store-Commits trägt die Klasse `storage` nicht

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md) (Übersetzungsverantwortung des Adapters) · `SPEC-008` (Klasse `storage`) · Fix-Message `46fbf33` („Treiber-Fehler gehen an der Adapter-Grenze ueber storageFailure in die Klasse storage")
- `pfad`: `internal/adapters/driven/postgresstorage/store.go:116` (`return tx.Commit(ctx)`) gegen die Kontrakt-Zeile am Port (`changestore.go:101-104`: „Treiber-Fehler gehen an der Adapter-Grenze in die Klassen des Pflichtenhefts") und den Adapter-Kommentar (`store.go:64-65`: „Treiber-Fehler gehen in die Klasse `storage` (storageFailure)")
- `befund`: Die F-3-Nacharbeit wickelt jeden Treiber-Ausgang in `storageFailure` — außer dem letzten: der Commit-Ausgang des Store-Commits liefert Treiber-Fehler (Verbindungsabbruch während COMMIT) **roh** an den Port. Die Klassen-Tests decken den Exec-Pfad (FK-Verstoß) und den Ping-Pfad; der Commit-Pfad ist ungetestet und trägt die behauptete Klasse nicht. Dieselbe Klasse wie review-slice-003 F-4 („unbenannte Grenze am öffentlichen Vertrag"): für den Capture-Pfad folgenlos (kein ACK bei Fehler), für die Verbräucher der Klassen (slice-006-Verdrahtung, `cdc_errors_total{class}`) ist der Pfad unbenannt.
- `verifizierbar`: ja — Code-Lese-Beleg `store.go:116` gegen die Wrapp-Liste (`New`, `Begin`, `Exec` ×2, `Query`, `Scan`, `rows.Err` — alle gewickelt, Commit nicht); `grep -n "tx.Commit" internal/adapters/driven/postgresstorage/store.go`
- `klasse`: F-3-Nacharbeit unvollständig am letzten Treiber-Ausgang (unbenannte Grenze am Port-Kontrakt, 2. Auftreten der Klasse)

### V-2 — Review-Basis `193716e` ist unerreichbar; `13abe88` ist ein verbundener Commit unter einer Message, die nur Tests nennt

- `kategorie`: MEDIUM
- `quelle`: Traceability-Regeln (`harness/README.md`) · Modul 8 (Übergabe-Artefakte) · Modul 9 (Commit-Messages beschreiben, was da ist)
- `pfad`: `docs/reviews/review-slice-004.md` („Range `c9a068f..193716e`, 3 Commits") gegen `git cat-file -t 193716e` (existiert) und `git log --all | grep 193716e` (**0 Treffer** — von keinem Ref erreichbar); `git show 13abe88 --name-only`
- `befund`: Der vom Review dokumentierte Range-Head `193716e` war ein Vorgänger-Zustand von `13abe88` (gleiche Code-Inhalte — `git diff 193716e..13abe88` berührt nur die sechs Doku-Dateien) und wurde beim Neuschreiben verdrängt: `13abe88` trägt unter der Message „test(storage): Adapter-Tests gegen reale PostgreSQL im Testcontainer" zusätzlich **Review-Report (429 Zeilen) · Plan-Nachzug §1/§3 · [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) (Accepted) + ADR-Index · `harness/README`-Werkzeuge-Zeilen · DoD-Edits an slice-005/006** — keines davon nennt die Message. Der Review-Beleg verweist damit auf einen Commit, den kein Ref trägt; die Review-Kette „Diff → Report" ist formal unterbrochen (inhaltlich gedeckt: der Code von `193716e` ist byte-identisch in `13abe88`). Ein Accepted-ADR im Implementer-Test-Commit berührt zudem die Rollen-Form (Modul 8: ADR schreibt der Architect) — die ADR trägt den Auftraggeber-Entscheid als Autor pt9912; Form-Befund, kein Inhaltswiderspruch.
- `verifizierbar`: ja — `git cat-file -t 193716e`; `git log --all --oneline | grep 193716e` (0 Treffer); `git diff --stat 193716e 13abe88`
- `klasse`: Beleg-Kette umgeschrieben (unerreichbare Review-Basis) · verbundener Commit mit unbenannten Artefakten

### V-3 — DoD-Aussage „kein öffentlicher Vertrag berührt (nur internal/**)" deckt die eigene Range nicht

- `kategorie`: MEDIUM
- `quelle`: DoD-Item 5 (Slice-Plan §2) · Minimal Agent Workflow Schritt 7 (`harness/README.md`) · `harness/README.md` als Rang-9-Quelle
- `pfad`: Plan §2 Item 5 („**entfällt**: kein öffentlicher Vertrag berührt (nur `internal/**`)", abgehakt in `13abe88`) gegen `git show 13abe88 -- harness/README.md` (+2 Werkzeuge-Zeilen) und `git show 59e3568 -- harness/README.md` (+2 Werkzeuge-Zeilen, `make schema-validate`/`schema-rollout`)
- `befund`: Die eigene Range berührt `harness/README.md` — den öffentlichen Harness-Einstieg (Rang 9) — an der Werkzeuge-Tabelle. Das ist genau die Schritt-7-Pflicht (neue Targets deklarieren) und korrekt getragen; falsch ist die DoD-Begründung: sie behauptet die Nicht-Berührung, die die eigene Range widerlegt. Das Doku-Update ist **geschehen**, nicht „entfallen" — der Häkchen-Text ist falsch, nicht das Verhalten.
- `verifizierbar`: ja — die beiden `git show`-Belege gegen die Plan-Zeile
- `klasse`: DoD-Behauptung ohne Deckung durch die eigene Range

### V-4 — `59e3568` trägt unbenannte Änderungen am Minimal-agent-workflow und bricht den Datei-Abschluss

- `kategorie`: LOW
- `quelle`: AGENTS §3.7 (Kommentar/Zustandsfeld trägt den Zustand, nicht die Chronik) · Traceability (Message deckt den Inhalt) · Datei-Abschluss-Klasse (3. Auftreten: review-slice-001 F-6, review-slice-004 F-5, jetzt hier)
- `pfad`: `git show 59e3568 -- harness/README.md` (Hunk 2: Schritte 2–4 und 8 des Minimal Agent Workflow umformuliert, Rollenwechsel-Absatz umgestellt) · `tail -c1 harness/README.md` → kein Zeilenumbruch · Commit-Message („harness/README.md: zwei Werkzeuge-Zeilen")
- `befund`: Die Message nennt zwei Werkzeuge-Zeilen; der Commit trägt zusätzlich eine Umschrift des Minimal-agent-workflow (Inhalte deckungsgleich mit `AGENTS.md` §6 — Alignment, kein Widerspruch) und lässt die Datei ohne Zeilenumbruch enden. Die Klasse „Datei-Abschluss" steht damit beim **dritten** Auftreten — Register-Kandidat für die Closure (neue Beobachtung, 1×).
- `verifizierbar`: ja — `git show`-Hunks; `tail -c1 | xxd`
- `klasse`: unbenannte Commit-Inhalte · Datei-Abschluss (3. Auftreten)

### V-5 — DoD-Items fremder Slices vorab abgehakt

- `kategorie`: LOW
- `quelle`: Modul 5 (Closure-Kriterien sind Bedingungen für `done/`) · AGENTS §3.7 (Zustandsfelder)
- `pfad`: `13abe88` — `docs/plan/planning/open/slice-005-replication-stream-adapter.md` und `slice-006-integrationstest-umgebung.md`: Item „Doku-Update" von der Vorlagen-Zeile auf „**entfällt** … [x]" gesetzt, **vor** deren Arbeit
- `befund`: Ein `[x]` auf einem DoD-Item eines Slices in `open/` ist eine Abschluss-Aussage über Arbeit, die nicht stattgefunden hat. Die Entfällt-Vorhersage kann sich als Plan-Änderung erweisen (beide Slices berühren Compose — ein öffentlicher Vertrag ist nicht ausgeschlossen), und der Häkchen-Stand lügt bis dahin über den Zustand. Korrektur wirkt vorwärts: bei der Eröffnung von slice-005/006 ist das Item zurückzusetzen oder der Entfällt-Grund je Zeitpunkts zu tragen.
- `verifizierbar`: ja — `git show 13abe88 -- docs/plan/planning/open/`
- `klasse`: Vorschul-Zustandsaussage über fremde Slices

### V-6 — d-migrate-Verdrahtung ist nicht in §3 budgetiert (und die [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)-Committierung trägt der Implementer-Commit)

- `kategorie`: INFO
- `quelle`: Modul 5 („Wer später mitnimmt …, hat den Plan geändert") · Modul 8 (ADR-Authorenschaft)
- `pfad`: `59e3568` (Makefile-Targets `schema-validate`/`schema-rollout` + README-Zeilen) gegen Slice-Plan §3 (keine Zeile zu d-migrate/Schema-Targets) · [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) erstmals committet in `13abe88`
- `befund`: Die d-migrate-Verdrahtung liegt innerhalb der Lebensdauer dieses Slices, aber außerhalb seines §3 — entweder ist sie eine §3-Plan-Erweiterung (F-1-Klasse, 5. Auftreten) oder ein eigenes wellenloses Slice (das [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)-Einbau-Artefakt). Die ADR selbst ist inhaltlich konform (oben) und trägt die Grenze konsistent mit dem Adapter-Code; die Committierung im Implementer-Commit ist Form, kein Inhalt.
- `verifizierbar`: ja — §3 gegen die Makefile-Targets; Commit-Geschichte der ADR-Datei
- `klasse`: Plan-Erweiterung ohne Plan-Nachzug (5. Auftreten der Klasse — Sequenz läuft weiter)

## Negativbefunde

- geprüft, ohne Befund: **Docker-only (AGENTS §3.1)** — alle go-Belege im gepinnten Toolchain-Container (`--network none` für `make test`, Modul-Cache im Volume), alle Image-Builds über `docker buildx build` aus `/tmp`-Worktrees; kein Host-Toolchain-Aufruf, kein Schreibzugriff auf den Hauptbaum (`git status` clean am Ende; `git worktree list` trägt nur den Hauptbaum)
- geprüft, ohne Befund: **Idempotenz am realen Adapter** — `TestPersistTransactionIsIdempotent` grün gegen die gepinnte PostgreSQL (eigener Lauf); Deduplizierung an beiden Kanten (`cdc.transaction`-PK, `UNIQUE (transaction_id, sequence)`)
- geprüft, ohne Befund: **Lesen verändert keine Position** — `TestReadLeavesPersistedStateUnchanged` grün am realen Treiber (Zeilen- und Transaktions-Zählung unverändert)
- geprüft, ohne Befund: **SQL-Injection/Credential-Grenze** — alle vier Queries parametrisiert (`$1…$5`), Tabellen-/Spaltennamen als Konstanten in `queries/`; der Adapter loggt nichts; pgx v5.11.0 redaktiert Passwörter in Parse-Fehlern (Review-Probe, nicht widersprochen)
- geprüft, ohne Befund: **Testcontainer-Hygiene** — Daten im Container, `/src` read-only gemountet, Netz-Rückbau im `trap` (F-5-Fix), kein Arbeitsbaum-Eintrag; nach meinen Läufen: PG-Container und Netz entfernt
- geprüft, ohne Befund: **§1-Abgrenzung** — keine Replication-Stream-Integration, keine Consumer-Use-Cases, kein Performance-Tuning in der Range; die DDL-Schmälerung (d) ist plan-nachgezogen und code-konsistent (Vorbemerkung `schema.sql:10-13`)
- geprüft, ohne Befund: **Traceability der Range** — `make doc-commits`/`make doc-immutable` über `c9a068f..96c47af` je 0 Befunde; jeder Commit trägt `LH-*`/`ADR-*`; keine `SPEC-*`/`ARC-*` in Commit-Messagen
- geprüft, ohne Befund: **a-check-Abdeckung** — der `adapters`-Glob (`internal/adapters/**`) matcht den echten Adapter-Content; `make gates` (a-check im Bündel) 0 Befunde; keine undeclarierte Kante
- geprüft, ohne Befund: **V-3-Rest aus verify-slice-003** — der Vorgänger-Worktree `/tmp/verify-slice-002-wt` ist nicht mehr registriert (`git worktree prune` in diesem Lauf; `git worktree list` clean)
- geprüft, ohne Befund: **Digest-Pins** — Toolchain- und PG-Digests konsistent zwischen Dockerfile, Makefile und `run-store-tests.sh`; d-migrate-Image per Digest gepinnt, kein Target an `GATE_CHECKS`

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 2 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** F-3-Nacharbeit unvollständig am letzten
Treiber-Ausgang · Beleg-Kette umgeschrieben (unerreichbare Review-Basis,
verbundener Commit) · DoD-Behauptung ohne Deckung durch die eigene Range ·
unbenannte Commit-Inhalte + Datei-Abschluss (3. Auftreten) ·
Vorschul-Zustandsaussage über fremde Slices · Plan-Erweiterung ohne
Plan-Nachzug (5. Auftreten — Sequenz läuft weiter)

**Zusammenfassung DoD:** **4/10 Punkte jetzt erfüllt** (Items 1–4, Belege
selbst gefahren, inkl. 10/10 Store-Tests gegen die reale gepinnte PostgreSQL)
· 1 nicht bestätigt in der Begründung (Item 5 — die Arbeit ist korrekt getan,
die Begründung falsch, V-3) · 1 entfällt (Item 7) · 3 erfüllen sich erst bei
Closure (Items 6, 8, 9-Endwert) · 1 delegiert an die Welle-2-Closure
(Item 10). **F-2-Schiedsspruch:** Review-Fassung abgewiesen (Belege: Doppel-
Build A/B, Binary-Identität, Provenance-Gegenprobe); Semantik mit `96c47af`
getragen; Residuum builder-gebundener Digest, nicht blockierend.

**Review-Blockierpunkte:** Blocker 2 (F-1 Plan-Nachzug) erfüllt (`13abe88`).
Blocker 3 (F-3) erfüllt mit Rest V-1 (Commit-Ausgang). Blocker 1 (F-2)
abgewiesen in der Review-Fassung; die Semantik-Frage ist mit `96c47af`
beantwortet und durch meine Builds bestätigt — als Residuum bleibt die
builder-Gebundenheit des Digests benannt (nicht blockierend).

## Verdikt

**Merge-blockierend:** nein — der Adapter ist am realen gepinnten Treiber
belegt (10/10 Store-Tests, eigener Lauf), die kritischen Zusagen (Idempotenz,
Bereichs-Ordnung, Lesen-ohne-Positionsänderung) sind grün und die
Klassen-Zusagen durch die replizierte rote Sentinel-Probe gedeckt; `make
gates` ist grün am HEAD (zweifach gefahren), `go vet`/`gofmt` clean im
gepinnten Container, die drei Paarungen-Vorarbeiten und Plan-Nachzüge stehen.

**Blockierend für Closure (der normale Zustand eines Slice in
`in-progress/`):**

1. **V-1:** der Commit-Ausgang des Store-Commits trägt die Klasse `storage`
   nicht — der Adapter wickelt `tx.Commit` in `storageFailure` (oder der Plan
   terminisiert die Grenze mit Ausgang), bevor der F-3-Fix als vollständig
   gemeldet wird.
2. **V-3:** die DoD-Begründung zu Item 5 ist zu korrigieren — das Doku-Update
   ist geschehen (korrekt), das Item trägt „entfällt" mit einer falschen
   Begründung; vor der Closure ist die Aussage zu berichtigen (Plan-Edit,
   keine Code-Nacharbeit).
3. **Closure-Pflichten (Items 6, 8, 9):** Closure-Notiz mit Lerneintrag
   (inkl. F-7-Klasse als Lerneintrag-Kandidat und den Finding-Klassen
   V-1…V-6), Register-Notiz „kein neues Auftreten" (plus die V-4-Klasse als
   neuen Register-Kandidaten), Endausgänge für beide §6-Risiken.
4. **V-4/V-5/V-6 als vor-Closure-Nacharbeit ohne Blockier-Charakter:**
   Zeilenumbruch am Ende von `harness/README.md`, Vorschul-`[x]` bei
   slice-005/006 zurücknehmen oder begründen, d-migrate-Verdrahtung als
   §3-Erweiterung oder eigenes wellenloses Slice adressieren (Planner).

**Übergabe:** Bericht an den Planner. Keine Reparaturen. Der F-2-Schiedsspruch
ist mit der `96c47af`-Semantik abgeschlossen; mein Residuum (builder-gebundener
Digest, Abschnitt oben) gehört als Hinweis in den Architect-Kontext, falls die
„genau dann"-Formulierung künftig als maschinelle Prüfung verwendet werden
soll — als Beleg-Semantik trägt sie.