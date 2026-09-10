# Slice slice-003: Capture-Persist-Pfad — Ports, Service, Persist-before-ACK

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-1.

**Bezug:** [`LH-QA-REL-001`](../../../../spec/lastenheft.md),
[`LH-QA-REL-002`](../../../../spec/lastenheft.md),
[`LH-FA-CAP-006`](../../../../spec/lastenheft.md),
[`LH-FA-CAP-007`](../../../../spec/lastenheft.md),
[`ADR-0011`](../../../../docs/plan/adr/README.md),
[`ADR-0012`](../../../../docs/plan/adr/README.md),
[`ADR-0007`](../../../../docs/plan/adr/README.md),
[`ADR-0009`](../../../../docs/plan/adr/README.md),
[`ADR-0027`](../../../../docs/plan/adr/README.md),
[`ADR-0028`](../../../../docs/plan/adr/README.md)

**Berührte Spec-Stellen:**
[`LH-QA-REL-001.a`](../../../../spec/pflichtenheft.md),
[`LH-FA-CAP-006.a`](../../../../spec/pflichtenheft.md),
[`SPEC-001`](../../../../spec/pflichtenheft.md),
[`ARC-002`](../../../../spec/architecture.md),
[`ARC-003`](../../../../spec/architecture.md),
[`ARC-004`](../../../../spec/architecture.md)

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

**Ziel:** Der Capture-Pfad im Application-Layer: `CaptureInboundPort`,
`CaptureService` und der `ChangeStorePort` (Interface), mit der
Persist-before-ACK-Ordnung aus [`LH-QA-REL-001.a`](../../../../spec/pflichtenheft.md)
— getestet gegen Fake Ports (Testpyramide,
[`ADR-0030`](../../../../docs/plan/adr/README.md)).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Reale PostgreSQL-Adapter (Store, ACK, Stream) — es wäre ein anderer
  Vorgang: Treiber-Integration mit realem PostgreSQL
  ([`LH-QA-POR-003`](../../../../spec/lastenheft.md)/[`SPEC-011`](../../../../spec/pflichtenheft.md)-Umfang); die Welle 2
  bündelt sie, sobald sie eröffnet wird (noch keine Lifecycle-Datei — die
  Slices von Welle 2 entstehen bei ihrer Eröffnung, Modul 6). Die
  Fake-Port-Tests dieser Welle tragen die Ordnungs-Logik schon.
- Große Transaktionen/Spooling — [`ADR-0021`](../../../../docs/plan/adr/README.md)
  ist Proposed; der TransactionBufferPort wird mit der realen Stream-
  Integration entschieden, nicht hier.
- Consumer-Verwaltung und Retention — andere Use Cases
  ([`ADR-0028`](../../../../docs/plan/adr/README.md)); dieser Slice hält den
  Capture-Pfad klein.

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

- [x] `CaptureInboundPort` und `CaptureService` tragen die
      Persist-before-ACK-Ordnung: ACK nur nach dauerhafter Persistenz —
      Test referenziert zu [`LH-QA-REL-001.a`](../../../../spec/pflichtenheft.md).
- [x] Application-Tests mit Fake Ports belegen: kein ACK bei
      Persistenzfehler; Crash zwischen Persistenz und ACK erzeugt höchstens
      erneute Verarbeitung — Teil-Beleg zu [`LH-QA-REL-002`](../../../../spec/lastenheft.md)
      und [`LH-FA-CAP-006.a`](../../../../spec/pflichtenheft.md) (Rollback-Case in
      [`LH-FA-CAP-007`](../../../../spec/lastenheft.md)-Semantik).
