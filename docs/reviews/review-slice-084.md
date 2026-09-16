# Review-Report: slice-084 — 2026-09-16

**Review-Art:** Code — geprüft gegen **Plan und Entscheidungen**
(Baseline-Regelwerk `v6.5.0` · `regelwerk/modul-10-review-harness.md`
§Drei Review-Arten). DoD-/Spec-Konformität (Verifier), §6-Risiko-Ausgänge,
Beobachtungs-Register und die drei Paarungen (Planner-Closure) sind **nicht**
Gegenstand dieses Reports.

**Gegenstand:** `slice-084` (`docs/plan/planning/in-progress/slice-084-postgresack-naht.md`),
Diff `fb6adf6..4035ee7` (vier Commits: `eecb1d9` Naht · `8e53d2a` Plan-Nachzug
und Null-Befund · `26d4ef0` Digest-Beleg · `4035ee7` §3-Zeile). Fünf Dateien:
`internal/adapters/driven/postgresack/{ack.go,seam.go,seam_test.go}`,
`harness/image-hash.txt`, der Slice-Plan selbst. **Kein** anderer Träger ist
berührt — `Dockerfile`, `tools/harness/db-coverage.sh`, `harness/mk/coverage.mk`,
`internal/bootstrap`, `.a-check.yml`, `spec/**` sind diff-frei (eigene Prüfung,
§Eigene Messungen).

**Skill:** `.harness/skills/reviewer.md` @ Repo-Stand (vier repo-spezifische
HIGH-Regeln, letzte Schärfung 2026-09-09) · **Modell:**
deepseek-v4.1-flash:cloud · **Datum:** 2026-09-16.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-084-postgresack-naht.md`
  vollständig (§1–§8)
- [`ADR-0080`](../plan/adr/0080-nahtform-pgconn-adapter-treiberhuelle.md)
  (Nahtform: Hülle + fake-fähige Fläche, **im Paket**; Festlegung 1/4/5, die
  „benannte Grenze", Re-Evaluierungs-Trigger (a)/(e))
- [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
  Punkt 3/4/5 · [`ADR-0078`](../plan/adr/0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md)
  (dreiteiliger Nachweis, hier als Null-Befund) ·
  [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md) (Digest ist
  Lauf-Beleg, `harness/image-hash.txt` ist sein Träger) ·
  [`ADR-0005`](../plan/adr/0005-sourceposition-abstrahiert-lsn.md)
  (SourcePosition/LSN) · [`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md),
  [`ADR-0032`](../plan/adr/0032-postgresql-adapterdetail.md)
- `AGENTS.md` §3.1, §3.5, §3.6, §3.7, §3.9, §3.11, §5 · `harness/conventions.md`
  `MR-000` (ID-Schema)
- Beobachtungs-Register `docs/plan/planning/observations/BEO-PGC/`
  (namentlich `zahl-in-traeger-driftet-gegen-die-messung`,
  `endstufe-unter-eigenem-messgegenstand-unerreichbar`,
  `adapter-unittest-verdeckt-bootstrap-luecke`,
  `kommentar-behauptet-nicht-getragenen-fehlerpfad`) ·
  `docs/reviews/review-slice-081.md` (Geschwister-Naht, vorherige Findings am
  gleichen Modul)
- `harness/sensors/db-adapter-coverage.md`, `harness/sensors/coverage-gate.md`,
  `tools/harness/db-coverage.sh`, `Dockerfile` (Stufe `coverage`)

---

## Findings

### F-1 — Die Naht macht `harness/sensors/db-adapter-coverage.md` §Zählbasis falsch: 650/23 sind gemessen 659/32

- `kategorie`: LOW
- `quelle`: Maintainability · `AGENTS.md` §3.7 („Ein Kommentar beschreibt, was
  da ist") · Klasse `zahl-in-traeger-driftet-gegen-die-messung` (Beobachtungs-Register,
  1× aus `review-slice-081` — dieser Vorgang ist der **zweite**)
