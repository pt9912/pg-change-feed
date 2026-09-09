# Review-Report: slice-003 Implementer-Diff — 2026-09-09

**Review-Art:** Diff-Review (Implementer-Range `a227033..1983e84`, 3 Commits) —
*wogegen*: Slice-Plan §1/§3 (Plan-Treue, Bewertung der zwei gemeldeten
Design-Entscheidungen), ADR-Bezüge ([`ADR-0011`](../plan/adr)/0012, 0007/0009,
0027/0028, 0039, 0029-Regel 1, 0023), Hard Rules (`AGENTS.md` §3.7
Kommentar-Klassen, §5 Dokumentations-Regeln), Traceability (keine
superseded-ADR-Referenzen), Test-Qualität (Fake-Ports, Mutationen). Keine
DoD-Prüfung — das ist der Verifier (Modul 11).

**Gegenstand:** `58f0192` (Capture-Ports: Inbound, ChangeStore,
ReplicationAck) · `cd1a6de` (CaptureService, Persist-before-ACK) ·
`1983e84` (Fake-Port-Tests) — Basis `a227033` (slice-003 `next` →
`in-progress`).

**Skill:** `.harness/skills/reviewer.md` (Stand 2026-09-09, geschärft: vier
repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen) · Gerüst:
`docs/reviews/review-report.template.md` (Form wie `review-slice-002.md`).

**Modell:** Claude Code (glm-5.3-flash) · **Datum:** 2026-09-09

**Eingangs-Kontext:**

- Diff `git diff a227033..1983e84` (5 neue Dateien, 388 Zeilen:
  `internal/application/port/{inbound,outbound}`,
  `internal/application/usecase/capture` inkl. Tests)
- `docs/plan/planning/in-progress/slice-003-capture-persist.md` (§1–§3, §6)
- `docs/plan/adr/README.md` · [`ADR-0007`](../plan/adr) · [`ADR-0009`](../plan/adr) · [`ADR-0011`](../plan/adr) ·
  [`ADR-0012`](../plan/adr) · [`ADR-0023`](../plan/adr) · [`ADR-0027`](../plan/adr) · [`ADR-0028`](../plan/adr) · [`ADR-0029`](../plan/adr) · [`ADR-0034`](../plan/adr) ·
  [`ADR-0039`](../plan/adr) · [`ADR-0041`](../plan/adr); Status `Superseded`: [`ADR-0036`](../plan/adr)/[`ADR-0038`](../plan/adr)
- `spec/lastenheft.md` ([`LH-QA-REL-001`](../../spec/lastenheft.md), [`LH-QA-REL-002`](../../spec/lastenheft.md),
  CAP-006/007), `spec/pflichtenheft.md` ([`LH-QA-REL-001.a`](../../spec/pflichtenheft.md),
  [`LH-FA-CAP-006.a`](../../spec/pflichtenheft.md), [`SPEC-001`](../../spec/pflichtenheft.md)/002/008), Commit-Messagen der Range
- `AGENTS.md` §3/§5; `harness/conventions.md` (MR-000); `harness/sensors/
  a-check.md` (Grenze 3), `.a-check.yml`, Beobachtungs-Register
  (`BEO-PGC/a-check-null-abdeckung`)
- Stand-alone-Prüfung je Commit: `git worktree` (Arbeitsbaum unverändert) +
  `go build/vet` im gepinnten Toolchain-Container
  (`golang:1.27-alpine@sha256:cf6fca…`, docker-only, `AGENTS.md` §3.1);
  `go test ./...` am Range-Head; Mutationen gegen die Test-Suite (eigene
  Proben, siehe unten); `make a-check` am Range-Head (0 Befunde)

---

## Findings

