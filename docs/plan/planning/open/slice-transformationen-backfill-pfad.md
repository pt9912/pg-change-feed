# Slice transformationen-backfill-pfad: Backfill-Pfad — die Regelauswertung im Run, Fail-closed um den Regelstand, Beleg am laufenden System

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-transformationen](../welle-transformationen.md).

**Bezug:** [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) (die ausgelieferte
Change trägt die transformierte Form, nicht die Rohform),
[`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Backfill des Bestands),
[`LH-QA-SEC-004`](../../../../spec/lastenheft.md) (Wert nirgends im Image),
[`LH-FA-CAP-008`](../../../../spec/lastenheft.md),
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Folgepflicht 7 (Bindung künftiger Erzeugungspfade),
[`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 2
(eine Funktion für WAL- und Backfill-Pfad) und Teilfrage 4 (Fail-closed vor dem
Commit),
[`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md) (Klasse
`schema` für eine im Run nicht anwendbare Regel, Folgepflichten 2 und 3).

**Berührte Spec-Stellen:** [`SPEC-008`](../../../../spec/pflichtenheft.md)
(Fehlerklassen), [`SPEC-029`](../../../../spec/pflichtenheft.md) (Run-Zustand,
durch `slice-backfill-spec-nachzug`) — gelesen; die Zeile `schema` von
[`SPEC-008`](../../../../spec/pflichtenheft.md) und der Satz zum Run in
[`LH-FA-CAP-009.a`](../../../../spec/pflichtenheft.md) sind Gegenstand von
`slice-transformationen-spec-nachzug` (Folgepflicht 1 von
[`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md)), nicht dieses
Plans.

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Welle-Eröffnung
[welle-transformationen](../welle-transformationen.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Backfill-Changes tragen dieselbe transformierte Form wie WAL-Changes:
der Run liest den Regelstand (Port aus `antragsweg-usecase`) je Block neu und
baut das Bild über die gemeinsame Row-Image-Funktion **mit** dem Regelsatz; die
Fail-closed-Prüfung vor dem Commit prüft zusätzlich, dass der Regelstand
demselben Stand entspricht, mit dem die Blöcke gebaut wurden; eine auf den
Bestand nicht anwendbare Regel endet den Run sichtbar; ein E2E-Beleg zeigt die
transformierte Form eines Backfill-Bestands am laufenden System. Das ist
Kopplung K2 der Welle [welle-backfill-bestand](../done/welle-backfill-bestand.md)
§5.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Änderung am WAL-Pfad, am `Assembler` oder am Antragsweg** — `kern-rename`,
  `antragsweg-*`; dieser Slice ruft die dort gelieferten Bausteine und ergänzt
  den Run.
- **Der zweite Regeltyp** — `map-value`; der Paritäts- und der Eigenschaftstest
  dieses Slice zählen die Regeltypen aus der Domänen-Menge auf und erfassen
  `map_value`, sobald er dort steht.
- **Ein zweiter Bild-Bau-Weg im Run** — der Run hat **eine** Stelle, an der
  Ausschluss- und Regelstand in das Bild eingehen (Welle
  [welle-backfill-bestand](../done/welle-backfill-bestand.md) §5 K2); ein zweiter
  Weg wäre ein zweiter Träger derselben Aussage.
- **Checkpoint, Wiederaufnahme, Live-Zustellung von Backfill-Changes** —
  Out-of-Scope der Welle
  [welle-backfill-bestand](../done/welle-backfill-bestand.md).
- **Eine Regelform, die nur der Backfill kennt** — Regeln gelten je Tabelle für
  beide Erzeugungspfade; eine Backfill-spezifische Regel ist keine Fähigkeit
  von
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md).

## 2. Definition of Done

- [ ] Regelauswertung im Run: jeder Block liest den Regelstand über den Port
      neu (neben dem Ausschlussstand) und baut das Bild über die gemeinsame
      Funktion mit dem Regelsatz; ein Backfill-Change trägt bei gleicher Zeile
      und gleicher Regelmenge ein byte-gleiches Bild wie die WAL-Change am
      Ausgang der gemeinsamen Bild-Konstruktion (Byte-Gleichheit ist die Aussage von
      [`ADR-0115`](../../adr/0115-backfill-spaltenwerte-text-ergebnisformat.md)
      Festlegung 4 an dieser Stelle; die Spec sagt über den `jsonb`-Lesepfad
      Schlüsselmenge und Werte zu:
      [`LH-FA-CAP-009.a`](../../../../spec/pflichtenheft.md) „Markierung“;
      Paritätstest, tabellengetrieben über die Regeltypen der Domäne); ein
      Backfill-Aufrufer der gemeinsamen Funktion mit leerer Regelmenge, den
      `kern-rename` hinterlassen hat, ist ersetzt. *Zu belegen durch:* `make
      test` (Race-Detector) und — wenn der Paritätstest der Backfill-Welle im
      Replication-Tier liegt (am Start gelesen) — `make test-replication`.
