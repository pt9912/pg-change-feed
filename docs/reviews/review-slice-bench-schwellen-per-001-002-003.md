# Review-Report: bench-schwellen-per-001-002-003 — 2026-09-19

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/bench-schwellen-per-001-002-003.md`),
`ADR-0104`/`ADR-0054` und `AGENTS.md` Hard Rules (nicht gegen die DoD —
das ist Verifier-Aufgabe, Modul 11).

**Gegenstand:** Commit `ab544949` (Parent `22a3b63f`), einziger Commit
dieses Slice.

**Skill:** `.harness/skills/reviewer.md` @ `70098d3` · **Modell:**
claude-sonnet-5 · **Datum:** 2026-09-19

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/bench-schwellen-per-001-002-003.md` (voll gelesen)
- `docs/plan/adr/0104-benchmark-schwellen-per-001-002-003.md` (neu, voll gelesen)
- `docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md` (voll gelesen, zur Prüfung der Supersede-Präzision)
- `LH-QA-PER-001`/`002`/`003`, `SPEC-013`, `SPEC-014`, `SPEC-025`
- `AGENTS.md` §3 (Hard Rules), insbes. §3.7, §3.12, §3.13
- `git show ab544949` (voller Diff, alle 11 geänderten Dateien einzeln gelesen)
- `make doc-trace` real ausgeführt (netzlos) zur Verifikation der RTM-Zahlen

---

## Findings

### F-1 — `harness/README.md`s `make doc-trace`-Zeile ist nach diesem Diff real falsch, unangefasst gelassen

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.13 (bewegte Eigenschaft nicht nachgezogen); zugleich `AGENTS.md` §3.12 / Reviewer-Skill „Zahl im Träger … driftet gegen die Messung"
- `pfad`: `harness/README.md:129` (Zeile unverändert von diesem Diff gelassen — nur `harness/README.md:138`, die direkt benachbarte `make bench`-Zeile, wurde bearbeitet)
- `befund`: Die `make doc-trace`-Zeile behauptet „76 Anforderungen, 55 Waisen ohne `trace.coverage`, 21 Waisen mit" und zählt `LH-QA-PER-001`…`004` namentlich zu den neun Waisen, deren Beleg „auf einer anderen Ebene als `test-integration`" liegt (`make bench`, `ADR-0054` §(b) — als weiterhin ohne Coverage-Dimension dargestellt). Nach diesem Diff ist das real falsch: `make doc-trace` ausgeführt liefert `76 Anforderung(en), 18 Waise(n)` — `LH-QA-PER-001`/`002`/`003` zeigen jetzt `Bench`/`ok`, nur `LH-QA-PER-004` bleibt in dieser Gruppe Waise. Der Implementer hat dieselbe Datei im selben Commit eine Zeile darunter bearbeitet, ohne die direkt betroffene Nachbarzeile zu prüfen oder ihre Staleness zu melden (kein Suchlauf-Ergebnis dazu in Plan, ADR oder Commit-Message).
- `verifizierbar`: ja — `make doc-trace` (netzlos), real ausgeführt im Rahmen dieses Reviews, bestätigt `76 Anforderung(en), 18 Waise(n)` statt der im Träger stehenden 76/55/21.
- `klasse`: „Arbeit überholt stehenden Nachbar-Träger im selben Diff"

