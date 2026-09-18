# Review-Report: slice-043 — 2026-09-13

**Review-Art:** Code — geprüft gegen Plan (`slice-043`, §1/§2/§3/§4/§6/§8),
`welle-13`, den Architect-Verdikt zur Retention-Löschausführung und
`ADR-0009`/`0011`/`0012`/`0014` (alle `Accepted`, nur umgesetzt) sowie
`AGENTS.md` §3 Hard Rules (§3.1, §3.5, §3.7) — Rollentrennung Modul 8: diese
Prüfung läuft gegen Plan/ADR/Hard Rules (Maintainability), nicht gegen DoD
(Verifier-Aufgabe).

**Gegenstand:** Commit `de3ff8ca55dfca17b814b035bd9f0a64c87b09be`
(`feat(retention): ChangeStorePort-Löschmethode und RunRetentionUseCase
(LH-FA-RET-002, ADR-0014)`) — neu: `internal/application/usecase/retention/`
(`service.go`, `service_test.go`), `internal/application/port/inbound/retention.go`;
geändert: `internal/application/port/outbound/changestore.go`,
`internal/application/port/outbound/consumerstate.go`,
`internal/adapters/driven/postgresstorage/{store,consumerstate,queries/queries}.go`
und deren Tests, fünf Fake-Nachzieh-Patches in bestehenden Use-Case-Tests
(`acknowledge`, `capture`, `position`, `register`, `remove`,
`internal/bootstrap/heartbeat_internal_test.go`), Plan-Nachzug in
`slice-043-changestoreport-loeschmethode.md`.

