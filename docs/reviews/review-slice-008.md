# Review-Report: slice-008 Implementer-Diff — 2026-09-10

**Review-Art:** Code-Review (Implementer-Range `296faee..HEAD`, 6 Implementer-Commits)
— *wogegen*: Slice-Plan §1/§2/§3/§6 (Plan-Treue, Bewertung der fünf gemeldeten
Risiken), ADR-Bezüge ([`ADR-0028`](../plan/adr/README.md), 0042, 0034, 0026, 0039
als Rest, 0010, 0044), [`ADR-0028`](../plan/adr/README.md)-Kanon (Use-Case-/Port-Listen), Hard Rules
(`AGENTS.md` §3.1 Docker-only, §3.7 Kommentar-Klassen), Traceability (LH-*/ADR-*
je Commit, keine Struktur-IDs), Test-Qualität, `SPEC-001`-Konsistenz der
Aktivierungs-Zeilen. Keine DoD-Prüfung — das ist der Verifier (Modul 11).

**Gegenstand:** `814a166` (Inbound-Ports `verwaltung.go` + Outbound-Port
`tableactivation.go`) · `4d06164` (Use-Case-Services `enable`/`disable`/`list`/
`status` je mit Unit-Tests) · `c7a725b` (Aktivierungs-Adapter gegen
`cdc.source_table`/`cdc.schema_version`/Publication-DDL, Real-PostgreSQL-Tests) ·
`54ab209` (Verdrahtung: EnableTable-Aufrufe vor dem Stream-Start) · `82d3250`
(Runner ohne Seed-SQL, `compose.yaml`-Kommentar, `TestMVPActivationState`) ·
`4d072c9` (Image-Beleg).

**Skill:** `.harness/skills/reviewer.md` (Stand 2026-09-09, vier repo-spezifische
HIGH-Regeln, drei MEDIUM-Klassen) · Gerüst: `docs/reviews/review-report.template.md`
(Form wie `review-slice-007.md`).

**Modell:** Claude Code (glm-5.3-flash) · **Datum:** 2026-09-10

**Eingangs-Kontext:**

- Diff `git diff 296faee..HEAD` (18 Dateien, +2042/−35: Ports, vier
  Use-Case-Pakete, Adapter inkl. `queries.go`, `wiring.go`, `compose.yaml`,
  Runner, `mvp_test.go`, `harness/image-hash.txt`)
- `docs/plan/planning/in-progress/slice-008-cdc-verwaltung-use-cases.md` (§1–§8;
  im Range **unverändert** — keine Plan-Berührung, siehe F-1)
- [`ADR-0028`](../plan/adr/README.md) (Inbound Use Cases — Kanon-Liste) ·
  [`ADR-0042`](../plan/adr/README.md) (Transport-Typen am Port, Aliase am Use
  Case, Re-Evaluierungs-Trigger) · [`ADR-0034`](../plan/adr/README.md)
  (Ports nach Fähigkeiten) · [`ADR-0039`](../plan/adr/README.md) (Struktur-Regeln
  als Rest, Baumskizze) · [`ADR-0026`](../plan/adr/README.md) (Composition Root) ·
  [`ADR-0010`](../plan/adr/README.md)
- `spec/lastenheft.md` ([`LH-FA-CFG-001`](../../spec/lastenheft.md)…004),
  `spec/pflichtenheft.md` ([`LH-FA-CFG-001.a`](../../spec/pflichtenheft.md),
  `SPEC-001`, `SPEC-004`, `SPEC-008`-Klassen-Tabelle),
  `spec/architecture.md` (`ARC-003`/`ARC-004`/`ARC-005`)
- `.a-check.yml` (Layer-Edges, `composition_root`), `AGENTS.md` §3/§5,
  `harness/conventions.md` (MR-000), Vorbestand `internal/adapters/driving/
  replication/receive/receive.go` (Bezeichner-Alphabet, Publication-Prüfung),
  `internal/adapters/driving/replication/mapper/mapper.go` (Binding-Semantik)
- vorherige Reports `review-slice-005.md` (F-6/F-9/F-11-Klassen),
  `review-slice-006.md` (F-3-Klasse), `review-slice-007.md`
  (F-1/F-2/F-4/F-6/F-7-Klassen-Zählstände)

**Gate- und Probe-Läufe (am Range-Head `4d072c9`):** `make test` grün
(treiberfrei, gepinntes Toolchain-Image; alle Use-Case-Fakes-Tests, Adapter- und
Integrationstests skippen ohne DSN) — `make test-store`/`test-replication`/
`test-integration` nicht vom Review rerun (Beleg-Pflicht des Implementer-Berichts;
DoD-Prüfung ist Verifier-Sache). `make gates` nach Report-Commit (Beleg unten).
Kein eigener Real-PostgreSQL-Probe — die Adapter-Semantik (Bestands-Entzug,
Publication-Idempotenz) ist am Test-Bestand gelesen, nicht nachgefahren.

