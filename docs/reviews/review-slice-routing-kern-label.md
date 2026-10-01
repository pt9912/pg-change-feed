# Review-Report: slice-routing-kern-label — 2026-10-01

**Review-Art:** Code — geprüft gegen Plan, `ADR-0137`, `ADR-0138`, `ADR-0139`, `ADR-0114`, `SPEC-032`
(sowie `SPEC-001`/`SPEC-002`/`SPEC-008`) und `AGENTS.md` Hard Rules (Modul 10). Kein DoD-Abgleich
(Verifier).

**Gegenstand:** Slice `slice-routing-kern-label` (`welle-routing`), Diff-Range `41fbd0ca..HEAD`:
`88476c94` (Domäne, `Assembler`, `classifyRunError`), `b15b4416` (Schema, Store, Alt-Tag-Skript,
Erzeugnisse `plan.yaml`/`down.sql`), `3861d679` (Plan-Nachzug). 25 Dateien, +1659/−67
(`git diff --stat 41fbd0ca HEAD`).

**Skill:** `.harness/skills/reviewer.md` @ `c5207cc1`. **Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-01.

**Ablage:** Der Reviewer-Lauf hat diesen Report selbst geschrieben (Write-Werkzeug). Alle Mutationen
liefen an einer `git archive HEAD`-Kopie im Scratchpad (Mutation per `sed … > Datei`, nie `sed -i`,
nie eine Umleitung auf eine Repo-Datei), Tests im gepinnten Toolchain-Image mit `--network none` und
`-race` (dieselbe Aufrufform wie `make test`, Quelle: die Kopie). Es gab keine verweigerte Aktion.

**Eingangs-Kontext:**

- Slice-Plan `slice-routing-kern-label` (§1–§3 mit Suchlauf-Feld, §6 Risiken), `welle-routing` §5 (V3)
- `ADR-0137`, `ADR-0138` (inkl. „Offen — V3“), `ADR-0139`, `ADR-0114`, `ADR-0111` (Präzedenz `origin`)
- `SPEC-032` (Regelform, Bildbasis, Abwesenheit, Auswertung, Anwendbarkeit), `SPEC-019` R3, `SPEC-001`/`SPEC-002`/`SPEC-008`
- `LH-FA-CFG-008`, `LH-FA-DAT-006`, `LH-FA-ADM-003`
- `AGENTS.md` §3.1, §3.7, §3.9, §3.12, §3.13

---

## Findings

### F-1 — Kommentar widerspricht dem Nachbarn: „origin steht als letzte Spalte“ neben „route_target steht danach als letzte Spalte“

- `kategorie`: MEDIUM
- `quelle`: Maintainability / Nachzug widerspricht dem Nachbarn im selben Träger (`AGENTS.md` §3.7)
- `pfad`: `internal/adapters/driven/postgresstorage/queries/queries.go:44-48`; `internal/adapters/driven/postgresstorage/sqlexec/translate_test.go:306`
- `befund`: Der Godoc von `SelectChanges` sagt „`origin` steht als letzte Spalte“ und im nächsten Satz „`route_target` steht danach als letzte Spalte“ — zwei „letzte“ im selben Block. Der Test-Kommentar `translate_test.go:306` („Die letzte Spalte der Projektion trägt die Herkunft“) in einer vom Diff geänderten Datei ist seit diesem Diff falsch; `sqlviews_test.go:278` hat die Berichtigung („vorletzte“) bekommen, diese beiden Stellen nicht.
- `verifizierbar`: ja — `git grep -n "letzte Spalte" -- internal` (Lese-Handlung am Treffer; kein Gate).
- `klasse`: Nachzug widerspricht dem Nachbarn im selben Träger

### F-2 — Alt-Tag-Skript: Kopfkommentar und benannte Mutationen beschreiben einen Upgrade-Beleg, den der Lauf für die Transformations-Objekte nicht mehr trägt

