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
   Stufe der Bereich (Minimum–Maximum) und der Median, dazu der Speicher 20 s und
   60 s nach dem letzten Run der Stufe und die Zahl der Zeilen in `cdc.change` vor
   der Stufe und nach ihr. Vor der ersten Stufe druckt das Skript die
   **Grundlinie ohne Run**: den Speicher des Feed-Containers in Ruhe 30 s nach dem
   Start und die Zahl der Zeilen in `cdc.change`. Die **WAL-Messung** liest je Statusabfrage zwei
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
Ressourcenbegrenzung. Zwei Variablen der Umgebung gelten für jedes Skript, das
`bench::start_feed` ruft: `BENCH_FEED_ENV` (zusätzliche Umgebungsvariablen des
Feed-Containers, leerzeichengetrennt, `KEY=VALUE`) und `BENCH_FEED_DOCKER_ARGS`
(zusätzliche Argumente von `docker run`, etwa `--memory 6g`). Ein Kill des
Feed-Containers durch die Grenze ist ein Messergebnis in den zwei Skripten dieses
Vertrags: beide enden mit Exit 1 und drucken den Zustand des Containers (`Status`,
`Exit`, `OOMKilled`), `tools/bench-backfill.sh` in der Statusabfrage eines Runs,
`tools/bench-backfill-memory.sh` beim nächsten Zähler-Zugriff. Die übrigen
`tools/bench-*.sh` prüfen den Zustand des Feed-Containers nicht.

## Speicher-Untersuchung — `tools/bench-backfill-memory.sh`

Ein zweites Skript, **kein Teil von `make bench`** (Minuten je Stufe):
dieselbe Umgebung
(`tools/bench-lib.sh`), dieselbe Haltung — eine Messung ohne Pass/Fail. Es misst,
wie sich der Speicher des Feed-Containers im Backfill-Run und danach zur Zeilenzahl
der Tabelle, zur Zeilenbreite, zur Zahl der gespeicherten Changes und zu einer
Einstellung des Feeds verhält.

1. **Tabellen.** Je Stufe aus `BENCH_MEM_STAGES` eine Tabelle: `narrow` mit den
   fünf schmalen Spalten des Skripts `tools/bench-backfill.sh`, `wide` mit zwei
   weiteren Spalten (`payload jsonb`, `body text`); die Ausgabe nennt die
   mittlere Zeilenbreite je Tabelle in `pg_column_size` und in der Textform.
2. **Einheit der Messung.** Je Stufe und Run startet das Skript einen **frischen**
   Feed-Container; damit ist `memory.peak` (siehe unten) die Spitze genau dieser
   Einheit. Mit `BENCH_MEM_RESET=1` (Default) leert es vor jeder Stufe
   `cdc.change` und `cdc.transaction`; Run `k` einer Stufe beginnt dann mit
   `(k−1) ×` Stufengröße Zeilen in `cdc.change`, die Ausgabe nennt die Zahl vor
   dem Run und nach dem Nachlauf.
3. **Grundlinie, Run, Nachlauf.** 10 s Ruhe nach dem Start (Beginn, Maximum,
   Minimum, Ende), dann der Run mit einer Probe etwa alle 0,2 s (Verlauf bei
   0 %, 25 %, 50 %, 75 %, 100 % der kopierten Zeilen, Spitze), dann 60 s ohne
   Schreibzugriff mit einer Probe alle 0,5 s (Beginn, Maximum, Minimum, Ende).
   Mit `BENCH_MEM_EMPTY_AFTER=1` leert das Skript danach `cdc.change` und
   `cdc.transaction` und tastet weitere 120 s ab.
4. **Zähler.** Vom Host aus gelesen, ohne Eingriff in den Container: die
   cgroup-v2-Datei des Containers (`memory.current`, `memory.stat` mit `anon`,
   `file`, `kernel`, `memory.peak`), `/proc/<pid>/status` des Prozesses (`RssAnon`,
   `VmHWM`) und `docker stats`. Vorbedingung ist cgroup v2 mit dem systemd-Treiber
   von Docker: der Pfad ist `/sys/fs/cgroup/system.slice/docker-<Container-ID>.scope`.
   Ist die cgroup-Datei oder `/proc/<pid>` nicht lesbar (etwa auf einem Host mit dem
   cgroupfs-Treiber), endet das Skript mit Exit 1 und der Meldung „cgroup (…) oder
   /proc/… nicht lesbar“.
5. **Laufzeit des Go-Prozesses.** Mit `BENCH_FEED_ENV="GODEBUG=gctrace=1"` druckt
   das Skript je Fenster die Zeilen der Garbage-Collection des Feeds: Zahl der
   Läufe, größter Heap zu Beginn eines Laufs, größtes Ziel, letzter Heap nach
   einem Lauf; außerdem die Zahl der gelaufenen und der fehlgeschlagenen
   Bereinigungs-Takte seit dem Start des Feeds.
