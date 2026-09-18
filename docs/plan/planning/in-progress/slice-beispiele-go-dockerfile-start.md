# Slice beispiele-go-dockerfile-start: Go-Beispiele — Dockerfile(s) + `make`-Bau-Target + `make`-Start-Target

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-beispiele-start-ueber-make.

**Bezug:** [`LH-FA-SST-006`](../../../../spec/lastenheft.md),
[`LH-FA-SST-007`](../../../../spec/lastenheft.md),
[`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Scope: die Oberflächen,
die die vier Go-Beispiele zeigen), die neue Supersedes-ADR aus
`slice-beispiele-start-architect-entscheidung`:
[`ADR-0098`](../../adr/0098-beispiel-clients-start-ueber-make-dockerfile.md)
(Go-Bauform: `examples/Dockerfile` + `examples/Dockerfile.dockerignore`,
Start-Target-Form: `example-run-go` mit Pflicht-`SURFACE=`,
Umgebungsdatei-Kontrakt: `examples/.env`),
[`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
(Rest bestätigt: Ort, Import-Grenze, Kompilierpfad; Startform-Bullet von
`ADR-0098` superseded),
[`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) (der
Go-gRPC-Client und seine Kante `examples → contract`, unverändert).

**Berührte Spec-Stellen:** [`SPEC-023`](../../../../spec/pflichtenheft.md)
(Beispiel-Clients — Startform-Zeile, falls die ADR sie dort verankert).

**Verantwortlich:** pt9912.

**Autor:** Planner-Rolle. **Datum:** 2026-09-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die vier Go-Beispiele (`examples/http-client`, `examples/sse-client`,
`examples/grpc-client`, `examples/nats-client`) werden über die in
`slice-beispiele-start-architect-entscheidung` entschiedene Bauform per
`make` gebaut **und** über ein `make`-Ziel real gestartet (Container- oder
Prozess-Lauf gegen Adresse/Token aus dem gemeinsamen Umgebungsdatei-Kontrakt);
`go run ./examples/<name>` bleibt als Fallback funktionsfähig (kein Bruch der
Standardbibliotheks-Importierbarkeit), ist aber nicht mehr die im Handbuch
zitierte Startform.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Entscheidung selbst** (welche Bauform, welche Start-Target-Form) —
  `slice-beispiele-start-architect-entscheidung` übernimmt; dieser Slice
  setzt sie nur um.
- **C#/Kotlin-Start-Targets** — `slice-beispiele-csharp-kotlin-start-target`
  übernimmt; andere Sprach-Wurzel, andere Datei, kein Überschneidungspunkt.
- **Die Demo-Compose-Datei und das Bootstrapping** —
  `slice-beispiele-compose-bootstrap` übernimmt; dieser Slice liest nur den
  von dort/der ADR entschiedenen Umgebungsdatei-Kontrakt, legt ihn nicht an.
- **Der Go-gRPC-Client selbst wird nicht neu geschrieben** — er existiert
  bereits ([`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md));
  dieser Slice ändert nur, wie er gebaut/gestartet wird, nicht seinen Inhalt.
- **`.a-check.yml`** — bleibt unberührt, sofern die ADR keinen neuen Import
  verlangt (nicht erwartet: ein Bau-/Start-Mechanismus ist keine
  Import-Frage).

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **LP1** — Dockerfile(s) für die vier Go-Beispiele in der von der ADR
      entschiedenen Form (ein gemeinsames oder je Programm), digest-gepinnte
      Basis, `make`-Bau-Ziel baut alle vier.
- [x] **LP2** — `make`-Start-Ziel(e) starten je Programm real gegen die im
      Umgebungsdatei-Kontrakt vereinbarten Variablen (Flag-Override bleibt
      erhalten); Exit-Code-Disziplin nach `AGENTS.md` §3.9.
- [x] **LP3** — Träger nachgezogen: `examples/README.md` (Go-Tabelle,
      Startform-Spalte), die vier Zugriffs-Abschnitte in
      `docs/user/benutzerhandbuch.md` (Go-Zeile der `**Beispiele:**`-Blöcke),
      `harness/README.md` §Werkzeuge (neue Ziele benannt, kein Gate).
