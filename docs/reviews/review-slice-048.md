# Review-Report: slice-048 — 2026-09-13

**Review-Art:** Code — geprüft gegen Plan (inkl. Plan-Nachzug) + Konventionen
(Modul 10 §Drei Review-Arten). Keine DoD-/Spec-Konformitätsprüfung — das ist
Verifier-Aufgabe (Modul 11).

**Gegenstand:** `slice-048` (Commit `2e4b24d`, Diff gegen Elter `089834f`) —
`docs/plan/planning/in-progress/slice-048-retention-lebenszyklus-kombinierter-e2e-rundlauf.md`,
`tools/harness/run-integration-tests.sh`.

**Skill:** `.harness/skills/reviewer.md` @ `587fb9c` (2026-09-13)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext:**

- Slice-Plan `slice-048` inkl. §3 Plan-Nachzug (4 Punkte)
- `LH-FA-RET-002`…`006` (`spec/lastenheft.md`)
- `docs/plan/planning/done/slice-043-changestoreport-loeschmethode.md`,
  `slice-044`, `slice-045`, `slice-046`, `slice-047` (Vorgänger-Belege,
  referenziert im Plan-Nachzug)
- `internal/domain/model/retention.go` (`RetentionPolicy.AllowsDeletion`),
  `internal/application/usecase/retention/service.go`
  (`RunRetentionService.Run`) — gelesen, um die zentrale Behauptung des
  Plan-Nachzugs zu prüfen
- `tools/schema/schema.yaml` (View `cdc.retention_blockers`, Zeile 328ff.)
- `AGENTS.md` §3 Hard Rules, `harness/conventions.md` (MR-000/MR-001)
- `docs/plan/planning/observations/BEO-PGC/*` (Isolation, Test-Runner)

---

## Findings

### F-1 — DELETE-Workaround steht vor, nicht nach dem realen Löschbeleg

