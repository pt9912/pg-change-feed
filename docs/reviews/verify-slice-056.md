# Verifikationsbericht: slice-056 — 2026-09-14

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-056` §1 Ziel/Abgrenzung, §2 DoD, §4 Trigger, §5 Closure-Trigger, §6
Risiken, §8 Sub-Area) und die bindende `ADR-0051` (Entscheidung 2) — nicht
gegen Diff (Reviewer-Aufgabe, bereits abgeschlossen:
`docs/reviews/review-slice-056.md`, vollständig gelesen, nicht wiederholt)
und nicht gegen realen Bedarf (Validator, hier nicht ausgelöst — reine
CI-Infrastruktur, kein neuer Architektur-Sicht-Meilenstein).

**Frischer Kontext:** Diese Prüfung liest den vollständigen Slice-Plan,
`ADR-0051` vollständig, den Review-Report vollständig, `.github/workflows/e2e.yml`
und den `harness/README.md`-Eintrag selbst — keine Implementer- oder
Reviewer-Behauptung wird ungeprüft übernommen. `make gates` wurde in dieser
Sitzung **eigenständig real ausgeführt** (Exit-Code separat geprüft, kein
Pipe-Masking, `AGENTS.md` §3.9). Der reale GitHub-Actions-Lauf wurde
**eigenständig** über `gh run list`/`gh run view --log` sowie die
Branch-Protection-Konfiguration **eigenständig** über `gh api` abgefragt —
keine Zahl aus dem Review-Report wurde ungeprüft übernommen.

**Gegenstand:**
`docs/plan/planning/in-progress/slice-056-e2e-workflow-in-ci.md` zum Stand
`HEAD = b88e791`. Slice-eigene Commits: `9666725`/`f2552e6`/`b1cba09`
(Vorbereitung), `b3f0793` (`next→in-progress`, reiner Move), `e6aa155`
(Implementierung: `.github/workflows/e2e.yml` neu, `harness/README.md`
§Werkzeuge-Eintrag), `9985523`/`d5d4d03` (Review-Report, 0 HIGH/1 LOW).
Nachfolgende Commits (`d779e30`, `b88e791`) betreffen eine unabhängige
Session-Beobachtung (`BEO-PGC/report-nackte-id-ohne-link`) und berühren
`slice-056`s eigenen Liefer-Umfang nicht — per `git diff b3f0793..HEAD
--stat` bestätigt (siehe §4).

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | `.github/workflows/e2e.yml` neu: Trigger, `permissions: {}`/`contents: read`, Action-Pinning, ein Job ruft ausschließlich `make`-Targets auf | **erfüllt, selbst reproduziert** | Datei selbst gelesen: `"on": pull_request` + `push`/`branches: ['**']`/`tags-ignore: ['**']` — zeichengleich zu `ci.yml`. `permissions: {}` auf Workflow-Ebene, `contents: read` auf Job-Ebene. Zwei Steps: `run: make image`, `run: make test-integration` — kein Inline-Shell. Action-Pin `actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1` **eigenständig gegen die GitHub-API verifiziert**: `gh api repos/actions/checkout/git/refs/tags/v7.0.1` liefert exakt denselben SHA (`3d3c42e5aac5ba805825da76410c181273ba90b1`) — der Kommentar lügt nicht. |
| 2 | Workflow real über Push/PR ausgelöst, sichtbares Check-Ergebnis | **erfüllt, selbst reproduziert** | `gh run list --workflow=e2e.yml --limit 5` (selbst ausgeführt): vier aufeinanderfolgende reale Läufe auf `main`, alle `completed`/`success` — `34788082848` (3m41–44s, Commit `e6aa155`), `34788361676` (3m33s), `34788439123` (3m46s), `34788866667` (3m29s). `gh run view 34788082848` bestätigt Job `image + test-integration` grün in 3m41s. `gh run view 34788082848 --job=103807124807 --log` (voller Job-Log, selbst gelesen): reale Testphasen sichtbar (Retention-Lebenszyklus, Black-Box-CLI-Rundlauf, SQL-Administration-Live-Reload, Publication-Entzug, drei NATS-Testabschnitte, Schema-Change-Beleg), endet mit `ok  .../test/integration 0.143s`. `BEO-PGC/github-actions-unverifizierbar-lokal` ist für diesen Workflow-Typ damit real (nicht nur behauptet) aufgelöst. |
| 3 | Real belegt: fehlschlagender Testlauf lässt PR weiterhin mergebar erscheinen (kein Required-Check), Beleg nicht eingecheckt | **erfüllt — strukturell, eigenes Urteil, siehe §3 unten** | `gh api repos/pt9912/pg-change-feed/branches/main/protection` (selbst ausgeführt) → `404 "Branch not protected"`. Kein dedizierter Rot-Lauf/Test-PR wurde erzeugt (weder von Implementer noch Reviewer noch mir) — Begründung und Einordnung in §3. |
| 4 | `make gates` grün | **erfüllt, selbst reproduziert** | Eigener Lauf, Exit-Code separat geprüft (kein Pipe): `EXIT_CODE=0`. Siehe §2 für den vollständigen Ausschnitt. |
| 5 | Review durchgeführt, Report liegt vor | **erfüllt** | `docs/reviews/review-slice-056.md` vollständig gelesen: 0 HIGH, 1 LOW (F-1, DoD-Wortlaut-Präzisierung ohne Rollen-Widerspruch, kein Fixrunden-Pfad). DoD-Zeile im selben Commit (`9985523`) korrekt nachgezogen (`git diff b3f0793..HEAD` zeigt exakt diese eine Checkbox-Änderung an der Slice-Datei). |
| 6 | Doku-Update `harness/README.md` §Werkzeuge | **erfüllt (Inhalt), Checkbox korrekt noch offen — Planner-Closure-Arbeit** | `harness/README.md` Zeile 132 real gelesen: neuer Eintrag in der **Werkzeuge**-Tabelle, Bindung `kein Gate, ADR-0051, LH-QA-POR-003` — korrekt in der Werkzeuge- statt Sensors-Tabelle (kein Gate). Inhalt bereits im Commit `e6aa155` geliefert; die Checkbox selbst wird laut Rollen-Sequenz (Modul 8) erst bei der Planner-Closure gesetzt — kein DoD-Mangel. |
| 7 | Closure-Notiz mit Steering-Loop-Lerneintrag (§7) | **korrekt offen** | §7 trägt noch ausschließlich Platzhalter — Planner-Closure-Arbeit, beginnt erst nach diesem Bericht. |
| 8 | Reconciliation-Register | **korrekt entfällt** | `docs/plan/planning/reconciliation.md` existiert nicht — Repo durchgehend GF (`harness/conventions.md` Modus-Deklaration `*`/`PGC`). |
| 9 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | `docs/plan/planning/observations/BEO-PGC/github-actions-unverifizierbar-lokal/evidence/` real aufgelistet: nur `slice-039.md` vorhanden, kein `slice-056.md`. Zähler steht real bei 1×, nicht bei 2× — konsistent mit „noch nicht fortgeschrieben". |
| 10 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen — eigenes Urteil siehe §5 unten** | Alle drei Zeilen in §6 tragen wörtlich `<bei Closure einzutragen>`, real per Lektüre bestätigt. |
| 11 | Drei Paarungen | **korrekt offen** | Slice liegt noch in `in-progress/`; die Paarungen suchen in `done/` und sind erst nach dem `git mv` sinnvoll prüfbar. |

**Zwischenbefund:** Die vier real prüfbaren technischen/funktionalen Punkte
(1, 2, 4, 5) sind **selbst reproduziert erfüllt**. Punkt 6 ist inhaltlich
erfüllt, Checkbox korrekt noch offen. Punkt 3 ist in der Substanz erfüllt,
aber über einen anderen Beleg-Weg als im DoD-Wortlaut vorgesehen — siehe §3,
eigenes Urteil, kein Ausweichen. Die verbleibenden Punkte (7, 9–11, plus der
bereits erledigte Reconciliation-Entfall bei 8) sind **korrekt noch offen** —
Planner-Closure-Arbeit nach Modul 8 §Rollen-Sequenz für einen Slice. Kein
DoD-Verstoß.

## 2. Sensor-Läufe (selbst ausgeführt)

**`make gates`** (vollständiger Lauf, Exit-Code separat und explizit
geprüft — `EXIT_CODE=0`, kein Pipe-Masking):

```
baseline-verify: v6.5.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)
coverage-gate: OK — Coverage 41.00% erfüllt Schwelle 35%
d-check: 450 Datei(en) geprüft, 0 Befund(e)
d-check (commits, HEAD~5..HEAD): 450 Datei(en) geprüft, 0 Befund(e)
commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID
a-check: gesamt: 0 Befund(e)
  Hinweis: tools/harness/natssub/main.go liegt in keiner Schicht (dokumentiertes,
  nicht-fatales Verhalten, unverändert seit vor diesem Slice)
