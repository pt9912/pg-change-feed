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

- [ ] **LP1** — Neue ADR mit `Supersedes ADR-0076` in genau der
      Startform-Klausel geschrieben, mindestens drei Alternativen für die
      Go-Bauform verglichen (Pro/Contra), `Accepted`, ADR-Index
      (`docs/plan/adr/README.md`) aktualisiert.
- [ ] **LP2** — Start-Make-Target-Form für alle drei Sprachen entschieden
      und begründet (mindestens drei Alternativen verglichen), inklusive der
      Fitness-Function-Zeile (welches `make`-Ziel/welcher Aufruf startet
      welches Programm).
- [ ] **LP3** — Umgebungsdatei-/Variablen-Kontrakt der Demo-Umgebung
      entschieden (Dateiname/-pfad, Variablen-Namen, wer sie liest:
      Compose-Datei und/oder Start-Targets) samt Bestätigung oder Korrektur
      der Slice-Schnitt-Empfehlung für die drei Folge-Slices.
- [ ] `make gates` grün.
- [ ] Review durchgeführt (Konsistenzprüfung der neuen ADR gegen
      `ADR-0076`/`ADR-0087`/`ADR-0090`), Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: keiner erwartet (kein öffentlicher Vertrag geändert — die
      ADR ist das Erzeugnis); falls die ADR `SPEC-023` schärft, trägt sie
      selbst die Folgepflicht-Zeile für den umsetzenden Zug.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
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
  **Ausgang:** <bei Closure ausfüllen; falls eingetreten, benennt die ADR
  selbst die zusätzlichen Slices in ihrer Slice-Schnitt-Empfehlung>.
- Ein einzelnes, parametrisiertes Start-Target könnte über drei so
  unterschiedliche Bauformen (Go-Binary, C#/Kotlin-Container-Images) nicht
  einheitlich funktionieren und müsste doch je Sprache aufgeteilt werden —
  **Ausgang:** <bei Closure ausfüllen>.
- Der Umgebungsdatei-Kontrakt könnte mit bestehenden, bereits dokumentierten
  `CDC_*`-Variablennamen kollidieren (Doppeldeutigkeit Demo-Default vs.
  Produktions-Konfiguration) — **Ausgang:** <bei Closure ausfüllen>.

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

*(bei Closure zu füllen — die Welle ist gerade erst eröffnet)*

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
