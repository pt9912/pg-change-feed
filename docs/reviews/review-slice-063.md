# Review-Report: slice-063 — 2026-09-14

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10
§Drei Review-Arten), geprüft gegen Plan/ADR/Hard Rules (Maintainability),
**nicht** gegen die DoD (Verifier-Aufgabe, Modul 11).

**Gegenstand:** `slice-063` — Implementierungs-Commit `94d7530`, Elter
`265eda8` (reiner `next→in-progress`-Move). Zweiter
Implementierungs-Versuch nach einem real ausgeführten Blocker-Durchlauf
(`docs/reviews/blocker-slice-063.md`) und dessen Auflösung über
`ADR-0064` (Supersedes `ADR-0058`, nur Entscheidung 3). `git diff 265eda8
94d7530 --stat` bestätigt genau drei geänderte Dateien:
`docs/plan/planning/in-progress/slice-063-…md`, `harness/README.md`,
`tools/harness/run-integration-tests.sh` — kein Griff in eine Datei
außerhalb der im Slice-Plan §3 genannten.

**Skill:** `.harness/skills/reviewer.md` (Stand 2026-09-09, geschärft
2026-09-13: vier repo-spezifische HIGH-Regeln plus Slice-/Wellen-Chronik-
und Handbuch-Versionshistorie-Regel)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-14

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-063-upgrade-sicherheit-schema-rollout-zyklus.md`
  (vollständig gelesen, inkl. §4 Nachtrag, §6, §8)
- `docs/plan/adr/0064-lh-qa-ops-005-testansatz-korrektur.md` (vollständig
  gelesen — bindend für den korrigierten `LH-QA-OPS-005`-Testansatz)
- `docs/plan/adr/0058-testansatz-fuenf-luecken.md` (Entscheidung 3, als
  Kontext — für `LH-QA-OPS-005` durch `ADR-0064` superseded)
- `docs/reviews/blocker-slice-063.md` (vollständig gelesen — reale
  Erst-Reproduktion des Exit-8-Blockers)
- `AGENTS.md` §3 Hard Rules (insb. 3.1, 3.3, 3.5, 3.7, 3.9), §5
  Dokumentations-Regeln
- `harness/conventions.md` (MR-000 ID-Schema)
- vorherige Findings am gleichen Modul: `docs/reviews/review-slice-062.md`
  (unmittelbarer Vorgänger-Diff in derselben Datei,
  `tools/harness/run-integration-tests.sh`, keine offenen HIGH/MEDIUM)

---

## Findings

Keine HIGH-, MEDIUM- oder LOW-Findings. Der Diff hält sich strukturell
und inhaltlich an `ADR-0064`s exakt vorgegebene Testform und an den
Slice-Plan.

## Negativbefunde

- geprüft, ohne Befund: `tools/harness/run-integration-tests.sh` — neue
  Phase (Zeilen ~1554–1649) entspricht `ADR-0064`s vorgegebenem
  Mechanismus wortgetreu: `$COMPOSE up -d --force-recreate --no-deps
  pg-change-feed` ersetzt den `docker stop`/`make schema-rollout`/`docker
  start`-Dreischritt vollständig; kein Rest dieses Dreischritts im Diff.
- geprüft, ohne Befund: Platzierung der neuen Phase — sie liegt nach dem
  HTTP-API-Rundlauf und **vor** `TestE2ESchemaChangeDropColumn`
  (Zeile 1651ff.), nicht danach. Eigenständig nachvollzogen (nicht aus der
  Implementer-Begründung übernommen): `TestE2ESchemaChangeDropColumn`
  entfernt real eine Spalte und beendet den Erfassungspfad des
  Feed-Containers dauerhaft (`error_class=schema`, Replication-Slot wird
  neu angelegt, kein Replay) — ein Container-Tausch danach träfe auf
  einen bereits durch einen eigenen Recovery-Mechanismus (Slot-Neuanlage,
  Schema-Version-Nachtrag) rekonstruierten Zustand und würde mit diesem
  Mechanismus kollidieren bzw. ihn redundant machen. Der reale
  `make test-integration`-Lauf dieses Reviews bestätigt die Platzierung
  operational: Der Upgrade-Sicherheits-Rundlauf lief erfolgreich direkt
  vor `TestE2ESchemaChangeDropColumn`, und Letzterer bestand im Anschluss
  unverändert (`--- PASS: TestE2ESchemaChangeDropColumn (0.24s)`). Die
  Platzierungsentscheidung trägt.
- geprüft, ohne Befund: Container-ID-Wechsel real im Code verankert
  (`upgrade_feed_id_before`/`upgrade_feed_id_after` über `docker inspect
  --format '{{.Id}}'`, harter `exit 1` bei Gleichheit) und im eigenen
  `make test-integration`-Lauf real beobachtet: Feed-Container-ID wechselte
  von `ef76ef71a7ea…` auf `1b84bd4eeba9…`. `postgres`/`nats`-Container-IDs
  werden ebenso vor/nach verglichen (`upgrade_pg_id_before/after`,
  `upgrade_nats_id_before/after`, harter `exit 1` bei Abweichung) — im
  Lauf blieben beide unverändert, `--no-deps` hält real.
- geprüft, ohne Befund: Health-Poll nach dem Tausch (Zeilen ~1610–1622)
  ist strukturell identisch mit dem bestehenden `docker restart`-Rundlauf
  (`LH-QA-REL-001`, Zeilen ~636–645): dieselbe Schleifenlänge (60 ×
  1 Sekunde), derselbe `docker inspect --format
  '{{.State.Health.Status}}'`-Aufruf, dieselbe Abbruchbedingung.
- geprüft, ohne Befund: Datenstand-Beleg im Code nachvollzogen — Zeile
  `id=250` wird vor dem Tausch über `cdc.changes` gelesen und deren
  Position festgehalten, nach dem Tausch erneut per `count(*)`-Abfrage auf
  Identität geprüft (harter `exit 1` bei Abweichung), danach wird `id=251`
  eingefügt und ihr Auftauchen in `cdc.changes` gepollt (harter `exit 1`
  bei Timeout). Der reale Lauf bestätigt beide Belege (Position 30803592
  vor dem Tausch identisch lesbar, Position 30804928 für die danach
  eingefügte Zeile).
- geprüft, ohne Befund: `harness/README.md` — die `make
  test-integration`-Zeile trägt jetzt den `ADR-0064`-Verweis (`Supersedes
  ADR-0058 Entscheidung 3`) statt eines Rückgriffs auf die überholte
  `ADR-0058`-Form, korrekten `· seit slice-063`-Anker.
- geprüft, ohne Befund: Kommentar-Disziplin im neuen Skript-Code
  (`AGENTS.md` §3.7). Der Header-Absatz (Zeilen 32–35) und der
  Phasen-Kommentar (Zeilen 1554–1565) beschreiben, was die Phase **tut**
  und **warum** (`LH-QA-OPS-005`, `ADR-0064`), nicht die Geschichte des
  Blockers als Konjunktiv oder abwesenden Text. Der Verweis „ersetzt den
  in `ADR-0058` vorgesehenen, real blockierten zweiten `make
  schema-rollout`-Lauf (`BEO-PGC/schema-rollout-fremdobjekte`, Exit 8 auf
  vier Fremdobjekten, `docs/reviews/blocker-slice-063.md`)" ist Indikativ
  über den aktuellen Zustand und dessen Begründung, keine Chronik ohne
  Träger — und folgt exakt dem bereits im selben Skript etablierten,
  akzeptierten Muster einer BEO-/Review-Doc-Referenz als Beleg-Anker
  (Zeilen 1010–1011: „Publication-Entzug-Wirksamkeit — isolierter Beleg
  (`BEO-PGC/walsender-wirksamkeit`, `LH-FA-CFG-002`, `ADR-0050`,
  `docs/reviews/architect-verdict-walsender-wirksamkeit.md`)"). Keine
  Slice-/Wellen-Nummer trägt hier die Aussage allein — die primäre
  Begründung ist `LH-QA-OPS-005`/`ADR-0064`, der Review-Doc-Pfad ist
  Beleg-Anker, kein Ersatz dafür.
- geprüft, ohne Befund: Scope-Fidelity — `git diff 265eda8 94d7530
  --name-only` zeigt ausschließlich die drei im Slice-Plan §3 genannten
  Dateien (plus die Slice-Plan-Datei selbst für die DoD-Nachzüge). Keine
  Berührung von `internal/`, `test/integration/integration_test.go`,
  `compose.yaml` oder anderen Bereichen.
- geprüft, ohne Befund: Traceability — Commit-Betreff `feat(test):
  Upgrade-Sicherheits-Rundlauf über realen Container-Tausch
  (LH-QA-OPS-005, ADR-0064)` trägt beide erforderlichen Kennungen, keine
  `SPEC-*`/`ARC-*`-Kennung im Betreff.
- geprüft, ohne Befund: ID-Wertebereich-Kollision — die neuen IDs 250/251
  auf `feed_e2e_full` kollidieren mit keiner der bereits im Skript
  verwendeten IDs (90/91, 95/96, 120/121, 200/201, 210, 230/231, 240/241
  — geprüft per Grep über alle `VALUES (…)`- und `'id' = '…'`-Stellen).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** keine (0 Findings) — kein
Steering-Loop-Zähler-Eintrag aus diesem Report.

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, 0 LOW/INFO-Findings.
Keine Fixrunde nötig.

**Übergabe:** Keine Findings gehen an den Implementer zurück. Da dieser
Reviewer-Lauf zu 0 HIGH/MEDIUM/LOW kommt und keine
Reviewer→Implementer-Rückgabe stattfindet, wird die DoD-Zeile „Review
durchgeführt, Report unter `docs/reviews/` liegt vor" im Slice-Plan
(`docs/plan/planning/in-progress/slice-063-upgrade-sicherheit-schema-rollout-zyklus.md`)
im selben Commit, der diesen Report anlegt, auf `[x]` nachgezogen
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde).
Dieser Report ist ein **Lauf-Beleg** und wird über Läufe hinweg nicht
wieder gelesen. Er ersetzt keine Verifikation — DoD-/Spec-Konformität
prüft der Verifier separat (Modul 11).

## Gate-Läufe dieses Reviews

Beide ungefiltert und mit unmittelbar geprüftem eigenem Exit-Code
(`AGENTS.md` §3.9) — kein Piping/Wrapper zwischen Lauf und
Exit-Code-Prüfung:

- `make gates` — **Exit 0.** Alle inneren Gates grün
  (`baseline-verify`, `docs-check`, `commit-traceability`, `a-check`,
  `coverage-gate` bei 44.40 % über Schwelle 35 %).
- `make test-integration` — **Exit 0.** Der neue Upgrade-Sicherheits-
  Rundlauf ist im Log sichtbar und real erfolgreich: „Upgrade-Sicherheits-
  Rundlauf (`LH-QA-OPS-005`, `ADR-0064`) belegt — realer Container-Tausch
  über \$COMPOSE up -d --force-recreate --no-deps
  (ef76ef71a7ea…→1b84bd4eeba9…), postgres/nats unberührt (--no-deps),
  Feed-Container danach healthy, Datenstand vor dem Tausch (id=250,
  Position 30803592) identisch lesbar, danach eingefügte Zeile (id=251)
  weiterhin erfasst (Position 30804928)"; anschließend `--- PASS:
  TestE2ESchemaChangeDropColumn` und `--- PASS:
  TestE2ESchemaChangeIncompatibleTypeChange`, gesamter Lauf grün.
