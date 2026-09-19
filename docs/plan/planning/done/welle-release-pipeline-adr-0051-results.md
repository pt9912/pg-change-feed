# Welle welle-release-pipeline-adr-0051 — Closure-Notiz

**Welle:** welle-release-pipeline-adr-0051
**Abschluss:** 2026-09-19
**Verantwortlich:** pt9912

## Was wurde geliefert?

`ADR-0051` (`Accepted`) vollständig umgesetzt — alle fünf noch fehlenden
Teilentscheidungen, in fünf Slices:

- **`release-version-und-workflow`**: `docs/user/version.md` (Quelle der
  Wahrheit für die Version), `make image` um `VERSION`/`LATEST` erweitert
  (ein Build, Content-Mirror-Push nach GHCR **und** Docker Hub, kein
  zweiter Build — real gegen zwei lokale Registries verifiziert:
  identischer Manifest-Digest über vier Tag/Registry-Kombinationen, ohne
  echte Zugangsdaten zu brauchen), `release.yml` (Trigger `push: tags:
  ['v*']`, Tag-Validierung über `tools/harness/release-tag-info.sh` +
  `make test-release-tag-info`, GitHub-Release mit Digest im
  Beschreibungstext).
- **`release-image-scan`**: `make image-cve` (digest-gepinntes
  `aquasec/trivy`) + `image-scan.yml`, advisory, nächtlich +
  `workflow_dispatch`. Löst die seit `slice-001` angekündigte,
  unimplementierte Zusage ein.
- **`release-upstream-drift`**: sieben neue `make pin-stale-*`-Targets für
  P3–P9 (`tools/harness/pin-stale.sh`, `pin-stale-dcheck.sh`,
  `pin-stale-baseline.sh`, `pin-stale-actions.sh`) + `upstream-drift.yml`,
  fail-open über alle neun Achsen. Real ausgeführt: P3
  (`golang:1.27`)/P4 (`postgres:18-alpine`)/P5
  (`ghcr.io/pt9912/d-migrate`) zeigen echten, vorher unbekannten
  Upstream-Drift — der erste reale Fund dieses Werkzeugs.
- **`release-hub-description`**: `hub-description.yml`
  (`workflow_call`/`workflow_dispatch`) synchronisiert `README.md` als
  Docker-Hub-Beschreibung, als zusätzlicher `needs: release`-Job in
  `release.yml`. Trägt den im Schwester-Repo d-check real dokumentierten
  Scope-Hinweis (`DOCKERHUB_TOKEN` braucht `read/write/delete`, nicht nur
  `read/write`).
- **`release-doku-releasing`**: `docs/user/releasing.md` — beschreibt den
  jetzt vollständig real existierenden Prozess für Betreiber/Maintainer;
  jede genannte Datei/jedes genannte Target real gegen den Baum geprüft.

`make gates` grün auf dem Endstand (778 Dateien, 0 `docs-check`-Befunde,
Coverage 82,80 % ≥ 80 %, `a-check`/`generated-sync`/`commit-traceability`
je ohne Befund). Kein realer Release-Tag wurde gesetzt (`git tag -l`
leer) — diese Welle liefert bewusst nur den **Mechanismus** (§6
Out-of-Scope der Welle-Datei), nicht seine erste Anwendung.

## Was hat funktioniert?

Reale, hands-on-Verifikation ohne echte Zugangsdaten war über alle fünf
Slices hinweg die tragende Technik und fand mehrfach echte Fehler, bevor
sie in den Code gelangten — nicht nur theoretisch plausibel gemacht:

- Zwei lokale `registry:2`-Instanzen bewiesen die Content-Mirror-
  Eigenschaft ohne GHCR-/Docker-Hub-Zugangsdaten.
- Ein realer Docker/Trivy-Regressionstest deckte auf, dass ein
  allgemeiner `~/.docker/config.json`-Credential-Mount Trivys eigenen,
  unauthentifizierten Vulnerability-DB-Bezug mitfärbt — die gewählte
  Alternative (Trivys eigene `--username`/`--password`-Flags) wurde
  dadurch zur begründeten, nicht nur plausiblen Entscheidung.
- Vier reale `curl`-Aufrufe gegen die echte, öffentliche Docker-Hub-API
  mit absichtlich ungültigen Zugangsdaten bestätigten Endpoint-URLs und
  Feldnamen (`identifier`/`secret`/`access_token`/`full_description`)
  und sogar, dass `pt9912/pg-change-feed` dort bereits existiert (401
  statt 404).
- Ein Nutzer-Hinweis auf das Schwester-Repo d-check
  (`packaging/dockerhub/README.md`) lieferte einen dort real
  dokumentierten Produktionsfehler (Token-Scope `read/write` reicht
  nicht für den Beschreibungs-`PATCH`) und ließ sich direkt übernehmen,
  bevor der Nutzer das eigene Token überhaupt angelegt hatte.

