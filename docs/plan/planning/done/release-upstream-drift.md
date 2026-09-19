# Slice release-upstream-drift: Pin-Freshness über das Neun-Achsen-Inventar

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-release-pipeline-adr-0051](welle-release-pipeline-adr-0051.md).

**Bezug:** [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md)
Entscheidung 7 (Upstream-Pin-Freshness), Pin-Inventar-Tabelle (P1–P9).

**Berührte Spec-Stellen:** — (Prozess-ADR ohne Spec-Stratum).

**Verantwortlich:** Implementer-Agent (priorisiert 2026-09-19, direkter
Auftrag: "dann mach weiter bis die welle geschlossen ist").

**Autor:** Planner-Agent, direkt beauftragt. **Datum:** 2026-09-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das vollständige, in `ADR-0051` real ermittelte
Neun-Achsen-Pin-Inventar (P1–P9) beobachtbar machen: P1/P2 nutzt bereits
`make image-stale`, P3–P9 (Race-Toolchain, PG-Testcontainer, d-migrate,
a-check, d-check, Kurs-Baseline, GitHub-Action-Pins) bekommen neue,
analog benannte Make-Targets; `upstream-drift.yml` fragt alle neun Achsen
in einem fail-open Nachtlauf ab — `if: always()` verhindert, dass ein
fehlgeschlagener Schritt (echter Drift-Fund **oder** ein Werkzeug-/
Netzausfall) die nachfolgenden Achsen-Schritte abbricht, sodass jede
Achse unabhängig sichtbar bleibt. Reviewer-Finding F-2 (siehe §6)
präzisiert: das unterscheidet nicht zwischen den beiden Ursachen auf
Lauf-Statusebene — GitHub Actions färbt Schritt und Gesamtlauf bei jedem
Nicht-Null-Exit gleich rot.

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

- [x] Sieben neue, analog zu `make image-stale` benannte Make-Targets für
      P3–P9 existieren (`make pin-stale-race`/`-pgtest`/`-dmigrate`/
      `-acheck`/`-dcheck`/`-baseline`/`-actions`): je Achse ein
      Tag-/Digest-Vergleich gegen den tatsächlichen Upstream-Stand (P3
      Race-Toolchain, P4 PG-Testcontainer, P5 d-migrate, P6 a-check, P7
      d-check — zwei Achsen, P8 Kurs-Baseline, P9 GitHub-Action-Pins über
      alle `.github/workflows/*.yml`), alle advisory (kein Gate), brauchen
      Netz. Real ausgeführt (nicht nur `make -n`): P3 (`golang:1.27`), P4
      (`postgres:18-alpine`) und P5 (`ghcr.io/pt9912/d-migrate:latest`)
      zeigen echten, unbekannten Drift (Upstream-Tags wurden seit dem Pin
      neu gebaut) — kein Fehler dieses Slice, sondern der erste reale
      Fund, den dieses Werkzeug liefern soll; P6/P7/P8/P9 sind aktuell
      (`OK`).
- [x] `upstream-drift.yml` existiert: ein Workflow, alle neun Achsen
      (P1/P2 über das bestehende `make image-stale`, P3–P9 über die neuen
      Targets), `if: always()` je Achsen-Schritt (fail-open — ein
      fehlgeschlagener Schritt bricht die nachfolgenden Achsen nicht ab,
      jede bleibt unabhängig sichtbar; `if: always()` unterscheidet dabei
      NICHT zwischen einem echten Drift-Fund und einem Werkzeug-/
      Netzausfall auf Lauf-Statusebene — Reviewer-Finding F-2, §6),
      `schedule` (nächtlich, versetzt zu `image-scan.yml`)
      + `workflow_dispatch`; ein real gefundener Drift bleibt sichtbar
      (roter Schritt), bleibt aber advisory. YAML-Struktur per
      Ruby-Stdlib-YAML-Parser geprüft.
- [x] `harness/README.md` §Werkzeuge trägt alle sieben neuen Targets und
      `upstream-drift.yml`.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      1 HIGH (F-1: drei Skripte riefen `curl` bare auf dem Host auf statt
      containerisiert, `AGENTS.md` §3.1) und 1 MEDIUM (F-2: `if: always()`
      unterscheidet nicht zwischen Drift-Fund und Werkzeugausfall auf
      Lauf-Statusebene) sowie 2 LOW (F-3: Exit-Code-Vertrag von
      `pin-stale.sh` bei fehlendem Argument; F-4: doppelte Plan-Tabellenzeile)
      und 1 INFO (F-5: P9 prüft nur Frische, keine Form) in derselben
      Fixrunde behoben bzw. dokumentiert, kein offenes HIGH.
