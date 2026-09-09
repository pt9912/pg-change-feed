# Review-Report: slice-002 Implementer-Diff — 2026-09-09

**Review-Art:** Diff-Review (Implementer-Range `9a4d8ad..f737a0a`, 3 Commits) —
*wogegen*: Slice-Plan §1/§3 (Plan-Treue, bewertete DoD-überschreitende
Ergänzung `ConsumerPosition.SourceID`), ADR-Bezüge ([`ADR-0004`](../plan/adr)/0005/0029/0039/0040,
0025, 0012/0013), Hard Rules (`AGENTS.md` §3.7 Kommentar-Klassen, §5
Dokumentations-Regeln), Traceability (keine superseded ADR-Referenzen),
Stand-alone-Build je Commit, Test-Qualität (GWT-Pfade). Keine DoD-Prüfung —
das ist der Verifier (Modul 11).

**Gegenstand:** `247d188` (Domänenmodelle, Invarianten-Konstruktoren) ·
`4001eb8` (Domänen-Tests) · `f737a0a` (ClockPort) — Basis `9a4d8ad`
(slice-002 `next` → `in-progress`). Implementer-Verlaufskorrektur: erster
Commit-Schnitt kompilierte nicht stand-alone, neu geschnitten von `9a4d8ad`.

**Skill:** `.harness/skills/reviewer.md` (Stand 2026-09-09, geschärft: vier
repo-spezifische HIGH-Regeln, drei MEDIUM-Klassen) · Gerüst:
`docs/reviews/review-report.template.md` (Form wie `review-slice-001.md`).

**Modell:** Claude Code (glm-5.3-flash) · **Datum:** 2026-09-09

**Eingangs-Kontext:**

- Diff `git diff 9a4d8ad..f737a0a` (19 neue Dateien: `internal/domain/{model,
  errors}`, `internal/application/port/outbound`, je Tests)
- `docs/plan/planning/in-progress/slice-002-domain-kern.md` (§1–§3, §6)
- `docs/plan/adr/README.md` · [`ADR-0004`](../plan/adr) · [`ADR-0005`](../plan/adr) · [`ADR-0029`](../plan/adr) ·
  [`ADR-0039`](../plan/adr) · [`ADR-0040`](../plan/adr) · [`ADR-0025`](../plan/adr) · [`ADR-0012`](../plan/adr)/0013 (Kurzlektüre) ·
  Status `Superseded`: [`ADR-0036`](../plan/adr)/[`ADR-0038`](../plan/adr)
- `spec/lastenheft.md` ([`LH-FA-DAT-001..004`](../../spec/lastenheft.md), CAP-004/005/006/008,
  CON-003/004, RET-003/004, REA-004, SCH-005), `spec/pflichtenheft.md`
  (SPEC-001..004, 008), Commit-Messagen der Range
- `AGENTS.md` §3/§5; `harness/conventions.md` (MR-000); `harness/sensors/
  a-check.md` (Grenze 3), `.a-check.yml`
- Stand-alone-Prüfung je Commit: `git worktree` host-seitig (Arbeitsbaum
  unverändert) + `go build/vet/test ./...` im gepinnten Toolchain-Container
  (`golang:1.26-alpine@sha256:ce864e…`, docker-only, `AGENTS.md` §3.1)

---

## Findings

### F-1 — ADR-0029: Invarianten 3 und 5 sind im Kommentar behauptet, im Typ-Design nicht getragen

