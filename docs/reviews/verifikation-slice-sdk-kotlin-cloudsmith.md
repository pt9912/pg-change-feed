# Verifikations-Report: slice-sdk-kotlin-cloudsmith — 2026-09-26

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates + Nachmessung der Fixrunde + Nachmessung des realen Post-Push-Laufs
([`AGENTS.md`](../../AGENTS.md) §3.10). Review-Artefakt des Reviewers:
[`review-slice-sdk-kotlin-cloudsmith.md`](review-slice-sdk-kotlin-cloudsmith.md); Formvorbild dieses
Reports: [`verifikation-slice-sdk-readme-nutzerdoku.md`](verifikation-slice-sdk-readme-nutzerdoku.md). Die
vendored Baseline trägt kein eigenes Verifikations-Template (`.harness/baseline/v6.9.0/templates/docs/reviews/`
enthält nur `review-report.template.md`), ein Skill `.harness/skills/verifier.md` liegt nicht vor; der Report
folgt dem Formvorbild.

**Gegenstand:** Slice-Plan `slice-sdk-kotlin-cloudsmith` (ohne Welle), Stand `HEAD` = `2e9d8db8` (gepusht,
Tag `sdk-kotlin-v0.2.2` zeigt auf denselben Commit), Diff-Range `dbb87155..HEAD`: 9 Commits, 12 Dateien
(+686/−229 laut `git diff --stat`). Slice-Inhalt sind `d11adc14`, `98c4437e`, `9b045c95` (Lifecycle,
Verantwortlicher), `4d91d9ef` (Build, Workflow, Probe), `e5070e1b` (Träger), `6f91e197`, `577876ac` (Plan)
und die Fixrunde `2e9d8db8`; nicht Slice-Inhalt ist der Review-Report `a8e765dd`. Bezug:
[`LH-FA-SST-009`](../../spec/lastenheft.md),
[`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md),
[`ADR-0109`](../plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md),
[`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md).
Dieser Lauf ändert weder Code noch Plan, Spec oder Doku; er schreibt nur diesen Report. Alle Mutationen
liefen im Arbeitsbaum (reine Bash-Ersetzung, kein `sed -i`, kein Host-Python) und sind je Lauf per
`git checkout` zurückgenommen (`git status --short` leer nach jedem Lauf); nichts gepusht, nichts getaggt,
kein Publish, kein Zugriff mit Zugangsdaten. Zugriffe auf Cloudsmith und GitHub waren lesend und anonym
(`curl`) bzw. über `gh run view`/`gh api`.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei; der Exit-Code wurde im selben Aufruf gesondert gesichert, das Log danach
gelesen. Stand aller unmutierten Läufe: `HEAD` = `2e9d8db8`, Arbeitsbaum sauber. Schwere Läufe liefen nacheinander;
`free -m` (verfügbar) vor den Läufen 17,2 bis 17,5 GB.

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` | **EXIT=0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `coverage-gate: OK — Coverage 83.40% erfüllt Schwelle 80%` · `d-check: 1182 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"` · `generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto-export)` · a-check `gesamt: 0 Befund(e)` |
| `make docs-check` | **EXIT=0** | `d-check: 1182 Datei(en) geprüft, 0 Befund(e)` (die Zahl liegt um eins über dem Plan-Stand 1181, weil der Review-Report seither existiert) |
| `make commit-traceability` / `make doc-commits` / `make doc-immutable`, je `RANGE=dbb87155..HEAD` | je **EXIT=0** | `commit-traceability: OK — 9 Commit(s) in "dbb87155..HEAD", Betreffs ohne Struktur-ID`; je `d-check: 1182 Datei(en) geprüft, 0 Befund(e)`. `make doc-immutable` ohne `RANGE` bricht mit `flag needs an argument: --range` ab (Bedienung, kein Befund) |
| `make sdk-public-doc-check` | **EXIT=0** | `sdk-public-doc-check: keine interne Kennung unter sdks` |
| `make test-sdk-kotlin-release-tag-info` | **EXIT=0** | `run-sdk-kotlin-release-tag-info-tests: alle Fälle bestanden` |
| `make sdk-pack-kotlin` (unmutiert) | **EXIT=0** | alle Docker-Schichten einschließlich der Probe-Schicht `#20` `CACHED` (identischer Bau-Kontext wie im Implementer- und im Review-Lauf, die Testzeilen stammen nicht aus diesem Lauf); `sdks/kotlin/dist` trägt `pgchangefeed-kotlin-0.2.2.jar` (144457 Bytes) und `pgchangefeed-kotlin-0.2.2-sources.jar` (39041 Bytes); die Probe wurde stattdessen durch die Mutationen in §4 real ausgeführt |

Nicht gefahren: `make test-sdk-*-integration`, ein Publish, ein Tag-Push, jeder Zugriff mit Zugangsdaten. Hygiene:
dangling-Volumes (`docker volume ls -q -f dangling=true | wc -l`) vor den Läufen **34**, nach allen Läufen **34**;
kein `prune`; keine eigenen Volumes oder Images zurückgeblieben (das Image `pg-change-feed:sdk-kotlin-pack-export`
ist das Erzeugnis des Make-Ziels selbst).

## 2. Post-Push-Beleg nachgemessen (Plan §2/§5, [`AGENTS.md`](../../AGENTS.md) §3.10)

Alle Abrufe dieser Tabelle sind von mir gefahren; die Zahlen stammen aus meinen Ausgaben, nicht aus der Zusage.

