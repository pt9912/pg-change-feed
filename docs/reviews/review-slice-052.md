# Review-Report: slice-052 — 2026-09-13

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10
§Drei Review-Arten), geprüft gegen Plan/ADR/Hard Rules (Maintainability),
**nicht** gegen die DoD (Verifier-Aufgabe, Modul 11).

**Gegenstand:** `slice-052` — vier Implementierungs-Commits
`2bd3fb4`..`2ca089c` (Diff gegen Elter `ab69fd2`, dem reinen
Lifecycle-`git mv` `next` → `in-progress`):

- `2bd3fb4` — `ChangeNotificationPort` + `ErrNotify`-Sentinel
- `5d362fa` — `NatsChangeNotificationAdapter` (`natsnotify`)
- `2538110` — `CaptureService`-Erweiterung (Notify nach ACK, best-effort)
- `2ca089c` — `make test-notify` (NATS-Testcontainer)

Ein fünfter, unabhängiger Commit (`15ee973`, `slice-056` angelegt) ist
zwischenzeitlich auf `main` gelandet, während dieser Review lief — er
berührt keine der oben genannten Dateien und ist nicht Gegenstand dieses
Reports.

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-09-09 (vier
repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-052-changenotification-port-natsnotify-adapter.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §6 Risiken)
- `ADR-0055` (`docs/plan/adr/0055-nats-change-notification-wecksignal.md`)
- `LH-FA-SST-007`, `SPEC-017` (`spec/pflichtenheft.md`)
- Referenz-Adapter `internal/adapters/driven/postgresack/ack.go`
- `AGENTS.md` §3 Hard Rules (insb. 3.1, 3.3, 3.7, 3.8), §6 Workflow
- `harness/conventions.md` (MR-000 ID-Schema)
- Vorherige Findings am gleichen Modul: Review zu `slice-050`,
  Review zu `slice-051`
- `docs/plan/planning/observations/BEO-PGC/slice-chronik-in-code-kommentar/`
  (`observation.md`, `state.md`) und
  der Architect-Verdikt zur Slice-Chronik in Code-Kommentaren
  (Architect-Verdikt, 3× gezählt, verkörpert als geschärfte
  Selbstprüf-Instruktion)

---

## Findings

