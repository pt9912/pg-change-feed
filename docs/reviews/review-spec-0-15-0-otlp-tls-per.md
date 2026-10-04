# Review-Report: Spec 0.15.0 (OTLP-Export, TLS, Token-Liste, PER-001 absolut) — 2026-10-04

**Review-Art:** Design (Spec-Text gegen Bestand, ADRs, Hard Rules; sinngemäße Anwendung der Code-/Plan-Kategorien des Skills)

**Gegenstand:** `git diff 5b2c7ced~1 6a2ff1a4` — vier Commits: 5b2c7ced (Lastenheft), f7420941 (Pflichtenheft), f0ccfacf (Architektur-Sicht), 6a2ff1a4 ([`ADR-0149`](../plan/adr/0149-otlp-metrik-export-mechanismus.md) bis [`ADR-0151`](../plan/adr/0151-per-001-absolute-commit-latenz.md) samt Index)

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-10-04 (Arbeitsbaum, HEAD 6a2ff1a4)
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-04

**Eingangs-Kontext:**

- [`spec/lastenheft.md`](../../spec/lastenheft.md) 0.15.0 mit [`LH-FA-SST-010`](../../spec/lastenheft.md), [`LH-FA-SST-011`](../../spec/lastenheft.md), [`LH-FA-SST-012`](../../spec/lastenheft.md), [`LH-QA-PER-001`](../../spec/lastenheft.md)
- [`spec/pflichtenheft.md`](../../spec/pflichtenheft.md): `SPEC-008`, `SPEC-009`, `SPEC-016`, `SPEC-025`, `SPEC-033` bis `SPEC-036`
- [`spec/architecture.md`](../../spec/architecture.md): `ARC-011`, neue Sequenz
- [`ADR-0149`](../plan/adr/0149-otlp-metrik-export-mechanismus.md), [`ADR-0150`](../plan/adr/0150-tls-und-mehrfach-token.md), [`ADR-0151`](../plan/adr/0151-per-001-absolute-commit-latenz.md); [`ADR-0104`](../plan/adr/0104-benchmark-schwellen-per-001-002-003.md), [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md), [`ADR-0101`](../plan/adr/0101-zugangsdaten-klasse-sieben-schluessel.md)
- Architect-Verdikt `architect-verdict-backfill-wal-rueckstand-und-bench-rot` (Zahlen von [`ADR-0151`](../plan/adr/0151-per-001-absolute-commit-latenz.md))
- `AGENTS.md` §3.4, §3.5, §3.7, §3.12, §3.13

**Nachgefahren (Exit direkt):** `make gates` 0 · `make commit-traceability` 0 · `make doc-immutable RANGE=5b2c7ced~1..6a2ff1a4` 0 · `make doc-commits RANGE=5b2c7ced~1..6a2ff1a4` 0 · `make doc-trace` Exit 0, gedruckte Zeile „83 Anforderung(en), 3 Waise(n).“ (die drei sind SST-010/011/012, wie erwartet).

---

## Findings

### F-1 — Vorschlagswert 0,10 ms widerspricht der Lesart derselben ADR an der Grenze der Messumgebung

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.12 (Aussage breiter als Messung) · [`ADR-0151`](../plan/adr/0151-per-001-absolute-commit-latenz.md)
- `pfad`: docs/plan/adr/0151-per-001-absolute-commit-latenz.md:48 und :74 ; spec/pflichtenheft.md `SPEC-025`/`SPEC-036` (§3)
- `befund`: [`ADR-0151`](../plan/adr/0151-per-001-absolute-commit-latenz.md) liest die Zusatzlatenz als ≈ 1,02 × `fdatasync` (ein Punkt); die Umgebung `SPEC-036` lässt `fdatasync` bis 0,5 ms zu, was nach dieser Lesart ≈ 0,5 ms Zusatzlatenz ergibt, der Vorschlag lautet aber ≤ 0,10 ms. Die 0,10 ms stammen aus einem Host, dessen `fdatasync` dort „unbekannt“ ist (0,085 ms, übernommen), die Grenze 0,5 ms aus keiner Messung; Vorschlagswert und Umgebungsgrenze sind nicht gegeneinander hergeleitet.
- `verifizierbar`: nein — erst die Messung der Referenzumgebung ([`ADR-0151`](../plan/adr/0151-per-001-absolute-commit-latenz.md) Folgepflicht 1) löst es auf; ohne Sensor.
- `klasse`: Zahl-Paar ohne gemeinsame Herleitung