- `kategorie`: MEDIUM
- `quelle`: `ADR-0114` Entscheidung 7 (Alt-Tag-Lauf); Skill-Klasse „Beleg trägt seinen Satz nicht“
- `pfad`: `tools/harness/run-schema-rollout-guard-test.sh:36-58` (Kopf Lauf 5), `harness/targets/schema-rollout.md:183`
- `befund`: Die Änderung entfernt die Vorbedingungen zu `rule_name`/`rule_spec`, den zwei Funktionen und den Antragsarten (zu Recht: `git show v0.4.0:tools/schema/schema.yaml` trägt `rule_name`, `rule_spec` und die Transformations-Funktionen bereits, die alten Vorbedingungen mussten rot sein). Der Kopf sagt aber weiter „nach dem Upgrade trägt die Tabelle die zwei nullable Spalten … eine vor dem Upgrade geschriebene Antragszeile trägt dort NULL, die sieben Werte … stehen“ und führt die „rot färbenden Eingabeseiten-Mutationen“ in `nacharbeit-administration.sql`/`schema.yaml` (`rule_spec`, GRANT-/REVOKE-/CHECK-Zeilen der Transformations-Funktionen) weiter. Der Tag hat diese Objekte schon, das Upgrade ergänzt sie nicht mehr; ob jene Mutationen den Lauf 5 noch rot färben, ist nach der Änderung nicht nachgefahren (der Implementer hat nur die Vorbedingung begründet, nicht die Mutationsliste). Die Belegkraft des Laufs für die Transformations-Objekte sinkt von „Upgrade ergänzt“ auf „Upgrade erhält“, der Text sagt es nicht.
- `verifizierbar`: ja — die genannten Mutationen im Alt-Tag-Lauf fahren (`bash tools/harness/run-schema-rollout-guard-test.sh`, Docker, nicht nachgefahren, s. Negativbefunde).
- `klasse`: Beleg trägt seinen Satz nicht

### F-3 — V3: Die Unit-Messung ist ehrlich benannt, aber ihr Ergebnis löst die Bedingung „Architect-Verdikt vor dem Negative-Beleg“ aus

- `kategorie`: MEDIUM
- `quelle`: `ADR-0138` „Offen, nicht entschieden — V3“; `welle-routing` §5 V3; `LH-FA-CFG-008`
- `pfad`: `internal/adapters/driving/replication/mapper/routing_test.go:243-290`; `docs/plan/planning/in-progress/slice-routing-kern-label.md` §6 (V3-Absatz)
- `befund`: Nachgefahren: `TestRoutingNotApplicableReachabilityThroughRelationCheck` trägt, was Plan §6 sagt — (a) bei bekannter Spaltenform endet die Relation ohne die Bedingungsspalte in `ErrIncompatibleSchemaChange`, bevor eine Change assembliert wird; (b) nur bei Erstaktivierung ohne Spaltenform entsteht `ErrRoutingNotApplicable`. Die Aussage ist als Unit-Ebene (Instanz `Assembler` mit `fakeSchemaStore`) und die Systemebene als ungemessen (`slice-routing-e2e`) benannt, das ist ehrlich. Offen bleibt, dass auf dem naheliegenden Weg (a) der in `ADR-0138` und `SPEC-032` zugesagte Abhilfe-Weg „`cdc.remove_route` und Neustart“ nicht greift: die Relation bleibt gegen die bekannte Spaltenform inkompatibel, auch ohne Regel. Der Slice-Plan §6 trägt den Ausgang „bei der Closure einzutragen“ und den Hinweis auf das Verdikt, aber noch kein Verdikt.
- `verifizierbar`: ja — `go test` des genannten Tests (grün in meiner Kopie); die Systemebene nur durch `slice-routing-e2e`.
- `klasse`: Offene Vorab-Bedingung ohne Verdikt

### F-4 — Handbuch nennt `origin` die letzte Spalte der View (ab diesem Commit falsch), Aufschub-Adresse trägt den Gegenstand

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.13 (Träger außerhalb des Diffs)
- `pfad`: `docs/user/benutzerhandbuch.md:516`, `:1462`
- `befund`: „`origin` ist die letzte Spalte der View“ (516) und die Feldliste der dreizehn Felder (1462) beschreiben `cdc.changes` ohne `route_target`. Der Aufschub nach `slice-routing-betriebsdoku` ist benannt, und die Probe der Adresse trifft: `slice-routing-betriebsdoku` Zeile 80 trägt den Gegenstand („Änderungen lesen“, Lesesemantik von `NULL`). Der Satz „letzte Spalte“ selbst steht dort nicht wörtlich; Adresse und Gegenstand decken ihn nicht ausdrücklich.
- `verifizierbar`: ja — `git grep -n "letzte Spalte" -- docs/user`.
- `klasse`: Träger-Nachzug außerhalb des Diffs

