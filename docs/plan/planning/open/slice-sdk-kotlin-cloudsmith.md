# Slice sdk-kotlin-cloudsmith: Das Kotlin-SDK-Release veröffentlicht zusätzlich nach Cloudsmith — anonym beziehbares zweites Ziel, ein Job je Ziel, Version 0.2.2

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — der Slice trägt keine Closure-Bedingung, die von seiner
DoD verschieden wäre (Baseline-Regelwerk `modul-06-roadmap.md` §Wann Arbeit eine
Welle braucht). Er hat keine Kante zu einer Welle; die Paarungen prüft die
Closure von [welle-transformationen](../welle-transformationen.md) (§5
„Ereignis-Adresse“).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (offizielle
Client-Bibliotheken; Akzeptanz: über den Paketmanager der Sprache einbinden),
[`ADR-0123`](../../adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md)
(Kotlin-SDK zusätzlich auf Cloudsmith; löst zwei Secret-Aussagen von
[`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
teilweise ab, der Rest von `ADR-0109` gilt),
[`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) (Registry-Secret-Muster:
der Betreiber legt das Secret manuell an), Architect-Verdikt
[`architect-verdict-sdk-kotlin-cloudsmith`](../../../reviews/architect-verdict-sdk-kotlin-cloudsmith.md)
(Zuschnitt, Betreiber-Voraussetzungen, Folge-Arbeit),
[`AGENTS.md`](../../../../AGENTS.md) §3.1, §3.8, §3.10, §3.12, §3.13.

**Berührte Spec-Stellen:** [`SPEC-028`](../../../../spec/pflichtenheft.md)
(Package-Zeile: Vertriebsziele und Version) und §1
[`LH-FA-SST-009.a`](../../../../spec/pflichtenheft.md) (Kotlin-Absatz) —
geändert; [`SPEC-018`](../../../../spec/pflichtenheft.md),
[`SPEC-020`](../../../../spec/pflichtenheft.md),
[`SPEC-021`](../../../../spec/pflichtenheft.md),
[`SPEC-022`](../../../../spec/pflichtenheft.md) und
[`SPEC-024`](../../../../spec/pflichtenheft.md) (die Drahtverträge, die das
Package deckt) — unberührt.

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Architect-Verdikt
[`architect-verdict-sdk-kotlin-cloudsmith`](../../../reviews/architect-verdict-sdk-kotlin-cloudsmith.md)
(Nutzerentscheidung 2026-09-25: Cloudsmith nur für Kotlin). **Datum:** 2026-09-25.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein `sdk-kotlin-v<SemVer>`-Tag veröffentlicht das Package
`io.github.pt9912:pgchangefeed-kotlin` zusätzlich in das Open-Source-Repository
`pt9912/pg-change-feed` auf Cloudsmith — ein Job je Ziel, Docker-only in der
Stufe `publish`, ohne Action —, sodass Kotlin-Anwender die Bibliothek ohne
GitHub-Token beziehen können; die Kotlin-README nennt Cloudsmith als ersten
Bezugsweg, und die Quell-Version steht auf `0.2.2`.

**Ausgangslage** (gelesen am Stand `37e825e2`, 2026-09-25):

- `.github/workflows/sdk-kotlin-release.yml` trägt **einen** Job
  (`sdk-kotlin-release`), der nach `make sdk-pack-kotlin` die Stufe `publish` baut
  und darin `./gradlew --no-daemon publish` mit `GITHUB_TOKEN` startet
  (Datei, 128 Zeilen laut `wc -l`). Der `publishing`-Block der `build.gradle.kts`
  trägt **ein** Repository (`GitHubPackages`).
- GitHub Packages verlangt zum Lesen ein GitHub-Konto und einen Token mit
  `read:packages` (Kotlin-README, Abschnitt „Installation“; Handbuch, vier Stellen
  laut Suchlauf §3). Das Package hat damit eine Hürde, die NuGet.org (C#) und PyPI
  (Python) nicht haben — der Anlass von `ADR-0123`.
- Die Tags `sdk-kotlin-v0.2.0` und `sdk-kotlin-v0.2.1` existieren
  (`git ls-remote --tags origin 'sdk-kotlin*'` druckte am 2026-09-25 beide, auf
  `eb3ba91e` und `7df775c8`). Ein bestehender Tag löst den Workflow nicht erneut
  aus; Cloudsmith bekommt deshalb die Version des **nächsten** Tags.

**Versionsentscheidung (der Slice hebt auf `0.2.2`).** Eine veröffentlichte
Maven-Version ist unveränderlich (`ADR-0123` Festlegung 2 und 5); ein Release nach
Cloudsmith braucht einen neuen Tag und damit eine noch nicht veröffentlichte
Quell-Version, deren Tag-Suffix der Workflow gegen die Versionszeile abgleicht.
`0.2.2` ist die nächste PATCH-Version der `0.x`-Reihe (SemVer,
[`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
Festlegung 4): Der Slice ändert keine Signatur und kein Verhalten der Bibliothek,
nur Veröffentlichungsweg, README und Metadaten. Gehoben wird an **beiden** Stellen
der `build.gradle.kts` — der Top-Level-`version`-Zeile, die der Workflow gegen den
Tag abgleicht, und der `version` der Publikation, die das Artefakt und die POM
trägt (Beleg: `grep -n 'version = "0.2.1"'` druckt am Stand `37e825e2` die Zeilen
79 und 203) — und in allen Trägern (§3, Suchlauf). Dasselbe Vorgehen trug der
Vorgänger-Slice `slice-sdk-readme-nutzerdoku` (`0.2.0` auf `0.2.1`, damit die neue
Beschreibung mit dem nächsten Tag erscheint). Ist beim Start bereits ein Tag
`sdk-kotlin-v0.2.2` gesetzt (Prüfung §4), trägt der Slice die nächste freie Version.
**Der Tag-Push ist eine Betreiber-Handlung und nicht Teil des Slices**
([`AGENTS.md`](../../../../AGENTS.md) §3.10, `ADR-0123` Betreiber-Voraussetzung 6);
die Hebung ist nur in Verbindung mit ihm sinnvoll, ihr Beleg steht in §5.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der Tag-Push und der reale Release-Lauf** — externe Handlung mit Repository-Secrets
  und Cloudsmith-Konto; ein Agent führt sie nicht aus. Die Adresse für den Beleg ist
  die Closure dieses Slice (§5), nicht ein Folge-Slice.
- **Cloudsmith für C# und Python** — NuGet.org und PyPI sind anonym lesbar, ein zweites
  Ziel schafft dort keinen Nutzen (`ADR-0123` Festlegung 1, Alternative C3);
  `sdks/csharp/**`, `sdks/python/**` und ihre Workflows bleiben unberührt.
- **Maven Central, OIDC, Signatur, Javadoc-Jar** — eigene Entscheidungen mit eigenem
  Aufwand (`ADR-0123` Festlegung 7, Re-Evaluierungs-Trigger 1 und 3).
- **Abschaltung von GitHub Packages, Änderung des Tag-Namensraums oder der
  Tag-Validierung** — `ADR-0123` Festlegung 7: beides bleibt.
- **Anlage und Pflege des Cloudsmith-Kontos, -Repositories, -Services und der
  Repository-Secrets** — externe, kontobezogene Betreiber-Handlung; der Slice
  prüft ihr Ergebnis (§4 Startbedingung), legt es nicht an.
- **Ein neues Gate in `make gates`** — die Probe (§3) hängt an `make sdk-pack-kotlin`,
  einem Werkzeug wie die übrigen `sdk-pack-*`-Ziele (Netzbezug des Gradle-Baus);
  keine Schwelle wird gesenkt ([`AGENTS.md`](../../../../AGENTS.md) §3.6).
- **Ein dauerhafter Wächter für den Text der Kotlin-README-Beispiele** — anderer Vorgang
  (Änderung der Test-Bau-Kontexte, `slice-sdk-readme-nutzerdoku` §1); der Bezug der
  README auf den Download-Pfad wird stattdessen im Post-Push-Beleg (§5) gegen den
  realen Abruf gehalten.
- **Server-Version, Server-Code und Handbuch-Abschnitte außerhalb des Kotlin-Bezugs** —
  `docs/user/version.md` und `internal/**` bleiben unberührt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die Closure-Pflichten darunter zählen nicht mit. Die Punkte, die
erst nach der Arbeit belegbar sind, stehen als **Zusage** („zu belegen durch …“),
nicht als Feststellung ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B).

Liefer-Punkte (drei):

- [ ] **1 — Publish-Kette mit zwei Zielen.** `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts`
      trägt ein zweites Repository `Cloudsmith` (Upload-URL
      `https://maven.cloudsmith.io/pt9912/pg-change-feed/`, Anmeldedaten allein aus
      `CLOUDSMITH_USERNAME`/`CLOUDSMITH_API_KEY` der Umgebung, kein Wert im Code) und
      die Version `0.2.2` an beiden Stellen; `.github/workflows/sdk-kotlin-release.yml`
      führt je Ziel einen eigenen Job mit der Einzel-Aufgabe des Ziels (nicht die
      Sammel-Aufgabe `publish`); eine Probe im Bau von `make sdk-pack-kotlin`, ohne
      Zugangsdaten und ohne Zugriff auf die Ziele, färbt rot, wenn eine der beiden
      Aufgaben oder die POM-Koordinate fehlt. Zu belegen durch: die
      Mutations-Tabelle §3 (je Zusage die Eingabeseiten-Mutation mit gesehenem Rot).
- [ ] **2 — Anwender-Sicht.** `sdks/kotlin/pgchangefeed-kotlin/README.md` nennt Cloudsmith
      als ersten Bezugsweg ohne Token (Repository-URL, `repositories`-Snippet,
      Koordinate mit `0.2.2`), GitHub Packages als weiteren Weg mit Token, und trägt die
      von Cloudsmith verlangte Namensnennung (Text mit Link auf `cloudsmith.com` und
      dem Hinweis auf kostenloses Hosting, Anker: `ADR-0123` §Kontext, Zeile
      „Regeln“); die README steht in Anwender-Sprache ohne interne Kennung
      (`make sdk-public-doc-check` grün). `docs/user/benutzerhandbuch.md` beschreibt an
      den vier Kotlin-Stellen (§3) den Bezug ohne Token; `Version:` und
      Änderungshistorie sind im selben Diff nachgezogen.
- [ ] **3 — Träger-Nachzug.** `spec/pflichtenheft.md` (`LH-FA-SST-009.a`, Kotlin-Absatz,
      und `SPEC-028`-Zeile: zwei Vertriebsziele, Version `0.2.2`),
      `harness/README.md` (Zeile `sdk-kotlin-release.yml`, Zeile `make sdk-pack-kotlin`),
      `docs/user/releasing.md` §4 (Secret-Tabelle, Ablauf mit zwei Jobs, Absatz zur
      Paketseite) und die Kommentare in `sdks/kotlin/Dockerfile`,
      `harness/mk/sdk.mk`, `tools/harness/sdk-pack-kotlin.sh` tragen den Ist-Zustand;
      der Suchlauf (§3) steht mit beiden Ständen und dem, was er **nicht** fand, im
      Plan.

Gates und Belege:

- [ ] `make gates`, `make sdk-pack-kotlin`, `make sdk-public-doc-check` und
      `make test-sdk-kotlin-release-tag-info` grün; jeder Exit-Code direkt gelesen
      ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Post-Push-Beleg ([`AGENTS.md`](../../../../AGENTS.md) §3.10 — der Workflow ändert
      seine Job-Struktur): der Betreiber setzt den Tag `sdk-kotlin-v0.2.2` (oder die
      nächste freie Version), **beide** Jobs enden `success`, und ein **anonymer**
      Abruf der POM-Datei am Download-Pfad des Repositories antwortet HTTP 200
      (erwartet, nach dem Sammelfenster und der asynchronen Verarbeitung ggf. mit
      Wiederholung). Der Beleg wird vor der Closure in §6/§7 eingetragen; bis dahin
      bleibt das Risiko „erster realer Lauf“ offen.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **mit** Wellen von der nächsten Welle-Closure geprüft (auch für Slices ohne Wellen-Zugehörigkeit); Adresse in §4.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

Reihenfolge im Lauf: erst die Spec (die Spec führt), dann Build und Workflow samt
Probe, dann README und Handbuch, dann die übrigen Träger; der Implementer trägt
jede Abweichung vor dem Sensor-Lauf in diese Tabelle nach.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/pflichtenheft.md` §1 `LH-FA-SST-009.a` (Kotlin-Absatz, Z. 312–316 am Stand `37e825e2`), §6 `SPEC-028`-Zeile (Z. 822), §7 (eine Historie-Zeile) | update | Der Kotlin-Absatz und die `SPEC-028`-Zeile nennen **zwei** Vertriebsziele (GitHub Packages, Cloudsmith) und die Version `0.2.2`; der Satz „ein vierter Vertriebsweg bleibt offen“ (Z. 316) wird auf den neuen Stand gezogen. Die Spec-Aussage trägt **keinen** ADR- und keinen Slice-Bezug (`ADR-0123` nennt die Spec in `Schärft:`). Das Lastenheft bleibt unberührt (`LH-FA-SST-009` verlangt den Bezug über den Paketmanager, nicht einen Vertriebsweg). |
| `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` | update | Zweites `maven`-Repository `Cloudsmith` im `publishing`-Block (`ADR-0123` Festlegung 4), Anmeldedaten aus `System.getenv("CLOUDSMITH_USERNAME")`/`System.getenv("CLOUDSMITH_API_KEY")`; `version` `0.2.2` an drei Stellen (Kommentar „Version `0.2.1`“ Z. 67, Top-Level Z. 79, Publikation Z. 203); Kopfkommentar (Z. 9–25, sagt heute „no repository secret is needed“) beschreibt beide Ziele und beide Umgebungs-Paare. Plugin bleibt `maven-publish`. |
| `.github/workflows/sdk-kotlin-release.yml` | update | Zwei Jobs statt einem (Form: `permissions` ist je Job festgelegt; ein Matrix-Zweig nur, wenn der Implementer belegt, dass die Berechtigungen je Zweig verschieden sein können): `sdk-kotlin-github-packages` (`contents: read`, `packages: write`) und `sdk-kotlin-cloudsmith` (`contents: read`), ohne `needs` und ohne `continue-on-error`. Je Job: Checkout (der vorhandene SHA-gepinnte Schritt, keine neue `uses:`-Zeile), Tag validieren, `build.gradle.kts` gegen den Tag abgleichen, `make sdk-pack-kotlin`, Publish-Stufe bauen, Publish-Lauf mit der **Einzel-Aufgabe** des Ziels. Cloudsmith-Werte über `env:` des Schritts und `docker run -e NAME` (ohne Wert in der Kommandozeile), nie als Build-Argument. Kopfkommentar (Z. 1–83; Z. 59–64 sagen „AUSSCHLIESSLICH das eingebaute `GITHUB_TOKEN`“ und „KEIN externes Repository-Secret“) beschreibt die zwei Ziele. |
| `sdks/kotlin/Dockerfile` | update | Probe-`RUN` in der Stufe `pack` (siehe Ansatz unten); Kommentare der Stufe `publish` (Z. 83–107 nennen nur GitHub Packages und `GITHUB_ACTOR`/`GITHUB_TOKEN`) tragen beide Ziele; die `CMD` der Stufe bleibt (`ADR-0123` Festlegung 4: der Workflow überschreibt sie ausdrücklich), ihr Kommentar sagt, dass die Sammel-Aufgabe ohne Cloudsmith-Werte scheitert; Artefaktnamen `0.2.1` in Z. 66, 67, 88 auf `0.2.2`. |
| `sdks/kotlin/pgchangefeed-kotlin/README.md` | update | Abschnitt „Installation“: Cloudsmith zuerst (Repository-URL `https://dl.cloudsmith.io/public/pt9912/pg-change-feed/maven/`, `settings.gradle.kts`-Snippet mit `mavenCentral()` **und** dem Cloudsmith-Repository ohne Zugangsdaten — die Abhängigkeiten des Packages löst der Anwender über Maven Central auf, `ADR-0123` Festlegung 6 —, Koordinate `io.github.pt9912:pgchangefeed-kotlin:0.2.2`), GitHub Packages als weiterer Weg mit Token-Hinweis (der heutige Absatz und das Snippet), Namensnennung (Text mit Link auf `https://cloudsmith.com` und Hinweis auf kostenloses Hosting für Open-Source-Projekte). Anwender-Sprache, **keine** Kennung `SPEC-`/`ADR-`/`LH-`/`ARC-`, kein Slice-/Welle-Name, keine Chronik. |
| `docs/user/benutzerhandbuch.md` | update | Die vier Kotlin-Stellen, die den Bezug über GitHub Packages mit Token-Pflicht beschreiben (Z. 1049–1055, 1148–1153, 1225–1227, 1379–1382 am Stand `37e825e2`): der Bezug ohne Token über Cloudsmith steht zuerst, GitHub Packages als Weg mit Token. `Version:` (Stand `37e825e2`: `1.59`) und eine Zeile in `### Änderungshistorie` im selben Diff (`.claude/commands/implement-slice.md` Schritt 17). Die Zeilenangaben sind Lokatoren vom Stand `37e825e2` und verschieben sich mit jeder Nachbaränderung ([`AGENTS.md`](../../../../AGENTS.md) §3.13 Grenze). Beim Start `Version:` erneut lesen: ein anderer Slice kann sie vorher gehoben haben. |
| `docs/user/releasing.md` | update | `Version:` `1.9` samt Historie-Zeile; §1 (Aufzählung der Release-Wege, Z. 21–28); §4 Kotlin-Abschnitt: Überschrift, „vierter, eigenständiger Release-Mechanismus“ (Z. 223), Ablauf mit zwei Jobs und Einzel-Aufgaben, der Absatz „Kein externes Repository-Secret nötig“ (Z. 266–277) und der Absatz „Ein realer Preis bleibt trotzdem bestehen“ (Z. 278–286); Secret-Tabelle (Z. 103–108) um `CLOUDSMITH_USERNAME` und `CLOUDSMITH_API_KEY`; der Satz „Alle drei Secrets sind eine externe, kontobezogene Handlung“ (Z. 120) nennt heute eine Zahl, die nicht zur Tabelle passt (vier Zeilen) — die Zahl wird aus der Tabelle abgeleitet oder entfällt; ein neuer Absatz zur Cloudsmith-Paketseite (ob die POM-Beschreibung dort angezeigt wird: nicht geprüft, Betreiber pflegt optional den Text in der Web-App); der Absatz zur GitHub-Paketseite bleibt. Der erste reale Lauf bleibt als **offen** markiert, bis §5 erfüllt ist. |
| `harness/README.md` (§Sensors) | update | Zeile `.github/workflows/sdk-kotlin-release.yml` (zwei Ziele, zwei Jobs, Cloudsmith-Secrets, Einzel-Aufgaben, Beleg-Stand); Zeile `make sdk-pack-kotlin` (Probe, Artefaktnamen `0.2.2`). Die Aufnahme geschieht erst, wenn der Workflow im Diff real existiert ([`AGENTS.md`](../../../../AGENTS.md) §4). |
| `harness/mk/sdk.mk`, `tools/harness/sdk-pack-kotlin.sh`, `tools/harness/run-sdk-kotlin-release-tag-info-tests.sh` | update (nur Kommentare, falls der Suchlauf trifft) | Die Kommentare nennen GitHub Packages als einziges Ziel bzw. den Artefaktnamen `0.2.1` (`sdk.mk` Z. 75; `sdk-pack-kotlin.sh` Z. 34; der Tabellentest Z. 6, dort als Abgrenzung „kein GitHub-Packages-Äquivalent“ zu `:latest` — bleibt richtig, der Implementer prüft und belässt). Die Skriptlogik bleibt unverändert. |
| `docs/plan/planning/welle-transformationen.md` | update (Planner-Zug dieser Planung) | Der Satz „Package-Versionen der SDKs bleiben unberührt“ (§6) nennt die Slices **dieser Welle**; er trägt jetzt den Zusatz, dass die Kotlin-Version durch den wellenlosen Slice `slice-sdk-kotlin-cloudsmith` gehoben wird, damit die Closure der Welle keine Abweichung liest. |

**Ansatz der Probe im Bau (Zusage, im Slice zu erproben — die Aussagen über das
Gradle-Verhalten sind hergeleitet, nicht gemessen):**

- Ort: ein `RUN` in der Stufe `pack` des Dockerfiles. `make sdk-pack-kotlin` baut
  `--target pack-export` (`tools/harness/sdk-pack-kotlin.sh`), und `pack-export`
  baut auf `pack` auf (`sdks/kotlin/Dockerfile`, Zeilen `FROM build AS pack`,
  `FROM pack AS pack-export`); ein roter Probe-Schritt bricht damit den Bau ab, bevor
  ein Artefakt exportiert wird — und läuft in **beiden** Release-Jobs vor dem Publish,
  weil jeder Job `make sdk-pack-kotlin` ruft.
- Prüfung (a): beide Einzel-Aufgaben existieren — `./gradlew --no-daemon --dry-run
  publishMavenPublicationToGitHubPackagesRepository
  publishMavenPublicationToCloudsmithRepository`. Erwartet: `--dry-run` führt keine
  Aufgabe aus und endet mit Exit ≠ 0 bei einem unbekannten Aufgabennamen. Sind die
  Namen anders (die Gradle-Namensregel ist in `ADR-0123` Festlegung 4 hergeleitet),
  wählt der Implementer die belegt existierende Form und trägt sie in Workflow, Probe
  und Doku gleich ein.
- Prüfung (b): die generierte POM trägt die Koordinate — `./gradlew --no-daemon
  generatePomFileForMavenPublication`, danach `groupId` `io.github.pt9912`,
  `artifactId` `pgchangefeed-kotlin` und `version` gleich dem Versionsteil des Jar-Namens
  unter `build/libs/` (das Jar trägt die Top-Level-`version`, die POM die `version` der
  Publikation; die Probe bindet damit **beide** Versions-Stellen aneinander). Der Befehl
  ist im Bau-Image bereits gelaufen (`slice-sdk-readme-nutzerdoku` §3, Zeile „Kotlin ·
  Installation“).
- Prüfung (c, optional): die Upload-URL des Cloudsmith-Repositories stimmt mit dem
  Literal `https://maven.cloudsmith.io/pt9912/pg-change-feed/` überein (ein `grep`
  gegen `build.gradle.kts`); der Implementer entscheidet und begründet, welche
  Zusagen die Probe bindet.
- Rahmen: keine Zugangsdaten und kein Zugriff auf GitHub Packages oder Cloudsmith. Die
  Konfiguration liest die beiden Umgebungs-Paare mit `System.getenv`; sind sie leer,
  konfiguriert der Bau trotzdem (Beleg: `make sdk-pack-kotlin` endete im Vorgänger-Slice
  `slice-sdk-readme-nutzerdoku` mit `BUILD SUCCESSFUL in 28s` ohne diese Variablen —
  gedruckte Zeile aus dessen §3; dass der `publishing`-Block dabei ausgewertet wird, ist
  aus dem Gradle-Verhalten hergeleitet und zu erproben).
- Kommentare der Probe stehen im Ist-Zustand ohne Kennung ([`AGENTS.md`](../../../../AGENTS.md)
  §3.7; `make sdk-public-doc-check` scannt das Dockerfile).

**Eingabeseiten-Mutationen je Zusage** (`.claude/commands/implement-slice.md`
Schritt 19; Erwartung des Plans — der Implementer trägt das **gesehene** Rot ein oder
den Grund, warum die Zeile leer bleibt):

| Zusage | mutierte Eingabe | erwartetes Rot (gesehenes Rot: Implementer) |
|---|---|---|
| Die Probe erkennt ein fehlendes Cloudsmith-Repository | `name = "Cloudsmith"` in `build.gradle.kts` zu einem anderen Namen | `make sdk-pack-kotlin` Exit ≠ 0 im Probe-Schritt (Aufgabe unbekannt) |
| Die Probe erkennt ein fehlendes GitHub-Packages-Repository | `name = "GitHubPackages"` zu einem anderen Namen | derselbe Schritt rot |
| Die POM trägt die Koordinate | `artifactId` der Publikation ändern | Prüfung (b) rot |
| Beide Versions-Stellen stimmen überein | Top-Level-`version` auf `0.2.2` lassen, `version` der Publikation auf `0.2.1` zurücksetzen (und umgekehrt) | Prüfung (b) rot; der Workflow-Abgleich gegen den Tag liest nur die Top-Level-Zeile und bliebe grün — diese Mutation ist die Eingabeseite, an der die Probe den Unterschied trägt |
| Die README trägt keine interne Kennung | `ADR-0123` in `sdks/kotlin/pgchangefeed-kotlin/README.md` einfügen | `make sdk-public-doc-check` Exit 1 (Vorstufe von `make sdk-pack-kotlin`) |
| Die Upload-URL des Repositories stimmt (nur wenn Prüfung (c) gebaut wird) | Slug `pg-change-feed` in der URL ändern | Prüfung (c) rot |
| Die Zugangsdaten stehen nur in der Umgebung, nie im Code, nie als Build-Argument | — (leer) | Kein Sensor: `ADR-0123` Fitness Function führt beides als Review-Punkt. Belegbefehle im Suchlauf unten (`git grep` nach `^ARG |--build-arg` und nach den beiden Namen) |
| Der Cloudsmith-Job trägt kein `packages: write` | — (leer) | Kein Sensor (Workflow lokal nur statisch prüfbar, [`AGENTS.md`](../../../../AGENTS.md) §3.10); Belegbefehl: `git grep -n 'packages: write'` druckt genau eine Zeile im GitHub-Packages-Job; die Wirkung belegt der Post-Push-Lauf |

**§3.13-Suchlauf** (committetes Feld; bewegte Eigenschaften: „Secret-Regime und
Vertriebsziel des Kotlin-Packages“ und „Version des Kotlin-Packages“; Spalte
*Parent* am Stand `37e825e2` gemessen, Zahlen wie von `git grep` gedruckt; Spalte
*Diff-Stand* trägt der Implementer am Ende ein, mit Zählwort und dem, was **nicht**
gefunden wurde):

| Träger / Eigenschaft | Suchbefehl | Parent (`37e825e2`) | Diff-Stand | Behandlung |
|---|---|---|---|---|
| Secret-Regime-Aussagen („kein externes Secret“, `GITHUB_TOKEN`, Token-Pflicht) | `git grep -c -E 'GitHub Packages\|GITHUB_TOKEN\|kein externes\|Personal-Access-Token\|read:packages' <Stand> -- spec/pflichtenheft.md harness/README.md docs/user/releasing.md docs/user/benutzerhandbuch.md sdks .github/workflows/sdk-kotlin-release.yml docs/plan/planning/in-progress/roadmap.md README.md` | **51** Zeilen in 9 Dateien: Workflow 10, `releasing.md` 14, Handbuch 10, `build.gradle.kts` 6, Dockerfile 4, Kotlin-README 2, `pflichtenheft.md` 2, `roadmap.md` 2, `harness/README.md` 1 (Summe abgeleitet) | *(Implementer)* | Träger nachgezogen laut Tabelle oben; Records bleiben: Historie-Zeilen von Handbuch und `releasing.md`, Wellenname in der Roadmap |
| Vertriebsziel „GitHub Packages“ als einziges Ziel, weiter gefasst (ohne ADR, `done/`, Reviews, Register) | `git grep -c -E 'GitHub Packages\|GitHub-Packages' <Stand> -- . ':!docs/plan/adr' ':!docs/plan/planning/done' ':!docs/reviews' ':!docs/plan/planning/observations'` | **36** Zeilen in 12 Dateien (Workflow 4, Roadmap 3, Handbuch 8, `releasing.md` 9, `harness/README.md` 1, `sdk.mk` 2, Dockerfile 1, Kotlin-README 1, `build.gradle.kts` 1, `pflichtenheft.md` 4, `run-sdk-kotlin-release-tag-info-tests.sh` 1, `sdk-pack-kotlin.sh` 1; Summe abgeleitet) | *(Implementer)* | wie oben; die zwei Skript-Kommentare prüfen |
| Version `0.2.1` der Kotlin-Träger | `git grep -n '0\.2\.1' <Stand> -- sdks/kotlin harness spec docs/user tools .github Makefile` gefiltert auf Kotlin | **11** Kotlin-Zeilen in 6 Dateien: Kotlin-README Z. 40; `build.gradle.kts` Z. 67, 79, 203; Dockerfile Z. 66, 67, 88; `sdk.mk` Z. 75; `harness/README.md` Z. 160; `pflichtenheft.md` Z. 315, 822 (Summe abgeleitet); dazu die Records Z. 865 und Z. 875 (Historie) | *(Implementer)* | alle elf auf `0.2.2`; Records bleiben. Ergebnis der C#-/Python-Zeilen: unberührt (`0.2.1` bleibt dort) |
| Job-Kennung `sdk-kotlin-release` (Umbenennung der Job-Id) | `git grep -n 'sdk-kotlin-release' <Stand> -- . ':!docs/plan/planning/done' ':!docs/reviews' ':!docs/plan/adr' ':!docs/plan/planning/observations'` | die Job-Id steht **einmal** (`.github/workflows/sdk-kotlin-release.yml` Z. 94); alle anderen Treffer sind der Dateiname des Workflows oder Skriptnamen (`sdk-kotlin-release-tag-info…`) | *(Implementer)* | keine weiteren Träger der Job-Id (Ergebnis dieser Planung); der Implementer bestätigt am Diff-Stand |
| Zählwort der Secret-Tabelle | `git grep -n 'Alle drei Secrets' <Stand> -- docs/user/releasing.md` und die Zeilen der Tabelle „Benötigte Repository-Secrets“ | 1 Treffer (Z. 120); Tabelle trägt vier Zeilen (`DOCKERHUB_USERNAME`, `DOCKERHUB_TOKEN`, `NUGET_API_KEY`, `PYPI_API_TOKEN`) — das Zählwort war schon am Parent nicht die Zeilenzahl | *(Implementer)* | Zahl aus der Tabelle ableiten oder das Zählwort streichen |
| Offene Sätze zum Vertriebsweg | `git grep -n -E 'vierter Vertriebsweg\|vierte Sprache' <Stand> -- spec` | Z. 315–316 (Kotlin-Absatz, **Träger**), Z. 863 und 865 (Historie, Records) | *(Implementer)* | Z. 316 nachziehen; der C#- und der Python-Absatz führen eigene „bleibt offen“-Sätze über **ihre** Sprache und bleiben unberührt |
| Aussagen über die SDK-Versionen und Kotlin in offenen Plänen | `git grep -n -E 'sdk-kotlin\|GitHub Packages\|GitHub-Packages\|pgchangefeed-kotlin\|Cloudsmith\|Kotlin-SDK\|Kotlin/GitHub' -- docs/plan/planning/open docs/plan/planning/welle-transformationen.md docs/plan/planning/next` | Treffer: `slice-transformationen-betriebsdoku` Z. 73 und 137 (liest die Kotlin-README auf Row-Image-Aussagen — kein Widerspruch, sieht aber den neuen Installationsabschnitt); `welle-transformationen.md` §6 „Package-Versionen der SDKs bleiben unberührt“ (Z. 394–395); sonst keine | *(Implementer)* | die Welle-Zeile ist in dieser Planung nachgezogen; `slice-transformationen-betriebsdoku` bleibt, sein Suchlauf am Start liest die dann neue README |
| Zugangsdaten nicht als Build-Argument | `git grep -n -E '^ARG \|--build-arg' <Stand> -- .github/workflows/sdk-kotlin-release.yml sdks/kotlin/Dockerfile` | **0** Zeilen | *(Implementer; erwartet 0)* | Belegbefehl der Review-Zusage |
| Kein Wert in Logs | `git grep -n -E '(^\|[^a-z])(echo\|printenv\|set -x)( \|$)' <Stand> -- .github/workflows/sdk-kotlin-release.yml` | **1** Zeile (Z. 112, `echo "::error::…"` mit Pfad und Versionen, ohne Zugangsdaten) | *(Implementer; erwartet 2, je Job eine, ohne Zugangsdaten)* | die Fehlerausgabe des Versions-Abgleichs steht je Job; kein `echo`, `printenv`, `set -x` auf Secret-Werte, kein `--info`/`--debug` beim Gradle-Aufruf |

**Was der Suchlauf nicht fand** (Stand der Planung): kein Träger außerhalb der Tabelle
oben, der `GitHub Packages` als Vertriebsweg des Packages beschreibt — die Ausschlüsse
des Suchraums (`docs/plan/adr`, `done/`, `docs/reviews`, Register) sind Records oder
unveränderliche ADR. Ob die Zahlen `36` und `51` am Diff-Stand sinken, ist eine
Messung des Implementers, kein Ziel: manche Treffer bleiben richtig (GitHub Packages
bleibt Ziel).

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`), drei beobachtbare Bedingungen:

1. **WIP-Limit 1:** `docs/plan/planning/done/slice-backfill-speicher-untersuchung.md`
   existiert, und `in-progress/` trägt nur `roadmap.md`. Stand der Planung: der Slice
   liegt in `in-progress/`, die Datei in `done/` fehlt (`ls`, 2026-09-25) — die
   Bedingung ist **nicht** erfüllt.
2. **Betreiber-Voraussetzungen** (`ADR-0123` Folgepflicht 3; Verdikt §Betreiber-Voraussetzungen
   1 bis 5) — Ergebnis der Prüfung dieser Planung, 2026-09-25:
   - Secrets: `gh secret list` druckte die Zeilen
     `CLOUDSMITH_API_KEY	2026-09-25T13:48:36Z` und
     `CLOUDSMITH_USERNAME	2026-09-25T13:48:25Z` — beide gesetzt (nur Namen und
     Zeitstempel, kein Wert; `gh secret list` zeigt Werte nie). **Erfüllt.**
   - Namen: Organisation `pt9912`, Repository `pg-change-feed`, Service mit
     Schreibrecht — **übernommen** aus der Mitteilung des Betreibers, in der Web-App
     nicht eingesehen. Der Wert von `CLOUDSMITH_USERNAME` ist nach der Mitteilung des
     Betreibers der Service-Slug; der Wert steht in keinem Dokument dieses Repos.
   - Anonyme Sichtbarkeit des Repositories — **nicht belegt**. Gemessen (`curl -s`,
     lesend, ohne Zugangsdaten): `https://api.cloudsmith.io/v1/repos/pt9912/` antwortet
     `[]` mit HTTP 200; `…/v1/repos/pt9912/pg-change-feed/` antwortet
     `{"detail":"Not found."}` mit HTTP 404; `https://dl.cloudsmith.io/public/pt9912/pg-change-feed/maven/`
     antwortet HTTP 404; ein nicht existierender Namespace
     (`…/v1/repos/nonexistent-zz9/none/`) antwortet ebenfalls HTTP 404. Gelesen
     (hergeleitet): der Namespace `pt9912` ist anonym auflösbar, ein öffentliches
     Repository ist anonym nicht sichtbar. Entweder ist das Repository (noch) nicht
     als „Open source“/Broadcast sichtbar, oder ein leeres öffentliches Repository
     erscheint anonym nicht — die Messung trennt beides nicht. **Der Slice-Bau,
     die Doku und die Probe hängen davon nicht ab; der Beleg in §5 und die
     Aussage der README („ohne Token beziehbar“) hängen davon ab.** Der Implementer
     misst die drei Abrufe am Start erneut; ist das Ergebnis unverändert, meldet er
     es dem Betreiber (Frage: Sichtbarkeit und Slug des Repositories in der
     Web-App) und arbeitet weiter.
3. **Version frei:** `git ls-remote --tags origin 'sdk-kotlin-v0.2.2'` druckt keine
   Zeile (Stand der Planung: 0 Zeilen); sonst trägt der Slice die nächste freie
   PATCH-Version durch alle Träger (§1 Versionsentscheidung).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wenn die Probe mehr verlangt
  als ein `RUN` in der Stufe `pack` (eine neue Datei im Bau-Kontext, eine Änderung an
  `tools/harness/sdk-pack-kotlin.sh`), oder wenn die Träger-Tabelle §3 um Dateien wächst,
  die der Suchlauf nicht erklärt.
- `in-progress` → `open` (blockiert): wenn eine Betreiber-Voraussetzung fehlt und der
  Betreiber sie nicht liefert (Secret, Repository-Sichtbarkeit) und der Bau daran
  scheitert. Die Begründung steht im Move-Commit.

**Ereignis-Adresse der drei Paarungen** (Anker · Folge-Slice · Register): die Closure
der Welle [welle-transformationen](../welle-transformationen.md). Die Adresse kann
eintreten: die Roadmap führt die Welle unter *Offene Wellen*
(`docs/plan/planning/in-progress/roadmap.md` Z. 41, gelesen 2026-09-25), und die
Welle-Datei hat einen Closure-Trigger (§3 „Closure-Trigger (Welle schließt)“; Regel
`BEO-PGC/aufschub-adresse-verfaellt`, 3×, verkörpert: ein Ereignis-Träger ohne
gesichertes Eintreten ist keine Adresse). Ist die Welle bei der Closure dieses Slice
bereits geschlossen (Datei in `done/`), prüft die Closure dieses Slice die drei
Paarungen selbst und nennt das in §7. **Roadmap:** keine Zeile — die Roadmap führt einen
wellenlosen Slice nur mit Kante zu einer Welle (Drift-Log Z. 276: zwei Kanten des
`welle-transformationen`); dieser Slice trägt keine.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Beobachtbar: (1) die DoD (§2) ist vollständig, `make gates` grün, der Review-Report
liegt vor und ist aufgelöst; (2) der **Post-Push-Beleg** aus §2 steht in §6/§7 — der
Tag `sdk-kotlin-v0.2.2` (oder die nächste freie Version) existiert auf `origin`
(`git ls-remote --tags origin 'sdk-kotlin-v*'`), `gh run list --workflow
sdk-kotlin-release.yml` nennt für ihn einen Lauf, beide Jobs stehen auf `success`
(`gh run view <Kennung> --json jobs`), und der anonyme Abruf der POM-Datei antwortet
HTTP 200. Der Abruf setzt die URL aus der **README** zusammen — Download-Basis aus dem
`repositories`-Snippet der Kotlin-README, danach das Maven-Standardlayout
`io/github/pt9912/pgchangefeed-kotlin/<Version>/pgchangefeed-kotlin-<Version>.pom`
(hergeleitet, im Slice zu belegen) —, damit die README-Aussage am realen Pfad gehalten
ist. Der Tag-Push ist Betreiber-Handlung; der Implementer nennt sie im Übergabe-Bericht
mit dem genauen Befehl und wartet mit dem Beleg. Zwischen Verifikation und Tag liegt der
Slice in `in-progress/` (WIP-Limit 1 bleibt belegt, weil die Closure den Lauf abwartet:
[`AGENTS.md`](../../../../AGENTS.md) §3.10, „bevor Closure erfolgt“). Der Lerneintrag
(Steering-Loop) gehört zur Closure-Notiz.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht. Die Ausgänge stehen bei Planung als **Erwartung**; gesetzt werden sie bei der
Closure.

- **Erster realer Lauf der Job-Struktur** ([`AGENTS.md`](../../../../AGENTS.md) §3.10): der
  Workflow ändert seine Struktur (zwei Jobs statt einem, ein neues Ziel); jede lokale
  Prüfung ist statisch (YAML, `make gates`), Laufzeit, Registry-Zugangspfad und
  Runner-Unterschiede zeigen sich erst im Tag-Lauf. Zu belegen durch: §5.
  **Ausgang:** offen bis zum Tag-Lauf, dann *entfallen* (grün) oder *eingetreten*
  (rot, Korrektur im Slice oder Folge-Slice mit Kennung).
- **Benutzername der Maven-Anmeldung: Service-Slug oder Service-Name** — eine der drei
  „nicht geprüft“-Aussagen (`ADR-0123` §Kontext; die Service-Seite nennt „Service
  username“ in der Spalte NAME). Gesetzt ist der Slug (Mitteilung des Betreibers).
  *Erwartet:* die Anmeldung gelingt. *Rückfall bei Authentifizierungsfehler
  (HTTP 401/403 im Upload-Schritt):* der Betreiber setzt `CLOUDSMITH_USERNAME` auf den
  Service-Namen `github-ci` (Betreiber-Handlung, kein Agent-Schritt) und wiederholt
  nur den roten Job („Re-run failed jobs“). Die Doku (`releasing.md` Secret-Tabelle)
  nennt danach den Werttyp, der gewirkt hat, nicht den Wert. Der reale Beleg entsteht
  erst mit dem Tag-Push. **Ausgang:** offen bis zum Tag-Lauf.
- **Moduldatei (`.module`) beim nativen Maven-Upload** — Gradle veröffentlicht sie
  standardmäßig mit (`slice-sdk-readme-nutzerdoku` §3: „Publikations-Metadaten verweisen
  darauf“); ob Cloudsmith sie annimmt, ist nicht geprüft (`ADR-0123` §Kontext).
  *Erwartet:* angenommen. *Rückfall (Build-Korrektur ohne neue ADR, `ADR-0123`
  Festlegung 6):* die Moduldatei abschalten; beide Ziele bekommen dieselben Dateien, die
  Korrektur trifft also auch GitHub Packages. **Ausgang:** offen bis zum Tag-Lauf.
- **Anonyme Sichtbarkeit des Repositories** — gemessen nicht bestätigt (§4 Start, Punkt 2).
  *Erwartet:* nach der Betreiber-Antwort sichtbar; der Beleg §5 (anonymer POM-Abruf,
  HTTP 200) entscheidet. **Ausgang:** offen bis zum Tag-Lauf; bleibt der Abruf `404`
  trotz grünem Upload, ist das *eingetreten* und geht als Betreiber-Frage in den
  Übergabe-Bericht (kein Code-Fehler des Slice).
- **Wiederholung und Doppel-Upload derselben Version** — ob ein zweiter Upload am Ziel
  abgelehnt wird, ist für Cloudsmith (nativer Maven-Upload) und GitHub Packages nicht
  geprüft (`ADR-0123` Festlegung 5). Der Ablauf verlässt sich nicht darauf: er
  wiederholt nur den roten Job; ein Teil-Upload wird in der Web-App gelöscht.
  *Erwartet:* im ersten Lauf tritt keine Wiederholung ein. **Ausgang:** *entfallen*,
  wenn der erste Lauf grün ist; *eingetreten* sonst, mit dem beobachteten Verhalten in
  `releasing.md`.
- **Kein Wert eines Secrets in einem Log oder Dokument** — die Zusage des Plans: der
  Workflow gibt keinen Secret-Wert aus (kein `echo`, `printenv`, `set -x` auf Werte;
  Werte über `env:` und `docker run -e NAME`; kein Build-Argument; kein `--info`/
  `--debug` beim Gradle-Aufruf); der Service-Slug steht in keinem Dokument dieses
  Repos. Die Maskierung durch GitHub Actions ist hier nicht gelesen und trägt die
  Zusage nicht. Zu belegen durch: die beiden letzten Zeilen des Suchlaufs (§3) und
  `git grep` nach dem Slug-Anfang am Diff-Stand (erwartet: 0 Zeilen). **Ausgang:**
  *entfallen*, wenn der Suchlauf am Diff-Stand die erwarteten Zahlen druckt.
- **Die README nennt `0.2.2`, bevor der Tag gesetzt ist** — zwischen Merge und Tag
  verweist die Anwender-Doku auf eine Version, die auf keinem Ziel liegt (dasselbe
  Fenster hatte der Vorgänger-Slice bei `0.2.1`). Zu belegen durch: §5 (die Version ist
  abrufbar). **Ausgang:** *entfallen* mit dem Tag-Lauf.
- **Parallele Änderungen an geteilten Trägern** — `docs/user/benutzerhandbuch.md`
  (Version 1.59, Änderungshistorie), `docs/user/releasing.md`, `harness/README.md` und
  `spec/pflichtenheft.md` sind Träger, die andere Slices berühren; der Start liest den
  Ist-Zustand und zieht `Version:` und Historie am dann geltenden Stand nach.
  **Ausgang:** *entfallen*, wenn der Diff keine fremde Zeile überschreibt (`git diff
  <Basis>` je Datei am Ende).
- **Zwei-Quellen-Drift Handbuch gegen Pflichtenheft** — beide Träger nennen denselben
  Sachverhalt (Vertriebsziele des Kotlin-Packages); die Klasse steht mit 3× offen im
  Register (`BEO-PGC/zwei-quellen-drift-handbuch-gegen-pflichtenheft`, gemessen
  `ls …/evidence | wc -l` → 3, 2026-09-25). Zu belegen durch: der Suchlauf-Vergleich
  beider Träger am Diff-Stand. **Ausgang:** *entfallen* bei Deckung, sonst
  *eingetreten* (im Slice behoben, Evidenz-Datei im Register).

*Kein Risiko dieses Slice, sondern benannte Trigger der ADR:* das langlebige Secret
(OIDC statt API-Key: `ADR-0123` Re-Evaluierungs-Trigger 3) und eine Änderung der
Cloudsmith-Bedingungen oder Aussetzung des Repositories nach der nachträglichen Prüfung
(Trigger 2); beide sind dort adressiert, der Slice legt keine zweite Adresse an.

## 7. Closure-Notiz

*(wird bei der Closure durch den Planner gefüllt)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt sind `sdks/kotlin` (Build, README,
Dockerfile), `.github/workflows` (Release-Workflow), `harness/` (Sensor-Zeilen,
Werkzeug-Kommentare), `docs/user/` und `spec/` (Träger). Die Modus-Deklaration führt nur
die Default-Sub-Area `*` (`PGC`, Greenfield) — keine Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register `BEO-PGC` durchgegangen
(Zähler = Zahl der `evidence/`-Dateien, gemessen 2026-09-25). Treffer je berührter
Sub-Area:

- `release-mechanismus-nicht-in-releasing-doku-nachgezogen` (1×, offen): der
  Release-Mechanismus ändert sich; `releasing.md` ist deshalb ein **vorab eingeplanter**
  Liefer-Punkt (§2 Punkt 3), nicht Nachzug nach dem Review.
- `github-actions-unverifizierbar-lokal` (8×, verkörpert in `AGENTS.md` §3.10): der
  Workflow wechselt seine Struktur; §2/§5 tragen den Post-Push-Beleg, §6 das Risiko.
- `zwei-quellen-drift-handbuch-gegen-pflichtenheft` (3×, offen): Handbuch und
  Pflichtenheft beschreiben denselben Sachverhalt; als Risiko in §6.
- `nachzug-laesst-ueberholten-text-stehen` (9×, verkörpert) und
  `arbeit-ueberholt-stehenden-traeger` (31×, verkörpert): die bewegten Eigenschaften
  (Secret-Regime, Vertriebsziel, Version) haben den committeten Suchlauf in §3.
- `zahl-in-traeger-driftet-gegen-die-messung` (21×, verkörpert): die Zahlen im Suchlauf
  stehen wie gedruckt und mit Ursprung; die zwei Summen (`51`, `36`) sind als abgeleitet
  gekennzeichnet.
- `negativtest-ohne-bindung-an-seine-eingabe` (12×, verkörpert): die Mutations-Tabelle
  in §3 führt je Zusage die Eingabeseite, darunter die zweite Versions-Stelle.
- `handbuch-versionshistorie-uebersprungen` (3×, verkörpert): `Version:` und Historie
  des Handbuchs stehen im Liefer-Punkt 2.
- `intern-kennungen-in-ausgelieferten-texten` (1×, offen) und
  `deutsches-fachwort-im-englischen-sdk-readme` (3×, verkörpert): die Kotlin-README ist
  ein ausgelieferter englischer Anwender-Text; `make sdk-public-doc-check` und der
  Sprach-Suchlauf des Vorgänger-Slice gelten für die neue Fassung.
- `kommentar-behauptet-nicht-getragenen-fehlerpfad` (3×, verkörpert): die Kommentare in
  `build.gradle.kts`, Dockerfile und Workflow behaupten nichts über die Anzeige auf den
  Paketseiten (nicht geprüft) und nichts über das Verhalten der Ziele bei Wiederholung.
- `pipe-maskiert-make-exit-code` (4×, verkörpert): jeder Gate-Lauf wird ungepiped
  gelesen und die Folgehandlung (Commit, Push) getrennt beauftragt (DoD, Gates).
- `aufschub-adresse-verfaellt` (3×, verkörpert): die Ereignis-Adresse der Paarungen und
  ihr Fallback stehen in §4.
- Kein Treffer: `drei-sprachen-kopie-divergiert-am-randfall` (2×) — die Änderung betrifft
  nur das Kotlin-Package, keine Aussage wird aus einer Schwester-README kopiert.

**Modus:** alle berührten Sub-Areas GF.

**Umfang: M** (Schätzung, kein Messwert). Ursprung: bis zu elf Dateien im Slice
(abgeleitet aus der Tabelle §3: `pflichtenheft.md`, `build.gradle.kts`, Workflow,
Dockerfile, Kotlin-README, Handbuch, `releasing.md`, `harness/README.md`, `sdk.mk` und
`sdk-pack-kotlin.sh` sind zehn sichere Träger, der Kommentar im Tabellentest-Skript ist
der elfte, nur bei Treffer), davon zwei mit inhaltlichem Umbau (Workflow, 128 Zeilen
laut `wc -l` am Stand `37e825e2`, wird zu zwei Jobs; `releasing.md` §4). Die zwölfte
Datei, `welle-transformationen.md`, ändert diese Planung. Erwartete Diff-Größe: rund 250
bis 350 Zeilen (Schätzung aus dem Umbau des Workflows von einem zu zwei Jobs, dem
Kotlin-Abschnitt von `releasing.md` und vier Handbuch-Stellen — keine gemessene
Vergleichsgröße). Das Verdikt nennt den Zuschnitt „klein“; die Zahl der Träger hebt ihn
auf M, weil drei Liefer-Punkte je mehrere Dateien tragen.