---

## Findings

### F-1 — Plan-Nachzug fehlt trotz Implementer-Meldung (8. Auftreten — Sequenz läuft weiter)

- `kategorie`: MEDIUM
- `quelle`: Slice-Plan §3 · Modul 5 („Wer später mitnimmt …, hat den Plan
  **geändert**, nicht nur ergänzt") · laufende Konflikt-Sequenz (Modul 8 — seit
  review-slice-003 F-1; 6. Auftreten review-slice-006 F-3, 7. Auftreten
  review-slice-007 F-1)
- `pfad`: `docs/plan/planning/in-progress/slice-008-cdc-verwaltung-use-cases.md`
  §3 (drei Zeilen: `verwaltung.go`, `usecase/{enable,disable}/*.go`,
  `run-integration-tests.sh`) gegen die gelieferte Datei-Menge des Ranges
- `befund`: Acht gelieferte Dateien stehen außerhalb der §3-Tabelle —
  `port/outbound/tableactivation.go`, `usecase/list/`, `usecase/status/`,
  `tableactivation.go`/`queries.go`/`tableactivation_test.go`,
  `wiring.go`, `compose.yaml`, `mvp_test.go` — und der Range enthält **keinen**
  Plan-Commit: `git diff 296faee..HEAD -- docs/` ist leer. Der Implementer
  *meldet* den Nachzug, aber kein Artefakt trägt ihn; der Plan trägt zusätzlich
  einen inneren Widerspruch (§2-DoD fordert die Status-/List-Use-Cases, §3 führt
  sie nicht — die Ausweitung ist teilweise vom eigenen DoD verlangt und nur in
  §3 vergessen). Achttes Auftreten der Klasse; der Nachzug geht als Übergabe-
  Artefakt an den **Planner**.
- `verifizierbar`: ja — Datei-Menge des Ranges gegen die §3-Tabelle; leerer
  Docs-Diff im Range
- `klasse`: Plan-Erweiterung ohne Plan-Nachzug (8. Auftreten — Sequenz läuft,
  diesmal ohne jeden Plan-Commit)

