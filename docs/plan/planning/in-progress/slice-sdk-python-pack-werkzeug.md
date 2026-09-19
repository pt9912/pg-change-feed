# Slice sdk-python-pack-werkzeug: `make sdk-pack-python` + Träger-Nachzug Pflichtenheft

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-python-lh-fa-sst-009](../welle-sdk-python-lh-fa-sst-009.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md),
[`ADR-0107`](../../adr/0107-python-pypi-zweites-sdk-package.md)
Festlegung 5 (Build-Mechanismus: Docker-only bauen),
§Konsequenzen Folgepflicht 2/3 (Träger-Nachzug Pflichtenheft/
`harness/README.md`),
[`ADR-0108`](../../adr/0108-python-sdk-uv-statt-build-twine.md)
§Entscheidung Festlegung 1/3 (Build-Frontend `uv build --no-sources`
statt `build`(PyPA), digest-gepinntes `ghcr.io/astral-sh/uv`-Kopier-Image)
— superseded die `build`+`twine`-Klausel aus `ADR-0107` Festlegung 5.

**Berührte Spec-Stellen:** [`LH-FA-SST-009.a`](../../../../spec/pflichtenheft.md)
(wird mit diesem Slice für Python/PyPI aufgelöst — die Kennung selbst
bleibt bestehen, ihr Satz „zweite Sprache/Vertriebsweg offen" wird durch
`ADR-0107` und die jetzt reale Paketierbarkeit falsch, `AGENTS.md` §3.13),
`spec/pflichtenheft.md` §6 Externe Verträge (neue Zeile für das Package).

**Verantwortlich:** Implementer-Agent (dietmar.burkard@nerdware.dev), ab
2026-09-19.

**Autor:** Planner-Agent, direkt beauftragt (`ADR-0107` §Konsequenzen
Folgepflicht 1/2/3). **Datum:** 2026-09-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein neues, netzlos **nicht** prüfbares Werkzeug-Ziel
`make sdk-pack-python` (Docker-only, `pip install`/`pytest`/
`uv build --no-sources`, gepinntes `python`-Image + digest-gepinntes
`ghcr.io/astral-sh/uv`-Kopier-Image,
[`ADR-0108`](../../adr/0108-python-sdk-uv-statt-build-twine.md)
§Entscheidung Festlegung 1/3) erzeugt reale `.whl`- und
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

- [x] `make sdk-pack-python` existiert (Docker-only, kein Gate — analog
      `make sdk-pack-csharp`): baut, testet (`pytest`) und paketiert
      (`uv build --no-sources`,
      [`ADR-0108`](../../adr/0108-python-sdk-uv-statt-build-twine.md)
      §Entscheidung Festlegung 1) `sdks/python/pgchangefeed/` im gepinnten
      `python`-Image; `uv` selbst wird **nicht** per `pip install`
      bezogen, sondern per digest-gepinntem Multi-Stage-`COPY` aus
      Astrals eigenem Werkzeug-Image
      (`COPY --from=ghcr.io/astral-sh/uv:0.12.17@sha256:10787c682e4184e4f290de1171fd4703dc63de99221f10fe1c99002ce7fa9acc
      /uv /uvx /bin/`, exakter Wortlaut aus `ADR-0108` §Entscheidung
      Festlegung 3 — Digest beim Schreiben gegen die ADR erneut
      verifizieren, nicht blind übernehmen); ein roter Test bricht den
      `docker build` mit Exit ≠ 0 ab. Die Artefakte (`.whl`, `.tar.gz`)
      werden über einen `tar`-Stream-Export aus dem Docker-Bau in ein
      lokales Verzeichnis (z. B. `sdks/python/dist/`, `.gitignore`t)
      abgelegt — analog dem Extraktionsmuster von `make sdk-pack-csharp`
      (`tools/harness/sdk-pack-csharp.sh`, host-seitiger Export statt
      Bind-Mount/`--user`-Workaround).
- [x] Real ausgeführt: ein `.whl` und ein `.tar.gz` mit der erwarteten
      Version (z. B. `pgchangefeed-0.1.0-py3-none-any.whl`,
      `pgchangefeed-0.1.0.tar.gz`) existieren nach dem Lauf und sind als
      Smoke-Beleg im Bericht dieses Slice genannt (Datei-Existenz, keine
      Behauptung — `AGENTS.md` §3.12 Instanz B).
- [x] `spec/pflichtenheft.md` §1 trägt bei
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
- [x] `spec/pflichtenheft.md` §6 Externe Verträge bekommt eine neue
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
- [x] `harness/README.md` §Werkzeuge bekommt die reale
      `make sdk-pack-python`-Zeile (kein Gate, Bindung auf
      [`ADR-0107`](../../adr/0107-python-pypi-zweites-sdk-package.md)) —
      erst jetzt zulässig, weil das Ziel jetzt real existiert
      (`AGENTS.md` §4).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Report: `docs/reviews/review-slice-sdk-python-pack-werkzeug.md`
      (0 HIGH, 0 MEDIUM, 1 LOW, 1 INFO — keine Fixrunde nötig,
      DoD-Checkbox-Nachzug ohne Fixrunde gemäß Reviewer-Skill).
