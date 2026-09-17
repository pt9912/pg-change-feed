# ADR-0090: Beispiel-Clients — die volle Matrix, drei Sprachen, der `.proto`-Weg

**Status:** Accepted — Supersedes [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md)
in **genau vier** Klauseln:

1. **§Entscheidung Festlegung 1** („Umfang: die HTTP-Familie je Sprache, zwei
   Clients je Sprache") samt ihrem dritten Bullet („**Was der Nutzerwunsch
   damit nicht bekommt — ausgesprochen:** in C# und Kotlin gibt es **kein**
   gRPC- und **kein** NATS-Beispiel, und in allen drei Sprachen gibt es
   weiterhin **kein** gRPC-Beispiel") und der einschränkenden Hälfte des
   §Entscheidung-Leitsatzes („… in der HTTP-Familie — HTTP-/JSON und SSE").
2. **§Entscheidung Festlegung 3**, soweit sie das Sprach-Wurzelverzeichnis als
   den **einzigen** Kontext festlegt („**Je Sprache ein eigenes Dockerfile im
   Sprach-Wurzelverzeichnis** … gebaut mit dem **Sprach-Wurzelverzeichnis als
   Kontext**"). In Kraft bleibt: der Kontext **ist** das Sprach-Wurzelverzeichnis
   — hinzu tritt für eine Sprache, deren Client die `.proto` liest, **ein**
   zusätzlicher, ausdrücklich **benannter** Kontext (Festlegung 2).
3. **§Entscheidung Festlegung 5**, letzter Bullet („**Der Tag, an dem die
   Grenze kippt, ist benannt:** mit dem gRPC-Client entsteht ein `gen/**`, das
   auch fremdsprachig gelesen wird"). Der Tag tritt ein; die vorhergesagte
   Folge tritt **nicht** ein (Festlegung 3).
4. **§Re-Evaluierungs-Trigger 2** — er tritt ein und ist mit dieser ADR
   **konsumiert**: die zwei dort vertagten Fragen (Ort der erzeugten
   Fremdsprachen-Stubs; sprachneutrale `contract`-Gruppe) sind hier
   entschieden.

Alles Übrige der [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) bleibt
**hiermit bestätigt** und wird nicht wiederholt: der Ort und die
Sprach-Wurzel (§2), die zwei-Image-Bauform mit eigenem Kontext und die
Digest-Pinnung (§3, außer der Klausel oben), der Werkzeug-Ziel-ohne-Gate-Weg
samt nicht-blockierendem Workflow (§4), die Zulässigkeitsregel der Quellen
(§5, außer dem letzten Bullet), der Träger-Nachzug nach seinem Objekt (§6) und
„Was diese ADR nicht ändert" (§7). Ebenfalls unberührt:
[`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
(Ort, Import-Grenze, `.a-check.yml`-Gruppen, der Kompilierpfad als
Anti-Verrottungs-Träger, das Verhältnis zu den Wegwerf-Clients) und
[`ADR-0079`](0079-nats-beispielclient-vierter-examples-client.md) (der vierte
Go-Client).

**Datum:** 2026-09-17

**Autor:** pt9912 (Architect-Rolle, unabhängiger Architect-Zug auf
Nutzerentscheidung vom 2026-09-17; anderer Kontext als der Planner-Lauf, der
den Umzug-Slice geschnitten hat, als die Architect-Züge, die
[`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) und
[`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
geschrieben haben, und als die Implementer-Läufe der Go-Beispiel-Clients —
Modul 8 §Rollen-Regeln: „Architect schreibt")

**Bezug:** [`LH-FA-SST-006`](../../../spec/lastenheft.md) (HTTP-/JSON-API),
[`LH-FA-SST-007`](../../../spec/lastenheft.md) (NATS-Wecksignal),
[`LH-FA-SST-008`](../../../spec/lastenheft.md) (Live-Change-Stream),
[`ARC-005`](../../../spec/architecture.md) (Driving Adapters),
[`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) (in vier Klauseln
superseded; sonst bestätigt),
[`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
(§Re-Evaluierungs-Trigger 2 — tritt ein und wird hier beantwortet),
[`ADR-0079`](0079-nats-beispielclient-vierter-examples-client.md),
[`ADR-0060`](0060-grpc-streaming-mechanismus.md) (der Server-Stream und die
`.proto`), [`ADR-0061`](0061-http-sse-zusaetzlich-zu-grpc.md) (SSE),
[`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md)
(§Entscheidung Festlegung 3: ein Pfad-Bereich mit Import-Berechtigung braucht
eine ADR), [`ADR-0084`](0084-sync-gate-fuer-generierte-artefakte.md)
(Sync-Gate nur für das committete Erzeugnis), `ADR-0085`
(Build-Kontext-Ausnahme-Klasse — **nicht** berührt, Festlegung 2),
[`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
(der Coverage-Nenner — **nicht** berührt), `ADR-0044` (Image-Beleg — nicht
berührt), [`ADR-0070`](0070-supersede-reichweite-und-klassengrenze.md)
(Reichweite eines teilweisen `Supersedes`),
[`SPEC-017`](../../../spec/pflichtenheft.md) / `SPEC-018` / `SPEC-020` /
`SPEC-021` / `SPEC-022` (die Draht-Festlegungen, die die Beispiele
**benutzen**), [`SPEC-023`](../../../spec/pflichtenheft.md) (Beispiel-Clients —
von dieser ADR **geschärft**), [`AGENTS.md`](../../../AGENTS.md)
§3.1/§3.5/§3.9/§3.10/§3.12/§3.13, `.a-check.yml`,
`proto/cdc/stream/v1/changestream.proto`, `Dockerfile`, `.dockerignore`,
`Makefile`, `harness/mk/*.mk`, `harness/README.md` §Sensors/§Werkzeuge,
`.github/workflows/{ci,e2e}.yml`, `examples/**`,
`docs/user/benutzerhandbuch.md` (§4 Zugriffs-Abschnitte, §5 Konfiguration),
`docs/plan/planning/next/slice-097-umzug-vertragsflaeche.md` <!-- d-check:status-provenance -->,
`docs/plan/planning/next/slice-095-beispiel-clients-drei.md` <!-- d-check:status-provenance -->

**Schärft:** [`SPEC-023`](../../../spec/pflichtenheft.md) — die ADR macht die
Ziel-Form der Beispiel-Clients für die **volle Matrix** verbindlich: Umfang je
Sprache und Oberfläche, der `.proto`-Weg der fremdsprachigen gRPC-Bauten, der
Ort der erzeugten Stubs, die Abhängigkeits-Pins und die Grenze der
`examples`-Gruppe. Das Pflichtenheft trägt die Festlegung (Rang 2), diese ADR
die Entscheidung und ihre Begründung; bei Widerspruch gilt der höhere Rang.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md`
§Ziel-Form: ADR (MADR).

---

## Kontext

**Die Nutzerentscheidung — drei Sätze, ein Umfang.** Der Auftraggeber hat am
2026-09-17 entschieden: „Die Client-Beispiele brauchen wir noch in zwei
weiteren Programmiersprachen: **C# und Kotlin**" · „Beachte die drei
Prog-Lang: **go, c# und kotlin**" · „Auch **grpc-Streaming**" — präzisiert zu
**„Bitte beachte matrix: (Alle clients) × (c#, go, kotlin)"**. Der Umfang ist
damit die **volle Matrix**: die vier Zugriffs-Oberflächen `http`, `sse`,
`grpc`, `nats` in **jeder** der drei Sprachen — zwölf Programme. Das
**erweitert** [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md), die Go die
vier Oberflächen und C#/Kotlin die HTTP-Familie zuweist und die weiteren
Oberflächen je Sprache ausdrücklich **„perspektivisch"** nennt.

**Der Ist-Stand — gemessen, nicht erinnert.**

| Prozedur | Ergebnis |
|---|---|
| `ls examples/*/` | **drei** Go-Programme: `http-client`, `sse-client`, `nats-client` — je ein `main`-Paket |
| `ls examples/csharp examples/kotlin` | existiert **nicht** — die vier Programme der [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) sind **nicht gebaut** |
| `ls examples/grpc-client` | existiert **nicht** |
| `ls gen/` | existiert **nicht** — der Stub liegt unter `internal/adapters/driving/grpc/streamv1` (`option go_package` in der `.proto`) |
| `.a-check.yml` | Gruppen `domain`/`ports`/`app`/`adapters`/`tooling`/`examples`; **keine** `contract`-Gruppe, **kein** `languages`-Schlüssel außer `go` |
| `grep -rn "examples-" Makefile harness/mk/ .github/workflows/` | **kein Treffer** — die Ziele `examples-csharp`/`examples-kotlin` existieren nicht |
| `ls .github/workflows/` | `ci.yml`, `e2e.yml` — **kein** Workflow für die Sprachziele |
| `[SPEC-023](../../../spec/pflichtenheft.md)` §2, Zeile *Sprachen und Umfang* | „Go führt die vier Zugriffs-Oberflächen …, C# und Kotlin die **HTTP-Familie** … Die weiteren Oberflächen je Sprache sind **perspektivisch**" |

Heißt: **der Anlass dieser ADR ist nicht der einzige offene Punkt.** Der
Umfang der [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) ist selbst
noch nicht geliefert; ihre vier Programme sind die Zellen 1–4 der neun
fehlenden. Diese ADR ändert daran nichts, sie **erweitert** ihn — und sie
nennt die zwei zuerst zu schneidenden Slices beim Namen (§Slice-Schnitt).

**Was [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) entschieden hat —
und was von ihren drei Gründen heute noch trägt.** Sie schneidet gRPC und NATS
je Sprache aus, mit drei prüfbaren Gründen:

1. *„Sie teilen die Werkzeugkette vollständig"* — das galt für die
   HTTP-Familie **innerhalb** einer Sprache und bleibt wahr; es trägt keine
   Aussage über den Umfang (die teure Hälfte — Basis-Image, Manifest, Bauziel
   — fällt **je Sprache** an, ob mit zwei oder mit vier Programmen).
2. *„Sie brauchen keinen Codegenerator"* — das war der **Grund für den
   Zuschnitt**, nicht gegen gRPC. Die vertagte Frage („wo die erzeugten
   Stubs liegen, ob committet, was `contract` dann bedeutet") ist mit dieser
   ADR entschieden (Festlegung 3).
3. *„Die zwei Oberflächen, die jede Integration braucht"* — eine
   Priorisierung, keine Grenze; die Nutzerentscheidung setzt sie neu.

**Keiner der drei Gründe ist widerlegt — sie sind erledigt.** Genau deshalb
ist die Erweiterung eine **Entscheidung** und keine Auslegung: sie kehrt eine
Umfangs-Festlegung um, deren Argumente im Umfang selbst lagen.

### Die eine Frage, die dieser Zug beantworten muss: erreichen die C#-/Kotlin-gRPC-Bauten die `.proto`?

Der Umzug-Slice (heute in `next/`) benennt sie als **Folgepflicht** und
beantwortet sie ausdrücklich nicht: „Die C#-/Kotlin-Clients erzeugen ihre
Stubs **aus der `.proto`** — die **nicht** umzieht —, aber ihre
Dockerfile-Kontexte sind **eigene** ([`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md));
ob sie die `.proto` **erreichen**, ist dort zu beantworten." Antwort, hier
und jetzt:

**Sie erreichen sie nicht.** Die [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md)
Festlegung 3 setzt den **Kontext** eines Sprach-Baus auf das
Sprach-Wurzelverzeichnis (`examples/csharp/`, `examples/kotlin/`). Die `.proto`
liegt unter `proto/cdc/stream/v1/` — **außerhalb** dieses Kontexts. Docker
kopiert nicht aus dem Kontext heraus; der Bau kann die Datei nicht sehen. Der
Umzug hilft ihnen **nicht**: er bewegt die **Go-Bindung** (`gen/**`), nicht
die `.proto` — er ist die Voraussetzung des **Go**-gRPC-Clients und **einer**
Zelle der neun, nicht der drei gRPC-Zellen.

**Und der Weg, der das löst, ist real gemessen, nicht angenommen** (dieser
Zug, `docker buildx version` → `v0.37.1`):

| Probe | Ergebnis |
|---|---|
| `docker buildx build --build-context proto=<proto-Verzeichnis> -f <Sprach-Dockerfile> <Sprach-Wurzel>` mit `COPY --from=proto cdc/stream/v1/changestream.proto …` in einer `scratch`-Stufe, `--output type=local` | **die Datei kommt an** (25 B im Ausgabeverzeichnis) — ein **benannter Zusatzkontext** erreicht eine Quelle außerhalb des Bau-Kontexts |
| Ignore-Auflösung des benannten Kontexts — **dokumentiert, nicht gemessen** (Docker-Doku, *What is a build context?*: die Ignore-Suche gilt für die Wurzel **des Kontexts**) | für den Kontext `proto/` wäre das `proto/.dockerignore` — **nicht vorhanden**, also liegt die `.proto` **im** Kontext, ohne dass die Wurzel-`.dockerignore` angefasst wird; die Probe darüber lief mit genau diesem Baum (kein `.dockerignore` in ihm) |
| Dokumentierte Alternative: `<Dockerfile>.dockerignore` neben dem Dockerfile | nimmt **Vorrang** vor der Wurzel-`.dockerignore` (Docker-Doku, *Filename and location*) — trägt, ist aber der zweite Weg (siehe §Verglichene Alternativen, B2) |

**Fünf Bindungen prägen den Lösungsraum — alle fünf aus dem Repo.**

1. **Docker-only** ([`AGENTS.md`](../../../AGENTS.md) §3.1): kein Host-`dotnet`,
   kein Host-`gradle`, kein Host-`protoc`. Jede Sprache bekommt ihre
   Werkzeugkette im gepinnten Image — die Form steht
   ([`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) Festlegung 3), nichts
   daran ändert sich.
2. **Kein Gate.** Jedes Target in `GATE_CHECKS` läuft netzlos; ein
   C#-/Kotlin-Bau braucht Paket-Bezug und damit Netz. Die Sprachziele bleiben
   **Werkzeuge** ([`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md)
   Festlegung 4); `make gates` bleibt, was es ist.
3. **Der Träger ist ein nicht-blockierender Workflow** — und nach
   [`AGENTS.md`](../../../AGENTS.md) §3.10 gilt ein neuer oder strukturell
   geänderter Workflow erst nach einem realen, grünen Post-Push-Lauf als
   abgeschlossen. Das betroffene Risiko bleibt bis dahin **offen**
   (Register: `BEO-PGC/github-actions-unverifizierbar-lokal`,
   `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit`).
4. **Der ausgelieferte Bau ist unantastbar.** Wurzel-`Dockerfile`,
   `.dockerignore`, `make image` und der Go-Build-Kontext bleiben
   ([`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) Festlegung 3). Der
   `.dockerignore`-Ausnahmeklasse der `ADR-0085` („genau **eine** Datei") ist
   eine **Verzeichnis**-Negation nicht zugänglich.
5. **Eine Arbeit, die eine beschriebene Eigenschaft bewegt, zieht ihre Träger
   nach** ([`AGENTS.md`](../../../AGENTS.md) §3.13): die `.proto` bekommt drei
   weitere **lesende** Bauwege, und `SPEC-023` trägt eine Zeile, die mit dieser
   ADR falsch wird.

**Zwei Register-Einträge treffen diesen Zug.** `BEO-PGC/aufschub-adresse-nimmt-sendung-nicht-an`
(**2×**, offen) — der Umzug ist seit
[`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
benannt; diese ADR fügt dem **keine** weitere Adresse hinzu, sie benennt die
**Bedingung** der neun Zellen (Sprach-Wurzel; benannter Kontext; Umzug nur für
eine). `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
(**3×**, verkörpert) — der Handbuch-Nachzug je Beispiel bleibt Pflicht
(Festlegung 7), und mit zwölf Programmen wird die **Form** dieses Nachzugs
selbst zur Entscheidung.

## Entscheidung

Wir wählen: **`examples/` führt die Beispiel-Clients über die volle Matrix —
die vier Oberflächen `http`, `sse`, `grpc`, `nats` in jeder der drei Sprachen
Go, C# und Kotlin; ein fremdsprachiger gRPC-Bau erreicht die `.proto` über
**einen** zusätzlichen, benannten Bau-Kontext, während das
Sprach-Wurzelverzeichnis der Bau-Kontext bleibt; die erzeugten
Fremdsprachen-Stubs entstehen **im Bau** und werden **nicht** committet, so
dass `gen/**` die Go-Bindung bleibt; NATS bekommt je Sprache eine gepinnte
öffentliche Client-Bibliothek; die `examples`-Gruppe in `.a-check.yml` bleibt
eine Go-Aussage mit benannter Grenze.**

Acht Festlegungen:

### 1 — Umfang: die volle Matrix

- **Aufgenommen — neun neue Programme** (drei existieren bereits):

  | Sprache | `http` | `sse` | `grpc` | `nats` |
  |---|---|---|---|---|
  | **Go** | vorhanden (`examples/http-client/`) | vorhanden (`examples/sse-client/`) | **neu** (`examples/grpc-client/`) | vorhanden (`examples/nats-client/`) |
  | **C#** | **neu** | **neu** | **neu** | **neu** |
  | **Kotlin** | **neu** | **neu** | **neu** | **neu** |

  Jedes Programm bleibt, was [`SPEC-023`](../../../spec/pflichtenheft.md)
  festlegt: ein eigenständiges Programm, das **eine** dokumentierte
  Oberfläche real anspricht und ihre Ausgabe ausgibt — Vorbild, kein
  Belegträger, keine Zustandsmaschine.
- **Der Zuschnitt je Sprache, den diese ADR nicht öffnet:** jede Sprache führt
  **ihr** Sprach-Wurzelverzeichnis und ihre Werkzeugkette
  ([`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) Festlegung 2). Für
  Go gilt `examples/<name>/` flach weiter — auch für den **neuen** Go-gRPC-Client.
- **Der Preis je Zelle — und was daran nicht trivial ist.** Die teure Hälfte
  (Basis-Image, Abhängigkeits-Pin, Bauziel, Workflow) fällt **je Sprache** an,
  nicht je Zelle; sie ist für C#/Kotlin bereits mit
  [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) entschieden und **noch
  nicht gebaut**. Darüber hinaus:

  | Zelle | Neuer Preis gegenüber der HTTP-Familie |
  |---|---|
  | `http` (C#/Kotlin) | **keiner** — die einzige Zelle, die [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) schon trägt |
  | `sse` (C#/Kotlin) | **keine neue Abhängigkeit** — beide Runtimes tragen Streaming-HTTP in ihrer Standardausstattung (C#: der `HttpClient` der Runtime samt ihrem SSE-Parser `System.Net.ServerSentEvents`; Kotlin/JVM: `java.net.http.HttpClient` mit einem streamenden Body-Handler, `java.net.http` seit JDK 11). Welcher der beiden Wege je Sprache gewählt wird, ist Detail des umsetzenden Zuges; die Zulässigkeitsregel (Festlegung 4) lässt beide zu |
  | `grpc` (C#/Kotlin) | **der benannte Zusatzkontext** (Festlegung 2) **und** der Protobuf-/gRPC-Generator der Sprache — gemessen existent (Festlegung 4) |
  | `nats` (C#/Kotlin) | **eine gepinnte öffentliche Client-Bibliothek je Sprache** — gemessen existent (Festlegung 4) |
  | `grpc` (Go) | **der Umzug der Vertragsfläche** (der Umzug-Slice) <!-- d-check:status-provenance --> — die einzige Zelle, an der er hängt |

- **Was der Nutzerwunsch damit vollständig bekommt:** zwölf Programme, vier
  Oberflächen, drei Sprachen. Es bleibt **keine** Zelle der Matrix offen.

### 2 — Der `.proto`-Weg: der benannte Zusatzkontext

- **Der Bau-Kontext bleibt das Sprach-Wurzelverzeichnis.** Für eine Sprache,
  deren Client die `.proto` liest — heute nur `grpc` —, kommt **genau ein**
  zusätzlicher, ausdrücklich **benannter** Kontext hinzu, der das
  `proto`-Verzeichnis trägt (`docker buildx build --build-context …`); der
  Bau kopiert die `.proto` daraus (`COPY --from=<name> …`). Real gemessen
  (Kontext): der Weg trägt, und die Wurzel-`.dockerignore` bleibt unberührt.
- **Warum genau ein Zusatzkontext und nicht mehr.** Ein weiterer benannter
  Kontext ist eine weitere Quelle, die driften kann. Der Zusatzkontext ist
  **eine** Datei-Baum-Wurzel (`proto/`) auf **einen** Leser hin (den
  Sprach-Generator) — er wächst nicht mit der Zahl der Clients, sondern mit
  der Zahl der Sprachen, die die `.proto` lesen.
- **Die zwei anderen Sprachen bleiben unberührt.** `http`, `sse` und `nats`
  brauchen **keinen** Zusatzkontext; ihr Bau bleibt exakt die Form der
  [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) Festlegung 3. Der
  **Go**-Bau (Wurzel-`Dockerfile`) wird von alldem nicht berührt.
- **Der Kontextname gehört dem umsetzenden Zug** (Vorschlag: `proto`);
  festgelegt ist, *dass* der Zusatzkontext benannt, einzeln und auf seinen
  Leser hin benannt ist — nicht, wie er heißt.
- **Die `.proto` bleibt die einzige Quelle des Draht-Vertrags.** Sie wird
  weder kopiert noch committet noch gespiegelt; die drei fremdsprachigen Bauten
  lesen dieselbe Datei, die der Go-Bau liest.

### 3 — Die erzeugten Stubs: im Bau erzeugt, nicht committet

- **Die C#- und Kotlin-Stubs aus der `.proto` werden im Sprach-Bau erzeugt und
  liegen nicht im Baum.** Damit bleibt **`gen/**` die Go-Bindung** — und die
  von [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) Trigger 2 vertagte
  Frage ist beantwortet: **kein** Sprach-Schlüssel, **keine** `resolution`,
  **keine** sprachneutrale `contract`-Gruppe.
- **Drei Gründe, jeder prüfbar:**
  1. **Ein committetes Erzeugnis braucht sein Sync-Gate.**
     [`ADR-0084`](0084-sync-gate-fuer-generierte-artefakte.md) bindet das
     Sync-Gate an das Erzeugnis mit **netzloser, deterministischer** Quelle.
     Drei weitere committete Bäume wären drei weitere Sync-Gates — für Code,
     den niemand von Hand liest.
  2. **Ein committeter Fremdsprachen-Stub wäre eine zweite Repräsentation des
     Vertrags.** Die `.proto` ist die Quelle; eine eingecheckte Ausgabe ist
     eine Kopie, die gegen sie driften kann
     (`AGENTS.md` §3.13 — die Arbeit bewegte die beschriebene Eigenschaft).
  3. **Er träte in die Gate-Fläche ein.** Eine deklarierte Sprache liest
     a-check; für erzeugten Fremdsprachen-Code gibt es in diesem Repo **keine**
     Schicht, gegen die er verstoßen könnte — das Ergebnis wäre der
     Auflösungs-Hinweis je Datei, also Lärm ohne Kante (Festlegung 6).
- **Die Folge, ausgesprochen:** die Vorhersage der
  [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) Festlegung 5 („mit dem
  gRPC-Client entsteht ein `gen/**`, das auch fremdsprachig gelesen wird")
  tritt **nicht** ein. Der benannte Tag ist da, die Grenze kippt **nicht** —
  und das ist eine Entscheidung, kein Unterlassen.
- **Die `.proto` wird dadurch zur gemeinsamen Berührungsfläche.** Der
  Umzug-Slice ändert `option go_package` in ihr <!-- d-check:status-provenance -->;
  ein gRPC-Client-Slice kann sie für Sprach-Codegen-Optionen berühren. Wer sie
  anfasst, nennt die andere Arbeit (Konflikt-Fläche, nicht Wettlauf).

### 4 — Was die Quellen benutzen dürfen

- **Die Zulässigkeitsregel bleibt wörtlich** („ausschließlich die
  Standardbibliothek/Runtime der Sprache und **öffentliche** Fremdmodule; kein
  Import eines privaten Baums dieses Repos"). Die drei fremdsprachigen
  gRPC-Bauten erfüllen sie: sie lesen die **öffentliche** `.proto` dieses
  Repos und erzeugen daraus ihren eigenen Stub.
- **Die Abhängigkeiten je Sprache — Kandidaten gemessen am 2026-09-17, der
  Pin gehört dem umsetzenden Zug** (Quelle je Zeile: die Registry-Abfrage, die
  die Existenz belegt — nicht eine Schätzung, `AGENTS.md` §3.12):

  | Sprache | Zweck | Gemessener Kandidat (Quelle) |
  |---|---|---|
  | C# | NATS-Client | `NATS.Net` — NuGet (`api.nuget.org/v3-flatcontainer/nats.net/index.json`), Versionen bis `2.1.2`+ |
  | C# | gRPC-Laufzeit | `Grpc.Net.Client` — NuGet, bis `2.83.0` |
  | C# | Stub-Erzeugung | `Grpc.Tools` (trägt `protoc` und das C#-Plugin) — NuGet, bis `2.84.0` |
  | C# | Protobuf-Laufzeit | `Google.Protobuf` — NuGet, bis `3.36.1` |
  | C# | SSE (Wahl) | `System.Net.ServerSentEvents` — NuGet, ab `9.0.0` |
  | Kotlin/JVM | NATS-Client | `io.nats:jnats` — Maven Central (`repo1.maven.org/maven2/io/nats/jnats/maven-metadata.xml`), bis `2.26.3` |
  | Kotlin/JVM | gRPC-Stub | `io.grpc:grpc-kotlin-stub` — Maven Central, `1.5.0` |
  | Kotlin/JVM | Stub-Erzeugung | `io.grpc:protoc-gen-grpc-kotlin` — Maven Central, `1.5.0` |
  | Kotlin/JVM | gRPC-Transport | `io.grpc:grpc-netty-shaded` — Maven Central, `1.84.0` |
  | Kotlin/JVM | SSE (Alternative) | `com.squareup.okhttp3:okhttp-sse` — Maven Central, `5.5.0` (nicht nötig; die Runtime trägt Streaming-HTTP) |

- **Die zwei Sprach-Basen existieren weiterhin** (real gemessen:
  `docker manifest inspect mcr.microsoft.com/dotnet/sdk:10.0` und
  `eclipse-temurin:21-jdk`, je Exit 0). Die Digest-**Zeile** bleibt die der
  [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md); dieser Zug hat die
  Existenz geprüft, nicht die Digests neu gemessen.
- **Erwartet, nicht gemessen:** Bau-Zeit und Image-Größe der Sprach-Bauten,
  und die Frage, ob ein C#-/Kotlin-gRPC-Bau im Runner-Zeitbudget des
  Workflows liegt. Eine Zahl hier wäre eine Schätzung; sie zu messen ist Sache
  des umsetzenden Zuges, der sie in seiner Closure-Notiz führt.

### 5 — Der Bau- und Prüfweg: je Sprache ein Werkzeug, kein Gate

- **Je Sprache ein `make`-Ziel** (`examples-csharp`, `examples-kotlin`), das
  die Beispiele der Sprache übersetzt und ihre **netzlos** prüfbaren Teile
  fährt. Unverändert gegenüber
  [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) Festlegung 4 — nur
  wächst der Umfang, den das Ziel je Sprache trägt (vier Programme statt zwei).
- **Der Go-gRPC-Client braucht kein neues Ziel.** Er liegt im Wurzelmodul; der
  Kompilierpfad `make test` trägt ihn
  ([`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  Festlegung 5).
- **Kein Gate, kein Compose-Lauf** — unverändert
  ([`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) Festlegung 4); die
  Beispiele bleiben Doku mit Bindungs-Pflicht und erscheinen nicht in
  `docs/user/e2e-abdeckung.md`.
- **Der Carrier wächst mit seinem Umfang:** der nicht-blockierende Workflow
  über die Sprachziele deckt nach dieser ADR vier Programme je Sprache statt
  zwei. Ob das im Zeitbudget eines Runner-Laufs bleibt, ist **erwartet, nicht
  gemessen** (Festlegung 4) — und nach
  [`AGENTS.md`](../../../AGENTS.md) §3.10 bleibt der Workflow **offen**, bis
  ein realer, grüner Post-Push-Lauf ihn bestätigt.
- **§3.9 gilt unverändert:** der Exit-Code eines Ziels wird direkt gelesen,
  nie durch eine Pipe; die abhängige Folgehandlung wartet auf ihn.

### 6 — Die `examples`-Gruppe bleibt eine Go-Aussage, mit benannter Grenze

- **Kein `languages`-Schlüssel für `csharp`/`kotlin`, kein `resolution`, keine
  neue Kante.** Der Grund ist derselbe wie in
  [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) Festlegung 5 und wird
  durch Festlegung 3 **nicht** überholt: es gibt in diesem Repo **keine**
  C#-/Kotlin-Schicht, gegen die ein fremdsprachiger Import verstoßen könnte;
  ein deklarierter Sprach-Schlüssel erzeugte je Datei einen
  Auflösungs-Hinweis, also Lärm ohne Kante
  (`harness/sensors/a-check.md` Grenze 2).
- **Die Grenze wird dabei breiter, und sie wird weiter benannt.** Der Kommentar
  der `examples`-Gruppe sagt künftig: eine **Go**-Aussage; die nicht-Go-Quellen
  sind mechanisch ungeprüft; und die erzeugten Fremdsprachen-Stubs liegen
  **nicht im Baum** (Festlegung 3), also liest a-check sie ohnehin nicht. Der
  Wächter bleibt das Review (`AGENTS.md` §3.12).
- **`gen/**` bleibt Go.** Die `contract`-Gruppe und ihre Kanten
  (`adapters`/`tooling`/`examples` → `contract`) sind von dieser ADR **nicht**
  berührt. Damit ist [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  §Re-Evaluierungs-Trigger 2 **beantwortet** — mit „keine Änderung, weil kein
  Sprach-Pfad gelesen wird", nicht mit einem stillen Unterlassen: die drei
  neuen Konsumenten lesen die `.proto`-**Quelle**, keinen Pfad-Bereich dieses
  Repos.
- **Was a-check damit nicht trägt und nie trug, bleibt die Klasse:** die
  Import-Grenze der Nicht-Go-Quellen (Review), der Ablauf „kompiliert ist
  nicht läuft" (Review), und die Frage, ob ein Beispiel noch zeigt, was es
  zeigen soll (Review).

### 7 — Die Träger, und die Form des Handbuch-Nachzugs

- **`[SPEC-023](../../../spec/pflichtenheft.md)`** trägt die Festlegung für die
  volle Matrix. Die Zeile *Sprachen und Umfang* verliert die Beschränkung
  („C# und Kotlin die **HTTP-Familie** … **perspektivisch**") und nennt statt
  ihrer die Matrix, den benannten Zusatzkontext (Festlegung 2), den Ort der
  erzeugten Stubs (Festlegung 3) und die Abhängigkeits-Pins (Festlegung 4).
  Die Spec ist Rang 2 und trägt **keine** ADR- oder Slice-Kennung; der Grund
  des Nachzugs steht deshalb hier. §6 (externe Werkzeugketten-Verträge) hat
  seine Zeile bereits und bleibt; §7 bekommt seine Änderungshistorie-Zeile.
- **Das Handbuch bekommt *keine* zwölf Absätze — es bekommt vier Blöcke mit je
  drei Zeilen.** Heute trägt jeder der vier Zugriffs-Abschnitte **einen**
  `**Beispiel:**`-Absatz (vier Absätze, vier Programme). Die naive Erweiterung
  wäre ein Absatz je Programm — **zwölf** Absätze, und die Sprache würde zur
  zweiten Gliederungsachse neben der Oberfläche. Wir wählen die andere Form:
  die Prosa des Abschnitts (Erreichbarkeit, Authentifizierung, Stream,
  Zustellsemantik) bleibt **unverändert und einfach**; an die Stelle des
  einzelnen `**Beispiel:**`-Absatzes tritt **ein** `**Beispiele:**`-Block mit
  **einer Zeile je Sprache** — Programm-Pfad und Startform je Zeile. Das sind
  **zwölf Zeilen in vier Blöcken**; die Oberfläche bleibt die Gliederung, die
  Sprache wird eine Liste darin. Bleibt es lesbar? Ja — und der Vergleich ist
  der Grund: dieselbe Information als zwölf Fließtext-Absätze wäre es nicht.
- **Die Handbuch-Zeile gehört in denselben Zug wie ihr Beispiel** — die
  verkörperte Regel aus
  `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` (**3×**),
  samt Zeile in der Änderungshistorie des Handbuchs
  (`BEO-PGC/handbuch-versionshistorie-uebersprungen`).
- **`harness/README.md` §Werkzeuge** bekommt seine Zeilen mit den Zielen,
  nicht vorher: die zwei Sprachziele und den Workflow. Dieselbe Disziplin wie
  [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md):
  kein Träger nennt ein Target, das es nicht gibt. `AGENTS.md` §4 bleibt
  unberührt — diese Ziele sind **keine** Gates.
- **`.a-check.yml`** wird nach Festlegung 6 nachgezogen: der Kommentar der
  `examples`-Gruppe auf die Go-Aussage samt der weiteren Grenze.
- **`README.md` bleibt unberührt** — es führt keine Beispiel-Clients.

### 8 — Was diese ADR nicht ändert

- **Der Go-Bestand bleibt.** Die drei vorhandenen Go-Clients, ihre flache Form
  (`examples/<name>/`), ihre Startform `go run ./examples/<name>` und ihre
  Handbuch-Zeilen bleiben; `examples/grpc-client/` bleibt der offene
  Go-Vorgänger und wird **eine** der neun Zellen.
- **Der ausgelieferte Bau bleibt.** Wurzel-`Dockerfile`, Wurzel-`.dockerignore`,
  `make image`, sein Digest-Beleg (`ADR-0044`) und der Go-Build-Kontext ändern
  sich **nicht**; `make gates` bleibt netzlos und schnell.
- **Der Coverage-Messgegenstand bleibt**
  ([`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)):
  diese ADR erzeugt **kein** committetes Erzeugnis außerhalb des
  Messgegenstands (Festlegung 3). Der Umzug-Slice bewegt ihn — diese ADR nicht.
- **`spec/architecture.md` bleibt unberührt.** `examples/**`, `gen/**` und die
  Sprach-Wurzeln sind Gate-Scope-Bereiche, keine `ARC-*`-Komponenten.
- **Kein Produktionscode, kein Eingriff in `internal/**` oder `cmd/**`.**
- **Keine Zusage des Lastenhefts wird berührt:** die Beispiele **zeigen**
  [`LH-FA-SST-006`](../../../spec/lastenheft.md)/`LH-FA-SST-007`/`LH-FA-SST-008`,
  sie ändern sie nicht.
- **Kein neues Gate, keine Schwellen-Senkung, keine Ausnahme in
  `.a-check.yml`.**

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

### A — Umfang: welche Zellen, in welchem Zug

| Option | Pro | Contra |
|---|---|---|
| A1 — nichts tun: die HTTP-Familie aus [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) bleibt der Umfang | kein neuer Träger, keine neue Abhängigkeit, keine `.proto`-Frage | die Nutzerentscheidung bleibt offen; gRPC und NATS fehlen in C#/Kotlin und gRPC in Go; die Matrix bliebe halb, und `SPEC-023` müsste die Beschränkung weiterschreiben, die niemand mehr will |
| A2 — die drei **Draht**-Oberflächen × zwei Sprachen (ohne NATS) | deckt alles daten-tragende ab; NATS ist ein Signal | lässt die ausdrücklich genannte Zelle `nats` offen; die Begründung wäre eine Priorisierung, die der Auftrag schon entschieden hat |
| **A3 — die volle Matrix: vier Oberflächen × drei Sprachen, neun neue Programme in mehreren Slices (gewählt)** | die Nutzerentscheidung ist vollständig eingelöst; keine Zelle bleibt offen; die teure Hälfte je Sprache fällt ohnehin an (Festlegung 1) | neun Programme über zwei (für Go: kein) Werkzeugketten, drei gepinnte Fremdbibliotheken, ein Generator je Sprache, ein neuer Bau-Kontext; nur als **geschnitter** Zug prüfbar |
| A4 — die volle Matrix in **einem** Slice | ein Zug, eine Entscheidung, ein Review | neun Programme, zwei Werkzeugketten, drei Abhängigkeiten, ein Generator und ein neuer Bau-Kontext in einer Sitzung: der Slice ist bei Beginn zu groß (Modul 5 §Ziel-Form) und verletzt die Drei-Liefer-Punkte-Grenze, bevor er beginnt |

### B — Der `.proto`-Weg: wie der fremdsprachige gRPC-Bau die Quelle erreicht

| Option | Pro | Contra |
|---|---|---|
| B1 — der Bau hängt am **Laufzeit-Baum** (Bind-Mount der Wurzel in den Bau) | kein Zusatzkontext, kürzester Weg | der Bau hinge am unversionierten Baum des Aufrufers (samt `.git`) statt an einem definierten Kontext — dieselbe Eigenschaft, die `ADR-0085` §Verglichene Alternativen Option D für die Messung verworfen hat und [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) Option B1 für den Bau; kein reproduzierbarer Bau |
| B2 — Kontext = **Repo-Wurzel**, Sprach-Dockerfile per `-f`, Ignore über `<Dockerfile>.dockerignore` | ein Kontext, ein dokumentierter Mechanismus (Vorrang vor der Wurzel-`.dockerignore`, Doku gemessen); die Wurzel-`.dockerignore` bleibt unangetastet | gibt die [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md)-Form („Kontext = Sprach-Wurzel") auf; **jede** Sprach-Dockerfile-Zeile wird wurzel-relativ; eine neue Ignore-Datei-Klasse je Sprache, deren Fehler still ist (zu viel oder zu wenig im Kontext); die Kontext-**Auswahl** liegt in der Ignore-Datei statt im Aufruf |
| B3 — die `.proto` wird in das Sprach-Wurzelverzeichnis **kopiert** und committet | einfacher Kontext, kein neuer Docker-Schalter | **zwei weitere Quellen** desselben Vertrags, die gegen ihn driften können: sie bräuchten je ein Sync-Gate ([`ADR-0084`](0084-sync-gate-fuer-generierte-artefakte.md)) und einen Pin-Hebungs-Trigger — für eine Datei, die schon existiert und schon öffentlich ist |
| **B4 — Bau-Kontext bleibt die Sprach-Wurzel, **ein** zusätzlicher, benannter Kontext trägt die `.proto` (gewählt)** | ändert die [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md)-Form **nicht**; die Wurzel-`.dockerignore` und `make image` bleiben unberührt; die `.proto` bleibt die **einzige** Quelle; real gemessen (Probe: die Datei kommt an); die Auswahl steht **im Aufruf**, nicht in einer Ignore-Datei | eine neue Bauform im Repo (benannter Zusatzkontext); sie hängt an `docker buildx`, das der Repo-Bau ohnehin benutzt (`make image` fährt `docker buildx build`) — der Zusatzkontext selbst ist noch nicht im offiziellen Bau gemessen, nur in einer Probe |
| B5 — die fremdsprachigen gRPC-Clients entfallen; die `.proto` bleibt für sie unerreichbar | kein neuer Bau-Mechanismus | gibt eine der vier Oberflächen in zwei Sprachen auf — genau die Zelle, die die Nutzeranweisung nachträglich und ausdrücklich verlangt hat |

### C — Die erzeugten Stubs: wohin sie gehen

| Option | Pro | Contra |
|---|---|---|
| C1 — **im Bau erzeugt, nicht committet** (gewählt) | `gen/**` bleibt die Go-Bindung; kein Sync-Gate je Sprache ([`ADR-0084`](0084-sync-gate-fuer-generierte-artefakte.md)); keine zweite Repräsentation der `.proto`; a-check liest nichts Zusätzliches | der Leser sieht den erzeugten Stub nicht im Baum; der Stub ist nur über den Bau reproduzierbar — und die Sprache braucht dafür ihren Generator im gepinnten Image |
| C2 — committet **mit** Sync-Gate je Sprache | der Stub ist lesbar und die Ausgabe gegen die Quelle gehalten | drei weitere committete Bäume, drei weitere Sync-Gates, drei weitere Pin-Hebungs-Trigger — für Code, den niemand von Hand liest; die `.proto` bekäme drei Kopien |
| C3 — committet **ohne** Sync-Gate | am wenigsten Aufwand | die Kopien driften still gegen die `.proto` — die Klasse, gegen die [`ADR-0084`](0084-sync-gate-fuer-generierte-artefakte.md) entschieden wurde |

### D — Die `examples`-Gruppe in `.a-check.yml`

| Option | Pro | Contra |
|---|---|---|
| D1 — Sprach-Schlüssel `csharp`/`kotlin` aufnehmen | die Quellen werden gelesen; die Abdeckung klingt vollständig | gemessen: **keine** Kante gewonnen (es gibt keine fremdsprachige Schicht); je Zieldatei ein Auflösungs-Hinweis, den die Sensor-Doku „zu lesen, nicht zu ignorieren" nennt — dauerhafter Lärm ohne Prüfwert |
| D2 — `resolution` für C#/Kotlin konfigurieren | wäre der stärkste Träger, wenn er trägt | die Konfigurationsform ließ sich in der Probe der [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) **nicht** belegen; sie auf Vorrat zu deklarieren hieße, ein Grün zu behaupten, das nichts prüft |
| **D3 — Go-Aussage, Grenze benannt und um die nicht-im-Baum-liegenden Stubs erweitert (gewählt)** | der Kommentar wird **wahr**; kein Vorrat; die Grenze wächst mit dem Umfang statt still zu bleiben | die Nicht-Go-Quellen sind mechanisch **nicht** gewächtert — Wächter ist das Review |

### E — Die Form der Erweiterung: Folge-ADR oder Spec-Zeile

| Option | Pro | Contra |
|---|---|---|
| E1 — **kein** ADR; nur `[SPEC-023](../../../spec/pflichtenheft.md)` fortschreiben | kleinster Schritt; die Spec ist Rang 2 und trägt „Sprachen und Umfang" ohnehin | ließe [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) mit einer Umfangs-Festlegung stehen, die die Spec **widerspricht** — ohne ADR, das den Widerspruch auflöst; der `.proto`-Weg (B) und der Stub-Ort (C) sind Entscheidungen, die die Spec nicht tragen kann: eine `.a-check.yml`-Erweiterung mit Import-Berechtigung braucht nach [`ADR-0068`](0068-wegwerf-clients-begrenzte-import-berechtigung.md) Festlegung 3 eine ADR |
| **E2 — Folge-ADR mit `Supersedes` in genau den getroffenen Klauseln (gewählt)** | die abgelösten Klauseln sind **benannt**, der Rest der [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) bleibt **bestätigt** und muss nicht doppelt gelesen werden; die Immutabilitäts-Regel (`AGENTS.md` §3.5) bleibt unangetastet; B/C/D bekommen ihren Träger | eine weitere ADR in der Kette; wer die Sprach-Form sucht, liest zwei Dokumente zusammen (Muster der übrigen engen Klausel-Korrekturen: [`ADR-0070`](0070-supersede-reichweite-und-klassengrenze.md)) |
| E3 — [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) in-place ändern | ein Dokument | verletzt `AGENTS.md` §3.5 an genau der Stelle, an der sie zählt; die Zitat-Korrektur-Ausnahme der `ADR-0073` deckt **nur** das Zitat-/Verweisgerüst — §Entscheidung, §Konsequenzen, §Verglichene Alternativen und §Status sind unberührbar |

### F — Der Schnitt

| Option | Pro | Contra |
|---|---|---|
| F1 — ein Slice je Sprache über **alle vier** Oberflächen | drei Slices, je Sprache einer | je Slice vier Programme, drei Abhängigkeiten und (bei gRPC) ein Generator: über der Drei-Liefer-Punkte-Grenze und in einer Sitzung nicht prüfbar |
| F2 — zwei Slices (je Sprache einer über die drei neuen Oberflächen) | zwei Slices | derselbe Verstoß, nur größer; bündelt die zwei Form-Fragen (Bau-Kontext, Generator) mit zwei Programmen je Sprache |
| **F3 — sieben Slices, die zwei Sprach-Wurzeln zuerst und je eigener Form (gewählt)** | jede Sprach-Wurzel steht **allein** und beweist die Form ([`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md)s eigener Satz); der neue Bau-Kontext und der Generator stehen je in **einem** Slice; jede Zelle ist einzeln lieferbar; ≤ 3 Liefer-Punkte je Slice | sieben Slices für neun Programme — die Zerlegung ist feiner als der Umfang, und die zwei Sprach-Wurzeln blockieren alles Weitere |

**Fazit:** A3, B4, C1, D3, E2, F3. A2/A4 scheitern an der
Nutzerentscheidung bzw. an der Slice-Größe; B1 an der Reproduzierbarkeit des
Baus, B2 am Aufgeben der [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md)-Form
und an einer stillen Ignore-Datei-Klasse, B3 an der zweiten Quelle des
Vertrags, B5 an der Nutzeranweisung; C2/C3 an drei Sync-Gates bzw. an der
stillen Drift; D1/D2 daran, dass sie eine Prüfung behaupten, die gemessen
nichts prüft; E1/E3 an der Trägerfrage bzw. an `AGENTS.md` §3.5; F1/F2 an der
Drei-Liefer-Punkte-Grenze.

## Konsequenzen

- Positiv: Die Nutzerentscheidung ist **vollständig** eingelöst — zwölf
  Programme, vier Oberflächen, drei Sprachen, keine Zelle offen.
- Positiv: **Die Sprach-Form der [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md)
  bleibt stehen.** Bau-Kontext = Sprach-Wurzel, eigenes Dockerfile, digest-
  gepinnte Basis, eigenes Werkzeug-Ziel, kein Gate — die Erweiterung fügt
  **einen** benannten Zusatzkontext hinzu und sonst nichts.
- Positiv: **Der ausgelieferte Bau ist nicht berührt.** Wurzel-`Dockerfile`,
  `.dockerignore`, `make image` und der Coverage-Nenner bleiben; `make gates`
  bleibt netzlos.
- Positiv: **Die vertagte Frage der [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md)
  ist entschieden statt weiter vertagt** — Stub-Ort, `contract`-Gruppe und die
  a-check-Frage haben eine Antwort, und [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  §Re-Evaluierungs-Trigger 2 ist mit „keine Änderung" beantwortet.
- Positiv: **Die `.proto` bleibt die einzige Quelle** des Draht-Vertrags; sie
  wird nicht kopiert, nicht gespiegelt, nicht committet — dreimal gelesen,
  einmal geschrieben.
- Negativ: Es gibt **eine Bauform mehr** im Repo (neben dem Wurzel-Bau und den
  zwei Sprach-Bauten der [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md)
  ein Sprach-Bau **mit** benanntem Zusatzkontext). Sie ist auf die Zellen
  begrenzt, die die `.proto` lesen, und benannt; wer sie zusammenlegt, ändert
  diese Entscheidung.
- Negativ mit Grenze: Die Nicht-Go-Quellen sind mechanisch ungeprüft, was ihre
  Import-Grenzen angeht (Festlegung 6) — Wächter ist das Review; die Grenze ist
  **breiter** als vorher, weil es drei Fremdsprachen statt zwei sind.
- Negativ mit Grenze: Der Träger hängt an einem **nicht-blockierenden**
  Workflow, dessen Umfang wächst und dessen Vollständigkeit nach
  [`AGENTS.md`](../../../AGENTS.md) §3.10 am realen Post-Push-Lauf hängt
  (`BEO-PGC/github-actions-unverifizierbar-lokal`); ein nicht-blockierender
  Lauf kann überlesen werden
  (`BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit`).
- Negativ mit Grenze: „übersetzt und die reinen Funktionen laufen" ist nicht
  „holt am laufenden Feed eine Änderung" — unverändert
  ([`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) Festlegung 4).
- Negativ mit Grenze: Die Sprach-Bauten brauchen **Netz** (Paket-Bezug). Sie
  können deshalb kein Gate werden; der Preis ist, dass ein roter Sprach-Bau
  nur **sichtbar**, nicht **blockierend** ist.
- Folgepflicht (Spec-Zug): `[SPEC-023](../../../spec/pflichtenheft.md)` Zeile
  *Sprachen und Umfang* auf die Matrix, §7 Historie; ohne ADR-/Slice-Kennung
  im Spec-Text.
- Folgepflicht (Bau-Zug, je Sprache mit gRPC): der benannte Zusatzkontext im
  Sprach-Ziel, die Generator-Stufe im Sprach-Dockerfile, die gepinnten
  Plugin-/Laufzeit-Versionen im Manifest.
- Folgepflicht (Handbuch-Zug, je Oberfläche): der `**Beispiele:**`-Block mit
  einer Zeile je Sprache **und** die Zeile in der Änderungshistorie — im
  selben Zug wie das jeweilige Programm (Festlegung 7).
- Folgepflicht (Werkzeug-Zug, je Sprache): die erweiterten Ziele, ihre Zeilen
  in `harness/README.md` §Werkzeuge (**kein** Gate, nicht in `GATE_CHECKS`).
- Folgepflicht (Träger-Zug): der Kommentar der `examples`-Gruppe in
  `.a-check.yml` auf die Go-Aussage samt der erweiterten Grenze.
- Folgepflicht (NATS-Zug, je Sprache): die gepinnte Client-Bibliothek im
  Manifest der Sprache, mit ihrer Leser-Zeile („wofür sie da ist").
- Hinweis: Der **Go**-gRPC-Client bleibt an den Umzug der Vertragsfläche
  gebunden <!-- d-check:status-provenance -->; die drei fremdsprachigen
  gRPC-Clients **nicht** (Festlegung 2, Kontext).

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `make examples-csharp` / `make examples-kotlin` (Werkzeug, **kein** Gate) | Die Beispiele der Sprache werden im digest-gepinnten Werkzeugketten-Image ihres Sprach-Wurzelverzeichnisses **übersetzt** und ihre netzlos prüfbaren Teile **gefahren**; für `grpc` zusätzlich: der Stub wird im Bau aus der `.proto` erzeugt (der Bau bricht ohne den benannten Zusatzkontext ab). Der Exit-Code des Ziels wird direkt gelesen (`AGENTS.md` §3.9) | — (Werkzeug; **nicht** in `GATE_CHECKS`) |
| Go-Toolchain (`go test ./...`, im gepinnten Container) | `examples/grpc-client/` liegt im Wurzelmodul und wird von `make test` übersetzt; ein Kompilierfehler bricht den Lauf | `make test` |
| a-check (Modulzeile `wrong-direction`) | Die Gruppe `examples` prüft **Go**-Importe (ein Import aus `domain`, `ports`, `app` oder `adapters` in `examples/**` ist `wrong-direction`); für Nicht-Go-Quellen trägt sie **nichts**, und die erzeugten Fremdsprachen-Stubs liegen nicht im Baum. Die Gruppe `contract` (`gen/**`) bleibt **Go**; `languages`/`resolution` bleiben unverändert | `make a-check` (im Gate-Bündel) |
| `docker build` (Wurzel) und `make image` | Der ausgelieferte Bau bleibt unberührt: die letzte Stufe bleibt `runtime`, die Wurzel-`.dockerignore` bleibt byte-gleich, die Beispiele liegen in ihrem eigenen Kontext; der Digest-Beleg folgt `ADR-0044` | `make image` (kein Gate) |
| Nicht-blockierender Workflow | Ein PR, der eine Beispiel-Quelle bricht, wird **sichtbar** rot (nicht blockierend); sein Abschluss hängt am realen Post-Push-Lauf (`AGENTS.md` §3.10) | — (kein Gate) |
| Review-Prüfpflicht (nicht maschinell) | Import-Grenzen der Nicht-Go-Quellen (drei Sprachen); ob ein Beispiel noch zeigt, was es zeigen soll; die Verdrahtung über die reine Funktion hinaus; ob die Sprach-Bauten weiterhin die **einzige** Quelle der `.proto` lesen (keine Kopie im Baum) | — |

## Slice-Schnitt-Empfehlung

Regeln dieser Sektion: Empfehlung, endgültiger Schnitt liegt bei der
Planner-Rolle (Baseline-Regelwerk `modul-05-planning-harness.md`).

**Sieben Slices**, je ≤ 3 Liefer-Punkte und ≤ 2 berührte Schichten
(`examples/**` bzw. Wurzel-Tooling + Doku). Die zwei Form-Träger zuerst: eine
Sprache, die zum ersten Mal auftritt, **steht allein** — sie beweist die Form
([`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md)s eigener Satz für den
ersten Zug einer Sprache).

| # | Slice | Liefer-Punkte | Abhängigkeit |
|---|---|---|---|
| 1 | **C#-Sprach-Wurzel + HTTP-Client** — `examples/csharp/Dockerfile` mit digest-gepinnter Basis und gepinntem Projekt-/Paket-Manifest, das `make`-Ziel `examples-csharp`, das Programm mit netzlosen Tests, die Träger (Handbuch-`**Beispiele:**`-Block der HTTP-Familie mit der C#-Zeile, README §Werkzeuge, nicht-blockierender Workflow) | LP1 Wurzel + Werkzeugkette + Ziel; LP2 der Client + seine netzlosen Tests; LP3 die Träger samt Workflow | keine |
| 2 | **Kotlin-Sprach-Wurzel + HTTP-Client** — dieselbe Form, eigene Werkzeugkette (Compiler/Projekt-Werkzeug, Versions- und Prüfsummen-Pin), Ziel `examples-kotlin` | wie 1 | keine (parallel zu 1) |
| 3 | **SSE-Client in C# und Kotlin** — beide Programme, ihre netzlosen Tests (Zerlegen eines SSE-Frames), die zwei Handbuch-Zeilen | LP1 C#-Client + Tests; LP2 Kotlin-Client + Tests; LP3 die zwei Handbuch-Zeilen | 1, 2 |
| 4 | **NATS-Client in C# und Kotlin** — beide Programme samt zweiseitigem Ablauf, ihre Tests (Subjekt-Aufbau), die gepinnten Client-Bibliotheken (`NATS.Net`, `io.nats:jnats`) mit Leser-Zeile, die zwei Handbuch-Zeilen | LP1 C#-Client + Pin + Tests; LP2 Kotlin-Client + Pin + Tests; LP3 die zwei Handbuch-Zeilen | 1, 2 |
| 5 | **gRPC-Client in C#** — der **benannte Zusatzkontext** im C#-Ziel, die Generator-Stufe im C#-Dockerfile (`Grpc.Tools` samt gepinnten Laufzeit-Paketen), das Programm, seine netzlosen Tests, die Handbuch-Zeile | LP1 der `.proto`-Weg im Bau (Zusatzkontext + Generator; Beleg: der Bau erzeugt den Stub); LP2 der Client + Tests; LP3 die Handbuch-Zeile | 1 |
| 6 | **gRPC-Client in Kotlin** — dieselbe, jetzt **gemessene** Form auf den Kotlin-Bau übertragen (`protoc-gen-grpc-kotlin`, `grpc-kotlin-stub`, `grpc-netty-shaded`), das Programm, seine Tests, die Handbuch-Zeile | wie 5 | 2, 5 (die Form ist in 5 belegt) |
| 7 | **Go-gRPC-Client** (`examples/grpc-client/`) — das Programm mit Tests, die Kante `examples → contract` **mit ihrem Objekt**, die Go-Zeile im gRPC-Handbuch-Absatz | LP1 der Client + Tests; LP2 die Kante in `.a-check.yml`; LP3 die Handbuch-Zeile | **der Umzug der Vertragsfläche** (heute in `next/`) <!-- d-check:status-provenance --> — und **nichts** sonst |

**Abhängigkeiten, explizit:** 3 und 4 setzen 1 und 2 voraus (ohne
Sprach-Werkzeugkette kein Programm); 5 setzt 1 voraus, 6 setzt 2 und die in 5
belegte Form voraus; 7 hängt **allein** am Umzug und ist von 1–6 unabhängig.
3 und 4 sind voneinander unabhängig, 5/6 und 3/4 ebenfalls — die Reihenfolge
ist eine Reihenfolge, keine Kette.

**Was der Planner beim Anlegen beachten muss:**
- Die zwei Form-Fragen (Bau-Kontext, Generator) **stehen je in genau einem**
  Slice (5 für C#, 6 für Kotlin) — nicht in 1/2, deren Form die
  [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) schon entschieden hat.
- Der Umzug trägt **eine** der neun Zellen <!-- d-check:status-provenance -->;
  wer ihn aus dem Go-gRPC-Slice heraushält, hält ihn aus **allen** anderen
  ebenfalls heraus (er verbaut keinen von ihnen — Festlegung 2).
- Die `.proto` ist gemeinsame Berührungsfläche (Festlegung 3): wer sie
  für Codegen-Optionen anfasst, nennt den Umzug-Slice.

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Vier benannte Trigger, sonst permanent:

1. **Eine vierte Fremdsprache wird angefragt.** Dann gilt unverändert
   [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) §Re-Evaluierungs-Trigger 1
   (trägt die Form noch? wandert die Werkzeugketten-Frage in eine gemeinsame,
   sprachneutrale Form?) — und zusätzlich ist zu prüfen, ob der benannte
   Zusatzkontext (Festlegung 2) mit der Zahl der Sprachen noch trägt.
2. **Ein fremdsprachiger gRPC-Bau kann seinen Stub nicht mehr aus der `.proto`
   erzeugen** — weil der Vertrag eine Einfuhr aus einem **anderen**
   Verzeichnis verlangt (etwa eine gemeinsame Typ-Datei) oder weil die
   Codegen-Optionen den Vertrag selbst berühren müssten. Dann ist **diese**
   Bauform neu zu bewerten, nicht still eine zweite Quelle anzulegen.
3. **Eine Sprach-Werkzeugkette verlangt eine nicht-öffentliche Quelle**
   (privater Feed, Vendor-Verzeichnis, Zugangsdaten in der Auflösung). Dann ist
   Festlegung 4 verletzt und der Client der Sprache ist neu zu schneiden —
   eine kopierbare, nachbaubare Quelle ist die Bedingung des Beispiels, nicht
   sein Beiwerk.
4. **Der nicht-blockierende Workflow trägt seinen Umfang nicht mehr** (Laufzeit
   über dem Runner-Budget, oder ein dauerhaft roter, überlesener Lauf). Dann
   ist der Schnitt des Workflows (je Sprache ein eigener? blockierend für
   welche Hälfte?) neu zu entscheiden — dieselbe Frage, die
   `BEO-PGC/nicht-blockierender-workflow-alarmmuedigkeit` und
   `BEO-PGC/github-actions-unverifizierbar-lokal` im Register führen.

Sonst permanent: `examples/` führt die vier Zugriffs-Oberflächen in Go, C# und
Kotlin; jede Sprache hat ihr Sprach-Wurzelverzeichnis, ihre digest-gepinnte
Werkzeugkette und ihr Werkzeug-Ziel; ein fremdsprachiger gRPC-Bau liest die
`.proto` über **einen** benannten Zusatzkontext; die erzeugten
Fremdsprachen-Stubs liegen nicht im Baum; `gen/**` bleibt die Go-Bindung.

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-17 | Accepted — Anlass: Nutzerentscheidung vom 2026-09-17 („Beachte matrix: (Alle clients) × (c#, go, kotlin)" samt „Auch grpc-Streaming"), die den Umfang der [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) von der HTTP-Familie auf die volle Matrix hebt. Unabhängiger Architect-Zug entscheidet Umfang (neun neue Programme), den `.proto`-Weg der fremdsprachigen gRPC-Bauten (ein benannter Zusatzkontext — real gemessen), den Ort der erzeugten Stubs (im Bau, nicht committet), die Abhängigkeits-Pins (Registry-Abfragen je Bibliothek), die unveränderte Go-Aussage der `examples`-Gruppe, die Form des Handbuch-Nachzugs und den Schnitt (sieben Slices); supersedes vier Klauseln der [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) und beantwortet [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md) §Re-Evaluierungs-Trigger 2 mit „keine Änderung"; Modul 8 §Rollen-Regeln: „Architect schreibt" | gemessen in diesem Zug: der `examples/`-Bestand, das fehlende `gen/**` und die fehlenden C#-/Kotlin-Sprach-Wurzeln; die benannte-Kontext-Probe (`docker buildx build --build-context …` mit `COPY --from=`, `buildx` v0.37.1); die Registry-Abfragen zu `NATS.Net`, `Grpc.Net.Client`, `Grpc.Tools`, `Google.Protobuf`, `System.Net.ServerSentEvents`, `io.nats:jnats`, `grpc-kotlin-stub`, `protoc-gen-grpc-kotlin`, `grpc-netty-shaded`, `okhttp-sse`; die Existenz der zwei Sprach-Basen (`docker manifest inspect`) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0090` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
