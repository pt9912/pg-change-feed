# ADR-0134: `sdk-public-doc-check` wird Gate in `make gates`

**Status:** Accepted

**Datum:** 2026-09-29

**Autor:** pt9912 (Rolleninhaber: Architect-Lauf, 2026-09-29)

**Bezug:** [`LH-FA-SST-009`](../../../spec/lastenheft.md) (offizielle
Client-Bibliotheken — ihre Kommentare, Docstrings, KDoc, XML-Doku und
Fehlertexte erreichen Anwender über die Pakete, deshalb trägt dort keine
interne Kennung), [ADR-0110](0110-python-sdk-umfang-erweitert-vollmatrix.md)
(Umfang der Packages — der Wächter prüft alle drei SDK-Bäume),
[ADR-0106](0106-csharp-nuget-erstes-sdk-package.md) /
[ADR-0107](0107-python-pypi-zweites-sdk-package.md) /
[ADR-0109](0109-kotlin-github-packages-drittes-sdk-package.md) (die drei
`sdk-pack-*`-Ziele, deren Vorgänger-Kante auf den Wächter die bisherigen
Fangstellen trägt), `AGENTS.md` §3.6 (Träger-Pflicht — Änderungen an Gate
und Gate-Liste sind ADR, kein PR-Kommentar)

**Schärft:** — (Prozess-ADR ohne Spec-Stratum; die Harness-Werkzeug-Regel
braucht keine Spec-Stelle, [`LH-FA-SST-009`](../../../spec/lastenheft.md)
bleibt unverändert und wird hier nur als Bezug genannt)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Seit `slice-sdk-readme-nutzerdoku` prüft
`tools/harness/sdk-public-doc-check.sh`, dass keine Datei unter `sdks/`
eine interne Kennung trägt — `SPEC-`/`ADR-`/`ARC-`/`LH-FA-`/`LH-QA-`
samt Slice-/Welle-Namen —, ausgenommen Bau-Ausgaben und erzeugten Code
(`obj`, `bin`, `build`, `dist`, `.gradle`, `__pycache__`, `.pytest_cache`,
`*.egg-info`, `grpc_gen`). Reines `find`/`grep`, kein Docker, kein Netz.
Die drei `make sdk-pack-*`-Ziele hängen an ihm als Vorgänger-Kante, die
drei `sdk-*-release.yml`-Workflows fahren sie vor jedem Publish.

Daraus folgt die heutige Fangstellen-Lage: Ein Rückfall — etwa eine
interne Kennung in einem SDK-Docstring — fällt beim Packen oder beim
Release auf, nicht bei jedem Push. Er wird nicht ausgeliefert, aber
er fällt nicht dort auf, wo die übliche Arbeit sichtbar ist (PR und
Push, `ci.yml` fährt `make gates` plus `make test`).

Gegen die Aufnahme als Gate stand bislang eine Begründung in
`harness/mk/sdk.mk` und `harness/README.md` §Sensors: „ein weiteres Gate
ändert die Gate-Liste und ihre Sensor-Bindung“. Sie benennt den Aufwand
einer Gate-Listen-Änderung; als Grund gegen einen netzlosen grep-Lauf
trägt sie nicht — gemessen (2026-09-25, Slice-Plan
`docs/plan/planning/open/slice-sdk-public-doc-check-gate.md`; 2026-09-29
durch diesen Lauf bestätigt, `real 0m0,017s`).

Der Wächter ist der siebte Gate-Kandidat neben den sechs bestehenden
Einträgen in `GATE_CHECKS` (baseline-verify, docs-check, a-check,
commit-traceability, coverage-gate, generated-sync — aus den
`GATE_CHECKS +=`-Zeilen der Make-Fragmente hergeleitet). Diese ADR ist
der Träger, den `AGENTS.md` §3.6 für eine Änderung der Gate-Liste
verlangt — sie entsteht **vor** der Verdrahtung, nicht als Nachzug zu
einem Config-Commit. Das ist die Lehre aus der Beobachtung
`BEO-PGC/gate-scope-erweiterung-ohne-adr-traeger` (offen, 2×): Gate-Umfang
wächst still, wenn der Träger nach der Änderung kommt. Der Umsetzungs-Slice
(`docs/plan/planning/open/slice-sdk-public-doc-check-gate.md`) startet
deshalb erst mit dieser `Accepted` ADR.

## Entscheidung

Wir wählen **die Aufnahme des Wächters in `GATE_CHECKS`** — netzlos, im
Repo-Wurzel-`make gates`, mit unverändertem Muster- und Prüfumfang — mit
drei Teilfragen.

### Teilfrage 1 — Aufnahme als Gate

