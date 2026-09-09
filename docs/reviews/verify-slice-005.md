# Verifier-Report: slice-005 — 2026-09-09

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of Done,
**10 Items** — das doppelte „make gates"-Item ist in `1ac3558` entfernt, der
Auftrag nannte 11), §3 (Plan-vs-Code, Range `cb021dc..HEAD` inkl.
Plan-Nachzug `1ac3558`), §6 (Risiko-Ausgänge) und ADR-Konformität
([`ADR-0006`](../plan/adr/0006-replication-stream-driving-adapter.md) ·
[`ADR-0007`](../plan/adr/0007-source-ack-outbound-port.md) ·
[`ADR-0008`](../plan/adr/0008-pgoutput-standard.md) ·
[`ADR-0011`](../plan/adr/0011-persist-before-ack.md) ·
[`ADR-0012`](../plan/adr/0012-at-least-once.md) ·
[`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md) ·
[`ADR-0026`](../plan/adr/0026-composition-root.md) ·
[`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)). Nicht geprüft:
Diff gegen Plan/Hard Rules im Detail (Reviewer, `review-slice-005.md`,
Verdikt dort), realer Bedarf (Validator).

**Gegenstand:** Implementer-Handoff zu
`docs/plan/planning/in-progress/slice-005-replication-stream-adapter.md` ·
Range `cb021dc..HEAD` (HEAD **`1ac3558`**) · Fix-Commits `8c5484f`
(F-8/F-10/F-14), `0bc3277` (F-3/F-4/F-9 + tabellenscopierte Publications),
`11d9472` (F-6), `7b3a52b` + `3140568` (F-7 → [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)), Plan-Nachzug
`1ac3558` (F-2/F-11).

**Grundsatz:** Es wurden **keine Behauptungen übernommen** — jeder Sensor
unten wurde in diesem Lauf selbst gefahren; Ausgaben sind Belege über stdout.
Der Hauptbaum blieb read-only (`git status` clean am Ende, `git worktree
list` trägt nur den Hauptbaum); Image-Builds und Mutations-Proben liefen in
`/tmp`-Worktrees (nach Abschluss entfernt), die reale Testcontainer-
PostgreSQL lief in einem eigenen Container/Netz (`wal_level=logical`,
`wal_sender_timeout=2s` belegt; nach Abschluss abgeräumt, Temp-Images
`tmp005:7b`/`tmp005:head` entfernt).

**Modell:** Claude Code (glm-5.3-flash) · **Datum:** 2026-09-09

**Eingangs-Kontext:**

- Slice-Plan §1–§8 (am HEAD `1ac3558`, inkl. Plan-Nachzug) ·
  `docs/reviews/review-slice-005.md` (F-1…F-16) · Fix-Commits im Volltext ·
  `docs/reviews/verify-slice-004.md` (Gerüst + F-2-Schiedsspruch als
  [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)-Basis)
- `spec/lastenheft.md` ([`LH-FA-CAP-001`](../../spec/lastenheft.md)…003, 004,
  [`LH-FA-CAP-006.a`](../../spec/pflichtenheft.md), [`LH-QA-REL-001`](../../spec/lastenheft.md)) ·
  `spec/pflichtenheft.md` (`SPEC-008`, [`SPEC-010`](../../spec/pflichtenheft.md), [`LH-FA-CFG-001.a`](../../spec/pflichtenheft.md)) ·
  `spec/architecture.md` (`ARC-005`/006/008)