Strukturell etablierte sich ein wiederverwendbares Muster dreifach:
Inline-Shell-Logik in einem GitHub-Actions-Workflow bekommt ein
eigenständiges, netzlos testbares Skript statt unbelegter Inline-Logik
(`tools/harness/release-tag-info.sh`, `tools/harness/dockerhub-token.sh`,
`tools/harness/lib-github-api.sh` für containerisierte GitHub-API-
Abfragen aus lokalen `make`-Targets).

## Was ging anders als geplant?

Jeder der fünf Slices ging durch mindestens eine Reviewer-Fixrunde,
keiner blieb beim ersten Durchlauf unbeanstandet:

- `release-version-und-workflow`: 1 HIGH + 4 MEDIUM — darunter zwei
  echte SemVer-2.0-Bugs (führende Null im Prerelease-Identifier
  akzeptiert; Bindestrich in der Build-Metadata fälschlich als
  Prerelease gewertet) und eine bewusst offen gelassene
  `ADR-0103`-Governance-Frage (F-5, siehe Trigger-Audit unten).
- `release-image-scan`: 2 MEDIUM + 2 LOW — fehlende GHCR-Authentifizierung
  für private Pakete.
- `release-upstream-drift`: 1 HIGH + 1 MEDIUM + 2 LOW + 1 INFO — drei
  neue Skripte riefen `curl` bare auf dem Host auf, im Widerspruch zum
  bereits etablierten Container-Muster im selben Verzeichnis
  (`AGENTS.md` §3.1).
- `release-hub-description`: 2 HIGH + 1 MEDIUM + 1 LOW + 2 INFO — beide
  HIGH waren „Zitat nennt die falsche Stelle" (ADR-Abschnitt bzw.
  Welle-Datei-Abschnitt falsch benannt, Aussagen selbst korrekt).
- `release-doku-releasing`: 1 HIGH + 1 INFO — erneut „Zitat nennt die
  falsche Stelle" (interner Selbst-Verweis auf einen falschen
  Abschnitt).

Zusätzlich trat der bereits verkörperte Commit-Traceability-Vorab-Hook-
Mangel (`BEO-PGC/commit-traceability-kein-vorab-hook`) einmal real auf:
ein Fixrunden-Commit (`release-image-scan`) trug zunächst keine
`ADR-*`-Kennung, weil der lokale, opt-in `commit-msg`-Hook in dieser
Arbeitskopie nicht aktiviert war (`core.hooksPath` leer) — behoben über
das etablierte `git reset --soft`-Muster, kein `--amend`/`rebase -i`.

## Trigger-Audit

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 2 — drei Artefaktklassen, je eine
belegte Feststellung.

- **Carveouts (Modul 7):** 0 offen — `docs/plan/carveouts/` ist leer,
  kein Carveout referenziert `ADR-0051` oder eine dieser fünf Slices.
- **Bootstrap-aware Gates (Modul 13):** 0 betroffen — `coverage-gate`
  steht bereits auf der Endstufe (80 %, ausgeschöpft) und wurde von
  keinem der fünf Slices berührt (reiner CI/CD-/Doku-Umfang, kein
  Go-Code).
