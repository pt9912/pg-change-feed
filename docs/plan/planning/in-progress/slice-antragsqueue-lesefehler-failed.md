# Slice antragsqueue-lesefehler-failed: Die Lesung der Antrags-Queue lehnt keine Zeile ab — eine vom Konstruktor verworfene Zeile endet `failed`, die Queue läuft weiter

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — der Slice trägt keine Closure-Bedingung, die von
seiner DoD verschieden wäre. Er geht `slice-transformationen-start-reihenfolge`
voraus (Start-Trigger dort, §4;
[welle-transformationen](../welle-transformationen.md) §5): die Antragsweg-Fläche
der Welle liest dieselbe Queue.

**Bezug:** [`LH-FA-ADM-001`](../../../../spec/lastenheft.md) (schreibende
SQL-Administration), [`LH-FA-CFG-005`](../../../../spec/lastenheft.md)
(Spaltenausschluss) und [`LH-FA-CFG-007`](../../../../spec/lastenheft.md)
(Transformationen) als die Antragsarten, deren Zeilen die Queue trägt,
[`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(Antrags-Queue),
[`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md) (keine
Domänenlogik in SQL),
[`ADR-0029`](../../adr/0029-domain-invarianten.md) (Konstruktor-Invarianten),
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
(Prüfung der Regelfelder in der Verarbeitung, nicht beim Lesen);
Entscheidung des Architect-Zugs im Beobachtungs-Register:
`BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue` (`state.md`, Option
„allgemeine Form“).

**Berührte Spec-Stellen:** [`SPEC-019`](../../../../spec/pflichtenheft.md)
(`cdc.administration_request`: die Zeilen zu `column_name` und zu den
Transformations-Antragsarten, die Fehlertext-Tabelle).

**Verantwortlich:** Implementer-Agent.

**Autor:** Planner-Agent, Auftrag des Auftraggebers (Architect-Zug zu
`BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue`). **Datum:** 2026-09-26.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Lesung der Antrags-Queue lehnt keine einzelne Zeile ab: eine Zeile,
die der Antrags-Konstruktor verwirft (leeres Schema, leerer Tabellenname,
`exclude_column`/`include_column` mit leerer Spalte, jeder künftige Grund), wird
als `failed` mit Kennung und Fehlertext vermerkt, und die Zeilen dahinter werden
verarbeitet.

**Ausgangslage — gelesen, Stand `7b70b34a`.** `sqlexec.ReadPendingRequests`
(`internal/adapters/driven/postgresstorage/sqlexec/translate.go`) gibt beim
ersten Konstruktor-Fehler den Fehler zurück; `processAdministrationRequests`
(`internal/bootstrap/wiring.go`) protokolliert ihn und kehrt zurück; der nächste
Durchlauf liest dieselbe Zeile erneut. Kein Antrag der Queue, auch keiner einer
anderen Art, wird verarbeitet, bis die Zeile entfernt oder ihr Status per
`UPDATE` geändert ist. Erprobt ist das Verhalten des Lesepfads
(`TestReadPendingRequestsRejectsRowWithEmptySchemaOrTable`); dass die
SQL-Funktionen eine solche Zeile schreiben, ist **hergeleitet, nicht erprobt**.
Die Regelfelder (`rule_name`, `rule_spec`) der beiden Transformations-Antragsarten
sind bereits so behandelt: der Konstruktor lässt sie durch, der Use Case prüft
und lehnt mit dem Fehlertext der Spec ab. Dieser Slice schließt den Rest und
jeden künftigen Ablehnungsgrund mit.

**Warum vor `slice-transformationen-start-reihenfolge`.** Der Abhilfe-Weg der
Transformationen (Antrag `cdc.remove_transformation`,
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Folgepflicht 5) läuft über dieselbe Queue: eine Zeile, die die Lesung anhält,
hielte den Abhilfe-Antrag hinter ihr an (**hergeleitet** aus
`processAdministrationRequests`, nicht erprobt) — das Kriterium „der Antrag ist
`applied`, bevor die erste Transaktion der Tabelle assembliert wird“ trüge nicht.