### F-5 — Die Vorbedingung des Alt-Tag-Laufs („der Tag trägt `route_target` noch nicht“) kippt beim nächsten `v*`-Tag

- `kategorie`: INFO
- `quelle`: `ADR-0114` Entscheidung 7; Maintainability
- `pfad`: `tools/harness/run-schema-rollout-guard-test.sh:279`
- `befund`: Die Vorbedingung hängt am „jüngsten `v*`-Tag“; sobald ein Release diesen Slice enthält, bricht Lauf 5 an derselben Stelle wie er jetzt an den alten Vorbedingungen gebrochen wäre. Das ist die zweite Wiederholung des Musters (Tag holt die Vorbedingung ein); der Tag-Wechsel ist ein bekanntes, hier nicht abgefangenes Ereignis.
- `verifizierbar`: nein — tritt erst mit dem nächsten Release auf.
- `klasse`: Vorbedingung an beweglichem Tag

### F-6 — Auswertungslast ungemessen, Rollout-Dauer nicht im Diff

- `kategorie`: INFO
- `quelle`: `AGENTS.md` §3.12
- `pfad`: `docs/plan/planning/in-progress/slice-routing-kern-label.md` §6 „Last der Auswertung“
- `befund`: Der Plan nennt die Last als „erwartet klein, nicht gemessen“; das ist korrekt gekennzeichnet. Der Nullpfad (Tabelle ohne Routen) ist am Code nachgefahren: `checkRoutes` und `EvaluateRoute` iterieren über eine leere Liste, `WithRouteTarget("")` macht keinen Allokations- oder Validierungsschritt über einen Wertkopf hinaus — kein Verhaltensbruch. Ein „Rollout-Fenster 7 s“ steht in keiner hinzugefügten Zeile des Diffs; ich habe es nicht gefunden und nicht beurteilt.
- `verifizierbar`: nein — Benchmark fehlt bewusst.
- `klasse`: Last ungemessen, korrekt gekennzeichnet

## Zu den Schwerpunkten

