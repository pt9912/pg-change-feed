# Releasing: Release-Prozess für Betreiber und Maintainer

Version: 1.6
Stand: 2026-09-21

## 1. Zweck und Zielgruppe

Dieses Dokument beschreibt den **Mechanismus**, mit dem ein Release von
PG Change Feed entsteht: wie ein Release ausgelöst wird, was dabei
automatisch passiert, und welche begleitenden, nicht-blockierenden
Prüfungen laufen. Es richtet sich an Maintainer, die einen Release-Tag
setzen, und an Betreiber, die verstehen wollen, woher ein bestimmtes
Image auf GHCR oder Docker Hub stammt.

**Drei reale Server-Release-Tags sind bereits gesetzt** — `v0.1.0`,
`v0.1.1` und `v0.1.2` liefen jeweils mit grünem
`.github/workflows/release.yml`-Lauf durch (GHCR- und Docker-Hub-Push,
GitHub-Release mit Image-Digest). Der hier beschriebene Server-Release-
Mechanismus ist damit End-zu-Ende mit echten Repository-Secrets bewiesen,
nicht nur implementiert. Zwei der drei separaten SDK-Release-Wege (§4
„SDK-Release") sind inzwischen ebenfalls real mit einem grünen Tag-Push
bewiesen: `sdk-csharp-v0.1.0` (`PgChangeFeed.Client` auf NuGet.org) und
`sdk-python-v0.1.0` (`pgchangefeed` auf PyPI). Der dritte, jüngste Weg
(`sdk-kotlin-v*`, `pgchangefeed-kotlin` auf GitHub Packages) ist
implementiert, aber bis zum ersten realen `sdk-kotlin-v*`-Tag-Push nach
[`AGENTS.md`](../../AGENTS.md) §3.10 noch unbewiesen (siehe unten).

Dieses Dokument ersetzt nicht `docs/user/benutzerhandbuch.md` — jenes
beschreibt den laufenden Betrieb des Feed-Containers (Umgebungsvariablen,
SQL-Zugriffe), dieses hier den Release-Prozess des Projekts selbst.

## 2. Versionierung

Die aktuelle Version steht in [`docs/user/version.md`](version.md) —
eine einzelne Zeile, reines SemVer 2.0 ohne führendes `v` (z. B. `0.1.0`).
Diese Datei ist die **Quelle der Wahrheit** für die Version: der
Release-Workflow (§4) gleicht sie gegen den gesetzten Git-Tag ab und
bricht bei jeder Abweichung ab, bevor irgendein Login, Build oder Push
läuft.

Eine neue Version anzukündigen heißt: `docs/user/version.md` auf die
neue Version setzen, committen, **dann** den passenden Tag setzen (§3).
Umgekehrt — taggen ohne vorherigen `version.md`-Commit — bricht der
Release-Workflow kontrolliert ab.

## 3. Einen Release auslösen

Ein Release entsteht durch das Pushen eines Git-Tags der Form
`v<SemVer>` (z. B. `v0.1.0`, `v1.2.0-rc.1`, `v2.0.0+build.5`):

```bash
git tag v0.1.0
git push origin v0.1.0
```

Der Tag wird strikt gegen SemVer 2.0 validiert
(`tools/harness/release-tag-info.sh`, netzlos testbar über
`make test-release-tag-info`) — ein ungültiger Tag (z. B. `1.0.0` ohne
`v`, `v1.0` ohne Patch-Segment, `v1.0.0-01` mit führender Null im
Prerelease) bricht den Workflow sofort ab, bevor irgendetwas gebaut
oder gepusht wird.

**Wer taggen darf**, entscheidet die GitHub-Repository-Konfiguration
(Branch-/Tag-Protection-Regeln) — das ist keine Eigenschaft, die dieser
Workflow selbst durchsetzt.

## 4. Was beim Release automatisch passiert

Trigger: `push: tags: ['v*']` in
[`.github/workflows/release.yml`](../../.github/workflows/release.yml)
(`ci.yml`/`e2e.yml` schließen Tag-Pushes per `tags-ignore` aus — kein
Doppellauf).

1. **Tag validieren, Version/Stabilität ermitteln** — siehe §3.
2. **`docs/user/version.md` gegen den Tag abgleichen** — Abbruch bei
   Abweichung (§2).
3. **Image bauen und pushen** über `make image VERSION=<Version>` (ein
   Build, kein separater `docker build`): Push nach **GHCR und Docker
   Hub** (Content-Mirror, derselbe Bau-Vorgang, kein zweiter Build) —
   `ghcr.io/pt9912/pg-change-feed:<Version>` und
   `docker.io/pt9912/pg-change-feed:<Version>`. Bei einem **stabilen**
   Tag (kein Prerelease-Anteil) wird zusätzlich `:latest` auf beiden
   Registries gesetzt. Jedes Tag trägt **beide** Plattformen —
   `linux/amd64` **und** `linux/arm64` — in einer einzigen
   Manifestliste: läuft nativ auf Intel/AMD-Hardware, Apple-Silicon-Macs
   und ARM-basierten Linux-/Windows-Hosts (Docker Desktop, WSL2-Backend),
   ohne QEMU-Emulation zur Laufzeit. Ein natives Windows-Container-Image
   (`os: windows`) gibt es bewusst nicht — Docker Desktop unter Windows
   nutzt für Linux-Container ohnehin die WSL2-Linux-Engine.
4. **GitHub-Release anlegen** mit dem Image-Digest (`sha256:…`) im
   Beschreibungstext, als nachvollziehbarer Beleg dafür, welcher exakte
   Bau hinter dem Tag steht.
5. **Docker-Hub-Beschreibung synchronisieren** (zusätzlicher
   `hub-description`-Job, `needs: release`): spiegelt
   [`README.md`](../../README.md) als Docker-Hub-Repository-Beschreibung
   von `pt9912/pg-change-feed`. Ein Fehlschlag hier ist reine
   Präsentation — er lässt den bereits abgeschlossenen `release`-Job
   unberührt und blockiert kein Release.

### Benötigte Repository-Secrets

| Secret | Zweck | Angelegt von |
|---|---|---|
| `DOCKERHUB_USERNAME` | Docker-Hub-Login (Push + Beschreibungs-Sync) | Betreiber, manuell in GitHub |
| `DOCKERHUB_TOKEN` | Docker-Hub-Personal-Access-Token | Betreiber, manuell in GitHub |
| `NUGET_API_KEY` | NuGet.org-API-Key für `dotnet nuget push` (C#-SDK-Release, siehe unten) | Betreiber, manuell in GitHub |
| `PYPI_API_TOKEN` | PyPI-API-Token für `uv publish` (Python-SDK-Release, siehe unten) | Betreiber, manuell in GitHub |

**Scope-Hinweis:** `DOCKERHUB_TOKEN` braucht den Scope
**`read/write/delete`** — ein Token mit nur `read/write` authentifiziert
den Image-Push (Schritt 3) erfolgreich, scheitert aber am
Beschreibungs-Update (Schritt 5) mit `403 Forbidden`. Real im
Schwester-Repo d-check dokumentiert (`packaging/dockerhub/README.md`
§Transport): dort blieb der Fehler monatelang unbemerkt, weil der
betroffene Schritt mit `continue-on-error` lief. In diesem Repo ist der
`hub-description`-Job ein eigener, nicht maskierter Job — ein
Scope-Fehler bleibt dort sichtbar rot.

Alle drei Secrets sind eine externe, kontobezogene Handlung — kein
technischer Bestandteil dieses Repos legt sie an.

### SDK-Release: NuGet.org-Publish für `PgChangeFeed.Client`

Unabhängig vom oben beschriebenen Server-Image-Release existiert ein
zweiter, eigenständiger Release-Mechanismus für das C#-SDK-Package
`PgChangeFeed.Client`
([`ADR-0106`](../plan/adr/0106-csharp-nuget-erstes-sdk-package.md)
Festlegung 4): ein eigener Tag-Namensraum `sdk-csharp-v<SemVer>` (z. B.
`sdk-csharp-v0.1.0`) — bewusst getrennt vom Server-Namensraum `v*` (§3),
weil die SDK-Versionierung unabhängig vom Server läuft (`ADR-0106`
Festlegung 3, Alternative E3 verworfen) und `sdk-csharp-v*` das Muster
`v*` in `release.yml` ohnehin nicht matcht (kein
Präfix-Überlappungs-Doppellauf mit `ci.yml`/`e2e.yml`/`release.yml`).

Trigger: `push: tags: ['sdk-csharp-v*']` in
[`.github/workflows/sdk-csharp-release.yml`](../../.github/workflows/sdk-csharp-release.yml).
Der Workflow:

1. validiert den Tag-Suffix strikt gegen SemVer 2.0
   (`tools/harness/sdk-csharp-release-tag-info.sh`, netzlos testbar über
   `make test-sdk-csharp-release-tag-info`);
2. gleicht die ermittelte Version gegen die `<Version>` in
   `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj` ab —
   Abbruch bei jeder Abweichung, vor jedem Build/Push (Muster analog dem
   Tag-vs-`version.md`-Abgleich in §2/§3, hier gegen die Projektdatei
   statt gegen eine Markdown-Datei);
3. baut/testet/paketiert Docker-only über `make sdk-pack-csharp`;
4. veröffentlicht das erzeugte `.nupkg` per `dotnet nuget push … --api-key
   ${{ secrets.NUGET_API_KEY }} --source
   https://api.nuget.org/v3/index.json`.

Kein `:latest`-Äquivalent (NuGet kennt keins) und kein
GitHub-Release-Eintrag für das SDK — beides bewusst außerhalb dieser
Folgepflicht (`ADR-0106`).

**Ein realer `sdk-csharp-v*`-Tag ist bereits gesetzt** — `sdk-csharp-v0.1.0`
lief mit grünem `sdk-csharp-release.yml`-Lauf durch und veröffentlichte
`PgChangeFeed.Client` 0.1.0 real auf NuGet.org (`dotnet nuget push`
bestätigte `201 Created`/„Your package was pushed"). Der C#-SDK-Release-Weg
ist damit End-zu-Ende mit echtem `NUGET_API_KEY`-Secret bewiesen, nicht nur
implementiert.

### SDK-Release: PyPI-Publish für `pgchangefeed`

Analog zum C#-SDK-Release-Weg oben existiert ein dritter, eigenständiger
Release-Mechanismus für das Python-SDK-Package `pgchangefeed`
([`ADR-0107`](../plan/adr/0107-python-pypi-zweites-sdk-package.md)
Festlegung 5, [`ADR-0108`](../plan/adr/0108-python-sdk-uv-statt-build-twine.md)
§Entscheidung Festlegung 1): ein eigener Tag-Namensraum
`sdk-python-v<PEP 440>` (im hier genutzten einfachen Fall
`MAJOR.MINOR.PATCH`, z. B. `sdk-python-v0.1.0`) — getrennt sowohl vom
Server-Namensraum `v*` (§3) als auch vom C#-SDK-Namensraum
`sdk-csharp-v*` (siehe oben), weil die Python-SDK-Versionierung
unabhängig von beiden läuft (`ADR-0107` Festlegung 4) und
`sdk-python-v*` keines der drei anderen Muster (`v*`, `sdk-csharp-v*`,
sowie `ci.yml`/`e2e.yml`s `tags-ignore: ['**']`) matcht (kein
Präfix-Überlappungs-Doppellauf).

Trigger: `push: tags: ['sdk-python-v*']` in
[`.github/workflows/sdk-python-release.yml`](../../.github/workflows/sdk-python-release.yml).
Der Workflow:

1. validiert den Tag-Suffix strikt gegen den einfachen PEP-440-Fall
   `MAJOR.MINOR.PATCH` (`tools/harness/sdk-python-release-tag-info.sh`,
   netzlos testbar über `make test-sdk-python-release-tag-info` — eine
   eigenständige Regex, **kein** Sourcing aus
   `tools/harness/semver-regex.sh`, weil PEP 440 nicht identisch mit
   SemVer 2.0 ist, auch wenn der hier genutzte einfache Fall zu beiden
   Grammatiken kompatibel bleibt);
2. gleicht die ermittelte Version gegen die `[project] version` in
   `sdks/python/pgchangefeed/pyproject.toml` ab — Abbruch bei jeder
   Abweichung, vor jedem Build/Push (Muster analog dem
   Tag-vs-`version.md`-Abgleich in §2/§3 bzw. dem
   Tag-vs-`.csproj`-Abgleich oben, hier gegen die `pyproject.toml`);
3. baut/testet/paketiert Docker-only über `make sdk-pack-python`
   (erzeugt Wheel `.whl` und Source-Distribution `.tar.gz`);
4. veröffentlicht beide Artefakte per `uv publish` — Token über die
   Umgebungsvariable `UV_PUBLISH_TOKEN: ${{ secrets.PYPI_API_TOKEN }}`
   (kein `-u __token__ -p`-Flag-Paar wie bei `twine`). `uv` ist auf dem
   GitHub-hosted Runner nicht vorinstalliert und wird über die
   SHA-gepinnte Action `astral-sh/setup-uv` bezogen, mit demselben
   `uv`-Versionsstand (`0.12.17`), den `sdks/python/Dockerfile` für den
   Docker-only-Bau-/Pack-Schritt per digest-gepinntem `COPY --from=`
   nutzt.

Kein `:latest`-Äquivalent (PyPI kennt keins) und kein
GitHub-Release-Eintrag für das SDK — beides bewusst außerhalb dieser
Folgepflicht (`ADR-0107`).

**Ein realer `sdk-python-v*`-Tag ist bereits gesetzt** —
`sdk-python-v0.1.0` lief mit grünem `sdk-python-release.yml`-Lauf durch
und veröffentlichte `pgchangefeed` 0.1.0 real auf PyPI (die PyPI-API
`https://pypi.org/pypi/pgchangefeed/json` bestätigt Version `0.1.0`
sofort sichtbar, ohne Indexierungsverzögerung wie bei NuGet). Der
Python-SDK-Release-Weg ist damit End-zu-Ende mit echtem
`PYPI_API_TOKEN`-Secret bewiesen, nicht nur implementiert.

### SDK-Release: GitHub-Packages-Publish für `pgchangefeed-kotlin`

Analog zu den beiden SDK-Release-Wegen oben existiert ein vierter,
eigenständiger Release-Mechanismus für das Kotlin-SDK-Package
`pgchangefeed-kotlin`
([`ADR-0109`](../plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md)
Festlegung 2/5): ein eigener Tag-Namensraum `sdk-kotlin-v<SemVer 2.0>`
(z. B. `sdk-kotlin-v0.1.0`) — getrennt vom Server-Namensraum `v*` (§3)
sowie vom C#-SDK-Namensraum `sdk-csharp-v*` und vom Python-SDK-Namensraum
`sdk-python-v*` (siehe oben), weil die Kotlin-SDK-Versionierung unabhängig
von allen dreien läuft (`ADR-0109` Festlegung 4) und `sdk-kotlin-v*` keines
der vier anderen Muster (`v*`, `sdk-csharp-v*`, `sdk-python-v*`, sowie
`ci.yml`/`e2e.yml`/`examples.yml`s `tags-ignore: ['**']`) matcht (kein
Präfix-Überlappungs-Doppellauf).

Trigger: `push: tags: ['sdk-kotlin-v*']` in
[`.github/workflows/sdk-kotlin-release.yml`](../../.github/workflows/sdk-kotlin-release.yml).
Der Workflow:

1. validiert den Tag-Suffix strikt gegen SemVer 2.0
   (`tools/harness/sdk-kotlin-release-tag-info.sh`, netzlos testbar über
   `make test-sdk-kotlin-release-tag-info` — teilt die SemVer-2.0-Regex mit
   `tools/harness/release-tag-info.sh`/`tools/harness/sdk-csharp-release-tag-info.sh`
   über `tools/harness/semver-regex.sh`, dieselbe Grammatik wie beim
   C#-SDK, **nicht** Pythons PEP-440-Sonderfall);
2. gleicht die ermittelte Version gegen die Top-Level-`version`-Zeile in
   `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` ab — Abbruch bei
   jeder Abweichung, vor jedem Build/Publish (Muster analog dem
   Tag-vs-`.csproj`-Abgleich bzw. dem Tag-vs-`pyproject.toml`-Abgleich
   oben, hier gegen die Gradle-Versionszeile);
3. baut/testet/paketiert Docker-only über `make sdk-pack-kotlin` — ein
   roter Test bricht den Workflow ab, bevor der Publish-Schritt erreicht
   wird;
4. veröffentlicht anschließend ebenfalls Docker-only, über eine eigene
   `publish`-Docker-Stufe (`sdks/kotlin/Dockerfile`, baut auf der bereits
   vorhandenen `build`-Stufe auf) — `docker build --build-context
   proto=proto --target publish` gefolgt von `docker run --rm -e
   GITHUB_ACTOR=… -e GITHUB_TOKEN=… <image> ./gradlew --no-daemon publish`,
   mit Netzwerkzugriff zur Laufzeit (kein `--network none`, anders als
   jede andere Stufe dieses Dockerfiles, `ADR-0109` Festlegung 5 wörtlich:
   Docker-only bis einschließlich `publish`). Die `.proto`-Quelle
   (`proto/cdc/stream/v1/changestream.proto`) fließt über denselben
   benannten Bau-Kontext `proto` ein wie bei `make sdk-pack-kotlin` — kein
   manueller Kopier-Schritt außerhalb von Docker.

**Kein externes Repository-Secret nötig** — der zentrale Unterschied zu
den beiden Release-Wegen oben: GitHub Packages authentifiziert
ausschließlich über das eingebaute `GITHUB_TOKEN`
(`permissions: contents: read` / `packages: write` auf Job-Ebene, real
dokumentiertes Minimalrezept, `ADR-0109` §Kontext Recherche). `GITHUB_ACTOR`
ist ein von GitHub Actions automatisch bereitgestellter
Default-Umgebungswert; `GITHUB_TOKEN` wird explizit aus
`secrets.GITHUB_TOKEN` als Umgebungsvariable gesetzt, gelesen von
`sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts`s
`publishing.repositories.maven.credentials`. Kein Betreiber-Schritt, kein
neuer Eintrag in der Secret-Tabelle oben.

**Ein realer Preis bleibt trotzdem bestehen — anders als bei NuGet und
PyPI**: GitHub Packages verlangt eine Authentifizierung zum **Lesen**, auch
für ein öffentliches Package (real dokumentiert, `ADR-0109`
§Entscheidung Festlegung 2). Ein Kotlin-Consumer, der `pgchangefeed-kotlin`
in sein eigenes Projekt einbindet, braucht ein GitHub-Konto und einen
klassischen Personal-Access-Token mit `read:packages`-Scope, den er in
seiner eigenen Build-Konfiguration hinterlegt — ein anonymer Bezug wie bei
Maven Central, NuGet oder PyPI ist bei GitHub Packages strukturell nicht
möglich.

Kein `:latest`-Äquivalent (GitHub Packages kennt keins) und kein
GitHub-Release-Eintrag für das SDK — beides bewusst außerhalb dieser
Folgepflicht (`ADR-0109`).

**Der reale, grüne Post-Push-Lauf steht noch aus** — dieser Release-Weg
ist implementiert und `make gates` läuft grün, aber nach
[`AGENTS.md`](../../AGENTS.md) §3.10 bleibt er bis zum ersten echten
`sdk-kotlin-v*`-Tag-Push unbewiesen: `./gradlew publish`-Verhalten im
`publish`-Docker-Image und das reale GitHub-Packages-Registry-Antwortverhalten
sind lokal strukturell nicht prüfbar (kein Docker-only-Sensor kann einen
echten Registry-Schreibzugriff gegen `maven.pkg.github.com` ersetzen).

## 5. Begleitende, nicht-blockierende Workflows

Zwei weitere Workflows laufen unabhängig vom Release-Trigger, nächtlich
und per `workflow_dispatch`. Beide sind **advisory**: ein roter Lauf ist
in der Actions-Übersicht sichtbar, blockiert aber kein Release und kein
anderes Gate.

- **[`image-scan.yml`](../../.github/workflows/image-scan.yml)**
  (`make image-cve`): Trivy CRITICAL/HIGH gegen das publizierte GHCR-
  `:latest`-Image. Scheitert strukturell, solange kein stabiler Release
  je gelaufen ist (kein `:latest`-Tag existiert dann).
- **[`upstream-drift.yml`](../../.github/workflows/upstream-drift.yml)**:
  fragt das vollständige Neun-Achsen-Pin-Inventar ab (Basis-Image-Digests,
  Test-/Werkzeug-Container-Pins, die adoptierte Kurs-Baseline-Version,
  alle SHA-gepinnten GitHub-Action-`uses:`-Zeilen) — meldet gefundenen
  Drift, hebt ihn nicht automatisch an. Eine Pin-Anhebung bleibt ein
  bewusster, separater Commit.

## 6. Rollback

Es gibt keinen automatisierten Rollback-Mechanismus. Ein fehlerhafter
Release wird durch einen neuen, höheren Tag mit einer korrigierten
Version ersetzt — bereits gepushte Image-Tags und GitHub-Releases werden
nicht rückwirkend verändert oder gelöscht.

### Änderungshistorie

| Version | Datum | Änderung |
|---|---|---|
| 1.0 | 2026-09-19 | Erste Fassung — dokumentiert den in `welle-release-pipeline-adr-0051` real implementierten Release-Prozess (`ADR-0051`) |
| 1.1 | 2026-09-19 | §4 um den unabhängigen SDK-Release-Weg (`sdk-csharp-v*`-Tag, `NUGET_API_KEY`) ergänzt — Fixrunde nach Review-Finding F-1 (`docs/reviews/review-slice-sdk-csharp-publish-workflow.md`, `LH-FA-SST-009`, `ADR-0106`) |
| 1.2 | 2026-09-19 | §1 korrigiert: drei reale Server-Release-Tags (`v0.1.0`–`v0.1.2`) sind bereits gesetzt und liefen grün durch — der Server-Release-Mechanismus ist End-zu-Ende bewiesen; der SDK-Release-Weg bleibt bis zum ersten realen `sdk-csharp-v*`-Tag-Push separat unbewiesen |
| 1.3 | 2026-09-19 | §4 um den dritten, unabhängigen SDK-Release-Weg (`sdk-python-v*`-Tag, `PYPI_API_TOKEN`, `uv publish`) ergänzt, Secret-Tabelle um `PYPI_API_TOKEN` erweitert — vorab eingeplanter DoD-Punkt von `slice-sdk-python-publish-workflow` (`LH-FA-SST-009`, `ADR-0107`, `ADR-0108`), nicht erst nach einem Reviewer-Finding (Lehre aus `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen`) |
| 1.4 | 2026-09-20 | §1/§4 korrigiert: `sdk-python-v0.1.0` real gesetzt, `pgchangefeed` 0.1.0 real auf PyPI veröffentlicht — der Python-SDK-Release-Weg ist damit wie der C#-Weg End-zu-Ende bewiesen |
| 1.5 | 2026-09-20 | §1/§4 um den vierten, unabhängigen SDK-Release-Weg (`sdk-kotlin-v*`-Tag, GitHub Packages, `GITHUB_TOKEN`, kein externes Secret) ergänzt — vorab eingeplanter DoD-Punkt von `slice-sdk-kotlin-publish-workflow` (`LH-FA-SST-009`, `ADR-0109`), nicht erst nach einem Reviewer-Finding (Lehre aus `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen`); der reale Post-Push-Lauf bleibt nach `AGENTS.md` §3.10 bis zum ersten echten Tag-Push offen |
| 1.6 | 2026-09-21 | §4 korrigiert: der Publish-Schritt (`./gradlew publish`) läuft jetzt Docker-only in einer eigenen `publish`-Stufe (`sdks/kotlin/Dockerfile`), nicht mehr direkt auf dem Runner — Fixrunde nach Review-Finding F-1 (`docs/reviews/review-slice-sdk-kotlin-publish-workflow.md`, `LH-FA-SST-009`, `ADR-0109` Festlegung 5) |
