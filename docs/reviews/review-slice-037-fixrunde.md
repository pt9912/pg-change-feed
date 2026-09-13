# Review-Report: slice-037 (Fixrunde) — 2026-09-13

**Review-Art:** Code — gezielte Bestätigungsprüfung der vier Findings aus
`docs/reviews/review-slice-037.md` (Modul 10), **kein** vollständiges
Re-Review des Slice. Geprüft gegen den ursprünglichen Befund-Text, `ADR-0050`
(zentraler Prüfmaßstab) und `AGENTS.md` §3 Hard Rules.

**Gegenstand:** Commit `c78aa1d` (Fixrunde auf `2a771bc`/`b763253`).

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext:**

- `docs/reviews/review-slice-037.md` (F-1 bis F-4, vollständig)
- `git show c78aa1d` (vollständiger Diff, alle sieben geänderten Dateien)
- `internal/domain/model/schema_version.go` (Vergleichsmuster für F-4)
- `internal/bootstrap/walretention_internal_test.go` (Vergleichsmuster für F-2)
- `git log --oneline -10`, `git show 48d03ba --stat` (Kollisionsprüfung
  gegen `ADR-0051`)

---

## F-1 — `//nolint:gosec`-Suppression (HIGH)

- **Verdikt: behoben.**
- `grep -rn "nolint" --include="*.go" internal/` liefert real keinen
  Treffer mehr (vorher genau einer).
- Diff in `mapper_test.go`: Schleifenzähler `i` läuft jetzt selbst als
  `uint32` (`for i := uint32(0); i < iterations; i++`), `xid := i + 1`
  ohne Konvertierung. `iterations` ist eine untypisierte Konstante (`= 200`),
  der Vergleich `i < iterations` bleibt gültig. Die zweite betroffene
  Zeile (`Commit{CommitLSN: uint64(i + 1), ...}`) wurde konsistent auf
  `uint64(xid)` umgestellt — `xid` ist bereits `i + 1`, also derselbe Wert,
  kein Verhaltensunterschied im Testfall.
- Kein Overflow-Risiko bei 200 Iterationen weit unter jedem `uint32`-Rand.
- `verifizierbar`: ja — `grep -rn "nolint" --include="*.go" internal/`, real
  ausgeführt, 0 Treffer.

## F-2 — fehlende Whitebox-Testabdeckung (MEDIUM)

- **Verdikt: behoben.**
- Neue Datei `internal/bootstrap/administration_internal_test.go` (Muster
  `walretention_internal_test.go`/`heartbeat_internal_test.go` bestätigt,
  `package bootstrap`, kein externer Testpfad). Fakes für
  `AdministrationRequestPort`, `administrationListener`,
  `EnableTableUseCase`, `DisableTableUseCase`, `TableActivationPort`,
  `SchemaStorePort`.
- Sechs neue Testfälle decken genau die im vorigen Review benannten
  Lücken: Enable-Happy-Path mit realer `Assembler`-Bindung
  (`assemblerCapturesQualified` über die öffentliche `Consume`-API, kein
  Zugriff auf unexportierte Felder), Disable-Gegenstück, `MarkFailed`-Zweig,
  unbekannter `Kind` (`default`-Zweig), ctx-Cancel-Austritt aus
  `runAdministration`, Listener-Fehler-Zweig mit Fallback-Poll.
- Real reproduziert: `go test -race -v -run
  "TestProcessAdministrationRequests|TestApplyAdministrationRequest|TestRunAdministration"
  ./internal/bootstrap/...` — alle sechs Fälle `PASS`, kein Data-Race.
- Der Kommentar an `administrationListener` in `wiring.go` verspricht
  „dasselbe Whitebox-Test-Muster wie `walRetentionMeasurer`" — die neue
  Datei löst das ein: `fakeAdministrationListener` implementiert exakt das
  Interface, keine reale PostgreSQL-Instanz nötig.
- `verifizierbar`: ja — realer Testlauf oben, sowie vollständiger
  `make test` (Race-Detector), grün.

## F-3 — kein Backoff im `LISTEN`-Reconnect (MEDIUM)

- **Verdikt: behoben.**
- `AdministrationListener` trägt neu `reconnectBackoff time.Duration`.
  `WaitForNotification`: bei Erfolg oder erfolgreichem Reconnect wird sie
  auf 0 zurückgesetzt; bei gescheitertem Reconnect verdoppelt
  `nextAdministrationReconnectBackoff` sie (Start 200ms, Deckel 30s) und
  die Rückkehr wartet über `select { case <-time.After(...): case
  <-ctx.Done(): }` — ein `ctx`-Abbruch unterbricht das Warten wirklich,
  statt es abzuwarten.
- Logik der reinen Funktion korrekt geprüft: `current <= 0` → Initial-Wert
  (auch für negative Ausgangswerte abgesichert, per Testfall belegt);
  sonst Verdopplung, gedeckelt am Maximum. Kein Overflow-Risiko, da
  `current` durch den Deckel strukturell nie über `administrationReconnectMaxBackoff`
  (30s) hinauswächst, bevor verdoppelt wird.
- Neue reine Unit-Test-Datei
  `internal/adapters/driven/postgresstorage/administration_backoff_internal_test.go`
  (Pfad weicht vom Commit-Text ab, der `internal/bootstrap/` nennt —
  tatsächlich liegt sie im `postgresstorage`-Paket, wo die Backoff-Konstanten
  auch definiert sind; inhaltlich passend). Sechs Tabellenfälle inklusive
  Randfälle (negativer Ausgangswert, knapp unter/an/über der Obergrenze).
  Real reproduziert: `go test -race -v -run
  "TestNextAdministrationReconnectBackoff"
  ./internal/adapters/driven/postgresstorage/...` — `PASS`, alle
  Unterfälle grün.