- [x] Application-Testsuite grün (`go test ./internal/application/...`).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update, falls ein öffentlicher Vertrag berührt ist — hier:
      keine Schnittstelle berührt, dann trägt der Bericht die begründete
      Aussage „kein öffentlicher Vertrag berührt".
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. (*§7*)
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
| `internal/application/port/inbound/capture.go` | neu | `CaptureInboundPort` ([`ADR-0028`](../../../../docs/plan/adr/README.md)) |
| `internal/application/port/outbound/changestore.go` | neu | `ChangeStorePort`-Interface ([`ADR-0009`](../../../../docs/plan/adr/README.md)) |
| `internal/application/port/outbound/replicationack.go` | neu | `ReplicationAckPort`-Interface ([`ADR-0007`](../../../../docs/plan/adr/README.md)) — der ACK ist der zweite Port des Capture-Services |
| `internal/application/port/outbound/changestore.go` (Kommentar) | update | *Grenz-Träger aus slice-002 (Lerneintrag):* [`ADR-0029`](../../../../docs/plan/adr/README.md)-Regel 1 (Persist-before-ACK) hat in `ARC-001` kein Subjekt — die Grenze wird hier als Port-Kontrakt-Zeile benannt (ACK nur über nach Persistenz gefragte Positionen) |
| `internal/application/usecase/capture/*.go` | neu | `CaptureService` + Command/Result über zwei Ports ([`ADR-0027`](../../../../docs/plan/adr/README.md), [`ADR-0039`](../../../../docs/plan/adr/README.md)) |
| `internal/application/usecase/capture/*_test.go` | neu | Fake-Ports (Store, ACK): Ordnung, Fehlermodi, Idempotenz |
| `internal/application/port/inbound/capture.go` (Transport-Typen) | update | *Plan-Nachzug (Review F-2, Konflikt-Sequenz Klasse A 3×):* `CaptureCommand`/`CaptureResult` werden am Inbound-Port definiert, der Use Case führt sie als Typ-Aliase — beides trägt die [`ADR-0039`](../../../../docs/plan/adr/README.md)-Hälften (Fitness + Use-Case-Lesart); Schärfung als [`ADR-0042`](../../../../docs/plan/adr/README.md) (Architect-Verdikt der Sequenz) |
| `internal/application/port/*.go` (Signaturen) | update | *Plan-Nachzug (Review F-2):* `context.Context` in allen Port-Signaturen; `ErrMissingTransaction` als Port-Fehler ([`SPEC-008`](../../../../spec/pflichtenheft.md)-Familie) |
| `internal/application/usecase/capture/service.go` (leere Transaktion) | update | *Plan-Nachzug (Review F-4):* Verhalten der leeren committed Transaktion (ablehnen oder Grenze am Port-Kontrakt) + Test — Fix-Runde des Implementers |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): slice-002 ist in `done/` (Domänenmodelle
existieren); kein anderes Slice in `in-progress/`.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wächst die
  Ordnungs-Logik um Consumer- oder Retention-Konzepte (übersteigt die drei
  Liefer-Punkte) → Consumer/Retention in eigene Slices.
- `in-progress` → `open` (blockiert — Carveout?): die Fake-Port-Tests
  offenbaren, dass die Ordnungs-Logik reale Treiber-Semantik voraussetzt
  (z. B. LSN-Streaming-Grenzen) → zurück zur Planung mit
  [`ADR-0021`](../../../../docs/plan/adr/README.md)-Entscheidung.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD (§2) vollständig abgehakt · Application-Testsuite grün · Review-Report
unter `docs/reviews/` · Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- Die Persist-before-ACK-Ordnung ist im Fake-Test stärker als im realen
  Treiber (Fak verhindern, was `pgoutput`-Streaming möglich macht) —
  **Ausgang:** weiter offen, solange kein realer Adapter existiert; trägt
  ins Register, falls er 3× auffällt.
- Idempotenz des Persistierens (Wiederholung nach Crash) braucht eine
  Deduplizierungsbasis ([`SPEC-002`](../../../../spec/pflichtenheft.md): change_id) —
  **Ausgang:** weiter offen; der Fake-Store belegt die Ordnung
  (ACK nach Persistenz), aber nicht die reale Idempotenz-Eigenschaft des
  Stores — die ist [`ADR-0011`](../../../../docs/plan/adr/README.md)-Pflicht
  und bekommt ihren vollen Beleg erst mit dem realen Adapter (Welle 2).

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

