# Slice backfill-speicher-untersuchung: Speicher-Untersuchung des Backfills — Spitze des Feed-Containers über dem Blockbedarf, Ursache und gemessene Grenze

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — der Slice trägt keine Closure-Bedingung, die von
seiner DoD verschieden wäre. Er hat keine Kante zu einer offenen Welle und ist
von den Slices der Welle [welle-transformationen](../welle-transformationen.md)
unabhängig.

**Bezug:** [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Backfill des
Bestands),
[`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md)
(Bestand-Backfill; sein Re-Evaluierungs-Trigger zur Ausbaustufe),
[`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
(Warnkriterium mit der Richtgröße von 4.000.000 geschätzten Zeilen),
[`ADR-0054`](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md) §(b)
(Benchmark-Infrastruktur), Architect-Verdikt
[`architect-verdict-welle-backfill-bestand-lese-schritt`](../../../reviews/architect-verdict-welle-backfill-bestand-lese-schritt.md)
§5 (h).

**Berührte Spec-Stellen:** [`SPEC-005`](../../../../spec/pflichtenheft.md)
(TransactionBufferPort; der Satz in `spec/pflichtenheft.md`, Zeile 72, gelesen
am 2026-09-25: „Große Transaktionen dürfen nicht unbegrenzt im RAM gehalten
werden“) — gelesen, nicht geändert.

**Verantwortlich:** Implementer-Agent, 2026-09-25.

**Autor:** Planner-Agent, Closure der Welle
[welle-backfill-bestand](../done/welle-backfill-bestand.md). **Datum:** 2026-09-25.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Ursache der Speicher-Spitze des Feed-Containers im Backfill-Run
ist benannt, und ihre Abhängigkeit von der Tabellengröße ist gemessen. Das
Benutzerhandbuch nennt unter „Grenzwerte“ (Speicher des Feed-Containers) den
Stand *übernommen* aus einem Lauf, der im Repository nicht auflösbar ist:
Stufe mit 200.000 Zeilen, Spitze im Run 401,7 bis 467 MiB bei 176,3 bis 196,0
MiB in Ruhe davor, höchster gemessener Wert 641,7 MiB; zwei Runs über je
1.000.000 Zeilen mit 435 MiB und **1.544 MiB**; der Bedarf eines Blocks liegt
bei etwa 74 KB, „die Ursache ist nicht untersucht, ein Zusammenhang mit der
Tabellengröße ist nicht belegt“ (Handbuch §9, Zeilen 1653 bis 1668, übernommen).
Die Warn-Richtgröße steht bei 4.000.000 geschätzten Zeilen
(`internal/application/usecase/backfill/warn.go`, gemessen am 2026-09-25); eine
Richtgröße, die auf einer ungeklärten Speicherlage steht, ist ein Risiko für den
ersten Betreiber.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Änderung am Backfill-Code** (Bytelimit, Zeilenbreiten-Wache, andere
  Blockgröße). Der Slice ist eine Untersuchung; ergibt sie ein Wachstum mit der
  Tabellengröße, folgen ein Befund und ein eigener Änderungs-Slice mit
  Entscheidung des Architects ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  Re-Evaluierung: Ausbaustufe für Durchsatz).
- **Ein Gate oder eine Pass/Fail-Schwelle für den Speicher.** Die Messung ist
  ohne Schwelle geführt (`harness/targets/bench-backfill.md`); eine Schwelle
  brauchte eine ADR ([`AGENTS.md`](../../../../AGENTS.md) §3.6).
- **Die Kopierdauer und der Durchsatz** — Gegenstand der bestehenden Messung;
  der Slice liest sie mit, ändert sie nicht.
- **Breite Zeilen als eigene Messreihe**, sofern die Ursache nicht an der
  Zeilenbreite hängt: die Zeilenbreite (etwa 74 Bytes) ist im Ursprung jeder
  Kopier-Zahl genannt, breitere Zeilen sind ungemessen und bleiben so, bis die
  Ursache einen Zusammenhang zeigt.

## 2. Definition of Done

- [x] Die Messreihe steht mit gedruckten Zeilen im Repository: `tools/bench-backfill.sh`
      (vorhanden, Vertrag `harness/targets/bench-backfill.md`) liefert je Stufe
      den Speicher des Feed-Containers **und** eine Grundlinie ohne Run
      (Ruhe-Speicher vor und nach), über mindestens die Stufen bis 1.000.000
      Zeilen (`--full`); jede Zahl trägt Host, Lauf und die gedruckte Zeile
      (`AGENTS.md` §3.12 Instanz A), der Bericht trennt Messung von Übernahme.
      *Zu belegen durch:* die gedruckten Zeilen im Bericht des Slice unter
      `docs/reviews/` (auflösbar, anders als der übernommene Lauf
      `20260924T233628Z`); der Lauf, aus dem die 1.544 MiB stammen, wird
      wiederholt oder bleibt als *übernommen* gekennzeichnet.
- [x] Die Ursache ist benannt und belegt: ein Lauf, der sie an- und abschaltet
      (eine Änderung genau einer Größe — Blockgröße, Zeilenbreite,
      Garbage-Collector-Einstellung des Feed-Containers oder der Zeitpunkt der
      Probe —, dieselbe Messung vor und nach), oder die Ursache steht als
      *ungeklärt* mit der gemessenen Grenze im Bericht. Der Blockbedarf von etwa
      74 KB (Handbuch, übernommen) erklärt Spitzen von einigen hundert MiB nicht;
      welche Größe sie trägt, ist offen und im Slice zu erproben.