- [x] [`LH-FA-SST-006`](../../../../spec/lastenheft.md)/[`LH-FA-SST-007`](../../../../spec/lastenheft.md)/[`LH-FA-SST-008`](../../../../spec/lastenheft.md)
      weiterhin gezeigt (kein Verhaltens-, nur
      Bau-/Start-Mechanismus-Wechsel) — `make test` übersetzt die vier
      Beispiele unverändert.
- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update für die Go-Startform (öffentlicher Vertrag: das Handbuch
      zitiert eine Befehlsform) — siehe LP3.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`, oder „keine Beobachtung
      angefallen" in §7.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind von der
      nächsten Welle-Closure getragen.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `examples/Dockerfile` (fest, `ADR-0098` Festlegung 1) | neu | Bauform für die vier Go-Beispiele (Wurzel-Kontext, vier `runtime-<surface>`-Stufen) |
| `examples/Dockerfile.dockerignore` (fest, `ADR-0098` Festlegung 1) | neu | isoliert den Bau-Kontext (`go.mod`/`go.sum`/`gen/`/`examples/`), Wurzel-`.dockerignore` bleibt byte-gleich |
| `harness/mk/examples.mk` | update | neues Ziel `example-run-go` mit Pflicht-Argument `SURFACE=http|sse|grpc|nats` (`ADR-0098` Festlegung 2) |
| `Makefile` | **nicht realisiert** — kein Änderungsbedarf | `harness/mk/*.mk` wird per `include` bereits vollständig eingebunden (Zeile 17 der Wurzel-`Makefile`); `example-run-go` steht damit ohne weitere Verdrahtung in `make help`/`make gates`-Kontext zur Verfügung, real geprüft (`make help` listet das Ziel, `make -n example-run-go SURFACE=<x>` löst korrekt auf) |
| `examples/README.md` | update | Go-Tabelle: Startform-Spalte auf die neue Form |
| `docs/user/benutzerhandbuch.md` | update | vier `**Beispiele:**`-Blöcke, Go-Zeile |
| `harness/README.md` | update | §Werkzeuge, neue Zeile(n) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-beispiele-start-architect-entscheidung`
liegt `done/`, mit `Accepted`-ADR, die die Go-Bauform, die Start-Target-Form
und den Umgebungsdatei-Kontrakt entscheidet.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die entschiedene
  Bauform verlangt mehr als drei Liefer-Punkte (z. B. vier separate
  Dockerfiles statt eines gemeinsamen).
- `in-progress` → `open` (blockiert — Carveout?): Der Umgebungsdatei-Kontrakt
  aus `slice-beispiele-compose-bootstrap` liegt noch nicht vor und die
  Start-Target-Implementierung kann ohne ihn nicht sinnvoll gegen einen
  echten Default laufen (Soft-Abhängigkeit, siehe Welle-Datei §4).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

- `make example-run-go SURFACE=<surface>` (`ADR-0098` Festlegung 2) startet
  real ein Go-Beispiel gegen eine laufende Umgebung (Demo-Compose oder
  Compose-Testumgebung).
- `make gates` grün, Review-Report vorliegt, Closure-Notiz mit Lerneintrag
  geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Ein neues Go-Dockerfile könnte versehentlich denselben Namen/Kontext wie
  die Wurzel-`Dockerfile` kollidieren lassen (Bau-Kontext-Verwechslung) —
  **Ausgang: entfallen.** `examples/Dockerfile` liegt unter `examples/`, nicht
  an der Wurzel; `examples/Dockerfile.dockerignore` wirkt laut Docker-Mechanismus
  ausschließlich für Bauten mit `-f examples/Dockerfile`. Real gemessen in
  diesem Zug: `git diff --stat -- Dockerfile .dockerignore` bleibt leer nach
  der Implementierung, und ein `docker build -f Dockerfile .` (der
  ausgelieferte Bau) lädt weiterhin dieselbe Wurzel-`.dockerignore`
  (`load .dockerignore: 1.97kB`, unverändert) — keine Kollision.
- Der Start-Mechanismus könnte den bestehenden `go run`-Kompilierpfad
  (`make test`) versehentlich aus `go test ./...` herausnehmen (z. B. durch
  ein zweites `go.mod`) — **Ausgang: entfallen.** Kein zweites `go.mod`
  entstanden; `examples/**` bleibt Teil des Wurzelmoduls. Real gemessen: ein
  frischer `make test`-Lauf (Race-Detector, gepinnter Toolchain-Container)
  endet Exit 0 und listet alle vier Beispiel-Pakete explizit als `ok`
  (`examples/http-client`, `examples/sse-client`, `examples/grpc-client`,
  `examples/nats-client`).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register · `grundlagen-traceability.md` §Herkunfts-Anker
für Steering-Loop-Regeln.

**Entwurf (Implementer-Rolle, vor Review/Verifikation — Planner übernimmt
oder korrigiert bei tatsächlicher Closure):**

**Gegenstand:** vollständig geliefert — `make gates` grün (dieser Zug), der
unabhängige Reviewer-Pass gegen
[`ADR-0098`](../../adr/0098-beispiel-clients-start-ueber-make-dockerfile.md)
steht noch aus; der `git mv` nach `done/` folgt erst danach.

**Ergebnis:** Die vier Go-Beispiele bauen/starten jetzt über
`examples/Dockerfile` (Wurzel-Bau-Kontext, isoliert über
`examples/Dockerfile.dockerignore`, Basis-Digests aus der Wurzel-`Dockerfile`
wiederverwendet — real gebaut, alle vier Ziel-Stufen, in diesem Zug) und
`make example-run-go SURFACE=http|sse|grpc|nats` (Pflicht-Argument, `$(error
…)` bei fehlendem/unbekanntem Wert, real per `make -n`/Fehlerprobe geprüft).
Die Wurzel-`Dockerfile`/`.dockerignore` bleiben byte-gleich (`git diff`
leer nach der Implementierung). `examples/README.md`, die vier
`**Beispiele:**`-Blöcke in `docs/user/benutzerhandbuch.md` (samt
Versionshistorie) und `harness/README.md` §Werkzeuge zitieren die neue
Startform; `go run ./examples/<name>` bleibt technisch funktionsfähig
(`make test` real gelaufen, alle vier Beispiel-Pakete `ok`), ist aber nicht
mehr die zitierte Form. `example-run-csharp`/`example-run-kotlin` und
`examples/.env` bleiben bewusst außerhalb dieses Slices (Geschwister-Slices).

**Steering-Loop-Lerneintrag:** Die dateispezifische
`<Dockerfile>.dockerignore`-Isolation, die
[`ADR-0098`](../../adr/0098-beispiel-clients-start-ueber-make-dockerfile.md)
für den Go-Bau real gemessen hat, trägt real auch in der Umsetzung —
`docker build -f Dockerfile .` (der ausgelieferte Bau) lädt weiterhin
dieselbe Wurzel-`.dockerignore` unverändert, ein `docker build -f
examples/Dockerfile .` sieht zusätzlich `examples/Dockerfile.dockerignore`.
Kein neuer Mechanismus nötig, keine Überraschung gegenüber der ADR-Probe.

**Beobachtungs-Register:** keine Beobachtung angefallen — die Umsetzung
folgte den vier Festlegungen der ADR ohne Abweichung, die eine neue oder
zitierbare `BEO-PGC/*`-Musterlücke begründen würde.

**Risiken (§6):** beide mit Ausgang „entfallen" versehen, real gemessen
(byte-gleiche Wurzel-Dateien, `make test` grün mit allen vier
Beispiel-Paketen) — siehe §6.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührte Sub-Areas sind
`examples/**` (Go-Anteil), `harness/mk/**`/`Makefile`, `docs/user/**` — alle
Default-Sub-Area des Repos, keine feinere Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`beispiel`/`example`/`start`/`make-target`) — **keine Treffer** (siehe
Welle-Datei §1).

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.

### Sub-Area: `*` (Default)

- **Modus:** GF — `harness/conventions.md` §Modus-Deklaration.