### F-2 — Träger der neuen PER-001-Messgröße: Liste unvollständig, Adresse ist kein committeter Plan

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.13 · Skill „Aufschub-Adresse deckt Gegenstand?“
- `pfad`: docs/plan/adr/0151-per-001-absolute-commit-latenz.md (Folgepflicht 2); Makefile:246 ; tools/bench-source-impact.sh:22-23,104-108 ; docs/user/bench-abdeckung.md:14 ; harness/README.md (Zeile `make bench`)
- `befund`: Folgepflicht 2 nennt Skript, `bench-abdeckung.md`, `harness/README.md` und die Sensor-Verträge, nicht `Makefile:246` („[`LH-QA-PER-001`](../../spec/lastenheft.md)…003 mit Schwelle“); die Adresse `bench-source-impact-absolut` ist ein Name ohne Plan-Datei (`git grep` im Baum: nur in [`ADR-0151`](../plan/adr/0151-per-001-absolute-commit-latenz.md)). Bis zu diesem Slice widerspricht das Skript (35-%-Schwelle, [`ADR-0104`](../plan/adr/0104-benchmark-schwellen-per-001-002-003.md) Accepted) der neu gefassten `SPEC-025` im selben Baum, und die Index-Zeile von [`ADR-0104`](../plan/adr/0104-benchmark-schwellen-per-001-002-003.md) verweist auf eine `Proposed`-ADR.
- `verifizierbar`: ja — `git grep -n -i '35 %\|35%\|SPEC-025\|bench-source-impact' -- . ':!docs/plan/done' ':!docs/reviews' ':!.harness'` (Treffer: Makefile, Skript, bench-abdeckung, harness/README, bench-lib.sh:160/190 nur [`ADR-0104`](../plan/adr/0104-benchmark-schwellen-per-001-002-003.md)-Verweis; [`ADR-0105`](../plan/adr/0105-ci-matrix-rtm-sichtbarkeit-por-001-002.md) Verweis auf [`ADR-0104`](../plan/adr/0104-benchmark-schwellen-per-001-002-003.md) unberührt).
- `klasse`: Träger-Nachzug mit verwaister Adresse

### F-3 — ARC-011 führt „strukturierte Logs“ am OTLP-Empfänger, das Lastenheft schließt Log-Übertragung aus

- `kategorie`: MEDIUM
- `quelle`: Skill „Nachzug widerspricht dem Nachbarn im selben Träger“ · [`LH-FA-SST-010`](../../spec/lastenheft.md) Out-of-Scope
- `pfad`: spec/architecture.md:98 (`ARC-011`) gegen spec/lastenheft.md (SST-010 Out-of-Scope: „keine Übertragung von Traces oder Logs“) und architecture.md Sequenz (nur Metriken)
- `befund`: Die geänderte Zeile `ARC-011` nennt als Zweck „Metriken (Push) und strukturierte Logs“ für den OTLP-Empfänger; die Anforderung und die neue Sequenz übertragen nur Metriken, Logs gehen nicht über OTLP. Zwei Aussagen im selben Träger, keine verweist auf die andere.
- `verifizierbar`: nein (Lese-Handlung).
- `klasse`: Nachzug widerspricht Nachbar

