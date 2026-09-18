# Welle nats-drittstream: NATS als dritter, paralleler Vollinhalts-Zustellweg für Live-Streaming

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-nats-drittstream-results.md`). Der Zustand ist die
Verzeichnis-Position — kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** —. **Datum:** 2026-09-18.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

[`LH-FA-SST-008`](../../../spec/lastenheft.md) bekommt einen **dritten**,
unabhängig nutzbaren Zustellweg für vollständige Change-Inhalte —
[`ADR-0100`](../adr/0100-nats-dritter-vollinhalts-zustellweg.md) hat die
Entscheidung bereits `Accepted` getroffen (Publisher als dritter
`Broadcaster`-Abonnent, eigener Subjekt-Namensraum `cdc.stream.<source_id>.<schema>.<table>`,
Core NATS ohne Replay, ein geteilter, opt-in-pflichtiger Verbindungs-Token
`CDC_NATS_STREAM_TOKEN`). Das *Mehr* gegenüber den einzelnen Slice-DoDs: Ein
NATS-Client mit gültigem Token empfängt real eine vollständige Change über
den neuen Subjekt-Namensraum, während ein Client ohne (oder mit falschem)
Token vom Server selbst abgelehnt wird — **und** dieselbe Fähigkeit steht
danach in allen drei Beispiel-Sprachen (Go, C#, Kotlin) bereit, mit
Doku-Nachzug in `examples/README.md` und
`docs/user/benutzerhandbuch.md`. Kein Einzel-Slice-DoD trägt beide Hälften
(realer Server-Beleg + volle Sprachmatrix) zusammen.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- [`ADR-0100`](../adr/0100-nats-dritter-vollinhalts-zustellweg.md) trägt
  Status `Accepted` (bereits eingetreten — ein anderer Mensch liest den
  Status-Header der Datei).
- `welle-beispiele-start-ueber-make` liegt in `done/` — die Beispiel-Client-
  Startform (`make example-run-go/-csharp/-kotlin SURFACE=…`,
  [`ADR-0098`](../adr/0098-beispiel-clients-start-ueber-make-dockerfile.md))
  existiert bereits und wird von dieser Welle nur um einen vierten
  `SURFACE`-Wert erweitert, nicht neu geschaffen (bereits eingetreten —
  [`welle-beispiele-start-ueber-make-results.md`](done/welle-beispiele-start-ueber-make-results.md)).

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Alle drei Slices (`slice-nats-drittstream-core`,
  `slice-nats-drittstream-example-go`,
  `slice-nats-drittstream-example-csharp-kotlin`) liegen in `done/`.
- `make gates` grün.
- Ein realer `make test-integration`-Lauf zeigt **beide** Belege aus
  [`ADR-0100`](../adr/0100-nats-dritter-vollinhalts-zustellweg.md)s Fitness
  Function in einem Durchlauf: ein Client mit gültigem
  `CDC_NATS_STREAM_TOKEN` empfängt eine vollständige Change über
  `cdc.stream.<source_id>.<schema>.<table>`, ein Client ohne/mit falschem
  Token wird vom NATS-Server abgelehnt, und das bestehende Wecksignal
  (`natssub`) funktioniert mit demselben Test-Token unverändert weiter.
- Closure-Notiz in `welle-nats-drittstream-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-nats-drittstream-core | `natsstream.Publisher` (dritter `Broadcaster`-Abonnent), Bootstrap-Verdrahtung (`CDC_NATS_STREAM_TOKEN`, Zwei-Bedingungen-Aktivierung), `compose.yaml`-NATS-Auth, Wegwerf-Belegträger, `make test-integration`-Erweiterung | [`LH-FA-SST-008`](../../../spec/lastenheft.md), [`ADR-0100`](../adr/0100-nats-dritter-vollinhalts-zustellweg.md) |
| slice-nats-drittstream-example-go | `examples/nats-stream-client` (Go), `harness/mk/examples.mk`-Erweiterung um `SURFACE=nats-stream` für `example-run-go`, Doku-Nachzug | [`LH-FA-SST-008`](../../../spec/lastenheft.md), [`ADR-0100`](../adr/0100-nats-dritter-vollinhalts-zustellweg.md) |
| slice-nats-drittstream-example-csharp-kotlin | `examples/csharp/nats-stream-client`, `examples/kotlin/nats-stream-client`, `harness/mk/examples.mk`-Erweiterung um `SURFACE=nats-stream` für `example-run-csharp`/`-kotlin`, `examples-csharp`/`examples-kotlin`-Bau-Erweiterung, Demo-Umgebung (`examples/compose.yaml`, `examples/.env`), Doku-Nachzug | [`LH-FA-SST-008`](../../../spec/lastenheft.md), [`ADR-0100`](../adr/0100-nats-dritter-vollinhalts-zustellweg.md) |

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- **Reihenfolge — explizit:** `slice-nats-drittstream-core` läuft **zuerst**
  und muss `done/` erreichen, bevor `slice-nats-drittstream-example-go` oder
  `slice-nats-drittstream-example-csharp-kotlin` beginnen — ohne einen
  realen `natsstream.Publisher` gäbe es nichts, gegen das ein Beispiel-Client
  laufen könnte (dieselbe Abhängigkeitsform wie `ADR-0100`s eigene
  Slice-Schnitt-Empfehlung, Slice A → Slice B).
