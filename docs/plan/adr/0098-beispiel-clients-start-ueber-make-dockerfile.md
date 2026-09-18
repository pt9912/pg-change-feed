# ADR-0098: Beispiel-Clients — Startform über `make`/Dockerfile (Go), echter Start-Make-Target über alle drei Sprachen, Umgebungsdatei-Kontrakt

**Status:** Accepted — Supersedes [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
in **genau einer** Klausel: der Startform-Bullet ihrer §Entscheidung
Festlegung 1 („**Startform (Form-Vorgabe, Detail der umsetzende Slice):** je
Programm ein `go run ./examples/<name>` … Damit zitiert das Handbuch **eine**
Befehlsform, die innerhalb und außerhalb des Compose-Netzes funktioniert.").
Alles Übrige der [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
bleibt **hiermit bestätigt** und wird nicht wiederholt: Ort und Verzeichnisform
(§1, restliche Bullets), Import-Grenze und ihr Träger (§2), der gRPC-Baustein
(§3), die `.a-check.yml`-Gruppen (§4), der Kompilierpfad als
Anti-Verrottungs-Träger (§5), das Verhältnis zu den Wegwerf-Clients (§6), das
Handbuch-Zitat (§7). Ebenfalls unberührt:
[`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) und
[`ADR-0090`](0090-beispiel-clients-volle-matrix.md) — ihre Sprach-Wurzel-,
Bau-Kontext- und `.proto`-Weg-Festlegungen für C#/Kotlin bleiben, was sie
sind; diese ADR ergänzt sie um einen **Start**-Mechanismus, ohne ihren
**Bau**-Mechanismus zu ändern.

**Datum:** 2026-09-18

**Autor:** pt9912 (Architect-Rolle, unabhängiger Architect-Zug auf explizite
Auftraggeber-Anweisung vom 2026-09-18 — „das habe ich nicht entschieden.
Bitte make und Dockerfile verwenden!" — sowie die Anweisung, allen drei
Sprachen einen echten Start-Make-Target und `examples/` eine
Demo-Umgebung mit Bootstrapping zu geben; anderer Kontext als der
Planner-Lauf, der `welle-beispiele-start-ueber-make` geschnitten hat, und als
die Implementer-Läufe, die die zwölf Beispiel-Programme gebaut haben — Modul 8
§Rollen-Regeln: „Architect schreibt")

**Bezug:** [`LH-FA-SST-006`](../../../spec/lastenheft.md),
[`LH-FA-SST-007`](../../../spec/lastenheft.md),
[`LH-FA-SST-008`](../../../spec/lastenheft.md) (Scope: die Oberflächen, die
die zwölf Beispiele zeigen — diese ADR ändert sie nicht),
[`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
(in der Startform-Klausel superseded; sonst bestätigt),
[`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md),
[`ADR-0090`](0090-beispiel-clients-volle-matrix.md) (Bestand: Sprach-Wurzel,
Bau-Kontext, `.proto`-Weg, `make examples-csharp`/`make examples-kotlin` als
Bau-/Test-Ziel — unberührt), [`ADR-0044`](0044-image-beleg-semantik.md)
(Image-Beleg-Semantik — der neue Go-Bau ist ein zweiter, unabhängiger
Docker-Bau und berührt `make image`/`harness/image-hash.txt` nicht),
[`ADR-0085`](0085-build-kontext-ausnahme-test-only-zweck.md)
(Build-Kontext-Ausnahme-Klasse — **nicht** einschlägig, siehe Kontext),
[`AGENTS.md`](../../../AGENTS.md) §3.1/§3.6/§3.9/§3.13,
[`SPEC-023`](../../../spec/pflichtenheft.md) (Beispiel-Clients — von dieser
ADR **geschärft**), [`SPEC-016`](../../../spec/pflichtenheft.md)
(`CDC_CONFIG_FILE` — **nicht** berührt, siehe Festlegung 3),
`examples/**`, `examples/Dockerfile` (neu), `examples/Dockerfile.dockerignore`
(neu), `examples/.env` (neu), `harness/mk/examples.mk`,
`docs/user/benutzerhandbuch.md` §4, `examples/README.md`,
`docs/plan/planning/open/slice-beispiele-compose-bootstrap.md`,
`docs/plan/planning/open/slice-beispiele-go-dockerfile-start.md`,
`docs/plan/planning/open/slice-beispiele-csharp-kotlin-start-target.md`

**Schärft:** [`SPEC-023`](../../../spec/pflichtenheft.md) — Zeile
„Laufzeit und Startform" wird auf **alle drei** Sprachen einheitlich
Container-Aufruf gezogen (bislang bereits so formuliert, aber durch
[`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)s
`go run`-Festlegung für Go faktisch unterlaufen — siehe Kontext, gemessener
Fund); die Spec bekommt zusätzlich eine Zeile für den
Umgebungsdatei-Kontrakt und die Start-Make-Ziele. Das Pflichtenheft trägt die
Festlegung (Rang 2), diese ADR die Entscheidung und ihre Begründung; bei
Widerspruch gilt der höhere Rang.

**Regeln:** Baseline-Regelwerk `modul-04-adrs.md` §Ziel-Form: ADR (MADR).

---

## Kontext

**Die Auftraggeber-Anweisung — wörtlich, in zwei Teilen.** Am 2026-09-18: (1)
„das habe ich nicht entschieden. Bitte make und Dockerfile verwenden!" —
explizite Zurückweisung der `go run`-Startform, die
[`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
§Entscheidung Festlegung 1 für die Go-Beispiele gesetzt hatte; (2) alle drei
Sprachen bekommen einen **echten** Start-Make-Target (nicht nur Bau/Test),
und `examples/` bekommt eine Demo-Compose-Umgebung mit Bootstrapping und
„das config-file".

**Eine `Accepted`-ADR wird nicht überschrieben** (`AGENTS.md` §3.5) — die
Umkehr der Go-Startform ist deshalb diese neue ADR mit `Supersedes ADR-0076`
in genau der einen Klausel, nicht eine Korrektur der bestehenden Datei.

**Der Ist-Stand — gemessen, nicht erinnert.**

| Prüfung | Ergebnis |
|---|---|
| `examples/README.md`, Go-Tabelle | Startform je Programm ist `go run ./examples/<name> ...` — kein Bau-Schritt, kein Image |
| `examples/README.md`, C#/Kotlin-Tabellen | Bau über `make examples-csharp`/`make examples-kotlin` (Docker, digest-gepinnt); **Start** ist Fließtext: „Start je Image per `docker run --rm <ENV-Variablen> <Image> <Argumente>`" — kein `make`-Ziel |
| `harness/mk/examples.mk` | trägt `examples-csharp`/`examples-kotlin` — beide bauen **und testen**, keines startet |
| `grep -rn "example-run" Makefile harness/mk/` | kein Treffer — kein Start-Ziel existiert in irgendeiner Sprache |
| `compose.yaml` (Wurzel) | Lauf-/CI-Vertrag-Seite (`LH-QA-POR-003`); Kopfkommentar nennt sie ausdrücklich als Testlauf-Umgebung, kein Leser-Erzeugnis; **keine** `examples/`-Compose-Datei existiert |
| `spec/pflichtenheft.md` §2 `SPEC-023`, Zeile „Laufzeit und Startform" | „Docker-only … Jede Sprache wird aus einem **digest-gepinnten** Basis-Image gebaut … die Startform ist ein **Container-Aufruf**" — **sprachneutral formuliert**, ohne Go-Ausnahme; das steht im Widerspruch zum tatsächlichen Go-`go run`-Zustand, den [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md) gesetzt hat. Diese ADR löst den Widerspruch zugunsten der bereits geltenden Spec-Zeile auf, nicht umgekehrt — ein zusätzlicher, unabhängig gefundener Grund für die Auftraggeber-Entscheidung, keine Rückwirkungs-Fiktion |
| `spec/pflichtenheft.md` §2 `SPEC-023`, Zeile „Verhältnis zur Konfigurationsdatei" | „Beispiele lesen `CDC_CONFIG_FILE` **nicht** … Die Konfigurationsdatei ist der Deployment-Eingang des Feeds, nicht der eines Integrator-Programms" — **bleibt unverändert geltend** (Festlegung 3 unten löst „das config-file" der Auftraggeber-Anweisung als **andere** Sache auf) |
| `docs/user/benutzerhandbuch.md` §5.1 | führt die vollständige `CDC_*`-Variablenliste; jeder `**Beispiele:**`-Block der vier Zugriffs-Abschnitte zitiert bereits dieselben Namen (`CDC_HTTP_ADDR`, `CDC_GRPC_ADDR`, `CDC_API_TOKEN_READER`, `CDC_NATS_URL`) |
| `gen/cdc/stream/v1/` | existiert (Umzug aus `slice-097` <!-- d-check:status-provenance --> ist committet) — der Go-gRPC-Client importiert bereits den öffentlichen Vertragspfad, kein offener Vorgänger mehr |

**Fünf Bindungen prägen den Lösungsraum — alle fünf gemessen.**

1. **Der ausgelieferte Bau ist unantastbar**
   ([`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) Festlegung 3): die
   Wurzel-`Dockerfile`, ihre fünf Stufen und die Wurzel-`.dockerignore` bauen
   `cmd/pg-change-feed`; ein Go-Beispiel-Bau darf ihren Cache-Verhalten und ihr
   Digest-Verhalten nicht berühren.
2. **Ein zweites `go.mod` ist verboten** ([`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
   Festlegung 5): der Anti-Verrottungs-Träger der Go-Beispiele ist der
   bestehende Kompilierpfad `go test ./...` im **Wurzelmodul** — „die
   Beispiele dürfen aus `go test ./...` nicht herausfallen (kein zweites
   Modul, kein Ausschluss)". Diese ADR fügt einen **zweiten, unabhängigen**
   Docker-Bauweg hinzu; sie ersetzt den bestehenden Kompilierpfad nicht, und
   der bleibt unverändert grün.
3. **Die Docker-Kontext-Ausnahme-Klasse aus [`ADR-0085`](0085-build-kontext-ausnahme-test-only-zweck.md)
   passt hier nicht**, weil sie eine **einzelne Datei** trägt (Merkmal ii),
   nicht ein Verzeichnis. Ein Go-Beispiel-Bau, der `go.mod`/`go.sum`/`gen/`
   liest, verlangt also einen anderen Mechanismus als eine `!datei`-Negation
   in der Wurzel-`.dockerignore`.
4. **Docker-only** (`AGENTS.md` §3.1): kein Host-`go build` für die
   Beispiele — auch der Go-Startweg läuft über ein Image, nicht über
   `go run` auf dem Host.
5. **Eine Arbeit, die eine beschriebene Eigenschaft bewegt, zieht ihre Träger
   nach** (`AGENTS.md` §3.13): `SPEC-023`s „Laufzeit und Startform"-Zeile
   und ihr Go-Zweig in `examples/README.md`/`docs/user/benutzerhandbuch.md`
   werden mit dieser Entscheidung falsch, wenn sie nicht nachgezogen werden.

### Die Kernfrage dieses Zugs: wie erreicht ein Go-Beispiel-Bau `go.mod`/`gen/`, ohne den ausgelieferten Bau zu berühren?

Der C#/Kotlin-Weg (eigenes Dockerfile, Kontext = Sprach-Wurzelverzeichnis,
[`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) Festlegung 3) trägt für
Go **nicht direkt**: Go-Beispiele sind Teil des **Wurzelmoduls** (Bindung 2)
und brauchen `go.mod`/`go.sum`/`gen/` — Dateien, die außerhalb eines
Sprach-Wurzelverzeichnisses `examples/` liegen. Ein Kontext = `examples/`
allein sieht diese Dateien nicht; ein Kontext = Repo-Wurzel bräuchte eine
Ignore-Anpassung, die — naiv als `!examples/` in der Wurzel-`.dockerignore`
gesetzt — **jede** `COPY . .`-Stufe der Wurzel-`Dockerfile` (`coverage`,
`build`) vergrößerte und deren Digest-Verhalten auf Beispiel-Änderungen
reagieren ließe (dieselbe Klasse, die
[`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) §Verglichene
Alternativen Option B2 für C#/Kotlin verworfen hat).

**Der Weg, der das löst, ist real gemessen, nicht angenommen** (dieser Zug,
`docker build` v29.8.0, Wegwerf-Baum unter dem Scratchpad-Verzeichnis dieses
Zugs — kein Repo-Pfad):

| Probe | Ergebnis |
|---|---|
| Wurzel-`.dockerignore` trägt `*` mit Ausnahmen für `cmd/`/`internal/`/`gen/`/`go.mod`, **kein** `examples/`-Eintrag. Bau mit `-f examples/Dockerfile .` (Kontext = Wurzel, **kein** dateispezifisches Ignore) | Wurzel-`.dockerignore` gilt (`load .dockerignore: 73B` — dieselbe Byte-Zahl wie ohne Zusatzdatei); `examples/`-Inhalt fehlt im Bau-Kontext, `cmd/`/`internal/`/`gen/`/`go.mod` sind da |
| Dieselbe Wurzel-`.dockerignore`, **zusätzlich** eine Datei `examples/Dockerfile.dockerignore` (Inhalt: `*` + `!go.mod`/`!gen/`/`!examples/`) neben dem Dockerfile, Bau weiterhin mit `-f examples/Dockerfile .` | **Die dateispezifische Ignore-Datei trägt Vorrang**: der Bau-Kontext enthält jetzt `go.mod`, `gen/`, `examples/` — und **nicht** `cmd/`/`internal/` (dort keine `!`-Zeile in der neuen Datei) |
| Derselbe Baum, Bau mit `-f Dockerfile .` (der **Wurzel**-Dateiname, nicht `examples/Dockerfile`) | Die dateispezifische Datei `examples/Dockerfile.dockerignore` **greift nicht** — der Bau-Kontext ist identisch mit der ersten Probe (`cmd/`/`internal/`/`gen/`/`go.mod`, kein `examples/`) |

Heißt: eine **dateispezifische** `.dockerignore` (`<Dockerfile>.dockerignore`,
neben dem jeweiligen Dockerfile, dokumentierter Docker-Mechanismus) wirkt
**ausschließlich** für Bauten, die genau diese Dockerfile-Datei über `-f`
referenzieren. Die Wurzel-`Dockerfile` und ihre `.dockerignore` bleiben davon
**vollständig unberührt** — ein stärkerer, weil **byte-gleich** messbarer
Nachweis als die C#/Kotlin-Form, die für die `.proto`-Frage einen
zusätzlichen, benannten Bau-Kontext einführt
([`ADR-0090`](0090-beispiel-clients-volle-matrix.md) Festlegung 2): hier
genügt eine zweite Ignore-Datei, weil Go — anders als C#/Kotlin — schon immer
mit Kontext = Repo-Wurzel gebaut wird (die Wurzel-`Dockerfile` selbst tut das
für `cmd/pg-change-feed`).

## Entscheidung

Wir wählen: **Die vier Go-Beispiele bekommen ein eigenes
`examples/Dockerfile` mit Bau-Kontext Repo-Wurzel, isoliert über eine
dateispezifische `examples/Dockerfile.dockerignore` (Wurzel-`Dockerfile` und
Wurzel-`.dockerignore` bleiben byte-gleich); Startform ist ein
Container-Aufruf wie bei C#/Kotlin. Alle drei Sprachen bekommen je einen
`make`-Start-Target (`example-run-go`, `example-run-csharp`,
`example-run-kotlin`), Pflicht-Argument `SURFACE=`, der einen
digest-gepinnten Client-Container real gegen eine feste, benannte
Docker-Netzwerk-Umgebung startet. Eine gemeinsame, committete Umgebungsdatei
`examples/.env` trägt die bereits im Handbuch dokumentierten `CDC_*`-Namen
als Compose-Netzwerk-Adressen (nicht `localhost`) und wird sowohl von der
Demo-Compose-Datei (`env_file:`) als auch von jedem Start-Target
(`--env-file`) gelesen.**

Fünf Festlegungen.

### 1 — Go-Bauform: eigenes Dockerfile, Wurzel-Kontext, dateispezifisches Ignore

- **Neue Datei `examples/Dockerfile`** (Repo-Wurzel-Kontext: der Bau läuft als
  `docker build -f examples/Dockerfile .` — der **Pfad** liegt unter
  `examples/`, der **Kontext** ist die Repo-Wurzel, weil Go-Beispiele
  Wurzelmodul-Code sind, Bindung 2). Multi-Stage nach demselben Muster wie
  die C#/Kotlin-Dockerfiles: eine gemeinsame `build`-Stufe kompiliert alle
  vier Programme, vier schlanke `runtime*`-Stufen liefern je ein startbares
  Image — `runtime` (letzte Stufe, HTTP-Client, mirror der
  [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md)-Konvention: `docker
  build` ohne `--target` liefert das HTTP-Image), `runtime-sse`,
  `runtime-grpc`, `runtime-nats`.
- **Isolation über eine dateispezifische Ignore-Datei, nicht über eine
  Wurzel-`.dockerignore`-Negation** (real gemessen, Kontext): `examples/Dockerfile.dockerignore`
  trägt `*` plus `!go.mod`, `!go.sum`, `!gen/`, `!examples/`. Diese Datei
  wirkt **ausschließlich**, wenn ein Bau `-f examples/Dockerfile` referenziert
  — ein Bau mit `-f Dockerfile` (der ausgelieferte Bau, `make image`) sieht
  sie nicht und bleibt exakt beim bisherigen, byte-gleichen Kontext.
  **Die Wurzel-`.dockerignore` wird nicht angefasst.**
- **Digest-Wiederverwendung, kein neuer Pin.** Beide Basis-Images sind bereits
  gepinnt und `Accepted` (Wurzel-`Dockerfile`): `golang:1.27-alpine@sha256:cf6fca6641884b8433441b2b0652976f975e1d0fdd26d177eaaf8596087f3125`
  (Build-Stufe) und `gcr.io/distroless/static-debian12:nonroot@sha256:afa5c872c891853ca7fcf1f12c3edb23f7eeef36189728842dd51042ff57f7ab`
  (Runtime-Stufen) — dieselben zwei Digests, wiederverwendet statt neu
  gemessen. Ein Pin-Hebungs-Trigger dieser Datei fällt automatisch mit dem
  der Wurzel-`Dockerfile` zusammen (`make image-stale`, unverändert).
- **Der Kompilierpfad `go test ./...` bleibt der Anti-Verrottungs-Träger**
  ([`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  Festlegung 5, unverändert bestätigt) — `examples/Dockerfile` ist ein
  **zweiter, unabhängiger** Bauweg für die **Startform**, kein Ersatz des
  Kompilierpfads und kein zweites `go.mod` (Bindung 2 bleibt erfüllt).
- **Warum kein Umzug auf das C#/Kotlin-Sprach-Wurzel-Muster** (die im
  Welle-Text offene Frage a): Go-Beispiele bleiben **flach**
  (`examples/<name>/`, [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md)
  Festlegung 2 „Go verlangt das nicht") — ein Umzug nach `examples/go/<name>`
  wäre Arbeit am Bestand ohne Adresse, weil Go weder ein
  Projektverzeichnis noch einen SDK/Runtime-Split braucht (anders als
  C#/Kotlin, deren Werkzeugketten das verlangen). Die **leichtere,
  Go-spezifische Bauform** ist gewählt: eine Datei, ein Kontext, vier
  Ziel-Stufen, keine neue Verzeichnisebene.

### 2 — Start-Make-Target: ein Ziel je Sprache, Pflicht-Argument `SURFACE=`

- **Drei neue `.PHONY`-Ziele** in `harness/mk/examples.mk`:
  `example-run-go`, `example-run-csharp`, `example-run-kotlin` — je eines
  bekommt als **Pflicht-Argument** `SURFACE=http|sse|grpc|nats` (Make-Variable
  auf der Kommandozeile, `make example-run-go SURFACE=sse ARGS="..."`); ein
  fehlendes oder unbekanntes `SURFACE` bricht mit `$(error …)` und einer
  lesbaren Meldung ab, **bevor** ein `docker run` versucht wird.
- **Warum ein Ziel je Sprache statt je Programm oder eines einzigen
  parametrisierten Ziels über alle drei** (die im Welle-Text offene Frage b,
  drei verglichene Formen unten): Ein Ziel je Programm wären **zwölf**
  nahezu identische `.PHONY`-Stanzas über drei `mk`-Fragmente — Wiederholung
  ohne Nutzen. Ein einziges, sprachübergreifendes Ziel müsste intern
  verzweigen (Go: Binary/Image-Namensschema ohne SDK-Suffix; C#/Kotlin:
  Image-Tags `pg-change-feed-examples:<sprache>[-<surface>]`) — dieselbe
  Verzweigung, nur in einer Datei statt in dreien, und ohne den Vorteil, den
  ein `mk`-Fragment je Sprache heute schon hat (jede Sprache **besitzt** ihr
  eigenes Ziel, siehe `examples-csharp`/`examples-kotlin`). **Ein Ziel je
  Sprache mit Pflicht-Argument** trägt dieselbe Granularität wie die
  bestehenden Bau-Ziele (Bau ist je Sprache **alle vier** Programme; Start
  ist je Sprache **ein** Programm zur Zeit, gewählt über `SURFACE=`) und
  folgt einem bereits im Repo gelebten Idiom (Variable auf der
  Kommandozeile überschreibt einen Default, wie `RANGE=` bei
  `make commit-traceability`, `harness/README.md` §Traceability rules).
- **Der Mechanismus je Sprache — ein Container-Aufruf, keine
  Host-Toolchain** (`AGENTS.md` §3.1):
  - `example-run-go`: baut (falls nötig) `pg-change-feed-examples:go[-<surface>]`
    aus `examples/Dockerfile` (Festlegung 1) und startet ihn.
  - `example-run-csharp`/`example-run-kotlin`: starten die **bereits** von
    `make examples-csharp`/`make examples-kotlin` gebauten Image-Tags
    (`pg-change-feed-examples:csharp[-<surface>]`/`:kotlin[-<surface>]`,
    [`ADR-0090`](0090-beispiel-clients-volle-matrix.md)) — kein neuer Bau,
    nur ein neuer Start-Aufruf gegen ein bestehendes Erzeugnis.
  - Jeder Aufruf ist strukturell `docker run --rm --network <fester
    Netzname, Festlegung 4> --env-file examples/.env <Image-Tag> $(ARGS)`;
    `$(ARGS)` trägt die Flag-Übersteuerung, die
    [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
    Festlegung 1 bereits für alle Beispiele vorsieht.
- **`AGENTS.md` §3.9 gilt unverändert:** der Exit-Code jedes Ziels wird
  direkt gelesen; ein Start-Ziel ist kein Gate und läuft nicht in
  `GATE_CHECKS` — es ist Werkzeug, wie `examples-csharp`/`examples-kotlin`
  ([`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) Festlegung 4,
  unverändert: netzloses Gate-Bündel, Paket-/Image-Bezug bräuchte Netz).
- **`go run ./examples/<name>` bleibt technisch möglich** (Go verbietet es
  nicht, und es bricht nichts), ist aber **nicht mehr die im Handbuch
  zitierte Startform** — die Auftraggeber-Anweisung war genau diese
  Verschiebung, keine Löschung einer Fähigkeit.

### 3 — Umgebungsdatei-Kontrakt: `examples/.env`, committet, bestehende `CDC_*`-Namen

- **Bestätigung der Planner-Lesart, mit einer Korrektur:** die gemeinsame
  Umgebungsdatei ist tatsächlich eine `--env-file`/`env_file:`-Datei, die
  **sowohl** die neue Demo-Compose-Datei **als auch** jedes Start-Target
  liest (Planner-Lesart in der Welle-Datei §1, bestätigt). Die Korrektur:
  **Dateiname ist `examples/.env`, nicht `examples/.env.example`** — sie
  wird **direkt committet und benutzt**, nicht als Vorlage zum Kopieren. Der
  Grund: diese Datei trägt **keine echten Zugangsdaten** — sie gilt
  ausschließlich innerhalb des isolierten Demo-Docker-Netzwerks
  (Festlegung 4) und enthält dieselbe Klasse von Demo-Werten, die
  `compose.yaml` bereits committet führt (`CDC_API_TOKEN_READER:
  e2e-reader-token`, Klartext, Zeile 108 dieser Datei). Ein Kopier-Schritt
  (`.env.example` → `.env`) wäre eine Ausweich-Zeremonie ohne Sicherheitswert
  und widerspräche der Welle-Closure-Anforderung „ohne dass der Leser das
  von Hand einrichtet" (Welle-Datei §1 Punkt 3).
- **Die Variablennamen sind die bereits im Handbuch dokumentierten `CDC_*`-Namen**
  — **keine** neuen, Demo-spezifischen Namen. Das schließt das in der
  Welle-Datei §6/Slice-Risiko benannte Kollisions-Risiko: Es ist **keine**
  Kollision, sondern die **Absicht** — dieselbe Vokabel für Demo und
  Produktion ist der Punkt (ein Integrator liest `CDC_HTTP_ADDR` einmal, im
  Handbuch, in `examples/.env` und in seiner eigenen Produktionsumgebung).
  Die Trennung von Demo und Produktion liegt **nicht** im Namen, sondern in
  der Netzwerk-Isolation (Festlegung 4) und im Datei-Kopfkommentar (siehe
  unten). Damit ist das Risiko aus
  `slice-beispiele-compose-bootstrap` §6 Punkt 2 aufgelöst: **kein
  Sicherheits-Anti-Pattern**, solange der Kommentar die Grenze nennt.
- **`examples/.env` trägt einen Klartext-Kopfkommentar**: „Demo-Zugangsdaten,
  gültig ausschließlich im isolierten Docker-Netzwerk `cdc-examples`
  (Festlegung 4) — niemals gegen eine Produktionsinstanz verwenden."
- **Warum das *nicht* die von [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md)
  Festlegung 5 verbotene `CDC_CONFIG_FILE`-Lesung ist — explizit
  festgestellt, weil dieselbe Zeile dieser ADR sonst als stillschweigend
  aufgeweicht gelesen werden könnte:** `CDC_CONFIG_FILE` ist eine
  **YAML-Datei, die der Feed-Prozess selbst zur Laufzeit parst**
  (`SPEC-016`) — ein zweiter Konfigurations-Eingang **innerhalb** des
  Programms. `examples/.env` ist das **Gegenteil**: eine Datei, die **die
  Shell/Docker/Compose** vor dem Programmstart liest und daraus **dieselben**
  Umgebungsvariablen setzt, die jedes Beispiel schon immer über `os.Getenv`
  bzw. seine Runtime-Umgebung liest. Kein Beispiel-Programm öffnet
  `examples/.env` selbst; kein neuer Lesepfad entsteht **im** Programm. Diese
  ADR **berührt** [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md)
  Festlegung 5 damit **nicht** und braucht dafür **keinen** Supersedes-Eintrag
  — die Klausel bleibt vollständig in Kraft.
- **Adressform: Compose-Netzwerk-Servicenamen, nicht `localhost`** — siehe
  Festlegung 4. `examples/.env` trägt z. B. `CDC_HTTP_ADDR=pg-change-feed:8090`,
  nicht `CDC_HTTP_ADDR=localhost:8090`.
- **Was diese Festlegung nicht entscheidet:** die konkreten Werte für
  Quelle/Tabelle/Publikation der Demo-Daten (welche Tabelle die
  Beispiel-Registrierung anlegt) — das ist Sache von
  `slice-beispiele-compose-bootstrap`, der den Kontrakt (Dateiname,
  Variablennamen, Leser) hier übernimmt, aber die Werte befüllt.

### 4 — Das feste Docker-Netzwerk der Demo-Umgebung

- **Ein fest benannter Docker-Netzwerk-Name, `cdc-examples`** (Analogie zum
  bestehenden Muster: `compose.yaml` trägt `networks: cdc-test: name:
  cdc-feed-test` mit Kommentar „Fester Netzname: `make schema-rollout` trägt
  ihn … ohne vom Compose-Projektnamen abzuleiten"). Die Demo-Compose-Datei
  (`slice-beispiele-compose-bootstrap`) legt dieses Netzwerk mit **genau
  diesem** Namen an; jedes `example-run-*`-Ziel hängt seinen
  `docker run --network cdc-examples`-Aufruf daran.
- **Warum Netzwerk-Beitritt statt Port-Publizierung:** ein `docker run
  --network cdc-examples` erreicht den Feed-Container über seinen
  Compose-Servicenamen (`pg-change-feed`), genau wie es
  `tools/harness/{httpclient,sseclient,grpcclient}` heute gegen die
  CI-Compose-Umgebung tun (`docker exec`/Netzwerk-Alias-Muster,
  `harness/README.md` §Sensors, `make test-integration`). Eine
  Host-Port-Publizierung (`ports: ["8090:8090"]` in der Demo-Compose-Datei)
  bräuchte zusätzlich Konflikt-Vermeidung mit anderen lokal laufenden
  Diensten und zwei Adressformen (`localhost:<port>` außerhalb, Servicename
  innerhalb) — genau die Doppeldeutigkeit, die
  [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  Festlegung 1 mit „**eine** Befehlsform, die innerhalb und außerhalb …
  funktioniert" vermeiden wollte. Ein Netzwerk-Beitritt ist diese **eine**
  Form: `example-run-*` läuft identisch, ob der Aufrufer selbst in einem
  Container sitzt oder auf dem Host — der Docker-Daemon löst den Servicenamen
  über das benannte Netzwerk auf, unabhängig davon.
- **Das löst das in `slice-beispiele-compose-bootstrap` §6 benannte
  Kollisions-Risiko** (zwei Compose-Dateien im selben Repo): die Demo-Umgebung
  bekommt ihren **eigenen** Netzwerknamen (`cdc-examples`), getrennt von
  `cdc-feed-test` der Wurzel-`compose.yaml` — beide können gleichzeitig
  laufen, ohne Namenskollision.

### 5 — Was diese ADR nicht ändert

- **`spec/architecture.md` bleibt unberührt** — `examples/**` bleibt
  Gate-Scope-Gruppe, kein `ARC-*` kommt hinzu.
- **`.a-check.yml` bleibt unberührt** — kein neuer Import-Pfad, kein neuer
  Konsument des öffentlichen Vertrags; der neue Go-Bau liest `gen/**`
  genauso, wie es das Wurzelmodul immer schon tat.
- **Der ausgelieferte Bau bleibt unberührt** — Wurzel-`Dockerfile`,
  Wurzel-`.dockerignore` (byte-gleich, real gemessen), `make image`, sein
  Digest-Beleg (`ADR-0044`) und der Coverage-Messgegenstand
  ([`ADR-0071`](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)).
- **`make examples-csharp`/`make examples-kotlin` bleiben unverändert** —
  diese ADR fügt einen Start-Aufruf **gegen** ihre Erzeugnisse hinzu, ändert
  aber nichts an Bau oder Test.
- **Kein neues Gate.** Alle drei Start-Ziele sind Werkzeuge wie ihre
  Bau-Geschwister; `make gates` bleibt netzlos.
- **Keine Zusage des Lastenhefts wird berührt** — die Beispiele zeigen
  [`LH-FA-SST-006`](../../../spec/lastenheft.md)/`LH-FA-SST-007`/`LH-FA-SST-008`
  weiterhin unverändert.

## Verglichene Alternativen

Regeln dieser Sektion: **mindestens drei Optionen mit Pro/Contra** — „nichts
tun" ist eine davon. Eine ADR ohne Alternativen ist ein Postulat, kein
Entscheidungsprotokoll, und im Review nicht verteidigbar (Baseline-Regelwerk
`modul-04-adrs.md` §Ziel-Form: ADR (MADR)).

### A — Go-Bauform

| Option | Pro | Contra |
|---|---|---|
| A1 — nichts tun: `go run ./examples/<name>` bleibt die Startform | kein Eingriff | widerspricht wörtlich der Auftraggeber-Anweisung; widerspricht der bereits geltenden `SPEC-023`-Zeile „Laufzeit und Startform" (sprachneutral „Container-Aufruf") |
| A2 — Wurzel-`Dockerfile` bekommt vier neue Stufen für die Go-Beispiele, Wurzel-`.dockerignore` bekommt `!examples/` | ein Dockerfile im Repo statt zwei | vergrößert **jede** `COPY . .`-Stufe der Wurzel-Datei (`coverage`, `build`) um `examples/**` — deren Digest-Verhalten reagiert dann auf Beispiel-Änderungen; verletzt „der ausgelieferte Bau bleibt unberührt" ([`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md) Festlegung 3, dieselbe Klasse, die dort für C#/Kotlin verworfen wurde) |
| A3 — Go-Beispiele ziehen auf das C#/Kotlin-Sprach-Wurzel-Muster um (`examples/go/<name>/`, eigenes `go.mod`) | strukturelle Symmetrie zu den zwei anderen Sprachen | verletzt [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md) Festlegung 5 (zweites `go.mod` nimmt die Beispiele aus `go test ./...` der Wurzel heraus, der bestehende Anti-Verrottungs-Träger bräuchte einen Ersatz); Arbeit am Bestand ohne Adresse — Go braucht weder Projektverzeichnis noch SDK/Runtime-Split |
| **A4 — eigenes `examples/Dockerfile`, Wurzel-Kontext, Isolation über dateispezifisches `examples/Dockerfile.dockerignore` (gewählt)** | Wurzel-`Dockerfile`/`.dockerignore` bleiben byte-gleich (real gemessen); kein zweites `go.mod`, `go test ./...` bleibt unverändert Anti-Verrottungs-Träger; Digest-Wiederverwendung statt neuer Pins; eine Datei, vier Ziel-Stufen — leichter als der C#/Kotlin-Weg, weil Go keinen SDK/Runtime-Split braucht | eine zusätzliche Ignore-Datei-Klasse im Repo (bislang nur für C#/Kotlin als *dokumentiert, nicht gemessen* diskutiert, [`ADR-0090`](0090-beispiel-clients-volle-matrix.md) §B) — hier real gemessen und für Go erstmals genutzt |

### B — Start-Make-Target-Form über drei Sprachen

| Option | Pro | Contra |
|---|---|---|
| B1 — ein Ziel je Programm (zwölf `.PHONY`-Stanzas: `example-run-go-http`, …, `example-run-kotlin-nats`) | keine Argument-Validierung nötig; jedes Ziel einzeln in `make help` sichtbar | zwölf nahezu identische Stanzas über drei `mk`-Fragmente — Wiederholung ohne strukturellen Nutzen; jede neue Oberfläche bräuchte drei neue Ziele statt einer Zeile |
| B2 — ein einziges, sprachübergreifendes Ziel `example-run LANG=… SURFACE=…` | ein Ziel, ein Einstiegspunkt | müsste intern nach Sprache verzweigen (Go: Binary-Name; C#/Kotlin: Image-Tag-Schema) — dieselbe Verzweigung wie B3, nur zentralisiert statt je Sprache; verliert die bestehende, gelebte Aufteilung „ein `mk`-Fragment-Abschnitt je Sprache" (`examples-csharp`/`examples-kotlin`) |
| **B3 — ein Ziel je Sprache, Pflicht-Argument `SURFACE=` (gewählt)** | trägt dieselbe Granularität wie die bestehenden Bau-Ziele; folgt dem bereits gelebten Repo-Idiom (Kommandozeilen-Variable überschreibt einen Default, `RANGE=` bei `make commit-traceability`); `$(error …)` bei fehlendem/unbekanntem `SURFACE` ist eine lesbare, lokale Prüfung je Sprache | drei Ziele statt eines; ein Nutzer muss `SURFACE=` kennen (in der `## `-Zeile von `make help` dokumentiert) |
| B4 — kein Start-Target, nur Doku des `docker run`-Aufrufs (Ist-Zustand für C#/Kotlin) | kein Aufwand | genau der Zustand, den die Auftraggeber-Anweisung ausdrücklich zurückweist („echten Start-Make-Target") |

### C — Umgebungsdatei-Kontrakt

| Option | Pro | Contra |
|---|---|---|
| C1 — kein gemeinsamer Kontrakt: jedes Start-Target trägt seine ENV-Werte inline im `mk`-Fragment | kein neues Artefakt | Werte stehen doppelt (Compose-Datei und `mk`-Fragment) und driften; genau die Kopier-Last, die die Auftraggeber-Anweisung („config-file") vermeiden wollte |
| C2 — `examples/.env.example` als **Vorlage**, vom Leser nach `examples/.env` zu kopieren | üblich für Repos mit echten Secrets | diese Datei trägt **keine** echten Secrets (isoliertes Demo-Netzwerk) — der Kopier-Schritt wäre reine Zeremonie und verletzt die Welle-Closure-Anforderung „ohne manuellen Zusatzschritt des Lesers" |
| **C3 — `examples/.env`, committet, sofort nutzbar, bestehende `CDC_*`-Namen, Kopfkommentar zur Netzwerk-Grenze (gewählt)** | ein Artefakt, eine Quelle der Wahrheit für Compose-Datei und Start-Targets; keine Kopier-Zeremonie; dieselbe Vokabel wie das Handbuch — ein Integrator lernt keine zweite Namensmenge | committete, wenn auch harmlose, Zugangsdaten im Baum — durch Kopfkommentar und Netzwerk-Isolation (Festlegung 4) begrenzt, nicht eliminiert |
| C4 — neue, Demo-spezifische Variablennamen (`EXAMPLE_HTTP_ADDR` etc.), um jede Verwechslung mit Produktionsnamen auszuschließen | maximale Trennschärfe im Namen | zwei Vokabulare für dieselbe Sache; ein Integrator müsste `CDC_*` aus dem Handbuch UND `EXAMPLE_*` aus der Demo-Datei lernen — mehr Fläche, kein zusätzlicher Schutz gegenüber der Netzwerk-Isolation, die ohnehin trägt |

**Fazit:** A4, B3, C3. A1 widerspricht der Auftraggeber-Anweisung und der
bereits geltenden Spec-Zeile; A2 verletzt „der ausgelieferte Bau bleibt
unberührt"; A3 verletzt den Anti-Verrottungs-Träger der Go-Beispiele. B1/B2
verlieren an Wiederholung bzw. Verzweigungs-Zentralisierung ohne echten
Gewinn; B4 widerspricht der Anweisung direkt. C1 dupliziert Werte, C2 fügt
Zeremonie ohne Sicherheitswert hinzu, C4 verdoppelt das Vokabular ohne
zusätzlichen Schutz.

## Konsequenzen

- Positiv: Die Auftraggeber-Anweisung ist **wörtlich** eingelöst — Go
  baut/startet über `make`+Dockerfile, alle drei Sprachen haben einen echten
  Start-Target, `examples/` bekommt einen Umgebungsdatei-Kontrakt.
- Positiv: **Der ausgelieferte Bau bleibt real unberührt** (byte-gleiche
  Wurzel-`.dockerignore`, gemessen) — dieselbe Disziplin, die
  [`ADR-0087`](0087-beispiel-clients-csharp-kotlin.md)/[`ADR-0090`](0090-beispiel-clients-volle-matrix.md)
  für C#/Kotlin schon durchgehalten haben, jetzt auch für Go, mit einem
  **leichteren** Mechanismus (eine Ignore-Datei statt eines eigenen
  Sprach-Wurzelverzeichnisses).
- Positiv: Der Anti-Verrottungs-Träger der Go-Beispiele
  (`go test ./...`) bleibt unverändert — kein zweites `go.mod`, keine neue
  Ausnahme.
- Positiv: Die `SPEC-023`-Zeile „Laufzeit und Startform" wird **wahr** statt
  weiter durch die Go-Ausnahme unterlaufen zu werden (gemessener Fund,
  Kontext).
- Positiv: Ein Integrator lernt **ein** Vokabular (`CDC_*`) für Handbuch,
  Demo-Umgebung und eigene Produktionsinstanz.
- Negativ: Der Go-Bau bekommt eine **zweite** Ignore-Datei-Klasse im Repo
  (`<Dockerfile>.dockerignore`) neben der bekannten `ADR-0085`-Ausnahmeklasse
  — real gemessen, aber ein weiterer Mechanismus, den ein Leser der
  `.dockerignore`-Landschaft kennen muss.
- Negativ mit Grenze: `examples/.env` committet Demo-Zugangsdaten im
  Klartext — begrenzt durch Netzwerk-Isolation (`cdc-examples`) und
  Kopfkommentar, nicht eliminiert; ein Leser, der die Datei unreflektiert
  gegen eine echte Instanz einsetzt, ist vom Kommentar abhängig (Wächter:
  Review/Doku, kein Sensor).
- Negativ mit Grenze: Drei Start-Ziele mit Pflicht-Argument sind weniger
  selbsterklärend in `make help` als zwölf benannte Ziele — der `## `-Kommentar
  trägt die `SURFACE=`-Form explizit, um das zu mildern.
- Folgepflicht (Spec-Zug, umsetzender Zug): `SPEC-023` Zeile „Laufzeit und
  Startform" bestätigt (keine Änderung nötig, sie war schon sprachneutral
  formuliert — siehe Kontext); neue Zeile für den Umgebungsdatei-Kontrakt
  und die Start-Make-Ziele; §7 Historie.
- Folgepflicht (Go-Zug, `slice-beispiele-go-dockerfile-start`):
  `examples/Dockerfile`, `examples/Dockerfile.dockerignore`,
  `harness/mk/examples.mk` (`example-run-go`), `examples/README.md`
  (Go-Tabelle, Startform-Spalte), die vier `**Beispiele:**`-Blöcke in
  `docs/user/benutzerhandbuch.md`, `harness/README.md` §Werkzeuge.
- Folgepflicht (C#/Kotlin-Zug, `slice-beispiele-csharp-kotlin-start-target`):
  `harness/mk/examples.mk` (`example-run-csharp`/`example-run-kotlin`),
  dieselben Doku-Träger für C#/Kotlin.
- Folgepflicht (Umgebungs-Zug, `slice-beispiele-compose-bootstrap`):
  `examples/.env` mit den konkreten Demo-Werten, das Netzwerk `cdc-examples`
  in der Demo-Compose-Datei, das Bootstrapping.

## Fitness Function (falls maschinell prüfbar)

| Tooling | Regel | Make-Target |
|---|---|---|
| `docker build` (Wurzel) und `make image` | Der ausgelieferte Bau bleibt unberührt: Wurzel-`Dockerfile` und Wurzel-`.dockerignore` bleiben byte-gleich; ein Bau mit `-f Dockerfile .` sieht `examples/Dockerfile.dockerignore` nicht (real gemessen, Kontext) | `make image` (kein Gate) |
| `docker build -f examples/Dockerfile .` (Werkzeug, **kein** Gate) | Baut die vier Go-Beispiele aus Wurzel-Kontext, isoliert über die dateispezifische Ignore-Datei; Basis-Digests identisch mit der Wurzel-`Dockerfile` | `example-run-go`/eigenes Bau-Ziel (Detail des umsetzenden Zuges) |
| Go-Toolchain (`go test ./...`, im gepinnten Container) | `examples/**` bleibt im Wurzelmodul und wird von `make test` übersetzt — unverändert gegenüber [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md) Festlegung 5 | `make test` |
| `make example-run-go` / `make example-run-csharp` / `make example-run-kotlin` (Werkzeug, **kein** Gate) | Startet real einen digest-gepinnten Client-Container gegen das benannte Docker-Netzwerk `cdc-examples`, mit `--env-file examples/.env`; fehlendes/unbekanntes `SURFACE=` bricht mit `$(error …)` ab, bevor `docker run` läuft; Exit-Code direkt gelesen (`AGENTS.md` §3.9) | — (Werkzeug; **nicht** in `GATE_CHECKS`) |
| Review-Prüfpflicht (nicht maschinell) | Ob `examples/.env` tatsächlich nur Demo-Werte trägt und der Kopfkommentar die Netzwerk-Grenze nennt; ob ein Start-Target versehentlich `localhost` statt des Servicenamens verwendet | — |

## Re-Evaluierungs-Trigger

Regeln dieser Sektion: **jede ADR trägt einen Trigger** — eine beobachtbare
Bedingung — oder ausdrücklich *permanent*. Ohne Trigger gilt die Entscheidung
unbefristet weiter, auch wenn ihre Voraussetzung weg ist (Baseline-Regelwerk
`modul-04-adrs.md` §Kernidee (Modul 4)).

Drei benannte Trigger, sonst permanent:

1. **Ein fünftes Go-Beispiel oder eine fünfte Oberfläche entsteht.** Dann
   bekommt `examples/Dockerfile` eine weitere `runtime-<surface>`-Stufe nach
   demselben Muster — keine neue Entscheidung nötig, solange die Oberfläche
   nur die Standardbibliothek und öffentliche Fremdmodule benutzt
   ([`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
   Festlegung 2).
2. **Die Demo-Umgebung braucht Host-Port-Zugriff** (z. B. weil ein externes
   Werkzeug außerhalb von Docker gegen sie laufen soll, nicht nur die drei
   Start-Targets). Dann ist Festlegung 4 (Netzwerk-Beitritt statt
   Port-Publizierung) neu zu bewerten — beide Formen können nebeneinander
   bestehen, aber das ist dann eine bewusste Erweiterung, keine stille.
3. **`examples/.env` wird versehentlich gegen eine echte Instanz verwendet**
   (ein real aufgetretener Vorfall, kein hypothetischer). Dann reicht der
   Kopfkommentar nicht mehr als Wächter, und Festlegung 3 ist neu zu
   bewerten (z. B. ein Datei-Namensschema, das die Demo-Zugehörigkeit
   struktureller macht als ein Kommentar).

Sonst permanent: Go-Beispiele bauen/starten über `make`+Dockerfile mit
Wurzel-Kontext und dateispezifischem Ignore; alle drei Sprachen tragen einen
Start-Make-Target mit Pflicht-`SURFACE=`-Argument; die Demo-Umgebung ist ein
benanntes, isoliertes Docker-Netzwerk mit einer gemeinsamen, committeten
Umgebungsdatei.

## Slice-Schnitt — Bestätigung mit einer Korrektur

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — der Architect-Zug bestätigt oder korrigiert den vom
Planner vorgezeichneten Schnitt der drei Folge-Slices dieser Welle.

**Bestätigt, unverändert:** der Schnitt in drei parallele, voneinander
unabhängige Implementer-Slices
(`slice-beispiele-compose-bootstrap`,
`slice-beispiele-go-dockerfile-start`,
`slice-beispiele-csharp-kotlin-start-target`) trägt — jeder bewegt eigene
Dateien (`examples/Dockerfile`+`harness/mk/examples.mk`-Go-Ziel vs.
`harness/mk/examples.mk`-C#/Kotlin-Ziele vs. `examples/.env`+Demo-Compose),
keiner braucht den anderen fertig, um zu **beginnen** (die
Soft-Abhängigkeit — ein Start-Target ohne `examples/.env` verweist auf den
hier entschiedenen Kontrakt statt auf die konkrete Datei — ist bereits in
der Welle-Datei §4 benannt und bleibt richtig).

**Korrektur — Arbeitsnamen, jetzt fest:**

- `slice-beispiele-go-dockerfile-start`: Arbeitsname `examples/Dockerfile`
  wird **verbindlich** (Festlegung 1); zusätzlich verbindlich:
  `examples/Dockerfile.dockerignore` (bislang nicht im Slice-Plan genannt —
  wird beim Öffnen des Slice als Zeile in §3 „Plan (vor Code)" nachgetragen);
  das Ziel heißt `example-run-go` (nicht `example-run-go-<protokoll>`, wie im
  Slice-Plan §5 Closure-Trigger als Arbeitsname stand — **eines** Ziel je
  Sprache mit `SURFACE=`, Festlegung 2).
- `slice-beispiele-csharp-kotlin-start-target`: die Ziele heißen
  `example-run-csharp`/`example-run-kotlin` (Festlegung 2) statt der im
  Slice-Plan §3 offen gelassenen `example-run-csharp-*`/`example-run-kotlin-*`-Arbeitsnamen
  oder eines vollständig parametrisierten Ziels.
- `slice-beispiele-compose-bootstrap`: Arbeitsname `examples/.env.example`
  wird zu **`examples/.env`** (Festlegung 3, committet statt Vorlage); das
  Docker-Netzwerk der Demo-Compose-Datei heißt **`cdc-examples`**
  (Festlegung 4, bislang im Slice-Plan nicht benannt).

Keine der drei Korrekturen ändert die **Liefer-Punkte-Zahl** oder die
**Abhängigkeits-Struktur** der drei Slices — sie ersetzen offen gelassene
Arbeitsnamen durch die hier getroffene, verbindliche Form. Die drei
Slice-Pläne werden mit dieser ADR-Kennung und den festen Namen
nachgezogen (dieser Zug, siehe Bericht).

## Geschichte

| Datum | Ereignis | Verweis |
|---|---|---|
| 2026-09-18 | Accepted — Anlass: explizite Auftraggeber-Anweisung vom 2026-09-18, die Go-Startform von `go run` auf `make`+Dockerfile umzukehren (Supersedes [`ADR-0076`](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md) Festlegung 1, Startform-Bullet), plus die Anweisung, allen drei Sprachen einen echten Start-Make-Target und `examples/` eine Demo-Umgebung mit Umgebungsdatei-Kontrakt zu geben. Unabhängiger Architect-Zug entscheidet die Go-Bauform (eigenes Dockerfile, Wurzel-Kontext, dateispezifisches Ignore — real gemessen), die Start-Target-Form (ein Ziel je Sprache, Pflicht-`SURFACE=`), den Umgebungsdatei-Kontrakt (`examples/.env`, bestehende `CDC_*`-Namen, committet) und das Netzwerkmodell der Demo-Umgebung (`cdc-examples`); bestätigt bzw. korrigiert den Planner-Schnitt der drei Folge-Slices; Modul 8 §Rollen-Regeln: „Architect schreibt" | gemessen in diesem Zug: der `examples/README.md`/`harness/mk/examples.mk`-Ist-Stand, die `SPEC-023`-Zeile „Laufzeit und Startform" (sprachneutral, im Widerspruch zum `go run`-Zustand), und die dateispezifische `<Dockerfile>.dockerignore`-Isolation (`docker build -f examples/Dockerfile .` mit/ohne `examples/Dockerfile.dockerignore`, `docker build -f Dockerfile .` zur Gegenprobe, `docker` v29.8.0) |

Nach `Accepted` wird diese Datei **nicht mehr inhaltlich überschrieben**.
Spätere Korrekturen oder Schärfungen entstehen als neue ADR mit
`Supersedes ADR-0098` (Baseline-Regelwerk `modul-04-adrs.md`
§Hard Rule für Accepted-ADRs).