- [ ] Fail-closed und Nichtanwendbarkeit: die Prüfung vor dem Commit vergleicht
      zusätzlich den Regelstand (Mengengleichheit unabhängig von der
      Reihenfolge; eine Zwischenabweichung — Regel gesetzt, dann entfernt —
      wird erkannt); eine Regel, die auf den Bestand nicht anwendbar ist
      (Spalte fehlt in der Spaltenliste des Snapshots, Zielname kollidiert),
      endet den Run `failed` mit der Klasse `schema` — einmal je Run, nachdem
      der Snapshot seine Spalten liefert und bevor die Schreibtransaktion
      öffnet, mit derselben Prüffunktion der Domäne wie der Erfassungspfad
      ([`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md)
      Festlegung 1 und 2); `classifyError` bildet den Sentinel der Prüfung auf
      `schema` ab (der Kommentar „vergibt `schema` nicht“ entfällt), ein
      Wechsel des Regelstands zwischen den Lesungen endet mit `configuration`
      (Festlegung 5); der Fehler ist run-lokal, ohne Heartbeat-Fehlerzustand
      und ohne den Capture-Pfad zu stoppen
      ([`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md)
      Teilfrage 5, Festlegung 3 von `ADR-0117`); ein Eigenschaftstest im Run
      (Regeltyp × ausgeschlossene Spalte) belegt
      [`LH-QA-SEC-004`](../../../../spec/lastenheft.md) für den Backfill-Pfad.
      *Zu belegen durch:* `make test` gegen Fakes, je
      Negativfall an seine Eingabe gebunden (Mutation der Prüfung färbt den
      Test rot); dazu der Lesefehler des Regelstands je Block und unmittelbar
      vor dem Commit — der Fake scheitert ab dem n-ten Aufruf, je Aufrufstelle
      eine Mutation, die ihren Fehler verwirft
      (`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`).