```

**`gh run list --workflow=e2e.yml --limit 5`** (selbst ausgeführt, Exit 0):
vier reale, aufeinanderfolgende Läufe auf `main`, alle `completed`/`success`,
Laufzeiten 3m29s–3m46s — deutlich und konsistent unter dem 60-Minuten-Timeout.

**`gh run view 34788082848`** (selbst ausgeführt, Exit 0): Job
`image + test-integration` grün in 3m41s.

**`gh run view 34788082848 --job=103807124807 --log`** (selbst ausgeführt,
voller Job-Log gelesen): reale Testphasen sichtbar (u. a. Retention-Sichtbarkeit,
Black-Box-CLI-Rundlauf, SQL-Administration-Live-Reload, Publication-Entzug-
Wirksamkeit, NATS-Happy-Path/Boundary/Negative-Belege, Schema-Change-Beleg),
endet mit `ok  github.com/pt9912/pg-change-feed/test/integration 0.143s`.

**`gh api repos/pt9912/pg-change-feed/branches/main/protection`** (selbst
ausgeführt): `404 "Branch not protected"` — `main` trägt aktuell **keine**
Branch-Protection-Regel jeder Art.

**`gh api repos/actions/checkout/git/refs/tags/v7.0.1`** (selbst ausgeführt):
SHA `3d3c42e5aac5ba805825da76410c181273ba90b1` — identisch zum Pin in
`e2e.yml`.

## 3. DoD-Punkt 3 — Rot-Beleg für Nicht-Blockierung: eigenes Urteil

**Befund:** Weder Implementer noch Reviewer haben einen dedizierten
Test-PR mit absichtlich scheiterndem Schritt erzeugt (der Review-Report
benennt das selbst explizit als „weiterhin offen für den Verifier"). Ich
habe ebenfalls **keinen** solchen Test-PR erzeugt — stattdessen die
Branch-Protection-Konfiguration von `main` direkt über die GitHub-API
abgefragt: `gh api repos/pt9912/pg-change-feed/branches/main/protection`
liefert `404 "Branch not protected"`.

**Urteil:** Das ist kein Ausweichen, sondern ein stärkerer Beleg als der im
DoD-Wortlaut vorgesehene Einzel-Testlauf, und er trägt die Substanz des
Punkts vollständig:

- Ein „Required Status Check" ist bei GitHub **ausschließlich** über eine
  Branch-Protection-Regel konfigurierbar. Ohne jede Protection-Regel auf
  `main` — nicht nur ohne eine, die `e2e.yml` nennt — kann **kein**
  Workflow, gleich welchen Namens und gleich welchen Ergebnisses, einen
  Merge blockieren. Ein einzelner Rot-Lauf von `e2e.yml` hätte nur *diesen
  einen Fall* gezeigt; die Konfigurationsabfrage zeigt die *Klasse* — sie
  gilt für jeden künftigen roten Lauf von `e2e.yml` ebenso wie für `ci.yml`
  oder einen noch nicht existierenden Workflow, nicht nur für den einen
  getesteten Moment.
  Genau das erwartet Plan §1 auch strukturell: die
  Required-Status-Check-Konfiguration ist dort ausdrücklich als
  Repository-Admin-Aktion außerhalb des Repo-Inhalts benannt, die dieser
  Slice bewusst nicht anfasst — es gibt sie heute schlicht nicht.
- Ein jetzt eigens erzeugter Rot-Lauf (temporärer Test-PR mit
  scheiterndem Schritt) hätte gegenüber diesem Befund **keine zusätzliche
  Information** geliefert — er hätte nur bestätigt, was die
  Repository-Konfiguration bereits kategorisch ausschließt. Er hätte
  zusätzlich einen echten, wenn auch kurzlebigen PR erzeugt, den die DoD
  gerade *nicht* eingecheckt sehen will.

**Ergebnis:** Ich stufe DoD-Punkt 3 als **in der Substanz erfüllt** ein,
über einen stärkeren, strukturellen Beleg statt über den im Wortlaut
vorgesehenen Einzel-Testlauf. **Für den Planner vorgemerkt:** Das ist
dieselbe Finding-Klasse wie Review-F-1 („DoD-Wortlaut deckt gelieferten
Umfang nicht exakt") — der Wortlaut verlangt einen Verhaltensbeleg
(Testlauf), belegbar und tatsächlich belegt ist die Eigenschaft aber über
einen Konfigurationsbeleg. Kein Fixrunden-Fall, keine neue
Beobachtungsregister-Kennung nötig (Einzelfall, keine wiederkehrende
Klasse) — aber in der Closure-Notiz sollte stehen, *welcher* der beiden
Beleg-Wege tatsächlich trug.

## 4. Plan-vs-Code-Diff

- **Datei-Liste (§3 des Plans) exakt getroffen:** `git show --stat e6aa155`
  zeigt genau zwei Dateien — `.github/workflows/e2e.yml` (neu, 63 Zeilen)
  und `harness/README.md` (1 Zeile) —, deckungsgleich mit der Plan-Tabelle.
- **Kein Out-of-Scope-Punkt (§1) berührt:** `git diff b3f0793..HEAD --stat`
  (selbst ausgeführt) zeigt zusätzlich zu den beiden Slice-Dateien nur
  Dateien der unabhängigen Session-Beobachtung
  `BEO-PGC/report-nackte-id-ohne-link` (Observation-Verzeichnis,
  Architect-Verdikt-Report, `docs/reviews/review-slice-056.md` selbst) —
  keine davon berührt `test-store`/`test-replication`,
  `tools/harness/run-integration-tests.sh`, `release.yml` o. ä., und keine
  Branch-Protection-Konfiguration wurde angefasst (bestätigt durch den
  404-Befund in §2/§3 — sie existiert nach wie vor nicht). `.github/dependabot.yml`
  unverändert.
- **Keine unbegründete Abweichung gefunden.** Die einzige Differenz zum
  Plan-Wortlaut ist die in Review-F-1 bereits benannte und in §3 dieses
  Berichts vertiefte DoD-Wortlaut-Frage — kein Code-Diff-Problem.

## 5. §6-Risiken — eigenes, unabhängiges Urteil (kein Ausgang eingetragen — Planner-Arbeit)

- **Risiko 1 — GitHub-Actions-Workflows lassen sich nicht in Docker
  simulieren; ob `make test-integration` auf dem echten Runner innerhalb
  des Zeitlimits durchläuft, blieb bis zum ersten realen Push unbewiesen.**
  Eigene Prüfung: vier unabhängige, reale Läufe (siehe §2), alle grün.
  Meine Einschätzung: **trägt für den Ausgang „entfallen", mit Begründung**
  — die konkrete Unsicherheit dieses Slice (läuft `test-integration` auf
  dem GitHub-hosted Runner überhaupt durch?) ist durch vier reale,
  reproduzierte grüne Läufe ausgeräumt. Die *allgemeinere* Beobachtung
  `BEO-PGC/github-actions-unverifizierbar-lokal` (jeder künftige neue
  Workflow-Typ bleibt vor seinem ersten Push ungeprüft) bleibt davon
  unberührt offen — das ist eine andere, umfassendere Aussage als dieses
  konkrete Risiko.
- **Risiko 2 — Timeout-Knappheit: `make test-integration` deutlich länger
  als `make gates`/`make test`, Gesamtlaufzeit könnte das 60-Minuten-Limit
  auf dem Standard-Runner knapp werden lassen.** Eigene Prüfung: vier reale
  Läufe zwischen 3m29s und 3m46s — durchgehend **unter 7 % des
  Zeitbudgets**, keiner davon ein Ausreißer nach oben. Meine Einschätzung:
  **trägt für den Ausgang „entfallen", mit Begründung** — die Knappheit,
  die das Risiko benennt, ist mit über 16-fachem Sicherheitsabstand
  empirisch widerlegt, nicht nur einmalig, sondern über vier
  unterschiedliche Commits hinweg konsistent.
- **Risiko 3 — ohne Required-Status-Check könnte ein dauerhaft roter
  `e2e.yml`-Lauf unbemerkt bleiben (Alarmmüdigkeits-Risiko, analog zu
  `ADR-0051`s advisory-Läufen).** Dieses Risiko wird durch die vier grünen
  Läufe **nicht** berührt — es ist ein strukturelles Dauerrisiko über die
  Zeit, kein einmalig entscheidbarer Zustand. Es kann durch diesen Slice
  allein nicht geschlossen werden. Meine Einschätzung: **„weiter offen"**
  — gehört damit laut Modul 5 ins Beobachtungs-Register. Ich habe das
  Register durchsucht (`grep -rl "Alarmm" docs/plan/planning/observations/`):
  kein bestehender Eintrag deckt diese konkrete Klasse (nicht-blockierender
  PR-Trigger-Workflow, im Unterschied zu den bereits in `ADR-0051`
  benannten Nachtläufen CVE-Scan/Pin-Freshness). Vorschlag an den Planner:
  neues Verzeichnis, z. B. `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit/`,
  mit `evidence/slice-056.md` als erstem Beleg (1×) — kein Folge-Slice
  nötig, solange der Zähler unter 3× bleibt.

## 6. Hard Rules

- **3.1 (Docker-only):** Beide Workflow-Steps rufen ausschließlich
  bestehende `make`-Targets auf (`make image`, `make test-integration`) —
  selbst in der Datei gelesen, kein Inline-Shell.
- **3.3 (git mv + Inhaltsänderung = zwei Commits):** `b3f0793`
  (`next→in-progress`) ist ein reiner Move-Commit, `9666725`/`f2552e6`
  ebenso reine Vorbereitungs-Commits ohne Move+Inhalt-Vermischung; im
  geprüften Implementierungs-Commit (`e6aa155`) kein `git mv` — nicht
  einschlägig.
- **3.7 (Kommentar-Disziplin/Slice-Chronik-Verbot):** Eigene Lektüre des
  Workflow-Kommentars: begründet mit `ADR-0051`/`LH-QA-POR-003`/`ADR-0044`,
  zitiert `BEO-PGC/github-actions-unverifizierbar-lokal` als offene
  Beobachtung, Indikativ über den Ist-Zustand, keine Slice-/Wellen-Chronik
  (`slice-056` kommt im Workflow-Kommentar selbst nicht vor — eigener
  `grep -n "slice-" .github/workflows/e2e.yml` liefert keinen Treffer).
- **3.8 (Action-Pinning):** SHA-Pin gegen die GitHub-API selbst verifiziert
  (§2) — trägt.
- **3.6 (Gates nicht ohne ADR lockern):** nicht einschlägig — keine
  Schwellen-Änderung in diesem Slice, `make gates` unverändert.
- **3.9 (Exit-Code-Prüfung):** In dieser Sitzung durchgehend beachtet — der
  erste `make gates`-Lauf wurde versehentlich mit `tee` kombiniert
  (Exit-Code hätte dann den von `tee` widergespiegelt, nicht den von
  `make`); vor jeder Bewertung wiederholt **ohne** Pipe, mit separat
  geprüftem `$?` (siehe §2). Kein Befund landete auf Basis des
  pipe-behafteten Laufs.

## 7. Explizit NICHT geprüft (korrekt außerhalb dieser Rolle)

Die drei Paarungen (DoD-Punkt 11) — Slice liegt noch in `in-progress/`.
Beobachtungs-Register-Eintrag und Closure-Notiz (Planner-Closure-Arbeit,
beginnt erst nach diesem Bericht). Validierung gegen realen Bedarf (kein
Validator-Zug ausgelöst — reine CI-Infrastruktur, kein neuer
Architektur-Sicht-Meilenstein).

## Verdikt

**DoD-Konformität: bestätigt.** Alle vier real prüfbaren technischen Punkte
(1, 2, 4, 5) sind selbst reproduziert erfüllt (`make gates`, `gh run
list`/`gh run view --log`, direkte Datei-/API-Prüfung für Trigger/
Permissions/Pinning). Punkt 6 ist inhaltlich erfüllt, Checkbox korrekt noch
offen. Punkt 3 ist in der Substanz erfüllt — über einen strukturell
stärkeren Beleg (Branch-Protection-Abfrage) statt über den im DoD-Wortlaut
vorgesehenen Einzel-Testlauf; siehe §3 für die vollständige Begründung.
Die verbleibenden Punkte (7, 9–11, Reconciliation-Entfall bei 8) sind
korrekt noch offen — Planner-Closure-Arbeit nach Modul 8 §Rollen-Sequenz.

**Plan-vs-Code-Diff: keine unbegründete Abweichung.** Datei-Liste (§3 des
Plans) exakt getroffen, kein Out-of-Scope-Punkt (§1) berührt.

**§6-Risiken — meine Einschätzung, als Vorschlag an den Planner:**

- Risiko 1 (Runner-Unsicherheit): **entfallen** — vier reale, unabhängige
  grüne Läufe.
- Risiko 2 (Timeout-Knappheit): **entfallen** — durchgehend unter 7 % des
  60-Minuten-Budgets über vier Läufe.
- Risiko 3 (Alarmmüdigkeit bei stillem Dauer-Rot): **weiter offen** —
  neuer Beobachtungs-Registereintrag empfohlen (Vorschlag: `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit/`,
  1× mit `evidence/slice-056.md`).

**Übergabe an Planner:** Der Slice kann an die Planner-Closure übergeben
werden. Für die Closure-Notiz vorzumerken: die drei §6-Risiko-Ausgänge
(Vorschlag oben), der DoD-Punkt-3-Beleg-Weg (Konfigurationsbeleg statt
Einzel-Testlauf, §3), der neue Beobachtungs-Registereintrag für Risiko 3,
der bestehende Registereintrag `BEO-PGC/github-actions-unverifizierbar-lokal`
(Zähler auf 2× durch `evidence/slice-056.md`, weiterhin unter der
3×-Schwelle), und der `git mv` nach `done/` mit den drei Paarungen danach.
Kein Validator-Zug ausgelöst.

---

*Dieser Bericht ist ein Lauf-Beleg (Modul 11) und wird über Läufe hinweg
nicht wieder gelesen.*
