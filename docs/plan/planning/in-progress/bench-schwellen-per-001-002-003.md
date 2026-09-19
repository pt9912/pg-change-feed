# Slice bench-schwellen-per-001-002-003: reale Pass/Fail-Schwellen für die drei Performance-Benchmarks

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** ohne Welle — kein Closure-Kriterium jenseits der eigenen DoD.

**Bezug:** [`LH-QA-PER-001`](../../../../spec/lastenheft.md),
[`LH-QA-PER-002`](../../../../spec/lastenheft.md),
[`LH-QA-PER-003`](../../../../spec/lastenheft.md);
[`ADR-0104`](../../adr/0104-benchmark-schwellen-per-001-002-003.md)
(Supersedes [`ADR-0054`](../../adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
§(b), teilweise).

**Berührte Spec-Stellen:** [`SPEC-025`](../../../../spec/pflichtenheft.md)
(neu, `CDC_BENCH_THRESHOLDS`), [`SPEC-013`](../../../../spec/pflichtenheft.md)
(bestehend, für `LH-QA-PER-002` wiederverwendet).

**Verantwortlich:** — (wellenlos, direkt umgesetzt, siehe Autor).

**Autor:** Implementer-Agent, direkt beauftragt durch den Auftraggeber im
Anschluss an die RTM-Waisen-Untersuchung (`e2e-drei-rtm-luecken`) und den
Dialog über `ADR-0054`s "kein Gate"-Begründung für die drei
Performance-Benchmarks ("Erst messen, dann Schwelle festlegen"). Kein
separater Priorisierungs-Schritt. **Datum:** 2026-09-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** `LH-QA-PER-001`/`002`/`003` waren RTM-Waisen — real durch
`make bench` geprüft, aber ohne Pass/Fail-Schwelle (`ADR-0054` §(b): "keine
einzelne vergleichbare Schwelle") und ohne für `make doc-trace` sichtbare
Beleg-Quelle. Dieser Slice:

1. Misst real (Median von N=5000/5 Läufen für `LH-QA-PER-001` — ein
   Einzellauf streute zwischen zwei Messungen um fast Faktor 10, 17,5 %
   vs. 1,8 %; ein erster Median-Nachzug bei N=1000 blieb mit 17,3 % vs.
   33,0 % über zwei Läufen weiterhin instabil, siehe §6 und `ADR-0104`
   §Kontext (2)/(3); ein Einzellauf für `LH-QA-PER-003`, dessen Marge zu
   groß für Rauschen ist).
2. Legt Schwellen fest — `SPEC-025` (neu) für `LH-QA-PER-001`/`003`,
   Wiederverwendung von `SPEC-013`s bestehender `cdc_capture_lag`-Grenze
   für `LH-QA-PER-002` (keine zweite Zahl für dieselbe Metrik).
3. Verdrahtet die Schwellen als echtes Pass/Fail in den drei
   `tools/bench-*.sh`-Skripten (`ADR-0104`, Supersedes `ADR-0054` §(b)
   nur in diesem Punkt — `make bench` bleibt außerhalb von `make gates`).
4. Macht die Abdeckung für `make doc-trace` sichtbar über eine neue,
   generierte Datei `docs/user/bench-abdeckung.md` (Vorbild:
   `docs/user/e2e-abdeckung.md`), verdrahtet über einen zweiten
   `trace.coverage`-Eintrag in `.d-check.yml`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- `LH-QA-PER-004` — hat bereits eine Schwelle (`SPEC-013`) und einen
  bestehenden Lasttest-Beleg (`tools/harness/run-integration-tests.sh`),
  **war aber sehr wohl** Teil der neun ursprünglichen RTM-Waisen
  (Korrektur einer fehlerhaften Erst-Klassifikation dieses Slice-Plans —
  ein `make doc-trace`-Nachmessen zeigte sie weiterhin als `WAISE`); nicht
  Gegenstand von `ADR-0054` §(b), deshalb bewusst nicht in diesem Slice,
  sondern über einen Tag-only-Eintrag in der bestehenden Lasttest-Phase
  (`tools/harness/run-integration-tests.sh`) im Folge-Vorgang gelöst.
- Aufnahme von `make bench` in `make gates`/`fullbuild` — bleibt
  strukturell ausgeschlossen: `make bench` braucht `make image` und
  DB-Zugang, kein netzloser Lauf (`ADR-0104` §Verglichene Alternativen,
  Option D verworfen).
- Die übrigen sechs RTM-Waisen aus derselben Neun-Klassifikation
  (`LH-FA-SST-001`/`005`, `LH-FA-CFG-006`, `LH-QA-PER-004`,
  `LH-QA-POR-001`/`002`) — eigener Folge-Vorgang, andere Lösungsform
  (Tag-only bzw. CI-Matrix-Sichtbarkeit statt Benchmark-Schwelle).

## 2. Definition of Done

- [x] `LH-QA-PER-001` erfüllt, Schwelle real durchgesetzt —
      `tools/bench-source-impact.sh` läuft mit N=5000 Zeilen und 5 Läufen
      je Phase (Median); scheitert bei Überschreitung von `SPEC-025`s
      35-%-Schwelle; real gemessen 25,0 % und 28,2 % über zwei
      unabhängige Läufe, Marge eingehalten. Der Weg dorthin (N=1000/3
      Läufe blieb instabil, 17,3 %/33,0 %) ist real untersucht —
      Hintergrundlast anderer Container ausgeschlossen (`docker stats`:
      0–5 % CPU), Feed-Container-Logs während eines vollen Laufs
      geprüft: keine Fehler/Warnungen, nur erwartetes WAL-Rückstand-
      Wachstum unter Burst-Last (siehe §6, `ADR-0104` §Kontext (2)/(3)).
- [x] `LH-QA-PER-002` erfüllt, Schwelle real durchgesetzt —
      `tools/bench-scaling.sh` scheitert, wenn `cdc_capture_lag` bei
      irgendeiner `SPEC-014`-Lastenstufe `SPEC-013`s 60-s-Fehlergrenze
      überschreitet; real gemessen ≈1s bei allen drei Stufen.
- [x] `LH-QA-PER-003` erfüllt, Schwelle real durchgesetzt —
      `tools/bench-batch-vs-single.sh` scheitert, wenn der Batch-Vorteil
      unter `SPEC-025`s 10×-Mindestwert fällt; real gemessen 194×/189,8×/
      191,7× über drei unabhängige Läufe, durchgehend weit über der
      Schwelle.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review — 3 HIGH (F-1/F-2/F-3)
      in derselben Fixrunde behoben (siehe §6, neuer Risiko-Eintrag/Korrektur
      und der neue Beobachtungs-Beleg
      `BEO-PGC/arbeit-ueberholt-stehenden-traeger/evidence/slice-bench-schwellen-per-001-002-003.md`),
      1 MEDIUM als weiter offenes Risiko in §6 übernommen, kein offenes HIGH.
- [x] Doku-Update: `harness/README.md` (`make bench`-Zeile nachgezogen —
      Pass/Fail-Schwellen und reale Messwerte für PER-001/002/003),
      `docs/plan/adr/README.md` ([`ADR-0104`](../../adr/0104-benchmark-schwellen-per-001-002-003.md)-Index-Zeile), `docs/user/bench-abdeckung.md`
      real generiert (`make bench`, drei Zeilen, `trace.coverage` erkennt
      sie).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Beobachtungs-Register fortgeschrieben.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/pflichtenheft.md` | update | `SPEC-025` (`CDC_BENCH_THRESHOLDS`) in §3 Defaults und Konstanten ergänzt — Quell-Overhead ≤25 %, Batch-Vorteil ≥10×, mit Verweis auf die Wiederverwendung von `SPEC-013` für `LH-QA-PER-002`. |
| `docs/plan/adr/0104-benchmark-schwellen-per-001-002-003.md` | neu | Supersedes `ADR-0054` §(b) nur im Punkt "kein Pass/Fail" — die Docker-only-/Netzlos-Begründung gegen `make gates`-Aufnahme bleibt unverändert. |
| `docs/plan/adr/README.md` | update | Index-Zeile für `ADR-0104`. |
| `tools/bench-lib.sh` | update | Zwei neue Helfer: `bench::record_row` (schreibt die Zeile eines Skripts nach `.tmp/bench-abdeckung-rows/<id>.row`, gitignored), `bench::render_abdeckung` (sammelt alle vorhandenen Zeilen zu `docs/user/bench-abdeckung.md`, vom letzten Skript in Ausführungsreihenfolge aufgerufen). |
| `tools/bench-source-impact.sh` | update | N=5000 Zeilen, Median von 5 Läufen je Phase (statt N=1000/1 Lauf) — ein Einzellauf streute real um fast Faktor 10 zwischen zwei Messungen, ein erster Median-Nachzug bei N=1000/3 Läufen blieb instabil (17,3 %/33,0 %); Pass/Fail gegen `SPEC-025` (≤35 %); `bench::record_row`-Aufruf. |
| `tools/bench-scaling.sh` | update | Pass/Fail je Lastenstufe gegen die wiederverwendete `SPEC-013`-Schwelle (≤60s); `bench::record_row`-Aufruf. |
| `tools/bench-batch-vs-single.sh` | update | Pass/Fail gegen `SPEC-025` (≥10×); `bench::record_row`- und `bench::render_abdeckung`-Aufruf (letztes Skript in Ausführungsreihenfolge). |
| `.d-check.yml` | update | `trace.coverage` um einen zweiten Eintrag (`docs/user/bench-abdeckung.md`, Label „Bench") ergänzt. |
| `docs/user/bench-abdeckung.md` | Erzeugnis | von `make bench` neu geschrieben, analog `docs/user/e2e-abdeckung.md`. |
| `harness/README.md` | update | `make bench`-Zeile nachgezogen (Pass/Fail-Aussage für PER-001/002/003 korrigiert). |

## 4. Trigger

**Start** (`next` → `in-progress`): sofort — direkter Auftrag, kein
externer Trigger.

**Rückführungen:**

- `in-progress` → `next` (zu groß): entfällt — Umfang bereits real
  umgesetzt (Skripte laufen, Schwellen verdrahtet).
- `in-progress` → `open` (blockiert): entfällt aus demselben Grund.

## 5. Closure-Trigger

DoD vollständig + Review-Report liegt vor + Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Ein erster Median-Nachzug (N=1000, 3 Läufe je Phase) blieb über zwei
  unabhängige Messreihen instabil (17,3 %/33,0 %) — **Ausgang: eingetreten**,
  real untersucht: Hintergrundlast anderer Container ausgeschlossen
  (`docker stats`: 0–5 % CPU), Feed-Container-Logs während eines vollen
  Laufs geprüft (keine Fehler/Warnungen, nur erwartetes WAL-Rückstand-
  Wachstum unter Burst-Last). Behoben durch N=5000/5 Läufe (verlängert
  die absolute Messdauer, verdünnt denselben absoluten Jitter relativ) —
  zwei unabhängige Läufe danach bei 25,0 %/28,2 %, deutlich engeres Band.
- Ein einmaliger, ungewöhnlich hoher Systemlast-Ausreißer kann
  `bench-source-impact.sh` trotz Median (5 statt z. B. 10 Läufe) falsch
  scheitern lassen — **Ausgang:** weiter offen, benannt in `ADR-0104`
  §Konsequenzen und §Re-Evaluierungs-Trigger (Folge-ADR bei wiederholtem
  False-Positive-Rauschen).
- `bench::record_row`/`bench::render_abdeckung` laufen ohne Sperre — ein
  paralleler `make bench`-Lauf (zwei Aufrufer gleichzeitig) könnte sich
  gegenseitig die `.tmp/bench-abdeckung-rows/`-Dateien überschreiben —
  **Ausgang:** entfallen, kein realer Anwendungsfall für parallele
  `make bench`-Läufe in diesem Repo (Docker-Container-/Netzwerknamen
  sind ohnehin nicht parallelisierbar, `bench::cleanup` würde kollidieren).
- `docs/user/bench-abdeckung.md` ist nach einem `make bench`-Lauf nur
  vollständig, wenn alle drei Skripte in dieser Reihenfolge liefen — ein
  isolierter Aufruf von nur `bench-source-impact.sh` schreibt zwar seine
  eigene `.row`-Datei, aber `render_abdeckung` läuft nicht. Real geprüft
  (Review-Report
  [`review-slice-bench-schwellen-per-001-002-003.md`](../../../reviews/review-slice-bench-schwellen-per-001-002-003.md)
  <!-- d-check:status-provenance -->): auch `make bench` selbst ruft nicht
  zuverlässig alle drei auf — bricht `bench-source-impact.sh` oder
  `bench-scaling.sh` mit `exit 1` (Schwellen-Verstoß) ab, stoppt GNU Make
  das Recipe, bevor `bench-batch-vs-single.sh` (und damit
  `render_abdeckung`) läuft — **Ausgang:** weiter offen, praktisch
  begrenzt, da die Zeilen der Datei statische Schwellen-Deklarationen
  sind, keine Lauf-Kennzahlen — ein veralteter Stand macht `make doc-trace`
  nicht falsch positiv `ok`, solange die betroffene Zeile aus einem
  früheren erfolgreichen Lauf stammt.

## 7. Closure-Notiz

- **Was hat funktioniert:** "Erst messen, dann Schwelle festlegen"
  (Auftraggeber-Vorgabe) hat eine reale Messreihe erzwungen, bevor eine
  Zahl in `SPEC-025` landete — dabei kam eine echte Instabilität
  zutage (N=1000 streute 17,3 %/33,0 %), die eine blind übernommene
  Schwelle hätte falsch scheitern lassen. Statt die Streuung nur
  statistisch zu glätten (mehr Läufe), wurde real untersucht, ob ein
  zweites, anderes Problem vorliegt (`docker stats`, Feed-Container-Logs)
  — beide Hypothesen widerlegt, die Ursache blieb Jitter bei kurzer
  absoluter Messdauer. N=5000/5 Läufe löste es messbar (25,0 %/28,2 %).
- **Was ging anders als geplant:** Der Reviewer fand drei HIGH (F-1/F-2/F-3,
  siehe `docs/reviews/review-slice-bench-schwellen-per-001-002-003.md`):
  eine durch denselben Diff bewegte Eigenschaft
  (`make doc-trace`-Waisenzahl in `harness/README.md`) blieb auf dem
  Vor-Slice-Stand stehen (`AGENTS.md` §3.13-Verstoß, in derselben Datei,
  eine Zeile neben einer bereits aktualisierten Zeile — eine
  Nachbarzeilen-Lücke, keine Enumerations-Lücke über mehrere Dateien wie
  in früheren Fällen), und zwei Kopfkommentare (`tools/bench-scaling.sh`,
  `tools/bench-batch-vs-single.sh`) behaupteten weiterhin "kein
  Pass/Fail", obwohl derselbe Diff echtes Pass/Fail einführte. Alle drei
  in einer Fixrunde (Commit `6a324d70`) behoben; der Verifier bestätigte
  danach unabhängig alle sieben DoD-Prüfpunkte real. Zusätzlich wurde in
  derselben Fixrunde ein eigener Klassifikationsfehler in §1 dieses Plans
  korrigiert (`LH-QA-PER-004` war fälschlich als "nicht Teil der neun
  RTM-Waisen" beschrieben).
- **Steering-Loop-Eintrag:** kein neuer Sensor — der gefundene Fehlermodus
  (bewegte Eigenschaft in Nachbarzeile derselben Datei übersehen, obwohl
  der Implementer dort bereits aktiv war) ist der 14. Beleg der bereits
  verkörperten Regel `AGENTS.md` §3.13, keine neue Beobachtungsklasse.
- **Beobachtungs-Register (`../observations/`):** neuer Beleg
  `evidence/slice-bench-schwellen-per-001-002-003.md` unter der
  bestehenden `BEO-PGC/arbeit-ueberholt-stehenden-traeger/` (bereits
  verkörpert als `AGENTS.md` §3.13, keine neue Schwellen-Prüfung nötig).
- **Folge-Slices:** `rtm-reste-sst-cfg-por` (in Bearbeitung, direkt im
  Anschluss beauftragt: "Jetzt noch die restlichen 6 Anforderungen
  abdecken") — löst die in §1 "Ausdrücklich NICHT in diesem Slice"
  genannten sechs verbleibenden RTM-Waisen (`LH-FA-SST-001`/`005`,
  `LH-FA-CFG-006`, `LH-QA-PER-004`, `LH-QA-POR-001`/`002`).
- **Risiken aus §6:** vier Risiken, vier Ausgänge — (1) N=1000-Instabilität
  → **eingetreten**, real untersucht und durch N=5000/5 Läufe behoben;
  (2) einmaliger Systemlast-Ausreißer trotz Median → **weiter offen**,
  in `ADR-0104` §Re-Evaluierungs-Trigger benannt; (3) `record_row`/
  `render_abdeckung` ohne Sperre bei parallelem `make bench` →
  **entfallen**, kein realer Anwendungsfall; (4) `bench-abdeckung.md`
  nur vollständig bei allen drei Skripten in Reihenfolge → **weiter
  offen**, durch den Reviewer verschärft (auch `make bench` selbst
  bricht bei einem Schwellen-Verstoß vor dem letzten Skript ab, GNU-Make-
  Default-Verhalten real bestätigt), praktisch begrenzt auf statische
  Deklarationszeilen ohne Lauf-Kennzahlen.
- **Drei Paarungen:** Anker — kein `liegt in`-Feld (kein neuer
  Sensor/keine neue Regel verkörpert, nur ein weiterer Beleg für eine
  bestehende). Folge-Slice — `rtm-reste-sst-cfg-por` existiert als Datei
  unter `docs/plan/planning/in-progress/` (Lifecycle-Datei vorhanden).
  Register — `BEO-PGC/arbeit-ueberholt-stehenden-traeger/` existiert mit
  nicht-leerem `evidence/` (jetzt 14 Belege, der neue oben).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area „Performance-Benchmarks"
(`tools/bench-*.sh`, `ADR-0054`) — bereits mehrfach berührt
(`slice-047`, die ursprüngliche `ADR-0054`-Einführung), Schwelle erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
kein Treffer zu diesem konkreten Gegenstand.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Modus-Deklaration
`harness/conventions.md`, Default `PGC`/Greenfield).
