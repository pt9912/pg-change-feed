# Verifikationsbericht: slice-050 — 2026-09-13

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-050` §1 Ziel/Abgrenzung, §2 DoD, §3 Plan-Nachzug, §4 Trigger, §6
Risiken, §8 Sub-Area) und `ADR-0054` §(b) im Wortlaut — nicht gegen Diff
(Reviewer-Aufgabe, bereits abgeschlossen: `review-slice-050.md`, vollständig
gelesen) und nicht gegen realen Bedarf (Validator, hier nicht ausgelöst —
kein MVP-Meilenstein-Slice, reine Mess-Infrastruktur ohne neue
Architektur-Sicht-Aussage).

**Frischer Kontext:** Diese Prüfung liest den vollständigen Slice-Plan
(inkl. Plan-Nachzug und §7 Closure-Notiz), `ADR-0054` vollständig, den
Review-Report vollständig, den tatsächlichen Diff aller drei Commits
(`5a2ccc8..7fe4b15`) selbst — keine Implementer- oder Reviewer-Behauptung
wird ungeprüft übernommen. `make gates` und `make bench` (alle drei
Skripte, Default-Modus) wurden in dieser Sitzung selbst real ausgeführt.

**Gegenstand:**
`docs/plan/planning/in-progress/slice-050-performance-benchmark-infrastruktur.md`
zum Stand `HEAD = 7fe4b15`. Commits (chronologisch): `5a2ccc8` (`next` →
`in-progress`, reiner `git mv`), `fe45065` (Implementierung: vier
`tools/bench-*.sh`-Skripte, `Makefile`-Target `bench`,
`harness/README.md`-Zeile), `51c0cb6` (DoD-Nachzug: Plan-Nachzug, §6-
Risiko-Ausgänge, Beobachtungs-Registerzeile), `7fe4b15` (Review-Report, 0
HIGH/MEDIUM/LOW, 2 INFO, DoD-Checkbox „Review durchgeführt" im selben
Commit nachgezogen, keine Fixrunde). Sequenz selbst geprüft: Implementierung
→ Review — kein Self-Review (Implementer- und Reviewer-Läufe sind getrennte
Kontexte laut Report-Kopf), keine Rolle springt rückwärts ohne Artefakt.

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | `LH-QA-PER-001` real erfüllt: messbarer Unterschied mit/ohne aktivem Replication-Slot | **erfüllt, selbst reproduziert** | Eigener Lauf `make bench` (Skript `tools/bench-source-impact.sh`, N=1000): real `ohne CDC — 3540 ms`, `mit CDC — 6777 ms`, Differenz `3237 ms (91.4%)`. Größenordnung deckt sich mit dem im Slice-Plan dokumentierten Wert (+86,4 %), Absolutwert weicht erwartungsgemäß ab (andere Docker-Umgebung/Laufzeitpunkt) — siehe §3 unten zum Streuungs-Risiko. |
| 2 | `LH-QA-PER-002` real erfüllt: alle drei `SPEC-014`-Lastenstufen durchlaufen, dokumentiertes Ergebnis je Stufe | **erfüllt, selbst reproduziert** | Eigener Lauf `tools/bench-scaling.sh` (Default-Modus): `klein — Ziel 10/s über 10s, 100 Zeilen (~10.0/s), cdc_capture_lag=1.011050s`; `mittel — Ziel 100/s über 15s, 1500 Zeilen (~100.0/s), cdc_capture_lag=1.009720s`; `gross — Ziel 1000/s über 15s, 15000 Zeilen (~1000.0/s), cdc_capture_lag=1.009612s`. Alle drei Ziel-Raten real erreicht, `cdc_capture_lag` in allen drei Stufen nahe 1 s — deckt sich mit der Plan-Aussage. |
| 3 | `LH-QA-PER-003` real erfüllt: Batch- vs. Einzelabruf über `cdc.changes` | **erfüllt, selbst reproduziert** | Eigener Lauf `tools/bench-batch-vs-single.sh` (M=200): `Batch-Abruf — 129 ms`, `Einzelabruf — 16538 ms`, Faktor `128.2x`. Reale Größenordnung bestätigt (Plan dokumentiert 178,8×) — dieselbe Streuungs-Einordnung wie oben. |
| 4 | `make bench` startet alle drei Skripte, kein Gate; `harness/README.md` §Werkzeuge trägt die neue Zeile | **erfüllt** | Eigener `make bench`-Lauf: alle drei Skripte liefen nacheinander im selben Aufruf, korrekte Reihenfolge (Source-Impact → Scaling → Batch-vs-Single), jeweils mit eigenständigem Umgebungsaufbau/-abbau (`docker ps -a`/`docker network ls` danach leer — Cleanup real bestätigt). `Makefile:75-79` real gelesen: `.PHONY: bench` / `bench: image` mit drei `@bash tools/…`-Zeilen, kein Eintrag in `GATE_CHECKS`. `harness/README.md:132` trägt die neue Zeile in der Werkzeug-Tabelle (nicht in der Sensors-Tabelle) mit `kein Gate`-Kennzeichnung. |
| 5 | `make gates` grün (unverändert, kein neues Gate) | **erfüllt, selbst reproduziert** | Eigener, vollständiger `make gates`-Lauf gegen `HEAD = 7fe4b15` (Exit 0): `baseline-verify: v6.5.0 OK — 54 Dateien`; `coverage-gate: OK — Coverage 39.60% erfüllt Schwelle 35%` (unverändert seit slice-049, dieser Slice fügt kein neues Gate hinzu); `d-check: 396 Datei(en) geprüft, 0 Befund(e)` (zweimal, docs+commits-Modul); `commit-traceability: OK — 5 Commit(s) …`; `a-check: gesamt: 0 Befund(e)`. `grep bench` gegen `.github/workflows/ci.yml` und `AGENTS.md` §4 selbst ausgeführt: kein Treffer — `bench` ist real nirgends in `gates`/`fullbuild`/CI verdrahtet. |
| 6 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **erfüllt** | `docs/reviews/review-slice-050.md` vollständig gelesen: 0 HIGH, 0 MEDIUM, 0 LOW, 2 INFO (F-1, F-2), Verdikt „nicht merge-blockierend". DoD-Zeile im Slice-Plan korrekt im selben Commit (`7fe4b15`) nachgezogen. |
| 7 | Doku-Update `harness/README.md` §Werkzeuge | **erfüllt** | Siehe Punkt 4. |
| 8 | Closure-Notiz mit Steering-Loop-Lerneintrag (§7) | **erfüllt** | §7 vollständig geschrieben: Was funktionierte / was anders lief / Steering-Loop-Eintrag (bewusst „keiner" — dokumentiert begründet: reine Umsetzung einer bereits in `ADR-0054` §(b) entschiedenen Infrastruktur) / Beobachtungs-Register / Folge-Slices / §6-Risiken. |
| 9 | Reconciliation-Register — entfällt | **korrekt entfällt** | `docs/plan/planning/reconciliation.md` real geprüft: Datei existiert nicht — Repo ist Greenfield (`harness/conventions.md` Modus-Deklaration `*`/`PGC` = Greenfield), Item entfällt korrekt. |
| 10 | Beobachtungs-Register fortgeschrieben | **erfüllt** | `docs/plan/planning/observations/BEO-PGC/schema-rollout-braucht-compose-init/` real geprüft: `observation.md`, `state.md` (Zähler 1×, unter der Schwelle), `evidence/slice-050.md` — Form deckt sich mit der Ziel-Form (`observation.template.md`); inhaltlich real nachvollzogen (fehlender `compose-init`-Mount, andere Fehlerklasse als `BEO-PGC/coverage-stage-dockerignore-blockiert-tooling` aus `slice-049` — eigenständig geprüft, nicht nur Reviewer-Zitat übernommen: unterschiedliche Ursache, `compose.yaml`-unabhängige Umgebung vs. neue Docker-Multi-Stage-Stufe). |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | **erfüllt, mit eigenem Urteil — siehe §3 unten** | Beide Risiken tragen einen der drei zulässigen Ausgänge (*entfallen* / *eingetreten*), beide mit Begründung. Eigenes Urteil zum ersten Ausgang unten. |
| 12 | Drei Paarungen — an `welle-14`-Closure delegiert | **korrekt delegiert** | `docs/plan/planning/welle-14.md` real gelesen: liegt flach (offen), kein Slice dieser Welle in `done/`. Delegation an die Welle-Closure ist zutreffend (Modul 6/8) — DoD-Checkbox bleibt korrekt unangehakt. |

## 2. Sensor-Läufe (selbst ausgeführt, `HEAD = 7fe4b15`)

**`make gates`** (vollständiger Lauf, Exit 0):

```
baseline-verify: v6.5.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)
… [coverage-Stage] … total: (statements) 39.6%
coverage-gate: OK — Coverage 39.60% erfüllt Schwelle 35%
d-check: 396 Datei(en) geprüft, 0 Befund(e)
d-check (commits, HEAD~5..HEAD): 396 Datei(en) geprüft, 0 Befund(e)
commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID
a-check: gesamt: 0 Befund(e)
```

**`make bench`** (voller Lauf, Default-Modus, alle drei Skripte in einem Aufruf):

```
bench-source-impact: ohne CDC — 3540 ms für 1000 Zeilen
bench-source-impact: mit CDC — 6777 ms für 1000 Zeilen
bench-source-impact: Ergebnis (LH-QA-PER-001) — ohne CDC 3540 ms, mit CDC 6777 ms, Differenz 3237 ms (91.4%) für je 1000 Schreibtransaktionen auf public.bench_source_impact
bench-scaling: Stufe klein — Ziel 10/s über 10s, 100 Zeilen in 10s eingefügt (~10.0/s), cdc_capture_lag=1.011050s
bench-scaling: Stufe mittel — Ziel 100/s über 15s, 1500 Zeilen in 15s eingefügt (~100.0/s), cdc_capture_lag=1.009720s
bench-scaling: Stufe gross — Ziel 1000/s über 15s, 15000 Zeilen in 15s eingefügt (~1000.0/s), cdc_capture_lag=1.009612s
bench-scaling: Ergebnis (LH-QA-PER-002) — alle drei SPEC-014-Lastenstufen real durchlaufen (Modus: verkürzt)
bench-batch-vs-single: Batch-Abruf — 129 ms für 200 Zeilen (1 Roundtrip)
bench-batch-vs-single: Einzelabruf — 16538 ms für 200 Zeilen (200 Roundtrips)
bench-batch-vs-single: Ergebnis (LH-QA-PER-003) — Batch 129 ms vs. Einzelabruf 16538 ms für 200 Zeilen (Einzelabruf 128.2x langsamer)
```

Nach dem Lauf: `docker ps -a`/`docker network ls` gegen `pgc-bench-*` real
leer — Cleanup (`bench::cleanup` als `EXIT`-Trap in jedem Skript) bestätigt.
Nebeneffekt festgestellt und **nicht** übernommen: `tools/schema/plan.yaml`
(ein von `make schema-rollout` geschriebener Pflicht-Report, `git`-getrackt)
trug nach dem Lauf ein anderes `target`-DSN (`pgc-bench-postgres` statt dem
zuletzt committeten CI-Wert) — mit `git checkout -- tools/schema/plan.yaml`
zurückgesetzt, da dieser Lauf-Zeitpunkt-Artefakt kein Bestandteil dieser
Verifikation ist und sonst fälschlich eine inhaltliche Änderung vorgetäuscht
hätte.

## 3. §6-Risiken — eigenes, unabhängiges Urteil

- **Risiko „Streuung" — Ausgang *entfallen*.** Die eigene Reproduktion
  liefert eine reale Bestätigung, dass Streuung tatsächlich auftritt: mein
  Lauf ergab 91,4 % CDC-Overhead statt der im Plan dokumentierten 86,4 %,
  und einen Faktor 128,2× statt 178,8× für Batch-vs-Einzelabruf — eine nicht
  triviale Differenz zwischen zwei unabhängigen Läufen derselben Skripte auf
  demselben Host. Das bestätigt die Prämisse des Risikos (reale Streuung in
  einer geteilten/virtualisierten Docker-Umgebung), stützt aber trotzdem den
  gewählten Ausgang *entfallen* — und zwar aus demselben Grund, den
  Implementer und Reviewer bereits nennen: Kein Skript setzt einen
  Pass/Fail-Schwellenwert (`ADR-0054` §(b) verlangt das ausdrücklich nicht),
  und in beiden Läufen liegt keine der drei Kennzahlen in einer Größenordnung,
  die eine andere qualitative Aussage nahelegen würde (CDC-Schreib-Overhead
  bleibt deutlich > 0, Batch bleibt zwei Größenordnungen schneller,
  `cdc_capture_lag` bleibt in allen drei Lastenstufen nahe 1 s). Ein
  „weiter offen" wäre hier eine Überzeichnung: Es gibt keinen Mechanismus im
  Repo, gegen den ein Ausgang „weiter offen" wirken könnte (kein Gate, keine
  gespeicherte Referenzzahl, gegen die künftige Läufe geprüft würden) — das
  Risiko ist an eine Schwelle gebunden, die es nicht gibt. **Eigenes Urteil:
  Der Ausgang trägt**, unabhängig bestätigt durch reale Zweitmessung mit
  spürbar abweichenden Absolutwerten. Der von Reviewer-F-2 benannte
  Lesbarkeits-Aspekt (dokumentierte Einzelwerte könnten künftig als stabile
  Fakten statt als Einzelmessung gelesen werden) ist real begründet — meine
  eigene, deutlich abweichende Zweitmessung ist der beste verfügbare Beleg
  dafür — bleibt aber zu Recht ein INFO- statt ein DoD-Mangel: Der DoD-Punkt
  verlangt „dokumentiertes Ergebnis", nicht „stabiler Referenzwert", und §3
  Plan-Nachzug formuliert die Zahlen bereits als „real gemessen" (Bezug auf
  diesen einen Lauf), nicht als Dauerzusage.
- **Risiko „groß-Stufe-Dauer" — Ausgang *eingetreten*, gelöst im Slice.**
  Eigene Reproduktion des Default-Modus bestätigt reale, kurze Laufzeiten
  (10 s/15 s/15 s) bei unveränderter Ziel-Rate — der Default-Lauf des
  gesamten `make bench` inklusive aller drei Skripte lief in dieser Sitzung
  in deutlich unter zwei Minuten durch. Der `--full`-Pfad wurde nicht selbst
  reproduziert (das wäre laut Plan ~80 Minuten allein für die
  `SPEC-014`-Dauern) — das ist konsistent mit der Trigger-Absicht: Ein
  schneller, wiederholbarer Implementer-/Reviewer-/Verifier-Lauf ist der
  Regelfall, `--full` ein bewusst seltener Sonderfall. **Ausgang trägt.**

## 4. Plan-vs-Code-Diff

- **Drei Skripte statt eines Kombi-Skripts:** `ADR-0054` §(b) real
  eingehalten — jedes der drei `tools/bench-*.sh` bleibt einzeln aufrufbar,
  berechnet seine Kennzahl (Zeitdifferenz/Prozent, Durchsatz/Lag, Faktor)
  jeweils selbst inline; `tools/bench-lib.sh` (eigenständig gelesen, nicht
  nur Reviewer-Befund übernommen) trägt ausschließlich Umgebungs-Setup
  (Netzwerk/Container/Schema-Rollout/`psql`-Wrapper/Warte-Funktion) — keine
  Mess-, Differenz- oder Faktor-Berechnung in dieser Datei. Die geteilte
  Datei ist eine ungeplante, aber im Plan-Nachzug begründete vierte Datei;
  sie verletzt die ADR-Grundintention nicht, weil „je Beleg ein eigenes
  Skript" sich auf die Messung bezieht, nicht auf jeden einzelnen
  Umgebungs-Handgriff.
- **Kein Gate-Charakter:** `Makefile`-Diff (`fe45065`) real gelesen —
  `bench: image` hängt an keinem `GATE_CHECKS`-Eintrag; `.github/workflows/ci.yml`
  trägt real keinen `bench`-Bezug; `AGENTS.md` §4 (Quality Gates) trägt
  real keine `make bench`-Zeile — die Zeile steht korrekt nur in
  `harness/README.md` unter der separaten Werkzeug-Tabelle, nicht unter
  Sensors.
- **`SPEC-014` unverändert:** `git diff 5a2ccc8..7fe4b15 -- spec/` real
  leer — keine Spec-Änderung in diesem Slice. Die drei Lastenstufen
  (10/100/1000 pro Sekunde) werden in `tools/bench-scaling.sh` unverändert
  als Eingabeparameter übernommen; die 60-s-Annahme für die „klein"-Stufe
  ist im Skript-Kommentar korrekt als bench-skript-eigene Annahme
  ausgewiesen, nicht als `SPEC-014`-Schärfung — `spec/pflichtenheft.md:264`
  real gegengelesen: für „klein" ist tatsächlich nur die Rate (≤ 10/s)
  festgelegt, keine Dauer.
- Kein weiterer, unbegründeter Abweichungspunkt zwischen Plan/`ADR-0054`
  und Code gefunden.

## 5. Scope-Treue gegen §1 (Ausdrücklich NICHT in diesem Slice)

- **Ausschluss „Aufnahme in `make gates`/`fullbuild`":** eingehalten (§2
  Punkt 5, §4 oben).
- **Ausschluss „Neue Lastenstufen-Definition":** eingehalten — `SPEC-014`
  unverändert (§4 oben).
- **Ausschluss „Test-Coverage-Gate" (`slice-049`):** kein Bezug zu
  `harness/mk/coverage.mk`/`Dockerfile`-Stage `coverage` im Diff dieses
  Slice — eingehalten, keine Überschneidung.
- **Ausschluss „`LH-QA-PER-004`-Ausbau":** kein neuer
  Commit→CDC-Latenz-Beleg über den bestehenden `cdc_capture_lag`-Lasttest
  hinaus im Diff — eingehalten.

## 6. Hard Rules

- **3.7 (Slice-Chronik-Verbot):** eigener `grep -rniE
  "slice-[0-9]+|welle-[0-9]+"` gegen alle vier neuen `tools/bench-*.sh`
  liefert keinen Treffer. Für den `Makefile`/`harness/README.md`-Diff
  gezielt gegen die tatsächlich **hinzugefügten** Zeilen geprüft (`git diff
  … | grep "^+"`, nicht der volle Kontext-Diff, der Alt-Zeilen mit
  legitimen `seit slice-<NNN>`-Herkunfts-Ankern anderer Sensoren enthält) —
  ebenfalls kein Treffer. Hard Rule gewahrt.
- **3.6 (keine Gate-Lockerung ohne ADR):** nicht einschlägig — dieser Slice
  fügt kein Gate hinzu und senkt keine Schwelle.
- **3.3 (git mv + Inhaltsänderung = zwei Commits):** `5a2ccc8` ist ein
  reiner `git mv` (`next` → `in-progress`, real per `git show --stat`
  bestätigt: nur Rename, keine Inhaltsänderung).

## 7. Explizit NICHT geprüft (korrekt außerhalb dieser Rolle)

Drei Paarungen (DoD-Punkt 12) — `welle-14` ist offen, die Prüfung ist
zutreffend an die Welle-14-Closure delegiert (Modul 6/8). Validierung gegen
realen Bedarf (kein MVP-Meilenstein-Slice, kein Validator-Zug ausgelöst,
reine Mess-Infrastruktur ohne neue Architektur-Sicht-Aussage).

## Verdikt

**DoD-Konformität: bestätigt.** Alle zwölf DoD-Punkte real geprüft, elf
davon materiell erfüllt (eigene Reproduktion von `make bench` und
`make gates`), einer (Reconciliation) korrekt entfallen, einer (drei
Paarungen) korrekt an die Welle-14-Closure delegiert. Keine offene
Diskrepanz zwischen Behauptung und Beleg.

**Plan-vs-Code-Diff:** keine unbegründete Abweichung — Drei-Skripte-Muster,
Kein-Gate-Charakter und `SPEC-014`-Unveränderlichkeit sind im Plan benannt
und durch den Diff real gedeckt. Die ungeplante vierte Datei
(`tools/bench-lib.sh`) ist im Plan-Nachzug begründet und wahrt die
Grundintention von `ADR-0054` (eigenständig gegen den Dateiinhalt geprüft,
nicht nur den Reviewer-Befund zitiert).

**§6-Risiken:** beide Ausgänge tragen ein eigenständiges Urteil — für
„Streuung" (*entfallen*) liefert die eigene Zweitmessung mit spürbar
abweichenden Absolutwerten (91,4 % statt 86,4 %; 128,2× statt 178,8×) den
stärksten verfügbaren Beleg dafür, dass Streuung real ist, aber am
gewählten Ausgang nichts ändert, weil kein Schwellenwert existiert, den sie
verfälschen könnte; für „groß-Stufe-Dauer" (*eingetreten*, gelöst) bestätigt
der eigene Default-Lauf die praktikable Kurzlaufzeit.

**Scope-Treue (§1):** eingehalten — kein Gate-Charakter, keine
`SPEC-014`-Änderung, keine Überschneidung mit `slice-049`.

**Hard Rules:** 3.7 gewahrt (kein Chronik-Kommentar in den neuen
Bench-Skripten oder den hinzugefügten `Makefile`/`README`-Zeilen, eigener
`grep` gegen die tatsächlichen `+`-Zeilen), 3.6 nicht einschlägig, 3.3
gewahrt.

**Übergabe an Planner:** Der Slice kann an die Planner-Closure übergeben
werden. Für die Closure-Notiz vorzumerken: F-1 (Einzelabruf-Overhead
konfliert mit Werkzeug-Overhead, INFO) und F-2 (Risiko-Ausgang deckt
Gate-Frage, nicht Lesbarkeits-Frage, INFO) — beide reine Notizen, keine
Rückführung nötig. Die drei Paarungen bleiben zutreffend an die
`welle-14`-Closure delegiert. Kein Validator-Zug ausgelöst.

---

*Dieser Bericht ist ein Lauf-Beleg (Modul 11) und wird über Läufe hinweg
nicht wieder gelesen.*
