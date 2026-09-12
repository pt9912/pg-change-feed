# Slice slice-021: Consumer-Registrierung: Zugriffsweg (ADR) + Verdrahtung

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-6`](../welle-6.md) — die Ende-zu-Ende-Belegpflicht (externer
Consumer durchläuft Registrierung + Bestätigung ausschließlich über den
neuen Zugriffsweg) übersteigt, was dieser Slice allein beweist.

**Bezug:** [`LH-FA-CON-001.a`](../../../../spec/pflichtenheft.md) (Zugriffsweg
offen), [`LH-FA-CON-001`](../../../../spec/lastenheft.md) (Registrierung
benannter Consumer), [`ADR-0019`](../../adr/0019-cli-driving-adapter.md)
(CLI-Muster, trägt den Zugriffsweg), [`ADR-0020`](../../adr/0020-http-grpc-optional.md)
(Netzwerkschnittstelle ohne Bedarf ausgeschlossen), [`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
(SQL-Funktions-Weg physisch ungelöst, deshalb ausgeschlossen). Architect-Verdikt:
[`architect-review-slice-021.md`](../../adr/architect-review-slice-021.md)
— kein neues ADR nötig, Zugriffsweg = CLI-Unterbefehl.

**Berührte Spec-Stellen:** [`ARC-003`](../../../../spec/architecture.md)
(Inbound Ports), [`ARC-005`](../../../../spec/architecture.md) (Driving
Adapters — nennt CLI/SQL-Funktionen/Views/HTTP-gRPC explizit als Optionen).
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Verantwortlich:** pt9912.

**Autor:** pt9912. **Datum:** 2026-09-12.

---

## 1. Ziel und Abgrenzung

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Ein von außen erreichbarer Zugriffsweg (Architect entscheidet
CLI-Unterbefehl vs. Netzwerkschnittstelle vs. anderer Mechanismus als ADR)
ruft den bestehenden `RegisterConsumerUseCase`
(`internal/application/usecase/register`) auf — ein externer Consumer kann
sich damit registrieren, ohne die CDC-Speichertabellen direkt zu
beschreiben.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Positions-Bestätigung über denselben Zugriffsweg** — Folge-Slice
  `slice-022`: eigene Fähigkeit (`LH-FA-CON-004.a`), eigener Use Case
  (`AcknowledgeConsumerUseCase`), baut aber auf dem hier entschiedenen
  Zugriffsweg auf.
- **Administrative Entfernung von Consumern**
  ([`LH-FA-CON-006`](../../../../spec/lastenheft.md)) über denselben
  Zugriffsweg — Bestand bleibt bewusst stehen: keine benannte Spec-Lücke,
  kein Bedarf in dieser Welle (§6 der Welle).
- **Least-Privilege-Rollen-Zuordnung für den neuen Zugriffsweg**
  (`BEO-PGC/rollen-verdrahtung`) — anderer Vorgang: Die bestehende
  Instanz-DSN-Verdrahtung in `internal/bootstrap/wiring.go` wird von
  diesem Slice nicht aufgelöst, der neue Zugriffsweg reiht sich dort ein.

**Keine Mindestzahl.** Ein Slice mit *einem* echten Ausschluss ist besser als
einer mit vier erfundenen; die vier Klassen sind ein Suchraster, keine
Ausfüll-Liste. Suchreihenfolge: Was übernimmt ein **Folge-Slice** (mit
Kennung — und die Kennung muss den Punkt auch annehmen)? Was bleibt als
**Bestand** bewusst stehen (mit Begründung)? Was wäre ein **anderer Vorgang**?
Welche **Schicht** rührt der Slice nicht an?

Was hier steht, ist die Grenze, an der ein wachsender Slice sich messen lässt:
Wer später etwas mitnimmt, das hier ausgeschlossen war, hat den Plan
**geändert**, nicht nur ergänzt.

## 2. Definition of Done

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] ADR entschieden (Architect): Zugriffsweg für Consumer-Registrierung
      (CLI-Unterbefehl / Netzwerkschnittstelle / anderer Mechanismus),
      referenziert [`LH-FA-CON-001.a`](../../../../spec/pflichtenheft.md).
      Architect-Verdikt [`architect-review-slice-021.md`](../../adr/architect-review-slice-021.md):
      CLI-Unterbefehl, kein neues ADR nötig.
- [x] [`LH-FA-CON-001.a`](../../../../spec/pflichtenheft.md) erfüllt: ein
      externer Aufruf über den entschiedenen Zugriffsweg registriert einen
      Consumer über `RegisterConsumerUseCase`, ohne die CDC-Speichertabellen
      direkt zu schreiben — Test referenziert:
      `internal/bootstrap/register_test.go::TestRegisterConsumerEndToEnd`
      (Happy Path + „bereits registriert"-Boundary, real gegen PostgreSQL
      über `make test-store`).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      `docs/reviews/review-slice-021.md`: 0 HIGH, 1 MEDIUM (F-1,
      Negativtest-Lücke) — geschlossen in dieser Fixrunde
      (`internal/bootstrap/register_test.go::TestRegisterConsumerReportsDomainFailure`).
- [x] Doku-Update für `docs/user/benutzerhandbuch.md` (neuer Zugriffsweg für
      Consumer-Registrierung) und den ADR-Index. ADR-Index unverändert —
      dieser Slice legt kein neues ADR an (Architect-Verdikt).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. Kein Eintrag verkörpert
      (Normalfall) — siehe §7.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item. Geprüft: Datei existiert nicht (GF-Repo) — entfällt.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
      `BEO-PGC/rollen-verdrahtung/evidence/slice-021.md` ergänzt, Zähler 2×.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen). Ausgang: weiter offen.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit). Dieser Slice gehört zu `welle-6` — die Paarungen prüft die **Welle-Closure**, nicht dieser Slice. Item entfällt hier bewusst.

