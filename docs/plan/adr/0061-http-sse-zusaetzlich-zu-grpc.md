# ADR-0061: HTTP/SSE als zweiter, paralleler Zustellweg zusätzlich zu gRPC

**Status:** Accepted

**Datum:** 2026-09-14

**Autor:** pt9912 (Rolleninhaber: Architect-Lauf, 2026-09-14)

**Bezug:** [`LH-FA-SST-008`](../../../spec/lastenheft.md) (Haupt-Bezug —
Live-Streaming vollständiger Change-Inhalte, seit Version 0.9.0
protokollneutral formuliert und „gRPC, HTTP/SSE, oder beide nebeneinander"
ausdrücklich als Architektur-/Spezifikationsfrage benennend),
[`LH-FA-SST-007`](../../../spec/lastenheft.md) (Abgrenzung — NATS bleibt
eigenständiges Wecksignal, unverändert durch diese ADR),
[`LH-FA-REA-001`](../../../spec/lastenheft.md) (bestehender
Lesezugriffsweg als Nachhol-Pfad), [`LH-FA-CON-003`](../../../spec/lastenheft.md),
[`LH-FA-CON-005`](../../../spec/lastenheft.md) (bestätigte Consumer-Position
als Fortsetzungs-Referenz), [ADR-0060](0060-grpc-streaming-mechanismus.md)
(Ergänzung, **keine** Supersedes — der dort gewählte gRPC-Weg, der neue
Outbound Port `ChangeStreamPort` und der In-Prozess-`Broadcaster` bleiben
unverändert gültig und werden von dieser ADR wiederverwendet, nicht
ersetzt), [ADR-0057](0057-http-grpc-api.md) (Adapter-Basis — der
bestehende HTTP-Driving-Adapter `internal/adapters/driving/http/` samt
Token-Middleware, auf dem der neue SSE-Endpunkt aufbaut),
[ADR-0034](0034-ports-nach-faehigkeiten.md) (Port-Zuschnitt nach Fähigkeit
— Grundlage dafür, dass ein zweiter Zustellweg für dieselbe Fähigkeit
keinen zweiten Port braucht), [ADR-0055](0055-nats-change-notification-wecksignal.md)
(Zustellsemantik-Präzedenzfall, den `ADR-0060` bereits für Live-Streaming
übernommen hat und den diese ADR ein zweites Mal bestätigt)

