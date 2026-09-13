# Verifikationsbericht: slice-043 — 2026-09-13

**Rolle:** Verifier (Modul 11) — Prüfung „Bauen wir es richtig?" gegen Plan
(`slice-043` §1/§2 DoD/§3 Plan-Nachzug/§4/§6/§8) und `ADR-0009`/`0011`/
`0012`/`0014`, nicht gegen Diff (Reviewer-Aufgabe, bereits abgeschlossen)
und nicht gegen realen Bedarf (Validator, hier nicht ausgelöst — kein
MVP-Meilenstein-Slice).

**Frischer Kontext:** Diese Prüfung liest den vollständigen, aktuellen
Slice-Plan, alle vier referenzierten ADRs vollständig, beide Review-Reports
(Erstlauf und Fixrunde) vollständig und den tatsächlichen Code selbst —
keine Behauptung aus einem Bericht wird ungeprüft übernommen; jeder unten
genannte Sensor-/Werkzeuglauf wurde in dieser Sitzung **selbst**
ausgeführt, nicht aus den Reports zitiert.

**Gegenstand:**
`docs/plan/planning/in-progress/slice-043-changestoreport-loeschmethode.md`
zum Stand `HEAD = 3177b9c`. Commits (chronologisch): `6587a2e`
(`next→in-progress`), `de3ff8c` (Implementierung), `d2b5e6d`
(Review-Erstlauf, 2 MEDIUM: F-1, F-2), `193c47b` (Fixrunde: F-1 real
behoben, F-2 bewusst beibehalten und begründet), `923afe5`
(Review-Bestätigung: F-1 geschlossen, F-2 akzeptiert, neuer Nebenbefund
F-3), `3177b9c` (F-3 behoben — reine Zitat-Korrektur in Prosa). Sequenz
selbst geprüft (`git log --oneline 6587a2e^..3177b9c`): Move →
Implementierung → Review → Fix → Fixrunden-Bestätigung → Zitat-Korrektur —
keine Rolle springt rückwärts ohne Übergabe-Artefakt.

---

## 1. DoD-Konformität, Punkt für Punkt

