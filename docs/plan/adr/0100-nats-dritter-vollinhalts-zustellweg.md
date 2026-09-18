# ADR-0100: NATS als dritter, paralleler Vollinhalts-Zustellweg für Live-Streaming

**Status:** Accepted

**Datum:** 2026-09-18

**Autor:** pt9912 (Rolleninhaber: Architect-Lauf, 2026-09-18)

**Bezug:** [`LH-FA-SST-008`](../../../spec/lastenheft.md) (Haupt-Bezug —
Live-Streaming vollständiger Change-Inhalte, protokollneutral formuliert,
Out-of-Scope nennt „gRPC, HTTP/SSE, oder beide nebeneinander" ausdrücklich
als **Beispiele**, keine abschließende Liste), [`LH-FA-SST-007`](../../../spec/lastenheft.md)
(Abgrenzung — NATS bleibt zusätzlich ein eigenständiges Wecksignal;
unverändert durch diese ADR), [`LH-FA-REA-001`](../../../spec/lastenheft.md)
(bestehender Lesezugriffsweg als Nachhol-Pfad), [`LH-FA-CON-003`](../../../spec/lastenheft.md),
[`LH-FA-CON-005`](../../../spec/lastenheft.md) (bestätigte Consumer-Position
als Fortsetzungs-Referenz), [ADR-0055](0055-nats-change-notification-wecksignal.md),
[ADR-0056](0056-nats-tabellen-granulares-subjekt.md) (NATS-Wecksignal —
**unverändert und vollständig in Kraft**, keine dieser beiden ADRs wird von
dieser Entscheidung berührt oder superseded), [ADR-0060](0060-grpc-streaming-mechanismus.md)
(Ursprung von `ChangeStreamPort`/`Broadcaster`, hier als dritter Abonnent
wiederverwendet, unverändert), [ADR-0061](0061-http-sse-zusaetzlich-zu-grpc.md)
(strukturelles Vorbild — „zweiter, paralleler Zustellweg zusätzlich" auf
genau demselben `Broadcaster`, hier auf einen dritten Weg fortgeschrieben),
[ADR-0034](0034-ports-nach-faehigkeiten.md) (Port-Zuschnitt nach Fähigkeit),
[ADR-0090](0090-beispiel-clients-volle-matrix.md), [ADR-0098](0098-beispiel-clients-start-ueber-make-dockerfile.md)
(Beispiel-Client-Matrix- und Startform-Präzedenzfälle, hier auf einen
vierten Zugriffsweg fortgeschrieben)

**Schärft:** [`SPEC-024`](../../../spec/pflichtenheft.md) (neu, diese ADR
legt den Inhalt fest), [`ARC-013`](../../../spec/architecture.md) (Rollen-
beschreibung um die zweite Fähigkeit erweitert)

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

[`ADR-0060`](0060-grpc-streaming-mechanismus.md) und
[`ADR-0061`](0061-http-sse-zusaetzlich-zu-grpc.md) haben `LH-FA-SST-008`
bereits über zwei unabhängige, parallele Zustellwege (gRPC, HTTP/SSE)
vollständig erfüllt — beide teilen sich denselben In-Prozess-Mechanismus:
`CaptureService` publiziert über den Outbound Port `ChangeStreamPort` in
einen einzigen `*grpcstream.Broadcaster` (`internal/adapters/driven/grpcstream/`),
und jeder Driving-Adapter (gRPC-Server, HTTP/SSE-Handler) abonniert ihn
über `Subscribe() (<-chan *model.Change, cancel func())`. Der Auftraggeber
(pt9912) hat am 2026-09-18 explizit den Bedarf benannt, NATS als **dritten**
Zustellweg für dieselbe Fähigkeit hinzuzufügen — nicht als Ersatz für die
beiden bestehenden Wege und nicht als Erweiterung von `LH-FA-SST-007`s
Wecksignal.

`LH-FA-SST-008`s Out-of-Scope-Klausel nennt „gRPC, HTTP/SSE, oder beide
nebeneinander" ausdrücklich als **Beispiele** für die offene
Architektur-/Spezifikationsfrage „konkretes Übertragungsprotokoll" — keine
abschließende Aufzählung. Ein dritter Weg ist damit dieselbe Art Frage wie
die bereits beantwortete zweite: `ADR-0061`s Kontext hat den kategorialen
Unterschied zwischen „welcher **einzige** Mechanismus" (`ADR-0060`
Teilfrage 1, unter der damals noch protokollspezifischen Fassung von
`LH-FA-SST-008`) und „wie kommt ein **weiterer, paralleler** Mechanismus
hinzu" bereits herausgearbeitet; dieselbe Logik trägt hier ein drittes Mal.

