# Gedruckte Zeilen der Messreihen

Zu [`messbericht-slice-backfill-speicher-untersuchung`](messbericht-slice-backfill-speicher-untersuchung.md):
die Ausgabe der Läufe, unverändert bis auf zwei Auslassungen — das Präfix
`bench-backfill-memory[<Lauf>]: ` (bzw. `bench-backfill[<Lauf>]: `) am Zeilenanfang steht
einmal in der Überschrift der Reihe, und die Hinweiszeile von `psql` zum Rollout
(„NOTICE: constraint … does not exist, skipping“) fehlt. Die Reihen sind im Messbericht
(Abschnitt 2) benannt; die Zahlen dort stammen aus diesen Zeilen. Die Zeilen tragen
die Zähler `memory.current`, `anon`, `file`, `kernel` (cgroup v2 des Containers),
`RssAnon` und `VmHWM` (aus `/proc/<pid>/status`) sowie `memory.peak` (Höchstwert
seit dem Start des Containers), alle in MiB bzw. KiB wie gedruckt.

## Reihe A — Lauf `20260925T123600Z`

Default-Image, schmal, Stufen 10.000, 100.000, 200.000, je 3 Runs.

```text
Host — Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU, 33362599936 Byte RAM, PostgreSQL-Image postgres:18-alpine, Feed-Image-ID sha256:969b7fad86909255d47fdf1bd5ada33edee15a171d6cc8b8a9849edb3965cf4e
Stufen 10000 100000 200000, 3 Runs je Stufe, Zeilenbreite narrow, Feed-Umgebung 'GODEBUG=gctrace=1', Feed-Docker-Argumente '', Feed-Image ghcr.io/pt9912/pg-change-feed:dev
Tabelle bench_mem_narrow_10000 — 10000 Zeilen, mittlere Zeilenbreite 73 B (pg_column_size), 64 B (Textform)
Tabelle bench_mem_narrow_100000 — 100000 Zeilen, mittlere Zeilenbreite 73 B (pg_column_size), 67 B (Textform)
Tabelle bench_mem_narrow_200000 — 200000 Zeilen, mittlere Zeilenbreite 74 B (pg_column_size), 69 B (Textform)
Stufe 10000 narrow, Run 1/3 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.2 MiB zu Beginn, Maximum 6.5, Minimum 6.2, Ende 6.2; anon Beginn 5.2, Maximum 5.2, Minimum 5.2, Ende 5.2 MiB
Stufe 10000 narrow, Run 1/3 — Kopierdauer 1307 ms, 7651 Zeilen/s; Proben 6; memory.current Beginn 6.2, bei 25% 10.6, 50% 10.1, 75% 9.8, Ende 9.8, Spitze 10.6 MiB; anon Beginn 5.2, bei 25% 8.5, 50% 8.8, 75% 8.7, Ende 8.7, Spitze 8.9 MiB
Stufe 10000 narrow, Run 1/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 9.7 MiB (anon 8.7, file 0.2, kernel 0.9), Prozess RssAnon 8.7 MiB, docker stats 9.7 MiB, VmHWM 24032 KiB, memory.peak seit Start des Feeds 12.3 MiB; im Run: GC-Läufe 17, größter Heap zu Beginn 5 MB, größtes Ziel 5 MB, letzter Heap nach GC 2 MB
Stufe 10000 narrow, Run 1/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 9.7 MiB zu Beginn, Maximum 22.4, Minimum 9.7, Ende 18.6; anon Beginn 8.7, Maximum 21.2, Minimum 8.7, Ende 17.4 MiB
Stufe 10000 narrow, Run 1/3 — Zeilen in cdc.change nach dem Nachlauf: 10000, VmHWM 35244 KiB, memory.peak seit Start des Feeds 23.9 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 23, größter Heap zu Beginn 14 MB, größtes Ziel 15 MB, letzter Heap nach GC 6 MB
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 201: active WalSenderWaitForWal) — Sitzung beendet
Stufe 10000 narrow, Run 2/3 — Zeilen in cdc.change vor dem Run: 10000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 18.3 MiB zu Beginn, Maximum 19.3, Minimum 18.3, Ende 19.3; anon Beginn 17.4, Maximum 18.3, Minimum 17.4, Ende 18.3 MiB
Stufe 10000 narrow, Run 2/3 — Kopierdauer 1050 ms, 9524 Zeilen/s; Proben 5; memory.current Beginn 19.9, bei 25% 12.7, 50% 11.6, 75% 12.2, Ende 12.2, Spitze 19.9 MiB; anon Beginn 18.3, bei 25% 11.5, 50% 10.1, 75% 10.9, Ende 10.9, Spitze 18.3 MiB
Stufe 10000 narrow, Run 2/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 12.4 MiB (anon 11.5, file 0.0, kernel 1.0), Prozess RssAnon 11.5 MiB, docker stats 12.4 MiB, VmHWM 35984 KiB, memory.peak seit Start des Feeds 25.0 MiB; im Run: GC-Läufe 17, größter Heap zu Beginn 13 MB, größtes Ziel 13 MB, letzter Heap nach GC 2 MB
Stufe 10000 narrow, Run 2/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 12.4 MiB zu Beginn, Maximum 38.0, Minimum 12.4, Ende 24.0; anon Beginn 11.5, Maximum 35.0, Minimum 11.5, Ende 23.1 MiB
Stufe 10000 narrow, Run 2/3 — Zeilen in cdc.change nach dem Nachlauf: 20000, VmHWM 51220 KiB, memory.peak seit Start des Feeds 40.4 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 24, größter Heap zu Beginn 29 MB, größtes Ziel 29 MB, letzter Heap nach GC 10 MB
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 814: active WalSenderWaitForWal) — Sitzung beendet
Stufe 10000 narrow, Run 3/3 — Zeilen in cdc.change vor dem Run: 20000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 24.9 MiB zu Beginn, Maximum 27.3, Minimum 24.9, Ende 27.0; anon Beginn 23.7, Maximum 25.8, Minimum 23.7, Ende 25.8 MiB
Stufe 10000 narrow, Run 3/3 — Kopierdauer 1049 ms, 9533 Zeilen/s; Proben 5; memory.current Beginn 27.3, bei 25% 10.6, 50% 11.0, 75% 10.6, Ende 10.6, Spitze 27.3 MiB; anon Beginn 25.8, bei 25% 9.4, 50% 9.1, 75% 9.5, Ende 9.5, Spitze 25.8 MiB
Stufe 10000 narrow, Run 3/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 10.6 MiB (anon 9.5, file 0.2, kernel 1.0), Prozess RssAnon 9.5 MiB, docker stats 10.6 MiB, VmHWM 44664 KiB, memory.peak seit Start des Feeds 34.8 MiB; im Run: GC-Läufe 16, größter Heap zu Beginn 20 MB, größtes Ziel 21 MB, letzter Heap nach GC 2 MB
Stufe 10000 narrow, Run 3/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 10.6 MiB zu Beginn, Maximum 40.1, Minimum 10.6, Ende 36.6; anon Beginn 9.5, Maximum 38.2, Minimum 9.5, Ende 35.3 MiB
Stufe 10000 narrow, Run 3/3 — Zeilen in cdc.change nach dem Nachlauf: 30000, VmHWM 61576 KiB, memory.peak seit Start des Feeds 50.2 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 29, größter Heap zu Beginn 38 MB, größtes Ziel 38 MB, letzter Heap nach GC 15 MB
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 1420: active WalSenderWaitForWal) — Sitzung beendet
Stufe 100000 narrow, Run 1/3 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 5.9 MiB zu Beginn, Maximum 5.9, Minimum 5.8, Ende 5.9; anon Beginn 5.0, Maximum 5.0, Minimum 5.0, Ende 5.0 MiB
Stufe 100000 narrow, Run 1/3 — Kopierdauer 10051 ms, 9949 Zeilen/s; Proben 39; memory.current Beginn 6.6, bei 25% 10.4, 50% 10.4, 75% 10.2, Ende 10.1, Spitze 11.5 MiB; anon Beginn 5.0, bei 25% 8.9, 50% 9.4, 75% 8.8, Ende 8.8, Spitze 9.8 MiB
Stufe 100000 narrow, Run 1/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 9.8 MiB (anon 8.8, file 0.0, kernel 1.0), Prozess RssAnon 8.8 MiB, docker stats 9.8 MiB, VmHWM 24344 KiB, memory.peak seit Start des Feeds 12.8 MiB; im Run: GC-Läufe 170, größter Heap zu Beginn 5 MB, größtes Ziel 5 MB, letzter Heap nach GC 2 MB
Stufe 100000 narrow, Run 1/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 9.8 MiB zu Beginn, Maximum 114.2, Minimum 9.8, Ende 114.2; anon Beginn 8.8, Maximum 112.8, Minimum 8.8, Ende 112.8 MiB
Stufe 100000 narrow, Run 1/3 — Zeilen in cdc.change nach dem Nachlauf: 100000, VmHWM 144308 KiB, memory.peak seit Start des Feeds 133.8 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 52, größter Heap zu Beginn 104 MB, größtes Ziel 104 MB, letzter Heap nach GC 49 MB
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 2039: active WalSenderWaitForWal) — Sitzung beendet
Stufe 100000 narrow, Run 2/3 — Zeilen in cdc.change vor dem Run: 100000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 98.0 MiB zu Beginn, Maximum 100.1, Minimum 98.0, Ende 100.0; anon Beginn 96.7, Maximum 98.7, Minimum 96.7, Ende 98.7 MiB
Stufe 100000 narrow, Run 2/3 — Kopierdauer 15445 ms, 6475 Zeilen/s; Proben 58; memory.current Beginn 100.3, bei 25% 12.9, 50% 109.6, 75% 14.2, Ende 12.6, Spitze 143.5 MiB; anon Beginn 98.7, bei 25% 11.6, 50% 108.2, 75% 12.3, Ende 11.2, Spitze 142.1 MiB
Stufe 100000 narrow, Run 2/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 12.6 MiB (anon 11.2, file 0.0, kernel 1.4), Prozess RssAnon 11.2 MiB, docker stats 12.6 MiB, VmHWM 158360 KiB, memory.peak seit Start des Feeds 145.1 MiB; im Run: GC-Läufe 130, größter Heap zu Beginn 113 MB, größtes Ziel 115 MB, letzter Heap nach GC 2 MB
Stufe 100000 narrow, Run 2/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 12.6 MiB zu Beginn, Maximum 242.8, Minimum 12.6, Ende 200.6; anon Beginn 11.2, Maximum 240.6, Minimum 11.2, Ende 198.8 MiB
Stufe 100000 narrow, Run 2/3 — Zeilen in cdc.change nach dem Nachlauf: 200000, VmHWM 260020 KiB, memory.peak seit Start des Feeds 245.0 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 53, größter Heap zu Beginn 235 MB, größtes Ziel 239 MB, letzter Heap nach GC 108 MB
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 2883: active WalSenderWaitForWal) — Sitzung beendet
Stufe 100000 narrow, Run 3/3 — Zeilen in cdc.change vor dem Run: 200000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 104.0 MiB zu Beginn, Maximum 341.3, Minimum 104.0, Ende 340.9; anon Beginn 102.1, Maximum 338.7, Minimum 102.1, Ende 338.7 MiB
Stufe 100000 narrow, Run 3/3 — Kopierdauer 11950 ms, 8368 Zeilen/s; Proben 46; memory.current Beginn 341.0, bei 25% 106.7, 50% 16.5, 75% 101.0, Ende 211.8, Spitze 341.5 MiB; anon Beginn 338.7, bei 25% 104.6, 50% 13.9, 75% 98.6, Ende 209.8, Spitze 339.3 MiB
Stufe 100000 narrow, Run 3/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 211.8 MiB (anon 209.8, file 0.2, kernel 1.8), Prozess RssAnon 209.8 MiB, docker stats 211.8 MiB, VmHWM 360716 KiB, memory.peak seit Start des Feeds 343.7 MiB; im Run: GC-Läufe 97, größter Heap zu Beginn 209 MB, größtes Ziel 213 MB, letzter Heap nach GC 100 MB
Stufe 100000 narrow, Run 3/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 211.8 MiB zu Beginn, Maximum 423.7, Minimum 101.7, Ende 263.4; anon Beginn 209.8, Maximum 421.4, Minimum 99.8, Ende 261.2 MiB
Stufe 100000 narrow, Run 3/3 — Zeilen in cdc.change nach dem Nachlauf: 300000, VmHWM 467688 KiB, memory.peak seit Start des Feeds 449.2 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 30, größter Heap zu Beginn 409 MB, größtes Ziel 410 MB, letzter Heap nach GC 120 MB
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 3876: active WalSenderWaitForWal) — Sitzung beendet
Stufe 200000 narrow, Run 1/3 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.0 MiB zu Beginn, Maximum 6.0, Minimum 6.0, Ende 6.0; anon Beginn 5.2, Maximum 5.2, Minimum 5.2, Ende 5.2 MiB
Stufe 200000 narrow, Run 1/3 — Kopierdauer 19008 ms, 10522 Zeilen/s; Proben 74; memory.current Beginn 6.3, bei 25% 10.6, 50% 11.2, 75% 10.3, Ende 10.4, Spitze 12.0 MiB; anon Beginn 5.2, bei 25% 9.3, 50% 9.2, 75% 9.1, Ende 9.4, Spitze 10.1 MiB
Stufe 200000 narrow, Run 1/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 10.4 MiB (anon 9.4, file 0.0, kernel 1.0), Prozess RssAnon 9.4 MiB, docker stats 10.4 MiB, VmHWM 24860 KiB, memory.peak seit Start des Feeds 12.7 MiB; im Run: GC-Läufe 335, größter Heap zu Beginn 5 MB, größtes Ziel 5 MB, letzter Heap nach GC 2 MB
Stufe 200000 narrow, Run 1/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 10.4 MiB zu Beginn, Maximum 258.0, Minimum 10.4, Ende 193.3; anon Beginn 9.4, Maximum 253.2, Minimum 9.4, Ende 191.6 MiB
Stufe 200000 narrow, Run 1/3 — Zeilen in cdc.change nach dem Nachlauf: 200000, VmHWM 288968 KiB, memory.peak seit Start des Feeds 273.9 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 52, größter Heap zu Beginn 236 MB, größtes Ziel 240 MB, letzter Heap nach GC 108 MB
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 4803: active WalSenderWaitForWal) — Sitzung beendet
Stufe 200000 narrow, Run 2/3 — Zeilen in cdc.change vor dem Run: 200000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 220.3 MiB zu Beginn, Maximum 220.3, Minimum 220.0, Ende 220.0; anon Beginn 218.6, Maximum 218.6, Minimum 217.4, Ende 218.2 MiB
Stufe 200000 narrow, Run 2/3 — Kopierdauer 21979 ms, 9100 Zeilen/s; Proben 85; memory.current Beginn 220.3, bei 25% 17.9, 50% 19.0, 75% 189.0, Ende 17.0, Spitze 283.4 MiB; anon Beginn 218.2, bei 25% 15.8, 50% 16.8, 75% 183.9, Ende 15.2, Spitze 281.4 MiB
Stufe 200000 narrow, Run 2/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 17.0 MiB (anon 15.2, file 0.0, kernel 1.8), Prozess RssAnon 15.2 MiB, docker stats 17.0 MiB, VmHWM 361268 KiB, memory.peak seit Start des Feeds 344.2 MiB; im Run: GC-Läufe 293, größter Heap zu Beginn 236 MB, größtes Ziel 240 MB, letzter Heap nach GC 2 MB
Stufe 200000 narrow, Run 2/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 17.0 MiB zu Beginn, Maximum 415.4, Minimum 17.0, Ende 278.9; anon Beginn 15.2, Maximum 412.8, Minimum 15.2, Ende 276.6 MiB
Stufe 200000 narrow, Run 2/3 — Zeilen in cdc.change nach dem Nachlauf: 400000, VmHWM 488380 KiB, memory.peak seit Start des Feeds 468.6 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 41, größter Heap zu Beginn 466 MB, größtes Ziel 471 MB, letzter Heap nach GC 100 MB
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 5907: active WalSenderWaitForWal) — Sitzung beendet
Stufe 200000 narrow, Run 3/3 — Zeilen in cdc.change vor dem Run: 400000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 396.8 MiB zu Beginn, Maximum 396.8, Minimum 218.5, Ende 218.5; anon Beginn 394.5, Maximum 394.6, Minimum 216.3, Ende 216.3 MiB
Stufe 200000 narrow, Run 3/3 — Kopierdauer 21027 ms, 9512 Zeilen/s; Proben 80; memory.current Beginn 218.8, bei 25% 277.2, 50% 489.5, 75% 23.0, Ende 278.3, Spitze 605.7 MiB; anon Beginn 216.3, bei 25% 274.9, 50% 486.6, 75% 20.2, Ende 275.5, Spitze 602.5 MiB
Stufe 200000 narrow, Run 3/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 278.0 MiB (anon 275.5, file 0.0, kernel 2.4), Prozess RssAnon 275.6 MiB, docker stats 278.0 MiB, VmHWM 630020 KiB, memory.peak seit Start des Feeds 605.9 MiB; im Run: GC-Läufe 91, größter Heap zu Beginn 463 MB, größtes Ziel 469 MB, letzter Heap nach GC 120 MB
Stufe 200000 narrow, Run 3/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 278.0 MiB zu Beginn, Maximum 511.6, Minimum 235.9, Ende 511.4; anon Beginn 275.6, Maximum 508.9, Minimum 233.4, Ende 508.9 MiB
Stufe 200000 narrow, Run 3/3 — Zeilen in cdc.change nach dem Nachlauf: 600000, VmHWM 630020 KiB, memory.peak seit Start des Feeds 605.9 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 9, größter Heap zu Beginn 492 MB, größtes Ziel 492 MB, letzter Heap nach GC 315 MB
Ende — Lauf 20260925T123600Z
```

