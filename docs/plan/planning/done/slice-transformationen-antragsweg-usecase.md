# Slice transformationen-antragsweg-usecase: Antragsweg, Verarbeitung — Use Cases `SetTransformation`/`RemoveTransformation` mit K1–K4, Regelstand-Port, Verdrahtung und Dauerhaftigkeit

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-transformationen](../welle-transformationen.md).

**Bezug:** [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) (Happy Path,
Boundary), [`LH-FA-ADM-001`](../../../../spec/lastenheft.md) (Administration
über SQL),
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Teilfrage 1/3/6 (Konfliktfreiheit, Dauerhaftigkeit) und Folgepflicht 3,
[`ADR-0065`](../../adr/0065-spaltenausschluss-dauerhafter-traeger.md)
(dauerhafter, tabellen-scoped Träger — Muster),
[`ADR-0059`](../../adr/0059-spaltenauswahl-mechanismus.md) Teilfrage 5
(Existenz der Spalte an der Quelle),
[`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(Live-Reload), [`ADR-0028`](../../adr/0028-inbound-use-cases.md) (Inbound Use
Cases), [`ADR-0034`](../../adr/0034-ports-nach-faehigkeiten.md) (Ports nach
Fähigkeiten),
[`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
(keine Domänenlogik in SQL).

**Berührte Spec-Stellen:** [`SPEC-019`](../../../../spec/pflichtenheft.md)
(Antrags-Datensatz, Fehlertexte, durch `slice-transformationen-spec-nachzug`),
[`ARC-002`](../../../../spec/architecture.md),
[`ARC-003`](../../../../spec/architecture.md),
[`ARC-004`](../../../../spec/architecture.md),
[`ARC-007`](../../../../spec/architecture.md) — gelesen, nicht geändert.

**Verantwortlich:** Implementer-Agent, 2026-09-26.

**Autor:** Planner-Agent, Welle-Eröffnung
[welle-transformationen](../welle-transformationen.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein `set_transformation`-/`remove_transformation`-Antrag wird
verarbeitet und wirkt dauerhaft: zwei Use Cases prüfen die
Konfliktfreiheits-Invarianten K1–K4 (ein Verstoß endet `failed` mit dem
Fehlertext der Spec, der Regelstand bleibt unverändert), ein neuer Outbound
Port liefert die Regelstand-Ableitung und die Spaltenliste der Quelltabelle,
`applyAdministrationRequest` trägt die Regel in die laufende
`Assembler`-Bindung nach, und der aus den `applied`-Zeilen abgeleitete
Regelstand geht bei jedem Pfad, der eine Bindung anlegt (Prozessstart über
`activatedTableBindings`, Aktivierungs-Zweig), in die Bindung ein — Muster von
[`ADR-0065`](../../adr/0065-spaltenausschluss-dauerhafter-traeger.md).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Schema, Funktionen, Grants, Guard** — `antragsweg-schema`; dieser Slice
  setzt die Spalten `rule_name`/`rule_spec` voraus.
- **Wirkung im Assembler** — `kern-rename` (die Methoden
  `SetTransformation`/`RemoveTransformation` existieren dort; dieser Slice ruft
  sie).
- **Der zweite Regeltyp** — `map-value`; der Use Case prüft gegen den
  Regeltyp-Satz der Domäne und bleibt für einen weiteren Typ unverändert (das
  ist die Aussage von
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Folgepflicht 4, ihr Beleg steht dort).
- **Der Backfill-Pfad** — `backfill-pfad` liest den hier gelieferten
  Regelstand-Port.
- **Die Startreihenfolge** — `start-reihenfolge`; hier bleibt es bei der
  bestehenden Reihenfolge (Goroutine neben `stream.Run`).
- **Ein weiterer Lese- oder Anzeige-Weg des Regelstands** (Sicht,
  `diagnose`-Ausgabe) —
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 6: nicht Teil dieser Entscheidung.

## 2. Definition of Done

- [x] K1–K4 gegen Fakes, je an ihre Eingabe gebunden: (K1) ein bereits
      vergebener `rule_name` derselben Tabelle endet `failed`; (K2) eine zweite
      Spaltenregel auf derselben Quellspalte endet `failed`; (K3) ein Zielname,
      der einem anderen Zielnamen oder einem Spaltennamen der Quelltabelle
      (ausgeschlossene Spalten eingeschlossen) gleicht, endet `failed`; (K4)
      eine `column`, die an der Quelle nicht existiert, und ein
      `remove_transformation` gegen einen nicht geführten Namen enden `failed`;
      ein unbekannter `kind` und ein unbekannter Schlüssel in `rule_spec` enden
      `failed` (strikte Dekodierung); der Fehlertext entspricht der Spec, der
      Regelstand bleibt nach jedem `failed` unverändert. Die Prüfreihenfolge
      ist die der Spec ([`SPEC-019`](../../../../spec/pflichtenheft.md): fünf
      Formzeilen, dann K1 bis K4): eine Verletzung der Form von `column`/`to`
      in der Domäne (`kern-rename` §1) wird auf `rule_spec ist ungültig`
      abgebildet, ein Sentinel für `to` gleich `column` auf den K3-Text `Zielname
      kollidiert mit einer Spalte der Tabelle` in der Stellung von K3 — nie auf
      `rule_spec ist ungültig`; ein leerer, fehlender oder ungültiger Regelname
      endet `failed` mit `Regelname ist ungültig`, die Stelle dieser Prüfung
      (Konstruktor-Sentinel oder Vorprüfung im Use Case) legt dieser Slice fest,
      und ein Test bindet je Text den Auslöser. Übergabe aus
      `antragsweg-schema` (gemessen im Review
      `review-slice-transformationen-antragsweg-schema`): die Funktionen
      schreiben eine Zeile mit `rule_name` NULL oder leer und mit `rule_spec`
      SQL-NULL, ohne zu prüfen; der Konstruktor-Sentinel in
      `ReadPendingRequests` lehnt diese Zeile **beim Lesen** ab, `ListPending`
      liefert den Fehler, und `processAdministrationRequests` liest dieselbe
      Zeile im nächsten Durchlauf erneut — kein Antrag der Queue wird
      verarbeitet, bis die Zeile entfernt ist. [`SPEC-019`](../../../../spec/pflichtenheft.md)
      sagt für diese Fälle `failed` mit `Regelname ist ungültig` bzw.
      `rule_spec ist ungültig`: die Zeile wird gelesen und verarbeitet statt
      abgelehnt. Der Slice legt fest, an welcher Stelle geprüft wird (Lesen oder
      Verarbeiten) und bindet den Fall mit einem Test, der eine Queue mit einer
      ungültigen und einer gültigen Zeile durchläuft (die gültige Zeile wird
      verarbeitet); derselbe Lese-Pfad lehnt am Parent eine Spaltenart mit
      leerer `column` ab (`cdc.exclude_column(…, NULL)`), der Slice nennt, ob
      seine Stelle der Prüfung sie mitführt oder die Grenze
      (`BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue`, geplant, 1×).
      Zu den Eingaben des Falls gehört `rule_spec` in den
      Formen, die der Aufruf annimmt ([`ADR-0126`](../../adr/0126-transformationen-annahmemenge-rule-spec.md)
      Festlegung 1): SQL-`NULL` (der Store liefert den leeren Text),
      JSON-`null` (der Store liefert den Text `null`), ein Wert ohne Objekt und
      ein doppelter Schlüssel — `jsonb` normalisiert vor Go
      (`{"kind":"a","kind":"b","z":1,"a":2}` wird als
      `{"a": 2, "z": 1, "kind": "b"}` gelesen, gemessen im Review), die strikte
      Dekodierung sieht einen doppelten Schlüssel deshalb nicht. *Zu belegen durch:*
      `make test` und je Invariante eine Mutation der Prüfung, die den Test rot
      färbt (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`, verkörpert,
      6×).
- [x] Regelstand-Ableitung und Spaltenliste im Store: der Adapter liefert den
      Regelstand je Tabelle einer Quelle aus den `applied`-Zeilen der zwei
      Arten in der Ordnung `requested_at`, bei gleichem Zeitstempel nach
      `administration_request_id` (`set_transformation` trägt ein,
      `remove_transformation` nimmt heraus), und die Spaltennamen der
      Quelltabelle aus dem Katalog; die Faltung der Zeilen ist eine reine
      Funktion an einer Stelle, die in der netzlos gemessenen Fläche liegt
      (`BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code`, offen, 2×; der
      Ausschluss des Coverage-Gates trifft nur die Pakete `postgresstorage`,
      `postgresack`, `postgressnapshot` und `replication/receive` selbst,
      gelesen am Dockerfile).
      *Zu belegen durch:* `make test-store` (gleicher Zeitstempel, Reihenfolge,
      Set/Remove-Zyklus, Katalog-Lesart) und `make test` (Faltung).
- [x] Verdrahtung und Dauerhaftigkeit: ein `set_transformation`-Antrag gegen
      eine aktivierte Tabelle wird `applied`, die laufende `Assembler`-Bindung
      trägt die Regel ohne Neustart; `remove_transformation` nimmt sie wieder
      heraus; der Prozessstart (`activatedTableBindings`) und der
      Aktivierungs-Zweig (Deaktivierung, dann Aktivierung) tragen den
      abgeleiteten Regelstand in die neue Bindung; die erneute Verarbeitung
      eines bereits nachgetragenen, noch `pending` stehenden Antrags ist
      idempotent (der Nachtrag ersetzt nach Namen, K1 prüft nur gegen
      `applied`-Zeilen); ein Antrag gegen eine Tabelle ohne laufende Bindung
      endet gemäß Spec. *Zu belegen durch:* `make test` (Whitebox in
      `internal/bootstrap` mit Fakes: Nachtrag, Prozessstart,
      Aktivierungs-Zweig, Wiederholung) und `make test-store` (Ableitung gegen
      reale Zeilen).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: entfällt — die Betreiber-Oberfläche ist mit
      `antragsweg-schema` entstanden und mit
      `slice-transformationen-betriebsdoku` adressiert; dieser Slice ändert
      keinen Vertrag, den das Handbuch beschreibt, ohne dass die Wirkung erst
      mit `e2e-wirkung` belegt ist.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure dieses Slice (§7) und zusätzlich von der Closure der Welle
      [welle-transformationen](../welle-transformationen.md) (die Roadmap führt
      sie unter *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/application/port/inbound/transformation.go` | neu | `SetTransformationUseCase` (liefert die geprüfte `model.Transformation`, die der Aufrufer in die Bindung einträgt) und `RemoveTransformationUseCase` samt `SetTransformationCommand`/`RemoveTransformationCommand` ([`ADR-0028`](../../adr/0028-inbound-use-cases.md), Transport-Typen am Port); die Namen der Use Cases stehen in [`ARC-003`](../../../../spec/architecture.md) (Antragsart-Tabelle). |
| `internal/application/port/outbound/transformation.go` | neu | **ein** Outbound Port `TransformationPort` mit zwei Lese-Methoden (`TransformationRules`: Regelstand je Tabelle einer Quelle; `SourceColumns`: Spaltennamen der Quelltabelle) — ein Port, wie [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) zählt („ein neuer Outbound Port“). Der Schnitt folgt dem Vorbild `ColumnExclusionPort` (`ColumnExists` aus dem Katalog und `ExcludedColumns` aus den Antrags-Zeilen in einem Port); beide Methoden beantworten dieselbe Frage („welche Regeln dürfen die Spalten dieser Tabelle tragen“) und laufen über dieselbe Adapter-Instanz. Ein Port mit zwei Lese-Methoden ist mit [`ADR-0034`](../../adr/0034-ports-nach-faehigkeiten.md) (Ports nach Fähigkeiten) vereinbar: beide Lesarten gehören zur Frage der Konfliktprüfung, `ColumnExclusionPort` schneidet beide Objektklassen ebenfalls in einen Port, und die Konsistenzgrenze der ADR trifft hier keine gemeinsame Transaktion, weil beide Methoden lesen (Review `review-slice-transformationen-antragsweg-usecase`, Negativbefund Port-Schnitt, Finding F-5). |
| `internal/application/usecase/settransformation/service.go`, `internal/application/usecase/removetransformation/service.go` (+ je `service_test.go`) | neu | die zwei Use Cases in der netzlos gemessenen Fläche: `Set` prüft die Formzeilen (Regelname, Regelform) **vor** dem ersten Lesen des Ports, dann K1 bis K3 über `TransformationSpec.CheckConflicts`, dann K4 gegen die Spaltenliste; `Remove` prüft den Namen und K4 (`Regelname nicht geführt`). Der Fehlertext ist der der Spec (Klartext, Doppelpunkt, Adresse), der Grund bleibt über `errors.Is` erreichbar. Beide Use Cases schreiben den Regelstand nicht (er ist die Ableitung aus den `applied`-Zeilen). |
| `internal/domain/model/transformationspec.go` (+ `transformationspec_test.go`) | neu | Parser `ParseTransformationSpec` (strikt: UTF-8, JSON-Objekt, `kind`, unbekannter Regeltyp, unbekannter Schlüssel in aufsteigender Ordnung, Pflichtschlüssel und Bezeichner-Form über `NewRenameColumn`), `CheckRuleName` (Alphabet `[a-z0-9_]{1,63}`), `TransformationSpec.CheckConflicts` (K1 bis K3, zeichengenau, in der Reihenfolge der Spec), `Build` und die reine Faltung `FoldTransformations` (`applied`-Zeilen → Regelstand; Set trägt ein und ersetzt nach Namen, Remove nimmt heraus; eine nicht mehr lesbare Zeile endet als Fehler). Ein Zielname gleich der Quellspalte ist keine Formverletzung des Parsers, sondern K3 (`ErrTargetCollidesWithColumn`) in der Stellung von K3. |
| `internal/domain/model/administrationrequest.go` (+ Test), `internal/domain/errors/errors.go` | update | der Konstruktor `NewAdministrationRequest` lehnt die Transformations-Antragsarten nicht mehr wegen leerem Regelnamen oder leerer Regelform ab (**Stelle der Prüfung: Verarbeiten im Use Case, nicht Lesen**); `AdministrationRequestKinds()` zählt die sieben Arten auf; neun Sentinels tragen die Ablehnungsgründe (Klartext der Spec-Zeilen). |
| `internal/adapters/driven/postgresstorage/queries/queries.go`, `sqlexec/translate.go`, `tableactivation.go` (+ `sqlexec/translate_test.go`, `administrationrequest_test.go`) | update | `SelectAppliedTransformationRequests` (Ordnung `requested_at`, dann `administration_request_id`), `SelectTableColumns` (`information_schema.columns`, `ORDER BY ordinal_position`); `sqlexec.ReadTransformationRules` (Zeilen je Tabelle, Faltung in der Domäne, eine Tabelle ohne verbleibende Regel trägt keinen Eintrag) und `ReadSourceColumns` in der netzlos gemessenen Fläche; `TableActivationAdapter` implementiert den Port (dieselbe Instanz über `CDC_ADMIN_DSN`, `SELECT` der Rolle `cdc_admin` auf der Antrags-Tabelle wie `ExcludedColumns`; `cdc_capture` trägt dort kein Recht). Die Faltung liegt an **einer** Stelle, `model.FoldTransformations`. |
| `internal/bootstrap/wiring.go` | update | Zweige `set_transformation`/`remove_transformation` in `applyAdministrationRequest` (Use Case, danach `Assembler.SetTransformation`/`RemoveTransformation`; der Kommentar am Set-Zweig nennt die Grenze von K3), drei Felder in `administrationDeps` (`transformations`, `setTransformations`, `removeTransformations`), `activatedTableBindings` und der Aktivierungs-Zweig tragen den abgeleiteten Regelstand neben dem Ausschlussstand, die Doc-Kommentare zählen die Arten und beschreiben die Zweige; `processedAdministrationKinds` ist eine Funktion über `model.AdministrationRequestKinds()` (eine Quelle statt einer handgeführten Zeichenkette). |
| `internal/bootstrap/administration_internal_test.go`, `wiring_rest_internal_test.go`, `backfill_internal_test.go` | update | Whitebox mit Fakes: Nachtrag live (Set, Remove), K1 bis K4 und die Formzeilen bis zum `failed`-Vermerk mit dem Fehlertext der Spec, ungültige Zeilen neben einer gültigen, Idempotenz der Wiederholung, Tabelle ohne Bindung, Prozessstart (`activatedTableBindings`), Aktivierungs-Zweig (Deaktivierung, dann Aktivierung), Lesefehler des Regelstandes, `default`-Zweig und die Bindung der Aufzählung an die Fälle des `switch` (`TestApplyAdministrationRequestHandlesEveryKindOfTheClosedSet`); die bestehenden Fixtures tragen den Fake des neuen Ports. Die zwei Tests, die das Fenster „Funktion vorhanden, Wirkung fehlt“ banden (`…RejectsUnprocessedKind`, `…MarksTransformationRequestsFailed`), sind ersetzt. |
| `internal/bootstrap/administration_roles_internal_test.go` | update (Übergabe aus `slice-backfill-run-store`) | der Login-Test zieht die zwei neuen Arten unter dem `cdc_admin`-/`cdc_capture`-Login durch die Queue: Set, Remove, K1 bis K4 mit Fehlertext, Regelstand und Spaltenliste über den realen Adapter, eine gültige Zeile neben Zeilen mit fehlendem Regelnamen und fehlender Regelform (die Queue bleibt lesbar). |
| `.dockerignore` | prüfen — keine Änderung | alle neuen Pfade liegen unter `internal/` (`!internal/`); keine Datei außerhalb des Go-Baums kommt hinzu. Belegt durch `make coverage-gate` und `make test` (Bau-Kontext trägt die Pakete). |
| `spec/pflichtenheft.md` `SPEC-019` (Zeile „`rule_name` ist … Pflicht … (Domänen-Invarianten des Antrags-Konstruktors)“) | gemeldet, mit der Closure nachgezogen | die Klammer nannte den Konstruktor als Ort der Prüfung; nach diesem Slice liegt sie im Use Case (Formzeile `Regelname ist ungültig`). Fremde Datei; Meldung an den Planner, Frist: Closure dieses Slice — der Planner ersetzte die Klammer durch die wahre Aussage (Verifikation V-4) und trug die Historie-Zeile nach (Commit der Closure). |
| Spalten-Antragsarten mit leerer Spalte (`cdc.exclude_column(…, NULL)`) | keine Änderung (Grenze, benannt) | der Konstruktor lehnt `exclude_column`/`include_column` mit leerer Spalte weiterhin beim Lesen ab (`SPEC-019`: „Domänen-Invariante des Antrags-Konstruktors“ für die Spalten-Antragsarten); die Stelle der Prüfung dieses Slice führt sie **nicht** mit. Dieselbe Grenze gilt für jede Antragsart mit leerem Schema oder leerem Tabellennamen (`ErrEmptyIdentifier`, Ist-Verhalten erprobt: `TestReadPendingRequestsRejectsRowWithEmptySchemaOrTable`, siehe §6). Adresse: `BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue` (geplant, 1×), Entscheidung beim Planner. |
| `internal/adapters/driven/postgresstorage/queries/queries.go`, `internal/application/port/outbound/administrationrequest.go` (Fixrunde, Review F-1) | update | `SelectPendingAdministrationRequests` ordnet `requested_at, administration_request_id` — dieselbe Ordnung wie die zwei Abfragen der Ableitung; Verhaltensänderung nur bei gleichem `requested_at`, für **alle** sieben Antragsarten (auch `exclude_column`/`include_column` des Parents); Kommentare der Abfrage und von `ListPending` nennen die Ordnung. |
| `internal/adapters/driven/postgresstorage/administrationrequest_order_test.go` (Fixrunde, Review F-1) | neu | zwei Store-Tests (`make test-store`): `TestAdministrationRequestListPendingOrdersTiesByRequestID` (vier gleichzeitige Anträge absteigend eingefügt, Ordnung = Kennung; Regelstand und Ausschlussstand der Verarbeitung gleich der Ableitung) und `TestAdministrationRequestSameTransactionCallsKeepCallOrder` (Fixrunde 2, [`ADR-0127`](../../adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md), Fixrunde 3 auf alle sieben Funktionen erweitert: über die realen Funktionen je Durchlauf zwei Transaktionen auf eigenen Tabellen — die Vorwärts-Folge `backfill_table`, `remove_transformation`, `set_transformation`, `exclude_column`, `include_column`, `disable_table`, `enable_table` und ihre Umkehrung —, 60 Durchläufe; `ListPending` liefert die Aufruf-Reihenfolge, Regelstand und Ausschlussstand der Verarbeitung gleichen der Ableitung und dem Stand der Aufrufe). |
| `internal/adapters/driven/postgresstorage/administrationrequest_order_test.go`, `internal/bootstrap/wiring.go` (Fixrunde 3, Verifikation V-2, V-3, V-5) | update | V-2: der Store-Test bindet `clock_timestamp()` an alle sieben Funktionen in beiden Richtungen (Folge und Umkehrung, `backfill_table` an erster und an letzter Stelle); V-3: der Doc-Kommentar des Tests nennt die von der Mutationstabelle getragene Aussage (die Mutation färbt den Test rot, sobald die Funktion einen Vorgänger hat); V-5: der Doc-Kommentar von `activatedTableBindings` nennt die Fehlerklasse `internal` der Faltung statt „wie ein Lesefehler“. Kein Produktivcode. Nicht Teil der Fixrunde: V-1 (`ADR-0127`, keine Änderung an der `Accepted` ADR), V-4 (`SPEC-019`, Meldung an den Planner, Frist Closure), V-6, V-7. |
| `tools/schema/nacharbeit-administration.sql` (Fixrunde 2, [`ADR-0127`](../../adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md) Folgepflicht 2) | update | die sieben Funktionen schreiben `requested_at` je Aufruf mit `clock_timestamp()` (zweite Spalte, zweiter Wert des `INSERT`; Signaturen unverändert, nur `CREATE OR REPLACE`-Text); der Kopfkommentar nennt die Zeitstempel-Vergabe. Kein zweiter Träger des Funktionstextes: `git grep -n gen_random_uuid -- tools internal examples cmd` trifft nur diese Datei, `plan.yaml`/`down.sql` tragen keine Funktion (schema.yaml und der Spalten-Default bleiben unverändert). |
| `internal/adapters/driven/postgresstorage/queries/queries.go`, `administrationrequest.go`, `sqlexec/translate.go`, `internal/application/port/outbound/administrationrequest.go` (Fixrunde 2, [`ADR-0127`](../../adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md) Folgepflicht 2) | update | die Kommentare nennen `requested_at` als Aufrufzeitpunkt der schreibenden Funktion (`SelectPendingAdministrationRequests`, `SelectAppliedColumnRequests`, `SelectAppliedTransformationRequests`) und die Ordnung „Aufruf-Reihenfolge“ (`ListPending`, `ReadPendingRequests`, Port). Kein Verhalten. |
| `internal/bootstrap/administration_callorder_internal_test.go` (Fixrunde 2, [`ADR-0127`](../../adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md)) | neu | Whitebox-Test mit realer Queue (`make test-store`): `TestAdministrationSameTransactionRemoveThenSetLeavesTheNewRuleLiveAndDerived` — je 50 Tabellen ein `remove_transformation`/`set_transformation` derselben Regel in einer Transaktion über die realen Funktionen, verarbeitet von `processAdministrationRequests`; die laufende Bindung trägt danach den Zielnamen der neuen Regel, und eine aus `TransformationRules` neu gebildete Bindung liefert dasselbe Bild. |
| `docs/plan/planning/open/slice-transformationen-betriebsdoku.md` (Fixrunde 2, [`ADR-0127`](../../adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md) Folgepflicht 4) | update | Übergabe an den Handbuch-Slice: dessen §2 trägt die Aussage zur Aufruf-Reihenfolge; das Handbuch selbst bleibt unberührt. |
| `internal/adapters/driven/postgresstorage/sqlexec/translate_test.go` (Fixrunde, Review F-6) | update | `TestReadPendingRequestsRejectsRowWithEmptySchemaOrTable`: Ist-Verhalten der Restklasse (leeres Schema, leere Tabelle enden beim Lesen als `ErrEmptyIdentifier`); kein Code-Fix, Grenze in §6. |
| `internal/bootstrap/wiring.go` (Fixrunde, Review F-2, F-3, F-7) | update | der Kommentar am Aufruf von `activatedTableBindings` in `Run` nennt den wahren Zustand (die Seeds erreichen den Assembler nie); der Kommentar am Set-Zweig nennt die Grenze des fehlgeschlagenen Vermerks mit kollidierendem Folgeantrag; `gofmt`-Form der Felder von `administrationDeps`. |
| `internal/bootstrap/administration_internal_test.go` (Fixrunde, Review F-4) | update | der Doc-Kommentar von `TestProcessAdministrationRequestsSetTransformationIsIdempotent` nennt, was der Test bindet (den Fake) und wo die Idempotenz getragen ist (Store-Abfrage, Login-Test). |
| `internal/domain/model/transformationspec.go`, `internal/adapters/driven/postgresstorage/administrationrequest_test.go` (Fixrunde, Review F-3) | update | reine `gofmt`-Form (Liste im Doc-Kommentar, Ausrichtung eines `const`-Blocks); Docker-only über das gepinnte Toolchain-Image, kein neues make-Ziel. |
| Handbuch `docs/user/benutzerhandbuch.md` | keine Änderung | wie DoD Punkt 8; Kandidatenlauf im Bericht. |

**Übergaben aus `slice-transformationen-kern-rename`** (gemeldet, kein zusätzlicher Umfang):

- **Regelliste ohne geteilten Speicher.** `NewAssembler` und `AddBinding` teilen die übergebene
  Liste mit dem Aufrufer; der Doc-Kommentar von `TableBinding` sagt „ab dem Schreiben
  unverändert“ zu. Die Ableitung aus den `applied`-Zeilen baut je Lesung eine frische Liste und
  schreibt sie nach der Übergabe nicht mehr.
- **Gültiges UTF-8 als Eintritts-Eigenschaft.** `NewRenameColumn` prüft Spalte und Zielname nicht
  auf UTF-8; `rule_spec` als `jsonb` liefert gültiges UTF-8 (hergeleitet aus dem Typ, nicht
  erprobt). Der Parser belegt es an einem Test oder führt die Prüfung selbst.
- **K3 zwischen Regeln.** Die Domäne trägt einen zweiten, zeilenabhängigen Wächter gegen zwei
  Regeln mit gleichem Zielnamen (`ErrTransformationTargetCollides` aus `BuildRowImage`); die
  Prüfung beim Antrag bleibt Sache dieses Slice, ihr Test „zwei Regeln mit gleichem Zielnamen
  endet `failed`“ zeigt, dass der Betrieb den Wächter nicht erreicht (Register
  `BEO-PGC/wertabhaengiger-zweiter-waechter-ohne-spec-zeile`).

**Übergabe aus `slice-transformationen-antragsweg-schema`** (gemeldet, kein zusätzlicher Umfang):

- **Die verarbeitete Menge steht zweimal.** `processedAdministrationKinds` in
  `internal/bootstrap/wiring.go` ist eine handgeführte Zeichenkette neben dem
  `switch` von `applyAdministrationRequest`; `TestApplyAdministrationRequestRejectsUnprocessedKind`
  bindet die Konstante an den Fehlertext, kein Test bindet sie an die Fälle des
  `switch`. Der Slice zieht die Konstante nach (der genannte Test färbt sich rot,
  sobald `set_transformation` verarbeitet wird) und bindet die Aufzählung an die
  Fälle des `switch` oder benennt die Grenze.
- **Die Mutationsangabe in der Fitness-Function-Zeile von
  [`ADR-0126`](../../adr/0126-transformationen-annahmemenge-rule-spec.md) ist
  falsch** (gemessen im Verifikations-Report: `p_rule_spec::jsonb` → `p_rule_spec`
  färbt keinen Test, weil der Zuweisungs-Cast `json` → `jsonb` dieselbe
  Umwandlung leistet; wirksam ist `NULL::jsonb`). Die `Accepted` ADR bleibt
  unverändert. Berührt dieser Slice oder ein Architect-Zug zu ihm `rule_spec`
  (Parametertyp, Cast, PostgreSQL-Hauptversion — die Re-Evaluierungs-Trigger von
  `ADR-0126`), trägt die dabei entstehende ADR die Berichtigung als eigene
  Klausel (`BEO-PGC/adr-aussage-breiter-als-ihre-messung`, siebtes Auftreten);
  ein eigener ADR-Zug allein dafür entfällt.

**Behandlung der Übergaben** (kein zusätzlicher Umfang; je Übergabe der Beleg):

- *Regelliste ohne geteilten Speicher* — `FoldTransformations` und `ReadTransformationRules` bauen je Lesung frische Listen; `activatedTableBindings` und der Aktivierungs-Zweig geben sie unverändert an die Bindung, niemand schreibt sie danach (Doc-Kommentare an beiden Stellen). Grenze: kein Test bindet die Nicht-Teilung, sie folgt aus dem Aufbau der Faltung (`replaceOrAppend`/`dropByName` legen neue Listen an).
- *Gültiges UTF-8* — der Parser prüft es selbst (`utf8.ValidString`); `TestParseTransformationSpecRejectsInSpecOrder` bindet den Fall „ungültiges UTF-8“.
- *K3 zwischen Regeln* — `TestSetTransformationRejectsWithTheSpecTexts` (Zeile „K3 Zielname gleicht dem Ziel einer anderen Regel“) und `TestProcessAdministrationRequestsRuleViolationsFailWithSpecTexts` (Zeile `req-k3a`) enden `failed`; im zweiten Test baut derselbe Assembler danach weiter Bilder (`assemblerRowImage` bricht bei einem Fehler des Erfassungspfads ab): der Betrieb erreicht den Wächter `ErrTransformationTargetCollides` nicht.
- *Die verarbeitete Menge steht zweimal* — gelöst: `processedAdministrationKinds()` leitet sie aus `model.AdministrationRequestKinds()` ab; `TestApplyAdministrationRequestHandlesEveryKindOfTheClosedSet` bindet jede Art der Menge an einen Zweig des `switch` (Grenze im Kommentar des Tests: eine Art, die der Konstruktor annimmt und die Aufzählung nicht nennt, sieht der Test nicht).
- *Die Mutationsangabe in `ADR-0126`* — berührt dieser Slice nicht: weder Parametertyp noch Cast noch PostgreSQL-Hauptversion ändern sich; kein ADR-Zug.
- *Fenster „Funktion vorhanden, Wirkung fehlt“* — mit diesem Slice geschlossen: die zwei Tests des `default`-Zweigs für `set_transformation` sind durch die Wirkungs-Tests ersetzt (Mutation „Zweig entfällt“ in der Tabelle unten).
- *Die Stelle der Prüfung (Übergabe aus `antragsweg-schema`)* — **Verarbeiten**: `NewAdministrationRequest` nimmt Regelname und Regelform, wie die Zeile sie hält; `CheckRuleName` und `ParseTransformationSpec` im Use Case bestimmen den Fehlertext. `TestReadPendingRequestsCarriesRuleRowsWithEmptyRuleFields` (netzlos), `TestAdministrationRequestListPendingCarriesRuleRowsWithMissingFields` (reale Zeilen mit NULL-Regelname, SQL-NULL, JSON-`null`), `TestProcessAdministrationRequestsInvalidRuleRowsDoNotStallTheQueue` und der Login-Test (gültige Zeile hinter ungültigen wird `applied`) binden den Fall. Die Formen von `rule_spec` der Übergabe: SQL-NULL (leerer Text), JSON-`null` (Text `null`) und Wert ohne Objekt enden `rule_spec ist ungültig`; ein doppelter Schlüssel erreicht Go nicht (`jsonb` normalisiert vor dem Lesen, Kommentar am Parser).

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „die Menge der
Antragsarten in `applyAdministrationRequest`“, „die Menge der Pfade, die eine
Bindung anlegen und den Ausschlussstand mitführen“, „die Felder von
`administrationDeps`“; beide Stände gemessen, Parent-Stand `80eefead`):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Stellen, die eine Bindung anlegen | `git grep -n 'TableBinding{' -- internal ':!*_test.go'` (Zeilen 1–2 des Blocks) | **Gefunden.** Parent 5, Diff 5 (Zeilenzahl unverändert, Inhalt geändert): `config_file.go:250` und `parseTables` in `wiring.go` (die Seeds aus Konfigurationsdatei und `CDC_TABLES`: Bindung ohne Ausschluss- und Regelstand; sie gehen ausschließlich in die Aktivierung `enableTables.Enable` und erreichen den Assembler nie — gemessen mit `git grep -n 'cfg.Tables' -- internal ':!*_test.go'`: die Zuweisung in `config_file.go`/`wiring.go`, die eine Schleife in `Run` und `receive.Config.Tables`, das der Prozessstart aus `activatedTableBindings` füllt), `activatedTableBindings` und der Aktivierungs-Zweig (beide tragen den Regelstand mit). **Nichtgefunden:** keine sechste Stelle, die eine Bindung anlegt; `Assembler.AddBinding` erhält Ausschluss- und Regelstand einer getragenen Bindung selbst (`kern-rename`). | Die Seeds tragen keinen Regelstand und brauchen keinen, weil sie den Assembler nicht erreichen (Kommentar am Aufruf von `activatedTableBindings` in `Run`); die zwei anderen Stellen tragen ihn. |
| Stellen, die den Ausschlussstand mitführen, und ihr Regelstand-Zwilling | `git grep -n ExcludedColumns -- internal ':!*_test.go'` (Zeilen 3–4) und `git grep -n 'Transformations:' -- internal ':!*_test.go'` (Zeilen 5–6) | **Gefunden.** `ExcludedColumns` Parent 25, Diff 26 (eine Zeile mehr: der Doc-Kommentar von `TransformationRules`, der auf `ExcludedColumns` verweist); `Transformations:` Parent 0, Diff 4 (`activatedTableBindings` und Aktivierungs-Zweig, dazu die Felder `setTransformations:`/`removeTransformations:` im Literal von `Run`, die das Muster mitzählt). Träger, die `ExcludedColumns` lesen oder setzen: `mapper.go` (`kern-rename` trägt dort den Regelstand-Zwilling), `wiring.go` (zwei Stellen, beide mit Zwilling), `backfill/service.go` (liest den Ausschlussstand, der Regelstand kommt mit `slice-transformationen-backfill-pfad`), `sqlexec`/`tableactivation.go`/`columnexclusion.go` (Ableitung, Zwilling `ReadTransformationRules`/`TransformationRules`). Test-Fixtures: die neuen Fixtures tragen `Transformations` über die Bindung. **Nichtgefunden:** kein Träger im Erfassungs- oder Zustellpfad, der den Ausschlussstand ohne Regelstand kopiert. | `backfill/service.go` bleibt (Adresse `backfill-pfad`, dessen §2). |
| Beschreibung der Antragsarten im Doc-Kommentar von `applyAdministrationRequest` | `git grep -n 'Spalten-Antragsarten' -- internal/bootstrap/wiring.go` (Zeilen 7–8) | **Gefunden.** Parent 1, Diff 2; der Doc-Kommentar beschreibt die zwei neuen Zweige (Use Case, Nachtrag, Tabelle ohne Bindung, Idempotenz, Zeilen mit fehlenden Regelfeldern) und trägt den Aktivierungs-Zweig mit Regelstand. **Nichtgefunden:** keine zweite Beschreibung der Antragsarten in `internal/bootstrap`, die die zwei Arten als „nicht verarbeitet“ führt (`git grep -n 'nicht verarbeitet' -- internal/bootstrap` ohne Test-Treffer zu diesem Gegenstand). | nachgezogen. |
| Port-Übersichten in Doku | `git grep -n ColumnExclusionPort -- docs spec harness internal ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/in-progress' ':!*_test.go'` (Zeilen 9–10) | **Gefunden.** Parent 21, Diff 22 (die eine Zeile mehr: der Doc-Kommentar von `TransformationPort.SourceColumns`, der die Katalog-Lesart von `ColumnExclusionPort.ColumnExists` nennt; Treffer sind Code-Träger, drei ADRs, ein Register-State und `spec/architecture.md:274`, das den Spaltenausschluss beschreibt und richtig bleibt). Stand der Closure: Diff 23 (eine Zeile mehr: die Übergabe im Plan von `slice-transformationen-backfill-pfad` nennt `ColumnExclusionPort.ExcludedColumns`; gemessen mit `make suchlauf-nachmessen` am Arbeitsbaum der Closure). **Nichtgefunden:** keine Aufzählung der Outbound Ports (Port-Übersicht) in `docs/`, `spec/` oder `harness/`: `spec/architecture.md` nennt Ports im Text und in den Diagrammen, `ARC-004` führt keine Liste; die zwei Use Cases und „ein Outbound Port“ stehen dort bereits (`spec/architecture.md`, Antragsart-Tabelle, Absatz „Transformations-Antragsarten“). | keine Änderung. |
| Zweite Quelle der verarbeiteten Menge | `git grep -n processedAdministrationKinds -- internal` (Zeilen 11–12) und das Zählwort „fünf“ in `internal/bootstrap` (Zeilen 13–14) | **Gefunden.** Parent 5 Treffer (Konstante, Fehlertext, drei Test-Zeilen), Diff 3 (Funktion, Fehlertext, Doc-Kommentar); „fünf“ in `internal/bootstrap` Parent 12, Diff 9: die drei verschwundenen Treffer sind der Test und die Kommentare der Fünf-Arten-Menge; die neun übrigen zählen fremde Mengen (fünf Umgebungsvariablen, fünf Lese-Views, fünf Tabellen). **Nichtgefunden:** keine handgeführte Aufzählung der verarbeiteten Antragsarten mehr außerhalb des einen Testliterals, das den Fehlertext bindet (`TestApplyAdministrationRequestRejectsKindOutsideTheClosedSet`). | Menge abgeleitet, an die Fälle des `switch` gebunden. |
| Felder von `administrationDeps` | `git grep -n 'administrationDeps{' -- internal` (Zeilen 15–16) und `git grep -n 'transformations:' -- internal/bootstrap` (Zeilen 17–18) | **Gefunden.** Literale von `administrationDeps` Parent 21, Diff 21 (Stand nach Fixrunde 2: der Whitebox-Test `administration_callorder_internal_test.go` trägt ein Literal; bis Fixrunde 1: 20; `administration_internal_test.go` 12 → 11: zwei Literale der ersetzten Tests entfallen, `ruleFixture` fügt eines hinzu; jede andere Datei trägt dieselbe Zahl, das Literal des Login-Tests ist um drei Felder erweitert); Setzungen des Feldes `transformations:` Parent 0, Diff 11 (Stand nach Fixrunde 2 mit dem Literal des Whitebox-Tests `administration_callorder_internal_test.go`; bis Fixrunde 1: 10; sieben Fixtures bestehender Tests, das Literal des Login-Tests, `ruleFixture` und die Verdrahtung in `Run`). **Nichtgefunden:** kein Literal, dessen Pfad den Aktivierungs-Zweig erreicht, ohne das Feld (der erste Lauf von `make test` fand zwei Literale in `wiring_rest_internal_test.go` als Nil-Dereferenzierung, Exit 2, nachgezogen). | nachgezogen. |
| Fehlerklassen-Abbildung | Lesen von `classifyRunError` in `internal/bootstrap/wiring.go` | **Gefunden.** `classifyRunError` bildet Fehler des Capture-Pfads ab (Assembler, Replikation, Speicher); Antrags-Fehler laufen durch `processAdministrationRequests` in den `failed`-Vermerk und erreichen `classifyRunError` nicht. **Nichtgefunden:** keine Stelle, an der ein Antrags-Fehler den Prozess beendet. | keine neue Abbildung nötig. |
| Handbuch und Spec (Meldung) | `git grep -n -e set_transformation -e remove_transformation -- docs/user` (Zeilen 19–20); `git grep -n 'Domänen-Invarianten des Antrags-Konstruktors' -- spec` (Zeilen 21–22) | **Gefunden.** `docs/user`: Parent 0, Diff 0 (die Betreiber-Oberfläche ist unbeschrieben, bis `betriebsdoku` sie beschreibt); `spec`: Parent 1, Diff 1 — `SPEC-019`, Zeile „`rule_name` ist … Pflicht … (Domänen-Invarianten des Antrags-Konstruktors)“ nennt den Konstruktor als Ort einer Prüfung, die dieser Slice in den Use Case legt. **Nichtgefunden:** kein weiterer Träger der Aussage „der Konstruktor lehnt leeren Regelnamen ab“ (Register-Records und Pläne nennen den Zustand am Parent). Nach dem Nachzug der Closure: `spec` Diff 0 (die Klammer ist ersetzt; die Zeile zur Spalte nennt „Domänen-Invariante des Antrags-Konstruktors“ im Singular und bleibt richtig). | Spec-Zeile gemeldet und mit der Closure nachgezogen (Planner, Frist eingehalten); Handbuch: Aufschub mit Adresse `slice-transformationen-betriebsdoku`. |

| Ordnung der offenen Anträge (Fixrunde, Review F-1: bewegte Eigenschaft „die Ordnung der Queue-Abfrage“) | `git grep -n 'ORDER BY requested_at' -- internal ':!*_test.go'` (Zeilen 23–24) und `git grep -n 'ListPending. heute nur' -- docs ':!docs/reviews' ':!docs/plan/planning/in-progress'` (Zeilen 25–26) | **Gefunden.** `ORDER BY requested_at` Parent 3, Diff 4 (Parent: die Abfrage der offenen Anträge, die Ableitung des Ausschlussstands und die Aufnahme-Abfrage der Backfill-Runs mit `run_id`; die vierte Zeile ist `SelectAppliedTransformationRequests` aus diesem Slice; die Abfrage der offenen Anträge trägt jetzt `, administration_request_id`, das Muster trifft sie weiter). Der Text „`ListPending` heute nur nach `requested_at` ordnet“ steht in `ADR-0113` §Festlegung 1 (Parent 1, Diff 2: die zweite Fundstelle ist das Zitat in `ADR-0127` §Kontext, das die Aussage als überholt einordnet; beide `Accepted`, unberührbar). **Nichtgefunden:** in `spec/`, `harness/` und `docs/user` keine Beschreibung der Queue-Ordnung außer `SPEC-019` (Absatz „Ordnung der Verarbeitung“ mit Zweitschlüssel, seit `ADR-0127`); Träger im Code: Doc-Kommentare der Abfrage und des Ports nachgezogen. | die Aussage in `ADR-0113` ist mit `ADR-0127` eingeordnet, der Wortlaut von `SPEC-019` steht (§6). |
| Aufrufzeitpunkt statt Transaktionsbeginn (Fixrunde 2, [`ADR-0127`](../../adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md): bewegte Eigenschaft „`requested_at` ist der Aufrufzeitpunkt“) | Zählwort und Beschreibung `Transaktionszeitstempel\|Transaktionsbeginn\|Anlage-Reihenfolge\|Anlage-Zeitpunkt` über den ganzen Baum (Zeilen 27–28); Symbol `clock_timestamp` in `tools` und `internal` ohne Tests (Zeilen 29–30); `requested_at` in `docs/user` (Zeilen 31–32) | **Gefunden.** Beschreibung: Parent 15, Diff 14 — verschwunden sind die Kommentare von `queries.go` (`Transaktionszeitstempel`, `Anlage-Reihenfolge`), `administrationrequest.go` und `translate.go`; hinzugekommen drei Zeilen, die den Transaktionsbeginn als Spalten-Default oder als Wert der Mutation beschreiben (Kopfkommentar von `nacharbeit-administration.sql`, Mutationsnotizen der zwei neuen Tests); Stand der Closure: Diff 16 (zwei Zeilen mehr in der `state.md` des neuen Register-Eintrags `queue-ordnung-transaktionsbeginn-und-zufallskennung`, die den früheren Zustand als Beschreibung der Klasse nennen — gemessen mit `make suchlauf-nachmessen` am Arbeitsbaum der Closure; die Datei `evidence/slice-transformationen-antragsweg-usecase.md` des Eintrags trägt eine weitere Zeile und liegt außerhalb des Suchraums, weil das Werkzeug jede Datei mit dem Namen des Plans ausschließt). Die übrigen Treffer sind fremde Gegenstände: fünf Zeilen in `ADR-0127` (Architect-Zug, `Accepted`), `queries.go:278` (Spalten-Reihenfolge der Tabelle), `decode.go:36` (Beginn der Quelltransaktion), `sqlviews_test.go`, `service_test.go:528` (Uhr des Backfill-Use-Case), `table_schema.go`, `spec/pflichtenheft.md:889` (`SPEC-029`, die eigene Spalte `cdc.backfill_run.requested_at`). Symbol: Parent 4 (vier Zeilen in `run-integration-tests.sh`, Frist einer Wartezeit), Diff 14 (acht Zeilen in `nacharbeit-administration.sql`, zwei Kommentare in `queries.go`). **Nichtgefunden:** `docs/user` trägt kein `requested_at` (Parent 0, Diff 0); keine zweite Stelle, die den Funktionstext der sieben Funktionen trägt (`git grep -n gen_random_uuid -- tools internal examples cmd` trifft nur `nacharbeit-administration.sql`); Handbuch-Träger der Ordnung: keiner. | Handbuch: Aufschub mit Adresse `slice-transformationen-betriebsdoku` (dessen §2 trägt die Aussage, [`ADR-0127`](../../adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md) Folgepflicht 4). |

| Reichweite der Bindung von `clock_timestamp()` an die Tests (Fixrunde 3, Verifikation V-2/V-3: bewegte Eigenschaft „welche der sieben Funktionen der Test bindet“ und „wann die Mutation an einer Funktion rot färbt“) | Zählwort und Beschreibung `sechs Aufrufe\|nur für die später\|früher aufgerufenen Funktion` über den ganzen Baum (Zeilen 33–34) | **Gefunden.** Parent `d8c9e78f` 2 Zeilen, beide im Doc-Kommentar von `TestAdministrationRequestSameTransactionCallsKeepCallOrder` („sechs Aufrufe einer Transaktion“, „In der früher aufgerufenen Funktion bleibt die Reihenfolge richtig und der Test grün“); Diff 0 (der Kommentar nennt sieben Funktionen und die Regel „rot, sobald die Funktion einen Vorgänger hat“). **Nichtgefunden:** kein weiterer Träger im Baum außerhalb der Records (`docs/reviews/**`, `done/`), der die Zusage als „sechs Funktionen“ oder die Mutation als „im Fall früherer Aufruf grün“ beschreibt; der zweite Test (`…RemoveThenSetLeavesTheNewRuleLiveAndDerived`) trägt die Aussage für `remove_transformation` als erste Funktion seiner Folge zutreffend und bleibt. | Kommentar nachgezogen; die Plan-Zeile der Mutationstabelle ist berichtigt. |

```suchlauf
80eefead 5 -n 'TableBinding{' -- internal ':!*_test.go'
diff 5 -n 'TableBinding{' -- internal ':!*_test.go'
80eefead 25 -n ExcludedColumns -- internal ':!*_test.go'
diff 26 -n ExcludedColumns -- internal ':!*_test.go'
80eefead 0 -n 'Transformations:' -- internal ':!*_test.go'
diff 4 -n 'Transformations:' -- internal ':!*_test.go'
80eefead 1 -n 'Spalten-Antragsarten' -- internal/bootstrap/wiring.go
diff 2 -n 'Spalten-Antragsarten' -- internal/bootstrap/wiring.go
80eefead 21 -n ColumnExclusionPort -- docs spec harness internal :!docs/reviews :!docs/plan/planning/done :!docs/plan/planning/in-progress :!*_test.go
diff 23 -n ColumnExclusionPort -- docs spec harness internal :!docs/reviews :!docs/plan/planning/done :!docs/plan/planning/in-progress :!*_test.go
80eefead 5 -n processedAdministrationKinds -- internal
diff 3 -n processedAdministrationKinds -- internal
80eefead 12 -n fünf -- internal/bootstrap
diff 9 -n fünf -- internal/bootstrap
80eefead 21 -n 'administrationDeps{' -- internal
diff 21 -n 'administrationDeps{' -- internal
80eefead 0 -n 'transformations:' -- internal/bootstrap
diff 11 -n 'transformations:' -- internal/bootstrap
80eefead 0 -n -e set_transformation -e remove_transformation -- docs/user
diff 0 -n -e set_transformation -e remove_transformation -- docs/user
80eefead 1 -n 'Domänen-Invarianten des Antrags-Konstruktors' -- spec
diff 0 -n 'Domänen-Invarianten des Antrags-Konstruktors' -- spec
80eefead 3 -n 'ORDER BY requested_at' -- internal ':!*_test.go'
diff 4 -n 'ORDER BY requested_at' -- internal ':!*_test.go'
80eefead 1 -n 'ListPending. heute nur' -- docs ':!docs/reviews' ':!docs/plan/planning/in-progress'
diff 2 -n 'ListPending. heute nur' -- docs ':!docs/reviews' ':!docs/plan/planning/in-progress'
a69d3853 15 -n -E -e 'Transaktionszeitstempel|Transaktionsbeginn|Anlage-Reihenfolge|Anlage-Zeitpunkt' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!tools/schema/plan.yaml
diff 16 -n -E -e 'Transaktionszeitstempel|Transaktionsbeginn|Anlage-Reihenfolge|Anlage-Zeitpunkt' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!tools/schema/plan.yaml
a69d3853 4 -n clock_timestamp -- tools internal ':!*_test.go'
diff 14 -n clock_timestamp -- tools internal ':!*_test.go'
a69d3853 0 -n requested_at -- docs/user
diff 0 -n requested_at -- docs/user
d8c9e78f 2 -n -e 'sechs Aufrufe' -e 'nur für die später' -e 'früher aufgerufenen Funktion' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
diff 0 -n -e 'sechs Aufrufe' -e 'nur für die später' -e 'früher aufgerufenen Funktion' -- . :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
```

**Mutationen des Implementer-Laufs** (Zusage · mutierte Eingabe · gesehenes Rot; jede Mutation
einzeln am Arbeitsbaum, mit dem Edit-Werkzeug zurückgenommen, am Ende `git diff` gegen den
Commit ohne Quell-Abweichung; `make test` je Lauf, sofern nicht `make test-store` genannt — `make test-store`
fährt zuerst `internal/bootstrap` und bricht dort ab, die Zeilen nennen den ersten roten Test):

| Zusage | mutierte Eingabe | gesehenes Rot |
|---|---|---|
| K1: ein vergebener Regelname endet `failed` | Vergleich `rule.Name() == name` gegen `s.column` ersetzt (`CheckConflicts`) | `TestCheckConflictsBindsK1ToK3AndTheirOrder` („K1 Name vergeben“), `TestSetTransformationRejectsWithTheSpecTexts` (K1) |
| K2: eine Quellspalte trägt höchstens eine Spaltenregel | `rule.Column() == s.column` gegen `s.to` ersetzt | `TestCheckConflicts…` („K2“), `TestSetTransformationRejectsWithTheSpecTexts` (K2), `TestProcessAdministrationRequestsRuleViolationsFailWithSpecTexts` (`req-k2`) |
| K3: Zielname gleich dem Ziel einer anderen Regel | `rule.To() == s.to` gegen `s.column` ersetzt | `TestCheckConflicts…` („K3 Ziel gleicht Ziel einer anderen Regel“), Use-Case- und Verdrahtungs-Test (`req-k3a`) |
| K3: Zielname gleich der Quellspalte (auch ohne Katalogeintrag) | den Teil `s.to == s.column ||` entfernt | `TestCheckConflicts…` („… auch ohne Katalogeintrag“), `TestSetTransformationRejectsWithTheSpecTexts` („K3 (Ziel gleich Quelle) vor K4 bei fehlender Spalte“) |
| K3: Zielname gleich einer Spalte der Tabelle | den Teil `containsName(columns, s.to)` entfernt | `TestCheckConflicts…` („K3 Ziel gleicht Spalte der Tabelle“), Use-Case-Test (Spalte, ausgeschlossene Spalte), Verdrahtungs-Test (`req-k3b`) |
| K4: `column` existiert an der Quelle; Adresse ist die Spalte | Katalogprüfung auf `spec.Target()` statt `spec.Column()` | `TestSetTransformationAcceptsAConflictFreeRule`, `…BindsTheTableAndSourceOfTheCommand`, `…ComparesNamesCharacterExact`, Verdrahtungs-Test (erste Regel nicht `applied`) |
| K4: Adresse des Fehlertextes ist die Spalte | Adresse auf `spec.Target()` | `TestSetTransformationRejectsWithTheSpecTexts` (K4, `t.Errorf` je Fall) |
| die Adressen von K1, K2, K3 im Fehlertext | K1 → Spalte, K2 → Regelname, K3 → Spalte | `TestSetTransformationRejectsWithTheSpecTexts`: neun Fälle rot (K1, K2, K3 ×3, K4, drei Prüfreihenfolge-Fälle), `TestProcessAdministrationRequestsRuleViolationsFailWithSpecTexts` (`req-k1`) |
| `remove_transformation` gegen einen nicht geführten Namen endet `failed` | Namensvergleich `rule.Name() == command.RuleName` gegen `rule.Name() != ""` | `TestRemoveTransformationRejectsWithTheSpecTexts` (K4), Verdrahtungs-Test (`req-k4r`) |
| Regelstand je Tabelle der Anfrage | `state[schema.table]` gegen `state["public.other"]` (Set bzw. Remove) | Set: `TestSetTransformationAcceptsAConflictFreeRule`, `…BindsTheTableAndSourceOfTheCommand`, K1-Fall; Remove: `TestRemoveTransformationAcceptsAKeptRule`, `…RejectsWithTheSpecTexts` („Nachbartabelle“) |
| Regelstand der Quelle der Anfrage | `TransformationRules(ctx, "src-1")` statt `command.Source` (Remove) | `TestRemoveTransformationRejectsWithTheSpecTexts` („Nachbarquelle“) |
| Spaltenliste der Tabelle der Anfrage | `SourceColumns(ctx, schema, "other")` | drei Use-Case-Tests und fünf Verdrahtungs-Tests (u. a. `…TakeEffectLive`: „Spalte existiert nicht an der Quelle“) |
| Regelname ist geprüft (Alphabet, leer) | `CheckRuleName`-Fehler ignoriert (Set bzw. Remove); Grenze `{1,63}` → `{1,64}` | `TestSetTransformationRejectsWithTheSpecTexts` („Regelname leer“), `…ChecksTheFormBeforeReadingTheStore`, `TestRemoveTransformationRejectsWithTheSpecTexts` (zwei Namensfälle), `TestProcessAdministrationRequestsInvalidRuleRowsDoNotStallTheQueue`; Alphabet: `TestCheckRuleNameBindsTheAlphabet` |
| die Formzeilen laufen vor dem ersten Lesen des Ports | ein Lesen vor die Prüfung gesetzt | `TestSetTransformationChecksTheFormBeforeReadingTheStore` (Port-Aufrufe 1/0 statt 0/0) |
| `rule_spec` ist gültiges UTF-8 | `ValidString`-Prüfung ausgeschaltet | `TestParseTransformationSpecRejectsInSpecOrder` („ungültiges UTF-8“: `encoding/json` nimmt es an) |
| unbekannter Schlüssel in `rule_spec` endet `failed` | Prüfung `containsName(allowed, key)` ausgeschaltet | Parser-, Use-Case- und Verdrahtungs-Test (`req-key`) |
| unbekannter `kind` endet vor dem Schlüssel-Test | `if !known`-Zweig ausgeschaltet | Parser (`unbekannter kind` mit „Schlüssel column“ statt Regeltyp), Use Case, Verdrahtung, `TestFoldTransformationsFailsVisiblyOnUnparsableRow`, `TestReadTransformationRulesEmptyAndUnparsable` |
| Bezeichner-Form von `column`/`to`; Nicht-String liest sich als leerer Name | die Prüfung über `NewRenameColumn` ausgeschaltet | `TestParseTransformationSpecRejectsInSpecOrder` („column fehlt“), `…AcceptsValidForms` (64 Byte in 32 Zeichen), Use-Case-Test („Form vor K1“) |
| Zielname gleich Quellspalte ist keine Formverletzung, sondern K3 | den Ausschluss von `ErrTransformationTargetIsColumn` im Parser entfernt | `TestParseTransformationSpecAcceptsValidForms`, `TestTransformationSpecBuild`, Use-Case-Test („K3 Zielname gleicht der Quellspalte“), Verdrahtungs-Test (`req-k3c`) |
| Faltung: Reihenfolge der Zeilen | Zeilen rückwärts gelesen (`FoldTransformations`); Zeilen je Tabelle vorangestellt statt angehängt (`ReadTransformationRules`) | `TestFoldTransformationsFollowsTheOrderOfTheRows` („zwei Sets in Reihenfolge“), `TestReadTransformationRulesFoldsPerTable` |
| Faltung: Remove nimmt heraus | `dropByName` behält die Regel | `TestFoldTransformationsFollowsTheOrderOfTheRows` („Set, Remove“), `TestReadTransformationRulesFoldsPerTable` (`public.gone`) |
| Faltung: ein Set unter vorhandenem Namen ersetzt an der Stelle | Zweig `existing.name == rule.name` ausgeschaltet | `TestFoldTransformationsFollowsTheOrderOfTheRows` („Set unter vorhandenem Namen ersetzt“) |
| Faltung: eine nicht lesbare Zeile endet sichtbar | Fehler der Regelform übersprungen (`continue`); Ursache im Fehler nicht mehr umhüllt (`%s` statt `%w`, `ReadTransformationRules`) | `TestFoldTransformationsFailsVisiblyOnUnparsableRow`, `TestReadTransformationRulesEmptyAndUnparsable` |
| eine Tabelle ohne verbleibende Regel trägt keinen Eintrag; Zeilen je Tabelle getrennt | `len(rules) > 0` durch `true`; Schlüssel `schema.table` auf `table` verkürzt | `TestReadTransformationRulesFoldsPerTable` (`public.gone`; Tabellen-Schlüssel) |
| die Zeile mit leerem Regelnamen oder leerer Regelform ist ein Antrag, kein Lesefehler | Prüfung `ruleName == ""` in `NewAdministrationRequest` zurückgelegt | `TestNewAdministrationRequestCarriesEmptyRuleFields`, `TestReadPendingRequestsCarriesRuleRowsWithEmptyRuleFields`; `make test-store`: `TestAdministrationPathRunsUnderLeastPrivilegeLogins` („Zeile mit ungültigen Regelfeldern: Status `pending`“ — die Queue steht) |
| die geschlossene Menge der Antragsarten ist aufgezählt | `RemoveTransformation` aus `AdministrationRequestKinds` gestrichen | `TestAdministrationRequestKindsEnumeratesTheClosedSet` (sechs statt sieben), `TestApplyAdministrationRequestRejectsKindOutsideTheClosedSet` (Fehlertext) |
| jede Art der Menge trägt einen Zweig (Fenster „Funktion vorhanden, Wirkung fehlt“ geschlossen) | `case AdministrationRequestSetTransformation` durch eine fremde Art ersetzt | `TestApplyAdministrationRequestHandlesEveryKindOfTheClosedSet` und fünf Wirkungs-Tests (Fehlertext des `default`-Zweigs) |
| ein `set_transformation`-Antrag wirkt ohne Neustart | Nachtrag `deps.assembler.SetTransformation` gestrichen | `TestProcessAdministrationRequestsSetAndRemoveTransformationTakeEffectLive` und drei weitere (Bild bleibt roh) |
| `remove_transformation` nimmt die Regel heraus | Nachtrag `RemoveTransformation` gestrichen | `TestProcessAdministrationRequestsSetAndRemoveTransformationTakeEffectLive` (Bild trägt den Zielnamen weiter) |
| ein Antrag gegen eine Tabelle ohne Bindung legt keine an | Nachtrag durch `AddBinding` mit Regel ersetzt | `TestProcessAdministrationRequestsRuleAgainstTableWithoutBindingIsApplied` |
| Aktivierungs-Zweig trägt den Regelstand | Schlüssel `rules[qualified]` auf eine fremde Tabelle | `TestProcessAdministrationRequestsDisableEnableCycleRestoresTransformations` |
| Prozessstart trägt den Regelstand je Tabelle | Schlüssel auf eine fremde Tabelle; Schlüssel auf die Tabelle der Regel festgelegt | `TestActivatedTableBindingsCarriesTransformations` (beide Richtungen: „Transformations = []“ bzw. „Nachbartabelle“) |
| ein Lesefehler des Regelstandes endet den Start bzw. den Antrag sichtbar | Fehler in beiden Stellen verworfen (`rules = nil`) | `TestActivatedTableBindingsFehlerpfade` („TransformationRules-Fehler“), `TestProcessAdministrationRequestsMarksFailedWhenRuleStateReadFails` |
| die Ableitung liest `applied`-Zeilen in der Ordnung `requested_at`, `administration_request_id` | `administration_request_id` aus dem `ORDER BY` gestrichen | `make test-store`: `TestTableActivationTransformationRulesDeriveAppliedRuleRequests` (beide Paare gleicher Zeitstempel gehen verloren) — die erste Fassung des Tests (ein Paar, Zeilen nicht in aufsteigender Ordnung eingefügt) blieb **grün**; er ist auf zwei Paare in aufsteigender Einfüge-Ordnung umgestellt, sonst liest PostgreSQL bei einer nicht vorsortierten Eingabe die Gleichzeitigen in beliebiger Ordnung |
| K1 prüft nur gegen `applied`-Zeilen (Grundlage der Idempotenz) | Statusfilter `status = 'applied'` gegen `status <> 'failed'` | `make test-store`: `TestAdministrationPathRunsUnderLeastPrivilegeLogins` („set_transformation … `failed`: Regelname bereits vergeben“ — der eigene, noch `pending` stehende Antrag zählte) |
| die Ableitung deckt beide Antragsarten | `request_kind IN ('set_transformation')` | `make test-store`: Login-Test („remove_transformation K4 … `applied` statt `failed`“) |
| die Ableitung gilt je Quelle | `source_id = $1` gegen `source_id <> $1` | `make test-store`: Login-Test („set_transformation K1 … `applied` statt `failed`“) |
| die Spaltenliste liest die Tabelle der Anfrage | `table_name = $2` gegen `table_name = $2 \|\| 'x'` bzw. gegen `$1` | `make test-store`: Login-Test („Spalte existiert nicht an der Quelle“ bzw. Argumentzahl) |
| `SourceColumns` prüft Schema und Tabelle gegen das Alphabet | `validateIdentifier(schema)` gegen `validateIdentifier("public")` | `make test-store`: `TestTableActivationSourceColumnsReadsTheCatalog` („Schema außerhalb des Alphabets“) |
| Idempotenz der Wiederholung eines nachgetragenen, noch `pending` stehenden Antrags | — | **kein Rot am Code auf Unit-Ebene erreichbar:** `TestProcessAdministrationRequestsSetTransformationIsIdempotent` bindet das Verhalten mit einem Store-Fake, der nur vermerkte Anträge ableitet; jede Mutation an `applyAdministrationRequest`, die die Wiederholung bricht, müsste K1 gegen nicht vermerkte Zeilen prüfen — das leistet nur die Store-Abfrage, und deren Mutation (Zeile „K1 prüft nur gegen `applied`-Zeilen“) färbt den realen Login-Test rot. Ein Nachtrag, der Regeln anhängt statt zu ersetzen (`withTransformation`), ist am Bild unsichtbar (die erste treffende Regel entscheidet) — der Ersatz nach Namen ist im Assembler-Paket gebunden (`kern-rename`). |
| Regelliste ohne geteilten Speicher (Übergabe aus `kern-rename`) | — | **ohne Mutation:** die Nicht-Teilung folgt aus dem Aufbau der Faltung, kein Test kann sie am Ergebnis unterscheiden; benannte Grenze. |
| Fixrunde F-1: die Queue ordnet bei gleichem `requested_at` nach `administration_request_id` | `administration_request_id` aus dem `ORDER BY` von `SelectPendingAdministrationRequests` gestrichen | `make test-store`: `TestAdministrationRequestListPendingOrdersTiesByRequestID` („ListPending-Ordnung der Gleichzeitigen = [queue-order-b-set queue-order-a-remove queue-order-d-include queue-order-c-exclude]“, gedruckt); `TestAdministrationRequestSameTransactionRequestsAgreeLiveAndDerived` färbt sich bei dieser Mutation nur, wenn die zufällige Kennung des Remove die größere ist (die Reihenfolge der Funktionsaufrufe entspräche sonst der Kennungs-Ordnung) — er bindet den Pfad der realen Funktionen, die Mutation bindet der erste Test. |
| Fixrunde F-6: leeres Schema oder leerer Tabellenname endet beim Lesen als `ErrEmptyIdentifier` (Ist-Verhalten der Restklasse) | `schema == "" \|\| table == ""` aus der Prüfung in `NewAdministrationRequest` entfernt | `make test`: `TestReadPendingRequestsRejectsRowWithEmptySchemaOrTable` (beide Fälle) und `TestNewAdministrationRequestRejectsInvariantViolations` (Fall „leere Kennung“) |
| Fixrunde F-2, F-4, F-7: Kommentare (Mechanismus der Seeds, Bindung der Idempotenz, Grenze des Vermerk-Fehlers) | — | **ohne Mutation:** Kommentare ohne Verhalten; F-2 gemessen mit `git grep -n 'cfg.Tables'`, F-7 aus dem Quelltext hergeleitet und im Kommentar so benannt. |
| Fixrunde 2 ([`ADR-0127`](../../adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md)): Aufrufe einer Transaktion tragen `requested_at` je Aufruf, `remove_transformation` vor `set_transformation` derselben Regel wird in dieser Folge verarbeitet | `clock_timestamp()` im `INSERT` von `cdc.set_transformation` durch `now()` ersetzt (`tools/schema/nacharbeit-administration.sql`; das Set trägt den Transaktionsbeginn und sortiert vor das Remove mit `clock_timestamp()`) | `make test-store`, Exit 2 im vorgezogenen `internal/bootstrap`-Lauf: `TestAdministrationSameTransactionRemoveThenSetLeavesTheNewRuleLiveAndDerived` — „Durchlauf 0: Row Image der laufenden Bindung = {"id":"1","secret":"geheim"}, erwartet {"id":"1","second":"geheim"} (Remove vor Set) — Log: [WARN: administration: Antrag fehlgeschlagen]“ (das Set endete an K1 `failed`, das Remove nahm die Regel heraus) |
| Fixrunde 2: `ListPending` liefert die Aufruf-Reihenfolge (`exclude_column`/`include_column` derselben Spalte in einer Transaktion) | `clock_timestamp()` im `INSERT` von `cdc.include_column` durch `now()` ersetzt | `make test-store`, Exit 2 im Lauf von `postgresstorage`: `TestAdministrationRequestSameTransactionCallsKeepCallOrder` — „Durchlauf 0: ListPending-Ordnung = [169ac791-… 44a8669b-… …], wollen die Aufruf-Reihenfolge [44a8669b-… d98d8716-… f928b233-… 169ac791-… …]“ (der Include stand vor dem Exclude); `internal/bootstrap` blieb bei dieser Mutation grün, weil sie `include_column` nicht berührt |
| Fixrunde 2 (Stand vor Fixrunde 3, `remove_transformation` als erster Aufruf der Folge): dieselbe Zusage an `cdc.remove_transformation` | `clock_timestamp()` im `INSERT` von `cdc.remove_transformation` durch `now()` ersetzt | **kein Rot, `make test-store` Exit 0** am Stand vor Fixrunde 3: die Folge des Tests begann mit dem Remove — es trägt den Transaktionsbeginn, ihm ging kein Aufruf mit `clock_timestamp()` voraus, die Ordnung blieb richtig, weil `now()` nie hinter einem `clock_timestamp()` derselben Transaktion liegt. Die Farbe einer Mutation an `clock_timestamp()` hängt an der **Stellung** der Funktion in der Folge des Tests, nicht an ihrer Rolle im Paar: rot, sobald ein Aufruf mit `clock_timestamp()` vorausgeht (Verifikation V-3: `disable_table` und `exclude_column` als früher aufgerufene Funktion ihres Paares färbten rot). Die Lücke schloss Fixrunde 3 (Zeilen darunter). Hinweis an `ADR-0127` §Fitness Function: „`clock_timestamp()` aus einer Funktion entfernen“ ist für jede Funktion mit Vorgänger in der Folge wirksam; die Alt-Fassung mit `now()` in allen sieben Funktionen färbt über den Zufall der Kennung (Messung der ADR, hier nicht wiederholt). |
| Fixrunde 3 (V-2): jede der sieben Funktionen trägt `clock_timestamp()` in beiden Richtungen (Vorwärts-Folge `backfill`, `remove`, `set`, `exclude`, `include`, `disable`, `enable` und ihre Umkehrung in je einer Transaktion) | je einzeln `clock_timestamp()` im `INSERT` durch `now()` ersetzt, danach byte-genau zurückgenommen; enger Lauf gegen eine frische PostgreSQL 18 mit dem Rollout des mutierten Textes: `go test -count=1 -run TestAdministrationRequestSameTransactionCallsKeepCallOrder ./internal/adapters/driven/postgresstorage` | unmutiert Exit 0 (`PASS`, 3,79 s). Je Mutation Exit 1, Ausgabe „Durchlauf 0 (…): ListPending-Ordnung … wollen die Aufruf-Reihenfolge …“: `cdc.backfill_table` rot in der Umkehrung (dort an letzter Stelle, in der Vorwärts-Folge als erster Aufruf grün); `cdc.remove_transformation`, `cdc.set_transformation`, `cdc.exclude_column`, `cdc.include_column`, `cdc.disable_table` und `cdc.enable_table` rot in der Vorwärts-Folge. Nach der letzten Rücknahme kein Unterschied an `tools/schema/nacharbeit-administration.sql` (`git diff --stat`). |
| Fixrunde 2: Kommentare und Plan-Zeilen (Aufrufzeitpunkt, Betriebsdoku-Übergabe) | — | **ohne Mutation:** Kommentare und Plan ohne Verhalten; gemessen mit dem Suchlauf oben (Zeilen 27–32). |

**Läufe des Implementer-Laufs** (Exit-Code je Lauf ungefiltert in eine Log-Datei geschrieben und
gesondert gelesen, [`AGENTS.md`](../../../../AGENTS.md) §3.9; ein schwerer Docker-Lauf zugleich):

- `make test` (Race-Detector) Exit 0, 44 Pakete `ok`; `make test-store` Exit 0, gedruckt:
  `DB-Adapter-Coverage: 82.56% (gedeckt 885 von 1072 Statements; Profile gemergt: store,replication)`,
  `db-coverage: OK — DB-Adapter-Coverage 82.56% erfuellt Schwelle 80%` (der Login-Test
  `TestAdministrationPathRunsUnderLeastPrivilegeLogins` läuft real und ist in den Mutationen rot gesehen).
- `make a-check` Exit 0, `gesamt: 0 Befund(e)`; `make coverage-gate` Exit 0,
  `coverage-gate: OK — Coverage 84.80% erfüllt Schwelle 80%` (drei Läufe des Implementers: 84.80 %,
  84.80 %, 84.70 %; zwei Läufe des Reviewers: 84.70 %, 84.90 % — gemessen, die Zahl schwankt zwischen
  Läufen über eine Spanne von 0,2 Punkten, die Schwelle steht mit Abstand).
- `make gates` Exit 0 (alle sechs Ziele), gedruckt u. a. `baseline-verify: v6.9.0 OK — 54 Dateien`,
  `d-check: 1231 Datei(en) geprüft, 0 Befund(e)`, `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"`,
  `generated-sync: OK`, `gesamt: 0 Befund(e)`.
- `make suchlauf-nachmessen PLAN=<diese Datei>` Exit 0, `suchlauf-nachmessen: 22 Zeilen stimmen`;
  `make commit-traceability RANGE=origin/main..HEAD` Exit 0.
- §3.7-Probe, diff-skopiert gegen `80eefead`: ein Treffer, `administration_internal_test.go:30`
  (`slice-037`, vorbestehende Testfall-Provenienz, Subjekt: der Test), kein Treffer in einer
  geänderten Produktionsdatei.
- Kandidatenlauf Handbuch (`git diff --name-only 80eefead -- internal/bootstrap/ tools/schema/
  internal/adapters/driving/`): fünf Dateien, alle in `internal/bootstrap/` (Verdrahtung und Tests);
  keine neue Umgebungsvariable, keine SQL-Funktion, kein Endpunkt — keine neue Betreiber-Oberfläche,
  `docs/user/benutzerhandbuch.md` liegt bewusst nicht im Diff (die Oberfläche ist mit
  `antragsweg-schema` entstanden, Aufschub mit Adresse `slice-transformationen-betriebsdoku`).
- Nicht gelaufen, mit Begründung: `make test-integration` und `make test-replication` — dieser Slice
  ändert weder den Replikations- noch den Compose-Pfad; die Wirkung am laufenden Feed-Container
  belegt `slice-transformationen-e2e-wirkung`. `make image` — keine Änderung am Build-Kontext
  außerhalb von `internal/`.
- Docker-Volumes: `docker volume ls -q -f dangling=true | wc -l` vor und nach den Läufen jeweils 34.

**Läufe der Fixrunde** (Review-Findings F-1 bis F-11; Exit-Code je Lauf ungefiltert in eine
Log-Datei geschrieben und gesondert gelesen, ein schwerer Docker-Lauf zugleich; gemessen am Stand
nach dem Commit `cd4dae0e`, Arbeitsbaum sauber außer dieser Datei):

- `make test` (Race-Detector) Exit 0, 44 Pakete `ok`; `make test-store` Exit 0, gedruckt
  `DB-Adapter-Coverage: 82.56% (gedeckt 885 von 1072 Statements; Profile gemergt: store,replication)`,
  `db-coverage: OK — DB-Adapter-Coverage 82.56% erfuellt Schwelle 80%` (die neuen Store-Tests laufen
  real: die Mutation der Ordnung färbt `TestAdministrationRequestListPendingOrdersTiesByRequestID`
  rot, Exit 2).
- `make a-check` Exit 0, `gesamt: 0 Befund(e)`; `make coverage-gate` Exit 0, drei Läufe der Fixrunde:
  `84.90 %` (`make coverage-gate`), `84.80 %` (in `make gates`, zweimal) — Spanne der acht Läufe
  dieses Slice 84,70–84,90 %.
- `make suchlauf-nachmessen PLAN=<diese Datei>` Exit 0, `suchlauf-nachmessen: 26 Zeilen stimmen`.
- `make gates` Exit 0 (alle sechs Ziele), gedruckt u. a. `baseline-verify: v6.9.0 OK — 54 Dateien`,
  `d-check: 1232 Datei(en) geprüft, 0 Befund(e)`, `commit-traceability: OK — 5 Commit(s) in "HEAD~5..HEAD"`,
  `generated-sync: OK`, `gesamt: 0 Befund(e)`; `make commit-traceability RANGE=origin/main..HEAD`
  Exit 0, `OK — 10 Commit(s) in "origin/main..HEAD"`.
- `gofmt` (kein make-Ziel; Docker-only über das gepinnte `TOOLCHAIN_IMAGE`, `gofmt -l` schreibgeschützt,
  `gofmt -w` mit `--user` auf die drei Dateien des Findings F-3): die drei Dateien sind formkonform.
  `gofmt -l internal` nennt weiter fünf Dateien, die nicht zu diesem Diff gehören
  (`queries/queries.go` — die Umschreibung von `''` in einem Doc-Kommentar der Retention-Abfrage —,
  `mapper/transformation_test.go`, `receive/seam_test.go`, `port/outbound/log_test.go`,
  `retention/service_test.go`); ihre Form ist vor diesem Slice abweichend und bleibt unberührt.
- §3.7-Probe, diff-skopiert gegen den Stand vor der Fixrunde (`c7796d83`): ein Treffer,
  `administration_internal_test.go:30` (`slice-037`, vorbestehende Testfall-Provenienz, Subjekt: der
  Test), kein Treffer in einer geänderten Produktionsdatei.
- Handbuch: `git diff --name-only 80eefead -- internal/bootstrap/ tools/schema/ internal/adapters/driving/`
  bleibt bei fünf Dateien in `internal/bootstrap/`; keine neue Betreiber-Oberfläche, das Handbuch liegt
  nicht im Diff.
- Docker-Volumes: `docker volume ls -q -f dangling=true | wc -l` vor und nach den Läufen jeweils 34.

**Läufe der Fixrunde 2** ([`ADR-0127`](../../adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md);
Exit-Code je Lauf ungefiltert in eine Log-Datei geschrieben und gesondert gelesen, ein schwerer
Docker-Lauf zugleich; gemessen am Stand nach dem Commit der Funktionen und Tests, Arbeitsbaum sauber
außer dieser Datei):

- `make test` (Race-Detector) Exit 0, 44 Pakete `ok`; `make test-store` Exit 0, gedruckt
  `DB-Adapter-Coverage: 82.56% (gedeckt 885 von 1072 Statements; Profile gemergt: store,replication)`,
  `db-coverage: OK — DB-Adapter-Coverage 82.56% erfuellt Schwelle 80%` — mit den sieben Funktionen im
  neuen Text (der Rollout des Laufs druckt sieben `CREATE FUNCTION`) und den zwei neuen Tests
  (`TestAdministrationSameTransactionRemoveThenSetLeavesTheNewRuleLiveAndDerived` im vorgezogenen
  `internal/bootstrap`-Lauf, `TestAdministrationRequestSameTransactionCallsKeepCallOrder` im Lauf von
  `postgresstorage`; beide in den Mutationen oben rot gesehen).
- `make a-check` Exit 0, `gesamt: 0 Befund(e)`; `make coverage-gate` Exit 0,
  `coverage-gate: OK — Coverage 84.90% erfüllt Schwelle 80%` (in `make gates` 84.70 % — die Zahl
  schwankt zwischen Läufen wie bisher).
- `tools/harness/run-schema-rollout-guard-test.sh` Exit 0, alle sechs Läufe, gedruckt
  `Lauf 5 OK — Tag v0.2.0: Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, ohne Vorlauf), Exit 0
  (Arbeitsbaum, zweiter Lauf)` (Alt-Tag-Lauf, `CREATE OR REPLACE` mit unveränderter Signatur: die
  Rollen-Rechte, `EXECUTE` allein für `cdc_admin`, die zwei Spalten und der Aufruf beider
  Transformations-Funktionen unter `cdc_admin` stehen unverändert; damit ist die als hergeleitet
  geführte Einordnung von `ADR-0127` §Konsequenzen für diesen Lauf erprobt) und
  `OK — alle Belege real erbracht (… Alt-Tag v0.2.0, Negativ-Abbruch: … make-Exit 2/2 mit
  d-migrate-Exit 8)`. Zwei aufeinanderfolgende `make schema-rollout` gegen eine Wegwerf-Datenbank sind
  die Läufe 1 (frisch) und 2 (Idempotenz) dieses Skripts, je Exit 0; `git status` nach dem Lauf: kein
  Unterschied an `tools/schema/plan.yaml` und `tools/schema/down.sql` — der Wrapper
  `tools/schema/rollout-restore.sh` stellt sie wieder her, und der neue Funktionstext liegt außerhalb
  des d-migrate-Modells (`schema.yaml` und Spalten-Default unverändert), der committete Stand bleibt
  also gültig ([`harness/targets/schema-rollout.md`](../../../../harness/targets/schema-rollout.md)
  §Erzeugnisse in Test-, Bench- und Beispiel-Läufen).
- `make suchlauf-nachmessen PLAN=<diese Datei>` Exit 0, `suchlauf-nachmessen: 32 Zeilen stimmen`
  (die ersten Läufe nach dem Nachzug der Zeilen 27–32 wichen an vier Soll-Werten ab — Literale und
  Setzungen des neuen Whitebox-Tests, das Zitat in `ADR-0127`, die Mutationsnotiz desselben Tests —,
  sie sind nachgezogen).
- `make gates` Exit 0 (alle sechs Ziele), gedruckt u. a. `baseline-verify: v6.9.0 OK — 54 Dateien`,
  `d-check: 1233 Datei(en) geprüft, 0 Befund(e)`, `generated-sync: OK`, `gesamt: 0 Befund(e)`; der
  erste Lauf endete Exit 2 an einer nackten Kennung in dieser Datei (`id-unlinked`), nachgezogen.
- `gofmt -l` (Docker-Only, gepinntes `TOOLCHAIN_IMAGE`) über die berührten Go-Dateien ohne
  `queries/queries.go`: keine Ausgabe (`queries.go` ist seit dem Retention-Kommentar formabweichend,
  siehe Fixrunde 1).
- §3.7-Probe, diff-skopiert gegen `a69d3853`: kein Treffer in einer geänderten `*.go`- oder
  `tools/schema/*.sql`-Datei.
- Handbuch: `git diff --name-only 80eefead -- internal/bootstrap/ tools/schema/ internal/adapters/driving/`
  nennt jetzt zusätzlich `tools/schema/nacharbeit-administration.sql` (der Funktionstext ändert den
  Zeitstempel, keine Signatur, keine neue Funktion, keine neue Umgebungsvariable, keinen Endpunkt) —
  keine neue Betreiber-Oberfläche; das Handbuch liegt nicht im Diff, der Aufschub mit Adresse
  `slice-transformationen-betriebsdoku` trägt die Aussage zur Aufruf-Reihenfolge.
- Nicht gelaufen, mit Begründung: `make test-integration` und `make test-replication` — weder
  Replikations- noch Compose-Pfad geändert; `make image` — kein Build-Kontext außerhalb von `internal/`
  und `tools/schema/` geändert (die Funktionen rollt `make schema-rollout` aus, den das Guard-Skript
  und `make test-store` real fahren).
- Docker-Volumes: `docker volume ls -q -f dangling=true | wc -l` vor den Läufen 34, nach den Läufen 34.

**Läufe der Fixrunde 3** (Verifikation V-2, V-3, V-5; Exit-Code je Lauf ungefiltert in eine
Log-Datei geschrieben und gesondert gelesen, ein schwerer Docker-Lauf zugleich; gemessen am
Arbeitsbaum über dem Commit `d8c9e78f`):

- Enger Lauf des geänderten Tests gegen eine frische PostgreSQL 18 mit dem Rollout des Arbeitsbaums
  (`go test -count=1 -v -run TestAdministrationRequestSameTransactionCallsKeepCallOrder`): Exit 0,
  `--- PASS: TestAdministrationRequestSameTransactionCallsKeepCallOrder (3.79s)`; dieselbe Zeile mit je
  einer Mutation (Zeile „Fixrunde 3 (V-2)“ der Mutationstabelle): siebenmal Exit 1.
- `make test` (Race-Detector) Exit 0, 44 Pakete `ok`; `make test-store` Exit 0, gedruckt
  `DB-Adapter-Coverage: 82.56% (gedeckt 885 von 1072 Statements; Profile gemergt: store,replication)`,
  `db-coverage: OK — DB-Adapter-Coverage 82.56% erfuellt Schwelle 80%`,
  `ok  github.com/pt9912/pg-change-feed/internal/adapters/driven/postgresstorage 13.386s`.
- `make a-check` Exit 0, `gesamt: 0 Befund(e)`; `make coverage-gate` Exit 0,
  `coverage-gate: OK — Coverage 84.80% erfüllt Schwelle 80%` (in `make gates` 84.70 %).
- `make suchlauf-nachmessen PLAN=<diese Datei>` Exit 0, `suchlauf-nachmessen: 32 Zeilen stimmen` vor dem
  Nachzug der Zeilen 33–34, danach 34 Zeilen (Lauf am Ende der Fixrunde).
- `make gates` Exit 0 (alle sechs Ziele), gedruckt u. a. `baseline-verify: v6.9.0 OK — 54 Dateien`,
  `d-check: 1234 Datei(en) geprüft, 0 Befund(e)`, `commit-traceability: OK — 5 Commit(s) in
  "HEAD~5..HEAD"`, `generated-sync: OK`, `gesamt: 0 Befund(e)`.
- `gofmt -l` (Docker-Only, gepinntes `TOOLCHAIN_IMAGE`, `--network none`) über
  `administrationrequest_order_test.go` und `wiring.go`: keine Ausgabe.
- §3.7-Probe, diff-skopiert gegen `d8c9e78f`: kein Treffer in einer geänderten `*.go`- oder
  `tools/schema/*.sql`-Datei.
- Handbuch: die Fixrunde ändert weder `internal/bootstrap/` (nur einen Doc-Kommentar) noch
  `tools/schema/` (der Funktionstext ist nach den Mutationen byte-gleich zurückgenommen) noch einen
  Adapter unter `internal/adapters/driving/`; das Handbuch liegt nicht im Diff.
- Docker-Volumes: `docker volume ls -q -f dangling=true | wc -l` vor den Läufen 34, nach den Läufen 34;
  keine Container und keine Netze mit dem Präfix `cdc-` nach den Läufen.

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `antragsweg-schema` in `done/` liegt
(die Spalten und Antragsarten bestehen im Store) und kein anderer Slice in
`in-progress/` liegt (WIP-Limit 1). Die Kopplung K3 der Welle
[welle-backfill-bestand](../done/welle-backfill-bestand.md) §5 ist mit
`antragsweg-schema` erfüllt; `applyAdministrationRequest` trägt zu diesem
Zeitpunkt bereits den Backfill-Zweig, dieser Slice erweitert ihn additiv.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Use Cases,
  Store-Adapter und Verdrahtung nicht in einem Review tragen — der abtrennbare
  Teil ist der zweite Liefer-Punkt (Store-Adapter für Regelstand und
  Spaltenliste, Tier `make test-store`) als eigener Slice mit Start nach der
  Use-Case-Hälfte.
- `in-progress` → `open` (blockiert): falls der Port-Schnitt (ein Port für
  Regelstand und Spaltenliste, wie
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  ihn zählt, gegen zwei Fähigkeiten nach
  [`ADR-0034`](../../adr/0034-ports-nach-faehigkeiten.md)) eine
  Architect-Entscheidung braucht, oder falls die Spaltenliste an der Quelle
  nicht ohne zusätzliche Rechte der Login-Identität von `CDC_ADMIN_DSN` lesbar
  ist (Rollen-Frage nach
  [`ADR-0047`](../../adr/0047-rollenspezifische-dsn-verdrahtung.md)).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test` (Race-Detector), `make
test-store` und `make coverage-gate` real grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Der abgeleitete Stand und die laufende Bindung laufen auseinander**, wenn
  der Vermerk `applied` nach dem Nachtrag scheitert
  ([`ADR-0065`](../../adr/0065-spaltenausschluss-dauerhafter-traeger.md)-Muster;
  `processAdministrationRequests` meldet nur „Erfolg nicht vermerkt“).
  *Erwartet, zu belegen durch:* der Idempotenz-Test aus DoD Punkt 3
  (Wiederholung desselben Antrags ist folgenlos). **Grenze (hergeleitet aus
  dem Quelltext, nicht erprobt):** der Test bindet nur die Wiederholung
  **desselben** Antrags. Verarbeitet ein Durchlauf dazwischen einen Antrag, der
  zum noch nicht vermerkten in K2 oder K3 steht (dieselbe Spalte, dasselbe
  Ziel), prüft er gegen einen Regelstand ohne den ersten und wird `applied`;
  die Wiederholung des ersten endet danach `failed`, die laufende Bindung trägt
  bis zum Neustart beide Regeln, und zwei Regeln mit gleichem Zielnamen in
  einer Bindung beenden den Erfassungspfad in der Fehlerklasse `schema`. Der
  Fall setzt einen fehlgeschlagenen Vermerk **und** einen kollidierenden
  Folgeantrag voraus; der Kommentar am Set-Zweig in
  `applyAdministrationRequest` nennt die Grenze. **Ausgang:** *weiter offen* — die
  Idempotenz der Wiederholung ist belegt (Store-Abfrage mit Statusfilter, gebunden im Login-Test;
  Mutation rot), die Grenze bleibt hergeleitet, nicht erprobt (Review F-7, LOW). Adresse: das
  Register, `BEO-PGC/wertabhaengiger-zweiter-waechter-ohne-spec-zeile` (`state.md`, Absatz „Stand
  nach `slice-transformationen-antragsweg-usecase`“): die Konstellation ist der Weg (1) zu dem
  zweiten Wächter, dessen Auflösung dort ansteht.
- **Verarbeitungs-Ordnung und Ableitungs-Ordnung sind dieselbe, `requested_at`
  ist der Aufrufzeitpunkt** ([`ADR-0127`](../../adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md)).
  `SelectPendingAdministrationRequests` ordnet `requested_at`, bei gleichem
  Zeitstempel nach `administration_request_id`, wie die zwei Abfragen der
  Ableitung; die sieben SQL-Funktionen schreiben `requested_at` je Aufruf mit
  `clock_timestamp()`, Aufrufe einer Transaktion tragen verschiedene
  Zeitstempel in der Aufruf-Reihenfolge (auch `exclude_column`/`include_column`
  und `disable_table`/`enable_table`). Bindung: die Store-Tests in
  `administrationrequest_order_test.go`, der Whitebox-Test in
  `administration_callorder_internal_test.go`. Verbleibende Grenzen
  (Festlegung 3 der ADR): die Ordnung ist der Zeitpunkt des Aufrufs, nicht der
  des `COMMIT` — Anträge auf dieselbe Regel oder Spalte aus zeitlich
  überlappenden Transaktionen mehrerer Sitzungen können live und abgeleitet
  bis zum nächsten Prozessstart abweichen (hergeleitet, nicht erprobt);
  Rückwärtssprung der Serveruhr; Zeilen vor der Änderung tragen den
  Transaktionsbeginn. Die Aussage von
  [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md)
  Festlegung 1 („`ListPending` heute nur nach `requested_at` ordnet“) ist nach
  [`ADR-0127`](../../adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md)
  §Kontext überholt, ihr Zweitschlüssel bleibt richtig. **Ausgang:**
  *eingetreten und im Slice aufgelöst* — Review F-1 (MEDIUM) fand die Lücke, die Fixrunde 1
  glich die Ordnung bei gleichem Zeitstempel an, [`ADR-0127`](../../adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md)
  legte die Aufruf-Reihenfolge fest (Fixrunde 2), die Fixrunde 3 band alle sieben Funktionen in beiden
  Richtungen. Die Klasse steht im Register (`BEO-PGC/queue-ordnung-transaktionsbeginn-und-zufallskennung`,
  1×). Die drei Grenzen der ADR (überlappende Transaktionen mehrerer Sitzungen, Rückwärtssprung der
  Serveruhr, Zeilen vor der Änderung) sind ein akzeptiertes Negativ **in der ADR** (Festlegung 3, mit
  Re-Evaluierungs-Trigger) und tragen **keinen eigenen Register-Eintrag**: ein Eintrag verlangt ein
  Auftreten in `evidence/`, und keine der drei Grenzen ist aufgetreten (hergeleitet, nicht
  erprobt); die Betreiberregel trägt `slice-transformationen-betriebsdoku` §2. Die Aussage von
  [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1
  („`ListPending` heute nur nach `requested_at` ordnet“) ist überholt (Kenntnis, `Accepted` und
  unberührbar; `ADR-0127` §Kontext trägt die Einordnung).
- **Ein Antrag mit leerem Schema oder leerem Tabellennamen stallt die Queue
  weiter** (Klasse `BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue`, zur
  Hälfte gelöst). Die Regelfelder werden verarbeitet statt beim Lesen
  abgelehnt; der Konstruktor lehnt beim Lesen weiter ab: die beiden
  Spalten-Antragsarten mit leerer Spalte und **jede** Antragsart mit leerem
  Schema oder leerem Tabellennamen (`ErrEmptyIdentifier`). Erprobt am
  Ist-Verhalten: `TestReadPendingRequestsRejectsRowWithEmptySchemaOrTable`
  (die Lesung endet am Fehler, die Zeile dahinter bleibt ungelesen; rot bei
  entfernter Prüfung im Konstruktor). Ob die SQL-Funktionen einen solchen
  Antrag schreiben (`cdc.enable_table(…, '', 't')`), ist aus
  `tools/schema/nacharbeit-administration.sql` und `schema.yaml` hergeleitet
  (`NOT NULL` ohne weitere Prüfung), nicht erprobt. Kein Code-Fix in diesem
  Slice: `SPEC-019` nennt für Schema und Tabelle keinen Fehlertext, ein `failed`-
  Vermerk ohne Adresse nähme eine Entscheidung vorweg. **Adresse:** der Planner
  zieht in der Closure dieses Slice nach — Register-Zustand des Eintrags und
  Meldung der Zeile `SPEC-019` „Domänen-Invarianten des Antrags-Konstruktors“.
  **Ausgang:** *weiter offen → Register* (`BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue`,
  2×, `state.md` nachgezogen: Regelfelder gelöst, Rest offen); die Spec-Zeile zu den Regelfeldern
  ist ersetzt, die Zeile zur Spalte bleibt richtig; die Entscheidung über den Rest ist ein Vorschlag
  an den Architect-Zug der Closure von `welle-transformationen` (im `state.md`).
- **Kenntnis aus dem Review (INFO, keine Aktion dieses Slice):** eine
  `applied`-Zeile, die die Ableitung nicht mehr in eine Regel führt, endet als
  Fehler und hält Prozessstart und jeden Regel-Antrag **jeder** Tabelle der
  Quelle an (bewusst: der Stand wird nie um eine Zeile verkürzt; die
  Fehlerklasse ist `internal`); und `TransformationRules` liest je Aufruf alle
  `applied`-Zeilen der Quelle — nach dem Plan von
  `slice-transformationen-backfill-pfad` („je Block neu“) je Block eines
  Backfill-Runs eine Lesung über die ganze Quelle. **Adresse:**
  `slice-transformationen-backfill-pfad` (Kosten der Lesung je Block). **Ausgang:** die Kosten
  (F-11) *weiter offen* — Adresse: der Übergabe-Block „Übergabe aus
  `slice-transformationen-antragsweg-usecase`“ im Plan von `backfill-pfad` (in der Closure angelegt);
  die Nicht-Lesbarkeit einer `applied`-Zeile (F-8) *entfallen als Risiko dieses Slice* — bewusste
  Eigenschaft der Ableitung, ihre Fehlerklasse steht im Doc-Kommentar von `activatedTableBindings`;
  die Rückfall-Grenze auf einen Binärstand ohne `map_value` ist in den Übergabe-Blöcken von
  `slice-transformationen-map-value` und im Handbuch-Punkt von `slice-transformationen-betriebsdoku`
  (Dauerhaftigkeit) getragen.
- **K3 prüft gegen die Spaltenliste zum Antragszeitpunkt**; eine spätere
  Spalten-Erweiterung kann den Zielnamen kollidieren lassen. Das ist der Fall,
  den die Prüfung im Assembler (`kern-rename`) fängt und den `e2e-abhilfe` real
  belegt; hier wird er nicht verhindert. *Erwartet, zu belegen durch:*
  Kommentar am K3-Zweig, der die Grenze nennt (`wiring.go`, gelesen), und der Verweis in §7.
  **Ausgang:** *weiter offen* — der Kommentar trägt die Grenze; die reale Kollision mit einer später
  hinzugefügten Spalte belegt `slice-transformationen-e2e-abhilfe` (dessen §2: `ADD COLUMN` mit dem
  Zielnamen der Regel und die Gegenprobe).
- **Fenster: Funktion vorhanden, Wirkung fehlt** (Übergabe aus
  `antragsweg-schema`). Seit dessen Closure schreiben `cdc.set_transformation`
  und `cdc.remove_transformation` einen Antrag, den `applyAdministrationRequest`
  im `default`-Zweig als `failed` mit dem Fehlertext der verarbeiteten
  Antragsarten vermerkt: sichtbar, keine Regel, keine Wirkung. Das Fenster endet
  mit der Closure dieses Slice. *Erwartet, zu belegen durch:* der Test des
  `default`-Zweigs färbt sich rot, sobald der Zweig für `set_transformation`
  entfällt (die zweite Quelle der verarbeiteten Menge, Übergabe oben).
  **Ausgang:** *entfallen* — die zwei Tests des `default`-Zweigs für `set_transformation` sind durch
  die Wirkungs-Tests ersetzt; `TestApplyAdministrationRequestHandlesEveryKindOfTheClosedSet` und
  fünf Wirkungs-Tests färben sich rot, sobald der Zweig für eine Art entfällt (Mutationstabelle in §3, Zeile „jede Art der Menge trägt einen
  Zweig“; Reviewer M24 für `remove_transformation`).
- **Fenster zwischen Regel-Setzbarkeit und Backfill-Bindung.** Ab diesem Slice
  sind Regeln setzbar; ein Backfill-Run, der vor `backfill-pfad` läuft,
  lieferte die Rohform (Welle §5). Der Slice ändert das Backfill-Verhalten
  nicht. *Erwartet, zu belegen durch:* die Reihenfolge der Welle
  (`backfill-pfad` folgt diesem Slice unmittelbar) und die Benennung im
  Bericht. **Ausgang:** *weiter offen* — das Fenster besteht ab dieser Closure (Regeln sind setzbar, ein Run
  liefert die Rohform, gemessen am Stand: `slice-transformationen-backfill-pfad` liegt in `open/`) und
  endet mit der Closure von `slice-transformationen-backfill-pfad`. Adresse: der Risiko-Punkt
  „Fenster zwischen Regel-Setzbarkeit und Backfill-Bindung“ in §6 des Plans von `backfill-pfad`
  (in der Closure angelegt, Ausgang dort „entfallen mit der Closure dieses Slice“).
- **Zwischenzustand der Regeltyp-Menge.** Der Use Case prüft gegen die
  Regeltyp-Menge der Domäne; bis `map-value` (Reihenfolge der Welle: nach
  `backfill-pfad`) besteht sie aus `rename_column`, ein `map_value`-Antrag endet
  in dieser Spanne `failed` mit `unbekannter Regeltyp`, während
  [`SPEC-019`](../../../../spec/pflichtenheft.md) ihn als Regeltyp führt.
  Festlegung: der Zwischenstand ist zulässig, weil der Antragsweg bis
  `map-value` weder veröffentlicht noch durch einen E2E-Beleg als Zusage
  gelesen wird (`e2e-wirkung` folgt `map-value`); die Abweichung endet mit
  `map-value`, dessen §2 den `map_value`-Antrag über denselben Use Case prüft.
  *Erwartet, zu belegen durch:* die Reihenfolge der Welle und der Use-Case-Test
  aus `map-value` §2. **Ausgang:** *weiter offen* — der Zwischenstand besteht ab dieser Closure (ein `map_value`-Antrag
  endet `failed` mit `unbekannter Regeltyp`) und endet mit der Closure von
  `slice-transformationen-map-value`. Adresse: der zweite DoD-Punkt von `map-value` (Use-Case-Test des
  `map_value`-Antrags) und der Übergabe-Block „Übergabe aus `slice-transformationen-antragsweg-usecase`“
  in dessen §3 (in der Closure angelegt).
- **Der Port-Schnitt folgt der ADR-Zählung** („ein neuer Outbound Port“,
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  §Schärft): ein Port mit zwei Lese-Methoden, wie `ColumnExclusionPort` zwei
  Objektklassen in einen Port schneidet; der Review nennt den Diff
  ohne Abweichung (kein Auftreten von
  `BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab`) und nennt den Schnitt
  mit [`ADR-0034`](../../adr/0034-ports-nach-faehigkeiten.md) vereinbar (Review
  `review-slice-transformationen-antragsweg-usecase`, Negativbefund
  Port-Schnitt und Finding F-5): beide Methoden lesen, die Konsistenzgrenze der
  ADR trifft keine gemeinsame Transaktion. **Ausgang:** *entfallen* —
  der Review trägt die Aussage (Negativbefund Port-Schnitt; F-5 berichtigte die Plan-Aussage, die
  fälschlich „weicht ab“ sagte); kein Auftreten von
  `BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` (`state.md` nennt es, Zähler bleibt 2×).
- **Store-Tests teilen Zustand** (`BEO-PGC/test-isolation-geteilter-zustand`,
  offen, 1×): unskopierte `DELETE`/`DROP SCHEMA CASCADE` in einem Store-Test
  könnten die neuen Zeilen zerstören. *Erwartet, zu belegen durch:* skopierte
  Bereinigung je Test. **Ausgang:** *entfallen* — die Store-Tests dieses Slice räumen skopiert auf
  (Kennungs-Vorsilbe, `source_id`, `table_name`; kein unskopiertes `DELETE`/`DROP SCHEMA` im Diff,
  Review-Negativbefund). Kein Auftreten: `BEO-PGC/test-isolation-geteilter-zustand` zählt 2×
  (gemessen `ls evidence | wc -l`; §8 nannte 1× am Planungsstand), Zähler unverändert.
- **Kosten und Rollen der Spaltenlisten-Abfrage**: sie läuft je
  `set_transformation`-Antrag über die Administrations-Verbindung
  (`cdc_admin`). *Erwartet, zu belegen durch:* `make test-store` unter der
  realen Rolle (gleiche Katalog-Lesart wie `ColumnExists`, gelesen am Bestand).
  **Ausgang:** *entfallen* — `TestAdministrationPathRunsUnderLeastPrivilegeLogins` liest die
  Spaltenliste unter dem realen `cdc_admin`-Login (`[id secret]`, `make test-store` Exit 0; die
  Mutationen der Tabellen- und Argument-Eingabe sind rot). Die Kosten sind eine Katalogabfrage je
  `set_transformation`-Antrag, ohne Messwert; die Kosten der Regelstand-Lesung je Block trägt
  `backfill-pfad`.

## 7. Closure-Notiz

- **Was hat funktioniert:** (1) Die Leser trugen, wo die Mutationen des Implementers nicht trugen: der Plan führt 49 Tabellenzeilen unter „Mutationen“ (gezählt an §3 am Stand der Closure), 45 mit Mutation und 4 ohne (benannt), und der Review fand F-1 (MEDIUM) trotzdem — durch Lesen zweier Abfragen gegeneinander, nicht durch eine Mutation; die 26 Mutationen des Reviewers (M1 bis M26, übernommen aus dem Review-Report) waren rot bis auf M4, das in `internal/bootstrap` dort äquivalent gedeckt ist. Der Verifier mutierte an Kopien (G1 bis G7, S1 bis S8, übernommen aus dem Verifikations-Report) und fand V-2 durch zwei grüne Mutationen (S4, S7), die er mit einem Zusatztest (R0 bis R2) bestätigte; die Tabelle des Implementers führte für beide Funktionen keine Mutation mit Rot. (2) F-1 wurde nicht als Grenze in Spec und Handbuch festgeschrieben, sondern an der Instanz entschieden: [`ADR-0127`](../../adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md) trägt eine Messtabelle an zwei PostgreSQL-Hauptversionen (494 bzw. 513 von 1000 Transaktionen mit falscher Folge, übernommen aus der ADR), die Fixrunde 2 setzt `clock_timestamp()` in sieben Funktionen, die Fixrunde 3 bindet alle sieben in beiden Richtungen (sieben Einzel-Mutationen je Exit 1, unmutiert Exit 0, gemessen im Plan §3). (3) Das Suchlauf-Feld trug: `make suchlauf-nachmessen` stimmt am Arbeitsbaum der Closure mit 34 Zeilen (gemessen, Exit 0); die Zeilen der Fixrunde 2 wichen zuerst an vier Soll-Werten ab und wurden vom Werkzeug vor dem Handoff gefangen (kein Reviewer-Fund). (4) Der Alt-Tag-Lauf des Guard-Skripts gegen `v0.2.0` (sechs Läufe, Exit 0, Verifikations-Report §1) erprobte die zuvor als hergeleitet geführte Einordnung von `ADR-0127` §Konsequenzen (`CREATE OR REPLACE` mit gleicher Signatur). (5) Der Review-Negativbefund trug die Port-Schnitt-Aussage: ein Port mit zwei Lese-Methoden entspricht „ein neuer Outbound Port“ ([`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)) und ist mit [`ADR-0034`](../../adr/0034-ports-nach-faehigkeiten.md) vereinbar.
- **Was ging anders als geplant:** (1) **F-1 brauchte zwei Züge.** Die Fixrunde 1 ordnet die Queue bei gleichem Zeitstempel nach der Kennung wie die Ableitung; erst der Architect-Zug legte die Aufruf-Reihenfolge fest, weil die Kennung zufällig ordnet. Der Slice änderte damit den Funktionstext (`tools/schema/nacharbeit-administration.sql`), den §1 ausschloss („Schema, Funktionen, Grants, Guard — `antragsweg-schema`“); der Architect-Zug ist die Grundlage, §3 nennt die Zeile, der Verifier las keine unbenannte Abweichung. (2) **Umfang:** 16 Commits, 34 Dateien, +4774/−559 (gemessen mit `git diff --shortstat 2a47cd9a..2c22334f`, beide Seiten der zwei Lifecycle-Moves eingeschlossen); drei Fixrunden, ein Architect-Zug, eine [`SPEC-019`](../../../../spec/pflichtenheft.md)-Anpassung. (3) **Die Plan-Aussagen trugen ihren Beleg-Anker nicht** an drei Stellen: „weicht von der ADR-Zählung ab“ (F-5), eine Mutationsangabe ohne setzbare Mutation (F-4) und die Abgrenzung „nur für die später aufgerufene Funktion“ (V-3) — die Instanz-B-Klasse ([`AGENTS.md`](../../../../AGENTS.md) §3.12), von Reviewer und Verifier gefunden. (4) **Zähler des Plans standen hinter dem Register:** §8 nannte `BEO-PGC/test-isolation-geteilter-zustand` „offen, 1×“, das Register zählt 2× (gemessen mit `ls evidence | wc -l`); `BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` 2× stimmt; die Stände stehen am Planungsstand. (5) **Die Coverage streut:** 84,70 bis 84,90 % über die Läufe in Plan und Reports (übernommen), die Schwelle 80 % ist in jedem erfüllt; die DB-Adapter-Coverage 82,56 % (885 von 1072 Statements, gedruckt in jedem Lauf von `make test-store`). (6) **Die `diff`-Zeilen des Suchlaufs bewegten sich mit den Closure-Nachzügen:** `ColumnExclusionPort` 22 → 23 (Übergabe im Plan von `backfill-pfad`), die Beschreibung des Transaktionsbeginns 14 → 16 (zwei Zeilen im neuen Register-Eintrag), die Spec-Klammer 1 → 0; die Sollwerte und Befund-Zellen sind nachgezogen. (7) **Kenntnis:** die Aussage von [`ADR-0113`](../../adr/0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) Festlegung 1 („`ListPending` heute nur nach `requested_at` ordnet“) ist überholt; die ADR ist `Accepted` und bleibt unberührt, `ADR-0127` §Kontext trägt die Einordnung, ihr Zweitschlüssel bleibt richtig.
- **Steering-Loop-Eintrag (Lerneintrag):** *(a) Neuer Sensor, gebaut und verkörpert.* Vier Tests binden, was vorher nur Zusage im Text war: `TestAdministrationRequestListPendingOrdersTiesByRequestID` (Ordnung der Queue bei gleichem Zeitstempel), `TestAdministrationRequestSameTransactionCallsKeepCallOrder` (alle sieben Funktionen in beiden Richtungen, 60 Durchläufe, `make test-store`), `TestAdministrationSameTransactionRemoveThenSetLeavesTheNewRuleLiveAndDerived` (live und abgeleitet gleich, 50 Tabellen, `make test-store`) und `TestApplyAdministrationRequestHandlesEveryKindOfTheClosedSet` (die verarbeitete Menge an die Fälle des `switch` gebunden, `make test`). Anker: `BEO-PGC/queue-ordnung-transaktionsbeginn-und-zufallskennung`, `state.md`, `seit slice-transformationen-antragsweg-usecase`; Grenze: die Ordnung ist der Aufrufzeitpunkt, nicht die Festschreibung (Festlegung 3 der ADR, hergeleitet, nicht erprobt). *(b) Geschärfte Regel, Vorschläge an den Architect (Verkörperung 3b, `v6.9.0` · `regelwerk/modul-08-agentenrollen.md` §Rollen-Sequenz für eine Welle, Schritt 3b: Planner → Architect → Planner; der Planner schärft nicht selbst — Regel-Verkörperung ist eine Entscheidung, und `AGENTS.md` ist Hard-Rule-Text).* Erstens `AGENTS.md` §3.12, Absatz „Verfasser einer ADR“, der Satz zur Mutation in einer Fitness-Function-Zeile — **jetzt nach dem achten Auftreten der Klasse `BEO-PGC/adr-aussage-breiter-als-ihre-messung`** (V-1: die Zeile von `ADR-0127` verallgemeinert die Mutation von „in allen sieben Funktionen“ auf „aus einer Funktion“). Die Klasse ist ≥3× **und** wiederholt trotz des Absatzes (siebtes und achtes Auftreten stehen nach dem Absatz „Verfasser einer ADR“), beide im Zug des Architects. Wortlaut des Vorschlags, geschärft gegenüber der Fassung im Vorgänger-Slice: „Eine in einer Fitness-Function-Zeile genannte Mutation nennt die Menge der Stellen, an denen sie erprobt ist (eine · alle · welche), und die Instanz der Messung; jede Verallgemeinerung darüber hinaus steht als hergeleitet, und ‚der Implementer fährt sie‘ ist eine Erwartung, keine Erprobung.“ Die Ergänzung „Menge der Stellen“ ist neu: die ADR trug die Grenze des Erprobten richtig, die erste Fassung („erprobt oder hergeleitet“) hätte die Zeile bestanden. **Adresse: der nächste Architect-Zug zu einer ADR mit Fitness-Function-Zeile — nicht erst die Closure der Welle**: der Vorschlag stand beim achten Auftreten im Register und nicht im Text, den jeder Zug liest; die Welle hatte in diesem Verlauf drei Architect-Züge mitten im Slice ([`ADR-0125`](../../adr/0125-transformationen-parametertyp-regelform-json.md), [`ADR-0126`](../../adr/0126-transformationen-annahmemenge-rule-spec.md), `ADR-0127`), der Lese-Schritt der Closure kommt nach ihnen. Der Orchestrator beauftragt den Zug (Zielort `AGENTS.md` §3.12 und `.claude/agents/architect.md`); die Gegenentscheidung „akzeptiertes Negativ, LOW, kein Verbraucher“ ist ein Verdikt des Architects; ein Sensor ist ausgeschlossen (Prosa über einen Test). Zweitens `BEO-PGC/formatierungs-drift-ohne-gate`, mit diesem Slice **3× (Schwelle erreicht)**: Vorschlag (1) ein Schritt im Implementer-Ablauf vor dem Handoff, `gofmt -l` im gepinnten Toolchain-Image über die Go-Dateien des eigenen Diffs (`.claude/commands/implement-slice.md`), Ausgabe im Bericht; Kandidaten (2) und (3) und die Wahl im `state.md`; Adresse: der Lese-Schritt der Closure von `welle-transformationen`. Drittens die Entscheidung zum Rest der Klasse `BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue` (2×): drei Optionen im `state.md`, Adresse: der Architect-Zug der Closure von `welle-transformationen`. *(c) Benannte Spec-Lücke:* [`SPEC-019`](../../../../spec/pflichtenheft.md) nennt für ein leeres Schema, einen leeren Tabellennamen und eine leere Spalte weder einen Fehlertext noch den Ort der Prüfung, und die Zeile zur Spalte sagt nicht, dass die Lesung der Queue am Konstruktor endet; Adresse: das Register (`state.md` des genannten Eintrags), Entscheidung beim Architect. **Nachgezogen mit dieser Closure (Verifikation V-4):** die Klammer „(Domänen-Invarianten des Antrags-Konstruktors)“ zu `rule_name`/`rule_spec` in `SPEC-019` sagte den Konstruktor als Ort der Prüfung; sie nennt jetzt den wahren Zustand (die Prüfung liegt in der Verarbeitung, der Antrag endet `failed` mit dem Fehlertext der Tabelle, die Queue lehnt ihn nicht beim Lesen ab), samt Zeile in der Historie der Spec; die Zeile zur Spalte bleibt richtig. *(d) Gelernt, bei 1× keine Regel:* eine Zusage über die **Übereinstimmung zweier Stellen** („Verarbeitung ordnet wie die Ableitung“) steht an keiner der beiden allein; eine Mutation an einer Eingabeseite findet sie nicht (45 Mutationszeilen des Implementers, F-1 blieb), das Lesen beider Stellen gegeneinander schon — Eintrag `BEO-PGC/queue-ordnung-transaktionsbeginn-und-zufallskennung`. Und: die Farbe einer Mutation an einer von N Funktionen hängt an der **Stellung** der Funktion im Ablauf des Tests, nicht an ihrer Rolle (V-2/V-3): eine Zusage über eine Menge von Funktionen bindet jede Funktion an jeder Stellung, an der sie die Aussage tragen soll.
- **Beobachtungs-Register (`../observations/`):** je Anfall eine Datei `evidence/slice-transformationen-antragsweg-usecase.md`, Zähler = Zahl der Dateien (gemessen mit `ls evidence | wc -l` am Stand dieser Closure). *Neue Belege:* `BEO-PGC/adr-aussage-breiter-als-ihre-messung` **8×** (V-1 LOW; verkörpert; Vorschlag (b) oben, `state.md` mit Wortlaut und Adresse), `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` **5×** (F-2 MEDIUM, V-5 INFO; Ausprägung Mechanismus; verkörpert), `BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue` **2×** (F-6 LOW; `state.md` nachgezogen, Rest offen, Adresse Architect-Zug der Welle-Closure), `BEO-PGC/formatierungs-drift-ohne-gate` **3×** (F-3 LOW, F-10 INFO als Ursache; **Schwelle erreicht**, Vorschlag im `state.md`). *Neuer Eintrag, 1×:* `BEO-PGC/queue-ordnung-transaktionsbeginn-und-zufallskennung` (F-1 MEDIUM; offen, Instanz behoben, Klasse unter der Schwelle). *`state.md` nachgezogen ohne neue Datei:* `BEO-PGC/wertabhaengiger-zweiter-waechter-ohne-spec-zeile` (Trigger für den Regelfall eingetreten, F-7 als Weg zur Klasse, Zähler bleibt 1×), `BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` (kein Auftreten, Zähler bleibt 2×). *Deckel-Fälle ohne Datei, Finding-Kennung hier* (verkörpert, ab 10×; vom Reviewer bzw. Verifier vor dem Merge gefunden, Schwere ≤ LOW, bekannter Träger-Typ): V-2 (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`, Deckel bei 14×: die Bindung von `clock_timestamp()` an zwei der sieben Funktionen fehlte, Ausprägung Stellung im Ablauf des Tests; nächstverwandt ist der Positions-Beleg von `slice-transformationen-kern-rename`; in der Fixrunde 3 gebunden); F-4 und V-3 (`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`, Deckel bei 14×: Testkommentar und Plan-Zeile mit einer Mutationsangabe, die nicht trägt — Träger-Typ Testkommentar ist dort belegt, `slice-091`); F-5 (`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`, 13×: die Plan-Aussage widersprach ihrem Nachbarn); F-9 (`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`, Deckel bei 23×: eine bewegliche Zahl mit zu enger Spanne). *Kein Register-Anfall:* F-7, F-8, F-11 (benannte Grenzen, Ausgang in §6), V-6 (Fixrunden ohne zweiten Reviewer-Lauf — **entschieden: keiner**: der Verifier mutierte auf Fixrunde 1 und 2 (S1 bis S8, G1 bis G7), Fixrunde 3 ändert nur Test und Kommentare und ist mit sieben Einzel-Mutationen gebunden; ein späterer Fund an dieser Fläche wäre ein weiterer Beleg in `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`), V-7 (übernommene Messungen, im Report benannt). *Kein Anfall:* `BEO-PGC/test-isolation-geteilter-zustand`, `BEO-PGC/rollen-test-abdeckungsluecken` (der Login-Test zieht die zwei Arten unter den realen Rollen durch die Queue), `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` (die Faltung liegt in der netzlos gemessenen Fläche, die Unterpaket-Regel hat gegriffen), `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (die Sollwert-Abweichungen fing das Werkzeug vor dem Handoff). *Lese-Schritt der Closure von `welle-transformationen`:* aus diesem Slice erreicht neu ein Eintrag 3× ohne Ausgang, `BEO-PGC/formatierungs-drift-ohne-gate`; er trägt Vorschlag und Adresse in seiner `state.md`; die vier mit neuer Datei und bestehendem Ausgang oder unter der Schwelle ändern ihren Ausgang nicht. Offen aus dem Vorgänger-Slice bleibt `BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` (3×, Vorschlag im `state.md`).
- **Folge-Slices:** keine angelegt. Übergaben mit Adresse (gemeldet, Frist: diese Closure, gezogen): `slice-transformationen-backfill-pfad` (`open/`) — der Übergabe-Block zu den Kosten der Regelstand-Lesung je Block (F-11) und zur Klasse des Faltungsfehlers im Run sowie neu das Risiko „Fenster zwischen Regel-Setzbarkeit und Backfill-Bindung“ in dessen §6; `slice-transformationen-map-value` (`open/`) — der Übergabe-Block zum Zwischenzustand der Regeltyp-Menge und zur Rückfall-Grenze auf einen älteren Binärstand; `slice-transformationen-betriebsdoku` (`open/`) — trägt bereits die Aufrufform `::json` und die Aussage zur Aufruf-Reihenfolge von `ADR-0127` Folgepflicht 4 (gelesen am Stand `2c22334f`), neu die Grenze der Dauerhaftigkeit; `slice-transformationen-e2e-abhilfe` (`open/`) belegt die Kollision mit einer später hinzugefügten Spalte (Risiko K3). Die **Start-Bedingung von `backfill-pfad`** ist mit diesem Move erfüllt (gelesen in dessen §4: `slice-backfill-run-usecase` und `slice-backfill-e2e` liegen in `done/`, `slice-transformationen-antragsweg-usecase` liegt nach dem Move in `done/`, in `in-progress/` liegt nur die Roadmap); der Plan ist nicht geändert.
- **Risiken aus §6:** je ein Ausgang, mit Beleg in §6. *Entfallen:* Fenster „Funktion vorhanden, Wirkung fehlt“ (mit dieser Closure) · Port-Schnitt · Store-Tests teilen Zustand · Kosten und Rollen der Spaltenlisten-Abfrage. *Eingetreten, im Slice aufgelöst:* Ordnung von Verarbeitung und Ableitung (F-1, `ADR-0127`; die drei Grenzen der ADR bleiben akzeptiertes Negativ in der ADR, ohne Register-Eintrag, weil keine aufgetreten ist). *Weiter offen:* Auseinanderlaufen bei Vermerk-Fehler mit kollidierendem Folgeantrag (Register `BEO-PGC/wertabhaengiger-zweiter-waechter-ohne-spec-zeile`) · leeres Schema, leere Tabelle, leere Spalte stallt die Queue (Register `BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue`) · Kosten der Regelstand-Lesung je Block (Nehmer `backfill-pfad`) · K3 gegen die Spaltenliste zum Antragszeitpunkt (Nehmer `e2e-abhilfe`) · Fenster Backfill (Nehmer `backfill-pfad` §6, endet mit dessen Closure) · Zwischenzustand der Regeltyp-Menge (Nehmer `map-value`, endet mit dessen Closure).
- **Drei Paarungen:** dieser Slice gehört zu [welle-transformationen](../welle-transformationen.md) (offen) — die Closure der Welle prüft sie mit; die Slice-Closure trägt sie zusätzlich jetzt: *Anker:* der Lerneintrag (a) trägt den Anker `BEO-PGC/queue-ordnung-transaktionsbeginn-und-zufallskennung` `state.md` mit `seit slice-transformationen-antragsweg-usecase` (am Ort existent); (b) sind Vorschläge an den Architect, kein `liegt in`-Feld eines verkörperten Ziels, ihre Zielorte (`AGENTS.md` §3.12, `.claude/agents/architect.md`, `.claude/commands/implement-slice.md`) existieren; *Folge-Slice:* keiner genannt; die vier Übergabe-Adressen `slice-transformationen-backfill-pfad`, `slice-transformationen-map-value`, `slice-transformationen-betriebsdoku` und `slice-transformationen-e2e-abhilfe` existieren als Dateien in `open/`; *Register:* jede genannte Kennung `BEO-PGC/<slug>` existiert als Verzeichnis mit nicht leerem `evidence/` (geprüft mit `ls docs/plan/planning/observations/BEO-PGC/<slug>/evidence`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Use Cases, Ports, Store-Adapter und Composition Root
sind keine eigenen Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger` (verkörpert — der Regelstand
ist ein dauerhafter, tabellen-scoped Träger, der Assembler-Cache ist Laufzeit),
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 6×, DoD Punkt
1), `BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` (offen, 2×,
einschlägig — DoD Punkt 2),
`BEO-PGC/dockerignore-default-deny-blockiert-neuen-pfad` (offen, 1×,
Plan-Zeile), `BEO-PGC/test-isolation-geteilter-zustand` (offen, 1×, Risiko §6),
`BEO-PGC/adapter-fehler-ausgang` (offen, 2×, gesichtet — Antrags-Fehler enden
`failed`, ein Retry ist nicht Teil),
`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` (offen, 2×: zwei Dateien
in `evidence/`, gemessen 2026-09-26; dieser Slice trägt kein Auftreten, §6),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf §3).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
