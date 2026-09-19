# Review-Report: release-multi-arch-image — 2026-09-19

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/release-multi-arch-image.md`), `ADR-0051`
(additive Erweiterung) und `AGENTS.md` Hard Rules (Modul 10 §Drei
Review-Arten).

**Gegenstand:** Commit `4c3e4ada30575286ac1e4342403ad9ae85915c05` (Parent
`caf8cd133dde194afe8228f12c24da398d4c3a88`), Slice
`release-multi-arch-image`, ohne Welle.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-19.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/release-multi-arch-image.md` (Slice-Plan
  §1–§3, §6 Risiken, §8 Sub-Area/Modus)
- `ADR-0051` (Accepted) — Entscheidung 1/2 (Multi-Registry-Push),
  „ein Build, Content-Mirror"-Mechanismus
- `ADR-0044`/`ADR-0103` — Image-Beleg-Semantik, lokaler `:dev`-Pfad
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentar-Disziplin), §3.9
  (Exit-Code-Disziplin), §3.12 (Herkunft von Aussagen)
- `docs/reviews/review-slice-release-upstream-drift.md` — Vorgänger-Review
  derselben CI/CD-Pipeline-Sub-Area, als Sorgfalts-Referenz
- `docs/plan/planning/observations/BEO-PGC/` — Register-Konventionen für
  1×-Einträge (Stichprobe über mehrere Geschwister-Verzeichnisse)
- `harness/conventions.md` (MR-000 ID-Schema, MR-002 Slice-Kennungen)

**Eigenständig durchgeführte Prüfungen (nicht nur Commit-Message
übernommen):**

- `docker buildx imagetools inspect` gegen beide gepinnten Basis-Images
  real ausgeführt: `golang:1.27-alpine@sha256:cf6fca66…` trägt
  `linux/arm64/v8` (neben sechs weiteren Plattformen);
  `gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872…` trägt
  `linux/arm64/v8` (neben `amd64`/`arm/v7`/`s390x`/`ppc64le`) — Implementer-
  Behauptung 1 bestätigt.
- Realer `docker buildx build --push --platform linux/amd64,linux/arm64`
  gegen eine lokale `registry:2`-Testinstanz (eigenes, minimales
  Test-Dockerfile) durchgeführt: `--metadata-file` liefert genau **einen**
  `containerimage.digest`-Eintrag, exakt identisch mit dem sha256 des über
  `docker buildx imagetools inspect --raw` abgerufenen rohen
  Index-Manifests — Implementer-Behauptung 4 (Digest-Extraktion)
  unabhängig bestätigt.
- Realer `docker buildx build --load --platform linux/amd64,linux/arm64`
  auf demselben Colima/buildx-Setup, das der Implementer für die eigene
  Verifikation genutzt hat, durchgeführt (siehe F-1) — widerlegt die im
  Diff wiederholt behauptete „harte Grenze".
- `make gates` real ausgeführt, Exit-Code direkt (ungepiped) geprüft: `0`.
  `baseline-verify` (v6.9.0, 54 Dateien), `docs-check`/d-check
  (783 Dateien, 0 Befunde), `commit-traceability` (5 Commits, alle mit
  Struktur-ID), `generated-sync` (Proto-Erzeugnis byte-gleich),
  `coverage-gate` (82,80 % ggü. Schwelle 80 %), `a-check` (0 Befunde) —
  Implementer-Behauptung 5 bestätigt.
- Working Tree nach allen lokalen Test-Läufen wieder sauber (`git status
  --short` leer außerhalb des nicht-committeten, lokalen
  `harness/image-hash.txt`).

---

## Findings

### F-1 — „`buildx --load` kann keine Multi-Platform-Manifestliste laden — harte Grenze" ist als unbedingte Tatsache formuliert, aber auf dem realen, für diesen Slice benutzten Docker-Setup nachweislich falsch

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.12 (Herkunft von Aussagen in Trägern, Instanz B —
  Tatsachenbehauptung ohne Beleg-Anker)
