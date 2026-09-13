# Review-Report: slice-050 — 2026-09-13

**Review-Art:** Code — Code-Review gegen Plan + Konventionen (Modul 10
§Drei Review-Arten), geprüft gegen Plan/ADR/Hard Rules (Maintainability),
**nicht** gegen die DoD (Verifier-Aufgabe, Modul 11).

**Gegenstand:** `slice-050` — Diff `5a2ccc8..51c0cb6` (Commits `fe45065`
Implementierung, `51c0cb6` DoD-Nachzug).

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-09-09 (vier
repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-050-performance-benchmark-infrastruktur.md`
  (inkl. §3 Plan-Nachzug)
- `ADR-0054` (`docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md`) §(b)
- `docs/plan/planning/done/slice-049-test-coverage-gate.md` (Referenz)
- `LH-QA-PER-001`…`003`, `SPEC-014` (`spec/pflichtenheft.md`)
- `AGENTS.md` §3 (Hard Rules), `harness/conventions.md` (MR-000)

---

## Findings

Keine HIGH-, MEDIUM- oder LOW-Findings. Zwei INFO-Hinweise.

### F-1 — Einzelabruf-Messung konfliert Lese-Ineffizienz mit Werkzeug-Overhead

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `tools/bench-batch-vs-single.sh:43-52,57`
- `befund`: Der "Einzelabruf"-Pfad ruft je Zeile einen frischen
  `docker exec`+`psql`-Prozess auf (`bench::psql_scalar`), nicht nur eine
  zusätzliche SQL-Anweisung über eine bereits offene Verbindung. Der
  dokumentierte Faktor (178,8× bei M=200, real reproduziert: 72 ms vs.
  12874 ms) enthält damit Prozess-Spawn- und Verbindungsaufbau-Overhead
  des Mess-Werkzeugs selbst zusätzlich zur eigentlichen
  Batch-vs.-Einzelabruf-Lesekosten-Differenz, die `LH-QA-PER-003`
  eigentlich adressiert (ein realer Consumer mit persistenter Verbindung
  würde diesen Overhead nicht in dieser Höhe tragen). Die Messung ist
  real und reproduzierbar (kein Messfehler im Sinne einer falschen
  Anweisung), deckt die Zielgröße aber nur indirekt ab.
- `verifizierbar`: nein — kein Gate-Lauf prüft Benchmark-Methodik; nur
  durch Wiederholung mit einer persistenten Verbindung (z. B. ein
  kleines Go-/psql-`\watch`-Skript, das die Verbindung offen hält)
  überprüfbar.
- `klasse`: Benchmark-Ergebnis deckt Zielgröße nur indirekt ab

### F-2 — Risiko "Streuung" als *entfallen* markiert, ohne den Lesbarkeits-Aspekt der dokumentierten Absolutwerte zu adressieren

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `docs/plan/planning/in-progress/slice-050-performance-benchmark-infrastruktur.md:228-242`
  (§6, Risiko 1) und `§7` (Closure-Notiz)
- `befund`: Die Begründung für den Ausgang *entfallen* trägt korrekt für
  die Gate-Frage (keine Pass/Fail-Schwelle betroffen, `ADR-0054` §(b)).
  Sie adressiert nicht explizit, dass einzelne dokumentierte
  Absolutwerte (z. B. "+86,4 %", "178,8×") von künftigen Lesern (Roadmap,
  Closure-Notiz, ein späterer Slice, der auf diese Zahlen verweist) als
  stabile Fakten statt als Einzelmessung in einer geteilten/virtualisierten
  Umgebung gelesen werden könnten — das ist kein Gate-Risiko, aber ein
  Lesbarkeits-/Interpretations-Risiko, das der gewählte Ausgang nicht
  vollständig abdeckt.
- `verifizierbar`: nein — Interpretationsfrage, kein Gate-Gegenstand.
- `klasse`: Risiko-Ausgang deckt Gate-Frage, nicht Lesbarkeits-Frage

## Negativbefunde

- geprüft, ohne Befund: `tools/bench-lib.sh` gegen `ADR-0054` §(b) „je
  Beleg ein eigenes Skript" — die Datei trägt ausschließlich
  Umgebungs-Setup (Docker-Netzwerk/Container, Schema-Rollout, `psql`-
  Wrapper, Erfassungs-Wartefunktion als Synchronisation); keine
  Mess-/Differenz-/Faktor-Berechnung. Die drei Bench-Skripte berechnen
  ihre Kennzahlen (Zeitdifferenz, Prozent, Faktor) jeweils selbst inline
  und bleiben einzeln lauffähig und lesbar (je eigener `trap`, eigene
  Tabellen-/Slot-/Publication-Namen, eigenes Ergebnis-Format). Die
  ADR-Grundintention (kein monolithisches Kombi-Skript, das drei
  Messungen vermischt) ist gewahrt.
- geprüft, ohne Befund: `SPEC-014`-Wortlaut (`spec/pflichtenheft.md:264`)
  — nennt für die "klein"-Stufe tatsächlich **nur** die Rate (≤ 10/s),
  keine Dauer; die 60-s-Annahme in `tools/bench-scaling.sh` ist korrekt
  als eigene, dokumentierte Bench-Skript-Annahme ausgewiesen, nicht als
  `SPEC-014`-Schärfung deklariert. Kein Spec-Stratum-Verstoß.
- geprüft, ohne Befund: Default-/`--full`-Split in `tools/bench-scaling.sh`
  — Skript-Header, Ergebnis-Zeile ("Modus: verkürzt/voll") und
  Plan-Nachzug §3 machen konsistent klar, dass der Default-Modus
  verkürzte Dauern bei unveränderter Ziel-Rate fährt und keine volle
  `SPEC-014`-Konformität beansprucht; nicht irreführend.
- geprüft, ohne Befund: `Makefile`-Diff (`fe45065`) — `bench:`-Target
  hängt an keinem Eintrag in `GATE_CHECKS`; `.github/workflows/ci.yml`
  enthält keinen `bench`-Bezug; `AGENTS.md` §4 (Quality Gates) trägt
  korrekt keine `make bench`-Zeile (kein Gate).
- geprüft, ohne Befund: `SPEC-014`-Lastenstufen (klein ≤10/s, mittel
  100/s×30min, groß 1000/s×60min) in `tools/bench-scaling.sh` und
  Slice-Kopf — unverändert als Eingabeparameter übernommen, keine
  Neu-Festlegung.
- geprüft, ohne Befund: `docs/plan/planning/observations/BEO-PGC/schema-rollout-braucht-compose-init/`
  — Form korrekt (`observation.md`/`state.md`/`evidence/slice-050.md`
  nach Ziel-Form `observation.template.md`); Abgrenzung zur
  `slice-049`-Beobachtung `BEO-PGC/coverage-stage-dockerignore-blockiert-tooling`
  real geprüft und bestätigt — andere Ursache (fehlender
  `compose-init`-Mount bei einer von `compose.yaml` unabhängigen
  PostgreSQL-Instanz vs. `.dockerignore`/Alpine-`bash`-Problem bei einer
  neuen Docker-Multi-Stage-Stufe); diese Skripte fügen keine neue Stage
  hinzu.
- geprüft, ohne Befund: Hard Rule 3.7 (Slice-Chronik-Verbot) — `grep -rniE
  "slice-[0-9]+|welle-[0-9]+"` gegen alle vier neuen `tools/bench-*.sh`
  und den `fe45065`-Diff von `Makefile`/`harness/README.md` liefert keinen
  Treffer; die vorhandenen `slice-*`-Zitate im bestehenden `Makefile` sind
  Altbestand außerhalb dieses Diffs (per `git diff 5a2ccc8..fe45065 --
  Makefile` verifiziert).
- geprüft, ohne Befund: Docs-Check-Falle (bare IDs) — `make docs-check`
  real ausgeführt: „395 Datei(en) geprüft, 0 Befund(e)".
- geprüft, ohne Befund: Commit-Traceability — `RANGE=5a2ccc8..51c0cb6
  make commit-traceability` real ausgeführt: „OK — 2 Commit(s) …,
  Betreffs ohne Struktur-ID".
- geprüft, ohne Befund: DoD-Checkboxen in §2 — materiell nachgezogen bis
  auf „Review durchgeführt" (korrekt offen, wird unten nachgezogen) und
  „Drei Paarungen" (korrekt an `welle-14`-Closure delegiert, `welle-14`
  ist offen).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Benchmark-Ergebnis deckt Zielgröße nur
indirekt ab · Risiko-Ausgang deckt Gate-Frage, nicht Lesbarkeits-Frage

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, 0 LOW; die zwei
INFO-Hinweise erwarten keine Implementer-Aktion (Skill §Klassifikation:
„Hinweis ohne erwartete Aktion").

**Zu den drei markiert geprüften Punkten:**

1. **`tools/bench-lib.sh`:** Trägt real nur Umgebungs-Setup, keine
   Mess-/Vergleichslogik. Die Grundintention von `ADR-0054` §(b) (drei
   eigenständige Belege statt eines vermischten Kombi-Skripts) ist
   gewahrt — die geteilte Datei ist Infrastruktur-Kopplung, keine
   Melde-/Ergebnis-Kopplung. Kein Finding.
2. **Default-/`--full`-Split + 60-s-Annahme:** `SPEC-014` legt für
   „klein" tatsächlich keine Dauer fest (nur ≤10/s) — real im
   Pflichtenheft-Wortlaut bestätigt. Die 60-s-Annahme ist als
   Bench-Skript-eigene, dokumentierte Annahme ausgewiesen, nicht als
   Spec-Schärfung. Der Default-/`--full`-Split ist im Skript-Header, in
   der Ergebniszeile und im Plan-Nachzug konsistent und nicht
   irreführend dokumentiert. Kein Finding.
3. **Risiko „Streuung" als *entfallen*:** Die Begründung (keine
   Pass/Fail-Schwelle betroffen) trägt für die Gate-Frage und ist damit
   im Sinne von Modul 5 ein zulässiger Ausgang. Sie deckt jedoch nicht
   vollständig ab, dass dokumentierte Einzelwerte künftig als stabile
   Fakten fehlinterpretiert werden könnten — dieser Aspekt ist als
   INFO (F-2) festgehalten, rechtfertigt aber keinen anderen Ausgang
   („weiter offen" wäre hier eine Überzeichnung, da kein Mechanismus im
   Repo existiert, der ein wiederholtes Auftreten dieser Art von
   Fehlinterpretation beobachtbar zählen könnte).

**Reale Bench-Ausgaben (Punkt 10):** Die 178,8×-Größenordnung für
Batch- vs. Einzelabruf ist plausibel und kein Messfehler im Sinne einer
falschen Abfrage — sie resultiert aus 200 tatsächlich einzeln
ausgeführten `docker exec`+`psql`-Aufrufen (F-1). Die 86,4 %-CDC-
Schreib-Overhead-Zahl ist unauffällig (zwei sequenzielle Phasen
innerhalb desselben Laufs, dieselbe Instanz).

**Übergabe:** Die DoD-Checkbox „Review durchgeführt, Report unter
`docs/reviews/` liegt vor" wird in diesem Commit selbst nachgezogen
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde,
Modul 8: keine Fixrunde nötig). Die Finding-Klassen (F-1, F-2) gehen in
den Zähler des Beobachtungs-Registers ein, sofern eine künftige
Slice-Closure sie als wiederkehrend erkennt — für `slice-050` selbst ist
§7 bereits geschlossen; sie werden hier als Lauf-Beleg festgehalten, ohne
eine weitere Registerzeile zu erzwingen (Erstauftreten, unter der
Schwelle).