- [ ] E2E-Beleg in `make test-integration`: für eine Tabelle mit
      `rename_column`-Regel trägt ein Backfill-Run den Bestand über
      `cdc.changes` und `GET /changes` mit umbenanntem Schlüssel und `origin =
      'backfill'`; mit zusätzlichem `exclude_column` auf der umbenannten Spalte
      trägt kein Backfill-Image Quellnamen, Zielnamen oder Wert
      ([`LH-QA-SEC-004`](../../../../spec/lastenheft.md)); eine im Run nicht
      anwendbare Regel endet den Run `failed`/`schema` ohne Change, der
      Erfassungspfad läuft weiter, und nach der Abhilfe (Regel entfernen, neuer
      Antrag) endet ein neuer Run `completed`
      ([`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md)
      Folgepflicht 3; Festlegung 4 führt „ohne Prozessneustart“ als *erwartet*).
      *Zu belegen durch:*
      ein realer, grüner `make test-integration`-Lauf am laufenden
      Feed-Container (Zeile im Runner-Erzeugnis
      [`docs/user/e2e-abdeckung.md`](../../../user/e2e-abdeckung.md)).
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: entfällt — keine neue Betreiber-Oberfläche; die Aussage,
      dass Regeln auch für einen Backfill gelten, steht mit den übrigen im
      Handbuch-Abschnitt von `slice-transformationen-betriebsdoku` (Adresse in
      dessen §2).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von
      der Closure der Welle
      [welle-transformationen](../welle-transformationen.md) (die Roadmap führt
      sie unter *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| Use Case des Runs (`internal/application/usecase/backfill/…`, aus `slice-backfill-run-usecase`; Ort am Start gemessen) | update | Regelstand je Block neu lesen, Bild über die gemeinsame Funktion mit dem Regelsatz, Fail-closed um den Regelstand erweitert, Nichtanwendbarkeit endet den Run. |
| Fähigkeits-Ports des Runs und ihre Fakes | update | der Regelstand-Port aus `antragsweg-usecase` geht in den Use Case ein. |
| `internal/bootstrap/wiring.go` (Backfill-Verdrahtung) | update | reicht den Regelstand-Port an den Run. |
| Paritätstest der Backfill-Welle (Ort am Start gelesen) | update | tabellengetrieben über die Regeltypen der Domäne, mit und ohne Regel. |
| `test/integration/integration_test.go`, `tools/harness/run-integration-tests.sh` | update | Backfill-Phase mit Regel; `-run`-Muster und Abdeckungs-Deklaration (`BEO-PGC/test-runner-stiller-ausschluss`, offen, 2×). |
| `docs/user/e2e-abdeckung.md` | Erzeugnis | kommt aus dem Runner, wird nicht von Hand geschrieben. |

**Übergaben aus `slice-backfill-run-usecase`** (gemeldet, kein zusätzlicher Umfang; alle
Stellen in `internal/application/usecase/backfill/service.go`):

- **Stelle des Bild-Baus.** `blockBuilder.build` baut das Bild je Zeile über
  `model.BuildRowImage(columns, row, excluded)`; das ist die einzige Stelle des Runs, an der
  ein Regelsatz eingeht.
- **Stelle der Fail-closed-Prüfung.** In `copyBlocks` liest `excludedColumns` je Block den
  Ausschlussstand neu und `sameNames` vergleicht ihn mit dem Stand des ersten Blocks; vor dem
  Commit prüfen `stillBound` die Bindung und derselbe Vergleich den Stand; ein nicht lesbarer
  Stand endet den Run wie eine Abweichung. Die Grenze der Prüfung (ein zwischen zwei
  Lesungen gesetzter und zurückgenommener Stand ist unsichtbar, der Stand trägt keine
  Historie) steht im Doc-Kommentar von `copyBlocks` und gilt für den Regelstand ebenso.
- **Port und Fake des Regelstands.** `Ports` bündelt die Pflicht-Ports des Use Cases; der
  Regelstand-Port kommt dort hinzu. Die Fakes des Use-Case-Tests lassen einen Aufruf ab dem
  n-ten scheitern (`fakeExclusion.errCall`); der Fake des Regelstands folgt diesem Muster.
- **Klassifikation.** `classifyError` (Abbildung des Sentinels auf `schema`:
  [`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md) Festlegung 1).

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „welche
Erzeugungspfade `model.Change`-Bilder bauen und ob sie den Regelstand tragen“;
beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Alle Bild-Erzeuger | `grep -rn` nach dem Funktionsnamen der gemeinsamen Funktion (am Start gemessen) über `internal --include=*.go` und `grep -rn 'json.Marshal' internal --include=*.go` | *(Implementer trägt ein)* | jede Fundstelle, die ein Row Image baut, ruft die gemeinsame Funktion mit dem Regelsatz; eine zweite Konstruktion ist ein Befund |
| Aufrufer mit leerer Regelmenge | Suchbefehl der Kern-Slice-Zeile „Aufrufer der gemeinsamen Funktion“ am Parent-Stand | *(Implementer trägt ein)* | keine leere Regelmenge im Backfill-Pfad mehr |
| Fail-closed-Aufzählungen (Bindung, Ausschlussstand) in Doc-Kommentaren und Docs | `grep -rn 'Ausschlussstand' internal docs spec harness` | *(Implementer trägt ein)* | Aufzählungen um den Regelstand ergänzen |
| Beschreibung des Backfill-Bild-Baus in Doku | `grep -rn 'Backfill' docs/user harness spec` | *(Implementer trägt ein)* | Aussagen über die Bildform des Backfills tragen die Regel-Bindung |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-backfill-run-usecase` in
`done/` liegt (Kopplung K2 der Welle
[welle-backfill-bestand](../done/welle-backfill-bestand.md) §5), **zusätzlich**
`slice-backfill-e2e` (der E2E-Beleg braucht das lauffähige System und die
Backfill-Phase im Runner) und `slice-transformationen-antragsweg-usecase`
(Regelstand-Port und Wirkung) in `done/` liegen und kein anderer Slice in
`in-progress/` liegt (WIP-Limit 1). Das Architect-Kurzverdikt zur
Nichtanwendbarkeit einer Regel im Run liegt vor: das Verdikt
[`architect-verdict-backfill-schema-klasse-rollen`](../../../reviews/architect-verdict-backfill-schema-klasse-rollen.md)
und [`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md) (`Accepted`,
Klasse `schema` im Run, run-lokal, Folgepflicht 7 von
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
für den Run ausgefüllt). **Der Übergangs-Commit `next` → `in-progress` nennt
beide** (`BEO-PGC/start-trigger-ohne-uebergabe-artefakt`, offen, 1×). Die
Abbildung, an der `ADR-0117` ansetzt, steht in `classifyError` am Use Case des Runs
(`internal/application/usecase/backfill/service.go`): ein nicht erkannter Fehler endet als
`internal`, `schema` vergibt der Run bis zur Umsetzung dieses Slice nicht
(Register: `BEO-PGC/run-fehlerklasse-schema-im-transformations-backfill`).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls
  Run-Erweiterung, Fail-closed und E2E-Beleg nicht in einem Review tragen — der
  abtrennbare Teil ist der E2E-Beleg (dritter Liefer-Punkt) als eigener Slice
  mit Start nach diesem.