- [ ] Doku-Update für `harness/README.md` §Werkzeuge — entfällt als
      eigener Punkt, da bereits §2 oben dieselbe Zeile explizit als
      DoD-Kriterium trägt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine Reconciliation-Datei in diesem Repo (kein Brownfield-Bootstrap).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — kein neuer Eintrag: F-2 ist eine andere Fehlerklasse als `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit` (kein Required-Status-Check-Bezug, siehe §6), keine weitere Beobachtung angefallen.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Dieser Slice gehört zu `welle-release-pipeline-adr-0051` (noch offen) — Prüfung folgt regelkonform bei deren Closure.

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
| `Makefile` | update | sieben neue P3–P9-Targets, alle `bash tools/harness/pin-stale*.sh`-Aufrufe. |
| `tools/harness/pin-stale.sh` | neu (Plan-Nachzug) | gemeinsames Skript für P3–P6 (Digest-Pin in Makefile-/`a-check.mk`-Variable, `docker buildx imagetools inspect` gegen Tag oder explizites drittes Argument bei tag-losem Digest-Pin — P5/P6). |
| `tools/harness/pin-stale-dcheck.sh` | neu (Plan-Nachzug) | P7, zwei Achsen (Digest-Drift über `pin-stale.sh` + Tag-Frische über die GitHub-Releases-API von `pt9912/d-check`) in einem Skript, weil `DCHECK_IMAGE`/`DCHECK_DIGEST` zwei getrennte Variablen sind. |
| `tools/harness/pin-stale-baseline.sh` | neu (Plan-Nachzug) | P8, `harness/conventions.md` §Baseline gegen die GitHub-Releases-API von `pt9912/ai-harness-course` (kein Docker-Pin, eigenes Parsing). |
| `tools/harness/pin-stale-actions.sh` | neu (Plan-Nachzug) | P9, scannt alle `uses:`-Zeilen über `.github/workflows/*.yml`, zwei Achsen je Repo@Tag (`git ls-remote` gegen Tag-Mutation, GitHub-Releases-API gegen Tag-Frische), dedupliziert Mehrfachnennungen. |
| `.github/workflows/upstream-drift.yml` | neu | neun Achsen, fail-open, `schedule` + `workflow_dispatch`. |
| `harness/README.md` §Werkzeuge | update (Plan-Nachzug) | acht neue Zeilen (vier P3–P6-Targets gebündelt, P7/P8/P9 einzeln, `upstream-drift.yml`) — ersetzt die ursprünglich grob geplante einzelne Zeile. |
| `tools/harness/lib-github-api.sh` | neu (Plan-Nachzug, Fixrunde) | Reviewer-Finding F-1 (HIGH): gemeinsame containerisierte GitHub-API-GET-Funktion, analog `tools/harness/ci-matrix-abdeckung.sh` (`AGENTS.md` §3.1) — `pin-stale-baseline.sh`/`pin-stale-dcheck.sh`/`pin-stale-actions.sh` riefen `curl` zuvor bare auf dem Host auf. |

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
- Reviewer-Finding F-2 (MEDIUM): `if: always()` liefert nicht die
  ursprünglich hier behauptete Eigenschaft „Werkzeug-/Netzausfall führt
  zu Skip dieser Achse, nicht zu Rot des Gesamtlaufs" — GitHub Actions
  unterscheidet auf Lauf-Statusebene nicht zwischen Exit 1 (echter Fund)
  und Exit 2 (unbestimmbar); beide färben Schritt und Gesamtlauf gleich
  rot. `continue-on-error: true` wäre keine Abhilfe (verschluckt echte
  Funde ebenso). **Ausgang:** eingetreten, Doku korrigiert (§1, §2 oben,
  `upstream-drift.yml`-Kopfkommentar) — kein Code-Fix möglich, ohne echte
  Funde zu verschlucken; die genauere Fassung des bereits benannten
  Alarmmüdigkeits-Risikos (siehe Bullet direkt oberhalb) bleibt weiter
  offen.
- Reviewer-Finding F-5 (INFO): `pin-stale-actions.sh` prüft nur die
  Frische bereits §3.8-konformer `uses:`-Zeilen — eine künftige,
  §3.8-verletzende Zeile (floatender Tag, SHA ohne Versionskommentar)
  erzeugt weder `DRIFT` noch `UNBESTIMMT`, sondern bleibt für P9
  unsichtbar; die Form-Konformität selbst bleibt Review-Aufgabe, kein
  Gate deckt sie. **Ausgang:** entfallen als Risiko dieses Slice — P9 war
  nie als Form-Prüfung geplant (§1 Ziel: „Pin-Freshness", nicht
  „Pin-Form"), aktuell real folgenlos (alle neun `uses:`-Zeilen
  §3.8-konform, siehe Review-Report Negativbefunde).

## 7. Closure-Notiz

