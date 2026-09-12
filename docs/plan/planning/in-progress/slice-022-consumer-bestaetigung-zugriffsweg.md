# Slice slice-022: Positions-Bestätigung über denselben Zugriffsweg

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** [`welle-6`](../welle-6.md) — die Ende-zu-Ende-Belegpflicht (externer
Consumer durchläuft Registrierung + Bestätigung ausschließlich über den
neuen Zugriffsweg) übersteigt, was dieser Slice allein beweist.

**Bezug:** [`LH-FA-CON-004.a`](../../../../spec/pflichtenheft.md) (Zugriffsweg
offen, Vorwärts-Invariante nicht durchgesetzt),
[`LH-FA-CON-004`](../../../../spec/lastenheft.md) (Bestätigung einer Position),
[`ADR-0019`](../../adr/0019-cli-driving-adapter.md) (CLI-Muster, trägt den
Zugriffsweg), [`ADR-0020`](../../adr/0020-http-grpc-optional.md)
(Netzwerkschnittstelle ohne Bedarf ausgeschlossen), [`ADR-0046`](../../adr/0046-sql-driving-adapter-lese-schreib-trennung.md)
(SQL-Funktions-Weg physisch ungelöst, deshalb ausgeschlossen).
Architect-Verdikt: [`architect-review-slice-021.md`](../../adr/architect-review-slice-021.md)
§3 — dieselbe, für `slice-021` bereits entschiedene Zugriffsweg-Mechanik
(CLI-Unterbefehl) gilt kanalgenerisch auch hier, kein erneuter
Architect-Rundlauf für dieselbe Frage.

**Berührte Spec-Stellen:** [`ARC-003`](../../../../spec/architecture.md)
(Inbound Ports), [`ARC-005`](../../../../spec/architecture.md) (Driving
Adapters).
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

**Ziel:** Derselbe, in `slice-021` entschiedene Zugriffsweg ruft zusätzlich
den bestehenden `AcknowledgeConsumerUseCase`
(`internal/application/port/inbound/consumer.go`) auf — ein externer
Consumer kann eine Position bestätigen, ohne die Consumer-Positions-Tabelle
direkt zu schreiben, und die Vorwärts-Invariante
(`ErrPositionRegression`) ist damit für jeden externen Aufruf durchgesetzt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Ein neuer Zugriffsweg-Mechanismus** — Bestand bleibt bewusst stehen:
  Dieser Slice nutzt den in `slice-021` entschiedenen und verdrahteten
  Zugriffsweg, trifft keine eigene ADR.
- **Härtung des Datenbankschemas selbst gegen direktes Schreiben**
  (z. B. Constraint/Trigger auf `cdc.consumer_position`) — anderer Vorgang:
  `LH-FA-CON-004.a` verlangt einen funktionierenden Zugriffsweg, keine
  zusätzliche Schema-Sperre; ein direktes Schreiben am Zugriffsweg vorbei
  bleibt möglich, wie bei jeder anderen Tabelle auch.

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

- [x] [`LH-FA-CON-004.a`](../../../../spec/pflichtenheft.md) erfüllt: ein
      externer Aufruf über den Zugriffsweg aus `slice-021` bestätigt eine
      Position über `AcknowledgeConsumerUseCase`, ohne die
      Consumer-Positions-Tabelle direkt zu schreiben — Test referenziert:
      `internal/bootstrap/acknowledge_test.go::TestAcknowledgeConsumerEndToEnd`
      (Happy Path, real gegen PostgreSQL über `make test-store`).
- [x] Vorwärts-Invariante (`ErrPositionRegression`) real über den externen
      Zugriffsweg getestet (Wiederholung/Rückschritt abgelehnt) —
      `internal/bootstrap/acknowledge_test.go::TestAcknowledgeConsumerEndToEnd`
      (Wiederholung derselben Position idempotent, echter Rückschritt
      abgelehnt, gespeicherter Stand bleibt unverändert; real gegen
      PostgreSQL, verbose Einzellauf bestätigt: `--- PASS:
      TestAcknowledgeConsumerEndToEnd`).
- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update für `docs/user/benutzerhandbuch.md` (neuer Zugriffsweg für
      Positions-Bestätigung).
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
| `cmd/pg-change-feed/main.go` | update | Zugriffsweg aus `slice-021` um Bestätigung erweitert |
| `internal/bootstrap/wiring.go` | update | ruft `AcknowledgeConsumerUseCase` über den Zugriffsweg |
| `docs/user/benutzerhandbuch.md` | update | Bestätigung über den Zugriffsweg dokumentiert |
| `internal/bootstrap/acknowledge_test.go` | neu | End-to-End-Test gegen reale PostgreSQL (Happy Path, Wiederholung/Idempotenz, echter Rückschritt gegen `ErrPositionRegression`, Registrierungs-Grenze) und zwei netzlose Verdrahtungsfehler-Tests — nicht in der ursprünglichen Plan-Tabelle, Plan-Nachzug im selben Lauf (dasselbe Muster wie `register_test.go` in `slice-021`) |
| `tools/harness/run-store-tests.sh` | update | `internal/bootstrap` läuft vorgezogen in einem eigenen `go test`-Aufruf, getrennt vom übrigen Paket-Bündel — nicht im ursprünglichen Plan, Plan-Nachzug im selben Lauf: mehrere `postgresstorage`-Tests räumen `cdc.consumer`/`cdc.consumer_position` tabellenweit ab bzw. per `DROP SCHEMA cdc CASCADE` neu auf; ein paralleler Lauf (Go-Default über Pakete hinweg) ließ den neuen Bestätigungs-Test real und reproduzierbar an einer weggerissenen bzw. entfernten Zeile scheitern, ohne dass der Zugriffsweg selbst fehlerhaft war |

## 4. Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): `slice-021` liegt in `done/`
(Zugriffsweg entschieden und verdrahtet), WIP-Limit (1 je Implementer) frei.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): Nicht erwartet —
  der Zugriffsweg-Mechanismus ist durch `slice-021` bereits entschieden,
  dieser Slice erweitert ihn nur um einen zweiten Aufruf.
- `in-progress` → `open` (blockiert — Carveout?): `slice-021` schließt
  ohne funktionierenden Zugriffsweg (z. B. wegen Carveout) — dann blockiert
  dieser Slice auf dessen Auflösung.

## 5. Closure-Trigger

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

Ein externer Aufruf bestätigt real eine Position über den Zugriffsweg, die
Vorwärts-Invariante ist real getestet, **und** `make gates` ist grün.

## 6. Risiken und offene Punkte

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Der Lesepfad `cdc.consumer_status` (SQL-View) und der Go-Use-Case-Lesepfad
  könnten bei der Verdrahtung des neuen Zugriffswegs erneut auseinanderlaufen
  (`BEO-PGC/lese-doppelquelle`, aktuell 2×) — falls dieser Slice den
  View-Lesepfad berührt, wäre das ein drittes Auftreten. **Ausgang:** wird
  bei Closure eingetragen.

## 7. Closure-Notiz

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <Guide oder Sensor> <geschärft/ergänzt>: <was genau>
  — liegt in `<AGENTS.md §X | Makefile:<target> | .harness/skills/…>`.
  Auslöser: `BEO-<NNN>` (<slice-NNN>, <slice-MMM>, <slice-KKK> — 3×).
  *(Wurde mit diesem Slice nichts verkörpert — der Normalfall —, entfällt die
  Teil-Zeile `— liegt in …` ersatzlos. Der Eintrag ist dann gezählt, nicht
  verkörpert.)*
- **Beobachtungs-Register (`../observations/`):** <`BEO-<KUERZEL>/<slug>/` neu angelegt, Beleg `evidence/slice-NNN.md` | `evidence/slice-NNN.md` in `BEO-<KUERZEL>/<slug>/` ergaenzt — Zaehler steht damit bei <N>x | keine Beobachtung angefallen>
- **Folge-Slices:** <slice-NNN (<Titel>) — ist eine Datei in `open/`>
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>
- **Drei Paarungen:** <nur im Repo ohne Wellen-Betrieb — Anker · Folge-Slice · Register, Ergebnis>

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
durchgegangen — `BEO-PGC/lese-doppelquelle` (2×) betrifft denselben
Consumer-Lesepfad, sofern dieser Slice ihn berührt (siehe §6-Risiko);
`BEO-PGC/rollen-verdrahtung` (1×) wie bereits bei `slice-021` notiert, kein
neuer Aspekt.

**Modus-Begründungsblock — Umfang.** Alle berührten Sub-Areas GF (neue
Fähigkeit auf bestehender Greenfield-Codebasis, keine Inventur-Diskrepanz
zu erwarten).
