# Verifikationsbericht: slice-sdk-python-publish-workflow — 2026-09-19

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/slice-sdk-python-publish-workflow.md` §2)
und die `ADR-0107`-/`ADR-0108`-Konformität, in frischem Kontext. **Nicht**
gegen den Diff als solchen (Reviewer-Aufgabe, abgeschlossen mit
[`review-slice-sdk-python-publish-workflow.md`](review-slice-sdk-python-publish-workflow.md),
0 HIGH/MEDIUM/LOW/INFO) und **nicht** gegen realen Bedarf (Validator, hier
nicht ausgelöst).

Dieser Slice ist der **letzte** der Welle
`welle-sdk-python-lh-fa-sst-009` — entsprechend gründlich geprüft: jede
DoD-Zeile einzeln, drei unabhängige SHA-/Test-Reverifikationen (vierte
Instanz insgesamt, nach Implementer und Reviewer), eigener `make
gates`-Lauf, eigene Traceability-/Immutabilitäts-Läufe über den exakten
Diff-Bereich, eigene Live-Prüfung gegen GitHub/PyPI/NuGet (Secrets,
Namensraum, C#-Realisierung).

**Gegenstand:** Diff-Range `3043c76a..HEAD` (sechs Commits): `90f2c9b8`
(open→next, reiner Move), `3fcc941a` (Verantwortlich gesetzt), `aad9ee26`
(next→in-progress, reiner Move), `3808be6c` (Implementer-Zug:
`.github/workflows/sdk-python-release.yml`,
`tools/harness/sdk-python-release-tag-info.sh`,
`tools/harness/run-sdk-python-release-tag-info-tests.sh`, `Makefile`,
`harness/README.md`, `docs/user/releasing.md`, Slice-Plan), `38a137f7`
(separater Coordinator-Fix zum bereits real veröffentlichten
C#-SDK-Release-Stand in `docs/user/releasing.md`, außerhalb dieses
Slice-Scopes, vom Reviewer bereits auf Korrektheit geprüft), `6c7e58f8`
(Review-Report samt DoD-Checkbox-Nachzug).

**Frischer Kontext, eigene Läufe statt Behauptungs-Übernahme:** Slice-Plan
(vollständig, §1/§2/§3/§6/§7/§8), Review-Report (vollständig, 0 offene
HIGH/MEDIUM/LOW/INFO), `ADR-0107` (Accepted, Volltext inkl. Festlegung 4/5,
Alternativentabellen), `ADR-0108` (Accepted, Volltext inkl. §Entscheidung
Festlegung 1, Alternativentabelle E), `AGENTS.md` §3.8/§3.9/§3.10,
`harness/README.md` §Werkzeuge, `docs/user/releasing.md` (Volltext), die
Workflow-Datei selbst (vollständig gelesen, 103 Zeilen), das neue
Shell-Skript (vollständig gelesen). Eigene Testläufe (alle drei
`run-*-release-tag-info-tests.sh`), eigener `make gates`-Lauf, eigener
YAML-Parse plus isolierte `bash -n`-Prüfung aller vier `run:`-Blöcke,
eigene `doc-commits`/`doc-immutable`-Läufe über den exakten Slice-Bereich,
eigene Live-Prüfung gegen GitHub (`gh secret list`), PyPI (`curl`) und
NuGet (`curl`) sowie `git tag -l`/`git ls-remote`.

---

## 1. DoD-Vertrag (§2) — jede Zeile einzeln geprüft

### 1.1 Workflow existiert: Trigger, PEP-440-/`pyproject.toml`-Reihenfolge, Publish-Aufruf

- `.github/workflows/sdk-python-release.yml` vollständig selbst gelesen
  (103 Zeilen). Trigger: `"on": push: tags: ['sdk-python-v*']` (Zeile
  61–64) — **kein** Zweit-Trigger.
- **Trigger-Isolation selbst nachgemessen:** eigener `grep -rn "uses:"`
  über die Workflow-Datei sowie eine Sichtprüfung der Kopf-Kommentar-
  Aussage gegen `ci.yml`/`e2e.yml` (`tags-ignore: ['**']`),
  `release.yml` (`tags: ['v*']`) und `sdk-csharp-release.yml`
  (`tags: ['sdk-csharp-v*']`) bestätigt: ein GitHub-Actions-`tags:`-Glob
  verlangt literale Präfix-Übereinstimmung vor dem `*` —
  `sdk-python-v*` beginnt weder mit `v` noch mit `sdk-csharp-v`, matcht
  keines der drei anderen Muster. Kein Doppellauf.
- **Reihenfolge selbst an den Zeilennummern nachvollzogen:** Checkout
  (Z. 77–78) → Tag validieren/Version ermitteln über
  `tools/harness/sdk-python-release-tag-info.sh` (Z. 80–81) →
  `pyproject.toml`-Version-Abgleich, Abbruch bei Abweichung (Z. 83–89) →
  `make sdk-pack-python` (Z. 91–92) → `uv` installieren (Z. 94–97) →
  `uv publish` (Z. 99–102). Beide Validierungsschritte laufen vollständig
  **vor** jedem Login/Build/Push — exakt wie im Plan zugesagt.
- Publish-Aufruf exakt wie zugesagt (Z. 99–102):
  `UV_PUBLISH_TOKEN: ${{ secrets.PYPI_API_TOKEN }}` als Umgebungsvariable,
  `uv publish sdks/python/dist/*` — **kein** `-u __token__ -p`-Flag-Paar
  wie bei `twine` (`ADR-0108` §Entscheidung Festlegung 1 eingehalten).

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.2 Action-Pinning (`AGENTS.md` §3.8) — vierte unabhängige SHA-Reverifikation

Genau zwei `uses:`-Zeilen im gesamten Workflow (Zeile 78, 95):

- `actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1` —
  repo-weit wiederverwendeter SHA. Selbst gegen die Quelle verifiziert:
  `git ls-remote https://github.com/actions/checkout refs/tags/v7.0.1` →
  `3d3c42e5aac5ba805825da76410c181273ba90b1`. Identisch.
- `astral-sh/setup-uv@bec219d24cd3e171d82865faccec33120bb574f4 # v10.1.0`
  — **dritte unabhängige Prüfung dieser SHA nach Implementer und
  Reviewer** (vierte Messung insgesamt inklusive dieser Sitzung). Selbst
  ausgeführt: `git ls-remote https://github.com/astral-sh/setup-uv
  refs/tags/v10.1.0` → `bec219d24cd3e171d82865faccec33120bb574f4`.
  Identisch mit dem im Workflow gepinnten SHA und Tag-Kommentar.

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.3 YAML-Struktur

Eigener Lauf: `ruby -ryaml -e "YAML.load_file('.github/workflows/sdk-python-release.yml')"`
→ kein Parse-Fehler. Zusätzlich alle vier `run:`-Blöcke einzeln über
einen eigenen Ruby-Extraktionslauf isoliert und jeweils mit `bash -n`
geprüft: alle vier syntaktisch fehlerfrei (Reviewer-Angabe selbst
nachvollzogen, nicht bloß übernommen).

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.4 `harness/README.md` §Werkzeuge — reale Zeile

Eigene Lektüre: Zeile 147 trägt die vollständige
`sdk-python-release.yml`-Zeile — Trigger, Isolationsbegründung,
Ablaufschritte (PEP-440-/`pyproject.toml`-Abgleich, `make
sdk-pack-python`, `uv publish` mit `UV_PUBLISH_TOKEN`), Secret-Referenz-
ohne-Anlage-Hinweis samt real geprüftem Nichtvorhandensein, `AGENTS.md`
§3.10-Verweis, Endung `kein Gate, [ADR-0107] Festlegung 5, [ADR-0108]
§Entscheidung Festlegung 1 · seit slice-sdk-python-publish-workflow` —
formkonsistent mit der Nachbarzeile `sdk-csharp-release.yml`. Referenzierte
Targets (`make test-sdk-python-release-tag-info`, `make
sdk-pack-python`) existieren real im `Makefile` (eigener `grep`, siehe
§1.6) — `AGENTS.md` §4 eingehalten, kein halluziniertes Target.

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.5 `make gates` grün

Eigener, ungepiped Lauf: `make gates > /tmp/gates-verifier2.log 2>&1; ec=$?`
→ `GATES_EXIT=0`. Log-Auszug:

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien` |
| `docs-check` | `d-check: 837 Datei(en) geprüft, 0 Befund(e)` |
| `commit-traceability` | `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators` |
| `a-check` | `gesamt: 0 Befund(e)` |

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.6 Review durchgeführt, Report vorhanden, 0 HIGH/MEDIUM/LOW/INFO

Review-Report vollständig gelesen: 0 HIGH/MEDIUM/LOW/INFO-Findings, keine
Fixrunde nötig, DoD-Checkbox-Nachzug ohne Fixrunde korrekt im selben
Commit vollzogen (`6c7e58f8`). Eigene Prüfungen bestätigen die im Report
behaupteten Läufe, statt sie zu übernehmen (siehe §1.1–§1.3, §1.7–§1.9,
§2).

**Permissions-Prüfung unabhängig wiederholt** (dritte unabhängige
Instanz nach Implementer und Reviewer): Top-Level `permissions: {}`
(Zeile 67) entzieht global alles; der einzige Job
`sdk-python-release` trägt **zusätzlich ein eigenes**
`permissions: { contents: read }`-Feld (Zeilen 74–75) — er erbt **nicht**
stillschweigend die leeren Top-Level-Rechte. Diese Lockerung ist auch
tatsächlich nötig: `actions/checkout` (Zeile 78) braucht `contents:
read`. `uv publish` (Zeile 99–102) authentifiziert ausschließlich über
`PYPI_API_TOKEN` (Secret), nicht über `GITHUB_TOKEN` — kein weiteres
GitHub-Permission-Feld nötig. Bestätigt: derselbe Fehlermodus wie beim
historischen `hub-description.yml`-`startup_failure`-Vorfall (fehlende
Job-Permissions unter einem restriktiven Top-Level-Feld) tritt hier
**nicht** auf.

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.7 Doku-Update — entfällt als eigener Punkt (bereits 1.4)

Korrekt als Duplikat der bereits oben geführten Zeile markiert.

### 1.8 Closure-Notiz

DoD-Zeile korrekt `[ ]` — Closure-Notiz mit Steering-Loop-Lerneintrag ist
ausdrücklich der Closure-Rolle vorbehalten, nicht diesem Verifikations-
Auftrag („Was du NICHT tust: … keine Closure-Notiz"). §7 des Slice-Plans
trägt bereits den Träger-Nachzug-Suchlauf (`AGENTS.md` §3.13) und den
Beobachtungs-Register-Check; kein Widerspruch zu §6/§1 gefunden.

### 1.9 Reconciliation entfällt

Bestätigt: keine Reconciliation-Datei in diesem Repo (wie bei allen
Vorgänger-Verifikationen dieser und der C#-Welle) — Entfall korrekt.

### 1.10 Drei Paarungen — bei Welle-Closure

DoD-Zeile korrekt `[ ]` — die Prüfung ist ausdrücklich der Welle-Closure
`welle-sdk-python-lh-fa-sst-009` vorbehalten, nicht diesem Slice. Dieser
Slice ist zugleich der letzte Anlass, die Welle-Closure anzustoßen —
außerhalb des heutigen Verifikations-Auftrags.

### 1.11 §6-Risiken tragen einen Ausgang — Checkbox-Semantik geprüft, alle drei Achsen selbst nachgemessen (dritte unabhängige Instanz)

Die DoD-Zeile bestätigt, dass zu **jedem** der vier §6-Risiken (Post-Push,
`PYPI_API_TOKEN`-Fehlen, PyPI-Namensraum, PEP-440-vs-SemVer-Regex-Teilung)
ein dokumentierter **Ausgang** existiert — für zwei davon lautet der
Ausgang „weiter offen, strukturell", für die anderen beiden „eingetreten,
real bestätigt" bzw. „entfallen"/„aufgelöst". Eigene Reverifikation aller
drei live prüfbaren Achsen:

- **`gh secret list --repo pt9912/pg-change-feed`** → `DOCKERHUB_TOKEN`,
  `DOCKERHUB_USERNAME`, `NUGET_API_KEY` — `PYPI_API_TOKEN` fehlt real,
  wie im Plan/§6 und im Review behauptet. **Dritte unabhängige
  Bestätigung** (nach Implementer und Reviewer).
- **`curl -s -o /dev/null -w '%{http_code}' https://pypi.org/pypi/pgchangefeed/json`**
  → `404` — Paketname real frei. **Dritte unabhängige Bestätigung.**
- **PEP-440-vs-SemVer-Regex-Isolation** — eigenständige Regex in
  `tools/harness/sdk-python-release-tag-info.sh` (Kopf-Kommentar
  vollständig gelesen, Begründung nachvollzogen); alle drei
  `run-*-release-tag-info-tests.sh`-Skripte real selbst ausgeführt (§2
  unten) — **vierte unabhängige Bestätigung**, dass keines der drei
  Skripte ein anderes gebrochen hat.

An keiner Stelle im Slice-Plan, im Review-Report oder in
`docs/user/releasing.md` wird behauptet, der reale Post-Push-Lauf sei
bereits erfolgt oder das Secret existiere — beide Aussagen bleiben
durchgängig als „weiter offen, strukturell" bzw. „fehlt real" markiert.

**Ergebnis: Checkbox-Semantik korrekt, kein fälschliches `[x]` für einen
tatsächlich noch bestehenden Blocker.** Der Slice liegt weiterhin unter
`in-progress/`, nicht `done/`.

## 2. Vierte unabhängige Test-Reverifikation — alle drei Skripte real selbst ausgeführt

```
$ bash tools/harness/run-sdk-python-release-tag-info-tests.sh
run-sdk-python-release-tag-info-tests: alle Fälle bestanden  (exit=0)

$ bash tools/harness/run-release-tag-info-tests.sh
run-release-tag-info-tests: alle Fälle bestanden  (exit=0)

$ bash tools/harness/run-sdk-csharp-release-tag-info-tests.sh
run-sdk-csharp-release-tag-info-tests: alle Fälle bestanden  (exit=0)
```

Alle drei Skripte liefen in derselben Sitzung, nacheinander, ohne
gegenseitige Beeinflussung — bestätigt, dass die neue, eigenständige
PEP-440-Regex weder das bestehende `tools/harness/semver-regex.sh` noch
dessen beide bisherigen Konsumenten (`release-tag-info.sh`,
`sdk-csharp-release-tag-info.sh`) verändert oder gebrochen hat. Dies ist
die **vierte** unabhängige Bestätigung dieser Isolation (Implementer,
Reviewer, diese Sitzung zweifach über alle drei Skripte).

## 3. Commit-Traceability und Immutabilität über den exakten Slice-Bereich

```
$ make doc-commits RANGE=3043c76a..HEAD
d-check: 837 Datei(en) geprüft, 0 Befund(e)

$ make doc-immutable RANGE=3043c76a..HEAD
d-check: 837 Datei(en) geprüft, 0 Befund(e)
```

Beide Läufe grün, `EXIT=0`. Zusätzlich eigene Prüfung aller sechs
Commit-Betreffe (`git log --format='%H %s' 3043c76a..HEAD`): alle sechs
tragen `LH-FA-SST-009` und mindestens eine `ADR-*`-Kennung
(`ADR-0107`/`ADR-0108`/`ADR-0106`), kein `SPEC-*`/`ARC-*` im Betreff
irgendeines Commits.

## 4. Live-Prüfung gegen GitHub/PyPI/NuGet (über die reine Repo-Textprüfung hinaus)

- **`gh secret list --repo pt9912/pg-change-feed`** → `DOCKERHUB_TOKEN`,
  `DOCKERHUB_USERNAME`, `NUGET_API_KEY` real vorhanden; `PYPI_API_TOKEN`
  fehlt real (§1.11).
- **`curl -s -o /dev/null -w '%{http_code}' https://pypi.org/pypi/pgchangefeed/json`**
  → `404` (§1.11).
- **`git tag -l "sdk-csharp-v*"`** → `sdk-csharp-v0.1.0` real vorhanden —
  bestätigt die Kernaussage des Coordinator-Fix-Commits `38a137f7`.
  **Dritte unabhängige Bestätigung** dieses Tags (nach Implementer/
  Reviewer der C#-Welle bzw. dem Reviewer dieses Slice).
- **`curl -s https://api.nuget.org/v3-flatcontainer/pgchangefeed.client/index.json`**
  → `{"versions": ["0.1.0"]}` — `PgChangeFeed.Client` 0.1.0 real auf
  NuGet.org veröffentlicht. **Dritte unabhängige Bestätigung**, dass der
  Coordinator-Fix-Text zum C#-SDK-Release-Stand in
  `docs/user/releasing.md` inhaltlich zutreffend ist.
- **Kein Tag wurde durch diese Verifikationssitzung gepusht, kein Secret
  angelegt** — bestätigt durch die eigene Befehlshistorie dieser Sitzung
  (ausschließlich Lese-/Prüfbefehle: `git ls-remote`, `git tag -l`,
  `gh secret list`, `curl`, `git status`, `make gates`,
  `make doc-commits`/`doc-immutable`).

## 5. ADR-0107/ADR-0108-Konformität und Scope-Grenze

- **Alle fünf Festlegungen aus `ADR-0107`** (Umfang v1, Vertriebsweg
  PyPI, Ort `sdks/python/`, Versionierung PEP 440, Build-/Publish-
  Mechanismus-Form) bleiben von `ADR-0108` unverändert bestätigt bis auf
  die eine benannte Klausel (Frontend `build`+`twine` → `uv`); der
  Workflow setzt exakt `ADR-0108`s Festlegung 1 um (`uv build`/`uv
  publish`, `UV_PUBLISH_TOKEN`), nicht `ADR-0107`s ursprüngliche
  `twine`-Zeile.
- **Scope-Grenze real geprüft:** eigener
  `git diff 3043c76a..HEAD --stat -- sdks/python .a-check.yml
  spec/architecture.md docs/user/version.md` → **keine Ausgabe** — alle
  vier genannten Pfade bleiben in diesem Diff unberührt. Der SDK-Code
  selbst (`sdks/python/pgchangefeed/**`, `sdks/python/Dockerfile`) stammt
  vollständig aus vorausgehenden, bereits `done/`-Slices; dieser Slice
  ändert an ihnen nichts.
- **Vollständiger Diff dieses Slice** (`git diff 3043c76a..HEAD --stat`):
  acht Dateien — Workflow, `Makefile`, Slice-Plan, Review-Report,
  `docs/user/releasing.md`, `harness/README.md`, das neue Tag-Info-Skript
  und sein Testskript. Kein unerwarteter Pfad.

## 6. Was diese Sitzung bewusst nicht erneut ausgeführt hat

- **`make sdk-pack-python`** — bereits im vorausgehenden Slice
  `slice-sdk-python-pack-werkzeug` real verifiziert; dieser Slice **ruft**
  das Werkzeug nur aus dem Workflow heraus auf, ändert sein Verhalten
  nicht. Kein zusätzlicher Erkenntniswert für die DoD *dieses* Slice.
- **Ein realer `sdk-python-v*`-Tag-Push** — ausdrücklich außerhalb des
  Verifier-Auftrags. Das zugehörige §6-Risiko bleibt nach `AGENTS.md`
  §3.10 strukturell offen.

## Verdikt

**DoD erfüllt** (für den aktuellen `in-progress`-Stand des Slice — die
Closure-Pflichten selbst, insbesondere die drei Paarungen bei
Welle-Closure, sind bewusst noch offen, siehe §1.10). Alle beauftragten
Prüfpunkte wurden real und unabhängig nachgemessen, keine Behauptung ohne
eigenen Beleg übernommen:

1. Trigger ausschließlich `push: tags: ['sdk-python-v*']`, kollisionsfrei
   mit `ci.yml`/`e2e.yml`/`release.yml`/`sdk-csharp-release.yml` —
   selbst gelesen und nachvollzogen.
2. PEP-440-/`pyproject.toml`-Abgleich läuft vollständig vor jedem
   Login/Build/Push — an den Zeilennummern selbst nachvollzogen.
3. Genau zwei `uses:`-Zeilen, beide SHA-gepinnt; beide SHAs selbst gegen
   die jeweilige Quelle verifiziert (`astral-sh/setup-uv`: dritte
   unabhängige Prüfung nach Implementer und Reviewer).
4. YAML-Struktur selbst geparst, fehlerfrei; alle vier `run:`-Blöcke
   isoliert `bash -n`-geprüft.
5. `harness/README.md` §Werkzeuge trägt die reale, formkonsistente
   Zeile, referenzierte Targets existieren real.
6. `make gates` real ausgeführt, ungepiped `GATES_EXIT=0`, alle sechs
   Gates grün.
7. Review durchgeführt, 0 HIGH/MEDIUM/LOW/INFO; Permissions-Kernpunkt
   (`contents: read` auf Job-Ebene trotz `permissions: {}` global) selbst
   nachvollzogen — dritte unabhängige Instanz.
8. Alle drei `run-*-release-tag-info-tests.sh`-Skripte real selbst
   ausgeführt: vierte unabhängige Bestätigung, dass die eigenständige
   PEP-440-Regex nichts gebrochen hat — alle drei `exit=0`.
9. §6-Risiken tragen korrekt einen Ausgang — die `[x]`-Checkbox bestätigt
   vollständige Buchführung, **nicht** Erledigung; `gh secret list` und
   der PyPI-404-Check wurden als dritte unabhängige Instanz real
   wiederholt, mit identischem Ergebnis.
10. Reconciliation entfällt korrekt (keine Reconciliation-Datei im Repo).
11. **Live-Prüfung:** `PYPI_API_TOKEN` fehlt real; `pgchangefeed` ist auf
    PyPI real frei; `sdk-csharp-v0.1.0` existiert real als Git-Tag,
    `PgChangeFeed.Client` 0.1.0 ist real auf NuGet.org veröffentlicht
    (Coordinator-Fix-Commit `38a137f7` inhaltlich korrekt, dritte
    unabhängige Bestätigung). Kein Tag/Secret wurde durch diese Sitzung
    angelegt.
12. Commit-Traceability und Immutabilität über den exakten Slice-Bereich
    (`3043c76a..HEAD`) beide grün, alle sechs Commits tragen
    `LH-FA-SST-009` und mindestens eine `ADR-*`-Kennung.
13. Scope-Grenze real geprüft: `sdks/python/**`, `.a-check.yml`,
    `spec/architecture.md`, `docs/user/version.md` bleiben in diesem Diff
    unberührt.

**Freigabe an den Planner:** Der Slice ist DoD-konform für seinen
aktuellen `in-progress`-Stand. Vor Closure der Welle
`welle-sdk-python-lh-fa-sst-009` bleiben die im Plan selbst bereits als
offen geführten Punkte zu erfüllen: die drei Paarungen (§1.10), die
Closure-Notiz mit Steering-Loop-Lerneintrag (§1.8) — und, strukturell
unabhängig von jeder Slice-/Welle-Closure, das Post-Push-Risiko sowie der
reale Nachweis eines erfolgreichen `sdk-python-v*`-Tag-Publish, die nach
`AGENTS.md` §3.10 bis zum ersten echten Tag-Push und der Anlage von
`PYPI_API_TOKEN` offen bleiben.