- [x] Der Ausgang ist gesetzt: wächst der Speicher mit der Tabellengröße, steht
      ein Befund im Bericht und ein Änderungs-Slice existiert als Datei in `open/`
      (Entscheidung des Architects, ob eine ADR nötig ist); sonst trägt das
      Benutzerhandbuch unter „Grenzwerte“ die gemessene Grenze samt Ursprung,
      und die Warn-Richtgröße von 4.000.000 geschätzten Zeilen ist gegen die
      Messung bewertet (bleibt, oder eine Nachschärfung nach dem Trigger von
      [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
      „Eine Messung liegt vor“ ist beantragt); die Handbuch-Version und die
      Änderungshistorie tragen eine Zeile. *Zu belegen durch:* Lesen des
      Handbuch-Abschnitts und `make docs-check`.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (die Fixrunde löst F-1 bis F-5; `review-slice-backfill-speicher-untersuchung`)
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: siehe dritter Liefer-Punkt (Handbuch §Grenzwerte,
      `harness/targets/bench-backfill.md`, falls der Vertrag der Messung eine
      Grundlinie führt).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure der nächsten Welle (die Roadmap führt
      [welle-transformationen](../welle-transformationen.md) unter *Offene
      Wellen*, das Ereignis kann eintreten; ein Slice ohne Welle wird von ihr
      mitgeprüft).

**Umfang:** M — Schätzung, nicht gemessen: die Messreihe fährt die 1.000.000-Stufe
(Laufzeit und Speicherbedarf des Messhosts sind ungemessen), die Ursachensuche
hat kein bekanntes Ende.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/bench-backfill.sh` (falls die Grundlinie fehlt) | update | Speicher ohne Run vor und nach der Stufe, im Ursprung jeder Zahl genannt. Geliefert: Grundlinie 30 s nach dem Start, Speicher 20 s und 60 s nach der Stufe, Zahl der Zeilen in `cdc.change` vor und nach der Stufe. |
| `harness/targets/bench-backfill.md` | update | Vertrag der Messung nennt die Grundlinie; Abschnitt zum zweiten Skript; Grenzen „Speicher“ und „Zeilenbreite“ nachgezogen. |
| `docs/user/benutzerhandbuch.md` (§Grenzwerte, Version, Änderungshistorie) | update | gemessene Grenze samt Ursprung statt der übernommenen Zahlen. Geliefert: Speicher-Bullet ersetzt (Ursache, Tabelle, Bemessung), Richtgröße-Bullet um den Speicher ergänzt, Vorbedingungs-Absatz in §4 „Bestand als Backfill überführen“ und Absatz in „Aufbewahrung (Retention)“ (Träger derselben Aussage), Version 1.60. |
| Bericht des Slice unter `docs/reviews/` | neu | gedruckte Zeilen der Messreihe, Ursache, Ausgang. Geliefert als `messbericht-slice-backfill-speicher-untersuchung.md` (Bericht) und `messbericht-slice-backfill-speicher-untersuchung-zeilen.md` (gedruckte Zeilen; 21 Überschriften `Reihe`/`Reihen <Buchstabe>` am Stand `989beef3`, gemessen mit `grep -c -E '^#{2,3} Reihen? [A-Z]'`; die Buchstaben reichen von A bis R). |
| `tools/bench-backfill-memory.sh` | neu | **über den Plan hinaus:** die Trennung von Run, Nachlauf und Zahl der Changes braucht einen Container je Run und die cgroup-Zähler vom Host aus; `tools/bench-backfill.sh` liefert das nicht, ohne seine Messung zu ändern. Kein Produktionscode nötig (Zähler von außen, `GODEBUG=gctrace=1`). |
| `tools/bench-lib.sh` | update | **über den Plan hinaus:** `BENCH_FEED_ENV` und `BENCH_FEED_DOCKER_ARGS` für `bench::start_feed` (Umgebungsvariablen und `docker run`-Argumente des Feed-Containers, Speichergrenze). |
| `harness/README.md` (Zeile `make bench`) | update | **über den Plan hinaus:** nennt das zweite Skript als Werkzeug außerhalb von `make bench`. |
| `docs/plan/planning/open/slice-retention-lauf-speicher-begrenzung.md` | neu | **Ausgang** (DoD, dritter Punkt): der Änderungs-Slice der Ursache, Datei in `open/`. |
| `docs/plan/planning/observations/BEO-PGC/blockgroesse-zaehlt-zeilen-nicht-bytes/state.md` | update | nur der Verweis auf den Slice: ein Link mit festem Lifecycle-Verzeichnis (`open/`) bricht mit dem `git mv`; er steht jetzt als Zitat der Kennung. Der Inhalt des Registers bleibt der Planner-Closure. |

**Fixrunde** (Findings des Reviews, alle **über den Plan hinaus**; Zeilen der Läufe im
Zeilen-Dokument, Reihen P1 bis R):

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/reviews/messbericht-slice-backfill-speicher-untersuchung.md` | update | F-1 (`GOGC=25` senkt um 11 bis 21 %, drei Paare genannt), F-2 (Höchstwert je Change 1,59 statt 1,57; Zeile der Runs 3 mit 2.000.000 Changes davor, Divisor genannt; Folgezahlen 3,9 bis 6,1 GiB und 181 bis 279 MiB), F-5 (§7: Reichweite `v0.1.0` bis `v0.1.2`, Ausgang „Änderungs-Slice zuerst, danach `v0.2.0`“ mit Anker), F-9 und F-11 (Herkunft der Live-Aussage, Reihe N als Gegenprobe, Hypothese des Architect-Verdikts), F-12 (Tag), Abschnitt 8 um die Mutationen der Fixrunde |
| `docs/reviews/messbericht-slice-backfill-speicher-untersuchung-zeilen.md` | update | die gedruckten Zeilen der Läufe der Fixrunde (Reihen P1, P2, P3, Q1, Q2, R) |
| `docs/user/benutzerhandbuch.md` | update | F-1, F-2 (Zeile und Höchstwert), F-3 (§4: Speicher im Run bei nicht leerem `cdc.change`), F-7 (zwei Sätze im Ist-Zustand), F-9 (Ursprung der Live-Aussage), F-11/F-12 (Reihe N, `v0.1.0` bis `v0.1.2`); Version 1.61 mit Historienzeile; die Zeile 1.60 trägt den korrigierten Höchstwert |
| `tools/bench-backfill.sh` | update | F-4: `feed_running_or_report` beendet die Warteschleife des Runs und die der Live-Phase mit dem Zustand des Feed-Containers |
| `tools/bench-backfill-memory.sh` | update | F-6 (`gc_summary` unterscheidet gesetztes und ungesetztes `GODEBUG=gctrace=1`), F-10 (die Zeitreihen-Datei eines Runs liegt unter `$WINDOW_FILE` und wird von der Falle entfernt) |
| `harness/targets/bench-backfill.md` | update | F-4 (der Satz „Kill ist ein Messergebnis“ nennt die zwei Skripte, die ihn tragen), F-13 (cgroup-v2 mit systemd-Treiber als Vorbedingung) |

Fixrunde, nicht realisiert: F-10 (b) — die Meldung „Zeile 58 …/memory.current: Datei oder
Verzeichnis nicht gefunden“ ist in der Fassung des Reviews nicht reproduzierbar (Shell-Test der
Konstruktion `{ S=$(<Datei); } 2>/dev/null` gegen eine fehlende Datei: keine Ausgabe;
Lauf `20260925T165631Z` mit `--memory 64m`: keine solche Zeile), keine Änderung. F-8 (Ursache und
Zahl im Plan des Folge-Slice) und der Suchbefehl mit `1,03 bis 1,57` im Plan des Folge-Slice
(Zeile 293) gehören in dessen Plan: **gemeldet**, nicht mitgeändert (Implementer); F-8 ist
im Plan des Folge-Slice gezogen, der Suchbefehl (`1,03 bis 1,5[79]`) und sein Suchlauf-Feld
sind in der Closure berichtigt (Planner). F-14: keine Aktion (Prozess-Hinweis, kein Rückstand).

Nicht realisiert: nichts vom Plan. Die Eingrenzung „Breite Zeilen als eigene Messreihe“
(§1) wurde ausgeübt, weil die Ursache an der Zeilenbreite hängt (Reihen H und I); der Plan
nannte diese Bedingung selbst.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Aussage über den
Speicher des Feed-Containers im Backfill und die Warn-Richtgröße“; beide Stände
gemessen; die Befehle stehen im Codeblock, der Implementer trägt Stand und
Trefferzahl ein):**