### F-2 — Zustands-Doppeldeutigkeit: Deaktivierung mit Change-Bestand meldet sich als „aktiviert"

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Implementer-Risiko bestätigt; „unklare
  Fehlerbehandlung am Rand des Spec-Bereichs") · `SPEC-001`
  (`cdc.source_table` = „aktivierte Tabellen je Quelle")
- `pfad`: `internal/application/port/inbound/verwaltung.go:73-78`
  (`GetStatusResult` trägt nur `Enabled`) ·
  `internal/adapters/driven/postgresstorage/queries/queries.go:90-100`
  (`SelectSourceTableChanges` → Retained-Pfad belässt die Bindungs-Zeile) ·
  `internal/application/usecase/disable/service.go:74-81`
- `befund`: Nach einer Deaktivierung mit Change-Bestand bleibt die
  Bindungs-Zeile als Herkunft bestehen — `GetStatus` liest genau diese Zeile und
  meldet `Enabled=true`, `ListTables` führt die Tabelle als aktiviert; die
  Erfassung ist durch den Publication-Entzug gestoppt. Das Result `Retained`
  (verwaltung.go:61-64) trägt die Unterscheidung nur an den Deaktivierungs-
  Aufrufer — der Status-/Listen-Lesepfad hat keine Möglichkeit, den
  Herkunfts-Bestand vom Erfassungs-Zustand zu unterscheiden.
- `verifizierbar`: ja — Adapter-/Use-Case-Test: Deaktivierung mit Bestand,
  anschließend `Status`-Aufruf (meldet `Enabled=true`)
- `klasse`: Aktivierungszustand doppeldeutig (Bindungs-Zeile trägt zwei Zustände)

### F-3 — Wirksamkeits-Grenze der Deaktivierung am laufenden Walsender nur als Port-Kommentar

- `kategorie`: MEDIUM
- `quelle`: Maintainability („unklare Fehlerbehandlung am Rand des
  Spec-Bereichs"; Implementer-Risiko bestätigt) ·
  [`LH-FA-CFG-002`](../../spec/lastenheft.md) Happy Path (Zeitverhalten von
  „fortan keine Änderungen" — Konformität prüft der Verifier)
- `pfad`: `internal/application/port/outbound/tableactivation.go:84-85`
  („seine Wirkung auf einen laufenden Stream-Container liegt am Walsender
  der Quelle") — einziger Träger der Grenze
- `befund`: Der Publication-Entzug wird an einem laufenden Walsender erst nach
  dessen Neuaufbau wirksam; der Erfolg-Rückgabe steht eine zeitweilig
  fortlaufende Erfassung gegenüber. Die Grenze steht ausschließlich im
  Port-Kommentar — kein Test, kein Vertrags-Text, kein Risiko-Ausgang im
  Plan §6 trägt sie (§6 führt nur das Aktivierungs-Risiko), und kein
  Deaktivierungs-Aufrufer existiert bisher am Produktionspfad.
- `verifizierbar`: ja — `grep -n "Walsender" internal/` (nur der Kommentar);
  kein Test-Beleg für das Zeitverhalten
- `klasse`: Wirksamkeits-Grenze der Deaktivierung unbenannt (nur Kommentar-Träger)

### F-4 — Bezeichner-Alphabet-Vertrag ohne Negativtest (3. Auftreten — MEDIUM-Stufe erreicht)

- `kategorie`: MEDIUM
- `quelle`: Skill-MEDIUM („Wiederholung eines Musters, das schon zweimal LOW
  war") · [`SPEC-008`](../../spec/pflichtenheft.md)-Klasse `configuration`
- `pfad`: `internal/adapters/driven/postgresstorage/tableactivation.go:27`
  (`identifierShape`), `:65-69`, `:273-287` (`validateIdentifier`,
  `ErrActivationConfiguration`) — kein `*_test.go` berührt den Pfad
- `befund`: Der neue `configuration`-Vertrag des Aktivierungs-Adapters
  (Bezeichner außerhalb `^[a-z0-9_]{1,63}$` enden vor dem ersten SQL-Aufruf)
  trägt keinen Test auf irgendeiner Ebene: die Adapter-Tests prüfen
  Register/Unregister/Publish/Unpublish/List/Missing-Table, die Use-Case-Tests
  nur die leeren Kennungen, der Integrationstest nur gültige Bezeichner
  (`grep` nach `ErrActivationConfiguration` in Tests: ohne Treffer). Dasselbe
  Muster lief bereits zweimal LOW (review-slice-005 F-9, review-slice-007 F-6) —
  drittes Auftreten auf der MEDIUM-Stufe des Skills.
- `verifizierbar`: ja — `go test` mit ungültigem Bezeichner gegen den Adapter
  (Klasse `configuration`, kein SQL-Aufruf)
- `klasse`: Neuer Vertrag ohne Negativtest (3. Auftreten — MEDIUM-Stufe)

### F-5 — ADR-Listen still erweitert: `ListTablesUseCase` und `TableActivationPort` ohne Übergabe-Artefakt

- `kategorie`: MEDIUM
- `quelle`: Modul 8 (Implementer darf Folge-ADR vorschlagen, nicht still über
  einen Accepted-ADR-Kanon hinausführen) ·
  [`ADR-0028`](../plan/adr/0028-inbound-use-cases.md)-Liste ·
  [`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md) (als Rest
  fortgeltende Struktur-Liste) · Maintainability
- `pfad`: `internal/application/port/inbound/verwaltung.go:111-114`
  (`ListTablesUseCase`) · `internal/application/port/outbound/tableactivation.go:44-86`
  (`TableActivationPort`) ·
  `docs/plan/adr/0028-inbound-use-cases.md:31-42` ·
  `docs/plan/adr/0039-paketstruktur-detaillierung-go.md:59-63` (outbound-Liste
  ohne `TableActivationPort`; usecase-Baum ohne `list/`)
- `befund`: `ADR-0028` nennt neun Use Cases, `ListTablesUseCase` gehört nicht
  dazu; die Struktur-Liste ([`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md)-Rest) führt keinen
  `TableActivationPort` und kein `list/`-Paket. Plan §1 („EnableTable/
  DisableTable/Status/Liste je [`ADR-0028`](../plan/adr/README.md)") und die Commit-Messagen zitieren
  [`ADR-0028`](../plan/adr/README.md) je Liste — der Bezug trägt mehr als der ADR-Wortlaut. Die
  Erweiterung widerspricht dem ADR nicht („Vorgesehen sind" ist offen, und
  [`ADR-0028`](../plan/adr/README.md)/[`ADR-0034`](../plan/adr/README.md) anticipieren neue Fähigkeits-Ports), aber sie hat keinen
  Ausgangs-Träger: kein Folge-ADR, kein Vermerk, keine Closure-Zeile.
- `verifizierbar`: ja — Use-Case-/Port-Baum gegen die zwei ADR-Listen (`grep`)
- `klasse`: ADR-Liste still erweitert (1. Auftreten)

### F-6 — Bezeichner-Alphabet doppelt geführt (Kopplung deklariert, Gewinner nicht)

- `kategorie`: LOW
- `quelle`: Maintainability · HIGH-Klasse „Zwei-Quellen-Drift" hier nicht
  erfüllt — die Kopplung ist benannt, aber keine Stelle ist als Träger erklärt
- `pfad`: `internal/adapters/driven/postgresstorage/tableactivation.go:27`
  gegen `internal/adapters/driving/replication/receive/receive.go:48` —
  dasselbe Muster `^[a-z0-9_]{1,63}$`, zweimal definiert
- `befund`: Dieselbe Alphabet-Regel steht jetzt in zwei Adaptern als
  unabhängige Definition; der Aktivierungs-Kommentar nennt die Kopplung
  („dieselbe Grenze wie im Stream-Adapter"), aber keine der beiden Definitionen
  ist als Quelle deklariert — eine Änderung an einer Stelle driftet still.
  Erweiterungspotenzial: die Regel trägt weder `SPEC-*` noch ADR.
- `verifizierbar`: ja — `grep -rn "a-z0-9_" internal/ | grep -v _test`
- `klasse`: Bezeichner-Alphabet doppelt geführt (1. Auftreten)

### F-7 — Konjunktiv-Klausel über die verworfene Reihenfolge im Enable-Kommentar

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §3.7 (Kommentar-Klassen — „Konjunktiv über die
  verworfene Alternative" ist die Falsch-Form; der Kommentar als Ganzes trägt
  Klassen, die Klausel nicht)
- `pfad`: `internal/application/usecase/enable/service.go:46-47` („umgekehrt
  würde eine Publication ohne Bindungs-Zeile die ersten Changes … verlieren")
- `befund`: Die Begründung der Reihenfolge (Bindung vor Publication) steht als
  Konjunktiv über die nicht gewählte Ordnung statt als Indikativ über den
  geltenden Kopplungs-Zustand; der Rest des Kommentars ist klassengetreu.
- `verifizierbar`: ja — Lese der Stelle gegen `AGENTS.md` §3.7
- `klasse`: Konjunktiv-Klausel über verworfene Alternative (1. Auftreten)

### F-8 — Plan-Form-Defekte (Planner): Dup-DoD, Platzhalter, Ausschluss ohne Kennung, ADR-Kanon-Abweichung im DoD-Namen

- `kategorie`: LOW
- `quelle`: Baseline-Regelwerk `modul-05-planning-harness.md` §Ziel-Form: Slice ·
  [`ADR-0028`](../plan/adr/0028-inbound-use-cases.md) (Kanon `GetStatusUseCase`)
- `pfad`: `docs/plan/planning/in-progress/slice-008-cdc-verwaltung-use-cases.md`
  §2 (Zeilen 68/69: `make gates` grün doppelt) · §2 Zeile 73 („Doku-Update für
  <Schnittstelle X>" — Vorlagen-Platzhalter) · §1 („die vollständige Schicht
  folgt in späteren Wellen" — Ausschluss-Klasse 1 ohne `slice-<NNN>`-Kennung) ·
  §2 Zeile 66 (`StatusUseCase` gegen den [`ADR-0028`](../plan/adr/README.md)-Kanon `GetStatusUseCase`)
- `befund`: Vier Plan-Form-Defekte, alle Planner-Sache, der Diff berührt den
  Plan nicht. Die Dup-DoD-Klasse stand bereits (review-slice-005 F-11), der
  Ausschluss-ohne-Adresse-Defekt ebenfalls (review-slice-007 F-7) — beide
  zählen hier ihr zweites Auftreten. Der Code folgt beim Use-Case-Namen
  korrekt dem ADR (der Plan-Text ist die Abweichung); `ListTablesUseCase`
  selbst ist F-5.
- `verifizierbar`: ja — Lese der Plan-Datei gegen die Vorlage und [`ADR-0028`](../plan/adr/README.md)
- `klasse`: Dup-DoD (2. Auftreten) · Vorlagen-Platzhalter im aktiven Plan ·
  Ausschluss-Klasse 1 ohne Adresse (2. Auftreten) · DoD-Name gegen ADR-Kanon

### F-9 — Zweiter Verbindungspool (+ Stream-Verbindung) gegen denselben DSN ohne benannte Grenze

- `kategorie`: LOW
- `quelle`: Maintainability (Implementer-Risiko bewertet)
- `pfad`: `internal/bootstrap/wiring.go:156` (`NewTableActivation` baut einen
  eigenen `pgxpool` neben dem Store-Pool und der Stream-Verbindung)
- `befund`: Die Verdrahtung hält jetzt drei Verbindungen gegen dieselbe Instanz
  (Store-Pool, Aktivierungs-Pool, Stream-Verbindung); der Adapter-Kommentar
  trägt die Ein-Instanz-Rationalisierung des MVP, aber nicht die
  Multi-Pool-Folge (Verbindungs-Budget, Lebensdauer dreier Pools am
  Composition Root). Semantisch folgenlos am heutigen Pfad.
- `verifizierbar`: ja — `grep -n "pgxpool.New\|NewTableActivation" internal/bootstrap/ internal/adapters/driven/postgresstorage/`
- `klasse`: Verbindungs-Pool-Grenze unbenannt (1. Auftreten)

### F-10 — Anfangs-Version hart 1: getragen, Restrisiko nur an der Schema-Evolutions-Kante

- `kategorie`: INFO
- `quelle`: Implementer-Risiko bewertet · `SPEC-004` · [`LH-FA-SCH-004.a`](../../spec/pflichtenheft.md)
- `pfad`: `internal/bootstrap/wiring.go:171-177` (Kommentar + `Version: 1`)
- `befund`: Die hart gesetzte Anfangs-Version trägt ihren Träger: der
  Verdrahtungs-Kommentar nennt die Kopplung (`SPEC-004`, spätere Versionen über
  den Metadata-Pfad `LH-FA-SCH-004.a`) und die SchemaVersionID läuft aus der
  Umgebung (`CDC_TABLES`), nicht erfunden — der Referenzpfad des Streams liest
  dieselbe Kennung (`mapper.go` Binding). Restrisiko: die Versions-Zeile
  behauptet Version 1 ohne Abgleich mit dem realen Relation-Stand der Quelle;
  der Punkt wird erst mit der Schema-Evolution wirksam. Hinweis ohne erwartete
  Aktion am Diff.
- `verifizierbar`: ja — Lese der Verdrahtung und Mapper-Bindung
- `klasse`: Schema-Version-Konvention nur im Kommentar (Vorab-Träger)

---

## Design-Entscheidungen des Implementers — Bewertung (Prüf-Fragen)

- **Plan-Nachzug (gemeldet):** *nicht getragen* — **F-1**. Die Meldung ist
  ehrlich (die Ausweitung ist in Ziel und DoD getragen — §2 verlangt die
  Status-/List-Use-Cases ausdrücklich), aber der Nachzug selbst fehlt
  vollständig: kein Plan-Commit im Range.
- **`GetStatusUseCase` statt `StatusUseCase`:** *konform, kein Code-Befund.*
  Der Code folgt dem [`ADR-0028`](../plan/adr/README.md)-Kanon exakt (Port-Interface
  `verwaltung.go:105-109`, Service `GetStatusService`, Paket `status/` gemäß
  [`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md)-Baum); die Abweichung steht im Plan-§2-Text und ist F-8. Kein
  stiller ADR-Widerspruch.
- **`ListTablesUseCase` als Ergänzung:** *Erweiterung ohne Ausgang* — **F-5**.
  Die Folgepflicht des [`ADR-0028`](../plan/adr/README.md) („Use-Case-Liste wird mit jedem neuen Driving
  Adapter auf Vollständigkeit geprüft") ist hier nicht getriggert (kein neuer
  Driving Adapter im Range), aber die Erweiterung des Kanons braucht einen
  Träger — Folge-ADR oder Closure-Vermerk.
- **`ADR-0042` (Typen am Port, Aliase am Use Case):** *voll getragen.* Alle
  vier Use-Case-Pakete definieren Typ-Aliase (`= inbound.X`), die Definition
  liegt am Inbound-Port; der Re-Evaluierungs-Trigger (Port-Typ ohne Alias) ist
  nicht eingetreten — `grep` über `usecase/**` ohne Neuedefinition. Kein Befund.
- **`ADR-0034` (Fähigkeits-Port):** *getragen.* `TableActivationPort` bildet
  eine Konsistenzgrenze als Ganzes ab (Bindungs-/Versions-Zeilen in einem
  Store-Commit, Publication als eigene Operation, drei Entzugs-Ausgänge); kein
  Repository-je-Entity-Muster. Die Listen-Erweiterung braucht den Ausgang F-5.
- **`ADR-0026` (Composition Root):** *getragen.* Nur `bootstrap` kennt den
  Aktivierungs-Adapter und verdrahtet EnableTable vor dem Stream-Start; `main`
  ist im Range unberührt; `a-check`-Edges unangetastet (`app → ports`-Kante
  getragen, keine neue Kente entstanden). Kein Befund.
- **Enable-Reihenfolge (Bindung vor Publication):** *begründet und testgetragen*
  — die Fehlordnungs-Rationale steht als Konjunktiv-Klausel (**F-7**, Form), die
  Semantik ist korrekt: fehlgeschlagene Publication hinterlässt die Bindung,
  ein erneuter Aufruf trägt nach; umgekehrte Reihenfolge wäre die FK-Kante der
  ersten Persistenz.
- **Runner ohne Seed-SQL (Aktivierung in der Verdrahtung):** *Plan §3-Zeile und
  Slice-Ziel exakt getragen.* Der Runner trägt nur noch die Vorbedingungen
  (Quell-Tabellen, Quelle-Zeile); die dreifache Kennungs-Fixtur aus slice-007
  reduziert sich auf zweifache (compose-ENV, Test-Konstanten) — die
  Bindungs-Zeilen entstehen an genau einer Stelle (Verdrahtung). Das tote
  `PUBLICATION`-Fixtur-Element aus review-slice-007 F-4 ist mit dem Seed-SQL
  entfallen — Klasse erledigt. Kein Befund.
- **`TestMVPActivationState` liest über Status-/Listen-Use-Cases statt SQL:**
  *Richtung richtig.* Der Test verdrahtet Adapter + Use Cases am
  Composition-Root-Pfad (`test/integration` ist `composition_root` in
  `.a-check.yml`) und liest den Verdrahtungs-Stand über den Port-Vertrag statt
  über SQL — dieselbe Kopplungs-Reduktion wie in slice-007 (b). Kein Befund.

## Implementer-Risiken — Bewertung

- **(a) Disable am laufenden Feed (Wirksamkeit erst bei Walsender-Neuaufbau):**
  *bestätigt, Grenze nur Kommentar-getragen* — **F-3**. Der Port-Kommentar
  benennt die Grenze korrekt, aber die Wirksamkeits-Lücke (Erfolg-Meldung vor
  Wirkung) braucht einen Ausgang (Risiko-Eintrag, Folge-Slice oder
  Beobachtungs-Register) vor der Closure.
- **(b) Zustands-Doppeldeutigkeit `cdc.source_table`:** *bestätigt als reale
  Semantik-Lücke, nicht nur Risiko* — **F-2**. Der Retained-Pfad ist adapter-
  und use-case-seitig korrekt implementiert und getestet; die doppeldeutige
  Lesart im Status-/Listen-Pfad ist die offene Hälfte.
- **(c) Bezeichner-Alphabet `^[a-z0-9_]{1,63}$`:** *Vertrag sichtbar getragen,
  zwei Residuen.* Die Verweigerung endet sichtbar über Klasse `configuration`
  vor dem ersten SQL-Aufruf (kein stiller Pfad) und ist dieselbe Grenze wie im
  Stream-Adapter — aber ohne Negativtest (**F-4**) und doppelt definiert
  (**F-6**). Die Enge des Alphabets (keine Groß-/Quoted-Bezeichner) ist
  bewusste Grenze, am Spec-Rand (Konformität: Verifier).
- **(d) Zwei Connection-Pools:** *semantisch folgenlos, Grenze unbenannt* —
  **F-9** (LOW).
- **(e) Version hart 1:** *getragen* — **F-10** (INFO), kein Befund.

## ADR-Deckung (Ports, Use Cases, Adapter, Verdrahtung)

| ADR | Aussage | Träger im Diff |
|---|---|---|
| [`ADR-0028`](../plan/adr/README.md) | Use Cases als explizite Inbound Ports; Kanon-Liste | getragen für Capture/Enable/Disable/GetStatus — `ListTablesUseCase` außerhalb der Liste (F-5) |
| [`ADR-0042`](../plan/adr/README.md) | Transport-Typen am Port, Aliase am Use Case | getragen — Definition in `verwaltung.go`/`tableactivation.go`, vier Alias-Blöcke je Use-Case-Paket; Trigger nicht eingetreten |
| [`ADR-0034`](../plan/adr/README.md) | Ports je Fähigkeit und Konsistenzgrenze | getragen — `TableActivationPort` als eine Fähigkeit (mit F-5-Ausgangspflicht für die Listen-Erweiterung) |
| [`ADR-0026`](../plan/adr/README.md) | Composition Root verdrahtet | getragen — `wiring.go` baut Adapter + Use Case und aktiviert vor dem Stream-Start; `main` unverändert |
| [`ADR-0039`](../plan/adr/README.md) (Rest) | Struktur-Regeln, Paketbaum | getragen bis auf `list/` + `TableActivationPort` (F-5); `status/` gemäß Baum |
| [`ADR-0010`](../plan/adr/README.md) | PostgreSQL-CDC-Store | getragen — Bindungs-Zeilen gegen `cdc.source_table`/`cdc.schema_version` nach `SPEC-001` |
| [`ADR-0023`](../plan/adr/README.md)/`SPEC-008` | Fehlerklassen je Kontrakt | getragen an den Sentinel (`ErrActivationConfiguration` = `configuration`, `storageFailure` = `storage`, `ErrSourceTableMissing` am Port) — Alphabet-Pfad ohne Test (F-4) |

## Test-Qualität

Unit-Tests je Use Case decken Happy/Boundary/Negative am Fake
(inkl. Aufruf-Zähler gegen stillen Erfassungs-Pfad); der Aktivierungs-Adapter
trägt Real-PostgreSQL-Tests über `make test-store` mit den drei Entzugs-Ausgängen,
der Retained-Semantik (inkl. FK-Schutz-Probe) und der Publication-Idempotenz;
`TestMVPActivationState` liest den Verdrahtungs-Stand Ende-zu-Ende über die
Use Cases. Lücken: der Alphabet-Verweigerungspfad (F-4) und die
Deaktivierungs-Wirksamkeit am laufenden Walsender (F-3, testbar als
Adapter-Probe mit folgendem Capture-Lauf).

## Negativbefunde

- geprüft, ohne Befund: **HIGH-Klassen über den Diff** — kein ADR-Verstoß auf
  Layer-/Tool-Ebene (a-check-Edges unangetastet; `app → ports`-Kante getragen;
  `bootstrap` bleibt Composition Root), kein Sicherheits-Anti-Pattern
  (Publication-DDL interpoliert nur über das Alphabet geprüfte Literale,
  `pgx`-Parameter sonst; Fehler-Zeilen nennen Bezeichner, keine Credentials),
  kein Korrektheitsfehler im kritischen Pfad (Persist-before-ACK-Ordnung
  unberührt; der Aktivierungs-Pfad läuft vor dem Stream-Start, die
  Bindungs-Zeile vor der Publication), keine Gate-Suppression, keine Norm nur
  im Template-Kommentar, kein Chronik-tragendes Zustandsfeld, **kein
  Docker-only-Verstoß** (alle Läufe über make/docker, gepinnte Digests), kein
  Zwei-Quellen-Drift im HIGH-Sinn (F-6 ist LOW-getragen: Kopplung deklariert)
- geprüft, ohne Befund: **Traceability aller sechs Commits** — jeder trägt
  mindestens eine `LH-*`-/`ADR-*`-Kennung, alle IDs lösen auf
  ([`LH-FA-CFG-001`](../../spec/lastenheft.md)…003, [`ADR-0028`](../plan/adr),
  0010, 0026, 0044); **keine Struktur-ID im Betreff** (Klasse bleibt beim
  Stand von review-slice-007: kein weiteres Auftreten);
  `commit-traceability`-Standing-Gate läuft grün
- geprüft, ohne Befund: **§3.3 mv/content-Trennung** — kein `git mv` im Range
  (Plan-Berührung fehlt, F-1); `296faee` liegt vor dem Range
- geprüft, ohne Befund: **Kommentar-Klassen im neuen Code (§3.7)** —
  `verwaltung.go`, `tableactivation.go` (Port + Adapter), vier Use-Case-Pakete,
  `wiring.go`-Aktivierungsblock, `compose.yaml`, Runner-Heredoc — Indikativ
  über den geltenden Zustand, Klassen Zusage/Kopplung/Abgrenzung/Grenze/
  Rang-Zeiger; die eine Konjunktiv-Klausel ist F-7
- geprüft, ohne Befund: **Datei-Abschluss-Zeilenumbrüche** — alle elf neu
  berührten Textdateien enden mit `\n` (`tail -c1`-Probe; die F-3-Klasse aus
  review-slice-007 bleibt nach der `11d9472`-Reparatur ohne Wiederholung)
- geprüft, ohne Befund: **totes Fixtur-Element** — die `PUBLICATION`-Variable
  ist aus dem Runner entfernt (review-slice-007 F-4-Klasse erledigt); `SLOT`
  bleibt im Wächter genutzt
- geprüft, ohne Befund: **Idempotenz-Semantik an Port und Adapter** —
  Register/Unregister/Publish/Unpublish je Boundary-Pfad, drei Entzugs-Ausgänge
  adapter-geprüft (Retained/Removed/Absent), Enable-Idempotenz über
  `AlreadyEnabled`
- geprüft, ohne Befund: **`compose.yaml`-Konsistenz** — `CDC_TABLES` trägt beide
  aktivierten Feed-Tabellen; `feed_mvp_idle` liegt außerhalb und ist der
  Boundary-Fall des neuen Tests; die List-Assertions decken beide aktivierten
  Tabellen und den Ausschluss
- geprüft, ohne Befund: **Schema-Version-Kette** — die Verdrahtung schreibt die
  Bindungs-Kennung aus der Umgebung (`binding.SchemaVersion`), der Mapper liest
  dieselbe Kennung in `cdc.change.schema_version`; kein Kennungs-Drift am Pfad
- geprüft, ohne Befund: **Gate-Läufe am Range-Head** — `make test` grün
  (treiberfrei, alle 15 Pakete); `make gates` nach Report-Commit (Beleg unten)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 5 |
| LOW | 4 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** Plan-Erweiterung ohne Plan-Nachzug (8.
Auftreten — Sequenz läuft, diesmal ohne jeden Plan-Commit) ·
Aktivierungszustand doppeldeutig (Bindungs-Zeile trägt zwei Zustände) ·
Wirksamkeits-Grenze der Deaktivierung unbenannt (nur Kommentar-Träger) ·
Neuer Vertrag ohne Negativtest (3. Auftreten — MEDIUM-Stufe) · ADR-Liste still
erweitert (1. Auftreten) · Bezeichner-Alphabet doppelt geführt ·
Konjunktiv-Klausel über verworfene Alternative · Dup-DoD (2. Auftreten) ·
Ausschluss-Klasse 1 ohne Adresse (2. Auftreten) ·
Verbindungs-Pool-Grenze unbenannt · Schema-Version-Konvention nur im Kommentar

**Sequenz-Beobachtung (Steering-Loop):** die Klasse „Plan-Erweiterung ohne
Plan-Nachzug" steht hier beim **achten** Auftreten — die Konflikt-Sequenz
(Modul 8, Pflicht ab dem dritten gleichen Konflikttyp) läuft seit
review-slice-003 F-1 und hat ihren Übergang über den **Architect** noch nicht
genommen. Dieser Lauf liefert das Übergabe-Artefakt (F-1 + F-8); bleibt der
Nachzug beim nächsten Vorgang erneut ungetragen, ist die Sequenz über den
Architect zu führen, nicht ein weiteres LOW/MEDIUM-Weiterzählen.

## Verdikt

**Merge-blockierend:** nein — kein HIGH-Finding; die Architektur trägt
([`ADR-0028`](../plan/adr/README.md), 0042, 0034, 0026), die Aktivierungs-Kette ist in Ziel und Reihenfolge
schlüssig, die Unit-Tests laufen grün, und die Runner-Umstellung entfernt den
Seed-SQL-Pfad exakt wie der Plan-Closure-Trigger es verlangt.

**Blockierend für Closure:** ja, in vier Punkten, bevor der Slice nach `done/`
geht:

1. **F-1 (8. Auftreten):** Plan-Nachzug in §3 (bzw. §7 „Was ging anders als
   geplant") über den **Planner** — Übergabe-Artefakt: F-1 + F-8 dieses
   Reports. Der Slice geht nicht mit einem §3-Stand, der 8 von 11 Dateien des
   eigenen Ranges nicht nennt, in `done/`.
2. **F-2 (Zustands-Doppeldeutigkeit):** Ausgang vor der Closure — Risiko-Ausgang
   („weiter offen" → Beobachtungs-Register oder Folge-Slice mit Kennung, die
   den Zustands-Vertrag von `cdc.source_table` annimmt) oder
   Status-/Listen-Erweiterung, die die Herkunfts-Zeile vom
   Erfassungs-Zustand trennt.
3. **F-3 + F-4 (Wirksamkeits-Grenze, Negativtest):** die `configuration`-
   Verweigerung braucht einen Test-Träger, die Walsender-Grenze einen
   benannten Ausgang — beides Vor-Closure-Nacharbeit; die gemeldeten Risiken
   (a) und (c) aus §6 des Plans brauchen je einen Ausgang (der Plan führt
   bisher nur Risiko (a) der Aktivierung).
4. **F-5 (ADR-Listen-Erweiterung):** Ausgang über Folge-ADR (Erweiterung der
   [`ADR-0028`](../plan/adr/README.md)-/[`ADR-0039`](../plan/adr/0039-paketstruktur-detaillierung-go.md)-Listen um `ListTablesUseCase`/`TableActivationPort`) oder
   Closure-Vermerk §7 mit ADR-Bezug — keine still weiterlaufende Kanon-Lücke.

F-6–F-9 sind Vor-Closure-Nacharbeit ohne Blockier-Charakter; F-10 geht ohne
Aktion in die Closure §7. DoD- und Spec-Konformität prüft der Verifier
separat (Modul 11) — insbesondere vollständige DoD-Häkchen, `make gates`
grün an der finalen Fassung und die `test-store`-/`test-integration`-Belege.