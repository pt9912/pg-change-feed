# Review-Report: slice-056 — 2026-09-14

**Review-Art:** Code — geprüft gegen Plan + Konventionen (Modul 10 §Drei
Review-Arten).

**Gegenstand:** Commit `e6aa155` (Elter `b3f0793`, reiner
`next→in-progress`-Move) — `.github/workflows/e2e.yml` (neu),
`harness/README.md` (§Werkzeuge-Eintrag).

**Skill:** `.harness/skills/reviewer.md` @ Accepted (Stand 2026-09-13,
DoD-Checkbox-Nachzug-Abschnitt)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-14

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-056-e2e-workflow-in-ci.md` (§1 Ziel/Abgrenzung, §2 DoD, §4 Trigger, §6 Risiken)
- `ADR-0051` (`docs/plan/adr/0051-cicd-pipeline-github-actions.md`), Entscheidung 2
- `.github/workflows/ci.yml` (Vorbild-Workflow)
- `compose.yaml` (Kommentar zum fehlenden `build:`-Block, `ADR-0044`)
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentar-Disziplin), §3.8 (Action-Pinning), §3.9 (Exit-Code-Prüfung)
- `LH-QA-POR-003`

---

## Findings

### F-1 — DoD-Wortlaut „ausschließlich make test-integration" deckt den zweiten Schritt nicht ab

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `docs/plan/planning/in-progress/slice-056-e2e-workflow-in-ci.md:89-90`
- `befund`: Der DoD-Punkt formuliert „Ein Job ruft ausschließlich
  `make test-integration` auf" — die Implementierung ruft zusätzlich zuerst
  `make image` auf. Die Parenthese direkt danach („kein Inline-Shell, das das
  Target umgeht") macht die eigentliche Absicht klar (Make-Target-Zwang statt
  Ein-Schritt-Zwang), aber der Wortlaut „ausschließlich" liest sich isoliert
  als Ein-Target-Aussage und weicht vom tatsächlich gelieferten Umfang ab.
  Weder §1 (Ziel) noch §3 (Plan-Tabelle) des Slice-Plans nennen den
  `make image`-Schritt explizit, obwohin er (siehe F-INFO-1) technisch
  zwingend ist.
- `verifizierbar`: nein — Wortlaut-Frage, kein Gate-Lauf entscheidet das.
- `klasse`: „DoD-Wortlaut deckt gelieferten Umfang nicht exakt"

## Negativbefunde

- geprüft, ohne Befund: Action-Pinning — `actions/checkout@3d3c42e5aac5ba805825da76410c181273ba90b1 # v7.0.1` ist identisch zum Pin in `ci.yml` (AGENTS.md §3.8 erfüllt).
- geprüft, ohne Befund: Trigger-Form — `pull_request` + `push`/`branches: ['**']`/`tags-ignore: ['**']` ist zeichengleich zu `ci.yml` übernommen.
- geprüft, ohne Befund: Permissions — `permissions: {}` auf Workflow-Ebene, `contents: read` auf Job-Ebene; kein weiteres Recht (z. B. `packages: write`) angefordert — korrekt, denn `make image` läuft mit `docker buildx build --load` (Makefile:27), nicht `--push`, braucht also keine Registry-Schreibrechte.
- geprüft, ohne Befund: Kein Inline-Shell — beide Steps sind reine `run: make <target>`-Aufrufe, kein Bypass bestehender Ziele (AGENTS.md §3.1).
- geprüft, ohne Befund: Reihenfolge `make image` vor `make test-integration` — technisch notwendig, eigenständig verifiziert: `compose.yaml:1-6`/`:60-63` dokumentiert explizit, dass die Datei bewusst **keinen** `build:`-Block trägt, weil ein Zweit-Build das geladene `:dev`-Tag vom Lauf-Beleg in `harness/image-hash.txt` entkoppeln würde (`ADR-0044`); der reale Lauf (Run `34788082848`) bestätigt das zusätzlich — der `Image bauen`-Step schloss erfolgreich ab, bevor `Compose-Integrationstest` startete. Die Implementer-Begründung trägt, nicht nur behauptet.
- geprüft, ohne Befund: Timeout 60 Minuten — plausibel gegenüber `ci.yml`s 30 Minuten; `tools/harness/run-integration-tests.sh` trägt allein aus festen `sleep`-Aufrufen (u. a. zweimal `sleep 25`, mehrfach `sleep 1`–`sleep 5`) mehrere Minuten zusätzliche Wartezeit gegenüber reinen Unit-Tests, plus Compose-Hochfahren, Image-Pulls und Schema-Rollout — 60 Minuten ist ein sinnvoller Sicherheitsabstand, keine willkürliche Verdopplung.
- geprüft, ohne Befund: `harness/README.md` — der neue Eintrag landete in der **Werkzeuge**-Tabelle (Zeile nach `make test-integration`), nicht in der Sensors-Tabelle — korrekt, da `e2e.yml` kein Gate ist (`kein Gate, ADR-0051, LH-QA-POR-003` in der Bindung-Spalte).
- geprüft, ohne Befund: Kommentar-Disziplin (AGENTS.md §3.7) — der neue Workflow-Kommentar begründet mit `ADR-0051`/`LH-QA-POR-003`/`ADR-0044` und zitiert `BEO-PGC/github-actions-unverifizierbar-lokal` als offene Beobachtung; keine Slice-/Wellen-Chronik (`slice-056` kommt im Workflow-Kommentar selbst nicht vor).
- geprüft, ohne Befund: Commit-Trailer — kein `Co-Authored-By:`, kein `Claude-Session:`.
- geprüft, ohne Befund: Traceability — Betreff nennt `LH-QA-POR-003` und `ADR-0051`, keine `SPEC-*`/`ARC-*`-ID im Betreff.
- geprüft, ohne Befund: `make gates` — lokal ausgeführt, Exit-Code direkt geprüft (`$?`), Ergebnis `EXIT_CODE=0` (baseline-verify, docs-check/d-check 0 Befunde, commit-traceability OK, a-check 0 Befunde).
- geprüft, ohne Befund: §1 Out-of-Scope-Disziplin, §6 Risiken, §8 Sub-Area-Begründung des Slice-Plans — vier Klassen mit Begründung, drei Risiken vorab benannt, Sub-Area-Sichtung (`BEO-PGC/github-actions-unverifizierbar-lokal`, Stand 1×) dokumentiert; keine Abweichung im Diff feststellbar (dieser Commit ändert nur Code, nicht den Plan — Risiko-Ausgänge sind Closure-Aufgabe, nicht Gegenstand dieses Reviews).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** DoD-Wortlaut deckt gelieferten Umfang nicht exakt

