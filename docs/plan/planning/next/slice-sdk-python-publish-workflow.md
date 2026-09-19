# Slice sdk-python-publish-workflow: PyPI-Publish-Workflow (`sdk-python-release.yml`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-python-lh-fa-sst-009](../welle-sdk-python-lh-fa-sst-009.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md),
[`ADR-0107`](../../adr/0107-python-pypi-zweites-sdk-package.md)
Festlegung 5 (Trigger, Secret, Tag-Präfix), [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
Entscheidung 3/8 (Release-Tag-Trigger, Registry-Secret-Muster — Vorbild,
wie `ADR-0106` es bereits für NuGet nutzte),
[`ADR-0108`](../../adr/0108-python-sdk-uv-statt-build-twine.md)
§Entscheidung Festlegung 1 (Publish-Frontend `uv publish` statt
`twine upload`, Token-Konsum über `UV_PUBLISH_TOKEN`/`--token`) —
superseded die `twine`-Publish-Zeile aus `ADR-0107` Festlegung 5.

**Berührte Spec-Stellen:** — (Prozess-/CI-Artefakt ohne eigene
`SPEC-*`-Kennung, analog `sdk-csharp-release.yml`).

**Verantwortlich:** Implementer-Agent (dietmar.burkard@nerdware.dev), ab 2026-09-19.

**Autor:** Planner-Agent, direkt beauftragt (`ADR-0107` §Konsequenzen
Folgepflicht 1/5). **Datum:** 2026-09-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein neuer, Netz-bindender GitHub-Actions-Workflow
`.github/workflows/sdk-python-release.yml`, Trigger ausschließlich
`push: tags: ['sdk-python-v*']`, der den Tag strikt gegen PEP 440
validiert, ihn gegen die in `pyproject.toml` geführte `version` abgleicht
(Abbruch bei Abweichung, vor jedem Build/Push — Muster `release.yml`s
Tag-gegen-`version.md`-Abgleich bzw. `sdk-csharp-release.yml`s
Tag-gegen-`.csproj`-Abgleich, hier gegen die `pyproject.toml`), die
Artefakte aus `make sdk-pack-python` erzeugt und per `uv publish`
(Token über `UV_PUBLISH_TOKEN=${{ secrets.PYPI_API_TOKEN }}`,
[`ADR-0108`](../../adr/0108-python-sdk-uv-statt-build-twine.md)
§Entscheidung Festlegung 1) veröffentlicht.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Anlage des `PYPI_API_TOKEN`-Repository-Secrets** — der Workflow
  **referenziert** das Secret, legt es nicht an (`ADR-0107` Festlegung 5:
  „ein Repository-Secret, das der Workflow referenziert, nicht anlegt" —
  dieselbe Klasse wie `NUGET_API_KEY`/`DOCKERHUB_USERNAME`/
  `DOCKERHUB_TOKEN`). Die reale Anlage ist eine externe, außerhalb des
  Repos liegende Handlung des Repository-Betreibers.
- **Denselben Tag-Namensraum wie `release.yml` (`v*`) oder
  `sdk-csharp-release.yml` (`sdk-csharp-v*`)** — `ADR-0107` Festlegung 5
  begründet ausdrücklich den eigenen Präfix `sdk-python-v*`, um die
  unabhängige SDK-Versionierung (Festlegung 4) nicht mit den beiden
  bestehenden Release-Räumen kollidieren zu lassen.
- **PyPI Trusted Publishing (OIDC)** — `ADR-0107` Festlegung 5 verwirft
  das ausdrücklich für das Erst-Release (Henne-Ei-Registrierungs-
  voraussetzung); das API-Token-Muster ist die gewählte Erstwahl.
  Trusted Publishing bleibt ein Re-Evaluierungs-Anlass für **nach** dem
  ersten erfolgreichen Publish (`ADR-0107` §Re-Evaluierungs-Trigger 4).
- **`:latest`-Äquivalent für das SDK** — PyPI kennt kein Analogon zu einem
  Docker-`:latest`-Tag; ein veröffentlichtes Package ist über seine exakte
  Version referenzierbar.
- **GitHub-Release-Eintrag für das SDK** — nicht Teil von `ADR-0107`
  Festlegung 5; kein Bestandteil dieser Folgepflicht.
- **Der reale, grüne Post-Push-Lauf gegen PyPI** — bleibt nach
  [`AGENTS.md`](../../../../AGENTS.md) §3.10 ein offenes Risiko dieses
  Slice bis zum ersten tatsächlichen Tag-Push (siehe §6); `make gates`
  grün und ein plausibler YAML-Aufbau sind keine Ersatz-Bestätigung.

## 2. Definition of Done

- [ ] `.github/workflows/sdk-python-release.yml` existiert: Trigger
      ausschließlich `push: tags: ['sdk-python-v*']` (`ci.yml`/`e2e.yml`/
      `release.yml`/`sdk-csharp-release.yml` schließen diesen
      Tag-Namensraum durch ihre jeweils eigenen `tags`/`tags-ignore`-Filter
      aus, kein Doppellauf — Beleg im Bericht dieses Slice); validiert den
      Tag-Suffix (nach `sdk-python-v`) strikt gegen PEP 440; liest
      `version` aus `sdks/python/pgchangefeed/pyproject.toml` (per
      `grep -oE`/`sed` direkt auf dem Runner, analog
      `sdk-csharp-release.yml`s `.csproj`-Lesart, `AGENTS.md` §3.1) und
      bricht bei Abweichung ab, **vor** jedem Login/Build/Push; ruft
      `make sdk-pack-python` auf; pusht mit `uv publish` (Token über die
      Umgebungsvariable `UV_PUBLISH_TOKEN: ${{ secrets.PYPI_API_TOKEN }}`,
      [`ADR-0108`](../../adr/0108-python-sdk-uv-statt-build-twine.md)
      §Entscheidung Festlegung 1 — kein `-u __token__ -p`-Flag-Paar wie
      bei `twine`).
- [ ] Jede `uses:`-Zeile ist auf einen vollständigen Commit-SHA gepinnt,
      mit Tag-Kommentar (`AGENTS.md` §3.8) — `actions/checkout`
      wiederverwendet denselben, bereits im ganzen Repo gepinnten SHA.
- [ ] YAML-Struktur geprüft (Muster `sdk-csharp-release.yml`: Ruby-Stdlib-
      YAML-Parser, netzlos — `ruby -ryaml -e "YAML.load_file(...)"`,
      zusätzlich jeder `run:`-Block einzeln mit `bash -n` auf
      Syntaxfehler geprüft).
- [ ] `harness/README.md` §Werkzeuge bekommt die reale Zeile für
      `.github/workflows/sdk-python-release.yml` (kein Gate, Bindung auf
      [`ADR-0107`](../../adr/0107-python-pypi-zweites-sdk-package.md)) —
      erst jetzt zulässig, weil der Workflow jetzt real existiert
      (`AGENTS.md` §4).
- [ ] `docs/user/releasing.md` bekommt einen eigenen Abschnitt/eine
      Tabellenzeile für den Python-SDK-Release-Weg (Tag-Präfix
      `sdk-python-v*`, Secret `PYPI_API_TOKEN`, Vertrag-Datei
      `pyproject.toml`), **analog dem bereits bestehenden C#-Eintrag** —
      als eigener, vorab eingeplanter DoD-Punkt dieses Slice, **nicht**
      erst nach einem Reviewer-Finding (`BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen`,
      1×, entstanden bei `slice-sdk-csharp-publish-workflow`; diese
      Planung nimmt die Lehre vorweg, um kein zweites Auftreten dieser
      Klasse zu erzeugen).
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `harness/README.md` §Werkzeuge und
      `docs/user/releasing.md` (siehe oben) — entfallen als eigene Punkte,
      da bereits oben als DoD-Kriterien geführt.
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
      [welle-sdk-python-lh-fa-sst-009](../welle-sdk-python-lh-fa-sst-009.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure — und
      ist zugleich der **letzte** Slice der Welle, also auch der letzte
      Anlass, die Welle-Closure selbst anzustoßen.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `.github/workflows/sdk-python-release.yml` | neu | Tag-Trigger, PEP-440-/`pyproject.toml`-Abgleich, `uv publish`. |
| `harness/README.md` §Werkzeuge | update | reale Workflow-Zeile. |
| `docs/user/releasing.md` | update | eigener Abschnitt/Tabellenzeile für den Python-SDK-Release-Weg, analog dem C#-Eintrag — **vorab eingeplant**, siehe §2. |
| `tools/harness/sdk-python-release-tag-info.sh` (Arbeitsname) | neu | validiert `sdk-python-v<PEP 440>`-Tags, gibt `version=` aus (kein `latest=`) — Struktur-Vorbild `tools/harness/sdk-csharp-release-tag-info.sh`; PEP-440-Grammatik ist **nicht** identisch mit SemVer 2.0 (Pre-/Post-Release-Suffixe unterscheiden sich), deshalb **kein** Sourcing aus `tools/harness/semver-regex.sh` ohne vorherige Prüfung, ob die dortige Regex den hier tatsächlich genutzten einfachen Fall (`MAJOR.MINOR.PATCH`, `ADR-0107` Festlegung 4) bereits abdeckt — diese Entscheidung fällt beim Schreiben, nicht hier. |
| `tools/harness/run-sdk-python-release-tag-info-tests.sh` + `Makefile`-Target `test-sdk-python-release-tag-info` | neu | netzloser Tabellentest, Muster `run-sdk-csharp-release-tag-info-tests.sh`/`test-sdk-csharp-release-tag-info` (kein Gate). |

**Ansatz:** Strukturvorbild ist `slice-sdk-csharp-publish-workflow`
(bereits `done/`) — Tag-Trigger + Tag-vs-Datei-Abgleich, Secret-
Referenz-ohne-Anlage-Muster, YAML-Struktur-/`bash -n`-Prüfung. Der
zentrale Unterschied: PEP 440 statt SemVer 2.0 als Ziel-Grammatik (`ADR-0107`
Festlegung 4) — der Implementer prüft beim Schreiben, ob die im einfachen
Fall (`MAJOR.MINOR.PATCH`) verwendete Teilmenge sich mit
`tools/harness/semver-regex.sh` deckt, bevor er eine zweite,
eigenständige Regex einführt oder die bestehende wiederverwendet.

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-sdk-python-pack-werkzeug`
in `done/` liegt (siehe Welle-Plan §4 Reihenfolge — der Workflow
veröffentlicht, was das Pack-Werkzeug real erzeugt).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): entfällt aus
  heutiger Sicht — ein Workflow, analog `sdk-csharp-release.yml`.
- `in-progress` → `open` (blockiert — Carveout?): `slice-sdk-python-pack-werkzeug`
  liegt noch nicht in `done/`.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben — **das Post-Push-Risiko (§6) bleibt unabhängig davon nach
`AGENTS.md` §3.10 offen**, bis ein realer Tag-Push bestätigt oder ein roter
Befund mit Folgemaßnahme dokumentiert ist; dieser Slice kann trotzdem nach
`done/` gehen, wenn das Risiko explizit als „weiter offen" geführt wird
(dieselbe Praxis wie `slice-sdk-csharp-publish-workflow`).

## 6. Risiken und offene Punkte

- **`AGENTS.md` §3.10 gilt unverändert:** `make gates` grün und ein
  plausibler YAML-Aufbau belegen nicht, dass der reale Post-Push-Lauf auf
  GitHub grün läuft (Registry-Zugangsdaten-Pfad zu PyPI, Runner-
  spezifisches `uv publish`-Verhalten — zusätzlich verschärft durch die
  in `ADR-0108` §Kontext benannte 0.x-Versionierung von `uv` selbst).
  **Ausgang:** weiter offen,
  strukturell — bestätigt oder widerlegt erst durch einen realen Tag-Push
  `sdk-python-v0.1.0` (oder gleichwertig), dieselbe Klasse wie
  `BEO-PGC/github-actions-unverifizierbar-lokal` (bereits verkörpert als
  `AGENTS.md` §3.10).
- **`PYPI_API_TOKEN` existiert zum Zeitpunkt der Slice-Planung
  voraussichtlich nicht als Repository-Secret** (analog dem
  `NUGET_API_KEY`-Blocker bei der C#-Welle) — ein realer Tag-Push würde am
  fehlenden Secret scheitern, kein Implementierungsfehler dieses Slice.
  **Ausgang:** weiter offen — bei der C#-Welle stellte sich real heraus,
  dass das Secret zum Zeitpunkt der Verifikation bereits (vermutlich extern)
  angelegt worden war; ob das für `PYPI_API_TOKEN` ebenso gilt, ist beim
  Implementieren/Verifizieren dieses Slice real zu prüfen (`gh secret
  list`), nicht anzunehmen.
- **Der PyPI-Paketname `pgchangefeed` könnte bereits vergeben sein** — im
  Gegensatz zu NuGet (real geprüft frei) ist die PyPI-Namensraum-Freiheit
  für diesen Arbeitsnamen noch nicht bestätigt. **Ausgang:** weiter offen
  — vor dem ersten realen Tag-Push prüft der Implementer
  `pip index versions pgchangefeed` bzw. die PyPI-Weboberfläche; ist der
  Name vergeben, braucht es einen alternativen Arbeitsnamen (z. B.
  `pgchangefeed-client`), was `pyproject.toml` und diesen Workflow gleich
  beträfe — kein Rückbau bereits gelieferter Fläche, nur ein Namens-Nachzug.
- **PEP-440- vs. SemVer-2.0-Regex-Teilung** — analog dem bei der C#-Welle
  real aufgetretenen Risiko (gemeinsames `tools/harness/semver-regex.sh`
  vs. eigenständige Kopie), hier verschärft durch die abweichende
  Grammatik (PEP 440 ≠ SemVer 2.0 im vollen Funktionsumfang, nur im hier
  genutzten einfachen Fall kompatibel, `ADR-0107` Festlegung 4). **Ausgang:**
  weiter offen, zu entscheiden beim Schreiben (siehe §3 Ansatz).

## 7. Closure-Notiz

<…>

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `.github/workflows/`
(CI/CD-Pipeline) — bereits mehrfach berührt (`ci.yml`, `e2e.yml`,
`release.yml`, `image-scan.yml`, `upstream-drift.yml`,
`hub-description.yml`, `sdk-csharp-release.yml`), Schwelle erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/github-actions-unverifizierbar-lokal` (bereits verkörpert als
`AGENTS.md` §3.10) betrifft diese Sub-Area unmittelbar (siehe §6);
`BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen` (offen,
1×) betrifft diesen Slice unmittelbar — als DoD-Punkt oben bereits
vorab aufgenommen, um ein zweites Auftreten zu vermeiden;
`BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit` geprüft — betrifft
mehrfach rote, aber nie behobene advisory Workflows, nicht diesen Slice
(der Workflow ist neu, keine bestehende Rot-Historie).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
