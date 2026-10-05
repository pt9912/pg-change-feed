# `make image`, `make image-stale`, `make image-cve` — Image bauen und prüfen

Ausführliche Fassung der Index-Zeilen aus [`harness/README.md` §Sensors](../README.md#sensors-feedback-gates); die Zeile dort trägt einen Satz und verlinkt hierher.

## `make image`

baut das OCI-Image (Multi-Stage, digest-gepinnt); **Image-Digest** nach `harness/image-hash.txt` — der Digest ist der **Lauf-Beleg** des letzten offiziellen `make image`-Laufs (builder- und lauf-gebunden, kein Inhalts-Fingerabdruck); ein Digest-Vergleich über Umgebungen/Läufe ist kein Staleness-Beweis, Inhalts-Streits entscheidet der sha256 des extrahierten Binaries (Container-Export). Die Datei ist **lokal, nicht committet** — der lokale `:dev`-Pfad (ohne `VERSION`) wird nie gepusht, also braucht keine über Commits hinweg nachschlagbare Historie zu existieren ([`ADR-0103`](../../docs/plan/adr/0103-image-hash-lokal-statt-committet.md)); ein Zug, der Build-Kontext-Dateien ändert, läuft `make image` vor seiner Closure, ohne Commit-Schritt. Mit `VERSION=<semver>` (optional `LATEST=true`, [`ADR-0051`](../../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 1/2) pusht **derselbe** Build-Aufruf stattdessen nach GHCR **und** Docker Hub (Content-Mirror, kein zweiter Build), als Multi-Arch-Manifestliste über `--platform linux/amd64,linux/arm64` (ein Buildx-Aufruf, kein zweiter Build je Plattform; die Digest-Extraktion liefert dabei den Index-Digest, übernommen aus der Prüfung gegen eine lokale Test-Registry, die der Kommentar am Ziel im `Makefile` nennt) — dieser Pfad braucht Registry-Login und läuft real nur innerhalb von `release.yml`; der Digest dieses Laufs bekommt dort einen eigenen, dauerhaften Beleg (GitHub-Release-Beschreibungstext), kein Widerspruch zu `ADR-0103`s Begründung für den bewusst einplattformigen `:dev`-Pfad (`--load` bleibt für eine konsistente, schnelle lokale Dev-Iteration einplattformig — ob `--load` überhaupt eine Multi-Platform-Manifestliste laden kann, hängt vom Storage-Treiber ab; übernommen aus demselben Kommentar im `Makefile`: mit aktiviertem containerd-Image-Store gelingt es, keine grundsätzliche buildx-Grenze)

**Bindung:** kein Gate, [`ADR-0044`](../../docs/plan/adr/0044-image-beleg-semantik.md), [`ADR-0103`](../../docs/plan/adr/0103-image-hash-lokal-statt-committet.md), [`ADR-0051`](../../docs/plan/adr/0051-cicd-pipeline-github-actions.md), [`ADR-0042`](../../docs/plan/adr/README.md)

## `make image-stale`

meldet Base-Image-Drift: FROM-Digests des Dockerfile gegen die aktuellen Registry-Digests derselben Tags, plus Existenz des nächsten Major-Tags; braucht Netz

**Bindung:** kein Gate, [`ADR-0051`](../../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 7 (Update = bewusster Digest-Commit, Modul 14)

## `make image-cve`

Trivy CRITICAL/HIGH gegen das publizierte GHCR-`:latest`-Image (digest-gepinntes `aquasec/trivy`), Exit 1 bei mindestens einem Befund dieser Schweregrade; braucht Netz, `GHCR_USERNAME`/`GHCR_PASSWORD` optional für ein privates Paket. Gemessen 2026-10-05 mit `make image-cve` (Exit 0): das Ziel-Image `ghcr.io/pt9912/pg-change-feed:latest` existiert (Release `v0.7.0`, `gh release list`), die gedruckte Zeile nennt `debian 12.15` mit `0` Befunden; automatisiert über `.github/workflows/image-scan.yml` (`schedule` nächtlich + `workflow_dispatch`, advisory, kein Gate, kein PR-Blocker)

**Bindung:** kein Gate, [`ADR-0051`](../../docs/plan/adr/0051-cicd-pipeline-github-actions.md) Entscheidung 6
