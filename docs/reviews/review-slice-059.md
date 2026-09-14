# Review-Report: slice-059 — 2026-09-14

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10
§Drei Review-Arten), geprüft gegen Plan/ADR/Hard Rules (Maintainability),
**nicht** gegen die DoD (Verifier-Aufgabe, Modul 11).

**Gegenstand:** `slice-059` — drei Commits auf `main`, Elter `1d0d172`
(reiner `next→in-progress`-Move):

- `7b6b253` — HTTP/JSON-Adapter-Grundgerüst (`internal/adapters/driving/http/`),
  Token-Middleware (401/403), `RegisterConsumer`-Handler, additive
  Bootstrap-Verdrahtung (`CDC_HTTP_ADDR`, `CDC_API_TOKEN_READER`,
  `CDC_API_TOKEN_ADMIN`)
- `40fe030` — `SPEC-018`-Eintrag in `spec/pflichtenheft.md`
- `8bcc8a5` — fünf Implementierungs-DoD-Punkte abgehakt

Vorausgegangener, bereits vor diesem Review-Auftrag vom
Planner-Koordinator korrigierter Commit `897809c` (`SPEC-018` im Betreff,
per `reset --soft` ersetzt durch `40fe030`) ist nicht mehr Teil des
geprüften Standes und kein Finding-Kandidat.

**Skill:** `.harness/skills/reviewer.md` @ `f07ba3d` (Stand 2026-09-13,
vier repo-spezifische HIGH-Regeln, Slice-/Wellen-Chronik-Regel und
Handbuch-Versionshistorie-Regel enthalten)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-14

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-059-http-api-adapter-grundgeruest-registerconsumer.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §6 Risiken, §8 Sub-Area) — vollständig gelesen
- `ADR-0057` (`docs/plan/adr/0057-http-grpc-api.md`) — bindend: Protokoll
  (Teilfrage 1), Umfang (Teilfrage 2), Token-Authn (Teilfrage 3),
  Adapter-Platzierung/Schichten-Konformität (Teilfrage 4)
- `LH-FA-SST-006`, `LH-FA-CON-001` (`spec/lastenheft.md`), `SPEC-018`
  (neu, `spec/pflichtenheft.md`)
- `AGENTS.md` §3 Hard Rules (insb. 3.1, 3.3, 3.7, 3.9), §5
  Traceability-Regeln
- `harness/conventions.md` (MR-000 ID-Schema)
- `.a-check.yml` (Layer-/Edge-Definition, zum Abgleich mit den
  tatsächlichen Imports des neuen Adapters)
- Reale Gate-Läufe dieses Reviews: `make gates` (baseline-verify,
  docs-check, commit-traceability, coverage-gate, a-check) — alle grün,
  Exit-Code direkt und ungefiltert geprüft (siehe Verdikt)

---

## Findings

### F-1 — Token-Middleware nutzt einfachen String-Vergleich statt konstant-zeitigem Vergleich

