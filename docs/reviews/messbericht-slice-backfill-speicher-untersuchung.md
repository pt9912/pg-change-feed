# Messbericht: Speicher des Feed-Containers im Backfill

**Rolle:** Implementer (Modul 9) — Messbericht, **kein** Review und keine
Verifikation. **Datum:** 2026-09-25. **Gegenstand:** Untersuchung der Speicher-Spitze
des Feed-Containers im Backfill-Run
([`LH-FA-CAP-009`](../../spec/lastenheft.md),
[`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md),
[`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)).
**Werkzeug:** `tools/bench-backfill-memory.sh` und `tools/bench-backfill.sh` (Vertrag:
[`harness/targets/bench-backfill.md`](../../harness/targets/bench-backfill.md)). Die
gedruckten Zeilen aller Läufe stehen in
[`messbericht-slice-backfill-speicher-untersuchung-zeilen`](messbericht-slice-backfill-speicher-untersuchung-zeilen.md);
jede Zahl dieses Berichts nennt Reihe und Lauf und ist dort auflösbar.

## 1. Ergebnis

1. **Ursache (belegt).** Der Speicher des Feed-Containers nach einem Backfill entsteht
   **nicht im Run**, sondern im periodischen Retention-Lauf: alle 10 Sekunden liest
   er **alle Changes der Quelle** samt beider Row Images in eine Liste
   (`RunRetentionService.Run` ruft `ReadChanges` mit der Quelle als einzigem Filter;
   die Abfrage `SelectChanges` trägt `LIMIT NULL`). Der Speicher hängt deshalb an der
   **Zahl der Changes in `cdc.change`** — nicht an der Größe der kopierten Tabelle
   und nicht daran, ob ein Backfill oder die laufende Erfassung sie geschrieben hat
   (gemessen ist nur der Backfill als Quelle; die laufende Erfassung ist aus dem Code
   gelesen, nicht gemessen, Abschnitt 9).
   Ein Backfill füllt `cdc.change` in einem Zug und macht die Abhängigkeit sichtbar.
   Der Beleg sind zwei Schalter (Abschnitt 4): Bereinigung aus → der Speicher bleibt
   flach; `cdc.change` leeren → der Live-Heap fällt.
2. **Im Run** bleibt der Speicher flach: Spitze des Go-Prozesses 8,9 bis 10,5 MiB
   (`anon`) bei 10.000 bis 1.000.000 Zeilen, solange `cdc.change` leer ist
   (Abschnitt 3.2). Der Blockbedarf von etwa 74 KB ist nicht die Ursache.
3. **Nach dem Run** steigt er mit den Takten der Bereinigung auf etwa 1,0 bis
   1,6 KiB je Change bei schmalen Zeilen (etwa 70 Bytes) und 2,65 bis 4,2 KiB bei
   Zeilen von etwa 1,3 KB (abgeleitet, Abschnitt 3.3 und 3.5): 1.000.000 Changes
   ergeben Spitzen von 1.083 und 1.274 MiB (Reihen B und C, Run 1).
