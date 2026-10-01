# Review-Report: slice-routing-betriebsdoku — 2026-10-01

**Review-Art:** Code (Dokumentations-Diff) — geprüft gegen Plan, ADRs, Spec und `AGENTS.md` Hard Rules (Modul 10), mit
Schwerpunkt auf [`AGENTS.md`](../../AGENTS.md) §3.12 Instanz A und B. Kein DoD-Abgleich (Verifier).

**Gegenstand:** Slice `routing-betriebsdoku` der Welle [welle-routing](../plan/planning/done/welle-routing.md), Diff-Range
`2e68966f~1..HEAD` (`HEAD` = `be178eb6`), begrenzt auf `docs/user/benutzerhandbuch.md` und den Slice-Plan
[`slice-routing-betriebsdoku`](../plan/planning/done/slice-routing-betriebsdoku.md); Commits `758ae7a3` (Plan,
Konkretisierung und Suchlauf-Feld) und `be178eb6` (Handbuch 1.84). Reine Dokumentation, kein Code.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09“ (seither um weitere Klassen ergänzt).
**Modell:** claude-sonnet-5-5 · **Datum:** 2026-10-01.

**Ablage:** Der Reviewer-Lauf hat diesen Report selbst geschrieben (Write-Werkzeug; kein Edit-Werkzeug im Lauf). Die
Wegwerf-Umgebung lag im Scratchpad (eigene Compose-Datei mit eigenen Container- und Netznamen, `docker compose … down -v`
am Ende); Mutationen am Handbuch liefen als `sed … Datei > Kopie` im Scratchpad. Es wurde kein Image gebaut (`:dev` ist
das vorhandene Image), die Repo-Dateien blieben außer diesem Report unverändert (`tools/schema/plan.yaml`, vom
Schema-Rollout geschrieben, wurde mit `git checkout` zurückgenommen). Zeilenangaben beziehen sich auf
`docs/user/benutzerhandbuch.md` am `HEAD`.

**Eingangs-Kontext:**