### F-4 — Zugangsdaten-Klasse wächst von sieben auf elf Schlüssel; ADR-0101 (Accepted) nennt sieben, ADR-0149/0150 deklarieren „kein Supersedes“

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.5, §3.13 · [`ADR-0101`](../plan/adr/0101-zugangsdaten-klasse-sieben-schluessel.md)
- `pfad`: spec/pflichtenheft.md:608-617 (`SPEC-016`) ; docs/plan/adr/0150-tls-und-mehrfach-token.md (Kopf, „kein Supersedes“) ; internal/bootstrap/config_file_internal_test.go (nach [`ADR-0101`](../plan/adr/0101-zugangsdaten-klasse-sieben-schluessel.md) Zeile 111: Iteration über sieben Schlüssel)
- `befund`: `SPEC-016` führt die Klasse jetzt mit `otlp_endpoint`, `otlp_headers` und zwei Token-Listen; [`ADR-0101`](../plan/adr/0101-zugangsdaten-klasse-sieben-schluessel.md) legt die Klassengröße ausdrücklich auf sieben fest und trug frühere Erweiterungen als eigene (supersedierende) ADR nach. [`ADR-0149`](../plan/adr/0149-otlp-metrik-export-mechanismus.md) und [`ADR-0150`](../plan/adr/0150-tls-und-mehrfach-token.md) nennen [`ADR-0101`](../plan/adr/0101-zugangsdaten-klasse-sieben-schluessel.md) weder im Bezug noch in der Folgepflicht (Test, der die Klasse iteriert), „Schärft `SPEC-016`“ ersetzt das nicht.
- `verifizierbar`: ja — `git grep -n 'sieben' -- docs/plan/adr/0101* internal/bootstrap`.
- `klasse`: Accepted-ADR-Zahl durch Erweiterung überholt, ohne Zeiger

### F-5 — Happy-Path-Kriterium von SST-010 behauptet Gleichheit „im selben Moment“ für Werte, die das nicht tragen

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12 Instanz B · [`LH-FA-SST-010`](../../spec/lastenheft.md)
- `pfad`: spec/lastenheft.md (SST-010 Happy Path) ; tools/schema/nacharbeit-observability.sql (`cdc_capture_lag`, `cdc_oldest_change_age_seconds` über `now()`) ; [`ADR-0149`](../plan/adr/0149-otlp-metrik-export-mechanismus.md) Fitness-Zeile 4
- `befund`: „Ein übertragener Wert entspricht dem Wert, den die SQL-Sicht im selben Moment liefert“ gilt für die zwei `now()`-abhängigen Kennzahlen nur mit Toleranz, und `cdc_wal_retention_bytes` hat keine Zeile in der Sicht. Die Fitness-Zeile von [`ADR-0149`](../plan/adr/0149-otlp-metrik-export-mechanismus.md) übernimmt „Wert im Export = Wert der Sicht“ ohne diese Eingrenzung (als *erwartet* gekennzeichnet, daher LOW).
- `verifizierbar`: nein — kein Lauf existiert.
- `klasse`: Aussage breiter als ihre Menge

### F-6 — Warn-Bereich 6: Träger außerhalb des Diffs, Herkunft in ADR-0144 geschärft statt überschrieben

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.13 · [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md)
- `pfad`: internal/domain/messagecode/codes.go:131 ; internal/domain/messagecode/messagecode_test.go:60-63 ; docs/user/benutzerhandbuch.md:2440
- `befund`: [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) bleibt unverändert (Festlegung zählt Bereiche 1 bis 5 und 9 reserviert); `SPEC-008` (höherer Rang) ergänzt 6, [`ADR-0149`](../plan/adr/0149-otlp-metrik-export-mechanismus.md) nennt das als Schärfung und überschreibt die ADR nicht — in Ordnung. Drei lebende Träger beschreiben die Bereichsmenge (Code-Kommentar, Test-Kommentar „Bereich 1 bis 5“, Handbuch „Bereich 5 … vorgesehen“) und sind im Slice nachzuziehen; [`ADR-0149`](../plan/adr/0149-otlp-metrik-export-mechanismus.md) Folgepflicht nennt sie nicht einzeln.
- `verifizierbar`: ja — `git grep -n 'Bereich 5\|9 reserviert' -- internal docs/user`.
- `klasse`: Träger-Nachzug (Bereichsmenge)

