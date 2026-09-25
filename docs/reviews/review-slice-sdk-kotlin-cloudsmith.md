# Review-Report: slice-sdk-kotlin-cloudsmith — 2026-09-25

**Review-Art:** Code — der Diff fügt dem Kotlin-SDK-Release ein zweites Publish-Ziel (Cloudsmith) hinzu:
zweites Repository in `build.gradle.kts`, Workflow mit einem Job je Ziel und Einzel-Aufgaben,
Probe der Publish-Konfiguration in der Docker-Stufe `pack`, Version `0.2.2`, Kotlin-README
(Anwender-Sicht, Namensnennung), Träger-Nachzug (Pflichtenheft, Handbuch, `releasing.md`,
`harness/README.md`, Werkzeug-Kommentare); geprüft gegen Plan, ADR und `AGENTS.md` Hard Rules
(Modul 10 §Drei Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Slice `slice-sdk-kotlin-cloudsmith` (ohne Welle), Diff-Range `dbb87155..577876ac`
(7 Commits, 11 Dateien, +417/−221 laut `git diff --stat`; HEAD `577876ac`, Baum sauber, gepusht,
kein Tag im Range).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere
HIGH-Klassen ergänzt, u. a. Zahl-im-Träger, Beleg-Satz, Zusage-ohne-Eingabeseite,
Form-Vorbild-Kopie, Handbuch-Zug).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-25.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** (`v6.9.0` ·
> `regelwerk/<datei>.md` §<Abschnitt>). Ein `pfad`-Feld auf den **geprüften
> Gegenstand** zitiert den Stand des Laufs und darf ihn festhalten.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-sdk-kotlin-cloudsmith` (§1 Ziel und Abgrenzung, §2 DoD, §3 Plan mit
  Plan-Nachzug, Mutations-Tabelle und Suchlauf-Feld, §4 Trigger, §5 Closure-Trigger, §6 Risiken)
- [`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md) (Accepted; Festlegungen
  1 bis 7, Fitness Function), [`ADR-0109`](../plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md),
  [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md); Architect-Verdikt
  `architect-verdict-sdk-kotlin-cloudsmith`
- [`LH-FA-SST-009`](../../spec/lastenheft.md); [`SPEC-028`](../../spec/pflichtenheft.md) und
  `LH-FA-SST-009.a` im Pflichtenheft
- `AGENTS.md` (Hard Rules §3.1 Docker-only, §3.3 Commit-Struktur, §3.4 Sicht ohne ADR-Bezug,
  §3.7 Kommentar-Klassen, §3.8 Action-Pinning, §3.9 Exit-Code, §3.10 realer Post-Push-Lauf,
  §3.11 host-lokale Pfade, §3.12 Herkunft von Aussagen, §3.13 Träger-Nachzug),
  `harness/conventions.md` (`MR-000`/`MR-001`)
- Cloudsmith-Bedingungen: Abschnitt „Attribution“ der Seite
  `docs.cloudsmith.com/resources/open-source-hosting-policy`, am 2026-09-25 per `curl` geladen und
  im Wortlaut gelesen
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-sdk-readme-nutzerdoku.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem Implementer-Bericht übernommen;
Exit-Codes ungepiped in Log-Dateien gesichert, gedruckte Zeilen zitiert):

- **`make sdk-pack-kotlin` am Stand `577876ac`:** Exit 0. Die Docker-Schichten (auch die
  Probe-Schicht `#19`) waren `CACHED` — identischer Bau-Kontext wie beim Implementer, die
  Test-Zeilen stammen deshalb nicht aus diesem Lauf; die Probe wurde stattdessen durch die
  Mutationen unten real ausgeführt. `sdks/kotlin/dist` trägt `pgchangefeed-kotlin-0.2.2.jar` und
  `pgchangefeed-kotlin-0.2.2-sources.jar`.
- **Weitere Läufe:** `make test-sdk-kotlin-release-tag-info` Exit 0, gedruckt „alle Fälle
  bestanden“; `bash tools/harness/sdk-kotlin-release-tag-info.sh sdk-kotlin-v0.2.2` Exit 0, gedruckt
  `version=0.2.2`; `make sdk-public-doc-check` Exit 0, gedruckt „keine interne Kennung unter sdks“;
  `make test-sdk-public-doc-check` Exit 0, „alle Fälle bestanden“; `make docs-check` Exit 0,
  gedruckt „d-check: 1181 Datei(en) geprüft, 0 Befund(e)“.