Diese ADR ist **keine Korrektur und keine Supersedes-ADR** — weder zu
`ADR-0060`/`ADR-0061` (deren Textkörper, Teilfragen und Verdikte bleiben
unangetastet und in Kraft; `ChangeStreamPort` und `Broadcaster` werden
unverändert wiederverwendet) noch zu `ADR-0055`/`ADR-0056` (das
NATS-Wecksignal für `LH-FA-SST-007` — Core NATS, leerer Payload, Subjekt-
Schema `cdc.changes.<source_id>.<schema>.<table>`, Fehlerklasse `transient`
nach ACK — bleibt **byte-identisch** in Kraft; kein Zeilen-Diff an
`internal/adapters/driven/natsnotify/` oder `ChangeNotificationPort` ist
Folge dieser ADR). Eine Option, die `ADR-0055` Punkt 3 (leerer Payload)
oder dessen Subjekt-Schema stillschweigend um Inhalt erweitert hätte, wäre
strukturell unzulässig auf diesem Weg — exakt die Begründung, mit der
`ADR-0060` Teilfrage 1 Option D bereits verworfen wurde (`AGENTS.md` §3.5:
eine `Accepted`-ADR wird nicht inhaltlich überschrieben, eine Korrektur
braucht `Supersedes`). Diese ADR fügt der bestehenden Entscheidungslage
etwas hinzu, sie ändert sie nicht.

Vier Konstraints prägen den Lösungsraum, drei davon bereits durch
`ADR-0055`/`ADR-0060`/`ADR-0061` vorgeprägt, einer neu:

- **`ChangeStreamPort`/`Broadcaster` sind bereits protokoll-agnostisch**
  (`ADR-0060` Teilfrage 2, von `ADR-0061` bereits ein zweites Mal bestätigt):
  `Subscribe()` liefert einen reinen `*model.Change`-Kanal, kein
  gRPC- oder HTTP-spezifischer Typ taucht darin auf. Ein dritter Abonnent
  braucht daran nichts zu ändern.
- **`LH-FA-SST-008`s Boundary- und Negative-Kriterien gelten für jeden
  Zustellweg, der die Fähigkeit für sich beansprucht — auch für einen
  dritten.** Insbesondere das Negative-Kriterium (Ablehnung eines
  Verbindungsversuchs ohne gültige Authentifizierung) ist **neu zu
  bewerten**, weil NATS in diesem Repo bislang **keine** Verbindungs-Auth
  trägt: `compose.yaml`s `nats`-Service startet ohne `--auth`/`--user`
  o. ä. (`command: ["-m", "8222"]`, ausschließlich der Monitor-Port für den
  Healthcheck), und `ADR-0055` hat das bewusst so gelassen — für ein leeres
  Wecksignal ohne Inhalt gibt es nichts zu schützen. Für einen
  **inhaltstragenden** dritten Weg gilt diese Begründung nicht mehr:
  Bliebe die NATS-Verbindung offen, entstünde eine unauthentifizierte
  Seitentür zu genau den Daten, die gRPC (`ADR-0060` Teilfrage 4) und SSE
  (`ADR-0061` Teilfrage 4) bereits hinter denselben zwei Token-Klassen
  schützen — ein Sicherheitsrückschritt, keine Detailfrage.
- **`ADR-0055`/`ADR-0056`s Subjekt-Schema `cdc.changes.<source_id>.<schema>.<table>`
  und sein leerer Payload sind `Accepted`-Inhalt, nicht verhandelbar auf
  diesem Weg** (`AGENTS.md` §3.5) — ein dritter Weg braucht einen eigenen
  Subjekt-Namensraum, keinen wiederverwendeten.
- **Bestehender Capture-kritischer Pfad ist unantastbar**
  (`ADR-0011`/`ADR-0027`) — dieselbe Grenze wie bei allen bisherigen
  Zustellwegen.

## Entscheidung

Wir wählen **NATS Core (kein JetStream) als dritten Abonnenten des
bestehenden `ChangeStreamPort`/`Broadcaster`, über einen eigenen,
opt-in-token-geschützten Subjekt-Namensraum** — mit sechs Teilfragen.

### Teilfrage 1 — Erzeugerpfad und Anknüpfung an den Broadcaster