**Entscheidung (Architect-Zug, `BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue`,
`state.md`):** Option (a) in der allgemeinen Form. (b) scheidet aus: eine Prüfung
in den SQL-Funktionen berührt
[`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md) (keine
Domänenlogik in SQL) und schützt nicht gegen einen direkten `INSERT` von
`cdc_admin`. (c) scheidet aus: eine Zeile einer vertrauten Rolle hielte den
Betrieb an, die Ursache stünde nur im Log. Der Fix ist Code (Lesepfad und Use
Case), keine ADR.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Eine Prüfung in den SQL-Funktionen** (Option (b)):
  [`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md), siehe
  oben. Die SQL-Funktionen, das Schema und die Grants bleiben unverändert; es
  läuft kein `make schema-rollout` für diesen Slice.
- **Die Prüfung der Regelfelder und K1 bis K4.** Sie liegen in der Verarbeitung
  und bleiben, wie sie sind
  ([`SPEC-019`](../../../../spec/pflichtenheft.md)); der Slice ändert nur, was mit
  einer vom **Konstruktor** verworfenen Zeile geschieht.
- **Ein Recovery-Weg für Schema-Fehler**
  (`BEO-PGC/kein-admin-weg-schema-fehler-recovery`): eine andere Ursache
  (Abbruch des Erfassungspfads), nicht die der Queue.
- **Eine Zeile ohne Kennung** (`administration_request_id` leer): sie lässt sich
  nicht `failed` vermerken (`MarkFailed` verlangt eine Kennung); die Lesung
  überspringt sie mit Warnung und hält die Queue nicht an. Benannte Grenze
  (§6), die der Text von `SPEC-019` nennt; ob eine solche Zeile entstehen kann,
  ist hergeleitet (Primärschlüssel, von der SQL-Funktion vergeben), nicht erprobt.

## 2. Definition of Done

- [x] **Liefer-Punkt 1 — Lesepfad und Verarbeitung.** Eine vom
      Antrags-Konstruktor verworfene Zeile (leeres Schema, leerer Tabellenname,
      `exclude_column`/`include_column` mit leerer Spalte, eine Antragsart
      außerhalb der geschlossenen Menge, jeder künftige Grund) wird mit ihrer
      Kennung und dem Fehlertext durchgereicht und in der Verarbeitung `failed`
      vermerkt (`error_message` trägt den Text); die Zeilen dahinter werden in
      der Ordnung der Queue verarbeitet. `ListPending` gibt für eine verworfene
      Zeile keinen Fehler mehr zurück, nur für einen Fehler der Lesung selbst
      (Anfrage, Scan, Iteration). Eine Zeile ohne Kennung wird mit Warnung
      übersprungen (Grenze, §6). Die Form der Durchreichung ist am Start zu
      entscheiden: Empfehlung ein Rückgabetyp, der je Zeile Antrag oder Verwurf
      trägt; das ändert den Port `outbound.AdministrationRequestPort`, den
      Adapter, `sqlexec` und die Fakes in `internal/bootstrap` (die Stellen
      nennt der Suchlauf in §3). *Zu belegen durch:* ein Store-Test mit realen
      Zeilen (`make test-store`, im Paket `postgresstorage` mit `-run` der neuen
      und geänderten Tests, danach der Lauf des Pakets): leeres Schema, leerer
      Tabellenname, `exclude_column` und `include_column` mit leerer Spalte, je
      gefolgt von einer gültigen Zeile, die `applied` endet — die **Bereinigung
      des Tests löscht nur seine eigenen Zeilen** (Kennungen des Tests, nie die
      Tabelle: eine `pending`-Zeile eines abgebrochenen Tests dürfte sonst andere
      Tests anhalten, `BEO-PGC/test-isolation-geteilter-zustand`); ein
      Whitebox-Test der Verarbeitung im Paket `internal/bootstrap` (Fakes): die
      verworfene Zeile endet `failed` mit Text, die Zeile dahinter wird
      verarbeitet, ein Fehler von `MarkFailed` wird protokolliert und bricht den
      Durchlauf nicht ab; der Login-Test
      `TestAdministrationPathRunsUnderLeastPrivilegeLogins` (Paket
      `internal/bootstrap`) führt eine verworfene Zeile unter `cdc_admin` bis
      `failed`; `TestReadPendingRequestsRejectsRowWithEmptySchemaOrTable`
      ist auf das neue Verhalten umgeschrieben (Name und Aussage tragen
      „durchgereicht“, nicht „abgelehnt“). Die Mutationen der Eingabeseite je
      einmal rot gesehen: der Konstruktor-Fehler wird wieder zurückgegeben statt
      durchgereicht, die Schleife bricht an der verworfenen Zeile ab, der
      Fehlertext des Vermerks wird durch einen festen Text ersetzt (der Test
      prüft den Text je Grund); ein Fehler der Lesung selbst
      (`TestReadPendingRequestsClassifiesScanFailure`) bleibt ein Fehler.
