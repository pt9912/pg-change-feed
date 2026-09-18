# Slice beispiele-csharp-kotlin-start-target: C#/Kotlin — echter Start-Make-Target ergänzt

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
die die acht C#/Kotlin-Beispiele zeigen), die neue Supersedes-ADR aus
`slice-beispiele-start-architect-entscheidung`:
[`ADR-0098`](../../adr/0098-beispiel-clients-start-ueber-make-dockerfile.md)
(Start-Target-Form: `example-run-csharp`/`example-run-kotlin` mit
Pflicht-`SURFACE=`, Umgebungsdatei-Kontrakt: `examples/.env`, festes
Docker-Netzwerk `cdc-examples`),
[`ADR-0087`](../../adr/0087-beispiel-clients-csharp-kotlin.md),
[`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) (Bestand:
Sprach-Wurzel, eigenes Dockerfile je Sprache, `make examples-csharp`/
`make examples-kotlin` als Bau-/Test-Ziel — bleibt unverändert; dieser Slice
ergänzt nur den Start).

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

**Ziel:** Die acht bereits gebauten C#/Kotlin-Beispiele (je vier Programme:
`http-client`, `sse-client`, `grpc-client`, `nats-client`) bekommen den in
`slice-beispiele-start-architect-entscheidung` entschiedenen echten
Start-Make-Target — ein `make`-Aufruf, der ein Programm real gegen Adresse/
Token aus dem gemeinsamen Umgebungsdatei-Kontrakt startet, statt dass der
Leser den `docker run --rm <ENV> <Image> <Argumente>`-Aufruf von Hand aus dem
Handbuch abtippt. `make examples-csharp`/`make examples-kotlin` (Bau + Test)
bleiben unverändert bestehen.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Entscheidung selbst** (welche Start-Target-Form) —
  `slice-beispiele-start-architect-entscheidung` übernimmt.
- **Go-Start-Target** — `slice-beispiele-go-dockerfile-start` übernimmt;
  andere Bauform (Binary statt Sprach-Wurzel-Image), anderer Dockerfile-Baum.
- **Die Demo-Compose-Datei und das Bootstrapping** —
  `slice-beispiele-compose-bootstrap` übernimmt; dieser Slice liest nur den
  Umgebungsdatei-Kontrakt, legt ihn nicht an.
- **Änderungen an den Dockerfiles/Werkzeugketten selbst** (Basis-Images,
  Paket-Pins) — bleibt Bestand aus
  [`ADR-0087`](../../adr/0087-beispiel-clients-csharp-kotlin.md)/
  [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md); dieser Slice
  fügt nur einen Start-Aufruf gegen die bereits gebauten Images hinzu.
- **Der nicht-blockierende CI-Workflow** — bleibt, was er ist
  ([`ADR-0087`](../../adr/0087-beispiel-clients-csharp-kotlin.md)
  Festlegung 4); ein Start-Target ist kein Bau-/Test-Schritt und braucht
  keinen Workflow-Eintrag, sofern die ADR nichts anderes verlangt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **LP1** — Start-Make-Ziel(e) für die vier C#-Programme in der von der
      ADR entschiedenen Form, gegen die im Umgebungsdatei-Kontrakt
      vereinbarten Variablen.
- [x] **LP2** — dieselbe Form für die vier Kotlin-Programme.
- [x] **LP3** — Träger nachgezogen: `examples/README.md` (C#/Kotlin-Tabellen,
      Startform-Spalte auf `make`-Aufruf statt rohem `docker run`), die vier
      Zugriffs-Abschnitte in `docs/user/benutzerhandbuch.md` (C#-/Kotlin-Zeilen
      der `**Beispiele:**`-Blöcke), `harness/README.md` §Werkzeuge.
- [x] [`LH-FA-SST-006`](../../../../spec/lastenheft.md)/[`LH-FA-SST-007`](../../../../spec/lastenheft.md)/[`LH-FA-SST-008`](../../../../spec/lastenheft.md)
      weiterhin gezeigt (kein Verhaltens-, nur
      Start-Mechanismus-Wechsel) — `make examples-csharp`/
      `make examples-kotlin` bleiben unverändert grün.
- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update für die C#/Kotlin-Startform (öffentlicher Vertrag) — siehe
      LP3.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag. *(Entwurf in §7 vorbereitet,
      wird nach dem Review-Pass finalisiert.)*
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
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
| `harness/mk/examples.mk` | update | neue Ziele `example-run-csharp`/`example-run-kotlin`, je mit Pflicht-Argument `SURFACE=http|sse|grpc|nats` (`ADR-0098` Festlegung 2) |
| `examples/README.md` | update | C#-/Kotlin-Tabellen, Startform-Spalte |
| `docs/user/benutzerhandbuch.md` | update | vier `**Beispiele:**`-Blöcke, C#-/Kotlin-Zeilen |
| `harness/README.md` | update | §Werkzeuge, neue Zeile(n) |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-beispiele-start-architect-entscheidung`
liegt `done/`, mit `Accepted`-ADR, die die Start-Target-Form und den
Umgebungsdatei-Kontrakt entscheidet.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die entschiedene
  Form verlangt getrennte Slices je Sprache (analog zum Muster in
  `ADR-0090` §Slice-Schnitt-Empfehlung).
- `in-progress` → `open` (blockiert — Carveout?): Der Umgebungsdatei-Kontrakt
  aus `slice-beispiele-compose-bootstrap` liegt noch nicht vor (Soft-
  Abhängigkeit, siehe Welle-Datei §4).

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

- Ein `make`-Aufruf startet real je ein C#- und ein Kotlin-Beispiel gegen
  eine laufende Umgebung (Demo-Compose oder Compose-Testumgebung).
- `make gates` grün, Review-Report vorliegt, Closure-Notiz mit Lerneintrag
  geschrieben.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Ein Start-Target, das `docker run` fest verdrahtet, könnte mit einem
  künftigen Netzwerk-Alias-Kontrakt der Demo-Compose-Datei kollidieren
  (Adresse `localhost` vs. Compose-Service-Name) — **Ausgang:** entfallen.
  `ADR-0098` Festlegung 3/4 legt `examples/.env` bereits auf
  Compose-Servicenamen fest (`CDC_HTTP_ADDR=pg-change-feed:8090`, nicht
  `localhost`); real gegen die laufende Demo-Umgebung geprüft —
  `make example-run-csharp SURFACE=http ARGS="--source demo-source
  --publication pub_demo"` und dieselbe Kotlin-Form liefern real die
  registrierte Tabelle über `cdc-examples`, kein Adress-Konflikt.
- Acht Programme über zwei Sprachen könnten die Drei-Liefer-Punkte-Grenze
  real sprengen, falls die Start-Target-Form pro Programm statt pro Sprache
  verdrahtet werden muss — **Ausgang:** entfallen. `ADR-0098` Festlegung 2
  entscheidet ein Ziel je Sprache mit Pflicht-`SURFACE=`; die Umsetzung
  bleibt bei zwei neuen `.PHONY`-Zielen (LP1/LP2), keine Zerlegung nötig.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register · `grundlagen-traceability.md` §Herkunfts-Anker
für Steering-Loop-Regeln.

**Entwurf — wird nach dem unabhängigen Reviewer-Pass finalisiert (Modul 8,
kein Self-Review).**

**Gegenstand:** Implementierung abgeschlossen, wartet auf den unabhängigen
Reviewer-Pass gegen
[`ADR-0098`](../../adr/0098-beispiel-clients-start-ueber-make-dockerfile.md)
(Festlegung 2); der `git mv` nach `done/` folgt erst danach.

**Ergebnis:** `harness/mk/examples.mk` trägt zwei neue `.PHONY`-Ziele,
`example-run-csharp` und `example-run-kotlin` — Pflicht-Argument `SURFACE=
http|sse|grpc|nats`, `$(error …)` bei fehlendem/unbekanntem Wert (real
geprüft: beide Ziele brechen mit Exit 2 ab, bevor ein `docker run`
versucht wird). Anders als `example-run-go` bauen sie **nichts** — sie
starten den bereits von `make examples-csharp`/`make examples-kotlin`
gebauten Image-Tag `pg-change-feed-examples:csharp[-<surface>]`/
`:kotlin[-<surface>]` real gegen das Docker-Netzwerk `cdc-examples` mit
`--env-file examples/.env`. Real gegen die laufende Demo-Umgebung
(`make example-demo-up`) verifiziert: `make example-run-csharp SURFACE=http
ARGS="--source demo-source --publication pub_demo"` und dieselbe
Kotlin-Form liefern beide real `{"tables":[{"table_id":"tbl-orders",…}]}`
von `GET /tables` gegen den laufenden Feed-Container; `SURFACE=sse` startet
für beide Sprachen real einen laufenden Container gegen `cdc-examples`
(Tag-Auflösung auf `*-sse` bestätigt). `make examples-csharp`/
`make examples-kotlin` real neu gebaut (Layer-Cache-Treffer, Exit 0) —
unverändert. `examples/README.md` (C#-/Kotlin-Tabellen, Spalte „Start" auf
`make example-run-*`-Aufruf statt rohem `docker run`), die vier
`**Beispiele:**`-Blöcke in `docs/user/benutzerhandbuch.md` (samt
Versionshistorie, 1.26 → 1.27) und `harness/README.md` §Werkzeuge (neue
Zeile, die vormalige „noch nicht implementiert"-Notiz bei `example-run-go`
entfernt) zitieren die neue Startform. `make gates` zweimal grün gelaufen
(vor und nach den Doku-Änderungen, Exit-Code beide Male direkt und
ungepiped geprüft).

**Steering-Loop-Lerneintrag:** Der von `example-run-go` etablierte
Pflicht-`SURFACE=`-Guard (`ifeq $(filter …) $(error …)`) trägt unverändert,
wenn das Ziel — wie hier — keinen eigenen Bau-Schritt hat: die Guard-Zeile
bleibt identisch, nur der `docker build`-Aufruf entfällt. Ein fehlender
Bau-Schritt scheitert erwartungsgemäß real und sichtbar am `docker
run`-eigenen „image not found", statt eines stillen Vorab-Baus — kein neuer
Mechanismus nötig, keine Überraschung gegenüber der ADR-Festlegung.

**Beobachtungs-Register:** keine Beobachtung angefallen — die Umsetzung
folgte Festlegung 2 der ADR ohne Abweichung, die eine neue oder zitierbare
`BEO-PGC/*`-Musterlücke begründen würde.

**Risiken (§6):** beide mit Ausgang „entfallen" versehen, real gemessen
(Compose-Servicename statt `localhost` in `examples/.env`, zwei Ziele statt
Zerlegung) — siehe §6.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung.

**Vorgelagert — Sub-Area-Wahl prüfen:** berührte Sub-Areas sind
`examples/csharp/**`, `examples/kotlin/**` (nur die Werkzeug-Ziel-Ebene, kein
Quellcode-Eingriff), `harness/mk/**`, `docs/user/**` — alle Default-Sub-Area
des Repos, keine feinere Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`beispiel`/`example`/`start`/`make-target`) — **keine Treffer** (siehe
Welle-Datei §1).

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.

### Sub-Area: `*` (Default)

- **Modus:** GF — `harness/conventions.md` §Modus-Deklaration.
