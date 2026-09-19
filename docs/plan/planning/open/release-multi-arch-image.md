# Slice release-multi-arch-image: Multi-Arch-Image (linux/amd64 + linux/arm64)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** ohne Welle — die Closure-Bedingung ist identisch mit der DoD
dieses einen Slice, kein *Mehr* gegenüber ihr.

**Bezug:** [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
(additive Erweiterung des dort entschiedenen „ein Build,
Content-Mirror"-Mechanismus — buildx baut mehrere Plattformen weiterhin
in einem Aufruf, kein zweiter Build, kein neuer Entscheidungspunkt).

**Berührte Spec-Stellen:** — (Prozess-/Infrastruktur-Änderung ohne
Spec-Stratum-Bezug).

**Verantwortlich:** — (direkt beauftragt, noch nicht priorisiert).

**Autor:** Planner-Agent, direkt beauftragt. **Datum:** 2026-09-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das veröffentlichte Image (`ghcr.io/pt9912/pg-change-feed`,
`docker.io/pt9912/pg-change-feed`) trägt künftig **zwei** native
Plattformen — `linux/amd64` **und** `linux/arm64` — statt nur `amd64`,
damit es auf Apple-Silicon-Macs, ARM-basierten Linux-Hosts und
ARM-basierten Windows-Geräten (Docker Desktop, WSL2-Backend) ohne
QEMU-Emulation läuft. Windows und macOS auf `amd64`-Hardware sind bereits
heute abgedeckt (Docker Desktop nutzt dort ohnehin eine Linux-VM/WSL2 —
kein natives Windows-Container-Image nötig, explizit mit dem
Auftraggeber geklärt).

Real vorab geprüft, bevor dieser Plan geschrieben wurde:

- Beide bereits gepinnten Basis-Images sind schon heute
  Multi-Arch-Manifest-Listen mit `linux/arm64/v8`-Eintrag —
  `golang:1.27-alpine@sha256:cf6fca66…` (`docker buildx imagetools
  inspect`, Plattform-Liste enthält `linux/arm64/v8`) und
  `gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872…`
  (ebenso). **Kein Re-Pin nötig**, nur Build-seitige Ziel-Plattform-Wahl.
- `Dockerfile`s `build`-Stufe ist bereits `CGO_ENABLED=0` — reines
  Go-Cross-Compiling für `GOARCH=arm64` ist ohne C-Toolchain pro
  Zielarchitektur möglich (Standard-Muster: `deps`/`proto`/
  `proto-export`/`coverage` mit `--platform=$BUILDPLATFORM` nativ auf
  dem Bau-Host bauen, nur die `build`-Stufe cross-kompiliert über
  `ARG TARGETOS TARGETARCH` — vermeidet QEMU/Binfmt in `release.yml`s
  GitHub-hosted `amd64`-Runner vollständig).
- `docker buildx build --load` (der lokale, `VERSION`-lose `:dev`-Pfad)
  kann **keine** Multi-Platform-Manifestliste laden — das ist eine
  harte buildx-Grenze, kein Implementierungsdetail. Der `:dev`-Pfad
  bleibt deshalb bewusst einplattformig (`linux/amd64`, unverändert);
  nur der `VERSION=`-Push-Pfad (`release.yml`) wird multi-arch.
- Auf diesem Entwicklungsrechner (Colima) listet der aktive
  buildx-Builder nur `linux/arm64, linux/386` als native Plattformen —
  **kein** `linux/amd64`, kein Binfmt für Cross-Emulation registriert
  (`docker buildx inspect --bootstrap`). Ein lokaler End-zu-Ende-Test
  des künftigen `linux/amd64,linux/arm64`-Push-Builds ist auf diesem
  Rechner ohne vorherige `docker run --privileged --rm tonistiigi/binfmt
  --install all` nicht direkt möglich — siehe §6.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Ein natives Windows-Container-Image (`os: windows`, z. B. Nano
  Server) — mit dem Auftraggeber explizit geklärt: Docker Desktop unter
  Windows nutzt für Linux-Container ohnehin die WSL2-Linux-Engine: kein
  zusätzliches Image nötig. Ein echtes Windows-Image bräuchte einen
  komplett eigenen Dockerfile-Pfad und einen Windows-GitHub-Runner —
  für einen reinen Backend-/CDC-Dienst kein sinnvoller Aufwand.
- Weitere Plattformen (`linux/arm/v7`, `linux/386`, `linux/riscv64`,
  `linux/s390x`, `linux/ppc64le`) — beide Basis-Images tragen sie zwar
  in ihren Manifest-Listen, aber ohne einen benannten Bedarf (kein
  Nutzer-/Betreiber-Wunsch) bleibt die Matrix bei den zwei praktisch
  relevanten Plattformen (Mehrbelastung der Bauzeit sonst ohne
  Gegenwert).
- Der lokale `:dev`-Build-Pfad (`make image` ohne `VERSION`) bleibt
  einplattformig — buildx' `--load` unterstützt keine
  Multi-Platform-Manifestliste (technische Grenze, siehe §1 oben), kein
  Gestaltungsspielraum.
- Multi-Arch für die Beispiel-Client-Images (`examples/Dockerfile`,
  `examples/csharp/Dockerfile`, `examples/kotlin/Dockerfile`) — reine
  Wegwerf-Entwicklungswerkzeuge, kein veröffentlichtes, für Endnutzer
  bestimmtes Artefakt wie das Haupt-Image.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] `Dockerfile`: `deps`/`proto`/`proto-export`/`coverage`-Stufen auf
      `FROM --platform=$BUILDPLATFORM …` umgestellt (native Bauzeit,
      kein Cross-Emulations-Bedarf für Codegenerierung/Tests); `build`-
      Stufe erhält `ARG TARGETOS` / `ARG TARGETARCH` und kompiliert über
      `GOOS=$TARGETOS GOARCH=$TARGETARCH` statt des fest verdrahteten
      `GOOS=linux`. `runtime`-Stufe unverändert (keine `RUN`-Schritte,
      braucht keine Emulation). Real geprüft: `make image` (ohne
      `VERSION`, unverändert einplattformig) baut weiterhin grün und
      liefert ein lauffähiges `:dev`-Image auf diesem Rechner.