- [x] **Liefer-Punkt 2 — die Spec.** [`SPEC-019`](../../../../spec/pflichtenheft.md)
      trägt im **selben Commit** wie der Code: der Ort der Prüfung (Verarbeitung,
      nicht Lesen), der Fehlertext je Grund (Klartext, gefolgt von Doppelpunkt,
      Leerzeichen und der Adresse — die Form der bestehenden Tabelle; die
      Adresse ist die Antrags-Kennung, weil Schema, Tabelle oder Spalte leer
      sind), die benannte Grenze „Zeile ohne Kennung“; die Nachbarzeilen sind
      gelesen und mitgezogen: die Spaltenzeile `column_name`
      („Domänen-Invariante des Antrags-Konstruktors“), die Zeilen
      `schema_name`/`table_name` und der Absatz „Transformations-Antragsarten“
      („die Queue lehnt … nicht beim Lesen ab“); keine ADR- und keine
      Slice-Bezüge in der Spec ([`AGENTS.md`](../../../../AGENTS.md) §3.4); die
      Änderungshistorie des Pflichtenhefts trägt eine Zeile. *Zu belegen durch:*
      Lesen von `SPEC-019` gegen den Code — der Text je Grund im Test ist
      wörtlich der der Spec — und der Suchlauf in §3; `make docs-check`.
