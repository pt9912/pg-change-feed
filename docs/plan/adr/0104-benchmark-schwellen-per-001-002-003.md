# ADR-0104: Benchmark-Schwellen für PER-001/002/003 — Supersedes ADR-0054 (nur Entscheidung (b) „kein Gate")

**Status:** Accepted — Supersedes [`ADR-0054`](0054-coverage-gate-und-benchmark-infrastruktur.md)
(nur deren Entscheidung (b) „Benchmark-Infrastruktur — Kein Gate"; die
Entscheidung (a) „Coverage-Gate" dieser ADR bleibt unverändert bestehen
und wird hier nicht wiederholt)

**Datum:** 2026-09-19

**Autor:** pt9912

**Bezug:** [`ADR-0054`](0054-coverage-gate-und-benchmark-infrastruktur.md)
(die teilweise supersedete ADR) ·
[`LH-QA-PER-001`](../../../spec/lastenheft.md),
[`LH-QA-PER-002`](../../../spec/lastenheft.md),
[`LH-QA-PER-003`](../../../spec/lastenheft.md) ·
[`SPEC-013`](../../../spec/pflichtenheft.md) (bestehende
`cdc_capture_lag`-Schwelle, wiederverwendet für `LH-QA-PER-002`) ·
[`SPEC-025`](../../../spec/pflichtenheft.md) (neue Schwellen für
`LH-QA-PER-001`/`003`, mit dieser ADR eingeführt) ·
`tools/bench-source-impact.sh`, `tools/bench-scaling.sh`,
`tools/bench-batch-vs-single.sh`

**Schärft:** [`SPEC-025`](../../../spec/pflichtenheft.md) (`CDC_BENCH_THRESHOLDS`)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**(1) Der Anlass.** `make doc-trace` führte `LH-QA-PER-001`…`004` als
„WAISE" — `ADR-0054` (b) begründet das bewusst: „für PER-001…003 gibt es
keine einzelne vergleichbare Schwelle" gegen die ein Pass/Fail sinnvoll
wäre, deshalb dokumentieren die drei Bench-Skripte Aufwand/Ergebnis auf
stdout, ohne jede Schwelle. Im Dialog mit dem Auftraggeber wurde diese
Prämisse geprüft: Für `LH-QA-PER-002` existiert bereits eine reale
Schwelle an anderer Stelle (`SPEC-013`s `cdc_capture_lag`-Fehlergrenze,
ursprünglich für `LH-QA-PER-004` festgelegt, misst aber dieselbe Metrik,
die auch `bench-scaling.sh` bei jeder Lastenstufe abfragt). Für
`LH-QA-PER-001`/`003` existierte keine Schwelle — der Auftraggeber hat
entschieden, real zu messen und danach eine Schwelle mit Marge
festzulegen, statt die Lücke offen zu lassen.

**(2) Die reale Messung, die die Schwellen trägt.** `tools/bench-source-impact.sh`
lief zunächst einfach (kein Median, N=1000) und lieferte an zwei
aufeinanderfolgenden Läufen auf derselben Maschine 17,5 % und 1,8 %
CDC-Overhead — eine Streuung von fast Faktor 10, keine belastbare Basis
für eine Schwelle. Ein erster Median-Nachzug (drei Läufe je Phase,
weiterhin N=1000) linderte das nicht durchgehend: zwei unabhängige
Median-Läufe lieferten 17,3 % und 33,0 %. Untersucht (siehe §Kontext (3)):
kein Fehler, keine Hintergrundlast — die Ursache liegt in der kurzen
Gesamtdauer (rund 200–300 ms für 1000 einzeln committete Transaktionen),
bei der gewöhnliches WAL-Fsync-/Scheduling-Jitter relativ stark
durchschlägt. Mit N=5000 und fünf Läufen je Phase (verlängert die absolute
Messdauer, verdünnt denselben absoluten Jitter relativ) lieferten zwei
unabhängige Läufe 25,0 % und 28,2 % — ein deutlich engeres Band. Die
gewählte Schwelle (35 %) trägt eine reale Marge über diesem Band, nicht
über den früheren, methodisch instabilen Einzelwerten.
`tools/bench-batch-vs-single.sh` lieferte an einem einzelnen Lauf 194×
(Batch- gegen Einzelabruf), eine so große Marge zu jeder sinnvollen
Schwelle, dass ein Median hier nicht nötig ist.

