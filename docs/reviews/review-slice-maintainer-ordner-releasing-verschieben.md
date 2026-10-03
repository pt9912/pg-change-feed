# Review-Report: slice-maintainer-ordner-releasing-verschieben — 2026-10-03

**Review-Art:** Code (Doku-Diff) gegen Plan, Entscheidungen und Hard Rules

**Gegenstand:** `git diff de24629d HEAD` — Commits `b26460a3` (reiner Move), `13b6919b` (Verweise, Index, Gate), `36b311b4` (Zitat-Korrektur der [`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md)), `2d78f702` (Reviews, Register), `0e7bbb6c` (Suchlauf-Messung im Plan). Der Planner-Commit `ca3f04fe` (SDK-Plan) liegt im Bereich und ist fremd; er ist nicht Gegenstand.

**Skill:** `.harness/skills/reviewer.md` · **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-03

**Eingangs-Kontext:**

- Slice-Plan `slice-maintainer-ordner-releasing-verschieben` (in-progress)
- [`ADR-0143`](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md), [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md), [`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md), [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)
- `AGENTS.md` §3.1, §3.3, §3.5, §3.7, §3.9, §3.11, §3.12, §3.13, §4
- Auftraggeber-Entscheidungen: Rang-6-Wortlaut freigegeben, `releasing.md` ungeteilt verschoben, `version.md` bleibt

---

## Findings

### F-1 — ADR-0123: Linktext nennt den alten Pfad, das Linkziel den neuen

- `kategorie`: LOW
- `quelle`: [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) Entscheidung 1 (Verweisgerüst), `AGENTS.md` §3.7 (Ist-Zustand)
- `pfad`: `docs/plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md:48`, `:139`, `:304`
- `befund`: Drei Links tragen den sichtbaren Text `docs/user/releasing.md` bei Linkziel `../../maintainer/releasing.md`; Lesende der ADR ohne `git` sehen einen Pfad, den es nicht mehr gibt. [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) nennt „Linkziele", sichtbare Linktexte sind Verweisgerüst derselben Klasse, und `docs-check` fängt die Diskrepanz nicht (Ziel löst auf).
- `verifizierbar`: nein (kein Gate liest Linktext gegen Ziel)
- `klasse`: Zitat-Korrektur halb gezogen

### F-2 — Plan-Prosa zählt Review-Links uneinheitlich

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 Instanz A
- `pfad`: Slice-Plan §1, Absatz „Messung am Stand `diff`"
- `befund`: Der Absatz erklärt, der Plan habe 3 Review-Links gezählt, tatsächlich sind es 2 Linkziel-Änderungen. Die Suchlauf-Zeile `diff 3 -E \]\([^)]*releasing\.md -- docs/reviews` zählt zusätzlich den Dateinamen-Link auf `review-slice-release-doku-releasing.md`. Die Zahlen stimmen (nachgemessen), die Erklärung ist nur schwer lesbar. Der Plan nennt 13 Suchlauf-Zeilen im Auftrag, `make suchlauf-nachmessen` meldet 12 (alle stimmen).
- `verifizierbar`: ja (`make suchlauf-nachmessen`)
- `klasse`: Zahl im Träger, lesbar aber nicht falsch

### F-3 — Hinweis für den Planner: SDK-Plan nennt den alten Pfad (kein Befund des Slice)

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.13 (fremde Datei wird gemeldet)
- `pfad`: `docs/plan/planning/open/slice-sdk-meldungscodes-in-fehlertypen.md:70`, `:386`
- `befund`: Zwei Klartext-Nennungen von `docs/user/releasing.md`; der Plan wurde nach dem Start dieses Slice angelegt und ist im Slice-Plan als gemeldet geführt (Frist Closure). Adressat: Planner.
- `verifizierbar`: ja (`git grep -n "docs/user/releasing"`)
- `klasse`: Träger-Nachzug fremde Datei

### F-4 — ADR-0143 Festlegung 2 nennt `docs/user/releasing.md` als ausgenommen (Klartext, unverändert)

- `kategorie`: INFO
- `quelle`: [`ADR-0143`](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md) Z. 87
- `pfad`: `docs/plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md:87`
- `befund`: Die Accepted-ADR beschreibt ihren Stand (Ausnahme in `docs/user/`); das Gate führt `releasing.md` nicht mehr, die Datei liegt räumlich außerhalb. Die Aussage „dauerhaft außerhalb" ist räumlich bestätigt, kein Lockern. Unveränderlichkeit gewahrt, Klartext-Rest als Record akzeptabel.
- `verifizierbar`: nein
- `klasse`: Record beschreibt früheren Stand

---

## Prüfpunkte (a) bis (j)

- **(a) Move rein:** `git diff --stat -M de24629d..b26460a3` = `docs/{user => maintainer}/releasing.md | 0`, 1 Datei, 0 Einfügungen, 0 Löschungen. `git log --follow` verfolgt die Datei über `b26460a3` hinaus bis `603413c5`. AGENTS §3.3 gewahrt.
- **(b) releasing.md:** Diff gegen den Parent zeigt genau Titel, Version 1.13, Stand, §1 (Zielgruppe), Link `../user/version.md` und die Historienzeile 1.13. Die Zeile trägt `ADR-0143`/`ADR-0051` wie ältere Zeilen (1.7 bis 1.12), Maintainer-Doku ist nicht Gegenstand des Handbuch-Gates (Mutation 3 unten). `docs/maintainer/README.md`: kurz, Index korrekt, Link auf `../user/` löst auf.
- **(c) AGENTS.md / harness/README.md:** `AGENTS.md` genau eine Zeile (Z. 66), Wortlaut wie freigegeben. `harness/README.md` Rang-6-Zeile gleicher Inhalt in Tabellenform. `git grep -i releasing -- AGENTS.md harness/README.md`: nur diese beiden Zeilen, keine Restnennung von `docs/user/releasing.md`. Zusätzlich geändert und im Plan benannt: `.harness/skills/nutzerdoku-schreiben.md` Z. 13.
- **(d) Wurzel-READMEs:** beide Links auf `docs/maintainer/releasing.md`. Der Hinweis zur Docker-Hub-Beschreibung (spiegelt `README.md` erst beim nächsten Release) ist plausibel: `hub-description.yml` läuft nur aus `release.yml` oder manuell.
- **(e) [`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md):** 3 Linkziele plus genau eine neue §Geschichte-Zeile mit Commit `b26460a3`; Entscheidung und Konsequenzen unberührt (`git diff -U0`, Zeile für Zeile). Klartext in Z. 275 bleibt. Sichtbarer Linktext: siehe F-1. [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md) und [`ADR-0143`](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md) unverändert (`git diff --name-only` ohne beide), Klartext-Rest akzeptabel.
- **(f) Reviews/Register:** beide Review-Dateien nur Linkziel. `state.md`: Zustandsfeld nennt Adresse `docs/maintainer/releasing.md`, keine Chronik, Anker löst auf.
- **(g) Handbuch-Gate:** `excluded` ohne `releasing.md`; Tabellentest 41 Fälle. Mutationen einzeln auf Kopien (`git archive` im Scratchpad, Mutation per `awk` nach Kopie, Echtrepo `git status --short` leer):
  1. `releasing.md` wieder in `excluded`: Gate Exit 2, „genannte Datei fehlt: docs/user/releasing.md".
  2. `bench-abdeckung.md` aus `excluded`: Tabellentest Exit 1, mehrere Fälle rot (u. a. „Kennung in bench-abdeckung.md ausgenommen", „fehlende ausgenommen-genannte Datei").
  3. Kennung an `docs/maintainer/releasing.md`: Gate Exit 0, unberührt.
  Sensor-Vertrag und Gate-Zeile in `harness/README.md` sind gegen das Skript wahr (drei geprüfte Dateien, vier Erzeugnisse ausgenommen, Exit 2 bei fehlender oder unklassifizierter Datei).
- **(h) `.d-check.yml`:** `git grep releasing -- .d-check.yml` leer, unverändert, `docs-check` Exit 0. `git grep -n "docs/user/releasing"` im aktiven Baum ohne Records: Treffer nur [`ADR-0051`](../plan/adr/0051-cicd-pipeline-github-actions.md), [`ADR-0123`](../plan/adr/0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md) (Klartext), [`ADR-0143`](../plan/adr/0143-handbuch-public-doc-check-gate-kennungsfreie-nutzerdoku.md), Plan selbst, Register (Klartext) und der SDK-Plan (F-3). `.github/**`, `Makefile`, `harness/mk`, `tools/**`, `.claude`, `.harness/skills`: kein Treffer auf den alten Pfad. Suchlauf: `make suchlauf-nachmessen PLAN=…` Exit 0, 12 Zeilen stimmen.
- **(i) Läufe (Exit direkt):** `make docs-check` 0 (1614 Dateien, 0 Befunde), `make handbuch-public-doc-check` 0, `make test-handbuch-public-doc-check` 0 (41 Fälle), `make sdk-public-doc-check` 0, `make gates` 0 (ungepiped, Exit in Datei gesichert), `make doc-immutable RANGE=de24629d..HEAD` 0, `make doc-commits RANGE=de24629d..HEAD` 0, `make commit-traceability` 0, `make fmt-check` 0, `make kommentar-kennungen DIFF=de24629d` 0.
- **(j) Reichweite:** `git diff --name-only de24629d HEAD -- .github docs/user .d-check.yml Makefile harness/mk` listet nur die verschobene `docs/user/releasing.md`; `version.md` unverändert, kein Release, kein Tag, `release.yml` unverändert.

## Geprüft, ohne Befund

- `docs/maintainer/`: ohne Befund
- `AGENTS.md`, `harness/README.md`, `README.md`, `README.de.md`: ohne Befund
- `tools/harness/handbuch-public-doc-check.sh` und Tabellentest: ohne Befund (Kommentar-Probe leer)
- `harness/sensors/handbuch-public-doc-check.md`: ohne Befund
- `docs/reviews/` (zwei Dateien) und Register `state.md`: ohne Befund
- `.github/`, `Makefile`, `harness/mk`, `.d-check.yml`: unverändert

## Architect-Fragen

1. F-1: Genügt für eine Zitat-Korrektur nach [`ADR-0073`](../plan/adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) das Linkziel, oder gehört der sichtbare Linktext (Verweisgerüst) in dieselbe Korrektur? Die Rückfrage ist nur nötig, falls F-1 nicht ohne Weiteres als weitere Zitat-Korrektur (mit einer weiteren §Geschichte-Zeile) gezogen werden soll.

## Verdikt

0 HIGH, 0 MEDIUM, 1 LOW, 3 INFO. Merge-blockierend: **nein**. Keine Fixrunde am Implementer nötig; F-1 kann der Planner oder Architect ohne Rückgabe-Pfeil entscheiden. Die DoD-Zeile „Review durchgeführt" im Plan ist mit diesem Report auf `[x]` gezogen.
