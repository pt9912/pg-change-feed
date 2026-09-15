# ADR-0081: Changes-Lesen über die HTTP-API (Supers. ADR-0057, teilw.)

**Status:** Accepted — Supersedes [`ADR-0057`](0057-http-grpc-api.md) in genau
**zwei Klauseln**, jeweils in ihrer Changes-Lesen-Hälfte: dem Ausschluss-Satz
in §Entscheidung Teilfrage 2 („Changes-Lesen und Diagnose/Health bleiben damit
bewusst ausgeschlossen — nicht verworfen, sondern vertagt") und der
zugehörigen Negativ-Konsequenz in §Konsequenzen („Changes-Lesen
(`LH-FA-REA-*`) und Diagnose/Health bleiben nach dieser ADR weiterhin
CLI-/SQL-exklusiv"). Alles Übrige aus `ADR-0057` bleibt unverändert in Kraft:
die Protokoll-Wahl HTTP/JSON (Teilfrage 1), der Umfang der *ersten* Ausbaustufe
(Teilfrage 2, Option B) als Beschluss über diese Stufe, die zwei
Token-Rechtsklassen (Teilfrage 3), die Adapter-Platzierung (Teilfrage 4) und
die Diagnose/Health-Hälfte beider abgelösten Klauseln.

**Datum:** 2026-09-15

**Autor:** pt9912 (Rolleninhaber: Architect-Lauf, 2026-09-15)

**Bezug:** [`LH-FA-SST-006`](../../../spec/lastenheft.md) (Haupt-Bezug —
die API), [`LH-FA-REA-001`](../../../spec/lastenheft.md)…[`006`](../../../spec/lastenheft.md),
[`LH-FA-SST-007`](../../../spec/lastenheft.md),
[`LH-FA-SST-008`](../../../spec/lastenheft.md), [`ADR-0057`](0057-http-grpc-api.md)
(Teilfrage 2 Option C, die hier ausgeführt wird), [`ADR-0046`](0046-sql-driving-adapter-lese-schreib-trennung.md)
(Lese-View-Direktzugriff ohne Anwendungsdienst), [`ADR-0056`](0056-nats-tabellen-granulares-subjekt.md)
(Klartext-Achse statt opaker Kennung an einem Draht-Vertrag), [`ADR-0042`](0042-transport-typen-am-port.md)
(Transport-Typen am Port), [`ADR-0034`](0034-ports-nach-faehigkeiten.md),
[`ADR-0028`](0028-inbound-use-cases.md), [`ADR-0079`](0079-nats-beispielclient-vierter-examples-client.md)
(der Beispiel-Client, dessen HTTP-Abfrage diese ADR erst möglich macht)

**Schärft:** [`ARC-005`](../../../spec/architecture.md) (der HTTP-Driving-Adapter
trägt zusätzlich die lesende Fähigkeit; die Lese-Sequenz der Sicht-§4 mit
`ReadChangesUseCase`/`ChangeStorePort` gilt damit auch für den
Netzwerkzugriffsweg) — die konkrete Endpunkt- und Rückgabeform führt ein
**neuer** `SPEC-*`-Eintrag des Pflichtenhefts (Folgepflicht dieser ADR).

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

**Der Befund, gemessen.** `internal/adapters/driving/http/server.go`
registriert **zehn** Routen — neun Port-gedeckte Fähigkeiten und
`GET /changes/stream` ([`SPEC-021`](../../../spec/pflichtenheft.md)); ein
**`GET /changes` gibt es nicht**. Ein Netzwerk-Consumer kann Changes heute
ausschließlich über den SQL-View-Direktzugriff `cdc.changes` lesen.

**Der Vertrag, der das anders verlangt.**
[`LH-FA-SST-006`](../../../spec/lastenheft.md) (Rang 1) verlangt in seiner
**Beschreibung**: „die bestehenden Lese- und Verwaltungsfähigkeiten (u. a.
Consumer-Registrierung, Positions-Bestätigung, **Changes lesen**,
Status-/Diagnoseabfragen) zusätzlich zu den bestehenden CLI-/SQL-Zugriffswegen
über einen Netzwerkzugriffsweg bereitstellen". Seine **Akzeptanzkriterien**
sind demgegenüber exemplarisch („eine unterstützte Fähigkeit (**z. B.**
Consumer-Registrierung)") und tragen die Verengung nicht; seine
**Out-of-Scope**-Klausel grenzt allein Protokoll-, Authn- und
Endpunkt-Details aus. **Die Beschreibung bindet.**

**Wie die Verengung entstand — und wo sie steht.** [`ADR-0057`](0057-http-grpc-api.md)
hat für die *erste* API-Ausbaustufe Option B gewählt (nur die bereits
Port-gedeckten Fähigkeiten) und das Changes-Lesen „bewusst ausgeschlossen —
nicht verworfen, sondern vertagt"; deren Option C ist dort wörtlich
vorgezeichnet: „zusätzlich Changes lesen (`LH-FA-REA-*`) über einen neuen
`ReadChangesUseCase`/Inbound Port". Diese Verengung ist eine Entscheidung der
ADR- und Spezifikations-Ebene; [`SPEC-018`](../../../spec/pflichtenheft.md)
hat sie übernommen („Changes-Lesen (`LH-FA-REA-*`) und Diagnose/Health bleiben
außerhalb") — **an der Beschreibung des Lastenhefts ist sie nie angekommen**.
Zwei weitere Anforderungen setzen den API-Leseweg ohnehin voraus:
[`LH-FA-SST-007`](../../../spec/lastenheft.md) und
[`LH-FA-SST-008`](../../../spec/lastenheft.md) sprechen beide vom „bestehenden
Lesezugriffsweg (SQL/**API**)".

**Was den Zug ausgelöst hat.** [`ADR-0079`](0079-nats-beispielclient-vierter-examples-client.md)
(Accepted, 2026-09-15) legt einen öffentlichen Beispiel-Client fest, der das
NATS-Wecksignal lauscht und „die Änderung **real** über die HTTP-API
(`LH-FA-SST-006`)" holt. Diese Prämisse ist beim Umsetzungsversuch am
gemessenen Bestand gescheitert: Es gibt keinen Lese-Endpunkt. Der Bedarf ist
damit konkret benannt — genau die Bedingung, die `ADR-0057` in ihrem
Re-Evaluierungs-Trigger als Auslöser einer Folge-ADR festgeschrieben hat.

**Vier Constraints prägen den Lösungsraum.**

- **Der Leseport existiert bereits — und ist heute fast ungenutzt.**
  `internal/application/port/outbound/changestore.go` trägt
  `ChangeQuery{Source, Start, End, Table, Limit}` und
  `ReadChanges(ctx, query) ([]ChangeRecord, error)`; produktiver Aufrufer ist
  allein der Retention-Lauf
  (`internal/application/usecase/retention/service.go`, nur `Source`). Die
  Sicht-§4 zeichnet diese Sequenz für den CLI-Kanal bereits:
  `ReadChangesUseCase` → `ChangeStorePort` → `PostgresChangeStoreAdapter`.
- **`cdc.changes` ist bewusst ein View-Direktzugriff ohne Anwendungsdienst**
  ([`ADR-0046`](0046-sql-driving-adapter-lese-schreib-trennung.md)): Der
  SQL-Kanal projiziert direkt, weil eine reine Lese-View keine
  Entscheidungslogik trägt. Der Netzwerk-Leseweg darf deshalb **kein zweiter
  Lesepfad** mit eigener Logik werden; er ist eine zweite *Fassade* über
  demselben Port (`LH-FA-SST-006` Boundary: dieselbe Domänenlogik, kein
  Zweitpfad).
- **Der Lesepfad trägt Schema- und Tabellenname heute nicht.**
  `mapper.ToChange` (`internal/adapters/driven/postgresstorage/mapper/mapper.go`)
  rekonstruiert einen Change ohne die Klartext-Bezeichner; `model.Change` führt
  sie, lässt sie auf diesem Pfad aber leer (Kommentar in
  `internal/domain/model/change.go`). Die View `cdc.changes` trägt sie
  dagegen (`tools/schema/schema.yaml`, `schema_name`/`table_name` über den
  Join auf `cdc.source_table`).
- **Die opake `SourceTableID` ist kein Draht-Bezeichner.**
  [`ADR-0056`](0056-nats-tabellen-granulares-subjekt.md) Festlegung 1 hat für
  die NATS-Subjekt-Achse genau diesen Punkt entschieden: Über `CDC_TABLES`
  aktivierte Tabellen tragen ein frei gewähltes Token, über die
  SQL-Administration aktivierte deterministisch den qualifizierten Namen —
  ein Draht-Vertrag, der sich auf diese Kennung stützt, gibt einem Consumer
  keine stabile Adresse.

## Entscheidung

Wir wählen **die API bekommt das Changes-Lesen** — als dünne, lesende Fassade
über dem bestehenden `ChangeStorePort`, mit Klartext-Filterachse
(Quelle/Schema/Tabelle), Positionsbereich und Limit, über einen neuen
Inbound Port `ReadChangesUseCase` und den Endpunkt `GET /changes` in der
`reader`-Rechtsklasse. Vier Teilfragen, vier Festlegungen.

### Teilfrage 1 — Port-Form

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (`ADR-0057` Option B bleibt in Kraft, kein Leseport) | kein Aufwand | die **Beschreibung** von `LH-FA-SST-006` bleibt unerfüllt; `ADR-0079`s Beispiel-Client bleibt unbaubar; die Verengung bliebe eine undokumentierte Abweichung vom Vertrag |
| B — Inbound Port als dünne Delegation, Filterachse und Rückgabe **unverändert** (Filter nur über die opake `SourceTableID`, Antwort ohne Schema-/Tabellenname) | kleinster Eingriff; kein Eingriff in den Lesepfad | die Antwort trägt die Tabellen-Identität nur als opakes Token — genau die Form, die `ADR-0056` für einen Draht-Vertrag verworfen hat; die Antwort ist nicht deckungsgleich mit der View, aus der dieselben Daten heute gelesen werden |
| **C — Inbound Port als dünne Delegation, Filterachse in **Klartext** (`source`/`schema`/`table`), Rückgabe trägt Schema- und Tabellenname (gewählt)** | der Consumer adressiert genau das, was das Wecksignal adressiert (`<source>.<schema>.<table>`, `ADR-0056`); die Antwort ist deckungsgleich mit der View-Projektion; kein zweiter Lesepfad — dieselbe Delegation an denselben Port | der Lesepfad zieht einen `cdc.source_table`-Join nach (eine Zeile SQL, ein Mapper-Feld), also ein Eingriff außerhalb der Application-Schicht |
| D — der HTTP-Adapter liest `cdc.changes` selbst (SQL-View-Direktzugriff, wie der SQL-Kanal) | kein neuer Port, keine Anwendungsschicht-Abstraktion | verletzt `ARC-003`/`ARC-005`: ein Driving Adapter spricht ausschließlich über Inbound Ports; der Lesepfad bekäme zwei Implementierungen (Adapter-SQL und `ChangeStorePort`) — der Zweitpfad, den `LH-FA-SST-006` Boundary ausschließt |

**Festlegungen (gewählt: C).**

- **Neuer Inbound Port** `internal/application/port/inbound/readchanges.go`
  (Fähigkeits-Port je [`ADR-0034`](0034-ports-nach-faehigkeiten.md)):
  `ReadChangesUseCase` mit `ReadChanges(ctx, ReadChangesQuery) (ReadChangesResult, error)`.
  Die Transport-Typen liegen **am Port** ([`ADR-0042`](0042-transport-typen-am-port.md));
  der Use Case (`internal/application/usecase/readchanges/`) führt sie als
  Aliase.
- **`ReadChangesQuery`** trägt: `Source` (**Pflicht**), `Schema` und `Table`
  (Klartext, je optional und **unabhängig** — dieselbe Filter-Grammatik wie
  die `WHERE`-Klausel des SQL-Kanals auf der View), `Start`/`End`
  (`*model.SourcePosition`, optional, **Start inklusiv, End exklusiv** —
  `LH-FA-REA-001`), `Limit` (`*int`, optional, gesetzt ≥ 1 — `LH-FA-REA-003`).
- **`ReadChangesResult`** trägt `Changes []ReadChange`; ein `ReadChange`
  bündelt `Position model.SourcePosition`, `Change model.Change` und
  `CommittedAt model.TimePoint` — dieselbe Bündelung, die der Outbound-Port
  als `ChangeRecord` führt (Ordnungsgröße an der Transaktion,
  `LH-FA-REA-004.a`).
- **Der Use Case ist eine Übersetzung, keine Logik:** Er bildet
  `ReadChangesQuery` auf `outbound.ChangeQuery` ab und ruft
  `ChangeStorePort.ReadChanges`. **Kein** neuer Outbound-Port, **keine**
  eigene Abfrage, **keine** zweite Sortierung — die Ordnung trägt der
  bestehende Pfad (`LH-FA-REA-004`).
- **Konsistenz zur View, ausdrücklich:** Derselbe Datenstand, dieselbe
  Bereichs-/Limit-/Filter-Semantik, dieselbe deterministische Sortierung. Der
  SQL-Kanal bleibt, was er nach `ADR-0046` ist — ein direkter
  View-Direktzugriff; der Netzwerk-Kanal ist die Port-Fassade über demselben
  Bestand.

### Teilfrage 2 — Endpunkt-Form

| Option | Pro | Contra |
|---|---|---|
| A — `GET /changes/{source}/{schema}/{table}` (Pfad-Parameter) | „resource-vollständig" | Schema/Tabelle sind **Filter**, nicht Identität der Ressource: die ungefilterte Quelle ist ein legitimer Aufruf, ein Pfad ohne Filter wäre nicht ausdrückbar; die bestehenden neun Endpunkte führen ihre Eingaben als Query-Parameter/JSON-Body (`SPEC-018`) |
| B — `POST /changes` mit Filter im JSON-Body | Filter beliebig erweiterbar | Lesen ist idempotent und nebenwirkungsfrei; `POST` macht daraus eine schreibende Form und weicht von den lesenden Bestandsendpunkten ab (`GET /tables`, `GET /tables/status`) |
| C — eigener Pfad je Filtrationsgrad (z. B. `GET /tables/changes`) | passt in die `/tables`-Familie | verdeckt, dass der Lesezugriff quellen-weit ist und die Tabelle nur ein optionaler Filter; die View heißt `cdc.changes`, nicht `cdc.table_changes` |
| **D — `GET /changes` mit Query-Parametern (gewählt)** | ein Endpunkt für Bereich, Limit und Tabellenfilter; Form der lesenden Bestandsendpunkte (`GET /tables`, `GET /tables/status`); neben `GET /changes/stream` die zweite Form desselben Gegenstands | Filter und Bereich reisen in der URL (kein Body) — Länge und Kodierung liegen beim Client |

**Festlegungen (gewählt: D).**

| Merkmal | Festlegung |
|---|---|
| Endpunkt / Methode | `GET /changes` — neben `GET /changes/stream` ([`SPEC-021`](../../../spec/pflichtenheft.md)) die zweite, **nicht** streamende Form |
| Rechtsklasse | **`reader`** oder `admin` (Teilfrage 3 aus `ADR-0057`; lesend, also die niedrigere Klasse — `admin` deckt sie implizit ab) |
| Query-Parameter | `source` (**Pflicht**), `schema`, `table` (optional, unabhängig), `from`, `to` (optional, `commit_position`-Werte ≥ 1, `from` inklusiv / `to` exklusiv), `limit` (optional, ≥ 1). **Kein** Default-Limit — ohne `limit` liest der Aufruf unbegrenzt, wie der View-Direktzugriff |
| Response | `200` mit `{"changes": [ … ]}`, je Eintrag: `commit_position` (int64), `change_id`, `transaction_id`, `source_table_id`, `schema`, `table` (`string`), `sequence` (int64), `operation` (`INSERT`/`UPDATE`/`DELETE`), `old_image`, `new_image` (eingebettete JSON-Werte, fehlendes Bild `null`), `schema_version`, `committed_at` (RFC 3339 mit Nanosekunden, UTC) |
| Reihenfolge | deterministisch nach (`commit_position`, `transaction_id`, `sequence`) — `LH-FA-REA-004`; die Antwort listet in dieser Ordnung, die Fortsetzung ist `from = <letzte gelieferte commit_position> + 1` |
| Leere Menge | `{"changes": []}`, nie `null` — dieselbe Form, die `listTablesResponse` für leere Listen wählt |
| Fehler-Antwortform | unverändert `{"error": "<Klartext>"}` ([`SPEC-018`](../../../spec/pflichtenheft.md)) |

Die Feldnamen folgen dem **API-Vokabular**, nicht dem Spaltenvokabular der
View: `schema`/`table` statt `schema_name`/`table_name` (so wie
`listTablesResponse` sie schon führt) und `old_image`/`new_image` statt
`old_data`/`new_data` (so wie der gRPC-Stream und der SSE-Stream sie führen,
`SPEC-020`/`SPEC-021`). `source_table_id` bleibt zusätzlich in der Antwort —
die API adressiert Tabellen an anderen Stellen über `table_id`
(`POST /tables/enable`), der Wert ist also kein Fremdkörper.

### Teilfrage 3 — Trägt der Lesepfad die Tabellen-Identität?

| Option | Pro | Contra |
|---|---|---|
| A — Schema/Tabellenname nur als **Filter** zulassen, die Antwort auf die `SourceTableID` beschränken; der Filter wird im Use Case über die Tabellenaktivierung aufgelöst | kein Eingriff in den Lesepfad | die Antwort nennt die Tabelle nur opak; der Use Case bekäme eine zweite Outbound-Abhängigkeit (`TableActivationPort`) für eine Frage, die der Lesepfad selbst beantworten kann; ungefiltert ließe sich der Name gar nicht auflösen |
| **B — der Lesepfad trägt Schema- und Tabellenname (gewählt)** | die Antwort ist deckungsgleich mit der View-Projektion; **eine** Filterachse; der Join auf `cdc.source_table` steht in der View schon genau so | eine Zeile SQL und ein Mapper-Feld werden nachgezogen — die Datei liegt außerhalb der Application-Schicht |
| C — Klartext-Filter im HTTP-Adapter auflösen (Port bleibt auf `SourceTableID`) | der Lesepfad bleibt unberührt | verlagert eine Adapterfrage in den Adapter, der dafür einen Outbound-Port benutzen müsste — der Driving Adapter spräche dann über die Tabellenaktivierung; die Auflösung hinge an einer Registrierung, die mit dem Lesen nichts zu tun hat |

**Festlegungen (gewählt: B).**

- `outbound.ChangeQuery` trägt die Filterachse als **Klartext**:
  `Schema string` und `Table string` (je optional) **ersetzen** das bisherige
  `Table *model.SourceTableID` — eine Filterachse, **eine** Form, kein
  zweiter Weg zum selben Prädikat. Quelle, Start/End und Limit bleiben
  unverändert.
- `SelectChanges` (`internal/adapters/driven/postgresstorage/queries/queries.go`)
  trägt den Join `cdc.source_table` und liest `st.schema_name`,
  `st.table_name` mit — dieselbe Projektion, die die View
  (`tools/schema/schema.yaml`) führt.
- `mapper.ToChange` setzt `Change.Schema` und `Change.Table`. Der
  Lesepfad liefert damit dieselbe Change-Gestalt wie der Capture-Pfad; der
  Kommentar in `internal/domain/model/change.go`, der den Lesepfad als
  schema-/tabellenlos beschreibt, ist nachzuziehen.
- **`ADR-0046` bleibt unberührt.** Die View behält ihren direkten
  Projektionszugriff; es entsteht **keine** Entscheidungslogik in SQL und
  **keine** zweite Leselogik in Go — nur eine zusätzliche Spalte im
  bestehenden Join.

### Teilfrage 4 — Boundary und Negative

| Option | Pro | Contra |
|---|---|---|
| A — ein nicht treffender Filter endet `404` (unbekannte Ressource) | macht eine falsche Adresse sichtbar | die View liefert für dieselbe Eingabe eine leere Menge; der API-Aufruf wäre **nicht** fachlich gleichwertig (`LH-FA-SST-006` Boundary) — und eine Tabelle ohne Changes ist kein Fehler |
| B — ein nicht treffender, aber syntaktisch gültiger Filter endet `200` mit leerer Menge; nur syntaktisch ungültige Eingaben enden `400` (gewählt) | deckungsgleich mit dem View-Verhalten (`LH-FA-REA-006` Boundary) | ein Tippfehler in `table` liefert still eine leere Menge — die Antwort ist nicht falsch, aber schmallippig |
| C — ein Parametername außerhalb der Liste wird ignoriert (Form der neun bestehenden Endpunkte) | einheitlich mit `SPEC-018` | bei einem **Filter**-Endpunkt ist das Ignorieren keine Auslassung, sondern eine **stille Änderung des Ergebnisstands**: ein Tippfehler in `table` liefert ungefiltert alle Changes — genau das „stillschweigende Ignorieren", das `LH-FA-SST-006` Negative ausschließt |

**Festlegungen (gewählt: B, mit der Verschärfung aus C).**

- **Pflichtfeld fehlt** (`source`) → `400`.
- **Unbekannter Query-Parameter** → `400`, mit dem Parameternamen im
  Klartext. Das ist **strenger** als die neun Bestandsendpunkte, und zwar aus
  dem in C genannten Grund: Dort ändert ein zusätzlicher Parameter das
  Ergebnis nicht, hier ändert ein unbekannter **Filter** es still.
- **Syntaktisch ungültige Werte** → `400`: `from`/`to`/`limit` nicht als
  Ganzzahl lesbar, `from`/`to` `< 1` (eine Position 0 existiert nicht,
  `model.SourcePosition`), `limit` `< 1` (`LH-FA-REA-003` Negative), `from >
  to` (`LH-FA-REA-001` Negative).
- **Kein Treffer** (unbekannte Quelle, unbekanntes Schema, unbekannte
  Tabelle, leerer Bereich) → `200` mit `{"changes": []}` (`LH-FA-REA-006`
  Boundary) — **kein** `404`. Der eine `404`-Fall des Bestands
  (`ErrSourceTableMissing` bei `EnableTable`/`DisableTable`/`GetStatus`) ist
  eine **physische** Prüfung an der Quelle und gilt hier nicht: Das Lesen
  prüft nichts an der Quelle, es liest einen Bestand.
- **Nicht autorisiert** → unverändert: kein oder unbekannter Bearer-Token
  `401`; die `reader`-Klasse genügt, ein `403`-Pfad entsteht auf diesem
  Endpunkt nicht.
- **Store-Fehler** (Klasse `storage`) → `500`, geloggt, nicht verworfen.

### Testabdeckung (Erwartung, keine abschließende Festlegung)

- **Whitebox-Unit-Tests des Adapters** (`internal/adapters/driving/http`):
  `401` ohne/unbekanntem Token; `200` mit `reader`-Token; `400` für
  unbekannten Parameter, unlesbare Zahl, `limit < 1`, `from > to`, fehlende
  Quelle; `200` mit leerer Liste ohne Treffer; Feldnamen und Form der Antwort.
- **Unit-Tests des Use Case** (`internal/application/usecase/readchanges`):
  die Übersetzung `ReadChangesQuery` → `outbound.ChangeQuery` und die
  Delegation an den Port (Fake) — der Nachweis, dass **kein** zweiter
  Lesepfad entsteht.
- **Unit-Tests des Lesepfads** (`internal/adapters/driven/postgresstorage`):
  die Abfrage filtert über `st.schema_name`/`st.table_name` und trägt die
  Namen in die Rückgabe; Ordnung unverändert.
- **Erweiterung des E2E-Rundlaufs** (`make test-integration`): ein realer
  `GET /changes` gegen den laufenden Feed-Container über den bestehenden
  HTTP-Wegwerf-Client — das ist der Netzwerk-Beleg zu `LH-FA-SST-006`, den
  die neun Bestandsendpunkte bereits tragen.

**Der geplante Beispiel-Client bleibt lauffähig.** `examples/nats-client`
([`ADR-0079`](0079-nats-beispielclient-vierter-examples-client.md)) baut seine
Abfrage als `GET /changes?source=<id>&schema=<s>&table=<t>`; genau diese Form
fixiert diese ADR. `from` und `limit` sind weglassbar, der Client braucht
**keinen** Nachzug.

## Verglichene Alternativen

Die vier Teilfragen oben tragen je mindestens drei Optionen samt „nichts tun"
(Teilfrage 1, Option A) und Pro/Contra (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

## Konsequenzen

- Positiv: `LH-FA-SST-006`s **Beschreibung** ist erfüllbar — die API trägt
  das Changes-Lesen über denselben Port, über den der CLI-Kanal die Sicht-§4
  bereits zeichnet; kein Zweitpfad, keine zweite Domänenlogik.
- Positiv: `ADR-0079`s Beispiel-Client ist baubar; die Verengung, die er
  vorausgesetzt hat, ist damit an der Wurzel behoben statt im Client umgangen.
- Positiv: Die Antwort ist deckungsgleich mit der View-Projektion
  (einschließlich `schema`/`table`) — der Netzwerk-Consumer bekommt keine
  geringere Sicht als der SQL-Consumer.
- Positiv: Keine neue Fremd-Abhängigkeit, keine neue Schichten-Kante,
  `ARC-005`/`ARC-003` unverändert.
- Negativ: Der Lesepfad wird außerhalb der Application-Schicht berührt (Join,
  Mapper-Feld, Filterfeld) — der Slice geht damit über eine Schicht hinaus.
  Er ist trotzdem **ein** Schnitt nach Lieferwert: Der Leseport ohne
  Gegenstand wäre ein Zombie-Slice (Modul 5), und ein Schicht-Schnitt
  (`…-store`, `…-usecase`, `…-http`) ist die dort ausdrücklich verworfene
  Form.
- Negativ: Ohne `limit` liest der Aufruf unbegrenzt; ein großer Bestand
  antwortet mit einem großen Body. Bewusste Wahl (Deckungsgleichheit mit dem
  View-Direktzugriff), benannt statt versteckt.
- Negativ: Der Endpunkt ist das erste Stück der API, das Filter in der URL
  trägt — die neun Bestandsendpunkte kennen nur Pflichtfelder.
- Folgepflicht (Spec-Nachzug, **nicht** Gegenstand dieser ADR):
  `spec/pflichtenheft.md` verliert in `SPEC-018` die
  Changes-Lesen-Hälfte seines Abgrenzungssatzes und erhält einen **neuen**
  `SPEC-*`-Eintrag mit der Endpunkt-, Parameter- und JSON-Tabelle aus
  Teilfrage 2/4; die Historie bekommt ihre Zeile. `spec/lastenheft.md`
  **braucht keine Änderung** — die Beschreibung von `LH-FA-SST-006` trug die
  Forderung von Anfang an, und ihre Out-of-Scope-Klausel grenzt nur
  Protokoll-, Authn- und Endpunkt-Details als Architektur-/
  Spezifikationsfrage aus.
- Folgepflicht (Benutzerhandbuch): `docs/user/benutzerhandbuch.md` bekommt
  eine Zeile in der Fähigkeits-Tabelle des Abschnitts `### Zugriff über die
  HTTP-/JSON-API` und die Parameter-/Antwort-Beschreibung; die dortige
  Aussage, der Live-Stream sei „die Ausnahme", ist auf zwei lesende
  Endpunkte nachzuziehen.
- Folgepflicht (Kommentar): Der Kommentar in
  `internal/domain/model/change.go`, der die Rekonstruktion eines
  persistierten Change als schema-/tabellenlos beschreibt, beschreibt nach
  diesem Zug nicht mehr, was ist (`AGENTS.md` §3.7).
- Folgepflicht (Fehler-Mapping): `writeDomainError`
  (`internal/adapters/driving/http/errors.go`) bildet heute die
  Domänen-Sentinels auf `400` ab; die Port-Kontrakt-Sentinels
  `outbound.ErrNonPositiveLimit`/`outbound.ErrRangeInverted` kommen hinzu.
- Folgepflicht (Slice-Schnitt-Empfehlung, endgültiger Schnitt liegt bei der
  Planner-Rolle): **ein** Slice, drei Liefer-Punkte — (1) der Lesepfad trägt
  die Tabellen-Identität (Filterachse Klartext, Join, Mapper, Port-Tests),
  (2) `ReadChangesUseCase` samt Inbound Port und Transport-Typen,
  (3) der Endpunkt `GET /changes` samt Fehler-Mapping und Adapter-Tests.
  Der Slice trägt **kein** DB-Schema und **keine** DDL.
- Offen bleibt: **Diagnose/Health** — die andere Hälfte der abgelösten
  `ADR-0057`-Klausel. Sie braucht weiterhin eine eigene
  Port-Design-Entscheidung und ist **nicht** Gegenstand dieser ADR.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Unit-Test (`internal/adapters/driving/http`, Whitebox) | `GET /changes` ohne/unbekanntem Bearer-Token ⇒ `401`; mit `reader`-Token ⇒ `200`; unbekannter Query-Parameter, unlesbare Zahl, `limit < 1`, `from > to`, fehlendes `source` ⇒ `400`; kein Treffer ⇒ `200` mit `{"changes": []}` | `make test` |
| Go-Unit-Test (`internal/application/usecase/readchanges`, Whitebox) | der Use Case ruft `ChangeStorePort.ReadChanges` mit der übersetzten Abfrage auf — Regressionstest gegen einen zweiten Lesepfad | `make test` |
| Go-Unit-Test (`internal/adapters/driven/postgresstorage`) | die Lese-Abfrage filtert über `st.schema_name`/`st.table_name` und trägt beide Namen in die Rückgabe | `make test` |
| `.a-check` | unverändert — die neuen Dateien liegen in den bestehenden Globs `app`/`ports`/`adapters`, keine neue Hexagon-Schichten-Kante | `make a-check` |
| `make test-integration` | realer `GET /changes` gegen den laufenden Feed-Container über den bestehenden HTTP-Wegwerf-Client | `make test-integration` |

## Re-Evaluierungs-Trigger

- **Ein Filter jenseits dieser Achsen wird benannt** (Operation, Zeilen-
  Prädikat, Consumer-Position als Startpunkt, Cursor-/Token-Form der
  Pagination): Die Filter-Grammatik aus Teilfrage 2 braucht dann eine eigene
  Folge-ADR — dieselbe Form, in der `ADR-0056` seine Subjekt-Erweiterung
  behandelt hat.
- **Ein Default-Limit oder eine harte Obergrenze wird gefordert**
  (Ressourcenschutz gegen große Antworten): eigene Folge-ADR — sie ändert den
  Ergebnisstand still, wenn sie nur eingebaut wird.
- **Eine zweite Antwortform wird gefordert** (z. B. ein reiner
  Positionen-/Zähler-Endpunkt ohne Row Images): eigene Folge-ADR.
- **Diagnose/Health** bleibt an die von `ADR-0057` benannte eigene
  Port-Design-Entscheidung gebunden; tritt dieser Bedarf ein, löst das eine
  weitere Folge-ADR aus.
- Sonst permanent — die Port- und Endpunktform gilt unabhängig davon, wie
  viele Changes der Bestand führt.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-15 | Accepted — Architect-Zug: führt `ADR-0057` Teilfrage 2 Option C aus (deren Re-Evaluierungs-Trigger „Bedarf an Changes-Lesen konkret benannt" ist eingetreten); Teil-Supersedes der Changes-Lesen-Hälfte zweier `ADR-0057`-Klauseln | [`ADR-0057`](0057-http-grpc-api.md), [`LH-FA-SST-006`](../../../spec/lastenheft.md), [`ADR-0079`](0079-nats-beispielclient-vierter-examples-client.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0081` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