### F-1 — Transport-Typen am Inbound-Port definiert: die Text-Hälfte von ADR-0039 bleibt ohne Architect-Verdikt

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0039`](../plan/adr) (Entscheidung: „Jeder Use Case trägt sein
  Service- und Transport-Typ-Paar (Command/Query, Result) im eigenen Paket";
  Baum-Skizze `usecase/capture/ # CaptureService, CaptureCommand,
  CaptureResult`) · Maintainability
- `pfad`: `internal/application/port/inbound/capture.go:12-40` (Definition von
  `CaptureCommand`/`CaptureResult`) ·
  `internal/application/usecase/capture/service.go:15-22` (Typ-Aliase) ·
  Commit `58f0192` (Message meldet die Platzierung)
- `befund`: Die Transport-Typen sind am Inbound-Port definiert und im
  Use-Case-Paket nur als Aliase geführt. Der [`ADR-0039`](../plan/adr)-Wortlaut legt
  die Typen beim Use Case ab; die Fitness-Zeile derselben ADR
  („driving importiert nur `port/inbound` und `domain`") verlangt sie
  effektiv am Port — ein Command, das nur im Use-Case-Paket definiert
  wäre, wäre für den Driving-Adapter nur über einen
  `driving → usecase`-Import erreichbar, den keine `.a-check.yml`-Kante
  deklariert und die Fitness-Funktion ausschließt. Der Implementer hat
  die Entscheidung im Handoff gemeldet (kein stiller Widerspruch, Modul 8)
  und die Alternativen benannt (`ports → usecase`-Import — nicht
  deklarierte Kante; rein domänengetypter Port — kein Command/Result-Träger);
  die gewählte Lösung trägt die Maschinen-Hälfte (a-check 0 Befunde, keine
  neue Kante entstanden) und erfüllt die Text-Hälfte wörtlich über die
  Aliase — aber die *Definitions*-Stelle widerspricht der natürlichen
  Lesart des Accepted-ADR. **Bewertung der Prüf-Frage: Plan-Ergänzung,
  nicht Plan-Änderung** — §1 verletzt keinen Ausschluss, die DoD-Punkte
  bleiben unberührt, §3 nannte „Command/Result über zwei Ports" ohne
  Definitions-Ort. Der Defekt ist der fehlende Verdikt: eine Abweichung
  vom Wortlaut eines Accepted-ADR braucht das Architect-Artefakt
  (Folge-ADR, das die Platzierungs-Frage schärft, oder
  Neu-Implementierung) — der Handoff ist kein Verdikt.
- `verifizierbar`: ja — Import-Menge der Range gegen die
  Fitness-Funktion; `make a-check` (0 Befunde) gegen die Kanten; ADR-Text
  gegen die Typ-Definitionen
- `klasse`: Transport-Typ-Platzierung abweichend vom ADR-Wortlaut
  (Architect-Verdikt ausstehend)

### F-2 — Design-Entscheidungen ohne Plan-Nachzug (3. Auftreten — Sequenz-Pflicht greift)

- `kategorie`: MEDIUM
- `quelle`: Slice-Plan §3/§7 · Modul 5 („Wer später mitnimmt …, hat den Plan
  geändert, nicht nur ergänzt") · Modul 8 („ab dem dritten gleichen
  Konflikttyp greift die Sequenz")
- `pfad`: `docs/plan/planning/in-progress/slice-003-capture-persist.md` (§3-
  Tabelle und §7 ohne Eintrag) gegen `internal/application/port/inbound/
  capture.go:12-40`, `internal/application/usecase/capture/service.go:20-28`
- `befund`: Die fachlichen Festlegungen dieses Diffs — (a)
  Transport-Typ-Platzierung am Inbound-Port mit Aliasen, (b)
  `ErrMissingTransaction` als neue öffentliche Grenze am Use Case, (c)
  `context.Context` in allen drei Port-Signaturen — stehen weder als
  §3-Nachzug noch als §7-Zeile im Plan; sie sind erst über den Handoff
  rekonstruierbar. Das ist das **dritte Auftreten der Klasse** (F-3 in
  review-slice-001, F-4 in review-slice-002); damit ist die
  Konflikt-Sequenz als Übergabe über den Architect **Pflicht** (Modul 8:
  1× notieren · 2× Symptom · 3× Lücke) — nicht mehr ein Hinweis, den der
  Implementer vor der Closure nachzieht. Übergabe-Artefakte der Sequenz:
  Plan-Nachzug (Planner) und das [`ADR-0039`](../plan/adr)-Verdikt aus F-1
  (Architect).
- `verifizierbar`: ja — Datei-Menge und Entscheidungen gegen die §3-Tabelle
- `klasse`: Plan-Erweiterung ohne Plan-Nachzug (3. Auftreten — Sequenz
  ausgelöst)

### F-3 — `SPEC-008` in der Commit-Message von `1983e84` (3. Auftreten — Sequenz-Pflicht greift)

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §5 („Struktur-IDs (`SPEC-<NNN>`, `ARC-<NNN>`) …
  gehören nicht in die Commit-Message") · Zählstand der Übergabe: die
  F-3-Klasse zählt 2× (review-slice-002, beide Commits gezählt)
- `pfad`: Commit `1983e84` (Body: „kein Source-ACK bei Persistenzfehler
  ([`SPEC-008`](../../spec/pflichtenheft.md)-Klasse storage über [`LH-QA-REL-001.a`](../../spec/pflichtenheft.md))")
- `befund`: Ein Commit der Range trägt eine `SPEC-*`-Kennung in der Message;
  Struktur-IDs adressieren innerhalb der Spec und gehören in Code-Kommentar
  (dort steht `SPEC-008` zutreffend als Rang-Zeiger) und Spec, nicht in die
  Commit-Message. Mit dem Übergabe-Zählstand (2×) erreicht die Klasse hier
  das **dritte Auftreten** — die Konflikt-Sequenz ist Pflicht (Modul 8);
  Übergabe-Artefakt: Konventions-Nachzug im Implementer-Briefing über den
  Architect. Die Messagen selbst bleiben unveränderlich in der Historie;
  Korrektur wirkt nur vorwärts. *(Streng je Vorgang gezählt wäre dies das
  zweite Auftreten; der Übergabe-Zählstand entscheidet — als Vorgabe
  übernommen und hier offen gelegt.)*
- `verifizierbar`: ja — `git log --format=%B a227033..1983e84` gegen
  `AGENTS.md` §5
- `klasse`: Struktur-ID in Commit-Message (3. Auftreten — Sequenz ausgelöst)

### F-4 — Grenze „leere committed Transaktion" unentschieden, kein Negativtest

- `kategorie`: MEDIUM
- `quelle`: Maintainability (fehlende Negativtests bei neuem öffentlichem
  Vertrag) · [`LH-FA-CAP-006.a`](../../spec/pflichtenheft.md) („vollständige,
  konsumierbare Quelltransaktionen")
- `pfad`: `internal/application/usecase/capture/service.go:54-72` (kein
  Guard, keine Behandlung) · `internal/application/usecase/capture/
  service_test.go:106-224` (kein Test für diese Eingabe)
- `befund`: Eine committed Transaktion ohne Changes (nur `Commit`, kein
  `AppendChange`) passiert den Pfad vollständig: sie wird persistiert und
  ihre Commit-Position bestätigt. Weder der Port-Kontrakt („vollständige,
  committed Quelltransaktion") noch der Service schließt die leere
  Transaktion aus, kein Test belegt das Verhalten, und keine Spec-Stelle
  entscheidet es — ein Driving-Adapter sieht in pgoutput BEGIN/COMMIT-Paare
  ohne relevante Changes als realen Input. Die Invariante
  `ACK(position) => durable(…)` ([`ADR-0011`](../plan/adr)) bleibt verletztungsfrei
  (dauerhaft(leer) trivial), aber die Grenze ist an einem neuen
  öffentlichen Vertrag unbenannt.
- `verifizierbar`: ja — `tx := NewOpenTransaction; tx.Commit(pos)` durch
  `Capture` persistiert und ackt (ohne Fehler)
- `klasse`: Unbenannte Grenze am öffentlichen Vertrag

### F-5 — `CommitPosition`-Flag ignoriert, Kopplung an den `Changes`-Guard nur per Kommentar

- `kategorie`: LOW
- `quelle`: Maintainability (zwei Quellen für denselben Zustand in einer
  Funktion)
- `pfad`: `internal/application/usecase/capture/service.go:67` gegen
  `internal/domain/model/transaction.go:38-41`
- `befund`: `position, _ := tx.CommitPosition()` verwirft das zweite
  Ergebnis („committed") und stützt sich dafür auf den früheren
  `Changes()`-Fehler — zwei Prüfungen für denselben Zustand, die Kopplung
  trägt nur der Kommentar („committed ist über Changes geprüft"). Ändert
  eine der beiden Semantiken, ackt der Pfad eine Null-Position still.
  Der Kommentar trägt die Kopplung im Indikativ (§3.7-konform), der Typ
  trägt sie nicht.
- `verifizierbar`: ja — `IsCommitted()`/Flag-Wege sind unabhängig
  durchbrechbar
- `klasse`: Kopplung über ignoriertes Ergebnis

### F-6 — ADR-0011-Idempotenz-Pflicht ohne Träger am Store-Port-Kontrakt

- `kategorie`: LOW
- `quelle`: [`ADR-0011`](../plan/adr) (Konsequenz: „Idempotente Persistenz … ist
  Pflicht, kein Optimismus") · [`ADR-0009`](../plan/adr) („Port-Vertrag ist eine
  öffentliche Stelle")
- `pfad`: `internal/application/port/outbound/changestore.go:10-27`
  (Kontrakt nennt Persistenz, nicht Wiederholung/Deduplizierung) ·
  `internal/application/usecase/capture/service_test.go:186-191`
  (Test-Kommentar trägt die Pflicht)
- `befund`: Die Idempotenz-Pflicht des Persistierens (Wiederholung nach
  Crash, [`SPEC-002`](../../spec/pflichtenheft.md) `change_id` als Deduplizierungsbasis) steht
  nur im Fake-Test-Kommentar und im Plan-§6-Ausgang — nicht in der
  Kontrakt-Zeile des Ports, an der die Grenz-Träger-Disziplin dieses
  Slices ([`ADR-0029`](../plan/adr)-Regel 1, plan-vermerkt) sonst konsequent gehalten wird.
  Terminisiert ist die Pflicht (realer Adapter, Welle 2, plan §6); der
  Kontrakt, den der reale Adapter dann implementiert, sagt es nicht.
- `verifizierbar`: nein — Urteil über Kontrakt-Vollständigkeit
- `klasse`: ADR-Pflicht ohne Kontrakt-Träger

### F-7 — Beobachtungs-Register: `state.md` führt einen gespeicherten, veralteten Zähler

- `kategorie`: MEDIUM
- `quelle`: Baseline-Regelwerk `modul-06-roadmap.md` §Das Beobachtungs-Register
  („Der Zähler wird abgeleitet, nicht geführt — es gibt kein Feld, in das
  man ihn schreibt") · Zwei-Quellen-Drift-Klasse — **Vorbestand, nicht Teil
  dieses Diffs**
- `pfad`: `docs/plan/planning/observations/BEO-PGC/a-check-null-abdeckung/
  state.md:1-5` („Zustand: offen — Zähler 1×") gegen `evidence/`
  (`slice-001.md`, `slice-002.md` — abgeleitet 2×) und die
  slice-002-Closure („Zähler steht damit bei 2×")
- `befund`: `state.md` schreibt den Zähler als Feld und steht damit gegen die
  eigenen Beleg-Dateien: abgeleitet sind 2×, geführt steht 1×. Genau der
  Zustand, den die Abgeleitet-Regel strukturell ausschließen will —
  die gespeicherte Zeile ist die zweite Quelle und driftet. Mit diesem
  Slice erreicht die Beobachtung 3× (siehe F-9); der Lese-Schritt der
  Closure liest den Ausgang gegen den Beleg-Stand, nicht gegen die
  gespeicherte Zeile — die Zeile ist bei der Closure zu berichtigen
  (Zustand ohne geführten Zähler, Beleg als Anker).
- `verifizierbar`: ja — Anzahl der `evidence/`-Dateien gegen die `state.md`-
  Zeile
- `klasse`: Zähler geführt statt abgeleitet (stale)

### F-8 — a-check-Grenze konkret: `adapters → ports` unterscheidet Inbound/Outbound nicht

- `kategorie`: INFO
- `quelle`: [`ADR-0041`](../plan/adr) („soweit pfadgetrieben möglich";
  Re-Evaluierungs-Trigger) · `harness/sensors/a-check.md` Grenze 3 ·
  Implementer-Risiko (c) des Handoffs
- `pfad`: `.a-check.yml:8-19` (Layer `ports` = `internal/application/port/**`
  für Inbound und Outbound; Edge `adapters → ports`) ·
  `internal/application/port/outbound/changestore.go` (erster Outbound-Port
  neben dem Inbound-Port im Baum)
- `befund`: Die [`ADR-0039`](../plan/adr)-Fitness-Zeile „driving importiert nur
  `port/inbound` und `domain`" ist pfadgetrieben nicht ausdrückbar — ein
  künftiger Driving-Adapter, der `ChangeStorePort` direkt ruft und den Use
  Case umgeht, bleibt a-check-grün. Die Grenze ist im Sensor-Doc benannt
  (Grenze 3 nennt exemplarisch die `time`-Regel); mit dem ersten
  Inbound-/Outbound-Port-Paar im selben Layer ist die
  Inbound/Outbound-Teilung jetzt die zweite konkrete Instanz. Notiz —
  keine Aktion im Diff; Zuständigkeit: Architect (Re-Evaluierungs-Trigger
  von [`ADR-0041`](../plan/adr), ergänzendes Tooling als Folge-ADR).
- `verifizierbar`: ja — `make a-check` über einen hypothetischen
  Driving→outbound-Import (bleibt grün)
- `klasse`: Sensor-Grenze konkretisiert

### F-9 — Beobachtung `BEO-PGC/a-check-null-abdeckung` erreicht mit diesem Slice 3×

- `kategorie`: INFO
- `quelle`: Beobachtungs-Register · Implementer-Risiko-Kontext des Handoffs
- `pfad`: `docs/plan/planning/observations/BEO-PGC/a-check-null-abdeckung/`
  · `internal/application/usecase/**` (neu in dieser Range)
- `befund`: Der `app`-Glob matcht seit diesem Slice existierende Dateien
  (`usecase/capture`); die Null-Abdeckung gilt nur noch für `adapters`
  (Welle 2). Mit einem dritten Beleg (`evidence/slice-003.md`) erreicht die
  Beobachtung die Schwelle — Zähler, Lese-Schritt und Ausgang sind
  Register-Sache der Slice-Closure (§7), nicht des Reviews; hier nur
  notiert. Vorbestand daneben: die gespeicherte Zähler-Zeile in `state.md`
  (F-7).
- `verifizierbar`: ja — Glob-Match gegen den Baum
- `klasse`: Register-Beobachtung erreicht 3× (Closure-Sache)

---

## Design-Entscheidungen des Implementers — Bewertung (Prüf-Fragen)

- **(a) `CaptureCommand`/`CaptureResult` am Inbound-Port, Aliase im
  Use-Case-Paket:** F-1. Ergebnis der Bewertung: **Plan-Ergänzung** (kein
  §1-Ausschluss berührt, kein DoD-Punkt verletzt, §3 legte den
  Definitions-Ort nicht fest) — aber eine Ergänzung, die die Text-Hälfte
  eines Accepted-ADR berührt und deshalb das **Architect-Verdikt** braucht
  (Folge-ADR, das [`ADR-0039`](../plan/adr) schärft, oder Neu-Implementierung), bevor
  der Slice still schließt. Die gewählte Variante ist die einzige der drei
  gemeldeten Alternativen, die beide Hälften von [`ADR-0039`](../plan/adr)
  (Fitness-Funktion *und* Use-Case-trägt-Formulierung via Alias) wörtlich
  trägt; das Architektur-Urteil darüber bleibt beim Architect.
- **(b) `context.Context` in allen Port-Signaturen:** kein Befund —
  konsistent über alle drei Ports; kein aktiver ADR-Kontrakt widerspricht
  (die `ctx`-Form stammt aus der superseded [`ADR-0038`](../plan/adr)-Skizze und
  wurde übernommen); a-check `tech-leak`/`port-impurity` läuft grün über
  die neuen Signaturen (0 Befunde).

## ADR-Deckung (Service und Ports)

| ADR | Aussage | Träger im Diff |
|---|---|---|
| [`ADR-0011`](../plan/adr) | Persist-before-ACK, Invariante an EINER Stelle | getragen — `CaptureService.Capture` ist die einzige Orchestrierungs-Stelle (persist → ack, kein zweiter ACK-Pfad); Grenz-Träger am Store-Port-Kontrakt plan-getreu gesetzt (`changestore.go:16-21`); Idempotenz-Pflicht ohne Kontrakt-Träger (F-6) |
| [`ADR-0012`](../plan/adr) | At-Least-Once, Crash-Fenster | getragen — `TestCrashBetweenPersistAndAckLeadsToReprocessing` belegt Persistenz bleibt, ACK fehlt, Wiederholung führt zu Ende; Doppel-Persist im Fake ehrlich dokumentiert |
| [`ADR-0007`](../plan/adr) | ACK als Core-gesteuerter Outbound-Port | getragen — `ReplicationAckPort` am Application Layer gerufen, Adapter entscheidet nicht |
| [`ADR-0009`](../plan/adr) | ChangeStore als Fähigkeits-Port | getragen — eine Persist-Grenze, kein Repository je Entity ([`ADR-0034`](../plan/adr) konsistent) |
| [`ADR-0023`](../plan/adr) | Fehlerklassen | kein Verstoß — Klassifikation bleibt Adapter-Sache; die Service-Logik (kein ACK bei *jedem* Persist-Fehler) hängt nicht an der Klasse |
| [`ADR-0027`](../plan/adr) | Service über zwei Ports | getragen — zwei Outbound-Ports, keine Adapter-Logik im Service |
| [`ADR-0028`](../plan/adr) | Inbound-Use-Case-Form | getragen — `CaptureInboundPort` benannt wie vorgesehen |
| [`ADR-0029`](../plan/adr) Regel 1 | Grenz-Träger am Store-Port-Kontrakt | getragen — plan-vermerkte Kontrakt-Zeile gesetzt; die Regel lebt jetzt am Port *und* im Service (drei konsistente Restate: Inbound-Port, Service, Store-Port — gleiche Ordnung, kein Drift) |
| [`ADR-0039`](../plan/adr) | Command/Result beim Use Case | F-1 — Text-Hälfte nur per Alias erfüllt, Definitions-Stelle am Port, Architect-Verdikt ausstehend |

## Test-Qualität — Fake-Ports und Mutations-Proben

Die Fakes teilen eine Ereignisliste (`service_test.go:41-79`) — die
Ordnungszusage hat einen Beleg, nicht zwei unabhängige Aufrufprotokolle.
Eigene Mutations-Proben gegen die Suite (gepinnter Container):

| Mutation | Erwartung | Ergebnis |
|---|---|---|
| ACK vor Persist (Reihenfolge getauscht) | `TestCapturePersistsBeforeAck` rot | **rot** (Crash-Test: „Store trägt 0 … wollen 1") |
| Persistenzfehler geschluckt (weiter zum ACK) | `TestCaptureDoesNotAckOnPersistenceError` rot | **rot** („Capture-Fehler = nil, wollen storage defekt") |
| `Changes()`-Guard entfernt | `TestCaptureRejectsOpenTransaction` rot | **rot** („wollen Transaktion ist nicht committed") |

Die Tests tragen die Ordnungs-Zusagen; der leere-Transaktion-Fall bleibt
ungetestet (F-4). `make a-check` am Range-Head: 0 Befunde (Edges
`app → ports`, `app → domain`, `ports → domain` erfüllt, keine undeclarierte
Kante entstanden).

## Implementer-Risiken — Bewertung

- **(a) Fake-Ordnung stärker als realer Treiber** (Welle 2): kein neuer
  Befund — Plan §6 trägt das Risiko mit Ausgang „weiter offen"; die
  Fake-Strenge (exakte Ereignisfolge) ist die richtige Richtung, bis der
  reale Adapter die Semantik belegt.
- **(b) Reale Idempotenz des Stores nicht belegt**: F-6 — die Pflicht ist
  terminisiert (Plan §6-Ausgang, [`ADR-0011`](../plan/adr)-Pflicht), trägt aber am
  Port-Kontrakt nichts.
- **(c) depguard-Lint-Lücke**: F-8 — Grenze im Sensor-Doc benannt; mit
  diesem Slice ist die Inbound/Outbound-Unterscheidung die zweite
  konkrete Lücke neben der `time`-Regel.

## Negativbefunde

- geprüft, ohne Befund: **HIGH-Klassen über den Diff** — kein ADR-Verstoß
  gegen eine Accepted-ADR auf Layer-Ebene (a-check 0 Befunde, keine
  `ports → app`-Kante, keine undeclarierte Import-Kante), keine
  Gate-Suppression, kein Sicherheits-Anti-Pattern, kein Korrektheitsfehler
  im kritischen Pfad (Persist-before-ACK ist am Service an einer Stelle
  getragen und mutations-geprüft; Retention/Lesen-Positionen im Diff nicht
  berührt), keine Norm nur im Template-Kommentar, kein
  Chronik-tragendes Zustandsfeld im Diff, kein Docker-only-Verstoß
  (Build/Test im gepinnten Container, worktree-basiert)
- geprüft, ohne Befund: **Stand-alone-Build der drei Commits** — je Commit
  `go build ./...` und `go vet ./...` im gepinnten Toolchain-Container
  (`golang:1.27-alpine@sha256:cf6fca…`), alles grün; `go test ./...` am
  Range-Head grün (3 Test-Pakete)
- geprüft, ohne Befund: **Traceability-Grundpflege der drei Commits** —
  jeder trägt mindestens eine `LH-*`-/`ADR-*`-Kennung; alle genannten IDs
  existieren (`LH-QA-REL-001/002`, `LH-FA-CAP-006/007`,
  `ADR-0007/0009/0011/0012/0027/0028/0029/0030/0039`); keine superseded
  Referenz ([`ADR-0036`](../plan/adr)/[`ADR-0038`](../plan/adr)) im Diff oder in den Messagen; alle Präfixe
  MR-000-deklariert (die `SPEC-008`-Platzierung in `1983e84` ist F-3,
  kein Präfix-Verstoß)
- geprüft, ohne Befund: **Spec-Stratum** — der Diff berührt keine
  Spec-Datei; keine Erweiterung des Technik-Stratums; die
  Fehlerklasse-`storage`-Behauptung im Service-Kommentar zitiert
  [`SPEC-008`](../../spec/pflichtenheft.md) über [`LH-QA-REL-001.a`](../../spec/pflichtenheft.md) als Rang-Zeiger, ohne den
  Technik-Stratum zu erweitern
- geprüft, ohne Befund: **Kommentar-Klassen im neuen Code (§3.7)** —
  `capture.go` (Inbound), `changestore.go` (inkl. der plan-vermerkten
  Grenz-Träger-Zeile: Grenze/Zusage im Indikativ), `replicationack.go`,
  `service.go` (inkl. Alias-Begründung und der `CommitPosition`-Kopplungs-
  zeile), `service_test.go` — Klassen Zusage/Kopplung/Abgrenzung/
  Rang-Zeiger/Grenze, kein abwesender Text, kein abgebrochener Satz, keine
  verworfene Alternative im Konjunktiv; die [`ADR-0011`](../plan/adr)-Ordnung ist in drei
  Kommentaren konsistent gleich wiederholt (kein Drift zwischen ihnen)
- geprüft, ohne Befund: **§3-Datei-Menge des Plans** — alle sechs
  §3-Zeilen geliefert (Inbound-Port, Store-Port inkl. Kommentar-Zeile,
  ACK-Port, Use-Case inkl. Tests); keine unbudgetierte Datei; `.a-check.yml`
  unverändert und ohne Änderungsbedarf (`app`-Glob war seit Bootstrap
  deklariert, Abdeckung folgt dem Baum)
- geprüft, ohne Befund: **Datei-Abschluss** — alle fünf neuen Dateien enden
  mit Zeilenumbruch (F-6-Muster aus slice-001 nicht wiederholt)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 5 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Transport-Typ-Platzierung abweichend vom
ADR-Wortlaut (Architect-Verdikt ausstehend) · Plan-Erweiterung ohne
Plan-Nachzug (3. Auftreten — Sequenz ausgelöst) · Struktur-ID in
Commit-Message (3. Auftreten — Sequenz ausgelöst) · Unbenannte Grenze am
öffentlichen Vertrag · Kopplung über ignoriertes Ergebnis · ADR-Pflicht
ohne Kontrakt-Träger · Zähler geführt statt abgeleitet (stale) ·
Sensor-Grenze konkretisiert · Register-Beobachtung erreicht 3×
(Closure-Sache)

## Verdikt

**Merge-blockierend:** nein — der Diff ist inhaltlich schlüssig, alle drei
Commits bauen stand-alone, die Persist-before-ACK-Invariante ist am Service
an einer Stelle getragen und mutations-geprüft, und die gemeldete
Typ-Platzierung trägt die Maschinen-Hälfte der [`ADR-0039`](../plan/adr)-Fitness.

**Blockierend für Closure:** ja, in drei Punkten, bevor der Slice nach
`done/` geht:

1. **F-1 + F-2 als Konflikt-Sequenz (Modul 8 — drittes Auftreten der
   Plan-Nachzug-Klasse):** der Plan-Nachzug (§3/§7 zu Typ-Platzierung,
   `ErrMissingTransaction`, `context`-Signaturen) und das
   Architect-Verdikt zu [`ADR-0039`](../plan/adr) (Folge-ADR-Schärfung oder
   Neu-Implementierung) laufen als Sequenz mit Übergabe-Artefakten über den
   Architect — nicht als informelle Notiz.
2. **F-3 als Konflikt-Sequenz (drittes Auftreten der Struktur-ID-Klasse):**
   Konventions-Nachzug im Implementer-Briefing über den Architect.
3. **F-4:** die leere-committed-Transaktion-Grenze ist am neuen öffentlichen
   Vertrag zu benennen (Test oder ausdrückliche Abgrenzung) — Ausgang im
   Plan §6 oder Nacharbeit vor der Closure.

F-7 geht in die Slice-Closure §7 (Register-Lese-Schritt: Ausgang gegen den
Beleg-Stand zuweisen, gespeicherte Zähler-Zeile berichtigen). DoD- und
Spec-Konformität prüft der Verifier separat (Modul 11) — insbesondere
`make gates` grün als beobachtbarer Beleg.