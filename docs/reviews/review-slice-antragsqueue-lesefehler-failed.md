# Review-Report: slice-antragsqueue-lesefehler-failed — 2026-09-26

**Review-Art:** Code — der Diff ändert den Port `AdministrationRequestPort.ListPending`
(Rückgabe je Zeile Antrag oder verworfene Zeile), die Lesung der Antrags-Queue (`sqlexec`), die
Verarbeitung `processAdministrationRequests` und den Absatz „Zeilen, die kein Antrag sind“ von
`SPEC-019`; geprüft gegen Plan, die Architect-Entscheidung „Option (a), allgemeine Form“ im
Register-Eintrag `BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue`,
[ADR-0050](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md),
[ADR-0046](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md),
[ADR-0029](../plan/adr/0029-domain-invarianten.md),
[ADR-0034](../plan/adr/0034-ports-nach-faehigkeiten.md),
[ADR-0028](../plan/adr/0028-inbound-use-cases.md),
[ADR-0112](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md),
[ADR-0127](../plan/adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md), `SPEC-019` in
[spec/pflichtenheft.md](../../spec/pflichtenheft.md) und die Hard Rules (Modul 10 §Drei
Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Slice `slice-antragsqueue-lesefehler-failed` (ohne Welle; geht
`slice-transformationen-start-reihenfolge` voraus), Diff-Range `fc107f42..5d8e0748` (5 Commits, 15 Dateien,
+617/−128; Baum sauber, nicht gepusht; der Code- und Spec-Commit ist `fd1a7ef2`).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ in der Form des Diff-Stands
`5d8e0748` (mit MEDIUM „Herkunft als mehrere Felder“, LOW „`make fmt-check`“, HIGH „Beleg trägt seinen
Satz nicht“ und „Zusage ohne Bindung an ihre Eingabeseite“).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-26.

> **Zitier-Form** *(dieser Block bleibt stehen — er ist Norm, kein
> Ausfüll-Hinweis; die `<Platzhalter>` darin sind Formbeispiele)*. Dieser
> Report friert ein; was er zitiert, bewegt sich
> weiter. Deshalb: **Kennung, nicht Adresse** — `slice-<Kennung>` statt seines
> Lifecycle-Pfads, `make <target>` statt eines Links auf die Sensor-Datei, eine
> Baseline-Stelle als **Tag + Pfad in Inline-Code** statt als Link
> (`v<X.Y.Z>` · `regelwerk/<datei>.md` §<Abschnitt>). Der vendored Baum trägt
> genau einen Tag; der Sprung löscht den alten, und ein Link darauf färbt beim
> nächsten Bump ein Artefakt rot, das niemand mehr anfassen darf. Ein `pfad`-Feld
> auf den **geprüften Gegenstand** ist davon nicht betroffen — es zitiert den
> Stand des Laufs und darf ihn festhalten (`v<X.Y.Z>` ·
> `regelwerk/grundlagen-harness-dateien.md` §harness/README.md als
> Einstiegspunkt — diese Zeile ist selbst ein Beispiel der Form).

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde — ohne
diese Liste ist der Lauf nicht reproduzierbar):

- Slice-Plan `slice-antragsqueue-lesefehler-failed` (§1 Ziel und Abgrenzung, §2 Liefer-Punkte als Bezug, §3 Plan-Tabelle,
  „Form der Durchreichung“, Suchlauf-Feld und Träger-Tabelle, §4 Rückführungen, §6 Risiken, §8 Beobachtungen);
  Register-Eintrag `BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue` (`state.md`, Option „allgemeine Form“)
- ADRs: [ADR-0050](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md) (Antrags-Queue),
  [ADR-0046](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md) (keine Domänenlogik in SQL),
  [ADR-0029](../plan/adr/0029-domain-invarianten.md) (Konstruktor-Invarianten),
  [ADR-0034](../plan/adr/0034-ports-nach-faehigkeiten.md) und
  [ADR-0028](../plan/adr/0028-inbound-use-cases.md) (Port-Schnitt),
  [ADR-0112](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) (Prüfung in der Verarbeitung),
  [ADR-0127](../plan/adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md) (Ordnung der Queue),
  [ADR-0083](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md) (Herkunft von Aussagen)
- Anforderungen: [LH-FA-ADM-001](../../spec/lastenheft.md), [LH-FA-CFG-005](../../spec/lastenheft.md),
  [LH-FA-CFG-007](../../spec/lastenheft.md); `SPEC-019` im Pflichtenheft (Spalten der Antrags-Tabelle, Absatz
  „Zeilen, die kein Antrag sind“, Änderungshistorie)
- `AGENTS.md` (Hard Rules §3.1, §3.2, §3.3, §3.4, §3.7, §3.9, §3.12, §3.13), `harness/conventions.md`
  (`MR-000`/`MR-001`)
- Baseline `v6.9.0` · `regelwerk/modul-10-review-harness.md` §Ziel-Form: Reviewer-Skill
- Report-Gerüst: `docs/reviews/review-report.template.md`, Formvorbild `docs/reviews/review-slice-harness-fmt-check.md`

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem Implementer-Bericht übernommen;
Exit-Codes ungepiped in Log-Dateien gesichert, gedruckte Zeilen zitiert; Scratchpad der Sitzung):

