# Slice sdk-kotlin-pack-werkzeug: `make sdk-pack-kotlin` + Träger-Nachzug Pflichtenheft

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-kotlin-lh-fa-sst-009](../welle-sdk-kotlin-lh-fa-sst-009.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md),
[`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
Festlegung 5 (Build-Mechanismus: Docker-only bauen, Gradle,
`maven-publish`), §Konsequenzen Folgepflicht 2/3 (Träger-Nachzug
Pflichtenheft/`harness/README.md`).

**Berührte Spec-Stellen:** [`LH-FA-SST-009.a`](../../../../spec/pflichtenheft.md)
(wird mit diesem Slice für Kotlin/GitHub Packages aufgelöst — die Kennung
selbst bleibt bestehen, ihr Satz „eine dritte Sprache oder ein dritter
Vertriebsweg bleibt offen" wird durch `ADR-0109` und die jetzt reale
Paketierbarkeit falsch, `AGENTS.md` §3.13), `spec/pflichtenheft.md` §6
Externe Verträge (neue Zeile für das Package, voraussichtlich `SPEC-028` —
real gegen den Bestand zu verifizieren, `SPEC-027` ist die zuletzt real
vergebene Nummer, siehe §3).

**Verantwortlich:** Implementer-Agent, 2026-09-20.

**Autor:** Planner-Agent, direkt beauftragt (`ADR-0109` §Konsequenzen
Folgepflicht 1/2/3). **Datum:** 2026-09-20.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein neues, netzlos **nicht** prüfbares Werkzeug-Ziel
`make sdk-pack-kotlin` (Docker-only, `./gradlew build`/`./gradlew test`/
ein Jar-Bau, gepinntes `eclipse-temurin:21-jdk`-Image) erzeugt ein reales
`.jar` (und ggf. Sources-/Javadoc-Jar, falls die reale GitHub-Packages-
Publish-Form das verlangt — real zu recherchieren, kein Vorgriff) aus
`sdks/kotlin/pgchangefeed-kotlin/` — dazu der Träger-Nachzug, den
`ADR-0109` §Konsequenzen Folgepflicht 2/3 fordert:
[`LH-FA-SST-009.a`](../../../../spec/pflichtenheft.md) (§1) und
`spec/pflichtenheft.md` §6 (Externe Verträge, neue `SPEC-<NNN>`-Zeile für
das Package) sowie `harness/README.md` §Werkzeuge (die reale
`make sdk-pack-kotlin`-Zeile).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **GitHub-Packages-Veröffentlichung** — `slice-sdk-kotlin-publish-workflow`
  übernimmt das; dieser Slice erzeugt das Artefakt, veröffentlicht es
  nicht (`ADR-0109` Festlegung 5: „Bauen … bleibt Docker-only …
  Veröffentlichung … braucht einen eigenen, Netz-bindenden Workflow").
- **Neues Gate** — `make sdk-pack-kotlin` bleibt außerhalb von
  `GATE_CHECKS`/`make gates` (`ADR-0109` Festlegung 5, dieselbe
  Begründung wie `make examples-kotlin`: Paket-Bezug braucht Netz).
- **Änderung der beiden Fläche-Slices** — dieser Slice paketiert, was
  `slice-sdk-kotlin-http-client-flaeche`/`slice-sdk-kotlin-grpc-client-flaeche`
  bereits geliefert haben; kein neuer API-Umfang.
- **`docs/user/version.md`-Änderung** — die Server-Version bleibt
  eigenständig geführt (`ADR-0109` Festlegung 4/6).
- **Anlage eines Secrets** — GitHub Packages braucht keins; entfällt
  strukturell (`ADR-0109` §Entscheidung Festlegung 2/5), nicht nur
  „nicht Teil dieses Slice" wie bei den Vorgänger-SDKs.

## 2. Definition of Done

- [x] `make sdk-pack-kotlin` existiert (Docker-only, kein Gate — analog
      `make examples-kotlin`): baut, testet (`./gradlew test` gegen die
      Test-Quellen beider Fläche-Slices) und paketiert (ein Jar-Bau-Task,
      real zu bestimmen anhand des Gradle-`java`/`kotlin("jvm")`-Plugins —
      `./gradlew jar` bzw. `assemble`) `sdks/kotlin/pgchangefeed-kotlin/`
      im gepinnten `eclipse-temurin:21-jdk`-Image; ein roter Test bricht
      den `docker build` mit Exit ≠ 0 ab (Muster `harness/mk/examples.mk`).
      Der gRPC-Anteil braucht denselben zusätzlichen, benannten
      Bau-Kontext wie die beiden Fläche-Slices
      (`--build-context proto=proto`). Das Artefakt wird über einen
      `tar`-Stream-Export aus dem Docker-Bau in ein lokales Verzeichnis
      (z. B. `sdks/kotlin/dist/`, `.gitignore`t) abgelegt — analog dem
      Extraktionsmuster von `make proto-generate`/`make sdk-pack-csharp`/
      `make sdk-pack-python` (`docker run --rm --network none <image> |
      tar -x -C .`, `AGENTS.md` §3.9 Pipe-Disziplin).
- [x] Real recherchiert und entschieden: ob GitHub Packages für dieses
      Package zwingend ein Sources- und/oder Javadoc-Jar neben dem
      Haupt-Jar verlangt, oder ob ein einzelnes Haupt-Jar genügt — die
      Entscheidung steht im Bericht dieses Slice mit ihrem Beleg
      (`AGENTS.md` §3.12 Instanz B), kein Vorgriff aus der ADR (`ADR-0109`
      lässt die Frage bewusst offen).
- [x] Real ausgeführt: ein `.jar` mit einem Dateinamen, der die Koordinate
      `pgchangefeed-kotlin` und die Version `0.1.0` trägt (exakte
      Gradle-Standard-Namensform, z. B. `pgchangefeed-kotlin-0.1.0.jar`),
      existiert nach dem Lauf und ist als Smoke-Beleg im Bericht dieses
      Slice genannt (Datei-Existenz, keine Behauptung — `AGENTS.md` §3.12
      Instanz B).
- [x] `spec/pflichtenheft.md` §1 trägt bei
      [`LH-FA-SST-009.a`](../../../../spec/pflichtenheft.md) einen
      Nachzug-Satz: für Kotlin/GitHub Packages ist die Sprachmatrix-/
      Vertriebsweg-Frage durch
      [`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
      beantwortet und das Package real paketierbar (`AGENTS.md` §3.13,
      kein Streichen der Kennung — sie bleibt die Adresse für eine
      künftige vierte Sprache/einen vierten Vertriebsweg) — **ohne**
      ADR-Bezug im Fließtext, dieselbe Struktur-Regel wie bei den
      bestehenden `LH-FA-SST-009.a`-Nachzugzeilen für C#/Python
      (`.d-check.yml` `matrix.rules` verbietet mechanisch jede Referenz
      `spec → adr`, real gegen `slice-sdk-csharp-pack-werkzeug` bereits
      einmal geprüft).
- [x] `spec/pflichtenheft.md` §6 Externe Verträge bekommt eine neue
      `SPEC-<NNN>`-Zeile für `pgchangefeed-kotlin` (System: GitHub-Packages-
      Gradle-/Maven-Package; Version: SemVer 2.0, `0.x.y`; Vertrag-Datei:
      Verweis auf `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` als
      Metadaten-Quelle, analog der bestehenden `SPEC-026`/`SPEC-027`-Zeilen)
      sowie §7 Historie-Zeile. Die reale nächste freie Nummer wird zum
      Bau-Zeitpunkt dieses Slice **neu verifiziert**
      (`grep -c "SPEC-0[0-9]*" spec/pflichtenheft.md` bzw. gleichwertig) —
      `SPEC-027` ist der zuletzt real vergebene Wert zum Zeitpunkt der
      Planung dieser Welle (2026-09-20), `SPEC-028` damit der
      voraussichtlich nächste freie, aber kein Vorgriff, falls zwischen
      Planung und Umsetzung ein anderer Zug bereits `SPEC-028` vergeben hat.
- [x] `harness/README.md` §Werkzeuge bekommt die reale
      `make sdk-pack-kotlin`-Zeile (kein Gate, Bindung auf
      [`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)) —
      erst jetzt zulässig, weil das Ziel jetzt real existiert
      (`AGENTS.md` §4).
- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: `harness/README.md` §Werkzeuge (siehe oben) — entfällt
      als eigener Punkt, da bereits oben als DoD-Kriterium geführt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`; keine Beobachtung angefallen
      ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-kotlin-lh-fa-sst-009](../welle-sdk-kotlin-lh-fa-sst-009.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/mk/sdk.mk` (bereits vorhanden aus der C#-/Python-Welle) | update | neues `sdk-pack-kotlin`-Target, Docker-only, kein `GATE_CHECKS`-Eintrag. |
| `sdks/kotlin/Dockerfile` | update | zusätzliche Bau-/Export-Stufe für den Jar-Bau, Extraktion analog `make proto-generate`/`make sdk-pack-csharp`. |
| `tools/harness/sdk-pack-kotlin.sh` (Plan-Nachzug — nicht explizit als eigene Zeile geplant, aber Sache der DoD/Ziel-Beschreibung, analog `sdk-pack-csharp.sh`/`sdk-pack-python.sh`) | neu | Host-seitiger Export-Wrapper (`docker build --build-context proto=proto --target pack-export`, `docker run --rm --network none … \| tar -x -C sdks/kotlin/dist/`). |
| `sdks/kotlin/.gitignore` (Plan-Nachzug) | update | `dist/`-Eintrag ergänzt, analog `sdks/csharp/.gitignore`/`sdks/python/.gitignore`. |
| `spec/pflichtenheft.md` §1 (`LH-FA-SST-009.a`) | update | Nachzug-Satz: Kotlin/GitHub Packages nicht mehr offen, vierte Sprache/Vertriebsweg bleibt offen. Zusätzlich (Plan-Nachzug, `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`): die stale gewordene Schlusszeile des Python-Absatzes („Eine dritte Sprache … bleibt offen") entfernt, damit sie dem neuen Kotlin-Absatz nicht widerspricht — der C#-Absatz-Altbestand („Eine zweite Sprache … bleibt offen", bereits vor diesem Slice stale) bleibt unangetastet, außerhalb des DoD-Umfangs. |
| `spec/pflichtenheft.md` §6 | update | neue `SPEC-028`-Zeile für `pgchangefeed-kotlin` (real verifiziert: `SPEC-027` letzte vergebene Nummer, kein `SPEC-028`-Treffer im Bestand vor diesem Zug). |
| `spec/pflichtenheft.md` §7 Historie | update | Historie-Zeile für `SPEC-028` und den `LH-FA-SST-009.a`-Nachzug. |
| `harness/README.md` §Werkzeuge | update | reale `make sdk-pack-kotlin`-Zeile. |

**Referenz für die Nummer-Verifikation:** `grep -n "SPEC-027" spec/pflichtenheft.md`
zeigte zum Zeitpunkt der Planung dieser Welle (2026-09-20) drei Treffer
(§2-Vertragszeile, §6-Zeile, §7-Historie), keine `SPEC-028`-Zeile irgendwo
im Baum — real gegen `spec/*.md` und `docs/` geprüft
(`grep -rn "SPEC-028" spec/*.md docs/` liefert 0 Treffer). Der
Implementer-Zug dieses Slice wiederholt diese Prüfung am eigenen
Bau-Zeitpunkt, statt den hier dokumentierten Stand unbesehen zu
übernehmen (`AGENTS.md` §3.12).

**Hinweise aus dem Beobachtungs-Register (proaktiv):**

- `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (19×): Der Nachzug-Satz bei
  `LH-FA-SST-009.a` macht seinerseits einen bereits bestehenden Satz
  falsch (wie schon zweimal bei C#/Python geschehen) — dieser Slice prüft
  vor dem Schreiben, dass der neue Satz an den bestehenden Wortlaut
  anschließt, statt neben ihm zu stehen (`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`,
  1×, dieselbe Vorsicht).
- **Backtick-Paritäts-Check** vor jedem Commit dieses Slice — besonders
  wichtig hier, weil die neue `SPEC-<NNN>`-Zeile und der Nachzug-Satz
  mehrere neue Kennungen und Backtick-Paare in kurzer Distanz einführen.

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-sdk-kotlin-http-client-flaeche`
**und** `slice-sdk-kotlin-grpc-client-flaeche` in `done/` liegen (siehe
Welle-Plan §4 Reihenfolge — `ADR-0109` Festlegung 1: ein Package aus beiden
Flächen).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): entfällt aus
  heutiger Sicht — ein Make-Target plus zwei Doku-Nachzüge.
- `in-progress` → `open` (blockiert — Carveout?): der Jar-Bau scheitert
  strukturell am gepinnten JDK-Image oder am Gradle-`maven-publish`-Setup —
  unwahrscheinlich, `./gradlew build` ist Standard-Tooling jeder
  Gradle-Version, `examples/kotlin/` belegt den JDK/Gradle-Bau bereits
  real.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + reales `.jar` als Smoke-Beleg +
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Der Export-Mechanismus für das `.jar` aus dem Docker-Bau ist für diesen
  Artefakt-Typ neu (die C#-/Python-Wellen belegen ihn bereits für
  `.nupkg`/`.whl`+`.tar.gz`, aber nicht für ein Gradle-Jar) — ein
  Bind-Mount-Workaround oder ein falsch gesetztes `--user` könnte das
  Artefakt mit falschen Dateirechten oder gar nicht auf den Host bringen.
  **Ausgang:** weiter offen, zu prüfen beim Schreiben — dasselbe
  `tar`-Stream-Export-Muster wie bei den beiden Vorgänger-SDKs sollte
  unverändert tragen, da es artefakt-unabhängig ist (ein `tar`-Stream
  kennt keinen Unterschied zwischen `.nupkg`/`.whl`/`.jar`).
- Ob GitHub Packages für dieses Package ein Sources-/Javadoc-Jar zwingend
  verlangt oder nicht, ist zum Planungszeitpunkt dieser Welle **nicht**
  abschließend recherchiert (`ADR-0109` lässt die Frage bewusst offen).
  **Ausgang:** weiter offen, zu klären beim Schreiben dieses Slice — real
  gegen die GitHub-Docs zu prüfen, mit Beleg im Bericht (`AGENTS.md`
  §3.12).
- Die neue `SPEC-<NNN>`-Nummer in `spec/pflichtenheft.md` §6/§2 muss
  fortlaufend vergeben werden (zum Planungszeitpunkt dieser Welle zuletzt
  `SPEC-027`) — ein Parallel-Slice, der ebenfalls eine neue `SPEC-*`-Nummer
  vergibt, könnte zu einer Kollision führen. **Ausgang:** weiter offen,
  real zu verifizieren zum Bau-Zeitpunkt dieses Slice (siehe §3
  Referenz-Prüfung).

## 7. Closure-Notiz

- **Was hat funktioniert:** <wird beim Abschluss ergänzt>
- **Was ging anders als geplant:** <wird beim Abschluss ergänzt>
- **Steering-Loop-Eintrag:** <wird beim Abschluss ergänzt>
- **Beobachtungs-Register (`../observations/`):** <wird beim Abschluss
  ergänzt>
- **Folge-Slices:** `slice-sdk-kotlin-publish-workflow` — bereits als
  Datei in `open/` vorhanden.
- **Risiken aus §6:** <wird beim Abschluss ergänzt>
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-kotlin-lh-fa-sst-009](../welle-sdk-kotlin-lh-fa-sst-009.md)
  (noch offen) — die Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Zwei Sub-Areas: `sdks/kotlin/`
(bereits eröffnet, GF) und `spec/pflichtenheft.md`/`harness/README.md`
(Default-Sub-Area `*`/`PGC`, Greenfield, `harness/conventions.md`).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/plan-nachzug`/`BEO-PGC/adr-folgepflicht-ohne-traeger-slice` geprüft
(Namen legen einen Bezug nahe): beide betreffen andere Fehlerklassen
(Slice-Plan-Nachzug bzw. ADR-Folgepflicht ohne jeden Träger-Slice) — dieser
Slice **ist** der Träger-Slice, kein Treffer, der ihn selbst beträfe.
`BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (1×) betrifft diesen Slice
über den in §3 aufgenommenen Hinweis.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
