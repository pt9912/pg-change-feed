# Review-Report: slice-049 — 2026-09-13

**Review-Art:** Code — Diff gegen Plan/ADR/Konventionen (Modul 10 §Drei
Review-Arten). Verifikation gegen DoD/Spec bleibt Verifier-Aufgabe.

**Gegenstand:** Commit `3f1054b` (Diff gegen Elter `6fcc0b0`) — Test-Coverage-Gate.

**Skill:** `.harness/skills/reviewer.md` @ `3f1054b`
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-049-test-coverage-gate.md` (inkl. §3
  Plan-Nachzug)
- `docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md`
- `AGENTS.md` §3 (Hard Rules, insb. 3.1 Docker-only, 3.6 Gate-Lockerung nur
  per ADR, 3.7 Kommentar-Disziplin)
- `harness/conventions.md` (MR-000 ID-Schema)
- Vorbild-Repo `/Development/d-check` (Dockerfile `coverage`-Stage,
  `tools/coverage-gate.sh`, `Makefile`)

---

## Findings

### F-1 — Sensors-Tabellen-Zeile weicht im Target-Feld vom etablierten Zeilenformat ab

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `harness/README.md:117`
- `befund`: Alle Geschwisterzeilen der Sensors-Tabelle (`baseline-verify`,
  `docs-check`, `a-check`, `commit-traceability`, `gates`, `image`, …) tragen
  den Target in Spalte 1 als reinen Inline-Code (` `make <target>` `); nur die
  neue `coverage-gate`-Zeile verlinkt den Target selbst
  (`[`make coverage-gate`](sensors/coverage-gate.md)`) und dupliziert damit
  denselben Link, der bereits in der Bindung-Spalte steht.
- `verifizierbar`: nein — kein Gate prüft Tabellen-Zellformat; `make
  docs-check` bleibt grün, da der Link auflöst.
- `klasse`: „Sensors-Zeile Format-Abweichung Target-Spalte“

### F-2 — Beobachtungs-Register: `Ausgang`-Feld vor Erreichen der 3×-Schwelle belegt

- `kategorie`: INFO
- `quelle`: Baseline-Regelwerk `modul-06-roadmap.md` §Das
  Beobachtungs-Register — „Unterhalb der Schwelle ist `offen` der Stand —
  dort ist er kein Ausgang, sondern der Normalzustand."
- `pfad`: `docs/plan/planning/observations/BEO-PGC/coverage-stage-dockerignore-blockiert-tooling/state.md:1`
- `befund`: `state.md` trägt bei Zähler 1× bereits „Zustand: offen — Ausgang:
  **weiter offen**" — kanonisch wird der Ausgang erst vom Lese-Schritt bei
  Erreichen von 3× zugewiesen. Das ist keine Neuerung dieses Slice: Dasselbe
  Muster liegt bereits in `BEO-PGC/cdc-capture-lag-real/state.md` (2×,
  „Ausgang: eingetreten") vor — eine bestehende, repo-weite Lokalkonvention,
  nicht ein von diesem Diff eingeführter Verstoß.
- `verifizierbar`: nein
- `klasse`: „Register-Ausgang vor Schwelle belegt"

---

## Zu den zehn Prüfpunkten (real nachvollzogen)

1. **Arithmetik/Rot-Grün-Paar:** 39,6 % abgerundet auf den nächsten
   5-%-Schritt = 35 % — korrekt. Real nachgebaut:
   `make coverage-gate THRESHOLD=45` → `FAIL — Coverage 39.60% unter
   Schwelle 45%`, Exit 1; `make coverage-gate` (Default `THRESHOLD=35`) →
   `OK — Coverage 39.60% erfüllt Schwelle 35%`, Exit 0. Beide Belege decken
   sich exakt mit §3 Plan-Nachzug und `harness/sensors/coverage-gate.md`.
2. **Hard Rule 3.6:** Die Einstiegsstufe 35 % ist durch die in `ADR-0054`
   §(a) vorab beschlossene Eskalationsklausel gedeckt (Ist-Stand, abgerundet
   auf 5 %, niemals über 80 %) — keine eigenmächtige Schwellen-Erfindung.
   Endstufe 80 % bleibt unverändert fest.
3. **`SHELL ["/bin/bash", "-eo", "pipefail", "-c"]`:** `apk add --no-cache
   bash` läuft im Dockerfile-Diff korrekt *vor* der `SHELL`-Direktive (noch
   mit Alpine-Default-`sh`). Real bestätigt: `pipefail` maskiert den
   Exit-Code von `go tool cover -func=… | tee …` nicht — der reale
   `THRESHOLD=45`-Lauf endet mit Exit 1 des gesamten `RUN`-Befehls trotz
   `tee` im letzten Gliedsegment davor.
4. **`.dockerignore`-Fix:** minimal-invasiv — nur `!tools/coverage-gate.sh`
   ergänzt, nicht `!tools/` pauschal. Kein weiteres, ungewolltes
   Wieder-Einschließen von Dateien.
5. **DB-Adapter-Tests skippen ohne DSN:** §3 Plan-Nachzug und §6 stellen den
   Zusammenhang ehrlich dar (Coverage-Zahl künstlich gedrückt, reale Tests
   existieren und laufen in `make test-store`/`make test-replication`
   grün) und benennen die drei `[no test files]`-Pakete korrekt als reine
   Konstanten-/Typ-Träger ohne Fachlogik-Lücke. Keine Verschleierung
   erkennbar.
6. **Slice-Chronik in Kommentar (Hard Rule 3.7):** `grep` gegen Dockerfile,
   `tools/coverage-gate.sh`, `harness/mk/coverage.mk`,
   `harness/sensors/coverage-gate.md` zeigt keine `slice-\d+`/`welle-\d+`-
   Zitate in Code/Config; die einzige Fundstelle `seit slice-049` liegt in
   `harness/sensors/coverage-gate.md` §Bindung — zulässiger
   Herkunfts-Anker in Doku, kein Chronik-Kommentar in Produktionscode.
7. **Neue Beobachtung `BEO-PGC/coverage-stage-dockerignore-blockiert-tooling`:**
   Form korrekt — `observation.md` (Identität), `state.md` (Stand + Zähler
   abgeleitet), `evidence/slice-049.md` (ein Beleg-Vorgang). Siehe F-2 zum
   `Ausgang`-Feld.
8. **Docs-Check:** `make docs-check` real ausgeführt — 390 Dateien geprüft,
   0 Befunde. Kein übrig gebliebener `id-unlinked`-Fund.
9. **Commit-Traceability/DoD:** `make commit-traceability` real ausgeführt
   (`RANGE=HEAD~5..HEAD`) — OK, 5 Commits, Betreffs ohne Struktur-ID. DoD
   ist bis auf „Review durchgeführt" und „drei Paarungen" korrekt
   nachgezogen; die Paarungen sind zutreffend an die offene `welle-14`-
   Closure delegiert (nicht an diesem Slice zu erfüllen).
10. **`harness/README.md` §Sensors / `AGENTS.md` §4:** `AGENTS.md` §4 folgt
    dem etablierten Zeilenformat (Target als Inline-Code, Beschreibung,
    Referenzen am Ende in Klammern). `harness/README.md` §Sensors weicht im
    Target-Feld vom Format aller Geschwisterzeilen ab — siehe F-1.

## Negativbefunde

- geprüft, ohne Befund: `Dockerfile` (neue `coverage`-Stage) — Reihenfolge,
  Kommentar-Disziplin, Scope-Treue zu `ADR-0054`
- geprüft, ohne Befund: `tools/coverage-gate.sh` — Awk-Vergleichslogik,
  Exit-Codes, Kommentarkopf
- geprüft, ohne Befund: `harness/mk/coverage.mk` — `GATE_CHECKS`-Verdrahtung,
  `NO_CACHE_FILTER_COV`, Kalibrierungs-Bindung-Kommentar
- geprüft, ohne Befund: `AGENTS.md` §4 — Zeilenformat, ADR-Bezug
- geprüft, ohne Befund: `docs/plan/planning/in-progress/slice-049-test-coverage-gate.md`
  §1–§8 — Ziel/Abgrenzung, DoD, Plan-Nachzug, Trigger, Risiken-Ausgänge,
  Closure-Notiz, Sub-Area-Begründung
- geprüft, ohne Befund: `.dockerignore` — Minimal-Invasivität des Fixes
- geprüft, ohne Befund: Commit-Message `3f1054b` — Traceability-Bezug
  (`ADR-0054` im Betreff), Inhalt deckt sich mit realen Belegen

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Sensors-Zeile Format-Abweichung
Target-Spalte · Register-Ausgang vor Schwelle belegt

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM; das einzige LOW ist eine
Format-Abweichung ohne Gate-Wirkung, das INFO ist eine bereits bestehende
Lokalkonvention (nicht neu eingeführt). Arithmetik und Rot-/Grün-Beleg sind
real nachvollzogen und decken sich mit der Slice-Dokumentation. Hard Rule
3.6 ist gewahrt (Eskalationsstufe durch `ADR-0054` gedeckt, keine
eigenmächtige Senkung), Hard Rule 3.7 ist gewahrt (kein Chronik-Kommentar in
Code/Config).

**Übergabe:** Keine Fixrunde nötig. F-1 und F-2 gehen als LOW/INFO ohne
Reviewer→Implementer-Rückgabe-Pfeil in die Slice-Closure §7 (Finding-Klassen
in den Zähler). Nach `.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne
Fixrunde zieht dieser Report die DoD-Zeile „Review durchgeführt" im
Slice-Plan im selben Commit nach. Die drei Paarungen bleiben zutreffend an
die `welle-14`-Closure delegiert. Dieser Report ist Lauf-Beleg und wird über
Läufe hinweg nicht erneut gelesen; DoD-/Spec-Konformität prüft der Verifier
separat (Modul 11).
