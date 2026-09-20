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
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8);
      Report: `docs/reviews/review-slice-sdk-kotlin-pack-werkzeug.md`
      (0 HIGH, keine Fixrunde, DoD-Checkbox-Nachzug im selben Commit,
      Reviewer-Skill §DoD-Checkbox-Nachzug ohne Fixrunde).
- [x] Doku-Update: `harness/README.md` §Werkzeuge (siehe oben) — entfällt
      als eigener Punkt, da bereits oben als DoD-Kriterium geführt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — zwei
      Belege: F-1 an `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`
      (2×, weiter offen), F-2 an
      `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (bereits verkörpert,
      weiterer Beleg an einem neuen Beleg-Typ: Web-Recherche) — siehe §7.
- [x] Jedes Risiko aus §6 trägt einen Ausgang — siehe §7.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-kotlin-lh-fa-sst-009](../welle-sdk-kotlin-lh-fa-sst-009.md)
      (noch offen, nur noch `slice-sdk-kotlin-publish-workflow` fehlt) —
      Anker gesetzt (§0 Bezug,
      [`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)),
      Folge-Slice benannt (§7), Register
      fortgeschrieben (siehe oben); die volle Drei-Paarungen-Prüfung der
      Welle selbst läuft regelkonform erst bei deren eigener Closure.

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
  **Ausgang: eingetreten wie erwartet, kein Problem.** Der bestehende
  `tar`-Stream-Export (`pack-export`-Stufe) trug unverändert für ein Jar —
  vier unabhängige Läufe (Implementer, Reviewer, Verifier ×2) erzeugten je
  ein reales, valides `sdks/kotlin/dist/pgchangefeed-kotlin-0.1.0.jar`
  (112591 Bytes, normale Owner-Rechte, kein `--user`-Workaround nötig).
- Ob GitHub Packages für dieses Package ein Sources-/Javadoc-Jar zwingend
  verlangt oder nicht, ist zum Planungszeitpunkt dieser Welle **nicht**
  abschließend recherchiert (`ADR-0109` lässt die Frage bewusst offen).
  **Ausgang: real recherchiert und entschieden — ein Haupt-Jar genügt.**
  Drei unabhängige Web-Recherchen (Implementer, Reviewer, Verifier) zeigen:
  die zitierte GitHub-Doku-Seite erwähnt Sources-/Javadoc-Jars an keiner
  Stelle (weder Pflicht noch Freistellung), trägt die im Diff verwendete
  Zuschreibung also nicht explizit — der praktische Schluss bleibt aber
  haltbar, gestützt auf das real bestätigte Gradle-Kern-Faktum
  (`from(components["java"])` ohne `withSourcesJar()`/`withJavadocJar()`
  packt nur den Haupt-Jar). Die Attributions-Unschärfe ist Review-Finding
  F-2 (MEDIUM) und dem Beobachtungs-Register zugeführt (§7); Risiko
  zusätzlich strukturell durch `ADR-0109` §Re-Evaluierungs-Trigger 4
  abgefangen.
- Die neue `SPEC-<NNN>`-Nummer in `spec/pflichtenheft.md` §6/§2 muss
  fortlaufend vergeben werden (zum Planungszeitpunkt dieser Welle zuletzt
  `SPEC-027`) — ein Parallel-Slice, der ebenfalls eine neue `SPEC-*`-Nummer
  vergibt, könnte zu einer Kollision führen. **Ausgang: keine Kollision.**
  `SPEC-028` real vergeben und dreifach verifiziert (Implementer, Reviewer,
  Verifier je per `grep -rn "SPEC-028"`) — kein Parallel-Zug hat die Nummer
  belegt, genau die drei erwarteten Treffer (§1, §6, §7).

## 7. Closure-Notiz

- **Was hat funktioniert:** Das etablierte Dreiklang-Muster
  (`harness/mk/sdk.mk`-Target, Dockerfile-`pack`/`pack-export`-Stufen,
  host-seitiger `tar`-Stream-Export-Wrapper) trug unverändert für einen
  dritten, neuen Artefakt-Typ (`.jar` statt `.nupkg`/`.whl`+`.tar.gz`) —
  der Reviewer bestätigt `tools/harness/sdk-pack-kotlin.sh` als „strukturell
  nahezu Zeile-für-Zeile identisch" mit `sdk-pack-csharp.sh`. Der
  Träger-Nachzug (Pflichtenheft `LH-FA-SST-009.a`/§6/§7,
  `harness/README.md`) folgte demselben, bereits zweimal geübten Muster.
  Alle drei §6-Risiken lösten sich real wie im Plan erwartet auf (siehe
  unten) — keine Überraschung im Bau selbst.