```text
git grep -n -i -E 'Speicher|MiB|1\.544|4\.000\.000|Richtgröße|nicht untersucht' -- docs/user harness spec internal tools
git grep -n -E 'estimatedRowsGuideline|DefaultBlockSize' -- internal
```

**Stand und Trefferzahl (gemessen mit den zwei Befehlen oben; Parent =
`493a28ad`, Diff = Arbeitsbaum samt Index dieses Slice, neue Dateien mit `git add`
aufgenommen; die Diff-Zahlen sind an zwei Ständen gemessen, weil die Fixrunde das Handbuch
weiter verändert).** Befehl 1: Parent 242 Zeilen (`docs/user` 58, `harness` 13, `spec` 36,
`internal` 111, `tools` 24); Stand `4f94f900` 280 (71, 23, 36, 111, 39); Stand `989beef3` 288
(`docs/user` 79, `harness` 23, `spec` 36, `internal` 111, `tools` 39; gemessen mit `git grep
-c` je Wurzel am 2026-09-25). Befehl 2: Parent 11, Stand `989beef3` 11.
Der zusätzliche Lauf auf die alten Zahlen und die Aussage „nicht untersucht“ (`git grep -n
-E '1\.544|641,7|467 MiB|401,7|416,5' -- docs/user harness` und
`git grep -n -i 'nicht untersucht' -- docs/user harness spec internal tools`): Parent 7
und 1 Treffer, Stand `989beef3` 2 und 0.

