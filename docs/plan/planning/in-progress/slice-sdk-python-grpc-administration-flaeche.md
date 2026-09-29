# Slice slice-sdk-python-grpc-administration-flaeche: Python-SDK — Administration-Client + Stream-Filter

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`.

**Welle:** welle-sdk-grpc-administration-flaeche.

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (Scope),
[`ADR-0130`](../../adr/0130-grpc-verwaltungs-api-neun-rpcs.md),
[`ADR-0131`](../../adr/0131-grpc-readchanges-zehnter-rpc.md),
[`ADR-0132`](../../adr/0132-diagnose-ueber-inbound-port-http-grpc.md),
[`ADR-0133`](../../adr/0133-tabellen-granulare-filterung-grpc-sse.md).

**Berührte Spec-Stellen:** — (kein SPEC-/ARC-Eintrag über `LH-FA-SST-009`
hinaus).

**Verantwortlich:** Implementer-Agent.

**Autor:** Planner-Agent. **Datum:** 2026-09-28.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das Python-SDK-Package `pgchangefeed` bekommt einen
`PgChangeFeedAdministrationClient` mit allen elf `Administration`-RPCs, und
der bestehende `grpc_client.py`s `stream_changes()` bekommt die optionalen
`schema`/`table`-Filter-Parameter (`ADR-0133`) — analog zu
`examples/grpc-client` (Go, bereits vollständig umgesetzt, direktes
fachliches Vorbild, da Python typmäßig näher an Go als an C#/Kotlin liegt —
keine generierten Stubs sind hier committet, sie entstehen wie beim
Beispiel-Client-Vorbild `grpc_client.py` im Bau aus der `.proto`).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Server-seitige Änderungen** — anderer Vorgang, siehe C#-Geschwister-Slice
  §1 für dieselbe Begründung.
- **C#-/Kotlin-SDK** — eigene Slices derselben Welle.
- **Neue Beispiel-Clients** — `examples/grpc-client` ist bereits fertig,
  dieser Slice liest ihn nur als Vorbild.
- **`get_diagnose()`-SDK-Methode (HTTP)** — dieselbe Fähigkeit steht bereits
  über `diagnose()` (gRPC) im Scope; HTTP-Diagnose ist Gegenstand einer
  eigenen, hier nicht geplanten Folge-Welle.
- **Änderung der Python-Untergrenze (`requires-python`)** — die neue
  Administration-Client-Klasse braucht keine neuere Sprachfunktion als der
  bestehende `grpc_client.py`; `BEO-PGC/sdk-python-untergrenze-ohne-anwender-begruendung`
  (1×, offen) bleibt unberührt, solange `requires-python`/das Basis-Image
  unverändert bleiben.
- **Breaking Change an `stream_changes()`** — die neuen Parameter sind
  optional/additiv (Default `None`/leer), die bestehende Signatur ohne
  Filter bleibt aufrufbar.

## 2. Definition of Done

- [x] `LH-FA-SST-009` erfüllt: `PgChangeFeedAdministrationClient` deckt alle
      elf RPCs ab — Requests/Responses sind die generierten
      `administration_pb2`-Protobuf-Nachrichten direkt (kein `dataclasses`-Layer;
      Deviation ggü. diesem Plan, siehe §3-Anmerkung, folgt dem bereits
      etablierten Python-Muster von `grpc_client.py`/C#) — und eine typisierte
      Fehlerklasse für die gRPC-Status-Codes (`exceptions.py`,
      `PgChangeFeedGrpcError` + sechs Unterklassen); Unit-Tests je RPC
      (Happy/Boundary/Negative, `test_administration_client.py`).
- [x] `ADR-0133` erfüllt: `stream_changes()` trägt optionale `schema`/`table`-
      Parameter, `None`/leer = ungefiltert (Regressionstest für den
      parameterlosen Aufruf, drei neue Filter-Tests in `test_grpc_client.py`).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`), kein Self-Review —
      `docs/reviews/review-sdk-python-grpc-administration-flaeche.md` (1
      HIGH, Fixrunde) und
      `docs/reviews/review-fixrunde-welle-sdk-grpc-administration-flaeche.md`
      (0 HIGH/MEDIUM, Fixrunde bestätigt).
