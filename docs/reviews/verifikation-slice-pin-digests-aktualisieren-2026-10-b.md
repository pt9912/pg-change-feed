# Verifikations-Report: slice-pin-digests-aktualisieren-2026-10-b — 2026-10-03

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität
([`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md) Festlegung 2 und 6,
[`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7) + Plan-vs-Code-Diff + Gates.
Review-Artefakt:
[`review-slice-pin-digests-aktualisieren-2026-10-b.md`](review-slice-pin-digests-aktualisieren-2026-10-b.md)
(Commit `26b953b1`, Haken-Commit `7fad7000`; 0 HIGH/0 MEDIUM, F-1 LOW, F-2 bis F-5 INFO). Formvorbild:
[`verifikation-slice-pin-digests-aktualisieren-2026-10.md`](verifikation-slice-pin-digests-aktualisieren-2026-10.md).

**Gegenstand:** Slice-Plan `slice-pin-digests-aktualisieren-2026-10-b` (wellenlos, Bezug
[`LH-QA-POR-001`](../../spec/lastenheft.md), [`LH-QA-POR-003`](../../spec/lastenheft.md)), Diff-Range
`7bc4aadd..HEAD` (`7fad7000`): Implementer-Commits `b010e244` (nats), `3912b9f0` (dotnet/sdk), `21d1d168`
(dotnet/runtime), `cc2cbddc` (temurin jdk), `f763bb1c` (temurin jre), `457d9fc6` (python), `f6acabc0` (trivy),
`a4134a96` (Kommentare postgres:18), `1c298431` (postgres:17), `726c085f` und `cfd818c7` (Plan), Review
`26b953b1`, Review-Haken `7fad7000`. Dieser Lauf ändert weder Code noch Plan; er schreibt nur diesen Report.
Es wurden keine DoD-Häkchen gesetzt und nicht gepusht.

**Methodenvermerke.**

1. *Scratchpad-Altlast:* das Scratchpad enthielt `*.exit`-/`ALLDONE`-Dateien früherer Läufe unter denselben
   Namen. Ich habe zwei frische Unterverzeichnisse (`v2`, `v3`) verwendet; ein erster Lauf im Wurzelverzeichnis
   wurde abgebrochen und nicht gewertet. Jeder Exit-Wert unten stammt aus `v2`/`v3` und trägt Start-/Endzeit.
2. *Build-Cache:* die Läufe in `v2` (`make image`, `make sdk-pack-*`, `make examples-*`) trafen den
   Docker-Build-Cache des Messhosts (0 bis 4 s Laufzeit). Ein Cache-Treffer trägt „diese Eingaben wurden mit dem
   neuen Pin schon einmal gebaut“, nicht „die Tests liefen jetzt“. Deshalb habe ich die fünf Bauketten
   zusätzlich mit `docker build --no-cache` (gleiche Argumente wie die Make-Ziele, ohne `-t`) neu gefahren
   (§1, Zeilen `nc-*`).
3. *Docker-Hub-Abruflimit:* ein früher Stapel von Registry-Abfragen lief in HTTP 429. Was dadurch nicht
   gelesen werden konnte, steht in §3 als nicht gemessen.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped, [`AGENTS.md`](../../AGENTS.md) §3.9)

