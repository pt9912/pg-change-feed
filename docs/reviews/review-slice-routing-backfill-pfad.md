# Review-Report: slice-routing-backfill-pfad — 2026-10-01

**Review-Art:** Code — geprüft gegen Plan, [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md),
[`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md) Festlegung 2,
[`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md) Festlegung 1,
[`LH-FA-CAP-009.a`](../../spec/pflichtenheft.md) („Fail-closed vor dem Commit“, „Ziel der Backfill-Changes“),
[`SPEC-008`](../../spec/pflichtenheft.md), [`SPEC-032`](../../spec/pflichtenheft.md) und `AGENTS.md` Hard Rules
(Modul 10). Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice [`slice-routing-backfill-pfad`](../plan/planning/done/slice-routing-backfill-pfad.md)
([`welle-routing`](../plan/planning/welle-routing.md)), Diff-Range `177bbac5~1..HEAD`: `d8c7cf4f` (Plan), `55384d2f`
(Code und Tests), `d459ba3a` (Suchlauf-Feld, Übergabe-Block); `013ece33` (Linkkorrektur) liegt in der Range und ist
reiner Link. 14 Dateien, +1019/−44 (`git diff --stat 177bbac5~1 HEAD`).

**Skill:** `.harness/skills/reviewer.md` @ `c5207cc1`. **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-01.

**Ablage:** Der Reviewer-Lauf hat diesen Report selbst geschrieben (Write-Werkzeug). Alle Mutationen liefen an einer
Kopie der getrackten Dateien im Scratchpad (Mutation per `sed … > Kopie`, nie `sed -i`, nie eine Umleitung auf eine
Repo-Datei); Unit-Mutationen im gepinnten Toolchain-Image mit `--network none` und `-race` (Aufrufform von `make test`,
Quelle: die Kopie), Store-Mutationen über `tools/harness/run-store-tests.sh` aus der Kopie. Es gab keine verweigerte
Aktion.

**Eingangs-Kontext:**