- `kategorie`: MEDIUM
- `quelle`: `LH-FA-RET-004`
- `pfad`: `tools/harness/run-integration-tests.sh:472-505`
- `befund`: Die zweite Bestätigung (`acknowledge-consumer … "$lifecycle_position"`,
  Zeile 472) und `DELETE FROM cdc.consumer_position …` (Zeile 481) laufen
  **vor** der Polling-Schleife, die die reale Löschung von `id=210` belegt
  (Zeile 493-505) — nicht danach, wie der Plan-Nachzug Punkt 2 es
  nahelegt ("die Entfernung dient ausschließlich der Sichtbarkeits-Prüfung …
  und der Wiederherstellung der Ausgangslage"). Per Code-Lektüre
  (`RetentionPolicy.AllowsDeletion`, `internal/domain/model/retention.go:35-44`)
  blockiert eine bestätigte Position, die nicht mehr vor der Change-Position
  liegt, die Löschung nicht — die zweite Bestätigung allein reicht aus,
  unabhängig vom späteren `DELETE`. Diese Unabhängigkeit ist damit als
  **Code-Tatsache** verifiziert, aber **nicht** durch den Testlauf selbst
  empirisch demonstriert: weil `DELETE` sofort (ohne Wartezeit) nach der
  zweiten Bestätigung läuft, ist zum Zeitpunkt jedes real beobachtbaren
  Lösch-Takts (`retentionInterval=10s`) die Consumer-Position bereits
  entfernt. Der Testlauf beobachtet real ausschließlich den Zustand
  "keine bestätigte Position vorhanden" — denselben Zustand, den bereits
  bestehende Abschnitte (nie bestätigender Consumer, `slice-045`) abdecken
  — nicht den Zustand "bestätigte, fortgeschrittene Position weiterhin
  vorhanden". Eine künftige Regression, die eine fortgeschrittene, aber
  noch vorhandene Position fälschlich weiter blockieren ließe, würde von
  diesem Abschnitt nicht erkannt.
- `verifizierbar`: ja — durch Umstellen der Reihenfolge (`DELETE` nach der
  bestehenden Polling-Schleife für `lifecycle_deleted`, mit einer
  Wiederholung der `blocker_after`-Sichtbarkeitsprüfung danach) und einen
  erneuten `make test-integration`-Lauf ließe sich die Unabhängigkeit auch
  empirisch statt nur über Code-Lektüre zeigen.
- `klasse`: Testbeleg-Reihenfolge verdeckt behauptete Kausal-Unabhängigkeit

## Negativbefunde

- geprüft, ohne Befund: DELETE-Workaround selbst (Frage 1) — die
  Kern-Behauptung "Bestätigung allein reicht aus" ist per Domain-Code
  (`RetentionPolicy.AllowsDeletion`) korrekt und keine Verschleierung; die
  `DELETE`-Zeile beeinflusst nachweislich nur die View-Sichtbarkeit
  (`cdc.retention_blockers` trägt sonst weiterhin eine Zeile mit
  `backlog = 0`, siehe `tools/schema/schema.yaml:328-345`, kein
  `WHERE`-Filter auf `backlog`) — kein Eingriff in den Löschmechanismus
  selbst. Einordnung als eigenständiges Finding unter F-1, keine
  zusätzliche Kategorie nötig.
- geprüft, ohne Befund: Interferenz mit dem nachfolgenden "Zustand 1 — kein
  Blocker"-Abschnitt (Frage 3) — nach dem `DELETE` trägt
  `cdc.consumer_position` für `src-mvp` real wieder null Zeilen; die
  Go-Integrationstest-Consumer sind zu diesem Zeitpunkt bereits per
  `t.Cleanup` entfernt (bestätigt über die Skript-Reihenfolge: `go test`
  mit `TestMVPRetentionBlockersViewShowsFurthestBehindConsumer` läuft vor
  Zeile 380, `CLI_CONSUMER`/`BACKLOG_CONSUMER` registrieren erst danach,
  Zeile 556/688) — keine Kollision.
- geprüft, ohne Befund: Isolation (`BEO-PGC/test-isolation-geteilter-zustand`)
  — `id=210` auf `feed_mvp_full` kollidiert mit keinem bestehenden
  Wertebereich (1/7, 90/91, 95/96, 120/121, 200/201); `cli-e2e-lifecycle-consumer`
  ist namentlich getrennt von `cli-e2e-consumer`/`cli-e2e-backlog-consumer`;
  keine Wiederverwendung der Variablennamen an anderer Stelle im Skript.
- geprüft, ohne Befund: Slice-Chronik in Produktions-/Skript-Kommentaren
  (Hard Rule 3.7) — kein `slice-\d+`/`welle-\d+` in den neu hinzugefügten
  Kommentarzeilen des Diffs.
- geprüft, ohne Befund: Docs-Check-Falle bei `[LH-FA-RET-002](pfad)…006` —
  etabliertes, bereits an mehreren Stellen (`slice-028`, `welle-13`)
  verwendetes Muster; `make docs-check` lief real grün (0 Befunde).
- geprüft, ohne Befund: `cdc_storage_bytes`-Erwartung — Skript prüft
  ausschließlich `^[0-9]+$` vor und nach der Kette, keine Erwartung einer
  Größenreduktion; konsistent mit Plan §1 Ausschluss.
- geprüft, ohne Befund: Commit-Traceability und DoD-Checkboxen — Commit
  nennt `LH-FA-RET-002…006` im Betreff, `commit-traceability.sh` und
  `d-check --enable commits` liefen real grün gegen `089834f..2e4b24d`;
  DoD-Checkboxen konsistent mit dem tatsächlichen Diff-Umfang (Doku-Update
  korrekt als "Geprüft: nein" begründet, keine `docs/user/benutzerhandbuch.md`-
  Änderung im Diff).
- geprüft, ohne Befund: Beobachtungs-Register — keine neue Datei unter
  `docs/plan/planning/observations/` im Diff; konsistent mit §7/§8 des
  Slice-Plans ("keine Beobachtung angefallen", beide gesichteten Treffer
  bleiben bei 1×).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Testbeleg-Reihenfolge verdeckt behauptete Kausal-Unabhängigkeit

## Verdikt

**Merge-blockierend:** nein — 0 HIGH. Das eine MEDIUM-Finding (F-1) betrifft
nicht die Korrektheit von Produktionscode oder einen ADR-/Hard-Rule-Verstoß,
sondern die empirische Reichweite eines Testbelegs: Der Slice erfüllt sein
eigenes Ziel (Blocker-Übergang real über `cdc.retention_blockers` sichtbar
gemacht, Zeile real gelöscht, `cdc_storage_bytes` durchgehend numerisch) und
die zentrale Behauptung des Plan-Nachzugs ("Bestätigung allein reicht aus,
`DELETE` ist rein kosmetisch für die Sichtbarkeits-Prüfung") ist durch
Code-Lektüre (`RetentionPolicy.AllowsDeletion`) durch den Reviewer
unabhängig bestätigt — nicht widerlegt, nur nicht durch die konkrete
Zeitfolge des Skripts selbst belegt. Kein Datenverlust-, Sicherheits- oder
Suppression-Risiko.

**DELETE-Workaround (Kernfrage des Auftrags):** sauber und ehrlich in der
Absicht — er greift nachweislich nicht in den Löschmechanismus ein, sondern
ausschließlich in die View-Sichtbarkeit (durch die `INNER JOIN`-Semantik von
`cdc.retention_blockers` ohne `backlog`-Filter erzwungen, kein Workaround um
eine sonst nicht erreichbare Sichtbarkeitsaussage). Die einzige Einschränkung:
die Reihenfolge im Skript belegt die Unabhängigkeit nicht selbst empirisch,
weil `DELETE` vor statt nach der Lösch-Polling-Schleife läuft — das ist ein
Präzisions-, kein Ehrlichkeitsmangel.

**Fixrunde:** nicht nötig. 0 HIGH, und das einzige MEDIUM-Finding wird ohne
Reviewer→Implementer-Rückgabe-Pfeil weitergereicht (zur Aufnahme ins
Beobachtungs-Register bei der nächsten Slice-Closure, falls ein weiteres
Vorkommen dieses Musters auftritt). Entsprechend wird die DoD-Zeile "Review
durchgeführt" in `slice-048` in diesem Commit selbst nachgezogen
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde).

**Übergabe:** Findings gehen an den Implementer als Notiz (kein
Rückgabe-Pfeil); die Finding-Klasse geht in die Slice-Closure §7 und von
dort in den Zähler. Dieser Report ist Lauf-Beleg und wird über Läufe hinweg
nicht wieder gelesen. Verifikation gegen DoD/Spec bleibt Verifier-Aufgabe.
