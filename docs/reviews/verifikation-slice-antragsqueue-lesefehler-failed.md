# Verifikations-Report: slice-antragsqueue-lesefehler-failed — 2026-09-26

**Rolle:** Verifier (Modul 11) — Prüfung der Belege, nicht der Behauptungen
([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B). DoD-Abgleich + Entscheidungs-Konformität +
Plan-vs-Code-Diff + Gates. Review-Artefakt des Reviewers:
[`review-slice-antragsqueue-lesefehler-failed.md`](review-slice-antragsqueue-lesefehler-failed.md); Formvorbild dieses
Reports: [`verifikation-slice-harness-fmt-check.md`](verifikation-slice-harness-fmt-check.md). Die vendored Baseline
(`v6.9.0`) trägt kein eigenes Verifikations-Template (`templates/docs/reviews/` enthält nur
`review-report.template.md`, Gerüst per `cp` übernommen und durch die Form des Formvorbilds ersetzt).

**Gegenstand:** Slice-Plan `slice-antragsqueue-lesefehler-failed` (wellenlos; geht `slice-transformationen-start-reihenfolge`
voraus), `HEAD` = `ad9eb61a`, Diff-Range `fc107f42..HEAD`, 8 Commits, 17 Dateien (+979/−137 einschließlich der Lifecycle-Moves,
des Review-Reports und der Träger-Meldung an `slice-transformationen-betriebsdoku`). Inhalt: Lifecycle (`7965076b`, `f7eace20`,
`47a646b5`), Implementer-Lauf (`fd1a7ef2` Code und Spec, `5d8e0748` Plan-Nachzug), Review-Report (`97a0be6f`: 2 HIGH, 1 MEDIUM,
1 LOW, 7 INFO), Fixrunde (`90384748` Code, `ad9eb61a` Plan). Der Stand ist nicht gepusht (`origin/main` = `fc107f42`). Dieser Lauf
ändert weder Code noch Plan noch Doku; er schreibt nur diesen Report. Die Mutationen liefen an der Arbeitskopie mit einer
Einzelersetzung durch ein Skript im Scratchpad der Sitzung (Ersetzung auf genau einen Treffer geprüft, kein `sed -i`,
`perl -pi`, `awk -i`), danach `git checkout -- <Datei>`; `git status --short` nach jeder Mutation und am Ende leer.

## 1. Eigene Sensor-Belege (dieser Lauf, ungepiped — [`AGENTS.md`](../../AGENTS.md) §3.9)

Jeder Lauf schrieb in eine Log-Datei (Scratchpad); der Exit-Code wurde im selben Aufruf gesichert und danach gelesen. Stand
aller Läufe ohne Mutation: `HEAD` = `ad9eb61a`, Arbeitsbaum sauber. Je ein schwerer Docker-Lauf zugleich (`free -m` vorher:
13,6 GB verfügbar, 3,2 GB frei).

| Sensor | Ausgang | Beleg aus meinem Lauf (gedruckt) |
|---|---|---|
| `make test` (Race-Detector) | **Exit 0** | 45 Zeilen `ok`, keine Zeile `FAIL`; letzte Zeile `ok …/tools/schema/rolloutguard 1.016s` |
| `make test-store` | **Exit 0** | `ok …/internal/bootstrap 5.395s`, `ok …/internal/adapters/driven/postgresstorage 13.797s coverage: 78.9% of statements …`, `DB-Adapter-Coverage: 82.56% (gedeckt 885 von 1072 Statements; Profile gemergt: store,replication)`, `db-coverage: OK — … erfuellt Schwelle 80%` |
| `make a-check` | **Exit 0** | `gesamt: 0 Befund(e)` |
| `make coverage-gate` | **Exit 0** | `coverage-gate: OK — Coverage 85.10% erfüllt Schwelle 80%` |
| `make fmt-check` | **Exit 0** | `fmt-check: 255 Go-Dateien geprüft, alle formatiert` |
| `make kommentar-kennungen DIFF=fc107f42 COUNT=1` | **Exit 0** | `0` |
| `make suchlauf-nachmessen PLAN=<Plan des Slice>` | **Exit 0** | `suchlauf-nachmessen: 13 Zeilen stimmen` |
| `make gates` (einmal) | **Exit 0** | `baseline-verify: v6.9.0 OK — 54 Dateien` · `coverage-gate: OK — Coverage 85.10%` · `d-check: 1271 Datei(en) geprüft, 0 Befund(e)` · `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"` · `generated-sync: OK` · `gesamt: 0 Befund(e)` |
| `make commit-traceability RANGE=fc107f42..HEAD` | **Exit 0** | `OK — 8 Commit(s) in "fc107f42..HEAD", Betreffs ohne Struktur-ID` |
| `make doc-commits RANGE=fc107f42..HEAD` | **Exit 0** | `d-check: 1271 Datei(en) geprüft, 0 Befund(e)` |
| `make doc-immutable RANGE=fc107f42..HEAD` | **Exit 0** | `d-check: 1271 Datei(en) geprüft, 0 Befund(e)` (ohne `RANGE` endet das Ziel mit Exit 2 „flag needs an argument: --range“) |
| `make doc-tracked` | **Exit 0** | — |
| Mutationen (§4) | 23 Läufe | 21 rot, 2 grün (beide als äquivalent begründet, §4) |

**Zahlen mit Ursprung ([`AGENTS.md`](../../AGENTS.md) §3.12 Instanz A).** Die Zahlen des Plans (45 Pakete, 255 Dateien,
85.10 %, 13 Zeilen) sind Belege des Fixrunden-Laufs des Implementers; jede stimmt mit meinem Lauf überein (gemessen, s. o.).
Die DB-Adapter-Coverage (82.56 %) und die Paket-Coverage `postgresstorage` (78.9 %) stehen nicht im Plan; sie sind Belege
dieses Laufs, keine Zustandsgrößen.

**Hygiene:** dangling Volumes (`docker volume ls -qf dangling=true | wc -l`) vor dem ersten Lauf **34**, nach allen Läufen
**34**; kein `prune`, kein `system prune`. Nicht gefahren: `make test-integration`, `make test-replication`, `make bench`
(der Diff berührt weder Schema noch Compose noch Runner noch den Erfassungspfad; `git diff --stat fc107f42..HEAD -- tools/schema`
ist leer).

## 2. DoD — Verdikt je Zeile (§2 des Plans)

Gezählt am Plan: acht `[x]`-Zeilen, vier `[ ]`-Zeilen, zusammen zwölf (`grep -c '^- \[x\]'` 8, `grep -c '^- \[ \]'` 4).

| # | DoD-Zeile | Verdikt | Realer Beleg |
|---|---|---|---|
| 1 | Liefer-Punkt 1 — Lesepfad und Verarbeitung (Verwurf durchreichen, `failed` mit Text, Zeilen dahinter laufen, Zeile ohne Kennung mit Warnung, `ListPending` nur bei Fehler der Lesung selbst) | **erfüllt** | Store-Test `TestAdministrationRequestListPendingPassesRejectedRowsThrough` (reale Zeilen über die SQL-Funktionen, vier Gründe, je eine gültige Zeile davor und dahinter, Vermerke, Bereinigung nur nach Kennungen des Tests) grün im `make test-store`-Lauf und rot unter der Mutation S-A/S-B (§4); Whitebox-Tests der Verarbeitung (`…FailsRejectedRowsInQueueOrder`, `…RejectedRowSurvivesMarkFailedError`) grün und rot unter W1–W6; Login-Test `TestAdministrationPathRunsUnderLeastPrivilegeLogins` führt vier verworfene Zeilen und die gültige dahinter unter `cdc_admin` bis `failed`/`applied`, rot unter S-A (Status `pending`); `TestReadPendingRequestsRejectsRowWithEmptySchemaOrTable` ist ersetzt durch `TestReadPendingRequestsPassesRejectedRowsThrough` (`git grep -n` in `*.go`: kein Treffer des alten Namens; er steht nur noch in Plan und Record); `TestReadPendingRequestsClassifiesScanFailure` und die Iterations-Tests bleiben Fehler (E3 rot) |
| 2 | Liefer-Punkt 2 — die Spec im selben Commit | **erfüllt** | `fd1a7ef2` trägt Code und `spec/pflichtenheft.md`; die sechs Klartexte stehen Zeichen für Zeichen gleich in Spec (Backtick-Zellen) und Code (`rejectionMessage`) — von Hand je Text mit `grep -c` in beiden Dateien 1/1; die hinzugefügten Spec-Zeilen tragen weder `ADR-`, `slice`, `welle` noch `ARC-` (`git diff -U0` + `grep` Exit 1); Zeilen `schema_name`/`table_name`, `column_name`, Absatz „Transformations-Antragsarten“ und die Änderungshistorie sind mitgezogen (Diff gelesen); `make docs-check` (Teil von `make gates`) Exit 0 |
| 3 | Liefer-Punkt 3 — Kommentare und Träger des alten Verhaltens | **erfüllt** | Suchlauf-Feld an beiden Ständen von Hand nachgefahren (§3); die 7 Fundstellen der Ablehnung beim Lesen sind nachgezogen, die eine verbleibende ist der Record der Änderungshistorie |
| 4 | `make gates` grün, Exit gesondert | **erfüllt** | eigener Lauf Exit 0 (§1) |
| 5 | Review durchgeführt, Report liegt vor | **erfüllt** | `review-slice-antragsqueue-lesefehler-failed` liegt vor (2 HIGH, 1 MEDIUM, 1 LOW, 7 INFO); Finding für Finding in §7 am Ist-Zustand nachgemessen — F-1 bis F-4 aufgelöst, F-5 bis F-11 als Kenntnis geführt; kein offenes HIGH oder MEDIUM |
| 6 | §3.13-Suchlauf: Feld mit Gefundenem und Nichtgefundenem je Träger, beide Stände, `make suchlauf-nachmessen` nach jeder Fixrunde | **erfüllt** | 13 Zeilen mit dem Werkzeug Exit 0; neun Zeilen von Hand nachgefahren (§3); Prosa und Träger-Tabelle nennen Gefundenes und Nichtgefundenes je Träger |
| 7 | Doku-Update: `SPEC-019`; Handbuch ohne Aussage zur Ablehnung beim Lesen, Versionshistorie unberührt | **erfüllt** | `git diff --name-only fc107f42..HEAD -- docs/user harness` leer; die Suchlauf-Zeilen 11 und 12 messen 0 an beiden Ständen; `git grep` der Log-Texte der Verarbeitung in `docs/user harness spec` ohne Treffer |
| 8 | Reconciliation-Register — entfällt | **erfüllt (entfällt)** | Greenfield, keine Datei |
| 9 | Closure-Notiz mit Steering-Loop-Lerneintrag | **korrekt offen** | Plan §7 trägt „*(zu tragen bei Closure)*“ (Planner) |
| 10 | Beobachtungs-Register fortgeschrieben | **korrekt offen** | Closure-Pflicht (Planner); der Diff berührt keine Register-Datei; `state.md` von `antrag-mit-leerem-regelnamen-stallt-die-queue` steht auf „entschieden, Fix offen“ |
| 11 | Jedes Risiko aus §6 trägt einen Ausgang | **korrekt offen** | Plan §6: acht Zeilen mit „Ausgang: *(bei Closure)*“ (Planner) |
| 12 | Die drei Paarungen | **korrekt offen** | hängen an der Closure des Slice (Planner) |

Kein `[x]` ohne Beleg, kein `[ ]`, das über die Rollen-Sequenz hinaus belegt wäre. Closure-Notiz, Register, Risiko-Ausgänge und
Paarungen sind Planner-Arbeit und nicht Teil dieser Prüfung. §5 (Closure-Trigger) verlangt einen realen, grünen
`make test-store`-Lauf mit gedruckter Zeile: die Zeile steht in §1 (`ok …/postgresstorage 13.797s`, `ok …/internal/bootstrap 5.395s`).

## 3. Plan-vs-Code-Diff und Zahlen mit Ursprung ([`AGENTS.md`](../../AGENTS.md) §3.12)

**§1 Ziel und „Ausdrücklich NICHT“.** Eingehalten: `git diff --stat fc107f42..HEAD -- tools/schema` leer (SQL-Funktionen, Schema und
Grants unverändert; kein `make schema-rollout` nötig); die Prüfung der Regelfelder und K1 bis K4 im Use Case bleibt
(`git diff --stat` nennt keinen Pfad unter `internal/application/usecase/`); kein Recovery-Weg für Schema-Fehler; die Zeile ohne
Kennung wird als benannte Grenze übersprungen.

**§3-Tabelle, Zeile für Zeile gegen den Diff:**

| Plan-Zeile | Ist im Diff |
|---|---|
| Port `administrationrequest.go` | `PendingAdministrationRequest` (`Request` oder `Rejected`) und `RejectedAdministrationRequest` (`ID`, `Message`); der Port importiert nur `model` (Ports-Regel, `a-check` 0 Befunde) |
| `sqlexec/translate.go` | Konstruktor-Fehler wird als `Rejected` an der Stelle der Ordnung angehängt (`continue`), `rows.Err()`/Scan/Anfrage bleiben Fehler; `rejectionMessage` bildet Klartext, Doppelpunkt, Leerzeichen, Kennung |
| `postgresstorage/administrationrequest.go` | `ListPending` folgt der Form (4 Zeilen) |
| `bootstrap/wiring.go` | `failRejectedAdministrationRequest`; `continue` an der Stelle der Queue; Zeile ohne Kennung: Warnung, kein Vermerk; Fehler von `MarkFailed` protokolliert |
| `model/administrationrequest.go` | Konstruktor unverändert, Kommentar im Indikativ; die Kennung `ARC-001` im Kommentar entfällt (eine Kennung je Kommentarblock, `make kommentar-kennungen` 0) |
| `spec/pflichtenheft.md` `SPEC-019` | Absatz „Zeilen, die kein Antrag sind“ mit Tabelle der sechs Klartexte, Grenze, Änderungshistorie |
| Tests: `translate_test.go`, `rejection_internal_test.go` (neu), `administrationrequest_test.go`, `administrationrequest_order_test.go`, `administration_internal_test.go`, `administration_endtoend_test.go`, `administration_roles_internal_test.go`, `model/administrationrequest_test.go` | alle vorhanden; Fakes und Aufrufer folgen der Form (`make test`, `make test-store` übersetzen sie) |
| `harness/README.md`, Handbuch | nicht berührt (`git diff --name-only` leer für `harness`, `docs/user`) |
| `open/slice-transformationen-betriebsdoku.md` §2 | Übergabe-Text „Zeilen, die kein Antrag sind“ vorhanden (`git grep -c -i 'Zeilen, die kein Antrag sind'` = 2); die Probe des Reviews (`git grep -c -i -E 'verworfen\|Konstruktor\|Lesung\|kein Antrag\|leeres Schema\|leere Spalte'`) trifft jetzt 2 Zeilen statt 0; der Text nennt Verhalten, Klartext und Kennung in `error_message`, die Grenze „Zeile ohne Kennung“, die Belege (Store-Test und Login-Test) und die Gründe ohne Weg über die SQL-Funktionen — tragfähig |

Über den Plan hinaus, im Plan als **Nachzug** und **Fixrunde** benannt: Form der Durchreichung (Entscheidung am Start),
`rejection_internal_test.go`, die Zeile im Plan der Adresse, die Reihenfolge-Paare. Keine `Accepted` ADR geändert
(`make doc-immutable` Exit 0), keine Schwelle und keine Gate-Konfiguration verändert ([`AGENTS.md`](../../AGENTS.md) §3.5, §3.6).
Die zwei Lifecycle-Commits (`7965076b` open → next, `47a646b5` next → in-progress) sind reine Renames
([`AGENTS.md`](../../AGENTS.md) §3.3).

**Suchlauf-Feld von Hand nachgefahren** (`git grep -n -E …` ohne das Werkzeug; Ist gleich Soll in allen neun Zeilen):

| Zeile des Blocks | Stand | Soll | Ist von Hand |
|---|---|---|---|
| 1 `ReadPendingRequests\|ListPending`, ohne Tests | `47a646b5` | 8 | 8 |
| 2 dasselbe, mit Tests | `47a646b5` | 73 | 73 |
| 3 `NewAdministrationRequest`, ohne Tests | `47a646b5` | 7 | 7 |
| 4 Beschreibung der Ablehnung beim Lesen | `47a646b5` | 8 | 8 (Test-Godocs 757, 1262, `translate_test.go` 1317, `administration_internal_test.go` 940, Konstruktor 85–86, Spec 676 und 1176) |
| 5 wie 1 | `diff` | 11 | 11 (Adapter 4, `sqlexec` 3, Port 2, Verarbeitung 2) |
| 6 wie 2 | `diff` | 91 | 91 |
| 7 wie 3 | `diff` | 7 | 7 |
| 8 wie 4 | `diff` | 1 | 1 (Änderungshistorie der Spec, Record) |
| 9/10 weiter gefasste Beschreibung | `47a646b5` / `diff` | 7 / 6 | 7 / 6 (die sechs: zwei Testmeldungen „keine Lesefehler“, ein Log-Assert, zwei Tests des Lesefehlers, die Änderungshistorie) |
| 13 Zeilen-Lokatoren `datei.go:zahl` in Plänen | `diff` | 0 | 0 |

Alle „Zeile N“-Verweise im Plan gegen den Block gelesen: „Zeilen 1 und 2“ (Umfang, Risiko 1), „Zeile 4“ (Risiko 5, Träger-Tabelle,
Handbuch-Zeile), „Zeilen 5 bis 7“, „Zeile 8“ (README: nicht berührt), „Zeilen 9 und 10“, „Zeilen 11 und 12“ (Handbuch), „Zeile 13“
(Lokatoren) — alle treffen die Zeile, die sie nennen. Das Feld schließt die Plan-Datei aus, ihr Suchraum ist `*.go`, `spec`,
`docs/user`, `harness` (ADRs und Records sind ausgenommen, im Plan benannt).

**Zahlen mit Ursprung.**

| Zahl | Stand und Ursprung im Plan | Mein Lauf |
|---|---|---|
| 8/73 und 11/91 Fundstellen | Plan: gemessen, `47a646b5` bzw. `diff`, Feld in §3 | von Hand gleich (s. o.) |
| 45 Pakete `ok`, 255 Dateien, 85.10 %, 13 Zeilen | Plan: Belege des Fixrunden-Laufs | mein Lauf gleich (§1) |
| Umfang „8 Nicht-Test- und bis zu 73 Fundstellen“ | Plan: gemessen, Schätzung des Umfangs als solche gekennzeichnet | stimmt |
| „sieben Fälle der Tabelle und beide Tests“ rot (Zeile „Antragsart hängt am Grund“) | Plan: Lauf der Fixrunde | E8 (§4): sieben Fälle und beide Tests rot, gleich |

Der Parent des Fix-Commits ist `47a646b5`; der Plan nennt ihn durchgängig (Feld-Stand, Umfang, §6).

## 4. Mutationen der Eingabeseite (dieser Lauf)

Eine Mutation je Lauf an der Arbeitskopie, `go test -race` im gepinnten Toolchain-Image (`--network none`) für die Pakete
`sqlexec` bzw. `bootstrap`; die DB-Läufe über `make test-store` (bei den Läufen S-A und S-B waren im Runner die Aufrufe der
beiden vorgezogenen Pakete mit `|| true` versehen, damit der Lauf `postgresstorage` erreicht; der Runner selbst wurde danach
per `git checkout` zurückgesetzt). Jeder Lauf: Exit ungepiped, danach `git checkout -- <Datei>`.

**Eingabeseite der Lesung (`sqlexec/translate.go`, Paket `sqlexec`):**

| # | Zusage | Mutation | Rote Fälle (gedruckt) |
|---|---|---|---|
| E1 | die Lesung lehnt keine Zeile ab | Konstruktor-Fehler wieder zurückgeben (`return nil, err` statt anhängen und `continue`) | `…PassesRejectedRowsThrough`, `…KeepsReadFailureAfterRejectedRow` |
| E2 | Text je Grund | `Tabellenname ist leer` → `Tabelle ist leer` | `…PassesRejectedRowsThrough` |
| E3 | Fehler der Iteration bleibt Fehler | Prüfung von `rows.Err()` wirkungslos (`err != nil && false`) | `…ClassifiesIterationFailure`, `…KeepsReadFailureAfterRejectedRow` |
| E4 | Schema vor Tabelle | Fälle `schema`/`table` vertauscht | `…/leeres_Schema_und_leerer_Tabellenname` |
| E5 | Quelle vor Schema | Fälle `source`/`schema` vertauscht | `…/leere_Quelle_und_leeres_Schema` |
| E6 | Kennung vor Quelle | Fall `id` hinter `source` gestellt | `…/leere_Kennung_und_leere_Quelle` |
| E7 | Zeile ohne Kennung trägt die Kennungs-Adresse | Fall `id` entfernt | `…/leere_Kennung_und_leere_Quelle`, `…/Zeile_ohne_Kennung` |
| E8 | Antragsart hängt am Grund des Konstruktors | Fall der Antragsart vor `schema`, Bedingung `cause != nil` | sieben Fälle (leeres Schema, leerer Tabellenname, beide Paare mit Schema/Tabelle und Antragsart, Schema+Tabelle, `exclude_column`, `include_column`) und beide Tests in `rejection_internal_test.go` (vier Teilfälle) |
| E9 | (Stellung des Antragsart-Falls) | Fall der Antragsart nur vor `schema` gestellt, Bedingung `Is(cause, ErrInvalidAdministrationRequestKind)` unverändert; ebenso mit vorangestelltem `cause != nil &&` | **grün — äquivalent**, beide Läufe (s. u.) |
| E10 | Spaltenname nur bei leerer Spalte | `column == "" &&` entfernt | `…ColumnCaseNeedsEmptyColumnAndEmptyIdentifierCause/gesetzte_Spalte,_Grund_leerer_Bezeichner` |
| E11 | Spaltenname nur mit Grund `ErrEmptyIdentifier` | Grundprüfung entfernt (`column == ""` allein) | `…/leere_Spalte,_fremder_Grund` |
| E12 | allgemeiner Klartext | Anfangswert `clear` auf die leere Zeichenkette | `…FallsBackToGeneralText` und zwei Teilfälle des Spalten-Tests |
| E13 | Ordnung der Queue | verworfene Zeilen in eine zweite Liste, ans Ende gehängt | alle zwölf Fälle von `…PassesRejectedRowsThrough` |

**Verarbeitung (`bootstrap/wiring.go`, Paket `bootstrap`, Whitebox):**

| # | Zusage | Mutation | Rote Fälle (gedruckt) |
|---|---|---|---|
| W1 | die Zeilen dahinter laufen | `continue` → `return` nach der verworfenen Zeile | `…FailsRejectedRowsInQueueOrder`, `…RejectedRowSurvivesMarkFailedError` |
| W2 | Text des Vermerks ist der der Lesung | fester Text `verworfen` statt `rejected.Message` | `…FailsRejectedRowsInQueueOrder` |
| W3 | Zeile ohne Kennung wird nicht vermerkt | Guard `rejected.ID == ""` entfernt | `…FailsRejectedRowsInQueueOrder` |
| W4 | ein Fehler von `MarkFailed` bricht nichts ab | Warnung durch `panic(markErr)` ersetzt | `…RejectedRowSurvivesMarkFailedError` |
| W5 | Zeile ohne Kennung trägt eine Warnung | Warnung entfernt | `…FailsRejectedRowsInQueueOrder` |
| W6 | Vermerk an der Stelle der Queue | verworfene Zeile per `defer` am Ende vermerkt | `…FailsRejectedRowsInQueueOrder` |

**Gegen die reale PostgreSQL (`make test-store`):**

| # | Zusage | Mutation | Rote Fälle (gedruckt) |
|---|---|---|---|
| S-A | die Lesung lehnt keine Zeile ab | wie E1 | Login-Test `TestAdministrationPathRunsUnderLeastPrivilegeLogins` (`verworfene Zeile unter cdc_admin-Login: Status "pending", … erwartet "failed" …`); Store-Test `…ListPendingPassesRejectedRowsThrough` (`ListPending = leere Kennung, wollen nil`), dazu die `sqlexec`-Fälle aus E1 |
| S-B | Ordnung der Queue (`requested_at`) | `ORDER BY administration_request_id` ohne `requested_at` in `SelectPendingAdministrationRequests` | `TestAdministrationRequestListPendingPassesRejectedRowsThrough` (`die verworfene Zeile … steht nicht hinter ihrer Vorgängerin`), `TestAdministrationRequestSameTransactionCallsKeepCallOrder`, im Paket `bootstrap` `TestAdministrationSameTransactionRemoveThenSetLeavesTheNewRuleLiveAndDerived` |

Damit sind die Mutationen der Eingabeseite der DoD (Konstruktor-Fehler zurückgeben; Schleife bricht ab; fester Text; Fehler der
Lesung bleibt Fehler) sowie Prüfreihenfolge je Paar, `rows.Err()`-Zweig, Guard `ID == ""`, Vermerk-Reihenfolge, Text des Vermerks
und Ordnung im Adapter je rot gesehen; E1 (das Lese-Ergebnis) und W1 (die Verarbeitung) färben jeweils Tests **beider** Schichten rot.

**Menge der Stellen (§3.12).** Die Mutationsangaben des Implementers nennen die Stelle je Zusage: die Tabelle der Fixrunde nennt
den Schalter von `rejectionMessage` (sieben Zeilen); die Godocs nennen `sqlexec.ReadPendingRequests` (Rückgabe des
Konstruktor-Fehlers, eine Stelle: `git grep -n 'return nil, err'` im Block der Lesung trifft sie), `processAdministrationRequests`
(`continue`) und den Text des Vermerks. Die Angaben stimmen: jede genannte Mutation ist rot (E1–E8, E10–E13, W1–W6), die eine
als äquivalent genannte (E9) ist grün (in beiden Varianten); über die genannten Stellen hinaus gibt es an diesen Zusagen keine weitere
(`rejected.Message` ist der einzige Text-Weg in den Vermerk, `Rejected` hat eine Erzeugerstelle).

**Bewertung der Erkenntnis „die Stellung des Antragsart-Falls im Schalter ist nicht beobachtbar“.** Wahr, und die Angabe im
Godoc und im Plan trägt. Der Konstruktor prüft Kennung, Quelle, Schema und Tabelle in **einem** Ausdruck und liefert dafür
`ErrEmptyIdentifier`, bevor er die Antragsart betrachtet; der Grund `ErrInvalidAdministrationRequestKind` entsteht daher nie
neben einer leeren Kennung, Quelle, Schema oder Tabelle, und der Fall `column == "" && Is(cause, ErrEmptyIdentifier)` schließt
ihn ebenso aus (unbekannte Art heißt Grund `ErrInvalidAdministrationRequestKind`, nicht `ErrEmptyIdentifier`). Eine Umstellung
des Antragsart-Falls ändert damit den Text keiner Zeile: E9 grün ist ein echtes Äquivalent, E8 (Bedingung auf den Grund gelöst)
färbt die Zeilen rot, bei denen die Bedingung trägt. Die Tabelle des Plans („Antragsart hängt am Grund …, nicht an der Stellung
im Schalter“) ist wahr; die Zusage „Reihenfolge über die Paare gebunden“ gilt für die beobachtbaren Paare (Kennung/Quelle,
Quelle/Schema, Schema/Tabelle — E4–E6 rot).

## 5. Entscheidungs-Konformität und Träger

**Architect-Entscheidung, Option (a) in der allgemeinen Form** (Register `BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue`,
`state.md`): eingehalten — die Lesung lehnt keine einzelne Zeile ab, jede vom Konstruktor verworfene Zeile endet `failed`, die
Zeilen dahinter werden verarbeitet; der Fix ist Code (Lesepfad und Verarbeitung) ohne neue ADR. Eine Nuance: `state.md` nennt „den
Fehlertext des Konstruktors“; der Konstruktor liefert Sentinel-Fehler ohne Text, der Klartext entsteht in `rejectionMessage` aus
Grund und Feldern der Zeile (Plan §3 „Form der Durchreichung“ hält das fest; die Spec trägt Klartext und Adresse je Grund). Ein
künftiger Konstruktor-Grund ohne eigenen Fall trägt bis dahin den allgemeinen Klartext (`Antrag ist ungültig`) — benannte Grenze
(Plan §6, Review F-5), keine Abweichung vom Zug.

**[`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md).** Antrags-Queue, Status `pending`/`applied`/`failed`,
Abfrage und Wecksignal unverändert (`SelectPendingAdministrationRequests` im Wortlaut gleich, `LISTEN`-Schleife nicht berührt);
der `failed`-Vermerk nutzt die bestehende Port-Fähigkeit `MarkFailed`.

**[`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md).** Keine Domänenlogik in SQL: `git diff --stat fc107f42..HEAD -- tools/schema`
leer; die Prüfung liegt im Konstruktor und in `sqlexec`/Verarbeitung.

**[`ADR-0029`](../plan/adr/0029-domain-invarianten.md).** Der Konstruktor `NewAdministrationRequest` ist unverändert (nur sein Kommentar); die
Invarianten bleiben im Domain-Core.

**[`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md) / [`ADR-0028`](../plan/adr/0028-inbound-use-cases.md).** Der Port trägt Domänentypen und Text, keinen Treibertyp;
`make a-check` 0 Befunde.

**[`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md).** Die Prüfung der Regelfelder liegt weiter in der Verarbeitung
(`git diff` nennt keine Datei unter `internal/application/usecase/`); `TestReadPendingRequestsCarriesRuleRowsWithEmptyRuleFields`
und `TestAdministrationRequestListPendingCarriesRuleRowsWithMissingFields` bleiben grün.

**[`ADR-0127`](../plan/adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md).** Ordnung `requested_at`, dann Kennung, unverändert; die Verarbeitung vermerkt die verworfene
Zeile an ihrer Stelle (S-B und E13 rot).

**`SPEC-019` gegen den Code.** Sechs Klartexte Zeichen für Zeichen gleich (§2 Zeile 2). Ort der Prüfung, Reihenfolge („die erste
verletzte Prüfung bestimmt den Text, in der Reihenfolge der Tabelle“ — belegt E4–E6, E10, E11), Grenze „Zeile ohne Kennung“
(überspringen, Warnung, `pending`; W3, W5) und „Fehler beim Vermerk“ (protokollieren, weiter; W4) entsprechen dem Code. Die Spec
trägt keinen ADR- und keinen Slice-Bezug ([`AGENTS.md`](../../AGENTS.md) §3.4).

**[`AGENTS.md`](../../AGENTS.md) §3.1 (Host-Werkzeuge).** Im Diff kein `sed -i`, `perl -pi`, `awk -i` in einer Kommandozeile; die
Mutationen dieses Laufs nutzten ein Einzelersetzungs-Skript auf Kopien der Arbeitsdateien, keine Datei blieb geändert.

**Handbuch-Kandidatenlauf** (`git diff --name-only 47a646b5 -- internal/bootstrap/ tools/schema/ internal/adapters/driving/`): vier
Dateien, alle in `internal/bootstrap/` — drei Test-Dateien und `wiring.go`; `tools/schema/` und `internal/adapters/driving/` leer. Das
Diff von `wiring.go` fügt eine Funktion (`failRejectedAdministrationRequest`) und die Verzweigung in der Schleife hinzu; neu für den
Betreiber sind zwei Log-Meldungen (`administration: Antrag verworfen`, `administration: Zeile ohne Kennung übersprungen`) und das
Verhalten „verworfene Zeile endet `failed` mit Text“. Weder Handbuch, `harness` noch `spec` nennen diese Log-Texte oder eine
Ablehnung beim Lesen (`git grep` ohne Treffer); keine neue Konfiguration, kein neuer Endpunkt, kein CLI-Sondermodus. Der
Handbuch-Aufschub liegt bei `slice-transformationen-betriebsdoku`, dessen §2 den Gegenstand jetzt als Text trägt (§3).

## 6. Hard Rules

Docker-only (alle Werkzeugläufe über `make` bzw. das gepinnte Image; Host: `git`, `python3` nur für die Einzelersetzung der Mutationen im
Scratchpad); kein `//nolint`; kein Gate ohne ADR (keine Gate-Änderung); Exit-Codes ungepiped gesichert (§1); jeder Betreff mit `ADR-0050`
und ohne `SPEC-`/`ARC-` (`make commit-traceability` Exit 0, 8 Commits); Code und Spec im selben Commit (`fd1a7ef2`); der Plan bleibt in
`in-progress/` (die Closure ist Planner-Arbeit); keine Chronik-Sprache in den neuen Doku-Absätzen der Spec.

## 7. Findings des Reviews — Auflösung am Ist-Zustand

Die Zählung folgt dem Review-Report (F-1 bis F-11).

| Finding | Kategorie | Stand am Ist-Zustand | Beleg aus meinem Lauf |
|---|---|---|---|
| F-1 | HIGH | **aufgelöst** | `git grep -n 'hielte' -- internal/adapters/driven/postgresstorage/sqlexec/translate.go` Exit 1 (kein Treffer); `git diff fc107f42..HEAD -U0 -- '*.go' \| grep -nE '^\+.*//.*(wäre\|würde\|hielte\|hätte\|sonst\|statt)'` trifft acht Zeilen: sieben Mutationsbeschreibungen in Test-Godocs (Subjekt der Test, „statt“ als Gegenüber zweier Zustände) und ein Indikativ-Satz in `rejectionMessage` („trägt statt der Kennung `schema.table` als Adresse“); der Kommentar von `ReadPendingRequests` nennt die Zusage („steht als `Rejected` … an ihrer Stelle der Ordnung; die Lesung endet nur bei einem Fehler …“) |
| F-2 | HIGH | **aufgelöst** | Handbuch-Zeile der Träger-Tabelle nennt „Suchlauf Zeilen 11 und 12“; im Block sind das die beiden Zeilen über `docs/user harness` (Soll/Ist 0); die README-Zeile nennt Zeile 8 (harness ohne Treffer) |
| F-3 | MEDIUM | **aufgelöst** | Übergabe-Text in `open/slice-transformationen-betriebsdoku.md` §2 vorhanden und tragfähig (§3); `git grep -c -i 'Zeilen, die kein Antrag sind'` = 2 (Plan) |
| F-4 | LOW | **aufgelöst** | Paare Kennung/Quelle, Quelle/Schema, Schema/Tabelle als Fälle der Tabelle in `TestReadPendingRequestsPassesRejectedRowsThrough` (E4–E6 rot); Antragsart-Stellung als äquivalent begründet und bestätigt (E9, §4); Spalten-Bedingungen in `rejection_internal_test.go` (E10, E11 rot) |
| F-5 | INFO | **als Kenntnis mit Adresse geführt** | Plan §6 letzter Punkt („Grenzen der Durchreichung“) und Behandlungsabsatz der Fixrunde nennen die Kopplung von `rejectionMessage` an die Reihenfolge und den allgemeinen Klartext; der Kommentar der Funktion nennt sie |
| F-6 | INFO | **als Kenntnis mit Adresse geführt** | Plan §6 (Invariante `Request`/`Rejected` gilt per Konvention einer Erzeugerstelle) |
| F-7 | INFO | **als Kenntnis geführt** (im Plan nicht namentlich) | Grenzen der Verarbeitung (Warnung je Durchlauf, Vermerk wiederholt sich) stehen in der Spec-Grenze und im Godoc von `failRejectedAdministrationRequest`; keine Aktion erwartet |
| F-8 | INFO | **aufgelöst** | `git grep -n 'würde jeden Durchlauf' -- internal/bootstrap/wiring.go` Exit 1; der Block über `processAdministrationRequests` sagt „die Abfrage liest nur `pending`, der Antrag wird nicht erneut versucht“ |
| F-9 | INFO | **aufgelöst** | Parent im Feld: `47a646b5` in allen sieben Stand-Zeilen; von Hand nachgefahren (§3). Der Plan bezeichnet diesen Befund als „F-7“ (V-1) |
| F-10 | INFO | **Kenntnis** | `make kommentar-kennungen DIFF=fc107f42 COUNT=1` Exit 0, Ausgabe 0; die Kennungen stehen anderswo im Baum weiter |
| F-11 | INFO | **Kenntnis** | im Diff keine Spur eines Host-Interpreters; `make fmt-check` Exit 0 |

## 8. Verifier-Findings

| ID | Kategorie | Befund | Verifizierbar |
|---|---|---|---|
| V-1 | LOW | Der Fixrunden-Absatz des Plans zählt die Findings des Reviews falsch: „F-7 (Stand ‚Parent‘ ist `47a646b5` …)“ meint im Review F-9 (F-7 sind die benannten Grenzen der Verarbeitung), und „F-9, F-10, F-11 — keine Aktion“ verdeckt, dass F-9 einen Nachzug trug. Die Behebung selbst ist im Plan vollständig; nur die Zuordnung Finding ↔ Behandlung ist verschoben. Kein Einfluss auf eine DoD-Zeile, kein Sensor meldet es. | Plan §3, Absatz „Fixrunde zum Review“, gegen `docs/reviews/review-slice-antragsqueue-lesefehler-failed.md` Zeilen 180 und 182 |
| V-2 | INFO | Der Godoc von `rejectionMessage` nennt „die Reihenfolge der Prüfungen des Konstruktors: Quelle, Schema, Tabelle, Antragsart, Spalte“. Der Konstruktor prüft Kennung, Quelle, Schema und Tabelle in einem Ausdruck mit einem Fehlerwert; die Reihenfolge dieser vier ist eine Konvention von `rejectionMessage` (durch Spec-Tabelle und die Paar-Fälle gebunden), nicht eine des Konstruktors, und die Stellung des Antragsart-Falls ist nicht beobachtbar (E9). Wahre Aussage über das Verhalten, ungenaue Herkunftsangabe im Kommentar. Kenntnis; ohne erwartete Aktion. | `internal/domain/model/administrationrequest.go` (`if id == "" \|\| source == "" \|\| schema == "" \|\| table == ""`) gegen den Godoc in `translate.go` |
| V-3 | INFO | Die Spec-Tabelle führt `Quelle ist leer` und `Antragsart ist unbekannt` als Gründe; über die SQL-Funktionen (Fremdschlüssel auf `cdc.source`, `CHECK` `chk_administration_request_kind`) entstehen sie nicht, ihr Text ist nur am Fake-Test belegt (Plan §6 nennt es; die Spec nennt es nicht, und ihr Satz „die SQL-Funktionen prüfen Quelle … nicht“ lässt den Fremdschlüssel weg). Kein Widerspruch zum Code; Kenntnis. | `tools/schema/schema.yaml` (`references`), `nacharbeit-administration.sql` Zeile 83; Plan §6 |
| V-4 | INFO | Eine verworfene Zeile der Antragsart `backfill` einer fremden Quelle wird von dieser Instanz `failed` vermerkt, während eine gültige `backfill`-Zeile der fremden Quelle `pending` bleibt (die Quellbindung der Schleife gilt nur für gültige Anträge; alle übrigen Antragsarten werden ohnehin quellenunabhängig verarbeitet). Für eine ungültige Zeile folgenlos (sie wäre für jede Instanz ungültig); Kenntnis. | `internal/bootstrap/wiring.go` Zeilen 1354–1361 |

Kein HIGH, kein MEDIUM aus Verifier-Sicht. V-1 ist ein Nachzug des Plans (Planner bei der Closure); V-2 bis V-4 ohne erwartete Aktion.

## 9. Verdikt

**DoD trägt:** ja. Acht `[x]`-Zeilen sind am Ist-Zustand belegt (die Läufe sind meine, nicht die des Implementers oder des Reviewers);
vier `[ ]`-Zeilen sind korrekt offen (Closure, Register, Risiko-Ausgänge, Paarungen — Planner). Aus Verifier-Sicht ist der Slice
closure-fähig; V-1 (Zuordnung Finding ↔ Behandlung im Plan) ist beim Closure-Nachzug mitzunehmen.

**Zusammenfassung der Gates:** `make test` (45 Pakete), `make test-store` (`postgresstorage` 13.797s, `bootstrap` 5.395s, DB-Adapter-Coverage
82.56 %), `make a-check`, `make coverage-gate` (85.10 %), `make fmt-check` (255 Dateien), `make kommentar-kennungen DIFF=fc107f42 COUNT=1`
(0), `make suchlauf-nachmessen` (13 Zeilen), `make gates`, `make commit-traceability`, `make doc-commits`, `make doc-immutable RANGE=…`,
`make doc-tracked` — alle Exit 0; 23 Mutationsläufe (E1–E13, W1–W6, S-A in zwei Läufen, S-B; E9 in zwei Varianten): 21 rot, 2 grün und als
äquivalent begründet (beide Läufe von E9: Stellung des Antragsart-Falls bei unverändertem Grund); Baum am Ende sauber.

**Übergabe an den Planner (Closure):**
1. V-1 im Fixrunden-Absatz des Plans nachziehen (Finding-Nummern des Reviews).
2. Träger-Meldungen mit Frist „Closure dieses Slice“: `state.md` von `BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue` („Fix
   geliefert“), Welle-Datei `welle-transformationen.md` (Träger benennen), Register `BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an`
   (fünfte Ausprägung, in F-3 des Reviews genannt, hier aufgelöst).
3. Risiko-Ausgänge §6 (acht Punkte; Belege: Suchlauf-Feld, Whitebox-Tests, Store-Test mit skopierter Bereinigung, Login-Test, Diff ohne
   `tools/schema`), Closure-Notiz mit Steering-Loop-Eintrag (Spec-Lücke von `SPEC-019` geschlossen), drei Paarungen.
4. Finding-Klassen des Reviews (Konjunktiv-Kommentar, Beleg-Verweis auf die falsche Zeile, Aufschub-Adresse ohne Text) in die Closure §7.

Dieser Report ist ein **Lauf-Beleg** (dieser Stand, dieser Lauf); er ersetzt weder Review noch Closure.