### F-1 — Slice-Chronik in Produktionscode-Kommentar (`natsnotify.New`)

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.7 (Hard Rule „Ein Kommentar beschreibt, was da
  ist") · Baseline-Regelwerk `grundlagen-harness-dateien.md` §Was ein
  Kommentar trägt · Präzedenzfall
  `docs/plan/planning/observations/BEO-PGC/slice-chronik-in-code-kommentar/`
  (3× gezählt, Architect-Verdikt vom 2026-09-13)
- `pfad`: `internal/adapters/driven/natsnotify/notify.go:59-61`
- `befund`: Der Godoc-Kommentar über der Produktionsfunktion `New` lautet:
  „New legt den Notify-Adapter auf eine bestehende NATS-Verbindung; die
  Composition Root verdrahtet den Verbindungsaufbau über `CDC_NATS_URL`
  (Folge-Slice `slice-053`)." Subjekt des Satzes ist der
  **Produktionscode-Pfad** (die künftige Composition-Root-Verdrahtung),
  nicht ein Testfall — exakt die Unterscheidung, die der
  Architect-Verdikt als tragenden Trennstrich zwischen zulässiger
  Testfall-Provenienz und unzulässiger Chronik benennt (§„Strukturell
  ununterscheidbar von den drei tatsächlichen Verstößen"). Der
  Referenz-Adapter `postgresack.New` löst dieselbe Aussage ohne
  Slice-Nummer: „die Composition Root verdrahtet beide Adapter über
  `receive.Stream.Conn` (`ADR-0026`)" — ADR-Bezug statt Slice-Bezug. Ein
  reiner `ADR-0055`-Verweis hätte hier ebenso getragen; die Slice-Nummer
  trägt nichts, was das ADR nicht bereits sagt, und veraltet, sobald
  `slice-053` in `done/` wandert oder umnummeriert wird.
- `verifizierbar`: nein — kein Sensor deckt diese Klasse (Architect-Verdikt
  hat einen repo-weiten Textmuster-Sensor geprüft und wegen der
  etablierten, weiterhin legitim wachsenden Testfall-Provenienz-Konvention
  verworfen; nur eine diff-skopierte Selbstprüf-Instruktion im
  Implementer-Workflow ist vorgesehen).
- `klasse`: Slice-Chronik in Code-Kommentar

Dies wäre das **vierte** gezählte Auftreten dieser Klasse (nach den drei,
die den bestehenden Architect-Verdikt auslösten) — nach der eigenen
Eskalationsnotiz in `BEO-PGC/slice-chronik-in-code-kommentar/state.md`
("Tritt die Klasse trotz der geschärften Instruktion ein viertes Mal
auf, ist das ein Signal, dass Enumeration allein nicht trägt") ein Signal
an den Planner, unabhängig vom Ausgang dieser Fixrunde.

### F-2 — Code-Kommentar referenziert nicht-committetes „Implementer-Bericht"-Artefakt

- `kategorie`: MEDIUM
- `quelle`: Maintainability (`AGENTS.md` §3.7 — ein Kommentar schreibt an
  den, der die Stelle ändert; ein künftiger, sessionsferner Leser kann
  das Ziel nicht auflösen)
- `pfad`: `internal/application/usecase/capture/service_test.go:335-338`
- `befund`: Der Kommentar über `TestCaptureSucceedsDespiteFailingNotification`
  lautet u. a. „Rot-Beleg (nicht eingecheckt, siehe Implementer-Bericht):
  entfernt man das Abfangen …, schlägt genau dieser Test fehl." Der
  Begriff „Implementer-Bericht" ist im Repo etabliert (50+ Treffer), aber
  **ausschließlich** in Lauf-gebundenen Markdown-Artefakten
  (`docs/reviews/*.md`, `docs/plan/planning/done/*.md`,
  Beobachtungs-Register-`evidence/*.md`) — Dokumente, die selbst als
  Audit-Beleg *eines bestimmten Laufs* gelesen werden und in deren Kontext
  ein Verweis auf den zugehörigen Sitzungsbericht plausibel ist. Dies ist
  der erste Fund dieses Begriffs in einem **committeten `.go`-Kommentar**
  — einem permanenten Artefakt, das unabhängig von der Session gelesen
  wird, in der es entstand. Ein späterer, sessionferner Leser (anderer
  Agent, anderer Mensch) kann „Implementer-Bericht" nicht auflösen.
- `verifizierbar`: nein — kein Gate prüft Kommentar-Referenzen auf
  Auflösbarkeit außerhalb der Slice-Chronik-Klasse (F-1).
- `klasse`: Code-Kommentar referenziert ephemeres Implementer-Bericht-Artefakt

## Negativbefunde

- geprüft, ohne Befund: `internal/application/port/outbound/changenotification.go`
  — Signatur `Notify(ctx, sourceID string) error` deckt sich wörtlich mit
  `ADR-0055`s Port-Design; `ErrNotify` klassifiziert `transient`
  (`SPEC-008`), analog zu `outbound.ErrReplication`; keine Domain-Typen
  nötig (Layer-Konformität, `.a-check.yml`: Port referenziert nur
  Standardbibliothek).
- geprüft, ohne Befund: `internal/adapters/driven/natsnotify/notify.go`
  — Konstruktions-/Options-Muster (`New`, `Option`, `WithLog`) ist
  strukturell identisch zu `postgresack.go` (bis auf F-1); Subjekt-Schema
  `cdc.changes.<source_id>` deckt sich mit `SPEC-017`
  (`spec/pflichtenheft.md:263`); Payload ist real leer (`a.conn.Publish(subject, nil)`),
  keine Positions-/Change-Daten — deckt `ADR-0055` Punkt 3 vollständig ab.
- geprüft, ohne Befund: `.a-check` real ausgeführt: „gesamt: 0 Befund(e)"
  — neue Port-/Adapter-Dateien liegen vollständig in den bestehenden
  Glob-Layern `ports`/`adapters`, kein neuer Edge nötig, wie in `ADR-0055`
  vorausgesagt.
- geprüft, ohne Befund: `docs-check` real ausgeführt: „410 Datei(en)
  geprüft, 0 Befund(e)".
- geprüft, ohne Befund: `CaptureService.Capture()` — die Notify-Fehlerbehandlung
  ist wasserdicht: der einzige Aufrufpfad zu `s.notify.Notify` liegt
  **nach** dem `return`-Punkt für Persistenz-/ACK-Fehler; sein `error`
  wird ausschließlich geloggt (`s.log.Warn`), nie in eine `return`-Anweisung
  übernommen — es gibt keinen zweiten Aufrufpfad. Regressionstest
  `TestCaptureSucceedsDespiteFailingNotification` und
  `TestCaptureNotifiesAfterAckOnSuccess` real ausgeführt (`make test`):
  beide grün, Ereignis-Reihenfolge `persist → ack → notify` belegt.
- geprüft, ohne Befund: Ordnung `Receive → Decode → Persist → COMMIT
  Store → ACK Source → Notify` — `internal/adapters/driving/replication/receive/receive.go`
  bleibt im Diff unverändert (kein neuer Fehlerpfad in `process`), wie im
  Plan (§3) vorgesehen.
- geprüft, ohne Befund: fünf bestehende `NewCaptureService`-Aufrufstellen
  (`internal/bootstrap/wiring.go:468`,
  `internal/bootstrap/replication_stream_test.go:123`,
  `internal/bootstrap/heartbeat_internal_test.go:129`, sowie der
  Test-Helper und zwei neue Aufrufstellen mit `WithChangeNotification` in
  `service_test.go`) — alle kompilieren unverändert, `make test` real
  grün über alle Pakete inkl. `internal/bootstrap`.
- geprüft, ohne Befund: `make test-notify` real ausgeführt — echter
  NATS-Testcontainer (`nats:2-alpine@sha256:065e8355…`), Test
  `TestNotifyPublishesEmptyPayloadOnSubject` grün: realer
  Publish-Erfolgsbeleg, Subjekt `cdc.changes.src-1`, Payload-Länge 0.
- geprüft, ohne Befund: `make gates`-Umfang — `test-notify` steht **nicht**
  in `GATE_CHECKS` (nur `baseline-verify`, `docs-check`,
  `commit-traceability`, `coverage-gate` tragen `GATE_CHECKS +=`) und
  nicht in `.github/workflows/ci.yml`; `harness/README.md` klassifiziert
  den neuen Eintrag korrekt als „kein Gate" — deckt sich mit `ADR-0055`
  und dem Plan.
- geprüft, ohne Befund: `tools/harness/run-notify-tests.sh` — Docker-only
  (`AGENTS.md` §3.1): kein lokales Toolchain-Install, gepinnte Digests
  für Toolchain und NATS-Image, realer TCP-Bereitschafts-Check statt
  Sleep-Raten; Kopfkommentar zitiert `slice-052`s §6-Risiko als
  Testfall-Provenienz (Subjekt: die Testinfrastruktur selbst) — dieselbe
  zulässige Zitierform wie in den vom Architect-Verdikt geprüften
  `_test.go`-Beispielen, kein Fund analog F-1.
- geprüft, ohne Befund: `go.mod`/`go.sum` — `github.com/nats-io/nats.go
  v1.53.1` als erste direkte Nicht-PostgreSQL-Abhängigkeit, wie in
  `ADR-0055` Punkt 5 angekündigt; keine unerwarteten transitiven
  Root-Pakete (nur `klauspost/compress`, `nats-io/nkeys`, `nats-io/nuid`,
  `golang.org/x/crypto`, `golang.org/x/sys` als `// indirect`).
- geprüft, ohne Befund: Commit-Traceability — `RANGE=ab69fd2..2ca089c
  make commit-traceability` real ausgeführt: „OK — 4 Commit(s) …,
  Betreffs ohne Struktur-ID"; alle vier Commit-Betreffs nennen `ADR-0055`,
  keine `SPEC-*`/`ARC-*`-Kennung im Betreff.
- geprüft, ohne Befund: Commit-Trailer — `git log -5 --format=%B` zeigt
  **keine** `Co-Authored-By:`- oder `Claude-Session:`-Trailer in den vier
  Commits.
- geprüft, ohne Befund: `AGENTS.md` §3.3 (git mv + Inhalt = zwei Commits)
  — `ab69fd2` (`next` → `in-progress`) ist ein reiner Rename (0
  Einfügungen/Löschungen laut `git show --stat`), Inhaltsänderungen liegen
  in den vier separaten Folge-Commits.
- geprüft, ohne Befund: Out-of-Scope §1 des Slice-Plans — weder
  `compose.yaml` noch `internal/bootstrap/wiring.go` (`CDC_NATS_URL`)
  noch Boundary-/Negative-Testbelege sind im Diff enthalten; alle drei
  bleiben korrekt bei `slice-053`/`054`/`055`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Slice-Chronik in Code-Kommentar ·
Code-Kommentar referenziert ephemeres Implementer-Bericht-Artefakt

## Verdikt

**Merge-blockierend:** ja — ein HIGH-Finding (F-1) blockiert die Übergabe
an den Verifier.

**Übergabe:** F-1 und F-2 gehen als Findings an den Implementer zurück
(Rückkante Review → Implementer, Modul 8) für eine Fixrunde. Beide sind
lokal begrenzte Kommentar-Korrekturen ohne Verhaltensänderung — keine
Architect-Beteiligung nötig (kein Rollen-Widerspruch, kein dritter
gleicher Konflikttyp im Sinne von Modul 8 §Konflikt-Pfad). Da eine
Fixrunde folgt, bleibt die DoD-Zeile „Review durchgeführt" im Slice-Plan
**offen** und wird regulär nach der Fixrunde nachgezogen
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde,
Umkehrschluss). Der Report ersetzt keine Verifikation — DoD-/
Spec-Konformität prüft der Verifier separat, nach Abschluss der Fixrunde.
