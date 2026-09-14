# Review-Report: slice-066 — 2026-09-14

**Review-Art:** Code — geprüft gegen Plan + Konventionen (Modul 10
§Drei Review-Arten); DoD-/Spec-Konformität ist Verifier-Aufgabe und nicht
Gegenstand dieses Reports.

**Gegenstand:** `slice-066`
(`docs/plan/planning/in-progress/slice-066-spaltenausschluss-sql-funktionen.md`),
Diff `529f021..2246f63` (Elter-Commit `529f021` ist ein reiner
`next→in-progress`-Move und trägt keinen Inhalt). Vier Commits: `82ce83e`
(Schema/Antrags-Queue), `0d2030f` (Verarbeitung), `131fd98`
(`spec/architecture.md`/`spec/pflichtenheft.md`), `2246f63` (DoD-Häkchen).

**Skill:** `.harness/skills/reviewer.md` @ `2246f63` (Stand zum Review-Zeitpunkt,
unverändert seit der letzten Schärfung 2026-09-13).
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-14.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `docs/plan/planning/in-progress/slice-066-spaltenausschluss-sql-funktionen.md` (vollständig)
- [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) (vollständig),
  insbesondere Teilfrage 3/5, §Konsequenzen (Folgepflichten) und
  §Fitness Function
