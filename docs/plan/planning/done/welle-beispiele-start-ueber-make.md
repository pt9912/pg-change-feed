# Welle beispiele-start-ueber-make: Beispiel-Clients — Start über `make`/Dockerfile statt `go run`, echter Start-Make-Target, Demo-Umgebung mit Bootstrapping

**Lifecycle:** Diese Datei entsteht bei der **Eröffnung** der Welle und liegt
flach unter `docs/plan/planning/`; bei Closure wandert sie per `git mv` nach
`done/` (neben ihre `welle-beispiele-start-ueber-make-results.md`). Der
Zustand ist die Verzeichnis-Position — kein Status-Feld.

**Zielmeilenstein:** kein Meilenstein-Bezug.

**Verantwortlich:** —. **Datum:** 2026-09-18.

---

## 1. Welle-Ziel

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht.

Der Auftraggeber hat drei Entscheidungen getroffen, die den bestehenden
Beispiel-Client-Bestand (`examples/**`,
[`ADR-0076`](../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md),
[`ADR-0087`](../adr/0087-beispiel-clients-csharp-kotlin.md),
[`ADR-0090`](../adr/0090-beispiel-clients-volle-matrix.md)) überholen:

1. **Die Go-Beispiele werden über `make` + Dockerfile gebaut und gestartet**,
   nicht über `go run ./examples/<name>`. Das widerspricht wörtlich
   [`ADR-0076`](../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
   §Entscheidung Festlegung 1 dritter Bullet („**Startform**: je Programm ein
   `go run ./examples/<name>` … Damit zitiert das Handbuch **eine**
   Befehlsform, die innerhalb und außerhalb des Compose-Netzes
   funktioniert."). Eine `Accepted`-ADR wird nach `AGENTS.md` §3.5 nicht
   inhaltlich überschrieben — die Umkehr braucht eine neue ADR mit
   `Supersedes ADR-0076` in genau dieser Klausel.
2. **Alle drei Sprachen (Go, C#, Kotlin) bekommen einen echten
   Start-Make-Target** — ein `make`-Ziel, das ein Beispiel-Programm
   tatsächlich **startet** (Container- bzw. Prozess-Lauf gegen Adresse/Token),
   nicht nur baut/testet. Heute bauen `make examples-csharp`/
   `make examples-kotlin` nur (`harness/README.md` §Werkzeuge, „kein
   Lauf-Beleg"); ein Start läuft manuell über
   `docker run --rm <ENV> <Image> <Argumente>` — real dokumentiert im
   `**Beispiele:**`-Block der vier Zugriffs-Abschnitte in
   `docs/user/benutzerhandbuch.md` (z. B. §4 „Zugriff über
   Server-Sent-Events", Zeilen „**C#:**"/„**Kotlin:**"). Für Go gibt es
   bislang gar keinen `make`-Bau; die zwölf Programme
   ([`ADR-0090`](../adr/0090-beispiel-clients-volle-matrix.md) §Entscheidung
   Festlegung 1) sind heute uneinheitlich baubar/startbar (Go: `go run`;
   C#/Kotlin: `make`-Bau + manueller `docker run`).
3. **Eine Demo-Umgebung unter `examples/` mit Bootstrapping.** Eine eigene
   Compose-Datei unter `examples/` (Arbeitsname `examples/compose.yaml` —
   der endgültige Name ist Detail des umsetzenden Zuges, siehe Slice 2)
   fährt PostgreSQL-Quelle, den Feed-Container und — für das NATS-Beispiel —
   NATS hoch, **und** treibt danach den Schema-Rollout (d-migrate,
   `tools/schema/schema.yaml`) sowie die Registrierung/Aktivierung einer
   Beispiel-Quelle und -Tabelle (`cdc.source`, `cdc.enable_table` oder
   gleichwertig), damit die Beispiel-Clients gegen echte Daten laufen, ohne
   dass der Leser das von Hand einrichtet. Diese Umgebung ist **bewusst
   getrennt** von der Wurzel-`compose.yaml` — deren Datei-Kopfkommentar
   nennt sie ausdrücklich als „Compose-Umgebung der Lauf- und CI-Vertrag-Seite
   … Die Umgebung läuft über `tools/harness/run-integration-tests.sh`" (Test-
   /CI-Harness, kein Leser-Erzeugnis). Die neue Datei ist eine
   **Leser-Quickstart-Umgebung** für `examples/` mit anderem Zweck und
   anderem Betreiber (der Integrator, nicht der Testlauf).

   **Das „config-file" — von Slice 1 bestätigt, mit einer Korrektur**
   ([`ADR-0098`](../adr/0098-beispiel-clients-start-ueber-make-dockerfile.md)
   Festlegung 3): die Planner-Lesart trifft zu — eine gemeinsame
   Umgebungsdatei, die **sowohl** die neue Compose-Datei (`env_file:`) **als
   auch** jeder Start-Make-Target-Aufruf (`--env-file`) liest, statt Werte
   von Hand zwischen Compose-Datei und jedem Beispiel-Aufruf zu kopieren. Die
   Korrektur: **fester Dateiname `examples/.env`** (nicht
   `examples/.env.example`) — die Datei wird **committet und sofort
   nutzbar**, keine Kopiervorlage, weil sie keine echten Zugangsdaten trägt
   (isoliertes Docker-Netzwerk `cdc-examples`,
   [`ADR-0098`](../adr/0098-beispiel-clients-start-ueber-make-dockerfile.md)
   Festlegung 4). Die
   Variablennamen sind die bereits im Handbuch dokumentierten `CDC_*`-Namen,
   bewusst identisch mit der Produktions-Vokabel — die Trennung liegt in der
   Netzwerk-Isolation, nicht im Namen. Die verworfene Lesart
   (`tools/schema/schema.yaml` als „config-file") ist damit ausgeschlossen:
   jene Datei bleibt der Schema-Rollout-Eingang, unberührt von dieser
   Entscheidung.

**Das *Mehr* gegenüber den einzelnen Slice-DoDs** ist die **gemeinsame
Design-Entscheidung**, die alle drei Hälften trägt und die noch nicht
getroffen ist: Ob Go auf das C#/Kotlin-Muster (eigenes Dockerfile je
Sprach-Wurzel) umgezogen wird oder ein leichteres, sprach-spezifisches
Bauform bekommt (Go-Binaries sind trivial statisch/cross-kompilierbar,
anders als C#/Kotlin mit SDK+Runtime-Split); welche Form der neue
Start-Target über **drei** Sprachen hinweg trägt (ein Ziel je Programm, ein
parametrisiertes Ziel, oder ein Ziel je Sprache mit Pflicht-Argumenten); und
das **gemeinsame Umgebungsdatei-Format/die Variablen-Namen**, gegen die
sowohl die Demo-Compose-Datei als auch beide Start-Target-Implementierungen
schreiben. Keine der drei Implementierungs-Hälften ist ohne diese
Entscheidung unabhängig lieferbar — sie ist der gemeinsame Träger, den kein
Einzel-Slice-DoD abdeckt, und der Grund, warum eine `Accepted`-Immutable-ADR
([`ADR-0076`](../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md))
hier nicht durch einen Implementer-Slice, sondern nur durch eine neue
Architect-Entscheidung bewegt werden kann (`AGENTS.md` §3.5).

**Register-Sichtung (Eröffnungs-Schritt 2):** `docs/plan/planning/observations/`
durchsucht auf `beispiel`/`example`/`start`/`make-target`/`compose` —
**keine Treffer**. Diese Welle eröffnet keinen bestehenden
Beobachtungs-Eintrag.

## 2. Trigger (Welle startet)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Regeln — ein Trigger ist **beobachtbar** dann, wenn ein *anderer*
Mensch ohne Rückfrage sagen kann, ob er eingetreten ist; ein Datum darf erwähnt
werden, aber nie Trigger sein. Und der **Start**-Trigger ist **kein Ergebnis
dieser Welle**: Steht er in der Slice-Liste unten, ist er falsch platziert.

- Explizite Auftraggeber-Entscheidung vom 2026-09-18 (drei Teile, siehe §1) —
  eingetreten, dies ist der Anlass der Eröffnung.

## 3. Closure-Trigger (Welle schließt)

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wann Arbeit eine Welle braucht — der Trigger muss das *Mehr* gegenüber den
einzelnen Slice-DoDs benennen; kann er das nicht, liegt keine Welle vor.

- Die Supersedes-ADR (§4, Slice 1) liegt `Accepted` vor und entscheidet die
  Go-Bauform, die Start-Target-Form für alle drei Sprachen **und** das
  Umgebungsdatei-/Variablen-Kontrakt der Demo-Umgebung.
- `examples/http-client`, `examples/sse-client`, `examples/grpc-client`,
  `examples/nats-client` (Go) sind über ein `make`-Ziel baubar **und** über
  ein `make`-Ziel startbar.
- C# und Kotlin bekommen denselben echten Start-Make-Target (nicht nur den
  bestehenden Bau/Test) für ihre je vier Programme.
- Die Demo-Compose-Datei unter `examples/` bringt PostgreSQL-Quelle,
  Feed-Container und NATS real hoch, rollt das Schema aus und
  registriert/aktiviert eine Beispiel-Quelle und -Tabelle — ohne manuellen
  Zusatzschritt des Lesers.
- Alle vier Slices (§4) liegen in `done/`.
- `make gates` grün.
- `examples/README.md`, die vier Zugriffs-Abschnitte in
  `docs/user/benutzerhandbuch.md` und `harness/README.md` §Werkzeuge nennen
  die neue Start-Form und die Demo-Umgebung konsistent über alle drei
  Sprachen.
- Closure-Notiz in `welle-beispiele-start-ueber-make-results.md`.

## 4. Slices in dieser Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Lifecycle als State Machine — der Zustand eines Slice ist sein
Lifecycle-Verzeichnis und wird hier **nicht** gespiegelt.

| Slice | Titel | Bezug |
|---|---|---|
| slice-beispiele-start-architect-entscheidung | Architect-Entscheidung: Supersedes-ADR zu `ADR-0076` Festlegung 1 (Go-Startform) + Ausgestaltung des Start-Make-Targets + Umgebungsdatei-Kontrakt der Demo-Umgebung | [`ADR-0076`](../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md), [`ADR-0087`](../adr/0087-beispiel-clients-csharp-kotlin.md), [`ADR-0090`](../adr/0090-beispiel-clients-volle-matrix.md) |
| slice-beispiele-compose-bootstrap | Demo-Umgebung: `examples/compose.yaml` (PostgreSQL/Feed/NATS) + gemeinsame Umgebungsdatei + Bootstrapping (Schema-Rollout, Beispiel-Quelle/-Tabelle) | [`LH-QA-OPS-001`](../../../spec/lastenheft.md) |
| slice-beispiele-go-dockerfile-start | Go-Beispiele: Dockerfile(s) + `make`-Bau-Target + `make`-Start-Target (vier Programme) | [`LH-FA-SST-006`](../../../spec/lastenheft.md), [`LH-FA-SST-007`](../../../spec/lastenheft.md), [`LH-FA-SST-008`](../../../spec/lastenheft.md) |
| slice-beispiele-csharp-kotlin-start-target | C#/Kotlin: echter Start-Make-Target ergänzt (acht Programme, zwei Sprachen) | [`LH-FA-SST-006`](../../../spec/lastenheft.md), [`LH-FA-SST-007`](../../../spec/lastenheft.md), [`LH-FA-SST-008`](../../../spec/lastenheft.md) |

**Reihenfolge — explizit:** `slice-beispiele-start-architect-entscheidung`
zuerst und blockierend für alle drei anderen — ohne die ADR gibt es weder
eine entschiedene Go-Bauform noch eine entschiedene Start-Target-Form noch
einen entschiedenen Umgebungsdatei-Kontrakt, und ein Implementer-Slice, der
vorgriffe, widerspräche stillschweigend einer künftigen ADR (`AGENTS.md`
§3.5/§3.6). Danach: `slice-beispiele-compose-bootstrap`,
`slice-beispiele-go-dockerfile-start` und
`slice-beispiele-csharp-kotlin-start-target` sind **voneinander unabhängig**
und parallel lieferbar (verschiedene Verzeichnisse, keine gemeinsame Datei);
alle drei referenzieren lediglich denselben, von Slice 1 entschiedenen
Umgebungsdatei-Kontrakt (Dateiname, Variablen-Namen). Ein Implementer-Slice,
der vor `slice-beispiele-compose-bootstrap` fertig wird, verweist auf den
**Kontrakt** aus der ADR, nicht auf die konkrete Datei — die Datei selbst
muss nicht zuerst existieren.

## 5. Abhängigkeiten

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Roadmap-Struktur: fünf Abschnitte.

- Blockiert: keine nachfolgende Welle ist derzeit geplant (*Nächste Wellen*
  in der Roadmap ist leer).
- Wird blockiert von: keine bestehende `Accepted`-ADR — die blockierende
  Instanz ist die **noch zu schreibende** Supersedes-ADR dieser Welle selbst
  (Slice 1), kein externes Vorwellen-Ergebnis.

## 6. Out-of-Scope für diese Welle

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Eröffnung Schritt 1 — Out-of-Scope gehört zur
Zielsetzung: Was nicht ausdrücklich ausgeschlossen ist, dehnt die Welle, bis
der Closure-Trigger unerreichbar wird.

- **Keine neuen Zugriffs-Oberflächen oder Sprachen.** Die Matrix bleibt, was
  [`ADR-0090`](../adr/0090-beispiel-clients-volle-matrix.md) festgelegt hat
  (vier Oberflächen × drei Sprachen, zwölf Programme) — diese Welle ändert
  **wie** gebaut/gestartet/demonstriert wird, nicht **was** existiert.
- **Kein Gate.** Die Bau-/Start-Ziele bleiben Werkzeuge
  ([`ADR-0087`](../adr/0087-beispiel-clients-csharp-kotlin.md) Festlegung 4,
  [`ADR-0090`](../adr/0090-beispiel-clients-volle-matrix.md) Festlegung 5) —
  ob das für Go unverändert gilt, entscheidet Slice 1, aber „neues Gate" ist
  in keinem Zuschnitt dieser Welle vorgesehen.
- **Kein realer Compose-/E2E-Lauf der Beispiele als Beleg für `make gates`.**
  Der Draht-Beleg bleibt bei `make test-integration` über die
  Wegwerf-Clients ([`ADR-0076`](../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
  Festlegung 5) — die neue Demo-Umgebung ist ein **Leser**-Erzeugnis, kein
  zusätzlicher Gate-Beleg, und ersetzt `tools/harness/run-integration-tests.sh`
  nicht.
- **Der `.proto`-Weg und `gen/**`** bleiben, was
  [`ADR-0090`](../adr/0090-beispiel-clients-volle-matrix.md) Festlegung 2/3
  entschieden hat — unverändert, sofern Slice 1 nichts anderes entscheidet
  (nicht erwartet: kein neuer Konsument der `.proto` entsteht durch einen
  Start-Mechanismus oder die Demo-Umgebung).
- **Der ausgelieferte Bau** (Wurzel-`Dockerfile`, `.dockerignore`,
  `make image`) bleibt unberührt — die Beispiele bauen in eigenen
  Kontexten, wie bisher; die Demo-Compose-Datei referenziert das per
  `make image` gebaute Image, ohne einen eigenen `build:`-Block
  (dieselbe Disziplin wie die Wurzel-`compose.yaml`, `ADR-0044`).
- **Kein eigener Nachweis in `docs/user/e2e-abdeckung.md`** für die
  Demo-Umgebung — sie ist Doku/Quickstart, kein E2E-Testtier.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-traceability.md`
§Herkunfts-Anker für Steering-Loop-Regeln.

Ergebnis: `welle-beispiele-start-ueber-make-results.md`, Geschwister im
Ruheort `done/`.
Zähler: `../observations/README.md`, eine Ebene über dem Ruheort.