| Option | Pro | Contra |
|---|---|---|
| A — nichts tun (nur gRPC/SSE bleiben) | kein Aufwand; `LH-FA-SST-008` ist bereits vollständig erfüllt | ein Consumer, dessen Betriebsumgebung ausschließlich NATS als Zustellmechanismus vorsieht (z. B. weil er bereits einen NATS-Broker als zentrale Integrationsschicht betreibt), müsste zusätzlich einen HTTP- oder gRPC-Client pflegen, obwohl der Auftrag einen dritten, NATS-nativen Weg ausdrücklich benennt |
| B — `ChangeNotificationPort`/`natsnotify` (`ADR-0055`) um Inhalt erweitern | ein Port statt zwei | verboten auf diesem Weg (`AGENTS.md` §3.5) — ändert `Accepted`-Inhalt von `ADR-0055` Punkt 3 ohne `Supersedes`; bricht zudem `LH-FA-SST-007`s eigene Abgrenzung („NATS bleibt eigenständiges Wecksignal") |
| **C — neuer Driven-Adapter `internal/adapters/driven/natsstream/` als dritter `Broadcaster`-Abonnent (gewählt)** | `Broadcaster.Subscribe()` liefert bereits einen protokoll-agnostischen Kanal — ein dritter Aufrufer braucht keine Änderung an `ChangeStreamPort`, `Broadcaster` oder `CaptureService`; **keine neue `CaptureService`-Konstruktions-Option nötig**, anders als bei `ADR-0060`s initialer Einführung von `WithChangeStream` — dieser Weg ist strukturell noch einfacher als `ADR-0061`s SSE-Ergänzung, weil er nicht einmal einen neuen Endpunkt im bestehenden HTTP-Adapter braucht, sondern einen reinen Hintergrund-Verteiler | ein drittes Paket, das aus dem `Broadcaster`-Kanal liest — dieselbe Art Kopplung wie gRPC/SSE, nur ohne eigenen externen Listener |
| D — eigener, vierter Outbound Port `ChangeNatsStreamPort` neben `ChangeStreamPort` | Isolation zwischen den Zustellwegen | dupliziert `ADR-0060` Teilfrage 2 Option C's bereits verworfenen Einwand: zwei parallele Fan-outs für dieselbe Fähigkeit ohne fachlichen Grund, verdoppelt Publish-Kosten und `CaptureService`-Konstruktions-Optionen |

`ChangeStreamPort`, `Broadcaster` und die bestehende `CaptureService`-Option
`WithChangeStream` bleiben durch diese Entscheidung **unverändert** — anders
als der gRPC-Server oder der SSE-Handler ist `natsstream.Publisher` kein
Driving-Adapter (er nimmt keine externe Verbindung entgegen, die ihn
auslöst) und kein neuer Inbound-Zugriffsweg; er ist ein Driven-Adapter, der
denselben lokal deklarierten `changeSubscriber`-Interface-Vertrag erfüllt
wie die Driving-Adapter aus `ADR-0060`/`ADR-0061` (`Subscribe() (<-chan
*model.Change, func())`), und den Bootstrap ein drittes Mal mit dem
konkreten `*grpcstream.Broadcaster` verdrahtet — kein Adapter-Paket
importiert ein anderes Adapter-Paket direkt, dieselbe Richtungskonvention
wie in `ADR-0060` Teilfrage 2 bereits etabliert. `.a-check.yml` braucht
**keine** Änderung — das neue Paket liegt vollständig im bestehenden Glob
`adapters: ["internal/adapters/**"]`.

### Teilfrage 2 — Subjekt- und Nachrichtenschema

| Option | Pro | Contra |
|---|---|---|
| A — dasselbe Subjekt wie das Wecksignal (`cdc.changes.<source_id>.<schema>.<table>`), jetzt mit Inhalt gefüllt | ein Subjekt-Namensraum statt zwei | verboten (`AGENTS.md` §3.5, siehe Kontext) — ändert `ADR-0055` Punkt 3/`ADR-0056` Punkt 2 stillschweigend; ein bestehender Wecksignal-Consumer, der laut `SPEC-017` einen leeren Payload erwartet, bekäme unangekündigt Inhalt |
| **B — eigener, paralleler Subjekt-Namensraum `cdc.stream.<source_id>.<schema>.<table>` (gewählt)** | keine Berührung von `ADR-0055`/`ADR-0056`s `Accepted`-Inhalt; dieselbe vierstufige, tabellen-granulare Struktur wie das Wecksignal (Wiedererkennungswert, dieselbe Wildcard-Ergonomie: `cdc.stream.<source_id>.>` für alle Tabellen einer Quelle, `cdc.stream.>` für alle Quellen) — nur der Wurzel-Token unterscheidet die beiden Fähigkeiten eindeutig | ein Consumer muss wissen, dass es zwei verschiedene NATS-Fähigkeiten mit unterschiedlichen Wurzel-Tokens gibt — Dokumentationspflicht, kein struktureller Nachteil |
| C — ein einziges, flaches Subjekt ohne Tabellen-Granularität (`cdc.stream.<source_id>`) | einfacheres Subjekt-Schema | inkonsistent zur bereits etablierten Tabellen-Granularität des Wecksignals (`ADR-0056`) und zur Nicht-Granularität von gRPC/SSE (`ADR-0060`/`ADR-0061` liefern unfiltriert alle Changes der konfigurierten Tabellen) — hier ist Granularität sogar günstiger zu haben, weil das Subjekt-Schema die Filterung bereits Broker-seitig trägt; ein Verzicht darauf wäre keine Vereinfachung, sondern ein ungenutzter Vorteil |

Für den **Nachrichteninhalt** wird das bereits für SSE definierte
JSON-Schema (`SPEC-021`) unverändert übernommen — dieselben zehn Felder wie
`model.Change` (`change_id`, `transaction_id`, `source_table_id`,
`sequence`, `operation`, `old_image`, `new_image`, `schema_version`,
`schema`, `table`), als NATS-Message-Payload (Bytes des JSON-Dokuments).
Kein drittes Nachrichtenschema für dieselben Daten — dieselbe
Wiederverwendungslogik wie `ADR-0061` Teilfrage 1, hier vom Broadcaster auf
das Nachrichtenformat übertragen. **Granularität: eine Nachricht je
Zeilen-Change**, wie bei gRPC/SSE (`ADR-0060` Teilfrage 2) — **nicht** die
Tabellen-Dedup-Kardinalität des Wecksignals (`ADR-0056` Punkt 3). Beide
NATS-Fähigkeiten haben unterschiedliche Kardinalität aus fachlichem Grund:
Das Wecksignal trägt keine Positionsinformation und braucht daher keine
höhere Auflösung als „etwas ist neu in dieser Tabelle"; der Vollinhalts-
Stream trägt jeden Change einzeln, weil ein Consumer sonst Inhalt verlöre.

### Teilfrage 3 — Zustellsemantik: Core NATS vs. JetStream

| Option | Pro | Contra |
|---|---|---|
| A — JetStream (Streams, Persistenz, Consumer-Gruppen, Replay) | eigene Zustellgarantie, Replay ohne Rückgriff auf den SQL-Pfad | baut exakt die zweite Nachvollziehbarkeits-Infrastruktur auf, die `ADR-0055` Option B und `ADR-0060` Teilfrage 3 Option B beide bereits verworfen haben — der `Broadcaster` selbst ist bewusst zustandslos (kein Puffer), und `LH-FA-SST-008`s Boundary-Kriterium verlangt kein Replay **innerhalb** eines Streams, nur Nachholbarkeit über den bestehenden Lesezugriffsweg; JetStream würde eine dritte parallele Wahrheit neben Store und den beiden bestehenden Fire-and-Forget-Wegen einführen, mit dem bekannten Konfliktauflösungsproblem |
| **B — Core NATS, Fire-and-Forget, kein Replay (gewählt)** | konsistent mit dem `Broadcaster`s bereits etablierter Zustellsemantik (`ADR-0060` Teilfrage 3, von `ADR-0061` Teilfrage 3 ein zweites Mal bestätigt) — ein Consumer sieht dasselbe Grundverhalten unabhängig vom gewählten Protokoll; kein zusätzlicher Zustand; identische Fehlerklasse-Isolation wie die beiden bestehenden Wege (ein Publish-Fehlschlag beeinflusst nichts stromaufwärts, weil der Publisher nur aus dem bereits isolierten `Broadcaster`-Kanal liest, nicht aus `CaptureService` selbst) | ein NATS-Stream-Consumer, der kurz getrennt war, verpasst Changes ersatzlos — muss wie ein gRPC-/SSE-Consumer selbst über den bestehenden Lesezugriffsweg nachholen |
| C — Core NATS mit begrenztem In-Memory-Ring-Buffer im neuen Publisher | überbrückt sehr kurze Trennungen | dieselbe bereits zweimal verworfene unvorhersagbare Zusicherung (`ADR-0060` Teilfrage 3 Option C, `ADR-0061` Teilfrage 3 Option C) — macht die drei Zustellwege für dieselbe Fähigkeit inkonsistent |

Der Unterschied zu `ADR-0055`s Core-NATS-Wahl für das Wecksignal ist rein
der Anlass: Dort folgte Core NATS aus „das Signal trägt ohnehin keinen
Zustand, den es zu bewahren gäbe". Hier folgt dieselbe Wahl aus einem
anderen, aber gleich tragfähigen Grund: Der `Broadcaster`, den dieser
dritte Weg wiederverwendet, ist bereits zustandslos und Fire-and-Forget
(`ADR-0060` Teilfrage 3) — JetStream einzuführen würde nicht die
Zustellsemantik *dieses* Wegs verbessern, sondern eine Eigenschaft
behaupten, die die gemeinsame Quelle (`Broadcaster`) strukturell nicht hat.

### Teilfrage 4 — Authentifizierung/Autorisierung

| Option | Pro | Contra |
|---|---|---|
| A — kein Schutz (nichts tun) | kein Zusatzaufwand | `LH-FA-SST-008`s Negative-Kriterium wird für diesen Weg unerfüllbar; öffnet eine unauthentifizierte Seitentür zu Daten, die gRPC/SSE bereits hinter Tokens schützen — realer Sicherheitsrückschritt, kein akzeptabler Blindfleck |
| B — NATS-Accounts/Permissions-Konfiguration (`nats-server.conf`, subjekt-scoped: ein Nutzer ohne Credential darf nur `cdc.changes.>` (Wecksignal), ein credentialierter Nutzer zusätzlich `cdc.stream.>`) | granularste Lösung — Wecksignal bliebe exakt so offen wie heute, nur der Vollinhalts-Namensraum wäre geschützt | deutlich höherer Konfigurationsaufwand (mounted Config-Datei statt einer CLI-Flag-Zeile, Accounts-/Permissions-Modell, das dieses Repo bislang nirgends führt) — unverhältnismäßig gegenüber dem Bedarf dieser ADR; bewusst als Re-Evaluierungs-Trigger vertagt, nicht verworfen |
| **C — einzelner, geteilter NATS-Verbindungs-Token, der die **gesamte** Aktivierung des Vollinhalts-Wegs voraussetzt (gewählt)** | nutzt Core NATS' einfachsten nativen Mechanismus (ein Server-CLI-Flag/eine Config-Zeile, kein Accounts-Modell); an **eine** neue, rein additive Umgebungsvariable gekoppelt (`CDC_NATS_STREAM_TOKEN`), die die Aktivierung des dritten Wegs *und* die Server-Auth-Pflicht in einem Schritt trägt — kein bestehendes, unverändert betriebenes Deployment ist betroffen, weil niemand diese Variable heute setzt | sobald ein Betreiber den dritten Weg aktiviert, verlangt der NATS-Server (Core NATS' einzige einfache Form) den Token **serverweit** — auch die bestehende, bislang anonyme Wecksignal-Verbindung muss ihn dann mitführen; ein bewusst benannter, aber nur beim expliziten Opt-in wirksamer Nebeneffekt (siehe Teilfrage 5 und Konsequenzen) |