- `kategorie`: MEDIUM
- `quelle`: [`ADR-0029`](../plan/adr) (Regel 3 „Offene Transaktionen sind nicht
  konsumierbar", Regel 5 „Retention löscht keine benötigten Changes";
  Konsequenz: „Verstöße sind zur Laufzeit unmöglich") · Maintainability
- `pfad`: `internal/domain/model/transaction.go:20-25, 82-85` (Changes() liefert
  die Changes einer **offenen** Transaktion unverändert aus) ·
  `internal/domain/model/retention.go:20-30` (AllowsDeletion verlässt sich auf
  den Parameter `allConsumersAcknowledged bool`, den der Aufrufer frei setzen
  kann) · `internal/domain/model/source.go:3-5` (Package-Kommentar „illegale
  Zustände sind zur Laufzeit unerreichbar")
- `befund`: Regel 3 ist im Kommentar der Transaktion behauptet, aber das
  Aggregat schützt sie nicht: `Changes()` ist an einer offenen Transaktion
  vollständig lesbar, ein Use Case kann sie konsumieren, ohne dass der Typ
  widerspricht. Regel 5 ist auf einen Bool-Parameter delegiert — der Wert ist
  keine Invariante, sondern eine Zusage an den Aufrufer; ein falsch gesetztes
  `true` löscht benötigte Changes. Dazu kommt der Package-Kommentar in
  `source.go`: alle Modell-Structs tragen exportierte Felder, illegale
  Zustände sind über Struktur-Literale konstruierbar (auch der Test-Helper
  `consumer_position_test.go:28-33` bedient sich so) — die Zusage
  „unerreichbar" trägt nur für die Konstruktor-Pfade, nicht für den Typ.
  Zum Vergleich getragen: Regel 2 (Advance), Regel 6 (AppendChange), Regel 7
  (NewChange) sind echte Konstruktor-/Methoden-Zwänge. Siehe auch die
  ADR-0029-Deckung-Tabelle unten.
- `verifizierbar`: ja — mechanisch: `Changes()` an offener Transaktion
  nicht leer; `AllowsDeletion(age, true)` ohne Consumer-Beleg
- `klasse`: Invarianten-Träger behauptet, nicht getragen

### F-2 — ConsumerPosition trägt die Quelle doppelt, ohne Konsistenz-Zwang

- `kategorie`: MEDIUM
- `quelle`: Maintainability (Zwei Quellen für denselben Zustand innerhalb
  eines Wertobjekts) · [`ADR-0029`](../plan/adr) (Konstruktor-Disciplin)
- `pfad`: `internal/domain/model/consumer.go:47-56` (Felder `SourceID` und
  `Position SourcePosition`; `Position` trägt seinerseits `SourceID`, siehe
  `position.go:15-18`)
- `befund`: `ConsumerPosition` hält die Quelle zweimal — als eigenes Feld und
  in `Position.SourceID`. Kein Konstruktor und keine Methode stellt die
  Kopplung her oder prüft sie: `NewConsumerPosition` lässt `Position` null,
  `Advance` überschreibt `Position`, ohne das eigene `SourceID`-Feld gegen
  `position.SourceID` zu führen (es bleibt vom Empfänger). Ein Struktur-Literal
  mit `SourceID: "src-1"` und `Position.SourceID: "src-2"` ist legal
  konstruierbar; `Advance` prüft dann gegen das freie Feld, `Acknowledged`
  gegen `Position` — zwei Verhaltensquellen, die divergieren können. Der
  Test-Helper baut das Objekt per Literal, nicht über den Konstruktor, und
  belegt damit den Umgehungs-Pfad.
- `verifizierbar`: ja — direktes Struktur-Literal mit widersprüchlicher Quelle
  kompiliert und läuft ohne Fehler
- `klasse`: Doppelter Zustandsträger ohne Zwang

### F-3 — Struktur-IDs (`ARC-001`, `ARC-004`) in Commit-Messagen

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §5 („Struktur-IDs (`SPEC-<NNN>`, `ARC-<NNN>`) …
  gehören nicht in die Commit-Message") · Negativbefund-Präzedenz
  `review-slice-001.md` (dort als eingehalten vermerkt)
- `pfad`: Commit `247d188` (Subject: „… (ARC-001, ADR-0029)") · Commit
  `f737a0a` (ID-Zeile: „ADR-0040 ADR-0039 ARC-004")
- `befund`: Zwei der drei Commits tragen `ARC-*`-Kennungen in der Message;
  `ARC-*` adressiert innerhalb der Spec (Sicht-Stratum) und ist laut
  `AGENTS.md` §5 der Commit-Message vorbehalten — LH-/ADR-Bezug ja,
  Struktur-ID nein. Einstufung MEDIUM, nicht HIGH: `ARC-*` ist in MR-000
  deklariert, die HIGH-Klasse „Traceability-/ID-Schema-Verstoß" trifft die
  Form nicht; es ist aber die erste positive Feststellung dieser Art nach der
  slice-001-Negativzeile und daher Zähl-Signal, nicht Stilfrage. Die
  Messagen bleiben unveränderlich in der Historie; Korrektur wirkt nur
  vorwärts (Konvention im Implementer-Briefing verankern).
- `verifizierbar`: ja — `git log --format=%B` gegen `AGENTS.md` §5
- `klasse`: Struktur-ID in Commit-Message

### F-4 — DoD-überschreitende Ergänzung `ConsumerPosition.SourceID` ohne Plan-Nachzug

- `kategorie`: MEDIUM
- `quelle`: Slice-Plan §3/§7 · Modul 5 („Wer später mitnimmt …, hat den Plan
  geändert, nicht nur ergänzt") · Implementer-Handoff (Ergänzung gemeldet)
- `pfad`: `docs/plan/planning/in-progress/slice-002-domain-kern.md` (§3-Tabelle
  und §7 ohne Eintrag) gegen `internal/domain/model/consumer.go:47-56`
- `befund`: Bewertung der Ergänzung selbst: **Plan-Ergänzung, nicht
  Plan-Änderung.** §1 schließt die Feldwahl nicht aus (Ausschluss-Punkte
  unberührt — Persistenz, pgoutput, Adapter, Domain-Events bleiben außen
  vor); die Quelle am `ConsumerPosition` ist aus [`SPEC-003`](../../spec/pflichtenheft.md)
  (Position trägt `source_id`) ableitbar und für die erste Bestätigung
  nötig (`Advance` bindet an die Quelle, bevor ein Vorgänger existiert) —
  der Plan hat sie nicht festgelegt, der Diff verletzt keinen §1-Ausschluss
  und keinen DoD-Punkt. Der Defekt ist die Form: die fachliche Festlegung
  (Datenstruktur-Entscheidung) steht weder als §3-Nachzug noch als §7-Zeile
  („Was ging anders als geplant") im Plan — sie ist erst über den Handoff
  rekonstruierbar. Zweites Auftreten der Klasse (F-3 in review-slice-001
  war das erste); bei einem dritten greift die Sequenz-Regel.
- `verifizierbar`: ja — Datei-Menge und Feld-Entscheidung gegen die §3-Tabelle
- `klasse`: Plan-Erweiterung ohne Plan-Nachzug (2. Auftreten)

### F-5 — Test referenziert `LH-FA-CAP-006`, Commit-Message nennt die ID nicht

- `kategorie`: LOW
- `quelle`: `AGENTS.md` §5 (Anforderungs-IDs in Commits referenzieren, was
  berührt ist)
- `pfad`: `internal/domain/model/transaction_test.go:196-214`
  (`TestLHFACAP006OpenTransactionHasNoCommitPosition`) gegen Commit `4001eb8`
  (ID-Zeile ohne `LH-FA-CAP-006`)
- `befund`: Der Commit trägt Tests zu einer Anforderung, deren ID seine
  Message nicht nennt — die umgekehrte Richtung der Traceability-Prüfung
  (Code nennt eine ID, der Commit führt sie nicht). Kein Datenverlust, aber
  die Commit-Messagen als Deckungs-Nachweis (DoD-Punkt „referenzierte
  Tests") bleiben lückenhaft.
- `verifizierbar`: ja — Test-Namen der Range gegen die Commit-Messagen
- `klasse`: Anforderungs-ID im Code ohne Commit-Referenz

### F-6 — `TimePoint.IsZero` konflatiert Unix-Epoche mit „nicht gesetzt"

- `kategorie`: LOW
- `quelle`: Maintainability (Nullwert-Semantik) · [`ADR-0040`](../plan/adr)
- `pfad`: `internal/domain/model/timepoint.go:7-11, 15-17`
- `befund`: `IsZero()` meldet `UnixNanos == 0` — der Wert 0 ist zugleich die
  Unix-Epoche (1970-01-01T00:00:00Z) und der „nicht gesetzt"-Marker.
  `NewTimePoint` ist unvalidiert (negative Nanos erlaubt), sodass beide
  Lesarten konstruierbar sind. Solange kein Use Case mit Epochen-zeitpunkten
  rechnet, ohne semantische Auswirkung — die Retention rechnet mit
  `Duration`-Abständen —, aber die Konflation ist ein Träger für einen
  künftigen „Alter des ältesten Changes"-Defekt ([`SPEC-009`](../../spec/pflichtenheft.md),
  `cdc_oldest_change_age`).
- `verifizierbar`: nein — Urteil über künftige Verwendung
- `klasse`: Nullwert-Konflation

### F-7 — Ordnungs-Test prüft seine eigene Sortierung

- `kategorie`: LOW
- `quelle`: Maintainability (Test-Qualität) · [`LH-FA-DAT-004`](../../spec/lastenheft.md)
- `pfad`: `internal/domain/model/transaction_test.go:139-159`
  (`TestLHFADAT004IntraTransactionOrderPreservedBySequence`)
- `befund`: Der Test hängt die Changes in der Reihenfolge 3, 1, 2 an, sortiert
  sie dann selbst per `sort.Slice` nach `Sequence` und prüft 1, 2, 3 — er
  testet `sort.Slice`, nicht den Anspruch des Kommentars („Reihenfolge
  innerhalb der Transaktion", [`LH-FA-DAT-004`](../../spec/lastenheft.md) Boundary). Die
  Anhang-Reihenfolge von `Changes()` (3, 1, 2) wird nie gegen die
  Sequenz-Ordnung ins Feld geführt.
- `verifizierbar`: ja — Assertion ohne Einfluss der Anhang-Reihenfolge
  nachvollziehbar
- `klasse`: Selbst-erfüllender Ordnungs-Test

### F-8 — ADR-0029 zählt „acht Invarianten", Entscheidung listet sieben

- `kategorie`: INFO
- `quelle`: [`ADR-0029`](../plan/adr) (Kontext: „acht Invarianten" ·
  Entscheidung: Regeln 1–7) — Vorbestand, nicht Teil dieses Diffs
- `pfad`: `docs/plan/adr/0029-domain-invarianten.md:24-39`
- `befund`: Zähl-Diskrepanz im Accepted-ADR; sie macht die
  Deckungs-Frage dieses Slices (welche Regeln trägt das Typ-Design?)
  unnötig mehrdeutig. Zuständig: Architect (Folge-ADR oder
  Kontext-Korrektur als Doku-Nachzug, je nach Lesart); kein Befund gegen
  den geprüften Code.
- `verifizierbar`: ja — Zählen der Listeneinträge
- `klasse`: Zähldiskrepanz im Regel-ADR

### F-9 — Beobachtung `BEO-PGC/a-check-null-abdeckung` durch diesen Slice teils entkräftet

- `kategorie`: INFO
- `quelle`: Implementer-Risiko (c) des Handoffs · Beobachtungs-Register
- `pfad`: `docs/plan/planning/observations/BEO-PGC/a-check-null-abdeckung/state.md`
  · `internal/domain/**`, `internal/application/port/**` (neu in dieser Range)
- `befund`: Die Layer-Globs `domain` und `ports` der `.a-check.yml` matchen
  seit diesem Slice existierende Dateien (Edges `ports → domain` erfüllt:
  `clock.go` importiert `model`); die Beobachtung „null Dateien gematcht"
  gilt damit nur noch für `app`/`adapters` (entstehen mit slice-003 ff.).
  Keine Aktion im Diff — Zähler und Ausgang sind Register-Sache der
  Slice-Closure (§7, Lese-Schritt), nicht des Reviews; hier nur notiert,
  damit die Closure die Teilenkräftung nicht still unterlässt.
- `verifizierbar`: ja — Glob-Match gegen den Baum
- `klasse`: Register-Beobachtung teils entkräftet (Closure-Sache)

---

## ADR-0029-Deckung (Typ-/Konstruktor-Design) — Antwort auf die Prüf-Frage

| Regel | Inhalt | Träger im Diff |
|---|---|---|
| 1 | Persist-before-ACK | **kein Typ-Träger** — kein Source-ACK-Typ im Modell-Satz; plausible Verlagerung in den Use Case (Persist-Sequenz, [`LH-QA-REL-001.a`](../../spec/pflichtenheft.md)), aber ADR-0029 sagt „erzwungen im Domain Core" — als Grenze benennen, siehe F-1-Umfeld |
| 2 | Consumer-ACK nur vorwärts | getragen — `ConsumerPosition.Advance` (`consumer.go:66-79`), idempotente Wiederholung erlaubt ([`LH-FA-CON-004`](../../spec/lastenheft.md)) |
| 3 | Offene Transaktion nicht konsumierbar | **nur behauptet** — F-1 |
| 4 | Rollbacks erzeugen keine committed Changes | implizit getragen — kein Weg zu `committed` außer `Commit`; kein Rollback-Pfad im Modell |
| 5 | Retention löscht keine benötigten Changes | **nur behauptet** — F-1 (Bool-Parameter) |
| 6 | Eindeutige Zuordnung/Reihenfolge | getragen — `AppendChange` (Transaktions-Match, doppelte Sequenz), Negativtests vorhanden |
| 7 | Change referenziert Schema-Version | getragen — `NewChange` verlangt nichtleere Referenz; Existenz-Prüfung ist Store-Sache (späterer Slice) |

Modell-Satz gegen [`ADR-0004`](../plan/adr): alle neun Kernobjekte vorhanden
(Source, SourceTable, Change, ChangeTransaction, SourcePosition, Consumer,
ConsumerPosition, SchemaVersion, RetentionPolicy). [`ADR-0005`](../plan/adr): kein
LSN-Typ in der Domain (Offset `uint64`, Mapping beim Adapter). [`ADR-0039`](../plan/adr):
`domain/model`, `domain/errors`, `application/port/outbound` genau wie
festgelegt. [`ADR-0025`](../plan/adr): kein `domain/event/` angelegt — §1-Abgrenzung
gehalten. [`ADR-0012`](../plan/adr)/0013: ACK vorwärts (Regel 2) und idempotent
(Advance erlaubt dieselbe Position) im Typ getragen.

## Implementer-Risiken — Bewertung

- **Positionsordnung als Quell-Tiebreak** (Faktor 6.1): kein Befund —
  [`LH-FA-CAP-004`](../../spec/lastenheft.md) Out-of-Scope verlangt ausdrücklich **keine**
  Totalordnung über Quellen hinweg, nur eine konsistente, eindeutige Ordnung
  ([`LH-FA-REA-004`](../../spec/lastenheft.md)); der Tiebreak ist als Sortier-Schlüssel
  dokumentiert („keine fachliche Aussage über Quellen hinweg") und mit
  Test getragen; der REA-004.a-Sortierschlüssel (Commit-Position,
  Transaktions-ID, Sequenz) bleibt Use-Case-Sache. Die fachliche Ordnung
  wird nur innerhalb der Quelle angewendet (Advance prüft Source-Mismatch
  **vor** `Before`).
- **TimePoint als Unix-Nanosekunden ohne `time`-Import** (Faktor 6.2):
  konsistent mit [`ADR-0040`](../plan/adr) (`Now() domain.TimePoint`, Fitness:
  Domain/Use-Cases importieren `time` nicht — grep: kein `time`-Import in
  `internal/`; a-check-Grenze 3 lässt die Regel ohnehin nur dem Review).
  Restrisiko als F-6.
- **BEO-PGC teils entkräftet** (Faktor 6.3): F-9, INFO — Register-Sache.

## Negativbefunde

- geprüft, ohne Befund: **HIGH-Klassen über den Diff** — kein ADR-Verstoß
  gegen eine Accepted-ADR ([`ADR-0004`](../plan/adr) Modell-Satz vollständig, [`ADR-0005`](../plan/adr)
  LSN-frei, [`ADR-0039`](../plan/adr) Paketpflege exakt, [`ADR-0025`](../plan/adr) Abgrenzung gehalten,
  [`ADR-0040`](../plan/adr) Port-Form und `time`-Import-Regel eingehalten), keine
  Gate-Suppression, kein Sicherheits-Anti-Pattern, kein Korrektheitsfehler
  im kritischen Pfad (Persist-before-ACK, Retention-Löschung, Lesen-verändert-
  Positionen sind im Diff nicht berührt), keine Norm nur im Template-Kommentar,
  kein Chronik-tragendes Zustandsfeld, kein Docker-only-Verstoß (Build-/Test-
  Läufe im gepinnten Container, worktree-basiert), kein Zwei-Quellen-Drift
  über Dateien (F-2 betrifft Felder innerhalb eines Typs, keine zwei Dateien)
- geprüft, ohne Befund: **Stand-alone-Build der drei Commits** — je Commit
  `go build ./...` und `go vet ./...` im Toolchain-Container, alles grün;
  `go test ./...` am Range-Head grün (2 Pakete mit Tests); die
  Implementer-Verlaufskorrektur (neu geschnitten ab `9a4d8ad`) ist
  nachvollziehbar und das Ergebnis trägt
- geprüft, ohne Befund: **Traceability der drei Commits** — jeder trägt
  mindestens eine `LH-*`-/`ADR-*`-Kennung; alle genannten IDs existieren
  (`LH-FA-DAT-001/002/004`, `CAP-004/005`, `CON-004`, `RET-003/004`,
  `ADR-0029/0030/0039/0040`); keine superseded Referenz
  ([`ADR-0036`](../plan/adr)/[`ADR-0038`](../plan/adr)) im Diff oder in den Messagen; alle Präfixe MR-000-deklariert
  (die `ARC-*`-Platzierung ist F-3, kein Präfix-Verstoß)
- geprüft, ohne Befund: **Spec-Stratum** — der Diff berührt keine Spec-Datei;
  keine Erweiterung des Technik-Stratums; die Feld-Entscheidung ist aus
  [`SPEC-003`](../../spec/pflichtenheft.md) ableitbar, keine Spec-Erweiterung
- geprüft, ohne Befund: **Kommentar-Klassen im neuen Code (§3.7)** —
  `errors.go`, `change.go` (inkl. der Implementer-Begründung „Row-Image-
  Präsenz je Operationstyp nicht erzwungen": Abgrenzung/Kopplung im
  Indikativ, beschreibt den Zustand, nicht die verworfene Alternative),
  `position.go`, `consumer.go`, `transaction.go`, `retention.go`,
  `timepoint.go`, `source.go`, `clock.go` — Klassen Zusage/Kopplung/
  Abgrenzung/Rang-Zeiger, kein abwesender Text, kein abgebrochener Satz;
  der Overclaim im Package-Kommentar `source.go` ist F-1 zugeordnet
- geprüft, ohne Befund: **GWT-Pfade der genannten Lastenhefts-Zeilen** —
  DAT-001 (Happy + Boundary), DAT-002 (Happy + Boundary), DAT-004 (Happy +
  Boundary), CAP-004 (Happy, intra/inter-transaktional), CAP-005 (Happy +
  Boundary), CON-004 (Boundary idempotent), RET-003 (Happy + Boundary am
  exakten Alter), RET-004 (Consumer-basiert) je mit referenzierbarem Test;
  Negative-Pfade je Konstruktor als Tabellen-/Sub-Tests; einzig CAP-006
  trägt eine Message-Lücke (F-5) und DAT-004-Boundary einen zirkulären
  Test (F-7); keine Suppression in Test- oder Produktcode
- geprüft, ohne Befund: **a-check-Abdeckung des neuen Baums** — Layer-Globs
  `domain`/`ports` matchen die neuen Pfade; Edge `ports → domain` erfüllt;
  Composition-Root-Globs unberührt (`cmd/**`); `.a-check.yml` in-range
  unverändert (Glob-Nachzug für `app`/`adapters` folgt mit slice-003)
- geprüft, ohne Befund: **Datei-Abschluss** — alle 19 neuen Dateien enden
  mit Zeilenumbruch (slice-001 F-6-Muster nicht wiederholt)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 4 |
| LOW | 3 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Invarianten-Träger behauptet, nicht
getragen · Doppelter Zustandsträger ohne Zwang · Struktur-ID in
Commit-Message · Plan-Erweiterung ohne Plan-Nachzug (2. Auftreten) ·
Anforderungs-ID im Code ohne Commit-Referenz · Nullwert-Konflation ·
Selbst-erfüllender Ordnungs-Test · Zähldiskrepanz im Regel-ADR ·
Register-Beobachtung teils entkräftet

## Verdikt

**Merge-blockierend:** nein — der Diff ist inhaltlich schlüssig, alle drei
Commits bauen stand-alone, die Invarianten-Träger für die Regeln 2, 4, 6, 7
sind echt, der `time`-Verzicht ist sauber, und die gemeldete Ergänzung ist
als Plan-Ergänzung gerechtfertigt.

**Blockierend für Closure:** ja, in einem Punkt, bevor der Slice nach
`done/` geht: F-4 — §3-Nachzug bzw. §7-Zeile („Was ging anders als geplant")
zur Feld-Entscheidung `ConsumerPosition.SourceID`, damit Plan und Diff vor
der Closure wieder deckungsgleich sind (zweites Auftreten der Klasse; ein
drittes löst die Sequenz aus, Modul 8).

**Übergabe:** F-1 und F-2 an den Implementer als Richtung (kein Code-Input
vom Reviewer): Invarianten 3 und 5 entweder wirklich in den Typ ziehen oder
die Ansage im Kommentar auf das zurücknehmen, was getragen wird; die
`SourceID`-Kopplung entweder über einen validierenden Konstruktor-Pfad
tragen oder auf den Doppel-Träger verzichten. F-3 geht als Konventions-
Hinweis an den Implementer (Struktur-IDs aus Commit-Messagen halten —
Richtung: `harness`-Konvention/§5-Nachzug, Architect-Entscheidung, da die
Messagen selbst unveränderlich sind). F-8 geht an den Architect. F-9 geht in
die Slice-Closure §7 (Register-Zeile zitieren, nicht neu formulieren).
Finding-Klassen gehen zusätzlich in den Closure-Eintrag und von dort in den
Zähler. DoD-/Spec-Konformität prüft der Verifier separat (Modul 11) —
insbesondere `make gates` grün als beobachtbarer Beleg.