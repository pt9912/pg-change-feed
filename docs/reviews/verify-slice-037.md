# Verifikationsbericht: slice-037 — 2026-09-13

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen
Plan (`slice-037` §1/§2 DoD/§3/§4/§6/§8) und DoD, nicht gegen Diff
(Reviewer-Aufgabe, bereits abgeschlossen) und nicht gegen realen Bedarf
(Validator, hier nicht ausgelöst — kein MVP-Meilenstein-Slice).

**Frischer Kontext:** Diese Prüfung liest den vollständigen, aktuellen
Slice-Plan, `ADR-0050` vollständig, beide Review-Reports und den
tatsächlichen Code selbst — keine Behauptung des Implementers oder
Reviewers wird ungeprüft übernommen; jeder unten genannte Testlauf und
Gate-Lauf wurde in dieser Sitzung **selbst** ausgeführt, nicht aus den
Reports zitiert.

**Gegenstand:** `docs/plan/planning/in-progress/slice-037-administrations-goroutine-live-reload.md`
zum Stand `HEAD = 3f77a85` (Commits `2a771bc`, `b763253`, `c78aa1d`;
`51a81e5`/`3f77a85` sind die Review-Reports; `48d03ba` ist der
themenfremde `ADR-0051`-Commit dazwischen, nicht Gegenstand dieser
Prüfung — bestätigt: `git show 48d03ba --stat` berührt ausschließlich
`AGENTS.md`, `docs/plan/adr/0051-*.md`, `docs/plan/adr/README.md`).

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | Administrations-Goroutine verarbeitet reale `pending`-Anträge, ruft Enable/Disable, schreibt `applied`/`failed` | **erfüllt** | `runAdministration`/`processAdministrationRequests`/`applyAdministrationRequest` gelesen (`internal/bootstrap/wiring.go:734-841`); Verhalten über sechs Whitebox-Tests (`administration_internal_test.go`) real reproduziert (Enable-, Disable-, MarkFailed-, default-Kind-, ctx-Cancel-, Listener-Fehler-Zweig) — alle sechs `PASS` mit `-race`, selbst ausgeführt |
| 2 | `Assembler` bekommt synchronisierte `AddBinding`/`RemoveBinding`; `a.tables`-Zugriff aus zwei Goroutinen race-frei | **erfüllt, eigenständig verifiziert (nicht nur aus Reports übernommen)** | Sperren in `lookupBinding`/`AddBinding`/`RemoveBinding` in dieser Sitzung testweise entfernt (Python-Patch), `go test -race -run TestAssemblerLiveReloadIsRaceFree` schlägt danach real fehl — `WARNING: DATA RACE` exakt an `mapper.go:378` (`lookupBinding`, `mapaccess2_faststr`) und `mapper.go:389` (`AddBinding`, `mapassign_faststr`), deckungsgleich mit Implementer- und Reviewer-Behauptung. Datei aus Backup wiederhergestellt, `git status` zeigt danach keinen Diff, derselbe Test läuft wieder grün |
| 3 | `wiring.go` baut `Assembler.tables` beim Start aus `cdc.source_table` (`TableActivationPort.List`); `CDC_TABLES` bleibt Seed | **erfüllt** | `activatedTableBindings` (`wiring.go:236-263`) gelesen: liest `activation.List` + `schemaStore.CurrentVersion`, `receive.Config.Tables` erhält diesen Rückgabewert (`wiring.go:406`), nicht `cfg.Tables` direkt; `cfg.Tables` bleibt ausschließlich Input der Erstaktivierungs-Schleife (`wiring.go:367-387`) |
| 4 | Realer End-zu-End-Nachweis Enable/Disable ohne Neustart | **erfüllt, dreifach selbst reproduziert** | `make test-integration` dreimal in Folge ausgeführt (Exit 0 in allen drei Läufen), jedes Mal beide Log-Zeilen „SQL-Administration Live-Reload-Beleg (enable)"/„(disable)" vorhanden; `grep -n "docker restart"` im Skript zeigt, dass die einzige `docker restart`-Stelle zum bestehenden Black-Box-CLI-Rundlauf gehört (Zeile 418), nicht zum SQL-Admin-Abschnitt (Zeilen 570–672) |
| 5 | `spec/architecture.md`-Sequenzdiagramm zu `LH-FA-CFG-001.a` in CLI/SQL getrennt, keine Wellen-/Slice-/ADR-Bezüge im Diagramm | **erfüllt** | Beide Diagramme gelesen (`spec/architecture.md:189-242`); Teilnehmer-Bezeichner ausschließlich `ARC-*`/Prosa (`ARC-005`, `ARC-003/002`, `ARC-004`, `ARC-006`, `ARC-007`, `PostgreSQL`) — kein `slice-*`/`welle-*`/`ADR-*` im Diagrammtext (Hard Rule 3.4) |
| 6 | `make gates` grün, `make test` (Race) und `make test-integration` dreimal grün | **erfüllt** | siehe §2 unten — alle real ausgeführt in dieser Sitzung |
| 7 | Review durchgeführt, Report liegt vor | **erfüllt** | `docs/reviews/review-slice-037.md` (1 HIGH/3 MEDIUM) und `docs/reviews/review-slice-037-fixrunde.md` (alle vier behoben) vorhanden, beide vollständig gelesen |
| 8 | Doku-Update bei öffentlichem Vertrag | **erfüllt** | `harness/README.md`-Diff zu `make test` (Race-Detector) und `make test-integration` (neuer Beleg) bestätigt (`git show 2a771bc --stat`) |
| 9–12 | Closure-Notiz, Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit (Modul 8); dieser Slice liegt noch in `in-progress/`, die vier Häkchen sind bewusst `[ ]`. Kein Verifikations-Gegenstand dieser Prüfung — der Planner trägt sie nach der Übergabe dieses Berichts nach |