- **Läufe am Stand `5d8e0748`** (einer zugleich; `free -m` vorher 14,0 GB verfügbar; dangling Volumes
  `docker volume ls -qf dangling=true | wc -l` 34 vor und 34 nach allen Läufen; kein `prune`;
  `git status --short` danach leer):
  `make test` (Race) Exit 0, 45 Pakete `ok`, kein `FAIL`;
  `make test-store` Exit 0, gedruckt „ok … internal/bootstrap 5.352s“, „ok … postgresstorage 12.255s“ und
  „db-coverage: OK — DB-Adapter-Coverage 82.56% erfuellt Schwelle 80%“;
  `make a-check` Exit 0, „gesamt: 0 Befund(e)“;
  `make coverage-gate` Exit 0, „coverage-gate: OK — Coverage 85.20% erfüllt Schwelle 80%“;
  `make fmt-check` Exit 0, „fmt-check: 255 Go-Dateien geprüft, alle formatiert“;
  `make gates` (einmal) Exit 0, in der Ausgabe „baseline-verify: v6.9.0 OK — 54 Dateien“, „d-check: 1270 Datei(en) geprüft, 0 Befund(e)“,
  „commit-traceability: OK — 5 Commit(s)“, „coverage-gate: OK — Coverage 85.20%“, „generated-sync: OK“,
  „gesamt: 0 Befund(e)“;
  `make commit-traceability RANGE=fc107f42..5d8e0748` Exit 0, „5 Commit(s) … Betreffs ohne Struktur-ID“;
  `make kommentar-kennungen DIFF=fc107f42 COUNT=1` Exit 0, Ausgabe **0**;
  `make suchlauf-nachmessen PLAN=<Plan-Datei>` Exit 0, „suchlauf-nachmessen: 13 Zeilen stimmen“.
- **Eingabeseiten-Mutationen** (Kopie des Stands `5d8e0748` im Scratchpad per `git archive`, Mutation durch
  Ersetzen eines Textstücks in der Kopie, kein `sed -i`; der Test lief im gepinnten Toolchain-Image gegen die Kopie,
  die Store- und Login-Tests gegen eine Wegwerf-PostgreSQL mit dem d-migrate-Rollout; die Kopie wurde nach jedem Lauf
  aus `git` zurückgestellt, `git status` im Repo blieb leer). Sqlexec = `internal/adapters/driven/postgresstorage/sqlexec`,
  Verarbeitung = `internal/bootstrap` ohne Datenbank, Store = `TestAdministrationRequestListPendingPassesRejectedRowsThrough`
  gegen reale PostgreSQL, Login = `TestAdministrationPathRunsUnderLeastPrivilegeLogins` gegen reale PostgreSQL:

  | Nr. | Mutation | Ergebnis (rote Tests) |
  |---|---|---|
  | M1 | `ReadPendingRequests` gibt den Konstruktor-Fehler zurück (`return nil, err`) | rot: `TestReadPendingRequestsPassesRejectedRowsThrough`, `TestReadPendingRequestsKeepsReadFailureAfterRejectedRow`; Store rot (isoliert gesehen); Login rot |
  | M2 | fester Text `Antrag ist ungültig: <id>` statt des Klartexts je Grund | rot: `TestReadPendingRequestsPassesRejectedRowsThrough`; Store rot |
  | M3 | Fälle `schema` und `table` in `rejectionMessage` vertauscht | rot: `TestReadPendingRequestsPassesRejectedRowsThrough` (Fall „leeres Schema und leerer Tabellenname“) |
  | M4 | Anfangswert von `clear` leer | rot: `TestRejectionMessageFallsBackToGeneralText` |
  | M5 | `rows.Err()`-Zweig gibt `requests, nil` zurück | rot: `TestReadPendingRequestsClassifiesIterationFailure`, `TestReadPendingRequestsKeepsReadFailureAfterRejectedRow` |
  | M6 | verworfene Zeilen nach den gültigen (Ordnung im Adapter) | rot: `TestReadPendingRequestsPassesRejectedRowsThrough`; Store rot |
  | M7 | `continue` → `return` nach der verworfenen Zeile in `processAdministrationRequests` | rot: `TestProcessAdministrationRequestsFailsRejectedRowsInQueueOrder`, `TestProcessAdministrationRequestsRejectedRowSurvivesMarkFailedError`; Login rot |
  | M8 | fester Text im Vermerk (`MarkFailed(…, "verworfen")`) | rot: `…FailsRejectedRowsInQueueOrder`; Login rot |
  | M9 | Guard `rejected.ID == ""` entfernt | rot: `…FailsRejectedRowsInQueueOrder` |
  | M10 | Vermerk der verworfenen Zeile erst am Ende (`defer`, Ordnung des Vermerks) | rot: `…FailsRejectedRowsInQueueOrder` |
  | M11 | Warnung der Zeile ohne Kennung entfernt | rot: `…FailsRejectedRowsInQueueOrder` |
  | M12 | Fehler von `MarkFailed` einer verworfenen Zeile bricht ab (`panic`) | rot: `…RejectedRowSurvivesMarkFailedError` |
  | M13 | Warnung „Fehlschlag nicht vermerkt“ entfernt | rot: `…RejectedRowSurvivesMarkFailedError` |
  | M14 | Zweig `row.Rejected != nil` abgeschaltet | rot: `…FailsRejectedRowsInQueueOrder` |
  | S1 | Fälle `source` und `schema` in `rejectionMessage` vertauscht | **grün** (F-4) |
  | S2 | Fall `Antragsart` vor `schema` | **grün** (F-4) |
  | S3/S7 | Fall `Kennung` nach `source` bzw. nach `schema` | **grün** (F-4) |
  | S5 | `column == "" &&` aus dem Fall `Spaltenname` entfernt | grün — äquivalent (F-5): jeder andere Grund mit `ErrEmptyIdentifier` wird vorher gefangen |
  | S6 | Fall `Spaltenname` vor `schema` | rot: der Fall „leeres Schema und leerer Tabellenname“ (kein Befund) |

  Die Angabe des Implementers zu „11 Mutationen“ ist eine Untermenge; die Zusage „ein Fehler von `MarkFailed` bricht
  den Durchlauf nicht ab“ ist nicht nur über eine Mutation gebunden: M7, M12 und M13 färben denselben Test rot,
  eine eigene Mutation ist nicht nötig. Der Store-Test ist bei M1 isoliert rot gesehen (ein Lauf des einen Tests
  gegen die Wegwerf-PostgreSQL; im Runner von `make test-store` bräche der vorgezogene Schritt `internal/bootstrap`
  zuerst ab).