6. **Slot.** Vor jedem Start wartet das Skript bis zu 10 s darauf, dass die
   Sitzung des entfernten Feeds den Slot freigibt; danach beendet es sie
   (`pg_terminate_backend`) und nennt es in der Ausgabe.

| Variable | Default | Bedeutung |
|---|---|---|
| `BENCH_MEM_STAGES` | `10000 100000 200000` | Zeilenzahlen der Stufen |
| `BENCH_MEM_RUNS` | `3` | Runs je Stufe |
| `BENCH_MEM_WIDTH` | `narrow` | `narrow` oder `wide` |
| `BENCH_MEM_RESET` | `1` | `cdc.change` vor jeder Stufe leeren |
| `BENCH_MEM_EMPTY_AFTER` | `0` | `1`: nach dem Nachlauf leeren und 120 s weitermessen |
| `BENCH_MEM_SERIES_DIR` | leer | Verzeichnis für die Zeitreihe je Run (`ms cur anon file kernel rssanon rows`) |
| `BENCH_MEM_RUN_TIMEOUT_S` | `3600` | Sekunden je Run, bevor die Messung mit Exit 1 abbricht |

Grenzen dieses Skripts: `memory.current` zählt `anon`, `file` und `kernel` des
Containers; die Spalte `anon` ist der Speicher des Go-Prozesses. `memory.peak` ist
der Höchstwert seit dem Start des Containers und schließt den Nachlauf ein, wenn
er nach dem Nachlauf gedruckt wird. Die Proben im Run liegen 0,2 s auseinander
(die Zeit einer Abfrage kommt hinzu); eine kürzere Spitze fängt nur `memory.peak`.
Die Zeilenzahl in `cdc.change` wird mit `count(*)` zwischen den Fenstern gelesen.

## Grenzen

Was das Skript **nicht** misst — eine Aussage darüber ist ungedeckt:

- **Host und Umgebung.** Jede Zahl gilt für den Host, den die erste Ausgabezeile
  nennt (Kernel, Docker-Version, CPU-Zahl, RAM), und für PostgreSQL-Defaults;
  andere Hardware, Plattengeschwindigkeit und parallele Last anderer Container
  sind nicht kontrolliert.
- **Zeilenbreite.** Die Tabellen dieses Skripts tragen fünf schmale Spalten (Ausgabe
  nennt die mittlere Zeilenbreite in Bytes). Kopierrate und Richtgröße bei breiten
  Zeilen (`jsonb`, `bytea`) sind ungemessen; eine Richtgröße in Zeilen gilt nur für
  diese Breite. Den Speicher bei breiten Zeilen misst
  `tools/bench-backfill-memory.sh` (`BENCH_MEM_WIDTH=wide`, eine Zeile von etwa
  1,3 KB).
- **Einfügeform.** Die Kopierdauer schließt die zeilenweise Einfügung eines Blocks
  in **eine** Schreibtransaktion ein (`AppendBlock`); eine andere Einfügeform
  verändert die Zahl.
- **Blockdauer.** Die Ausgabe nennt die **mittlere** Blockdauer (Median-Dauer durch
  Blockzahl, abgeleitet); die längste Einzeldauer eines Blocks und die Dauer des
  Commits sind nicht gemessen.
- **Speicher.** Dieses Skript misst den Feed-Container über `docker stats` im Takt
  der Statusabfrage (etwa 1 bis 2 Sekunden); die Spitze kann zwischen zwei Proben
  liegen, und die Probe 20 s oder 60 s nach dem letzten Run der Stufe kann darüber
  liegen (der höchste gedruckte Wert einer Stufe ist der höchste gemessene Wert,
  nicht die Spitze allein), weil der Speicher des Feeds an der Zahl der Changes in
  `cdc.change` hängt und nicht nur am Run. Die Untersuchung des Speichers
  (`tools/bench-backfill-memory.sh`, oben) trennt beides. Der Speicher der
  PostgreSQL-Instanz ist nicht gemessen.
- **WAL.** Die Messung ist ohne Pass/Fail und ohne Ursachenzuordnung: sie
  liest Rückstand und gehaltenes WAL des Slots, nicht die Ursache eines Werts.
  Die Wartezeit ohne Schreibzugriff ist auf 120 s begrenzt. Das Skript schreibt
  zwischen den Runs nichts, was den Slot bestätigen ließe: eine Regression der
  Leerlauf-Bestätigung bleibt als hoher Rückstand in der Ausgabe sichtbar.
- **Live-Wirkung.** Eine Live-Last, eine Live-Tabelle, ein Run zugleich; die
  Referenz ohne Run hat dieselbe Dauer, aber keine Wiederholung.
- **Extrapolation.** Die Richtgröße rechnet die Rate der größten Stufe auf die
  Toleranz hoch; keine Stufe kopiert die volle Toleranzdauer.
