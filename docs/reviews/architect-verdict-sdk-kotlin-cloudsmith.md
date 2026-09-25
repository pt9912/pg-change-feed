# Architect-Verdikt: Kotlin-SDK zusätzlich auf Cloudsmith

**Rolle:** Architect (Modul 8)

**Anlass:** Der Auftraggeber will das Kotlin-SDK-Package `pgchangefeed-kotlin`
zusätzlich auf Cloudsmith (Open-Source-Angebot) veröffentlichen — „Für Kotlin
bringt es viel — das hört sich gut an", nur für Kotlin, nicht für C# und Python;
PyPI und NuGet.org bleiben der Standard. Grund: Kotlin-Anwender sollen die
Bibliothek ohne GitHub-Token beziehen können. Der Zug prüft, ob Cloudsmith das
trägt, und entscheidet den Weg. Der Rollenwechsel Planner → Architect → Planner
folgt Modul 8 §Konflikt-Pfad.

**Rolleninhaber:** pt9912 (Architect-Zug, frischer Kontext)

**Datum:** 2026-09-25

**Bezug:** [`LH-FA-SST-009`](../../spec/lastenheft.md),
[`LH-FA-SST-009.a`](../../spec/pflichtenheft.md),
[`SPEC-028`](../../spec/pflichtenheft.md),
[`ADR-0109`](../plan/adr/0109-kotlin-github-packages-drittes-sdk-package.md),
[`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md); der Slice
ist als Kennung genannt, nicht als Pfad-Link (ein Slice wechselt die
Lifecycle-Ablage, ein Pfad-Link bräche mit)

**Erzeugte Artefakte dieses Zugs:**

- [`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md)
  (Accepted, `Supersedes ADR-0109` teilweise — nur die Aussagen zum
  Secret-Regime des Publish-Weges); Index-Zeile in `docs/plan/adr/README.md`
- dieses Dokument

Pläne, Spec, Handbuch, Workflow und Code bleiben in diesem Zug unberührt.

