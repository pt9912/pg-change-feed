# ADR-0087: Beispiel-Clients in C# und Kotlin — Sprach-Wurzel, eigene Werkzeugkette, Spec-Anker

**Status:** Accepted — Supersedes [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
in **einer** Klausel: der Gestalt-Aussage ihrer §Entscheidung Festlegung 1
(„**Ein Verzeichnis je Protokoll, je ein `main`-Paket:** `examples/http-client/`,
`examples/sse-client/`, `examples/grpc-client/`"), soweit sie `examples/` als
**flache** Protokoll-Ablage festlegt. Diese Klausel bleibt für die
**Go**-Clients in Kraft; für Sprachen, deren Werkzeugkette ein
Projektverzeichnis verlangt, tritt das Sprach-Wurzelverzeichnis daneben
(Festlegung 2). Alles Übrige der [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
bleibt **hiermit bestätigt** und wird nicht wiederholt: die Import-Grenze und
ihr Träger (§2), der gRPC-Baustein (§3), die `.a-check.yml`-Gruppen und Kanten
(§4), die Bindung an den `make test`-Pfad (§5), das Verhältnis zu den
Wegwerf-Clients (§6), das Handbuch-Zitat (§7) und „Was diese ADR nicht ändert".
Ebenfalls unberührt: [`ADR-0079`](0079-nats-beispielclient-vierter-examples-client.md)
(der vierte Go-Client) — die Zahl der **Go**-Clients ändert sich nicht.

**Datum:** 2026-09-17

**Autor:** pt9912 (Architect-Rolle, unabhängiger Architect-Zug auf
Nutzerentscheidung vom 2026-09-17; anderer Kontext als die Implementer-Läufe
der drei Go-Beispiel-Clients und als der Planner-Lauf, der den
`examples`-Nachzug führte — Modul 8 §Rollen-Regeln: „Architect schreibt")

**Bezug:** [`LH-FA-SST-006`](../../../spec/lastenheft.md) (HTTP-/JSON-API — der
Draht des HTTP-Beispiels), [`LH-FA-SST-008`](../../../spec/lastenheft.md)
(Live-Change-Stream — der Draht des SSE-Beispiels),
[`ARC-005`](../../../spec/architecture.md) (Driving Adapters),
[`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
(in der Gestalt-Klausel superseded; sonst bestätigt),
[`ADR-0079`](0079-nats-beispielclient-vierter-examples-client.md) (der vierte
Go-Client; unberührt), [`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md)
(die Wegwerf-Clients und die Gate-Scope-Gruppen),
[`ADR-0041`](0041-a-check-maschinenform-architekturpruefung.md) (Maschinenform
der §2-Constraints), [`ADR-0044`](0044-image-beleg-semantik.md)
(Image-Beleg-Semantik), [`ADR-0085`](0085-build-kontext-ausnahme-test-only-zweck.md)
(die Build-Kontext-Ausnahme-Klasse — von der Bauform berührt),
[`SPEC-023`](../../../spec/pflichtenheft.md) (Beispiel-Clients — von dieser ADR
**geschärft**), [`AGENTS.md`](../../../AGENTS.md) §3.1/§3.8/§3.9/§3.12,
`.a-check.yml`, `a-check.mk`, `Makefile`, `Dockerfile`, `.dockerignore`,
`harness/README.md` §Sensors/§Werkzeuge, `.github/workflows/e2e.yml`,
`examples/**`, `docs/user/benutzerhandbuch.md` (§4 Zugriffs-Abschnitte,
§5 Konfiguration), `proto/cdc/stream/v1/changestream.proto`

**Schärft:** [`SPEC-023`](../../../spec/pflichtenheft.md) — die ADR macht die
Ziel-Form der Beispiel-Clients verbindlich: Sprach-Wurzel, Werkzeugketten-Form,
Bau-/Prüfweg und die Grenze der `examples`-Gruppe. Das Pflichtenheft trägt die
Festlegung (Rang 2), diese ADR die Entscheidung und ihre Begründung; bei
Widerspruch gilt der höhere Rang.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

**Die Nutzerentscheidung.** Der Auftraggeber hat am 2026-09-17 entschieden:
„Die Client-Beispiele brauchen wir noch in zwei weiteren Programmiersprachen:
**C#** und **Kotlin**." Dazu die Nutzeranweisung, dass das **Pflichtenheft**
die Beispiel-Clients tragen muss.

**Der Ist-Stand — gemessen, nicht erinnert.**

- `examples/` auf der Repo-Wurzel trägt **drei** Go-Clients:
  `examples/http-client/`, `examples/sse-client/`, `examples/nats-client/`
  (je ein `main`-Paket, `go run ./examples/<name>`). Der vierte benannte
  Go-Client, `examples/grpc-client/`, **existiert nicht**: sein Import des
  internen Stubs ist `wrong-direction`, und die öffentliche Bindungsfläche
  `gen/cdc/stream/v1` **existiert nicht** (`no required module provides
  package …`) — beide Wege real gemessen im Zug, der den `examples`-Nachzug
  zurückgeführt hat.
- `spec/pflichtenheft.md` führt `SPEC-001`…`SPEC-022`; **kein** `SPEC` nennt
  `examples/` oder eines der Client-Programme. Die drei **existierenden**
  Go-Clients sind damit unverankert: entschieden von
  [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)/[`ADR-0079`](0079-nats-beispielclient-vierter-examples-client.md),
  gebaut von zwei Slices, aber die Spec weiß nichts von ihnen.
- `.a-check.yml` führt die Gruppe `examples: ["examples/**"]` mit **keiner**
  Kante und dem Kommentar, ihre Clients benutzten „ausschließlich die
  Standardbibliothek und öffentliche Fremdmodule".
- `.dockerignore` ist eine geschlossene Liste (`*` + Negationen für `cmd/`,
  `internal/`, `go.mod`, `go.sum` und zwei `tools/`-Dateien): `examples/**`
  liegt **nicht** im Build-Kontext; `make image` baut ohne `--target` die
  letzte Stufe (`runtime`, distroless) und damit nur deren Vorfahren.
- `docs/user/benutzerhandbuch.md` §5.1 führt alle `CDC_*`-Variablen; die
  Zugriffs-Abschnitte nennen je Oberfläche ein Go-Beispiel.

**Was [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
entschieden hat — und was nicht.** Sie entscheidet `examples/` als Ort
öffentlicher Beispiel-Clients, ihre Import-Grenze, ihren nicht vorhandenen
Lauf-Beleg und ihr Handbuch-Zitat. **Jede ihrer Begründungen ist
Go-förmig:** sie argumentiert mit dem Verbot des `internal/`-Imports (eine
Go-Sprachregel über Import-Pfade), mit dem **Go**-Binding aus der `.proto`
(`option go_package`) und mit `go run ./examples/<name>` als Startform
(Festlegung 1). Über **andere Sprachen sagt sie nichts.** Die Frage, die
dieser Zug entscheidet, ist deshalb nicht „gilt `ADR-0076`?", sondern
**„was von ihr gilt fremdsprachig, und was nicht?"**

**Fünf Bindungen prägen den Lösungsraum — alle fünf gemessen.**

1. **Docker-only** (`AGENTS.md` §3.1): der Host hat Docker und GNU `make`.
   Kein Host-`dotnet`, kein Host-`gradle`, kein Host-`kotlinc`. Jede Sprache
   braucht ein **gepinntes Basis-Image**; die Digest-Pinnung ist gelebte
   Haus-Kultur (jede der fünf `Dockerfile`-Stufen ist digest-gepinnt).
2. **Der Anti-Verrottungs-Träger der Go-Beispiele trägt fremdsprachig nicht.**
   `ADR-0076` Festlegung 5 lehnt ein eigenes Gate ab, *weil* `make test`
   (`go test -race ./...`) die Beispiele ohnehin übersetzt und ihre Tests
   fährt. Für C#/Kotlin existiert **kein** solcher Pfad: es gibt heute keinen
   Lauf, der sie übersetzt.
3. **Die `.a-check.yml`-Aussage ist für Nicht-Go-Quellen nicht getragen.**
   Real gemessen in diesem Zug (a-check-Binary, Wegwerf-Baum):

   | Probe | Ergebnis |
   |---|---|
   | `.cs`/`.kt` unter `examples/**`, `languages` führt nur `go` | **0 Befunde**, kein Hinweis — die Dateien werden nicht gelesen |
   | `languages: {csharp: ["**/*.cs"]}` aufgenommen, ein `.cs`-`using` vorhanden | Datei wird gelesen, das Symbol **löst auf keine Schicht auf**; a-check meldet einen Hinweis auf stderr („0 von 1 Import-Symbolen lösen auf eine Schicht auf") — **kein** Befund, Exit 0 |
   | `using Cdc.Domain;` gegen einen `resolution`-Block (`mode: fixed-root`, `package_base`, `roots`) | Symbol bleibt in dieser Probe unaufgelöst; die Konfigurationsform für C# ist mit diesem Zug **nicht** belegt |
   | unbekannter Sprach-Schlüssel (z. B. `cs:`) | **Exit 2** („unbekannte Sprache"), gültig sind `cpp|csharp|go|java|kotlin|python|rust|typescript` |

   Heißt: die Gruppe kann für Nicht-Go-Quellen heute **keine Kante tragen**,
   und der Kommentar über ihr liest sich trotzdem wie ein geprüfter Beleg.
   Das ist eine `AGENTS.md` §3.12-Frage (eine Aussage, die als Beleg gelesen
   wird), kein Stil.
4. **Die Bauform ist nicht beliebig.** Ebenfalls real gemessen:
   ein `docker build` **ohne** `--target` baut eine Stufe, die nicht Vorfahre
   der letzten ist, **nicht** (Probe: eine `unused`-Stufe zwischen `build` und
   `runtime` erschien nicht im Bauplan); und ein Build-Kontext, der ein
   **Unterverzeichnis** ist, unterliegt der `.dockerignore` der Wurzel
   **nicht** (Probe: Kontext `ctxsub/cs` kopierte `Program.cs`, Kontext
   `ctxsub` mit `*`-Ignore kopierte nichts).
5. **Die verkörperte Handbuch-Regel.** Wächst eine benutzer-sichtbare
   Oberfläche, ohne dass das Handbuch im selben Zug mitgezogen wird, entsteht
   eine Lücke, die kein Sensor fängt.

**Keiner der vier Re-Evaluierungs-Trigger `ADR-0076`s tritt ein** — und das
ist für den Zuschnitt dieser Entscheidung entscheidend. Insbesondere
**nicht** Trigger 2 („Der öffentliche Vertrag bekommt einen weiteren
Konsumenten **außerhalb** von `adapters`, `tooling` und `examples`"): die
C#- und Kotlin-Clients liegen **in** `examples`, und der gewählte Zuschnitt
erreicht `gen/**` gar nicht. Hätte er es getan — ein gRPC-Beispiel, das seine
Stubs aus der `.proto` erzeugt —, wäre genau dieser Trigger berührt und die
`contract`-Kanten wären per Folge-ADR neu zu ziehen. Deshalb schneidet diese
ADR den gRPC-Client **aus** (Festlegung 1) und benennt ihn als eigene
Entscheidung statt ihn mitzunehmen.

## Entscheidung

Wir wählen: **`examples/` führt die Beispiel-Clients jeder Sprache in der
HTTP-Familie — HTTP-/JSON und SSE —; jede Sprache bekommt ein eigenes
Sprach-Wurzelverzeichnis, ein eigenes digest-gepinntes Werkzeugketten-Image
und ein eigenes Bau-/Testziel, das Werkzeug bleibt und kein Gate wird; die
Go-Clients behalten Form und Ort, die Nicht-Go-Quellen bekommen einen
Spec-Anker, und die `.a-check.yml`-Aussage wird auf Go beschränkt, mit
benannter Grenze.**

Sieben Festlegungen:

### 1 — Umfang: die HTTP-Familie je Sprache, zwei Clients je Sprache

- **Aufgenommen:** je Sprache ein **HTTP-/JSON-Client** (`GET /tables` und
  `GET /changes` mit dem `reader`-Token) und ein **SSE-Client**
  (`GET /changes/stream`, jedes Event ausgegeben) — zusammen vier Programme:
  `examples/csharp/http-client/`, `examples/csharp/sse-client/`,
  `examples/kotlin/http-client/`, `examples/kotlin/sse-client/`.
- **Warum die HTTP-Familie und nicht alle vier Oberflächen.** Drei Gründe,
  jeder prüfbar:
  1. **Sie teilen die Werkzeugkette vollständig.** Beide laufen auf
     `CDC_HTTP_ADDR`, benutzen denselben Transport und dieselbe Token-Klasse;
     der marginale Preis des zweiten Clients einer Sprache ist eine Datei, ein
     Handbuch-Absatz und ein Test — die teure Hälfte (Basis-Image,
     Abhängigkeits-Pin, Bauziel) fällt **je Sprache** an, nicht je Client.
  2. **Sie brauchen keinen Codegenerator.** Damit kommt die Entscheidung
     ohne die Frage aus, wo erzeugte Fremdsprachen-Stubs liegen, ob sie
     committet werden und was die `contract`-Gruppe dann bedeutet — genau
     die Frage, die `ADR-0076` Trigger 2 aufwirft (Kontext).
  3. **Sie sind die Oberflächen, die jede Integration braucht.** Die zwei
     Streams und das Wecksignal sind optionale Features, die per
     Umgebungsvariable aktiviert werden; die HTTP-Familie ist es nicht.
- **Nicht aufgenommen, je mit Adresse für das Wachstum:**
  - **gRPC je Sprache.** Er braucht protoc samt Sprach-Plugin und eine
    Festlegung für die erzeugten Stubs (Ort, committet oder im Bau erzeugt,
    Verhältnis zu `gen/**` und zur `contract`-Gruppe). Das ist eine eigene
    Entscheidung mit eigener Beweislast und wird hier **nicht** mitgenommen.
  - **NATS je Sprache.** Er braucht je Sprache eine Fremdsprachen-Client-
    Bibliothek und zeigt den zweiseitigen Ablauf; sein Vorbild-Wert ist
    zweitrangig gegenüber der Basis-Oberfläche und sein Preis je Sprache
    höher (eine zweite Abhängigkeit).
  - **Die Go-Clients.** Sie bleiben unverändert; insbesondere bleibt
    `examples/grpc-client/` der offene Go-Vorgänger. Diese ADR schneidet
    nichts um.
- **Was der Nutzerwunsch damit nicht bekommt — ausgesprochen:** in C# und
  Kotlin gibt es **kein** gRPC- und **kein** NATS-Beispiel, und in allen drei
  Sprachen gibt es weiterhin **kein** gRPC-Beispiel. Wer das ändert, ändert
  den Umfang per Folgeschritt, nicht per Auslegung.

### 2 — Ort und Sprach-Wurzel (die abgelöste Klausel)

- **`examples/` bleibt die Wurzel**; je Sprache entsteht ein
  **Sprach-Wurzelverzeichnis**: `examples/csharp/` und `examples/kotlin/`.
  Darin liegt je Client ein eigenes Verzeichnis mit **einem** Programm und
  **einem** Einstiegspunkt.
- **Die Sprach-Wurzel existiert, wo die Werkzeugkette ein Projektverzeichnis
  verlangt** — C#: eine Projektdatei je Client, Kotlin: ein
  Werkzeugketten-Projekt im Sprach-Wurzelverzeichnis mit einem Client je
  Modul/Einstiegsklasse. Für **Go** verlangt die Werkzeugkette das nicht:
  die drei bestehenden Clients bleiben **flach** (`examples/<name>/`), ihre
  Startform `go run ./examples/<name>` bleibt gültig, und die Handbuch-Zeilen
  der drei werden **nicht** angefasst. Ein Umzug der Go-Clients wäre
  Arbeit am Bestand ohne Adresse (die Sprach-Wurzel für Go trüge keine
  Aussage, die `go run` nicht schon trägt).
- **Die Namen tragen weiter `-client`** (`ADR-0076` Festlegung 1); das Wort
  *Consumer* bleibt dem Domänenbegriff (`cdc.consumer`) vorbehalten.
- **Genau diese Gestalt-Aussage ist die superseded Klausel** (§Status): was
  `ADR-0076` als *flache* Protokoll-Ablage festlegt, gilt für die Go-Clients;
  für die anderen Sprachen tritt die Sprach-Wurzel daneben.
- **Ein Verzeichnis mit Beispiel-Quellen ist keine Doku-Ablage.** Die
  Beispiele bleiben ausführbarer Code; ein Beispiel ohne seine Handbuch-Zeile
  ist undokumentiert, nicht „später dokumentiert" (Kontext, Bindung 5).

### 3 — Werkzeugkette und Bauform

- **Je Sprache ein eigenes Dockerfile im Sprach-Wurzelverzeichnis**
  (`examples/csharp/Dockerfile`, `examples/kotlin/Dockerfile`), gebaut mit
  dem **Sprach-Wurzelverzeichnis als Kontext**. Das ist die Form, die
  drei Dinge zugleich leistet (alle drei gemessen, Kontext Bindung 4):
  1. **Die Wurzel-`Dockerfile` bleibt unberührt** — kein fremder
     Werkzeugketten-Block in der Datei, die das Auslieferungs-Image baut.
  2. **`.dockerignore` bleibt unberührt.** Ein Kontext, der ein
     Unterverzeichnis ist, unterliegt der Wurzel-`.dockerignore` nicht; die
     Beispiele liegen daher in ihrem eigenen Kontext, ohne dass die zwei
     `!`-Zeilen der Wurzel um ein **Verzeichnis** wachsen müssten. Eine
     Verzeichnis-Negation wäre an der Klasse der Build-Kontext-Ausnahme
     ([`ADR-0085`](0085-build-kontext-ausnahme-test-only-zweck.md), Merkmal
     (ii): „genau eine Datei") vorbei, und sie vergrößerte den Kontext
     **jeder** Stufe (`COPY . .` in `coverage` und `build`) — die
     `examples/**`-Änderung veränderte dann das Image-Digest-Verhalten des
     ausgelieferten Baus.
  3. **`make image` bleibt unverändert.** Die letzte Stufe bleibt `runtime`;
     ein `docker build` ohne `--target` baut fremde Werkzeugketten-Stufen
     ohnehin nur, wenn sie Vorfahren der Zielstufe sind — mit eigenen
     Dockerfiles ist die Frage gar nicht gestellt. Der Runtime-Graph des
     Images, sein Digest-Beleg ([`ADR-0044`](0044-image-beleg-semantik.md))
     und der Bau-Kontext des Go-Baus bleiben, was sie sind.
- **Jede Basis ist digest-gepinnt**, jede Abhängigkeit versionsgepinnt. Die
  Pinnung ist Form, kein Detail: sie folgt derselben Haus-Kultur wie die
  fünf `Dockerfile`-Stufen. **Gemessen am 2026-09-17** (Befehl:
  `docker manifest inspect <image>`, amd64-linux-Manifest-Digest) als
  Kandidaten — die **Wahl** und der Pin gehören dem umsetzenden Zug, der
  Pin-Commit ist ein bewusster Commit:

  | Sprache | Kandidat (Tag) | gemessener amd64-Digest |
  |---|---|---|
  | C# | `mcr.microsoft.com/dotnet/sdk:10.0` | `sha256:60a2b2230a0d052bc54c0d453e97e331219ed503c5411470ee226859420e693c` |
  | C# (Alternative) | `mcr.microsoft.com/dotnet/sdk:9.0` | `sha256:a19efb5ab1d4017030ff791d463825f9a64009a30c522a4239947d0431fc9589` |
  | Kotlin | `eclipse-temurin:21-jdk` | `sha256:085eb93e049c7397f725bd8be31c4f52ba4777a75168e851428508234aa224e` |
  | Kotlin (Alternative) | `eclipse-temurin:25-jdk` | `sha256:6fc323be479ffa5874d1d686f4fc49aa0e11d415c79d2257592ff3d6937221eb` |

  Die Kotlin-Werkzeugkette (Compiler/Projekt-Werkzeug, JSON-Zugriff) ist
  **kein** Digest, sondern eine Versions- und Prüfsummen-Pinnung im
  Sprach-Wurzelverzeichnis; die Kotlin-Standardbibliothek führt **keinen**
  JSON-Parser, C# dagegen liest JSON aus dem SDK-Rahmen — der Zugriff selbst
  (öffentliches Fremdmodul oder Rahmen) ist Detail des umsetzenden Zuges,
  die Zulässigkeit folgt Festlegung 5.
- **Das Image ist startbar.** Startform ist ein Container-Aufruf, nicht
  `dotnet run`/`gradle run` auf dem Host — dieselbe „eine Befehlsform, innen
  wie außen"-Regel, die `ADR-0076` Festlegung 1 für `go run` gewählt hat.
  Ob ein Image je Sprache oder je Client entsteht, ist Form-Detail des
  umsetzenden Zuges; **dass** kein Host-Werkzeug nötig ist, ist Vorgabe.
- **Erwartet, nicht gemessen:** Bau-Zeit und Image-Größe. Beide hängen an
  der Werkzeugkette und am Cache; eine Zahl hier wäre eine Schätzung
  (`AGENTS.md` §3.12). Sie zu messen ist Sache des umsetzenden Zuges, der sie
  in seiner Closure-Notiz führt.

### 4 — Laufen und Prüfen: ein Werkzeug je Sprache, kein Gate

- **Je Sprache ein `make`-Ziel** (Arbeitsname `examples-csharp`,
  `examples-kotlin`), das die Beispiele der Sprache **übersetzt und ihre
  netzlos prüfbaren Tests fährt**. Das ist `ADR-0079` Festlegung 3 auf einen
  neuen Träger angewandt: nicht nur Kompilat, sondern Ausführung der reinen
  Funktionen (Aufbau der Anfrage, Zerlegen eines SSE-Frames) — die
  nicht-offensichtliche Logik bekommt einen Sensor, der netzlos läuft.
- **Kein Gate.** Die Begründung ist strukturell und aus dem Repo gemessen:
  **jedes** Target in `GATE_CHECKS` läuft netzlos (`docker run --rm
  --network none`, in `d-check.mk`, `a-check.mk`, `harness/mk/*.mk`), und
  `make gates` ist das Pflicht-Bündel vor jedem PR. Ein C#-/Kotlin-Bau
  braucht Paket-Bezug (NuGet/Maven) und damit Netz; er wäre das **erste**
  netzgebundene Gate und würde das Bündel strukturell verändern. Die Ziele
  sind **Werkzeuge**: sie hängen nicht an `GATE_CHECKS`, und `make gates`
  bleibt, was es ist. Ein Gate wäre überdies eine Schwellen-/Bindungs-Frage
  eigener Beweislast — diese ADR will keines (siehe §Verglichene
  Alternativen, Option C).
- **Kein Compose-Lauf.** `ADR-0076` Festlegung 5 lehnt einen eigenen realen
  Lauf gegen die Compose-Umgebung als **Duplikation** ab (der Draht ist
  einmal real belegt); das gilt fremdsprachig unverändert. Die Beispiele
  sind **Doku mit Bindungs-Pflicht**, kein Lauf-Beleg: sie erscheinen nicht
  in `docs/user/e2e-abdeckung.md`.
- **Der Carrier, ausgesprochen statt angenommen.** Ein Ziel, das niemand
  fährt, ist kein Carrier. Deshalb bekommt die Sprache einen **eigenen,
  nicht-blockierenden Workflow** (Muster `.github/workflows/e2e.yml`: sichtbar
  im PR, blockiert keinen Merge) — **ein** Workflow über beide Ziele, nicht
  zwei. Er ist ein **eigener Schritt** (Festlegung: der Workflow steht
  allein), weil ein neuer Workflow nach `AGENTS.md` §3.10 erst nach einem
  realen, grünen Post-Push-Lauf als abgeschlossen gilt: das betroffene
  Risiko bleibt bis dahin **offen**.
- **Der Aufruf-Disziplin gilt §3.9**: der Exit-Code des Ziels wird direkt
  gelesen, nie durch eine Pipe; die abhängige Folgehandlung wartet auf ihn.
- **Die Grenze, benannt:** „übersetzt und die reinen Funktionen laufen" ist
  nicht „der Client holt am laufenden Feed eine Änderung". Der Rest-Fall ist
  derselbe wie `ADR-0076` Festlegung 5 und liegt beim Review, bis
  `ADR-0076` §Re-Evaluierungs-Trigger 3 ihn schließt.

### 5 — Was die Quellen benutzen dürfen — und was die `examples`-Gruppe trägt

- **Die Zulässigkeitsregel `ADR-0076`s gilt wörtlich:** ausschließlich die
  Standardbibliothek/Runtime der Sprache und **öffentliche** Fremdmodule;
  kein Import eines privaten Baums dieses Repos; die Quelle muss außerhalb
  dieses Repos kopierbar und übersetzbar bleiben. Für ein Beispiel, das nur
  die HTTP-Familie benutzt, ist das erreichbar: es liest JSON, baut URLs und
  liest Zeilen — alles Rahmen bzw. öffentliches Modul.
- **Die `.a-check.yml`-Aussage wird auf Go beschränkt, und die Lücke wird
  benannt statt still.** Der Kommentar der `examples`-Gruppe sagt künftig,
  was er trägt: eine **Go**-Aussage. Der Grund ist gemessen (Kontext,
  Bindung 3): a-check liest Nicht-Go-Quellen nur, wenn ihre Sprache
  deklariert ist; ohne `resolution`-Block löst ein `using`/`import` auf
  keine Schicht auf, und es gibt in diesem Repo **keine** C#-/Kotlin-Schicht,
  gegen die etwas verstoßen könnte. Die zwei ehrlichen Wege sind:
  den Sprach-Schlüssel **setzen** (dann liest a-check die Dateien und meldet
  für jede Zieldatei einen Auflösungs-Hinweis — Lärm ohne Kante), oder die
  Aussage **auf Go beschränken** und die Grenze nennen. Wir wählen das
  Zweite; **der Wächter dort ist das Review**, und ein Beleg ist lesbar:
  die gepinnten Abhängigkeitsmanifeste der Sprache (Projektdatei bzw.
  Versionskatalog) tragen die Abhängigkeitsliste — „nur öffentliche Module"
  ist damit **nachlesbar**, nicht bloß behauptet.
- **Der Tag, an dem die Grenze kippt, ist benannt:** mit dem gRPC-Client
  entsteht ein `gen/**`, das auch fremdsprachig gelesen wird. Dann ist die
  Sprach-Deklaration samt `resolution` **mit** der `contract`-Frage zu
  entscheiden (Trigger 2), nicht vorher auf Vorrat.
- **Der Spec-Anker.** `spec/pflichtenheft.md` bekommt `SPEC-023`
  (Beispiel-Clients) und trägt damit **alle** Clients der Sprache Go, C# und
  Kotlin — die Nachverankerung der drei bestehenden Go-Clients ist Teil
  dieses Zuges, nicht ein Folgeschritt: ein Eintrag nur für die zwei neuen
  Sprachen ließe die drei bestehenden ohne Anker und führte den Leser der
  Spec zu den Sprachen, die es noch nicht gibt, statt zu denen, die es gibt.
- **Die Beispiele lesen keine Konfigurationsdatei.** Sie beziehen Adresse
  und Token aus denselben Umgebungsvariablen, die das Handbuch führt, und
  lassen sie per Flag übersteuern — **kein** `CDC_CONFIG_FILE`. Der Grund
  steht in der Sache: `CDC_CONFIG_FILE` ist der Deployment-Eingang des
  **Feeds**, nicht der eines Integrator-Programms; eine dritte Herkunft
  verdreifachte die Lese-Regel jedes Beispiels und lehrte eine Kopplung, die
  es für den Leser nicht gibt. Die Aussage steht als Zeile in `SPEC-023`,
  damit sie nicht als Annahme gelesen wird.

### 6 — Die Träger, und die Reihenfolge ihres Nachzugs

- **Im selben Zug wie ein Beispiel:** der Handbuch-Absatz seiner Oberfläche
  samt Startform **und** die Zeile in der Änderungshistorie des Handbuchs
  (verkörperte Regel, Kontext Bindung 5).
- **`SPEC-023`** trägt die Festlegung für alle drei Sprachen (Festlegung 5);
  die ADR trägt die Begründung. Die Spec ist Rang 2 und trägt **keine**
  ADR- oder Slice-Kennung — die Richtung Spec → ADR ist mechanisch
  verboten, und der Grund des Nachzugs steht deshalb in dieser ADR.
- **`.a-check.yml`**: der Kommentar der `examples`-Gruppe wird auf die
  Go-Aussage gezogen und nennt die Grenze (Festlegung 5).
- **`harness/README.md` §Werkzeuge bekommt seine Zeilen mit den Zielen,
  nicht vorher** — dieselbe Disziplin, die `ADR-0076` für die Kante
  `examples → contract` gewählt hat („die Kante erst mit ihrem Objekt").
  Kein Träger nennt ein Target, das es nicht gibt; `AGENTS.md` §4 listet
  Gates, und diese Ziele sind keine.
- **`README.md` bleibt unberührt** — es führt keine Beispiel-Clients.

### 7 — Was diese ADR nicht ändert

- **Die drei Go-Clients bleiben, wo sie sind**, mit ihrer Startform und
  ihren Handbuch-Zeilen; `examples/grpc-client/` bleibt der offene
  Go-Vorgänger.
- **`ADR-0076`s Form bleibt** — Ort auf der Wurzel, das Verbot des
  `internal/`-Imports, das Nebeneinander mit den Wegwerf-Clients, der
  Abgrenzungssatz im Handbuch, die vier Trigger.
- **`spec/architecture.md` bleibt unberührt**; `examples/**` ist eine
  Gate-Scope-Gruppe, kein `ARC-*` kommt hinzu.
- **`Dockerfile`, `.dockerignore`, `make image`, `make gates` und der
  Coverage-Messgegenstand bleiben unverändert** (Festlegung 3/4).
- **Kein Produktionscode, kein Eingriff in `internal/**` oder `cmd/**`.**
- **Keine Zusage des Lastenhefts wird berührt:** die Beispiele zeigen
  [`LH-FA-SST-006`](../../../spec/lastenheft.md)/[`LH-FA-SST-008`](../../../spec/lastenheft.md),
  sie ändern sie nicht.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

### A — Umfang: wieviele Clients, in welchen Oberflächen

| Option | Pro | Contra |
|---|---|---|
| A1 — nichts tun: keine fremdsprachige Beispiel-Quelle | kein Eingriff, keine neue Werkzeugkette, kein neuer Träger | die Nutzerentscheidung bleibt offen; die zwei Sprachen bekommen kein Vorbild, und die zwei Register-Einträge zum Handbuch-Nachzug hätten keinen neuen Träger; die Spec-Lücke (kein Anker für `examples/`) bliebe |
| A2 — **alle vier** Oberflächen × zwei Sprachen (acht Programme) | maximale Abdeckung; jede Handbuch-Schnittstelle bekäme ein fremdsprachiges Vorbild | acht Programme über zwei Werkzeugketten, dazu zwingend protoc samt zwei Sprach-Plugins und eine Festlegung für erzeugte Stubs (Ort, committet?, Verhältnis zu `gen/**` und zur `contract`-Gruppe = `ADR-0076` Trigger 2); je Sprache eine zweite Fremdbibliothek (NATS); der Zug wäre in einer Review-Sitzung nicht prüfbar — er verletzt die Slice-Größen-Regel, bevor er beginnt |
| A3 — die **drei Draht-Oberflächen** × zwei Sprachen (sechs Programme), ohne NATS | deckt alles Daten-tragende ab; NATS ist ein Signal, kein Draht | zieht den gRPC-Zwang mit (protoc, Stub-Festlegung, `gen/**`-Frage) und damit die größte offene Entscheidung dieses Zugs in einen Zug, der sie nicht braucht; der Umfang bleibt über der Drei-Liefer-Punkte-Grenze |
| **A4 — die HTTP-Familie je Sprache (vier Programme), gRPC und NATS als benannte Folge-Entscheidung (gewählt)** | teilt die teure Hälfte (Basis-Image, Bauziel, Träger) je Sprache, nicht je Client; kein Codegenerator und damit **kein** Berühren der `contract`-Frage; die zwei Oberflächen, die jede Integration braucht; vier Programme, zwei Werkzeugketten, in einer Sitzung je Sprache prüfbar | in C#/Kotlin gibt es kein gRPC- und kein NATS-Beispiel, in keiner Sprache ein gRPC-Beispiel (ausgesprochen, Festlegung 1); ein späteres gRPC-Beispiel braucht eine eigene Entscheidung |

### B — Bauform: wo die Werkzeugkette lebt und wo der Kontext herkommt

| Option | Pro | Contra |
|---|---|---|
| B1 — Stufen in der **Wurzel-`Dockerfile`**, Quellen per **Bind-Mount** zur Laufzeit (Muster des Generierungs-Ziels) | ein Dockerfile; kein Kontext-Zuwachs; `make image` baut die Stufen nicht (sie sind keine Vorfahren der Zielstufe — gemessen) | die Wurzel-Datei, die das Auslieferungs-Image baut, trägt zwei fremde Werkzeugketten, die dieses Image **nie** braucht; der Bau hinge am Laufzeit-Baum des Aufrufers (unversioniert, samt `.git`) statt an einem definierten Kontext — dieselbe Eigenschaft, die [`ADR-0085`](0085-build-kontext-ausnahme-test-only-zweck.md) Option D für die Messung verworfen hat; Bau-Ausgaben landen im Arbeitsbaum |
| B2 — Stufen in der Wurzel-`Dockerfile`, Quellen per **Kontext-Negation** (`!examples/**`) | ein Dockerfile, hermetischer Bau | eine **Verzeichnis**-Negation in `.dockerignore`; sie gilt kontextweit für **jede** Stufe (so sagt es die Datei selbst), vergrößert also `COPY . .` in `coverage` und `build` und lässt das Digest-Verhalten des ausgelieferten Baus auf Beispiel-Änderungen reagieren; sie ist an [`ADR-0085`](0085-build-kontext-ausnahme-test-only-zweck.md) Merkmal (ii) vorbei (genau **eine Datei**) |
| **B3 — eigenes Dockerfile je Sprache, Kontext = Sprach-Wurzelverzeichnis (gewählt)** | hermetisch (Quellen sind im Kontext, `COPY . .`); Wurzel-`Dockerfile` und `.dockerignore` bleiben unberührt; `make image`, sein Digest-Beleg und der Go-Bau-Kontext bleiben unverändert; der Datei-Kommentar nennt beim Lesen den Ort der Werkzeugkette (gemessen: die Wurzel-`.dockerignore` greift für einen Unterverzeichnis-Kontext nicht) | zwei weitere Dockerfiles; der `.proto` liegt nicht im Sprach-Kontext — für die HTTP-Familie irrelevant, für ein künftiges gRPC-Beispiel eine eigene Frage (und genau die, die Festlegung 1 vertagt) |
| B4 — Werkzeugkette auf dem **Host** installieren | kein Image, kein Digest, kürzester Bau | verletzt `AGENTS.md` §3.1 (Docker-only) und die Reproduzierbarkeit; kein Weg dieses Repos |

### C — Prüfweg: womit die Quellen am Verrotten gehindert werden

| Option | Pro | Contra |
|---|---|---|
| C1 — kein Träger (nur die Quellen) | kein Aufwand | die Beispiele verrotten still: es gibt keinen Lauf, der sie übersetzt; der schlechteste Fall — ein Vorbild, das nicht baut, ist schlechter als keines (`ADR-0076` Option C1) |
| C2 — ein **Gate** in `make gates` | sichtbar im Pflicht-Bündel, blockiert einen roten PR | das Bündel ist heute **vollständig netzlos** (gemessen: `--network none` in jedem Gate-Fragment); ein C#-/Kotlin-Bau braucht Paket-Bezug und würde das erste netzgebundene Gate einführen — eine strukturelle Änderung des Bündels, mit eigener Beweislast, für eine Aussage, die ein Werkzeug ebenso trägt |
| C3 — ein **Werkzeug**-Ziel je Sprache, ohne CI-Anbindung | klein, kein Gate-Eingriff | ein Ziel, das niemand fährt, ist kein Carrier — der Rest-Fall wäre nicht benannt, sondern verdeckt |
| **C4 — Werkzeug-Ziel je Sprache **und** ein eigener, nicht-blockierender Workflow über beide (gewählt)** | der Bau ist lokal fahrbar **und** wird im PR sichtbar; das Pflicht-Bündel bleibt netzlos und schnell; das Muster existiert im Repo (der E2E-Workflow ist bewusst nicht-blockierend, damit er `ci.yml`s kurzes Signal nicht verzögert) | ein neuer Workflow, dessen Vollständigkeit nach `AGENTS.md` §3.10 am realen Post-Push-Lauf hängt — das Risiko bleibt bis dahin offen und ist benannt; ein nicht-blockierender Lauf kann überlesen werden (`BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit`) |

### D — Die `examples`-Gruppe in `.a-check.yml`

| Option | Pro | Contra |
|---|---|---|
| D1 — die Sprach-Schlüssel **auf Vorrat** setzen (`csharp`, `kotlin`) | die Dateien werden gelesen; die Abdeckungs-Aussage der Gruppe klingt vollständig | real gemessen: **keine Kante** gewonnen — es gibt keine Schicht, gegen die ein fremdsprachiger Import verstoßen könnte; je Zieldatei entsteht ein Auflösungs-Hinweis auf stderr, den die Sensor-Doku „zu lesen, nicht zu ignorieren" nennt — dauerhafter Lärm ohne Prüfwert |
| D2 — `resolution` für C#/Kotlin konfigurieren | wäre der stärkste Träger, wenn er trägt | die Konfigurationsform für C# ließ sich in dieser Probe **nicht** belegen (das Symbol blieb unaufgelöst); sie auf Vorrat zu deklarieren hieße, ein Grün zu behaupten, das nichts prüft — dieselbe Klasse, die `harness/sensors/a-check.md` Grenze 2 („Auflösungs-Hinweis ist kein Grün") benennt |
| **D3 — die Aussage auf Go beschränken, die Lücke benennen, den kippenden Tag datieren (gewählt)** | der Kommentar wird **wahr** und bleibt lesbar; kein Vorrat; die Abhängigkeitsliste ist über die gepinnten Manifeste nachlesbar; der Tag, an dem die Grenze neu zu entscheiden ist, steht fest (gRPC) | die Nicht-Go-Quellen sind mechanisch **nicht** gewächtert — Wächter ist das Review (`AGENTS.md` §3.12: eine Aussage, die als Beleg gelesen wird, trägt ihren Ursprung) |

**Fazit:** A4, B3, C4, D3. A2/A3 scheitern an der Slice-Größe und daran,
dass sie die `contract`-Frage mitziehen, obwohl der Anlass sie nicht stellt;
B1/B2 an der Kontext- und Image-Semantik des ausgelieferten Baus; B4 an
`AGENTS.md` §3.1; C1/C3 daran, dass sie keinen bzw. einen unwirksamen Träger
setzen; C2 an der netzlosen Eigenschaft des Gate-Bündels; D1/D2 daran, dass
sie eine Prüfung behaupten, die gemessen nichts prüft.

## Konsequenzen

- Positiv: Die Nutzerentscheidung ist eingelöst, ohne die Form der Go-Clients
  anzutasten: zwei Sprachen, je ein Sprach-Wurzelverzeichnis, je eine
  Werkzeugkette, je ein Bauziel — die drei Go-Clients und ihre Handbuch-Zeilen
  bleiben, wie sie sind.
- Positiv: **Der ausgelieferte Bau ist nicht berührt.** Wurzel-`Dockerfile`,
  `.dockerignore`, `make image`, sein Digest-Beleg und der Go-Build-Kontext
  bleiben; die neue Werkzeugkette lebt in ihrem eigenen Kontext
  (Festlegung 3) — und `make gates` bleibt netzlos und schnell (Festlegung 4).
- Positiv: **Die zwei `AGENTS.md` §3.12-Stellen sind saniert statt gestillt:**
  der `examples`-Kommentar wird auf die Aussage gezogen, die er trägt
  (Festlegung 5), und die Spec bekommt mit `SPEC-023` einen Anker, der auch
  die drei **bestehenden** Go-Clients trägt — der Leser der Spec findet, was
  es gibt, nicht nur, was kommt.
- Positiv: Kein neues Gate, keine Schwellen-Senkung, kein Produktionscode,
  kein Umbau am Hexagon; die Coverage-Messfläche
  ([`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md))
  ist unberührt.
- Negativ: Es gibt **zwei** Bauwege im Repo (Go über den Wurzel-Bau, die
  übrigen Sprachen über ihre eigenen Dockerfiles und Ziele). Das ist die
  Setzung: der Wurzel-Bau ist der **Auslieferungs**-Bau, die Sprach-Bauten
  sind Werkzeuge. Wer sie zusammenlegt, ändert diese Entscheidung, nicht ihre
  Auslegung.
- Negativ mit Grenze: Die Nicht-Go-Quellen sind mechanisch ungeprüft, was
  ihren Import-Grenzen angeht (Festlegung 5) — Wächter ist das Review. Der
  Tag, an dem das kippt, ist benannt (gRPC).
- Negativ mit Grenze: Der Träger hängt an einem neuen, **nicht-blockierenden**
  Workflow, dessen Vollständigkeit nach `AGENTS.md` §3.10 am realen
  Post-Push-Lauf hängt; bis dahin ist das betroffene Risiko offen. Ein
  nicht-blockierender Lauf kann überlesen werden — dieselbe Beobachtung, die
  das Register als `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit`
  führt.
- Negativ mit Grenze: „übersetzt und die reinen Funktionen laufen" ist nicht
  „holt am laufenden Feed eine Änderung" (Festlegung 4, benannte Grenze).
- Folgepflicht (Handbuch-Zug, je Beispiel): der Zugriffs-Absatz seiner
  Oberfläche mit Programm-Pfad und Startform **und** die Zeile in der
  Änderungshistorie — im selben Zug, nicht danach.
- Folgepflicht (Spec-Zug, dieser Zug): `SPEC-023` in `spec/pflichtenheft.md`
  samt Zeile in §6 (externe Werkzeugketten-Verträge) und Änderungshistorie;
  ohne ADR-/Slice-Kennung im Spec-Text.
- Folgepflicht (Werkzeug-Zug, je Sprache): das `make`-Ziel, seine Zeile in
  `harness/README.md` §Werkzeuge (**kein** Gate, nicht in `GATE_CHECKS`),
  die Ziele selbst im `Makefile`-Fragment.
- Folgepflicht (CI-Zug, eigener Schritt): der nicht-blockierende Workflow
  über beide Ziele, mit dem §3.10-Risiko bis zum ersten realen grünen Lauf.
- Folgepflicht (Träger-Zug): der Kommentar der `examples`-Gruppe in
  `.a-check.yml` auf die Go-Aussage samt benannter Grenze.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `make examples-csharp` / `make examples-kotlin` (Werkzeug, **kein** Gate) | Die Beispiele der Sprache werden im digest-gepinnten Werkzeugketten-Image des Sprach-Wurzelverzeichnisses **übersetzt** und ihre netzlos prüfbaren Tests **gefahren**; der Exit-Code des Ziels wird direkt gelesen (`AGENTS.md` §3.9) | — (Werkzeug; **nicht** in `GATE_CHECKS`) |
| a-check (Modulzeile `wrong-direction`) | Die Gruppe `examples` (`examples/**`) prüft **Go**-Importe: ein Import aus `domain`, `ports`, `app` oder `adapters` in `examples/**` ist `wrong-direction`. Für Nicht-Go-Quellen trägt sie **nichts** (gemessen, Festlegung 5) — die Grenze ist benannt, nicht still | `make a-check` (im Gate-Bündel) |
| `docker build` (Wurzel) und `make image` | Der ausgelieferte Bau bleibt unberührt: die letzte Stufe bleibt `runtime`, die Beispiele liegen in ihrem eigenen Kontext, `.dockerignore` und `Dockerfile` der Wurzel ändern sich nicht; der Digest-Beleg folgt [`ADR-0044`](0044-image-beleg-semantik.md) | `make image` (kein Gate) |
| Nicht-blockierender Workflow | Ein PR, der eine Beispiel-Quelle bricht, wird **sichtbar** rot (nicht blockierend); sein Abschluss hängt am realen Post-Push-Lauf (`AGENTS.md` §3.10) | — (kein Gate) |
| Review-Prüfpflicht (nicht maschinell) | Import-Grenzen der Nicht-Go-Quellen; die Frage, ob ein Beispiel noch das zeigt, was es zeigen soll; die Verdrahtung über die reine Funktion hinaus | — |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Vier benannte Trigger, sonst permanent:

1. **Eine dritte Fremdsprache wird angefragt.** Dann ist zu prüfen, ob die
   Form trägt (Sprach-Wurzel, eigenes Dockerfile, eigenes Ziel) — und ob die
   Werkzeugketten-Frage in eine gemeinsame, sprachneutrale Form wandert.
   Die Form selbst bleibt bis dahin unverändert.
2. **Der gRPC-Client einer der Sprachen wird angefragt.** Dann tritt
   `ADR-0076` §Re-Evaluierungs-Trigger 2 ein (ein weiterer Konsument der
   `.proto`), und die zwei Fragen, die Festlegung 1 vertagt, sind zusammen zu
   entscheiden: wo die erzeugten Fremdsprachen-Stubs liegen (committet oder
   im Bau erzeugt), und ob die `contract`-Gruppe sprachneutral wird — samt
   `languages`/`resolution` in `.a-check.yml`. **Nicht** vorher auf Vorrat.
3. **Ein Beispiel verrottet trotz Bauziel** — es baut, bricht aber im Ablauf,
   und das fällt erst einem Leser auf. Dann ist der in Festlegung 4 benannte
   Rest-Fall zu **schließen** (eigener Smoke-Lauf gegen eine laufende
   Umgebung mit eigener Bindung), nicht weiter zu benennen — dieselbe Folge,
   die `ADR-0076` §Re-Evaluierungs-Trigger 3 vorzeichnet.
4. **Die gepinnte Basis einer Sprache wird abgekündigt / ihr Tag verschwindet**
   (Base-Image-Drift, wie ihn `make image-stale` für den Wurzel-Bau meldet).
   Dann ist der Pin zu heben — ein bewusster Commit, keine stille Auflösung —
   und die Form der Werkzeugkette neu zu bewerten, wenn kein Nachfolge-Tag die
   Rolle trägt.

Sonst permanent: `examples/` führt je Sprache ein Sprach-Wurzelverzeichnis mit
eigener, digest-gepinnter Werkzeugkette und eigenem Bau-/Prüfziel; die Go-Form
bleibt; die Zahl der Clients je Sprache und die Zahl der Sprachen ändern daran
nichts ([`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-17 | Accepted — Anlass: Nutzerentscheidung vom 2026-09-17, die Client-Beispiele in **zwei weiteren Sprachen** (C# und Kotlin) zu führen, samt Nutzeranweisung, dass das Pflichtenheft sie tragen muss. Unabhängiger Architect-Zug entscheidet Umfang (HTTP-Familie je Sprache), Ort (Sprach-Wurzel), Werkzeugketten- und Bauform (eigenes Dockerfile je Sprache, eigener Kontext), Prüfweg (Werkzeug-Ziel je Sprache plus nicht-blockierender Workflow, kein Gate), die Grenze der `examples`-Gruppe und den Spec-Anker `SPEC-023` samt Nachverankerung der drei bestehenden Go-Clients; supersedes eine Gestalt-Klausel der [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md); Modul 8 §Rollen-Regeln: „Architect schreibt" | gemessen in diesem Zug: die drei `examples/`-Go-Clients und ihr Fehlen in `spec/pflichtenheft.md`; der a-check-Probe-Baum (Sprach-Schlüssel, Auflösungs-Hinweis, Exit 2 bei unbekanntem Schlüssel); Kontext- und Stufen-Proben des Baus; `docker manifest inspect` für die vier Basis-Kandidaten |
| 2026-09-17 | Digest-Korrektur — `eclipse-temurin:21-jdk`-Zeile der Digest-Pinning-Tabelle (Zeile 254) per Folge-ADR ersetzt, 63- statt 64-stelliger Transkriptionsfehler behoben (`ADR-0093`, Anlass Review zu `slice-099`, F-1) | `aa22937` <!-- d-check:status-provenance --> |
| 2026-09-18 | Zitat-Korrektur — `docs/reviews/**`-Pfade durch Kennung ersetzt (`ADR-0073`) | `c2bc868` |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0087` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
