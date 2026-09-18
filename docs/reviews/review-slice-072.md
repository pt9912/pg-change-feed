# Review-Report: slice-072 — 2026-09-15

**Review-Art:** Code — geprüft gegen Plan + Entscheidungen (Baseline-Regelwerk
`v6.5.0` · `regelwerk/modul-10-review-harness.md` §Drei Review-Arten);
DoD-/Spec-Konformität ist Verifier-Aufgabe und **nicht** Gegenstand dieses
Reports.

**Gegenstand:** `slice-072`, Commit `89d31d1` (= `HEAD`). Der Elter-Commit
`f4468fc` ist der reine `next → in-progress`-Move; `git diff --name-only
f4468fc..89d31d1` nennt genau die elf erwarteten Pfade —
`docs/plan/planning/in-progress/slice-072-http-sse-streaming-endpunkt.md`,
`harness/README.md`, `harness/image-hash.txt`,
`internal/adapters/driving/http/server.go`,
`internal/adapters/driving/http/sse.go`,
`internal/adapters/driving/http/sse_test.go`,
`internal/bootstrap/changestream_internal_test.go`,
`internal/bootstrap/wiring.go`, `spec/pflichtenheft.md`,
`tools/harness/run-integration-tests.sh`, `tools/harness/sseclient/main.go`.
Kein Fremd-Commit im Range.

**Skill:** `.harness/skills/reviewer.md` @ `68d2ebd` (letzte Schärfung
2026-09-14 — die Regel `Neue Betreiber-Oberfläche ohne Handbuch-Zug` liegt
**vor** diesem Implementer-Lauf).
**Modell:** deepseek-v4.1-flash:cloud[1m] · **Datum:** 2026-09-15.

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan
  `docs/plan/planning/in-progress/slice-072-http-sse-streaming-endpunkt.md`
  vollständig (§1–§8), inklusive des §3-Nachzugs aus `89d31d1`
