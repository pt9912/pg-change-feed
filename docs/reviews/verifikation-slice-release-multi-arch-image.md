# Verifikationsbericht: slice-release-multi-arch-image — 2026-09-19

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/release-multi-arch-image.md` §2) und die
§6-Risiko-Ausgänge, in frischem Kontext. **Nicht** gegen den Diff als
solchen (Reviewer-Aufgabe, abgeschlossen mit
[`review-slice-release-multi-arch-image.md`](review-slice-release-multi-arch-image.md))
und **nicht** gegen realen Bedarf (Validator, hier nicht ausgelöst).

**Gegenstand:** zwei Commits auf `main`, Slice `release-multi-arch-image`,
ohne Welle:

- `4c3e4ada` — ursprünglicher Implementer-Commit (`Dockerfile`
  Multi-Arch-Umbau, `Makefile` `--platform`-Erweiterung, Doku-Updates,
  neuer Beobachtungs-Registereintrag `BEO-PGC/binfmt-werkzeuge-inkompatibel`).
- `8d86eb2d` — Fixrunde nach unabhängigem Review (1 HIGH F-1, 1 MEDIUM
  F-2, 1 LOW F-3), inklusive dem Review-Report selbst.

**Frischer Kontext:** Slice-Plan (vollständig, §1–§8), Review-Report
(vollständig), `Dockerfile`-/`Makefile`-/`harness/README.md`-Diffs beider
Commits im Volltext selbst gelesen. Nichts aus Slice-Plan, Commit-Message
oder Review-Report ungeprüft übernommen: eigener `docker info`-Lauf gegen
den realen Storage-Treiber dieses Rechners, eigener realer
`docker buildx build --load --platform linux/amd64,linux/arm64`-Testlauf
samt Start beider Plattformen, eigener realer
`docker buildx build --push --platform linux/amd64,linux/arm64`-Lauf
gegen eine selbst gestartete lokale `registry:2`-Instanz mit
`--metadata-file`-Digest-Extraktion über denselben `grep`-Befehl wie im
`Makefile` und Vergleich gegen den sha256 des rohen Index-Manifests,
eigener `git diff` über beide Commits für `Dockerfile`/`Makefile`/
`harness/README.md`/Slice-Plan, eigener ungepipter `make gates`-Lauf mit
direkter Exit-Code-Prüfung (`AGENTS.md` §3.9), eigene `git log`-Prüfung
beider Commit-Messages.

---

## 1. DoD-Vertrag (§2) — jede Checkbox einzeln geprüft

### 1.1 `Dockerfile`-Umbau (`--platform=$BUILDPLATFORM`, `TARGETOS`/`TARGETARCH`)

Eigener `git diff caf8cd13 8d86eb2d -- Dockerfile`: `deps` trägt jetzt
`FROM --platform=$BUILDPLATFORM golang:1.27-alpine@sha256:cf6fca66…`;
`build` (`FROM deps`) deklariert `ARG TARGETOS`/`ARG TARGETARCH` und
kompiliert über `GOOS=$TARGETOS GOARCH=$TARGETARCH` statt des vormals
fest verdrahteten `GOOS=linux`. `runtime` (`FROM gcr.io/distroless/...`)
unverändert, kein `--platform`. Der reale
`docker buildx build --load --platform linux/amd64,linux/arm64`-Testlauf
(Abschnitt 2 unten) durchlief exakt diese Kette für beide Plattformen ohne
Fehler und lieferte funktionsfähige Binaries für beide Architekturen —
**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.2 `Makefile`s `image:`-Target (`--platform linux/amd64,linux/arm64` nur im `VERSION=`-Zweig)

Eigene Lektüre `Makefile:47-56`: `VERSION`-Zweig trägt
`--platform linux/amd64,linux/arm64`; `VERSION`-loser `:dev`-Zweig bleibt
unverändert bei `docker buildx build --load` ohne `--platform`-Flag. Beide
Teilbehauptungen real nachgemessen:

- **Digest-Extraktion bei Multi-Platform-Push** — eigener Lauf gegen eine
  selbst gestartete `registry:2`-Instanz (`docker run -d --rm
  --name verify-registry -p 15050:5000 registry:2`):
  ```
  $ docker buildx build --push --platform linux/amd64,linux/arm64 \
      --metadata-file /tmp/verify-meta.json -t localhost:15050/verify:v1 .
  … exporting manifest list sha256:c20bdadbb6a3f3f08ab26ec6d95d5f1a6ecce767b9a007bead5d907167e3adf1 done
  $ grep -o '"containerimage.digest":[[:space:]]*"sha256:[0-9a-f]*' /tmp/verify-meta.json \
      | head -1 | grep -o 'sha256:[0-9a-f]*'
  sha256:c20bdadbb6a3f3f08ab26ec6d95d5f1a6ecce767b9a007bead5d907167e3adf1
  $ docker buildx imagetools inspect localhost:15050/verify:v1 --raw | sha256sum
  c20bdadbb6a3f3f08ab26ec6d95d5f1a6ecce767b9a007bead5d907167e3adf1  -
  ```
  Exakt **ein** `containerimage.digest`-Treffer (`grep … | wc -l` → `1`),
  identisch mit dem sha256 des rohen, von der Registry abgerufenen
  Index-Manifests. Kein Code-Fix nötig — die DoD-Behauptung ist real
  bestätigt, unabhängig vom Implementer- **und** vom Reviewer-Lauf.
- **`--load`-Verhalten auf diesem Rechner** — siehe Abschnitt 2 unten
  (F-1-Fix).

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.3 `docs/user/releasing.md` §4 / `harness/README.md`s `make image`-Zeile

Eigener `grep`: `docs/user/releasing.md:78-79` nennt „Jedes Tag trägt
**beide** Plattformen — `linux/amd64` **und** `linux/arm64` — in einer
einzigen …"; `harness/README.md:126`s `make image`-Zeile nennt
`--platform linux/amd64,linux/arm64` und die Digest-Semantik
widerspruchsfrei zum `Makefile`-Diff. **Ergebnis: Checkbox berechtigt auf
`[x]`.**

### 1.4 „`make gates` grün" — siehe eigenständiger Abschnitt 4 unten

### 1.5 „Review durchgeführt …" — siehe eigenständiger Abschnitt 6 unten

### 1.6 Die verbleibenden `[ ]`-Checkboxen (Doku-Update entfällt, Closure-Notiz, drei Paarungen)

Der Plan liegt weiterhin unter `docs/plan/planning/in-progress/` (eigener
`ls`-Beleg), §7 „Closure-Notiz" trägt weiterhin ausschließlich den
Platzhalter „*(wird bei Bearbeitung gefüllt.)*". Für einen Slice, der
noch nicht nach `done/` gewandert ist, ist das korrekt — diese Punkte
sind Closure-Pflichten, keine Liefer-Punkte, und ihr `[ ]`-Zustand
widerspricht sich nicht mit dem übrigen DoD-Bild. Die
„Doku-Update entfällt"-Zeile bleibt konsistent `[ ]` — Präzedenzfälle
(`release-doku-releasing.md`, `release-hub-description.md`,
`release-image-scan.md`, `release-upstream-drift.md`) nutzen exakt
dasselbe Muster für dieselbe Zeilenklasse.

## 2. F-1-Fix real wirksam — keine unbedingte „harte Grenze"-Behauptung mehr

Eigener `docker info`-Lauf auf diesem Rechner:

```
$ docker info 2>&1 | grep -i "driver-type\|Storage Driver"
 Storage Driver: overlayfs
  driver-type: io.containerd.snapshotter.v1
