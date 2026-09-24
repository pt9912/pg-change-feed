# ADR-0119: Backfill — Wirkung der Lesesperre auf wartende DDL und Beleg für `RENAME COLUMN` (Supersedes ADR-0118, teilweise)

**Status:** Accepted — Supersedes [`ADR-0118`](0118-backfill-umschreiben-im-snapshot-fenster.md)
in genau drei Stellen: §Kontext „Nicht gemessen“ (der Teilsatz zur
Sperr-Warteschlange), §Konsequenzen (der zweite Negativ-Punkt, Klammer
„Sperr-Warteschlange, hergeleitet, nicht gemessen“) und §Entscheidung
Festlegung 5 (der Satzteil „E2E-belegt“). Alles Übrige von `ADR-0118` bleibt in
Kraft, insbesondere die Festlegungen 1 bis 4 und der Ablauf Import → Sperre →
Vergleich → Spaltenliste → Cursor.

**Datum:** 2026-09-25

**Autor:** Planner-Agent (Modul 8), Closure von `slice-backfill-e2e`; die Aussagen
sind an den Messungen des Review-Reports und des Verifikations-Reports belegt
(siehe §Kontext).

**Bezug:** [`LH-FA-CAP-009`](../../../spec/lastenheft.md) (Haupt-Bezug),
[`ADR-0118`](0118-backfill-umschreiben-im-snapshot-fenster.md) (teilweise
superseded — Haupt-Bezug), [`ADR-0111`](0111-backfill-bestand-snapshot-bulk-copy.md)
(nur zur Abgrenzung: Ablauf des Runs)

