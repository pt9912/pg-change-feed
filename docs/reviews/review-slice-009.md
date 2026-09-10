# Review-Report: slice-009 — 2026-09-10

**Review-Art:** Code — geprüft gegen Slice-Plan + ADRs (Maintainability).

**Gegenstand:** Implementer-Commits von slice-009, Range `60cfd3e..HEAD`
(be191af = `next` → `in-progress`; 7 Implementer-Commits `becafdf…ce2900e`).

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-09-09 (vier
repo-spezifischen HIGH-Regeln, drei MEDIUM-Klassen) ·
**Modell:** Claude Code (GLM) · **Datum:** 2026-09-10

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- Slice-Plan slice-009 (Consumer-Verwaltung, `LH-FA-CON-001…006`, Welle 3)
- `spec/lastenheft.md` §[`LH-FA-CON-001`](../../spec/lastenheft.md)…006 · `spec/pflichtenheft.md` §2
  (`SPEC-001`, `SPEC-008`-Klassen-Tabelle)
- [`ADR-0013`](../plan/adr) (Consumer-Domänenkonzept), [`ADR-0028`](../plan/adr) (Inbound Use Cases),
  [`ADR-0029`](../plan/adr) (Domain-Invarianten, Regel 2), [`ADR-0034`](../plan/adr) (Fähigkeits-Ports),
  [`ADR-0042`](../plan/adr) (Transport-Typen am Port), [`ADR-0043`](../plan/adr) (d-migrate-Kette)
- `AGENTS.md` §3 Hard Rules · `harness/conventions.md` (MR-000/MR-001)
- Register-Sichtung: `BEO-PGC/plan-nachzug` (1×), `BEO-PGC/adapter-fehler-ausgang`
  (1×), `BEO-PGC/walsender-wirksamkeit` (1×), `BEO-PGC/a-check-null-abdeckung`
  (verkörpert), `BEO-PGC/d-migrate-nacharbeit` (1×)
- vorherige Reports: review-slice-005 (F-6/F-9/F-11), review-slice-006 (F-3),
  review-slice-007 (F-1/F-2/F-4), review-slice-008 (F-1/F-2/F-3/F-5/F-6/F-7/F-8)

---

## Findings

### F-1 — Plan-Nachzug erneut ungetragen: 7 Erweiterungen und die reset-Streichung ohne jeden Plan-Commit (9. Auftreten — Architekt-Sequenz läuft jetzt)

- `kategorie`: MEDIUM
- `quelle`: Slice-Plan §3 · Modul 5 („Wer später mitnimmt …, hat den Plan
  **geändert**, nicht nur ergänzt") · laufende Konflikt-Sequenz (Modul 8 —
  1. Auftreten review-slice-003 F-1, … 8. Auftreten review-slice-008 F-1;
  dort als Architekt-Sequenz-Aufruf angekündigt) ·
  `BEO-PGC/plan-nachzug` (1×; mit diesem Vorgang 2×)
- `pfad`: `docs/plan/planning/in-progress/slice-009-consumer-verwaltung.md`
  §3 (zwei Zeilen) gegen die gelieferte Datei-Menge des Ranges