| Behauptung | Meine Messung (gedruckt) | Verdikt |
|---|---|---|
| Tag `sdk-kotlin-v0.2.2` existiert auf `origin` und zeigt auf `HEAD` | `git ls-remote --tags origin 'sdk-kotlin*'`: `sdk-kotlin-v0.2.0` auf `eb3ba91e`, `sdk-kotlin-v0.2.1` auf `7df775c8`, `sdk-kotlin-v0.2.2` auf `2e9d8db89595f6fc228a83502315d93402daf200` (gleich `git log -1`) | bestätigt |
| Lauf `36201941235`: `success`, beide Jobs `success` | `gh run view 36201941235 --json conclusion,status,headSha,event,jobs`: `conclusion: success`, `status: completed`, `event: push`, `headBranch: sdk-kotlin-v0.2.2`, `headSha` gleich dem Tag-Commit; Job „…nach GitHub Packages veroeffentlichen“ `success` (23:40:47Z bis 23:43:16Z), Job „…nach Cloudsmith veroeffentlichen“ `success` (23:40:47Z bis 23:43:00Z); je alle Schritte `success`, darunter „Nach GitHub Packages veroeffentlichen (Gradle-Aufgabe des Ziels, im Docker-Image)“ und „Nach Cloudsmith veroeffentlichen (Gradle-Aufgabe des Ziels, im Docker-Image)“ | bestätigt |
| Gradle-Aufgabe des Ziels lief | Log: `> Task :publishMavenPublicationToCloudsmithRepository` (23:42:45 bis 23:42:58, `BUILD SUCCESSFUL in 20s`) bzw. `> Task :publishMavenPublicationToGitHubPackagesRepository` (23:42:41 bis 23:43:11, `BUILD SUCCESSFUL in 40s`); kein `Task … not found`, kein `401`/`403` im Log | bestätigt |
| Anonymer Abruf am README-Pfad | `curl -s -o … -w '%{http_code} %{size_download}' -L` gegen die Download-Basis aus dem `repositories`-Snippet der Kotlin-README plus Maven-Layout `io/github/pt9912/pgchangefeed-kotlin/0.2.2/…`: `.pom` **200 3095**, `.jar` **200 144457**, `-sources.jar` **200 39041**, `.module` **200 4711**, `maven-metadata.xml` **200 342** | bestätigt (Zeitpunkt meiner Abrufe 2026-09-26, nach der Indexierung) |
| `maven-metadata.xml` | `<latest>0.2.2</latest>`, `<release>0.2.2</release>`, genau eine `<version>0.2.2</version>`, `lastUpdated 20260925234427` | bestätigt |
| POM-Koordinate und Metadaten | die heruntergeladene POM trägt `groupId io.github.pt9912`, `artifactId pgchangefeed-kotlin`, `version 0.2.2`, `name`, `description`, `url https://github.com/pt9912/pg-change-feed`, Lizenz `MIT`, `scm`; Abhängigkeiten `kotlin-stdlib` `compile`, acht weitere `runtime`, dazu `grpc-bom` als `import` | bestätigt |
| Moduldatei | die `.module` ist ein gültiges Gradle-Metadatenmodell (`formatVersion 1.1`, `module pgchangefeed-kotlin`, `version 0.2.2`, Varianten `apiElements`/`runtimeElements`/`sourcesElements`); Cloudsmith hat sie angenommen und liefert sie aus | bestätigt |
| Artefakte gleich dem Bau | Größen von Jar (144457) und Sources-Jar (39041) sind **gleich** denen meines lokalen `make sdk-pack-kotlin`-Erzeugnisses; die SHA-256-Summen sind **verschieden** (ein Jar trägt Bauzeitstempel), Byte-Gleichheit behaupte ich deshalb nicht | Größe bestätigt, Byte-Gleichheit nicht geprüft |
| GitHub Packages listet die Versionen | `gh api /users/pt9912/packages/maven/io.github.pt9912.pgchangefeed-kotlin/versions --jq '.[].name'`: `0.2.2`, `0.2.1`, `0.2.0` | bestätigt |
| Zeit bis zur Auslieferung | Der Publish-Task selbst lief 13 s, der anonyme Abruf ist nach der Indexierung 200; die Zeit dazwischen liegt bei rund zwei Minuten (Angabe des Auftrags, `lastUpdated` 23:44:27Z gegen Task-Ende 23:42:58Z gibt rund 90 Sekunden als Untergrenze der Verarbeitung an, abgeleitet) | übernommen, plausibel |
| Weitere Läufe zum Tag-Commit | `gh run list --commit 2e9d8db8`: `sdk-kotlin-release` (push, completed, success), `ci` success, `examples` success, `e2e` success | bestätigt |
| Base-URL selbst | `curl` gegen `https://dl.cloudsmith.io/public/pt9912/pg-change-feed/maven/` liefert **404** (ein Verzeichnis ohne Index, kein Artefakt); Gradle löst Artefaktpfade auf, nicht die Basis | kein Befund, für Anwender ohne Folge |

### Secret-Hygiene im Lauf-Log

`gh run view 36201941235 --log` (1006 Zeilen, in einer Log-Datei gesichert, Exit 0) durchsucht:

- `env:`-Blöcke: `CLOUDSMITH_USERNAME: ***`, `CLOUDSMITH_API_KEY: ***`, `GITHUB_TOKEN: ***` — alle drei maskiert.
- Die `docker run`-Zeilen der Schritte tragen nur Namen: `-e CLOUDSMITH_USERNAME`, `-e CLOUDSMITH_API_KEY` bzw.
  `-e GITHUB_ACTOR`, `-e GITHUB_TOKEN`, kein Wert, keine Zuweisung.
- `CLOUDSMITH_API_KEY=` mit Klartext: **0** Treffer; `Authorization:`-Header: **0**; `Basic `: **0**; `password`: **0**
  (die einzigen `token`-Treffer sind maskierte Checkout-Zeilen, `git-credentials` des Checkout-Schritts und die
  Registry-Auth-Zeile des Docker-Bau-Schritts, ohne Wert).
- Das vierstellige Suffix und der Name des Service-Slugs stehen **nicht** im Log (0 Treffer); gesucht habe ich
  sie über die Angabe des Auftrags, ohne sie hier zu wiederholen.
- Kein `set -x`, kein `printenv`, kein `--info`/`--debug` im Gradle-Aufruf.

Ergebnis: das Risiko „Kein Wert eines Secrets in einem Log oder Dokument“ ist für das **Log** belegt. Für die
**Dokumente** siehe V-4.