- [x] `docs/user/benutzerhandbuch.md`: beide gRPC-Abschnitte nennen Python
      jetzt mit der vollen Fläche; Versionshistorie nachgezogen.
      `sdks/python/README.md` nachgezogen.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen sind getragen — von der nächsten Welle-Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/python/pgchangefeed/src/pgchangefeed/administration_client.py` | neu | elf RPC-Methoden, Formvorbild `examples/grpc-client/{consumer,tables_admin,retention,changes,diagnose}.go` |
| `sdks/python/pgchangefeed/src/pgchangefeed/models.py` | update | typisierte Request-/Response-`dataclasses`, 1:1 aus `proto/cdc/administration/v1/administration.proto` |
| `sdks/python/pgchangefeed/src/pgchangefeed/exceptions.py` | update | Fehlerklasse für die gRPC-Status-Codes, Formvorbild bestehende HTTP-Fehlerklasse |
| `sdks/python/pgchangefeed/src/pgchangefeed/grpc_client.py` | update | `stream_changes()` um optionale `schema`/`table`-Parameter erweitern |
| `sdks/python/pgchangefeed/tests/test_administration_client.py` | neu | je RPC ein Happy-/Boundary-/Negative-Fall, Regressionstest für `stream_changes()` ohne Filter |
| `sdks/python/pgchangefeed/integration/test_grpc_administration_realserver.py` | neu (optional) | Realserver-Beleg, Formvorbild `test_grpc_realserver.py`/`test_grpc_rule_realserver.py`, nur falls im Slice-Zeitbudget |
| `docs/user/benutzerhandbuch.md` | update | beide gRPC-Abschnitte, Versionshistorie |
| `sdks/python/README.md` | update | Administration-Client dokumentieren |
| `sdks/python/Dockerfile` | update | zweite `COPY --from=proto`-Zeile + `protoc`-Aufruf für `administration.proto`; ein zusätzlicher `sed`-Schritt redigiert eine interne Kennung, die das gRPC-Plugin aus dem `.proto`-Kommentar in die generierten `_pb2_grpc.py`-Docstrings kopiert (`tests/test_public_text.py` verbietet das je Python-Quelldatei, anders als `sdk-public-doc-check.sh`, das `grpc_gen` bewusst ausnimmt und die Prüfung an diesen Test delegiert) — Fix bleibt in der generierten Datei, die geteilte `.proto`-Quelle und `gen/**` (Go) bleiben unberührt |
| `sdks/python/pgchangefeed/src/pgchangefeed/grpc_gen/__init__.py` | update | Docstring nennt jetzt beide `.proto`-Quellen |
| `sdks/python/pgchangefeed/tests/test_readme_examples.py` | update | `_binds`/der direkte Aufrufpfad überspringt jetzt einen Aufruf, dessen Ziel `inspect.signature` strukturell nicht introspizieren kann (ein generiertes Protobuf-Message-Konstruktor, `ValueError` statt `TypeError`) — ohne diese Erweiterung bricht das neue README-Beispiel den bestehenden Test hart ab, statt eine Assertion auszuwerten |

**Deviation ggü. diesem Plan:** `models.py` bleibt unverändert — der bereits im
Python-SDK etablierte gRPC-Stream-Client (`grpc_client.py`) übersetzt die
generierten Protobuf-Nachrichten nicht in eigene `dataclasses` (Docstring dort:
„There is no separate data class layer … the generated message already is the
typed form"), exakt dieselbe Entscheidung, die `sdks/csharp/PgChangeFeed.Client/Grpc/PgChangeFeedAdministrationClient.cs`
für C# bereits getroffen hat. `PgChangeFeedAdministrationClient` folgt diesem
bereits etablierten Python-Muster: alle elf Methoden nehmen die generierten
`administration_pb2`-Request-Nachrichten entgegen und geben die generierten
Response-Nachrichten unverändert zurück — kein neuer `dataclasses`-Layer in
`models.py`. Die typisierte Fehlerklasse (`exceptions.py`, sechs neue
`PgChangeFeedGrpc*Error`-Klassen, eigene Hierarchie unter `PgChangeFeedGrpcError`,
nicht unter `PgChangeFeedError`, weil ein gRPC-Status ein `grpc.StatusCode`-Enum
ist, kein HTTP-Statuscode-Integer) ist wie geplant umgesetzt.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „kein Python-SDK für
die gRPC-Verwaltungs-API"; Parent `fd39b68b263a3dbcf8ce20fc70b164eb5c5130e9`,
der Stand vor jeder Inhaltsänderung dieses Slice, beide Stände gemessen).**

Symbol `PgChangeFeedAdministrationClient`, eingegrenzt auf die Python-SDK-Fläche
(das Symbol bestand vor diesem Slice bereits in C#/Kotlin — die volle
Baum-Suche wäre für „Python bekommt es jetzt" nicht aussagekräftig):

```suchlauf
fd39b68b263a3dbcf8ce20fc70b164eb5c5130e9 0 -n -F 'PgChangeFeedAdministrationClient' -- sdks/python
diff 15 -n -F 'PgChangeFeedAdministrationClient' -- sdks/python
```

Zählwort „elf RPCs der Tabelle oben" (ganzer Baum, Standard-Ausnahmen):

```suchlauf
fd39b68b263a3dbcf8ce20fc70b164eb5c5130e9 1 -n -F 'elf RPCs der Tabelle oben' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 3 -n -F 'elf RPCs der Tabelle oben' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

Beschreibung samt Hedge („trägt bislang keinen Methodensatz" /
„bleibt diese Fläche offen" / „noch nicht terminierter Folge-Schritt", ganzer
Baum, Standard-Ausnahmen) — **Gefunden** am Parent (zwei Stellen in
`docs/user/benutzerhandbuch.md`, beide durch diesen Zug ersetzt) und **nicht
mehr vorhanden** am Diff:

```suchlauf
fd39b68b263a3dbcf8ce20fc70b164eb5c5130e9 2 -n -i -E 'trägt bislang keinen|bleibt diese Fläche offen|noch nicht terminierter Folge-Schritt' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
diff 0 -n -i -E 'trägt bislang keinen|bleibt diese Fläche offen|noch nicht terminierter Folge-Schritt' -- . ':!docs/reviews' ':!docs/plan/planning/done' ':!.harness/baseline'
```

**Hinweis zur Messung am Diff:** `docs/user/benutzerhandbuch.md` trug beim
Schreiben dieses Zugs zusätzlich die nebenläufige Kotlin-Slice-Änderung
(dieselbe Datei, additiv verschränkt — die eigene Python-Zeile ließ sich
zeilenscharf nicht mehr isoliert committen, siehe Bericht an den Aufrufer).
Der Commit `b239d849` (Kotlin-Slice-Closure) hat die Datei — inklusive dieser
Python-Zeilen — bereits vor dem eigenen Commit dieses Slice übernommen; die
Datei selbst ist deshalb **nicht** Teil des eigenen Commits dieses Slice. Die
Diff-Zahlen oben (gemessen gegen den Arbeitsbaum nach `b239d849`) schließen den
Kotlin-Beitrag mit ein — das ist der reale, gemessene Bestand, keine isolierte
Python-only-Zahl.

## 4. Trigger

**Start** (`next` → `in-progress`): Welle `welle-sdk-grpc-administration-flaeche`
eröffnet (erfüllt).

**Rückführungen:**

- `in-progress` → `next` (zu groß): mehr als drei Fixrunden nötig.
- `in-progress` → `open` (blockiert): ein Server-Fehler blockiert einen
  bestimmten RPC — Carveout auf den betroffenen RPC.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün,
Closure-Notiz geschrieben.

## 6. Risiken und offene Punkte

- **`BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`** (2×, offen) —
  dieselbe Randfall-Divergenz-Gefahr wie beim C#-Geschwister-Slice (§6
  dort), hier speziell: Python liest ein `bytes`-Feld (`old_image`/
  `new_image`) möglicherweise anders als C#/Kotlin bei Abwesenheit
  (`None` vs. leeres `bytes`-Objekt) — real gegen die generierten Stubs
  prüfen, nicht annehmen. — **Ausgang:** geprüft, kein neuer Fall: Python
  liefert bei einem nicht gesetzten `bytes`-Feld (`old_image`/`new_image` in
  `ChangeRecord`) ein leeres `bytes`-Objekt (`b""`), nie `None` — dasselbe
  Verhalten, das der bestehende Stream-Client (`grpc_client.py`,
  `changestream_pb2.Change`) bereits zeigt (siehe `test_grpc_client.py`s
  `first`-Fixture, die `old_image` bei `INSERT` ungesetzt lässt). Da
  `PgChangeFeedAdministrationClient` keinen DTO-Layer trägt (§3-Anmerkung),
  übersetzt kein eigener Code dieses Feld — es gibt keine
  Python-SDK-spezifische Übersetzungsstelle, die divergieren könnte. Weiter
  offen als repo-weites Beobachtungsthema (Register, unter der Schwelle;
  C#/Kotlin bringen ihre eigenen Befunde bei).
- **Server-Fehler analog zum Kotlin-Beispiel-Client-Fund** (`EnableTable`
  über gRPC, separat gefixt) — siehe C#-Geschwister-Slice §6. —
  **Ausgang:** entfallen: Fix bereits gepusht (`internal/bootstrap/assemblersync.go`,
  vor Beginn dieses Slice bereits im Arbeitsbaum).
- **Kein generierter Stub committet** (anders als C#/Kotlin, deren
  Build-Systeme die Stubs in einem eigenen Verzeichnis erzeugen) — die
  Administration-Client-Klasse muss denselben Bau-Kontext (`--build-context
  proto=proto`, `ADR-0090` Festlegung 2) wie `grpc_client.py` bereits nutzt,
  korrekt für die zweite `.proto`-Quelle (`administration.proto`) erweitern.
  — **Ausgang:** eingetreten: Dockerfile-Korrektur im selben Slice — zweite
  `COPY --from=proto`-Zeile, `protoc`-Aufruf trägt jetzt beide `.proto`-Dateien,
  zweiter `sed`-Import-Fix für `administration_pb2_grpc.py`; zusätzlich ein
  dritter, so nicht geplanter `sed`-Schritt gegen eine interne Kennung, die das
  gRPC-Plugin aus dem `.proto`-Servicekommentar in die generierten
  Docstrings kopiert (real erst beim `make sdk-pack-python`-Lauf gefunden,
  siehe §3-Anmerkung) — real belegt: `make sdk-pack-python` erfolgreich
  (Build+130 Tests+Pack), `make gates` Exit 0.

## 7. Closure-Notiz

- **Was hat funktioniert:** <wird bei Closure gefüllt>
- **Was ging anders als geplant:** <wird bei Closure gefüllt>
- **Steering-Loop-Eintrag:** <wird bei Closure gefüllt, falls einer entsteht>
- **Beobachtungs-Register (`../observations/`):** <wird bei Closure gefüllt>
- **Folge-Slices:** <wird bei Closure gefüllt, falls einer entsteht>
- **Risiken aus §6:** <wird bei Closure gefüllt — siehe §6>
- **Drei Paarungen:** von der Welle-Closure (`welle-sdk-grpc-administration-flaeche`).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** eine Sub-Area, `sdks/python/`
(bestehendes, aktiv gepflegtes Package mit eigenem Release-Workflow).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen.
Treffer unter der 3×-Schwelle, als Risiko in §6 aufgenommen:
`drei-sprachen-kopie-divergiert-am-randfall` (2×),
`sdk-python-untergrenze-ohne-anwender-begruendung` (1×, nur relevant, falls
dieser Slice `requires-python`/das Basis-Image anfasst — laut §1 nicht
geplant). Kein Treffer ≥ 3× offen und einschlägig.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (neue Dateien in
einer bestehenden, bereits Greenfield geführten Sub-Area).