Messhost: `uname -m` = `x86_64`, Linux 6.8.0-139-generic.

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make image` | Exit 0 | Lauf bis `#15 DONE` (Cache-Treffer, 0 s) |
| `make sdk-pack-csharp` | Exit 0 | Cache-Treffer (2 s); uncached: `docker build --no-cache … --target pack-export sdks/csharp` Exit 0, `Passed! - Failed: 0, Passed: 201, Skipped: 0, Total: 201` |
| `make sdk-pack-python` | Exit 0 | Cache-Treffer; uncached Exit 0, `195 passed, 5 warnings` |
| `make sdk-pack-kotlin` | Exit 0 | Cache-Treffer; uncached Exit 0, `BUILD SUCCESSFUL in 28s` (Tests und Publish-Probe in der Kette) |
| `make examples-csharp` | Exit 0 | Cache-Treffer; uncached (Hauptstufe) Exit 0, `Passed!` 54 (Http), 25 (Sse), 9 (Nats), 52 (Grpc), 20 (NatsStream) |
| `make examples-kotlin` | Exit 0 | Cache-Treffer; uncached Exit 0, `BUILD SUCCESSFUL in 24s` und `in 5s` |
| `make test-notify` | Exit 0 | letzte Zeile `ok … /driven/natsstream 2.132s`, 0 Zeilen `FAIL`; Lauf gegen den neuen NATS-Pin |
| `make test-store` | Exit 0 | `DB-Adapter-Coverage: 83.03% (gedeckt 964 von 1161 Statements; Profile gemergt: store,replication)`; `db-coverage: OK` |
| `make test-replication` (PG 18, Standard-Pin) | Exit 0 | `PostgreSQL 18.6: Keepalive inmitten der Transaktion nach 11531 von 400000 Änderungen …`; `--- PASS: TestSourceKeepaliveInsideTransactionDeliversWholeTransaction (40.41s)` |
| `make test-replication PG_TEST_IMAGE=postgres:17-alpine@sha256:b0f9560a…2b24` (Pin aus `e2e.yml`) | Exit 0 | `PostgreSQL 17.11: Keepalive inmitten der Transaktion nach 11078 von 400000 Änderungen …`; `--- PASS … (40.48s)` |
| `make test-sdk-csharp-integration` | Exit 0 | Schlusszeilen `Regel-Belege … grün`, `Routing-Belege … grün`, `Fehlercode-Belege … grün`, `Filter-Belege … grün` |
| `make test-sdk-python-integration` | Exit 0 | dieselben vier Schlusszeilen |
| `make test-sdk-kotlin-integration` | Exit 0 | dieselben vier Schlusszeilen |
| `make test-integration` | Exit 0 | 0 Zeilen `FAIL`; `E2E-Abdeckungstabelle unverändert …`; `Lauf abgeschlossen — E2E-Abdeckungstabelle aus 21 Go-Zeilen und 54 Bash-Zeilen`; danach `docker ps -a` ohne Test-Container (nur das fremde `gitea`), `git status --short` leer |
| `make gates` | Exit 0 | `baseline-verify: v6.13.0 OK — 54 Dateien`; `d-check: 1630 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability: OK — 5 Commit(s)`; `coverage-gate: OK — Coverage 82.40% erfüllt Schwelle 80%`; `generated-sync: OK`; `gesamt: 0 Befund(e)` (a-check); `sdk-public-doc-check`, `handbuch-public-doc-check`, `ausgabe-kennungen-check`, `meldungscodes-check` (97 Codes) ohne Befund |
| `make docs-check` | Exit 0 | `d-check: 1630 Datei(en) geprüft, 0 Befund(e)` (vor Anlage dieses Reports) |
| `make test` | Exit 0 | 51 `ok`-Zeilen, 0 `FAIL` |
| `make fmt-check` | Exit 0 | `fmt-check: 334 Go-Dateien geprüft, alle formatiert` |
| `make doc-immutable RANGE=7bc4aadd..HEAD` | Exit 0 | `0 Befund(e)` |
| `make doc-commits RANGE=7bc4aadd..HEAD` | Exit 0 | `0 Befund(e)` |
| `make commit-traceability` | Exit 0 | `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `make suchlauf-nachmessen PLAN=…` | Exit 0 | `suchlauf-nachmessen: 24 Zeilen stimmen` (12 Parent-Zeilen am `39e27242`, 12 `diff`-Zeilen) |
| `make doc-trace` | Exit 0 | `80 Anforderung(en), 0 Waise(n).` |
| `make image-stale` | Exit 0 | `OK golang:1.27-alpine == sha256:8a5910f3…`; `OK gcr.io/distroless/static-debian12:nonroot == sha256:afa5c872…` |
| `make pin-stale-race` / `-pgtest` / `-dmigrate` / `-acheck` / `-dcheck` / `-baseline` / `-actions` | je Exit 0 | je `OK` (`TOOLCHAIN_RACE_IMAGE`, `PG_TEST_IMAGE == 77f58511…`, `D_MIGRATE_IMAGE`, `A_CHECK_IMAGE`, `DCHECK_DIGEST` und Tag-Frische, Kurs-Baseline v6.13.0, vier Actions je Tag-Mutation und Tag-Frische) — die Nachher-Messung P1 bis P9 ist reproduziert |
| `make image-cve` (neuer Pin `af6acf9a…`) | Exit 0 | `debian 12.15: 0`, `gobinary: 0` Befunde |
| `make image-cve TRIVY_IMAGE=aquasec/trivy@sha256:62b1e65e…` (Parent-Pin aus `git show 7bc4aadd:Makefile`) | Exit 0 | `debian 12.15: 0`, `gobinary: 0` — **identischer Befund** am Parent-Werkzeug |

**Nicht gefahren / nicht gelesen:** die Vorher-Messung P1 bis P9 am Parent (der Drift ist behoben, das Vorher ist
nicht mehr herzustellen; reproduziert ist nur das Nachher); die Registry-Form (`mediaType`) von
`eclipse-temurin:21-jre`, `python:3.14-slim`, `postgres:17-alpine` und `postgres:18-alpine` (HTTP 429, siehe §3).
Der Post-Push-Lauf von `e2e.yml` und `upstream-drift.yml` (§7) steht aus.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

| # | DoD-Zeile | Haken | Verdikt | Realer Beleg |
|---|---|---|---|---|
| 1 | Messung (Liefer-Punkt 1) | `[ ]` | **korrekt offen** | Vorab-Messung P1 bis P9 fehlt (Abruflimit, im Plan §7 ehrlich vermerkt); nur das Nachher ist belegt und von mir reproduziert (§1). Die acht Drift-Messungen stehen als Tabelle in §7, nicht als gedruckte Zeilen (V-1) |
| 2 | Hebung und Nachzug (Liefer-Punkt 2) | `[x]` | **bestätigt** | neun Pin-/Kommentar-Commits, Betreffs nennen [`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md), kein `SPEC-*`/`ARC-*`; Träger-Gleichheit §3, Kommentar-Suchlauf §4 |
| 3 | Belege (Liefer-Punkt 3) | `[ ]` | **korrekt offen** | alle lokalen Läufe von mir bestätigt (§1); offen bleibt allein der Post-Push-Lauf (§7) — der Haken darf erst danach fallen |
| 4 | `make gates` / `make docs-check` / `make test` | `[x]` | **bestätigt** | eigene Läufe (§1); `82.40%` und `gesamt: 0 Befund(e)` decken sich mit dem Plan |
| 5 | Review durchgeführt | `[x]` | **bestätigt** | Report liegt vor; seine als übernommen geführten Langläufe (F-5) habe ich selbst gefahren (§1) |
| 6 | Doku-Update: kein öffentlicher Vertrag berührt | `[x]` | **bestätigt** | `git diff 7bc4aadd HEAD --stat -- spec docs/user harness README.md AGENTS.md` ist leer (§5); `sha256:` in `spec`, `docs/user`, `README.md`, `AGENTS.md`: 0 Treffer (Suchlauf-Zeile) |
| 7 | Closure-Notiz mit Lerneintrag | `[ ]` | **korrekt offen** | §7 trägt „(bei Closure)“ |
| 8 | Beobachtungs-Register | `[ ]` | **korrekt offen** | Vorschläge §8 |
| 9 | Risiken §6 mit Ausgang | `[ ]` | **korrekt offen** | alle sieben Zeilen „offen bis Closure“; Vorschläge §8 |
| 10 | Drei Paarungen | `[ ]` | **korrekt offen** | Closure-Pflicht |

