# Slice sdk-python-projektgeruest: SDK-Projektgerüst `sdks/python/pgchangefeed/`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-python-lh-fa-sst-009](../welle-sdk-python-lh-fa-sst-009.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (Scope),
[`ADR-0107`](../../adr/0107-python-pypi-zweites-sdk-package.md)
Festlegung 1/3/4 (Umfang, Ort, Versionierung).

**Berührte Spec-Stellen:** [`LH-FA-SST-009.a`](../../../../spec/pflichtenheft.md)
(bleibt nach diesem Slice noch offen für Python — der Träger-Nachzug folgt
erst mit `slice-sdk-python-pack-werkzeug`, wenn das Package real
paketierbar ist).

**Verantwortlich:** Implementer-Agent (direkt beauftragt).

**Autor:** Planner-Agent, direkt beauftragt (`ADR-0107` §Konsequenzen
Folgepflicht 1). **Datum:** 2026-09-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein neues, eigenständiges Python-Projekt unter
`sdks/python/pgchangefeed/` anlegen — `pyproject.toml` mit PyPA-Standard-
Metadaten (`name = "pgchangefeed"`, `version = "0.1.0"`, `description`,
`license`, `urls`), ein eigenes, digest-gepinntes Docker-Bau-Setup (analog
`sdks/csharp/Dockerfile`, aber `python`-Basis statt `dotnet/sdk`), ein
englischsprachiges `README.md` und ein minimales Konfigurations-/
Optionen-Skelett (Adresse, Token) — ohne jeden Import aus
`internal/**`/`cmd/**` dieses Repos.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Öffentliche HTTP-API-Fläche** — `slice-sdk-python-http-client-flaeche`
  übernimmt das; dieser Slice liefert nur den Ort, in den sie schreibt (das
  Skelett bleibt minimal — höchstens eine gemeinsame Konfigurationsklasse
  für Adresse/Token, kein Vorgriff auf Endpunkt-Methoden).
- **`make sdk-pack-python` (Pack-Werkzeug)** —
  `slice-sdk-python-pack-werkzeug` übernimmt das; dieser Slice liefert nur
  ein Docker-Bau-Setup, das `pip install`/`pytest` trägt, kein
  `python -m build`-Aufruf und kein `make`-Ziel.
- **PyPI-Publish-Workflow** — `slice-sdk-python-publish-workflow`
  übernimmt das; dieser Slice pflegt nur die `version` in
  `pyproject.toml`, kein Tag-Trigger, kein Secret-Bezug.
- **gRPC/SSE/NATS-Anteile** — `ADR-0107` Festlegung 1 grenzt v1
  ausdrücklich auf HTTP-API ein (Welle-Plan §6 Out-of-Scope); dieser
  Slice legt kein Codegenerierungs-Gerüst dafür an.
- **Ein `examples/python/`-Referenz-Client** — nicht Teil dieser ADR/Welle
  (`ADR-0107` §Kontext Ist-Stand stellt sein Fehlen ausdrücklich fest, ohne
  seine Nachlieferung zu verlangen).
- **Träger-Nachzug in `spec/pflichtenheft.md`** — `ADR-0107` §Konsequenzen
  Folgepflicht 2 bindet den Nachzug an den Zug, der das Package **real**
  existent macht (`python -m build` läuft real); das ist
  `slice-sdk-python-pack-werkzeug`, nicht dieser Slice, der nur ein
  Projekt-Gerüst ohne Pack-Fähigkeit liefert.

## 2. Definition of Done

- [x] `sdks/python/pgchangefeed/pyproject.toml` existiert (PyPA-Standard,
      `[project]`-Tabelle): `name = "pgchangefeed"`,
      `version = "0.1.0"` (`ADR-0107` Festlegung 4, Start bei `0.x.y`,
      PEP 440), `description`, `authors`, `license` (`MIT`, wie das
      Repo-Root-`LICENSE`), `[project.urls]` auf dieses Repo verweisend,
      `readme` auf das SDK-eigene `README.md`. Kein Import eines privaten
      Baums dieses Repos, keine Fremdabhängigkeit über die
      HTTP-Client-Bibliothek hinaus (`ADR-0107` §Entscheidung Festlegung 1:
      `httpx`, keine schwere Fremdabhängigkeit).
- [x] `sdks/python/Dockerfile` (Bau-Kontext `sdks/python/`, eigenständig
      von `sdks/csharp/Dockerfile` und der Wurzel-`Dockerfile`) mit
      digest-gepinnter `python`-Basis (real gemessener Digest zum
      Bau-Zeitpunkt, Kommentar-Pflicht analog `sdks/csharp/Dockerfile`);
      `pip install -e .`/`pytest` laufen darin, kein Runtime-Server-Start
      nötig (ein SDK ist keine startbare Anwendung).