- [x] **Liefer-Punkt 3 — Kommentare und Träger des alten Verhaltens.** Die
      Kommentare, die die Ablehnung beim Lesen beschreiben (Konstruktor,
      `ReadPendingRequests`, `processAdministrationRequests`, Test-Godocs), sind
      nachgezogen; ein Träger in einer fremden Datei wird gemeldet, nicht still
      mitgeändert. *Zu belegen durch:* das Suchlauf-Feld in §3 (beide Stände).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=<Plan-Datei>` läuft nach jeder Fixrunde
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: [`SPEC-019`](../../../../spec/pflichtenheft.md)
      (Liefer-Punkt 2); das Benutzerhandbuch trägt zur Queue keine Aussage über
      eine Ablehnung beim Lesen (Suchlauf in §3, Befund am Start prüfen); die
      Handbuch-Versionshistorie bleibt unberührt, solange das Handbuch nicht
      geändert wird.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — der Ausgang
      von `BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue` (`state.md`:
      Rest geschlossen, Träger benannt) und eine weitere `evidence/`-Datei, falls
      ein Auftreten anfällt; kein Anfall ist ebenfalls eine Antwort und wird in
      §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Slice-Closure selbst (der Slice hat keine Welle; das Ereignis kann
      eintreten).

**Umfang:** S bis M — Schätzung, nicht gemessen: ein Rückgabetyp über Port,
Adapter, `sqlexec`, Verarbeitung und die Fakes von 8 Nicht-Test- und bis zu 73
Fundstellen des Symbols (gemessen, Suchlauf in §3, Stand `7b70b34a`), ein
Store-Test, ein Whitebox-Test, ein Spec-Absatz.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/application/port/outbound/administrationrequest.go` | update | `ListPending` trägt je Zeile Antrag oder Verwurf (Kennung, Fehlertext); der Port-Kommentar nennt es. |
| `internal/adapters/driven/postgresstorage/sqlexec/translate.go` | update | `ReadPendingRequests`: der Konstruktor-Fehler wird durchgereicht statt zurückgegeben; die Ordnung der Zeilen bleibt; Kommentar. |
| `internal/adapters/driven/postgresstorage/administrationrequest.go` | update | `ListPending` folgt der neuen Form. |
| `internal/bootstrap/wiring.go` | update | `processAdministrationRequests`: eine verworfene Zeile wird `failed` vermerkt (`MarkFailed` mit dem Text), die Zeilen dahinter laufen; eine Zeile ohne Kennung wird mit Warnung übersprungen; Kommentare. |
| `internal/domain/model/administrationrequest.go` | update | der Konstruktor bleibt; sein Kommentar („keine Ablehnung beim Lesen der Queue — ein abgelehnter Lesevorgang hielte jeden Antrag dahinter an“) beschreibt den alten Zustand und wird nachgezogen. |
| `spec/pflichtenheft.md` `SPEC-019` | update | Ort der Prüfung, Fehlertext je Grund, Grenze „Zeile ohne Kennung“, Änderungshistorie; im selben Commit wie der Code. |
| `internal/adapters/driven/postgresstorage/sqlexec/translate_test.go` | update | `TestReadPendingRequestsRejectsRowWithEmptySchemaOrTable` auf das neue Verhalten (Negative nach [`LH-FA-ADM-001`](../../../../spec/lastenheft.md)); neue Fälle je Grund und die Zeile dahinter. |
| `internal/adapters/driven/postgresstorage/administrationrequest_test.go` | update | Store-Test mit realen Zeilen (`make test-store`): je Grund eine verworfene Zeile mit gültiger Zeile dahinter; skopierte Bereinigung; die Mutations-Kommentare an zwei Tests (Zeilen 757, 1262) folgen dem neuen Ausgang. |
| `internal/bootstrap/administration_internal_test.go` und die übrigen Aufrufer von `ListPending` | update | Fake-Port folgt der Form; Whitebox-Test der Verarbeitung; der Kommentar an Zeile 940 („beim Lesen abgelehnt; hier bilden das die Fakes nicht nach“). |
| `harness/README.md` | prüfen | Beschreibung von `make test-store`: sie nennt die Queue nicht (Zeile 4 des Feldes: keine Fundstelle in `harness`); keine Änderung erwartet. **Geliefert:** nicht berührt, Suchlauf Zeile 8. |
| `internal/application/port/outbound/administrationrequest.go` — Form der Durchreichung | geliefert | Entscheidung (kleinste tragende Form): `ListPending` liefert `[]PendingAdministrationRequest`; eine Zeile trägt entweder den Antrag (`Request`) oder `Rejected` (`RejectedAdministrationRequest`: `ID` und `Message`). Der Konstruktor bleibt unverändert; den Fehlertext bildet `sqlexec.rejectionMessage` aus dem Konstruktor-Fehler und den Feldern der Zeile (Klartext, Doppelpunkt, Leerzeichen, Antrags-Kennung), weil dort Zeile und Fehler zusammentreffen. Eine Zeile ohne Kennung reist als `Rejected` mit leerer `ID` bis in die Verarbeitung und wird dort mit Warnung übersprungen (Ort der Warnung: das Log-Port liegt in der Verarbeitung, nicht in `sqlexec`). Die Ports-Regel (`ADR-0034`, `ADR-0028`) bleibt: der Port trägt Domänentypen und Text, keinen Treibertyp. |
| `internal/adapters/driven/postgresstorage/sqlexec/rejection_internal_test.go` | neu | Test der Auffangzeile `Antrag ist ungültig` von `rejectionMessage` (über die Zeilen der Abfrage nicht herstellbar; Grund: der interne Test ruft die Funktion mit einem Fehler außerhalb der Konstruktor-Gründe). |
| `internal/bootstrap/administration_endtoend_test.go`, `internal/adapters/driven/postgresstorage/administrationrequest_order_test.go` | update | Aufrufer von `ListPending` folgen der Form (Lesehilfe überspringt verworfene Zeilen). |
| `internal/adapters/driven/postgresstorage/administrationrequest_test.go` | update, Nachzug | zusätzlich zur Zeile oben: `TestAdministrationRequestListPendingPassesRejectedRowsThrough` (reale Zeilen über die SQL-Funktionen: leeres Schema, leerer Tabellenname, `exclude_column`/`include_column` mit leerer Spalte, je eine gültige Zeile davor bzw. dahinter, Vermerke `failed`/`applied`, Bereinigung nach Kennungen des Tests). |
| `internal/bootstrap/administration_internal_test.go` | update, Nachzug | Fake-Port: Feld `rows` (Zeilen mit Verwurf) und `marked` (Reihenfolge der Vermerke); `TestProcessAdministrationRequestsFailsRejectedRowsInQueueOrder` (verworfene Zeilen vor, zwischen und hinter gültigen; Zeile ohne Kennung), `TestProcessAdministrationRequestsRejectedRowSurvivesMarkFailedError`. |
| `internal/bootstrap/administration_roles_internal_test.go` | update, Nachzug | `TestAdministrationPathRunsUnderLeastPrivilegeLogins`: vier verworfene Zeilen und eine gültige dahinter unter dem `cdc_admin`-Login bis `failed` bzw. `applied`. |
| `internal/domain/model/administrationrequest_test.go` | update | Test-Godoc zu den Regelfeldern trug die Umschreibung „lehnte das Lesen der Queue die Zeile ab“, die das Muster in Zeile 4 des Feldes nicht trifft; gefunden im ergänzenden Lauf des Feldes (Zeile 9, Muster auf Antrag und Lesen); nachgezogen. |
| `docs/user/benutzerhandbuch.md` | geprüft, nicht berührt | keine Aussage zu einer Ablehnung oder einem Anhalten der Queue durch eine fehlerhafte Zeile (Suchlauf Zeile 7 ohne Treffer; die Zeile des Handbuchs mit „Konstruktor“ betrifft ein Konstruktor-Argument des SDK); keine neue Betreiber-Oberfläche (Kandidatenlauf `git diff --name-only 47a646b5 -- internal/bootstrap/ tools/schema/ internal/adapters/driving/` trifft nur Tests und `wiring.go`-Kommentar/Ablauf). Aufschub der Aussage „eine verworfene Zeile endet `failed` mit Text“ im Handbuch: Adresse `slice-transformationen-betriebsdoku`. |