## 3. DoD — Verdikt je Zeile (§2 des Plans)

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Publish-Kette mit zwei Zielen (`build.gradle.kts` zweites Repository `Cloudsmith`, Anmeldedaten allein aus der Umgebung, Version `0.2.2` an beiden Stellen; Workflow ein Job je Ziel mit Einzel-Aufgabe; Probe im Bau rot bei fehlender Aufgabe oder POM-Koordinate) | **erfüllt** | `build.gradle.kts`: Zeile 86 und 210 `version = "0.2.2"`, Zeile 233/241 die Repositories `GitHubPackages`/`Cloudsmith`, Upload-URL `https://maven.cloudsmith.io/pt9912/pg-change-feed/`, Zugangsdaten nur über `System.getenv`. Workflow gelesen: zwei Jobs, kein `needs`, kein `continue-on-error`, `packages: write` genau einmal (Zeile 71, GitHub-Packages-Job), der Cloudsmith-Job trägt nur `contents: read`; Werte über `env:` und `docker run -e NAME`. Die Jobs riefen im realen Lauf die Einzel-Aufgaben (§2). Mutationen §4: fünf Eingabeseiten rot |
| 2 | Anwender-Sicht: Kotlin-README (Cloudsmith zuerst, GitHub Packages mit Token, Namensnennung, kennungsfrei); Handbuch an vier Stellen, `Version:` und Historie | **erfüllt** | README gelesen (§5); `make sdk-public-doc-check` Exit 0, Mutation M5 rot. Handbuch: `Version: 1.65` (Parent `dbb87155`: `1.64`), Historie-Zeile 1.65 vorhanden, die vier Kotlin-Stellen (Zeilen 1078, 1180, 1257, 1412) nennen Cloudsmith zuerst |
| 3 | Träger-Nachzug (Pflichtenheft `LH-FA-SST-009.a` und [`SPEC-028`](../../spec/pflichtenheft.md), `harness/README.md`, `releasing.md` §4, Kommentare in Dockerfile, `sdk.mk`, `sdk-pack-kotlin.sh`; Suchlauf im Plan) | **erfüllt, mit V-1 und V-5** | Pflichtenheft Zeilen 317 bis 324 und Zeile 830 tragen zwei Vertriebsziele und `0.2.2`, **ohne** ADR- und Slice-Bezug ([`AGENTS.md`](../../AGENTS.md) §3.4, `git grep` im Absatz gelesen); `releasing.md` `Version: 1.11`; `harness/README.md` Zeile 148 und `make sdk-pack-kotlin`-Zeile nennen `0.2.2` und die Probe; die Kommentare nennen `<Version>` statt einer Zahl. Suchlauf: §6 |
| G | `make gates`, `make sdk-pack-kotlin`, `make sdk-public-doc-check`, `make test-sdk-kotlin-release-tag-info` grün | **erfüllt** | eigene Läufe §1, alle Exit 0; die im Plan gedruckten Zeilen (`Coverage 83.40%`, `commit-traceability: OK`, `generated-sync: OK`, a-check `gesamt: 0 Befund(e)`) decken sich mit meinen. Die Test-/Bau-Zeilen von `make sdk-pack-kotlin` (`BUILD SUCCESSFUL in 24s`/`4s`) sind nicht von mir gemessen (`CACHED`), im Lauf-Log des realen Laufs stehen `BUILD SUCCESSFUL in 37s`/`7s` (Job Cloudsmith) für denselben Bau |
| P | Post-Push-Beleg: Tag, beide Jobs `success`, anonymer POM-Abruf HTTP 200 | **erfüllt (Beleg liegt vor), Haken im Plan noch `[ ]`** | §2. Die Zeile im Plan (Zeile 170) trägt weiter die Zusage-Form („erwartet …“); das Eintragen des Belegs in §6/§7 und das Setzen des Hakens ist Aufgabe des Planners bei der Closure |
| R | Review durchgeführt, Report liegt vor | **erfüllt, mit ehrlicher Benennung** | Report vorhanden (0 HIGH, 2 MEDIUM, 2 LOW, 3 INFO); Nachmessung der Fixrunde in §7 |
| C | Closure-Notiz mit Lerneintrag | **korrekt offen** | §7 des Plans trägt „*(wird bei der Closure durch den Planner gefüllt)*“ |
| B | Beobachtungs-Register | **korrekt offen** | Closure-Pflicht |
| Ri | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | §6 trägt 12 Einträge mit „offen“/Erwartung; die Ausgänge setzt die Closure — Vorschlag in §8 |
| Pa | Die drei Paarungen | **korrekt offen** | hängen an der nächsten Welle-Closure ([welle-transformationen](../plan/planning/welle-transformationen.md), Adresse in Plan §4) |

`[x]` sind vier Zeilen (drei Liefer-Punkte, Gates und Belege) und die Review-Zeile, `[ ]` fünf (Post-Push-Beleg,
Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge, Paarungen), zusammen 10 — gezählt am Plan (Zeilen 128 bis
191). Kein `[x]` ohne Beleg; das `[ ]` der Post-Push-Zeile ist inhaltlich belegt und formal offen.

## 4. Mutationen der Eingabeseite (dieser Lauf)

Aufbau: die Datei per Bash-Ersetzung in sich selbst geschrieben, `make sdk-pack-kotlin` ungefiltert in eine
Log-Datei (Exit 2 ist der Fehlercode von `make`, das Rezept endet mit Exit 1), Rot-Meldung aus dem Log gelesen,
Datei per `git checkout` zurückgenommen, danach `git status --short` leer.

| # | Datei · Mutation | Rot gesehen (gedruckt) |
|---|---|---|
| M1 | `build.gradle.kts`: `name = "Cloudsmith"` → `"CloudsmithZ"` (Repository-Name) | Exit 2, `Probe: Gradle-Aufgabe publishMavenPublicationToCloudsmithRepository fehlt` (Abkürzungs-Falle der ersten Probe-Fassung bleibt geschlossen: der aufgelöste Name `…CloudsmithZ…` erfüllt die ganzzeilige Prüfung nicht) |
| M2 | `artifactId = "pgchangefeed-kotlin"` → `"pgchangefeed-kotlinz"` | Exit 2, `Probe: POM trägt nicht <artifactId>pgchangefeed-kotlin</artifactId>` |
| M3 | Publikations-Version `0.2.2` → `0.2.1` (Top-Level bleibt `0.2.2`) | Exit 2, `Probe: POM trägt nicht <version>0.2.2</version>` |
| M4 | Top-Level-Version `0.2.2` → `0.2.1` (Publikation bleibt `0.2.2`) — die umgekehrte Richtung; der Workflow-Abgleich gegen den Tag liest nur diese Zeile und bliebe grün | Exit 2, `Probe: POM trägt nicht <version>0.2.1</version>` |
| M5 | Kotlin-README: „Requires Java 21 or newer.“ → „… (see <ADR-Kennung>).“ (eine ADR-Kennung mit vier Ziffern) | `make sdk-public-doc-check` Exit 2 (Rezept Exit 1), Treffer `README.md:17`, `sdk-public-doc-check: interne Kennung in den SDK-Dateien`; nach `git checkout` und `cmp` gegen die Sicherung wieder Exit 0 |

