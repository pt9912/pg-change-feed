# Review-Report: slice-043 — Fixrunde — 2026-09-13

**Review-Art:** Code — Bestätigungslauf zu einer Fixrunde nach eigenem
Vorbefund. Geprüft gegen den eigenen vorherigen Report
(`review-slice-043.md`, Findings F-1/F-2), den Plan-Nachzug in
`slice-043-changestoreport-loeschmethode.md` §3 Punkt 7/8, `ADR-0014`,
`ADR-0029`, `ADR-0005` und `AGENTS.md` §3 Hard Rules (§3.5, §3.7) —
Rollentrennung Modul 8: diese Prüfung läuft eigenständig gegen Code und
Belege, nicht als Übernahme der Implementer-Zusammenfassung.

**Gegenstand:** Commit `193c47be66e77f47b0363e82bb4584fe0456305e`
(`fix(retention): Transaktions-Waisen und Consumer-Abwesenheits-Lesart
geklärt (LH-FA-RET-004, review-slice-043 F-1/F-2)`) — geändert:
`internal/adapters/driven/postgresstorage/store.go`,
`internal/adapters/driven/postgresstorage/queries/queries.go`,
`internal/adapters/driven/postgresstorage/store_test.go`,
`internal/application/port/outbound/consumerstate.go`, Plan-Nachzug in
`slice-043-changestoreport-loeschmethode.md`.

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/reviews/review-slice-043.md` (vollständig — eigener Vorbefund F-1/F-2)
- Vollständiger `git show 193c47b` (alle 5 Dateien)
- `internal/adapters/driven/postgresstorage/store.go`,
  `internal/adapters/driven/postgresstorage/queries/queries.go` (vollständig
  gelesen, nicht nur der Diff-Ausschnitt)
- `internal/domain/model/consumer.go` (`ConsumerPosition.Advance`,
  unverändert in diesem Diff — Prüfgrundlage für F-2)
- `docs/plan/adr/0029-domain-invarianten.md`,
  `docs/plan/adr/0005-sourceposition-abstrahiert-lsn.md` (vollständig)
- `tools/schema/schema.yaml` (`transaction`, `change`, FK-Kanten, keine
  `on_delete`-Angabe)
- `tools/schema/nacharbeit-observability.sql` (referenziert, nicht Teil des
  Diffs)
- reales `make test-store` (Testcontainer, PostgreSQL) und `make gates`

---

## Findings

### F-1 (Vorbefund) — verwaiste `cdc.transaction`-Zeilen — bestätigt behoben

- `kategorie`: MEDIUM (Vorbefund, jetzt geschlossen)
- `quelle`: Maintainability (Out-of-Scope-Disziplin, wie im Vorbefund)
- `pfad`: `internal/adapters/driven/postgresstorage/store.go:169-234`
  (`DeleteChanges`), `internal/adapters/driven/postgresstorage/queries/queries.go:62-91`
  (`DeleteChanges`, `DeleteOrphanedTransactions`)
- `befund`: `DeleteChanges` läuft jetzt als eine DB-Transaktion
  (`pool.Begin`/`defer tx.Rollback`/`tx.Commit`): Die erste Abfrage löscht
  die freigegebenen `cdc.change`-Zeilen und liefert per `RETURNING
  transaction_id` die betroffenen Transaktions-Kennungen zurück; eine
  zweite Abfrage (`DeleteOrphanedTransactions`) löscht aus genau dieser
  Menge nur die Transaktionen, für die `NOT EXISTS (SELECT 1 FROM
  cdc.change c WHERE c.transaction_id = cdc.transaction.transaction_id)`
  gilt. Beide Statements teilen sich dieselbe `pgx.Tx`; ein Fehler in
  irgendeinem Schritt lässt `defer tx.Rollback` greifen, bevor `Commit`
  erreicht wird — kein Teilzustand möglich. Der reale Test
  `TestDeleteChangesRemovesOrphanedTransactionOnly` legt zwei
  Transaktionen an (`t-1` mit einer Change-Zeile, `t-2` mit zwei), löscht
  `t-1-1` und `t-2-1` und bestätigt: `t-1` verschwindet (0 Zeilen), `t-2`
  bleibt bestehen (1 Zeile) mit ihrer verbleibenden Change-Zeile intakt (1
  Zeile). Das ist genau der Kollateralschaden-Fall, den die Aufgabe
  prüfen sollte, und er tritt nicht ein.
- `verifizierbar`: ja — `make test-store` real ausgeführt (eigener Lauf,
  nicht übernommen): `ok
  github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage
  4.308s`, kein Build-Tag schließt den Test aus, er ist Teil desselben
  Pakets und lief gegen den echten Testcontainer.
- `klasse`: „Out-of-Scope-Lücke ohne Adresse" (aus dem Vorbefund, jetzt
  geschlossen statt offen)

### F-2 (Vorbefund) — Consumer-Abwesenheits-Lesart — Entscheidung bestätigt,
Dokumentation klar, mit einer neuen Randbeobachtung

- `kategorie`: MEDIUM (Vorbefund) → Entscheidung akzeptiert (kein Code-Fix
  nötig, war als Interpretationsfrage klassifiziert)
- `quelle`: `LH-FA-RET-004`, `ADR-0029`
- `pfad`: `internal/application/port/outbound/consumerstate.go:82-94`
  (`Positions`-Dokumentation), `internal/domain/model/consumer.go:50-66`
  (`Advance`)
- `befund`: Die Entscheidung selbst — Abwesenheits-Lesart bleibt
  unverändert, ein registrierter, aber gegen eine Quelle noch nie
  bestätigender Consumer blockiert `AllowsDeletion` für diese Quelle
  nicht — ist am Code nachvollziehbar: `ConsumerPosition.Advance` prüft
  `SourceMismatch`/`PositionRegression` nur, wenn `c.Acknowledged()`
  bereits wahr ist; vor der ersten Bestätigung gibt es keine
  Quellen-Bindung, die `Positions` melden könnte. Die geschärfte
  Dokumentation in `consumerstate.go:82-94` ist für `slice-044` klar
  genug: Sie benennt die Grenze explizit („blockiert damit keine
  Löschung", „Schutz … entsteht erst mit … `Acknowledge`") statt sie nur
  implizit aus der Abwesenheits-Beschreibung ableiten zu lassen — ein
  künftiger Implementer muss die Konsequenz nicht selbst herleiten.
- `verifizierbar`: nein (Interpretationsfrage, wie im Vorbefund)
- `klasse`: „Consumer-Sichtbarkeitslücke vor Erstbestätigung" (aus dem
  Vorbefund; Ausgang: Entscheidung dokumentiert, operative Konsequenz für
  `slice-044` benannt statt aufgelöst — bewusst, siehe Plan-Nachzug §3
  Punkt 8)

### F-3 — `ADR-0029` Regel 2 belegt nicht das zitierte Verhalten

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Traceability-Genauigkeit einer ADR-Referenz)
- `pfad`: Commit-Message `193c47b` (Absatz „F-2"), Plan-Nachzug
  `slice-043-changestoreport-loeschmethode.md` §3 Punkt 8 (Zeile mit „ein
  Consumer bindet sich erst mit seiner ersten `Advance`
  (`ADR-0029` Regel 2) an eine Quelle")
- `befund`: `ADR-0029`s Regel 2 lautet wörtlich „Consumer-ACK regulär nur
  vorwärts" — eine Monotonie-Zusage, keine Aussage über
  Quellen-Bindung. Die im Fix zitierte Eigenschaft („ein Consumer bindet
  sich erst mit der ersten `Advance` an eine Quelle") steht stattdessen im
  bestehenden Code-Kommentar zu `ConsumerPosition.Advance`
  (`internal/domain/model/consumer.go:53-55`, unverändert in diesem Diff)
  und wird dort — selbst mit loser Passung — `ADR-0005` zugeordnet, nicht
  `ADR-0029`. Repo-weit ist die Zuordnung „Regel 2 = Monotonie" etabliert
  und konsistent (`consumer.go:51`, `review-slice-009.md:225/299`,
  `retention_test.go:11` nutzt Regel 5 für eine andere Zusage) — der
  neue Fix ist die erste Stelle, die Regel 2 für die Bindungs-Eigenschaft
  statt für die Vorwärts-Monotonie beansprucht. Die *Entscheidung* selbst
  (F-2, oben) bleibt davon unberührt — nur ihre zitierte Begründung trägt
  nicht das, was sie zu tragen behauptet. Betroffen sind ausschließlich
  Prosa-Artefakte (Commit-Message, Plan-Nachzug); die tatsächlich von
  `slice-044` konsultierte Quelle, die `Positions`-Dokumentation in
  `consumerstate.go`, zitiert `ADR-0029` an dieser Stelle gar nicht
  (nur `LH-FA-CON-001`) und ist von diesem Fehlzitat nicht betroffen.
- `verifizierbar`: ja — `docs/plan/adr/0029-domain-invarianten.md` Zeile 33
  („2. Consumer-ACK regulär nur vorwärts.") gegen den zitierten Satz im
  Plan-Nachzug.
- `klasse`: „ADR-Zitat trägt eine andere Regel als die referenzierte
  Nummer" (erstes Auftreten dieser Klasse)

## Negativbefunde

- geprüft, ohne Befund: **Atomarität von `DeleteChanges`.** Beide
  SQL-Anweisungen laufen auf derselben `pgx.Tx`
  (`store.go:188-230`); `defer tx.Rollback(ctx)` deckt jeden
  Fehlerpfad vor `tx.Commit(ctx)` ab — kein Zwischenzustand erreichbar,
  in dem Changes gelöscht, aber die Waisen-Prüfung nicht gelaufen wäre
  (oder umgekehrt).
- geprüft, ohne Befund: **`DeleteOrphanedTransactions`-Scope ist eng auf
  die tatsächlich betroffene Menge begrenzt.** Die Abfrage nimmt nur die
  von `touchedTransactions` gesammelten (deduplizierten) Kennungen entgegen
  — eine Transaktion, die nie eine Change-Zeile trug, taucht in
  `RETURNING transaction_id` nie auf und wird von dieser Abfrage nie
  berührt (deckt sich mit der Begründung im Plan-Nachzug §3 Punkt 7,
  „ein anderer Fall").
- geprüft, ohne Befund: **`make test-store` real ausgeführt, grün.**
  Kompletter Lauf inkl. Schema-Rollout, alle Pakete `ok`, insbesondere
  `internal/adapters/driven/postgresstorage` (4.308s) und
  `internal/application/usecase/retention` (0.006s); kein Build-Tag
  schließt den neuen Test aus.
- geprüft, ohne Befund: **`make gates` real ausgeführt, grün.**
  `baseline-verify` (54 Dateien), `docs-check` (345 Dateien, 0 Befunde,
  zweimal — Struktur- und Commit-Traceability-Lauf), `commit-traceability`
  (5 Commits, Betreffs ohne Struktur-ID), `a-check` (0 Befunde).
- geprüft, ohne Befund: **Kommentar-Disziplin (`AGENTS.md` §3.7), keine
  neue Slice-Chronik in den vier geänderten `.go`-Dateien.** Eigener
  vollständiger `grep`-Durchlauf über die **hinzugefügten** Zeilen (nicht
  nur den Dateibestand) nach
  `slice-\d+|review-slice|welle-\d+|vor diesem|seit |nur noch|jetzt |nicht
  mehr` in allen vier `.go`-Dateien: keine Treffer. Ein Grep über den
  gesamten Dateibestand (statt nur der Diff-Zeilen) meldet zwei
  Alt-Treffer in `queries.go` (Zeilen 207/294) — beide vor diesem Diff
  bestehend und beide keine Slice-Chronik („vor diesem Upsert" = zeitlich
  vor der SQL-Operation, „nicht mehr `pending`" = Zustandsbeschreibung
  eines Antrags, keine Diff-Historie).
- geprüft, ohne Befund: **Neue Kommentare tragen zulässige Klassen
  (`AGENTS.md` §3.7).** Die drei neuen/geänderten Kommentarblöcke
  (`queries.go:62-75`, `queries.go:77-91`, `store.go:169-178`)
  beschreiben ausschließlich den geltenden Zustand und seine Kopplung an
  `cdc.metrics` (Klasse Kopplung/Zusage) — kein Konjunktiv über eine
  verworfene Alternative, kein abwesender Text, kein abgebrochener Satz.
- geprüft, ohne Befund: **Kein ADR verändert (`AGENTS.md` §3.5).**
  `git show 193c47b --name-only` trifft keine Datei unter
  `docs/plan/adr/`.
- geprüft, ohne Befund: **Docker-only (`AGENTS.md` §3.1).** Kein neues
  Build-/Test-Skript, beide Tests laufen über die bestehende `make
  test-store`-Kette.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 (F-3, neu) |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** „ADR-Zitat trägt eine andere Regel als
die referenzierte Nummer" (1×, erstes Auftreten, F-3) — F-1 und F-2 aus
dem Vorbefund sind hier keine neuen Klassen, sondern deren Bestätigung
(F-1 geschlossen, F-2 als Entscheidung akzeptiert); kein
Steering-Loop-Eintrag fällig, keine Klasse erreicht 3×.

## Verdikt

**F-1: bestätigt behoben.** Reale, atomare Bereinigung verwaister
`cdc.transaction`-Zeilen; kein Kollateralschaden an Transaktionen mit
verbleibenden Change-Zeilen; real gegen PostgreSQL belegt
(`TestDeleteChangesRemovesOrphanedTransactionOnly`, selbst ausgeführt statt
aus dem Implementer-Bericht übernommen).

**F-2: Entscheidung akzeptiert, kein Code-Fix erforderlich.** Die
Abwesenheits-Lesart bleibt bewusst unverändert; die Begründung über
`ConsumerPosition.Advance` trägt am Code, die geschärfte Dokumentation in
`consumerstate.go` ist für `slice-044` klar genug. Eine neue, unabhängige
Randbeobachtung (F-3) betrifft nur die *Zitier-Genauigkeit* der
Begründung (`ADR-0029` Regel 2 belegt Monotonie, nicht Quellen-Bindung),
nicht die Entscheidung selbst.

**Merge-blockierend:** nein — der Commit ist bereits gepusht, `make
gates`/`make test-store` sind grün, F-3 ist eine Prosa-Ungenauigkeit ohne
Rollen-Widerspruch (kein HIGH, keine dritte Wiederholung, kein
bestrittenes Verdikt) und braucht keine Architect-Sequenz nach Modul 8.

**Übergabe:** F-3 kann als Korrektur der Zitatzeile im Plan-Nachzug §3
Punkt 8 vor der Closure nachgezogen werden (Prosa-Korrektur, kein
Code-Diff); die Commit-Message selbst ist unveränderlich (`git`). Für
`slice-043`s eigene DoD-Punkte ist das kein Mangel — DoD-/Spec-Konformität
prüft der Verifier separat; dieser Report ist ein Lauf-Beleg und wird über
Läufe hinweg nicht wieder gelesen, die Summary-Zeile speist bei Bedarf den
Closure-Eintrag (Modul 5).