- **Was ging anders als geplant:** Der Kotlin-Nachzug bereinigte korrekt
  den überholten Python-Absatz in `LH-FA-SST-009.a` (kein Wiederholungsfall
  von `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` innerhalb dieses
  eigenen Zuges), ließ dabei aber sichtbar werden, dass der **C#-Absatz**
  bereits seit dem vorangegangenen Python-Nachzug denselben Fehler trägt —
  unbemerkt in beiden vorangegangenen Reviews. Zusätzlich fand der
  Reviewer, dass die im Diff verwendete Sources-/Javadoc-Jar-Freistellungs-
  Aussage ihren zitierten Beleg (die GitHub-Doku-Seite) nicht explizit
  trägt — die Seite schweigt zur Frage, statt sie zu beantworten; der
  praktische Schluss bleibt trotzdem haltbar (drittfach unabhängig
  nachrecherchiert: Implementer, Reviewer, Verifier).
- **Steering-Loop-Eintrag:** Zwei Finding-Klassen dieses Laufs, beide nicht
  neu geschaffen, sondern an bestehenden Beobachtungen angedockt (siehe
  unten) — kein neuer BEO-Eintrag nötig, die bestehende Klassifikation
  greift für beide Fälle, jeweils an einem für die Klasse neuen Träger-Typ
  (Pflichtenheft-Fließtext bzw. Web-Recherche statt Befehl/Codestelle).
- **Beobachtungs-Register (`../observations/`):**
  - F-1 (LOW, C#-Absatz bleibt stale) →
    `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`, neuer Beleg
    (`evidence/slice-sdk-kotlin-pack-werkzeug.md`), Zähler 1× → 2×, weiter
    unter der Schwelle (3×), Zustand bleibt **offen**. Ursprung und
    Vorkommen sind getrennt: Ursprung ist `slice-sdk-python-pack-werkzeug`
    (dort entstand die Staleness), gefunden wurde sie erst beim Review
    dieses Slice — derselbe Ursprung/Vorkommen-Split, den
    `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` bereits einmal
    dokumentiert. Erster Beleg dieser Klasse an einem
    Pflichtenheft-Fließtext-Träger statt einem Slice-Plan-Dokument.
  - F-2 (MEDIUM, Beleg trägt die Aussage nur durch Schweigen) →
    `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (bereits `verkörpert`
    seit `welle-20`, Reviewer-Skill-HIGH-Punkt „Beleg trägt seinen Satz
    nicht"), neuer Beleg (`evidence/slice-sdk-kotlin-pack-werkzeug.md`) —
    erster Beleg dieser Klasse an einer **Web-Recherche** statt einem
    Befehl/einer Codestelle/einer Adresse/einer Assertion. Kein neuer
    Ausgang nötig (Klasse bereits verkörpert); state.md bleibt unverändert
    (Präzedens: `evidence/slice-release-version-und-workflow.md` wurde
    ebenfalls ohne state.md-Änderung nach der Verkörperung ergänzt).
  - Keine dritte Beobachtung angefallen.
- **Folge-Slices:** `slice-sdk-kotlin-publish-workflow` — bereits als
  Datei in `open/` vorhanden; um einen proaktiven Hinweis auf F-2 ergänzt
  (§3 „Hinweise aus dem Beobachtungs-Register", optional, kein eigener
  DoD-Punkt dort).
- **Risiken aus §6:** alle drei real aufgelöst — Export-Mechanismus trug
  unverändert (reales `.jar`, vier unabhängige Läufe), Sources-/
  Javadoc-Jar-Frage real recherchiert und entschieden (ein Haupt-Jar
  genügt, Attributions-Unschärfe als F-2 dem Register zugeführt),
  `SPEC-028` kollisionsfrei vergeben (dreifach verifiziert). Kein Risiko
  bleibt strukturell offen für diesen Slice selbst — anders als bei
  `slice-sdk-kotlin-publish-workflow`, wo `AGENTS.md` §3.10 bis zum realen
  Post-Push-Lauf weiter offen bleibt.
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-kotlin-lh-fa-sst-009](../welle-sdk-kotlin-lh-fa-sst-009.md)
  (noch offen — dies ist der **vorletzte** der fünf Slices; nur noch
  `slice-sdk-kotlin-publish-workflow` fehlt für die Wellen-Closure). Anker
  (`ADR-0109`, `LH-FA-SST-009`), Folge-Slice
  (`slice-sdk-kotlin-publish-workflow`, bereits in `open/`) und Register
  (siehe oben) sind für diesen Slice selbst getragen; die volle
  Drei-Paarungen-Prüfung **der Welle** läuft regelkonform erst bei deren
  eigener, separater Closure — nicht Teil dieses Slice-Zugs.

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