- `in-progress` → `open` (blockiert): falls der Run-Use-Case der
  Backfill-Welle keine einzige Stelle für den Bild-Bau trägt (dann Plan-Nachzug
  an die Backfill-Welle statt einer Zweitkopie), oder falls die Prüffunktion der
  Regelanwendbarkeit in der Domäne keinen Ort trägt, den Erfassungspfad und Run
  gemeinsam rufen
  ([`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md)
  Festlegung 2; Architect-Frage zum Zuschnitt).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test` (Race-Detector) grün + ein
realer, grüner `make test-integration`-Lauf + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Die Prüfung der Regelanwendbarkeit steht zweimal** (eine Kopie für den Run
  neben der Domänen-Funktion des Erfassungspfads). *Erwartet, zu belegen
  durch:* der Suchlauf über die Aufrufer der Prüffunktion und Review
  ([`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md)
  Fitness Function, Zeile „Review-Prüfpflicht“). **Ausgang:** *(bei Closure)*
- **Der Bild-Bau des Runs hat zwei Wege** (Ausschluss über Bild-Funktion, Regel
  über eine zweite Stelle). *Erwartet, zu belegen durch:* der Suchlauf (Zeile
  1) und Review. **Ausgang:** *(bei Closure)*
- **Fail-closed ist zu lasch für den Regelstand**: eine Regel, die zwischen
  zwei Blöcken gesetzt und wieder entfernt wird, hinterlässt Blöcke mit
  abweichender Form. *Erwartet, zu belegen durch:* der Negativtest der
  Zwischenabweichung mit Mutation. **Ausgang:** *(bei Closure)*
- **Der Paritätstest sitzt in einem Tier, das dieser Slice nicht fährt**
  (DB-gestützt). *Erwartet, zu belegen durch:* Lesen des Ortes am Start; liegt
  er im Replication-Tier, gehört `make test-replication` zur Closure.
  **Ausgang:** *(bei Closure)*
- **Laufzeit von `make test-integration`** wächst mit der Backfill-Phase
  (`BEO-PGC/test-integration-retention-timing-flake`, verkörpert, 3×).
  *Erwartet, zu belegen durch:* ein realer Lauf; die Zahl trägt ihren Lauf
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A). **Ausgang:** *(bei
  Closure)*
- **Kommentare zum Fehlerpfad des Runs behaupten mehr, als der Code trägt**
  (`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`, offen, 2×).
  *Erwartet, zu belegen durch:* Review liest die Kommentare der neuen Zweige.
  **Ausgang:** *(bei Closure)*

## 7. Closure-Notiz

- **Was hat funktioniert:** *(zu tragen bei Closure)*
- **Was ging anders als geplant:** *(zu tragen bei Closure)*
- **Steering-Loop-Eintrag (Lerneintrag):** *(zu tragen bei Closure —
  geschärfte Regel · neuer Sensor · benannte Spec-Lücke; ohne ihn kein
  `done/`-Übergang)*
- **Beobachtungs-Register (`../observations/`):** *(je Anfall Beleg oder
  „keine Beobachtung angefallen“ als notierte Antwort)*
- **Folge-Slices:** *(zu tragen bei Closure)*
- **Risiken aus §6:** *(je ein Ausgang)*
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-transformationen](../welle-transformationen.md) (offen) — die Prüfung
  läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area
`*`/`PGC` (Greenfield); Use Case des Runs, Composition Root und Test-Runner
sind keine eigenen Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/start-trigger-ohne-uebergabe-artefakt` (offen, 1×, einschlägig —
Start-Trigger), `BEO-PGC/vorab-bedingung-nach-umsetzung-geprueft` (offen, 2×,
einschlägig — das Verdikt ist Vorab-Bedingung im Start-Trigger, nicht
Rückführungs-Bedingung), `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`
(verkörpert, 6×), `BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×,
Plan-Zeile Runner), `BEO-PGC/test-integration-retention-timing-flake`
(verkörpert, 3×, Risiko §6),
`BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (offen, 2×),
`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` (offen, 1×, für die
Run-Klasse `schema` nicht einschlägig: sie ist mit
[`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md)
(`Supersedes` für einen Satzteil von
[`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md)
Teilfrage 5) entschieden),
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf §3).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
