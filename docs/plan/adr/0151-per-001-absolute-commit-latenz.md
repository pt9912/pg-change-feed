# ADR-0151: PER-001 als absolute Commit-Latenz in einer benannten Messumgebung — Supersedes ADR-0104 (nur die PER-001-Schwelle)

**Status:** Proposed — Supersedes [`ADR-0104`](0104-benchmark-schwellen-per-001-002-003.md)
(nur deren Schwelle für `LH-QA-PER-001`, den Quell-Overhead ≤ 35 %; die Schwellen für
`LH-QA-PER-002` und `LH-QA-PER-003` und das Pass/Fail-Prinzip aus `ADR-0104` bleiben
unverändert bestehen). Wird `Accepted`, sobald die Referenzumgebung gemessen ist
(Folgepflicht 1).

**Datum:** 2026-10-04

**Autor:** pt9912 (Architect-Rolle, Modul 8)

**Bezug:** [`ADR-0104`](0104-benchmark-schwellen-per-001-002-003.md) (die teilweise
supersedete ADR),
[`ADR-0054`](0054-coverage-gate-und-benchmark-infrastruktur.md),
[`LH-QA-PER-001`](../../../spec/lastenheft.md),
[ADR-0011](0011-persist-before-ack.md),
[ADR-0010](0010-postgresql-cdc-store.md),
[ADR-0083](0083-herkunft-von-aussagen-in-traegern.md), `AGENTS.md` §3.5, §3.6, §3.12,
Architect-Verdikt zu WAL-Rückstand und rotem Bench (Review-Bericht, ohne Verweis),
`tools/bench-source-impact.sh`.

**Schärft:** [`SPEC-025`](../../../spec/pflichtenheft.md) (`CDC_BENCH_THRESHOLDS`) und
[`SPEC-036`](../../../spec/pflichtenheft.md) (`BENCH_SOURCE_ENV`), beide in
[`spec/pflichtenheft.md`](../../../spec/pflichtenheft.md) §3 (Tabellenzeilen ohne eigenen
Anker).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

`ADR-0104` legte für `LH-QA-PER-001` eine relative Schwelle fest: Overhead ≤ 35 %
des Medians ohne CDC. Auf dem Messhost von `slice-backfill-bench-richtgroesse` endet
`make bench` an diesem Skript rot.

**Gemessen** (Quelle: Architect-Verdikt, Abschnitt der Messungen; Messhost Linux
6.8.0-139-generic, `postgres:18-alpine`):

| Größe | Wert | Herkunft |
|---|---|---|
| 5.000 Einzeltransaktionen ohne CDC | 15.870 ms | gemessen, `tools/bench-source-impact.sh`, Median von 5 |
| dieselbe Last mit CDC | 30.885 ms (94,6 %) | gemessen |
| `fdatasync` | 2.956 µs je Operation | gemessen, `pg_test_fsync` |
| Commit-Latenz ohne CDC | 3,17 ms | abgeleitet (15.870 / 5.000) |
| Commit-Latenz mit CDC | 6,18 ms | abgeleitet (30.885 / 5.000) |
| zusätzliche Latenz je Transaktion | 3,00 ms | abgeleitet ((30.885 − 15.870) / 5.000) |
| Verhältnis zusätzliche Latenz zu `fdatasync` | ≈ 1,02 | abgeleitet, ein Punkt, ein Host |
| Kontrolle `synchronous_commit=off` | 223 ms ohne, 208 ms mit CDC | gemessen, Differenz abgeleitet |
| Host von `ADR-0104`, zusätzliche Latenz | ≈ 0,05 bis 0,085 ms | übernommen, nur grob hergeleitet, Mischung zweier Werte von N, `fdatasync` dort unbekannt |

Lesart (hergeleitet): der Feed schreibt je Quelltransaktion einen eigenen Commit in
den Store; das Verhältnis ≈ 1 zu `fdatasync` an **einem** Messpunkt stützt, dass die
zusätzliche Latenz am Festschreiben des Datenträgers hängt. Mit einem Punkt ist das
keine Regel, und die Kontrolle ohne Festschreiben zeigt keinen Overhead.

Eine relative Schwelle misst darum den Datenträger mit: dieselbe Software liegt je nach
Festschreib-Latenz bei 25 % oder 95 %. **Die Daten reichen nicht für eine einzelne
absolute Zahl, die für jeden Host gilt.**

## Entscheidung

Wir wählen **Option A: eine absolute Obergrenze der zusätzlichen Commit-Latenz je
Quelltransaktion, gültig nur in einer benannten Messumgebung**.

1. **Messgröße im Skript:** Gesamtdauer von N = 5000 einzeln committeten `INSERT`s,
   Median von 5 Läufen, mit und ohne Feed; zusätzliche Latenz je Transaktion =
   (mit − ohne) / N.
