# ADR-0041: a-check als Maschinenform der Architektur-Prüfung

**Status:** Accepted — Supersedes [`ADR-0036`](0036-architekturpruefung-ci.md)

**Datum:** 2026-09-09

**Autor:** pt9912

**Bezug:** [`LH-QA-POR-002`](../../../spec/lastenheft.md)

**Schärft:** [`architecture.md §2`](../../../spec/architecture.md)
(Schichten-Constraints) und [`ARC-001`](../../../spec/architecture.md),
[`ARC-002`](../../../spec/architecture.md),
[`ARC-003`](../../../spec/architecture.md),
[`ARC-004`](../../../spec/architecture.md),
[`ARC-005`](../../../spec/architecture.md),
[`ARC-006`](../../../spec/architecture.md) (Komponenten der Sicht);
die Maschinenform ist [`.a-check.yml`](../../../.a-check.yml).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0036`](0036-architekturpruefung-ci.md) (Proposed) forderte die
Prüfung der Abhängigkeitsrichtung im Gate-Lauf und nannte depguard
(statisches Import-Linting) als Beispiel-Tooling. Operationalisiert ist
das Architektur-Gate **a-check**: pfadbasierte Hexagon-Regeln über
[`.a-check.yml`](../../../.a-check.yml) (zehn Regelmodule,
sprachübergreifend), hermetisch (`--network none` + read-only-Mount,
[`AGENTS.md`](../../../AGENTS.md) §3.1) gegen ein digest-gepinntes
Release-Image (`a-check.mk`, Pin-Hebung = bewusster Commit, Modul 14).
Der erste grüne Lauf über `cmd/pg-change-feed` ist belegt (0 Befunde);
die Wrong-direction-Regel war zuvor unter einer Mutation rot gesehen
worden (Import `internal/domain` → `internal/application/port`, laut
Implementer-Handoff zum Bootstrap-Slice). Die Layer-Globs der
`.a-check.yml` (`internal/**`) werden mit den nächsten Slices
(Domain-Kern, Capture-Persist) belastet.

Die Werkzeug-Frage wurde im Review des Bootstrap-Slice (Findings F-1/F-2)
offengelegt: ADR-0036 und die Fitness-Tabellen der
[`ADR-0002`](0002-abhaengigkeitsrichtung.md),
[`ADR-0003`](0003-physische-modulgrenzen.md),
[`ADR-0031`](0031-paketgrenzen.md) und
[`ADR-0039`](0039-paketstruktur-detaillierung-go.md) nennen depguard,
aktiviert ist a-check. Der ADR-0036-Hochschalt-Trigger („Gate im
Gate-Lauf existiert und läuft grün") ist mit der GATE_CHECKS-Aufnahme
erfüllbar — aber nur mit einer Lesart-Entscheidung, welches Werkzeug die
ADR-0036-Regel trägt. Diese ADR ist sie.

## Entscheidung

Wir wählen **a-check als die Maschinenform der Schichten-Constraints der
Sicht** ([`architecture.md §2`](../../../spec/architecture.md)). Die
[`.a-check.yml`](../../../.a-check.yml) ist ihr deklarativer Stand (Schichten
und Edges, `composition_root`); Änderungen an Schichten-Globs und Edges
sind Änderungen dieser Datei, keine neuen ADRs, solange die §2-Constraints
unverändert abgebildet bleiben. `make a-check` wird in das
`make gates`-Bündel aufgenommen (`GATE_CHECKS += a-check` in
`a-check.mk`); damit ist der [`ADR-0036`](0036-architekturpruefung-ci.md)-
Trigger („Import-Linting-Gate im Gate-Lauf") als Maschinenform erfüllt.
depguard bleibt **nicht** im Werkzeug-Bestand dieses Repos — die in
ADR-0002/0003/0031/0039 benannten **Regeln** (Domain ⊬ Application/
Adapters, Adapter ⊬ Bootstrap, Driving nur über Inbound-Ports) gelten
fort; die Werkzeug-Benennung depguard ist durch diese ADR ersetzt und
bleibt in den Accepted-ADRs als Historie stehen.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

| Option | Pro | Contra |
|---|---|---|
| A — beim ADR-0036-Wortlaut bleiben (depguard einführen) | wörtliche ADR-Treue; Go-idiomatisch, prüft echte Import-Kanten | zweites Toolchain-Bit neben a-check; nur Go, wo die Schicht-Regel sprachübergreifend formuliert ist; lint-Plugins laufen nicht hermetisch Docker-only (AGENTS.md §3.1) |
| B — nichts tun; Schicht-Regel bleibt Review-Prüfpflicht | kein Werkzeug, kein Pflegeaufwand | genau der Zustand, den ADR-0036 beanstandet: bei wachsender Adapterzahl die unzuverlässigste Verteidigung |
| **C — a-check als Maschinenform (gewählt)** | pfadbasiert, sprachübergreifend, hermetisch und digest-gepinnt; deklarativer Stand in einer Datei; deckt die §2-Edges maschinell | pfadgetrieben — feingranulare Import-Verbote (z. B. `time`) sind nicht ausdrückbar; externes Release-Image mit Pin-Pflege; Layer-Globs folgen der Paketstruktur |

## Konsequenzen

- Positiv: Die Schichten-Constraints der Sicht (`architecture.md §2`,
  `ARC-001…006`) sind maschinell geprüft und laufen im `make gates`-Bündel;
  der ADR-0036-Hochschalt-Trigger ist erfüllt und die ADR-0036-Regel
  verliert den Proposed-Zustand. Die Prüfung ist hermetisch (netzlos,
  read-only) und digest-gepinnt — kein Netz- und Pin-Drift im Gate-Lauf.
- Positiv mit Grenze: a-check deckt die Fitness-Zeilen der
  ADR-0002/0003/0031/0039 und [`ADR-0040`](0040-clockport.md) **soweit
  pfadgetrieben möglich** ab (Schichten-Edges, composition_root);
  die `time`-Import-Regel aus ADR-0040 ist **nicht** pfadausdrückbar
  (Pfad-Globs unterscheiden kein importiertes Paket) und bleibt bis zu
  ergänzendem Tooling Review-Prüfpflicht — Grenze hier benannt, nicht
  still.
- Negativ: Das Gate hängt am externen Release-Image
  (`ghcr.io/pt9912/a-check`, Digest-Pin); Pin-Hebung ist ein bewusster
  Commit (Modul 14). Nur die Go-Pfad-Extraktion ist aktiviert; die
  Layer-Globs treffen die `internal/**`-Struktur erst mit den laufenden
  Slices.
- Folgepflicht: `GATE_CHECKS += a-check` in `a-check.mk` (dieser Zug);
  `harness/README.md` trägt `make a-check` in der Sensors-Tabelle mit
  ADR-0041-Bindung statt in der Werkzeuge-Zeile. Die Abdeckungs- und
  Auflösungs-Hinweise des Moduls bei null erfassten Quellen sind **kein
  grüner Gate**: Aktivierungsbedingung bleibt „erste `.go`-Datei im Baum",
  erfüllt seit dem Bootstrap-Slice (`cmd/pg-change-feed/main.go`,
  `composition_root: cmd/**`).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| a-check (Modulzeilen wrong-direction, tech-leak, port-impurity, core/app-impurity, lateral-adapter) | Schichten-Edges der `.a-check.yml`: `ports→domain`, `app→domain`, `app→ports`, `adapters→ports`, `adapters→domain`; `composition_root` außerhalb der Schichten | `make a-check` (im Gate-Bündel) |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Beobachtbarer Trigger: **a-check kann eine geforderte Regel nicht
maschinell ausdrücken** — konkret, wenn das `time`-Import-Verbot aus
[`ADR-0040`](0040-clockport.md) (`internal/domain` und
`internal/application/usecase/*` importieren `time` nicht) maschinell
geprüft werden soll: dann ergänzendes Tooling (z. B. ein
Import-Linting-Gate neben a-check) als Folge-ADR, nicht als Lockerrung
dieser Entscheidung. Andernfalls permanent; Änderungen an Schichten-Globs
und Edgen sind `.a-check.yml`-Änderungen, keine ADRs.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-09 | Accepted — Anlass: Review F-1/F-2 (Review-Report des Bootstrap-Slice) und erster grüner a-check-Lauf im Implementer-Handoff | [`ADR-0036`](0036-architekturpruefung-ci.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-NNNN` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).