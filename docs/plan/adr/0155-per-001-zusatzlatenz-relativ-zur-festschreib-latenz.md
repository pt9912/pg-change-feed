# ADR-0155: PER-001 — Obergrenze der Zusatzlatenz je Commit relativ zur gemessenen Festschreib-Latenz — Supersedes ADR-0151 und ADR-0104 (nur die PER-001-Schwelle)

**Status:** Accepted — Supersedes [`ADR-0151`](0151-per-001-absolute-commit-latenz.md)
(`Proposed`, nie `Accepted`; sein Status wird auf „Superseded by ADR-0155“ gesetzt, sein
Inhalt bleibt) und [`ADR-0104`](0104-benchmark-schwellen-per-001-002-003.md) (nur deren
Schwelle für `LH-QA-PER-001`, den Quell-Overhead ≤ 35 %; die Schwellen für
`LH-QA-PER-002` und `LH-QA-PER-003` und das Pass/Fail-Prinzip aus `ADR-0104` bleiben
unverändert bestehen).

**Datum:** 2026-10-05

**Autor:** pt9912 (Architect-Rolle, Modul 8)

**Bezug:** [`ADR-0151`](0151-per-001-absolute-commit-latenz.md),
[`ADR-0104`](0104-benchmark-schwellen-per-001-002-003.md),
[`ADR-0054`](0054-coverage-gate-und-benchmark-infrastruktur.md),
[`LH-QA-PER-001`](../../../spec/lastenheft.md),
[ADR-0011](0011-persist-before-ack.md),
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

`ADR-0104` legte für `LH-QA-PER-001` eine relative Schwelle fest (Overhead ≤ 35 % des
Medians ohne CDC); `ADR-0151` (`Proposed`) wollte stattdessen ≤ 0,10 ms zusätzliche
Latenz je Transaktion in einer Messumgebung mit `fdatasync` ≤ 0,5 ms. `ADR-0151` trug
einen benannten Widerspruch: nach ihrem eigenen Modell läge die Zusatzlatenz an der Grenze
dieses Bands bei ≈ 0,5 ms, nicht bei ≤ 0,10 ms. Dieser ADR löst den Widerspruch auf,
indem er das Band streicht und die Grenze **an die im selben Lauf gemessene
Festschreib-Latenz bindet**.

**Gemessen** (Quelle: Architect-Verdikt, Abschnitt der Messungen; Messhost Linux
6.8.0-139-generic, `postgres:18-alpine`):

| Größe | Wert | Ursprung |
|---|---|---|
| 5.000 Einzeltransaktionen ohne CDC | 15.870 ms | gemessen, `tools/bench-source-impact.sh`, Median von 5 |
| dieselbe Last mit CDC | 30.885 ms (94,6 %) | gemessen |
| `fdatasync` | 2,956 ms (2.956 µs) je Operation | gemessen, `pg_test_fsync -s 2` |
| zusätzliche Latenz je Transaktion Δ | 3,00 ms | abgeleitet ((30.885 − 15.870) / 5.000) |
| Verhältnis Δ / `fdatasync` | 1,015 | abgeleitet, ein Messpunkt |
| Läufe 87,5 % bis 95,8 % (sechs) | Verhältnis 0,94 bis 1,03 | abgeleitet aus den Prozentwerten mal 3,174 ms (Commit-Latenz ohne CDC) geteilt durch 2,956 ms; die Prozentwerte sind teils gemessen, teils übernommen, `fdatasync` nur einmal gemessen |
| Kontrolle `synchronous_commit=off` | 223 ms ohne, 208 ms mit CDC | gemessen; Δ ≈ −0,003 ms je Transaktion (abgeleitet), also kein messbarer Software-Anteil ohne Festschreiben |
| Host von `ADR-0104`, Δ | ≈ 0,05 bis 0,085 ms | übernommen, nur grob hergeleitet, `fdatasync` dort unbekannt |

Lesart (*hergeleitet*): der Feed schreibt je Quelltransaktion einen eigenen, fsync-gebundenen
Commit in dieselbe Instanz (`ADR-0011`); die Zusatzlatenz ist deshalb ≈ eine
Festschreib-Latenz, und der Software-Anteil ist neben ihr nicht messbar. Mit einem
gemessenen Punkt ist das keine Regel, trägt aber eine Grenze **mit Spielraum**: die
Grenze unten liegt 46 % über dem höchsten abgeleiteten Verhältnis dieses Hosts.