- [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md)
  (Antrags-Queue, Erweiterungsbasis), [`ADR-0028`](../plan/adr/0028-inbound-use-cases.md),
  [`ADR-0034`](../plan/adr/0034-ports-nach-faehigkeiten.md),
  [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)
  (Ausweichform/Re-Evaluierungs-Trigger), [`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
- `AGENTS.md` §3.3/§3.7/§3.9, §5 (Traceability), §6 (Rollenwechsel)
- `harness/conventions.md` (`MR-000` ID-Schema),
  `docs/plan/planning/observations/BEO-PGC/d-migrate-nacharbeit` (Registerstand,
  `observation.md`/`state.md`), `BEO-PGC/slice-chronik-in-code-kommentar`
  (`state.md`, `evidence/`), `BEO-PGC/dod-checkbox-nachzug`
- vorherige Findings am gleichen Modul: `docs/reviews/review-slice-060.md`
  (F-1 Chronik-Klasse), `docs/reviews/review-slice-062.md`,
  `docs/reviews/review-slice-065.md`
- `spec/lastenheft.md` (`LH-FA-CFG-005`, `LH-QA-SEC-004`, `LH-FA-ADM-001`),
  `spec/architecture.md` (`ARC-002`/`ARC-004`/`ARC-005`), `spec/pflichtenheft.md`
  (`SPEC-019`)

**Eigene Sensor-Läufe (Exit-Code jeweils ungepiped in einem eigenen Schritt
ermittelt, `AGENTS.md` §3.9):**

| Lauf | Exit | Bemerkung |
|---|---|---|
| `make gates` | **0** | baseline-verify, d-check (517 Dateien, 0 Befunde), commit-traceability (5 Commits), a-check (0 Befunde), coverage-gate |
| `make test` | **0** | vollständige Suite im Race-Container, alle Pakete `ok`, inkl. der zwei neuen `usecase`-Pakete |
| `make test-store` | **0** | `internal/bootstrap` vorgezogen (`0.725s`, alle realen PostgreSQL-Tests grün, inkl. des neuen E2E-Tests), danach alle übrigen Pakete |

**Eigene Nachmessung zu Punkt 1 (d-migrate/CHECK) — Aufbau und Ergebnis:**

Auf einer wegwerfbaren `postgres:18-alpine`-Instanz (Digest wie im Makefile),
`cdc`-Schema + `search_path` wie in `tools/harness/run-store-tests.sh`:

1. **Erstanlage (frisch, aktueller Baum):** `make schema-rollout` → Exit 0;
   `nacharbeit-administration.sql:58` meldet `NOTICE: constraint … does not
   exist, skipping`, die Klausel entsteht anschließend mit den vier Arten.
2. **Wiederholte Ausweichform:** dreimaliges `psql -f
   tools/schema/nacharbeit-administration.sql` → dreimal Exit 0, Endzustand
   unverändert (`chk_administration_request_kind` mit
   `enable`/`disable`/`exclude_column`/`include_column`). Idempotent.
3. **Messung des Implementers nachgestellt:** ein Modell mit *deklariertem*
   CHECK (`schema.yaml`-Kopie plus der vierwertigen Klausel) per
   `schema migrate --execute` gegen eine Instanz, deren Tabelle die alte
   zweiwertige Klausel trägt → **Exit 5**,
   `[ERROR] Post-execute compare detected drift`; danach ist
   `chk_administration_request_kind` **verschwunden**, die neue Klausel
   existiert nicht, `column_name` **ist** angelegt. Das ist exakt das im
   Plan-Nachzug beschriebene Bild — die Messung ist reproduziert.
4. **Instanz auf dem Stand *vor* dem Slice** (Eltern-`schema.yaml` +
   Eltern-`nacharbeit-administration.sql` + Rollen-/Observability-/Heartbeat-
   Nacharbeit ausgerollt, Tabelle trägt die zweiwertige Klausel, keine
   `column_name`), dann `make schema-rollout` **mit dem aktuellen Baum** →
   **Exit 8** (`DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION`; `make` bricht
   mit „Fehler 8" ab, die psql-Nacharbeitsschritte laufen nicht mehr).
   Dieselbe Instanz mit dem **Eltern**-`schema.yaml` → ebenfalls **Exit 8**.
5. **Zuordnung der neuen destruktiven Operation:** auf einer aktuell
   ausgerollten Instanz enthält der Plan 11 Operationen, darunter
   `AlterColumnType:COLUMN … ["administration_request","request_kind"]`, dessen
   erste Anweisung `ALTER TABLE "administration_request" DROP CONSTRAINT IF
   EXISTS "chk_administration_request_kind";` lautet. Wird die Klausel
   vorher manuell entfernt, sinkt der Plan auf 10 Operationen und genau diese
   Operation verschwindet.

---

## Findings

### F-1 — Slice-Chronik in Produktionscode-Kommentar (`NewAdministrationRequest`)

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.7 (Hard Rule „Ein Kommentar beschreibt, was da
  ist") · `.harness/skills/reviewer.md` §Klassifikation (HIGH-Bullet
  „Slice-/Wellen-Chronik in Produktionscode-Kommentar") · Präzedenzfall
  `docs/reviews/review-slice-052.md` F-1 und `docs/reviews/review-slice-060.md`
  F-1 · `BEO-PGC/slice-chronik-in-code-kommentar` (5 Belege, „verkörpert")
- `pfad`: `internal/domain/model/administrationrequest.go:51`
- `befund`: Der Godoc-Block über dem **Produktionscode**-Konstruktor
  `NewAdministrationRequest` (nicht ein `Test*`-Godoc) begründet die
  Verlagerung der Mengen-Prüfung an den Domain-Core-Rand mit einem
  Review-Fund samt Slice-Kennung: „… wie es die Architektur-Sicht für
  Domänenobjekte und ihre Invarianten vorsieht (`Review-Finding F-4`,
  `review-slice-037.md`)". Satzsubjekt ist der Produktionscode-Pfad, nicht ein
  Testfall — genau die in `review-slice-052.md` F-1 etablierte Probe. Die
  Zeile ist in diesem Diff eine **`+`-Zeile**: derselbe Block wurde
  umgeschrieben (zwei andere Chronik-Stellen — „bislang die einzige Ausnahme
  davon" und „lag zuvor … trägt sie jetzt" — sind dabei korrekt entfernt
  worden), die Fund-Referenz ist als einzige stehen geblieben. Kein Gate
  fängt das; der repo-weite Textmuster-Sensor wurde geprüft und verworfen
  (Architect-Verdikt `architect-verdict-slice-chronik-in-code-kommentar.md`).
- `verifizierbar`: nein — kein Gate prüft Kommentar-Klassen
- `klasse`: „Slice-Chronik in Produktionscode-Kommentar"

### F-2 — Plan-Nachzug belegt einen „Exit 0"-Ausrolllauf, der auf einer bestehenden Instanz nicht reproduzierbar ist; die destruktive Upgrade-Semantik der Ausweichform bleibt unbeschrieben

- `kategorie`: MEDIUM
- `quelle`: Maintainability · [`ADR-0043`](../plan/adr/0043-schemamigrationen-mit-d-migrate.md)
  §Re-Evaluierungs-Trigger (Ausweichform) · Slice-Plan §3 (Plan-Nachzug)
- `pfad`: `docs/plan/planning/in-progress/slice-066-spaltenausschluss-sql-funktionen.md:175-180`
  (§3, Satz „Beleg: `make schema-rollout` gegen eine Instanz mit dem Stand
  *vor* dieser Änderung läuft Exit 0 …"); betroffene Artefakte:
  `tools/schema/schema.yaml:226-240` (Constraint bewusst nicht deklarativ),
  `tools/schema/nacharbeit-administration.sql:15-29`, `:58-59`
- `befund`: Nachgestellt gegen eine reale PostgreSQL-18-Instanz ergibt der
  zweite Aufzählungspunkt der Plan-Nachzug-Schlussfolgerung sich nicht: eine
  Instanz auf dem Stand *vor* dieser Änderung (Eltern-Modell ausgerollt,
  nacharbeit-Objekte vorhanden) endet bei `make schema-rollout` mit dem
  aktuellen Baum in **Exit 8**
  (`DESTRUCTIVE_OPERATION_REQUIRES_CONFIRMATION`), nicht in Exit 0 — der
  Lauf bricht vor den psql-Nacharbeitsschritten ab, der versprochene
  Endzustand (Spalte + vierwertige Klausel) entsteht so nicht. Der Blocker
  ist **nicht neu** (derselbe Lauf mit dem *Eltern*-`schema.yaml` gegen
  dieselbe Instanz endet ebenfalls Exit 8; die Drop-Funktion-/Drop-View-Objekte
  sind die seit `slice-011`/`slice-012`/`slice-036` undeklarierten
  Nacharbeitsobjekte). **Neu** durch diesen Diff ist die zusätzliche
  destruktive Plan-Operation auf `administration_request.request_kind`
  (`ALTER TABLE … DROP CONSTRAINT IF EXISTS "chk_administration_request_kind"`
  plus No-op-`ALTER COLUMN … TYPE TEXT`, geführt als
  `AlterColumnType`): sie entsteht erst dadurch, dass die Klausel aus dem
  deklarativen Modell genommen wurde, während sie im Katalog liegt; manuelles
  Entfernen der Klausel lässt den Plan von 11 auf 10 Operationen fallen und
  genau diese Operation verschwinden. Damit dokumentiert der Plan-Nachzug
  weder den realen Ausgang der Upgrade-Konstellation noch die neue
  Konsequenz der Aufteilung (das Modell erklärt ein Objekt, das es im
  Katalog vorfindet, zur Entfernung). Die Messung *zum CHECK selbst* (Exit 5,
  bestehende Klausel entfällt, `column_name` konvergiert) und die
  Idempotenz der Ausweichform sind unabhängig nachgestellt und bestätigt —
  siehe Eigene Nachmessung oben.
- `verifizierbar`: ja — `make schema-rollout` gegen eine Instanz auf dem
  Eltern-Stand (Wiederholung der Schritte 4/5 der eigenen Nachmessung)
- `klasse`: „Beleg-Aussage im Plan nicht reproduzierbar (Bestands-Instanz)"

### F-3 — Neuer Ablehnungszweig des Domain-Konstruktors ohne Test

- `kategorie`: MEDIUM
- `quelle`: `.harness/skills/reviewer.md` §Klassifikation (MEDIUM-Bullet
  „fehlende Negativtests bei neuem öffentlichem Vertrag") ·
  `internal/domain/model/validation_test.go` (Repo-Muster: jeder Konstruktor
  bekommt seinen Invarianten-Ablehnungstest)
- `pfad`: `internal/domain/model/administrationrequest.go:55-68`
- `befund`: Der Diff erweitert den öffentlichen Konstruktor
  `model.NewAdministrationRequest` um eine neue Invariante (die beiden
  Spalten-Antragsarten brauchen eine nichtleere Spalte, sonst
  `ErrEmptyIdentifier`) und um zwei neue Werte der geschlossenen Menge
  (der `default`-Zweig trägt jetzt die Ablehnung aller übrigen Werte). Für
  `internal/domain/model` existiert **keine** Testdatei zu diesem Typ —
  `grep -rn "Administration" internal/domain/model/*_test.go` liefert keinen
  Treffer, `model.NewAdministrationRequest` wird produktiv ausschließlich aus
  `ListPending` aufgerufen, und die einzigen Tests, die den Konstruktor
  erreichen, führen ausschließlich gültige Zeilen (die SQL-Funktionen setzen
  `column_name` immer). Der neue Leer-Spalten-Zweig und der
  Mengen-`default`-Zweig sind damit ungetestet, obwohl der Diff im selben
  Zug vier Service-Tests für den darüberliegenden Use Case ergänzt.
- `verifizierbar`: ja — `make test` mit einem ergänzten Testfall
- `klasse`: „Neuer Invarianten-Zweig am öffentlichen Konstruktor ohne Negativtest"

### F-4 — Plan-§3- und DoD-Beleg-Angabe nicht auf die tatsächlich tragenden Dateien nachgezogen

- `kategorie`: LOW
- `quelle`: Maintainability · Modul 5 §Ziel-Form: Slice („Plan ist die
  Stelle, an der Spec und ADR auf einen Code-Diff zusammenfallen")
- `pfad`: `docs/plan/planning/in-progress/slice-066-spaltenausschluss-sql-funktionen.md:113-123`
  (§2, zweiter DoD-Punkt) und `:146-156` (§3-Tabelle)
- `befund`: Der zweite DoD-Punkt benennt
  `internal/adapters/driven/postgresstorage/administrationrequest_test.go`
  als Beleg für „Happy Path `applied`, Negative-Fall nicht existierende
  Spalte `failed` mit Fehlertext". Der dort neue Testfall ruft die beiden
  SQL-Funktionen real auf und liest `column_name`/`request_kind` zurück,
  setzt den `applied`/`failed`-Übergang aber selbst über
  `MarkApplied`/`MarkFailed` mit einer **konstruierten**
  `ErrSourceColumnMissing`-Meldung — die nicht existierende Spalte wird dort
  nicht geprüft. Der real ausgeübte Pfad (`applyAdministrationRequest` →
  Inbound Port → `ColumnExclusionPort` gegen den Katalog) liegt in einer
  Datei, die weder die DoD-Zeile noch die §3-Tabelle nennt:
  `internal/bootstrap/administration_endtoend_test.go` (neu, 176 Zeilen).
  Auch `internal/domain/model/administrationrequest.go`,
  `internal/domain/errors/errors.go` und
  `internal/adapters/driven/postgresstorage/{administrationrequest,queries}.go`
  stehen nicht in §3, obwohl alle fünf produktiv zur Erfüllung der
  gelisteten Punkte nötig sind. Ein Plan-Nachzug existiert nur für die
  Migrationsform (§3, 2026-09-14).
- `verifizierbar`: ja — Abgleich `git diff 529f021 2246f63 --name-only`
  gegen die §3-Tabelle und die Beleg-Angabe
- `klasse`: „Plan-Tabelle/Beleg-Angabe nicht auf den tatsächlichen Umfang nachgezogen"

### F-5 — Kopplungs-Kommentar im `Makefile` beschreibt den Inhalt der Nacharbeit-Datei unvollständig

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Kommentar-Klasse „Kopplung" beschreibt, was da
  ist)
- `pfad`: `Makefile:143-149` (Kommentar über dem
  `nacharbeit-administration.sql`-Schritt)
- `befund`: Der Kommentar führt den Schritt als Träger „der schreibenden
  SQL-Funktionen `cdc.enable_table`/`cdc.disable_table`" und begründet ihn
  ausschließlich mit der d-migrate-Grenze für Funktionen. Die Datei trägt
  seit diesem Diff vier Funktionen und zusätzlich die
  `request_kind`-CHECK-Klausel (`nacharbeit-administration.sql:6-32` nennt
  selbst zwei Objektklassen); der `Makefile`-Kommentar wurde nicht
  mitgezogen und nennt die zweite Objektklasse gar nicht.
- `verifizierbar`: nein — kein Gate prüft Kommentar-Inhalte
- `klasse`: „Kopplungs-Kommentar neben geändertem Artefakt nicht nachgezogen"

### F-6 — Sicht-Aussage über das eigene Sequenzdiagramm trägt nicht für alle vier Antragsarten

- `kategorie`: INFO
- `quelle`: Maintainability · `spec/architecture.md` ist Rang 3 (Sicht)
- `pfad`: `spec/architecture.md:237`
- `befund`: Der neue Satz „Die Antragsarten-Wahl im Diagramm unten steht für
  alle vier" steht über dem unveränderten Diagramm, dessen Schritt
  `AP->>EUC: Aktiviere Tabelle t` und dessen Teilnehmer
  `EUC as EnableTableUseCase` den Enable-Pfad festlegen; für
  `exclude_column`/`include_column` ruft der Hintergrund-Zug laut der neuen
  Tabelle denselben Pfads, aber einen anderen Inbound Port auf. Die
  Folgepflicht aus [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md)
  (die beiden neuen Antragsarten „auf demselben Pfad") ist durch die neue
  Tabelle erfüllt; die Zusatz-Behauptung über das Diagramm ist eine
  Prosa-Genauigkeitsfrage ohne Verhaltenswirkung. Kein Gate fängt das.
- `verifizierbar`: nein — Diagramm-Prosa wird von keinem Sensor geprüft
- `klasse`: „Sicht-Prosa behauptet mehr, als das Diagramm zeigt"

### F-7 — Die geschlossene `request_kind`-Menge liegt nach dieser Änderung an drei Orten, keiner davon generiert den DB-CHECK

- `kategorie`: INFO
- `quelle`: Maintainability · `AGENTS.md` §3.7 (Zwei-Quellen-Disziplin)
- `pfad`: `tools/schema/nacharbeit-administration.sql:58-59` ·
  `internal/domain/model/administrationrequest.go:21-26` ·
  `spec/pflichtenheft.md:310-332` (`SPEC-019`)
- `befund`: Der DB-CHECK wird nicht mehr aus dem deklarativen Modell
  erzeugt, sondern ist handgeschrieben; Domain-Konstanten,
  `SPEC-019` und die Nacharbeit-Datei müssen übereinstimmen. Gedeckt ist die
  Übereinstimmung heute nur mittelbar über Tests, die Zeilen mit allen vier
  Arten schreiben (die beiden Adapter-/E2E-Tests) — einen Sensor auf diese
  Paarung gibt es nicht (und nach `BEO-PGC/d-migrate-nacharbeit` ist der
  `functions:`-Weg als Ausweichform bewusst offen). Hinweis ohne erwartete
  Aktion an diesem Slice; relevant, falls die Menge ein weiteres Mal wächst.
- `verifizierbar`: teilweise — `make test`/`make test-store` würden eine
  divergente Klausel bei einer Einfügung sichtbar machen
- `klasse`: „Handgeschriebener DB-CHECK neben Modell-Konstanten und SPEC-Festlegung"

### F-8 — Test-Asymmetrie zwischen den beiden neuen Use-Case-Paketen

- `kategorie`: INFO
- `quelle`: Maintainability · `.harness/skills/reviewer.md` §Klassifikation
  (MEDIUM-Bullet „fehlende Negativtests", hier nicht erreicht — der Pfad ist
  in einem der beiden Pakete belegt)
- `pfad`: `internal/application/usecase/excludecolumn/service_test.go:76-83`
  gegenüber `internal/application/usecase/includecolumn/service_test.go`
- `befund`: `excludecolumn` trägt drei Testfälle (Happy, fehlende Spalte,
  durchgereichter Port-Fehler), `includecolumn` zwei — der
  Port-Fehler-Zweig (`return err`, sonst identischer Rumpf) hat im
  Einschluss-Paket kein Pendant. Der Zweig ist in der Ausschluss-Variante
  belegt, die beiden Dienste sind bis auf den Methodennamen gleich; keine
  erwartete Aktion an diesem Slice.
- `verifizierbar`: ja — `make test`
- `klasse`: „Asymmetrische Testabdeckung zwischen spiegelbildlichen Use Cases"

---

## Negativbefunde

- **geprüft, ohne Befund: Punkt 2 (der zusätzliche E2E-Test) — Scope und
  Pfad-Ausübung.** `internal/bootstrap/administration_endtoend_test.go` liegt
  im Scope: er testet `applyAdministrationRequest`, das dieser Slice laut
  §1/§3 selbst ändert, in genau dem Paket, dem die Funktion gehört; er
  duplziert die Fake-Ebene (`administration_internal_test.go`) nicht, sondern
  ergänzt die reale Ebene (SQL-Funktion → `ListPending` → echter
  `ColumnExclusionPort` → Status in der DB) und übt den versprochenen
  `applied`- **und** `failed`-Pfad real aus (zweite Hälfte: nicht
  existierende Spalte → `ErrSourceColumnMissing` → `failed` samt Fehlertext,
  gelesen aus `cdc.administration_request`). Er läuft ausschließlich unter
  `make test-store` (`CDC_STORE_TEST_DSN`) und skippt unter `make test`
  netzlos. Die im Dateikommentar benannte Ordnungs-Kopplung
  (`administration_endtoend_test.go` vor `replication_stream_test.go` und
  `walretention_endtoend_test.go`, die `DROP SCHEMA cdc CASCADE` fahren)
  trifft real zu: beide genannten Dateien existieren und setzen das Schema
  per `DROP SCHEMA … CASCADE` neu auf (`replication_stream_test.go:82`,
  `walretention_endtoend_test.go:92`); `make test-store` läuft mit dem
  gesetzten DSN grün. Die Abweichung von der im Plan genannten Belegdatei ist
  als F-4 klassifiziert, nicht als Scope-Creep.
- **geprüft, ohne Befund: Punkt 3 (DoD-Häkchen).** `2246f63` setzt genau vier
  Zeilen auf `[x]` — Schema-Erweiterung, SQL-Funktionen/Ports/Antrags-Verarbeitung,
  `make gates` grün, Doku-Update. Review, Closure-Notiz,
  Reconciliation-/Beobachtungs-Register, Risiko-Ausgänge und die drei
  Paarungen bleiben offen; keine überschrittene Checkbox
  (`BEO-PGC/dod-checkbox-nachzug` trifft hier nicht).
- **geprüft, ohne Befund: Punkt 4 (neuer Outbound Port, Pool, Layering).**
  `ColumnExclusionPort` deklariert nur `ColumnExists`
  (`internal/application/port/outbound/columnexclusion.go`), die
  Implementierung sitzt am `*TableActivationAdapter` (`tableactivation.go:71-93`)
  und arbeitet auf demselben `a.pool`: `activation` wird in
  `internal/bootstrap/wiring.go:417` einmal aus `cfg.AdminDSN` (`cdc_admin`,
  [`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md)) gebaut
  und dieselbe Instanz an die beiden neuen Use Cases gereicht
  (`wiring.go:471-472`, `:641-642`). Kein zweiter Pool, keine neue
  Berechtigungsklasse, kein neuer Grant; die Application-Schicht sieht
  ausschließlich das Port-Interface (`_ outbound.ColumnExclusionPort`-Assertion
  in der Adapterschicht). Die `REVOKE`/`GRANT`-Zeile der Nacharbeit nimmt die
  beiden neuen Signaturen mit (`nacharbeit-administration.sql:137-138`), und
  `administrationrequest_test.go` prüft die `cdc_reader`-Ablehnung real für
  `cdc.exclude_column` (SQLSTATE 42501). `make a-check` grün (0 Befunde).
- **geprüft, ohne Befund: Punkt 6 (Scope-Fidelity).**
  `git diff 529f021 2246f63 --stat` berührt 23 Dateien; `internal/adapters/driving/http`
  ist **nicht** darunter, ebenso wenig `compose.yaml`, `harness/README.md`,
  `.github/workflows/`, `docs/user/` oder ein Slice-Plan/-Wellen-Dokument
  außer dem eigenen. Die Dateien außerhalb der §3-Tabelle sind
  Folgeänderungen gelisteter Punkte (Domänenfeld/-fehlertext, zwei
  Testdateien, die generierten Rollout-Artefakte `down.sql`/`plan.yaml`) —
  die fehlende Nachziehung ist F-4.
- **geprüft, ohne Befund: Punkt 7 (`spec/architecture.md`, `spec/pflichtenheft.md`).**
  Die Sichtergänzung nennt kein `ADR-*`, keinen Slice und keine Welle
  (`AGENTS.md` §3.4); die neue Vierzeilen-Tabelle nennt Antragsart,
  SQL-Funktion und Inbound Port und deckt damit die Folgepflicht aus
  [`ADR-0059`](../plan/adr/0059-spaltenauswahl-mechanismus.md) real ab
  (die Verhaltensprüfung `applyAdministrationRequest` schaltet tatsächlich
  über `request.Kind` auf die vier Ports, `wiring.go:1026-1087`).
  `SPEC-019` beschreibt die Feldform deckungsgleich mit Modell und Nacharbeit
  (`column_name` nullable, Pflicht für die beiden Spaltenarten; vierwertige
  `request_kind`-Menge; `applied`/`failed`/`pending`; Fehlertextform
  „Spalte existiert nicht an der Quelle" + `schema.table.column`), nennt
  **keinen** `ADR-*`-Rückverweis (Spec→ADR-Verbot eingehalten) und
  die Änderungshistorie trägt die neue Zeile (`pflichtenheft.md:432`).
  Ein `Version:`-Kopf existiert in dieser Datei nicht (die
  Handbuch-HIGH-Regel des Skills betrifft `docs/user/benutzerhandbuch.md`).
- **geprüft, ohne Befund: Punkt 5 (Kommentar-Disziplin) außer F-1.** Die
  Prüfung lief als `grep` über die **hinzugefügten Zeilen** aller geänderten
  `.go`- und `tools/schema/`-Dateien
  (`slice-[0-9]+|welle-[0-9]+|bislang|zuvor|früher|jetzt|nicht mehr|vorher|ehemals|erstmals`
  sowie Konjunktiv-Muster für die verworfene Alternative): genau zwei
  Treffer. Der eine ist F-1; der zweite ist
  `docs/plan/planning/…/slice-066-…md:184` („die seit `slice-015` deklarativ
  gelöste Erstanlage-Konvergenz") — eine Herkunfts-Anker-Formulierung im
  Plan-Dokument, kein Produktionscode-Kommentar. Unberührt-gebliebene
  Alt-Treffer in denselben Dateien (`administrationrequest.go:138/214`,
  `wiring.go:150/773`, `administration_internal_test.go:27`) liegen
  außerhalb des Diff-Zeilenbereichs und sind nicht Gegenstand dieses
  Reports. Die neue `ColumnExists`-Doku (`tableactivation.go:71-83`) trägt
  Zusage und Grenze (Bezeichner-Alphabet für Schema/Tabelle, Spaltenname als
  Wert statt DDL-Text), keinen Konjunktiv über eine verworfene Alternative.
- **geprüft, ohne Befund: Traceability/ID-Schema.**
  `82ce83e`, `0d2030f`, `131fd98`, `2246f63` tragen je `LH-FA-CFG-005` und
  `ADR-0059` im Betreff, kein `SPEC-*`/`ARC-*` im Betreff **und** keins im
  Rumpf; `make commit-traceability` grün. Die neuen Kennungen folgen `MR-000`
  (`BEO-PGC/…`-Pfadform unverändert).
- **geprüft, ohne Befund: `tools/schema/{schema.yaml,nacharbeit-administration.sql}` —
  Idempotenz und Konvergenz der Ausweichform.** Siehe Eigene Nachmessung
  Schritte 1–3: Erstanlage Exit 0, dreimalige Wiederholung Exit 0 mit
  unverändertem Endzustand, CHECK-Messung exakt reproduziert. Die Aufteilung
  „Spalte deklarativ, Constraint in Nacharbeit" ist konsistent mit der
  etablierten Tradition (die vier Funktionen und die drei Views liegen aus
  demselben Grund dort); ihre *Konsequenz* für den Bestands-Rollout ist F-2.
- **geprüft, ohne Befund: Verhalten der Antrags-Verarbeitung gegen die
  Reihenfolge-Zusagen.** `applyAdministrationRequest` trägt für die beiden
  neuen Arten nichts in den `Assembler` nach und fällt nicht in den
  `default`-Fehlerzweig (`wiring.go:1069-1082`; der Fehlertext nennt jetzt die
  vierwertige Menge); der Fehlerpfad bleibt beim `MarkFailed`-Muster der
  bestehenden Arten (`BEO-PGC/adapter-fehler-ausgang` bekommt damit keinen
  dritten Fund in diesem Diff — die Use Cases reichen Adapter-Fehler
  unverändert durch). Der Zwischenzustand „Antrag `applied`, aber noch keine
  Filterwirkung" ist in §1 des Plans ausdrücklich als Zustand dieses Slice
  benannt (`slice-067`/`slice-068` schließen ihn) und daher kein Finding.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 2 |
| LOW | 2 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Slice-Chronik in Produktionscode-Kommentar ·
Beleg-Aussage im Plan nicht reproduzierbar (Bestands-Instanz) · Neuer
Invarianten-Zweig am öffentlichen Konstruktor ohne Negativtest ·
Plan-Tabelle/Beleg-Angabe nicht auf den tatsächlichen Umfang nachgezogen ·
Kopplungs-Kommentar neben geändertem Artefakt nicht nachgezogen ·
Sicht-Prosa behauptet mehr, als das Diagramm zeigt · Handgeschriebener
DB-CHECK neben Modell-Konstanten und SPEC-Festlegung · Asymmetrische
Testabdeckung zwischen spiegelbildlichen Use Cases

## Verdikt

**Merge-blockierend:** ja — F-1 ist HIGH, F-2/F-3 sind MEDIUM. Die
**Richtung** des Slice ist nicht bestritten: die d-migrate-Grenze ist real
und unabhängig nachgestellt (Exit 5, Klauselverlust, konvergente Spalte), die
Ausweichform ist idempotent, der neue Port nutzt denselben `cdc_admin`-Pool,
die DoD-Häkchen sitzen genau auf den vier umsetzungsbezogenen Zeilen, und
`make gates`/`make test`/`make test-store` laufen grün (Exit 0).

**Urteil zur Ausweichform (Punkt 1) und zur „neuen Objektklasse":**
Die Ausweichform trägt und ist begründet; die Messung des Implementers zum
CHECK ist reproduziert und die Aufteilung „Spalte deklarativ, Constraint in
Nacharbeit" deckt sich mit der seit `slice-011`/`slice-012`/`slice-036`
etablierten Tradition. Ich teile die Einschätzung als **neue Objektklasse**
mit einer Präzisierung: es ist kein neuer *Mechanismus* — es ist derselbe
`POST_EXECUTE_DRIFT`-Familienfall wie die Funktionsklasse, aber ein anderes
Objekt (Änderung/Anlage an einer **bestehenden** Tabelle statt der seit
`slice-015` deklarierten Erstanlage); das ist bisher an keiner Stelle belegt
und verlangt den Beleg-Eintrag zu `BEO-PGC/d-migrate-nacharbeit` bei der
Closure, wie §3 des Plans ihn bereits ankündigt. **Ein Problem verdeckt die
Aufteilung allerdings dort, wo der Plan sie als „Exit 0"-Beleg führt**: die
Ausweichform macht ein Objekt, das d-migrate im Katalog vorfindet und im
Modell nicht mehr deklariert, für den Bestands-Rollout zum *destruktiven*
Plan-Eintrag (F-2), und der reale Ausgang dieser Konstellation (Exit 8,
Abbruch vor den psql-Schritten) steht nirgends — während er für einen
frischen Rollout unverändert Exit 0 liefert.

**Übergabe:** Fixrunde am Implementer (F-1, F-2, F-3). F-4/F-5 sind
Plan-/Kommentar-Nachzüge, F-6 bis F-8 sind Hinweise ohne erwartete Aktion.
Da eine Rückgabe an den Implementer erfolgt, bleibt die DoD-Zeile „Review
durchgeführt, Report unter `docs/reviews/` liegt vor" in §2 offen
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde greift
nicht); sie wird regulär bei Schritt 21 des Implementer-Workflows nachgezogen.
Dieser Report ist Lauf-Beleg und wird über Läufe hinweg nicht erneut gelesen;
Verifikation gegen DoD/Spec bleibt Aufgabe des Verifiers.
