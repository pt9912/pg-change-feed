# Review-Report: slice-routing-antragsweg — 2026-10-01

**Review-Art:** Code — der Diff ist Produktivcode (Domäne, Use Cases, Ports,
Store-Adapter, Verdrahtung), Schema-Nacharbeit, Rollout-Wache, Alt-Tag-Skript,
Tests und Slice-Plan; geprüft gegen Plan, ADRs, Spec und `AGENTS.md` Hard Rules
(Modul 10 §Drei Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe
(Modul 11).

**Gegenstand:** `git diff 461ba1bf~1 HEAD` (Basis `30fd6cb5`), vier Commits
`7507f20a` (Plan konkretisiert), `2551e995` (Go: Regelform, Use Cases,
Verdrahtung), `622cb55d` (Schema, Funktionen, Wache, Alt-Tag-Lauf) und
`fcd6ab70` (Suchlauf-Feld); 40 Dateien, +3065/−132, Slice-Plan im Lifecycle
`in-progress`.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(seither um weitere HIGH-Klassen ergänzt).
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-01.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan `slice-routing-antragsweg` (Ziel, Umfang, §3 samt Umsetzungs-
  entscheidungen und Suchlauf-Feld, §6) und Welle `welle-routing`
- [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  (Teilfragen 3 und 6, R1–R6, Folgepflicht 3),
  [`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md),
  [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md),
  [`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md)
  (nur als Nachbarn), Muster
  [`ADR-0050`](../plan/adr/0050-sql-administration-antragsqueue-und-live-reload.md),
  [`ADR-0065`](../plan/adr/0065-spaltenausschluss-dauerhafter-traeger.md),
  [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md),
  [`ADR-0125`](../plan/adr/0125-transformationen-parametertyp-regelform-json.md),
  [`ADR-0127`](../plan/adr/0127-antrags-queue-requested-at-aufrufzeitpunkt.md),
  [`ADR-0046`](../plan/adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
- Lastenheft [`LH-FA-CFG-008`](../../spec/lastenheft.md),
  [`LH-FA-ADM-001`](../../spec/lastenheft.md),
  [`LH-QA-SEC-002`](../../spec/lastenheft.md),
  [`LH-QA-SEC-004`](../../spec/lastenheft.md) und die Pflichtenheft-Stellen
  [`SPEC-019`](../../spec/pflichtenheft.md) (Antrags-Datensatz, R1–R6,
  Fehlertext-Tabelle, Dauerhaftigkeit),
  [`SPEC-032`](../../spec/pflichtenheft.md) (Regelform),
  [`SPEC-008`](../../spec/pflichtenheft.md),
  [`ARC-005`](../../spec/architecture.md)
- `AGENTS.md` (Hard Rules §3.1, §3.7, §3.9, §3.12, §3.13, §3.15),
  `harness/conventions.md`
- Vorgänger-Reviews der Welle: `review-slice-routing-spec-nachzug` (die
  Rücknahme der `order`-Obergrenze in der Fixrunde)

**Eigenständig durchgeführte Prüfungen** (gemessen, nicht aus dem
Implementer-Bericht übernommen):

- `make kommentar-kennungen DIFF=461ba1bf~1` → Exit 0, keine Ausgabe (kein
  Kandidat); `make fmt-check` → Exit 0, „313 Go-Dateien geprüft, alle
  formatiert"; `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-routing-antragsweg.md`
  → Exit 0, „20 Zeilen stimmen"; `RANGE=461ba1bf~1..HEAD make commit-traceability`
  → Exit 0, sechs Commits.
- Volllauf der Unit-Tests auf einer Kopie des Arbeitsbaums im Scratchpad
  (`go test -race ./...` im gepinnten `TOOLCHAIN_RACE_IMAGE`): grün, kein
  FAIL, kein Build-Fehler.
- `tools/harness/run-schema-rollout-guard-test.sh` auf der Kopie, Exit 0, alle
  sechs Läufe; Lauf 5 „OK — Tag v0.4.0": der Alt-Tag-Lauf kennt die neun Arten
  und die beiden neuen Funktionen mit den Rechten `cdc_admin t` /
  `cdc_capture f` / `cdc_reader f`, Lauf 2 nimmt den `--allow-destructive`-Pfad
  gegen den realen d-migrate-Report mit den elf Fremdobjekten.
- Store- und Verdrahtungs-Tests gegen eine Wegwerf-PostgreSQL 18 (gepinnter
  Digest, Schema-Stand über `tools/schema/apply-rollout.sh`, eigener
  Läufer im Scratchpad mit `-run`-Filter auf Routing, Antrags-Queue und
  Rollen): grün (Basislauf), danach die Mutationen unten.
- Spec-Abgleich Satz für Satz: [`SPEC-019`](../../spec/pflichtenheft.md)
  Fehlertext-Tabelle und R1–R6 gegen `routespec.go`, `setroute`,
  `removeroute`, `excludecolumn` und die Sentinel-Texte in `errors.go` — alle
  15 Klartexte wortgleich; Prüfreihenfolge Regelname · Formzeilen (Objekt,
  unbekannter Schlüssel, Form, Alphabet) · R1 · R2 · R3 (fehlt, ausgeschlossen)
  · R4 · R5 stimmt mit der Tabelle von oben nach unten; `remove_route` nur
  Regelname und R6.
- Messung an PostgreSQL 18 (gepinnter Digest, einmalig): `select
  ('{"order": 1e1, "b": 10.0, "c": 1E+1, "d": 10}')::json::jsonb::text` druckt
  `{"b": 10.0, "c": 10, "d": 10, "order": 10}` — Grundlage von F-1.

**Mutationen, selbst gefahren** (Kopie im Scratchpad, Änderung als `sed … >
Datei-in-der-Kopie`, kein `-i`, Rücknahme per `git checkout` in der Kopie;
Unit im gepinnten Race-Image, Store über die Wegwerf-PostgreSQL). 23 Mutationen,
Farbe je Mutation:

| Nr | Stelle | Mutation | Ergebnis |
|---|---|---|---|
| M1 | `setroute/service.go` | R3-Ausschlussprüfung gegen die falsche Tabelle (`excluded["x"]`) | rot (Use-Case- und Verdrahtungs-Test) |
| M2 | `excludecolumn/service.go` | Regelstand-Schlüssel nur `Table` statt `schema.table` | rot (beide Ebenen) |
| M3 | `routespec.go` R4 | `>` → `>=` | **grün** — äquivalent: gleiche `order` endet vorher an R2 (INFO F-8) |
| M4 | `routespec.go` | `order > Max` → `>= Max` | rot (Grenzwert-Test) |
| M5 | `routespec.go` R5 | `equals` aus dem Paarvergleich entfernt | rot |
| M6 | `wiring.go` Prozessstart | `Routes: routes["x"]` | rot (`TestActivatedTableBindingsCarriesRoutes`) |
| M7 | `assemblersync.go` | `Routes: routes["x"]` | rot (Disable/Enable-Zyklus) |
| M8 | `routespec.go` | `order < 1` → `order < 0` | rot |
| M9 | `setroute/service.go` | R3 „Spalte fehlt" abgeschaltet | rot (beide Ebenen) |
| M10 | `wiring.go` | `assembler.SetRoute` entfernt | rot (drei Verdrahtungs-Tests) |
| M11 | `routespec.go` | UTF-8-Prüfung abgeschaltet | rot |
| M12 | `routespec.go` | `when.column`: Leer- und NUL-Prüfung entfernt | **grün** (F-8) |
| M13 | `routespec.go` | nur NUL-Prüfung entfernt | **grün** (F-8) |
| M14 | `routespec.go` | nur Leer-Prüfung entfernt | **grün** (F-8) |
| G1 | `rolloutguard/guard.go` | Eintrag `remove_route` gestrichen | rot (zwei Tests) |
| G2 | `rolloutguard/guard.go` | `set_route`-Signatur `in:jsonb` | rot (drei Tests) |
| S1 | `queries.go` | `ORDER BY requested_at` ohne Zweitschlüssel | rot (Tie-Break-Test) |
| S2 | `nacharbeit-administration.sql` | `set_route` aus dem `REVOKE … FROM PUBLIC` | rot (Rollen-Test) |
| S3 | `nacharbeit-administration.sql` | `remove_route` aus der CHECK-Menge | rot (sechs Tests) |
| S4 | `nacharbeit-administration.sql` | `set_route`: `clock_timestamp()` → `now()` | rot (Aufruf-Ordnung) |
| S5 | `nacharbeit-administration.sql` | `set_route` ohne `SET search_path` | rot mit dem Filter, der den Katalog-Test einschließt (mit meinem ersten, zu engen Filter grün — Filter-Fehler, nicht Test-Lücke) |
| S6 | `nacharbeit-administration.sql` | `set_route` ohne `SECURITY DEFINER` | rot (Rollen-Test) |
| S7 | `nacharbeit-administration.sql` | `remove_route` schreibt `rule_name` NULL | rot |

Die vom Implementer gemeldeten Zahlen (23 Unit-, 10 Store- und zwei weitere
Mutationen) habe ich nicht einzeln nachgefahren; sie sind **übernommen**. Von
den benannten Grenzen ohne gesehene Mutation sind der `search_path`-Pin und
`SECURITY DEFINER` durch S5/S6 nun rot belegt, der Guard-Eintrag durch G1/G2
im Unit-Test und durch den realen Alt-Tag-Lauf (Lauf 2 gegen den gemessenen
Report); offen bleiben `UTF-8/NUL in when.column` (F-8) und die
API-Aktivierung (F-4).

---

## Findings

<!-- Kein Fließtext, kein Lösungsvorschlag im Befund. -->

### F-1 — Plan-Satz „`1e1` endet `rule_spec ist ungültig`; `jsonb` bewahrt die Schreibweise" trägt nicht

- `kategorie`: HIGH
- `quelle`: Hard-Rule-Klasse „Beleg trägt seinen Satz nicht" (Reviewer-Skill);
  `AGENTS.md` §3.12 Instanz B
- `pfad`: `docs/plan/planning/in-progress/slice-routing-antragsweg.md:200-203`
  (Umsetzungsentscheidung „Schreibweise der Zahl"),
  `internal/domain/model/routespec.go:82` (Godoc: „ein Exponent, auch `10.0`,
  endet hier"), `internal/application/usecase/setroute/service_test.go:222`
  (Testfall „order mit Exponent")
- `befund`: Der Plan stützt „`10.0` und `1e1` enden `rule_spec ist ungültig`"
  auf die Aussage, `jsonb` bewahre die Schreibweise. Gemessen an PostgreSQL 18
  bewahrt `jsonb` den Bruch (`10.0` bleibt `10.0`), normalisiert den Exponenten
  aber (`1e1` und `1E+1` werden `10`): `1e1` erreicht Go über die Antrags-Queue
  als `10` und wird angenommen; der Testfall „order mit Exponent" belegt nur den
  direkten Text-Aufruf von `ParseRouteSpec`, einen Pfad, den kein Antrag
  erreicht.
- `verifizierbar`: ja — `select ('{"order": 1e1}')::json::jsonb::text` an der
  gepinnten PostgreSQL; ein Store-Test, der `1e1` durch `cdc.set_route` und
  `ListPending` schickt, färbte die Aussage rot (der Annahmemengen-Test
  `TestAdministrationRequestSetRouteAcceptanceSet` deckt nur `10.0`).
- `klasse`: Beleg trägt seinen Satz nicht

### F-2 — Schreibweise von `order` enger als `SPEC-032`; Setzung ohne Spec-Nachzug oder Architect-Entscheid, Aufschub-Adresse nimmt sie nicht an

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Setzung des Implementers an der Spec-Grenze);
  [`SPEC-032`](../../spec/pflichtenheft.md) Zeile `order` („JSON-Zahl, positive
  ganze Zahl"); [`SPEC-019`](../../spec/pflichtenheft.md) Fehlertext-Tabelle
  („`order` fehlt oder ist keine positive Ganzzahl"); Review-Probe „Aufschub
  deckt die Adresse den Gegenstand?"
- `pfad`: `internal/domain/model/routespec.go:162-171` (`jsonRouteOrder`),
  `docs/plan/planning/in-progress/slice-routing-antragsweg.md:196-203`
- `befund`: `10.0` ist eine JSON-Zahl und mathematisch eine positive ganze
  Zahl; der Antrag endet dennoch `failed` mit `rule_spec ist ungültig`. Die
  Fixrunde von `slice-routing-spec-nachzug` hat den Ausschluss von Bruchteil
  und Exponent ausdrücklich aus der Spec genommen („keinen Träger in den
  Entscheidungen und keinen genannten Grund"); die Umsetzung führt ihn im Code
  wieder ein, der Plan meldet ihn „an den Planner". Für die neue Setzung
  (nur Ziffernfolge) steht im Plan der Adresse `slice-routing-betriebsdoku`
  kein Treffer (`git grep -E '2147483647|Obergrenze|10\.0|Ganzzahl|Schreibweise'`
  auf `slice-routing-betriebsdoku.md` und `welle-routing.md`: 0), das
  Handbuch würde die Grenze also nicht nennen.
- `verifizierbar`: nein — Spec-Konformität ist Lese-Handlung; `git grep` auf
  die Adresse bestätigt nur die fehlende Aufnahme.
- `klasse`: Setzung enger als die Spec ohne Träger

### F-3 — Begründung der Obergrenze `MaxRouteOrder` nennt einen Verbraucher, den es nicht gibt

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.12 Instanz B; [`SPEC-032`](../../spec/pflichtenheft.md)
  (keine Obergrenze, die Fixrunde des Spec-Nachzugs delegiert sie an diesen
  Slice „mit ihrer Begründung")
- `pfad`: `internal/domain/model/routespec.go:15-20`,
  `docs/plan/planning/in-progress/slice-routing-antragsweg.md:196-199`
- `befund`: Als Grund der Grenze 2147483647 nennen Plan und Godoc, die `order`
  bleibe „in jedem Ganzzahltyp eines Verbrauchers darstellbar"; die `order` ist
  weder Feld einer Change noch einer Nachricht noch einer View, ein Verbraucher
  liest sie nirgends (`git grep` der Ausgabeseite: kein Treffer). Die Grenze
  als Umsetzungsentscheidung ist durch die Delegation der Spec gedeckt; ihr
  genannter Grund ist es nicht belegt.
- `verifizierbar`: nein
- `klasse`: Begründung ohne Beleg-Anker

### F-4 — Pfad „API-Aktivierung" trägt den Regelstand ungebunden

- `kategorie`: MEDIUM
- `quelle`: Hard-Rule-Klasse „Zusage ohne Bindung an ihre Eingabeseite"
  (Reviewer-Skill); [`SPEC-019`](../../spec/pflichtenheft.md) („Jeder Pfad, der
  eine Erfassungs-Bindung anlegt, trägt den abgeleiteten Stand mit");
  DoD-Punkt (C) des Plans
- `pfad`: `internal/bootstrap/assemblersync.go:97` (Weitergabe `e.routing`),
  `internal/bootstrap/wiring.go:992` (`routing: activation` der
  `enableTableAPI`), `internal/bootstrap/assemblersync_internal_test.go:55,100,134`
- `befund`: Alle drei Tests des Dekorators `enableTableWithAssemblerSync`
  übergeben `&fakeRoutingPort{}` mit leerem Regelstand; keine Mutation an der
  Weitergabe oder an der Feldbelegung in `Run` kann sie rot färben. Die
  gemeinsame Funktion `syncAssemblerAddBinding` ist über den Aktivierungs-Zweig
  gebunden (M7 rot), die Zusage „in jedem Pfad" ist es für die HTTP-/gRPC-
  Aktivierung nicht; `slice-routing-e2e` und `welle-routing` nennen diesen Pfad
  nicht (`git grep` auf `enableTableWithAssemblerSync`, „API-Aktivierung":
  kein Treffer).
- `verifizierbar`: ja — `make test` mit einem Dekorator-Test, dessen
  `fakeRoutingPort` eine Regel trägt, färbt Weitergabe und Belegung rot.
- `klasse`: Zusage ohne Bindung an ihre Eingabeseite

### F-5 — `applyAdministrationRequest`: Wiederholungs-Zusage von `set_route` weder markiert noch getestet

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 / §3.12 (Kommentar der Klasse Zusage);
  Reviewer-Skill „Kommentar trägt keine der Kommentar-Klassen" (Zusage)
- `pfad`: `internal/bootstrap/wiring.go:1787-1799`
- `befund`: Der Kommentar sagt zu, eine Wiederholung nach fehlgeschlagenem
  `applied`-Vermerk sei „folgenlos" und beschreibt eine Grenze bei einem
  dazwischen verarbeiteten Antrag, ohne beides als hergeleitet zu kennzeichnen.
  Das Nachfahren im Code trägt beide Sätze (R1 bis R5 lesen nur `applied`-Zeilen,
  `Assembler.SetRoute` ersetzt nach Namen); die Transformations-Antragsarten
  tragen dafür einen eigenen Test (`TestProcessAdministrationRequestsSetTransformationIsIdempotent`),
  `internal/bootstrap/routing_internal_test.go` hat keinen Wiederholungs-Test.
- `verifizierbar`: ja — ein Wiederholungs-Test für `set_route` und die
  Mutation „`SetRoute` hängt immer an".
- `klasse`: Zusage im Kommentar ohne Marker und ohne Test

### F-6 — Interpunktion und Umbruch nach Teilersetzung in drei Kommentarblöcken

- `kategorie`: LOW
- `quelle`: Maintainability (dieselbe Klasse wie F-4/F-6 der Vorgänger-Reviews
  zu den Spec-Nachzügen); `make fmt-check` fängt Kommentare nicht
- `pfad`: `internal/adapters/driven/postgresstorage/queries/queries.go:325-326`
  (die Aufzählung „SelectAppliedColumnRequests SelectAppliedTransformationRequests
  und …" ohne Komma), `queries.go:354-356` (eine Zeile mit rund 120 Zeichen
  nach der Ersetzung), `internal/domain/model/administrationrequest.go:81-89`
  (Gedankenstrich am Zeilenende und ein umgebrochener Klammerzusatz mitten im
  Satz)
- `befund`: Die Sätze sind lesbar und inhaltlich richtig, tragen aber die Spur
  einer Teilersetzung (fehlendes Komma, nicht umbrochene Zeilen).
- `verifizierbar`: nein
- `klasse`: Umbruch/Interpunktion nach Teilersetzung

### F-7 — Godoc des Aufruf-Ordnungs-Tests nennt „alle sieben Funktionen"

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.13 (Träger-Nachzug außerhalb des Plan-Suchmusters)
- `pfad`: `internal/adapters/driven/postgresstorage/administrationrequest_order_test.go:162`
- `befund`: Der Godoc von `TestAdministrationRequestSameTransactionCallsKeepCallOrder`
  sagt „alle sieben Funktionen", es gibt jetzt neun; die beiden neuen deckt der
  eigene Test `TestAdministrationRequestRoutingCallsKeepCallOrderInOneTransaction`
  (S4 rot). Das Muster „sieben Funktionen" des Plan-Suchlaufs trifft die Wendung
  „alle sieben" nicht.
- `verifizierbar`: ja — `git grep -n -E 'alle sieben Funktionen'`.
- `klasse`: Träger-Nachzug außerhalb des Suchmusters

### F-8 — `ParseRouteSpec`: Prüfungen von `when.column` redundant zu `NewRouteRule`; äquivalente und ungebundene Mutanten

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `internal/domain/model/routespec.go:129` (nicht leer, ohne U+0000),
  `internal/domain/model/route.go:89` (dieselbe Prüfung im Konstruktor), M3,
  M12–M14
- `befund`: Entfällt die Prüfung in `ParseRouteSpec`, endet der Antrag über
  `Build`/`NewRouteRule` mit demselben Text (`rule_spec ist ungültig`); M12–M14
  bleiben auf Modell- und Use-Case-Ebene grün, und die einzige Wirkung wäre
  die Reihenfolge gegenüber `Zielname ist ungültig` bei gleichzeitiger
  Verletzung beider, die kein Test setzt. M3 (R4 `>=`) ist äquivalent, weil R2
  vorangeht. U+0000 in `when.column` erreicht Go über `jsonb` nicht. Keine
  Erwartung an eine Aktion.
- `verifizierbar`: nein
- `klasse`: äquivalenter Mutant

### F-9 — `ADR-0137` schreibt `rule_spec jsonb`, `SPEC-019` und der Bestand `json`

- `kategorie`: INFO
- `quelle`: [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  Teilfrage 3; [`SPEC-019`](../../spec/pflichtenheft.md) („nimmt die Regelform
  wie `cdc.set_transformation` als `json`-Parameter an");
  [`ADR-0125`](../plan/adr/0125-transformationen-parametertyp-regelform-json.md)
- `pfad`: `tools/schema/nacharbeit-administration.sql` (Funktion `set_route`),
  `tools/schema/rolloutguard/guard.go`
- `befund`: Die Abweichung des Parametertyps von der Notation der ADR ist
  gerechtfertigt: die Spec nennt `json` ausdrücklich, die Wache kennt nur
  `in:json`, ein `jsonb`-Parameter ließe den Abbau im zweiten Rollout scheitern
  (G2 rot; Alt-Tag-Lauf 2 und 5 grün). Die `Accepted`-ADR bleibt unberührt
  (`AGENTS.md` §3.5).
- `verifizierbar`: ja — `bash tools/harness/run-schema-rollout-guard-test.sh`.
- `klasse`: ADR-Notation gegen Bestand

### F-10 — Umfang des Diffs

- `kategorie`: INFO
- `quelle`: Plan §4 (vorab benannte Rückführung)
- `pfad`: gesamter Diff
- `befund`: 40 Dateien, rund 2.500 der 3.065 Zeilen sind Tests; die Produktiv-
  und Skriptänderungen liegen bei rund 700 Zeilen und schneiden sich sauber in
  die Commits `2551e995` (Go, B und C) und `622cb55d` (Schema, A). Die in §4
  benannte Trennung von (A) ist nicht ausgelöst; der Diff war in diesem
  Zuschnitt prüfbar.
- `verifizierbar`: nein
- `klasse`: Umfang

## Negativbefunde

- geprüft, ohne Befund: `internal/application/usecase/setroute/`,
  `removeroute/`, `excludecolumn/` — Prüfreihenfolge, Adressen und Klartexte
  stimmen mit der Fehlertext-Tabelle von [`SPEC-019`](../../spec/pflichtenheft.md)
  überein; R3 in beide Richtungen (M1, M2, M9 rot); der Bestandspfad von
  `exclude_column` ohne Routing-Bedingung ist unverändert (Regressionstest
  `TestExcludeColumnChecksSourceColumn`, rot-grün gelesen), der Regelstand wird
  erst nach der Spaltenexistenz gelesen.
- geprüft, ohne Befund: `internal/domain/model/routespec.go` (außer F-1 bis F-3,
  F-8) — Reihenfolge der Formzeilen und des unbekannten Schlüssels, `FoldRoutes`
  (Ersatz nach Namen, Entfernen, frische Liste), R1/R2/R4/R5 (M4, M5, M8 rot).
- geprüft, ohne Befund: `internal/adapters/driven/postgresstorage/`
  (`queries.go`, `sqlexec/translate.go`, `tableactivation.go`) — Ordnung
  `requested_at`, dann Antrags-ID (S1 rot), `ReadRoutingRules` faltet je
  Tabelle, Fehler eines unlesbaren Satzes endet sichtbar; außer F-6.
- geprüft, ohne Befund: `internal/bootstrap/wiring.go`, `assemblersync.go` —
  Prozessstart (M6 rot), Aktivierungs-Zweig (M7 rot), `applyAdministrationRequest`
  (M10 rot), `processedAdministrationKinds` leitet sich aus
  `AdministrationRequestKinds()` ab; Persist-before-ACK und Queue-Semantik
  (Lesefehler enden `failed`, kein Anhalten der Queue) sind unberührt; außer
  F-4 und F-5.
- geprüft, ohne Befund: `tools/schema/` (`nacharbeit-administration.sql`,
  `nacharbeit-roles.sql`, `schema.yaml`, `rolloutguard/`) — Funktionen
  `SECURITY DEFINER` mit `search_path`-Pin (S5, S6 rot), `REVOKE … FROM PUBLIC`
  und `GRANT` allein an `cdc_admin` (S2 rot, Alt-Tag-Lauf), CHECK mit neun
  Werten (S3 rot), `clock_timestamp()` (S4 rot), elf Fremdobjekte (G1, G2 rot,
  Alt-Tag-Lauf 2 grün); `plan.yaml` und `down.sql` unverändert.
- geprüft, ohne Befund: `tools/harness/run-schema-rollout-guard-test.sh`,
  `run-integration-tests.sh` — hart verdrahtete Liste der Arten (neun),
  Vorbedingung „Funktionen fehlen im Stand des Tags", Rechte unter den drei
  Rollen, Abschlusszeile; Lauf real grün.
- geprüft, ohne Befund: `harness/README.md`, `harness/targets/schema-rollout.md`
  — Zählwörter (sieben → neun Funktionen/Arten, neun → elf Fremdobjekte)
  nachgezogen; im ganzen Baum außerhalb der Records bleibt „sieben übrigen
  Antragsarten" in `queries.go:354` (sieben ist dort richtig: 9 − 2) und
  `spec/pflichtenheft.md:732/:734` (richtig: 9 − 2); Handbuch bleibt unberührt
  mit Aufschub-Adresse `slice-routing-betriebsdoku`, deren §2 den Gegenstand
  der Betreiber-Oberfläche trägt (`git grep -c -E
  'cdc\.set_route|rule_spec|Rollen-Tabelle'` → 8); ausgenommen die
  `order`-Setzung (F-2).
- geprüft, ohne Befund: Kommentare des Diffs — `make kommentar-kennungen` ohne
  Kandidat, keine Slice-/Wellen-Chronik in Produktionscode, keine
  Vorher/Nachher-Sprache; Hard Rule §3.1: `excludecolumn/service.go` entspricht
  im Endstand dem Plan (nur die beabsichtigten Hunks gegenüber `30fd6cb5`, 24
  Zeilen, Test 113); keine Datei ohne Zeilenende-Marker, `make fmt-check` grün.
  Ob eine Repo-Datei auf anderem Weg per Umleitung entstand, ist aus dem
  Endstand nicht ablesbar (`verifizierbar: nein`); ein Hinweis darauf liegt
  nicht vor. Die Selbstauskunft des Implementers zu `excludecolumn/service.go`
  (einmal geleert, per `git checkout` wiederhergestellt) ist **übernommen**.
- geprüft, ohne Befund: Commit-Traceability — alle vier Commits nennen
  [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  und [`LH-FA-CFG-008`](../../spec/lastenheft.md), kein Betreff trägt eine
  Struktur-Kennung.

## Beantwortung der Prüfauftrags-Punkte

**(a) Treue zu [`SPEC-019`](../../spec/pflichtenheft.md) und
[`SPEC-032`](../../spec/pflichtenheft.md).** Fehlertexte wörtlich, Reihenfolge
und Adressen stimmen, R1–R6 je eine Verletzung (Use-Case-Tabellentest mit
Eingabe-Bindung, M1/M4/M5/M9 rot); R3 in beide Richtungen; Antrag gegen eine
Tabelle ohne Bindung endet `applied`. Die Abweichung ist die Schreibweise von
`order` (F-1, F-2).

**(b) Setzungen.** `MaxRouteOrder = 2147483647` ist als Umsetzungsentscheidung
durch die Delegation der Spec-Fixrunde gedeckt, ihre Begründung trägt nicht
(F-3); die Beschränkung auf eine Ziffernfolge ist eine neue Einschränkung gegen
den Wortlaut der Spec, die der Spec-Nachzug bewusst entfernt hat, und ihr Beleg
trägt nur zur Hälfte (F-1, F-2) — Architect-Frage A-1. Konsistent mit dem
Slice-Plan, der die Grenze als „gemeldet" führt, nicht mit der Fixrunde des
Spec-Nachzugs.

**(c) Verdrahtung.** Zweige, Ableitung, Aktivierungs-Zweig und Prozessstart
gebunden (M6, M7, M10 rot). Die Lücke „API-Aktivierung" bewertet F-4 als
MEDIUM: gemeinsame Funktion gebunden, die Weitergabe im Dekorator und die
Belegung in `Run` nicht.

**(d) Schema.** Parametertyp `json` ist durch die Spec gedeckt (F-9); Rechte,
Pin, CHECK, Wache und Alt-Tag-Skript sind rot-grün belegt.

**(e) §3.1.** Siehe Negativbefund; kein Hinweis auf einen Verstoß im Endstand.

**(f) Zusagen-Klausel.** Ein Befund: F-5. Der Satz zur Wiederholungs-Grenze
trägt im Code, ist aber nicht als hergeleitet gekennzeichnet — entgegen der
Annahme im Prüfauftrag.

**(g) Mutationen.** Siehe Tabelle; 23 selbst gefahren, 19 rot, 4 grün
(vier äquivalente beziehungsweise redundante Mutanten, F-8); der
Implementer-Bestand ist übernommen.

**(h) Umfang.** F-10.

**(i) Sensoren und Zählwörter.** Siehe Eingangsprüfungen und Negativbefund.

## Architect-Fragen

- **A-1 (zu F-1 bis F-3).** Gilt für `order` die Annahmemenge „positive ganze
  Zahl" nach dem Wortlaut von [`SPEC-032`](../../spec/pflichtenheft.md) —
  also `10.0` angenommen — oder soll die Spec die Ziffernfolge und die
  Obergrenze 2147483647 festlegen? Welche der beiden Aussagen trägt den
  Fehlertext `rule_spec ist ungültig` für `order`? Der Entscheid bestimmt, ob
  der Code (`jsonRouteOrder`) oder die Spec nachgezogen wird und welche
  Aussage `slice-routing-betriebsdoku` ins Handbuch übernimmt.
- **A-2 (zu F-4).** Trägt `slice-routing-e2e` die HTTP-/gRPC-Aktivierung mit
  Regelstand, oder bleibt der Dekorator-Test die einzige Bindung dieses Pfads?

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 2 |
| LOW | 4 |
| INFO | 3 |

**Finding-Klassen dieses Laufs:** Beleg trägt seinen Satz nicht · Setzung enger
als die Spec ohne Träger · Zusage ohne Bindung an ihre Eingabeseite · Zusage im
Kommentar ohne Marker und ohne Test · Umbruch/Interpunktion nach Teilersetzung ·
Träger-Nachzug außerhalb des Suchmusters · Begründung ohne Beleg-Anker ·
äquivalenter Mutant

## Verdikt

**Merge-blockierend:** ja — ein HIGH (F-1) und zwei MEDIUM (F-2, F-4) sind
offen; die DoD-Zeile „Review durchgeführt … kein offenes HIGH/MEDIUM" bleibt
bis zur Fixrunde ungezogen. Der Code selbst ist belastbar (Spec-Texte,
Ordnung, R3 in beide Richtungen, Rechte, Wache und Alt-Tag-Lauf rot-grün
belegt); die Blockade liegt an einer Belegaussage im Plan und einer
Spec-Grenze, die der Architect entscheiden muss, und an einem ungebundenen
Pfad.

**Übergabe:** Findings F-1, F-4 bis F-7 gehen an den Implementer; F-2 und F-3
gehen über den Architect (A-1) und danach an Planner (Spec-/Plan-Nachzug) und
Implementer; die **Finding-Klassen** gehen zusätzlich in die Slice-Closure §7
und von dort in den Zähler. Dieser Report ist ein **Lauf-Beleg** und ersetzt
keine Verifikation — DoD-/Spec-Konformität prüft der Verifier separat
(Modul 11).