| Träger | Befund | Behandlung |
|---|---|---|
| Handbuch §9 „Grenzwerte“, Bullet Speicher des Feed-Containers | Parent: Zahlen 467, 401,7, 416,5, 641,7, 435 und 1.544 MiB, alle übernommen oder aus Reviews, dazu „nicht untersucht“ und „Zusammenhang mit der Tabellengröße nicht belegt“ | ersetzt durch gemessene Zahlen mit Reihe und Lauf; die eine übernommene Zahl (1.544 MiB) steht als übernommen mit dem vereinbaren Messwert (1.538,7 MiB) |
| Handbuch §9, Bullet Richtgröße | Parent: Aussage „Breite Zeilen … ungemessen“ ohne Bezug zum Speicher | um den Speicher des Feed-Containers ergänzt, „ungemessen“ auf die Kopierrate eingegrenzt |
| Handbuch §4 „Bestand als Backfill überführen“ und „Aufbewahrung (Retention)“ | beschrieben die Bereinigung und die Vorbedingungen ohne den Speicher (Parent: keine Treffer auf „Speicher des Feed“ dort) | je ein Absatz ergänzt; das Zählwort „vier Betriebs-Vorbedingungen“ bleibt, der Speicher steht als eigener Absatz außerhalb der Liste |
| `harness/targets/bench-backfill.md` | Ablauf Nr. 2 ohne Grundlinie; Grenze „Speicher“ nur `docker stats`; Grenze „Zeilenbreite“ „breite Zeilen ungemessen“ | Ablauf, Abschnitt zum zweiten Skript, beide Grenzen nachgezogen |
| `harness/README.md` (Zeile `make bench`) | nennt die Speichermessung ohne das zweite Skript | Verweis ergänzt |
| Doc-Kommentare an `DefaultBlockSize`, `estimatedRowsGuideline` (Befehl 2: 11 Treffer an Definition und Verwendungen, Parent und Diff gleich) sowie an `BackfillTableService` und `NextBlock` (gelesen) | tragen „Zeilen, nicht Bytes“ und „Richtgröße aus Rate mal Toleranz“; keine Aussage über den Speicher des Containers | **unverändert**: beide Aussagen bleiben wahr (Messbericht Abschnitt 3.4 und 6) |
| `internal/bootstrap/wiring.go`, Kommentar an `retentionInterval` | nennt „Jeder Takt liest alle Changes der Quelle“ — die Ursache, ohne Bezug zum Speicher | **unverändert**: wahr; der Änderungs-Slice trägt den Kommentar |
| `docs/plan/planning/done/**`, `docs/reviews/**` (Architect-Verdikt, Review-Reports, Results-Notiz der Welle) | tragen die übernommenen Zahlen und „nicht untersucht“ (Treffer im Lauf auf `docs`, außerhalb der zwei Befehle) | **nicht geändert**: Records ([`AGENTS.md`](../../../../AGENTS.md) §3.5) |
| Register `BEO-PGC/blockgroesse-zaehlt-zeilen-nicht-bytes/state.md` | Feld „Gegenstand“ trägt „die Ursache ist nicht untersucht“ | **gemeldet**, nicht mitgeändert (Planner-Closure); nur der gebrochene Link ist ein Zitat |
| `spec/pflichtenheft.md` Zeile 72 (`SPEC-005`) | „Große Transaktionen dürfen nicht unbegrenzt im RAM gehalten werden“ | gelesen: betrifft den Transaktionspuffer des Capture-Pfads, nicht die Retention-Lesung; nicht geändert |
| **Nicht gefunden** | kein Träger außerhalb der genannten beschreibt die Retention-Lesung als Speicherquelle; kein Träger in `spec/` und `docs/user/` (außer dem Handbuch) trägt Zahlen zum Speicher des Feed-Containers | — |

**§3.13-Suchlauf der Fixrunde** (bewegte Eigenschaften: der Höchstwert je Change und die
Folgezahlen, die Prozentangabe zu `GOGC`, die Aussage „Speicher nach dem Run, nicht im Run“,
der Vertragssatz zum Kill des Feed-Containers, die Empfehlung zum Server-Release; Suchform
nach [`AGENTS.md`](../../../../AGENTS.md) §3.13: Symbolname, Zählwort, Beschreibung; Träger:
`docs/user`, `harness`, `tools`, die zwei Messberichte und `docs/plan/planning/open`; Parent =
`5b1f7762`, Stand vor der Fixrunde, mit `git grep <Befehl> 5b1f7762 -- <Wurzeln>`; Diff =
Arbeitsbaum der Fixrunde ohne diesen Plan):

```text
git grep -n -E '1,57|1[.]57|1,51 |1,35 |6,0 GiB|18 bis 21' <Stand> -- docs/user harness tools <Messberichte> docs/plan/planning/open
git grep -n -E 'nach dem Run, nicht|nicht während er|Kill durch die Grenze|jedes Skript' <Stand> -- docs/user harness tools <Messberichte> docs/plan/planning/open
git grep -n -E 'vor dem Release|vor einem Server-Release|bekannte, gemessene|Änderungs-Slice vor' <Stand> -- docs/user harness tools <Messberichte> docs/plan/planning/open
git grep -n -E 'feed_mem_mib|feed_running_or_report|BENCH_FEED_DOCKER_ARGS|gc_summary' <Stand> -- docs/user harness tools <Messberichte> docs/plan/planning/open
git grep -n -E 'Reihen? [A-Z] (bis|und)' <Stand> -- docs/user harness tools <Messberichte> docs/plan/planning/open
```

`<Messberichte>` sind `docs/reviews/messbericht-slice-backfill-speicher-untersuchung.md` und
`docs/reviews/messbericht-slice-backfill-speicher-untersuchung-zeilen.md`. **Trefferzahl**
(gemessen mit den fünf Befehlen): Befehl 1 Parent 16, Diff 7; Befehl 2 Parent 2, Diff 1;
Befehl 3 Parent 2, Diff 0; Befehl 4 Parent 17, Diff 27; Befehl 5 Parent 8, Diff 13.