- **Was hat funktioniert:** die Ordnungs-Invariante liegt an **einer**
  Orchestrierungs-Stelle (`service.Capture` — persist → ack, kein zweiter
  ACK-Pfad); die Grenz-Träger-Zeilen stehen am Port-Kontrakt (Regel 1,
  Idempotenz), wie der slice-002-Lerneintrag sie vorsah; alle drei
  Ordnungs-Zusagen trugen Mutations-Proben (rot gesehen, Fix-Zug
  `9bfcd7a`, Verifier-Replikation).
- **Was ging anders als geplant:** die Transport-Typen mussten am
  Inbound-Port definiert werden (Plan §3 legte den Ort nicht fest) — der
  Konflikt (beide [`ADR-0039`](../../../../docs/plan/adr/README.md)-Hälften) ging in die **Konflikt-Sequenz**
  (drittes Auftreten der Klasse) und endete in [`ADR-0042`](../../../../docs/plan/adr/README.md)
  (Supersedes [`ADR-0039`](../../../../docs/plan/adr/README.md)); die
  Klasse „Struktur-IDs in Commit-Messagen" erreichte ebenfalls 3× und ist
  als Regel in `.claude/commands/implement-slice.md` verkörpert
  (`· seit slice-003`). Der Review-F-7-Verlauf (stale Zähler im Register)
  war ein §3.7-Defekt, vor der Welle-Closure berichtigt.
- **Steering-Loop-Eintrag:** *Klasse A — Konflikt-Sequenz ausgegangen in
  [`ADR-0042`](../../../../docs/plan/adr/README.md) (Transport-Typen am
  Port; der bestehende Code ist die Ziel-Form) · Klasse B: Commit-
  Message-Regel verkörpert in `.claude/commands/implement-slice.md`
  (Struktur-IDs aus Messagen) · seit slice-003.*
- **Beobachtungs-Register (`../observations/`):** `evidence/slice-003.md`
  in `BEO-PGC/a-check-null-abdeckung/` ergänzt — **Zähler 3×**; der
  Ausgang wird im Lese-Schritt der Welle-1-Closure zugewiesen (Regel-
  Material steht bereit: [`ADR-0042`](../../../../docs/plan/adr/README.md)
  trägt die Alias-Pflicht, a-check bleibt die Maschinenform).
- **Folge-Slices:** keiner — welle-1 schließt mit diesem Slice; Welle 2
  (reale PostgreSQL-Integration) wird bei ihrer Eröffnung geschnitten.
- **Risiken aus §6:** Risiko (a) Fake-Ordnung stärker als realer Treiber →
  **weiter offen** (Beleg erst mit dem realen Adapter, Welle 2); Risiko
  (b) reale Idempotenz des Stores → **weiter offen** ([`ADR-0011`](../../../../docs/plan/adr/README.md)-Pflicht,
  Kontrakt-Zeile getragen; Beleg mit dem realen Adapter).
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

**Vorgelagert — Sub-Area-Wahl prüfen:** berührte Sub-Area: Application-Layer
(`internal/application/`; [`ARC-002`](../../../../spec/architecture.md)–[`ARC-004`](../../../../spec/architecture.md)).
Achsen: (1) Konventionen-Dichte — Ports/Use-Cases in [`ADR-0028`](../../../../docs/plan/adr/README.md)/[`ADR-0039`](../../../../docs/plan/adr/README.md),
Ordnung in [`LH-QA-REL-001.a`](../../../../spec/pflichtenheft.md) verankert;
(2) Phase-Reife — Phase 3; (3) Evidenz-Risiko niedrig (GF). Schwelle ≥ 2
erfüllt.

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
