# Slice sdk-csharp-publish-workflow: NuGet.org-Publish-Workflow (`sdk-csharp-release.yml`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-csharp-lh-fa-sst-009](../welle-sdk-csharp-lh-fa-sst-009.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md),
[`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md)
Festlegung 4 (Trigger, Secret, Tag-Präfix), Alternative E2 (gewählt),
[`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
Entscheidung 3/8 (Release-Tag-Trigger, Docker-Hub-Secret-Muster — Vorbild).

**Berührte Spec-Stellen:** — (Prozess-/CI-Artefakt ohne eigene
`SPEC-*`-Kennung, analog `release.yml`/`hub-description.yml`).

**Verantwortlich:** Implementer-Agent (priorisiert 2026-09-19, direkter
Auftrag: NuGet.org-Publish-Workflow für das C#-SDK-Package umsetzen).

**Autor:** Planner-Agent, direkt beauftragt (`ADR-0106` §Konsequenzen
Folgepflicht 1/5). **Datum:** 2026-09-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein neuer, Netz-bindender GitHub-Actions-Workflow
`.github/workflows/sdk-csharp-release.yml`, Trigger ausschließlich
`push: tags: ['sdk-csharp-v*']`, der den Tag strikt gegen SemVer 2.0
validiert, ihn gegen die im `.csproj` geführte `<Version>` abgleicht
(Abbruch bei Abweichung, vor jedem Build/Push — Muster `release.yml`s
Tag-gegen-`version.md`-Abgleich, hier gegen die Projektdatei statt gegen
eine Markdown-Datei), das `.nupkg` aus `make sdk-pack-csharp` erzeugt und
per `dotnet nuget push … --api-key $NUGET_API_KEY --source
https://api.nuget.org/v3/index.json` veröffentlicht.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Anlage des `NUGET_API_KEY`-Repository-Secrets** — der Workflow
  **referenziert** das Secret, legt es nicht an (`ADR-0106` Festlegung 4:
  „ein Repository-Secret, das der Workflow referenziert, nicht anlegt" —
  dieselbe Klasse wie `DOCKERHUB_USERNAME`/`DOCKERHUB_TOKEN`,
  [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
  Entscheidung 8). Die reale Anlage ist eine externe, außerhalb des Repos
  liegende Handlung des Repository-Betreibers.
- **Denselben Tag-Namensraum wie `release.yml` (`v*`)** —
  `ADR-0106` Festlegung 4/Alternative E3 (verworfen) begründet ausdrücklich
  den eigenen Präfix `sdk-csharp-v*`, um die unabhängige SDK-Versionierung
  (Festlegung 3) nicht mit dem Server-Release-Raum kollidieren zu lassen.
- **`:latest`-Äquivalent für das SDK** — NuGet kennt kein Analogon zu
  einem Docker-`:latest`-Tag; ein veröffentlichtes Package ist über seine
  exakte Version referenzierbar, kein zusätzlicher Mechanismus nötig.
- **GitHub-Release-Eintrag für das SDK** (analog `release.yml`s
  `gh release create`) — nicht Teil von `ADR-0106` Festlegung 4; kann ein
  späterer, eigenständiger Nachzug sein, wenn Bedarf entsteht, aber kein
  Bestandteil dieser Folgepflicht.
- **Der reale, grüne Post-Push-Lauf gegen NuGet.org** — bleibt nach
  [`AGENTS.md`](../../../../AGENTS.md) §3.10 ein offenes Risiko dieses
  Slice bis zum ersten tatsächlichen Tag-Push (siehe §6); `make gates`
  grün und ein plausibler YAML-Aufbau sind keine Ersatz-Bestätigung.

## 2. Definition of Done

- [x] `.github/workflows/sdk-csharp-release.yml` existiert: Trigger
      ausschließlich `push: tags: ['sdk-csharp-v*']` (`ci.yml`/`e2e.yml`/
      `release.yml` schließen diesen Tag-Namensraum durch ihre jeweils
      eigenen `tags`/`tags-ignore`-Filter aus, kein Doppellauf — Beleg im
      Bericht dieses Slice); validiert den Tag-Suffix (nach `sdk-csharp-v`)
      strikt gegen SemVer 2.0; liest `<Version>` aus
      `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj` (per
      `grep -oE`/`sed` direkt auf dem Runner statt im Toolchain-Container —
      Plan-Nachzug §3, `AGENTS.md` §3.1) und bricht bei Abweichung ab,
      **vor** jedem Login/Build/Push (Muster `release.yml`); ruft
      `make sdk-pack-csharp` auf; pusht mit `dotnet nuget push … --api-key
      ${{ secrets.NUGET_API_KEY }} --source
      https://api.nuget.org/v3/index.json`.
- [x] Jede `uses:`-Zeile ist auf einen vollständigen Commit-SHA gepinnt,
      mit Tag-Kommentar (`AGENTS.md` §3.8) — `actions/checkout` wiederverwendet
      denselben, bereits im ganzen Repo gepinnten SHA
      (`3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1`).
- [x] YAML-Struktur geprüft (Muster `release-image-scan`: Ruby-Stdlib-
      YAML-Parser, netzlos — `ruby -ryaml -e "YAML.load_file(...)"`, zusätzlich
      jeder `run:`-Block einzeln mit `bash -n` auf Syntaxfehler geprüft,
      Muster `release-version-und-workflow`).
- [x] `harness/README.md` §Werkzeuge bekommt die reale Zeile für
      `.github/workflows/sdk-csharp-release.yml` (kein Gate, Bindung auf
      [`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md)) —
      erst jetzt zulässig, weil der Workflow jetzt real existiert
      (`AGENTS.md` §4).
- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `harness/README.md` §Werkzeuge (siehe oben) — entfällt
      als eigener Punkt, da bereits oben als DoD-Kriterium geführt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — keine
      Beobachtung angefallen (siehe §7); bestehende Einträge geprüft und
      unverändert zutreffend.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (das Post-Push-Risiko UND das
      fehlende `NUGET_API_KEY`-Secret bleiben nach `AGENTS.md` §3.10 bzw.
      strukturell **weiter offen**, auch bei grüner DoD im Übrigen; Risiko 3
      ist eingetreten und aufgelöst, siehe §6).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-csharp-lh-fa-sst-009](../welle-sdk-csharp-lh-fa-sst-009.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.github/workflows/sdk-csharp-release.yml` | neu | Tag-Trigger, SemVer-/`.csproj`-Abgleich, `dotnet nuget push`. |
| `harness/README.md` §Werkzeuge | update | reale Workflow-Zeile. |

**Plan-Nachzug (vor dem Gate-Lauf, `AGENTS.md` §6 Schritt 5):** Beim
Schreiben kamen drei weitere Dateien hinzu, die im ursprünglichen Plan
oben nicht standen — Konsequenz der §6-Risiko-3-Entscheidung
(gemeinsames Skript statt zweier unabhängig driftender Regex-Kopien,
siehe §6/§7):

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/semver-regex.sh` | neu | gemeinsame SemVer-2.0-Regex, sourcebar, geteilt zwischen `release-tag-info.sh` und dem neuen `sdk-csharp-release-tag-info.sh` (§6 Risiko 3). |
| `tools/harness/release-tag-info.sh` | update | inline-Regex durch `source tools/harness/semver-regex.sh` ersetzt — Verhalten unverändert, `make test-release-tag-info` bleibt grün. |
| `tools/harness/sdk-csharp-release-tag-info.sh` | neu | validiert `sdk-csharp-v<SemVer>`-Tags, gibt `version=` aus (kein `latest=` — kein NuGet-`:latest`-Äquivalent, §1). |
| `tools/harness/run-sdk-csharp-release-tag-info-tests.sh` + `Makefile`-Target `test-sdk-csharp-release-tag-info` | neu | netzloser Tabellentest, Muster `run-release-tag-info-tests.sh`/`test-release-tag-info` (kein Gate). |

**Weitere Abweichung vom DoD-Beispieltext:** §2 nennt als Beispiel „per
`xmllint`/`grep` im gepinnten Toolchain-Container". Umgesetzt wurde
stattdessen ein bloßes `grep -oE`/`sed` **direkt auf dem Runner**, ohne
Docker-Umweg — dieselbe Form, mit der `release.yml` bereits
`docs/user/version.md` per `tr` direkt auf dem Runner liest (Zeile 74 dort),
keine neue Klasse. Der Grund gegen den Toolchain-Container: Der dort
gepinnte `golang:1.27-alpine` trägt kein `xmllint`; es nachzuinstallieren
bräuchte einen zusätzlichen `apk add`-Netzzugriff für eine Ein-Zeilen-Extraktion,
die `grep -oE '<Version>[^<]+</Version>'` gegen das feste, selbst
committete `<Version>…</Version>`-Tag-Muster ebenso zuverlässig leistet —
kein AGENTS.md-§3.1-Verstoß, weil es sich um eine bereits auf dem
ephemeren GitHub-Runner vorinstallierte Coreutils-Nutzung handelt, keine
lokale Toolchain-/SDK-Installation.

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-sdk-csharp-pack-werkzeug`
in `done/` liegt (siehe Welle-Plan §4 Reihenfolge — der Workflow
veröffentlicht, was das Pack-Werkzeug real erzeugt).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): entfällt aus
  heutiger Sicht — ein Workflow, analog `release.yml`/`image-scan.yml`.
- `in-progress` → `open` (blockiert — Carveout?): `slice-sdk-csharp-pack-werkzeug`
  liegt noch nicht in `done/`.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben — **das Post-Push-Risiko (§6) bleibt unabhängig davon nach
`AGENTS.md` §3.10 offen**, bis ein realer Tag-Push bestätigt oder ein roter
Befund mit Folgemaßnahme dokumentiert ist; dieser Slice kann trotzdem nach
`done/` gehen, wenn das Risiko explizit als „weiter offen" geführt wird
(dieselbe Praxis wie `release-image-scan`/`release-version-und-workflow`).

## 6. Risiken und offene Punkte

- **`AGENTS.md` §3.10 gilt unverändert:** `make gates` grün und ein
  plausibler YAML-Aufbau belegen nicht, dass der reale Post-Push-Lauf auf
  GitHub grün läuft (Registry-Zugangsdaten-Pfad zu NuGet.org, Runner-
  spezifisches `dotnet nuget push`-Verhalten). **Ausgang:** weiter offen,
  strukturell — bestätigt oder widerlegt erst durch einen realen Tag-Push
  `sdk-csharp-v0.1.0` (oder gleichwertig), dieselbe Klasse wie
  `BEO-PGC/github-actions-unverifizierbar-lokal` (bereits verkörpert als
  `AGENTS.md` §3.10).
- **`NUGET_API_KEY` existiert zum Zeitpunkt dieses Slice nicht als
  Repository-Secret** (analog `release-image-scan`s Blocker mit dem
  fehlenden GHCR-`:latest`-Tag) — ein realer Tag-Push würde am fehlenden
  Secret scheitern, kein Implementierungsfehler dieses Slice. **Ausgang:**
  weiter offen, strukturell — löst sich, sobald der Repository-Betreiber
  das Secret real anlegt (externe Handlung, außerhalb des Umfangs dieses
  Slice, siehe §1 Ausschluss).
- **Ein SemVer-Tag-Parser für `sdk-csharp-v*` teilt sich möglicherweise
  Code mit `release.yml`s bereits bestehendem `v*`-Parser** — eine
  naive Kopie würde beide Stellen unabhängig driften lassen. **Ausgang:**
  eingetreten, real entschieden beim Schreiben — **gemeinsames Skript**:
  Die SemVer-2.0-Regex selbst liegt jetzt in `tools/harness/semver-regex.sh`
  (sourcebar), das sowohl `tools/harness/release-tag-info.sh` (unverändertes
  Verhalten, `make test-release-tag-info` bleibt grün) als auch das neue
  `tools/harness/sdk-csharp-release-tag-info.sh` einbinden. Präfix-Handling
  (`v` vs. `sdk-csharp-v`) und Rückgabewert (`latest=` vs. kein `latest=`)
  bleiben bewusst je Skript eigenständig, weil beide Tag-Räume
  unterschiedliche Semantik tragen (Server-Stabilität vs. kein
  `:latest`-Äquivalent) — nur die tatsächlich identische Regel (die
  SemVer-2.0-Grammatik) wandert in die gemeinsame Datei, nicht die
  gesamten Skripte.

## 7. Closure-Notiz

- **Was hat funktioniert:** Die drei bestehenden Vorbild-Workflows
  (`release.yml`, `hub-description.yml`, `image-scan.yml`) trugen jeweils
  ein direkt übertragbares Formstück — Tag-Trigger + Tag-vs-Datei-Abgleich
  von `release.yml`, `workflow_call`-freie Minimalstruktur von
  `image-scan.yml`, Secret-Referenz-ohne-Anlage-Muster von beiden — die
  Zusammensetzung ging ohne Rückfragen. Der Nachweis, dass `sdk-csharp-v*`
  von keinem der drei bestehenden Workflows mitgetriggert wird, ließ sich
  direkt an den `tags`/`tags-ignore`-Zeilen ablesen (`ci.yml`/`e2e.yml`:
  `tags-ignore: ['**']` schließt jeden Tag-Push ein, `release.yml`:
  `tags: ['v*']` matcht ein `sdk-csharp-`-Präfix nicht).
- **Was ging anders als geplant:** Zwei Abweichungen vom DoD-Beispieltext,
  beide als Plan-Nachzug in §3 begründet: (1) drei zusätzliche Dateien
  (`semver-regex.sh`, `sdk-csharp-release-tag-info.sh`,
  `run-sdk-csharp-release-tag-info-tests.sh` + Makefile-Target) als
  Auflösung von §6 Risiko 3 — nicht im ursprünglichen §3-Plan gelistet,
  weil die Entscheidung explizit erst „beim Schreiben" fallen sollte; (2)
  die `.csproj`-Versionsprüfung läuft als bloßes `grep -oE`/`sed` direkt auf
  dem Runner statt im gepinnten Toolchain-Container (Begründung: der dort
  gepinnte `golang:1.27-alpine` trägt kein `xmllint`, ein `apk add`-Umweg
  für eine Ein-Zeilen-Extraktion wäre unverhältnismäßig — dieselbe Form wie
  `release.yml`s eigener `tr`-Aufruf gegen `docs/user/version.md`).
- **Steering-Loop-Eintrag:** kein neuer Sensor — die vorhandene
  Ruby-Stdlib-YAML-Struktur-/`bash -n`-Prüfung (Muster
  `release-version-und-workflow`) reichte erneut aus. Kein Guide/Sensor
  geschärft.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung
  angefallen. `BEO-PGC/github-actions-unverifizierbar-lokal` und
  `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit` wurden geprüft
  (§8) und bleiben unverändert zutreffend, kein neuer Beleg nötig — das
  Post-Push-Risiko dieses Slice ist dieselbe, bereits verkörperte Klasse
  (`AGENTS.md` §3.10), kein drittes Auftreten einer neuen Verstoßklasse.
- **Folge-Slices:** keine aus diesem Slice selbst erwartet — letzter Slice
  der Welle.
- **Risiken aus §6:** Post-Push-Lauf gegen NuGet.org — **weiter offen**,
  strukturell (`AGENTS.md` §3.10), bestätigt/widerlegt erst durch einen
  realen Tag-Push. `NUGET_API_KEY`-Secret fehlt — **weiter offen**,
  strukturell, externe Kontohandlung außerhalb dieses Slice-Umfangs.
  SemVer-Parser-Code-Teilung — **eingetreten und aufgelöst**: gemeinsames
  Skript `tools/harness/semver-regex.sh`, siehe §6 für die volle
  Begründung.
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-csharp-lh-fa-sst-009](../welle-sdk-csharp-lh-fa-sst-009.md)
  (noch offen) — die Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `.github/workflows/`
(CI/CD-Pipeline) — bereits mehrfach berührt (`ci.yml`, `e2e.yml`,
`release.yml`, `image-scan.yml`, `upstream-drift.yml`,
`hub-description.yml`), Schwelle erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/github-actions-unverifizierbar-lokal` (bereits verkörpert als
`AGENTS.md` §3.10) betrifft diese Sub-Area unmittelbar (siehe §6);
`BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit` geprüft — betrifft
mehrfach rote, aber nie behobene advisory Workflows, nicht diesen Slice
(der Workflow ist neu, keine bestehende Rot-Historie); kein weiterer
Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