- **Port und Adapter:** `git grep -n 'ListPending(' -- internal` nennt Definition (Port, Adapter, Fake) und alle
  Aufrufer (Verarbeitung, drei Test-Pakete); jede Stelle folgt der Form `[]outbound.PendingAdministrationRequest`,
  die Übersetzung schlägt bei `make test`/`make test-store` nicht fehl. `sqlexec.ReadPendingRequests` reicht einen
  Konstruktor-Fehler als `Rejected` an der Stelle der Ordnung durch und gibt bei Abfrage, Scan und `rows.Err()` die Klasse
  des Aufrufers (`ErrAdministrationStorage`) zurück. Der Text ist wahr gegen den Konstruktor
  `model.NewAdministrationRequest` (unverändert): die Prüfung `id, source, schema, table` leer steht vor der
  Antragsart, die Antragsart vor der Spalte der beiden Spalten-Antragsarten — die Reihenfolge der Fälle in
  `rejectionMessage` folgt ihr. Erreichbarkeit über die SQL-Funktionen: kein `RAISE` in
  `tools/schema/nacharbeit-administration.sql` (die Funktionen prüfen Schema, Tabelle, Spalte nicht); `source_id`
  trägt einen Fremdschlüssel auf `cdc.source`, `request_kind` einen CHECK (`chk_administration_request_kind`) — ein
  leeres `source_id` und eine unbekannte Antragsart entstehen über den Weg der Funktionen nicht.
- **Fehlertext-Injektion:** die Kennung stammt aus der Spalte `administration_request_id` (`text`, Primärschlüssel); der
  Vermerk läuft parametrisiert (`UPDATE … SET error_message = $2`), der Log-Adapter ist ein JSON-Handler
  (`internal/adapters/driven/telemetry/slog.go`), der Werte escaped. Keine Injektion; die Länge des Textes ist wie bei den
  Fehlertexten der Verarbeitung (Regelname, Schema, Tabelle „wie beantragt“) unbegrenzt (F-7).
- **Verarbeitung:** `wiring.go:1338-1391` gelesen: verworfene Zeile → `failRejectedAdministrationRequest`
  → `MarkFailed(id, Message)` an ihrer Stelle (M10 rot), Fehler des Vermerks protokolliert, `continue` (M7/M12/M13 rot);
  Zeile ohne Kennung: Warnung, kein Vermerk, `pending` (M9/M11 rot). Der Statusfilter der Abfrage
  (`WHERE status = 'pending'`, `queries.go:337-341`) und `UPDATE … WHERE status = 'pending'` tragen die Idempotenz: eine
  `failed` gemarkte Zeile wird nicht erneut gelesen. Takt: `administrationPollInterval` = `heartbeatInterval` = 5 s. Eine
  Zeile ohne Kennung ist wegen des Primärschlüssels höchstens einmal vorhanden (eine Warnung je Durchlauf); eine Zeile, deren
  Vermerk dauerhaft scheitert, erzeugt je Durchlauf ihre Log-Zeilen — dieselbe Klasse wie der bestehende `failed`-Pfad
  (F-7). Ordnung nach [ADR-0127](../plan/adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md): die Abfrage ordnet
  `requested_at, administration_request_id`, die Durchreichung erhält die Reihenfolge (M6 rot, Store-Test prüft Zeile
  davor/hinter je Fall).