**§3.13-Suchlauf (committetes Feld).** Bewegte Eigenschaft: „was geschieht mit
einer Zeile, die der Antrags-Konstruktor verwirft“ — Symbolnamen
(`ReadPendingRequests`, `ListPending`, `NewAdministrationRequest`), Beschreibung
(„Ablehnung beim Lesen“, „nicht beim Lesen ab“, „abgelehnter Lesevorgang“).
Suchraum: der ganze Baum für Go, die Spec, das Handbuch und `harness`;
ausgenommen sind `docs/reviews/**`, die Records unter `done/` und
`.harness/baseline/**`. Stand ist der Parent `7b70b34a`; der Implementer ergänzt
die Zeilen mit Stand `diff`:

```suchlauf
7b70b34a 8 -E 'ReadPendingRequests|ListPending' -- '*.go' ':!*_test.go'
7b70b34a 73 -E 'ReadPendingRequests|ListPending' -- '*.go'
7b70b34a 7 -E 'NewAdministrationRequest' -- '*.go' ':!*_test.go'
7b70b34a 8 -E 'Ablehnung beim Lesen|abgelehnter Lesevorgang|nicht beim Lesen ab|Tabellennamen beim Lesen ab|endet mit der Konstruktor-Invariante|beim Lesen abgelehnt|ListPending. endet mit' -- '*.go' spec docs/user harness
diff 11 -E 'ReadPendingRequests|ListPending' -- '*.go' ':!*_test.go'
diff 91 -E 'ReadPendingRequests|ListPending' -- '*.go'
diff 7 -E 'NewAdministrationRequest' -- '*.go' ':!*_test.go'
diff 1 -E 'Ablehnung beim Lesen|abgelehnter Lesevorgang|nicht beim Lesen ab|Tabellennamen beim Lesen ab|endet mit der Konstruktor-Invariante|beim Lesen abgelehnt|ListPending. endet mit' -- '*.go' spec docs/user harness
7b70b34a 7 -i -E 'Antr[a-z]*.*(beim Lesen|Lesefehler|Lesung.*(ablehn|Fehler))|(beim Lesen|Lesefehler).*Antr' -- '*.go' spec docs/user harness
diff 6 -i -E 'Antr[a-z]*.*(beim Lesen|Lesefehler|Lesung.*(ablehn|Fehler))|(beim Lesen|Lesefehler).*Antr' -- '*.go' spec docs/user harness
7b70b34a 0 -i -E 'Lesefehler|Ablehnung beim Lesen|abgelehnt beim Lesen|Queue anhalt' -- docs/user harness
diff 0 -i -E 'Lesefehler|Ablehnung beim Lesen|abgelehnt beim Lesen|Queue anhalt' -- docs/user harness
diff 0 -n -E 'administrationrequest.go:[0-9]+|translate.go:[0-9]+|wiring.go:[0-9]+' -- docs/plan/planning/open docs/plan/planning/next docs/plan/planning/in-progress
```

