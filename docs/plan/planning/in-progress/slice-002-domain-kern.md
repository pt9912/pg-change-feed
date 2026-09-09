# Slice slice-002: Domänenkern — Modelle, Invarianten, ClockPort

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-1.

**Bezug:** [`LH-FA-DAT-001`](../../../../spec/lastenheft.md),
[`LH-FA-DAT-004`](../../../../spec/lastenheft.md),
[`LH-FA-CAP-004`](../../../../spec/lastenheft.md),
[`LH-FA-CAP-005`](../../../../spec/lastenheft.md),
[`ADR-0004`](../../../../docs/plan/adr/README.md),
[`ADR-0005`](../../../../docs/plan/adr/README.md),
[`ADR-0029`](../../../../docs/plan/adr/README.md),
[`ADR-0039`](../../../../docs/plan/adr/README.md),
[`ADR-0040`](../../../../docs/plan/adr/README.md)

**Berührte Spec-Stellen:** [`SPEC-002`](../../../../spec/pflichtenheft.md),
[`SPEC-003`](../../../../spec/pflichtenheft.md),
[`ARC-001`](../../../../spec/architecture.md)

**Verantwortlich:** pt9912 (Implementer-Rolle, Agent-Lauf).

**Autor:** pt9912. **Datum:** 2026-09-09.

---

## 1. Ziel und Abgrenzung

<!-- BEDIENHINWEIS: Ziel = ein Satz, Liefer-Fokus, kein "wir machen
aufraeumen". Abgrenzung = je Punkt eine Begruendung, nicht nur eine Nennung:
ein Ausschluss ohne Grund ist eine Behauptung. Keine Mindestzahl — ein echter
Ausschluss ist besser als vier erfundene. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — Schnitt nach Lieferwert, nicht nach Schichten; jeder Slice
ist einzeln lieferbar. **§1 nennt Ziel und Abgrenzung** (Out-of-Scope-Disziplin
des Lastenhefts, auf den Slice-Plan angewandt); die vier Klassen des
Ausschlusses stehen in **eben diesem Abschnitt** des Baseline-Regelwerks,
zusammen mit der Begründungs-Pflicht je Punkt.

**Ziel:** Die Domänenmodelle je [`ARC-001`](../../../../spec/architecture.md) —
Source, SourceTable, Change, ChangeTransaction, SourcePosition, Consumer,
ConsumerPosition, SchemaVersion, RetentionPolicy — mit Invarianten-
Validierung in Konstruktoren und Domain-Tests; dazu das `ClockPort`-Interface
([`ADR-0040`](../../../../docs/plan/adr/README.md)). `SourcePosition` folgt
[`ADR-0005`](../../../../docs/plan/adr/README.md) (abstrahiert die LSN).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Persistenz und Store-Vertrag — ein Folge-Slice übernimmt es (slice-003
  setzt die Ports; der reale Adapter folgt in einer späteren Welle);
  Domänen-Tests brauchen keinen Store.
- `pgoutput`-Dekodierung und Replication-Stream — es wäre ein anderer
  Vorgang (Treiber-Integration); die Modelle sind technologieunabhängig
  ([`ADR-0032`](../../../../docs/plan/adr/README.md)).
- SQL-/CLI-Adapter — Schicht-Abgrenzung: dieser Slice berührt nur die
  Domain-Schicht und den Port-Typ; Adapter folgen nach den Use Cases.
- Domain-Events (`domain/event/`) — [`ADR-0025`](../../../../docs/plan/adr/README.md)
  erlaubt sie, aber kein Use Case verbraucht sie hier; erst mit dem ersten
  Consumer- oder Lag-Event.

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

<!-- BEDIENHINWEIS: je Zeile ein pruefbares Kriterium. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Ziel-Form: Slice — **≤ 3 Liefer-Punkte**; mehr heißt: der Slice ist zu groß und
gehört zurück zur Zerlegung. Gezählt wird nur, was mit dem Umfang wächst — die
Gate-Läufe und die fünf Closure-Pflichten darunter zählen nicht mit.

- [x] Domänenmodelle je [`ARC-001`](../../../../spec/architecture.md) tragen ihre
      Invarianten; Domain-Tests prüfen [`LH-FA-DAT-001`](../../../../spec/lastenheft.md)
      (eindeutige Identifikation) und [`LH-FA-DAT-004`](../../../../spec/lastenheft.md)
      (sortierbare Position) mit referenzierten Tests.
