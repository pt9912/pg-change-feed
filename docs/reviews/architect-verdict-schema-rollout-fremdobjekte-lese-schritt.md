# Review-Report: Architect-Verdikt `schema-rollout-fremdobjekte`-Lese-Schritt — 2026-09-18

**Review-Art:** Governance/Entscheidungs-Review (Diff gegen Register-
Konvention + Konventionen, Modul 10 §Drei Review-Arten), nicht gegen DoD
(Verifier-Aufgabe).

**Gegenstand:** Commit `7d0bf05` — Architect-Entscheidung zum 3×-Lese-
Schritt von `BEO-PGC/schema-rollout-fremdobjekte` (Ausgang zunächst
`verkörpert`, neue Hard Rule `AGENTS.md` §3.14), ausgelöst am Ende von
`welle-beispiele-start-ueber-make`s Closure (alle vier Slices `done/`).

**Skill:** `.harness/skills/reviewer.md`
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-18

**Eingangs-Kontext:**

- `docs/plan/planning/observations/BEO-PGC/schema-rollout-fremdobjekte/`
  (`observation.md`, `state.md`, alle drei Evidence-Dateien)
- `docs/plan/planning/observations/BEO-PGC/d-migrate-nacharbeit/`
- `AGENTS.md` §3.13 (Präzedenzfall für die Stilform), §3.14 (neu)
- `docs/plan/planning/observations/README.md` (Register-Konvention)
- `Makefile` (Zeilen ~181-186, `schema-rollout`-Target)
- Eigene Messung: `git show 7d0bf05 --stat` (Scope-Grenze),
  Kreuz-Verifikation der drei Evidence-Zitate gegen `slice-016`/`slice-063`/
  `slice-beispiele-compose-bootstrap`, `make gates` (ungepiped Exit 0).

---

## Findings

### F-1 — Systemische Alternative nicht erwogen

- `kategorie`: MEDIUM
- `quelle`: Maintainability / Architect-Reasoning-Vollständigkeit
- `pfad`: `AGENTS.md` §3.14 (wie in Commit `7d0bf05` geschrieben),
  `docs/plan/planning/observations/BEO-PGC/schema-rollout-fremdobjekte/state.md`
  (Fassung `7d0bf05`)
- `befund`: Die Begründung erwog nur zwei Gegenmaßnahmen — (a) Root-
  Cause-Fix ins neutrale Modell (blockiert an `POST_EXECUTE_DRIFT`,
  korrekt verifiziert) und (b) ein generisches `--allow-destructive`-
  Handling im Makefile-Target (zurecht als zu grobkörnig verworfen). Eine
  dritte, vom `POST_EXECUTE_DRIFT`-Blocker unabhängige Alternative fehlte
  in der Abwägung: den Existenz-Check zentral im `schema-rollout`-
  Makefile-Target selbst zu verankern statt in jedem Aufrufer einzeln —
  das Target hat bereits `$(PG_TEST_IMAGE)`/`$(SCHEMA_TARGET)` zur Hand.
  Das ist genau das Risiko, das die Begründung selbst benennt („ein
  viertes Auftreten ist wahrscheinlich") und mit „verkörpert"
  (Governance-Pflicht für Menschen/Agenten) statt „geplant" (ein
  konkreter Engineering-Task) beantwortet, ohne die Alternative sichtbar
  erwogen/verworfen zu haben.
- `verifizierbar`: nein (Architektur-Abwägung, kein Gate-Lauf)
- `klasse`: „Systemische Alternative nicht erwogen"

### F-2 — Fehlender ADR-Anker trotz Commit-Zitat

- `kategorie`: LOW
- `quelle`: Konsistenz mit Nachbar-Regeln
- `pfad`: `AGENTS.md` §3.14 (Fassung `7d0bf05`)
- `befund`: §3.14 trug keinen „Träger und Anker"-Absatz mit explizitem
  ADR-Verweis, obwohl die Commit-Message `ADR-0043` nennt und mehrere
  Nachbar-Regeln (§3.11, §3.12, §3.13) diese Form pflegen.
- `verifizierbar`: nein
- `klasse`: „Fehlender ADR-Anker trotz Commit-Zitat"

### F-3 — Zitationsform weicht vom etablierten Muster ab (fremde Datei, nur zur Kenntnis)

- `kategorie`: INFO
- `quelle`: Stilkonvention
- `pfad`: `docs/plan/planning/done/welle-beispiele-start-ueber-make-results.md`
  (zum Prüfzeitpunkt noch nicht committet, Orchestrator-Dokument)
- `befund`: Verwendet an einer Stelle einen echten Markdown-Link auf
  `AGENTS.md` statt der sonst in `done/welle-*-results.md`-Dateien
  üblichen reinen Inline-Code-Zitation. Nicht Gegenstand dieses
  Architect-Zugs, nur zur Vollständigkeit vermerkt.
- `verifizierbar`: ja (Grep-Vergleich)
- `klasse`: „Zitationsform weicht vom etablierten Muster ab"

## Negativbefunde

- geprüft, ohne Befund: alle drei Evidence-Zitate (Objektzahlen 2→4→6,
  konkrete View-/Funktionsnamen) stimmen exakt mit den Evidence-Texten
  überein.
- geprüft, ohne Befund: `BEO-PGC/d-migrate-nacharbeit`-Zitat
  (`POST_EXECUTE_DRIFT`/Exit 5) wortgleich zum dortigen `state.md`
  belegt.
- geprüft, ohne Befund: Stilform der Hard Rule (Falsch/Richtig/
  Begründung, Nummerierung 3.14 direkt vor „## 4. Quality Gates")
  konform zu §3.1–§3.11.
- geprüft, ohne Befund: `AGENTS.md` §3.7/§3.12-Konformität des neuen
  Textes — durchgehend Indikativ, keine Chronik-Sprache.
- geprüft, ohne Befund: Aktionabilität — ein Implementer, der nur §3.14
  liest, hat eine konkrete, ausführbare Anweisung.
- geprüft, ohne Befund: Register-Konvention — `observation.md`
  unverändert, `state.md`-Ausgang „verkörpert" mit auflösbarem Anker im
  selben Format wie Präzedenzfälle.
- geprüft, ohne Befund: Scope von `7d0bf05` — exakt `AGENTS.md` + eine
  `state.md`-Datei.
- geprüft, ohne Befund: `make gates` — ungepiped Exit 0.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Systemische Alternative nicht erwogen ·
Fehlender ADR-Anker trotz Commit-Zitat · Zitationsform weicht vom
etablierten Muster ab

## Verdikt

**Merge-blockierend:** nein — 0 HIGH. Das eine MEDIUM-Finding (F-1) löst
eine Fixrunde am Architect-Zug aus (Modul 8 §Konflikt-Pfad nicht nötig,
kein Rollen-Widerspruch — der Architect reicht die Alternative selbst
nach). Die Fixrunde ist in Commit `5f7b28f` erfolgt (Ausgang geändert auf
„geplant", `AGENTS.md` §3.14 um den Träger-und-Anker-Absatz ergänzt,
neuer Folge-Slice terminiert) und ist Gegenstand eines eigenen,
nachfolgenden Reviews.

Dieser Report ersetzt keine Verifikation — DoD-/Spec-Konformität prüft der
Verifier separat (Modul 11), sofern dieser Lese-Schritt eine eigene
Verifikation vorsieht.