## Reihe B — Lauf `20260925T130048Z`

Default-Image, schmal, 1.000.000, 3 Runs, `--memory 6g`.

```text
Host — Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU, 33362599936 Byte RAM, PostgreSQL-Image postgres:18-alpine, Feed-Image-ID sha256:969b7fad86909255d47fdf1bd5ada33edee15a171d6cc8b8a9849edb3965cf4e
Stufen 1000000, 3 Runs je Stufe, Zeilenbreite narrow, Feed-Umgebung 'GODEBUG=gctrace=1', Feed-Docker-Argumente '--memory 6g', Feed-Image ghcr.io/pt9912/pg-change-feed:dev
Tabelle bench_mem_narrow_1000000 — 1000000 Zeilen, mittlere Zeilenbreite 74 B (pg_column_size), 70 B (Textform)
Stufe 1000000 narrow, Run 1/3 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.0 MiB zu Beginn, Maximum 6.0, Minimum 6.0, Ende 6.0; anon Beginn 5.1, Maximum 5.1, Minimum 5.1, Ende 5.1 MiB
Stufe 1000000 narrow, Run 1/3 — Kopierdauer 95133 ms, 10512 Zeilen/s; Proben 364; memory.current Beginn 6.3, bei 25% 10.2, 50% 10.0, 75% 11.2, Ende 10.6, Spitze 12.7 MiB; anon Beginn 5.1, bei 25% 9.0, 50% 8.7, 75% 9.2, Ende 8.8, Spitze 10.5 MiB
Stufe 1000000 narrow, Run 1/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 9.9 MiB (anon 8.8, file 0.1, kernel 1.0), Prozess RssAnon 8.8 MiB, docker stats 9.9 MiB, VmHWM 24996 KiB, memory.peak seit Start des Feeds 13.5 MiB; im Run: GC-Läufe 1661, größter Heap zu Beginn 5 MB, größtes Ziel 6 MB, letzter Heap nach GC 2 MB
Stufe 1000000 narrow, Run 1/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 9.9 MiB zu Beginn, Maximum 1080.1, Minimum 9.9, Ende 428.3; anon Beginn 8.8, Maximum 1076.6, Minimum 8.8, Ende 423.9 MiB
Stufe 1000000 narrow, Run 1/3 — Zeilen in cdc.change nach dem Nachlauf: 1000000, VmHWM 1115812 KiB, memory.peak seit Start des Feeds 1082.7 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 68, größter Heap zu Beginn 1021 MB, größtes Ziel 1021 MB, letzter Heap nach GC 369 MB
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 155: active WalSenderWaitForWal) — Sitzung beendet
Stufe 1000000 narrow, Run 2/3 — Zeilen in cdc.change vor dem Run: 1000000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 1203.0 MiB zu Beginn, Maximum 1203.5, Minimum 476.1, Ende 1202.8; anon Beginn 1198.1, Maximum 1198.2, Minimum 470.2, Ende 1197.9 MiB
Stufe 1000000 narrow, Run 2/3 — Kopierdauer 105023 ms, 9522 Zeilen/s; Proben 399; memory.current Beginn 1202.8, bei 25% 1171.8, 50% 41.6, 75% 1090.2, Ende 43.3, Spitze 1466.7 MiB; anon Beginn 1197.9, bei 25% 1165.9, 50% 36.1, 75% 1085.0, Ende 38.0, Spitze 1460.8 MiB
Stufe 1000000 narrow, Run 2/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 43.0 MiB (anon 38.0, file 0.1, kernel 4.9), Prozess RssAnon 38.0 MiB, docker stats 43.0 MiB, VmHWM 1580524 KiB, memory.peak seit Start des Feeds 1538.7 MiB; im Run: GC-Läufe 414, größter Heap zu Beginn 1236 MB, größtes Ziel 1236 MB, letzter Heap nach GC 2 MB
Stufe 1000000 narrow, Run 2/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 43.0 MiB zu Beginn, Maximum 2251.1, Minimum 43.0, Ende 1973.0; anon Beginn 38.0, Maximum 2244.4, Minimum 38.0, Ende 1965.9 MiB
Stufe 1000000 narrow, Run 2/3 — Zeilen in cdc.change nach dem Nachlauf: 2000000, VmHWM 2327392 KiB, memory.peak seit Start des Feeds 2269.2 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 35, größter Heap zu Beginn 2282 MB, größtes Ziel 2282 MB, letzter Heap nach GC 937 MB
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 3287: active WalSenderWaitForWal) — Sitzung beendet
Stufe 1000000 narrow, Run 3/3 — Zeilen in cdc.change vor dem Run: 2000000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 2439.3 MiB zu Beginn, Maximum 2439.3, Minimum 2439.3, Ende 2439.3; anon Beginn 2431.8, Maximum 2431.8, Minimum 2431.8, Ende 2431.8 MiB
Stufe 1000000 narrow, Run 3/3 — Kopierdauer 102209 ms, 9784 Zeilen/s; Proben 388; memory.current Beginn 2439.3, bei 25% 3063.2, 50% 69.9, 75% 70.4, Ende 70.4, Spitze 3096.6 MiB; anon Beginn 2431.8, bei 25% 3055.7, 50% 62.1, 75% 62.1, Ende 62.0, Spitze 3088.1 MiB
Stufe 1000000 narrow, Run 3/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 69.6 MiB (anon 62.0, file 0.1, kernel 7.6), Prozess RssAnon 62.0 MiB, docker stats 69.6 MiB, VmHWM 3175496 KiB, memory.peak seit Start des Feeds 3096.6 MiB; im Run: GC-Läufe 1213, größter Heap zu Beginn 2205 MB, größtes Ziel 2220 MB, letzter Heap nach GC 2 MB
Stufe 1000000 narrow, Run 3/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 69.6 MiB zu Beginn, Maximum 69.6, Minimum 69.6, Ende 69.6; anon Beginn 62.0, Maximum 62.0, Minimum 62.0, Ende 62.0 MiB
Stufe 1000000 narrow, Run 3/3 — Zeilen in cdc.change nach dem Nachlauf: 3000000, VmHWM 3175496 KiB, memory.peak seit Start des Feeds 3096.6 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt)
Ende — Lauf 20260925T130048Z
```

## Reihe C — Lauf `20260925T131359Z`

Wie B mit `--memory 8g`, mit Zählung der Bereinigungs-Takte.

```text
Host — Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU, 33362599936 Byte RAM, PostgreSQL-Image postgres:18-alpine, Feed-Image-ID sha256:969b7fad86909255d47fdf1bd5ada33edee15a171d6cc8b8a9849edb3965cf4e
Stufen 1000000, 3 Runs je Stufe, Zeilenbreite narrow, Feed-Umgebung 'GODEBUG=gctrace=1', Feed-Docker-Argumente '--memory 8g', Feed-Image ghcr.io/pt9912/pg-change-feed:dev
Tabelle bench_mem_narrow_1000000 — 1000000 Zeilen, mittlere Zeilenbreite 74 B (pg_column_size), 70 B (Textform)
Stufe 1000000 narrow, Run 1/3 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.1 MiB zu Beginn, Maximum 6.2, Minimum 6.1, Ende 6.2; anon Beginn 5.4, Maximum 5.4, Minimum 5.4, Ende 5.4 MiB
Stufe 1000000 narrow, Run 1/3 — Kopierdauer 118965 ms, 8406 Zeilen/s; Proben 444; memory.current Beginn 6.7, bei 25% 10.2, 50% 10.3, 75% 11.1, Ende 10.1, Spitze 12.3 MiB; anon Beginn 5.4, bei 25% 8.9, 50% 9.3, 75% 9.7, Ende 9.6, Spitze 10.1 MiB
Stufe 1000000 narrow, Run 1/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 10.1 MiB (anon 9.6, file 0.0, kernel 1.0), Prozess RssAnon 9.2 MiB, docker stats 10.2 MiB, VmHWM 24860 KiB, memory.peak seit Start des Feeds 13.3 MiB; im Run: GC-Läufe 1662, größter Heap zu Beginn 5 MB, größtes Ziel 6 MB, letzter Heap nach GC 2 MB
Stufe 1000000 narrow, Run 1/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 10.1 MiB zu Beginn, Maximum 1191.4, Minimum 10.1, Ende 1023.5; anon Beginn 9.2, Maximum 1185.1, Minimum 9.2, Ende 1019.1 MiB
Stufe 1000000 narrow, Run 1/3 — Zeilen in cdc.change nach dem Nachlauf: 1000000, VmHWM 1310740 KiB, memory.peak seit Start des Feeds 1273.5 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 71, größter Heap zu Beginn 1252 MB, größtes Ziel 1252 MB, letzter Heap nach GC 482 MB; seit dem Start des Feeds: Bereinigung gelaufen 25, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 162: active WalSenderWaitForWal) — Sitzung beendet
Stufe 1000000 narrow, Run 2/3 — Zeilen in cdc.change vor dem Run: 1000000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 1000.4 MiB zu Beginn, Maximum 1108.5, Minimum 458.0, Ende 1108.5; anon Beginn 996.8, Maximum 1105.0, Minimum 451.8, Ende 1105.0 MiB
Stufe 1000000 narrow, Run 2/3 — Kopierdauer 114202 ms, 8756 Zeilen/s; Proben 431; memory.current Beginn 1108.6, bei 25% 39.6, 50% 1312.5, 75% 564.2, Ende 45.6, Spitze 1525.4 MiB; anon Beginn 1105.0, bei 25% 35.4, 50% 1307.6, 75% 556.9, Ende 40.1, Spitze 1519.8 MiB
Stufe 1000000 narrow, Run 2/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 44.9 MiB (anon 40.1, file 0.0, kernel 4.9), Prozess RssAnon 40.1 MiB, docker stats 44.9 MiB, VmHWM 1701660 KiB, memory.peak seit Start des Feeds 1656.0 MiB; im Run: GC-Läufe 707, größter Heap zu Beginn 1555 MB, größtes Ziel 1555 MB, letzter Heap nach GC 2 MB
Stufe 1000000 narrow, Run 2/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 44.9 MiB zu Beginn, Maximum 2771.1, Minimum 44.9, Ende 2416.3; anon Beginn 40.1, Maximum 2761.9, Minimum 40.1, Ende 2408.7 MiB
Stufe 1000000 narrow, Run 2/3 — Zeilen in cdc.change nach dem Nachlauf: 2000000, VmHWM 3134492 KiB, memory.peak seit Start des Feeds 3058.0 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 54, größter Heap zu Beginn 2283 MB, größtes Ziel 2283 MB, letzter Heap nach GC 1128 MB; seit dem Start des Feeds: Bereinigung gelaufen 23, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 3874: active WalSenderWaitForWal) — Sitzung beendet
Stufe 1000000 narrow, Run 3/3 — Zeilen in cdc.change vor dem Run: 2000000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 2453.7 MiB zu Beginn, Maximum 2453.7, Minimum 2453.7, Ende 2453.7; anon Beginn 2446.8, Maximum 2446.8, Minimum 2446.8, Ende 2446.8 MiB
Stufe 1000000 narrow, Run 3/3 — Kopierdauer 118836 ms, 8415 Zeilen/s; Proben 446; memory.current Beginn 2454.5, bei 25% 60.3, 50% 61.1, 75% 60.9, Ende 62.2, Spitze 2711.0 MiB; anon Beginn 2446.8, bei 25% 53.1, 50% 54.0, 75% 53.8, Ende 54.8, Spitze 2703.9 MiB
Stufe 1000000 narrow, Run 3/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 61.7 MiB (anon 54.8, file 0.0, kernel 6.9), Prozess RssAnon 54.8 MiB, docker stats 61.7 MiB, VmHWM 2820796 KiB, memory.peak seit Start des Feeds 2751.7 MiB; im Run: GC-Läufe 1531, größter Heap zu Beginn 2395 MB, größtes Ziel 2410 MB, letzter Heap nach GC 2 MB
Stufe 1000000 narrow, Run 3/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 61.7 MiB zu Beginn, Maximum 61.7, Minimum 61.7, Ende 61.7; anon Beginn 54.8, Maximum 54.8, Minimum 54.8, Ende 54.8 MiB
Stufe 1000000 narrow, Run 3/3 — Zeilen in cdc.change nach dem Nachlauf: 3000000, VmHWM 2820796 KiB, memory.peak seit Start des Feeds 2751.7 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt); seit dem Start des Feeds: Bereinigung gelaufen 5, fehlgeschlagen 0
Ende — Lauf 20260925T131359Z
```

## Reihe D — Lauf `20260925T132801Z`

Default-Image, schmal, 100.000, 3 Runs, danach `cdc.change` leeren und 120 s messen.