**(3) Die Untersuchung, warum kein zweites Problem vorliegt.** Vor der
Median-/N-Anpassung wurde geprüft, ob die Streuung auf etwas anderes als
Mess-Jitter zurückgeht: `docker stats` gegen alle zum Zeitpunkt laufenden
Fremd-Container zeigte 0–5 % CPU-Auslastung — keine relevante
Hintergrundlast. Die Logs des Feed-Containers während eines vollen Laufs
(`docker logs`, mitgeschnitten während der „mit CDC"-Phase) enthalten
ausschließlich Start-/Verdrahtungs-Meldungen und periodische
`cdc_wal_retention_bytes`-Messungen (alle 5 s) — keine Fehler, keine
Warnungen, keine Reconnects. Der WAL-Rückstand wächst während des Bursts
real (erwartetes Verhalten: kein Consumer hält mit einem plötzlichen
5000-Zeilen-Burst in Echtzeit mit, er drainiert danach ab) — das erklärt
kein verstecktes Fehlverhalten, bestätigt aber, dass die Streuung ein
reales Charakteristikum der gleichzeitigen WAL-Dekodierung unter Burst-
Last ist, kein Mess- oder Software-Fehler.

**(4) Warum `LH-QA-PER-002` keine neue Zahl braucht.** `bench-scaling.sh`
misst bei jeder `SPEC-014`-Lastenstufe real `cdc_capture_lag`
(durchgehend ≈1 s in allen drei Stufen, verkürzter Modus) — genau die
Metrik, die `SPEC-013` bereits mit einer Fehlergrenze (> 60 s) versieht.
Eine zweite, eigene Schwelle für dieselbe Metrik unter demselben Namen
wäre eine zweite Quelle für dieselbe Aussage; `SPEC-025`s Begründung
verweist deshalb ausdrücklich auf `SPEC-013` statt eine eigene Zahl zu
führen.

## Entscheidung

Wir wählen: **Die drei Bench-Skripte für `LH-QA-PER-001`…`003` bekommen
einen echten Pass/Fail-Ausgang gegen die in `SPEC-025` (neu) bzw.
`SPEC-013` (bestehend, wiederverwendet) festgelegten Schwellen — bei
Überschreitung endet das jeweilige Skript mit einem sichtbaren Fehler
und Exit-Code ungleich 0. `make bench` selbst bleibt außerhalb von
`make gates`/`fullbuild` (dieser Teil von `ADR-0054` (b) bleibt
unverändert — der Grund ist unverändert: `make bench` braucht `make
image` und DB-Zugang, ist also strukturell kein netzloser Gate-Kandidat,
unabhängig von der Schwellen-Frage).**

1. **`bench-source-impact.sh` (`LH-QA-PER-001`):** läuft mit N=5000 Zeilen
   und fünf Läufen je Phase (Median-Bildung, siehe §Kontext (2)); scheitert
   mit Exit 1, wenn die Overhead-Prozentzahl (Median „mit CDC" gegen
   Median „ohne CDC") `SPEC-025`s 35-%-Schwelle überschreitet.
2. **`bench-scaling.sh` (`LH-QA-PER-002`):** scheitert mit Exit 1, wenn
   `cdc_capture_lag` bei **irgendeiner** der drei `SPEC-014`-Lastenstufen
   `SPEC-013`s Fehlergrenze (> 60 s) überschreitet — dieselbe Schwelle,
   kein neuer Wert.
3. **`bench-batch-vs-single.sh` (`LH-QA-PER-003`):** scheitert mit Exit 1,
   wenn der gemessene Batch-Vorteil `SPEC-025`s 10×-Mindestwert
   unterschreitet.
4. **RTM-Sichtbarkeit:** Jedes der drei Skripte schreibt zusätzlich eine
   Zeile in eine neue, generierte Datei `docs/user/bench-abdeckung.md`
   (Vorbild: `docs/user/e2e-abdeckung.md`s Erzeuger-Vertrag — stabile
   Deklaration, kein Lauf-Beleg, nur bei inhaltlicher Abweichung
   neu geschrieben). `trace.coverage` in `.d-check.yml` bekommt einen
   zweiten Eintrag, der diese Datei als Coverage-Dimension „Bench" liest
   — `make doc-trace` erkennt `LH-QA-PER-001`/`002`/`003` damit als `ok`,
   ohne dass ein `test-integration`-Lauf sie tragen müsste.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun; Benchmarks bleiben reine Dokumentation ohne Schwelle | kein Eingriff, `ADR-0054` bleibt unangetastet | die reale Prämisse von `ADR-0054` (b) — „keine einzelne vergleichbare Schwelle" — ist für alle drei Anforderungen inzwischen falsch (siehe §Kontext); `LH-QA-PER-001`…`003` bleiben dauerhaft RTM-Waisen |
| B — Schwellen nur informell in den Skript-Kommentaren, kein SPEC-Eintrag | kleinster Eingriff, keine Pflichtenheft-Änderung | widerspricht dem bereits etablierten Muster (`SPEC-013`/`014` für `LH-QA-PER-002`/`004`) — eine Schwelle ohne `SPEC-<NNN>`-Anker ist eine Zahl ohne Herkunfts-Träger (`AGENTS.md` §3.12) und für eine Folge-ADR nicht „schärfbar" |
| **C — `SPEC-025` für die beiden neuen Schwellen, `SPEC-013`-Wiederverwendung für die bestehende, echtes Pass/Fail in den drei Skripten, `make bench` bleibt außerhalb von `make gates` (gewählt)** | folgt dem etablierten `SPEC-<NNN>`-Muster; keine doppelte Schwellen-Führung für `cdc_capture_lag`; die Endstufe ist real gemessen (Median), nicht geraten; `make bench` bleibt strukturell kein Gate-Kandidat (DB-/Image-Abhängigkeit unverändert) | die Skripte tragen jetzt eine echte Fehlklassifikations-Möglichkeit (falsch-positiv bei einmaligem Ausreißer) — durch den Median bei PER-001 und die großzügige Marge bei PER-002/003 gemindert, nicht ausgeschlossen |
| D — Benchmarks direkt in `make gates`/`fullbuild` aufnehmen | maximale Durchsetzung | widerspricht der unveränderten strukturellen Begründung aus `ADR-0054`: `make bench` braucht `make image` und DB-Zugang, kein netzloser Lauf — dieselbe Grenze wie bei `test-integration` |

**Fazit:** C. A lässt eine inzwischen falsche Prämisse stehen; B verzichtet
auf den etablierten Herkunfts-Anker; D widerspricht der unveränderten
strukturellen Docker-only-/Netzlos-Grenze.

## Konsequenzen

- Positiv: `LH-QA-PER-001`/`002`/`003` sind nicht mehr RTM-Waisen, ohne
  `ADR-0054`s strukturelle Docker-only-Begründung für „kein `make
  gates`-Gate" anzutasten.
- Positiv: Keine doppelte Schwellen-Führung — `LH-QA-PER-002` liest
  dieselbe `SPEC-013`-Zahl wie `LH-QA-PER-004`.
- Positiv: Die neuen Schwellen sind an eine reale, Median-gestützte
  Messung gebunden (`AGENTS.md` §3.12), nicht geraten.
- Negativ mit benannter Grenze: Ein einmaliger, ungewöhnlich hoher
  Systemlast-Ausreißer kann `bench-source-impact.sh` trotz Median falsch
  scheitern lassen (fünf statt zehn oder mehr Läufe — ein pragmatischer
  Kompromiss gegen Laufzeit, kein statistisch belastbares
  Konfidenzintervall). Die 35-%-Schwelle (gegenüber real gemessenen
  25,0 % und 28,2 % über zwei unabhängige N=5000-Läufe) trägt eine reale
  Marge, schließt einen erneuten Ausreißer aber nicht aus.
- Folgepflicht (Implementer-Zug): die drei Bench-Skripte (Schwellen-Check
  + `docs/user/bench-abdeckung.md`-Erzeugung), `.d-check.yml`
  (`trace.coverage`-Zweiteintrag), `harness/README.md` (`make bench`-Zeile:
  „kein Pass/Fail-Schwellenwert" ist ab jetzt falsch für PER-001/002/003).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `tools/bench-source-impact.sh` | Overhead-Median (N=5000, 5 Läufe) ≤ 35 % (`SPEC-025`) | `make bench` (kein Gate — Docker-/DB-Abhängigkeit unverändert, `ADR-0054`) |
| `tools/bench-scaling.sh` | `cdc_capture_lag` ≤ 60 s bei jeder `SPEC-014`-Lastenstufe (`SPEC-013`, wiederverwendet) | `make bench` (kein Gate) |
| `tools/bench-batch-vs-single.sh` | Batch-Vorteil ≥ 10× (`SPEC-025`) | `make bench` (kein Gate) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbarer Trigger: Ein realer `make bench`-Lauf schlägt wiederholt
(≥ 3× über unabhängige Läufe) an einer der drei Schwellen fehl, ohne dass
eine reale Regression vorliegt (False-Positive-Rauschen trotz Median/Marge)
— dann Folge-ADR: Schwelle neu kalibrieren oder Lauf-Anzahl je Phase
erhöhen. Sonst permanent.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-19 | Accepted — Anlass: `make doc-trace` führte `LH-QA-PER-001`…`003` als Waise; `ADR-0054` (b) hatte "keine vergleichbare Schwelle" als Grund genannt. Reale Messung im Dialog mit dem Auftraggeber: N=1000/3 Läufe streute real zwischen 17,3 % und 33,0 %; Untersuchung (Hintergrundlast 0–5 % CPU, Feed-Container-Logs fehlerfrei) schloss ein verstecktes Problem aus; N=5000/5 Läufe lieferte darauf zwei unabhängige Werte 25,0 % und 28,2 % — Schwelle 35 % trägt reale Marge darüber. `LH-QA-PER-003` (194×, Einzellauf) trägt die zweite `SPEC-025`-Schwelle; `LH-QA-PER-002` nutzt die bestehende `SPEC-013`-Schwelle | Bench-Läufe dieses Zugs (siehe Slice-Plan) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0104` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
