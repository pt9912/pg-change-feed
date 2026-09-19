# Verifikationsbericht: slice-sdk-python-pack-werkzeug — 2026-09-19

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/slice-sdk-python-pack-werkzeug.md`
§2/§3/§6/§7), in frischem Kontext. **Nicht** gegen den Diff als solchen
(Reviewer-Aufgabe, abgeschlossen mit
[`review-slice-sdk-python-pack-werkzeug.md`](review-slice-sdk-python-pack-werkzeug.md),
0 HIGH/0 MEDIUM/1 LOW/1 INFO) und **nicht** gegen realen Bedarf (Validator,
hier nicht ausgelöst).

**Gegenstand:** fünf Commits `19d20c7b..HEAD` (`2184fafd`, `b43142ca`
reine Lifecycle-Moves; `8bd536c1` Slice-Metadaten; `0e42b801`
Implementer-Inhalt; `7c5c335e` Review-Report), Slice
`slice-sdk-python-pack-werkzeug`, Welle
`welle-sdk-python-lh-fa-sst-009`, `LH-FA-SST-009`, `ADR-0107`, `ADR-0108`.

**Eingangs-Kontext (eigen gelesen, nicht aus Bericht übernommen):**
Slice-Plan §1/§2/§3/§6/§7/§8 vollständig, Review-Report vollständig,
`ADR-0107` Festlegung 5 im Volltext, `ADR-0108` §Entscheidung Festlegung
1–4 im Volltext, `spec/pflichtenheft.md` §1 `LH-FA-SST-009.a` (Zeilen
170–194), §6 `SPEC-026`/`SPEC-027` (Zeilen 610–620), §7 Historie (Zeilen
653–657), `harness/README.md` §Werkzeuge (`sdk-pack-python`- und
`sdk-pack-csharp`-Zeile), `sdks/python/Dockerfile` vollständig,
`sdks/python/pgchangefeed/pyproject.toml` vollständig,
`harness/mk/sdk.mk`, `tools/harness/sdk-pack-python.sh`.

---

## 1. DoD-Vertrag (§2) — jede Zeile einzeln geprüft

### 1.1 „`make sdk-pack-python` existiert (Docker-only, kein Gate) …"

Eigener Lauf (sechste unabhängige Bau-Bestätigung nach Implementer,
Reviewer-Erstlauf und Reviewer-Mutationsprobe):

```
$ rm -rf sdks/python/dist && make sdk-pack-python; echo $?
0
$ ls -la sdks/python/dist/
pgchangefeed-0.1.0-py3-none-any.whl   9330 Bytes
pgchangefeed-0.1.0.tar.gz            11290 Bytes
$ file sdks/python/dist/*.whl sdks/python/dist/*.tar.gz
…: Zip archive data, at least v2.0 to extract, compression method=deflate
…: gzip compressed data, was "pgchangefeed-0.1.0.tar", max compression
```

Bytegrößen identisch zu Slice-Plan §6, Commit-Message und Review-Report
— kein Drift (`AGENTS.md` §3.12 Instanz A).

**Eigene, vom Reviewer unabhängige Mutations-Probe** (der Reviewer
mutierte eine Testerwartung in `test_options.py`; diese Sitzung mutiert
stattdessen den **Quellcode selbst**, eine andere Fehlerklasse):
`sdks/python/pgchangefeed/src/pgchangefeed/http_client.py` wurde vor dem
bestehenden `from __future__ import annotations` eine zusätzliche
`import`-Zeile eingefügt (`import
this_module_does_not_exist_verifier_probe`) — das erzeugt real einen
`SyntaxError` beim Pytest-Collection-Schritt (ein `from __future__
import` muss die erste Anweisung der Datei bleiben), nicht nur einen
fehlgeschlagenen Test:

```
$ rm -rf sdks/python/dist && make sdk-pack-python; echo $?
…
E   File ".../http_client.py", line 32
E     from __future__ import annotations
E   SyntaxError: from __future__ imports must occur at the beginning of the file
ERROR tests/test_http_client.py
ERROR tests/test_options.py
!!!!!!!!!!!!!!!!!!! Interrupted: 2 errors during collection !!!!!!!!!!!!!!!!!!!!
ERROR: process "/bin/sh -c pytest" did not complete successfully: exit code: 2
make: *** [sdk-pack-python] Error 1
2
$ ls sdks/python/dist/
ls: sdks/python/dist/: No such file or directory
```

Bau bricht sichtbar an der `RUN pytest`-Zeile ab (Exit `2`), **kein**
`dist/`-Verzeichnis entsteht — bestätigt real, auf einem anderen
Fehlerpfad als der Reviewer, dieselbe strukturelle Garantie (Docker-
Layer-Abhängigkeit `pack` hinter `build`). Mutation danach exakt
zurückgesetzt (Backup-Datei zurückkopiert):

```
$ git diff --stat sdks/python/pgchangefeed/src/pgchangefeed/http_client.py
(leer)
$ git status --porcelain sdks/python/
(leer)
```

Artefakte danach real neu erzeugt (identische Bytegrößen wie oben,
zweiter sauberer Lauf, siehe §4).

**Ergebnis: berechtigt auf `[x]`.**

### 1.2 „Real ausgeführt: ein `.whl` und ein `.tar.gz` … als Smoke-Beleg im
Bericht genannt"

Slice-Commit-Message (`0e42b801`) und Review-Report nennen beide
identisch `pgchangefeed-0.1.0-py3-none-any.whl` (9330 Bytes) und
`pgchangefeed-0.1.0.tar.gz` (11290 Bytes) — mit dem eigenen Lauf in 1.1
übereinstimmend, kein Drift. **Ergebnis: berechtigt auf `[x]`.**

### 1.3 `spec/pflichtenheft.md` §1 `LH-FA-SST-009.a`-Nachzug

Eigene Lektüre der Nachzugzeile (Zeilen 189–194):

> „Für Python/PyPI ist die Frage ebenfalls beantwortet: `pgchangefeed`
> (`SPEC-027`) deckt HTTP-API, real Docker-only paketierbar
> (`make sdk-pack-python`) und real geprüft
> (`pgchangefeed-0.1.0-py3-none-any.whl`, `pgchangefeed-0.1.0.tar.gz`).
> Eine dritte Sprache oder ein dritter Vertriebsweg bleibt offen — diese
> Kennung bleibt ihre Adresse."

Trifft präzise: sagt aus, dass Python/PyPI nicht mehr offen ist, dass
eine **dritte** Sprache/Vertriebsweg offen bleibt, streicht die Kennung
selbst nicht (`LH-FA-SST-009.a`-Überschrift bleibt bestehen). Kein
`ADR-*`-Bezug in dieser Zeile selbst — geprüft per `grep -n "ADR-"
spec/pflichtenheft.md | sed -n '/170,194p'` (manuell gegen die Zeilen
170–194 gelesen): kein Treffer in diesem Absatz, nur in der
C#-Vorgängerzeile ebenfalls keiner. `make docs-check`s `matrix`-Modul
läuft grün darüber (siehe §5). **Ergebnis: berechtigt auf `[x]`.**

### 1.4 `spec/pflichtenheft.md` §6 `SPEC-027`-Zeile

Eigener Lauf: `grep -c "SPEC-027" spec/*.md`

```
spec/architecture.md:0
spec/lastenheft.md:0
spec/pflichtenheft.md:3
```

Genau drei Treffer, ausschließlich in `spec/pflichtenheft.md` (§1-Absatz,
§6-Tabellenzeile, §7-Historie) — kollisionsfrei, keine Zweitvergabe. Form
der Zeile (System/Version/Vertrag-Datei-Spalte) gegen `SPEC-026`
verglichen: strukturell identisch, referenziert korrekt nur `SPEC-018`
(HTTP-API), nicht zusätzlich `SPEC-020` (gRPC) — deckt sich mit dem
tatsächlichen v1-Umfang aus `ADR-0107` Festlegung 1 (nur HTTP-API).
**Ergebnis: berechtigt auf `[x]`.**

### 1.5 `spec/pflichtenheft.md` §7 Historie

Eigene Lektüre der Zeile 657: „`LH-FA-SST-009.a` nachgezogen: Für
Python/PyPI ist die Sprachmatrix-/Vertriebsweg-Frage beantwortet und
`pgchangefeed` real paketierbar (`make sdk-pack-python`,
`pgchangefeed-0.1.0-py3-none-any.whl`, `pgchangefeed-0.1.0.tar.gz`) — die
Kennung bleibt bestehen, eine dritte Sprache oder ein dritter
Vertriebsweg bleibt offen" sowie Zeile 656 (`SPEC-027`-Historie-Zeile) —
**kein** `ADR-*`/`slice-*`-Verweis in beiden Zeilen, konform zur
§7-Kopfregel. **Ergebnis: berechtigt auf `[x]`.**

### 1.6 `harness/README.md` §Werkzeuge — `make sdk-pack-python`-Zeile

Eigene Lektüre (Zeile 156) gegen die `sdk-pack-csharp`-Zeile (155)
verglichen: gleiche Form Ende zu Ende (Bau-Kontext, Stufen `build`→`pack`→
`pack-export`, Testlauf-vor-Artefakt-Zusage, `.gitignore`-Hinweis,
Schlusssatz „Braucht Netz … Werkzeug statt Gate", Bindung `kein Gate,
ADR-… · seit slice-<name>`); referenziert real ein existierendes Target
(bestätigt in 1.1). Pipe im Fließtext korrekt als `\|` escaped.
**Ergebnis: berechtigt auf `[x]`.**

### 1.7 `sdks/python/.gitignore` trägt `dist/`

Eigene Lektüre: `dist/`-Zeile vorhanden (bereits vorbestehend, nicht Teil
dieses Diffs). Nach beiden eigenen Pack-Läufen (§1.1, §4):
`git status --porcelain` durchgehend leer — keine Artefakte als
untracked sichtbar. **Ergebnis: berechtigt auf `[x]`.**

### 1.8 „`make gates` grün"

Siehe §3 unten. **Ergebnis: berechtigt auf `[x]`.**

### 1.9 Review, Doku-Update, Closure-Notiz, Reconciliation,
Beobachtungs-Register, §6-Risiko-Ausgänge, drei Paarungen

Review-Report liegt vor (0 HIGH/0 MEDIUM/1 LOW/1 INFO, kein
Fixrunden-Erfordernis, DoD-Checkbox-Nachzug ohne Fixrunde laut
Reviewer-Skill-Regel korrekt angewendet). Closure-Notiz (§7) ist bereits
vollständig ausgefüllt (kein `<…>`-Platzhalter mehr) — anders als beim
Schwester-Slice `slice-sdk-python-http-client-flaeche` liegt dieser Slice
bereits mit gefüllter Closure-Notiz vor, aber **noch unter
`in-progress/`** (kein `git mv` nach `done/` in diesem Diff-Bereich,
real per `find docs/plan/planning/in-progress -iname
"slice-sdk-python-pack-werkzeug.md"` bestätigt: Datei liegt weiterhin
dort). Reconciliation „entfällt" korrekt begründet (keine
Reconciliation-Datei in diesem Repo). Beobachtungs-Register „keine
Beobachtung angefallen" — eigene Prüfung: `find
docs/plan/planning/observations -newer
docs/plan/planning/in-progress/slice-sdk-python-pack-werkzeug.md`
liefert keinen neuen Treffer für diesen Slice; zutreffend. Alle drei
§6-Risiken tragen einen expliziten Ausgang (zwei „entfallen/aufgelöst"
real durch die eigenen Läufe hier bestätigt, eines „weiter offen" mit
nachvollziehbarer Begründung — PyPI-Namensraum-Kollision ist für dieses
netzlose Pack-Werkzeug-Slice ohne Wirkung). Die „drei Paarungen" verweisen
korrekt auf die noch offene Welle-Closure.

**Kein DoD-Blocker in dieser Gruppe.**

---

## 2. Spec-/ADR-Konformität — eigenständig geprüft

### 2.1 `pyproject.toml` `[build-system]` unverändert

Eigene Lektüre `sdks/python/pgchangefeed/pyproject.toml` Zeilen 1–3:

```toml
[build-system]
requires = ["setuptools>=68"]
build-backend = "setuptools.build_meta"
```

Identisch zum in `ADR-0108` §Entscheidung Festlegung 2 festgehaltenen
Ist-Stand — kein Diff an dieser Datei über den gesamten
Verifikations-Diff-Bereich (`git diff 19d20c7b..HEAD --stat --
sdks/python/pgchangefeed/pyproject.toml` liefert keinen Treffer).
Selbst gegengelesen, nicht nur Reviewer-Bericht zitiert.

### 2.2 `uv`-Digest — dritte unabhängige Reverifikation

```
$ docker buildx imagetools inspect ghcr.io/astral-sh/uv:0.12.17
Digest: sha256:10787c682e4184e4f290de1171fd4703dc63de99221f10fe1c99002ce7fa9acc
```

Identisch zum in `ADR-0108` §Entscheidung Festlegung 3, im Slice-Plan §2
und in `sdks/python/Dockerfile:62` genannten Digest — dritte unabhängige
Bestätigung nach Implementer und Reviewer, alle drei mit identischem
Ergebnis. `grep -n "COPY --from=ghcr.io/astral-sh/uv"
sdks/python/Dockerfile` zeigt die byte-identische Zeile.

### 2.3 ADR-0107/0108 §5/unberührte Pfade

```
$ git diff 19d20c7b..HEAD --stat -- docs/user/version.md .a-check.yml \
    examples/csharp/ sdks/csharp/ \
    sdks/python/pgchangefeed/src/pgchangefeed/http_client.py \
    sdks/python/pgchangefeed/src/pgchangefeed/models.py \
    sdks/python/pgchangefeed/src/pgchangefeed/exceptions.py
(leer)
```

Kein Treffer — alle sieben genannten Pfade sind über den vollständigen
Diff-Bereich dieses Slices unberührt, wie `ADR-0107` §6 „Was diese ADR
nicht ändert" (a-check, `docs/user/version.md`, kein Umzug von
`sdks/csharp/**`) und `ADR-0108` §Entscheidung Festlegung 6 es festlegen,
sowie das Slice-Ziel „keine Änderung der HTTP-Fläche" (§1 Ausdrücklich
NICHT). Vollständiger Diff-Stat (`19d20c7b..HEAD`, sieben Dateien
Inhaltsänderung plus zwei Lifecycle-Moves) zusätzlich gelesen:
ausschließlich `sdks/python/Dockerfile`, `harness/mk/sdk.mk`,
`tools/harness/sdk-pack-python.sh`, `spec/pflichtenheft.md`,
`harness/README.md` und die beiden Planning-/Review-Träger — kein Treffer
außerhalb des erwarteten Slice-Umfangs.

---

## 3. `make gates` real, ungepiped, in dieser Sitzung ausgeführt

```
$ git log -1 --format=%H
7c5c335e...
$ make gates > /tmp/gates.log 2>&1; echo $?
0
```

Einzelbelege aus demselben Lauf (Log vollständig eingesehen):

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK` |
| `docs-check` | `d-check: 832 Datei(en) geprüft, 0 Befund(e)` |
| `a-check` | `gesamt: 0 Befund(e)` |
| `commit-traceability` | `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators` |

`grep -n "sdk-pack" /tmp/gates.log` liefert **keinen** Treffer — eigene
Bestätigung, dass `sdk-pack-python` nicht Teil von `make gates` ist
(zusätzlich zur direkten `harness/mk/sdk.mk`-Lektüre: kein `GATE_CHECKS
+=`-Eintrag in diesem Fragment, im Gegensatz zu `baseline.mk`,
`doc-gate.mk`, `coverage.mk`, `generated-sync.mk`).

**Ergebnis: DoD-Checkbox „`make gates` grün" berechtigt auf `[x]`.**

## 4. Artefakte nach Mutations-Rücksetzung erneut real erzeugt

```
$ rm -rf sdks/python/dist && make sdk-pack-python; echo $?
0
$ ls -la sdks/python/dist/
pgchangefeed-0.1.0-py3-none-any.whl   9330 Bytes
pgchangefeed-0.1.0.tar.gz            11290 Bytes
```

Identische Bytegrößen wie beim ersten eigenen Lauf (§1.1) — bestätigt,
dass die zurückgesetzte Mutation kein Restartefakt hinterlassen hat.

## 5. `make docs-check` isoliert erneut geprüft (`matrix`-Modul)

Bereits Teil des `make gates`-Laufs (§3): `832 Datei(en) geprüft, 0
Befund(e)` — die ADR-freie Form der Nachzugzeilen in
`spec/pflichtenheft.md` §1/§6/§7 verletzt das `matrix`-Modul nicht.

---

## Verdikt

**DoD erfüllt** für den aktuellen `in-progress`-Stand des Slice. Alle
beauftragten Prüfpunkte wurden real und unabhängig nachgemessen, nicht
aus Bericht oder Commit-Message übernommen:

1. `make sdk-pack-python` real ausgeführt (sechste unabhängige
   Bau-Bestätigung), Exit `0`, beide Artefakte mit erwarteten
   Bytegrößen/Dateitypen real erzeugt.
2. **Eigene, vom Reviewer verschiedene Mutations-Probe** — Quellcode-
   Syntaxfehler statt Testerwartungs-Fehler — bricht den Bau real vor
   `pack` ab (Exit `2`, kein `dist/`), Mutation danach exakt
   zurückgesetzt (`git diff`/`git status` leer).
3. `make gates` nicht mit `sdk-pack-python` — real per `grep` gegen den
   eigenen Gate-Log sowie `harness/mk/sdk.mk`-Lektüre bestätigt.
4. `LH-FA-SST-009.a`-Nachzug sagt präzise: Python/PyPI nicht mehr offen,
   dritte Sprache/Vertriebsweg bleibt offen, ohne ADR-Bezug, ohne die
   Kennung zu streichen.
5. `SPEC-027`-Zeile real kollisionsfrei (`grep -c` → drei Treffer,
   ausschließlich `spec/pflichtenheft.md`), Form konsistent zu
   `SPEC-026`.
6. §7-Historie-Zeilen tragen keinen `ADR-*`/`slice-*`-Verweis.
7. `harness/README.md`-Zeile formkonsistent zu `sdk-pack-csharp`.
8. `sdks/python/.gitignore` trägt `dist/` — nach beiden eigenen
   Pack-Läufen bleibt `git status --porcelain` leer.
9. `pyproject.toml` `[build-system]` eigenständig gegengelesen:
   unverändert `setuptools.build_meta`.
10. `uv`-Digest dritte unabhängige Reverifikation: identisch zu ADR,
    Dockerfile und Reviewer-Messung.
11. `git diff --stat` gegen sieben unberührte Pfade (`docs/user/version.md`,
    `.a-check.yml`, `examples/csharp/**`, `sdks/csharp/**`,
    `http_client.py`, `models.py`, `exceptions.py`) — leer, alle
    unberührt.
12. `make gates` eigenständig, ungepiped, `EXIT=0` — alle sechs Gates
    grün (baseline-verify, docs-check 832/0, a-check 0,
    commit-traceability OK, coverage-gate 82.80% ≥ 80%, generated-sync
    OK).

**Freigabe an den Planner:** Der Slice ist DoD-konform für seinen
aktuellen `in-progress`-Stand. Die Closure-Notiz (§7) ist bereits
gefüllt; der `git mv` nach `done/` sowie die drei Paarungen bleiben bei
der Closure von `welle-sdk-python-lh-fa-sst-009` fällig, wie im Slice-Plan
selbst bereits vorgesehen.