- [x] `sdks/python/README.md` (Englisch) beschreibt Zweck,
      Installationsweg (`pip install pgchangefeed`) und verweist auf das
      Repo-Root-`README.md` für den vollen Kontext — über absolute
      GitHub-Blob-URLs statt relativer Pfade, weil ein PyPI-Package
      außerhalb dieses Repository-Checkouts gelesen wird (Analogie
      `sdks/csharp/README.md`, siehe dessen Closure-Notiz); kein Duplikat
      der Draht-Doku (`SPEC-018` bleibt die kanonische Quelle).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Erst-Lauf fand 1 HIGH (F-1, Slice-Chronik in `__init__.py`-Docstring),
      Fixrunde behoben (Commit `70aa64ce`), Fixrunden-Nachprüfung im selben
      Report freigegeben (`docs/reviews/review-slice-sdk-python-projektgeruest.md`
      §Fixrunden-Nachprüfung, 2026-09-19) — 0 offenes HIGH.
- [x] Doku-Update für `harness/README.md` entfällt in diesem Slice — kein
      neues `make`-Target entsteht hier (Pack-Werkzeug folgt in
      `slice-sdk-python-pack-werkzeug`).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo (kein Brownfield-Bootstrap).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis `BEO-PGC/<slug>/` oder eine weitere Datei in dessen
      `evidence/`; **kein Zähler wird gesetzt**, er folgt aus den Dateien.
      Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in
      §7 notiert. **Vor dem ersten Test-/Bau-Lauf dieser neuen Sub-Area:**
      Suchlauf gegen `sdks/csharp/**` auf ein bereits gelöstes, analoges
      Problem prüfen, bevor ein neues Workaround-Muster erfunden wird
      (`BEO-PGC/workaround-uebersieht-etabliertes-muster-im-bestand`,
      2×, siehe Welle-Plan §6 Eröffnungs-Sichtung und §6 dieses Plans).
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-python-lh-fa-sst-009](../welle-sdk-python-lh-fa-sst-009.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/python/pgchangefeed/pyproject.toml` | neu | PyPA-Metadaten, `version = "0.1.0"` (`ADR-0107` Festlegung 3/4). |
| `sdks/python/pgchangefeed/src/pgchangefeed/__init__.py` (Arbeitsname) | neu | Package-Wurzel, minimales Skelett. |
| `sdks/python/pgchangefeed/src/pgchangefeed/options.py` (Arbeitsname) | neu | gemeinsamer Konfigurations-Nenner (Adresse, Token) für den Folge-Slice, falls sich beim Schreiben ein echter zeigt — kein Vorgriff auf Endpunkt-Methoden. |
| `sdks/python/Dockerfile` | neu | digest-gepinnter Bau, `pip install -e .`/`pytest`, keine Runtime-Stufe. |
| `sdks/python/README.md` | neu | Englisch, Installationsweg, Verweis auf Repo-Root-`README.md` über absolute GitHub-Blob-URLs. |
| `sdks/python/.gitignore` | neu | `__pycache__/`, `*.egg-info/`, `dist/`, `build/`, analog `sdks/csharp/.gitignore`. |
| Testdatei (`tests/test_options.py`, Arbeitsname) | neu | Konstruktions-/Validierungstest der Optionsklasse (kein Draht-Verhalten, das kommt mit dem Folge-Slice), `pytest`. |

**Plan-Nachzug (Nutzer-Entscheidung nach Verifikation):** `requires-python`
auf `>=3.14` (statt der ursprünglich gewählten `>=3.11`) und die
Docker-Basis auf `python:3.14-slim` (`3.14.7-slim-trixie`, real per
`docker buildx imagetools inspect` gemessener Index-Digest
`sha256:caaf356f40667c496d405780745b9ac25771c189a51dfcc42430d531ea09f8a2`,
Python 3.14.7 real als aktuelle stabile Version verifiziert, released
2026-08-05) angehoben — bewusst die aktuelle Version statt einer breiten
Mindestversion, analog der `net10.0`-Festlegung des C#-SDK. Real
gegengebaut: `docker build --no-cache -f sdks/python/Dockerfile
sdks/python` → Exit 0, `pytest` 3/3 grün unter Python 3.14.7.

**Ansatz:** Referenzmaterial für die Bau-/Projektstruktur ist
`sdks/csharp/PgChangeFeed.Client/`/`sdks/csharp/Dockerfile` (Formvorbild,
kein Code-Import — andere Sprache, anderes Paket-Ökosystem). Der
Implementer prüft **vor** dem ersten Bau-Lauf, ob `sdks/csharp/**` bereits
ein analoges Docker-Bau-/Test-Problem gelöst hat (siehe DoD-Punkt oben,
`BEO-PGC/workaround-uebersieht-etabliertes-muster-im-bestand`).

**Plan-Nachzug (Implementer, vor dem Gate-Lauf):**

- **Vorbild-Suchlauf (`sdks/csharp/**`) durchgeführt** —
  `sdks/csharp/Dockerfile` und `sdks/csharp/.gitignore` vollständig
  gelesen. Übernommen: Layer-Reihenfolge (Metadaten-/Projektdatei zuerst,
  dann Quellcode, damit der Docker-Layer-Cache über Quelltext-Änderungen
  stabil bleibt), keine Runtime-Stufe (ein SDK ist keine startbare
  Anwendung), Digest-Kommentar-Form (real gemessener Digest zum
  Bau-Zeitpunkt samt Pin-Hebungs-Hinweis), `.gitignore`-Kommentarform
  (Docker-only-Hinweis + Grund für den `dist/`-Eintrag). Nicht übertragbar:
  das `.proto`-Zusatzkontext-COPY-Muster (kein gRPC in diesem Erst-Scope,
  `ADR-0107` Festlegung 1) — aber das **Grundmuster** „eine vom
  Projekt an einer bestimmten relativen Stelle erwartete Datei per
  `COPY` dorthin bringen, ohne sie im committeten Baum zu duplizieren"
  wurde real für das README-Problem unten wiederverwendet (siehe
  nächster Punkt) statt ein neues Workaround zu erfinden.
- **Abweichung von der ursprünglich angenommenen Form (real geprüft, kein
  Vorgriff):** `[project] readme` in `pyproject.toml` kann — anders als
  `PgChangeFeed.Client.csproj`s `<None Include="../README.md" .../>` bei
  C# (MSBuild erlaubt beliebige relative Pfade) — **nicht** auf eine Datei
  außerhalb des Projekt-Wurzelverzeichnisses zeigen: `pip install -e
  ".[test]"` scheiterte real mit
  `distutils.errors.DistutilsOptionError: Cannot access
  '/src/pgchangefeed/../README.md' (or anything outside
  '/src/pgchangefeed')` (setuptools `_assert_local`). Lösung: `readme =
  "README.md"` (rein lokaler Pfad) in `pyproject.toml`, und
  `sdks/python/Dockerfile` kopiert die committete
  `sdks/python/README.md` zusätzlich an die von `pyproject.toml`
  erwartete lokale Stelle (`pgchangefeed/README.md`, neben der
  `pyproject.toml`) — dasselbe wiederverwendete Muster wie oben, keine
  Duplizierung im committeten Baum, nur im Docker-Bau-Layer. Kein Wechsel
  der Baumstruktur (`sdks/python/pgchangefeed/pyproject.toml` und
  `sdks/python/README.md` bleiben an den in §2 genannten Orten).
- **Python-Mindestversion: `requires-python = ">=3.11"`.** Real recherchiert
  (`docker buildx imagetools inspect python:3.13-slim`, 2026-09-19): Python
  3.9 ist bereits EOL (Oktober 2025), 3.10 erreicht sein EOL im Oktober
  2026 (also binnen eines guten Monats nach diesem Zug) — beide damit keine
  sinnvolle Untergrenze für ein neu veröffentlichtes Package. 3.11 bleibt
  bis Oktober 2027 im Security-Support, ist breit verfügbar (jede aktuelle
  Linux-Distribution/jedes offizielle `python`-Docker-Image führt es) und
  bringt bereits verbesserte Fehlermeldungen/Exception-Groups — eine
  Analogie zur `net10.0`-Wahl bei C# (`ADR-0106`): die jeweils aktuell
  unterstützte, nicht die älteste technisch mögliche Version, ohne
  Nutzungsdaten, die eine andere Wahl rechtfertigen (§6 dieses Plans).
- **Basis-Image-Wahl:** `python:3.13-slim` (statt der vollen `python:3.13`-
  Variante) — kleineres Image, keine C-Compiler-Abhängigkeit nötig (`httpx`
  und seine Abhängigkeiten `httpcore`/`h11`/`certifi`/`idna`/`anyio`/
  `sniffio` sind reine Python-Pakete, kein Bau von C-Extensions). Digest
  real gemessen als **Manifest-Index-Digest** (alle Plattformen), analog
  zur C#-Basis, die ebenfalls den vollständigen, plattformübergreifenden
  Digest trägt.

## 4. Trigger

**Start** (`next` → `in-progress`): sofort — `ADR-0107` ist `Accepted`,
keine weitere Vorbedingung.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): entfällt aus
  heutiger Sicht — Umfang ist eine `pyproject.toml`, ein Dockerfile, ein
  README, ein minimales Optionsmodul.
- `in-progress` → `open` (blockiert — Carveout?): kein gepinntes
  `python`-Basisimage mit dem für das Skelett benötigten Python-Standard
  auffindbar — unwahrscheinlich, offizielle `python`-Images sind für jede
  aktuelle Minor-Version verfügbar.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- Das minimale Konfigurations-Skelett könnte beim Schreiben des
  Folge-Slices (`slice-sdk-python-http-client-flaeche`) doch keinen echten
  gemeinsamen Nenner zeigen (anders als bei C#, wo zwei Flächen — HTTP und
  gRPC — ihn teilen mussten; hier gibt es nur eine Fläche, der Nenner
  könnte deshalb trivial oder überflüssig sein). **Ausgang:** weiter
  offen, entschieden beim Schreiben des Folge-Slices.
- **Kein Python-Referenz-Client existiert** (`ADR-0107` §Kontext „Was das
  ändert") — für dieses reine Projektgerüst-Slice ohne Draht-Verhalten hat
  das keine unmittelbare Wirkung, ist aber die Grundbedingung, unter der
  der nachfolgende Fläche-Slice arbeitet; hier nur zur Vollständigkeit
  benannt, die eigentliche Risiko-Tragweite steht in
  `slice-sdk-python-http-client-flaeche` §6. **Ausgang:** weiter offen,
  strukturell (kein Ausgang möglich, bevor die Fläche geschrieben ist).
- Die Wahl der Python-Mindestversion (z. B. `>=3.10`) für
  `[project.requires-python]` schränkt den Konsumentenkreis ein — ein
  SDK-Package hat potenziell breitere Zielgruppen als ein
  Docker-only-Beispielprogramm. **Ausgang:** weiter offen — Entscheidung
  bleibt bei diesem Slice (Analogie zur C#-`net10.0`-Entscheidung, keine
  Nutzungsdaten, die eine andere Wahl rechtfertigen); ein Wechsel wäre eine
  spätere, eigenständige Entscheidung (`ADR-0107` §Re-Evaluierungs-Trigger,
  Nutzungsdaten-getrieben).

## 7. Closure-Notiz

<…>

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Neue Sub-Area `sdks/python/` —
noch nie berührt, entsteht mit diesem Slice erstmalig. Analogie zur
bereits etablierten Sub-Area `sdks/csharp/` (`ADR-0106`, Docker-only,
digest-/paket-gepinnt) — dieselbe Konventionen-Dichte übertragen, keine
eigene Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/workaround-uebersieht-etabliertes-muster-im-bestand` (offen, 2×,
siehe Welle-Plan §6 Eröffnungs-Sichtung) betrifft diesen Slice unmittelbar
(erster Docker-Bau eines neuen Sprach-Ökosystems in diesem Repo, nach der
C#-Sub-Area die zweite ihrer Art) — als DoD-Punkt und Plan-Hinweis oben
aufgenommen. Kein weiterer Treffer für „Docker-only-Bau eines neuen
Sprach-Ökosystems" oder „PyPA-Metadaten" über die im Welle-Plan §6
dokumentierte `grep`-Suche hinaus.

**Modus-Begründungsblock:**

### Sub-Area: `sdks/python/`

- **Modus:** Greenfield — neuer Baum, Doku (`ADR-0107`) führt vor Code.
- **Konventionen-Dichte:** übernommen von `sdks/csharp/`
  (`harness/conventions.md` §Modus-Deklaration, Default-Sub-Area `*`/`PGC`
  Greenfield) — Docker-only, digest-/paket-gepinnt, kein Host-`python`/`pip`.
- **Phase-Reife:** Phase 0 (Erstanlage) — Doc (`ADR-0107`) vollständig vor
  dem ersten Code dieser Sub-Area.
- **Evidenz-/Diskrepanz-Risiko:** niedrig (GF, keine Inventur nötig).
- **Reconciliation-Aufwand:** keiner — kein Brownfield-Bestand in diesem
  Baum.