```text
Host — Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU, 33362599936 Byte RAM, PostgreSQL-Image postgres:18-alpine, Feed-Image-ID sha256:969b7fad86909255d47fdf1bd5ada33edee15a171d6cc8b8a9849edb3965cf4e
Stufen 100000, 3 Runs je Stufe, Zeilenbreite narrow, Feed-Umgebung 'GODEBUG=gctrace=1', Feed-Docker-Argumente '', Feed-Image ghcr.io/pt9912/pg-change-feed:dev
Tabelle bench_mem_narrow_100000 — 100000 Zeilen, mittlere Zeilenbreite 73 B (pg_column_size), 67 B (Textform)
Stufe 100000 narrow, Run 1/3 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.2 MiB zu Beginn, Maximum 6.2, Minimum 6.2, Ende 6.2; anon Beginn 5.4, Maximum 5.4, Minimum 5.4, Ende 5.4 MiB
Stufe 100000 narrow, Run 1/3 — Kopierdauer 12125 ms, 8247 Zeilen/s; Proben 46; memory.current Beginn 7.0, bei 25% 10.2, 50% 10.2, 75% 10.6, Ende 10.5, Spitze 11.5 MiB; anon Beginn 5.4, bei 25% 9.2, 50% 9.2, 75% 9.4, Ende 9.3, Spitze 9.9 MiB
Stufe 100000 narrow, Run 1/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 10.3 MiB (anon 9.3, file 0.0, kernel 1.0), Prozess RssAnon 9.3 MiB, docker stats 10.2 MiB, VmHWM 24712 KiB, memory.peak seit Start des Feeds 12.3 MiB; im Run: GC-Läufe 168, größter Heap zu Beginn 5 MB, größtes Ziel 5 MB, letzter Heap nach GC 2 MB
Stufe 100000 narrow, Run 1/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 10.3 MiB zu Beginn, Maximum 142.3, Minimum 10.3, Ende 98.3; anon Beginn 9.3, Maximum 140.4, Minimum 9.3, Ende 96.9 MiB
Stufe 100000 narrow, Run 1/3 — Zeilen in cdc.change nach dem Nachlauf: 100000, VmHWM 165676 KiB, memory.peak seit Start des Feeds 151.3 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 34, größter Heap zu Beginn 135 MB, größtes Ziel 135 MB, letzter Heap nach GC 51 MB; seit dem Start des Feeds: Bereinigung gelaufen 14, fehlgeschlagen 0
Stufe 100000 narrow, Run 1/3 — cdc.change geleert, Nachlauf 120 s: Nachlauf nach der Leerung: memory.current 98.3 MiB zu Beginn, Maximum 98.4, Minimum 55.4, Ende 55.4; anon Beginn 96.9, Maximum 97.0, Minimum 54.0, Ende 54.0 MiB; im Fenster: GC-Läufe 1, größter Heap zu Beginn 64 MB, größtes Ziel 104 MB, letzter Heap nach GC 1 MB
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 158: active WalSenderWaitForWal) — Sitzung beendet
Stufe 100000 narrow, Run 2/3 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 5.8 MiB zu Beginn, Maximum 5.8, Minimum 5.8, Ende 5.8; anon Beginn 5.0, Maximum 5.0, Minimum 5.0, Ende 5.0 MiB
Stufe 100000 narrow, Run 2/3 — Kopierdauer 12906 ms, 7748 Zeilen/s; Proben 49; memory.current Beginn 5.9, bei 25% 10.1, 50% 10.0, 75% 10.2, Ende 9.9, Spitze 11.6 MiB; anon Beginn 5.0, bei 25% 8.9, 50% 9.0, 75% 8.8, Ende 8.9, Spitze 9.6 MiB
Stufe 100000 narrow, Run 2/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 9.9 MiB (anon 8.9, file 0.0, kernel 1.0), Prozess RssAnon 9.0 MiB, docker stats 9.9 MiB, VmHWM 24396 KiB, memory.peak seit Start des Feeds 13.1 MiB; im Run: GC-Läufe 168, größter Heap zu Beginn 5 MB, größtes Ziel 6 MB, letzter Heap nach GC 2 MB
Stufe 100000 narrow, Run 2/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 9.9 MiB zu Beginn, Maximum 105.0, Minimum 9.9, Ende 104.5; anon Beginn 9.0, Maximum 103.8, Minimum 9.0, Ende 103.1 MiB
Stufe 100000 narrow, Run 2/3 — Zeilen in cdc.change nach dem Nachlauf: 100000, VmHWM 157736 KiB, memory.peak seit Start des Feeds 145.3 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 52, größter Heap zu Beginn 104 MB, größtes Ziel 104 MB, letzter Heap nach GC 44 MB; seit dem Start des Feeds: Bereinigung gelaufen 14, fehlgeschlagen 0
Stufe 100000 narrow, Run 2/3 — cdc.change geleert, Nachlauf 120 s: Nachlauf nach der Leerung: memory.current 104.5 MiB zu Beginn, Maximum 104.9, Minimum 88.0, Ende 88.0; anon Beginn 103.1, Maximum 103.5, Minimum 86.7, Ende 86.7 MiB; im Fenster: GC-Läufe 1, größter Heap zu Beginn 85 MB, größtes Ziel 89 MB, letzter Heap nach GC 1 MB
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 1069: active WalSenderWaitForWal) — Sitzung beendet
Stufe 100000 narrow, Run 3/3 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.0 MiB zu Beginn, Maximum 6.0, Minimum 6.0, Ende 6.0; anon Beginn 5.2, Maximum 5.2, Minimum 5.2, Ende 5.2 MiB
Stufe 100000 narrow, Run 3/3 — Kopierdauer 11339 ms, 8819 Zeilen/s; Proben 44; memory.current Beginn 6.5, bei 25% 10.1, 50% 10.4, 75% 10.6, Ende 9.8, Spitze 11.6 MiB; anon Beginn 5.2, bei 25% 9.1, 50% 9.4, 75% 9.1, Ende 8.8, Spitze 9.8 MiB
Stufe 100000 narrow, Run 3/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 9.8 MiB (anon 8.8, file 0.0, kernel 1.0), Prozess RssAnon 8.8 MiB, docker stats 9.8 MiB, VmHWM 24180 KiB, memory.peak seit Start des Feeds 12.7 MiB; im Run: GC-Läufe 168, größter Heap zu Beginn 5 MB, größtes Ziel 5 MB, letzter Heap nach GC 2 MB
Stufe 100000 narrow, Run 3/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 9.8 MiB zu Beginn, Maximum 130.3, Minimum 9.8, Ende 130.1; anon Beginn 8.8, Maximum 128.7, Minimum 8.8, Ende 128.7 MiB
Stufe 100000 narrow, Run 3/3 — Zeilen in cdc.change nach dem Nachlauf: 100000, VmHWM 160060 KiB, memory.peak seit Start des Feeds 147.7 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 39, größter Heap zu Beginn 112 MB, größtes Ziel 112 MB, letzter Heap nach GC 53 MB; seit dem Start des Feeds: Bereinigung gelaufen 14, fehlgeschlagen 0
Stufe 100000 narrow, Run 3/3 — cdc.change geleert, Nachlauf 120 s: Nachlauf nach der Leerung: memory.current 130.1 MiB zu Beginn, Maximum 131.2, Minimum 108.1, Ende 108.1; anon Beginn 128.7, Maximum 128.8, Minimum 106.7, Ende 106.7 MiB; im Fenster: GC-Läufe 1, größter Heap zu Beginn 105 MB, größtes Ziel 107 MB, letzter Heap nach GC 1 MB
Ende — Lauf 20260925T132801Z
```

## Reihe E — Lauf `20260925T134215Z`

Image ohne Bereinigung (`retentionInterval` = 1000 h), schmal, Stufen wie A.

