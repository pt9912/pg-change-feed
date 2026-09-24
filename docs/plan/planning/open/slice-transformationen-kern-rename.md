# Slice transformationen-kern-rename: Kern — Regeltyp `rename_column` in der Domäne, Regelstand und Nichtanwendbarkeits-Prüfung im Assembler

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-transformationen](../welle-transformationen.md).

**Bezug:** [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) (Happy Path,
Negative), [`LH-FA-CFG-005`](../../../../spec/lastenheft.md) (Ausschluss gilt
zuerst), [`LH-QA-SEC-004`](../../../../spec/lastenheft.md) (Wert nirgends im
Image), [`LH-FA-DAT-005`](../../../../spec/lastenheft.md) (Abwesenheit
erkennbar), [`LH-FA-CAP-008`](../../../../spec/lastenheft.md) (fehlendes Bild
bleibt fehlend), [`LH-FA-ADM-003`](../../../../spec/lastenheft.md) (sichtbarer
Fehlerzustand), [`LH-FA-SCH-004`](../../../../spec/lastenheft.md) (Fehlerpfad
der Schemaänderung),
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Teilfrage 2/3/4/5/6 und Folgepflicht 2,
[`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) Teilfrage 2
(die gemeinsame Row-Image-Funktion),
[`ADR-0059`](../../adr/0059-spaltenauswahl-mechanismus.md) (Muster
`ExcludeColumn`), [`ADR-0023`](../../adr/0023-fehlerklassifikation.md)
(Fehlerklassen).

**Berührte Spec-Stellen:** [`SPEC-002`](../../../../spec/pflichtenheft.md) (Row
Images), [`SPEC-008`](../../../../spec/pflichtenheft.md) (Fehlerklasse
`schema`), [`ARC-001`](../../../../spec/architecture.md) (Domain Core: reine
Funktionen), [`ARC-005`](../../../../spec/architecture.md) (der Assembler ist
Teil des Driving Adapters Replication Stream),
[`ARC-007`](../../../../spec/architecture.md) (`classifyRunError` in der
Composition Root) — gelesen, nicht geändert (die Spec trägt
`slice-transformationen-spec-nachzug`).

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Welle-Eröffnung
[welle-transformationen](../welle-transformationen.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** `rename_column` wirkt im Erfassungspfad: die Domäne trägt den
Regeltyp mit Konstruktor-Invarianten und einer reinen Auswertung, die **in der
gemeinsamen Row-Image-Funktion** nach dem Ausschluss ansetzt; der `Assembler`
hält den Regelstand je Bindung als unveränderliche Liste, ersetzt sie unter
`tablesMu` und erhält sie bei `AddBinding`-Merge und `setSchemaVersion`; eine
auf eine Change nicht anwendbare Regel endet den Erfassungspfad sichtbar
(`ErrTransformationNotApplicable`, Klasse `schema`), bevor ein Wert
serialisiert wird. Die Regel ist in diesem Slice nur über die
Assembler-Methoden setzbar — kein SQL-Weg.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Antragsweg, Use Cases, Schema, Dauerhaftigkeit** — `antragsweg-schema` und
  `antragsweg-usecase`; hier gibt es keine Datenbank und keinen Antrag. Ein
  Regelstand entsteht nur über `SetTransformation`.
- **K1–K4 je Tabelle** — die Konfliktfreiheits-Invarianten sind Prüfungen beim
  Antrag (Use Case, ihr Wortlaut steht in
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 3); der `Assembler` prüft die **Anwendbarkeit** an einer konkreten
  Change, nicht die Konfliktfreiheit der Regelmenge. Die Domäne trägt nur die
  Invarianten einer einzelnen Regel (Spalten- und Zielname nicht leer,
  verschieden).
- **`map_value`** — `slice-transformationen-map-value`; der Regeltyp-Satz ist
  hier ein Satz aus einem Typ, aber als Menge geführt, damit ein zweiter Typ
  keine Umbau-Kante braucht.
- **Der Backfill-Pfad** — `slice-transformationen-backfill-pfad`; ein bereits
  vorhandener Aufrufer der gemeinsamen Funktion im Backfill-Pfad übergibt bis
  dahin die leere Regelmenge (benannt in §6).
- **Die Startreihenfolge** — `slice-transformationen-start-reihenfolge`.

## 2. Definition of Done

- [ ] `rename_column` wirkt auf beide Images: der Schlüssel `column` steht
      unter `to`, an der Position seiner Quellspalte
      (Relation-Spaltenreihenfolge), der Wert bleibt unverändert; Abwesenheit
      bleibt Abwesenheit (NULL, unverändertes TOAST, ausgeschlossene Spalte,
      fehlendes Bild); `change_id`, `transaction_id`, `source_table_id`,
      `sequence`, `operation`, `schema_version`, `schema` und `table` bleiben
      unverändert; eine leere Regelmenge liefert Bytes, die dem Stand vor
      diesem Slice gleich sind. *Zu belegen durch:* `make test` (Race-Detector)
      — bestehende Mapper- und Domänen-Tests ohne geänderte Erwartungswerte,
      neue Tests je Eigenschaft.
- [ ] Der `Assembler` trägt den Regelstand und die Prüfung:
      `SetTransformation`/`RemoveTransformation` ersetzen die Liste unter
      `tablesMu`, `AddBinding`-Merge und `setSchemaVersion` erhalten sie, ein
      Leser hält seinen Schnappschuss ohne eigene Sperre; ist die `column`
      einer Regel in `event.Relation.Columns` nicht enthalten oder kollidiert
      ihr Zielname mit einer Spalte der Relation, meldet `change`
      `mapper.ErrTransformationNotApplicable` **vor** jeder Serialisierung, die
      Transaktion wird weder persistiert noch bestätigt, und `classifyRunError`
      bildet den Fehler auf `model.ErrorClassSchema` ab. *Zu belegen durch:*
      `make test` — je ein Negativtest an seine Eingabe gebunden (Mutation der
      Prüfung färbt den Test rot,
      `BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe`), darunter der Fall,
      den erst diese Prüfung fängt: Zielname kollidiert nach einer kompatiblen
      Spalten-Erweiterung; die spalten-entfernenden Fälle enden weiterhin
      vorher an `relationOther`.
- [ ] Die Fitness Function von
      [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
      steht: ein Eigenschaftstest im `mapper`-Paket (Regeltyp × ausgeschlossene
      Spalte: das Image trägt weder den Quellschlüssel noch einen Zielnamen
      noch den Quellwert, [`LH-QA-SEC-004`](../../../../spec/lastenheft.md)),
      der die Regeltypen aus der Domänen-Menge aufzählt statt aus einer zweiten
      Liste; ein Determinismus-Test (gleiche Regelmenge und Relation →
      byte-gleiches Image); ein Nebenläufigkeits-Test (die Regelliste eines
      Lesers bleibt unter gleichzeitigem `SetTransformation` stabil). *Zu
      belegen durch:* `make test` mit Race-Detector; `make a-check` grün (die
      Regeltypen liegen in `internal/domain/**` und importieren aus keiner
      anderen Schicht); `make coverage-gate` grün (Domäne und
      `replication/mapper` liegen in der netzlos gemessenen Fläche — die
      `coverage`-Stufe des Dockerfile schließt nur die Pakete
      `postgresstorage`, `postgresack`, `postgressnapshot` und
      `replication/receive` selbst aus; kein neues Paket).
- [ ] Der Kommentar-Träger folgt: der Doc-Kommentar von `TableBinding` und
      `AddBinding` nennt den Regelstand neben `ExcludedColumns`;
      `ErrTransformationNotApplicable` trägt einen Kommentar, der nur zusagt,
      was der Code trägt (der Erfassungspfad endet, kein ACK) — keine
      Behauptung, die Abhilfe wirke im gescheiterten Prozess
      ([`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
      führt sie als erwartet;
      `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad`, offen, 2×).
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: entfällt — kein öffentlicher Vertrag berührt (das
      Nachrichtenschema bleibt, keine Betreiber-Oberfläche); die Spec trägt
      `slice-transformationen-spec-nachzug`, das Handbuch
      `slice-transformationen-betriebsdoku`.
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
| `internal/domain/model/transformation.go` (Arbeitsname; + Test) | neu | Regeltyp `rename_column`, Konstruktor-Invarianten, reine Auswertungsfunktion, keine geteilten Zustände; Menge der Regeltypen als eine Quelle. |
| die gemeinsame Row-Image-Funktion in `internal/domain/model` (aus `slice-backfill-row-image-gemeinsam`; Datei und Name am Start gemessen) | update | nimmt den Regelsatz als weiteren Parameter und wertet nach dem Ausschluss aus, innerhalb derselben Schleife; alle Aufrufer nachgezogen (Suchlauf). |
| `internal/adapters/driving/replication/mapper/mapper.go` | update | `TableBinding.Transformations`, `Assembler.SetTransformation`/`RemoveTransformation`, Erhalt bei `AddBinding`-Merge und `setSchemaVersion`, Anwendbarkeits-Prüfung in `change`, `ErrTransformationNotApplicable`. |
| `internal/adapters/driving/replication/mapper/mapper_test.go` (+ Datei für den Eigenschaftstest) | update / neu | Regeltests, Eigenschaftstest [`LH-QA-SEC-004`](../../../../spec/lastenheft.md), Determinismus, Nebenläufigkeit — nach dem Muster der `ExcludeColumn`-Tests. |
| `internal/bootstrap/wiring.go` (+ Test) | update | `classifyRunError` bildet `ErrTransformationNotApplicable` auf `model.ErrorClassSchema` ab. |
| Aufrufer der gemeinsamen Funktion im Backfill-Pfad (falls vorhanden) | update | reicht bis `backfill-pfad` die leere Regelmenge, ausdrücklich benannt. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaften: „Signatur und
Aufrufer der gemeinsamen Row-Image-Funktion“, „der Feldsatz von
`TableBinding`“, „die Menge der Sentinels, die `classifyRunError` auf `schema`
abbildet“; beide Stände gemessen):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Aufrufer der gemeinsamen Funktion | `grep -rn` nach dem Funktionsnamen der gemeinsamen Funktion (am Start gemessen) über `internal --include=*.go` | *(Implementer trägt ein)* | jeden Aufrufer nachziehen; der Backfill-Pfad übergibt die leere Menge, mit Kommentar-Bezug auf `backfill-pfad`. Übergabe aus `slice-backfill-row-image-gemeinsam` (Kopplung K1), am Stand `89053d3b` nachgemessen (`git grep -n 'BuildRowImage' HEAD -- 'internal/*.go'`, ohne Tests): die Funktion ist `model.BuildRowImage` in `internal/domain/model/rowimage.go` mit positionaler Signatur `(columns []string, values []*string, excluded []string) ([]byte, error)`; sie hat zwei Aufrufstellen in `Assembler.change` (`mapper.go`, je Bild eine, `columns` einmal je Änderung über `columnNames`); ein Regel-Parameter ändert diese zwei und den künftigen Backfill-Aufrufer. Am Start neu messen |
| Anlage- und Merge-Stellen von Bindungen | `grep -rn 'TableBinding{' internal --include=*.go` | *(Implementer trägt ein)* | Stellen, die eine Bindung anlegen oder mergen, tragen den Regelstand mit (in diesem Slice: `AddBinding`-Merge, `setSchemaVersion`); die Anlage aus der Datenbank trägt `antragsweg-usecase` |
| Feldsatz-Beschreibungen von `ExcludedColumns` | `grep -rn 'ExcludedColumns' internal docs spec harness` | *(Implementer trägt ein)* | Doc-Kommentare, die den Feldsatz aufzählen, nennen den Regelstand; Doku-Träger melden an `betriebsdoku` |
| Aufzählungen der Sentinels der Klasse `schema` | `grep -rn 'ErrIncompatibleSchemaChange' internal docs spec harness` | *(Implementer trägt ein)* | Aufzählungen (Kommentar an `classifyRunError`, Doku) um den neuen Sentinel ergänzen; die Handbuch-Zeile `schema` (§6) an `betriebsdoku` melden |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-backfill-row-image-gemeinsam`
in `done/` liegt (Kopplung K1 der Welle
[welle-backfill-bestand](../welle-backfill-bestand.md) §5: die Regelauswertung
hängt an der **einen** gemeinsamen Funktion, nicht an zwei Bild-Erzeugern),
`slice-transformationen-spec-nachzug` in `done/` liegt und kein anderer Slice
in `in-progress/` liegt (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls Domäne,
  Assembler und Fitness-Tests nicht in einem Review tragen — der abtrennbare
  Teil ist der dritte Liefer-Punkt (Eigenschafts-, Determinismus- und
  Nebenläufigkeits-Tests) als eigener Slice mit Start nach diesem.
- `in-progress` → `open` (blockiert): falls die gemeinsame Funktion der
  Backfill-Welle keinen Erweiterungspunkt ohne Signaturbruch für ihre Aufrufer
  erlaubt (Architect-Frage zum Zuschnitt), oder falls `make a-check` die
  Regeltypen in der Domäne ablehnt.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + `make test` (Race-Detector) grün +
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Byte-Abweichung ohne Regel.** Der Regel-Zweig verändert das Image auch
  dann, wenn keine Regel gesetzt ist (Escaping, Schlüsselreihenfolge, Leerraum)
  — jedes erfasste Bild wäre betroffen. *Erwartet, zu belegen durch:* die
  bestehenden Mapper- und Domänen-Tests mit den Referenz-Bytes der
  Backfill-Welle ohne geänderte Erwartung. **Ausgang:** *(bei Closure)*
- **Schlüsselposition nach der Umbenennung.**
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 3 legt fest, dass der umbenannte Schlüssel die Position seiner
  Quellspalte behält; ein Umbau der Schleife könnte ihn ans Ende setzen.
  *Erwartet, zu belegen durch:* ein Test mit Regel auf einer mittleren Spalte.
  **Ausgang:** *(bei Closure)*
- **Auswertung sieht eine ausgeschlossene Spalte.** Die Struktur der Auswertung
  (gleiche Schleife, Ausschluss zuerst) ist die tragende Zusage von
  [`LH-QA-SEC-004`](../../../../spec/lastenheft.md); eine Umsortierung der
  Schritte bräche sie still. *Erwartet, zu belegen durch:* der Eigenschaftstest
  mit Mutation (Ausschluss-Check entfernen → rot). **Ausgang:** *(bei Closure)*
- **Data Race auf der Regelliste** zwischen Capture-Goroutine und
  Administrations-Goroutine. *Erwartet, zu belegen durch:* der
  Nebenläufigkeits-Test unter `-race`, Muster der `ExcludeColumn`-Tests.
  **Ausgang:** *(bei Closure)*
- **Kosten der Anwendbarkeits-Prüfung im heißen Pfad.** Sie läuft je Change und
  vergleicht Regeln gegen `Relation.Columns`; die Größenordnung (Regeln je
  Tabelle klein) ist angenommen, nicht gemessen. *Erwartet, zu belegen durch:*
  ein `go test -bench` gegen den Parent-Stand oder eine begründete
  Nicht-Messung im Bericht. **Ausgang:** *(bei Closure)*
- **Zwischenzustand im Backfill-Pfad.** Ein Aufrufer im Backfill-Pfad übergibt
  bis `backfill-pfad` die leere Regelmenge; ein Regelstand kann erst mit
  `antragsweg-usecase` überhaupt entstehen, `backfill-pfad` folgt diesem
  unmittelbar (Welle §5). *Erwartet, zu belegen durch:* der Kommentar an der
  Aufrufstelle und die Reihenfolge der Welle. **Ausgang:** *(bei Closure:
  entfallen mit der Closure von `slice-transformationen-backfill-pfad`, dessen
  §2 die Stelle nennt)*
- **Der Ort der Auswertung weicht vom ADR-Wortlaut ab.**
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  nennt „in `rowImage`“, das private `rowImage` des Mappers entfällt mit
  `slice-backfill-row-image-gemeinsam`; die Auswertung sitzt in der gemeinsamen
  Domänen-Funktion (`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab`,
  offen, 1×). Die Auslegung ist in der Kopplung K1 begründet; sie ist im Review
  prüfbar und wird nicht als Entscheidung dargestellt. **Ausgang:** *(bei
  Closure)*

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
`*`/`PGC` (Greenfield); Domäne, Replication-Mapper und Composition Root sind
keine eigenen Sub-Areas — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×, Suchlauf §3),
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 6×, DoD Punkt
2), `BEO-PGC/kommentar-behauptet-nicht-getragenen-fehlerpfad` (offen, 2×,
einschlägig — DoD Punkt 4), `BEO-PGC/slice-chronik-in-code-kommentar`
(verkörpert, 9×, die neuen Kommentare tragen keine Slice-/Wellen-Chronik),
`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` (offen, 1×, einschlägig —
Risiko §6, letzter Punkt), `BEO-PGC/a-check-null-abdeckung` (verkörpert — die
Layer-Globs sind besetzt, `make a-check` prüft die neue Domänen-Datei),
`BEO-PGC/endstufe-unter-eigenem-messgegenstand-unerreichbar` (verkörpert —
neuer Code liegt in der netzlos gemessenen Fläche, Beleg in DoD Punkt 3).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