- **Was hat funktioniert:** Der reale, hands-on-Testlauf aller sieben
  neuen `make pin-stale-*`-Targets vor der Übergabe lieferte echte,
  unbekannte Funde (P3 `golang:1.27`, P4 `postgres:18-alpine`, P5
  `ghcr.io/pt9912/d-migrate:latest` — alle drei zeigen Upstream-Drift,
  der vorher nirgends sichtbar war) — ein wertvolles Nebenprodukt des
  Slice über seinen eigentlichen Liefergegenstand hinaus. Ebenso die
  bewusste Wiederverwendung des etablierten `tools/harness/
  ci-matrix-abdeckung.sh`-Musters für die Fixrunde (`lib-github-api.sh`)
  statt einer neuen, abweichenden Lösung.
- **Was ging anders als geplant:** Der Reviewer fand 1 HIGH (F-1: drei
  neue Skripte riefen `curl` bare auf dem Host auf, im Widerspruch zum
  bereits etablierten Container-Muster im selben Verzeichnis —
  `AGENTS.md` §3.1) und 1 MEDIUM (F-2: die im Plan/
  [`ADR-0051`](../../adr/0051-cicd-pipeline-github-actions.md) übernommene
  Formulierung „Werkzeug-/Netzausfall führt zu Skip, nicht zu Rot des
  Gesamtlaufs" hält für die Lauf-Statusanzeige nicht — `if: always()`
  unterscheidet nicht zwischen einem echten Fund und einem transienten
  Ausfall) sowie 2 LOW (F-3: `pin-stale.sh`s Exit-Code-Vertrag stimmte
  bei fehlendem Argument nicht; F-4: doppelte Plan-Tabellenzeile) und
  1 INFO (F-5: P9 prüft nur Frische, keine Form-Konformität). Alle fünf
  in einer Fixrunde behoben bzw. dokumentiert (`eb4e7d3f`) — F-2 ließ
  sich nicht code-seitig beheben, ohne echte Funde zu verschlucken
  (`continue-on-error: true` wäre keine Abhilfe), deshalb Korrektur der
  Doku-Behauptung statt eines Mechanismus-Fixes.
- **Steering-Loop-Eintrag:** kein neuer Sensor, keine geschärfte Regel —
  F-1 bestätigt das bereits etablierte Muster (Container-Kapselung für
  Netzwerkwerkzeuge in `tools/harness/*.sh`, `AGENTS.md` §3.1) und macht
  es über `lib-github-api.sh` für drei weitere Skripte wiederverwendbar,
  statt es dreifach zu duplizieren. F-2 ist eine genauere technische
  Fassung eines bereits bekannten, aber anderer Fehlerklasse zugehörigen
  Alarmmüdigkeits-Risikos (`BEO-PGC/nicht-blockierender-workflow-
  alarmmuedigkeit`, kein Required-Status-Check-Bezug hier, siehe §6) —
  kein neuer Registereintrag, da die Klasse nicht deckungsgleich ist.
- **Beobachtungs-Register (`../observations/`):** kein neuer Eintrag
  (siehe Steering-Loop-Eintrag oben und §6).
- **Folge-Slices:** `release-hub-description`, `release-doku-releasing`
  — beide bereits als Dateien in `open/` vorhanden (Welle
  `welle-release-pipeline-adr-0051`). Kein neuer Folge-Slice aus diesem
  Slice selbst.
- **Risiken aus §6:** fünf Risiken, fünf Ausgänge — (1) `AGENTS.md` §3.10
  realer Post-Push-Lauf → **weiter offen**, bereits verkörpert; (2) P9
  deckt künftige `uses:`-Zeilen automatisch → **entfallen**; (3)
  Alarmmüdigkeit bei wiederholtem Werkzeugausfall → **weiter offen**,
  strukturell verwandt mit `BEO-PGC/nicht-blockierender-workflow-
  alarmmuedigkeit` (1×, unter der Schärfungsschwelle); (4) F-2 „if:
  always() unterscheidet nicht zwischen Fund und Ausfall" → **eingetreten**,
  Doku korrigiert, kein Code-Fix möglich; (5) F-5 „P9 prüft nur Frische,
  keine Form" → **entfallen**, war nie als Form-Prüfung geplant.
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-release-pipeline-adr-0051](welle-release-pipeline-adr-0051.md)
  (noch offen) — die Prüfung läuft regelkonform bei deren Closure, nicht
  hier (§2 DoD-Zeile „im Repo mit Wellen von der nächsten
  Welle-Closure"). Vorab-Hinweis für diese spätere Prüfung: kein
  `liegt in`-Feld in diesem Slice; beide Folge-Slices existieren bereits
  als Dateien unter `docs/plan/planning/open/`; kein neuer
  Beobachtungs-Beleg in diesem Slice.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area „CI/CD-Pipeline"
(`.github/workflows/*.yml`, `Makefile`) — bereits mehrfach berührt,
Schwelle erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/github-actions-unverifizierbar-lokal` (7×, verkörpert) und
`BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit` (1×, siehe §6)
betreffen diese Sub-Area; kein weiterer Treffer.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
