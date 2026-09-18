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

- [ ] **LP1** — Start-Make-Ziel(e) für die vier C#-Programme in der von der
      ADR entschiedenen Form, gegen die im Umgebungsdatei-Kontrakt
      vereinbarten Variablen.
- [ ] **LP2** — dieselbe Form für die vier Kotlin-Programme.
- [ ] **LP3** — Träger nachgezogen: `examples/README.md` (C#/Kotlin-Tabellen,
      Startform-Spalte auf `make`-Aufruf statt rohem `docker run`), die vier
      Zugriffs-Abschnitte in `docs/user/benutzerhandbuch.md` (C#-/Kotlin-Zeilen
      der `**Beispiele:**`-Blöcke), `harness/README.md` §Werkzeuge.
- [ ] [`LH-FA-SST-006`](../../../../spec/lastenheft.md)/[`LH-FA-SST-007`](../../../../spec/lastenheft.md)/[`LH-FA-SST-008`](../../../../spec/lastenheft.md)
      weiterhin gezeigt (kein Verhaltens-, nur
      Start-Mechanismus-Wechsel) — `make examples-csharp`/
      `make examples-kotlin` bleiben unverändert grün.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für die C#/Kotlin-Startform (öffentlicher Vertrag) — siehe
      LP3.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`, oder „keine Beobachtung
      angefallen" in §7.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
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
  (Adresse `localhost` vs. Compose-Service-Name) — **Ausgang:** <bei Closure
  ausfüllen>.
- Acht Programme über zwei Sprachen könnten die Drei-Liefer-Punkte-Grenze
  real sprengen, falls die Start-Target-Form pro Programm statt pro Sprache
  verdrahtet werden muss — **Ausgang:** <bei Closure ausfüllen; siehe
  Rückführung in §4>.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register · `grundlagen-traceability.md` §Herkunfts-Anker
für Steering-Loop-Regeln.

*(bei Closure zu füllen)*

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
