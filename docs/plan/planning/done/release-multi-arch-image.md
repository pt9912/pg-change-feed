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

**Verantwortlich:** Implementer-Agent (priorisiert 2026-09-19, direkter
Auftrag: "ja, fang an").

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
  bleibt bewusst einplattformig — nicht wegen einer harten buildx-Grenze
  (Reviewer-Finding F-1 zu diesem Slice: mit aktiviertem
  containerd-Image-Store kann `--load` real eine Multi-Platform-
  Manifestliste laden, real gegenprüft), sondern für eine konsistente,
  schnelle lokale
  Dev-Iteration über unterschiedliche Docker-Setups hinweg — nicht
  jeder Entwicklerrechner hat den containerd-Image-Store aktiv. Der
  `:dev`-Pfad baut weiterhin für die native Host-Plattform (auf diesem
  Rechner `linux/arm64`, nicht `linux/amd64`); nur der
  `VERSION=`-Push-Pfad (`release.yml`) wird multi-arch.
- Auf diesem Entwicklungsrechner (Colima) listete der aktive
  buildx-Builder ursprünglich nur `linux/arm64, linux/386` als native
  Plattformen — **kein** `linux/amd64`, kein Binfmt für Cross-Emulation
  registriert (`docker buildx inspect --bootstrap`). Real während der
  Implementierung gelöst (siehe §6) — ein echter
  `linux/amd64,linux/arm64`-Push-Build gelang danach vollständig.

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
  einplattformig — eine bewusste Design-Entscheidung für konsistente,
  schnelle lokale Dev-Iteration über unterschiedliche Docker-Setups
  hinweg (siehe §1 oben), keine technisch erzwungene Grenze.
- Multi-Arch für die Beispiel-Client-Images (`examples/Dockerfile`,
  `examples/csharp/Dockerfile`, `examples/kotlin/Dockerfile`) — reine
  Wegwerf-Entwicklungswerkzeuge, kein veröffentlichtes, für Endnutzer
  bestimmtes Artefakt wie das Haupt-Image.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] `Dockerfile`: `deps`-Stufe (und die davon abgeleiteten `proto`/
      `proto-export`/`coverage`/`build`) auf `FROM --platform=$BUILDPLATFORM
      …` umgestellt (native Bauzeit, kein Cross-Emulations-Bedarf für
      Codegenerierung/Tests); `build`-Stufe erhält `ARG TARGETOS`/
      `ARG TARGETARCH` und kompiliert über `GOOS=$TARGETOS
      GOARCH=$TARGETARCH` statt des fest verdrahteten `GOOS=linux`.
      `runtime`-Stufe unverändert (keine `RUN`-Schritte, braucht keine
      Emulation). Real geprüft: `make image` (ohne `VERSION`,
      unverändert einplattformig) baut weiterhin grün, liefert ein
      lauffähiges `:dev`-Image, und der reale Build-Log zeigt
      `GOARCH=arm64` korrekt automatisch übernommen (kein Regressions-
      Risiko für den bestehenden Pfad).
- [x] `Makefile`s `image:`-Target: der `VERSION=`-Zweig bekommt
      `--platform linux/amd64,linux/arm64`; der `VERSION`-lose `:dev`-
      Zweig bleibt unverändert einplattformig (bewusste Design-
      Entscheidung für konsistente, schnelle lokale Dev-Iteration, keine
      technisch erzwungene Grenze — Reviewer-Finding F-1: mit
      aktiviertem containerd-Image-Store kann `--load` real eine
      Multi-Platform-Manifestliste laden, real gegengeprüft). Real
      gegen eine lokale Test-Registry (`registry:2`) mit echtem
      `--push --platform linux/amd64,linux/arm64` verifiziert (nach
      einmaliger, danach wieder zurückgenommener Binfmt-Registrierung
      für `amd64`-Emulation via Rosetta, siehe §7): `docker buildx
      imagetools inspect` zeigt beide Plattform-Manifeste real in der
      Index-Manifestliste; **beide** Images real gezogen und gestartet
      (`docker run --platform linux/amd64 …`/`--platform linux/arm64 …`)
      — beide liefern identische, korrekte Anwendungsausgabe. Zusätzlich
      real geprüft: die bestehende Digest-Extraktion aus
      `--metadata-file` liefert bei Multi-Platform-Builds unverändert
      genau **einen** `containerimage.digest`-Wert, der exakt dem
      sha256 des von der Registry abgerufenen rohen Index-Manifests
      entspricht (kein Code-Fix nötig — das in §6 benannte Risiko ist
      damit entfallen, nicht nur vermutet).