- **Entscheidung/ADR:** `ADR-0103`s Re-Evaluierungs-Trigger („ein
  Lauf-Zweig braucht einen historischen, über Commits hinweg
  nachschlagbaren Image-Beleg — etwa ein `docker push`-Workflow") ist
  durch diese Welle **noch nicht ausgelöst** — Architect-Prüfung real
  durchgeführt (nicht selbst entschieden): die Archiv-Bedingung verlangt
  einen tatsächlich gepushten, real vorgehaltenen Image-Zustand, keine
  bloße Code-Existenz des Push-Mechanismus (dieselbe Messlatte, die
  `ADR-0103` selbst am 2026-09-18 gegen das damals bereits `Accepted`
  `ADR-0051` anlegte). `git tag -l` ist leer, `release.yml` lief nie
  real. **Beobachtbare Auslösebedingung:** der erste reale `git tag
  vX.Y.Z && git push --tags`, der einen grünen `release.yml`-Lauf
  erzeugt — prüfbar über `git tag -l` (nicht mehr leer) plus einen
  grünen Actions-Lauf. Bis dahin bleibt `ADR-0103` unverändert in Kraft,
  `harness/image-hash.txt` unverändert lokal/nicht-committet. Sobald der
  Trigger real feuert, steht eine weitere Architect-Entscheidung an: ob
  der GitHub-Release-Beschreibungstext als „historisch nachschlagbarer
  Beleg" ausreicht, oder eine Folge-ADR (`Supersedes ADR-0103`) einen
  committeten Digest/Binary-Hash für den Release-Pfad verlangt.

## Steering-Loop-Einträge

Kein neuer Eintrag über der 3×-Schwelle **erstmals durch diese Welle** —
beide in dieser Welle mehrfach getroffenen Beobachtungen waren bereits
vor Wellen-Beginn verkörpert:

- `BEO-PGC/zitat-nennt-die-falsche-stelle` (bereits verkörpert seit
  `welle-d-check`): in dieser Welle allein **drei** neue Fundstellen
  (`release-hub-description` F-1/F-2, `release-doku-releasing` F-1),
  Zähler jetzt 6 Belegdateien/7 Fundstellen insgesamt. Bemerkenswert:
  alle drei Fundstellen dieser Welle sind interne Selbst-Verweise
  innerhalb eines neu geschriebenen Dokuments (Abschnitts-/ADR-Zitat auf
  dieselbe bzw. eine benachbarte Datei) — anders als die vier früheren
  Belege, die auf ein fremdes Dokument verwiesen. Eine mögliche
  Verfeinerung der Beobachtung (Selbst-Verweis vs. Fremd-Verweis als
  getrennte Unterklassen), hier nicht entschieden — Kandidat für den
  nächsten Lese-Schritt, falls die Konzentration sich wiederholt.
- `BEO-PGC/commit-traceability-kein-vorab-hook` (bereits verkörpert seit
  `slice-073`): ein neuer Beleg in dieser Welle (`release-image-scan`),
  Zähler jetzt 5×. Trifft weiterhin dieselbe bereits benannte schwächste
  Stelle der Verkörperung (Hook ist Opt-in, `core.hooksPath` war in
  dieser Arbeitskopie nicht gesetzt).

## Beobachtungs-Register (Zeiger)

Der Zähler steht in
[`docs/plan/planning/observations/`](../observations/). Kein neuer
Ausgang in dieser Welle fällig (beide oben genannten Einträge bereits
verkörpert vor Wellen-Beginn).

## Folge-Slices

Keine. `ADR-0051` ist mit dieser Welle vollständig umgesetzt; was
ansteht, sind keine weiteren Slices, sondern:

- Der erste tatsächliche Release (`git tag` + `git push --tags`) — eine
  irreversible, extern sichtbare Aktion, die nur nach expliziter,
  gesonderter Rückfrage beim Auftraggeber läuft (Welle-Datei §6
  Out-of-Scope), kein technischer Slice-Gegenstand.
- Danach zwei bereits vorbenannte Folge-Entscheidungen: die
  `ADR-0103`-Architect-Entscheidung (siehe Trigger-Audit oben) und die
  Re-Evaluierung von `BEO-PGC/kein-echter-versionswechsel-upgrade-test`
  (`ADR-0064` Re-Evaluierungs-Trigger 1 — durch diese Welle nur zur
  Hälfte erfüllt, der erste echte Tag fehlt noch).

## Archivierung

Dieses Repo führt kein Archivierungs-Werkzeug für Wellen-Zeitdokumente
(kein `archiv`-Ziel in `Makefile`/`harness/mk/*.mk`) — die Bedingung für
Schritt 4 ist nicht eingetreten, keine Handarbeit als Ersatz.

## Verifikation

- `docs/reviews/review-slice-release-version-und-workflow.md` (nach
  Fixrunde: 1 HIGH + 4 MEDIUM real behoben, F-5 bewusst als offenes
  Risiko übernommen) + `docs/reviews/verifikation-slice-release-version-und-workflow.md`.
- `docs/reviews/review-slice-release-image-scan.md` (2 MEDIUM + 2 LOW
  real behoben) + `docs/reviews/verifikation-slice-release-image-scan.md`.
- `docs/reviews/review-slice-release-upstream-drift.md` (1 HIGH + 1
  MEDIUM + 2 LOW + 1 INFO real behoben) +
  `docs/reviews/verifikation-slice-release-upstream-drift.md`.
- `docs/reviews/review-slice-release-hub-description.md` (2 HIGH + 1
  MEDIUM + 1 LOW + 2 INFO real behoben) +
  `docs/reviews/verifikation-slice-release-hub-description.md`.
- `docs/reviews/review-slice-release-doku-releasing.md` (1 HIGH + 1 INFO
  real behoben) + `docs/reviews/verifikation-slice-release-doku-releasing.md`.
- `make gates`: grün auf dem Endstand (778 Dateien, 0 Befunde; Coverage
  82,80 %), ungepiped geprüft (`AGENTS.md` §3.9) nach jedem Commit dieser
  Welle.
- `git tag -l`: leer — bestätigt, dass kein realer Release stattfand
  (Welle-Datei §6 Out-of-Scope eingehalten).