Eine relative Schwelle gegen die Gesamtdauer (`ADR-0104`) misst den Datenträger mit; eine
absolute Zahl in ms gilt nur an einem Host mit bekannter Festschreib-Latenz. Beide sind
für ein Verdikt auf **jedem** Host untauglich.

## Entscheidung

Wir wählen **Option E: eine Obergrenze der zusätzlichen Commit-Latenz je Quelltransaktion,
formuliert als Funktion der im selben Lauf gemessenen Festschreib-Latenz `t_sync`, mit
einer absoluten Untergrenze**:

> Δ ≤ max(0,10 ms ; 1,5 × `t_sync`)

1. **Δ** ist die zusätzliche Latenz je Quelltransaktion: Gesamtdauer von N = 5000 einzeln
   committeten `INSERT`s, Median von 5 Läufen je Phase, mit minus ohne Feed, geteilt
   durch N (unverändert gegenüber `ADR-0151` Entscheidung 1).
2. **`t_sync`** ist die Latenz einer Festschreib-Operation auf dem Datenträger des
   Datenverzeichnisses der Messinstanz, im **selben Lauf** vor den Phasen gemessen mit
   `pg_test_fsync` im PostgreSQL-Container, und zwar die Zeile der Methode, die die
   Instanz nutzt (`SHOW wal_sync_method`; unter Linux im Regelfall `fdatasync`). Die
   Messung läuft im Lauf, nicht als gepflegte Konstante: Das Verdikt reist mit dem Host.
3. **Faktor 1,5** ist eine Festlegung: die Zusatzlatenz ist ≈ eine Festschreibung
   (Verhältnis 0,94 bis 1,03 an einem Host, abgeleitet), der halbe Festschreib-Anteil
   darüber ist Spielraum für Jitter und einen künftigen kleinen Software-Anteil.
4. **Untergrenze 0,10 ms** ist eine Festlegung, **übernommen** aus dem Vorschlag von
   `ADR-0151`, der aus der Obergrenze 0,085 ms des Hosts von `ADR-0104` gerundet war.
   Sie greift nur bei `t_sync` < 0,067 ms (abgeleitet: 0,10 / 1,5), also bei sehr
   schnellen Datenträgern, wo 1,5 × `t_sync` unter dem Rauschen der Messung läge.
5. **Verdikt auf jedem Host.** Das Skript endet mit Exit 1, wenn Δ die Grenze
   überschreitet. Der frühere „kein Verdikt außerhalb der Umgebung“-Zweig von `ADR-0151`
   entfällt; mit der Grenze am gemessenen `t_sync` gibt es keine Außerhalb-Umgebung.
   Das Skript druckt immer Δ, `t_sync`, das Verhältnis Δ / `t_sync` und die Grenze.
6. **Fail-closed.** Ist `t_sync` nicht messbar (Werkzeug fehlt, Zeile nicht lesbar), endet
   das Skript mit Exit 2 und einer Meldung; es fällt nicht auf die alte relative
   Schwelle und nicht auf einen gedachten Wert zurück.
7. **Lastenheft** nennt nur die Form (Obergrenze der Zusatzlatenz je Commit, bezogen auf
   die gemessene Festschreib-Latenz des Datenträgers der Messung; Wert in
   `SPEC-025`/`SPEC-036`), nicht die Zahlen (Muster `LH-QA-PER-002`/`LH-QA-PER-004`).
8. **`SPEC-036`** verliert das Band `fdatasync` ≤ 0,5 ms und trägt statt dessen die
   Messvorschrift für `t_sync` (Punkt 2).

**Nicht Teil dieses Zuschnitts:** die Latenz selbst zu senken (Store-Commits über mehrere
Quelltransaktionen bündeln, CDC-Store auf eine andere Instanz oder einen anderen
Datenträger legen) — berührt `ADR-0011` und `ADR-0010`, eigene spätere Entscheidung.

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nur absolute Zahl in ms für eine benannte Messumgebung (`ADR-0151`) | einfach zu formulieren | Zahl und Band stehen im Widerspruch (Modell: ≈ 0,5 ms an der Bandgrenze); auf dem Entwicklungshost kein Verdikt |
| B — relative Schwelle gegen die Gesamtdauer anheben (`ADR-0104`) | kleinster Eingriff | misst den Datenträger mit: dieselbe Software 25 % oder 95 % |
| C — nur Verhältnis zu `t_sync`, ohne Untergrenze | host-unabhängig | auf sehr schnellen Datenträgern liegt die Grenze unter dem Messrauschen (falsch-rot) |
| D — nichts tun, `make bench` bleibt auf dem Messhost rot | kein Aufwand | ein dauerhaft rotes Skript trägt keine Information (`AGENTS.md` §3.6) |
| **E — Δ ≤ max(0,10 ms; 1,5 × `t_sync`), `t_sync` im selben Lauf gemessen (gewählt)** | trennt Software-Anteil vom Datenträger; Verdikt auf jedem Host; behält die absolute Zahl des Lastenhefts als Untergrenze; löst den Widerspruch von `ADR-0151` | Faktor 1,5 stützt sich auf einen gemessenen Punkt (mit 46 % Spielraum); das Skript braucht `pg_test_fsync` im Messcontainer |