```text
Host — Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU, 33362599936 Byte RAM, PostgreSQL-Image postgres:18-alpine, Feed-Image-ID sha256:11c0107110859a525df6ef1ab36cd885f8cc4d8384cc7a7020d0e4fb3dbc83c5
Stufen 10000 100000 200000, 3 Runs je Stufe, Zeilenbreite narrow, Feed-Umgebung 'GODEBUG=gctrace=1', Feed-Docker-Argumente '', Feed-Image pg-change-feed:bench-noret
Tabelle bench_mem_narrow_10000 — 10000 Zeilen, mittlere Zeilenbreite 73 B (pg_column_size), 64 B (Textform)
Tabelle bench_mem_narrow_100000 — 100000 Zeilen, mittlere Zeilenbreite 73 B (pg_column_size), 67 B (Textform)
Tabelle bench_mem_narrow_200000 — 200000 Zeilen, mittlere Zeilenbreite 74 B (pg_column_size), 69 B (Textform)
Stufe 10000 narrow, Run 1/3 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 22.1 MiB zu Beginn, Maximum 22.1, Minimum 22.1, Ende 22.1; anon Beginn 5.2, Maximum 5.2, Minimum 5.2, Ende 5.2 MiB
Stufe 10000 narrow, Run 1/3 — Kopierdauer 1157 ms, 8643 Zeilen/s; Proben 6; memory.current Beginn 22.4, bei 25% 26.3, 50% 26.9, 75% 26.3, Ende 25.8, Spitze 26.9 MiB; anon Beginn 5.2, bei 25% 8.9, 50% 8.9, 75% 8.8, Ende 8.8, Spitze 8.9 MiB
Stufe 10000 narrow, Run 1/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 25.8 MiB (anon 8.8, file 16.0, kernel 0.9), Prozess RssAnon 8.8 MiB, docker stats 25.6 MiB, VmHWM 23432 KiB, memory.peak seit Start des Feeds 29.0 MiB; im Run: GC-Läufe 18, größter Heap zu Beginn 5 MB, größtes Ziel 5 MB, letzter Heap nach GC 2 MB
Stufe 10000 narrow, Run 1/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 25.8 MiB zu Beginn, Maximum 26.1, Minimum 25.8, Ende 25.8; anon Beginn 8.8, Maximum 8.8, Minimum 8.8, Ende 8.8 MiB
Stufe 10000 narrow, Run 1/3 — Zeilen in cdc.change nach dem Nachlauf: 10000, VmHWM 23440 KiB, memory.peak seit Start des Feeds 29.0 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt); seit dem Start des Feeds: Bereinigung gelaufen 0, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 195: active WalSenderWaitForWal) — Sitzung beendet
Stufe 10000 narrow, Run 2/3 — Zeilen in cdc.change vor dem Run: 10000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.0 MiB zu Beginn, Maximum 6.0, Minimum 6.0, Ende 6.0; anon Beginn 5.1, Maximum 5.1, Minimum 5.1, Ende 5.1 MiB
Stufe 10000 narrow, Run 2/3 — Kopierdauer 1232 ms, 8117 Zeilen/s; Proben 6; memory.current Beginn 6.6, bei 25% 10.2, 50% 9.2, 75% 9.3, Ende 9.7, Spitze 10.2 MiB; anon Beginn 5.1, bei 25% 9.3, 50% 8.3, 75% 8.1, Ende 8.8, Spitze 9.3 MiB
Stufe 10000 narrow, Run 2/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 9.7 MiB (anon 8.8, file 0.0, kernel 0.9), Prozess RssAnon 8.8 MiB, docker stats 9.7 MiB, VmHWM 23500 KiB, memory.peak seit Start des Feeds 11.8 MiB; im Run: GC-Läufe 17, größter Heap zu Beginn 5 MB, größtes Ziel 6 MB, letzter Heap nach GC 2 MB
Stufe 10000 narrow, Run 2/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 9.7 MiB zu Beginn, Maximum 9.7, Minimum 9.7, Ende 9.7; anon Beginn 8.8, Maximum 8.8, Minimum 8.8, Ende 8.8 MiB
Stufe 10000 narrow, Run 2/3 — Zeilen in cdc.change nach dem Nachlauf: 20000, VmHWM 23504 KiB, memory.peak seit Start des Feeds 11.8 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt); seit dem Start des Feeds: Bereinigung gelaufen 0, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 808: active WalSenderWaitForWal) — Sitzung beendet
Stufe 10000 narrow, Run 3/3 — Zeilen in cdc.change vor dem Run: 20000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.1 MiB zu Beginn, Maximum 6.1, Minimum 6.1, Ende 6.1; anon Beginn 5.2, Maximum 5.2, Minimum 5.2, Ende 5.2 MiB
Stufe 10000 narrow, Run 3/3 — Kopierdauer 1048 ms, 9542 Zeilen/s; Proben 5; memory.current Beginn 6.5, bei 25% 10.5, 50% 10.8, 75% 9.7, Ende 9.7, Spitze 10.8 MiB; anon Beginn 5.2, bei 25% 9.0, 50% 9.0, 75% 8.3, Ende 8.3, Spitze 9.1 MiB
Stufe 10000 narrow, Run 3/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 9.4 MiB (anon 8.3, file 0.1, kernel 0.9), Prozess RssAnon 8.4 MiB, docker stats 9.4 MiB, VmHWM 23232 KiB, memory.peak seit Start des Feeds 12.1 MiB; im Run: GC-Läufe 17, größter Heap zu Beginn 5 MB, größtes Ziel 5 MB, letzter Heap nach GC 2 MB
Stufe 10000 narrow, Run 3/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 9.4 MiB zu Beginn, Maximum 9.6, Minimum 9.4, Ende 9.4; anon Beginn 8.4, Maximum 8.5, Minimum 8.4, Ende 8.5 MiB
Stufe 10000 narrow, Run 3/3 — Zeilen in cdc.change nach dem Nachlauf: 30000, VmHWM 23296 KiB, memory.peak seit Start des Feeds 12.1 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt); seit dem Start des Feeds: Bereinigung gelaufen 0, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 1426: active WalSenderWaitForWal) — Sitzung beendet
Stufe 100000 narrow, Run 1/3 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.0 MiB zu Beginn, Maximum 6.0, Minimum 6.0, Ende 6.0; anon Beginn 5.1, Maximum 5.1, Minimum 5.1, Ende 5.1 MiB
Stufe 100000 narrow, Run 1/3 — Kopierdauer 11345 ms, 8814 Zeilen/s; Proben 44; memory.current Beginn 6.0, bei 25% 11.3, 50% 10.8, 75% 10.8, Ende 10.0, Spitze 12.0 MiB; anon Beginn 5.1, bei 25% 9.8, 50% 9.0, 75% 9.4, Ende 8.8, Spitze 10.2 MiB
Stufe 100000 narrow, Run 1/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 10.0 MiB (anon 8.8, file 0.1, kernel 1.0), Prozess RssAnon 8.8 MiB, docker stats 9.9 MiB, VmHWM 23708 KiB, memory.peak seit Start des Feeds 12.5 MiB; im Run: GC-Läufe 167, größter Heap zu Beginn 5 MB, größtes Ziel 5 MB, letzter Heap nach GC 2 MB
Stufe 100000 narrow, Run 1/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 9.9 MiB zu Beginn, Maximum 9.9, Minimum 9.9, Ende 9.9; anon Beginn 8.8, Maximum 8.8, Minimum 8.8, Ende 8.8 MiB
Stufe 100000 narrow, Run 1/3 — Zeilen in cdc.change nach dem Nachlauf: 100000, VmHWM 23708 KiB, memory.peak seit Start des Feeds 12.5 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt); seit dem Start des Feeds: Bereinigung gelaufen 0, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 2041: active WalSenderWaitForWal) — Sitzung beendet
Stufe 100000 narrow, Run 2/3 — Zeilen in cdc.change vor dem Run: 100000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.2 MiB zu Beginn, Maximum 6.2, Minimum 6.2, Ende 6.2; anon Beginn 5.3, Maximum 5.3, Minimum 5.3, Ende 5.3 MiB
Stufe 100000 narrow, Run 2/3 — Kopierdauer 11183 ms, 8942 Zeilen/s; Proben 43; memory.current Beginn 6.2, bei 25% 10.0, 50% 10.7, 75% 10.3, Ende 9.0, Spitze 11.5 MiB; anon Beginn 5.3, bei 25% 8.8, 50% 9.2, 75% 9.0, Ende 8.0, Spitze 9.6 MiB
Stufe 100000 narrow, Run 2/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 9.0 MiB (anon 8.0, file 0.0, kernel 1.0), Prozess RssAnon 8.0 MiB, docker stats 9.0 MiB, VmHWM 23392 KiB, memory.peak seit Start des Feeds 12.2 MiB; im Run: GC-Läufe 170, größter Heap zu Beginn 5 MB, größtes Ziel 5 MB, letzter Heap nach GC 1 MB
Stufe 100000 narrow, Run 2/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 9.0 MiB zu Beginn, Maximum 9.2, Minimum 9.0, Ende 9.0; anon Beginn 8.0, Maximum 8.0, Minimum 8.0, Ende 8.0 MiB
Stufe 100000 narrow, Run 2/3 — Zeilen in cdc.change nach dem Nachlauf: 200000, VmHWM 23392 KiB, memory.peak seit Start des Feeds 12.2 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt); seit dem Start des Feeds: Bereinigung gelaufen 0, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 2916: active WalSenderWaitForWal) — Sitzung beendet
Stufe 100000 narrow, Run 3/3 — Zeilen in cdc.change vor dem Run: 200000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.0 MiB zu Beginn, Maximum 6.0, Minimum 6.0, Ende 6.0; anon Beginn 5.3, Maximum 5.3, Minimum 5.3, Ende 5.3 MiB
Stufe 100000 narrow, Run 3/3 — Kopierdauer 11374 ms, 8792 Zeilen/s; Proben 45; memory.current Beginn 6.1, bei 25% 9.8, 50% 11.2, 75% 9.9, Ende 10.0, Spitze 11.6 MiB; anon Beginn 5.3, bei 25% 8.8, 50% 9.5, 75% 7.8, Ende 8.7, Spitze 9.8 MiB
Stufe 100000 narrow, Run 3/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 10.0 MiB (anon 8.7, file 0.2, kernel 1.0), Prozess RssAnon 8.8 MiB, docker stats 10.0 MiB, VmHWM 23924 KiB, memory.peak seit Start des Feeds 12.7 MiB; im Run: GC-Läufe 169, größter Heap zu Beginn 5 MB, größtes Ziel 5 MB, letzter Heap nach GC 2 MB
Stufe 100000 narrow, Run 3/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 10.0 MiB zu Beginn, Maximum 10.3, Minimum 10.0, Ende 10.1; anon Beginn 8.8, Maximum 8.9, Minimum 8.8, Ende 8.9 MiB
Stufe 100000 narrow, Run 3/3 — Zeilen in cdc.change nach dem Nachlauf: 300000, VmHWM 23924 KiB, memory.peak seit Start des Feeds 12.7 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt); seit dem Start des Feeds: Bereinigung gelaufen 0, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 3794: active WalSenderWaitForWal) — Sitzung beendet
Stufe 200000 narrow, Run 1/3 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.1 MiB zu Beginn, Maximum 6.1, Minimum 6.1, Ende 6.1; anon Beginn 5.3, Maximum 5.3, Minimum 5.3, Ende 5.3 MiB
Stufe 200000 narrow, Run 1/3 — Kopierdauer 21396 ms, 9348 Zeilen/s; Proben 82; memory.current Beginn 6.6, bei 25% 10.5, 50% 11.9, 75% 10.5, Ende 10.3, Spitze 11.9 MiB; anon Beginn 5.3, bei 25% 9.4, 50% 9.4, 75% 9.5, Ende 9.3, Spitze 9.9 MiB
Stufe 200000 narrow, Run 1/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 10.3 MiB (anon 9.3, file 0.0, kernel 1.0), Prozess RssAnon 9.3 MiB, docker stats 10.3 MiB, VmHWM 24352 KiB, memory.peak seit Start des Feeds 12.6 MiB; im Run: GC-Läufe 334, größter Heap zu Beginn 5 MB, größtes Ziel 5 MB, letzter Heap nach GC 2 MB
Stufe 200000 narrow, Run 1/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 10.3 MiB zu Beginn, Maximum 10.3, Minimum 8.7, Ende 8.7; anon Beginn 9.3, Maximum 9.3, Minimum 7.8, Ende 7.8 MiB
Stufe 200000 narrow, Run 1/3 — Zeilen in cdc.change nach dem Nachlauf: 200000, VmHWM 23944 KiB, memory.peak seit Start des Feeds 12.6 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 1, größter Heap zu Beginn 4 MB, größtes Ziel 5 MB, letzter Heap nach GC 1 MB; seit dem Start des Feeds: Bereinigung gelaufen 0, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 4695: active WalSenderWaitForWal) — Sitzung beendet
Stufe 200000 narrow, Run 2/3 — Zeilen in cdc.change vor dem Run: 200000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.0 MiB zu Beginn, Maximum 6.0, Minimum 6.0, Ende 6.0; anon Beginn 5.2, Maximum 5.2, Minimum 5.2, Ende 5.2 MiB
Stufe 200000 narrow, Run 2/3 — Kopierdauer 22939 ms, 8719 Zeilen/s; Proben 88; memory.current Beginn 6.3, bei 25% 9.8, 50% 10.3, 75% 10.4, Ende 9.8, Spitze 11.8 MiB; anon Beginn 5.2, bei 25% 8.5, 50% 8.6, 75% 9.1, Ende 8.8, Spitze 9.8 MiB
Stufe 200000 narrow, Run 2/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 9.8 MiB (anon 8.8, file 0.0, kernel 1.0), Prozess RssAnon 8.8 MiB, docker stats 9.8 MiB, VmHWM 23928 KiB, memory.peak seit Start des Feeds 12.8 MiB; im Run: GC-Läufe 334, größter Heap zu Beginn 5 MB, größtes Ziel 5 MB, letzter Heap nach GC 2 MB
Stufe 200000 narrow, Run 2/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 9.8 MiB zu Beginn, Maximum 10.1, Minimum 9.8, Ende 9.8; anon Beginn 8.8, Maximum 8.8, Minimum 8.8, Ende 8.8 MiB
Stufe 200000 narrow, Run 2/3 — Zeilen in cdc.change nach dem Nachlauf: 400000, VmHWM 23928 KiB, memory.peak seit Start des Feeds 12.8 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt); seit dem Start des Feeds: Bereinigung gelaufen 0, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 5844: active WalSenderWaitForWal) — Sitzung beendet
Stufe 200000 narrow, Run 3/3 — Zeilen in cdc.change vor dem Run: 400000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.1 MiB zu Beginn, Maximum 6.1, Minimum 6.1, Ende 6.1; anon Beginn 5.2, Maximum 5.2, Minimum 5.2, Ende 5.2 MiB
Stufe 200000 narrow, Run 3/3 — Kopierdauer 22093 ms, 9053 Zeilen/s; Proben 85; memory.current Beginn 6.4, bei 25% 9.5, 50% 10.3, 75% 10.4, Ende 9.4, Spitze 11.6 MiB; anon Beginn 5.2, bei 25% 8.5, 50% 8.8, 75% 8.9, Ende 8.4, Spitze 10.1 MiB
Stufe 200000 narrow, Run 3/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 9.4 MiB (anon 8.4, file 0.0, kernel 1.0), Prozess RssAnon 8.4 MiB, docker stats 9.4 MiB, VmHWM 23904 KiB, memory.peak seit Start des Feeds 12.6 MiB; im Run: GC-Läufe 336, größter Heap zu Beginn 5 MB, größtes Ziel 5 MB, letzter Heap nach GC 2 MB
Stufe 200000 narrow, Run 3/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 9.4 MiB zu Beginn, Maximum 9.7, Minimum 9.4, Ende 9.5; anon Beginn 8.4, Maximum 8.5, Minimum 8.4, Ende 8.5 MiB
Stufe 200000 narrow, Run 3/3 — Zeilen in cdc.change nach dem Nachlauf: 600000, VmHWM 23904 KiB, memory.peak seit Start des Feeds 12.6 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt); seit dem Start des Feeds: Bereinigung gelaufen 0, fehlgeschlagen 0
Ende — Lauf 20260925T134215Z
```

## Reihe F — Lauf `20260925T140625Z`

Wie E, Blockgröße 100, Stufe 200.000.

```text
Host — Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU, 33362599936 Byte RAM, PostgreSQL-Image postgres:18-alpine, Feed-Image-ID sha256:b4398d62f046c3822ace76d4cf7e18b5e811dd2700450a9418ded64d97001e59
Stufen 200000, 3 Runs je Stufe, Zeilenbreite narrow, Feed-Umgebung 'GODEBUG=gctrace=1', Feed-Docker-Argumente '', Feed-Image pg-change-feed:bench-noret-b100
Tabelle bench_mem_narrow_200000 — 200000 Zeilen, mittlere Zeilenbreite 74 B (pg_column_size), 69 B (Textform)
Stufe 200000 narrow, Run 1/3 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 22.1 MiB zu Beginn, Maximum 22.1, Minimum 22.1, Ende 22.1; anon Beginn 5.3, Maximum 5.3, Minimum 5.3, Ende 5.3 MiB
Stufe 200000 narrow, Run 1/3 — Kopierdauer 30047 ms, 6656 Zeilen/s; Proben 115; memory.current Beginn 22.1, bei 25% 25.6, 50% 25.8, 75% 25.9, Ende 25.6, Spitze 26.8 MiB; anon Beginn 5.3, bei 25% 7.9, 50% 7.8, 75% 8.2, Ende 8.2, Spitze 8.4 MiB
Stufe 200000 narrow, Run 1/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 25.6 MiB (anon 8.2, file 16.4, kernel 1.0), Prozess RssAnon 8.2 MiB, docker stats 25.5 MiB, VmHWM 22988 KiB, memory.peak seit Start des Feeds 27.9 MiB; im Run: GC-Läufe 467, größter Heap zu Beginn 4 MB, größtes Ziel 4 MB, letzter Heap nach GC 1 MB
Stufe 200000 narrow, Run 1/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 25.6 MiB zu Beginn, Maximum 25.9, Minimum 25.6, Ende 25.6; anon Beginn 8.2, Maximum 8.2, Minimum 8.2, Ende 8.2 MiB
Stufe 200000 narrow, Run 1/3 — Zeilen in cdc.change nach dem Nachlauf: 200000, VmHWM 22988 KiB, memory.peak seit Start des Feeds 27.9 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt); seit dem Start des Feeds: Bereinigung gelaufen 0, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 153: active WalSenderWaitForWal) — Sitzung beendet
Stufe 200000 narrow, Run 2/3 — Zeilen in cdc.change vor dem Run: 200000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.1 MiB zu Beginn, Maximum 6.3, Minimum 6.1, Ende 6.1; anon Beginn 5.2, Maximum 5.2, Minimum 5.2, Ende 5.2 MiB
Stufe 200000 narrow, Run 2/3 — Kopierdauer 30847 ms, 6484 Zeilen/s; Proben 117; memory.current Beginn 6.4, bei 25% 9.2, 50% 10.0, 75% 10.0, Ende 9.6, Spitze 10.5 MiB; anon Beginn 5.2, bei 25% 7.9, 50% 8.2, 75% 8.3, Ende 8.2, Spitze 8.6 MiB
Stufe 200000 narrow, Run 2/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 9.2 MiB (anon 8.1, file 0.1, kernel 1.0), Prozess RssAnon 8.1 MiB, docker stats 9.2 MiB, VmHWM 23228 KiB, memory.peak seit Start des Feeds 11.6 MiB; im Run: GC-Läufe 466, größter Heap zu Beginn 4 MB, größtes Ziel 4 MB, letzter Heap nach GC 1 MB
Stufe 200000 narrow, Run 2/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 9.2 MiB zu Beginn, Maximum 9.2, Minimum 9.2, Ende 9.2; anon Beginn 8.1, Maximum 8.1, Minimum 8.1, Ende 8.1 MiB
Stufe 200000 narrow, Run 2/3 — Zeilen in cdc.change nach dem Nachlauf: 400000, VmHWM 23232 KiB, memory.peak seit Start des Feeds 11.6 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt); seit dem Start des Feeds: Bereinigung gelaufen 0, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 1532: active WalSenderWaitForWal) — Sitzung beendet
Stufe 200000 narrow, Run 3/3 — Zeilen in cdc.change vor dem Run: 400000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.4 MiB zu Beginn, Maximum 6.4, Minimum 6.4, Ende 6.4; anon Beginn 5.5, Maximum 5.5, Minimum 5.5, Ende 5.5 MiB
Stufe 200000 narrow, Run 3/3 — Kopierdauer 30015 ms, 6663 Zeilen/s; Proben 114; memory.current Beginn 6.6, bei 25% 8.8, 50% 9.4, 75% 9.4, Ende 9.5, Spitze 10.4 MiB; anon Beginn 5.5, bei 25% 7.8, 50% 8.2, 75% 8.2, Ende 8.3, Spitze 8.5 MiB
Stufe 200000 narrow, Run 3/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 9.2 MiB (anon 8.3, file 0.0, kernel 1.0), Prozess RssAnon 8.3 MiB, docker stats 9.2 MiB, VmHWM 23244 KiB, memory.peak seit Start des Feeds 11.2 MiB; im Run: GC-Läufe 470, größter Heap zu Beginn 4 MB, größtes Ziel 4 MB, letzter Heap nach GC 1 MB
Stufe 200000 narrow, Run 3/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 9.2 MiB zu Beginn, Maximum 9.2, Minimum 9.2, Ende 9.2; anon Beginn 8.3, Maximum 8.3, Minimum 8.3, Ende 8.3 MiB
Stufe 200000 narrow, Run 3/3 — Zeilen in cdc.change nach dem Nachlauf: 600000, VmHWM 23244 KiB, memory.peak seit Start des Feeds 11.2 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt); seit dem Start des Feeds: Bereinigung gelaufen 0, fehlgeschlagen 0
Ende — Lauf 20260925T140625Z
```

## Reihe G — Lauf `20260925T141522Z`

Wie E, Blockgröße 10.000, Stufe 200.000.