- [ ] `Makefile`s `image:`-Target: der `VERSION=`-Zweig bekommt
      `--platform linux/amd64,linux/arm64`; der `VERSION`-lose `:dev`-
      Zweig bleibt unverändert einplattformig. Real gegen eine
      Test-Registry (`registry:2`, analog dem bereits etablierten
      Content-Mirror-Testmuster aus `release-version-und-workflow`)
      mit `--push` verifiziert: beide Plattform-Manifeste real
      vorhanden, `docker buildx imagetools inspect` zeigt
      `linux/amd64` **und** `linux/arm64/v8` in der Index-Manifestliste,
      und ein real gezogenes `linux/arm64`-Image startet und beantwortet
      eine reale Anfrage (kein bloßer `docker manifest inspect`-Beleg,
      echter Container-Start pro Plattform).
- [ ] `docs/user/releasing.md` §4 „Was beim Release automatisch
      passiert" nennt die Plattform-Abdeckung; `harness/README.md`s
      `make image`-Zeile aktualisiert.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für `docs/user/releasing.md`/`harness/README.md` —
      entfällt als eigener Punkt, da bereits §2 oben dieselbe Zeile
      explizit als DoD-Kriterium trägt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `Dockerfile` | update | `--platform=$BUILDPLATFORM` auf den Bau-Stufen, `TARGETOS`/`TARGETARCH`-Cross-Compile in `build`. |
