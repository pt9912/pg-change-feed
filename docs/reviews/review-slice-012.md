# Review-Report: slice-012 — 2026-09-10

**Review-Art:** Code — geprüft gegen Slice-Plan + ADRs (Maintainability).

**Gegenstand:** Implementer-Commits von slice-012, Range `c30c624..HEAD`
(Claim-Commit `c30c624` ausgenommen) — sieben Commits: `91ee305`
(Plan-Nachzug), `0362506` (`HeartbeatPort`), `1614c29`
(`PostgresHeartbeatAdapter`), `8c47515` (Schema + `cdc.heartbeat`-View +
`SPEC-001`-Korrektur), `8270fe4` (Bootstrap-Timer + Whitebox-Test),
`20e97d9` (`--healthcheck`-CLI + Compose + Integrationstest-Script),
`7a480f8` (Image-Hash-Beleg).

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-09-09 (vier
repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen) ·
**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-10

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-012-health-endpoint-heartbeat.md`
  (mit vollständigem Plan-Nachzug §3/§8)
- `spec/lastenheft.md` §`LH-FA-ADM-002`, §`LH-QA-OPS-002` ·
  `spec/pflichtenheft.md` §`SPEC-001` (`cdc.process_heartbeat`-Zeile),
  §`SPEC-007` (`HEALTH_STATES`, unverändert)
- [`ADR-0020`](../plan/adr/0020-http-grpc-optional.md) (HTTP/gRPC optional,
  permanent), [`ADR-0024`](../plan/adr/0024-observability-ausserhalb-der-domain.md)
  (Observability außerhalb der Domain, permanent),
  [`ADR-0027`](../plan/adr/0027-capture-application-service.md)
  (Capture Application Service, permanent),
  [`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
  (SQL-Driving-Adapter: Lese-Views direkt, Kategorie C), und das
  Übergabe-Artefakt [`architect-review-slice-011.md`](../plan/adr/architect-review-slice-011.md)
  (Architect-Verdikt: Heartbeat-Muster ohne Folge-ADR trägt)
- `AGENTS.md` §3 Hard Rules (§3.1 Docker-only, §3.2 Suppression-Verbot, §3.7
  Kommentar-Klassen) · `harness/conventions.md` (MR-000/MR-001, genau eine
  Sub-Area `PGC`, Greenfield) · `.a-check.yml` (Layer-Edges,
  `composition_root: internal/bootstrap/**, cmd/**, test/integration/**`)