- `pfad`: `Makefile:33-35` (Kommentarblock über dem `image:`-Target,
  VERSION-loser Zweig), `harness/README.md` (`make image`-Zeile, Klausel
  „kein Widerspruch zu `ADR-0103`s Begründung … (buildx' `--load` kann
  keine Multi-Platform-Manifestliste laden)"),
  `docs/plan/planning/in-progress/release-multi-arch-image.md:52-56` (§1),
  ebd. `:109-110` (§2 DoD zweiter Punkt)
- `befund`: Alle vier Stellen behaupten unbedingt und im Indikativ, `docker
  buildx build --load` könne grundsätzlich keine
  Multi-Platform-Manifestliste laden („harte buildx-Grenze", „kein
  Gestaltungsspielraum"). Real auf genau dem Colima/buildx-Setup geprüft,
  das der Implementer für alle übrigen realen Verifikationen dieses Slice
  benutzt hat (`docker buildx build --load --platform
  linux/amd64,linux/arm64 -t mp-test:load .`): Der Build gelang, `docker
  image inspect` zeigt einen `application/vnd.oci.image.index.v1+json`-
  Descriptor (echte Multi-Platform-Manifestliste, lokal geladen), und
  sowohl `docker run --platform linux/amd64 mp-test:load` (`x86_64`) als
  auch `--platform linux/arm64` (`aarch64`) liefern die jeweils korrekte
  native Ausgabe. Grund: Dieser Docker-Daemon nutzt den
  `io.containerd.snapshotter.v1`-Treiber (`docker info` →
  `driver-type io.containerd.snapshotter.v1`) — mit aktiviertem
  Containerd-Image-Store unterstützt `--load` seit mehreren Docker-/
  Buildx-Versionen sehr wohl Multi-Platform-Ergebnisse; die Einschränkung
  gilt nachweislich nur für den klassischen (nicht-Containerd)
  Image-Store. Die Behauptung ist also nicht per se falsch, aber als
  unbedingte, umgebungsunabhängige Tatsache formuliert, obwohl sie
  bedingt ist — und ausgerechnet auf der Umgebung widerlegbar, mit der
  der Implementer alle übrigen Behauptungen dieses Slice real geprüft
  hat, ohne dass diese eine, ebenso leicht mit vorhandenem Werkzeug
  prüfbare Behauptung getestet wurde.
- `verifizierbar`: ja — `docker info --format '{{.DriverStatus}}'` bzw.
  `docker info | grep -i "driver-type"` gegen den aktiven Daemon, plus ein
  realer `docker buildx build --load --platform
  linux/amd64,linux/arm64`-Testlauf.
- `klasse`: Beleg trägt seinen Satz nicht (unbedingte Tatsachenbehauptung
  ohne Beleg-Anker, bei eigener Prüfung umgebungsabhängig falsch)

**Einordnung — was dieser Fund nicht ist:** Die Design-Entscheidung, den
`:dev`-Pfad einplattformig zu belassen, bleibt dadurch nicht zwangsläufig
falsch — sie lässt sich unabhängig davon rechtfertigen (Konsistenz mit dem
Verhalten auf GitHub-hosted Runnern ohne Containerd-Image-Store,
Einfachheit, kein Nutzerbedarf). Der Fund betrifft ausschließlich die
**Begründung**: Eine als „harte Grenze" bezeichnete, unbedingte
technische Tatsache, verteilt über vier Träger, hält der eigenen Prüfung
auf der eigenen Testumgebung nicht stand.

### F-2 — DoD-Checkbox „Reconciliation-Register" nicht nach dem etablierten Repo-Muster nachgezogen, obwohl die beiden Nachbar-Checkboxen im selben Commit bearbeitet wurden

- `kategorie`: MEDIUM
- `quelle`: Maintainability / DoD-Konsistenz mit etabliertem Muster
- `pfad`: `docs/plan/planning/in-progress/release-multi-arch-image.md:136`
- `befund`: Die Zeile „Reconciliation-Register (`../reconciliation.md`)
  fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst**"
  bleibt unverändert `[ ]`, ohne Anmerkung. `docs/plan/planning/
  reconciliation.md` existiert in diesem Repo nicht (real geprüft,
  `ls` → „No such file or directory"). Mindestens drei bereits
  abgeschlossene Slices (`slice-nats-drittstream-example-csharp-kotlin`,
  `slice-beispiele-start-architect-entscheidung`,
  `slice-nats-drittstream-example-go`) markieren exakt diese Zeile in
  genau diesem Fall als `[x]` mit einer expliziten
  „entfällt: Datei existiert … nicht"-Anmerkung. Der Implementer hat im
  selben Commit die beiden direkt benachbarten Zeilen
  (Beobachtungs-Register, Risiko-Ausgang) bearbeitet, diese Zeile aber
  unangetastet gelassen — kein neuer Fehler, aber eine Inkonsistenz zu
  einem bereits dreifach etablierten Muster in unmittelbarer Nähe der
  tatsächlich bearbeiteten Zeilen.
- `verifizierbar`: ja — `ls docs/plan/planning/reconciliation.md` (Datei
  fehlt) und `grep -rn "Reconciliation-Register" docs/plan/planning/done/*.md`
  (Präzedenzfälle) sind beide mechanisch nachvollziehbar.
- `klasse`: DoD-Checkbox nicht nach etabliertem Muster nachgezogen

### F-3 — Grenzwertige Vorher/Nachher-Sprache im `build`-Stufen-Kommentar

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Kommentar-Disziplin) / Maintainability
- `pfad`: `Dockerfile:112-113`
- `befund`: „… statt des vormals fest verdrahteten `GOOS=linux`: …" nennt
  den entfernten Wert explizit, ähnlich dem in §3.7 als unzulässig
  genannten Muster „die frühere Fassung prüfte nur die Länge". Der Rest
  des Kommentars beschreibt aber überwiegend den geltenden Zustand
  (Kopplung an `--platform=$BUILDPLATFORM` von `deps`, Begründung für
  natives Cross-Compiling) — die historische Nebenbemerkung ist kurz und
  dient der Einordnung, nicht der alleinigen Begründung. Es existiert
  bereits ein nahezu identisches, unbeanstandetes Vorbild im selben
  Dockerfile (`Dockerfile:49`, „vormals `docker run -v` in den
  Bind-Mount …, siehe [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md)"), das denselben Stil verwendet und in
  einem früheren Slice unwidersprochen blieb — kein neues, sondern ein
  wiederkehrendes Grenzmuster, deshalb LOW statt HIGH trotz der im Skill
  genannten HIGH-Klasse für „abwesenden Text".
- `verifizierbar`: nein — reine Stilfrage, kein Gate prüft
  Kommentar-Zeitform.
- `klasse`: Wiederholung eines Musters, das schon zweimal LOW war (hier:
  Erstauftreten dieser konkreten Sub-Klasse, mit einem bereits
  bestehenden, strukturell identischen Repo-Vorbild)

## Negativbefunde

- geprüft, ohne Befund: Dockerfile-Korrektheit für native Cross-Compilation
  — `deps`/`proto`/`proto-export`/`coverage`/`build` sind alle transitiv
  `FROM deps` (interne Stufe, kein externes Image-`FROM`); sie erben
  `--platform=$BUILDPLATFORM` von `deps`, ohne eigene `--platform`-Angabe
  zu brauchen — dokumentiertes Buildx-Verhalten für von einer bereits
  plattform-gepinnten internen Stufe abgeleitete `FROM`-Zeilen (kein
  erneutes Pull/Resolve, Fortsetzung derselben Stufe). `runtime`
  (`FROM gcr.io/distroless/...`, kein `--platform`) zieht automatisch die
  Ziel-Plattform-Variante — Default-Verhalten für ein externes Image ohne
  `--platform`. `ARG TARGETOS`/`ARG TARGETARCH` sind in `build` korrekt
  vor ihrer Nutzung deklariert (Buildx' vordefinierte, global verfügbare,
  aber pro Stufe re-deklarationspflichtige Build-Args).
- geprüft, ohne Befund: Basis-Image-Multi-Arch-Behauptung (Implementer-
  Behauptung 1) — beide gepinnten Digests real per `docker buildx
  imagetools inspect` bestätigt (`linux/arm64/v8` vorhanden).
- geprüft, ohne Befund: Digest-Extraktion bei Multi-Platform-Builds
  (Implementer-Behauptung 4) — real gegen eine eigene Test-Registry
  reproduziert, exakt ein `containerimage.digest`, identisch mit dem
  Index-Digest.
- geprüft, ohne Befund: `make gates` (Implementer-Behauptung 5) — real
  ausgeführt, Exit 0, alle sechs Segmente grün, Zahlen oben genannt.
- geprüft, ohne Befund: Doku-Widerspruchsfreiheit — `docs/user/releasing.md`
  §4 und `harness/README.md`s `make image`-Zeile beschreiben denselben
  Mechanismus (ein Buildx-Aufruf, zwei Plattformen, Index-Digest) ohne
  Widerspruch zueinander oder zum `Makefile`-Diff; die Windows-Abgrenzung
  in §4 deckt sich mit Plan §1 „Ausdrücklich NICHT in diesem Slice".
- geprüft, ohne Befund: Beobachtungs-Register-Formalia
  (`BEO-PGC/binfmt-werkzeuge-inkompatibel/`) — Verzeichnisstruktur
  (`observation.md`/`state.md`/`evidence/<datei>.md`), Zähler-Text „1×"
  konsistent zwischen `state.md` und der einen `evidence/`-Datei, Format
  deckungsgleich mit mehreren Geschwister-Einträgen bei 1× (Stichprobe:
  `adr-folgepflicht-ohne-traeger-slice`, `dod-checkbox-nachzug-
  architect-pfad`, `geschaetzter-wert-als-grenze`).
- geprüft, ohne Befund: Plan-Deckung — die in `docs/plan/planning/
  in-progress/release-multi-arch-image.md` §3 gelistete Datei-Menge (inkl.
  Plan-Nachzug-Zeile für den neuen Beobachtungs-Eintrag) deckt sich 1:1
  mit den im Commit tatsächlich geänderten Dateien; keine
  Scope-Erweiterung.
- geprüft, ohne Befund: Traceability — Commit-Message nennt `ADR-0051`,
  kein `SPEC-*`/`ARC-*` im Betreff; `harness/conventions.md` MR-002
  (Slice-Kennungen sind Namen) — `release-multi-arch-image` ist bereits
  ein Name.
- geprüft, ohne Befund: §8 Sub-Area-Wahl „CI/CD-Pipeline" — deckt sich mit
  der bereits mehrfach etablierten Benennung in
  `release-hub-description`/`release-version-und-workflow`/
  `release-image-scan`/`release-upstream-drift`; Modus GF konsistent mit
  `harness/conventions.md`s einziger Wildcard-Zeile.
- geprüft, ohne Befund: keine neue Betreiber-Oberfläche eingeführt
  (kein `CDC_*`, keine `cdc.*`-SQL-Funktion, kein neuer Endpunkt) —
  Handbuch-Nachzugsregeln (Skill-HIGH-Bullets) greifen nicht.
- geprüft, ohne Befund: Working Tree nach allen eigenen lokalen
  Verifikations-Läufen (Test-Registry, `--load`-Gegenprobe, `make gates`)
  wieder sauber; keine Test-Artefakte im Repo zurückgelassen.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Beleg trägt seinen Satz nicht (unbedingte
Tatsachenbehauptung ohne Beleg-Anker, bei eigener Prüfung
umgebungsabhängig falsch) · DoD-Checkbox nicht nach etabliertem Muster
nachgezogen · Wiederholung eines Musters, das schon zweimal LOW war
(Vorher/Nachher-Sprache mit bestehendem Repo-Vorbild)

## Verdikt

**Merge-blockierend:** ja — F-1 ist ein Hard-Rule-Verstoß (`AGENTS.md`
§3.12) mit konkretem, auf der eigenen Testumgebung reproduziertem
Gegenbeweis gegen eine als unbedingt formulierte technische Tatsache, die
sich über vier Träger zieht (`Makefile`, `harness/README.md`, Slice-Plan
§1 und §2). Die dahinterliegende Design-Entscheidung (`:dev`-Pfad bleibt
einplattformig) muss dadurch nicht fallen — aber die Begründung „harte
Grenze" ist in dieser unbedingten Form falsch und braucht entweder eine
korrekte, bedingte Formulierung (abhängig vom Containerd-Image-Store) oder
eine andere, tragfähige Begründung für dieselbe Entscheidung.

**Übergabe:** F-1 (HIGH) geht als Rückmeldung an den Implementer
(Fixrunde: die vier Träger — `Makefile`-Kommentar, `harness/README.md`,
Slice-Plan §1/§2 — auf eine bedingte, real geprüfte Formulierung
umstellen; die Entscheidung selbst kann unverändert bleiben, wenn sie
anders begründet wird). Da mit F-1 ein HIGH-Finding mit
Rollen-Rückgabe vorliegt, ziehe ich die DoD-Checkbox „Review
durchgeführt" im Slice-Plan **nicht** nach (Skill-Regel
„DoD-Checkbox-Nachzug ohne Fixrunde" greift nur bei 0 HIGH oder wenn alle
Findings ohne Implementer-Rückgabe weitergereicht werden — hier liegt
beides nicht vor). F-2 und F-3 gehen als zusätzliche, nicht blockierende
Beobachtungen mit in dieselbe Rückmeldung. Die Finding-Klassen gehen bei
Slice-Closure ins Beobachtungs-Register. Dieser Report ersetzt keine
Verifikation gegen die DoD — das bleibt Verifier-Aufgabe (Modul 11).