| # | DoD-Punkt | Verdikt | Beleg (eigene Prüfung) |
|---|---|---|---|
| 1 | `ChangeStorePort` trägt neue Löschmethode; `PostgresChangeStoreAdapter` implementiert sie real gegen PostgreSQL, `cdc.change`-Zeilen werden tatsächlich entfernt, `make test-store` real belegt | **erfüllt** | `internal/application/port/outbound/changestore.go` gelesen: `DeleteChanges(ctx, changeIDs []model.ChangeID) error` am Port. `internal/adapters/driven/postgresstorage/store.go:169-234` implementiert sie real: SQL `DELETE FROM cdc.change WHERE change_id = ANY($1) RETURNING transaction_id` (`queries.go:73-75`). Selbst ausgeführt: `make test-store` gegen realen Testcontainer (Schema-Rollout via d-migrate, dann `go test ./...`) — alle Pakete `ok`, insbesondere `postgresstorage` (6.255s). Testcode selbst gelesen (`store_test.go:467-542`): `TestDeleteChangesRemovesOnlyGivenChanges` zählt reale verbleibende Zeilen nach echtem SQL-DELETE, `TestDeleteChangesIsIdempotent` und `TestDeleteChangesRemovesOrphanedTransactionOnly` ebenso — kein Fake-Store beteiligt |
| 2 | `RunRetentionUseCase` ruft `RetentionPolicy.AllowsDeletion` je betrachtetem Change real auf (Alter, Change-Position, alle bestätigten Consumer-Positionen), übergibt ausschließlich freigegebene Changes an die Port-Methode — Fake-Port-Test mit gemischter Menge | **erfüllt** | `internal/application/usecase/retention/service.go:53-81` (`Run`) gelesen: liest `s.state.Positions(ctx, command.Source)` (alle bestätigten Consumer-Positionen der Quelle), liest `s.store.ReadChanges`, berechnet `age := now.Sub(record.CommittedAt)` je Record, ruft **je Change in der Schleife** `command.Policy.AllowsDeletion(age, record.Position, consumerPositions)` auf und sammelt nur `true`-Treffer in `eligible`; `s.store.DeleteChanges(ctx, eligible)` erhält ausschließlich diese Menge — kein zweiter Aufrufpfad. `internal/domain/model/retention.go` (`AllowsDeletion`) bestätigt die Signatur (Alter, Change-Position, `[]ConsumerPosition`). `TestRunDistinguishesEligibleChangesFromMixedSet` (`service_test.go`) real gelesen und via `make test` selbst ausgeführt: vier Changes (löschbar, zu jung, Consumer zurückhängend, Boundary) — Assertion erzwingt `deletedIDs == [c-eligible, c-boundary]`, exakt die erwartete Unterscheidung |
| 3 | `LH-FA-RET-002` Negative (ungültige Konfiguration → expliziter Fehler statt stiller Übernahme) real erfüllt, Test für Fehlerpfad | **erfüllt** | Zwei unabhängige Ebenen selbst nachvollzogen: (a) Domänenebene — `NewRetentionPolicy` weist negative `MinAge` über `ErrNegativeDuration` zurück (bereits vor diesem Slice, unverändert); (b) Use-Case-Ebene — `Run` weist `command.Source == ""` über `domainerrors.ErrEmptyIdentifier` zurück, **bevor** irgendein Port berührt wird (`service.go:54-56`). `TestRunRejectsEmptySource` (`service_test.go`) real gelesen: prüft sowohl den Fehler (`errors.Is(err, ErrEmptyIdentifier)`) als auch `store.readCalls == 0 && !store.deleteCalled && state.positionCalls == 0` — die Null-Berührung ist Teil der Assertion, nicht nur behauptet. Test lief grün in `make test` (eigener Lauf) |
| 4 | `make gates` grün, `make test`/`make test-store` grün | **erfüllt, selbst reproduziert** | Alle drei Läufe in dieser Sitzung selbst ausgeführt gegen `HEAD = 3177b9c`: `make gates` → `baseline-verify` (v6.5.0, 54 Dateien) OK, `d-check` Struktur (346 Dateien, 0 Befunde), `d-check` Commits (`commits`-Modul, `HEAD~5..HEAD`, 0 Befunde), `commit-traceability.sh` (5 Commits, Betreffs ohne Struktur-ID) OK, `a-check` (0 Befunde). `make test` → alle Pakete `ok`, inkl. `internal/application/usecase/retention`. `make test-store` → realer Testcontainer, Schema-Rollout via d-migrate erfolgreich, alle Pakete `ok`, `postgresstorage` real gegen PostgreSQL (6.255s) |
| 5 | Review durchgeführt, Report unter `docs/reviews/` liegt vor | **inhaltlich erfüllt, Formular-Diskrepanz** | Beide Reports vollständig gelesen: `docs/reviews/review-slice-043.md` (2 MEDIUM: F-1, F-2) und `docs/reviews/review-slice-043-fixrunde.md` (F-1 bestätigt behoben, F-2 als Entscheidung akzeptiert, neuer Nebenbefund F-3 — inzwischen mit `3177b9c` behoben). Die Bedingung ist damit tatsächlich erfüllt. **Aber:** Checkbox in §2 der Plan-Datei (Zeile 108) steht weiterhin auf `- [ ]` — keiner der Commits `d2b5e6d`/`193c47b`/`923afe5`/`3177b9c` hat sie auf `[x]` gesetzt (`git show <commit> --stat` je geprüft: keiner berührt den Checkbox-Bereich). Siehe Finding V-1 unten |
| 6 | Doku-Update, falls öffentlicher Vertrag entsteht — geprüft: keiner, `DeleteChanges`/`Positions`/`RunRetentionUseCase` sind interne Go-Schnittstellen | **erfüllt** | Selbst nachvollzogen: Keine der drei neuen/geänderten Schnittstellen hat eine CLI-/SQL-Außenfläche — `RunRetentionUseCase` hat in diesem Slice **keinen Aufrufer** (Out-of-Scope-Punkt 1, `slice-044` liefert ihn); `grep -rn "RunRetentionUseCase" internal --include=*.go -l` trifft nur die Definitions- und Implementierungsdatei, keinen Composition-Root-/CLI-Code. Begründung im Plan-Nachzug deckt sich mit dem Befund |
| 7 | Closure-Notiz mit Steering-Loop-Lerneintrag | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit (Modul 8); §7 des Plans ist noch die Bedienhinweis-Vorlage. Kein Verifikations-Gegenstand dieser Prüfung |
| 8 | Reconciliation-Register fortgeschrieben, falls einschlägig | **entfällt strukturell** | `find docs/plan/planning -maxdepth 1 -iname "reconciliation*"` ohne Treffer — `harness/conventions.md` §Modus-Deklaration führt ausschließlich die Sub-Area `*`/`PGC` im Modus Greenfield, kein Brownfield-Bootstrap. Das Item nennt diese Bedingung selbst; kein offener Planner-Punkt, sondern strukturell nie zutreffend für dieses Repo |
| 9 | Beobachtungs-Register fortgeschrieben | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit. Zur Einordnung selbst geprüft: `docs/plan/planning/observations/BEO-PGC/retention-keine-loeschausfuehrung/` existiert bereits (0×, im Plan §8 als „benannt nicht gezählt" dokumentiert) — deckt sich mit der im Plan dokumentierten Sichtung |
| 10 | Jedes Risiko aus §6 trägt einen Ausgang | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit; beide Risiken in §6 tragen noch `<bei Closure einzutragen>`. Der Plan-Nachzug §3 Punkte 6 und 2 liefert bereits die inhaltliche Grundlage für beide Ausgänge (strukturell ausgeschlossener Konflikt bzw. gebauter Lesezugriffsweg), aber das Eintragen selbst ist Planner-Arbeit |
| 11 | Drei Paarungen (Anker · Folge-Slice · Register) | **offen — korrekt unbeansprucht** | Planner-Closure-Arbeit, wellenlos hier statt bei einer Welle-Closure fällig (`welle-13` ist noch nicht geschlossen), aber ohnehin erst **nach** dem `git mv` nach `done/` sinnvoll prüfbar |

## 2. Finding V-1 — DoD-Checkbox „Review durchgeführt" nicht nachgezogen

- **Klasse:** Verifier-only — für Tests und Review unsichtbar, weil beide
  Rollen ihre eigene Arbeit erledigt haben; nur ein Blick auf den
  *Formular-Zustand nach* der vollständigen Review-Sequenz (Erstlauf +
  Fixrunde + Zitat-Korrektur) deckt die Lücke auf.
- **Befund:** DoD-Punkt 5 in
  `slice-043-changestoreport-loeschmethode.md:108` steht auf `- [ ]`,
  obwohl die Bedingung — Review durchgeführt, Report liegt vor — seit
  Commit `d2b5e6d` faktisch erfüllt ist und seit `923afe5`/`3177b9c`
  zusätzlich vollständig abgeschlossen (alle drei Findings mit Ausgang)
  vorliegt.
- **Einordnung:** kein inhaltlicher Mangel an Implementierung oder
  Review-Substanz — beide sind, wie oben belegt, real und reproduzierbar
  erfüllt. Es ist eine Diskrepanz zwischen dem DoD-Formular und der
  tatsächlichen Sachlage — dieselbe Finding-Klasse wie `V-1` in
  `verify-slice-039.md`, exakt die Klasse, für die der Verifier existiert.
- **Erwartete Korrektur:** Checkbox auf `[x]` setzen, in einem eigenen,
  kleinen Commit vor dem `git mv` nach `done/` (Inhalt vor Move, Modul 5
  §git mv + Inhaltsänderung). Kein Rollback, keine Rückführung
  (`in-progress→next`/`open`) — die Korrektur ist ein Ein-Zeilen-Nachzug.

## 3. Scope-Treue gegen §1 (Ausdrücklich NICHT in diesem Slice)

Eigenständig geprüft, nicht aus dem Review übernommen:

- **Ausschluss 1 (Hintergrund-Job/CLI-Trigger, `slice-044`):**
  `grep -rn "RunRetentionUseCase" internal --include="*.go" -l` trifft nur
  `retention/service.go` und `port/inbound/retention.go` — kein
  Composition-Root-, CLI- oder Cron-Code ruft den Use Case auf. Eingehalten.
- **Ausschluss 2 (Sichtbarkeit blockierender Consumer, `LH-FA-RET-005`,
  `slice-045`):** `grep -rn "blockierend" internal --include="*.go"` trifft
  nur einen Dokumentationskommentar am Inbound-Port, der den Ausschluss
  selbst benennt („Kein blockierender Consumer wird über diesen Use Case
  sichtbar gemacht"), keinen Code, der das täte. Eingehalten.
- **Ausschluss 3 (`cdc_storage_bytes`-Metrik, `slice-046`):**
  `grep -rn "cdc_storage_bytes" internal tools/schema` ohne Treffer.
  Eingehalten.
- **Ausschluss 4 (Laufzeit-Konfigurierbarkeit von `RetentionPolicy.MinAge`):**
  `NewRetentionPolicy`/`RetentionPolicy` unverändert seit vor diesem Slice
  (`git log -p 6587a2e..3177b9c -- internal/domain/model/retention.go` ohne
  Treffer — die Datei ist in keinem Commit dieses Slice berührt).
  Eingehalten.
- **Vollständiger Diff seit dem `next→in-progress`-Move** (`git diff
  6587a2e..3177b9c --stat`): 17 Dateien im Implementierungs-Commit plus
  Plan-Datei- und Review-Report-Änderungen der Folgecommits — keine Datei
  außerhalb des in §3 der Plan-Tabelle benannten Umfangs (Port, Adapter,
  Use-Case, Consumer-State-Port-Erweiterung, Fake-Nachzieh-Patches in
  bestehenden Use-Case-Tests, die durch die neue `ConsumerStatePort`-Methode
  strukturell nötig wurden).

## 4. Weitere Hard-Rule-/ADR-Prüfungen (eigenständig)

- **`ADR-0009` (Change Store als Outbound Port):** `ChangeStorePort` bleibt
  der einzige Persistenz-/Löschzugang; `DeleteChanges` ist eine weitere
  Fähigkeit desselben Fähigkeits-Ports, kein zweiter Port für dieselbe
  Konsistenzgrenze — konsistent mit der in `ADR-0009` getroffenen
  Options-C-Entscheidung.
- **`ADR-0011` (Persist-before-ACK) / `ADR-0012` (At-Least-Once) — kein
  Bypass am Domain Core:** `DeleteChanges` ist die einzige Löschfähigkeit
  am Port und wird in `service.go` ausschließlich mit der über
  `AllowsDeletion` freigegebenen Teilmenge aufgerufen — selbst am Code
  nachvollzogen, kein zweiter Aufrufpfad existiert (`grep -rn
  "DeleteChanges(" internal --include="*.go"` trifft genau die
  Port-Definition, die Adapter-Implementierung, den einen Aufruf in
  `service.go` und die Testaufrufe). Das Idempotenz-Argument aus dem
  Plan-Nachzug §3 Punkt 6 (Crash-Replay kann laut `LH-QA-REL-001.a` nur vor
  dem Source-ACK auftreten, also lange bevor ein Consumer sie verarbeitet
  und `AllowsDeletion` sie zur Löschung freigibt) ist eine Reihenfolge-
  Argumentation, die dieser Bericht nachvollzieht, aber nicht durch einen
  eigenen Determinismus-Test verifizieren kann — sie bleibt ein Struktur-
  Argument, kein Automat-Beleg. Kein Widerspruch zu `ADR-0011`/`ADR-0012`
  gefunden.
- **`ADR-0014` (Retention als Domain Policy) — zentraler Prüfmaßstab:**
  Die Freigabe-Entscheidung liegt vollständig in
  `RetentionPolicy.AllowsDeletion` (Domain, unverändert in diesem Diff);
  `RunRetentionService.Run` (Application) orchestriert nur, trifft keine
  eigene Freigabe-Logik; `PostgresChangeStoreAdapter.DeleteChanges`
  (Adapter) führt nur aus, was ihm übergeben wird — exakt die in `ADR-0014`
  festgelegte Trennung von Entscheidung und Ausführung. Bestätigt.
- **`AGENTS.md` §3.5 (ADR-Immutabilität):** `git log --oneline --all --
  docs/plan/adr/0009-*.md docs/plan/adr/0011-*.md docs/plan/adr/0012-*.md
  docs/plan/adr/0014-*.md` liefert für alle vier je zwei Commits (Anlage +
  ein früherer `.d-check.yml`-Schärfungs-Commit, der die ADR-Dateien nicht
  inhaltlich berührt) — keiner der `slice-043`-Commits taucht auf; alle vier
  bleiben seit ihrer Erstanlage unverändert.
- **`AGENTS.md` §3.7 (Kommentar-Disziplin), keine Slice-Chronik:**
  Eigener vollständiger `grep`-Durchlauf über alle 16 seit `6587a2e`
  geänderten `.go`-Dateien nach `slice-[0-9]+|welle-[0-9]+|review-slice`
  → **keine Treffer**. Ein weiterer Lauf nach
  `früher|vorherig|ehemals|zuvor stand|wäre gewesen|hätte` trifft nur
  Domänen-Ordnungssprache („eine frühere Position" = eine Position, die in
  der Ordnungsrelation vor einer anderen liegt, keine Diff-Historie) und
  einen bereits vor diesem Slice bestehenden Testkommentar
  (`store_test.go:37`, „vorherigen Tests" im Sinn von Testisolation) — kein
  neuer Treffer, `BEO-PGC/slice-chronik-in-code-kommentar` bewegt sich mit
  diesem Slice nicht.
- **`AGENTS.md` §3.1 (Docker-only):** Der Diff fügt keine lokale
  Toolchain-Installation hinzu; alle neuen Tests laufen über die
  bestehende `make test`/`make test-store`-Kette (selbst ausgeführt, siehe
  oben).

## 5. Explizit NICHT geprüft (korrekt außerhalb dieser Rolle)

Wie im Auftrag benannt und durch §Träger im Repo ohne Wellen-Betrieb
(Modul 6/8) gedeckt: **Closure-Notiz** (DoD 7), **Beobachtungs-Register-
Fortschreibung** (DoD 9), **Risiko-Ausgänge §6** (DoD 10), **die drei
Paarungen** (DoD 11). Alle vier zugehörigen §2-Häkchen sind **korrekt
unbeansprucht** — kein Mangel, sondern der vorgesehene Zustand vor dem
nächsten Rollenwechsel an den Planner. Das Reconciliation-Register-Item
(DoD 8) entfällt strukturell (Greenfield-Repo, keine gesonderte Zählung in
dieser Übersicht nötig). Auch nicht Gegenstand: Validierung gegen realen
Bedarf (kein MVP-Meilenstein-Slice, kein Validator-Zug ausgelöst).

## Verdikt

**DoD-Konformität: bestätigt**, mit einer benannten, nicht
merge-blockierenden Formular-Diskrepanz (V-1: Checkbox 5 nicht
nachgezogen — inhaltlich längst erfüllt). Alle sechs substanziellen
Implementer-DoD-Punkte (1–6) sind durch eigene, unabhängige Reproduktion
gedeckt: `DeleteChanges` real gegen PostgreSQL getestet (inkl. der
atomaren Waisen-Bereinigung aus der Fixrunde, `TestDeleteChangesRemovesOrphanedTransactionOnly`
selbst über `make test-store` gegen den echten Testcontainer gelaufen),
`RunRetentionService.Run` ruft `AllowsDeletion` je Change real mit Alter,
Change-Position und allen bestätigten Consumer-Positionen auf und übergibt
ausschließlich die freigegebene Menge, `LH-FA-RET-002` Negative real über
zwei unabhängige Schichten erfüllt, `make gates`/`make test`/`make
test-store` selbst grün gelaufen, kein öffentlicher Vertrag fällig. Die vier
Planner-Closure-Punkte (7, 9, 10, 11) sind korrekt offen; das
Reconciliation-Register-Item (8) entfällt strukturell.

**Zu `ADR-0009`/`0011`/`0012`/`0014` (zentraler Prüfmaßstab):** bestätigt —
keine ADR verändert, die Domain-Policy-Trennung aus `ADR-0014` ist
strukturell eingehalten (Domain entscheidet, Application orchestriert,
Adapter führt aus), kein Bypass des Persist-before-ACK-/At-Least-Once-Pfads.

**Keine Rückführung nötig.** Weder `in-progress→next` (der Slice ist
nicht zu groß — die drei substanziellen Liefer-Punkte aus §2 sind sauber
erfüllt, kein vierter Liefer-Punkt und keine dritte Schicht entstanden;
die Fixrunde hat den bestehenden Scope vertieft, nicht erweitert) noch
`in-progress→open` (kein Blocker). Vor dem `git mv` nach `done/` ist
lediglich die Checkbox-Korrektur aus V-1 fällig — ein Ein-Zeilen-Commit,
kein Zerlegungs- oder Blocker-Fall.

**Übergabe an Planner:** Dieser Bericht bestätigt DoD-Konformität für die
Closure-Entscheidung, mit dem Hinweis V-1 zur Nachbesserung vor dem
`git mv`. Kein Validator-Zug ausgelöst — `slice-043` ist kein
MVP-Meilenstein-Slice im Sinn von Modul 8.

---

*Dieser Bericht ist ein Lauf-Beleg (Modul 11) und wird über Läufe hinweg
nicht wieder gelesen.*
