# Review-Report: slice-028 — 2026-09-12

**Review-Art:** Code — geprüft gegen Plan (`slice-028`, §1/§2/§3) und
`AGENTS.md` §3 Hard Rules (Modul 10 §Drei Review-Arten).

**Gegenstand:** Commits `632ddca`, `ba508ed`, `9a84407`, `f8ce37b`
(letzter Slice von `welle-8`).

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-12

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-028-integrationstest-nachzug-rollen-mvp.md` (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan)
- `docs/plan/adr/0047-rollenspezifische-dsn-verdrahtung.md` (Accepted)
- `docs/plan/planning/done/slice-023-rollen-spezifische-dsn-verdrahtung.md` und dessen Closure-Notiz §7
- `docs/reviews/review-slice-023.md` (Vorgänger-Review, INFO-1)
- `docs/plan/planning/observations/BEO-PGC/rollen-test-abdeckungsluecken/` (observation.md, state.md, evidence/slice-023.md)
- `spec/lastenheft.md` (`LH-QA-SEC-001`…`003`)
- `compose.yaml`, `internal/bootstrap/roles_wiring_test.go`
- `AGENTS.md` §3 Hard Rules

---

## Findings

### F-1 — Closure-Beleg für Punkt (2) breiter formuliert als der tatsächliche Testumfang

- `kategorie`: MEDIUM
- `quelle`: Maintainability — Genauigkeit des Beobachtungs-Register-Belegs
  (`BEO-PGC/rollen-test-abdeckungsluecken`), DoD-Checkbox-Ehrlichkeit
  (dieselbe Prüfkategorie, die `review-slice-023.md` bereits anwandte)
- `pfad`: `tools/harness/run-integration-tests.sh:67-142` (neuer Abschnitt
  „Rollen-DSN-Verifikation gegen den Compose-Stack", Commit `ba508ed`);
  `docs/plan/planning/in-progress/slice-028-integrationstest-nachzug-rollen-mvp.md:88-95`
  (DoD §2, erstes Häkchen, Commit `f8ce37b`)
- `befund`: Die DoD-Zeile behauptet, ein „Replication-Stream-/ACK-Aufruf
  mit vertauschter Rolle" werde „am laufenden, containerisierten System
  zurückgewiesen" und schließe `BEO-PGC/rollen-test-abdeckungsluecken`
  Punkt (2) — Punkt (2) benennt wörtlich die *Adapter*
  (`receive.NewStream`/`postgresack.New`, `cdc_capture`-gebunden). Der
  neue Testabschnitt ruft diese Adapter jedoch nie auf: Er legt drei
  ephemere, eigens angelegte Login-Identitäten an und prüft per rohem
  `psql`-Aufruf (`replication=database`, `IDENTIFY_SYSTEM;`) ausschließlich
  PostgreSQL-serverseitig, ob das `REPLICATION`-Attribut fehlt/vorhanden
  ist — entkoppelt vom laufenden Feed-Container (`cdc-test-feed`), der für
  alle drei DSNs unverändert mit dem Superuser `postgres` verbunden bleibt
  (siehe F-2). Das ist derselbe Testmuster-Typ, den
  `internal/bootstrap/roles_wiring_test.go` bereits für andere Aufrufer
  fährt (dort teils über die realen Adapter-Konstruktoren, teils —
  dokumentiert benannt, Zeilen 262-266 — als offene Lücke gerade für
  Replication-Stream/ACK), aber nicht dasselbe wie „der Adapter ist gegen
  Rollen-Vertauschung testgesichert". Die Aussage im DoD/Register-Beleg
  ist damit enger zu fassen, als sie derzeit formuliert ist: belegt ist
  „PostgreSQL verweigert einer Login-Identität ohne `REPLICATION`-Attribut
  eine Replication-Verbindung auf dem Compose-Stack" — nicht „der
  Replication-Stream-/ACK-Adapter verhält sich bei Rollen-Vertauschung
  korrekt".
- `verifizierbar`: ja — Code-Lektüre von `tools/harness/run-integration-tests.sh`
  zeigt, dass der neue Abschnitt keinen Aufruf von `receive.NewStream`
  oder `postgresack.New` (bzw. eines Go-Testbinaries, das diese
  Adapter-Konstruktoren nutzt) enthält, sondern ausschließlich
  `docker exec … psql …`.
- `klasse`: „Closure-Beleg breiter formuliert als tatsächlicher Testumfang"

### F-2 — Compose-Referenzumgebung bleibt Superuser-verdrahtet (bestätigt, aber vorbestehend und bereits klassifiziert)

- `kategorie`: INFO
- `quelle`: Maintainability — Nachprüfung einer bereits im Vorgänger-Review
  klassifizierten Beobachtung, kein neuer Befund
- `pfad`: `compose.yaml:53-64` (Env-Block des `pg-change-feed`-Diensts)
- `befund`: Bestätigt: `CDC_CAPTURE_DSN`/`CDC_ADMIN_DSN`/`CDC_READER_DSN`
  zeigen in `compose.yaml` unverändert alle drei auf denselben Superuser
  `postgres` — der Feed-Container läuft in der Referenzumgebung also
  faktisch mit Superuser-Rechten auf allen drei Kanälen, wie im
  Compose-Kommentar selbst dokumentiert. Das ist jedoch **kein neuer
  Befund dieses Slices**: `git log --follow -- compose.yaml` zeigt, dass
  genau dieser Zustand mit Commit `a32a2c9` (`slice-023`) entstand; keiner
  der vier hier geprüften Commits (`632ddca`/`ba508ed`/`9a84407`/`f8ce37b`)
  fasst `compose.yaml` an. Der Zustand wurde bereits bei `slice-023`
  real geprüft und dort als `INFO-1` klassifiziert
  (`docs/reviews/review-slice-023.md`, „Rollentrennung im
  Compose-Integrationslauf nicht real exerziert"). Eine Eskalation auf
  HIGH ist nach eigener Prüfung nicht gedeckt: `ADR-0047` selbst schließt
  die Frage „wie der Betreiber Login-Identitäten anlegt" ausdrücklich aus
  seinem Geltungsbereich aus (Kontext-Befund 1) — es liegt also kein
  ADR-Verstoß vor, den `slice-023`/`ADR-0047` „bereits als `done`
  geschlossene Zusage" widerspräche; die dortige Zusage („Konfigurationsvertrag
  erweitert … drei Variablen statt einer") ist eingehalten, drei
  *getrennte* DSN-Variablen existieren, nur ihre Werte sind (bewusst,
  begründet) identisch. `LH-QA-SEC-001`…`003` sind
  Fähigkeits-Anforderungen („sollen … können"), deren Messmethode über
  `make test-store` (`internal/bootstrap/roles_wiring_test.go`) real
  erfüllt ist — die Compose-Referenzumgebung ist dafür nicht die
  vorgesehene Nachweis-Ebene. Keine der vier repo-spezifischen
  HIGH-Kategorien (ADR-Verstoß, Traceability, Spec-Stratum, Docker-only,
  Zwei-Quellen-Drift) noch eine der generischen (Sicherheits-Anti-Pattern,
  Korrektheitsfehler im kritischen Pfad, Suppression) trifft auf diesen
  unveränderten, bereits dokumentierten und bereits bewerteten Zustand zu.
  Diese Einschätzung wird hier unabhängig wiederholt, nicht vom
  Implementer-Bericht übernommen.
- `verifizierbar`: ja — `compose.yaml` Zeilen 62-64 zeigen identische
  DSN-Werte; `git log --follow -- compose.yaml` bestätigt Ursprung in
  `a32a2c9` (slice-023), nicht in einem der vier geprüften Commits.
- `klasse`: „Least-Privilege nur in Unit-/Rollen-Tests, nicht im
  End-to-End-Referenzlauf belegt" (identisch mit der Klasse aus
  `review-slice-023.md` INFO-1 — dasselbe Muster, kein neues Auftreten)

## Negativbefunde

- geprüft, ohne Befund: `632ddca` — reiner `git mv`
  (`test/integration/mvp_test.go` → `test/integration/integration_test.go`),
  0 Insertions/0 Deletions, kein Inhalt im selben Commit (`AGENTS.md` §3.3
  eingehalten).
- geprüft, ohne Befund: `9a84407` — Referenz-Nachzug vollständig:
  `Makefile`, `harness/README.md`, `integration_test.go`-Kommentare/
  Skip-Meldungen, `run-integration-tests.sh`-Kopf-/Kreuzverweis-Kommentar,
  `welle-8.md`; `grep -rn mvp_test\.go` über das Repo zeigt danach nur noch
  den vor-/nachher-erzählenden Slice-Plan-Text selbst, keine lebende
  Referenz. `AGENTS.md` trug keine „MVP"-Erwähnung — Behauptung „kein
  Änderungsbedarf" bestätigt.
- geprüft, ohne Befund: Go-Testnamen (`TestMVPCaptureFlow` u. a.) bewusst
  unverändert gelassen — sie beschreiben weiterhin korrekt das Verhalten
  an den `feed_mvp_*`-Fixture-Tabellen; keine Namens-Inkonsistenz.
  `spec/lastenheft.md` unberührt (`Berührte Spec-Stellen: —` im Slice-Kopf
  korrekt).
- geprüft, ohne Befund: `ba508ed` — die drei behaupteten Prüfungen real im
  Skript nachvollzogen: (1) `cdc_reader`-Login scheitert am `INSERT` auf
  `cdc.change` (erwartet `permission denied`), (2) `cdc_admin`-Login ohne
  `REPLICATION`-Attribut scheitert am Verbindungsaufbau im
  Replication-Protokoll (erwartet „permission denied to start WAL
  sender"), (3) `cdc_capture`-Login mit `REPLICATION`-Attribut gelingt als
  Kontrast — alle drei Fälle korrekt mit `set +e`/Status-Prüfung
  abgefangen, kein stiller `set -e`-Abbruch möglich. Einschränkung siehe
  F-1.
- geprüft, ohne Befund: Aufräumen der drei Test-Login-Identitäten
  (`DROP ROLE` am Abschnittsende) — zusätzlich unabhängig von diesem
  Nachlauf abgesichert durch `trap cleanup EXIT` (`$COMPOSE down -v`), das
  bei jedem Ausgang den kompletten Postgres-Container inkl. aller
  Rollen-Test-Fixtures entfernt; kein Leck bei frühem `exit 1`.
- geprüft, ohne Befund: `run-replication-tests.sh` — Plan-Nachzug in §3
  begründet plausibel und real nachvollziehbar, warum diese Datei
  unangetastet bleibt (Compose-Instanz trägt die Gruppenrollen bereits
  über `schema-rollout`); `git show --stat` über alle vier Commits
  bestätigt: die Datei wird von keinem berührt.
- geprüft, ohne Befund: Hard Rule 3.7 (Kommentar-Klassen) — die neuen/
  geänderten Kommentare in `run-integration-tests.sh` und
  `integration_test.go` beschreiben den geltenden Zustand indikativ
  (Zusage-/Kopplungs-/Grenz-Klasse); die Formulierung „unverändert seit
  slice-011/-023" ist ein auflösbarer Herkunfts-Anker (dieselbe Form wie
  `· seit welle-<NN>`, AGENTS.md §3.7-Begründung), keine Chronik-Erzählung
  einer verworfenen Alternative.
- geprüft, ohne Befund: Traceability — alle vier Commit-Betreffs tragen
  keine `SPEC-*`/`ARC-*`-Kennung; `LH-QA-SEC-001` steht im Bodytext jedes
  der vier Commits; `make commit-traceability` lief in dieser Sitzung über
  `HEAD~5..HEAD` grün mit (0 Befunde).
- geprüft, ohne Befund: `f8ce37b` — reine Doku-/Plan-Aktualisierung
  (DoD-Häkchen, Plan-Nachzug), keine Code-/Verhaltensänderung im selben
  Commit; §6/§7 und die Review-Zeile bewusst unangetastet gelassen,
  konsistent mit dem Lifecycle-Zustand `in-progress`.
- geprüft, ohne Befund: `AGENTS.md` §3.1 (Docker-only), §3.2
  (Suppression-Verbot), §3.5 (ADR-Immutabilität), §3.6 (Gate-Lockerung) —
  keine Verstöße; `ADR-0047` wird von keinem der vier Commits editiert.
- geprüft, ohne Befund: `make gates` (lokal, dieser Review-Lauf) —
  `baseline-verify`, `docs-check`, `commit-traceability`, `a-check` alle
  grün, 0 Befunde.
- nicht reproduziert in dieser Review-Sitzung: `make test-integration`
  (3× grün) — der Implementer-Bericht behauptet reale, wiederholte grüne
  Läufe gegen den Compose-Stack; eine eigenständige Reproduktion dieses
  Laufs ist Verifier-Aufgabe (Modul 8/11), nicht Teil dieses
  Maintainability-Reviews.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Closure-Beleg breiter formuliert als
tatsächlicher Testumfang · Least-Privilege nur in Unit-/Rollen-Tests,
nicht im End-to-End-Referenzlauf belegt

## Verdikt

**Merge-blockierend:** ja (F-1) — nicht, weil der neue Testabschnitt
falsch oder wertlos wäre (er belegt real ein echtes serverseitiges
Sicherheitsverhalten), sondern weil die DoD-Zeile und der geplante
Beobachtungs-Register-Beleg für `BEO-PGC/rollen-test-abdeckungsluecken`
Punkt (2) eine engere, korrektere Formulierung brauchen, bevor die
Planner-Rolle bei der Closure (§7) einen Ausgang für dieses Risiko
einträgt — sonst dokumentiert das Register eine stärkere Zusage, als der
Beleg tatsächlich trägt, und das Muster wiederholt sich unbemerkt. F-2
ist nicht merge-blockierend (INFO, unverändert und bereits an anderer
Stelle bewertet).

**Ausdrücklich geprüft und nicht bestätigt:** Der Auftrag benannte die
Möglichkeit, F-2 als eigenständiges HIGH-Finding zu werten. Nach Prüfung
von `ADR-0047` (Kontext-Befund 1, Geltungsbereichs-Ausschluss),
`git log --follow -- compose.yaml` (Ursprung in `slice-023`, nicht in
diesem Diff) und `docs/reviews/review-slice-023.md` (bereits als INFO-1
klassifiziert) wird diese Einschätzung hier **nicht** übernommen — siehe
Begründung unter F-2. Dies ist eine eigenständige, gegenprüfte
Einschätzung dieses Laufs, keine unreflektierte Übernahme der
Implementer- oder Auftrags-Einschätzung.

**Übergabe:** Findings gehen an den Implementer/Planner (F-1: DoD-Zeile
und geplanter Register-Beleg für Punkt (2) vor Closure präzisieren). Die
**Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort
in den Zähler des Beobachtungs-Registers. Dieser Report ist ein
Lauf-Beleg und wird über Läufe hinweg nicht wieder gelesen. Er ersetzt
keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat.