- `pfad`: `harness/sensors/db-adapter-coverage.md:65` („Der gemergte Nenner ist
  **650 Statements** (`postgresstorage` 472 · `postgresack` 23 ·
  `replication/receive` 155)") gegen `:79` („real gemessener Ist-Stand
  **73,38 %** (477 von 650 Statements)")
- `befund`: Der DB-Gegenstand trägt nach diesem Diff **659** Statements mit
  **491** gedeckten (74,51 %), `postgresack` **32** statt 23 — eigene Messung
  des gemergten Profils (§Eigene Messungen). Das Sensor-Dokument ist der
  lebende Vertrag über genau diesen Gegenstand und trägt die alten Zahlen
  weiter, ohne Datierung und ohne „am Stand"; der Zug hat sie bewegt und
  `harness/sensors/**` nach §3 ausdrücklich **nicht angefasst**. Die Kategorie
  bleibt **LOW**, nicht HIGH wie `review-slice-081` F-1: dort stand der Träger
  in der **deklarierten** Datei-Liste des Slice und wurde nur an der Zahl
  daneben liegen gelassen; hier steht er in „**Nicht angefasst**", die
  Drift ist also Folge der erklärten Grenze, nicht eines unvollständigen
  Zugriffs — aber auch nicht INFO wie `review-slice-081` F-6, weil dieser
  Träger ein **lebendes Sensor-Dokument** mit Gegenwarts-Aussagen ist, das
  dieser Diff nachweislich falsch gemacht hat. Die Zeilen `:159-162`
  (§Rot-/Grün-Beleg, „477 von 650 … 73.38 %") sind dagegen als **zitierter
  Einzellauf** korrekt eingefroren und nicht betroffen.
  *Nebenbefund, nicht diesem Vorgang zuzurechnen:* `coverage-gate.md:56`
  nennt denselben Zustand als „1831 Statements, davon 1306 gedeckt = 71,33 %";
  der Gegenstand dieses Baums trägt **1903** (1369 gedeckt, 71,94 % dedup,
  von der Stufe 71.90 % gedruckt). Der Unit-Nenner ist von diesem Diff
  **unberührt** (die Naht liegt im ausgenommenen Paket, eigenes Messprofil
  enthält keine `postgresack`-Zeile) — die Abweichung stammt aus Vorgängen
  zwischen `slice-081` und diesem Zug und wird hier nur benannt, damit der
  Zug nicht als „alle Träger nachgezogen" gelesen wird.
- `verifizierbar`: ja — eigenes `-coverprofile` des Tier-Laufs, per
  `awk` über Block-Position dedupliziert (§Eigene Messungen, Zeile
  „DB-Adapter-Coverage gemergt")

### F-2 — §3(b) zitiert als Beleg seinen **Befehl** nicht mit: `git diff --name-only` trägt die Aussage „isoliert auf ein Paket" nicht

- `kategorie`: LOW
- `quelle`: Maintainability (Reproduzierbarkeit des Belegs) — Geschwister der
  Klasse aus `docs/reviews/review-slice-081.md` F-5 („Der in §3(b) zitierte
  Beleg-Befehl löst nach dem Landen nicht mehr auf")
- `pfad`: `docs/plan/planning/in-progress/slice-084-postgresack-naht.md:206-209`
- `befund`: Der Plan belegt den Null-Befund (b) mit „(`git diff --name-only`
  listet ausschließlich Dateien unter `internal/adapters/driven/postgresack/`)".
  Ohne Argumente zeigt der Befehl den **Arbeitsbaum** gegen den Index — auf
  sauberem Baum **nichts** (mit liegengebliebenem Schema-Lauf-Artefakt genau
  `tools/schema/plan.yaml`); über die Range `fb6adf6..4035ee7` listet er
  **fünf** Pfade, darunter `docs/plan/planning/…` und
  `harness/image-hash.txt`. Die **Aussage** stimmt (mit
  `-- internal/adapters/driven/postgresack/` bzw. gegen `--stat` nachgeprüft,
  §Eigene Messungen), die **zitierte Form** trägt sie nicht. Dieselbe
  Nachlässigkeit in der Folgklammer: „`git diff` gegen `Dockerfile` und
  `tools/harness/db-coverage.sh` ist leer" nennt ebenfalls keinen Range —
  geprüft ist die Aussage (diff-frei über die ganze Range), die Form ist
  zweideutig.
- `verifizierbar`: ja — die beiden Befehle ausführen
  (`git diff --name-only` → leer bzw. nur `tools/schema/plan.yaml`;
  `git diff --name-only fb6adf6..4035ee7` → fünf Pfade)

### F-3 — Der führende Kommentar in `seam_test.go` ist syntaktisch ein **zweiter** Package-Doc-Kommentar

- `kategorie`: INFO
- `quelle`: Maintainability (Doc-Kommentar-Konvention; kein Gate)
- `pfad`: `internal/adapters/driven/postgresack/seam_test.go:1-8` gegen
  `internal/adapters/driven/postgresack/ack.go:1-9`
- `befund`: `seam_test.go` erklärt `package postgresack` (paket-intern, nötig
  für die unexportierten Naht-Bezeichner) und trägt seine Test-Erklärung ohne
  Leerzeile unmittelbar vor der Paket-Klausel — damit ist sie eine
  Package-Doc-Kandidatin neben der in `ack.go`. Die Geschwister-Naht aus
  `slice-081` vermeidet das strukturell, weil ihre Testdatei im **externen**
  Testpaket `sqlexec_test` liegt (`postgresstorage/sqlexec/translate_test.go:1-6`);
  die Abweichung ist der Sache nach unvermeidlich, die **Form** ist es nicht.
  Gemessen: sichtbar wird nur der `ack.go`-Text; der `seam_test.go`-Text
  erscheint in keiner Paket-Doku (`go list -f '{{.Doc}}'`), `go vet` ist still.
  Kein Defekt, kein Gate, keine erwartete Aktion.
- `verifizierbar`: ja — `go list -f '{{.Doc}}' ./internal/adapters/driven/postgresack`
  (im gepinnten Toolchain-Container, `GOPROXY=off`) zeigt den `ack.go`-Text;
  `go vet` Exit 0

### F-4 — Die neue nil-Grenze von `newOnSender` fängt nur das **untypisierte** nil

- `kategorie`: INFO
- `quelle`: Maintainability (Go-Interface-nil-Semantik) gegen `slice-084` §1
  („`New(conn, opts…)` samt nil-Grenze … bleiben **unverändert**" — die
  zugesagte `New`-Grenze ist unverändert, das gilt für die **neue** Grenze
  daneben nicht)
- `pfad`: `internal/adapters/driven/postgresack/ack.go:79-82`
- `befund`: `newOnSender` prüft `sender == nil`; ein **typisiertes** nil
  (`(*fakeSender)(nil)` in der Schnittstelle) passiert die Prüfung, weil die
  Schnittstelle dann einen Typ trägt. Gemessen: `newOnSender((*nilProbeSender)(nil))`
  liefert `err = <nil>` und einen nicht-nil Adapter (§Eigene Messungen) — der
  Fehler träte erst als Panik im `Acknowledge`. Erreichbar ist das nur
  paket-intern (der Einstieg ist unexportiert, `New` prüft seinen eigenen Weg
  davor); die Reichweite bleibt deshalb klein. Keine erwartete Aktion.
- `verifizierbar`: ja — Wegwerf-Test gegen `newOnSender` mit typisiertem nil
  (Probe in einer Wegwerf-Kopie, §Eigene Messungen)

---

## Negativbefunde

- geprüft, ohne Befund: **kein Verhaltens-Change** — Zeile-für-Zeile-Vergleich
  von `ack.go` gegen den Vorstand: `New`s Signatur, nil-Grenze und
  Fehlertext (`errNoConnection` trägt wörtlich „%w: keine
  Replication-Verbindung"), `ackLSN` als die alte `IsZero`-Grenze mit
  demselben Text, `standbyStatus` als genau die drei LSN-Felder,
  `replicationClass` als die alte Wrappung, `replicationFailure` unverändert
  loggend — die Ausführung ist deckungsgleich, nur die Aufrufstelle wechselt
  von `pglogrepl.SendStandbyStatusUpdate(ctx, a.conn, …)` auf
  `a.sender.SendStandbyStatusUpdate(ctx, …)`. Belegt zusätzlich durch die
  eigenen Mutationen (M1/M2/M3 rot, §Eigene Messungen) und die realen Tiere.
- geprüft, ohne Befund: **keine abgeschwächte Zusicherung, keine Maskierung** —
  `ack_test.go` ist **byte-identisch** (`git diff` über die Range leer,
  `git log -1 -- ack_test.go` = `0bc3277`, vor dem Slice); `*_test.go`-Diff
  zeigt genau **eine neue** Datei (`seam_test.go`), **keine** gelöschte; kein
  neues `t.Skip`/`SkipNow`, kein `//nolint`/`#noqa`, kein `|| true` im Diff
  (Muster-Grep leer). Die zwei `t.Skip`-Treffer des Pakets sind die
  **vorbestehenden** `CDC_REPLICATION_TEST_DSN`-Wächter in `ack_test.go`.
- geprüft, ohne Befund: **die realen Tests sind der Wächter, nicht die Fakes** —
  eigene Probe: eine Hülle, die **nicht** 1:1 delegiert (`connSender` gibt
  `nil` zurück), ist netzlos **grün** (Exit 0 — die Naht-Tests fahren den Fake,
  nie die Hülle) und wird vom **realen** Tier gefangen (`FAIL:
  TestAcknowledgeOnClosedConnection`, Exit 1), §Eigene Messungen M4. Genau
  das ist die Arbeitsteilung, die `ADR-0080` §Fitness Function zuschreibt.
- geprüft, ohne Befund: **die Hülle trägt den Verdacht aus `slice-085` §6
  nicht** — `connSender` hält `*pgconn.PgConn` als Feld (nicht in der
  Schnittstelle), `standbySender` spricht genau die Operation in der
  **Wertform** des Treibers aus (eine Methode, Signatur wörtlich wie
  `ADR-0080` Festlegung 5), und Fake wie Hülle erfüllen dieselbe Schnittstelle
  (`var _ standbySender = connSender{}` im **Paket**, nicht im Test;
  `var _ standbySender = (*fakeSender)(nil)` im Test). `pglogrepl`-Paket­funktion
  und `*pgconn.PgConn` bleiben innen. Der in `ADR-0080` §„Bestätigt im
  Buchstaben, widerlegt in der Folge" benannte Nicht-Blocker tritt ein.
- geprüft, ohne Befund: **der Null-Befund (a)/(b)/(c) ist gemessen, nicht
  behauptet** — eigene Zahlen decken die Plan-Tabelle Zeile für Zeile
  (1903 unverändert; DB 650 → 659; `postgresstorage` 472 und
  `replication/receive` 155 byte-stabil; `postgresack` 23 → 32, 32/32 gedeckt;
  74,51 %). `k_ab = 0` trägt: das Unit-Profil dieses Baums enthält **keine
  einzige** `postgresack`- oder `replication/receive`-Zeile (Grep = 0).
- geprüft, ohne Befund: **kein anderer Träger berührt** — `git diff --name-status`
  über die Range listet fünf Pfade; die Gegenstands-Listen sind diff-frei
  (`Dockerfile`, `tools/harness/db-coverage.sh` → `DB_COVERAGE_PKGS` und
  `DB_COVERAGE_THRESHOLD`, `harness/mk/coverage.mk` → `THRESHOLD` bleibt 70,
  beide Endstufen 80 %, keine Schwellen-ADR fällig), ebenso
  `internal/bootstrap/wiring.go` (die Composition Root ruft unverändert
  `postgresack.New(stream.Conn(), …)`, `wiring.go:568`), `.a-check.yml`,
  `spec/**`, `replication/receive/**` (→ `slice-085`). Kein Paketwechsel,
  kein Unterpaket.
- geprüft, ohne Befund: **`ADR-0044`-Form** — der Zug ändert Build-Kontext
  (`.dockerignore` lässt `cmd/`, `internal/`, `go.mod`, `go.sum` durch), der
  Digest wechselte `f4475e1e…` → `4bd43435…`, `harness/image-hash.txt` ist der
  deklarierte Träger (`ADR-0044` §1), der Commit `26d4ef0` ist ein eigener
  `chore(image)`-Commit mit `ADR-*` im Betreff. Eigene Gegenprobe: ein
  `docker buildx build` dieses Baums (Metadaten nach `/tmp`, der getrackte
  Träger unberührt) liefert **denselben** Digest `4bd43435…` (§Eigene
  Messungen). Das ist der Befund „unveränderter Digest im selben Builder" —
  kein Inhalts- und kein Umgebungs-Beweis (`ADR-0044` §1/§3).
- geprüft, ohne Befund: **Hygiene** — kein Lauf-Artefakt in einem der vier
  Commits (jede Datei-Liste: Naht, Plan, `harness/image-hash.txt`);
  `tools/schema/plan.yaml` liegt in **keinem** Commit und wurde nach den
  Schema-fahrenden Zielen dieses Laufs zurückgenommen (`git status` leer);
  keine Trailer (`Co-Authored-By`/`Signed-off-by`/„Generated with"),
  Betreffe ohne `STRUCT`-IDs (`commit-traceability: OK — 5 Commit(s),
  Betreffs ohne Struktur-ID`), jede Message mit `ADR-*`; keine host-lokalen
  Pfade (d-check inkl. `hostpaths`-Modul: 704 Dateien, 0 Befunde); keine
  unaufgelösten `BEO-`/`CO-`/`MR-`/`slice-<NNN>`-Kennungen im Diff;
  Historie **linear auf `main`** (`git log --oneline -6` ohne Merge,
  `git rev-parse --abbrev-ref HEAD` = `main`).
- geprüft, ohne Befund: **Vorlagen-Reste** — §2 des Plans trägt keine doppelte
  DoD-Zeile und keinen Platzhalter; die vier `BEDIENHINWEIS`-Kommentare (§1,
  §2, §3, §7) und die `<…>`-Felder in §6/§7 stehen **vor** dem `git mv` nach
  `done/` (dort fängt der `forbid-pattern` der fünften `structure`-Regel sie);
  dieser Stand ist erwartet (dieselbe Feststellung wie `review-slice-081`,
  Negativbefunde).
- geprüft, ohne Befund: **Slice-/Wellen-Chronik in Produktionscode** — der Diff
  **entfernt** einen Fall der Klasse: `ack.go` sagte „wie an jeder anderen
  ctx-losen Konstruktionsstelle dieses **Slices**" und sagt jetzt „dieses
  **Adapters**". In `seam.go`/`seam_test.go` steht weder `slice-0NN` noch
  `seit slice-`; alle Herkunfts-Anker sind `ADR-*`/`LH-*`/`SPEC-*`.
- geprüft, ohne Befund: **`seam_test.go` prüft wirklich** (nicht nur
  „grün, ohne etwas zu prüfen") — drei eigene Mutationen, alle gefangen
  (§Eigene Messungen M1–M3); die Aufrufzähler-Zusage (`calls`/`sent`) ist
  scharf.
- geprüft, ohne Befund: **§8-Sichtung** — der Plan nennt `adapter-fehler-ausgang`
  (2×) und `dod-begruendung-unzutreffende-tatsachenbehauptung` (2×) als
  „benachbart, kein Treffer" bzw. „kein Treffer"; die Registerstände sind am
  Verzeichnis nachgezählt (beide 2×, keine erreicht 3×) — die Einschätzung
  trägt.

---

## Eigene Messungen (Exit-Codes ungepiped, Docker-only)

| Lauf | Exit | Ausgabe / Ergebnis |
|---|---|---|
| `make gates` (Log-Datei, Exit danach separat gelesen) | **0** | `baseline-verify: v6.5.0 OK — 54 Dateien` · d-check 704 Dateien / 0 Befunde (inkl. `hostpaths`) · d-check `commits` 0 Befunde · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"` · `coverage-gate: OK — Coverage 71.90% erfüllt Schwelle 70%` · a-check 0 Befunde |
| `make test` (Race-Detector, netzlos) | **0** | **31 Pakete `ok`**, kein `FAIL`, kein `SKIP`; die übrigen `[no test files]` (`cmd/pg-change-feed`, `queries`, `streamv1`, `port/inbound`, `domain/errors`, die Wegwerf-Clients) |
| `make test-store` | **0** | `…/postgresstorage 4.249s` · `DB-Adapter-Coverage (Teilzahl, nur store-Bestand): 73.52% (gedeckt 347 von 472 Statements)` |
| `make test-replication` (frisches Store-Profil aus demselben Lauf) | **0** | `DB-Adapter-Coverage: 74.51% (gedeckt 491 von 659 Statements; Profile gemergt: store,replication)` · `db-coverage: OK — … erfuellt Schwelle 70%` |
| Dedup-Zählung `merged.coverprofile` (eigene `awk`, Block-Position) | — | `postgresack/ack.go 31/31` · `postgresack/seam.go 1/1` → **32/32** · `postgresstorage/** 347/472` · `receive/** 112/155` · **GESAMT 491/659 = 74,51 %** |
| Dedup-Zählung `/out/coverage.out` der `coverage`-Stufe | — | **1369 gedeckt von 1903** = 71,94 % (Stufe druckt 71.90 %); `grep -c postgresack` = **0**, `grep -c replication/receive` = **0** — `k_ab = 0` für die Unit-Fläche |
| `go test -coverpkg=./internal/adapters/driven/postgresack/... ./internal/adapters/driven/postgresack/...` **ohne DSN** (netzlos, Wegwerf-Kopie) | 0 | **30/32 = 93,75 %** dieses Pakets sind **ohne jede PostgreSQL-Verbindung** gedeckt — die Verdünnungs-Zahl (Schwerpunkt 4) |
| `go test ./internal/adapters/driven/postgresack/...` ohne DSN (Baseline der Mutations-Proben) | 0 | 8 Tests PASS, 2 SKIP (`CDC_REPLICATION_TEST_DSN`-Wächter, vorbestehend) |
| **M1** — nil-Grenze in `newOnSender` entfernt | **1** | `--- FAIL: TestNewOnSenderRejectsNilSender` |
| **M2** — `Acknowledge` setzt **zweimal** ab (Aufrufzähler) | **1** | `--- FAIL: TestAcknowledgeSendsStandbyStatus` |
| **M3** — `New`s nil-Grenze entfernt (typisiertes Durchreichen) | **1** | `--- FAIL: TestNewRequiresConnection` |
| **M4** — `connSender` delegiert **nicht** mehr (`return nil`) | netzlos **0** / real **1** | netzlos grün (die Fakes fahren die Hülle nie) · `FAIL: TestAcknowledgeOnClosedConnection` im `measure`-Lauf gegen reale PostgreSQL |
| Probe typisiertes nil: `newOnSender((*nilProbeSender)(nil))` | 0 | `err=<nil>, adapter-nicht-nil=true` (F-4) |
| `go vet ./internal/adapters/driven/postgresack/...` | 0 | keine Ausgabe |
| `go list -f '{{.Doc}}' ./internal/adapters/driven/postgresack` (`GOPROXY=off`) | 0 | nur der `ack.go`-Text (F-3) |
| `docker buildx build --load --metadata-file /tmp/… -t rev084:img .` | **0** | `containerimage.digest = sha256:4bd4343533…` — identisch mit `harness/image-hash.txt`; Träger blieb unberührt (`git status` leer) |
| `git diff --name-only` (bar) / `git diff --name-only fb6adf6..4035ee7` | — | leer auf sauberem Baum / **fünf** Pfade (F-2) |

Alle Mutationen liefen in **Wegwerf-Kopien** (`/tmp`), der Arbeitsbaum wurde
nie mutiert; das Schema-Lauf-Artefakt `tools/schema/plan.yaml` ist
zurückgenommen, `git status` ist am Ende leer.

---

## Antwort auf die Schwerpunkte

1. **Kein Verhaltens-Change — die tragende Zusage hält.** Die reale
   Verdrahtung geht **unverändert** durch beide Tiere (eigene Läufe:
   `make test-store` Exit 0, `make test-replication` Exit 0, `make test`
   Exit 0, dazu die netzlosen Paketläufe). Keine Zusicherung abgeschwächt:
   `ack_test.go` byte-identisch, keine Testdatei entfernt, kein neues
   `t.Skip`, kein `|| true`, kein `//nolint`. Und die Fakes sind
   **Zusatz, nicht Ersatz**: die M4-Probe zeigt schwarz auf weiß, dass die
   netzline Suite eine **nicht delegierende Hülle** nicht sieht und der reale
   Tier sie fängt — die Wächter-Rolle liegt dort, wo der Plan sie verortet.
2. **Der Null-Befund (a)/(b)/(c) trägt — alle drei Teile selbst gemessen.**
   (a) `k_ab = 0`: Unit-Nenner 1903, und das Unit-Profil dieses Baums enthält
   keine einzige `postgresack`-Zeile; `postgresstorage` 472 und
   `replication/receive` 155 byte-stabil; der DB-Nenner wächst um genau +9
   (23 → 32). (b) der Paket-Diff ist auf **ein** Paket isoliert; `Dockerfile`,
   `db-coverage.sh`, `coverage.mk`, `bootstrap`, `.a-check.yml` sind
   diff-frei — mit einer Einschränkung an der **Form** des Belegs (F-2, LOW),
   nicht an der Aussage. (c) kein Testfall entfernt, die realen Läufe grün. Die
   Plan-Tabelle stimmt in jeder Zeile mit meiner Messung überein — mit **einer**
   Nuancierung: die Spalte „nach dem Zug" der Unit-Zeile nennt 1371 gedeckte
   Statements (72,00 %), mein Lauf misst **1369** (71,94 %/71.90 %). Das ist
   die im Plan selbst benannte Lauf-Varianz (±2) und berührt die Aussage
   nicht — der **Nenner** (1903) ist die tragende Zahl, und er ist stabil.
3. **Die Hülle trägt.** Der `slice-085`-§6-Verdacht („die Paketfunktionen
   verlangen den konkreten `*pgconn.PgConn`") ist **bestätigt und
   folgenlos**: der konkrete Typ bleibt als **Feld** der Hülle innen, die
   Fläche nach außen ist eine Methode in der **Wertform** der Bibliothek, und
   der Fake erfüllt dieselbe Schnittstelle wie die Hülle (beide Zusicherungen
   kompilieren, eine davon im Produktionspfad `seam.go` — wie `ADR-0080`
   §Fitness Function es verlangt). Die Schnittstelle deckt sich wörtlich mit
   `ADR-0080` Festlegung 5; `New`, nil-Grenze und Composition-Root-Verdrahtung
   sind unverändert.
4. **Der Kandidat des Implementers ist keine Harmlosigkeit — aber auch (noch)
   keine undokumentierte.** Eigene Messung: `postgresack` steht im
   DB-Gegenstand bei **32/32**, davon **30/32 (93,75 %) ohne jede
   PostgreSQL-Verbindung** gedeckt; die reale Verbindung hebt die Zahl dieses
   Pakets nur noch um **2** Statements. Das ist das exakte Spiegelbild zu
   `slice-081`: dort hat ein Subjekt-Transfer den DB-Nenner **gedrängt** (ein
   Träger wechselte den Gegenstand), hier **bläht netzlos geprüfter Code den
   DB-Nenner auf**, ohne dass sich der Gegenstand ändert. **Register-Urteil:**
   eine eigene Beobachtung — die Klasse **existiert nicht** (Verzeichnis-Stand
   geprüft: `endstufe-unter-eigenem-messgegenstand-unerreichbar` ist
   *verkörpert* und betrifft ein **unerreichbares Ziel**, nicht einen
   **falsch qualifizierten Ist**; `zahl-in-traeger-driftet-gegen-die-messung`
   betrifft **veraltete Zahlen**, nicht **mitgezählte Herkunft**;
   `adapter-unittest-verdeckt-bootstrap-luecke` betrifft eine **Verdrahtungs**-Lücke,
   nicht die Zählbasis). Vorgeschlagener Pfad:
   `docs/plan/planning/observations/BEO-PGC/<slug>` mit einem `evidence/slice-084.md`;
   Kategorie des Findings: **INFO** (der Diff selbst ist ADR-konform — die
   Grenze ist in `ADR-0080` §„Die benannte Grenze" **entschieden** und trägt
   einen Trigger; der Implementer durfte die Zahl nicht „heilen"). Die
   Handlung liegt bei der Closure: entweder Eintrag anlegen (Zähler 1×) oder
   in §7 ausdrücklich benennen, dass diese Messung das Datum des
   `ADR-0080`-Triggers (a) ist. **Nicht** stillschweigend fallen lassen: 93,75 %
   sind die erste Zahl, an der sich „materiell" überhaupt entscheiden lässt.
5. **`ADR-0044` ist form-korrekt.** Der Zug ändert Build-Kontext-Dateien
   (`internal/**` liegt über `.dockerignore` im Kontext), also läuft `make image`
   vor der Closure, und `harness/image-hash.txt` ist der **richtige Träger**
   (`ADR-0044` §1 — Digest als Lauf-Beleg, nicht als Inhalts-Fingerabdruck).
   Der Digest-Wechsel ist als eigener `chore(image)`-Commit mit `ADR-0044`/`ADR-0080`
   im Betreff abgelegt. Eigene Gegenprobe: ein frischer `docker buildx build`
   dieses Baums ergibt **denselben** Digest — der eingetragene Wert stammt
   damit nachweislich aus einem realen Bau **dieses** Baums (gültiger Befund
   im selben Builder; **kein** Inhalts- oder Umgebungs-Beweis, `ADR-0044` §3).
6. **Hygiene sauber, Historie linear auf `main`.** Vier Commits, kein
   Lauf-Artefakt, keine Trailer, Betreffe mit `ADR-*` und ohne
   Struktur-IDs, `commit-traceability` OK über den Standing-Range, d-check
   inkl. `hostpaths` 0 Befunde, keine unaufgelösten Kennungen, kein
   Branch/Merge. Reste, die **vor** der Closure erwartet sind: die
   `BEDIENHINWEIS`-Blöcke und die §6/§7-Platzhalter — sie gehören dem
   `done/`-Sensor, der sie beim `git mv` fängt.

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** „Stale Zahl in lebendem Träger nach
Gegenstands-Änderung" (F-1, Klasse `zahl-in-traeger-driftet-gegen-die-messung`
— **zweiter Vorgang**, Register damit 2×) · „Beleg-Kommando trägt die zitierte
Aussage nicht" (F-2) · „Test-Datei trägt Package-Doc-Kommentar" (F-3) ·
„Interface-nil-Grenze fängt nur das untypisierte nil" (F-4).

**Beobachtung ohne Finding-Status (Schwerpunkt 4):** „dienst-qualifizierter
Gegenstand zählt netlos gedecktes mit" — Kandidat für das
Beobachtungs-Register, Klasse existiert noch nicht; Messung 30/32 netzlos.

## Verdikt

**Merge-blockierend:** **nein** — 0 HIGH, 0 MEDIUM. Die tragenden Zusagen des
Slices halten gegen jede eigene Messung: kein Verhaltens-Change, keine
abgeschwächte Zusicherung, der Null-Befund in allen drei Teilen belegt, die
Nahtform wörtlich nach `ADR-0080`, der Digest form-korrekt. Die zwei
LOW-Befunde sind **Nachzüge an einem lebenden Träger und an einem zitierten
Beleg**, keine Substanz-Fehler.

**Rückgabe-Pfeil an den Implementer: ja** (klein). F-1 (Zählbasis-Zahlen in
`harness/sensors/db-adapter-coverage.md` auf 659/32/491/74,51 % ziehen oder
als datierten Stand deklarieren) und F-2 (den Beleg-Befehl in §3(b) mit
seinem Range/Pathspec zitieren) — beide sind Ein-Zeilen-Nachzüge und können im
selben Lauf mit der noch offenen §7-Closure gehen. **Kein** Rückweg zur
Plan-Korrektur, **kein** Architect-Zug: F-2 ist ein Zitat, kein Schnitt; F-1
ist eine Zahl, kein Gegenstand (der Gegenstand ist per `ADR-0080` entschieden).

**DoD-Häkchen:** Das Kästchen „Review durchgeführt, Report unter `docs/reviews/`
liegt vor" wird **nicht** von mir nachgezogen — mit dem Rückgabe-Pfeil ist die
Bedingung des Skill-Nachzugs („keine Fixrunde nötig") nicht erfüllt; der
reguläre Nachzug läuft über den Implementer-Lauf (Schritt 21).

**Für die Closure §7 mitzunehmen:** die vier Finding-Klassen oben **und** die
Messung aus Schwerpunkt 4 (30/32 netzlos gedeckt) als Datum des
`ADR-0080`-Triggers (a) — entweder als neuer Register-Eintrag
(`BEO-PGC/…`, `evidence/slice-084.md`) oder als benannter Befund ohne
Register-Zeile; beides ist zulässig, nur das Schweigen nicht (§6-Risiko 3 des
Plans bekommt damit seinen Ausgang).

Dieser Report ist ein **Lauf-Beleg** (dieser Diff, dieser Skill, dieses Modell,
dieses Verdikt); er ersetzt **keine** Verifikation — DoD-/Spec-Konformität
prüft der Verifier separat (Modul 11, anderer Eingabe-Kontext).