**Zentrale `ADR-0050`-Eigenschaft (Fitness Function Zeile 2, `go test -race`):
bestätigt, eigenständig reproduziert** — sowohl im grünen Fall (Sperren
vorhanden, gesamte Suite grün) als auch im roten Kontrollfall (Sperren
entfernt, Data Race exakt an den behaupteten Stellen). Das ist die
stärkste verfügbare Prüfung: nicht nur „Test ist grün", sondern „Test
unterscheidet den korrekten vom fehlerhaften Zustand".

## 2. Eigene Reproduktion der Gate-/Testläufe

Alle vier Kommandos in dieser Sitzung selbst ausgeführt, keine Übernahme
aus Implementer- oder Reviewer-Bericht:

- **`make gates`** — grün: `baseline-verify` (54 Dateien), `d-check`
  (305 Dateien, 0 Befunde, `HEAD`-Stand), `commit-traceability`
  (`HEAD~5..HEAD`, 5 Commits, „Betreffs ohne Struktur-ID" — OK),
  `a-check` (0 Befunde).
- **`make test`** (Race-Detector, volle Suite) — grün, alle Pakete `ok`,
  einschließlich `internal/bootstrap` und
  `internal/adapters/driving/replication/mapper`; zusätzlich gezielt
  mit `-v` gegen die Fixrunden-Testdateien gelaufen:
  `internal/bootstrap/administration_internal_test.go` (6/6 `PASS`) und
  `internal/adapters/driven/postgresstorage/administration_backoff_internal_test.go`
  (6/6 Unterfälle `PASS`).
- **Race-Freiheit, Kontrollprobe** — siehe Tabelle oben, Punkt 2: Sperren
  testweise entfernt → real rot mit `DATA RACE` an den behaupteten
  Stellen → Datei wiederhergestellt → real wieder grün. `git status` vor
  und nach dem Eingriff zeigt keinen verbleibenden Diff.
- **`make test-integration`** — **dreimal in Folge**, jedes Mal Exit 0,
  jedes Mal beide SQL-Administration-Live-Reload-Belegzeilen (enable +
  disable) im Log, jedes Mal die vollständige `TestMVP*`-Suite grün
  (u. a. `TestMVPSchemaChangeAddColumn`, `TestMVPHeartbeatHealthy`,
  `TestMVPSchemaChangeIncompatibleTypeChange`). Nebenwirkung
  `tools/schema/plan.yaml` (von `schema-rollout` neu geschrieben) nach
  jedem Lauf auf `HEAD` zurückgesetzt (`git status` zeigt am Ende der
  Sitzung keinen Diff).

Kein Kommando lief nur einmal und wurde für „genügend" erklärt — die
DoD verlangt „dreimal in Folge" für `make test-integration`, und genau
das wurde erfüllt.

## 3. Hard Rule 3.4 (Architektur sprach-/meilensteinfrei)

Bestätigt am Diagrammtext selbst (§1, Punkt 5), nicht nur am
Review-Zitat: beide neuen Sequenzdiagramme unter `LH-FA-CFG-001.a`
referenzieren ausschließlich `ARC-*`-Kennungen. Der neue Teilnehmer
„Administrations-Hintergrundzug (`ARC-007`)" liegt korrekt unter `ARC-007`
(Composition Root/Bootstrap), da `runAdministration` in
`internal/bootstrap` lebt — konsistent mit der Schicht-Zuordnung der
übrigen Hintergrund-Züge (`runHeartbeat`, `runWALRetentionCheck`).

## 4. Welle-12-Zwischenstand (Zusatzinformation für den Planner, kein
Verifikations-Gegenstand dieses Slice)