```text
Host — Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU, 33362599936 Byte RAM, PostgreSQL-Image postgres:18-alpine, Feed-Image-ID sha256:eaa3b44ffadb125d3bb55a592391d96107652ea0a7476121674fd10b6807466a
Stufen 200000, 3 Runs je Stufe, Zeilenbreite narrow, Feed-Umgebung 'GODEBUG=gctrace=1', Feed-Docker-Argumente '', Feed-Image pg-change-feed:bench-noret-b10000
Tabelle bench_mem_narrow_200000 — 200000 Zeilen, mittlere Zeilenbreite 74 B (pg_column_size), 69 B (Textform)
Stufe 200000 narrow, Run 1/3 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 22.4 MiB zu Beginn, Maximum 22.4, Minimum 22.4, Ende 22.4; anon Beginn 5.5, Maximum 5.5, Minimum 5.5, Ende 5.5 MiB
Stufe 200000 narrow, Run 1/3 — Kopierdauer 20538 ms, 9738 Zeilen/s; Proben 78; memory.current Beginn 22.7, bei 25% 41.3, 50% 36.8, 75% 33.5, Ende 34.2, Spitze 42.6 MiB; anon Beginn 5.5, bei 25% 23.8, 50% 18.9, 75% 16.1, Ende 16.8, Spitze 24.4 MiB
Stufe 200000 narrow, Run 1/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 34.2 MiB (anon 16.8, file 16.4, kernel 1.0), Prozess RssAnon 16.8 MiB, docker stats 34.0 MiB, VmHWM 38092 KiB, memory.peak seit Start des Feeds 43.3 MiB; im Run: GC-Läufe 137, größter Heap zu Beginn 19 MB, größtes Ziel 20 MB, letzter Heap nach GC 5 MB
Stufe 200000 narrow, Run 1/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 34.2 MiB zu Beginn, Maximum 34.4, Minimum 34.2, Ende 34.2; anon Beginn 16.8, Maximum 16.8, Minimum 16.8, Ende 16.8 MiB
Stufe 200000 narrow, Run 1/3 — Zeilen in cdc.change nach dem Nachlauf: 200000, VmHWM 38092 KiB, memory.peak seit Start des Feeds 43.3 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt); seit dem Start des Feeds: Bereinigung gelaufen 0, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 154: active WalSenderWaitForWal) — Sitzung beendet
Stufe 200000 narrow, Run 2/3 — Zeilen in cdc.change vor dem Run: 200000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.3 MiB zu Beginn, Maximum 6.3, Minimum 6.3, Ende 6.3; anon Beginn 5.4, Maximum 5.4, Minimum 5.4, Ende 5.4 MiB
Stufe 200000 narrow, Run 2/3 — Kopierdauer 18138 ms, 11027 Zeilen/s; Proben 70; memory.current Beginn 6.6, bei 25% 21.8, 50% 17.9, 75% 22.2, Ende 17.6, Spitze 25.6 MiB; anon Beginn 5.4, bei 25% 20.7, 50% 16.1, 75% 21.0, Ende 16.4, Spitze 24.4 MiB
Stufe 200000 narrow, Run 2/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 17.5 MiB (anon 16.4, file 0.1, kernel 1.0), Prozess RssAnon 16.4 MiB, docker stats 17.6 MiB, VmHWM 38300 KiB, memory.peak seit Start des Feeds 27.3 MiB; im Run: GC-Läufe 141, größter Heap zu Beginn 19 MB, größtes Ziel 20 MB, letzter Heap nach GC 5 MB
Stufe 200000 narrow, Run 2/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 17.5 MiB zu Beginn, Maximum 17.8, Minimum 17.5, Ende 17.5; anon Beginn 16.4, Maximum 16.4, Minimum 16.4, Ende 16.4 MiB
Stufe 200000 narrow, Run 2/3 — Zeilen in cdc.change nach dem Nachlauf: 400000, VmHWM 38300 KiB, memory.peak seit Start des Feeds 27.3 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt); seit dem Start des Feeds: Bereinigung gelaufen 0, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 1276: active WalSenderWaitForWal) — Sitzung beendet
Stufe 200000 narrow, Run 3/3 — Zeilen in cdc.change vor dem Run: 400000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.0 MiB zu Beginn, Maximum 6.0, Minimum 6.0, Ende 6.0; anon Beginn 5.2, Maximum 5.2, Minimum 5.2, Ende 5.2 MiB
Stufe 200000 narrow, Run 3/3 — Kopierdauer 22871 ms, 8745 Zeilen/s; Proben 86; memory.current Beginn 6.1, bei 25% 18.6, 50% 18.9, 75% 20.0, Ende 17.4, Spitze 25.4 MiB; anon Beginn 5.2, bei 25% 17.4, 50% 17.4, 75% 17.8, Ende 16.3, Spitze 24.3 MiB
Stufe 200000 narrow, Run 3/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 17.4 MiB (anon 16.3, file 0.0, kernel 1.1), Prozess RssAnon 16.3 MiB, docker stats 17.4 MiB, VmHWM 41080 KiB, memory.peak seit Start des Feeds 29.4 MiB; im Run: GC-Läufe 139, größter Heap zu Beginn 19 MB, größtes Ziel 20 MB, letzter Heap nach GC 5 MB
Stufe 200000 narrow, Run 3/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 17.4 MiB zu Beginn, Maximum 17.4, Minimum 17.4, Ende 17.4; anon Beginn 16.3, Maximum 16.3, Minimum 16.3, Ende 16.3 MiB
Stufe 200000 narrow, Run 3/3 — Zeilen in cdc.change nach dem Nachlauf: 600000, VmHWM 41080 KiB, memory.peak seit Start des Feeds 29.4 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt); seit dem Start des Feeds: Bereinigung gelaufen 0, fehlgeschlagen 0
Ende — Lauf 20260925T141522Z
```

## Reihe H — Lauf `20260925T142349Z`

Default-Image, breit (Zeile etwa 1,3 KB), Stufen wie A.

```text
Host — Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU, 33362599936 Byte RAM, PostgreSQL-Image postgres:18-alpine, Feed-Image-ID sha256:969b7fad86909255d47fdf1bd5ada33edee15a171d6cc8b8a9849edb3965cf4e
Stufen 10000 100000 200000, 3 Runs je Stufe, Zeilenbreite wide, Feed-Umgebung 'GODEBUG=gctrace=1', Feed-Docker-Argumente '', Feed-Image ghcr.io/pt9912/pg-change-feed:dev
Tabelle bench_mem_wide_10000 — 10000 Zeilen, mittlere Zeilenbreite 1273 B (pg_column_size), 1258 B (Textform)
Tabelle bench_mem_wide_100000 — 100000 Zeilen, mittlere Zeilenbreite 1277 B (pg_column_size), 1263 B (Textform)
Tabelle bench_mem_wide_200000 — 200000 Zeilen, mittlere Zeilenbreite 1279 B (pg_column_size), 1265 B (Textform)
Stufe 10000 wide, Run 1/3 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.4 MiB zu Beginn, Maximum 6.4, Minimum 5.9, Ende 6.0; anon Beginn 5.1, Maximum 5.1, Minimum 5.1, Ende 5.1 MiB
Stufe 10000 wide, Run 1/3 — Kopierdauer 1367 ms, 7315 Zeilen/s; Proben 7; memory.current Beginn 6.3, bei 25% 14.5, 50% 13.8, 75% 14.1, Ende 12.6, Spitze 16.0 MiB; anon Beginn 5.1, bei 25% 12.2, 50% 11.6, 75% 11.7, Ende 11.7, Spitze 14.4 MiB
Stufe 10000 wide, Run 1/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 12.6 MiB (anon 11.7, file 0.0, kernel 1.0), Prozess RssAnon 11.7 MiB, docker stats 12.6 MiB, VmHWM 28328 KiB, memory.peak seit Start des Feeds 17.1 MiB; im Run: GC-Läufe 42, größter Heap zu Beginn 9 MB, größtes Ziel 10 MB, letzter Heap nach GC 3 MB
Stufe 10000 wide, Run 1/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 12.6 MiB zu Beginn, Maximum 46.0, Minimum 12.6, Ende 42.0; anon Beginn 11.7, Maximum 45.0, Minimum 11.7, Ende 41.0 MiB
Stufe 10000 wide, Run 1/3 — Zeilen in cdc.change nach dem Nachlauf: 10000, VmHWM 59668 KiB, memory.peak seit Start des Feeds 49.5 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 19, größter Heap zu Beginn 39 MB, größtes Ziel 39 MB, letzter Heap nach GC 16 MB; seit dem Start des Feeds: Bereinigung gelaufen 13, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 210: active WalSenderWaitForWal) — Sitzung beendet
Stufe 10000 wide, Run 2/3 — Zeilen in cdc.change vor dem Run: 10000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 30.2 MiB zu Beginn, Maximum 38.4, Minimum 30.2, Ende 38.4; anon Beginn 29.2, Maximum 37.4, Minimum 29.2, Ende 37.4 MiB
Stufe 10000 wide, Run 2/3 — Kopierdauer 1327 ms, 7536 Zeilen/s; Proben 6; memory.current Beginn 38.9, bei 25% 14.4, 50% 16.4, 75% 14.3, Ende 18.3, Spitze 38.9 MiB; anon Beginn 37.4, bei 25% 12.9, 50% 14.9, 75% 13.0, Ende 16.1, Spitze 37.4 MiB
Stufe 10000 wide, Run 2/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 17.1 MiB (anon 16.1, file 0.0, kernel 1.0), Prozess RssAnon 16.1 MiB, docker stats 17.1 MiB, VmHWM 54524 KiB, memory.peak seit Start des Feeds 43.3 MiB; im Run: GC-Läufe 42, größter Heap zu Beginn 28 MB, größtes Ziel 29 MB, letzter Heap nach GC 3 MB
Stufe 10000 wide, Run 2/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 17.1 MiB zu Beginn, Maximum 77.0, Minimum 17.1, Ende 74.9; anon Beginn 16.1, Maximum 74.6, Minimum 16.1, Ende 73.8 MiB
Stufe 10000 wide, Run 2/3 — Zeilen in cdc.change nach dem Nachlauf: 20000, VmHWM 93484 KiB, memory.peak seit Start des Feeds 81.0 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 26, größter Heap zu Beginn 67 MB, größtes Ziel 69 MB, letzter Heap nach GC 30 MB; seit dem Start des Feeds: Bereinigung gelaufen 13, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 835: active WalSenderWaitForWal) — Sitzung beendet
Stufe 10000 wide, Run 3/3 — Zeilen in cdc.change vor dem Run: 20000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 50.8 MiB zu Beginn, Maximum 51.1, Minimum 50.8, Ende 51.1; anon Beginn 49.7, Maximum 50.0, Minimum 49.7, Ende 50.0 MiB
Stufe 10000 wide, Run 3/3 — Kopierdauer 1448 ms, 6906 Zeilen/s; Proben 7; memory.current Beginn 51.7, bei 25% 14.7, 50% 14.0, 75% 15.4, Ende 16.8, Spitze 51.7 MiB; anon Beginn 50.0, bei 25% 13.0, 50% 12.1, 75% 13.6, Ende 15.5, Spitze 50.0 MiB
Stufe 10000 wide, Run 3/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 16.6 MiB (anon 15.5, file 0.1, kernel 1.1), Prozess RssAnon 15.5 MiB, docker stats 16.6 MiB, VmHWM 79672 KiB, memory.peak seit Start des Feeds 67.5 MiB; im Run: GC-Läufe 40, größter Heap zu Beginn 51 MB, größtes Ziel 52 MB, letzter Heap nach GC 3 MB
Stufe 10000 wide, Run 3/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 16.6 MiB zu Beginn, Maximum 109.7, Minimum 16.6, Ende 105.8; anon Beginn 15.5, Maximum 108.1, Minimum 15.5, Ende 104.5 MiB
Stufe 10000 wide, Run 3/3 — Zeilen in cdc.change nach dem Nachlauf: 30000, VmHWM 128336 KiB, memory.peak seit Start des Feeds 116.6 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 25, größter Heap zu Beginn 100 MB, größtes Ziel 103 MB, letzter Heap nach GC 41 MB; seit dem Start des Feeds: Bereinigung gelaufen 13, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 1441: active WalSenderWaitForWal) — Sitzung beendet
Stufe 100000 wide, Run 1/3 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.2 MiB zu Beginn, Maximum 6.2, Minimum 6.2, Ende 6.2; anon Beginn 5.4, Maximum 5.4, Minimum 5.4, Ende 5.4 MiB
Stufe 100000 wide, Run 1/3 — Kopierdauer 15276 ms, 6546 Zeilen/s; Proben 58; memory.current Beginn 6.5, bei 25% 15.4, 50% 15.3, 75% 13.7, Ende 13.3, Spitze 17.3 MiB; anon Beginn 5.4, bei 25% 14.0, 50% 13.8, 75% 11.9, Ende 12.2, Spitze 14.7 MiB
Stufe 100000 wide, Run 1/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 13.3 MiB (anon 12.2, file 0.0, kernel 1.0), Prozess RssAnon 12.2 MiB, docker stats 170.1 MiB, VmHWM 192768 KiB, memory.peak seit Start des Feeds 182.3 MiB; im Run: GC-Läufe 441, größter Heap zu Beginn 155 MB, größtes Ziel 159 MB, letzter Heap nach GC 120 MB
Stufe 100000 wide, Run 1/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 223.9 MiB zu Beginn, Maximum 359.9, Minimum 223.4, Ende 359.9; anon Beginn 222.2, Maximum 355.8, Minimum 221.8, Ende 355.8 MiB
Stufe 100000 wide, Run 1/3 — Zeilen in cdc.change nach dem Nachlauf: 100000, VmHWM 378272 KiB, memory.peak seit Start des Feeds 361.0 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 21, größter Heap zu Beginn 328 MB, größtes Ziel 332 MB, letzter Heap nach GC 125 MB; seit dem Start des Feeds: Bereinigung gelaufen 15, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 2073: active WalSenderWaitForWal) — Sitzung beendet
Stufe 100000 wide, Run 2/3 — Zeilen in cdc.change vor dem Run: 100000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 277.9 MiB zu Beginn, Maximum 277.9, Minimum 244.0, Ende 244.0; anon Beginn 276.1, Maximum 276.1, Minimum 242.2, Ende 242.2 MiB
Stufe 100000 wide, Run 2/3 — Kopierdauer 13099 ms, 7634 Zeilen/s; Proben 51; memory.current Beginn 244.3, bei 25% 20.9, 50% 143.1, 75% 22.7, Ende 19.5, Spitze 312.7 MiB; anon Beginn 242.2, bei 25% 18.3, 50% 139.5, 75% 20.0, Ende 17.2, Spitze 310.1 MiB
Stufe 100000 wide, Run 2/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 19.1 MiB (anon 17.2, file 0.0, kernel 1.8), Prozess RssAnon 17.2 MiB, docker stats 19.0 MiB, VmHWM 387484 KiB, memory.peak seit Start des Feeds 373.1 MiB; im Run: GC-Läufe 394, größter Heap zu Beginn 296 MB, größtes Ziel 300 MB, letzter Heap nach GC 3 MB
Stufe 100000 wide, Run 2/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 19.0 MiB zu Beginn, Maximum 799.2, Minimum 19.0, Ende 680.9; anon Beginn 17.2, Maximum 796.4, Minimum 17.2, Ende 678.1 MiB
Stufe 100000 wide, Run 2/3 — Zeilen in cdc.change nach dem Nachlauf: 200000, VmHWM 844396 KiB, memory.peak seit Start des Feeds 818.1 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 26, größter Heap zu Beginn 724 MB, größtes Ziel 728 MB, letzter Heap nach GC 296 MB; seit dem Start des Feeds: Bereinigung gelaufen 14, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 3050: active WalSenderWaitForWal) — Sitzung beendet
Stufe 100000 wide, Run 3/3 — Zeilen in cdc.change vor dem Run: 200000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 692.2 MiB zu Beginn, Maximum 692.3, Minimum 388.0, Ende 578.8; anon Beginn 689.4, Maximum 689.4, Minimum 385.1, Ende 576.0 MiB
Stufe 100000 wide, Run 3/3 — Kopierdauer 13884 ms, 7203 Zeilen/s; Proben 54; memory.current Beginn 578.8, bei 25% 637.7, 50% 243.9, 75% 37.0, Ende 32.7, Spitze 830.1 MiB; anon Beginn 576.0, bei 25% 633.8, 50% 240.6, 75% 32.8, Ende 29.8, Spitze 826.1 MiB
Stufe 100000 wide, Run 3/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 32.7 MiB (anon 29.8, file 0.0, kernel 2.8), Prozess RssAnon 29.8 MiB, docker stats 32.7 MiB, VmHWM 868320 KiB, memory.peak seit Start des Feeds 840.6 MiB; im Run: GC-Läufe 249, größter Heap zu Beginn 801 MB, größtes Ziel 805 MB, letzter Heap nach GC 3 MB
Stufe 100000 wide, Run 3/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 32.7 MiB zu Beginn, Maximum 1015.2, Minimum 32.7, Ende 736.3; anon Beginn 29.8, Maximum 1012.0, Minimum 29.8, Ende 733.0 MiB
Stufe 100000 wide, Run 3/3 — Zeilen in cdc.change nach dem Nachlauf: 300000, VmHWM 1061104 KiB, memory.peak seit Start des Feeds 1033.0 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 33, größter Heap zu Beginn 980 MB, größtes Ziel 980 MB, letzter Heap nach GC 365 MB; seit dem Start des Feeds: Bereinigung gelaufen 14, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 3998: active WalSenderWaitForWal) — Sitzung beendet
Stufe 200000 wide, Run 1/3 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.2 MiB zu Beginn, Maximum 6.2, Minimum 6.2, Ende 6.2; anon Beginn 5.4, Maximum 5.4, Minimum 5.4, Ende 5.4 MiB
Stufe 200000 wide, Run 1/3 — Kopierdauer 25148 ms, 7953 Zeilen/s; Proben 97; memory.current Beginn 6.7, bei 25% 15.0, 50% 16.3, 75% 13.4, Ende 14.3, Spitze 17.3 MiB; anon Beginn 5.4, bei 25% 13.7, 50% 14.1, 75% 11.7, Ende 12.4, Spitze 14.9 MiB
Stufe 200000 wide, Run 1/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 13.5 MiB (anon 12.4, file 0.0, kernel 1.0), Prozess RssAnon 12.4 MiB, docker stats 13.5 MiB, VmHWM 29384 KiB, memory.peak seit Start des Feeds 19.0 MiB; im Run: GC-Läufe 858, größter Heap zu Beginn 10 MB, größtes Ziel 10 MB, letzter Heap nach GC 3 MB
Stufe 200000 wide, Run 1/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 13.5 MiB zu Beginn, Maximum 691.1, Minimum 13.5, Ende 229.4; anon Beginn 12.4, Maximum 685.3, Minimum 12.4, Ende 225.9 MiB
Stufe 200000 wide, Run 1/3 — Zeilen in cdc.change nach dem Nachlauf: 200000, VmHWM 749076 KiB, memory.peak seit Start des Feeds 725.5 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 39, größter Heap zu Beginn 686 MB, größtes Ziel 725 MB, letzter Heap nach GC 289 MB; seit dem Start des Feeds: Bereinigung gelaufen 16, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 4968: active WalSenderWaitForWal) — Sitzung beendet
Stufe 200000 wide, Run 2/3 — Zeilen in cdc.change vor dem Run: 200000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 490.5 MiB zu Beginn, Maximum 631.4, Minimum 459.1, Ende 536.5; anon Beginn 487.7, Maximum 628.4, Minimum 455.6, Ende 533.7 MiB
Stufe 200000 wide, Run 2/3 — Kopierdauer 29878 ms, 6694 Zeilen/s; Proben 115; memory.current Beginn 536.8, bei 25% 494.7, 50% 30.8, 75% 32.2, Ende 708.5, Spitze 830.6 MiB; anon Beginn 533.7, bei 25% 491.4, 50% 27.2, 75% 28.6, Ende 705.6, Spitze 827.3 MiB
Stufe 200000 wide, Run 2/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 708.5 MiB (anon 705.6, file 0.0, kernel 2.9), Prozess RssAnon 705.6 MiB, docker stats 708.5 MiB, VmHWM 864068 KiB, memory.peak seit Start des Feeds 836.2 MiB; im Run: GC-Läufe 498, größter Heap zu Beginn 801 MB, größtes Ziel 805 MB, letzter Heap nach GC 362 MB
Stufe 200000 wide, Run 2/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 708.5 MiB zu Beginn, Maximum 711.4, Minimum 318.4, Ende 711.4; anon Beginn 705.6, Maximum 706.7, Minimum 314.6, Ende 706.7 MiB
Stufe 200000 wide, Run 2/3 — Zeilen in cdc.change nach dem Nachlauf: 400000, VmHWM 1095472 KiB, memory.peak seit Start des Feeds 1058.8 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 8, größter Heap zu Beginn 719 MB, größtes Ziel 726 MB, letzter Heap nach GC 521 MB; seit dem Start des Feeds: Bereinigung gelaufen 11, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 6234: active WalSenderWaitForWal) — Sitzung beendet
Stufe 200000 wide, Run 3/3 — Zeilen in cdc.change vor dem Run: 400000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 975.6 MiB zu Beginn, Maximum 975.6, Minimum 975.6, Ende 975.6; anon Beginn 971.4, Maximum 971.4, Minimum 971.4, Ende 971.4 MiB
Stufe 200000 wide, Run 3/3 — Kopierdauer 33347 ms, 5998 Zeilen/s; Proben 127; memory.current Beginn 975.9, bei 25% 42.0, 50% 611.5, 75% 768.8, Ende 773.0, Spitze 1231.9 MiB; anon Beginn 971.4, bei 25% 36.6, 50% 607.1, 75% 764.2, Ende 768.6, Spitze 1227.7 MiB
Stufe 200000 wide, Run 3/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 773.0 MiB (anon 768.6, file 0.2, kernel 4.2), Prozess RssAnon 768.6 MiB, docker stats 773.0 MiB, VmHWM 1593492 KiB, memory.peak seit Start des Feeds 1550.6 MiB; im Run: GC-Läufe 291, größter Heap zu Beginn 1130 MB, größtes Ziel 1136 MB, letzter Heap nach GC 368 MB
Stufe 200000 wide, Run 3/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 773.0 MiB zu Beginn, Maximum 1441.1, Minimum 773.0, Ende 1441.1; anon Beginn 768.6, Maximum 1436.8, Minimum 768.6, Ende 1436.8 MiB
Stufe 200000 wide, Run 3/3 — Zeilen in cdc.change nach dem Nachlauf: 600000, VmHWM 1593492 KiB, memory.peak seit Start des Feeds 1550.6 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 2, größter Heap zu Beginn 931 MB, größtes Ziel 937 MB, letzter Heap nach GC 716 MB; seit dem Start des Feeds: Bereinigung gelaufen 6, fehlgeschlagen 0
Ende — Lauf 20260925T142349Z
```

