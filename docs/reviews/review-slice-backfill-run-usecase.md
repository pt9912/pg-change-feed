# Review-Report: slice-backfill-run-usecase — 2026-09-24

**Review-Art:** Code — der Diff führt die Domäne `BackfillRun`, den Inbound Port
`BackfillTableUseCase`, drei Outbound-Fähigkeits-Ports (Annahme, Run-Zustand,
Schreiber), den Use Case `usecase/backfill` samt Fakes-Tests und `model.ChangeIDFor`
(vom WAL-Mapper gerufen) ein; geprüft gegen Plan, ADRs und `AGENTS.md` Hard Rules
(Modul 10 §Drei Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe
(Modul 11).

**Gegenstand:** Slice `slice-backfill-run-usecase`, Diff-Range `455bcdef..HEAD`
(`f4ae3b2d`). Slice-Commits `d2e24986`, `3df8565d` (Lifecycle-Moves), `643582b0`
(Verantwortlich), `510b247e` (Domäne, Ports, Use Case, Tests, `ChangeIDFor`),
`9d4e9de8` (Klasse für unbekannte Schema-Version), `a1ffb9be`, `f4ae3b2d` (Plan).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(seither um weitere HIGH-Klassen ergänzt, u. a. Zahl-im-Träger, Beleg-Satz,
Zusage-ohne-Eingabeseite, Kommentar-Chronik, Handbuch-Zug).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-24.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis)*. Dieser Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-NNN` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** (`v6.9.0` ·
> `regelwerk/<datei>.md` §<Abschnitt>). Ein `pfad`-Feld auf den **geprüften
> Gegenstand** zitiert den Stand des Laufs und darf ihn festhalten.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-backfill-run-usecase` (§1 Ziel, §2 DoD als Prüfmaßstab für
  Plan-Zusagen, §3 Plan, Festlegungen ohne Vorgabe, Suchlauf-Feld, §6 Risiken)
  und Welle `welle-backfill-bestand`; die Pläne `slice-backfill-run-store`,
  `slice-backfill-sql-administration`, `slice-backfill-bench-richtgroesse`,
  `slice-backfill-e2e` und `slice-transformationen-backfill-pfad` als Träger
  (§3.13)
