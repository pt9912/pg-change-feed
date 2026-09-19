# Review-Report: slice-sdk-python-publish-workflow — 2026-09-19

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-python-publish-workflow.md`),
`ADR-0107`/`ADR-0108` (beide Accepted) und `AGENTS.md` Hard Rules (Modul
10 §Drei Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe
(Modul 11).

**Gegenstand:** Diff-Range `3043c76a..HEAD` (Abschluss von
`slice-sdk-python-pack-werkzeug` bis Kopf), zwei Commits: `3808be6c`
(Implementer-Zug dieses Slice — plus drei vorausgehende reine
Lifecycle-/Feld-Commits `90f2c9b8`/`3fcc941a`/`aad9ee26`, alle real als
0/0- bzw. 1-Zeilen-Diffs geprüft, `AGENTS.md` §3.3 sauber) und `38a137f7`
(separater Coordinator-Fix, außerhalb dieses Slice-Scopes — Fund beim
Gegenlesen, behebt eine veraltete Aussage zum C#-SDK-Release-Stand in
`docs/user/releasing.md`, hier nur auf Korrektheit geprüft, nicht als Teil
der Slice-DoD).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, seither um mehrere weitere HIGH-Klassen
ergänzt, u. a. `AGENTS.md` §3.9 Pipe-Disziplin, §3.10 Post-Push-Risiko,
§3.13 Träger-Nachzug).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-19.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-sdk-python-publish-workflow.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan inkl. „Beim Schreiben getroffene
  Entscheidungen", §6 Risiken, §7 Closure-Notiz, §8 Sub-Area/Modus)
- `docs/plan/adr/0107-python-pypi-zweites-sdk-package.md` (Accepted) —
  Festlegung 4 (Versionierung), Festlegung 5 (Trigger, Secret,
  Tag-Präfix, Publish-Muster)
- `docs/plan/adr/0108-python-sdk-uv-statt-build-twine.md` (Accepted,
  Supersedes `ADR-0107` in genau einer Klausel) — §Entscheidung
  Festlegung 1 (`uv build`/`uv publish`, `UV_PUBLISH_TOKEN`)
- `AGENTS.md` §3.1 (Docker-only), §3.3 (git mv + Inhalt = zwei Commits),
  §3.7 (Kommentar-Disziplin), §3.8 (Action-Pinning), §3.9
  (Pipe-/Exit-Code-Disziplin), §3.10 (Post-Push-Risiko GitHub Actions),
  §3.12 (Herkunft von Aussagen), §3.13 (Träger-Nachzug), §4 (kein Träger
  nennt ein nicht existentes Target)
- `harness/README.md` §Werkzeuge (Gate-Index, Vergleichszeilen
  `release.yml`/`sdk-csharp-release.yml`)
- `docs/user/releasing.md` (bestehender Release-Prozess-Träger, Vorbild
  für Secret-Tabelle und Trigger-Beschreibung; hier zusätzlich Gegenstand
  des separaten Coordinator-Fix-Commits)
- `docs/reviews/review-slice-sdk-csharp-publish-workflow.md` —
  Formvorbild für diesen Report

**Eigenständig durchgeführte Prüfungen (nicht nur Commit-Message/Bericht
übernommen):**

- `.github/workflows/sdk-python-release.yml` vollständig gelesen:
  Trigger ausschließlich `push: tags: ['sdk-python-v*']`, kein
  Zweit-Trigger.
- **Permissions-Prüfung (gezielter Auftrag, bekannte historische
  Fehlerklasse):** Top-Level `permissions: {}` (Zeile 67), der einzige
  Job `sdk-python-release` trägt **zusätzlich** ein eigenes
  `permissions: { contents: read }`-Feld (Zeile 74f.) — er erbt nicht
  stillschweigend die leeren Top-Level-Rechte. `actions/checkout` (Zeile
  78) hat damit die nötige Berechtigung; `uv publish` authentifiziert
  über `PYPI_API_TOKEN`, kein GitHub-`GITHUB_TOKEN`-Pfad. Kein Befund.
- `grep -n "tags:"` real über alle fünf Workflow-Dateien
  (`ci.yml`/`e2e.yml`/`release.yml`/`sdk-csharp-release.yml`/
  `sdk-python-release.yml`) laufen lassen: `ci.yml`/`e2e.yml` tragen
  `tags-ignore: ['**']` (schließt **jeden** Tag-Push aus),
  `release.yml` trägt `tags: ['v*']`, `sdk-csharp-release.yml` trägt
  `tags: ['sdk-csharp-v*']`, `sdk-python-release.yml` trägt
  `tags: ['sdk-python-v*']`. Ein GitHub-Actions-`tags:`-Glob-Muster
  verlangt eine exakte literale Präfix-Übereinstimmung vor dem `*` —
  `sdk-python-v*` matcht keines der drei anderen Muster (beginnt weder
  mit `v` noch mit `sdk-csharp-v`), `tags-ignore: ['**']` schließt es
  ebenfalls aus. Kein Doppellauf über alle fünf Dateien hinweg.
- **PEP-440-/`pyproject.toml`-Abgleich-Reihenfolge real an den
  Zeilennummern geprüft:** Checkout (Z. 77f.) → Tag validieren/Version
  ermitteln (Z. 80f.) → `pyproject.toml`-Version abgleichen, Abbruch bei
  Abweichung (Z. 83–89) → `make sdk-pack-python` (Z. 91f.) → `uv`
  installieren (Z. 94–97) → `uv publish` (Z. 99–102) — beide
  Validierungsschritte laufen **vor** jedem Login/Build/Push, exakt wie
  im Plan zugesagt.
- **Eigenständige PEP-440-Regex-Entscheidung nachvollzogen:** Der
  Kopf-Kommentar von `tools/harness/sdk-python-release-tag-info.sh`
  begründet konkret und geprüft (nicht behauptet), warum kein Sourcing
  aus `tools/harness/semver-regex.sh` erfolgt — dessen Regex akzeptiert
  SemVer-Pre-Release-/Build-Metadata-Syntax, die kein gültiges PEP 440
  ist; ein Sourcing würde die Regex als „PEP 440" beschriften, obwohl sie
  SemVer-Grammatik prüft. Nachvollziehbar, real durch die Testfälle
  `assert_invalid "sdk-python-v1.0.0-alpha.1"`/
  `assert_invalid "sdk-python-v1.0.0+build.5"` belegt (genau die Syntax,
  die eine geteilte SemVer-Regex fälschlich akzeptiert hätte). `bash
  tools/harness/run-sdk-python-release-tag-info-tests.sh` real
  ausgeführt: „alle Fälle bestanden", Exit `0`, alle 18 Fälle (3 gültig,
  4 Präfix-Negativfälle, 11 PEP-440-Grammatik-Negativfälle) real gezählt.
- **`astral-sh/setup-uv`-SHA real gegen die Quelle verifiziert:**
  `git ls-remote https://github.com/astral-sh/setup-uv refs/tags/v10.1.0`
  liefert `bec219d24cd3e171d82865faccec33120bb574f4` — identisch mit dem
  im Workflow gepinnten SHA und Tag-Kommentar (Zeile 95). `actions/checkout`
  wiederverwendet den bereits repo-weit gepinnten SHA (`3d3c42e5aac5…`, `v7.0.1`).
