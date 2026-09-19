# Review-Report: release-hub-description — 2026-09-19

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/release-hub-description.md`), `ADR-0051`
Entscheidung 8 und `AGENTS.md` Hard Rules (Modul 10 §Drei Review-Arten).
Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Commit `36910422` ("feat(release): Docker-Hub-Beschreibungs-Sync
nach Release ([`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md))"), Slice `release-hub-description`, Welle
`welle-release-pipeline-adr-0051`.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, HIGH-Nachbarform „Zitat nennt die
falsche Stelle" seit welle-d-check gefasst).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-19.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/release-hub-description.md` (Slice-Plan, §1–§3, §6)
- `docs/plan/planning/welle-release-pipeline-adr-0051.md` (Welle-Datei, §1 und §6 Out-of-Scope)
- `ADR-0051` (Accepted) — Entscheidung 8, §Konsistenz-Prüfung (AGENTS §3.1/§3.6)
- `AGENTS.md` §3.1 (Docker-only), §3.6 (Gates ohne ADR), §3.7
  (Kommentar-Disziplin), §3.8 (Action-Pinning), §3.10 (Post-Push-Verifikationspflicht),
  §3.12 (Herkunft von Aussagen)
- `harness/conventions.md` (MR-000 ID-Schema)
- Nachbar-Review `docs/reviews/review-slice-release-version-und-workflow.md`
  (Präzedenzfall: `tools/harness/release-tag-info.sh` als committeter,
  netzloser Test für Inline-Shell-Logik derselben Welle)

**Eigenständig nachgefahrene Prüfungen (nicht vom Implementer-Bericht
übernommen):**

- Vier reale `curl`-Aufrufe gegen die echte, öffentliche Docker-Hub-API
  (`POST /v2/auth/token` mit zu kurzem Fake-Secret → `400` mit
  `{"secret":"must be at least 9 characters"}`; mit korrekt geformtem,
  falschem Secret → `401 unauthorized`; mit leerem Body →
  `400` mit `{"identifier":"value is required","secret":"value is
  required"}`, bestätigt die Feldnamen `identifier`/`secret`; mit
  falschen Feldnamen `username`/`password` → generisches `400` ohne
  Feld-Hinweis, bestätigt, dass `identifier`/`secret` die tatsächlich
  erkannten Feldnamen sind). `PATCH /v2/repositories/pt9912/pg-change-feed`
  mit Fake-Token → `401 unauthorized` (bestätigt Pfad und Repo-Existenz).
- Den vollständigen `run:`-Block aus `hub-description.yml` extrahiert,
  `bash -n` (syntaktisch fehlerfrei) und real zweimal mit den beiden
  oben genannten Fake-Zugangsdaten-Paaren ausgeführt — beide Läufe
  erreichen exakt den behaupteten Fehlerpfad
  (`::error::Docker-Hub-Login fehlgeschlagen — kein access_token in der
  Antwort`, Exit 1).
- `ruby -ryaml` gegen beide geänderten Workflow-Dateien (`hub-description.yml`,
  `release.yml`) — beide strukturell fehlerfrei.
- Realer `make gates`-Lauf auf dem Commit selbst, Exit-Code ungepiped
  geprüft: `0`. `baseline-verify` (v6.9.0, 54 Dateien), `docs-check`/d-check
  (771 Dateien, 0 Befunde, volle Modul-Liste inkl. `hostpaths`),
  `commit-traceability` (5 Commits, alle mit Struktur-ID),
  `coverage-gate` (82.70 % ggü. Schwelle 80 %), `generated-sync`
  (Proto-Erzeugnis byte-gleich), `a-check` (0 Befunde).

---

## Findings

### F-1 — Falsch zitierte ADR-Stelle im Kopfkommentar von `hub-description.yml`

- `kategorie`: HIGH
- `quelle`: Hard-Rule „Beleg trägt seinen Satz nicht" / Nachbarform „Zitat
  nennt die falsche Stelle" (`.harness/skills/reviewer.md` §Klassifikation
  HIGH, `BEO-PGC/zitat-nennt-die-falsche-stelle`)