- [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  Teilfrage 2 (Herkunft, Row-Image-Funktion), 3 (Position), 4 (Atomarität,
  Fail-closed), 5 (Fehlerklassen, Wecksignal), 6 (Ordnung, `change_id`), 7
  (Retention) und die akzeptierten Negative,
  [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  Festlegung 1 (Rollenschnitt, `Admit`), 2 (erneute Prüfung, Aufnahme), 3
  (Schätzung „unbekannt"),
  [`ADR-0115`](../plan/adr/0115-backfill-spaltenwerte-text-ergebnisformat.md),
  [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md),
  [`ADR-0028`](../plan/adr/0028-inbound-use-cases.md),
  [`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md),
  [`ADR-0002`](../plan/adr/0002-abhaengigkeitsrichtung.md),
  [`ADR-0003`](../plan/adr/0003-physische-modulgrenzen.md),
  [`ADR-0012`](../plan/adr/0012-at-least-once.md),
  [`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md),
  [`ADR-0040`](../plan/adr/0040-clockport.md),
  [`ADR-0055`](../plan/adr/0055-nats-change-notification-wecksignal.md),
  [`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md)
- [`LH-FA-CAP-009`](../../spec/lastenheft.md),
  [`LH-FA-CAP-004`](../../spec/lastenheft.md),
  [`LH-QA-SEC-004`](../../spec/lastenheft.md),
  [`LH-FA-SCH-005`](../../spec/lastenheft.md),
  [`LH-FA-CAP-006.a`](../../spec/pflichtenheft.md),
  [`SPEC-008`](../../spec/pflichtenheft.md) (Fehlerklassen),
  [`SPEC-029`](../../spec/pflichtenheft.md) (Feldform des Run-Zustands),
  [`ARC-003`](../../spec/architecture.md),
  [`ARC-004`](../../spec/architecture.md) (Sequenz „Bestand als Backfill überführen")
- `AGENTS.md` (Hard Rules §3.1, §3.2, §3.3, §3.6, §3.7, §3.9, §3.11, §3.12,
  §3.13), `harness/conventions.md` (MR-000/MR-001)
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild
  `docs/reviews/review-slice-backfill-snapshot-reader.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem
Implementer-Bericht übernommen; Exit-Codes ungepiped gesichert):

- **Gates am Stand `f4ae3b2d`:** `make test` (Race-Detector) Exit 0;
  `make a-check` Exit 0 („gesamt: 0 Befund(e)"); `make coverage-gate` Exit 0,
  gedruckt „coverage-gate: OK — Coverage 84.40% erfüllt Schwelle 80%". Kein
  DB-Sensor nötig, bestätigt: `git diff --name-only 455bcdef..HEAD` nennt außer
  Plan-Datei nur `internal/domain/…`, `internal/application/…` und eine Zeile in
  `internal/adapters/driving/replication/mapper/mapper.go` — kein
  `postgres*`-Paket, kein Skript, kein Dockerfile; die Tests des Mappers laufen in
  `make test` (Mutation 23 unten färbt sie rot).
- **Byte-Gleichheit der WAL-Change-IDs:** Parent-Stand
  `fmt.Sprintf("%s-%d", a.open.tx.ID, sequence)` gegen `ChangeIDFor(tx, sequence)`
  (`fmt.Sprintf("%s-%d", tx, sequence)`): dieselbe Formatzeichenfolge,
  `sequence` ist in beiden `int64`, die Transaktions-Kennung ein `string`-Typ;
  gebunden über `TestConsumeFullTransaction` und `TestConsumeTwoTransactions`
  (Mutation 23: `_` statt `-` färbt beide rot).
- **Ordnung „vor jeder WAL-Kennung":** WAL-Kennungen bilden
  `strconv.FormatUint(uint64(event.XID), 10)` (`mapper.go:155`), also
  Ziffernfolgen ohne führende Null (XID 0 ist ungültig); `SelectChanges` ordnet
  `ORDER BY t.commit_position, c.transaction_id, c.sequence`
  (`queries.go:73`), `transaction_id` als Text. Das erste Zeichen `0` sortiert
  vor `1`…`9` in jeder Text-Kollation, weil Ziffern vor Buchstaben und
  Satzzeichen den ersten Vergleich entscheiden; die Blöcke desselben Runs
  unterscheiden sich nur in der achtstelligen Nummer. Die Behauptung trägt.
  `TestBackfillTransactionID` bindet sie mit dem Byte-Vergleich von Go, nicht mit
  der Sortierung der Datenbank (siehe F-13).
- **Suchlauf-Feld (Plan §3) nachgemessen:** `git grep -n -i backfill` über
  `internal/**/*.go`: Parent `643582b0` 80 Zeilen in 18 Dateien, Stand `9d4e9de8`
  362 Zeilen in 26 Dateien; `git grep -n 'ports/outbound\|port/outbound'` über
  `docs spec harness`: an beiden Ständen 44 Zeilen in 26 Dateien; die drei
  Fundstellen der alten Bedeutung (`mapper.go:334`, `schemastore.go:106`,
  `queries.go:286`) stehen an den genannten Zeilen; `classifyRunError` steht in
  `internal/bootstrap/wiring.go:1378`, der Diff berührt `internal/bootstrap` nicht.
  Die Zahlen stimmen; die vierte Zeile: F-6.
- **§3.13-Träger (Suche selbst gefahren, Parent `643582b0` und Stand):** die
  Pläne `slice-backfill-run-store` (Zeilen 43–66: drei Adapter ohne
  Methodennamen, Grants, `Admit` in einer Transaktion),
  `slice-backfill-sql-administration` (51–54: `Request`, danach Wecksignal),
  `slice-backfill-bench-richtgroesse` (63, 155: Warnung (1) über `Admit`,
  Warnung (2) über das Fortschritts-Update) und
  `slice-transformationen-backfill-pfad` (135, 150) beschreiben nichts, was der
  Diff falsch macht; die drei gemeldeten Übergaben (`SchemaStorePort` in `Ports`,
  `Publication` im Command, `RecordProgress`/`Finish` tragen den ganzen Run)
  stimmen mit dem Code überein. „Keine fremde Datei geändert" ist gedeckt
  (Dateiliste des Diffs: ein Plan, dieser Slice). Nicht gemeldet: F-4, F-5.
- **Kommentare (§3.7):** hinzugefügte Zeilen in `*.go` per Textsuche auf
  `slice-`/`welle-`/„seit "/„jetzt"/„vorher"/„früher"/„nur noch"/„bisher"/
  „nicht mehr"/„wäre"/„würde"/„künftig"/„später"/„hätte"/„sonst": keine Chronik,
  kein Konjunktiv über eine verworfene Alternative in Produktionscode; die
  Treffer sind `ADR-`/`SPEC-`-Anker und das Wort „nicht mehr" in einer
  Fehlermeldung (`service.go:361`). Test-Godoc trägt Provenienz mit Subjekt
  `TestXyz` (zulässig).
- **Commit-Struktur:** `git log --name-status 455bcdef..HEAD`: `d2e24986` und
  `3df8565d` sind Renames mit Ähnlichkeit 100 % (rein, §3.3 erfüllt); Inhalt
  getrennt in `643582b0`. Alle Betreffs tragen `LH-FA-CAP-009`/`ADR-*`, keine
  `SPEC-`/`ARC-`-Kennung im Betreff; `make commit-traceability` läuft in
  `make gates` (siehe Abschluss).

### Mutationen (Eingabeseite, selbst ausgeführt)

Je Mutation ein Lauf `make test`, Datei danach per `git checkout` zurückgenommen;
Endstand `git status` und `git diff` leer. 38 wirksame Mutationen (drei Fassungen
scheiterten zunächst am Übersetzen und sind in korrigierter Form gezählt, eine war
wirkungslos und ist nicht gezählt): 30 rot, **8 grün**.

| # | Mutation | Ort | Ergebnis |
|---|---|---|---|
| 1 | `Admit` vor die Vorbedingungen gezogen | `service.go` (`Request`) | rot: `TestRequestAdmitsLastAfterPreconditionsAndEstimate`, `…PreconditionFailuresLeaveNoRun`, `…ErrorsBeforeAdmit` |
| 2 | `known` ignoriert (`known \|\| true`) | `service.go` (`Request`) | rot: `TestRequestCarriesEstimate` |
| 3 | Ausschluss-Vergleich je Block entfernt | `service.go` (`copyBlocks`) | rot: `TestExecuteFailClosed` |
| 4 | Ausschluss-Vergleich vor dem Commit entfernt | `service.go` (`copyBlocks`) | rot: `TestExecuteFailClosed` |
| 5 | Vergleich der Tabellen-Kennung in `stillBound` entfernt | `service.go` | rot: `TestExecuteFailClosed` |
| 6 | Fehler von `stillBound` verworfen | `service.go` | rot: `TestExecuteFailClosed` |
| 7 | `defer closeSnapshot` entfernt | `service.go` | rot: fünf Tests |
| 8 | `Close` auf `ctx` statt `WithoutCancel` | `service.go` | rot: `TestExecuteContextEndedInterrupts` |
| 9 | `Rollback` auf `ctx` statt `WithoutCancel` | `service.go` | rot: `TestExecuteContextEndedInterrupts` |
| 10 | `Finish` auf `ctx` statt `WithoutCancel` | `service.go` (`conclude`) | rot: `TestExecuteContextEndedInterrupts` |
| 11 | `committed_at` aus `run.StartedAt` statt der Uhr nach dem Öffnen | `service.go` | rot: `TestExecuteCommittedAtIsSnapshotTime` |
| 12 | Sentinel `ErrSnapshotTransient` → `storage` | `service.go` (`classifyError`) | rot: `TestExecuteClassifiesFailures` |
| 13 | `ErrSchemaVersionUnknown` aus der Abbildung | `service.go` | rot: `TestExecuteClassifiesFailures` |
| 14 | Zweig „`queued` bleibt bei Kontextende `queued`" entfernt | `service.go` (`conclude`) | rot: `TestExecuteContextEndedBeforeStartLeavesQueued` |
| 15 | Wecksignal auch bei 0 Zeilen | `service.go` | rot: `TestExecuteEmptyTable` |
| 16 | Wecksignal: Schema und Tabelle vertauscht | `service.go` (`wake`) | rot: `TestExecuteWritesAllBlocksInOneTransaction` |
| 17 | Sequenz ab 0 | `service.go` (`build`) | rot: acht Tests |
| 18 | Herkunft `wal` statt `backfill` | `service.go` (`build`) | rot: `TestExecuteWritesAllBlocksInOneTransaction` |
| 19 | leere Publication in der erneuten Prüfung | `service.go` (`Execute`) | rot: `TestExecuteRechecksPreconditions` |
| 20 | erneute Prüfung ohne `Published` | `service.go` (`Execute`) | rot: `TestExecuteRechecksPreconditions` |
| 21 | Ausschluss-Schlüssel `run.Table` statt `schema.table` | `service.go` (`excludedColumns`) | rot: `TestExecuteWritesAllBlocksInOneTransaction`, `TestExecuteFailClosed` |
| 22 | Fortschritt je Block entfernt | `service.go` | rot: drei Tests |
| 23 | Präfix `bf-`; `ChangeIDFor` mit `_`; `Start` aus jedem Zustand; `Interrupt` aus `queued`; Rückschritt des Zählers in `RecordProgress` und `Complete`; `IsActive` ohne `queued`; negative Schätzung | `backfillrun.go`, `change.go` | rot, je an `TestBackfillTransactionID`, `TestChangeIDFor`, `TestConsume…` (Mapper), `TestBackfillRunTransitions`, `TestBackfillRunProgressGuards`, `TestRowEstimate` |
| 24 | `finished_at` von `failed` auf `TimePoint{}` | `service.go` (`conclude`) | **grün** — F-2 |
| 25 | `finished_at` von `interrupted` auf `TimePoint{}` | `service.go` (`conclude`) | **grün** — F-2 |
| 26 | `finished_at` der leeren Tabelle auf `TimePoint{}` | `service.go` (`copyBlocks`) | **grün** — F-2 |
| 27 | Fehler der Ausschluss-Lesung **unmittelbar vor dem Commit** verworfen | `service.go:256` | **grün** — F-1 |
| 28 | Fehler der Ausschluss-Lesung **je Block** verworfen | `service.go:216` | **grün** (von der Lesung vor dem Commit gedeckt) — F-1 |
| 29 | Fehler des Fortschritts je Block verworfen | `service.go:237` | **grün** — F-3 |
| 30 | Fehler des Fortschritts mit der Anfangsposition verworfen | `service.go:192` | **grün** (vom Fortschritt je Block gedeckt) — F-3 |
| 31 | Fehler von `Finish` bei leerer Tabelle verworfen | `service.go:247` | **grün** — F-3 |

Gebunden an ihrer Eingabeseite sind: Reihenfolge der Annahme, Bedeutung von
`known`, Vergleich je Block und vor dem Commit (Inhalt, nicht Lesefehler), Bindung
mit Kennung, alle drei `WithoutCancel`-Stellen, Schließen des Snapshots, Quelle
von `committed_at`, Sentinel-Abbildung, Wecksignal (Bedingung und Argumente),
Sequenz und Herkunft, erneute Prüfung, Streaming-Reihenfolge, die Domäne.
Nicht gebunden: die Lesefehler-Zweige der Fail-closed-Prüfung (F-1), `finished_at`
der drei Endzustands-Wege (F-2), zwei Fortschritts-Fehlerzweige und der
Leer-Endzustand (F-3).

**Fakes gegen die echten Ports:** `fakeActivation.Published` prüft Schema und
Tabelle gegen `Registered` und meldet bei Abweichung einen Fehler; jeder Fake
hält seine Argumente (`gotSource`, `gotID`, `gotRunID`, …), und die Tests
vergleichen sie. `fakeSnapshot.NextBlock` liefert bei beendetem Kontext
`ErrSnapshotTransient` wie der Adapter. Schwäche: `fakeExclusion` und `fakeRuns`
scheitern nur dauerhaft (`err`, `progressErr` gelten für jeden Aufruf); ein Fehler
„ab Aufruf n" ist nicht ausdrückbar — Ursache von F-1 und F-3.

---

## Findings

### F-1 — Lesefehler des Ausschlussstands vor dem Commit ist an seiner Eingabe nicht gebunden

- `kategorie`: HIGH
- `quelle`: Zusage ohne Bindung an ihre Eingabeseite (Reviewer-Skill, HIGH-Liste);
  [`LH-QA-SEC-004`](../../spec/lastenheft.md);
  [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 4
  („unmittelbar vor dem Commit prüft der Run … jede Abweichung rollt den Run
  zurück"); Plan §2 DoD Punkt 2 („je Test eine Mutation der Prüfung, die den Test
  rot färbt")
- `pfad`: `internal/application/usecase/backfill/service.go:256-262` (Lesung vor
  dem Commit), `:216-219` (Lesung je Block);
  `internal/application/usecase/backfill/service_test.go:968` (einziger Fall mit
  Lesefehler), `:1029-1112` (Tabelle `TestExecuteFailClosed`, kein Lesefehler-Fall)
- `befund`: Der Zweig „Ausschlussstand nicht lesbar ⇒ kein Commit" hält nur, weil
  der einzige Test jede Lesung dauerhaft scheitern lässt und schon die erste den
  Run beendet. Verwirft man den Fehler der Lesung unmittelbar vor `Commit`
  (Mutation 27) oder den der Lesung je Block (Mutation 28), bleibt `make test`
  grün; bei leerem Ausgangsstand vergleicht Mutation 27 dann `nil` mit `nil` und
  committet, ohne dass der Stand vor dem Commit je gelesen wurde.
- `verifizierbar`: ja — `make test` mit der Mutation (`excluded, _ :=` an
  `service.go:256`); der Fake kann einen Fehler erst ab dem n-ten Aufruf nicht
  ausdrücken.
- `klasse`: Zusage ohne Bindung an ihre Eingabeseite

### F-2 — `finished_at` der Endzustände ist an der Uhr-Eingabe nicht gebunden

- `kategorie`: MEDIUM
- `quelle`: fehlende Negativtests bei neuem öffentlichem Vertrag (Reviewer-Skill);
  Port-Doku `internal/application/port/outbound/backfillrun.go:35-38`
  („`failed` und `interrupted` mit `finished_at` und Fehlertext");
  [`SPEC-029`](../../spec/pflichtenheft.md) (`finished_at`: Übergang nach
  `completed`, `failed`, `interrupted`);
  [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  Festlegung 3 Punkt 2 (Kopierdauer = `finished_at − started_at`)
- `pfad`: `internal/application/usecase/backfill/service.go:291`, `:293` (`conclude`),
  `:243` (leere Tabelle); `internal/application/usecase/backfill/service_test.go:983`,
  `:1164` (vergleichen den vom Fake festgehaltenen Run mit dem Ergebnis — beide
  stammen aus derselben Zeile des Use Case)
- `befund`: Kein Test liest `FinishedAt` von `failed`, `interrupted` oder dem
  `completed` der leeren Tabelle; `TimePoint{}` statt der Uhr in `conclude` (Mutationen
  24, 25) und in der leeren Tabelle (26) lässt `make test` grün. Nur der Commit-Pfad
  bindet `FinishedAt` (`service_test.go:752`).
- `verifizierbar`: ja — `make test` mit den drei Mutationen.
- `klasse`: fehlende Negativtests bei neuem öffentlichem Vertrag

### F-3 — Fehlerzweige der Fortschritts-Schreibung und des Leer-Endzustands nur im Aggregat oder gar nicht gebunden

- `kategorie`: MEDIUM
- `quelle`: Zusage ohne Bindung an ihre Eingabeseite (Reviewer-Skill; nach
  Wirkung als MEDIUM eingeordnet, kein Sicherheitspfad); Port-Doku
  `backfillrun.go:35-38` (`Finish` „hält den Endzustand … `completed` für eine leere
  Tabelle" fest); Inbound-Doku `port/inbound/backfill.go` (`Execute`: „Fehler …
  dass der Run-Zustand nicht festgehalten werden konnte")
- `pfad`: `internal/application/usecase/backfill/service.go:237-239` (Fortschritt je
  Block), `:192-194` (Fortschritt mit Anfangsposition), `:247-249` (`Finish` der leeren
  Tabelle); `internal/application/usecase/backfill/service_test.go:965`
  (`progressErr` gilt für jeden Aufruf), `:1219-1230` (`Finish`-Fehler nur im
  `conclude`-Pfad)
- `befund`: Verwirft man den Fehler des Fortschritts je Block (Mutation 29) oder den
  der Anfangsposition (30), bleibt `make test` grün, weil der jeweils andere Aufruf
  den dauerhaft scheiternden Fake auffängt; verwirft man den Fehler von `Finish`
  bei leerer Tabelle (31), meldet `Execute` `completed` ohne dass der Endzustand
  festgehalten ist, und kein Test färbt sich rot.
- `verifizierbar`: ja — `make test` mit den drei Mutationen.
- `klasse`: Zusage ohne Bindung an ihre Eingabeseite

### F-4 — Zeitbegrenzung der Adapter unter `WithoutCancel` hat keinen Träger

- `kategorie`: MEDIUM
- `quelle`: unklare Fehlerbehandlung am Rand des Spec-Bereichs (Reviewer-Skill);
  Modul 8 „Kein Pfeil ohne benennbares Artefakt" (`v6.9.0` ·
  `regelwerk/modul-08-agentenrollen.md` §Die neun Übergaben);
  [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md)
- `pfad`: `internal/application/usecase/backfill/service.go:298` (`Finish`), `:417`
  (`Close`), `:428` (`Rollback`); `internal/application/port/outbound/backfillrun.go:35-39`,
  `internal/application/port/outbound/backfillwriter.go:34-36`,
  `internal/application/port/outbound/tablesnapshot.go:84-86` (keine Aussage zur
  Dauer); Plan `slice-backfill-run-store` §1 (keine Aussage zur Dauer)
- `befund`: Rollback, `Finish` und `Close` laufen auf einem vom Abbruch gelösten
  Kontext ohne Zeitgrenze; ein hängender Adapter hält `Execute` dadurch über das
  Ende des Aufrufer-Kontexts hinaus. Die Zusage „die Dauer begrenzt der Adapter"
  steht allein in Plan §3 (Festlegungen ohne Vorgabe) dieses Slice — weder in einer
  Port-Doku noch im Plan des Slice, der den Run-Zustands- und den Schreiber-Adapter
  baut; der Snapshot-Adapter begrenzt seine Schließ-Dauer selbst
  (`postgressnapshot/snapshot.go:319`, `closeTimeout`), die beiden anderen Adapter
  existieren nicht.
- `verifizierbar`: nein — kein Gate liest die Adapter-Pflicht; nachzumessen durch
  Lesen der Port-Doku und der Pläne.
- `klasse`: Übergabe ohne Träger (Adapter-Pflicht steht nur im Plan des Aufrufers)

### F-5 — Vertrag von `Finish` gegenüber einem beendeten Run lässt Fehler oder Erfolg offen

- `kategorie`: LOW
- `quelle`: unklare Fehlerbehandlung am Rand (Maintainability);
  [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  Festlegung 1
- `pfad`: `internal/application/port/outbound/backfillrun.go:16-20` („ändert keinen
  beendeten Run"), `:35-39` (`Finish`); `internal/application/usecase/backfill/service.go:298-300`
- `befund`: Die Port-Doku sagt, dass keine schreibende Operation einen beendeten Run
  ändert, nennt aber nicht, was `Finish` für einen bereits beendeten Run zurückgibt
  (Fehler oder wirkungsloser Erfolg); `conclude` liest jeden Fehler als „Run-Zustand
  nicht festgehalten". Trifft ein `Commit` mit Verbindungsabbruch nach dem
  serverseitigen Commit ein, liegt genau dieser Fall vor: die Zeile ist
  `completed`, der Use Case meldet `failed` oder `interrupted` bzw. einen Fehler.
- `verifizierbar`: nein — der Adapter existiert nicht.
- `klasse`: unklare Fehlerbehandlung am Rand

### F-6 — Beleg-Befehl des Suchlaufs liefert kopiert keinen Treffer

- `kategorie`: LOW
- `quelle`: Beleg trägt seinen Satz nicht (Reviewer-Skill); `AGENTS.md` §3.13
- `pfad`: `docs/plan/planning/in-progress/slice-backfill-run-usecase.md:225` (§3,
  Suchlauf-Feld, vierte Zeile)
- `befund`: Der genannte Befehl `git grep -n -E 'Admit\|Annahme-Port\|Run-Zustands\|Schreiber' …`
  ergibt mit `-E` und dem Zeichen `\|` wörtlich aus der Zelle kopiert 0 Treffer
  (nachgefahren am Stand `9d4e9de8`); mit unmaskiertem `|` liefert er die genannten
  Zeilen von `run-store`, `sql-administration` und `bench-richtgroesse`, dazu Treffer
  auf „Schreibern" in `slice-backfill-e2e`, die die Zeile nicht nennt, und **keinen**
  in `slice-transformationen-backfill-pfad` — die dort genannten Zeilen 135 und 150
  stehen an diesen Stellen, stammen aber aus dem Lesen, nicht aus dem Befehl. Die
  Zeilenangaben selbst und die Zahlen der übrigen drei Zeilen stimmen.
- `verifizierbar`: ja — Befehl aus der Zelle ausführen.
- `klasse`: Beleg trägt seinen Satz nicht

### F-7 — Grenze der Fail-closed-Prüfung ist nirgends benannt

- `kategorie`: LOW
- `quelle`: [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  Teilfrage 4 (Vergleich auf Stand-Gleichheit);
  [`LH-QA-SEC-004`](../../spec/lastenheft.md); Plan §6 zweites Risiko
- `pfad`: `internal/application/usecase/backfill/service.go:166-172`
  (Doku `copyBlocks`), `:208-262`; Plan §6
- `befund`: Die Prüfung erkennt jeden Stand, den eine Lesung sieht (Zwischenblock,
  Ende, auch „Ausschluss, dann Einschluss" bei einer Lesung dazwischen). Ein
  Ausschluss, der zwischen zwei Lesungen gesetzt und wieder zurückgenommen wird,
  ohne dass eine Lesung ihn sieht, bleibt unsichtbar, weil der Stand
  (`ExcludedColumns`, abgeleitet aus den `applied`-Zeilen) keine Historie trägt;
  weder Code-Kommentar noch §6 nennen diese Grenze. Bewertung: die Zusage der ADR
  (Stand-Gleichheit) ist erfüllt; die Wirkung entspricht dem WAL-Pfad, in dem ein
  Change vor dem Ausschluss den Wert behält — kein Befund höherer Schwere.
- `verifizierbar`: nein — kein Test kann die Lücke ohne Historie im Port schließen.
- `klasse`: Grenze einer Zusage nicht benannt

### F-8 — `schema_version` der Backfill-Changes: akzeptiertes Negativ deckt den Fall vor dem Run nicht wörtlich

- `kategorie`: INFO
- `quelle`: [`LH-FA-SCH-005`](../../spec/lastenheft.md);
  [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) §Konsequenzen,
  akzeptiertes Negativ „Schema-Version-Verweis"
- `pfad`: `internal/application/usecase/backfill/service.go:132-139`, `:366-377`;
  `internal/adapters/driving/replication/mapper/mapper.go:340-388`
- `befund`: Der Use Case liest `CurrentVersion` vor dem Snapshot. Die Version wechselt
  im Bestand allein in `observeRelation` mit der nächsten Relation-Nachricht des
  WAL-Pfads, und die statische Erstaktivierung legt die Versions-Zeile ohne
  `TableSchema` an (nachgetragen erst mit der ersten Relation-Nachricht,
  `mapper.go:331-334`). Eine Tabelle, die seit der Aktivierung oder seit einem
  `ALTER TABLE … ADD COLUMN` keine WAL-Änderung hatte, bekommt Backfill-Changes,
  deren Bild neuere Spalten trägt als ihre `schema_version`, oder die auf eine
  Version ohne `TableSchema` verweisen. Das akzeptierte Negativ der ADR nennt „eine
  kompatible Spalten-Erweiterung **während** des Runs" und die Selbstbeschreibung
  des Bilds; die Fälle davor sind dort nicht genannt, das Argument (Bild ist
  selbstbeschreibendes JSON) trägt sie ebenso. Die Boundary von
  [`LH-FA-SCH-005`](../../spec/lastenheft.md) (Versionen vor und nach einer
  Änderung unterscheidbar) bleibt erfüllt. Die Grenze steht weder im Code noch in
  Plan §6.
- `verifizierbar`: nein — kein Gate; Adresse: Architect-Sicht auf die Reichweite
  des akzeptierten Negativs, kein Fixrunden-Anlass.
- `klasse`: Grenze einer Zusage nicht benannt

### F-9 — `internal` als Rückfall der Run-Klassen und `classifyError` als Übergabepunkt

- `kategorie`: INFO
- `quelle`: [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  Teilfrage 5 (fünf Klassen des Runs);
  [`ADR-0023`](../plan/adr/0023-fehlerklassifikation.md),
  [`SPEC-008`](../../spec/pflichtenheft.md) (sieben Klassen)
- `pfad`: `internal/application/usecase/backfill/service.go:304-330`
  (`internal` in Zeile 328); `internal/domain/model/errorstate.go:14-21`
- `befund`: Die ADR nennt für den Run fünf Klassen; `internal` steht in der
  geschlossenen Menge von `SPEC-008`/`ADR-0023` (`NewErrorClass` akzeptiert alle
  sieben) und ist im Plan als Festlegung ohne Vorgabe genannt — zulässig.
  `ErrSchemaVersionUnknown` → `configuration` ist eine Setzung ohne ADR-Text, die
  zur Klasse „Konfiguration" der Tabelle passt. Die Klasse `schema` vergibt der Run
  nicht; ob der Backfill-Pfad der Transformationen sie für „Regel nicht anwendbar"
  braucht, ist eine offene Frage an den Architect
  (`slice-transformationen-backfill-pfad` nennt die Nichtanwendbarkeit ohne Klasse).
- `verifizierbar`: nein.
- `klasse`: Übergabepunkt ohne Entscheidung (Architect)

### F-10 — Vergleich der Ausschlussstände doppelt; je Block eine Abfrage der ganzen Quelle

- `kategorie`: INFO
- `quelle`: Maintainability;
  [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 4
- `pfad`: `internal/application/usecase/backfill/service.go:225-227`, `:260-262`,
  `:381-387`
- `befund`: Der Vergleich `sameNames(baseline, excluded)` mit demselben Fehler steht
  an zwei Stellen (je Block, vor dem Commit); eine Änderung der Vergleichs-Semantik
  braucht beide, eine Bindung besteht für beide (Mutationen 3, 4). Kein Befund über
  die Doppelung hinaus. `ExcludedColumns` liefert je Aufruf den Stand der ganzen
  Quelle (Port-Doku), und der Use Case ruft ihn je Block; die ADR verlangt „jeder
  Block liest neu", die Kosten bei großer Blockzahl sind nicht gemessen.
- `verifizierbar`: nein.
- `klasse`: Wiederholung im Vergleich (Maintainability)

### F-11 — `Verantwortlich` gesetzt nach `next → in-progress`

- `kategorie`: INFO
- `quelle`: `v6.9.0` · `regelwerk/modul-05-planning-harness.md` §Lifecycle als State
  Machine („`open → next` setzt den Verantwortlichen")
- `pfad`: `docs/plan/planning/in-progress/slice-backfill-run-usecase.md:28`
  (Commits `d2e24986`, `3df8565d`, `643582b0`)
- `befund`: Das Feld stand bei `open → next` und bei `next → in-progress` auf „—
  (noch nicht priorisiert)" und wurde erst im Commit `643582b0` in `in-progress/`
  gesetzt; der Vorgänger `slice-backfill-snapshot-reader` setzte es in `next/`.
  Das Feld steht vor dem ersten Code-Commit (`510b247e`); die Moves sind rein.
- `verifizierbar`: nein — Deklaration, kein Sensor.
- `klasse`: Lifecycle-Reihenfolge (Deklaration)

### F-12 — Kommentar in `conclude` nennt Verhalten, das ein anderer Slice trägt

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.7 (Klasse Zusage/Rang-Zeiger);
  `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`
- `pfad`: `internal/application/usecase/backfill/service.go:278-281`
- `befund`: Der Kommentar sagt, ein `queued`-Run „wird nach einem Neustart
  aufgenommen" und ein nicht festgehaltener Run bleibe „bis zum Abgleich beim
  Prozessstart `running`"; beide Aussagen trägt der Code dieses Diffs nicht (Aufnahme
  und Abgleich liegen bei `slice-backfill-sql-administration`, der Port-Vertrag
  `InterruptRunning` steht hier). Sie stimmen mit der Architektur-Sicht
  ([`ARC-004`](../../spec/architecture.md), Sequenz „Bestand als Backfill
  überführen") und
  [`ADR-0113`](../plan/adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  Festlegung 2 überein; der Kommentar trägt keinen Rang-Zeiger darauf.
- `verifizierbar`: nein.
- `klasse`: Kommentar behauptet Verhalten außerhalb des Diffs

### F-13 — Ordnungs-Beleg in der Domäne ist ein Go-Vergleich; `TestPortsCarryNoCapturePathPort` bindet über Typnamen

- `kategorie`: INFO
- `quelle`: [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md)
  Teilfrage 6; [`LH-FA-CAP-004`](../../spec/lastenheft.md)
- `pfad`: `internal/domain/model/backfillrun_test.go:242-249`;
  `internal/application/usecase/backfill/service_test.go:676-686`
- `befund`: Die Sortier-Zusage „vor jeder WAL-Kennung" ist mit dem Byte-Vergleich von
  Go belegt, nicht mit `ORDER BY` an einer Datenbank; sie trägt (siehe oben,
  erstes Zeichen), der Nachweis an der realen Sortierung ist Sache der Slices mit
  Datenbank (`slice-backfill-run-store`, `slice-backfill-e2e`). Der Test „run-lokal,
  kein Heartbeat" prüft Feld-Typnamen der `Ports`-Struktur auf vier Teilstrings; ein
  Capture-Pfad-Port unter anderem Typnamen würde nicht auffallen. „Heartbeat
  unberührt ist strukturell" ist als Antwort tragfähig, weil der Use Case keinen
  Zugang besitzt.
- `verifizierbar`: nein.
- `klasse`: Beleg mit engerer Reichweite als die Aussage

## Negativbefunde

- geprüft, ohne Befund: `internal/domain/model/backfillrun.go` — Übergänge
  (`queued → running → completed | failed | interrupted`, `queued → failed`),
  `Interrupt` nur aus `running`, `Fail` nur aus `queued`/`running`, Wert-Semantik
  (jede Methode liefert einen neuen Run), `RowEstimate` (Nullwert unbekannt, `0` nur
  über `NewRowEstimate(0)`, negativ verboten), `BackfillTransactionID` (Blocknummer ab
  1, achtstellig, Grenze als Fehler), Feldform gegen `SPEC-029` (`started_at` leer bis
  `running`, `error_message` nur bei `failed`); sieben Mutationen rot (Nr. 23)
- geprüft, ohne Befund: `model.ChangeIDFor` und der Aufruf im WAL-Mapper — Formatzeichenfolge
  und Typen unverändert gegen den Parent, die Byte-Gleichheit ist vom Mapper-Test
  gebunden (Nr. 23); Eindeutigkeit von `<Tx>-<Sequenz>` (die Sequenz ist rein
  numerisch, nach dem letzten `-` gelesen, damit trennscharf)
- geprüft, ohne Befund: Ports (`port/inbound/backfill.go`, `port/outbound/backfilladmission.go`,
  `backfillrun.go`, `backfillwriter.go`) — Fähigkeits-Schnitt nach
  [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md), `Admit` als einzige
  Anlage-Operation (`TestPortSchnitt` bindet die Methodenmengen), keine
  Adapter-Kante (`make a-check` Exit 0), Transport-Typen am Inbound Port
  ([`ADR-0042`](../plan/adr/0042-transport-typen-am-port.md)), Doku nennt Fehler-
  Sentinels (`ErrBackfillRunActive`, `ErrBackfillStorage`); Kopf-Kommentar der
  Inbound-Datei folgt der Form der Nachbardateien; `Ports.Schemas` als achter
  Pflicht-Port ist gerechtfertigt: `model.NewChange` verlangt eine
  `SchemaVersionID`, der Port ist der bestehende Fähigkeits-Port — kein Schnittfehler
  (Reichweite des Werts: F-8)
- geprüft, ohne Befund: `usecase/backfill/service.go` `Request` — Vorbedingungen,
  Schätzung, `Admit` als letzter Schritt, keine Nebenwirkung davor (Mutationen 1, 2),
  „unbekannt" gegen bekannte `0` (drei Fälle in `TestRequestCarriesEstimate`)
- geprüft, ohne Befund: `Execute` — erneute Prüfung ohne Snapshot (Mutationen 19, 20),
  Streaming der Blöcke (`TestExecuteStreamsBlocks`), Position `X`, `committed_at` =
  Snapshot-Zeit (11), `INSERT`/`backfill`/kein `old_data`/Sequenz und Blocknummern
  ab 1 (17, 18), ein Wecksignal je Tabelle nur bei Zeilen (15, 16), Snapshot-Schließen
  und Rollback auf jedem Pfad (7–9), `interrupted` gegen `failed` am eigenen Kontext
  (10, 14), `queued` bleibt `queued`; Fehler des Aufrufs gegen Ergebnis im Endzustand
  ist im Inbound-Port eindeutig dokumentiert; keine Goroutine, keine geteilten Daten
- geprüft, ohne Befund: Fehlerklassen-Abbildung — fünf Sentinels des Snapshot-Ports,
  die Storage-Sentinels, `ErrSchemaVersionUnknown` → `configuration` gebunden
  (Nr. 12, 13); die Adapter der Aktivierung und des Schema Stores liefern
  `outbound.ErrStorage` bzw. `ErrSchemaStoreStorage` (`postgresstorage/store.go:63`,
  `schemastore.go:57`), also `storage` statt `internal`; Einordnung von `internal`
  und `schema`: F-9
- geprüft, ohne Befund: `AGENTS.md` §3.1 (kein lokales Toolchain), §3.2 (kein
  `//nolint`), §3.6 (kein Gate gelockert), §3.11 (kein host-lokaler Pfad im Diff),
  §3.12 (Zahlen des Suchlauf-Felds nachgemessen, siehe oben), Spec-Stratum (kein
  Diff in `spec/`), Traceability (alle Betreffs), Handbuch (kein Diff in
  `docs/user/`, keine `CDC_*`-Variable, keine SQL-Funktion, kein Endpunkt — keine
  Betreiber-Oberfläche), Persist-before-ACK
  ([`ADR-0012`](../plan/adr/0012-at-least-once.md)): der Run berührt weder
  `ReplicationAck` noch Heartbeat noch Live-Stream, Fehler sind run-lokal
- geprüft, ohne Befund: Abhängigkeitsrichtung
  ([`ADR-0002`](../plan/adr/0002-abhaengigkeitsrichtung.md),
  [`ADR-0003`](../plan/adr/0003-physische-modulgrenzen.md)) — `usecase/backfill`
  importiert nur `port/inbound`, `port/outbound`, Domäne; `make a-check` Exit 0
- geprüft, ohne Befund: `docs/plan/planning/in-progress/slice-backfill-run-usecase.md`
  Suchlauf-Zeilen 1–3 (Zahlen und Zeilen nachgemessen), Festlegungen gegen den Code
  (Blocknummern ab 1, `queued → failed`, Kontext, Ergebnis, Klassen, Wecksignal,
  Publication im Command), §6-Risiken tragen „bei Closure" als Ausgang — kein
  Fremdträger stehen geblieben (siehe §3.13-Zeile oben)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 3 |
| LOW | 3 |
| INFO | 6 |

**Finding-Klassen dieses Laufs:** Zusage ohne Bindung an ihre Eingabeseite (F-1, F-3) ·
fehlende Negativtests bei neuem öffentlichem Vertrag (F-2) · Übergabe ohne Träger
(F-4) · unklare Fehlerbehandlung am Rand (F-5) · Beleg trägt seinen Satz nicht (F-6) ·
Grenze einer Zusage nicht benannt (F-7, F-8) · Übergabepunkt ohne Entscheidung (F-9) ·
Wiederholung im Vergleich (F-10) · Lifecycle-Reihenfolge (F-11) · Kommentar behauptet
Verhalten außerhalb des Diffs (F-12) · Beleg mit engerer Reichweite als die Aussage (F-13)

## Verdikt

**Merge-blockierend:** ja — ein HIGH (F-1: der Lesefehler-Zweig der Fail-closed-Prüfung
vor dem Commit ist ungebunden; die Klasse steht in der HIGH-Liste des Skills, und der
Pfad ist der Sicherheitspfad von
[`LH-QA-SEC-004`](../../spec/lastenheft.md)) und drei MEDIUM (F-2 bis F-4). Die Gates am
Stand laufen grün (`make test`, `make a-check`, `make coverage-gate`); die Befunde
liegen an Stellen, die die Tests nicht an ihrer Eingabe prüfen, und an einer
Adapter-Pflicht ohne Träger. Die Domäne, die Port-Schnitte und die Ordnungs- und
Kennungs-Regeln tragen.

**Übergabe:** F-1 bis F-3, F-5, F-6, F-10 bis F-12 gehen an den Implementer (F-1 und
F-3 setzen einen Fake voraus, der einen Fehler erst ab dem n-ten Aufruf liefert). F-4
geht an den Implementer (Port-Doku) **und** an den Planner (Pflicht im Plan
`slice-backfill-run-store`); das Übergabe-Artefakt ist dieser Report. F-8 und F-9 gehen
an den Architect (Reichweite des akzeptierten Negativs der ADR; Klasse `schema` im
Backfill-Pfad der Transformationen). Die DoD-Zeile „Review durchgeführt, Report unter
`docs/reviews/` liegt vor" bleibt **offen**, weil eine Fixrunde nötig ist (Skill
§DoD-Checkbox-Nachzug). Die **Finding-Klassen** gehen in die Slice-Closure §7 und von
dort in den Zähler; „Zusage ohne Bindung an ihre Eingabeseite"
(`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`) tritt in diesem Lauf zweimal
auf, „Beleg trägt seinen Satz nicht" einmal. Dieser Report ist ein **Lauf-Beleg**;
DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).
