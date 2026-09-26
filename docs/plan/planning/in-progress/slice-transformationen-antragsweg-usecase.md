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
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
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
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure der Welle
      [welle-transformationen](../welle-transformationen.md) (die Roadmap führt
      sie unter *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/application/port/inbound/transformation.go` | neu | `SetTransformationUseCase` (liefert die geprüfte `model.Transformation`, die der Aufrufer in die Bindung einträgt) und `RemoveTransformationUseCase` samt `SetTransformationCommand`/`RemoveTransformationCommand` ([`ADR-0028`](../../adr/0028-inbound-use-cases.md), Transport-Typen am Port); die Namen der Use Cases stehen in [`ARC-003`](../../../../spec/architecture.md) (Antragsart-Tabelle). |
| `internal/application/port/outbound/transformation.go` | neu | **ein** Outbound Port `TransformationPort` mit zwei Lese-Methoden (`TransformationRules`: Regelstand je Tabelle einer Quelle; `SourceColumns`: Spaltennamen der Quelltabelle) — die Zählung von [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) („ein neuer Outbound Port“). Kein Abweichen: der Schnitt folgt dem Vorbild `ColumnExclusionPort` (`ColumnExists` aus dem Katalog und `ExcludedColumns` aus den Antrags-Zeilen in einem Port); beide Methoden beantworten dieselbe Frage („welche Regeln dürfen die Spalten dieser Tabelle tragen“) und laufen über dieselbe Adapter-Instanz. Offene Frage an den Architect, keine Entscheidung: ob [`ADR-0034`](../../adr/0034-ports-nach-faehigkeiten.md) (Ports nach Fähigkeiten) den Katalog-Zugriff (Spaltenliste) und die Antrags-Ableitung (Regelstand) als zwei Fähigkeiten zählt; `ColumnExclusionPort` schneidet beide Objektklassen ebenfalls in einen Port. |
| `internal/application/usecase/settransformation/service.go`, `internal/application/usecase/removetransformation/service.go` (+ je `service_test.go`) | neu | die zwei Use Cases in der netzlos gemessenen Fläche: `Set` prüft die Formzeilen (Regelname, Regelform) **vor** dem ersten Lesen des Ports, dann K1 bis K3 über `TransformationSpec.CheckConflicts`, dann K4 gegen die Spaltenliste; `Remove` prüft den Namen und K4 (`Regelname nicht geführt`). Der Fehlertext ist der der Spec (Klartext, Doppelpunkt, Adresse), der Grund bleibt über `errors.Is` erreichbar. Beide Use Cases schreiben den Regelstand nicht (er ist die Ableitung aus den `applied`-Zeilen). |
| `internal/domain/model/transformationspec.go` (+ `transformationspec_test.go`) | neu | Parser `ParseTransformationSpec` (strikt: UTF-8, JSON-Objekt, `kind`, unbekannter Regeltyp, unbekannter Schlüssel in aufsteigender Ordnung, Pflichtschlüssel und Bezeichner-Form über `NewRenameColumn`), `CheckRuleName` (Alphabet `[a-z0-9_]{1,63}`), `TransformationSpec.CheckConflicts` (K1 bis K3, zeichengenau, in der Reihenfolge der Spec), `Build` und die reine Faltung `FoldTransformations` (`applied`-Zeilen → Regelstand; Set trägt ein und ersetzt nach Namen, Remove nimmt heraus; eine nicht mehr lesbare Zeile endet als Fehler). Ein Zielname gleich der Quellspalte ist keine Formverletzung des Parsers, sondern K3 (`ErrTargetCollidesWithColumn`) in der Stellung von K3. |
| `internal/domain/model/administrationrequest.go` (+ Test), `internal/domain/errors/errors.go` | update | der Konstruktor `NewAdministrationRequest` lehnt die Transformations-Antragsarten nicht mehr wegen leerem Regelnamen oder leerer Regelform ab (**Stelle der Prüfung: Verarbeiten im Use Case, nicht Lesen**); `AdministrationRequestKinds()` zählt die sieben Arten auf; neun Sentinels tragen die Ablehnungsgründe (Klartext der Spec-Zeilen). |
| `internal/adapters/driven/postgresstorage/queries/queries.go`, `sqlexec/translate.go`, `tableactivation.go` (+ `sqlexec/translate_test.go`, `administrationrequest_test.go`) | update | `SelectAppliedTransformationRequests` (Ordnung `requested_at`, dann `administration_request_id`), `SelectTableColumns` (`information_schema.columns`, `ORDER BY ordinal_position`); `sqlexec.ReadTransformationRules` (Zeilen je Tabelle, Faltung in der Domäne, eine Tabelle ohne verbleibende Regel trägt keinen Eintrag) und `ReadSourceColumns` in der netzlos gemessenen Fläche; `TableActivationAdapter` implementiert den Port (dieselbe Instanz über `CDC_ADMIN_DSN`, `SELECT` der Rolle `cdc_admin` auf der Antrags-Tabelle wie `ExcludedColumns`; `cdc_capture` trägt dort kein Recht). Die Faltung liegt an **einer** Stelle, `model.FoldTransformations`. |
| `internal/bootstrap/wiring.go` | update | Zweige `set_transformation`/`remove_transformation` in `applyAdministrationRequest` (Use Case, danach `Assembler.SetTransformation`/`RemoveTransformation`; der Kommentar am Set-Zweig nennt die Grenze von K3), drei Felder in `administrationDeps` (`transformations`, `setTransformations`, `removeTransformations`), `activatedTableBindings` und der Aktivierungs-Zweig tragen den abgeleiteten Regelstand neben dem Ausschlussstand, die Doc-Kommentare zählen die Arten und beschreiben die Zweige; `processedAdministrationKinds` ist eine Funktion über `model.AdministrationRequestKinds()` (eine Quelle statt einer handgeführten Zeichenkette). |
| `internal/bootstrap/administration_internal_test.go`, `wiring_rest_internal_test.go`, `backfill_internal_test.go` | update | Whitebox mit Fakes: Nachtrag live (Set, Remove), K1 bis K4 und die Formzeilen bis zum `failed`-Vermerk mit dem Fehlertext der Spec, ungültige Zeilen neben einer gültigen, Idempotenz der Wiederholung, Tabelle ohne Bindung, Prozessstart (`activatedTableBindings`), Aktivierungs-Zweig (Deaktivierung, dann Aktivierung), Lesefehler des Regelstandes, `default`-Zweig und die Bindung der Aufzählung an die Fälle des `switch` (`TestApplyAdministrationRequestHandlesEveryKindOfTheClosedSet`); die bestehenden Fixtures tragen den Fake des neuen Ports. Die zwei Tests, die das Fenster „Funktion vorhanden, Wirkung fehlt“ banden (`…RejectsUnprocessedKind`, `…MarksTransformationRequestsFailed`), sind ersetzt. |
| `internal/bootstrap/administration_roles_internal_test.go` | update (Übergabe aus `slice-backfill-run-store`) | der Login-Test zieht die zwei neuen Arten unter dem `cdc_admin`-/`cdc_capture`-Login durch die Queue: Set, Remove, K1 bis K4 mit Fehlertext, Regelstand und Spaltenliste über den realen Adapter, eine gültige Zeile neben Zeilen mit fehlendem Regelnamen und fehlender Regelform (die Queue bleibt lesbar). |
| `.dockerignore` | prüfen — keine Änderung | alle neuen Pfade liegen unter `internal/` (`!internal/`); keine Datei außerhalb des Go-Baums kommt hinzu. Belegt durch `make coverage-gate` und `make test` (Bau-Kontext trägt die Pakete). |
| `spec/pflichtenheft.md` `SPEC-019` (Zeile „`rule_name` ist … Pflicht … (Domänen-Invarianten des Antrags-Konstruktors)“) | gemeldet, nicht geändert | die Klammer nennt den Konstruktor als Ort der Prüfung; nach diesem Slice liegt sie im Use Case (Formzeile `Regelname ist ungültig`). Fremde Datei; Meldung an den Planner, Frist: Closure dieses Slice. |
| Spalten-Antragsarten mit leerer Spalte (`cdc.exclude_column(…, NULL)`) | keine Änderung (Grenze, benannt) | der Konstruktor lehnt `exclude_column`/`include_column` mit leerer Spalte weiterhin beim Lesen ab (`SPEC-019`: „Domänen-Invariante des Antrags-Konstruktors“ für die Spalten-Antragsarten); die Stelle der Prüfung dieses Slice führt sie **nicht** mit. Adresse: `BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue` (geplant, 1×), Entscheidung beim Planner. |
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
| Stellen, die eine Bindung anlegen | `git grep -n 'TableBinding{' -- internal ':!*_test.go'` (Zeilen 1–2 des Blocks) | **Gefunden.** Parent 5, Diff 5 (Zeilenzahl unverändert, Inhalt geändert): `config_file.go:250` und `parseTables` in `wiring.go` (die Seeds aus Konfigurationsdatei und `CDC_TABLES`: Bindung ohne Ausschluss- und Regelstand, der Prozessstart ersetzt sie über `activatedTableBindings`), `activatedTableBindings` und der Aktivierungs-Zweig (beide tragen den Regelstand mit). **Nichtgefunden:** keine sechste Stelle, die eine Bindung anlegt; `Assembler.AddBinding` erhält Ausschluss- und Regelstand einer getragenen Bindung selbst (`kern-rename`). | Seeds bleiben ohne Regelstand (Begründung in §3 Kommentar an `activatedTableBindings` und am Aufruf in `Run`); die zwei anderen Stellen tragen ihn. |
| Stellen, die den Ausschlussstand mitführen, und ihr Regelstand-Zwilling | `git grep -n ExcludedColumns -- internal ':!*_test.go'` (Zeilen 3–4) und `git grep -n 'Transformations:' -- internal ':!*_test.go'` (Zeilen 5–6) | **Gefunden.** `ExcludedColumns` Parent 25, Diff 26 (eine Zeile mehr: der Doc-Kommentar von `TransformationRules`, der auf `ExcludedColumns` verweist); `Transformations:` Parent 0, Diff 4 (`activatedTableBindings` und Aktivierungs-Zweig, dazu die Felder `setTransformations:`/`removeTransformations:` im Literal von `Run`, die das Muster mitzählt). Träger, die `ExcludedColumns` lesen oder setzen: `mapper.go` (`kern-rename` trägt dort den Regelstand-Zwilling), `wiring.go` (zwei Stellen, beide mit Zwilling), `backfill/service.go` (liest den Ausschlussstand, der Regelstand kommt mit `slice-transformationen-backfill-pfad`), `sqlexec`/`tableactivation.go`/`columnexclusion.go` (Ableitung, Zwilling `ReadTransformationRules`/`TransformationRules`). Test-Fixtures: die neuen Fixtures tragen `Transformations` über die Bindung. **Nichtgefunden:** kein Träger im Erfassungs- oder Zustellpfad, der den Ausschlussstand ohne Regelstand kopiert. | `backfill/service.go` bleibt (Adresse `backfill-pfad`, dessen §2). |
| Beschreibung der Antragsarten im Doc-Kommentar von `applyAdministrationRequest` | `git grep -n 'Spalten-Antragsarten' -- internal/bootstrap/wiring.go` (Zeilen 7–8) | **Gefunden.** Parent 1, Diff 2; der Doc-Kommentar beschreibt die zwei neuen Zweige (Use Case, Nachtrag, Tabelle ohne Bindung, Idempotenz, Zeilen mit fehlenden Regelfeldern) und trägt den Aktivierungs-Zweig mit Regelstand. **Nichtgefunden:** keine zweite Beschreibung der Antragsarten in `internal/bootstrap`, die die zwei Arten als „nicht verarbeitet“ führt (`git grep -n 'nicht verarbeitet' -- internal/bootstrap` ohne Test-Treffer zu diesem Gegenstand). | nachgezogen. |
| Port-Übersichten in Doku | `git grep -n ColumnExclusionPort -- docs spec harness internal ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/in-progress' ':!*_test.go'` (Zeilen 9–10) | **Gefunden.** Parent 21, Diff 22 (die eine Zeile mehr: der Doc-Kommentar von `TransformationPort.SourceColumns`, der die Katalog-Lesart von `ColumnExclusionPort.ColumnExists` nennt; Treffer sind Code-Träger, drei ADRs, ein Register-State und `spec/architecture.md:274`, das den Spaltenausschluss beschreibt und richtig bleibt). **Nichtgefunden:** keine Aufzählung der Outbound Ports (Port-Übersicht) in `docs/`, `spec/` oder `harness/`: `spec/architecture.md` nennt Ports im Text und in den Diagrammen, `ARC-004` führt keine Liste; die zwei Use Cases und „ein Outbound Port“ stehen dort bereits (`spec/architecture.md`, Antragsart-Tabelle, Absatz „Transformations-Antragsarten“). | keine Änderung. |
| Zweite Quelle der verarbeiteten Menge | `git grep -n processedAdministrationKinds -- internal` (Zeilen 11–12) und das Zählwort „fünf“ in `internal/bootstrap` (Zeilen 13–14) | **Gefunden.** Parent 5 Treffer (Konstante, Fehlertext, drei Test-Zeilen), Diff 3 (Funktion, Fehlertext, Doc-Kommentar); „fünf“ in `internal/bootstrap` Parent 12, Diff 9: die drei verschwundenen Treffer sind der Test und die Kommentare der Fünf-Arten-Menge; die neun übrigen zählen fremde Mengen (fünf Umgebungsvariablen, fünf Lese-Views, fünf Tabellen). **Nichtgefunden:** keine handgeführte Aufzählung der verarbeiteten Antragsarten mehr außerhalb des einen Testliterals, das den Fehlertext bindet (`TestApplyAdministrationRequestRejectsKindOutsideTheClosedSet`). | Menge abgeleitet, an die Fälle des `switch` gebunden. |
| Felder von `administrationDeps` | `git grep -n 'administrationDeps{' -- internal` (Zeilen 15–16) und `git grep -n 'transformations:' -- internal/bootstrap` (Zeilen 17–18) | **Gefunden.** Literale von `administrationDeps` Parent 21, Diff 20 (`administration_internal_test.go` 12 → 11: zwei Literale der ersetzten Tests entfallen, `ruleFixture` fügt eines hinzu; jede andere Datei trägt dieselbe Zahl, das Literal des Login-Tests ist um drei Felder erweitert); Setzungen des Feldes `transformations:` Parent 0, Diff 10 (sieben Fixtures bestehender Tests, das Literal des Login-Tests, `ruleFixture` und die Verdrahtung in `Run`). **Nichtgefunden:** kein Literal, dessen Pfad den Aktivierungs-Zweig erreicht, ohne das Feld (der erste Lauf von `make test` fand zwei Literale in `wiring_rest_internal_test.go` als Nil-Dereferenzierung, Exit 2, nachgezogen). | nachgezogen. |
| Fehlerklassen-Abbildung | Lesen von `classifyRunError` in `internal/bootstrap/wiring.go` | **Gefunden.** `classifyRunError` bildet Fehler des Capture-Pfads ab (Assembler, Replikation, Speicher); Antrags-Fehler laufen durch `processAdministrationRequests` in den `failed`-Vermerk und erreichen `classifyRunError` nicht. **Nichtgefunden:** keine Stelle, an der ein Antrags-Fehler den Prozess beendet. | keine neue Abbildung nötig. |
| Handbuch und Spec (Meldung) | `git grep -n -e set_transformation -e remove_transformation -- docs/user` (Zeilen 19–20); `git grep -n 'Domänen-Invarianten des Antrags-Konstruktors' -- spec` (Zeilen 21–22) | **Gefunden.** `docs/user`: Parent 0, Diff 0 (die Betreiber-Oberfläche ist unbeschrieben, bis `betriebsdoku` sie beschreibt); `spec`: Parent 1, Diff 1 — `SPEC-019`, Zeile „`rule_name` ist … Pflicht … (Domänen-Invarianten des Antrags-Konstruktors)“ nennt den Konstruktor als Ort einer Prüfung, die dieser Slice in den Use Case legt. **Nichtgefunden:** kein weiterer Träger der Aussage „der Konstruktor lehnt leeren Regelnamen ab“ (Register-Records und Pläne nennen den Zustand am Parent). | Spec-Zeile gemeldet (Planner, Frist: Closure dieses Slice); Handbuch: Aufschub mit Adresse `slice-transformationen-betriebsdoku`. |

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
diff 22 -n ColumnExclusionPort -- docs spec harness internal :!docs/reviews :!docs/plan/planning/done :!docs/plan/planning/in-progress :!*_test.go
80eefead 5 -n processedAdministrationKinds -- internal
diff 3 -n processedAdministrationKinds -- internal
80eefead 12 -n fünf -- internal/bootstrap
diff 9 -n fünf -- internal/bootstrap
80eefead 21 -n 'administrationDeps{' -- internal
diff 20 -n 'administrationDeps{' -- internal
80eefead 0 -n 'transformations:' -- internal/bootstrap
diff 10 -n 'transformations:' -- internal/bootstrap
80eefead 0 -n -e set_transformation -e remove_transformation -- docs/user
diff 0 -n -e set_transformation -e remove_transformation -- docs/user
80eefead 1 -n 'Domänen-Invarianten des Antrags-Konstruktors' -- spec
diff 1 -n 'Domänen-Invarianten des Antrags-Konstruktors' -- spec
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

**Läufe des Implementer-Laufs** (Exit-Code je Lauf ungefiltert in eine Log-Datei geschrieben und
gesondert gelesen, [`AGENTS.md`](../../../../AGENTS.md) §3.9; ein schwerer Docker-Lauf zugleich):

- `make test` (Race-Detector) Exit 0, 44 Pakete `ok`; `make test-store` Exit 0, gedruckt:
  `DB-Adapter-Coverage: 82.56% (gedeckt 885 von 1072 Statements; Profile gemergt: store,replication)`,
  `db-coverage: OK — DB-Adapter-Coverage 82.56% erfuellt Schwelle 80%` (der Login-Test
  `TestAdministrationPathRunsUnderLeastPrivilegeLogins` läuft real und ist in den Mutationen rot gesehen).
- `make a-check` Exit 0, `gesamt: 0 Befund(e)`; `make coverage-gate` Exit 0,
  `coverage-gate: OK — Coverage 84.80% erfüllt Schwelle 80%`.
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
  (Wiederholung desselben Antrags ist folgenlos). **Ausgang:** *(bei Closure)*
- **K3 prüft gegen die Spaltenliste zum Antragszeitpunkt**; eine spätere
  Spalten-Erweiterung kann den Zielnamen kollidieren lassen. Das ist der Fall,
  den die Prüfung im Assembler (`kern-rename`) fängt und den `e2e-abhilfe` real
  belegt; hier wird er nicht verhindert. *Erwartet, zu belegen durch:*
  Kommentar am K3-Zweig, der die Grenze nennt, und der Verweis in §7.
  **Ausgang:** *(bei Closure)*
- **Fenster: Funktion vorhanden, Wirkung fehlt** (Übergabe aus
  `antragsweg-schema`). Seit dessen Closure schreiben `cdc.set_transformation`
  und `cdc.remove_transformation` einen Antrag, den `applyAdministrationRequest`
  im `default`-Zweig als `failed` mit dem Fehlertext der verarbeiteten
  Antragsarten vermerkt: sichtbar, keine Regel, keine Wirkung. Das Fenster endet
  mit der Closure dieses Slice. *Erwartet, zu belegen durch:* der Test des
  `default`-Zweigs färbt sich rot, sobald der Zweig für `set_transformation`
  entfällt (die zweite Quelle der verarbeiteten Menge, Übergabe oben).
  **Ausgang:** *(bei Closure: entfallen mit der Closure dieses Slice)*
- **Fenster zwischen Regel-Setzbarkeit und Backfill-Bindung.** Ab diesem Slice
  sind Regeln setzbar; ein Backfill-Run, der vor `backfill-pfad` läuft,
  lieferte die Rohform (Welle §5). Der Slice ändert das Backfill-Verhalten
  nicht. *Erwartet, zu belegen durch:* die Reihenfolge der Welle
  (`backfill-pfad` folgt diesem Slice unmittelbar) und die Benennung im
  Bericht. **Ausgang:** *(bei Closure: entfallen mit der Closure von
  `slice-transformationen-backfill-pfad`)*
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
  aus `map-value` §2. **Ausgang:** *(bei Closure: entfallen mit der Closure von
  `slice-transformationen-map-value`)*
- **Der Port-Schnitt weicht von der ADR-Zählung ab** („ein neuer Outbound
  Port“,
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  §Schärft; `BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab`, offen, 1×).
  *Erwartet, zu belegen durch:* Begründung im Plan-Nachzug und Review-Prüfung;
  eine Abweichung wird als offene Frage geführt, nicht als Entscheidung.
  **Ausgang:** *(bei Closure)*
- **Store-Tests teilen Zustand** (`BEO-PGC/test-isolation-geteilter-zustand`,
  offen, 1×): unskopierte `DELETE`/`DROP SCHEMA CASCADE` in einem Store-Test
  könnten die neuen Zeilen zerstören. *Erwartet, zu belegen durch:* skopierte
  Bereinigung je Test. **Ausgang:** *(bei Closure)*
- **Kosten und Rollen der Spaltenlisten-Abfrage**: sie läuft je
  `set_transformation`-Antrag über die Administrations-Verbindung
  (`cdc_admin`). *Erwartet, zu belegen durch:* `make test-store` unter der
  realen Rolle (gleiche Katalog-Lesart wie `ColumnExists`, gelesen am Bestand).
  **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen“ als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-transformationen](../welle-transformationen.md) (offen) — die Prüfung
  läuft regelkonform bei deren Closure.

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
`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` (offen, 1×, Risiko §6),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf §3).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
