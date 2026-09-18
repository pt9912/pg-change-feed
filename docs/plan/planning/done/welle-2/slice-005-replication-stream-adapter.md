# Slice slice-005: Replication-Stream-Adapter (pgoutput, real)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.

**Welle:** welle-2.

**Bezug:** [`LH-FA-CAP-001`](../../../../spec/lastenheft.md)…003, [`LH-QA-REL-001`](../../../../spec/lastenheft.md), [`ADR-0006`](../../../../docs/plan/adr/README.md), [`ADR-0008`](../../../../docs/plan/adr/README.md)

**Berührte Spec-Stellen:** [`LH-FA-CAP-006.a`](../../../../spec/pflichtenheft.md), [`SPEC-010`](../../../../spec/pflichtenheft.md), [`ARC-005`](../../../../spec/architecture.md), [`ARC-008`](../../../../spec/architecture.md)
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

**Ziel:** Der reale `Replication-Stream-Adapter`: `pgoutput` dekodieren (`receive/`, `decode/`, `mapper/`), in Aufrufe des `CaptureInboundPort` übersetzen, den `ReplicationAckPort` real anbinden — Persist-before-ACK am realen Treiber; die Aktivierung inklusive Publication/Slot ([`LH-FA-CFG-001.a`](../../../../spec/pflichtenheft.md)).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- Große Transaktionen/Spooling — [`ADR-0021`](../../../../docs/plan/adr/README.md) ist
  Proposed; der TransactionBufferPort wird entschieden, wenn der
  Integrationstest ein RAM-Problem belegt (nicht vorab).
- HTTP-/gRPC-API — [`ADR-0020`](../../../../docs/plan/adr/README.md) bleibt optional ohne
  beobachtbaren Bedarf.
- Consumer-Verwaltung — nicht MVP.
- d-migrate-Rollout — nicht nötig in diesem Slice: die Stream-Adapter-
  Tests brauchen ein Fixture-Schema (Publikation über eine Fixture-Tabelle),
  keinen CDC-Schema-Rollout ([`ADR-0043`](../../../../docs/plan/adr/README.md):
  Erstversatz ist slice-006, Compose-Umgebung).

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

- [x] Der Stream-Adapter übersetzt `pgoutput` in Aufrufe des
      `CaptureInboundPort` (INSERT/UPDATE/DELETE) — Teil-Beleg zu
      [`LH-FA-CAP-001`](../../../../spec/lastenheft.md)…003. *Beleg:
      verify-slice-005.md Item 1.*
- [x] ACK real: nur Positionen nach dauerhafter Persistenz —
      [`LH-QA-REL-001.a`](../../../../spec/pflichtenheft.md) am realen
      Treiber. *Beleg: `TestRealPersistBeforeAck` + Keepalive-Regel mit
      roter Probe (verify-slice-005.md Item 2).*
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      *Beleg: review-slice-005.md committet (`22e9b93`).*
- [x] Doku-Update — getragen: die Werkzeuge-Zeile `test-replication`
      (`harness/README.md`) ist über den Verifier bestätigt
      (verify-slice-005.md Item 5).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. (*§7*)
- [x] Reconciliation-Register — **entfällt**: Repos ohne
      Brownfield-Bootstrap haben die Datei nicht.
- [x] Beobachtungs-Register fortgeschrieben — **kein neues Auftreten**
      (a-check 0 Befunde am HEAD; die verkörperte Beobachtung
      `BEO-PGC/a-check-null-abdeckung` ist unverändert).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (siehe §7).
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
| `internal/adapters/driving/replication/{receive,decode,mapper}/*.go` | neu | Stream-Adapter je [`ADR-0006`](../../../../docs/plan/adr/README.md)/[`ADR-0039`](../../../../docs/plan/adr/README.md) |
| `internal/adapters/driven/postgresack/*.go` | neu | Realer ACK-Adapter ([`ADR-0007`](../../../../docs/plan/adr/README.md)) |
| `internal/application/port/outbound/replicationack.go` (Kontrakt) | update | *Plan-Nachzug (V-2):* Port-Sentinel `ErrReplication` (Klasse `replication`, [`ADR-0023`](../../../../docs/plan/adr/README.md)) als Kontrakt-Zeile am Port — Träger-Muster wie `ErrStorage` (slice-004) |
| `internal/adapters/driving/replication/` (Keepalive-Regel) | update | *Plan-Nachzug (Review F-2, Implementer-Entscheidung):* Keepalive-Antworten melden ausschließlich die bestätigte Position — never über den Empfangsstand hinaus ([`LH-QA-REL-001.a`](../../../../spec/pflichtenheft.md), Datenverlust-Fenster); Test mit roter Probe am realen Pfad |
| `internal/adapters/driving/replication/receive/receive.go` (BindCapture) | update | *Plan-Nachzug (Review F-2):* Verbindungsaufbau getrennt von Port-Verdrahtung — der ACK-Adapter braucht die Verbindung erst nach deren Aufbau ([`ADR-0007`](../../../../docs/plan/adr/README.md) Option C) |
| `internal/bootstrap/replication_stream_test.go` | neu | *Plan-Nachzug (Review F-2, a-check-Befunde):* Verdrahtungstests im Composition-Root-Layer ([`ADR-0026`](../../../../docs/plan/adr/README.md)) — Adapter-Tests importieren keine fremden Adapter/Use Cases |
| `Makefile` (test-replication) | update | *Plan-Nachzug:* Testcontainer-Target gegen reale PostgreSQL (`wal_level=logical`, gepinnte Digests, tabellenscopierte Publications) + `tools/harness/run-replication-tests.sh` |
| `github.com/jackc/pglogrepl` | neu (deps) | Replication-Protokoll-Handshake — [`ADR-0032`](../../../../docs/plan/adr/README.md): Infrastrukturdetail; pgx/v5 trägt den Fall nicht (grep-Beleg im Modul-Cache) |