| Option | Pro | Contra |
|---|---|---|
| A — Status quo: Werkzeug mit Vorgänger-Kante an `sdk-pack-*` | keine Gate-Listen-Änderung | der Rückfall fällt erst beim Pack oder Release auf — Orte, die nicht an jedem Push laufen; die bisherige Gegen-Begründung (Gate-Listen-Aufwand) ist an einer 0,017-s-Messung (gemessen 2026-09-29) nicht mehr tragfähig |
| **B — Aufnahme in `GATE_CHECKS` (gewählt)** | deterministischer, netzloser grep ohne Docker — die Klasse von Prüflauf, für die `make gates` gebaut ist; der Rückfall fällt beim PR/Push auf, dem Ort der üblichen Arbeit; der Umfang des Wächters bleibt exakt derselbe (keine Muster- oder Reichweiten-Änderung, siehe Slice-Plan §1) | die Gate-Liste wächst um einen Eintrag — mit `real 0m0,017s` (gemessen 2026-09-29) der billigste Eintrag der Liste |
| C — eigener, advisory-GitHub-Workflow statt Gate-Aufnahme | sichtbar im PR, ohne die Gate-Liste zu ändern | schafft einen zweiten, driftenden Ort (Workflow-Datei) neben `make gates` für eine Prüfung, die deterministisch und billig ist; „rot, aber blockiert nicht“ ist die falsche Klasse für einen Befund, dessen Prüfkosten gegen null fallen |

### Teilfrage 2 — der Tabellentest

`make test-sdk-public-doc-check` belegt die Wächter-Logik (elf Fälle: je
Kennungsart ein Treffer, dazu die Ausnahmen). Er ist nicht der Gegenstand
des Gates — der Gegenstand ist der Baum.

| Option | Pro | Contra |
|---|---|---|
| **A — Test bleibt Werkzeug (gewählt)** | der Gegenstand des Gates ist der Baum, nicht der Wächter; ein Defekt der Wächter-Logik ist ein Werkzeug-Defekt — dieselbe Klasse wie die Tabellentests der übrigen Harness-Werkzeuge (`test-fmt-check`, `test-suchlauf-nachmessen`, `test-kommentar-kennungen`, alle kein Gate); die Gate-Liste bleibt eine Liste von Baum-Prüfungen | die Wächter-Logik ist im Gate-Lauf nicht gegen ihren Test gesichert — akzeptiertes Negativ, siehe Konsequenzen |
| B — zweiter `GATE_CHECKS`-Eintrag für den Test | die Wächter-Logik wird je Gate-Lauf mitgeprüft | bläht die Gate-Liste um einen zweiten Eintrag für dasselbe Anliegen auf — ein Gate prüft den Baum, der andere die Logik eines Werkzeugs; verwischt die Trennung zwischen Gate (Baum) und Werkzeug-Test |
| C — Test als Vorgänger-Kante am Gate-Target (`sdk-public-doc-check: test-sdk-public-doc-check`) | ein Gate-Eintrag, die Logik wird je Lauf mitgeprüft | koppelt jeden manuellen Aufruf — auch die drei `sdk-pack-*`-Ziele — an den Test; der Prüflauf färbt dann aus zwei Gründen rot (Baum oder Logik), was die Diagnose verwässert, ohne dass ein Akzeptanzkriterium sie verlangt |

### Festlegungen

1. **Verdrahtung.** `harness/mk/sdk.mk` hängt das Ziel an
   `GATE_CHECKS += sdk-public-doc-check`; `make gates` fährt es in jedem
   Lauf. Netzlos, kein Docker; das Skript steigt über
   `git rev-parse --show-toplevel` auf die Baum-Wurzel und prüft
   `sdks/`. Muster, Ausnahmen und Reichweite des Wächters bleiben
   unverändert.
2. **Tabellentest.** `make test-sdk-public-doc-check` bleibt Werkzeug
   (Teilfrage 2, Option A).
3. **Sensor-Bindung — Pflicht des Implementer-Zugs.** Künftige Träger:
   `harness/sensors/sdk-public-doc-check.md` (Vertrag, Grenze, Bindung —
   in der Form der übrigen Sensor-Dateien, z. B.
   `harness/sensors/generated-sync.md`) und die Gate-Zeile in
   `harness/README.md` §Sensors. Die „kein Gate“-Begründung in
   `harness/mk/sdk.mk` und `harness/README.md` entfällt bzw. wird auf den
   Ist-Zustand gezogen (Suchlauf nach `AGENTS.md` §3.13 über
   `git grep -n 'sdk-public-doc-check'`). Der Slice-Plan §3 führt beide
   Dateien.