- `slice-036` liegt in `done/` (`docs/plan/planning/done/slice-036-antragsqueue-sql-funktionen.md`).
- `slice-037` liegt noch in `in-progress/` — DoD inhaltlich vollständig
  erfüllt (siehe §1), die vier Planner-Closure-Punkte (Closure-Notiz,
  Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen) stehen noch
  aus und sind Gegenstand des nächsten Rollenwechsels (Planner).
- `slice-038` (CLI-Diagnose, `LH-FA-SST-003`) ist als letzte Zeile der
  Welle-12-Slice-Liste benannt (`docs/plan/planning/welle-12.md:93`),
  aber noch nicht als eigene Datei im Lifecycle angelegt (`open/`,
  `next/`, `in-progress/`, `done/` — keine trifft).
- `docs/plan/planning/next/slice-039-ci-workflow-dependabot.md`
  (`ADR-0051`) ist unabhängig von `welle-12` und nicht Teil dieser
  Prüfung.
- **Beobachtungs-Register:** `BEO-PGC/verwaltung-keine-sql-administration`
  steht weiter `offen` mit Ausgang „weiter offen → Welle 12" (Zähler
  0×, `evidence/` leer) — konsistent mit §8 des Slice-Plans. Die
  Auflösung dieser Beobachtung ist Sache der Welle-Closure, nicht
  dieses Einzel-Slice.

Mit `slice-037` inhaltlich erfüllt (Closure-Formalitäten ausstehend) und
`slice-036` bereits `done/`, fehlt für `welle-12`s Closure-Trigger
(„der Nachweis, dass eine über SQL beantragte Aktivierung/Deaktivierung
den laufenden Erfassungspfad real erreicht") inhaltlich nichts mehr —
dieser Nachweis liegt bereits vor (§1, Punkt 4). Formal offen bleiben:
`slice-037`s Planner-Closure und `slice-038` vollständig.

## Verdikt

**DoD-Konformität: bestätigt.** Alle acht inhaltlichen DoD-Punkte (1–8)
sind durch eigene Reproduktion gedeckt, nicht nur durch Bericht-Zitat.
Die vier verbleibenden Punkte (9–12) sind korrekt der Planner-Rolle
vorbehalten und kein Verifikations-Mangel.

**Zentrale `ADR-0050`-Eigenschaft (Race-Freiheit): bestätigt**, mit
eigenem Kontrollversuch (rot ohne Sperren, grün mit Sperren) — die
stärkste in dieser Sitzung mögliche Prüfung.

**End-zu-End-Beleg „ohne Neustart": bestätigt**, dreifach real
reproduziert, kein `docker restart` im betreffenden Abschnitt.

**Hard Rule 3.4: bestätigt** am Diagrammtext selbst.

**Übergabe an Planner:** Dieser Bericht bestätigt DoD-Konformität für
die Closure-Entscheidung. Kein offener Befund, keine Rückführung
erforderlich. Kein Validator-Zug ausgelöst — `slice-037` ist kein
MVP-Meilenstein-Slice im Sinn von Modul 8.

---

*Dieser Bericht ist ein Lauf-Beleg (Modul 11) und wird über Läufe hinweg
nicht wieder gelesen.*