Fünf Mutationen (gefordert waren drei), je genau die Zusage gefärbt, die sie bindet: Repository-Name (M1), Koordinate
(M2), beide Versionsstellen (M3, M4), Kennungsfreiheit der README (M5). Der Reviewer hatte M1 bis M6 an derselben
Probe gefahren (Repository-Name des GitHub-Ziels, Upload-URL-Literal zusätzlich); meine Läufe sind davon
unabhängig und mit anderen Ersatzwerten. Nicht von mir gefahren: der Slug in der Upload-URL (Prüfung (c)) — der
Reviewer hat sie gesehen (M5 des Reviews), der reale Lauf belegt die URL: der Upload landete unter
`pt9912/pg-change-feed`.

## 5. Kotlin-README (Anwender-Sicht) gegen die real veröffentlichten Artefakte

| Aussage der README | Beleg |
|---|---|
| Koordinate `io.github.pt9912:pgchangefeed-kotlin:0.2.2` | POM und `.module` tragen genau diese Gruppe, dieses Artefakt und diese Version; `maven-metadata.xml` nennt `latest`/`release` `0.2.2` |
| Repository-URL `https://dl.cloudsmith.io/public/pt9912/pg-change-feed/maven/` (`maven { url = uri(…) }`) | Die Artefaktpfade unter dieser Basis antworten anonym 200 (§2); das Snippet nennt `mavenCentral()` daneben |
| „The Cloudsmith repository is public: no account and no token are needed to read it“ | anonymer `curl` ohne Kopfzeile: 200 auf POM, Jar, Sources-Jar, `.module`, `maven-metadata.xml` — die Erwartung (Plan §3 Fixrunde, Zeile „F-1“, Plan §6) ist eingelöst |
| „GitHub Packages always requires authentication to read a package“ | anonymer `curl` auf die POM-URL bei GitHub Packages: **401** |
| „Apart from the Kotlin standard library, the library’s own dependencies are not passed on to your compile classpath“ | `.module`: Variante `apiElements` trägt nur `kotlin-stdlib`, alle übrigen stehen in `runtimeElements`; POM: `kotlin-stdlib` `compile`, alle anderen `runtime` |
| Snippet-Versionen der Zusatzabhängigkeiten (`kotlinx-coroutines-core:1.11.0`, `protobuf-java:4.36.2`, `gson:2.14.0`) | nicht geändert vom Slice, nicht neu gemessen (übernommen aus dem Vorgänger-Slice) |
| Namensnennung Cloudsmiths (Text mit Link auf `https://cloudsmith.com`, „free of charge for open-source projects“) | Zeile vorhanden; der Wortlaut folgt dem Beispieltext der Cloudsmith-Bedingungen, den der Reviewer am 2026-09-25 im Wortlaut gelesen hat (übernommen, nicht neu geladen) |
| Anwender-Sprache, keine interne Kennung | `make sdk-public-doc-check` Exit 0; M5 rot; die Datei ist rein englisch |

Ergebnis: die drei Snippets (Repository, Koordinate, GitHub-Packages-Weg) sind gegen die reale Veröffentlichung
gehalten; kein Snippet ist kompiliert (kein Konsumenten-Projekt gebaut; das war nie Teil der DoD, der Plan nennt
es als abgeleitet).

## 6. Suchlauf-Feld (§3 des Plans), an beiden Ständen nachgefahren

Parent `dbb87155` bzw. Fixrunden-Parent `a8e765dd` per `git grep … <sha>`, Diff-Stand `HEAD` = `2e9d8db8`.
Gedruckte Zahlen (`git grep -c`, Summen per `awk` aus den gedruckten Zählwerten, also abgeleitet):