## 3. Plan (vor Code)

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| ~~`docs/plan/adr/NNNN-consumer-zugriffsweg.md`~~ | **entfällt** | Architect-Verdikt [`architect-review-slice-021.md`](../../adr/architect-review-slice-021.md): drei bereits `Accepted`-ADRs (`ADR-0019`, `ADR-0020`, `ADR-0046`) schließen den Optionsraum erschöpfend — kein neues ADR nötig |
| `cmd/pg-change-feed/main.go` | update | neuer Sondermodus `register-consumer <name>` (Unterbefehl), verdrahtet bis `RegisterConsumerUseCase`, Exit-Codes 0/1/2 |
| `internal/bootstrap/wiring.go` | update | `RegisterConsumer(ctx, cfg, name)` baut `postgresstorage.NewConsumerState` + `register.NewRegisterConsumerService` und ruft `Register` auf |
| `internal/bootstrap/register_test.go` | neu | End-to-End-Test gegen reale PostgreSQL (Happy Path + „bereits registriert"-Boundary) und ein netzloser Verdrahtungsfehler-Test — nicht in der ursprünglichen Plan-Tabelle, Plan-Nachzug im selben Lauf |
| `docs/user/benutzerhandbuch.md` | update | neuer Zugriffsweg für Consumer-Registrierung dokumentiert |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): Priorisiert, `Verantwortlich:` gesetzt,
WIP-Limit (1 je Implementer) frei, `welle-6` eröffnet.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Falls die
  ADR-Entscheidung eine Netzwerkschnittstelle mit eigenem Protokoll
  verlangt (statt CLI/bestehendem SQL-Kanal) — dann berührt der Slice mehr
  als zwei Schichten und wird neu geschnitten.
- `in-progress` → `open` (blockiert — Carveout?): Architect trifft keine
  Entscheidung innerhalb dieses Laufs (Konflikt-Sequenz, Modul 8) —
  unwahrscheinlich, da nur ein Rolleninhaber beteiligt ist.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

ADR accepted **und** ein externer Aufruf über den entschiedenen Zugriffsweg
registriert real einen Consumer (Test grün) **und** `make gates` ist grün.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der neue Zugriffsweg könnte dieselbe gemeinsame Instanz-DSN nutzen wie
  Store-/Aktivierungs-/Stream-Verbindung, statt der `cdc_admin`-Rolle
  (`BEO-PGC/rollen-verdrahtung`) — Folge des bestehenden, unveränderten
  Verdrahtungsstands, keine Neuverschärfung durch diesen Slice.
  **Ausgang: weiter offen** → `BEO-PGC/rollen-verdrahtung`
  (`evidence/slice-021.md` ergänzt, Zähler jetzt 2×).

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** Der Architect-Verdikt-Weg (Verdikt statt neues
  ADR, wenn bereits Accepted-ADRs den Optionsraum abdecken) hat einen ganzen
  Zug ohne Folge-ADR-Overhead getragen — drei bereits bestehende ADRs
  (`ADR-0019`/`0020`/`0046`) reichten aus, die Zugriffsweg-Frage
  erschöpfend zu entscheiden. Der Verifier hat den vom Implementer
  offengelassenen Mutationstest-Vorbehalt eigenständig geschlossen (real
  reproduziert, real rot).
- **Was ging anders als geplant:** `internal/bootstrap/register_test.go`
  war nicht im ursprünglichen §3-Plan (Plan-Nachzug im selben
  Implementer-Lauf). Der Reviewer fand F-1 (MEDIUM, Negativtest-Lücke für
  den `Register`-Domänenfehler-Pfad), geschlossen in einer eigenen
  Fixrunde (`ef99df8`). Ein themenfremder Lastenheft-CR (`9936e82`,
  `LH-FA-SST-006`) landete zeitlich zwischen den Slice-Commits — der
  Verifier hat das als V-1 (LOW, non-blocking) korrekt nicht diesem Slice
  angelastet, sondern der `welle-6`-Trigger-Audit zugeordnet (siehe
  `welle-6.md` §Vermerk für den Trigger-Audit bei Closure).
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-021.md`
  in `BEO-PGC/rollen-verdrahtung` ergänzt — Zähler steht jetzt bei 2×.
- **Folge-Slices:** `slice-022` (Positions-Bestätigung über denselben
  Zugriffsweg) — liegt in `open/`.
- **Risiken aus §6:** einziges Risiko (gemeinsame Instanz-DSN statt
  `cdc_admin`) — Ausgang *weiter offen* → `BEO-PGC/rollen-verdrahtung`.
- **Drei Paarungen:** Dieser Slice gehört zu `welle-6` — die Paarungen
  prüft die Welle-Closure, nicht dieser Slice (Modul 6/8).

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

**Vorgelagert — Sub-Area-Wahl prüfen:** Berührte Sub-Area ist `*` (Default,
`PGC`, Greenfield laut `harness/conventions.md`) — die einzige deklarierte
Sub-Area dieses Repos, Schwelle also trivial erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register (`../observations/`)
durchgegangen — `BEO-PGC/rollen-verdrahtung` (1×) betrifft dieselbe
Bootstrap-Verdrahtung, die der neue Zugriffsweg mitnutzt (siehe §6-Risiko);
kein Eintrag betrifft Consumer-Zugriffswege direkt.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (neue
Fähigkeit auf bestehender Greenfield-Codebasis, keine Inventur-Diskrepanz
zu erwarten).
