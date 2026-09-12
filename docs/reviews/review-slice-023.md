# Review-Report: slice-023 — 2026-09-12

**Review-Art:** Code — geprüft gegen Plan (`slice-023`) + `ADR-0047`
(Rollen-spezifische DSN-Verdrahtung, Accepted) + Konventionen (`AGENTS.md`,
`ADR-0026`) + Vorgänger-Finding-Klassen/Beobachtungen
(`BEO-PGC/rollen-verdrahtung`, `BEO-PGC/test-isolation-geteilter-zustand`).

**Gegenstand:** Commit `a32a2c9` — `feat(bootstrap): rollen-spezifische
DSN-Verdrahtung nach ADR-0047 (LH-QA-SEC-001)`

**Skill:** `.harness/skills/reviewer.md` @ `e9159ff` (HEAD zum Review-Zeitpunkt)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-12

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-023-rollen-spezifische-dsn-verdrahtung.md`
  §1–§8 (inkl. Plan-Nachzug)
- `docs/plan/adr/0047-rollenspezifische-dsn-verdrahtung.md` (Accepted)
- `docs/plan/adr/0048-heartbeat-grant-korrektur-select-ergaenzung.md`
  (Accepted, Supersedes `ADR-0047` teilweise — Commit `e486e1b`, unmittelbar
  nach dem Review-Gegenstand gelandet; als Kontext gelesen, nicht selbst
  Review-Gegenstand dieses Reports)
- `LH-QA-SEC-001`/`002`/`003` (`spec/lastenheft.md`)
- `AGENTS.md` §3 (Hard Rules, insbesondere §3.5, §3.7)
- `harness/conventions.md` (MR-000/MR-001)
- `docs/plan/planning/observations/BEO-PGC/rollen-verdrahtung`,
  `docs/plan/planning/observations/BEO-PGC/test-isolation-geteilter-zustand`

---

## Findings

### F-1 — Grant-Korrektur weicht vom in `ADR-0047` §Konsequenzen genannten Text ab, ohne dass der Commit selbst eine Folge-ADR mitführt

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.5 (ADR-Immutabilität) / Baseline-Regelwerk
  `v6.5.0` · `regelwerk/modul-08-agentenrollen.md` §Konflikt-Pfad als
  Rollen-Sequenz, Verdikt „Lockerung legitim, aber undokumentiert"
- `pfad`: `tools/schema/nacharbeit-roles.sql` (neuer Grant `GRANT SELECT,
  INSERT, UPDATE ON cdc.process_heartbeat TO cdc_admin`) vs.
  `docs/plan/adr/0047-rollenspezifische-dsn-verdrahtung.md` §Konsequenzen
  (nennt `GRANT INSERT, UPDATE ...`)
- `befund`: Der Commit korrigiert real und nachvollziehbar einen SQL-Text,
  den `ADR-0047` als Folgepflicht benennt (siehe F-1-Verifikation unten),
  ohne die ADR-Datei selbst zu verändern (Hard Rule 3.5 korrekt
  eingehalten) — trägt aber selbst keine begleitende Folge-ADR, sondern nur
  einen Verweis „Kandidat für einen Steering-Loop-Eintrag bei der
  Planner-Closure (§7)" (Slice-Plan §3 Plan-Nachzug). Für die Dauer
  zwischen diesem Commit und der ADR-Korrektur widerspricht der
  tatsächliche Datenbank-Zustand dem, was das noch geltende `ADR-0047`
  wörtlich als Konsequenz vorschreibt. Der unmittelbar folgende Commit
  `e486e1b` (`ADR-0048`, Supersedes `ADR-0047` teilweise) schließt diese
  Lücke bereits als Sofort-Folgecommit — exakt das im Baseline-Regelwerk
  vorgesehene Muster —, ist aber nicht Teil des hier geprüften Commits.
  Eingestuft als LOW statt HIGH, weil (a) die Abweichung im Commit selbst
  vollständig transparent dokumentiert ist (Commit-Message, Plan-Nachzug),
  (b) `ADR-0047`s Kern-Entscheidung (Drei-DSN-Vertrag, Rollen-Zuordnung)
  unberührt bleibt, und (c) die Lücke real bereits im nächsten Commit
  geschlossen ist.
- `verifizierbar`: ja — `git log --oneline` zeigt `e486e1b` unmittelbar
  nach `a32a2c9`; `docs/plan/adr/README.md`/`0048-...md` bestätigen die
  Supersedes-Beziehung.
- `klasse`: „ADR-Konsequenz real widerlegt ohne Folge-ADR im selben Commit"

## Negativbefunde

- geprüft, ohne Befund: `internal/bootstrap/wiring.go` — alle acht Aufrufer
  aus `ADR-0047`s Zuordnungstabelle (Store-Adapter, Replication-Stream,
  ACK-Adapter, Tabellen-Aktivierung, Heartbeat-Adapter, `RegisterConsumer`,
  `AcknowledgeConsumer`, `Healthcheck`) binden exakt an die dort
  vorgesehene Rolle/DSN — Zeile-für-Zeile gegen die Tabelle geprüft
  (`grep` auf `cfg.CaptureDSN`/`cfg.AdminDSN`/`cfg.ReaderDSN`), keine
  verbliebene Referenz auf ein gemeinsames DSN-Feld.
- geprüft, ohne Befund: `cmd/pg-change-feed/main.go` — `--healthcheck`
  zieht korrekt auf `cfg.ReaderDSN` nach; `register-consumer`/
  `acknowledge-consumer` unverändert (nutzen weiterhin `cfg` als Ganzes,
  intern über `cfg.AdminDSN` in `wiring.go` aufgelöst).
- geprüft, ohne Befund: `internal/bootstrap/roles_wiring_test.go` — real
  gegen PostgreSQL-Testcontainer reproduziert (dreimal `make test-store`
  in dieser Review-Sitzung selbst ausgeführt, alle drei grün): eine
  `cdc_reader`-Login-Identität scheitert real an `INSERT` auf
  `cdc.change`, eine `cdc_capture`-Login-Identität real an `INSERT` auf
  `cdc.consumer`, und der Heartbeat-Schreibpfad gelingt nachweislich erst
  nach dem ergänzenden Grant (Vorher/Nachher, Grant real entzogen und
  wiederhergestellt). Test-Fixtures (Login-Rollen, Consumer-/Change-Zeilen)
  tragen eindeutige, testspezifische Kennungen und werden über
  `t.Cleanup` aufgeräumt (Login-Rollen) bzw. mit dem Grant wiederhergestellt
  — keine unskopierten `DELETE`/`DROP`-Operationen, kein neuer Beitrag zu
  `BEO-PGC/test-isolation-geteilter-zustand`; der Test läuft zudem im
  bereits isolierten, vorgezogenen `internal/bootstrap`-Lauf von
  `tools/harness/run-store-tests.sh`.
- geprüft, ohne Befund: **Grant-Korrektur, technische Prüfung** — real
  gegen einen frischen PostgreSQL-16-Container reproduziert
  (`INSERT ... ON CONFLICT (id) DO UPDATE ...` scheitert mit `permission
  denied for table t` bei `GRANT INSERT, UPDATE` allein und gelingt erst
  nach zusätzlichem `GRANT SELECT`) — der im Commit behauptete
  PostgreSQL-Mechanismus ist bestätigt korrekt.
- geprüft, ohne Befund: Hard Rule 3.5 (ADR-Immutabilität) — `git show
  a32a2c9 --stat` enthält keine Änderung an
  `docs/plan/adr/0047-rollenspezifische-dsn-verdrahtung.md`; die Datei
  wurde nicht editiert.
- geprüft, ohne Befund: bestehende `postgresstorage`-Adapter-Tests — sie
  verbinden sich direkt über `CDC_STORE_TEST_DSN`, nicht über
  `bootstrap.Config`, und sind von der Feldumbenennung strukturell
  unberührt; alle Pakete liefen in dieser Review-Sitzung real grün mit.
- geprüft, ohne Befund: `wiring_test.go`/`register_test.go`/
  `acknowledge_test.go`/`welle6_endtoend_test.go` — Feldumbenennung
  (`DSN` → `AdminDSN`/Drei-Feld-Form) korrekt und vollständig nachgezogen,
  keine übersehene Stelle.
- geprüft, ohne Befund: `compose.yaml`, `docs/user/benutzerhandbuch.md`,
  `harness/README.md` — kein verbliebener `CDC_SOURCE_DSN`-Bezug außer der
  historischen Erwähnung im Changelog-Eintrag und im erklärenden
  Quellcode-Kommentar (beide sachlich korrekt: der Name entfällt).
  `compose.yaml` bindet alle drei neuen Variablen bewusst an denselben
  Superuser-DSN der Testumgebung — dokumentiert begründet (Testrollen
  bleiben `NOLOGIN`, reale Rollentrennung wird in `internal/bootstrap`
  belegt, nicht im Compose-Integrationslauf); siehe INFO-1.
- geprüft, ohne Befund: Hard Rule 3.7 (Kommentare) — die neuen/geänderten
  Kommentare in `wiring.go`, `main.go`, `nacharbeit-roles.sql` nennen den
  geltenden Zustand indikativ (Zusage-/Kopplungs-/Grenz-Klasse), keine
  Beschreibung verworfener Alternativen oder abgebrochener Sätze; keine
  Slice-Chronik-Referenzen im SQL/Produktcode.
- geprüft, ohne Befund: Slice-Plan §1 Abgrenzung — kein Neu-Zuschnitt der
  drei Rollen (nur die eine Grant-Lücke geschlossen, wie in `ADR-0047`
  selbst begründet), keine Auth/Autorisierung des externen Zugriffswegs
  berührt, keine Verdrahtung künftiger Zugriffswege.
- geprüft, ohne Befund: DoD-Checkbox-Ehrlichkeit — alle als `[x]`
  markierten Punkte (SEC-001…003 real getestet, Konfigurationsvertrag,
  `make gates`, Doku-Update) sind durch reale, in dieser Sitzung
  reproduzierte Läufe gedeckt; die als `[ ]` offen belassenen Punkte
  (Review, Closure-Notiz, Reconciliation-/Beobachtungs-Register, Risiken,
  Paarungen) sind konsistent mit dem Lifecycle-Zustand `in-progress`.
- geprüft, ohne Befund: Traceability — Commit-Betreff nennt `ADR-0047` und
  `LH-QA-SEC-001`, kein `SPEC-*`/`ARC-*` im Betreff; `make
  commit-traceability` lief in dieser Sitzung grün mit.
- geprüft, ohne Befund: Mutationstest (eigenständig durchgeführt) —
  `RegisterConsumer` testweise auf `cfg.CaptureDSN` statt `cfg.AdminDSN`
  umgestellt: drei Tests (`TestAcknowledgeConsumerEndToEnd`,
  `TestRegisterConsumerEndToEnd`, `TestWelle6ConsumerFullCycleEndToEnd`)
  wurden real rot (`make test-store`); Rücksetzung verifiziert (`git diff`
  danach leer), anschließender Lauf wieder grün.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** ADR-Konsequenz real widerlegt ohne
Folge-ADR im selben Commit

### INFO-1 — Rollentrennung im Compose-Integrationslauf nicht real exerziert

- `kategorie`: INFO
- `quelle`: Maintainability (Test-Abdeckungshinweis, kein Mangel — Weiterleitung
  an Verifier/Validator für die Frage, ob `test-integration` künftig eigene
  Login-Identitäten je Rolle braucht)
- `pfad`: `compose.yaml` (`CDC_CAPTURE_DSN`/`CDC_ADMIN_DSN`/`CDC_READER_DSN`
  zeigen alle auf denselben Superuser)
- `befund`: Die reale Rechte-Trennung wird ausschließlich durch
  `internal/bootstrap/roles_wiring_test.go` (`make test-store`) belegt,
  nicht durch `make test-integration` — dort läuft der Feed-Container
  faktisch mit Superuser-Rechten auf allen drei Kanälen. Das ist explizit
  begründet (`ADR-0047` Kontext-Befund 1: die ADR entscheidet nicht, wie
  der Betreiber Login-Identitäten anlegt) und im Compose-Kommentar
  dokumentiert, kein stiller Mangel.
- `verifizierbar`: ja — `docker inspect`/`compose.yaml` zeigen identische
  DSN-Werte für alle drei Variablen.
- `klasse`: „Least-Privilege nur in Unit-nahen Tests, nicht im
  End-to-End-Lauf belegt"

## Verdikt

**Merge-blockierend:** nein — kein HIGH, kein MEDIUM. Die einzige LOW-
Einstufung (F-1) beschreibt eine bereits durch den unmittelbar folgenden
Commit (`e486e1b`, `ADR-0048`) geschlossene Lücke; kein Nacharbeits-Bedarf
am geprüften Commit selbst.

**Übergabe:** Findings gehen an den Implementer. Die Finding-Klasse aus F-1
geht zusätzlich in die Slice-Closure §7 und von dort in den Zähler
(`docs/plan/planning/observations/`). Dieser Report ist ein Lauf-Beleg und
ersetzt keine Verifikation — DoD-/Spec-Konformität (inkl. der noch offenen
DoD-Punkte: Closure-Notiz, Beobachtungs-/Reconciliation-Register, Risiken-
Ausgänge, drei Paarungen) prüft der Verifier separat.
