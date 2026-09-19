# ADR-0105: CI-Matrix-RTM-Sichtbarkeit für POR-001/002 — generierte Coverage-Dimension über reale GitHub-Actions-Läufe

**Status:** Accepted

**Datum:** 2026-09-19

**Autor:** pt9912

**Bezug:** [`LH-QA-POR-001`](../../../spec/lastenheft.md),
[`LH-QA-POR-002`](../../../spec/lastenheft.md) ·
[`ADR-0104`](0104-benchmark-schwellen-per-001-002-003.md) (dasselbe
Trägermuster — generierte Coverage-Datei + `trace.coverage`-Eintrag —,
hier auf CI-Matrix-Läufe statt Benchmark-Läufe angewandt) ·
`.github/workflows/e2e.yml`, `.github/workflows/ci.yml`,
`tools/harness/ci-matrix-abdeckung.sh`

**Schärft:** — (Prozess-/Werkzeug-ADR, kein Spec-Stratum berührt: die
Messmethoden von `LH-QA-POR-001`/`002` — „Testumgebung je unterstützter
Version" / „CI-/Testlauf unter Linux" — sind bereits real erfüllt, siehe
`.github/workflows/e2e.yml`/`ci.yml`; diese ADR macht das nur für
`make doc-trace` sichtbar, ändert an der Messmethode nichts.)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

`make doc-trace` führte `LH-QA-POR-001` (mehrere aktiv unterstützte
PostgreSQL-Major-Versionen) und `LH-QA-POR-002` (Linux als primäre
Zielplattform) als „WAISE" — nicht, weil die Anforderungen unbewiesen
wären, sondern weil ihr Beleg strukturell außerhalb dessen liegt, was
`make doc-trace` heute sieht: `.github/workflows/e2e.yml` fährt bereits
eine reale, parallele Zwei-Legs-Matrix (PostgreSQL 17/18, `SPEC-012`,
`fail-fast: false`) gegen den vollen Compose-Integrationstest, und
`.github/workflows/ci.yml` trägt seit `AGENTS.md` §3.10 einen eigenen
Schritt `Linux-Plattform-Assertion (LH-QA-POR-002)`
(`uname -s`/`go env GOOS` gegen `Linux`/`linux`). Beide laufen
ausschließlich auf einem gehosteten GitHub-Actions-Runner — dieselbe
strukturelle Grenze, die `AGENTS.md` §3.10 für Workflow-Verhalten
allgemein benennt: ein Docker-only-Sensor (`AGENTS.md` §3.1) kann echtes
Runner-Verhalten nicht lokal nachbilden, `make gates` sieht diese
Läufe nicht.

