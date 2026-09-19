# ADR-0108: `uv` als Bau-/Publish-Werkzeug für das Python-SDK

**Status:** Accepted — Supersedes [`ADR-0107`](0107-python-pypi-zweites-sdk-package.md)
in **einer** Klausel: §Entscheidung Festlegung 5, erster Aufzählungspunkt
(„**Paket-Werkzeug: `build` (PyPA-Referenzwerkzeug) + `setuptools`-Backend,
nicht Poetry oder Hatch**" samt der zugehörigen `twine`-Publish-Zeile und der
Bewertung in §Verglichene Alternativen Tabelle E) — ausschließlich die Wahl
des **Build-/Publish-Frontends**. Alles Übrige von
[`ADR-0107`](0107-python-pypi-zweites-sdk-package.md) bleibt **hiermit
bestätigt** und wird nicht wiederholt: Festlegung 1 (Umfang v1 nur HTTP-API),
Festlegung 2 (Vertriebsweg PyPI, Arbeitsname `pgchangefeed`), Festlegung 3
(Ort `sdks/python/`, Import-Grenze), Festlegung 4 (Versionierung PEP 440,
Start `0.x.y`), der übrige Teil von Festlegung 5 (Docker-only bauen, kein
Gate, Secret-**Klasse** `PYPI_API_TOKEN`, eigener Tag-Namensraum
`sdk-python-v*`, eigener Workflow `sdk-python-release.yml`), §Kontext (die
Ist-Stand-Messung, die Risiko-Asymmetrie ggü. `ADR-0106`), §Verglichene
Alternativen Tabelle A/B/C/D vollständig sowie Tabelle E nur in ihrer
E2/E3-Bewertung (Poetry/Hatch bleiben aus denselben Gründen verworfen), alle
fünf Punkte in §Entscheidung Festlegung 6 („Was diese ADR nicht ändert"),
§Konsequenzen Folgepflicht 1–5, §Re-Evaluierungs-Trigger 1–4 und §Geschichte.
`sdks/python/pgchangefeed/pyproject.toml`s `[build-system]
build-backend = "setuptools.build_meta"` bleibt **unverändert** — diese ADR
ersetzt ein Frontend, kein Backend-Format (§Entscheidung Festlegung 2).

**Datum:** 2026-09-19

**Autor:** pt9912 (Architect-Rolle, Nutzerentscheidung vom 2026-09-19 im Chat
nach einer Rückfrage per `AskUserQuestion` mit Empfehlung: eine echte Lücke
in [`ADR-0107`](0107-python-pypi-zweites-sdk-package.md) §Verglichene
Alternativen Tabelle E — `uv` (Astral) wurde dort **nicht** geprüft. Der
Nutzer hat sich für den Wechsel entschieden, nicht für den Verbleib bei
`build`+`twine`. Modul 8 §Rollen-Regeln: „Architect schreibt"; die
Immutabilitäts-Regel für `Accepted`-ADRs erzwingt eine **neue** ADR mit
`Supersedes`, keine Nachbesserung an
[`ADR-0107`](0107-python-pypi-zweites-sdk-package.md) selbst, `AGENTS.md`
§3.5)

**Bezug:** [`LH-FA-SST-009`](../../../spec/lastenheft.md) (Client-Bibliotheken/
SDKs), [`ADR-0107`](0107-python-pypi-zweites-sdk-package.md) (teilweise
superseded — siehe Kopf), [`ADR-0106`](0106-csharp-nuget-erstes-sdk-package.md)
(C#/NuGet-SDK — von dieser ADR **unberührt**, eigener Bau-/Publish-Mechanismus),
[`AGENTS.md`](../../../AGENTS.md) §3.1 (Docker-only — gilt unverändert: `uv`
läuft ausschließlich im Container, kein Host-`uv`/`pip`/`python`) §3.5
(Accepted-ADR-Immutabilität — der Grund für diese Datei statt einer
In-Place-Korrektur) §3.6 (Gates nur per ADR gelockert — hier irrelevant,
kein Gate berührt) §3.12 (Herkunft von Aussagen — die Reifegrad-Einschätzung
von `uv publish` unten trägt ihren Beleg), `sdks/python/pgchangefeed/pyproject.toml`
(das unveränderte `[build-system]`), `sdks/python/Dockerfile` (die Bau-Stufe,
die den Frontend-Wechsel real trägt, sobald sie entsteht),
`docs/plan/planning/open/slice-sdk-python-pack-werkzeug.md` und
`docs/plan/planning/open/slice-sdk-python-publish-workflow.md` (Planungsträger
mit textuellem `build`+`twine`-Verweis — Nachzug ist Folgepflicht des
Planners, siehe §Konsequenzen), `docs/plan/planning/welle-sdk-python-lh-fa-sst-009.md`
(dieselbe Klasse Verweis auf Wellen-Ebene)

**Schärft:** [`LH-FA-SST-009.a`](../../../spec/pflichtenheft.md) — dieselbe
Pflichtenheft-Stelle bleibt für den Frontend-Mechanismus weiterhin über diese
ADR-Kette (jetzt: [`ADR-0107`](0107-python-pypi-zweites-sdk-package.md) +
diese ADR) beantwortet; an der Antwort selbst — Python/PyPI als zweite
Sprache/zweiter Vertriebsweg — ändert diese ADR nichts.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**Die Lücke.** [`ADR-0107`](0107-python-pypi-zweites-sdk-package.md)
§Verglichene Alternativen Tabelle E vergleicht drei Optionen für das
Python-Paket-Werkzeug: `build`(PyPA)+`setuptools`-Backend+`twine` (gewählt),
Poetry, Hatch. `uv` (Astral) — inzwischen ein etabliertes,
aktiv entwickeltes Python-Bau-/Paket-/Publish-Werkzeug — wurde in dieser
Tabelle **nicht** genannt. Das ist eine echte Alternativen-Lücke, keine
Verwerfung: eine Option, die nicht in der Tabelle steht, wurde nicht
gegen die anderen abgewogen, sie wurde übersehen.

**Der Ist-Stand — gemessen, nicht erinnert (heutiger Zug).**

| Prozedur | Ergebnis |
|---|---|
| `find sdks/python -maxdepth 3 -type f` | `sdks/python/Dockerfile`, `sdks/python/README.md`, `sdks/python/.gitignore`, `sdks/python/pgchangefeed/pyproject.toml`, `sdks/python/pgchangefeed/tests/test_options.py` — das Projektgerüst existiert bereits (`slice-sdk-python-projektgeruest`, real reviewed und verifiziert, siehe Commit-Historie) |
| `cat sdks/python/pgchangefeed/pyproject.toml` | `[build-system] requires = ["setuptools>=68"]` / `build-backend = "setuptools.build_meta"` — **genau** das Backend, das [`ADR-0107`](0107-python-pypi-zweites-sdk-package.md) Festlegung 5 gewählt hat; kein `python -m build`-Aufruf irgendwo im committeten Baum |
| `cat sdks/python/Dockerfile` | trägt nur `pip install -e ".[test]"` + `pytest`; die Pack-Stufe (`python -m build`) ist laut eigenem Kommentar explizit **noch nicht** implementiert — Folgepflicht von `slice-sdk-python-pack-werkzeug` |
| `grep -rn "sdk-pack-python\|twine\|python -m build"` über `docs/plan/planning/` | drei Fundstellen, alle in **noch offenen** Planning-Dokumenten (`docs/plan/planning/open/slice-sdk-python-pack-werkzeug.md`, `docs/plan/planning/open/slice-sdk-python-publish-workflow.md`, `docs/plan/planning/welle-sdk-python-lh-fa-sst-009.md`) — **kein** bereits gebauter/gemergter Code referenziert `build` oder `twine` |

**Was das für den Zuschnitt dieser ADR bedeutet.** Es gibt real **nichts zu
migrieren** — nur Planungstext, der auf die jetzt abgelöste Klausel
verweist, bevor er umgesetzt wird. Diese ADR ändert deshalb keinen Code und
keine bereits geschlossene Slice-Datei; sie ändert die **Form**, die der
noch offene Umsetzungs-Zug einschlägt (§Konsequenzen Folgepflicht).

**Recherche zur Reife von `uv build`/`uv publish` — real durchgeführt, mit
Beleg (`AGENTS.md` §3.12, keine Behauptung ohne Anker).**

- **Live abgerufen** (`curl`, heute, `main`-Branch von `astral-sh/uv` auf
  GitHub): `docs/guides/package.md`. Zentrale, wörtlich geprüfte Aussagen:
  - „If your project does not include a `[build-system]` definition …, `uv`
    will not build it during `uv sync` … but will fall back to the legacy
    setuptools build system during `uv build`." — im Umkehrschluss, real
    geprüft an der eigenen `pyproject.toml`: **existiert** ein
    `[build-system]`-Block (unser Fall, `setuptools.build_meta`), verwendet
    `uv build` genau **dieses** Backend über die PEP-517-Schnittstelle,
    ohne es zu ersetzen. Das ist der tragende Unterschied zu Poetry/Hatch
    (§Verglichene Alternativen).
  - „Set a PyPI token with `--token` or `UV_PUBLISH_TOKEN` … For publishing
    to PyPI from GitHub Actions or another Trusted Publisher, you don't
    need to set any credentials." — Trusted Publishing (OIDC) ist **nativ**
    in `uv publish` eingebaut, kein separates Plugin nötig (anders als bei
    `twine`, wo Trusted Publishing über einen zusätzlichen
    OIDC-Austausch-Schritt vor dem Upload verdrahtet werden muss).
  - Ein eigener Abschnitt verweist auf eine vollständige „GitHub Guide:
    publishing to PyPI" — ein offiziell dokumentierter, kein experimenteller
    Pfad.
  - `uv build --no-sources` wird für Publish-Zwecke ausdrücklich empfohlen
    („we recommend running `uv build --no-sources` to ensure that the
    package builds correctly when `tool.uv.sources` is disabled").
- **Live abgerufen** (`curl`, GitHub Releases API): aktuelle `uv`-Version
  `0.12.17` (`tag_name`). `uv` selbst ist **weiterhin 0.x-versioniert** —
  Astral hat, soweit in der abgerufenen Doku sichtbar, **keine** explizite
  „`uv publish` ist seit Version X stabil"-Zusage formuliert.
- **Offener Punkt, ausdrücklich als solcher markiert (keine Tatsachenbehauptung):**
  Der Reifegrad von `uv publish` in einem echten, produktiven PyPI-Release
  bleibt bis zum ersten realen `sdk-python-release.yml`-Lauf unbewiesen —
  dieselbe Klasse Unsicherheit, die `AGENTS.md` §3.10 für jeden neuen
  GitHub-Actions-Workflow ohnehin schon benennt. Diese ADR behandelt die
  Doku-Lage (nativ dokumentiert, offizielle CI-Anleitung vorhanden) als
  hinreichenden Grund für die Wahl, **nicht** als Beweis eines fehlerfreien
  Produktivlaufs.

**Digest-Recherche für den Docker-Bezug — real gemessen, heute:**

| Kandidat | Befehl | Ergebnis |
|---|---|---|
| `ghcr.io/astral-sh/uv:0.12.17` (amd64-Manifest) | `docker buildx imagetools inspect` | Index-Digest `sha256:10787c682e4184e4f290de1171fd4703dc63de99221f10fe1c99002ce7fa9acc`; amd64-Plattform-Manifest `sha256:1194b357d63b7bea6c121d8eef5d08d29c26ffbcfc64ec5ebcefd642ea759edb` |
| `ghcr.io/astral-sh/uv:latest` | dieselbe Prozedur | **identischer** Index-Digest — bestätigt `0.12.17` als aktuellen Release-Stand |

## Entscheidung

Wir wählen: **`uv build` + `uv publish` ersetzen `build`(PyPA)+`twine` als
Bau-/Publish-**Frontend** für das Python-SDK. Das `setuptools`-Backend in
`pyproject.toml` bleibt unverändert.** Vier Festlegungen (nur die Klausel,
die diese ADR ändert):

### 1 — Frontend-Wechsel: `uv build --no-sources` + `uv publish`

- **Bauen:** `uv build --no-sources` statt `python -m build` — erzeugt
  identisch `.tar.gz` (sdist) und `.whl` (Wheel) im `dist/`-Unterverzeichnis;
  `--no-sources` folgt der eigenen Empfehlung der `uv`-Dokumentation für
  Publish-Korrektheit (§Kontext).
- **Publizieren:** `uv publish` statt `twine upload --repository pypi dist/*
  -u __token__ -p $PYPI_API_TOKEN`. Die Secret-**Klasse** bleibt exakt die
  aus [`ADR-0107`](0107-python-pypi-zweites-sdk-package.md) Festlegung 5
  (Arbeitsname `PYPI_API_TOKEN`, ein Repository-Secret, das der Workflow
  referenziert, nicht anlegt) — nur die Konsum-Form ändert sich: `uv publish`
  liest den Token über `UV_PUBLISH_TOKEN` (oder `--token`) statt über
  `twine`s `-u __token__ -p`-Flags.
- **Trusted Publishing bleibt vertagt, aus demselben Grund wie in
  [`ADR-0107`](0107-python-pypi-zweites-sdk-package.md) Festlegung 5**
  (Henne-Ei: erfordert eine vorherige PyPI-Projekt-Registrierung, die es für
  ein noch nie veröffentlichtes Package nicht geben kann) —
  [`ADR-0107`](0107-python-pypi-zweites-sdk-package.md) §Re-Evaluierungs-Trigger
  4 bleibt unverändert gültig, jetzt mit einem Werkzeug, das native
  OIDC-Unterstützung **bereits mitbringt**, sobald der Trigger eintritt
  (kein zusätzliches Plugin nötig).

### 2 — Das Backend ändert sich nicht

`sdks/python/pgchangefeed/pyproject.toml`s `[build-system]
build-backend = "setuptools.build_meta"` bleibt **exakt**, wie es real im
Baum steht (§Kontext Ist-Stand) — diese ADR fordert **keine** Änderung an
dieser Datei. `uv build` benutzt ein vorhandenes PEP-517-Backend über die
Standard-Schnittstelle, es ersetzt es nicht (§Kontext, Beleg aus
`docs/guides/package.md`). **Das ist der zentrale Unterschied zu Poetry und
Hatch**, den [`ADR-0107`](0107-python-pypi-zweites-sdk-package.md)
§Verglichene Alternativen E2/E3 gegen genau diese beiden Werkzeuge vorbrachte
(eigenes Lockfile-/Environment-/Projektstruktur-Format) — dieser Einwand
**trifft `uv` nicht**, weil `uv` hier ausschließlich als Frontend auftritt,
nicht als Backend- oder Projektformat-Ersatz. Aus demselben Grund bleibt
Poetry/Hatch weiterhin verworfen (§Verglichene Alternativen E4).

### 3 — Docker-Bezug von `uv`: digest-gepinntes Kopier-Image, kein `pip install uv`

`sdks/python/Dockerfile` bezieht `uv` über einen zusätzlichen
Multi-Stage-`COPY`-Schritt aus Astrals eigenem, offiziellem Werkzeug-Image:

```dockerfile
COPY --from=ghcr.io/astral-sh/uv:0.12.17@sha256:10787c682e4184e4f290de1171fd4703dc63de99221f10fe1c99002ce7fa9acc /uv /uvx /bin/
```

— **nicht** `pip install uv` innerhalb des bereits gepinnten
`python:3.14-slim`-Basis-Images. Begründung:

- **Kein Henne-Ei-Bezug über PyPI zur Build-Zeit für das Werkzeug, das
  selbst gegen PyPI publiziert.** `pip install uv` zöge `uv` als
  gewöhnliche PyPI-Abhängigkeit zur Build-Zeit — funktional zulässig, aber
  eine unnötige zusätzliche PyPI-Vertrauenskette für ein Werkzeug, das
  Astral bereits als fertiges, digest-pinnbares Binary-Image ausliefert.
- **Folgt der bestehenden Haus-Kultur unverändert:** jede Stufe der
  Wurzel-`Dockerfile` und jedes bereits bestehende Sprach-Dockerfile
  (`sdks/csharp/Dockerfile`, `examples/csharp/Dockerfile`,
  `examples/kotlin/Dockerfile`) ist digest-gepinnt; ein `COPY --from=`
  gegen ein weiteres digest-gepinntes Image ist dieselbe Form, keine neue
  Klasse.
- **Bleibt Docker-only** (`AGENTS.md` §3.1) — kein Host-`uv`, kein
  Host-`pip`, kein Host-`python`; der Bezug läuft ausschließlich innerhalb
  des Containerbaus.
- **Digest real gemessen, heute** (§Kontext): `ghcr.io/astral-sh/uv:0.12.17`
  == `ghcr.io/astral-sh/uv:latest` zum Zeitpunkt dieser ADR — die **Wahl**
  und der tatsächliche Pin-Commit gehören dem umsetzenden Zug
  (`slice-sdk-python-pack-werkzeug`), analog wie
  [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) §Entscheidung
  Festlegung 3 es für die C#/Kotlin-Basis-Images vorgemacht hat
  (Kandidat + gemessener Digest in der ADR, Pin-Commit beim umsetzenden
  Zug).

### 4 — Was unverändert bleibt

- **Kein Gate.** `make sdk-pack-python` und der Publish-Workflow bleiben
  Werkzeuge, wie in [`ADR-0107`](0107-python-pypi-zweites-sdk-package.md)
  Festlegung 5/Bindung 3 entschieden — der Wechsel des Frontends ändert
  nichts an dieser Einstufung (weiterhin Netz-bindend: PyPI-Paketbezug für
  Test-Abhängigkeiten und der Publish-Schritt selbst).
- **Der Tag-Namensraum (`sdk-python-v*`) und der Workflow-Name
  (`sdk-python-release.yml`)** bleiben unverändert — nur der darin
  aufgerufene Befehl wechselt.
- **Kein Produktionscode, kein Eingriff in `internal/**`/`cmd/**`.**
- **`ADR-0106` und `sdks/csharp/**` bleiben unberührt.**
- **`docs/user/version.md` bleibt unberührt.**

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

### E — Paket-Werkzeug (Python) — Fortsetzung von `ADR-0107` Tabelle E

| Option | Pro | Contra |
|---|---|---|
| E1 — `build`(PyPA)+`setuptools`-Backend+`twine` (bisherige Wahl, jetzt abgelöst) | von packaging.python.org als Minimalweg dokumentiert; bereits einmal begründet gewählt ([`ADR-0107`](0107-python-pypi-zweites-sdk-package.md)); `twine` ist seit Jahren der etablierte PyPI-Upload-Weg | zwei getrennte Werkzeuge für Bauen und Publizieren statt eines; kein natives Trusted-Publishing (braucht einen zusätzlichen OIDC-Austausch-Schritt vor dem `twine upload`, den `ADR-0107` Trigger 4 ohnehin als spätere Stufe vorsah); explizite Nutzerentscheidung heute gegen den Verbleib |
| **E2 — `uv build`+`uv publish` (gewählt)** | ein Werkzeug für Build **und** Publish; nativ PEP-517-kompatibel und respektiert das bestehende `setuptools`-Backend unverändert (real geprüft, §Kontext); natives, dokumentiertes Trusted-Publishing ohne Zusatz-Plugin; offizielle GitHub-Actions-Integrationsanleitung vorhanden (Astral, live abgerufen); aktiver Entwicklungstakt (aktuelle Version `0.12.17`, real gemessen); explizite Nutzerentscheidung dafür | `uv` selbst ist 0.x-versioniert, keine explizite Stabilitätszusage für `uv publish` in der abgerufenen Doku gefunden — offener Punkt, eigener Re-Evaluierungs-Trigger; ein weiteres digest-gepinntes Image kommt zur Pin-Pflege hinzu; der reale Produktiv-Lauf bleibt bis zum ersten `sdk-python-release.yml`-Durchlauf unbewiesen (`AGENTS.md` §3.10-Analogie) |
| E3 — Mischform: `uv build` fürs Bauen, `twine` weiterhin fürs Publizieren | kleinerer Wechsel; `twine` bleibt der etablierte, seit Jahren geprüfte Upload-Weg | verschenkt genau den Teil des Werkzeug-Wechsels, für den natives Trusted-Publishing spricht; **zwei** Werkzeuge bleiben zu pflegen (uv-Digest-Pin **und** twine-Versions-Pin) statt eines — kein Netto-Vereinfachungsgewinn, der den Wechsel überhaupt rechtfertigt |
| E4 — All-in-one-Wechsel zu Poetry oder Hatch, jetzt da ohnehin eine Alternative geprüft wird | würde dieselbe Gelegenheit nutzen, den Paket-Werkzeug-Stack einmalig neu zu bewerten | dieselben Einwände wie in [`ADR-0107`](0107-python-pypi-zweites-sdk-package.md) §Verglichene Alternativen E2/E3 gelten unverändert: eigenes Lockfile-/Environment-Format, Ersatz des Projektformats/Backends — genau der Unterschied zu `uv`, das das bestehende Backend **nicht** anfasst; ein Formatwechsel wäre zusätzlich ein inhaltlicher Eingriff in die bereits real existierende, reviewte `pyproject.toml` (`slice-sdk-python-projektgeruest`) ohne neuen Grund gegenüber der ursprünglichen Abwägung |
| E5 — nichts tun: bei `build`+`twine` bleiben | kein Aufwand, keine neue ADR nötig gewesen | widerspricht der expliziten, heute im Chat getroffenen Nutzerentscheidung für den Wechsel; die real existierende Alternativen-Lücke (`uv` fehlte in `ADR-0107` Tabelle E) bliebe unadressiert |

**Fazit:** E2. E1 verliert gegen die explizite Nutzerentscheidung und den
fehlenden nativen Trusted-Publishing-Support; E3 verschenkt den
Haupt-Vorteil des Wechsels ohne Pflegeaufwand zu senken; E4 trifft
unverändert die in [`ADR-0107`](0107-python-pypi-zweites-sdk-package.md)
bereits gegen Poetry/Hatch vorgebrachten Einwände; E5 ignoriert die
Nutzerentscheidung.

## Konsequenzen

- **Positiv:** ein Werkzeug (`uv`) trägt Bauen und Publizieren statt zweier
  getrennter (`build`+`twine`); das bestehende, bereits real existierende
  `setuptools`-Backend bleibt vollständig unangetastet — kein Migrationsaufwand
  an `sdks/python/pgchangefeed/pyproject.toml`; natives Trusted-Publishing
  steht bereit, sobald [`ADR-0107`](0107-python-pypi-zweites-sdk-package.md)
  §Re-Evaluierungs-Trigger 4 eintritt, ohne ein zusätzliches Werkzeug dafür
  einzuführen.
- **Negativ:** `uv`s Reifegrad für `uv publish` in einem realen PyPI-Release
  ist bis zum ersten produktiven Lauf unbewiesen (0.x-Versionierung, keine
  gefundene explizite Stabilitätszusage) — benannter offener Punkt, eigener
  Trigger unten; ein zusätzliches digest-gepinntes Image
  (`ghcr.io/astral-sh/uv`) kommt zur Pin-Pflege dieses Repos hinzu, mit
  derselben Freshness-Pflicht wie die bestehenden `pin-stale-*`-Achsen
  (Träger-Nachzug ist Sache des Zuges, der den Pin real einführt).
- **Folgepflicht:**
  1. Die noch **offenen** Planungsträger
     `docs/plan/planning/open/slice-sdk-python-pack-werkzeug.md`,
     `docs/plan/planning/open/slice-sdk-python-publish-workflow.md` und
     `docs/plan/planning/welle-sdk-python-lh-fa-sst-009.md` verweisen
     textuell auf `build`+`twine` — real geprüft existiert dazu noch kein
     Code (§Kontext Ist-Stand). Sie sind vor ihrer Umsetzung auf diese ADR
     zu ziehen (`AGENTS.md` §3.13) — Sache des Planners, nicht dieses
     ausschließlich ADR-schreibenden Zuges.
  2. Sobald `sdks/python/Dockerfile` real um die Pack-Stufe erweitert wird,
     zieht der umsetzende Zug den in §Entscheidung Festlegung 3 benannten
     `ghcr.io/astral-sh/uv`-Digest-Pin ein (Pin-Commit ist ein bewusster
     Commit, `AGENTS.md` §3.6/Modul 14).
  3. `harness/README.md` §Werkzeuge bekommt seine `make sdk-pack-python`-/
     Publish-Workflow-Zeilen weiterhin erst, wenn sie real existieren
     (`AGENTS.md` §4) — unverändert aus
     [`ADR-0107`](0107-python-pypi-zweites-sdk-package.md).
  4. Der PyPI-Publish-Workflow gilt nach `AGENTS.md` §3.10 unverändert erst
     nach einem realen, grünen Post-Push-Lauf als abgeschlossen — jetzt mit
     `uv publish` statt `twine upload` als geprüftem Befehl.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Review (kein Sensor) | `sdks/python/pgchangefeed/pyproject.toml`s `[build-system]`-Block bleibt `setuptools.build_meta` — ein Wechsel des Backends wäre eine neue Entscheidung, keine Umsetzung dieser ADR | — |
| künftiges Bau-Werkzeug (Folgepflicht) | `make sdk-pack-python` ruft `uv build --no-sources` auf, der Publish-Workflow `uv publish`; beide Docker-only, kein Gate | `make sdk-pack-python` (noch nicht existent) |

## Re-Evaluierungs-Trigger

1. **`uv publish` erweist sich im ersten realen `sdk-python-release.yml`-Lauf
   (`AGENTS.md` §3.10) als instabil oder inkompatibel mit PyPIs
   Upload-API** — Rückkehr zu `twine` als eigene, weitere Folge-ADR mit
   `Supersedes ADR-0108`.
2. **PyPA/packaging.python.org erklärt `uv` (oder ein Nachfolgewerkzeug)
   offiziell zum empfohlenen Minimalweg für ein einzelnes Package** — keine
   weitere Prüfung nötig, diese Entscheidung bleibt bestätigt.
3. **`uv` erreicht Version 1.0 mit einer expliziten Stabilitätszusage für
   `uv publish`** — der in §Kontext benannte offene Punkt schließt sich;
   kein neuer Trigger nötig, aber als Ereignis in dieser ADRs Geschichte
   nachzutragen (Zitat-Korrektur-fähig nach
   [`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md), da reine
   Ergänzung eines Beleg-Datums, kein Überschreiben von §Entscheidung).
4. **[`ADR-0107`](0107-python-pypi-zweites-sdk-package.md) §Re-Evaluierungs-Trigger
   4 tritt ein** (Umstellung auf Trusted Publishing nach dem ersten
   erfolgreichen Token-Publish) — läuft jetzt über `uv publish`s native
   OIDC-Unterstützung, kein zusätzliches Werkzeug nötig.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-19 | Proposed | dieser Architect-Zug (kein vorausgehender Slice-Plan — ADR vor Umsetzung des Frontend-Wechsels, Modul 8) |
| 2026-09-19 | Accepted | Nutzerentscheidung im Chat vom 2026-09-19 nach Rückfrage per `AskUserQuestion` mit Empfehlung, samt Alternativenvergleich und real recherchiertem Reifegrad-Beleg |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0108` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
