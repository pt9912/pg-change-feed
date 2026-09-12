# Verifier-Report: slice-035 — 2026-09-13

**Review-Art:** Verifikation — *wogegen*: Slice-Plan §1 (Ziel/Abgrenzung),
§2 (Definition of Done, aktueller Stand nach `9d6dfd8`), §6 (Risiken,
noch ohne Ausgang), §8 (Sub-Area-Sichtung), sowie DoD-/Spec-Konformität
gegen [`LH-FA-ADM-002`](../../spec/lastenheft.md), [`LH-FA-ADM-005`](../../spec/lastenheft.md)
(`spec/lastenheft.md`) und den Review-Report
[`review-slice-035.md`](review-slice-035.md) (1 LOW, Verdikt
merge-blockierend: keins). Zusätzlich, außerhalb der formalen Slice-DoD:
eine Vorab-Prüfung der Closure-Trigger-Reife von `welle-11` als Ganzes
(§5 des Auftrags), da `slice-035` ihr letzter Slice ist.

**Grundsatz:** Keine Behauptung übernommen — jeder Beleg unten wurde in
diesem Lauf selbst gelesen oder ausgeführt: vollständiger Slice-Plan
(§1–§8), Review-Report (Volltext), Testdatei-Ausschnitt
`TestMVPHeartbeatHealthy` (Zeilen 747–789), der neue
Verarbeitungsrückstand-Abschnitt in `tools/harness/run-integration-tests.sh`
(Zeilen ~485–567), das `-run`-Muster (Zeile 241), `spec/lastenheft.md`
(`LH-FA-ADM-002`/`005`-Akzeptanzkriterien), `welle-11.md` vollständig,
`make gates` einmal und `make test-integration` **dreimal in Folge**
selbst gestartet, `git status`/`git diff` am Ende sauber.

**Gegenstand:** `aa88ed5` (neue Testfälle `TestMVPHeartbeatHealthy` und
neuer Rückstands-Beleg-Abschnitt), `9d6dfd8` (DoD-Häkchen, Plan-Nachzug
§3), `6cce840` (Review-Report, 1 LOW).