- [`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md) vollständig
  — insbesondere Teilfrage 1 (gemeinsame Quelle), Teilfrage 2
  (Adapter-Platzierung, lokal deklarierte `changeSubscriber`-Schnittstelle),
  Teilfrage 3 (Fire-and-Forget, `Last-Event-ID` ungenutzt), Teilfrage 4
  (Authn), Teilfrage 5 (Bootstrap-Aktivierung, `503`-Folgepflicht),
  Teilfrage 6, §Konsequenzen (Folgepflichten), §Fitness Function
- [`ADR-0066`](../plan/adr/0066-broadcaster-begrenzte-empfangswarteschlange.md)
  vollständig — §Entscheidung Festlegungen 1–4, §Verhältnis zu `ADR-0061`
- [`ADR-0060`](../plan/adr/0060-grpc-streaming-mechanismus.md)
  §Teilfrage 2/3/6 (Port-/Adapter-Design, Fire-and-Forget, Bootstrap),
  [`ADR-0068`](../plan/adr/0068-wegwerf-clients-begrenzte-import-berechtigung.md)
  §Entscheidung (Schicht `tooling`, Kante `tooling → adapters`),
  [`ADR-0041`](../plan/adr/0041-a-check-maschinenform-architekturpruefung.md),
  [`ADR-0044`](../plan/adr/0044-image-beleg-semantik.md),
  [`ADR-0057`](../plan/adr/0057-http-grpc-api.md) §Teilfrage 3/4
- [`SPEC-018`](../../spec/pflichtenheft.md) vollständig (Authn-Header-Form,
  Fehler-Antwortform, Endpunkt-Tabelle, Streaming-Erweiterung, Historie),
  [`SPEC-020`](../../spec/pflichtenheft.md) (Abgrenzung — gRPC-Weg),
  [`LH-FA-SST-008`](../../spec/lastenheft.md) (Happy Path/Boundary/Negative),
  [`LH-FA-CAP-008`](../../spec/lastenheft.md),
  [`LH-FA-REA-001`](../../spec/lastenheft.md),
  [`LH-FA-CON-003`](../../spec/lastenheft.md),
  [`LH-FA-CON-005`](../../spec/lastenheft.md)
- `AGENTS.md` §3.1/§3.3/§3.5/§3.7/§3.9/§5/§6; `harness/conventions.md`
  (`MR-000`); `.a-check.yml`; `harness/sensors/docs-check.md`;
  `harness/sensors/a-check.md`
- Vorgänger: `docs/plan/planning/done/slice-069-grpc-streaming-adapter-grundgeruest.md`
  (lokal deklariertes `changeSubscriber`), `…/done/slice-070-grpc-capture-integration.md`,
  `…/done/slice-071-grpc-beispielclient-e2e.md`, das Review zu `slice-071`
  (F-2 Assertions-Bindung, F-4 Aufschub-Adresse)
- `docs/plan/planning/open/slice-077-handbuch-betreiber-stand.md` (§1 Ziel,
  §1 Abgrenzung, §2 DoD — die benannte Aufschub-Adresse),
  `docs/user/benutzerhandbuch.md` (Version 1.12, kein Netzwerk-Zugriffsweg
  beschrieben)
- Beobachtungs-Register `BEO-PGC`: `adapter-unittest-verdeckt-bootstrap-luecke`
  (1×, offen), `handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
  (3×, verkörpert seit slice-077), `report-nackte-id-ohne-link` (5×,
  verkörpert), `slice-chronik-in-code-kommentar` (verkörpert),
  `adapter-fehler-ausgang` (2×), `rollen-test-abdeckungsluecken` (2×)
- Code: `internal/bootstrap/wiring.go`, `internal/adapters/driving/http/{server,middleware,sse,sse_test}.go`,
  `internal/bootstrap/changestream_internal_test.go`,
  `internal/adapters/driving/grpc/server.go`,
  `internal/adapters/driven/grpcstream/{broadcaster,broadcaster_test}.go`,
  `internal/adapters/driving/replication/mapper/mapper.go` (`rowImage`),
  `internal/domain/model/change.go`, `tools/harness/run-integration-tests.sh`,
  `tools/harness/sseclient/main.go`

---

## Findings

### F-1 — Die benannte Aufschub-Adresse `slice-077` zählt zwei Netzwerk-Zugriffswege und nennt den SSE-Endpunkt nicht

- `kategorie`: LOW
- `quelle`: Baseline-Regelwerk `v6.5.0` ·
  `regelwerk/modul-05-planning-harness.md` §Ziel-Form: Slice
  (Out-of-Scope-Klasse 1: „Die Adresse muss die Sendung annehmen") ·
  `.harness/skills/reviewer.md` §Klassifikation (HIGH „Neue Betreiber-Oberfläche
  ohne Handbuch-Zug" — hier **nicht** erfüllt, s. u.) ·
  [`LH-FA-SST-008`](../../spec/lastenheft.md)
- `pfad`:
  `docs/plan/planning/in-progress/slice-072-http-sse-streaming-endpunkt.md:143`;
  `docs/plan/planning/open/slice-077-handbuch-betreiber-stand.md:44`, `:89-92`;
  `docs/plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md:162`;
  `docs/user/benutzerhandbuch.md` (Version 1.12)
- `befund`: Der Plan verweist die Betreiber-Dokumentation des neuen Endpunkts
  als benannten Aufschub auf `slice-077`; dessen §2-DoD fordert aber „die
  **zwei** Netzwerk-Zugriffswege (HTTP/JSON-API, gRPC-Stream)" in §4 — und
  [`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md) zählt den
  SSE-Stream an seiner Authn-Stelle selbst als **dritten**, eigenen Weg
  („HTTP-Verwaltungs-API `ADR-0057`, gRPC-Stream `ADR-0060`, jetzt
  SSE-Stream"). Der Empfänger ist damit strukturell richtig gewählt, seine
  Formulierung nimmt die Sendung aber nicht ausdrücklich an; ob „HTTP/JSON-API"
  den `text/event-stream`-Endpunkt mit seinen eigenen Merkmalen
  (Zustellsemantik, kein Replay) einschließt, bleibt der Lesart überlassen.
  Die Klasse „Neue Betreiber-Oberfläche ohne Handbuch-Zug" ist dennoch **nicht**
  ausgelöst: der Diff führt kein neues ENV-Feld ein (der Endpunkt hängt an der
  bereits dokumentierten `CDC_HTTP_ADDR`), und ein benannter Aufschub mit
  Adresse liegt vor.
- `verifizierbar`: ja — `slice-077` §2 gegen `slice-072` §3 lesen;
  `grep -n "Stream\|HTTP" docs/user/benutzerhandbuch.md` zeigt, dass kein
  Netzwerk-Zugriffsweg dort bislang beschrieben ist.
- `klasse`: „Aufschub-Adresse nimmt die Sendung nicht an" (2× — erstes
  Auftreten Review zu `slice-071` F-4)

### F-2 — Der `503`-Pfad des SSE-Endpunkts ist über `Run` strukturell unerreichbar; DoD und `SPEC-018` beschreiben ihn als Bootstrap-Zustand, der Unit-Test erreicht ihn nur über die direkte `Config`-Konstruktion

- `kategorie`: LOW
- `quelle`: [`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md)
  Teilfrage 5 (Folgepflicht: „Ist `CDC_HTTP_ADDR` gesetzt, aber der
  `Broadcaster` aus einem anderen Grund nicht verdrahtet, antwortet der
  Endpunkt mit `503`") gegen `internal/bootstrap/wiring.go` ·
  Beobachtungs-Register `BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke`
- `pfad`: `internal/bootstrap/wiring.go:376-378` (`changeStreamEnabled`),
  `:564-569` (Broadcaster-Konstruktion), `:674`, `:676`, `:695`;
  `internal/adapters/driving/http/sse.go:89-92`;
  `internal/adapters/driving/http/sse_test.go:276-291`;
  `docs/plan/planning/in-progress/slice-072-http-sse-streaming-endpunkt.md:99`
- `befund`: `Config.Subscriber` wird in `Run` genau einmal gesetzt, innerhalb
  von `if cfg.HTTPAddr != ""` (`:676`, `:695`); der Wert ist
  `changeBroadcaster`, und der entsteht nach `:566` genau dann, wenn
  `changeStreamEnabled(cfg.GRPCAddr, cfg.HTTPAddr)` gilt — also **immer**,
  wenn `cfg.HTTPAddr != ""`. Über den regulären Bootstrap-Pfad ist
  `Subscriber` damit nie `nil`; den `503`-Zweig erreicht nur, wer den Adapter
  direkt konstruiert (wie `TestStreamOhneBroadcasterAntwortetMit503`). Der
  Guard ist damit **kein toter Code** — er ist eine Grenz-Bedingung an einem
  exportierten `Config`-Feld, die einen Nil-Panic in eine sichtbare Antwort
  überführt, und `ADR-0061` Teilfrage 5 verlangt sie ausdrücklich —, aber die
  von der ADR und der DoD-Zeile genannte Auslöse-Bedingung („`CDC_HTTP_ADDR`
  gesetzt, Broadcaster aus einem anderen Grund nicht verdrahtet") hat in der
  heutigen Verdrahtung keinen Erzeuger; kein Lauf über `Run` kann sie
  herstellen. Dasselbe gilt für den zweiten Guard des Handlers
  (`sse.go:93-96`, `http.Flusher` nicht vorhanden → `500`).
- `verifizierbar`: ja — die zwei Verdrahtungsstellen sind zeichengleich lesbar;
  ein Bootstrap-Lauf-Test bräuchte eine reale Strecke (`make test-integration`
  belegt diesen Zustand nicht).
- `klasse`: „Adapter-Unit-Test konstruiert einen Zustand, den der Bootstrap
  nicht herstellt" (2× — `BEO-PGC/adapter-unittest-verdeckt-bootstrap-luecke`,
  erstes Auftreten Review zu `slice-061` F-1)

### F-3 — `TestStreamOhneVerbundenenClientBlockiertNicht` kann in seiner zweiten Hälfte nicht fehlschlagen

- `kategorie`: LOW
- `quelle`: [`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md)
  §Fitness Function („ein `Publish` ohne verbundenen SSE-Client blockiert
  nicht") · Maintainability · `AGENTS.md` §3.7 (ein Kommentar beschreibt, was
  da ist)
- `pfad`: `internal/adapters/driving/http/sse_test.go:293-322`
- `befund`: Der Test konstruiert den Adapter, belegt real, dass ohne laufenden
  Request **keine** Subskription entsteht (`:296-301`, eine echte
  Adapter-Aussage), und stellt danach den nicht-blockierenden Handoff lokal
  nach: `select { case subscriber.changes <- &change: default: }` über den
  Kanal des eigenen Fakes mit Kapazität 1 und leerem Puffer (`:304-315`).
  Diese Hälfte kann für keine Änderung am Produktionscode fehlschlagen — sie
  prüft `select`/`default` der Laufzeit, nicht den Adapter; Name und Kommentar
  („die Fire-and-Forget-Hälfte der Fitness Function aus `ADR-0061`", „Der
  Handoff des Broadcaster") beanspruchen mehr als die Bindung.
- `verifizierbar`: ja — der Testkörper enthält keinen Aufruf in den
  Produktionspfad; ein Entfernen jeder Adapter-Zeile außer `streamChangesHandler`
  lässt ihn grün.
- `klasse`: „Test-Zusage ohne Bindung an den Produktionspfad" (1×)

### F-4 — `rowImage`s Leer-Verzweigung ist durch keinen Test gebunden; ihre Entfernung lässt den `null`-Boundary-Test grün

- `kategorie`: LOW
- `quelle`: [`LH-FA-CAP-008`](../../spec/lastenheft.md) Boundary ·
  Maintainability
- `pfad`: `internal/adapters/driving/http/sse.go:48-53`;
  `internal/adapters/driving/http/sse_test.go:240-274`;
  `internal/adapters/driving/replication/mapper/mapper.go:524-527`
- `befund`: `rowImage` ersetzt ein leeres Bild durch `json.RawMessage("null")`;
  `TestStreamLeeresRowImageTraegtNull` reicht `nil` durch. Entfernt man die
  Verzweigung, bleibt die gesamte SSE-Testmenge grün (gemessen, s.
  §Mutations- und Probenläufe P1): `encoding/json` macht aus einem `nil`
  `RawMessage` selbst `null`. Der Produzent liefert für ein fehlendes Bild
  ohnehin `nil` (`replication/mapper/mapper.go:525-527`) und sonst einen mit
  `{` beginnenden Puffer — der Fall, den die Verzweigung wirklich trägt
  (leer, aber nicht `nil`, wäre für `json.Marshal` ungültig), wird von keinem
  Test erzeugt.
- `verifizierbar`: ja — P1 unten ist der Beleg.
- `klasse`: „Test-Zusage ohne Bindung an den Produktionspfad" (1×)

### F-5 — Der generische Fehler-Antwortform-Satz von `SPEC-018` nimmt den neuen `503` nicht auf; die Historien-Zeile trägt einen Tippfehler

- `kategorie`: LOW
- `quelle`: [`SPEC-018`](../../spec/pflichtenheft.md) (Fehler-Antwortform) ·
  [`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md) Teilfrage 5
- `pfad`: `spec/pflichtenheft.md:290-296` gegen `:326`; `spec/pflichtenheft.md:476`
- `befund`: Der Satz „**Fehler-Antwortform** (für alle Endpunkte gleich):
  `{"error": "<Klartext>"}` bei `400`/`401`/`403`/`404`/`500`" ist mit der
  Erweiterung unvollständig geworden: der neue Endpunkt antwortet über
  `writeError` auch mit `503`, und die `Aktivierung`-Zeile derselben Sektion
  (`:326`) deklariert genau das. Die Information ist damit im Abschnitt
  vorhanden — die generische Aufzählung, die „für alle Endpunkte gleich" gilt,
  nennt den Code aber nicht. Zusätzlich trägt die Historien-Zeile `:476` das
  Wort „Streamings-Endpunkt".
- `verifizierbar`: ja — beide Stellen sind zeichengleich lesbar; `grep -n "503"
  spec/pflichtenheft.md` findet `:326` und `:476`, aber nicht `:291`.
- `klasse`: „`SPEC-018`-Textstellen der Erweiterung nicht nachgezogen" (1×)

### F-6 — `ADR-0061`s Folgepflicht nennt einen **neuen** `SPEC-*`-Eintrag; der Diff folgt dem Plan und erweitert `SPEC-018` — die Abweichung ist nirgends benannt

- `kategorie`: LOW
- `quelle`: [`ADR-0061`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md)
  §Konsequenzen, Folgepflicht (`:237-240`) · Baseline-Regelwerk `v6.5.0` ·
  `regelwerk/modul-08-agentenrollen.md` §Konflikt-Pfad (Verdikt 1: „ADR gilt,
  Slice-Plan hat falsch behauptet" → Plan-Korrektur)
- `pfad`: `docs/plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md:237-240`
  gegen `spec/pflichtenheft.md:310-326`;
  `docs/plan/planning/in-progress/slice-072-http-sse-streaming-endpunkt.md:102-104`
- `befund`: Die ADR verlangt „einen **neuen `SPEC-*`-Eintrag** für das
  konkrete SSE-Nachrichtenschema (… analog zur bereits in `ADR-0060` benannten
  Folgepflicht für das Protobuf-Schema)". Der umsetzende Slice erweitert
  stattdessen `SPEC-018` (so steht es in seiner §2-DoD), und der Diff folgt
  dem Plan. In der Sache trägt die Form des Plans: der Endpunkt teilt
  Horch-Adresse, Token-Klassen und Aktivierungsvariable mit `SPEC-018`, und
  das Fehlen einer §6-Zeile (kein externer Vertrag — `net/http` ist keine
  Fremdabhängigkeit, anders als NATS in `SPEC-017` und die gRPC-Toolchain in
  `SPEC-020`) nimmt der Analogie genau den Grund, der dort eine eigene
  Kennung rechtfertigte. Was fehlt, ist allein die Benennung der Abweichung
  als Abweichung: kein Artefakt sagt, dass hier bewusst von der in der
  `Accepted`-ADR genannten Form abgewichen wird.
- `verifizierbar`: ja — die drei Textstellen sind zeichengleich lesbar;
  `grep -c "ADR-" ` gegen die neu hinzugefügten Pflichtenheft-Zeilen ergibt 0
  (Spec→ADR-Verbot gewahrt, kein Rückverweis).
- `klasse`: „ADR-Folgepflicht in abweichender Form umgesetzt, Abweichung nicht
  benannt" (1×)

---

## Mutations- und Probenläufe

Alle Proben am **Blatt** `89d31d1` vorgenommen, jede in einem eigenen
Werkzeug-Aufruf gesetzt und ausgewertet; zurückgenommen wurde jeweils per
`git checkout -- <pfad>`. Die Blatt-Identität ist danach einzeln geprüft:
`git hash-object <pfad>` gegen `git rev-parse HEAD:<pfad>` — beide Größen
paaren sich zeichengleich (`wiring.go`
`472785188daf3f139fee9388674cae6e53530c1a`,
`sse.go` `e4805be74bdbd36d0a21bf2fc3c43114935e81b9`) — und
`git status --porcelain` ist leer.

| Probe | Mutation | Erwartung | Ergebnis |
|---|---|---|---|
| M1 | `wiring.go:377` `grpcAddr != "" \|\| httpAddr != ""` → `&&` | `TestChangeStreamEnabled` rot („nur gRPC gesetzt", „nur HTTP gesetzt") | **rot** — zwei Sub-Fälle `FAIL`, `:25` „= false, Erwartung true" |
| M2 | `sse.go:89-92` `503`-Guard entfernt | `TestStreamOhneBroadcasterAntwortetMit503` rot | **rot** — Nil-Interface-Panic in `streamChangesHandler.func1` `sse.go:94`, Client sieht `EOF` |
| M3 | `sse.go` `Content-Type: text/event-stream` entfernt | Happy Path rot | **rot** — `sse_test.go:178` „Content-Type: \"\" (Erwartung: text/event-stream)" |
| M4 | `sse.go` `event: %s\ndata: %s\n\n` → `data: %s\n\n` | Happy Path rot | **rot** — `sse_test.go:178` „event-Typ: \"\" (Erwartung: \"change\")" |
| P1 | `sse.go:49-52` Leer-Verzweigung in `rowImage` entfernt | — | **grün** (`ok`, `-run TestStream` über das ganze Paket) → F-4 |

Die vier rot gesehenen Mutationen decken die drei tragenden Zusagen des
Slices ab: die Oder-Bedingung der Bootstrap-Entkopplung (M1 — sie färbt genau
auf den zwei Fällen rot, die `ADR-0061` Teilfrage 5 fordert), den
`503`-Guard gegen Nil (M2) und die von `SPEC-018` deklarierte Drahtform
(M3/M4: `text/event-stream` und `event: change`). P1 ist die eine Probe, die
**nicht** rot wurde, und trägt F-4.

## Sensoren dieses Laufs (Exit-Code je eigenem Schritt, mechanisch konditioniert, nie gepiped)

| Lauf | Exit | Beleg |
|---|---|---|
| `make a-check` | **0** | `gesamt: 0 Befund(e)` |
| `make gates` | **0** | `baseline-verify` grün · `docs-check` 566 Dateien / 0 Befunde · `a-check` 0 Befunde · `commit-traceability` OK (5 Commits, Betreffs ohne Struktur-ID) · `coverage-gate` „Coverage 49.40% erfüllt Schwelle 35%" |
| `make test` | **0** | 28 Pakete `ok`, kein `FAIL` (Race-Detector, netzlos im gepinnten Toolchain-Container) |
| `bash -n tools/harness/run-integration-tests.sh` | **0** | Syntaxprobe des erweiterten Skripts |

**Nicht ausgeführt: `make test-integration`.** Der reale SSE-E2E-Rundlauf ist
nicht Gegenstand dieses Laufs (der Auftrag nennt die drei inneren Sensoren);
seine Aussagen in §Geprüfte Zusagen unten sind am Skript gelesen, nicht real
nachgefahren. Die DoD-/Beleg-Konformität dieses Punktes prüft der Verifier.

## Geprüfte Zusagen ohne Befund

- **Bootstrap-Entkopplung (`ADR-0061` Teilfrage 5), real nachgeprüft:**
  `changeStreamEnabled` (`wiring.go:376-378`) ist die Oder-Bedingung über
  beide Adressen; die Konstruktion hängt an genau dieser Funktion (`:566-569`),
  und **derselbe** Wert `changeBroadcaster` geht an **beide** Driving-Adapter
  (`:695` HTTP, `:723` gRPC). Sind beide Adressen leer, entsteht kein
  `Broadcaster` und `capture.WithChangeStream` wird nicht angehängt — genau
  das Bestandsverhalten, das die ADR zusagt. Der gRPC-only-Pfad bleibt
  bit-identisch (dieselbe Bedingung, die vorher `cfg.GRPCAddr != ""` lautete,
  ist in der Oder-Bedingung enthalten). `TestChangeStreamEnabled`
  (`internal/bootstrap/changestream_internal_test.go:17-29`) deckt alle vier
  Adress-Kombinationen ab und ist rot-fähig (M1).
- **SSE-Endpunkt gegen `ADR-0061` Teilfrage 2/3/4 und `SPEC-018`:**
  `Content-Type: text/event-stream`, `Cache-Control: no-cache`,
  `event: change` je Change, `http.Flusher` nach jedem Event — deckungsgleich
  mit der neuen `SPEC-018`-Tabelle. Der `Last-Event-ID`-Header wird **nicht**
  gelesen und **nicht** gesendet (kein `id:` im Frames-String) — die
  Begründung der ADR (Teilfrage 3: identische Zustellsemantik zu
  `ADR-0060`, Store bleibt Nachhol-Pfad) trägt und ist im Handler-Godoc
  referenziert. Auth läuft unverändert über `withToken(…, roleReader, …)`
  (`server.go:97-98`); `401` fällt vor dem Handler und damit vor jedem Event,
  belegt durch zwei Tests (fehlender und unbekannter Token) plus die
  Assertion „keine Subskription trotz `401`".
- **Adapter-Importrichtung (`ADR-0068`/`.a-check.yml`):** der neue Handler
  importiert ausschließlich `port/outbound` und `domain/model`; die lokal
  deklarierte `changeSubscriber`-Schnittstelle (`sse.go:24-26`) ist dieselbe
  Form wie in `internal/adapters/driving/grpc/server.go:29-31` und folgt dem
  Muster aus `slice-069` — ein Import des `grpcstream`-Pakets wäre die
  unerlaubte `adapters → adapters`-Kante. `.a-check.yml` braucht keine
  Änderung (Glob `adapters: ["internal/adapters/**"]`), und der Wegwerf-Client
  `tools/harness/sseclient` importiert kein internes Paket, fällt also nicht
  unter die `tooling → adapters`-Kante. `make a-check` grün, `0 Befunde`.
- **Die doppelte lokale Deklaration** (gRPC- und HTTP-Adapter) ist die von
  `ADR-0060` Teilfrage 2 vorgesehene Form; der neue Kommentar nennt die
  Gegenstelle namentlich (`sse.go:30-32`), die Abwesenheit eines
  „beide Fassungen zusammen ändern"-Satzes (wie in `middleware.go:28-33` für
  `classifyToken`) ist damit getragen — kein Befund.
- **Realer Rundlauf (`run-integration-tests.sh:1905-2041`), gelesen:**
  der Block liegt **nach** dem gRPC-Rundlauf und **vor** dem
  Upgrade-Sicherheits-Container-Tausch (`:2043`) sowie **vor** der
  Container-Ende-Grenze (`TestE2ESchemaChangeDropColumn`, `:2140-2160`) — die
  Zusage des Blockkommentars stimmt mit der Position überein. Die Assertion
  bindet die empfangene `change_id` wie bei `slice-071` an `cdc.changes`
  (`source_id`, `change_id`, `table_name`, `new_data->>'name'`, `:2025-2032`),
  und die `401`-Ablehnung wird real gefordert (Client-Ausgang 0 **und**
  `REJECTED code=401` in den Logs, `:1968-1975`, `:1993-1996`). Der Client
  abonniert vor dem Header-Flush (Handler: `Subscribe()` vor `WriteHeader`),
  die `READY`-Zeile ist damit kein Frühstart — das Review zu `slice-071`
  F-6 findet hier keine Entsprechung.
- **Betriebs-Oberfläche:** kein neues `CDC_*`-Feld, keine Änderung an
  `compose.yaml`, keine Änderung an `docs/user/benutzerhandbuch.md` — die
  Versionshistorie-Regel ist damit nicht ausgelöst (sie greift erst beim
  Anfassen des Handbuchs); die Mitnahme-Frage trägt F-1.
- **Kommentar-Disziplin (`AGENTS.md` §3.7), `grep`-basiert:** kein
  `slice-`/`welle-`-Verweis und keine Vorher/Nachher-Sprache in einer
  hinzugefügten Produktionscode-Zeile des Diffs
  (`git show 89d31d1 -- internal/ tools/ | grep -E "^\+"` gegen
  `slice-[0-9]|welle-[0-9]|früher|vorher|bislang|nicht mehr|wäre|würde|ehemals|vormals`
  → ein Treffer, und der ist die Zustandsbeschreibung
  „Feed-Container lief … nicht mehr weiter", die dem Muster der bestehenden
  Blöcke folgt). Die Herkunfts-Anker `· seit slice-072` in
  `harness/README.md:133` und in der `SPEC-018`-Historie sind die
  zulässige Form. Kein `//nolint`/`#noqa`/`[SuppressMessage]` im Diff. Kein
  `ADR-`-Rückverweis in den hinzugefügten Pflichtenheft-Zeilen
  (Spec→ADR-Verbot gewahrt; `grep -c "ADR-"` → 0).
- **Traceability/ID-Schema:** Commit-Betreff trägt `LH-FA-SST-008` und
  `ADR-0061`, keine `SPEC-*`/`ARC-*`-Kennung; `make commit-traceability` grün.
  `harness/image-hash.txt` ist um genau eine Digest-Zeile fortgeschrieben
  (`ADR-0044`: der Zug ändert Build-Kontext-Dateien).
- geprüft, ohne Befund: `internal/adapters/driving/http/` (Handler, Routing,
  Middleware-Wiederverwendung, 8 Tests)
- geprüft, ohne Befund: `internal/bootstrap/` (Oder-Bedingung, beide
  Adapter-Verdrahtungen, Shutdown-Pfad)
- geprüft, ohne Befund: `tools/harness/` (sseclient, Rundlauf-Position und
  -Bindung, Aufräumpfade)
- geprüft, ohne Befund: `harness/` (README-Sensorzeile, image-hash.txt)
- geprüft, ohne Befund: `docs/plan/planning/in-progress/slice-072-…md`
  (§3-Nachzug, DoD-Stand)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 6 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** „Aufschub-Adresse nimmt die Sendung nicht an"
(F-1, 2×) · „Adapter-Unit-Test konstruiert einen Zustand, den der Bootstrap
nicht herstellt" (F-2, 2×) · „Test-Zusage ohne Bindung an den Produktionspfad"
(F-3, F-4, 1×) · „`SPEC-018`-Textstellen der Erweiterung nicht nachgezogen"
(F-5, 1×) · „ADR-Folgepflicht in abweichender Form umgesetzt, Abweichung nicht
benannt" (F-6, 1×)

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM. Kein Finding berührt die
Richtung des Slice: die Bootstrap-Entkopplung ist real umgesetzt und mit dem
selben Wert an beide Adapter gebunden (M1), der `503`-Guard trägt als
Grenz-Bedingung (M2), die Drahtform des Endpunkts ist von zwei Mutationen
gedeckt (M3/M4), der reale Rundlauf steht an der richtigen Stelle und bindet
die empfangene `change_id` an den Lesezugriffsweg.

**Zu den drei Fragen, die der Implementer offen benannt hat:**

1. **Bootstrap-Entkopplung — trägt.** `changeStreamEnabled` ist die eine
   Bedingung, ihre Wirkung ist an beiden Adaptern dieselbe Variable, und der
   Fall „beide leer" führt zum unveränderten Bestand. Der `503`-Guard ist
   **kein toter Code** — er ist die von `ADR-0061` Teilfrage 5 verlangte
   Grenz-Bedingung an einem exportierten `Config`-Feld und verhindert einen
   Nil-Panic —, aber seine Auslöse-Bedingung hat in `Run` keinen Erzeuger
   (F-2). Die Einordnung des Implementers trägt; die Grenze ist mit F-2
   benannt statt still.
2. **Der schwächere Punkt (Abweichung 2) — die DoD-Zeile trägt, der Testname
   nicht.** Die Zusage „ein `Publish` ohne verbundenen SSE-Client blockiert
   nicht" hat zwei Träger: am Adapter selbst die (echte, geprüfte) Aussage,
   dass ohne Request keine Subskription besteht, und in der Tiefe
   `internal/adapters/driven/grpcstream` (`TestPublishOhneAbonnentenBlockiertNicht`,
   `TestNichtLesenderAbonnentHaeltPublishNichtAn`) die eigentliche
   Nicht-Blockade. Zusammen ist die Anforderung belegt; was nicht trägt, ist
   der Anspruch des Testnamens, diese Hälfte selbst zu binden (F-3).
3. **Handbuch-Adresse — bedingt tragfähig (F-1).** `slice-077` ist der
   richtige Empfänger, seine §2 zählt aber zwei Netzwerk-Zugriffswege und
   nennt den SSE-Endpunkt nicht; die Sendung wird damit nicht ausdrücklich
   angenommen.

**Übergabe:** Keine Rückgabe an den Implementer — keine Fixrunde. F-1 ist eine
Plan-Zeile in `slice-077` (Planner), F-5 und F-6 sind Ein-Zeilen-Nachzüge im
`spec/pflichtenheft.md`; alle drei sind textlich, isoliert und ohne
Verhaltenswirkung — für sie wird die Rückkante Review → Implementer nicht
geöffnet (Modul 8 §Konflikt-Pfad, `v6.5.0` ·
`regelwerk/modul-08-agentenrollen.md`: „bei isolierten LOW/INFO-Findings ist
die Sequenz Overkill"). F-2, F-3 und F-4 sind benannte Grenzen ohne
Handlungsbedarf am Diff. Die **Finding-Klassen** gehen in die Slice-Closure §7
und von dort in den Zähler.

**DoD-Nachzug:** Da keine Fixrunde folgt, zieht dieser Report die Zeile „Review
durchgeführt, Report unter `docs/reviews/` liegt vor" in §2 des Slice-Plans
selbst auf `[x]` nach, im selben Commit wie dieser Report
(`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde). Nur diese
eine Zeile; Closure-Notiz, Risiko-Ausgänge, Beobachtungs-Register und die drei
Paarungen bleiben Planner-Arbeit. Der Report ersetzt keine Verifikation —
DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).