| Plan-Zeile | Meine Messung | Ergebnis |
|---|---|---|
| Secret-Regime-Aussagen (fünf Suchwörter, acht Träger-Pfade) | Parent `dbb87155`: **51** in 9 Dateien (Workflow 10, Handbuch 10, `releasing.md` 14, `harness/README.md` 1, Dockerfile 4, Kotlin-README 2, `build.gradle.kts` 6, Pflichtenheft 2, Roadmap 2) — gleich der Plan-Angabe; `HEAD`: **47** (Workflow 7, Handbuch 9, `releasing.md` **16**, `harness/README.md` 1, Dockerfile 2, Kotlin-README 4, `build.gradle.kts` 2, Pflichtenheft 4, Roadmap 2) | Parent bestätigt; Diff-Stand im Plan (**45**) ist der Stand vor der Fixrunde, `releasing.md` trägt seither zwei Zeilen mehr (V-5) |
| GitHub Packages weit gefasst (ohne ADR, `done/`, Reviews, Register) | Parent `dbb87155`: **54** minus 18 (Plan-Datei) gleich **36**, wie im Plan; `HEAD` ohne Plan-Datei: **40** in 10 Dateien (Workflow 4, Roadmap 3, Handbuch 7, `releasing.md` **14**, `harness/README.md` 1, `sdk.mk` 1, Dockerfile 1, Kotlin-README 3, Pflichtenheft 5, Tabellentest-Skript 1) | Parent bestätigt; Plan-Diff-Stand **38** ist der Stand vor der Fixrunde (`releasing.md` 12 zu 14, V-5) |
| Version `0.2.1` der Kotlin-Träger | Parent: **11** Kotlin-Zeilen; `HEAD`: **2** (`releasing.md` Zeile 360, dort die Tag-Namen `sdk-kotlin-v0.2.0`/`sdk-kotlin-v0.2.1` der bewiesenen Läufe, und die Historie-Zeile 1.11) | Parent bestätigt; die zwei Treffer sind Records bzw. Lauf-Belege, keine Träger einer Version (V-5) |
| Job-Kennung `sdk-kotlin-release` | Job-Ids `sdk-kotlin-github-packages` und `sdk-kotlin-cloudsmith` im Workflow; `sdk-kotlin-release` bleibt Workflow-Name und Dateiname | bestätigt |
| Zählwort „Alle drei Secrets“ | Parent: 1 Treffer (`releasing.md` Zeile 128); `HEAD`: 1 Treffer, in der Historie-Zeile 1.10 (Zitat) | bestätigt |
| Offene Sätze („vierter Vertriebsweg“, „vierte Sprache“) in `spec` | Parent: Zeilen 320/321 (Träger) und 868, 870 (Historie); `HEAD`: Zeile 323 („vierte Sprache oder ein weiterer Vertriebsweg“, Träger), 871, 873, 885 (Historie) | bestätigt |
| Zugangsdaten nicht als Build-Argument | `git grep -n -E '^ARG \|--build-arg' HEAD -- .github/workflows/sdk-kotlin-release.yml sdks/kotlin/Dockerfile`: 0 Zeilen, Exit 1 | bestätigt |
| `packages: write` | `git grep`: genau eine Zeile im Workflow (Zeile 71, GitHub-Packages-Job); die zweite Fundstelle ist die Historie-Zeile 1.10 von `releasing.md` | bestätigt |
| Suchlauf der Fixrunde (vier Befehle) | „End-zu-Ende bewiesen“ am `HEAD` 2 Zeilen (Historie-Records), alter Schrittname `./gradlew publish, im Docker-Image` steht nirgends mehr; „anonym lesbar“-Familie 13 Zeilen in 6 Dateien | stichprobenhaft bestätigt |

Nicht gefunden: eine weitere Datei außerhalb der Tabellen, die GitHub Packages als **einziges** Vertriebsziel des
Kotlin-Packages führt; kein verbliebener Satz „kein externes Secret“/„ausschließlich `GITHUB_TOKEN`“ in Workflow,
Dockerfile, `build.gradle.kts`, `releasing.md`, `harness/README.md`.

## 7. Review-Findings nachgemessen

| Finding | Was ich gemessen habe | Verdikt |
|---|---|---|
| F-1 (MEDIUM) „anonym lesbar“ vor dem Beleg | Der Plan hat es in der Fixrunde ehrlich als **Erwartung** benannt (Plan §3 Fixrunde, Zeile „F-1“; §6 Risiko „Tatsachenbehauptung vor dem Beleg“) und die Träger unverändert gelassen. Der Post-Push-Beleg (§2) macht die Aussage wahr: anonyme 200 auf alle fünf Dateien. Welche Stellen ihre Erwartung eingelöst haben: Kotlin-README Zeilen 19 und 21; Handbuch Zeilen 1078, 1180, 1257, 1412; Pflichtenheft Zeile 321; `releasing.md` Zeilen 324 bis 325; Workflow-Kopfkommentar Zeile 6; `harness/README.md` Zeile 148 — alle **bleiben unverändert**, sie sind jetzt belegt | **entfallen** (Risiko-Ausgang), keine Textänderung nötig |
| F-2 (MEDIUM) `releasing.md` „End-zu-Ende bewiesen“ und alter Schrittname | `releasing.md` Zeilen 355 bis 370: der Absatz heißt „Belegt ist der Publish nach GitHub Packages, nicht die Zwei-Job-Struktur“, nennt die Läufe `36117929552`/`36129672852` und die heutigen Schrittnamen (wie im Workflow, Zeilen 93 und 130); die Aussage „unbewiesen“ dort ist mit dem Lauf `36201941235` **überholt** (V-1) | in der Fixrunde behoben, seit dem Post-Push-Lauf nachzuziehen |
| F-3 (LOW) „Name und API-Key“ gegen Service-Slug | `git grep -n 'dessen Name und$' HEAD`: 0 Zeilen; Zeile 316 sagt Slug mit Rückfall; der reale Lauf hat den Slug angenommen (§2: Task `success`) | behoben; der Rückfall ist **entfallen** (V-1) |
| F-4 (LOW) Verweis auf ein Kontingent-Risiko, das §6 nicht führte | Plan §6 trägt jetzt „Kontingent-Angabe“ | behoben |
| F-5 (INFO) Probe bindet nicht den Aufrufer | Plan §6 „Probe bindet nicht den Aufrufer“ als benannte Grenze; im realen Lauf trugen beide Aufruf-Literale des Workflows dieselben Namen wie die Probe (kein `Task … not found`); das Risiko bleibt strukturell offen | **weiter offen** → Register, wie im Plan |
| F-6 (INFO) fünf wortgleiche Schritte in zwei Jobs | keine Aktion | unverändert |
| F-7 (INFO) erster realer Lauf der Job-Struktur | (a) Der Cloudsmith-Job lief ohne `GITHUB_ACTOR`/`GITHUB_TOKEN`-Werte und der GitHub-Job ohne `CLOUDSMITH_*`, beide Gradle-Läufe `BUILD SUCCESSFUL` — Gradle bemängelt die Zugangsdaten des jeweils anderen Repositories nicht. (b) `docker run -e GITHUB_ACTOR` ohne Wert erbte den Standardwert des Runners (Zeile `-e GITHUB_ACTOR` im Log, Publish nach GitHub Packages `success`). (c) `.module` angenommen, Slug angenommen, anonym sichtbar; Doppel-Upload nicht ausgelöst | **belegt** |

## 8. Risiken §6 — Zuordnung zu den Post-Push-Belegen (Vorschlag für die Ausgänge, die der Planner setzt)

