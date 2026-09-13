# Verifikationsbericht: slice-045 — 2026-09-13

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-045` §1/§2 DoD/§3 Plan-Nachzug/§4/§6/§8), `welle-13` und
`ADR-0046`/den Architect-Verdikt
`architect-verdict-retention-loeschausfuehrung.md`, nicht gegen Diff
(Reviewer-Aufgabe, bereits abgeschlossen: `review-slice-045.md`, vollständig
gelesen) und nicht gegen realen Bedarf (Validator, hier nicht ausgelöst —
kein MVP-Meilenstein-Slice).

**Frischer Kontext:** Diese Prüfung liest den vollständigen Slice-Plan, das
Lastenheft (`LH-FA-RET-005` im Wortlaut), den Architect-Verdikt, den
Review-Report und den tatsächlichen Code/Diff selbst — keine Behauptung aus
dem Implementer- oder Reviewer-Bericht wird ungeprüft übernommen. Jeder unten
genannte Sensor-Lauf (`make gates`, `make test-integration`, `make
test-store`) wurde in dieser Sitzung **selbst** ausgeführt.

**Gegenstand:** `docs/plan/planning/in-progress/slice-045-blockierende-consumer-sichtbarkeit.md`
zum Stand `HEAD = dfe02b0`. Commits (chronologisch): `2a0e4fa`
(`Verantwortlich` gesetzt), `667222d` (`open→next`), `5d36fea`
(`next→in-progress`), `134bc0a` (Implementierung: View, Grant, Tests, Doku,
Plan-Nachzug), `dfe02b0` (Review-Report, 0 HIGH/MEDIUM/LOW, 1 INFO). Sequenz
selbst geprüft (`git log --oneline 5d36fea^..dfe02b0`): Move → Implementierung
→ Review — keine Rolle springt rückwärts ohne Übergabe-Artefakt.

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | Neue SQL-View zeigt je Quelle real den am weitesten zurückliegenden unbestätigten Consumer, real gegen ≥2 Consumer getestet (einer blockiert, einer nicht) | **erfüllt** | `tools/schema/schema.yaml` selbst gelesen: `retention_blockers` liest `cdc.consumer_position INNER JOIN cdc.consumer`, `DISTINCT ON (cp.source_id) ORDER BY cp.source_id, cp.acknowledged_position ASC, c.consumer_id ASC` — liefert real höchstens eine Zeile je Quelle, die mit der kleinsten `acknowledged_position` (genau die Position, die `RetentionPolicy.AllowsDeletion` als bindende Untergrenze behandelt). `TestMVPRetentionBlockersViewShowsFurthestBehindConsumer` (`test/integration/integration_test.go`) selbst gelesen und real ausgeführt (siehe Punkt 3): registriert zwei Consumer direkt über `ConsumerStatePort`, bestätigt beide auf unterschiedliche reale Positionen, prüft `len(blockers) == 1`, `consumerID == behindConsumer`, `backlog > 0` und explizit, dass der weiter bestätigende Consumer **nicht** erscheint — reale Zwei-Consumer-Unterscheidung, kein Mock |
| 2 | `LH-FA-RET-005` real erfüllt | **erfüllt — mit Präzisierung, siehe §2 unten** | Lastenheft selbst gelesen (`spec/lastenheft.md:719-730`): Happy Path „gegeben `c1` blockiert, ist erkennbar welcher Consumer blockiert und bis zu welcher Position" — durch `cdc.retention_blockers` + o.g. Test real belegt. Boundary „mehrere blockierende Consumer … alle einzeln erkennbar" wird **nicht** von der neuen View allein getragen (sie liefert bewusst nur eine Zeile je Quelle), sondern von der bereits bestehenden `cdc.consumer_status`/`cdc_consumer_lag`-Kette — eigenständig nachgeprüft, siehe §2 |
| 3 | `make gates` grün, `make test-integration` grün; `make test-store` grün | **erfüllt, selbst reproduziert** | `make gates` selbst ausgeführt gegen `HEAD = dfe02b0`: `baseline-verify` (v6.5.0, 54 Dateien) OK, `d-check` Struktur (358 Dateien, 0 Befunde), `d-check` Commits (`HEAD~5..HEAD`, 0 Befunde), `commit-traceability.sh` (5 Commits, Betreffs ohne Struktur-ID) OK, `a-check` (0 Befunde). `make test-integration` selbst ausgeführt: alle Go-Testfälle inkl. `TestMVPRetentionBlockersViewShowsFurthestBehindConsumer` `PASS`, anschließender Compose-Rundlauf inkl. des bestehenden Retention-Belegs (`RetentionOld` real entfernt nach Freigabe, `RetentionYoung` blieb erhalten) ebenfalls grün — kein Interferenz-Effekt durch den neuen Testfall (Cleanup-Fix aus Plan-Nachzug Punkt 5 wirkt). `make test-store` selbst ausgeführt: alle Pakete `ok`, insbesondere `postgresstorage` (4.12s, enthält `TestCdcReaderRoleReadsViewsNotBaseTables` mit der neuen `cdc.retention_blockers`-Zeile) |
| 4 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **inhaltlich erfüllt, Formular-Diskrepanz — siehe Finding V-1** | `review-slice-045.md` vollständig gelesen: 0 HIGH/MEDIUM/LOW, 1 INFO (F-1), Verdikt „nicht merge-blockierend". Die Bedingung ist damit tatsächlich erfüllt. Checkbox in §2 (Zeile 91) steht aber weiterhin auf `- [ ]` — der Review-Commit `dfe02b0` ändert ausschließlich `docs/reviews/review-slice-045.md` (`git show --stat dfe02b0`), die Plan-Datei nicht |
| 5 | Doku-Update `docs/user/benutzerhandbuch.md` | **erfüllt** | Abschnitt „Blockierende Consumer erkennen" (Zeile 362ff.) selbst gelesen: korrekt unter „Aufbewahrung (Retention)" platziert (Überschriften-Reihenfolge geprüft), SQL-Beispiel stimmt mit der realen View-Signatur überein, Abwesenheits-Lesart korrekt dokumentiert; `cdc_reader`-Zeile der Rollen-Tabelle (Zeile 77) nennt `cdc.retention_blockers` |
| 6 | Closure-Notiz mit Steering-Loop-Lerneintrag | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit (Modul 8); §7 ist noch die Bedienhinweis-Vorlage. Kein Verifikations-Gegenstand dieser Prüfung |
| 7 | Reconciliation-Register fortgeschrieben, falls einschlägig | **entfällt strukturell** | Kein `docs/plan/planning/reconciliation.md` — `harness/conventions.md` §Modus-Deklaration führt ausschließlich `*`/`PGC` im Modus Greenfield |
| 8 | Beobachtungs-Register fortgeschrieben | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit. Zur Einordnung selbst geprüft: `BEO-PGC/retention-keine-loeschausfuehrung` (0×, `evidence/` real leer) und `BEO-PGC/d-migrate-nacharbeit` (5×, `evidence/` real: `slice-006/010/015/016/036.md`) — Zählerstände decken sich exakt mit der Plan-Aussage in §8, keiner erreicht mit diesem Slice 3× neu |
| 9 | Jedes Risiko aus §6 trägt einen Ausgang | **offen — korrekt unbeansprucht, inhaltlich vorbereitet** | Planner-Closure-Arbeit; beide Risiken tragen noch `<bei Closure einzutragen>`. Reale Grundlage für beide selbst geprüft — siehe §3 unten |
| 10 | Drei Paarungen (Anker · Folge-Slice · Register) | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit, wellenlos hier bzw. bei `welle-13`-Closure fällig, erst **nach** dem `git mv` nach `done/` sinnvoll prüfbar |

## 2. Antwort auf Reviewer-Finding F-1 (LH-FA-RET-005-Boundary über zwei Views)

Der Reviewer hat den Punkt korrekt als „für die Verifier-Rolle" markiert
(DoD-/Spec-Konformitätsfrage, kein Review-Gegenstand). Eigenständig
nachgeprüft, nicht nur der Plan-Nachzug-Argumentation geglaubt:

- **`cdc.consumer_status`** (`tools/schema/schema.yaml:268-289`, unverändert
  in diesem Diff) liefert **eine Zeile je (Consumer, Quelle)-Paar** über
  `LEFT JOIN cdc.consumer_position` — jeder Consumer mit oder ohne bestätigte
  Position erscheint einzeln, mit `acknowledged_position` und
  `latest_commit_position`. Bei **mehreren** gleichzeitig zurückliegenden
  Consumern derselben Quelle zeigt diese View **jeden einzeln** mit seinem
  eigenen Rückstand — unabhängig davon, ob er aktuell die Löschgrenze trägt.
  Das ist strukturell verschieden von `cdc.retention_blockers`
  (`DISTINCT ON`, bewusst genau eine Zeile je Quelle).
- **Real am Code nachvollzogen, dass diese Sichtbarkeit bereits vor diesem
  Slice existierte und getestet ist:** `internal/bootstrap/wiring.go:1193-1230`
  (`Diagnose`) liest `cdc.metrics WHERE metric_name = 'cdc_consumer_lag'`
  (eine Zeile je Consumer) und gibt **jeden Consumer einzeln** mit seinem
  Rückstand aus (`fmt.Printf("    %s: %.0f\n", consumer, *lag)`). Der
  bestehende E2E-Beleg in `tools/harness/run-integration-tests.sh` (Zeilen
  554–608, unverändert in diesem Diff, aus `slice-038`/`slice-039`) zeigt real
  **zwei gleichzeitig existierende Consumer** (`CLI_CONSUMER` mit
  Rückstand ≠ 0, `BACKLOG_CONSUMER` mit Rückstand 0) in derselben
  `diagnose`-Ausgabe — bei diesem `make test-integration`-Lauf selbst erneut
  reproduziert (`run-integration-tests: CLI-Diagnose-Beleg (Normalbetrieb) —
  alle vier Signale … sichtbar`, Ausgabe siehe §1 Punkt 3).
- **Fazit:** Die Boundary „mehrere blockierende Consumer … alle einzeln
  erkennbar" war bereits **vor** `slice-045` durch `cdc.consumer_status`/
  `cdc_consumer_lag` real erfüllt und real getestet — dieselbe Feststellung,
  die der Architect-Verdikt vor der Welle-Eröffnung bereits traf. `slice-045`
  liefert keine Lücken-Schließung für die Boundary, sondern eine **engere,
  zusätzliche** Sicht: welcher Consumer **aktuell die Löschentscheidung**
  bindet (Happy Path mit Bezug zur tatsächlichen `RetentionPolicy`-Grenze,
  nicht nur „irgendein Rückstand"). Beide Aussagen sind wahr und widerspruchsfrei;
  DoD-Punkt 2 ist damit **inhaltlich erfüllt**, aber so formuliert, dass er
  ausschließlich den neuen Test zitiert und die tragende Rolle von
  `consumer_status` nicht sichtbar macht. Das ist kein Substanz-Mangel — die
  Anforderung ist real gedeckt, nachprüfbar mit bestehenden, unveränderten
  Testartefakten — sondern eine Formulierungs-Unschärfe im DoD-Punkt selbst.
  Keine Rückführung nötig; ein Hinweis für die Closure-Notiz (§7 „Was ging
  anders als geplant") ist die angemessene Form, kein Plan-Defekt.

## 3. Risiken aus §6 — reale Grundlage geprüft (Ausgang bleibt Planner-Arbeit)

- **Risiko 1 („Consumer nie bestätigt")**: `tools/schema/schema.yaml` zeigt
  real `FROM cdc.consumer_position cp JOIN cdc.consumer c` — **INNER JOIN**,
  keine `LEFT JOIN`. Ein Consumer ohne jede bestätigte Position gegen eine
  Quelle trägt strukturell keine Zeile in `cdc.consumer_position` und kann
  deshalb in `cdc.retention_blockers` weder fälschlich als „kein Blocker"
  noch mit einem irreführenden Wert erscheinen — das Risiko ist durch das
  Design ausgeschlossen, nicht nur behauptet. Trägt eine Grundlage für den
  Ausgang *entfallen*.
- **Risiko 2 (`POST_EXECUTE_DRIFT` bei d-migrate)**: Selbst reproduziert —
  sowohl `make test-integration` als auch `make test-store` haben in dieser
  Sitzung `schema migrate --execute` real gegen eine leere DB ausgeführt;
  beide Läufe endeten mit `CREATE VIEW`/`GRANT` ohne
  `VIEW_SIGNATURE_UNKNOWN`/`VIEW_SIGNATURE_INCOMPATIBLE`-Blocker. Trägt eine
  Grundlage für den Ausgang *entfallen*.

Beide Zuweisungen sind Planner-Urteil (Modul 5 §Offene Risiken werden bei
Closure aufgelöst); diese Prüfung liefert nur den realen Befund, auf dem das
Urteil aufsetzen kann.

## 4. Plan-vs-Code-Diff

- **View-Name/-Form (`cdc.retention_blockers`):** deckt sich mit
  Plan-Nachzug Punkt 1; Implementer-Entscheidung, wie im Slice-Ziel
  vorgesehen (§1: „Implementer entscheidet den Namen und begründet ihn").
- **`cdc_reader`-Grant-Erweiterung:** deckt sich mit Plan-Nachzug Punkt 4 und
  dem Architect-Verdikt Frage 2 (View-Owner-Muster) — `git diff` gegen
  `tools/schema/nacharbeit-roles.sql` selbst gelesen: nur die
  `cdc_reader`-Grant-Zeile erweitert, kein Grant auf eine Basistabelle.
- **CLI-Erweiterung (`diagnose`) bewusst unterlassen:** Plan §1 formuliert
  dies als *optional*, „falls ohne Bruch möglich" — kein Pflichtpunkt der
  DoD. Plan-Nachzug Punkt 2 begründet den Verzicht mit der
  Slice-Größen-Obergrenze (3 Liefer-Punkte bereits ausgeschöpft) **und**
  einem inhaltlichen Argument (kein Informationsgewinn ggü. direktem
  `CDC_READER_DSN`-Zugriff, der bereits real über die neue View belegt ist —
  siehe §1 Punkt 1/3). Das ordnet sich sauber der §1-Klasse „anderer Vorgang
  / Bestand bleibt bewusst stehen" zu, mit benanntem Folge-Slice-Vorbehalt
  („falls ein Betriebs-Bedarf … entsteht"). **Kein** unbegründeter Verzicht,
  **keine** übersehene DoD-Anforderung — die DoD selbst verlangt die
  CLI-Erweiterung nicht.
- Kein weiterer Abweichungspunkt gefunden: `git diff 5d36fea..134bc0a`
  vollständig gegen die Plan-Tabelle (§3) und den Plan-Nachzug gehalten —
  alle geänderten Dateien (`schema.yaml`, `nacharbeit-roles.sql`,
  `plan.yaml`/`down.sql` d-migrate-generiert, `roles_test.go`,
  `integration_test.go`, `run-integration-tests.sh` [nur Kommentar-Anpassung
  wegen des neuen Go-Testfalls], `benutzerhandbuch.md`) sind im Plan-Nachzug
  benannt oder folgen mechanisch aus einer benannten Änderung.

## 5. Scope-Treue gegen §1 (Ausdrücklich NICHT in diesem Slice)

- **Ausschluss 1 (neue Berechnung des Safe Watermark):** `git diff
  5d36fea..134bc0a -- internal/domain internal/application` ist leer — keine
  Go-Domain-/Use-Case-Änderung in diesem Diff. Eingehalten.
- **Ausschluss 2 (automatische Warnung/Alarmierung):** kein Schwellenwert-
  Code, keine neue Fehlerklasse in diesem Diff. Eingehalten.
- **Ausschluss 3 (Löschausführung selbst):** `RunRetentionUseCase`
  unverändert (`git diff` gegen `internal/application/usecase/retention/`
  leer). Eingehalten.

## 6. Explizit NICHT geprüft (korrekt außerhalb dieser Rolle)

Wie durch §Träger im Repo ohne Wellen-Betrieb (Modul 6/8) gedeckt:
**Closure-Notiz** (DoD 6), **Beobachtungs-Register-Fortschreibung für diesen
Slice** (DoD 8), **Risiko-Ausgänge §6** (DoD 9, reale Grundlage aber unter
§3 dieses Berichts geprüft), **die drei Paarungen** (DoD 10) — alle
zugehörigen §2-Häkchen sind **korrekt unbeansprucht**. Das
Reconciliation-Register-Item (DoD 7) entfällt strukturell (Greenfield-Repo).
Auch nicht Gegenstand: Validierung gegen realen Bedarf (kein
MVP-Meilenstein-Slice, kein Validator-Zug ausgelöst).

## Finding V-1 — DoD-Checkbox „Review durchgeführt" nicht nachgezogen

- **Klasse:** Verifier-only — für Tests und Review unsichtbar; nur ein Blick
  auf den Formular-Zustand *nach* dem Review deckt die Lücke auf. Dieselbe
  Klasse wie `V-1` in `verify-slice-043.md`/`verify-slice-044.md` und
  ursprünglich `verify-slice-039.md` — bereits als `BEO-PGC/dod-checkbox-nachzug`
  verkörpert (Pflicht-Zeile in `.claude/commands/implement-slice.md` Schritt
  18, seit `welle-5`) **und** trotzdem wiederholt aufgetreten (`slice-039`,
  `slice-043`, `slice-044`, jetzt `slice-045`) — ein Hinweis, dass die
  verkörperte Regel den Nachzug **nach** dem Reviewer-Rollenwechsel (statt
  nur im Implementer-Lauf selbst) nicht abdeckt, analog zur bereits
  benannten Lücke in `BEO-PGC/dod-checkbox-nachzug-architect-pfad`
  (Architect-Pfad). Kein neuer Zähler-Beitrag durch diese Prüfung selbst
  (das ist Planner-Einordnung bei Closure), aber die Beobachtung gehört in
  die Sichtung.
- **Befund:** DoD-Punkt 4 in
  `slice-045-blockierende-consumer-sichtbarkeit.md:91` steht auf `- [ ]`,
  obwohl die Bedingung — Review durchgeführt, Report liegt vor, 0
  HIGH/MEDIUM/LOW — seit `dfe02b0` faktisch erfüllt ist.
- **Einordnung:** kein inhaltlicher Mangel — das Review ist real und
  vollständig; reine Formular-Diskrepanz.
- **Erwartete Korrektur:** Checkbox auf `[x]` setzen, in einem eigenen,
  kleinen Commit vor dem `git mv` nach `done/` (Inhalt vor Move, Modul 5
  §git mv + Inhaltsänderung). Kein Rollback, keine Rückführung.

## Verdikt

**DoD-Konformität: bestätigt**, mit einer benannten, nicht
merge-blockierenden Formular-Diskrepanz (V-1) und einer Formulierungs-
Unschärfe in DoD-Punkt 2 (§2 dieses Berichts), die inhaltlich **nicht** zu
einer Rückführung führt: `LH-FA-RET-005` ist real und vollständig erfüllt —
Happy Path durch die neue `cdc.retention_blockers`-View (real gegen zwei
Consumer getestet, dieser Lauf), Boundary durch die bereits bestehende,
unveränderte `cdc.consumer_status`/`cdc_consumer_lag`-Kette (real in
diesem Lauf reproduziert). Alle vier substanziellen Implementer-DoD-Punkte
(1, 2, 3, 5) sind durch eigene, unabhängige Reproduktion gedeckt:
`make gates`, `make test-integration` und `make test-store` liefen in dieser
Sitzung selbst grün, das Benutzerhandbuch ist sachlich korrekt fortgeschrieben.

**Plan-vs-Code-Diff:** keine unbegründete Abweichung. View-Name, Grant-
Erweiterung und der bewusste CLI-Verzicht sind plankonform und im
Plan-Nachzug sauber begründet; der CLI-Verzicht ist zudem laut Plan §1
ohnehin optional, keine übersehene Pflicht.

**§6-Risiken:** beide Risiken tragen eine reale, in dieser Sitzung
nachgeprüfte Grundlage für den Ausgang *entfallen* (INNER JOIN schließt
Risiko 1 strukturell aus; realer, wiederholter Migrationslauf ohne Drift
widerlegt Risiko 2) — die Zuweisung selbst bleibt Planner-Urteil.

**Scope-Treue (§1):** eingehalten, kein Ausschluss verletzt.

**Die vier Planner-Closure-Punkte** (Closure-Notiz, Beobachtungs-Register-
Fortschreibung für diesen Slice, Risiken-Ausgänge, drei Paarungen) sind
**korrekt unbeansprucht** — bewusst nicht Teil dieser Prüfung.

**Keine Rückführung nötig.** Weder `in-progress→next` (drei Liefer-Punkte
sauber erfüllt, keine unerwartete Komplexität) noch `in-progress→open` (kein
Blocker). Vor dem `git mv` nach `done/` sind lediglich die
Checkbox-Korrektur aus V-1 und die Risiko-Ausgänge (§6, mit der unter §3
belegten Grundlage) fällig — Planner-Closure-Arbeit, kein Zerlegungs- oder
Blocker-Fall.

**Übergabe an Planner:** Dieser Bericht bestätigt DoD-Konformität für die
Closure-Entscheidung, mit dem Hinweis V-1 zur Nachbesserung vor dem
`git mv` und der eigenständig geprüften Antwort auf Reviewer-Finding F-1
(§2). Kein Validator-Zug ausgelöst — `slice-045` ist kein
MVP-Meilenstein-Slice im Sinn von Modul 8.

---

*Dieser Bericht ist ein Lauf-Beleg (Modul 11) und wird über Läufe hinweg
nicht wieder gelesen.*
