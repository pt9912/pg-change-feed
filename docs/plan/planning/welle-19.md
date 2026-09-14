# Welle 19: Live-Change-Streaming — gRPC-Server-Streaming und HTTP/SSE (`LH-FA-SST-008`)

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-19-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** —. **Datum:** 2026-09-14.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

[`LH-FA-SST-008`](../../../spec/lastenheft.md) wird über **zwei** parallele,
unabhängig funktionsfähige Zustellwege erfüllbar gemacht — gRPC
Server-Streaming ([ADR-0060](../adr/0060-grpc-streaming-mechanismus.md))
und HTTP/SSE ([ADR-0061](../adr/0061-http-sse-zusaetzlich-zu-grpc.md)) —,
beide gespeist vom selben In-Prozess-`Broadcaster` hinter dem neuen
Outbound Port `ChangeStreamPort`. Das *Mehr* gegenüber den einzelnen
Slice-DoDs: Ein realer E2E-Beleg zeigt, dass eine einzelne committed
Änderung **gleichzeitig** über beide Wege einen verbundenen Client mit
vollständigem Inhalt erreicht — kein einzelner Slice-DoD belegt das
Zusammenspiel beider Wege über denselben `Broadcaster`.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- [ADR-0060](../adr/0060-grpc-streaming-mechanismus.md) liegt `Accepted` vor
  — erfüllt.
- [ADR-0061](../adr/0061-http-sse-zusaetzlich-zu-grpc.md) liegt `Accepted`
  vor — erfüllt.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle vier Slices (`slice-069`…`slice-072`) liegen in `done/`.
- `make gates` grün.
- Ein real belegter gRPC-E2E-Rundlauf (`slice-071`, `make test-integration`)
  zeigt: eine committed Änderung erreicht einen verbundenen gRPC-Client mit
  vollständigem Inhalt; ein Verbindungsversuch ohne gültiges Token wird
  abgelehnt.
- Ein real belegter SSE-E2E-Rundlauf (`slice-072`, `make test-integration`)
  zeigt dasselbe über `GET /changes/stream`.
- Closure-Notiz in `welle-19-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-069 | gRPC-Streaming-Adapter-Grundgerüst | [LH-FA-SST-008](../../../spec/lastenheft.md), [ADR-0060](../adr/0060-grpc-streaming-mechanismus.md) |
| slice-070 | gRPC-Capture-Integration | [LH-FA-SST-008](../../../spec/lastenheft.md), [ADR-0060](../adr/0060-grpc-streaming-mechanismus.md) |
| slice-071 | gRPC-Beispiel-Client und E2E-Beleg | [LH-FA-SST-008](../../../spec/lastenheft.md), [ADR-0060](../adr/0060-grpc-streaming-mechanismus.md) |
| slice-072 | HTTP/SSE-Streaming-Endpunkt | [LH-FA-SST-008](../../../spec/lastenheft.md), [ADR-0061](../adr/0061-http-sse-zusaetzlich-zu-grpc.md) |

**Reihenfolge — explizit:** `slice-069` (Port/Broadcaster/Server-Grundgerüst)
vor `slice-070` (Capture-Integration, `WithChangeStream`) vor `slice-071`
(gRPC-Beispiel-Client/E2E). `slice-072` hängt **ausschließlich** von
`slice-069` **und** `slice-070` ab — Port und `Broadcaster` müssen
existieren **und** in `CaptureService` verdrahtet sein, sonst sieht der
SSE-Endpunkt nie ein Event. `slice-072` ist **nicht** von `slice-071`
abhängig: Der gRPC-Beispiel-Client/E2E-Beleg und der SSE-Beleg sind
voneinander unabhängige Nachweisformen für zwei verschiedene Konsumenten
desselben `Broadcaster` ([ADR-0061](../adr/0061-http-sse-zusaetzlich-zu-grpc.md)
§Slice-Schnitt-Empfehlung). `slice-072` kann also parallel zu `slice-071`
laufen, sobald `slice-070` in `done/` liegt.

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine nachfolgende Welle ist derzeit geplant (*Nächste Wellen*
  in der Roadmap ist leer).
- Wird blockiert von: keine — beide tragenden ADRs
  ([ADR-0060](../adr/0060-grpc-streaming-mechanismus.md),
  [ADR-0061](../adr/0061-http-sse-zusaetzlich-zu-grpc.md)) liegen `Accepted`
  vor.
- Intern (siehe §4 *Reihenfolge — explizit*): `slice-070` setzt `slice-069`
  voraus; `slice-071` setzt `slice-070` voraus; `slice-072` setzt `slice-069`
  **und** `slice-070` voraus, nicht `slice-071`.
- Nebenläufigkeit außerhalb dieser Welle: `slice-061` (HTTP-API-Beispiel-
  Client/E2E, [ADR-0057](../adr/0057-http-grpc-api.md)) läuft parallel in
  `in-progress/` und berührt
  `internal/adapters/driving/http/`, `compose.yaml`,
  `tools/harness/run-integration-tests.sh`, `harness/README.md`,
  `internal/bootstrap/wiring.go`, `tools/harness/httpclient/` — `slice-072`
  berührt dasselbe Paket (`internal/adapters/driving/http/`) und dieselben
  Dateien; sein Start ist zusätzlich an das Ende von `slice-061` gebunden
  (Merge-Konflikt-Vermeidung, kein Trigger im Sinne von §2/§3).

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Stream-internes Replay** (gRPC- oder SSE-seitig, Fortsetzung ab einer
  Consumer-Position innerhalb des Streams selbst) — beide ADRs verwerfen
  das bewusst (Teilfrage 3 in beiden) und benennen je eine eigene
  Re-Evaluierungs-Trigger-Bedingung; nicht Gegenstand dieser Welle.
- **Tabellen-granulare Filterung** des gRPC- oder SSE-Streams (analog zu
  [ADR-0056](../adr/0056-nats-tabellen-granulares-subjekt.md)s
  NATS-Granularität) — dieselbe Vertagung, dieselbe
  Re-Evaluierungs-Bedingung in beiden ADRs.
- **Änderungen an `ChangeNotificationPort`/`natsnotify`** (NATS bleibt
  unverändert eigenständiges Wecksignal,
  [LH-FA-SST-007](../../../spec/lastenheft.md)).
- **Änderungen an der administrativen HTTP-Verwaltungs-API**
  ([ADR-0057](../adr/0057-http-grpc-api.md)) jenseits des einen neuen
  `GET /changes/stream`-Endpunkts — bestehende
  Endpunkte, Middleware-Verhalten und Fehler-Antwortform bleiben
  unangetastet.
- **`.a-check.yml`** — beide ADRs stellen fest, dass keine Änderung nötig
  ist (bestehende Globs decken beide neuen Pakete ab); ein Slice, der
  dennoch eine Änderung vornimmt, hat den Plan geändert, nicht nur
  ergänzt.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: noch offen — erst nach Welle-Abschluss zu füllen (Zeiger auf
`welle-19-results.md`, Geschwister im Ruheort `done/`).
Zähler: noch offen — erst nach Welle-Abschluss zu füllen (Zeiger auf
`../observations/README.md`, eine Ebene über dem Ruheort).
