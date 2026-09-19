# Review-Report: slice-sdk-csharp-publish-workflow — 2026-09-19

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-csharp-publish-workflow.md`),
`ADR-0106` (Accepted) und `AGENTS.md` Hard Rules (Modul 10 §Drei
Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Diff-Range `f48abff2..HEAD` (Abschluss von
`slice-sdk-csharp-pack-werkzeug` bis Implementer-Commit dieses Slice),
Slice `slice-sdk-csharp-publish-workflow`, Welle
`welle-sdk-csharp-lh-fa-sst-009`. Fünf Commits: `377dd7fa` (open→next,
reiner Move, 0 Insertions/Deletions), `0f93e673` (Verantwortlich gesetzt,
2 Insertions/1 Deletion), `1d329a9d` (next→in-progress, reiner Move, 0
Insertions/Deletions), `b0fbfa91` (Inhalt: `.github/workflows/sdk-csharp-release.yml`,
`Makefile`, `tools/harness/semver-regex.sh`, `tools/harness/release-tag-info.sh`
(refactor), `tools/harness/sdk-csharp-release-tag-info.sh`,
`tools/harness/run-sdk-csharp-release-tag-info-tests.sh`,
`harness/README.md`, Slice-Plan-Datei selbst — kein echter Tag-Push
ausgelöst, wie beauftragt).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, seither um mehrere weitere HIGH-Klassen
ergänzt, u. a. `AGENTS.md` §3.9 Pipe-Disziplin, §3.10 Post-Push-Risiko,
§3.13 Träger-Nachzug).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-19.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-sdk-csharp-publish-workflow.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan inkl. Plan-Nachzug, §6 Risiken, §8
  Sub-Area/Modus)
- `docs/plan/adr/0106-csharp-nuget-erstes-sdk-package.md` (Accepted) —
  Festlegung 4 (Trigger/Secret/Tag-Präfix, Alternative E2 gewählt/E3
  verworfen)
- `docs/plan/adr/0051-cicd-pipeline-github-actions.md` Entscheidung 3/8
  (Vorbild Release-Tag-Trigger, Secret-Referenz-ohne-Anlage-Muster)
- `AGENTS.md` §3.1 (Docker-only), §3.3 (git mv + Inhalt = zwei Commits),
  §3.7 (Kommentar-Disziplin), §3.8 (Action-Pinning), §3.9
  (Pipe-/Exit-Code-Disziplin), §3.10 (Post-Push-Risiko GitHub Actions),
  §3.12 (Herkunft von Aussagen), §3.13 (Träger-Nachzug), §4 (kein Träger
  nennt ein nicht existentes Target)
- `harness/README.md` §Sensors/§Werkzeuge (Gate-Index, Vergleichszeilen
  `release.yml`/`hub-description.yml`/`image-scan.yml`)
- `docs/user/releasing.md` (bestehender Release-Prozess-Träger, Vorbild
  für Secret-Tabelle und Trigger-Beschreibung)
- `docs/reviews/review-slice-sdk-csharp-pack-werkzeug.md` — Formvorbild
  für diesen Report

**Eigenständig durchgeführte Prüfungen (nicht nur Commit-Message/Bericht
übernommen):**

- `cat .github/workflows/sdk-csharp-release.yml` vollständig gelesen:
  Trigger ausschließlich `push: tags: ['sdk-csharp-v*']`, kein
  Zweit-Trigger.
- **Permissions-Prüfung (gezielter Auftrag):** Der Workflow ruft
  **keinen** lokalen Workflow per `uses: ./...` auf (kein
  `workflow_call`-Reuse-Pfad wie beim historischen
  `hub-description.yml`-`startup_failure`-Vorfall). Trotzdem real
  geprüft, ob derselbe Fehlermodus ohne den `uses:`-Trigger auftritt: der
  einzige Job `sdk-csharp-release` trägt **selbst** ein
  `permissions: { contents: read }`-Feld (Zeile 57f.), zusätzlich zum
  Top-Level `permissions: {}` (Zeile 50) — er erbt **nicht** stillschweigend
  die leeren Top-Level-Rechte. `actions/checkout` (Zeile 61) hat damit die
  nötige `contents: read`-Berechtigung; `dotnet nuget push` (Zeile 80)
  authentifiziert über `NUGET_API_KEY`, nicht über `GITHUB_TOKEN`, braucht
  also kein zusätzliches GitHub-Permission-Feld. Kein Befund — die Lektion
  aus dem `hub-description.yml`-Vorfall wurde korrekt angewendet, obwohl der
  konkrete Auslöser (`uses:`-Reuse) hier gar nicht vorliegt.
- `grep -n "tags" .github/workflows/ci.yml .github/workflows/e2e.yml
  .github/workflows/release.yml .github/workflows/sdk-csharp-release.yml`
  real ausgeführt und die umliegenden Zeilen gelesen: `ci.yml`/`e2e.yml`
  tragen `tags-ignore: ['**']` (schließt jeden Tag-Push ein), `release.yml`
  trägt `tags: ['v*']`. Ein GitHub-Actions-`tags:`-Glob-Muster verlangt eine
  exakte literale Präfix-Übereinstimmung vor dem `*` — `v*` matcht nur Tags,
  die mit dem Zeichen `v` **beginnen**. `sdk-csharp-v0.1.0` beginnt mit `s`,
  nicht `v`, matcht `v*` also **nicht** (explizit bestätigt, kein impliziter
  Trugschluss) — kein Doppellauf über alle vier Dateien hinweg.
- **SemVer-/`.csproj`-Abgleich-Reihenfolge real an den Zeilennummern
  geprüft** (nicht nur Existenz der Schritte): Schritt-Reihenfolge im YAML
  ist Checkout (Z. 60) → Tag validieren/Version ermitteln (Z. 63) →
  `.csproj`-Version abgleichen, Abbruch bei Abweichung (Z. 66) →
  `make sdk-pack-csharp` (Z. 74) → `dotnet nuget push` (Z. 77) — beide
  Validierungsschritte laufen **vor** jedem Build/Push, exakt wie im Plan
  zugesagt; kein Login-Schritt nötig (Authentifizierung erfolgt inline über
  `--api-key` im Push-Aufruf selbst, kein separater `docker/login-action`-
  Analogon für NuGet).
- **Action-Pinning real gegen die GitHub-API verifiziert** (Netz war
  verfügbar): `curl https://api.github.com/repos/actions/checkout/git/refs/tags/v7.0.1`
  liefert `sha: 3d3c42e5aac5ba805825da76410c181273ba90b1` — identisch mit
  dem im Workflow gepinnten SHA und Tag-Kommentar (Zeile 61); zusätzlich
  intern konsistent mit `release.yml`/`e2e.yml`, die denselben SHA
  wiederverwenden.
