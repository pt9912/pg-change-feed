# Verifier-Report: slice-012 — 2026-09-10

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §2 (Definition of Done,
10 Punkte), §3 (Plan-vs-Code, Range `c30c624..64c90c1`, inkl. beider
Nachzüge `91ee305` und `4eb2298`), §6 (Risiko-Ausgang) und
Entscheidungs-Konformität ([`ADR-0020`](../plan/adr/0020-http-grpc-optional.md) ·
[`ADR-0024`](../plan/adr/0024-observability-ausserhalb-der-domain.md) ·
[`ADR-0027`](../plan/adr/0027-capture-application-service.md) ·
[`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md) ·
[`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md) ·
Architect-Verdikt [`architect-review-slice-011.md`](../plan/adr/architect-review-slice-011.md)).
Nicht geprüft: Diff gegen Plan/Hard Rules im Detail über die
F-1…F-4-Dispositionen hinaus (Reviewer, `review-slice-012.md`, Verdikt
dort), realer Bedarf (Validator).

**Gegenstand:** Implementer-/Fix-Range `c30c624..64c90c1` (Claim-Commit
`c30c624` ausgenommen) — 13 Commits: `91ee305` (Plan-Nachzug 1),
`0362506` (`HeartbeatPort`), `1614c29` (`PostgresHeartbeatAdapter`),
`8c47515` (Schema + `cdc.heartbeat`-View + `SPEC-001`-Korrektur),
`8270fe4` (Bootstrap-Timer + Whitebox-Test), `20e97d9`
(`--healthcheck`-CLI + Compose + Integrationstest-Script), `7a480f8`
(Image-Hash-Beleg 1), `7cdf220` (Review-Report), `b809cae`
(Review-Report-Linkfix), `5904fff` (Fix F-2), `22e79c6` (Fix F-1/F-4),
`4eb2298` (Fix F-3), `64c90c1` (Image-Hash-Beleg 2).

**Grundsatz:** Es wurden **keine Behauptungen übernommen** — jeder Sensor
unten wurde in diesem Lauf selbst gefahren, inklusive einer eigenen
Mutationsprobe gegen `TestHeartbeatDoesNotBlockCapturePersistAck` und
einer eigenen `docker inspect`-Beobachtung des realen Compose-Health-Status
(nicht nur des Runner-Skript-Exit-Codes). Plan-Datei und Code blieben
unberührt (die Mutationsprobe lief an einer temporären Kopie, danach
zurückgesetzt und per `git status`/`git diff --stat` als folgenlos
bestätigt); die einzige Schreibaktion dieses Laufs ist dieser Report.

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-10

**Eingangs-Kontext:**

- Slice-Plan §1–§8 am geprüften Stand (`in-progress/slice-012-…md`) ·
  `review-slice-012.md` (F-1 HIGH, F-2 MEDIUM, F-3/F-4 LOW, committet
  `7cdf220`/`b809cae`)
- `spec/lastenheft.md` (`LH-FA-ADM-002`, `LH-QA-OPS-002`, `LH-FA-ADM-003`,
  `LH-QA-REL-001.a`) · `spec/pflichtenheft.md` (`SPEC-001`-Zeile
  `cdc.process_heartbeat`, `SPEC-007` unverändert)
- `ADR-0020`, `ADR-0024`, `ADR-0027`, `ADR-0046` im Volltext,
  `architect-review-slice-011.md`