**Modell:** Claude Sonnet 5 · **Datum:** 2026-09-13

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-035-black-box-observability-konsolidierung.md`
  (vollständig, §1–§8, aktueller Stand)
- `docs/reviews/review-slice-035.md` (1 LOW, kein Merge-Blocker)
- `spec/lastenheft.md` (`LH-FA-ADM-002`, `LH-FA-ADM-003`, `LH-FA-ADM-004`,
  `LH-FA-ADM-005`, vollständige Akzeptanzkriterien-Blöcke)
- `test/integration/integration_test.go` (`TestMVPHeartbeatHealthy`
  Zeile 747–789, `awaitHeartbeatErrorClass` Zeile 791ff. als Referenz
  für den bereits belegten Fehlerfall)
- `tools/harness/run-integration-tests.sh` (`-run`-Muster Zeile 241,
  Rückstands-Beleg-Abschnitt Zeile ~485–567, Lasttest-Beleg-Block
  Zeile 244ff.)
- `docs/plan/planning/welle-11.md` (vollständig, insbesondere §3
  Closure-Trigger)
- `docs/plan/planning/done/slice-034-black-box-status-liste.md` und
  `docs/reviews/verify-slice-034.md` (Nachweis, dass der
  `cdc.active_tables`-Closure-Trigger-Punkt bereits real erfüllt ist)
- `docs/plan/planning/observations/BEO-PGC/test-runner-stiller-ausschluss/`
  (Zähler-Stand, `evidence/`)
- `Makefile` (Ziel-Existenz-Prüfung für `doc-commits`/`doc-immutable`)

---

## Sensor-Belege (in diesem Lauf selbst gefahren)

| Sensor | Ausgabe-Kernaussage | Exit |
|---|---|---|
| `make gates` | `baseline-verify`: v6.5.0, 54 Dateien OK · `d-check` Standardlauf: 291 Dateien, 0 Befunde · `d-check --range HEAD~5..HEAD` (Modul `commits`): 291 Dateien, 0 Befunde · `commit-traceability.sh`: OK, 5 Commits, Betreffs ohne Struktur-ID · `a-check`: 0 Befunde | **0** |
| `make test-integration` (Lauf 1/3) | alle 8 Compose-Go-Testfälle PASS, `TestMVPHeartbeatHealthy` PASS (0,01s), Rückstands-Beleg „Rückstand vor der zweiten Bestätigung 1464, danach 0“, `TestMVPSchemaChangeIncompatibleTypeChange` PASS | **0** |
| `make test-integration` (Lauf 2/3) | alle 8 PASS, `TestMVPHeartbeatHealthy` PASS (0,01s), Rückstand 1464 → 0 | **0** |
| `make test-integration` (Lauf 3/3) | alle 8 PASS, `TestMVPHeartbeatHealthy` PASS (0,01s), Rückstand 1464 → 0 | **0** |
| `git status`/`git diff` (am Ende dieses Laufs) | einzige Restspur `tools/schema/plan.yaml` (erwartetes d-migrate-Rollout-Artefakt, bereits im Review-Report vermerkt) — per `git checkout --` zurückgesetzt, danach sauber | — |

**Zu `make doc-commits`/`make doc-immutable`:** Beide Targets existieren
in diesem Repos `Makefile` **nicht** (eigene Prüfung: `grep -n
"^[a-zA-Z_-]*:" Makefile`) — nach `AGENTS.md` §4 werden nur real
existierende Targets genannt. Die inhaltliche Traceability-Deckung läuft
über `commit-traceability` (Teil von `make gates`, oben grün); eine
dedizierte ADR-Immutabilitäts-Prüfung ist hier ohne Berührungspunkt, da
keiner der drei slice-035-Commits eine ADR-Datei ändert (eigene `git
show --stat` über alle drei Commits, s. u.).

`make test-integration` wurde entgegen der bloßen Implementer-/
Reviewer-Behauptung in diesem Lauf **selbst dreimal in Folge**
ausgeführt (nicht übernommen) — das ist die häufigste Verifier-Lücke,
und genau sie wird hier geschlossen. Der Rückstandswert-Wechsel
(1464 → 0) war in allen drei Läufen identisch reproduzierbar.

## Prüfpunkt 1 — `LH-FA-ADM-002` (Betriebsstatus, Happy Path) real erfüllt?

**Ja.** `TestMVPHeartbeatHealthy` (`integration_test.go:761–789`) liest
`cdc.heartbeat` real über `pgxpool` gegen die reale Sicht (kein
Go-Adapter-Import), pollt bis zu 30s auf `error_class = ''` **und**
`age_seconds < 15`. Die Schwelle 15s ist eigenständig gegen
`internal/bootstrap/wiring.go` nachvollzogen:
`heartbeatInterval = 5 * time.Second`, `heartbeatStaleAfter = 3 *
heartbeatInterval` = 15s — derselbe Wert, den der Compose-Healthcheck
bereits gegen dieselbe Zeile prüft, keine neu erfundene Zahl. Der
Testfall ist im `-run`-Muster des ersten `go test`-Aufrufs enthalten
(Zeile 241) und lief in allen drei eigenen Reproduktionen PASS (0,01s
je Lauf) — vor `TestMVPSchemaChangeIncompatibleTypeChange`, der
`error_class` dauerhaft setzt. Damit ist der Happy-Path-Wortlaut
(„Given ein betriebsbereites System, when der Status abgefragt wird,
dann wird der aktuelle Betriebszustand gemeldet“) real black-box belegt.
Boundary (`LH-FA-ADM-003`, Fehlerzustand) ist außerhalb des §1-Scopes
dieses Slice und bereits über `slice-033`s `awaitHeartbeatErrorClass`
belegt — eigenständig im Code bestätigt (Zeile 798ff., referenziert
dieselbe `cdc.heartbeat`-Projektion).

## Prüfpunkt 2 — `LH-FA-ADM-005` (Verarbeitungsrückstand) real erfüllt?

**Ja — und die im Review geprüfte Plan-Präzisierung trägt.** Eigenständig
nachvollzogen (nicht aus Review-Report übernommen, sondern selbst am
Runner-Skript und Testlauf-Output geprüft): Der neue Bash-Abschnitt
(`run-integration-tests.sh:485–567`) registriert einen Consumer, lässt
eine erste Änderung real erfassen, bestätigt sie extern (`docker exec
… acknowledge-consumer`), lässt eine zweite Änderung real erfassen und
liest `latest_commit_position - acknowledged_position` aus
`cdc.consumer_status` — in allen drei eigenen Läufen konsistent `1464`
vor und `0` nach der zweiten Bestätigung. Das deckt beide
Lastenheft-Kriterien: Happy Path („Given unbestätigte Changes jenseits
der bestätigten Position eines Consumers, … dann ist er erkennbar“) über
den Zustand vor der zweiten Bestätigung, Boundary („Given kein
Rückstand, … dann zeigt sie das auch“) über den Zustand danach.

Die im Review bereits geprüfte Abweichung von der ursprünglichen
DoD-Wortwahl („niemals bestätigt“ → Rückstand über eine bereits einmal
bestätigte Position) habe ich **eigenständig erneut** anhand des
Lastenheft-Wortlauts nachvollzogen, nicht nur den Review-Befund
übernommen: `LH-FA-ADM-005`s Happy-Path-Satz selbst spricht von Changes
„jenseits der bestätigten Position eines Consumers“ — das setzt
sprachlich bereits eine bestätigte Position voraus. Ein Consumer, der
*nie* bestätigt hat, hat begrifflich keine „bestätigte Position“, gegen
die etwas „jenseits“ liegen könnte; der real umgesetzte Ablauf trifft
den Wortlaut damit exakter als die ursprüngliche Plan-Paraphrase. Kein
Akzeptanzkriterium bleibt unbelegt.

## Prüfpunkt 3 — Eigene Reproduktion, Stabilität über drei Läufe

**Bestanden.** Drei unabhängige `make test-integration`-Läufe (siehe
Sensor-Tabelle) lieferten identische Ergebnisse: `TestMVPHeartbeatHealthy`
PASS in 0,01s je Lauf, Rückstandswert-Wechsel 1464 → 0 in allen drei
Läufen. Beide §6-Risiken (Timing-Flakiness beim Rückstandstestfall,
Testreihenfolge-Interferenz beim Heartbeat-Testfall) sind in keinem der
drei Läufe eingetreten — das ist ein starkes, aber kein hinreichendes
Indiz für „entfallen“ (drei Läufe sind kein Beweis der Abwesenheit von
Flakiness); das Urteil über den konkreten Risiko-Ausgang bleibt beim
Planner (Modul 5 — Ausgang ist Menschen-Urteil, nicht Sensor-Ergebnis).

## DoD-Prüfung (Slice-Plan §2, Punkt für Punkt, Stand nach `9d6dfd8`)

| # | DoD-Punkt | Verdikt | Beleg-Kernaussage |
|---|---|---|---|
| 1 | `LH-FA-ADM-005` erfüllt (mit Präzisierung) | **bestätigt** | Prüfpunkt 2 — eigenständig gegen Lastenheft-Wortlaut und Testlauf-Output nachvollzogen |
| 2 | `LH-FA-ADM-002` erfüllt (Happy Path, Schwelle 15s) | **bestätigt** | Prüfpunkt 1 — Schwelle eigenständig gegen `wiring.go` nachgerechnet |
| 3 | `make gates` grün, `make test-integration` dreimal in Folge grün | **eigenständig real reproduziert** | Sensor-Belege oben — nicht aus Implementer-/Reviewer-Bericht übernommen, selbst dreimal gefahren |
| 4 | Review durchgeführt, Report liegt vor | **bestätigt** | `docs/reviews/review-slice-035.md` (`6cce840`), 1 LOW, Verdikt „keins“ merge-blockierend — Checkbox im Plan-Text selbst noch `[ ]`, korrekt (Planner-Closure-Schritt zieht sie nach) |
| 5 | Doku-Update, falls öffentlicher Vertrag berührt | **bestätigt, entfällt korrekt** | reine Testdatei-/Runner-Skript-Änderung, keine API-/CLI-/Schema-Änderung — eigene `git show --stat` bestätigt |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **noch offen — korrekt, Planner-Closure-Schritt** | §7 ist Platzhalter |
| 7 | Reconciliation-Register, falls einschlägig | **entfällt korrekt** | `../reconciliation.md` existiert nicht (Greenfield, `harness/conventions.md` MR-000) |
| 8 | Beobachtungs-Register fortgeschrieben | **noch offen — korrekt, Planner-Closure-Schritt** | §8 zitiert zwei existierende Verzeichnisse korrekt (eigene Prüfung: beide vorhanden, Zähler-Stände stimmen mit Plan-Angabe überein) |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **noch offen — korrekt, Planner-Closure-Schritt** | beide §6-Risiken tragen noch `<bei Closure einzutragen>`; Prüfpunkt 3 liefert Evidenz (kein Auftreten in drei Läufen), Urteil bleibt beim Planner |
| 10 | Drei Paarungen getragen | **noch offen — korrekt, Welle-11-Closure vorbehalten** | Slice-Kopf nennt `welle-11` als Träger; Paarungen sind laut Modul 6 Sache der nächsten Welle-Closure |

**Zwischenstand:** Die drei Kern-DoD-Punkte, die der Implementer als
erledigt beansprucht (1–3), sind eigenständig real nachgeprüft und
tragen vollständig, ohne Einschränkung. Die fünf offenen Punkte (6, 8, 9,
10 sowie das noch fehlende Checkbox-Häkchen für Punkt 4 im Plan-Text
selbst) sind korrekt der Planner-Closure vorbehalten (Modul 5/6/8) und
kein Implementer- oder Reviewer-Versäumnis.

## Negativbefunde

- geprüft, ohne Befund: **Kein ADR-Inhalt geändert** — `git show --stat`
  über alle drei slice-035-Commits (`aa88ed5`, `9d6dfd8`, `6cce840`):
  keiner berührt eine Datei unter `docs/plan/adr/` (Hard Rule 3.5
  gewahrt).
- geprüft, ohne Befund: **Kein öffentlicher Vertrag geändert** — Diff
  beschränkt sich auf `test/integration/integration_test.go`,
  `tools/harness/run-integration-tests.sh` und den Review-Report; kein
  Pfad unter `internal/` oder `cmd/`, also kein Image-Rebuild-Trigger
  (deckt sich mit Reviewer-Negativbefund, eigenständig anhand `git show
  --stat` nachvollzogen).
- geprüft, ohne Befund: **Traceability** — alle drei Commit-Betreffs
  tragen `LH-FA-ADM-002`/`LH-FA-ADM-005`, kein `SPEC-*`/`ARC-*` im
  Betreff; `commit-traceability` (Teil von `make gates`) lief in diesem
  Lauf über `HEAD~5..HEAD` grün.
- geprüft, ohne Befund: **§8-Register-Zitate existieren und stimmen** —
  `docs/plan/planning/observations/BEO-PGC/test-runner-stiller-ausschluss/`
  (Zähler 1×, `evidence/slice-033.md`) und
  `.../verwaltung-keine-sql-administration/` beide vorhanden, konsistent
  mit dem im Slice-Kopf §8 genannten Stand.
- geprüft, ohne Befund: **`TestMVPHeartbeatHealthy` ist im `-run`-Muster
  erfasst** — `run-integration-tests.sh:241` enthält den Testfall; ohne
  diese Zeile liefe er nie in `make test-integration`
  (`BEO-PGC/test-runner-stiller-ausschluss`-Risiko korrekt vermieden,
  nicht wiederholt).
- geprüft, ohne Befund: **`make gates`** — ein eigener Lauf, 0 Befunde
  (291 Dateien Standardlauf, 291 Dateien Range-Lauf, `commit-traceability`
  OK, `a-check` 0 Befunde).
- geprüft, ohne Befund: **`make test-integration` dreimal in Folge** —
  drei eigene, unabhängige Läufe, alle acht Go-Testfälle des ersten
  Compose-Aufrufs PASS in jedem Lauf, `TestMVPHeartbeatHealthy` PASS in
  allen drei (0,01s), Rückstandswert-Wechsel 1464 → 0 in allen drei
  identisch, `TestMVPSchemaChangeIncompatibleTypeChange` PASS im
  jeweils letzten Aufruf.
- geprüft, ohne Befund: **Arbeitsverzeichnis am Ende dieses Laufs** —
  einzige Restspur war das erwartete `tools/schema/plan.yaml`-Artefakt
  (Rollout-Nebenwirkung von `make test-integration`, bereits im
  Review-Report vermerkt), per `git checkout --` zurückgesetzt;
  `git status --porcelain` danach leer.

## Zusätzliche Prüfung: Closure-Trigger-Reife von `welle-11` (kein formaler Teil der Slice-DoD)

Regeln: Baseline-Regelwerk `modul-06-roadmap.md` §Wellen-Closure-Prozedur,
Schritt 1 (Trigger prüfen). `welle-11.md` §3 nennt fünf Kriterien:

| Closure-Trigger (`welle-11.md` §3) | Stand |
|---|---|
| Alle Slices dieser Welle liegen in `done/` | **noch nicht** — `slice-034` liegt bereits in `done/`; `slice-035` liegt noch in `in-progress/` und wird erst nach dieser Verifikation vom Planner geschlossen |
| `make gates` grün | **erfüllt** — eigener Lauf oben, 0 Befunde |
| Neuer Black-Box-Testfall liest `cdc.active_tables` real über SQL, belegt `LH-FA-CFG-003`/`004` | **erfüllt** — durch `slice-034` (`TestMVPActiveTablesViewMatchesActivationState`), bereits eigenständig verifiziert in `docs/reviews/verify-slice-034.md` |
| Mindestens ein neuer Black-Box-Testfall je Observability-Signal (`LH-FA-ADM-002`…`005`) | **erfüllt** — `002` durch `TestMVPHeartbeatHealthy` (dieser Slice), `003` bereits vorher durch `slice-033`s `awaitHeartbeatErrorClass` (eigenständig im Code nachvollzogen, Zeile 798ff.), `004` bereits vorher durch den `cdc_capture_lag`-Lasttest-Beleg (eigenständig im Skript nachvollzogen, Zeile 244ff.), `005` durch den neuen Rückstands-Beleg-Abschnitt (dieser Slice) |
| Closure-Notiz in `welle-11-results.md` | **noch nicht** — Datei existiert noch nicht (erwartet: entsteht im Rahmen der Welle-Closure-Prozedur, Schritt 3, nach Abschluss aller Slices) |

**Einschätzung für den Planner:** Inhaltlich sind alle fünf
Closure-Trigger-Kriterien von `welle-11` bereits jetzt entweder erfüllt
oder unmittelbar erfüllbar — die beiden noch offenen Punkte
(„alle Slices in `done/`“, Closure-Notiz) sind strukturell genau die
beiden Schritte, die die Welle-Closure-Prozedur selbst als *nächste*
Handlung vorsieht (erst `slice-035` nach `done/`, dann Schritt 3 der
Prozedur inkl. Closure-Notiz), kein zusätzlicher Arbeitsbedarf über die
reguläre Closure hinaus. `welle-11` ist damit — vorbehaltlich der
Planner-Closure-Schritte an `slice-035` selbst (§6-Risiko-Ausgänge,
Beobachtungs-Register, Closure-Notiz) — bereit für ihre eigene
Closure-Prozedur unmittelbar im Anschluss.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 (eigene Befunde) |
| INFO | 0 |

**Zusammenfassung DoD:** Die drei implementierten Kern-Punkte (1–3) sind
eigenständig real nachgeprüft und tragen vollständig. Die fünf noch
offenen Punkte (4-Häkchen im Plan-Text, 6, 8, 9, 10) sind korrekt der
anstehenden Planner-Closure vorbehalten und kein Implementer-/
Reviewer-Versäumnis. Kein DoD-Defekt im Sinn eines unbelegten
„bestätigt“-Punkts.

## Verdikt

**DoD-Konformität (soweit vom Implementer beansprucht): bestätigt, ohne
Einschränkung.** Beide neuen Testfälle liegen real gegen die
Compose-Umgebung, lesen ausschließlich über externe SQL-Sichten (kein
Go-Adapter-Import), und beide Lastenheft-Anforderungen
(`LH-FA-ADM-002` Happy Path, `LH-FA-ADM-005` Happy Path + Boundary) sind
damit jetzt real black-box belegt. Die im Review geprüfte
Plan-Präzisierung bei `LH-FA-ADM-005` wurde in diesem Lauf eigenständig
erneut gegen den Lastenheft-Wortlaut geprüft und trägt.

**Eigene Reproduktion:** `make gates` einmal, `make test-integration`
**dreimal in Folge**, alle vier Läufe grün, real in diesem Kontext
ausgeführt — nicht aus dem Implementer- oder Reviewer-Bericht
übernommen. Der Rückstandswert-Wechsel (1464 → 0) und
`TestMVPHeartbeatHealthy`s Erfolg wurden dabei in jedem der drei Läufe
einzeln bestätigt.

**Plan-vs-Code-Diff:** Der Diff (`aa88ed5`) hält sich an §3 des
Slice-Plans (Testdatei-Ergänzung + Runner-Skript-Nachzug), inklusive der
im Plan-Nachzug ehrlich benannten Implementer-Entscheidungen (Bash statt
Go-Test wegen fehlendem Docker-Socket, 15s-Schwelle aus Produktionswert,
Rückstandsszenario-Präzisierung). Keine unangekündigte Schicht- oder
Umfangs-Erweiterung; §1s Abgrenzung (SQL-Administration, CLI-Diagnose,
Schwellenwert-Klassifikation, `LH-FA-ADM-003`/`004` als bereits
gedeckt) wird nicht berührt.

**Closure-Bereitschaft:** Die Verifikation der real gelieferten Punkte
ist abgeschlossen und positiv. Vor dem `git mv` nach `done/` fehlen noch
die Planner-Closure-Schritte: DoD-Häkchen 4 im Plan-Text nachziehen,
§7-Closure-Notiz schreiben, Beobachtungs-Register fortschreiben
(voraussichtlich „keine Beobachtung angefallen“, da beide §8-Treffer
unter der Schwelle bleiben), §6-Risiko-Ausgänge eintragen (Prüfpunkt 3
liefert Evidenz für „entfallen“, Urteil bleibt beim Planner) — die drei
Paarungen bleiben regulär der `welle-11`-Closure vorbehalten, die
unmittelbar im Anschluss ansteht (siehe §Zusätzliche Prüfung oben).

**Übergabe:** Bericht an den Planner. Keine Reparaturen — Slice-Plan,
Register und Roadmap wurden von diesem Lauf nicht verändert
(`git status`/`git diff` am Ende sauber, bis auf den erwarteten,
zurückgesetzten `tools/schema/plan.yaml`-Rollout-Rest).

---

**Gate-Beleg:** `make gates` einmal in diesem Lauf ausgeführt, Exit 0, 0
Befunde (`baseline-verify`: v6.5.0, 54 Dateien OK; `d-check`
Standardlauf: 291 Dateien, 0 Befunde; `d-check --range HEAD~5..HEAD`
Modul `commits`: 291 Dateien, 0 Befunde; `commit-traceability.sh`: OK, 5
Commits; `a-check`: 0 Befunde). `make test-integration` dreimal in Folge
ausgeführt, Exit 0 in jedem Lauf, `TestMVPHeartbeatHealthy` PASS in
allen drei, Rückstandswert-Wechsel 1464 → 0 in allen drei identisch.
`make doc-commits`/`make doc-immutable` existieren in diesem Repo nicht
(eigene Makefile-Prüfung) — Traceability-Deckung läuft über
`commit-traceability` (Teil von `make gates`), ADR-Immutabilität ist
hier ohne Berührungspunkt (kein Commit ändert eine ADR-Datei). `git
status`/`git diff` am Ende dieses Laufs sauber (keine
Arbeitsverzeichnis-Änderung durch die Verifikation selbst).
