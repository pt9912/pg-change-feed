# Review-Report: slice-036 — Fixrunden-Bestätigung — 2026-09-13

**Review-Art:** Fixrunden-Bestätigung — gezielte Nachprüfung der drei
MEDIUM-Findings aus `docs/reviews/review-slice-036.md` (F-1/F-2/F-3) gegen
Commit `3f57e8d`, kein erneutes Vollreview des Slice. Prüfmaßstab: `ADR-0050`,
`AGENTS.md` §3, Modul 6 (Beobachtungs-Register).

**Gegenstand:** Commit `3f57e8d` (`state.md` des Beobachtungs-Registers
nachgezogen; neuer Testfall
`TestAdministrationRequestEnableTableRequiresCdcAdminMembership`;
DoD-Begründung im Slice-Plan korrigiert).

**Skill:** `.harness/skills/reviewer.md` @ HEAD · **Modell:** claude-sonnet-5
· **Datum:** 2026-09-13

---

## F-1 — Beobachtungs-Register-Notiz nicht mit neuem Beleg nachgezogen

**Verdikt: behoben.**

Selbst gelesen: `state.md` und alle fünf `evidence/*.md`-Dateien
(`slice-006`, `slice-010`, `slice-015`, `slice-016`, `slice-036`) in
`docs/plan/planning/observations/BEO-PGC/d-migrate-nacharbeit/`.

- Zähler in `state.md` nennt jetzt „5×" mit allen fünf Dateinamen — deckt
  sich mit dem tatsächlichen Verzeichnisinhalt (`ls evidence/` liefert genau
  diese fünf Dateien).
- Die neue Formulierung nennt **drei** Objektklassen (CHECK-Constraint,
  Views, Funktionen) statt der vorherigen zwei — inhaltlich korrekt gegen
  die Evidence-Dateien geprüft:
  - CHECK-Constraint (`chk_change_operation`) „aufgelöst seit slice-015":
    deckungsgleich mit `evidence/slice-015.md` („konvergiert. Exit 0 …
    Ausweichform ist zurückgebaut").
  - Views „aufgelöst seit slice-016": deckungsgleich mit
    `evidence/slice-016.md` („technisch aufgelöst … `nacharbeit-views.sql`
    ist gelöscht").
  - Funktionen „bleibt offen … `POST_EXECUTE_DRIFT`": deckungsgleich mit
    `evidence/slice-036.md` (bricht für jede über den `functions:`-Knoten
    deklarierte Funktion mit `POST_EXECUTE_DRIFT`/Exit 5 ab, Ausweichform
    `nacharbeit-administration.sql` bleibt bestehen).
- Keine Über- oder Untertreibung: `state.md` behauptet nicht, dass die
  Funktionsklasse aufgelöst sei, und benennt korrekt die Bedingung für ihre
  künftige Auflösung (d-migrate blendet die Objektklasse analog zu Views aus
  dem Post-Compare-Fingerabdruck aus).

`verifizierbar`: ja — Dateiabgleich `evidence/*.md` gegen `state.md`-Text
(kein automatisierter Sensor, aber deterministisch nachvollziehbar).

## F-2 — fehlende Negativtests bei neuem öffentlichem Vertrag

**Verdikt: behoben.**

Neuer Testfall gelesen und real reproduziert — **nicht** nur den
Implementer-Bericht übernommen:

- Eigener, isolierter Testcontainer (PostgreSQL 18, gepinnter Digest wie
  `run-store-tests.sh`), eigener Schema-Rollout über `make schema-rollout`
  gegen dieselbe DB.
- Manuelle Gegenprobe unabhängig vom Go-Test: `SET ROLE cdc_reader; SELECT
  cdc.enable_table(...)` liefert `ERROR: permission denied for function
  enable_table` — das ist SQLSTATE 42501 (`insufficient_privilege`), exakt
  die Fehlerklasse, die `permissionDenied()` (`roles_test.go`) prüft. Kein
  Verbindungsfehler, kein anderer Fehlergrund.
- Anschließend der tatsächliche Go-Testfall isoliert ausgeführt
  (`go test -run TestAdministrationRequestEnableTableRequiresCdcAdminMembership
  -v ./internal/adapters/driven/postgresstorage/...`) gegen dieselbe
  Instanz: `PASS`.
- Der Test nutzt echtes `SET ROLE` auf einer dedizierten `Acquire()`-
  Verbindung (nicht über den Pool) — dasselbe, korrekte Muster wie
  `roles_test.go`; `RESET ROLE` im `defer` verhindert Rollen-Leck auf die
  gepoolte Verbindung.

`verifizierbar`: ja — reproduziert in dieser Sitzung, zweifach (manuelle
psql-Probe + realer Go-Testlauf).

## F-3 — DoD-Begründung mit unzutreffender Tatsachenbehauptung

**Verdikt: behoben.**

Eigener Grep gegen `harness/README.md`:

```
grep -n "nacharbeit" harness/README.md
131: … Die Ausweichform `nacharbeit-views.sql` ist zurückgebaut · seit slice-016 …
```

- `nacharbeit-views.sql` wird tatsächlich genannt — ausschließlich im
  Kontext ihrer *Auflösung* („zurückgebaut · seit slice-016"), nicht als
  aktive Ausweichform.
- `nacharbeit-roles.sql`, `nacharbeit-observability.sql`,
  `nacharbeit-heartbeat.sql` kommen im gesamten `harness/README.md` **kein
  einziges Mal** vor — bestätigt die neue Begründung, dass die drei
  aktiven Ausweichformen dort nicht namentlich genannt werden.
- Die korrigierte Begründung im Slice-Plan zieht daraus den richtigen
  Schluss: `nacharbeit-administration.sql` fällt unter dasselbe Muster wie
  die drei aktiven, nicht genannten Ausweichformen — der Präzedenzfall
  trägt nur Auflösungen. Kein Doku-Update nötig, Begründung trägt jetzt.

`verifizierbar`: ja — Grep-Ergebnis oben ist der vollständige Beleg.

## Gate-Lauf

`make gates` (dieser Review-Lauf, eigenständig ausgeführt, HEAD =
`3f57e8d`): `baseline-verify` (54 Dateien) · `d-check` (298 Dateien, 0
Befunde) · `commit-traceability` (5 Commits, OK) · `a-check` (0 Befunde) —
alle grün.

## Summary

| Finding | Verdikt |
|---|---|
| F-1 (Register-Notiz) | behoben |
| F-2 (fehlender Negativtest) | behoben |
| F-3 (DoD-Begründung) | behoben |

Kein neues Finding aus dieser Fixrunde. `make gates` grün.

## Verdikt

**Merge-blockierend:** keins. Alle drei MEDIUM-Findings aus
`review-slice-036.md` sind durch Commit `3f57e8d` behoben, jeweils
eigenständig gegen den tatsächlichen Datei-/Datenbankzustand nachgeprüft,
nicht aus dem Implementer-Bericht übernommen. Kein Rollen-Konflikt, keine
Architect-Sequenz nötig (Modul 8).

**Übergabe:** Der Verifier prüft DoD-/Spec-Konformität separat (Modul 8);
dieser Report bestätigt nur die Findings-Behebung und ist Lauf-Beleg wie
sein Vorgänger.