**Gefunden und gezogen:** die Treffer des Parent bei Befehl 1 im Handbuch (§9, Historienzeile
1.60), im Messbericht (§1, §3.1, §3.3, §3.5, §3.6, §4, §6) und im Wert der Runs 3; bei Befehl 2 der
Handbuch-Absatz „nach dem Run, nicht während er läuft“ (§4); bei Befehl 3 die zwei Zeilen der
Empfehlung in §7 des Messberichts. **Gefunden und nicht gezogen:** die
sieben Treffer des Befehls 1 im Diff sind wahr — zwei Zeilen der Tabelle in §3.3 des Messberichts
(1,51 aus Reihe D, 1,57 aus Reihe C, Run 2), zwei Treffer auf „16,0 GiB“ (der Suchausdruck `6,0 GiB`
trifft die Endung von „16,0 GiB“), die Versionszeile `1.57` und die Historienzeile 1.61 des
Handbuchs (nennt den früheren Wert als Gegenstand der Korrektur), und die Zeile 293 des Plans
des Folge-Slice (Stand `3405c9a1`; sie nennt `1,03 bis 1,57` als Suchausdruck; **gemeldet**, Träger in fremder
Datei: der Ausdruck trifft im Handbuch nicht mehr, dort steht `1,03 bis 1,59`; in der Closure
auf `1,03 bis 1,5[79]` berichtigt); Befehl 2 im Diff: der Satz „für jedes Skript, das
`bench::start_feed` ruft“ im Vertrag gilt für die zwei Umgebungsvariablen und bleibt wahr.
**Nicht gefunden** (innerhalb der Suchwurzeln der fünf Befehle): kein weiterer Träger in
`docs/user`, `harness`, `tools`, den zwei Messberichten und `docs/plan/planning/open` nennt
den Höchstwert 1,57 (ohne die Selbstverweise der Pläne, die den früheren Wert als Gegenstand
der Berichtigung nennen), die Prozentangabe oder die Reihenfolge „Release vor Änderungs-Slice“;
`spec/` und `README.md` tragen keine dieser Zahlen. **Außerhalb der Suchwurzeln gefunden**
(Nachmessung der Verifikation V-1 mit `git grep -n -E '1,03 bis 1,57' <Stand>` über den ganzen
Baum, Stand `989beef3`): `docs/plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md`
(Zeilen 68, 191, 199, `Accepted`), das Architect-Verdikt
`docs/reviews/architect-verdict-retention-lauf-speicher-begrenzung.md` (Zeilen 152, 155, 167,
Record) und die Zitate des Befunds in Review- und Verifikations-Report (Records). **Behandlung:**
keine Änderung an ADR, Verdikt und Reports; die Entscheidung trägt mit der Abweichung 1,57
gegen 1,59 (Höchstwert je Change, Reihe B Run 3: 3.096,6 MiB × 1.024 / 2.000.000 = 1,585 KiB,
abgeleitet) unverändert — der Faktor der Kandidaten-Größe („sieben bis zehn“ gegen etwa 0,15 KiB)
wird 6,9 bis 10,6, die Bedarfsrechnung 864.000 × 1,03 bis 1,59 KiB ergibt 0,85 bis 1,31 GiB
(abgeleitet, Verifikation V-1) und keine Festlegung der ADR hängt an der zweiten Stelle der
Zahl; deshalb keine Berichtigungs-ADR. Bei Bedarf schreibt die Closure des Folge-Slice
`slice-retention-lauf-speicher-begrenzung` die Berichtigung in ihre Suchlauf-Zeile. Der Träger
`state.md` des Registers und der Plan des Folge-Slice (Trigger, Zahl „3.000.000“, Ursache des
Ausbleibens der Takte) sind in der Closure nachgezogen.

## 4. Trigger

**Start** (`next` → `in-progress`): kein weiterer Slice in `in-progress/`
(WIP-Limit 1); sonst unabhängig von jedem anderen Slice. Der Slice muss `done`
sein, **bevor** ein Server-Release veröffentlicht wird, dessen Commit den
Backfill trägt — beobachtbar: `git tag -l 'v*'` nennt ein `v*`-Tag über dem
höchsten Tag `v0.1.2`, dessen Commit die Backfill-Änderungen enthält (gemessen am
2026-09-25: kein Server-Tag trägt sie). Ein mechanischer Wächter für diese
Bedingung existiert nicht; der Träger ist dieser Start-Trigger und die
Validator-Feststellung der Results-Notiz von `welle-backfill-bestand`
(Benannte Grenze, §6).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Messreihe und
  Ursachensuche nicht in einem Review tragen — der abtrennbare Teil ist die
  Messreihe (erster Liefer-Punkt) als eigener Slice, die Ursachensuche folgt.
- `in-progress` → `open` (blockiert): falls die Messreihe an der Kapazität des
  Messhosts scheitert (1.000.000 Zeilen brauchen Speicher und Zeit; kein
  Ersatz-Host verfügbar) — dann steht die Grenze der Messung im Bericht statt
  einer Ursache.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + die Messreihe mit gedruckten Zeilen im
Bericht + Ausgang gesetzt (Befund und Änderungs-Slice als Datei in `open/`,
oder Handbuch-Grenze) + Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Die Messung hängt am Host** (Speicher und Laufzeit des Messhosts). *Erwartet,
  zu belegen durch:* jede
  Zahl nennt Host und Lauf; das DoD-Kriterium „Ursache benannt“ ist hostunabhängig
  gefasst (`BEO-PGC/dod-kriterium-haengt-am-messhost`, 1×, offen).
  **Ausgang:** *weiter offen* → Register `BEO-PGC/dod-kriterium-haengt-am-messhost`
  (1×, offen; der Eintrag trägt das Risiko „Messung hängt am Host“) und das Risiko „Die
  Nachmessung hängt am Host“ im Folge-Slice `slice-retention-lauf-speicher-begrenzung`
  (§6). Beleg: jede Zahl trägt Host, Reihe und Lauf und ist als gemessen, abgeleitet oder
  übernommen gekennzeichnet (Verifikation §7); die Ursache ist hostunabhängig belegt
  (Schalter „Bereinigung aus“, Code-Lesung an drei Tags, Verifikation §2 Zeile 2 und §5);
  offen bleibt die Streuung: ein Host, `n` = 2 bis 3 je Bedingung (Messbericht §9). Das
  Kriterium hat keine zweite Ausprägung der Klasse geliefert, der Zähler des Eintrags
  bleibt 1×.