### F-2 — `tools/bench-scaling.sh`-Kopfkommentar widerspricht der im selben Diff eingeführten Pass/Fail-Logik

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.7 („Ein Kommentar beschreibt, was da ist")
- `pfad`: `tools/bench-scaling.sh:1-6`
- `befund`: Der Kopfkommentar sagt wortwörtlich „Kein Pass/Fail — dokumentiertes Ergebnis je Stufe." Dieser Text ist von diesem Diff unangetastet geblieben (identisch zu Parent `22a3b63f`), obwohl derselbe Diff 60 Zeilen weiter unten `THRESHOLD_LAG_SECONDS`/`THRESHOLD_BREACHED` und einen realen `exit 1`-Pfad bei Schwellenüberschreitung einführt. Der Kommentar beschreibt damit nicht mehr, was in der Datei steht.
- `verifizierbar`: ja — Diff-Lektüre (`git show 22a3b63f:tools/bench-scaling.sh` vs. Arbeitsstand, Zeilen 1-6 identisch) plus Code-Lektüre der neuen `if [ "$THRESHOLD_BREACHED" ... ]`-Zeilen im selben File.
- `klasse`: „Kopfkommentar nicht mit eigener Pass/Fail-Änderung nachgezogen"

### F-3 — `tools/bench-batch-vs-single.sh`-Kopfkommentar widerspricht derselben, im selben Diff eingeführten Pass/Fail-Logik

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.7 („Ein Kommentar beschreibt, was da ist")
- `pfad`: `tools/bench-batch-vs-single.sh:2-5`
- `befund`: Der Kopfkommentar endet mit „… dokumentiertes Ergebnis, kein Pass/Fail-Schwellenwert." Unverändert von diesem Diff, obwohl derselbe Diff `THRESHOLD_FACTOR` und einen realen `exit 1`-Pfad bei Unterschreitung des Batch-Vorteils einführt. Dieselbe Fehlerklasse wie F-2, zweites Vorkommen in diesem Diff — nur `tools/bench-source-impact.sh`s Kopfkommentar wurde korrekt nachgezogen.
- `verifizierbar`: ja — Diff-Lektüre plus Code-Lektüre der `THRESHOLD_FACTOR`/`exit 1`-Zeilen im selben File.
- `klasse`: „Kopfkommentar nicht mit eigener Pass/Fail-Änderung nachgezogen" (2. Vorkommen dieses Diffs — Steering-Loop-relevant bei einem dritten Vorkommen in einem Folge-Slice)

### F-4 — Risiko-Ausgang „`make bench` ruft immer alle drei auf" hält der eigenen neuen Fehlklassifikations-Möglichkeit nicht stand

- `kategorie`: MEDIUM
- `quelle`: Maintainability / Slice-Plan §6, vierter Punkt
- `pfad`: `docs/plan/planning/in-progress/bench-schwellen-per-001-002-003.md:155-160`; `Makefile:82-86` (`bench:`-Target); `tools/bench-lib.sh:141-183` (`bench::render_abdeckung`, nur vom letzten Skript aufgerufen)
- `befund`: Der Risiko-Ausgang schließt die Unvollständigkeits-Gefahr für `docs/user/bench-abdeckung.md` mit der Begründung „kein Sensor erzwingt die Reihenfolge, `make bench` selbst ruft immer alle drei auf" — das gilt aber nur, wenn alle drei Skripte mit Exit 0 enden. Reale Prüfung: GNU Make bricht ein Recipe beim ersten nicht-null Exit-Code eines Rezept-Schritts ab (empirisch bestätigt), und `Makefile:84-86` verkettet die drei `bash tools/bench-*.sh`-Aufrufe ohne `-`/`|| true`. Da dieser Diff genau neue `exit 1`-Pfade in `bench-source-impact.sh`/`bench-scaling.sh` einführt, bricht ein realer Schwellen-Fehlschlag in einem der ersten beiden Skripte den `make bench`-Lauf ab, bevor `bench-batch-vs-single.sh` (und damit `bench::render_abdeckung`) läuft — die genau als Gegenargument genannte Eigenschaft trifft in diesem Fall nicht zu. Die praktische Auswirkung ist begrenzt, weil `bench::record_row`s Zeileninhalt eine stabile Deklaration ist (Schwellentext/Ort), keine Lauf-Kennzahl, und weil `record_row` je Skript vor dessen eigenem Exit-1-Pfad läuft — aber `docs/user/bench-abdeckung.md` wird in diesem Fall in dieser Ausführung nicht neu geschrieben.
- `verifizierbar`: ja — reales Minimalbeispiel bestätigt Makes Abbruchverhalten; ein realer `make bench`-Lauf mit einem erzwungenen Fehlschlag in `bench-source-impact.sh` (z. B. `BENCH_SOURCE_IMPACT_N` extrem klein) würde es am Zielobjekt bestätigen.
- `klasse`: „Beleg trägt seinen Satz nicht (Risiko-Ausgang)"

### F-5 — Drei-Liefer-Punkte-Frage im Plan nicht explizit adressiert

- `kategorie`: LOW
- `quelle`: Maintainability / Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice
- `pfad`: `docs/plan/planning/in-progress/bench-schwellen-per-001-002-003.md` §1
- `befund`: §1 zählt vier Ziel-Punkte (Messen, Schwellen festlegen, Pass/Fail verdrahten, RTM-Sichtbarkeit) über drei `LH-QA-PER-*`-IDs, ohne die „≤ 3 Liefer-Punkte"-Leitlinie (Modul 5) explizit zu adressieren. Gemildert durch echtes Repo-Präzedenz (der unmittelbar vorangehende Slice `e2e-drei-rtm-luecken` bündelte ebenfalls drei `LH-*`-IDs unter einem Slice) und dadurch, dass alle vier Punkte an einer einzigen, unteilbaren ADR-Entscheidung (`ADR-0104`) hängen und der Diff in einer Sitzung vollständig prüfbar war.
- `verifizierbar`: nein — Ermessensfrage, kein Gate-Lauf bestätigt „zu groß".
- `klasse`: „Liefer-Punkte-Bündelung ohne explizite Modul-5-Prüfung im Plan"

### F-6 — Untersuchungstiefe der Rausch-Ursache eng, aber proportional

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `docs/plan/adr/0104-benchmark-schwellen-per-001-002-003.md` §Kontext (3)
- `befund`: Die Untersuchung (`docker stats` gegen Fremd-Container, `docker logs` des Feed-Containers) deckt externe Störquellen und Anwendungsfehler ab, prüft aber keine PostgreSQL-internen Signale (Checkpoint-Aktivität, Autovacuum, Host-I/O-Warteschlange) als Alternativerklärung für die Streuung. Kein Blocker: Die gewählte Marge (35 % gegenüber real 25,0–28,2 %; 10× gegenüber real 189–194×) ist großzügig genug, dass eine engere, ungeprüfte Rausch-Quelle das Ergebnis kaum in Frage stellt.
- `verifizierbar`: nein — Hinweis ohne erwartete Aktion.
- `klasse`: —

### F-7 — `ADR-0104`s `Autor:`-Feld ohne Rolleninhaber-Anmerkung

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `docs/plan/adr/0104-benchmark-schwellen-per-001-002-003.md:10`
- `befund`: `ADR-0054` (die teilweise supersedete ADR) trägt „Autor: pt9912 (Rolleninhaber: Architect-Lauf, …)"; `ADR-0104` trägt nur „Autor: pt9912" ohne Rollen-Anmerkung. Kosmetische Inkonsistenz, kein Hard-Rule-Bezug.
- `verifizierbar`: nein.
- `klasse`: —

---

## Positiv-Befunde (was real geprüft wurde und trägt)

- **Supersede-Präzision (`ADR-0104` vs. `ADR-0054`):** `ADR-0054` §(b) selbst gelesen — der Satz „keine Aufnahme in `make gates`/`fullbuild` als Pass/Fail-Bedingung — Zweck ist der dokumentierte Beleg … nicht Durchsetzung" vermengt zwei trennbare Aussagen (Gate-Aufnahme vs. Schwellen-Existenz). `ADR-0104` trennt sie sauber: Titel, Status-Zeile und Entscheidung benennen ausdrücklich nur die „kein Pass/Fail"-Hälfte als supersedet, zitieren `ADR-0054`s Docker-/DB-Abhängigkeits-Begründung gegen `make gates`-Aufnahme wörtlich und lassen sie unangetastet stehen (auch in `Makefile:75-86`, `Fitness Function`-Tabelle beider ADRs konsistent). `ADR-0054` selbst bleibt unverändert (Immutable-Regel `AGENTS.md` §3.5 eingehalten — kein Diff auf die Datei).
- **`SPEC-025`-Platzierung:** `spec/pflichtenheft.md:504` neue Zeile geprüft — kein `ADR-0104`-Verweis darin, nur `LH-*`-Links und der generische, bereits bei `SPEC-013`/`SPEC-014` etablierte Zusatz „über ADR schärfbar". Nummer (`SPEC-025`) ist real die nächste freie in der Tabelle (bis `SPEC-024` vergeben). `ADR-0104`s `Schärft:`-Feld zeigt korrekt auf `SPEC-025` (Referenz-Richtung ADR → Pflichtenheft, nicht umgekehrt).
- **`bench::record_row`/`bench::render_abdeckung`-Mechanik selbst funktioniert real:** `make doc-trace` real ausgeführt bestätigt `LH-QA-PER-001`/`002`/`003` als `Bench`/`ok` — der Kern-Mechanismus dieses Slice ist nachweislich wirksam (unabhängig von F-4s Sequenz-Randfall). `.tmp/bench-abdeckung-rows/` liegt unter der bereits bestehenden, repo-weiten `.tmp/`-Ignorierregel (`.gitignore:11`) — kein Leak-Risiko.
- **`median_of`/`run_phase_median`-Arithmetik korrekt:** `median_of` mit `mid=(count+1)/2` und `sort -n | sed -n "${mid}p"` liefert für 5 Werte korrekt den 3.-kleinsten (echter Median). `run_phase_median`s Offsets (`phase_offset_base + (run-1)*count`) belegen für Phase 1 `[0, 5N)` und für Phase 2 `[5N, 10N)` disjunkte ID-Bereiche — keine PK-Kollision über 5 Läufe × 2 Phasen. `bench::wait_captured … $((RUNS*N))` passt zur tatsächlich in Phase 2 eingefügten Zeilenzahl.
- **Investigations-Kette real nachvollziehbar** (nicht nur behauptet): Die in ADR/Plan genannten Zwischenwerte (17,5 %/1,8 % → 17,3 %/33,0 % bei N=1000 → 25,0 %/28,2 % bei N=5000/5) sind in `ADR-0104` §Kontext (2)/(3) und im Slice-Plan §6 konsistent wiedergegeben, keine widersprüchlichen Zahlen zwischen beiden Trägern gefunden.
- **AGENTS.md §3.7-Konformität der neuen Kommentare:** Alle neu hinzugefügten Kommentare (`THRESHOLD_PCT`, `THRESHOLD_FACTOR`, `THRESHOLD_LAG_SECONDS`, `median_of`, `run_phase_median`, `bench::record_row`, `bench::render_abdeckung`, `.d-check.yml`s `trace.coverage`-Kommentar) beschreiben den aktuellen Zustand und zitieren `ADR-0104`/`SPEC-025`/`SPEC-013` als Rang-Zeiger — keine Slice-/Wellen-Chronik, kein Konjunktiv über verworfene Alternativen im Code. Einzige Ausnahme sind die zwei **unveränderten** Alt-Kopfkommentare aus F-2/F-3.
- **Traceability:** Commit-Betreff nennt `ADR-0104` (und implizit `ADR-0054` im Text), keine `SPEC-*`/`ARC-*`-Kennung im Betreff — regelkonform.
- **Keine Betreiber-Oberfläche berührt:** Keine neue `CDC_*`-Umgebungsvariable, keine `cdc.*`-SQL-Funktion, kein Endpunkt — `docs/user/benutzerhandbuch.md` zu Recht nicht angefasst.

## Negativbefunde

- geprüft, ohne Befund: `spec/pflichtenheft.md` (Spec-Stratum-Grenze, Referenz-Richtung, ID-Vergabe)
- geprüft, ohne Befund: `docs/plan/adr/0104-benchmark-schwellen-per-001-002-003.md` (MADR-Vollständigkeit, ≥3-Alternativen-Regel, Immutabilität von `ADR-0054`)
- geprüft, ohne Befund: `docs/plan/adr/README.md` (Index-Zeile)
- geprüft, ohne Befund: `.d-check.yml` (`trace.coverage`-Verdrahtung, Kommentar-Aktualisierung)
- geprüft, ohne Befund: `docs/user/bench-abdeckung.md` (generierter Inhalt konsistent mit den drei Skripten und mit `bench::record_row`/`render_abdeckung`)
- geprüft, ohne Befund: `tools/bench-lib.sh` (neue Helfer, Kommentar-Klassen, `.tmp`-Ignorierung)
- geprüft, ohne Befund: `tools/bench-source-impact.sh` (Median-/Offset-Arithmetik, Kopfkommentar korrekt nachgezogen)
- geprüft, ohne Befund: Docker-only-Disziplin (§3.1), Suppression-Verbot (§3.2), `git mv`-Trennung (§3.3), Action-Pinning (§3.8) — keine dieser Klassen im Diff berührt
- nicht ausgeführt: voller `make bench`-Lauf (Docker-Image-Bau + DB) — Logik stattdessen statisch verifiziert (Median-/Offset-Arithmetik, Schwellen-Vergleiche); `make doc-trace` dagegen real ausgeführt (netzlos)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 3 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Arbeit überholt stehenden Nachbar-Träger
im selben Diff · Kopfkommentar nicht mit eigener Pass/Fail-Änderung
nachgezogen (2×) · Beleg trägt seinen Satz nicht (Risiko-Ausgang) ·
Liefer-Punkte-Bündelung ohne explizite Modul-5-Prüfung im Plan

## Verdikt

**Merge-blockierend:** ja — drei HIGH-Findings (F-1, F-2, F-3), alle
mechanisch einfach zu beheben (Nachzug einer Träger-Zeile bzw. zweier
Kopfkommentare), aber real falsch bzw. real widersprüchlich zum
eigenen Diff-Inhalt.

**Übergabe:** F-1 bis F-4 gehen an den Implementer zurück (Fixrunde).
F-5 bis F-7 sind Hinweise ohne Rückgabe-Zwang; F-5 kann bei der
Slice-Closure im Beobachtungs-Register vermerkt werden, falls sich das
Bündelungsmuster wiederholt. Da dieses Verdikt eine Fixrunde vorsieht,
wird die DoD-Checkbox „Review durchgeführt" in
`docs/plan/planning/in-progress/bench-schwellen-per-001-002-003.md` **nicht**
in diesem Report nachgezogen (Reviewer-Skill §DoD-Checkbox-Nachzug ohne
Fixrunde — Grenzfall trifft hier nicht zu) — der reguläre Nachzug läuft
nach der Fixrunde am Implementer-Workflow-Schritt 21.

Dieser Report ist ein Lauf-Beleg (Audit: dieser Diff, dieser Skill,
dieses Modell, dieses Verdikt) und ersetzt keine Verifikation — DoD-/
Spec-Konformität prüft der Verifier separat (Modul 11).