- `harness/README.md` (Werkzeuge-Zeilen `test-replication`, `make image`) ·
  `Makefile`, `tools/harness/run-replication-tests.sh`, `Dockerfile`,
  `.dockerignore`, `harness/image-hash.txt` ·
  Beobachtungs-Register (`BEO-PGC/a-check-null-abdeckung`)

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` (HEAD `1ac3558`) | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 94 Datei(en) geprüft, 0 Befund(e)` · `a-check … gesamt: 0 Befund(e)` | 0 |
| `go build ./...` + `go vet ./...` + `gofmt -l .` im gepinnten Toolchain-Container (`golang:1.27-alpine@sha256:cf6fca…`, Modul-Cache-Volume `pg-change-feed-gomodcache`, `/src:ro`, `--network none`) | `BUILD_OK` · `VET_OK` · `GOFMT_CLEAN` | 0 |
| `make test` (netzlos, gepinnter Container) | **10 Pakete `ok`** (`postgresack`, `postgresstorage`, `postgresstorage/mapper`, `replication/decode`, `replication/mapper`, `replication/receive`, `port/outbound`, `usecase/capture`, `bootstrap`, `domain/model`) — Implementer-Claim „10/10" bestätigt | 0 |
| `make test-replication` (reale PostgreSQL `postgres:18-alpine@sha256:63bdc97d…`, `wal_level=logical`, gepinnt, Daten im Container) | **grün** — dieselben 10 Pakete, `replication/receive 4.0s` und `bootstrap 8.4s` (reale Läufe) | 0 |
| Eigene PG-Instanz, verbose | `pg_isready` ok, `SHOW wal_level` = `logical`, `SHOW wal_sender_timeout` = `2s`; **5/5 Stream-Tests + 3/3 ACK-Tests + Verdrahtungs-Test PASS**: `TestStreamTranslatesRealChanges`, `TestStreamOpenTransactionNotConsumable`, `TestStreamTruncateUnsupported`, **`TestStreamKeepaliveReportsAcknowledgedPosition`**, **`TestStreamRestartsOnExistingSlot`**, `TestNewRequiresConnection`, `TestAcknowledgeRejectsZeroPosition`, `TestAcknowledgeOnClosedConnection`, `TestRealPersistBeforeAck` (8.4 s, echte Persist-before-ACK-Verdrahtung) | 0 |
| Mutations-Probe 1 (Keepalive: `WALWritePosition/Flush/Apply: s.lastAcked` → `0`; `/tmp`-Worktree, revertiert und entfernt) | `--- FAIL: TestStreamKeepaliveReportsAcknowledgedPosition` — `confirmed_flush_lsn erreicht die bestätigte Position nicht (1d162b0 != 1d163d0)` — **erwartet rot, rot gesehen** (Review-F-3-Variante mit der alten Code-Basis war grün; die Zusage trägt jetzt einen Test) | 1 |
| Mutations-Probe 2 (Doppel-BEGIN-Guard `mapper.go:102-104` entfernt; `/tmp`-Worktree, revertiert und entfernt) | `--- FAIL: TestConsumeBeginWithoutCommit` — **erwartet rot, rot gesehen** | 1 |
| Image-Builds aus `/tmp`-Worktrees (`docker buildx build --load --metadata-file`, Hauptbaum unberührt) | `7b3a52b` → **`sha256:9ac4a9fb812aa57c14df104843517091075a8d8afdf24c688f489df68fbc2344`** (exakt der eingetragene Beleg) · HEAD `1ac3558` → **dasselbe** · Binary-Extraktion beider Images: `sha256 43c3aec0…` byte-identisch | 0 |
| `make doc-commits RANGE=cb021dc..HEAD` | `d-check: 94 Datei(en) geprüft, 0 Befund(e)` (Traceability je Commit, volle Range) | 0 |
| `make doc-immutable RANGE=cb021dc..HEAD` | `d-check: 94 Datei(en) geprüft, 0 Befund(e)` (Immutabilität) | 0 |
| `git worktree list` / `git status --porcelain` | nur der Hauptbaum; Arbeitsbaum clean — keine Einträge dieses Laufs | — |

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt — 10 Items nach Entfernung des Duplikats in `1ac3558`)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | Stream-Adapter übersetzt `pgoutput` in Aufrufe des `CaptureInboundPort` (INSERT/UPDATE/DELETE) — Teil-Beleg zu [`LH-FA-CAP-001`](../../spec/lastenheft.md)…003 | **bestätigt** | `TestStreamTranslatesRealChanges` grün am realen Pfad (eigener Lauf, verbose): INSERT mit exaktem Image `{"id":"1","name":"Alpha"}`, UPDATE, DELETE mit Alt-Image `{"id":"1"}`, Transaktion der nicht aktivierten Tabelle als leere Transaktion, Commit-Positionen in Reihenfolge; Unit-Seite (`decode_test.go`, `mapper_test.go`) grün gegen die Binärcodes |
| 2 | ACK real: nur Positionen nach dauerhafter Persistenz — [`LH-QA-REL-001.a`](../../spec/pflichtenheft.md) am realen Treiber | **bestätigt** | `TestRealPersistBeforeAck` grün am echten Capture Service (persistiert über den ChangeStore, bestätigt über den ACK-Adapter — 8.4 s an der realen Instanz); Keepalive-Regel am realen Pfad (`TestStreamKeepaliveReportsAcknowledgedPosition` + rote Probe, unten) |
| 3 | `make gates` grün | **bestätigt** | Drei Gates grün am HEAD (Tabelle oben); `[x]` am Item in `1ac3558` korrekt |
| 4 | Review durchgeführt, Report unter `docs/reviews/` | **bestätigt** | `docs/reviews/review-slice-005.md` committet (in `22e9b93`, im Range); Rollenwechsel nach Schritt 8 eingehalten |
| 5 | Doku-Update bzw. begründete Aussage „kein öffentlicher Vertrag berührt" | **getan, Häkchen noch offen** | `harness/README.md` trägt die Werkzeuge-Zeile `test-replication` (Zeile 127, [`ADR-0030`](../plan/adr/0030-testpyramide.md)) und die `make image`-Zeile in der [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)-Semantik (Zeile 123) — **getragen**; das Häkchen steht noch nicht (fällig vor dem `git mv`) |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **offen bis Closure (notiert)** | §7 trägt Platzhalter; Slice korrekt in `in-progress/` — fällig vor dem `git mv` (Lerneintrag-Kandidaten unten: Datei-Abschluss-Klasse ≥3×, F-1/F-5-Sequenz) |
| 7 | Reconciliation-Register | **entfällt** | `docs/plan/planning/reconciliation.md` existiert nicht (Repo ohne Brownfield-Bootstrap) — vom DoD-Wortlaut ausdrücklich vorgesehen; kein Inventur-Fund im Range |
| 8 | Beobachtungs-Register fortgeschrieben | **kein neues Auftreten — notiert (Closure-Antwort fällig)** | `BEO-PGC/a-check-null-abdeckung` trägt den Ausgang `verkörpert · seit welle-1`; `evidence/` trägt `slice-001/002/003.md` — unverändert; **kein neues Auftreten dieser Klasse** (a-check 0 Befunde am HEAD). §8 des Plans notiert die Sichtung korrekt („keine offenen Treffer"). Neue Register-Kandidaten für die Closure: die Datei-Abschluss-Klasse (V-1, jetzt ≥3×) und ggf. die Commit-Kennung-Klasse |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **Voreinstellung getragen — Endwert fällig bei Closure** | (a) „Fake-Ordnung stärker als realer Treiber" — bewertet am realen Treiber: die Ordnungszusagen (`TestStreamOpenTransactionNotConsumable`, Commit-Reihenfolge) sind am realen Pfad belegt; Endausgang fällig. (b) „Reale Idempotenz" — ([`ADR-0011`](../plan/adr/0011-persist-before-ack.md))-Kontrakt, belegt im Integrationstest (slice-006) — Endausgang fällig. **Zusätzlich:** die Review-Risiken (F-3 Keepalive, F-4 Restart) sind durch `0bc3277` **mit Tests am realen Pfad belegt** — bei der Closure als entstanden/eingetreten mit Ausgang zu führen |
| 10 | Drei Paarungen (Anker · Folge-Slice · Register) | **delegiert — korrekt** | Repo **mit** Wellen-Betrieb (`docs/plan/planning/welle-2.md` flach vorhanden, Slice trägt `Welle: welle-2`): der DoD-Wortlaut weist die Prüfung der Welle-2-Closure zu; hier nicht fällig |

**Zwischenstand: 4/10 jetzt erfüllt** (Items 1–4, Belege selbst gefahren) ·
1 getan, Häkchen offen (Item 5) · 1 entfällt (Item 7) · 3 erfüllen sich erst
bei Closure (Items 6, 8 als Notiz, 9 Endwert) · 1 delegiert (Item 10). Der
Zustand ist der **normale eines Slice in `in-progress/`** — die offenen Items
sind Closure-Pflichten, keine Befunde gegen die Implementer-Arbeit.

## Plan-vs-Code-Diff (Range `cb021dc..HEAD`, Plan-Nachzug `1ac3558`)

**Plan-Nachzug (`1ac3558`) gegen den Review-Blocker F-2:** bestätigt — §3
trägt jetzt alle drei nachgezogenen Entscheidungen: Keepalive-Regel (mit
[`LH-QA-REL-001.a`](../../spec/pflichtenheft.md) und roter Probe als Testpflicht),
`BindCapture`-Trennung ([`ADR-0007`](../plan/adr/0007-source-ack-outbound-port.md) Option C), Verlagerung der
Verdrahtungstests nach `internal/bootstrap` ([`ADR-0026`](../plan/adr/0026-composition-root.md)), plus
`Makefile`/`test-replication`-Verkabelung und `pglogrepl`-Deps ([`ADR-0032`](../plan/adr/0032-postgresql-adapterdetail.md)).
Auch das Dup-DoD-Item (F-11) ist entfernt und das getragene abgehakt. Die
Message trägt Vertrags-Kennungen ([`ADR-0007`](../plan/adr/0007-source-ack-outbound-port.md)/0026/0032, `LH-QA-REL`).

**Deckung §3 (am HEAD):** alle gelieferten Code-Dateien sind durch die
(nachgezogenen) §3-Zeilen gedeckt: `replication/{receive,decode,mapper}/*.go`
inkl. drei Testdateien, `postgresack/*.go` (inkl. `ack_test.go`),
`internal/bootstrap/replication_stream_test.go`, `Makefile`-Block +
`tools/harness/run-replication-tests.sh`, `go.mod`/`go.sum`
(pglogrepl v0.0.0-20260824, pgx/v5 v5.11.0 bleibt). Der Keepalive-Regel- und
BindCapture-Nachzug deckt `receive.go` inhaltlich.

**Unbudgetiert — V-2:** `internal/application/port/outbound/replicationack.go`
(neu, 32 Zeilen: Sentinel + `ReplicationAckPort`) ist in **keiner** §3-Zeile
aufgetaucht — der Review-F-2-Befund nannte den Port-Kontrakt ausdrücklich,
der Nachzug hat ihn nicht aufgenommen. Klein (Contract-Belegung des
[`ADR-0007`](../plan/adr/0007-source-ack-outbound-port.md)-Ports, keine
Eigenentscheidung), aber der Nachzug ist damit **unvollständig** — Rest der
F-2-Klasse.

**Nicht in §3, zulässig (Lauf-Belege bzw. Planner-/Decision-Artefakte):**
Review-Report `docs/reviews/review-slice-005.md` (Übergabe-Artefakt, Modul 5)
· [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) + ADR-Index-Zeile
(Entscheidung, `7b3a52b`/`3140568`) · `Dockerfile`/`harness/README.md`
Verkabelungsträger der [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) · Plan-Nachzug und `31e1b41`-Berührungen an
slice-006 (Planner) · `spec/pflichtenheft.md` ([`SPEC-015`](../../spec/pflichtenheft.md), Planner-Commits
vorangehend).

## ADR-Konformität

- **[`ADR-0006`](../plan/adr/0006-replication-stream-driving-adapter.md)** (Stream als Driving; „dekodieren ja, entscheiden nein"): **konform** — `receive.go` hält keinen Ordnungs- und keinen Persistenz-Entscheid: die bestätigte Position kommt ausschließlich als `CaptureResult.Acknowledged` entgegen (`receive.go:348-350`), der Stream bestätigt nichts selbst; Keepalive meldet nur `lastAcked`.
- **[`ADR-0007`](../plan/adr/0007-source-ack-outbound-port.md)** (ACK real, Option C, BindCapture-Split): **konform** — `PostgresReplicationAckAdapter` als Driven-Rolle an derselben technischen Verbindung (`stream.Conn()`), getrennte Pakete und Sentinels (`receive.ErrReplication` vs. `outbound.ErrReplication`), Verdrahtung in der Composition-Root (`replication_stream_test.go:123`, `BindCapture(capture.NewCaptureService(store, ack))`); der ACK-Adapter setzt Write/Flush/Apply auf die bestätigte Position (`ack.go:59-63`).
- **[`ADR-0008`](../plan/adr/0008-pgoutput-standard.md)** (`pgoutput` als Standard-Output-Plugin): **konform** — `outputPlugin`-Konstante, `proto_version '1'`, `publication_names` (`receive.go:113-118`).
- **[`ADR-0011`](../plan/adr/0011-persist-before-ack.md)** (Deduplizierung über Transaktions-ID): **konform in der Ordnung** — die Kennung läuft über den BEGIN-Stand des Assemblers (`mapper.go:105`); Restart liest `confirmed_flush_lsn` als Fortsetzungs-Stand (`receive.go:203-218`), belegt durch `TestStreamRestartsOnExistingSlot`. **Rest (V-3):** die Wraparound-Grenze (XID-Wiederholung nach 2^32) ist **weiterhin unbenannt** — kein Eintrag in [`ADR-0011`](../plan/adr/0011-persist-before-ack.md), kein Code-Kommentar, kein Plan-§6-Ausgang (`grep -i wraparound` → 0 Treffer; F-13 des Reviews läuft weiter).
- **[`ADR-0012`](../plan/adr/0012-at-least-once.md)** (Restart trägt): **konform und belegt** — `TestStreamRestartsOnExistingSlot` übt den bestehende-Slot-Zweig am realen Pfad (Stream-Ende → weitere Changes → Restart über denselben Slot; Wiederholung bestätigter Transaktionen als At-Least-Once-Fall akzeptiert). Review-F-4 ist damit **beantwortet**.
- **[`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md)** (Klassen je Port-Kontrakt, Sentinel-Träger): **konform** — `decode.ErrSchema`, `receive.ErrReplication`/`ErrConfiguration`, Mapper-Sentinels (`ErrTruncateUnsupported` `schema`, `ErrChangeWithoutBegin`/`ErrCommitWithoutBegin`/`ErrBeginWithoutCommit` `replication`), `outbound.ErrReplication` am ACK-Port; `replicationFailure` wickelt Treiber-Fehler, `TestAcknowledgeOnClosedConnection` trägt die Übersetzung; die neue Doppel-BEGIN- und leeres-CopyData-Grenze enden im Fehlerklassen-Pfad (`receive.go:278-279`).
- **[`ADR-0026`](../plan/adr/0026-composition-root.md)** (Verdrahtungstests im bootstrap-Layer): **konform** — `internal/bootstrap/replication_stream_test.go`; a-check 0 Befunde im `make gates`-Lauf.
- **[`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)** (Image-Beleg-Semantik): **konform, eingetragen und durch meine Builds bestätigt** — Dockerfile-Kopf (Zeilen 5-11) und Werkzeuge-Zeile (`harness/README.md:123`) tragen die korrigierte Semantik; ADR-Index-Zeile vorhanden; die `96c47af`-Formel ist dort ausdrücklich verworfen. **Eigener Reproduktions-Beleg:** zwei `buildx`-Builds (7b3a52b und HEAD) lieferten exakt den eingetragenen `sha256:9ac4a9fb…`; beide exportierten Binaries sind byte-identisch (`43c3aec0…`). **Erklärung (nicht in [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md), aber konsistent):** `cmd/pg-change-feed` ist noch der Bootstrap-Stub (`import "fmt"/"os"` — main.go:7-10), das Binary ist damit invariant gegen die Adapter-Änderungen; kein Stale-Beleg, **kein Digest-Commit fällig** (Binary unverändert). Notiz für den Planner: der Image-Beleg bezeugt weiterhin den Bootstrap-Binärstand, nicht die slice-005-Funktionalität — erst die Verdrahtung in `cmd` (slice-006-Sicht) bringt sie in das Image.

## Fix-Commits gegen die Review-Findings

| Finding | getragen im Fix? | Beleg |
|---|---|---|
| F-3 (Keepalive ohne Test-Träger) | **ja** | `TestStreamKeepaliveReportsAcknowledgedPosition` am realen Pfad; Mutation `lastAcked` → `0` **rot** (eigener Lauf); `confirmed_flush_lsn` wird gegen `ackedPosition` gehalten, nie darüber hinaus — die Zusage „nie über die bestätigte Position hinaus" trägt ein Test am realen Pfad (Prüfpunkt a) |
| F-4 (Restart-Zweig ohne Beleg) | **ja** | `TestStreamRestartsOnExistingSlot` übt `ensureSlot` am bestehenden Slot (liest `confirmed_flush_lsn`) und akzeptiert die Wiederholung bestätigter Transaktionen |
| F-6 (Datei-Abschluss go.mod/go.sum) | **ja, mit Rest V-1** | `go.mod`/`go.sum` enden auf `\n` (`0a`-Beleg), Makefile und Test-Skript sauber; **aber** `harness/README.md`, [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) und `Dockerfile` enden ohne Umbruch — V-1 |
| F-8 (Doppel-BEGIN still) | **ja** | Sentinel `ErrBeginWithoutCommit` (`mapper.go:102-104`), `TestConsumeBeginWithoutCommit`; Guard-Mutation **rot** (eigener Lauf) — Prüfpunkt (b) |
| F-9 (postgresack ohne Tests) | **ja** | `ack_test.go` neu (90 Zeilen): `TestNewRequiresConnection`, `TestAcknowledgeRejectsZeroPosition`, `TestAcknowledgeOnClosedConnection` — alle grün am realen Treiber |
| F-10 (totes Feld `Change.XID`) | **ja** | Feld entfernt; `XID` lebt nur noch an `decode.Begin` (`decode.go:37,133`) — konsistent mit der Assembler-Zuordnung |
| F-14 (leeres CopyData → Panik) | **ja** | `receive.go:278-279` — Längenprüfung, Ausgang in `ErrReplication`; kein eigener Test (INFO-Klasse, akzeptabel) |
| F-7 (Beleg-Semantik-Formel widerlegt) | **ja, über [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)** | Decision-ADR Accepted, Verkabelung getragen (Dockerfile-Kopf, Werkzeuge-Zeile), Reproduktions-Claim durch meine zwei Builds bestätigt |
| F-2 (Plan-Nachzug) | **ja, mit Rest V-2** | `1ac3558`; Port-Datei `replicationack.go` fehlt in §3 |
| F-1/F-5 (Kennungen / Struktur-IDs in Messagen) | **kein Wiederauftreten** — wirkt vorwärts | alle fünf Fix-Commits tragen `LH-*`/`ADR-*`; keine `SPEC-*`/`ARC-*` in den Messagen; `doc-commits` 0 Befunde über die volle Range |
| F-11 (Dup-DoD) | **ja** | in `1ac3558` entfernt |
| F-11/F-12/F-13/F-16 (Planner-Residuen) | **offen, korrekt bei Planner** | [`SPEC-008`](../../spec/pflichtenheft.md)-Klassenrand TRUNCATE→`schema` unverändert (F-12); Wraparound weiter unbenannt (V-3); `PH-DEP-002`/`PH-TST-001` in slice-006 weiter undeklariert (F-16) |

## Befunde

### V-1 — Die Datei-Abschluss-Klasse (≥3. Auftreten) ist im Fix-Zug nicht geschlossen

- `kategorie`: MEDIUM
- `quelle`: `review-slice-005.md` F-6 (3. Auftreten, MEDIUM-Stufe) ·
  `verify-slice-004.md` V-4 (dieselbe Klasse, `harness/README.md`) ·
  Modul 5/6 (Steering-Loop-Pflicht ab 3×)
- `pfad`: `harness/README.md`, `docs/plan/adr/0044-image-beleg-semantik.md`,
  `Dockerfile` — je `tail -c1 … | od -An -c` → kein `\n`; dagegen sauber:
  `go.mod`, `go.sum`, `Makefile`, `run-replication-tests.sh` und alle neun
  neuen Go-Dateien des Ranges
- `befund`: Der Fix `11d9472` (F-6) schließt die vom deps-Zug betroffenen
  Lock-Files — aber dieselbe Klasse tritt im selben Zug an drei Stellen
  wieder auf, die der Fix- und Entscheidungs-Zug selbst berührte: die
  Werkzeuge-Zeile-Ausgabe und die [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)-Verkabelung (`7b3a52b` bzw.
  `3140568`). Die Klasse steht damit bei einem weiteren Auftreten über der
  3×-Schwelle — bei der Closure als Beobachtung zu führen (Kennung vergeben
  oder zitieren, Ausgang zuweisen), nicht still zu lassen.
- `verifizierbar`: ja — `tail -c1 <Datei> | od -An -c` je Datei (Beleg oben)
- `klasse`: Datei-Abschluss (Fix-Zug unvollständig — Klasse ≥3×, weiter offen)

### V-2 — Plan-Nachzug unvollständig: der ACK-Port liegt in keiner §3-Zeile

- `kategorie`: LOW
- `quelle`: Modul 5 („Wer später mitnimmt …, hat den Plan **geändert**") ·
  Review-F-2 (nannte den Port-Kontrakt ausdrücklich)
- `pfad`: `internal/application/port/outbound/replicationack.go` (neu, 32
  Zeilen: Sentinel + Interface) gegen die §3-Tabelle am HEAD (7 Zeilen —
  keine nennt den Port)
- `befund`: Der F-2-Nachzug deckt die drei benannten Entscheidungen und die
  Verkabelung, aber genau die vom Review als „Vertragserweiterung am
  Outbound-Port" genannte Datei fehlt. Klein und ohne Eigenentscheidung,
  aber die F-2-Klasse ist damit nicht vollständig geschlossen — §3-Zeile
  nachziehen (Planner, vor der Closure).
- `verifizierbar`: ja — §3-Tabelle gegen `git show --name-only`
- `klasse`: Plan-Erweiterung ohne Plan-Nachzug (Rest — Nachzug unvollständig)

### V-3 — XID-Wraparound-Grenze weiterhin unbenannt (F-13 läuft weiter)

- `kategorie`: INFO
- `quelle`: [`ADR-0011`](../plan/adr/0011-persist-before-ack.md) (Deduplizierung) · Review-F-13 (INFO)
- `pfad`: `mapper.go:105` (Transaktions-Kennung = dezimaler XID) · kein
  Eintrag in [`ADR-0011`](../plan/adr/0011-persist-before-ack.md), kein Plan-§6-Ausgang, kein
  Code-Kommentar (`grep -i wraparound` → 0 Treffer in ADR, Mapper, Plan)
- `befund`: Nach 2^32 Quelltransaktionen wiederholen sich XIDs; die
  Idempotenz-Form dedupliziert dann still, was keine Wiederholung ist. Im
  MVP-Fenster fern, aber die Grenze lebt weiterhin nur im Review-Report —
  Planner-Item (Kontrakt-Notiz oder §6-Ausgang), kein DoD-Bruch.
- `verifizierbar`: ja — grep-Beleg
- `klasse`: Unbenannte Grenze am Spec-Rand (F-13 läuft weiter)

### V-4 — Undeklarierte ID-Präfixe in slice-006 weiter unkorrigiert (F-16)

- `kategorie`: INFO
- `pfad`: `docs/plan/planning/open/slice-006-integrationstest-umgebung.md:12`
  (`PH-DEP-002`, `PH-TST-001`) — der Plan-Nachzug `1ac3558` berührte
  slice-006, ohne die Präfixe zu korrigieren
- `befund`: Vorbestand-Residue, Planner-Sache; kein Befund gegen den
  Implementer-Diff.
- `verifizierbar`: ja — MR-000-Deklaration gegen die Plan-Zeile
- `klasse`: Undeklariertes ID-Präfix (Vorbestand, berührt und unkorrigiert)

## Negativbefunde

- geprüft, ohne Befund: **Docker-only (AGENTS §3.1)** — alle go-Belege im
  gepinnten Toolchain-Container (`--network none` für Build/Vet/`make test`,
  Modul-Cache im Volume), Image-Builds aus `/tmp`-Worktrees; kein
  Host-Toolchain-Aufruf, kein Schreibzugriff auf den Hauptbaum
- geprüft, ohne Befund: **Keepalive-Zusage am realen Pfad (Prüfpunkt a)** —
  der Test hält `confirmed_flush_lsn` an der bestätigten Position und meldet
  das Überschreiten sofort rot; die rote Probe (Empfangsstand → 0) wird am
  realen Treiber gefangen (Mutation 1, Beleg oben)
- geprüft, ohne Befund: **Doppel-BEGIN (Prüfpunkt b)** — der Sentinel
  `ErrBeginWithoutCommit` endet sichtbar statt still zu überschreiben;
  der Guard-Träger wird von der Unit-Suite bewacht (Mutation rot)
- geprüft, ohne Befund: **Persist-before-ACK am realen Treiber** —
  `TestRealPersistBeforeAck` (8.4 s) verdrahtet echten Capture Service +
  Store + ACK-Adapter; die bestätigte Position folgt der Persistenz
- geprüft, ohne Befund: **Image-Digest reproduziert unverändert** — zwei
  Buildx-Builds (7b3a52b, HEAD) lieferten exakt `9ac4a9fb…`; Binary byte-
  identisch (`43c3aec0…`) — der Claim aus `7b3a52b` ist belegt, und der
  Beleg bleibt auch am HEAD gültig (Binary invariant, Bootstrap-Stub)
- geprüft, ohne Befund: **Testcontainer-Hygiene** — Daten/Publication/Slot
  im Container, tabellenscopierte Publications (Fix in `0bc3277`), Slot-Rückbau
  mit Retry, Netz- und Container-Rückbau in jedem Ausgang; nach meinen
  Läufen: PG-Container, Netz und Temp-Images entfernt
- geprüft, ohne Befund: **Traceability der Fix-Commits** — jeder trägt
  `LH-*`/`ADR-*`-Kennungen; keine `SPEC-*`/`ARC-*` in Messagen; `doc-commits`/
  `doc-immutable` über `cb021dc..HEAD` je 0 Befunde
- geprüft, ohne Befund: **superseded-Referenzen / ADR-Index** —
  [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) im Index (Accepted); keine
  tragenden Referenzen auf abgelöste ADRs im Range
- geprüft, ohne Befund: **a-check-Abdeckung** — `make gates` (a-check im
  Bündel) 0 Befunde am HEAD; Verdrahtungs-Test liegt in der
  Composition-Root
- geprüft, ohne Befund: **§1-Abgrenzung** — kein Spooling, keine HTTP-/gRPC-
  API, keine Consumer-Verwaltung, kein d-migrate-Rollout im Range

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Datei-Abschluss (Klasse ≥3×, im Fix-Zug
nicht geschlossen — harness/README, [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md), Dockerfile) · Plan-Nachzug
unvollständig (ACK-Port ohne §3-Zeile) · unbenannte Grenze
(XID-Wraparound, F-13 läuft) · undeklariertes ID-Präfix (Vorbestand
slice-006, F-16).

**Zusammenfassung DoD:** **4/10 Punkte jetzt erfüllt** (Items 1–4, Belege
selbst gefahren, inkl. 10/10 Test-Pakete netzlos und der Replikations-Suite
gegen die reale gepinnte PostgreSQL mit `wal_level=logical`) · 1 getan,
Häkchen offen (Item 5) · 1 entfällt (Item 7) · 3 erfüllen sich erst bei
Closure (Items 6, 8, 9-Endwert) · 1 delegiert an die Welle-2-Closure
(Item 10). **Alle Review-Blockierpunkte für die Merge-Form sind beantwortet:**
F-2 (Nachzug `1ac3558`, Rest V-2 klein), F-3/F-4/F-8/F-9/F-10/F-14 (Tests
am realen Pfad bzw. Sentinel, zwei Mutationen rot repliziert), F-6 (go.mod/
go.sum; Rest V-1), F-7 ([`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md), Reproduktion bestätigt).

## Verdikt

**Merge-blockierend:** nein — der Adapter ist am realen gepinnten Treiber
belegt (Replikations-Suite grün, Verdrahtungs-Test 8.4 s am echten Capture
Service), die kritischen Zusagen (Keepalive-Position, Doppel-BEGIN,
Restart, Persist-before-ACK) tragen Tests und sind durch zwei replizierte
Mutations-Proben bewacht; `make gates` ist grün am HEAD, Build/Vet/Gofmt
clean im gepinnten Container, der Image-Digest reproduziert unverändert und
die [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)-Semantik trägt Dockerfile und Werkzeuge-Zeile.

**Blockierend für Closure (der normale Zustand eines Slice in
`in-progress/`):**

1. **V-1:** die Datei-Abschluss-Klasse ist im Fix-Zug nicht geschlossen —
   `harness/README.md`, [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) und `Dockerfile`
   enden ohne Zeilenumbruch; bei der Closure als Beobachtung zu führen
   (≥3×-Schwelle, Kennung zitieren oder vergeben).
2. **V-2:** §3-Zeile für `replicationack.go` nachziehen (Planner, vor der
   Closure).
3. **Closure-Pflichten (Items 5-Häkchen, 6, 8-Notiz, 9-Endausgänge):**
   Doku-Häkchen, Closure-Notiz mit Lerneintrag (Kandidaten: Commit-Kennungs-
   Klasse F-1, Struktur-ID-Klasse F-5 — beide Sequenzen laufen; Datei-Abschluss
   als Register-Kandidat V-1), Endausgänge für §6 (a)/(b) plus die im Range
   entstandenen Review-Risiken (Keepalive/Restart — durch `0bc3277` belegt).
4. **V-3/V-4 als Planner-Hinweise ohne Blockier-Charakter:** Wraparound-
   Grenze benennen (Kontrakt-Notiz oder slice-006), undeklarierte Präfixe in
   slice-006 korrigieren.

**Übergabe:** Bericht an den Planner. Keine Reparaturen. Der
Review-Report-Verdikt-Punkt 2 (F-3/F-4) ist mit `0bc3277` **erfüllt** —
beide Zusagen tragen Tests am realen Pfad; Punkt 1 (F-2) mit `1ac3558`
erfüllt (Rest V-2); Punkt 4 (F-7) über [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)
adjudiziert und durch meine Builds bestätigt; Punkt 3 (F-5-Sequenz) läuft
weiter — Konventions-Nachzug über den Architect bleibt fällig.