- **[SPEC-019](../../spec/pflichtenheft.md) gegen den Code:** die sechs Klartexte der neuen Tabelle stehen je genau einmal in
  `spec/pflichtenheft.md` und in `sqlexec/translate.go` (`git grep -c -F` je Text, Zeichen für Zeichen gleich); die Form
  „Klartext, Doppelpunkt, Leerzeichen, Adresse“ entspricht der bestehenden Tabelle (K1–K4); die Zeilen `schema_name`/`table_name`/
  `column_name` und der Absatz „Transformations-Antragsarten“ sind mitgezogen; Spec und Code stehen im selben Commit
  (`fd1a7ef2`); der Spec-Diff trägt keinen `ADR-`, `slice-`, `welle-` oder `ARC-` Bezug (`git diff … -- spec | grep`,
  0 Treffer); die Änderungshistorie trägt eine Zeile. Die Grenze „Zeile ohne Kennung“ ist getragen (M9/M11).
- **Suchlauf-Feld:** `make suchlauf-nachmessen` Exit 0 (13 Zeilen). Von Hand mit `git grep` an **beiden** Ständen
  nachgefahren (`7b70b34a` und dem tatsächlichen Parent `47a646b5`, Ist gleich Soll): Zeile 1 (Symbole, ohne Tests)
  8/8; Zeile 2 (mit Tests) 73/73; Zeile 3 (`NewAdministrationRequest`, ohne Tests) 7/7; Zeile 4 (Beschreibung der
  Ablehnung beim Lesen) 8/8; Zeile 9 (Beschreibung weit) 7/7; am Stand `diff`: Zeile 5 = 11, Zeile 6 = 91, Zeile 7 = 7,
  Zeile 8 = 1, Zeile 10 = 6, Zeilen 11–13 = 0 (das Werkzeug meldet sie `OK`). Der Parent der Plan-Zeilen heißt
  `7b70b34a`, der Parent des Fix-Commits `47a646b5` (F-9). Träger selbst gesucht: `git grep -n -E
  '(administrationrequest|translate|wiring|administration_internal_test|…)\.go:[0-9]+'` außerhalb `docs/reviews` und
  `done/`: nur [ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) (`wiring.go:981`, `:414`) und [ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) (`wiring.go:282-286`); der Diff ändert `wiring.go` erst ab
  Zeile 1336, die Lokatoren sind unberührt. `docs/user/benutzerhandbuch.md`: die Zeilen 218/249/282 („`failed` (Fehlertext in
  `error_message`)“) und 287 bleiben wahr, keine Aussage zu einer Ablehnung beim Lesen; „Konstruktor“ (Zeile 1420) ist das
  SDK-Konstruktor-Argument. `harness/README.md` `make test-store`: keine Aussage, die der Diff ändert. Der Träger der
  Aufschub-Adresse ist **nicht** vollständig (F-3).
- **Handbuch-Kandidatenlauf:** `git diff --name-only 47a646b5 -- internal/bootstrap/ tools/schema/
  internal/adapters/driving/` nennt vier Dateien in `internal/bootstrap` (drei Tests und `wiring.go`); keine
  `CDC_*`-Variable, keine SQL-Funktion, keine Adresse, kein Endpunkt — keine neue Betreiber-Oberfläche; das Handbuch bleibt
  unberührt, die Versionshistorie damit ohne Zug.
- **Hard Rules:** `git diff fc107f42..5d8e0748 | grep -nE 'sed -i|perl -pi|awk -i'` ohne Treffer;
  `git diff --stat -- tools/schema` leer (keine SQL-Funktion, kein Schema, kein Grant geändert, kein
  `make schema-rollout` nötig); keine Suppression (`//nolint`), `.a-check.yml` unberührt, `make a-check` grün.
  Commit-Struktur: `7965076b` (open → next) und `47a646b5` (next → in-progress) sind reine Renames
  (`git show -M --stat`: 0 Zeilen), der Inhalt steht in eigenen Commits (`f7eace20`, `fd1a7ef2`, `5d8e0748`).

---

## Findings