| Risiko (Plan §6) | Beleg aus dem realen Lauf und meinen Abrufen | Vorschlag |
|---|---|---|
| Erster realer Lauf der Job-Struktur | Lauf `36201941235`, beide Jobs `success`, keine Berechtigungsfehler (`packages: write` nur im GitHub-Job) | *entfallen* |
| Benutzername: Service-Slug oder Service-Name | Slug gesetzt, Upload `success`, kein `401`/`403` im Log | *entfallen* — der Slug wirkt; ob der Name auch wirkte, bleibt ungeprüft und ist ohne Folge |
| Moduldatei (`.module`) | `.module` anonym 200, 4711 Bytes, gültiges Gradle-Metadatenmodell | *entfallen* — angenommen |
| Anonyme Sichtbarkeit | anonym 200 auf POM, Jar, Sources-Jar, `.module`, `maven-metadata.xml`; erst nach rund zwei Minuten (Indexierung, `releasing.md` nennt das Sammelfenster) | *entfallen* |
| Wiederholung und Doppel-Upload derselben Version | **nicht belegt**: der erste Lauf war grün, es gab keine Wiederholung; ein zweiter Upload derselben Version wurde nie ausgelöst | *entfallen* nur nach der Regel des Plans („grüner erster Lauf“); das Verhalten bei Wiederholung bleibt **unbeobachtet** und gehört ehrlich als Grenze in die Closure-Notiz |
| Kein Wert eines Secrets in einem Log oder Dokument | Log: belegt (§2); Dokumente: siehe V-4 | Log *entfallen*; Dokumente mit V-4 zu entscheiden |
| README nennt `0.2.2`, bevor der Tag gesetzt ist | Version `0.2.2` abrufbar (§2) | *entfallen* |
| Parallele Änderungen an geteilten Trägern | `git diff dbb87155..HEAD --stat`: 12 Dateien, zehn Träger der Plan-Tabelle §3, die Plan-Datei und der Review-Report; kein Move außer den zwei Lifecycle-Renames (`open`→`next` in `d11adc14`, `next`→`in-progress` in `9b045c95`, je `0 0` in `git show --numstat`, also reine Renames) | *entfallen* |
| Zwei-Quellen-Drift Handbuch gegen Pflichtenheft | beide Träger nennen Cloudsmith zuerst (anonym), GitHub Packages mit Token, `0.2.2` — gelesen | *entfallen* bei Deckung |
| Kontingent-Angabe | **nicht geprüft** (die Usage-Seite des Repositories in der Web-App habe ich nicht eingesehen; Paketgröße rund 190 kB je Version, gemessen als Summe der drei Dateien) | *weiter offen*, Betreiber-Blick auf die Usage-Seite; ohne Folge bei dieser Größe |
| Tatsachenbehauptung vor dem Beleg (F-1) | §7 | *entfallen* |
| Probe bindet nicht den Aufrufer (F-5) | §7 | *weiter offen* → Register |

## 9. Entscheidungs-Konformität