**Schärft:** [`ARC-005`](../../../spec/architecture.md) (macht das dort
bereits generisch genannte „später HTTP-/gRPC" für **denselben**
Live-Streaming-Bedarf ein zweites Mal konkret — nach `ADR-0060`s gRPC-Wahl
jetzt zusätzlich HTTP/SSE über den bereits bestehenden `ADR-0057`-Adapter)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0060`](0060-grpc-streaming-mechanismus.md) traf für
`LH-FA-SST-008` — zum damaligen Zeitpunkt in Lastenheft-Version 0.8.0 noch
gRPC-spezifisch formuliert („Live-Streaming neuer Changes über gRPC") —
die Entscheidung für gRPC Server-Streaming und verwarf in Teilfrage 1
Option B (HTTP/SSE) explizit, mit der Begründung, SSE sei „reine
Ersatzkonstruktion" ohne typisiertes Nachrichtenschema. Diese Verwerfung
galt für eine **Entweder-oder-Frage**: Unter der damals gRPC-spezifischen
Fassung von `LH-FA-SST-008` konnte es nur einen Mechanismus geben, und
gRPC gewann den Vergleich.

Am selben Tag (2026-09-14) wurde `LH-FA-SST-008` auf Version 0.9.0
protokollneutral umformuliert: Titel, Beschreibung und Out-of-Scope
nennen jetzt ausdrücklich „gRPC, HTTP/SSE, oder beide nebeneinander" als
offene Architektur-/Spezifikationsfrage. Der Auftraggeber (pt9912) hat
danach explizit den Bedarf benannt, HTTP/SSE **zusätzlich** zu `ADR-0060`s
gRPC-Weg einzuführen — nicht als Ersatz.

Diese ADR ist **keine Korrektur und keine Supersedes-ADR zu `ADR-0060`**.
Der Unterschied zur ursprünglichen Entweder-oder-Frage ist kategorial,
nicht graduell: `ADR-0060` Teilfrage 1 beantwortete „welcher **einzige**
Mechanismus implementiert `LH-FA-SST-008`" unter einer Anforderung, die
zum Entscheidungszeitpunkt nur einen Mechanismus zuließ. Diese ADR
beantwortet eine andere Frage — „wie kommt ein **zweiter, paralleler**
Mechanismus zur bereits getroffenen gRPC-Entscheidung hinzu" —, die als
Option unter der alten Anforderungsfassung gar nicht existierte. `ADR-0060`s
Textkörper, seine sechs Teilfragen und ihre Verdikte bleiben unangetastet
und in Kraft; kein Zeilen-Diff an dieser Datei ist Folge dieser ADR
(Baseline-Regelwerk `modul-08-agentenrollen.md` §Konflikt-Pfad als
Rollen-Sequenz, Verdikt 2: Folge-ADR statt stiller Korrektur — hier
sogar strenger, weil diese ADR nicht einmal behauptet, `ADR-0060`s
Vergleich sei falsch gewesen; sie fügt eine Option hinzu, die zum
Vergleichszeitpunkt nicht zur Wahl stand).

Drei Konstraints prägen den Lösungsraum, alle drei bereits durch
`ADR-0060` bzw. `ADR-0057` vorgeprägt:

- **Der neue Outbound Port `ChangeStreamPort` und der `Broadcaster`
  existieren bereits als architektonische Entscheidung** (`ADR-0060`
  Teilfrage 2): `Publish(ctx, change)` je einzelnem `Change`,
  `Subscribe() (<-chan *model.Change, cancel func())` je Abonnent,
  Fire-and-Forget ohne Puffer. Beide Signaturen sind bereits
  protokoll-agnostisch — kein gRPC-spezifischer Typ taucht in ihnen auf.
- **Der HTTP-Driving-Adapter `internal/adapters/driving/http/` existiert
  bereits** (`ADR-0057`) mit einer Token-Middleware (`withToken` in
  `middleware.go`), die zwei Rechtsklassen gegen `CDC_API_TOKEN_READER`/
  `ADMIN` prüft und danach unverändert an den nächsten `http.Handler`
  delegiert (`next.ServeHTTP(w, r)`) — sie liest oder puffert den
  Response-Body nicht selbst.
- **`LH-FA-SST-008`s Boundary- und Negative-Akzeptanzkriterien gelten für
  jeden Zustellweg gleichermaßen** — kein Weg darf Changes bei
  Unterbrechung verlustbehaftet machen (Store bleibt Nachhol-Pfad) und
  keiner darf einen Verbindungsversuch ohne gültige Authentifizierung
  stillschweigend fortsetzen.

## Entscheidung

Wir wählen **HTTP/SSE als zweiten Endpunkt im bereits bestehenden
`internal/adapters/driving/http/`-Adapter, gespeist vom selben
`ChangeStreamPort`/`Broadcaster` aus `ADR-0060` als zweiter Abonnent** —
mit sechs Teilfragen.

### Teilfrage 1 — Gemeinsame Quelle für beide Zustellwege

| Option | Pro | Contra |
|---|---|---|
| A — SSE-Endpunkt pollt `ChangeStorePort` periodisch statt Push vom Capture-Pfad (nichts Neues bauen) | keine Änderung an `CaptureService` oder `Broadcaster` nötig | reintroduziert genau das Pollen, das `LH-FA-SST-008` ausdrücklich vermeiden will — derselbe Einwand wie `ADR-0060` Teilfrage 2 Option C |
| **B — SSE nutzt denselben `Broadcaster`/`ChangeStreamPort` aus `ADR-0060` als zweiten Abonnenten (gewählt)** | `Broadcaster.Subscribe()` liefert bereits nur `<-chan *model.Change, cancel func()` — kein gRPC-spezifischer Typ, also bereits protokoll-agnostisch; ein SSE-Client registriert sich über genau denselben Aufruf wie ein gRPC-Stream-Client; `CaptureService` publiziert weiterhin genau einmal je `Change` (`ADR-0060` Teilfrage 2), unabhängig davon, wie viele Abonnenten — gRPC- und SSE-Clients zusammen — gerade lauschen; keine Änderung an `ADR-0060`s Kern-Design nötig | keine wesentlichen — die einzige neue Kopplung ist ein weiterer `Subscribe()`-Aufrufer im `driving/http`-Paket |
| C — eigener, zweiter Broadcaster ausschließlich für SSE, `CaptureService` erhält eine zweite Konstruktions-Option (`WithSSEStream` neben `WithChangeStream`) | Isolation zwischen den beiden Zustellwegen | dupliziert einen bereits protokoll-agnostischen Mechanismus ohne fachlichen Grund — zwei In-Prozess-Fan-outs, die dieselbe Aufgabe an denselben Change-Strom lösen; verdoppelt die Publish-Kosten je `Change` und die Anzahl der Konstruktions-Optionen an `CaptureService`, ohne dass ein Akzeptanzkriterium diese Trennung verlangt |
| D — SSE läuft über einen gRPC-Gateway-Übersetzer (z. B. grpc-gateway), der den bestehenden gRPC-Stream in SSE-Framing übersetzt | ein einziger interner Erzeuger-Pfad (nur gRPC) | zieht die Protobuf-/gRPC-Toolchain-Abhängigkeit in den HTTP-Adapter, den `ADR-0057` bewusst ohne Fremd-Abhängigkeit hielt; widerspricht dem Zweck eines zweiten, *einfacheren* Zustellwegs — ein Consumer, der SSE gerade wählt, um Codegen zu vermeiden, würde über den Umweg doch wieder an das gRPC-Schema gekoppelt |

Die gewählte Option ändert **nichts** an `ADR-0060`s Port-/Adapter-Design
(Teilfrage 2 dort): `ChangeStreamPort`, `Broadcaster` und die
`CaptureService`-Option `WithChangeStream` bleiben exakt wie dort
festgelegt. Diese ADR fügt lediglich einen zweiten Aufrufer von
`Broadcaster.Subscribe()` hinzu.

### Teilfrage 2 — Adapter-Platzierung

| Option | Pro | Contra |
|---|---|---|
| A — neues, drittes Adapter-Paket (z. B. `internal/adapters/driving/sse/`) mit eigenem Server/Port | saubere Paket-Trennung nach Protokoll | dritte Horch-Adresse und dritte Konfigurationsvariable für einen Endpunkt, der denselben Trägermechanismus (`net/http`) und dieselbe Token-Middleware wie der bestehende Adapter braucht — reine Verdopplung von Infrastruktur, die bereits existiert; ein SSE-Client müsste einen zweiten Port kennen, obwohl der HTTP-Server aus `ADR-0057` bereits läuft |
| **B — neuer Endpunkt (`GET /changes/stream`) im bestehenden `internal/adapters/driving/http/`-Paket, über denselben `http.ServeMux` und dieselbe `withToken`-Middleware verdrahtet (gewählt)** | keine neue Horch-Adresse, kein neues Konfigurationsfeld für die Adresse; `withToken` (`middleware.go`) prüft das Token und delegiert danach unverändert an `next.ServeHTTP(w, r)` — für eine lang laufende Verbindung ändert sich daran nichts, die Middleware liest oder puffert den Response-Body nicht; das eigentliche Event-Framing (`data: …\n\n`) und das Flushing nach jedem Event (`http.Flusher`) sind Sache des SSE-Handlers selbst, nicht der Middleware — dieselbe Trennung wie bei jedem anderen Handler in diesem Paket | ein Handler in diesem Paket hält jetzt eine lang laufende Verbindung offen, während alle anderen Handler Request/Response sind — dokumentationspflichtiger Unterschied innerhalb desselben Pakets, aber kein struktureller |
| C — eigene `http.Server`-Instanz auf demselben Port, aber in einem separaten Package, das den bestehenden Adapter importiert | vermeintliche Trennung | ein zweiter `http.Server` kann nicht denselben Port binden wie der erste — entweder Reverse-Proxy-Aufwand ohne Gegenwert, oder es ist ohnehin derselbe `ServeMux`, nur über einen unnötigen Import-Umweg — kein Vorteil gegenüber Option B, zusätzliche Indirektion |

Die Middleware-Wiederverwendung ist geprüft, nicht angenommen:
`withToken` (`internal/adapters/driving/http/middleware.go`) klassifiziert
das Bearer-Token und ruft bei Erfolg `next.ServeHTTP(w, r)` — der Handler
dahinter erhält den vollen `http.ResponseWriter` und kann selbst per
`http.Flusher`-Type-Assertion nach jedem Event flushen. Es gibt keinen
Schritt in `withToken`, der eine lang laufende Verbindung stört (kein
Response-Buffering, kein Timeout-Wrapping). Der Adapter importiert dafür
weiterhin ausschließlich Inbound Ports und Domain-Typen — der neue
Handler importiert zusätzlich die lokal deklarierte
`changeSubscriber`-Schnittstelle aus `ADR-0060` Teilfrage 2 (kein Import
des `grpcstream`-Pakets selbst; Bootstrap verdrahtet den konkreten
`*grpcstream.Broadcaster` in beide Adapter, exakt wie in `ADR-0060`
Teilfrage 2 für den gRPC-Server beschrieben).

### Teilfrage 3 — Zustellsemantik und Last-Event-ID

| Option | Pro | Contra |
|---|---|---|
| A — Replay ab `Last-Event-ID` beim (Re-)Connect, gegen den Store aufgelöst | SSE-natives Reconnect-Resume-Feature würde genutzt, Consumer bekäme lückenlose Zustellung ohne selbst nachzuholen | baut exakt die zweite Nachvollziehbarkeits-Infrastruktur auf, die `ADR-0060` Teilfrage 3 Option B bereits verworfen hat (dupliziert den `ChangeStorePort`-Lesepfad); der `Broadcaster` hat bewusst keinen Puffer (`ADR-0060` Teilfrage 3) — ein `Last-Event-ID`-Replay bräuchte zusätzlich einen Store-Lookup samt Übergang in den Live-Strom, den kein Akzeptanzkriterium von `LH-FA-SST-008` verlangt; würde zudem die beiden Zustellwege inkonsistent machen (gRPC ohne Replay, SSE mit) |
| **B — Fire-and-Forget ohne Replay, `Last-Event-ID` bewusst ungenutzt (gewählt)** | identische Zustellsemantik wie `ADR-0060` Teilfrage 3 für den gRPC-Weg — ein Consumer sieht dasselbe Verhalten unabhängig vom gewählten Protokoll; erfüllt `LH-FA-SST-008`s Boundary-Kriterium exakt (Store bleibt Nachhol-Pfad über `LH-FA-REA-001` ff. und die bestätigte Consumer-Position `LH-FA-CON-003`/`005`); kein zusätzlicher Zustand im `Broadcaster` | ein SSE-Client, der kurz getrennt war, verpasst Changes über den Stream ersatzlos — muss wie ein gRPC-Client selbst über den bestehenden Lesezugriffsweg nachholen; das native `Last-Event-ID`-Feature des Protokolls bleibt ungenutzt |
| C — kleiner In-Memory-Ring-Buffer nur für den SSE-Weg (asymmetrisch zu gRPC) | überbrückt sehr kurze Verbindungslücken für SSE-Clients | macht die beiden parallelen Zustellwege für **dieselbe** Fähigkeit inkonsistent — ein Consumer, der zwischen gRPC und SSE wechselt, sähe unterschiedliches Verhalten ohne fachlichen Grund; dieselbe unvorhersagbare Zusicherung, die `ADR-0060` Teilfrage 3 Option C bereits verworfen hat |

Der Unterschied zu `ADR-0060` ist rein die Existenz des
`Last-Event-ID`-Headers als Protokoll-Feature — er wird von Server und
Handler ignoriert (kein Aussenden einer `id:`-Zeile, kein Auswerten eines
eingehenden `Last-Event-ID`-Headers), damit ein SSE-Client sich nicht auf
ein Resume-Verhalten verlässt, das dieser Adapter nicht anbietet. Ein
späterer Bedarf an echtem Replay bräuchte, wie bei `ADR-0060`, eine
eigene Folge-ADR (siehe Re-Evaluierungs-Trigger).

### Teilfrage 4 — Authentifizierung/Autorisierung

| Option | Pro | Contra |
|---|---|---|
| A — kein Schutz für den SSE-Endpunkt (nichts tun) | kein Zusatzaufwand | `LH-FA-SST-008`s Negative-Akzeptanzkriterium wird strukturell unerfüllbar; inkonsistent mit dem geschützten gRPC-Weg für dieselbe Fähigkeit |
| B — eigener, neuer Token-Mechanismus (z. B. `CDC_SSE_TOKEN`) | vom übrigen HTTP-Tokenraum entkoppelt | eine dritte Secret-Klasse neben `CDC_API_TOKEN_READER`/`ADMIN` und der bereits für den gRPC-Weg wiederverwendeten Klasse (`ADR-0060` Teilfrage 4) ohne fachlichen Grund — dieselbe Einwand-Struktur wie dort |
| **C — Wiederverwendung derselben `CDC_API_TOKEN_READER`/`CDC_API_TOKEN_ADMIN`-Klassen über den bestehenden `Authorization: Bearer`-Header und `withToken` (gewählt)** | konsistentes Least-Privilege-Modell über alle drei Netzwerk-Zugriffswege hinweg (HTTP-Verwaltungs-API `ADR-0057`, gRPC-Stream `ADR-0060`, jetzt SSE-Stream); keine neue Secret-Klasse; ein `reader`-Token genügt, da Streaming rein lesend ist (dieselbe Einordnung wie bei `ADR-0060` Teilfrage 4); der bestehende `withToken`-Wrapper wird unverändert eingesetzt, kein neuer Middleware-Pfad | Tokens bleiben Shared Secrets ohne Ablauf-/Widerrufsmechanismus — dieselbe, bereits in `ADR-0057`/`ADR-0060` benannte Einschränkung; kein Browser-`EventSource`-Header-Limit relevant, da kein Browser-Konsument gefordert ist und ein selbstgebauter Client den `Authorization`-Header frei setzen kann |

Ein Verbindungsversuch ohne `Authorization`-Header oder mit einem Token,
das keiner der beiden konfigurierten Klassen entspricht, endet mit `401`
— vor jedem Schreiben eines SSE-Events, also bevor der Stream überhaupt
geöffnet wird; sichtbar, nicht still mit leeren Daten fortgesetzt
(`LH-FA-SST-008` Negative).

### Teilfrage 5 — Bootstrap-Aktivierung

| Option | Pro | Contra |
|---|---|---|
| A — eigenes Feature-Gate/ENV (z. B. `CDC_HTTP_SSE_ENABLED`) zusätzlich zu `CDC_HTTP_ADDR` | expliziter Opt-in | ein zusätzlicher Konfigurationsparameter für einen Endpunkt, der ohnehin nur erreichbar ist, wenn `CDC_HTTP_ADDR` bereits gesetzt und der Server bereits gestartet ist; kein anderer Endpunkt in diesem Adapter (`ADR-0057`) trägt ein eigenes Pro-Endpunkt-Gate — inkonsistent zum etablierten Muster „additive Route im selben `mux`" |
| **B — Endpunkt läuft automatisch mit, sobald `CDC_HTTP_ADDR` gesetzt ist; kein zusätzliches ENV (gewählt)** | folgt demselben Muster wie jeder andere Endpunkt in `internal/adapters/driving/http/` — eine zusätzliche `mux.Handle(...)`-Zeile, kein neues Konfigurationsfeld; jedes bestehende Deployment, das `CDC_HTTP_ADDR` bereits nutzt, bekommt den Endpunkt additiv dazu — kein Breaking Change, keine bestehende Route ändert sich | ein Betreiber, der `CDC_HTTP_ADDR` ausschließlich für die Verwaltungs-API gesetzt hatte, bekommt den Streaming-Endpunkt ungefragt mit — bewusst in Kauf genommen, da derselbe Token-Schutz gilt wie für jeden anderen Endpunkt dieses Adapters und kein zusätzlicher Ressourcenverbrauch entsteht, solange kein Client den Endpunkt tatsächlich verbindet |
| C — Endpunkt registriert, antwortet aber mit `503`, solange kein `ChangeStreamPort`/`Broadcaster` verdrahtet ist | robust gegen eine Bootstrap-Reihenfolge, in der `CDC_HTTP_ADDR` ohne aktive Streaming-Fähigkeit gesetzt ist | das ist keine Alternative zu B, sondern eine Ergänzung dazu — als Verhaltensregel unter Konsequenzen/Folgepflicht festgehalten, nicht als eigenständige Option, da sie dieselbe Grundentscheidung (kein eigenes ENV) nicht verändert |

Damit `LH-FA-SST-008` per SSE auch nutzbar ist, ohne dass ein Betreiber
zwingend auch den gRPC-Weg aktiviert (`CDC_GRPC_ADDR` gesetzt), muss die
Konstruktion des `Broadcaster` in `internal/bootstrap` von `CDC_GRPC_ADDR`
entkoppelt werden: Der `Broadcaster` wird konstruiert und über
`CaptureService.WithChangeStream` verdrahtet, sobald **mindestens einer**
der beiden Zustellwege aktiv ist (`CDC_GRPC_ADDR` gesetzt **oder**
`CDC_HTTP_ADDR` gesetzt) — nicht nur, wenn beide oder nur der gRPC-Weg
aktiv sind. Ist keiner der beiden gesetzt, bleibt die Fähigkeit
vollständig deaktiviert, bestehendes Verhalten bit-identisch (dieselbe
Zusicherung wie `ADR-0060` Teilfrage 6 für den bislang einzigen Weg).
Ist `CDC_HTTP_ADDR` gesetzt, aber der `Broadcaster` aus einem anderen
Grund nicht verdrahtet, antwortet der Endpunkt mit `503` statt eine leere
oder blockierende Verbindung offenzuhalten (Option C oben, als
Folgepflicht des umsetzenden Slices).

### Teilfrage 6 — Verhältnis der beiden Zustellwege zueinander

| Option | Pro | Contra |
|---|---|---|
| A — keine Empfehlung dokumentiert, beide vollständig gleichwertig ohne Hinweis | kein zusätzlicher Text | ein Betreiber ohne Kontext muss die Abwägung (Codegen vs. Einfachheit) selbst aus den beiden ADRs rekonstruieren, obwohl `LH-FA-SST-008` selbst keine Präferenz verlangt |
| **B — beide unabhängig nutzbar, ohne Präferenz- oder Fallback-Kopplung; dokumentierte Empfehlung als Orientierung (gewählt)** | `LH-FA-SST-008` verlangt keine Bevorzugung — beide Wege bleiben architektonisch gleichrangig, kein automatischer Wechsel, kein gemeinsamer Reconnect-Mechanismus zwischen ihnen; die Empfehlung („SSE für einfache HTTP-only-Clients ohne Codegen-Bedarf, gRPC für typisierte/produktionsnahe Consumer mit bereits vorhandener Protobuf-Toolchain") hilft ohne etwas zu erzwingen | Empfehlung kann veralten, wenn sich Consumer-Muster ändern — dafür trägt sie keinen eigenen Mechanismus, nur Dokumentationstext |
| C — strikte Präferenz/Deprecation-Pfad (SSE als Übergangslösung, gRPC als Zielzustand oder umgekehrt) | klare Ansage für künftige Consumer | widerspricht dem expliziten Auftrag „zusätzlich, nicht als Ersatz"; würde einen der beiden von `ADR-0060` bzw. dieser ADR getroffenen, gleichrangigen Entscheidungen ohne fachlichen Grund abwerten |

Ein Consumer wählt genau einen der beiden Wege für seine jeweilige
Verbindung; nichts in dieser ADR verlangt, dass ein Consumer beide
gleichzeitig nutzt oder dass ein Server einen Wechsel zwischen ihnen
vermittelt.

## Konsequenzen

- Positiv: `LH-FA-SST-008` wird über einen zweiten Zustellweg erfüllbar,
  ohne `ADR-0060`s gRPC-Entscheidung, den `ChangeStreamPort`-Vertrag oder
  den `Broadcaster` zu ändern — beide Wege teilen sich denselben
  Erzeuger-Pfad (`ADR-0034`: ein Port pro Fähigkeit, nicht pro
  Zustellweg).
- Positiv: Kein neuer Fremd-Abhängigkeits-Fußabdruck — SSE läuft
  vollständig über die bereits vorhandene `net/http`-Basis (`ADR-0057`)
  und den bereits vorhandenen `Broadcaster` (`ADR-0060`); im Unterschied
  zu `ADR-0060`s neuer Protobuf-/gRPC-Toolchain-Abhängigkeit entsteht hier
  keine.
- Positiv: Das Least-Privilege-Modell (`ADR-0047`, `ADR-0057`, `ADR-0060`)
  setzt sich über einen vierten Zugriffsweg fort (Verwaltungs-API,
  gRPC-Stream, jetzt SSE-Stream), weiterhin ohne eine dritte Secret-Klasse.
- Negativ: Ohne Stream-internes Replay muss auch jeder SSE-Consumer
  eigenständig einen Resync-Pfad über den bestehenden Lesezugriffsweg
  implementieren — dieselbe bewusste Entscheidung wie bei `ADR-0060`,
  jetzt für einen zweiten Weg bestätigt, keine übersehene Lücke.
- Negativ: Zwei parallele Zustellwege für dieselbe Fähigkeit bedeuten
  zwei Testpfade in `make test-integration` und zwei Beispiel-Clients zu
  pflegen — der Preis für „zusätzlich, nicht als Ersatz".
- Negativ: `internal/bootstrap` verdrahtet den `Broadcaster` jetzt gegen
  eine Oder-Bedingung zweier Umgebungsvariablen (`CDC_GRPC_ADDR` oder
  `CDC_HTTP_ADDR`) statt einer einzigen — geringfügig höhere
  Verdrahtungs-Komplexität, dokumentiert in Teilfrage 5.
- Folgepflicht: Der umsetzende Slice implementiert den Endpunkt
  `GET /changes/stream` (`text/event-stream`) im bestehenden
  `internal/adapters/driving/http/`-Paket, die Bootstrap-Entkopplung der
  `Broadcaster`-Konstruktion von `CDC_GRPC_ADDR` (Teilfrage 5) und den
  `503`-Pfad bei fehlender Verdrahtung.
- Folgepflicht: `spec/pflichtenheft.md` erhält einen neuen `SPEC-*`-Eintrag
  für das konkrete SSE-Nachrichtenschema (Event-Feld-Layout, `event:`-Typ,
  falls verwendet) — analog zur bereits in `ADR-0060` benannten
  Folgepflicht für das Protobuf-Schema, hier für SSE.
- Folgepflicht: `spec/architecture.md` braucht **keine** Änderung —
  `ARC-005` nennt HTTP/gRPC bereits generisch als Driving-Adapter-
  Technologie; diese ADR macht die Wahl für einen zweiten Zustellweg
  derselben Fähigkeit konkret, ohne die Sicht selbst zu ändern.
- Folgepflicht: `.a-check.yml` braucht **keine** Änderung — der neue
  Endpunkt liegt vollständig im bestehenden Glob
  `adapters: ["internal/adapters/**"]`; keine neue Kante zwischen den
  Paketen `driving/http` und `driven/grpcstream` jenseits der bereits in
  `ADR-0060` Teilfrage 2 beschriebenen lokal deklarierten Schnittstelle.
- Folgepflicht: Erweiterung um `Last-Event-ID`-Replay oder
  tabellen-granulare Filterung des SSE-Streams braucht je eine eigene
  Folge-ADR — nicht Gegenstand dieser ADR (siehe Re-Evaluierungs-Trigger).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Unit-Test (`internal/adapters/driving/http`, Whitebox) | Ein `GET /changes/stream`-Aufruf ohne oder mit unbekanntem Bearer-Token endet mit `401`, vor jedem geschriebenen SSE-Event; ein gültiges `reader`- oder `admin`-Token öffnet den Stream erfolgreich | `make test` |
| Go-Unit-Test (`internal/adapters/driving/http`, Whitebox) | Ein über `Broadcaster.Publish` verteilter `Change` erreicht einen verbundenen SSE-Testclient als vollständiges Event; ein `Publish` ohne verbundenen SSE-Client blockiert nicht (Fire-and-Forget-Regressionstest, spiegelt `ADR-0060`s Test für den gRPC-Weg) | `make test` |
| Go-Unit-Test (`internal/bootstrap`, Whitebox) | Der `Broadcaster` wird genau dann konstruiert und verdrahtet, wenn `CDC_GRPC_ADDR` oder `CDC_HTTP_ADDR` gesetzt ist — keines von beiden gesetzt lässt `CaptureService` ohne `ChangeStreamPort`, bestehendes Verhalten bit-identisch | `make test` |
| `.a-check` | unverändert — der neue Endpunkt liegt vollständig im bestehenden Glob `adapters: ["internal/adapters/**"]`, kein neuer Hexagon-Schichten-Edge nötig | `make a-check` |
| `make test-integration` | realer SSE-Rundlauf gegen den laufenden Feed-Container: eine committed Änderung erreicht einen verbundenen SSE-Client mit vollständigem Inhalt; ein Verbindungsversuch ohne gültiges Token wird mit `401` abgelehnt — Umsetzung Gegenstand des umsetzenden Slice-Schnitts | `make test-integration` |

## Slice-Schnitt-Empfehlung

Regeln dieser Sektion: Empfehlung, endgültiger Schnitt liegt bei der
Planner-Rolle (Baseline-Regelwerk `modul-05-planning-harness.md`).

Ein einzelner Slice genügt: SSE baut vollständig auf bereits an anderer
Stelle entschiedener bzw. vorhandener Infrastruktur auf (bestehender
HTTP-Adapter aus `ADR-0057`, `ChangeStreamPort`/`Broadcaster` aus
`ADR-0060`) und bringt selbst keine neue Toolchain-Abhängigkeit mit —
anders als `ADR-0060`s Dreiteilung, die eine neue Protobuf-Toolchain
schrittweise einführen musste.

**Abhängigkeit — explizit benannt:** Dieser Slice kann erst beginnen,
wenn `ChangeStreamPort` und der `Broadcaster` **existieren und in
`CaptureService` verdrahtet sind** — das ist `ADR-0060`s Slice A
(Adapter-Grundgerüst: Port, `Broadcaster`, Bootstrap-Verdrahtung) **und**
Slice B (Capture-Integration: `WithChangeStream`, tatsächliche
Publish-Kette). Ohne Slice B veröffentlicht `CaptureService` nie einen
`Change` über den `Broadcaster` — ein SSE-Endpunkt ohne diese Verdrahtung
wäre ein Stream, der nie ein Event sieht. Dieser Slice braucht dagegen
**nicht** `ADR-0060`s Slice C (gRPC-Beispiel-Client, gRPC-E2E-Erweiterung)
— der gRPC-Beleg und der SSE-Beleg sind voneinander unabhängige
Nachweisformen für zwei verschiedene Konsumenten desselben
`Broadcaster`.

Inhalt des einen Slices: Endpunkt `GET /changes/stream` im bestehenden
Adapter, Bootstrap-Entkopplung der `Broadcaster`-Konstruktion von
`CDC_GRPC_ADDR` (Teilfrage 5), Unit-Tests (Auth-Pfad, Fire-and-Forget),
Beispiel-Client (`tools/harness/sseclient/`, Wegwerf-Werkzeug außerhalb
der Produktionsschichten, analog zu `tools/harness/httpclient/`/
`tools/harness/grpcclient/`) und E2E-Erweiterung
(`run-integration-tests.sh`-Rundlauf über einen echten SSE-Verbindungsaufbau).

## Re-Evaluierungs-Trigger

Wird ein Bedarf an `Last-Event-ID`-gestütztem Replay (Fortsetzung ab
einer Consumer-Position **innerhalb** des SSE-Streams selbst, ohne
Rückgriff auf den SQL-Lesezugriffsweg) oder an tabellen-granularer
Filterung des SSE-Streams (analog zu `ADR-0056` für NATS, bzw. zu
`ADR-0060`s gleichlautendem Trigger für den gRPC-Weg) konkret benannt:
Beide widersprechen der hier getroffenen Wahl (Teilfrage 3 bzw. der unter
Konsequenzen benannten Vertagung) fundamental genug, dass sie eine eigene
Folge-ADR brauchen, die die jeweilige Option neu bewertet — nicht als
Korrektur dieser ADR, sondern als eigene Entscheidung, weil sich die
Prämisse (kein belegter Bedarf über das Boundary-Kriterium hinaus)
geändert hätte. Sonst permanent — die Wahl HTTP/SSE als zweiter,
paralleler Zustellweg neben gRPC gilt unabhängig vom Zeitpunkt der
Umsetzung.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-14 | Accepted — Architect-Entscheidung nach explizitem Auftraggeber-Bedarf (2026-09-14); ergänzt `ADR-0060` um einen zweiten, parallelen Zustellweg für dieselbe Fähigkeit, keine Supersedes-ADR | [`ADR-0060`](0060-grpc-streaming-mechanismus.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0061` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