| ID | Kategorie | Befund | Quelle | Pfad | Verifizierbar | Klasse |
|---|---|---|---|---|---|---|
| F-1 | HIGH | Der Doc-Kommentar von `ReadPendingRequests` begründet die Durchreichung mit der verworfenen Alternative im Konjunktiv: „statt die Lesung zu beenden — ein Fehler der Lesung hielte jeden Antrag dahinter an“. Derselbe Halbsatz stand am Parent im Konstruktor-Kommentar (`internal/domain/model/administrationrequest.go:86`, „ein abgelehnter Lesevorgang hielte jeden Antrag dahinter an“), wurde dort entfernt und hier neu geschrieben. Das ist die Form „Ohne X wäre …“ (`AGENTS.md` §3.7: Konjunktiv über die verworfene Alternative statt Indikativ über den Zustand); der Skopus des Skills nennt sie ausdrücklich. | `AGENTS.md` §3.7; Reviewer-Skill HIGH „Kommentar trägt keine der Kommentar-Klassen“ | `internal/adapters/driven/postgresstorage/sqlexec/translate.go:375-376` | ja — `git grep -n 'hielte' -- internal/adapters/driven/postgresstorage/sqlexec/translate.go` (ein Treffer) | Kommentar beschreibt die verworfene Alternative |
| F-2 | HIGH | Die Träger-Tabelle des Plans stützt „keine Aussage zu einer Ablehnung oder einem Anhalten der Queue“ im Handbuch auf „Suchlauf Zeile 7 ohne Treffer“. Zeile 7 des Feldes zählt die Fundstellen von `NewAdministrationRequest` in Nicht-Test-Dateien (`diff 7`, sieben Treffer); die Messung über `docs/user` und `harness` mit Ergebnis 0 stehen in den Zeilen 11 und 12. Die Aussage ist wahr, der genannte Beleg trägt sie nicht (die Zeilennummern verschoben sich mit den am Ende ergänzten Zeilen des Feldes). | Reviewer-Skill HIGH „Beleg trägt seinen Satz nicht“ (Zitat nennt die falsche Stelle); `AGENTS.md` §3.13 | `docs/plan/planning/in-progress/slice-antragsqueue-lesefehler-failed.md:214` (Zeilen 11 und 12 des `suchlauf`-Blocks: `:236-237`) | ja — `make suchlauf-nachmessen` zeigt Zeile 7 als `NewAdministrationRequest`; Zeilen des Blocks nachzählen | Beleg-Verweis auf die falsche Zeile des Suchlauf-Felds |
| F-3 | MEDIUM | Der Plan schiebt „eine verworfene Zeile endet `failed` mit Text“ im Handbuch mit der Adresse `slice-transformationen-betriebsdoku` auf. Der Plan der Adresse trägt die Aussage nicht: `git grep -c -i -E 'verworfen\|Konstruktor\|Lesung\|kein Antrag\|leeres Schema\|leere Spalte' -- <Plan der Adresse>` findet keinen Treffer; sein §2 nennt Status und die Fehlertexte von K1–K4, nicht die Zeilen, die der Antrags-Konstruktor verwirft. Die Meldung an die Adresse steht in einem Träger des Absenders, in keinem committeten Text der Adresse. Der Register-Eintrag `BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` trägt vier Belege (`evidence/`, Zustand „Schwelle erreicht“); dies ist die fünfte Ausprägung. | Reviewer-Skill (Zusatz „Neue Betreiber-Oberfläche ohne Handbuch-Zug“, Probe `git grep` der Kernbegriffe im Plan der Adresse); `AGENTS.md` §3.12 Instanz B (Beleg-Anker); `AGENTS.md` §3.13 (Meldung mit Frist) | `docs/plan/planning/in-progress/slice-antragsqueue-lesefehler-failed.md:214`; `docs/plan/planning/open/slice-transformationen-betriebsdoku.md` §2 (Zeilen 85-125) | ja — der genannte `git grep` gegen den Plan der Adresse (Exit 1, keine Zeile) | Aufschub-Adresse deckt den Gegenstand nicht |
| F-4 | LOW | Die Zusage von `SPEC-019` „die erste verletzte Prüfung bestimmt den Text, in der Reihenfolge der Tabelle von oben nach unten“ ist an einem Paar gebunden (Schema vor Tabellenname, M3 rot) und an der Spalte (S6 rot). Das Vertauschen von Quelle und Schema (S1), das Vorziehen der Antragsart vor das Schema (S2) und das Nachordnen der Kennung (S3, S7) färbt keinen Test rot; kein Testfall trägt zwei dieser Gründe zugleich. Erreichbarkeit: Quelle (Fremdschlüssel) und Antragsart (CHECK) entstehen über die SQL-Funktionen nicht, der Text der Kennung-Zeile geht nur ins Log; die Zusage ist auf diesen Paaren nur per Fake belegt, die Gründe selbst sind es auch (`Quelle ist leer`, `Antragsart ist unbekannt` nur im Fake-Test). Hier **fehlt** Abdeckung, die vorhandene Zusage ist auf dem einen Paar gebunden — deshalb LOW. | `SPEC-019`; Reviewer-Skill „Zusage ohne Bindung an ihre Eingabeseite“ (Abgrenzung: Teil-Bindung) | `internal/adapters/driven/postgresstorage/sqlexec/translate.go:420-437`; `internal/adapters/driven/postgresstorage/sqlexec/translate_test.go` (Fälle der Tabelle in `TestReadPendingRequestsPassesRejectedRowsThrough`) | ja — S1, S2, S3 an einer Kopie im Scratchpad, `go test` grün | Reihenfolge-Zusage nur an einem Paar gebunden |
| F-5 | INFO | `rejectionMessage` bildet den Klartext der Spec im **driven Adapter** aus Feldinspektion (`source == ""`, `schema == ""`, …) und wiederholt damit die Prüfreihenfolge des Konstruktors; die Kopplung steht nur im Kommentar („in der Reihenfolge der Prüfungen des Konstruktors“). Ändert der Konstruktor Reihenfolge oder Gründe, färbt kein Test sie rot (der Test speist die Felder, nicht den Konstruktor); ein neuer Grund fällt auf `Antrag ist ungültig`, eine umgestellte Reihenfolge auf einen falschen Klartext. Alternativen (Fehlerwert je Zeile, zwei Listen) sind gegen die Form bewertet: zwei Listen verlören die Ordnung (M6), ein Fehlerwert je Zeile trüge den Grund, nicht den Text — die Form `Request`/`Rejected` ist die kleinste tragende. Der Fall `column == "" &&` ist redundant (S5 äquivalent). | Maintainability; [ADR-0029](../plan/adr/0029-domain-invarianten.md), [ADR-0034](../plan/adr/0034-ports-nach-faehigkeiten.md) | `sqlexec/translate.go:412-437` | nein | Kopplung an den Konstruktor nur im Kommentar |
| F-6 | INFO | Die Zustandsinvariante von `PendingAdministrationRequest` („entweder `Request` oder `Rejected`“) steht im Port-Kommentar und wird nur durch die eine Erzeugerstelle (`sqlexec`) eingehalten; der Typ erzwingt sie nicht. Ein Wert `PendingAdministrationRequest{}` (beide leer) läuft als Antrag mit leerer Kennung in `applyAdministrationRequest`, dessen `default`-Zweig ihn ablehnt, und `MarkFailed` endet am Guard mit `ErrEmptyIdentifier` — ohne Test. Für die heutige Erzeugerstelle nicht erreichbar. | Maintainability; [ADR-0034](../plan/adr/0034-ports-nach-faehigkeiten.md) | `internal/application/port/outbound/administrationrequest.go:19-36` | nein | Summentyp per Konvention statt per Typ |
| F-7 | INFO | Benannte Grenzen der Verarbeitung, bewertet: (a) eine Zeile ohne Kennung ist wegen des Primärschlüssels höchstens einmal vorhanden und erzeugt eine Warnung je Durchlauf (5 s, `administrationPollInterval`), keine Flut; (b) scheitert der Vermerk einer verworfenen Zeile dauerhaft, erzeugt jede solche Zeile je Durchlauf ihre Warnungen („Antrag verworfen“, „Fehlschlag nicht vermerkt“, dazu die Zeile des Adapters) — dieselbe Klasse wie der bestehende `failed`-Pfad; (c) `error_message` trägt die Kennung ungekürzt (`text`, keine Längengrenze), wie die bestehenden Fehlertexte. Keine Endlosschleife ohne Fortschritt: der Durchlauf endet je Runde, die Zeile bleibt `pending` bis der Vermerk gelingt. | Maintainability; [ADR-0050](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md) | `internal/bootstrap/wiring.go:1338-1391` | nein | Wiederholung im Fehlerfall ohne Drosselung (Kenntnis) |
| F-8 | INFO | Der Kommentar über `processAdministrationRequests` trägt im geänderten Block weiterhin „ein `pending` bleibender Antrag würde jeden Durchlauf erneut versuchen“ (Konjunktiv über den Zustand ohne Vermerk; Wortlaut am Parent unverändert, `wiring.go:1342`, im Diff nur neu umbrochen). Nicht Gegenstand dieses Diffs; Bestand für die Bereinigung der Kommentare. | `AGENTS.md` §3.7 (Bestand) | `internal/bootstrap/wiring.go:1342` | ja — `git grep -n 'würde jeden Durchlauf' -- internal/bootstrap/wiring.go` | Konjunktiv im Bestand des geänderten Blocks |
| F-9 | INFO | Die Plan-Zeilen des Suchlaufs tragen als Stand des Parents `7b70b34a`; der Parent des Fix-Commits ist `47a646b5` (`fc107f42` der Basis dieses Reviews), dazwischen liegen Commits, die `.go`-Dateien ändern (Kommentar-Werkzeug, Format). Die gemessenen Zahlen sind an beiden Ständen gleich (Zeilen 1, 2, 3, 4, 9 nachgemessen), der Stand ist damit nicht falsch, aber nicht der Parent. | `AGENTS.md` §3.12 (Stand einer Messung) | `docs/plan/planning/in-progress/slice-antragsqueue-lesefehler-failed.md:222` | ja — `git grep` der Zeilen 1–4 an `7b70b34a` und `47a646b5` | Stand-Bezeichnung „Parent“ nicht der Parent |
| F-10 | INFO | Zur Reduktion der Kennungen auf eine je Kommentarblock (`make kommentar-kennungen DIFF=fc107f42 COUNT=1` Ausgabe 0): entfernt sind `ADR-0029` (2×, Kopf und `ReadPendingRequests`), `ADR-0047`, `LH-QA-SEC-001`, `LH-QA-SEC-002`, `ADR-0112`, `ARC-001` und drei `LH-FA-CFG-007`; jede ist im Baum weiter vorhanden (`git grep -l` je Kennung am Parent/heute: `LH-QA-SEC-001` 9→8 Dateien, `ADR-0047` 18→17, `ADR-0112` 17→16, `ARC-001` 2→1). Eine Kennung trug zwei Aussagen im Godoc von `TestAdministrationPathRunsUnderLeastPrivilegeLogins` (Antragspfad `ADR-0050` und Rollen `ADR-0047`/`LH-QA-SEC-001`/`-002`); der Rollenbezug bleibt über Testnamen und die Grants-Verweise im Test erkennbar, eine Aktion ist nicht erwartet. Nebenbeobachtung: die Godoc-Zeile „zuweist: die Antrags-Queue“ bricht kurz um (Reste der Umformulierung), ohne Wirkung. | Reviewer-Skill MEDIUM „Herkunft als mehrere Felder“ | `internal/bootstrap/administration_roles_internal_test.go:53-58` | ja — `make kommentar-kennungen DIFF=fc107f42 COUNT=1` | Reduktion der Kennungen ohne Informationsverlust (Kenntnis) |
| F-11 | INFO | Der Bericht des Implementers nennt einen versehentlichen `python3`-Aufruf mit leerem Heredoc ohne Wirkung; im Diff gibt es keine Spur eines Host-Interpreters (`git diff … \| grep -nE 'sed -i\|perl -pi\|awk -i'` 0 Treffer, die Testdateien folgen `gofmt`, `make fmt-check` Exit 0). Kenntnis, nach `AGENTS.md` §3.1 kein Verstoß am Erzeugnis. | `AGENTS.md` §3.1 | — | ja — `make fmt-check` | Host-Interpreter ohne Wirkung (Kenntnis) |