- [x] Transaktionszusammengehörigkeit und eindeutige Sequenz modelliert —
      Teil-Beleg zu [`LH-FA-CAP-004`](../../../../spec/lastenheft.md) und
      [`LH-FA-CAP-005`](../../../../spec/lastenheft.md).
- [x] `ClockPort`-Interface (Outbound, [`ADR-0040`](../../../../docs/plan/adr/README.md))
      mit Fake Clock in den Tests.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update, falls ein öffentlicher Vertrag berührt ist — hier:
      keine Schnittstelle berührt, dann trägt der Bericht die begründete
      Aussage „kein öffentlicher Vertrag berührt".
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — im Repo **ohne** Wellen-Betrieb hier geprüft, im Repo **mit** Wellen von der nächsten Welle-Closure (auch für Slices ohne Wellen-Zugehörigkeit).

## 3. Plan (vor Code)

<!-- BEDIENHINWEIS: Datei- oder Komponenten-Ebene reicht; der
Implementer-Agent erweitert die Liste in seinem ersten Lauf. -->

Regeln dieser Sektion: Baseline-Regelwerk `grundlagen-bootstrap.md`
§Was ist eine Sub-Area? — diese Liste liefert die **Pfad-Kandidaten** für §8,
nicht die Antwort: Pfad-Berührung ist nicht hinreichend, und eine
Aussagen-Berührung steht hier gar nicht.

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/domain/model/*.go` | neu | Domänenobjekte je [`ARC-001`](../../../../spec/architecture.md); Invarianten in Konstruktoren ([`ADR-0029`](../../../../docs/plan/adr/README.md)) |
| `internal/domain/errors/*.go` | neu | Domänen-Fehler (Invarianten-Verletzungen), [`ADR-0039`](../../../../docs/plan/adr/README.md) |
| `internal/application/port/outbound/clock.go` | neu | `ClockPort` ([`ADR-0040`](../../../../docs/plan/adr/README.md)) |
| `internal/domain/model/consumer.go` | update | *Plan-Nachzug (Review F-4, Endstand nach F-2-Fix `3250232`):* die Quelle einer `ConsumerPosition` trägt allein `Position.SourceID` ([`SPEC-003`](../../../../spec/pflichtenheft.md)) — die erste Bestätigung bindet sie, jede weitere muss sie wiederholen; die zunächst doppelt geführte Quelle (Implementer-Ergänzung 247d188, „Plan-**Ergänzung**" bewertet) ist im F-2-Fix entfernt |
| `internal/domain/model/*_test.go` | neu | Domain-Tests (Testpyramide-Basis, [`ADR-0030`](../../../../docs/plan/adr/README.md)) |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): slice-001 ist in `done/` (go.mod und
Baum existieren); kein anderes Slice in `in-progress/`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): übersteigt der
  Modell-Umfang die drei Liefer-Punkte (z. B. wenn RetentionPolicy eine
  eigene Policy-Engine verdient) → aufteilen.
- `in-progress` → `open` (blockiert — Carveout?): eine Invariante lässt sich
  nicht ohne Quellentscheidung modellieren (z. B. LSN-Semantik) → Carveout
  oder Architekt-Rückfrage.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD (§2) vollständig abgehakt · Domain-Testsuite grün · Review-Report unter
`docs/reviews/` · Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Die Positionsordnung (SourcePosition) braucht Semantik-Entscheidungen
  (Totale Ordnung vs. Transaktions-Kohärenz), die [`ADR-0029`](../../../../docs/plan/adr/README.md)
  nur als Invariante nennt, nicht als Typ-Entwurf — **Ausgang:** offen;
  falls die Entscheidung schieft, Architekt-Rückfrage, sonst folgt sie aus
  [`LH-FA-REA-004.a`](../../../../spec/pflichtenheft.md).
- `ClockPort`-Typ (domänengetragener `TimePoint`) kollidiert mit Go-Konventionen
  (Standard-`time.Time` an der Port-Kante) — **Ausgang:** offen; bewertet bei
  Closure (Interface kann `time.Time` durchreichen, Domain bleibt pur).

## 7. Closure-Notiz

<!-- BEDIENHINWEIS — keine Norm; faellt beim Kopieren weg (README.md
§Verwendung, Schritt 5) und darf deshalb nichts Tragendes halten. Reihenfolge:
diese Sektion vor dem `git mv` nach done/ fuellen — einzige Ausnahme ist das
letzte DoD-Item in §2 (die Paarungen suchen in `done/`, also nach dem `git mv`).
Im Repo ohne Wellen-Betrieb braucht die Closure dadurch drei Commits: Inhalt,
`git mv`, Haekchen — das folgt aus der Hard Rule, es widerspricht ihr nicht. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Das Beobachtungs-Register (vorhandene `BEO-<NNN>` **zitieren** statt neu
formulieren — sonst zählt das Register zwei Namen getrennt) ·
`grundlagen-traceability.md` §Herkunfts-Anker für Steering-Loop-Regeln (das
Feld `liegt in` steht **nur**, wenn mit diesem Slice wirklich etwas verkörpert
wurde; Feld und Zielort auf **einer** Zeile, Sektionsangabe innerhalb der
Backticks).

- **Was hat funktioniert:** die Layer-Globs des a-check matchen mit diesem
  Slice erstmals echten Domain-/Ports-Content (0 Befunde, Kante
  `ports → domain` per `a-check-graph` bestätigt); die Invarianten liegen
  in Konstruktoren statt nur in Kommentaren (nach Review F-1); der
  Verlaufskorrektur-Schnitt (3 Commits, je stand-alone) hielt der
  Stand-alone-Prüfung je Commit stand.
- **Was ging anders als geplant:** der erste Commit-Schnitt kompilierte
  nicht stand-alone und wurde von `9a4d8ad` aus neu geschnitten — dabei
  landete der Erstversuch (`e509c09`) über den Auto-Push der Umgebung auf
  origin und erzeugte eine Divergenz, die per Force-Push mit
  Ours-Auflösung (Vorgabe pt9912) aufgelöst wurde. Lesson: der
  Verlaufskorrektur-Schnitt muss vor jedem Push stehen, nicht nach ihm.
- **Steering-Loop-Eintrag:** geschärfte Regel: *[`ADR-0029`](../../../../docs/plan/adr/README.md)-Invarianten gelten
  als getragen nur, wenn der Typ sie erzwingt; ein Kommentar-Claim ohne
  Typ-Träger ist ein Review-Befund (F-1), und Invarianten ohne Subjekt im
  [`ARC-001`](../../../../spec/architecture.md)-Baum werden als Grenze mit
  Folgetermin benannt (Regel 1: Persist-before-ACK folgt mit dem
  Store-Port, slice-003 §3)* —
  *(gezählt, nicht verkörpert; kein Zielort trägt sie als Artefakt).*
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-002.md`
  in `BEO-PGC/a-check-null-abdeckung/` ergänzt — Zähler steht damit bei
  2× (teils entkräftet: `domain`/`ports` matchen, `app`/`adapters`
  weiter null bis slice-003).
- **Folge-Slices:** slice-003 (Capture-Persist-Pfad) — ist eine Datei in
  `open/`.
- **Risiken aus §6:** Risiko 1 (Positionsordnung als Quell-Tiebreak) →
  **entfallen** (Review bestätigt: CAP-004 fordert keine Totalordnung über
  Quellen, Source-Mismatch-Check vor `Before`); Risiko 2 (ClockPort-Typ) →
  **entfallen** (domänengetragener `TimePoint`, konsistent mit
  [`ADR-0040`](../../../../docs/plan/adr/README.md)-Fitness).
- **Drei Paarungen:** im Wellen-Betrieb an die Welle-1-Closure delegiert.

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

**Vorgelagert — Sub-Area-Wahl prüfen:** berührte Sub-Area: Domain-Kern
(`internal/domain/`; [`ARC-001`](../../../../spec/architecture.md)). Achsen:
(1) Konventionen-Dichte — Paketstruktur in [`ADR-0039`](../../../../docs/plan/adr/README.md),
Invarianten in [`ADR-0029`](../../../../docs/plan/adr/README.md), Modelle in
[`SPEC-002`](../../../../spec/pflichtenheft.md)/[`SPEC-003`](../../../../spec/pflichtenheft.md)
verankert; (2) Phase-Reife — Phase 3 (Spec committet, Code folgt); (3)
Evidenz-Risiko niedrig (GF). Schwelle ≥ 2 erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** Register trägt nur seine
`README.md` — **keine Treffer**; notiert.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

<!-- Block für jede berührte Sub-Area duplizieren. Format identisch
mit dem im Baseline-Regelwerk §Ziel-Form: Sub-Area-Modus-Begründung
abgedruckten Block. -->

*Reiner GF-Hinweis genügt (siehe oben); kein Sub-Area-Block.*