4. **Ausgang.** Ein behebbarer Defekt im Code: Änderungs-Slice
   `slice-retention-lauf-speicher-begrenzung` (Datei in `open/`); die Lösung steht in
   [`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md).
   Das Benutzerhandbuch trägt unter „Grenzwerte“ die gemessene Grenze des Standes dieses
   Berichts.
5. **Warn-Richtgröße** (4.000.000 geschätzte Zeilen, `warn.go`): sie bezieht den
   Speicher nicht ein. Ein Backfill dieser Größe hinterlässt 4.000.000 Changes; bei
   1,03 bis 1,59 KiB je Change sind das 3,9 bis 6,1 GiB Speicher des Feed-Containers
   (abgeleitet), bei Zeilen von 1,3 KB 10,1 bis 16,0 GiB. Der Wert im Code bleibt in
   diesem Slice unverändert (eine Änderung ist Code); die Neubemessung ist ein
   Liefer-Punkt des Änderungs-Slice, nach dem Trigger von
   [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
   „Eine Messung liegt vor“ (Abschnitt 6).
6. **Server-Release.** Der Defekt steht in `v0.1.0` bis `v0.1.2`; Änderungs-Slice zuerst,
   danach `v0.2.0` (Nutzer-Entscheidung, Abschnitt 7).

## 2. Umgebung und Herkunft der Zahlen

**Host** (erste Ausgabezeile jeder Reihe): Linux 6.8.0-139-generic, Docker 29.8.1,
20 CPU, 33.362.599.936 Byte RAM, cgroup v2 mit systemd-Treiber; PostgreSQL 18
(`postgres:18-alpine`, Digest aus `tools/bench-lib.sh`) mit den Startparametern von
`tools/bench-lib.sh`. **Feed-Image:** `ghcr.io/pt9912/pg-change-feed:dev`, Image-ID
`sha256:969b7fad…cf4e` (gebaut mit `make image` zu Beginn der Untersuchung, der
Quelltext der Produktion entspricht dem Stand des Ausgangs-Commits); für die
Schalter lokal gebaute Varianten (Abschnitt 4), deren Image-IDs die Ausgabezeile der
Reihe nennt.

**Herkunft.** *Gemessen* heißt: gedruckt in einer Zeile der Reihe (Lauf-Kennung
`bench-backfill-memory[<Lauf>]`), aufgelöst im Zeilen-Dokument. *Abgeleitet*: aus
gemessenen Zahlen gerechnet, die Rechnung steht daneben. *Übernommen*: aus dem
Benutzerhandbuch, der Lauf ist im Repository nicht auflösbar — nur die Zahl
1.544 MiB (Lauf `20260924T233628Z`, Abschnitt 3.8).

**Wie gemessen wird.** Je Stufe und Run ein frischer Feed-Container (`memory.peak`
gilt damit je Einheit); der Speicher wird vom Host aus gelesen: cgroup-v2-Zähler des
Containers (`memory.current`, `memory.stat`, `memory.peak`) und `/proc/<pid>/status`;
`GODEBUG=gctrace=1` liefert die Garbage-Collection-Zeilen des Go-Prozesses
(Laufzeit-Ausgabe, kein Eingriff in den Code). `cdc.change` und `cdc.transaction`
werden vor jeder Stufe geleert (`BENCH_MEM_RESET=1`); Run `k` einer Stufe beginnt
mit `(k−1) ×` Stufengröße Changes. Jede Zahl der Tabellen ist die eines Runs;
Bereiche stehen über die Runs derselben Bedingung, `n` nennt die Zahl.

| Reihe | Lauf | Bedingung |
|---|---|---|
| A | `20260925T123600Z` | Default-Image, schmal, Stufen 10.000, 100.000, 200.000, je 3 Runs |
| B | `20260925T130048Z` | Default-Image, schmal, 1.000.000, 3 Runs, `--memory 6g` |
| C | `20260925T131359Z` | wie B, `--memory 8g` (Wiederholung, mit Zählung der Bereinigungs-Takte) |
| D | `20260925T132801Z` | Default-Image, schmal, 100.000, 3 Runs, danach `cdc.change` leeren und 120 s messen |
| E | `20260925T134215Z` | Image ohne Bereinigung (`retentionInterval` = 1000 h), schmal, wie A |
| F | `20260925T140625Z` | wie E, Blockgröße 100, Stufe 200.000 |
| G | `20260925T141522Z` | wie E, Blockgröße 10.000, Stufe 200.000 |
| H | `20260925T142349Z` | Default-Image, **breit** (Zeile etwa 1,3 KB), Stufen wie A |
| I | `20260925T144833Z` | wie E, breit, Stufe 200.000 |
| J | `20260925T145726Z` | wie E, schmal, 1.000.000, `--memory 6g` |
| K | `20260925T151031Z` | Default-Image mit `GOGC=25`, schmal, Stufe 200.000 |
| L1, L2 | `20260925T151925Z` (L1, mit `GODEBUG`), `20260925T152159Z` (L2, ohne) | Default-Image, schmal, 100.000, 2 Runs, `--memory 64m` (Mutationen, Abschnitt 8) |
| M | `20260925T152516Z` | Default-Image, schmal, Stufen 10.000 und 20.000, 1 Run, `BENCH_MEM_RESET=0` (Mutation) |
| N | `20260925T153145Z` | `tools/bench-backfill.sh --full`, Default-Image, 3 Runs je Stufe, `--memory 10g` |
| O | `20260925T152357Z` | Kopie von `bench-backfill-memory.sh` mit verändertem cgroup-Pfad (Mutation) |
| P1 bis P3, Q1, Q2, R | `20260925T163946Z`, `20260925T170101Z`, `20260925T170249Z`, `20260925T165323Z`, `20260925T165457Z`, ohne Lauf | Läufe der Fixrunde: Mutationen der zwei Skripte (Abschnitt 8), Zeilen im Zeilen-Dokument |

Die Varianten E bis G, I und J ändern **genau eine Konstante** am Quelltext des
Arbeitsbaums (`retentionInterval`, bei F und G zusätzlich `DefaultBlockSize`) und
sind lokal gebaut (`docker buildx build`, Tag `pg-change-feed:bench-…`); der
Arbeitsbaum ist danach zurückgenommen, nichts davon ist committet.

## 3. Messreihe

Mittlere Zeilenbreite laut Ausgabe: schmal 73 bis 74 Bytes (`pg_column_size`),
64 bis 69 Bytes in der Textform; breit 1.273 bis 1.279 Bytes, 1.258 bis 1.265 Bytes
in der Textform (Reihe H).

### 3.1 Grundlinie ohne Run

In den 19 Einheiten mit leerem `cdc.change` (Reihen A bis K) steht der Feed-Container
10 s nach dem Start bei `anon` 5,0 bis 5,5 MiB; `memory.current` liegt in 15 Einheiten
bei 5,8 bis 6,8 MiB und in vier Einheiten bei 18,9 bis 22,4 MiB (jeweils der erste
Container einer lokal gebauten Variante, Reihen E, F, G und I: Datei-Seiten des
Image-Zugriffs, in der Zeile „Run-Ende“ der Reihe E, Stufe 10.000, Run 1 mit `file`
16,0 MiB; `anon` unverändert). Beispiel: Reihe A, Stufe 10.000, Run 1:
`memory.current` 6,2 MiB, `anon` 5,2 MiB. Die „Ruhe“ von 176 bis 196 MiB vor der
Stufe mit 200.000 Zeilen im Handbuch (übernommen) ist mit dem Zustand nach Runs
vereinbar, nicht mit dem eines frischen Containers: vor dieser Stufe standen aus den
Stufen 10.000 und 50.000 (je 3 Runs) 180.000 Changes in `cdc.change`
(abgeleitet: 3 × 10.000 + 3 × 50.000), bei 1,03 bis 1,59 KiB je Change 181 bis
279 MiB (abgeleitet, Abschnitt 3.3).

### 3.2 Verlauf im Run (Speicher über die Blöcke)

Spitze im Run bei leerem `cdc.change` (Run 1 der Stufe; `anon` = Go-Prozess):

| Zeilen der Tabelle | Reihe, Lauf | `anon` Spitze im Run | `memory.current` Spitze im Run |
|---|---|---|---|
| 10.000 | A | 8,9 | 10,6 |
| 100.000 | A; D (3 Runs) | 9,8; 9,9, 9,6, 9,8 | 11,5; 11,5, 11,6, 11,6 |
| 200.000 | A | 10,1 | 12,0 |
| 1.000.000 | B; C | 10,5; 10,1 | 12,7; 12,3 |

Über 10.000 bis 1.000.000 Zeilen liegt die Spitze zwischen 8,9 und 10,5 MiB (`anon`,
`n` = 8): kein Wachstum mit der Zeilenzahl. Innerhalb dieser Runs bleibt `anon` bei
0 %, 25 %, 50 %, 75 % und 100 % der kopierten Zeilen im Bereich von 5,0 bis 10,5 MiB
(Zeilen „Kopierdauer … Proben …“ der Runs mit leerem `cdc.change`, Reihen A bis D):
ein Plateau ab dem ersten Block, kein Wachstum je Block. `GODEBUG=gctrace=1` zeigt im Run einen Heap von höchstens 5
bis 6 MB zu Beginn eines Laufs (Reihe B, Run 1: „größter Heap zu Beginn 5 MB“,
1.661 GC-Läufe im Run).

### 3.3 Nach dem Run: Abhängigkeit von den Changes in `cdc.change`

Speicher-Spitze (`memory.peak` seit dem Start des Containers, gedruckt nach dem
60-s-Nachlauf) gegen die Zahl der Changes, die nach dem Run in `cdc.change` stehen
(schmale Zeilen; `je Change` abgeleitet: Spitze in KiB durch Changes; die letzte Zeile
nennt ihre Zahl selbst):

| Changes in `cdc.change` | Spitze in MiB (Reihe, Run) | je Change in KiB |
|---|---|---|
| 10.000 | 23,9 (A) | 2,45 |
| 20.000 | 40,4 (A) | 2,07 |
| 30.000 | 50,2 (A) | 1,71 |
| 100.000 | 133,8 (A); 151,3, 145,3, 147,7 (D) | 1,37; 1,55, 1,49, 1,51 |
| 200.000 | 245,0, 273,9 (A) | 1,25; 1,40 |
| 300.000 | 449,2 (A) | 1,53 |
| 400.000 | 468,6 (A) | 1,20 |
| 600.000 | 605,9 (A) | 1,03 |
| 1.000.000 | 1.082,7 (B); 1.273,5 (C) | 1,11; 1,30 |
| 2.000.000 | 2.269,2 (B, Run 2); 3.058,0 (C, Run 2) | 1,16; 1,57 |
| 2.000.000 vor dem Run, 3.000.000 danach (Run 3, Spitze im Run) | 3.096,6 (B, Run 3); 2.751,7 (C, Run 3) | 1,59; 1,41 |

Die zwei Runs 3 der Reihen B und C tragen 2.000.000 Changes **vor** dem Run und 3.000.000
danach; ihre Spitze liegt im Run (bei 25 % der kopierten Zeilen, Abschnitt 3.3), nicht im
Nachlauf. Der Wert je Change teilt deshalb durch 2.000.000, die Zahl vor dem Run
(abgeleitet: 3.096,6 × 1.024 / 2.000.000 = 1,585 und 2.751,7 × 1.024 / 2.000.000 = 1,409;
die Zahl der Changes zum Zeitpunkt der Spitze liegt zwischen 2.000.000 und 2.250.000, der
Wert ist damit eine Obergrenze).

Ab 100.000 Changes liegt der Wert je Change zwischen 1,03 und 1,59 KiB (`n` = 15);
darunter überwiegt die Grundlinie von etwa 6 MiB. Die Spitze tritt in den Takten der
Bereinigung auf (alle 10 s, ein Sägezahn: Reihe A, Stufe 200.000, Run 3: `anon` im
Nachlauf Minimum 233,4 und Maximum 508,9 MiB); ein Lauf nach dem anderen lässt die
Grundlinie des nächsten Runs steigen (Reihe B, Run 2: `anon` zu Beginn 1.198 MiB bei
1.000.000 Changes vor dem Run).

**Zum Vergleich mit dem Handbuch:** die dortigen 1.544 MiB im zweiten Run über
1.000.000 Zeilen (Lauf `20260924T233628Z`, übernommen) sind mit den Messungen dieser
Reihen vereinbar: Reihe B, Run 2 (1.000.000 Changes vor dem Run) 1.538,7 MiB
`memory.peak` bis zum Ende des Runs; Reihe C, Run 2 1.656,0 MiB.

Bei 2.000.000 Changes vor dem Run (Run 3 der Reihen B und C) endet die Beobachtung
anders: die Spitze im Run liegt bei 3.096,6 (B) und 2.751,7 MiB (C) `memory.current`
(Reihe B: 3.063,2 MiB bei 25 % der Zeilen, 69,9 MiB bei 50 %; Reihe C: Spitze vor 25 %
der Zeilen, 60,3 MiB bei 25 %), danach steht `anon` bei 62,0 und 54,8 MiB und
bleibt dort 60 s ohne Schreibzugriff: es laufen keine weiteren Bereinigungs-Takte
(Reihe C: „Bereinigung gelaufen 5, fehlgeschlagen 0“ seit dem Start des Feeds, kein
GC-Lauf im Nachlauf). Die Ursache ist **nicht belegt**. Die Abfrage in PostgreSQL erklärt
das Ausbleiben nach dem Messstand des Architect-Verdikts
(`architect-verdict-retention-lauf-speicher-begrenzung`) nicht: 1,9 bis 2,0 s je
1.000.000 Changes im Client-Lauf der bisherigen Abfrage (übernommen aus diesem Verdikt,
Einzelläufe), linear etwa 4 s je 2.000.000 (abgeleitet), gegen ein Fenster von mehr als
60 s ohne Takt; die Hypothese dort liegt im Feed-Prozess (Dekodierung und Garbage
Collection bei einem Heap von etwa 3 GiB) und ist nicht gemessen. Auch in diesem Bericht
fehlt ein `pg_stat_activity`-Lauf.

**Gegenprobe (Reihe N).** Ein Feed-Container über alle Stufen (`--memory 10g`) fährt
Run 3 der Stufe 1.000.000 mit 2.330.000 Changes davor und steht 60 s nach dem Run bei
2.448,4 MiB (`docker stats`, Abschnitt 3.8), 20 s nach dem Run bei 2.737,2 MiB: dort
laufen die Bereinigungs-Takte bei mehr als 2.000.000 Changes weiter, `anon` fällt nicht
auf etwa 60 MiB wie in den Reihen B und C. Das Ausbleiben hängt damit nicht allein an
der Zahl von 2.000.000 Changes; der Unterschied der Anordnung (frischer Container je Run
gegen ein Container über alle Stufen, `docker stats` gegen `anon`) ist nicht als Ursache
belegt.

### 3.4 Blockgröße (Bereinigung aus)

Reihen E, F, G, Stufe 200.000, Run 1 bis 3 (Median, Bereich der Spitze im Run,
`anon`); `n` = 3 je Größe:

| Blockgröße B | `anon` Spitze im Run in MiB | Reihe |
|---|---|---|
| 100 | 8,4 bis 8,6 (Median 8,5) | F |
| 1.000 | 9,8 bis 10,1 (Median 9,9) | E |
| 10.000 | 24,3 bis 24,4 (Median 24,4) | G |

Die Spitze im Run wächst mit der Blockgröße (rund 1,6 KiB je zusätzliche Zeile,
abgeleitet aus (24,4 − 9,9) MiB durch 9.000 Zeilen) und liegt bei der Blockgröße 1.000
zwischen 9,8 und 10,1 MiB: der Blockbedarf ist klein und flach über der Grundlinie
von etwa 5 MiB.

### 3.5 Zeilenbreite

Reihe H (Default-Image, breit, ohne Leerung zwischen den Runs einer Stufe wie A):

| Changes in `cdc.change` | Spitze in MiB (Run) | je Change in KiB |
|---|---|---|
| 10.000 | 49,5 | 5,07 |
| 100.000 | 361,0 | 3,70 |
| 200.000 | 725,5; 818,1 | 3,71; 4,19 |
| 300.000 | 1.033,0 | 3,53 |
| 400.000 | 1.058,8 | 2,71 |
| 600.000 | 1.550,6 | 2,65 |

Ab 100.000 Changes 2,65 bis 4,19 KiB je Change (`n` = 6), gegenüber 1,03 bis 1,59 KiB
bei schmalen Zeilen: rund das Zwei- bis Dreifache bei etwa dem Siebzehnfachen der
Zeilenbreite. Die Spitze im Run bei leerem `cdc.change` bleibt flach (Run 1 der Stufen
10.000, 100.000, 200.000: `anon` 14,4, 14,7 und 14,9 MiB); ohne Bereinigung (Reihe I,
Stufe 200.000, `n` = 3) liegt sie bei 15,2 bis 15,7 MiB gegenüber 9,8 bis 10,1 MiB bei
schmalen Zeilen (Reihe E): die Zeilenbreite hebt die Spitze im Run um etwa 5,5 MiB, der
Speicher nach dem Run aber mit der Zahl der Changes.

### 3.6 Feed-Einstellung `GOGC`

Reihe K (`GOGC=25`, Stufe 200.000, schmal, `n` = 3): `memory.peak` 224,5, 415,5 und
477,4 MiB bei 200.000, 400.000 und 600.000 Changes gegenüber 273,9, 468,6 und 605,9 MiB
mit der Standardeinstellung (Reihe A): 11 bis 21 % weniger (abgeleitet, drei Paare:
224,5 gegen 273,9 MiB = −18,0 %, 415,5 gegen 468,6 MiB = −11,3 %, 477,4 gegen 605,9 MiB
= −21,2 %; Reihe K und Reihe A, Zeilen „Zeilen in cdc.change nach dem Nachlauf“ der Stufe
200.000, Run 1 bis 3). Die
Abhängigkeit von der Zahl der Changes bleibt unverändert; die Einstellung senkt den
Zuschlag der Garbage Collection, nicht den Bestand.

### 3.7 Nach dem Run: gibt der Prozess Speicher frei?

Reihe D (`cdc.change` nach dem Nachlauf geleert, 120 s ohne Schreibzugriff, `n` = 3):
in jedem der drei Fenster läuft genau eine Garbage Collection („im Fenster: GC-Läufe 1“)
und lässt einen Heap von 1 MB zurück („letzter Heap nach GC 1 MB“); `anon` sinkt in
120 s von 96,9, 103,1 und 128,7 MiB auf 54,0, 86,7 und 106,7 MiB: der Prozess gibt den
Speicher nur langsam an das Betriebssystem zurück (das Fenster endet mit 56 bis 84 % des
Werts zu Beginn, abgeleitet).

### 3.8 `tools/bench-backfill.sh --full`

Reihe N (Lauf `20260925T153145Z`, Stufen 10.000, 100.000, 1.000.000, je 3 Runs, ein
Feed-Container über alle Stufen, `--memory 10g`, Messung über `docker stats` im Takt
der Statusabfrage, Zeilen im Zeilen-Dokument). Die Ausgabe trägt seit diesem Slice die
**Grundlinie ohne Run** und die Zahl der Changes in `cdc.change`:

| Stufe (Zeilen) | Ruhe vor der Stufe in MiB | Spitze der Stufe in MiB | 20 s / 60 s nach dem letzten Run in MiB | Changes in `cdc.change` nach der Stufe |
|---|---|---|---|---|
| Grundlinie (30 s nach dem Start) | 6,1 | — | — | 0 (davor) |
| 10.000 | 6,1 | 25,5 | 37,8 / 36,1 | 30.000 |
| 100.000 | 41,7 | 231,9 | 349,6 / 243,6 | 330.000 |
| 1.000.000 | 275,7 | 2.738,2 | 2.737,2 / 2.448,4 | 3.330.000 |

Die Kopierdauer der Stufe 1.000.000 liegt bei 109.110 bis 133.442 ms (Median 129.620
ms, `n` = 3), der Durchsatz bei 7.494 bis 9.165 Zeilen/s (Median 7.715). Das Ergebnis
des Bench stimmt mit den Reihen A bis D überein: „Ruhe vor der Stufe“ folgt der Zahl der
Changes aus den früheren Stufen (6,1 → 41,7 → 275,7 MiB bei 0, 30.000 und 330.000
Changes), und die Probe 20 s nach dem Run liegt über der Spitze im Run
(Stufe 100.000: 349,6 gegen 231,9 MiB). Der Speicher der Stufe 1.000.000 steigt im
dritten Run auf 2.738,2 MiB bei 2.330.000 Changes vor dem Run (abgeleitet: 330.000
plus 2 × 1.000.000), das sind 1,2 KiB je Change (abgeleitet), innerhalb des Bereichs
aus Abschnitt 3.3.

## 4. Ursache

**Code (gelesen).** `internal/application/usecase/retention/service.go`, `Run`:
`ReadChanges(ctx, ChangeQuery{Source: command.Source})` ohne Limit, danach eine
Schleife über alle Records mit `AllowsDeletion`; `internal/adapters/driven/postgresstorage/sqlexec/translate.go`,
`ReadChanges`: alle Zeilen werden in eine Liste gelesen und durch die
Domänen-Konstruktoren geführt; `internal/adapters/driven/postgresstorage/queries/queries.go`,
`SelectChanges`: beide Row Images (`old_data`, `new_data`), `LIMIT $6` mit `NULL`;
`internal/bootstrap/wiring.go`, `retentionInterval` = 10 s (der Kommentar dort nennt
„Jeder Takt liest alle Changes der Quelle“).

**Schalter 1 — die Bereinigung aus (Reihen E, F, G, I, J gegen A, B, C, H).** Ein
Image mit `retentionInterval` = 1000 h, sonst dasselbe: bei Blockgröße 1.000 und
schmalen Zeilen (Reihen E und J, `n` = 12) liegt die Spitze im Run bei 8,9 bis
10,6 MiB (`anon`) **und der Speicher im Nachlauf bei höchstens 8,0 bis 10,0 MiB
(`anon`, Maximum im 60-s-Fenster)**, obwohl `cdc.change` bis 3.000.000 Changes wächst
(Reihe J: `memory.peak` 13,3, 13,2 und 13,7 MiB bei 1.000.000, 2.000.000 und 3.000.000
Changes; Reihe E: 11,8 bis 29,0 MiB `memory.peak`, der Wert 29,0 MiB ist die
Datei-Seiten-Zugabe des ersten Containers, Abschnitt 3.1). Bei Blockgröße 100 (Reihe F)
8,4 bis 8,6 und 8,1 bis 8,3 MiB, bei Blockgröße 10.000 (Reihe G) 24,3 bis 24,4 und 16,3
bis 16,8 MiB, bei breiten Zeilen (Reihe I) 15,2 bis 15,7 und 8,7 bis 11,9 MiB: jeweils
flach über die Zahl der Changes. Dieselbe Messung mit Bereinigung: 1.082,7 MiB
`memory.peak` bei 1.000.000 Changes (Reihe B, Run 1). Die Bereinigung ist damit für die
Spitze nach dem Run **notwendig** (aus: flach) und bei ausreichend vielen Changes
**hinreichend** (an: die Werte der Abschnitte 3.3 und 3.5); eine Änderung, dieselbe
Messung, an fünf Reihen.

**Schalter 2 — die Zahl der Changes (Reihe D).** Leert der Lauf `cdc.change` bei
laufendem Feed, bleibt nach der nächsten Garbage Collection ein Heap von 1 MB
(Abschnitt 3.7): der Bestand des Heaps sind die gelesenen Changes.

**Zeitpunkt der Probe.** Die Spitze im Run (Abschnitt 3.2) und die Spitze nach dem Run
liegen bei 1.000.000 Changes um das Achtzigfache auseinander (`memory.peak` 13,5 MiB am
Ende des Runs gegen 1.082,7 MiB nach dem Nachlauf, Reihe B, Run 1, abgeleitet); die
Handbuch-Aussage, die Probe 20 s nach dem Run liege über der Spitze im Run, folgt
daraus (Bereinigungs-Takte im Nachlauf).

**Was ausgeschlossen ist.** Der Snapshot-Leser (Cursor-Block von `B` Zeilen), der Bau
des Row Image im Use Case, die Schreibtransaktion mit `AppendBlock` je Zeile und die
Garbage Collection erklären die Spitze nicht: die Spitze im Run bleibt bei 1.000.000
Zeilen bei 10,5 MiB (Reihe B, Run 1), mit `GOGC=25` sinkt die Spitze nach dem Run nur
um 11 bis 21 % (Reihe K, Abschnitt 3.6). Die eine ungeklärte Größe ist das Ausbleiben der
Bereinigungs-Takte ab 2.000.000 Changes (Abschnitt 3.3).

**Grad der Sicherheit.** Hoch für die Aussage „die Bereinigung erzeugt die Spitze nach
dem Run und ihre Abhängigkeit von der Zahl der Changes“ (Code gelesen, zwei Schalter,
Schalter 1 an fünf Reihen und drei Blockgrößen wiederholt). Mittel für die Zahl je
Change (1,0 bis 1,6 KiB): sie streut mit dem Zeitpunkt der Garbage Collection.
Niedrig für das Verhalten ab 2.000.000 Changes.

## 5. Ausgang

1. **Änderungs-Slice.** Die Ursache ist ein behebbarer Defekt im Code, nicht eine
   Grenze der Zeilenzahl des Backfills: `slice-retention-lauf-speicher-begrenzung`
   (Datei in `open/`). Er ist **nicht** Teil dieses Slice (Untersuchung,
   `ADR`-Trigger beim Architect).
2. **Benutzerhandbuch.** §9 „Grenzwerte“ trägt die gemessene Grenze samt Ursprung, eine
   Bemessungsregel für den Speicher des Feed-Containers und den Hinweis, dass
   `docker --memory` unterhalb des Bedarfs den Container beendet: ein Feed mit 64 MiB
   Grenze wird nach einem Backfill über 100.000 Zeilen vom Kernel beendet (Reihen L1
   und L2, Abschnitt 8); wie sich ein durch den Neustart-Mechanismus des Betreibers wieder
   gestarteter Feed verhält, ist nicht gemessen (erwartet: derselbe Takt liest die
   Changes erneut).
3. **Grundlinie und Messwerkzeug.** `tools/bench-backfill.sh` druckt die Grundlinie ohne
   Run und die Zeilen in `cdc.change`; `tools/bench-backfill-memory.sh` trennt die
   Spitze im Run von der nach dem Run.

## 6. Warn-Richtgröße

Der Wert im Code bleibt: 4.000.000 geschätzte Zeilen
(`internal/application/usecase/backfill/warn.go`). Bewertung gegen die Messung: die
Richtgröße folgt der **Kopierdauer** (Rate mal Toleranz), und diese Größe ist von der
Speicherlage unberührt; der Lauf N (Stufe 1.000.000, Median 7.715 Zeilen/s mal 600 s =
4.629.000, abgerundet 4.000.000, Ausgabezeile „Richtgröße (abgeleitet)“) ergibt
denselben Wert wie der Code, und die Kopierdauer der Stufe (129.620 ms) liegt weit
unter der Toleranz von 600 s; der Speicherbedarf nach einem Backfill dieser Größe liegt bei
3,9 bis 6,1 GiB für schmale Zeilen (4.000.000 × 1,03 bis 1,59 KiB, abgeleitet) und
10,1 bis 16,0 GiB für Zeilen von 1,3 KB (4.000.000 × 2,65 bis 4,19 KiB, abgeleitet). Eine Nachschärfung der Konstante beseitigte das Problem nicht: die
Warnung `warn_estimated_size` nennt keine Speichergröße, und eine Richtgröße in Zeilen
kann die Zeilenbreite nicht tragen. Deshalb (a) bleibt die Konstante bis zum
Änderungs-Slice, (b) führt das Handbuch die Bemessungsregel je Change statt einer
Zeilengrenze, und (c) ist die Neubemessung nach dem Trigger „Eine Messung liegt vor“
([`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)) ein
Liefer-Punkt des Änderungs-Slice: erst nach der Behebung sagt die Messung, welche
Größe die Richtgröße trägt.