- Slice-Plan `routing-betriebsdoku` (§1 Abgrenzung, §2 DoD und Übergabe-Block, §3 Konkretisierung und Suchlauf-Feld)
- [`ADR-0137`](../plan/adr/0137-routing-zustellziele-persistiertes-ziel-label.md),
  [`ADR-0138`](../plan/adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md),
  [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md),
  [`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md),
  [`ADR-0141`](../plan/adr/0141-run-regelstand-lesefehler-klasse-der-ursache.md),
  [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) (Teilfrage 4),
  [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- [`LH-FA-CFG-008`](../../spec/lastenheft.md), [`LH-FA-SST-006`](../../spec/lastenheft.md);
  [`SPEC-019`](../../spec/pflichtenheft.md), [`SPEC-020`](../../spec/pflichtenheft.md),
  [`SPEC-021`](../../spec/pflichtenheft.md), [`SPEC-022`](../../spec/pflichtenheft.md),
  [`SPEC-024`](../../spec/pflichtenheft.md), [`SPEC-031`](../../spec/pflichtenheft.md),
  [`SPEC-032`](../../spec/pflichtenheft.md)
- [`AGENTS.md`](../../AGENTS.md) (§3.7, §3.9, §3.12, §3.13), [`harness/conventions.md`](../../harness/conventions.md)
- Berichte der Quell-Slices: [`review-slice-routing-nats-subjekt`](review-slice-routing-nats-subjekt.md),
  [`verifikation-slice-routing-nats-subjekt`](verifikation-slice-routing-nats-subjekt.md),
  [`review-slice-routing-e2e`](review-slice-routing-e2e.md),
  [`verifikation-slice-routing-e2e`](verifikation-slice-routing-e2e.md),
  [`slice-routing-lesewege`](../plan/planning/done/slice-routing-lesewege.md) (§6 Altserver),
  [`slice-routing-e2e`](../plan/planning/done/slice-routing-e2e.md)

**Eigene Messungen** (Exit-Codes ungefiltert gesichert, je als eigener Schritt ausgewertet):

- `make docs-check` Exit 0, gedruckt: „d-check: 1514 Datei(en) geprüft, 0 Befund(e)“.
- `make suchlauf-nachmessen PLAN=docs/plan/planning/in-progress/slice-routing-betriebsdoku.md` Exit 0, gedruckt:
  „suchlauf-nachmessen: 16 Zeilen stimmen“ (alle 16 Zeilen `OK`, Soll gleich Ist; der Stand `diff` ist der saubere
  Arbeitsbaum gleich `HEAD`).
- `make gates` Exit 0 (ungefiltert in eine Datei, Exit-Datei gesondert gelesen; letzte Zeilen:
  „sdk-public-doc-check: keine interne Kennung unter sdks“, „gesamt: 0 Befund(e)“ der Architekturprüfung).
- `make test-notify` Exit 0 (selbst gefahren, siehe (e)): `TestRealServerRouteSubjects` PASS, gedruckt
  „natsstream-kosten: je 10000 Changes, Median von 5 Läufen: ohne Ziel 18.399382ms (543497 Changes/s), mit Ziel 19.436453ms
  (514497 Changes/s)“ und „(abgeleitet): mit/ohne = 1.06“.
- Image-Stand: `ghcr.io/pt9912/pg-change-feed:dev` erzeugt 2026-10-01 18:00:12; letzter Commit unter `internal`, `cmd`,
  `proto`, `gen`: `d8371372` um 17:34:11; `git diff --stat d8371372 HEAD -- internal cmd proto gen` ist leer. Das Image trägt
  den Code des Handbuch-Stands.

**Wegwerf-Umgebung (Wahrheitsprobe, PostgreSQL 18.6, Image `:dev`, NATS):** Schema über `make schema-rollout`, Login-Rollen
`rv_admin` (Mitglied `cdc_admin`) und `rv_reader` (Mitglied `cdc_reader`), alle Aufrufe von `cdc.set_route` unter
`rv_admin` (`select current_user, pg_has_role('cdc_admin','member')` gedruckt `rv_admin|t`), alle Lesezugriffe unter
`rv_reader`, HTTP mit dem Token der Klasse `reader`. Der Feed-Container selbst lief mit der Superuser-Identität (mit den
Login-Rollen scheiterte der Start an „Tabelle existiert nicht an der Quelle“, weil die Katalogabfrage die Spalten ohne
`SELECT`-Recht nicht sieht; die Rollen-Trennung der Feed-Verdrahtung ist nicht Gegenstand dieses Slice, `compose.yaml` fährt
ebenfalls mit dem Superuser).

| Probe | Handbuch-Aussage | Gedrucktes Ergebnis | Urteil |
|---|---|---|---|
| Beispiel: Zeilen 1 bis 3 vor den Regeln, `eu_orders` (10) und `rest` (100), danach Zeilen 4 bis 6; Abfrage `new_data->>'id'`, `region`, `route_target` | Tabelle des Beispiels | `eu_orders -> applied`, `rest -> applied`; `1 eu NULL`, `2 us NULL`, `3 NULL NULL`, `4 eu eu`, `5 us sonstige`, `6 NULL sonstige` | stimmt, Zeile für Zeile |
| `WHERE route_target = 'eu'` / `IS NULL` | „liefert die Zeile 4“ / Changes ohne Ziel | `4` / `1,2,3` | stimmt |
| `::jsonb`-Aufruf | `function cdc.set_route(unknown, unknown, unknown, unknown, jsonb) does not exist` | wörtlich gleich | stimmt |
| R1 | `Regelname bereits vergeben: public.orders.eu_orders` | `failed`, wörtlich gleich | stimmt |
| R2 | `order bereits vergeben: public.orders.10` | wörtlich gleich | stimmt |
| R3, Spalte fehlt | `Spalte existiert nicht an der Quelle: public.orders.regio` | wörtlich gleich | stimmt |
| R3, ausgeschlossen (nach `cdc.exclude_column` auf `name`) | `Spalte ist ausgeschlossen: public.orders.name` | wörtlich gleich | stimmt |
| R3 umgekehrt (`exclude_column` auf `region`) | `Spalte trägt eine Routing-Bedingung: public.orders.region` | wörtlich gleich | stimmt |
| R4, zweite Abschlussregel | `Regel ohne when bereits vorhanden: public.orders.rest` | wörtlich gleich | stimmt |
| R4, Abschlussregel nicht höchste (`rest` entfernt, Regel mit `order` 2147483647 gesetzt, `rest` mit 100 gesetzt) | `Regel ohne when trägt nicht die höchste order: public.orders.rest` | wörtlich gleich | stimmt |
| R4, Regel mit `when` hinter der Abschlussregel (`order` 150) | `order liegt hinter der Regel ohne when: public.orders.rest` | wörtlich gleich | stimmt |
| R5 | `Bedingung bereits vergeben: public.orders.region` | wörtlich gleich | stimmt |
| R6 | `Regelname nicht geführt: public.orders.nope` | wörtlich gleich | stimmt |
| Formprüfungen | `rule_spec ist ungültig`, `unbekannter Schlüssel in rule_spec`, `Zielname ist ungültig` | `rule_spec ist ungültig: public.orders.bad_o` (fehlendes `target`), `unbekannter Schlüssel in rule_spec: foo` (auch innerhalb `when`), `Zielname ist ungültig: public.orders.A` | stimmt (Adresse des unbekannten Schlüssels: F-4) |
| `order` 2147483647 / 2147483648 | angenommen / `rule_spec ist ungültig` | `applied` / `failed`, `rule_spec ist ungültig: public.orders.o_over` | stimmt |
| `order` `10.0`, `1e1` | als 10 gelesen | beide `failed`, `order bereits vergeben: public.orders.10` | stimmt |
| `order` `5.0`, `1e0` | angenommen | beide `applied` | stimmt |
| `order` 0, -3, 0.5, `"7"` (Zeichenkette) | `rule_spec ist ungültig` | alle `failed`, `rule_spec ist ungültig` | stimmt |
| Namensraum der Transformationsregeln getrennt | „von dem der Transformationsregeln getrennt“ | `set_transformation` mit dem Namen `id_rule` neben einer gleichnamigen Routing-Regel `applied` | stimmt |
| Bedingung liest den Quellwert vor der Transformation | „eine `rename_column`-Regel ändert das Ergebnis nicht“ | Bild `{"id": "11", "area": "eu"}` (Spalte umbenannt), Ziel `reg` (Bedingung auf `region`) | stimmt |
| Antrag bei stehendem Prozess | „`pending`, solange der Prozess steht“ | `pending`, `pending`; nach `docker restart` `applied`, `applied` | stimmt |
| `DELETE` ohne volle Replica-Identität, Zeile 5 (Region `us`) | Alt-Bild `{"id": "5"}`, Ziel `sonstige` | `DELETE | {"id": "5"} | sonstige` | stimmt |
| Bedingung auf die Schlüsselspalte beim `DELETE` (Handbuch: *hergeleitet*) | trifft | Regel auf `id` gleich `9`: `INSERT … idt`, `DELETE old_data {"id": "9"} … idt` | die Herleitung hält; vom Reviewer am System gefahren |
| Backfill-Run mit den Regeln des Beispiels, Zeile 5 gelöscht | fünf Backfill-Changes `eu`/`sonstige`, WAL-Changes derselben Zeilen unverändert | `completed`, `rows_copied 5`; `backfill INSERT 1 eu eu`, `2 us sonstige`, `3 NULL sonstige`, `4 eu eu`, `6 NULL sonstige`; WAL-Zeilen 1 bis 3 ohne Ziel | stimmt |
| `GET /changes?source=…&target=eu` | die eine `eu`-Change | `200`, eine Change (Zeile 4) | stimmt |
| `?target=Gross` | `200` mit `{"changes":[]}` | `200`, `{"changes":[]}` | stimmt |
| `?target=Gross&limit=0` | `400` | `HTTP/1.1 400 Bad Request` | stimmt |
| `?target=eu&schema=other` / `schema=public` | leer / die Change (Konjunktion) | `200 {"changes":[]}` / `200` mit der Change | stimmt |
| `?target=` (leer) und ein unbekannter Parameter | kein Filter / `400` | `200` ungefiltert / `400` | stimmt |
| Spalte `route_target` in `cdc.changes` | letzte Spalte, `origin` davor | `\d cdc.changes`: `origin`, danach `route_target` | stimmt |
| Zusatz-Subjekt `eu-west_1`, Wildcards, Kosten | Tests in `make test-notify` | Exit 0, `TestRealServerRouteSubjects` PASS, Kosten-Verhältnis 1,06 | stimmt, Spanne bleibt 0,92 bis 1,30 |
| Publication mit Spaltenliste bei **bekannter** Spaltenform (`ALTER PUBLICATION … SET TABLE public.orders (id, name)`, danach eine Zeile; Regel auf `region` geführt) | Ursache 1 nennt die Publication mit Spaltenliste als Fall der nicht anwendbaren Regel | Log `Fehlerklasse schema: Relation-Änderung nicht sicher als Obermenge interpretierbar: public.orders`, `cdc.heartbeat` `error_class = schema`, Exit 1; nach `ALTER PUBLICATION … SET TABLE public.orders` und Neustart dasselbe (die erneut gelesene Relation fehlt eine bekannte Spalte); nach Entfernen **beider** Regeln und Neustart (beide Anträge `applied`) erneut Exit 1 mit demselben Text | weicht ab, F-2 |

**Nicht selbst gefahren (benannt, nicht übernommen):** die Phasen von `make test-integration` (Happy Path an gRPC, SSE und
NATS mit Ruhefenster von 15 s, zwei `docker restart`, API-Aktivierung, Erreichbarkeit (b) und (c) an PostgreSQL 17.11 und
18.6, `DELETE`-Messung an 17.11), der gRPC-`ReadChanges`-Weg (kein gRPC-Client in der Umgebung), das Verhalten eines
ausgelieferten Alt-Servers, unverändertes TOAST und generierte Spalten als Bedingungsspalte, die Unit-Test-gestützten
Backfill-Aussagen (Stand-Wechsel `configuration`, nicht anwendbare Regel im Run, Lesefehler). Ihre Markierung im Handbuch
steht unter (c).

## Findings

### F-1 — Ursprung „gemessen im Lauf von `make test-integration`“ ohne Lauf-Anker und ohne „übernommen“

- `kategorie`: LOW
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.12 Instanz A und B, [`ADR-0083`](../plan/adr/0083-herkunft-von-aussagen-in-traegern.md)
- `pfad`: `docs/user/benutzerhandbuch.md:502`–`507` (Neustart, API-Aktivierung), `:600`–`:604` und `:612`–`:620`
  (Nichtanwendbarkeit), `:670`–`:673` (Replay), `:702`–`:708` (`DELETE`), `:1643` und `:1822`–`:1823` („gemessen an der
  Protobuf-Bibliothek“), `:1646`–`:1649` (Happy Path gRPC), `:1930`–`:1933` (Happy Path SSE), `:2424` (Grenzwerte:
  „gemessen: `2147483647` wird angenommen“)
- `befund`: Werte und Aussagen, die der Handbuch-Zug nicht selbst gefahren hat (E2E-Phasen, Protobuf-Bibliothek, Obergrenze
  von `order`), stehen als „gemessen im Lauf von `make test-integration`“ bzw. „gemessen“ ohne Lauf-Kennung, Datum oder
  Bericht und ohne den Vermerk „übernommen“; der Phasen-Anker auf `docs/user/e2e-abdeckung.md` (`:502`–`:507`) zeigt auf ein
  Erzeugnis, das `harness/README.md` (Zeile `make test-integration`) ausdrücklich „kein Lauf-Beleg“ nennt. Der Nachbarbestand
  desselben Handbuchs nennt bei Läufen die Lauf-Kennung (z. B. „Lauf `20260925T032925Z` von `tools/bench-backfill.sh`“) oder
  verlinkt den Bericht. Jede nachgeprüfte Aussage ist wahr (Tabelle oben; die Replay-Zeile steht wörtlich in beiden
  E2E-Berichten, die Phasen-Namen stehen im Runner); beanstandet ist die Form des Ursprungs.
- `verifizierbar`: nein (Lese-Handlung; die Wahrheit der Aussagen ist durch die Proben oben und die Berichte belegt)
- `klasse`: Ursprung ohne Lauf-Anker (übernommene Messung als „gemessen“)

### F-2 — Fall „Publication mit Spaltenliste“ bei bekannter Spaltenform führt in die Lage der Ursache 2; deren Überschrift nennt nur die entfernte Spalte

- `kategorie`: LOW
- `quelle`: [`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md) Entscheidung 2 und 3,
  [`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B
- `pfad`: `docs/user/benutzerhandbuch.md:596`–`:620` (Ursache 1 und Lösung), `:622`–`:636` (Ursache 2)
- `befund`: Ursache 1 führt „eine Publication mit Spaltenliste, die die Bedingungsspalte nicht enthält“ als Fall der nicht
  anwendbaren Regel auf, mit der Abhilfe „Regel entfernen, Neustart“. Am Wegwerf-System endete der Prozess bei einer
  Spaltenliste auf eine Tabelle mit **bekannter** Spaltenform als inkompatible Schemaänderung (Text der Ursache 2), und das
  Entfernen beider Regeln mit Neustart half nicht (Exit 1, derselbe Text). Die einleitende Bedingung von Ursache 1 („ohne
  dass die bekannte Spaltenform eine Entfernung zeigt“) ordnet den Fall nach der Lesart zu Ursache 2, deren Überschrift aber
  nur „die Spalte wurde an der Quelle entfernt“ nennt. Die E2E-Messung (c) lief an einer Tabelle ohne Spaltenform.
- `verifizierbar`: nein (eine Phase von `make test-integration` mit einer Spaltenliste an einer Tabelle mit bekannter
  Spaltenform würde es bestätigen; keine vorhanden)
- `klasse`: Ursache-Zuordnung unvollständig

### F-3 — Transformations-Fehlerblock übernimmt den Fall „Spaltenform noch nicht bekannt“ aus der Routing-Messung ohne Ursprung

- `kategorie`: LOW
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.12 Instanz B, [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 4
- `pfad`: `docs/user/benutzerhandbuch.md:389`–`:402`
- `befund`: Der Satz „oder ihre `column` fehlt in der Relation einer Tabelle, deren Spaltenform noch nicht bekannt ist“
  beschreibt für Transformationsregeln einen Fall, den nur die Routing-Phase von `make test-integration` am System gemessen hat;
  [`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 4 nennt ihn nicht. Der Satz
  trägt keinen Ursprung (weder gemessen noch hergeleitet noch „übertragen von Routing“); nur der Teil zur entfernten Spalte ist
  als *hergeleitet* markiert.
- `verifizierbar`: nein
- `klasse`: Tatsachenbehauptung ohne Ursprung (Übertragung zwischen Regelarten)

### F-4 — Adressliste der Fehlertexte nennt die Adresse des unbekannten Schlüssels nicht

- `kategorie`: INFO
- `quelle`: [`SPEC-019`](../../spec/pflichtenheft.md)
- `pfad`: `docs/user/benutzerhandbuch.md:557`–`:563` (Einleitung), `:579`–`:581` (Formprüfungen)
- `befund`: Der Einleitungssatz zählt als Adressen `schema.tabelle.regelname`, `…spalte`, `…zielname` und bei `order`
  `…<Wert>` auf; die Zeile `unbekannter Schlüssel in rule_spec` trägt nach Spec und Messung den bloßen Schlüsselnamen
  (`unbekannter Schlüssel in rule_spec: foo`). Die Formprüfungen stehen ohne Adresse.
- `verifizierbar`: nein
- `klasse`: Aufzählung unvollständig

### F-5 — Test- und Dateinamen in der Betreiber-Prosa

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `docs/user/benutzerhandbuch.md:646` (`routing_test.go`), `:1818`–`:1819`
  (`TestReadChangesTargetOutsideAlphabetAnswersEmptyWithoutStore`, `TestReadChangesContractErrorsWinOverInvalidTarget`),
  `:2095` (`publisher_nats_test.go`)
- `befund`: Der Elternstand führt Testfunktionsnamen nur in einer Zeile der Änderungshistorie (`git show 30fd6cb5:…`, ein
  Treffer auf `Test[A-Z]…`); der Text der Abschnitte nennt Phasen von `make test-integration` und Skripte, keine
  Unit-Testnamen oder Testdateien. Die neuen Stellen sind als Ursprung auflösbar (alle drei Testnamen und beide Dateien
  existieren) und tragen die Markierung „nicht am System gefahren“; sie sind eine neue Form im Bestand.
- `verifizierbar`: nein
- `klasse`: neue Ursprungsform im Handbuch-Text

## Prüfung nach den Schwerpunkten der Beauftragung

**(a) Wahrheit gegen Quellen.** R1 bis R6 (alle zehn Zeilen der Tabelle) und die Formprüfungen stimmen in Wortlaut, Adresse
und Reihenfolge mit [`SPEC-019`](../../spec/pflichtenheft.md) und mit dem System überein (Tabelle oben). Die `rule_spec`-Form
stimmt mit [`SPEC-032`](../../spec/pflichtenheft.md) überein (Schlüssel, Alphabet 1 bis 63 Zeichen, `when` mit beiden
Pflichtschlüsseln, zeichengenauer Vergleich, Quellwert vor Transformation, Abwesenheit, erster Treffer in aufsteigender
`order`). `order`: positive ganze Zahl; die Obergrenze 2147483647 steht als Setzung der Umsetzung (die Spec nennt keine:
`SPEC-032` führt „positive ganze Zahl“, im Code `MaxRouteOrder` in `internal/domain/model/routespec.go`), und jede JSON-Zahl
mit positivem ganzzahligem Wert wird angenommen (`10.0`, `1e1` am System als 10 gelesen, `5.0`, `1e0` angenommen; 0, -3, 0.5,
`"7"` und 2147483648 enden `rule_spec ist ungültig`). Abschlussregel (ohne `when`, höchste `order`): alle drei R4-Zeilen am
System gefahren; die Lösung („Abschlussregel entfernen, nach der neuen Regel mit der dann höchsten `order` neu setzen“) ist
bis zum Fehlschlag des Neusetzens mit zu kleiner `order` gefahren; das erfolgreiche Neusetzen mit höchster `order` folgt aus R4
und ist nicht gefahren. Lese-Kontrakt-Vorrang („bei sonst gültiger Anfrage“) stimmt mit [`SPEC-022`](../../spec/pflichtenheft.md)
(Zeile Zustellziel) überein; am HTTP-Weg gefahren (`limit=0` mit `target` gleich `400`). Den gRPC-Weg trägt das Handbuch als
Unit-Test-gestützt (beide Testnamen existieren in `internal/application/usecase/readchanges/service_test.go`).

**(b) Ausgeführte Beispiele.** Alle Beispiele und Ausgaben der Tabelle oben wurden unter `cdc_admin`- bzw.
`cdc_reader`-Login-Identitäten gefahren; der Poll auf `applied` ist gedruckt; die Ausgabe des Beispiels stimmt zeichengleich.

**(c) Ursprungs-Kennzeichnung.**

- *Gemessen im Lauf von `make test-integration`* (aus dem E2E-Slice): Neustart und API-Aktivierung (Phasen „Routing-Neustart
  und Ausschluss-Sperre“, „Routing-Aktivierung über die API“ — beide Namen stehen im Runner `tools/harness/run-integration-tests.sh`),
  Erreichbarkeit (b) und (c) samt Abhilfe (Phase „Routing-Nichtanwendbarkeit und Abhilfe“, Log-Sentinel im Runner vorhanden),
  Replay-Zeile (steht wörtlich in beiden E2E-Berichten), `DELETE` an PostgreSQL 17.11 und 18.6, Happy Path an gRPC und SSE mit
  15 s Ruhefenster (`RT_WINDOW=15s` im Runner). Keine dieser Aussagen tritt als im Handbuch selbst gemessen auf; die Form des
  Ursprungs ist F-1.
- *Am Handbuch-Stand selbst gefahren* (der Text sagt „in einer Compose-Umgebung“): Beispiel, Fehlertexte, `DELETE`-Alt-Bild,
  Backfill-Beispiel, `GET /changes`-Proben. Alle vom Reviewer reproduziert.
- *Unit-Test-gestützt, als „nicht am System gefahren“ markiert*: Stand-Wechsel im Run (`configuration`), nicht anwendbare Regel
  im Run (`schema`, vor der ersten Zeile), Lesefehler mit Klasse der Ursache (`storage`), gRPC-`ReadChanges` leerer Treffer und
  Vorrang. Die Tests `TestExecuteRoutingStateChangeEndsRunAsConfiguration`, `TestExecuteInapplicableRoutingRuleEndsRunAsSchema`,
  `TestExecuteRoutingReadFailureEndsRun` und die zwei `ReadChanges`-Tests existieren.
- *Hergeleitet, markiert*: Alt-Server (gRPC ignoriert das Feld, HTTP und SSE `400`), „Entfernen der Regel genügt bei entfernter
  Spalte nicht“, „Bedingung auf die Schlüsselspalte trifft beim `DELETE`“, unverändertes TOAST und generierte Spalten (Zusage der
  Spec, nicht gefahren), Wecksignal-Wurzel nur negativ belegt. Die Herleitung zum `DELETE` mit Schlüsselspalte hält (Probe oben).
  Zu „Entfernen genügt nicht“: für den Fall der entfernten Spalte nicht gefahren (der Text sagt es); der Fall Spaltenliste bei
  bekannter Spaltenform zeigte dasselbe Verhalten, siehe F-2.
- Die Aussage „proto3 verwirft unbekannte Felder; gemessen an der Protobuf-Bibliothek dieses Repositorys“ (`:1643`, `:1822`)
  ist **übernommen** aus [`slice-routing-lesewege`](../plan/planning/done/slice-routing-lesewege.md) §6 (die Tests
  `TestStreamChangesRequestIgnoriertUnbekannteFelder`, `TestReadChangesRequestTraegtTargetAlsFeldSieben`); vom Reviewer nicht
  nachgefahren, im Handbuch nicht als übernommen gekennzeichnet (F-1).

**(d) Zählwort-Nachzüge.** „Zwei optionale“ wurde an den vier Roh-Wegen ersetzt: HTTP (`:1455` „`schema`, `table` und
`target`“), gRPC-Stream (`:1630` „drei optionale, unabhängig setzbare Felder“), `ReadChanges` (`:1799`–`:1802`, Feldnummer 7
gleich `string target = 7` in `proto/cdc/administration/v1/administration.proto:137`), SSE (`:1918` „Drei optionale
Query-Parameter“). Der Stream-Request trägt `target` als Feldnummer 3 (`proto/cdc/stream/v1/changestream.proto:37`). An den
SDK-Abschnitten C#, Kotlin und Python (`:1708`, `:1721`, `:1746`) bleibt „zwei“: **bestätigt** — `git grep -n -i -w target --
sdks/csharp/PgChangeFeed.Client sdks/python/pgchangefeed/pgchangefeed sdks/kotlin/pgchangefeed-kotlin/src/main` findet nur zwei
Kommentarzeilen („JSON needs a target type“) in den Verwaltungs-Clients, keinen Parameter; die Signatur
`StreamChangesAsync(schema, table, cancellationToken)` und die Gegenstücke in Kotlin und Python tragen zwei Parameter. Die
Abweichung vom Wortlaut des Übergabe-Blocks ist wahr und im Plan (§3 Konkretisierung) begründet; der Satz „das Feld `target`
folgt mit dem Package“ ist der benannte Aufschub zu
[`slice-routing-sdk-beispiel-target`](../plan/planning/done/slice-routing-sdk-beispiel-target.md) (der Plan existiert in
`open/`, sein Titel nennt „der Parameter `target` an allen Zustellweg-Flächen“). Die Suchläufe im Plan (`zwei optionale` 4 am
Elternstand, 2 am Diff; `drei optionale` 2) stimmen nach `make suchlauf-nachmessen`; die C#-Stelle bricht zwischen „zwei“ und
„optionale“ um und fällt aus dem Zeilenmuster, der Plan sagt das.

**(e) NATS-Kostenabsatz (`:2101`ff.).** Sechs Werte im Handbuch, Spanne 0,92 bis 1,30; jeder Wert gegen die Quelle:

| Wert | Herkunft | Nachgerechnet / nachgelesen |
|---|---|---|
| 1,16 | Umsetzung, „ohne Ziel 18.242683ms, mit Ziel 21.239865ms“ | 21,239865 / 18,242683 = 1,164 |
| 1,30 | Umsetzung, „ohne Ziel 17.479304ms, mit Ziel 22.739539ms“ | 22,739539 / 17,479304 = 1,301 |
| 1,05 | Review, „ohne Ziel 20.132253ms, mit Ziel 21.076066ms“ | im Review-Bericht gedruckt „mit/ohne = 1.05“ |
| 1,14 | Fixrunde, „ohne Ziel 18.450212ms, mit Ziel 21.111558ms“ | 21,111558 / 18,450212 = 1,144 |
| 0,92 | Verifikation, „ohne Ziel 20.856066ms, mit Ziel 19.271615ms“ | im Verifikations-Bericht gedruckt „mit/ohne = 0.92“ |
| 1,02 | Handbuch-Lauf, „ohne Ziel 19.656453ms, mit Ziel 20.129363ms“ | 20,129363 / 19,656453 = 1,024; die Zeile steht nur im Handbuch, der Lauf ist nicht auflösbar |

Die Ursprungs-Kennzeichnung (gemessen und abgeleitet in diesem Zug, übernommen aus den Berichten) ist korrekt getrennt; „kein
Aufschlag auflösbar“ ist als Aussage mit der Spanne begründet, nicht als Zusage. Ein siebter Wert vom Reviewer: 1,06
(`make test-notify`, Exit 0), innerhalb der Spanne; der Hinweis V-3 der Verifikation (Spanne statt „etwa 1,05 bis 1,30“) ist
übernommen. Kein Befund. Für „gemessen“ trägt der Wert 1,02 die gedruckte Zeile im Text, wie §3.12 es verlangt.

**(f) Historie und Version.** Kopf `Version: 1.84` und `Stand: 2026-10-01`; die Zeile 1.84 (`:2722`) steht hinter der letzten
Zeile 1.83 (Tabellenende); Suchlauf-Zeilen `Version: 1.84` einmal und Reihe 1.8x fünf Zeilen stimmen. Keine Chronik-Sprache in
der neuen Prosa: gelesen und per Muster (`jetzt|zuvor|früher|bisher|vorher|nachher|wurde|…`) über alle hinzugefügten Zeilen
abgesucht; die Treffer sind Ist-Beschreibungen („wie zuvor jeden Change“ für den ungefilterten Altaufruf, eine Formel des
Bestands; „die zuvor nicht bestätigte Transaktion“ als Zustand). Die Kennungen `SPEC-0xx` (17 Zeilen am Elternstand), `ADR-0xxx`
und `LH-*` kommen im Nachbarbestand vor; neu im Handbuch ist `LH-FA-SCH-004` (eine Lastenheft-Kennung wie die übrigen). Glossar
(`:2400`ff., zwei Zeilen), Rollen-Tabelle (zwei Funktionen), Grenzwerte (`:2415`ff.: Alphabet, `order`, keine Obergrenze der
Regelzahl), §6 Fehlerklassen-Zeile `schema` (`:2283`), „Neustart nach einem Fehler“ (`:2352`ff.) und „Schema aktualisieren“
(`:1339`, `:1365`: zwei zusätzliche View-Spalten, `cdc.change.route_target`): gelesen, ohne Befund. Die Aussage der Grenzwerte
„die Spezifikation führt keine Obergrenze für die Zahl der Regeln“ ist eine Absenz; `git grep -n -E "Keine
Obergrenze|Obergrenze für die Zahl|Zahl der Regeln" -- spec/pflichtenheft.md` liefert keinen Treffer (die Suche ist eng; sie
ersetzt kein Lesen von `SPEC-019` und `SPEC-032`, die ich gelesen habe: keine Zahl).

