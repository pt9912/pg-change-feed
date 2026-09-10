# Review-Report: slice-011 — 2026-09-10

**Review-Art:** Code — geprüft gegen Slice-Plan + ADRs (Maintainability).

**Gegenstand:** Implementer-Commits von slice-011, Range `950922a..HEAD`
(slice-010-Closure-Ende) — zwei Implementer-Commits: `958be4f`
(Plan-Nachzug, Health-Endpoint zurückgestellt) und `1cee5e9`
(Least-Privilege-Rollen + `cdc.metrics`-View + Tests). Die drei
vorangehenden Lifecycle-Commits (`fa59f6e`, `857703e`, `f6246e5`) sind reine
Feld-/Move-Commits, nicht Review-Gegenstand.

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-09-09 (vier
repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen) ·
**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-10

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `docs/plan/planning/in-progress/slice-011-sicherheit-observability.md`
  (mit vollständigem Plan-Nachzug §1/§3/§6/§8)
- `spec/lastenheft.md` §`LH-QA-SEC-001`…003, §`LH-FA-ADM-002`/003,
  §`LH-FA-SST-003`/004, §`LH-QA-OPS-002`/003 · `spec/pflichtenheft.md`
  §`SPEC-007` (`HEALTH_STATES`), §`SPEC-009` (Metriken-Katalog)
