# Slice routing-antragsweg: Antragsweg und Dauerhaftigkeit — `cdc.set_route`/`cdc.remove_route`, Validierung R1–R6, Regelstand aus den `applied`-Zeilen

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-routing](../welle-routing.md).

**Bezug:** [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (Routing —
Haupt-Bezug), [`LH-FA-ADM-001`](../../../../spec/lastenheft.md) (SQL-Administration),
[`LH-QA-SEC-002`](../../../../spec/lastenheft.md) (getrennte Berechtigbarkeit),
[`LH-QA-SEC-004`](../../../../spec/lastenheft.md) (ausgeschlossene Spaltenwerte
erscheinen nicht in den Changes),
[`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
Folgepflicht 3 (Antragsweg und Dauerhaftigkeit — Träger dieses Slice),
[`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)
(Antrags-Queue), [`ADR-0065`](../../adr/0065-spaltenausschluss-dauerhafter-traeger.md)
(Dauerhaftigkeit aus den `applied`-Zeilen),
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
(Formvorbild der zwei Antragsarten),
[`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
(Validierung in Go).

**Berührte Spec-Stellen:**
[`SPEC-019`](../../../../spec/pflichtenheft.md) (Antrags-Datensatz, Fehlertexte),
die Routing-Regelform (neue Kennung aus `slice-routing-spec-nachzug`),
[`LH-FA-CFG-008.a`](../../../../spec/pflichtenheft.md),
[`ARC-005`](../../../../spec/architecture.md) (SQL-Funktionen). Die Spec führt:
der Slice setzt `slice-routing-spec-nachzug` voraus und ändert sie nicht.

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-routing](../welle-routing.md).
**Datum:** 2026-10-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Betreiber richtet eine Routing-Regel über die SQL-Antrags-Queue ein
und entfernt sie wieder; eine verletzte Invariante R1–R6 endet `failed` mit dem
Klartext der Spec, der Regelstand bleibt; der Regelstand ist dauerhaft aus den
`applied`-Zeilen abgeleitet und wirkt am laufenden `Assembler` ohne Neustart. Drei
Liefer-Punkte:

- (A) **Schema und Funktionen:** die `request_kind`-Menge um `set_route` und
  `remove_route` (neun Werte), die SQL-Funktionen `cdc.set_route(source_id,
  schema_name, table_name, rule_name, rule_spec)` und `cdc.remove_route(source_id,
  schema_name, table_name, rule_name)`, ausführbar nur für `cdc_admin`
  ([`LH-QA-SEC-002`](../../../../spec/lastenheft.md)), im Idempotenz-Guard und im
  Rollen-Test; die Spalten `rule_name`/`rule_spec` bestehen (Transformationen);
- (B) **Use Cases:** `SetRoute`/`RemoveRoute` mit der Validierung der `rule_spec` in
  Go (Form, unbekannte Schlüssel enden `failed`), den Invarianten R1 (Regelname
  eindeutig, eigener Namensraum), R2 (`order` eindeutig), R3 (`when.column`
  existiert und ist nicht ausgeschlossen — **und** umgekehrt: `exclude_column` gegen
  eine Spalte, die eine Routing-Bedingung trägt, endet `failed`), R4 (höchstens
  eine Regel ohne `when`, mit der höchsten `order`), R5 (Paar (`column`, `equals`)
  höchstens einmal) und R6 (`remove_route` gegen einen nicht geführten Namen);
  Regelstand-Port und Store-Adapter;
- (C) **Dauerhaftigkeit und Verdrahtung:** der Regelstand einer Tabelle wird aus den
  `applied`-Zeilen beider Antragsarten abgeleitet (Ordnung `requested_at`, dann
  `administration_request_id`), in jedem Pfad, der eine Bindung anlegt
  (`applyAdministrationRequest`, Aktivierungs-Zweig, Prozessstart), und am
  laufenden `Assembler` gesetzt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Auswertung selbst** — `slice-routing-kern-label` (Domäne, `Assembler`,
  Spalte); dieser Slice setzt ihre Schnittstelle (Setzen der Regelliste an der
  Bindung) voraus und ändert sie nicht.
- **Der Backfill-Pfad** — `slice-routing-backfill-pfad`; zwischen diesem Slice und
  jenem sind Regeln setzbar, ein Backfill-Run liefert `route_target = NULL`
  (benanntes Fenster, Welle §5, ein Slice lang).
- **Lesewege und NATS-Subjekt** — `slice-routing-lesewege`, `slice-routing-nats-subjekt`.
- **Handbuch-Beschreibung der Betreiber-Oberfläche** — `slice-routing-betriebsdoku`
  (Aufschub-Adresse, §2): die Wirkung liegt fünf Slices hinter der Oberfläche, und
  das Handbuch darf nur beschreiben, was `slice-routing-e2e` belegt hat.
- **gRPC-Verwaltungs-API für die Antragsarten** — der `Administration`-Service
  trägt keinen RPC für Regeln (gemessen am Parent `30fd6cb5`: `grep -n 'rpc '
  proto/cdc/administration/v1/administration.proto` druckt elf RPCs, darunter
  `EnableTable`/`DisableTable`, keinen für Transformations- oder Spaltenregeln); ein
  RPC für Routing wäre eine eigene Entscheidung.
- **Operatoren über Gleichheit hinaus, mehrere Ziele, Standardziel** — mit der ADR
  ausgeschlossen.

## 2. Definition of Done

- [ ] [`LH-FA-ADM-001`](../../../../spec/lastenheft.md) und
      [`LH-QA-SEC-002`](../../../../spec/lastenheft.md) (A): `cdc.set_route` und
      `cdc.remove_route` legen unter `cdc_admin` einen `pending`-Antrag der Art
      `set_route`/`remove_route` an; ein Aufruf unter `cdc_reader` und
      `cdc_capture` endet mit SQLSTATE `42501`; die `request_kind`-Menge trägt neun
      Werte; der Rollout zweimal hintereinander endet mit Exit 0; der Alt-Tag-Lauf
      (`tools/harness/run-schema-rollout-guard-test.sh`) kennt die neuen Arten und
      Funktionen (seine hart verdrahtete Liste der Arten ändert der Slice mit).
      *Zu belegen durch:* `make test-store` (Rollen-Test, Funktionen),
      `make schema-rollout` (zweimal), Alt-Tag-Lauf.
- [ ] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) Boundary (B): je eine
      Verletzung von R1 bis R6 endet `failed` mit dem Klartext der Spec und der
      Adresse, der Regelstand bleibt, die Gegenprobe ohne die Verletzung endet
      `applied`; ein unbekannter Schlüssel und ein ungültiger Zielname enden
      `failed`; [`LH-QA-SEC-004`](../../../../spec/lastenheft.md): R3 sperrt in
      **beide** Richtungen (Regel gegen ausgeschlossene Spalte; `exclude_column`
      gegen Routing-Spalte), der bestehende Pfad von `exclude_column` gegen eine
      Spalte ohne Routing-Bedingung bleibt unverändert (Regressionstest).
      *Zu belegen durch:* Use-Case-Tabellentests je Verletzung mit Eingabe-Bindung
      (die Eingabe ist die Verletzung, nicht ein Nachbarfall,
      `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`), `make test`,
      `make test-store`.
- [ ] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) Happy Path (C): ein
      `applied`-Antrag setzt die Regelliste am laufenden `Assembler` ohne Neustart;
      nach einem Prozessstart ist der Regelstand aus den `applied`-Zeilen abgeleitet
      (Ordnung `requested_at`, dann `administration_request_id`, auch bei gleichem
      Zeitpunkt), und derselbe Stand entsteht im Aktivierungs-Zweig einer neu
      aktivierten Tabelle; ein Antrag gegen eine Tabelle ohne laufende Bindung endet
      wie bei den Spalten- und Transformations-Antragsarten (Spec-Zusage).
      *Zu belegen durch:* `internal/bootstrap`-Test mit realen Store-Zeilen
      (`make test-store`), Unit-Test der Ableitung (`make test`).
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-routing-antragsweg.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: `harness/targets/schema-rollout.md` (Zählwort der SQL-Funktionen
      und der Arten) und `harness/README.md` (nur bewegte Beschreibungen, laut
      Suchlauf). **Aufschub mit Adresse:** die Beschreibung der Betreiber-Oberfläche
      im Benutzerhandbuch — `cdc.set_route`/`cdc.remove_route` mit Parametern und
      Rolle `cdc_admin` (Rollen-Tabelle, Abschnitt „Zugriff und Rollen"), die
      `rule_spec`-Form, R1–R6 mit Fehlertexten (führende Stelle: die Spec),
      Abhilfe — geht an `slice-routing-betriebsdoku`, dessen §2 den Gegenstand
      vollständig nennt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der
      Closure der Welle [welle-routing](../welle-routing.md) (die Roadmap führt sie
      unter *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/schema/nacharbeit-administration.sql` | update | `chk_administration_request_kind` um zwei Werte (Stand `30fd6cb5`: sieben Werte, Zeile 82), zwei Funktionen (Formvorbild `cdc.set_transformation`, Zeile 179), `REVOKE`/`GRANT EXECUTE` nur für `cdc_admin`, Idempotenz-Guard-Eintrag. Der Parametertyp von `rule_spec` ist `json` wie bei `set_transformation` (der Wortlaut der ADR sagt `jsonb`, §6). |
| `tools/schema/nacharbeit-roles.sql` | update | Rechte-Kommentar und Grants für die zwei Funktionen (Formvorbild Zeile 157). |
| `tools/schema/rolloutguard/guard.go`, `guard_test.go` | lesen / bei Bedarf update | die Wache kennt die Schreibweise der Funktions-Signaturen im Report (`set_transformation(in:json…)`); für die neuen Funktionen nur, wenn ihre Signatur sie braucht. |
| `tools/harness/run-schema-rollout-guard-test.sh` | update | der Alt-Tag-Lauf trägt eine hart verdrahtete Liste der Arten (Lauf 5) und die Signatur von `cdc.set_transformation`; Liste und Funktionen der neuen Arten nachziehen. |
| `internal/domain/model/administrationrequest.go` | update | die zwei Arten in der geschlossenen Menge, Konstruktor-Prüfung (`rule_name` Pflicht, `rule_spec` nur bei `set_route`); Formvorbild: die Transformations-Arten. |
| `internal/application/port/inbound/` (neue Datei neben `transformation.go`), `internal/application/port/outbound/` (neue Datei neben `transformation.go`) | neu | Inbound-Ports `SetRouteUseCase`/`RemoveRouteUseCase`, Outbound-Port des Regelstands (Lesen der `applied`-Zeilen). |
| `internal/application/usecase/setroute/`, `internal/application/usecase/removeroute/` (neue Pakete; Namen des Implementers) | neu | Validierung und R1–R6; Formvorbild `usecase/settransformation`, `usecase/removetransformation`. |
| `internal/application/usecase/excludecolumn/service.go` | update | R3 in der Gegenrichtung: `exclude_column` gegen eine Spalte, die eine Routing-Bedingung trägt, endet `failed`. |
| `internal/adapters/driven/postgresstorage/` (`tableactivation.go`, `queries/queries.go`, `sqlexec/translate.go`) | update | Regelstand-Lesen aus den `applied`-Zeilen, Annahme der zwei Arten in der Queue-Lesung; Formvorbild der Transformations-Arten (`BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue`: Lesefehler enden `failed`, keine Zeile hält die Queue an). |
| `internal/bootstrap/wiring.go`, `internal/bootstrap/assemblersync.go` | update | `applyAdministrationRequest` (Verdrahtung), `activatedTableBindings` und Aktivierungs-Zweig (Regelstand), Prozessstart, Setzen der Regelliste am `Assembler`. |
| `internal/bootstrap/*_test.go`, `internal/adapters/driven/postgresstorage/*_test.go`, Use-Case-Tests | neu / update | Happy/Boundary/Negative nach [`LH-FA-CFG-008`](../../../../spec/lastenheft.md); Rollen-Test unter den drei Logins; Ableitung mit gleichem `requested_at`. |
| `harness/targets/schema-rollout.md`, `harness/README.md` | update | Zählwörter und Beschreibungen, soweit der Suchlauf sie findet. |

**Konkretisierung vor dem Code** (Ist-Zustand am Parent `461ba1bf~1` gemessen, bevor
editiert wurde; Namen des Implementers):

| Datei / Komponente | Änderungs-Art | Konkretisierung |
|---|---|---|
| `internal/domain/model/routespec.go` (+ `routespec_test.go`) | neu | `ParseRouteSpec` (Form, unbekannte Schlüssel, Alphabet des Ziels; Reihenfolge der Fehlertext-Tabelle), `RouteSpec.Build`, `RouteRule.CheckIdentity` (R1, R2), `RouteRule.CheckArrangement` (R4, R5), `RoutesConditionOn` (R3 Gegenrichtung), `RouteRecord`/`FoldRoutes` (Ableitung aus den `applied`-Zeilen), `RouteSpecError` (Grund plus Adresswert), `MaxRouteOrder`. R3 (Spalte existiert, nicht ausgeschlossen) liegt im Use Case: Katalog und Ausschlussstand sind Port-Wissen. |
| `internal/domain/errors/errors.go` | update | sechs Sentinels für die Texte der Spec, die kein Transformations-Grund sind: `ErrRouteOrderTaken`, `ErrRouteConditionTaken`, `ErrRouteWithoutWhenTaken`, `ErrRouteWithoutWhenNotLast`, `ErrRouteBehindWithoutWhen`, `ErrColumnHasRouteCondition`; R1/R6/R3 tragen `ErrRuleNameTaken`/`ErrRuleNotKept`/`ErrRoutingColumnExcluded` (der Kern-Slice legte den letzten an) und `inbound.ErrSourceColumnMissing`. |
| `internal/application/port/inbound/routing.go`, `.../outbound/routing.go` | neu | `SetRouteUseCase`/`RemoveRouteUseCase` (Namen aus `spec/architecture.md`), `outbound.RoutingPort` mit `RoutingRules(ctx, source)`; die Spaltenexistenz und der Ausschlussstand kommen aus dem bestehenden `ColumnExclusionPort`. |
| `internal/application/usecase/setroute/`, `removeroute/` | neu | Prüfreihenfolge: Regelname · Formzeilen · Zielalphabet · R1 · R2 · R3 (fehlt, ausgeschlossen) · R4 · R5; `remove_route`: Regelname · R6. |
| `internal/application/usecase/excludecolumn/service.go` | update | Konstruktor trägt den `RoutingPort`; nach der Prüfung der Spaltenexistenz endet eine Spalte mit Routing-Bedingung als `Spalte trägt eine Routing-Bedingung: schema.table.column`. Aufrufer: `wiring.go` und drei Test-Stellen. |
| `internal/adapters/driven/postgresstorage/` | update | `RoutingRules` am `TableActivationAdapter` (`queries.SelectAppliedRoutingRequests`, `sqlexec.ReadRoutingRules`); die Lesung der offenen Anträge (`ListPending`) braucht keine Änderung, die Art kommt über den Domänen-Konstruktor. |
| `internal/bootstrap/wiring.go`, `assemblersync.go` | update | `activatedTableBindings` und `syncAssemblerAddBinding` tragen einen weiteren Port `routing` (Prozessstart, Aktivierungs-Zweig, API-Aktivierung); `applyAdministrationRequest` zwei Zweige (`SetRoute`/`RemoveRoute` am `Assembler`). `processedAdministrationKinds` leitet sich aus `model.AdministrationRequestKinds()` ab (gemessen: keine handgeführte Liste). |
| `internal/adapters/driven/postgresstorage/schema.sql` | keine Änderung (Nicht-Realisierung) | die zweite Schema-Beschreibung trägt weder `administration_request` noch `request_kind` (gemessen: `git grep -n -E 'administration\|request_kind' -- internal/adapters/driven/postgresstorage/schema.sql` druckt keine Zeile); die Antrags-Tabelle steht allein in `tools/schema/schema.yaml` und `nacharbeit-administration.sql`. |
| `tools/schema/plan.yaml`, `tools/schema/down.sql` | keine Änderung (Nicht-Realisierung) | die Erzeugnisse des Rollouts ändern sich durch diesen Slice nicht: weder Tabellen noch Spalten noch Views kommen hinzu (die Funktionen und die CHECK-Menge liegen in der Nacharbeit-Datei, nicht im d-migrate-Plan); gemessen am ersten Rollout gegen ein leeres Ziel: Differenz allein die Ziel-Bezeichnung in `plan.yaml`, `down.sql` unverändert. Beide Dateien werden nach dem Lauf aus dem Index wiederhergestellt (`rollout-restore`-Regel). |
| `internal/bootstrap/roles_rollout_file_internal_test.go`, `administration_roles_internal_test.go`, `routing_internal_test.go` (neu) | update / neu | die Signaturen der zwei Funktionen in der Datei-Prüfung; der Antragsweg unter den Rollen-Logins zieht die Routing-Anträge und die Sperre gegen `exclude_column` durch; Verdrahtungs-Tests mit abgeleitetem Stand. |
| `internal/application/port/outbound/administrationrequest.go`, `internal/domain/errors/errors.go`, `tools/harness/run-integration-tests.sh`, `harness/README.md` | update (Kommentar / Beschreibung) | Aufzählungen der Arten und Zählwörter („neun“ → „elf“ Fremdobjekte) nachgezogen, die der Suchlauf in §3 findet; keine Verhaltensänderung. |
| `internal/adapters/driven/postgresstorage/administrationrequest_routing_test.go` (neu), `administrationrequest_test.go` | neu / update | Funktionen schreiben `pending` und senden `NOTIFY`; Rollen-Test (42501 unter `cdc_reader`/`cdc_capture`, Gegenprobe `cdc_admin`); Annahmemenge von `rule_spec`; Ableitung mit gleichem `requested_at`; Aufruf-Ordnung in einer Transaktion; Menge der Arten (neun) und der Funktionen (neun). |

**Umsetzungsentscheidungen** (die Spec lässt sie offen; gemeldet, die Spec ändert dieser
Slice nicht):

- **Parametertyp `json`** (nicht `jsonb` wie im Wortlaut der ADR): der Bestand ist der Beleg —
  `cdc.set_transformation` trägt `p_rule_spec json` (`tools/schema/nacharbeit-administration.sql`),
  der Precheck-Report nennt die Funktion `…in:json`, die Spalte `rule_spec` ist `jsonb`.
  `cdc.set_route` folgt dem Bestand; die Messung am Ziel steht im Bericht.
- **Obergrenze von `order`:** `MaxRouteOrder` = 2147483647 (32-Bit-Ganzzahl). Grund: weit
  über jeder Zahl von Regeln einer Tabelle, und der Wert bleibt in jedem Ganzzahltyp eines
  Verbrauchers (SQL-`integer`, SDK-Sprachen) darstellbar; eine größere Zahl endet
  `rule_spec ist ungültig`. `SPEC-032` sagt „positive ganze Zahl“ ohne Obergrenze.
- **Schreibweise der Zahl:** nur eine Ganzzahl-Schreibweise (Ziffern ohne Bruch und Exponent).
  `10.0` und `1e1` enden `rule_spec ist ungültig`, obwohl sie mathematisch ganz sind;
  `jsonb` bewahrt die Schreibweise (`10.0` bleibt `10.0`). Die Spec nennt die Schreibweise
  nicht — Meldung an die Spec, nicht Änderung.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die `request_kind`-Menge
trägt zwei Werte mehr; es gibt zwei SQL-Funktionen mehr; der Regelstand hat eine
zweite Ableitung"; Parent ist `30fd6cb5`; der Implementer ergänzt die `diff`-Zeilen
und trägt Gefundenes und Nichtgefundenes ein):**

```suchlauf
30fd6cb5 40 -n set_transformation -- tools
30fd6cb5 43 -n remove_transformation -- tools internal ':!*_test.go'
30fd6cb5 55 -n -E 'RemoveTransformation|SetTransformation' -- internal ':!*_test.go'
30fd6cb5 3 -n -E 'sieben SQL-Funktionen|sieben Funktionen|sieben deklarierten' -- harness docs/user spec README.md tools
f7242d44 39 -n set_transformation -- tools
diff 39 -n set_transformation -- tools
f7242d44 44 -n remove_transformation -- tools internal ':!*_test.go'
diff 45 -n remove_transformation -- tools internal ':!*_test.go'
f7242d44 56 -n -E 'RemoveTransformation|SetTransformation' -- internal ':!*_test.go'
diff 56 -n -E 'RemoveTransformation|SetTransformation' -- internal ':!*_test.go'
f7242d44 2 -n -E 'sieben SQL-Funktionen|sieben Funktionen|sieben deklarierten' -- harness docs/user spec README.md tools
diff 0 -n -E 'sieben SQL-Funktionen|sieben Funktionen|sieben deklarierten' -- harness docs/user spec README.md tools
diff 3 -n -E 'neun SQL-Funktionen|neun Funktionen|neun deklarierten' -- harness docs/user spec README.md tools
f7242d44 11 -n -E 'neun bekannten|aktuell neun' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!.harness' ':!docs/plan/adr'
diff 0 -n -E 'neun bekannten|aktuell neun' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!docs/plan/planning/observations' ':!.harness' ':!docs/plan/adr'
diff 29 -n set_route -- tools
diff 36 -n remove_route -- tools internal ':!*_test.go'
diff 53 -n -E 'RemoveRoute|SetRoute' -- internal ':!*_test.go'
f7242d44 7 -n 'CREATE OR REPLACE FUNCTION cdc\.' -- tools/schema
diff 9 -n 'CREATE OR REPLACE FUNCTION cdc\.' -- tools/schema
```

| Träger | Messung am Parent (`30fd6cb5`, gemessen am 2026-10-01) | Behandlung und Befund am Diff |
|---|---|---|
| Stellen, an denen die Transformations-Antragsarten stehen (Muster für die zwei neuen) | Zeilen 1 bis 3: 40 Zeilen in `tools` (u. a. `nacharbeit-administration.sql`, `nacharbeit-roles.sql`, `rolloutguard`, `run-schema-rollout-guard-test.sh`, `run-integration-tests.sh`, `lib-sdk-rule-fixture.sh`), 43 Nicht-Test-Zeilen in `tools` und `internal`, 55 Nicht-Test-Zeilen zu den beiden Use-Case-Namen in `internal` | jede Stelle lesen: ist die Aufzählung der Arten eine geschlossene Liste, die um die zwei neuen wachsen muss (Guard, Alt-Tag-Lauf, Queue-Lesung, Verdrahtung), oder ein Beleg zu den Transformationen (bleibt). **Befund am Diff** (Zeilen 5 bis 10; die Dateiliste der Treffer mit `set_transformation` gegen die mit `set_route`/`SetRoute` verglichen): *gefunden und um die zwei Arten erweitert* — `nacharbeit-administration.sql` (CHECK, zwei Funktionen, `REVOKE`/`GRANT`), `nacharbeit-roles.sql` (Kommentar), `schema.yaml` (Kommentar, Beschreibung), `rolloutguard/guard.go`+`guard_test.go`, `run-schema-rollout-guard-test.sh`, `queries.go`, `model/administrationrequest.go`, `outbound/administrationrequest.go` und `errors.go` (Aufzählungen in Kommentaren), `wiring.go`, `harness/targets/schema-rollout.md`. *Gefunden, Beleg zu den Transformationen, unverändert:* `tools/harness/lib-sdk-rule-fixture.sh`, `run-integration-tests.sh` (E2E-Phasen; nur das Zählwort „elf“ im Kommentar nachgezogen), `examples/transformation-demo.sh`, `docs/user/e2e-abdeckung.md` (Erzeugnis), `harness/README.md` (Beschreibung der E2E-Phasen von `make test-integration`; nur das Zählwort „elf“ nachgezogen). Die Route-Gegenstücke liegen in den neuen Dateien (`routespec.go`, `routing.go`, `setroute/`, `removeroute/`). *Nicht gefunden:* eine handgeführte Liste der Arten außerhalb dieser Dateien — `processedAdministrationKinds` leitet sich aus `AdministrationRequestKinds()` ab; kein Träger in `compose.yaml`, `Dockerfile`, `.github/workflows`, `examples/` (außer dem Transformations-Demo). Zählung: 39 → 39 (`set_transformation`, `tools`), 44 → 45 (`remove_transformation`; die Differenz ist die Umbruch-Form der erweiterten Kommentarzeilen, gegen den Parent zeilenweise verglichen), 56 → 56 (Namen der Use Cases). Die Zeilen mit `set_route` (Zeile 16: 29) und `SetRoute` (Zeile 18: 53) sind die neuen Stellen. |
| Zählwörter der Funktionen und Arten | Zeile 4: 3 Zeilen (`harness/targets/schema-rollout.md:74` und `:90`, `spec/pflichtenheft.md:628`); am Parent dieses Diffs (`f7242d44`) 2, weil die Spec-Zeile `slice-routing-spec-nachzug` schon verändert hat | auf die gemessene Menge ziehen. **Befund am Diff:** `CREATE OR REPLACE FUNCTION cdc.` zählt 9 (Parent 7, gemessen); „sieben …“ → „neun …“ in `schema-rollout.md:74`/`:90` (Zeile 12: 0, Zeile 13: 3 mit der Spec-Zeile); das Zählwort der Fremdobjekte „neun“ → „elf“ in `guard.go`, `guard_test.go`, `run-schema-rollout-guard-test.sh`, `run-integration-tests.sh`, `schema-rollout.md` und `harness/README.md` (Zeile 14: 11 Zeilen am Parent, Zeile 15: 0 am Diff — der Träger `harness/README.md` stand im Suchlauf des Plans nicht, die Zeile 14 hat ihn gefunden); „sieben Arten“ → „neun Arten“ in `internal/domain/model/administrationrequest_test.go`, `internal/bootstrap/administration_roles_internal_test.go`, `queries.go` und dem Store-Test `administrationrequest_test.go`. *Nicht nachgezogen:* die Chronik-Zeile in `spec/pflichtenheft.md` (Änderungshistorie „sieben Werte“, datierter Eintrag) |
| Handbuch (Betreiber-Oberfläche) | `docs/user` im Suchraum von Zeile 4: bewegt sich erst mit `slice-routing-betriebsdoku` | gemeldet (Aufschub mit Adresse, §2), nicht mitgeändert. **Befund:** `docs/user/benutzerhandbuch.md` führt die Transformations-Funktionen in der Rollen-Tabelle (Zeile 76), im Abschnitt zur Regel-Konfiguration (Zeilen 361 bis 378) und im Glossar; der Gegenstand steht als committeter Text in §2 von `slice-routing-betriebsdoku` (Zeilen 39 bis 57 und 83 bis 85; `git grep -c -E 'cdc\.set_route\|rule_spec\|Rollen-Tabelle'` auf den Plan der Adresse: 8). Die Umsetzungsentscheidungen zu `order` (Obergrenze 2147483647, nur Ganzzahl-Schreibweise) sind neuer Gegenstand: gemeldet an den Planner |

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-routing-kern-label` liegt in `done/` und
kein anderer Slice liegt in `in-progress/` (WIP-Limit 1). Grund der ersten Bedingung:
der Antragsweg setzt die Regelliste an der Bindung und die Spalte voraus, die der
Kern liefert; ohne Auswertung wäre der Antrag ein Eintrag ohne Wirkung. Kollisionen:
der Slice ändert `internal/bootstrap/wiring.go` (`applyAdministrationRequest`) und
`tools/schema/nacharbeit-administration.sql`; beide stehen dann fest.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): der Slice ist mit drei
  Liefer-Punkten der größte der Welle. Zeigt sich nach der ersten Fixrunde, dass (B)
  und (C) zusammen mehr als drei Fixrunden brauchen, trennt sich (A) als
  `slice-routing-antragsweg-schema` ab (Schema, Funktionen, Guard, Rollen-Test);
  dazwischen gilt der Zwischenzustand der Transformations-Welle: die Funktion
  existiert, ein Antrag endet `failed` mit Text (Naht benannt, nicht erst im
  Nachhinein gesucht).
- `in-progress` → `open` (blockiert): das Verdikt zur Spec (z. B. Verhalten eines
  Antrags ohne Bindung) fehlt, oder der Alt-Tag-Lauf ist rot an der erweiterten
  CHECK-Menge — Carveout oder Architect-Zug, kein stilles Rot in `done/`.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code
ungefiltert), `make test-store` und der Alt-Tag-Lauf real grün, `make schema-rollout`
zweimal Exit 0, Suchlauf-Block nachgemessen, Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Zwischenzustand: Regeln setzbar, Backfill ohne Label.** Ab diesem Slice trägt
  jede neu erfasste WAL-Change das Ziel, ein Backfill-Run noch `NULL`; dieselbe
  Tabelle hat zwei Formen, bis `slice-routing-backfill-pfad` (unmittelbar danach)
  liegt. — **Ausgang:** bei der Closure einzutragen (Fenster ein Slice lang).
- **Parametertyp `json` gegen `jsonb`.** Die ADR nennt `rule_spec jsonb`;
  `cdc.set_transformation` trägt `json`, und die Wache `rolloutguard` kennt die
  Schreibweise der Report-Signatur (`in:json…`, `tools/schema/rolloutguard/guard.go`).
  Ein anderer Typ ändert die Signatur, die der Rollout-Report über einen Alt-Bestand
  meldet. — **Ausgang:** bei der Closure einzutragen (Typ nach Bestand, Wortlaut der
  Spec; der Bestand ist der Beleg, die ADR-Notation ist kein Gegenbeleg).
- **R3 in der Gegenrichtung ändert einen bestehenden Use Case.** `exclude_column`
  trägt bisher keine Kenntnis von Routing-Regeln; die Sperre ist eine zweite
  Abhängigkeit des Use Case. Das Queue-Lesen der Anträge läuft in **einer**
  Administrations-Goroutine, die Prüfung ist dadurch seriell (*erwartet*; die Annahme
  „genau eine Instanz" ist nicht erzwungen, `BEO-PGC/ein-instanz-annahme-ohne-erzwingung`,
  offen, 3×). — **Ausgang:** bei der Closure einzutragen (Regressionstest der
  bestehenden Pfade; Aussage zur Serialität mit ihrem Ursprung).
- **R4 und die Reihenfolge der Anträge.** Eine Abschlussregel ohne `when` muss die
  höchste `order` tragen; ein später gesetzter Antrag mit höherer `order` endet
  `failed`, auch wenn er ein anderer Zweck wäre. Der Betreiber löst es über
  `cdc.remove_route` und ein neues Setzen. — **Ausgang:** bei der Closure
  einzutragen (Handbuch-Adresse `slice-routing-betriebsdoku` für die Abhilfe).
- **Aufschub mit Adresse** (Handbuch-Beschreibung, §2). — **Ausgang:** bei der
  Closure einzutragen; der Gegenstand steht als committeter Text im §2 von
  `slice-routing-betriebsdoku` (Prüfung: `git grep` der Kernbegriffe
  `cdc.set_route`, `rule_spec`, `Rollen-Tabelle` im Plan der Adresse).
- **Coverage-Messgegenstand und Paketlisten.** Zwei neue Use-Case-Pakete liegen in der
  netzlos gemessenen Fläche; der Store-Teil im DB-Gegenstand
  (`tools/harness/db-package-lists-check.sh` hält die Paketlisten gleich). —
  **Ausgang:** bei der Closure einzutragen (`make coverage-gate`).
- **Rollen-Test.** Die Rechte der neuen Funktionen sind eine Wirkung, die kein Sensor
  außer dem Rollen-Test liest (`BEO-PGC/rollen-test-abdeckungsluecken`, gestrichen,
  4×; `BEO-PGC/rollen-verdrahtung`). — **Ausgang:** bei der Closure einzutragen.

## 7. Closure-Notiz

Wird bei der Closure gefüllt (vor dem `git mv` nach `done/`).

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit den
Pfaden `internal/domain/`, `internal/application/`, `internal/adapters/driven/`,
`internal/bootstrap/` und `tools/schema/` — eine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** das Register wurde durchgegangen
(Zähler = Zahl der `evidence/`-Dateien, ausgezählt am 2026-10-01); die Treffer
stehen in §6 und in der Eröffnungs-Sichtung der Welle
[welle-routing](../welle-routing.md) §6:
`BEO-PGC/schema-rollout-fremdobjekte` (verkörpert, 3×),
`BEO-PGC/d-migrate-nacharbeit` (verkörpert, 8×),
`BEO-PGC/antrag-mit-leerem-regelnamen-stallt-die-queue` (verkörpert),
`BEO-PGC/rollen-test-abdeckungsluecken` (gestrichen, 4×),
`BEO-PGC/rollen-verdrahtung` (4×), `BEO-PGC/ein-instanz-annahme-ohne-erzwingung`
(offen, 3×), `BEO-PGC/laufzeitzustand-ohne-dauerhaften-traeger` (offen, 1×),
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 21×). Keiner der
offenen Einträge erreicht mit diesem Slice 3× (`ein-instanz-annahme-ohne-erzwingung`
steht bereits bei 3×, Ausgang beim Lese-Schritt der Welle-Closure).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