- [x] `docs/user/releasing.md` §4 „Was beim Release automatisch
      passiert" nennt die Plattform-Abdeckung; `harness/README.md`s
      `make image`-Zeile aktualisiert.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      1 HIGH (F-1: die Behauptung „`--load` kann keine Multi-Platform-
      Manifestliste laden" war als unbedingte Tatsache formuliert, real
      aber vom Storage-Treiber abhängig — mit containerd-Image-Store
      gelingt es; die Design-Entscheidung selbst blieb unverändert
      richtig, nur ihre Begründung wurde korrigiert) und 1 MEDIUM (F-2:
      Reconciliation-Register-Checkbox ohne „entfällt"-Vermerk) sowie
      1 LOW (F-3: grenzwertige Vorher/Nachher-Sprache im Dockerfile-
      Kommentar) in derselben Fixrunde behoben, kein offenes HIGH.
- [ ] Doku-Update für `docs/user/releasing.md`/`harness/README.md` —
      entfällt als eigener Punkt, da bereits §2 oben dieselbe Zeile
      explizit als DoD-Kriterium trägt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine Reconciliation-Datei in diesem Repo (kein Brownfield-Bootstrap).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-PGC/binfmt-werkzeuge-inkompatibel/` angelegt (1. Beleg, unter der 3×-Schärfungsschwelle): `tonistiigi/binfmt` gefolgt von `multiarch/qemu-user-static --reset` beschädigte real den lokalen Docker-Runtime, `colima restart` behob es.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft (wellenloser Slice): kein `liegt in`-Feld in diesem Slice (Anker vacuously erfüllt), kein Folge-Slice genannt, alle drei referenzierten `BEO-PGC`-Kennungen (`binfmt-werkzeuge-inkompatibel`, `image-digest-nichtdeterminismus-erzeugt-merge-konflikt`, `github-actions-unverifizierbar-lokal`) existieren real mit nicht-leerem `evidence/`.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `Dockerfile` | update | `--platform=$BUILDPLATFORM` auf den Bau-Stufen, `TARGETOS`/`TARGETARCH`-Cross-Compile in `build`. |
| `Makefile` (`image:`-Target) | update | `--platform linux/amd64,linux/arm64` nur im `VERSION=`-Zweig. |
| `docs/user/releasing.md` | update | §4 nennt die Plattform-Abdeckung. |
| `harness/README.md` | update | `make image`-Zeile nennt die Plattform-Abdeckung des `VERSION=`-Pfads. |
| `docs/plan/planning/observations/BEO-PGC/binfmt-werkzeuge-inkompatibel/` | neu (Plan-Nachzug) | realer Docker-Runtime-Zwischenfall bei der lokalen Verifikation, siehe §6. |

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
  eindeutigen Wert liefert. **Ausgang:** entfallen — real gegen eine
  lokale Test-Registry geprüft: `--metadata-file` trägt bei einem
  Multi-Platform-Build genau **einen** `containerimage.digest`-Eintrag,
  identisch mit dem sha256 des von der Registry abgerufenen rohen
  Index-Manifests. Kein Code-Fix nötig.
- **Lokale Verifizierbarkeit eingeschränkt:** Der aktive
  Entwicklungsrechner (Colima) hatte `linux/amd64` zunächst nicht als
  natives Ziel und kein Binfmt für Cross-Emulation registriert. **Ausgang:**
  eingetreten, real gelöst — mit realem Zwischenfall: `docker run
  --privileged --rm tonistiigi/binfmt --install all` registrierte
  zunächst erfolgreich alle Emulatoren (inkl. `rosetta`), der
  Buildx-Builder zeigte `linux/amd64` danach aber weiterhin nicht als
  unterstützt. Ein zweiter Versuch mit `multiarch/qemu-user-static
  --reset -p yes` (eine gängige, aber mit `tonistiigi/binfmt`
  unverträgliche Alternative) beschädigte danach real den gesamten
  lokalen Docker-Runtime — jeder Container-Start scheiterte mit
  `exec format error`, auch für triviale native Images
  (`hello-world`). Behoben durch `colima restart` (stellt den
  VM-Kernelzustand inkl. `binfmt_misc` sauber zurück, keine manuelle
  Handarbeit nötig) — danach lief `linux/amd64` real über Rosetta
  (`docker run --platform linux/amd64 alpine uname -m` → `x86_64`), und
  ein echter `docker buildx build --platform linux/amd64,linux/arm64
  --push`-Lauf gegen eine lokale Registry gelang vollständig, obwohl
  `docker buildx inspect --bootstrap` `linux/amd64` weiterhin nicht in
  der Platform-Liste des Builders zeigt (Laufzeit-Emulation über
  Rosetta funktioniert unabhängig von dieser Anzeige). **Lehre:**
  `tonistiigi/binfmt` und `multiarch/qemu-user-static` nicht
  nacheinander auf demselben Colima-Setup ausführen — ein neuer Beleg
  für das Beobachtungs-Register (siehe §7).
- **Bauzeit:** Multi-Platform-Builds über QEMU-freie Cross-Compilation
  erhöhen die Bauzeit real nur moderat — der lokale Test-Build zeigte
  den amd64-Cross-Compile-Schritt (`GOARCH=amd64` auf dem nativen
  arm64-Host) in 4,5s, praktisch identisch mit dem nativen
  arm64-Compile. **Ausgang:** entfallen als Sorge für die Compile-Stufe
  selbst; die tatsächliche Gesamt-Bauzeit in `release.yml` (inkl.
  `apk add`/`protoc`/Tests je Plattform-Durchlauf) bleibt real erst mit
  dem ersten Post-Push-Lauf messbar — dieser Rest-Anteil bleibt
  **weiter offen** (`AGENTS.md` §3.10 analog).
- `AGENTS.md` §3.10 gilt für die geänderte `release.yml`-Struktur
  (mittelbar über `make image`) unverändert — ein realer, grüner
  Post-Push-Lauf (nächster echter Release-Tag) bleibt der einzige volle
  Beleg. **Ausgang:** weiter offen, strukturell (derselbe Fall wie
  `BEO-PGC/github-actions-unverifizierbar-lokal`, bereits verkörpert).
- Reviewer-Finding F-1 (HIGH, `AGENTS.md` §3.12): Die ursprüngliche
  Begründung für den einplattformigen `:dev`-Pfad („buildx' `--load`
  kann keine Multi-Platform-Manifestliste laden — harte Grenze") war
  als unbedingte Tatsache formuliert, real aber vom Storage-Treiber
  abhängig — mit aktiviertem containerd-Image-Store gelingt es (real
  auf diesem Rechner gegengeprüft, eigene Reproduktion unabhängig vom
  Reviewer). **Ausgang:** eingetreten, real behoben — alle vier
  betroffenen Träger (`Makefile`-Kommentar, `harness/README.md`,
  dieser Slice-Plan §1/§2) auf die zutreffende, bedingte Formulierung
  umgestellt; die Design-Entscheidung selbst (`:dev` bleibt
  einplattformig) blieb unverändert richtig — nur ihre Begründung war
  falsch (Konsistenz über unterschiedliche Docker-Setups statt einer
  nicht existierenden harten Grenze).

## 7. Closure-Notiz

- **Was hat funktioniert:** Reale, hands-on-Verifikation vor UND während
  der Implementierung deckte zwei genuine, sonst unentdeckte Fakten auf:
  (1) beide gepinnten Basis-Images sind bereits Multi-Arch-
  Manifestlisten — kein Re-Pin nötig, real per `docker buildx imagetools
  inspect` vor Code-Änderung geprüft; (2) der reale
  `--platform linux/amd64,linux/arm64 --push`-Test gegen eine lokale
  Test-Registry bewies nicht nur, dass der Mechanismus funktioniert,
  sondern auch, dass die bestehende Digest-Extraktion unverändert
  korrekt bleibt (Index-Digest statt Einzelmanifest-Digest) — ein Risiko,
  das ohne diesen Test bis zum ersten echten Release unbewiesen
  geblieben wäre.
- **Was ging anders als geplant:** Der Reviewer fand 1 HIGH (F-1: die
  Begründung für den einplattformigen `:dev`-Pfad war als unbedingte
  Tatsache formuliert — real hängt es vom Storage-Treiber ab, mit
  containerd-Image-Store gelingt `--load` auch für Multi-Platform-Builds;
  die Design-Entscheidung selbst blieb richtig, nur ihre Begründung war
  falsch) und 1 MEDIUM (F-2: Reconciliation-Checkbox ohne
  „entfällt"-Vermerk) sowie 1 LOW (F-3: grenzwertige Kommentar-Sprache).
  Alle drei in einer Fixrunde behoben. Zusätzlich, außerhalb der
  Reviewer-Funde: während der lokalen Verifikation beschädigte eine
  Kombination aus `tonistiigi/binfmt` und `multiarch/qemu-user-static
  --reset` real den lokalen Docker-Runtime (jeder Container-Start
  scheiterte, auch für triviale native Images) — behoben durch
  `colima restart`, dokumentiert als neue Beobachtung.
- **Steering-Loop-Eintrag:** kein neuer Sensor, keine geschärfte Regel —
  F-1 ist eine Instanz des bereits bestehenden `AGENTS.md` §3.12
  (Herkunft von Aussagen), keine neue Fehlerklasse; der Docker-Runtime-
  Zwischenfall ist der 1. Beleg einer neuen, noch unter der
  3×-Schärfungsschwelle liegenden Beobachtung
  (`BEO-PGC/binfmt-werkzeuge-inkompatibel`).
- **Beobachtungs-Register (`../observations/`):** neues Verzeichnis
  `BEO-PGC/binfmt-werkzeuge-inkompatibel/` angelegt, 1. Beleg
  (`evidence/slice-release-multi-arch-image.md`).
- **Folge-Slices:** keine.
- **Risiken aus §6:** fünf Risiken, fünf Ausgänge — (1)
  Digest-Semantik bei Multi-Platform-Builds → **entfallen**, real
  geprüft und unauffällig; (2) lokale Verifizierbarkeit eingeschränkt →
  **eingetreten**, real gelöst (mit dokumentiertem Zwischenfall); (3)
  Bauzeit → **entfallen** für die Compile-Stufe selbst, **weiter offen**
  für die Gesamt-CI-Bauzeit (real erst mit dem ersten Post-Push-Lauf
  messbar); (4) `AGENTS.md` §3.10 → **weiter offen**, bereits
  verkörpert; (5) Reviewer-Finding F-1 (unbedingte Tatsachenbehauptung)
  → **eingetreten**, real behoben.
- **Drei Paarungen:** wellenloser Slice, hier geprüft (nicht bei einer
  Welle-Closure) — kein `liegt in`-Feld in diesem Slice (Anker vacuously
  erfüllt), kein Folge-Slice genannt, alle drei referenzierten
  `BEO-PGC`-Kennungen (`binfmt-werkzeuge-inkompatibel`,
  `image-digest-nichtdeterminismus-erzeugt-merge-konflikt`,
  `github-actions-unverifizierbar-lokal`) existieren real mit
  nicht-leerem `evidence/`.

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