**Werkzeug-Hinweis:** Der Auftrag nennt `WebFetch`; dieser Lauf hat das Werkzeug
nicht. Die Docs-Seiten sind mit `curl` als HTML geladen und in Text gewandelt
(die Docs-Seiten von `docs.cloudsmith.com` laufen ohne die 429-Sperre, die
`cloudsmith.com/open-source` liefert). Jede Aussage über Cloudsmith trägt ihren
Anker in [`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md)
§Kontext; was nicht gelesen ist, steht dort als „nicht geprüft".

---

## Verdikt

**Cloudsmith trägt den Bedarf. Empfehlung: jetzt als zweites Publish-Ziel für
Kotlin, Maven Central später und nur bei Bedarf.** Konflikt-Pfad-Verdikt: die
Entscheidung wird per Folge-Entscheidung geschärft (`ADR-0123`, teilweise
ablösend — eine Aussage von `ADR-0109` über die Secrets wird für das zweite Ziel
falsch, deshalb `Supersedes` teilweise und nicht bloß `Schärft`).

**1 — Anonymer Lesezugriff auf Maven-Repositories: ja, im Open-Source-Angebot.**
Ein Open-Source-Repository ist ein Broadcast mit Sichtbarkeit „Open source",
öffentlich ohne Authentifizierung lesbar; der Maven-Bezug läuft über
`https://dl.cloudsmith.io/public/<Organisation>/<Repository>/maven/`. Konditionen
(alle mit Anker in `ADR-0123`):

| Punkt | Stand |
|---|---|
| Kontingent | mindestens 50 GB Artefakte und 200 GB Auslieferung über alle Open-Source-Repositories; dieses Package liegt bei unter 1 MB je Version (abgeleitet aus `du -sk`, entpackt 548 + 260 kB) — das Kontingent begrenzt nichts |
| Freigabe | keine vorherige Genehmigung; Anlage im Web-Formular, Prüfung nachträglich |
| Lizenz | SPDX-Liste; `LICENSE` des Repositories ist MIT (gemessen); ob `MIT` im Formular erscheint: nicht geprüft |
| Bedingungen | freies, quelloffenes Projekt, Antragsteller Maintainer, Namensnennung Cloudsmiths (Snippet oder Text mit Link auf `cloudsmith.com` und Hinweis auf kostenloses Hosting); für gewinnorientierte Unternehmen mit Umsatz oder Finanzierung ist ein bezahlter Plan angemessener |
| Anlage | Organisation (Namespace) und Repository (Slug) legt der Betreiber an; **Open-Source-Repositories lassen sich nicht über die API anlegen** — nur Web-App |
| Nicht geprüft | ob `cloudsmith.com/open-source` weitere Konditionen nennt (HTTP 429); der Inhalt der Seite „Open-source hosting policy" in `docs.cloudsmith.com` ist gelesen |

**2 — Veröffentlichungsweg aus der CI: Gradle `maven-publish`, API-Key eines
Services, zwei Repository-Secrets, keine Action.**

- Gradle veröffentlicht per `maven-publish` an `https://maven.cloudsmith.io/<Organisation>/<Repository>/`
  mit Benutzername und API-Key (Beleg: Cloudsmith-Seite „Gradle repository").
- Auth: der API-Key eines **Services** (Konto ohne Nutzerbindung für CI/CD).
  Secrets `CLOUDSMITH_USERNAME` und `CLOUDSMITH_API_KEY`, vom Betreiber manuell
  angelegt wie `NUGET_API_KEY` und `PYPI_API_TOKEN`
  ([`docs/user/releasing.md`](../user/releasing.md) §4).
- OIDC aus GitHub Actions ist dokumentiert, verlangt aber die Action
  `cloudsmith-io/cloudsmith-cli-action` (Pinning `AGENTS.md` §3.8) und eine
  Provider-Einstellung mit Manager-/Owner-Rolle; ob es im OSS-Angebot verfügbar
  ist, ist nicht geprüft. **Der API-Key trägt ohne diese Voraussetzungen**; OIDC
  ist Re-Evaluierungs-Trigger 3 der ADR.
- Docker-only: der Publish läuft unverändert in der Docker-Stufe `publish`
  (`sdks/kotlin/Dockerfile`), mit einem zweiten Repository-Ziel; die Werte
  fließen zur Laufzeit über `docker run -e`, nie als Build-Argument.

**3 — Koordinaten und Anwender-Sicht.** Die Koordinate bleibt
`io.github.pt9912:pgchangefeed-kotlin:<Version>`; neu ist ein
`repositories`-Eintrag `maven { url = uri("https://dl.cloudsmith.io/public/<Organisation>/<Repository>/maven/") }`
ohne Zugangsdaten. Arbeitsnamen: Organisation `pt9912`, Repository `pg-change-feed`
(der Betreiber kann abweichen, der Slice trägt die tatsächlichen Namen). Die
Abhängigkeiten des Packages löst der Consumer über Maven Central auf.
`sdks/kotlin/pgchangefeed-kotlin/README.md` nennt heute nur GitHub Packages samt
Token-Pflicht (Abschnitt „Installation"); im Folge-Slice steht Cloudsmith als
erster Bezugsweg, GitHub Packages als zweiter, dazu die Namensnennung Cloudsmiths.

**4 — Konsequenzen.**

- **`ADR-0109` (Accepted, unberührbar):** bleibt vollständig lesbar; abgelöst
  werden nur zwei Aussagen (Festlegung 5, „kein zusätzliches Repository-Secret",
  und Zeile 3 der Fitness Function „ausschließlich `GITHUB_TOKEN`") — für den
  Cloudsmith-Job gilt die neue Aussage. Das Publish-Ziel GitHub Packages bleibt.
- **Aufruf:** die Sammel-Aufgabe `publish` führt beide Repositories aus und
  scheitert ohne Cloudsmith-Werte; der Workflow ruft je Ziel die
  Einzel-Aufgabe auf (Namen aus Repository-Name und Publikation abgeleitet, im
  Slice zu erproben).
- **Release-Ablauf, ein Tag, zwei Ziele:** ein Job je Ziel (oder
  Matrix-Zweig, `fail-fast: false`). Scheitert ein Ziel, bleibt das andere
  veröffentlicht, der Lauf ist rot, der Betreiber wiederholt nur den roten Job
  („Re-run failed jobs", gelesen). Damit entsteht keine erneute Übertragung an
  das erfolgreiche Ziel; ob ein zweiter Upload derselben Version am Ziel
  abgelehnt wird, ist für den nativen Maven-Upload bei Cloudsmith und bei GitHub
  Packages **nicht geprüft** — der Ablauf verlässt sich nicht darauf. Die
  Cloudsmith-Einstellung „Republishing" bleibt aus (Standard). Teil-Upload nach
  Fehlschlag: Paket in der Web-App löschen, Job wiederholen.
- **Test- und Beleg-Form:** netzlos prüfbar ist die Konfiguration (Aufgabe und
  POM im Docker-Bau, ohne Zugangsdaten); der reale Beleg entsteht erst mit dem
  ersten Tag-Push (`AGENTS.md` §3.10): grüner Lauf **und** ein anonymer Abruf der
  POM-Datei am Download-Pfad mit HTTP 200 (erwartet; die Verarbeitung ist
  asynchron, „usually around one minute").
- **Versionsstand:** die Tags `sdk-kotlin-v0.2.0` und `sdk-kotlin-v0.2.1`
  existieren; Cloudsmith bekommt die Version des nächsten Tags, die älteren
  bleiben allein auf GitHub Packages.

**5 — Maven Central: später, nur bei Bedarf.** Aufwand gegenüber Cloudsmith
(heute gelesen, `central.sonatype.org/publish/requirements/` und
`…/register/namespace/`): Central-Portal-Konto, Namespace-Verifikation (für
`io.github.<Nutzer>` bei GitHub-Anmeldung in den meisten Fällen automatisch),
Sources- **und** Javadoc-Jar (der Build führt nur ein Sources-Jar), GPG-Signatur
und Prüfsummen je Datei. Nutzen gegenüber Cloudsmith: kein zusätzlicher
`repositories`-Eintrag beim Consumer. Der Nutzen rechtfertigt den Aufwand erst
bei nachgewiesener Nachfrage oder mit `1.0.0` — Re-Evaluierungs-Trigger 1 der
ADR. Die Reihenfolge „Cloudsmith jetzt, Maven Central später" ist die
Empfehlung; „Maven Central zuerst" würde Signatur, Javadoc und Portal vor den
ersten anonymen Bezug stellen.

**6 — Zuschnitt: ein kleiner Slice `slice-sdk-kotlin-cloudsmith`,** siehe
§Zuschnitt. Die Betreiber-Voraussetzungen stehen in §Betreiber-Voraussetzungen.

---

## Zuschnitt

**Slice `slice-sdk-kotlin-cloudsmith`** (Empfehlung: ein Slice, Implementer-Umfang
klein, keine Wellen-Bindung nötig; der Planner ordnet ihn in die Roadmap ein).

**Startbedingung** (beobachtbar, kein Agenten-Schritt): der Betreiber hat Konto,
Organisation und Repository angelegt und die tatsächlichen Namen genannt, und
`gh secret list` nennt `CLOUDSMITH_USERNAME` und `CLOUDSMITH_API_KEY`.

**Berührte Dateien (Vorschlag):**

| Datei | Änderung |
|---|---|
| `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` | zweites `maven`-Repository `Cloudsmith`; Kopfkommentar („kein Repository-Secret") auf beide Ziele nachziehen |
| `.github/workflows/sdk-kotlin-release.yml` | je Ziel ein Job oder Matrix-Zweig; explizite Gradle-Aufgaben; Berechtigungen je Job minimal; Kopfkommentar |
| `sdks/kotlin/Dockerfile` | Kommentare der Stufe `publish` (nur GitHub Packages genannt) |
| `sdks/kotlin/pgchangefeed-kotlin/README.md` | Abschnitt „Installation": Cloudsmith zuerst, GitHub Packages zweitens, Namensnennung |
| Test-/Beleg-Form | netzlose Probe im Docker-Bau, dass die zwei Aufgaben und die POM existieren (`make sdk-pack-kotlin` erweitert) |

**Definition of Done (Skizze, Zusagen statt Feststellungen):** die Probe im Bau
färbt rot, wenn eines der Repositories fehlt (Mutation belegt); `make gates` grün;
die Träger aus §Folge-Arbeit gezogen; Risiko „erster realer Lauf" (`AGENTS.md`
§3.10) bleibt offen, bis der Post-Push-Lauf grün ist **und** der anonyme
POM-Abruf HTTP 200 liefert. Drei „nicht geprüfte" Aussagen klärt dieser Lauf
(Moduldatei, Benutzername des Services, Anzeige der POM-Beschreibung auf der
Paketseite); jede führt zu einer Build-/Workflow-Korrektur im Slice, nicht zu
einer neuen ADR.

**§3.13-Suchlauf (Träger der bewegten Eigenschaft „Vertriebsweg/Secret-Regime des
Kotlin-Packages"):** `grep` nach `GitHub Packages`, `GITHUB_TOKEN`, `kein externes
Secret`, `Personal-Access-Token` über `spec/pflichtenheft.md`, `harness/README.md`,
`docs/user/releasing.md`, `docs/user/benutzerhandbuch.md`,
`sdks/kotlin/**`, `.github/workflows/sdk-kotlin-release.yml`,
`docs/plan/planning/in-progress/roadmap.md`; beide Stände messen, Ergebnis im
Slice-Plan festhalten.

---

## Betreiber-Voraussetzungen

Nicht durch einen Agenten ausführbar (externe Konten, Secrets):

1. **Cloudsmith-Konto** anlegen und eine **Organisation** (Namespace) erstellen;
   Arbeitsname `pt9912` (bei Vergabe: anderer Name, dem Planner nennen).
2. **Open-Source-Repository per Web-App** anlegen: Name `pg-change-feed`,
   Broadcast-Schalter an, Sichtbarkeit „Open source", Lizenz `MIT`, Projekt-URL
   `https://github.com/pt9912/pg-change-feed`, Zustimmung zu den Hosting-Bedingungen.
   Über die API ist das nicht möglich. „Republishing" auf dem Standard lassen.
3. **Service** (Owner-/Manager-Rolle nötig) anlegen, ihm Schreibrecht auf dieses
   Repository geben und den **API-Key sofort kopieren** (er wird nur einmal
   angezeigt; Erneuern ist möglich).
4. **Zwei Repository-Secrets** in GitHub setzen: `CLOUDSMITH_USERNAME` (Name des
   Services), `CLOUDSMITH_API_KEY` (Key des Services).
5. **Namen mitteilen** (Organisation, Repository-Slug), damit Slice und
   Anwender-Doku sie tragen.
6. **Nach dem Slice:** einen neuen Tag `sdk-kotlin-v<nächste Version>` pushen (die
   Tags `sdk-kotlin-v0.2.0`/`sdk-kotlin-v0.2.1` lösen den Workflow nicht erneut aus); die
   Version wird zuvor in `build.gradle.kts` und den Trägern gehoben.
7. **Optional:** Beschreibungstext des Repositories/der Paketseite in der
   Web-App pflegen (ob die POM-Beschreibung angezeigt wird, ist nicht geprüft).

---

## Folge-Arbeit

Jede Pflicht hat einen Träger; die Pläne ändert dieses Dokument nicht — der
Planner zieht sie nach (`AGENTS.md` §3.13).

| Träger | Nachzug | Herkunft |
|---|---|---|
| Planner | Slice `slice-sdk-kotlin-cloudsmith` anlegen (Zuschnitt oben), Startbedingung eintragen, in die Roadmap ordnen | `ADR-0123` Folgepflicht 1 |
| `slice-sdk-kotlin-cloudsmith` | `build.gradle.kts`, Workflow, `sdks/kotlin/Dockerfile`-Kommentare, Kotlin-README, netzlose Probe im Bau | `ADR-0123` Festlegung 4/5/6 |
| `spec/pflichtenheft.md` (Planner-Spec-Zug, im Slice) | `LH-FA-SST-009.a` Absatz Kotlin und `SPEC-028`-Zeile: zwei Vertriebsziele | `ADR-0123` Folgepflicht 2 |
| `harness/README.md` | Zeile `sdk-kotlin-release.yml`: zwei Ziele, zwei Secrets; erst wenn der Workflow real existiert (`AGENTS.md` §4) | `ADR-0123` Folgepflicht 2 |
| `docs/user/releasing.md` §4 | Secret-Tabelle um zwei Zeilen; Abschnitt „GitHub-Packages-Publish" um den zweiten Weg; Absatz zur Paketseite | `ADR-0123` Folgepflicht 2 |
| `docs/user/benutzerhandbuch.md` | SDK-Hinweis zur Token-Pflicht: anonymer Bezug über Cloudsmith | `ADR-0123` Folgepflicht 2, `ADR-0109` Folgepflicht 4 |

**Reihenfolge:** Betreiber-Voraussetzungen 1–5 → Slice → Tag-Push (Voraussetzung
6) → Post-Push-Lauf und anonymer POM-Abruf (Beleg, `AGENTS.md` §3.10).

---

## Auftraggeber-Frage

**Keine offen.** `ADR-0123` steht `Accepted` (Vollmacht zum Architect-Zug, jede
Cloudsmith-Aussage an ihrer Docs-Seite belegt oder als „nicht geprüft"
gekennzeichnet, `AGENTS.md` §3.12). Die einzigen Betreiber-Handlungen stehen in
§Betreiber-Voraussetzungen.