- **(a) Domäne.** Alphabet gegen `SPEC-032` Zielname: stimmt (erstes Zeichen `a-z0-9`, Folgezeichen plus `_`/`-`, 1 bis 63 Byte; Nicht-ASCII fällt durch). `order >= 1`, `when.column` nicht leer und ohne U+0000, leere `equals` zulässig: stimmt. `EvaluateRoute`: kleinste `order`, bei Gleichstand die zuerst genannte (R2 hält sie eindeutig), kein Treffer leer; abwesender Wert (`nil`, Index außerhalb) ist Nicht-Treffer. `CheckApplicable` hängt an der Spaltenmenge, nie am Wert (Spec: Anwendbarkeit). `CheckNotExcluded` deckt die Zeile „R3, `when.column` ist ausgeschlossen“ (Klartext „Spalte ist ausgeschlossen“) in Richtung Regel → Ausschluss; die Gegenrichtung („Spalte trägt eine Routing-Bedingung“) ist `antragsweg`. `make a-check`: 0 Befunde (Domäne importiert nichts aus anderen Schichten).
- **(b) Mapper.** Auswertung gegen `event.New`/`event.Old` (Quellwerte) nach dem Bau von Neu- und Alt-Bild, aber unabhängig von ihnen; `checkRoutes` steht vor `a.open.sequence++`; `AddBinding` übernimmt `Routes`; Listen werden bei jedem Nachtrag neu aufgebaut, der Leser hält seinen Schnappschuss aus `lookupBinding`. `-race`-Test vorhanden.
- **(c) Persist-before-ACK.** Nachgefahren: `Consume` gibt bei `ErrRoutingNotApplicable` `nil, err` zurück (`mapper.go:221-224`), die Transaktion bleibt offen und erreicht `Commit`/`CaptureCommand` nie; `receive.Stream.process` (`receive.go:615-618`) gibt den Fehler zurück, der Lauf endet vor `Capture` und vor jeder Bestätigung. `classifyRunError` bildet auf `schema` ab. Kein Befund.
- **(d) Schema.** `route_target` in `schema.yaml` und `postgresstorage/schema.sql` (nullable, kein DEFAULT, kein CHECK), View ohne `COALESCE` als letzte Spalte in Abfrage und `columns:`-Signatur. `plan.yaml`-Token-Diff (Kommas als Trenner): nur geänderte Operations-/Gruppen-IDs, Hash, Länge, die neue Spalte und die neue View-Spalte; kein Host-Name und kein Fremdinhalt. `down.sql`: nur der `artifactHash` im Kopf. Die Rechte der drei Rollen auf der View sind durch den DB-Test und Lauf 5 belegt (letzterer nicht von mir gefahren).
- **(e) Store.** Beide Insert-Wege (`InsertChange`, `InsertBackfillChange` mit `$10`), `ChangeRow.RouteTarget *string`, `SelectChanges` ohne `COALESCE`, `ReadChanges` scannt die Spalte, Alphabet in beide Richtungen der Zeilen-Übersetzung. Ein außerhalb des Alphabets gespeicherter Wert lässt `ReadChanges` mit Domänenfehler enden (wie `origin`); kein Befund.
- **(f) Alt-Tag-Skript.** Die entfernten Vorbedingungen waren am Tag `v0.4.0` nachweislich falsch (`git show v0.4.0:tools/schema/schema.yaml` trägt `rule_name`, `rule_spec`, `set_transformation`); die Änderung ist gerechtfertigt. Zugesetzte Prüfungen (Rechte auf `cdc.changes`, `route_target IS NULL` der Alt-Zeile) stärken den Lauf für diesen Slice. Der Lauf selbst: nicht nachgefahren. Belegkraft-Verlust und Text: F-2, Wiederholung: F-5.
- **(g) V3.** F-3. Gehört jetzt zu einem Architect-Verdikt (Frage unten).
- **(h) Mutationen (Kopie, `go test -race`, Ziel-Pakete).** Nachgefahren, jeweils rot: M1 DELETE liest `event.New` statt `event.Old` (`TestConsumeRoutesByBildbasisAndOrder`); M2 `values[i] != nil` entfernt (Nullzeiger-Panic, `TestEvaluateRoute`); M3 `binding.Routes = existing.Routes` entfernt (`TestAssemblerRoutesSurviveAddBindingAndSchemaVersionBump`); M4 `Lock()`/`Unlock()` in `SetRoute` entfernt (`TestAssemblerRoutesAreRaceFree`); M5 `checkRoutes` hinter das Vorrücken der Sequenz verschoben (`TestConsumeRoutingRuleOnMissingColumnIsNotApplicable`; diese Mutation steht nicht im Kommentar, die dort benannte war „Aufruf entfernen“); M6 `order`-Vergleich gedreht (`TestEvaluateRoute`, `TestEvaluateRouteIsIndependentOfListOrder`); M7 `order < 1` zu `order < 0` (`TestNewRouteRuleInvariants`); M8 Eingabeseite: Bildbasis für UPDATE auf `event.Old` (vier Tests rot); D1 `SelectChanges` liefert `NULL::text AS route_target` (`make test-store`-Runner auf der Kopie: `TestChangesViewCarriesRouteTargetLikeReadChanges` und `TestPersistAndReadCarryRouteTarget` rot, sonst grün). Nicht nachgefahren: die übrigen 13 der behaupteten 22 (u. a. D2–D5, Alphabet-Mutationen, `classifyRunError`-Mutation); die Aussage „22 Mutationen“ steht als **übernommen**, nicht als gemessen. Eine Mutation (`EvaluateRoute(…, nil)` statt `routeValues`) endete als Build-Fehler (unbenutzte Variable) und zählt nicht.
- **(i) Formale Sensoren.** `make fmt-check`: 303 Dateien formatiert (Exit 0). `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-routing-kern-label.md`: 12 Zeilen stimmen. `make kommentar-kennungen DIFF=910629bb`: Exit 0, kein Kandidat (Probe, kein Beleg; die Form-Lesung ergab F-1). `make a-check`: 0 Befunde.
- **(j) Umfang.** +1659/−67 ist tragfähig: die Test-Dateien tragen 1219 Zeilen der `--stat`-Ausgabe (abgeleitet: Summe der Einzelzahlen, nicht am Diff neu gezählt), `plan.yaml`/`down.sql` sind Erzeugnisse, Produktionscode und Schema sind der kleinere Rest; drei Commits, getrennt nach Schicht. Ein Review des Diffs ist ohne Aufteilung möglich. Die im Plan §4 benannte Rückführung (Persistenz als eigener Slice) ist nicht eingetreten.

## Frage an den Architect (zu F-3)

