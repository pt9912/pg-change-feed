# Slice sdk-tls-optionen: die drei SDKs verbinden über TLS mit eigenem Vertrauensanker

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** ohne Welle — es gibt keine Closure-Bedingung jenseits der DoD
dieses Slice, siehe Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht (Modul 6).

**Abhängigkeit:** harte Voraussetzung: Lastenheft 0.16.0 mit `LH-FA-SST-013` und
`SPEC-037` im Pflichtenheft sind committet (die Folge-Anforderung, auf die
`slice-examples-grpc-tls` wartete). Server-seitig liegt TLS vor
([`slice-tls-http-grpc-server`](../done/slice-tls-http-grpc-server.md)); die
Beispiele unter `examples/` hängen nicht an den SDKs und sind nicht Teil dieses
Slice.

**Bezug:** [`LH-FA-SST-013`](../../../../spec/lastenheft.md) (Scope),
[`LH-FA-SST-009`](../../../../spec/lastenheft.md) (die Packages),
[`LH-FA-SST-011`](../../../../spec/lastenheft.md) (Serverseite),
[`ADR-0150`](../../adr/0150-tls-und-mehrfach-token.md),
[`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md) (Mechanik der
Realserver-Phasen).

**Berührte Spec-Stellen:** `SPEC-037` · `SPEC-026` · `SPEC-027` · `SPEC-028`
(Package-Versionen, Minor-Sprung).

**Verantwortlich:** — bis zur Priorisierung.

**Autor:** pt9912 (Architect-Zuschnitt). **Datum:** 2026-10-05.

---

## 1. Ziel und Abgrenzung

**Ziel:** Jedes der drei SDK-Packages (C#, Python, Kotlin) nimmt eine Option für
einen eigenen Vertrauensanker (PEM-Datei) an und bedient damit HTTP, SSE,
gRPC-Stream und gRPC-Verwaltung gegen einen Server mit TLS; ein Realserver-Lauf
je SDK belegt es gegen ein Zertifikat, das nicht im Systemspeicher steht.

**Entscheidungen des Zuschnitts** (Quelle: Pflichtenheft `SPEC-037`):

- Eine Option, ein Inhalt: Pfad einer PEM-Datei; mit Anker vertraut die
  Verbindung genau diesen Zertifikaten, ohne Anker gilt der Systemspeicher.
- Die Prüfung ist nie abschaltbar; kein Überschreiben des Servernamens.
- NATS-Fläche bleibt unberührt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Client-Zertifikate (mTLS) — nicht gefordert (`LH-FA-SST-013` Out-of-Scope);
  der Server verlangt keines (`SPEC-034`).
- Eine „unsichere“ Option (Prüfung aus) und ein Servernamen-Override — bewusst
  nicht angeboten; der Override wäre in der JDK-HTTP-Schicht nur mit einem
  eigenen Trust-Manager machbar und öffnet den Weg zur abgeschalteten Prüfung.
  Wer den Namen nicht treffen kann, stellt das Server-Zertifikat auf den Namen
  der Adresse aus. Ein Folge-Wunsch ist eine eigene Anforderung.
- TLS der NATS-Verbindung — Sache der NATS-Konfiguration (wie `LH-FA-SST-011`).
- Die Beispiele und die Server-seitige Konfiguration — anderer Gegenstand
  (`slice-examples-grpc-tls`, `slice-tls-http-grpc-server`).
- Veröffentlichung der neuen Versionen (Tags, Registries) — braucht eine neue
  Freigabe des Auftraggebers; der Slice endet mit gebauten Artefakten
  (`make sdk-pack-*`).

## 2. Definition of Done

- [ ] **C#:** `PgChangeFeedClientOptions` trägt die Anker-Option; der HTTP-, SSE-
      und gRPC-Weg (Stream und Verwaltung) nutzt sie; Unit-Tests: Happy
      (Verbindung über TLS mit gültigem Anker), Boundary (kein Anker → Systemspeicher,
      `http://` unverändert), Negative (fremdes Zertifikat, falscher Name,
      nicht lesbarer Pfad, Datei ohne PEM → sichtbarer Fehler, kein Rückfall).
      `LH-FA-SST-013` erfüllt, Test referenziert.
- [ ] **Kotlin** und **Python:** wie C#, je Sprache dieselben Testfälle und
      derselbe Inhalt der Option (`SPEC-037`).
- [ ] **Realserver-Beleg:** je SDK eine TLS-Phase im Runner
      (`make test-sdk-csharp-integration`, `…-kotlin-…`, `…-python-…`): der Feed-Container
      läuft mit Zertifikat aus `tools/harness/certgen` (Muster der TLS-Phase von
      `make test-integration`, Override-Datei im Temp-Verzeichnis, `compose.yaml`
      bleibt unverändert); je Fläche empfängt/ruft der Client mit Anker; ohne
      Anker scheitert dieselbe Verbindung sichtbar. Der Abdeckungs-Träger
      `docs/user/sdk-e2e-abdeckung.md` trägt die neuen Phasen (Runner schreibt ihn).
