# Review-Report: slice-005 Implementer-Diff — 2026-09-09

**Review-Art:** Diff-Review (Implementer-Range `cb021dc..4ca6658`, 6 Implementer-Commits
+ 2 vorgelagerte Planner-Commits des Ranges) — *wogegen*: Slice-Plan §1/§2/§3/§6
(Plan-Treue, Bewertung der vier gemeldeten Entscheidungen), ADR-Bezüge
([`ADR-0006`](../plan/adr/README.md), 0007, 0008, 0011, 0012, 0016, 0023, 0026, 0027,
0030, 0032, 0042), Hard Rules (`AGENTS.md` §3.1 Docker-only, §3.7
Kommentar-Klassen, §5 Dokumentations-Regeln), Traceability (LH-*/ADR-* je
Commit, keine Struktur-IDs, keine superseded-Referenzen als tragende Anker),
neue Angriffsfläche (erster Replication-Protokoll-Code: SQL-Katalogabfragen,
Bezeichner-Interpolation, Keepalive-/ACK-Positionen, Testcontainer-Hygiene,
Digest-Pins), Test-Qualität mit eigenen Mutations-Proben, Implementer-Risiken
(a)–(d), [`ADR-0043`](../plan/adr)-Abgrenzung, Zweitschreiber-Frage. Keine DoD-Prüfung — das
ist der Verifier (Modul 11).

**Gegenstand:** `d1a84eb` (pglogrepl) · `1e4d89f` (decode + mapper) · `49ea22f`
(receive + postgresack + Port) · `cf3c57a` (Replikations-Tests, test-replication,
image-hash-Renewal) · `bffcea0` (Verlagerung Verdrahtungs-Test nach
`internal/bootstrap`) · `4ca6658` (Kommentar-Korrekturen) — im Range zusätzlich
`3125c0a`/`31e1b41` (Planner: [`SPEC-015`](../../spec/pflichtenheft.md) + [`ADR-0043`](../plan/adr)-Plan-Nachzug, vor den
Implementer-Commits; Bewertung siehe Zweitschreiber-Abschnitt).

**Skill:** `.harness/skills/reviewer.md` (Stand 2026-09-09, geschärft: vier
repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen) · Gerüst:
`docs/reviews/review-report.template.md` (Form wie `review-slice-004.md`).

**Modell:** Claude Code (glm-5.3-flash) · **Datum:** 2026-09-09

**Eingangs-Kontext:**

- Diff `git diff cb021dc..4ca6658` (18 Dateien, +2262/−10: `decode/`+`mapper/`+
  `receive/` inkl. Tests, `postgresack/`, `replicationack.go` (Port),
  `replication_stream_test.go` (bootstrap), `go.mod`/`go.sum`, Makefile,
  `tools/harness/run-replication-tests.sh`, harness-README-Zeile,
  `harness/image-hash.txt`, Planning-Dateien der Planner-Commits)
- `docs/plan/planning/in-progress/slice-005-replication-stream-adapter.md`
  (§1–§8, am Range-Start `cb021dc` und am Range-Head inkl. der
  Planner-Berührung `31e1b41`)
- `docs/plan/adr/README.md` · [`ADR-0006`](../plan/adr/README.md) · [`ADR-0007`](../plan/adr/README.md) · [`ADR-0008`](../plan/adr/README.md) · [`ADR-0011`](../plan/adr/README.md) · [`ADR-0012`](../plan/adr/README.md) · [`ADR-0016`](../plan/adr/README.md) ·
  [`ADR-0023`](../plan/adr/README.md) · [`ADR-0026`](../plan/adr/README.md) · [`ADR-0027`](../plan/adr/README.md) · [`ADR-0030`](../plan/adr/README.md) · [`ADR-0032`](../plan/adr/README.md) · [`ADR-0042`](../plan/adr/README.md) · [`ADR-0043`](../plan/adr/README.md) — Status:
  Superseded sind [`ADR-0038`](../plan/adr/README.md) (→ 0039) und [`ADR-0039`](../plan/adr/README.md) (→ 0042,
  Rest-Fortgeltung ausdrücklich); Proposed sind [`ADR-0021`](../plan/adr/README.md), 0022
- `spec/lastenheft.md` ([`LH-FA-CAP-001`](../../spec/lastenheft.md)…008, [`LH-FA-CAP-004`](../../spec/lastenheft.md)/006,
  [`LH-QA-REL-001`](../../spec/lastenheft.md), [`LH-QA-REL-004`](../../spec/lastenheft.md)), `spec/pflichtenheft.md`
  ([`SPEC-001`](../../spec/pflichtenheft.md)/002/003/008/010, [`LH-FA-CAP-006.a`](../../spec/pflichtenheft.md), [`LH-FA-SCH-004.a`](../../spec/pflichtenheft.md)/005,
  Fehlerklassen-Tabelle §4),
  `spec/architecture.md` ([`ARC-003`](../../spec/architecture.md)/005/006)
- Commit-Messagen der Range; `AGENTS.md` §3/§5; `harness/conventions.md`
  (MR-000); `harness/README.md` §Sensors (Werkzeuge-Zeile `test-replication`:
  kein Gate; `image`-Zeile: Beleg gilt am HEAD); `96c47af` (Beleg-Semantik-Formel);
  `.a-check.yml` (composition_root: bootstrap, cmd)
- Stand-alone-Prüfung am Range-Head im gepinnten Toolchain-Container
  (`golang:1.27-alpine@sha256:cf6fca…`, docker-only): `gofmt -l` leer,
  `go vet ./...` ohne Befund, treiberfreie `go test ./...` grün; **`make
  test-replication` am Range-Head grün** (Stream-Tests gegen reale
  `postgres:18-alpine` mit `wal_level=logical`, Verdrahtungs-Test in
  `internal/bootstrap`); drei Gates am Range-Head (`baseline-verify` OK 54
  Dateien, `d-check` 92 Dateien/0 Befunde, `a-check` 0 Befunde); drei
  Mutations-Proben in Worktrees außerhalb des Arbeitsbaums (Tabelle unten);
  Binary-Vergleich an Range-Basis und Range-Head mit den Dockerfile-Flags
  (bit-identisch, siehe F-7); Re-Build-Probe `make image` am Range-Head
  (Digest reproduziert `9ac4a9…`)

