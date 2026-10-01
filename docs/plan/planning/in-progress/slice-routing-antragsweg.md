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

- [x] [`LH-FA-ADM-001`](../../../../spec/lastenheft.md) und
      [`LH-QA-SEC-002`](../../../../spec/lastenheft.md) (A): `cdc.set_route` und
      `cdc.remove_route` legen unter `cdc_admin` einen `pending`-Antrag der Art
      `set_route`/`remove_route` an; ein Aufruf unter `cdc_reader` und
      `cdc_capture` endet mit SQLSTATE `42501`; die `request_kind`-Menge trägt neun
      Werte; der Rollout zweimal hintereinander endet mit Exit 0; der Alt-Tag-Lauf
      (`tools/harness/run-schema-rollout-guard-test.sh`) kennt die neuen Arten und
      Funktionen (seine hart verdrahtete Liste der Arten ändert der Slice mit).
      *Zu belegen durch:* `make test-store` (Rollen-Test, Funktionen),
      `make schema-rollout` (zweimal), Alt-Tag-Lauf. *Beleg (Verifier):*
      Verifikations-Report §1 und §2 Zeile 1 (Rollout zweimal Exit 0, neun Werte in der
      CHECK-Menge, Rechte-Matrix und `permission denied` an einer Wegwerf-DB,
      `make test-store` Exit 0, Alt-Tag-Lauf Exit 0; Mutationen M5 und M10 rot).
- [x] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) Boundary (B): je eine
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
      `make test-store`. *Beleg (Verifier):* Verifikations-Report §2 Zeile 2 (Klartexte
      und Adressen am Quelltext gegen die Fehlertext-Tabelle von
      [`SPEC-019`](../../../../spec/pflichtenheft.md) gelesen, `make test` Exit 0;
      Mutationen M6a, M6b und M7 rot).
- [x] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) Happy Path (C): ein
      `applied`-Antrag setzt die Regelliste am laufenden `Assembler` ohne Neustart;
      nach einem Prozessstart ist der Regelstand aus den `applied`-Zeilen abgeleitet
      (Ordnung `requested_at`, dann `administration_request_id`, auch bei gleichem
      Zeitpunkt), und derselbe Stand entsteht im Aktivierungs-Zweig einer neu
      aktivierten Tabelle; ein Antrag gegen eine Tabelle ohne laufende Bindung endet
      wie bei den Spalten- und Transformations-Antragsarten (Spec-Zusage).
      *Zu belegen durch:* `internal/bootstrap`-Test mit realen Store-Zeilen
      (`make test-store`), Unit-Test der Ableitung (`make test`). *Beleg (Verifier):*
      Verifikations-Report §2 Zeile 3 (Antrag setzt die Regelliste live, Ableitung gegen
      reale PostgreSQL mit Zweitschlüssel, Aktivierungs-Zweig, Antrag ohne Bindung
      `applied`; Mutationen M8 und M9c rot, M9a rot über Panik). **Rest V-2:** die
      Feldbelegung `routing: activation` des HTTP-/gRPC-Aktivierungs-Pfads in
      `internal/bootstrap/wiring.go` bindet kein Unit-Test (M9b grün, Verifikations-Report
      §4); der Rest ist durch den DoD-Punkt „API-Aktivierung“ von
      [`slice-routing-e2e`](../open/slice-routing-e2e.md) §2 getragen (Aktivierung per
      HTTP bzw. gRPC mit `applied`-Routing-Regel, Gegenprobe ohne Regel,
      `make test-integration`, gedruckte Zeile je Weg), bis dahin trägt die Zeile an diesem
      Pfad die Weitergabe im Dekorator (M9a), nicht die Belegung.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9). *Beleg
      (Verifier):* Verifikations-Report §1 (Exit 0, `a-check`: 0 Befunde).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8). *Beleg:* Review
      [`review-slice-routing-antragsweg`](../../../reviews/review-slice-routing-antragsweg.md)
      (1 HIGH F-1, 2 MEDIUM F-2 und F-4, 4 LOW, 3 INFO) **und** Re-Review der Fixrunde
      `cc653a85`
      [`review-slice-routing-antragsweg-fixrunde-1`](../../../reviews/review-slice-routing-antragsweg-fixrunde-1.md)
      (0 HIGH, 0 MEDIUM, 2 LOW F-N1 und F-N2, 2 INFO F-N3 und F-N4); F-N1 und F-N2 sind
      im Nachzug `836e64d2` geschlossen (Store-Test-Godoc benennt die
      `jsonb`-Normalisierung, Kommentarzeile in `queries.go` umbrochen), F-N3 und F-N4
      tragen keine Aktion. Der Nachzug `836e64d2` ist von keinem Reviewer gelesen (nur
      Godoc, Fallname und Umbruch; keine Anweisung). Kein offenes HIGH/MEDIUM.
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-routing-antragsweg.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13). *Beleg
      (Verifier):* Verifikations-Report §1 und §2 Zeile 6 (22 Zeilen stimmen).