- `kategorie`: INFO
- `quelle`: Maintainability · `ADR-0057` Teilfrage 3 ("minimaler Aufwand
  (Header-Vergleich gegen konfigurierte Tokens …)")
- `pfad`: `internal/adapters/driving/http/middleware.go:27-38` (`classifyToken`)
- `befund`: `classifyToken` vergleicht den Aufruf-Token mit `==`
  (Nicht-konstant-zeitiger String-Vergleich), was bei einem
  Netzwerk-exponierten Bearer-Token-Schema theoretisch ein
  Timing-Seitenkanal-Risiko wäre. `ADR-0057` beschreibt exakt diesen
  einfachen Header-Vergleich als gewählte Ausgestaltung — die Umsetzung
  weicht von der Entscheidung nicht ab, sie setzt sie wörtlich um. Kein
  Konformitäts-Defekt, daher keine erwartete Aktion an diesem Slice.
- `verifizierbar`: nein — kein Gate prüft Timing-Verhalten
- `klasse`: „Timing-Seitenkanal bei Token-Vergleich, ADR-konform akzeptiert“

### F-2 — `Server.Start`/`Server.Shutdown` ohne Unit-Testabdeckung (0 %)

- `kategorie`: INFO
- `quelle`: Maintainability · `ADR-0057` §Testabdeckung ("E2E-Beleg …
  Gegenstand von Slice C"), `slice-059` §1 Out-of-Scope (E2E-Rundlauf →
  `slice-061`)
- `pfad`: `internal/adapters/driving/http/server.go:82-98`
- `befund`: Der Coverage-Report dieses Gate-Laufs weist `Start` und
  `Shutdown` mit 0 % aus — beide Pfade öffnen/schließen einen realen
  Socket und sind laut Plan/ADR bewusst erst über den E2E-Rundlauf von
  `slice-061` real belegt, nicht über die Whitebox-Tests dieses Slice
  (`httptest.NewServer` umgeht `Start` gezielt). Konsistent mit dem
  vorab dokumentierten Schnitt, keine erwartete Aktion an diesem Slice.
- `verifizierbar`: ja — `make gates` (Coverage-Ausgabe von `go test -cover`)
- `klasse`: „Geplante Coverage-Lücke, durch Folge-Slice bereits benannt“

## Negativbefunde

- geprüft, ohne Befund: `internal/adapters/driving/http/middleware.go` —
  `classifyToken` prüft `token == ""` **zuerst** und bedingungslos
  (`return roleNone` vor jedem Klassen-Vergleich); ein leer
  konfiguriertes `readerToken`/`adminToken` kann dadurch nie gegen einen
  leeren/fehlenden Aufruf-Token matchen — das in `slice-059` §6 benannte
  Risiko ist strukturell ausgeschlossen, nicht nur getestet. Real
  verifiziert über `TestClassifyTokenLeereKonfigurationTrifftKeinToken`
  (beide Fälle: leerer Token gegen leere Konfiguration, beliebiger Token
  gegen leere Konfiguration) — Test lokal nachvollzogen (Codepfad
  gelesen, nicht nur Testname vertraut)
- geprüft, ohne Befund: `internal/adapters/driving/http/middleware.go` —
  401 (kein/unbekannter Token) und 403 (bekannter, aber unzureichender
  Token) sind über `withToken` sauber getrennt (`roleNone` → 401,
  `callerRole < required` → 403); beide Fitness-Function-Fälle aus
  `ADR-0057` real getestet (`TestRegisterConsumerOhneTokenEndetMit401`,
  `TestRegisterConsumerUnbekannterTokenEndetMit401`,
  `TestRegisterConsumerReaderTokenEndetMit403`,
  `TestRegisterConsumerAdminTokenRegistriert`)
- geprüft, ohne Befund: `internal/adapters/driving/http/server.go`,
  `registerconsumer.go` — Imports ausschließlich
  `application/port/inbound`, `application/port/outbound` (`LogPort`,
  identisches Muster wie im bestehenden Driving-Adapter
  `internal/adapters/driving/replication/receive`), `domain/errors`,
  `domain/model`; keine Driven-Adapter-Interna, keine
  Application-Interna — exakt `ADR-0057` Teilfrage 4
- geprüft, ohne Befund: `.a-check.yml` — unverändert (kein Diff);
  `make a-check` real ausgeführt, 0 Befunde, der neue Adapter liegt
  vollständig im bestehenden `adapters`-Glob wie von `ADR-0057`
  vorausgesagt
- geprüft, ohne Befund: `internal/bootstrap/wiring.go` — die drei neuen
  Umgebungsvariablen sind in `ConfigFromEnv` rein additiv gelesen, keine
  der sechs bestehenden Vorbedingungs-Prüfungen (`CaptureDSN` …
  `Slot`) wurde berührt; `Run()` konstruiert `apihttp.Server` nur
  innerhalb von `if cfg.HTTPAddr != ""`, geordneter Shutdown vor
  Rückkehr aus `Run` (`httpServer.Shutdown` + `httpDone.Wait()`, analog
  zum bestehenden `administrationDone`/`retentionDone`-Muster)
- geprüft, ohne Befund: `internal/bootstrap/wiring_test.go` —
  `TestConfigFromEnvHTTPAddrBleibtOptional` belegt real (nicht nur
  behauptet), dass eine vollständige Bestandsverdrahtung ohne die drei
  neuen Variablen unverändert durch `ConfigFromEnv` läuft und die neuen
  Felder leer bleiben; Regressionstest-Namensgebung konsistent mit dem
  DoD-Wortlaut
- geprüft, ohne Befund: `spec/pflichtenheft.md` `SPEC-018` — per
  `grep`/Sichtprüfung real bestätigt: **keine** `ADR-*`-Referenz
  innerhalb des Eintrags (nur `LH-*`-Bezüge im Fließtext), damit
  konform zum Spec→ADR-Referenzverbot (Rang-Disziplin der drei
  Spec-Straten); Fehler-Antwortform, Endpunkt, JSON-Schema und
  Token-Header-Form im Eintrag stimmen mit der tatsächlichen
  Implementierung überein (Statuscodes, Feldnamen `consumer_id`/`name`/
  `already_registered`/`error`)
- geprüft, ohne Befund: `spec/pflichtenheft.md` §6 „Externe Verträge“ —
  bewusst nicht um `SPEC-018` ergänzt; dieselbe Abgrenzung wie bei den
  bestehenden CLI-/SQL-Zugriffswegen (auch `cdc.changes`, CLI-Register
  sind dort nicht gelistet) — §6 führt Verträge zu externen Systemen,
  auf die diese Anwendung als Client zugreift (PostgreSQL, OCI, NATS),
  nicht die von dieser Anwendung selbst angebotenen Schnittstellen
- geprüft, ohne Befund: Scope-Treue gegen `slice-059` §1 — `git diff
  1d0d172 8bcc8a5 --stat` bestätigt: keine Änderung an `compose.yaml`,
  `harness/README.md`, `tools/harness/run-integration-tests.sh`,
  `test-store`/`test-replication`-Skripten; kein zusätzlicher Endpunkt
  über `POST /consumers` hinaus; keine gRPC-Artefakte
- geprüft, ohne Befund: Commit-Betreffs aller drei Commits — keine
  `SPEC-*`/`ARC-*`-ID im Betreff, jede Message trägt mindestens eine
  `LH-*`- oder `ADR-*`-Kennung; keine `Co-Authored-By:`- oder
  Session-Trailer
- geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7) in allen
  neuen/geänderten Go-Dateien — Zusage-/Kopplungs-/Grenz-Kommentare mit
  `ADR-*`/`LH-*`-Bezug; keine Slice-Chronik im Produktionscode (einzige
  `slice-059`-Erwähnungen stehen in Godoc-Kommentaren, die den *aktuellen
  Umfang* dieses Slice als Grenze benennen — `Config.RegisterConsumer`,
  `slice-059 §1` —, nicht eine Vorher/Nachher-Begründung; das
  Test-Godoc `TestClassifyTokenLeereKonfigurationTrifftKeinToken` zitiert
  „§6-Risiko dieses Slice“ als zulässige Testfall-Provenienz, Subjekt ist
  der Test); keine konjunktivische Aussage über eine verworfene
  Alternative, kein abwesender Text