Ein Verbindungsversuch auf den `cdc.stream.>`-Namensraum ohne den
konfigurierten Token wird vom NATS-Server selbst abgelehnt (Verbindungsebene,
nicht Anwendungsebene) — sichtbar als Verbindungsfehler, nicht als
stillschweigend leerer Stream. Das ist strukturell dieselbe Aussage wie
gRPCs `Unauthenticated` (`ADR-0060` Teilfrage 4) und SSEs `401`
(`ADR-0061` Teilfrage 4), nur eine Ebene tiefer verortet (Transport-
Verbindung statt Anwendungs-Handshake), weil Core NATS keine
Anwendungsebene zwischen Verbindungsaufbau und Subjekt-Zustellung kennt.

Option B (Accounts/Permissions) bleibt die architektonisch sauberere
Lösung — sie würde das Wecksignal beim Opt-in unangetastet anonym lassen —,
ist aber ein eigener, größerer Konfigurationsvorgang ohne Präzedenzfall in
diesem Repo und wird bewusst nicht mit dieser ADR miterledigt (siehe
Re-Evaluierungs-Trigger), analog zu `AGENTS.md` §3.14s Umgang mit einer
erkannten, aber vertagten saubereren Lösung.

### Teilfrage 5 — Verbindungs-Wiederverwendung und Aktivierung

