# Releasing: Release-Prozess für Betreiber und Maintainer

Version: 1.10
Stand: 2026-09-25

## 1. Zweck und Zielgruppe

Dieses Dokument beschreibt den **Mechanismus**, mit dem ein Release von
PG Change Feed entsteht: wie ein Release ausgelöst wird, was dabei
automatisch passiert, und welche begleitenden, nicht-blockierenden
Prüfungen laufen. Es richtet sich an Maintainer, die einen Release-Tag
setzen, und an Betreiber, die verstehen wollen, woher ein bestimmtes
Image auf GHCR oder Docker Hub stammt.

**Vier reale Server-Release-Tags sind bereits gesetzt** — `v0.1.0`,
`v0.1.1`, `v0.1.2` und `v0.2.0` liefen jeweils mit grünem
`.github/workflows/release.yml`-Lauf durch (GHCR- und Docker-Hub-Push,
GitHub-Release mit Image-Digest). Der Lauf von `v0.2.0` (Lauf-Kennung
36190768475, 2026-09-25) endete mit `success` in beiden Jobs (Release und
`hub-description`); `ghcr.io/pt9912/pg-change-feed` trägt die Tags `0.2.0` und
`latest`, Docker Hub trägt `0.2.0` und `latest` mit `linux/amd64` und
`linux/arm64`, und das GitHub-Release nennt den Digest
`sha256:fab4a53d96ac414578739319307dc80b7aeacbb6ba88659b88fbf48bb8c5b58e`. Die
Release-Beschreibung eines Tags lässt sich nach dem Lauf mit
`gh release edit <Tag> --notes-file <Datei>` um Release-Hinweise ergänzen; der
Workflow schreibt nur den Digest und die Image-Namen. Der hier beschriebene
Server-Release-Mechanismus ist damit End-zu-Ende mit echten
Repository-Secrets bewiesen, nicht nur implementiert. Alle drei separaten SDK-Release-Wege (§4
„SDK-Release") sind ebenfalls real mit einem grünen Tag-Push bewiesen:
`sdk-csharp-v0.1.0` und `sdk-csharp-v0.2.0` (`PgChangeFeed.Client` auf
NuGet.org), `sdk-python-v0.1.0` und `sdk-python-v0.2.0` (`pgchangefeed` auf
PyPI) sowie `sdk-kotlin-v0.2.0` (`pgchangefeed-kotlin` auf GitHub
Packages). Die `0.2.0`-Läufe der drei Workflows (Lauf-Kennungen 36117929191,
36117929298 und 36117929552, abgefragt mit `gh run list --workflow
<datei>.yml`, 2026-09-25) endeten je mit `success`. Der Kotlin-Weg
veröffentlicht zusätzlich nach Cloudsmith (zwei Jobs, ein Ziel je Job);
diese Struktur ist bis zu ihrem ersten Tag-Lauf nicht bewiesen (§4
„SDK-Release: GitHub-Packages- und Cloudsmith-Publish“).

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
| `CLOUDSMITH_USERNAME` | Anmeldename des Cloudsmith-Service-Kontos für den Upload (Kotlin-SDK-Release, Job `sdk-kotlin-cloudsmith`, siehe unten) | Betreiber, manuell in GitHub |
| `CLOUDSMITH_API_KEY` | API-Key desselben Service-Kontos (Schreibrecht auf das Repository `pt9912/pg-change-feed`) | Betreiber, manuell in GitHub |

**Scope-Hinweis:** `DOCKERHUB_TOKEN` braucht den Scope
**`read/write/delete`** — ein Token mit nur `read/write` authentifiziert
den Image-Push (Schritt 3) erfolgreich, scheitert aber am
Beschreibungs-Update (Schritt 5) mit `403 Forbidden`. Real im
Schwester-Repo d-check dokumentiert (`packaging/dockerhub/README.md`
§Transport): dort blieb der Fehler monatelang unbemerkt, weil der
betroffene Schritt mit `continue-on-error` lief. In diesem Repo ist der
`hub-description`-Job ein eigener, nicht maskierter Job — ein
Scope-Fehler bleibt dort sichtbar rot.

Alle Secrets dieser Tabelle sind eine externe, kontobezogene Handlung — kein
technischer Bestandteil dieses Repos legt sie an. Der Werttyp von
`CLOUDSMITH_USERNAME` ist der Service-Slug des Cloudsmith-Kontos; ob
Cloudsmith beim Upload stattdessen den Service-Namen verlangt, zeigt erst der
erste Tag-Lauf (bei einem Authentifizierungsfehler HTTP 401/403 im
Upload-Schritt setzt der Betreiber das Secret auf den Service-Namen und
wiederholt nur den roten Job).

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

**Zwei reale `sdk-csharp-v*`-Tags sind gesetzt** — `sdk-csharp-v0.1.0`
lief mit grünem `sdk-csharp-release.yml`-Lauf durch und veröffentlichte
`PgChangeFeed.Client` 0.1.0 real auf NuGet.org (`dotnet nuget push`
bestätigte `201 Created`/„Your package was pushed"); `sdk-csharp-v0.2.0`
lief ebenfalls grün (Lauf 36117929191). Der C#-SDK-Release-Weg ist damit
End-zu-Ende mit echtem `NUGET_API_KEY`-Secret bewiesen, nicht nur
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

**Zwei reale `sdk-python-v*`-Tags sind gesetzt** —
`sdk-python-v0.1.0` lief mit grünem `sdk-python-release.yml`-Lauf durch
und veröffentlichte `pgchangefeed` 0.1.0 real auf PyPI (die PyPI-API
`https://pypi.org/pypi/pgchangefeed/json` bestätigt Version `0.1.0`
sofort sichtbar, ohne Indexierungsverzögerung wie bei NuGet);
`sdk-python-v0.2.0` lief ebenfalls grün (Lauf 36117929298). Der
Python-SDK-Release-Weg ist damit End-zu-Ende mit echtem
`PYPI_API_TOKEN`-Secret bewiesen, nicht nur implementiert.

### SDK-Release: GitHub-Packages- und Cloudsmith-Publish für `pgchangefeed-kotlin`

Analog zu den beiden SDK-Release-Wegen oben existiert ein vierter,
eigenständiger Release-Mechanismus für das Kotlin-SDK-Package
`pgchangefeed-kotlin`
([`ADR-0109`](../plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md)
Festlegung 2/5, erweitert um Cloudsmith als zweites Ziel durch
[`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md)): ein eigener Tag-Namensraum `sdk-kotlin-v<SemVer 2.0>`
(z. B. `sdk-kotlin-v0.2.0`) — getrennt vom Server-Namensraum `v*` (§3)
sowie vom C#-SDK-Namensraum `sdk-csharp-v*` und vom Python-SDK-Namensraum
`sdk-python-v*` (siehe oben), weil die Kotlin-SDK-Versionierung unabhängig
von allen dreien läuft (`ADR-0109` Festlegung 4) und `sdk-kotlin-v*` keines
der vier anderen Muster (`v*`, `sdk-csharp-v*`, `sdk-python-v*`, sowie
`ci.yml`/`e2e.yml`/`examples.yml`s `tags-ignore: ['**']`) matcht (kein
Präfix-Überlappungs-Doppellauf).

Trigger: `push: tags: ['sdk-kotlin-v*']` in
[`.github/workflows/sdk-kotlin-release.yml`](../../.github/workflows/sdk-kotlin-release.yml).
Ein Tag löst **zwei unabhängige Jobs** aus, einen je Ziel
(`sdk-kotlin-github-packages` und `sdk-kotlin-cloudsmith`, ohne `needs` und
ohne `continue-on-error`). Scheitert ein Ziel, bleibt das andere
veröffentlicht, der Lauf ist rot, und „Re-run failed jobs“ wiederholt nur
den roten Job. Beide Ziele erhalten dieselben Artefakte. Jeder Job
durchläuft dieselben Schritte:

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
   wird; im selben Bau prüft eine Probe ohne Zugangsdaten und ohne Zugriff
   auf ein Ziel die Publish-Konfiguration (beide Einzel-Aufgaben unter ihrem
   Namen, POM-Koordinate und POM-Version gegen den Jar-Namen, Upload-URL
   des Cloudsmith-Repositories) und bricht bei einer Abweichung den Bau ab;
4. veröffentlicht anschließend ebenfalls Docker-only, über eine eigene
   `publish`-Docker-Stufe (`sdks/kotlin/Dockerfile`, baut auf der bereits
   vorhandenen `build`-Stufe auf) — `docker build --build-context
   proto=proto --target publish` gefolgt von `docker run --rm -e NAME …
   <image> ./gradlew --no-daemon <Aufgabe des Ziels>`, mit Netzwerkzugriff
   zur Laufzeit (kein `--network none`, anders als jede andere Stufe dieses
   Dockerfiles, `ADR-0109` Festlegung 5 wörtlich: Docker-only bis
   einschließlich `publish`). Die Aufgabe des Ziels ist
   `publishMavenPublicationToGitHubPackagesRepository` bzw.
   `publishMavenPublicationToCloudsmithRepository`, nicht die Sammel-Aufgabe
   `publish` (sie führt beide Ziele aus und scheitert ohne die
   Cloudsmith-Werte). Die `.proto`-Quelle
   (`proto/cdc/stream/v1/changestream.proto`) fließt über denselben
   benannten Bau-Kontext `proto` ein wie bei `make sdk-pack-kotlin` — kein
   manueller Kopier-Schritt außerhalb von Docker.

**Secrets je Ziel.** Der GitHub-Packages-Job authentifiziert ausschließlich
über das eingebaute `GITHUB_TOKEN` (`permissions: contents: read` /
`packages: write` auf Job-Ebene, real dokumentiertes Minimalrezept,
`ADR-0109` §Kontext Recherche); `GITHUB_ACTOR` ist ein von GitHub Actions
automatisch bereitgestellter Default-Umgebungswert, `GITHUB_TOKEN` wird aus
`secrets.GITHUB_TOKEN` als Umgebungsvariable des Schritts gesetzt. Der
Cloudsmith-Job (`permissions: contents: read`, kein `packages: write`) liest
ausschließlich `CLOUDSMITH_USERNAME` und `CLOUDSMITH_API_KEY`
(Secret-Tabelle oben; Betreiber-Schritt, `ADR-0123` Festlegung 3). Beide
Werte fließen über `env:` des Schritts und `docker run -e NAME` in den
Container — nie als Build-Argument und nie als Wert in der Kommandozeile;
Gradle liest sie in `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts`
(`publishing.repositories.maven.credentials`).

**Betreiber-Voraussetzungen für Cloudsmith** (extern, kein Agenten-Schritt):
ein Cloudsmith-Konto mit der Organisation `pt9912`, ein Open-Source-Repository
`pg-change-feed` (Broadcast, Sichtbarkeit „Open source“, in der Web-App
angelegt — Open-Source-Repositories lassen sich nicht über die API anlegen)
und ein Service-Konto mit Schreibrecht auf dieses Repository, dessen Name und
API-Key die zwei Repository-Secrets tragen. Die Einstellung „Republishing“
bleibt aus: eine veröffentlichte Version ist unveränderlich. Der Upload
sammelt die Einzeldateien einer Version in einem Sammelfenster (60 Sekunden)
und verarbeitet sie danach asynchron; der Job meldet Erfolg mit dem Upload,
nicht mit der Verfügbarkeit.

**Was Anwender sehen.** Der Bezug über Cloudsmith
(`https://dl.cloudsmith.io/public/pt9912/pg-change-feed/maven/`) braucht kein
Konto und keinen Token; GitHub Packages verlangt dagegen weiter eine
Authentifizierung zum **Lesen**, auch für ein öffentliches Package (real
dokumentiert, `ADR-0109` §Entscheidung Festlegung 2): ein GitHub-Konto und
einen klassischen Personal-Access-Token mit `read:packages`-Scope. Die
Kotlin-README nennt Cloudsmith deshalb als ersten Bezugsweg und trägt die von
Cloudsmith für Open-Source-Repositories verlangte Namensnennung.

**Wiederholung.** Ob ein zweiter Upload derselben Version am Ziel abgelehnt
wird, ist für Cloudsmith (nativer Maven-Upload) und GitHub Packages nicht
geprüft; der Ablauf verlässt sich nicht darauf, sondern wiederholt nur den
roten Job. Bleibt nach einem Fehlschlag ein Teil-Upload zurück (Einzeldateien
ohne vollständiges Paket), löscht der Betreiber das Paket in der Web-App des
Ziels und wiederholt den Job.

**Offen bis zum ersten Tag-Lauf.** Die Job-Struktur mit zwei Zielen ist
lokal nur statisch geprüft (`AGENTS.md` §3.10): der reale Lauf mit dem ersten
`sdk-kotlin-v*`-Tag, der beide Ziele trägt (`sdk-kotlin-v0.2.2`), muss beide Jobs `success`
zeigen, und ein anonymer Abruf der POM-Datei am Cloudsmith-Download-Pfad
(`https://dl.cloudsmith.io/public/pt9912/pg-change-feed/maven/io/github/pt9912/pgchangefeed-kotlin/0.2.2/pgchangefeed-kotlin-0.2.2.pom`)
muss HTTP 200 antworten (nach dem Sammelfenster, ggf. mit Wiederholung).
Ebenfalls bis dahin nicht geprüft: ob Cloudsmith die von Gradle
mitveröffentlichte Moduldatei annimmt, ob der Anmeldename ein Service-Slug
oder ein Service-Name sein muss und ob die Paketseite die POM-Beschreibung
anzeigt (optional pflegt der Betreiber dort den Text in der Web-App). Die
Versionen `0.2.0` und `0.2.1` liegen nur auf GitHub Packages, auf Cloudsmith
liegen erst die Versionen ab `0.2.2`.

Kein `:latest`-Äquivalent (weder GitHub Packages noch Cloudsmith kennen
eines) und kein GitHub-Release-Eintrag für das SDK — beides bewusst außerhalb
dieser Folgepflicht (`ADR-0109`, `ADR-0123`).

**Ein realer `sdk-kotlin-v*`-Tag ist gesetzt** — `sdk-kotlin-v0.2.0` lief
mit grünem `sdk-kotlin-release.yml`-Lauf durch (Lauf 36117929552; alle
Schritte `success`, darunter „Nach GitHub Packages veroeffentlichen
(./gradlew publish, im Docker-Image)"), und die Paketversion ist über die
GitHub-API abfragbar (`gh api
/users/pt9912/packages/maven/io.github.pt9912.pgchangefeed-kotlin/versions`
nennt `0.2.0`, 2026-09-25). Der Kotlin-SDK-Release-Weg ist damit
End-zu-Ende bewiesen, nicht nur implementiert.

**Beschreibung der Kotlin-Paketseite (manuell).** GitHub Packages zeigt bei
Maven-Paketen weder eine README noch die `<description>` der POM an; die
Paketseite trägt ein Beschreibungsfeld („Write a description"), das nur in der
Weboberfläche gepflegt wird — die REST- und die GraphQL-Schnittstelle bieten
dafür keine schreibende Operation (geprüft mit `gh api`: `PackageVersion.readme`
ist nur lesbar, die einzige Paket-Mutation ist `deletePackageVersion`). Nach
jedem Kotlin-Release wird der Inhalt von
[`sdks/kotlin/pgchangefeed-kotlin/README.md`](../../sdks/kotlin/pgchangefeed-kotlin/README.md)
in dieses Feld eingefügt. Ob ein eingetragener Text bei der nächsten
Veröffentlichung erhalten bleibt, ist nicht geprüft.

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
| 1.7 | 2026-09-25 | §1/§4 auf den Ist-Stand gezogen (`LH-FA-SST-009`, `ADR-0110`, slice-sdk-readme-nutzerdoku Fixrunde): alle drei SDK-Release-Wege sind mit realen Tags bewiesen (`sdk-csharp-v0.1.0`/`0.2.0`, `sdk-python-v0.1.0`/`0.2.0`, `sdk-kotlin-v0.2.0`; `0.2.0`-Läufe 36117929191, 36117929298, 36117929552 je `success`, Kotlin-Paketversion `0.2.0` über die GitHub-API abgefragt) |
| 1.8 | 2026-09-25 | §4 um den manuell gepflegten Beschreibungstext der Kotlin-Paketseite ergänzt (`LH-FA-SST-009`, `ADR-0109`): GitHub Packages zeigt bei Maven weder README noch POM-Beschreibung, die Schreib-Schnittstelle fehlt in REST und GraphQL |
| 1.9 | 2026-09-25 | §1 auf den realen Server-Release `v0.2.0` gezogen (`ADR-0051`): Lauf 36190768475 `success`, GHCR und Docker Hub tragen `0.2.0` und `latest` (amd64, arm64), Release-Hinweise lassen sich mit `gh release edit` ergänzen |
| 1.10 | 2026-09-25 | §4 Kotlin-Abschnitt auf zwei Vertriebsziele gezogen (`LH-FA-SST-009`, `ADR-0123`, slice-sdk-kotlin-cloudsmith): ein Job je Ziel (GitHub Packages, Cloudsmith) mit den Einzel-Aufgaben statt der Sammel-Aufgabe, Secret-Tabelle um `CLOUDSMITH_USERNAME` und `CLOUDSMITH_API_KEY` erweitert, Betreiber-Voraussetzungen, Wiederholung je Job und die Offen-Punkte bis zum ersten Tag-Lauf beschrieben; „Alle drei Secrets“ leitet die Zahl nicht mehr aus einem Zählwort ab |
