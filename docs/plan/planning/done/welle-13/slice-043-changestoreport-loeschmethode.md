# Slice slice-043: ChangeStorePort-Löschmethode und RunRetentionUseCase

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-13 — erster Slice, liefert die eigentliche
Löschausführung, auf der `slice-044` (Hintergrund-Job/CLI-Trigger)
aufbaut.

**Bezug:**
[`LH-FA-RET-002`](../../../../spec/lastenheft.md),
[`LH-FA-RET-003`](../../../../spec/lastenheft.md),
[`LH-FA-RET-004`](../../../../spec/lastenheft.md),
[`ADR-0009`](../../../../docs/plan/adr/0009-change-store-outbound-port.md),
[`ADR-0011`](../../../../docs/plan/adr/0011-persist-before-ack.md),
[`ADR-0012`](../../../../docs/plan/adr/0012-at-least-once.md),
[`ADR-0014`](../../../../docs/plan/adr/0014-retention-domain-policy.md)
(alle nur umgesetzt — keine aktive ADR wird geändert).

**Berührte Spec-Stellen:**
[`LH-FA-RET-004.a`](../../../../spec/pflichtenheft.md) (Safe Watermark
der Retention, bereits vorhanden — dieser Slice implementiert das dort
beschriebene Verfahren real). Der Verweis zeigt **aufwärts**: Die Spec
nennt diesen Slice nie (Baseline-Regelwerk
`grundlagen-referenz-richtung.md` §Referenz-Richtung (SDP),
`grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-13.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** `ChangeStorePort` (`internal/application/port/outbound/changestore.go`)
bekommt eine neue Löschmethode (z. B. `DeleteChangesBefore` oder
äquivalent — Implementer entscheidet Namen/Signatur und begründet im
Plan-Nachzug), implementiert vom `PostgresChangeStoreAdapter`. Ein neuer
`RunRetentionUseCase` (Application-Schicht, Muster analog zu
bestehenden Use-Cases wie `EnableTableUseCase`) liest für eine Quelle
alle bestätigten Consumer-Positionen, ruft `RetentionPolicy
.AllowsDeletion` je betrachtetem Change real auf und übergibt die
freigegebene Menge an die neue Port-Methode. Kein Aufrufer dieses
Use-Case in dieser Ebene (Hintergrund-Job/CLI) — das liefert
`slice-044`.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Hintergrund-Job/CLI-Trigger, der `RunRetentionUseCase` real
  aufruft** — Folge-Slice `slice-044`; dieser Slice liefert den
  Use-Case und die Port-Fähigkeit, nicht den Auslösemechanismus
  (Schicht-Abgrenzung: Application/Adapter-Schicht hier, Composition
  Root/CLI in `slice-044`).
- **Sichtbarkeit blockierender Consumer** (`LH-FA-RET-005`) — Folge-Slice
  `slice-045`; ein anderer Liefer-Fokus (Observability, nicht
  Löschausführung selbst).
- **`cdc_storage_bytes`-Metrik** (`LH-FA-RET-006`) — Folge-Slice
  `slice-046`; berührt eine SQL-View, nicht den Go-Store-Port.
- **Konfigurierbarkeit von `RetentionPolicy.MinAge` zur Laufzeit** (z. B.
  über `CDC_CONFIG_FILE`, `slice-041`) — Bestand bleibt bewusst stehen:
  `RetentionPolicy` wird weiterhin so konstruiert, wie es der aufrufende
  Code heute vorsieht; eine Laufzeit-Konfigurationsanbindung ist ein
  anderer Vorgang, den `welle-13` §6 bereits ausschließt.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] `ChangeStorePort` trägt eine neue Löschmethode; der
      `PostgresChangeStoreAdapter` implementiert sie real gegen
      PostgreSQL (`cdc.change`-Zeilen werden tatsächlich entfernt,
      `make test-store` real belegt).
- [x] `RunRetentionUseCase` (Application-Schicht) ruft für eine Quelle
      `RetentionPolicy.AllowsDeletion` je betrachtetem Change real auf
      (Alter, Change-Position, alle bestätigten Consumer-Positionen) und
      übergibt ausschließlich freigegebene Changes an die neue
      Port-Methode — mit Fake-Port-Test, der eine gemischte Menge
      (löschbar/nicht löschbar wegen Alter, löschbar/nicht löschbar wegen
      eines zurückhängenden Consumers) real unterscheidet.
- [x] `LH-FA-RET-002` (Negative: ungültige Konfiguration liefert
      expliziten Fehler statt stiller Übernahme) real erfüllt — Test für
      den Fehlerpfad.
- [x] `make gates` grün, `make test`/`make test-store` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Beleg: [`docs/reviews/review-slice-043.md`](../../../reviews/review-slice-043.md)
      (0 HIGH, 2 MEDIUM), Fixrunde in Commit `193c47b` (F-1 real
      behoben — Transaktions-Waisen-Bereinigung neu gebaut —, F-2
      Semantik bewusst beibehalten und begründet), bestätigt in
      [`docs/reviews/review-slice-043-fixrunde.md`](../../../reviews/review-slice-043-fixrunde.md)
      (dabei neuer Nebenbefund F-3 LOW — falsche `ADR-0029`-Zitierung in
      Prosa —, in Commit `3177b9c` behoben). Verifikation in
      [`docs/reviews/verify-slice-043.md`](../../../reviews/verify-slice-043.md)
      (DoD eigenständig nachgeprüft, keine Rückführung nötig).
- [x] Doku-Update: falls ein öffentlicher Vertrag entsteht (neue
      Port-Methode ist intern, kein CLI/SQL-Vertrag in diesem Slice) —
      Implementer prüft und begründet im Plan-Nachzug. **Geprüft:** kein
      öffentlicher Vertrag entsteht — `DeleteChanges`, `Positions` und
      `RunRetentionUseCase` sind interne Go-Schnittstellen ohne
      CLI-/SQL-Außenfläche; kein Doku-Update fällig (Details im
      Plan-Nachzug).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. Siehe §7.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Entfällt: Repo ist Greenfield, `../reconciliation.md` existiert nicht.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert. Siehe §7 — keine Beobachtung angefallen; `BEO-PGC/retention-keine-loeschausfuehrung` bleibt bei 0× bis zur `welle-13`-Closure (analog zu `welle-12`s Slices).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen). Siehe §6 — beide entfallen.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Repo mit Wellen-Betrieb (`welle-13` offen) — Prüfung läuft bei der `welle-13`-Closure.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/application/port/outbound/changestore.go` | update | neue Löschmethode am `ChangeStorePort` |
| `internal/adapters/driven/postgresstorage/*.go` | update | `PostgresChangeStoreAdapter` implementiert die Löschmethode real |
| `internal/application/usecase/retention/*.go` (neu) | neu | `RunRetentionUseCase`, ruft `AllowsDeletion` real auf |
| `internal/application/port/outbound/consumerposition.go` (falls nötig) | prüfen | ob ein bestehender Port bereits alle Consumer-Positionen einer Quelle liefert, oder eine neue Lesefähigkeit nötig ist |
| `*_test.go` | neu/update | Fake-Port-Tests (Use-Case), reale PostgreSQL-Tests (Adapter) |

### Plan-Nachzug (nach Implementierung)

Regeln dieser Sektion: Implementierungsentscheidungen, die über die Tabelle
oben hinausgehen — Namen, Signaturen und die Auflösung der beiden §6-Risiken.

**1. Name/Signatur der Löschmethode.** `ChangeStorePort.DeleteChanges(ctx
context.Context, changeIDs []model.ChangeID) error` — keine
`DeleteChangesBefore(position)`-Variante. Begründung: Die DoD verlangt, dass
`RunRetentionUseCase` `AllowsDeletion` **je betrachtetem Change** real
aufruft und **ausschließlich die freigegebene Menge** an den Port übergibt.
Eine explizite ID-Liste ist die wörtlichste Umsetzung dieses Vertrags — sie
verlangt keine Zusatz-Annahme über die Kontiguität der freigegebenen Menge
(die zwar aus der Policy-Struktur folgt, aber vom Use Case nicht behauptet
werden muss) und ist an der realen Test-Assertion
(`TestDeleteChangesRemovesOnlyGivenChanges`) direkt ablesbar: eine Teilmenge
geht, der Rest bleibt unangetastet, unabhängig von der Positions-Ordnung.
`DeleteChanges` ist idempotent (SQL `DELETE ... WHERE change_id = ANY($1)`);
eine leere Menge ist ein gültiger Aufruf ohne Datenbank-Rundlauf.

**2. „Alle bestätigten Consumer-Positionen einer Quelle" (§6, Risiko 2).**
Kein bestehender Port lieferte das — `ConsumerStatePort.Position` liest
genau einen Consumer über seine Kennung, es gab keine quellen-skopierte
Variante. Statt eines neuen Ports (`consumerposition.go` aus der
Plan-Tabelle) wurde `ConsumerStatePort` um eine vierte Methode erweitert:
`Positions(ctx context.Context, source model.SourceID)
([]model.ConsumerPosition, error)` — Begründung: Die Fähigkeit gehört
fachlich zum bestehenden `ConsumerStatePort` (`ARC-004`, Fähigkeits-Port je
`ADR-0034`), der bereits die einzige Konsistenzgrenze der Consumer-Zustände
trägt; ein zweiter Port für dieselbe Tabelle (`cdc.consumer_position`) hätte
keine eigene Konsistenzgrenze, nur eine zweite Adresse für dieselbe Zeile.
SQL: `SELECT consumer_id, acknowledged_position FROM
cdc.consumer_position WHERE source_id = $1` — neue Konstante
`SelectConsumerPositionsBySource`.

**Abwesenheits-Lesart (wichtige Nebenentscheidung):** `Positions` liest
ausschließlich Zeilen aus `cdc.consumer_position`. Ein registrierter, aber
noch nie gegen *irgendeine* Quelle bestätigender Consumer erscheint in
**keiner** Quelle und blockiert damit auch keine — dieselbe Lesart, die der
bestehende Code für `Position`/`Remove` bereits dokumentiert („die
Zeilen-Abwesenheit trägt der Rolle des Consumers in der Retention
Rechnung", `ConsumerStatePort.Remove`-Doku). Das ist konsistent, weil das
Domänenmodell einen Consumer erst über seine erste Bestätigung an eine
Quelle bindet (`internal/domain/model/consumer.go`,
`ConsumerPosition.Advance`: „eine bestätigte Quelle bleibt gebunden") —
vor der ersten Bestätigung besteht keine Quellen-Zuordnung, die
`Positions` melden könnte.

**3. Age-Berechnung — `ChangeRecord` trägt jetzt `CommittedAt`.**
`RetentionPolicy.AllowsDeletion` verlangt das Alter eines Changes; der
reale Quell-Commit-Zeitpunkt lag bereits in `cdc.transaction.committed_at`
(`LH-FA-ADM-004`), aber `ChangeRecord`/`ReadChanges` gaben ihn nicht an den
Aufrufer weiter. Ergänzt: `ChangeRecord.CommittedAt model.TimePoint`,
`SelectChanges` liest zusätzlich `t.committed_at`, `collectRecords` scannt
und mappt ihn. Rückwärtskompatibel (additive Struct-Erweiterung, einzige
bestehende `ChangeRecord{...}`-Konstruktion mit benannten Feldern in
`store.go`, per `grep` verifiziert). `RunRetentionService.Run` berechnet
`age := clock.Now().Sub(record.CommittedAt)` über den injizierten
`ClockPort` (`ADR-0040`) — Domain und Application rufen keine Systemzeit
direkt auf.

**4. `RunRetentionUseCase` als neue Inbound-Port-Kategorie.** Neue Datei
`internal/application/port/inbound/retention.go` (`RunRetentionCommand`,
`RunRetentionResult`, `RunRetentionUseCase`) statt Erweiterung einer
bestehenden Datei — `ARC-002` führt „Retention" bereits als eigene
Use-Case-Kategorie neben Capture/Consumer/Konfiguration
(Architect-Verdikt), dieselbe Eins-Datei-je-Kategorie-Konvention wie
`capture.go`/`consumer.go`/`verwaltung.go`.

**5. `RunRetentionCommand`-Validierung — Auflösung von `LH-FA-RET-002`
Negative auf dieser Schicht.** Die domänenseitige Negative
(`NewRetentionPolicy` weist eine negative `MinAge` zurück) bestand bereits
vor diesem Slice (`TestNewRetentionPolicyRejectsNegativeDuration`) und
bleibt unverändert. Auf der Use-Case-Schicht ist eine `RetentionPolicy`
immer schon gültig konstruiert (Wertobjekt ohne Nachträglich-Invalidierung)
— die einzige an dieser Schicht neu mögliche „ungültige Konfiguration" ist
ein `RunRetentionCommand` ohne Quelle. `Run` weist eine leere
`command.Source` über `domainerrors.ErrEmptyIdentifier` zurück, **bevor**
einer der drei Ports berührt wird (`TestRunRejectsEmptySource` belegt beide
Hälften: den Fehler und die Null-Berührung der Ports) — derselbe explizite
Fehlerpfad statt stiller Übernahme, den `LH-FA-RET-002` verlangt, nur auf
der Konfigurations-Fläche dieser Schicht statt der Policy-Konstruktion.

**6. Idempotenz-Risiko (§6, Risiko 1) — kein Konflikt, kein
Implementierungs-Zusatz nötig.** Das Architect-Verdikt hatte bereits
geprüft, dass `AllowsDeletion` als einzige Freigabe-Instanz einen
zeitlichen Vorrang von Persist-before-ACK/At-Least-Once vor jeder
physischen Löschung erzwingt. Die Implementierung bestätigt das strukturell:
`RunRetentionService.Run` löscht ausschließlich Changes, die
`ReadChanges` bereits als persistiert zurückgibt und für die
`AllowsDeletion` zusätzlich alle übergebenen Consumer-Positionen als
bestätigt und mindestens auf Höhe der Change-Position verlangt — ein
Bypass am Domain Core vorbei existiert nicht (`DeleteChanges` ist die
einzige Löschfähigkeit am Port, und sie wird ausschließlich aus der
freigegebenen Teilmenge heraus aufgerufen). Der Fall „gelöschte Transaktion muss erneut persistiert werden" tritt unter
dieser Reihenfolge nicht ein: Das Source-ACK einer Transaktion erfolgt
bereits im Capture-Pfad direkt nach ihrer Persistenz
(`LH-QA-REL-001.a`), unabhängig von der Retention — ein Crash-Replay
derselben Quelltransaktion (`ADR-0012`, At-Least-Once) kann also nur
auftreten, *bevor* die Quelle ihr Commit bestätigt bekam, also lange bevor
irgendein Consumer sie verarbeitet und `AllowsDeletion` sie je zur
physischen Löschung freigibt. Der Konflikt, den das Risiko benannte, setzt
eine Reihenfolge voraus, die am Domain Core (`AllowsDeletion` als einzige
Freigabe-Instanz) strukturell ausgeschlossen ist.

**7. Fixrunde nach `review-slice-043` — F-1 (verwaiste `cdc.transaction`-Zeilen):
behoben, nicht ausgeschlossen.** Realer Befund vor der Entscheidung:
`tools/schema/nacharbeit-observability.sql`s bereits ausgerollte
`cdc.metrics`-View berechnet `cdc_oldest_change_age_seconds` und
`cdc_transactions_total` direkt über `cdc.transaction` — ungefiltert nach
verbliebenen `cdc.change`-Zeilen. Eine verwaiste Transaktion (alle ihre
Changes bereits über `DeleteChanges` bereinigt) hätte
`cdc_oldest_change_age_seconds` weiterhin auf ihrem alten
`committed_at`-Wert gehalten, obwohl der tatsächlich älteste verbliebene
Change längst jünger ist — genau die „fälschlich als bestehende Aktivität"
gelesene Zeile, vor der die Fixrunden-Anweisung warnte. Kein
`ON DELETE CASCADE` hätte geholfen: Die Fremdschlüssel-Kante
(`change.transaction_id → transaction.transaction_id`) kaskadiert nur beim
Löschen der Transaktion, nicht beim Löschen ihrer Changes — die
Elternwaisen-Bereinigung braucht die umgekehrte Richtung. `DeleteChanges`
(`internal/adapters/driven/postgresstorage/store.go`,
`internal/adapters/driven/postgresstorage/queries/queries.go`) läuft jetzt
als eine DB-Transaktion: `DeleteChanges`-SQL mit `RETURNING transaction_id`
liefert die betroffenen Transaktions-Kennungen, eine zweite Abfrage
(`DeleteOrphanedTransactions`) entfernt daraus genau die, die keine
Change-Zeile mehr referenzieren. Scope bewusst eng auf die von diesem
Aufruf betroffene Menge — eine Transaktion, die von Geburt an nie eine
Change-Zeile trug, ist ein anderer Fall (keine Rolle der Retention) und
bleibt unberührt; dafür bräuchte es einen eigenen Slice, der diese Frage
erst aufwirft. Real belegt:
`TestDeleteChangesRemovesOrphanedTransactionOnly` (verwaiste Transaktion
verschwindet, weiterhin referenzierte bleibt stehen).

**8. Fixrunde nach `review-slice-043` — F-2 (Abwesenheits-Lesart von
`ConsumerStatePort.Positions`): bestätigt, mit dokumentierter Konsequenz.**
`LH-FA-RET-004` und `LH-FA-CON-005` vollständig gelesen: Keines der beiden
entscheidet die konkrete Konstellation „registriert, aber für **diese**
Quelle noch nie bestätigend" — beide Akzeptanzkriterien setzen einen
Consumer voraus, der bereits irgendeine Position für die Quelle trägt oder
zumindest *irgendwann* liest. Entscheidung: Die Abwesenheits-Lesart bleibt
unverändert. Begründung: `model.Consumer` (Registrierung) und
`model.ConsumerPosition` (Quellen-Bindung) sind im Domänenmodell bewusst
getrennt — ein Consumer bindet sich erst mit seiner ersten `Advance`
(`internal/domain/model/consumer.go`) an eine Quelle. Vor dieser ersten Bindung besteht
keine Quellen-Zuordnung, die `Positions` melden könnte, und ein
registrierter, aber quellen-seitig untätiger Consumer ist retentionsseitig
so zu behandeln, als hätte er noch nie erklärt, dass ihn *diese* Quelle
betrifft — dieselbe Lesart, die der bestehende `Remove`-Kommentar bereits
für die administrative Entfernung trägt, jetzt einheitlich auch für die
Erstbestätigung. Eine geänderte Semantik (jeder registrierte Consumer
blockiert jede Quelle bis zur ersten Bestätigung) würde eine
Quellen-Bindung bereits bei der Registrierung voraussetzen, die
`LH-FA-CON-001`/`RegisterConsumerService.Register` nicht kennt — das wäre
eine Spec-Erweiterung, kein Bugfix dieses Slice.
**Konsequenz für `slice-044`:** Ein frisch registrierter, aber gegen eine
Quelle noch nie bestätigender Consumer blockiert `RunRetentionUseCase.Run`
für diese Quelle **nicht** — ihre gesamte bisherige Historie kann gelöscht
werden, bevor dieser Consumer je gelesen hat. Schutz entsteht erst mit der
ersten `Acknowledge`-Bestätigung gegen diese Quelle (auch eine, die nur die
dokumentierte Anfangsposition aus `LH-FA-CON-005` Boundary bestätigt).
Operativ heißt das: ein neu angebundener Consumer, der vor seinem ersten
Lesezugriff Schutz vor Retention braucht, muss diese erste Bestätigung
setzen, *bevor* ein Retention-Lauf gegen dieselbe Quelle läuft — `slice-044`
(Hintergrund-Job/CLI-Trigger) dokumentiert diese Reihenfolge-Abhängigkeit
in seinem eigenen Plan, sie ist keine neue Erfindung dieses Slice, aber
erstmals operativ scharf.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `welle-13` liegt flach unter
`docs/plan/planning/` (eröffnet), `Verantwortlich:` gesetzt, WIP-Limit
(1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Zeigt sich,
  dass die Ermittlung „alle bestätigten Consumer-Positionen einer
  Quelle" einen bisher nicht existierenden, größeren Lesezugriffsweg
  braucht als erwartet, gehört das zurück zur Zerlegung (eigener
  Vorbereitungs-Slice).
- `in-progress` → `open` (blockiert — Carveout?): Kein bekannter Blocker.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD vollständig **und** `make gates` grün **und** `make test-store` grün
**und** Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Eine Löschoperation am `ChangeStorePort` könnte mit `ADR-0011`s
  Idempotenz-Pflicht (`PersistTransaction` ist deduplizierbar über die
  interne Transaktions-ID) kollidieren, wenn eine gelöschte Transaktion
  erneut persistiert werden müsste (Crash-Replay) — **Ausgang: entfallen.**
  Das Kollisionsfenster von `PersistTransaction`s Replay-Pflicht liegt
  strukturell vor der Quell-ACK (ein Crash-Restart wiederholt eine
  jüngst persistierte, noch nicht quellenseitig bestätigte Transaktion).
  `AllowsDeletion` gibt eine Transaktion dagegen erst frei, wenn sie das
  konfigurierte Mindestalter erreicht hat **und** alle betrachteten
  Consumer-Positionen sie bereits bestätigt haben — beide Fenster
  (Replay-nah vs. Retention-alt-und-bestätigt) überschneiden sich in
  keinem realen Ablauf.
- Es könnte bereits einen etablierten Lesezugriffsweg für „alle
  bestätigten Consumer-Positionen einer Quelle" geben (z. B. über
  `cdc.consumer_status`s Go-Pendant), oder er könnte fehlen und müsste
  neu gebaut werden — der Aufwand ist vor der Implementierung nicht
  exakt bekannt. **Ausgang: entfallen.** Real geprüft (Plan-Nachzug
  Punkt 2): kein bestehender Port lieferte das; `ConsumerStatePort`
  wurde um `Positions` erweitert, statt einen neuen Port anzulegen —
  begrenzter, klar begründeter Aufwand, kein unerwartet größerer
  Umbau nötig.

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** Der Implementer und der Reviewer trieben
  beide dieselbe Disziplin real bis zum Ende: der Implementer prüfte in
  der Fixrunde tatsächlich, ob eine verwaiste `cdc.transaction`-Zeile
  irgendwo gelesen wird, statt das Risiko nur als „theoretisch" abzutun
  — und fand einen echten, bereits ausgerollten Fehler in
  `cdc.metrics` (`cdc_oldest_change_age_seconds`/`cdc_transactions_total`
  hätten eine verwaiste Transaktion weiter mitgezählt). Der Reviewer
  reproduzierte den Fix eigenständig gegen echtes PostgreSQL und fand
  zusätzlich noch einen kleinen, eigenständigen Zitier-Fehler (F-3).
- **Was ging anders als geplant:** Der Reviewer fand 2 MEDIUM (verwaiste
  Transaktions-Zeilen, Consumer-Abwesenheits-Lesart gegen
  `LH-FA-RET-004` ungeprüft) — beide in der Fixrunde behoben bzw.
  bewusst und begründet beibehalten. Die Fixrunden-Bestätigung fand
  einen kleinen Nebenbefund F-3 (falsche `ADR-0029`-Regel-Zitierung in
  der Prosa des Plan-Nachzugs — die Consumer-Quellen-Bindung ist eine
  Domänenmodell-Eigenschaft, keine `ADR-0029`-Regel), bei der Closure
  direkt korrigiert.
- **Steering-Loop-Eintrag:** Kein Eintrag erreicht mit diesem Slice 3×.
- **Beobachtungs-Register (`../observations/`):** keine Beobachtung
  angefallen; `BEO-PGC/retention-keine-loeschausfuehrung` bleibt bei 0×
  bis zur `welle-13`-Closure.
- **Folge-Slices:** keine neuen — `slice-044`/`045`/`046` stehen bereits
  in `welle-13` §4 als vorgesehene nächste Slices.
- **Risiken aus §6:** beide *entfallen* — siehe §6.
- **Drei Paarungen:** Repo **mit** Wellen-Betrieb (`welle-13` offen) —
  Prüfung läuft bei der `welle-13`-Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** Einzige berührte Sub-Area ist die
Repo-weite Default-Sub-Area `*`/`PGC`.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen.
Treffer für `PGC`: `BEO-PGC/retention-keine-loeschausfuehrung` (0×,
benannt nicht gezählt — dieser Slice liefert den ersten und
architektonisch entscheidenden Baustein der Auflösung). Keiner der
übrigen Treffer erreicht mit diesem Slice 3×.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

Alle berührten Sub-Areas GF (nur `*`/`PGC`).
