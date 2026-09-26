# Slice transformationen-antragsweg-schema: Antragsweg, Schema — Antragsarten `set_transformation`/`remove_transformation`, Spalten `rule_name`/`rule_spec`, zwei SQL-Funktionen, Grants, Idempotenz-Guard

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-transformationen](../welle-transformationen.md).

**Bezug:** [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) (Konfiguration der
Transformation), [`LH-FA-ADM-001`](../../../../spec/lastenheft.md)
(Administration über SQL), [`LH-QA-SEC-002`](../../../../spec/lastenheft.md)
(Least-Privilege der Rollen),
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Teilfrage 1 und Folgepflicht 3,
[`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(Antrags-Queue und Live-Reload),
[`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md)
(Schemamigrationen — Idempotenz-Guard),
[`ADR-0047`](../../adr/0047-rollenspezifische-dsn-verdrahtung.md) (Rollen),
[`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
(keine Domänenlogik in SQL).

**Berührte Spec-Stellen:** [`SPEC-019`](../../../../spec/pflichtenheft.md)
(Antrags-Datensatz, durch `slice-transformationen-spec-nachzug`),
[`ARC-005`](../../../../spec/architecture.md) (SQL-Funktionen als Driving
Adapter) — gelesen, nicht geändert.

**Verantwortlich:** Implementer-Agent, 2026-09-26.

**Autor:** Planner-Agent, Welle-Eröffnung
[welle-transformationen](../welle-transformationen.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die Antrags-Queue nimmt die zwei neuen Antragsarten entgegen: das
Schema trägt die Spalten `rule_name` und `rule_spec`, die geschlossene
`request_kind`-Menge trägt zwei Werte mehr, und die SQL-Funktionen
`cdc.set_transformation(source_id, schema_name, table_name, rule_name,
rule_spec json)` und `cdc.remove_transformation(source_id, schema_name,
table_name, rule_name)` — die Parameterlisten legt dieser Slice fest, die Spec
([`SPEC-019`](../../../../spec/pflichtenheft.md)) nennt die Funktionen ohne
Parameter; `rule_spec` ist ein `json`-Parameter, der als `jsonb` in die Spalte
`rule_spec` geht (§3 Abweichung von der Wortwahl `jsonb` in
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Teilfrage 1) —
schreiben **nur** den Antrag und senden `pg_notify` —
beide ausschließlich der Rolle `cdc_admin` ausführbar (`SECURITY DEFINER`,
gepinnter `search_path`, `REVOKE … FROM PUBLIC`). Die Domäne kennt die
Antragsarten, der Store liest die zwei Spalten. Die Verarbeitung eines solchen
Antrags trägt der Folge-Slice.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Use Cases, K1–K4, Regelstand-Ableitung, Verdrahtung in
  `applyAdministrationRequest`** — `antragsweg-usecase`. Bis dahin endet ein
  Antrag der neuen Arten im `default`-Zweig von `applyAdministrationRequest`
  als `failed` mit Text: sichtbar, nicht still (§6).
- **Regel-Parsing in SQL** — die Funktionen prüfen keine Regelform; die
  Validierung liegt vollständig in Go
  ([`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md),
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 1).
- **Wirkung im Assembler** — `kern-rename`.
- **Betreiber-Handbuch** — der Slice liefert Betreiber-Oberfläche (zwei
  SQL-Funktionen), die noch nichts bewirkt; ihre Beschreibung ist aufgeschoben
  mit Adresse `slice-transformationen-betriebsdoku`, dessen §2 den
  aufgeschobenen Gegenstand vollständig nennt (SQL-Funktionen, Rollen, Wirkung,
  Abhilfe, Glossar, Änderungshistorie;
  `BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an`, offen, 2×). Aussagen über
  eine Wirkung stehen erst, wenn sie belegt ist
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12).
- **Eine Regel-Sicht (View)** —
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 6: der Regelstand ist über `cdc.administration_request` lesbar,
  eine eigene Sicht ist nicht Teil.

## 2. Definition of Done