- Slice-Plan (§1–§3 mit Suchlauf-Feld, §6 Risiken), [`welle-routing`](../plan/planning/welle-routing.md)
- [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md),
  [`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md),
  [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md),
  [`ADR-0111`](../plan/adr/0111-backfill-bestand-snapshot-bulk-copy.md),
  [`ADR-0117`](../plan/adr/0117-backfill-run-fehlerklasse-schema.md)
- [`LH-FA-CFG-008`](../../spec/lastenheft.md), [`LH-FA-CAP-009`](../../spec/lastenheft.md),
  [`SPEC-008`](../../spec/pflichtenheft.md), [`SPEC-032`](../../spec/pflichtenheft.md)
- `AGENTS.md` §3.1, §3.7, §3.12, §3.13, §3.15; Vorgänger-Review
  [`review-slice-routing-kern-label`](review-slice-routing-kern-label.md)

---

## Findings

### F-1 — Lesefehler des Routing-Regelstands endet `storage`; Spec und ADR sagen `configuration`

- `kategorie`: HIGH
- `quelle`: [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md) Festlegung 1 („Eine Abweichung — auch ein nicht lesbarer Stand — rollt den Run zurück: `failed`, Klasse `configuration`“); [`LH-FA-CAP-009.a`](../../spec/pflichtenheft.md) „Fail-closed vor dem Commit“ („jede Abweichung … Fehlerklasse `configuration`“ und „Ein nicht lesbarer Stand gilt als Abweichung“); Plan §1 und DoD-Punkt 2
- `pfad`: `internal/application/usecase/backfill/service.go:230-235, 270-274, 328-334` (`routingRules` reicht den Fehler des Ports unverändert durch, `classifyError` bildet ihn auf die Klasse der Ursache ab); `internal/application/usecase/backfill/routing_test.go` (`TestExecuteRoutingReadFailureEndsRun`, erwartet `storage: `)
- `befund`: Ein Port-Fehler `ErrStorage` bei der Lesung 1 bis 5 des Routing-Standes endet den Run mit der Klasse `storage`, nicht mit `configuration`. Das ist das Verhalten des Ausschlussstandes (`TestExecuteExclusionReadFailure`) und trägt die Spec-Zeile „nicht lesbar = Abweichung = `configuration`“ nicht; der Test legt die Abweichung vom Wortlaut fest, statt sie zu benennen. Der Plan-Text wiederholt „nicht lesbarer Stand … Klasse `configuration`, wie bei den Transformationsregeln“, die Transformationsregeln verhalten sich aber ebenfalls anders.
- `verifizierbar`: ja — `make test` (Test grün, aber gegen die Klasse `storage` geschrieben); Gegenprobe ist das Lesen von `LH-FA-CAP-009.a` und `ADR-0139`.
- `klasse`: Implementierung weicht vom ADR-Wortlaut ab (`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab`, bisher 2×, dies wäre das dritte Auftreten)

### F-2 — Zitat nennt die falsche Stelle: Fehlerklassen-Tabelle steht in Handbuch-§6, nicht §5

- `kategorie`: LOW
- `quelle`: Maintainability / Beleg trägt seinen Satz nicht (Abschnittsnummer, `AGENTS.md` §3.13 Aufschub-Adresse)
- `pfad`: `docs/plan/planning/open/slice-routing-betriebsdoku.md:101-105` (Übergabe-Block: „Zeile `schema` der Fehlerklassen-Tabelle in §5“); `docs/plan/planning/in-progress/slice-routing-backfill-pfad.md` §3 letzte Tabellenzeile („§5 Fehlerklassen“)
- `befund`: Aufgeschlagen: `docs/user/benutzerhandbuch.md:1876` (Zeile `schema`) liegt unter `### Fehlerklassen` in `## 6. Fehlerbehebung` (Zeile 1864); `## 5.` ist „Konfiguration“. Der Zeilen-Lokator 1876 und der Lokator 756 („Ausgeschlossene Spalten“ unter `### Bestand als Backfill überführen`, Zeile 541) sind richtig, die Abschnittsnummer nicht; der Betriebsdoku-Plan selbst nennt an anderer Stelle §6.
- `verifizierbar`: ja — `git grep -n -E "Ausgeschlossene Spalten|Transformationsregel ist auf die Relation" -- docs/user/benutzerhandbuch.md` plus Überschriften-Lesung.
- `klasse`: Zitat nennt die falsche Stelle

### F-3 — `sameSet`-Kommentar nennt nur die Transformationsregel als Vergleichstyp

- `kategorie`: LOW
- `quelle`: Maintainability / Nachzug widerspricht dem Nachbarn im selben Träger (`AGENTS.md` §3.7, §3.13)
- `pfad`: `internal/application/usecase/backfill/service.go:647-649`
- `befund`: Der Godoc von `sameSet` begründet die Typschranke mit „`model.Transformation` ist über alle seine Felder vergleichbar“; der Diff ruft die Funktion jetzt auch mit `model.RouteRule` (sechs Felder, ebenfalls vergleichbar) auf, der Kommentar nennt das nicht. Der Suchlauf in Plan §3 hat die Stelle nicht erfasst (Muster `Transformation` im Nicht-Test-Code ohne den Zusatz „vergleichbar“).
- `verifizierbar`: ja — Lesen von `service.go:647-662`.
- `klasse`: Nachzug widerspricht dem Nachbarn im selben Träger

### F-4 — Store-Test belegt die Kette „echter Snapshot-Adapter → Run → `route_target`“ nicht

- `kategorie`: INFO
- `quelle`: Plan §3 („Abweichung vom Plan (Store-Test)“), DoD-Punkt 1 („Store- und Snapshot-Test“)
- `pfad`: `internal/bootstrap/backfill_routing_store_internal_test.go`; `internal/adapters/driven/postgressnapshot/snapshot_test.go` (`checkRouteParity`)
- `befund`: Die Abweichung ist ehrlich benannt und begründet: `make test-store` fährt ohne `wal_level=logical`, der Snapshot-Adapter öffnet seinen temporären Slot nur unter `make test-replication` — nachgelesen in `tools/harness/run-store-tests.sh` und am Lauf. Die Teile sind einzeln belegt (Run, Regelstand, Writer, `cdc.changes` im Store-Test; Gleichheit der Auswertung über Werte des echten Snapshot-Adapters im Typ-Satz-Test, `make test-replication`, PostgreSQL 18). Keine Zeile belegt sie zusammen (Adapter-Zeilen als Eingabe des Run bis zur Persistenz); diese Lücke gehört dem E2E-Slice [`slice-routing-e2e`](../plan/planning/open/slice-routing-e2e.md). Der Typ-Satz-Test läuft nur an einer Version (18), nicht an 17 — in der Beschreibung als Grenze benennbar.
- `verifizierbar`: nein (Zuschnitt, kein Gate-Lauf).
- `klasse`: Belegkette in Teilen, nicht als Ganzes

### F-5 — Mutationsangaben in Test-Godocs: Store-Mutation nicht in der benannten Form nachgefahren

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12 (Ursprung einer Aussage), Skill „Beleg trägt seinen Satz nicht“
- `pfad`: `internal/bootstrap/backfill_routing_store_internal_test.go` (Godoc: „die Spalte `route_target` aus `InsertBackfillChange` streichen — die Ziele in `cdc.changes` sind leer“)
- `befund`: Nachgefahren als Variante (`$10` im `VALUES` durch `NULL` ersetzt): rot, aber mit „mismatched param and argument count“ und Run `failed` statt leerer Ziele; das Streichen von Spalte und Argument zugleich habe ich nicht gefahren. Die zweite Store-Mutation (`'applied'` → `'pending'` in `SelectAppliedRoutingRequests`) ist rot wie beschrieben. Die Godoc-Aussagen der Unit-Tests sind alle nachgefahren (siehe Schwerpunkt (g)); keine Mutation ist in den Godocs behauptet, die ich nicht rot sah.
- `verifizierbar`: ja — `make test-store` auf einer Kopie mit der Mutation.
- `klasse`: Beleg trägt seinen Satz (nur in Variante) nicht

---

## Zu den Schwerpunkten

- **(a) Auswertung.** Eine Stelle: `blockBuilder.build` ruft `model.EvaluateRoute(routes, b.columns, row)` mit den Roh-Zeilenwerten (vor `BuildRowImage`/Transformation), Neu-Bild, ohne Bezug zum Block; die Domänen-Funktion wird gerufen, nicht nachgebaut (`git grep -n 'EvaluateRoute('` im Nicht-Test-Code: 3 Treffer, Definition, Mapper, `build`). Der Vertragstest `TestBackfillAndWALRouteTargetsAreEqual` trägt seinen Satz: echter `mapper.Assembler` und echter `BackfillTableService`, Erwartung als feste Tabelle (`tc.want`), nicht als Aufruf der Auswertung; die Eingabe ist gebunden (rot bei Mutation von Regelliste, Spaltenliste, Anwendbarkeitsprüfung, WAL-Seite — siehe (g)). Die Eingabe „Quellwert vor Transformation“ ist durch den Fall `rename_column an der Bedingungsspalte` im Unit- und im Paritätstest getragen; eine Mutation, die die Regel gegen das transformierte Bild auswertet, habe ich nicht gefahren. Kein Befund.
- **(b) Fail-closed.** Reihenfolge im Code wie `ADR-0139` Festlegung 1: Öffnen des Snapshots → Transformationsstand lesen/prüfen → Routing-Stand lesen/prüfen (beide vor `NextBlock` und vor `Writer.Begin`) → je Block Mengenvergleich (`sameSet`) → Vergleich vor dem Commit. Fehlertext der Anwendbarkeit trägt Regelname und Spalte, `classifyError` bildet `ErrRoutingColumnMissing` auf `schema`, `ErrRoutingStateChanged` auf `configuration` ab; `schema` geht `configuration` voraus (Test `TestExecuteInapplicabilityPrecedesStateChange`, Mutation rot). Die Grenze „gesetzt und wieder entfernt zwischen zwei Lesungen unsichtbar“ steht im Code-Godoc von `copyBlocks` und in [`LH-FA-CAP-009.a`](../../spec/pflichtenheft.md) mit gleichem Inhalt. **Abweichung:** der nicht lesbare Stand (F-1). Der leere Run (kein Block) prüft den Routing-Stand nicht erneut — wie der Transformationsstand; die Spec nennt den „Zustand unmittelbar vor dem Commit“, der bei einer leeren Tabelle nicht entsteht; kein Befund.
- **(c) Bestandsverhalten.** Ohne Routing-Regeln läuft der Run unverändert: der Fall „keine Regel“ liefert leere Ziele, die bestehenden Tests sind bis auf `Routing` in den Rigs unverändert (`service_test.go`-Diff: nur neue Zeilen). Die Lesekosten (Blockzahl plus zwei Lesungen je Regelstand) stehen im Godoc als Zählung des Codes (`TestExecuteBuildsRouteTargetsFromTheRuleSet` prüft fünf Lesungen bei drei Blöcken); im Plan §6 ist die Verdopplung der Last gegenüber der Messung des Architect-Verdikts als *hergeleitet*, die Zahlen als *übernommen* gekennzeichnet — korrekt gekennzeichnet, Ausgang für die Closure offen. Backfill-Changes erreichen keinen Live-Weg: der Diff berührt keinen Zustell-/Notify-Code.
- **(d) Tests.** Store-Test in `internal/bootstrap`: Abweichung ehrlich benannt (Plan §3) und gerechtfertigt; der Test läuft real (`make test-store` Exit 0, `internal/bootstrap` ok; die Mutationen (g) färben ihn rot, also wird er nicht übersprungen). Lücke der Gesamtkette: F-4. Änderung bestehender Tests: Pflicht-Port `Routing` in vier Rigs (`service_test.go` zweifach, `backfill_endtoend_test.go`, `backfill_image_parity_test.go`, `administration_roles_internal_test.go`) — nur Ergänzungen, keine Erwartung gelockert. `make test-replication`: Exit 0, `postgressnapshot` ok (PostgreSQL 18).
- **(e) Fremddatei.** Die Änderung an [`slice-routing-betriebsdoku`](../plan/planning/open/slice-routing-betriebsdoku.md) ist als Aufschub-Adresse nach `AGENTS.md` §3.13 gerechtfertigt: der Gegenstand steht als committeter Text im Plan der Adresse (Kernbegriffe `Routing-Regelstand`, `remove_route`, `schema`, Regelname und Spalte per `git grep` gefunden). Zeilen-Lokatoren aufgeschlagen: `benutzerhandbuch.md:756` und `:1876` stimmen; die Abschnittsnummer „§5“ nicht (F-2). Das Handbuch selbst ist unberührt, damit entfällt die Versionshistorie-Pflicht, und es entsteht keine neue Betreiber-Oberfläche (Run-Verhalten, kein neuer Parameter).
- **(f) Kommentare.** `make kommentar-kennungen DIFF=177bbac5~1`: Exit 0, kein Kandidat (Probe, kein Beleg). Gelesen: Godocs von `copyBlocks`, `classifyError`, `build`, `checkRoutesApplicable` tragen Zusagen, die der Code an der Stelle erfüllt (Pfade nachgefahren); keine Chronik, kein Konjunktiv über Verworfenes. Test-Godocs behaupten Mutationen, die ich bis auf F-5 rot gesehen habe (siehe (g)). Rest: F-3.
- **(g) Mutationen (selbst gefahren).** 18 Läufe, alle rot wie behauptet: M1 `nil` statt Regelliste an `EvaluateRoute` (`TestExecuteBuildsRouteTargetsFromTheRuleSet`); M2 `nil` statt Zeile (dito); M3 Anwendbarkeitsprüfung gegen `nil` statt Snapshot-Spalten (anwendbare Regel scheitert; mehrere Tests); M4 Vergleich je Block entfernt (sechs Fälle „zwischen Block“ rot); M5 Vergleich vor dem Commit entfernt (die zwei Fälle „vor dem Commit“ rot); M6 beide Vergleiche entfernt bei bleibendem Transformationsvergleich (`TestExecuteOnlyRoutingStateChangesEndsRunAsConfiguration` rot); M7 Anwendbarkeitsprüfung entfernt (`…AsSchema`, `…PrecedesStateChange` rot); M8 `ErrRoutingColumnMissing` → `configuration` abgebildet (rot); M9 Regelstand aller Tabellen statt der Tabelle des Runs (zwei Fälle rot); M10, M11, M12 am Paritätstest (`nil` statt Regelliste, Prüfung gegen `nil`, verschobene Spaltenliste: rot), M13 WAL-Seite (`mapper.go`, `nil` statt `binding.Routes`: Paritätstest rot); M14 Längenvergleich statt Mengenvergleich (Fälle „Ordnung und Doppelung“, „Regel ersetzt“ rot); M15 Lesefehler verworfen (`TestExecuteRoutingReadFailureEndsRun` Lesung 1 bis 5 rot); M16 Anwendbarkeit gegen leere Regelliste (rot); S1 Store: `$10` → `NULL` im Insert (rot, als Variante, F-5); S2 Store: `'applied'` → `'pending'` in `SelectAppliedRoutingRequests` (Store-Test und `TestAdministrationPathRunsUnderLeastPrivilegeLogins` rot). **Gefahren: 18. Die übrigen der laut Plan/Godocs etwa 20 genannten Mutationen: übernommen.**
- **(h) Sensoren.** `make kommentar-kennungen DIFF=177bbac5~1`: Exit 0. `make fmt-check`: 316 Dateien formatiert, Exit 0. `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-routing-backfill-pfad.md`: 23 Zeilen stimmen, Exit 0 (Zahlen und Stände; Suchraum und Muster gelesen: der Lauf hat `sameSet` nicht erfasst, F-3; Nicht-Gefundenes ist je Träger benannt). `make test`: Exit 0. `make test-store`: Exit 0. `make test-replication`: Exit 0.
- **(i) Umfang.** +1019/−44 über 14 Dateien, davon rund 800 Zeilen Tests (abgeleitet aus `--stat`, nicht am Diff neu gezählt); Produktionscode: `service.go` (`git diff --numstat`: +85/−25), ein Sentinel, eine Verdrahtungszeile. Drei Commits nach Gegenstand. Review ohne Aufteilung tragfähig.

## Frage an den Architect (zu F-1)

[`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md) Festlegung 1 und
[`LH-FA-CAP-009.a`](../../spec/pflichtenheft.md) sagen: ein nicht lesbarer Routing-Regelstand ist eine Abweichung und
endet den Run `failed` mit der Klasse `configuration`. Der Bestand für Ausschluss- und Transformationsstand endet bei
einem Lesefehler mit der Klasse der Ursache (`storage` bzw. was der Port meldet). Präzise Frage: Gilt für den
Routing-Stand (1) der Wortlaut — dann wird der Lesefehler dieses Standes auf `configuration` abgebildet und die Spec-Zeile
bleibt (der Bestand bei den anderen Ständen weicht dann ab und wäre eigens zu führen), oder (2) die bestehende
Behandlung „Klasse der Ursache“ — dann ist die Aussage „auch ein nicht lesbarer Stand … `configuration`“ in der Spec und in
`ADR-0139` zu berichtigen (Spec-Nachzug und, für die `Accepted` ADR, eine Folge-ADR nach `AGENTS.md` §3.5)? Die
Implementierung darf nicht stillschweigend die Auslegung (2) tragen, solange der Text (1) sagt.

## Negativbefunde

- geprüft, ohne Befund: `internal/domain/errors/errors.go` (Sentinel `ErrRoutingStateChanged`, Text und Klasse nach Muster)
- geprüft, ohne Befund: `internal/bootstrap/wiring.go` (eine Zeile `Routing: activation`; einzige Konstruktionsstelle im Nicht-Test-Code, `git grep -n 'backfill\.Ports{'`)
- geprüft, ohne Befund: `internal/application/usecase/backfill/routing_test.go`, `service_test.go` (alle Godoc-Mutationen rot, kein Test gelockert)
- geprüft, ohne Befund: `internal/bootstrap/backfill_route_parity_test.go` (Eingabe-Bindung, feste Erwartung)
- geprüft, ohne Befund: `internal/adapters/driven/postgressnapshot/snapshot_test.go` (`checkRouteParity`; läuft unter `make test-replication`, PostgreSQL 18; Version 17 nicht gefahren)
- geprüft, ohne Befund: `internal/bootstrap/backfill_routing_store_internal_test.go` (läuft real, Mutationen rot; Lücke F-4, Variante F-5)
- geprüft, ohne Befund: `internal/bootstrap/backfill_endtoend_test.go`, `backfill_image_parity_test.go`, `administration_roles_internal_test.go` (nur `Routing` im Port-Satz ergänzt)
- geprüft, ohne Befund: `docs/plan/planning/done/slice-routing-antragsweg.md` (reine Linkkorrektur auf `in-progress/`)
- Docker-only-Disziplin: im Diff kein `sed -i`, keine Host-Interpreter, keine Umleitung auf Repo-Dateien erkennbar; Traceability der Commit-Messages (`ADR-0137`/`LH-FA-CFG-008`, keine `SPEC-*`-Kennung im Betreff) in Ordnung.
- nicht geprüft: `make gates` vollständig, `make test-integration`, `make coverage-gate`, Replikationstest an PostgreSQL 17, Lesekosten-Messung (im Plan *hergeleitet*) — Gate-/Messläufe, Gegenstand des Verifiers.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 2 |
| INFO | 2 |

## Verdikt

**Merge-blockierend: ja — wegen F-1 (HIGH), nicht wegen eines Korrektheitsfehlers im Normalpfad.** Auswertung (eine Stelle, Quellwerte, Parität), Anwendbarkeit (`schema`, einmal, vor Schreibtransaktion), Mengenvergleich je Block und vor dem Commit sind am Code und an 18 selbst gefahrenen Mutationen getragen. Offen: F-1 hängt an einer Architect-Entscheidung (Wortlaut oder Bestand); F-2 und F-3 sind kleine Fixrunden-Punkte (Planner bzw. Implementer).

**DoD-Checkbox:** nicht nachgezogen — es gibt eine Fixrunde (F-1); die Zeile „Review durchgeführt“ wird nach der Fixrunde regulär gesetzt.

**Übergabe:** F-1 als Architect-Frage (Artefakt: diese Frage in diesem Report; Träger: Planner), danach an den Implementer; F-2 an den Planner (Plan §3 und Übergabe-Block in [`slice-routing-betriebsdoku`](../plan/planning/open/slice-routing-betriebsdoku.md): „§6“ statt „§5“); F-3 an den Implementer; F-4, F-5 zur Kenntnis.
