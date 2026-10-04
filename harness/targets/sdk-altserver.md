# `make test-sdk-altserver` — SDK 0.6.0 gegen einen Server vor 0.6.0

## Vertrag

`make test-sdk-altserver` fährt die Fehlerfälle der SDK-Packages (C#, Kotlin,
Python, Quelle des Arbeitsstands) gegen einen Server, der noch keine
Meldungscodes sendet, und misst die Aussage von
[`ADR-0145`](../../docs/plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md)
„der Code gilt nur dort, wo der Server einen setzt; ein Server vor 0.6.0
liefert leer“: die Eigenschaft bleibt leer, Fehlertext und Fehlertyp bleiben,
nichts stürzt ab. Der Server vor 0.6.0 ist das veröffentlichte Image
`ghcr.io/pt9912/pg-change-feed:0.5.0` in der Form `<Tag>@<Index-Digest>`; der
Tag bleibt in der Referenz, damit `make pin-stale-all` den Digest gegen den
Tag `0.5.0` vergleicht
([`ADR-0146`](../../docs/plan/adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)).
Der Digest steht an genau einer Stelle, dem Standardwert von
`SDK_ALTSERVER_IMAGE` in `tools/harness/run-sdk-altserver-tests.sh`.

Der Runner fährt dieselbe Compose-Umgebung wie die drei SDK-Realserver-Runner
([`ADR-0110`](../../docs/plan/adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
Festlegung 2): PostgreSQL und NATS aus `compose.yaml`, Schema-Rollout des
**Arbeitsbaums** über `make schema-rollout`; nur der Feed-Container ist das
Altserver-Image. Das Image kommt über eine Override-Datei im Temp-Verzeichnis
(`docker compose -f compose.yaml -f <Temp>/override.yaml`); `compose.yaml` und
die drei bestehenden Runner bleiben unverändert. Schritte:

| Schritt | Messung | Gedruckt |
|---|---|---|
| B0 | Der Server vor 0.6.0 läuft gegen das aktuelle Schema: Slot da, Health `healthy`, eine eingefügte Zeile einer aktivierten Tabelle liegt über `cdc.changes` | `ALTSERVER B0: healthy, Slot <slot>, change_id=<id>, Feed-Image <Referenz>` |
| B1 | Rohdraht-Probe der HTTP-Seite ohne SDK (`tools/harness/sdk-kompat/altserver/rohdraht.py`, per `python -` im Python-Image des SDK im Compose-Netz): die Aktivierung einer fehlenden Tabelle endet mit Status 404, der Körper trägt kein Feld `code` | `ALTSERVER B1: BODY <roher Körper>`, `STATUS 404`, `TEXT <Fehlertext>` |
| B2 | Je Sprache eine Test-Klasse (`ErrorCodeAltServerTests` C#, `ErrorCodeAltServerTest` Kotlin, `test_error_code_altserver.py` Python) im bestehenden Integrations-Quellsatz: HTTP und gRPC enden mit dem typisierten `NotFound`-Fehler, der Meldungscode ist leer, der Fehlertext ist nicht leer (HTTP: gleich dem Text aus B1); ein Reader-Token endet mit 403 bzw. `PermissionDenied` ohne Code; die Diagnose (gRPC) liefert den Bericht ohne Absturz mit leerem `error_code` — im Normalbetrieb und in einem Fehlerzustand (`error_class = 'schema'` ohne `error_code`), den der Runner nach dem Marker `NORMAL_DONE` in `cdc.process_heartbeat` schreibt (der periodische Heartbeat setzt ihn zurück, der Runner schreibt wiederholt, bis der Test endet) | `ALTSERVER B2 <sprache>: RECEIVED code=none http=404 grpc=NOT_FOUND text=<n> diag_error_code=leer (Exit 0)` |
| B3 | Negativprobe: `SDK_ALTSERVER_IMAGE=ghcr.io/pt9912/pg-change-feed:dev` (der geladene 0.6.0-Bau, `make image` vorher) färbt den Lauf an B1 rot | `ALTSERVER B1 FEHLER: der Fehlerkörper trägt das Feld code …`, Exit ≠ 0 |

Die gRPC-Seite hat keine Rohdraht-Probe: dass der Server vor 0.6.0 kein
`ErrorInfo` setzt, ist aus dem Quelltext abgeleitet (`git grep` auf `ErrorInfo`
am Tag `v0.5.0` unter `internal` und `cmd`). Gemessen sind drei Aussagen, mit
unterschiedlicher Trägerschaft: HTTP (B1 am Rohdraht, B2 am SDK) und der gRPC-Fehler
`NotFound` (B2; die Probe, die alle `message_code`-Prüfungen bis auf die gRPC-Prüfung
entfernt, färbt B2 am `:dev`-Server rot) tragen als *gemessen*. Die Diagnose trägt nur
als „das SDK stürzt bei leerem `error_code` nicht ab“ (*gemessen*); dass der Server vor
0.6.0 den `error_code` nicht setzt, ist *hergeleitet* (Schema- und Quelltext-Diff): der
Runner schreibt den Fehlerzustand selbst ohne `error_code`, und der Normalzustand trägt auch
am `:dev`-Server keinen, die Prüfung der Diagnose kann an keinem Server rot werden, den der
Runner so fährt (Probe des Verifiers, nur in Python gefahren, für C# und Kotlin *hergeleitet*).

Der Runner schreibt nichts in den Arbeitsbaum, insbesondere nicht den
Abdeckungs-Träger `docs/user/sdk-e2e-abdeckung.md` (`git status --porcelain` ist
vor und nach dem Lauf gleich).

## Aufruf

```text
make test-sdk-altserver
SDK_ALTSERVER_IMAGE=ghcr.io/pt9912/pg-change-feed:dev make test-sdk-altserver                          # B3
SDK_ALTSERVER_WEITER=1 SDK_ALTSERVER_IMAGE=ghcr.io/pt9912/pg-change-feed:dev make test-sdk-altserver   # B3, alle Phasen
```

Host-Werkzeuge: `bash`, `git`, `mktemp`, `sed` (ohne `-i`) und `docker`
([`AGENTS.md`](../../AGENTS.md) §3.1). Alles andere (`dotnet`, `gradle`,
`python`) läuft im Container.

## Overrides

| Variable | Bedeutung | Standard |
|---|---|---|
| `SDK_ALTSERVER_IMAGE` | Image des Feed-Containers | `ghcr.io/pt9912/pg-change-feed:0.5.0@<Index-Digest>` (Standardwert im Runner) |
| `SDK_ALTSERVER_WEITER` | `1`: der Runner läuft nach einem roten Schritt weiter und endet am Schluss mit Exit 1 (Diagnose der Negativprobe: zeigt, dass auch die Test-Klassen von B2 an einem Server mit Meldungscodes rot werden) | `0` (Abbruch beim ersten roten Schritt) |

## Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | B0 bis B2 grün: der Server trägt kein `code`, die Eigenschaft bleibt in allen drei Sprachen leer |
| 1 | ein Schritt ist rot (Meldung nennt Schritt und Befund), oder ein Container- bzw. Netzname der anderen Runner ist belegt (Abbruch vor jedem Start, aufgeräumt wird nur, was der Runner selbst angelegt hat) |

Über `make` kommt jeder Ausgang ≠ 0 als 2 an. Kein Gate: der Lauf braucht DB-Zugang,
Docker-Pulls und Netz; das Ziel steht in keinem Gate-Bündel (`make gates`,
[`AGENTS.md`](../../AGENTS.md) §3.6, §4).

## Wer es aufruft

- **Implementer, Reviewer, Verifier** einer SDK-Version, die Fehlertypen oder die
  Diagnose berührt; Reviewer und Verifier als Probe der Aussage „ein Server vor
  0.6.0 liefert leer“.

## Grenze

1. **Ein Server-Stand.** Gemessen ist 0.5.0 als der letzte Stand vor den
   Meldungscodes; ein früherer Server und früher gebaute SDKs sind nicht Gegenstand
   (ein weiterer Stand ist ein Wert von `SDK_ALTSERVER_IMAGE`, kein neuer Runner).
2. **Schema des Arbeitsbaums.** Der Rollout stammt aus dem Arbeitsbaum, nicht aus
   dem Stand 0.5.0: das Paar „Server 0.5.0, Schema 0.6.0“ ist die Rollout-Reihenfolge
   „Schema vor Server“, nicht ein vollständig alter Betrieb.
3. **Eine Fläche je Aussage.** Die Stream-Clients und NATS tragen keine Fehlertypen mit
   Code ([`ADR-0145`](../../docs/plan/adr/0145-sdk-meldungscodes-eigenschaft-der-fehlertypen.md)
   Festlegung 5) und stehen nicht in diesem Lauf; die Gegenseite (SDK liest den Code
   am Server 0.6.0) trägt die Fehlercode-Phase der drei Realserver-Runner.
4. **Namen und Netz der anderen Runner.** Der Runner teilt Container-Namen
   (`cdc-test-postgres`, `cdc-test-nats`, `cdc-test-feed`) und das Netz
   `cdc-feed-test` mit `make test-integration` und den SDK-Realserver-Runnern; ein
   gleichzeitiger Lauf bricht mit einer Meldung ab.
5. **Registry.** Die Pulls von `postgres` und `nats` (Docker Hub) und des
   Altserver-Images (ghcr) müssen möglich sein; ein Abruflimit endet als Docker-Fehler,
   nicht als grüner Lauf.

## Test

Das Ziel ist selbst die Messung. Je Zusage die Mutation ihrer Eingabeseite (der
Server) und das gesehene Rot:

| Zusage | mutierte Eingabe | gesehenes Rot |
|---|---|---|
| der Fehlerkörper des Servers vor 0.6.0 trägt kein `code` (B1) | Feed-Image `:dev` (Server mit Meldungscodes) statt 0.5.0 | B1: `ALTSERVER B1 FEHLER: der Fehlerkörper trägt das Feld code`, Exit ≠ 0 |
| die Eigenschaft bleibt leer, in allen drei Sprachen (B2) | Feed-Image `:dev` mit `SDK_ALTSERVER_WEITER=1` | B2 C#, Kotlin und Python rot (Marker `NORMAL_DONE` blieb aus: die Prüfung `MessageCode` leer scheitert vor ihm), Exit ≠ 0 |

Die Diagnose-Prüfung (`error_code` leer) hat keine Mutation, die sie rot färbt: auch
am `:dev`-Server bleibt sie grün (Absturzfreiheit, keine Aussage über den Server).

Menge der Erprobung: das Image `:dev` als einzige Mutation der Eingabe; ein
Server, der nur das `ErrorInfo` (gRPC) oder nur `error_code` (Diagnose) trägt,
ist nicht gefahren — die Übertragung auf diese Einzelpfade ist *hergeleitet*
aus den Unit-Tabellen der SDKs.
