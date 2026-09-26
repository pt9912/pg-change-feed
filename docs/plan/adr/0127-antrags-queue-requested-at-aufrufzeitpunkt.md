# ADR-0127: Antrags-Queue — `requested_at` ist der Aufrufzeitpunkt (ergänzt ADR-0050)

**Status:** Accepted — Supersedes: keine. Ergänzt
[`ADR-0050`](0050-sql-administration-antragsqueue-und-live-reload.md) (dort steht
„Zeitstempel“ ohne Bedeutung) und legt die Ordnung fest, auf die
[`ADR-0065`](0065-spaltenausschluss-dauerhafter-traeger.md) und
[`ADR-0112`](0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 6
sich beziehen. [`ADR-0113`](0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) bleibt
unberührt (siehe §Kontext).

**Datum:** 2026-09-26

**Autor:** Architect-Agent (Modul 8), Architect-Zug zu Finding F-1 (MEDIUM) des Reviews
`review-slice-transformationen-antragsweg-usecase`; jede Tatsachenaussage trägt ihren
Beleg-Anker (gedruckte Messzeile, §Gemessen) oder ist als hergeleitet gekennzeichnet.

**Bezug:** [`LH-FA-CFG-007`](../../../spec/lastenheft.md) (Transformationen),
[`LH-FA-CFG-005`](../../../spec/lastenheft.md) (Spaltenausschluss),
[`LH-FA-ADM-001`](../../../spec/lastenheft.md) (SQL-Administration),
[`ADR-0050`](0050-sql-administration-antragsqueue-und-live-reload.md),
[`ADR-0065`](0065-spaltenausschluss-dauerhafter-traeger.md),
[`ADR-0112`](0112-transformationsform-deklarative-regeln-vor-persistenz.md),
[`ADR-0043`](0043-schemamigrationen-mit-d-migrate.md) (Ausweichform `nacharbeit-*.sql`)

**Schärft:** [`SPEC-019`](../../../spec/pflichtenheft.md) — die Spalte `requested_at`, die
Ordnung der Verarbeitung und die Klausel K1 („erst entfernen, dann neu setzen“).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

Die Antrags-Queue (`cdc.administration_request`) trägt zwei Leser, die dieselbe Menge von
Anträgen ordnen müssen: die **Verarbeitung** (die Administrations-Goroutine wendet offene
Anträge an) und die **Ableitung** der dauerhaften Stände (Ausschlussstand und Regelstand
aus den `applied`-Zeilen, beim Prozessstart und im Aktivierungs-Zweig). Beide ordnen nach
`(requested_at, administration_request_id)`; stimmt die Ordnung nicht überein, gilt eine
Regel live und fehlt nach dem Neustart, ohne Fehler und ohne Meldung.

`requested_at` hat den Spalten-Default `current_timestamp` (`tools/schema/schema.yaml`), und die
sieben SQL-Funktionen (`tools/schema/nacharbeit-administration.sql`) lassen die Spalte aus. Der
Wert ist damit der **Beginn der Transaktion**. Zwei Aufrufe derselben Transaktion tragen denselben
`requested_at`, und die Kennung (`gen_random_uuid()`) entscheidet — zufällig. Der Fehlertext von K1
(„erst entfernen, dann neu setzen“) legt genau diese Folge nahe; ein Betreiber, der sie in eine
Transaktion legt (`BEGIN … COMMIT`, `psql -1`, eine Migrationsdatei), bekommt sie in beliebiger
Reihenfolge.

`ADR-0113` Festlegung 2 Punkt 1 begründet den Zweitschlüssel `run_id` der Backfill-Aufnahme damit,
dass „`ListPending` heute nur nach `requested_at` ordnet“. Das war eine Aussage über den Code-Stand
zum Zeitpunkt der ADR. Die Entscheidung selbst (Aufnahme in `(requested_at, run_id)`) hängt an
`cdc.backfill_run.requested_at`, einer eigenen Spalte, die der Use Case bei der Annahme aus der Uhr des
Prozesses setzt (`internal/application/usecase/backfill/service.go`, `Clock.Now()` in
`NewQueuedBackfillRun`) und die von dieser ADR nicht berührt wird; die Funktion `cdc.backfill_table`
schreibt nur die Antrags-Zeile (`SPEC-019`, „Für die Antragsart `backfill` heißt `applied` angenommen“). Ein Supersede von `ADR-0113` ist nicht
nötig; der Satz ist überholt, seine Folge (der Zweitschlüssel) bleibt richtig.

### Gemessen

**Die Funktionen in ihrem Text aus `tools/schema/nacharbeit-administration.sql` („alt“) und mit
einer Änderung („neu“: in jedem `INSERT` steht `requested_at` als zweite Spalte, der Wert
`clock_timestamp()` als zweiter Wert), eine Tabelle `cdc.administration_request` mit den Spalten der
Funktionen (Default `current_timestamp` auf `requested_at`), je eine frische PostgreSQL-Instanz:
PostgreSQL 18.6 (`PG_TEST_IMAGE` aus dem `Makefile`) und PostgreSQL 17.11 (Digest aus
`.github/workflows/e2e.yml`), unter dem Superuser; der Arbeitsbaum des Repositories ist
unberührt.** Ordnung der Zeilen je Regel/Spalte: `(requested_at, administration_request_id)`.
Gedruckt:

| Lauf | PostgreSQL 18.6 alt | PostgreSQL 18.6 neu | PostgreSQL 17.11 alt | PostgreSQL 17.11 neu |
|---|---|---|---|---|
| 1000 Transaktionen `BEGIN; remove_transformation(rN); set_transformation(rN); COMMIT` — Zeilen mit falscher Folge (Remove nach Set) / je Regel gleicher `requested_at` | 494 / 1000 | 0 / 0 | 513 / 1000 | 0 / 0 |
| `psql -1`: 300 mal `exclude_column(cN); include_column(cN)` in einer Transaktion — Include vor Exclude | 158 von 300 | 0 von 300 | 142 von 300 | 0 von 300 |
| 500 Transaktionen `disable_table; enable_table; backfill_table` — Transaktionen mit exakter Aufruffolge | 72 von 500 | 500 von 500 | — | — |
| 300 Autocommit-Paare `remove; set` (getrennte Transaktionen) — Remove nach Set | 0 von 300 | 0 von 300 | 0 von 300 | 0 von 300 |
| 20000 Aufrufe `remove_transformation` in einer Transaktion (`DO`-Schleife) — Paare mit nicht streng steigendem `requested_at` / kleinster Abstand | 19999 / 0 | 0 / 8 µs | 19999 / 0 | 0 / 8 µs |

Die Zeile 3 ist nur an PostgreSQL 18.6 gemessen; die Zeilen 1, 2, 4 und 5 an beiden Versionen. Die
Werte 494/513 und 158/142 sind Zufallsstichproben (die Kennung ist eine UUID); die Aussage über sie
ist „etwa die Hälfte“, nicht eine feste Zahl. Der kleinste Abstand von 8 µs ist die Schärfe **dieser**
Maschine bei einer Schleife ohne Netz; dass zwei Aufrufe nie denselben `clock_timestamp()` tragen, ist
damit **nicht bewiesen** (hergeleitet: jeder Aufruf schreibt eine Zeile und sendet `pg_notify`, das dauert
länger als die Auflösung von 1 µs). Deshalb bleibt die Kennung als Zweitschlüssel bestehen.

### Konstraints

- Keine Domänenlogik in SQL ([`ADR-0046`](0046-sql-driving-adapter-lese-schreib-trennung.md)) und
  die Funktionen schreiben ausschließlich die Antrags-Zeile ([`ADR-0050`](0050-sql-administration-antragsqueue-und-live-reload.md)):
  eine Zeitstempel-Vergabe im Insert ist keine Domänenlogik.
- Kein `INSERT` auf die Tabelle außer über die Funktionen (Grants, `SPEC-019`): der Wert von
  `requested_at` entsteht nirgends sonst (hergeleitet aus den Grants; im Go-Code schreibt nur der
  Backfill-Pfad seine eigene Tabelle `cdc.backfill_run`, `git grep -n requested_at internal cmd`
  ohne Testdateien).
- d-migrate 1.3.1 konvergiert eine Änderung an einer bestehenden Tabelle nicht immer
  (`ADR-0043` §Re-Evaluierungs-Trigger, für den CHECK gemessen; für einen geänderten Spalten-Default
  nicht gemessen). Die sieben Funktionen liegen bereits außerhalb des Modells in
  `nacharbeit-administration.sql` als `CREATE OR REPLACE`.

## Entscheidung

Wir wählen **`requested_at` je Funktionsaufruf auf `clock_timestamp()` zu setzen** (in allen sieben
Funktionen, explizit im `INSERT`), die Verarbeitung in **derselben Ordnung** wie die Ableitung
`(requested_at, administration_request_id)` zu führen und die verbleibenden Grenzen zu benennen.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, Grenze „gleicher Transaktionsbeginn: Reihenfolge undefiniert“ in `SPEC-019` und Handbuch festschreiben (Betreiber führt Remove und Set in getrennten Transaktionen aus) | kein SQL-Eingriff | die Folge, die der K1-Text nahelegt, ist in einer Transaktion ein Würfelwurf (494 von 1000 falsch, §Gemessen); `psql -1` und Migrations-Werkzeuge laufen in einer Transaktion; die Grenze wäre eine Falle, die still bleibt (kein Fehler, nur eine fehlende Regel) |
| **B — `requested_at` per `clock_timestamp()` im `INSERT` der sieben Funktionen (gewählt)** | nur `CREATE OR REPLACE` mit unveränderter Signatur; kein Schema-, kein Go-, kein Guard-Eingriff; die Aufruffolge einer Transaktion bleibt erhalten (0 Fehler in 1000, 300, 500 Transaktionen); Autocommit unverändert (0 von 300) | ordnet nach Aufrufzeitpunkt, nicht nach Festschreibung (siehe Konsequenzen, Grenze 1); Zeilen vor der Änderung behalten ihren Transaktionsbeginn |
| C — Spalten-Default in `schema.yaml` auf `clock_timestamp()` | dieselbe Wirkung ohne Änderung der Funktionen | Änderung an einer bestehenden Tabelle über d-migrate, dessen Konvergenz eines geänderten Defaults nicht gemessen ist; der Default greift nur für einen `INSERT` ohne Wert, den außer den Funktionen niemand ausführen darf — kein Gewinn |
| D — monotone Ordnungsspalte (Sequenz oder `identity`) statt der Zeit | ordnet unabhängig von Uhr und Auflösung | neue Spalte auf der bestehenden Tabelle (Konvergenz für `column_name` gemessen, für diese Spalte nicht), Alt-Zeilen brauchen einen Ordnungswert, beide Abfragen, drei Ports und `SPEC-019` ändern sich; dieselbe Grenze 1 (eine Sequenz ordnet nach `nextval`, nicht nach Festschreibung) |
| E — die laufende Bindung nach jedem Antrag aus der Ableitung neu bilden | Live-Stand und abgeleiteter Stand sind gleich per Konstruktion, auch bei überlappenden Transaktionen | Eingriff in die Verarbeitung (Regelprüfung K1–K4 und Bindung), größer als die Lücke; kein Bedarf, solange Grenze 1 benannt ist |

### Festlegung 1 — `requested_at` ist der Aufrufzeitpunkt

Jede der sieben Funktionen (`cdc.enable_table`, `cdc.disable_table`, `cdc.exclude_column`,
`cdc.include_column`, `cdc.backfill_table`, `cdc.set_transformation`,
`cdc.remove_transformation`) schreibt `requested_at` mit `clock_timestamp()` im `INSERT`; der
Spalten-Default `current_timestamp` bleibt unverändert. Aufrufe **derselben Transaktion** tragen
dadurch verschiedene Zeitstempel in der Reihenfolge des Aufrufs.

### Festlegung 2 — die Ordnung der Verarbeitung

Die offenen Anträge werden in `(requested_at, administration_request_id)` verarbeitet — dieselbe
Ordnung, in der die dauerhaften Stände abgeleitet werden. Die Kennung ordnet nur bei gleichem
Zeitstempel und ist dann deterministisch, nicht zeitlich.

### Festlegung 3 — Grenzen

1. **Aufrufzeitpunkt, nicht Festschreibung.** Ordnet die Verarbeitung nach dem Aufruf, aber die
   Transaktion des früher aufgerufenen Antrags schreibt später fest als die eines später aufgerufenen,
   sieht die Queue den früheren erst danach: die Verarbeitung wendet ihn nach dem späteren an, die
   Ableitung davor. Live-Stand und abgeleiteter Stand weichen dann bis zum nächsten Prozessstart ab.
   Das betrifft nur Anträge auf **dieselbe Regel oder Spalte** aus **zeitlich überlappenden**
   Transaktionen mehrerer Sitzungen; die Betreiberregel lautet, solche Anträge nicht überlappend
   abzusetzen. Akzeptiertes Negativ: die Lücke ist schmal (zwei Sitzungen, dieselbe Adresse,
   gegenläufig, offene Transaktion), ihre Schließung ist Option E und größer als der Nutzen; nicht
   erprobt (hergeleitet).
2. **Uhr.** `clock_timestamp()` ist die Serveruhr; ein Rückwärtssprung der Uhr zwischen zwei Aufrufen
   kehrt deren Ordnung um. Dieselbe Eigenschaft hatte `now()` zwischen Transaktionen; nicht erprobt.
3. **Zeilen vor der Änderung** tragen den Transaktionsbeginn; ihre Ties ordnet die Kennung, wie die
   Ableitung sie bisher ordnete. `applied`-Zeilen werden nicht neu geordnet.

## Konsequenzen

- Positiv: `SELECT cdc.remove_transformation(…); SELECT cdc.set_transformation(…);` in einer
  Transaktion wird in dieser Folge verarbeitet und abgeleitet; dasselbe gilt für `exclude_column` und
  `include_column` derselben Spalte und für `disable_table` vor `enable_table`.
- Positiv: kein Eingriff in Schema, Go-Abfragen, Ports, `rolloutguard` (Bekannt-Liste ist nach
  Signaturen geführt, `tools/schema/rolloutguard/guard.go`; die Signaturen bleiben) und den
  Alt-Tag-Lauf (`CREATE OR REPLACE` mit gleicher Signatur; Einordnung hergeleitet, der Lauf
  wird vom Implementer gefahren, Folgepflicht 3).
- Negativ (Grenzen): Festlegung 3 Punkt 1 bis 3.
- Folgepflicht 1: **`SPEC-019` Nachzug** (Spalte `requested_at`, Absatz „Ordnung der Verarbeitung“,
  K1-Klausel) in demselben Commit wie diese ADR.
- Folgepflicht 2: **Funktions-Text** — die sieben Funktionen in
  `tools/schema/nacharbeit-administration.sql` erhalten `requested_at` mit `clock_timestamp()` (der
  Kopfkommentar der Datei nennt die Zeitstempel-Vergabe); die Kommentare in
  `internal/adapters/driven/postgresstorage/queries/queries.go`, die `requested_at` als
  „Transaktionszeitstempel“ beschreiben (`SelectAppliedColumnRequests`,
  `SelectPendingAdministrationRequests`), werden angepasst.
- Folgepflicht 3: **Tests** (Träger: Implementer, Fixrunde des Slice
  `slice-transformationen-antragsweg-usecase`, Zeile der Fitness Function).
- Folgepflicht 4: **Handbuch** — `docs/user/benutzerhandbuch.md` §4 („Spalte vom Ausschluss
  konfigurieren“, künftig der Transformations-Abschnitt) trägt die Aussage „Aufrufe einer Transaktion
  werden in Aufrufreihenfolge verarbeitet; Anträge auf dieselbe Regel oder Spalte nicht aus überlappenden
  Transaktionen absetzen“ (Träger: Planner, beim Handbuch-Slice der Transformationen).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Test, reale PostgreSQL (Paket `postgresstorage`), **noch nicht vorhanden — Folgepflicht 3** | ruft die Funktionen (nicht `INSERT` von Hand) in **einer** Transaktion auf: `remove_transformation` + `set_transformation` derselben Regel, `exclude_column` + `include_column` derselben Spalte, `disable_table` + `enable_table`; über eine Schleife von mindestens 50 Transaktionen (die Kennung ordnet zufällig — ein einzelner Lauf färbt bei falscher Ordnung nur mit Wahrscheinlichkeit 1/2 rot); erwartet: `ListPending` liefert je Transaktion die Aufruffolge, `TransformationRules`/`ExcludedColumns` nach dem Vermerk `applied` in derselben Folge dasselbe Ergebnis wie die Verarbeitung. Rot färbende Mutation: `requested_at`/`clock_timestamp()` aus einer Funktion entfernen (die Zeilen der Funktion tragen wieder den Transaktionsbeginn) — **am SQL-Text erprobt** (Spalte „alt“ in §Gemessen: 494 bzw. 513 von 1000 Transaktionen falsch), am Go-Test **noch nicht gefahren** | `make test-store` |
| Go-Test (Whitebox `internal/bootstrap`, mit realer Queue) | eine Transaktion `remove`+`set` derselben Regel führt zu einem Regelstand, in dem die neue Regel steht, live (Assembler) und nach Neubildung aus `TransformationRules` (Prozessstart-Pfad) gleich; **noch nicht vorhanden**, Träger wie oben | `make test-store` |

## Re-Evaluierungs-Trigger

- **Zwei Sitzungen setzen Anträge auf dieselbe Regel oder Spalte überlappend ab** und der Betreiber
  meldet einen Live-Stand, der nach dem Neustart abweicht: Option E oder eine Ordnung nach
  Festschreibung (`track_commit_timestamp`) prüfen.
- **Die Antrags-Queue erhält eine Antragsart, deren Wirkung von der Reihenfolge abhängt** und die nicht
  über die sieben Funktionen entsteht: die Regel auf sie erstrecken.
- Sonst permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-26 | Accepted — Architect-Zug (Vollmacht des Auftraggebers) zu Finding F-1 des Reviews `review-slice-transformationen-antragsweg-usecase`; die Ordnung an PostgreSQL 18.6 und 17.11 gemessen (§Gemessen) | [`LH-FA-CFG-007`](../../../spec/lastenheft.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0127` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