- **Die Ursache bleibt ungeklärt.** *Erwartet, zu belegen durch:* der Ausgang
  „ungeklärt mit gemessener Grenze“ ist zulässig und steht im Bericht und im
  Handbuch; er ist kein stiller Abschluss. **Ausgang:** *entfallen* — die Ursache ist
  benannt und belegt: der Retention-Lauf liest alle Changes der Quelle samt Row Images;
  Bereinigung aus hält `memory.peak` bei 13,3, 13,2 und 13,7 MiB (1.000.000, 2.000.000,
  3.000.000 Changes, Reihe J) gegen 1.082,7 MiB mit Bereinigung (Reihe B, Run 1;
  Verifikation §2 Zeile 2). Ein Teilbefund bleibt **benannte Grenze**, kein Risiko dieses
  Slice: das Ausbleiben der Bereinigungs-Takte ab 2.000.000 Changes ist beobachtet und
  nicht erklärt (Messbericht §9, Handbuch „Grenzwerte“); Adresse: Folge-Slice
  `slice-retention-lauf-speicher-begrenzung` §6, zweiter Punkt (Nachmessung bei 1.000.000,
  2.000.000 und 3.000.000 Changes).
- **Der Slice läuft nicht vor dem ersten Server-Release mit Backfill** (kein
  Wächter). *Erwartet, zu belegen durch:* der Start-Trigger (§4) und die
  Results-Notiz von `welle-backfill-bestand` nennen die Bedingung; die Prüfung liegt
  beim Planner der Release-Vorbereitung. **Ausgang:** *entfallen* — die Bedingung ist
  eingehalten: `git tag -l 'v*'` nennt `v0.1.0`, `v0.1.1` und `v0.1.2` (gemessen am
  2026-09-25), kein Server-Tag über `v0.1.2`, und dieser Slice liegt in `done/`. Die
  Reihenfolge „Änderungs-Slice vor `v0.2.0`“ bleibt ohne mechanischen Wächter; ihr Träger
  ist die Vorbedingung im Folge-Slice `slice-retention-lauf-speicher-begrenzung` (§4) und
  das Risiko dort (§6, letzter Punkt); Adresse der Prüfung: der Planner der
  Release-Vorbereitung.
- **Die Warn-Richtgröße bleibt unbegründet, weil kein Wachstum gemessen wird**
  (`BEO-PGC/blockgroesse-zaehlt-zeilen-nicht-bytes`, 1×, Ausgang: dieser Slice).
  *Erwartet, zu belegen durch:* die Bewertung im dritten Liefer-Punkt.
  **Ausgang:** *entfallen* — das Wachstum ist gemessen (1,03 bis 1,59 KiB je Change bei
  schmalen Zeilen, Messbericht §3.3) und die Richtgröße bewertet: sie bezieht den Speicher
  nicht ein, 4.000.000 Changes entsprechen 3,9 bis 6,1 GiB (abgeleitet, Messbericht §6);
  der Wert im Code bleibt, die Neubemessung nach dem Trigger von
  [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) „Eine Messung
  liegt vor“ ist Liefer-Punkt des Folge-Slice. Das Register trägt den Eintrag mit dem
  Ausgang *geplant* → `slice-retention-lauf-speicher-begrenzung`.
- **Speicher des Messhosts während der Messreihe.** Die Stufe mit 1.000.000
  Zeilen belastet Docker-Volumes und den Feed-Container; kein `docker volume
  prune` und kein `docker system prune` im Lauf. *Erwartet, zu belegen durch:*
  der Bericht nennt die Aufräum-Schritte des Runners. **Ausgang:** *entfallen* — der
  Messbericht nennt die Schritte (§9: `docker rm -fv` und `docker network rm` je Lauf, kein
  `prune`) und den Bestand (§10: 34 freie Volumes vor und nach den Läufen, kein Container
  und kein Netz `pgc-bench-*` zurück); die Verifikation zählt in ihren Läufen ebenfalls 34
  vor und 34 nach (Verifikation §1, Hygiene).

## 7. Closure-Notiz

- **Was hat funktioniert:** Die Messung war der Sensor, und der Schalter trennte Ursache
  von Wirkung: die Untersuchung des Backfills fand die Ursache im Retention-Lauf — Bereinigung
  aus hält den Speicher flach, `cdc.change` leeren senkt den Live-Heap (Messbericht §4). Die
  Rollen-Kette lief in getrennten Kontexten: der Reviewer (2 HIGH · 3 MEDIUM · 5 LOW · 4 INFO)
  rechnete die Handbuch-Zahlen gegen das Zeilen-Dokument nach und fand F-1 (11,3 statt 18 %
  Untergrenze) und F-2 (1,585 statt 1,51 KiB je Change); der Verifier bestätigte die DoD mit
  `make gates` (Exit 0) und eigenen Läufen der Eingabeseite (Container mit `--memory 64m`
  endet mit Exit 1 und dem Zustand `OOMKilled true`; ohne die Prüfung im Skript nennt der
  Abbruch den Zustand nicht) und fand V-1 durch den Lauf des Suchbefehls über den ganzen
  Baum.