---

## Findings

### F-1 — Commit `4ca6658` trägt keine `LH-*`-/`ADR-*`-Kennung (erste Auftreten der direkten Form)

- `kategorie`: HIGH
- `quelle`: `harness/README.md` §Traceability rules („PRs/Commits **müssen**
  mindestens eine `LH-*` oder `ADR-*`-ID nennen") · Skill-Klasse
  „Traceability-/ID-Schema-Verstoß" — direkte Form: Commit nennt keine
  Kennung · `AGENTS.md` §5
- `pfad`: Commit `4ca6658` (Subject „docs(comments): Slice-Bezug aus
  Zustandskommentaren entfernt"; Body nennt nur `AGENTS.md §3.7` — keine
  `LH-*`- oder `ADR-*`-Kennung)
- `befund`: Fünf der sechs Implementer-Commits tragen je mindestens eine
  `LH-*`/`ADR-*`-Kennung; der Kommentar-Korrektur-Commit nennt nur eine
  Hard-Rule-Adresse. Die Hard-Rule-Nennung ist kein Ersatz: Die
  Traceability-Regel verlangt die Kennung des berührten Vertrags, und die
  berührten Kommentare tragen `LH-FA-SCH-004.a`/`LH-FA-CFG-001.a` (Code-Seite).
  Erste Auftreten der *direkten* Form in diesem Repo (Vorgänger slice-002 F-5
  war die schwächere Form „ID im Code, Commit trägt sie nicht", LOW — der
  Commit trug dort andere Kennungen). Die Message bleibt in der Historie
  unveränderlich; die Korrektur wirkt nur vorwärts.
- `verifizierbar`: ja — `git log --format=%B 4ca6658` gegen
  `harness/README.md` §Traceability rules
- `klasse`: Commit ohne Vertrags-Kennung (1. Auftreten der direkten Form)

### F-2 — Plan-Erweiterung ohne Plan-Nachzug (5. Auftreten — laufende Sequenz)

- `kategorie`: MEDIUM
- `quelle`: Slice-Plan §1/§3/§6 · Modul 5 („Wer später mitnimmt …, hat den
  Plan **geändert**, nicht nur ergänzt") · laufende Konflikt-Sequenz (Modul 8
  — seit review-slice-003 F-1, 4. Auftreten in review-slice-004 F-1)
- `pfad`: `docs/plan/planning/in-progress/slice-005-replication-stream-adapter.md`
  (§3 unverändert — zwei Adapter-Zeilen aus dem Plan-Entwurf, keine Erweiterung
  im Range; §6 ohne neue Risiken) gegen `internal/application/port/outbound/
  replicationack.go` (neuer öffentlicher Port-Kontrakt + Sentinel),
  `internal/bootstrap/replication_stream_test.go`, `tools/harness/
  run-replication-tests.sh`, `Makefile:50-51`, `harness/README.md` (Werkzeuge-
  Zeile), Commit-Messagen der Range
- `befund`: Von den vier gemeldeten Entscheidungen haben drei keinen
  Plan-Träger: **(b) Keepalive-Position-Regel** — die Zusage „Keepalive meldet
  ausschließlich die bestätigte Position" ist eine Konkretisierung des
  §2-DoD-Punkts (ACK nur nach Persistenz), deckt sich mit keinem §1-Ausschluss
  und berührt keinen — *Plan-Ergänzung, nicht Plan-Änderung*; sie lebt aber
  nur im Handoff und Commit-Text. **(c) `BindCapture`** (Verbindungsaufbau
  getrennt von Port-Verdrahtung) — ebenfalls nur im Handoff; §3 nennt keine
  Verdrahtungs-Zeile. **(d) Verlagerung des Verdrahtungs-Tests nach
  `internal/bootstrap`** — eine Grenzziehung, die nur in `bffcea0` steht; §3
  kennt `internal/bootstrap` nicht. Der neue Port-Kontrakt
  (`replicationack.go`) liegt außerhalb beider §3-Zeilen-Pfade und ist eine
  Vertragserweiterung am Outbound-Port (slice-004-(c)-Muster). *(a) pglogrepl*
  ist dagegen Plan-konform: die Treiber-Wahl ist Infrastrukturdetail am
  Adapter ([`ADR-0032`](../plan/adr/README.md)), die §3-Zeile „Stream-Adapter je
  [`ADR-0006`](../plan/adr/README.md)" trägt sie; keine go.mod-Zeile entspricht der
  slice-004-Präzedenz. Die Klasse steht damit beim **fünften** Auftreten — die
  Konflikt-Sequenz läuft weiter; die Nachzüge (b/c/d) gehen als Übergabe-
  Artefakte an den Planner.
- `verifizierbar`: ja — Datei-Menge und Entscheidungen gegen die §3-Tabelle
- `klasse`: Plan-Erweiterung ohne Plan-Nachzug (5. Auftreten — Sequenz läuft)

### F-3 — Keepalive-Zusage (`LH-QA-REL-001.a`) ohne Test-Träger (Mutation belegt)

- `kategorie`: MEDIUM
- `quelle`: `LH-QA-REL-001.a` (Fehlermodi: „kein stilles Überspringen";
  confirmed_flush_lsn rückt nur über bestätigte Positionen) · Maintainability
  (Zusage des kritischen Pfads ohne Negativtest)
- `pfad`: `internal/adapters/driving/replication/receive/receive.go:77-82`
  (`lastAcked`-Mechanik), `receive.go:287-304` (Keepalive-Zweig) · keine
  Keepalive-Berührung in `receive/stream_test.go`,
  `decode/decode_test.go`, `mapper/mapper_test.go`,
  `internal/bootstrap/replication_stream_test.go`
- `befund`: Die Zusage „confirmed_flush_lsn rückt nie durch den Empfangsstand"
  (das Datenverlust-Fenster des Persist-before-ACK) trägt kein Test. Eigene
  Mutations-Probe: `WALWritePosition: s.lastAcked` → `0` in der Keepalive-
  Antwort — beide Test-Ebenen (`make test`, `make test-replication` inkl.
  Verdrahtungs-Test) bleiben **grün**. Der Implementer meldet das Risiko
  selbst; es steht aber in keinem Plan-§6-Ausgang.
- `verifizierbar`: ja — Mutation am Keepalive-Zweig, `make test-replication`
  bleibt grün; `grep -n keepalive internal/adapters/driving/replication/
  receive/stream_test.go` ohne Treffer
- `klasse`: Zusage am kritischen Pfad ohne Test-Träger

### F-4 — Restart auf bestehendem Slot: zentraler ADR-0012-Pfad ohne Beleg und ohne Plan-Träger

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0012`](../plan/adr/README.md)/[`ADR-0011`](../plan/adr/README.md) („Neustart setzt ohne Datenlücke
  fort") · `LH-QA-REL-001.a` Fehlermodi („Nach Crash wird anhand des
  Replication Slots und persistierter Zustände fortgesetzt") · Maintainability
- `pfad`: `internal/adapters/driving/replication/receive/receive.go:198-231`
  (`ensureSlot` — der bestehende-Slot-Zweig liest `confirmed_flush_lsn` als
  Startposition) gegen `receive/stream_test.go` (alle Tests laufen über den
  Neu-Anlege-Zweig; kein Test übt `confirmed_flush_lsn > 0`) und beide
  Plan-Dateien (slice-005 §6 nennt das Risiko nicht; slice-006 nennt weder
  Restart noch Slot)
- `befund`: Der bestehende-Slot-Zweig ist der [`ADR-0012`](../plan/adr)-Träger (At-Least-Once
  lebt von diesem Pfad) und läuft im gesamten Range ungetestet; der
  Implementer meldet ihn als Risiko mit slice-006-Verweis, aber slice-006
  nimmt ihn nicht an (kein Restart-/Slot-Beleg in dessen Ziel, DoD oder
  Risiken). Damit ist der Pfad aktuell ohne Träger und ohne terminierten
  Ausgang.
- `verifizierbar`: ja — Test-Menge gegen `ensureSlot`; grep in slice-006
- `klasse`: Kritischer Pfad ohne Beleg, Risiko ohne Ausgangs-Träger

### F-5 — Struktur-IDs in Commit-Messagen (4. Auftreten — ausgelöste Sequenz läuft nicht)

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §5 („Struktur-IDs (`SPEC-<NNN>`, `ARC-<NNN>`) …
  gehören nicht in die Commit-Message") · Zählstand der Übergabe: die Klasse
  zählt seit review-slice-003 F-3 beim 3. Auftreten (Sequenz-Pflicht ausgelöst)
- `pfad`: Commit `1e4d89f` (Body: „`SPEC-008`", „`SPEC-002`", „`SPEC-003`") ·
  Commit `49ea22f` (Body: „`SPEC-008`") — zusätzlich die zwei Planner-Commits
  des Ranges (`3125c0a`, `31e1b41`: `SPEC-015`/`SPEC-011` im Subject/Body)
- `befund`: Zwei Implementer-Commits des Ranges tragen `SPEC-*`-Kennungen in
  der Message; Struktur-IDs gehören in Code-Kommentar (dort zutreffend als
  Rang-Zeiger) und Spec. Mit dem Übergabe-Zählstand ist dies das **vierte**
  Auftreten — und die in review-slice-003 als Pflicht ausgelöste
  Konflikt-Sequenz (Konventions-Nachzug im Implementer-Briefing über den
  Architect) hat nicht sichtbar gelaufen (`AGENTS.md` §5 unverändert, kein
  Konventionseintrag). Die Messagen bleiben historisch; Korrektur wirkt nur
  vorwärts.
- `verifizierbar`: ja — `git log --format=%B cb021dc..4ca6658` gegen
  `AGENTS.md` §5
- `klasse`: Struktur-ID in Commit-Message (4. Auftreten — Sequenz läuft nicht)

### F-6 — Datei-Abschluss: `go.mod`/`go.sum` ohne Zeilenumbruch (3. Auftreten — MEDIUM-Stufe)

- `kategorie`: MEDIUM
- `quelle`: Maintainability — Muster war zweimal LOW (F-6 review-slice-001,
  F-5 review-slice-004); Skill-Regel „Wiederholung eines Musters, das schon
  zweimal LOW war"
- `pfad`: `go.mod` (letzte Zeile `)` ohne `\n`, eingeführt in `d1a84eb` —
  `git show cb021dc:go.mod` endet mit Umbruch) · `go.sum` (ebenfalls ohne
  Umbruch)
- `befund`: Drittens dieselbe Klasse: die deps-Commits lassen die beiden
  Lock-Files ohne abschließenden Zeilenumbruch zurück, während alle neun
  neuen Go- und Shell-Dateien des Ranges sauber abschließen. Die Klasse
  erreicht damit die MEDIUM-Stufe des Skills.
- `verifizierbar`: ja — `tail -c1 go.mod | od -c`
- `klasse`: Datei-Abschluss (3. Auftreten — MEDIUM)

### F-7 — Image-Beleg: die `96c47af`-Semantik-Formel ist am Verhalten des Ranges widerlegt

- `kategorie`: MEDIUM
- `quelle`: `96c47af` (deklarierte Beleg-Semantik: „deps-Layer-Änderungen
  ohne Binary-Änderung lassen ihn unverändert") · `harness/README.md`
  Werkzeuge-Zeile `make image` („der Beleg gilt am HEAD") · [`ADR-0039`](../plan/adr/README.md)
  (Reproduzierbarkeits-Anker)
- `pfad`: `harness/image-hash.txt` (`447eab…` → `9ac4a9…` in `cf3c57a`) gegen
  `Dockerfile:5-10` (Semantik-Formel) und den Binary-Vergleich beider
  Range-Grenzen
- `befund`: Probe: das Binary (mit exakt den Dockerfile-Flags
  `-trimpath -ldflags="-s -w"`, `CGO_ENABLED=0`) ist an Range-Basis und
  Range-Head **bit-identisch** (`sha256 43c3aec…` in beiden Builds) — und der
  exportierte Digest wechselte trotzdem (`447eab` → `9ac4a9`). Die Formel
  „deps-Änderungen ohne Binary-Import lassen den Digest unverändert" trifft
  den Mechanismus nicht; der Beleg selbst reproduziert am HEAD (Re-Build-
  Probe lief exakt `9ac4a9…`). Der Konflikt zwischen Deklaration
  (Dockerfile-Kommentar, Werkzeuge-Zeile) und beobachtetem Verhalten ist
  unentschieden — Zuständigkeit **Architect/Planner** (Semantik-Korrektur ist
  eine Entscheidung, kein Implementer-Schritt). Das Renewal in `cf3c57a`
  (Re-Build + Commit) war konservativ korrekt. Der Implementer-Risiko-Meldung
  zufolge (c): die Binary-These ist probe-bestätigt, die
  Provenanz-Varianz-These auf dem Review-Daemon **nicht** reproduziert
  (zwei Builds, gleicher Digest) — bleibt offen, welcher Builder variiert.
- `verifizierbar`: ja — Binary-Builds an beiden Range-Grenzen (bit-gleich),
  `make image` am Range-Head (Digest reproduziert)
- `klasse`: Beleg-Semantik-Formel widerlegt (Deklaration vs. Verhalten)

### F-8 — Zweites BEGIN überschreibt still die offene Transaktion

- `kategorie`: MEDIUM
- `quelle`: `LH-QA-REL-001.a` Fehlermodi („kein stilles Überspringen") ·
  `SPEC-008`/`ADR-0023` („nicht sicher interpretierbar → sichtbarer Fehler") ·
  Maintainability (fehlender Negativtest)
- `pfad`: `internal/adapters/driving/replication/mapper/mapper.go:93-98`
  (`case decode.Begin:` setzt `a.open` ohne Prüfung auf bereits offene
  Transaktion)
- `befund`: Ein BEGIN bei bereits offener Transaktion — derselbe
  Stream-Vertragsverstoß-Typ wie `ErrChangeWithoutBegin`/
  `ErrCommitWithoutBegin` — verwirft die offene Transaktion samt ihrer
  gesammelten Changes still, statt als sichtbarer Fehler zu enden. Die
  Paket-Disziplin (Kommentar-Header, Protokollverstoß-Klassen) nennt genau
  diesen Fall als fehlerpflichtig; kein Test übt den Doppel-Begin.
- `verifizierbar`: ja — Test mit zwei aufeinanderfolgenden BEGINs (fehlt);
  `grep -n "a.open = &openTransaction" mapper/mapper.go` gegen den
  Fehlerzweig
- `klasse`: Stiller Vertragsverstoß am Spec-Rand (kein Negativtest)

### F-9 — `postgresack` ohne direkte Tests

- `kategorie`: LOW
- `quelle`: Maintainability (fehlende Negativtests bei neuem öffentlichem
  Vertrag — `ReplicationAckPort`/`outbound.ErrReplication`) · [`ADR-0030`](../plan/adr/README.md)
- `pfad`: `internal/adapters/driven/postgresack/ack.go` — kein `*_test.go` im
  Paket (`go test ./...` meldet „[no test files]"); die IsZero-Grenze
  (`ack.go:55-57`) und der Fehlerpfad (`SendStandbyStatusUpdate` →
  `replicationFailure`) sind nur über den Happy-Pfad des
  Verdrahtungs-Tests belastet
- `befund`: Der erste Driven-Adapter des ACK-Wegs trägt seine Grenzen
  (Null-Position-Verbot, Fehlerklassen-Wrapping) ohne eine eigene
  Test-Datei; die Failure-Übersetzung, die der slice-004-F-3-Klasse hier
  ordentlich getragen wird, ist unbelegt.
- `verifizierbar`: ja — `go test ./internal/adapters/driven/postgresack/...`
- `klasse`: Neuer Port-Kontrakt ohne Negativtest

### F-10 — `decode.Change.XID` ist ein totes Feld

- `kategorie`: LOW
- `quelle`: Maintainability (totes Feld als falsche Vertrags-Lesart)
- `pfad`: `internal/adapters/driving/replication/decode/decode.go:93` (Feld
  `XID uint32` in `Change`) — `Decode` setzt es nie; `mapper.Assembler.Consume`
  liest es nie (die XID-Zuordnung läuft über den BEGIN-Stand des Assemblers);
  nur `mapper_test.go:216` setzt es
- `befund`: Das Feld liest sich als pro-Change-Transaktions-Kennung; die
  tatsächliche Zuordnung ([`ADR-0011`](../plan/adr)-Basis) trägt ausschließlich der
  Assembler-`open`-Stand aus BEGIN. Ein Leser, der dem Feld traut, liest
  immer 0.
- `verifizierbar`: ja — `grep -n "XID" decode/decode.go mapper/mapper.go`
- `klasse`: Totes Vertrags-Feld

### F-11 — slice-005 §2: „`make gates` grün" doppelt (Planner hat die Dopplung in slice-006 gestrichen, hier nicht)

- `kategorie`: LOW
- `quelle`: review-slice-004 (a) — Vorlagen-Glitch, Planner-Korrektur
  zugesagt
- `pfad`: `docs/plan/planning/in-progress/slice-005-replication-stream-adapter.md`
  §2 (zwei aufeinanderfolgende Zeilen „`make gates` grün") gegen
  `31e1b41` (streicht die identische Dopplung in slice-006 §2)
- `befund`: Der Planner-Zug, der denselben Glitch in slice-006 korrigierte,
  berührte slice-005 (§1) im selben Commit — die Dopplung in §2 blieben
  stehen. Plan-Defekt, Planner-Sache, keine Plan-Änderung des Implementers.
- `verifizierbar`: ja — `grep -n "make gates" <Plan-Datei>`
- `klasse`: Vorlagen-Glitch unkorrigiert bei berührter Datei

### F-12 — TRUNCATE trägt Fehlerklasse `schema`, ist aber eine korrekt dekodierte Nicht-Unterstützung

- `kategorie`: LOW
- `quelle`: `SPEC-008` (Kategorie `schema`: „nicht sicher interpretierbare
  Schemaänderung/Dekodierfehler") · `LH-FA-CAP-003` Out-of-Scope
- `pfad`: `internal/adapters/driving/replication/mapper/mapper.go:24-28` (`ErrTruncateUnsupported`
  — Klasse `schema`) gegen die `SPEC-008`-Kategorientabelle
- `befund`: Eine TRUNCATE-Nachricht ist fehlerfrei dekodiert; die Grenze ist
  die MVP-Nicht-Unterstützung, kein Dekodierfehler. Die Zuordnung zur Klasse
  `schema` ist am Rand des Kategorie-Vertrags und nur im Code-Kommentar
  begründet (die sichtbare-Fehler-Wirkung stimmt unabhängig von der Klasse —
  es geht um `cdc_errors_total{class}`).
- `verifizierbar`: ja — `SPEC-008`-Tabelle gegen die Sentinel-Texte
- `klasse`: Fehlerklassen-Zuordnung am Kategorie-Rand

### F-13 — XID-Wraparound-Abstand nicht in der Transaktions-Kennung, Grenze unbenannt

- `kategorie`: INFO
- `quelle`: [`ADR-0011`](../plan/adr/README.md) (Deduplizierung wiederholter WAL-Daten) ·
  Implementer-Risiko (d) · Maintainability
- `pfad`: `internal/adapters/driving/replication/mapper/mapper.go:94`
  (`TransactionID` = dezimaler XID, ohne Epochen-/Wraparound-Abstand) ·
  Plan §6/§7 ohne Eintrag; kein Code-Kommentar
- `befund`: Nach 2^32 Quelltransaktionen wiederholen sich XIDs — eine neue
  Transaktion trägt die Kennung einer alten, und die Idempotenz-Form des
  Stores würde still deduplizieren, was keine Wiederholung ist. Die Grenze
  ist im MVP-Fenster real fern, aber sie ist weder im Code noch im Plan
  benannt — sie lebt nur im Handoff. Zuständigkeit: Plan-§6-Ausgang oder
  Kontrakt-Notiz (Planner). (Hinweis ohne erwartete Aktion am Diff —
  INFO-Kanal.)
- `verifizierbar`: ja — Kennungs-Form gegen [`ADR-0011`](../plan/adr/README.md)-Kontrakt;
  `grep -n Wraparound` ohne Treffer
- `klasse`: Unbenannte Grenze am Spec-Rand (Implementer im Handoff benannt)

### F-14 — `CopyData` mit leerem Payload: Index-Zugriff statt Fehler

- `kategorie`: INFO
- `quelle`: Maintainability (defensive Grenze am Treiber-Rand)
- `pfad`: `internal/adapters/driving/replication/receive/receive.go:277-278`
  (`switch message.Data[0]` ohne Längenprüfung)
- `befund`: Ein leeres `CopyData` (stiller Protokollverstoß der Quelle) endet
  in einer Panik statt im Fehlerklassen-Pfad. Serverseitig nicht konstruktiv;
  die Notiz gehört zur Fehler-Übersetzungs-Disziplin des Adapters.
- `verifizierbar`: ja — Test mit leerem CopyData (Panik)
- `klasse`: Panik statt Fehlerklasse am Rand

### F-15 — Keepalive-Antwort vor erster Bestätigung meldet Position 0

- `kategorie`: INFO
- `quelle`: Implementer-Risiko (a)-Umfeld · `LH-QA-REL-001.a`
- `pfad`: `receive.go:81` (`lastAcked` Startstand 0), `receive.go:297-300`
- `befund`: Der Startstand ist vom Zustand „keine bestätigte Position" nicht
  unterscheidbar; korrekt nach der Zusage (nie Empfangsstand), aber der
  fehlende Keepalive-Test (F-3) verdeckt genau diese Randlage. Beleg-Kandidat
  für den ausstehenden Test, kein Befund am Verhalten.
- `verifizierbar`: ja — als Assertion im künftigen Keepalive-Test
- `klasse`: Zustands-Verschmelzung ohne Test-Blick

### F-16 — `PH-DEP-002`/`PH-TST-001` in slice-006 (undeklarierte Präfixe) von `31e1b41` berührt, unberührt gelassen

- `kategorie`: INFO
- `quelle`: Skill-HIGH-Klasse „Traceability-/ID-Schema-Verstoß" (Erstauftreten
  review-slice-002 F-2, `PH-*`/`TST-*`) — Vorbestand, nicht vom
  Implementer-Diff
- `pfad`: `docs/plan/planning/open/slice-006-integrationstest-umgebung.md:12`
  („Berührte Spec-Stellen: `PH-DEP-002`, `PH-TST-001`") — der Planner-Commit
  `31e1b41` editierte dieselbe Datei (§1/§3), ohne die Kennungen zu
  korrigieren
- `befund`: Die undeklarierten Präfixe stehen seit dem Erstauftreten im
  Zähler; die Berührung im Range korrigierte sie nicht. Planner-Sache —
  kein Befund gegen den Implementer-Diff.
- `verifizierbar`: ja — MR-000-Deklaration gegen die Plan-Kopf-Zeile
- `klasse`: Undeklariertes ID-Präfix (Vorbestand, berührt und unkorrigiert)

---

## Design-Entscheidungen des Implementers — Bewertung (Prüf-Fragen)

- **(a) pglogrepl v0.0.0-20260824 als Replication-Treiber:** *Plan-konform,
  kein Befund.* Die grep-Begründung (pgconn führt keinen CopyBoth-Zug, kein
  START_REPLICATION-Handshake, keine Standby-Status-Updates) trägt die
  Wahl-Aussage am Commit (`d1a84eb`); die konkrete Treiber-Bibliothek ist
  Infrastrukturdetail am Adapter ([`ADR-0032`](../plan/adr/README.md)-Konsequenz), und `pgoutput`
  bleibt das Output-Plugin ([`ADR-0008`](../plan/adr/README.md) — `outputPlugin`-Konstante,
  `PluginArgs`). Kein ADR-Verstoß, kein Plan-Änderungs-Bedarf für die
  Dependency selbst (Präzedenz slice-004: auch dort keine go.mod-Zeile in §3).
- **(b) Keepalive meldet nur die bestätigte Position:** inhaltlich
  [`LH-QA-REL-001.a`](../../spec/pflichtenheft.md)-deckungsgleich mit dem §2-DoD-Punkt (ACK nur nach
  Persistenz) — *Plan-Ergänzung, nicht Plan-Änderung*; kein §1-Ausschluss
  berührt. Aber: F-2 (Nachzug fällig) und F-3 (die Zusagen trägt kein Test —
  Mutation belegt).
- **(c) `BindCapture` / `Conn()` (Option C):** [`ADR-0007`](../plan/adr/README.md)-konform —
  der Stream baut die technische Verbindung, der ACK-Adapter ist die
  Driven-Rolle an derselben Verbindung; die Rollen bleiben durch getrennte
  Pakete (driving/driven), getrennte Sentinels (`receive.ErrReplication` vs.
  `outbound.ErrReplication`) und die Composition-Root-Verdrahtung benannt —
  genau die „Benennung (Adapter-Namen)", die die ADR-Konsequenz fordert. Die
  Verlagerung des Verdrahtungs-Tests in `internal/bootstrap` entspricht der
  Maschinenform (`.a-check.yml` `composition_root: bootstrap, cmd` — a-check
  0 Befunde am Range-Head). Plan-Nachzug siehe F-2.
- **(d) Verlagerung nach `internal/bootstrap` (nach a-check-Befunden):**
  sachlich richtig — der Adapter-Test importierte Use Case + fremde Adapter
  und erzeugte lateral-adapter-/wrong-direction-Befunde; die Rolle der
  Verdrahtung liegt in der Composition-Root ([`ADR-0026`](../plan/adr/README.md)). Kein
  Architekt-Verdikt nötig. Plan-Nachzug siehe F-2.

## ADR-Deckung (Stream, ACK, Dekodierung, Mapper)

| ADR | Aussage | Träger im Diff |
|---|---|---|
| [`ADR-0006`](../plan/adr/README.md) | Stream ist Driving Adapter; „dekodieren ja, entscheiden nein" | getragen — `receive/` ruft `CaptureInboundPort`; Persistenz/ACK-Entscheidung liegt beim Service; Keepalive meldet nur das Capture-Ergebnis (`lastAcked`), nie den Empfangsstand |
| [`ADR-0007`](../plan/adr/README.md) Option C | ACK als Outbound-Port, dieselbe technische Verbindung, getrennte Rollen | getragen — `postgresack.New(stream.Conn())`, getrennte Pakete und Sentinels, Verdrahtung im Composition-Root (bootstrap-Test) |
| [`ADR-0008`](../plan/adr/README.md) | `pgoutput` als Standard-Output-Plugin | getragen — `outputPlugin`-Konstante, `proto_version 1`; [`SPEC-010`](../../spec/pflichtenheft.md) als Rang-Zeiger |
| [`ADR-0011`](../plan/adr/README.md) | Persist-before-ACK; idempotente Persistenz | getragen — Ordnung trägt der Capture Service (unverändert, `service.go:59-68`); der Adapter meldet nur die bestätigte Position; Restart liest `confirmed_flush_lsn` (F-4: ohne Test) |
| [`ADR-0012`](../plan/adr/README.md) | At-Least-Once; Wiederholung ist Betrieb | getragen — `ensureSlot` startet am `confirmed_flush_lsn`-Stand des bestehenden Slots; Beleg fehlt (F-4) |
| [`ADR-0016`](../plan/adr/README.md) | JSON-Row-Images als Text in Relation-Reihenfolge | getragen — `rowImage` ohne Typ-Interpretation, Relation-Spalten-Reihenfolge; exakte JSON-Tests |
| [`ADR-0023`](../plan/adr/README.md)/`SPEC-008` | Fehlerklassen je Port-Kontrakt | getragen — Sentinels je Kontrakt (`decode.ErrSchema`, `receive.ErrReplication`/`ErrConfiguration`, `mapper`-Sentinels, `outbound.ErrReplication`); die slice-004-F-3-Klasse (rohe Treiberfehler) tritt hier **nicht** wieder auf — beide neuen Adapter übersetzen; Grenze: F-12 |
| [`ADR-0027`](../plan/adr/README.md)/[`ADR-0029`](../plan/adr/README.md) | Orchestrierung am Service; offene Transaktionen unkonsumierbar | getragen — `TestStreamOpenTransactionNotConsumable` am realen Pfad (`LH-FA-CAP-006.a`) |
| [`ADR-0026`](../plan/adr/README.md) | Verdrahtung nur im Composition-Root | getragen — Verlagerung konform; a-check 0 Befunde |
| [`ADR-0030`](../plan/adr/README.md) | Testpyramide | getragen — Unit gegen `pgoutput`-Binärcodes (Handform, kein Netz), Real-Tests gegen gepinnte Instanz |
| [`ADR-0043`](../plan/adr/README.md) §5 | Testloader (DDL + `ApplySchema`) bleibt, bis d-migrate ihn ersetzt | getragen — der Verdrahtungs-Test nutzt `ApplySchema` in der ausdrücklich benannten Grenze; Erstversatz slice-006; die §1-Ausschluss-Zeile trägt Begründung und Folge-Slice-Verweis — **konform** |

## Test-Qualität — reale Treiber-Läufe und Mutations-Proben

`make test` am Range-Head: treiberfrei grün (11 Pakete mit Tests; `postgresack`
ohne Testdateien — F-9). `make test-replication` am Range-Head: **grün**
(`receive` 3.0 s — vier Stream-Tests gegen reale PostgreSQL mit
Publication/Slot; `bootstrap` 0.5 s — Verdrahtungs-Test am echten Capture
Service). Eigene Mutations-Proben (Worktrees außerhalb des Arbeitsbaums,
gepinnte Container):

| Mutation | Erwartung | Ergebnis |
|---|---|---|
| Commit-LSN falsch gemappt (`CommitLSN` → `TransactionEndLSN`) | Decode-Tests rot | **rot** (Position-Zusage trägt Tests) |
| Sequenz-Zähler entfernt (`a.open.sequence++` gestrichen) | Mapper-Tests rot | **rot** (Sequenz-Ordnung trägt Tests) |
| Keepalive meldet 0 statt `lastAcked` | Replication-Suite rot | **grün — Zusage ohne Test-Träger (F-3)** |

Die kritischen Zusagen auf der Dekodier-/Mapper-Seite (Commit-Position,
Sequenz-Reihenfolge, Row-Images, TRUNCATE, Aktivierungs-Filter,
Sichtbarkeit `LH-FA-CAP-006.a` am realen Pfad) tragen ihre Tests. Die
Keepalive- und Restart-Zusagen tragen keine (F-3/F-4). Die drei Gates am
Range-Head: `baseline-verify` OK (54 Dateien) · `d-check` 92 Dateien/0
Befunde · `a-check` 0 Befunde (Verlagerung deckt sich mit
`composition_root`).

## Implementer-Risiken — Bewertung

- **(a) Keepalive-Position-Nicht-Abdeckung:** **F-3** — bestätigt und
  mutations-belegt; Ausgang im Plan §6 fehlt.
- **(b) Restart auf bestehendem Slot ohne Test:** **F-4** — bestätigt;
  zusätzlich trägt auch slice-006 das Risiko nicht (kein Ausgangs-Träger).
- **(c) Image-Digest manifest-instabil:** **F-7** — die Binary-These ist
  probe-bestätigt (bit-identisch), die Varianz-These auf dem Review-Daemon
  nicht reproduziert; die `96c47af`-Formel ist in jedem Fall widerlegt —
  Architect/Planner.
- **(d) XID-Wraparound als MVP-Detail:** **F-13** — inhaltlich richtig
  eingestuft, aber die Grenze steht nur im Handoff; Plan-§6-Ausgang oder
  Kontrakt-Kandidat für slice-006.

## Zweitschreiber-Frage — `3125c0a`/`31e1b41` (Planner)

Bewertung: **zulässige Planning-Arbeit, kein Modul-5-Disziplin-Befund gegen
den Implementer-Lauf** — mit drei benannten Residuen.

- **Kein Produktcode, kein Slice-Produktions-Commit:** beide Commits berühren
  nur `spec/pflichtenheft.md` und Planning-Dateien; der Implementer-Baum
  (`internal/`, `tools/`, Makefile) bleibt unberührt.
- **Zeitpunkt:** beide liegen *vor* den Implementer-Commits (19:47/19:49
  gegen 20:57 ff.); die Berührung des slice-005-Plans (§1 Out-of-Scope
  [`ADR-0043`](../plan/adr)) ist eine Grenz-**Verengung** vor der Umsetzung, die der
  gelieferte Diff einhält (kein CDC-Schema-Rollout im Slice; der Verdrahtungs-
  Test nutzt den von [`ADR-0043`](../plan/adr/README.md) §5 ausdrücklich geduldeten
  Testloader). Die Übergabe ist als Commit beobachtbar — Artefakt vorhanden.
- **Kein WIP-Limit-Befund:** das WIP-Limit misst Slices in `in-progress/` je
  Rolleninhaber; Planning-Commits beanspruchen keinen Slice, und das
  `Verantwortlich:`-Feld bleibt beim Implementer.
- **Residuen:** (1) `31e1b41` mischt Stratum-Änderung (`pflichtenheft`) und
  Planning-Dateien in einem Commit — zwei Vorgänge, ein Beleg. (2) Beide
  Planner-Commits tragen `SPEC-*`-Kennungen in der Message — die
  F-5-Klasse (Struktur-ID in Commit-Message) zählt damit über den Vorgang
  hinaus auch auf der Planner-Seite. (3) Die Plan-Änderung während
  `in-progress/` durch einen zweiten Schreiber ist formal beobachtbar und
  materiell konform, sollte aber künftig als eigener Plan-Nachzugs-Commit
  *vor* dem nächsten Umsetzungs-Commit sichtbar bleiben (war hier der Fall).

## Negativbefunde

- geprüft, ohne Befund: **HIGH-Klassen über den Diff (abgesehen von F-1)** —
  kein ADR-Verstoß auf Layer-Ebene (a-check 0 Befunde; [`ADR-0006`](../plan/adr/README.md)
  „dekodieren ja, entscheiden nein" getragen: kein ACK- und kein
  Persistenz-Entscheid im Adapter, `a-check`-Edges erfüllt), kein
  Sicherheits-Anti-Pattern (Slot-/Publication-Namen über das
  Bezeichner-Alphabet begrenzt, bevor sie als Literale interpoliert werden —
  Katalog-Query der Tests parametrisiert; kein Credential-Pfad im Log), kein
  Korrektheitsfehler im kritischen Pfad (Persist-before-ACK am realen Treiber
  belegt; die Stream-Seite bestätigt nichts selbst), keine Gate-Suppression,
  keine Norm nur im Template-Kommentar, kein Chronik-tragendes Zustandsfeld
  im Diff, kein Docker-only-Verstoß (beide Images per Digest gepinnt,
  Modul-Cache im Volume, Daten im Container, Netz-Aufräumpfad inklusive),
  kein Zwei-Quellen-Drift im Diff
- geprüft, ohne Befund: **Stand-alone-Build am Range-Head** — `gofmt -l` leer,
  `go vet ./...` ohne Befund, `go test ./...` treiberfrei grün,
  `make test-replication` grün, drei Gates grün
- geprüft, ohne Befund: **Traceability-Grundpflege der übrigen fünf
  Implementer-Commits** — jeder trägt mindestens eine `LH-*`-/`ADR-*`-Kennung;
  alle genannten IDs existieren ([`LH-FA-CAP-001`](../../spec/lastenheft.md)…003,
  [`LH-FA-CAP-004`](../../spec/lastenheft.md), [`LH-FA-CAP-006.a`](../../spec/pflichtenheft.md), [`LH-FA-CFG-001.a`](../../spec/pflichtenheft.md),
  [`LH-QA-REL-001.a`](../../spec/pflichtenheft.md), [`LH-QA-POR-003`](../../spec/lastenheft.md), [`ADR-0006`](../plan/adr)/0007/0008/0023/0026/0027/0030/0032/0042);
  die `SPEC-*`-Nennungen in Messagen sind F-5 (Klassen-Zählung), kein
  undeklariertes Präfix
- geprüft, ohne Befund: **superseded-ADR-Referenzen** — die neuen Code- und
  Plan-Referenzen nennen [`ADR-0006`](../plan/adr)/0007/0008/0011/0012/0016/0023/0026/0030/0032/0042
  (Accepted) und [`ADR-0021`](../plan/adr) (Proposed, als Proposed benannt); keine
  Referenz auf [`ADR-0038`](../plan/adr)/0039 als tragenden Anker (die
  slice-004-F-4-Klasse tritt hier nicht wieder auf)
- geprüft, ohne Befund: **Kommentar-Klassen im neuen Code (§3.7)** —
  `decode.go`, `mapper.go`, `receive.go`, `ack.go`, `replicationack.go`,
  drei Testdateien, `run-replication-tests.sh`, Makefile-Block,
  harness-README-Zeile — Indikativ, Klassen
  Zusage/Kopplung/Abgrenzung/Rang-Zeiger/Grenze; die `4ca6658`-Korrektur hat
  den Lifecycle-Bezug aus den Zustandskommentaren entfernt (das Gegenstück
  zur F-1-Meldung, inhaltlich richtig)
- geprüft, ohne Befund: **Testcontainer-Hygiene** — Daten und Publication/
  Slot im Container, je Test frisch aufgesetzt, Slot-Rückbau mit Retry,
  Netz- und Container-Rückbau in jedem Ausgang (der slice-004-F-5-Netz-Defekt
  ist hier nicht wiederholt); gepinnte Digests konsistent mit slice-004
  (Toolchain `cf6fca…`, PostgreSQL `63bdc97…`)
- geprüft, ohne Befund: **`SPEC-008`-Übersetzung am neuen Treiber-Adapter**
  (slice-004-F-3-Klasse) — `postgresack` und `receive` tragen die
  Übersetzung in Sentinels am Port-Kontrakt; keine rohen Treiberfehler am
  Port (Grenze: F-12 Kategorie-Rand)
- geprüft, ohne Befund: **§3-Datei-Menge** — die beiden §3-Zeilen decken
  `receive/decode/mapper` (inkl. Tests, `*.go`-Muster) und `postgresack/` ab;
  die darüber hinaus gelieferten Dateien sind F-2 (Nachzug), kein stiller
  Umfang
- geprüft, ohne Befund: **Beobachtungs-Register** — keine neue Beobachtung im
  Diff angefallen; die Register-Pflichten sind Closure-Sache (§7 steht aus)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 7 |
| LOW | 4 |
| INFO | 4 |

**Finding-Klassen dieses Laufs:** Commit ohne Vertrags-Kennung (1. Auftreten
der direkten Form) · Plan-Erweiterung ohne Plan-Nachzug (5. Auftreten —
Sequenz läuft) · Zusage am kritischen Pfad ohne Test-Träger (Keepalive,
Mutation belegt) · Kritischer Pfad ohne Beleg und ohne Ausgangs-Träger
(Restart) · Struktur-ID in Commit-Message (4. Auftreten — Sequenz läuft
nicht) · Datei-Abschluss (3. Auftreten — MEDIUM-Stufe) · Beleg-Semantik-Formel
widerlegt (`96c47af` vs. Digest-Verhalten) · stiller Vertragsverstoß ohne
Negativtest (Doppel-BEGIN) · neuer Port-Kontrakt ohne Negativtest (postgresack)
· totes Vertrags-Feld (`Change.XID`) · Vorlagen-Glitch unkorrigiert (§2-Dopplung)
· Fehlerklasse am Kategorie-Rand (TRUNCATE → `schema`) · unbenannte Grenze
(XID-Wraparound) · Panik statt Fehlerklasse (leeres CopyData) · Keepalive-0-Stand
· undeklariertes ID-Präfix (Vorbestand slice-006, Planner)

## Verdikt

**Merge-blockierend:** nein — der Diff ist inhaltlich schlüssig, kompiliert,
vet- und gofmt-sauber stand-alone, die Dekodier-/Mapper-Zusagen sind
mutations-geprüft, die Translation und die Persist-before-ACK-Verdrahtung
sind am realen Treiber belegt, die vier gemeldeten Entscheidungen sind
ADR- und Plan-seitig tragfähig (drei davon brauchen Nachzug, keine
Korrektur), und die drei Gates laufen grün.

**Blockierend für Closure:** ja, in vier Punkten, bevor der Slice nach
`done/` geht:

1. **F-2 als Sequenz-Fortsetzung (Modul 8, 5. Auftreten):** Plan-Nachzug für
   (b) Keepalive-Regel, (c) `BindCapture`/`Conn`-Trennung und (d)
   Verlagerung/Port-Datei/`test-replication`-Verkabelung über den **Planner**
   (§3-Zeilen bzw. §7-Zeile „Was ging anders als geplant"); Übergabe-Artefakt:
   dieser Report-Abschnitt. Die Sequenz läuft seit review-slice-003 — der
   Plan-Nachzug ist nicht mehr Implementer-Vor-Closure-Arbeit.
2. **F-3 + F-4:** die Keepalive-Zusage und der Restart-Zweig brauchen einen
   Beleg — entweder Test (Keepalive-Position gegen `confirmed_flush_lsn`;
   Restart am bestehenden Slot mit deduplizierter Wiederholung) oder
   terminierter Ausgang im Plan §6 **mit Kennung eines Slices, der den Punkt
   annimmt** (slice-006 nimmt beides derzeit nicht an — vor der Closure
   klären, nicht still lassen).
3. **F-5 als Konflikt-Sequenz (4. Auftreten der Struktur-ID-Klasse):** die in
   review-slice-003 ausgelöste Sequenz ist weiterhin nicht gelaufen —
   Übergabe-Artefakt Konventions-Nachzug über den **Architect**; zählt
   zusammen mit F-1 (Kennungspflicht) in den Steering-Loop.
4. **F-6 + F-7:** Datei-Abschluss in `go.mod`/`go.sum` (MEDIUM-Stufe erreicht)
   und die Beleg-Semantik-Frage über den **Architect/Planner** adjudizieren —
   die `96c47af`-Formel beschreibt den Mechanismus nicht; der Beleg selbst
   reproduziert am HEAD.

F-1 (HIGH) wirkt nur vorwärts: künftige Commits tragen je mindestens eine
`LH-*`/`ADR-*`-Kennung; die Klasse zählt als erstes Auftreten der direkten
Form. F-8 (Doppel-BEGIN) und F-9 (postgresack-Negativtests) sind Vor-Closure-
Nacharbeit ohne Blockier-Charakter. F-11/F-13–F-16 gehen an Planner/Verifier
bzw. in die Closure §7. DoD- und Spec-Konformität prüft der Verifier separat
(Modul 11) — insbesondere `make gates` grün und der Keepalive-/Restart-Beleg.