## Reihe I — Lauf `20260925T144833Z`

Wie E, breit, Stufe 200.000.

```text
Host — Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU, 33362599936 Byte RAM, PostgreSQL-Image postgres:18-alpine, Feed-Image-ID sha256:11c0107110859a525df6ef1ab36cd885f8cc4d8384cc7a7020d0e4fb3dbc83c5
Stufen 200000, 3 Runs je Stufe, Zeilenbreite wide, Feed-Umgebung 'GODEBUG=gctrace=1', Feed-Docker-Argumente '', Feed-Image pg-change-feed:bench-noret
Tabelle bench_mem_wide_200000 — 200000 Zeilen, mittlere Zeilenbreite 1279 B (pg_column_size), 1265 B (Textform)
Stufe 200000 wide, Run 1/3 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 18.9 MiB zu Beginn, Maximum 19.2, Minimum 18.9, Ende 18.9; anon Beginn 5.2, Maximum 5.2, Minimum 5.2, Ende 5.2 MiB
Stufe 200000 wide, Run 1/3 — Kopierdauer 26727 ms, 7483 Zeilen/s; Proben 103; memory.current Beginn 19.7, bei 25% 29.1, 50% 29.2, 75% 27.8, Ende 26.8, Spitze 30.3 MiB; anon Beginn 5.2, bei 25% 14.7, 50% 14.3, 75% 12.6, Ende 11.8, Spitze 15.7 MiB
Stufe 200000 wide, Run 1/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 26.0 MiB (anon 11.8, file 13.1, kernel 1.0), Prozess RssAnon 11.8 MiB, docker stats 25.8 MiB, VmHWM 28420 KiB, memory.peak seit Start des Feeds 31.5 MiB; im Run: GC-Läufe 864, größter Heap zu Beginn 10 MB, größtes Ziel 10 MB, letzter Heap nach GC 3 MB
Stufe 200000 wide, Run 1/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 25.9 MiB zu Beginn, Maximum 25.9, Minimum 23.0, Ende 23.0; anon Beginn 11.8, Maximum 11.8, Minimum 8.8, Ende 8.8 MiB
Stufe 200000 wide, Run 1/3 — Zeilen in cdc.change nach dem Nachlauf: 200000, VmHWM 28420 KiB, memory.peak seit Start des Feeds 31.5 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 1, größter Heap zu Beginn 6 MB, größtes Ziel 7 MB, letzter Heap nach GC 1 MB; seit dem Start des Feeds: Bereinigung gelaufen 0, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 162: active WalSenderWaitForWal) — Sitzung beendet
Stufe 200000 wide, Run 2/3 — Zeilen in cdc.change vor dem Run: 200000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.2 MiB zu Beginn, Maximum 6.2, Minimum 6.2, Ende 6.2; anon Beginn 5.4, Maximum 5.4, Minimum 5.4, Ende 5.4 MiB
Stufe 200000 wide, Run 2/3 — Kopierdauer 30298 ms, 6601 Zeilen/s; Proben 116; memory.current Beginn 6.5, bei 25% 13.8, 50% 15.5, 75% 16.2, Ende 13.5, Spitze 17.1 MiB; anon Beginn 5.4, bei 25% 12.5, 50% 14.5, 75% 15.2, Ende 12.1, Spitze 15.2 MiB
Stufe 200000 wide, Run 2/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 13.1 MiB (anon 12.1, file 0.0, kernel 1.0), Prozess RssAnon 12.1 MiB, docker stats 9.8 MiB, VmHWM 29396 KiB, memory.peak seit Start des Feeds 18.4 MiB; im Run: GC-Läufe 867, größter Heap zu Beginn 10 MB, größtes Ziel 11 MB, letzter Heap nach GC 1 MB
Stufe 200000 wide, Run 2/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 9.8 MiB zu Beginn, Maximum 10.0, Minimum 9.8, Ende 9.8; anon Beginn 8.7, Maximum 8.7, Minimum 8.7, Ende 8.7 MiB
Stufe 200000 wide, Run 2/3 — Zeilen in cdc.change nach dem Nachlauf: 400000, VmHWM 29396 KiB, memory.peak seit Start des Feeds 18.4 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt); seit dem Start des Feeds: Bereinigung gelaufen 0, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 1462: active WalSenderWaitForWal) — Sitzung beendet
Stufe 200000 wide, Run 3/3 — Zeilen in cdc.change vor dem Run: 400000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.2 MiB zu Beginn, Maximum 6.2, Minimum 6.2, Ende 6.2; anon Beginn 5.3, Maximum 5.3, Minimum 5.3, Ende 5.3 MiB
Stufe 200000 wide, Run 3/3 — Kopierdauer 27104 ms, 7379 Zeilen/s; Proben 104; memory.current Beginn 7.0, bei 25% 13.4, 50% 16.5, 75% 14.9, Ende 12.9, Spitze 17.1 MiB; anon Beginn 5.3, bei 25% 12.1, 50% 14.7, 75% 13.3, Ende 11.9, Spitze 15.4 MiB
Stufe 200000 wide, Run 3/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 12.9 MiB (anon 11.9, file 0.0, kernel 1.0), Prozess RssAnon 11.9 MiB, docker stats 12.9 MiB, VmHWM 29260 KiB, memory.peak seit Start des Feeds 18.3 MiB; im Run: GC-Läufe 868, größter Heap zu Beginn 10 MB, größtes Ziel 10 MB, letzter Heap nach GC 3 MB
Stufe 200000 wide, Run 3/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 12.9 MiB zu Beginn, Maximum 12.9, Minimum 12.9, Ende 12.9; anon Beginn 11.9, Maximum 11.9, Minimum 11.9, Ende 11.9 MiB
Stufe 200000 wide, Run 3/3 — Zeilen in cdc.change nach dem Nachlauf: 600000, VmHWM 29260 KiB, memory.peak seit Start des Feeds 18.3 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt); seit dem Start des Feeds: Bereinigung gelaufen 0, fehlgeschlagen 0
Ende — Lauf 20260925T144833Z
```

## Reihe J — Lauf `20260925T145726Z`

Wie E, schmal, 1.000.000, `--memory 6g`.

