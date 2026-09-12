# Verifier-Report: slice-036 — 2026-09-13

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §1 (Ziel/Abgrenzung),
§2 (Definition of Done, aktueller Stand nach `3f57e8d`), §6 (Risiken, noch
ohne Ausgang), §8 (Sub-Area-Sichtung), Plan-Nachzug, sowie DoD-/
ADR-Konformität gegen [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(Accepted, zentraler Prüfmaßstab), [`LH-FA-ADM-001`](../../spec/lastenheft.md)
und die beiden Review-Reports
[`review-slice-036.md`](review-slice-036.md) (3 MEDIUM, kein Merge-Blocker)
und [`review-slice-036-fixrunde.md`](review-slice-036-fixrunde.md) (alle drei
behoben, kein neues Finding).

**Grundsatz:** Keine Behauptung übernommen — jeder Beleg unten wurde in
diesem Lauf selbst gelesen oder ausgeführt: vollständiger Slice-Plan
(§1–§8, Plan-Nachzug), `ADR-0050` vollständig, beide Review-Reports
vollständig, der tatsächliche SQL-Quelltext
(`tools/schema/nacharbeit-administration.sql`), der Testdatei-Volltext
(`internal/adapters/driven/postgresstorage/administrationrequest_test.go`),
der `administration_request`-Knoten in `tools/schema/schema.yaml`, das
Beobachtungs-Register (`observation.md`, `state.md`, alle fünf
`evidence/*.md`), `harness/README.md`s `nacharbeit`-Zeile, die
Makefile-Rollout-Reihenfolge, `make gates` einmal und `make test-store`
einmal selbst ausgeführt, zusätzlich die drei genannten Testfälle
**einzeln mit `-v`** in einem eigenen, frisch aufgesetzten
PostgreSQL-18-Testcontainer reproduziert (nicht der Implementer-/
Reviewer-Containerlauf), `git status` am Ende sauber.

**Gegenstand:** `b6a1c9f` (neue Tabelle `cdc.administration_request` in
`tools/schema/schema.yaml`; neue SQL-Funktionen in
`tools/schema/nacharbeit-administration.sql`; Testdatei; Makefile-Anpassung;
Beobachtungs-Beleg), `3cd1f64` (DoD-Häkchen, Plan-Nachzug), `c78f743`
(Review, 3 MEDIUM), `3f57e8d` (Fixrunde: Register-Notiz, Negativtest,
DoD-Begründung), `6f4dbb2` (Fixrunden-Bestätigung, alle drei behoben).

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-13

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-036-antragsqueue-sql-funktionen.md` (vollständig, §1–§8, Plan-Nachzug, aktueller Stand nach `3f57e8d`)
- `docs/reviews/review-slice-036.md` (3 MEDIUM, kein Merge-Blocker)
- `docs/reviews/review-slice-036-fixrunde.md` (alle drei Findings behoben)
- `docs/plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md` (vollständig — zentraler Prüfmaßstab)
- `tools/schema/nacharbeit-administration.sql` (vollständig)
- `internal/adapters/driven/postgresstorage/administrationrequest_test.go` (vollständig)
- `tools/schema/schema.yaml` (`administration_request`-Knoten, Kommentar-Block darüber)
- `docs/plan/planning/observations/BEO-PGC/d-migrate-nacharbeit/observation.md`, `state.md`, alle fünf `evidence/*.md` (`slice-006`, `slice-010`, `slice-015`, `slice-016`, `slice-036`)
- `harness/README.md` (`nacharbeit`-Zeile, `make schema-rollout`)
- `Makefile` (Rollout-Reihenfolge `nacharbeit-*.sql`)
- `docs/plan/planning/welle-12.md` (vollständig — `slice-036` ist der erste, nicht der letzte Slice der Welle; keine Welle-Closure-Reife-Prüfung fällig)

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` | `baseline-verify`: v6.5.0, 54 Dateien OK · `d-check` Standardlauf: 299 Dateien, 0 Befunde · `d-check --range HEAD~5..HEAD` (Modul `commits`): 299 Dateien, 0 Befunde · `commit-traceability.sh`: OK, 5 Commits, Betreffs ohne Struktur-ID · `a-check`: 0 Befunde | **0** |
| `make test-store` | vollständiger Rollout (`nacharbeit-roles.sql` → `nacharbeit-observability.sql` → `nacharbeit-heartbeat.sql` → `nacharbeit-administration.sql`, je `CREATE FUNCTION` ×2/`REVOKE`/`GRANT` ohne Fehler) + komplettes Go-Testpaket grün, einschließlich `internal/adapters/driven/postgresstorage` (2.931s) | **0** |
| eigener, isolierter Testcontainer + `go test -v -run 'TestAdministrationRequest…'` | alle drei benannten Testfälle einzeln `--- PASS` (`EnableTableWritesPendingRequestAndNotifies` 0.03s, `EnableTableRequiresCdcAdminMembership` 0.01s, `DisableTableWritesPendingRequestAndNotifies` 0.02s) | **0** |
| `git status` (am Ende dieses Laufs) | sauber, keine Restspur | — |

Die drei DoD-benannten Testfälle wurden **nicht aus dem Implementer- oder
Reviewer-Bericht übernommen**, sondern in dieser Sitzung in einem eigenen,
frisch erzeugten Docker-Netz/-Container (`cdc-verify-036`) real gegen
PostgreSQL 18 (derselbe gepinnte Digest wie `run-store-tests.sh`)
ausgeführt und danach abgeräumt — genau die Verifier-Lücke, die Modul 11
benennt, wird hier geschlossen.

## Prüfpunkt 1 — Schreiben die SQL-Funktionen wirklich *nur* Antrags-Datensatz + `pg_notify` (`ADR-0050`s zentrale Entscheidung)?

**Ja — eigenständig im SQL-Quelltext nachvollzogen.**
`tools/schema/nacharbeit-administration.sql` enthält für
`cdc.enable_table`/`cdc.disable_table` je genau eine `INSERT INTO
cdc.administration_request (...)` und ein `PERFORM pg_notify('cdc_administration',
v_id)`. Kein `INSERT`/`UPDATE`/`DELETE` auf `cdc.source_table` oder
`cdc.schema_version`, kein `ALTER PUBLICATION`, kein Aufruf einer anderen
schreibenden Funktion oder eines Inbound Ports. Das deckt sich mit `ADR-0050`s
Entscheidungstext („schreibt **ausschließlich** einen Antrags-Datensatz …
und sendet `pg_notify`") und mit der ADR-Fitness-Function-Zeile 1. Real
gegen PostgreSQL bestätigt (`TestAdministrationRequestEnableTable…`/
`DisableTable…WritesPendingRequestAndNotifies`, beide in diesem Lauf
eigenständig PASS): Antrag mit Status `pending` entsteht, `pg_notify` wird
über `LISTEN` empfangen, Payload = Antrags-ID.

## Prüfpunkt 2 — Least-Privilege (`ADR-0047`) wirksam?

**Ja.** `REVOKE EXECUTE ON FUNCTION cdc.enable_table(...), cdc.disable_table(...)
FROM PUBLIC;` läuft nach beiden `CREATE OR REPLACE FUNCTION`-Anweisungen
und vor dem gezielten `GRANT … TO cdc_admin` — kein impliziter
PUBLIC-Zugriff bleibt zwischen den Anweisungen bestehen (ein einziger,
sequenziell abgearbeiteter `psql -f`-Lauf). Der in der Fixrunde neu
hinzugekommene Testfall
`TestAdministrationRequestEnableTableRequiresCdcAdminMembership` in diesem
Lauf eigenständig reproduziert: `SET ROLE cdc_reader`, dann `SELECT
cdc.enable_table(...)` scheitert mit SQLSTATE 42501 (`permissionDenied`) —
`--- PASS` in der eigenen, isolierten Reproduktion.

## Prüfpunkt 3 — d-migrate-Ausweichform (`ADR-0043`-Re-Evaluierungs-Trigger) gerechtfertigt?

**Ja, ohne eigene Neu-Reproduktion des d-migrate-Fehlers** — der Review-Report
hat den `POST_EXECUTE_DRIFT`-Befund bereits zweifach unabhängig bestätigt
(Quellcode-Analyse gegen `MigrationFingerprint.kt`/`RawSqlTextProjection.kt`
sowie eigene isolierte Reproduktion mit einer trivialen No-Arg-Funktion).
Diese Verifikation prüft stattdessen die **Konsequenz** der Ausweichform
eigenständig real: Der `make test-store`-Lauf in dieser Sitzung zeigt den
`schema-rollout` real inklusive `nacharbeit-administration.sql` nach
`nacharbeit-roles.sql` — die Reihenfolge ist notwendig (`cdc_admin` muss
existieren, bevor `GRANT … TO cdc_admin` laufen kann) und real ohne Fehler
bestätigt. Die Antrags-Tabelle selbst läuft unverändert über
`tools/schema/schema.yaml` (`schema validate`: 10 Tabellen gefunden,
inklusive `administration_request`).

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt, Stand nach `3f57e8d`)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | Neue Antrags-Tabelle, real ausgerollt | **bestätigt** | `schema.yaml` (`administration_request`-Knoten, eigenständig gelesen), real ausgerollt in dieser Sitzung (`make test-store`, `schema validate`: 10 Tabellen) |
| 2 | `cdc.enable_table`/`cdc.disable_table` real, schreiben ausschließlich Antrag+Notify | **bestätigt** | Prüfpunkt 1 — Quelltext eigenständig gelesen, alle drei Testfälle einzeln real reproduziert |
| 3 | `make gates` grün | **eigenständig real reproduziert** | Sensor-Tabelle oben |
| 4 | Review durchgeführt, Report liegt vor | **bestätigt** | `review-slice-036.md` (`c78f743`) + `review-slice-036-fixrunde.md` (`6f4dbb2`), alle drei MEDIUM behoben — Checkbox im Plan-Text selbst noch `[ ]`, korrekt (Planner-Closure-Schritt zieht sie nach) |
| 5 | Doku-Update `harness/README.md`/`AGENTS.md`, falls neuer Sensor/Vertrag | **bestätigt, entfällt korrekt** | eigener `grep -n "nacharbeit" harness/README.md`: nur `nacharbeit-views.sql` genannt, ausschließlich im Auflösungs-Kontext („zurückgebaut · seit slice-016“); die drei aktiven Ausweichformen (`roles`, `observability`, `heartbeat`) kommen kein einziges Mal vor — `nacharbeit-administration.sql` fällt unter dasselbe Muster, kein neues Gate/Target entstanden |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **noch offen — korrekt, Planner-Closure-Schritt** | §7 ist Platzhalter |
| 7 | Reconciliation-Register, falls einschlägig | **entfällt korrekt** | `../reconciliation.md` existiert nicht (Greenfield, `harness/conventions.md` MR-000) |
| 8 | Beobachtungs-Register fortgeschrieben | **bestätigt** | `evidence/slice-036.md` existiert, `state.md` zieht den Zähler korrekt auf 5× nach (eigenständig gegen `ls evidence/` geprüft — siehe unten) |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **noch offen — korrekt, Planner-Closure-Schritt** | beide §6-Risiken tragen noch `<bei Closure einzutragen>` |
| 10 | Drei Paarungen getragen | **noch offen — korrekt, Welle-12-Closure vorbehalten** | Slice-Kopf nennt `welle-12`; `slice-036` ist deren *erster*, nicht letzter Slice — Paarungen sind Sache der Welle-Closure |

**Zwischenstand:** Die drei Kern-DoD-Punkte, die der Implementer als
erledigt beansprucht (1–3), sind eigenständig real nachgeprüft und tragen
vollständig. Punkt 5 wurde durch die Fixrunde korrigiert und trägt jetzt
selbst mit eigenständig geprüfter Begründung. Punkt 8 ist erfüllt. Die vier
offen gelassenen Punkte (4-Häkchen im Plan-Text, 6, 9, 10) sind korrekt der
Planner-Closure vorbehalten (Modul 5/6/8) und kein Implementer- oder
Reviewer-Versäumnis.

## Eigenständige Prüfung des Beobachtungs-Registers (nicht aus den Review-Reports übernommen)

Alle fünf `evidence/*.md`-Dateien selbst gelesen:

- `slice-006.md` — CHECK-Constraint-Fund (`raw-sql-text-drift`, d-migrate 1.2.0).
- `slice-010.md` — dieselbe Grenze trifft die drei Views.
- `slice-015.md` — CHECK-Constraint konvergiert (1.3.0), Views bleiben offen.
- `slice-016.md` — Views konvergieren (1.3.1, `source_dialect`+`columns:`-Signatur), technisch aufgelöst.
- `slice-036.md` — dritte Objektklasse (Funktionen): DDL-Generierung korrekt, `schema migrate --execute` bricht mit `POST_EXECUTE_DRIFT` ab; Ausweichform bleibt.

`state.md` fasst das korrekt zusammen: „Zähler (abgeleitet): 5×" mit allen
fünf Dateinamen — deckt sich mit dem tatsächlichen Verzeichnisinhalt
(`ls evidence/` liefert exakt diese fünf Dateien, eigenständig geprüft).
Die Drei-Klassen-Aufteilung (CHECK aufgelöst seit slice-015, Views aufgelöst
seit slice-016, Funktionen weiterhin offen) ist gegen die fünf
Evidence-Dateien deckungsgleich, keine Über- oder Untertreibung. Kein
Widerspruch zwischen `observation.md` (Identität, unverändert seit Anlage)
und `state.md` (veränderlicher Stand).

## Negativbefunde

- geprüft, ohne Befund: **§1-Abgrenzung gewahrt** — kein Zugriff auf
  `cdc.source_table`/`cdc.schema_version`/Publication im Diff, keine
  Verarbeitung der Anträge (keine Goroutine, kein Aufruf von
  `EnableTableUseCase`/`DisableTableUseCase`), keine CLI-Diagnose/
  Consumer-Verwaltung — exakt der in §1 benannte Scope.
- geprüft, ohne Befund: **Kein Accepted-ADR-Inhalt geändert** (Hard Rule
  3.5) — `ADR-0050` bleibt unverändert `Accepted`, keiner der fünf
  slice-036-Commits ändert eine Datei unter `docs/plan/adr/`.
- geprüft, ohne Befund: **Kommentar-Disziplin (`AGENTS.md` §3.7)** — die
  Kommentare in `nacharbeit-administration.sql` und `schema.yaml`
  beschreiben den geltenden Zustand im Indikativ mit auflösbaren Ankern
  (`ADR-0050`, `ADR-0043`, `LH-FA-ADM-001`, `BEO-PGC/d-migrate-nacharbeit`),
  keine Chronik, kein Vorwärtsverweis auf `slice-037`/`slice-038` namentlich.
- geprüft, ohne Befund: **Traceability** — Commit-Betreffs tragen `ADR-0050`,
  kein `SPEC-*`/`ARC-*` im Betreff; `commit-traceability` (Teil von
  `make gates`) lief in diesem Lauf über `HEAD~5..HEAD` grün.
- geprüft, ohne Befund: **`make gates`** — eigener Lauf, 0 Befunde (299
  Dateien Standardlauf, 299 Dateien Range-Lauf, `commit-traceability` OK,
  `a-check` 0 Befunde). Die Differenz zu den im DoD-Text genannten „54/296"
  Dateien ist erwartete Drift durch die seither hinzugekommenen
  Review-Report-Dateien, kein Befund dadurch berührt.
- geprüft, ohne Befund: **`make test-store` und die drei benannten
  Testfälle einzeln** — real reproduziert in einem eigenen, isolierten
  Testcontainer (nicht der Implementer-/Reviewer-Lauf), alle drei `PASS`.
- geprüft, ohne Befund: **Rollout-Reihenfolge im Makefile** —
  `nacharbeit-roles.sql` läuft vor `nacharbeit-administration.sql` (eigener
  `grep -n nacharbeit Makefile`), notwendig weil `cdc_admin` dort entsteht.
- geprüft, ohne Befund: **Arbeitsverzeichnis am Ende dieses Laufs** —
  `git status --porcelain` leer; eigener Testcontainer/-netz
  (`cdc-verify-036`/-pg) nach Gebrauch abgeräumt.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 (eigene Befunde) |
| INFO | 0 |

**Zusammenfassung DoD:** Die drei implementierten Kern-Punkte (1–3) sowie
Punkt 5 (nach Fixrunde) und Punkt 8 sind eigenständig real nachgeprüft und
tragen vollständig. Die vier noch offenen Punkte (4-Häkchen im Plan-Text, 6,
9, 10) sind korrekt der anstehenden Planner-Closure vorbehalten und kein
Implementer- oder Reviewer-Versäumnis. Kein DoD-Defekt im Sinn eines
unbelegten „bestätigt"-Punkts.

## Verdikt

**DoD-Konformität (soweit vom Implementer beansprucht): bestätigt, ohne
Einschränkung.** `ADR-0050`s zentrale Entscheidung — die beiden
SQL-Funktionen schreiben ausschließlich einen Antrags-Datensatz und senden
`pg_notify`, ohne `cdc.source_table`/`cdc.schema_version`/die Publication zu
berühren — ist eigenständig im Quelltext verifiziert und real getestet
bestätigt. Die Least-Privilege-Durchsetzung (`ADR-0047`) ist wirksam und seit
der Fixrunde durch einen realen Negativtest abgesichert. Die
d-migrate-Ausweichform ist gerechtfertigt (bereits im Review zweifach
unabhängig belegt, hier auf Konsequenz-Ebene — realer Rollout-Erfolg —
eigenständig nachgeprüft).

**Eigene Reproduktion:** `make gates` einmal, `make test-store` einmal, plus
alle drei DoD-benannten Testfälle **einzeln mit `-v`** in einem eigenen,
frisch aufgesetzten Testcontainer — alle vier Läufe grün, real in diesem
Kontext ausgeführt, nicht aus dem Implementer- oder Reviewer-Bericht
übernommen.

**Plan-vs-Code-Diff:** Der Diff (`b6a1c9f`, `3f57e8d`) hält sich an §3 des
Slice-Plans (Schema-Tabelle, SQL-Funktionen, Testdatei) und an die im
Plan-Nachzug ehrlich benannten Implementer-Entscheidungen (Tabellenname,
d-migrate-Einschränkung mit Ausweichform, Rollenwahl `cdc_admin`,
`SECURITY DEFINER`, Kanal-Name). Keine unangekündigte Schicht- oder
Umfangs-Erweiterung; §1s Abgrenzung (keine Antrags-Verarbeitung, kein
direktes Schreiben der Zielzeilen, keine CLI-Diagnose/Consumer-Verwaltung)
wird nicht berührt.

**Closure-Bereitschaft:** Die Verifikation der real gelieferten Punkte ist
abgeschlossen und positiv. Vor dem `git mv` nach `done/` fehlen noch die
Planner-Closure-Schritte: DoD-Häkchen 4 im Plan-Text nachziehen,
§7-Closure-Notiz schreiben, §6-Risiko-Ausgänge eintragen. Die drei
Paarungen bleiben regulär der `welle-12`-Closure vorbehalten (`slice-036`
ist deren erster, nicht letzter Slice — keine Welle-Closure-Reife-Prüfung
in diesem Lauf fällig).

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Slice-Plan,
Register und Roadmap wurden von diesem Lauf nicht verändert (`git status`
am Ende sauber).

---

**Gate-Beleg:** `make gates` einmal in diesem Lauf ausgeführt, Exit 0, 0
Befunde (`baseline-verify`: v6.5.0, 54 Dateien OK; `d-check` Standardlauf:
299 Dateien, 0 Befunde; `d-check --range HEAD~5..HEAD` Modul `commits`: 299
Dateien, 0 Befunde; `commit-traceability.sh`: OK, 5 Commits; `a-check`: 0
Befunde). `make test-store` einmal ausgeführt, Exit 0, komplettes
Go-Testpaket grün. Zusätzlich alle drei DoD-benannten Testfälle einzeln mit
`-v` in einem eigenen, isolierten Testcontainer reproduziert (Exit 0, alle
drei `PASS`). `git status` am Ende dieses Laufs sauber.