- `verifizierbar`: ja — realer Testlauf oben.

## F-4 — `AdministrationRequest` ohne validierenden Konstruktor (MEDIUM)

- **Verdikt: behoben.**
- Neuer Konstruktor `model.NewAdministrationRequest(id, source, schema,
  table, kind) (AdministrationRequest, error)` in
  `internal/domain/model/administrationrequest.go` — Mustervergleich mit
  `NewSchemaVersion` (`schema_version.go`) bestätigt: gleiche Leer-Prüfung
  über `domainerrors.ErrEmptyIdentifier`, gleiche Struktur (Validierung,
  dann Struct-Literal, `nil`-Fehler). Zusätzlich ein neuer Sentinel
  `domainerrors.ErrInvalidAdministrationRequestKind` für die geschlossene
  Menge `enable`/`disable` (dieselbe Fehlerklasse-Konvention wie die
  übrigen Konstruktoren, die bei einer verletzten Invariante einen eigenen
  Sentinel statt eines generischen Fehlers zurückgeben).
  `verifizierbar`: ja — realer Testlauf oben.
- Aufrufstelle real umgestellt, nicht nur definiert: `ListPending`
  (`internal/adapters/driven/postgresstorage/administrationrequest.go`)
  baut den Satz jetzt über `model.NewAdministrationRequest(...)` statt
  über ein Struct-Literal und reicht einen Konstruktor-Fehler nach oben
  durch. `grep -rn "model.AdministrationRequest{" internal/ --include="*.go" |
  grep -v _test.go` liefert real keinen Treffer — kein verbliebener
  Bypass am Domain-Core-Rand in Produktionscode.
- INFO (kein neues Finding, nur Beobachtung ohne erwartete Aktion): Es
  gibt jetzt zwei gleichnamige `NewAdministrationRequest`-Funktionen in
  zwei Paketen (`model.NewAdministrationRequest` — Domain-Konstruktor,
  neu; `postgresstorage.NewAdministrationRequest` — Adapter-Konstruktor,
  Bestand seit `2a771bc`). Beide sind durch ihr Paket eindeutig
  qualifiziert, kein Compile- oder Laufzeitproblem; an den beiden
  Aufrufstellen (`wiring.go:426`, `administrationrequest.go:83`) ist die
  Zuordnung im Kontext eindeutig lesbar.

## Regressionsprüfung

- `make test` (Race-Detector, vollständige Suite, real ausgeführt):
  alle Pakete `ok`, einschließlich `internal/bootstrap`,
  `internal/adapters/driven/postgresstorage`,
  `internal/adapters/driving/replication/mapper`.
- `make gates` (real ausgeführt): `baseline-verify` (54 Dateien),
  `d-check` (304 Dateien, 0 Befunde), `commit-traceability` (5 Commits,
  „Betreffs ohne Struktur-ID" — OK), `a-check` (0 Befunde) — alle grün.
- Kollisionsprüfung mit `ADR-0051` (Commit `48d03ba`, zwischen
  Review-Report und Fixrunde committet): `git show 48d03ba --stat` zeigt
  ausschließlich `AGENTS.md`, `docs/plan/adr/0051-*.md`,
  `docs/plan/adr/README.md`. Die Fixrunde `c78aa1d` berührt ausschließlich
  Slice-Plan, Domain-, Adapter- und Bootstrap-Testdateien — keine
  überlappenden Pfade, kein Merge-Konflikt-Risiko, keine inhaltliche
  Berührung (CI/CD-Pipeline vs. Administrations-Goroutine-Fixes).

## Negativbefunde

- geprüft, ohne Befund: kein neuer `//nolint`/`#noqa`/`[SuppressMessage]`
  im gesamten Diff außer der entfernten Zeile.
- geprüft, ohne Befund: keine der drei neuen Testdateien führt einen
  Netzwerk- oder DB-Zugriff (reine Unit-/Whitebox-Tests mit Fakes) —
  konsistent mit dem versprochenen „ohne reale PostgreSQL-Instanz".
- geprüft, ohne Befund: keine Änderung an einer `Accepted`-ADR; `ADR-0050`
  bleibt inhaltlich unverändert (nur der Slice-Plan referenziert die
  Fixrunde neu).
- geprüft, ohne Befund: Commit-Betreff trägt `ADR-0050`, kein
  `SPEC-*`/`ARC-*` im Betreff.

## Summary

| Finding | Verdikt |
|---|---|
| F-1 (HIGH) | behoben |
| F-2 (MEDIUM) | behoben |
| F-3 (MEDIUM) | behoben |
| F-4 (MEDIUM) | behoben |

**Neue Findings dieser Fixrunde:** keine (eine INFO-Beobachtung ohne
erwartete Aktion, siehe F-4).

## Verdikt

**Merge-blockierend:** nein — alle vier Findings aus
`review-slice-037.md` sind real verifiziert behoben, keine Regression in
`make test` (Race) oder `make gates`, keine Kollision mit der
zwischenzeitlich committeten `ADR-0051`.

**Übergabe:** Slice-037 kann aus Reviewer-Sicht zur Verifikation
(Modul 11) weitergereicht werden. Dieser Report ist ein Lauf-Beleg und
wird über Läufe hinweg nicht wieder gelesen; er ersetzt keine
Verifikation gegen DoD/Spec.