## Live-Status des echten `e2e`-Workflow-Laufs

Erster realer Lauf nach dem Push von `e6aa155` — bislang nur statische
Prüfung möglich (`BEO-PGC/github-actions-unverifizierbar-lokal`, Stand vor
diesem Slice: 1×).

- **Run:** `34788082848` (Trigger: Push von `e6aa155` auf `main`), Job
  `image + test-integration` (ID `103807124807`).
- **Zwischenstand während der Erstfassung dieses Reports (2026-09-14, ca.
  22:56 UTC, Laufzeit ~3–4 Min.):** `Checkout` ✓, `Image bauen (Lauf-Beleg
  harness/image-hash.txt)` ✓ erfolgreich abgeschlossen, `Compose-
  Integrationstest (Black-Box-E2E)` lief zu diesem Zeitpunkt noch — bewusst
  nicht abgewartet, nicht geraten.
- **Endergebnis (nachgetragen im selben Review-Lauf, 2026-09-14, ca.
  22:57 UTC):** Run `34788082848` ist `completed`/`success` in 3m44s —
  deutlich unter dem 60-Minuten-Timeout. Der Job-Log zeigt alle erwarteten
  Testphasen real durchlaufen (Rollen-DSN, Retention-Lebenszyklus,
  SQL-Administration-Live-Reload, Publication-Entzug-Wirksamkeit,
  Black-Box-CLI-Rundlauf, drei NATS-Testabschnitte, Schema-Change-Beleg) und
  endet mit `ok  	.../test/integration	0.143s`. Das ist der erste reale,
  grüne Beleg für `e2e.yml` — `BEO-PGC/github-actions-unverifizierbar-lokal`
  ist für diesen Workflow damit aufgelöst, nicht mehr nur statisch geprüft.
- **Zusatzbeobachtung:** Der parallele `ci`-Lauf auf demselben Commit
  (`34788082830`) war bereits `completed`/`success` in 1m36s, deutlich vor
  dem `e2e`-Endergebnis — das belegt real, dass der neue `e2e`-Workflow
  `ci.yml`s kurzes Signal tatsächlich nicht verzögert (getrennte Workflows,
  wie im Plan vorgesehen).
- **Verifikations-Konsequenz:** Das DoD-Item „real ausgelöst und sichtbares
  Check-Ergebnis" ist mit diesem grünen Lauf erfüllt. Weiterhin offen für den
  Verifier (nicht Gegenstand dieses Reviews): der separate, nicht
  einzucheckende Rot-Beleg für die Nicht-Blockierung (DoD-Punkt 3) sowie die
  Risiko-Ausgänge aus §6 des Slice-Plans.

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, ein LOW ohne
Rollen-Widerspruch (Wortlaut-Präzisierung, die der Planner bei Gelegenheit
nachziehen kann, keine Fixrunde am Implementer nötig).

**Übergabe:** Keine Reviewer→Implementer-Rückgabe nötig (siehe
`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde) — die
DoD-Zeile „Review durchgeführt" wird in diesem Commit selbst nachgezogen.
Der `e2e`-Lauf schloss während dieses Review-Laufs grün ab (§Live-Status
oben) — der Verifier prüft davon unabhängig weiterhin den separaten
Rot-Beleg (Nicht-Blockierung) und die Risiko-Ausgänge aus §6 vor Closure;
dieser Report ist ein Lauf-Beleg und wird über Läufe hinweg nicht erneut
gelesen. Kein Ersatz für Verifikation gegen DoD/Spec (Modul 11).