Gemessen am Stand `diff` (Arbeitsbaum nach der Lieferung):

- Zeilen 5 bis 7 (Symbole): 11 Zeilen in Nicht-Test-Dateien (Port 2, Adapter 4, `sqlexec` 3, Verarbeitung 2; am Parent 8 — die Zunahme sind neue Kommentarzeilen an den geänderten Stellen), 91 mit Tests (73 am Parent), 7 Fundstellen des Konstruktors; die Aufrufer, Fakes und Definitionen folgen alle der neuen Form (`make test` und `make test-store` übersetzen sie).
- Zeile 8 (Beschreibung der Ablehnung beim Lesen): **1** von 8 am Parent; die eine Fundstelle ist die Zeile der Änderungshistorie von `SPEC-019` (ein Record, unverändert). Die sieben übrigen sind nachgezogen: Konstruktor-Kommentar (zwei Zeilen), Absatz „Transformations-Antragsarten“ der Spec, zwei Test-Godocs am Store-Test, ein Test-Godoc in `translate_test.go`, ein Kommentar in `administration_internal_test.go`.
- Zeilen 9 und 10 (Beschreibung, weiter gefasst): 7 am Parent, 6 am Diff; die sechs verbleibenden Fundstellen tragen zutreffende Aussagen (Fehlerpfad der Lesung selbst, Testmeldungen, der Record der Änderungshistorie). Zeilen 11 und 12: `docs/user` und `harness` tragen weder am Parent noch am Diff eine Aussage dazu. Zeile 13: keine Zeile der Slice-Pläne zitiert eine Zeilennummer der geänderten Dateien.

**Nicht gefunden, weil nicht durchsucht** (Grenze des Musters, §3.13): eine Formulierung ohne die Wörter „Lesen“, „Lesung“, „Lesefehler“ oder „Queue anhalt…“; Lokatoren in einer anderen Form als `datei.go:zahl`.

| Träger | Befund | Behandlung |
|---|---|---|
| Aufrufer und Definitionen von `ListPending`/`ReadPendingRequests` | 8 Zeilen in Nicht-Test-Dateien, 73 mit Tests (Zeilen 1 und 2 des Feldes) | Port, Adapter, `sqlexec`, Verarbeitung und die Fakes folgen der neuen Form. |
| Beschreibungen der Ablehnung beim Lesen | 8 Zeilen (Zeile 4 des Feldes): Konstruktor-Kommentar (2), `SPEC-019` Absatz „Transformations-Antragsarten“ (1), die Änderungshistorie der Spec (1, ein Record, unverändert), zwei Test-Godocs an `ListPending`, ein Test-Godoc in `translate_test.go`, ein Kommentar in `administration_internal_test.go` | die Zeilen außerhalb des Records werden nachgezogen (Liefer-Punkt 3); die Zeile der Änderungshistorie bleibt. |
| `docs/user/benutzerhandbuch.md` | keine Fundstelle zur Ablehnung beim Lesen (Zeile 4 des Feldes, Suchraum `docs/user`) | am Start prüfen; trägt das Handbuch eine Aussage zu „Anträge, die die Queue anhalten“, ist es ein Träger für `slice-transformationen-betriebsdoku` (Meldung mit Adresse). |
| Beobachtungs-Register `BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue` | Zustand „entschieden, Fix offen (Adresse dieser Slice)“ | Träger der Planner-Closure: `state.md` auf „Fix geliefert“; fremde Datei, deshalb Meldung. |