```

Bestätigt: containerd-Image-Store aktiv, exakt wie im Fix behauptet.
Eigener realer Multi-Platform-`--load`-Testlauf, unabhängig von
Implementer und Reviewer:

```
$ docker buildx build --load --platform linux/amd64,linux/arm64 \
    -t verify-multiarch-load:v1 .
… exporting manifest list sha256:675c423f8126a0776ab626813dd9159e7946c7b45a0ce7ada521f1d097a042fb done
… naming to docker.io/library/verify-multiarch-load:v1 done
$ docker run --rm --platform linux/amd64 --entrypoint /pg-change-feed \
    verify-multiarch-load:v1 --version
pg-change-feed 0.2.0-verdrahtung
$ docker run --rm --platform linux/arm64 --entrypoint /pg-change-feed \
    verify-multiarch-load:v1 --version
pg-change-feed 0.2.0-verdrahtung
```

Der Build gelang, lud real eine Multi-Platform-Manifestliste
(`docker image inspect` zeigt die geladene Image-ID identisch mit dem
gemeldeten `manifest list`-Digest), und **beide** Plattformen liefern
eine reale, korrekte Anwendungsausgabe (`--version`). Test-Image danach
mit `docker rmi -f verify-multiarch-load:v1` entfernt.

Eigener `git diff 4c3e4ada 8d86eb2d` über alle vier vom Reviewer
benannten Träger bestätigt: keiner behauptet mehr eine unbedingte „harte
buildx-Grenze".

- `Makefile:32-46` (Kommentarblock über `image:`): „ob `--load` überhaupt
  eine Multi-Platform-Manifestliste laden kann, haengt vom
  Storage-Treiber ab (real bestaetigt: mit aktiviertem
  containerd-Image-Store gelingt es, es ist keine grundsaetzliche
  buildx-Grenze)" — bedingt formuliert, mit Beleg-Anker.
- `harness/README.md:126` (`make image`-Zeile): „`--load` bleibt für eine
  konsistente, schnelle lokale Dev-Iteration einplattformig — ob `--load`
  überhaupt eine Multi-Platform-Manifestliste laden kann, hängt vom
  Storage-Treiber ab, real bestätigt: mit aktiviertem
  containerd-Image-Store gelingt es, keine grundsätzliche buildx-Grenze"
  — ebenso bedingt.
- Slice-Plan §1 (Zeilen 52-62): „nicht wegen einer harten buildx-Grenze
  (Reviewer-Finding F-1 … mit aktiviertem containerd-Image-Store kann
  `--load` real eine Multi-Platform-Manifestliste laden, real
  gegenprüft), sondern für eine konsistente, schnelle lokale
  Dev-Iteration …" — bedingt, mit Reviewer-Finding-Anker.
- Slice-Plan §2 (Zeilen 112-124): „bewusste Design-Entscheidung für
  konsistente, schnelle lokale Dev-Iteration, keine technisch erzwungene
  Grenze — Reviewer-Finding F-1: mit aktiviertem containerd-Image-Store
  kann `--load` real eine Multi-Platform-Manifestliste laden, real
  gegengeprüft" — bedingt.

Alle vier Stellen benennen jetzt korrekt: Design-Entscheidung
(Konsistenz über heterogene Docker-Setups) getrennt von der zuvor falsch
als unbedingt formulierten technischen Begründung. **Ergebnis: F-1 ist
real und vollständig behoben, unabhängig auf diesem Rechner
reproduziert — nicht nur behauptet.**

## 3. Digest-Extraktion — erneut selbst geprüft (unabhängig von Implementer und Reviewer)

Bereits in Abschnitt 1.2 vollständig dokumentiert: eigener
`registry:2`-Testlauf, `grep`-Extraktion exakt nach demselben Muster wie
`Makefile:52`, ein einziger Treffer, identisch mit dem sha256 des rohen
Index-Manifests (`c20bdadbb6a3f3f08ab26ec6d95d5f1a6ecce767b9a007bead5d907167e3adf1`
in beiden Fällen). Test-Registry (`docker rm -f verify-registry`) und
Test-Image (`docker rmi -f localhost:15050/verify:v1`) danach entfernt,
Metadaten-Datei gelöscht.

## 4. `make gates` real, ungepiped ausgeführt

```
$ git log -1 --format=%H
8d86eb2dd37e6e0e02d7cb7984b15677caec1fae
$ git status --short
(leer)
$ make gates > <log> 2>&1; ec=$?; echo "EXIT_CODE=$ec"
EXIT_CODE=0
```

Einzelbelege aus demselben Lauf (Log vollständig eingesehen):

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 784 Datei(en) geprüft, 0 Befund(e)` |
| `commit-traceability` | `d-check`-Modul `commits`: `784 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability.sh`: `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto)` |
| `a-check` | `gesamt: 0 Befund(e)` |