| Option | Pro | Contra |
|---|---|---|
| A — zweite, eigenständige NATS-Verbindung (eigene `nats.Conn`, eigene URL-Variable) | Isolation von der Wecksignal-Verbindung | zwei Verbindungen zum selben Server für zwei Fähigkeiten desselben Prozesses, ohne fachlichen Grund; bei Option C in Teilfrage 4 ohnehin hinfällig, weil ein serverweiter Token beide Verbindungen gleichermaßen beträfe |
| **B — dieselbe `nats.Conn` wie das Wecksignal, Aktivierung über zwei Bedingungen (gewählt)** | eine Verbindung, ein Fehlerpfad bei Verbindungsproblemen (`ErrConfiguration`, analog zu `ADR-0055`s bestehendem Muster in `internal/bootstrap/wiring.go`); folgt `ADR-0061` Teilfrage 5s bereits akzeptiertem Muster „ein zusätzlicher Bezieher auf einer bereits geteilten Ressource ist ein additiver, bewusst benannter Nebeneffekt, kein Bruch" | `CDC_NATS_URL` allein aktiviert weiterhin nur das Wecksignal — der dritte Weg braucht **zusätzlich** `CDC_NATS_STREAM_TOKEN`, sonst bliebe der bestehende, unveränderte Betrieb (`CDC_NATS_URL` gesetzt, `CDC_NATS_STREAM_TOKEN` ungesetzt) plötzlich implizit spezifikationswidrig, sobald ein Betreiber irrtümlich einen unauthentifizierten dritten Weg erwartet — dagegen schützt exakt die Zwei-Bedingungen-Form |
| C — eigenes Feature-Gate zusätzlich zu `CDC_NATS_URL`, ohne dass der Token gleichzeitig die Server-Auth-Pflicht trägt | trennt Aktivierung von Auth-Konfiguration | löst Teilfrage 4 nicht — bräuchte eine dritte Variable, ohne einen Mehrwert gegenüber der bereits gewählten Kopplung „ein Token, zwei Zwecke" |