- `ruby -ryaml -e "YAML.load_file('.github/workflows/sdk-python-release.yml')"`
  real ausgeführt: kein Parse-Fehler. Zusätzlich alle vier `run:`-Blöcke
  einzeln extrahiert und mit `bash -n` geprüft: alle vier syntaktisch
  fehlerfrei.
- `grep -rn "PYPI_API_TOKEN\|NUGET_API_KEY"` über `.yml`/`.sh`/`.md` real
  ausgeführt: nur Referenzen (`${{ secrets.PYPI_API_TOKEN }}`,
  `UV_PUBLISH_TOKEN`, Prosa-Erwähnungen in Plan/ADR/README) — keine
  hartkodierte Zeichenkette, keine Anlage.
- `gh secret list --repo pt9912/pg-change-feed` real ausgeführt: zeigt
  `DOCKERHUB_TOKEN`, `DOCKERHUB_USERNAME`, `NUGET_API_KEY` —
  `PYPI_API_TOKEN` fehlt real, wie im Plan/§6 behauptet.
- `curl -s -o /dev/null -w '%{http_code}' https://pypi.org/pypi/pgchangefeed/json`
  real ausgeführt: `404` — Paketname real frei, wie im Plan/§6 behauptet.
- `git tag -l "sdk-csharp-v*"` real ausgeführt: `sdk-csharp-v0.1.0`
  existiert real — trägt die Kernaussage des Coordinator-Fix-Commits.
  Zusätzlich `curl https://api.nuget.org/v3-flatcontainer/pgchangefeed.client/index.json`
  real ausgeführt: `{"versions": ["0.1.0"]}` — das Package ist real auf
  NuGet.org veröffentlicht, bestätigt auch die im Fix-Text zitierte
  Aussage „`dotnet nuget push` bestätigte 201 Created" indirekt (Existenz
  auf der Registry), nicht nur den Tag selbst. **Coordinator-Fix
  inhaltlich korrekt.**
