# Slice release-upstream-drift: Pin-Freshness über das Neun-Achsen-Inventar

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-release-pipeline-adr-0051](../welle-release-pipeline-adr-0051.md).

**Bezug:** [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
Entscheidung 7 (Upstream-Pin-Freshness), Pin-Inventar-Tabelle (P1–P9).

**Berührte Spec-Stellen:** — (Prozess-ADR ohne Spec-Stratum).

**Verantwortlich:** — (wellenlos priorisiert, siehe Welle-Auftrag).

**Autor:** Planner-Agent, direkt beauftragt. **Datum:** 2026-09-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das vollständige, in `ADR-0051` real ermittelte
Neun-Achsen-Pin-Inventar (P1–P9) beobachtbar machen: P1/P2 nutzt bereits
`make image-stale`, P3–P9 (Race-Toolchain, PG-Testcontainer, d-migrate,
a-check, d-check, Kurs-Baseline, GitHub-Action-Pins) bekommen neue,
analog benannte Make-Targets; `upstream-drift.yml` fragt alle neun Achsen
in einem fail-open Nachtlauf ab — ein Werkzeug-/Netzausfall einer
einzelnen Achse führt zu Skip dieser Achse, nicht zu Rot des gesamten
Laufs, ein tatsächlich gefundener Drift bleibt sichtbar.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Automatisches Anheben eines gefundenen Pins — bleibt bewusster Commit
  (dieselbe Disziplin wie das bestehende `make image-stale`,
  `ADR-0051` Entscheidung 7); der Nachtlauf meldet, hebt nicht.
- P9 (GitHub-Action-Pins) real gegen die Action-Releases prüfen, bevor
  `release.yml` (Slice `release-version-und-workflow`) überhaupt erste
  `uses:`-Zeilen enthält — die dort neu gepinnten Actions gehen in
  dasselbe P9-Target ein, sobald sie existieren; die Reihenfolge der
  beiden Slices innerhalb der Welle ist Planungsentscheidung (WIP-Limit
  1), P9 deckt am Ende ALLE `uses:`-Zeilen über alle Workflow-Dateien,
  unabhängig davon, wann sie entstanden.
- Ein blockierendes Gate aus einer der neun Achsen machen — `ADR-0051`
  Entscheidung 7/§Re-Evaluierungs-Trigger 5 verlangt dafür ausdrücklich
  eine eigene Folge-ADR (`AGENTS.md` §3.6, Gate-Verschärfung ist
  ADR-pflichtig).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [ ] Sieben neue, analog zu `make image-stale` benannte Make-Targets für
      P3–P9 existieren (Namensschema Implementierungsdetail, z. B.
      `make pin-stale-<achse>`): je Achse ein Tag-/Digest-Vergleich gegen
      den tatsächlichen Upstream-Stand (P3 Race-Toolchain, P4
      PG-Testcontainer, P5 d-migrate, P6 a-check, P7 d-check, P8
      Kurs-Baseline, P9 GitHub-Action-Pins über alle
      `.github/workflows/*.yml`), alle netzlos-unfähig (brauchen Netz,
      wie `make image-stale`), advisory (kein Gate).
- [ ] `upstream-drift.yml` existiert: ein Workflow, alle neun Achsen
      (P1/P2 über das bestehende `make image-stale`, P3–P9 über die neuen
      Targets), `if: always()` je Achsen-Schritt (fail-open — Werkzeug-/
      Netzausfall einer Achse führt zu Skip dieser Achse, nicht zu Rot des
      Gesamtlaufs), `schedule` (nächtlich) + `workflow_dispatch`; ein
      real gefundener Drift bleibt sichtbar (roter Lauf), bleibt aber
      advisory.
- [ ] `harness/README.md` §Werkzeuge trägt alle sieben neuen Targets und
      `upstream-drift.yml`.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für <Schnittstelle X> falls öffentlicher Vertrag berührt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

- Sieben neue Make-Targets, je Ziel-Pin ein Tag-/Digest-Vergleich analog
  `tools/harness/image-stale.sh`s Muster (`docker manifest inspect`/
  `docker buildx imagetools inspect` gegen den jeweils referenzierten
  Tag): P3 `TOOLCHAIN_RACE_IMAGE`, P4 `PG_TEST_IMAGE`, P5
  `D_MIGRATE_IMAGE`, P6 `A_CHECK_IMAGE`, P7 `DCHECK_IMAGE`/`DCHECK_DIGEST`
  (zwei Achsen — Tag-Frische **und** Digest-Drift, siehe `ADR-0051`
  Pin-Inventar-Tabelle), P8 Kurs-Baseline-Version in
  `harness/conventions.md` §Baseline (Release-API des Kurs-Repos, kein
  Docker-Pin), P9 alle `uses:`-SHA-Pins in `.github/workflows/*.yml`
  gegen die jeweilige Action-Repository-Releases.
- `upstream-drift.yml`: neun Achsen-Schritte, `if: always()` je Schritt,
  P1/P2 rufen `make image-stale` auf, P3–P9 die neuen Targets.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `Makefile` | update | sieben neue P3–P9-Targets. |
| `.github/workflows/upstream-drift.yml` | neu | neun Achsen, fail-open, `schedule` + `workflow_dispatch`. |
| `harness/README.md` §Werkzeuge | update | sieben neue Targets + `upstream-drift.yml`. |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): sofort — Welle eröffnet, unabhängig von
den übrigen Slices (siehe Welle-Datei §5 Abhängigkeiten).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): sieben
  strukturell unterschiedliche Achsen (Docker-Digests, Kurs-Release-API,
  GitHub-Action-SHAs) erweisen sich als zu heterogen für einen Slice —
  dann Aufteilung nach Achsen-Klasse (Docker-Pins P3–P7 vs.
  Nicht-Docker-Pins P8/P9).
- `in-progress` → `open` (blockiert — Carveout?): keine bekannte
  Blockade-Bedingung.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- `AGENTS.md` §3.10 gilt für `upstream-drift.yml` als neuen Workflow
  unverändert. **Ausgang:** weiter offen, strukturell (derselbe Fall wie
  `BEO-PGC/github-actions-unverifizierbar-lokal`, bereits verkörpert).
- P9 (GitHub-Action-Pins) deckt zum Zeitpunkt dieses Slice ggf. noch
  nicht die `uses:`-Zeilen aus `release.yml`/`image-scan.yml`/
  `hub-description.yml`, falls diese Slices noch nicht liefen — **Ausgang:**
  entfallen als Risiko: P9 ist ein Scan über den **jeweils aktuellen**
  Workflow-Baum, kein statisches Inventar; neue `uses:`-Zeilen werden ab
  ihrer Existenz automatisch mitgedeckt, keine Nacharbeit nötig.
- Ein fail-open-Design kann bei wiederholtem Werkzeugausfall (nicht
  Drift) unbemerkt bleiben, wenn niemand die Skip-Meldungen liest.
  **Ausgang:** weiter offen, dasselbe Muster wie
  `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit` (1×, noch unter
  der 3×-Schärfungsschwelle) — dieser Slice trägt keinen neuen Beleg
  (kein Required-Status-Check-Bezug, andere Fehlerklasse), aber die
  strukturelle Verwandtschaft wird hier benannt.

## 7. Closure-Notiz

*(wird bei Bearbeitung gefüllt.)*

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area „CI/CD-Pipeline"
(`.github/workflows/*.yml`, `Makefile`) — bereits mehrfach berührt,
Schwelle erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/github-actions-unverifizierbar-lokal` (7×, verkörpert) und
`BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit` (1×, siehe §6)
betreffen diese Sub-Area; kein weiterer Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
