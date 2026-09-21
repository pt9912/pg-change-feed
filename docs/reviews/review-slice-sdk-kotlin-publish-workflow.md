# Review-Report: slice-sdk-kotlin-publish-workflow — 2026-09-21

**Review-Art:** Code — geprüft gegen Plan + `ADR-0109` + Hard Rules
(`AGENTS.md` §3, Modul 10 §Drei Review-Arten).

**Gegenstand:** Commit `546a7503eb8615b1a2479c60b89e915ca1d112ab`
(`feat(sdk-kotlin): GitHub-Packages-Publish-Workflow sdk-kotlin-release.yml`)
gegen Elternstand `eaf6ded5` — `git diff eaf6ded5..546a7503`.

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-09-13 (Handbuch-
Versionshistorie-HIGH nachgetragen; Repo-Commit `HEAD` zum Zeitpunkt dieses
Laufs).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-21.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-sdk-kotlin-publish-workflow.md`
  (vollständig gelesen, §1–§8)
- `ADR-0109` (Kotlin/GitHub Packages als drittes SDK-Package,
  `Accepted`, vollständig gelesen)
- `ADR-0051` Entscheidung 3/8 (Release-Tag-Trigger, Secret-Muster)
- `ADR-0106`/`ADR-0107` (Vorgänger-SDK-Publish-Workflows — Formvorbild und
  Kontrastfolie)
- `LH-FA-SST-009` (`spec/lastenheft.md`)
- `AGENTS.md` §3.1 (Docker-only), §3.5 (Accepted-ADR-Immutabilität), §3.7
  (Kommentar-Disziplin), §3.8 (Action-Pinning), §3.9 (Exit-Code-Disziplin),
  §3.10 (Post-Push-Risiko), §3.12/§3.13 (Beleg-Herkunft, Träger-Nachzug)
- `harness/conventions.md` (MR-000/MR-001/MR-002)
- Vergleichsdateien: `.github/workflows/sdk-csharp-release.yml`,
  `.github/workflows/sdk-python-release.yml`,
  `tools/harness/sdk-csharp-release-tag-info.sh`,
  `tools/harness/sdk-python-release-tag-info.sh`,
  `tools/harness/semver-regex.sh`, `docs/user/releasing.md`
- Vorherige Findings am selben Modul: `docs/reviews/review-slice-sdk-kotlin-pack-werkzeug.md`,
  `docs/reviews/review-slice-sdk-csharp-publish-workflow.md`,
  `docs/reviews/review-slice-sdk-python-publish-workflow.md`

---

## Findings

### F-1 — `./gradlew publish` läuft direkt auf dem Runner statt, wie von `ADR-0109` Festlegung 5 wörtlich festgelegt, im gepinnten Docker-Image — echte Rebuild-Duplikation, nicht bloß eine CLI-Zeile

- `kategorie`: HIGH
- `quelle`: `ADR-0109` §Entscheidung Festlegung 5 (`Accepted`), `AGENTS.md`
  §3.1 (Docker-only), §3.5 (Accepted-ADR-Immutabilität)
- `pfad`: `.github/workflows/sdk-kotlin-release.yml:111–124` (Schritte
  „Kotlin-SDK bauen, testen, paketieren (make sdk-pack-kotlin)" und „Nach
  GitHub Packages veroeffentlichen (./gradlew publish)"),
  `docs/plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md:424–430`,
  `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts:26–29` (Kommentar
  bestätigt die Abweichung selbst)
- `befund`: `ADR-0109` §Entscheidung Festlegung 5 legt wörtlich fest: „ein
  neues, netzlos nicht prüfbares Werkzeug-Ziel, Arbeitsname
  `make sdk-pack-kotlin`, baut/testet/paketiert (`./gradlew build`,
  `./gradlew publish`) **im gepinnten** `eclipse-temurin:21-jdk`**-Image**."
  Real umgesetzt ist stattdessen: `make sdk-pack-kotlin` ruft nur
  `./gradlew test`/`build` im Docker-Image auf (`sdks/kotlin/Dockerfile:65-66`,
  bereits in `slice-sdk-kotlin-pack-werkzeug` so gebaut und dort explizit im
  Dockerfile-Kommentar begründet: „Kein `publish`-Aufruf in diesem Slice");
  `./gradlew --no-daemon publish` läuft in diesem Diff als **eigener,
  separater Schritt direkt auf dem GitHub-hosted Runner** — außerhalb jedes
  Docker-Containers. Der neue Kommentar in `build.gradle.kts:26-29` bestätigt
  das selbst ausdrücklich: „Außerhalb dieses Workflows (lokal, in
  `make sdk-pack-kotlin`) bleiben beide Umgebungsvariablen leer —
  `./gradlew publish` läuft dort ohnehin nicht." Das ist keine
  Formalie: `publish` hängt (über die Standard-`maven-publish`-Task-Graphen
  des `from(components["java"])`-Publications-Blocks) von `jar` ab, `jar`
  von `compileKotlin`, `compileKotlin` vom `generateProto`-Task des
  `com.google.protobuf`-Plugins — dieselbe Kette, die `make sdk-pack-kotlin`
  bereits vollständig **im Docker-Image** durchläuft. Der Runner-Schritt
  wiederholt diese gesamte Bau-Kette ein zweites Mal, außerhalb von Docker,
  mit einem frischen, ungecachten Gradle-Dependency-Bezug (Kotlin-Gradle-Plugin
  2.4.20, `com.google.protobuf` 0.10.0, `protoc` 4.36.2, mehrere
  `io.grpc:*`-Koordinaten) direkt vom Runner aus über das offene Netz — und
  das bereits von `make sdk-pack-kotlin` erzeugte, im selben Workflow einen
  Schritt zuvor produzierte `sdks/kotlin/dist/pgchangefeed-kotlin-0.1.0.jar`
  wird dabei **nicht** verwendet; Gradle baut auf dem Runner sein eigenes
  `build/libs/*.jar` neu. Das unterscheidet sich strukturell von
  `sdk-csharp-release.yml`s `dotnet nuget push sdks/csharp/dist/*.nupkg`
  und `sdk-python-release.yml`s `uv publish sdks/python/dist/*`
  (real gegengelesen, siehe Eingangs-Kontext) — beide laden ausschließlich
  das bereits im Docker-Bau erzeugte, unveränderte Artefakt aus `dist/`
  hoch, ohne irgendeinen Compile-/Codegenerierungs-Schritt auf dem Runner
  auszulösen. Der Kotlin-Workflow-Kommentar rechtfertigt die Abweichung mit
  einer eigenen Lesart von `AGENTS.md` §3.1 („bindet die
  Host-Toolchain-Sperre an lokale Entwicklung, nicht an den ephemeren
  GitHub-hosted Runner") — das ist eine plausible Position, aber sie steht
  im Widerspruch zu dem, was `ADR-0109` als **Accepted**-Entscheidung bereits
  wörtlich festgelegt hat (Docker-only bis einschließlich `publish`), ohne
  dass eine Folge-ADR mit `Supersedes ADR-0109` diese Änderung trägt
  (`AGENTS.md` §3.5). Weder der Slice-Plan noch die Commit-Message benennen
  diesen Wortlaut-Widerspruch explizit — der Plan spricht nur von einer
  „notwendigen Konkretisierung beim Schreiben", nicht von einer Abweichung
  von einem bereits im ADR-Text festgelegten Ort.
- `verifizierbar`: teilweise — die Textabweichung ist durch Gegenlesen von
  `ADR-0109` Zeile 424–430 gegen `sdks/kotlin/Dockerfile:20-24` und
  `build.gradle.kts:26-29` mechanisch nachvollziehbar (hier real getan); ob
  der reale Runner-seitige `./gradlew publish`-Lauf tatsächlich erfolgreich
  gegen GitHub Packages schreibt, bleibt strukturell erst durch den ersten
  echten Tag-Push geprüft (`AGENTS.md` §3.10, im Slice-Plan §6 bereits als
  offenes Risiko geführt — dieses Finding ist davon unabhängig: es betrifft
  den **Ort** der Ausführung, nicht ihren Erfolg).
- `klasse`: „Accepted-ADR-Festlegung beim Implementieren stillschweigend
  umgangen, ohne Supersedes-ADR"

## Negativbefunde

- geprüft, ohne Befund: Workflow-Trigger — `push: tags: ['sdk-kotlin-v*']`
  ist der einzige Trigger; `ci.yml`/`e2e.yml`/`examples.yml` tragen je
  `tags-ignore: ['**']`, `release.yml`/`sdk-csharp-release.yml`/
  `sdk-python-release.yml` tragen je ihr eigenes `tags: [...]`-Muster
  (`v*`/`sdk-csharp-v*`/`sdk-python-v*`) — real per `grep -n "tags" -A2`
  über alle sechs Dateien gegengelesen, kein Muster überlappt
  `sdk-kotlin-v*`.
- geprüft, ohne Befund: `tools/harness/sdk-kotlin-release-tag-info.sh` —
  teilt sich `tools/harness/semver-regex.sh` als dritter Konsument (Kopf-
  Kommentar dort korrekt auf drei Konsumenten erweitert), gleicht gegen die
  Top-Level-`version`-Zeile in `build.gradle.kts` ab, bricht **vor**
  `make sdk-pack-kotlin` ab (Schrittreihenfolge im Workflow: Tag-Validierung
  → `build.gradle.kts`-Abgleich → `make sdk-pack-kotlin` → Publish). Eigener
  Lauf: `make test-sdk-kotlin-release-tag-info` → „alle Fälle bestanden"
  (16 Fälle, 6 gültig/10 ungültig, inkl. Präfix-Kollisionsfälle
  `sdk-csharp-v*`/`sdk-python-v*`).
- geprüft, ohne Befund: `permissions` — Top-Level `permissions: {}`
  (alles entzogen), Job-Ebene `contents: read`/`packages: write`, kein
  `secrets.<NAME>`-Verweis außer `secrets.GITHUB_TOKEN` (dem eingebauten
  Token) — kein externes Repository-Secret, wie geplant.
- geprüft, ohne Befund: `build.gradle.kts`-`publishing.repositories.maven`-
  Block — Registry-URL `https://maven.pkg.github.com/pt9912/pg-change-feed`,
  Credentials aus `System.getenv("GITHUB_ACTOR")`/
  `System.getenv("GITHUB_TOKEN")`; eigener Lauf von `make sdk-pack-kotlin`
  bestätigt, dass Gradle die Datei inklusive des neuen `publishing`-Blocks
  fehlerfrei parst und den bereits vorher erzeugten Jar-Bau-Cache-Layer
  unverändert bestätigt (`sdks/kotlin/dist/pgchangefeed-kotlin-0.1.0.jar`,
  Zeitstempel identisch mit dem Implementer-Lauf, `git status` danach
  sauber).
- geprüft, ohne Befund: Action-Pinning — die einzige `uses:`-Zeile
  (`actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1`)
  real gegen `api.github.com/repos/actions/checkout/git/refs/tags/v7.0.1`
  reverifiziert — SHA identisch, Tag-Kommentar korrekt.
- geprüft, ohne Befund: `harness/README.md` §Werkzeuge — neue Zeile für
  `sdk-kotlin-release.yml`, thematisch/alphabetisch zwischen den
  Nachbar-SDK-Workflow-Zeilen positioniert, Bindung auf `ADR-0109`
  Festlegung 2/5, Format konsistent mit den Nachbarzeilen.
- geprüft, ohne Befund: `docs/user/releasing.md` — neuer Abschnitt „SDK-
  Release: GitHub-Packages-Publish für `pgchangefeed-kotlin`", **ohne**
  neue Secret-Tabellenzeile, mit explizitem Hinweis auf die
  Lese-Authentifizierungspflicht (PAT, `read:packages`); §1-Intro korrekt
  auf „vierter … Weg … implementiert, aber … noch unbewiesen" gezogen;
  Versionshistorie korrekt von 1.4 auf 1.5 mit eigener, inhaltlich
  zutreffender Zeile fortgeschrieben (`AGENTS.md` §3.7-Zustandsfeld-Analogon
  eingehalten — Zeile nennt Zustand + Anker, keine Chronik-Prosa über den
  Diff hinaus).
- geprüft, ohne Befund: YAML-Struktur — `ruby -ryaml -e YAML.load_file`
  parst `.github/workflows/sdk-kotlin-release.yml` fehlerfrei; jeder der
  fünf `run:`-Blöcke einzeln extrahiert und mit `bash -n` geprüft — alle
  fünf syntaktisch fehlerfrei.
- geprüft, ohne Befund: Träger-Nachzug (`AGENTS.md` §3.13) — eigener
  Suchlauf über `grep -rln "sdk-csharp-v\|sdk-python-v\|NUGET_API_KEY\|
  PYPI_API_TOKEN\|SDK-Release"` und gezielt über
  `spec/pflichtenheft.md` (`LH-FA-SST-009.a`, `SPEC-028`, §6 Externe
  Verträge): alle Stellen außerhalb des Diffs, die die bewegte Eigenschaft
  „Anzahl/Zustand der SDK-Release-Wege" beschreiben, waren bereits vor
  diesem Diff korrekt (durch den vorangehenden `slice-sdk-kotlin-pack-
  werkzeug`-Zug nachgezogen) — kein durch diesen Diff neu falsch gewordener
  Träger gefunden, außer den im Diff selbst bereits mitgezogenen
  (`docs/user/releasing.md`).
- geprüft, ohne Befund: Backtick-Parität — alle vier von diesem Diff
  inhaltlich berührten Markdown-/Konfigurationsdateien (Slice-Plan 474,
  `docs/user/releasing.md` 470, `harness/README.md` 1782, Workflow-Datei
  128 Backticks) sind je `% 2 == 0`.
- geprüft, ohne Befund: Out-of-Scope-Einhaltung — `git diff --stat`
  zeigt exakt neun Dateien, alle innerhalb des geplanten Umfangs; kein
  echter Tag wurde gesetzt (`git tag -l | grep kotlin` liefert keinen
  Treffer); keine der vier bereits geschlossenen Vorgänger-Slice-Dateien
  (`slice-sdk-kotlin-projektgeruest`/`-pack-werkzeug`/
  `-http-client-flaeche`/`-grpc-client-flaeche`) ist im Diff enthalten;
  kein `docs/plan/adr/`-Diff.
- geprüft, ohne Befund: `make gates` — eigener Lauf, Exit-Code direkt
  (ungepiped) geprüft: `0` (grün), inklusive `generated-sync`/`a-check`
  (`gesamt: 0 Befund(e)`), `AGENTS.md` §3.9 eingehalten.
- geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7) — die
  Slice-ID-Nennung im Workflow-Kopfkommentar
  (`.github/workflows/sdk-kotlin-release.yml:6`) folgt dem bereits in
  `review-slice-sdk-kotlin-pack-werkzeug.md` als „etablierte
  Herkunfts-Anker-Form" akzeptierten Muster, kein Konjunktiv über eine
  verworfene Alternative, kein abwesender Text, kein mitten im Satz
  abgebrochener Kommentar.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Accepted-ADR-Festlegung beim
Implementieren stillschweigend umgangen, ohne Supersedes-ADR

## Verdikt

**Merge-blockierend:** ja — ein HIGH-Finding (F-1). Es handelt sich um
einen realen, mechanisch nachvollziehbaren Widerspruch zwischen dem
Wortlaut einer `Accepted`-ADR (`ADR-0109` §Entscheidung Festlegung 5:
Docker-only bis einschließlich `publish`) und der tatsächlichen
Implementierung (`./gradlew publish` direkt auf dem Runner, mit echter
Rebuild-/Codegenerierungs-Duplikation außerhalb von Docker) — kein
Stilproblem, sondern ein Hard-Rule-naher ADR-/Docker-only-Konflikt
(`AGENTS.md` §3.1/§3.5). Da es sich um ein HIGH mit potenziellem
Rollen-Widerspruch handelt (der Implementer hat die Abweichung im
Slice-Plan selbst als bewusste, begründete Entscheidung dargestellt, nicht
als offene Frage), empfiehlt dieser Report **nicht** die einseitige
Rückgabe an den Implementer als alleinige Auflösung: Der Diskussionspunkt
sollte — wie vom Koordinator bereits vorgesehen — über eine Fixrunde
laufen, die entweder (a) `./gradlew publish` tatsächlich in den
Docker-Bau zieht (dem ADR-Wortlaut folgend) oder (b) eine begründete
Klarstellung/Korrektur zu `ADR-0109` nachträgt (`Supersedes ADR-0109`,
`AGENTS.md` §3.5), bevor dieser Slice nach `done/` geht. Wegen des HIGH
bleibt die DoD-Checkbox „Review durchgeführt" in
`slice-sdk-kotlin-publish-workflow.md` bewusst **offen** — sie wird erst
nach einer erfolgreichen Fixrunde vom Implementer nachgezogen (Reviewer-Skill
§DoD-Checkbox-Nachzug ohne Fixrunde greift hier explizit **nicht**, da eine
Fixrunde erforderlich ist).

**Übergabe:** Finding F-1 geht an den Implementer zurück (Fixrunde nötig,
kein Self-Review durch diesen Reviewer-Lauf); bei Bedarf eskaliert der
Koordinator zusätzlich an den Architect (Modul 8), falls die Auflösung eine
ADR-Korrektur statt eines reinen Workflow-Fixes verlangt. Dieser Report
selbst ist ein Lauf-Beleg und ersetzt keine Verifikation (Modul 11, separater
Kontext).