**Bindende Festlegung:** Der `natsstream.Publisher` (und der `Broadcaster`,
falls er sonst nicht bereits über `CDC_GRPC_ADDR`/`CDC_HTTP_ADDR`
konstruiert würde — `changeStreamEnabled` aus `ADR-0061` Teilfrage 5 wird
um eine dritte Oder-Bedingung erweitert: `cfg.NatsURL != "" &&
cfg.NatsStreamToken != ""`) werden **nur** konstruiert, wenn **sowohl**
`CDC_NATS_URL` **als auch** `CDC_NATS_STREAM_TOKEN` gesetzt sind. Ist nur
`CDC_NATS_URL` gesetzt (heutiger Zustand jeder bestehenden Installation):
bestehendes Verhalten bleibt bit-identisch — Wecksignal aktiv, kein dritter
Weg, keine Server-Auth-Pflicht. Ist `CDC_NATS_STREAM_TOKEN` gesetzt, aber
`CDC_NATS_URL` leer: Konfigurationsfehler beim Start (`ErrConfiguration`,
analog zum bestehenden Umgang mit einer fehlgeschlagenen
`CDC_NATS_URL`-Verbindung) — ein Betreiber, der den dritten Weg aktiviert,
aber kein Verbindungsziel angibt, soll das beim Start bemerken. Die
bestehende `nats.Connect(cfg.NatsURL)`-Aufrufstelle in
`internal/bootstrap/wiring.go` erhält bei gesetztem
`CDC_NATS_STREAM_TOKEN` eine zusätzliche Client-Option
(`nats.Token(cfg.NatsStreamToken)` oder gleichwertig) — dieselbe
Verbindung, nicht zwei.

### Teilfrage 6 — Beispiel-Client-Matrix und Namenskollision

`SURFACE=nats` ist bereits vergeben — `make example-run-go/-csharp/-kotlin
SURFACE=nats` startet den bestehenden Wecksignal-Client
(`examples/nats-client`, `examples/csharp/nats-client`,
`examples/kotlin/nats-client`, [`ADR-0079`](0079-nats-beispielclient-vierter-examples-client.md),
[`ADR-0098`](0098-beispiel-clients-start-ueber-make-dockerfile.md)). Ein
vierter Zugriffsweg braucht einen eigenen, kollisionsfreien `SURFACE`-Wert.

| Option | Pro | Contra |
|---|---|---|
| A — `SURFACE=nats2` oder `SURFACE=nats-full` | kurz | trägt keine erkennbare Bedeutung; „2" suggeriert Versionierung statt einer zweiten Fähigkeit |
| **B — `SURFACE=nats-stream` (gewählt)**, Client-Verzeichnisse `examples/nats-stream-client`, `examples/csharp/nats-stream-client`, `examples/kotlin/nats-stream-client` | benennt exakt den fachlichen Unterschied zum bestehenden `nats`-Wert (Wecksignal vs. Vollinhalts-Stream); folgt `SPEC-023`s bestehender Namenskonvention „die Namen tragen `-client`"; keine Kollision mit `nats`, `http`, `sse`, `grpc` | ein weiterer, vierter Eintrag in der `SURFACE=`-Fehlermeldung von `harness/mk/examples.mk` |
| C — den bestehenden `nats-client` um den Vollinhalts-Modus erweitern (ein Programm, zwei Modi über ein Flag) | kein neues Client-Verzeichnis | widerspricht `SPEC-023`s Klasse „ein Beispiel-Client spricht **eine** dokumentierte Zugriffs-Oberfläche an" — Wecksignal und Vollinhalts-Stream sind zwei verschiedene Oberflächen mit unterschiedlichem Nachrichtenschema und Subjekt-Namensraum, kein Modus desselben Programms |

Nach `SPEC-023`s bereits etablierter Festlegung „Sprachen und Umfang: volle
Matrix" ([`ADR-0090`](0090-beispiel-clients-volle-matrix.md)) bekommt
dieser vierte Zugriffsweg von Anfang an **alle drei Sprachen** (Go, C#,
Kotlin) — kein Go-only-Pilot. Startform folgt
[`ADR-0098`](0098-beispiel-clients-start-ueber-make-dockerfile.md)
unverändert: `make example-run-go/-csharp/-kotlin SURFACE=nats-stream`.
`harness/mk/examples.mk`s drei `$(error …)`-Zeilen (`SURFACE muss http, sse,
grpc oder nats sein`) werden um `nats-stream` als vierten zulässigen Wert
erweitert — Folgepflicht des umsetzenden Slices, hier bereits benannt,
damit der Planner sie nicht neu herleiten muss.

## Konsequenzen

- Positiv: `LH-FA-SST-008` bekommt einen dritten, unabhängig nutzbaren
  Zustellweg, ohne `ADR-0060`/`ADR-0061`s Entscheidungen, den
  `ChangeStreamPort`-Vertrag oder den `Broadcaster` zu ändern — alle drei
  Wege teilen sich denselben Erzeuger-Pfad (`ADR-0034`).
- Positiv: Kein Zeilen-Diff an `CaptureService` nötig — anders als
  `ADR-0060`s initiale Einführung von `WithChangeStream` ist dieser dritte
  Weg ein reiner zusätzlicher `Broadcaster`-Abonnent; der bereits isolierte
  Kanal trägt die Fehlerisolation automatisch mit.
- Positiv: `ADR-0055`/`ADR-0056`s Wecksignal-Entscheidung bleibt vollständig
  unberührt und byte-identisch in Kraft — kein Zeilen-Diff an
  `natsnotify`/`ChangeNotificationPort`.