**Commit-Traceability speziell für beide Slice-Commits geprüft** (Auftrag:
„tragen beide Commits `ADR-0051` im Betreff?"):

```
$ git log -1 --format=%s 4c3e4ada
feat(release): Multi-Arch-Image linux/amd64+linux/arm64 (ADR-0051)
$ git log -1 --format=%s 8d86eb2d
fix(release): Reviewer-Fixrunde für release-multi-arch-image (ADR-0051, 1 HIGH, 1 MEDIUM, 1 LOW)
```

Beide Betreffs tragen `ADR-0051`, kein `SPEC-*`/`ARC-*` im Betreff.
**Ergebnis: Die DoD-Checkbox „`make gates` grün." ist berechtigt auf
`[x]` gesetzt** — real, ungepiped, `EXIT_CODE=0`.

## 5. §6-Risiken — jedes mit zulässigem Ausgang?

Die drei zulässigen Ausgänge (Baseline-Regelwerk
`modul-05-planning-harness.md` §Offene Risiken werden bei Closure
aufgelöst): *eingetreten* → Carveout/Folge-Slice · *entfallen* →
gestrichen mit Begründung · *weiter offen* → Beobachtungs-Register.

| # | Risiko (Kurzform) | Ausgang im Plan | Zulässige Klasse? | Eigene Einschätzung |
|---|---|---|---|---|
| 1 | Digest-Semantik bei Multi-Platform-Builds | „entfallen — real gegen eine lokale Test-Registry geprüft: exakt **ein** `containerimage.digest`, identisch mit dem Index-Digest" | ✓ entfallen | Eigen bestätigt in Abschnitt 1.2/3 oben, identischer Digest-Wert. |
| 2 | Lokale Verifizierbarkeit eingeschränkt (Binfmt) | „eingetreten, real gelöst — mit realem Zwischenfall … `colima restart` behoben" | ✓ eingetreten | Deckt sich mit dem neuen Beobachtungs-Registereintrag (Abschnitt 7 unten); der reale Push-Test in Abschnitt 1.2/3 oben gelang auf demselben Rechner. |
| 3 | Bauzeit | „entfallen als Sorge für die Compile-Stufe selbst; Gesamt-Bauzeit … bleibt weiter offen" | ✓ entfallen / weiter offen (zwei Teilaussagen, beide zulässig) | Nachvollziehbar formuliert, keine Übertreibung. |
| 4 | `AGENTS.md` §3.10 — realer Post-Push-Lauf unverifiziert | „weiter offen, strukturell (derselbe Fall wie `BEO-PGC/github-actions-unverifizierbar-lokal`)" | ✓ weiter offen | Registerverzeichnis existiert real (`ls docs/plan/planning/observations/BEO-PGC/github-actions-unverifizierbar-lokal/`). |
| 5 (F-1) | Reviewer-Finding F-1 — unbedingte Tatsachenbehauptung zu `--load` | „eingetreten, real behoben — alle vier betroffenen Träger … umgestellt" | ✓ eingetreten | Eigen vollständig nachgeprüft, siehe Abschnitt 2 oben. |

**Ergebnis: Alle fünf §6-Risiken tragen eine der drei zulässigen
Ausgangsklassen**, inklusive des neu von der Fixrunde hinzugefügten
F-1-Eintrags.

## 6. DoD-Checkbox „Review durchgeführt" — inhaltlich korrekt gegen den Report?

Review-Report-Summary-Tabelle (eigene Lektüre): HIGH 1, MEDIUM 1, LOW 1,
INFO 0 — Findings F-1 (HIGH, unbedingte Tatsachenbehauptung ohne
Beleg-Anker, `AGENTS.md` §3.12), F-2 (MEDIUM,
Reconciliation-Register-Checkbox ohne „entfällt"-Vermerk), F-3 (LOW,
grenzwertige Vorher/Nachher-Sprache im Dockerfile-Kommentar).

DoD-Checkbox-Text (Plan §2, Zeilen 140-147) nennt exakt dieselbe
Verteilung („1 HIGH … und 1 MEDIUM … sowie 1 LOW … in derselben Fixrunde
behoben, kein offenes HIGH") mit denselben Kurzbeschreibungen je
Finding-ID — Zahlen und Zuordnung stimmen 1:1 mit dem Report überein.

Eigene Prüfung, ob „behoben" für jedes Finding zutrifft:

- **F-1** (HIGH): real behoben, siehe Abschnitt 2 oben — eigenständig auf
  diesem Rechner reproduziert.
- **F-2** (MEDIUM): eigene Lektüre `git diff 4c3e4ada 8d86eb2d` — die
  Reconciliation-Register-Zeile wechselt von `[ ]` ohne Anmerkung zu
  `[x] … entfällt: keine Reconciliation-Datei in diesem Repo (kein
  Brownfield-Bootstrap)" — deckt sich mit dem vom Reviewer zitierten,
  bereits dreifach etablierten Muster (`slice-nats-drittstream-*`).
- **F-3** (LOW): real behoben, siehe `Dockerfile:109-113`-Diff — die
  Formulierung „statt des vormals fest verdrahteten `GOOS=linux`" ist
  entfernt, der Kommentar beschreibt jetzt ausschließlich den geltenden
  Zustand im reinen Indikativ.

**Kein offenes HIGH, kein unbehandeltes MEDIUM/LOW. Ergebnis: Checkbox
berechtigt auf `[x]` gesetzt.**

## 7. Beobachtungs-Registereintrag `BEO-PGC/binfmt-werkzeuge-inkompatibel/` — formal geprüft

Eigener `find`-Lauf: Verzeichnis trägt exakt die drei erwarteten Dateien
(`observation.md`, `state.md`, `evidence/slice-release-multi-arch-image.md`)
— deckungsgleich mit der Struktur aller geprüften Geschwister-Einträge
(`a-check-null-abdeckung/`, `adr-folgepflicht-ohne-traeger-slice/`,
`dod-checkbox-nachzug-architect-pfad/`, `geschaetzter-wert-als-grenze/`).
`state.md` trägt „Zustand: offen — Ausgang: **weiter offen** (unter der
3×-Schärfungsschwelle)" und „Zähler (abgeleitet): 1×
(evidence/slice-release-multi-arch-image.md)" — Zähler-Text konsistent
mit dem Muster der Geschwister-Einträge (`1× (evidence/…)`), Zahl
deckungsgleich mit der tatsächlichen Anzahl der `evidence/`-Dateien
(genau eine). Inhaltlich deckt sich der Beleg mit §6 des Slice-Plans
(Reihenfolge `tonistiigi/binfmt` → `multiarch/qemu-user-static --reset`
→ `colima restart`). **Ergebnis: formal korrekt.**

---

## Verdikt

**DoD erfüllt** — für den aktuellen `in-progress`-Stand des Slice. Alle
sieben beauftragten Prüfpunkte wurden real und unabhängig nachgemessen,
nicht aus Bericht oder Commit-Message übernommen:

1. Jede `[x]`-Checkbox in §2 ist gegen reale Läufe/Dateien geprüft und
   berechtigt gesetzt; die verbleibenden `[ ]`-Checkboxen (Doku-Update
   entfällt, Closure-Notiz, drei Paarungen) sind konsistent mit dem
   `in-progress`-Zustand (Closure-Pflichten, nicht Liefer-Punkte).
2. F-1 ist real und vollständig behoben — eigenständig auf diesem
   Rechner reproduziert: `docker info` bestätigt den containerd-Image-
   Store, ein realer `--load --platform linux/amd64,linux/arm64`-Build
   gelang und lieferte für beide Plattformen korrekte Anwendungsausgabe;
   alle vier betroffenen Träger tragen jetzt eine bedingte, belegte
   Formulierung statt einer unbedingten „harten Grenze".
3. Die Digest-Extraktion ist real gegen eine eigene Test-Registry
   reproduziert: exakt ein `containerimage.digest`-Wert, identisch mit
   dem sha256 des rohen Index-Manifests.
4. `make gates` lief eigenständig, ungepiped, mit `EXIT_CODE=0`; beide
   Slice-Commits tragen `ADR-0051` im Betreff (eigene `git log`-Prüfung).
5. Alle fünf §6-Risiken (inkl. des von der Fixrunde neu hinzugefügten
   F-1-Eintrags) tragen eine der drei zulässigen Ausgangsklassen.
6. Die DoD-Checkbox „Review durchgeführt" ist inhaltlich korrekt gegen
   den tatsächlichen Report (1 HIGH/1 MEDIUM/1 LOW, identische
   Zuordnung) und berechtigt gesetzt — kein offenes HIGH, kein
   unbehandeltes MEDIUM/LOW.
7. Der neue Beobachtungs-Registereintrag
   `BEO-PGC/binfmt-werkzeuge-inkompatibel/` ist formal korrekt
   (Verzeichnisstruktur, Zähler-Text, deckungsgleicher Inhalt zu §6).

**Aufgeräumt:** Test-Image `verify-multiarch-load:v1`
(`docker rmi -f`), Test-Registry-Container `verify-registry`
(`docker rm -f`) und Test-Image `localhost:15050/verify:v1`
(`docker rmi -f`) sowie die Metadaten-Temp-Datei wurden nach jeweiliger
Prüfung entfernt. Ein vorgefundenes, nicht von dieser Verifikation
erzeugtes Image `verify-build-ctx:latest` (Erstellungsdatum 2026-09-18,
vor dieser Sitzung) wurde unverändert belassen — es liegt außerhalb des
Auftragsumfangs dieser Verifikation.

**Freigabe an den Planner:** Der Slice ist DoD-konform für seinen
aktuellen `in-progress`-Stand. Vor Closure (`git mv` nach `done/`)
bleiben die im Plan selbst bereits als offen geführten Closure-Pflichten
zu erfüllen (§7 Closure-Notiz, formaler Risiko-Ausgangs-Nachzug in §6
bereits vollständig, die drei Paarungen bei der nächsten Welle-Closure —
dieses Repo betreibt Wellen an anderer Stelle, auch für wellenlose
Slices gilt deshalb der Welle-Closure-Pfad, nicht der wellenlose Pfad).
