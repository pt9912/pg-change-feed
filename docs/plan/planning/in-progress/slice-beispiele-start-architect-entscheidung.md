# Slice beispiele-start-architect-entscheidung: Architect-Entscheidung — Supersedes-ADR zu `ADR-0076` Festlegung 1 (Go-Startform) + Ausgestaltung des Start-Make-Targets

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** welle-beispiele-start-ueber-make — die gemeinsame
Design-Entscheidung (Go-Bauform, Start-Target-Form über drei Sprachen) ist das
*Mehr*, das kein Einzel-Slice-DoD trägt (siehe Welle-Datei §1).

**Bezug:** [`LH-FA-SST-006`](../../../../spec/lastenheft.md),
[`LH-FA-SST-007`](../../../../spec/lastenheft.md),
[`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Scope: die Oberflächen, die
die Beispiele zeigen — diese ADR ändert sie nicht),
[`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
(in der Startform-Klausel zu superseden), [`ADR-0087`](../../adr/0087-beispiel-clients-csharp-kotlin.md),
[`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) (Bestand, den die
neue ADR einordnen muss).

**Berührte Spec-Stellen:** [`SPEC-023`](../../../../spec/pflichtenheft.md)
(Beispiel-Clients — trägt ggf. eine Folgepflicht-Zeile, wenn die Startform
Teil der dortigen Festlegung ist; die neue ADR entscheidet und benennt das).

**Verantwortlich:** — (Architect-Rolle, unabhängiger Zug — siehe §3.5-Hinweis
unten).

**Autor:** Planner-Rolle. **Datum:** 2026-09-18.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Eine neue, `Accepted`-ADR entscheidet (a) die Rücknahme von
[`ADR-0076`](../../adr/0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md)
§Entscheidung Festlegung 1, Startform-Bullet — Go-Beispiele werden künftig
über `make` + Dockerfile gebaut/gestartet statt über `go run` — samt
Begründung, ob dafür das C#/Kotlin-Sprach-Wurzel-Muster (eigenes Dockerfile
je Sprach-Wurzel) übernommen oder eine leichtere, Go-spezifische Bauform
gewählt wird; (b) die Form eines echten Start-Make-Targets über alle drei
Sprachen (Go, C#, Kotlin) — ein Ziel je Programm, ein parametrisiertes Ziel,
oder ein Ziel je Sprache mit Pflicht-Argumenten; und (c) den
Umgebungsdatei-/Variablen-Kontrakt einer neuen Demo-Umgebung unter
`examples/` (Compose-Datei + gemeinsame Umgebungsdatei), gegen den sowohl
die Demo-Umgebung als auch beide Start-Target-Implementierungen schreiben.
Die ADR schließt mit einer Slice-Schnitt-Bestätigung oder -Korrektur für die
drei Folge-Slices dieser Welle (§4 der Welle-Datei).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Der tatsächliche Go-Dockerfile-Bau und das Start-Target selbst** —
  `slice-beispiele-go-dockerfile-start` übernimmt die Umsetzung der hier
  entschiedenen Form; dieser Slice trifft die Entscheidung, nicht den Code.
- **Die Start-Target-Implementierung für C#/Kotlin** —
  `slice-beispiele-csharp-kotlin-start-target` übernimmt.
- **Die Demo-Compose-Datei, die gemeinsame Umgebungsdatei und das
  Bootstrapping (Schema-Rollout, Beispiel-Quelle/-Tabelle) selbst** —
  `slice-beispiele-compose-bootstrap` übernimmt die Umsetzung des hier
  entschiedenen Kontrakts.
- **Änderungen am `.proto`-Weg oder an `gen/**`** — bleibt Bestand aus
  [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) Festlegung 2/3,
  sofern die Analyse in diesem Slice keinen neuen `.proto`-Konsumenten
  findet (nicht erwartet: ein Start-Mechanismus und eine Demo-Umgebung lesen
  die `.proto` nicht).
- **Neue Zugriffs-Oberflächen oder Sprachen** — die Matrix bleibt, was
  [`ADR-0090`](../../adr/0090-beispiel-clients-volle-matrix.md) festgelegt
  hat; dieser Slice ändert **wie**, nicht **was**.
- **Der ausgelieferte Bau** (Wurzel-`Dockerfile`, `.dockerignore`,
  `make image`) — bleibt unberührt, wie in allen drei bestehenden
  Beispiel-ADRs festgelegt; kein Anlass, das hier zu öffnen.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] **LP1** — Neue ADR [`ADR-0098`](../../adr/0098-beispiel-clients-start-ueber-make-dockerfile.md)
      mit `Supersedes ADR-0076` in genau der Startform-Klausel geschrieben,
      vier Alternativen für die Go-Bauform verglichen (Pro/Contra), `Accepted`,
      ADR-Index (`docs/plan/adr/README.md`) aktualisiert.
- [x] **LP2** — Start-Make-Target-Form für alle drei Sprachen entschieden
      und begründet (vier Alternativen verglichen, §B der ADR), inklusive der
      Fitness-Function-Zeile (welches `make`-Ziel/welcher Aufruf startet
      welches Programm).
- [x] **LP3** — Umgebungsdatei-/Variablen-Kontrakt der Demo-Umgebung
      entschieden (Dateiname/-pfad, Variablen-Namen, wer sie liest:
      Compose-Datei und/oder Start-Targets) samt Bestätigung (mit einer
      Namens-Korrektur, `examples/.env` statt `examples/.env.example`) der
      Slice-Schnitt-Empfehlung für die drei Folge-Slices.
- [x] `make gates` grün (ungepiped Exit-Code geprüft, `AGENTS.md` §3.9 —
      siehe Bericht).
- [ ] Review durchgeführt (Konsistenzprüfung der neuen ADR gegen
      `ADR-0076`/`ADR-0087`/`ADR-0090`), Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      **Offen** — dieser Zug ist die Architect-Rolle selbst; der
      Rollenwechsel zum Reviewer ist der nächste Schritt, kein
      Self-Review-Ersatz. Der `git mv` nach `done/` steht deshalb noch aus.
- [x] Doku-Update: keiner erwartet (kein öffentlicher Vertrag geändert — die
      ADR ist das Erzeugnis); die ADR trägt selbst die Folgepflicht-Zeile für
      `SPEC-023` (umsetzender Zug).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag — siehe §7 (Notiz
      vorbereitet; die eigentliche Lifecycle-Closure folgt nach dem
      Reviewer-Pass).
- [x] Reconciliation-Register (`../reconciliation.md`) — Datei existiert in
      diesem Repo nicht (kein Brownfield-Bootstrap); Item entfällt.
- [x] Beobachtungs-Register (`../observations/`) — keine Beobachtung
      angefallen, siehe §7.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `docs/plan/adr/00XX-<slug>.md` | neu | die Supersedes-ADR selbst (LP1/LP2/LP3) |
| `docs/plan/adr/README.md` | update | Index-Zeile für die neue ADR |

**Ansatz** (kein Code, kein Test — die Architect-Rolle liefert eine
Entscheidung, keine Implementierung):

- Real messen statt annehmen, wo möglich: Ist ein Go-Binary tatsächlich
  billiger cross-kompilierbar als ein C#/Kotlin-Werkzeugketten-Image (Digest,
  Bauzeit-Größenordnung)? Trägt ein einzelnes, parametrisiertes
  `make example-run`-Ziel über drei so unterschiedliche Bauformen (Go-Binary
  vs. Container-Images), oder ist ein Ziel je Sprache robuster? Reicht ein
  `--env-file` gegen dieselbe Umgebungsdatei für Go (`go run`, kein
  Container) genauso wie für C#/Kotlin (`docker run --env-file`)?
- Mindestens drei Alternativen je Entscheidungsachse (Go-Bauform,
  Start-Target-Form, Umgebungsdatei-Kontrakt) mit Pro/Contra, wie es
  `ADR-0076`/`ADR-0087`/`ADR-0090` vormachen.

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): sofort bei Welle-Eröffnung — dieser Slice
ist der erste, blockierende Slice der Welle; kein externer Trigger
erforderlich außer der Welle-Eröffnung selbst.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Die Entscheidung
  verlangt einen größeren strukturellen Umbau als die drei Folge-Slices
  tragen können (z. B. ein sprachneutrales Meta-Bau-Framework über alle drei
  Sprachen) — dann wird die ADR selbst kleiner geschnitten (z. B. Go-Startform
  separat von der Demo-Umgebung).
- `in-progress` → `open` (blockiert — Carveout?): Der Auftraggeber revidiert
  eine der drei Entscheidungen, bevor die ADR `Accepted` ist.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

- Die neue ADR liegt `Accepted` vor, mit `Supersedes ADR-0076` in genau der
  Startform-Klausel und einer Slice-Schnitt-Bestätigung/-Korrektur für die
  drei Folge-Slices.
- Review (Konsistenzprüfung) liegt vor, keine offene HIGH-Finding.
- Closure-Notiz mit Lerneintrag (§7).

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Die Go-Bauform-Entscheidung könnte einen vollständigen Strukturumbau der
  drei bestehenden Go-Beispiele auf das C#/Kotlin-Sprach-Wurzel-Muster nach
  sich ziehen (breiter als ein einzelner Folge-Slice tragen kann) —
  **Ausgang: entfallen.** `ADR-0098` Festlegung 1/§A wählt bewusst die
  **leichtere**, Go-spezifische Form (ein `examples/Dockerfile`, Wurzel-Kontext,
  dateispezifisches Ignore) statt des Sprach-Wurzel-Umzugs (Option A3
  verworfen) — kein Strukturumbau der drei bestehenden Go-Beispiele nötig.
- Ein einzelnes, parametrisiertes Start-Target könnte über drei so
  unterschiedliche Bauformen (Go-Binary, C#/Kotlin-Container-Images) nicht
  einheitlich funktionieren und müsste doch je Sprache aufgeteilt werden —
  **Ausgang: eingetreten**, und in der ADR selbst aufgelöst. Die reale Probe
  im Zug bestätigte genau dieses Risiko (§B der ADR, Option B2 verworfen):
  `ADR-0098` Festlegung 2 wählt **ein Ziel je Sprache** mit Pflicht-Argument
  `SURFACE=`, kein sprachübergreifend parametrisiertes Ziel.
- Der Umgebungsdatei-Kontrakt könnte mit bestehenden, bereits dokumentierten
  `CDC_*`-Variablennamen kollidieren (Doppeldeutigkeit Demo-Default vs.
  Produktions-Konfiguration) — **Ausgang: entfallen.** `ADR-0098`
  Festlegung 3 stellt fest: die Wiederverwendung derselben `CDC_*`-Namen ist
  **Absicht**, keine Kollision — die Trennung zwischen Demo und Produktion
  liegt in der Netzwerk-Isolation (Festlegung 4, eigenes Docker-Netzwerk
  `cdc-examples`), nicht im Variablennamen.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks). Ging der Gegenstand an einen anderen Slice oder entfiel er, trägt
diese Sektion die Zeile `Gegenstand:` mit Kennung oder Grund und jedes Risiko
aus §6 seinen Ausgang; die Liefer-Punkte der DoD bleiben leer
(`modul-05-planning-harness.md` §Ein Slice, dessen Gegenstand ein anderer
übernimmt).

**Gegenstand:** vollständig geliefert — dieser Slice bleibt zunächst in
`in-progress/`, nicht `done/`: die Lifecycle-Closure (§5) verlangt einen
Review-Report, und dieser Zug ist die Architect-Rolle selbst
(Modul 8 §Rollen-Regeln, kein Self-Review). Der `git mv` nach `done/` folgt,
sobald der Reviewer-Pass gegen
[`ADR-0098`](../../adr/0098-beispiel-clients-start-ueber-make-dockerfile.md)
vorliegt.

**Ergebnis:** [`ADR-0098`](../../adr/0098-beispiel-clients-start-ueber-make-dockerfile.md)
(`Accepted`, `Supersedes ADR-0076` in genau der Startform-Klausel)
entscheidet: (1) Go-Beispiele bauen/starten künftig über ein eigenes
`examples/Dockerfile` (Wurzel-Kontext, dateispezifisches
`examples/Dockerfile.dockerignore` — Isolationsmechanismus real gemessen in
diesem Zug, `docker` v29.8.0, kein Eingriff in die Wurzel-`Dockerfile`/
`.dockerignore`); (2) alle drei Sprachen bekommen einen echten
Start-Make-Target — ein Ziel je Sprache (`example-run-go`,
`example-run-csharp`, `example-run-kotlin`) mit Pflicht-Argument `SURFACE=`,
gegen ein festes, benanntes Docker-Netzwerk `cdc-examples`; (3) die
gemeinsame Umgebungsdatei heißt `examples/.env` (committet, sofort nutzbar,
bestehende `CDC_*`-Namen — die Planner-Lesart war richtig, der Dateiname
wurde korrigiert). Die drei Folge-Slices (`slice-beispiele-compose-bootstrap`,
`slice-beispiele-go-dockerfile-start`,
`slice-beispiele-csharp-kotlin-start-target`) sind mit den festen Namen aus
der ADR nachgezogen; ihr Schnitt (drei unabhängige, parallele Slices) bleibt
bestätigt.

**Steering-Loop-Lerneintrag:** Ein `Accepted`-ADR-Widerruf durch den
Auftraggeber ist kein Sonderfall, den die Immutabilitäts-Regel schwächt —
er ist der Regelfall, für den `Supersedes` existiert (`AGENTS.md` §3.5). Die
reale Docker-Probe (dateispezifisches `.dockerignore`, real gemessen statt
nur dokumentiert zitiert) war der Unterschied zwischen einer plausiblen und
einer geprüften Bauform-Entscheidung — dieselbe Disziplin, die
`ADR-0090` bereits für den `.proto`-Zusatzkontext vorgemacht hat, hier zum
ersten Mal auf den Go-Bauweg angewendet.

**Beobachtungs-Register:** keine Beobachtung angefallen — die drei
entschiedenen Fragen sind Einzelfall-Entscheidungen dieses Zugs, keine
wiederkehrende Musterlücke, die einen neuen `BEO-PGC/*`-Eintrag rechtfertigt.

**Risiken (§6):** alle drei mit Ausgang versehen (zwei entfallen, eines
eingetreten und in der ADR selbst aufgelöst) — siehe §6.

## 8. Sub-Area-Prüfungen und Modus-Begründung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Sub-Area-Modus-Begründung — dort die **zwei vorgelagerten
Schritte** (sie stehen in jedem Slice-Plan, unabhängig von Modus und
Slice-Typ) und die **vier Pflichtkriterien** (Konventionen-Dichte ·
Phase-Reife · Evidenz-/Diskrepanz-Risiko · Reconciliation-Aufwand), vier und
nicht mehr.

**Der Abschnitt selbst entfällt nie.** Die zwei vorgelagerten Prüfungen laufen
in **jedem** Slice-Plan — sie hängen weder am Modus noch am Slice-Typ. Bedingt
ist allein der Modus-Begründungsblock am Ende; deshalb nennt der Titel beide
Hälften.

**Vorgelagert — Sub-Area-Wahl prüfen:** einzige berührte Sub-Area ist
`docs/plan/adr/**` (ADR-Verzeichnis) — Default-Sub-Area des gesamten Repos
nach `harness/conventions.md` §Modus-Deklaration, keine feinere
Ausdifferenzierung nötig (kein Code-Pfad berührt).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen
(`beispiel`/`example`/`start`/`make-target`/`compose`) — **keine Treffer**
(siehe Welle-Datei §1).

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF.

### Sub-Area: `*` (Default, `docs/plan/adr/**`)

- **Modus:** GF — `harness/conventions.md` §Modus-Deklaration führt das
  gesamte Repo als Greenfield; ein ADR-Verzeichnis ohne Code-Bestand-Frage.
