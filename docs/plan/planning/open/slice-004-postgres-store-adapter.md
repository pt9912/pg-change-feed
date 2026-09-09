# Slice slice-004: PostgreSQL-ChangeStore-Adapter

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-2.

**Bezug:** [`LH-QA-REL-001`](../../../../spec/lastenheft.md), [`LH-QA-REL-002`](../../../../spec/lastenheft.md), [`ADR-0010`](../../../../docs/plan/adr/README.md), [`ADR-0032`](../../../../docs/plan/adr/README.md)

**Berührte Spec-Stellen:** [`SPEC-001`](../../../../spec/pflichtenheft.md), [`SPEC-002`](../../../../spec/pflichtenheft.md), [`ARC-009`](../../../../spec/architecture.md), [`ARC-006`](../../../../spec/architecture.md)
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

**Verantwortlich:** — bis zur Priorisierung.
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

**Ziel:** Der reale `PostgresChangeStoreAdapter`: Persist mit Deduplizierung über die interne Transaktions-ID (Idempotenz-Kontrakt, [`ADR-0011`](../../../../docs/plan/adr/README.md)), deterministisches Lesen, Position-Verwaltung — gegen reale PostgreSQL-Tabellen je [`SPEC-001`](../../../../spec/pflichtenheft.md)/[`SPEC-002`](../../../../spec/pflichtenheft.md), mit Adapter-Tests gegen eine reale Testcontainer-PostgreSQL.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Reale Replication-Stream-Integration — ein Folge-Slice übernimmt es
  (slice-005); der Store wird hier unabhängig von der Stream-Seite real.
- Consumer-Use-Cases und Retention-Adapter — es wäre ein anderer Vorgang
  (nicht MVP, [`ADR-0028`](../../../../docs/plan/adr/README.md)-Rest).
- Schreibpfad-Performance-Tuning — Messungen folgen mit dem
  MVP-Integrationstest (slice-006).

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

- [ ] `PostgresChangeStoreAdapter` persistiert committed Transaktionen
      deduplizierbar (Idempotenz-Kontrakt [`ADR-0011`](../../../../docs/plan/adr/README.md):
      dieselbe Transaktion erneut → höchstens erneute Verarbeitung) —
      Teil-Beleg zu [`LH-FA-RET-001`](../../../../spec/lastenheft.md).
- [ ] Deterministisches Lesen über [`SPEC-001`](../../../../spec/pflichtenheft.md)-Tabellen
      (Bereich/Limit/Filter) — Teil-Beleg zu
      [`LH-FA-REA-001`](../../../../spec/lastenheft.md)…006.
- [ ] `make gates` grün.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update für <Schnittstelle X> falls öffentlicher Vertrag berührt.
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben, **falls dieser Slice einen Inventur-Fund auflöst** — Zeile mit Datum und auflösendem Artefakt nach *Aufgelöste Einträge* verschoben. Repos ohne Brownfield-Bootstrap haben die Datei nicht; dann entfällt das Item.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues Verzeichnis `BEO-<KUERZEL>/<slug>/` oder eine weitere Datei in dessen `evidence/`; **kein Zaehler wird gesetzt**, er folgt aus den Dateien. Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter offen).
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
| `internal/adapters/driven/postgresstorage/*.go` | neu | row-Typen, `queries/`, `mapper/` je [`ADR-0039`](../../../../docs/plan/adr/README.md); [`SPEC-001`](../../../../spec/pflichtenheft.md)-Schema-DDL |
| `internal/adapters/driven/postgresstorage/*_test.go` | neu | Adapter-Tests gegen reale PostgreSQL (Testcontainer im gepinnten Docker-Netz) |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): slice-003 liegt in `done/` (Ports und Service existieren); kein anderes
Slice in `in-progress/` (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): wächst der
  Adapter-Scope über die drei Liefer-Punkte hinaus (z. B. die Leseseite
  verdient einen eigenen Slice) → aufteilen.
- `in-progress` → `open` (blockiert — Carveout?): Treiber-/Testcontainer-
  Faktoren (PostgreSQL-Image nicht erreichbar) blockieren den realen Lauf.

## 5. Closure-Trigger

<!-- BEDIENHINWEIS: z.B. "DoD vollstaendig + PR gemerged + Closure-Notiz
geschrieben." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Closure- und Lerneintrag-Regeln — zwei beobachtbare Kriterien **und** ein
Lerneintrag; ohne ihn ist der Slice nur abgelegt.

DoD (§2) vollständig abgehakt · Tests im Toolchain-Container grün ·
Review-Report unter `docs/reviews/` · Closure-Notiz mit Lerneintrag.

## 6. Risiken und offene Punkte

<!-- BEDIENHINWEIS: Was koennte schief gehen? Welche Carveouts entstehen
ggf.? Die drei Ausgaenge stehen als Form in der Zeile darunter. -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Offene Risiken werden bei Closure aufgelöst — **jedes** Risiko bekommt genau
**einen** Ausgang, und kein Slice geht nach `done/`, während eines ohne Ausgang
dasteht.

- (a) Die reale Idempotenz-Deduplizierung ist stärker als der Fake —
  **Ausgang:** weiter offen; der Adapter-Test belegt sie, der
  slice-003-Ausgang trägt den Termin (Welle 2).
- (b) Row-Images liefern nicht alle Werte je Operationstyp (REPLICA
  IDENTITY) — **Ausgang:** weiter offen; Bewertung mit dem
  Integrationstest (slice-006).

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

**Vorgelagert — Sub-Area-Wahl prüfen:** berührte Sub-Area:
PostgreSQL-Treiber-Integration (`internal/adapters/driven/**` bzw.
`driving/replication/**`, Compose-Umgebung); Achsen: (1) Konventionen-
Dichte — Adapter-Form in [`ADR-0039`](../../../../docs/plan/adr/README.md), Schichten-Constraints in
[`ADR-0002`](../../../../docs/plan/adr/README.md), Maschinenform in `.a-check.yml` (Edges aktiv);
(2) Phase-Reife — Phase 4 (Spec + ADRs + Domain/Ports committet, Adapter
entstehen); (3) Evidenz-Risiko niedrig (GF). Schwelle ≥ 2 erfüllt.

**Vorgelagert — offene Beobachtungen sichten:** das Register trägt nur
  den verkörperten Eintrag `BEO-PGC/a-check-null-abdeckung` (Ausgang
  `verkörpert · seit welle-1`) — **keine offenen Treffer**; notiert.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Deklaration
`harness/conventions.md`, Default-Zeile) — kein BF/Hybrid, kein Block.
