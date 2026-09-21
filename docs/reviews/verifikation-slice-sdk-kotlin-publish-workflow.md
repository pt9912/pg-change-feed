# Verifikationsbericht: slice-sdk-kotlin-publish-workflow — 2026-09-21

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/slice-sdk-kotlin-publish-workflow.md` §2)
und die §6-Risiko-Ausgänge, in frischem Kontext. **Nicht** gegen den Diff
als solchen (Reviewer-Aufgabe, zweifach abgeschlossen mit
[`review-slice-sdk-kotlin-publish-workflow.md`](review-slice-sdk-kotlin-publish-workflow.md)
und
[`review-slice-sdk-kotlin-publish-workflow-fixrunde.md`](review-slice-sdk-kotlin-publish-workflow-fixrunde.md))
und **nicht** gegen realen Bedarf (Validator, hier nicht ausgelöst). Dies
ist der **letzte** Slice der Welle
[`welle-sdk-kotlin-lh-fa-sst-009`](../plan/planning/welle-sdk-kotlin-lh-fa-sst-009.md).

**Gegenstand:** vier Commits auf `main` (`LH-FA-SST-009`, `ADR-0109`):

- `546a7503` — ursprünglicher Implementer-Commit (`.github/workflows/sdk-kotlin-release.yml`,
  `build.gradle.kts`, Tag-Info-Werkzeug, Doku-Nachzug).
- `24077306` — erster Review-Report: 1 HIGH (F-1 — `./gradlew publish`
  lief auf dem Runner statt Docker-only, gegen `ADR-0109` §Entscheidung
  Festlegung 5 wörtlich).
- `21872c3d` — Fixrunde: neue Docker-Stufe `publish` in
  `sdks/kotlin/Dockerfile`, Publish-Schritt jetzt Docker-only.
- `2af9db25` — Fixrunden-Review: 0 HIGH/MEDIUM, 1 INFO — F-1 real
  aufgelöst, per eigenem `docker inspect` bestätigt.

**Frischer Kontext:** Diese Sitzung hat `harness/README.md`,
`AGENTS.md`, `harness/conventions.md`, den vollständigen Slice-Plan
(§1–§8), `ADR-0109` §Entscheidung Festlegung 5 (wörtlich, `Accepted`) und
beide Review-Reports vollständig gelesen. Nichts aus Slice-Plan,
Commit-Message oder Review-Report wurde ungeprüft übernommen: eigener
`docker build --build-context proto=proto --target publish`-Lauf +
`docker inspect` gegen das Ergebnis-Image, eigener `grep`-Abgleich aller
sechs Workflow-Trigger-Muster, eigener Lauf von
`make test-sdk-kotlin-release-tag-info`, eigener ungepipter
`make gates`-Lauf mit direkter Exit-Code-Prüfung (`AGENTS.md` §3.9),
eigene isolierte `make doc-commits`/`make doc-immutable`-Läufe über den
exakten Vier-Commit-Bereich, eigene Backtick-Paritäts-Zählung, eigene
`git diff --stat`-Prüfung auf unbeabsichtigte Berührungen, eigene
`git tag`-Prüfung.

---

## 1. DoD-Zeilen einzeln gegen den Stand geprüft

| DoD-Zeile (§2) | Prüfung | Ergebnis |
|---|---|---|
| Workflow existiert, Trigger `sdk-kotlin-v*`, kein Doppellauf | `grep -n "tags"` über alle sechs `.github/workflows/*.yml`: `ci.yml`/`e2e.yml`/`examples.yml` je `tags-ignore: ['**']`; `release.yml`/`sdk-csharp-release.yml`/`sdk-python-release.yml` je eigenes `tags: [...]`-Muster (`v*`/`sdk-csharp-v*`/`sdk-python-v*`); nur `sdk-kotlin-release.yml` trägt `sdk-kotlin-v*` | **bestätigt** — kein Musterüberlapp |
| Tag-Validierung gegen SemVer 2.0, eigener Testlauf | `tools/harness/sdk-kotlin-release-tag-info.sh` existiert; `make test-sdk-kotlin-release-tag-info` selbst ausgeführt | **bestätigt** — `run-sdk-kotlin-release-tag-info-tests: alle Fälle bestanden`, Exit 0 |
| `permissions: contents: read`/`packages: write`, kein externes Secret | `grep -n "permissions"`/`"secrets\."` im Workflow | **bestätigt** — Top-Level `permissions: {}`, Job-Ebene `contents: read`/`packages: write`; einziger `secrets.*`-Verweis ist `secrets.GITHUB_TOKEN` (eingebaut) |
| Publish läuft Docker-only, vierter unabhängiger Build der `publish`-Stufe | Eigener `docker build --build-context proto=proto --target publish -f sdks/kotlin/Dockerfile sdks/kotlin` + `docker inspect` | **bestätigt** — Build lief durch (voller Cache-Hit auf `build`), `docker inspect`: `Cmd=["./gradlew","--no-daemon","publish"]`, `Entrypoint=["/__cacert_entrypoint.sh"]` (von der Basis `eclipse-temurin:21-jdk` vererbt, nicht von der Stufe selbst gesetzt), `WorkingDir=/src/pgchangefeed-kotlin` — deckungsgleich mit beiden Review-Reports, unabhängig nachvollzogen |
| Alle `uses:`-Zeilen SHA-gepinnt mit Tag-Kommentar | `grep -rn "uses:"` über alle Workflow-Dateien | **bestätigt** — einzige `uses:`-Zeile ist der repo-weit wiederverwendete `actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1` |
| `harness/README.md` §Werkzeuge, `docs/user/releasing.md` (Version 1.6) nachgezogen | Beide Dateien gelesen | **bestätigt** — `harness/README.md`-Zeile beschreibt korrekt die Docker-only-`publish`-Stufe (nicht mehr den Runner-Aufruf); `docs/user/releasing.md` §4 Punkt 4 ebenso, Versionshistorie 1.5→1.6 mit anker-tragender Zeile |
| `make gates` grün | Eigener, ungepipter Lauf | **bestätigt** — Exit 0 (siehe §3) |
| Review durchgeführt, kein offenes HIGH | Beide Reports gelesen | **bestätigt** — Lauf 1: 1 HIGH (F-1), Lauf 2 (Fixrunde): 0 HIGH/MEDIUM, 1 INFO |

Alle sieben geprüften DoD-Zeilen sind **im aktuellen Repo-Zustand real
erfüllt**, nicht nur behauptet.

## 2. Fußnoten-Punkt aus dem Fixrunden-Review — `publish` als letzte
Dockerfile-Stufe, Default-Target-Risiko

`sdks/kotlin/Dockerfile` trägt `publish` jetzt als letzte definierte
Stufe (`FROM build AS publish` nach `pack-export`). Ein
`docker build sdks/kotlin` **ohne** `--target` würde standardmäßig
`publish` bauen (Docker-Semantik: ohne `--target` wird die zuletzt
definierte Stufe gebaut).

Eigene Prüfung:

```
$ grep -rn "docker build.*sdks/kotlin" .
tools/harness/sdk-pack-kotlin.sh:43: … --target pack-export -t "$SDK_PACK_KOTLIN_IMAGE" sdks/kotlin
.github/workflows/sdk-kotlin-release.yml:120: … --target publish -t pg-change-feed:sdk-kotlin-publish sdks/kotlin
```

**Bestätigt:** Es existieren real genau zwei Aufrufstellen
(`tools/harness/sdk-pack-kotlin.sh`, `.github/workflows/sdk-kotlin-release.yml`),
beide setzen explizit `--target`. Kein dritter Aufrufer im Repo. Weder
`.a-check.yml` noch `.d-check.yml` referenzieren `sdks/kotlin/Dockerfile`.
Die Fußnote bleibt zu Recht eine Fußnote, kein Befund — der latente Pfad
(ein künftiger dritter Aufrufer ohne `--target`) ist real nicht vorhanden.

## 3. `make gates` — eigener, ungepipter Lauf

```
$ make gates > /tmp/verifier_gates_kotlin.log 2>&1; echo $?
0
```

Einzelbelege aus demselben Log:

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK` |
| `docs-check` | `d-check: 870 Datei(en) geprüft, 0 Befund(e)` |
| `commit-traceability` | `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators` |
| `a-check` | `gesamt: 0 Befund(e)` |

Alle sechs Gate-Ziele grün, Exit-Code direkt (nicht durch eine Pipe)
geprüft (`AGENTS.md` §3.9).

Isolierte Bestätigung über den exakten Vier-Commit-Slice-Bereich
(`546a7503^..2af9db25`):

```
$ make doc-commits RANGE=546a7503^..2af9db25
d-check: 870 Datei(en) geprüft, 0 Befund(e)   → Exit 0

$ make doc-immutable RANGE=546a7503^..2af9db25
d-check: 870 Datei(en) geprüft, 0 Befund(e)   → Exit 0
```

Kein Commit-Traceability-Verstoß, keine ADR-Immutabilitätsverletzung über
diesen Slice — alle vier Commit-Messages tragen `LH-FA-SST-009`/`ADR-0109`
im Betreff (eigene `git log`-Prüfung, keine `SPEC-*`/`ARC-*`-Kennung im
Betreff).

## 4. `git diff` gegen Elternstand (`eaf6ded5..HEAD`) — Out-of-Scope-Prüfung

```
$ git diff --stat eaf6ded5..HEAD -- examples/kotlin .a-check.yml \
    spec/architecture.md docs/user/version.md sdks/csharp sdks/python
(keine Ausgabe)

$ git diff --stat eaf6ded5..HEAD -- docs/plan/planning/done/slice-sdk-kotlin-projektgeruest.md \
    docs/plan/planning/done/slice-sdk-kotlin-pack-werkzeug.md \
    docs/plan/planning/done/slice-sdk-kotlin-http-client-flaeche.md \
    docs/plan/planning/done/slice-sdk-kotlin-grpc-client-flaeche.md
(keine Ausgabe)

$ git tag -l | grep -i kotlin
(keine Ausgabe)
```

**Bestätigt:** Keine unbeabsichtigte Berührung von `examples/kotlin/**`,
`.a-check.yml`, `spec/architecture.md`, `docs/user/version.md`,
`sdks/csharp/**`, `sdks/python/**` oder den vier bereits geschlossenen
Kotlin-Vorgänger-Slices. `git diff --stat eaf6ded5..HEAD` zeigt exakt 13
Dateien, alle innerhalb des geplanten §3-Umfangs (inkl. einer
planungs-eigenen, dieser Welle fremden Zwischen-Commit-Datei
`slice-generated-sync-tar-export.md`, `3f754976` — ein eigenständiger,
unabhängig committeter Slice-Plan außerhalb dieser Welle, kein
Scope-Verstoß dieses Slice). `git tag -l | grep kotlin` liefert **keinen**
Treffer — kein realer Tag wurde gesetzt.

## 5. Backtick-Parität — eigenständig nachgezählt

| Datei | Backticks | gerade? |
|---|---|---|
| `.github/workflows/sdk-kotlin-release.yml` | 118 | ja |
| Slice-Plan | 506 | ja |
| `docs/user/releasing.md` | 494 | ja |
| `harness/README.md` | 1796 | ja |
| `sdks/kotlin/Dockerfile` | 112 | ja |
| `build.gradle.kts` | 136 | ja |
| `docs/reviews/review-slice-sdk-kotlin-publish-workflow.md` | 354 | ja |
| `docs/reviews/review-slice-sdk-kotlin-publish-workflow-fixrunde.md` | 482 | ja |

Alle acht in diesem Slice geänderten Markdown-/Konfigurationsdateien
tragen eine gerade Backtick-Anzahl.

## 6. §6-Risiken — Ausgänge geprüft

| Risiko | Ausgang im Plan | Zulässige Klasse? | Eigene Einschätzung |
|---|---|---|---|
| `AGENTS.md` §3.10 — realer Post-Push-Lauf unverifiziert | „weiter offen, strukturell … bis zum ersten realen Tag-Push" | ✓ weiter offen | Korrekt geführt — kein Closure-Blocker, kein `make gates`-Ersatz behauptet; das INFO-Finding F-1 des Fixrunden-Reviews (Kommentar-Behauptung zur Gradle-Up-to-date-Prüfung) ist explizit als Detailschärfung dieses selben, bereits offenen Risikos eingeordnet, kein eigenständiges neues Risiko |
| Kein Secret-Anlage-Risiko | „entfällt strukturell" | ✓ entfallen | Begründung trägt — `GITHUB_TOKEN` ist immer verfügbar, real bestätigt (kein `secrets.<NAME>`-Verweis außer dem eingebauten) |
| SemVer-Tag-Parser teilt sich Code mit `semver-regex.sh` | „aufgelöst — dritter Konsument" | ✓ aufgelöst | Real bestätigt: `make test-sdk-kotlin-release-tag-info` grün, Kopf-Kommentar in `semver-regex.sh` auf drei Konsumenten erweitert |
| Reales `publishing.repositories.maven`-Minimalrezept könnte Gradle-Eigenheit zeigen | „teilweise geprüft, im Kern weiter offen" | ✓ weiter offen | Korrekt — `make sdk-pack-kotlin` bestätigt fehlerfreies Parsen, der reale Schreibzugriff gegen GitHub Packages bleibt nach `AGENTS.md` §3.10 bis zum ersten Tag-Push offen |

Alle vier Risiken tragen eine der drei zulässigen Ausgangsklassen
(*eingetreten/aufgelöst* · *entfallen* · *weiter offen*). Kein Risiko
wurde stillschweigend geschlossen, ohne dass sein Beleg dafür ausreicht.

## 7. Kein echter Tag, kein echter Push

- `git tag -l | grep -i kotlin` → **keine Ausgabe** (kein `sdk-kotlin-v*`-Tag
  im Repo).
- Der eigene Docker-Build/-Inspect-Lauf in §1/§2 dieses Berichts lief
  ausschließlich gegen den lokalen Layer-Cache, ohne `docker run … publish`
  auszuführen — kein Netzwerkkontakt zu `maven.pkg.github.com`, kein
  GitHub-Token verwendet.
- Kein `git push`, keine Änderung an `origin` durch diesen Verifikations-Lauf.

## Verdikt

**DoD konform: ja.** Alle sieben real geprüften DoD-Zeilen entsprechen
dem tatsächlichen Repo-Zustand, nicht nur der Implementer-/Reviewer-
Behauptung. Das HIGH-Finding F-1 des ersten Reviews (`./gradlew publish`
lief entgegen `ADR-0109` §Entscheidung Festlegung 5 wörtlich auf dem
Runner) ist real aufgelöst — eigener `docker build --target publish` +
`docker inspect` bestätigt `Cmd=["./gradlew","--no-daemon","publish"]`
in der neuen, gepinnten Docker-Stufe `publish`, kein Aufruf mehr auf dem
GitHub-hosted Runner. Die vom Fixrunden-Review genannte Fußnote (Default-
Target-Risiko der letzten Dockerfile-Stufe) ist real ungefährlich — beide
existierenden Aufrufstellen setzen explizit `--target`. `git diff` gegen
den Elternstand zeigt keine unbeabsichtigte Berührung von
`examples/kotlin/**`, `.a-check.yml`, `spec/architecture.md`,
`docs/user/version.md`, `sdks/csharp/**`, `sdks/python/**` oder den vier
bereits geschlossenen Kotlin-Vorgänger-Slices. Backtick-Parität aller acht
geänderten Dateien ist gegeben. Alle vier §6-Risiken tragen eine
zulässige Ausgangsklasse; das AGENTS.md §3.10-Risiko bleibt korrekt
strukturell **weiter offen**, kein Closure-Blocker. `make gates` lief
eigenständig und ungepiped grün (Exit 0), ebenso isoliert
`make doc-commits`/`make doc-immutable` über den exakten
Vier-Commit-Bereich. Kein echter `sdk-kotlin-v*`-Tag wurde gesetzt, kein
echter Push gegen `maven.pkg.github.com` erfolgte.

**Was dieser Bericht NICHT tut:** keine Fixes, kein `git mv` nach
`done/`, keine Closure-Notiz-Fertigstellung dieses Slice — die
verbleibenden `[ ]`-DoD-Zeilen (Closure-Notiz, Beobachtungs-Register,
Risiko-Ausgang-Formalisierung, drei Paarungen) sind bewusst dem
Rollenwechsel zum Planner/Architect vorbehalten, nicht Teil dieser
Verifikation.

**Übergabe:** Dieser Bericht geht an den Planner — DoD-Konformität für
`slice-sdk-kotlin-publish-workflow` ist bestätigt; der Planner kann die
verbleibenden Closure-Schritte (Closure-Notiz, Beobachtungs-Register,
§6-Formalisierung, drei Paarungen inkl. Welle-Closure für
`welle-sdk-kotlin-lh-fa-sst-009`) einleiten. Das AGENTS.md §3.10-Risiko
bleibt auch nach dieser Verifikation strukturell offen — es wird erst
durch den ersten realen `sdk-kotlin-v*`-Tag-Push aufgelöst oder mit rotem
Befund und Folgemaßnahme dokumentiert.
