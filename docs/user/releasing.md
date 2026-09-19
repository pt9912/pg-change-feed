# Releasing: Release-Prozess für Betreiber und Maintainer

Version: 1.0
Stand: 2026-09-19

## 1. Zweck und Zielgruppe

Dieses Dokument beschreibt den **Mechanismus**, mit dem ein Release von
PG Change Feed entsteht: wie ein Release ausgelöst wird, was dabei
automatisch passiert, und welche begleitenden, nicht-blockierenden
Prüfungen laufen. Es richtet sich an Maintainer, die einen Release-Tag
setzen, und an Betreiber, die verstehen wollen, woher ein bestimmtes
Image auf GHCR oder Docker Hub stammt.

**Zum Zeitpunkt dieses Dokuments wurde noch kein realer Release-Tag
gesetzt** — der hier beschriebene Prozess ist implementiert und real
gegen die beteiligten APIs/Endpunkte geprüft (siehe die einzelnen
Abschnitte), aber der End-zu-Ende-Ablauf mit echten Repository-Secrets
ist strukturell erst nach dem ersten echten Tag-Push bewiesen
(`AGENTS.md` §3.10). Dieses Dokument behauptet keinen abgeschlossenen
ersten Release.

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
   Registries gesetzt.
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

**Scope-Hinweis:** `DOCKERHUB_TOKEN` braucht den Scope
**`read/write/delete`** — ein Token mit nur `read/write` authentifiziert
den Image-Push (Schritt 3) erfolgreich, scheitert aber am
Beschreibungs-Update (Schritt 5) mit `403 Forbidden`. Real im
Schwester-Repo d-check dokumentiert (`packaging/dockerhub/README.md`
§Transport): dort blieb der Fehler monatelang unbemerkt, weil der
betroffene Schritt mit `continue-on-error` lief. In diesem Repo ist der
`hub-description`-Job ein eigener, nicht maskierter Job — ein
Scope-Fehler bleibt dort sichtbar rot.

Beide Secrets sind eine externe, kontobezogene Handlung — kein
technischer Bestandteil dieses Repos legt sie an.

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