- [x] Doku-Update: `harness/README.md` §Werkzeuge (siehe oben) — entfällt
      als eigener Punkt, da bereits oben als DoD-Kriterium geführt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`; keine Beobachtung angefallen
      ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-python-lh-fa-sst-009](../welle-sdk-python-lh-fa-sst-009.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `harness/mk/sdk.mk` (bereits durch `slice-sdk-csharp-pack-werkzeug` angelegt) | update | `sdk-pack-python`-Target ergänzen, Docker-only, kein `GATE_CHECKS`-Eintrag. |
| `sdks/python/Dockerfile` | update | zusätzliche Bau-/Export-Stufe für `uv build --no-sources` (digest-gepinntes `ghcr.io/astral-sh/uv`-Kopier-Image, `ADR-0108` §Entscheidung Festlegung 3), Extraktion analog `make sdk-pack-csharp`. |
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
- `in-progress` → `open` (blockiert — Carveout?): `uv build --no-sources`
  scheitert strukturell am digest-gepinnten `ghcr.io/astral-sh/uv`-Bezug
  (z. B. Digest inzwischen zurückgezogen) — unwahrscheinlich, der Digest
  wurde real gemessen (`ADR-0108` §Kontext); ein alternativer, aktuell
  gültiger Digest wäre im Umsetzungszug real nachzumessen.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + reale `.whl`/`.tar.gz`-Artefakte als
Smoke-Beleg + Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Der Export-Mechanismus für `.whl`/`.tar.gz` aus dem Docker-Bau ist für
  diesen Artefakt-Typ neu (das C#-Vorbild exportierte ein einzelnes
  `.nupkg`, hier sind es zwei Dateien) — ein Bind-Mount-Workaround oder
  ein falsch gesetztes `--user` könnte die Artefakte mit falschen
  Dateirechten oder gar nicht auf den Host bringen. **Ausgang:**
  entfallen/aufgelöst — der `tar`-Stream-Export
  (`docker run --rm --network none <image> | tar -x -C sdks/python/dist/`)
  trug beim ersten Versuch real für beide Dateien in einem Stream
  (`pgchangefeed-0.1.0-py3-none-any.whl`, 9330 Bytes;
  `pgchangefeed-0.1.0.tar.gz`, 11290 Bytes) — kein Bind-Mount, kein
  `--user`-Workaround nötig, keine Dateirechte-Auffälligkeit
  (`git status --porcelain sdks/python/dist/` leer, Verzeichnis vollständig
  `.gitignore`t).
- Die neue `SPEC-<NNN>`-Nummer in `spec/pflichtenheft.md` §6/§2 muss
  fortlaufend vergeben werden — zum Planungszeitpunkt dieser Welle ist
  `SPEC-026` die zuletzt vergebene Nummer (real verifiziert,
  `grep -oE "SPEC-[0-9]+" spec/pflichtenheft.md | sort -t- -k2 -n -u |
  tail -3` → `SPEC-024`, `SPEC-025`, `SPEC-026`); zwischen Planung und
  Umsetzung dieses Slice könnte eine weitere Nummer vergeben worden sein
  (Kollisionsrisiko). **Ausgang:** entfallen — der Implementer verifizierte
  unmittelbar vor dem Schreiben erneut (`grep -oE "SPEC-[0-9]+"
  spec/pflichtenheft.md | sort -t- -k2 -n -u | tail -5` → unverändert
  `SPEC-026` als höchste Nummer), vergab `SPEC-027` und bestätigte danach
  Kollisionsfreiheit (`grep -c "SPEC-027" spec/*.md` → genau drei Treffer,
  ausschließlich in `spec/pflichtenheft.md`: §1-Fließtext, §6-Zeile,
  §7-Historie; `0` in `spec/lastenheft.md`/`spec/architecture.md`).
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

- **Was hat funktioniert:** Ein glatter Durchlauf ohne Plan-Nachzug — anders
  als bei `slice-sdk-csharp-pack-werkzeug` (dort musste die ADR-Referenz in
  der `LH-FA-SST-009.a`-Nachzugzeile nachträglich gestrichen werden, weil
  `make docs-check`s `matrix`-Modul sie verbot) sah dieser Slice-Plan die
  ADR-freie Form der Nachzugzeilen bereits von Anfang an vor (§2 DoD-Punkt
  „Kein ADR-Bezug in der Nachzug-Zeile selbst") — die gelernte Regel aus dem
  Vorgänger-Slice trug direkt beim Schreiben, kein zweiter Anlauf nötig. Der
  `uv`-Digest wurde real gegen die Registry reverifiziert
  (`docker buildx imagetools inspect ghcr.io/astral-sh/uv:0.12.17`) und war
  identisch zum in `ADR-0108` genannten Digest
  (`sha256:10787c682e4184e4f290de1171fd4703dc63de99221f10fe1c99002ce7fa9acc`).
  Der `tar`-Stream-Export trug beim ersten Versuch für **beide** Artefakte
  in einem Stream — die im Slice-Plan §6 offen benannte Unsicherheit
  (zwei Dateien statt einer) bestätigte sich als unbegründet.
- **Was ging anders als geplant:** Nichts Wesentliches. Eine Nebenbeobachtung
  ohne Konsequenz: `uv build` legt selbständig eine `.gitignore` (Inhalt
  `*`) im Ausgabeverzeichnis (`/out` im Container, `sdks/python/dist/` nach
  dem Export) ab — harmlos, weil `sdks/python/.gitignore` bereits `dist/`
  vollständig ausschließt (`git status --porcelain sdks/python/dist/`
  bleibt leer).
- **Steering-Loop-Eintrag:** kein neuer Sensor, keine geschärfte Regel — der
  reale Fund unten (unclosed-backtick-Fernwirkung) bleibt unter der
  3×-Schwelle (Erstauftreten), Kandidat für eine künftige Schärfung des
  Reviewer-Skills ist bereits im Beobachtungs-Eintrag benannt.
- **Beobachtungs-Register (`../observations/`):** eine neue Beleg-Datei
  angelegt —
  [`BEO-PGC/unclosed-backtick-taeuscht-nackte-id-vor`](../observations/BEO-PGC/unclosed-backtick-taeuscht-nackte-id-vor/observation.md)
  (Erstauftreten). Der Reviewer dieses Slices bemerkte im eigenen
  Report-Entwurf einen unclosed-backtick-Fehler (Codespan-Parität
  verschoben), korrigierte ihn vor dem eigenen Commit — der committete
  Report trägt ihn nicht mehr — und vermutete, dasselbe Muster könnte auch
  im bereits gemergten `docs/reviews/review-slice-sdk-csharp-pack-werkzeug.md`
  vorliegen, ohne dies selbst zu prüfen. Die Planner-Closure hat die
  Vermutung real nachgemessen: ein Backtick-Gesamtzahl-Zähllauf
  (`grep -o` + `wc -l`) gegen diese Datei ergab 403 (ungerade) — ein realer
  unclosed-backtick-Defekt an Zeile 122f. (`kein Gate,` ohne schließendes
  Backtick vor der Fortsetzung auf der Folgezeile), unabhängig vom
  Inhalt/Diff dieses Slices, in einem bereits gemergten Nachbar-Report.
  Behoben in einem eigenen, separat committeten Fix (ein eingefügtes
  schließendes Backtick, Gesamtzahl danach 404/gerade). `make gates`/
  `make docs-check` liefen über den gesamten Zeitraum des Defekts
  wiederholt grün — kein `id-unlinked`-Fehlalarm in diesem konkreten Fall,
  siehe Beobachtungs-Eintrag für die Abgrenzung zu
  `BEO-PGC/report-nackte-id-ohne-link` (dort: vergessene Verlinkung an der
  eigenen Stelle; hier: Fernwirkung eines entfernten Syntaxfehlers).
- **Folge-Slices:** `slice-sdk-python-publish-workflow` — bereits als Datei
  in `open/` vorhanden.
- **Risiken aus §6:**
  - „Export-Mechanismus für `.whl`/`.tar.gz` neu, zwei Dateien statt einer" —
    **Ausgang: entfallen/aufgelöst** — real bestätigt, beide Artefakte im
    selben `tar`-Stream, keine Dateirechte-Auffälligkeit.
  - „Neue `SPEC-<NNN>`-Nummer könnte kollidieren" — **Ausgang: entfallen** —
    `SPEC-027` real kollisionsfrei vergeben (`grep -c "SPEC-027"
    spec/*.md` → genau drei Treffer, ausschließlich in
    `spec/pflichtenheft.md`).
  - „PyPI-Namensraum-Kollision für `pgchangefeed`" — **Ausgang: weiter
    offen**, wie im Plan bereits vorgesehen — ohne Wirkung für dieses
    netzlose Pack-Werkzeug-Slice, zu prüfen vor dem ersten realen
    Publish-Versuch (`slice-sdk-python-publish-workflow`).
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-python-lh-fa-sst-009](../welle-sdk-python-lh-fa-sst-009.md)
  (noch offen) — die Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Zwei Sub-Areas: `sdks/python/`
(bereits eröffnet, GF) und `spec/pflichtenheft.md`/`harness/README.md`
(Default-Sub-Area `*`/`PGC`, Greenfield, `harness/conventions.md`).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
kein Treffer, der speziell dieses Pack-Werkzeug-Slice beträfe, über die
bereits in der Welle-Plan §6 benannten Einträge hinaus.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