2. **Messumgebung `SPEC-036`:** PostgreSQL 18 im Digest-Pin des Skripts, Feed
   co-located, N = 5000, Median von 5, `fdatasync` des Datenträgers des
   Datenverzeichnisses ≤ 0,5 ms (gemessen mit `pg_test_fsync`).
3. **Vorschlagswert:** ≤ 0,10 ms, aufgerundet aus der Altzahl-Obergrenze von
   0,085 ms (*abgeleitet, nicht erprobt*). Der Wert bleibt Vorschlag, bis die
   Referenzumgebung real gemessen ist; bis dahin bleibt diese ADR `Proposed`.
4. **Außerhalb der Umgebung** (der Messhost von `slice-backfill-bench-richtgroesse` ist
   *außerhalb*, `fdatasync` 2,956 ms) druckt das Skript die Messung samt Verhältnis zu
   `fdatasync` und trägt **kein Verdikt**; es endet dort nicht rot.
5. Das Lastenheft nennt nur die Form („absolute Obergrenze, Wert in `SPEC-025`/`SPEC-036`“),
   nicht die Zahl (Muster `LH-QA-PER-002`/`LH-QA-PER-004`).

**Nicht Teil dieses Zuschnitts** (hergeleitet, nicht erprobt): die Latenz selbst zu
senken, etwa Store-Commits über mehrere Quelltransaktionen zu bündeln oder den
CDC-Store auf eine andere Instanz oder einen anderen Datenträger zu legen. Beides
berührt `ADR-0011` (Persist-before-ACK) und `ADR-0010` und ist eine eigene spätere
Entscheidung.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| **A — absolute Zahl nur für eine benannte Messumgebung (gewählt)** | trennt Software-Anteil vom Datenträger; Verdikt dort, wo es etwas aussagt | die Referenzumgebung muss real gemessen werden; außerhalb gibt es kein Verdikt |
| B — relative Schwelle anheben | kleinster Eingriff | misst weiter den Datenträger mit; jede Zahl ist an einem Host rot, am anderen grün |
| C — nur Verhältnis zu `fdatasync` als Schwelle | host-unabhängig | belegt durch **einen** Punkt (≈ 1,02), keine zweite Umgebung gemessen; eine Regel aus einem Punkt wäre ungeprüft |
| D — nichts tun, `make bench` bleibt auf dem Messhost rot | kein Aufwand | ein dauerhaft rotes Skript trägt keine Information (`AGENTS.md` §3.6 verlangt eine ADR für jede Lockerung, nicht ein stilles Rot) |

## Konsequenzen

- Positiv: das Verdikt hängt nicht mehr am Datenträger des Messhosts; die Zahl hat
  eine benannte Umgebung.
- Negativ: ohne Referenzumgebung bleibt die ADR `Proposed` und der Wert ein
  Vorschlag; auf dem Messhost der Entwicklung gibt es kein Verdikt für
  `LH-QA-PER-001` (akzeptiert, begründet in Entscheidung 4).
- Folgepflicht: (1) die Referenzumgebung real messen (zusätzliche Latenz,
  `fdatasync`, Median von 5) und den Wert bestätigen oder korrigieren — danach
  `Accepted`; (2) Slice `bench-source-impact-absolut` (Code-Name ohne Link):
  `tools/bench-source-impact.sh`, `docs/user/bench-abdeckung.md`,
  `harness/README.md` §Sensors (Zeile `make bench`) und
  `harness/sensors`-Verträge auf die neue Messgröße ziehen; (3) die Index-Zeile von
  `ADR-0104` trägt den Zusatz „→ ADR-0151, teilweise“, die Datei von `ADR-0104`
  bleibt unverändert.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `tools/bench-source-impact.sh` in der Umgebung `SPEC-036` | zusätzliche Latenz ≤ 0,10 ms, sonst Exit ≠ 0 — *erwartet*, nicht gefahren | `make bench` |
| `tools/bench-source-impact.sh` außerhalb der Umgebung | druckt Wert und Verhältnis zu `fdatasync`, Exit 0 ohne Verdikt — *erwartet*, nicht gefahren | `make bench` |

## Re-Evaluierungs-Trigger

Die Messung der Referenzumgebung liegt über oder weit unter 0,10 ms (Wert
korrigieren), eine zweite Umgebung widerspricht dem Verhältnis zu `fdatasync`, oder
eine Entscheidung zu Bündelung oder Store-Trennung ändert die Latenz.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-04 | Proposed — Anlass: Auftraggeber-Entscheidung zu einer absoluten Zielgröße je Transaktion; Daten reichen für eine einzelne Zahl nicht | `LH-QA-PER-001` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0151` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