```text
Host — Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU, 33362599936 Byte RAM, PostgreSQL-Image postgres:18-alpine, Feed-Image-ID sha256:11c0107110859a525df6ef1ab36cd885f8cc4d8384cc7a7020d0e4fb3dbc83c5
Stufen 1000000, 3 Runs je Stufe, Zeilenbreite narrow, Feed-Umgebung 'GODEBUG=gctrace=1', Feed-Docker-Argumente '--memory 6g', Feed-Image pg-change-feed:bench-noret
Tabelle bench_mem_narrow_1000000 — 1000000 Zeilen, mittlere Zeilenbreite 74 B (pg_column_size), 70 B (Textform)
Stufe 1000000 narrow, Run 1/3 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.2 MiB zu Beginn, Maximum 6.2, Minimum 6.2, Ende 6.2; anon Beginn 5.4, Maximum 5.4, Minimum 5.4, Ende 5.4 MiB
Stufe 1000000 narrow, Run 1/3 — Kopierdauer 118966 ms, 8406 Zeilen/s; Proben 450; memory.current Beginn 6.7, bei 25% 10.7, 50% 11.3, 75% 11.4, Ende 10.6, Spitze 12.7 MiB; anon Beginn 5.4, bei 25% 9.7, 50% 10.0, 75% 10.0, Ende 9.6, Spitze 10.6 MiB
Stufe 1000000 narrow, Run 1/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 10.6 MiB (anon 9.6, file 0.0, kernel 1.0), Prozess RssAnon 9.6 MiB, docker stats 10.6 MiB, VmHWM 25016 KiB, memory.peak seit Start des Feeds 13.3 MiB; im Run: GC-Läufe 1689, größter Heap zu Beginn 5 MB, größtes Ziel 6 MB, letzter Heap nach GC 2 MB
Stufe 1000000 narrow, Run 1/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 10.6 MiB zu Beginn, Maximum 10.6, Minimum 9.8, Ende 9.8; anon Beginn 9.6, Maximum 9.6, Minimum 8.7, Ende 8.8 MiB
Stufe 1000000 narrow, Run 1/3 — Zeilen in cdc.change nach dem Nachlauf: 1000000, VmHWM 25016 KiB, memory.peak seit Start des Feeds 13.3 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 1, größter Heap zu Beginn 5 MB, größtes Ziel 5 MB, letzter Heap nach GC 1 MB; seit dem Start des Feeds: Bereinigung gelaufen 0, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 156: active WalSenderWaitForWal) — Sitzung beendet
Stufe 1000000 narrow, Run 2/3 — Zeilen in cdc.change vor dem Run: 1000000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.1 MiB zu Beginn, Maximum 6.1, Minimum 6.1, Ende 6.1; anon Beginn 5.3, Maximum 5.3, Minimum 5.3, Ende 5.3 MiB
Stufe 1000000 narrow, Run 2/3 — Kopierdauer 100515 ms, 9949 Zeilen/s; Proben 384; memory.current Beginn 6.4, bei 25% 10.1, 50% 11.3, 75% 10.7, Ende 10.5, Spitze 11.9 MiB; anon Beginn 5.3, bei 25% 8.8, 50% 9.4, 75% 9.5, Ende 9.2, Spitze 10.2 MiB
Stufe 1000000 narrow, Run 2/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 10.2 MiB (anon 9.2, file 0.0, kernel 1.0), Prozess RssAnon 9.2 MiB, docker stats 10.2 MiB, VmHWM 24280 KiB, memory.peak seit Start des Feeds 13.2 MiB; im Run: GC-Läufe 1672, größter Heap zu Beginn 5 MB, größtes Ziel 6 MB, letzter Heap nach GC 2 MB
Stufe 1000000 narrow, Run 2/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 10.2 MiB zu Beginn, Maximum 10.5, Minimum 10.2, Ende 10.2; anon Beginn 9.2, Maximum 9.2, Minimum 9.2, Ende 9.2 MiB
Stufe 1000000 narrow, Run 2/3 — Zeilen in cdc.change nach dem Nachlauf: 2000000, VmHWM 24280 KiB, memory.peak seit Start des Feeds 13.2 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt); seit dem Start des Feeds: Bereinigung gelaufen 0, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 3883: active WalSenderWaitForWal) — Sitzung beendet
Stufe 1000000 narrow, Run 3/3 — Zeilen in cdc.change vor dem Run: 2000000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.2 MiB zu Beginn, Maximum 6.5, Minimum 6.2, Ende 6.2; anon Beginn 5.4, Maximum 5.4, Minimum 5.4, Ende 5.4 MiB
Stufe 1000000 narrow, Run 3/3 — Kopierdauer 115530 ms, 8656 Zeilen/s; Proben 440; memory.current Beginn 6.8, bei 25% 10.4, 50% 10.1, 75% 10.4, Ende 11.0, Spitze 12.1 MiB; anon Beginn 5.4, bei 25% 9.4, 50% 8.9, 75% 9.4, Ende 10.0, Spitze 10.3 MiB
Stufe 1000000 narrow, Run 3/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 11.0 MiB (anon 10.0, file 0.0, kernel 1.0), Prozess RssAnon 10.0 MiB, docker stats 11.0 MiB, VmHWM 25060 KiB, memory.peak seit Start des Feeds 13.7 MiB; im Run: GC-Läufe 1667, größter Heap zu Beginn 5 MB, größtes Ziel 6 MB, letzter Heap nach GC 2 MB
Stufe 1000000 narrow, Run 3/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 11.0 MiB zu Beginn, Maximum 11.2, Minimum 11.0, Ende 11.0; anon Beginn 10.0, Maximum 10.0, Minimum 10.0, Ende 10.0 MiB
Stufe 1000000 narrow, Run 3/3 — Zeilen in cdc.change nach dem Nachlauf: 3000000, VmHWM 25060 KiB, memory.peak seit Start des Feeds 13.7 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt); seit dem Start des Feeds: Bereinigung gelaufen 0, fehlgeschlagen 0
Ende — Lauf 20260925T145726Z
```

## Reihe K — Lauf `20260925T151031Z`

Default-Image mit `GOGC=25`, schmal, Stufe 200.000.

```text
Host — Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU, 33362599936 Byte RAM, PostgreSQL-Image postgres:18-alpine, Feed-Image-ID sha256:969b7fad86909255d47fdf1bd5ada33edee15a171d6cc8b8a9849edb3965cf4e
Stufen 200000, 3 Runs je Stufe, Zeilenbreite narrow, Feed-Umgebung 'GODEBUG=gctrace=1 GOGC=25', Feed-Docker-Argumente '', Feed-Image ghcr.io/pt9912/pg-change-feed:dev
Tabelle bench_mem_narrow_200000 — 200000 Zeilen, mittlere Zeilenbreite 74 B (pg_column_size), 69 B (Textform)
Stufe 200000 narrow, Run 1/3 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.8 MiB zu Beginn, Maximum 6.8, Minimum 6.8, Ende 6.8; anon Beginn 5.5, Maximum 5.5, Minimum 5.5, Ende 5.5 MiB
Stufe 200000 narrow, Run 1/3 — Kopierdauer 22992 ms, 8699 Zeilen/s; Proben 88; memory.current Beginn 7.3, bei 25% 9.8, 50% 9.5, 75% 9.4, Ende 9.1, Spitze 11.1 MiB; anon Beginn 5.5, bei 25% 7.7, 50% 7.4, 75% 7.8, Ende 7.7, Spitze 8.2 MiB
Stufe 200000 narrow, Run 1/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 9.2 MiB (anon 7.7, file 0.5, kernel 1.0), Prozess RssAnon 7.7 MiB, docker stats 9.2 MiB, VmHWM 23460 KiB, memory.peak seit Start des Feeds 12.4 MiB; im Run: GC-Läufe 801, größter Heap zu Beginn 3 MB, größtes Ziel 3 MB, letzter Heap nach GC 2 MB
Stufe 200000 narrow, Run 1/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 9.2 MiB zu Beginn, Maximum 222.6, Minimum 9.2, Ende 181.9; anon Beginn 7.7, Maximum 219.9, Minimum 7.7, Ende 179.8 MiB
Stufe 200000 narrow, Run 1/3 — Zeilen in cdc.change nach dem Nachlauf: 200000, VmHWM 240012 KiB, memory.peak seit Start des Feeds 224.5 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 151, größter Heap zu Beginn 200 MB, größtes Ziel 202 MB, letzter Heap nach GC 162 MB; seit dem Start des Feeds: Bereinigung gelaufen 15, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 153: active WalSenderWaitForWal) — Sitzung beendet
Stufe 200000 narrow, Run 2/3 — Zeilen in cdc.change vor dem Run: 200000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 222.2 MiB zu Beginn, Maximum 222.2, Minimum 133.1, Ende 162.7; anon Beginn 220.5, Maximum 220.5, Minimum 129.8, Ende 160.9 MiB
Stufe 200000 narrow, Run 2/3 — Kopierdauer 22460 ms, 8905 Zeilen/s; Proben 88; memory.current Beginn 163.5, bei 25% 19.0, 50% 13.8, 75% 191.0, Ende 13.2, Spitze 205.8 MiB; anon Beginn 160.9, bei 25% 16.5, 50% 11.5, 75% 189.3, Ende 11.5, Spitze 203.9 MiB
Stufe 200000 narrow, Run 2/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 13.2 MiB (anon 11.5, file 0.1, kernel 1.6), Prozess RssAnon 11.5 MiB, docker stats 13.2 MiB, VmHWM 239176 KiB, memory.peak seit Start des Feeds 224.3 MiB; im Run: GC-Läufe 770, größter Heap zu Beginn 199 MB, größtes Ziel 202 MB, letzter Heap nach GC 2 MB
Stufe 200000 narrow, Run 2/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 13.2 MiB zu Beginn, Maximum 377.9, Minimum 13.2, Ende 184.5; anon Beginn 11.5, Maximum 375.8, Minimum 11.5, Ende 182.4 MiB
Stufe 200000 narrow, Run 2/3 — Zeilen in cdc.change nach dem Nachlauf: 400000, VmHWM 431084 KiB, memory.peak seit Start des Feeds 415.5 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 89, größter Heap zu Beginn 396 MB, größtes Ziel 396 MB, letzter Heap nach GC 161 MB; seit dem Start des Feeds: Bereinigung gelaufen 11, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 1354: active WalSenderWaitForWal) — Sitzung beendet
Stufe 200000 narrow, Run 3/3 — Zeilen in cdc.change vor dem Run: 400000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 361.3 MiB zu Beginn, Maximum 361.3, Minimum 192.4, Ende 192.4; anon Beginn 359.2, Maximum 359.2, Minimum 190.2, Ende 190.2 MiB
Stufe 200000 narrow, Run 3/3 — Kopierdauer 23322 ms, 8576 Zeilen/s; Proben 90; memory.current Beginn 192.6, bei 25% 178.2, 50% 17.3, 75% 189.2, Ende 178.0, Spitze 397.4 MiB; anon Beginn 190.2, bei 25% 176.0, 50% 14.6, 75% 187.0, Ende 175.9, Spitze 394.2 MiB
Stufe 200000 narrow, Run 3/3 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 178.0 MiB (anon 175.9, file 0.0, kernel 2.1), Prozess RssAnon 175.9 MiB, docker stats 178.0 MiB, VmHWM 431552 KiB, memory.peak seit Start des Feeds 413.2 MiB; im Run: GC-Läufe 348, größter Heap zu Beginn 389 MB, größtes Ziel 393 MB, letzter Heap nach GC 122 MB
Stufe 200000 narrow, Run 3/3 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 178.0 MiB zu Beginn, Maximum 452.9, Minimum 172.8, Ende 452.0; anon Beginn 175.9, Maximum 449.7, Minimum 169.5, Ende 449.7 MiB
Stufe 200000 narrow, Run 3/3 — Zeilen in cdc.change nach dem Nachlauf: 600000, VmHWM 494992 KiB, memory.peak seit Start des Feeds 477.4 MiB, OOMKilled false, Neustarts 0; im Nachlauf: GC-Läufe 32, größter Heap zu Beginn 405 MB, größtes Ziel 405 MB, letzter Heap nach GC 321 MB; seit dem Start des Feeds: Bereinigung gelaufen 8, fehlgeschlagen 0
Ende — Lauf 20260925T151031Z
```

## Reihe L1 — Lauf `20260925T151925Z`

Mutation: `--memory 64m`, 100.000 Zeilen, 2 Runs, mit `GODEBUG=gctrace=1`.

```text
Host — Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU, 33362599936 Byte RAM, PostgreSQL-Image postgres:18-alpine, Feed-Image-ID sha256:969b7fad86909255d47fdf1bd5ada33edee15a171d6cc8b8a9849edb3965cf4e
Stufen 100000, 2 Runs je Stufe, Zeilenbreite narrow, Feed-Umgebung 'GODEBUG=gctrace=1', Feed-Docker-Argumente '--memory 64m', Feed-Image ghcr.io/pt9912/pg-change-feed:dev
Tabelle bench_mem_narrow_100000 — 100000 Zeilen, mittlere Zeilenbreite 73 B (pg_column_size), 67 B (Textform)
Stufe 100000 narrow, Run 1/2 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.0 MiB zu Beginn, Maximum 6.0, Minimum 6.0, Ende 6.0; anon Beginn 5.2, Maximum 5.2, Minimum 5.2, Ende 5.2 MiB
Stufe 100000 narrow, Run 1/2 — Kopierdauer 9811 ms, 10193 Zeilen/s; Proben 39; memory.current Beginn 6.3, bei 25% 11.3, 50% 10.1, 75% 9.5, Ende 10.3, Spitze 12.5 MiB; anon Beginn 5.2, bei 25% 9.9, 50% 8.9, 75% 8.4, Ende 9.3, Spitze 10.1 MiB
Stufe 100000 narrow, Run 1/2 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 10.3 MiB (anon 9.3, file 0.0, kernel 1.0), Prozess RssAnon 9.3 MiB, docker stats 10.3 MiB, VmHWM 24736 KiB, memory.peak seit Start des Feeds 12.6 MiB; im Run: GC-Läufe 168, größter Heap zu Beginn 5 MB, größtes Ziel 5 MB, letzter Heap nach GC 2 MB
tools/bench-backfill-memory.sh: Zeile 58: /sys/fs/cgroup/system.slice/docker-b223a9643905b3dd7d18825478d8375381a1e8565d9de1c5fe19af25031fc9fa.scope/memory.current: Datei oder Verzeichnis nicht gefunden
Feed-Container nicht mehr lesbar: Status exited, Exit 137, OOMKilled true
Stufe 100000 narrow, Run 1/2 — Nachlauf 60 s ohne Schreibzugriff: 
```

## Reihe L2 — Lauf `20260925T152159Z`

Mutation: `--memory 64m`, 100.000 Zeilen, 2 Runs, ohne `BENCH_FEED_ENV`.

```text
Host — Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU, 33362599936 Byte RAM, PostgreSQL-Image postgres:18-alpine, Feed-Image-ID sha256:969b7fad86909255d47fdf1bd5ada33edee15a171d6cc8b8a9849edb3965cf4e
Stufen 100000, 2 Runs je Stufe, Zeilenbreite narrow, Feed-Umgebung '', Feed-Docker-Argumente '--memory 64m', Feed-Image ghcr.io/pt9912/pg-change-feed:dev
Tabelle bench_mem_narrow_100000 — 100000 Zeilen, mittlere Zeilenbreite 73 B (pg_column_size), 67 B (Textform)
Stufe 100000 narrow, Run 1/2 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.2 MiB zu Beginn, Maximum 6.2, Minimum 6.2, Ende 6.2; anon Beginn 5.4, Maximum 5.4, Minimum 5.4, Ende 5.4 MiB
Stufe 100000 narrow, Run 1/2 — Kopierdauer 9687 ms, 10323 Zeilen/s; Proben 39; memory.current Beginn 6.2, bei 25% 10.0, 50% 10.1, 75% 10.8, Ende 8.8, Spitze 11.2 MiB; anon Beginn 5.4, bei 25% 8.8, 50% 9.1, 75% 9.1, Ende 7.8, Spitze 9.5 MiB
Stufe 100000 narrow, Run 1/2 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 8.8 MiB (anon 7.8, file 0.0, kernel 1.0), Prozess RssAnon 7.8 MiB, docker stats 8.8 MiB, VmHWM 24320 KiB, memory.peak seit Start des Feeds 12.5 MiB; im Run: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt)
Feed-Container nicht mehr lesbar: Status exited, Exit 137, OOMKilled true
Stufe 100000 narrow, Run 1/2 — Nachlauf 60 s ohne Schreibzugriff: 
```

## Reihe M — Lauf `20260925T152516Z`

Mutation: `BENCH_MEM_RESET=0`, Stufen 10.000 und 20.000, 1 Run.

