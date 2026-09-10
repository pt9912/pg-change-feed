# Verifier-Report: slice-006 — 2026-09-10

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of Done,
**10 Items** — der Auftrag nannte 11; §2 trägt zehn, kein Duplikat), §3
(Plan-vs-Code, Range `c286339..HEAD` inkl. Plan-Nachzüge), §6
(Risiko-Ausgänge) und ADR-Konformität ([`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) ·
[`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) · [`ADR-0026`](../plan/adr/0026-composition-root.md) ·
[`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md)/0029). Nicht geprüft:
Diff gegen Plan/Hard Rules im Detail (Reviewer, `review-slice-006.md`,
Verdikt dort), realer Bedarf (Validator).

**Gegenstand:** Implementer-Handoff zu
`docs/plan/planning/in-progress/slice-006-integrationstest-umgebung.md` ·
Range `c286339..HEAD` · Fix-Commits `32b9e81` (F-1+F-7), `da96157`
(Linkpflicht), `97986ae` (Plan-Korrekturen F-3/F-4/F-6/F-11 + slice-007),
`00a7c2a` (Beleg-Erneuerung). **Während dieser Verifikation wanderte der
HEAD** durch die parallele F-2-Verkörperung des Architect: `57e20af`
([`ADR-0045`](../plan/adr/README.md)) → `154d9b4` (Sensor-Mechanik) → `6b3ae0c` (Doku-Bindung). Die
Sensor-Läufe zu Gates/Test/Rollout liefen am Handoff-Head `00a7c2a`; die
Traceability- und Standing-Gate-Läufe am dann aktuellen HEAD.

**Grundsatz:** Es wurden **keine Behauptungen übernommen** — jeder Sensor
unten wurde in diesem Lauf selbst gefahren; Ausgaben sind Belege über
stdout. Der Hauptbaum blieb read-only (`git status` clean nach jedem Lauf);
die Rollout-Proben (frische Instanz, Idempotenz, Katalog- und Rote-Probe)
liefen in einem `/tmp`-Worktree (nach Abschluss entfernt), die Compose-DB
dort nach Abschluss abgeräumt (`down -v`, Netz entfernt, rote Probe-Schema
gedroppt).

**Modell:** Claude Code (glm-5.3-flash) · **Datum:** 2026-09-10

**Eingangs-Kontext:**

- Slice-Plan §1–§8 (am HEAD) · `docs/reviews/review-slice-006.md`
  (F-1…F-12, Verdikt: vier Closure-Blockierpunkte) · Fix-Commits im
  Volltext · `docs/reviews/verify-slice-005.md` (Gerüst)
- `spec/lastenheft.md` (MVP-Schnitt §1, [`LH-FA-CAP-001`](../../spec/lastenheft.md)…004/008,
  [`LH-FA-REA-002`](../../spec/lastenheft.md)/004.a/005, [`LH-FA-CFG-001`](../../spec/lastenheft.md),
  [`LH-QA-POR-003`](../../spec/lastenheft.md)) · `spec/pflichtenheft.md` ·
  `spec/architecture.md` ([`ARC-007`](../../spec/architecture.md))
- `compose.yaml`, `Makefile`, `tools/harness/run-integration-tests.sh`,
  `tools/schema/{schema.yaml,plan.yaml,down.sql,nacharbeit-operation-check.sql,compose-init/}`,
  `test/integration/mvp_test.go`, `.dockerignore`, `Dockerfile`,
  `harness/image-hash.txt`, `.a-check.yml`, `d-check.mk` ·
  Beobachtungs-Register · `docs/plan/planning/welle-2.md` ·
  `docs/plan/planning/open/slice-007-bootstrap-verdrahtung.md`

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` (HEAD `00a7c2a`) | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 97 Datei(en) geprüft, 0 Befund(e)` · `a-check … gesamt: 0 Befund(e)` | 0 |
| `make test` (netzlos, gepinnter Toolchain-Container `golang:1.27-alpine@sha256:cf6fca…`) | **11 Pakete `ok`** — u. a. `test/integration 0.003s` (kompiliert, skippt ohne DSN), `bootstrap`, `replication/receive`, `postgresack` | 0 |
| `make test-integration` (HEAD `00a7c2a`) | **grün** — Compose frisch (postgres + feed), Feed-Smoke Exit 0, `schema-validate` 5 Tabellen/21 Spalten/1 Index/5 Constraints, Rollout `--execute`, Nacharbeit-Schritt sichtbar (NOTICE „constraint … does not exist, skipping" — die F-5-Grenze im eigenen Lauf gesehen), `test/integration ok 0.696s`; **Arbeitsbaum danach clean** — `plan.yaml`/`down.sql` bit-stabil | 0 |
| `make schema-validate` (standalone, `/tmp`-Worktree) | 5 Tabellen, 21 Spalten, 1 Index, 5 Constraints, 0 Warnings | 0 |
| `make schema-rollout` (frische Instanz, `/tmp`-Worktree, gepinnte Digests: d-migrate `8d1433…`, postgres `63bdc97…`) | Exit 0, Nacharbeit-Schritt läuft; **Katalogprobe: `committed_at` → `is_nullable = NO`, `timestamp with time zone`, `CURRENT_TIMESTAMP`** | 0 |
| Rote Probe (eigene, gegen dieselbe Instanz) | `CREATE TABLE` mit der **Vor-Zustand-Form** (exakt der Statement-Text des Vor-Fix-Belegs `f79af17:tools/schema/plan.yaml`: `committed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP` **ohne NOT NULL**) im Wegwerf-Schema `redprobe` → **`is_nullable = YES`**; derselbe Katalog zeigt für den korrigierten Rollout `cdc.committed_at` → **`NO`** | — |
| Zweiter `make schema-rollout` gegen die **bestückte** Instanz | **scheitert**: `E012 Check expression 'chk_change_operation' references unknown column 'ARRAY'/'text'` → `make: *** Fehler 3` — Idempotenz trägt nur über die Runner-Räumung (frische Compose-Instanz je Lauf; Grenze im Makefile-Kommentar benannt) | 1 |
| `make doc-commits RANGE=c286339..HEAD` | **3 Befunde** (`commit-untraceable`): `3ea9244`, `779bc64` (bekannt, F-2) **und `97986ae`** — der Plan-Korrektur-Commit selbst trägt keine `LH-*`-/`ADR-*`-Kennung | 1 |
| `make commit-traceability` (Standing-Gate, HEAD~5..HEAD, am HEAD `6b3ae0c`) | **1 Befund: `97986ae`** — das neue Gate fängt den Commit, der es begleitete, selbst | 1 |
| `make doc-immutable RANGE=c286339..HEAD` | `0 Befund(e)` (MR-Immutabilität über die volle Range) | 0 |
| `make doc-planning` | `0 Befund(e)` | 0 |
| `git status` / `git worktree list` / `docker ps` | Hauptbaum clean; nur der Hauptbaum; `/tmp`-Worktree entfernt; keine `cdc-*`-Container/-Netze übrig | — |

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt — 10 Items)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | Compose-Umgebung startet PostgreSQL + PG Change Feed reproduzierbar aus dokumentierten Schritten — Teil-Beleg zu [`LH-QA-POR-003`](../../spec/lastenheft.md) | **bestätigt** | `compose.yaml` (Postgres Digest-Pin, `wal_level=logical`, Healthcheck, Init-SQL; Feed-Container am `:dev`-Image) — mein `test-integration`-Lauf fuhr die Umgebung frisch hoch, Feed-Smoke Exit 0, Kette dokumentiert (`harness/README.md:128` Werkzeuge-Zeile) |
| 2 | MVP-Integrationstest läuft automatisiert grün (MVP-Schnitt §1) | **bestätigt** | eigener Lauf `TestMVPCaptureFlow` + `TestMVPUpdateOldImageWithFullReplicaIdentity` — grün gegen die Compose-DB; Abfolge, Reihenfolge, Inhalt, Wiederlesen, Bereichs-Grenze (Details unten, Prüfpunkt c) |
| 3 | `make gates` grün | **bestätigt** | drei Gates grün am HEAD (Tabelle oben) |
| 4 | Review durchgeführt, Report unter `docs/reviews/` | **bestätigt** | `docs/reviews/review-slice-006.md` committet (`da96157`); Rollenwechsel nach Schritt 8 eingehalten |
| 5 | Doku-Update bzw. begründete Aussage | **getan, Häkchen noch offen** | `harness/README.md` trägt die Werkzeuge-Zeile `test-integration` (Zeile 128, [`ADR-0030`](../plan/adr/0030-testpyramide.md)/[`LH-QA-POR-003`](../../spec/lastenheft.md)) — getragen; Häkchen fällig vor dem `git mv` |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **offen bis Closure (notiert)** | §7 trägt Platzhalter; fällig vor dem `git mv` (Lerneintrag-Kandidaten: Commit-Kennung-Klasse — jetzt 4× mit V-3 · Nacharbeit-Retirement V-1 · Datei-Abschluss-Klasse) |
| 7 | Reconciliation-Register | **entfällt** | `docs/plan/planning/reconciliation.md` existiert nicht (Repo ohne Brownfield-Bootstrap) — vom DoD-Wortlaut ausdrücklich vorgesehen |
| 8 | Beobachtungs-Register fortgeschrieben | **kein neues Auftreten — notiert (Closure-Antwort fällig)** | `BEO-PGC/a-check-null-abdeckung` unverändert (`verkörpert · seit welle-1`, evidence `slice-001/002/003.md`); kein neues Auftreten im Range. **Neue Register-Kandidaten für die Closure:** die Datei-Abschluss-Klasse (≥3×, im Fix-Zug `32b9e81` geschlossen, aber kein Sensor trägt sie) und die Commit-Kennung-Klasse (4× mit V-3 — Verkörperung läuft, Ausgang „geplant" mit [`ADR-0045`](../plan/adr/README.md)-Verkabelung) |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **(a)/(b) am Beleg — Endausgänge fällig; die vom Review geforderten Zusatz-Ausgänge FEHLEN (V-1)** | (a) reale Ordnung: strenge Positionsordnung im eigenen Lauf belegt; (b) Row-Images: beide Identity-Formen im Test (`TestMVPUpdateOldImageWithFullReplicaIdentity` trägt Alt-Bild beider Spalten). **Aber:** der vom Review-Verdikt (Punkt 2) geforderte Nacharbeit-Retirement-Ausgang (F-4) ist **nirgends** getragen — `97986ae` behauptet ihn, liefert ihn nicht (V-1 unten) |
| 10 | Drei Paarungen (Anker · Folge-Slice · Register) | **delegiert — korrekt** | Repo **mit** Wellen-Betrieb (`welle-2.md` flach, Slice trägt `Welle: welle-2`): Prüfung durch die Welle-2-Closure |

**Zwischenstand: 4/10 jetzt erfüllt** (Items 1–4, Belege selbst gefahren) ·
1 getan, Häkchen offen (Item 5) · 1 entfällt (Item 7) · 3 erfüllen sich erst
bei Closure (Items 6, 8 als Notiz, 9-Endwert) · 1 delegiert (Item 10) —
**plus zwei ungelöste Review-Closure-Blockierpunkte, deren Behebung in
`97986ae` behauptet, aber nicht geliefert wurde (V-1/V-2).**

## Plan-vs-Code-Diff (Range `c286339..HEAD`)

**Deckung §3:** die vier §3-Zeilen sind geliefert — `compose.yaml` (neu),
`test/integration/mvp_test.go` (neu), `Makefile` (schema-rollout-Verkabelung),
`tools/schema/schema.yaml` (neu, [`ADR-0043`](../plan/adr/README.md)-Erstlieferung). Die
Rollout-Belege (`plan.yaml`, `down.sql`), `nacharbeit-operation-check.sql`,
`compose-init/` und `tools/harness/run-integration-tests.sh` sind
Beleg-/Runner-Artefakte des im §3-Makefile-Zeile nachgezogenen
[`ADR-0043`](../plan/adr/README.md)-Erstversatzes — zulässige Deckung über die Plan-Nachzüge.

**Unbudgetiert am HEAD — F-3-Nachzug fehlt (V-2):** `.a-check.yml`
(Glob-Erweiterung `test/integration/**`) und die `harness/README.md`-Zeile
stehen **nicht** in §3 — obwohl die `97986ae`-Message genau diesen Nachzug
behauptet („F-3 … als §3-Erweiterung nachgezogen"). Der Commit berührt die
slice-006-Plan-Datei **gar nicht** (nur `done/slice-003`, `open/slice-007`,
`welle-2.md`); §3 hat am HEAD unverändert vier Zeilen. Die F-3-Klasse
(6. Auftreten laut Review) ist damit **behauptet, aber nicht geschlossen**.

**Getragen über die Review-Alternative:** F-6 (Verdrahtungs-/Aktivierungs-Grenze)
— `slice-007` in `open/` nimmt die Verdrahtung an (§6-Ausgang „eingetreten;
Träger ist dieser Slice") und `welle-2.md` §1 trägt die Grenze als
benannte Notiz. Das ist die vom Review ausdrücklich zugelassene Alternative
(„Ausgang im Plan §6 **oder** Welle-§1-Formulierung schärfen") — getragen.

**Nicht in §3, zulässig (Lauf- bzw. Decision-Artefakte):** Review-Report
(`da96157`) · [`ADR-0045`](../plan/adr/README.md) + Index (`57e20af`, parallele F-2-Verkörperung) ·
slice-007-Plan (Folge-Slice, Planner) · Plan-Nachzüge/Korrekturen
(`da96157`, `97986ae`) · `.claude/settings.json` (Zweitschreiber-Abschnitt
des Reviews).

## ADR-Konformität

- **[`ADR-0043`](../plan/adr/README.md)** (Erstversatz, Pflicht-Report, Rollback-Artefakt): **konform und
  belegt** — Ersteinsatz ist der Compose-Rollout vor jedem E2E-Lauf (Runner-Kette
  im eigenen Lauf gefahren); Pflicht-Report `tools/schema/plan.yaml` (`status ok`,
  `postUpVerified: true`, 5 Operationen) und Rollback-Artefakt `down.sql` sind
  committet und in meinem Lauf bit-stabil reproduziert (Arbeitsbaum clean).
  **Die F-1-Überführungs-Abweichung ist geschlossen und durch meine eigene
  Katalog- und Rote-Probe belegt** (siehe Prüfpunkt a) — die Überführung
  verliert **keine** Constraint-Form mehr. Grenzen, die bleiben: F-5 (Beleg-Kette
  endet vor der Nacharbeit — im eigenen Lauf als NOTICE sichtbar) und die
  E012-Idempotenz-Grenze (Prüfpunkt b).
- **[`ADR-0044`](../plan/adr/README.md)** (Image-Beleg-Semantik): **konform** — `compose.yaml` trägt
  keinen `build:`-Block (Kommentar-Klasse Abgrenzung nennt die Begründung);
  der Range berührt **keine** Build-Kontext-Datei (`.dockerignore`: nur `cmd/`,
  `internal/`, `go.mod`, `go.sum` — keine davon im Range) → kein `make image`-Lauf
  fällig, `harness/image-hash.txt` (`sha256:9ac4a9fb…`) bleibt der Lauf-Beleg
  unverändert. Grenze F-8 (`:dev`-Tag ohne Beleg-Abgleich) bleibt LOW offen.
- **[`ADR-0026`](../plan/adr/README.md)** (Composition Root): **konform** — der Test verdrahtet Store,
  Stream, ACK und Capture-Service über die Composition-Root-Ebene
  (`startCapture`); a-check 0 Befunde im Bündel-Lauf.
- **[`ADR-0023`](../plan/adr/README.md)/[`ADR-0029`](../plan/adr/README.md)** (Fehlerklassen, Regel 1/3/6): **unberührt** — der Range
  trägt keinen Produkt-Code; der Rollout-Pfad trägt sichtbare Fehler
  (`ON_ERROR_STOP=1`, Exit-Ketten).

## Fix-Commits gegen die Review-Findings

| Finding | getragen im Fix? | Beleg |
|---|---|---|
| F-1 (Schema-Drift `committed_at` NOT NULL) | **ja, mit eigenem Beleg** | `32b9e81` trägt exakt `required: true` an `schema.yaml`; `00a7c2a` regeneriert Belege — **meine eigene Katalogprobe** an frischer Instanz: `is_nullable = NO`; **Rote Probe**: Vor-Zustand-Form (committeter Vor-Fix-Beleg + Wegwerf-`CREATE TABLE`) → `YES` |
| F-2 (Kennungspflicht, 3. Auftreten) | **Verkörperung läuft — siehe eigene Sektion** | [`ADR-0045`](../plan/adr/README.md) am HEAD (`57e20af`), Sensor-Mechanik am HEAD (`154d9b4`, `6b3ae0c`); Reste: V-3 und die Bündel-Verkabelung |
| F-3 (Plan-Nachzug `.a-check.yml`/`harness/README`) | **NEIN — behauptet, nicht geliefert (V-2)** | §3 am HEAD unverändert vier Zeilen; `97986ae` berührt die Plan-Datei nicht |
| F-4 (Nacharbeit-Retirement ohne Ausgangs-Träger) | **NEIN — behauptet, nicht geliefert (V-1)** | `grep -rni "nacharbeit\|retirement" docs/plan/planning/` → **0 Treffer**; slice-006 §6 unverändert, slice-007/welle-2 ohne Erwähnung, kein BEO-Eintrag, kein Folge-Slice dazu |
| F-4-Nebenbefund F-5 (Beleg-Kette vor Nacharbeit) | **unverändert offen** | eigener Rollout-Lauf: NOTICE „constraint … does not exist, skipping" — Report kennt den Constraint nicht; trägt die F-4-Verankerung mit |
| F-6 (Verdrahtungs-/Aktivierungs-Grenze) | **ja, über slice-007** | Folge-Slice in `open/` mit Ausgang „eingetreten; Träger ist dieser Slice"; `welle-2.md` §1 Grenz-Notiz |
| F-7 (Datei-Abschluss, 5 Dateien) | **ja** | `32b9e81` trägt exakt die fünf genannten Dateien; `tail -c1`-Probe: alle fünf enden auf `\n` |
| F-11 (`PH-*`-Präfixe) | **ja** | Plan-Kopf korrigiert (`da96157`), ebenso `done/slice-003` (`97986ae`) |
| F-12 (welle-2 §3 Platzhalter) | **teilweise — überzeichnet (V-4)** | Platzhalter sind gestrippt, aber §3 steht jetzt **ohne** Closure-Trigger da; der §2-Glitch „bereits eingetreten.>" bleibt |
| Kennungspflicht der Fix-Commits selbst | **Rest: `97986ae` (V-3)** | `doc-commits`: 3ea9244, 779bc64, **97986ae** `commit-untraceable`; `32b9e81`/`da96157`/`00a7c2a` tragen [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)/`LH-QA-POR-003` |

## Explizite Prüfpunkte

**(a) F-1-Rote-Probe (Vorzustand `is_nullable = YES`):** **bestätigt und
über die Behauptung hinaus belegt.** Der Vor-Zustand ist dreifach belegt: der
committete Vor-Fix-Beleg `f79af17:tools/schema/plan.yaml` trägt den
Statement-Text ohne NOT NULL; mein Wegwerf-`CREATE TABLE` mit exakt diesem
Text ergab `is_nullable = YES`; der korrigierte Rollout (mein eigener Lauf,
frische Instanz) ergibt `NO`. Die Korrektur ist damit nicht nur
behauptet, sondern gegen den gemeldeten Befund gerichtet belegt.

**(b) `make schema-rollout` idempotent:** **nicht in der wörtlichen Lesart
bestätigt.** Der zweite Lauf gegen die **bestückte** Instanz scheitert mit
E012 (d-migrate-Introspektion des Nacharbeit-Constraints liest
`'INSERT','UPDATE','DELETE'` als Spaltennamen `ARRAY`/`text`); der fehlgeschlagene
Lauf überschreibt zudem `plan.yaml` (`status`/`exitCode`-Felder). Tragend ist
die Grenze nur in der **Runner-Form**: `run-integration-tests.sh` räumt die
Umgebung vor jedem Rollout ab, so dass der automatisierte Pfad stets gegen
eine frische Instanz läuft — genau diese Grenze ist im Makefile-Kommentar
benannt („scheitert an der Katalogform des Constraints (E012) — die
Runner-Kette … räumt die Umgebung vorher ab"). Die Implementer-Claim
„`plan.yaml`/`down.sql` bit-stabil über beide Läufe" ist richtig **für zwei
frische Läufe** (mein `test-integration`-Lauf: Arbeitsbaum clean) — sie ist
keine Idempotenz-Zusage über bestückte Instanzen, und die sollte auch keiner
lesen.

**(c) MVP-Schnitt-Abfolge und Aktivierungs-Form:** **bestätigt.** Die Abfolge
PostgreSQL starten → CDC aktivieren → INSERT/UPDATE/DELETE → lesen →
Reihenfolge/Inhalt läuft im eigenen Lauf am realen Treiber: drei getrennte
Quelltransaktionen, strenge Positionsordnung, INSERT-Neu-Bild/DELETE-Alt-Bild,
UPDATE ohne Alt-Bild bei Default-Identity, Gegenprobe mit `REPLICA IDENTITY
FULL` (beide Spalten), Wiederlesen in derselben Ordnung, Bereich hinter
letzter Position leer. Die **Aktivierung ist direktes SQL**
(`mvp_test.go:91-105`: `CREATE PUBLICATION` + Bindungs-Zeilen in
`cdc.source`/`source_table`/`schema_version`) — marker-konsistent
([`LH-FA-CFG-001`](../../spec/lastenheft.md) trägt keinen MVP-Marker; welle-2 deklariert die
SQL-/CLI-Adapter als spätere Wellen). **Die Verdrahtungslücke trägt der Plan
offen:** `slice-007` (Bootstrap-Verdrahtung) liegt in `open/` und nimmt den
F-6-Ausgang an; `welle-2.md` §1 nennt die Grenze am M1-Claim. Der
Plan-§6-Träger ist damit der Folge-Slice, nicht slice-006 selbst — die vom
Review zugelassene Form.

## F-2-Verkörperung (parallele Architect-Arbeit) — Zustand am HEAD

Der Auftrag: notieren als getragen, prüfen, ob sie schon am HEAD sitzt.
Stand am Ende dieses Laufs (HEAD `6b3ae0c`): **mechanisch am HEAD, mit zwei
Residuen.**

- **Am HEAD:** `ADR-0045` (Accepted, `57e20af`), Sensor-Mechanik (`154d9b4`:
  `tools/harness/commit-traceability.sh`, Makefile-Target, `.d-check.yml`
  commits-Abschnitt) und die Doku-Bindung (`6b3ae0c`:
  `harness/README.md` Sensors-Zeile mit Herkunfts-Anker `seit slice-006`,
  `AGENTS.md` §4-Zeile). Eigener Lauf: `make commit-traceability` fährt
  (d-check commits-Modul + Betreff-Grenze) und meldet **1 Befund: `97986ae`**
  — das neue Gate fängt den Commit, der die F-2-Klasse während der
  Verkörperung wiederholte, selbst (V-3).
- **Residuum 1:** `harness/README.md:117` behauptet, `make gates` führe
  „baseline-verify, docs-check, a-check, **commit-traceability**" im Bündel —
  `GATE_CHECKS` trägt **nur** `baseline-verify`, `docs-check`, `a-check`;
  es gibt keine `GATE_CHECKS += commit-traceability`-Zeile. Mein
  `make gates`-Lauf am neuen HEAD ist grün, **ohne** dass das neue Gate im
  Bündel lief — der Doku-Claim ist an der Verkabelung vorbei. (Parallele
  Architect-Arbeit, nicht slice-006-Diff — Notiz an den Planner/Architect.)
- **Residuum 2:** das Standing-Gate ist am HEAD **rot** (`97986ae` im
  HEAD~5..HEAD-Fenster) — die Verkörperung ist damit funktionsfähig und
  fängt den eigenen Nachlauf; sie ist **nicht grün abgeschlossen**.

## Befunde

### V-1 — F-4 (Nacharbeit-Retirement-Ausgang) behauptet, nirgends geliefert

- `kategorie`: HIGH
- `quelle`: Review-Verdikt Punkt 2 („F-4 … braucht einen Ausgangs-Träger
  **vor** der Closure") · Modul 5 („kein Slice geht nach `done/`, während
  eines ohne Ausgang dasteht") · Modul 11 (Behauptung ohne Bestätigung)
- `pfad`: Commit-Message `97986ae` („F-4: d-migrate-Nacharbeit-Retirement
  als §6-Ausgang getragen (die gemeldete Fix-Release retiriert die
  Ausweichform; Rückbau-Prüfung im Rollout-Beleg)") gegen
  `docs/plan/planning/in-progress/slice-006-integrationstest-umgebung.md` §6
  (unverändert: nur Risiken (a)/(b)), `open/slice-007-bootstrap-verdrahtung.md`
  (keine Nacharbeit-Erwähnung), `welle-2.md`, `observations/` (unverändert)
- `befund`: Die Message behauptet die Behebung des Review-Blockierpunkts 2 —
  der Commit berührt die slice-006-Plan-Datei **nicht**, und ein
  abschließender `grep -rni "nacharbeit|retirement" docs/plan/planning/`
  liefert **null Treffer**; auch `chk_change_operation` wird im gesamten
  Planning-/ADR-Stratum nirgends als Ausgangs-Träger genannt. Der vom Review
  geforderte Ausgang (§6-Risiko mit „weiter offen" → `BEO-PGC/<slug>` oder
  Folge-Slice-ID) existiert nicht. Der Slice kann mit diesem Stand nicht
  nach `done/` — und die Message-Behauptung ist die Verifier-Klasse:
  eine Zusage über einen Zustand, den der Baum nicht trägt.
- `verifizierbar`: ja — Commit-Stat (`97986ae` berührt slice-006 nicht) ·
  grep-Beleg · §6-Lese am HEAD
- `klasse`: Behauptete Plan-Korrektur ohne Lieferung (Review-Blocker ungelöst)

### V-2 — F-3 (Plan-Nachzug `.a-check.yml`/`harness/README`) behauptet, nicht geliefert

- `kategorie`: MEDIUM
- `quelle`: Review-Verdikt („F-3 … gehen als Vor-Closure-Nacharbeit") ·
  Modul 5 (Plan-Erweiterung ohne Nachzug, 6. Auftreten)
- `pfad`: slice-006 §3 am HEAD (vier Zeilen — keine nennt `.a-check.yml`
  oder die `harness/README.md`-Werkzeuge-Zeile) gegen die im Range
  gelieferten Berührungen beider Dateien
- `befund`: Dieselbe Message-Behauptung („F-3 … als §3-Erweiterung
  nachgezogen (Architect-Fixes getragen)") ohne Lieferung; die zwei
  gelieferten Dateien bleiben unbudgetiert.
- `verifizierbar`: ja — §3-Tabelle gegen `git show --stat 97986ae`
- `klasse`: Plan-Erweiterung ohne Plan-Nachzug (Klasse ungeschlossen)

### V-3 — `97986ae` selbst commit-untraceable (4. Auftreten der F-2-Klasse)

- `kategorie`: MEDIUM (wirkt vorwärts)
- `quelle`: `harness/README.md` §Traceability rules · review-slice-006 F-2
  (3. Auftreten, Sequenz-Pflicht fällig)
- `pfad`: Commit-Message `97986ae` — Body nennt `F-3/F-4/F-6/F-11`,
  `Architect-Fixes`, `welle-2`, `slice-003/005/006`, aber **keine**
  `LH-*`-/`ADR-*`-Kennung
- `befund`: Mein `doc-commits`-Lauf über die Range meldet drei
  `commit-untraceable`-Befunde: die beiden bekannten (`3ea9244`, `779bc64`)
  **und `97986ae`**. Das neue Standing-Gate flaggt denselben Commit über
  `HEAD~5..HEAD` — der erste eigene Fang der frischen Verkörperung ist ihr
  Geschwister-Commit. Die Klasse steht damit beim vierten Auftreten, und
  zwar am HEAD **während** die Verkörperung gegen genau diese Klasse
  angenommen wurde. Die Messagen bleiben historisch; das Gate wirkt vorwärts
  (im Fünf-Commit-Fenster rot, bis das Fenster weiterschiebt).
- `verifizierbar`: ja — `make doc-commits RANGE=c286339..HEAD` ·
  `make commit-traceability` (Belege oben)
- `klasse`: Commit ohne Vertrags-Kennung (4. Auftreten — Gate getragen,
  Lauf rote)

### V-4 — `welle-2.md` §3 trägt jetzt GAR KEINEN Closure-Trigger (F-12 überzeichnet behoben)

- `kategorie`: MEDIUM
- `quelle`: Modul 6 (Welle schließt durch beobachtbare Closure-Kriterien) ·
  review-slice-006 F-12 · Commit-Message `97986ae` („F-12: welle-2 §3-
  Vorlagen-Platzhalter gestrippt")
- `pfad`: `docs/plan/planning/welle-2.md` §3 — seit `97986ae` **leer** (die
  vier `<z.B. …>`-Platzhalter sind entfernt, kein Trigger-Text geschrieben);
  §2 endet weiter mit dem Glitch „**bereits eingetreten**.>"
- `befund`: Die Message behauptet die Behebung; tatsächlich wurde der
  Symptom-Text entfernt, ohne den Befund zu schließen — die Welle, die
  slice-006 einsammelt und den M1-Claim trägt, hat jetzt **gar keine**
  Trigger-Zeile mehr statt der Platzhalter. Der Review-Befund
  („trägt keinen benannten Closure-Trigger") ist unverändert wahr, jetzt
  ohne die Platzhalter-Sichtbarkeit. Vor der Welle-2-Closure zu schließen
  (Planner): beobachtbarer Trigger benennen (z. B. „alle drei Slices in
  `done/` und `make gates` grün" — Start-Trigger-Disziplin beachten).
- `verifizierbar`: ja — Lese der Welle-Datei §3 am HEAD
- `klasse`: Closure-Trigger unbenannt (Vorbestand, Fehlbehebung)

### V-5 — `slice-007` trägt ein Dup-DoD und unfilled Vorlagen-Platzhalter

- `kategorie`: LOW
- `quelle`: F-11-Klasse (Dup-DoD, review-slice-005 — dort entfernt in
  `1ac3558`) · Modul 5 (Template kopiert und ausgefüllt)
- `pfad`: `docs/plan/planning/open/slice-007-bootstrap-verdrahtung.md` §2
  („`make gates` grün." **zweimal**), §4 Rückführungen (`<Bedingung>`),
  §5 Closure-Trigger (`<…>`), §8 vorgelagerte Prüfungen (`<…>`)
- `befund`: Der neue Folge-Slice trägt dieselbe Dup-DoD-Klasse, die in
  slice-005 als F-11 entfernt wurde, und die kopierten Vorlagen-Plätze sind
  nicht ausgefüllt — slice-006 lag als vollständiger Plan in `open/`.
  Kein Blockier-Befund gegen slice-006; Planner-Arbeit, bevor slice-007 in
  `next/` geht.
- `verifizierbar`: ja — Lese der Datei
- `klasse`: Dup-DoD (Klasse wiederholt) · Vorlagen-Platzhalter im neuen Slice

### V-6 — Rollout-Idempotenz nur in der Runner-Form; F-5 bleibt offen

- `kategorie`: INFO
- `quelle`: Makefile-Kommentar (E012-Grenze, Runner-Räumung) · review-slice-006 F-5
- `pfad`: eigener zweiter Rollout-Lauf gegen die bestückte Instanz (Beleg
  oben) · `tools/schema/plan.yaml` (Report ohne den Nacharbeit-Constraint)
- `befund`: Der zweite Lauf gegen eine bestückte Instanz scheitert (E012);
  der automatisierte Pfad umgeht das über die frische Instanz — die Grenze
  ist im Makefile-Kommentar deklariert und damit getragen, aber der Claim
  „idempotent" trägt nur für den Runner-Pfad. F-5 (Beleg-Kette endet vor der
  Nacharbeit) bleibt unverändert — mein Lauf sah das NOTICE; die Grenze
  gehört in dieselbe Verankerung wie F-4 (V-1).
- `verifizierbar`: ja — eigener E012-Lauf (Beleg oben)
- `klasse`: Grenze deklariert, Claim enger als der Beleg

## Negativbefunde

- geprüft, ohne Befund: **Docker-only (AGENTS §3.1)** — alle Läufe über
  `make`/docker; Toolchain, PostgreSQL und d-migrate je Digest gepinnt
  (`cf6fca…`, `63bdc97…`, `8d1433…`); kein Host-Toolchain-Aufruf; Hauptbaum
  read-only (Rollout-Proben im `/tmp`-Worktree)
- geprüft, ohne Befund: **`make test` netzlos** — 11 Pakete `ok`; der
  MVP-Test skippt ohne DSN sauber und kompiliert im Toolchain-Container
  (`.PHONY`-Fix trägt: das Rezept läuft)
- geprüft, ohne Befund: **Beleg-Determinismus** — mein
  `test-integration`-Lauf am HEAD ließ `plan.yaml`/`down.sql` byte-identisch
  (Arbeitsbaum clean); die committeten Belege dirtien nicht
- geprüft, ohne Befund: **[`ADR-0044`](../plan/adr/README.md)-Scope** — keine Build-Kontext-Datei des
  Ranges (`.dockerignore`-Menge: `cmd/`, `internal/`, `go.mod`, `go.sum`)
  berührt → kein `make image` fällig; `image-hash.txt` unverändert
- geprüft, ohne Befund: **`doc-immutable`/`doc-planning`** über die volle
  Range je 0 Befunde (MR-Immutabilität; Lifecycle-Konsistenz)
- geprüft, ohne Befund: **Feed-Smoke-Hygiene** — Compose-Container und -Netz
  in jedem Ausgang abgeräumt (Runner-`trap`); nach meinen Läufen: keine
  `cdc-*`-Container/-Netze; `/tmp`-Worktree entfernt
- geprüft, ohne Befund: **§1-Abgrenzung** — kein Produktions-Deployment, kein
  Benchmark, kein Exportadapter im Range
- geprüft, ohne Befund: **a-check-Abdeckung** — `make gates` (a-check im
  Bündel) 0 Befunde am HEAD; `test/integration/**` in der
  Composition-Root-Abdeckung

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 3 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Behauptete Plan-Korrektur ohne Lieferung
(F-4-Retirement, V-1) · Plan-Nachzug ungeschlossen (F-3, V-2) · Commit ohne
Vertrags-Kennung (4. Auftreten, V-3 — Standing-Gate am HEAD rot und
funktionsfähig) · Welle-Closure-Trigger unbenannt (V-4) · Dup-DoD/
Vorlagen-Platzhalter im neuen Folge-Slice (V-5) · Idempotenz-Grenze
E012/F-5 offen (V-6).

**Zusammenfassung DoD:** **4/10 Punkte jetzt erfüllt** (Items 1–4, Belege
selbst gefahren: `make gates` grün, `make test` 11 Pakete netzlos,
`make test-integration` grün gegen die gepinnte Compose-Umgebung mit
d-migrate-Rollout) · 1 getan, Häkchen offen (Item 5) · 1 entfällt (Item 7) ·
3 erfüllen sich erst bei Closure (Items 6, 8-Notiz, 9-Endausgänge) ·
1 delegiert an die Welle-2-Closure (Item 10). **Die F-1-Kette ist
geschlossen und gegen den gemeldeten Befund belegt** (`32b9e81` + `00a7c2a`,
Katalogprobe `NO`, Rote Probe `YES`), die F-7-Klasse ist geschlossen, F-6
trägt einen Ausgangs-Träger (slice-007), F-11 ist korrigiert.

## Verdikt

**Merge-blockierend:** nein — der Slice-Inhalt ist am realen gepinnten Pfad
belegt: Compose-Umgebung mit Rollout-Erstversatz ([`ADR-0043`](../plan/adr/README.md), Pflicht-Report
und Rollback-Artefakt bit-stabil reproduziert), MVP-Schnitt-Abfolge am
realen Treiber mit beiden Identity-Formen, drei Gates grün, die F-1-Abweichung
geschlossen (eigene Katalog- und Rote-Probe) und die F-7-Klasse geschlossen.

**Blockierend für Closure:**

1. **V-1 (F-4):** der vom Review geforderte Nacharbeit-Retirement-Ausgang
   existiert **nicht** — die `97986ae`-Message behauptet ihn, der Baum trägt
   ihn nicht (grep null Treffer). Plan-§6-Risiko mit Ausgang „weiter offen"
   → `BEO-PGC/<slug>` (Beleg je Closure) **oder** Folge-Slice-ID, die ihn
   annimmt — vor dem `git mv` (Planner).
2. **V-2 (F-3):** §3-Zeilen für `.a-check.yml` und die
   `harness/README`-Werkzeuge-Zeile nachziehen — dieselbe Vor-Closure-Nacharbeit.
3. **V-4 (welle-2 §3):** Closure-Trigger benennen (Vor der Welle-Closure;
   solange leer, schließt die Welle, die slice-006 einsammelt, ohne
   benannten Trigger).
4. **V-3 (Kennungspflicht, 4. Auftreten):** wirkt vorwärts über das neue
   Standing-Gate — bei der Closure §7 als wiederkehrende Finding-Klasse
   führen; die Verkörperung ([`ADR-0045`](../plan/adr/README.md)) ist angenommen, **aber die
   gates-Bündel-Verkabelung fehlt gegen den `harness/README`-Claim
   (Residuum 1)** — an den Architect zurückzugeben, nicht still zu lassen.

**Notiert (Closure-Pflichten, normaler Zustand eines Slice in
`in-progress/`):** Doku-Häkchen (Item 5), Closure-Notiz mit Lerneintrag
(Item 6 — Kandidaten: Commit-Kennung-Klasse 4× mit [`ADR-0045`](../plan/adr/README.md)-Ausgang
„geplant", Datei-Abschluss-Klasse ≥3×, Nacharbeit-Retirement), Register-Notiz
(kein neues Auftreten), Endausgänge für §6 (a)/(b) — beide am
Integrationstest-Beleg. Item 10 delegiert an die Welle-2-Closure.

**Übergabe:** Bericht an den Planner. Keine Reparaturen.