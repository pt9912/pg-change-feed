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

**Verantwortlich:** — bis zur Priorisierung.

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

- [ ] `.github/workflows/sdk-csharp-release.yml` existiert: Trigger
      ausschließlich `push: tags: ['sdk-csharp-v*']` (`ci.yml`/`e2e.yml`/
      `release.yml` schließen diesen Tag-Namensraum durch ihre jeweils
      eigenen `tags`/`tags-ignore`-Filter aus, kein Doppellauf — Beleg im
      Bericht dieses Slice); validiert den Tag-Suffix (nach `sdk-csharp-v`)
      strikt gegen SemVer 2.0; liest `<Version>` aus
      `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj` (z. B.
      per `xmllint`/`grep` im gepinnten Toolchain-Container, `AGENTS.md`
      §3.1) und bricht bei Abweichung ab, **vor** jedem Login/Build/Push
      (Muster `release.yml`); ruft `make sdk-pack-csharp` auf; pusht mit
      `dotnet nuget push … --api-key ${{ secrets.NUGET_API_KEY }} --source
      https://api.nuget.org/v3/index.json`.
- [ ] Jede `uses:`-Zeile ist auf einen vollständigen Commit-SHA gepinnt,
      mit Tag-Kommentar (`AGENTS.md` §3.8).
- [ ] YAML-Struktur geprüft (Muster `release-image-scan`: Ruby-Stdlib-
      YAML-Parser oder gleichwertig, netzlos).
- [ ] `harness/README.md` §Werkzeuge bekommt die reale Zeile für
      `.github/workflows/sdk-csharp-release.yml` (kein Gate, Bindung auf
      [`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md)) —
      erst jetzt zulässig, weil der Workflow jetzt real existiert
      (`AGENTS.md` §4).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `harness/README.md` §Werkzeuge (siehe oben) — entfällt
      als eigener Punkt, da bereits oben als DoD-Kriterium geführt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`; keine Beobachtung angefallen
      ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (das Post-Push-Risiko bleibt
      nach `AGENTS.md` §3.10 strukturell **weiter offen**, auch bei
      grüner DoD im Übrigen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-csharp-lh-fa-sst-009](../welle-sdk-csharp-lh-fa-sst-009.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.github/workflows/sdk-csharp-release.yml` | neu | Tag-Trigger, SemVer-/`.csproj`-Abgleich, `dotnet nuget push`. |
| `harness/README.md` §Werkzeuge | update | reale Workflow-Zeile. |

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
  weiter offen, zu entscheiden beim Schreiben (gemeinsames Skript unter
  `tools/harness/` vs. bewusst getrennte, unabhängige Prüf-Schritte je
  Workflow — beide Formen sind mit `AGENTS.md` vereinbar).

## 7. Closure-Notiz

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <Guide oder Sensor geschärft/ergänzt, oder
  „kein neuer Sensor" — je nach Lauf>.
- **Beobachtungs-Register (`../observations/`):** <neu angelegt | Beleg
  ergänzt | keine Beobachtung angefallen>.
- **Folge-Slices:** keine aus diesem Slice selbst erwartet — letzter Slice
  der Welle.
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6; die
  beiden ersten bleiben nach `AGENTS.md` §3.10 strukturell weiter offen>.
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