- `docs/plan/planning/observations/BEO-PGC/health-endpoint-heartbeat/`
  (Zustand `offen`/„weiter offen", 1×) und `BEO-PGC/d-migrate-nacharbeit/`
  (2×) im Volltext
- `internal/application/port/outbound/heartbeat.go`,
  `internal/adapters/driven/postgresstorage/heartbeat.go` +
  `heartbeat_test.go`, `internal/bootstrap/wiring.go` +
  `heartbeat_internal_test.go` + `healthcheck_test.go`,
  `cmd/pg-change-feed/main.go`, `compose.yaml`,
  `tools/schema/nacharbeit-heartbeat.sql`, `tools/schema/schema.yaml`
  (Ausschnitt), `tools/harness/run-integration-tests.sh` im Volltext

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` | `baseline-verify: v6.5.0 OK — 54 Dateien` · `d-check: 151 Datei(en) geprüft, 0 Befund(e)` (voll und `--range HEAD~5..HEAD`) · `commit-traceability: OK — 5 Commit(s), Betreffs ohne Struktur-ID` · `a-check gesamt: 0 Befund(e)` | **0** |
| Traceability der übrigen 9 Range-Commits (über das 5er-Gate-Fenster hinaus, per Lese) | alle 13 Commits `c30c624..64c90c1` tragen im Betreff oder Body mindestens eine `LH-*`-/`ADR-*`-Kennung, kein Betreff trägt `SPEC-*`/`ARC-*` | — |
| `make test` (netzlos, gepinnter Container `golang:1.27-alpine@sha256:cf6fca…`) | 19 Pakete `ok`/`[no test files]`, keine Fehlschläge | **0** |
| `make test-store` (reale PostgreSQL, gepinnte Digests, Rollout inkl. `nacharbeit-heartbeat.sql`) | grün — `postgresstorage 3.218s`; Rollout-Log zeigt alle vier Nacharbeit-Schritte inkl. `cdc.heartbeat`-View; `-v`-Wiederholung zeigt `TestBeatWritesHeartbeatRow`, `TestBeatIsIdempotentAndAdvancesTimestamp`, `TestBeatRejectsEmptySource`, `TestHeartbeatViewProjectsAge` sowie `TestHealthcheckReportsConnectionFailure`/`TestHealthcheckReportsMissingHeartbeatRow` real ausgeführt (nicht geskippt — `CDC_STORE_TEST_DSN` gesetzt) | **0** |
| `make test-integration` (Compose-Umgebung, Feed-Container real gestartet) | 4/4 Tests grün (`TestMVPCaptureFlow`, `TestMVPUpdateOldImageWithFullReplicaIdentity`, `TestMVPActivationState`, `TestMVPDisableRetainedState`) | **0** |
| Eigene `docker inspect`-Beobachtung während des laufenden `test-integration` (Polling-Schleife, nicht nur der Runner-Skript-Exit-Code) | Container `cdc-test-feed`: `starting` (t=6…15s) → **`healthy`** (t=16s); `.State.Health.Log` trägt real zwei Einträge — erster Check `ExitCode 1`, `Output: "pg-change-feed: healthcheck: kein Lebenszeichen für Quelle \"src-mvp\" — Instanz hat noch nie geschlagen"` (vor dem ersten Timer-Tick), zweiter Check `ExitCode 0` (nach dem ersten `Beat`) | — |
| Eigene Mutationsprobe (temporäre Kopie von `heartbeat_internal_test.go`, synchroner `hb.Beat`-Aufruf im Testgoroutine-Kontext vor dem Start der Capture-Goroutine, wie vom Implementer als rot-färbende Mutation beschrieben) | Baseline (unverändert): `--- PASS (0.00s)`. Nach Mutation: `panic: test timed out after 5s` — `goroutine … blockiert in blockingHeartbeat.Beat … heartbeat_internal_test.go:34/83`, Test schlägt real fehl | **1** |
| Dieselbe Probe nach Rückbau der Mutation (Datei aus Sicherung wiederhergestellt) | `--- PASS (0.00s)`, `git status`/`git diff --stat` zeigen keine Restspur | **0** |
| `git status`/`git diff --stat` nach allen Läufen | leer — keine Plan-/Code-Änderung durch diesen Verifikationslauf | — |

**Mutationsprobe im Detail:** Kopie der Original-Testdatei gesichert;
direkt nach der Initialisierung von `hb` und vor dem Start der
`runHeartbeat`-Goroutine ein synchroner `hb.Beat(ctx, "src-1")`-Aufruf im
Testgoroutine-Kontext eingefügt — Stellvertreter für „Heartbeat und
Capture-Persist-ACK laufen in derselben Goroutine/demselben Lock", die
vom Implementer im Testkommentar benannte rot-färbende Mutation. `Beat`
blockiert real auf `<-b.release` (nichts gibt frei), der Testlauf mit
`-timeout 5s` endet in einem echten `panic: test timed out`, sichtbar in
genau der Goroutine der Mutation. Das bestätigt: Der Test ist kein
zufällig grüner Test — er hängt tatsächlich an der Goroutine-Trennung,
die `runHeartbeat` (eigene Goroutine, eigener Pool, `wiring.go:249-259`)
bereitstellt. Datei danach vollständig zurückgesetzt (Kopie überschrieben),
erneuter Lauf grün, kein Diff im Arbeitsbaum.

Nicht selbst gefahren: `make schema-validate` isoliert (läuft als Teil von
`make test-store`/`make test-integration` ohnehin mit) ·
`make test-replication` (kein DoD-Item dieses Slice, der Heartbeat-Timer
berührt den Replication-Stream-Pfad nicht) · `make image`-Neubau zur
Digest-Verifikation (kein Gate, `harness/image-hash.txt` wurde im
geprüften Range zweimal fortgeschrieben — `7a480f8`, `64c90c1` — Inhalt
gelesen und plausibilisiert, nicht neu gebaut).

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | Heartbeat-Schreiber im Capture-Prozess (Timer, Bootstrap-Verdrahtung) — Teil-Beleg `LH-FA-ADM-002` | **bestätigt** | `runHeartbeat` (`wiring.go:270-281`) läuft über einen `time.Ticker` in eigener Goroutine (`wiring.go:249-259`), eigener Verbindungspool (`postgresstorage.NewHeartbeat`, getrennt von Store-/Aktivierungs-Pool); `HeartbeatPort`/`PostgresHeartbeatAdapter` real vorhanden und über `make test-store` grün belegt (`TestBeatWritesHeartbeatRow`, `TestBeatIsIdempotentAndAdvancesTimestamp`) |
| 2 | `cdc.process_heartbeat`-View (`cdc.heartbeat`) + `make test-store`-Beleg — Teil-Beleg `LH-QA-OPS-002` | **bestätigt** | `tools/schema/nacharbeit-heartbeat.sql` erzeugt `cdc.heartbeat` als reine Projektion (`source_id`, `heartbeat_at`, `age_seconds`), real ausgerollt im `test-store`-Lauf, `TestHeartbeatViewProjectsAge` grün gegen reale PostgreSQL |
| 3 | `make gates` grün | **bestätigt** | Exit 0 in diesem Lauf (Sensor-Tabelle oben) |
| 4 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **bestätigt** | `review-slice-012.md` committet (`7cdf220`, Linkfix `b809cae`); alle vier Findings F-1…F-4 in Folge-Commits disponiert (siehe Dispositions-Prüfung unten) |
| 5 | Doku-Update für `compose.yaml` (Healthcheck liest den Heartbeat statt nur den Prozess-Start) | **bestätigt** | `compose.yaml:69-79` trägt den `--healthcheck`-Aufruf und einen indikativen Kommentar (nach F-1-Fix ohne Historien-Satz); real am Compose-Container bestätigt (`.State.Health.Status` `starting`→`healthy`, Log-Einträge zeigen den `cdc.heartbeat`-Bezug) |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **noch offen** | §7 trägt am geprüften Stand weiterhin die Vorlagen-Platzhalter (`<…>`) — korrekt, da die Verifikation vor der Closure läuft; kein Defekt, regulärer Vor-Closure-Zustand |
| 7 | Reconciliation-Register, falls Inventur-Fund | **entfällt — korrekt geprüft** | `docs/plan/planning/reconciliation.md` existiert nicht (Repo ohne Brownfield-Bootstrap, `harness/conventions.md`: Sub-Area `PGC` durchgehend Greenfield) |
| 8 | Beobachtungs-Register fortgeschrieben | **noch offen (regulär)** | `BEO-PGC/health-endpoint-heartbeat/state.md` trägt am geprüften Stand weiterhin `Zustand: offen` / „weiter offen", 1× (`evidence/slice-011.md`) — der Ausgang *eingetreten* (dieser Slice liefert) ist erst mit der Closure zu setzen; `BEO-PGC/d-migrate-nacharbeit` (2×) korrekt unangetastet, `nacharbeit-heartbeat.sql` ist ihr vierter Treffer und noch nicht als `evidence/slice-012.md` eingetragen — Closure-Sache, kein Implementer-Defekt (Register-Sichtung §8 hatte das bereits korrekt so vorgemerkt) |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **noch offen (regulär), mit Verifier-Befund zur Disposition** | Risiko 1 trägt im Plan bereits „Ausgang: eingetreten (Träger: dieser Slice)" — korrekt, dieser Slice liefert den Heartbeat. Risiko 2 trägt noch „Wird bei Closure bewertet", kein Ausgang aus der geschlossenen Drei-Menge; die eigene Mutationsprobe dieses Laufs liefert dem Planner die Evidenzbasis für **entfallen** (das Risiko ist nicht eingetreten — Timer- und Persist-ACK-Pfad sind strukturell und per rot-färdender Mutation belegt getrennt) — Entscheidung selbst bleibt Planner-Sache |
| 10 | Drei Paarungen (Anker · Folge-Slice · Register) | **korrekt delegiert** | `welle-4.md` liegt flach (offen), slice-013/014 noch nicht in `done/` — Paarungen prüft die nächste welle-4-Closure, nicht dieser Slice allein (Repo **mit** Wellen-Betrieb) |

**Zwischenstand: 6/10 Kriterien materiell erfüllt und in diesem Lauf
selbst geprüft, 1 korrekt entfallen (Item 7), 1 korrekt an die
Welle-4-Closure delegiert (Item 10), 3 korrekt noch offen, weil
Verifikation vor Closure läuft (Items 6, 8, 9).** Kein DoD-*Defekt* im
Sinn eines unbelegten „bestätigt"-Punkts.

## Timer/ACK-Unabhängigkeit — Mutationsproben-Nachvollzug

Siehe Sensor-Tabelle oben für den vollen Ablauf. Ergebnis: **bestätigt,
kein Zufallsgrün.** `TestHeartbeatDoesNotBlockCapturePersistAck` besteht
mit der realen Implementierung (`wiring.go` startet `runHeartbeat` in
einer eigenen Goroutine über einen eigenen `HeartbeatPort`/Pool) und
schlägt real fehl, sobald die vom Implementer benannte
Struktur-Verletzung (`Beat` und Capture-Persist-ACK teilen sich
Goroutine/Lock) nachgebildet wird — nicht als Race, sondern als
deterministischer 5s-Timeout mit sichtbarem Stack an genau der
blockierenden Stelle. Die reale Compose-Umgebung bestätigt dasselbe
Muster am Integrationstest: Der erste `--healthcheck`-Aufruf lief, bevor
der erste Timer-Tick geschrieben hatte (`ExitCode 1`, Diagnose „kein
Lebenszeichen … noch nie geschlagen"), der zweite bereits grün — der
Heartbeat-Zug lief unabhängig vom restlichen Verdrahtungsstart weiter,
kein Blockieren des Feed-Containers.

## F-1…F-4-Dispositionsprüfung (`review-slice-012.md`)

| Finding | Fix-Commit | Verdikt |
|---|---|---|
| F-1 (HIGH — Kommentar erzählt abwesenden Vorzustand, `compose.yaml`) | `22e79c6` | **getragen** — der Satz „Vorher (kein Healthcheck-Block) meldete Compose nur den Prozess-Start…" ist entfernt, die verbleibenden Zeilen sind indikativ über den geltenden Zustand (Zusage/Kopplung). Diff geprüft, kein Rest-Konjunktiv |
| F-2 (MEDIUM — kollabierte Fehlerklassen ohne Diagnose) | `5904fff` | **getragen, real geprüft.** Fünf Fehlerklassen tragen fünf unterscheidbare `stderr`-Zeilen (Substring-eindeutig, keine Überschneidung): „DSN ungültig" (`pgxpool.New`-Fehler) · „Instanz nicht erreichbar" (`pool.Ping`-Fehler) · „kein Lebenszeichen für Quelle … — Instanz hat noch nie geschlagen" (`pgx.ErrNoRows`) · „`cdc.heartbeat` nicht lesbar (Schema-Rollout gelaufen?)" (sonstiger SQL-Fehler) · „Lebenszeichen veraltet (…, Schwelle …)" (Alters-Schwelle). Zwei neue Tests (`TestHealthcheckReportsConnectionFailure` netzlos, `TestHealthcheckReportsMissingHeartbeatRow` via `make test-store`) prüfen je die **An**- und **Ab**wesenheit der jeweils anderen Diagnose-Zeile — in diesem Lauf beide real (nicht geskippt) grün. Zusätzlich am realen Compose-Container beobachtet: die „kein Lebenszeichen"-Diagnose erscheint tatsächlich im `.State.Health.Log` des ersten fehlschlagenden Checks. Trägt zugleich `LH-QA-OPS-002`s eigene Messmethode („inklusive Fehlerzustandsfall, `LH-FA-ADM-003`") besser als vor dem Fix |
| F-3 (LOW — Plan-Nachzug ohne `plan.yaml`/`down.sql`) | `4eb2298` | **getragen für die benannten zwei Dateien** — beide Zeilen jetzt in §3. Siehe aber **V-1** unten: derselbe Fix-Lauf hinterließ eine weitere, in §3 ungetragene Datei |
| F-4 (LOW — Timeout ohne Sicherheitsabstand) | `22e79c6` | **getragen** — `compose.yaml`-`timeout` von 3s auf 5s angehoben, internes `context.WithTimeout`-Budget unverändert bei 3s; damit besteht jetzt ein 2s-Abstand, in dem die interne Diagnose (F-2) vor einem externen Docker-Timeout greifen kann |

## Plan-vs-Code-Diff (Range `c30c624..64c90c1`, beide Nachzüge)

**Deckung §3 gegen den vollen Datei-Diff:**

```
git diff --name-only c30c624..64c90c1
```

liefert: `cmd/pg-change-feed/main.go`, `compose.yaml`,
`docs/plan/planning/in-progress/slice-012-health-endpoint-heartbeat.md`,
`docs/reviews/review-slice-012.md`, `harness/image-hash.txt`,
`internal/adapters/driven/postgresstorage/heartbeat.go`,
`internal/adapters/driven/postgresstorage/heartbeat_test.go`,
`internal/adapters/driven/postgresstorage/queries/queries.go`,
`internal/application/port/outbound/heartbeat.go`,
**`internal/bootstrap/healthcheck_test.go`**,
`internal/bootstrap/heartbeat_internal_test.go`,
`internal/bootstrap/wiring.go`, `Makefile`, `spec/pflichtenheft.md`,
`tools/harness/run-integration-tests.sh`, `tools/schema/down.sql`,
`tools/schema/nacharbeit-heartbeat.sql`,
`tools/schema/nacharbeit-observability.sql`, `tools/schema/plan.yaml`,
`tools/schema/schema.yaml`.

Die sechs ursprünglichen §3-Zeilen und die zehn Nachzug-Zeilen (nach
beiden Nachzügen `91ee305`/`4eb2298`) decken 19 der 20 tatsächlich
berührten Dateien — `docs/plan/planning/…md` und
`docs/reviews/review-slice-012.md` selbst zählen nicht als Liefer-Punkt
(Plan-/Review-Artefakte), `harness/image-hash.txt` ist Lauf-Beleg
(`ADR-0044`, kein eigener Liefer-Punkt). **Eine Datei fehlt in §3 vollständig:
`internal/bootstrap/healthcheck_test.go`** — siehe **V-1**.

## Befunde

### V-1 — `internal/bootstrap/healthcheck_test.go` fehlt vollständig in §3 (weiteres Auftreten derselben Klasse wie F-3, in dessen eigenem Fix-Lauf entstanden)

- `kategorie`: LOW
- `quelle`: Baseline-Regelwerk `modul-09-implementierung.md` (Plan-Nachzug-Pflicht) · dieselbe Klasse wie `review-slice-012.md` F-3
- `pfad`: `docs/plan/planning/in-progress/slice-012-health-endpoint-heartbeat.md` §3 (weder Original- noch Nachzug-Tabelle) vs. `internal/bootstrap/healthcheck_test.go` (neu, Commit `5904fff`)
- `befund`: Der F-2-Fix-Commit `5904fff` legt die neue Datei
  `internal/bootstrap/healthcheck_test.go` an (81 Zeilen, zwei neue
  Tests). Der spätere Plan-Nachzug-Commit `4eb2298` — der exakt für
  F-3 „Plan-Nachzug unvollständig" geschrieben wurde — trägt nur die
  beiden in F-3 benannten Dateien (`plan.yaml`, `down.sql`) nach, nicht
  aber diese neue Testdatei aus demselben Fix-Lauf. §3 nennt an keiner
  Stelle `healthcheck_test.go`. Kein DoD-Blocker (die Datei ist real,
  getestet, gehört fachlich zum bereits geplanten Liefer-Punkt „Doku-
  Update `compose.yaml`" bzw. der F-2-Disposition selbst), aber
  strukturell dieselbe Lücke, die F-3 im Review bereits als „2.
  Auftreten" benannt hatte — mit diesem Fund wäre es review-übergreifend
  das dritte benannte Auftreten der Klasse „Plan-Nachzug unvollständig",
  allerdings **innerhalb desselben Slice**, nicht über Slices hinweg wie
  bei F-3 selbst — die 3×-Schwelle des Beobachtungs-Registers zählt über
  *Vorgänge* (Slices), nicht über Commits desselben Slices; ob dieser
  Fund als weiterer Beleg für `d-migrate-nacharbeit`-artige
  Register-Klassen zählt, ist eine eigene Frage, kein automatischer
  dritter Treffer.
- `verifizierbar`: ja — Diff-Abgleich Commit-Dateiliste gegen §3-Tabelle
  (wie bei F-3)
- `klasse`: Plan-Nachzug unvollständig, diesmal im Fix-Round-Commit
  selbst statt im ursprünglichen Feature-Commit
- **Für die Closure:** eine weitere §3-Nachzug-Zeile für
  `internal/bootstrap/healthcheck_test.go` ergänzen, vor dem `git mv`
  nach `done/`.

### V-2 — `schema.sql`-Kommentar nennt weiterhin den durch `SPEC-001` abgelösten Namen `cdc.capture_state`

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Kommentar nennt die geltende Zusage, keine
  zwei Namen für denselben Zweck) · Commit `8c47515` selbst (Begründung
  der `SPEC-001`-Umbenennung: „sonst zwei Namen für denselben Zweck")
- `pfad`: `internal/adapters/driven/postgresstorage/schema.sql:11`
  (Bestand, von diesem Slice nicht berührt) vs.
  `spec/pflichtenheft.md:137` (`SPEC-001`, jetzt `cdc.process_heartbeat`)
- `befund`: Der Store-Adapter-DDL-Kommentar zählt weiterhin
  „`cdc.capture_state` (Betriebs-/Capture-Zustand)" als eine der drei
  Tabellen, die die handgeschriebene DDL bewusst nicht trägt. Mit der
  `SPEC-001`-Umbenennung dieses Slices trägt das Pflichtenheft diese
  Tabelle jetzt unter dem Namen `cdc.process_heartbeat` — der
  Store-Adapter-Kommentar ist damit sachlich weiterhin korrekt (die DDL
  trägt die Tabelle tatsächlich nicht), nennt aber einen Namen, den kein
  anderes Artefakt im Repo mehr für dieses Konzept führt. Das ist exakt
  das Risiko, das die eigene Begründung von `8c47515` benennt („sonst
  zwei Namen für denselben Zweck") — nur zwischen Spec und Code statt
  innerhalb der Spec. Kein DoD-Blocker (der Slice hat `schema.sql` nicht
  in §3 geplant, Bestand bleibt bewusst außerhalb seines Umfangs), aber
  ein realer Nachzieh-Bedarf für die Closure oder ein benanntes
  Folge-Slice.
- `verifizierbar`: ja — Grep auf `capture_state` in Produktcode/-kommentaren
  außerhalb `docs/`/ADRs (ADRs und alte Review-Reports sind laut Hard
  Rule 3.5/immutable korrekt unverändert und zählen nicht)
- `klasse`: Kommentar mit abgelöstem Namen nach Spec-Umbenennung (1.
  Auftreten dieser Sub-Klasse in diesem Repo)
- **Für die Closure:** entweder `schema.sql:11` in diesem Slice
  nachziehen (Kommentar-Fix, kein DDL-Wandel, kein neuer Liefer-Punkt im
  Sinne von §2) oder als benannte Beobachtung ins Register aufnehmen.

## Negativbefunde

- geprüft, ohne Befund: **`SPEC-001`-Umbenennung selbst
  (`spec/pflichtenheft.md`)** — keine verwaiste `cdc.capture_state`-Zeile
  oder -Referenz mehr innerhalb der Spec-Straten oder im aktuellen
  Produktcode/-tests; `SPEC-007` unverändert, keine neue bindende
  Anforderung eingeführt (bestätigt Review-Negativbefund)
- geprüft, ohne Befund: **`a-check`-Konformität** — 0 Befunde für die
  neuen/geänderten Go-Dateien inkl. `healthcheck_test.go`
  (`internal/bootstrap`, `composition_root`)
- geprüft, ohne Befund: **Docker-only** — alle Sensor- und
  Mutationsproben-Läufe dieses Verifikationslaufs liefen über
  `docker run`/`make`, kein lokales Toolchain-Install
- geprüft, ohne Befund: **Suppression-Verbot** — keine
  `nolint`/`noqa`/`SuppressMessage`-Marker im gesamten Range
- geprüft, ohne Befund: **Kommentar-Klassen (`AGENTS.md` §3.7)** außerhalb
  V-2 — die neuen/geänderten Kommentare in `wiring.go`, `main.go`,
  `heartbeat.go` (Port + Adapter), `queries.go`,
  `nacharbeit-heartbeat.sql`, `healthcheck_test.go` sind indikativ über
  den geltenden Zustand
- geprüft, ohne Befund: **WIP-Limit und Lifecycle** — genau ein Slice
  (`slice-012`) in `in-progress/` am Prüfzeitpunkt
- geprüft, ohne Befund: **`ADR-0027`/Persist-before-ACK strukturell** —
  `capture.NewCaptureService`-Aufruf unverändert, Heartbeat-Timer läuft
  über eigenen Pool/eigene Goroutine (siehe Mutationsprobe oben)
- geprüft, ohne Befund: **Arbeitsbaum nach allen Sensor-/Probenläufen** —
  `git status`/`git diff --stat` leer; die einzige Schreibaktion dieses
  Laufs ist dieser Report

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Plan-Nachzug unvollständig, diesmal im
Fix-Round-Commit selbst (V-1) · Kommentar mit abgelöstem Namen nach
Spec-Umbenennung (V-2, 1. Auftreten).

**Zusammenfassung DoD:** **6/10 Kriterien materiell erfüllt und in diesem
Lauf selbst geprüft** (`make gates` grün, `make test` grün netzlos,
`make test-store` grün am realen PostgreSQL inkl. beider F-2-Tests real
ausgeführt, `make test-integration` grün mit eigener `docker
inspect`-Bestätigung des realen `healthy`-Übergangs), 1 Item korrekt
entfallen (Reconciliation-Register), 1 Item korrekt an die
Welle-4-Closure delegiert (drei Paarungen), 3 Items regulär noch offen,
weil Verifikation vor Closure läuft (Closure-Notiz, Beobachtungs-Register,
Risiko-Ausgang). **Kein DoD-Defekt.**

**F-1…F-4 sind geschlossen:** alle vier Review-Findings real im Code
disponiert und in diesem Lauf am Artefakt nachgeprüft, nicht nur am
Commit-Text übernommen — F-2 zusätzlich real am laufenden
Compose-Container beobachtet.

**Timer/ACK-Unabhängigkeit ist real bewiesen, kein Zufallsgrün:** eigene
Mutationsprobe reproduziert die vom Implementer beschriebene
rot-färbende Mutation und lässt den Test real fehlschlagen (Timeout mit
sichtbarem Blockier-Stack), die unveränderte Implementierung besteht
sowohl isoliert als auch end-to-end am realen Container.

**Neue Verifier-only-Funde (V-1, V-2):** zwei LOW-Findings, keiner
DoD-blockierend, beide vor dem `git mv` nach `done/` sinnvoll
nachzuziehen — genau die Art Befund, für die die Verifier-Rolle vom
Reviewer getrennt ist (V-1 entstand im Fix-Round-Commit, den kein Review
mehr sah; V-2 ist eine Kreuz-Artefakt-Konsistenzfrage, die Review als
reinen Code-vs-Plan-Diff nicht adressiert).

## Verdikt

**DoD-/Entscheidungs-Konformität:** DoD materiell erfüllt, kein
HIGH-/MEDIUM-Finding offen, keine ADR- oder Hard-Rule-Verletzung
gefunden. Heartbeat-Schreiber, `cdc.heartbeat`-View und der
`--healthcheck`-Compose-Vertrag sind real und sensor-belegt — inklusive
einer eigenen, am realen Container ausgeführten `docker
inspect`-Beobachtung und einer eigenen, real ausgeführten
Mutationsprobe der Timer/ACK-Unabhängigkeit.

**Blockierend für Closure (`git mv` nach `done/`):**

1. **V-1:** `internal/bootstrap/healthcheck_test.go` in §3 nachtragen.
2. **V-2:** `schema.sql:11`-Kommentar auf `cdc.process_heartbeat`
   nachziehen oder als benannte Beobachtung registrieren.
3. **§6-Risiko 2** disponieren (Empfehlung dieses Reports: **entfallen**
   — durch die eigene Mutationsprobe belegt) und `BEO-PGC/health-endpoint-heartbeat`
   den Ausgang **eingetreten** zuweisen (Plan trägt das bereits vor).
4. **§7 Closure-Notiz** mit Steering-Loop-Lerneintrag füllen —
   `nacharbeit-heartbeat.sql` als vierten Beleg für
   `BEO-PGC/d-migrate-nacharbeit` erwägen (Zähler würde auf 3× steigen,
   Register-Lese-Schritt der nächsten Welle-4-Closure).
5. **Paarungen (Item 10 der Vorlage, delegiert):** Anker-, Folge-Slice-
   und Register-Paarung prüft die welle-4-Closure — auch für diesen
   Slice, erst wenn slice-013/014 ebenfalls in `done/` liegen.

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Plan-Datei und
Code wurden von diesem Lauf nicht verändert; die temporäre
Mutationsproben-Kopie und alle Testcontainer/-netze dieses Laufs wurden
vollständig zurückgesetzt bzw. entfernt.