## 4. Trigger

<!-- BEDIENHINWEIS: Beispiele — "Wenn Welle X done." / "Wenn Carveout CO-NN
aufgeloest." -->

Regeln dieser Sektion: Baseline-Regelwerk `modul-05-planning-harness.md`
§Trigger je Lifecycle-Übergang und WIP-Limit.

**Start** (`next` → `in-progress`): slice-004 liegt in `done/` (Store persistiert real); kein anderes Slice
in `in-progress/` (WIP-Limit 1).

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

- (a) Fake-Ordnung stärker als realer Treiber (slice-003-Ausgang) —
  **Ausgang:** bewertet mit diesem Slice am realen Treiber (Träger des
  slice-003-Ausgangs).
- (b) Reale Idempotenz — **Ausgang:** getragen durch den
  [`ADR-0011`](../../../../docs/plan/adr/README.md)-Kontrakt in slice-004 (Deduplizierung über
  Transaktions-ID), belegt im Integrationstest (slice-006).

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

- **Was hat funktioniert:** die Keepalive-Regel („nur die bestätigte
  Position") trägt einen Test **am realen Pfad** mit roter Probe — genau
  das Datenverlust-Fenster, das der Review als nicht gedeckt meldete; der
  Verdrahtungstest zog in den Composition-Root-Layer (a-check-Befunde
  getragen statt still gefixt); der Restart-Zweig (`ensureSlot`) trägt
  seinen Test.
- **Was ging anders als geplant:** der Review-F-7-Schiedsspruch
  widerlegte die `96c47af`-Beleg-Formel (Digest wechselt bei deps-
  Änderungen ohne Binary-Import) — die Korrektur ging als
  [`ADR-0044`](../../../../docs/plan/adr/README.md) über den Architect,
  nicht als Implementer-Nachzug; der doppelte BEGIN verschwieg still
  (F-8) und trägt jetzt einen Sentinel. Die Test-Flakiness (shared
  Container, FOR ALL TABLES) wurde mit tabellenscopierten Publications
  gelöst.
- **Steering-Loop-Eintrag:** *F-6-Klasse „Datei-Ende-Zeilenumbrüche
  fehlen" erreicht 3× (slice-001 F-6, slice-004 F-5, slice-005 F-6) —
  Ausgang folgt im Lese-Schritt der Welle-2-Closure; die Fix-Runde
  schloss alle berührten Dateien.* *(Verkörperung steht aus; gezählt,
  nicht verkörpert.)*
- **Beobachtungs-Register (`../observations/`):** **keine Beobachtung
  angefallen** — die verkörperte Beobachtung
  `BEO-PGC/a-check-null-abdeckung` ist unverändert (a-check 0 Befunde am
  vollen Baum).
- **Folge-Slices:** slice-006 (Compose + MVP-Integrationstest) — ist
  eine Datei in `open/`.
- **Risiken aus §6:** Risiko (a) Fake-Ordnung stärker als realer Treiber
  → **entfallen** (Ordnungszusagen am realen Pfad belegt);
  Risiko (b) reale Idempotenz → **weiter offen** ([`ADR-0011`](../../../../docs/plan/adr/README.md)-Kontrakt
  trägt; Beleg im Integrationstest, slice-006).
- **Drei Paarungen:** im Wellen-Betrieb an die Welle-2-Closure delegiert.

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