- **Was ging anders als geplant:** (1) Der Plan rechnete mit einem Ausgang „wächst der
  Speicher mit der Tabellengröße“; die Ursache liegt nicht im Backfill, sondern im
  Nachbarsystem Retention und besteht in allen veröffentlichten Server-Versionen (`v0.1.0` bis
  `v0.1.2`, Verifikation §5). Der Ausgang ist ein Änderungs-Slice
  (`slice-retention-lauf-speicher-begrenzung`) mit der Entscheidung
  [`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md); der
  Release-Rahmen ändert sich: Änderungs-Slice zuerst, danach `v0.2.0` (Entscheidung des
  Nutzers, Messbericht §7). (2) Das Werkzeug wuchs über den Plan hinaus: ein zweites Skript
  (`tools/bench-backfill-memory.sh`) und zwei Umgebungsvariablen in `tools/bench-lib.sh`
  (§3, „über den Plan hinaus“). (3) Die Fixrunde lief ohne zweiten Reviewer-Durchgang —
  **benannte Grenze** (V-5), wie bei `slice-backfill-bench-richtgroesse` (dort V-4) und
  `slice-backfill-e2e` (dort V-3): der Verifier maß die Fixrunde nach (Verifikation §3, §4,
  §7), eine Reviewer-Lesung der neuen Texte ist nicht gefahren.
- **Verifier-Beobachtungen (V-1 bis V-5):** *V-1* (MEDIUM): das Feld „Nicht gefunden“ der
  Fixrunde ist berichtigt (§3: Suchwurzeln genannt, Fundstellen außerhalb mit Behandlung).
  [`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md) (`Accepted`)
  trägt in den Zeilen 68, 191 und 199 „1,03 bis 1,57 KiB“ (n = 15), der Messbericht den
  Höchstwert 1,59 (Reihe B, Run 3). **Bekannte, für die Entscheidung folgenlose Abweichung**;
  die ADR und das Architect-Verdikt (Zeilen 152, 155, 167) bleiben unverändert. Eine
  Berichtigungs-ADR setzt der Planner nicht an, weil keine Festlegung der ADR an der zweiten
  Stelle der Zahl hängt: der Faktor der Kandidaten-Größe („sieben bis zehn“ gegen etwa
  0,15 KiB) wird 6,9 bis 10,6, die Bedarfsrechnung 864.000 × 1,03 bis 1,59 KiB ergibt 0,85
  bis 1,31 GiB (abgeleitet, Verifikation V-1). Die Closure des Folge-Slice schreibt bei Bedarf
  die Berichtigung in ihre Suchlauf-Zeile (der Plan trägt die Zeile zur ADR).
  *V-2* (LOW): der Suchausdruck im Plan des Folge-Slice ist auf `1,03 bis 1,5[79]`
  berichtigt, der Stand mit den Trefferzahlen an beiden Ständen genannt (Parent 12, Stand
  `989beef3` 13). *V-3* (LOW): das erste Zählfeld nennt beide Stände (280 am Stand
  `4f94f900`, 288 am Stand `989beef3`), und „15 Reihen“ ist auf 21 Überschriften gemessen
  (`grep -c`) umgestellt. *V-4* (INFO): **benannte Grenze** — die Varianten-Images der
  Reihen E, F, G, I und J (aus dem Arbeitsbaum mit einer geänderten Konstante gebaut) sind
  nicht committet, das Rezept steht als Prosa (Messbericht §2 und §4); der Lauf
  `20260925T165631Z` steht nicht im Zeilen-Dokument (die Nicht-Reproduzierbarkeit von F-10 (b)
  bestätigt der Verifier mit eigenem Lauf, Verifikation §4 Zeile D). Wer die Reihen
  nachfährt, baut nach dem Prosa-Rezept. *V-5* (INFO): **benannte Grenze**, siehe „Was ging
  anders“ (3).
- **Steering-Loop-Eintrag (Lerneintrag):** *Geschärfte Regel (Kandidat, nicht entschieden,
  1×):* Findet ein Untersuchungs- oder Bench-Slice die Ursache eines Symptoms im
  Nachbarsystem, so verschiebt der Befund den **Release-Rahmen**: die Messung mit einem
  Schalter (hier „Bereinigung aus“) trennte Ursache von Wirkung und deckte einen in allen
  veröffentlichten Versionen vorhandenen Defekt auf (`v0.1.0` bis `v0.1.2`); der
  Untersuchungs-Slice trägt deshalb im selben Zug die Reichweite an den Tags, die
  Reihenfolge „Änderungs-Slice vor Server-Release“ als Vorbedingung im Trigger des
  Folge-Slice und die Release-Aussage in dessen Closure —
  `BEO-PGC/vorbestehender-defekt-durch-messung-im-nachbarsystem-gefunden` (neu, 1×, offen,
  Träger nicht gewählt). *Zahlen im Messbericht (bestehende Klasse):* die Handbuch-Zahlen
  wichen im Review zweimal (HIGH) von den gedruckten Zeilen ab (F-1, F-2), dazu zwei
  Zählwörter im Suchlauf-Feld (V-3):
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (22×, verkörpert, Deckel); die Regel
  „Messung zuerst, dann Text“ trägt `AGENTS.md` §3.12 Instanz A — kein neuer Träger.
  *Neuer Sensor:* `tools/bench-backfill-memory.sh` (Vertrag
  `harness/targets/bench-backfill.md`) trennt Run, Nachlauf und Zahl der Changes und liest
  die cgroup-Zähler vom Host aus; ein durch die Speichergrenze beendeter Feed-Container ist
  ein Messergebnis mit Exit 1 in beiden Bench-Skripten (Mutationen: Messbericht §8,
  Verifikation §4). Kein Gate, keine Schwelle (`AGENTS.md` §3.6). *Benannte Spec-Lücke:* das
  Pflichtenheft führt keine Aussage zum Speicher der Retention-Lesung
  ([`LH-FA-RET-004.a`](../../../../spec/pflichtenheft.md) trägt drei Punkte;
  [`SPEC-005`](../../../../spec/pflichtenheft.md) betrifft den Transaktionspuffer des
  Erfassungspfads); der Folge-Slice trägt den Nachzug (Punkt 4 der Verfeinerung).
  *Sensor-Dokumente:* keine Änderung nötig — der Diff berührt weder Code noch Coverage-Profile
  (`git diff --stat 493a28ad..989beef3 -- internal cmd gen proto` leer, gemessen), und
  `git grep -n -i -E 'bench-backfill|Speicher des Feed|memory\.peak|bench-lib|BENCH_FEED'
  -- harness/sensors` trifft nichts (gemessen am Stand `989beef3`).