**(g) Querverweise und Fehlerblock der Transformationsregeln.** `make docs-check` Exit 0 (Anker und Links), Suchlauf Exit 0.
Der Fehlerblock der Transformationen trennt die entfernte Spalte von der nicht anwendbaren Regel und sagt für die entfernte
Spalte nicht „Regel entfernen“ als Abhilfe; das stimmt mit
[`ADR-0112`](../plan/adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md) Teilfrage 4 („die spalten-entfernenden
Fälle enden dort schon vorher mit derselben Klasse“) und
[`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md) Entscheidung 2 überein; der
zusätzliche Fall „Spaltenform noch nicht bekannt“ ist F-3.

**(h) Vollständigkeit gegen den Übergabe-Block** (Plan §2, Zeile für Zeile mit `git grep` am Handbuch):

| Aufgeschobener Gegenstand | Fundstelle im Handbuch | Befund |
|---|---|---|
| `kern-label`: Spalte `route_target` in „Änderungen lesen“ und in den Zugriffswegen, `NULL` heißt „nicht geroutet“ | `:793`–`:813` (SELECT-Liste, Satz zur letzten Spalte), `:1465` (HTTP), `:1802`–`:1810` (`ReadChanges`-Feldliste) | getragen |
| `antragsweg`: `set_route`/`remove_route`, Rolle, `rule_spec`, R1–R6, Abhilfe, R4-Abschlussregel, `order` | Abschnitt „Routing-Regel konfigurieren“ (`:439`ff.), Rollen-Tabelle (Zeile `cdc_admin`, `:76`) | getragen, Wortlaut am System bestätigt |
| `backfill-pfad`: Aufzählung „Ausgeschlossene Spalten“ im Abschnitt „Bestand als Backfill überführen“ | neuer Punkt „Routing-Regeln“ (`:1057`–`:1066`) | getragen |
| `backfill-pfad`: Zeile `schema` der Fehlerklassen-Tabelle | `:2283` | getragen |
| `lesewege`: Parent-Zeilen `:1149`, `:1301`, `:1364`–`:1365`, `:1378`, `:1403`, `:1456`–`:1457`, `:1557` | jetzt `:1455`, `:1630`, `:1708`, `:1721`, `:1746`, `:1799`–`:1802`, `:1918` | je Stelle gelesen, getragen; SDK-Stellen bleiben „zwei“ mit Aufschub-Vermerk |
| `lesewege`: Alt-Server, Konjunktion, „Auswahl kein Zugriffsschutz“, „bei sonst gültiger Anfrage“ | HTTP `:1482`ff., gRPC-Stream, `ReadChanges`, SSE | getragen, als *hergeleitet* markiert |
| `nats-subjekt`: Zusatz-Subjekt, fire-and-forget, Kosten, Wecksignal | `:2084`ff. (Subjekt), `:2101`ff. (Kosten), `:2115`ff. (Zustellsemantik), `:2015`–`:2019` (Wecksignal) | getragen |
| `e2e`: `DELETE`-Messung (Alt-Bild, Abschlussregel, `REPLICA IDENTITY FULL`, 17.11 und 18.6) | Hinweis „DELETE ohne volle Replica-Identität“ (`:698`ff.) | getragen, Menge benannt |
| `e2e`: Erreichbarkeit (a) bis (c), Abhilfe, Replay/Backfill | Ursache 1 und 2 (`:596`ff.), Hinweise (`:664`ff.) | getragen (F-1 zur Form, F-2 zu (c)) |
| Aufschub: Flags der Beispiel-Clients, SDK-Abschnitte | `:1651`–`:1657`, `:1933`–`:1937` | benannte Zwischenzeit mit „folgt mit den SDK-Packages“; Adresse `slice-routing-sdk-beispiel-target` existiert in `open/` mit dem Gegenstand |

**(i) Gates.** `make docs-check` Exit 0, `make suchlauf-nachmessen` Exit 0, `make gates` Exit 0 (siehe oben).

**Mutations-Analog für Doku** (drei Tatsachenaussagen auf Kopien im Scratchpad, `sed … > Kopie`; das Handbuch selbst blieb
unverändert):

| Mutation | Probe, die sie entlarvt | Ergebnis der Probe am System |
|---|---|---|
| R1-Fehlertext `Regelname bereits vergeben` zu `Regelname existiert bereits` | R1-Aufruf `set_route` mit vorhandenem Namen, `error_message` lesen | gedruckt `Regelname bereits vergeben: public.orders.eu_orders` — die Mutation wäre falsch |
| Obergrenze `2147483647` zu `2147483648` | `order` 2147483648 setzen | `failed`, `rule_spec ist ungültig`; 2147483647 `applied` — die Mutation wäre falsch |
| `?target=Gross&limit=0` `400` zu `200` | `GET /changes?…&target=Gross&limit=0` | `HTTP/1.1 400 Bad Request` — die Mutation wäre falsch |

**Aussagen ohne mögliche Probe in diesem Lauf** (benannt): E2E-Messwerte an PostgreSQL 17.11 und die Phasen von `make
test-integration` (nicht gefahren; getragen durch die Berichte der E2E-Slices), Alt-Server-Verhalten (kein ausgelieferter
Alt-Server vorhanden), unverändertes TOAST und generierte Spalten, gRPC-`ReadChanges` am System, die Unit-Test-gestützten
Backfill-Aussagen (nur die Unit-Tests existieren), die Aussage „Entfernen der Regel genügt bei entfernter Spalte nicht“ für
`DROP COLUMN` (der Fall Spaltenliste zeigte das Verhalten; der Fall der entfernten Spalte selbst ist nicht gefahren, der Text
sagt *hergeleitet*), und die Absenz einer Obergrenze der Regelzahl in der Spec (nur Lesen).

## Negativbefunde

- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` Abschnitt „Routing-Regel konfigurieren“ — R1 bis R6, `rule_spec`-Form,
  `order`, Beispiel, `jsonb`-Fehler, Abschlussregel, Backfill-Beispiel, `DELETE` (am System nachgefahren; Ausnahmen F-1, F-2, F-4)
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` Zugriffswege (HTTP, gRPC-Stream, gRPC-Verwaltungs-API, SSE, NATS-Wecksignal
  und NATS-Vollinhalt) — Zählwörter, Feldnummern, Konjunktion, Vorrang, Subjekt-Form; Kosten-Absatz gegen die Berichte
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` Kopf, Änderungshistorie, Rollen-Tabelle, Glossar, Grenzwerte, §6
  Fehlerklassen, „Neustart nach einem Fehler“, „Schema aktualisieren“
- geprüft, ohne Befund: `docs/plan/planning/in-progress/slice-routing-betriebsdoku.md` — Suchlauf-Feld (16 Zeilen
  nachgemessen), Konkretisierung, Übergabe-Block (alle fünf Quell-Slices abgedeckt); die DoD-Zeilen stehen vor dem Review offen
- geprüft, ohne Befund: Rückbezüge außerhalb des Handbuchs — der Plan misst 0 Treffer für `route_target`/`cdc.set_route` in
  `README.md`, `examples`, `sdks`; eigene Stichprobe `git grep -n -i target -- examples sdks` findet nur Build-Dateien,
  Kommentare und Integrationstest-Konstanten (`RuleTargetKey`), keinen Aufrufparameter `target`

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Ursprung ohne Lauf-Anker (übernommene Messung als „gemessen“) · Ursache-Zuordnung
unvollständig · Tatsachenbehauptung ohne Ursprung (Übertragung zwischen Regelarten) · Aufzählung unvollständig · neue
Ursprungsform im Handbuch-Text

## Fragen an den Architect

1. **Fall „Spaltenliste an einer Tabelle mit bekannter Spaltenform“ (F-2).**
   [`ADR-0140`](../plan/adr/0140-routing-nichtanwendbarkeit-erreichbarkeit-und-abhilfe-grenze.md) Entscheidung 3 (c) beschreibt
   die Publication mit Spaltenliste als Fall des Sicherheitsnetzes (die erste Relation legt die Spaltenform an). Am System
   endet dieselbe Spaltenliste an einer Tabelle, deren Spaltenform schon angelegt ist, als inkompatible Schemaänderung, und das
   Entfernen der Regeln hilft nicht. Gehört dieser Fall zur Ursache 2 des Handbuchs (Überschrift erweitert auf „Relation verliert
   eine bekannte Spalte“), oder will die ADR ihn selbst nennen? Eine neue ADR wäre nur nötig, wenn die Abhilfe sich ändern soll;
   die Frage betrifft zunächst die Zuordnung im Handbuch.
2. **Form des Ursprungs „gemessen im Lauf von `make test-integration`“ (F-1).** Soll das Handbuch bei übernommenen E2E-Werten
   den Bericht ([`review-slice-routing-e2e`](review-slice-routing-e2e.md),
   [`verifikation-slice-routing-e2e`](verifikation-slice-routing-e2e.md)) verlinken, wie der Nachbarbestand es für andere
   Messläufe tut, oder genügt die Phase als Anker? Der Plan schreibt die Formel vor; die Frage ist, ob die Formel §3.12 genügt.

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM. Die drei LOW- und zwei INFO-Findings betreffen die Form des Ursprungs, die
Zuordnung eines Randfalls zur Ursache und eine unvollständige Aufzählung. Keine Aussage des Handbuchs war gegen Spec oder System
falsch; in F-2 hilft die Abhilfe der Ursache 1 im gemessenen Randfall nicht, der Text ordnet den Fall nach seiner einleitenden
Bedingung aber ohnehin der Ursache 2 zu (deshalb LOW).

**Übergabe:** Findings an den Implementer (F-1 bis F-3 als kurze Handbuch-Fixrunde empfohlen, kein Plan-Defekt); F-2 zusätzlich
als Architect-Frage 1; die Finding-Klassen gehen in die Slice-Closure §7. Die DoD-Zeile „Review durchgeführt“ im Slice-Plan
bleibt offen, weil eine Fixrunde empfohlen wird; wird keine gefahren, zieht der Planner die Zeile nach (der Lauf hatte kein
Edit-Werkzeug, der Plan ist 376 Zeilen lang und wurde nicht überschrieben). Dieser Report ersetzt keine Verifikation.