- [`ADR-0020`](../plan/adr/0020-http-grpc-optional.md) (HTTP/gRPC optional,
  permanent, Re-Evaluierungs-Trigger „neuer Entscheidungsanlass"),
  [`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
  (SQL-Driving-Adapter: Lese-Views direkt, schreibende Funktionen über
  Inbound Ports)
- `AGENTS.md` §3 Hard Rules (§3.1 Docker-only, §3.2 Suppression-Verbot, §3.7
  Kommentar-Klassen) · `harness/conventions.md` (MR-000/MR-001, genau eine
  Sub-Area `PGC`, Greenfield) · `.a-check.yml` (Layer-Edges, nur Go-Globs)
- Register-Sichtung: `docs/plan/planning/observations/BEO-PGC/` — alle sieben
  Einträge gegen den Plan-Nachzug (§8) geprüft: `lese-doppelquelle` (1×,
  `evidence/slice-010.md`) und `d-migrate-nacharbeit` (2×,
  `evidence/slice-006.md` + `evidence/slice-010.md`) real bestätigt; die
  übrigen fünf ohne Bezug zu diesem Slice
- vorheriger Report: `review-slice-010.md` (F-1 HIGH → `ADR-0046` als
  Folge-ADR, Konflikt-Pfad Modul 8 — erster Praxistest dieses Pfads)

---

## Findings

### F-1 — Health-Endpoint-Rückzug: Begründung deckt nur den HTTP/gRPC-Zweig, nicht die Schicht-Abgrenzung, die tatsächlich trägt

- `kategorie`: MEDIUM
- `quelle`: „Maintainability" (Plan-Argumentationsqualität) ·
  [`ADR-0020`](../plan/adr/0020-http-grpc-optional.md) ·
  [`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
- `pfad`: `docs/plan/planning/in-progress/slice-011-sicherheit-observability.md:59-71`
  (§1) und `:185-195` (§6) · `tools/schema/nacharbeit-observability.sql:20-26`
  (derselbe Wortlaut im SQL-Kommentar) · Commit `958be4f`
- `befund`: Der Plan-Nachzug begründet die Nicht-Realisierung des
  Health-Endpoints damit, dass „eine echte Prozess-Liveness-Antwort einen
  neuen Driving-Adapter-Zuschnitt (HTTP-Listener, Heartbeat-Datei o. ä.)"
  brauche und `ADR-0020` genau das zurückstelle — und leitet daraus eine
  Architect-Pflicht (Folge-ADR) ab. Das trägt nur für den Zweig, den
  `ADR-0020` tatsächlich sperrt (HTTP/gRPC). Ein Heartbeat-*Schreiben* aus
  dem laufenden Prozess in eine Tabelle, gelesen über eine vierte SQL-View
  nach demselben, in diesem Slice bereits verwendeten Muster (`ADR-0046`),
  berührt `ADR-0020` nicht — es ist kein HTTP/gRPC-Driving-Adapter. Der
  tatsächliche Grund, warum dieser Slice keinen Heartbeat liefert, ist die
  im selben §1 für die Rollen-Verdrahtung explizit gezogene
  Schicht-Abgrenzung („dieser Slice hält sich auf die DB-Schicht … die
  Anwendungs-/Verdrahtungsschicht bleibt unberührt"): Ein Heartbeat-Schreiber
  müsste den laufenden Prozess selbst ändern (Anwendungs-/Driven-Schicht) —
  dieselbe Grenze, nicht `ADR-0020`. Die Eskalation an den Architect
  (Folge-ADR nötig) ist damit breiter als die Begründung trägt; der
  eigentliche Grund ist eine Umfangs-/Schicht-Entscheidung, die der Plan an
  anderer Stelle (wiring.go-Ausschluss) bereits selbst so einordnet, hier
  aber „anderer Vorgang" statt „Schicht-Abgrenzung" nennt und daraus eine
  Architect-Pflicht statt eines einfachen Folge-Slice-Verweises macht.
- `verifizierbar`: nein — Plan-Argumentation, kein Gate-Lauf
- `klasse`: Architektur-Rückzug ohne vollständige Alternativen-Prüfung
  (1. Auftreten)

**Konflikt-Pfad (Modul 8):** Kein HIGH-Rollenwiderspruch — der Implementer
hat korrekt nicht selbst entschieden, sondern eskaliert (§6, Reviewer-Hinweis
statt Implementer-Ermessen). Das Finding richtet sich gegen die *Reichweite*
der Eskalation, nicht gegen den Eskalations-Weg selbst. Empfehlung an
Planner/Architect: vor einer Folge-ADR prüfen, ob ein Heartbeat-Muster
(Schreiben: Folge-Slice in der Anwendungsschicht; Lesen: weitere SQL-View
nach `ADR-0046`) den `LH-FA-ADM-002`/`LH-QA-OPS-002`-Bedarf ohne neue ADR
deckt — dann wird aus der Architect-Pflicht ein einfacher Folge-Slice.

### F-2 — §5 Closure-Trigger nach dem Plan-Nachzug nicht mehr erfüllbar, aber unverändert

- `kategorie`: MEDIUM
- `quelle`: Baseline-Regelwerk `modul-05-planning-harness.md` §Closure- und
  Lerneintrag-Regeln
- `pfad`: `docs/plan/planning/in-progress/slice-011-sicherheit-observability.md:166-168`
  (§5, unverändert seit Planner-Anlage) vs. `:59-71` (§1, Plan-Nachzug)
- `befund`: §5 verlangt weiterhin „Health- und Metriken-Endpoint am
  verdrahteten System belegt (beobachtbar am Integrationstest-Diff)" als
  Closure-Bedingung. Nach dem Plan-Nachzug ist der Health-Teil dieser
  Bedingung mit dem aktuellen Scope nicht erfüllbar — §6 benennt das Risiko
  und verweist korrekt an Planner/Reviewer, aber §5 selbst bleibt
  unangepasst. Der Slice steht damit mit einer intern widersprüchlichen
  Closure-Bedingung in `in-progress/`.
- `verifizierbar`: ja — Textvergleich §1/§6 gegen §5
- `klasse`: Stale Closure-Trigger nach Plan-Nachzug (1. Auftreten)

**Empfehlung (explizit angefragt):** Sequenz, kein Parallel-Fix. Erst
**Architect-Verdikt** zur Health-Frage aus F-1 (Heartbeat-Muster ausreichend
ohne Folge-ADR, oder tatsächlich `ADR-0020`-Re-Evaluierung nötig) — davon
hängt ab, ob der Health-Teil per einfachem Folge-Slice nachgeliefert werden
kann oder eine neue ADR braucht. **Danach** passt der **Planner** §5 an
(Metriken-Closure von Health-Closure trennen, oder Carveout mit
Folge-Slice-ID). Ein reiner Planner-Fix ohne das Verdikt liefe Gefahr, den
Closure-Trigger auf Basis einer nicht geprüften Architektur-Prämisse (F-1)
neu zu schreiben.

### F-3 — `cdc_admin`-Grant für Publication-Verwaltung ungetestet und PostgreSQL-Eigentümer-Regel plausibel nicht gedeckt

- `kategorie`: MEDIUM
- `quelle`: `LH-QA-SEC-001`/002 (Least-Privilege je Rolle muss für ihre
  Aufgabe *ausreichen*) · „Maintainability"
- `pfad`: `tools/schema/nacharbeit-roles.sql:65` (`GRANT CREATE ON DATABASE
  … TO cdc_admin`) · `internal/adapters/driven/postgresstorage/tableactivation.go:218,231`
  (`CREATE PUBLICATION …`, `ALTER PUBLICATION … ADD TABLE …`) ·
  `internal/adapters/driven/postgresstorage/roles_test.go` (kein Test deckt
  diesen Pfad)
- `befund`: Der Kommentar in `nacharbeit-roles.sql` begründet das
  `CREATE ON DATABASE`-Grant an `cdc_admin` mit dem Bedarf von
  `CREATE PUBLICATION`/`ALTER PUBLICATION` in `tableactivation.go`. PostgreSQL
  verlangt für `ALTER PUBLICATION … ADD TABLE` (und für `CREATE PUBLICATION …
  FOR TABLE`) zusätzlich zum Datenbank-Privileg Eigentümerrechte an der
  hinzuzufügenden Tabelle — die zu aktivierenden Quelltabellen liegen
  außerhalb des `cdc`-Schemas und gehören `cdc_admin` nach dieser Lieferung
  nicht. `roles_test.go` prüft die drei Rollen ausschließlich gegen
  `cdc`-Schema-Objekte (`source_table`, `change`, `transaction`,
  `cdc.metrics`); kein Test führt `CREATE`/`ALTER PUBLICATION` unter
  `SET ROLE cdc_admin` aus. Ob das Grant für den genannten Zweck tatsächlich
  ausreicht, ist damit unbelegt — und wird erst relevant, sobald
  `wiring.go` die Rolle tatsächlich einsetzt (§6, offenes Risiko), ist aber
  bereits jetzt eine unbewiesene Kern-Behauptung des Kommentars.
- `verifizierbar`: ja — Test, der unter `SET ROLE cdc_admin` eine reale
  `CREATE PUBLICATION …`/`ALTER PUBLICATION … ADD TABLE`-Sequenz gegen eine
  fremde Tabelle ausführt (`make test-store`)
- `klasse`: Ungetesteter Least-Privilege-Kernpfad (1. Auftreten)

---

## Negativbefunde

- geprüft, ohne Befund: **Plan-Nachzug-Vollständigkeit** — `958be4f`
  (13:02:43) liegt vor `1cee5e9` (13:02:48), Plan vor Code im selben Lauf.
  §1/§3/§6/§8 tragen *beide* Abweichungen (Health-Endpoint nicht realisiert,
  `wiring.go` unverändert) mit Begründung und Klasse; §3 deckt exakt die vier
  gelieferten Dateien (`nacharbeit-roles.sql`, `nacharbeit-observability.sql`,
  `Makefile`, `roles_test.go`) — kein ungedeckter Datei-Umfang, zweiter
  Praxistest der Regel seit slice-009/010 besteht strukturell
- geprüft, ohne Befund: **Least-Privilege-Rollentrennung (cdc-Schema)** —
  `cdc_capture` (INSERT `transaction`/`change`, SELECT `source_table`/
  `schema_version`), `cdc_admin` (DML `source_table`/`schema_version`/
  `consumer`/`consumer_position`, SELECT `transaction`/`change`, kein INSERT
  auf `transaction`/`change`), `cdc_reader` (SELECT nur auf die vier Views,
  kein Grant auf eine Basistabelle) sind innerhalb des `cdc`-Schemas
  wechselseitig sauber getrennt; die vier `roles_test.go`-Tests decken je
  eine Grenze in beide Richtungen (Capture kann nicht administrieren, Admin
  nicht erfassen, Reader nicht schreiben/Basistabellen lesen) — strukturell
  plausibel gegen die Grants gelesen, real bestätigt ist Verifier-Sache
- geprüft, plausibel: **Mutationsproben-Behauptung (SQLSTATE 42501)** — die
  vier Tests prüfen konsistent über `strings.Contains(err.Error(), "SQLSTATE
  42501")`, dieselbe Prüfform wie `store_test.go`; die erwarteten
  Fehlpfade (INSERT ohne Grant, SELECT ohne Grant) entsprechen exakt den
  Grants in `nacharbeit-roles.sql` — strukturell plausibel, nicht in diesem
  Review-Lauf erneut ausgeführt (Bestätigung ist Verifier-Sache, `make
  test-store`)
- geprüft, ohne Befund: **`cdc.metrics` als reine Projektion (`ADR-0046`)** —
  alle fünf `UNION ALL`-Zweige aggregieren/joinen bereits persistierte Zeilen
  ohne Schwellenwert- oder Autorisierungsentscheidung; die Klassifikation
  (`SPEC-007` `HEALTH_STATES`) bleibt explizit beim lesenden System (Kommentar
  `nacharbeit-observability.sql:6-8`). `WHERE cs.acknowledged_position IS NOT
  NULL` ist ein Null-Schutz für die Subtraktion, keine Geschäftsregel — siehe
  INFO-Hinweis unten
- geprüft, ohne Befund: **Deckung von `SPEC-009` gegen den Plan-Ausschluss** —
  die fünf gelieferten Metriken (`cdc_transactions_total`,
  `cdc_changes_processed`, `cdc_oldest_change_age_seconds`,
  `cdc_consumer_position`, `cdc_consumer_lag`) und die fünf explizit als
  „nicht abgedeckt" benannten (`cdc_changes_pending`, `cdc_capture_lag`,
  `cdc_errors_total`, `cdc_wal_retention_bytes`, `cdc_storage_bytes`) ergeben
  zusammen genau die zehn `SPEC-009`-Zeilen — keine stillschweigend
  ausgelassene Metrik
- geprüft, ohne Befund: **Kommentar-Klassen** (`AGENTS.md` §3.7) in
  `nacharbeit-roles.sql`, `nacharbeit-observability.sql`, `Makefile`,
  `roles_test.go` — indikativ über den geltenden Zustand (Zusage/Kopplung/
  Abgrenzung/Rang-Zeiger/Grenze), keine Konjunktiv-Klausel über eine
  verworfene Alternative, kein abgebrochener Satz
- geprüft, ohne Befund: **Docker-only** — beide neuen `schema-rollout`-Schritte
  laufen über `docker run … $(PG_TEST_IMAGE) psql …`, kein lokales
  Toolchain-Install; `roles_test.go` läuft über `make test-store`
- geprüft, ohne Befund: **Suppression-Verbot** — keine `nolint`/`noqa`/
  `SuppressMessage`-Marker in den beiden Commits
- geprüft, ohne Befund: **Traceability der zwei Implementer-Commits** — je
  Betreff mindestens eine `LH-*`-/`ADR-*`-Kennung
  (`LH-FA-ADM-002`/`ADR-0020` bzw. `LH-QA-SEC-001`/`LH-FA-SST-004`), keine
  Struktur-ID (`SPEC-*`/`ARC-*`) im Betreff
- geprüft, ohne Befund: **a-check-Edges** — `roles_test.go` liegt im
  `adapters`-Layer-Glob, importiert nur `pgxpool`/Standardbibliothek, keine
  Rückwärts-Kante; die beiden SQL-Dateien liegen außerhalb der Go-Layer-Globs
  von `.a-check.yml` und sind ihm unsichtbar (wie bereits in review-slice-010
  F-1 für die Views festgestellt — hier greift dieselbe Grenze für Rollen/DDL)
- geprüft, ohne Befund: **`internal/bootstrap/wiring.go` unverändert** —
  bestätigt (`git diff 950922a..HEAD -- internal/bootstrap/wiring.go` ist
  leer), deckungsgleich mit der Plan-Nachzug-Behauptung in §1/§3
- geprüft, ohne Befund: **§8 Register-Sichtung, Genauigkeit** —
  `lese-doppelquelle` (1×, `evidence/slice-010.md`) und `d-migrate-nacharbeit`
  (2×, `evidence/slice-006.md` + `evidence/slice-010.md`) real gegen die
  Register-Dateien bestätigt; kein dritter Evidence-Eintrag für
  `d-migrate-nacharbeit` fällig, weil die beiden neuen Nacharbeit-Dateien den
  d-migrate-Rollout-Pfad gar nicht durchlaufen (Rollen sind kein
  `schema.yaml`-Objekt; die View folgt dem etablierten Muster ohne neuen
  Drift-Fund), kein neuer Evidence-Eintrag durch den Implementer (korrekt
  Planner-Sache bei Closure, Modul 8)
- geprüft, ohne Befund: **Sub-Area-Wahl (§8)** — `harness/conventions.md`
  deklariert genau eine Sub-Area (`PGC`, Greenfield); die berührte Fläche
  fällt vollständig darunter, reiner GF-Hinweis korrekt ohne
  Modus-Begründungsblock
- geprüft, ohne Befund: **WIP-Limit und Lifecycle** — genau ein Slice in
  `in-progress/`; die drei vorangehenden Lifecycle-Commits sind reine
  Feld-/Move-Commits mit Kennung im Betreff

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 0 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Architektur-Rückzug ohne vollständige
Alternativen-Prüfung (1. Auftreten) · Stale Closure-Trigger nach
Plan-Nachzug (1. Auftreten) · Ungetesteter Least-Privilege-Kernpfad
(1. Auftreten)

**INFO-Hinweis (zusätzlich zu den drei MEDIUM, siehe Negativbefunde):**
`cdc.metrics`s `WHERE cs.acknowledged_position IS NOT NULL`-Klausel ist ein
Null-Schutz, kein Geschäftsentscheid — aber `ADR-0046`s Fitness-Function-Zeile
benennt genau diese Art View-interner `WHERE`-Klausel als künftige
Review-Prüfpflicht (kein Gate deckt SQL). Für den Architect/nächsten Reviewer
festgehalten, keine Aktion für diesen Slice.

## Verdikt

**Merge-blockierend:** nein für den bereits gelieferten Teil (Rollen +
`cdc.metrics`, F-3 ausgenommen), aber die Slice-Closure ist blockiert, bis
F-1/F-2 aufgelöst sind — §5 kann in seiner jetzigen Form nicht als erfüllt
gelten, ohne dass entweder der Health-Teil doch geliefert wird oder Architect
+ Planner die Closure-Bedingung anpassen. Kein HIGH; die drei MEDIUM-Findings
sind kein Merge-Stopp für die bereits vorliegenden Dateien, aber ein
Closure-Stopp für den Slice als Ganzes.

**Übergabe:** F-1 und F-2 gehen an Planner **und** Architect (Sequenz:
Architect-Verdikt zur Heartbeat-Frage zuerst, Planner passt §5 danach an —
kein Rollenwechsel ohne dieses Übergabe-Artefakt, Modul 8). F-3 geht an den
Implementer (Folge-Zug: Test + ggf. Grant-Korrektur), sobald die
Rollen-Verdrahtung tatsächlich aufgenommen wird (§6, offenes Risiko).

---

**Gate-Beleg:** `make gates` nach diesem Report-Commit (Lauf 2026-09-10,
Range-Head); Ergebnis im Commit-Text vermerkt.
