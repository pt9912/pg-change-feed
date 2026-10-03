# Slice pin-digests-aktualisieren-2026-10-b: Gedriftete Digest-Pins ohne Inventar-Eintrag auf den Index-Digest heben

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung jenseits der DoD
dieses Slice (gemessen: `docs/plan/planning/` trägt keine flache Welle-Datei,
`ls docs/plan/planning/*.md` nennt nur `README.md`), siehe Baseline-Regelwerk
`modul-06-roadmap.md` §Wann Arbeit eine Welle braucht (Modul 6).

**Bezug:** [`LH-QA-POR-001`](../../../../spec/lastenheft.md),
[`LH-QA-POR-003`](../../../../spec/lastenheft.md) (Scope),
[`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
(Folgepflicht Reihenfolge, Slice 1; Festlegung 2: Index-Digest als Pin-Form),
[`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7
(„Digest-Hebung bleibt bewusster Commit“),
[`ADR-0058`](../../adr/0058-testansatz-fuenf-luecken.md).

**Berührte Spec-Stellen:** — ([`SPEC-012`](../../../../spec/pflichtenheft.md) nennt nur die Major-Versionen
17 und 18, keinen Digest; zu belegen durch den Suchlauf in §3, der keinen Digest in
`spec/` findet — Erwartung, am Parent von der Planung nicht gemessen). Eine
Spec-Berührung entstünde nur bei einem Tag- oder Major-Wechsel; das ist die
Stopp-Regel in §4, nicht Gegenstand dieses Slice.

**Verantwortlich:** —
<!-- BEDIENHINWEIS: Verantwortlich hält die Arbeit — der Rolleninhaber der
Implementer-Rolle, gesetzt beim Übergang open→next (Baseline-Regelwerk
modul-05-planning-harness.md §Lifecycle als State Machine). Der Autor schrieb
den Plan; zwei Felder, zwei Fragen. Kein Statuswert: der Zustand bleibt das
Verzeichnis. Kein Sensor prüft das Feld — es ist Deklaration. -->

**Autor:** Planner-Agent, direkt beauftragt (kein Architect: die Entscheidung
liegt mit [`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
vor). **Datum:** 2026-10-03.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar.

**Anlass (übernommen aus
[`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
Kontext 3, dort gemessen am Stand `28a6242a`):** acht Digest-Pins außerhalb des
Inventars P1–P9 driften oder tragen die falsche Form. **Vom Planner
nachgemessen** (2026-10-03, Arbeitsstand `39e27242`, je Referenz
`docker buildx imagetools inspect <tag> --format '{{.Manifest.Digest}}'` gegen
die gepinnten Werte; das ist eine Registry-Abfrage, keine Toolchain):

| Pin (Tag) | gepinnt (Präfix) | aktueller Index-Digest (Präfix) | Befund |
|---|---|---|---|
| `nats:2-alpine` | `065e8355` | `ac8f88a6` | weicht ab |
| `mcr.microsoft.com/dotnet/sdk:10.0` | `60a2b223` | `e70cdb7f` | weicht ab |
| `mcr.microsoft.com/dotnet/runtime:10.0` | `e6541e52` | `b89586dc` | weicht ab (in [`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md) Kontext 3 nicht gemessen) |
| `eclipse-temurin:21-jdk` | `085eb93e` | `3e3c176f` | weicht ab |
| `eclipse-temurin:21-jre` | `ca7551d4` | `cff19e62` | weicht ab |
| `python:3.14-slim` | `caaf356f` | `0741d101` | weicht ab |
| `aquasec/trivy` (gegen `:latest`) | `62b1e65e` | `af6acf9a` | weicht ab |
| `postgres:17-alpine` | `aa90e97e` | `b0f9560a` | Einzelplattform-Digest, Mitglied des aktuellen Index (`imagetools inspect --raw`, Treffer 1) — unzulässige Pin-Form |

Für die ersten sieben Pins ist **nicht gemessen**, ob der gepinnte Wert ein
Index-Digest oder ein Einzelplattform-Digest ist: er ist kein Mitglied des
*aktuellen* Index (`--raw`, kein Treffer), das trennt „alter Index“ nicht von
„alte Einzelplattform“. Die Kommentare in den Dockerfiles und Skripten nennen
„amd64“ und „`docker manifest inspect`“ (Suchlauf in §3), also die
Einzelplattform-Form. Die Digests sind Momentaufnahmen und werden mit jedem
Registry-Push anders; der Implementer misst im Liefer-Punkt 1 neu, die Werte
oben sind kein Soll ([`AGENTS.md`](../../../../AGENTS.md) §3.12).

**Ziel:** Jede der acht Referenzen trägt an jeder Fundstelle den Index-Digest
ihres Tags, die Kommentare zur Digest-Gewinnung nennen den Weg, der ihn liefert
(`docker buildx imagetools inspect`), und die Bau- und Testläufe belegen, dass
die gehobenen Basis-Images die SDKs, Beispiele und Tests gleich bewerten wie
vorher.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Sensor P10** — Folge-Slice
  [`slice-pin-stale-alle-digest-pins`](slice-pin-stale-alle-digest-pins.md); die
  Reihenfolge „erst Hebung, dann Sensor“ ist Teil der Entscheidung
  ([`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
  Frage (d)), der erste Lauf mit Sensor soll grün sein.
- **Ein Tag- oder Major-Wechsel** (`dotnet` 10.0 → 11.0, `eclipse-temurin` 21 →
  25, `python` 3.14 → 3.15, `postgres` 17/18, `-alpine`): das ist eine
  Änderung am Werkzeug-Stand
  ([`SPEC-012`](../../../../spec/pflichtenheft.md) für PostgreSQL) und kein
  Digest-Commit; Grenze von
  [`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
  Festlegung 6. Die Stopp-Regel steht in §4.
- **Versionsbumps der SDKs und jede Veröffentlichung** (`sdk-*-v*`-Tags,
  `v*`-Release): die Hebung ändert das Bau-Image, nicht die Quelle eines
  Package; ein Release braucht eine eigene Nutzerfreigabe. Die
  Versionsdateien (`.csproj`, `pyproject.toml`, `build.gradle.kts`) bleiben
  unberührt.
- **Pins, die P1 bis P9 bereits führen** (`golang`, `distroless`,
  `postgres:18-alpine`, `d-migrate`, `a-check`, `d-check`, `uv`): sie stehen
  seit `slice-pin-digests-aktualisieren-2026-10` auf dem Stand des Messtags;
  der Lauf von `upstream-drift.yml` auf `38c9ff93` war grün. Der Slice misst
  sie im Liefer-Punkt 1 einmal mit — zeigt eine Messung Drift, steht der
  Befund in §7 und die Achse wird als eigener Commit mitgehoben (derselbe
  Vorgang), sonst nicht angefasst.
- **Ein Alarm- oder Eigentümer-Mechanismus für künftigen Drift:** die
  Beobachtung
  [`BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit`](../observations/BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit/observation.md)
  bleibt Beobachtung
  ([`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
  Frage (d)).
- **Die in `Accepted` ADRs, `docs/reviews/`, `done/` und `observations/`
  stehenden alten Digests und Kommentare** (unter anderem `ADR-0055`,
  `ADR-0087`, `ADR-0093`): sie halten den Stand ihrer Messung
  ([`AGENTS.md`](../../../../AGENTS.md) §3.5); eine Zitat-Korrektur nach
  [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
  trifft sie nicht, weil der Referent unverändert ist. Der `diff`-Suchlauf in
  §3 führt sie als erwartete Resttreffer.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] **Messung (Liefer-Punkt 1).** Alle acht Referenzen aus §1 sind am
      Arbeitsstand des Implementers neu gemessen: je Referenz der
      Index-Digest (`docker buildx imagetools inspect <tag> --format
      '{{.Manifest.Digest}}'`) und die **Form des gepinnten Werts** (Index oder
      Einzelplattform: `docker buildx imagetools inspect <tag> --raw` auf den
      gepinnten Digest, und `docker buildx imagetools inspect
      <image>@<gepinnter Digest> --raw` gelesen auf `mediaType`). Dazu die
      Achsen-Messung der Pins P1 bis P9 (`make image-stale`,
      `make pin-stale-race|-pgtest|-dmigrate|-acheck|-dcheck|-baseline|-actions`)
      am selben Stand. Die gedruckten Zeilen mit vollen Digests und der
      Parent-Kennung stehen in §7; die Werte aus §1 sind **nicht** zu
      übernehmen. Zeigt die Messung einen Major-/Tag-Wechsel im Image
      (Stopp-Regel §4), endet der Slice dort.
- [ ] **Hebung und Nachzug (Liefer-Punkt 2).** Je Pin ein eigener Commit
      (Betreff nennt
      [`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md),
      keine `SPEC-`/`ARC-`-Kennung), jeder Pin auf den **Index-Digest** und an
      **jeder** Fundstelle, die der `git grep` am Arbeitsstand des Implementers
      findet (§3 nennt die am Parent gemessenen; maßgeblich ist der Befund am
      Arbeitsstand, nicht diese Liste): `nats:2-alpine` (`compose.yaml`,
      `examples/compose.yaml`, `tools/harness/run-notify-tests.sh` und weitere),
      `dotnet/sdk:10.0` (`sdks/csharp/Dockerfile`, `examples/csharp/Dockerfile`),
      `dotnet/runtime:10.0` (`examples/csharp/Dockerfile`, fünf Zeilen),
      `eclipse-temurin:21-jdk` (`sdks/kotlin/Dockerfile`,
      `examples/kotlin/Dockerfile`), `eclipse-temurin:21-jre`
      (`examples/kotlin/Dockerfile`, fünf Zeilen), `python:3.14-slim`
      (`sdks/python/Dockerfile`), `aquasec/trivy` (`Makefile`
      `TRIVY_IMAGE`) und `postgres:17-alpine`
      (`.github/workflows/e2e.yml`). Die Kommentare, die einen Digest oder
      seinen Gewinnungsweg nennen (Dockerfile-Köpfe, `e2e.yml`, `Makefile`,
      `tools/harness/run-*-tests.sh`), sind im selben Commit auf
      `docker buildx imagetools inspect` und die Index-Form gezogen — kein
      verbliebenes „`docker manifest inspect`“ und kein „amd64“ als Gewinnungsweg
      eines Pins in einem lebenden Träger (die Resttreffer-Erwartung steht in §3).
- [ ] **Belege (Liefer-Punkt 3).** Je gehobene Achse läuft der engste
      betroffene Lauf (§3 Reihenfolge) und seine Ausgabe wird gegen den Parent
      verglichen (Befunde, Meldungsform); `make sdk-pack-csharp`,
      `make sdk-pack-python`, `make sdk-pack-kotlin`, `make examples-csharp`,
      `make examples-kotlin`, `make test-notify`, `make test-integration`
      und die drei `make test-sdk-*-integration` enden mit Exit 0 (die
      Integrationsläufe setzen ein geladenes `:dev`-Image voraus: `make image`
      zuerst); danach läuft die Achsen-Messung aus Liefer-Punkt 1 erneut
      und druckt für jede gehobene Achse `OK`. Ein Lauf, der am Parent grün und
      am Diff rot ist, ist ein Befund (§6), kein Anlass, den Lauf zu lockern.
      Der reale Post-Push-Lauf von `.github/workflows/e2e.yml` (beide
      Matrix-Legs) und `.github/workflows/upstream-drift.yml`
      (`workflow_dispatch`) ist nach dem Push gelesen
      ([`AGENTS.md`](../../../../AGENTS.md) §3.10; der Push gehört nicht zu
      diesem Plan, die Lesung ist Closure-Pflicht).

Gate- und Lauf-Pflichten (zählen nicht zu den Liefer-Punkten):

- [ ] `make gates` grün (Exit direkt ausgewertet, am Endstand;
      [`AGENTS.md`](../../../../AGENTS.md) §3.9), `make docs-check` Exit 0,
      `make test` Exit 0.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: kein öffentlicher Vertrag berührt, zu belegen durch den
      `spec/`-/`docs/user/`-Suchlauf in §3 (Zeile 0 am `diff`-Stand);
      [`harness/README.md`](../../../../harness/README.md) §Sensors und die
      Sensor-Verträge unter `harness/sensors/` nennen die Digest-Gewinnung für
      keinen der acht Pins als Einzelplattform-Lesung — der Implementer liest
      die Zeilen zu `make test-store`, `make test-replication`,
      `make test-notify` und `make image-cve` und zieht nach oder meldet.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der Slice-Closure selbst, weil die Roadmap unter *Offene Wellen* keine Welle führt (gemessen: `docs/plan/planning/` trägt keine flache Welle-Datei) und „die nächste Welle-Closure“ damit keine Adresse ist. Die Folge-Slice-Adresse ist
      [`slice-pin-stale-alle-digest-pins`](slice-pin-stale-alle-digest-pins.md)
      (liegt in `open/`); sie nimmt die Sendung „erster Lauf mit Sensor grün“ an:
      ihr Plan trägt die Abhängigkeit „erst nach diesem Slice“.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

Fundstellen am Parent `39e27242` (gemessen mit `git grep -n -F '<image>@'`
je Pin, ohne `docs/reviews`, `done/`, `observations/`, `.harness/baseline`;
maßgeblich bleibt der Befund am Arbeitsstand des Implementers):

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `compose.yaml` (Zeile 49), `examples/compose.yaml` (56), `tools/harness/run-notify-tests.sh` (14, Kommentar Zeile 6) | update | `nats:2-alpine`: drei Fundstellen des Pins, eine des Gewinnungs-Kommentars („linux/amd64-Manifest“) |
| `sdks/csharp/Dockerfile` (25, 27), `examples/csharp/Dockerfile` (18, 34) | update | `dotnet/sdk:10.0`: Pin und Kommentar-Digest je Datei |
| `examples/csharp/Dockerfile` (20, 97, 105, 116, 125, 134) | update | `dotnet/runtime:10.0`: fünf `FROM`-Zeilen und ein Kommentar-Digest; Kopfkommentar (Zeile 14: „`docker manifest inspect`, amd64/linux“) |
| `sdks/kotlin/Dockerfile` (31, 33), `examples/kotlin/Dockerfile` (19, 38) | update | `eclipse-temurin:21-jdk`; Kopfkommentare (Zeile 27 bzw. 16) |
| `examples/kotlin/Dockerfile` (21, 81, 89, 100, 109, 118) | update | `eclipse-temurin:21-jre`: fünf `FROM`-Zeilen und ein Kommentar-Digest |
| `sdks/python/Dockerfile` (31, 33) | update | `python:3.14-slim` |
| `Makefile` (`TRIVY_IMAGE`, Zeile 92) | update | `aquasec/trivy` gegen `:latest` |
| `.github/workflows/e2e.yml` (Zeilen 49–50, 76–79) | update | PostgreSQL-17-Leg: Digest auf den Index-Digest, zwei Kommentare zur Gewinnung („`docker manifest inspect`“, „amd64“); Matrix-Struktur und PostgreSQL-18-Leg unverändert |
| `Makefile` (Zeilen 184–185), `tools/harness/run-store-tests.sh` (4), `tools/harness/run-replication-tests.sh` (15) | update | Kommentare zum Gewinnungsweg des `postgres:18-alpine`-Pins; der Pin selbst ist nach der Messung des Implementers Index oder wird mitgehoben (§1, Punkt „Pins, die P1 bis P9 führen“) |
| `docs/plan/planning/in-progress/slice-pin-digests-aktualisieren-2026-10-b.md` | update | §7 Messzeilen, Vergleich gegen den Parent (die Datei wandert beim Start dorthin) |

- **Reihenfolge:** (1) messen, Zeilen in §7; (2) je Pin ein Commit, davor den
  engsten Lauf: `nats` → `make test-notify` und `make test-integration`;
  `dotnet` → `make sdk-pack-csharp`, `make examples-csharp`; `eclipse-temurin` →
  `make sdk-pack-kotlin`, `make examples-kotlin`; `python` → `make
  sdk-pack-python`; `trivy` → `make image-cve` (advisory; scheitert es
  strukturell am Ziel-Image, steht das in §7); `postgres:17` → `make
  test-replication` mit `PG_TEST_IMAGE` auf dem Index-Digest und
  `make test-store`; (3) am Ende `make gates`, `make image`, `make
  test-integration`, die drei `make test-sdk-*-integration` und die
  Achsen-Messung. Eine Achse nach der anderen, damit ein roter Lauf einer Achse
  zuzuordnen ist.
- **Digest-Gewinnung:** ausschließlich `docker buildx imagetools inspect <tag>
  --format '{{.Manifest.Digest}}'` (Index-Digest,
  [`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
  Festlegung 2); der Implementer trägt den Wert aus der gedruckten Zeile ein,
  nicht aus diesem Plan. Ein Einzelplattform-Digest ist keine zulässige Form.
- **Stopp-Regel — Major-Tag-Wechsel und Hauptversion:** vor der Hebung eines
  Pins liest der Implementer die Version, die das neue Image meldet
  (`docker run --rm --entrypoint <Werkzeug> <image>@<neuer Digest> --version`
  bzw. das Label `org.opencontainers.image.version`, je Image das, was es
  trägt) und vergleicht sie mit der Version am gepinnten Digest. Ändert sich die
  **Hauptversion oder die Minor-Linie, die der Tag nennt** (`dotnet` 10.0,
  `temurin` 21, `python` 3.14, `nats` 2.x-Linie, `postgres` 17), hält der
  Implementer **diesen** Pin an, hebt die übrigen und meldet den Befund; ein
  Patch-Wechsel innerhalb der Linie ist der Normalfall dieses Slice. Ein
  Tag-Wechsel (`10.0` → `11.0`, `21` → `25`) wird nie in diesem Slice
  vorgenommen.
- **Suchlauf (§3.13, bewegte Eigenschaft: der Digest-Wert eines Pins und der
  Satz zu seiner Gewinnung).** Der Parent ist
  `39e27242eb58ed1212b381a15da55c6ee1f5033c` (`git rev-parse HEAD` am
  Planungsstand, vor dem Plan-Commit; nie `HEAD`). Suchraum: ganzer Baum ohne
  `docs/reviews`, `done/`, `observations/`, `.harness/baseline` und die offenen
  Pläne (Records, vendored Baseline und die Pläne, die die alten Werte als
  Suchmuster tragen). Jede Zeile ist am Parent gemessen (Soll = Trefferzeilen);
  die `diff`-Zeilen tragen die **Erwartung** nach der Hebung (hergeleitet aus
  der Aufteilung der Parent-Treffer in lebende Träger und Accepted-ADRs, noch
  nicht gemessen) — der Implementer misst nach und trägt Gefundenes und
  Nichtgefundenes in §7 ein:

```suchlauf
39e27242eb58ed1212b381a15da55c6ee1f5033c 38 -n -E '065e8355|60a2b223|e6541e52|085eb93e|ca7551d4|caaf356f|62b1e65e|aa90e97e' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
39e27242eb58ed1212b381a15da55c6ee1f5033c 5 -n -E '065e8355' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
39e27242eb58ed1212b381a15da55c6ee1f5033c 5 -n -E '60a2b223' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
39e27242eb58ed1212b381a15da55c6ee1f5033c 6 -n -E 'e6541e52' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
39e27242eb58ed1212b381a15da55c6ee1f5033c 11 -n -E '085eb93e' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
39e27242eb58ed1212b381a15da55c6ee1f5033c 6 -n -E 'ca7551d4' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
39e27242eb58ed1212b381a15da55c6ee1f5033c 2 -n -E 'caaf356f' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
39e27242eb58ed1212b381a15da55c6ee1f5033c 2 -n -E '62b1e65e' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
39e27242eb58ed1212b381a15da55c6ee1f5033c 2 -n -E 'aa90e97e' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
39e27242eb58ed1212b381a15da55c6ee1f5033c 0 -n 'sha256:' -- spec docs/user README.md AGENTS.md
39e27242eb58ed1212b381a15da55c6ee1f5033c 7 -n -E 'manifest inspect' -- . ':!docs' ':!.claude' ':!.harness' ':!tools/harness/image-stale.sh'
39e27242eb58ed1212b381a15da55c6ee1f5033c 10 -n -E 'amd64/linux|linux/amd64,? *$|amd64-Manifest|\(amd64|amd64\)|, amd64|inspect.*amd64|linux/amd64-' -- . ':!docs' ':!.claude' ':!.harness' ':!tools/harness/image-stale.sh'
diff 11 -n -E '065e8355|60a2b223|e6541e52|085eb93e|ca7551d4|caaf356f|62b1e65e|aa90e97e' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
diff 2 -n -E '065e8355' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
diff 1 -n -E '60a2b223' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
diff 0 -n -E 'e6541e52' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
diff 7 -n -E '085eb93e' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
diff 0 -n -E 'ca7551d4' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
diff 0 -n -E 'caaf356f' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
diff 0 -n -E '62b1e65e' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
diff 1 -n -E 'aa90e97e' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!docs/plan/planning/open' ':!.harness/baseline'
diff 0 -n 'sha256:' -- spec docs/user README.md AGENTS.md
diff 0 -n -E 'manifest inspect' -- . ':!docs' ':!.claude' ':!.harness' ':!tools/harness/image-stale.sh'
diff 0 -n -E 'amd64/linux|linux/amd64,? *$|amd64-Manifest|\(amd64|amd64\)|, amd64|inspect.*amd64|linux/amd64-' -- . ':!docs' ':!.claude' ':!.harness' ':!tools/harness/image-stale.sh'
```

  **Erwartung am `diff`-Stand** (hergeleitet): die Resttreffer sind ausschließlich
  `Accepted` ADRs ([`ADR-0055`](../../adr/0055-nats-change-notification-wecksignal.md) 1,
  [`ADR-0087`](../../adr/0087-beispiel-clients-csharp-kotlin.md) 2,
  [`ADR-0093`](../../adr/0093-digest-korrektur-adr-0087-kotlin-basis-image.md) 6,
  [`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md) 2,
  zusammen 11 — am Parent aus `git grep -c` je ADR-Datei gemessen). `.claude/settings.json`
  (`Bash(docker manifest inspect:*)`, eine Berechtigung, kein Gewinnungsweg) und
  `tools/harness/image-stale.sh` (ein echter `docker manifest inspect`-Aufruf auf
  einen Major-Tag, kein Pin) sind aus dem Suchraum der letzten beiden Zeilen
  ausgenommen. Ein Kommentar, der nach der Hebung den neuen Gewinnungsweg nennt,
  darf die Wörter der letzten Zeile nicht tragen (sonst Soll 0 verfehlt).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `in-progress/` trägt keinen anderen Slice
(WIP-Limit 1), und der Implementer hat die Messung aus Liefer-Punkt 1 am
aktuellen Arbeitsstand gefahren — zeigt sie für keine der acht Referenzen
Drift oder Fehlform mehr, entfällt der Slice (Datei nach `done/`, §7 Zeile
`Gegenstand:`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): die Messung zeigt
  mehr als die acht Referenzen (eine weitere Referenz der Form
  `<image>@sha256:…` außerhalb von `docs/` und `.harness/`, die driftet), oder
  ein gehobener Pin verlangt Änderungen über den Pin hinaus (neue Befunde eines
  Bau-Werkzeugs, geänderte Projektdateien). Der Slice hebt dann die
  unbeteiligten Pins, der Rest wird ein eigener Slice.
- `in-progress` → `open` (blockiert — Carveout?): ein Registry- oder
  Netzausfall hindert Messung oder Digest-Gewinnung; kein Carveout, weil kein
  Gate rot ist. Gleiches gilt für die **Stopp-Regel** aus §3 (Hauptversion oder
  Tag-Linie ändert sich): der betroffene Pin bleibt auf dem Parent, die
  Entscheidung darüber ist Sache des Planners/Architects, nicht des
  Implementers.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Die DoD aus §2 ist vollständig, `make gates` endet mit Exit 0 am Stand der
Closure, die Achsen-Messung druckt für jede gehobene Achse `OK`, der reale
Post-Push-Lauf von `e2e.yml` (beide Legs) ist gelesen, und die Closure-Notiz in
§7 trägt einen Lerneintrag (geschärfte Regel, neuer Sensor oder benannte
Spec-Lücke). Ein Gate, das am Stand der Closure rot ist, geht nur mit
dokumentiertem Carveout nach `done/`. Der Slice
[`slice-pin-stale-alle-digest-pins`](slice-pin-stale-alle-digest-pins.md) startet
frühestens, wenn dieser in `done/` liegt.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

Ausgangsform je Risiko: eingetreten (CO-NNN oder Folge-Slice) · entfallen
(Grund) · weiter offen (BEO-Eintrag im Register). Alle Ausgänge sind bis zur
Closure **offen** (Platzhalter `Ausgang: offen bis Closure`).

- **Eine Hebung ändert das Verhalten eines Bau- oder Testlaufs** (neue
  .NET-/JDK-/Python-/NATS-Patchstände; `make test-notify` und der
  NATS-Vollinhalts-Rundlauf in `make test-integration` sind am wahrscheinlichsten
  betroffen, danach die drei SDK-Integrationsläufe): ein Lauf, der am Parent grün
  und am Diff rot ist, ist ein Befund, die Schwelle bleibt
  ([`AGENTS.md`](../../../../AGENTS.md) §3.6); der Pin der Achse bleibt dann auf
  dem Parent. — **Ausgang:** offen bis Closure.
- **Ein Pin war bisher ein Einzelplattform-Digest und der Bau läuft
  auf einer anderen Plattform als die Messung:** ein Index-Digest löst auf die
  Plattform des Bau-Hosts auf; ein arm64-Host baute bisher amd64-Images oder
  umgekehrt. — **Ausgang:** offen bis Closure (zu belegen: Plattform des
  Messhosts in §7 genannt, `uname -m` am Lauf-Host).
- **Eine Kopie des Pins bleibt auf dem alten Digest** (zwei Toolchains
  koexistieren, der Sensor aus dem Folge-Slice meldete sie als `DRIFT`): —
  **Ausgang:** offen bis Closure (zu belegen: der `diff`-Suchlauf aus §3 trifft
  die Erwartung, `make suchlauf-nachmessen PLAN=…` Exit 0).
- **Der Digest eines Tags bewegt sich zwischen Messung und Push** (Registry-Push
  des Upstream zwischen Liefer-Punkt 1 und dem Post-Push-Lauf): die erste
  Messung des Folge-Slice wäre dann nicht grün. — **Ausgang:** offen bis
  Closure; der Folge-Slice misst am Start neu und hebt einen frisch gedrifteten
  Pin zuerst.
- **CI-Matrix `e2e.yml` PostgreSQL 17/18:** die Datei wird berührt (Pin-Wert und
  Kommentare, keine Struktur); ein lokal grüner Lauf zeigt nicht, dass der Runner
  grün läuft ([`AGENTS.md`](../../../../AGENTS.md) §3.10). Die Pins ändern sich,
  gegen die die Matrix läuft; der Post-Push-Lauf beider Legs wird gelesen
  (`gh run list --workflow e2e.yml`, `gh run view`). — **Ausgang:** offen bis
  der Post-Push-Lauf gelesen ist; Träger der Nachverfolgung bei „weiter offen“:
  [`BEO-PGC/github-actions-unverifizierbar-lokal`](../observations/BEO-PGC/github-actions-unverifizierbar-lokal/observation.md).
- **`make image-cve` mit dem neuen `TRIVY_IMAGE`:** die README nennt das Ziel
  advisory und, bis zum ersten Release, strukturell scheiternd; seit dem Release
  `v0.2.0` ist das nicht mehr gemessen. Ein roter Lauf ist kein Gate-Rot. —
  **Ausgang:** offen bis Closure (gedruckte Zeile und Exit in §7).
- **Alarmmüdigkeit:** das Register führt
  [`BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit`](../observations/BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit/observation.md)
  mit zwei Evidenz-Dateien (gezählt am Planungsstand); ein weiterer roter
  Nachtlauf in der Zeit zwischen Hebung und Sensor-Einführung wäre eine dritte
  Datei und macht die Beobachtung zur Lücke mit eigenem Slice. —
  **Ausgang:** offen bis Closure.

## 7. Closure-Notiz

Wird mit der Closure gefüllt (Inhalt, dann `git mv`, dann Häkchen der
Paarungs-Zeile — [`AGENTS.md`](../../../../AGENTS.md) §3.3). Die Messzeilen des
Liefer-Punkts 1 (volle Digests aus den gedruckten Zeilen, Parent-Kennung,
Plattform des Messhosts), der Vergleich gegen den Parent je Lauf, die
`diff`-Suchlauf-Befunde (Gefundenes und Nichtgefundenes) und die Post-Push-Läufe
stehen hier.

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

- **Was hat funktioniert:** — (bei Closure)
- **Was ging anders als geplant:** — (bei Closure)
- **Steering-Loop-Eintrag:** — (bei Closure; Lerneintrag Pflicht, siehe §5)
- **Beobachtungs-Register (`../observations/`):** — (bei Closure; Kandidat:
  eine weitere Evidenz-Datei zu
  [`BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar`](../observations/BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar/observation.md),
  wenn die Messung zeigt, dass der Drift seit der letzten Hebung unbeobachtet
  entstand)
- **Folge-Slices:** [`slice-pin-stale-alle-digest-pins`](slice-pin-stale-alle-digest-pins.md)
  (Sensor P10) — ist eine Datei in `open/`
- **Risiken aus §6:** — (bei Closure: je Risiko genau ein Ausgang)
- **Drei Paarungen:** — (bei Closure: Anker · Folge-Slice · Register, Ergebnis)

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist allein die Default-Sub-Area
`*`/`PGC` (Build-, Test- und Gate-Konfiguration, Basis-Images der SDK- und
Beispiel-Dockerfiles; kein Domänen- oder Adapter-Code). Eine feinere Aufteilung
trägt der Slice nicht: er ändert Pin-Werte und Kommentare, keine
Verhaltens-Pfade.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`docs/plan/planning/observations/BEO-PGC/`, Namen gelesen, Evidenz-Dateien am
Planungsstand gezählt). Treffer:
[`BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar`](../observations/BEO-PGC/pin-ohne-inventar-eintrag-driftet-unsichtbar/observation.md)
(Zähler 1, eine Datei unter `evidence/`) — der Gegenstand dieses Slice; hat mit
[`ADR-0146`](../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)
einen Träger. [`BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit`](../observations/BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit/observation.md)
(Zähler 2): siehe §6. [`BEO-PGC/github-actions-unverifizierbar-lokal`](../observations/BEO-PGC/github-actions-unverifizierbar-lokal/observation.md)
(Zähler 8, hat einen eigenen Träger in `AGENTS.md` §3.10): betrifft `e2e.yml`,
§6. Kein Treffer für die übrigen Einträge.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (die Modus-Deklaration in
[`harness/conventions.md`](../../../../harness/conventions.md) führt `*` als
Greenfield).