## 7. Empfehlung zum Server-Release

**Reichweite.** Der Retention-Lauf mit unbegrenzter Lesung steht in allen drei
veröffentlichten Server-Versionen `v0.1.0` bis `v0.1.2` (Beleg: `git diff v0.1.2 HEAD --
internal/application/usecase/retention internal/application/port/outbound/changestore.go`
ist leer; `git show v0.1.0:internal/application/usecase/retention/service.go` und die
Fassungen von `v0.1.1` und `v0.1.2` tragen in Zeile 63 dieselbe Zeile `ReadChanges(ctx,
outbound.ChangeQuery{Source: command.Source})`, `retentionInterval = 10 * time.Second`
steht in jedem `internal/bootstrap/wiring.go` in Zeile 172; gelesen, nicht am Bild der
Version gemessen; Quelle des Befunds:
[`architect-verdict-retention-lauf-speicher-begrenzung`](architect-verdict-retention-lauf-speicher-begrenzung.md)
Verdikt 6). Der Defekt ist deshalb nicht auf den Backfill beschränkt: bei laufender
Erfassung hält das Mindestalter von 24 Stunden die Changes von 24 Stunden in
`cdc.change`, und derselbe Lauf liest sie in jedem Takt (abgeleitet, nicht gemessen).

**Kosten ohne Behebung.** Ein Server-Release mit dem Backfill **ohne** die Behebung
setzt die Betreiber dem Verhalten aus Abschnitt 3.3 aus: ein Backfill über 1.000.000
Zeilen hält den Feed-Container in jedem Takt der Bereinigung auf Spitzen von 1,1 bis
1,2 GiB (bei breiten Zeilen mehr), und ein Container-Limit unterhalb des Bedarfs beendet
ihn.