- `harness/README.md` §Werkzeuge neue Zeile gegen `sdk-csharp-release.yml`-
  Zeile verglichen: gleiche Form (Trigger, Ablauf, Secret-Referenz-ohne-
  Anlage-Hinweis, `AGENTS.md` §3.10-Verweis, `kein Gate, [ADR-Link] · seit
  slice-<name>`); referenziertes Target (`make
  test-sdk-python-release-tag-info`) existiert real im `Makefile`
  (`AGENTS.md` §4); referenziertes `make sdk-pack-python` existiert real
  in `harness/mk/sdk.mk` (aus dem vorausgehenden, bereits `done` Slice).
- `docs/user/releasing.md` vollständig gelesen (nicht nur den Diff-Hunk):
  neuer Abschnitt „SDK-Release: PyPI-Publish für `pgchangefeed`" trägt
  Tag-Namensraum, vierstufigen Ablauf (Tag validieren → `pyproject.toml`
  abgleichen → `make sdk-pack-python` → `uv publish`) und den
  `PYPI_API_TOKEN`-Zeileneintrag in der Secret-Tabelle — Ablauf-
  Beschreibung gegen die reale Workflow-Datei gehalten, kein Drift.
  Versionskopf `1.2` → `1.3`, neue Zeile in der Änderungshistorie-Tabelle
  mit Bezug auf diesen Slice/`LH-FA-SST-009`/`ADR-0107`/`ADR-0108` — vorab
  eingeplant, nicht erst als Reviewer-Nachzug (Lehre aus der C#-Welle
  korrekt angewendet).
- Kommentar-Disziplin (`AGENTS.md` §3.7) in allen neuen/geänderten
  Dateien gelesen: Workflow-Kopfkommentar,
  `sdk-python-release-tag-info.sh`,
  `run-sdk-python-release-tag-info-tests.sh` — durchweg indikativ über
  den Ist-Zustand, ADR-/`slice-`-Herkunfts-Anker in der akzeptierten
  Form, kein Konjunktiv über eine verworfene Alternative, kein
  abwesender Text, kein Satzabbruch.
- `make gates` real ausgeführt, Exit-Code direkt (nicht durch Pipe)
  geprüft: `0` (`generated-sync: OK`, `a-check: gesamt: 0 Befund(e)` —
  vollständiger Log unter dem Lauf dieses Reviews).
- Import-/Scope-Grenze geprüft: `sdks/python/Dockerfile` und
  `tools/harness/sdk-pack-python.sh` (von `ADR-0108` referenziert) sind
  real **nicht** Teil dieses Diffs — bereits im vorausgehenden Slice
  `slice-sdk-python-pack-werkzeug` (Commit `0e42b801`) gelandet; dieser
  Slice ändert an ihnen nichts.

---

## Findings

Keine HIGH-, MEDIUM- oder LOW-Findings in diesem Lauf.

## Negativbefunde

- geprüft, ohne Befund: Permissions — der Job trägt ein eigenes,
  ausreichendes `permissions: { contents: read }`-Feld, erbt nicht
  stillschweigend das leere Top-Level `permissions: {}` — dieselbe
  Lektion aus dem historischen `hub-description.yml`-Vorfall korrekt
  angewendet.
- geprüft, ohne Befund: Trigger-Isolation — `sdk-python-v*` überlappt mit
  keinem der drei Muster in `ci.yml`/`e2e.yml`/`release.yml`/
  `sdk-csharp-release.yml` (real per `grep` über alle fünf Dateien
  bestätigt, kein impliziter Trugschluss).
- geprüft, ohne Befund: Validierungsreihenfolge — PEP-440- und
  `pyproject.toml`-Abgleich laufen beide vor jedem Login/Build/Push, per
  Zeilennummer im YAML bestätigt, nicht nur per Existenz der Schritte.
- geprüft, ohne Befund: eigenständige PEP-440-Regex-Entscheidung —
  nachvollziehbar begründet und durch die 18 Testfälle (insbesondere die
  SemVer-only-Negativfälle) real belegt, kein Sourcing-Verzicht ohne
  Substanz.
- geprüft, ohne Befund: Action-Pinning — `astral-sh/setup-uv`-SHA real
  gegen `git ls-remote` verifiziert (Tag `v10.1.0` zeigt exakt auf den im
  Kommentar genannten SHA); `actions/checkout` wiederverwendet den
  bereits repo-weit gepinnten SHA.