Kein `[x]` ohne Beleg; kein `[ ]`, das bereits belegt wäre, außer der Messung selbst (nicht herstellbar, s. o.).

## 3. Pins: Index-Form, Auflösbarkeit, Träger-Gleichheit, Stopp-Regel

`docker buildx imagetools inspect <Tag> --format '{{.Manifest.Digest}}'` am Messtag, verglichen mit den Pins im
HEAD; Träger-Zählung `git grep -h -o -E '<image>@sha256:[0-9a-f]{64}' -- . ':!docs' ':!.harness'`.

| Image | Tag-Digest = Pin (Präfix) | Träger (Zeilen, verschiedene Werte) | Form des Pins | Version alt → neu (gemessen: `--version`) |
|---|---|---|---|---|
| `nats:2-alpine` | `ac8f88a6` ja | 3, 1 (`compose.yaml`, `examples/compose.yaml`, `run-notify-tests.sh`) | Index (`mediaType` `image.index.v1`) | nats-server v2.14.6 → v2.15.0 |
| `dotnet/sdk:10.0` | `e70cdb7f` ja | 2, 1 | Index (`manifest.list.v2`) | 10.0.401 → 10.0.401 (übernommen aus §7) |
| `dotnet/runtime:10.0` | `b89586dc` ja | 5, 1 | Index (`manifest.list.v2`) | 10.0.12 → 10.0.12 (übernommen aus §7) |
| `eclipse-temurin:21-jdk` | `3e3c176f` ja | 2, 1 | Index (`image.index.v1`) | 21.0.12+8 → 21.0.12.1+1 (neu: `java --version` am JRE-Pin gemessen; JDK-Alt übernommen) |
| `eclipse-temurin:21-jre` | `cff19e62` ja | 5, 1 | Index **hergeleitet** (429 beim Lesen; Tag-Digest eines Multi-Arch-Tags ist der Index-Digest) | 21.0.12+8 → 21.0.12.1+1 (neu gemessen: `Temurin-21.0.12.1+1`) |
| `python:3.14-slim` | `0741d101` ja | 1, 1 | Index **hergeleitet** (429) | 3.14.7 → 3.14.8 (neu gemessen: `Python 3.14.8`) |
| `aquasec/trivy` (`:latest`) | `af6acf9a` ja | 1, 1 | Index (`image.index.v1`) | 0.74.0 → 0.75.0 (Trivy-Hinweis-URL im Lauf nennt `v0.75`) |
| `postgres:17-alpine` | `b0f9560a` ja | 1, 1 (`e2e.yml`) | Index **hergeleitet** (429) | 17.11 → 17.11 (gemessen: `PostgreSQL 17.11` im Replikationstest) |
| `postgres:18-alpine` (Kontrolle, nicht gehoben) | `77f58511` ja | 8, 1 | Index **hergeleitet** | 18.6 (gemessen im Replikationstest) |

