# Architect-Verdikt: Coverage-Gate — Messgegenstand und Endstufe 80 %

**Rolle:** Architect (Modul 8)

**Anlass:** Nutzerfrage an den Planner — „Soll dieses Gate messen, was
netzlos prüfbar ist, oder den ganzen Baum?" — und die Planner-Antwort, 80 %
sei „unter der eigenen Definition" dieses Gates unerreichbar. Der Planner
hat die Frage aufbereitet und die Zahlen gemessen; die Entscheidung selbst
ist eine Architektur-/Prozess-Entscheidung und gehört in die Architect-Rolle
(`modul-08-agentenrollen.md` §Rollen-Regeln: „ADR-Änderung: Architect
schreibt").

**Rolleninhaber:** pt9912 (Architect-Zug; anderer Kontext als der
Planner-Lauf, der die 80-%-Frage aufbereitet und die Inventur gemessen hat,
und als die Implementer-/Reviewer-/Verifier-Läufe von `slice-049`/`slice-076`)

**Datum:** 2026-09-15

**Bezug:**
[`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(dieses Zugs Entscheidung — Supersedes `ADR-0054`, teilweise) ·
[`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
(in drei Klauseln korrigiert) ·
[`ADR-0030`](../plan/adr/0030-testpyramide.md) (Unit- vs. Integrations-Tier) ·
[`harness/sensors/coverage-gate.md`](../../harness/sensors/coverage-gate.md)
§Vertrag/§Kalibrierungs-Bindung/§Grenze · `harness/mk/coverage.mk` ·
`Makefile` (`test-store`/`test-replication`) ·
`.github/workflows/e2e.yml` · `AGENTS.md` §3.1, §3.5, §3.6, §3.7 ·
`docs/plan/planning/in-progress/slice-076-coverage-gate-reifestufe-40.md` <!-- d-check:status-provenance -->
(dessen Kalibrierung dieser Schnitt neu bemisst).

**Erzeugtes Artefakt dieses Zugs:** die Entscheidung `ADR-0071` + ihr
Index-Eintrag; die Folgearbeit (Zielrolle, unten) ist als Adresse benannt,
nicht als Slice angelegt — Priorisierung und Schnitt führt der Planner.

---

## Frage

Zwei Fragen, in dieser Reihenfolge:

1. **Liegt der Selbstwiderspruch vor?** Sagt `ADR-0054` §(a) mit der
   Endstufe **80 %** eine Zahl zu, die das von ihm selbst festgelegte
   Verfahren (netzloser Docker-Lauf über `./internal/... ./cmd/...`)
   strukturell nicht erreichen kann?
2. **Was ist dieses Gate** — eine Messung der netzlos prüfbaren Fläche, oder
   ein Proxy für System-Korrektheit? Die Antwort bestimmt, welche Option
   trägt.

## Befund — eigene Messung dieses Zugs

Der Planner hat geliefert; dieser Zug hat **nachgemessen**, nicht übernommen.
Verfahren: `docker build --no-cache-filter coverage --target coverage` mit dem
gepinnten Toolchain-Image `golang:1.27-alpine@sha256:cf6fca66…` (derselbe
Digest wie `TOOLCHAIN_IMAGE`/`Dockerfile`), danach Extraktion des rohen
Profils aus der gebauten Stage (`docker run --rm … cat /out/coverage.out`) und
eigene Auswertung, **dedupliziert über die Block-Position** (jede
Paket-Testbinary instrumentiert mit `-coverpkg` den ganzen Baum;
ungededupliziert zählt `numStmt` mehrfach).

| Größe | eigener Wert | Planner | Gate |
|---|---|---|---|
| Statements gesamt | **2467** | 2467 | — |
| netzlos gedeckt | **1215** | 1217 | — |
| Prozent (dedupliziert) | **49,25 %** | 49,33 % | **49,30 %** |

Die Spanne 49,25/49,30/49,33 liegt innerhalb der dokumentierten
Lauf-zu-Lauf-Schwankung (`verify-slice-049.md` §2: 39,6 % vs. 39,8 %) — die
zwei Statement-Differenzen sitzen ausschließlich in `internal/bootstrap`
(ein goroutine-naher Block, dessen Coverage-Lauf schwankt). **Die
Grundaussage ist damit unabhängig bestätigt.**

| Größe | Wert |
|---|---|
| Statements in den drei DB-gestützten Paketen (`postgresstorage` 610, `replication/receive` 155, `postgresack` 23) | **788** |
| Decke des heutigen Verfahrens (100 % von allem außer den drei Paketen) | **1679 / 2467 = 68,06 %** |

**80 % > 68,06 %.** Der Selbstwiderspruch liegt vor. Die drei Aussagen in
`ADR-0054` §(a) — Endstufe 80 %, Umfang `./internal/... ./cmd/...`, und der
eigene Satz, die DB-Adapter-Tests „zählen nicht gegen die Schwelle" — können
nicht gleichzeitig wahr sein.

Nebenbefunde, nachgeprüft: SQL ist bereits ausgelagert
(`postgresstorage/queries` trägt SQL-Text ohne ausführbare Statements, das
Unterpaket `mapper` ist mit 80 % gedeckt); die Naht zur Ausführung fehlt —
die Adapter hängen am konkreten `*pgxpool.Pool`, es gibt kein Interface.

## Verdikt 1 — Selbstwiderspruch: liegt vor → Folge-ADR mit `Supersedes`

Von den drei legitimen Ausgängen (Modul 8 §Konflikt-Pfad) trägt **Verdikt 2:
Folge-ADR mit `Supersedes`**:

- **Plan-Korrektur trägt nicht.** Der Fehler sitzt nicht in einem Slice-Plan,
  sondern im Text einer `Accepted`-ADR.
- **Undokumentierte Lockerung trägt nicht.** Es gibt keine Lockerung, die
  nachzuziehen wäre; die Endstufe ist nicht gesenkt, sondern unter einen
  falschen Gegenstand gestellt.
- **Folge-ADR trägt** (`AGENTS.md` §3.5: „Eine ADR mit Status `Accepted` wird
  nicht inhaltlich überschrieben. Korrekturen entstehen als neue ADR mit
  `Supersedes`"). Sie ist `ADR-0071` und supersedet `ADR-0054` in genau drei
  Klauseln (§(a) Scope-Bullet, §(a) Schwelle-Bullet, Fitness-Function-Zeile);
  alles Übrige — Docker-Stage, Gate-Skript, Suppression-Vollverbot,
  Eskalations-Mechanismus, §(b) — bleibt in Kraft.

`ADR-0054` selbst wird **nicht angefasst** (`AGENTS.md` §3.5).

## Verdikt 2 — Scope: das Gate misst die netzlos prüfbare Fläche

**Das Coverage-Gate ist eine Messung der netzlos prüfbaren Fläche, kein
Proxy für System-Korrektheit.** Ein netzloser Docker-Lauf kann nur
beurteilen, was netzlos läuft; die DB-gestützte Korrektheit hat ihren eigenen
Beleg-Träger (die Integrations-Ebene). Der Satz „80 %" bekommt damit einen
präzisen, erreichbaren Gegenstand:

- **Der Messgegenstand** wird auf `./internal/... ./cmd/...` **ohne** die
  Pakete geschnitten, deren Tests ohne externen Dienst überspringen
  (`postgresstorage` ohne `mapper`, `postgresack`,
  `replication/receive`). Tragende Regel ist die **Eigenschaft**, nicht die
  Liste: *Ein Paket, dessen Testlauf einen externen Dienst voraussetzt, ist
  nicht Gegenstand dieses Gates.*
- **Die Endstufe 80 % bleibt** — über der neuen Fläche (1679 Statements).
  Der Ist-Stand steigt auf **69,74 %** (dieselbe gedeckte Menge, kleinerer
  Nenner); die Rampe wird auf diesen Nenner neu kalibriert. Das ist **keine
  Schwellen-Senkung** (§3.6): nicht weniger wird verlangt, der Nenner wird
  ehrlich.
- **Die ausgenommene Fläche bekommt eine eigene Messung (Option F)** — dort,
  wo ihre Tests ohnehin gegen einen echten PostgreSQL laufen
  (`make test-store`/`make test-replication`) mit eigenem `-coverprofile`
  und eigener Schwelle, getragen vom **nicht-blockierenden**
  `.github/workflows/e2e.yml` — **nicht** in `make gates` (das läuft bei
  jedem Commit; Option C wäre der verworfenen Preis).
- **„Die Coverage" des Repos ist die Gate-getragene Zahl** (die Unit-Zahl);
  die DB-Adapter-Coverage trägt ihr Subjekt immer mit. Zwei Zahlen, zwei
  Gegenstände, ein Name je Gegenstand — keine Doppelquelle.

**Option D (die Naht) ist zulässig, aber nicht Teil dieser Entscheidung:** Sie
hebt die Decke ohne Messänderung und bessert das Design (schmale
Abhängigkeit, Fehlerklassifikation als reine Funktion). Ihre Rechtfertigung
ist **Design, nicht die Zahl** — als Zweck „die Zahl" wäre sie Goodhart und
damit unzulässig. Sie gehört als eigener, design-begründeter Vorgang nach
`next/`.

## Folgearbeit mit Zielrolle

Größenmaß ist der **Statement-Anteil** (Stmts gesamt / ungedeckt) aus der
Inventur; „Ziel" nennt die Rolle, die den Vorgang trägt.

| Vorgang | Zielrolle | Größenmaß | Art |
|---|---|---|---|
| **Coverage-Scope-Schnitt + Neukalibrierung** — `-coverpkg`/Profil-Filter ohne die drei Pakete, Rampe neu bemessen, Sensor-Doc/`README.md`/`AGENTS.md` §4 nachgezogen, Grün- und Rot-Beleg | Implementer (wellenlos) | kein Produkt-Code; Tooling/Doku | eigener Slice |
| **DB-Adapter-Coverage** — `-coverprofile` in `test-store`/`test-replication`, eigenes Target, Träger `e2e.yml`, Erstkalibrierung | Implementer (wellenlos) | Tooling | eigener Slice |
| **Wave „80 % auf der netzlos prüfbaren Fläche"** — Slices nach Paketgruppen; Schneide-Vorschlag unten | Planner schneidet, Implementer baut | +172 Statements nötig (1171 → 1343 von 1679) | Wave |
| **Executor-Naht (Option D)** — schmale Abhängigkeit, Logik pur ziehen; design-begründet, Zahl als Folge | Implementer | 8 Adapterdateien | eigener Vorgang in `next/` |

**Schneide-Vorschlag für die 80-%-Wave** (Paketgruppen, nach Statement-Anteil;
ungedeckt = der Hebel):

| Slice-Kandidat | Stmts gesamt | ungedeckt | Rolle im 80-%-Ziel |
|---|---|---|---|
| `internal/bootstrap` | **591** | **327** | der größte Einzelhebel; allein ~136 gedeckte Statements hier + `cmd` erreichen 80 % |
| `cmd/pg-change-feed` | 49 | 49 | klein, isoliert, 0 % — vollständig schließbar |
| Rest-Tail (`http` 28, `mapper` 25, `streamv1` 25, `decode` 14, `natsnotify` 5, Use-Cases) | ~400 | ~130 | füllt die Restlücke, wenn `bootstrap`+`cmd` nicht reichen |

Die Wave ist nicht zwingend, wenn `bootstrap` + `cmd` die 80 % tragen; ob sie
geschnitten wird, entscheidet der Planner an der realen Größe des
`bootstrap`-Anteils (bei > 3 Liefer-Punkten zerfällt `bootstrap` in zwei
Slices). Ein einzelner Slice je Paketgruppe ist einzeln lieferbar.

**Beobachtung zum laufenden `slice-076`:** dessen Hochschaltung (35 → 40 %)
ist unter dem **alten** Nenner korrekt und seine DoD trägt dort; nach dem
Schnitt aus `ADR-0071` bemisst der Scope-Slice die Rampe neu, und die Stufe
40 % ist gegen den neuen Nenner (~69,7 %) real trivial grün. Das ist ein
Übergangszustand, kein Defekt; `slice-076` wird von diesem Zug **nicht
angefasst**.

## Was dieser Zug geändert hat — und was nicht

**Geändert:**

- `docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md` (neu — die Entscheidung)
- `docs/plan/adr/README.md` (neue Zeile [`ADR-0071`](../plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md); [`ADR-0054`](../plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)-Zeile auf
  `(→ ADR-0071, teilweise)`)
- diese Verdikt-Datei

**Nicht geändert:** `docs/plan/adr/0054-…` (bleibt `Accepted`, unangetastet),
`harness/mk/coverage.mk`, `harness/sensors/coverage-gate.md`,
`harness/README.md`, `AGENTS.md` §4, `Makefile`, `.github/workflows/e2e.yml`,
`internal/**`, `spec/**`, jeder Slice-Plan (auch `slice-076`). Die Umsetzung
ist Folgearbeit, nicht Teil dieses Zugs.