## 4. Trigger

**Start** (`next` → `in-progress`): kein anderer Slice liegt in `in-progress/`
(WIP-Limit 1). **Reihenfolge (Empfehlung an den Orchestrator):** 1.
`slice-code-kommentare-kennungen`, 2. `slice-harness-fmt-check`, 3. dieser Slice,
danach die Welle [welle-transformationen](../welle-transformationen.md). Dieser
Slice muss `done` sein, **bevor** `slice-transformationen-start-reihenfolge`
startet (Kante, dort §4): beide Slices ändern die Verarbeitung der Queue in
`internal/bootstrap/wiring.go`, und der Vorlauf vor `stream.Run` liest dieselbe
Queue — eine Zeile, die die Lesung anhält, hielte auch den Vorlauf an
(hergeleitet aus `processAdministrationRequests`, nicht erprobt). Zu den
übrigen Slices der Welle gibt es keine Kante; `slice-harness-fmt-check` hat
keine der hier berührten Dateien geändert (die sechs von ihm formatierten
Go-Dateien stehen nicht in §3; `make fmt-check` endet am Baum mit Exit 0).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht erwartet — S bis
  M. Sprengen die Fakes der Fundstellen den Umfang, ist der abtrennbare Teil der
  Umbau des Ports und der Fakes ohne Verhaltensänderung (Zwischenzustand:
  `ListPending` liefert für eine verworfene Zeile weiter den Fehler), danach das
  Verhalten.
- `in-progress` → `open` (blockiert): (a) die Durchreichung braucht eine
  Entscheidung, die der Architect-Zug nicht trägt — eine Zeile ohne Kennung
  (überspringen oder Fehler), eine Antragsart außerhalb der geschlossenen Menge;
  (b) der Fix verlangt eine Änderung an einer SQL-Funktion, am Schema oder an der
  `LISTEN`-Schleife — dann berühren
  [`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md) oder
  [`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)
  den Zug, und die Aussage des Architect-Zugs „Code, keine ADR“ trüge nicht.
  Beide Fälle: Architect-Frage.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + ein realer, grüner `make test-store`-Lauf
(die gedruckte Zeile im Bericht) + Review-Report ohne offenes HIGH oder MEDIUM +
Verifikation, dass die DoD trägt (Spec-Text je Grund gegen den Code) +
Closure-Notiz mit Lerneintrag geschrieben (Rest der Klasse
`BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue` geschlossen, benannte
Spec-Lücke in `SPEC-019` geschlossen).

## 6. Risiken und offene Punkte

- **Der Umbau des Ports trifft viele Fundstellen** (Fakes und Aufrufer, 73
  Zeilen mit dem Symbol, gemessen am Stand `7b70b34a`). *Erwartet, zu belegen
  durch:* das Suchlauf-Feld (Zeilen 1 und 2) an beiden Ständen; die Rückführung
  §4 ist der Ausgang bei Überlauf. **Ausgang:** *(bei Closure)*
- **Ein nicht vermerkbarer Fall hält die Queue an oder wird endlos wiederholt** —
  eine Zeile ohne Kennung, ein Fehler von `MarkFailed`. *Erwartet, zu belegen
  durch:* der Whitebox-Test (Fake): ein Fehler von `MarkFailed` wird protokolliert
  und bricht den Durchlauf nicht ab, eine Zeile ohne Kennung wird übersprungen;
  das Verhalten der Zeile hinter dem Fall ist im Test gebunden. **Ausgang:**
  *(bei Closure)*
- **Der Store-Test hinterlässt `pending`-Zeilen, die andere Tests der Tabelle
  stören** (`BEO-PGC/test-isolation-geteilter-zustand`, 2×). *Erwartet, zu
  belegen durch:* die Bereinigung nach den Kennungen des Tests und ein Lauf des
  ganzen Pakets `postgresstorage` nach dem Slice (grün). **Ausgang:** *(bei
  Closure)*