- `pfad`: `.github/workflows/hub-description.yml:7`
- `befund`: Der Kopfkommentar schreibt „Praesentation, kein Bestandteil
  der Distributions-Zusage ([`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md)
  §Konsistenz-Pruefung)". Der Wortlaut
  „Präsentation, kein Bestandteil der Distributions-Zusage" steht jedoch
  wörtlich in `ADR-0051` **Entscheidung 8** (Zeile 173–174), nicht im
  Abschnitt „Konsistenz-Prüfung gegen bestehende Hard Rules/ADRs" (Zeilen
  93–122, dort geht es um `AGENTS §3.1`/`§3.6`/`ADR-0043`/`ADR-0044`). Ich
  habe beide Stellen im Original aufgeschlagen (`grep -n` gegen die volle
  ADR-Datei): „Präsentation"/„Distributions-Zusage" kommt in der gesamten
  Datei genau einmal vor, und zwar in Entscheidung 8 — §Konsistenz-Prüfung
  enthält diese Formulierung nicht. Der zitierte Anker trägt die Aussage
  nicht, die er stützen soll.
- `verifizierbar`: nein — kein Gate prüft Zitat-Ziel-Treue innerhalb einer
  ADR; nur Aufschlagen der zitierten Stelle zeigt die Abweichung.
- `klasse`: Zitat nennt die falsche Stelle

### F-2 — Falsch zitierte Welle-Datei-Stelle in der DoD-Zeile des Slice-Plans

- `kategorie`: HIGH
- `quelle`: Hard-Rule „Beleg trägt seinen Satz nicht" / Nachbarform „Zitat
  nennt die falsche Stelle" (`.harness/skills/reviewer.md` §Klassifikation
  HIGH)
- `pfad`: `docs/plan/planning/in-progress/release-hub-description.md:59`
  gegen `docs/plan/planning/welle-release-pipeline-adr-0051.md` §1 und §6
- `befund`: Die (im selben Commit auf `[x]` gesetzte) DoD-Zeile schreibt:
  „`DOCKERHUB_USERNAME`/`DOCKERHUB_TOKEN` referenziert, nicht angelegt —
  siehe §1 Abgrenzung der Welle-Datei." Die Welle-Datei hat sieben
  Abschnitte (`## 1. Welle-Ziel` … `## 7. Closure-Notiz`); keiner heißt
  „Abgrenzung". §1 „Welle-Ziel" erwähnt `DOCKERHUB_USERNAME`/
  `DOCKERHUB_TOKEN` nicht. Die tatsächlich zutreffende Stelle ist **§6
  „Out-of-Scope für diese Welle"**, zweiter Punkt: „Anlage der
  `DOCKERHUB_USERNAME`/`DOCKERHUB_TOKEN`-Repository-Secrets — eine
  externe, kontobezogene Handlung … kein technischer Slice-Gegenstand."
  Auch der Slice-Plan selbst hat ein eigenes §1 („Ziel und Abgrenzung"),
  aber auch dessen zwei „Ausdrücklich NICHT in diesem Slice"-Punkte
  erwähnen die Secrets nicht — beide naheliegenden Lesarten von „§1"
  scheitern, der Anker zeigt auf keine Stelle, die die Aussage trägt.
  Diese exakte Formulierung stand bereits vor diesem Commit im Slice-Plan
  (Planner-Text vom selben Tag) und wurde in diesem Commit unverändert
  übernommen, während die Implementer-Änderung die umgebende DoD-Zeile
  von `[ ]` auf `[x]` setzt und damit den Anker faktisch bestätigt.
- `verifizierbar`: nein — kein Gate prüft Zitat-Ziel-Treue in
  Planning-Dokumenten; nur Aufschlagen der Welle-Datei zeigt die
  Abweichung.
- `klasse`: Zitat nennt die falsche Stelle

### F-3 — Keine committete, netzlose Negativtest-Abdeckung für die neue Inline-Shell-Logik

- `kategorie`: MEDIUM
- `quelle`: `.harness/skills/reviewer.md` §Klassifikation MEDIUM
  („fehlende Negativtests bei neuem öffentlichem Vertrag")
- `pfad`: `.github/workflows/hub-description.yml:64-75` (`run:`-Block)
- `befund`: Die einzige nicht-triviale Logik dieses Slices — die
  Fehlererkennung „kein `access_token` in der Antwort → `exit 1`" — ist
  ausschließlich durch Ad-hoc-Läufe (Implementer und dieser Review) gegen
  die echte, öffentliche Docker-Hub-API belegt, nicht durch ein
  committetes, wiederholbares, netzloses Testskript. Dieselbe Welle hat
  für eine strukturell vergleichbare Situation (Inline-Shell-Logik in
  `release.yml`, dort die SemVer-Validierung) bereits das Präzedenzmuster
  `tools/harness/release-tag-info.sh` +
  `tools/harness/run-release-tag-info-tests.sh` + `make
  test-release-tag-info` etabliert (netzlose Tabellentests, real im
  Repo vorhanden). Der Token-Extraktions-/Fehlerpfad ließe sich mit
  kanonischen, canned JSON-Antworten (kein Netzzugriff nötig) ebenso
  abdecken, ist es hier aber nicht.
- `verifizierbar`: nein — kein Gate erzwingt Tests für
  Workflow-Inline-Shell-Logik.
- `klasse`: fehlende Negativtests bei neuem öffentlichem Vertrag

### F-4 — `curl -f`-Asymmetrie (Login ohne, PATCH mit) nur in der Commit-Message erklärt, nicht im Code

- `kategorie`: LOW
- `quelle`: Maintainability (angrenzend an `AGENTS.md` §3.7, aber keine
  Verletzung — es gibt keinen irreführenden Kommentar, nur die Abwesenheit
  eines erklärenden)
- `pfad`: `.github/workflows/hub-description.yml:65-75`
- `befund`: Die Commit-Message begründet, warum `curl -f` beim
  Login-Aufruf bewusst fehlt („bricht sonst vor der eigenen
  Fehlerbehandlung ab") und beim PATCH-Aufruf steht. Diese Begründung
  lebt ausschließlich in der Commit-Message, nicht als Kommentar an der
  Stelle selbst — ein künftiger Editor, der die scheinbare Inkonsistenz
  „beheben" will (z. B. `-f` beim Login ergänzen), hat im Diff keinen
  Hinweis, dass das die Fehlerbehandlung sabotiert. Ich habe die
  Begründung unabhängig nachvollzogen und für zutreffend befunden
  (GitHub Actions ruft `run:`-Schritte standardmäßig mit `bash -eo
  pipefail` auf — `curl -f` hätte bei `4xx` also einen sofortigen,
  unkontrollierten Skript-Abbruch vor der eigenen `::error::`-Meldung
  zur Folge).
- `verifizierbar`: nein.
- `klasse`: Erklärung lebt in Commit-Message statt im Code

### F-5 — Erfolgspfad-Feldnamen (`access_token`, `full_description`) bleiben bis zum ersten realen Lauf unbewiesen

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.10 (analog) / Maintainability
- `pfad`: `.github/workflows/hub-description.yml:66-75`
- `befund`: Alle real durchgeführten Prüfungen (meine eigenen und die im
  Commit behaupteten) treffen ausschließlich Fehlerpfade (`400`/`401`
  ohne gültige Zugangsdaten). Die tatsächliche Form einer *erfolgreichen*
  Antwort (`{"access_token": "…"}` bei `POST /v2/auth/token`,
  Erfolgs-Statuscode bei `PATCH .../repositories/...`) bleibt strukturell
  unbewiesen, bis der Workflow einmal mit echten
  `DOCKERHUB_USERNAME`/`DOCKERHUB_TOKEN`-Secrets läuft — dieselbe
  Kategorie Risiko, die §6 des Slice-Plans bereits als „weiter offen"
  führt, hier nur auf die konkreten Feldnamen zugespitzt.
- `verifizierbar`: nein.
- `klasse`: Erfolgspfad strukturell unbeweisbar vor erstem realen Lauf

### F-6 — Permissions-Vererbung beim `uses:`-Aufruf ohne expliziten `permissions:`-Block im Caller-Job

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `.github/workflows/release.yml:108-113` (`hub-description`-Job)
- `befund`: Der aufrufende Job `hub-description` in `release.yml` setzt
  keinen eigenen `permissions:`-Block; er erbt damit den
  Workflow-Top-Level-Default `permissions: {}` von `release.yml`. Die
  aufgerufene, wiederverwendbare Workflow-Datei `hub-description.yml`
  deklariert ihrerseits `contents: read` auf Job-Ebene für ihren
  `actions/checkout`-Schritt. Ob GitHub Actions bei einem lokalen
  `uses: ./…yml`-Aufruf ohne explizite `permissions:`-Angabe am
  Aufrufer die im aufgerufenen Workflow deklarierten Job-Permissions
  tatsächlich gewährt oder den Aufrufer-Default (`{}`) durchreicht, lässt
  sich nicht netzlos/lokal abschließend klären — es ist derselbe
  strukturelle Vorbehalt wie `AGENTS.md` §3.10 (erst der reale
  Post-Push-Lauf zeigt, ob `Checkout` in `hub-description.yml` mit
  ausreichenden Rechten läuft), hier nur auf eine konkretere, bislang
  nicht benannte Einzel-Fragestellung zugespitzt.
- `verifizierbar`: nein — nur ein realer Workflow-Lauf zeigt das
  tatsächliche Verhalten.
- `klasse`: Reusable-Workflow-Permissions-Vererbung ungeklärt

## Negativbefunde

- geprüft, ohne Befund: Docker-Hub-API-Endpunkte, Feldnamen
  (`identifier`/`secret`/`access_token`/`full_description`) im
  Fehlerpfad — vier eigene reale `curl`-Aufrufe gegen die echte API
  bestätigen Pfad, Feldnamen und Repo-Existenz (siehe „Eigenständig
  nachgefahrene Prüfungen" oben).
- geprüft, ohne Befund: Shell-Logik-Korrektheit — `set -uo pipefail`
  ohne `-e` ist konsistent mit GitHub Actions' eigenem
  `bash -eo pipefail`-Default für `run:`-Schritte (die lokale
  `set`-Zeile ändert daran nichts, ist aber auch nicht schädlich); `curl
  -sS` (ohne `-f`) beim Login lässt `4xx`-Antworten den eigenen
  `access_token`-Check erreichen statt das Skript vorzeitig abzubrechen —
  real mit zwei Fake-Zugangsdaten-Paaren reproduziert, beide erreichen
  exakt den behaupteten Fehlerpfad; `curl -fsS` beim PATCH-Aufruf ist
  korrekt, da dieser der letzte Schritt ist. `jq -n --arg desc "$(cat
  README.md)"` kapselt den README-Inhalt korrekt als JSON-String-Wert
  ohne Injection-Risiko (kein `eval`, keine Wortaufspaltung durch die
  doppelten Anführungszeichen um die Kommandosubstitution).
- geprüft, ohne Befund: Struktur des neuen `hub-description`-Jobs in
  `release.yml` — `needs: release`, `uses:
  ./.github/workflows/hub-description.yml`, `secrets: inherit` sind
  syntaktisch und semantisch korrekt für einen lokalen
  Reusable-Workflow-Aufruf; `hub-description.yml`s eigenes `on:
  workflow_call: secrets: {DOCKERHUB_USERNAME, DOCKERHUB_TOKEN}` (beide
  `required: true`) passt dazu. Die Behauptung „ein Fehlschlag lässt den
  bereits abgeschlossenen `release`-Job unberührt, GitHub Actions ändert
  dessen Status nicht rückwirkend" ist technisch zutreffend — Job-Status
  in GitHub Actions sind nach Abschluss fixiert und werden durch
  nachgelagerte `needs`-Jobs nicht rückwirkend geändert (die
  Gesamt-Run-Konklusion des Workflow-Laufs selbst kann trotzdem
  „failure" zeigen, das ist aber eine andere, von der DoD-Aussage nicht
  behauptete Ebene, und ohne Branch-Protection-Bindung an `release.yml`
  ohne blockierende Wirkung).
- geprüft, ohne Befund: Docker-only-Grenze (`AGENTS.md` §3.1) — der
  direkte `curl`/`jq`-Aufruf im `run:`-Block läuft auf dem GitHub-hosted
  Runner ohne Docker-Kapselung. `ADR-0051` selbst zieht diese Grenze
  explizit in seiner „Konsistenz-Prüfung gegen bestehende Hard
  Rules/ADRs": „Wo ein Schritt reine Workflow-Mechanik ist (Tag-Parsing,
  Registry-Login, Digest-Vergleich zwischen zwei Registries), bleibt das
  Workflow-Skript — es ersetzt kein `make`-Target, es orchestriert sie."
  Der Docker-Hub-Beschreibungs-Sync dupliziert kein bestehendes
  `make`-Target und ist derselben Kategorie „reine Workflow-Mechanik"
  wie `release.yml`s bereits etablierte, unkapselte `gh`-CLI-Nutzung
  (dort im Vorgänger-Slice ohne Beanstandung reviewt). Die
  Docker-only-Pflicht für `tools/harness/*.sh`-Dev-Skripte
  (`tools/harness/lib-github-api.sh`, Container-Wrapper für lokale
  `make pin-stale-*`-Aufrufe auf einem Entwicklerrechner) ist ein
  anderer Ausführungskontext und bleibt von diesem Slice unberührt —
  kein neues `tools/harness/*.sh`-Skript wurde hier eingeführt.
- geprüft, ohne Befund: Action-Pinning (`AGENTS.md` §3.8) —
  `hub-description.yml`s einzige `uses:`-Zeile
  (`actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1`)
  ist SHA-gepinnt mit Tag-Kommentar, identisch mit dem bereits in
  `ci.yml`/`e2e.yml`/`release.yml` verwendeten Pin (im Vorgänger-Review
  bereits gegen `git ls-remote --tags` verifiziert).
- geprüft, ohne Befund: `AGENTS.md` §3.6 (Gates dürfen nicht ohne ADR
  gelockert werden) — kein bestehendes Gate wird berührt;
  `hub-description.yml` ist explizit „kein Gate" dokumentiert.
- geprüft, ohne Befund: `harness/README.md` — beide neuen/geänderten
  Zeilen (`release.yml`, `hub-description.yml`) verlinken `ADR-0051`
  korrekt (`../docs/plan/adr/0051-cicd-pipeline-github-actions.md`,
  Ziel existiert, Pfad von `harness/README.md` aus korrekt aufgelöst);
  Markdown-Tabellenstruktur bleibt intakt (gleiche Feldanzahl wie der
  Tabellenkopf).
- geprüft, ohne Befund: kein host-lokaler absoluter Pfad (`AGENTS.md`
  §3.11) in den vier geänderten Dateien.
- geprüft, ohne Befund: Traceability — Commit-Message nennt `ADR-0051`;
  keine `LH-*`-ID nötig, da Prozess-ADR ohne Spec-Stratum (`ADR-0051`
  §Berührte Spec-Stellen).
- geprüft, ohne Befund: `make gates` — realer Lauf auf dem Commit,
  Exit-Code ungepiped `0` (siehe „Eigenständig nachgefahrene
  Prüfungen").
- geprüft, ohne Befund: YAML-Struktur beider geänderter Workflow-Dateien
  (`ruby -ryaml`) und `bash -n` auf dem extrahierten `run:`-Block.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zitat nennt die falsche Stelle (2×,
F-1/F-2) · fehlende Negativtests bei neuem öffentlichem Vertrag ·
Erklärung lebt in Commit-Message statt im Code · Erfolgspfad strukturell
unbeweisbar vor erstem realen Lauf · Reusable-Workflow-Permissions-Vererbung
ungeklärt.

## Verdikt

**Merge-blockierend:** ja — F-1 und F-2 sind beide Instanzen der
etablierten HIGH-Klasse „Zitat nennt die falsche Stelle": ein genannter
Beleg, der die Aussage nicht trägt, ist unabhängig von der inhaltlichen
Richtigkeit der Aussage selbst ein Befund (die Aussagen dahinter — „kein
Bestandteil der Distributions-Zusage" und „Secrets werden extern
angelegt" — sind beide sachlich zutreffend, nur die zitierten Anker
falsch). Beide sind mechanisch billig zu beheben (Zitat auf `ADR-0051`
Entscheidung 8 statt „§Konsistenz-Prüfung" bzw. auf Welle-Datei §6 statt
„§1" korrigieren). F-2 betrifft Text, der ursprünglich vom
Planner-Agenten stammt und im Implementer-Commit nur unverändert
mitgeführt und per DoD-Checkbox bestätigt wurde — die Korrektur kann
trotzdem im nächsten Implementer- oder Planner-Zug erfolgen, je nachdem
wer als Nächstes an dieser Datei arbeitet; das ist eine
Zuständigkeitsfrage, keine Entwarnung. F-3 (MEDIUM) ist kein Blocker für
die Kern-Mechanik (alle real getesteten Pfade verhalten sich wie
behauptet), gehört aber vor Closure entschieden, da dieselbe Welle
bereits den strengeren, netzlos-testbaren Präzedenzfall
(`release-tag-info.sh`) etabliert hat.

**Übergabe:** Da mindestens ein HIGH-Finding eine Fixrunde auslöst, ziehe
ich die DoD-Checkbox „Review durchgeführt" im Slice-Plan **nicht** nach
(Skill-Regel „DoD-Checkbox-Nachzug ohne Fixrunde" greift nur bei 0 HIGH).
Die Finding-Klassen gehen bei Slice-Closure ins Beobachtungs-Register —
„Zitat nennt die falsche Stelle" ist bereits ein bekanntes,
mehrfach-belegtes Muster (`BEO-PGC/zitat-nennt-die-falsche-stelle`,
jetzt ein fünftes/sechstes Vorkommen über zwei neue Fundstellen in
diesem Slice). Dieser Report ersetzt keine Verifikation gegen die DoD —
das bleibt Verifier-Aufgabe (Modul 11).