- `bash tools/harness/run-release-tag-info-tests.sh` real ausgeführt:
  „alle Fälle bestanden", Exit `0` — der Refactor auf
  `tools/harness/semver-regex.sh` bricht `release-tag-info.sh` **nicht**.
- `bash tools/harness/run-sdk-csharp-release-tag-info-tests.sh` real
  ausgeführt: „alle Fälle bestanden", Exit `0`.
- `ruby -ryaml -e "YAML.load_file('.github/workflows/sdk-csharp-release.yml')"`
  real ausgeführt: kein Parse-Fehler. Zusätzlich jeden der vier
  `run:`-Blöcke einzeln extrahiert und mit `bash -n` geprüft: alle vier
  syntaktisch fehlerfrei.
- `make gates` real ausgeführt, Exit-Code direkt (nicht durch Pipe)
  geprüft: `0` (`d-check: 808 Datei(en) geprüft, 0 Befund(e)`,
  `commit-traceability: OK`, `generated-sync: OK`, `a-check: 0 Befund(e)`).
- `grep -rn "NUGET_API_KEY"` über `.yml`/`.sh`/`.md` real ausgeführt: nur
  Referenzen (`${{ secrets.NUGET_API_KEY }}`, Prosa-Erwähnungen in Plan/ADR/
  README) — keine hartkodierte Zeichenkette, keine Anlage.
- `harness/README.md` §Werkzeuge neue Zeile gegen `release.yml`/
  `hub-description.yml`-Zeilen verglichen: gleiche Form (Trigger, Ablauf,
  Secret-Referenz-ohne-Anlage-Hinweis, `AGENTS.md` §3.10-Verweis, `kein
  Gate, [ADR-Link] · seit slice-<name>`); referenziertes Target
  (`make test-sdk-csharp-release-tag-info`) existiert real im Makefile
  (`AGENTS.md` §4).
- `git show --stat` auf allen fünf Commits geprüft: beide Lifecycle-Moves
  (`377dd7fa`, `1d329a9d`) tragen 0 Insertions/Deletions, das
  `Verantwortlich`-Feld (`0f93e673`) ändert ausschließlich die Slice-Datei
  — `AGENTS.md` §3.3 sauber eingehalten. `git status --porcelain` am Ende
  leer, kein unkommitteter Rest.
