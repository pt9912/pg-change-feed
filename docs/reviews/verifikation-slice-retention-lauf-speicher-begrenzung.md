# Verifikations-Report: slice-retention-lauf-speicher-begrenzung — 2026-09-25

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich, Entscheidungs-Konformität, Plan-vs-Code-Diff,
Gates, Eingabeseiten-Mutationen und Nachmessung. Review-Artefakt des Reviewers:
[`review-slice-retention-lauf-speicher-begrenzung.md`](review-slice-retention-lauf-speicher-begrenzung.md); Formvorbild
dieses Reports: [`verifikation-slice-backfill-speicher-untersuchung.md`](verifikation-slice-backfill-speicher-untersuchung.md).
Die vendored Baseline trägt kein eigenes Verifikations-Template (`.harness/baseline/v6.9.0/templates/docs/reviews/`
enthält nur `review-report.template.md`), ein Skill `.harness/skills/verifier.md` liegt nicht vor; der Report folgt dem
Formvorbild.

**Gegenstand:** Slice-Plan `slice-retention-lauf-speicher-begrenzung` (ohne Welle, Lifecycle `in-progress`), Stand `HEAD` =
`1d169429` (gepusht), Diff-Range `8cd39719..HEAD` (23 Dateien, +2301/−163). Slice-Inhalt sind `6d8e5d3f` (Umsetzung),
`d49e248f` und `7572051c` (Fortschrittsprüfung), `0e6b1b30` (Plan-Nachzug), `2636ca46` (Nachmessung, Messbericht, Handbuch,
Pflichtenheft), `a583afe5` (Gate-Zeilen), `b2484add` (DoD-Haken), `f1aeeae1` (Fixrunde) sowie die drei Plan-Commits
`e486beb8`, `8a68d2b2`, `ff387a69` (`e486beb8` und `ff387a69` sind reine Renames, `git show --stat -M`: 0 Einfügungen,
0 Löschungen). Nicht Slice-Inhalt sind `c21c43b9` (Review-Report) und `1d169429` (Architect-Verdikt). Bezug:
[`LH-FA-RET-002`](../../spec/lastenheft.md), [`LH-FA-RET-003`](../../spec/lastenheft.md),
[`LH-FA-RET-004`](../../spec/lastenheft.md), [`LH-FA-CAP-009`](../../spec/lastenheft.md),
[`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md).
Dieser Lauf ändert weder Produktionscode noch Plan, Spec oder Doku; er schreibt nur diesen Report. Alle Mutationen
liefen im Arbeitsbaum und sind per `git checkout` zurückgenommen (`git status --short` danach: nur diese Datei);
nichts gepusht, nichts getaggt.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei des Scratchpads; der Exit-Code wurde gesondert gesichert und danach gelesen. Schwere
Docker-Läufe liefen nacheinander; `free -m` (verfügbar) vor den Läufen 18,4 bis 19,1 GB.

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make gates` (Stand `1d169429`) | **EXIT=0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `coverage-gate: OK — Coverage 83.40% erfüllt Schwelle 80%` · `d-check: 1177 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"` · `generated-sync: OK` · a-check `gesamt: 0 Befund(e)` |
| `make test` | **EXIT=0** | 42 Zeilen `ok`, keine `FAIL` |
| `make a-check` | **EXIT=0** | `gesamt: 0 Befund(e)` |
| `make docs-check` | **EXIT=0** | `d-check: 1177 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-commits RANGE=8cd39719..HEAD` | **EXIT=0** | `d-check: 1177 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=8cd39719..HEAD` | **EXIT=0** | `d-check: 1177 Datei(en) geprüft, 0 Befund(e)` |
| `make test-store` (PostgreSQL 18, 45 s) | **EXIT=0** | `DB-Adapter-Coverage: 82.61% (gedeckt 879 von 1064 Statements; Profile gemergt: store,replication)` · `db-coverage: OK — DB-Adapter-Coverage 82.61% erfuellt Schwelle 80%` |
| `tools/bench-scaling.sh` (verkürzt, 51 s) | **EXIT=0** | `cdc_capture_lag=1.000198s` (klein), `1.004965s` (mittel), `1.008498s` (gross) gegen die Grenze von 60 s |
| `tools/bench-backfill-memory.sh` (Nachmessung, 751 s) | **EXIT=0** | Abschnitt 4 |
| CI am Commit `1d169429` (`gh run list --commit 1d169429b72ae336d20d2fda7dfb718e4b96663d`, nach Abschluss abgefragt) | success | `examples` Lauf 36185281841, `ci` Lauf 36185281845, `e2e` Lauf 36185281855 — alle `completed success` |

Hygiene: dangling-Volumes (`docker volume ls -q -f dangling=true | wc -l`) vor den Läufen **34**, nach allen Läufen
**34**; kein `prune`; kein Container `cdc-store-test-pg` und kein Netz `cdc-store-test` oder `pgc-bench-*` zurück (Bench-
Skripte über ihren Abräum-Pfad, Store-Tier-Mutationsläufe mit `docker rm -fv`); `tools/schema/plan.yaml` und
`tools/schema/down.sql` nach jedem Rollout-Lauf per `git checkout` zurückgenommen.

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Zeilen des Plans gezählt: 13 `[x]`, 5 `[ ]`, zusammen 18.

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Umsetzung nach [`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md) (Port, `Run` je Seite, Adapter, kein Löschprädikat in SQL) | **erfüllt** | Diff gegen die sechs Festlegungen gelesen (Abschnitt 3); `git diff 8cd39719..HEAD` zu `queries.go`, `translate.go`, `store.go`, `changestore.go`, `service.go` |
| 2 | Unit-Tests (`make test`): Menge bei Begrenzung 1, 2, 3, N, größer N gegen handgeschriebene Erwartung; `limit`/`after` je Aufruf; Löschung je Seite; Fehler ab Aufruf `n`; kurze Seiten; bestehende Erwartungen | **erfüllt** | `make test` EXIT=0; `TestRunReleasesSameSetAtEveryPageSize` trägt die Erwartung `[c1 c4 c5 c7]` als Literal und sechs Seitengrößen; sieben Use-Case-Mutationen per Assertion rot (Abschnitt 5) |
| 3 | Store-Tier (`make test-store`): Abdeckung, Position/Zeit gleich `ReadChanges`, Cursor-Grenze, 25.000 Changes gegen unabhängige SQL-Zählung, `cdc_admin`-Login, Eingabeseiten-Mutationen | **erfüllt** | `make test-store` EXIT=0; sechs Tests vorhanden (`retentioncandidates_test.go`, `retention_pages_test.go`); vier Mutationen selbst rot gesehen (Abschnitt 5) |
| 4 | Gate-Zuordnung (Übersetzung in `sqlexec`, Methode in `postgresstorage`, Konstante in `queries`) | **erfüllt** | DB-Nenner 1064 und gedeckt 879 aus meinem `make test-store`-Lauf gedruckt (gleich der Zeile in `harness/sensors/db-adapter-coverage.md`); Nenner des Unit-Gates 2601 nachgezählt (Abschnitt 6) |
| 5 | Nachmessung liegt vor (`BENCH_MEM_STAGES=1000000`, drei Runs, gedruckte Zeilen, Host, Lauf, Zeile) | **erfüllt, mit V-3** | Messbericht Abschnitt 9 trägt zwei Läufe mit gedruckten Zeilen; mein eigener Lauf (Abschnitt 4) liegt im Band; das Ergebnis ist der Form nach erfüllt, der Zahl nach nicht (Abschnitt 4, V-3) |
| 6 | Regressions-Beleg der Live-Last (`tools/bench-scaling.sh`, `cdc_capture_lag` unter 60 s) | **erfüllt** | eigener Lauf EXIT=0, drei Stufen gedruckt (Abschnitt 1) |
| 7 | Retention-Lebenszyklus-Rundlauf von `make test-integration` unverändert grün | **erfüllt (CI)** | nicht selbst gefahren (Grenze, Abschnitt 8); `e2e` am Commit `1d169429` success, Rundlauf-Zeile „Retention-Beleg“ steht im Messbericht Abschnitt 9 („Tests und Gates“) für `b2484add`; `tools/harness/run-integration-tests.sh` ist im Diff unverändert |
| 8 | Träger nachgezogen (Handbuch, Pflichtenheft, Kommentar an `retentionInterval`, Version/Historie) | **erfüllt** | Abschnitt 6 |
| 9 | Bewertung der Richtgröße | **erfüllt** | Messbericht Abschnitt 7.2; Kopierraten 9.802 bis 11.056 Zeilen/s aus den gedruckten Dauern nachgerechnet, × 600 s = 5,88 bis 6,63 Millionen; mein Lauf 9.193 bis 10.612 Zeilen/s, × 600 s = 5,5 bis 6,4 Millionen (abgeleitet); `warn.go` im Diff unverändert |
| 10 | Release-Aussage in §7 | **erfüllt** | Reichweite `v0.1.0` bis `v0.1.2` (Abschnitt 7); Zahlen tragen Ursprung |
| 11 | §3.13-Suchlauf, Gefundenes und Nichtgefundenes, beide Stände | **erfüllt, mit V-2** | Abschnitt 6 |
| 12 | `make gates` grün, Exit ungefiltert gesichert | **erfüllt** | eigener Lauf EXIT=0 |
| 13 | Review durchgeführt, Report liegt vor | **erfüllt, Nachzug offen (V-1)** | Report vorhanden (0 HIGH, 3 MEDIUM, 3 LOW, 7 INFO); F-1 und F-2 in der Fixrunde behoben (Abschnitt 5); F-3 und F-13 vom Architect als akzeptierte Negative bewertet, der Nachzug in Plan und Handbuch steht aus (V-1) |
| 14 | Closure-Notiz mit Lerneintrag | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*“ |
| 15 | Reconciliation-Register | **korrekt offen** | Zeile „entfällt“ trägt die Begründung; das Häkchen setzt der Planner |
| 16 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht des Planners |
| 17 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | acht Risiken tragen „*(bei Closure)*“; zwei Ausgänge liegen als Verdikt vor (V-1) |
| 18 | Drei Paarungen | **korrekt offen** | Closure-Pflicht des Planners |

Kein `[x]` ohne Beleg; kein `[ ]`, das über die Rollen-Sequenz hinaus belegt wäre.

## 3. Entscheidungs-Konformität

| Entscheidung | Verdikt | Beleg |
|---|---|---|
| [`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md) Festlegung 1 (Lese-Vertrag am bestehenden Port, `RetentionCandidate` mit genau `ChangeID`, `Position`, `CommittedAt`) | konform | `changestore.go`: neue Methode und Typ am `ChangeStorePort`; kein neuer Port; `ReadChanges` unverändert |
| Festlegung 2 (Seitenvertrag: `limit` kleiner 1 → `ErrNonPositiveLimit`, leere Quelle → `ErrEmptyIdentifier`, Ordnung des Schlüssels, leere Seite ist das Ende) | konform | `store.go`; Port-Kommentar nennt Ordnung und Ende; Use Case beendet nur an der leeren Seite (`TestRunContinuesPastShortPages`) |
| Festlegung 3 (`PageSize` = 10.000 als Konstante) | konform | `service.go` `const PageSize = 10_000`; keine Konfiguration |
| Festlegung 4 (Positionen und Uhr einmal, je Seite ein `DeleteChanges` über die Freigaben der Seite, `Deleted` über alle Seiten, Fehler → leeres Ergebnis) | konform | `service.go` gelesen; `TestRunReadsPositionsOncePerRun`, `TestRunDeletesOnlyTheReleasedIDsOfEachPage`; Mutationen M3, M4 rot |
| Festlegung 5 (Regel bleibt Domain Policy, kein Löschprädikat in SQL) | konform | `SelectRetentionCandidates` trägt nur `t.source_id = $1 AND c.change_id > $2`; `AllowsDeletion` im Use Case je Kandidat |
| Festlegung 6 (nicht atomar über Seiten, Lauf verschiebt nur, löscht nie früher) | konform | Uhr vor der ersten Seite gelesen (Alter höchstens unterschätzt); `TestRetentionCandidateBehindCursorAppearsInNextRun`; Mutation `>=` rot |
| Kandidatenmenge gleich der bisherigen | konform | `SelectChanges` verbindet zusätzlich `cdc.source_table` per `JOIN`; die Spalte `source_table_id` ist durch den Fremdschlüssel gebunden, der `JOIN` filtert keine Zeile (Review-Befund nachgelesen, am Schema-Ziel konsistent); der 25.000-Changes-Test zählt unabhängig |
| [`ADR-0014`](../plan/adr/0014-retention-domain-policy.md) (Retention als Domain Policy) | konform | siehe Festlegung 5; `retention.go` im Diff unverändert |
| [`ADR-0009`](../plan/adr/0009-change-store-outbound-port.md), [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md), [`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md) (Fähigkeits-Port, Transport-Typ am Port, Ports referenzieren nur die Domain) | konform | `RetentionCandidate` trägt nur `model`-Typen; `make a-check` EXIT=0, `gesamt: 0 Befund(e)` |
| [`ADR-0007`](../plan/adr/0007-source-ack-outbound-port.md), [`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md) | nicht berührt | der Diff trägt weder Ack-Port noch SQL-Driving-Adapter; `tools/schema/nacharbeit-roles.sql` unverändert; das Lesen unter der `cdc_admin`-Login-Identität belegt `TestRetentionCandidatesRunUnderTheirRoles` (Rot bei gekürztem Grant laut Messbericht Abschnitt 5) |
| [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md), [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) (Trigger, Richtgröße) | konform | Richtgröße bewertet und nicht geändert (`warn.go` im Diff unverändert), keine Schwelle, kein neues Gate ([`AGENTS.md`](../../AGENTS.md) §3.6) |
| [`AGENTS.md`](../../AGENTS.md) §3.2, §3.3, §3.5, §3.7, §3.9 | konform | kein `nolint` im Diff; die zwei Lifecycle-Commits sind reine Renames, die Inhaltsänderung `8a68d2b2` steht getrennt; keine `Accepted`-ADR im Diff (`make doc-immutable` EXIT=0); Kommentare in `service.go` benennen Zustand und Vertrag, kein Konjunktiv über verworfene Alternativen (Fixrunde); alle Exit-Codes ungepiped |

**Plan-vs-Code:** jede der 23 Dateien gehört zu einer Zeile der Plan-Tabellen (§3 samt den Fixrunde-Zeilen) oder ist
Record (Review-Report, Architect-Verdikt, Messbericht, Plan selbst). Über den Plan hinaus liegt die Fortschrittsprüfung
`last == after` in `service.go` (im Plan §3 als Abweichung benannt, Beleg: Mutation `>=` ohne sie lief bis zur Zeitüberschreitung).

## 4. Nachmessung selbst gefahren

Ein Lauf `tools/bench-backfill-memory.sh` nach dem Vertrag von `harness/targets/bench-backfill.md`:
`BENCH_MEM_STAGES=1000000 BENCH_FEED_ENV='GODEBUG=gctrace=1' BENCH_FEED_DOCKER_ARGS='--memory 6g'`, drei Runs, frischer
Feed-Container je Run, Lauf `20260925T203047Z`, Exit 0, 751 s. Host: Linux 6.8.0-139-generic, Docker 29.8.1, 20 CPU,
33.362.599.936 Byte RAM; Feed-Image `sha256:8c8dea000428448c0ff22a95f66061b32634beff80eeef1cf5c70e94c4a82a40` — dieselbe
ID wie im Messbericht; der ausführbare Code seit diesem Image ist unverändert (`git diff 0e6b1b30..HEAD -U0` zu
`service.go` und `wiring.go`: nur Kommentarzeilen). Gedruckt (Run 1/2/3, Bestand in `cdc.change` vor dem Run 0 / 1.000.000 /
2.000.000):

| Größe | Run 1 | Run 2 | Run 3 |
|---|---|---|---|
| `memory.peak` seit Start des Feeds, nach dem Nachlauf, in MiB | 15,5 | 16,9 | 17,3 |
| `anon` Spitze im Run in MiB | 10,1 | 12,8 | 13,2 |
| Kopierdauer (Zeilen/s) | 94.229 ms (10.612) | 100.421 ms (9.958) | 108.778 ms (9.193) |
| Bereinigung gelaufen seit dem Start / fehlgeschlagen | 21 / 0 | 23 / 0 | 24 / 0 |
| Zeilen in `cdc.change` nach dem Nachlauf | 1.000.000 | 2.000.000 | 3.000.000 |
| `OOMKilled` / Neustarts | false / 0 | false / 0 | false / 0 |
| `file` am Run-Ende in MiB | 0,0 | 0,0 | 0,0 |

Einordnung (abgeleitet): die Spitze liegt bei 15,5 bis 17,3 MiB, innerhalb der Spanne 14,9 bis 17,6 MiB des Messberichts
und flach gegen die 1.082,7 bis 3.096,6 MiB der Untersuchung (Reihe B); die Takte laufen bis zum Ende weiter (21 bis 24,
keine fehlgeschlagene); der Anstieg zwischen 1.000.000 und 3.000.000 Changes ist mit 1,8 MiB (17,3 − 15,5) von der
Größenordnung des Messberichts (2,4 und 2,5 MiB). Der Wert von Run 1 (15,5) liegt 0,4 MiB über dem Band der zwei Läufe des
Implementers (14,9 bis 15,1), unter dem Wert des Review-Laufs bei gleicher Stufe (31,8 mit 16,5 MiB Seiten-Cache); die
Streuung derselben Stufe über die Läufe ist damit größer als die Bandbreite von zwei Läufen — ein Argument für die
Formulierung „nicht mit nennenswertem Betrag“ statt einer Zahl als Grenze.

**Nachgerechnet gegen die gedruckten Zeilen des Messberichts** (Zeilen extrahiert, dann gerechnet; alle bestätigt):
(1) 1.000.000 / 91,565 s = 10.921 Zeilen/s (gedruckt 10921); (2) 17,3 − 14,9 = 2,4 und 17,6 − 15,1 = 2,5 MiB;
(3) 2,5 × 1.048.576 / 2.000.000 = 1,31 Bytes je Change; (4) 1,3 / 1.628 = 0,08 % und 1,3 / 1.055 = 0,12 %
(1,03 × 1.024 = 1.055, 1,59 × 1.024 = 1.628); (5) Differenzen zu Reihe J (13,3; 13,2; 13,7): 1,6 bis 1,8, 3,3 bis 3,7,
3,6 bis 3,9 MiB; (6) Zähler „Bereinigung gelaufen“ 20/22/23 und 20/23/23; (7) GC-Läufe im Nachlauf 560 bis 2.467;
(8) `anon` Spitze Run 3 minus Run 1: 12,6 − 10,2 = 2,4 und 12,8 − 10,0 = 2,8 MiB; (9) Grundlinie Maximum der Runs 2 und 3:
10,2; 10,1; 10,5; 10,7; (10) Schritt 2.000.000 → 3.000.000: 0,7 und 0,8 MiB. Die Vergleichswerte aus der Untersuchung
stimmen mit den Zeilen dort: Reihe B 1.082,7 / 2.269,2 / 3.096,6, Reihe C 1.273,5 / 3.058,0 / 2.751,7, Reihe J
13,3 / 13,2 / 13,7 MiB.

**Plan-Klausel (F-2).** Die Auswertung steht in Messbericht Abschnitt 7.1: die Klausel löst der Sache nach nicht aus.
Gegen die Zeilen geprüft: Belege 1 (Takte 20 bis 24, keine fehlgeschlagene, nun auch in meinem Lauf), 2 (0,08 bis 0,12 %),
3 (die Grundlinie springt einmal von 5,1/5,4 auf 10,1 bis 10,7 MiB beim Wechsel von leerem zu gefülltem `cdc.change` und
ist danach flach; der Schritt 2.000.000 → 3.000.000 ist 0,7/0,8/0,2 MiB und liegt in der Streuung derselben Stufe von
0,2 bis 0,8 MiB) tragen. Der Wortlaut der Klausel („die Spitze wächst mit der Zahl der Changes“) ist an den Bereichen
formal erfüllt; der Plan benennt das ehrlich („dem Wortlaut nach steigen die Bereiche … ohne Überlappung“) und hält die
Entscheidung als die des Auftraggebers fest. Das ist eine dokumentierte Lesart, keine stille Anpassung — **erfüllt**.

## 5. Review-Findings nachgemessen und Eingabeseiten-Mutationen

**Findings des Reviews:**

| Finding | Was ich gemessen habe | Verdikt |
|---|---|---|
| F-1 (MEDIUM) Wortlaut „hängt nicht an der Zahl der Changes“ | `git grep -n -i 'nicht an der Zahl\|hängt nicht'` in `docs/user spec harness internal docs/plan/adr`: 7 Treffer am Stand `c21c43b9`, 5 im Arbeitsbaum; die verbliebenen: zwei fremde Gegenstände (`ADR-0065`, Handbuch „Dauerhaftigkeit“), `ADR-0124` Zeile 209 (`Accepted`, [`AGENTS.md`](../../AGENTS.md) §3.5), Handbuch Zeile 1679 in der neuen Form („nicht mit nennenswertem Betrag“), Pflichtenheft Zeile 673 (fremder Gegenstand). Handbuch Zeilen 667 und 1679, Pflichtenheft Zeilen 143 bis 145 tragen „nicht mit nennenswertem Betrag an der Zahl“; der Godoc von `PageSize` trägt keine Speicher-Aussage mehr, nur „ein Lauf hält je Zeitpunkt eine Seite im Speicher“ | **behoben** |
| F-2 (MEDIUM) Auswertung der Plan-Klausel | Messbericht Abschnitt 7.1 gelesen und an den Zeilen nachgerechnet (Abschnitt 4); der Plan (DoD-Zeile „Nachmessung“) trägt die Formulierung samt Grenze (drei Stufen, schmale Zeilen, `n` = 2 plus `n` = 1) | **behoben** |
| F-3 (MEDIUM) Bewertung „neuer Consumer im Lauf“ | Architect-Verdikt `architect-verdict-retention-neue-consumer-und-seitengrenze` liegt vor: akzeptiertes Negativ, keine ADR; die Begründung trägt am Code (`Positions` liest `cdc.consumer_position`; ein Consumer ohne Zeile ist in jedem Lauf ungeschützt). Der Nachzug in Plan §6 und Handbuch steht aus | **bewertet, Nachzug offen (V-1)** |
| F-4 (LOW) Fortschrittsprüfung nur per Zeitüberschreitung rot | Mutation M1 (unten) selbst gefahren: Rot per Assertion in 0,00 s, kein Timeout | **behoben** |
| F-5 (LOW) roter Zwischencommit `d49e248f` | `git archive d49e248f`, `go test` im gepinnten Toolchain-Image: `--- FAIL: TestRunRejectsPageWithoutProgress … gelöschte Menge = [c1 c1], wollen genau die erste Seite [c1]`; Plan §7 benennt die Falle und den Ausweg `git bisect skip d49e248f` | **benannt** (Historie bleibt) |
| F-6 (LOW) Vorher-Nachher-Sprache im Handbuch | Suche der hinzugefügten Handbuch-Zeilen nach früher, bisher, vorher, zuvor, jetzt, nicht mehr, stattdessen: kein Prosa-Treffer (der eine Treffer „Gemessen wurde mit“ ist eine Sachaussage); keine `ADR-`/`LH-`/`SPEC-`/`slice-`-Kennung im Fließtext (nur Historienzeilen und Linkziele) | **behoben** |
| F-7 (INFO) 25.000-Changes-Test bindet die Grenze „Position gleich“ nicht | Plan §3 trägt den Hinweis samt Verweis auf die tragenden Tests | **kenntlich gemacht** |
| F-8 (INFO) Kommentar an der Klassengrenze | `service.go`: „endet als Fehler der Klasse `storage`“ ohne „nicht als Endlosschleife“ | **behoben** |
| F-9 bis F-11 (INFO) | Messbericht Abschnitt 3 „Kennzahl“, Grundlinie-Absatz und Abschnitt 1 Ergebnis 1 gelesen; der Cache-Anteil in `memory.peak` steht auch in `harness/targets/bench-backfill.md` Zeile „Speicher“ | **behoben** |
| F-12 (INFO) Technik-Stratum an der Präzisierungsgrenze | Pflichtenheft Punkt 4 ohne ADR- und Slice-Bezug, Änderungshistorie trägt eine Zeile | zur Kenntnis (Planner) |
| F-13 (INFO) a: Transaktion über zwei Seiten | Architect-Verdikt: akzeptiertes Negativ; Nachzug im Handbuch (Nicht-Atomarität) steht aus | **bewertet, Nachzug offen (V-1)** |

**Mutationen der Eingabeseite selbst gefahren** (Änderung im Arbeitsbaum, `make test` bzw. Store-Tier-Läufe, danach
`git checkout`; Store-Tier über eine Kopie von `tools/harness/run-store-tests.sh` im Scratchpad, gleiche Images, gleicher
d-migrate-Rollout, Aufruf auf `TestRetention…` und `Retention` beschränkt; Baseline-Lauf der Kopie grün):

| # | Mutierte Eingabe | Gesehenes Rot (gedruckt) |
|---|---|---|
| M1 | Fortschrittsprüfung abgeschaltet (`if false && last == after`) | `TestRunRejectsPageWithoutProgress`: „Lauf las über die Seite ohne Fortschritt hinaus: Lese-Aufruf über dem Budget des Fakes (Lese-Aufrufe = 3, erlaubt 2)“ — Assertion in 0,00 s, **kein** `timed out` |
| M2 | `if true \|\| command.Policy.AllowsDeletion(…)` | neun Tests rot, u. a. `TestRunDistinguishesEligibleChangesFromMixedSet`, `TestRunReleasesSameSetAtEveryPageSize` |
| M3 | Consumer-Positionen ignoriert (`consumerPositions[:0]`) | sieben Tests rot, u. a. `TestRunDistinguishesEligibleChangesFromMixedSet` |
| M4 | `DeleteChanges` über die ganze Seite statt der Freigaben | neun Tests rot, u. a. `TestRunDeletesOnlyTheReleasedIDsOfEachPage` |
| M5 | Cursor `after = page[0].ChangeID` | sieben Tests rot, u. a. `TestRunReadsEachPageAtPageSizeFromTheLastKey` |
| M6 | `limit = PageSize+1` | `TestRunReadsEachPageAtPageSizeFromTheLastKey` |
| M7 | Alter mit vertauschtem Vorzeichen | acht Tests rot |
| S1 | Cursor `c.change_id >= $2` | `TestRetentionCandidatePagesCoverExactlyTheSourceChanges`, `TestRetentionCandidateBehindCursorAppearsInNextRun` |
| S2 | Quellfilter ersetzt durch `$1::text IS NOT NULL` (Parameter bleibt referenziert) | `TestRetentionCandidatePagesCoverExactlyTheSourceChanges` |
| S4 | `ORDER BY c.change_id` entfernt | `TestRetentionCandidatePagesCoverExactlyTheSourceChanges`, `TestRetentionCandidateBehindCursorAppearsInNextRun` |
| S8 | `|| true` vor `AllowsDeletion` im Use Case, Lauf über den echten Store | `TestRetentionRunOverManyPagesDeletesWhatAnIndependentCountNames` (Paket `bootstrap`): „Deleted = 25000, unabhängige Zählung nennt 7710“ |

Jede der elf Mutationen färbt mindestens einen Test rot, alle per Assertion. Der Baseline-Lauf am unveränderten Stand ist
grün (`ok` für `postgresstorage` und `bootstrap`). Nicht gefahren: die Store-Tier-Mutationen `LIMIT` entfernt, Position
und Zeitpunkt verschoben, `cdc_admin`-Grant gekürzt (im Messbericht Abschnitt 5 gedruckt; vom Reviewer reproduziert).

## 6. Träger, Zahlen, Sprache, Suchlauf

- **Handbuch:** `Version: 1.63` und die Historienzeile 1.63 stehen (`grep -n '^Version:'`); Zeile 1.62 trägt den
  Slice-Bezug; die Abschnitte „Aufbewahrung (Retention)“, „Bestand als Backfill überführen“ und „Grenzwerte“ tragen die
  Nachmessung samt Ursprung (Lauf-IDs `20260925T183239Z`, `20260925T184503Z`, `20260925T193207Z`), die Bemessungsregel
  „2 KiB je Change plus 64 MiB“ und der `docker --memory`-Hinweis sind entfallen, die Vorversionen stehen im Ist-Zustand.
- **Pflichtenheft:** [`LH-FA-RET-004.a`](../../spec/pflichtenheft.md) Punkt 4 ohne ADR- und Slice-Bezug
  (das Doc-Gate verbietet die Kante Spec → ADR und Slice; `make docs-check` EXIT=0), eine Zeile in der
  Änderungshistorie; Wortlaut „nicht mit nennenswertem Betrag an der Zahl der gespeicherten Changes“.
- **Sensor-Doku, Zahlen mit Ursprung:** `harness/sensors/db-adapter-coverage.md` nennt Nenner **1064** (692 · 32 · 130 ·
  210 = 1064, gerechnet) mit Lauf und gedruckter Zeile — mein `make test-store` druckt `gedeckt 879 von 1064` (**gleich**).
  `harness/sensors/coverage-gate.md` nennt Nenner **2601**, gedeckt 2167, 83,31 %: das Profil `/out/coverage.out` des
  gebauten Images (`pg-change-feed:coverage`, aus meinem `make gates`) per Awk über die Block-Position dedupliziert:
  **2601** Statements (Nenner gleich), gedeckt **2170** = 83,43 %; die gedeckte Zahl ist die eines Laufs und streut um drei
  Statements (die Zeile nennt ihren Lauf über `0e6b1b30`, das Gate druckte dort 83,30 %, bei mir 83,40 %) — Instanz A
  korrekt gekennzeichnet. Die Differenz 34 (2601 − 2567) ist als abgeleitet benannt.
- **`harness/README.md`** (Zeile `make test-store`): nennt die Kandidaten-Seiten-Tests; `harness/targets/bench-backfill.md`:
  der Halbsatz „bis in den GiB-Bereich“ entfällt, die Zeile „Speicher“ nennt den Seiten-Anteil und den Cache-Anteil in
  `memory.peak`.
- **Kommentar an `retentionInterval`** (`internal/bootstrap/wiring.go`): zeilenneutral — `git diff 8cd39719..HEAD --numstat`
  nennt 4 Einfügungen, 4 Löschungen; die Dateilänge ist am Parent und am `HEAD` je 1800 Zeilen; die Zeilen 1183/1184
  (Fehlerzweig der WAL-Rückstands-Messung) und 1287/1288 (`return`) tragen an beiden Ständen denselben Inhalt; die
  Lokatoren `:1287.4,1288.1` und `:1183.5,1184.13` in `coverage-gate.md` stehen unverändert. Zeilen-Lokatoren anderer
  Träger auf `wiring.go` (`ADR-0082`, `ADR-0088`, Beobachtungs-Belege) liegen hinter dem geänderten Kommentar bzw. sind
  Records.
- **Sprache:** keine Chronik/Forensik in den hinzugefügten Handbuch-Zeilen (Suche in Abschnitt 5, F-6); kein host-lokaler
  absoluter Pfad (`make docs-check` EXIT=0).

**Suchlauf-Feld (§3 des Plans), an beiden Ständen nachgefahren** (`git grep -n … <Baum> -- <Wurzeln>`, Zeilen gezählt):

| Plan-Zeile | Meine Messung | Ergebnis |
|---|---|---|
| Befehl 1, Parent `8cd39719` | 298 (`internal` 263, `test` 24, `tools` 2, `spec` 7, `docs/user` 2, `harness` 0) | bestätigt |
| Befehl 2, Parent `8cd39719` | 13 | bestätigt |
| Befehl 3, Parent `8cd39719` | 87 (`docs/user` 25, `harness` 5, `spec` 5, `internal` 16, `tools` 9, `README.md` 0, `docs/plan/adr` 27) | bestätigt |
| `5b1f7762` | 298 / 12 / 80 (Befehl 3: `docs/user` 18) | bestätigt |
| Diff-Stand, Befehl 1 | Arbeitsbaum 309 (`internal` 274, `test` 24, `tools` 2, `spec` 7, `docs/user` 2, `harness` 0) | bestätigt |
| Diff-Stand, Befehl 2 und 3 | Am Stand `b2484add` **6** und **85** (Feld bestätigt); am `HEAD` **7** und **88** (`docs/user` 20, `harness` 8, `spec` 7, `internal` 17) | Feld gilt für `b2484add`, nicht für `HEAD` (V-2) |
| Fixrunde, Befehl A | `c21c43b9` 7, `HEAD` 5 | bestätigt |
| Fixrunde, Befehl B | `c21c43b9` 10, `HEAD` 11 | bestätigt |
| „Nicht gefunden“ (`ADR-0124` Zeile 209, `ADR-0065`, Handbuch „Dauerhaftigkeit“, Pflichtenheft Zeile 673) | alle vier an den genannten Zeilen vorhanden und Gegenstand-fremd bzw. `Accepted` | trägt |
| Lokatoren in `coverage-gate.md` (`git grep -n -E 'wiring\.go:[0-9]+\|:1[0-9]{3}\.[0-9]'`) | keine weiteren Träger außerhalb `docs/reviews` und `done/` mit Lokator hinter dem Kommentar | trägt |

## 7. Release-Aussage

- Server-Tags: `git tag -l 'v*'` nennt `v0.1.0`, `v0.1.1`, `v0.1.2`; `git ls-remote --tags origin` kennt daneben nur Tags des
  Namensraums `sdk-*`; **nichts** über `v0.1.2` ist getaggt.
- Reichweite: `git grep -n ReadChanges v0.1.0 v0.1.1 v0.1.2 -- internal/application/usecase/retention` (ohne Tests) nennt in
  allen drei Tags Zeile 63 `s.store.ReadChanges(ctx, outbound.ChangeQuery{Source: command.Source})`; `retentionInterval`
  steht in `internal/bootstrap/wiring.go` Zeile 172 als `10 * time.Second`; `git ls-tree -r --name-only <Tag> | grep -c
  usecase/backfill` ist 0 in allen drei Tags; `git diff v0.1.2 8cd39719 --stat -- internal/application/usecase/retention
  internal/application/port/outbound/changestore.go` ist leer. Die Aussage „`v0.1.0` bis `v0.1.2` tragen den Defekt“ trägt.
  Der Bedarf in `v0.1.x` steht in der Aussage als hergeleitet (10 Changes/s × 86.400 s = 864.000 Changes, × 1,03 bis 1,59 KiB
  = 0,85 bis 1,31 GiB, gerechnet).
- Vorbedingung: der Slice ist die Vorbedingung des Server-Release `v0.2.0` ([`LH-FA-RET-004`](../../spec/lastenheft.md)); sie
  hat keinen mechanischen Wächter (Plan §4 benennt es). Ob der Slice vor `v0.2.0` in `done/` liegt, hängt an der Closure
  (V-1) und dem Planner.
- `git log`: der rote Zwischencommit `d49e248f` ist in Plan §7 („Bekannte `git bisect`-Falle“) benannt, samt Folgecommit
  `7572051c` und `git bisect skip`.

## 8. Harte Regeln und Grenzen

- **§3.1** — nur `make`, Repo-Skripte und `docker`; `awk`/`sed` (ohne `-i`)/`git`-Auswertungen auf dem Host, kein Host-Go,
  kein Host-Python.
- **§3.2, §3.3, §3.5** — kein `nolint` im Diff; die zwei `git mv`-Commits sind rein; `make doc-immutable` EXIT=0.
- **§3.9** — Exit-Codes nie durch eine Pipe; Gate-Lauf und Folgehandlung getrennt.
- **§3.12** — jede Zahl dieses Reports ist **gemessen** (Befehl oder Lauf genannt) oder als **abgeleitet** gekennzeichnet;
  übernommen sind nur die Findings-Zahlen des Review-Reports (0/3/3/7).
- **§3.13** — Suchlauf-Feld an beiden Ständen nachgefahren (Abschnitt 6).
- **§3.10** — der Slice trägt keinen Workflow; die CI-Läufe am `HEAD` sind trotzdem grün.
- **Nicht gefahren:** `make test-integration` (Retention-Rundlauf: gelesen über CI `e2e`, Abschnitt 1),
  `make test-replication`, `make bench` als Ganzes (Auftrag), Nachmessungen mit anderer Seitengröße, breiten Zeilen,
  Live-Last bei gefülltem `cdc.change`, Läufe über 3.000.000 Changes; die Store-Tier-Mutationen `LIMIT` entfernt, Zeitpunkt/
  Position verschoben, Grant gekürzt (Abschnitt 5).

## 9. Befunde und Restrisiken

Kein Befund blockiert die Verifikation des Codes und der Messung; V-1 (MEDIUM) ist die Bedingung für den Übergang nach
`done/`, V-2 bis V-3 sind Nachzüge, V-4 bis V-5 zur Kenntnis.

- **V-1 (MEDIUM) — Nachzug des Architect-Verdikts steht aus.** Das Verdikt
  `architect-verdict-retention-neue-consumer-und-seitengrenze` benennt drei Träger. (a) Plan §6, Risiken „Ein Consumer
  bestätigt erstmals während eines längeren Laufs“ und „Ein Abbruch zwischen zwei Seiten“: Ausgang „akzeptiertes Negativ,
  Verdikt …“ fehlt (beide tragen „*(bei Closure)*“); (b) im selben Risiko trägt der Plan „nicht breiter geworden (Dauer des
  Laufs nicht länger)“ (Plan-Zeile 510), das Verdikt ersetzt es durch abgeleitete Grenzen (Regelbetrieb 0,38 bis 2,0 s,
  Extremlauf etwa 26 s gegen etwa 3,4 s, das Fenster wächst im Extremlauf um das Siebenfache); (c) Handbuch Zeilen 682/683:
  der Satz „Bestätigt ein Consumer erstmals, während ein Durchlauf läuft, gilt seine Position erst im nächsten Durchlauf“
  trägt die Folge nicht („Changes, die dieser Durchlauf nach dem Lesen der Positionen löscht, sind für diesen Consumer
  verloren“) und der Satz zur Nicht-Atomarität über Transaktionen fehlt; eine Handbuch-Änderung braucht danach eine Version
  1.64 mit Historienzeile. Zuständig: Planner bei der Closure (Plan §6 und §7), Handbuch-Wortlaut über den
  Implementer-Weg. Messbericht Abschnitt 6 Befund 3 („bleibt dem Architect“, „nicht länger“) ist Record und bleibt.
- **V-2 (LOW) — Zählwörter des Suchlauf-Feldes ohne Stand.** Die Zeile „Diff (Arbeitsbaum am Ende der Arbeit …)“ nennt für
  Befehl 2 sechs und für Befehl 3 85 Treffer (`docs/user` 18, `harness` 7, `spec` 7, `internal` 17); das gilt für
  `b2484add`. Am `HEAD` sind es 7 und 88 (`docs/user` 20, `harness` 8), weil die Fixrunde Handbuch- und Vertragszeilen
  ergänzte. Instanz A: eine Zustandsgröße nennt ihren Stand (Commit), nicht „Ende der Arbeit“. Nachzug an den Planner.
- **V-3 (INFO) — Fitness-Function-Zeile 3 der ADR und die Messung.** Die Zeile trägt die Erwartung „die Spitze nach dem
  Run liegt … im Bereich der Spitze im Run (8,9 bis 10,5 MiB) plus dem Bedarf einer Seite“. Gemessen sind 14,9 bis 17,6 MiB
  (`memory.peak`, zwei Läufe), 15,5 bis 17,3 MiB in meinem Lauf; der `anon`-Wert der Spitze im Run liegt bei 10,0 bis 13,2 MiB.
  Wörtlich ist die Erwartung nicht erfüllt, der Form nach ja; der Plan hält die Lesart als Entscheidung des Auftraggebers
  fest (Messbericht Abschnitt 7.1). Die ADR ist `Accepted` und bleibt; für die Closure gehört die Lesart in die
  Closure-Notiz, eine Korrektur wäre eine neue ADR ([`AGENTS.md`](../../AGENTS.md) §3.5).
- **V-4 (INFO) — Cache-Anteil in `memory.peak`.** Der Review-Lauf misst bei 1.000.000 Changes 31,8 MiB mit 16,5 MiB `file`;
  meine drei Runs und die zwei Läufe des Implementers zeigen `file` 0,0 MiB am Run-Ende. Die Kennzahl streut mit einem Anteil,
  dessen Ursache nicht untersucht ist; Handbuch, Messbericht und Vertrag nennen es.
- **V-5 (INFO) — gedeckte Zahl des Unit-Gates.** Nenner 2601 gleich; gedeckt 2167 (Lauf über `0e6b1b30`) gegen 2170 in
  meinem `HEAD`-Lauf: die Zeile nennt ihren Lauf, keine Drift.
- **Restrisiko (weiter offen, im Plan getragen):** die Dauer eines Bereinigungslaufs am Feed ist ungemessen; die Live-
  Erfassung als Speicherquelle ist aus dem Code gelesen, nicht gemessen; breite Zeilen und Läufe über 3.000.000 Changes sind
  ungemessen; ein Host, `n` = 2 bis 3; die Release-Reihenfolge (Slice in `done/` vor `v0.2.0`) hat keinen mechanischen
  Wächter; die Datenbank-Arbeit je Takt wächst weiter linear mit der Zahl der Changes; der Anstieg von 2,4 und 2,5 MiB
  zwischen 1.000.000 und 3.000.000 Changes ist unerklärt.

## Verdikt

**DoD des Plans: erfüllt** — alle 13 Liefer- und Belegzeilen (`[x]`) tragen einen realen Beleg (eigene Läufe Abschnitt 1, 4,
5), die fünf offenen Zeilen sind korrekt offen (Closure durch den Planner). **Entscheidungs-Konformität:** konform mit
[`ADR-0124`](../plan/adr/0124-retention-kandidaten-seitenweise-ohne-row-images.md) Festlegung 1 bis 6,
[`ADR-0014`](../plan/adr/0014-retention-domain-policy.md), [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md),
[`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) und den Port-ADRs;
[`ADR-0007`](../plan/adr/0007-source-ack-outbound-port.md) und
[`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md) sind nicht berührt. **Löschsicherheit:** elf
Eingabeseiten-Mutationen rot per Assertion, `make gates`, `make test`, `make test-store` und CI am `HEAD` grün. **Findings des
Reviews:** F-1, F-2, F-4 bis F-6, F-8 bis F-11 behoben bzw. kenntlich; F-3 und F-13 bewertet (akzeptierte Negative), der
Nachzug steht aus. **Freigabe für `done/` und für `v0.2.0`:** vom Code und von der Messung aus tragfähig; bedingt durch V-1
(Plan §6 Ausgänge, Handbuch-Sätze) und die Closure-Notiz des Planners.

**Übergabe an den Planner:** V-1 (Ausgänge in Plan §6 samt abgeleiteten Grenzen statt „nicht breiter geworden“; Handbuch-Sätze
zu „neuer Consumer“ und Nicht-Atomarität über den Implementer-Weg, Version 1.64), V-2 (Stand des Zählfeldes), Lesart der
Fitness-Function-Zeile (V-3) in die Closure-Notiz, dann Closure. Dieser Report ist ein Lauf-Beleg; er ändert nichts.