```text
Host — Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU, 33362599936 Byte RAM, PostgreSQL-Image postgres:18-alpine, Feed-Image-ID sha256:969b7fad86909255d47fdf1bd5ada33edee15a171d6cc8b8a9849edb3965cf4e
Stufen 10000 20000, 1 Runs je Stufe, Zeilenbreite narrow, Feed-Umgebung '', Feed-Docker-Argumente '', Feed-Image ghcr.io/pt9912/pg-change-feed:dev
Tabelle bench_mem_narrow_10000 — 10000 Zeilen, mittlere Zeilenbreite 73 B (pg_column_size), 63 B (Textform)
Tabelle bench_mem_narrow_20000 — 20000 Zeilen, mittlere Zeilenbreite 73 B (pg_column_size), 65 B (Textform)
Stufe 10000 narrow, Run 1/1 — Zeilen in cdc.change vor dem Run: 0; Grundlinie, Ruhe 10 s nach dem Start: memory.current 6.2 MiB zu Beginn, Maximum 6.2, Minimum 6.2, Ende 6.2; anon Beginn 5.4, Maximum 5.4, Minimum 5.4, Ende 5.4 MiB
Stufe 10000 narrow, Run 1/1 — Kopierdauer 991 ms, 10091 Zeilen/s; Proben 5; memory.current Beginn 7.0, bei 25% 11.2, 50% 10.6, 75% 10.0, Ende 10.0, Spitze 11.2 MiB; anon Beginn 5.4, bei 25% 9.3, 50% 8.9, 75% 9.0, Ende 9.0, Spitze 9.3 MiB
Stufe 10000 narrow, Run 1/1 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 10.0 MiB (anon 9.0, file 0.0, kernel 0.9), Prozess RssAnon 9.0 MiB, docker stats 10.0 MiB, VmHWM 24412 KiB, memory.peak seit Start des Feeds 11.7 MiB; im Run: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt)
Stufe 10000 narrow, Run 1/1 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 10.0 MiB zu Beginn, Maximum 21.4, Minimum 10.0, Ende 19.6; anon Beginn 9.0, Maximum 20.3, Minimum 9.0, Ende 18.7 MiB
Stufe 10000 narrow, Run 1/1 — Zeilen in cdc.change nach dem Nachlauf: 10000, VmHWM 37040 KiB, memory.peak seit Start des Feeds 26.0 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt); seit dem Start des Feeds: Bereinigung gelaufen 13, fehlgeschlagen 0
Slot slot_bench_mem nach 10 s noch aktiv (Sitzung 165: active WalSenderWaitForWal) — Sitzung beendet
Stufe 20000 narrow, Run 1/1 — Zeilen in cdc.change vor dem Run: 10000; Grundlinie, Ruhe 10 s nach dem Start: memory.current 18.4 MiB zu Beginn, Maximum 21.0, Minimum 18.4, Ende 21.0; anon Beginn 17.4, Maximum 20.0, Minimum 17.4, Ende 20.0 MiB
Stufe 20000 narrow, Run 1/1 — Kopierdauer 2086 ms, 9588 Zeilen/s; Proben 9; memory.current Beginn 21.3, bei 25% 13.2, 50% 13.0, 75% 12.8, Ende 10.7, Spitze 21.3 MiB; anon Beginn 20.0, bei 25% 10.7, 50% 11.2, 75% 10.8, Ende 9.2, Spitze 20.0 MiB
Stufe 20000 narrow, Run 1/1 — unmittelbar nach dem Run: cgroup Run-Ende: memory.current 10.2 MiB (anon 9.2, file 0.0, kernel 1.0), Prozess RssAnon 9.2 MiB, docker stats 10.2 MiB, VmHWM 35728 KiB, memory.peak seit Start des Feeds 24.1 MiB; im Run: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt)
Stufe 20000 narrow, Run 1/1 — Nachlauf 60 s ohne Schreibzugriff: Nachlauf: memory.current 10.2 MiB zu Beginn, Maximum 42.7, Minimum 10.1, Ende 38.4; anon Beginn 9.2, Maximum 40.9, Minimum 9.2, Ende 37.3 MiB
Stufe 20000 narrow, Run 1/1 — Zeilen in cdc.change nach dem Nachlauf: 30000, VmHWM 61132 KiB, memory.peak seit Start des Feeds 50.1 MiB, OOMKilled false, Neustarts 0; im Nachlauf: keine GC-Zeilen (GODEBUG=gctrace=1 nicht gesetzt); seit dem Start des Feeds: Bereinigung gelaufen 13, fehlgeschlagen 0
Ende — Lauf 20260925T152516Z
```

## Reihe O — Lauf `20260925T152357Z`

Mutation: Kopie des Skripts mit verändertem cgroup-Pfad, Stufe 10.000.

```text
Host — Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU, 33362599936 Byte RAM, PostgreSQL-Image postgres:18-alpine, Feed-Image-ID sha256:969b7fad86909255d47fdf1bd5ada33edee15a171d6cc8b8a9849edb3965cf4e
Stufen 10000, 1 Runs je Stufe, Zeilenbreite narrow, Feed-Umgebung '', Feed-Docker-Argumente '', Feed-Image ghcr.io/pt9912/pg-change-feed:dev
Tabelle bench_mem_narrow_10000 — 10000 Zeilen, mittlere Zeilenbreite 73 B (pg_column_size), 64 B (Textform)
cgroup (/sys/fs/cgroup/system.slicex/docker-59f54ed666ec434b29eb9d1f1992b771255de78fc5c7436c296ca731da07bf42.scope) oder /proc/3952832 nicht lesbar
```

## Reihe N — Lauf `20260925T153145Z`

`tools/bench-backfill.sh --full`, Default-Image, `BENCH_BACKFILL_RUNS=3`, `BENCH_FEED_DOCKER_ARGS="--memory 10g"`.

```text
Host — Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU, 33362599936 Byte RAM, PostgreSQL-Image postgres:18-alpine
Umgebung wird aufgebaut (Stufen: 10000 100000 1000000 Zeilen, 3 Läufe je Stufe, Blockgröße B=1000 Zeilen, Toleranz 10 min) …
Schätz-Messung — Tabellen mit 100000 Zeilen …
Schätzung — frisch befüllt, Autovacuum an (100000 Zeilen tatsächlich): reltuples=-1, geschätzt unbekannt
Schätzung — frisch befüllt, Autovacuum aus (nie analysiert) (100000 Zeilen tatsächlich): reltuples=-1, geschätzt unbekannt
Schätzung — Autovacuum an, 35 s nach dem Befüllen (100000 Zeilen tatsächlich): reltuples=100000, geschätzt 100000, Abweichung +0.0%
Schätzung — Autovacuum aus, 35 s nach dem Befüllen (nie analysiert) (100000 Zeilen tatsächlich): reltuples=-1, geschätzt unbekannt
Schätzung — nach ANALYZE (100000 Zeilen tatsächlich): reltuples=100000, geschätzt 100000, Abweichung +0.0%
Schätzung — nach ANALYZE und 20000 weiteren Zeilen (+20 %) ohne erneutes ANALYZE (120000 Zeilen tatsächlich): reltuples=100000, geschätzt 100000, Abweichung -16.7%
Grundlinie ohne Run — Feed-Container in Ruhe 30 s nach dem Start 6.1 MiB, Zeilen in cdc.change: 0
Stufe 10000 Zeilen (Zeilenbreite ~73 B gemittelt über pg_column_size, B=1000, Einfügeform zeilenweise in einer Transaktion), Feed-Container in Ruhe 6.1 MiB, Zeilen in cdc.change vor der Stufe: 0
Stufe 10000, Lauf 1/3 — Kopierdauer 1112 ms (finished_at − started_at), 8993 Zeilen/s, geschätzt unbekannt, Warnung Größe f, Warnung Dauer f, Feed-Speicher-Spitze 10.1 MiB, WAL-Rückstand des Slots (confirmed_flush_lsn): Spitze 0 MiB im Run, 0 MiB unmittelbar danach, 0 MiB 0 s danach ohne Schreibzugriff; vom Slot gehaltenes WAL (restart_lsn): Spitze 6 MiB im Run
Stufe 10000, Lauf 2/3 — Kopierdauer 1061 ms (finished_at − started_at), 9425 Zeilen/s, geschätzt unbekannt, Warnung Größe f, Warnung Dauer f, Feed-Speicher-Spitze 25.5 MiB, WAL-Rückstand des Slots (confirmed_flush_lsn): Spitze 0 MiB im Run, 0 MiB unmittelbar danach, 0 MiB 0 s danach ohne Schreibzugriff; vom Slot gehaltenes WAL (restart_lsn): Spitze 6 MiB im Run
Stufe 10000, Lauf 3/3 — Kopierdauer 982 ms (finished_at − started_at), 10183 Zeilen/s, geschätzt unbekannt, Warnung Größe f, Warnung Dauer f, Feed-Speicher-Spitze 10.3 MiB, WAL-Rückstand des Slots (confirmed_flush_lsn): Spitze 0 MiB im Run, 0 MiB unmittelbar danach, 0 MiB 0 s danach ohne Schreibzugriff; vom Slot gehaltenes WAL (restart_lsn): Spitze 7 MiB im Run
Stufe 10000 Ergebnis — Kopierdauer 982–1112 (Median 1061, n=3) ms, Durchsatz 8993–10183 (Median 9425, n=3) Zeilen/s, Feed-Speicher-Spitze 25.5 MiB (Ruhe vor der Stufe 6.1 MiB, 20 s nach dem letzten Lauf 37.8 MiB, 60 s nach dem letzten Lauf 36.1 MiB, Zeilen in cdc.change danach: 30000), mittlere Blockdauer 106 ms (abgeleitet: Median-Dauer / 10 Blöcke; die längste Einzeldauer ist nicht gemessen), WAL-Rückstand-Spitze im Run Median 0 MiB (~0 B je Zeile, abgeleitet), vom Slot gehaltenes WAL Median 6 MiB
Stufe 100000 Zeilen (Zeilenbreite ~73 B gemittelt über pg_column_size, B=1000, Einfügeform zeilenweise in einer Transaktion), Feed-Container in Ruhe 41.7 MiB, Zeilen in cdc.change vor der Stufe: 30000
Stufe 100000, Lauf 1/3 — Kopierdauer 9895 ms (finished_at − started_at), 10106 Zeilen/s, geschätzt unbekannt, Warnung Größe f, Warnung Dauer f, Feed-Speicher-Spitze 40.9 MiB, WAL-Rückstand des Slots (confirmed_flush_lsn): Spitze 0 MiB im Run, 0 MiB unmittelbar danach, 0 MiB 0 s danach ohne Schreibzugriff; vom Slot gehaltenes WAL (restart_lsn): Spitze 63 MiB im Run
Stufe 100000, Lauf 2/3 — Kopierdauer 10187 ms (finished_at − started_at), 9816 Zeilen/s, geschätzt unbekannt, Warnung Größe f, Warnung Dauer f, Feed-Speicher-Spitze 65.3 MiB, WAL-Rückstand des Slots (confirmed_flush_lsn): Spitze 0 MiB im Run, 0 MiB unmittelbar danach, 0 MiB 0 s danach ohne Schreibzugriff; vom Slot gehaltenes WAL (restart_lsn): Spitze 63 MiB im Run
Stufe 100000, Lauf 3/3 — Kopierdauer 10284 ms (finished_at − started_at), 9724 Zeilen/s, geschätzt unbekannt, Warnung Größe f, Warnung Dauer f, Feed-Speicher-Spitze 231.9 MiB, WAL-Rückstand des Slots (confirmed_flush_lsn): Spitze 0 MiB im Run, 0 MiB unmittelbar danach, 0 MiB 0 s danach ohne Schreibzugriff; vom Slot gehaltenes WAL (restart_lsn): Spitze 70 MiB im Run
Stufe 100000 Ergebnis — Kopierdauer 9895–10284 (Median 10187, n=3) ms, Durchsatz 9724–10106 (Median 9816, n=3) Zeilen/s, Feed-Speicher-Spitze 231.9 MiB (Ruhe vor der Stufe 41.7 MiB, 20 s nach dem letzten Lauf 349.6 MiB, 60 s nach dem letzten Lauf 243.6 MiB, Zeilen in cdc.change danach: 330000), mittlere Blockdauer 102 ms (abgeleitet: Median-Dauer / 100 Blöcke; die längste Einzeldauer ist nicht gemessen), WAL-Rückstand-Spitze im Run Median 0 MiB (~0 B je Zeile, abgeleitet), vom Slot gehaltenes WAL Median 63 MiB
Stufe 1000000 Zeilen (Zeilenbreite ~74 B gemittelt über pg_column_size, B=1000, Einfügeform zeilenweise in einer Transaktion), Feed-Container in Ruhe 275.7 MiB, Zeilen in cdc.change vor der Stufe: 330000
Stufe 1000000, Lauf 1/3 — Kopierdauer 133442 ms (finished_at − started_at), 7494 Zeilen/s, geschätzt unbekannt, Warnung Größe f, Warnung Dauer f, Feed-Speicher-Spitze 439.7 MiB, WAL-Rückstand des Slots (confirmed_flush_lsn): Spitze 1 MiB im Run, 0 MiB unmittelbar danach, 0 MiB 0 s danach ohne Schreibzugriff; vom Slot gehaltenes WAL (restart_lsn): Spitze 631 MiB im Run
Stufe 1000000, Lauf 2/3 — Kopierdauer 129620 ms (finished_at − started_at), 7715 Zeilen/s, geschätzt unbekannt, Warnung Größe f, Warnung Dauer f, Feed-Speicher-Spitze 439.6 MiB, WAL-Rückstand des Slots (confirmed_flush_lsn): Spitze 1 MiB im Run, 0 MiB unmittelbar danach, 0 MiB 0 s danach ohne Schreibzugriff; vom Slot gehaltenes WAL (restart_lsn): Spitze 1019 MiB im Run
Stufe 1000000, Lauf 3/3 — Kopierdauer 109110 ms (finished_at − started_at), 9165 Zeilen/s, geschätzt unbekannt, Warnung Größe f, Warnung Dauer f, Feed-Speicher-Spitze 2738.2 MiB, WAL-Rückstand des Slots (confirmed_flush_lsn): Spitze 2 MiB im Run, 0 MiB unmittelbar danach, 0 MiB 0 s danach ohne Schreibzugriff; vom Slot gehaltenes WAL (restart_lsn): Spitze 1030 MiB im Run
Stufe 1000000 Ergebnis — Kopierdauer 109110–133442 (Median 129620, n=3) ms, Durchsatz 7494–9165 (Median 7715, n=3) Zeilen/s, Feed-Speicher-Spitze 2738.2 MiB (Ruhe vor der Stufe 275.7 MiB, 20 s nach dem letzten Lauf 2737.2 MiB, 60 s nach dem letzten Lauf 2448.4 MiB, Zeilen in cdc.change danach: 3330000), mittlere Blockdauer 130 ms (abgeleitet: Median-Dauer / 1000 Blöcke; die längste Einzeldauer ist nicht gemessen), WAL-Rückstand-Spitze im Run Median 1 MiB (~1 B je Zeile, abgeleitet), vom Slot gehaltenes WAL Median 1019 MiB
WAL-Rückstand der größten Stufe (1000000 Zeilen) — höchste Spitze im Run 2 MiB, unter der Warnschwelle von 100 MiB (SPEC-013)
Live-Erfassung — 100 Zeilen/s in public.bench_backfill_live, Run über Stufe 1000000 …
Live-Erfassung, Referenz ohne Run (114 s) …
cdc_capture_lag während des Runs (114 s, Stufe 1000000, Live-Last 100/s): 0.028904–1.033201 (Median 0.533628, n=93) s; in den 10 s danach: 0.129317–0.644176 (Median 0.344454, n=10) s; Referenz ohne Run: 0.028792–1.030879 (Median 0.497159, n=114) s (Live-Zeilen eingefügt: 13500 im Lauf mit Run, 12100 in der Referenz); im Lauf mit Run: WAL-Rückstand des Slots (confirmed_flush_lsn) Spitze 1 MiB, vom Slot gehaltenes WAL (restart_lsn) Spitze 740 MiB
Richtgröße (abgeleitet) — 7715 Zeilen/s (Median, Stufe 1000000) × Toleranz 600 s = 4629000 Zeilen, abgerundet auf eine Stelle: 4000000 Zeilen (Startwert der Toleranz: Setzung ohne Messung; die Rate ist über die Stufe hinaus hochgerechnet)
Ende — Lauf 20260925T153145Z
```
