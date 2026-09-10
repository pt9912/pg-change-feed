# Slice slice-006: Docker-Compose-Umgebung und MVP-Integrationstest

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-2.

**Bezug:** [`LH-QA-POR-003`](../../../../spec/lastenheft.md), [`LH-QA-OPS-001`](../../../../spec/lastenheft.md), [`LH-QA-OPS-002`](../../../../spec/lastenheft.md)

**Berührte Spec-Stellen:** [`SPEC-011`](../../../../spec/pflichtenheft.md)/[`SPEC-015`](../../../../spec/pflichtenheft.md) (Deployment-Formen, Umfang), [`LH-QA-POR-003`](../../../../spec/lastenheft.md) (reproduzierbare Umgebung), [`ARC-007`](../../../../spec/architecture.md) — *Korrektur (Review F-11): die Verweise `PH-DEP-002`/`PH-TST-001` adressieren Kennungen, die kein Stratum trägt (Vorbestand aus dem zurückgezogenen Pflichtenheft-Entwurf); die echten Anker sind die Deployment-Verträge und die Testpyramide-Regeln (`.harness/skills`-Muster, [`ADR-0030`](../../../../docs/plan/adr/0030-testpyramide.md)-Bahn).*
Der Verweis zeigt **aufwärts**: Die Spec nennt diesen Slice nie
(Baseline-Regelwerk `grundlagen-referenz-richtung.md`
§Referenz-Richtung (SDP), `grundlagen-source-precedence.md` §ID-Schema als Klammer).

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

**Ziel:** Die reproduzierbare Docker-Compose-Umgebung (PostgreSQL + PG Change Feed + Test-Consumer) und der automatisierte MVP-Integrationstest: PostgreSQL starten → CDC aktivieren → INSERT/UPDATE/DELETE → Changes lesen → Reihenfolge und Inhalt prüfen.

*Plan-Nachzug ([`ADR-0043`](../../../../docs/plan/adr/README.md)):* der **Schema-Rollout in die Compose-Test-DB läuft über d-migrate** (`make schema-rollout` mit Pflicht-Report und Rollback-Artefakt) **vor jedem E2E-Lauf** — das ist der Erstversatz der ADR und löst die handgeschriebene DDL + `ApplySchema`-Grenze aus slice-004 ab.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Produktions-Deployment (HA/Kubernetes) — Out-of-Scope des MVP.
- Performance-Benchmarks mit Laststufen — die
  [`LH-QA-PER-002`](../../../../spec/lastenheft.md)-Messmethode braucht das
  Benchmark-Design; erst nach dem grünen E2E.
- Exportadapter (Kafka/NATS/…) — keine Anforderung dieses Lastenhefts.
- d-migrate-Nacharbeit (CHECK `chk_change_operation` lebt als psql-Schritt
  statt im Rollout — [`ADR-0043`](../../../../docs/plan/adr/README.md)-Ausweichform
  am Blocker `raw-sql-text-drift`, d-migrate 1.2.0) — **Ausgang: weiter
  offen** → Retirement mit dem d-migrate-Fix-Release (Pin-Hebung +
  `chk_change_operation` direkt ins YAML + Rückbau des psql-Schritts);
  Träger: der Beobachtungs-Register-Eintrag
  `BEO-PGC/d-migrate-nacharbeit` (angelegt mit dieser Closure), Beleg
  `evidence/slice-006.md`.

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

- [ ] Die Compose-Umgebung startet PostgreSQL + PG Change Feed
      reproduzierbar aus dokumentierten Schritten — Teil-Beleg zu
      [`LH-QA-POR-003`](../../../../spec/lastenheft.md).
- [ ] Der MVP-Integrationstest läuft automatisiert grün: aktivieren →
      INSERT/UPDATE/DELETE → lesen → Reihenfolge/Inhalt (Abschnitt 1,
      MVP-Schnitt) — Teil-Beleg zu [`LH-FA-CAP-001`](../../../../spec/lastenheft.md)…003.
- [ ] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update bzw. begründete Aussage „kein öffentlicher Vertrag
      berührt" — fällig bei der Implementierung.
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
| `compose.yaml` | neu | Lauf- und CI-Vertrag je Modul 14 (Devcontainer = Komfort, nicht hier) |
| `test/integration/*` | neu | MVP-Integrationstest gegen die Compose-Umgebung |
| `Makefile` (schema-rollout-Verkabelung) | update | *Plan-Nachzug ([`ADR-0043`](../../../../docs/plan/adr/README.md)):* Schema-Rollout in die Compose-Test-DB vor jedem E2E-Lauf (`make schema-rollout` mit Pflicht-Report und Rollback-Artefakt) — Erstversatz der ADR, löst die slice-004-Loader-Grenze ab |
| `.a-check.yml` (Composition-Root-Glob) | update | *Plan-Nachzug (F-3, jetzt tatsächlich getragen):* `test/integration/**` in `composition_root` — der Verdrahtungstest-Layer als eigene a-check-Klasse (Architect-Fix, a-check 0 Befunde) |
| `harness/README.md` (Werkzeuge-Zeilen) | update | *Plan-Nachzug (F-3):* `test-integration`-Target als Werkzeuge-Zeile (gepinnte Digests, Compose-Kette) — Archiver-Fixes getragen |
| `tools/schema/schema.yaml` | neu | *Plan-Nachzug ([`ADR-0043`](../../../../docs/plan/adr/README.md)-Erstlieferung):* das neutrale Schema-YAML, überführt aus `internal/adapters/driven/postgresstorage/schema.sql` (Stand slice-004) — die Quelle, ohne die `schema-rollout` keinen Input hat; Platzierung (`tools/schema/` vs. co-loziert am Adapter) ist die offene Platzierungs-Frage aus [`ADR-0043`](../../../../docs/plan/adr/README.md) |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): slice-005 liegt in `done/` (Stream + ACK real); kein anderes Slice in
`in-progress/` (WIP-Limit 1).

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

- §6-Risiko (a): die reale Ordnung weicht vom Fake ab — **Ausgang:**
  bewertet am Integrationstest (Träger des slice-003-Ausgangs).
- §6-Risiko (b): Row-Images unvollständig (REPLICA IDENTITY) —
  **Ausgang:** bewertet am Integrationstest.

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

**Vorgelagert — offene Beobachtungen sichten:** Register trägt nur den
  verkörperten Eintrag `BEO-PGC/a-check-null-abdeckung` (verkörpert ·
  seit welle-1) — **keine offenen Treffer**; notiert.

**Modus-Begründungsblock — Umfang.** Pflicht, sobald mindestens eine berührte
Sub-Area BF oder Hybrid ist — einer pro Sub-Area. Bei reinem GF genügt der
Hinweis *"alle berührten Sub-Areas GF"*; bei reinem Refactor ohne neue
Sub-Area-Berührung entfällt **er** — nicht der Abschnitt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Deklaration
`harness/conventions.md`, Default-Zeile) — kein BF/Hybrid, kein Block.