| `Makefile` (`image:`-Target) | update | `--platform linux/amd64,linux/arm64` nur im `VERSION=`-Zweig. |
| `docs/user/releasing.md` | update | §4 nennt die Plattform-Abdeckung. |
| `harness/README.md` | update | `make image`-Zeile nennt die Plattform-Abdeckung des `VERSION=`-Pfads. |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): sofort — keine Abhängigkeit von
anderen offenen Slices; `ADR-0051` ist bereits vollständig umgesetzt
(`welle-release-pipeline-adr-0051`, `done/`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der
  `metadata-file`-Digest-Extraktionsschritt in `make image` (siehe §6)
  erweist sich als grundlegend anders für Multi-Platform-Builds, als
  hier angenommen, und braucht eine eigene Analyse-Runde — dann
  Zerlegung in „Dockerfile-Cross-Compile" und „Digest-Semantik
  Multi-Arch" als zwei Slices.
- `in-progress` → `open` (blockiert — Carveout?): das lokale
  buildx-Setup lässt sich ohne Systemänderung (Binfmt-Registrierung)
  nicht für einen realen Multi-Platform-Push-Test nutzen, und ein
  alleiniger CI-Beleg (ohne jede lokale Verifikation vor dem Push) wird
  als nicht ausreichend bewertet — dann Carveout mit der Begründung,
  dass die reale Erstverifikation auf den ersten CI-Lauf verschoben
  wird (`AGENTS.md` §3.10-analog).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Digest-Semantik bei Multi-Platform-Builds:** `make image`s
  bestehende Digest-Extraktion (`grep -o '"containerimage.digest"...'`
  aus `--metadata-file`) liest heute den Digest **eines** Manifests. Bei
  `--platform linux/amd64,linux/arm64` liefert buildx stattdessen (oder
  zusätzlich) den Digest der **Index-Manifestliste** — ob der bestehende
  `grep`/`head -1`-Mechanismus dabei weiterhin genau **einen**,
  eindeutigen Wert liefert (den Index-Digest, nicht zufällig den eines
  Einzelmanifests), ist vor der Implementierung nicht bewiesen. **Ausgang:**
  weiter offen bis zur Implementierung, dort real mit `docker buildx
  imagetools inspect` gegen den tatsächlich geschriebenen
  `harness/image-hash.txt`-Wert gegenzuprüfen.
- **Lokale Verifizierbarkeit eingeschränkt:** Der aktive
  Entwicklungsrechner (Colima) hat `linux/amd64` nicht als natives
  Ziel und kein Binfmt für Cross-Emulation registriert — ein realer
  Multi-Platform-Push-Test hier braucht entweder eine einmalige
  Systemänderung (`docker run --privileged --rm tonistiigi/binfmt
  --install all`, außerhalb des Repo-Umfangs, Auswirkung auf die
  Docker-Umgebung des Nutzers) oder verlässt sich auf den ersten realen
  CI-Lauf. **Ausgang:** weiter offen, strukturell — Implementer
  entscheidet zwischen beiden Wegen und benennt die Wahl explizit im
  Bericht.
- **Bauzeit:** Multi-Platform-Builds über QEMU-freie Cross-Compilation
  sollten die CI-Bauzeit nur moderat erhöhen (Go-Compile ist die einzige
  echte Cross-Stufe, `apk add`/`protoc`/Tests bleiben nativ), das ist
  aber unbewiesen, bis ein realer `release.yml`-Lauf die Bauzeit zeigt.
  **Ausgang:** weiter offen, real erst mit dem ersten Post-Push-Lauf
  messbar (`AGENTS.md` §3.10 analog).
- `AGENTS.md` §3.10 gilt für die geänderte `release.yml`-Struktur
  (mittelbar über `make image`) unverändert — ein realer, grüner
  Post-Push-Lauf (nächster echter Release-Tag) bleibt der einzige volle
  Beleg. **Ausgang:** weiter offen, strukturell (derselbe Fall wie
  `BEO-PGC/github-actions-unverifizierbar-lokal`, bereits verkörpert).

## 7. Closure-Notiz

*(wird bei Bearbeitung gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area „CI/CD-Pipeline"
(`Dockerfile`, `Makefile`, `release.yml`-Verhalten mittelbar über
`make image`) — bereits mehrfach berührt (`welle-release-pipeline-adr-0051`),
Schwelle erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/image-digest-nichtdeterminismus-erzeugt-merge-konflikt` (bereits
verkörpert, `ADR-0103`/`ADR-0044`) betrifft dieselbe Digest-Extraktions-
Mechanik wie das oben benannte Risiko, aber aus anderem Anlass
(Merge-Konflikt durch `--load`-vs-`--push`-Digest-Unterschied, nicht
Multi-Platform-Index-Digest) — kein direkter Treffer, aber strukturell
verwandt, im Blick zu behalten. `BEO-PGC/github-actions-unverifizierbar-lokal`
(bereits verkörpert, siehe §6) betrifft diese Sub-Area ebenfalls. Kein
weiterer Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