- **Mutationen der Eingabeseite** (einzeln nacheinander per reiner Bash-Ersetzung, `make` ungepiped,
  danach `git checkout` und `cmp` gegen die Sicherung — jeweils „restored“, `git status` sauber):

  | Nr. | Mutation | Ergebnis |
  |---|---|---|
  | M1 | `name = "Cloudsmith"` zu `"Cloudsmithy"` in `build.gradle.kts` | `make sdk-pack-kotlin` Exit 2, gedruckt `Probe: Gradle-Aufgabe publishMavenPublicationToCloudsmithRepository fehlt` |
  | M2 | `name = "GitHubPackages"` zu `"GitHubPackagesX"` | Exit 2, `Probe: Gradle-Aufgabe publishMavenPublicationToGitHubPackagesRepository fehlt` |
  | M3 | `artifactId` zu `pgchangefeed-kotlinx` | Exit 2, `Probe: POM trägt nicht <artifactId>pgchangefeed-kotlin</artifactId>` |
  | M4 | Version der Publikation `0.2.2` zu `0.2.1` (Top-Level bleibt `0.2.2`) | Exit 2, `Probe: POM trägt nicht <version>0.2.2</version>` |
  | M5 | Slug in der Upload-URL zu `pg-change-feed2` | Exit 2, `Probe: Cloudsmith-Upload-URL fehlt in build.gradle.kts` |
  | M6 | `(see ADR-0123)` in die Kotlin-README eingefügt | `make sdk-public-doc-check` Exit 2 (Rezept Exit 1), Treffer `README.md:17` |

  Alle sechs decken sich mit den im Plan §3 eingetragenen gesehenen Rot-Ständen. Nicht gefahren:
  Top-Level-Version umgekehrt (`0.2.1`, Publikation `0.2.2`); dieselbe Bindung wie M4, im Plan als
  gesehen eingetragen.
- **Suchlauf-Feld des Plans an beiden Ständen nachgefahren** (Parent `dbb87155` per `git grep … dbb87155`,
  Diff-Stand am Arbeitsbaum): Zeile 1 (Secret-Regime, fünf Suchwörter): Parent 51 Zeilen (Workflow 10,
  `releasing.md` 14, Handbuch 10, `build.gradle.kts` 6, Dockerfile 4, Kotlin-README 2,
  Pflichtenheft 2, Roadmap 2, `harness/README.md` 1), Diff-Stand 45 (7/14/9/2/2/4/4/2/1) — gleich der
  Plan-Angabe. Zeile 2 (GitHub Packages weit, ohne ADR/`done`/Reviews/Register): Parent 54 minus 18
  (Plan-Datei) gleich 36 in 12 Dateien, Diff-Stand 38 in 10 Dateien ohne Plan-Datei — gleich. Zeile 3
  (Version `0.2.1` der Kotlin-Träger): Parent 11, Diff-Stand 0; `0.2.2` steht in Kotlin-README (1),
  `build.gradle.kts` (3), `harness/README.md` (2), Pflichtenheft (Kotlin-Absatz, `SPEC-028`, Historie),
  `releasing.md` (3), Handbuch-Historie. Zeilen 4 (`sdk-kotlin-release` als Job-Kennung: keine
  Fundstelle außer Workflow-Name, Dateiname, Skriptnamen, Plan), 6 (`vierter Vertriebsweg|vierte
  Sprache`: Träger nachgezogen Z. 323), 8 (kein `^ARG `/`--build-arg`, Exit 1), 9 (zwei `echo "::error::…"`
  ohne Werte; `packages: write` genau eine Zeile, Z. 71) stimmen. Suchraum-Erweiterung (nicht im
  Plan): `git grep` nach `GitHub-Konto|personal access|Authentifizierung zum Lesen` und
  `maven.pkg.github.com|read:packages|pgchangefeed-kotlin` außerhalb von ADR/`done`/Reviews/Register
  findet keinen weiteren stehenden Träger der Token-Pflicht-Aussage.