- geprüft, ohne Befund: `PYPI_API_TOKEN`/`NUGET_API_KEY` — ausschließlich
  referenziert, nirgends hartkodiert oder angelegt.
- geprüft, ohne Befund: YAML-Struktur — Ruby-Stdlib-Parser lädt die
  Datei fehlerfrei; alle vier `run:`-Blöcke einzeln `bash -n`-sauber.
- geprüft, ohne Befund: `harness/README.md` §Werkzeuge — neue Zeile
  formkonsistent mit `sdk-csharp-release.yml`-Zeile, referenzierte
  Targets (`make test-sdk-python-release-tag-info`, `make
  sdk-pack-python`) existieren real.
- geprüft, ohne Befund: `docs/user/releasing.md` — neuer
  Python-Abschnitt vorab eingeplant (kein Reviewer-Nachzug nötig),
  Ablauf-Beschreibung deckungsgleich mit der realen Workflow-Datei,
  Versionshistorie korrekt fortgeschrieben.
- geprüft, ohne Befund: Coordinator-Fix-Commit `38a137f7` (C#-Absatz) —
  `sdk-csharp-v0.1.0` real als Git-Tag vorhanden, `PgChangeFeed.Client`
  0.1.0 real auf NuGet.org veröffentlicht; die korrigierte Aussage ist
  zutreffend.
- geprüft, ohne Befund: §6-Risiken — alle vier benannten Risiken
  (Post-Push, fehlendes Secret, PyPI-Namensraum, PEP-440-Regex-Teilung)
  tragen einen korrekten, real nachgeprüften Ausgang; keines fälschlich
  als erledigt markiert, wo es strukturell offen bleiben muss
  (Post-Push/Secret-Fehlen bleiben „weiter offen").
- geprüft, ohne Befund: `AGENTS.md` §3.3 — alle drei Lifecycle-/Feld-
  Commits (`90f2c9b8`, `3fcc941a`, `aad9ee26`) sind reine Moves bzw.
  Ein-Zeilen-Feldänderungen, getrennt vom Inhalts-Commit `3808be6c`.
- geprüft, ohne Befund: `AGENTS.md` §3.7 — Kommentare in Workflow-Kopf
  und beiden neuen Shell-Skripten sind indikativ über den Ist-Zustand,
  keine Chronik, kein Konjunktiv über eine verworfene Alternative.
- geprüft, ohne Befund: Traceability — beide Commits im Slice-Scope
  nennen `LH-FA-SST-009` und mindestens eine `ADR-*`-Kennung im Betreff,
  kein `SPEC-*`/`ARC-*` im Betreff.
- geprüft, ohne Befund: `make gates` real gefahren, Exit-Code direkt
  (ungepiped) `0`.
- geprüft, ohne Befund: Scope-Grenze — `sdks/python/Dockerfile` und
  `tools/harness/sdk-pack-python.sh` sind real nicht Teil dieses Diffs,
  keine stille Vermischung mit dem Vorgänger-Slice.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** keine.

## Verdikt

**Merge-blockierend:** nein — 0 HIGH/MEDIUM/LOW-Findings.

**DoD-Checkbox-Nachzug ohne Fixrunde** greift: Da dieses Verdikt keine
Fixrunde am Implementer auslöst, wird die DoD-Zeile „Review durchgeführt,
Report unter `docs/reviews/` liegt vor" im Slice-Plan
(`docs/plan/planning/in-progress/slice-sdk-python-publish-workflow.md`)
im selben Commit, der diesen Report anlegt, auf `[x]` nachgezogen, mit
Verweis auf diesen Report (`.harness/skills/reviewer.md` §DoD-Checkbox-
Nachzug ohne Fixrunde).

**Übergabe:** Kein Reviewer→Implementer-Rückgabe-Pfeil nötig. Der
Coordinator-Fix-Commit `38a137f7` wurde ausschließlich auf Korrektheit
geprüft (nicht als Teil dieser Slice-DoD) und ist inhaltlich zutreffend.
Dieser Report ist ein Lauf-Beleg; er ersetzt keine Verifikation gegen die
volle DoD — das bleibt Verifier-Aufgabe (Modul 11). Das Post-Push-Risiko
(realer Tag-Push gegen PyPI) bleibt nach `AGENTS.md` §3.10 strukturell
offen und ist damit **kein** Review-Finding, sondern korrekt als
weiter offenes §6-Risiko geführt.