- **Der Login-Test deckt die verworfene Zeile nicht**
  (`BEO-PGC/rollen-test-abdeckungsluecken`, gestrichen, 4×). *Erwartet, zu belegen durch:*
  `TestAdministrationPathRunsUnderLeastPrivilegeLogins` führt eine verworfene Zeile
  unter `cdc_admin` bis `failed` (Recht `UPDATE` auf der Antrags-Tabelle,
  hergeleitet aus den Grants, im Test zu erproben). **Ausgang:** *(bei Closure)*
- **Spec und Kommentare widersprechen sich nach dem Fix**
  (`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`, 13×). *Erwartet, zu belegen
  durch:* das Suchlauf-Feld (Zeile 4) an beiden Ständen und der Reviewer, der den
  Kontext um jede geänderte Zeile liest (`git diff -U20`). **Ausgang:** *(bei
  Closure)*
- **Der Fix berührt
  [`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md) oder
  [`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md).**
  *Erwartet, nicht belegt:* nein — der Diff ändert weder eine SQL-Funktion noch
  das Schema noch die `LISTEN`-Schleife. *Zu belegen durch:* `git diff --stat`
  ohne `tools/schema/**` und der Reviewer, der beide ADRs gegen den Diff liest.
  **Ausgang:** *(bei Closure)*
- **Die Ordnung der Queue kippt**: die verworfene Zeile wird am Ende statt an
  ihrer Stelle vermerkt. *Erwartet, zu belegen durch:* ein Test mit einer
  verworfenen Zeile **vor** und einer **hinter** gültigen Zeilen, der die
  Verarbeitungsreihenfolge prüft. **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag:** *(zu tragen bei Closure — geschärfte Regel · neuer
  Sensor · benannte Spec-Lücke; erwartet: die benannte Spec-Lücke von
  [`SPEC-019`](../../../../spec/pflichtenheft.md) (Fehlertext und Ort der Prüfung
  für leeres Schema, leeren Tabellennamen und leere Spalte) geschlossen, Auslöser
  `BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue`; ohne Eintrag kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(zu tragen bei Closure —
  `state.md` von `BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue`: Fix
  geliefert)*
- **Folge-Slices:** keine erwartet.
- **Risiken aus §6:** *(je ein Ausgang, zu tragen bei Closure)*
- **Drei Paarungen:** dieser Slice hat keine Welle; die Slice-Closure selbst trägt
  die drei Paarungen (Anker · Folge-Slice · Register), nach dem `git mv` nach
  `done/`.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Port, Store-Adapter, Composition Root und Spec-Absatz
sind keine eigenen Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen (Zähler
gemessen am 2026-09-26 mit `ls evidence | wc -l` je Eintrag):

- `BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue` (entschieden, Fix
  offen, 2×): der Gegenstand dieses Slice.
- `BEO-PGC/test-isolation-geteilter-zustand` (offen, 2×): Risiko in §6 (die
  Bereinigung des Store-Tests).
- `BEO-PGC/rollen-test-abdeckungsluecken` (gestrichen, 4×): der Login-Test bleibt
  die Probe der Rechte auf der Antrags-Tabelle (Risiko in §6).
- `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (verkörpert, 13×),
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 32×): Spec und
  Kommentare, die das alte Verhalten beschreiben — Suchlauf-Feld in §3.
- `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 16×): die
  Mutationen der Eingabeseite in Liefer-Punkt 1.
- `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (verkörpert, 14×),
  `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert, 23×): Befehle
  in Codeblöcken, Zahlen mit Ursprung und Stand.
- `BEO-PGC/kein-admin-weg-schema-fehler-recovery` (offen, 1×): abgegrenzt (§1).
- `BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×): der Runner von `make
  test-store` trägt keinen `-run`-Filter (`grep -n -e '-run'
  tools/harness/run-store-tests.sh` ohne Treffer, gemessen am Stand
  `7b70b34a`); ein neuer Test im Paket ist damit erfasst.
- Gesichtet, ohne Bezug zu diesem Slice: die übrigen Einträge des Registers.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