- **Cloudsmith-Attribution** im Wortlaut der Seite gelesen: verlangt wird ein Link auf
  `cloudsmith.com` (oder `app.cloudsmith.io`) und die Aussage, dass Cloudsmith freies Hosting
  bereitstellt; Beispieltext „Package repository hosting is graciously provided by
  [Cloudsmith](https://cloudsmith.com)“. Die README-Zeile trägt beides.
- **Secret-Hygiene:** `git grep -n -i 'github-ci\|i7j3'` außerhalb der Plan-Datei: 0 Treffer;
  `git grep -n -i cloudsmith -- .github sdks tools harness spec docs/user README.md examples` nennt
  nur Dateien mit Namen (`CLOUDSMITH_USERNAME`/`CLOUDSMITH_API_KEY`), keinen Wert.
- **Umgebung:** `free -m` verfügbar 18,4 GB vor, 18,5 GB nach den Läufen; dangling Volumes
  (`docker volume ls -q -f dangling=true | wc -l`) 34 vor und 34 nach allen Läufen; kein `prune`,
  keine eigenen Volumes oder Images zurückgeblieben.
- **Nicht gefahren (Grenze):** `make gates` als Ganzes (Gate-Zeilen des Plans §2 wie `Coverage 83.40%`
  nicht nachgemessen, nur `docs-check` und die Probe-Läufe); `make test-sdk-*-integration`; ein realer
  Publish, ein Tag-Push, ein Zugriff auf GitHub Packages oder Cloudsmith mit Zugangsdaten. Der
  Workflow ist nur **gelesen** (YAML-Struktur, `needs`, `permissions`, Ausdrücke, `if:`), nicht
  ausgeführt.

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | MEDIUM | „anonym lesbar / ohne Konto und Token beziehbar“ steht in fünf Trägern als Tatsache, ohne Beleg-Anker oder „erwartet“: Kotlin-README („The Cloudsmith repository is public: no account and no token are needed to read it“), Handbuch (§4, Kotlin-Absatz), Pflichtenheft (`LH-FA-SST-009.a`), `releasing.md` (§4 „Was Anwender sehen“) und Workflow-Kopfkommentar. Der Plan §4 misst am selben Tag `…/v1/repos/pt9912/pg-change-feed/` und `…/public/pt9912/pg-change-feed/maven/` je mit HTTP 404 und schreibt, die README-Aussage „hänge davon ab“; die Aussage ist bis zum Post-Push-Beleg (§5) unbelegt. Die Zeile „auf Cloudsmith liegen erst die Versionen ab `0.2.2`“ in `releasing.md` nimmt zudem einen Zustand vorweg, den erst der Tag-Lauf herstellt. | [`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B; [`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md) Festlegung 6 | `sdks/kotlin/pgchangefeed-kotlin/README.md:21`, `docs/user/benutzerhandbuch.md:1078`, `spec/pflichtenheft.md:321`, `docs/user/releasing.md:349`, `.github/workflows/sdk-kotlin-release.yml:6` | ja — `curl -s -o /dev/null -w '%{http_code}' https://dl.cloudsmith.io/public/pt9912/pg-change-feed/maven/` (am 2026-09-25 vor dem Upload 404 laut Plan §4) bzw. der anonyme POM-Abruf des §5-Belegs | Tatsachenbehauptung vor dem Beleg |
| F-2 | MEDIUM | `releasing.md` schließt den Kotlin-Abschnitt weiter mit „Der Kotlin-SDK-Release-Weg ist damit End-zu-Ende bewiesen, nicht nur implementiert“ und zitiert den Schrittnamen „Nach GitHub Packages veroeffentlichen (./gradlew publish, im Docker-Image)“. Die neuen Absätze direkt davor und §1 sagen, die Zwei-Job-Struktur sei bis zum ersten Tag-Lauf unbewiesen; der zitierte Schritt heißt im Workflow jetzt anders und ruft die Einzel-Aufgabe, nicht `publish`. Zwei Aussagen im selben Abschnitt, keine verweist auf die andere. | Reviewer-Skill „Nachzug widerspricht dem Nachbarn im selben Träger“ · [`AGENTS.md`](../../AGENTS.md) §3.13 | `docs/user/releasing.md:355-362` gegen `docs/user/releasing.md:335-350`, `.github/workflows/sdk-kotlin-release.yml:93` | ja — die Zeilen nebeneinander lesen; `git grep -n 'Nach GitHub Packages' -- .github` | Nachzug widerspricht dem Nachbarn |
| F-3 | LOW | `releasing.md` sagt in den Betreiber-Voraussetzungen, das Service-Konto sei eines, „dessen Name und API-Key die zwei Repository-Secrets tragen“; der Absatz direkt darüber (Secret-Tabelle) nennt den Werttyp von `CLOUDSMITH_USERNAME` den Service-**Slug**, und ob Name oder Slug gilt, lässt er offen. Der Plan §6 führt genau diese Frage als offenes Risiko. | Reviewer-Skill „Nachzug widerspricht dem Nachbarn im selben Träger“ | `docs/user/releasing.md:315-316` gegen `docs/user/releasing.md:134-140` | nein — Lese-Handlung | Nachzug widerspricht dem Nachbarn |
| F-4 | LOW | Der Plan §4 schließt die Diskrepanz der Kontingent-Angaben (Kontoanzeige 500 MB/1 GB gegen 50 GB/200 GB der ADR) mit „Die ADR bleibt unberührt (Risiko §6)“; §6 führt kein Kontingent-Risiko, und die DoD-Zeile „Jedes Risiko aus §6 trägt einen Ausgang“ erfasst es damit nicht. | `AGENTS.md` §3.12 Instanz B (Verweis auf eine Stelle, die den Satz nicht trägt) | `docs/plan/planning/in-progress/slice-sdk-kotlin-cloudsmith.md:351`, `…:403-467` | ja — `grep -n -i 'kontingent' ` über §6 druckt nichts | Verweis auf nicht vorhandenen Eintrag |
| F-5 | INFO | Die Probe bindet die Task-Namen nur am Literal im Dockerfile; die zwei Aufruf-Literale im Workflow (`…GitHubPackagesRepository`, `…CloudsmithRepository`) stehen davon getrennt, ein Auseinanderlaufen färbt keinen Sensor. Ein falscher Name endet im ersten Tag-Lauf laut mit „Task not found“ (Gradle löst nur eindeutige Abkürzungen auf); heute sind beide Stellen wortgleich. Weiter: Prüfung (c) ist ein `grep -F` auf `url = uri("…")` in `build.gradle.kts` und würde von einer Kommentarzeile mit dem Literal erfüllt; Prüfung (a) hängt am Ausgabeformat `:<Name> SKIPPED` von `--dry-run`. | Reviewer-Skill „Zusage ohne Bindung an ihre Eingabeseite“ (Nachbarform, nicht gebrochen: M1 bis M5 färben die Probe) · Maintainability | `sdks/kotlin/Dockerfile:84-102`, `.github/workflows/sdk-kotlin-release.yml:101,139` | ja — Workflow-Literal ändern, `make sdk-pack-kotlin` bleibt Exit 0 | Probe bindet nicht den Aufrufer |
| F-6 | INFO | Fünf Schritte (Checkout, Tag validieren, Versions-Abgleich, `make sdk-pack-kotlin`, Publish-Stufe bauen) stehen wortgleich in zwei Jobs; eine spätere Änderung der Tag-Validierung ist an zwei Stellen zu ziehen. Die ADR nimmt die Doppelung in Kauf (Alternative D1, „Bau und Tests je Job“). | [`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md) Festlegung 5 · Maintainability | `.github/workflows/sdk-kotlin-release.yml:73-91,110-128` | nein | Zwei Stellen für dieselbe Änderung |
| F-7 | INFO | Was der reale Tag-Lauf über das statisch Gelesene hinaus zeigen muss (lokal nicht ausführbar, `AGENTS.md` §3.10): (a) dass Gradle 8.x beim Task des einen Ziels die fehlenden Anmeldedaten des anderen Repositories nicht bemängelt (im Cloudsmith-Job sind `GITHUB_ACTOR`/`GITHUB_TOKEN` leer, im GitHub-Job die `CLOUDSMITH_*`; nur das Konfigurieren ohne alle vier ist über `make sdk-pack-kotlin` belegt, nicht das Ausführen eines Publish-Tasks); (b) dass `docker run -e GITHUB_ACTOR` (ohne Wert) den Standardwert des Runners erbt — der Vorgänger setzte ihn explizit über `${{ github.actor }}`; (c) Annahme der `.module`-Datei, Slug gegen Name, anonyme Sichtbarkeit und Doppel-Upload (alle vier im Plan §6 als offen geführt). | [`AGENTS.md`](../../AGENTS.md) §3.10 | `.github/workflows/sdk-kotlin-release.yml:93-101,130-139`, `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts:231-248` | nein — nur der Tag-Lauf `sdk-kotlin-v0.2.2` | Erster realer Lauf der Job-Struktur |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `.github/workflows/sdk-kotlin-release.yml` — Struktur | geprüft, ohne Befund: zwei Jobs (`sdk-kotlin-github-packages`, `sdk-kotlin-cloudsmith`), kein `needs`, kein `continue-on-error`, kein `if:`; ein roter Job lässt den anderen unberührt, „Re-run failed jobs“ trifft nur den roten (dieselbe `GITHUB_SHA`, Secrets frisch gelesen) |
| Workflow — Berechtigungen und Pinning | geprüft, ohne Befund: `permissions: {}` oben, `packages: write` genau einmal (Z. 71, GitHub-Packages-Job), Cloudsmith-Job nur `contents: read`; beide `uses:` `actions/checkout` mit SHA und Tag-Kommentar (§3.8), keine neue Action; Secret-Referenzen nur `secrets.GITHUB_TOKEN` und `secrets.CLOUDSMITH_*` |
| Workflow — Tag-Validierung und Versions-Abgleich | geprüft, ohne Befund: in beiden Jobs, vor `make sdk-pack-kotlin` und vor jedem Publish; der Regex `^version = "…"` trifft nur die Top-Level-Zeile (die Publikation ist eingerückt); `tag-info` druckt `version=0.2.2`; die Publikations-Version bindet die Probe (M4) |
| Workflow — Umgang mit Werten | geprüft, ohne Befund: Werte nur über `env:` des Schritts und `docker run -e NAME`, nie in der Kommandozeile, kein `${{ secrets… }}` im `run:`-Text (die Vorgängerfassung setzte `${{ … }}` in die Kommandozeile), kein `set -x`/`printenv`, zwei `echo "::error::…"` ohne Werte, kein `--info`/`--debug`, kein `ARG`/`--build-arg` |
| Workflow — Ausdrücke, Einrückung, Kommentare | geprüft, ohne Befund: `"on":` quotiert, Schrittnamen ohne `: `, Kopfkommentar im Indikativ (Klassen Zusage/Abgrenzung/Grenze, keine Chronik); Aussage „anonym lesbares“ siehe F-1 |
| Gradle-Task-Namen (`publishMavenPublicationTo<Repository>Repository`) | geprüft, ohne Befund: Publikation `maven`, Repositories `GitHubPackages`/`Cloudsmith`; M1/M2 färben die Probe, die aufgelösten Namen der `--dry-run`-Ausgabe werden ganzzeilig verglichen |
| `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` | geprüft, ohne Befund: Repositories, URLs, Anmeldedaten allein aus `System.getenv`, kein Wert im Code; `make sdk-pack-kotlin` läuft mit allen vier Variablen leer (Exit 0); Version `0.2.2` an beiden Stellen (Z. 86, 210) und im Kommentar (Z. 72); Kommentare im Indikativ ohne Chronik |
| `sdks/kotlin/Dockerfile` — Probe (a) bis (c) | geprüft, ohne Befund am Vertrag; vier Eingabeseiten-Mutationen plus GitHub-Packages-Name und README gesehen rot (M1 bis M6). Restpunkte als F-5 |
| Dockerfile — Stufe `publish`, Docker-only | geprüft, ohne Befund: kein `ARG` für die Namen, `CMD` bleibt Sammel-Aufgabe und scheitert ohne Cloudsmith-Werte (Kommentar sagt es), der Workflow überschreibt sie je Ziel; Publish läuft in der Stufe, nicht auf dem Runner (§3.1) |
| Kotlin-README (Anwender-Sicht) | geprüft, ohne Befund: Download-Basis `https://dl.cloudsmith.io/public/pt9912/pg-change-feed/maven/` wie [`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md) Festlegung 2, `repositories`-Snippet (Kotlin DSL, `dependencyResolutionManagement`) syntaktisch stimmig und mit `mavenCentral()`, Koordinate `io.github.pt9912:pgchangefeed-kotlin:0.2.2`, GitHub-Packages-Weg mit Token; Namensnennung erfüllt die Bedingungen wörtlich (Link `cloudsmith.com` und „free of charge for open-source projects“); Sprache rein englisch; `make sdk-public-doc-check` grün (M6 rot). Die Aussage „no account, no token“ siehe F-1 |
| `spec/pflichtenheft.md` | geprüft, ohne weiteren Befund: `LH-FA-SST-009.a` und `SPEC-028` präzisieren den Vertriebsweg (kein neuer bindender Anspruch), ohne ADR- und Slice-Bezug (§3.4), Historie-Zeile ohne Kennung des Slice; „vierte Sprache oder ein weiterer Vertriebsweg“ nachgezogen |
| `docs/user/benutzerhandbuch.md` | geprüft, ohne weiteren Befund: alle vier Kotlin-Stellen (HTTP, gRPC, SSE, NATS) nennen Cloudsmith zuerst und GitHub Packages mit Token; `Version:` 1.64 zu 1.65 und Historie-Zeile im selben Diff; kein verbliebener Satz mit alleiniger Token-Pflicht |
| `docs/user/releasing.md` außer F-2/F-3 | geprüft, ohne weiteren Befund: Version 1.10 samt Historie; Secret-Tabelle sechs Zeilen; „Alle Secrets dieser Tabelle“ statt Zählwort; Abschnitt „Offen bis zum ersten Tag-Lauf“ nennt POM-Abrufpfad, Moduldatei, Slug/Name, Paketseite; Wiederholung je Job; §1 nennt die Struktur unbewiesen |
| `harness/README.md`, `harness/mk/sdk.mk`, `tools/harness/sdk-pack-kotlin.sh` | geprüft, ohne Befund: Zeilen `sdk-kotlin-release.yml` und `make sdk-pack-kotlin` tragen den Ist-Zustand; „real erzeugt (`ls sdks/kotlin/dist`)“ mit `0.2.2` nachgemessen; Kommentare im Indikativ; Skriptlogik unverändert |
| Plan, Suchlauf-Feld (§3) | geprüft, ohne Befund: die Zeilen 1 bis 4, 6, 8, 9 an beiden Ständen nachgefahren, Zahlen wie gedruckt und mit Ursprung/Lauf; die Zeile „abgeleitet“ ist als abgeleitet gekennzeichnet; Lokatoren der Träger am Diff-Stand stimmen |
| Commit-Struktur und Traceability | geprüft, ohne Befund: die zwei Lifecycle-Moves (`open→next`, `next→in-progress`) sind reine Renames (0 Zeilen), jede Message nennt `LH-FA-SST-009` und `ADR-0123`, keine `SPEC-`/`ARC-`-Kennung im Betreff |
| Secret-Hygiene, Docker-only, host-lokale Pfade | geprüft, ohne Befund: kein Slug/Wert im Repo außerhalb des Plan-Rückfalltexts, keine Host-Toolchain in Skripten, `make docs-check` (hostpaths) grün |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Tatsachenbehauptung vor dem Beleg · Nachzug widerspricht dem Nachbarn ·
Verweis auf nicht vorhandenen Eintrag · Probe bindet nicht den Aufrufer · Zwei Stellen für dieselbe Änderung ·
Erster realer Lauf der Job-Struktur

## Verdikt

**Merge-blockierend:** nein. Kein HIGH; der Workflow, `build.gradle.kts`, die Probe und die
Secret-Behandlung tragen im Lesen und in den Läufen ohne Befund. Die beiden MEDIUM-Findings sind
Text-Findings in Anwender- und Betreiber-Doku: F-1 löst sich mit dem Post-Push-Beleg (§5), den der
Slice ohnehin offen hält (die Aussage ist vor dem ersten Upload nicht anders zu belegen; wer sie
vorher abschwächen will, formuliert sie als „erwartet“); F-2 ist eine Zeilen-Korrektur im selben
Träger. Ein Tag-Push `sdk-kotlin-v0.2.2` ist durch keines der Findings gehindert; der Beleg, dass
die Job-Struktur läuft, entsteht erst dort (F-7).

**Übergabe:** F-1 bis F-4 gehen an den Planner (Träger-Text, Plan §6) beziehungsweise, wenn er die
Fixrunde bevorzugt, an den Implementer; F-5 bis F-7 sind Hinweise ohne erwartete Aktion. Die
**Finding-Klassen** gehen zusätzlich in die Slice-Closure §7 und von dort in den Zähler. Die
DoD-Zeile „Review durchgeführt“ im Plan habe ich **nicht** gesetzt: der Auftrag ließ nur den Report
zu; ob eine Fixrunde für F-2 stattfindet, entscheidet der Aufrufer. Dieser Report ist ein
**Lauf-Beleg** und ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11).