4. **Reichweite am Arbeitsbaum (Grenze).** Das Gate läuft an jedem
   Arbeitsbaum-Zustand — auch an einem mit unversionierten Bau-Ausgaben
   unter `sdks/`, denn das Skript liest den Baum, nicht den git-Index.
   Die Prune-Liste des Skripts trägt diese Klasse: erprobt an den
   Tabellentest-Fällen `obj`, `dist` und `grpc_gen` (Exit 0 je Fall) und
   am realen Baum (gemessen 2026-09-29: die ignorierten Ausgaben
   `sdks/csharp/dist`, `sdks/python/dist`, `sdks/kotlin/dist` und
   `sdks/python/pgchangefeed/src/pgchangefeed/grpc_gen` sind vorhanden,
   `make sdk-public-doc-check` endet Exit 0). Eine **neue** Klasse
   erzeugter Verzeichnisse, die nicht in der Prune-Liste steht, färbt das
   Gate falsch rot — die Richtung des Fehlers ist die sichere (rot statt
   grün, dieselbe Begründungsform wie `harness/sensors/generated-sync.md`
   §Grenze 2); die Erweiterung der Prune-Liste ist eine Änderung am
   Wächter, die die Sensor-Datei nachzieht. Eine **Milderung** des
   Wächters (mehr Ausnahmen, weniger Muster) ist eine Umfangs-Änderung
   dieses Gates im Sinn von `AGENTS.md` §3.6 und braucht ihren eigenen
   Träger — kein stiller Skript-Commit.

## Konsequenzen

- Positiv: Ein Rückfall fällt bei jedem Push und PR auf
  (`ci.yml` fährt `make gates`), nicht erst beim Packen oder Release.
- Positiv: Der Wächter bekommt einen benannten Vertrag
  (`harness/sensors/sdk-public-doc-check.md`) — bislang führt nur die
  Werkzeug-Zeile in `harness/README.md` §Sensors ihn.
- Positiv: An den drei `sdk-pack-*`-Zielen und ihren Release-Workflows
  ändert sich nichts — die Vorgänger-Kante bleibt, der Rückfall fällt
  weiterhin vor dem Bau auf.
- Negativ: Die Gate-Liste wächst um einen Eintrag — mit
  `real 0m0,017s` (gemessen 2026-09-29) der billigste der sieben.
- Negativ (Grenze): Falsch rot an einem Arbeitsbaum mit einer neuen,
  nicht ausgenommenen Klasse erzeugter Verzeichnisse unter `sdks/` —
  fail-closed, Behandlung in Festlegung 4.
- Negativ (akzeptiertes Negativ): Der Gate-Lauf prüft die Wächter-Logik
  nicht gegen ihren Tabellentest (Teilfrage 2, Option A). Der Fall tritt
  nur ein, wenn jemand am Wächter ändert — und genau dann liest der
  Reviewer den Diff gegen die Sensor-Datei; ein schleichender Defekt ohne
  eine Änderung am Wächter ist nicht der zu erwartende Weg.
- Folgepflicht (Implementer-Zug, DoD des Slices): `GATE_CHECKS +=` in
  `harness/mk/sdk.mk`, Sensor-Datei, Gate-Zeile in `harness/README.md`
  §Sensors, „kein Gate“-Begründung ziehen; Beleg der Eingabeseite — eine
  eingefügte Kennung unter `sdks/` färbt `make gates` rot (Exit ≠ 0,
  Wächter-Meldung gedruckt), die Rücknahme grün.
- Folgepflicht: `spec/` und `spec/architecture.md` brauchen **keine**
  Änderung — Prozess-ADR ohne Spec-Stratum; [`LH-FA-SST-009`](../../../spec/lastenheft.md)
  und `ARC-*` bleiben unverändert.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Make-Verdrahtung (erwartet — zu belegen durch den Implementer-Zug) | `sdk-public-doc-check` steht in `GATE_CHECKS`; `make gates` fährt es; eine eingefügte Kennung unter `sdks/` färbt den Lauf rot (Exit ≠ 0), die Rücknahme grün | `make gates` |
| Tabellentest (erprobt) | elf Fälle: je Kennungsart (SPEC, ADR, ARC, LH, slice, welle) ein Treffer mit Exit 1; die Ausnahmen `obj`, `dist`, `grpc_gen` und eine saubere Datei mit Exit 0 — gemessen 2026-09-29, alle Fälle bestanden (`real 0m0,128s`) | `make test-sdk-public-doc-check` |
| Wächter am realen Baum (erprobt) | Arbeitsbaum mit den ignorierten Bau-Ausgaben `sdks/csharp/dist`, `sdks/python/dist`, `sdks/kotlin/dist`, `sdks/python/pgchangefeed/src/pgchangefeed/grpc_gen` — Exit 0; gemessen 2026-09-29 (`real 0m0,017s`) | `make sdk-public-doc-check` |

## Re-Evaluierungs-Trigger

Permanent — die Aufnahme selbst hat keine Prämisse, die entfallen kann:
solange `sdks/` Quellen trägt, die Anwender über die Pakete lesen, bleibt
die Prüfung Gegenstand. Ändert sich der Prüfumfang des Wächters (Muster,
Prune-Liste), ändert das den Vertrag dieses Gates und braucht seinen
eigenen Träger (Festlegung 4).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-29 | Accepted — Architect-Entscheidung als Vorstufe des Umsetzungs-Slices | `docs/plan/planning/open/slice-sdk-public-doc-check-gate.md` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0134` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