## Konsequenzen

- Positiv: das Verdikt hängt nicht mehr am Datenträger des Messhosts; der Entwicklungshost
  (Δ = 3,00 ms, Grenze 4,43 ms, abgeleitet) fällt unter die Grenze statt dauerhaft rot;
  der Host von `ADR-0104` (Δ ≈ 0,05 bis 0,085 ms, übernommen) fällt unter die
  Untergrenze.
- Negativ: die Grenze ist an zwei Messpunkten gestützt (ein gemessener, ein übernommener);
  ein Host mit wesentlich höherem Software-Anteil neben dem Festschreiben würde rot
  (gewollt: das ist eine Regression der Software) — oder, wenn die hergeleitete Lesart
  dort nicht trägt, der Anlass einer Folge-ADR (Re-Evaluierungs-Trigger).
- Folgepflicht (ein Slice, `slice-bench-source-impact-absolut`, Plan-Datei in `open/`):
  `tools/bench-source-impact.sh` (Messung `t_sync`, Grenze, Meldetexte, Exit 1/2,
  `bench::record_row`-Text), Hilfetext des Ziels `bench` im `Makefile`, `tools/bench-lib.sh`
  (Kommentar- und Kopfzeilen-Verweise auf `ADR-0104` als Quelle der PER-001-Schwelle),
  `docs/user/bench-abdeckung.md` (vom Generator geschrieben), `harness/README.md` §Sensors
  (Zeile `make bench`); `ADR-0105` (`Accepted`) bleibt unverändert. Die Index-Zeilen:
  `ADR-0104` „→ ADR-0155, teilw.“ (statt `ADR-0151`), `ADR-0151` Status „Superseded“.
  **Bekannte Lücke bis zu diesem Slice:** das Skript prüft weiter die relative
  35-%-Schwelle; das widerspricht der neu gefassten `SPEC-025` und ist ein benannter
  Zwischenstand.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `tools/bench-source-impact.sh` auf jedem Host | Δ ≤ max(0,10 ms; 1,5 × `t_sync`), sonst Exit 1 — *hergeleitet*: die Rechnung auf den vorhandenen Messwerten (Δ = 3,00 ms, `t_sync` = 2,956 ms, Grenze 4,43 ms) ergibt Pass; das Skript ist mit dieser Regel **nicht gefahren**, keine Mutation erprobt | `make bench` |
| dasselbe Skript, `t_sync` nicht messbar | Exit 2 mit Meldung, kein Rückfall — *erwartet*, nicht gefahren | `make bench` |
| dasselbe Skript | druckt Δ, `t_sync`, Verhältnis, Grenze in jedem Lauf — *erwartet*, nicht gefahren | `make bench` |

Eine Mutation (etwa Faktor oder Untergrenze im Skript geändert, ein Lauf, der die Grenze
reißt) ist **nicht erprobt**; „der Implementer fährt sie“ ist eine Erwartung, die der Slice
einlöst.

## Re-Evaluierungs-Trigger

Ein Lauf auf einem weiteren Host zeigt ein Verhältnis Δ / `t_sync` über 1,5 ohne
Regression der Software (Faktor korrigieren) oder deutlich unter 0,5 (Faktor zu weit),
`pg_test_fsync` ist im Pin des Bench-Images nicht verfügbar, oder eine Entscheidung zu
Bündelung oder Store-Trennung ändert die Latenz.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-10-05 | Accepted — löst den benannten Widerspruch von `ADR-0151` (Modell, Band, Altzahl) durch Bindung an `t_sync`; Auftraggeber beauftragte „offene Punkte lösen“ | `LH-QA-PER-001` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0155` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