- Positiv: Die neue Auth-Kopplung (Teilfrage 4/5) schließt den einzigen
  realen Sicherheits-Unterschied zu gRPC/SSE, ohne die schwerere
  Accounts/Permissions-Lösung vorwegzunehmen, die kein aktueller Bedarf
  rechtfertigt.
- Negativ: Sobald ein Betreiber `CDC_NATS_STREAM_TOKEN` setzt, verlangt der
  NATS-Server serverweit einen Token — auch für die bislang anonyme
  Wecksignal-Verbindung. Das ist ein bewusst benannter Konsequenz-Punkt
  dieses **freiwilligen** Opt-ins (kein bestehendes Deployment ist
  betroffen, das die Variable nicht setzt), aber ein realer operativer
  Unterschied gegenüber „nur den dritten Weg einschalten" — dokumentiert,
  nicht versteckt.
- Negativ: Core NATS ohne Replay bedeutet, dass auch der NATS-Stream-Weg
  wie gRPC/SSE keine Zustellgarantie trägt — dieselbe bewusste Entscheidung
  wie bei den beiden bestehenden Wegen, zum dritten Mal bestätigt.
- Negativ: Drei parallele Zustellwege für dieselbe Fähigkeit bedeuten jetzt
  drei Testpfade in `make test-integration` und drei Beispiel-Client-
  Familien (Go/C#/Kotlin je Weg) zu pflegen — derselbe Preis, den `ADR-0061`
  bereits für den Sprung von eins auf zwei benannt hat, hier für den Sprung
  auf drei.
- Folgepflicht: Der umsetzende Slice-Schnitt implementiert
  `internal/adapters/driven/natsstream/` (`Publisher`), die Bootstrap-
  Verdrahtung (`CDC_NATS_STREAM_TOKEN`, die Zwei-Bedingungen-Aktivierung aus
  Teilfrage 5, die Client-Options-Erweiterung an der bestehenden
  `nats.Connect`-Aufrufstelle), die Erweiterung von `changeStreamEnabled`
  um die dritte Oder-Bedingung, und die `compose.yaml`-Erweiterung des
  `nats`-Service um eine Server-Auth-Konfiguration für den Integrationstest
  (ein fester Test-Token, damit sowohl ein Erfolgs- als auch ein
  Ablehnungs-Beleg real geführt werden kann).
- Folgepflicht: `harness/mk/examples.mk`s drei `SURFACE=`-Prüfungen werden
  um `nats-stream` erweitert (Teilfrage 6); `examples/nats-stream-client`,
  `examples/csharp/nats-stream-client`, `examples/kotlin/nats-stream-client`
  entstehen nach demselben Muster wie die bestehenden vier Zugriffswege.
  `examples/README.md` bekommt einen neuen Zugriffs-Abschnitt
  (`SPEC-023`s Handbuch-Bindung).
- Folgepflicht: `spec/pflichtenheft.md` erhält `SPEC-024` (Subjekt-Schema,
  Nachrichtenform, Auth-Kopplung, Aktivierung) — **bereits Teil dieser
  ADR**, siehe unten — sowie eine `§6 Externe Verträge`-Zeile und eine
  Historie-Zeile; die `SPEC-016`/`ADR-0091`/`ADR-0092`-Zugangsdaten-Klasse
  (Feldmengen-Paarung Konfigurationsdatei/Env) prüft, ob
  `CDC_NATS_STREAM_TOKEN` als env-exklusiver, zugangsdaten-tragender
  Schlüssel in ihre bestehende Liste aufzunehmen ist — Gegenstand des
  umsetzenden Slices, nicht dieser ADR, weil diese Feldmengen-Paarungs-
  Maschinerie eine eigene, bereits etablierte Zuständigkeit ist
  (`ADR-0089`: „kein Sensor, Wächter Review").
- Folgepflicht: `spec/architecture.md` `ARC-013` ist mit dieser ADR bereits
  aktualisiert (siehe unten) — analog zur Aufgabenteilung aus `ADR-0055`.
- Folgepflicht: Erweiterung um NATS-Accounts/Permissions-Konfiguration
  (Teilfrage 4 Option B) braucht eine eigene Folge-ADR — nicht Gegenstand
  dieser ADR (siehe Re-Evaluierungs-Trigger).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| Go-Unit-Test (`internal/adapters/driven/natsstream`, Whitebox) | Ein über `Broadcaster.Publish` verteilter `Change` erreicht einen abonnierten NATS-Test-Client als vollständiges JSON-Event auf `cdc.stream.<source_id>.<schema>.<table>`; ein `Publish` ohne verbundenen NATS-Server/Publisher blockiert nicht (Fire-and-Forget-Regressionstest, spiegelt die bestehenden gRPC-/SSE-Tests) | `make test` |
| Go-Unit-Test (`internal/adapters/driven/natsstream`, Whitebox) | Das publizierte Subjekt trägt den Wurzel-Token `cdc.stream`, niemals `cdc.changes` — Regressionstest gegen eine versehentliche Kollision mit dem Wecksignal-Namensraum aus `ADR-0056` | `make test` |
| Go-Unit-Test (`internal/bootstrap`, Whitebox) | Der `natsstream.Publisher` wird genau dann konstruiert, wenn `CDC_NATS_URL` **und** `CDC_NATS_STREAM_TOKEN` gesetzt sind; ist nur `CDC_NATS_URL` gesetzt, bleibt der bestehende Wecksignal-Pfad unverändert aktiv, ohne den neuen Publisher — Regressionstest gegen eine versehentliche Ein-Bedingungs-Aktivierung | `make test` |
| `.a-check` | unverändert — das neue Paket liegt vollständig im bestehenden Glob `adapters: ["internal/adapters/**"]`, kein neuer Hexagon-Schichten-Edge nötig | `make a-check` |
| `make test-integration` | ein NATS-Client, der mit gültigem `CDC_NATS_STREAM_TOKEN` verbindet und `cdc.stream.<source_id>.<schema>.<table>` abonniert, empfängt eine danach committete Änderung mit vollständigem Inhalt; ein Verbindungsversuch ohne oder mit falschem Token wird vom NATS-Server abgelehnt; das bestehende Wecksignal (`natssub`) funktioniert unverändert weiter, sofern derselbe Test-Token auch dort mitgeführt wird — Umsetzung Gegenstand des umsetzenden Slice-Schnitts | `make test-integration` |
| `make examples-csharp` / `make examples-kotlin` | die neuen `nats-stream-client`-Programme bauen und testen netzlos wie die bestehenden vier Zugriffswege | `make examples-csharp`, `make examples-kotlin` |

## Slice-Schnitt-Empfehlung

Regeln dieser Sektion: Empfehlung, endgültiger Schnitt liegt bei der
Planner-Rolle (Baseline-Regelwerk `modul-05-planning-harness.md`).

1. **Slice A** — Kernfähigkeit: `internal/adapters/driven/natsstream/`
   (`Publisher`, JSON-Serialisierung nach `SPEC-021`s Schema, Subjekt-
   Aufbau), Bootstrap-Verdrahtung (`CDC_NATS_STREAM_TOKEN`,
   Zwei-Bedingungen-Aktivierung, Client-Options-Erweiterung),
   `changeStreamEnabled`-Erweiterung, Unit-Tests, `compose.yaml`-
   NATS-Auth-Konfiguration für den Integrationstest, ein Wegwerf-
   Belegträger (`tools/harness/`, analog zu `tools/harness/natssub/`) und
   die zugehörige `make test-integration`-Erweiterung. Diese Slice braucht
   `ADR-0060`s Slice A **und** B (Broadcaster muss existieren und beliefert
   werden), aber **nicht** `ADR-0060`s Slice C oder `ADR-0061`s Slice —
   der dritte Weg ist von den beiden bestehenden Beispiel-Clients
   unabhängig.
2. **Slice B** — Beispiel-Client-Matrix: `examples/nats-stream-client`
   (Go), `examples/csharp/nats-stream-client`, `examples/kotlin/nats-stream-client`,
   `harness/mk/examples.mk`-Erweiterung um `SURFACE=nats-stream`,
   `examples/README.md`-Zugriffs-Abschnitt — braucht Slice A als
   Voraussetzung (ohne einen echten Publisher gäbe es nichts zu
   demonstrieren).

## Re-Evaluierungs-Trigger

Wird ein Bedarf an echter Subjekt-scoped NATS-Zugriffskontrolle formuliert
— insbesondere, dass das Wecksignal (`cdc.changes.>`) auch nach Aktivierung
des Vollinhalts-Wegs anonym erreichbar bleiben soll, statt den geteilten
Serverweiten Token mitzuführen (Teilfrage 4 Option B, Accounts/Permissions-
Konfiguration) —, oder ein Bedarf an Stream-internem Replay bzw.
JetStream-Zustellgarantien für diesen dritten Weg: Beide widersprechen der
hier getroffenen Wahl fundamental genug, dass sie eine eigene Folge-ADR
brauchen, die die jeweilige Option neu bewertet — nicht als Korrektur
dieser ADR, sondern als eigene Entscheidung, weil sich die Prämisse (kein
belegter Bedarf über die drei Akzeptanzkriterien von `LH-FA-SST-008`
hinaus) geändert hätte. Sonst permanent — die Wahl Core NATS als dritter,
paralleler Zustellweg mit geteiltem Verbindungs-Token gilt unabhängig vom
Zeitpunkt der Umsetzung.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-18 | Accepted — Architect-Entscheidung nach explizitem Auftraggeber-Bedarf (2026-09-18); ergänzt `ADR-0060`/`ADR-0061` um einen dritten, parallelen Zustellweg für dieselbe Fähigkeit, keine Supersedes-ADR; berührt `ADR-0055`/`ADR-0056` (NATS-Wecksignal) nicht inhaltlich | [`ADR-0060`](0060-grpc-streaming-mechanismus.md), [`ADR-0061`](0061-http-sse-zusaetzlich-zu-grpc.md) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0100` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