`ADR-0138` „Offen — V3“ und `welle-routing` §5 V3 verlangen bei einem unerreichbaren Fall ein Architect-Verdikt vor dem Negative-Beleg. Die Unit-Messung zeigt: Auf dem Weg „Spalte entfernt“ erreicht der Erfassungspfad `ErrRoutingNotApplicable` nicht, weil `observeRelation` die inkompatible Relation vorher als `ErrIncompatibleSchemaChange` (Klasse `schema`) meldet; `ErrRoutingNotApplicable` ist nur bei Erstaktivierung ohne gespeicherte Spaltenform erzeugbar. Präzise Frage: Soll der in `LH-FA-CFG-008.a` und `SPEC-032` (Anwendbarkeit) zugesagte Abhilfe-Weg („Regel entfernen, dann läuft der Erfassungspfad wieder“) (1) für den Fall `ErrIncompatibleSchemaChange` bei vorhandener Routing-Bedingung auf die entfernte Spalte gelten (dann verlangt das eine Entscheidung, wie die Relation-Prüfung die Regel berücksichtigt), (2) auf den Erstaktivierungsfall beschränkt bleiben und der Negative-Beleg von `slice-routing-e2e` auf diesen Fall (oder auf den Unit-Beleg) zugeschnitten werden, oder (3) der Wortlaut der Spec angepasst werden, weil `ErrRoutingNotApplicable` am laufenden System ohne künstliche Konstruktion nicht entsteht?

## Negativbefunde

- geprüft, ohne Befund: `internal/domain/model/route.go`, `route_test.go`, `change.go` (Alphabet, Invarianten, Auswertung, Bildbasis, Abwesenheit, Mutationen M2/M6/M7 rot)
- geprüft, ohne Befund: `internal/domain/errors/errors.go` (vier Sentinels, Texte wie `SPEC-019`)
- geprüft, ohne Befund: `internal/adapters/driving/replication/mapper/mapper.go`, `routing_test.go` (Schnappschuss, `AddBinding`, Race, Reihenfolge der Prüfung; M1/M3/M4/M5/M8 rot)
- geprüft, ohne Befund: `internal/bootstrap/wiring.go`, `heartbeat_internal_test.go` (`classifyRunError`)
- geprüft, ohne Befund: `internal/adapters/driven/postgresstorage/` (Schema, Queries, Store, Backfill-Writer, `sqlexec`; D1 rot)
- geprüft, ohne Befund: `tools/schema/schema.yaml`, `plan.yaml`, `down.sql` (View-Signatur, Erzeugnisse ohne Fremdinhalt)
- nicht geprüft: `make test-store`/`make test-replication` vollständig als Gate, `make schema-rollout` zweimal, Alt-Tag-Lauf (Lauf 5), `make gates`, Benchmark der Auswertung — Docker-/DB-Läufe außerhalb dieser Prüfung, Gegenstand des Verifiers
- Docker-only-Disziplin: im Diff kein `sed -i`, keine Host-Interpreter, keine Umleitung in Repo-Dateien erkennbar; Traceability der drei Commit-Messages (`LH-FA-CFG-008`/`ADR-0137`, keine `SPEC-*` im Betreff) in Ordnung.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 3 |
| LOW | 0 |
| INFO | 3 |

## Verdikt

**Merge-blockierend: ja, nach der Closure-Regel des Plans („kein offenes HIGH/MEDIUM“), nicht wegen eines Korrektheitsfehlers.** Der kritische Pfad (Persist-before-ACK bei `ErrRoutingNotApplicable`), die Auswertungssemantik gegen `SPEC-032`, die Schicht-Konformität und die Spalte samt View sind am Code und an Mutationen getragen. Offen: F-1 und F-2 (Implementer, Fixrunde klein), F-3 (Architect-Verdikt, Träger: Planner/Architect; blockiert den Negative-Beleg von `slice-routing-e2e`, nicht die Auslieferung des Kerns).

**DoD-Checkbox:** nicht nachgezogen — es gibt eine Fixrunde (F-1, F-2); die Zeile „Review durchgeführt“ wird beim Schritt 21 des Implementer-Workflows gesetzt.

**Übergabe:** F-1, F-2 an den Implementer; F-3 als Architect-Frage über den Planner; F-4 an den Planner (Adresse `slice-routing-betriebsdoku` um den Satz „letzte Spalte“ ergänzen); F-5, F-6 zur Kenntnis.
