# ADR-0132: Diagnose über einen neuen Inbound Port — HTTP und gRPC (Supersedes ADR-0057, teilweise)

**Status:** Accepted — Supersedes [`ADR-0057`](0057-http-grpc-api.md) in genau
**zwei Klauseln**, jeweils in ihrer Diagnose/Health-Hälfte (die
Changes-Lesen-Hälfte derselben zwei Klauseln hat bereits
[`ADR-0081`](0081-changes-lesen-ueber-die-http-api.md) abgelöst): dem
Ausschluss-Satz in §Entscheidung Teilfrage 2 („Changes-Lesen und
Diagnose/Health bleiben damit bewusst ausgeschlossen — nicht verworfen,
sondern vertagt") und der zugehörigen Negativ-Konsequenz in §Konsequenzen
(„Changes-Lesen (`LH-FA-REA-*`) und Diagnose/Health bleiben nach dieser ADR
weiterhin CLI-/SQL-exklusiv"). Damit ist keine der beiden Klauseln aus
`ADR-0057` mehr in Kraft — die zweite (Changes-Lesen-)Hälfte löste
`ADR-0081`, die erste (Diagnose/Health-)Hälfte löst diese ADR. Alles Übrige
aus `ADR-0057` bleibt unverändert in Kraft: die Protokoll-Wahl HTTP/JSON
(Teilfrage 1), die zwei Token-Rechtsklassen (Teilfrage 3), die
Adapter-Platzierung (Teilfrage 4).

**Datum:** 2026-09-28

**Autor:** Architect-Agent (Modul 8), Architect-Zug nach zweifach explizit
bestätigtem Auftraggeber-Auftrag (2026-09-28): Diagnose/Health soll über
eine neue Port-Abstraktion sowohl per HTTP als auch per gRPC erreichbar
werden — Antwort auf die vom Vorgänger-Agenten gestellte Frage „Soll ich
das … als weiteren Auftrag aufnehmen?" war „ja". Jede Aussage über eine
Menge nennt die Menge, an der sie geprüft ist (`AGENTS.md` §3.12); vom
Architect selbst gefahren wurde nichts — diese ADR entscheidet die
Umsetzungsform vor der Implementierung, ihre Fitness-Function-Zeilen sind
Erwartungen an den umsetzenden Slice, keine Messungen.

**Bezug:** [`LH-FA-SST-003`](../../../spec/lastenheft.md) (Haupt-Bezug —
der bestehende Diagnose-Sondermodus, dessen Fähigkeit diese ADR über zwei
weitere Zugriffswege öffnet), [`LH-FA-ADM-002`](../../../spec/lastenheft.md)…[`005`](../../../spec/lastenheft.md),
[`LH-FA-RET-005`](../../../spec/lastenheft.md), [`LH-FA-RET-006`](../../../spec/lastenheft.md),
[`LH-FA-CAP-009`](../../../spec/lastenheft.md) (die sechs Diagnose-Signale,
unverändert dieselben Felder, jetzt zusätzlich über HTTP/gRPC),
[`LH-FA-SST-006`](../../../spec/lastenheft.md) (Boundary-Kriterium
„fachliche Gleichwertigkeit über alle Zugriffswege"),
[`LH-QA-SEC-001`](../../../spec/lastenheft.md)…[`003`](../../../spec/lastenheft.md),
[`ADR-0057`](0057-http-grpc-api.md) (teilweise superseded — Haupt-Bezug,
siehe oben), [`ADR-0081`](0081-changes-lesen-ueber-die-http-api.md) (Vorbild:
dieselbe Klausel-Ablösungsform, für die andere Hälfte derselben zwei
Klauseln), [`ADR-0130`](0130-grpc-verwaltungs-api-neun-rpcs.md) (der zu
erweiternde `Administration`-Service), [`ADR-0131`](0131-grpc-readchanges-zehnter-rpc.md)
(Vorbild: derselbe Erweiterungsschnitt — ein weiterer unärer RPC im
bestehenden Service statt eines neuen), [`ADR-0034`](0034-ports-nach-faehigkeiten.md)
(Ports nach Fähigkeit — Vorbild für den neuen Outbound Port),
[`ADR-0028`](0028-inbound-use-cases.md) (Inbound Use Cases),
[`ADR-0047`](0047-rollenspezifische-dsn-verdrahtung.md) (Rollen-Modell,
`cdc_reader`, das der neue Adapter unverändert verwendet), [`ADR-0023`](0023-fehlerklassifikation.md)
(Fehlerklassifikation, deren Übersetzungsmuster diese ADR fortsetzt)

**Schärft:** — keine neue Architektur-Sicht-Schärfung: `ARC-005` ist für
HTTP bereits durch `ADR-0057`/`ADR-0081` und für gRPC bereits durch
`ADR-0130` konkret gemacht; diese ADR öffnet über beide bereits
konkretisierten Driving-Adapter eine weitere, bereits im CLI-Sondermodus
bestehende Fähigkeit. `SPEC-018` (HTTP-API) und `SPEC-031` (gRPC-
Administration-API) erhalten die Erweiterung durch die Folgepflicht des
umsetzenden Slices, nicht durch diese ADR selbst.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

Der bestehende `diagnose`-Sondermodus
(`docs/user/benutzerhandbuch.md` §„Diagnose ausführen", `LH-FA-SST-003`,
deckt `LH-FA-ADM-002`…`005`, `LH-FA-RET-005`/`006`, `LH-FA-CAP-009`) läuft
ausschließlich über `docker exec`/einen direkten, kurzlebigen
Prozessaufruf — kein Netzwerk-Zugriffsweg. Intern ist `Diagnose`
(`internal/bootstrap/wiring.go` Zeile 1927) eine **freie Funktion ohne
Anwendungsschicht-Abstraktion**: Sie öffnet selbst einen `pgxpool.Pool`
gegen `cfg.ReaderDSN` und stellt fünf einzelne SQL-Abfragen direkt gegen
`cdc.heartbeat`, `cdc.metrics` (dreimal) und `cdc.retention_blockers` bzw.
`cdc.backfill_status`, deren Ergebnisse sie unmittelbar als Text auf
`stdout` schreibt. Sie liegt damit außerhalb des Hexagons: Weder ein
Inbound Port (`ARC-003`) noch ein Outbound Port (`ARC-004`) vermitteln
diesen Zugriff — dieselbe Lücke, die `ADR-0057` bereits 2026-09-14 explizit
benannt hat: „`Diagnose` und `Healthcheck` … sind freie Funktionen ohne
Anwendungsschicht-Abstraktion — für einen Driving Adapter, der laut
`ARC-002`/`ARC-003` nur über Inbound Ports mit der Application-Schicht
sprechen darf, fehlt hier die Vermittlung." `ADR-0057` schloss die
Fähigkeit deshalb bewusst aus der ersten HTTP-API-Version aus (Teilfrage 2,
Option D) und trug einen eigenen Re-Evaluierungs-Trigger dafür; `ADR-0130`s
Konsequenzen-Abschnitt zitiert dieselbe Lücke wörtlich als weiterhin offen
(„[Changes-Lesen/Diagnose] muss auf einen Folge-Slice warten"). Die
Changes-Lesen-Hälfte dieser Lücke hat `ADR-0081` bereits für HTTP
geschlossen (mit `ReadChangesUseCase`/`ChangeStorePort` als sauberem
Vorbild); die Diagnose-Hälfte blieb offen. Der Auftraggeber (pt9912) hat
diesen Trigger am 2026-09-28 zweifach explizit ausgelöst: Diagnose/Health
soll über eine neue Port-Abstraktion sowohl per HTTP als auch per gRPC
erreichbar werden.

**Abgrenzung — was „Health" hier nicht bedeutet.** Neben `Diagnose` trägt
`internal/bootstrap/wiring.go` eine zweite freie Funktion, `Healthcheck`
(Zeile 1870): Sie liefert **keinen Bericht**, sondern einen binären
Exit-Code für Dockers `HEALTHCHECK`-Direktive — ein Mechanismus, der einen
Prozess-Exit-Code verlangt, keine Netzwerk-Antwort. Ein HTTP- oder
gRPC-Endpunkt ersetzt diesen Mechanismus strukturell nicht: Docker ruft
`HEALTHCHECK` als Kommando im Container auf, nicht als externe Netzwerk-
Prüfung. Die von `Healthcheck` getroffene binäre Verdikt-Entscheidung
(gesund/ungesund anhand `cdc.heartbeat.age_seconds`) ist zudem eine echte
Teilmenge der Felder, die `Diagnose` bereits unverändert roh ausgibt
(Betriebsstatus). Diese ADR entscheidet deshalb: **„Health" in „Diagnose/
Health" ist hier die im Diagnose-Bericht bereits enthaltene
Betriebsstatus-Information, keine zweite, eigene Fähigkeit.**
`Healthcheck` selbst bleibt unverändert CLI-/Docker-exklusiv (siehe
Teilfrage 3) — dieselbe Grenze, die auch `ADR-0057`s eigener Wortlaut zog,
ohne sie damals aufzulösen.

Drei Konstraints prägen den Lösungsraum:

- **Kein bestehender Port zu vermitteln — anders als bei `ADR-0130`/`ADR-0131`.**
  Beide vorangegangenen gRPC-ADRs verdrahteten **bereits existierende**,
  saubere Inbound Ports (`consumer.go`, `verwaltung.go`, `retention.go`,
  `readchanges.go`) über einen neuen Transportweg. Für Diagnose existiert
  dieser Port nicht — die Aufgabe dieser ADR ist größer: Sie muss zuerst den
  fehlenden Inbound Port **und** einen fehlenden Outbound Port entwerfen,
  bevor HTTP/gRPC ihn erreichen können.
- **`LH-FA-SST-006`s Boundary-Kriterium** verlangt fachliche
  Gleichwertigkeit über alle Zugriffswege — der neue HTTP-/gRPC-Weg liefert
  exakt dieselben sechs Signalgruppen wie der bestehende CLI-Sondermodus,
  kein zweiter Domänenpfad.
- **`.a-check.yml`s bestehende Globs sind bereits rekursiv**
  (`ports: ["internal/application/port/**"]`,
  `adapters: ["internal/adapters/**"]`, `contract: ["gen/**"]`) — ein neuer
  Inbound/Outbound Port und ein neuer Driven-Adapter-Baustein liegen bereits
  in den bestehenden Globs, keine neue Hexagon-Schichten-Kante.

## Entscheidung

Wir schaffen einen neuen Inbound Port `DiagnoseUseCase` mit einem
vollständig strukturierten Rückgabewert, einen neuen Outbound Port
`DiagnosticsPort`, der die fünf bestehenden SQL-Lesezugriffe aus einem
neuen Driven Adapter (`postgresstorage`) bündelt, und exponieren diesen
Use Case über **drei** Zugriffswege: den bestehenden CLI-Sondermodus (jetzt
als dünner Wrapper), einen neuen `GET /diagnose`-Endpunkt (`reader`-Klasse)
und einen elften RPC `Diagnose` im bestehenden `Administration`-gRPC-
Service. Sechs Teilfragen.

### Teilfrage 1 — Schnitt und Rückgabeform des neuen Inbound Ports

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (freie Funktion bleibt) | kein Aufwand | widerspricht dem zweifach expliziten Auftrag; verstößt strukturell weiter gegen `ARC-002`/`ARC-003` |
| B — mehrere kleine Ports, einer je Signalgruppe (`HeartbeatQuery`, `MetricsQuery`, `RetentionBlockerQuery`, `BackfillStatusQuery`) | feingranular, jede Fähigkeit einzeln testbar | fragmentiert das, was CLI/HTTP/gRPC als **einen** Bericht liefern, in vier bis fünf Aufrufstellen je Zugriffsweg — höherer Koordinationsaufwand ohne belegten Bedarf; kein bestehendes Vorbild verlangt getrennte Aufrufe (`GetStatus`/`ListTables` sind je eine eigene, aber inhaltlich disjunkte Fähigkeit, keine Zerlegung eines einzelnen Berichts) |
| **C — ein Inbound Port `DiagnoseUseCase.Diagnose(ctx, DiagnoseQuery{Source}) (DiagnoseResult, error)` mit vollständig strukturierter Rückgabe (gewählt)** | ein Aufruf, ein Bericht — dieselbe Semantik wie der bestehende CLI-Sondermodus und wie `ReadChangesUseCase` (`ADR-0081`); `DiagnoseResult` trägt exakt die sechs Signalgruppen, die der heutige Text bereits ausgibt — keine neue Fähigkeit erfunden, nur strukturiert statt direkt gedruckt | ein einzelner Rückgabetyp mit mehreren optionalen Unterfeldern — abgefedert durch benannte Verschachtelung (siehe unten) |

`DiagnoseResult` (`internal/application/port/inbound/diagnose.go`), Felder
1:1 aus dem heutigen Text-Bericht (`internal/bootstrap/wiring.go` Zeilen
1941–2042) übernommen:

```go
type DiagnoseQuery struct {
    Source model.SourceID
}

// HeartbeatAgeSeconds == nil bedeutet „kein Lebenszeichen — Instanz hat
// noch nie geschlagen" (LH-FA-ADM-002 Boundary); in diesem Fall trägt
// ErrorClass ebenfalls nil, gelesen als „unbekannt", nicht als
// „Normalbetrieb" — die Unterscheidung liegt an HeartbeatAgeSeconds, nicht
// an einem dritten Feld.
type DiagnoseResult struct {
    HeartbeatAgeSeconds *float64
    ErrorClass          *string
    CaptureLag          float64
    ConsumerLags        []ConsumerLag
    RetentionBlocker    *RetentionBlocker
    StorageBytes        float64
    Backfill             []BackfillTableStatus
}

// Lag == nil bedeutet „Quelle trug noch nie eine Transaktion"
// (LH-FA-ADM-005 Fußnote im bestehenden Text).
type ConsumerLag struct {
    ConsumerID string
    Lag        *float64
}

// RetentionBlocker == nil bedeutet „kein Blocker" (LH-FA-RET-005
// Boundary); Backlog == nil bedeutet „Rückstand unbekannt, Quelle trug
// noch nie eine Transaktion" — dieselbe Fußnote wie bei ConsumerLag.
type RetentionBlocker struct {
    ConsumerID           string
    Name                 string
    AcknowledgedPosition int64
    Backlog              *int64
}

// EstimatedRows == nil bedeutet „unbekannt", nie 0 (LH-FA-CAP-009,
// diagnoseBackfillStatus-Kommentar).
type BackfillTableStatus struct {
    Schema            string
    Table             string
    Status            string
    RowsCopied        int64
    EstimatedRows     *int64
    WarnEstimatedSize bool
    WarnDuration      bool
    ErrorMessage      string
}

type DiagnoseUseCase interface {
    Diagnose(ctx context.Context, query DiagnoseQuery) (DiagnoseResult, error)
}
```

### Teilfrage 2 — Neuer Outbound Port, kein direkter DB-Zugriff aus der Application-Schicht

| Option | Pro | Contra |
|---|---|---|
| A — der neue Application Service hält selbst einen `*pgxpool.Pool`/DSN und fragt die fünf Views direkt ab (wie die heutige freie Funktion) | minimal invasiv, kein neuer Adapter-Baustein | verstößt gegen `ARC-002` („Application … darf NICHT importieren: konkrete Adapter, Treiber-Frameworks") — exakt der Verstoß, den `ADR-0057` bereits als Grund dafür nannte, dass `Diagnose` bisher keine Portabstraktion trug; löst die benannte Lücke nicht, verschiebt sie nur um eine Schicht |
| **B — neuer Outbound Port `DiagnosticsPort` (`internal/application/port/outbound/diagnostics.go`), implementiert von einem neuen Driven Adapter `postgresstorage.NewDiagnostics(ctx, dsn, opts...)` (gewählt)** | volle Schichten-Konformität (`ARC-002`/`ARC-004`/`ARC-006`); derselbe Konstruktions-Stil wie jeder bestehende `postgresstorage`-Adapter (`NewHeartbeat`, `NewBackfillRun`, …, alle `func New*(ctx, dsn, opts...) (*Adapter, error)`); die fünf SQL-Abfragen wandern mechanisch vom `bootstrap`-Paket in den neuen Adapter — kein neuer SQL-Text, nur eine neue Adresse; schließt die von `ADR-0057` benannte architektonische Lücke tatsächlich, statt sie zu verschieben | ein weiterer Adapter-Baustein/eine weitere Datei |
| C — die fünf Lesezugriffe in den bestehenden `HeartbeatPort` (Schreib-Port des Capture-Pfads) aufnehmen | kein neuer Port-Typ | `HeartbeatPort` ist nach `ADR-0034`s „Ports nach Fähigkeit" bewusst eng auf die Schreib-Fähigkeit des Capture-Pfads geschnitten (`Beat`/`Fault`); eine Lese-Fähigkeit für einen anderen Konsumenten (Diagnose-Aufrufer, nicht der Capture-Prozess) in denselben Port zu packen vermischt zwei Fähigkeiten hinter einem Namen — derselbe Einwand, den `ADR-0130` Teilfrage 1 gegen das Vermischen von Streaming und Verwaltung vorbrachte |

`DiagnosticsPort` (Outbound, eigener Typ analog zu `outbound.ChangeRecord`
vs. `inbound.ReadChange` — zwei Schichten, zwei Typen, 1:1 gemappt vom
Application Service):

```go
var ErrDiagnosticsStorage = stderrors.New("Fehlerklasse storage: Persistenzfehler beim Lesen der Diagnose-Views")

type DiagnosticsSnapshot struct {
    HeartbeatAgeSeconds *float64
    ErrorClass          *string
    CaptureLag          float64
    ConsumerLags        []ConsumerLagSnapshot
    RetentionBlocker    *RetentionBlockerSnapshot
    StorageBytes        float64
    Backfill            []BackfillTableSnapshot
}
// ConsumerLagSnapshot, RetentionBlockerSnapshot, BackfillTableSnapshot:
// dieselben Felder wie ihre inbound-Entsprechungen (Teilfrage 1).

type DiagnosticsPort interface {
    Read(ctx context.Context, source model.SourceID) (DiagnosticsSnapshot, error)
}
```

Der neue Application Service (`internal/application/usecase/diagnose/service.go`,
Paketname `diagnose`, Muster identisch zu `readchanges.ReadChangesService`)
implementiert `inbound.DiagnoseUseCase` als dünne Fassade über
`DiagnosticsPort` — keine eigene Abfrage, keine zweite Sortierung. Der
Driven Adapter (`postgresstorage.NewDiagnostics`) trägt intern dieselben
fünf SQL-Texte, die heute in `internal/bootstrap/wiring.go` stehen
(mechanische Verschiebung, kein neuer SQL-Text) und wird — wie jeder andere
`cdc_reader`-gebundene Adapter — mit `cfg.ReaderDSN` konstruiert (`ADR-0047`).

### Teilfrage 3 — Verhältnis zum bestehenden CLI-`diagnose`-Sondermodus und zu `Healthcheck`

| Option | Pro | Contra |
|---|---|---|
| A — CLI-Sondermodus bleibt unverändert, unabhängige zweite Implementierung für HTTP/gRPC | kein Risiko für den bestehenden Text-Bericht | zwei unabhängige Domänenpfade für dieselbe Fähigkeit — verstößt gegen `LH-FA-SST-006`s Boundary-Kriterium und gegen den in `ADR-0130`/`ADR-0131` durchgehaltenen Grundsatz „kein zweiter Domänenpfad"; verewigt die von `ADR-0057` benannte Lücke für den CLI-Weg |
| **B — CLI-Sondermodus wird zum dünnen Wrapper um den neuen `DiagnoseUseCase` (gewählt)** | ein einziger Domänenpfad für alle drei Zugriffswege (CLI, HTTP, gRPC); der bestehende Text-Bericht bleibt byte-gleich im Format (`DiagnoseResult` trägt exakt die heute gedruckten Felder), nur die Datenbeschaffung wandert hinter den neuen Port; schließt nebenbei die von `ADR-0057` benannte Hexagon-Lücke für den CLI-Weg, den sie am längsten betraf | die bestehende `wiring.go`-Funktion `Diagnose` wird durch eine Konstruktion (`postgresstorage.NewDiagnostics` → `diagnose.NewDiagnoseService` → CLI-Formatierung) ersetzt — mechanische Umstellung, kein Verhaltensunterschied, aber ein Diff an bestehendem, produktivem Code |
| C — CLI-Sondermodus wird entfernt, nur noch HTTP/gRPC | ein Zugriffsweg weniger zu pflegen | verletzt `AGENTS.md`s Grundsatz additiver Erweiterung; betrifft Betriebspersonal ohne Netzwerkzugriff auf `CDC_HTTP_ADDR`/`CDC_GRPC_ADDR` (beide additiv-optional, `ADR-0057`/`ADR-0060`) — ein Deployment ohne beide bleibt sonst ganz ohne Diagnose-Zugriffsweg; kein Auftrag verlangt die Entfernung |

**`Healthcheck` bleibt unverändert CLI-/Docker-exklusiv** (siehe
Kontext-Abgrenzung oben) — kein neuer Port, kein neuer Zugriffsweg, weil
Dockers `HEALTHCHECK`-Mechanismus strukturell einen Prozess-Exit-Code
verlangt und die von ihm gelesene Information (Betriebsstatus) bereits
Teil von `DiagnoseResult` ist.

### Teilfrage 4 — HTTP-Exposition

Neuer Endpunkt `GET /diagnose?source=<id>`, Rechtsklasse **`reader`** —
analog zu `GetConsumerPosition`/`GetTableStatus`/`ListTables`/`GET /changes`
(`ADR-0057` Teilfrage 3, `ADR-0081` Teilfrage 4): ein reiner Lesezugriff auf
Diagnosesignale, keine administrative Wirkung. `source` ist Pflicht (fehlt
er oder ist er leer, `400` über `domainerrors.ErrEmptyIdentifier` — derselbe
Fehlerpfad wie `ReadChangesService.ReadChanges` bei leerer Quelle); wie bei
`GET /changes` (`ADR-0081` Teilfrage 4) ist die Parametermenge **geschlossen**
— ein Parameter außerhalb von `{source}` endet mit `400`, damit ein
Tippfehler im Query-String nicht still einen anderen Bericht liefert, als
der Aufrufer erwartet. Die JSON-Antwort spiegelt `DiagnoseResult` 1:1
(`snake_case`-Feldnamen, dieselbe `null`-für-„nicht gesetzt"-Konvention wie
`readChangeResponse.OldImage`/`NewImage`).

### Teilfrage 5 — gRPC-Exposition

| Option | Pro | Contra |
|---|---|---|
| A — neuer, eigenständiger Service (z. B. `Diagnostics`) in einem neuen `.proto`-Paket | saubere begriffliche Trennung „Diagnose lesen" vs. „Verwaltung ausführen" | dieselbe Abwägung wie `ADR-0131` Teilfrage 1, Option B: derselbe Server, dieselbe Adresse `CDC_GRPC_ADDR`, dieselben zwei Token-Klassen, dieselbe RPC-Form (unär) — die zwei Achsen, über die `ADR-0130` Teilfrage 1 `ChangeStream`/`Administration` getrennt hat (Stream- vs. Unary-Form, unabhängige Versionierung), unterscheiden sich hier nicht; ein zweiter `protoc`-Lauf und eine zweite generierte Go-Paket-Wurzel für einen einzelnen RPC ist Aufwand ohne belegten Bedarf |
| **B — elfter, unärer RPC `Diagnose` im bestehenden `Administration`-Service, derselben `.proto`-Datei (gewählt)** | folgt exakt dem mit `ADR-0131` bereits einmal getroffenen und bewährten Muster (`ReadChanges` als zehnter RPC im selben Service); `authUnaryInterceptor`, `administrationRPCRoles` und `administrationError` werden unverändert erweitert — ein Tabelleneintrag, kein neuer Interceptor; das HTTP-Vorbild (Teilfrage 4) legt `GET /diagnose` ebenfalls in denselben Adapter neben die zehn bestehenden Endpunkte, ohne eigenen Server | `Administration` trägt jetzt elf RPCs, davon nur neun im ursprünglich wörtlichen Sinn „Verwaltung" (`ADR-0131` hatte diesen Bedeutungs-Zuwachs bereits einmal akzeptiert) |
| C — RPC im `ChangeStream`-Service (`cdc.stream.v1`) | ein Client-Stub weniger | `ChangeStream` ist laut `ADR-0060`s eigenem Wortlaut exakt auf Live-Streaming beschränkt; ein unärer Request/Response-RPC dort widerspricht der in `ADR-0130` Teilfrage 1 bereits getroffenen und seither zweimal (`ADR-0131`, hier) bestätigten Trennung |

**Festlegung.** Wie in Option B — dieselbe Abwägung, die `ADR-0131` bereits
für `ReadChanges` getroffen hat, gilt hier unverändert.

```proto
message HeartbeatStatus {
  bool known = 1;          // false = kein Lebenszeichen
  double age_seconds = 2;  // nur gültig wenn known = true
  string error_class = 3;  // leer = Normalbetrieb (nur gültig wenn known = true)
}

message ConsumerLag {
  string consumer_id = 1;
  bool known = 2;   // false = Quelle trug noch nie eine Transaktion
  double lag = 3;
}

message RetentionBlocker {
  bool present = 1; // false = kein Blocker
  string consumer_id = 2;
  string name = 3;
  int64 acknowledged_position = 4;
  bool backlog_known = 5;
  int64 backlog = 6;
}

message BackfillTableStatus {
  string schema = 1;
  string table = 2;
  string status = 3;
  int64 rows_copied = 4;
  bool estimated_rows_known = 5;
  int64 estimated_rows = 6;
  bool warn_estimated_size = 7;
  bool warn_duration = 8;
  string error_message = 9;
}

message DiagnoseRequest {
  string source = 1;
}

message DiagnoseResponse {
  HeartbeatStatus heartbeat = 1;
  double capture_lag = 2;
  repeated ConsumerLag consumer_lags = 3;
  RetentionBlocker retention_blocker = 4;
  double storage_bytes = 5;
  repeated BackfillTableStatus backfill = 6;
}

service Administration {
  // … neun bestehende RPCs (ADR-0130), ReadChanges (ADR-0131) …
  rpc Diagnose(DiagnoseRequest) returns (DiagnoseResponse);
}
```

Explizite Präsenz-Flags (`known`/`present`/`*_known`) statt proto3
`optional` oder eines 0-als-„nicht gesetzt"-Sentinels: Anders als bei
`ADR-0131`s `from`/`to` (wo `0` als Positions-/Limit-Wert nie gültig ist)
ist `0.0` für `age_seconds`/`lag` ein real möglicher, gültiger Wert (ein
frisches Lebenszeichen kann näherungsweise `0` Sekunden alt sein) — ein
Sentinel wäre hier mehrdeutig. Die Flag-Form hält sich an den
bestehenden Stil des `.proto`-Bestands, der `optional` bisher nirgends
verwendet (geprüft durch Volltextsuche über beide `.proto`-Dateien bei
Abfassung dieser ADR).

### Teilfrage 6 — Rechtsklasse und Fehlerform (gRPC)

**`roleReader`** — wie `GetConsumerPosition`/`GetTableStatus`/`ListTables`/
`ReadChanges` (`ADR-0130` Teilfrage 4, `ADR-0131` Teilfrage 4). Ein
`admin`-Token deckt implizit ab. Zehnter → elfter Tabelleneintrag in
`administrationRPCRoles`:

```go
"Diagnose": roleReader,
```

Fehlerform — `administrationError` erweitert sich um genau einen neuen Fall
(alle übrigen Sentinels existieren bereits im Switch):

| HTTP-Status (Ursache) | gRPC-Code |
|---|---|
| `400` — fehlende/leere Quelle (`domainerrors.ErrEmptyIdentifier`) | `codes.InvalidArgument` (bereits im bestehenden Switch) |
| `400` — unbekannter Query-Parameter (HTTP-exklusiv, entfällt strukturell bei gRPC — Teilfrage 5) | entfällt strukturell |
| `401` — fehlender/unbekannter Bearer-Token | `codes.Unauthenticated` (unverändert `authUnaryInterceptor`) |
| `403` — Rechtsklasse unzureichend | entfällt praktisch: `Diagnose` verlangt nur `roleReader`, jedes gültige Token genügt (dieselbe Beobachtung wie `ADR-0131` Teilfrage 5) |
| `500` — Lesefehler an einer der fünf Diagnose-Views (`outbound.ErrDiagnosticsStorage`) | `codes.Internal` (neuer Fall in `administrationError`) |

Der Handler ruft `s.diagnose.Diagnose(ctx, inbound.DiagnoseQuery{Source: …})`
und übersetzt Erfolg/Fehler über dieselbe `administrationError`-Funktion,
die die zehn bestehenden Handler bereits verwenden.

## Verglichene Alternativen

„Nichts tun" ist hier keine ernsthafte vierte Option in den Teilfragen 3–6:
Der Auftrag benennt die Fähigkeit zweifach explizit als beauftragt — dieselbe
Lage wie bei `ADR-0131`, anders als bei `ADR-0130`/`ADR-0081`, wo der
Trigger erst neu bewertet werden musste, bevor eine positive Entscheidung
überhaupt zur Debatte stand.

## Konsequenzen

- Positiv: Ein Netzwerk-Consumer (HTTP oder gRPC) erreicht jetzt dieselben
  sechs Diagnose-Signalgruppen wie der bisherige CLI-Sondermodus — schließt
  die von `ADR-0057` benannte, seither zweimal (`ADR-0130`, `ADR-0131`)
  zitierte Lücke vollständig.
- Positiv: Die Umsetzung schließt zugleich eine bestehende Hexagon-Lücke
  (`ARC-002`/`ARC-003`): `Diagnose` verlässt seine Stellung als freie
  Funktion ohne Portabstraktion — alle drei Zugriffswege (CLI, HTTP, gRPC)
  laufen danach über denselben `DiagnoseUseCase`, keine zwei Domänenpfade.
- Positiv: Keine neue Secret-Klasse, keine neue ENV-Variable für HTTP/gRPC
  — `CDC_HTTP_ADDR`/`CDC_GRPC_ADDR` und die beiden bestehenden Token-Klassen
  tragen die Erweiterung vollständig.
- Positiv: `.a-check.yml` bleibt unverändert — der neue Inbound/Outbound
  Port und der neue Driven-Adapter-Baustein liegen bereits in den
  bestehenden, rekursiven Globs.
- Negativ: Der bestehende CLI-Sondermodus (`internal/bootstrap/wiring.go`
  `Diagnose`) wird umgebaut, nicht nur ergänzt — ein Diff an bestehendem,
  produktivem Code ohne Verhaltensänderung für den Aufrufer (Risiko: eine
  Formatierungsabweichung im Text-Bericht wäre eine reale Regression,
  Fitness Function unten).
- Negativ: `Administration` trägt jetzt elf RPCs — derselbe Bedeutungs-
  Zuwachs des Namens, den `ADR-0131` bereits einmal akzeptiert hat, jetzt
  ein zweites Mal.
- Folgepflicht: Der umsetzende Slice-Schnitt implementiert
  `internal/application/port/inbound/diagnose.go`,
  `internal/application/port/outbound/diagnostics.go`,
  `internal/adapters/driven/postgresstorage/diagnostics.go` (die fünf
  bestehenden SQL-Texte, aus `wiring.go` verschoben, kein neuer SQL-Text),
  `internal/application/usecase/diagnose/service.go`, den Umbau von
  `internal/bootstrap/wiring.go`s `Diagnose`-Funktion zu einem dünnen
  Formatierungs-Wrapper um den neuen Use Case, den HTTP-Handler
  `internal/adapters/driving/http/diagnose.go` samt Routen-Registrierung,
  die Erweiterung von `proto/cdc/administration/v1/administration.proto`
  um den elften RPC, die Regeneration des Go-Codes
  (`tools/harness/proto-generate.sh`, derselbe `protoc`-Lauf, keine neue
  Build-Stufe), den elften Handler in `administration.go`, den elften
  Tabelleneintrag in `interceptor.go` und die Config-/Wiring-Erweiterung in
  `server.go`.
- Folgepflicht: `spec/pflichtenheft.md` erweitert `SPEC-018` um
  `GET /diagnose` (Endpunkt, Rechtsklasse, JSON-Schema, Fehlercodes) und
  `SPEC-031` um den elften RPC (Nachrichtenschema, Fehlercode-Tabelle) —
  kein neuer `SPEC-*`-Eintrag nötig, beide bestehenden Einträge beschreiben
  bereits die jeweilige API als Ganzes; Gegenstand des umsetzenden Slices.
- Folgepflicht: `docs/user/benutzerhandbuch.md` §„Diagnose ausführen" trägt
  zusätzlich zum bestehenden CLI-Aufruf die beiden neuen Zugriffswege
  (`GET /diagnose`, gRPC `Diagnose`) als gleichwertige Alternativen.
- Folgepflicht: Beispiel-Clients (Go/C#/Kotlin, volle Matrix laut
  `SPEC-018`§„Sprachen und Umfang") und die Erweiterung der drei
  SDK-Packages (`PgChangeFeed.Client`, `pgchangefeed`,
  `pgchangefeed-kotlin`) um `Diagnose` — je eigener Folge-Schritt/Slice,
  nicht Gegenstand dieser ADR.
- Folgepflicht: E2E-Beleg (`make test-integration`) — ein realer
  `GET /diagnose`-Aufruf mit einem `reader`-Token und ein realer
  gRPC-`Diagnose`-Aufruf über `tools/harness/grpcadminclient` liefern
  denselben Betriebsstatus, der auch über den bestehenden
  `docker exec … diagnose`-Rundlauf sichtbar ist (Querabgleich derselben
  Quelle über alle drei Zugriffswege) — analog zum bestehenden HTTP-/gRPC-
  Rundlauf-Vorbild.
- Folgepflicht: `spec/architecture.md` braucht **keine** Änderung — `ARC-005`
  ist bereits durch `ADR-0057`/`ADR-0081`/`ADR-0130` konkret gemacht; diese
  ADR öffnet über beide bereits konkretisierten Driving-Adapter eine
  weitere, bereits bestehende Fähigkeit.
- Nicht Gegenstand dieser ADR: `Healthcheck` (Dockers binäres Verdikt) bleibt
  CLI-/Docker-exklusiv (Teilfrage 3) — kein neuer Port, kein neuer
  Zugriffsweg dafür.

## Fitness Function (falls maschinell prüfbar; Erwartung an den umsetzenden Slice, nicht durch diesen Architect-Lauf gefahren)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Unit-Test (`internal/application/usecase/diagnose`, Whitebox) | `DiagnoseService.Diagnose` ruft exakt `DiagnosticsPort.Read` mit der Quelle aus der Query auf (Fake-Port im Test) und gibt dessen Ergebnis unverändert als `DiagnoseResult` zurück — keine eigene Abfrage, keine zweite Sortierung | `make test` |
| Go-Unit-Test (`internal/bootstrap`, Whitebox, `Diagnose`) | Der Text-Bericht der umgebauten `Diagnose`-Funktion ist für dieselbe Eingabe byte-gleich mit dem heutigen, vor dieser ADR erzeugten Text (mindestens eine Instanz: derselbe Testfall wie `main_test.go`s bestehender Diagnose-Test) — belegt, dass der Umbau (Teilfrage 3) keine sichtbare Verhaltensänderung für den CLI-Aufrufer trägt | `make test` |
| Go-Unit-Test (`internal/adapters/driving/http`, Whitebox, `diagnoseHandler`) | Ein Aufruf ohne `source`-Parameter oder mit einem zusätzlichen, nicht gelisteten Parameter endet mit `400`; ein Aufruf ohne oder mit unbekanntem Bearer-Token endet mit `401`; ein gültiges `reader`- oder `admin`-Token erreicht den Endpunkt | `make test` |
| Go-Unit-Test (`internal/adapters/driving/grpc`, Whitebox, `authUnaryInterceptor`/`Diagnose`-Handler) | Ein gültiges `reader`-Token erreicht `Diagnose`; ein gültiges `admin`-Token ebenfalls; ein Aufruf ohne oder mit unbekanntem Token endet mit `codes.Unauthenticated`; der Handler ruft exakt `inbound.DiagnoseUseCase.Diagnose` auf (Fake-Use-Case im Test) | `make test` |
| `.a-check` | unverändert — der neue Inbound/Outbound Port und der neue Driven-Adapter-Baustein liegen vollständig in den bestehenden, rekursiven Globs `ports`/`adapters`/`contract`, kein neuer Hexagon-Schichten-Edge nötig (geprüft durch Lesen von `.a-check.yml` bei Abfassung dieser ADR) | `make a-check` |
| `make test-integration` | realer `GET /diagnose`-Aufruf mit `reader`-Token und realer gRPC-`Diagnose`-Aufruf gegen den laufenden Feed-Container liefern denselben Betriebsstatus wie der bestehende `docker exec … diagnose`-Rundlauf (Querabgleich); ein Aufruf ohne Token endet `401`/`Unauthenticated` — Umsetzung Gegenstand des umsetzenden Slice-Schnitts | `make test-integration` |

## Re-Evaluierungs-Trigger

Permanent — die Wahl „ein Inbound Port `DiagnoseUseCase`, ein Outbound Port
`DiagnosticsPort`, CLI als dünner Wrapper, `GET /diagnose` mit `reader`-
Klasse, elfter RPC im bestehenden `Administration`-Service" gilt unabhängig
vom Umsetzungszeitpunkt. Ausnahmen, die eine eigene Folge-ADR brauchen: (1)
ein konkret benannter Bedarf, `Healthcheck` (das binäre Docker-Verdikt)
ebenfalls über HTTP/gRPC zu exponieren — die in Teilfrage 3/Kontext
getroffene Abgrenzung neu bewerten; (2) ein Bedarf an einer eigenen
Rechtsklasse für `Diagnose` jenseits `reader`/`admin`; (3) ein praktisches
Wachstums-Limit des `Administration`-Service macht eine nachträgliche
Aufspaltung nötig (dieselbe Ausnahme wie bereits in `ADR-0131`
vorgezeichnet).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-28 | Accepted — Architect-Entscheidung nach zweifach explizit bestätigtem Auftraggeber-Auftrag; Supersedes `ADR-0057` Teilfrage 2/Konsequenzen in ihrer Diagnose/Health-Hälfte, teilweise (die Changes-Lesen-Hälfte derselben Klauseln hatte bereits `ADR-0081` abgelöst) | [`ADR-0057`](0057-http-grpc-api.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0132` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