- Register-Sichtung: `docs/plan/planning/observations/BEO-PGC/` — neun
  Einträge geprüft gegen §8-Behauptung: `health-endpoint-heartbeat` (1×,
  `evidence/slice-011.md`, Zustand `offen`/„weiter offen"),
  `d-migrate-nacharbeit` (2×, `evidence/slice-006.md` +
  `evidence/slice-010.md`), `rollen-verdrahtung` (1×, nicht berührt) —
  Plan-Behauptung real bestätigt
- vorherige Reports: `review-slice-011.md` (F-1/F-2 → Architect-Verdikt,
  Konflikt-Pfad Modul 8, zweiter Praxistest), `review-slice-010.md`
  (Ursprung von `ADR-0046`)

---

## Findings

### F-1 — Compose-Healthcheck-Kommentar beschreibt abwesenden Text statt nur die geltende Zusage

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.7 (Kommentar-Klassen) · Reviewer-Skill-HIGH-Bucket
  „Kommentar trägt keine der Kommentar-Klassen"
- `pfad`: `compose.yaml:69-75` (Commit `20e97d9`)
- `befund`: Der neue Kommentar über dem `feed`-Service-Healthcheck-Block
  schließt mit dem Satz „Vorher (kein Healthcheck-Block) meldete Compose nur
  den Prozess-Start, kein Lebenszeichen des laufenden Captures." Das ist
  exakt das in `AGENTS.md` §3.7 als „Falsch" benannte Muster — ein
  Kommentar, der den abwesenden/vorigen Zustand derselben Stelle erzählt
  („früher stand hier …"), statt ausschließlich die geltende Zusage zu
  nennen. Die vorige Fassung hält `git`; der Satz trägt keine der fünf
  zulässigen Klassen (Zusage/Kopplung/Abgrenzung/Rang-Zeiger/Grenze) — er
  ist reine Historie an einer Stelle, die jeder künftige Lauf mitliest. Die
  ersten fünf Zeilen desselben Kommentarblocks (Grund für `CMD` statt
  `CMD-SHELL`, Bezug auf `bootstrap.Healthcheck`) sind davon nicht
  betroffen und bleiben eine korrekte Kopplungs-Aussage.
- `verifizierbar`: nein — kein Gate fängt das (Grep auf „vorher"/„früher"
  ist kein bestehendes Sensor-Target)
- `klasse`: Kommentar beschreibt abwesenden Text/Historie statt geltende
  Zusage (1. Auftreten in diesem Repo — kein Treffer für „vorher"/„früher"
  in anderen `.go`/`.yaml`/`.sql`-Dateien außerhalb `docs/`)

### F-2 — `bootstrap.Healthcheck` kollabiert Verbindungsfehler, Schema-Drift und echte Staleness in denselben Exit-Code, ohne Diagnose-Ausgabe

- `kategorie`: MEDIUM
- `quelle`: „Maintainability" (Fehlerbehandlung am Rand des Spec-Bereichs,
  MEDIUM-Klasse „unklare Fehlerbehandlung am Rand des Spec-Bereichs")
- `pfad`: `internal/bootstrap/wiring.go:301-315` (`func Healthcheck`) ·
  `cmd/pg-change-feed/main.go:25-33` (Aufrufer)
- `befund`: `Healthcheck()` gibt bei jedem Fehlerfall — Verbindungsaufbau
  scheitert, die Zeile für die Quelle fehlt (nie geschlagen), oder ein
  SQL-/Schema-Fehler (z. B. eine künftig umbenannte Spalte `age_seconds`)
  — denselben Exit-Code `1` zurück, ohne irgendeine Diagnose auf `stderr`.
  `main.go` druckt nur bei einem `ConfigFromEnv`-Fehler eine Meldung; der
  `--healthcheck`-Pfad selbst bleibt bei jedem seiner drei möglichen
  Fehlerursachen stumm. Ein Docker-Healthcheck-Fail zeigt damit nicht, ob
  die Instanz tatsächlich veraltet ist oder ob z. B. die View fehlt
  (Schema-Rollout nicht gelaufen) — beides sieht von außen identisch aus.
- `verifizierbar`: ja — ein Test, der `Healthcheck()` gegen eine
  absichtlich falsche DSN und gegen eine Instanz ohne `cdc.heartbeat`
  laufen lässt und die Ausgabe (nicht nur den Exit-Code) prüft, würde das
  belegen
- `klasse`: unklare Fehlerbehandlung am Rand des Spec-Bereichs (1.
  Auftreten dieser Sub-Klasse — passend zur bereits deklarierten
  MEDIUM-Kategorie des Skills)

### F-3 — Plan-Nachzug (§3) nennt die geänderten Rollout-Belege `plan.yaml`/`down.sql` nicht

- `kategorie`: LOW
- `quelle`: Baseline-Regelwerk `modul-09-implementierung.md` (verkörpert in
  `.claude/commands/implement-slice.md` Schritt 14: „neu gelieferte Dateien
  und Artefakte ebenso")
- `pfad`: `docs/plan/planning/in-progress/slice-012-health-endpoint-heartbeat.md:113-122`
  (Plan-Nachzug-Tabelle) vs. `tools/schema/plan.yaml`, `tools/schema/down.sql`
  (beide in Commit `8c47515` geändert)
- `befund`: Commit `8c47515` ändert neben den acht in der Nachzug-Tabelle
  gelisteten Dateien auch `tools/schema/plan.yaml` und
  `tools/schema/down.sql` (Report-/Rollback-Beleg des
  `schema-rollout`-Laufs, laut Commit-Body selbst benannt: „plan.yaml/
  down.sql sind der Report-/Rollback-Beleg des Rollout-Laufs"). Beide
  fehlen in der §3-Nachzug-Tabelle. Dieselbe Lücke bestand bereits bei
  früheren Slices, die `schema.yaml` änderten (z. B. Commit `90cbfcb` im
  slice-Bestand vor slice-012) — kein neuer Verstoß, aber die Lücke
  besteht in der jetzt dritten praktischen Anwendung der Regel weiterhin.
- `verifizierbar`: ja — Diff-Abgleich Commit-Dateiliste gegen §3-Tabelle
- `klasse`: Plan-Nachzug unvollständig bei generierten Rollout-Artefakten
  (2. benanntes Auftreten dieser Sub-Klasse, informell — nicht im
  Beobachtungs-Register erfasst)

### F-4 — Compose-Healthcheck-Timeout und internes Query-Zeitbudget sind identisch (3s), kein Sicherheitsabstand

- `kategorie`: LOW
- `quelle`: „Maintainability"
- `pfad`: `compose.yaml:79` (`timeout: 3s`) vs.
  `internal/bootstrap/wiring.go:302` (`context.WithTimeout(ctx, 3*time.Second)`)
- `befund`: Der externe Docker-`timeout` und das interne
  Verbindungs-/Query-Zeitbudget von `Healthcheck()` sind beide exakt 3
  Sekunden. Braucht der interne Pfad nahe an sein volles Budget (z. B. bei
  einer langsam antwortenden Instanz), kann Docker den Prozess bereits
  parallel per Timeout als „unhealthy" werten, bevor der interne
  Context-Timeout selbst greift — beobachtbar wird dann ein externer
  Timeout-Abbruch statt der internen, differenzierteren Fehlerbehandlung
  (vgl. F-2).
- `verifizierbar`: nein — Timing-Verhalten, kein deterministischer Gate-Lauf
- `klasse`: Zeitbudget ohne Sicherheitsabstand zwischen innerem und äußerem
  Timeout (1. Auftreten)

---

## Negativbefunde

- geprüft, kein HIGH-Befund: **`--healthcheck`-CLI-Modus als Ermessens-Zug
  (Schwerpunkt 1)** — trägt. Drei unabhängige Prüfungen: (1) *Physikalische
  Notwendigkeit*: `Dockerfile:31-32` bestätigt `gcr.io/distroless/
  static-debian12` (keine Shell, kein Paketmanager, kein `psql`) — ein
  `HEALTHCHECK CMD-SHELL` oder ein eingebettetes `psql` sind damit real
  ausgeschlossen, kein erfundenes Motiv. (2) *`ADR-0020`-Sperrfeld*: Der
  `--healthcheck`-Modus öffnet keinen Netz-Listener, keine neue
  konsumierbare Schnittstelle — er ist ein `CMD`-Exec-Aufruf desselben
  Binarys durch Docker selbst, kein Driving Adapter im Sinne der ADR
  (HTTP/gRPC). Das deckt sich mit dem Architect-Verdikt
  (`architect-review-slice-011.md` §„Berührt das ADR-0020-Sperrfeld?
  Nein."), das exakt diese Grenze bereits für das Heartbeat-Muster
  gezogen hat. (3) *Scope*: Der Compose-Healthcheck-Umbau stand bereits in
  der **ursprünglichen** (nicht nachgetragenen) §3-Zeile „`compose.yaml` |
  update | Healthcheck liest den Heartbeat-Zustand" und in DoD-Punkt 5 —
  `main.go` ist damit kein neu erfundener Umfang, sondern das technische
  Mittel, ein bereits geplantes Ziel unter einer realen Image-Beschränkung
  zu erreichen; der Implementer hat das transparent **vor** dem
  Sensor-Lauf nachgetragen (§3 Plan-Nachzug), nicht verschwiegen. Die
  Lese-Seite (`Healthcheck()` liest `cdc.heartbeat` direkt per SQL) ist
  architektonisch deckungsgleich mit `ADR-0046` Kategorie C (reine
  Projektion, kein Port nötig) — der Compose-Aufruf tritt hier funktional
  als der SQL-Client auf, den `ADR-0046` für Lese-Views bereits vorsieht,
  nur eingebettet statt extern (`psql`), weil das Runtime-Image keinen
  externen SQL-Client trägt. **Kein Konflikt-Pfad nach Modul 8 nötig** —
  kein Rollenwiderspruch, kein heimliches Wachstum.
- geprüft, ohne Befund: **Timer/ACK-Unabhängigkeit (Schwerpunkt 2)** —
  `TestHeartbeatDoesNotBlockCapturePersistAck`
  (`internal/bootstrap/heartbeat_internal_test.go:74`) ist aussagekräftig:
  Der Heartbeat-Mock blockiert real in `Beat()` (Channel-Empfang ohne
  Timeout), während `captureSvc.Capture(...)` in einer zweiten Goroutine
  läuft und innerhalb von 1s durchlaufen muss — ein struktureller Fehler
  (geteilte Goroutine/geteilter Lock) würde den Test zuverlässig in den
  1s-Timeout laufen lassen, kein Race, der zufällig grün wird. Die
  geteilte `events`-Slice wird ausschließlich innerhalb der
  `Capture`-Goroutine beschrieben, der Heartbeat-Mock schreibt nur in den
  gepufferten `calls`-Channel — kein Datenwettlauf zwischen den beiden
  Goroutinen. Die im Kommentar behauptete rot-färbende Mutation (`Beat`
  und Capture-Persist-ACK unter demselben Lock/derselben Goroutine) ist
  strukturell plausibel und deckt sich mit der tatsächlichen Verdrahtung
  in `wiring.go:236-254` (eigene Goroutine, eigener Pool); reale
  Ausführung der Mutation ist Verifier-Sache.
- geprüft, ohne Befund: **`SPEC-001`-Umbenennung `cdc.capture_state` →
  `cdc.process_heartbeat` (Schwerpunkt 3)** — `spec/pflichtenheft.md:137`
  vor dieser Änderung trug nur einen Platzhalter-Tabelleneintrag ohne
  eigene `SPEC-<NNN>`-Detailsektion (anders als z. B. `SPEC-002` für
  `cdc.change`); kein bindendes Akzeptanzkriterium hing an dem alten
  Namen. Die Änderung präzisiert eine bereits vorgesehene, nie detaillierte
  Zeile auf den jetzt real gelieferten Namen — sie führt keine neue
  bindende Anforderung ein und erweitert das Technik-Stratum nicht über
  eine bestehende `LH-*`-ID hinaus (Pflichtenheft ist Rang 2, „technisch
  fortschreibbar"). Kein Spec-Stratum-Verstoß.
- geprüft, ohne Befund: **`cdc.heartbeat`-View gegen `ADR-0046` (Schwerpunkt
  4)** — `tools/schema/nacharbeit-heartbeat.sql:16-21` projiziert nur
  `source_id`, `heartbeat_at` und ein berechnetes `age_seconds`
  (`extract(epoch FROM (now() - heartbeat_at))`); keine
  Schwellenwert-/Autorisierungsentscheidung, keine `HEALTH_STATES`-
  Klassifikation in SQL (`SPEC-007` bleibt beim lesenden System, wie bei
  `cdc.metrics`). Deckungsgleich mit Kategorie C.
- geprüft, ohne Befund: **`a-check`-Konformität für `cmd/` und
  `internal/bootstrap/`** — beide liegen im deklarierten
  `composition_root` (`.a-check.yml:24-27`), unverändert seit vor diesem
  Slice; die neuen Dateien `internal/application/port/outbound/heartbeat.go`
  (importiert nur `context`, `errors`, `domain/model` — Edge `ports→domain`)
  und `internal/adapters/driven/postgresstorage/heartbeat.go` (importiert
  `queries`, `outbound`, `domain/errors`, `domain/model` — Edges
  `adapters→ports`, `adapters→domain`) verletzen keine deklarierte Kante.
- geprüft, ohne Befund: **Docker-only** — der neue `schema-rollout`-Schritt
  läuft über `docker run … $(PG_TEST_IMAGE) psql …` (`Makefile`); kein
  lokales Toolchain-Install; `run-integration-tests.sh` nutzt nur
  `docker inspect`.
- geprüft, ohne Befund: **Suppression-Verbot** — keine `nolint`/`noqa`/
  `SuppressMessage`-Marker in den sieben Commits.
- geprüft, ohne Befund: **Traceability der sieben Commits** — jede Message
  trägt mindestens eine `LH-*`-/`ADR-*`-Kennung im Body (`Bezug:`-Zeile),
  keine Struktur-ID (`SPEC-*`/`ARC-*`) in einem Betreff.
- geprüft, mit einer Ausnahme (F-1): **Kommentar-Klassen (`AGENTS.md`
  §3.7)** — alle übrigen neuen/geänderten Kommentare in `wiring.go`,
  `main.go`, `heartbeat.go` (Port + Adapter), `queries.go`,
  `nacharbeit-heartbeat.sql` sind indikativ über den geltenden Zustand
  (Zusage/Kopplung/Abgrenzung/Rang-Zeiger/Grenze); der aktualisierte
  Vier-Verbindungen-Kommentar in `wiring.go:174-181` nennt korrekt nur den
  jetzt geltenden Zustand, kein Konjunktiv über die vorige Fassung. Der
  aktualisierte Kommentar in `nacharbeit-observability.sql:20-30` erklärt
  eine dauerhafte physikalische Abgrenzung (Rang-Zeiger auf den
  Heartbeat-Mechanismus), keine Historie derselben Datei — anders
  eingeordnet als F-1, weil er nicht den vorigen Zustand *dieser Stelle*
  erzählt, sondern auf den jetzt zuständigen Mechanismus verweist.
- geprüft, ohne Befund: **Plan-Nachzug-Timing und -Struktur (Schwerpunkt 6,
  dritter Praxistest seit slice-009)** — Commit `91ee305` (Plan-Nachzug)
  liegt vor allen Code-Commits desselben Laufs und damit vor dem
  Sensor-Lauf; die Tabelle deckt sieben der neun tatsächlich über die
  ursprüngliche §3-Liste hinaus berührten Dateien korrekt mit Begründung
  (Lücke bei zweien: F-3). Keine stille Streichung eines geplanten Punkts;
  kein Punkt wurde nachträglich aus der ursprünglichen Liste entfernt.
- geprüft, ohne Befund: **§8 Register-Sichtung, Genauigkeit** —
  `health-endpoint-heartbeat` (1×, `evidence/slice-011.md`,
  `weiter offen`), `d-migrate-nacharbeit` (2×, `evidence/slice-006.md` +
  `evidence/slice-010.md`), `rollen-verdrahtung` (1×, korrekt als nicht
  berührt eingeordnet) — alle drei Zählerstände real gegen die
  Register-Dateien bestätigt; die vierte Nacharbeit-Datei dieses Slice
  (`nacharbeit-heartbeat.sql`) ist ein weiterer Treffer der
  `d-migrate-nacharbeit`-Klasse, bleibt aber unter der 3×-Schwelle — ein
  Evidence-Eintrag dafür ist Closure-Sache (§7), kein Implementer-Zug.
- geprüft, ohne Befund: **Sub-Area-Wahl (§8)** — `harness/conventions.md`
  deklariert genau eine Sub-Area (`PGC`, Greenfield); reiner GF-Hinweis
  korrekt ohne Modus-Begründungsblock.
- geprüft, ohne Befund: **WIP-Limit und Lifecycle** — genau ein Slice in
  `in-progress/`; Claim-Commit `c30c624` trägt Kennung im Betreff und ist
  ein reiner `git mv`.
- geprüft, ohne Befund: **`ADR-0027`/Persist-before-ACK** — der
  Heartbeat-Timer berührt `capture.NewCaptureService` nicht; die
  Orchestrierung von Persistenz/ACK bleibt vollständig im Application
  Layer, unverändert.

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Kommentar beschreibt abwesenden
Text/Historie statt geltende Zusage (1. Auftreten) · unklare
Fehlerbehandlung am Rand des Spec-Bereichs (1. Auftreten dieser
Sub-Klasse) · Plan-Nachzug unvollständig bei generierten
Rollout-Artefakten (2. Auftreten, informell) · Zeitbudget ohne
Sicherheitsabstand zwischen innerem und äußerem Timeout (1. Auftreten)

## Verdikt

**Merge-blockierend:** ja, wegen F-1 (HIGH) — ein einzeiliger Fix
(Konjunktiv-/Historien-Satz aus dem `compose.yaml`-Kommentar entfernen,
die Zusage-Sätze davor bleiben stehen). Kein Rollenwiderspruch, kein
Konflikt-Pfad nach Modul 8 nötig — reine Formkorrektur, keine
Architektur- oder Plan-Frage. F-2 bis F-4 sind kein Merge-Stopp, gehören
aber vor Closure adressiert oder mit Begründung zurückgestellt.

**Übergabe:** F-1 bis F-4 gehen an den Implementer (Rückkante Review →
Implementierung; F-1 ist reine Kommentarkorrektur, kein Plan-Defekt). Die
**Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort
in den Zähler. Dieser Report selbst ist ein **Lauf-Beleg** und wird über
Läufe hinweg nicht wieder gelesen. Der Report ersetzt keine Verifikation
— DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).

---

**Gate-Beleg:** `make gates` nach diesem Report-Commit (Lauf 2026-09-10,
Range-Head); Ergebnis im Commit-Text vermerkt.