- `befund`: Sieben gelieferte Erweiterungen stehen außerhalb der §3-Tabelle —
  `port/outbound/consumerstate.go`, `port/inbound/consumer.go`,
  `usecase/position/`, `usecase/remove/`, `queries/queries.go` (7 Konstanten
  neu), `run-store-tests.sh` (Rollout-Vorlauf + Readiness-Änderung) und der
  Spalten-Rename `position` → `acknowledged_position` — und der Range enthält
  **keinen** Plan-Commit: `git diff be191af..HEAD -- docs/plan/planning/**`
  ist leer. Die Ausweitung ist teilweise vom eigenen Plan verlangt (§1-Ziel
  nennt „Position" als Use Case, §3 führt es nicht) und trägt einen inneren
  Widerspruch: §3 nennt `usecase/reset`, realisiert wurde er nicht. Die
  Nicht-Realisierung ist **legitim** — das Lastenheft führt keine
  Reset-Anforderung (`LH-FA-CON-001…006` kennen keinen Reset), [`ADR-0028`](../plan/adr) nennt
  `ResetConsumerUseCase` nur unter „Vorgesehen sind", und [`ADR-0013`](../plan/adr) verlangt
  nur, dass ein Reset nicht als ACK durchkommt (erfüllt: der ACK verläuft nur
  vorwärts) — aber der Plan trägt die Streichung nicht: Die §3-Zeile bleibt
  als Behauptung stehen, die Reduktion ist dokumentiert nirgends. Neuntes
  Auftreten der Klasse; nach der Ankündigung in review-slice-008 ist die
  Konflikt-Sequenz jetzt über den **Architect** zu führen, nicht ein zehntes
  MEDIUM-Weiterzählen.
- `verifizierbar`: ja — Datei-Menge des Ranges gegen die §3-Tabelle; leerer
  Docs-Diff im Implementer-Range
- `klasse`: Plan-Erweiterung ohne Plan-Nachzug (9. Auftreten — Sequenz geht
  an den Architect)

### F-2 — Zwei neue Use Cases fehlen in der ADR-0028-Liste (2. Auftreten)

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0028`](../plan/adr) (Inbound-Use-Case-Liste) · review-slice-008 F-6
- `pfad`: `internal/application/port/inbound/consumer.go:89-98`
  (`GetConsumerPositionUseCase`, `RemoveConsumerUseCase`) ·
  `docs/plan/adr/0028-inbound-use-cases.md` §Entscheidung (9-Elemente-Liste)
- `befund`: Die Liste des ADR führt `RegisterConsumerUseCase`,
  `AcknowledgeConsumerUseCase`, `ResetConsumerUseCase` — `GetConsumerPositionUseCase`
  und `RemoveConsumerUseCase` stehen nicht darin; die Commit-Messagen zitieren
  [`ADR-0013`](../plan/adr)/0028 je Port. Die Erweiterung widerspricht dem ADR nicht („Vorgesehen
  sind" ist offen, [`ADR-0034`](../plan/adr) anticipiert neue Fähigkeits-Ports), aber sie hat
  wieder keinen Ausgangs-Träger — kein Folge-ADR, kein Vermerk; zweite
  Wiederholung der Klasse aus review-slice-008 F-6.
- `verifizierbar`: ja — Use-Case-/Port-Baum gegen die [`ADR-0028`](../plan/adr)-Liste (`grep`)
- `klasse`: ADR-Liste still erweitert (2. Auftreten)

### F-3 — Monotonie-Wächter am realen Adapter: die Zeilen-Sperre ist benannt, ihre Wirkung nicht demonstriert

- `kategorie`: MEDIUM
- `quelle`: Maintainability („fehlende Negativtests bei neuem öffentlichem
  Vertrag") · [`ADR-0029`](../plan/adr) (Regel 2) · `LH-FA-CON-002`
- `pfad`: `internal/adapters/driven/postgresstorage/queries/queries.go:143-151`
  (`SelectConsumerPositionLocked`, „die Zeilen-Sperre hält die
  konkurrierenden Bestätigungen desselben Consumers in Ordnung") ·
  `internal/adapters/driven/postgresstorage/consumerstate_test.go:183` (`TestAcknowledgeIsMonotonic`)
- `befund`: Der Monotonie-Vertrag ist am realen Adapter demonstriert —
  gesperrter Lese, Domänen-Vergleich (`Advance`), FK-Ende, Quell-Bindung,
  Idempotenz, Neustart-Lese, alles im Test-Diff sichtbar — aber ausschließlich
  sequenziell. Die Behauptung, gegen die die Sperre existiert (konkurrierende
  Bestätigungen desselben Consumers), wird von keinem Test getragen: Ein
  paralleler ACK-Lauf (Rückläufer gegen Vorwärts-Lauf) fehlt, damit bleibt
  die Sperre-Kette (Lese unter Sperre → Vergleich → Upsert) nur benannt.
- `verifizierbar`: ja — Store-Test mit zwei konkurrierenden ACK-Aufrufen
  (Vorwärts- und Rückläufigkeit) gegen dieselbe Consumer-Zeile
- `klasse`: Wächter benannt, nicht demonstriert (Monotonie-Sperre;
  1. Auftreten unter diesem Namen)

### F-4 — Fehlerklasse `storage` über den Träger erweitert: ChangeStore-Sentinel trägt die ConsumerState-Fehler

- `kategorie`: MEDIUM
- `quelle`: `SPEC-008` (Klasse `storage` = „Persistenzfehler im
  ChangeStore", Aktion „Kein Source-ACK") · [`ADR-0023`](../plan/adr) · Maintainability
  („unklare Fehlerbehandlung am Rand des Spec-Bereichs")
- `pfad`: `internal/application/port/outbound/changestore.go:56` (Sentinel-Text
  „Persistenzfehler im ChangeStore") ·
  `internal/adapters/driven/postgresstorage/consumerstate.go:33-34, 59-65` ·
  `internal/application/port/outbound/consumerstate.go:29-32` („endet
  sichtbar über die Klasse `storage` (`SPEC-008`)")
- `befund`: Der `ConsumerStatePort` meldet seine Fehler über
  `outbound.ErrStorage` — denselben Sentinel, dessen Text und Aktion
  (`LH-QA-REL-001.a`, „kein Source-ACK") am ChangeStore hängen. Am Rand
  liest sich die FK-Grenze „Bestätigung ohne Registrierung" — ein
  Aufrufvertrags-Verstoß des Consumers — als Infrastruktur-Persistenzfehler
  der Klasse `storage`; die Klassen-Tabelle des Pflichtenhefts nennt nur den
  ChangeStore. Die Grenze ist dokumentiert und demonstriert, aber die
  Klasse-Aktion (kein Source-ACK) hat am Consumer-ACK keinen Träger.
- `verifizierbar`: ja — Lese der Stelle gegen die `SPEC-008`-Tabelle
- `klasse`: Fehlerklasse über Träger erweitert (1. Auftreten)

### F-5 — Die zwei vorgelagerten §8-Prüfungen stehen unausgefüllt im aktiven Plan

- `kategorie`: MEDIUM
- `quelle`: Modul 5 §Zwei Schritte vor der Modus-Begründung ·
  `BEO-PGC/plan-nachzug` (Sub-Area Planning-Harness)
- `pfad`: `docs/plan/planning/in-progress/slice-009-consumer-verwaltung.md`
  §8 („**Vorgelagert — Sub-Area-Wahl prüfen:** <je berührter Sub-Area …>" ·
  „**Vorgelagert — offene Beobachtungen sichten:** <Register durchgegangen …>")
- `befund`: Beide vorgelagerten Prüfungen — die *unbedingte* Hälfte des
  Abschnitts, unabhängig von Modus und Slice-Typ — tragen noch den
  Vorlagen-Platzhalter. Im wellenlosen Betrieb ist die Plan-Sichtung die
  einzige Leser-Instanz für alles unter der 3×-Schwelle; der Plan selbst
  zitiert in §4 sogar `BEO-PGC/d-migrate-nacharbeit` — die Sichtung wurde
  offenbar getan, aber sie ist nirgends getragen, und `§8` sichtet damit
  einen ungetanen Bestand.
- `verifizierbar`: ja — Lese der Plan-Datei gegen die Vorlage
- `klasse`: Vorgelagerte §8-Prüfungen unausgefüllt (1. Auftreten)

### F-6 — Plan-Form-Defekte: Dup-DoD (3. Auftreten) und Vorlagen-Platzhalter im DoD (2. Auftreten)

- `kategorie`: MEDIUM
- `quelle`: Modul 5 §Ziel-Form: Slice · review-slice-005 F-11 (Dup-DoD,
  1. Auftreten), review-slice-008 F-8 (2. Auftreten; Platzhalter 1. Auftreten)
- `pfad`: `docs/plan/planning/in-progress/slice-009-consumer-verwaltung.md`
  §2 (Zeilen 66/67: „`make gates` grün" doppelt) · §2 Zeile 71
  („Doku-Update für <Schnittstelle X>" — Vorlagen-Platzhalter)
- `befund`: Dieselben zwei Form-Defekte wie in slice-005 und slice-008
  stehen im aktiven Plan. Das Dup-DoD erreicht damit sein drittes Auftreten
  und steigt nach der Skill-Regel („Wiederholung eines Musters, das schon
  zweimal LOW war") auf MEDIUM; der Platzhalter im DoD zählt sein zweites.
- `verifizierbar`: ja — Lese der Plan-Datei
- `klasse`: Dup-DoD (3. Auftreten — MEDIUM-Stufe) · Vorlagen-Platzhalter im
  aktiven Plan (2. Auftreten)

### F-7 — Konjunktiv-Klausel im Readiness-Kommentar (2. Auftreten)

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Kommentar-Klassen — „Konjunktiv über die
  verworfene Alternative" ist die Falsch-Form) · review-slice-008 F-7
- `pfad`: `tools/harness/run-store-tests.sh:36-38` („… und der Folgeschritt
  würde sonst gegen ihn fahren")
- `befund`: Die Begründung der echten Abfrage endet in einer Konjunktiv-Klausel
  über das nicht gewählte Verhalten („würde sonst … fahren"); der Rest des
  Kommentars beschreibt indikativ den geltenden Zustand. Die Diagnose der
  Initdb-Phase ist klassengetreu, die Klausel nicht.
- `verifizierbar`: ja — Lese der Stelle gegen `AGENTS.md` §3.7
- `klasse`: Konjunktiv-Klausel über verworfene Alternative (2. Auftreten)

### F-8 — Register-Adapter trägt keine Kennungs-Grenze (Asymmetrie am Port-Vertrag)

- `kategorie`: LOW
- `quelle`: Maintainability · `ADR-0029` (Kennungs-Invariante am
  Domänen-Konstruktor)
- `pfad`: `internal/adapters/driven/postgresstorage/consumerstate.go:54-65`
  (`Register` ohne `ErrEmptyIdentifier`-Prüfung) gegen Zeilen 74-76, 107-112,
  173-175 (Position/Acknowledge/Remove prüfen vor dem SQL-Aufruf)
- `befund`: Drei der vier Port-Methoden enden bei leerer Kennung vor dem
  ersten SQL-Aufruf über die Domänen-Invariante; `Register` lässt die leere
  Kennung zur Zeile zu (die DDL trägt `consumer_id` als `TEXT NOT NULL`, nicht
  als nicht-leer). Der Use-Case-Konstruktor (`model.NewConsumer`) fängt es —
  am Port-Vertrag ist die Grenze ungleich getragen.
- `verifizierbar`: ja — Adapter-Test `Register` mit leerer Kennung
- `klasse`: Adapter-Grenze-Asymmetrie (1. Auftreten)

### F-9 — test-store-Beschreibung in der Sensors-Tabelle hinter dem Runner-Verhalten

- `kategorie`: LOW
- `quelle`: Maintainability · `AGENTS.md` §6 Schritt 7 (Doku-Update bei
  berührtem Vertrag)
- `pfad`: `harness/README.md` §Sensors (Zeile `make test-store`) gegen
  `tools/harness/run-store-tests.sh:52-68` (d-migrate-Rollout vor dem Testlauf)
- `befund`: Der Runner rollt jetzt das CDC-Schema über `make schema-rollout`
  vor dem Testlauf aus; die Sensors-Zeile des README beschreibt nur
  „Adapter-Tests gegen reale PostgreSQL im Testcontainer" und nennt den
  Rollout-Schritt nicht — der Ziel-Vertrag hat eine Vorbedingung bekommen,
  die die Beschreibung nicht trägt.
- `verifizierbar`: ja — Lese der README-Zeile gegen den Runner
- `klasse`: Doku-Beschreibung hinter Runner-Verhalten (1. Auftreten)

### F-10 — Composition Root trägt die neuen Use Cases nicht; die Boundary steht in keinem Abschnitt

- `kategorie`: INFO
- `quelle`: Maintainability · [`ADR-0026`](../plan/adr) (Composition Root) — Verweis an
  Planner/Verifier
- `pfad`: `internal/bootstrap/` (im Range unberührt) · Slice-Plan §1
  (einziger Ausschluss: „Performance-Tuning")
- `befund`: Die vier neuen Inbound-Use-Cases sind an keinem Driving-Adapter
  verdrahtet — `internal/bootstrap` und `cmd` referenzieren sie nicht; die
  Fähigkeit ist an keiner Aufruf-Grenze erreichbar. §1 deklariert die
  Verdrahtung in keinem Ausschluss (welche Kennung übernimmt sie), §3 führt
  `internal/bootstrap/` nicht. Plan und Code sind hier konsistent (beide
  schweigen), aber die Boundary ist undeklariert — ob das Erreichbarkeit
  erst ein späterer Vorgang trägt, ist DoD-/Verifikations-Sache; der
  Plan-Nachzug aus F-1 sollte die Stelle mit deklarieren.
- `verifizierbar`: nein — kein Gate-Lauf; Lese des Baums gegen die §1-Grenze
- `klasse`: Verdrahtungs-Boundary undeklariert (1. Auftreten)

---

## Negativbefunde

- geprüft, ohne Befund: **ADR-Verstoß über den Diff** — [`ADR-0029`](../plan/adr) Regel 2
  (Monotonie läuft über `model.ConsumerPosition.Advance` im Domänen-Vergleich,
  kein SQL-Ersatz und kein administrativer Reset als ACK), [`ADR-0013`](../plan/adr)
  (Registrierung/ACK/Entfernung/Position als Fähigkeit, Lesen verändert
  Positionen nicht), [`ADR-0034`](../plan/adr) (Fähigkeits-Port), [`ADR-0042`](../plan/adr) (Transport-Typen am
  Inbound-Port, Use-Case-Adressen als Aliase), [`ADR-0043`](../plan/adr) (Rollout-Kette:
  `schema-validate` als Vorlauf über die make-Abhängigkeit, Pflicht-Report
  `plan.yaml` + Rollback-Artefakt `down.sql` + Nacharbeit-Schritt im
  Rollout-Target committet)
- geprüft, ohne Befund: **Monotonie-Vertrag am realen Adapter, sequenzieller
  Teil** — `TestAcknowledgeIsMonotonic` (Rückläufigkeit über
  `ErrPositionRegression`, Idempotenz, Fortschritt bleibt bei Abweisung
  stehen), `TestAcknowledgeCarriesSourceBinding` (`ErrSourceMismatch`),
  `TestAcknowledgeRejectsUnregisteredConsumer` (FK-Ende sichtbar über Klasse
  `storage`), `TestConsumersCarryIndependentPositions` (`LH-FA-CON-002`),
  `TestPositionCarriesRestartPersistence` (`LH-FA-CON-005` über
  Neuaufbau des Pools); der Spalten-Rename ist in der Schema-Beschreibung
  indikativ getragen (ce2900e) und in der Rollout-Kette (`plan.yaml`,
  `down.sql`) konsistent
- geprüft, ohne Befund: **Docker-only** — `run-store-tests.sh` läuft nur
  über Docker und `make`; kein lokales Toolchain-Install
- geprüft, ohne Befund: **Traceability der 7 Implementer-Commits** — je
  Betreff mindestens eine `LH-*`-/`ADR-*`-Kennung, keine Struktur-ID im
  Betreff
- geprüft, ohne Befund: **Kommentar-Klassen im neuen Code** (`AGENTS.md`
  §3.7) — Port-, Use-Case- und Adapter-Kommentare tragen Klassen
  (Zusage/Kopplung/Abgrenzung/Grenze) über den geltenden Zustand, ausgenommen
  F-7; **Datei-Abschluss-Zeilenumbrüche** (F-6-Klasse) in allen zehn
  berührten Textdateien des Ranges
- geprüft, ohne Befund: **a-check-Edges für die neuen Pakete** —
  `internal/application/port/**` referenziert nur Domain,
  `internal/application/usecase/**` nur Ports und Domain,
  `internal/adapters/**` Ports, Domain und Adapter-Interne (`queries`,
  `mapper`); keine Rückwärts-Kante
- geprüft, ohne Befund: **Schema-Rename-Begründung** — das reservierte Wort
  `position` ist in der Schema-Beschreibung, im Commit becafdf und in der
  Rollout-Kette getragen; der d-migrate-Verhalten-Bezug („textlicher Vergleich
  liest den Unterschied als Drift") ist Tool-Behauptung, nicht Gate-fällig
- geprüft, ohne Befund: **WIP-Limit und Lifecycle** — genau ein Slice in
  `in-progress/`; die drei Planungs-Commits des Ranges sind reine
  Feld-/Move-Commits mit Kennung im Betreff

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 6 |
| LOW | 3 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Plan-Erweiterung ohne Plan-Nachzug (9.
Auftreten — Sequenz geht an den Architect) · ADR-Liste still erweitert
(2. Auftreten) · Wächter benannt, nicht demonstriert (Monotonie-Sperre;
1. Auftreten) · Fehlerklasse über Träger erweitert (1. Auftreten) ·
Vorgelagerte §8-Prüfungen unausgefüllt (1. Auftreten) · Dup-DoD (3.
Auftreten — MEDIUM-Stufe) · Vorlagen-Platzhalter im aktiven Plan (2.
Auftreten) · Konjunktiv-Klausel über verworfene Alternative (2. Auftreten) ·
Adapter-Grenze-Asymmetrie (1. Auftreten) · Doku-Beschreibung hinter
Runner-Verhalten (1. Auftreten) · Verdrahtungs-Boundary undeklariert
(1. Auftreten)

**Sequenz-Beobachtung (Steering-Loop):** die Klasse „Plan-Erweiterung ohne
Plan-Nachzug" steht hier beim **neunten** Auftreten; review-slice-008 hat
den Architekt-Übergang für genau diesen Fall angekündigt. Dieser Report ist
das Übergabe-Artefakt (F-1 + F-5 + F-6) an den Architect; ein zehntes
Zählen ist kein Übergang. `BEO-PGC/plan-nachzug` zählt mit diesem Vorgang
2× — der Beleg (`evidence/slice-009.md`) entsteht bei der Slice-Closure.

## Verdikt

**Merge-blockierend:** nein — kein HIGH-Finding. Die Architektur trägt:
der Monotonie-Vertrag liegt im Domänen-Vergleich ([`ADR-0029`](../plan/adr) Regel 2), die
Use-Cases an den Ports mit Transport-Typen ([`ADR-0028`](../plan/adr)/0042), der Adapter auf
der d-migrate-Kette ([`ADR-0043`](../plan/adr)), und die Tests demonstrieren den sequenziellen
Vertrag am realen PostgreSQL-Pfad.

**Blockierend für Closure:** ja, in drei Punkten, bevor der Slice nach
`done/` geht:

1. **F-1 (9. Auftreten):** Plan-Nachzug in §3 — sieben Erweiterungen
   aufnehmen, die reset-Zeile streichen mit Begründung (legitime Reduktion),
   die Verdrahtungs-Boundary (F-10) deklarieren — über den **Planner**; die
   Sequenz geht danach über den **Architect** (Modul 8, Pflicht ab dem
   dritten gleichen Konflikttyp, seit review-slice-003 F-1 angekündigt).
2. **F-5/F-6 (Planner):** die zwei vorgelagerten §8-Prüfungen ausfüllen,
   Dup-DoD und Platzhalter im §2 bereinigen — dieselbe Plan-Korrektur wie
   unter 1.
3. **F-4:** die Fehlerklassen-Grenze des `ConsumerStatePort` klären (eigener
   Sentinel am Port oder deklarierter Träger des `storage`-Sentinels über den
   ChangeStore hinaus) — Implementer-Entscheidung mit ADR-Bezug.

F-2 (ADR-Liste), F-3 (Konkurrenz-Demonstration) und F-7 bis F-9 sind
Rückkante an den Implementer ohne Closure-Blocker; F-10 geht an Planner und
Verifier.

---

**Gate-Beleg:** `make gates` nach diesem Report-Commit (Lauf 2026-09-10,
Range-Head); Ergebnis im Commit-Text vermerkt.