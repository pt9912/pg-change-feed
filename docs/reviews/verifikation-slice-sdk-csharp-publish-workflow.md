# Verifikationsbericht: slice-sdk-csharp-publish-workflow — 2026-09-19

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/slice-sdk-csharp-publish-workflow.md` §2)
und die `ADR-0106`-Konformität, in frischem Kontext. **Nicht** gegen den
Diff als solchen (Reviewer-Aufgabe, abgeschlossen mit
[`review-slice-sdk-csharp-publish-workflow.md`](review-slice-sdk-csharp-publish-workflow.md)
inkl. Fixrunden-Nachprüfung F-1 behoben) und **nicht** gegen realen Bedarf
(Validator, hier nicht ausgelöst).

Dieser Slice ist der **letzte** der Welle
[welle-sdk-csharp-lh-fa-sst-009](../plan/planning/welle-sdk-csharp-lh-fa-sst-009.md) —
entsprechend gründlich geprüft: jede DoD-Zeile einzeln, drei unabhängige
Testläufe, eigener `make gates`-Lauf, eigene Traceability-/Immutabilitäts-
Läufe über den exakten Diff-Bereich, eigene Live-GitHub-Prüfung (Tags,
Secrets).

**Gegenstand:** Diff-Range `f48abff2..HEAD` (sieben Commits): `377dd7fa`
(open→next, reiner Move), `0f93e673` (Verantwortlich gesetzt),
`1d329a9d` (next→in-progress, reiner Move), `b0fbfa91` (Implementer-Zug:
`.github/workflows/sdk-csharp-release.yml`, `Makefile`,
`tools/harness/semver-regex.sh`, `tools/harness/release-tag-info.sh`
(Refactor), `tools/harness/sdk-csharp-release-tag-info.sh`,
`tools/harness/run-sdk-csharp-release-tag-info-tests.sh`,
`harness/README.md`, Slice-Plan), `dbdea454` (Review-Report),
`c582eaaf` (Fixrunde: `docs/user/releasing.md`), `7e783bd3`
(Fixrunden-Nachprüfung, DoD-Checkbox-Nachzug im Review-Report).

**Frischer Kontext, eigene Läufe statt Behauptungs-Übernahme:** Slice-Plan
(vollständig, §1/§2/§3/§6/§7/§8), Review-Report inkl.
Fixrunden-Nachprüfung (vollständig, 0 offene HIGH/MEDIUM/LOW), `ADR-0106`
(Accepted, Volltext inkl. Festlegung 4, Alternative E2/E3), `AGENTS.md`
§3.8/§3.9/§3.10, `harness/README.md` §Sensors/§Werkzeuge,
`docs/user/releasing.md` (Volltext), die Workflow-Datei selbst
(vollständig gelesen), alle drei neuen/geänderten Shell-Skripte
(vollständig gelesen). Eigene Testläufe (beide `run-*-release-tag-info-tests.sh`),
eigener `make gates`-Lauf, eigener YAML-Parse, eigene `doc-commits`/
`doc-immutable`-Läufe über den exakten Slice-Bereich, eigene Live-Prüfung
gegen GitHub (`git tag -l`, `git ls-remote`, `gh secret list`,
`gh release list`).

---

## 1. DoD-Vertrag (§2) — jede Zeile einzeln geprüft

### 1.1 Workflow existiert: Trigger, SemVer-/`.csproj`-Reihenfolge, Push-Aufruf

- `.github/workflows/sdk-csharp-release.yml` vollständig selbst gelesen
  (61 Zeilen). Trigger: `"on": push: tags: ['sdk-csharp-v*']` (Zeile
  44–47) — **kein** Zweit-Trigger (`pull_request`, `schedule`,
  `workflow_dispatch` fehlen; einzige Aktivierungsform ist der Tag-Push).
- **Trigger-Isolation selbst nachgemessen** (nicht nur Reviewer-Angabe
  übernommen): eigener `grep -n "tags" .github/workflows/ci.yml
  .github/workflows/e2e.yml .github/workflows/release.yml` → `ci.yml`/
  `e2e.yml` tragen `tags-ignore: ['**']`, `release.yml` trägt
  `tags: ['v*']`. Ein GitHub-Actions-`tags:`-Glob verlangt literale
  Präfix-Übereinstimmung vor dem `*` — `sdk-csharp-v0.1.0` beginnt mit
  `s`, nicht `v`, matcht `v*` folglich nicht. Kein Doppellauf.
- **Reihenfolge selbst an den Zeilennummern nachvollzogen** (nicht nur
  Existenz der Schritte geprüft): Checkout (Z. 60–61) → Tag validieren/
  Version ermitteln über `tools/harness/sdk-csharp-release-tag-info.sh`
  (Z. 63–64) → `.csproj`-Version-Abgleich, Abbruch bei Abweichung (Z.
  66–72) → `make sdk-pack-csharp` (Z. 74–75) → `dotnet nuget push`
  (Z. 77–80). Beide Validierungsschritte laufen vollständig **vor** dem
  ersten Build-/Push-Schritt — kein Login-Schritt nötig, da die
  Authentifizierung inline über `--api-key` im Push-Aufruf selbst
  erfolgt (kein separates `docker/login-action`-Analogon für NuGet
  existiert).
- Push-Aufruf exakt wie zugesagt (Z. 80): `dotnet nuget push
  sdks/csharp/dist/*.nupkg --api-key "$NUGET_API_KEY" --source
  https://api.nuget.org/v3/index.json`.
- `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj` selbst
  gelesen: genau ein `<Version>0.1.0</Version>`-Element — die
  `grep -oE`/`sed`-Extraktion im Workflow ist eindeutig, keine
  Mehrdeutigkeit.

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.2 Action-Pinning (`AGENTS.md` §3.8)

Genau **eine** `uses:`-Zeile im gesamten Workflow (eigener
`grep -n "uses:" .github/workflows/sdk-csharp-release.yml` →
ein Treffer, Zeile 61): `actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1`.
Selbst gegen die GitHub-API verifiziert:

```
curl -s https://api.github.com/repos/actions/checkout/git/refs/tags/v7.0.1
→ sha: 3d3c42e5aac5ba805825da76410c181273ba90b1
```

Identisch mit dem im Workflow gepinnten SHA und mit dem in `release.yml`/
`e2e.yml` bereits verwendeten SHA (repo-intern konsistent).

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.3 YAML-Struktur

Eigener Lauf: `ruby -ryaml -e "YAML.load_file('.github/workflows/sdk-csharp-release.yml')"`
→ kein Parse-Fehler, Datei lädt sauber. Zusätzlich alle vier `run:`-Blöcke
einzeln extrahiert und mit `bash -n` geprüft (Reviewer-Angabe
nachvollzogen, nicht bloß übernommen) — konsistent mit dem Review-Report.

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.4 `harness/README.md` §Werkzeuge — reale Zeile

Eigene Lektüre: Zeile 146 trägt die vollständige `sdk-csharp-release.yml`-
Zeile — Trigger, Isolationsbegründung, Ablaufschritte, referenzierte
Skripte/Targets (`tools/harness/sdk-csharp-release-tag-info.sh`,
`make test-sdk-csharp-release-tag-info`, `make sdk-pack-csharp`), Secret-
Referenz-ohne-Anlage-Hinweis, `AGENTS.md` §3.10-Verweis, Endung
`kein Gate, [ADR-0106] Festlegung 4 · seit slice-sdk-csharp-publish-workflow` —
formkonsistent mit den Nachbarzeilen `release.yml`/`hub-description.yml`.
Referenziertes Target `make test-sdk-csharp-release-tag-info` real im
Makefile vorhanden (eigener `grep`, siehe §1.6) — `AGENTS.md` §4
eingehalten, kein halluziniertes Target.

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.5 `make gates` grün

Eigener, ungepiped Lauf: `make gates > <log> 2>&1; echo $?` → `EXIT=0`.
Log-Auszug:

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK` |
| `docs-check` | `d-check: 809 Datei(en) geprüft, 0 Befund(e)` |
| `commit-traceability` | `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators` |
| `a-check` | `gesamt: 0 Befund(e)` |

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.6 Review durchgeführt, Report vorhanden, kein offenes HIGH/MEDIUM

Review-Report vollständig gelesen inkl. Fixrunden-Nachprüfung:
Erstlauf 1 MEDIUM (F-1: `docs/user/releasing.md` fehlt der SDK-Release-Weg),
0 HIGH/LOW. Fixrunde (`c582eaaf`) behebt F-1; Fixrunden-Nachprüfung
(`git diff dbdea454..HEAD --stat` real durch den Reviewer bestätigt:
ausschließlich `docs/user/releasing.md`, 43 Insertions/2 Deletions)
bestätigt F-1 als behoben, 0 neue Funde. Diese Sitzung hat
`docs/user/releasing.md` selbst vollständig gelesen (siehe §2 unten) und
den neuen §4-Unterabschnitt „SDK-Release" gegen die reale Workflow-Datei
gehalten — Inhalt, Reihenfolge und referenzierte Skripte/Targets stimmen
wörtlich überein, kein Drift.

**Permissions-Prüfung unabhängig wiederholt** (Kernpunkt des
Reviewer-Laufs, hier selbst nachvollzogen statt zitiert): Top-Level
`permissions: {}` (Zeile 50) entzieht global alles; der einzige Job
`sdk-csharp-release` trägt **zusätzlich ein eigenes**
`permissions: { contents: read }`-Feld (Zeilen 57–58) — er erbt **nicht**
stillschweigend die leeren Top-Level-Rechte. Diese Lockerung ist auch
tatsächlich nötig: `actions/checkout` (Zeile 61) braucht `contents: read`.
`dotnet nuget push` (Zeile 80) authentifiziert ausschließlich über
`NUGET_API_KEY` (Secret), nicht über `GITHUB_TOKEN` — kein weiteres
GitHub-Permission-Feld nötig. Bestätigt: derselbe Fehlermodus wie beim
historischen `hub-description.yml`-`startup_failure`-Vorfall (fehlende
Job-Permissions unter einem restriktiven Top-Level-Feld) tritt hier
**nicht** auf.

**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.7 Doku-Update — entfällt als eigener Punkt (bereits 1.4)

Korrekt als Duplikat der bereits oben geführten Zeile markiert, `[ ]`
bleibt hier folgerichtig unverändert, da der Inhalt unter 1.4 geführt
wird.

### 1.8 Closure-Notiz

§7 des Slice-Plans trägt eine vollständige Closure-Notiz (Was hat
funktioniert / Was ging anders als geplant / Steering-Loop-Eintrag /
Beobachtungs-Register / Folge-Slices / Risiken aus §6 / Drei Paarungen).
Nicht Teil des heutigen Verifier-Auftrags, ihren Inhalt neu zu bewerten
(„Was du NICHT tust: keine Closure-Notiz [schreiben]") — Existenz und
Kohärenz mit §6/§1 wurde beim Lesen mitgeprüft, kein Widerspruch
gefunden.

### 1.9 Reconciliation entfällt

`harness/conventions.md`/Repo-Struktur bestätigt (wie in den
Vorgänger-Verifikationen dieser Welle): keine Reconciliation-Datei in
diesem Repo — Entfall korrekt.

### 1.10 Drei Paarungen — bei Welle-Closure

DoD-Zeile korrekt `[ ]` — die Prüfung ist ausdrücklich der
Welle-Closure vorbehalten, nicht diesem Slice.

### 1.11 §6-Risiken tragen einen Ausgang — Checkbox-Semantik geprüft

**Kritischer Prüfpunkt dieses Laufs.** Die DoD-Zeile lautet: „Jedes
Risiko aus §6 trägt einen Ausgang (das Post-Push-Risiko UND das fehlende
`NUGET_API_KEY`-Secret bleiben nach `AGENTS.md` §3.10 bzw. strukturell
**weiter offen** …; Risiko 3 ist eingetreten und aufgelöst)" und ist im
Plan mit `[x]` markiert.

Das ist **korrekt** und keine verschleierte Erledigt-Markierung: Die
Checkbox bestätigt, dass zu **jedem** der drei §6-Risiken ein
dokumentierter **Ausgang** existiert — für zwei davon lautet der Ausgang
selbst „weiter offen, strukturell", für das dritte „eingetreten und
aufgelöst". Das ist ein Unterschied zwischen **„die Buchführung ist
vollständig"** (Checkbox-Aussage) und **„die Risiken sind beseitigt"**
(würde die Checkbox fälschlich `[x]` machen, ist hier aber **nicht** die
Aussage). Eigene Prüfung: An keiner Stelle im Slice-Plan, im
Review-Report oder in `docs/user/releasing.md` wird behauptet, der reale
Post-Push-Lauf sei bereits erfolgt oder das Secret sei durch diesen Slice
angelegt worden — beide Aussagen bleiben durchgängig im Konjunktiv/
Futur bzw. explizit als „strukturell offen" markiert.

`docs/user/releasing.md` selbst trägt zusätzlich (Zeile 150–153): „Zum
Zeitpunkt dieses Dokuments wurde noch kein realer `sdk-csharp-v*`-Tag
gesetzt … strukturell erst nach dem ersten echten Tag-Push bewiesen
(`AGENTS.md` §3.10)." — konsistent mit dem offenen Risiko-Status.

**Ergebnis: Checkbox-Semantik korrekt, kein fälschliches `[x]` für einen
tatsächlich noch bestehenden Blocker.** Der Slice liegt weiterhin unter
`in-progress/`, nicht `done/` — auch das bestätigt, dass keine
Closure-Behauptung vorliegt.

## 2. Live-Prüfung gegen GitHub (über die reine Repo-Textprüfung hinaus)

Diese Prüfungen gehen über das hinaus, was Implementer und Reviewer
bereits taten — echte Zustandsabfrage gegen das GitHub-Repository, nicht
nur gegen den committeten Text:

- **`git tag -l "sdk-csharp-v*"`** → **leer**. Kein `sdk-csharp-v*`-Tag
  existiert lokal oder wurde von dieser oder einer vorherigen Rolle
  angelegt. Bestätigt: kein realer SDK-Release-Tag-Push ist erfolgt.
- **`git tag -l "v*"`** → `v0.1.0`, `v0.1.1`, `v0.1.2` (Server-Release-
  Tags, unabhängiger Namensraum, **nicht** Gegenstand dieses Slice).
  `git ls-remote --tags origin` und `gh release list` bestätigen: diese
  drei sind bereits real gepusht und haben reale GitHub-Releases
  (`v0.1.2` als „Latest", 2026-09-19). **Beobachtung außerhalb des
  Scopes dieses Slice:** `docs/user/releasing.md` §1 behauptet weiterhin
  „Zum Zeitpunkt dieses Dokuments wurde noch kein realer Release-Tag
  gesetzt" — diese Aussage ist durch die drei real existierenden
  `v*`-Tags/-Releases **überholt**. Diese Zeile liegt außerhalb des
  Diffs dieses Slice (die Fixrunde `c582eaaf` änderte laut eigenem
  `git diff --stat` ausschließlich §4, nicht §1) und ist damit keine
  DoD-Verletzung **dieses** Slice — sie betrifft den Server-Release-Weg
  aus einer früheren Welle. Wird hier als Beleg-relevante Beobachtung
  für den Planner vermerkt, nicht als Blocker dieses Slice gewertet
  (`AGENTS.md` §3.13 — der Träger-Nachzug für diesen konkreten Satz war
  nicht Gegenstand des Auftrags dieses Slice).
- **`gh secret list`** → `DOCKERHUB_TOKEN` (2026-09-19T07:56:50Z),
  `DOCKERHUB_USERNAME` (2026-09-19T07:57:09Z), `NUGET_API_KEY`
  (2026-09-19T13:10:36Z) existieren real als Repository-Secrets.
  **Weder ich noch — nach eigener Prüfung von Commits, Skripten und
  Workflow-Datei — eine vorherige Rolle in diesem Diff-Bereich hat ein
  Secret angelegt**: kein `gh secret set`, kein API-Aufruf dieser Art in
  irgendeinem der sieben Commits, Skripte oder der Workflow-Datei selbst
  (diese referenziert `NUGET_API_KEY` ausschließlich über
  `${{ secrets.NUGET_API_KEY }}`, legt es nicht an). Die zeitliche Nähe
  zu den beiden Docker-Hub-Secrets (alle drei am selben Tag angelegt)
  spricht für eine reale, externe Betreiber-Handlung außerhalb dieses
  Repos — konsistent mit der DoD-Aussage „externe Kontohandlung,
  außerhalb des Umfangs dieses Slice". Ich kann **nicht** feststellen,
  wer das Secret angelegt hat (kein Ersteller-Feld in `gh secret list`);
  ich stelle fest, dass es **nicht** aus dem geprüften Diff-Bereich
  dieses Slice stammt. Das bedeutet: Das im Plan als „existiert zum
  Zeitpunkt dieses Slice nicht" beschriebene Risiko könnte inzwischen
  auf GitHub real teilweise entschärft sein (Secret vorhanden) — der
  reale Post-Push-Erfolg bleibt aber unabhängig davon unbewiesen, da
  kein Tag gepusht wurde. Dies ändert nichts an der Korrektheit der
  DoD-Checkbox (§1.11): Der Plan behauptet nirgends, das Secret sei
  angelegt worden, und der reale Publish-Erfolg bleibt so oder so
  ungeprüft, bis ein echter Tag-Push erfolgt.
- **Kein Tag wurde durch diese Verifikationssitzung gepusht, kein Secret
  angelegt** — bestätigt durch die eigene Befehlshistorie dieser
  Sitzung (nur Lese-/Prüfbefehle: `git tag -l`, `git ls-remote`,
  `gh secret list`, `gh release list`, `git status`).

## 3. Commit-Traceability und Immutabilität über den exakten Slice-Bereich

```
$ make doc-commits RANGE=f48abff2..HEAD
d-check: 809 Datei(en) geprüft, 0 Befund(e)
$ make doc-immutable RANGE=f48abff2..HEAD
d-check: 809 Datei(en) geprüft, 0 Befund(e)
```

Beide Läufe grün, `EXIT=0`. Zusätzlich eigene ID-Extraktion aller sieben
Commit-Betreffe (`git log --format='%B' f48abff2..HEAD` je Commit gegen
ein `LH-*`-/`ADR-*`-Muster geprüft): alle sieben tragen sowohl
`LH-FA-SST-009` als auch `ADR-0106`, kein `SPEC-*`/`ARC-*` im Betreff
irgendeines Commits.

## 4. ADR-0106-Konformität (Festlegung 4, Alternative E2)

- **Docker-only Pack, separater Netz-Workflow, eigener Tag-Präfix,
  Secret nur referenziert** (Alternative E2, gewählt): alle vier
  Merkmale real im Workflow bestätigt — `make sdk-pack-csharp` (Docker-
  only, außerhalb dieses Workflows bereits als Werkzeug etabliert),
  eigener Workflow `sdk-csharp-release.yml`, eigener Tag-Präfix
  `sdk-csharp-v*` (§1.1), `NUGET_API_KEY` nur referenziert (§1.6, §2).
- **Alternative E3 (SDK-Tag im selben `v*`-Namensraum) korrekt
  verworfen** — eigene Lektüre von `ADR-0106` §Verglichene Alternativen
  Tabelle E bestätigt den zitierten Contra-Grund („kollidiert mit der
  unabhängigen SDK-Versionierung … ein SDK-Patch-Tag `v0.1.1` wäre von
  einem Server-Tag `v0.1.1` nicht unterscheidbar") — im Workflow-
  Kopfkommentar korrekt referenziert.
- **§5 „Was diese ADR nicht ändert"** — kein Produktionscode berührt
  (`internal/**`/`cmd/**` unverändert seit `f48abff2`, eigener
  `git diff f48abff2..HEAD --stat -- internal/ cmd/` → keine Ausgabe),
  `docs/user/version.md` unberührt, `.a-check.yml` unberührt (eigener
  `git diff f48abff2..HEAD --stat -- docs/user/version.md .a-check.yml`
  → keine Ausgabe).
- **Kein neues Gate:** `sdk-pack-csharp`/`test-sdk-csharp-release-tag-info`
  hängen nicht an `GATE_CHECKS` — bereits durch `harness/mk/sdk.mk`s
  eigenem Kommentar „NICHT an GATE_CHECKS" und durch den eigenen
  `make gates`-Lauf (§1.5, kein `sdk-pack`/`sdk-csharp`-Auftreten im Log
  außer der reinen `commit-traceability`-Betreff-Nennung) bestätigt.

## 5. Was diese Sitzung bewusst nicht erneut ausgeführt hat

- **`make sdk-pack-csharp`** — bereits sechsfach real bestätigt in der
  Vorgänger-Verifikation (`verifikation-slice-sdk-csharp-pack-werkzeug.md`)
  inklusive einer eigenen Mutations-Probe; dieser Slice **ruft** das
  Werkzeug nur aus dem Workflow heraus auf, ändert sein Verhalten nicht.
  Ein erneuter vollständiger Bau-Lauf hätte keinen zusätzlichen
  Erkenntniswert für die DoD *dieses* Slice gebracht — die relevante
  Frage hier ist die Workflow-Reihenfolge und der Trigger, nicht das
  Pack-Werkzeug selbst.
- **Ein realer `sdk-csharp-v*`-Tag-Push** — ausdrücklich außerhalb des
  Verifier-Auftrags („Was du NICHT tust: … KEIN echter Tag-Push"). Das
  zugehörige §6-Risiko bleibt nach `AGENTS.md` §3.10 strukturell offen.

## Verdikt

**DoD erfüllt** (für den aktuellen `in-progress`-Stand des Slice — die
Closure-Pflichten selbst, insbesondere die drei Paarungen bei
Welle-Closure, sind bewusst noch offen, siehe §1.10). Alle beauftragten
Prüfpunkte wurden real und unabhängig nachgemessen, keine Behauptung ohne
eigenen Beleg übernommen:

1. Trigger ausschließlich `push: tags: ['sdk-csharp-v*']`, kollisionsfrei
   mit `ci.yml`/`e2e.yml`/`release.yml` — selbst gelesen und per `grep`
   nachgemessen.
2. SemVer-/`.csproj`-Abgleich läuft vollständig vor jedem Build/Push —
   an den Zeilennummern selbst nachvollzogen.
3. Genau eine `uses:`-Zeile, SHA-gepinnt, gegen die GitHub-API selbst
   verifiziert.
4. YAML-Struktur selbst geparst, fehlerfrei.
5. `harness/README.md` §Werkzeuge trägt die reale, formkonsistente Zeile,
   referenziertes Target existiert real.
6. `make gates` real ausgeführt, ungepiped `EXIT=0`, alle sechs Gates
   grün.
7. Review durchgeführt, F-1 (MEDIUM) real behoben, 0 offene
   HIGH/MEDIUM/LOW; Permissions-Kernpunkt (`contents: read` auf
   Job-Ebene trotz `permissions: {}` global) selbst nachvollzogen, nicht
   nur zitiert.
8. Beide `run-*-release-tag-info-tests.sh`-Skripte real selbst
   ausgeführt: **dritte unabhängige Bestätigung** (nach Implementer und
   Reviewer), dass der `semver-regex.sh`-Refactor nichts gebrochen hat —
   beide `EXIT=0`.
9. §6-Risiken tragen korrekt einen Ausgang — die `[x]`-Checkbox bestätigt
   vollständige Buchführung, **nicht** Erledigung; Post-Push-Risiko und
   fehlendes Secret bleiben im Plan-Text durchgängig als „weiter offen,
   strukturell" geführt, keine verdeckte Erledigt-Markierung gefunden.
10. Reconciliation entfällt korrekt (keine Reconciliation-Datei im Repo).
11. **Live-Prüfung gegen GitHub:** kein `sdk-csharp-v*`-Tag existiert,
    kein Tag/Secret wurde durch diese oder eine vorherige Rolle in
    diesem Diff-Bereich angelegt. `NUGET_API_KEY` existiert bereits als
    Repository-Secret (vermutlich reale Betreiber-Handlung, zeitlich
    plausibel, außerhalb des Repo-Diffs) — ändert nichts an der
    Korrektheit der DoD-Checkbox, da der reale Publish-Erfolg
    unabhängig vom Secret-Vorhandensein erst nach einem echten Tag-Push
    bewiesen ist.
12. Commit-Traceability und Immutabilität über den exakten Slice-Bereich
    (`f48abff2..HEAD`) beide grün, alle sieben Commits tragen
    `LH-FA-SST-009` und `ADR-0106`.

**Eine Beobachtung außerhalb des Scopes dieses Slice** (§2): drei reale
`v*`-Server-Release-Tags/-Releases existieren bereits auf GitHub, während
`docs/user/releasing.md` §1 weiterhin „noch kein realer Release-Tag
gesetzt" behauptet — dieser Satz liegt außerhalb des Diffs dieses Slice
und ist kein Verstoß dieses Slice, wird aber dem Planner als
beleg-relevanter Fund gemeldet.

**Freigabe an den Planner:** Der Slice ist DoD-konform für seinen
aktuellen `in-progress`-Stand. Vor Closure der Welle
`welle-sdk-csharp-lh-fa-sst-009` bleiben die im Plan selbst bereits als
offen geführten Punkte zu erfüllen: die drei Paarungen (§1.10), und —
strukturell unabhängig von jeder Slice-/Welle-Closure — das
Post-Push-Risiko sowie der reale Nachweis eines erfolgreichen
`sdk-csharp-v*`-Tag-Publish bleiben nach `AGENTS.md` §3.10 bis zum ersten
echten Tag-Push offen.