- `slice-nats-drittstream-example-go` und
  `slice-nats-drittstream-example-csharp-kotlin` sind **voneinander
  unabhängig** — beide hängen ausschließlich an
  `slice-nats-drittstream-core`, nicht aneinander (dieselbe Form wie
  `slice-beispiele-go-dockerfile-start`/`slice-beispiele-csharp-kotlin-start-target`
  in [`welle-beispiele-start-ueber-make`](done/welle-beispiele-start-ueber-make.md)
  §5) und können parallel laufen.
- Wird blockiert von: nichts außerhalb dieser Welle — `ADR-0100` liegt bereits
  `Accepted` vor, `ADR-0060`/`ADR-0061` (Broadcaster-Ursprung) liegen in
  `done/` (`welle-19`).
- Blockiert: nichts — kein bekannter Folge-Vorgang hängt an dieser Welle.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Jede Änderung an `ADR-0055`/`ADR-0056` (NATS-Wecksignal) oder an
  `internal/adapters/driven/natsnotify/`** — [`ADR-0100`](../adr/0100-nats-dritter-vollinhalts-zustellweg.md)
  hält beide ausdrücklich unverändert und byte-identisch in Kraft; ein
  Zeilen-Diff dort außerhalb der reinen Client-Options-Erweiterung an der
  gemeinsamen `nats.Connect`-Aufrufstelle (Teilfrage 5) wäre ein
  Scope-Bruch.
- **NATS-Accounts/Permissions-Konfiguration** (`ADR-0100` Teilfrage 4 Option
  B, subjekt-scoped Zugriffskontrolle) — bewusst vertagter
  Re-Evaluierungs-Trigger der ADR, kein Bestandteil dieser Welle; der
  serverweite, geteilte Token bleibt der einzige Auth-Mechanismus.
- **JetStream/Replay für den neuen Stream** — `ADR-0100` Teilfrage 3 wählt
  Core NATS ohne Replay bewusst; ein Re-Evaluierungs-Trigger existiert dort,
  nicht hier.
- **Erweiterung des bestehenden `nats-client`-Beispiels um einen zweiten
  Modus** — verworfen in `ADR-0100` Teilfrage 6 Option C zugunsten
  eigenständiger `nats-stream-client`-Verzeichnisse je Sprache.
- **`CaptureService`-Änderungen** — `ADR-0100` Teilfrage 1 stellt fest, dass
  der dritte Weg ohne neue Konstruktions-Option auskommt; jede
  `CaptureService`-Änderung in dieser Welle wäre ein Scope-Bruch gegen die
  eigene Grundlage.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln — dort die **Ruheort-Regel**: Die
beiden Zeiger unten sind so zu schreiben, wie sie vom Ruheort `done/` auflösen,
nicht vom Schreibort.

Ergebnis: `welle-nats-drittstream-results.md`
Zähler: `../observations/`