| Entscheidung | Festlegung | Beleg | Verdikt |
|---|---|---|---|
| [`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md) Festlegung 1 | zwei Ziele, ein Tag, Kotlin allein | `sdks/csharp/**`, `sdks/python/**` und ihre Workflows nicht im Diff (`git diff --stat`), Tag-Namensraum unverändert | konform |
| Festlegung 2 | Open-Source-Repository, Namen `pt9912`/`pg-change-feed`; Republishing bleibt aus | Upload- und Download-URL wie festgelegt, der reale Lauf und die anonymen Abrufe belegen den Namespace; der Schalter „Republishing“ ist nicht eingesehen | konform, Schalter ungeprüft |
| Festlegung 3 | zwei Repository-Secrets, keine Action, Publish in der Docker-Stufe zur Laufzeit | Workflow: keine neue `uses:`-Zeile (nur der SHA-gepinnte Checkout mit Tag-Kommentar, [`AGENTS.md`](../../AGENTS.md) §3.8), Werte nur über `env:` und `docker run -e NAME`, kein Build-Argument. **Der Wortlaut „`CLOUDSMITH_USERNAME` (Name des Services)“ der ADR ist durch den Lauf widerlegt**: im Secret steht der Slug (V-3) | konform, mit V-3 |
| Festlegung 4 | zweites `maven`-Repository, Einzel-Aufgaben je Ziel, `CMD` der Stufe bleibt | Namen `publishMavenPublicationTo<…>Repository` im realen Log gelesen | konform |
| Festlegung 5 | ein Job je Ziel, Berechtigungen minimal, keine Reihenfolge | Workflow gelesen; im Lauf starteten beide Jobs um 23:40:47Z | konform |
| Festlegung 6 | Anwender-Sicht, Namensnennung, Beleg | §5 | konform |
| Festlegung 7 | GitHub Packages und Tag-Namensraum bleiben | GitHub Packages listet die Version `0.2.2` (§2) | konform |
| Fitness Function, Zeile 3-Ersatz | GitHub-Job nur `GITHUB_TOKEN`, Cloudsmith-Job nur die zwei Cloudsmith-Secrets | Workflow gelesen: je Job genau diese `env:`-Werte | konform |
| Fitness Function, Docker-Bau | Bau rot ohne Repository-Block oder mit falscher POM | §4, M1 bis M4 | konform |
| [`ADR-0109`](../plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md) Festlegung 4/5 | SemVer `0.x`, Docker-only bis `publish`, `maven-publish` | `0.2.2` PATCH; Publish in der Stufe `publish` (Log: `pg-change-feed:sdk-kotlin-publish`, `docker run` mit gepinntem Image aus dem Bau) | konform |
| [`AGENTS.md`](../../AGENTS.md) §3.1, §3.3, §3.4, §3.5, §3.6, §3.7 | Docker-only, reine Moves, Sicht ohne ADR-Bezug, Immutabilität, keine Gate-Lockerung, Kommentar-Klassen | keine Host-Toolchain in den geänderten Skripten; `make doc-immutable RANGE=dbb87155..HEAD` Exit 0 (keine ADR im Diff); Pflichtenheft-Absätze ohne ADR-/Slice-Bezug; kein `nolint`; Kommentare in Workflow, `build.gradle.kts`, Dockerfile, `sdk.mk` im Indikativ ohne Chronik gelesen | konform |
| [`AGENTS.md`](../../AGENTS.md) §3.10 | neuer Workflow nur mit realem, grünem Post-Push-Lauf | Lauf `36201941235` (§2) | **erfüllt** |
| [`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B, §3.13 | Beleg-Anker; Suchlauf beider Stände | Plan-Felder mit Befehl und Zahl; Suchlauf nachgefahren (§6); zwei Stände der Diff-Zahlen tragen keinen Stand-Vermerk (V-5) | konform mit V-5 |
| [`AGENTS.md`](../../AGENTS.md) §3.11 | kein host-lokaler Pfad | `make docs-check` 0 Befunde; dieser Report trägt keinen | konform |

## 10. Befunde

Kein Befund blockiert. V-1 bis V-5 sind Nachzüge für den Planner (kein Code, keine Spec-Änderung): der Post-Push-Lauf
hat mehrere Aussagen von „unbewiesen/erwartet“ in den Ist-Zustand gehoben, die im Text noch als Erwartung stehen.

- **V-1 (LOW) — Träger-Sätze „unbewiesen“ nach dem Post-Push-Lauf überholt.** Der Lauf `36201941235` hat sie
  widerlegt. Nachzuziehen (Zeilen am Stand `HEAD`):
  `docs/user/releasing.md` Zeile 37 („diese Struktur ist bis zu ihrem ersten Tag-Lauf nicht bewiesen“);
  Zeilen 133 bis 140 (Werttyp `CLOUDSMITH_USERNAME`: „ob Cloudsmith … stattdessen den Service-Namen verlangt, zeigt
  erst der erste Tag-Lauf“ — der Slug hat gewirkt, der Rückfall-Satz entfällt); Zeilen 339 bis 352 (Absatz „Offen bis
  zum ersten Tag-Lauf“: Job-Struktur, anonymer Abruf, Moduldatei, Slug/Name sind jetzt belegt; offen bleibt allein „ob
  die Paketseite die POM-Beschreibung anzeigt“, der Kontingent-Blick und die Doppel-Upload-Frage; „liegen nach dem
  ersten Tag-Lauf erwartet die Versionen ab `0.2.2`“ wird Ist-Zustand); Zeilen 365 bis 370 („Belegt ist der Publish nach
  GitHub Packages, nicht die Zwei-Job-Struktur … unbewiesen“: mit dem Lauf `36201941235` belegt, beide Läufe von
  `0.2.0`/`0.2.1` bleiben als frühere Belege); Zeile 316 (Betreiber-Voraussetzung: Slug statt „Rückfall“);
  `harness/README.md` Zeile 148 („… sind bis zum ersten Tag-Lauf mit beiden Zielen (`sdk-kotlin-v0.2.2`)
  unbewiesen“); `releasing.md` `Version:` und Historie-Zeile dazu. Das ist derselbe Vorgang wie die Erwartung der
  Fixrunde (Plan §3, Zeile „F-1“), nur auf die vier weiteren „unbewiesen“-Sätze ausgedehnt, die der Plan dort nicht
  nannte (er nannte die „anonym lesbar“-Stellen).
- **V-2 (INFO) — DoD-Haken und §6/§7 tragen den Beleg noch nicht.** Die Post-Push-Zeile (Plan Zeile 170) steht in
  Zusage-Form und mit `[ ]`; §6 trägt „offen bis zum Tag-Lauf“; §7 ist leer. Vorschlag für die Ausgänge: §8 dieses
  Reports. Die Zeile des Review-Hakens (Zeilen 177 bis 187) trägt „F-1 ist nicht behoben, sondern als Erwartung
  benannt … das Kriterium ‚kein offenes MEDIUM‘ ist damit bis zu diesem Beleg … getragen, nicht erfüllt“: mit dem
  Beleg ist F-1 aufgelöst (§7), der Satz gehört in der Closure auf den Ist-Zustand.
- **V-3 (INFO) — [`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md) Festlegung 3 sagt „Name des
  Services“, der Lauf zeigt: der Slug.** Die ADR ist `Accepted` und unberührbar ([`AGENTS.md`](../../AGENTS.md) §3.5),
  `releasing.md` trägt die Wahrheit (nach V-1). Eine Korrektur wäre eine neue ADR mit `Supersedes` oder ein Vermerk in
  der Closure-Notiz; das entscheidet der Planner (Kosten: eine Aussage über einen Werttyp, kein Verhalten).
- **V-4 (LOW) — Secret-Hygiene der Dokumente.** Die Zusage des Plans („der Service-Slug steht in keinem Dokument
  dieses Repos“, Plan §6) ist im **Log** belegt (§2). In den **Dokumenten** steht der Name des Services im Plan
  (Zeile 462, Rückfalltext), und der Review-Report nennt in Zeile 94 ein Suchmuster, das neben dem Namen auch das
  vierstellige Suffix des Slugs enthält (`git grep -n -i` über das Repo: genau diese zwei Dateien). Zusammen lässt
  sich der Slug (der Anmeldename, nicht der API-Key) rekonstruieren; der Slug ist im Log maskiert, weil er als
  Secret geführt wird. Vorschlag: das Suffix in der Review-Zeile durch einen Platzhalter ersetzen (Zitat-Korrektur an
  einem Record, Beleg ist die Commit-Kennung, [`AGENTS.md`](../../AGENTS.md) §3.5) und den Ausgang „entfallen“ des
  Risikos nur für das Log setzen. Ich habe das Suffix hier nicht wiederholt.
- **V-5 (INFO) — Suchlauf-Zahlen ohne Stand-Vermerk.** In Plan §3 stehen die Diff-Stand-Zahlen **45** (Secret-Regime)
  und **38** (GitHub Packages) und „0 Kotlin-Zeilen mit `0.2.1`“ als Messung am Arbeitsbaum vor der Fixrunde. Am `HEAD`
  messe ich **47**, **40** und **2** (`releasing.md`: Fixrunde-Zeilen zu F-2/F-3, die Historie-Zeile 1.11 nennt die
  Tag-Namen). Die Parent-Zahlen (51, 36, 11) stimmen exakt. Instanz A (Zahl mit Lauf): der Plan nennt den Stand
  nicht ausdrücklich; die Fixrunden-Tabelle trägt ihren eigenen Stand (`a8e765dd` gegen Arbeitsbaum). Der Planner
  vermerkt bei der Closure „Diff-Stand vor der Fixrunde“ oder zieht die Zahlen nach.
- **V-6 (INFO) — Indexierungsverzögerung ist für Anwender sichtbar.** Der Publish-Task endet nach 13 s, der anonyme
  Abruf antwortete erst nach rund zwei Minuten mit 200; in dem Fenster liefert Cloudsmith 404. Die README nennt das
  Fenster nicht (Anwender-Text, sie verlässt sich auf die Veröffentlichung nach dem Tag); `releasing.md` nennt das
  „Sammelfenster, ggf. mit Wiederholung“. Kein Befund, nur die gemessene Größe für den Nachzug V-1.
- **V-7 (INFO) — Byte-Gleichheit.** Jar und Sources-Jar sind größengleich mit dem lokalen Bau, die SHA-256-Summen
  weichen ab (Zeitstempel im Archiv). Der Beleg „dieselben Artefakte in beiden Zielen“ ([`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md)
  Festlegung 4) ist hier nur über die Größen und über den gemeinsamen Bau im Job belegt, nicht per Hash gegen GitHub
  Packages (dort ist der Abruf authentifiziert; nicht gefahren).
- **V-8 (INFO) — Grenzen der Verifikation.** Nicht geprüft: die Anzeige der POM-Beschreibung auf der Cloudsmith-
  Paketseite, die Usage-/Kontingent-Anzeige, der Schalter „Republishing“, ein Doppel-Upload derselben Version, ein
  Konsumenten-Projekt gegen das Repository (die Snippets sind gegen die Artefakte gelesen, nicht gebaut).

## 11. Ergebnis

| Prüfpunkt | Ergebnis |
|---|---|
| DoD §2 — Liefer-Punkte 1 bis 3, Gates und Belege | **erfüllt** (je mit eigenem Beleg; Punkt 3 mit V-1, V-5) |
| DoD §2 — Post-Push-Beleg | **erfüllt durch eigene Messung** (Tag, Lauf `36201941235` beide Jobs `success`, anonym 200), Haken im Plan formal offen (V-2) |
| DoD §2 — Review-Zeile | **erfüllt**; F-1 mit dem Beleg aufgelöst |
| DoD §2 — Closure-Notiz, Register, Risiko-Ausgänge, Paarungen | **korrekt offen** (Vorschlag §8) |
| Review-Findings F-1 (MEDIUM), F-2 (MEDIUM) | **behoben bzw. durch den Beleg entfallen** (V-1 verlangt den Nachzug der überholten „unbewiesen“-Sätze) |
| Review-Findings F-3, F-4 (LOW) | **behoben** |
| Review-Findings F-5 bis F-7 (INFO) | F-7 **belegt**, F-5 **weiter offen** (Register), F-6 unverändert |
| Mutationen der Eingabeseite | fünf gefahren (M1 bis M5), alle rot, alle zurückgenommen |
| Kotlin-README gegen die Veröffentlichung | Koordinate, Repository-URL, Token-Aussagen, Klassenpfad-Aussage bestätigt |
| Suchlauf-Feld | Parent-Zahlen exakt bestätigt (51, 36, 11), Diff-Stand-Zahlen mit Stand-Vermerk nachzuziehen (V-5) |
| Secret-Hygiene | Log **sauber** (drei Werte maskiert, keine Klartext-Zeile); Dokumente mit V-4 |
| Gates | **`make gates` Exit 0**, `make docs-check`, `make sdk-public-doc-check`, `make test-sdk-kotlin-release-tag-info`, `make sdk-pack-kotlin`, `commit-traceability`/`doc-commits`/`doc-immutable` je Exit 0 |
| CI zum Tag-Commit | `sdk-kotlin-release`, `ci`, `examples`, `e2e` success |

## Verdikt

**Bestätigt.** Der Slice liefert, was seine DoD zusagt, und der reale Post-Push-Lauf ([`AGENTS.md`](../../AGENTS.md)
§3.10) liegt vor und trägt: beide Jobs enden `success`, die Einzel-Aufgaben liefen, das Package ist unter der
Cloudsmith-Download-Basis der README anonym abrufbar (POM, Jar, Sources-Jar, Moduldatei, `maven-metadata.xml`) und
GitHub Packages listet `0.2.2`. Damit lösen sich vier der Erwartungen des Plans ein: Slug statt Name, Moduldatei
angenommen, anonyme Sichtbarkeit, Job-Struktur. Die Erwartung der Fixrunde zu F-1 ist eingelöst — alle „anonym
lesbar“-Stellen bleiben unverändert. Offen bleiben Wiederholung/Doppel-Upload, Kontingent und Paketseiten-Anzeige als
benannte Grenzen (V-8) sowie das strukturelle Risiko F-5.

**Übergabe:** Verifier → Planner. Vor der Closure: (1) V-1 — die „unbewiesen/erwartet“-Sätze in `releasing.md`
(Zeilen 37, 133 bis 140, 316, 339 bis 352, 365 bis 370) und `harness/README.md` Zeile 148 auf den Ist-Zustand ziehen,
`releasing.md` `Version:` heben; (2) Post-Push-Beleg samt Lauf-Kennung `36201941235`, Tag-Commit und den vier
Abruf-Zeilen in Plan §6/§7 eintragen, den Haken setzen, die Risiko-Ausgänge nach §8 setzen; (3) V-4 entscheiden
(Suffix in der Review-Zeile); (4) V-3 im Closure-Vermerk entscheiden (ADR-Wortlaut „Name des Services“); (5) V-5 den
Stand der Zahlen vermerken; (6) Closure-Notiz mit Lerneintrag und Register (die Finding-Klassen des Reviews, darunter
„Tatsachenbehauptung vor dem Beleg“ und „Probe bindet nicht den Aufrufer“); danach der reine `git mv` nach `done/`
([`AGENTS.md`](../../AGENTS.md) §3.3, Fall 2). Dieser Report ändert nur `docs/reviews/`.