- [x] Doku-Update: `harness/targets/schema-rollout.md` (Zählwort der SQL-Funktionen
      und der Arten) und `harness/README.md` (nur bewegte Beschreibungen, laut
      Suchlauf). **Aufschub mit Adresse:** die Beschreibung der Betreiber-Oberfläche
      im Benutzerhandbuch — `cdc.set_route`/`cdc.remove_route` mit Parametern und
      Rolle `cdc_admin` (Rollen-Tabelle, Abschnitt „Zugriff und Rollen"), die
      `rule_spec`-Form, R1–R6 mit Fehlertexten (führende Stelle: die Spec),
      Abhilfe — geht an `slice-routing-betriebsdoku`, dessen §2 den Gegenstand
      vollständig nennt. *Beleg (Verifier):* Verifikations-Report §2 Zeile 7
      (`harness/targets/schema-rollout.md` und `harness/README.md` im Diff gelesen
      und gegen die gemessenen Zahlen gehalten; der Plan der Adresse trägt die
      `order`-Obergrenze und die Annahmemenge).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (§7).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert (§7: vier neue `evidence/`-Dateien in vier
      bestehenden Einträgen).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
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
- **Obergrenze von `order`:** `MaxRouteOrder` = 2147483647 (größter Wert einer 32-Bit-Ganzzahl),
  eine Setzung dieses Slice. Grund: weit über jeder Zahl von Regeln einer Tabelle; eine größere
  Zahl endet `rule_spec ist ungültig`. `SPEC-032` sagt „positive ganze Zahl“ ohne Obergrenze;
  das Handbuch nennt die Grenze (Aufschub-Adresse `slice-routing-betriebsdoku`).
- **Schreibweise der Zahl:** `SPEC-032` gilt wörtlich. Den Entscheid zu A-1 (Review-Frage an den
  Architect) hat der Hauptlauf der Sitzung getroffen, **kein Architect-Verdikt**: ein Report dazu
  liegt unter `docs/reviews/` nicht vor (Re-Review F-N4, `ls docs/reviews | grep -i routing` ohne
  Verdikt-Datei zu `order`). Die Richtung ist der Wortlaut der Spec und schließt die Abweichung des
  Erstentwurfs, statt eine neue zu setzen; die Obergrenze bleibt Umsetzungsgrenze. Das Handbuch
  nennt Obergrenze und Annahmemenge (Adresse `slice-routing-betriebsdoku` §2). Jede
  JSON-Zahl mit dem Wert einer positiven ganzen Zahl wird angenommen (`10`, `10.0`, `1e1`,
  `1E+1`); ein Bruchteil ungleich 0 (`1.5`), 0, ein negativer Wert, ein Nicht-Zahlen-Wert und
  ein Wert über `MaxRouteOrder` enden `rule_spec ist ungültig`. Der Wert wird exakt als
  rationale Zahl geprüft (`math/big`), nicht über eine Gleitkommazahl. Gemessen an PostgreSQL 18
  (`select (…)::json::jsonb::text`): `jsonb` bewahrt `10.0` und normalisiert `1e1`/`1E+1` zu
  `10`; die Annahme hängt an keiner der beiden Schreibweisen.

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
diff 54 -n -E 'RemoveRoute|SetRoute' -- internal ':!*_test.go'
f7242d44 7 -n 'CREATE OR REPLACE FUNCTION cdc\.' -- tools/schema
30fd6cb5 2 -n -E 'alle sieben Funktionen|sieben Antrags-Funktionen' -- internal
diff 0 -n -E 'alle sieben Funktionen|sieben Antrags-Funktionen' -- internal
diff 9 -n 'CREATE OR REPLACE FUNCTION cdc\.' -- tools/schema
```

| Träger | Messung am Parent (`30fd6cb5`, gemessen am 2026-10-01) | Behandlung und Befund am Diff |
|---|---|---|
| Stellen, an denen die Transformations-Antragsarten stehen (Muster für die zwei neuen) | Zeilen 1 bis 3: 40 Zeilen in `tools` (u. a. `nacharbeit-administration.sql`, `nacharbeit-roles.sql`, `rolloutguard`, `run-schema-rollout-guard-test.sh`, `run-integration-tests.sh`, `lib-sdk-rule-fixture.sh`), 43 Nicht-Test-Zeilen in `tools` und `internal`, 55 Nicht-Test-Zeilen zu den beiden Use-Case-Namen in `internal` | jede Stelle lesen: ist die Aufzählung der Arten eine geschlossene Liste, die um die zwei neuen wachsen muss (Guard, Alt-Tag-Lauf, Queue-Lesung, Verdrahtung), oder ein Beleg zu den Transformationen (bleibt). **Befund am Diff** (Zeilen 5 bis 10; die Dateiliste der Treffer mit `set_transformation` gegen die mit `set_route`/`SetRoute` verglichen): *gefunden und um die zwei Arten erweitert* — `nacharbeit-administration.sql` (CHECK, zwei Funktionen, `REVOKE`/`GRANT`), `nacharbeit-roles.sql` (Kommentar), `schema.yaml` (Kommentar, Beschreibung), `rolloutguard/guard.go`+`guard_test.go`, `run-schema-rollout-guard-test.sh`, `queries.go`, `model/administrationrequest.go`, `outbound/administrationrequest.go` und `errors.go` (Aufzählungen in Kommentaren), `wiring.go`, `harness/targets/schema-rollout.md`. *Gefunden, Beleg zu den Transformationen, unverändert:* `tools/harness/lib-sdk-rule-fixture.sh`, `run-integration-tests.sh` (E2E-Phasen; nur das Zählwort „elf“ im Kommentar nachgezogen), `examples/transformation-demo.sh`, `docs/user/e2e-abdeckung.md` (Erzeugnis), `harness/README.md` (Beschreibung der E2E-Phasen von `make test-integration`; nur das Zählwort „elf“ nachgezogen). Die Route-Gegenstücke liegen in den neuen Dateien (`routespec.go`, `routing.go`, `setroute/`, `removeroute/`). *Nicht gefunden:* eine handgeführte Liste der Arten außerhalb dieser Dateien — `processedAdministrationKinds` leitet sich aus `AdministrationRequestKinds()` ab; kein Träger in `compose.yaml`, `Dockerfile`, `.github/workflows`, `examples/` (außer dem Transformations-Demo). Zählung: 39 → 39 (`set_transformation`, `tools`), 44 → 45 (`remove_transformation`; die Differenz ist die Umbruch-Form der erweiterten Kommentarzeilen, gegen den Parent zeilenweise verglichen), 56 → 56 (Namen der Use Cases). Die Zeilen mit `set_route` (Zeile 16: 29) und `SetRoute` (Zeile 18: 54, eine Kommentarzeile mehr nach der Fixrunde zu F-5) sind die neuen Stellen. |
| Zählwörter der Funktionen und Arten | Zeile 4: 3 Zeilen (`harness/targets/schema-rollout.md:74` und `:90`, `spec/pflichtenheft.md:628`); am Parent dieses Diffs (`f7242d44`) 2, weil die Spec-Zeile `slice-routing-spec-nachzug` schon verändert hat | auf die gemessene Menge ziehen. **Befund am Diff:** `CREATE OR REPLACE FUNCTION cdc.` zählt 9 (Parent 7, gemessen); „sieben …“ → „neun …“ in `schema-rollout.md:74`/`:90` (Zeile 12: 0, Zeile 13: 3 mit der Spec-Zeile); das Zählwort der Fremdobjekte „neun“ → „elf“ in `guard.go`, `guard_test.go`, `run-schema-rollout-guard-test.sh`, `run-integration-tests.sh`, `schema-rollout.md` und `harness/README.md` (Zeile 14: 11 Zeilen am Parent, Zeile 15: 0 am Diff — der Träger `harness/README.md` stand im Suchlauf des Plans nicht, die Zeile 14 hat ihn gefunden); „sieben Arten“ → „neun Arten“ in `internal/domain/model/administrationrequest_test.go`, `internal/bootstrap/administration_roles_internal_test.go`, `queries.go` und dem Store-Test `administrationrequest_test.go`. *Nicht nachgezogen:* die Chronik-Zeile in `spec/pflichtenheft.md` (Änderungshistorie „sieben Werte“, datierter Eintrag) |
| Handbuch (Betreiber-Oberfläche) | `docs/user` im Suchraum von Zeile 4: bewegt sich erst mit `slice-routing-betriebsdoku` | gemeldet (Aufschub mit Adresse, §2), nicht mitgeändert. **Befund:** `docs/user/benutzerhandbuch.md` führt die Transformations-Funktionen in der Rollen-Tabelle (Zeile 76), im Abschnitt zur Regel-Konfiguration (Zeilen 361 bis 378) und im Glossar; der Gegenstand steht als committeter Text in §2 von `slice-routing-betriebsdoku` (Zeilen 39 bis 57 und 83 bis 85; `git grep -c -E 'cdc\.set_route\|rule_spec\|Rollen-Tabelle'` auf den Plan der Adresse: 8). Die Umsetzungsentscheidungen zu `order` (Obergrenze 2147483647; jede JSON-Zahl mit dem Wert einer positiven ganzen Zahl wird angenommen) sind neuer Gegenstand: als Handbuch-Gegenstand in den Übergabe-Block von `slice-routing-betriebsdoku` §2 eingetragen (Fixrunde zu F-2) |
| Wendungen „alle sieben Funktionen“ / „sieben Antrags-Funktionen“ in Godoc von `internal` (Fixrunde zu F-7) | Zeile 20: 2 Zeilen am Parent (`administrationrequest_order_test.go:162`, `administration_endtoend_test.go:37`) | auf den Ist-Stand gezogen. **Befund am Diff:** Zeile 21: 0; „sieben der neun Funktionen“ (der Test fährt sieben der neun) und „jede der sieben“ (die sieben Funktionen der zwei Folgen) bleiben als zutreffend stehen. *Nicht gefunden:* weitere Wendungen dieser Form in `internal` |

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
  liegt. — **Ausgang:** weiter offen. Das Fenster ist das benannte Fenster der Welle
  ([`welle-routing`](../welle-routing.md) §5), nicht gemessen; Adresse:
  [`slice-routing-backfill-pfad`](../open/slice-routing-backfill-pfad.md), dessen Start-Trigger
  diesen Slice in `done/` voraussetzt.
- **Parametertyp `json` gegen `jsonb`.** Die ADR nennt `rule_spec jsonb`;
  `cdc.set_transformation` trägt `json`, und die Wache `rolloutguard` kennt die
  Schreibweise der Report-Signatur (`in:json…`, `tools/schema/rolloutguard/guard.go`).
  Ein anderer Typ ändert die Signatur, die der Rollout-Report über einen Alt-Bestand
  meldet. — **Ausgang:** entfallen. Typ `json` nach Bestand und Wortlaut von
  [`SPEC-019`](../../../../spec/pflichtenheft.md); Alt-Tag-Lauf Exit 0 (Verifikations-Report
  §1), die Mutation `set_route`-Signatur `in:jsonb` färbt drei Tests rot (Review G2), die
  `Accepted`-ADR bleibt unberührt (Review F-9, Verifikations-Report §6).
- **R3 in der Gegenrichtung ändert einen bestehenden Use Case.** `exclude_column`
  trägt bisher keine Kenntnis von Routing-Regeln; die Sperre ist eine zweite
  Abhängigkeit des Use Case. Das Queue-Lesen der Anträge läuft in **einer**
  Administrations-Goroutine, die Prüfung ist dadurch seriell (*erwartet*; die Annahme
  „genau eine Instanz" ist nicht erzwungen, `BEO-PGC/ein-instanz-annahme-ohne-erzwingung`,
  offen, 3×). — **Ausgang:** Regression **entfallen**: die Bestandstests von
  `exclude_column` bleiben grün, der Regelstand wird erst nach der Spaltenexistenz
  gelesen, Mutationen M6a und M6b rot (Verifikations-Report §2 Zeile 2, §3, §4). Die
  Serialität bleibt **erwartet**, nicht gemessen (Ursprung: Lesen des Codes, eine
  Administrations-Goroutine); weiter offen unter
  `BEO-PGC/ein-instanz-annahme-ohne-erzwingung`, Ausgang beim Lese-Schritt der Closure
  von [`welle-routing`](../welle-routing.md).
- **R4 und die Reihenfolge der Anträge.** Eine Abschlussregel ohne `when` muss die
  höchste `order` tragen; ein später gesetzter Antrag mit höherer `order` endet
  `failed`, auch wenn er ein anderer Zweck wäre. Der Betreiber löst es über
  `cdc.remove_route` und ein neues Setzen. — **Ausgang:** weiter offen (Aufschub
  mit Adresse). Das Verhalten ist durch R4 gebunden (Verifikations-Report §2 Zeile 2); die
  Abhilfe beschreibt das Handbuch: Adresse [`slice-routing-betriebsdoku`](../open/slice-routing-betriebsdoku.md)
  §2 (Fall R4 steht dort als committeter Text, `git grep -n -E 'R4|höchste'` im Plan der
  Adresse: Zeile 86).
- **Aufschub mit Adresse** (Handbuch-Beschreibung, §2). — **Ausgang:** weiter offen;
  der Gegenstand steht als committeter Text im §2 von `slice-routing-betriebsdoku`
  (Prüfung: `git grep -n -c -E 'cdc\.set_route|rule_spec|Rollen-Tabelle'` im Plan der
  Adresse druckt 9, gemessen bei der Closure; der Verifier bestätigt die Aufnahme,
  Verifikations-Report §2 Zeile 7).
- **Coverage-Messgegenstand und Paketlisten.** Zwei neue Use-Case-Pakete liegen in der
  netzlos gemessenen Fläche; der Store-Teil im DB-Gegenstand
  (`tools/harness/db-package-lists-check.sh` hält die Paketlisten gleich). —
  **Ausgang:** entfallen: `make gates` Exit 0 (enthält `coverage-gate`), der DB-Teil im
  `make test-store`-Lauf des Verifiers gedruckt `db-coverage: OK — DB-Adapter-Coverage
  82.99% erfuellt Schwelle 80%` (Verifikations-Report §1; der Wert ist lauf-gebunden,
  keine Zustandsgröße).
- **Rollen-Test.** Die Rechte der neuen Funktionen sind eine Wirkung, die kein Sensor
  außer dem Rollen-Test liest (`BEO-PGC/rollen-test-abdeckungsluecken`, gestrichen,
  4×; `BEO-PGC/rollen-verdrahtung`). — **Ausgang:** entfallen: der Store-Test
  `TestAdministrationRequestRoutingFunctionsRequireCdcAdminMembership` und der Alt-Tag-Lauf
  tragen die Rechte, die Mutation `REVOKE` ohne `cdc.set_route` (M5) färbt den Store-Test
  rot (Verifikations-Report §2 Zeile 1, §4). Die Rollen-Tests in `internal/bootstrap` binden
  das `REVOKE` nicht (V-5, Re-Review F-N3): ein Schnitt zwischen zwei Testebenen, beide
  tragen, keine Lücke.
- **API-Aktivierung: Feldbelegung `routing: activation` in `Run` ungebunden** (Review F-4,
  Verifikation V-2; neu in §6 aufgenommen bei der Closure). Die Weitergabe im Dekorator ist
  gebunden (M9a rot, über Panik), die Belegung in `internal/bootstrap/wiring.go` nicht (M9b
  grün). Ein `nil`-Port wäre im Betrieb ein lauter Fehler beim ersten API-Aktivieren, kein
  stilles Ausbleiben (*hergeleitet*). — **Ausgang:** weiter offen. Adresse:
  [`slice-routing-e2e`](../open/slice-routing-e2e.md) §2, DoD-Punkt „API-Aktivierung“ (im Plan
  der Adresse gelesen, Verifikations-Report §5 Zeile F-4; `git grep -n 'API-Aktivierung'` im Plan
  der Adresse: Zeile 108).
- **Wiederholungs-Grenze der Queue (Review F-5).** Der Kommentar in `applyAdministrationRequest`
  (`internal/bootstrap/wiring.go`) sagt, eine Wiederholung eines `set_route` nach einem
  fehlgeschlagenen `applied`-Vermerk sei folgenlos; der Satz ist als **hergeleitet** gekennzeichnet
  (R1 liest nur `applied`-Zeilen, `Assembler.SetRoute` ersetzt nach dem Regelnamen), ein
  Wiederholungs-Test fehlt (die Transformations-Antragsarten tragen einen:
  `TestProcessAdministrationRequestsSetTransformationIsIdempotent`). — **Ausgang:** weiter
  offen. Adresse: der Kommentar an dieser Stelle (der Marker „Hergeleitet aus dem Code, ohne
  Wiederholungs-Test“ ist der Wächter); die Entscheidung gegen einen Test ist offen benannt,
  ein Träger-Slice ist nicht angelegt. Der Planner von
  [`slice-routing-e2e`](../open/slice-routing-e2e.md) liest diesen Punkt beim Start und
  entscheidet, ob der Neustart-Lauf dort ihn belegen kann.
- **Heredoc-Selbstauskunft des Implementers (Verifikation V-7).** Der Implementer meldete
  einen früheren Heredoc-Schreibvorgang auf `internal/application/usecase/excludecolumn/service.go`
  (Verstoß gegen [`AGENTS.md`](../../../../AGENTS.md) §3.1, einmal geleert, per `git checkout`
  wiederhergestellt). Am Endstand ist er nicht erkennbar (`git grep -E '^EOF$|<<.?EOF' -- internal`:
  0, `make fmt-check` Exit 0, Diff von Review und Verifier gelesen), der Guard liest
  Umleitungen nicht. — **Ausgang:** weiter offen, **übernommen, nicht gemessen**: eine
  Messung am Endstand ist strukturell unmöglich; kein Träger-Slice, keine
  `evidence/`-Datei (ohne Spur kein nachprüfbares Auftreten, Entscheidung in §7).

## 7. Closure-Notiz

- **Was hat funktioniert:** Die drei Liefer-Punkte tragen: zwei Antragsarten und zwei
  Funktionen allein für `cdc_admin`, Validierung R1 bis R6 in Go mit den Klartexten der Spec
  (R3 in beide Richtungen), Regelstand aus den `applied`-Zeilen in Prozessstart,
  Aktivierungs-Zweig und `applyAdministrationRequest`. Gemessen im Verifier-Lauf: `make test`,
  `make test-store`, `make test-replication`, zweimal `make schema-rollout`, Alt-Tag-Lauf,
  `make gates`, `make docs-check`, Suchlauf (22 Zeilen) alle Exit 0 (Verifikations-Report §1).
  Das Muster der Transformations-Antragsarten (Formvorbild, Guard-Liste, Alt-Tag-Lauf) trug den
  Schnitt ohne Rückführung nach §4; die in §4 vorab benannte Abtrennung von (A) ist nicht
  ausgelöst (Review F-10).
- **Was ging anders als geplant:** Eine Fixrunde (`cc653a85`) und ein Nachzug (`836e64d2`).
  Das HIGH F-1 war eine Tatsachenbehauptung des Plans: die Annahme, `jsonb` bewahre die
  Schreibweise der Zahl, sodass `1e1` abgelehnt werde. Gemessen an PostgreSQL 18 bewahrt `jsonb`
  `10.0` und normalisiert `1e1`/`1E+1` zu `10`; ein Testfall, der `ParseRouteSpec` mit Text
  aufrief, belegte einen Pfad, den kein Antrag erreicht. Mit F-2 (Setzung enger als
  [`SPEC-032`](../../../../spec/pflichtenheft.md)) führte das zur Annahme jeder JSON-Zahl mit dem
  Wert einer positiven ganzen Zahl, exakt über `math/big`. Den Entscheid zu `order` hat der
  Hauptlauf der Sitzung getroffen, **nicht ein Architect-Verdikt**: ein Report dazu liegt nicht
  vor (Re-Review F-N4); die Richtung ist der Wortlaut der Spec, die Obergrenze 2147483647 bleibt
  Umsetzungsgrenze, das Handbuch trägt sie über die Adresse `slice-routing-betriebsdoku`. F-4
  (Dekorator-Test mit leerem Regelstand) wurde bis auf die Feldbelegung in `Run` gebunden; der
  Rest liegt bei `slice-routing-e2e`. Anders als in den beiden Vorgänger-Slices wurde die Fixrunde
  **von einem Reviewer gelesen** (Re-Review, Verifier-Empfehlung V-1); der Nachzug `836e64d2`
  (Godoc, Fallname, Umbruch) ist wieder ungelesen, er enthält keine Anweisung. Die Zahl der
  Mutationen trägt ihren Ursprung (Instanz A von [`AGENTS.md`](../../../../AGENTS.md) §3.12):
  Implementer 23 + 10 + 2 = 35, **übernommen** (Bericht des Implementers, nicht nachgefahren);
  Reviewer 23 selbst gefahren, **gemessen** (19 rot, 4 grün: vier äquivalente beziehungsweise
  redundante Mutanten, F-8); Verifier 12 selbst gefahren, **gemessen** (11 rot, eine grün: M9b,
  die Feldbelegung); Re-Review 5 selbst gefahren, **gemessen** (R1 bis R5, jede in mindestens
  einer Ebene rot). Die Mutationen des Reviews übernahm der Verifier, die des Verifiers der
  Re-Review: Wer die Zahl liest, liest die Läufe eines unabhängigen Lesers, nicht die des
  Autors.
- **Steering-Loop-Eintrag:** geschärfte Regel, kein neuer Sensor. Ein Plan-Satz über das
  Verhalten des Fremdsystems („was PostgreSQL bewahrt oder normalisiert“) ist eine
  **Tatsachenbehauptung** ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B) und gehört vor
  dem Schreiben an der gepinnten Instanz gemessen, mit dem Messort im Satz (hier
  `select (…)::json::jsonb::text`); eine Annahme dieser Art steht sonst als *erwartet*. F-1 (HIGH)
  entstand, weil `1e1` als „bleibt erhalten“ galt und ein Testfall den Direktaufruf des Parsers
  statt des Antragswegs fuhr; die Messung dauert eine Anweisung, der Reviewer fand sie vor dem
  Merge. Zweite Hälfte: ein Test, der eine Schreibweise „bindet“, fährt den Weg, den die Eingabe
  real geht (Funktion, `jsonb`, Queue, Parser); der Store-Test tat es und band `1e1` trotzdem
  nicht, weil `jsonb` den Exponenten vor dem Parser normalisiert (F-N1, im Godoc benannt).
  Träger: die Lese-Handlung des Planners beim Schreiben von Plan-Sätzen und die verkörperten
  Einträge im Register (unten), keine neue Regel im Text von `AGENTS.md`: §3.12 Instanz B
  trägt den Satz bereits, der Slice liefert seine Ausprägung „Verhalten des Fremdsystems“.
- **Beobachtungs-Register (`../observations/`):** vier `evidence/`-Dateien in vier
  bestehenden Einträgen, Zähler real ausgezählt (`ls …/evidence | wc -l` je Eintrag):
  - **`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`** (verkörpert): F-1 (HIGH, daher Datei
    trotz Deckel), Form **Plan-Satz und Testfall**; Zähler **19×** → **20×**.
  - **`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung`** (verkörpert): F-1 ist
    dieselbe Aussage in ihrer Instanz-B-Ausprägung („Verhalten des Fremdsystems“, Plan-Satz
    als Begründung einer Umsetzungsentscheidung); eine eigene Datei, weil HIGH; Zähler
    **12×** → **13×**. Beide Einträge zählen denselben Vorgang mit verschiedenem Mechanismus
    (dort trägt der Beleg den Satz nicht, hier trägt die Behauptung nicht).
  - **`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`** (verkörpert): F-4 (MEDIUM, daher
    Datei trotz Deckel), Form **Dekorator-Test mit leerem Fake-Port**; Zähler **21×** → **22×**.
  - **`BEO-PGC/fixrunde-ohne-reviewer-lesung`** (offen): Verifikation V-1, drittes Auftreten;
    Zähler **2×** → **3×**, damit Schwelle erreicht, Ausgang beim Lese-Schritt der Closure
    von [`welle-routing`](../welle-routing.md). **Gegenbeleg im Eintrag vermerkt:** hier hat der
    Verifier das Muster benannt und einen Re-Review empfohlen, der Re-Review wurde gefahren und
    fand F-N1 (LOW), das dem Verifier entgangen war; anders als die beiden Vorgänger änderte
    diese Fixrunde Anweisungen (den Parser).
  - Ohne eigene Datei (Deckel, ≤ LOW, vor dem Merge gefunden, bekannter Träger-Typ):
    F-N1 (LOW, Beleg trägt seinen Satz nur zur Hälfte, `beleg-befehl-traegt-seinen-satz-nicht`),
    F-6 und F-N2 (LOW, Umbruch nach Teilersetzung, Zähler läuft unter
    `BEO-PGC/nachzug-laesst-ueberholten-text-stehen`), F-7 (LOW, Godoc „alle sieben Funktionen“
    außerhalb des Suchmusters, `arbeit-ueberholt-stehenden-traeger`), F-5 (LOW, Zusage im Kommentar
    ohne Marker, `negativtest-ohne-bindung-an-seine-eingabe`; Marker nachgezogen).
  - V-7 (Heredoc-Selbstauskunft): **keine** Datei bei `inplace-textwerkzeug-am-repo-trotz-nutzerregel`.
    Der Zähler misst Auftreten mit nachprüfbarer Spur; am Endstand gibt es keine, die Aussage ist
    **übernommen** (§6). Wird eine Spur sichtbar, folgt die Datei.
  - Kein Eintrag steht bei 3× oder mehr ohne Ausgang an, den dieser Slice neu erreichte, außer
    `fixrunde-ohne-reviewer-lesung` (3×, Ausgang beim Lese-Schritt der Welle-Closure).
- **Folge-Slices:** keine neuen. Übergaben: `slice-routing-e2e` (API-Aktivierung, Neustart-Lauf
  mit realen Zeilen V-6, Wiederholungs-Grenze), `slice-routing-betriebsdoku` (Betreiber-Oberfläche,
  Abhilfe R4, `order`-Obergrenze und Annahmemenge), `slice-routing-backfill-pfad` (Fenster
  Backfill ohne Label).
- **Risiken aus §6:** Zwischenzustand Backfill ohne Label **weiter offen** (Adresse
  `slice-routing-backfill-pfad`); Parametertyp `json` **entfallen**; R3-Gegenrichtung
  **entfallen**, Serialität **weiter offen** (erwartet, Register); R4-Abhilfe **weiter offen**
  (Adresse `slice-routing-betriebsdoku`); Aufschub Handbuch **weiter offen** (gleiche Adresse);
  Coverage **entfallen**; Rollen-Test **entfallen**; API-Aktivierung **weiter offen** (Adresse
  `slice-routing-e2e` §2); Wiederholungs-Grenze **weiter offen** (hergeleitet); Heredoc V-7 **weiter
  offen** (übernommen, nicht gemessen). Befunde ohne Risiko-Eintrag: V-3 (nicht reproduzierbare
  Aussage zu `run-store-tests.sh`, `make test-store` Exit 0), V-4 (durch F-N2 im Nachzug behoben),
  V-5 (Testebenen-Schnitt), V-6 (Adresse `slice-routing-e2e`).
- **Drei Paarungen:** dieser Slice gehört zu [welle-routing](../welle-routing.md) (offen) — die
  Prüfung läuft regelkonform bei deren Closure; die DoD-Zeile bleibt deshalb `[ ]`. (a) Anker: der
  Lerneintrag verkörpert nichts neu (Ausprägung unter §3.12 Instanz B, Register); (b) Folge-Slice:
  keiner neu, die genannten Pläne liegen unter `open/`; (c) Register: die genannten Kennungen
  existieren als Verzeichnis, jede trägt ein nicht leeres `evidence/`.
- **Validator (Modul 8):** entfällt — der Antragsweg ohne Lesewege und ohne Backfill-Label ist
  für Betreiber noch nicht als Ganzes nutzbar; der Nutzer-Bedarf
  ([`LH-FA-CFG-008`](../../../../spec/lastenheft.md)) wird erst durch den Wellen-Beleg
  validierbar.

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

