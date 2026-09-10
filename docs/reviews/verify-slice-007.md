# Verifier-Report: slice-007 — 2026-09-10

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of Done,
**10 Items** — der Auftrag nannte 11; die §2-Liste am HEAD zählt 10, ohne
Duplikat), §3 (Plan-vs-Code, Range `de340b6..5dc855a` inkl. Plan-Nachzug
`ead5c01`), §6 (Risiko-Ausgänge) und ADR-Konformität
([`ADR-0012`](../plan/adr/0012-at-least-once.md) ·
[`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md) ·
[`ADR-0026`](../plan/adr/0026-composition-root.md) ·
[`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) ·
[`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md)).
Nicht geprüft: Diff gegen Plan/Hard Rules im Detail (Reviewer,
`review-slice-007.md`, Verdikt dort), realer Bedarf (Validator).

**Gegenstand:** Implementer-Handoff zu
`docs/plan/planning/in-progress/slice-007-bootstrap-verdrahtung.md` (wellenlos) ·
Range `de340b6..5dc855a` (**HEAD am Prüfzeitpunkt `e971ed4`** — der
Review-Report-Commit liegt über dem im Auftrag genannten `5dc855a`; alle
Sensor-Läufe unten fahren am HEAD `e971ed4`) · Fix-Commits `62fc9b9`
(F-2/F-3-Teil), `a184663` (F-4/F-5/F-3-Teil), `10953e2` (F-6), `dae053d`
(F-3-Rest), `5dc855a` (Image-Beleg-Nachzug), Plan-Nachzug `ead5c01` (F-1/F-2-
Ausgang), Review-Report-Commit `e971ed4` (DoD-Punkt 4).

**Grundsatz:** Es wurden **keine Behauptungen übernommen** — jeder Sensor
unten wurde in diesem Lauf selbst gefahren; Ausgaben sind Belege über stdout.
Der Hauptbaum blieb read-only (`git status --porcelain` leer am Ende,
`git worktree list` trägt nur den Hauptbaum); Mutations- und Sensor-Proben
liefen in einem `/tmp`-Worktree (nach Abschluss revertiert und entfernt); die
Compose-Umgebung meines E2E-Probe wurde nach Abschluss abgeräumt
(`compose down -v`, keine `cdc-test-*`-Container/Volumes zurück).

**Modell:** Claude Code (glm-5.3-flash) · **Datum:** 2026-09-10

**Eingangs-Kontext:**

- Slice-Plan §1–§8 (am HEAD `e971ed4`, inkl. Plan-Nachzug `ead5c01`) ·
  `docs/reviews/review-slice-007.md` (F-1…F-10, committet in `e971ed4`) ·
  Fix-Commits im Volltext
- `spec/pflichtenheft.md` ([`SPEC-008`](../../spec/pflichtenheft.md),
  [`SPEC-015`](../../spec/pflichtenheft.md), [`LH-FA-CFG-001`](../../spec/lastenheft.md)) ·
  `spec/lastenheft.md` (MVP-Schnitt §1, [`LH-QA-OPS-001`](../../spec/lastenheft.md),
  [`LH-QA-POR-003`](../../spec/lastenheft.md)) · `spec/architecture.md` (`ARC-007`)
- `docs/plan/planning/welle-2.md` (§3 Rest-Verdrahtungs-Ausgang) ·
  `harness/README.md` (Werkzeuge-Zeilen `test-integration`/`commit-traceability`,
  §Traceability rules) · Beobachtungs-Register (`BEO-PGC/*`) ·
  [`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md)

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` (HEAD `e971ed4`) | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 103 Datei(en) geprüft, 0 Befund(e)` (voll) · `d-check … --range HEAD~5..HEAD`: 0 Befunde · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` · `a-check gesamt: 0 Befund(e)` | 0 |
| `go build ./...` + `go vet ./...` + `gofmt -l .` im gepinnten Toolchain-Container (`golang:1.27-alpine@sha256:cf6fca…`, `--network none`, `/src:ro`, Modul-Cache-Volume) | `BUILD_OK` · `VET_OK` · **`gofmt -l` meldet `internal/bootstrap/wiring_test.go`** (zwei Ausrichtungs-Blöcke, `gofmt -d` belegt) — V-1 | 1 (Befund-Liste) |
| `make test` (netzlos, gepinnter Container) | **10 Pakete `ok`** — inkl. `internal/bootstrap` (mit den neuen Negativtests aus `10953e2`) und `test/integration` (Skip ohne DSN) | 0 |
| `make test-integration` (Compose-Umgebung, gepinnte Digests: Toolchain `cf6fca…`, PostgreSQL `63bdc97d…`) | **grün**: d-migrate-Validierung + Rollout → Aktivierungs-SQL (`CREATE PUBLICATION`, Bindungs-Zeilen) → Feed-Container gestartet → `TestMVPCaptureFlow PASS (0.14s)` · `TestMVPUpdateOldImageWithFullReplicaIdentity PASS (0.13s)` → End-Wächter grün | 0 |
| Manueller E2E-Probe (eigene Compose-Instanz, nach Abschluss abgeräumt) | Feed-Container streamt selbst: nach 2 INSERTs an `public.feed_mvp_flow` enthält `cdc.change` genau 2 Zeilen (`INSERT`, `new_data {"id": "1", "name": "Alpha"}` / `{"id": "2", "name": "Bravo"}`, `schema_version sv-mvp-flow`) — der einzige Schreiber von `cdc.change` am Pfad ist das Binary im Feed-Container; Slot `slot_pgc_mvp active=t`, `confirmed_flush_lsn 0/1C7AA18`, Container `running exit=0` | — |
| Verhaltens-Probe am geladenen Image (`ghcr.io/pt9912/pg-change-feed:dev`) | `--version` → `pg-change-feed 0.2.0-verdrahtung` (Exit 0) · ohne ENV → `pg-change-feed: Fehlerklasse configuration: Verdrahtung ohne vollständige Vorbedingung: CDC_SOURCE_DSN fehlt`, **Exit 2** (Klasse `configuration` am ENV-Verweigerungspfad, [`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md)) | 0/2 |
| `docker image inspect ghcr.io/pt9912/pg-change-feed:dev` | Image-Id = **`sha256:899f39541e1939c6a8b10934afd6e179173048d9d38986f77e20e2a5ed6d890d`** — exakt der am HEAD committete Lauf-Beleg (`harness/image-hash.txt`, `5dc855a`) | 0 |
| Mutations-Probe 1 (`if cfg.DSN == ""`-Verweigerungspfad entfernt; `/tmp`-Worktree, revertiert und entfernt) | `--- FAIL: TestConfigFromEnvOhneVorbedingung … CDC_SOURCE_DSN fehlt: Verdrahtung startet ohne Vorbedingung` — **erwartet rot, rot gesehen** | 1 |
| Mutations-Probe 2 (Bindungs-Form-Prüfung in `parseTables` entfernt) | `--- FAIL: TestConfigFromEnvTabellenForm … CDC_TABLES "public.t1=tbl-1": Verdrahtung startet trotz Formverletzung` — **erwartet rot, rot gesehen** | 1 |
| Rote Probe commit-traceability ([`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md)): Probe-Commit ohne Kennung (`177670c`) in `/tmp`-Worktree, beide Sensor-Hälften über `e971ed4..177670c` | d-check-Modul `commits`: `177670c:1 … commit-untraceable probe: Commit ohne Vertrags-Kennung` · **d-check Exit 1** · Grenz-Hälfte (`commit-traceability.sh`) korrekt OK (keine Struktur-ID im Betreff) — das Standing-Gate greift an beiden Hälften gemäß seiner Teilung | 1 |
| End-Wächter-Probe (`a184663`-Logik gegen abwesenden Feed-Container) | `Feed-Container endete während des Testlaufs (Lauf: false, Ausgang: fehlt)`, Exit 1 — **erwartet rot, rot gesehen** | 1 |
| `make doc-commits RANGE=de340b6..HEAD` | `d-check: 103 Datei(en) geprüft, 0 Befund(e)` — Traceability je Commit über die **volle** 14-Commit-Range (das Standing-Gate-Fenster `HEAD~5..HEAD` deckt `ead5c01`/`62fc9b9` nicht; diese Run-Deckung schließt die Lücke) | 0 |
| `make doc-immutable RANGE=de340b6..HEAD` | `d-check: 103 Datei(en) geprüft, 0 Befund(e)` (Immutabilität) | 0 |
| Datei-Abschluss (`tail -c1 | od -c`) über alle sieben im Fix-Zug berührten Dateien | alle sieben enden auf `\n` — F-3-Klasse im Range geschlossen | — |
| `git worktree list` / `git status --porcelain` | nur der Hauptbaum; Arbeitsbaum clean — keine Einträge dieses Laufs | — |

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt — 10 Items)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | `cmd/pg-change-feed` verdrahtet die reale Pipeline (Store, Stream, Service, ACK) je [`ADR-0026`](../plan/adr/0026-composition-root.md) — der Feed-Container fährt CDC-Runtime statt `--version`-Smoke | **bestätigt** | `wiring.go` verdrahtet `postgresstorage.New` → `receive.NewStream` → `postgresack.New(stream.Conn())` → `BindCapture(capture.NewCaptureService(...))` → `stream.Run` an genau einer Stelle; `main` importiert nur `internal/bootstrap` (kein Adapter-Konstruktor, a-check 0 Befunde). Verhaltens-Probe am Image: ohne ENV Exit 2 über `bootstrap.ErrConfiguration` — das Binary trägt die Verdrahtung; **CDC-Runtime am Container:** manueller E2E-Probe (Tabelle oben) — `cdc.change` wird vom Feed-Container gefüllt, Slot aktiv |
| 2 | MVP-Integrationstest (slice-006) fährt gegen den verdrahteten Feed-Container: INSERT/UPDATE/DELETE landen im Store (Ende-zu-Ende durch das Binary) | **bestätigt** | `make test-integration` grün am HEAD (beide MVP-Tests PASS); der Test verdrahtet keinen Adapter (`mvp_test.go` importiert nur den Store-Adapter-Lese-Pfad), liest die Port-Kennung aus den CDC-Referenztabellen und verlangt den Slot als Vorbedingung des Container-Starts; „kein anderer Schreiber als das Binary" durch den manuellen Probe bestätigt (Aktivierung schreibt nur Referenztabellen) |
| 3 | `make gates` grün | **bestätigt** | Vier Gates grün am HEAD `e971ed4` (Tabelle oben), inkl. `commit-traceability` ([`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md)-Standing-Gate, 5 Commits OK); Häkchen in §2 noch offen (Closure-Pflicht) |
| 4 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **bestätigt** | `docs/reviews/review-slice-007.md` committet in `e971ed4` (Übergabe-Artefakt, Message trägt `LH-QA-OPS-001`/`ADR-0026`/`ADR-0044`); Rollenwechsel nach Schritt 8 eingehalten |
| 5 | Doku-Update für <Schnittstelle X> falls öffentlicher Vertrag berührt | **nicht bestätigt — V-2** | Der ENV-Container-Vertrag (`CDC_*`) ist de-facto-Vertragsform (als benannter Vertrag an drei Orten deklariert: `compose.yaml`, Runner, Test), aber keine öffentliche Doku-Stelle trägt ihn (grep in `harness/` ohne Treffer); die `harness/README.md:129`-Werkzeuge-Zeile zu `test-integration` beschreibt die slice-006-Kette („frisch hochfahren → Rollout → Toolchain-Container") **ohne** Aktivierungs-Schritt und ohne die CDC-Runtime-Form des Feed-Containers. Item-Häkchen offen |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **offen bis Closure (notiert)** | §7 trägt Platzhalter; Slice korrekt in `in-progress/` — fällig vor dem `git mv`. Lerneintrag-Kandidaten unten (Datei-Abschluss-Klasse, F-2-Ausgang ins Register) |
| 7 | Reconciliation-Register fortgeschrieben (falls Inventur-Fund) | **entfällt** | `docs/plan/planning/reconciliation.md` existiert nicht (Repo ohne Brownfield-Bootstrap) — vom DoD-Wortlaut ausdrücklich vorgesehen; kein Inventur-Fund im Range |
| 8 | Beobachtungs-Register fortgeschrieben | **kein neues Auftreten — notiert (Closure-Antwort fällig)** | `BEO-PGC/a-check-null-abdeckung` trägt den Ausgang `verkörpert · seit welle-1` (evidence slice-001/002/003, unverändert); `BEO-PGC/d-migrate-nacharbeit` offen 1× (evidence slice-006) — kein neues Auftreten, **keine evidence-Datei für slice-007 angelegt**. Register-Kandidaten für die Closure: die **Datei-Abschluss-Klasse** (Review F-3, 4.–7. Auftreten im Range, durch `62fc9b9`/`a184663`/`dae053d` geschlossen — aber die Klasse hat die 3×-Schwelle längst überschritten) und die **Commit-Kennung-Klasse** (durch [`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md) verkörpert — Ausgang `verkörpert`, Herkunfts-Anker `seit slice-006`/`seit slice-007` je Wortlaut der Verkörperung; die rote Probe oben belegt das Gate) |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **Voreinstellung getragen — Register-Belege fällig bei Closure** | (a) Verdrahtungs-Lücke — **eingetreten**, Träger ist dieser Slice (belegt: DoD 1/2). (b) Fehlerbehandlungs-Grenze — **weiter offen** (Träger: Kommentare `wiring.go`/`compose.yaml` aus `62fc9b9`; Retry/Backoff „folgt mit der Konfigurationsschicht (späterer Slice)") — bei Closure in das Beobachtungs-Register zu überführen; der genannte Träger trägt **keine slice-Kennung** (F-7-Klasse, V-4). (c) MVP-Claim an der Go-Baum-Ebene — **eingetreten**, Träger ist dieser Slice (DoD 2). Alle drei Ausgänge sind von der geschlossenen Menge; die Endbelege (Register-Zitat für (b)) fallen bei der Closure an |
| 10 | Drei Paarungen (Anker · Folge-Slice · Register) | **delegiert — korrekt** | Repo **mit** Wellen-Betrieb (`welle-2.md` flach vorhanden; der Slice selbst ist wellenlos): der DoD-Wortlaut weist die Prüfung der nächsten Welle-Closure zu; `welle-2.md` §3 verlangt den Rest-Verdrahtungs-Ausgang **vor** der Welle-Closure — genau dieser Ausgang ist durch diesen Slice getragen (§6 (a) eingetreten), die Welle-2-Closure kann nach der Slice-Closure folgen |

**Zwischenstand: 4/10 Punkte jetzt erfüllt** (Items 1–4, Belege selbst
gefahren) · 1 nicht bestätigt (Item 5, V-2) · 1 entfällt (Item 7) · 3
erfüllen sich erst bei Closure (Items 6, 8 als Notiz, 9-Endbelege) · 1
delegiert (Item 10). Der Zustand ist der **normale eines Slice in
`in-progress/`** — alle §2-Häkchen stehen noch offen; die offenen Items sind
Closure-Pflichten. V-2 und V-3 sind dagegen **vor** der Closure zu tragen.

## Plan-vs-Code-Diff (Range `de340b6..5dc855a`, Plan-Nachzug `ead5c01`)

**Plan-Nachzug (`ead5c01`) gegen den Review-Blocker F-1:** bestätigt — §3
trägt jetzt beide nachgezogenen Zeilen (`tools/harness/run-integration-tests.sh`
mit Aktivierung **vor** dem Container-Start und zweigeteiltem Wächter;
`test/integration/mvp_test.go` als Neuschreiben mit Store-Adapter-Lese-Pfad).
Die Message nennt die Vertrags-Kennungen ([`ADR-0026`](../plan/adr/0026-composition-root.md)/0023/0045,
`LH-QA-POR-003`). Im Gegensatz zum slice-005-Nachzug (dort Rest V-2) ist der
Nachzug diesmal **vollständig**.

**Deckung §3 (am HEAD):** alle im Range gelieferten Produkt-Dateien sind durch
die (nachgezogenen) §3-Zeilen gedeckt — `cmd/pg-change-feed/main.go`,
`internal/bootstrap/*.go` (deckt `wiring.go` **und** `wiring_test.go` über den
Glob), `compose.yaml`, `tools/harness/run-integration-tests.sh`,
`test/integration/mvp_test.go`. Kein Rest.

**Nicht in §3, zulässig:** `harness/image-hash.txt` (Lauf-Beleg,
[`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)-Semantik, zwei
Digest-Commits `789b76e`/`5dc855a`) · Plan-`git mv` open→in-progress
(Planner, reine Move-Kette) · Review-Report `docs/reviews/review-slice-007.md`
in `e971ed4` (Übergabe-Artefakt, nach dem Range-Schluss) · Fix-Commits
berühren nur die fünf §3-Dateien plus Beleg.

## ADR-Konformität

- **[`ADR-0026`](../plan/adr/0026-composition-root.md)** (Composition Root, `main` dünn): **konform** — `main.go` importiert nur `internal/bootstrap` (plus Std-Lib), trägt ENV-Lesen, Signale und Prozess-Ausgang; `bootstrap.Run` kennt die konkreten Adapter und verdrahtet an genau einer Stelle; `a-check` 0 Befunde im `make gates`-Lauf.
- **[`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md)** (Lauf-Beleg): **konform, diesmal trägt die Beleg-Kette** — der Range ändert Build-Kontext-Dateien (`cmd/`, `internal/` — `COPY . .`), `make image` lief vor dem Digest-Commit, und **das Binary trägt wirklich die Verdrahtung** (CDC-Runtime statt `--version`-Smoke): Digest-Kette `9ac4a9fb…` (Range-Basis `de340b6`) → `ee795a88…` (`789b76e`) → `899f3954…` (`5dc855a`, am HEAD). **Eigener Beleg:** `docker image inspect` des geladenen `:dev`-Images liefert exakt `sha256:899f3954…`, und das Image verhält sich als verdrahtetes Binary (`0.2.0-verdrahtung`; fehlende ENV → Exit 2, Klasse `configuration`).
- **[`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md)/[`SPEC-008`](../../spec/pflichtenheft.md)** (Fehlerklassen): **konform** — der ENV-Verweigerungspfad endet über `bootstrap.ErrConfiguration` (Exit 2 am Prozess, belegt am Image); Adapter-Fehler laufen klassen-ge Sentinel-ge (`ErrConfiguration`/`ErrReplication` je `receive`, `outbound.ErrStorage` am Store) auf Ausgang 1; meine Verhaltens-Probe zeigt die Klassen-Zeile ohne Credentials (ENV-Namen, nie Werte). Die Klassen-**Aktionen** (Retry, kontrollierte Fortsetzung) bleiben bewusst ungetragen — Grenze benannt (62fc9b9, siehe [`ADR-0012`](../plan/adr/0012-at-least-once.md)).
- **[`ADR-0012`](../plan/adr/0012-at-least-once.md)** (Neustart-Vertrag beim Aufrufer): **Grenze getragen** — `62fc9b9` benennt beide Seiten indikativ: `wiring.go` („der Prozess-Aufrufer endet auf jeden Adapter-Fehler mit Ausgang 1 — die Fortsetzung … trägt der Prozess-Neustart, der Slot liest seinen Start über confirmed_flush_lsn") und `compose.yaml` („Kein Neustart im Container-Vertrag … der Neustart-Vertrag liegt beim Aufrufer"). Die `transient`-Aktion ist als ungetragen benannt — derselbe Ausgang wie Plan §6 (b). **Konform als benannte Grenze** (Review-F-2-Verdikt-Punkt 2 erfüllt).
- **[`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md)** (Standing-Gate): **greift** — im `make gates`-Bündel am HEAD grün (5 Commits, `HEAD~5..HEAD`); **rote Probe repliziert**: ein Probe-Commit ohne Vertrags-Kennung färbt die d-check-Hälfte rot (`commit-untraceable`, Exit 1); die Grenz-Hälfte (keine Struktur-ID im Betreff) bleibt für denselben Commit korrekt grün — die Teilung der Hälften trägt genau wie im ADR beschrieben. Fenster-Grenze benannt: `ead5c01`/`62fc9b9` liegen außerhalb des 5-Commit-Fensters und wurden von mir über `doc-commits RANGE=de340b6..HEAD` (0 Befunde) gedeckt.
- **[`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md) Exit-Code-Form:** konsistent — `ConfigFromEnv`-Fehler → Ausgang 2, Adapter-Fehler → Ausgang 1 (`main.go:34-41`); die Kommentare in `wiring.go` benennen die Zuordnung; die Probe am Image belegt Ausgang 2 am ENV-Pfad.

## Fix-Commits gegen die Review-Findings

| Finding | getragen im Fix? | Beleg |
|---|---|---|
| F-1 (Plan-Nachzug Runner/Test) | **ja, vollständig** | `ead5c01` — beide §3-Zeilen, Message mit Vertrags-Kennungen; Deckung geprüft (oben) |
| F-2 (Fehlerbehandlungs-Grenze unbenannt) | **ja** | `ead5c01` (§6-Risiko mit Ausgang „weiter offen") + `62fc9b9` — Grenz-Kommentare in `wiring.go`/`compose.yaml`, indikativ („endet …", „liegt beim Aufrufer"), Kommentar-Klassen Zusage/Abgrenzung; die `SPEC-008`-`transient`-Aktion ist als ungetragen benannt, nicht verworfen |
| F-3 (Datei-Abschluss, 4 Dateien) | **ja** | `62fc9b9` (wiring.go, compose.yaml) · `a184663` (Runner) · `dae053d` (main.go, mvp_test.go) — `tail -c1`-Beleg: alle sieben berührten Dateien enden auf `\n` |
| F-4 (totes `PUBLICATION`) | **ja** | `a184663` — Variable entfernt; der Kommentar trägt nur noch den gelesenen `SLOT` (im Start- und End-Wächter genutzt); der Publication-Name trägt das Aktivierungs-SQL allein |
| F-5 (End-Beleg des Feed-Containers) | **ja** | `a184663` — End-Wächter (`State.Running` + `ExitCode` nach dem Testlauf); eigene Probe: abwesender Container → Exit 1 (Tabelle oben) |
| F-6 (`ConfigFromEnv`/`parseTables` ohne Negativtests) | **ja, mit Rest V-1** | `10953e2` — `wiring_test.go` (117 Zeilen): vier fehlende ENVs, Grenze ohne Aktivierung, sechs Form-Verletzungen, vollständiger Lese-Pfad; **zwei Mutationen rot repliziert** (DSN-Guard, Bindungs-Form — Tabelle oben); **aber** die Datei ist nicht gofmt-clean (V-1) |
| F-9 (Plan §4/§5/§8 Platzhalter) | **nicht getragen — V-3** | §4-Rückführungen (`<Bedingung>`), §5-Closure-Trigger (`<…>`) und §8-Sichtungs-Schritt (Stub) stehen am HEAD unverändert; nur §6 (Risiko-Ausgänge) wurde in `ead5c01` gefüllt. §5 ist Closure-Voraussetzung (Review-Verdikt) |
| F-7/F-8/F-10 (Planner-Residuen / INFO) | **offen, korrekt platziert** | §1-Ausschlüsse weiterhin ohne `slice-<NNN>`-Kennung (F-7); Trigger-Wortlaut §4 unverändert (F-8 — materiell folgenlos für den wellenlosen Zug); Signal-Ende-Vertrag weiter ohne automatisierten Träger (F-10 → V-5) |

## Befunde

### V-1 — `wiring_test.go` ist nicht gofmt-clean (Fix-Commit `10953e2`)

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Format-Abweichung im Fix-Zug) ·
  `gofmt -l`-Beleg im gepinnten Toolchain-Container
- `pfad`: `internal/bootstrap/wiring_test.go` — zwei Ausrichtungs-Blöcke
  (ENV-Map-Literal, Form-Verletzungs-Liste; `gofmt -d`-Beleg oben)
- `befund`: Die Negativtest-Datei des F-6-Fixes endet sauber auf `\n`, ist
  aber unformatiert — `go build`/`go vet` sind clean, `gofmt -l .` meldet
  genau diese Datei. Kein Gate im Bündel fährt `gofmt` (die vier
  `make gates`-Gates prüfen nicht das Format), daher bleibt der Defekt
  sensor-stumm — er ist von mir im Toolchain-Container belegt, nicht von
  einem Gate. Mechanisch reparabel, keine Vertragsverletzung.
- `verifizierbar`: ja — `gofmt -l .` / `gofmt -d` im gepinnten Container
- `klasse`: Format-Abweichung im Fix-Zug (Sensor-Lücke: kein gofmt-Gate)

### V-2 — Doku-Aussage zum Container-Vertrag unvollständig: `harness/README`-Zeile veraltet, ENV-Vertrag ohne öffentliche Stelle (DoD-Item 5)

- `kategorie`: MEDIUM
- `quelle`: Slice-Plan §2 DoD-Punkt 5 · `harness/README.md` §Sensors (Werkzeuge-Zeilen) · Review F-7 (ENV-Vertrag ohne Kompatibilitäts-Pflicht-Stelle)
- `pfad`: `harness/README.md:129` (`test-integration`-Zeile: „Umgebung frisch
  hochfahren → Schema-Rollout über d-migrate → Toolchain-Container gegen das
  Compose-Netz") · `compose.yaml` (ENV-Vertrag `CDC_*`)
- `befund`: Die Werkzeuge-Zeile beschreibt die slice-006-Kette und kennt weder
  den Aktivierungs-Schritt vor dem Feed-Container-Start noch die neue
  Vertragsform (Feed-Container als CDC-Runtime, Start-/End-Wächter); der
  ENV-Vertrag (`CDC_SOURCE_DSN`…`CDC_TABLES`) ist als Container-Vertrag
  deklariert (drei Code-Stellen nennen `compose.yaml` als Vertrag), aber
  keine öffentliche Doku-Stelle trägt ihn. Der DoD-Punkt 5 (Doku-Update bei
  berührtem Vertrag) ist damit **nicht getragen** — vor der Closure entweder
  die Zeile nachziehen oder die begründete Aussage „kein öffentlicher Vertrag"
  explizit gegen diese beiden Stellen fassen.
- `verifizierbar`: ja — `harness/README.md:129` gegen
  `tools/harness/run-integration-tests.sh` (Kette) · grep `CDC_SOURCE_DSN`
  in `harness/` ohne Treffer
- `klasse`: Doku-Vertrags-Aussage unvollständig (ENV-Vertrag ohne Adresse —
  F-7-Klasse läuft weiter)

### V-3 — F-9 nicht getragen: Plan §4/§5/§8 stehen weiterhin als Platzhalter

- `kategorie`: MEDIUM
- `quelle`: `review-slice-007.md` F-9 (INFO, mit Closure-Voraussetzung) ·
  Modul 5 (Rückführungen „vorab benennen"; Closure-Trigger-Pflicht) ·
  §8-Sichtungs-Schritt (einziger Register-Leser im Planungs-Zug)
- `pfad`: `docs/plan/planning/in-progress/slice-007-bootstrap-verdrahtung.md`
  §4 (zwei `<Bedingung>`-Zeilen) · §5 (`<…>`) · §8 (Sub-Area-Wahl- und
  Sichtung-Stand, Stub-Block)
- `befund`: Der Fix-Zug trägt F-1 (§3-Nachzug) und F-2 (§6-Ausgang), aber
  F-9 ist nicht angegangen: Der Closure-Trigger (§5) steht als `<…>` — er ist
  **Closure-Voraussetzung** und muss vor dem `git mv` gefüllt werden; die
  §8-Sichtungs-Prüfung (offene Beobachtungen gegen das Register) ist nicht
  aufgezeichnet. Planner-Arbeit, keine Implementer-Arbeit — der Auftrag
  nannte F-9 unter den getragenen Findings; dem ist nicht so.
- `verifizierbar`: ja — Lese der Plan-Datei am HEAD (Beleg oben)
- `klasse`: Vorlagen-Platzhalter im eröffneten Slice-Plan (F-9 läuft weiter)

### V-4 — Der „weiter offen"-Ausgang der Fehlerbehandlungs-Grenze braucht bei Closure seinen Register-Beleg; Träger ohne Kennung

- `kategorie`: LOW
- `quelle`: Modul 5 (Risiko-Ausgänge: „weiter offen" → Beobachtungs-Register)
  · Modul 5 §Ziel-Form (Klasse 1: „Ein Folge-Slice übernimmt es — mit
  Kennung")
- `pfad`: Slice-Plan §6 (b): „Ausgang: weiter offen → Retry/Backoff folgt mit
  der Konfigurationsschicht (späterer Slice)"
- `befund`: Der Ausgang ist von der richtigen geschlossenen Menge gewählt
  („weiter offen"), aber der genannte Träger („späterer Slice") trägt keine
  `slice-<NNN>`-Kennung, und der Register-Beleg (Verzeichnis
  `BEO-<KUERZEL>/<slug>/` mit `evidence/slice-007.md`) ist noch nicht
  angelegt — beides fällt bei der Closure an (Item 8/9-Endbelege). Notiert,
  nicht blockierend über die ohnehin offene Closure hinaus.
- `verifizierbar`: ja — Plan §6 gegen `docs/plan/planning/observations/`
- `klasse`: Ausschluss-/Träger-Klasse 1 ohne Adresse (F-7-Klasse, 2. Stelle
  in diesem Plan)

### V-5 — Signal-Ende-Vertrag weiter ohne automatisierten Träger (F-10 läuft weiter)

- `kategorie`: INFO
- `quelle`: `review-slice-007.md` F-10 · Commit `3c1de11` („Der Lauf endet
  kontrolliert auf SIGINT/SIGTERM")
- `pfad`: `cmd/pg-change-feed/main.go:28-32` (`signal.NotifyContext`) — kein
  Test sendet ein Signal (grep in Testdateien ohne Treffer)
- `befund`: Die Zusage bleibt dünn besetzt; die Implementer-Begründung
  („Binary-Lauf-Test-Aufbau wächst über den Nachzug") ist tragfähig — der
  Aufbau eines Binary-Lauf-Tests (Signal → Ausgang 0) ist ein eigener Zug.
  Grenze notiert; Beleg-Kandidat für einen künftigen Binary-Lauf-Test.
- `verifizierbar`: ja — Test mit Signal an das Binary (Ausgang 0)
- `klasse`: Zusage am Prozess-Rand ohne Test-Träger (INFO-Kanal, benannt)

## Negativbefunde

- geprüft, ohne Befund: **Docker-only (AGENTS §3.1)** — alle go-Belege im
  gepinnten Toolchain-Container (`--network none` für Build/Vet/`make test`,
  Modul-Cache im Volume, `/src:ro`), Integrationstest über `make
  test-integration`; kein Host-Toolchain-Aufruf, kein Schreibzugriff auf den
  Hauptbaum
- geprüft, ohne Befund: **Das Binary fährt die reale Pipeline (Prüfpunkt a)** —
  der manuelle E2E-Probe zeigt die CDC-Tabellen-Befüllung durch das Binary im
  Feed-Container: 2 INSERT-Zeilen in `cdc.change` mit korrektem
  `new_data`-Image, Slot aktiv mit `confirmed_flush_lsn`, Container
  `running exit=0`; der Runner-Chain-Beleg (`make test-integration` grün)
  belegt dieselbe Kette mit Start- und End-Wächter
- geprüft, ohne Befund: **ENV-Verweigerungspfade (Klasse `configuration`,
  Exit 2)** — Verhaltens-Probe am geladenen Image (ohne ENV → Klassen-Zeile
  mit ENV-Namen, Exit 2); zwei Mutations-Proben rot (DSN-Guard, Bindungs-Form)
  — die Zusage aus `10953e2` trägt
- geprüft, ohne Befund: **End-Wächter (F-5-Fix)** — Logik-Probe gegen
  abwesenden Container rot (Exit 1); im grünen Integrationslauf passierte der
  Wächter sichtbar (Script-Ende Exit 0)
- geprüft, ohne Befund: **Image-Beleg-Kette ([`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md))** —
  lokale Image-Id == committeter Beleg `899f3954…`; Dockerfile-Pins
  unverändert (kein FROM-Wechsel im Range); der Digest-Wechsel deckt sich mit
  der Binary-Änderung (Verdrahtung statt Stub)
- geprüft, ohne Befund: **Traceability aller 14 Commits der Range
  `de340b6..HEAD`** — `doc-commits` 0 Befunde; die drei Planner-Commits
  (`c537b74`/`6d0cd35`/`bb7b6b7`) reine `git mv`-Kette (Review-Beleg);
  `doc-immutable` 0 Befunde über die volle Range
- geprüft, ohne Befund: **superseded-Referenzen** — die Referenzen des Ranges
  nennen nur Accepted-ADRs (0007/0012/0023/0026/0044/0045)
- geprüft, ohne Befund: **Kommentar-Klassen (§3.7) im Fix-Zug** — die neuen
  Grenz-Kommentare in `wiring.go` und `compose.yaml` sind indikativ über den
  geltenden Zustand (Klasse Zusage/Abgrenzung), kein Konjunktiv über
  verworfene Alternativen; der compose-Block nennt den geltenden
  Neustart-Vertrag, nicht die Chronik
- geprüft, ohne Befund: **Runner-Hygiene** — `down -v` in jedem Ausgang
  (trap EXIT) und vor dem Start, PostgreSQL-Readiness, `-v ON_ERROR_STOP=1`
  am Aktivierungs-SQL, gepinnte Digests konsistent; meine Compose-Instanz
  nach Abschluss abgeräumt (keine `cdc-test-*`-Rückstände)
- geprüft, ohne Befund: **§1-Abgrenzung** — keine Config-Format-Festlegung
  (ENV-Minimal-Form, TOML/YAML offen), keine CLI-/SQL-Adapter, keine
  Observability-Metriken im Range
- geprüft, ohne Befund: **WIP-Limit 1** — nur slice-007 in `in-progress/`

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Format-Abweichung im Fix-Zug (gofmt,
Sensor-Lücke: kein gofmt-Gate) · Doku-Aussage zum Container-Vertrag
unvollständig (ENV-Vertrag ohne öffentliche Adresse — F-7-Klasse läuft
weiter) · Vorlagen-Platzhalter im Slice-Plan (F-9 läuft weiter — §5
Closure-Trigger fällig) · Folge-Slice-Adresse ohne Kennung (F-7-Klasse,
2. Stelle) · Signal-Ende ohne Test-Träger (F-10 läuft weiter, INFO).

**Zusammenfassung DoD:** **4/10 Punkte jetzt erfüllt** (Items 1–4, Belege
selbst gefahren: `make gates` grün am HEAD, 10/10 Test-Pakete netzlos,
`make test-integration` grün am verdrahteten Feed-Container, manueller
E2E-Probe belegt die `cdc.change`-Befüllung durch das Binary, Image-Beleg
`899f3954…` lokal verifiziert und verhaltens-belegt) · Item 5 **nicht
bestätigt** (V-2) · Item 7 entfällt · Items 6/8/9 erfüllen sich bei Closure
(Register-Kandidaten: Datei-Abschluss-Klasse ≥3×, Commit-Kennungs-Klasse
durch [`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md) verkörpert) · Item 10 delegiert an
die Welle-2-Closure (deren §3-Rest-Verdrahtungs-Ausgang dieser Slice trägt).

**Alle Review-Blockierpunkte sind beantwortet:** F-1 (Nachzug `ead5c01`,
vollständig), F-2 (Ausgang + Grenz-Kommentare `62fc9b9`), F-3 (alle vier
Dateien + Fixtur-Dateien schließen sauber), F-4/F-5 (`a184663`, End-Wächter
durch meine Probe rot-geprüft), F-6 (`10953e2`, zwei Mutationen rot
repliziert). Die rote Probe gegen das [`ADR-0045`](../plan/adr/0045-commit-traceability-standing-gate.md)-Gate
(Probe-Commit ohne Kennung → `commit-untraceable`, Exit 1) belegt, dass das
Standing-Gate am HEAD greift.

## Verdikt

**Merge-blockierend:** nein — die Verdrahtung ist am realen Pfad belegt
(Integrationstest grün, manueller E2E-Probe zeigt die Befüllung von
`cdc.change` durch das Binary im Feed-Container), das Image trägt den
verdrahteten Binary-Stand mit dem am HEAD committeten Lauf-Beleg
(`899f3954…`, lokal verifiziert), die Fix-Commits tragen F-2/F-3/F-4/F-5/F-6
nachweisbar (inkl. replizierter roter Proben), und alle vier Gates laufen
grün am HEAD.

**Blockierend für Closure (der normale Zustand eines Slice in
`in-progress/`, plus zwei echte Punkte):**

1. **V-2:** DoD-Punkt 5 vor dem Häkchen tragen — `harness/README.md:129`
   (Ketten-Beschreibung ohne Aktivierungs-Schritt und CDC-Runtime) nachziehen
   **oder** die begründete Aussage „kein öffentlicher Vertrag berührt" explizit
   gegen die ENV-Vertrags-Stellen fassen. Der ENV-Vertrag ist als benannter
   Vertrag an drei Code-Stellen deklariert; eine öffentliche Adresse braucht
   er trotzdem (F-7-Klasse).
2. **V-3:** Plan §5 (Closure-Trigger) füllen — Closure-Voraussetzung;
   §4-Rückführungen und §8-Sichtungs-Schritt gehören davor (F-9, Planner).
3. **V-1:** `wiring_test.go` gofmt-clean stellen (mechanisch, vor der
   Closure); die Sensor-Lücke (kein gofmt-Gate) als Register-/Lerneintrag-
   Kandidat führen.
4. **Closure-Pflichten (Items 5-Häkchen, 6, 8, 9-Endbelege):** Häkchen,
   Closure-Notiz mit Lerneintrag, Register-Belege (Datei-Abschluss-Klasse
   ≥3×; „weiter offen"-Ausgang der Fehlerbehandlungs-Grenze mit
   Register-Eintrag — V-4; Signal-Ende-INFO V-5), Risiko-Ausgänge je §6,
   Paarungen an die Welle-2-Closure delegiert (Item 10).

**Übergabe:** Bericht an den Planner. Keine Reparaturen. Review-Verdikt
geprüft: alle drei blockierenden Punkte (F-1, F-2, F-3) sind getragen; die
F-4–F-6-Vor-Closure-Nacharbeit ist getragen (F-6 mit Rest V-1).