Alte Digests in lebenden Trägern: **keine** (`git grep -l` über acht Alt-Präfixe ohne `docs/reviews`, `done/`,
`observations/`, `.harness/baseline` trifft nur `Accepted` ADRs ([`ADR-0055`](../plan/adr/0055-nats-change-notification-wecksignal.md)
und weitere) und die Plan-Datei).

**Stopp-Regel (Urteil, [`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md) Festlegung 6):**
Kein Pin wechselt eine Hauptversion oder die vom Tag genannte Linie.
- *nats* 2.14.6 → 2.15.0: der Tag `2-alpine` nennt nur die Major-Linie 2; die Minor-Hebung liegt innerhalb der
  Linie — **Regel nicht ausgelöst**, Behandlung als Normalfall regelkonform (F-4 des Reviews bestätigt). Das
  Laufzeitverhalten trägt: `make test-notify`, `make test-integration` und die drei SDK-Integrationsläufe
  (NATS-Vollinhalt und Routing-Subjekt) liefen grün gegen v2.15.0.
- *temurin* und *python*: Patch-Wechsel, Normalfall.
- *trivy* 0.74 → 0.75: Pin gegen `:latest`, keine Linie genannt; Befund-Vergleich gegen den Parent (§1) ist
  **gleich** (0/0 in beiden Zielen). Damit ist der in F-4 des Reviews offene Parent-Vergleich geschlossen.

## 4. Kommentare zur Digest-Gewinnung (Auftrag 4)

`git grep -n -i "manifest inspect\|amd64" -- . ':!docs' ':!.harness' ':!.claude'` am HEAD: Treffer nur in
`Dockerfile`/`Makefile` (Plattformliste `linux/amd64,linux/arm64`, kein Gewinnungsweg), README.md/README.de.md
(Distribution), `spec/pflichtenheft.md` (Cross-Compile), `tools/harness/image-stale.sh` (echter
`docker manifest inspect`-Aufruf auf einen Major-Tag, kein Pin) und `run-dockerfile-from-tests.sh`
(Test-Fixture) — **kein** Gewinnungsweg eines Pins. Die Plan-Suchlauf-Zeilen `manifest inspect` (Soll 0) und
`amd64/linux|…` (Soll 0) stimmen. Die neuen Kommentare in `e2e.yml` (Zeilen 49–52, 76–78), `Makefile`
(Zeilen 91–92, 184–187), den Dockerfile-Köpfen und den drei `run-*-tests.sh` nennen übereinstimmend
`docker buildx imagetools inspect <Tag> --format '{{.Manifest.Digest}}'`. Den Weg habe ich an neun Referenzen
selbst gefahren (§3); er druckt für jede den Digest, den der Träger führt. Die Kommentare sind Indikativ-Zustand
(keine Chronik), konsistent mit [`AGENTS.md`](../../AGENTS.md) §3.7. Das Datum „ermittelt 2026-10-03“ in `e2e.yml`
ist der Messtag der Implementer-Messung; mein Lauf am selben Tag bestätigt ihn.

## 5. Plan-vs-Code-Diff (`git diff 7bc4aadd HEAD`)

14 Dateien: zwölf Pin-/Kommentar-Träger (`e2e.yml`, `Makefile`, `compose.yaml`, `examples/compose.yaml`, fünf
Dockerfiles unter `sdks/` und `examples/`, drei `run-*-tests.sh`) und zwei Dokumente (Plan, Review-Report). Plan §3
nennt exakt diese Menge; keine Zusatzdatei, kein Produktions- oder Testcode.

- **`e2e.yml` (Auftrag 5):** der Nicht-Kommentar-Diff besteht aus **einer** Zeile — dem Wert
  `pg_test_image` des PostgreSQL-17-Legs (`aa90e97e…` → `b0f9560a…`); alles übrige sind Kommentare. Matrix,
  Trigger, Permissions, Actions unverändert. PG-18-Leg `77f58511…` ist byte-gleich `Makefile` `PG_TEST_IMAGE`
  (Zeile 203) und `compose.yaml`-Default (`git grep`: acht Träger, ein Wert).
- **§3.6 (Auftrag 6):** `git diff 7bc4aadd HEAD --stat -- harness .d-check.yml .a-check.yml spec docs/user` ist
  **leer**; keine Schwelle, kein Modul, keine Regel berührt (`harness/mk/coverage.mk` inbegriffen).
- **Keine SDK-Versionsbumps:** der Diff der Dateien `PgChangeFeed.Client.csproj`, `pyproject.toml`,
  `build.gradle.kts` und `docs/user/version.md` ist leer (0 Zeilen). Kein `sdk-*-v*`- oder `v*`-Tag.
- **Doku-Verträge:** der Plan fordert, dass `harness/README.md` und die Sensor-Verträge die Digest-Gewinnung der
  acht Pins nicht als Einzelplattform-Lesung führen. `git grep -n -i "manifest inspect\|amd64"` über
  `harness/README.md` und `harness/sensors/*.md` zeigt keine solche Lesung; die Zeile zu `make image` nennt
  `linux/amd64,linux/arm64` als Build-Plattformliste, nicht als Pin-Gewinnung. Tragend.

## 6. Entscheidungs-Konformität

- [`ADR-0146`](../plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md) Festlegung 2 (Index-Digest
  als Pin-Form): erfüllt für alle acht gehobenen Referenzen (fünf am `mediaType` gelesen, drei hergeleitet, §3).
- Festlegung 6 (Tag-/Major-Wechsel nicht in diesem Zug): eingehalten (§3).
- Reihenfolge „erst Hebung, dann Sensor“ (Folge-Slice `slice-pin-stale-alle-digest-pins` liegt in `open/`):
  eingehalten; der Sensor-Slice startet erst nach der Closure dieses Slice.
- [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7 („Digest-Hebung bleibt bewusster
  Commit“): je Pin ein eigener Commit, Betreffs mit Anker. Keine neue ADR nötig.

## 7. §3.10 — Post-Push-Läufe (Auftrag 9; nur benannt)

`e2e.yml` ist berührt (PG-17-Pin und Kommentare). Ein lokal grüner Lauf zeigt nicht, dass der Runner grün läuft.
**Offen, gehört dem Hauptlauf nach dieser Verifikation:** Push, danach `gh run list --workflow e2e.yml` und
`gh run view` für **beide** Matrix-Legs (PG 17 mit `b0f9560a…`, PG 18 mit `77f58511…`) und ein
`workflow_dispatch` von `upstream-drift.yml` (alle neun Achsen `OK`, kein roter Schritt). Bis dahin bleibt das
§6-Risiko „CI-Matrix `e2e.yml`“ des Plans *weiter offen*, und der Haken von Liefer-Punkt 3 darf nicht fallen.
Hinweis: ein Index-Digest löst auf dem amd64-Runner auf dasselbe Image auf, das ich auf dem `x86_64`-Messhost
lokal gefahren habe (PG 17.11 Replikationstest grün); am Runner selbst ist das nicht gemessen.

## 8. Vorschläge für Closure und Register (Auftrag 8, nur Vorschläge)

1. **`pin-ohne-inventar-eintrag-driftet-unsichtbar`** (Zähler 1, abgeleitet aus `evidence/`): dieser Slice ist der
   zweite reale Beleg derselben Klasse und der erste für die fünf Basis-Images der SDK-/Beispiel-Dockerfiles, `nats`
   und `trivy` — jeweils driftend und mit `make image-stale`/`make pin-stale-*` nicht gesehen (alle acht Referenzen
   wichen vom Index-Digest ihres Tags ab, vom Planner und vom Implementer unabhängig gemessen, von mir am Endstand
   als gleich bestätigt). Vorschlag: weitere Datei `evidence/slice-pin-digests-aktualisieren-2026-10-b.md`; der
   Zähler folgt aus den Dateien (dann 2).
2. **`nicht-blockierender-workflow-alarmmuedigkeit`** (zwei Evidenz-Dateien): dieser Slice liefert keinen neuen
   roten Nachtlauf (nichts gelesen; `upstream-drift.yml` nicht abgefragt). Kein Beleg zu zählen; nicht anfassen,
   bis der Post-Push-Lauf gelesen ist.
3. **Neue Beobachtung (Kandidat, Planner entscheidet):** *Der Pin-Kommentar nennt den falschen Gewinnungsweg* —
   vor dem Slice trugen sieben lebende Träger „`docker manifest inspect` (amd64)“ als Gewinnungsweg, der
   Einzelplattform-Digests liefert (Parent-Suchlauf Zeile 11: 7 Treffer). Der Slice hat das behoben; ein Sensor
   dafür ist der Folge-Slice. Kein eigener Eintrag nötig, wenn Punkt 1 den Satz trägt.
4. **§6-Ausgänge (Vorschlag je Risiko):**
   - Hebung ändert das Verhalten eines Bau-/Testlaufs → *entfallen* (alle Läufe grün, uncached bestätigt; Trivy-Befund
     gegen den Parent gleich).
   - Einzelplattform-Pin und Plattform des Bau-Hosts → *entfallen* (Messhost `x86_64`, `uname -m` gemessen; die
     Plattform-Auflösung folgt dem Host, kein Plattform-Wechsel gegenüber dem Parent beobachtet).
   - Kopie des Pins auf altem Digest → *entfallen* (§3: je Image ein Wert; `suchlauf-nachmessen` 24 Zeilen Exit 0).
   - Digest bewegt sich zwischen Messung und Push → *weiter offen* (Folge-Slice misst am Start neu; Träger: dessen
     erste Messung).
   - CI-Matrix `e2e.yml` → *weiter offen* bis der Post-Push-Lauf gelesen ist (§7); Träger der Nachverfolgung bleibt
     `github-actions-unverifizierbar-lokal`.
   - `make image-cve` mit neuem Pin → *entfallen* (Exit 0, 0/0, Parent gleich).
   - Alarmmüdigkeit → *weiter offen* (Eintrag im Register, Zähler 2, unverändert).
5. **Release-Folge (nur melden):** die Hebungen in `sdks/*/Dockerfile` (csharp, python, kotlin) wirken erst mit dem
   nächsten SDK-Release (`sdk-csharp-v*`, `sdk-python-v*`, `sdk-kotlin-v*`), weil die Release-Workflows über
   `make sdk-pack-*` bauen. Das Server-`Dockerfile` ist unberührt (kein Pin des Servers gewechselt), der nächste
   Server-Release (`v*`) trägt damit keine Basis-Image-Änderung aus diesem Slice. Kein Release ist Teil dieses
   Slice, keine Versionsdatei berührt.

## 9. Abweichungen und Findings des Verifiers

| # | Klasse | Schwere | Befund |
|---|---|---|---|
| V-1 | Zahl im Träger ohne gedruckte Messzeile (F-1 des Reviews) | LOW | **Bestätigt, mit Einordnung nach [`AGENTS.md`](../../AGENTS.md) §3.12.** Plan §7 trägt die acht Pin-Messungen als Tabelle mit Befehlsangabe und Parent-Kennung (`7bc4aadd`), die **neuen** Digests voll ausgeschrieben, die **alten** gekürzt (`065e8355…ccc`). Der Ursprung „gemessen“ ist genannt, die gedruckte Zeile fehlt; der Liefer-Punkt 1 verlangt sie ausdrücklich. Die Neu-Seite ist von mir reproduziert (§3): alle acht Werte stimmen, und die Versionen nats v2.14.6/v2.15.0, Temurin 21.0.12.1+1, Python 3.14.8 und PostgreSQL 17.11 sind nachgemessen. Die Alt-Digests sind vollständig aus dem Parent herstellbar (`git show 7bc4aadd:<Datei>`). **Nicht** nachgemessen und im Plan als gemessen geführt: die Alt-Form der Pins („Einzelplattform“, `docker.distribution.manifest.v2`) und die Alt-Versionen von .NET und JDK — sie sind am Endstand nicht mehr herstellbar und im Plan als **übernommen** zu kennzeichnen, nicht als gemessen. Nachtrag durch den Planner in §7 (ein Satz), kein Code |
| V-2 | Parent-Kennung uneinheitlich (F-2 des Reviews) | INFO | **Erklärbar, kein Befund.** `39e27242` ist der Planungsstand, an dem Planner den Suchlauf mass (Parent des Plan-Commits `18d34cbb`); `7bc4aadd` ist der Stand, an dem die Implementer-Commits ansetzen (Plan-Move und Verweis-Nachzug `1eca9b35`, `7bc4aadd`). Zwischen beiden liegen nur Dokument-Commits: `git diff --stat 39e27242 7bc4aadd` nennt drei Dateien unter `docs/plan/planning/` (Plan, Register-`state.md`, Folge-Plan), **kein** Code und keine Konfiguration. Die Parent-Zeilen des Suchlaufs (12) stimmen am Planungsstand, die Hebungen am Implementer-Stand. Für die Closure-Notiz ein Satz genügt |
| V-3 | Vorab-Messung P1 bis P9 unvollständig (F-3 des Reviews) | INFO | **Bestätigt, ehrlich vermerkt.** Ausgang: Haken `[ ]` zu Recht offen; das Nachher ist reproduziert (§1). Der Plan-Satz „die sieben Pins stehen auf dem Stand des Messtags“ ist am Endstand bestätigt, am Parent nicht gemessen |
| V-4 | Stopp-Regel (F-4 des Reviews) | INFO | Urteil in §3: nats-Minor innerhalb der Linie, Regel nicht ausgelöst; Trivy-Befund gegen den Parent jetzt gemessen und gleich |
| V-5 | Langläufe als übernommen (F-5 des Reviews) | INFO | **Geschlossen:** alle in F-5 genannten Läufe von mir selbst gefahren, Exit 0 (§1); die Bauketten zusätzlich uncached |
| V-6 | Registry-Form nicht überall gelesen | INFO | `mediaType` von drei gehobenen Referenzen (`jre`, `python`, `postgres:17`) und der Kontrolle `postgres:18` wegen HTTP 429 nicht gelesen (§3, „hergeleitet“); die Gleichheit Pin = Tag-Digest ist gemessen. Mit dem Folge-Slice (Sensor) wird die Form ohnehin je Lauf geprüft |

Keine HIGH-, keine MEDIUM-Befunde. Keine DoD-Verletzung.

## 10. Verdikt

**bestanden** — mit den folgenden Bedingungen für die Closure (keine verlangt eine Änderung am Pin-Bestand):

1. **Post-Push-Läufe lesen** ([`AGENTS.md`](../../AGENTS.md) §3.10): `e2e.yml` beide Legs grün und
   `upstream-drift.yml` (`workflow_dispatch`) ohne `DRIFT` und ohne roten Schritt; Ergebnis in Plan §7. Bis dahin:
   Liefer-Punkt 3 `[ ]`, §6-Risiko CI-Matrix *weiter offen*.
2. **V-1:** Plan §7 ergänzen: Alt-Form und Alt-Versionen von .NET/JDK als *übernommen* kennzeichnen (oder die
   gedruckten Parent-Zeilen aus `git show` nachtragen); ein Satz zu V-2 (zwei Parent-Kennungen).
3. **Register/§6:** Evidenz-Datei zu `pin-ohne-inventar-eintrag-driftet-unsichtbar` (§8.1); je Risiko ein Ausgang
   (§8.4); Lerneintrag in der Closure-Notiz (Kandidat: ein Gewinnungsweg-Kommentar, der den Einzelplattform-Digest
   nennt, hält die falsche Pin-Form am Leben — die Hebung braucht Wert **und** Kommentar im selben Commit).
4. **Folge-Slice** `slice-pin-stale-alle-digest-pins` erst nach der Closure dieses Slice starten.

Offene Punkte für den Planner: 1 bis 4 oben; das Setzen der verbleibenden DoD-Häkchen (Liefer-Punkte 1 und 3,
Closure-Zeilen) liegt beim Planner nach dem Post-Push-Lauf.