- `ADR-0106` §Verglichene Alternativen Tabelle E gelesen: E2 („Docker-only
  Pack, separater Netz-Workflow mit eigenem Tag-Präfix und Secret")
  gewählt, E3 („SDK-Tag im selben `v*`-Namensraum") explizit verworfen mit
  der im Workflow-Kommentar zitierten Begründung — Zitat trägt seinen Satz.
- `docs/user/releasing.md` vollständig gelesen (Kandidat für §3.13-
  Träger-Nachzug) — siehe F-1 unten.
- `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj` real
  gelesen: `<Version>0.1.0</Version>`, exakt das Format, das
  `grep -oE '<Version>[^<]+</Version>'`/`sed` extrahiert — kein zweites
  `<Version>`-Element im Dokument, keine Mehrdeutigkeit der Extraktion.
- Kommentar-Disziplin (`AGENTS.md` §3.7) in allen neuen/geänderten Dateien
  gelesen: Workflow-Kopfkommentar, `semver-regex.sh`,
  `sdk-csharp-release-tag-info.sh`, `run-sdk-csharp-release-tag-info-tests.sh`
  — durchweg indikativ über den aktuellen Zustand, Herkunfts-Anker in der
  bereits im Vorgänger-Review akzeptierten Form
  („slice-sdk-csharp-publish-workflow §6 Risiko 3"), kein Konjunktiv über
  eine verworfene Alternative, kein abwesender Text.

---

## Findings

### F-1 — Neuer Release-Mechanismus (`sdk-csharp-v*`-Tag, `NUGET_API_KEY`) fehlt in `docs/user/releasing.md`

- `kategorie`: MEDIUM
- `quelle`: Maintainability
- `pfad`: `docs/user/releasing.md` (ganzes Dokument, insb. §3/§4/§4
  Secret-Tabelle) — kein Treffer für `sdk-csharp`, `NUGET_API_KEY`,
  `PgChangeFeed.Client`
- `befund`: `docs/user/releasing.md` ist laut eigenem §1-Zweck der
  Träger, der beschreibt „wie ein Release ausgelöst wird" und trägt in §4
  bereits eine Tabelle „Benötigte Repository-Secrets" für den bestehenden
  `v*`-Server-Release sowie in §5 die begleitenden Release-Workflows
  (`image-scan.yml`, `upstream-drift.yml`). Der neue, zweite, unabhängige
  Release-Trigger `sdk-csharp-v*` samt dem neuen Secret `NUGET_API_KEY`
  ist in diesem Dokument nirgends erwähnt — ein Maintainer, der
  `docs/user/releasing.md` liest, um zu erfahren, wie ein Release
  entsteht und welche Secrets dafür nötig sind, erfährt nichts vom
  SDK-Release-Weg. `sdks/csharp/README.md` ist die NuGet-Paket-Doku für
  Consumer, keine Maintainer-Anleitung zum Auslösen eines Releases. Der
  Slice-Plan (§1 „Ausdrücklich NICHT in diesem Slice") listet diesen
  Doku-Träger nicht als bewusst zurückgestellt — die Auslassung wirkt
  unbedacht, nicht begründet.
- `verifizierbar`: nein — kein Sensor prüft Doku-Vollständigkeit dieser
  Art (`docs-check` prüft Referenzen, nicht Abdeckung); Prüfung ist Lesen.
- `klasse`: „Neuer Release-Mechanismus ohne Doku-Zug in
  `docs/user/releasing.md`" (erstes Auftreten)

## Negativbefunde

- geprüft, ohne Befund: Trigger-Isolation — `sdk-csharp-v*` überlappt mit
  keinem der drei Muster in `ci.yml`/`e2e.yml`/`release.yml` (real per
  `grep`/Glob-Semantik bestätigt, kein impliziter Trugschluss).
- geprüft, ohne Befund: Permissions — der Job trägt ein eigenes,
  ausreichendes `permissions: { contents: read }`-Feld, erbt nicht
  stillschweigend `permissions: {}`; derselbe Fehlermodus wie beim
  `hub-description.yml`-`startup_failure`-Vorfall tritt hier nicht auf
  (weder über einen `uses:`-Reuse-Pfad, den es nicht gibt, noch über
  fehlende Job-Permissions, die vorhanden sind).
- geprüft, ohne Befund: Validierungsreihenfolge — SemVer- und
  `.csproj`-Abgleich laufen beide vor jedem Build/Push, per Zeilennummer
  im YAML bestätigt, nicht nur per Existenz der Schritte.
- geprüft, ohne Befund: Action-Pinning — `actions/checkout`-SHA real
  gegen die GitHub-API auf Tag `v7.0.1` verifiziert, zusätzlich intern
  konsistent mit `release.yml`/`e2e.yml`.
- geprüft, ohne Befund: `NUGET_API_KEY` — ausschließlich referenziert
  (`secrets.NUGET_API_KEY`), nirgends hartkodiert oder angelegt.
- geprüft, ohne Befund: gemeinsames `semver-regex.sh`-Refactor — beide
  abhängigen Tabellentests (`run-release-tag-info-tests.sh`,
  `run-sdk-csharp-release-tag-info-tests.sh`) real grün, kein
  Regressionsschaden am bestehenden Server-Release-Pfad.
- geprüft, ohne Befund: YAML-Struktur — Ruby-Stdlib-Parser lädt die Datei
  fehlerfrei; alle vier `run:`-Blöcke einzeln `bash -n`-sauber.
- geprüft, ohne Befund: `harness/README.md` §Werkzeuge — neue Zeile
  formkonsistent mit `release.yml`/`hub-description.yml`, referenziertes
  Target existiert real.
- geprüft, ohne Befund: §6-Risiken — Post-Push-Lauf gegen NuGet.org und
  fehlendes `NUGET_API_KEY`-Secret sind beide korrekt als „weiter offen,
  strukturell" geführt, nicht fälschlich als erledigt markiert; Risiko 3
  (Regex-Code-Teilung) korrekt als „eingetreten und aufgelöst" von den
  beiden anderen unterschieden.
- geprüft, ohne Befund: `AGENTS.md` §3.3 — beide Lifecycle-Moves
  (`377dd7fa`, `1d329a9d`) reine Renames (0/0), `Verantwortlich`-Feld
  eigener Commit (`0f93e673`).
- geprüft, ohne Befund: `AGENTS.md` §3.7 — Kommentare in Workflow-Kopf
  und allen drei neuen/geänderten Shell-Skripten sind indikativ über den
  Ist-Zustand, keine Chronik, keine verworfene-Alternative-Konjunktiv-Form
  abseits des im Vorgänger-Review bereits akzeptierten
  Herkunfts-Anker-Musters.
- geprüft, ohne Befund: Gate-Scope — weder `test-release-tag-info` noch
  `test-sdk-csharp-release-tag-info` hängen an `GATE_CHECKS`
  (`grep`-Prüfung gegen `Makefile`/`harness/mk/*.mk`/`a-check.mk`).
- geprüft, ohne Befund: `make gates` real gefahren, Exit-Code direkt
  (ungepiped) `0`.
- geprüft, ohne Befund: Traceability — alle fünf Commits im Diff-Bereich
  nennen `LH-FA-SST-009` und `ADR-0106` im Betreff, kein `SPEC-*`/`ARC-*`
  im Betreff.
- geprüft, ohne Befund: ADR-Bezug — Alternative E2 (gewählt)/E3
  (verworfen) korrekt zitiert, Zitat trägt seinen Satz (Original in
  `docs/plan/adr/0106-csharp-nuget-erstes-sdk-package.md` §Verglichene
  Alternativen aufgeschlagen, nicht nur aus dem Plan übernommen).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** „Neuer Release-Mechanismus ohne
Doku-Zug in `docs/user/releasing.md`" (erstes Auftreten, kein
Steering-Loop-Zähler-Eintrag nötig unter 3x).

## Verdikt

**Merge-blockierend:** nein, aber **eine Fixrunde ist angezeigt** — 1
MEDIUM-Finding (F-1). Der Reviewer kategorisiert und empfiehlt keine
Lösung; ob `docs/user/releasing.md` in diesem Slice oder als benannter
Aufschub (Folge-Slice-Adresse) nachgezogen wird, entscheidet der
Implementer/Planner.

**DoD-Checkbox-Nachzug ohne Fixrunde** greift hier **nicht**: Da dieses
Verdikt ein MEDIUM-Finding trägt, das regulär an den Implementer
zurückgeht (Rückgabe-Pfeil Reviewer → Implementer), bleibt die DoD-Zeile
„Review durchgeführt, Report unter `docs/reviews/` liegt vor" im
Slice-Plan bis zur Fixrunde bzw. bis zur expliziten Zurückstellung von
F-1 offen — Nachzug erfolgt dann regulär bei Schritt 21 des
Implementer-Workflows.

**Übergabe:** F-1 geht an den Implementer zurück (Fixrunde: entweder
`docs/user/releasing.md` im selben Slice ergänzen, oder den Aufschub mit
einer benannten Folge-Slice-Adresse in §1 „Ausdrücklich NICHT in diesem
Slice" nachtragen). Dieser Report ist ein Lauf-Beleg; er ersetzt keine
Verifikation gegen die volle DoD — das bleibt Verifier-Aufgabe (Modul 11).