Im Dialog mit dem Auftraggeber ("Auch das kann man testen") wurde geprüft,
ob „andere Test-Ebene als `test-integration`" als Grund gelten darf, eine
Anforderung dauerhaft ungedeckt zu lassen — verneint. `ADR-0104` hat
dasselbe Muster bereits für `LH-QA-PER-001`/`002`/`003` gelöst: eine vom
jeweiligen Träger-Lauf generierte Coverage-Datei plus ein
`trace.coverage`-Eintrag in `.d-check.yml`. Diese ADR wendet dasselbe
Muster auf CI-Matrix-Läufe an, mit einem Unterschied: Ihr Träger-Lauf ist
kein lokal ausführbares `make`-Target, das Docker-only bleibt, sondern eine
**abfragende** Kontrolle über die GitHub-REST-API — dieselbe Bewusstseins-
Grenze wie `make image-stale` (Modul 14: „kein Gate, braucht Netz").

## Entscheidung

Wir wählen: ein neues, advisory `tools/harness/ci-matrix-abdeckung.sh`
(aufgerufen über `make doc-ci-matrix`, kein Gate) fragt über die
öffentliche GitHub-REST-API den jeweils letzten **erfolgreichen** Lauf von
`e2e.yml` und `ci.yml` ab, bestätigt real (nicht angenommen):

- `e2e.yml`: beide Matrix-Jobs (`PostgreSQL 17`, `PostgreSQL 18`) tragen
  `conclusion: success` (`LH-QA-POR-001`).
- `ci.yml`: der Schritt `Linux-Plattform-Assertion (LH-QA-POR-002)` trägt
  `conclusion: success` (`LH-QA-POR-002`).

und schreibt bei Erfolg eine neue, generierte Datei
`docs/user/ci-matrix-abdeckung.md` — je Zeile mit Lauf-ID, Commit-SHA und
Lauf-URL als Beleg (`AGENTS.md` §3.12: eine Zahl/Aussage trägt ihren
Ursprung). Ein zweiter `trace.coverage`-Eintrag (Label `CI-Matrix`) in
`.d-check.yml` macht die Datei für `make doc-trace` sichtbar. Scheitert
die Abfrage (kein erfolgreicher Lauf, Netz nicht erreichbar, ein
erwarteter Job/Schritt fehlt oder ist nicht `success`), bricht das Skript
mit `exit 1` ab und schreibt die Datei **nicht** — kein still veralteter
oder erfundener Beleg.

Die API-Abfrage selbst läuft — wie jedes andere Werkzeug in diesem Repo
(`AGENTS.md` §3.1) — in einem Container, nicht auf dem Host: Der bereits
gepinnte `TOOLCHAIN_IMAGE` (`golang:1.27-alpine`, Alpine-Basis) installiert
`curl`+`jq` zur Laufzeit (`apk add`, das Werkzeug braucht ohnehin Netz) und
fragt die öffentliche, unauthentifizierte GitHub-REST-API ab (das Repo ist
öffentlich, `GET .../actions/workflows/.../runs` braucht kein Token) —
kein neuer Host-Toolchain-Bedarf über Docker/`make` hinaus, kein neuer
Digest-Pin (Wiederverwendung von `TOOLCHAIN_IMAGE`).

## Verglichene Alternativen

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun, `LH-QA-POR-001`/`002` bleiben dauerhafte RTM-Waisen | kein neuer Code | genau die Prämisse, die der Auftraggeber verworfen hat ("Auch das kann man testen"); die neun-Klasse aus `harness/README.md`s `make doc-trace`-Zeile bliebe zu einem Viertel unaufgelöst |
| B — `gh`-CLI direkt auf dem Host aufrufen (bereits authentifiziert, kürzester Weg) | am wenigsten Code, `gh` ist hier bereits eingerichtet | verletzt `AGENTS.md` §3.1 wörtlich ("Host braucht nur Docker und GNU make") — ein frischer Klon ohne installiertes/authentifiziertes `gh` könnte das Ziel nicht ausführen; Auth-Token-Handhabung zusätzlich ungeklärt |
| C — RTM-Ausnahmeliste: `LH-QA-POR-001`/`002` fest als "strukturell unmessbar lokal" markieren, kein Beleg | kein neuer Code, keine Netzabhängigkeit | genau die stille Ausnahmeliste, die `AGENTS.md` §3.2 für Coverage bereits ablehnt ("kein Coverage-Pendant zu `//nolint`"); verschiebt das Problem, statt es zu lösen |
| **D — gewählt: Docker-gekapselte GitHub-REST-API-Abfrage + generierte Coverage-Datei (dieses ADR)** | real geprüfter, zitierbarer Beleg (Lauf-ID/SHA/URL); Docker-only respektiert; dasselbe, bereits akzeptierte Muster wie `ADR-0104`; kein neuer Digest-Pin | braucht Netz (wie `make image-stale`); der zitierte Lauf ist der letzte **verfügbare**, nicht zwingend der zu aktuellem HEAD gehörige — bei frisch gepushtem, noch nicht durchgelaufenem Commit zeigt die Datei auf den Vorgänger-Lauf |

## Konsequenzen

- Positiv: `LH-QA-POR-001`/`002` werden für `make doc-trace` real (nicht
  behauptet) `ok`; der Beleg ist ein zitierbarer, öffentlich einsehbarer
  GitHub-Actions-Lauf, keine Behauptung im Fließtext.
- Positiv: kein neuer Host-Toolchain-Bedarf, kein neuer Digest-Pin
  (`TOOLCHAIN_IMAGE` wiederverwendet).
- Negativ: `make doc-ci-matrix` bleibt außerhalb von `make gates` — es
  braucht Netz, dieselbe strukturelle Grenze wie `make image-stale`
  und `make bench` (`ADR-0104`).
- Negativ: Die generierte Datei zitiert den letzten **erfolgreichen**
  Lauf, nicht zwingend einen, der zu HEAD gehört — bei jeder erneuten
  Ausführung überschreibt sie sich mit dem dann aktuellsten Lauf (wie
  `docs/user/bench-abdeckung.md`, „stabile Abdeckungs-Deklaration, kein
  Lauf-Beleg für den aktuellen Stand").
- Folgepflicht: `docs/user/ci-matrix-abdeckung.md` bleibt — wie
  `e2e-abdeckung.md`/`bench-abdeckung.md` — ausschließlich
  skript-generiert; keine manuelle Bearbeitung.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `tools/harness/ci-matrix-abdeckung.sh` | beide `e2e.yml`-Matrix-Jobs und der `ci.yml`-Linux-Assertion-Schritt müssen im zuletzt abgefragten erfolgreichen Lauf `success` tragen, sonst `exit 1` ohne Datei-Schreibversuch | `make doc-ci-matrix` (kein Gate) |
| `.d-check.yml` `trace.coverage` | `docs/user/ci-matrix-abdeckung.md` (Label `CI-Matrix`) macht `LH-QA-POR-001`/`002` für `make doc-trace` als `ok` sichtbar | `make doc-trace` (advisory) |

## Re-Evaluierungs-Trigger

Wird `e2e.yml`s Matrix um weitere PostgreSQL-Versionen erweitert oder
`ci.yml`s Linux-Assertion-Schritt umbenannt, zieht derselbe Diff
`tools/harness/ci-matrix-abdeckung.sh`s Job-/Schritt-Namen nach
(`AGENTS.md` §3.13) — sonst meldet das Skript einen fehlenden Job/Schritt
und schlägt sichtbar fehl, statt still zu veralten.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-19 | Accepted — direkt umgesetzt im Anschluss an die RTM-Waisen-Untersuchung (`e2e-drei-rtm-luecken`, `bench-schwellen-per-001-002-003`) und den Auftrag "Jetzt noch die restlichen 6 Anforderungen abdecken" | `docs/plan/planning/in-progress/rtm-reste-sst-cfg-por.md` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0105` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