## Negativbefunde

| Bereich | Ergebnis |
|---|---|
| `internal/application/port/outbound/administrationrequest.go` (Form `PendingAdministrationRequest`/`RejectedAdministrationRequest`, Kommentare) | geprüft, ohne Befund über F-5, F-6 hinaus: kleinste tragende Form, Ordnung bleibt, der Port trägt Domänentypen und Text, keinen Treibertyp ([ADR-0034](../plan/adr/0034-ports-nach-faehigkeiten.md), [ADR-0028](../plan/adr/0028-inbound-use-cases.md)); der Kommentar nennt eine Kennung (`SPEC-019`) und trägt Zusage und Grenze, das Verhalten liegt im Adapter und in der Verarbeitung (Rang-Zeiger) |
| `internal/adapters/driven/postgresstorage/sqlexec/translate.go` (`ReadPendingRequests`, `rejectionMessage`) | geprüft, ohne Befund über F-1, F-4, F-5 hinaus: Konstruktor-Fehler als `Rejected` an der Stelle, Fehler von Abfrage/Scan/Iteration als Klasse des Aufrufers (M1, M5 rot), Reihenfolge der Gründe wahr gegen den Konstruktor, Text Zeichen für Zeichen wie `SPEC-019`, Kennung als Adresse ohne Injektionsfläche |
| `internal/adapters/driven/postgresstorage/administrationrequest.go` | geprüft, ohne Befund: `ListPending` reicht die neue Form mit einem Rang-Zeiger im Kommentar durch, `MarkFailed`/`MarkApplied` unverändert (Guard `id == ""`, Idempotenz durch `WHERE status = 'pending'`) |
| `internal/bootstrap/wiring.go` (`processAdministrationRequests`, `failRejectedAdministrationRequest`) | geprüft, ohne Befund über F-7, F-8 hinaus: Vermerk an der Stelle der Queue (M10 rot), Fehler des Vermerks protokolliert (M7/M12/M13 rot), Zeile ohne Kennung mit Warnung übersprungen (M9/M11 rot), Quellbindung des `backfill`-Antrags unverändert; Zusage-Kommentare (`lässt den Durchlauf weiterlaufen`, `bleibt dann pending`) sind vom Code getragen |
| `internal/domain/model/administrationrequest.go` | geprüft, ohne Befund: Konstruktor unverändert, der Kommentar trägt eine Kennung und beschreibt den Ist-Zustand („endet ebenfalls `failed`“, Rang-Zeiger auf `SPEC-019`); die Konjunktiv-Form des Parents ist entfernt (F-1 nennt ihre neue Stelle) |
| `sqlexec/translate_test.go`, `sqlexec/rejection_internal_test.go`, `administration_internal_test.go`, `administrationrequest_test.go` (Store), `administration_roles_internal_test.go` (Login), `administrationrequest_order_test.go`, `administration_endtoend_test.go`, `internal/domain/model/administrationrequest_test.go` | geprüft, ohne Befund über F-4, F-10 hinaus: 14 von 14 Mutationen der Eingabeseite rot (Tabelle oben), Store-Test mit realen Zeilen (vier verworfene Zeilen je Grund mit gültiger Zeile davor und dahinter, Vermerke `failed`/`applied`, Bereinigung nur nach den Kennungen des Tests), Login-Test unter `cdc_admin` führt die vier verworfenen Zeilen bis `failed` mit Text und die Zeile dahinter bis `applied`, Aufrufer der Fakes und Lesehilfen überspringen `Rejected` mit Bedingung, `TestReadPendingRequestsKeepsReadFailureAfterRejectedRow` bindet den Fehlerpfad der Lesung; die Mutationsangaben in den Godocs (M1, M2, M3, M4, M5, M6, M7, M8, M10) stimmen mit meinen Läufen überein |
| `spec/pflichtenheft.md` `SPEC-019` | geprüft, ohne Befund: Absatz, Tabelle, Grenze und Nachzug der Zeilen `schema_name`/`table_name`/`column_name` und des Absatzes „Transformations-Antragsarten“ sind wahr gegen den Code; keine Kennung außerhalb der Spec-Strate (§3.4), die Änderungshistorie trägt eine Zeile, Spec und Code im selben Commit; „die SQL-Funktionen prüfen Quelle, Schema, Tabelle und Spalte nicht“ trägt (kein `RAISE`); die Gründe „Quelle ist leer“ und „Antragsart ist unbekannt“ stehen als Zeilen der Tabelle, sind über die Funktionen nicht erreichbar (Fremdschlüssel, CHECK) und nur per Fake erprobt (Kenntnis, kein Widerspruch) |
| Plan `slice-antragsqueue-lesefehler-failed` (Nachzug, Form der Durchreichung, Suchlauf-Feld, Träger-Tabelle, DoD-Haken) | geprüft, ohne Befund über F-2, F-3, F-9 hinaus: die Träger-Tabelle nennt Gefundenes und Nichtgefundenes je Träger, beide Stände; die Zählungen der Beschreibungen (acht am Parent, eine am Diff, sechs zutreffende) sind nachgemessen; DoD-Haken stimmen mit dem Diff, die offenen Punkte (Review, Closure-Notiz, Register, Risiko-Ausgänge, Paarungen) bleiben `[ ]` |
| `harness/README.md`, `docs/user/benutzerhandbuch.md`, `tools/schema/**` | geprüft, ohne Befund: keine neue Betreiber-Oberfläche, Handbuch und Versionshistorie unberührt und ohne Zug, `git diff --stat -- tools/schema` leer, Zeilen-Lokatoren auf `wiring.go` ([ADR-0082](../plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md), [ADR-0088](../plan/adr/0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md)) unberührt |
| Hard Rules, Commit-Struktur | geprüft, ohne Befund: Docker-only (Läufe über `make`), kein `sed -i`/`perl -pi`/`awk -i` im Diff, keine Suppression, `make a-check`/`make coverage-gate`/`make gates` grün, fünf Betreffs mit `LH-FA-ADM-001` und `ADR-0050` ohne `SPEC-`/`ARC-` im Betreff, zwei Lifecycle-Moves rein, Inhalt in eigenen Commits |

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 7 |