**Skill:** `.harness/skills/reviewer.md` @ HEAD (Accepted, geschärft 2026-09-09)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-13

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-043-changestoreport-loeschmethode.md`
  (vollständig: §1 Ziel/Abgrenzung, §2 DoD, §3 Plan + Plan-Nachzug, §4
  Trigger, §6 Risiken, §8)
- `docs/plan/planning/welle-13.md` (vollständig — Welle-Ziel,
  Closure-Trigger, Out-of-Scope, Abhängigkeiten)
- der Architect-Verdikt zur Retention-Löschausführung
  (vollständig — bindende Architektur-Grundlage dieser Welle)
- `AGENTS.md` §3.1 (Docker-only), §3.5 (ADR-Immutabilität), §3.7
  (Kommentar-Disziplin — `BEO-PGC/slice-chronik-in-code-kommentar`, 2× vor
  diesem Lauf, besonders geprüft)
- `docs/plan/adr/0009-change-store-outbound-port.md`,
  `0011-persist-before-ack.md`, `0012-at-least-once.md`,
  `0014-retention-domain-policy.md`, `0029-domain-invarianten.md`,
  `0034-ports-nach-faehigkeiten.md` (vollständig)
- `internal/domain/model/retention.go` (`RetentionPolicy.AllowsDeletion`,
  unverändert in diesem Diff — Prüfgrundlage, nicht Gegenstand)
- `spec/lastenheft.md` — `LH-FA-RET-002`…`006`, `LH-FA-CON-005` (vollständig
  gelesen, nicht nur zitiert)
- `tools/schema/schema.yaml` — `cdc.change`, `cdc.transaction`,
  `cdc.consumer_position` (FK-Kanten, Indizes)
- Vollständiger `git show de3ff8c` (alle 17 Dateien, nicht nur die
  Implementer-Zusammenfassung)
- Review zu `slice-039` (Format-Vorlage)

---

## Findings

### F-1 — `cdc.transaction`-Zeilen bleiben nach `DeleteChanges` verwaist zurück

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Out-of-Scope-Disziplin, Baseline-Regelwerk
  `modul-05-planning-harness.md` §Ziel-Form: Slice — „was nicht ausdrücklich
  ausgeschlossen ist, wandert im Zweifel hinein“)
- `pfad`: `internal/adapters/driven/postgresstorage/store.go:169-187`
  (`DeleteChanges`) i.V.m. `tools/schema/schema.yaml:103-126` (`transaction`)
  und `:127-157` (`change`, FK `transaction_id → transaction.transaction_id`,
  keine `ON DELETE CASCADE`-Angabe)
- `befund`: `DeleteChanges` entfernt ausschließlich `cdc.change`-Zeilen. Die
  zugehörige `cdc.transaction`-Zeile bleibt in jedem Fall bestehen — auch
  wenn nach der Löschung keine einzige `cdc.change`-Zeile mehr auf sie
  verweist. `cdc.transaction` wächst damit unabhängig von jeder
  Retention-Ausführung monoton mit jeder committeten Quelltransaktion.
  Weder `slice-043`s §1 noch `welle-13`s §6 (Out-of-Scope) nennen die
  `transaction`-Tabelle als bewusst ausgeschlossen — die Frage „wird die
  Elternzeile mitbereinigt, oder ist das ein anderer Vorgang?“ ist an
  keiner Stelle beantwortet, weder mit Kennung (Folge-Slice) noch mit
  Begründung (Bestand bleibt bewusst stehen). `ADR-0014` selbst spricht
  durchgängig von „Changes“, nie von „Transaktionen“, was diese Lücke
  stützt, aber nicht schließt — die Architect-Verdikt-Prüfung (Frage 1)
  untersucht ebenfalls nur die `AllowsDeletion`/`ChangeStorePort`-Seite,
  nicht das Transaktions-Wachstum. Für `LH-FA-RET-006`
  („Kontrolle des Datenwachstums“) ist das keine DoD-Verletzung dieses
  Slice (dessen Metrik `cdc_storage_bytes` laut Architect-Verdikt ohnehin
  nur `pg_relation_size('cdc.change')` misst — `cdc.transaction`s Wachstum
  bliebe damit auch für die künftige Metrik unsichtbar), aber ein
  unbenannter Rest-Wachstumsvektor, den Retention insgesamt nicht
  vollständig „kontrolliert“, obwohl kein Plan das je ausgeschlossen hat.
- `verifizierbar`: ja — `grep -n "references" tools/schema/schema.yaml`
  zeigt die FK ohne Cascade-Angabe; `TestDeleteChangesIsIdempotent`/
  `TestDeleteChangesRemovesOnlyGivenChanges` prüfen ausschließlich
  `cdc.change`-Zeilenzahlen, keine gegen `cdc.transaction`.
- `klasse`: „Out-of-Scope-Lücke ohne Adresse“ (verwandt, aber nicht
  identisch mit den vier etablierten Ausschluss-Klassen — hier fehlt die
  Frage selbst, nicht nur ihre Antwort; erstes Auftreten dieser Klasse)

### F-2 — `ConsumerStatePort.Positions`s Abwesenheits-Lesart lässt einen registrierten, aber noch nie gegen die Quelle bestätigenden Consumer retentionsseitig unsichtbar werden

- `kategorie`: MEDIUM
- `quelle`: `LH-FA-RET-004` (Consumer-basierte Retention) — Maintainability/
  Spec-Interpretationsrisiko, kein ADR-Verstoß
- `pfad`: `internal/application/port/outbound/consumerstate.go:82-89`
  (`Positions`-Dokumentation), `internal/adapters/driven/postgresstorage/consumerstate.go:193-236`
  (`Positions`-Implementierung), `internal/application/usecase/retention/service.go:60-73`
  (`Run` befragt ausschließlich die von `Positions` gelieferte Liste)
- `befund`: `model.Consumer` bindet sich laut Domänenmodell erst mit der
  ersten `Acknowledge`-Bestätigung an eine Quelle
  (`ConsumerPosition.Advance`, `ADR-0029` Regel 2) — die Registrierung
  selbst (`RegisterConsumerService.Register` → `InsertConsumer`) trägt
  keine Quellenzuordnung. Ein Consumer, der registriert ist, aber gegen
  eine bestimmte Quelle noch **kein einziges** Mal bestätigt hat, erscheint
  in `Positions(source)` nicht und blockiert `AllowsDeletion` für diese
  Quelle folglich überhaupt nicht (`AllowsDeletion`: „Ohne übergebene
  Positionen blockiert keine Position die Bereinigung“). `LH-FA-RET-004`s
  Happy Path („`c1` hat `p` noch nicht bestätigt … Changes ab `p` werden
  nicht entfernt“) und `LH-FA-CON-005`s Boundary („`c` hatte nie bestätigt
  … startet an einer definierten Anfangsposition“) beschreiben beide einen
  Consumer, der *irgendwann* liest — keines der beiden Akzeptanzkriterien
  trifft explizit die Konstellation „registriert, aber für **diese**
  Quelle bislang komplett unbeteiligt“ oder entscheidet, ob ein solcher
  Consumer „relevant“ im Sinne von `LH-FA-RET-004` ist. Diese
  Abwesenheits-Lesart ist **nicht neu** — sie war bereits vor diesem Slice
  am `Remove`-Dokumentationskommentar festgehalten und vom Architect-Verdikt
  nicht beanstandet —, aber `slice-043` ist der erste Punkt, an dem sie
  operative Konsequenz bekommt: `DeleteChanges` führt jetzt real aus, was
  vorher nur eine unbenutzte Domänenregel war. Praktisch folgenlos ist das
  aktuell noch, weil kein Aufrufer existiert (`slice-044` liefert ihn
  erst) — das Fenster, das offen bleibt, ist real, aber noch nicht scharf
  gestellt.
- `verifizierbar`: nein — das ist keine Gate-prüfbare Aussage, sondern eine
  Spec-Interpretationsfrage; ein Test kann das gewählte Verhalten belegen
  (`TestPositionsWithoutAcknowledgementsCarriesEmptySlice` tut das bereits),
  aber nicht entscheiden, ob es das *richtige* ist.
- `klasse`: „Consumer-Sichtbarkeitslücke vor Erstbestätigung“ (erstes
  Auftreten dieser Klasse)

## Design-Entscheidungen — Verdikt zu den drei selbst geflaggten Punkten

**1. `DeleteChanges(ctx, changeIDs []model.ChangeID) error` statt
Positions-Cutoff — bestätigt korrekt.** Konsistent mit dem DoD-Wortlaut
„ausschließlich freigegebene Changes“: Eine explizite ID-Liste macht keine
Zusatzannahme über Kontiguität der freigegebenen Menge, die
`AllowsDeletion` nicht liefert (die Boundary-Test-Fälle
`c-too-young`/`c-consumer-behind` liegen zwischen zwei löschbaren Changes —
ein Cutoff-Parameter könnte das nicht abbilden, ohne selbst wieder eine
Menge zu berechnen). Die SQL `DELETE FROM cdc.change WHERE change_id =
ANY($1)` ist korrekt (setbasiert, keine Sequenzabhängigkeit) und real
idempotent belegt (`TestDeleteChangesIsIdempotent`: bereits gelöschte und
nie vorhandene Kennungen bleiben wirkungslos, keine leere Menge löst einen
DB-Aufruf aus). Kein Konflikt mit `ADR-0011`s Idempotenz-Pflicht — jene
betrifft `PersistTransaction`, nicht `DeleteChanges`; die beiden Operationen
sind unabhängige Vorgänge auf unterschiedlichen Zeilen.

**2. `ConsumerStatePort.Positions`-Erweiterung statt neuer Port —
bestätigt korrekt nach `ADR-0034`.** `ADR-0034`s eigene Folgepflicht prüft
neue Ports genau gegen die Frage, „ob sie eine Konsistenzgrenze komplett
abbilden“ — ein zweiter Port für dieselbe Tabelle (`cdc.consumer_position`)
hätte keine eigene Konsistenzgrenze, nur eine zweite Adresse für dieselbe
Zeile. Die Erweiterung ist damit die architekturell vorgesehene Antwort,
keine Verletzung von Interface-Segregation im Sinne dieses Repos (die
Fähigkeits-Port-Wahl aus `ADR-0034` nimmt diese Abwägung bereits vorweg).
Die inhaltliche Frage — ob die *Abwesenheits-Lesart selbst* für
`LH-FA-RET-004` korrekt ist — ist von der Port-Platzierung unabhängig und
oben als F-2 gesondert geführt.

**3. `LH-FA-RET-002`-Negative auf zwei Ebenen — keine Duplikation,
zwei unterschiedliche Konfigurationsflächen.** Die domänenseitige Prüfung
(`NewRetentionPolicy` verwirft negative `MinAge`) und die
use-case-seitige Prüfung (`Run` verwirft leere `Source`) prüfen
unterschiedliche Eingaben an unterschiedlichen Schichten — eine
`RetentionPolicy` ist als Wertobjekt immer schon gültig konstruiert, bevor
sie den Use Case erreicht, die einzige an dieser Schicht neu mögliche
Ungültigkeit ist die fehlende Quelle. Das ist keine Wiederholung derselben
Prüfung, sondern zwei verschiedene Instanzen derselben `LH-FA-RET-002`-
Negative-Anforderung an zwei verschiedenen Eingaben.

## Negativbefunde

- geprüft, ohne Befund: **Kein ADR verändert (`AGENTS.md` §3.5).**
  `git show de3ff8c --name-only` trifft keine Datei unter
  `docs/plan/adr/`; `0009`/`0011`/`0012`/`0014` bleiben `Accepted` und
  inhaltlich unverändert — konsistent mit dem Architect-Verdikt „alle nur
  umgesetzt“.
- geprüft, ohne Befund: **Kommentar-Disziplin (`AGENTS.md` §3.7), keine
  Slice-/Wellen-Chronik in den geänderten Go-Dateien.** Vollständiger
  `grep` über alle 17 in `de3ff8c` geänderten/neuen Dateien nach
  `slice-\d+|welle-\d+|früher|vorherig|ehemals|zuvor stand|wäre gewesen|
  hätte` liefert nur Treffer, die bereits vor diesem Diff bestanden
  (`consumerstate.go:131`, `consumerstate_test.go:181/199`,
  `outbound/consumerstate.go:43`, alle „frühere Position“ im Sinne einer
  Domäneninvariante, nicht einer Slice-Chronik) und einen Treffer im
  Plan-Dokument selbst (Prosa, keine Norm). Kein neuer Treffer in
  `*.go`-Dateien — `BEO-PGC/slice-chronik-in-code-kommentar` bleibt bei
  2×, erreicht mit diesem Slice **nicht** die 3×-Schwelle.
- geprüft, ohne Befund: **`ChangeRecord.CommittedAt`-Erweiterung ist real
  additiv.** `grep -rn "ChangeRecord{" internal/` (ohne Tests) trifft nur
  die eine benannte Konstruktion in `store.go:244` — keine positionale
  Struct-Initialisierung, die durch das neue Feld brechen würde.
- geprüft, ohne Befund: **Docker-only (`AGENTS.md` §3.1).** Der Diff fügt
  keine Build-/Test-Skripte hinzu, die eine lokale Toolchain installieren;
  alle neuen Tests laufen über die bestehende `make test`/`make
  test-store`-Kette.
- geprüft, ohne Befund: **`RunRetentionService.Run` respektiert Persist-
  before-ACK/At-Least-Once strukturell.** `DeleteChanges` ist die einzige
  Löschfähigkeit am `ChangeStorePort`; sie wird ausschließlich mit der über
  `AllowsDeletion` freigegebenen Teilmenge aufgerufen (`service.go:60-73`)
  — kein Bypass am Domain Core vorbei, wie vom Plan-Nachzug §3 Punkt 6
  behauptet und hier am Code nachvollzogen.
- geprüft, ohne Befund: **Rot-Kontrollproben-Testnamen sind plausibel und
  strukturell konsistent mit ihrer Behauptung.**
  `TestRunDistinguishesEligibleChangesFromMixedSet` prüft real eine
  gemischte Vier-Elemente-Menge (löschbar/zu jung/Consumer zurück/Boundary)
  gegen eine konkrete `deletedIDs`-Reihenfolge;
  `TestDeleteChangesRemovesOnlyGivenChanges`/`TestDeleteChangesIsIdempotent`
  laufen gegen eine reale PostgreSQL-Instanz (`newTestStore`,
  `seedReference`) und zählen tatsächliche Zeilen; `TestRunRejectsEmptySource`
  belegt sowohl den Fehler als auch die Null-Berührung aller drei Ports.
  Keine dieser Test-Strukturen deutet auf einen nachträglich grün
  geschriebenen statt vorher rot gelaufenen Test hin.
- geprüft, ohne Befund: **§1 Out-of-Scope, Klassenzuordnung der vier
  genannten Punkte korrekt** (Folge-Slice `slice-044`/`045`/`046` mit
  Kennung, Bestand `slice-041`-Konfigurierbarkeit mit Begründung) — mit der
  in F-1 benannten Lücke als einzige unbenannte Ausnahme.
- geprüft, ohne Befund: **§8 Sub-Area-Prüfung.** Einzige berührte Sub-Area
  `*`/`PGC`, korrekt GF; Beobachtungs-Register-Sichtung dokumentiert
  (`BEO-PGC/retention-keine-loeschausfuehrung`, 0×, benannt nicht gezählt,
  nachvollziehbar als „dieser Slice liefert den ersten Baustein“).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 2 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** „Out-of-Scope-Lücke ohne Adresse“ (1×,
erstes Auftreten, F-1) · „Consumer-Sichtbarkeitslücke vor Erstbestätigung“
(1×, erstes Auftreten, F-2) — kein Steering-Loop-Eintrag fällig, keine
Klasse erreicht 3×.

## Verdikt

**Merge-blockierend:** nein — der Commit ist bereits gepusht, keines der
beiden MEDIUM-Findings trägt einen Rollen-Widerspruch (beide sind
Beobachtungen gegen eine Lücke, nicht gegen eine bestrittene
Implementer-Aussage), und keines betrifft den kritischen Pfad *in diesem
Slice selbst* — `DeleteChanges` hat noch keinen Produktions-Aufrufer
(`slice-044` fehlt), beide Findings werden erst mit dessen realer
Verdrahtung operativ scharf.

**Zu den drei selbst geflaggten Design-Entscheidungen:** (1) ID-Listen-Signatur
statt Cutoff — bestätigt korrekt, SQL idempotent und richtig. (2)
Port-Erweiterung statt neuem Port — bestätigt korrekt nach `ADR-0034`; die
*inhaltliche* Abwesenheits-Lesart bleibt als F-2 offen, unabhängig von der
Port-Platzierung. (3) Zwei-Ebenen-Negative — keine Duplikation, zwei
unterschiedliche Konfigurationsflächen derselben Anforderung.

**Übergabe:** Zwei MEDIUM-Findings ohne Rollen-Widerspruch — Implementer
entscheidet über Annahme (Folge-Slice/Beobachtungs-Register-Eintrag) oder
Begründung vor der Closure; keine Architect-Sequenz nach Modul 8
erforderlich (kein HIGH, kein Widerspruch, keine dritte Wiederholung).
Empfehlung: F-1 und F-2 vor oder mit `slice-044` (der erste reale Aufrufer
von `RunRetentionUseCase`) explizit auflösen — entweder als benannter
Ausschluss mit Begründung, als Folge-Slice-Adresse, oder über einen
Architect-Zug, der die offene Interpretationsfrage (F-2) entscheidet. Für
`slice-043`s eigene DoD-Punkte ist das kein Mangel — DoD-/Spec-Konformität
prüft der Verifier separat; dieser Report ist ein Lauf-Beleg und wird über
Läufe hinweg nicht wieder gelesen, die Summary-Zeile speist bei Bedarf den
Closure-Eintrag (Modul 5).
