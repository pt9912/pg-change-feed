# ADR-0131: gRPC-`ReadChanges` — zehnter RPC im bestehenden `Administration`-Service (ergänzt ADR-0130)

**Status:** Accepted

**Datum:** 2026-09-28

**Autor:** Architect-Agent (Modul 8), Architect-Zug nach explizitem
Auftraggeber-Auftrag (2026-09-28): `ADR-0130`s Out-of-Scope-Klarstellung
wörtlich zitiert („ein zusätzlicher unärer Range-Read-RPC über gRPC wäre
eine eigene, hier nicht beauftragte Fähigkeit") — jetzt ausdrücklich
beauftragt. Jede Aussage über eine Menge nennt die Menge, an der sie
geprüft ist (`AGENTS.md` §3.12); vom Architect selbst gefahren wurde
nichts — diese ADR entscheidet die Umsetzungsform vor der Implementierung,
ihre Fitness-Function-Zeilen sind Erwartungen an den umsetzenden Slice,
keine Messungen.

**Bezug:** [`LH-FA-SST-006`](../../../spec/lastenheft.md) (Haupt-Bezug —
dieselbe Anforderung, jetzt über einen dritten Zugriffsweg neben SQL und
HTTP), [`LH-FA-REA-001`](../../../spec/lastenheft.md)…[`006`](../../../spec/lastenheft.md)
(die gespiegelte Lese-Fähigkeit, unverändert dieselben Boundary-/
Negative-Kriterien wie beim HTTP-Adapter), [`ADR-0130`](0130-grpc-verwaltungs-api-neun-rpcs.md)
(Haupt-Bezug — der zu erweiternde Service, dessen Out-of-Scope-Klausel
und Re-Evaluierungs-Trigger 2 hier eintritt), [`ADR-0081`](0081-changes-lesen-ueber-die-http-api.md)
(der HTTP-Kontrakt, `SPEC-022`, den diese ADR 1:1 auf gRPC überträgt),
[`ADR-0060`](0060-grpc-streaming-mechanismus.md) (Abgrenzung, unberührt —
`ChangeStream` deckt neue Changes bereits live ab, diese ADR deckt den
begrenzten historischen Bereich), [`ADR-0057`](0057-http-grpc-api.md)
(Ursprung des Token-Modells, das hier unverändert wiederverwendet wird),
[`ADR-0034`](0034-ports-nach-faehigkeiten.md) (Ports nach Fähigkeit),
[`ADR-0028`](0028-inbound-use-cases.md) (Inbound Use Cases — derselbe
`ReadChangesUseCase`, den der HTTP-Adapter bereits aufruft, kein zweiter
Domänenpfad), [`ADR-0046`](0046-sql-driving-adapter-lese-schreib-trennung.md)
(SQL-Lese-Trennung — Kontext für denselben Lesepfad über einen dritten
Zugriffsweg)

**Schärft:** — keine neue Architektur-Sicht-Schärfung: `ARC-005` ist
bereits durch `ADR-0130` für den `Administration`-Service konkret
gemacht; diese ADR erweitert dieselbe, bereits konkretisierte Wahl um
einen zehnten RPC. `SPEC-031` (die gRPC-Administration-API-Beschreibung)
erhält die Erweiterung durch die Folgepflicht des umsetzenden Slices,
nicht durch diese ADR selbst.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0130`](0130-grpc-verwaltungs-api-neun-rpcs.md) hat für die neun
Verwaltungsfähigkeiten gRPC eingeführt und dabei ausdrücklich
festgehalten, was **nicht** Gegenstand war: „`LH-FA-REA-*`/`SPEC-022`
(`GET /changes`, Bestandslesen über einen Bereich) ist **nicht**
Gegenstand dieser ADR: Für gRPC existiert für neue Changes bereits ein
Vollinhalts-Livepfad (`StreamChanges`, `ADR-0060`) — ein zusätzlicher
unärer Range-Read-RPC über gRPC wäre eine eigene, hier nicht beauftragte
Fähigkeit." Genau dieser Fall ist eingetreten: Der Auftraggeber (pt9912)
hat am 2026-09-28 diesen Bedarf explizit benannt — der in `ADR-0130`s
eigenem Re-Evaluierungs-Trigger vorgezeichnete Fall 2 („ein konkret
benannter Bedarf an `LH-FA-REA-*`/`SPEC-022` über gRPC"). Das ist keine
Korrektur von `ADR-0130`s damaliger Abgrenzung — die Prämisse „kein
Auftrag benennt diese Fähigkeit" galt zum Entscheidungszeitpunkt real —,
sondern eine neue Tatsache, die eine neue Entscheidung verlangt (Baseline-
Regelwerk `modul-08-agentenrollen.md` §Konflikt-Pfad als Rollen-Sequenz,
Verdikt 2: Folge-ADR statt stiller Korrektur — dieselbe Form, in der
`ADR-0130` selbst `ADR-0057`s Re-Evaluierungs-Trigger behandelt hat).

`ADR-0081` hat dieselbe Fähigkeit für HTTP bereits entschieden und
umgesetzt: ein Inbound Port `ReadChangesUseCase`
(`internal/application/port/inbound/readchanges.go`) mit
`ReadChangesQuery{Source, Schema, Table, Start, End, Limit}` und
`ReadChangesResult{Changes []ReadChange}`, aufgerufen vom HTTP-Adapter
über `GET /changes` (`internal/adapters/driving/http/readchanges.go`,
`SPEC-022`). Dieser Port existiert bereits vollständig und ist gemessen
in Betrieb — die Aufgabe dieser ADR ist ausschließlich, ihn über einen
dritten Zugriffsweg (gRPC) zu erreichen, nicht, ihn neu zu entwerfen.

**Drei Konstraints prägen den Lösungsraum, alle bereits gesetzt:**

- **Derselbe Inbound Use Case, kein zweiter Domänenpfad.** `LH-FA-SST-006`s
  Boundary-Kriterium verlangt fachliche Gleichwertigkeit über alle
  Zugriffswege. Der neue gRPC-Handler ruft exakt
  `inbound.ReadChangesUseCase.ReadChanges` auf — denselben Port, den der
  HTTP-Adapter bereits verwendet.
- **Bestehende Auth-/Fehler-Infrastruktur des `Administration`-Service ist
  bereits generisch genug.** `internal/adapters/driving/grpc/interceptor.go`s
  `authUnaryInterceptor` liest die Rechtsklasse aus einer Tabelle
  (`administrationRPCRoles map[string]role`) — ein zusätzlicher Eintrag
  genügt, kein neuer Interceptor. `administration.go`s
  `administrationError` bildet bereits **alle** Fehler-Sentinels ab, die
  `outbound.ChangeQuery.Validate()` zurückgeben kann
  (`domainerrors.ErrEmptyIdentifier`, `domainerrors.ErrSourceMismatch`,
  `domainerrors.ErrInvalidPosition`, `outbound.ErrNonPositiveLimit`,
  `outbound.ErrRangeInverted` → `codes.InvalidArgument`, geprüft durch
  Lesen von `administration.go` Zeilen 46–63) — kein neuer Zweig nötig.
- **`.a-check.yml`s bestehende Globs sind bereits rekursiv** (unverändert
  gegenüber `ADR-0130`s Feststellung): Ein zehnter RPC in derselben Datei
  und demselben generierten Paket ändert keine Hexagon-Schichten-Kante.

## Entscheidung

Wir erweitern den bestehenden `Administration`-Service
(`proto/cdc/administration/v1/administration.proto`, `ADR-0130`) um einen
zehnten, unären RPC `ReadChanges` — dieselbe `.proto`-Datei, derselbe
`grpc.Server`, dieselbe Adresse `CDC_GRPC_ADDR`, dieselben zwei
Token-Klassen. Sechs Teilfragen.

### Teilfrage 1 — Erweiterung des bestehenden Service oder neuer Service/neues Paket?

| Option | Pro | Contra |
|---|---|---|
| **A — zehnter RPC im bestehenden `Administration`-Service, derselben `.proto`-Datei (gewählt)** | `ReadChanges` teilt mit den neun bestehenden RPCs Form (unär), Trust-Domäne (dieselben zwei Token-Klassen) und Server (`CDC_GRPC_ADDR`) — genau die zwei Achsen, über die `ADR-0130` Teilfrage 1 `ChangeStream`/`Administration` getrennt hat (Stream- vs. Unary-Form, unabhängige Versionierung), unterscheiden sich hier **nicht**; `authUnaryInterceptor`, `administrationRPCRoles` und `administrationError` werden **unverändert** wiederverwendet — ein Eintrag in der Rollen-Tabelle, keine neue Infrastruktur; `tools/harness/grpcadminclient` (der bestehende E2E-Wegwerf-Client) erweitert sich um einen zehnten Aufruf statt einen neuen Client zu bekommen | `Administration` (benannt für Verwaltung: Register/Enable/Disable/Retention) trägt mit `ReadChanges` erstmals eine reine Dateninhalts-Lese-Fähigkeit statt einer Status-/Metadaten-Lese-Fähigkeit (`GetConsumerPosition`/`GetTableStatus`/`ListTables` lesen Zustand, nicht Change-Inhalt) — ein leichter Bedeutungs-Zuwachs des Namens, aber kein Bruch: alle vier sind bereits `roleReader`-RPCs desselben Service |
| B — neuer Service `ChangeRead` in einem neuen Paket `cdc.changeread.v1`, eigene `.proto`-Datei | trennt „Dateninhalt lesen" sauber von „Consumer/Tabelle verwalten" nach `ADR-0034`s Prinzip „Ports nach Fähigkeit"; eigene Versionierung für einen künftig womöglich wachsenden Lese-Vertrag (Filter, Pagination) | doppelte Infrastruktur ohne fachlichen Gewinn: derselbe Server, dieselben zwei Token-Klassen, dieselbe RPC-Form — die Trennungsgründe aus `ADR-0130` Teilfrage 1 (Stream vs. Unary, unabhängige Versionierung) treffen hier nicht zu; ein zweiter `protoc`-Lauf und eine zweite generierte Go-Paket-Wurzel für eine einzelne RPC ist Aufwand ohne belegten Bedarf — das HTTP-Vorbild (`ADR-0081`) hat `GET /changes` ebenfalls **nicht** in einen eigenen Server/Adapter ausgelagert, sondern in denselben `http`-Adapter neben die neun Verwaltungsendpunkte gestellt |
| C — neuer Service `ChangeRead` im selben Paket/derselben Datei `cdc.administration.v1`, aber als eigener `service`-Block neben `Administration` | hält den Namen „Administration" semantisch eng, ohne eine zweite `.proto`-Datei | löst den in B benannten Aufwand nicht, sondern nur einen Teil davon (eine zweite generierte Service-Registrierung, eine zweite Rollen-Tabelle oder ein Merge in dieselbe) — ohne den in B fehlenden fachlichen Grund für eine Trennung zu liefern: ein unnötiger Mittelweg |

**Festlegung.** Wie in Option A: `ReadChanges` wird der `service
Administration`-Block um einen zehnten Eintrag erweitert, in derselben
`.proto`-Datei. Der leichte Bedeutungs-Zuwachs des Namens „Administration"
(Contra in A) wird bewusst in Kauf genommen — dieselbe Abwägung, die das
HTTP-Vorbild bereits getroffen hat: `GET /changes` liegt im selben
`internal/adapters/driving/http/`-Paket, hinter demselben Server, neben
den neun Verwaltungsendpunkten (`ADR-0081`), ohne dass HTTP dafür einen
zweiten Adapter oder Port bekommen hätte.

### Teilfrage 2 — RPC-Form

Unär (Request/Response) — wie im Auftrag benannt und wie alle neun
bestehenden RPCs des `Administration`-Service. Ein begrenzter
Bereichs-Read mit endlicher Antwort ist kein Server-Streaming-Fall
(`ADR-0060` Teilfrage 1 hätte hier keinen strukturellen Vorteil zu bieten:
kein „ein Aufruf, viele Antworten über die Zeit" — die Antwort ist eine
einzelne, wenn auch potenziell große, Liste, exakt wie bei
`GET /changes`).

```proto
rpc ReadChanges(ReadChangesRequest) returns (ReadChangesResponse);
```

### Teilfrage 3 — Nachrichtenschema

Feldnamen/Typen 1:1 aus dem bestehenden HTTP-Kontrakt (`SPEC-022`,
`internal/adapters/driving/http/readchanges.go`) übernommen — dieselbe
Übernahmeregel wie `ADR-0130` Teilfrage 3; `snake_case` in proto3 wird zu
`camelCase` im generierten Go-Code, Formsache der Sprache. `from`/`to`
tragen den Positions-Offset als `uint64` mit `0` als „nicht gesetzt" —
kein `optional`-Feld nötig, weil eine gültige Quellposition laut
`model.SourcePosition`/`domainerrors.ErrInvalidPosition` ohnehin `≥ 1`
sein muss (`0` ist ein Wert, den keine gültige Position je annimmt, also
eindeutig als „unbegrenzt" lesbar — dieselbe Beobachtung gilt für `limit`
als `int64` mit `0` = „kein Limit", `LH-FA-REA-003`).

| Nachricht | Felder |
|---|---|
| `ReadChangesRequest` | `string source` (Pflicht), `string schema`, `string table` (je optional, unabhängig), `uint64 from`, `uint64 to` (je optional, `0` = nicht gesetzt, Start inklusiv/Ende exklusiv), `int64 limit` (optional, `0` = unbegrenzt) |
| `ChangeRecord` (Hilfsnachricht, ein Element von `ReadChangesResponse.changes`) | `int64 commit_position`, `string change_id`, `string transaction_id`, `string source_table_id`, `string schema`, `string table`, `int64 sequence`, `string operation`, `bytes old_image`, `bytes new_image`, `string schema_version`, `string committed_at`, `string origin` — dieselben dreizehn Felder wie `readChangeResponse` im HTTP-Adapter, in derselben Reihenfolge |
| `ReadChangesResponse` | `repeated ChangeRecord changes` — eine leere Menge trägt eine leere, gesetzte Liste, nie ein `nil`-Slice (`LH-FA-REA-006` Boundary) |

`old_image`/`new_image` als `bytes` — dieselbe Form wie
`streamv1.Change` (`changestream.proto`, Zeilen 20–21): die Row Images
reisen als rohe JSON-Bytes unverändert, keine zweite Repräsentation.
`committed_at` bleibt ein `string` im RFC-3339-Nanosekunden-Format
(dieselbe Form wie die HTTP-Antwort) statt `google.protobuf.Timestamp` —
Konsistenz mit dem 1:1-Übernahmeprinzip aus `ADR-0130` Teilfrage 3 wiegt
hier höher als protobuf-idiomatische Typwahl; eine Umstellung auf
`Timestamp` wäre eine eigene, hier nicht getroffene Entscheidung.

`ChangeRecord` heißt bewusst wie `outbound.ChangeRecord`
(`internal/application/port/outbound/changestore.go`) — dieselbe
Bündelung (Change + Commit-Position + Commit-Zeitpunkt seiner
Quelltransaktion), zwei verschiedene Schichten (Go-Domänentyp vs.
Draht-Nachricht), keine Namenskollision, weil unterschiedliche Pakete.

Kein Feld trägt einen Decode-Fehlerfall („unbekannter Query-Parameter",
„nicht lesbare Zahl" aus der HTTP-Fehlertabelle): Das Nachrichtenschema
ist bereits typisiert (Protobuf), diese beiden Fehlerpfade entfallen
strukturell — dieselbe Feststellung wie `ADR-0130` Teilfrage 3.

### Teilfrage 4 — Rechtsklasse

**`roleReader`** — lesend, wie `GetConsumerPosition`/`GetTableStatus`/
`ListTables` (`ADR-0130` Teilfrage 4). Ein `admin`-Token deckt sie
implizit ab (dieselbe Hierarchie `roleAdmin ≥ roleReader`). Der Eintrag
tritt der bestehenden Tabelle `administrationRPCRoles`
(`internal/adapters/driving/grpc/interceptor.go`) als zehnte Zeile bei:

```go
"ReadChanges": roleReader,
```

Kein neuer Interceptor, keine neue Konstruktions-Option — derselbe
`authUnaryInterceptor`, dieselbe Registrierung an `grpc.NewServer(...)`.

### Teilfrage 5 — Fehlerform

`administrationError` (`internal/adapters/driving/grpc/administration.go`)
bleibt **unverändert**: Jeder Fehler-Sentinel, den
`outbound.ChangeQuery.Validate()` zurückgeben kann, ist bereits im
bestehenden Switch abgedeckt (geprüft durch Lesen der Datei bei Abfassung
dieser ADR).

| HTTP-Status (Ursache, `ADR-0081` Teilfrage 4) | gRPC-Code |
|---|---|
| `400` — unbekannter Query-Parameter / nicht lesbare Zahl | entfällt strukturell (Teilfrage 3) |
| `400` — fehlende Quelle (`domainerrors.ErrEmptyIdentifier`) | `codes.InvalidArgument` (bereits im bestehenden Switch) |
| `400` — Positions-Offset `< 1` (`domainerrors.ErrInvalidPosition`) | `codes.InvalidArgument` (bereits im bestehenden Switch) |
| `400` — `limit < 1` (`outbound.ErrNonPositiveLimit`) | `codes.InvalidArgument` (bereits im bestehenden Switch) |
| `400` — `from > to` (`outbound.ErrRangeInverted`) | `codes.InvalidArgument` (bereits im bestehenden Switch) |
| `401` — fehlender/unbekannter Bearer-Token | `codes.Unauthenticated` (unverändert `authUnaryInterceptor`) |
| `403` — Rechtsklasse unzureichend | entfällt praktisch: `ReadChanges` verlangt nur `roleReader`, jedes gültige Token (`reader` oder `admin`) genügt — dieselbe Beobachtung wie bei `ADR-0060` Teilfrage 4 für `StreamChanges` |
| kein Treffer (`LH-FA-REA-006` Boundary) | **kein** `codes.NotFound` — eine leere `ReadChangesResponse.changes`-Liste bei Erfolg, dieselbe Festlegung wie `ADR-0081` Teilfrage 4 Option B |
| `500` — unerwarteter interner Fehler | `codes.Internal` (unverändert `administrationError`) |

Der Handler ruft `s.readChanges.ReadChanges(ctx, inbound.ReadChangesQuery{…})`
und übersetzt Erfolg/Fehler über exakt dieselbe `administrationError`-Funktion,
die die neun bestehenden Handler bereits verwenden — kein neuer Code-Pfad
für Fehlerklassifikation.

### Teilfrage 6 — Server-/Adapter-Platzierung

Derselbe `grpc.Server`, dieselbe Adresse `CDC_GRPC_ADDR`
(`ADR-0130` Teilfrage 6, unverändert). Additiv:

- `administrationService` (`administration.go`) erhält ein zehntes Feld
  `readChanges inbound.ReadChangesUseCase`.
- `Config` (`server.go`) erhält ein zehntes Feld
  `ReadChanges inbound.ReadChangesUseCase`, verdrahtet in `New()` in die
  bestehende `administrationv1.RegisterAdministrationServer(...)`-Aufrufstelle.
- Der Handler importiert wie die neun bestehenden Handler ausschließlich
  Inbound Ports (`internal/application/port/inbound/`) und Domain-Typen
  zur Übersetzung — keine Driven-Adapter-Interna, dieselbe Grenze wie beim
  bestehenden Service (`ARC-002`/`ARC-003`/`ARC-005`).

## Verglichene Alternativen

Teilfrage 1 trägt drei Optionen samt „nichts tun" wäre hier keine
sinnvolle vierte Option (der Auftrag benennt die Fähigkeit explizit als
beauftragt — Ablehnung stünde im Widerspruch zum Auftrag selbst und ist
deshalb nicht als ernsthafte Alternative geführt, anders als in
`ADR-0130`/`ADR-0081`, wo der Trigger erst neu bewertet werden musste).

## Konsequenzen

- Positiv: Ein Consumer, der ausschließlich gRPC spricht, erreicht jetzt
  denselben begrenzten Bereichs-Lesezugriff wie über HTTP/SQL — schließt
  die von `ADR-0130`s Out-of-Scope-Klausel offen gelassene Lücke.
- Positiv: Keine neue Infrastruktur — `authUnaryInterceptor`,
  `administrationRPCRoles`, `administrationError` und der bestehende
  E2E-Wegwerf-Client (`tools/harness/grpcadminclient`) werden alle
  unverändert erweitert statt neu geschaffen.
- Positiv: `ChangeStream`/`changestream.proto` und der bereits bestehende
  HTTP-Lesepfad (`ReadChangesUseCase`, `ADR-0081`) bleiben durch diese
  Entscheidung unverändert — kein Zeilen-Diff an bestehendem, produktivem
  Code dieser Dateien ist Folge dieser ADR (nur `administration.proto`,
  `administration.go`, `interceptor.go`, `server.go` wachsen additiv).
- Negativ: `Administration` trägt jetzt zehn RPCs — der Name bleibt
  streng genommen nur für neun davon wörtlich zutreffend (siehe Teilfrage
  1, Contra A); akzeptiert, weil eine Aufspaltung ohne fachlichen Gewinn
  wäre (Option B/C).
- Folgepflicht: Der umsetzende Slice-Schnitt ergänzt
  `proto/cdc/administration/v1/administration.proto` um
  `ReadChangesRequest`/`ChangeRecord`/`ReadChangesResponse` und den
  zehnten RPC, regeneriert den Go-Code (derselbe `protoc`-Lauf,
  `tools/harness/proto-generate.sh`, keine neue Build-Stufe), implementiert
  den Handler in `administration.go`, den zehnten Tabelleneintrag in
  `interceptor.go` und die Config-/Wiring-Erweiterung in `server.go`.
- Folgepflicht: `spec/pflichtenheft.md` erweitert `SPEC-031` um den
  zehnten RPC, sein Nachrichtenschema und die (unveränderte)
  Fehlercode-Tabelle — kein neuer `SPEC-*`-Eintrag nötig, weil `SPEC-031`
  bereits den `Administration`-Service als Ganzes beschreibt; Gegenstand
  des umsetzenden Slices, nicht dieser ADR.
- Folgepflicht: `tools/harness/grpcadminclient` bekommt einen `ReadChanges`-
  Aufruf (analog zu seinem bestehenden `ListTables`-Aufruf für die
  `reader`-Klasse); `tools/harness/run-integration-tests.sh` erweitert den
  bestehenden gRPC-Administration-Rundlauf um den neuen Beleg.
- Folgepflicht: Beispiel-Clients (Go/C#/Kotlin, volle Matrix laut
  `SPEC-018`§„Sprachen und Umfang") und die Erweiterung der drei
  SDK-Packages (`PgChangeFeed.Client`, `pgchangefeed`,
  `pgchangefeed-kotlin`) um `ReadChanges` — je eigener Folge-Schritt/
  Slice, nicht Gegenstand dieser ADR.
- Folgepflicht: E2E-Beleg (`make test-integration`) — ein realer
  `ReadChanges`-Aufruf mit einem `reader`-Token liefert eine zuvor
  committete Änderung, über ihre `change_id` gegen `cdc.changes`
  gehalten — analog zum bestehenden `GET /changes`-Beleg (`ADR-0081`).
- Folgepflicht: `spec/architecture.md` braucht **keine** Änderung —
  `ARC-005` ist bereits durch `ADR-0130` konkret gemacht; diese ADR
  erweitert dieselbe, bereits konkretisierte Wahl.

## Fitness Function (falls maschinell prüfbar; Erwartung an den umsetzenden Slice, nicht durch diesen Architect-Lauf gefahren)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Unit-Test (`internal/adapters/driving/grpc`, Whitebox, `authUnaryInterceptor`) | Ein gültiges `reader`-Token erreicht `ReadChanges`; ein gültiges `admin`-Token ebenfalls (implizite Abdeckung, `roleAdmin ≥ roleReader`); ein Aufruf ohne oder mit unbekanntem Token endet mit `codes.Unauthenticated` | `make test` |
| Go-Unit-Test (`internal/adapters/driving/grpc`, Whitebox, `ReadChanges`-Handler) | Der Handler ruft exakt `inbound.ReadChangesUseCase.ReadChanges` mit der aus dem Request übersetzten Abfrage auf (Fake-Use-Case im Test) — kein zweiter Domänenpfad; die vier Fehler-Sentinels aus Teilfrage 5 (`ErrEmptyIdentifier`, `ErrInvalidPosition`, `ErrNonPositiveLimit`, `ErrRangeInverted`) laufen unverändert über `administrationError` auf `codes.InvalidArgument`; eine leere Trefferliste liefert `codes.OK` mit leerem, aber gesetztem `changes`-Feld | `make test` |
| `.a-check` | unverändert — der zehnte RPC liegt in denselben, bereits rekursiven Globs `adapters: ["internal/adapters/**"]`/`contract: ["gen/**"]`, kein neuer Hexagon-Schichten-Edge nötig (geprüft durch Lesen von `.a-check.yml` bei Abfassung dieser ADR) | `make a-check` |
| `make test-integration` | realer `ReadChanges`-Aufruf über `tools/harness/grpcadminclient` gegen den laufenden Feed-Container: eine zuvor committete Änderung erscheint in der Antwort, über ihre `change_id` gegen `cdc.changes` gehalten; ein Aufruf ohne Token endet mit `Unauthenticated` — Umsetzung Gegenstand des umsetzenden Slice-Schnitts | `make test-integration` |

## Re-Evaluierungs-Trigger

Permanent — die Wahl „zehnter RPC im bestehenden `Administration`-Service,
unär, `roleReader`, dieselbe `.proto`-Datei" gilt unabhängig vom
Umsetzungszeitpunkt. Ausnahmen, die eine eigene Folge-ADR brauchen: (1)
ein konkret benannter Bedarf an einer Stream-Form für Bereichs-Reads
(z. B. sehr große Antworten, die eine paginierte oder gestreamte Zustellung
verlangen) — Teilfrage 2 neu bewerten; (2) ein Bedarf an einer eigenen
Rechtsklasse für `ReadChanges` jenseits `reader`/`admin`; (3) ein
praktisches Wachstums-Limit des `Administration`-Service (weitere,
fachlich fremde RPCs) macht eine nachträgliche Aufspaltung entlang der in
Teilfrage 1 Option B skizzierten Grenze nötig.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-28 | Accepted — Architect-Entscheidung nach explizitem Auftraggeber-Auftrag; erweitert `ADR-0130`s `Administration`-Service um den dort als Out-of-Scope benannten, jetzt beauftragten Range-Read-RPC | [`ADR-0130`](0130-grpc-verwaltungs-api-neun-rpcs.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0131` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