**Finding-Klassen dieses Laufs:** Kommentar beschreibt die verworfene Alternative · Beleg-Verweis auf die falsche
Zeile des Suchlauf-Felds · Aufschub-Adresse deckt den Gegenstand nicht · Reihenfolge-Zusage nur an einem Paar gebunden ·
Kopplung an den Konstruktor nur im Kommentar · Summentyp per Konvention statt per Typ · Wiederholung im Fehlerfall ohne
Drosselung (Kenntnis) · Konjunktiv im Bestand des geänderten Blocks · Stand-Bezeichnung „Parent“ nicht der Parent ·
Reduktion der Kennungen ohne Informationsverlust (Kenntnis) · Host-Interpreter ohne Wirkung (Kenntnis)

## Verdikt

**Merge-blockierend:** ja — wegen F-1 und F-2 (beide HIGH nach der Liste des Reviewer-Skills) und F-3 (MEDIUM): ein
Doc-Kommentar formuliert die verworfene Alternative im Konjunktiv, ein Beleg-Verweis im Plan nennt die falsche Zeile des
Suchlauf-Felds, und die Aufschub-Adresse trägt den Gegenstand nicht als Text. Der Code selbst ist in Ordnung: Form der
Durchreichung, Ordnung, Fehlerklassen, Text je Grund (wahr gegen Konstruktor und `SPEC-019`) und die Verarbeitung halten
14 von 14 Mutationen der Eingabeseite stand (auch an realer PostgreSQL und unter dem `cdc_admin`-Login), alle Läufe
(`make test`, `make test-store`, `make a-check`, `make coverage-gate`, `make fmt-check`, `make gates`) sind grün. Die
Korrekturen betreffen einen Halbsatz (F-1), eine Zeilennummer (F-2) und einen Übergabe-Text im Plan der Adresse (F-3,
fremde Datei: Meldung an den Planner mit Frist Closure dieses Slice); F-4 geht in derselben Fixrunde mit oder wird mit
Grund angenommen; F-5 bis F-11 sind Hinweise ohne erwartete Aktion.

**Übergabe:** Findings gehen an den Implementer (Fixrunde); die **Finding-Klassen** gehen zusätzlich in die
Slice-Closure §7 und von dort in den Zähler (F-3 ist die fünfte Ausprägung von
`BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an`, F-1 gehört in den Zähler zur Konjunktiv-Klasse der Kommentare). Die
DoD-Zeile „Review durchgeführt“ im Slice-Plan bleibt offen (eine Fixrunde folgt; sie wird bei Schritt 21 des
Implementer-Ablaufs nachgezogen). Dieser Report selbst ist ein **Lauf-Beleg** (Audit: dieser Diff, dieser Skill, dieses
Modell, dieses Verdikt) — er wird über Läufe hinweg nicht wieder gelesen, und muss es nicht. Der Report ersetzt keine
Verifikation — DoD-/Spec-Konformität prüft der Verifier separat (Modul 11; anderes Prüf-Artefakt, anderer
Eingabe-Kontext).
