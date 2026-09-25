# ADR-0123: Kotlin-SDK zusätzlich auf Cloudsmith — anonym beziehbares zweites Publish-Ziel (Supersedes ADR-0109, teilweise)

**Status:** Accepted — Supersedes [`ADR-0109`](0109-kotlin-github-packages-drittes-sdk-package.md)
**teilweise**: genau zwei Aussagen über das Secret-Regime des Publish-Weges
(§Entscheidung Festlegung 5, Punkt „Kein zusätzliches Repository-Secret", und
Zeile 3 der Fitness Function, „ausschließlich `GITHUB_TOKEN`"), soweit sie das
neue zweite Ziel betreffen; alles Übrige von `ADR-0109` bleibt in Kraft,
insbesondere GitHub Packages als Publish-Ziel mit dem eingebauten
`GITHUB_TOKEN`.

**Datum:** 2026-09-25

**Autor:** pt9912 (Architect-Rolle, Modul 8; Nutzerentscheidung im Chat vom
2026-09-25: „Für Kotlin bringt es viel — das hört sich gut an" zu Cloudsmith
nur für Kotlin, nicht für C# und Python; Nutzer-Anweisung „ohne Nachfragen
weiter")

**Bezug:** [`LH-FA-SST-009`](../../../spec/lastenheft.md) (Client-Bibliotheken;
Akzeptanz: „über den Paketmanager seiner Sprache einbinden"),
[`ADR-0109`](0109-kotlin-github-packages-drittes-sdk-package.md) (Haupt-Bezug:
Vertriebsweg GitHub Packages, Re-Evaluierungs-Trigger 5 — Lese-Authentifizierung
als Consumer-Problem), [`ADR-0106`](0106-csharp-nuget-erstes-sdk-package.md) und
[`ADR-0107`](0107-python-pypi-zweites-sdk-package.md) (NuGet.org und PyPI,
anonym lesbar, unberührt), [`ADR-0051`](0051-cicd-pipeline-github-actions.md)
(Registry-Secret-Muster: Betreiber legt das Secret manuell an),
[`AGENTS.md`](../../../AGENTS.md) §3.1 (Docker-only), §3.5 (Accepted-ADR-Immutabilität),
§3.8 (Action-Pinning), §3.10 (realer Post-Push-Lauf), §3.12 (Herkunft von
Aussagen), §3.13 (Träger-Nachzug)

**Schärft:** [`LH-FA-SST-009.a`](../../../spec/pflichtenheft.md) (Vertriebsweg
des Kotlin-Packages), [`SPEC-028`](../../../spec/pflichtenheft.md) (Package-Zeile:
Vertriebsweg `GitHub Packages` wird zu zwei Zielen)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**Das Problem.** [`ADR-0109`](0109-kotlin-github-packages-drittes-sdk-package.md)
wählt GitHub Packages und nimmt bewusst in Kauf, dass ein Consumer zum Lesen ein
GitHub-Konto und einen klassischen Token mit `read:packages` braucht (übernommen
aus `ADR-0109` §Kontext, dort live gegen
`docs.github.com/…/working-with-the-gradle-registry` gelesen; hier nicht neu
gelesen). Der Betreiber beobachtet zusätzlich: die GitHub-Paketseite zeigt zu
einem Maven-Paket weder die README noch die POM-Beschreibung
([`docs/user/releasing.md`](../../user/releasing.md) §4, Absatz „Beschreibung
der Kotlin-Paketseite"). Damit erfüllt der Vertriebsweg die Akzeptanz von
`LH-FA-SST-009` (Bezug über den Paketmanager) nur mit einer Hürde, die NuGet.org
(C#) und PyPI (Python) nicht haben. `ADR-0109` Re-Evaluierungs-Trigger 5 nennt
Maven Central als Rückweg; er ist für diesen Bedarf schwerer als nötig
(Festlegung 4).

**Belege zu Cloudsmith.** Alle Anker sind Seiten von `docs.cloudsmith.com`, heute
per `curl` als HTML geladen und in Text gewandelt (die Seite
`cloudsmith.com/open-source` antwortet mit HTTP 429 und ist nicht gelesen; die
Docs-Seite `resources/open-source-hosting-policy` trägt denselben Sachverhalt).

| Aussage | Beleg-Anker (Stand 2026-09-25) |
|---|---|
| Ein Open-Source-Repository wird ohne vorherige Genehmigung angelegt („We don't require advance permission"); es ist ein „Broadcast" mit Sichtbarkeit „Open source", „publicly available"; eine Prüfung findet nachträglich statt | `docs.cloudsmith.com/resources/open-source-hosting-policy`, Abschnitte „Rules", „Creating an open-source repository" |
| Kontingent: „at least 50GB of artifact data + 200GB of package delivery for free, across all open-source repositories", Nutzung getrennt von anderen Repositories gezählt | dieselbe Seite, Kopfabsatz |
| Regeln: Projekt frei und quelloffen, Antragsteller Maintainer, kein Minor-Fork, Artefakte gehören zum Projekt, Cloudsmith wird „adequately and fairly" genannt (Snippet oder Text mit Link auf `cloudsmith.com` und dem Hinweis auf kostenloses Hosting); für gewinnorientierte Unternehmen mit Umsatz oder Finanzierung ist ein bezahlter Plan „more appropriate", die Basis-OSS-Nutzung bleibt verfügbar | dieselbe Seite, „Rules", „Attribution", „Commercial entities" |
| Anlage im Web-Formular: Repository-Name, Broadcast-Schalter, Sichtbarkeit „Open source", Lizenz aus der SPDX-Liste, URL des Projekts, Zustimmung zu den Hosting-Bedingungen; **Open-Source-Repositories lassen sich nicht über die API anlegen** | `…/resources/open-source-hosting-policy` „Creating an open-source repository"; `…/repositories/create-a-repository` (Satz „It is not possible to create Open Source repositories via the Cloudsmith API") |
| Broadcast: „public page where users can browse, search, and download repository packages without requiring authentication" | `…/repositories/create-a-repository`, Abschnitt „Broadcasts" |
| Maven-Lesezugriff öffentlicher Repositories ohne Zugangsdaten: `https://dl.cloudsmith.io/public/OWNER/REPOSITORY/maven/` als `repositories`-Eintrag, Abhängigkeit `implementation 'GROUP_ID:ARTIFACT_ID:PACKAGE_VERSION'`; private Repositories verlangen Entitlement-Token oder Basic-Auth | `…/formats/gradle-repository` und `…/formats/maven-repository`, Abschnitte „Download / installing", „Public repositories" |
| Upload per Gradle: `maven-publish`, `url = "https://maven.cloudsmith.io/OWNER/REPOSITORY/"`, `credentials { username = 'USERNAME'; password = 'API-KEY' }`; Endpunkt-Alternative `https://api-g.cloudsmith.io/maven` bei sporadischen 443-Fehlern von Gradle mit Java 8 | `…/formats/gradle-repository`, „Upload via gradle publish", „Troubleshooting" |
| Ein „Service" ist ein Konto mit eigenem API-Key ohne Nutzerbindung, gedacht für CI/CD; der Key wird nur einmal angezeigt, lässt sich erneuern; Anlegen verlangt Owner- oder Manager-Rolle | `…/accounts-and-teams/service-accounts` |
| OIDC für GitHub Actions ist dokumentiert: Service, OIDC-Provider-Einstellung (Manager oder Owner), `permissions: id-token: write`, die Action `cloudsmith-io/cloudsmith-cli-action` tauscht das Token; mit `export-auth-token: 'true'` setzt sie `CLOUDSMITH_API_KEY` und `CLOUDSMITH_USERNAME`; ob OIDC im OSS-Angebot verfügbar ist, ist **nicht geprüft** | `…/authentication/setup-cloudsmith-to-authenticate-with-oidc-in-github-actions` |
| Wiederveröffentlichung derselben Version: Meldung „Republishing was not enabled for this package"; Schalter `--republish` oder je Repository unter „Miscellaneous Settings" (Angabe zur CLI; für den nativen Maven-Upload per Gradle **nicht geprüft**) | `…/artifact-management/package-upload`, FAQ |
| Der native Maven-Upload sammelt die Einzeldateien über Group-, Artifact-ID und Version mit einem 60-Sekunden-Sammelfenster; die Verarbeitung ist asynchron, „usually around one minute" bis zur Verfügbarkeit | `…/formats/maven-repository` „Native upload staging window"; `…/artifact-management/package-upload` FAQ |
| Löschen eines Pakets ist eine dokumentierte Aktion (Seite „Delete a package"; Inhalt nicht gelesen) | `docs.cloudsmith.com/llms.txt` |

**Nicht geprüft:** ob Cloudsmith die von Gradle standardmäßig mitveröffentlichte
Moduldatei (`.module`) beim Maven-Upload annimmt; ob der Benutzername eines
Services beim Gradle-Upload frei wählbar ist oder der Service-Name sein muss
(die Service-Seite nennt „Service username" in der Spalte NAME); ob die
Paketseite eines Broadcasts die POM-Beschreibung anzeigt; ob die Lizenz `MIT`
im Formular als SPDX-Eintrag erscheint (`LICENSE` des Repositories ist MIT,
`LICENSE` Zeile 1, gemessen). Die ersten drei klärt der erste reale Lauf
(Festlegung 6), die letzte das Anlegen durch den Betreiber.

**Messung zum Kontingent.** Die Artefakte einer Version (Jar plus Sources-Jar)
liegen entpackt unter 1 MB (`du -sk` der entpackten Verzeichnisse: 548 kB und
260 kB, abgeleitet als obere Schranke für die komprimierten Dateien; Messung am
Stand `0.2.1`); gegen 50 GB Speicher ist das um den Faktor 10⁴ kleiner
(abgeleitet), das Kontingent begrenzt dieses Package nicht.

**Zustand des Repositories.** Der Build führt Version `0.2.1`
(`sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts`); die Tags
`sdk-kotlin-v0.2.0` und `sdk-kotlin-v0.2.1` existieren (`git tag -l
'sdk-kotlin*'`, `git ls-remote --tags origin`). Ein bestehender Tag löst den
Workflow nicht erneut aus; Cloudsmith bekommt deshalb erst die Version des
nächsten Tags.

## Entscheidung

Wir wählen: **Das Kotlin-Package `io.github.pt9912:pgchangefeed-kotlin` wird bei
jedem `sdk-kotlin-v*`-Tag zusätzlich in ein Open-Source-Repository auf
Cloudsmith veröffentlicht; GitHub Packages bleibt Ziel. Nur Kotlin — C# (NuGet.org)
und Python (PyPI) bleiben unverändert.** Sieben Festlegungen:

### 1 — Zwei Ziele, ein Tag, Kotlin allein

- Quelle bleibt der Tag `sdk-kotlin-v<SemVer>`; der Workflow
  `sdk-kotlin-release.yml` veröffentlicht dieselben, in Docker gebauten
  Artefakte in **beide** Ziele.
- Cloudsmith ist der **anonym lesbare** Bezugsweg (Belege: Kontext); GitHub
  Packages bleibt bestehen, weil es ohne Zusatzkonto veröffentlichbar ist, mit
  dem `GITHUB_TOKEN` arbeitet und `ADR-0109` es nicht zurücknimmt. Eine
  Abschaltung von GitHub Packages ist eine eigene Entscheidung, keine Folge
  dieser.
- **Kein Cloudsmith für C# und Python.** NuGet.org und PyPI sind für ihr
  Ökosystem der anonym lesbare Standardweg; ein zweites Ziel dort hätte keinen
  Bedarf, den die Nutzerentscheidung nennt.

### 2 — Repository: Open-Source-Repository, Namen legt der Betreiber fest

- Ein Broadcast-Repository mit Sichtbarkeit „Open source", Lizenz `MIT`,
  Projekt-URL `https://github.com/pt9912/pg-change-feed` (Belege: Kontext).
  Weil sich Open-Source-Repositories nicht über die API anlegen lassen, ist die
  Anlage eine **Betreiber-Handlung** in der Web-App.
- Arbeitsnamen: Organisation (Namespace) `pt9912`, Repository (Slug)
  `pg-change-feed`, also Upload-URL `https://maven.cloudsmith.io/pt9912/pg-change-feed/`
  und Download-URL `https://dl.cloudsmith.io/public/pt9912/pg-change-feed/maven/`.
  Weichen die tatsächlichen Namen ab (der Namespace kann vergeben sein), tragen
  der Slice und die Anwender-Doku die tatsächlichen Namen; das ist keine neue
  Entscheidung.
- Die Einstellung „Republishing" bleibt auf ihrem Standard (aus, Belege:
  Kontext): eine veröffentlichte Version ist unveränderlich, wie es die
  Maven-Konvention verlangt.

### 3 — Authentifizierung: API-Key eines Services, zwei Repository-Secrets, keine Action

- Der Publish-Schritt liest `CLOUDSMITH_USERNAME` (Name des Services) und
  `CLOUDSMITH_API_KEY` (API-Key des Services) als Repository-Secrets; der
  Betreiber legt sie **manuell** an, wie `NUGET_API_KEY` und `PYPI_API_TOKEN`
  ([`docs/user/releasing.md`](../../user/releasing.md) §4, Tabelle „Benötigte
  Repository-Secrets").
- Der Service bekommt Schreibrecht nur auf dieses eine Repository (die
  OIDC-Seite nennt „repository access controls" als Einstellung des Services,
  `…/authentication/setup-cloudsmith-to-authenticate-with-oidc-in-github-actions`,
  Schritt 1).
- **Keine Action, kein OIDC.** OIDC verlangt die Action
  `cloudsmith-io/cloudsmith-cli-action` (Pinning `AGENTS.md` §3.8), eine
  OIDC-Provider-Einstellung mit Manager-/Owner-Rolle und `id-token: write`; die
  Verfügbarkeit im OSS-Angebot ist nicht geprüft. Ein API-Key in einem
  Repository-Secret ist derselbe Mechanismus wie bei NuGet.org und PyPI und
  braucht nichts davon (Alternativen B).
- Der Publish läuft wie bisher **in der Docker-Stufe** `publish`
  (`sdks/kotlin/Dockerfile`); die Werte fließen wie `GITHUB_TOKEN` zur
  **Laufzeit** über `docker run -e`, nie als Build-Argument (Layer-Cache;
  `AGENTS.md` §3.1). Zugriff auf Cloudsmith braucht dort Netz, wie der
  bestehende Lauf für GitHub Packages.

### 4 — Build: zweites `maven`-Repository, Gradle-Aufgaben je Ziel

- Der `publishing`-Block von `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts`
  bekommt ein zweites Repository `Cloudsmith` mit der Upload-URL aus Festlegung 2
  und den Anmeldedaten aus `System.getenv("CLOUDSMITH_USERNAME")` und
  `System.getenv("CLOUDSMITH_API_KEY")`. Das Plugin bleibt `maven-publish`
  (kein Dritt-Plugin, `ADR-0109` Festlegung 5).
- **Folge für den Aufruf:** die Sammel-Aufgabe `publish` führt beide Repositories
  aus; ein Lauf ohne Cloudsmith-Werte scheitert an ihr. Der Workflow ruft
  deshalb die **zwei Einzel-Aufgaben** auf
  (`publishMavenPublicationToGitHubPackagesRepository` und
  `publishMavenPublicationToCloudsmithRepository`, Namen aus Repository-Name und
  Publikation `maven`, abgeleitet aus der Gradle-Namensregel; im Slice erprobt).
  Die `CMD` der Stufe `publish` im `Dockerfile` bleibt (der Workflow überschreibt
  sie ausdrücklich, das ist der bestehende Zustand).
- POM-Felder, `groupId`, `artifactId`, Version und das Sources-Jar bleiben
  unverändert; beide Ziele erhalten dieselben Dateien.

### 5 — Release-Ablauf: unabhängige Ziele, Wiederholung je Ziel

- Jedes Ziel ist ein **eigener Job** desselben Workflows (oder ein Matrix-Zweig
  mit `fail-fast: false`; die Form wählt der Slice): Tag-Validierung,
  Versions-Abgleich, `make sdk-pack-kotlin` und die Publish-Stufe je Ziel. Die
  Berechtigungen bleiben je Job minimal (`packages: write` nur beim GitHub-
  Packages-Job).
- **Fehlerverhalten:** scheitert ein Ziel, bleibt das andere veröffentlicht; der
  Workflow-Lauf ist rot. Der Betreiber wiederholt **nur den roten Job**
  („Re-run failed jobs", `docs.github.com/…/manage-workflow-runs/re-run-workflows-and-jobs`,
  gelesen). Eine erneute Übertragung an das bereits erfolgreiche Ziel entsteht
  dadurch nicht.
- **Idempotenz:** eine Maven-Version ist unveränderlich. Ob ein zweiter Upload
  derselben Version am Ziel abgelehnt wird, ist für Cloudsmith bei nativem
  Maven-Upload und für GitHub Packages **nicht geprüft**; der Ablauf verlässt
  sich nicht darauf, weil er ein erfolgreiches Ziel nicht erneut anstößt. Bleibt
  nach einem Fehlschlag ein Teil-Upload (Einzeldateien ohne vollständiges Paket)
  im Ziel zurück, löscht der Betreiber das Paket in der Web-App und wiederholt
  den Job.
- **Reihenfolge:** keine; die Jobs laufen unabhängig. Cloudsmith verarbeitet
  Uploads asynchron (Kontext); der Job meldet Erfolg mit dem Upload, nicht mit
  der Verfügbarkeit.

### 6 — Anwender-Sicht und Beleg

- **Koordinaten:** `io.github.pt9912:pgchangefeed-kotlin:<Version>`
  (unverändert) plus ein `repositories`-Eintrag
  `maven { url = uri("https://dl.cloudsmith.io/public/pt9912/pg-change-feed/maven/") }`
  ohne Zugangsdaten (Belege: Kontext). Die Abhängigkeiten des Packages
  (gRPC, Gson, jnats, Coroutines) löst der Consumer wie bisher über Maven
  Central auf; das Repository bietet sie nicht an.
- **Doku:** die README `sdks/kotlin/pgchangefeed-kotlin/README.md` nennt
  Cloudsmith als **ersten** Bezugsweg ohne Token und GitHub Packages als zweiten
  mit Token-Hinweis; sie trägt die von der Open-Source-Regel verlangte
  Namensnennung (Belege: Kontext, „Attribution"). Die Anpassung ist Teil des
  Folge-Slices.
- **Beleg:** netzlos prüfbar ist die Publish-Konfiguration (Aufgabe und POM im
  Docker-Bau vorhanden, ohne Zugangsdaten). Der **reale** Beleg entsteht erst
  mit dem ersten echten Tag-Push ([`AGENTS.md`](../../../AGENTS.md) §3.10): der
  Lauf ist grün, und ein **anonymer** Abruf der veröffentlichten POM-Datei am
  Download-Pfad des Repositories antwortet mit HTTP 200 (erwartet, nach dem
  Sammelfenster und der asynchronen Verarbeitung ggf. mit Wiederholung; die
  URL-Form folgt dem Maven-Standardlayout unter dem `…/maven/`-Pfad,
  hergeleitet und im Slice zu belegen). Bleibt der Beleg aus, weil eine der drei
  „nicht geprüft"-Aussagen (Moduldatei, Service-Benutzername, Paketseite) nicht
  trägt, ist das eine Build- oder Workflow-Korrektur ohne neue ADR (z. B. das
  Abschalten der Moduldatei); eine Änderung des Publish-Mechanismus selbst wäre
  eine Folge-ADR.

### 7 — Was diese ADR nicht ändert

- Kein Produktionscode, kein Eingriff in `internal/**` und `cmd/**`.
- [`ADR-0109`](0109-kotlin-github-packages-drittes-sdk-package.md) Festlegungen 1
  (Umfang), 3 (Ort), 4 (SemVer), 6 und das Publish-Ziel GitHub Packages samt
  `GITHUB_TOKEN` bleiben.
- Der Tag-Namensraum `sdk-kotlin-v*`, die Tag-Validierung und der Abgleich mit
  der Versionszeile bleiben.
- Kein neues Gate; ein Publish braucht Netz und ist Werkzeug.
- `sdks/csharp/**`, `sdks/python/**` und ihre Workflows.
- Keine Signierung, kein Javadoc-Jar (Maven Central verlangt beides,
  Verglichene Alternativen C).

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

### A — Vertriebsweg für den anonymen Bezug

| Option | Pro | Contra |
|---|---|---|
| A1 — nur GitHub Packages (nichts tun) | kein Konto, kein Secret | Token-Pflicht zum Lesen bleibt (`ADR-0109`); die Paketseite zeigt README und POM-Beschreibung nicht; die Akzeptanz „über den Paketmanager einbinden" ist mit Hürde erfüllt, die die Schwester-SDKs nicht haben |
| **A2 — GitHub Packages und Cloudsmith-OSS-Repository (gewählt)** | anonymer Bezug ohne GitHub-Konto (Belege: Kontext); Anlage ohne Genehmigung, Kontingent 50 GB/200 GB gegen unter 1 MB je Version; derselbe Gradle-Mechanismus (`maven-publish`), nur ein zweites Repository; GitHub Packages bleibt | ein zusätzliches externes Konto und zwei Secrets; Namensnennung Cloudsmiths im Projekt verpflichtend; ein `repositories`-Eintrag beim Consumer bleibt nötig; Abhängigkeit von einem Drittanbieter mit eigenen Nutzungsbedingungen (Prüfung nachträglich, Aussetzung möglich, Belege: Kontext) |
| A3 — Maven Central statt Cloudsmith | Standardweg der JVM-Welt; kein zusätzlicher `repositories`-Eintrag beim Consumer | Central-Portal-Konto, Namespace-Verifikation (für `io.github.<Nutzer>` bei GitHub-Anmeldung „in most cases" automatisch, sonst über den Portal-Dialog), Sources- **und** Javadoc-Jar, GPG-Signatur und Prüfsummen für jede Datei (`central.sonatype.org/publish/requirements/`, `…/register/namespace/`, heute gelesen); der Build führt heute nur ein Sources-Jar; deutlich größerer Aufwand für denselben Nutzen |
| A4 — JitPack aus dem Git-Tag | kein Konto | Bau beim Consumer aus dem Quelltext, keine geprüfte Lieferung; in `ADR-0109` (B5) bereits abgelehnt |

### B — Authentifizierung gegen Cloudsmith

| Option | Pro | Contra |
|---|---|---|
| **B1 — API-Key eines Services als Repository-Secret (gewählt)** | keine Action (kein Pinning-Aufwand), derselbe Mechanismus wie NuGet.org/PyPI, plan-unabhängig; Publish bleibt in der Docker-Stufe | ein langlebiges Secret, vom Betreiber erneuerbar (Service-Seite) |
| B2 — OIDC aus GitHub Actions | kein langlebiges Secret | verlangt `cloudsmith-io/cloudsmith-cli-action` (Pinning §3.8), OIDC-Provider-Einstellung mit Manager-/Owner-Rolle und `id-token: write`; Verfügbarkeit im OSS-Angebot nicht geprüft; das Token müsste in die Docker-Stufe weitergereicht werden |
| B3 — persönlicher API-Key des Betreibers | ein Konto weniger | Bindung an eine Person, Rechte über das Repository hinaus; die Service-Seite empfiehlt Services für CI/CD |

### C — Umfang der Änderung

| Option | Pro | Contra |
|---|---|---|
| **C1 — Cloudsmith jetzt, Maven Central später bei Bedarf (gewählt)** | kleinster Schritt, der den anonymen Bezug herstellt; kein GPG-Schlüssel, kein Javadoc-Jar, kein Portal-Konto | zusätzlicher `repositories`-Eintrag beim Consumer bleibt bis Maven Central |
| C2 — Maven Central jetzt, Cloudsmith nie | Standardweg | Aufwand siehe A3; Nutzerentscheidung nennt Cloudsmith |
| C3 — alle drei SDKs auf Cloudsmith | ein Betreiber-Konto für alle | NuGet.org und PyPI sind anonym lesbar; das zweite Ziel schafft dort keinen Nutzen und drei Publish-Wege mehr; Nutzerentscheidung: nur Kotlin |

### D — Ablauf bei zwei Zielen

| Option | Pro | Contra |
|---|---|---|
| **D1 — ein Job je Ziel, Wiederholung je Job (gewählt)** | ein Fehlschlag blockiert das andere Ziel nicht; „Re-run failed jobs" trifft nur das rote Ziel | Bau und Tests laufen je Job (zweimal `make sdk-pack-kotlin`) |
| D2 — ein Job, ein `./gradlew publish` | einfachster Workflow | scheitert das eine Ziel, ist das andere abhängig von der Aufgaben-Reihenfolge; ein Neulauf sendet erneut an das erfolgreiche Ziel |
| D3 — ein Job, zwei Schritte mit `continue-on-error` | ein Bau | maskiert den Fehler im Schritt-Ergebnis (`docs/user/releasing.md` §4 beschreibt dieselbe Klasse beim `continue-on-error` des Schwester-Repos d-check); ein Neulauf wiederholt beide Ziele |

**Fazit:** A2, B1, C1, D1. A1 lässt die Hürde stehen, A3/C2 verlangen für
denselben Nutzen mehr (Signatur, Javadoc, Portal), A4 senkt die
Lieferzusicherung, B2 braucht eine Action und eine nicht geprüfte
Plan-Voraussetzung, B3 bindet an eine Person, C3 hat für NuGet.org/PyPI keinen
Bedarf, D2/D3 koppeln die Ziele.

## Konsequenzen

- **Positiv:** Kotlin-Anwender beziehen die Bibliothek ohne GitHub-Konto und
  Token; der Publish-Weg bleibt Docker-only und action-frei; ein Ziel ist ohne
  das andere wiederholbar; die Kotlin-README kann am Broadcast-Repository den
  Bezug ohne Vorbedingung beschreiben.
- **Negativ:** ein weiterer Drittanbieter (Konto, zwei Secrets, Pflicht zur
  Namensnennung, Eignungsprüfung nach der Anlage); Consumer tragen weiterhin
  einen `repositories`-Eintrag; der Release-Workflow baut je Ziel; die Versionen
  `0.2.0` und `0.2.1` bleiben auf GitHub Packages allein (Tags bestehen).
- **Folgepflicht:**
  1. **Slice `slice-sdk-kotlin-cloudsmith`** (ein kleiner Slice; Zuschnitt in
     der Verdikt-Datei): `build.gradle.kts` (zweites Repository), Workflow
     (ein Job je Ziel, explizite Gradle-Aufgaben), `sdks/kotlin/Dockerfile`
     (Kommentare der Stufe `publish`), Anwender-README samt Namensnennung.
     **Startbedingung** (beobachtbar): Betreiber hat Konto, Organisation und
     Repository angelegt, dem Planner die tatsächlichen Namen genannt, und
     `gh secret list` zeigt `CLOUDSMITH_USERNAME` und `CLOUDSMITH_API_KEY`.
  2. **Träger-Nachzug** ([`AGENTS.md`](../../../AGENTS.md) §3.13): `spec/pflichtenheft.md`
     (`LH-FA-SST-009.a`, Absatz Kotlin; `SPEC-028`-Zeile, Vertriebsweg),
     `harness/README.md` (Zeile `sdk-kotlin-release.yml`),
     [`docs/user/releasing.md`](../../user/releasing.md) §4 (Secret-Tabelle,
     Abschnitt „GitHub-Packages-Publish", Absatz zur Paketseite),
     `docs/user/benutzerhandbuch.md` (SDK-Hinweis zur Token-Pflicht),
     Kopfkommentare von Workflow und `build.gradle.kts`, die „kein externes
     Secret" sagen.
  3. **Betreiber-Voraussetzungen** (nicht durch einen Agenten ausführbar; Liste
     in der Verdikt-Datei): Cloudsmith-Konto und Organisation, Open-Source-
     Repository per Web-App, Service mit Schreibrecht, zwei Repository-Secrets.
  4. Der Workflow gilt nach [`AGENTS.md`](../../../AGENTS.md) §3.10 erst nach
     einem realen, grünen Post-Push-Lauf als abgeschlossen; der anonyme Abruf
     aus Festlegung 6 gehört zu diesem Beleg.

## Fitness Function (falls maschinell prüfbar)

Die Zeilen sind **hergeleitet**, nicht erprobt (die Umsetzung existiert nicht);
der Slice erprobt sie.

| Tooling | Regel | Make-Target |
|---|---|---|
| Docker-Bau (Folgepflicht 1) | im Bau ohne Zugangsdaten existieren die Gradle-Aufgaben beider Ziele und die generierte POM trägt `groupId`/`artifactId`/Version; fehlt ein Repository-Block, färbt der Bau rot | `make sdk-pack-kotlin` (Erweiterung) |
| Review (kein Sensor) | der Cloudsmith-Job erhält `CLOUDSMITH_USERNAME`/`CLOUDSMITH_API_KEY` nur zur Laufzeit über `docker run -e`, nie als Build-Argument | — |
| Review (kein Sensor) | jeder Job trägt die minimalen Berechtigungen (`packages: write` nur beim GitHub-Packages-Job); keine neue `uses:`-Zeile ohne SHA (`AGENTS.md` §3.8) | — |
| Ersetzt Zeile 3 der Fitness Function von `ADR-0109` | GitHub-Packages-Job: ausschließlich `GITHUB_TOKEN`; Cloudsmith-Job: ausschließlich die zwei Cloudsmith-Secrets; kein weiteres Secret | — |

## Re-Evaluierungs-Trigger

1. **Maven Central wird für Kotlin gebraucht** (Nachfrage nach einem Bezug ohne
   zusätzlichen `repositories`-Eintrag, oder `1.0.0` mit
   Rückwärtskompatibilitätszusage): eigene ADR; Aufwand: Portal-Konto, Sources-
   und Javadoc-Jar, Signatur, Prüfsummen (Alternativen A3).
2. **Cloudsmith ändert die Bedingungen des Open-Source-Angebots** (Kontingent,
   Regeln, Sichtbarkeit) oder setzt das Repository aus: Bewertung, ob GitHub
   Packages plus Maven Central genügt.
3. **OIDC im OSS-Angebot verfügbar und der Betreiber will kein langlebiges
   Secret:** Folge-ADR zu Festlegung 3.
4. **GitHub Packages wird für Kotlin abgeschaltet** (kein Bedarf mehr am Weg mit
   Token): Folge-ADR, die `ADR-0109` Festlegung 2 ablöst.
5. **Der erste reale Lauf zeigt, dass ein Mechanismus (Upload, Anmeldung,
   Moduldatei) nicht trägt:** Korrektur im Slice, bei geändertem Mechanismus
   Folge-ADR mit `Supersedes ADR-0123`.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-25 | Accepted — Nutzerentscheidung im Chat (Cloudsmith nur für Kotlin), Belege aus `docs.cloudsmith.com`, `central.sonatype.org` und `docs.github.com`, heute gelesen | [`LH-FA-SST-009`](../../../spec/lastenheft.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0123` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
