# Slice sdk-python-pack-werkzeug: `make sdk-pack-python` + Träger-Nachzug Pflichtenheft

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-python-lh-fa-sst-009](../welle-sdk-python-lh-fa-sst-009.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md),
[`ADR-0107`](../../adr/0107-python-pypi-zweites-sdk-package.md)
Festlegung 5 (Build-Mechanismus: Docker-only bauen, `build`+`twine`),
§Konsequenzen Folgepflicht 2/3 (Träger-Nachzug Pflichtenheft/
`harness/README.md`).

**Berührte Spec-Stellen:** [`LH-FA-SST-009.a`](../../../../spec/pflichtenheft.md)
(wird mit diesem Slice für Python/PyPI aufgelöst — die Kennung selbst
bleibt bestehen, ihr Satz „zweite Sprache/Vertriebsweg offen" wird durch
`ADR-0107` und die jetzt reale Paketierbarkeit falsch, `AGENTS.md` §3.13),
`spec/pflichtenheft.md` §6 Externe Verträge (neue Zeile für das Package).

**Verantwortlich:** — bis zur Priorisierung.

**Autor:** Planner-Agent, direkt beauftragt (`ADR-0107` §Konsequenzen
Folgepflicht 1/2/3). **Datum:** 2026-09-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein neues, netzlos **nicht** prüfbares Werkzeug-Ziel
`make sdk-pack-python` (Docker-only, `pip install`/`pytest`/
`python -m build`, gepinntes `python`-Image) erzeugt reale `.whl`- und
`.tar.gz`-Artefakte aus `sdks/python/pgchangefeed/` — dazu der
Träger-Nachzug, den `ADR-0107` §Konsequenzen Folgepflicht 2/3 fordert:
[`LH-FA-SST-009.a`](../../../../spec/pflichtenheft.md) (§1) und
`spec/pflichtenheft.md` §6 (Externe Verträge, neue `SPEC-<NNN>`-Zeile für
das Package — **nächste freie Nummer real verifizieren**, siehe §3) sowie
`harness/README.md` §Werkzeuge (die reale `make sdk-pack-python`-Zeile).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **PyPI-Veröffentlichung** — `slice-sdk-python-publish-workflow`
  übernimmt das; dieser Slice erzeugt die Artefakte, veröffentlicht sie
  nicht (`ADR-0107` Festlegung 5: „Bauen und Paketieren bleiben
  Docker-only … Veröffentlichung … braucht einen eigenen,
  Netz-bindenden Workflow").
- **Neues Gate** — `make sdk-pack-python` bleibt außerhalb von
  `GATE_CHECKS`/`make gates` (`ADR-0107` Festlegung 5, dieselbe
  Begründung wie `make sdk-pack-csharp`/`make examples-csharp`:
  Paket-/Test-Abhängigkeits-Bezug braucht Netz).
- **Änderung der HTTP-Fläche** — dieser Slice paketiert, was
  `slice-sdk-python-http-client-flaeche` bereits geliefert hat; kein neuer
  API-Umfang.
- **`docs/user/version.md`-Änderung** — die Server-Version bleibt
  eigenständig geführt (`ADR-0107` Festlegung 4/6).

## 2. Definition of Done

- [ ] `make sdk-pack-python` existiert (Docker-only, kein Gate — analog
      `make sdk-pack-csharp`): baut, testet (`pytest`) und paketiert
      (`python -m build`) `sdks/python/pgchangefeed/` im gepinnten
      `python`-Image; ein roter Test bricht den `docker build` mit Exit ≠ 0
      ab. Die Artefakte (`.whl`, `.tar.gz`) werden über einen `tar`-Stream-
      Export aus dem Docker-Bau in ein lokales Verzeichnis (z. B.
      `sdks/python/dist/`, `.gitignore`t) abgelegt — analog dem
      Extraktionsmuster von `make sdk-pack-csharp`
      (`tools/harness/sdk-pack-csharp.sh`, host-seitiger Export statt
      Bind-Mount/`--user`-Workaround).
- [ ] Real ausgeführt: ein `.whl` und ein `.tar.gz` mit der erwarteten
      Version (z. B. `pgchangefeed-0.1.0-py3-none-any.whl`,
      `pgchangefeed-0.1.0.tar.gz`) existieren nach dem Lauf und sind als
      Smoke-Beleg im Bericht dieses Slice genannt (Datei-Existenz, keine
      Behauptung — `AGENTS.md` §3.12 Instanz B).
- [ ] `spec/pflichtenheft.md` §1 trägt bei
      [`LH-FA-SST-009.a`](../../../../spec/pflichtenheft.md) einen
      Nachzug-Satz: für Python/PyPI ist die Sprachmatrix-/Vertriebsweg-
      Frage durch [`ADR-0107`](../../adr/0107-python-pypi-zweites-sdk-package.md)
      beantwortet und das Package real paketierbar; eine dritte Sprache
      oder ein dritter Vertriebsweg bleibt offen (`AGENTS.md` §3.13, kein
      Streichen der Kennung). **Kein ADR-Bezug in der Nachzug-Zeile selbst**
      — `make docs-check`s `matrix`-Modul verbietet mechanisch jede
      Referenz `spec → adr` (real erprobtes Muster aus
      `slice-sdk-csharp-pack-werkzeug`, siehe dessen Closure-Notiz „Was
      ging anders als geplant"); die Beziehung steht bereits umgekehrt in
      `ADR-0107`s eigenem `Schärft:`-Feld.
- [ ] `spec/pflichtenheft.md` §6 Externe Verträge bekommt eine neue
      `SPEC-<NNN>`-Zeile für `pgchangefeed` (System: PyPI-Package;
      Version: PEP 440, `0.x.y`; Vertrag-Datei: Verweis auf
      `sdks/python/pgchangefeed/pyproject.toml` als Metadaten-Quelle,
      analog der bestehenden `SPEC-026`-Zeile für das C#-Package) sowie §7
      Historie-Zeile — **ohne** ADR-Bezug (dieselbe `matrix`-Regel wie
      oben). Die nächste freie `SPEC-<NNN>`-Nummer ist **real zu
      verifizieren**, nicht aus der ADR zu übernehmen: `grep -oE
      "SPEC-[0-9]+" spec/pflichtenheft.md | sort -t- -k2 -n -u | tail -3`
      — zum Planungszeitpunkt dieser Welle ist `SPEC-026` die höchste
      vergebene Nummer (das C#-Package selbst), aber zwischen Planung und
      Umsetzung dieses Slice könnten weitere Nummern vergeben worden sein.
- [ ] `harness/README.md` §Werkzeuge bekommt die reale
      `make sdk-pack-python`-Zeile (kein Gate, Bindung auf
      [`ADR-0107`](../../adr/0107-python-pypi-zweites-sdk-package.md)) —
      erst jetzt zulässig, weil das Ziel jetzt real existiert
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
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-python-lh-fa-sst-009](../welle-sdk-python-lh-fa-sst-009.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/mk/sdk.mk` (bereits durch `slice-sdk-csharp-pack-werkzeug` angelegt) | update | `sdk-pack-python`-Target ergänzen, Docker-only, kein `GATE_CHECKS`-Eintrag. |
| `sdks/python/Dockerfile` | update | zusätzliche Bau-/Export-Stufe für `python -m build`, Extraktion analog `make sdk-pack-csharp`. |
| `tools/harness/sdk-pack-python.sh` (Arbeitsname) | neu | Host-seitiger `tar`-Stream-Export, analog `tools/harness/sdk-pack-csharp.sh`. |
| `spec/pflichtenheft.md` §1 (`LH-FA-SST-009.a`) | update | Nachzug-Satz: Python/PyPI nicht mehr offen, dritte Sprache/Vertriebsweg bleibt offen. |
| `spec/pflichtenheft.md` §6 | update | neue `SPEC-<NNN>`-Zeile für `pgchangefeed` (Nummer real verifizieren, siehe §2). |
| `spec/pflichtenheft.md` §7 Historie | update | Historie-Zeile für die neue `SPEC-<NNN>` und den `LH-FA-SST-009.a`-Nachzug. |
| `harness/README.md` §Werkzeuge | update | reale `make sdk-pack-python`-Zeile. |

**Ansatz:** Strukturvorbild ist `slice-sdk-csharp-pack-werkzeug` (bereits
`done/`) — dasselbe Export-Muster (`docker run --rm --network none <image>
| tar -x -C .`, `AGENTS.md` §3.9 Pipe-Disziplin), dieselbe
`matrix`-Modul-Umgehung für die Spec-Nachzugzeilen (kein ADR-Bezug),
dasselbe Vorgehen zur Kollisionsfreiheit der neuen `SPEC-<NNN>`-Nummer
(`grep -c "SPEC-<NNN>" spec/*.md` nach dem Schreiben, genau ein Treffer in
`spec/pflichtenheft.md`).

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-sdk-python-http-client-flaeche`
in `done/` liegt (siehe Welle-Plan §4 Reihenfolge).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): entfällt aus
  heutiger Sicht — ein Make-Target plus zwei Doku-Nachzüge.
- `in-progress` → `open` (blockiert — Carveout?): `python -m build`
  scheitert strukturell am gepinnten Image (z. B. fehlendes
  `build`-Paket) — unwahrscheinlich, `build` ist Standard-PyPA-Tooling,
  über `pip install build` netzgebunden nachinstallierbar im Bau-Schritt.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + reale `.whl`/`.tar.gz`-Artefakte als
Smoke-Beleg + Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Der Export-Mechanismus für `.whl`/`.tar.gz` aus dem Docker-Bau ist für
  diesen Artefakt-Typ neu (das C#-Vorbild exportierte ein einzelnes
  `.nupkg`, hier sind es zwei Dateien) — ein Bind-Mount-Workaround oder
  ein falsch gesetztes `--user` könnte die Artefakte mit falschen
  Dateirechten oder gar nicht auf den Host bringen. **Ausgang:** weiter
  offen, zu entscheiden beim Schreiben — Erwartung: derselbe
  `tar`-Stream-Export wie bei `make sdk-pack-csharp` trägt auch für zwei
  Dateien statt einer, da `tar` mehrere Dateien in einem Stream
  transportiert; real zu bestätigen.
- Die neue `SPEC-<NNN>`-Nummer in `spec/pflichtenheft.md` §6/§2 muss
  fortlaufend vergeben werden — zum Planungszeitpunkt dieser Welle ist
  `SPEC-026` die zuletzt vergebene Nummer (real verifiziert,
  `grep -oE "SPEC-[0-9]+" spec/pflichtenheft.md | sort -t- -k2 -n -u |
  tail -3` → `SPEC-024`, `SPEC-025`, `SPEC-026`); zwischen Planung und
  Umsetzung dieses Slice könnte eine weitere Nummer vergeben worden sein
  (Kollisionsrisiko). **Ausgang:** weiter offen — der Implementer
  verifiziert die nächste freie Nummer **erneut** unmittelbar vor dem
  Schreiben, nicht nur einmal bei der Planung (§2 DoD-Punkt oben).
- Python-Paketnamen unterliegen PyPI-Namensraum-Kollisionen — der
  Arbeitsname `pgchangefeed` könnte zum Zeitpunkt der Umsetzung bereits
  von einem fremden Package belegt sein (anders als bei NuGet, wo
  `PgChangeFeed.Client` bereits real geprüft frei war). **Ausgang:**
  weiter offen — der Implementer prüft `pip index versions pgchangefeed`
  bzw. die PyPI-Weboberfläche vor dem ersten realen Publish-Versuch
  (`slice-sdk-python-publish-workflow`); für dieses Pack-Werkzeug-Slice
  selbst (netzlos, kein Publish) hat eine Namensraum-Kollision keine
  unmittelbare Wirkung.

## 7. Closure-Notiz

<…>

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Zwei Sub-Areas: `sdks/python/`
(bereits eröffnet, GF) und `spec/pflichtenheft.md`/`harness/README.md`
(Default-Sub-Area `*`/`PGC`, Greenfield, `harness/conventions.md`).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
kein Treffer, der speziell dieses Pack-Werkzeug-Slice beträfe, über die
bereits in der Welle-Plan §6 benannten Einträge hinaus.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
