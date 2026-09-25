# Review-Report: slice-backfill-bench-richtgroesse — 2026-09-25

**Review-Art:** Code — der Diff führt die Auswertung der beiden Backfill-Warnungen im Use Case des
Runs (`warn.go`, Aufrufe in `service.go`), das vierte Bench-Skript `tools/bench-backfill.sh` samt
Bibliotheks-Nachzug (`bench::median_of`), den Vertrag `harness/targets/bench-backfill.md`, die
Träger (Handbuch 1.53, `harness/README.md`, `docs/user/bench-abdeckung.md`, `Makefile`, Plan mit
Suchlauf-Feld) und, als letzten Commit, das Architect-Verdikt samt
[`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md) ein; geprüft gegen Plan, ADRs,
Pflichtenheft und `AGENTS.md` Hard Rules (Modul 10 §Drei Review-Arten). Kein DoD-Abgleich — das ist
Verifier-Aufgabe (Modul 11).

**Gegenstand:** Slice `slice-backfill-bench-richtgroesse`, Diff-Range `933ab054..a90555b6` (8 Commits,
19 Dateien, +1570/−64; zwei reine `git mv`-Commits `82eba253` und `6f9b391c`; der letzte Commit
`a90555b6` ist das Architect-Verdikt und kein Slice-Inhalt).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere
HIGH-Klassen ergänzt, u. a. Zahl-im-Träger, Beleg-Satz, Zusage-ohne-Eingabeseite,
Kommentar-Chronik, Handbuch-Zug).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-25.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** (`v6.9.0` ·
> `regelwerk/<datei>.md` §<Abschnitt>). Ein `pfad`-Feld auf den **geprüften
> Gegenstand** zitiert den Stand des Laufs und darf ihn festhalten.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-backfill-bench-richtgroesse` (§1 Ziel, §2 DoD als Prüfmaßstab für Plan-Zusagen,
  §3 Plan samt Suchlauf-Feld und „Befunde der Messung“, §6 Risiken) und Welle
  `welle-backfill-bestand`; Architect-Verdikt `architect-verdict-backfill-wal-rueckstand-und-bench-rot`
- [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) (Festlegung 3, Teilfrage 4),
  [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) (Festlegung 3,
  Punkte 1–6),
  [`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md) §(b),
  [`ADR-0104`](../plan/adr/0104-benchmark-schwellen-per-001-002-003.md),
  [`ADR-0114`](../plan/adr/0114-schema-rollout-vorlauf-view-signatur.md),
  [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md) (nur Kontext)
- [`LH-FA-CAP-009`](../../spec/lastenheft.md),
  [`SPEC-029`](../../spec/pflichtenheft.md), [`SPEC-013`](../../spec/pflichtenheft.md),
  [`SPEC-025`](../../spec/pflichtenheft.md)
- `AGENTS.md` (Hard Rules §3.1–§3.13), `harness/conventions.md` (`MR-000`/`MR-001`/`MR-002`)
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-backfill-e2e.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem Implementer-Bericht
übernommen; Exit-Codes ungepiped in Log-Dateien gesichert):

- **Gates am Stand `a90555b6`:** `make test` Exit 0 (`-race`); `make a-check` Exit 0 (gedruckt
  „gesamt: 0 Befund(e)“); `make coverage-gate` Exit 0, gedruckt „coverage-gate: OK — Coverage 83.20%
  erfüllt Schwelle 80%“; `make test-store` Exit 0, gedruckt „DB-Adapter-Coverage: 82.13% (gedeckt 850
  von 1035 Statements; Profile gemergt: store,replication)“ (der Store-Teil trägt die neue
  Assertion des View-Tests); `make commit-traceability RANGE=933ab054..a90555b6` Exit 0, gedruckt
  „commit-traceability: OK — 8 Commit(s) … Betreffs ohne Struktur-ID“.
- **Mutationen der Eingabeseite** (Paket `usecase/backfill`, `go test` im Toolchain-Container, je danach
  `git checkout` der Datei; Ergebnis: Exit 1 = rot, Exit 0 = Mutation überlebt):

  | Nr. | Mutation | Ergebnis |
  |---|---|---|
  | M1 | `warn.go`: Warnung (1) `rows >` zu `rows >=` | rot |
  | M2 | `warn.go`: Warnung (2) `Nanos >` zu `Nanos >=` | rot |
  | M3 | `warn.go`: Sticky-Guard `run.WarnDuration \|\| …` entfernt | rot |
  | M4 | `warn.go`: Guard „Beginn der Kopie nicht gesetzt“ (`StartedAt.Unset()`) entfernt | rot |
  | M5 | `warn.go`: `known &&` entfernt | Kompilierfehler „declared and not used: known“ (kein Test-Rot) |
  | M5b | wie M5, mit `_ = known` (kompiliert) | **überlebt** (Exit 0), siehe F-8 |
  | M6 | `warn.go`: Toleranz-Konstante 10 zu 11 | **überlebt** (Exit 0), siehe F-7 |
  | M7 | `warn.go`: Beginn der Kopie aus `RequestedAt` statt `StartedAt` | rot |
  | M8 | `service.go`: Auswertung in `progress` entfernt | rot |
  | M9 | `service.go`: Auswertung vor `writer.Commit` entfernt | rot |
  | M10 | `service.go`: Auswertung in `conclude` (`failed`/`interrupted`) entfernt | rot |
  | M11 | `service.go`: `warnedEstimatedSize` im Antrag entfernt | rot |
  | M12 | `service.go`: Auswertung im Abschluss ohne Schreibtransaktion (leere Tabelle) entfernt | rot |
  | M13 | `tools/schema/schema.yaml`: View-Spalte `r.warn_estimated_size` zu `false AS warn_estimated_size` | rot durch den neuen Test `TestBackfillStatusViewShowsTheLatestRunPerTable` („vwst_e: … Warnungen false/false — erwartet 7000000, true/false“; gefahren mit einer Scratch-Kopie des Runners, die nur den `postgresstorage`-Aufruf ausführt, weil der `bootstrap`-Aufruf zuvor durch `TestDiagnoseReportsTheLatestBackfillRunPerTable` rot wird) |

- **Ein Lauf von `tools/bench-backfill.sh`** (Default-Stufen, `free -m` vor dem Start 20.460 MB
  verfügbar): Exit 0, Lauf `20260925T011036Z`, gedruckte Zeilen unter „Gemessene Läufe“. Kein
  `make bench` (endet laut Verdikt bekannt rot bei `tools/bench-source-impact.sh`).
- **Suchläufe des Plans an beiden Ständen** (Parent `933ab054` per `git grep … 933ab054`, Diff-Stand
  per `git grep … HEAD`): Zählwörter `drei|vier|zwei|beiden|sechs|sieben|fünf` nahe „bench“ (Parent:
  `docs/user/bench-abdeckung.md:3`, `harness/README.md:150`, `tools/bench-lib.sh:9` und `:171` „drei“;
  Diff-Stand: „drei … mit Pass/Fail-Schwelle“ bzw. „vier“, Zeilen nachgezogen); Ziel-Wörter
  „keine Auswertung|beide bleiben|solange keine“ (Parent: `docs/user/benutzerhandbuch.md:524`,
  `internal/domain/model/backfillrun.go:64` — beide nachgezogen; `tools/schema/schema.yaml:273/408`
  unverändert); `4_000_000`/`4.000.000` (nur `warn.go:23`, eine Handbuch-Angabe und das Plan-Feld);
  `bench-abdeckung.md` trägt drei Zeilen (`LH-QA-PER-001`…`003`), der Kopftext „drei … mit
  Pass/Fail-Schwelle, eines ohne“ zählt richtig; der Diff der Datei berührt nur den erzeugten
  Kopftext, das `Ort`-Feld trägt keine Zeilennummern (kein Lokator durch das Entfernen von
  `median_of` aus `bench-source-impact.sh` verschoben).
- **Hard Rules im Diff:** kein host-lokaler absoluter Pfad in Doku, Skripten oder Code (Suchlauf ohne
  Treffer); kein `slice-`/`welle-`/Vorher-Nachher-Satz und kein Konjunktiv in den hinzugefügten
  Kommentaren von `internal/`, `tools/`, `Makefile`, `harness/` (Suchlauf ohne Treffer); Betreffs
  ohne `SPEC-`/`ARC-`-Kennung, jede Message trägt `LH-*` und `ADR-*`; die zwei Lifecycle-Übergänge
  sind reine `git mv` (0 Zeilen Inhaltsänderung); `tools/schema/schema.yaml` im Diff unberührt
  (`git diff --stat -- tools/schema` leer), die View-Signatur bleibt.
- **Umgebung:** `free -m` verfügbar 20.460 MB vor, 18.949 MB nach dem Bench-Lauf; dangling Volumes
  (`docker volume ls -q -f dangling=true | wc -l`) 34 vor und 34 nach allen Läufen; kein `prune`; nach
  dem Bench-Lauf kein Container und kein Netz `pgc-bench-*`; die Tier-Nebeneffekte
  `tools/schema/plan.yaml`/`down.sql` per `git checkout` zurückgenommen.
- **Nicht gefahren (Grenze):** `make bench` als Ganzes (Verdikt); `--full` (1.000.000-Zeilen-Stufe; laut Handbuch und Verdikt überschreitet der dritte Run
  die WAL-Fehlerschwelle); die
  Läufe `20260924T233628Z` und `20260925T000439Z`, aus denen Handbuch und Konstante stammen (nicht im
  Repository, nur der eigene Lauf ist nachgemessen); der Store-Tier gegen PostgreSQL 17; ein realer
  Ende-zu-Ende-Beleg einer gesetzten Warnung über `make test-integration` (der Plan verlangt ihn
  nicht).

**Gemessene Läufe** (Lauf `20260925T011036Z`, Host Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU,
PostgreSQL 18 `postgres:18-alpine`; gedruckte Zeilen gekürzt):

- „Stufe 10000 Ergebnis — Kopierdauer 1107–1262 (Median 1205, n=3) ms, Durchsatz 7924–9033 (Median
  8299, n=3) Zeilen/s“; „Stufe 50000 … Durchsatz 8763–10509 (Median 9086)“; „Stufe 200000 Ergebnis —
  Kopierdauer 22124–23670 (Median 22388, n=3) ms, Durchsatz 8450–9040 (Median 8933, n=3) Zeilen/s,
  Feed-Speicher-Spitze 416.5 MiB (Ruhe vor der Stufe 176.3 MiB, 20 s nach dem letzten Lauf 641.7
  MiB)“.
- „Schätzung — frisch befüllt, Autovacuum an (100000 Zeilen tatsächlich): reltuples=-1, geschätzt
  unbekannt“; „Autovacuum an, 35 s nach dem Befüllen … Abweichung +0.0%“; „nach ANALYZE und 20000
  weiteren Zeilen (+20 %) ohne erneutes ANALYZE … Abweichung -16.7%“.
- „cdc_capture_lag während des Runs (25 s, Stufe 200000, Live-Last 100/s): 0.123086–1.074292 (Median
  0.590983, n=21) s … Referenz ohne Run: 0.040525–1.010614 (Median 0.404084, n=25) s“.
- „Richtgröße (abgeleitet) — 8933 Zeilen/s (Median, Stufe 200000) × Toleranz 600 s = 5359800 Zeilen,
  abgerundet auf eine Stelle: 5000000 Zeilen“; „WAL-Rückstand-Schwellen (abgeleitet) — bei ~659 B je
  Zeile … Fehlerschwelle (1024 MiB) bei ~1629350 Zeilen“; alle „Warnung Größe f, Warnung Dauer f“
  (keine Stufe erreicht die Toleranz).

---

## Findings

### F-1 — Plan und Handbuch tragen den Stand vor dem Architect-Verdikt; der Nachzug ist benannt, nicht gezogen

- `kategorie`: MEDIUM
- `quelle`: Architect-Verdikt `architect-verdict-backfill-wal-rueckstand-und-bench-rot` (Tabelle
  „Folge-Arbeit“, Zeile `slice-backfill-bench-richtgroesse` (Plan)) · `AGENTS.md` §3.13
  (Träger einer bewegten Eigenschaft nachziehen) · Reviewer-Skill „Zwei-Quellen-Drift“
- `pfad`: `docs/plan/planning/in-progress/slice-backfill-bench-richtgroesse.md:93-102` und
  `:214-219`; `docs/user/benutzerhandbuch.md:452-461`
- `befund`: Der Plan (§3 „Befunde der Messung“, Punkt 1) führt die WAL-Fehlerschwelle als
  Entscheidung offen („ob die Warnung (1) oder die Richtgröße sie berücksichtigt … Folge-ADR“), Punkt 3
  ohne Ursache, der DoD-Punkt „Bench“ nennt weiter einen „realer `make bench`-Lauf“ als Beleg; das
  Verdikt hat sie beantwortet (Capture-Pfad-Defekt, nicht an den Backfill gebunden, Ursache
  Flush-Latenz des Hosts) und dem Planner den Nachzug samt einem Satz Abhilfe im Handbuch-Punkt
  „WAL-Rückstand des Capture-Slots“ (Live-Commit, `wal_retention_error_bytes`) zugewiesen; im Handbuch
  fehlt dieser Satz noch, der Punkt rahmt den Rückstand als Wirkung des Runs. Die Zuweisung ist das
  Übergabe-Artefakt; zum Stand des Laufs ist sie nicht erfüllt.
- `verifizierbar`: ja — `grep -n 'wal_retention_error_bytes' docs/user/benutzerhandbuch.md` ohne
  Treffer im Punkt; Plan-Text gegen die Tabelle des Verdikts lesen.
- `klasse`: Träger vom Verdikt überholt, Nachzug nicht gezogen

### F-2 — Handbuch nennt für den Feed-Speicher eine „Spitze“, die kleiner ist als ein Wert derselben Messung

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.12 Instanz A (Zahl trägt Ursprung, Nachmessen ist die Probe) ·
  Reviewer-Skill „Zahl im Träger ohne Ursprung — oder gegen die Messung driftend“
- `pfad`: `docs/user/benutzerhandbuch.md:1607-1613`; `tools/bench-backfill.sh:266-280` (Zeile
  „Stufe … Ergebnis“)
- `befund`: Das Handbuch führt „Spitze 467 MiB in der Stufe mit 200.000 Zeilen (185 MiB in Ruhe davor)“.
  Nachgemessen (Lauf `20260925T011036Z`, gleiche Stufe): „Feed-Speicher-Spitze 416.5 MiB (Ruhe vor der
  Stufe 176.3 MiB, 20 s nach dem letzten Lauf 641.7 MiB)“ — die Probe 20 s nach dem letzten Run liegt
  54 % über der als „Spitze“ gedruckten. Die Zeile des Skripts, aus der das Handbuch schöpft, nennt
  beide Werte; das Handbuch trägt nur die kleinere und sagt dazu, die Ursache sei „nicht untersucht“.
  Ein Betreiber, der den Speicher des Feed-Containers bemisst, liest 467 MiB als obere Zahl. Für den
  Lauf `20260925T000439Z` ist der Wert 20 s nach dem letzten Run nicht nachgemessen (Lauf nicht im
  Repository).
- `verifizierbar`: ja — `bash tools/bench-backfill.sh` (Default-Stufen) druckt beide Werte je Stufe.
- `klasse`: Spitzenwert nicht der höchste gemessene Wert

### F-3 — Vertrag und `make bench`-Zeile beschreiben die WAL-Messung nicht, die das Skript druckt

- `kategorie`: LOW
- `quelle`: Maintainability · `AGENTS.md` §3.13 (Träger nachziehen) · [`ADR-0120`](../plan/adr/0120-capture-slot-leerlauf-bestaetigung.md)
  Folgepflicht 4 (Bench-Zeile und Vertrag)
- `pfad`: `harness/targets/bench-backfill.md:26-57` (§Ablauf); `harness/README.md:150`
  (`make bench`-Zeile); `tools/bench-backfill.sh:259-281,315-322`
- `befund`: Das Skript druckt je Run den WAL-Rückstand des Slots (`confirmed_flush_lsn`), das vom
  Slot gehaltene WAL (`restart_lsn`), einen Live-Commit zur Freigabe und die Zeile „WAL-Rückstand-
  Schwellen (abgeleitet)“; §Ablauf des Vertrags und die README-Zeile nennen davon nichts (Schritt 2
  listet Dauer, Rate, Schätzung, Warnungen, Speicher). Der Vertrag soll laut Kopf „Ablauf und Grenzen“
  an einem Ort tragen.
- `verifizierbar`: ja — Ausgabe des Skripts gegen §Ablauf lesen.
- `klasse`: Vertrag beschreibt nicht, was das Skript druckt

### F-4 — `BENCH_BACKFILL_RUN_TIMEOUT_S` trägt eine Sekunden-Einheit, die die Schleife nicht zählt

- `kategorie`: LOW
- `quelle`: Maintainability · Reviewer-Skill „Beleg trägt seinen Satz nicht“ (Zusage nennt eine
  Einheit)
- `pfad`: `tools/bench-backfill.sh:25,167-168`; `harness/targets/bench-backfill.md:65`
- `befund`: Der Zähler wächst je Abfrage-Durchlauf, jeder Durchlauf ruft `psql` und `docker stats`
  (Handbuch: Abstand der Proben „etwa 1 bis 2 s“); der Vertrag sagt „Abfragen (im Sekundentakt)“, die
  Variable heißt `_S`. Die reale Wartegrenze liegt über 1800 s und hängt an der Abfragedauer.
- `verifizierbar`: ja — Zeitstempel vor/nach einer Schleife lesen.
- `klasse`: Einheit im Namen ohne Bindung an die Schleife

### F-5 — Zwei Median-Berechnungen in einem Skript, der Plan sagt „einmal in der Bibliothek“

- `kategorie`: LOW
- `quelle`: Maintainability · Plan §3 Zeile `tools/bench-lib.sh` („`bench::median_of` steht einmal in der
  Bibliothek“)
- `pfad`: `tools/bench-backfill.sh:102-104` (`range_of`, awk `v[int((NR + 1) / 2)]`) und
  `:271-274` (`bench::median_of`)
- `befund`: `range_of` berechnet den Median mit eigener Indexformel neben `bench::median_of`; beide
  liefern dasselbe nur, solange sie dieselbe Formel behalten. Ein Lauf druckt beide: den Bereich mit
  Median aus `range_of`, die Rate für die Richtgröße aus `bench::median_of`.
- `verifizierbar`: ja — `grep -n 'NR + 1' tools/bench-backfill.sh tools/bench-lib.sh`.
- `klasse`: doppelte Berechnung derselben Größe

### F-6 — Neue Assertion des View-Tests druckt einen Zeiger statt der Schätzung

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `internal/adapters/driven/postgresstorage/backfillstatusview_test.go:59`
- `befund`: `t.Errorf("vwst_e: Schätzung %v, …", estimate, …)` mit `estimate` vom Typ `*int64`; die
  Mutationsmessung M13 druckte „Schätzung 0xeb0a3abd498“, der Wert ist in der Meldung nicht lesbar
  (die Zeilen 53 und 56 tragen dasselbe Muster).
- `verifizierbar`: ja — M13.
- `klasse`: Fehlermeldung nennt Zeiger statt Wert

### F-7 — Die Werte 10 Minuten und 4.000.000 sind durch keinen Test und keinen Sensor an Handbuch und ADR gebunden

- `kategorie`: INFO
- `quelle`: [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 3
  Punkt 1 (Wert an Konstante und Handbuch) · Reviewer-Skill „Zusage ohne Bindung an ihre Eingabeseite“
- `pfad`: `internal/application/usecase/backfill/warn.go:12,23`
- `befund`: Die Tests lesen beide Konstanten über `export_test.go` und prüfen die Grenzen relativ zu
  ihnen; M6 (Toleranz 10 zu 11) überlebt. Das ist die Folge der „eine Stelle“-Regel und gewollt;
  der Wert im Handbuch (10 Minuten, 4.000.000) und die Konstante driften nur, wenn ein Leser sie
  vergleicht (der Suchlauf des Plans ist die einzige Bindung). Kein Handlungsbedarf, solange die Werte
  sich nicht ändern.
- `verifizierbar`: nein — kein Gate-Lauf; die Mutation M6 zeigt es.
- `klasse`: Konstanten-Wert ohne Bindung an Träger

### F-8 — Mutation `known &&` ist bezüglich der Implementierung äquivalent, der Guard trägt die Zusage von `Rows()`

- `kategorie`: INFO
- `quelle`: Maintainability · Meldung des Implementers („grün-äquivalent“) nachgeprüft
- `pfad`: `internal/application/usecase/backfill/warn.go:29`; `internal/domain/model/backfillrun.go:52-54`
- `befund`: Die nackte Entfernung ist ein Kompilierfehler („declared and not used: known“), kein
  Test-Rot; mit `_ = known` überlebt sie (M5b): `Rows()` liefert für „unbekannt“ `0`, und `0` liegt nie
  über der Richtgröße. Der Guard sichert die dokumentierte Zusage („ist sie unbekannt, ist `rows` ohne
  Aussage“) gegen eine künftige Konstruktion mit `known == false` und `rows != 0`; der Test
  „unbekannt neben großer Zahl“ in `service_test.go` erreicht ihn nicht, weil der Service die Zahl bei
  „unbekannt“ verwirft. Die Meldung des Implementers stimmt.
- `verifizierbar`: ja — M5/M5b.
- `klasse`: äquivalente Mutation bestätigt

### F-9 — Die Richtgröße reproduziert sich nicht: ein zweiter Lauf ergibt 5.000.000, die Konstante trägt 4.000.000

- `kategorie`: INFO
- `quelle`: [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  Festlegung 3 Punkt 2 · Plan §6 (Risiko „Die Richtgröße wird zur Grenze“) · `AGENTS.md` §3.12
- `pfad`: `internal/application/usecase/backfill/warn.go:23`; `docs/user/benutzerhandbuch.md:1573-1585`
- `befund`: Der eigene Lauf `20260925T011036Z` druckt „8933 Zeilen/s … × Toleranz 600 s = 5359800
  Zeilen, abgerundet auf eine Stelle: 5000000“; die Konstante stammt aus einem Median von 7.693
  Zeilen/s (4.615.800, abgerundet 4.000.000). Handbuch und Plan nennen die Streuung (Faktor etwa 2,
  „hängt am Lauf“), und jede Nennung trägt „abgeleitet“, „Orientierung“ und Lauf; die Rundungsregel
  „eine Stelle“ macht aus 16 % Ratenunterschied 25 % Unterschied der Konstante. Kein Fehler des Diffs,
  ein Befund für den Re-Evaluierungs-Fall von [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md).
- `verifizierbar`: ja — `bash tools/bench-backfill.sh`.
- `klasse`: abgeleiteter Wert hängt am Einzellauf

### F-10 — Der Vergleich der Live-Wirkung ist ein Einzellauf; „kein Unterschied gemessen“ gilt für den Lauf, aus dem er stammt

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 · `harness/targets/bench-backfill.md` §Grenzen („keine Wiederholung“)
- `pfad`: `docs/user/benutzerhandbuch.md:1628-1633`
- `befund`: Das Handbuch nennt Median 0,48 s im Run gegen 0,50 s ohne Run; der eigene Lauf druckt
  Median 0,59 s (n=21) gegen 0,40 s (n=25), beide Spannen 0,04–1,07 s. Die Aussage bleibt an den
  gemessenen Lauf gebunden, der Vertrag benennt die fehlende Wiederholung; die Streuung der Mediane
  über Läufe steht im Handbuch nicht.
- `verifizierbar`: ja — `bash tools/bench-backfill.sh`.
- `klasse`: Einzellauf als Vergleichsbasis

### F-11 — Läufe und übernommene Zahlen des Handbuchs sind im Repository nicht auflösbar

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 (Ursprung: übernommen mit auflösbarem Anker) — Träger im Diff, Anker
  außerhalb
- `pfad`: `docs/user/benutzerhandbuch.md:1573-1585` („übernommen aus den Lauf-Berichten, nicht im
  Repository“)
- `befund`: Die Läufe `20260924T233628Z` und `20260925T000439Z` sowie die sechs weiteren Läufe
  (4.088 bis 9.425 Zeilen/s) sind als Ursprung genannt und als „übernommen“ gekennzeichnet, ihre
  gedruckten Zeilen liegen nirgends im Repository. Die Kennzeichnung erfüllt die Form; nachmessen kann
  ein Leser nur mit einem neuen Lauf (F-9).
- `verifizierbar`: nein
- `klasse`: übernommener Wert ohne auflösbaren Anker

### F-12 — `make bench` führt das neue Skript vor `bench-batch-vs-single.sh`; ein Fehlschlag der Messung verdeckt `LH-QA-PER-003` und die Regeneration von `bench-abdeckung.md`

- `kategorie`: INFO
- `quelle`: [`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md) §(b) · Plan §3
  Zeile `Makefile`
- `pfad`: `Makefile:191-194`; `tools/bench-lib.sh:167-171`
- `befund`: Die Reihenfolge steht nicht begründet im Plan. Das Skript endet mit Exit 1 bei ungültiger
  Messung (`rows_copied` ≠ Stufe, Toleranz nicht lesbar), `bench::render_abdeckung` läuft im letzten
  Schwellen-Skript danach; ein Abbruch des Backfill-Skripts verhindert die Schwellen-Prüfung
  `LH-QA-PER-003`. Am Stand des Laufs bricht `make bench` schon an `LH-QA-PER-001` ab (Verdikt), das
  neue Skript läuft nur einzeln.
- `verifizierbar`: nein — `make bench` endet vor dem Skript (Verdikt).
- `klasse`: Reihenfolge der Bench-Skripte ohne Begründung

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/application/usecase/backfill/warn.go` — Warn-Auswertung: Grenzfälle `>` (beide Warnungen: M1/M2 rot), Sticky (M3), `queued` warnt nie (M4), Beginn der Kopie aus `started_at` (M7), beide Konstanten an einer Stelle (`grep` über `internal/`, `tools/`, `docs/`), Uhr über `ClockPort` | geprüft, ohne Befund (F-7/F-8 sind INFO) |
| `service.go` — Aufrufstellen (Antrag, `progress`, Abschluss mit und ohne Schreibtransaktion, `conclude` für `failed`/`interrupted`): M8–M12 je rot; die Warnung ändert weder Status noch Ablauf (Ablaufvergleich gegen den Lauf ohne Warnung im Test), kein Antrag wird wegen der Größe abgelehnt (`Admit` als letzter Schritt in jedem Fall) | geprüft, ohne Befund |
| Tests `warn_test.go`, `export_test.go`, `service_test.go`: Fakes scheitern nur permanent (kein Retry-Pfad im Fake), `scriptClock` injiziert, Eingabeseite mutiert (M1–M12) | geprüft, ohne Befund |
| Adapter-Kette der Warn-Spalten (`backfilladmission.go`, `backfillrun.go`, `backfillwriter.go`, `queries.go`: `warn_duration = (warn_duration OR $n)`, sticky auch in SQL) und `internal/bootstrap` (`diagnose` liest die View, keine Auswertung, keine Konstante, im Diff unberührt) | geprüft, ohne Befund |
| View-Signatur: `tools/schema/schema.yaml` im Diff unberührt; Store-Test `backfillstatusview_test.go` trägt `warn_estimated_size` (M13 rot), Rollen-Bindung im Test unverändert | geprüft, ohne Befund (F-6 LOW) |
| `tools/bench-backfill.sh`: Docker-Hygiene (`docker rm -fv` in `bench::cleanup` als EXIT-Trap, keine Container/Netze/Volumes nach dem Lauf, dangling Volumes 34 → 34), Exit-Codes und Pipes (`set -euo pipefail`, keine Gate-Pipe, `docker stats | awk` unter `pipefail`), Toleranz und Blockgröße aus dem Code gelesen (Lauf druckt „B=1000 Zeilen, Toleranz 10 min“), Gültigkeitsprüfung `rows_copied` (Exit 1 bei Abweichung), Docker-only (nur `docker`, `make`, Standard-Werkzeuge des Hosts wie in den drei Schwester-Skripten) | geprüft, ohne Befund (F-4/F-5 LOW) |
| `tools/bench-lib.sh`/`tools/bench-source-impact.sh`: `bench::median_of` einmal in der Bibliothek, Aufruf im Schwester-Skript, Kopf- und Kommentartexte zählen vier bzw. drei Schwellen-Skripte | geprüft, ohne Befund |
| `Makefile`-Kommentar (kurz, ohne Chronik, Verweis auf den Vertrag) und Hilfetext „vier Skripte“; `docs/user/bench-abdeckung.md`: Entscheidung „keine Zeile für `LH-FA-CAP-009`“ trägt (die Datei bindet Kennungen an durchgesetzte Schwellen, das Skript hat keine; Kopftext zählt drei Zeilen richtig) | geprüft, ohne Befund (F-12 INFO) |
| Handbuch 1.53: `Version:`-Kopf und neue Zeile in `### Änderungshistorie` im selben Diff; keine neue Betreiber-Oberfläche (keine `CDC_*`-Variable, keine neue `cdc.*`-Funktion, kein Endpunkt); Wortwahl „geschätzt“, „Richtgröße“/„Orientierung“, „Startwert, Setzung ohne Messung“, kein „Grenze“/„Limit“/„maximal“ am Begriff (nur „keine Grenze“); Tabelle und Rate gegen die eigene Messung plausibel (7.693 × 600 = 4.615.800, 200.000/26,0 s ≈ 7.692, Schätz-Werte 35 s und −16,7 % im eigenen Lauf identisch) | geprüft, ohne Befund (F-1/F-2/F-9–F-11) |
| Plan-Suchlauf-Feld: fünf der sieben Zeilen an Parent und Diff-Stand nachgefahren (Zählwörter, Makefile, Abdeckungs-Träger, Toleranz/Richtgröße, zweite bewegte Eigenschaft), Handbuch-Grenzwerte und Build-Kontext nur am Diff-Stand gelesen; die Befunde stimmen mit den Fundstellen überein | geprüft, ohne Befund |
| Kommentare nach `AGENTS.md` §3.7 in `warn.go`, `service.go`, `export_test.go`, `backfillrun.go`, den Skripten und dem Makefile: Klassen Zusage/Kopplung/Abgrenzung/Rang-Zeiger, kein Konjunktiv über eine verworfene Alternative, keine Slice-/Wellen-Nummer | geprüft, ohne Befund |
| Hard Rules: kein host-lokaler Pfad (§3.11), kein `//nolint` (§3.2), Coverage-Schwelle ohne Ausnahme (83,20 % ≥ 80 %), Commit-Struktur (§3.3: zwei reine `git mv`), Traceability | geprüft, ohne Befund |
| Spec-Stratum: kein Wert in `spec/pflichtenheft.md` §3, `spec/` im Diff unberührt | geprüft, ohne Befund |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 4 |
| INFO | 6 |

**Finding-Klassen dieses Laufs:** Träger vom Verdikt überholt, Nachzug nicht gezogen ·
Spitzenwert nicht der höchste gemessene Wert · Vertrag beschreibt nicht, was das Skript druckt ·
Einheit im Namen ohne Bindung an die Schleife · doppelte Berechnung derselben Größe ·
Fehlermeldung nennt Zeiger statt Wert · Konstanten-Wert ohne Bindung an Träger · äquivalente
Mutation bestätigt · abgeleiteter Wert hängt am Einzellauf · Einzellauf als Vergleichsbasis ·
übernommener Wert ohne auflösbaren Anker · Reihenfolge der Bench-Skripte ohne Begründung

## Verdikt

**Merge-blockierend:** nein für den Code-Diff — 0 HIGH; Auswertung, Aufrufstellen, Adapter-Kette und
Bench-Skript tragen (13 Mutationen, alle bis auf die zwei benannten äquivalenten/relativen rot; alle
Gates grün). Die zwei MEDIUM sind Träger-Texte, keine Code-Änderungen: F-1 hat das Verdikt bereits dem
Planner zugewiesen (kein Rückgabe-Pfeil an den Implementer), F-2 ist ein Handbuch-Satz (Nennung des
Werts „20 s nach dem letzten Run“ oder Umbenennung der „Spitze“). Beide sind **vor der Closure** des
Slice zu ziehen; die DoD-Zeile „Review durchgeführt“ bleibt in diesem Lauf unberührt (der Auftrag
schließt Plan-Änderungen aus) und wird von der aufrufenden Stelle nach dem Nachzug gesetzt.

**Übergabe:** F-1 an den Planner (Plan §3 „Befunde der Messung“ Punkt 1 und 3, DoD-Punkt „Bench“,
Handbuch-Satz Abhilfe); F-2 an den Planner oder Implementer (Handbuch, ein Satz); F-3 bis F-6 als
LOW an den nächsten Slice, der dieselben Dateien berührt (`slice-backfill-slot-leerlauf-bestaetigung`
ändert `tools/bench-backfill.sh` und den Vertrag laut `ADR-0120` Folgepflicht 4 ohnehin); F-7 bis F-12
ohne Aktion. Die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den
Zähler. Dieser Report ist ein **Lauf-Beleg**; DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11).
