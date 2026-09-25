# `tools/bench-backfill.sh` — Messung des Backfills (Teil von `make bench`)

## Vertrag

Das Skript misst einen Backfill-Run gegen eine reale PostgreSQL und einen
laufenden Feed-Container ([`LH-FA-CAP-009`](../../spec/lastenheft.md),
[`ADR-0111`](../../docs/plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md),
[`ADR-0113`](../../docs/plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
Festlegung 3). Es ist eine **Messung ohne Pass/Fail**: die Richtgröße ist sein
Ergebnis, nicht seine Vorbedingung, und keine gedruckte Zahl wird gegen eine
Schwelle gehalten. Der Exit-Code sagt nur, ob die Messung gültig war (jeder Run
endet `completed` und kopiert genau die Zeilen der Tabelle). Es hängt an keinem
`GATE_CHECKS`-Eintrag und ist ein **Werkzeug** wie die drei übrigen
`tools/bench-*.sh`: eigene PostgreSQL-/Feed-Umgebung über
`tools/bench-lib.sh` (`docker network`/`docker run`, Schema-Rollout über
d-migrate, `docker rm -fv` beim Abbau), braucht ein geladenes `:dev`-Image
(`make image`).

Die Datei `docs/user/bench-abdeckung.md` bindet Kennungen an durchgesetzte
Schwellen; dieses Skript schreibt dorthin keine Zeile.

## Ablauf

Jede gedruckte Zeile beginnt mit `bench-backfill[<Lauf>]:` (der Lauf ist der
UTC-Zeitstempel des Starts) und nennt die Größe, auf die sie sich bezieht.

1. **Schätzung der Zeilenzahl.** Zwei Tabellen mit `BENCH_BACKFILL_EST_ROWS`
   Zeilen; gelesen wird `pg_class.reltuples` mit derselben Abbildung wie der
   Snapshot-Leser (ein Wert unter 0 heißt „unbekannt"): frisch befüllt
   (Autovacuum an, Autovacuum aus), nach Autovacuum (Wartezeit bis zur ersten
   Schätzung, höchstens 150 s), nach `ANALYZE`, und nach 20 % weiteren Zeilen
   ohne erneutes `ANALYZE`.
2. **Kopierdauer je Tabellengröße.** Je Stufe aus `BENCH_BACKFILL_STAGES` eine
   Tabelle (`id`, `name`, `amount`, `created_at`, `note`, Autovacuum aus, damit
   die Schätzung des Runs „unbekannt" und die Messung frei von Autovacuum ist)
   und `BENCH_BACKFILL_RUNS` Runs über `cdc.backfill_table`. Gemessen wird
   `finished_at − started_at` der Run-Zeile — die Wartezeit in `queued` zählt
   nicht, die Dauer des Commits liegt hinter `finished_at`. Je Run stehen
   Dauer, Zeilen je Sekunde, die geschätzte Zeilenzahl, beide Warn-Kennzeichnungen,
   die Speicher-Spitze des Feed-Containers und die WAL-Messung in der Ausgabe; je
   Stufe der Bereich (Minimum–Maximum) und der Median, dazu der Speicher 20 s nach
   dem letzten Run der Stufe. Die **WAL-Messung** liest je Statusabfrage zwei
   Größen aus `pg_replication_slots`: den WAL-Rückstand des Capture-Slots
   (`pg_current_wal_lsn()` minus `confirmed_flush_lsn`, dieselbe Größe wie
   `cdc_wal_retention_bytes`) und das vom Slot gehaltene WAL (`pg_current_wal_lsn()`
   minus `restart_lsn`); gedruckt werden je Run die Spitze, der Rückstand
   unmittelbar nach dem Run und der Rückstand nach der Wartezeit ohne jeden
   Schreibzugriff (bis unter die Warnschwelle, höchstens 120 s — die
   Leerlauf-Bestätigung des Streams senkt ihn,
   [`ADR-0120`](../../docs/plan/adr/0120-capture-slot-leerlauf-bestaetigung.md)),
   je Stufe der Median der Spitzen samt abgeleiteter Bytes je Zeile, nach der
   letzten Stufe die höchste Spitze der größten Stufe im Vergleich zur
   Warnschwelle (100 MiB, [`SPEC-013`](../../spec/pflichtenheft.md)). Die
   Blockgröße `B` und die Toleranz liest das
   Skript aus dem Code (`DefaultBlockSize` in
   `internal/adapters/driven/postgressnapshot/snapshot.go`,
   `copyDurationToleranceMinutes` in
   `internal/application/usecase/backfill/warn.go`), damit keine zweite Kopie
   entsteht; ist eine nicht lesbar, endet das Skript mit Exit 1.
3. **Wirkung auf die Live-Erfassung.** Ein Schreiber fügt
   `BENCH_BACKFILL_LIVE_RATE` Zeilen je Sekunde in eine eigene, aktivierte
   Tabelle ein; `cdc_capture_lag` (`cdc.metrics`) wird im Sekundentakt gelesen,
   einmal während eines Runs über die größte Stufe (und 10 s danach), einmal
   als Referenz gleicher Dauer ohne Run; die Zeile nennt zusätzlich die
   WAL-Spitzen (Rückstand und gehaltenes WAL) im Lauf mit Run.
4. **Richtgröße (abgeleitet).** Rate der größten Stufe (Median) mal Toleranz in
   Sekunden, abgerundet auf eine Stelle (die erste Ziffer bleibt, der Rest wird
   `0`). Die Zahl ist **abgeleitet**, solange keine Stufe die Toleranz ausfüllt;
   die Ausgabe sagt es.

## Parameter

| Variable | Default | Bedeutung |
|---|---|---|
| `BENCH_BACKFILL_STAGES` | `10000 50000 200000` (`--full`: `10000 100000 1000000`) | Zeilenzahlen der Stufen, aufsteigend |
| `BENCH_BACKFILL_RUNS` | `3` | Runs je Stufe (Bereich und Median) |
| `BENCH_BACKFILL_EST_ROWS` | `100000` | Zeilen der Tabellen der Schätz-Messung |
| `BENCH_BACKFILL_LIVE_RATE` | `100` | Live-Zeilen je Sekunde der Live-Phase |
| `BENCH_BACKFILL_RUN_TIMEOUT_S` | `1800` | Sekunden (Uhr der Shell) je Run, bevor die Messung mit Exit 1 abbricht |

Die Umgebung ist die von `tools/bench-lib.sh`: PostgreSQL-Image aus
`PG_TEST_IMAGE` mit den Startparametern `wal_level=logical`,
`max_replication_slots=10`, `max_wal_senders=10`, `wal_sender_timeout=2000`,
sonst PostgreSQL-Defaults; Feed-Container aus `FEED_IMAGE` ohne
Ressourcenbegrenzung.

## Grenzen

Was das Skript **nicht** misst — eine Aussage darüber ist ungedeckt:

- **Host und Umgebung.** Jede Zahl gilt für den Host, den die erste Ausgabezeile
  nennt (Kernel, Docker-Version, CPU-Zahl, RAM), und für PostgreSQL-Defaults;
  andere Hardware, Plattengeschwindigkeit und parallele Last anderer Container
  sind nicht kontrolliert.
- **Zeilenbreite.** Die Tabellen tragen fünf schmale Spalten (Ausgabe nennt die
  mittlere Zeilenbreite in Bytes). Breite Zeilen (`jsonb`, `bytea`) sind ohne
  eigenen Lauf ungemessen; eine Richtgröße in Zeilen gilt nur für diese Breite.
- **Einfügeform.** Die Kopierdauer schließt die zeilenweise Einfügung eines Blocks
  in **eine** Schreibtransaktion ein (`AppendBlock`); eine andere Einfügeform
  verändert die Zahl.
- **Blockdauer.** Die Ausgabe nennt die **mittlere** Blockdauer (Median-Dauer durch
  Blockzahl, abgeleitet); die längste Einzeldauer eines Blocks und die Dauer des
  Commits sind nicht gemessen.
- **Speicher.** Gemessen wird der Feed-Container über `docker stats` im Takt der
  Statusabfrage (etwa 1 bis 2 Sekunden); die Spitze kann zwischen zwei Proben
  liegen, und die Probe 20 s nach dem letzten Run der Stufe kann darüber liegen
  (der höchste gedruckte Wert einer Stufe ist der höchste gemessene Wert, nicht
  die Spitze allein). Der Speicher der PostgreSQL-Instanz ist nicht gemessen.
- **WAL.** Die Messung ist ohne Pass/Fail und ohne Ursachenzuordnung: sie
  liest Rückstand und gehaltenes WAL des Slots, nicht die Ursache eines Werts.
  Die Wartezeit ohne Schreibzugriff ist auf 120 s begrenzt. Das Skript schreibt
  zwischen den Runs nichts, was den Slot bestätigen ließe: eine Regression der
  Leerlauf-Bestätigung bleibt als hoher Rückstand in der Ausgabe sichtbar.
- **Live-Wirkung.** Eine Live-Last, eine Live-Tabelle, ein Run zugleich; die
  Referenz ohne Run hat dieselbe Dauer, aber keine Wiederholung.
- **Extrapolation.** Die Richtgröße rechnet die Rate der größten Stufe auf die
  Toleranz hoch; keine Stufe kopiert die volle Toleranzdauer.