- **Beobachtungs-Register (`../observations/`):** je Anfall eine Datei
  `evidence/slice-backfill-speicher-untersuchung.md`, Zähler = Zahl der Dateien.
  *Bestehende Klassen:* `BEO-PGC/blockgroesse-zaehlt-zeilen-nicht-bytes` **2×** (Ursache
  untersucht; Ausgang *geplant* → `slice-retention-lauf-speicher-begrenzung`),
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` **22×** (F-1, F-2, V-3),
  `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` **10×** (F-3, F-5, V-2; ab diesem
  Beleg gilt der Deckel bei 10×), `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` **14×**
  (V-1) und `BEO-PGC/arbeit-ueberholt-stehenden-traeger` **32×** (V-1, V-2); die vier
  verkörperten Einträge tragen bereits einen Ausgang. *Neue Klasse:*
  `BEO-PGC/vorbestehender-defekt-durch-messung-im-nachbarsystem-gefunden` **1×**, offen.
  *Ohne Anfall:* `BEO-PGC/dod-kriterium-haengt-am-messhost` (das Kriterium war
  hostunabhängig gefasst, Zähler bleibt 1×), `BEO-PGC/backfill-adapter-startwerte-ohne-messung`
  (1×, keine neue Messung zur Frist oder Einfügeform) und
  `BEO-PGC/geschaetzter-wert-als-grenze` (1×). F-4, F-6, F-8, F-9, F-10, F-11, F-13 und F-14
  sind kein eigener Register-Anfall (Werkzeug- und Text-Findings ohne wiederkehrende Klasse
  im Register). *Lese-Schritt der nächsten Welle-Closure (`welle-transformationen`):* kein
  Eintrag erreicht mit diesem Beleg 3× ohne Ausgang — die vier bestehenden Klassen sind
  verkörpert, die neue steht bei 1×, `blockgroesse-zaehlt-zeilen-nicht-bytes` bei 2× —,
  es entsteht kein Vermerk.
- **Folge-Slices:** `slice-retention-lauf-speicher-begrenzung` (Datei in `open/`;
  Umsetzung von [`ADR-0124`](../../adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md);
  Start-Trigger: dieser Slice in `done/`, Vorbedingung des Server-Release `v0.2.0`). Die
  Zahlen des Handbuchs unter „Grenzwerte“ sind Zahlen vor der Behebung; ihre Ersetzung und
  die Neubemessung der Warn-Richtgröße tragen die Liefer-Punkte des Folge-Slice (Adresse:
  Closure des Folge-Slice `slice-retention-lauf-speicher-begrenzung`).
  Übergaben an offene Pläne (`AGENTS.md` §3.13): der Folge-Slice trägt V-2, den Stand der
  Register-Zähler und die Zeile zur Abweichung 1,57 gegen 1,59; `slice-sdk-kotlin-cloudsmith`
  trägt die Handbuch-`Version:` am Stand `989beef3` (1.61) und die Register-Zähler; die
  Slices von `welle-transformationen` nennen keine Handbuch-Versionsnummer und keine Aussage
  zum Speicher (gemessen mit `git grep` auf Versionsnummern und „Speicher“/„Grenzwerte“ in
  `slice-transformationen-betriebsdoku`, Stand `989beef3`), kein Nachzug.
- **Risiken aus §6:** je ein Ausgang, mit Beleg in §6. *Entfallen:* Die Ursache bleibt
  ungeklärt (die Ursache ist belegt; das Ausbleiben der Takte ab 2.000.000 Changes bleibt
  benannte Grenze mit Adresse im Folge-Slice) · Der Slice läuft nicht vor dem ersten
  Server-Release mit Backfill · Die Warn-Richtgröße bleibt unbegründet · Speicher des
  Messhosts während der Messreihe. *Weiter offen:* Die Messung hängt am Host → Register
  `BEO-PGC/dod-kriterium-haengt-am-messhost` und Risiko der Nachmessung im Folge-Slice.
  *Restrisiken der Verifikation (keine Risiken des Plans):* die Live-Erfassung als
  Speicherquelle ist aus dem Code gelesen, nicht gemessen (Adresse: Handbuch „Grenzwerte“ und
  Messbericht §9, der Folge-Slice übernimmt die Formulierung); die Handbuch-Zahlen sind
  Zahlen vor der Behebung (Adresse: Closure des Folge-Slice, siehe oben).
- **Drei Paarungen:** dieser Slice hat keine Welle. Anker, Folge-Slice und Register trägt
  diese Closure selbst: Folge-Slice `slice-retention-lauf-speicher-begrenzung` in `open/`,
  Register fortgeschrieben (siehe oben). Die Prüfung läuft zusätzlich bei der Closure der
  nächsten Welle ([welle-transformationen](../welle-transformationen.md); die Roadmap führt
  sie unter *Offene Wellen*, das Ereignis kann eintreten).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Bench-Skripte, Handbuch und Use Case des Runs sind
keine eigenen Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(Zähler gemessen am 2026-09-25 mit `ls evidence | wc -l` je Eintrag) —
`BEO-PGC/blockgroesse-zaehlt-zeilen-nicht-bytes` (1×, Ausgang: dieser Slice),
`BEO-PGC/backfill-adapter-startwerte-ohne-messung` (offen, 1×, verwandt: dieselbe
Messung, andere Größe), `BEO-PGC/dod-kriterium-haengt-am-messhost` (offen, 1×,
Risiko §6), `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert, 21×,
die übernommenen Zahlen des Handbuchs).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