### F-7 — Kleinigkeiten an Zahlen und Einheiten

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 · [`ADR-0149`](../plan/adr/0149-otlp-metrik-export-mechanismus.md) und [`ADR-0151`](../plan/adr/0151-per-001-absolute-commit-latenz.md)
- `pfad`: [`ADR-0151`](../plan/adr/0151-per-001-absolute-commit-latenz.md) Kontext-Tabelle („2.956 µs“) und Entscheidung 4 („2,956 ms“); [`ADR-0149`](../plan/adr/0149-otlp-metrik-export-mechanismus.md) Entscheidung 4 (`cdc_consumer_lag` Einheit `1`); [`SPEC-033`](../../spec/pflichtenheft.md) Tabelle (`cdc_consumer_position`)
- `befund`: Dieselbe Messung steht als „2.956 µs“ (deutsche Tausendertrennung) und „2,956 ms“ — korrekt, aber leicht als 2,956 µs zu lesen. `cdc_consumer_lag` ist laut SQL-Text ein LSN-Byte-Abstand (Einheit `1` ist als bewusst vermerkt), und `cdc_consumer_position`/`-lag` als OTLP-Gauge tragen eine 64-Bit-Position in einer Gleitkomma- oder int64-Zahl (Präzision/Vorzeichen ungeprüft).
- `verifizierbar`: nein.
- `klasse`: Einheit/Notation

### F-8 — Index-Titel von ADR-0104 geändert

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.5
- `pfad`: docs/plan/adr/README.md:119
- `befund`: Der Index ist kein ADR-Inhalt, die Titeländerung ist zulässig und `git diff … -- docs/plan/adr/0104-*` ist leer. Der Titel verweist jedoch auf eine `Proposed`-ADR und wechselt unbegründet „Benchmark“ zu „Bench“.
- `verifizierbar`: ja — `git diff 5b2c7ced~1 6a2ff1a4 -- docs/plan/adr/0104-benchmark-schwellen-per-001-002-003.md` leer.
- `klasse`: Index-Zeile eilt Status voraus

## Negativbefunde