**Schärft:** —

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0118`](0118-backfill-umschreiben-im-snapshot-fenster.md) ist `Accepted` und
unberührbar (`AGENTS.md` §3.5). Zwei seiner Aussagen tragen eine andere Reichweite
als die Messungen, die danach am Gegenstand gefahren wurden; die Berichtigung ändert
den Referenten und ist keine Zitat-Korrektur.

### Gemessen

**Wirkung der Sperre auf wartende DDL.** Zwei unabhängige Läufe an PostgreSQL 18
(Pin von `PG_TEST_IMAGE` im `Makefile`, Wegwerf-Container): Sitzung A hält
`BEGIN; LOCK TABLE … IN ACCESS SHARE MODE`, Sitzung B führt
`ALTER TABLE … ALTER COLUMN name TYPE varchar(64)` aus und wartet. Ein danach
gestarteter `INSERT INTO t …`, ein `SELECT count(*) FROM t` und
`SELECT count(*) FROM pg_publication_tables WHERE pubname = …` laufen je in eine
4-s-Grenze (Review-Report F-1, gedruckt „exit=143 after 4063 ms“, „4053 ms“,
„4056 ms“; Verifikations-Report §5 F-1, gedruckt „exit=124 after 4005 ms“ für alle
drei); `pg_stat_activity` zeigt alle vier Sitzungen mit `wait_event_type='Lock'`.
Schreiber, Leser und die Publication-Abfrage der Administration (im Antrag und am
Run-Start) stehen also hinter der wartenden DDL. PostgreSQL 17: nicht gemessen.

**`RENAME COLUMN` im Fenster.** Ein Scratch-Test des Reviews (PostgreSQL 18, nicht
committet) endet am Cursor mit der Klasse `storage`
(„column "name" does not exist (SQLSTATE 42703)“); Review-Report F-7. Ein Test oder
eine Runner-Phase des Repos fährt `RENAME COLUMN` nicht: die Phase DDL-Fenster und
`TestNoRewriteInWindowReadsTheSnapshot` fahren `DROP COLUMN`.

Reports (Verzeichnis `docs/reviews/`): `review-slice-backfill-e2e`,
`verifikation-slice-backfill-e2e`.

### Konstraints

- Die Träger, die Betreiber lesen, tragen die gemessene Wirkung: Handbuch §4
  „Bestand als Backfill überführen“, Absatz „Sperre der Tabelle“.
- Eine `Accepted`-ADR ist Beleg für nachgeordnete Träger; ihre Aussage über den
  Gegenstand deckt sich mit der Messung, auf die sie sich stützt
  (`AGENTS.md` §3.12 Instanz B).

## Entscheidung

Wir wählen **eine neue ADR mit teilweisem `Supersedes`**, die die zwei Aussagen auf
die gemessene Reichweite setzt. Zwei Festlegungen.

### Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, Handbuch und Plan tragen die Messung | keine Änderung an der ADR | die ADR bleibt die einzige Stelle, die „nur Leser, nicht gemessen“ und „E2E-belegt“ sagt; jeder, der sich auf sie stützt, liest die engere bzw. breitere Aussage |
| B — die Sätze in `ADR-0118` in-place ändern | ein Ort | ändert §Konsequenzen und §Entscheidung, keine Zitat-Korrektur (`AGENTS.md` §3.5, [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)) |
| C — einen E2E-Lauf mit `RENAME COLUMN` ergänzen, damit „E2E-belegt“ wahr wird | Aussage und Beleg decken sich ohne Berichtigung | erweitert den Lauf von `make test-integration` um Minuten; `DROP COLUMN` belegt dieselbe Klasse (SQLSTATE `42703`, `storage`) |
| **D — neue ADR mit teilweisem `Supersedes` (gewählt)** | die Aussagen tragen ihre Messung; `ADR-0118` bleibt unberührt | eine weitere ADR für zwei Sätze |

### Festlegung 1 — Wirkung der Sperre auf wartende DDL

Solange eine DDL, die `ACCESS EXCLUSIVE` verlangt, auf die Lesesperre des Runs wartet,
warten alle nach ihr angefragten Zugriffe auf die Tabelle hinter ihr: Schreiber
(`INSERT`), Leser (`SELECT`) und die Abfrage der Publication-Mitgliedschaft
(`pg_publication_tables`) der Administration. Die Wartezeit reicht bis zum Ende der
Sperre des Runs oder zum Abbruch der DDL (hergeleitet aus der Reihenfolge der
Sperranforderungen; die Messung endet an der 4-s-Grenze). Gemessen an
PostgreSQL 18; PostgreSQL 17 ist nicht gemessen. Diese Festlegung ersetzt den Teilsatz zur Sperr-Warteschlange in
§Kontext „Nicht gemessen“ und den zweiten Negativ-Punkt in §Konsequenzen von
`ADR-0118`.

### Festlegung 2 — `RENAME COLUMN` im Fenster

`DROP COLUMN` im Fenster endet am `DECLARE` mit `42703`, Klasse `storage`, und ist im
E2E-Lauf belegt (Phase DDL-Fenster). `RENAME COLUMN` im Fenster endet am `DECLARE`
mit derselben Klasse; dieser Ausgang ist im Review gemessen und im Repo ohne Test.
Festlegung 5 von `ADR-0118` gilt in der Sache unverändert; ihr Satzteil „E2E-belegt“
gilt für `DROP COLUMN`.

## Konsequenzen

- Positiv: die Aussagen der ADR-Kette zur Sperre und zu `RENAME COLUMN` decken sich
  mit den Messungen und mit dem Handbuch.
- Negativ: `ADR-0118` trägt die berichtigten Sätze weiter im Text; wer sie liest,
  liest diese ADR über den Index (Titel „Supers. ADR-0118, teilw.“).
- Folgepflicht: keine — Handbuch §4 „Sperre der Tabelle“ und der Plan von
  `slice-backfill-e2e` (§6) tragen die Wirkung bereits.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| — | nicht maschinell prüfbar; die Wirkung trägt die Messung (Review-/Verifikations-Report), das Handbuch trägt sie für Betreiber | — |

## Re-Evaluierungs-Trigger

- **Eine Zeitgrenze der Sperranweisung wird eingeführt**
  (`BEO-PGC/lesesperre-ohne-zeitgrenze`): die Wartezeit aus Festlegung 1 neu bewerten.
- **`RENAME COLUMN` bekommt einen E2E-Beleg**: Festlegung 2 auf „E2E-belegt“ setzen.
- Sonst permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-25 | Accepted — Berichtigung zweier Aussagen von `ADR-0118` an den Messungen des Review- und des Verifikations-Reports zu `slice-backfill-e2e` | [`LH-FA-CAP-009`](../../../spec/lastenheft.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0119` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