- geprüft, ohne Befund: reale Sensor-Läufe dieses Reviews — `make gates`
  ungefiltert ausgeführt: `baseline-verify` (v6.5.0 OK, 54 Dateien),
  `docs-check` (471 Dateien, 0 Befunde, zweimal — Vollset und
  Commit-Modul), `commit-traceability` (OK, 5 Commits im Fenster,
  Betreffs ohne Struktur-ID), `a-check` (0 Befunde, unveränderter
  Abdeckungs-Hinweis für `tools/harness/natssub/main.go`),
  `coverage-gate` (41,90 % ≥ 35 %) — Exit-Code unmittelbar nach dem
  ungepipten `make gates`-Aufruf geprüft: `0`

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Timing-Seitenkanal bei Token-Vergleich,
ADR-konform akzeptiert · Geplante Coverage-Lücke, durch Folge-Slice
bereits benannt

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, 0 LOW. Beide INFO-Findings
sind bereits durch `ADR-0057` bzw. den vorab dokumentierten Slice-Schnitt
(`slice-059` §1 Out-of-Scope, Folge-Slice `slice-061`) abgedeckt und lösen
keine Aktion an diesem Slice aus.

**Übergabe:** Kein Fixrunden-Pfad zum Implementer nötig. Nach
`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde wird die
DoD-Zeile „Review durchgeführt, Report unter `docs/reviews/` liegt vor“ im
Slice-Plan in diesem Commit mitgezogen. Die Finding-Klassen gehen
zusätzlich in die Slice-Closure §7 und von dort in den
Beobachtungs-Register-Zähler. Dieser Report selbst ist ein **Lauf-Beleg**
(Audit: dieser Diff, dieser Skill, dieses Modell, dieses Verdikt) und wird
über Läufe hinweg nicht wieder gelesen. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11).
