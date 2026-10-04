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
| U | Versionswechsel bei konstantem Schema ([`ADR-0148`](../../docs/plan/adr/0148-kotlin-sdk-grpc-api-readme-und-upgrade-trigger-erfuellt.md)), läuft zuletzt: der Server 0.5.0 erfasst INSERT, UPDATE und DELETE auf `feed_e2e_full` (volle Replica-Identität, Zeile `id=9101`); vom Datenstand der Quelle über `cdc.changes` (Zeilenzahl und `md5` über `change_id`, `commit_position`, `old_data`, `new_data` in der Reihenfolge von `change_id`) wird eine Momentaufnahme gehalten. Ein `$COMPOSE up -d --force-recreate --no-deps pg-change-feed` mit einer Override-Datei **ohne** `image:`-Zeile ersetzt den Container durch das Image des Dienstes in `compose.yaml` (`:dev`). Geprüft: Container-ID neu, `postgres`/`nats` unberührt, Image-Referenz und Image-ID des neuen Containers gleich dem geladenen Ziel-Image und verschieden vom Start-Image, Health `healthy`, Datenstand gleich der Momentaufnahme, eine danach eingefügte Zeile (`id=9102`) über `cdc.changes` erfasst, Slot da mit nicht kleinerer `confirmed_flush_lsn`, Zeilenzahl der Quelle danach genau eine mehr. Ist die Start-Referenz gleich der Ziel-Referenz (B3), entfällt der Tausch | `ALTSERVER U: Tausch <ID alt> -> <ID neu>, Image <Referenz alt> -> <Referenz neu>, Datenstand vor dem Tausch (<n> Zeilen, Prüfsumme <h>) identisch lesbar, danach eingefügte Zeile erfasst (Position <p>), Phase U <s> s`; bei Start gleich Ziel `ALTSERVER U ÜBERSPRUNGEN: …` |

Das Ziel-Image des Tauschs steht in `compose.yaml`; der Runner liest es dort
(`docker compose config`) und führt kein zweites Image-Literal. Vor jedem Start
prüft er, dass es geladen ist; fehlt es, endet der Lauf mit Exit 1 und der Meldung
„make image vorher“ (kein Überspringschalter). Ein veraltetes `:dev`-Image ist
ein gültiges Ziel, das der Runner nicht prüft.

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
make image                                                                                            # Vorbedingung: das geladene :dev-Image
make test-sdk-altserver
SDK_ALTSERVER_IMAGE=ghcr.io/pt9912/pg-change-feed:dev make test-sdk-altserver                          # B3
SDK_ALTSERVER_WEITER=1 SDK_ALTSERVER_IMAGE=ghcr.io/pt9912/pg-change-feed:dev make test-sdk-altserver   # B3, alle Phasen
```

Host-Werkzeuge: `bash`, `git`, `mktemp`, `awk`, `sed` (ohne `-i`) und `docker`
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
| 0 | B0 bis B2 und U grün: der Server trägt kein `code`, die Eigenschaft bleibt in allen drei Sprachen leer, der Tausch auf das Ziel-Image erhält den Datenstand (bei Start gleich Ziel: U übersprungen, gedruckt) |
| 1 | ein Schritt ist rot (Meldung nennt Schritt und Befund), das Ziel-Image ist nicht geladen („make image vorher“, vor jedem Start), oder ein Container- bzw. Netzname der anderen Runner ist belegt (Abbruch vor jedem Start, aufgeräumt wird nur, was der Runner selbst angelegt hat) |

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
6. **Phase U: Schema konstant, ein Alt-Stand, keine Zeilen während des Tauschs.**
   Das Schema stammt aus dem Arbeitsbaum (Grenze 2); ein Schemawechsel über Versionen
   ist ungemessen ([`ADR-0148`](../../docs/plan/adr/0148-kotlin-sdk-grpc-api-readme-und-upgrade-trigger-erfuellt.md),
   akzeptiertes Negativ). Während des Tauschs läuft kein Schreiber; die Fortsetzung ab
   `confirmed_flush_lsn` trägt der Neustart-Rundlauf von `make test-integration`. Die
   Übertragung der Aussage auf andere Tabellen und Operationen als `feed_e2e_full`
   mit INSERT, UPDATE und DELETE ist *hergeleitet*.

## Test

Das Ziel ist selbst die Messung. Je Zusage die Mutation ihrer Eingabeseite (der
Server) und das gesehene Rot:

| Zusage | mutierte Eingabe | gesehenes Rot |
|---|---|---|
| der Fehlerkörper des Servers vor 0.6.0 trägt kein `code` (B1) | Feed-Image `:dev` (Server mit Meldungscodes) statt 0.5.0 | B1: `ALTSERVER B1 FEHLER: der Fehlerkörper trägt das Feld code`, Exit ≠ 0 |
| die Eigenschaft bleibt leer, in allen drei Sprachen (B2) | Feed-Image `:dev` mit `SDK_ALTSERVER_WEITER=1` | B2 C#, Kotlin und Python rot (Marker `NORMAL_DONE` blieb aus: die Prüfung `MessageCode` leer scheitert vor ihm), Exit ≠ 0 |

Phase U (Mutationen an Kopien des Runners im Scratchpad, Instanz: die PostgreSQL und
der Docker-Daemon des Runners, gesehen am Lauf von `slice-upgrade-versionswechsel-alt-image`):

| Zusage | mutierte Eingabe | gesehenes Rot |
|---|---|---|
| der Datenstand ist nach dem Tausch identisch lesbar | eine Zeile von `cdc.change` wird nach dem Tausch geändert (`UPDATE cdc.change SET new_data = …`, als Superuser, kein Schutz hielt ab) | `U ROT — der Datenstand der Quelle über cdc.changes ist nach dem Tausch nicht identisch`, Exit 1 |
| die Erfassung setzt nach dem Tausch fort | `docker stop` des Feed-Containers nach dem Tausch | `U ROT — die nach dem Tausch eingefügte Zeile (id=9102) erscheint nicht über cdc.changes`, Exit 1 |
| der Tausch wechselt den Build | der Override des Tauschs trägt weiter `image: <SDK_ALTSERVER_IMAGE>` | `U ROT — der neue Container läuft nicht den Ziel-Build`, Exit 1 |

Die Negativprobe B3 (`SDK_ALTSERVER_IMAGE` gleich dem Ziel) bleibt an B1 rot; mit
`SDK_ALTSERVER_WEITER=1` druckt Phase U `ÜBERSPRUNGEN`. Ein Lauf je Mutation; die
Übertragung auf andere Tabellen und Operationen ist *hergeleitet*.

Die Diagnose-Prüfung (`error_code` leer) hat keine Mutation, die sie rot färbt: auch
am `:dev`-Server bleibt sie grün (Absturzfreiheit, keine Aussage über den Server).

Menge der Erprobung: das Image `:dev` als einzige Mutation der Eingabe; ein
Server, der nur das `ErrorInfo` (gRPC) oder nur `error_code` (Diagnose) trägt,
ist nicht gefahren — die Übertragung auf diese Einzelpfade ist *hergeleitet*
aus den Unit-Tabellen der SDKs.