**Ausgang.** Der Nutzer hat entschieden: der Änderungs-Slice
`slice-retention-lauf-speicher-begrenzung` wird **zuerst** umgesetzt, danach folgt der
Server-Release `v0.2.0`; die Lösung steht in
[`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md)
(Projektion ohne Row Images, Seiten zu 10.000 Kandidaten) und im Architect-Verdikt
(Abschnitt „Anlass“). Die Reihenfolge steht als Vorbedingung des Server-Release `v0.2.0` im
Plan des Änderungs-Slice (Abschnitt 4). Ein Release vor der Umsetzung ist nicht vorgesehen.

## 8. Mutationen der Zusagen des Werkzeugs

Je Zusage eine Änderung an der **Eingabe** des geprüften Werkzeugs und das gesehene Rot
(die Exit-Codes sind beim Aufruf gelesen und stehen nicht in den gedruckten Zeilen):

| Zusage | mutierte Eingabe | gesehenes Rot |
|---|---|---|
| `BENCH_FEED_ENV` erreicht den Feed | Aufruf ohne `BENCH_FEED_ENV` (Reihe L2) gegen den Aufruf mit `GODEBUG=gctrace=1` (Reihe L1) | „im Run: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt)“ statt „GC-Läufe 168“ |
| `BENCH_FEED_DOCKER_ARGS` gilt für den Container | `--memory 64m` bei 100.000 Zeilen (Reihen L1 und L2) | Exit 1, „Feed-Container nicht mehr lesbar: Status exited, Exit 137, OOMKilled true“ |
| Der Zähler-Zugriff meldet einen nicht lesbaren cgroup-Pfad | Kopie des Skripts mit verändertem cgroup-Pfad (`system.slicex`, Reihe O) | Exit 1, „cgroup (…) oder /proc/… nicht lesbar“ |
| `BENCH_MEM_RESET` bestimmt die Zahl der Changes vor der Stufe | `BENCH_MEM_RESET=0` bei Stufen 10.000 und 20.000 (Reihe M) | „Zeilen in cdc.change vor dem Run: 10000“ in Stufe 20.000 statt 0 |
| `BENCH_MEM_EMPTY_AFTER` leert und misst weiter | Wert 1 (Reihe D) gegen 0 (Reihe A) | Zeile „cdc.change geleert, Nachlauf 120 s“ nur mit 1 |
| Die Grundlinie und die Zeilen in `cdc.change` in `bench-backfill.sh` folgen dem Inhalt der Ablage | der Inhalt von `cdc.change` (Eingabe der Zeile) ändert sich von Stufe zu Stufe derselben Reihe N | die Zeilen nennen 0, 30.000, 330.000 und 3.330.000 Changes und Ruhewerte von 6,1, 6,1, 41,7 und 275,7 MiB |
| Ein durch die Grenze beendeter Feed-Container beendet die Messung in `tools/bench-backfill.sh` mit dem Zustand des Containers (`feed_running_or_report`, Warteschleifen des Runs und der Live-Phase) | `BENCH_FEED_DOCKER_ARGS="--memory 64m"`, Stufe 100.000; vor der Änderung Reihe P1 (2 Runs), nach ihr Reihen P2 (3 Runs) und P3 (2 Runs) | vor: keine Zeile und kein Exit-Code, Abbruch von Hand nach mehr als 600 s, obwohl der Container beendet war (die Proben nennen 0,0 MiB); nach: Exit 1 mit „Feed-Container läuft nicht mehr, Run … nicht beendet (Run-Status queued): Status exited, Exit 137, OOMKilled true“ (P2) bzw. „…, Live-Phase, Run … nicht beendet (Run-Status unbekannt): …“ (P3) |
| `gc_summary` nennt `GODEBUG=gctrace=1` nur als ungesetzt, wenn es in `BENCH_FEED_ENV` fehlt | `BENCH_FEED_ENV` von `GODEBUG=gctrace=1` auf `GOGC=25`, `XGODEBUG=gctrace=1` und leer (Reihe R) | vor: „nicht gesetzt“ für alle fünf Werte; nach: „im Fenster (GODEBUG=gctrace=1 gesetzt)“ nur bei den zwei Werten mit `GODEBUG=gctrace=1` |
| Die Zeitreihen-Datei eines Runs bleibt nach einem Abbruch nicht zurück | `BENCH_MEM_RUN_TIMEOUT_S=1` (der Run endet ohne Abschluss, Reihen Q1 und Q2) | vor: eine Datei `tmp.…` in `TMPDIR` nach dem Lauf; nach: keine |
| `feed_mem_mib` liest den Container, den die Messung startet | keine Mutation gefahren | ohne Antwort: bei einem nicht lesbaren Container liefert die Funktion 0,0 statt eines Fehlers (vorhandene Eigenschaft, nicht Teil dieses Slice; die Warteschleifen der Runs prüfen den Container seither selbst) |

## 9. Aufräum-Schritte und Grenzen

- **Aufräumen:** jeder Lauf endet mit `docker rm -fv` der beiden Container und
  `docker network rm` (Falle `bench::cleanup`); kein `docker volume prune`, kein
  `docker system prune`. Gezählte freie Volumes: vor der ersten Messung 34, nach der
  letzten siehe Abschnitt 10.
- **Selbst gebaute Images** `pg-change-feed:bench-noret`, `…-b100`, `…-b10000`
  (Varianten der Schalter) verbleiben lokal und sind nicht committet.
- **Grenzen:** eine Umgebung (ein Host, PostgreSQL 18, Standardeinstellungen); `n` = 2
  bis 3 je Bedingung, 1.000.000 Changes zweimal, 2.000.000 Changes zweimal; die
  Zeilenbreite hat zwei Werte (73 und 1.273 Bytes); die Zahl je Change streut; die
  Bereinigungs-Takte ab 2.000.000 Changes sind beobachtet und nicht erklärt; der
  Speicher der PostgreSQL-Instanz ist nicht gemessen; Live-Erfassung (statt Backfill)
  als Quelle der Changes ist aus dem Code gelesen, nicht gemessen.

## 10. Bestand nach den Läufen

Nach dem letzten Lauf (Reihe N): `docker volume ls -q -f dangling=true` zählt **34**
(vor dem ersten Lauf ebenfalls 34, keine vom Runner erzeugte Leiche); kein Container
und kein Netz `pgc-bench-*` vorhanden; freier Platz der Wurzel-Platte 88 GB vor und nach
den Läufen (68 % belegt vor, 69 % nach den Läufen). Der Arbeitsbaum trägt nach den Läufen weder
`tools/schema/plan.yaml` noch `tools/schema/down.sql` als Änderung (jeweils per
`git checkout` zurückgenommen).