- [ ] `make gates` grün (inkl. `make sdk-public-doc-check`: keine interne Kennung in
      Quellen, Tests, README der SDKs).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein Self-Review.
- [ ] Doku-Update: README der drei Packages (Englisch, Abschnitt zu TLS; die
      Ist-Aussagen „plaintext gRPC“ in Docstrings und README werden nachgezogen),
      Benutzerhandbuch nur, falls es die SDK-Nutzung beschreibt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register fortgeschrieben oder „keine Beobachtung“ in §7.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen.

(Liefer-Punkte: drei — je SDK Option, Tests, Realserver-Phase.)

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/csharp/PgChangeFeed.Client/PgChangeFeedClientOptions.cs` und `Http/`, `Sse/`, `Grpc/` | update | Option und Verdrahtung in `HttpClientHandler`/`SocketsHttpHandler` (benutzerdefinierte Vertrauenskette) und `GrpcChannel` (`HttpHandler`) — [`LH-FA-SST-013`](../../../../spec/lastenheft.md) |
| `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/PgChangeFeedClientOptions.kt` und `http/`, `sse/`, `grpc/` | update | Option; `SSLContext`/Trust-Manager für `HttpTransport`/`SseTransport`, `TlsChannelCredentials` für gRPC statt `usePlaintext()` bei TLS-Adresse |
| `sdks/python/pgchangefeed/src/pgchangefeed/` (`http_client.py`, SSE-Client, `grpc_client.py`, `administration_client.py`) | update | Option; `ssl`-Kontext für HTTP/SSE, `grpc.ssl_channel_credentials(root_certificates=…)` für gRPC |
| Tests je SDK (Unit, neue Testdateien) | neu | Happy/Boundary/Negative nach `LH-FA-SST-013`; Zertifikate zur Laufzeit erzeugt oder als kleine Test-Fixtures außerhalb der Veröffentlichung |
| `tools/harness/run-sdk-{csharp,kotlin,python}-integration-tests.sh` und `tools/harness/lib-*` | update | TLS-Phase je Runner, Wiederverwendung von `tools/harness/certgen`; **kein Teil des aktuellen parallelen Zuges an `tools/harness`** — erst nach dessen Abschluss anfassen |
| SDK-READMEs, Docstrings | update | Ist-Zustand TLS, keine interne Kennung (`make sdk-public-doc-check`) |
| Package-Versionen (`.csproj`, `build.gradle.kts`, `pyproject.toml`) | update | additive Option → Minor-Sprung; die Anhebung gehört in den Slice, das Veröffentlichen nicht (§1) |

```suchlauf
71c6a886 9 -i -E 'plaintext|no TLS' -- sdks
```

Stand der Zeile: **gemessen** am Parent `71c6a886` mit `git grep -n` (9 Trefferzeilen);
der Implementer wiederholt sie mit `bash tools/harness/suchlauf-nachmessen.sh <Plan-Datei>`
und prüft am Endstand, dass die Trefferzeilen der Ist-Aussage „plaintext“ in Docstrings
und README nachgezogen sind (Zahlen und Stände, nicht die Vollständigkeit des Musters).

## 4. Trigger

**Start** (`next` → `in-progress`): Lastenheft 0.16.0 und `SPEC-037` committet;
der parallele Zug an `tools/harness` ist geschlossen oder berührt die Runner
nicht; ein geladenes `:dev`-Image (`make image`).

**Rückführungen — vorab benannt:**

- `in-progress` → `next` (zu groß): die drei Liefer-Punkte tragen nicht in einem
  Review — dann je Sprache ein Slice (`slice-sdk-tls-optionen-csharp`, `…-kotlin`,
  `…-python`), dieser Plan wird zum ersten, die beiden anderen entstehen in `open/`.
- `in-progress` → `open` (blockiert): eine Sprache kann den Anker ohne
  Abschalten der Prüfung nicht an HTTP **und** gRPC binden — dann Befund an
  `SPEC-037`, Folge-ADR.

## 5. Closure-Trigger

DoD vollständig, `make gates` grün, Review- und Verifier-Bericht liegen vor,
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Die JDK-/Java-HTTP-Schicht und das Python-SSE brauchen jeweils einen eigenen
  Weg zum Trust-Store — **Ausgang:** weiter offen bis zur Umsetzung; nicht lösbar
  heißt Rückführung (§4).
- Die Realserver-Runner liegen unter `tools/harness`, das ein anderer Zug gerade
  ändert — **Ausgang:** weiter offen; entfällt, wenn der Slice nach dessen
  Abschluss startet.
- Der Suchlauf trifft Symbolnamen, nicht verschobene Zahlen (`AGENTS.md` §3.13
  Grenze) — **Ausgang:** Lese-Handlung des Reviewers.

## 7. Closure-Notiz

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register:** —
- **Folge-Slices:** —
- **Risiken aus §6:** —
- **Drei Paarungen:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist `*` (Default; `sdks/`, `tools/harness`,
`docs/user/`); eine Ausdifferenzierung ist nicht nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; kein Treffer
für TLS in den SDKs.

Alle berührten Sub-Areas GF.