- [x] Antragsweg im Schema: `cdc.set_transformation` und
      `cdc.remove_transformation` schreiben ausschließlich einen Antrag (Art,
      `rule_name`, `rule_spec` nur bei `set_transformation`, `column_name`
      leer) und senden `pg_notify`; `PUBLIC` hat kein `EXECUTE`, ein Login ohne
      `cdc_admin`-Mitgliedschaft scheitert mit „permission denied for
      function“; die `request_kind`-Menge trägt genau die Werte des
      Parent-Stands plus die zwei neuen; `rule_name` und `rule_spec` sind
      nullable; der zweite Schema-Träger
      `internal/adapters/driven/postgresstorage/schema.sql` trägt
      `cdc.administration_request` nicht (am Start gemessen, §3) und bleibt
      unverändert. *Zu belegen durch:* `make test-store`, `make schema-rollout`
      zweimal hintereinander (Exit 0) und
      `tools/harness/run-schema-rollout-guard-test.sh` (alle Läufe);
      `plan.yaml` und `down.sql` regeneriert, falls der Rollout sie verändert.
      Der **Alt-Tag-Lauf** desselben Skripts
      ([`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md)
      Entscheidung 7: das Schema des jüngsten `v*`-Tags per `git archive`
      ausrollen, danach den Arbeitsbaum — Exit 0 zweimal, der zuvor eingefügte
      Datenstand über `cdc.changes` lesbar; der Bericht nennt den Tag und die
      gedruckten Exit-Codes) trägt den Upgrade-Beleg für die zwei nullable
      Spalten (`jsonb` statt `text` ist über einen Alt-Bestand ungemessen) und
      die zwei Funktionen.
- [x] Idempotenz-Guard und Rollen-Test: `knownForeignObjects` in
      `tools/schema/rolloutguard/guard.go` trägt beide Funktionen
      (Signatur-Schreibweise am realen `--plan-only`-Report gemessen, nicht
      angenommen — der Parameter der Regelform ist neu gegenüber den
      bestehenden Funktionen; sein Typ ist `json`, nicht `jsonb`, §3
      Abweichung), `guard_test.go` deckt sie;
      `internal/bootstrap/roles_rollout_file_internal_test.go` liest den
      Grant-Text beider Funktionen und färbt sich bei entferntem
      `GRANT`/fehlendem `REVOKE` rot. *Zu belegen durch:* `make test` (beide
      Tests) und je eine Mutation, die den Test rot färbt.
- [x] Domäne und Store: `model.AdministrationRequestKind` trägt die zwei Arten;
      der Antrags-Konstruktor erzwingt `rule_name` für beide Arten und
      `rule_spec` für `set_transformation` (Invariante wie bei `column_name`);
      der `AdministrationRequestAdapter` liest die zwei Spalten (`ListPending`
      liefert sie); ein Antrag der neuen Arten wird in
      `applyAdministrationRequest` als `failed` mit Text vermerkt, und der
      Fehlertext nennt die tatsächliche Menge. *Zu belegen durch:* `make test`
      (Whitebox in `internal/bootstrap`, Domänen-Test) und `make test-store`
      (Round-Trip der zwei Spalten, skopierte Bereinigung —
      `BEO-PGC/test-isolation-geteilter-zustand`, offen, 1×).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: aufgeschoben mit Adresse — die Betreiber-Oberfläche entsteht
      hier, ihre Beschreibung trägt `slice-transformationen-betriebsdoku`
      (siehe §1); die Änderungshistorie des Handbuchs trägt dieser Slice
      **nicht**, weil er das Handbuch nicht berührt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
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
| `tools/schema/schema.yaml` | update | Spalten `rule_name text`, `rule_spec jsonb` (nullable) im neutralen Modell von `administration_request`; die `description` nennt sie. |
| `tools/schema/nacharbeit-administration.sql` | update | CHECK-Menge um zwei Werte (dort steht die Menge außerhalb des neutralen Modells), die zwei Funktionen, `REVOKE … FROM PUBLIC`, `GRANT EXECUTE … TO cdc_admin`, Kopfkommentar. **Abweichung vom Plan-Wortlaut:** `cdc.set_transformation` trägt den Parameter `rule_spec` als `json`, nicht als `jsonb` (`ADR-0112` Teilfrage 1 und `SPEC-019` nennen `jsonb`); die Funktion schreibt ihn als `jsonb` in die Spalte. Grund, gemessen: d-migrate 1.3.1 meldet jede Funktion mit `json`- oder `jsonb`-Parameter als `in:json` und rendert ihren Abbau im zweiten Rollout als `DROP FUNCTION … (…, json)`; mit einem `jsonb`-Parameter endete der zweite `make schema-rollout` mit d-migrate-Exit 5 (`function set_transformation(text, text, text, text, json) does not exist`), mit `json` endet er mit Exit 0. Folge für Aufrufer: ein Literal und ein `::json`-Wert werden angenommen, ein `::jsonb`-Wert nicht (gemessen, „function … does not exist“). Der Widerspruch zwischen Wortlaut und Rollout-Mechanik gehört einer Architect-Entscheidung; dieser Slice hält den kleinsten Stand, der den Rollout idempotent lässt (§6). |
| `internal/adapters/driven/postgresstorage/schema.sql` | keine Änderung (Nicht-Realisierung) | Der Plan nahm an, der zweite Schema-Träger trage die Spalten. Am Start gemessen trägt die Datei `cdc.administration_request` nicht (sie führt allein `source`, `source_table`, `schema_version`, `transaction`, `change`; Kopfkommentar der Datei): `git grep -n administration_request 2f5ed3dc -- internal/adapters/driven/postgresstorage/schema.sql` liefert keine Zeile. Es gibt keinen zweiten Träger der Spaltenform (Suchlauf-Zeile unten); das Risiko „zwei Schema-Träger driften“ (§6) entfällt damit. |
| `tools/schema/plan.yaml`, `tools/schema/down.sql` | regeneriert | Beide sind der Beleg des Betriebs-Rollouts gegen eine frische Datenbank (`harness/targets/schema-rollout.md` §Erzeugnisse); der `CREATE TABLE "administration_request"` im Report trägt die zwei neuen Spalten, der Artefakt-Hash im Kopf von `down.sql` folgt der Schema-Form. Entscheidung: regenerieren, aus dem ersten von zwei aufeinanderfolgenden `make schema-rollout`-Läufen gegen eine frische Wegwerf-Datenbank (der zweite überschreibt beide Dateien mit dem Report des Folgelaufs); die Zielzeile im Report bleibt die des Compose-Rollouts, wie am Parent ([`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md)). |
| `.dockerignore` | update | Eine Negation für `tools/schema/nacharbeit-administration.sql`: der neue Rollen-Test liest sie, und die `coverage`-Stufe sieht sonst nur die Dateien der Allow-List des Build-Kontexts (Klasse [`ADR-0085`](../../adr/0085-build-kontext-ausnahme-test-only-zweck.md) Festlegung 3: gelesen, nicht gebaut; genau eine Datei; der Kommentar nennt den Leser; Image unberührt, §6). Gefunden durch `make coverage-gate`: `open /src/tools/schema/nacharbeit-administration.sql: no such file or directory`. |
| `tools/schema/rolloutguard/guard.go` (+ `guard_test.go`) | update | Einträge beider Funktionen in der Schreibweise des Reports (gemessen am `--plan-only`-Report gegen ein migriertes Ziel: `set_transformation(in:text,in:text,in:text,in:text,in:json)`, `remove_transformation(in:text,in:text,in:text,in:text)`; `guard_test.go` bindet die Schreibweise mit einem Test, der `in:jsonb` ablehnt); der Kommentar, der die Objekte zählt, zählt neu (neun); der Kommentar am Feld `allowDestructive` („nur bekannte Fremdobjekte blockieren“) nennt neben den bekannten Fremdobjekten die Klasse „View-Signatur“ (Adresse der Meldung aus `slice-backfill-change-origin`, `BEO-PGC/gemeldete-ungenauigkeit-ohne-traeger`). |
| `internal/bootstrap/roles_rollout_file_internal_test.go` | update | Grant-Text der zwei Funktionen (Rolle `cdc_admin`, `PUBLIC` ohne Recht) als **neuer Test** `TestAdministrationDateiTraegtDieFunktionsRechte`, der `tools/schema/nacharbeit-administration.sql` liest: der bestehende Parser der Datei liest die `GRANT`-Zeilen von `nacharbeit-roles.sql` und zerlegt Objektlisten an jedem Komma, das an den Klammern einer Funktionssignatur (`(text, text)`) scheitert; der neue Test zerlegt außerhalb der Klammern und leitet den Umfang aus den `CREATE FUNCTION`-Zeilen ab (jede definierte Funktion in `REVOKE … FROM PUBLIC` und `GRANT EXECUTE … TO cdc_admin`, keine andere Rolle). Die Prüfung (4b) aus `slice-backfill-run-store` hält die Tabellen-Grants von `cdc_admin` auf `cdc.administration_request` (`SELECT`, `UPDATE`; kein `INSERT`, kein `DELETE`; kein Recht der beiden anderen Rollen); die zwei neuen Spalten fallen unter diese Tabellen-Grants und brauchen keinen eigenen Grant. |
| `internal/domain/model/administrationrequest.go` (+ Test) | update | zwei Antragsarten, Felder `RuleName`/`RuleSpec` (`RuleSpec` als JSON-Text), Konstruktor-Invarianten (Signatur trägt zwei Parameter mehr); der Doc-Kommentar zählt die Menge auf. |
| `internal/adapters/driven/postgresstorage/queries/queries.go`, `sqlexec/translate.go` (+ Tests `translate_test.go`, `administrationrequest_test.go`) | update | Lesen der zwei Spalten (`ReadPendingRequests` scannt acht statt sechs Spalten; `COALESCE(rule_name, '')`, `COALESCE(rule_spec::text, '')`), Abbildung der Antragsarten. Die Datei `postgresstorage/administrationrequest.go` selbst bleibt unverändert: der Adapter reicht die Abfrage durch. |
| `internal/bootstrap/wiring.go`, `administration_internal_test.go` | update | der Fehlertext im `default`-Zweig von `applyAdministrationRequest` nennt die verarbeiteten Antragsarten (Konstante `processedAdministrationKinds`) statt einer „geschlossenen Menge“ — die geschlossene Menge der Datenbank hat jetzt sieben Werte, verarbeitet werden fünf; zwei Tests binden Antragsart und Menge an den Fehlertext, einer davon über `processAdministrationRequests` bis zum `failed`-Vermerk. |
| `internal/domain/errors/errors.go`, `internal/application/port/outbound/administrationrequest.go`, `internal/adapters/driven/postgresstorage/tableactivation.go`, `tools/schema/nacharbeit-roles.sql` | update (Kommentare) | die Doc-Kommentare, die die Antragsarten oder Funktionen aufzählen bzw. „die beiden Tabellen-Antragsarten“ als Gegenstück nennen, tragen die zwei neuen Arten (Suchlauf unten). |
| `tools/harness/run-integration-tests.sh`, `harness/README.md` | update (Zahl) | die Zahl der bekannten Fremdobjekte im Kommentar bzw. in der Zeile `make example-demo-up` (sieben → neun). |
| `tools/harness/run-schema-rollout-guard-test.sh`, `harness/targets/schema-rollout.md` | prüfen / update | Zahl der Fremdobjekte und Beschreibung der Läufe; der Alt-Tag-Lauf ([`ADR-0114`](../../adr/0114-schema-rollout-vorlauf-view-signatur.md) Entscheidung 7) wird ausgeführt **und geändert** (Übergabe aus `slice-backfill-sql-administration`): Lauf 5 vergleicht die `request_kind`-Menge des Alt-Bestands nach dem Upgrade exakt mit den **fünf** Werten (`backfill,disable,enable,exclude_column,include_column`, samt Fehlermeldung „fünf Antragsarten") und prüft `EXECUTE` allein für die Funktion `cdc.backfill_table`; der Slice zieht Menge und Meldung auf sieben Werte nach und ergänzt `EXECUTE` (nur `cdc_admin`, nicht `PUBLIC`) je neue Funktion — sonst färbt Lauf 5 mit der Erweiterung rot. **Zusätzlich nötig, am Start gemessen:** der jüngste `v*`-Tag ist `v0.2.0` und trägt die Antragsart `backfill`, die Funktion `cdc.backfill_table` und die Rechte auf `cdc.backfill_run` schon (`git show v0.2.0:tools/schema/nacharbeit-administration.sql` nennt `backfill_table` fünfmal); die drei Vorbedingungen des Alt-Tag-Laufs am Stand vor `backfill` (`cdc_admin` ohne `UPDATE` auf die Antrags-Tabelle, `backfill_table` fehlt, die Menge trägt kein `backfill`) sind gegen diesen Tag falsch — Lauf 5 bräche vor dem Upgrade ab. Die Vorbedingungen wandern auf das Delta dieses Slice (der Tag trägt weder die zwei Spalten noch die zwei Funktionen noch die zwei Antragsarten), eine Antragszeile des Alt-Bestands wird vor dem Upgrade geschrieben und trägt danach NULL in den zwei neuen Spalten, und beide Funktionen werden nach dem Upgrade unter `cdc_admin` (Antrag `pending`) und `cdc_reader` („permission denied for function“) aufgerufen; die Rechte-Prüfungen für `backfill_table` bleiben. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „die geschlossene
`request_kind`-Menge (Parent-Stand plus zwei)“, „die Menge der Fremdobjekte
außerhalb des neutralen Modells“, „die Spaltenform von
`administration_request`“; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Aufzählungen der Antragsarten | `git grep -n exclude_column` über `internal tools docs spec harness Makefile` ohne `docs/reviews`, `done/`, Baseline (Block unten, Zeilen 1–2); Zählwort „fünf“ (Zeilen 3–4); die Aufzählung mit Schrägstrich, die bei `backfill` endet (Zeilen 9–10) | **Gefunden.** Parent: 95 Trefferzeilen, Diff: 96 (Symbolname); Zählwort „fünf Antragsarten/Arten/SQL-Funktionen/Funktionen/Werte“ Parent 13, Diff 7; die Schrägstrich-Aufzählung, die bei `backfill` endet, Parent 6, Diff 5. Code-Träger, die die Menge aufzählen oder „die beiden Tabellen-Antragsarten“ als Gegenstück nennen: `errors.go`, `model/administrationrequest.go`, `queries.go` (zwei Kommentare), `translate.go`, `tableactivation.go`, `outbound/administrationrequest.go`, `wiring.go` (Fehlertext), `nacharbeit-administration.sql`, `nacharbeit-roles.sql`, `schema.yaml` (Kopf und `description`), die Test-Doc-Kommentare in `administration_internal_test.go` und `administrationrequest_test.go`. Spec: `spec/pflichtenheft.md` (sieben Treffer, `SPEC-019` und Änderungshistorie) und `spec/architecture.md` (Tabelle) tragen die sieben Arten bereits (Träger `spec-nachzug`). **Nichtgefunden:** kein weiterer Träger, der die Menge als geschlossene Liste von fünf führt; der Rest der Zählwort-Treffer (Diff 7) sind zwei `Accepted` ADRs mit anderem Gegenstand, zwei Register-Records und die Spec-Historie, dazu der Kommentar in `administration_roles_internal_test.go:58` („einen Antrag jeder der fünf Antragsarten“ — der Login-Test fährt die fünf verarbeiteten Arten) und mein eigener Kommentar am Spaltenform-Test („die fünf Antragsarten ohne Regel“, richtig). | Code-Träger im Diff nachgezogen. Handbuch: `docs/user/benutzerhandbuch.md` trägt zwei Treffer (Beispielaufruf Zeile 265, Änderungshistorie) — an `slice-transformationen-betriebsdoku` gemeldet (deren §3 führt die Zeile „Aufzählungen der Antragsarten im Handbuch“), nicht geändert. `administration_roles_internal_test.go:58` bleibt: die Zeile ist der Login-Test der verarbeiteten Arten, Adresse `slice-transformationen-antragsweg-usecase` (dessen §3 führt „der Login-Test zieht je Antragsart“). `Accepted` ADRs (0059, 0065, 0111, 0112) und Records nicht geändert. |
| Zahl der Fremdobjekte (Wortform des Parent-Stands, am Start gemessen: „sieben“) | `git grep -n` mit den Wortformen „sieben bekannten/Objekte/Fremdobjekte“ (Zeilen 5–6) und „neun bekannten/Objekte“ (Zeilen 7–8) über `internal tools docs spec harness Makefile` ohne `docs/reviews`, `done/`, Baseline | **Gefunden.** Parent: 12 Trefferzeilen „sieben …“ (Guard-Kommentar, `guard_test.go`, `harness/targets/schema-rollout.md`, `harness/README.md`, `run-integration-tests.sh`, `run-schema-rollout-guard-test.sh`), Diff: 0; „neun …“ Parent 0, Diff 12. **Nichtgefunden:** keine Wortform „acht“ oder „sechs“ als Zahl der Fremdobjekte in einem Träger; `Accepted` ADRs nennen ältere Zahlen (`ADR-0111`: „die sechs Fremdobjekte“, `ADR-0064`: „vier“) — Aussage-Text zu ihrem Stand, unberührbar. | alle zwölf Träger nachgezogen; ADRs und Register-Records unverändert. |
| Läufe des Guard-Tests | Zählwort „sechs Läufe“ (Zeilen 19–20) und Lesen von `tools/harness/run-schema-rollout-guard-test.sh` und `harness/targets/schema-rollout.md` §Belege | **Gefunden.** „sechs Läufe“: Parent 6, Diff 6 — die Zahl der Läufe ändert sich nicht (der Alt-Tag-Lauf ist Lauf 5 und bleibt einer). Geändert hat sich der Inhalt von Lauf 5 (Vorbedingungen, Prüfungen), den Kopfkommentar des Skripts und `harness/targets/schema-rollout.md` §Belege beschreiben. **Nichtgefunden:** kein Träger, der Lauf 5 mit „fünf Werten“ oder „`backfill_table`“ allein beschreibt. | Skript-Kopf, `schema-rollout.md` §Belege nachgezogen. |
| Spaltenform von `administration_request` in weiteren Trägern | `git grep -n administration_request` über `internal tools` ohne `tools/schema/plan.yaml`/`down.sql` (Zeilen 11–12); dasselbe in `internal/adapters/driven/postgresstorage/schema.sql` (Zeilen 13–14) | **Gefunden.** Parent 137, Diff 152 (Tests und Kommentare; die Spaltenform selbst steht allein in `tools/schema/schema.yaml`). **Nichtgefunden:** `internal/adapters/driven/postgresstorage/schema.sql` trägt `administration_request` weder am Parent noch im Diff (0/0) — es gibt keinen zweiten Träger der Spaltenform, und keine Test-Fixture legt die Tabelle mit eigener DDL an (`git grep -n 'CREATE TABLE.*administration_request'` trifft allein `tools/schema/plan.yaml`, das Erzeugnis des Rollouts). | keine Änderung an `schema.sql`; `plan.yaml`/`down.sql` regeneriert (§3 oben). |
| Fehlertext der geschlossenen Menge im `default`-Zweig | `git grep -n 'geschlossene Menge' -- internal/bootstrap/wiring.go` (Zeilen 15–16) und „verarbeiteten Antragsarten“ (Zeilen 17–18) | **Gefunden.** „geschlossene Menge“ Parent 1 (der Fehlertext), Diff 0; „verarbeiteten Antragsarten“ Parent 0, Diff 2 (Fehlertext, Konstante). **Nichtgefunden:** kein zweiter Träger dieses Wortlauts in `internal/bootstrap`. | der Fehlertext nennt die verarbeiteten Antragsarten (fünf), nicht die geschlossene Menge der Datenbank (sieben); zwei Tests lesen ihn (Antragsart und Menge), ein dritter geht bis zum `failed`-Vermerk. |
| Handbuch und Spec (Meldung) | `git grep -n` mit `exclude_column`/`enable_table`/`backfill_table` in `docs/user` (Zeilen 21–22); `set_transformation` in `spec` (Zeilen 23–24) | **Gefunden.** `docs/user`: Parent 11, Diff 11 — kein Diff-Eintrag, Handbuch unberührt; `spec`: Parent 14, Diff 14 — die Spec trägt die zwei Antragsarten seit `spec-nachzug`. **Nichtgefunden:** keine Handbuch-Stelle, die `set_transformation` oder `remove_transformation` nennt (die Betreiber-Oberfläche ist unbeschrieben, bis `betriebsdoku` sie beschreibt). | gemeldet an `slice-transformationen-betriebsdoku`; Kandidatenlauf (`git diff --name-only 2f5ed3dc -- internal/bootstrap/ tools/schema/ internal/adapters/driving/`) trifft die neue Oberfläche, `docs/user/benutzerhandbuch.md` liegt bewusst nicht im Diff (Aufschub mit Adresse, §2). |

```suchlauf
2f5ed3dc 95 -n exclude_column -- internal tools docs spec harness Makefile :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
diff 96 -n exclude_column -- internal tools docs spec harness Makefile :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
2f5ed3dc 13 -n -e 'fünf Antragsarten' -e 'fünf Arten' -e 'fünf SQL-Funktionen' -e 'fünf Funktionen' -e 'fünf Werte' -- internal tools docs spec harness Makefile :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
diff 7 -n -e 'fünf Antragsarten' -e 'fünf Arten' -e 'fünf SQL-Funktionen' -e 'fünf Funktionen' -e 'fünf Werte' -- internal tools docs spec harness Makefile :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
2f5ed3dc 12 -n -e 'sieben bekannten' -e 'sieben Objekte' -e 'sieben Fremdobjekte' -- internal tools docs spec harness Makefile :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
diff 0 -n -e 'sieben bekannten' -e 'sieben Objekte' -e 'sieben Fremdobjekte' -- internal tools docs spec harness Makefile :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
2f5ed3dc 0 -n -e 'neun bekannten' -e 'neun Objekte' -- internal tools docs spec harness Makefile :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
diff 12 -n -e 'neun bekannten' -e 'neun Objekte' -- internal tools docs spec harness Makefile :!docs/reviews :!docs/plan/planning/done :!.harness/baseline
2f5ed3dc 6 -n -e 'include_column/backfill' -e 'include_column`/`backfill' -e 'include_column`, `backfill' -- internal tools docs spec harness Makefile :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/planning/observations
diff 5 -n -e 'include_column/backfill' -e 'include_column`/`backfill' -e 'include_column`, `backfill' -- internal tools docs spec harness Makefile :!docs/reviews :!docs/plan/planning/done :!.harness/baseline :!docs/plan/planning/observations
2f5ed3dc 137 -n administration_request -- internal tools :!tools/schema/plan.yaml :!tools/schema/down.sql
diff 152 -n administration_request -- internal tools :!tools/schema/plan.yaml :!tools/schema/down.sql
2f5ed3dc 0 -n administration_request -- internal/adapters/driven/postgresstorage/schema.sql
diff 0 -n administration_request -- internal/adapters/driven/postgresstorage/schema.sql
2f5ed3dc 1 -n 'geschlossene Menge' -- internal/bootstrap/wiring.go
diff 0 -n 'geschlossene Menge' -- internal/bootstrap/wiring.go
2f5ed3dc 0 -n 'verarbeiteten Antragsarten' -- internal/bootstrap/wiring.go
diff 2 -n 'verarbeiteten Antragsarten' -- internal/bootstrap/wiring.go
2f5ed3dc 6 -n 'sechs Läufe' -- harness tools/harness docs/user
diff 6 -n 'sechs Läufe' -- harness tools/harness docs/user
2f5ed3dc 11 -n -e exclude_column -e enable_table -e backfill_table -- docs/user
diff 11 -n -e exclude_column -e enable_table -e backfill_table -- docs/user
2f5ed3dc 14 -n -e set_transformation -- spec
diff 14 -n -e set_transformation -- spec
```

**Mutationen des Implementer-Laufs** (Zusage · mutierte Eingabe · gesehenes Rot; jede Mutation
am Arbeitsbaum, danach byte-gleich zurückgenommen, geprüft mit `cmp`):

| Zusage | mutierte Eingabe | gesehenes Rot |
|---|---|---|
| Konstruktor verlangt `rule_name` für `remove_transformation` | Bedingung `ruleName == ""` im Zweig durch `false` ersetzt | `TestNewAdministrationRequestRejectsInvariantViolations/Transformations-Antragsart_ohne_Regelname` (`make test`-Container, `-run`) |
| Konstruktor verlangt `rule_spec` für `set_transformation` | Teilbedingung `ruleSpec == ""` entfernt | `…/set_transformation_ohne_Regelform` |
| die Antragsart `set_transformation` gehört zur Menge | den Zweig entfernt | drei Teiltests, „unbekannte Antragsart“ statt `ErrEmptyIdentifier` |
| `ListPending`-Übersetzung trägt den Regelnamen | Argument `ruleName` durch `column` ersetzt (`translate.go`) | `TestReadPendingRequestsTranslatesRequests` („leere Kennung“) |
| die Abfrage liest `rule_spec` | `COALESCE(rule_spec::text, '')` durch `''` ersetzt (`queries.go`) | `TestAdministrationRequestTransformationRequestsCarryRuleAndNotify` (`ListPending: leere Kennung`; `make test-store` Exit 2) |
| `cdc.set_transformation` schreibt die Art `set_transformation` | Art in der Funktion durch `'remove_transformation'` ersetzt | derselbe Test (Antrags-Zeile trägt die falsche Art; `make test-store` Exit 2) |
| die `request_kind`-Menge trägt sieben Werte | `'set_transformation'` aus dem CHECK gestrichen | `TestAdministrationRequestKindCheckCarriesExactlyTheSevenKinds` (SQLSTATE 23514; `make test-store` Exit 2) |
| `rule_spec` ist `jsonb` | `rule_spec: { type: text }` in `schema.yaml` | `TestAdministrationRequestRuleColumnsAreNullable` (`rule_spec:text:YES`; `make test-store` Exit 2) |
| der Fehlertext nennt die Antragsart | `request.Kind` durch ein Literal ersetzt (`wiring.go`) | `TestApplyAdministrationRequestRejectsUnprocessedKind`, `TestProcessAdministrationRequestsMarksTransformationRequestsFailed` |
| der Fehlertext nennt die verarbeiteten Antragsarten | `/backfill` aus `processedAdministrationKinds` gestrichen | `TestApplyAdministrationRequestRejectsUnprocessedKind` (nennt `backfill` nicht) |
| der Guard trägt die gemessene Schreibweise | `in:json` durch `in:jsonb` ersetzt | drei Tests in `tools/schema/rolloutguard` (`…AllowsDestructiveWhenAllBlockersKnown`, `…ViewSignatureWithKnownForeignObjects`, `…RefusesOtherSignatureSpelling`) |
| der Guard trägt beide Funktionen | den Eintrag `remove_transformation` gestrichen | zwei Tests (`…AllowsDestructiveWhenAllBlockersKnown`, `…ViewSignatureWithKnownForeignObjects`) |
| `cdc_admin` darf `set_transformation` aufrufen | Signatur aus der `GRANT`-Liste gestrichen | `TestAdministrationDateiTraegtDieFunktionsRechte` („fehlt in einer `GRANT EXECUTE … TO cdc_admin`-Anweisung“) |
| `PUBLIC` hat kein `EXECUTE` | Signatur aus der `REVOKE`-Liste gestrichen | derselbe Test („fehlt in einer `REVOKE EXECUTE … FROM PUBLIC`-Anweisung“); Alt-Tag-Lauf: `cdc_capture` trägt über `PUBLIC` `EXECUTE` (`run-schema-rollout-guard-test.sh` Exit 1) |
| nur `cdc_admin` trägt `EXECUTE` | eine zweite `GRANT EXECUTE … TO cdc_reader`-Zeile angehängt | derselbe Test („cdc_reader trägt EXECUTE auf …“) |
| der Alt-Tag-Lauf belegt die sieben Antragsarten | `'set_transformation'` aus dem CHECK gestrichen | Lauf 5 („trägt … die Menge …remove_transformation statt der sieben Antragsarten“, Exit 1) |
| der Alt-Tag-Lauf belegt die zwei Spalten | die Spalte `rule_spec` aus `schema.yaml` gestrichen | Lauf 5 („… die Spalten 'rule_name:text:YES' statt …“, Exit 1) |
| der zweite Rollout endet mit Exit 0 | Parametertyp `jsonb` statt `json` in der Funktion | zweiter `make schema-rollout` Exit 2 mit d-migrate-Fehler 5 (`does not exist`), mit `json` Exit 0 (gemessen; nicht mit dem Guard-Skript wiederholt) |
| die Ausnahme in `.dockerignore` trägt den Leser | die Negation nicht gesetzt | `make coverage-gate` Exit 2 (`open /src/tools/schema/nacharbeit-administration.sql: no such file or directory`) |

Leer bleibt keine Zusage. Nicht mutiert (Grund): die reinen Kommentar-, Zahl- und Beschreibungs-Nachzüge
(Suchlauf oben) — sie tragen keine Zusage, die ein Test bindet.

**Belege des Implementer-Laufs** (Exit-Codes je Lauf ungefiltert in eine Log-Datei geschrieben und
gesondert gelesen, [`AGENTS.md`](../../../../AGENTS.md) §3.9; gedruckte Zeilen wörtlich):

- `make test` (Race-Detektor) Exit 0; `make test-store` Exit 0, gedruckt: `db-coverage: OK —
  DB-Adapter-Coverage 82.61% erfuellt Schwelle 80%`; `make a-check` Exit 0, gedruckt: `gesamt: 0
  Befund(e)`; `make coverage-gate` Exit 0, gedruckt: `coverage-gate: OK — Coverage 83.90% erfüllt
  Schwelle 80%`.
- Zwei aufeinanderfolgende `make schema-rollout` gegen eine frische Wegwerf-Datenbank: Exit 0 und
  Exit 0; `plan.yaml` und `down.sql` stammen aus dem ersten Lauf.
- `bash tools/harness/run-schema-rollout-guard-test.sh` Exit 0, alle sechs Läufe; gedruckt: `Lauf 5 OK
  — Tag v0.2.0: Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, ohne Vorlauf), Exit 0 (Arbeitsbaum,
  zweiter Lauf); Zeile alttag-ch über cdc.changes lesbar; 19 Tabellen-/View-Rechte der drei Rollen
  auf administration_request, backfill_run und backfill_status, EXECUTE auf 3 Funktionen … allein
  für cdc_admin (nicht PUBLIC), Spalten rule_name:text:YES,rule_spec:jsonb:YES (Alt-Zeile
  alttag-req: NULL), request_kind-Menge
  backfill,disable,enable,exclude_column,include_column,remove_transformation,set_transformation,
  cdc.set_transformation/cdc.remove_transformation unter cdc_admin schreiben pending-Anträge,
  cdc_reader: permission denied for function`; Schlusszeile: `OK — alle Belege real erbracht
  (Idempotenz-Allow, echte Änderung bleibt wirksam, View-Signatur-Vorlauf, Alt-Tag v0.2.0,
  Negativ-Abbruch: unbekannte Funktion bleibt bestehen, make-Exit 2/2 mit d-migrate-Exit 8)`.
- `make suchlauf-nachmessen PLAN=<diese Datei>` Exit 0, gedruckt: `suchlauf-nachmessen: 24 Zeilen
  stimmen`.
- `make gates` (ohne Pipe in eine Log-Datei, Exit gesondert gelesen): der erste Lauf Exit 2 —
  `docs-check` fand die unverlinkte Kennung im Suchlauf-Feld (`id-unlinked`), im Feld korrigiert; der
  zweite Lauf Exit 0, gedruckt: `d-check: 1216 Datei(en) geprüft, 0 Befund(e)`, `commit-traceability:
  OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID`, `generated-sync: OK`,
  `coverage-gate: OK — Coverage 83.90% erfüllt Schwelle 80%`, `gesamt: 0 Befund(e)` (`a-check`).
- Image-Neutralität der `.dockerignore`-Zeile: siehe §6 (Digest vor und nach der Zeile gleich).

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-backfill-sql-administration` in
`done/` liegt (Kopplung K3 der Welle
[welle-backfill-bestand](../done/welle-backfill-bestand.md) §5: beide Umsetzungen
berühren die geschlossene `request_kind`-Menge in
`tools/schema/nacharbeit-administration.sql`,
[`SPEC-019`](../../../../spec/pflichtenheft.md), `applyAdministrationRequest`
und den Idempotenz-Guard — dieser Slice erweitert die dort stehende Menge
additiv), `slice-transformationen-kern-rename` in `done/` liegt und kein
anderer Slice in `in-progress/` liegt (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Schema, Guard,
  Rollen-Test, Domäne und Store nicht in einem Review tragen — der abtrennbare
  Teil ist der dritte Liefer-Punkt (Domäne und Store-Lesen) als eigener Slice
  `antragsweg-store` mit Start nach diesem; die Welle nimmt ihn dann in §4 auf.
- `in-progress` → `open` (blockiert): falls das neutrale Modell die zwei
  nullable Spalten samt CHECK-Änderung nicht konvergiert (dann Nacharbeit-SQL
  und ein zweiter Guard-Eintrag — Architect-Frage nach
  [`ADR-0043`](../../adr/0043-schemamigrationen-mit-d-migrate.md)) oder der
  Idempotenz-Guard die Funktion mit `jsonb`-Parameter nicht erkennt (Exit 8).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test`, `make test-store` und `make
schema-rollout` (zweimal) real grün + der Alt-Tag-Lauf von
`tools/harness/run-schema-rollout-guard-test.sh` real grün + Closure-Notiz mit
Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Der Idempotenz-Guard erkennt die Funktion mit `jsonb`-Parameter nicht** —
  ein zweiter Rollout bräche mit Exit 8. Die Signatur-Schreibweise der
  bestehenden Funktionen (`enable_table(in:text,in:text,in:text)` in
  `guard.go`) ist gelesen, die der neuen mit `jsonb` ist ungemessen. *Erwartet,
  zu belegen durch:* zweiter `make schema-rollout` und
  `run-schema-rollout-guard-test.sh`. *Beleg des Implementer-Laufs:* der Guard
  erkennt die Funktion — der `--plan-only`-Report gegen ein migriertes Ziel
  nennt sie als `set_transformation(in:text,in:text,in:text,in:text,in:json)`,
  und `guard.go` trägt genau diese Schreibweise. Das Risiko trat trotzdem in
  anderer Form ein: mit einem `jsonb`-Parameter endete der zweite Rollout nicht
  bei der Wache (Exit 8), sondern in `--execute` mit d-migrate-Exit 5 — d-migrate
  rendert den Abbau als `DROP FUNCTION "set_transformation"(text, text, text,
  text, json)`, die Signatur existiert nicht (siehe den folgenden Punkt).
  **Ausgang:** *(bei Closure)*
- **Der Parameter `rule_spec` ist `json`, nicht `jsonb` — Abweichung vom
  Wortlaut von `ADR-0112` Teilfrage 1 und `SPEC-019`.** Gemessen (zwei
  aufeinanderfolgende `make schema-rollout` gegen eine frische Wegwerf-
  Datenbank): mit `rule_spec jsonb` Exit 0, dann Exit 5
  (`executionError: function set_transformation(text, text, text, text, json)
  does not exist`, Wiederherstellung `FULL_ROLLBACK_CONFIRMED`); mit `rule_spec
  json` Exit 0 und Exit 0. Folge für Aufrufer: ein Literal und ein `::json`-Wert
  werden angenommen, ein `::jsonb`-Wert nicht (gemessen, „function … does not
  exist“). Der Slice hält den kleinsten Stand, der den Rollout idempotent lässt;
  die Wahl zwischen dem Parametertyp `json`, einem Erratum von `ADR-0112`/
  `SPEC-019` (ein Text, der den Typ nennt) oder einer anderen Behandlung der
  Fremdobjekte in `make schema-rollout` ist eine Architect-Entscheidung
  (`ADR-0043` Re-Evaluierungs-Trigger; angenommen ist die Nacharbeit-Form, nicht
  ihre Grenze bei `jsonb`). Ungeklärt, ein Versuch: `--routine-capability
  'function:enabled=false'` ließ `--plan-only` gegen ein Ziel mit einer
  `jsonb`-Funktion weiter mit Exit 8 enden (der Report wurde nicht erfasst,
  Ausgabepfad außerhalb des Mounts). *Erwartet, zu belegen durch:* Architect-
  Verdikt; der Betreiber-Text von `slice-transformationen-betriebsdoku` nennt
  die Aufrufform (Literal oder `::json`). **Ausgang:** *(bei Closure: weiter
  offen — Adresse Architect)*
- **Die zwei Spalten, besonders `rule_spec jsonb`, konvergieren nicht über
  einen Alt-Bestand.** Eine nullable `text`-Spalte an einer bestehenden Tabelle
  rollt über einen Alt-Bestand (gemessen im Architect-Verdikt
  `architect-verdict-schema-rollout-view-signatur`, Szenario 5); `jsonb` statt
  `text` ist nicht gemessen. *Erwartet, zu belegen durch:* der Alt-Tag-Lauf von
  `run-schema-rollout-guard-test.sh`. *Beleg des Implementer-Laufs:* Alt-Tag-Lauf
  gegen `v0.2.0` Exit 0 (Rollout des Tags), Exit 0 (Arbeitsbaum, ohne Vorlauf),
  Exit 0 (zweiter Lauf); danach `rule_name:text:YES,rule_spec:jsonb:YES`, die
  Alt-Zeile `alttag-req` trägt NULL in beiden. **Ausgang:** *(bei Closure)*
- **Zwei Schema-Träger driften** (`tools/schema/schema.yaml` und
  `postgresstorage/schema.sql` tragen dieselbe Spaltenform, jeder für einen
  anderen Tier). *Erwartet, zu belegen durch:* der Suchlauf-Eintrag und `make
  test-store`, das die Spalten real liest. *Beleg des Implementer-Laufs:* es
  gibt keinen zweiten Träger — `schema.sql` führt `cdc.administration_request`
  nicht (Suchlauf-Zeilen 13–14, 0/0), die Annahme des Plans war falsch; `make
  test-store` liest die zwei Spalten über den Rollout. **Ausgang:** *(bei
  Closure: entfallen)*
- **Der Rollen-Test deckt den neuen Grant nicht.** Der Parser liest literale
  `GRANT … ;`-Anweisungen (Kopf von `roles_rollout_file_internal_test.go`);
  eine Funktion mit `jsonb`-Parameter ist gegen ihn ungemessen
  (`BEO-PGC/rollen-test-abdeckungsluecken`, offen, 2×). *Erwartet, zu belegen
  durch:* die Mutation aus DoD Punkt 2. *Beleg des Implementer-Laufs:* der
  bestehende Parser trifft die Funktions-Grants nicht (Kommas in der
  Signatur, andere Datei); der neue Test liest `nacharbeit-administration.sql`
  und ist gegen drei Mutationen rot gesehen (Grant-Liste ohne
  `set_transformation`, `REVOKE`-Liste ohne `set_transformation`, zusätzlicher
  `GRANT … TO cdc_reader`). **Ausgang:** *(bei Closure)*
- **Zwischenzustand: Funktion vorhanden, Wirkung fehlt.** Ein Betreiber, der
  `cdc.set_transformation` zwischen diesem Slice und `antragsweg-usecase`
  aufruft, erhält einen `failed`-Antrag statt einer Regel; der Fehlertext nennt
  die Ursache. Der Slice ist unveröffentlicht und die Folge-Slices folgen
  unmittelbar. *Erwartet, zu belegen durch:* der Test des `default`-Zweigs.
  *Beleg des Implementer-Laufs:* `TestApplyAdministrationRequestRejectsUnprocessedKind`
  und `TestProcessAdministrationRequestsMarksTransformationRequestsFailed`
  (Antragsart und Menge im Fehlertext, `failed`-Vermerk statt `applied`),
  je an der Eingabe rot gesehen (Tabelle „Mutationen“ in §3). **Ausgang:** *(bei Closure: entfallen
  mit der Closure von `slice-transformationen-antragsweg-usecase`)*
- **Das Handbuch nennt die Funktionen nicht**
  (`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`,
  verkörpert, 3×): die Adresse `slice-transformationen-betriebsdoku` muss den
  aufgeschobenen Gegenstand tragen. *Erwartet, zu belegen durch:* Review liest
  §2 von `betriebsdoku` gegen §1 dieses Slice. *Beleg des Implementer-Laufs:*
  §2 von `betriebsdoku` nennt SQL-Funktionen, Rolle `cdc_admin`, Status und
  Fehlertexte, Wirkung, Abhilfe, Glossar und Änderungshistorie und ist damit
  gegen §1 dieses Slice vollständig; zwei Sachverhalte dieses Slice stehen dort
  nicht und gehören in den Abschnitt „Transformationsregel konfigurieren“: die
  Aufrufform des Parameters `rule_spec` (Literal oder `::json`, ein
  `::jsonb`-Wert wird abgelehnt) und dass ein Antrag mit NULL oder leerem
  `rule_name`/`rule_spec` in dieser Fassung die Lesung der Antrags-Queue stört
  (nächster Punkt; Adresse `antragsweg-usecase`, nach dessen Closure gilt
  `SPEC-019`). Beides ist als Meldung im Bericht an den Planner genannt, der
  Plan von `betriebsdoku` bleibt unverändert. **Ausgang:** *(bei Closure)*
- **Idempotenz gegen die Fremdobjekte.** Jede neue `nacharbeit-*`-Funktion
  vergrößert die Menge, die ein zweiter Rollout als Blocker sieht
  (`BEO-PGC/schema-rollout-fremdobjekte`, verkörpert, 3×;
  `BEO-PGC/d-migrate-nacharbeit`, verkörpert, 6×). *Erwartet, zu belegen
  durch:* zweiter Rollout, `run-schema-rollout-guard-test.sh`. *Beleg des
  Implementer-Laufs:* zwei aufeinanderfolgende `make schema-rollout` gegen eine
  frische Wegwerf-Datenbank Exit 0 und Exit 0; `run-schema-rollout-guard-test.sh`
  Exit 0 mit sechs Läufen (Lauf 2 nimmt den `--allow-destructive`-Pfad über neun
  bekannte Blocker). **Ausgang:** *(bei Closure)*
- **Ein Antrag mit NULL oder leerem `rule_name`/`rule_spec` stört die Lesung
  der Antrags-Queue.** Die SQL-Funktionen prüfen nichts
  ([`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md)),
  `cdc.set_transformation(…, NULL, NULL)` schreibt die Zeile (gemessen: die
  Zeile mit NULL in `rule_spec` entsteht). Der Antrags-Konstruktor lehnt sie ab
  (DoD Punkt 3, Invariante wie bei `column_name`); `ReadPendingRequests` gibt
  den Fehler unverändert zurück, `processAdministrationRequests` protokolliert
  ihn und versucht es im nächsten Durchlauf erneut — kein Antrag der Queue wird
  verarbeitet, bis die Zeile entfernt ist (hergeleitet aus dem Quelltext; am
  Konstruktor-Pfad gemessen: `ListPending: leere Kennung`). Dieselbe Klasse
  besteht am Parent-Stand für `cdc.exclude_column(…, NULL)`. `SPEC-019` verlangt
  für diese Fälle `failed` mit den Texten `Regelname ist ungültig`/`rule_spec
  ist ungültig`; die Stelle der Prüfung legt `antragsweg-usecase` fest (dessen
  §2, „Konstruktor-Sentinel oder Vorprüfung im Use Case“). *Erwartet, zu
  belegen durch:* dessen Test je Auslöser. **Ausgang:** *(bei Closure: weiter
  offen — Adresse `slice-transformationen-antragsweg-usecase`)*
- **Der Alt-Tag-Lauf koppelt seine Vorbedingung an den jüngsten Tag.** Lauf 5
  prüft vor dem Upgrade, dass der Stand des Tags weder die zwei Spalten noch die
  zwei Funktionen noch die zwei Antragsarten trägt; ein späterer Release-Tag mit
  diesen Objekten lässt die Vorbedingung laut scheitern („Vorbedingung
  fehlgeschlagen“, kein stilles Grün). Dieselbe Kopplung trug der Lauf für die
  Antragsart `backfill`; sie zerbrach mit `v0.2.0` (Vorbedingungen falsch).
  *Erwartet, zu belegen durch:* der nächste Release, der das Delta neu
  bestimmt. **Ausgang:** *(bei Closure: weiter offen — Adresse Planner, der
  Release-Slice bestimmt das Delta des neuen Tags)*
- **Die `.dockerignore`-Ausnahme verändert das Image.** Sie verlangt
  [`ADR-0085`](../../adr/0085-build-kontext-ausnahme-test-only-zweck.md)
  Merkmal (iv): das Image bleibt unberührt, mit Beleg. *Beleg des
  Implementer-Laufs:* `make image` vor und nach der Zeile am selben Quellstand,
  Digest beide Male `sha256:b9cfd65527f004eead067ab4bbc7102b19f9236b1eda6e03ded2d74557f39fe4`
  ([`ADR-0044`](../../adr/0044-image-beleg-semantik.md): Digest innerhalb
  desselben Builders; der Digest-Commit entfällt bei unverändertem Digest,
  `harness/image-hash.txt` ist lokal und nicht committet).
  **Ausgang:** *(bei Closure: entfallen)*

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
`*`/`PGC` (Greenfield); Schema-Skripte, Store-Adapter und Composition Root sind
keine eigenen Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/schema-rollout-fremdobjekte` (verkörpert, 3×) und
`BEO-PGC/d-migrate-nacharbeit` (verkörpert, 6×) — Guard-Eintrag und Rückführung
§4, `BEO-PGC/rollen-verdrahtung` (eingetreten, Grants der Rolle `cdc_admin`),
`BEO-PGC/rollen-test-abdeckungsluecken` (offen, 2×, einschlägig — Risiko §6),
`BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an` (offen, 2×, einschlägig — der
Handbuch-Aufschub),
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` und
`BEO-PGC/handbuch-versionshistorie-uebersprungen` (verkörpert, je 3×, Aufschub
mit Adresse), `BEO-PGC/test-isolation-geteilter-zustand` (offen, 1×,
einschlägig — Store-Tests skopiert bereinigen),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf §3).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