- geprüft, ohne Befund: spec/lastenheft.md — SST-010/011/012 vollständig nach Bestandsform (Beschreibung, Happy/Boundary/Negative, Out-of-Scope); Kennungen fortlaufend und eindeutig; keine MVP-Kennzeichnung, wie bei SST-009 (stimmt); Version 0.15.0 und Historienzeile in der Form von 0.12.0–0.14.0; OPS-003 und SST-004 laut `git diff` unverändert; §5: „High Availability“ bleibt, Bestand nicht verfälscht; PER-001 nennt keine Zahl, Messmethode (Differenz der Gesamtdauern / N) trägt
- geprüft, ohne Befund: spec/pflichtenheft.md — [`SPEC-033`](../../spec/pflichtenheft.md) bis 036 vollständig; Herkunft je Zahl in [`SPEC-025`](../../spec/pflichtenheft.md) und 036 gekennzeichnet (Vorschlag, abgeleitet); `SPEC-016`-Präzedenz bleibt (env schlägt Datei, env-exklusive Liste ergänzt); `SPEC-009`-Name `cdc_oldest_change_age_seconds` stimmt mit Sicht und Handbuch (Zeile 1199); [`SPEC-018`](../../spec/pflichtenheft.md)-Verweise und §6/§7 konsistent
- geprüft, ohne Befund: Kennzahlen-Menge — der Befehl `grep -o` auf `SELECT 'cdc_…'` in tools/schema/nacharbeit-observability.sql mit `sort -u | wc -l` ergibt 9 (nachgefahren), plus `cdc_wal_retention_bytes` = 10 Zeilen der Tabelle in `SPEC-033`; Gauge-Behauptung am SQL-Text hergeleitet und als *hergeleitet* gekennzeichnet (`count(*)` auf `cdc.transaction`/`cdc.change`, `cdc_errors_total` zählt Quellen je Klasse); `runWALRetentionCheck` in internal/bootstrap/wiring.go:1345 misst den WAL-Wert im Prozess (Kommentar 1325 ff. nachgelesen)
- geprüft, ohne Befund: Token-Bestand — `==` in internal/adapters/driving/http/middleware.go:40/43 und zweite Fassung in internal/adapters/driving/grpc/interceptor.go:51-60; admin zuerst, leeres Token matcht nie; [`ADR-0150`](../plan/adr/0150-tls-und-mehrfach-token.md) „Gemessen“-Zeilenbereiche (30-50, 45-65) enthalten beide Funktionen
- geprüft, ohne Befund: spec/architecture.md — `git diff … | grep -n -i` auf ADR-/Slice-/Welle-Bezüge leer (§3.4); Sequenz konsistent mit `SPEC-033` (eigene Goroutine, Frist, kein Nachholen, nur Log); keine neue ARC-Komponente
- geprüft, ohne Befund: [`ADR-0149`](../plan/adr/0149-otlp-metrik-export-mechanismus.md), [`ADR-0150`](../plan/adr/0150-tls-und-mehrfach-token.md) und [`ADR-0151`](../plan/adr/0151-per-001-absolute-commit-latenz.md) Form und Status (erste zwei Accepted, dritte Proposed mit Supersedes [`ADR-0104`](../plan/adr/0104-benchmark-schwellen-per-001-002-003.md) nur der PER-001-Schwelle); alle Fitness-Function-Zeilen tragen „erwartet, nicht gefahren“, keine behauptete Erprobung; Zahlen von [`ADR-0151`](../plan/adr/0151-per-001-absolute-commit-latenz.md) gegen das Architect-Verdikt (15.870 ms, 30.885 ms, 94,6 %, 2.956 µs, 223/208 ms) gemessen-gleich, abgeleitete Werte (3,17 / 6,18 / 3,00 ms, ≈ 1,02) nachgerechnet; [`ADR-0104`](../plan/adr/0104-benchmark-schwellen-per-001-002-003.md), [`ADR-0144`](../plan/adr/0144-meldungscodes-nutzerseitige-kennungen.md) und [`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md) unberührt
- geprüft, ohne Befund: Commit-Reihenfolge (Lastenheft, Pflichtenheft, Architektur, ADRs) und Betreffzeilen (keine Struktur-Kennung im Betreff, `make commit-traceability` Exit 0)
- geprüft, ohne Befund: Träger-Suchlauf (Prometheus) — im lebenden Baum nur [`ADR-0024`](../plan/adr/0024-observability-ausserhalb-der-domain.md) (Accepted, Zitat-Korrektur nach [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) möglich, nicht nötig, da „Prometheus/OpenTelemetry“ dort Beispiel ist); [`ADR-0040`](../plan/adr/0040-clockport.md) trägt keinen Treffer auf Prometheus/OpenTelemetry (nachgemessen, widerspricht der Meldung des Implementers)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 4 |
| LOW | 1 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Zahl-Paar ohne gemeinsame Herleitung · Träger-Nachzug mit verwaister Adresse · Nachzug widerspricht Nachbar · Accepted-ADR-Zahl durch Erweiterung überholt · Aussage breiter als ihre Menge

## Verdikt

**Merge-blockierend:** nein für den Push der Spec-Commits (kein HIGH, `make gates` und alle vier Einzelgates Exit 0); die vier MEDIUM sind vor dem ersten umsetzenden Slice bzw. vor `Accepted` von [`ADR-0151`](../plan/adr/0151-per-001-absolute-commit-latenz.md) zu schließen: F-3 als Korrektur an `spec/architecture.md`, F-4 als Zeiger oder Folge-ADR zu [`ADR-0101`](../plan/adr/0101-zugangsdaten-klasse-sieben-schluessel.md), F-1 vor `Accepted` von [`ADR-0151`](../plan/adr/0151-per-001-absolute-commit-latenz.md), F-2 durch einen committeten Plan für die Adresse. Dass die Spec [`SPEC-025`](../../spec/pflichtenheft.md) mit einem Vorschlagswert führt, während [`ADR-0151`](../plan/adr/0151-per-001-absolute-commit-latenz.md) `Proposed` ist und `make bench` noch 35 % prüft, ist benannt, aber eine bewusst offene Zwischenlage.

**Übergabe:** Findings gehen an den Planner/Architect (Spec-Diff, kein Implementer-Code); ein Rollen-Widerspruch liegt bisher nicht vor, daher kein Konflikt-Pfad. Dieser Report ersetzt keine Verifikation.
