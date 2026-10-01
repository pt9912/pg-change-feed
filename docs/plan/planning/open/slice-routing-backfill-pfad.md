# Slice routing-backfill-pfad: Backfill-Pfad — Backfill-Changes tragen das Label des Regelstands zum Run

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-routing](../welle-routing.md).

**Bezug:** [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (Routing —
Haupt-Bezug), [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Backfill des
Bestands),
[`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
Folgepflicht 6 und Teilfrage 6 (Backfill durchläuft dieselbe Auswertung),
[`ADR-0111`](../../adr/0111-backfill-bestand-snapshot-bulk-copy.md) (Backfill),
[`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md) (Fehlerklasse
`schema` des Runs für eine **Transformationsregel**),
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Folgepflicht 7 (Formvorbild: Bindung künftiger Erzeugungspfade).

**Berührte Spec-Stellen:**
[`LH-FA-CAP-009.a`](../../../../spec/pflichtenheft.md) (Absätze Markierung,
Fail-closed vor dem Commit, Sichtbarkeit und Fehler des Runs),
[`SPEC-008`](../../../../spec/pflichtenheft.md) (Zeile `schema`, Absatz „Nicht
anwendbare Regel"), die Routing-Regelform (neue Kennung aus
`slice-routing-spec-nachzug`). Die Spec führt: der Slice setzt
`slice-routing-spec-nachzug` voraus; die Zeilen zum Run (Vorab-Bedingung V2) zieht
dieser Slice nach, falls der Spec-Nachzug sie mangels Verdikt nicht trug
(Übergabe-Block in §2).

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-routing](../welle-routing.md).
**Datum:** 2026-10-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein Backfill-Run bestimmt für jede Zeile des Bestands das Ziel der
ersten treffenden Regel (derselben Auswertung wie der WAL-Pfad, **eine**
Auswertungsstelle) und schreibt es als `route_target` in die Backfill-Change; der
Regelstand ist Teil der Fail-closed-Prüfung des Runs (eine Abweichung des Regelstands
zwischen den Blöcken endet den Run `failed`, Klasse `configuration`, wie bei den
Transformationsregeln); eine im Run nicht anwendbare Regel endet den Run nach dem
Kurzverdikt zu Vorab-Bedingung V2 der Welle. Das Label gehört nicht zum
Zeilenzustand: die Replay-Invariante (das Log ab dem Log-Anfang ergibt den
Quellstand) bleibt unberührt.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Entscheidung der Run-Behandlung selbst** — V2 ist eine Lücke im Text von
  [`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  (`ADR-0117` gilt der Transformationsregel). Der Implementer legt sie nicht aus; der
  Slice startet mit dem Architect-Kurzverdikt (§4), wie
  `slice-transformationen-backfill-pfad` mit `ADR-0117` startete.
- **Eine neue Run-Fehlerklasse** — die Klassenmenge des Runs bleibt (sieben Klassen
  des Prozesses, [`ADR-0023`](../../adr/0023-fehlerklassifikation.md)); ein Verdikt, das
  eine achte Klasse verlangte, wäre eine Folge-ADR.
- **Die Auswertung im WAL-Pfad** — `slice-routing-kern-label`; der Slice **ruft**
  dieselbe Domänen-Funktion auf und dupliziert sie nicht.
- **Umetikettieren erfasster Changes** — Entscheidung 6 der ADR: der Altbestand wird
  über einen **neuen** Backfill-Run mit dem aktuellen Regelstand erzeugt, nichts wird
  umgeschrieben.
- **E2E-Beleg am laufenden System** — `slice-routing-e2e` (Backfill-Bestand mit
  Label); dieser Slice belegt auf Unit- und Store-Ebene (`make test`,
  `make test-store`, `make test-replication` für den Snapshot-Pfad).

## 2. Definition of Done

- [ ] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) und
      [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Auswertung): eine
      Backfill-Change trägt das Ziel, das die Auswertung für den Quellwert der Zeile
      bestimmt (Herkunfts- und Inhaltsregel; Reihenfolge `order`; abwesender Wert ist
      Nicht-Treffer; keine Regel: `NULL`), unabhängig vom Block in dem die Zeile liegt;
      das Label einer Zeile ist dasselbe, das der WAL-Pfad für dieselbe Zeile und
      dieselbe Regelliste bestimmte (Vertragstest über die gemeinsame
      Domänen-Funktion); die Auswertungsstelle ist eine — der Suchlauf in §3 findet
      keine zweite. *Zu belegen durch:* Use-Case-Tabellentest (`make test`), Store-
      und Snapshot-Test (`make test-store`, `make test-replication`).
- [ ] [`LH-FA-CAP-009`](../../../../spec/lastenheft.md) (Fail-closed): der Regelstand
      der Routing-Regeln wird beim Start des Runs gelesen und vor jedem Block
      verglichen; eine Abweichung zwischen den Blöcken endet den Run `failed` mit
      der Klasse der Transformationsregeln für diesen Fall (`configuration`); die
      Replay-Invariante bleibt belegt (die bestehende Prüfung läuft unverändert
      grün). *Zu belegen durch:* Use-Case-Test mit wechselndem Regelstand zwischen
      zwei Blöcken (Eingabe-Bindung: die Eingabe ist der Wechsel), `make test`,
      Replay-Test des bestehenden E2E (`make test-integration`, nicht Teil dieses
      Slice, hier nur nicht-brechend).
- [ ] [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) Negative im Run: eine im Run
      nicht anwendbare Regel (`when.column` fehlt in den Spalten des Snapshots)
      endet nach dem Kurzverdikt zu V2 (Welle §5): erwartet — *nicht entschieden* —
      Run `failed`, Klasse `schema`, run-lokal, ohne Change, der Erfassungspfad läuft
      weiter (Muster `ADR-0117`). *Zu belegen durch:* Use-Case-Test mit einem
      Snapshot, dem die Spalte fehlt; der Testname nennt die Eingabe.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-routing-backfill-pfad.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: `spec/pflichtenheft.md` (`LH-FA-CAP-009.a`, `SPEC-008` Zeile
      `schema` und Absatz zur Nichtanwendbarkeit im Run) **soweit** der Spec-Nachzug
      sie mangels Verdikt V2 nicht trug — **Übergabe-Block** (committeter Text, den
      der Slice abarbeitet): Gegenstand ist der Satz, dass eine im Run nicht anwendbare
      Routing-Regel den Run mit der im Verdikt entschiedenen Klasse beendet, ohne
      Change und ohne den Erfassungspfad zu berühren; das Benutzerhandbuch (Abschnitt
      „Bestand als Backfill überführen") bleibt unberührt — Adresse:
      `slice-routing-betriebsdoku` §2 (Backfill-Bestand trägt das Label des
      Regelstands zum Run; Altbestand-Neuerzeugung über einen neuen Run).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der
      Closure der Welle [welle-routing](../welle-routing.md) (die Roadmap führt sie
      unter *Offene Wellen*, das Ereignis kann eintreten).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `internal/application/usecase/backfill/service.go` | update | Regelstand der Routing-Regeln lesen (neben `transformationRules`), Auswertung je Zeile über die Domänen-Funktion, Label in die Change des Blocks, Vergleich des Regelstands vor jedem Block (Fail-closed), `classifyError` für die im Verdikt V2 entschiedene Klasse. |
| `internal/application/port/outbound/` (Regelstand-Port der Routing-Regeln aus `slice-routing-antragsweg`) | lesen | der Run nutzt den Port, den der Antragsweg anlegt; ein tabellenbezogener Lesezugriff ist nicht Teil (Risiko Lesekosten, §6). |
| `internal/adapters/driven/postgressnapshot/` | lesen / bei Bedarf update | der Snapshot-Reader liefert Spaltenwerte als Text (Ergebnisformat `ADR-0115`); die Auswertung liest dieselbe Textform wie der WAL-Pfad. |
| `internal/adapters/driven/postgresstorage/` (Backfill-Writer) | lesen | der Writer schreibt `route_target` seit `slice-routing-kern-label` (`InsertBackfillChange`); hier nur ein Test, dass ein gesetztes Feld ankommt. |
| `internal/bootstrap/wiring.go` | update | der Backfill-Dienst erhält den Regelstand-Port; Bindung des Worker-Pfads. |
| `internal/application/usecase/backfill/routing_test.go` (neu), `service_test.go` | neu / update | Happy/Boundary/Negative nach [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) und [`LH-FA-CAP-009`](../../../../spec/lastenheft.md); Vertragstest WAL-Label gleich Backfill-Label. |
| `internal/adapters/driven/postgresstorage/*_test.go`, `internal/adapters/driven/postgressnapshot/*_test.go` | update | Label überlebt den Backfill-Insert; Bild-Parität des Typ-Satzes bleibt grün (Bedingung liest Textwerte). |
| `spec/pflichtenheft.md` (`LH-FA-CAP-009.a`, `SPEC-008`) | update, soweit offen | Übergabe-Block in §2. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „ein Backfill-Run
liefert Changes ohne Ziel", „der Regelstand des Runs besteht aus den
Transformationsregeln"; Parent ist `30fd6cb5`; der Implementer ergänzt die
`diff`-Zeilen und trägt Gefundenes und Nichtgefundenes ein):**

```suchlauf
30fd6cb5 46 -n -E 'Transformation' -- internal/application/usecase/backfill internal/bootstrap ':!*_test.go'
30fd6cb5 18 -n -E 'ErrTransformationColumnMissing|ErrTransformationTargetCollides' -- internal ':!*_test.go'
30fd6cb5 5 -n -E 'ErrTransformationStateChanged' -- internal ':!*_test.go'
30fd6cb5 5 -n -i -E 'fail-closed' -- spec docs/user
30fd6cb5 25 -n -E 'Regelstand' -- spec docs/user
```

| Träger | Messung am Parent (`30fd6cb5`, gemessen am 2026-10-01) | Behandlung und Befund am Diff |
|---|---|---|
| Stellen, an denen der Run den Regelstand der Transformationen führt | Zeile 1: 46 Nicht-Test-Zeilen in Backfill-Dienst und Composition Root | jede Stelle lesen: ist sie Aufzählung der Regel-Quellen des Runs (Start, Block, Abschluss, Klassifikation), die um die Routing-Regeln wachsen muss — Befund am Diff: einzutragen |
| Klassen-Abbildung des Runs | Zeilen 2 und 3: 18 bzw. 5 Zeilen | die Abbildung der Fehler auf `schema`/`configuration` folgt dem Verdikt V2; Befund: einzutragen |
| Spec-Aussagen zum Run | Zeilen 4 und 5: 5 bzw. 25 Zeilen („Fail-closed", „Regelstand") | Aussagen über „der Regelstand" im Backfill nennen heute nur Transformationen und Ausschluss; Nachzug nach dem Übergabe-Block; Befund: einzutragen |
| Handbuch Abschnitt „Bestand als Backfill überführen" | liegt außerhalb dieses Suchraums | gemeldet an `slice-routing-betriebsdoku` (§2 dort nennt den Gegenstand), nicht mitgeändert |

## 4. Trigger

**Start** (`next` → `in-progress`): `slice-routing-antragsweg` liegt in `done/`,
**das Architect-Kurzverdikt zu V2** (Welle §5) liegt vor (eine Datei unter
`docs/reviews/` oder eine Folge-ADR; ohne sie legte der Implementer eine
Entscheidung aus) und kein anderer Slice liegt in `in-progress/` (WIP-Limit 1).
Grund der ersten Bedingung: der Run liest den Regelstand über den Port des
Antragswegs, und dieser Slice folgt **unmittelbar** auf den Antragsweg, damit das
Fenster, in dem Regeln setzbar sind, ein Backfill-Run aber `NULL` liefert, ein Slice
lang ist (Welle §4 Reihenfolge).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht erwartet — die
  Auswertung ist eine Zeile je Zeile des Bestands und die Fail-closed-Prüfung folgt
  dem Muster der Transformationsregeln; sprengte die Klassen-Abbildung den Umfang,
  trennt sich die Nichtanwendbarkeit im Run (Liefer-Punkt 3) ab.
- `in-progress` → `open` (blockiert): das Verdikt V2 fordert eine achte Fehlerklasse
  oder ein anderes Run-Ende als `failed` — dann eine Folge-ADR, kein Weiterbau.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code
ungefiltert), `make test-store` und `make test-replication` real grün, Suchlauf-Block
nachgemessen, Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Lesekosten je Block.** Der Regelstand-Port der Transformationen liest je Aufruf
  die `applied`-Zeilen aller Tabellen einer Quelle, der Run liest sie je Block; die
  Messung des Architect-Verdikts
  [`architect-verdict-welle-transformationen-offene-fragen`](../../../reviews/architect-verdict-welle-transformationen-offene-fragen.md)
  (§3, synthetisch: bis rund 10 000 Zeilen der Queue etwa 8 % einer Blockdauer, bei
  100 000 Zeilen etwa 70 %, **übernommen**) gilt für **einen** Port; der zweite Port
  des Routings verdoppelt die Lesungen je Block (*hergeleitet*, nicht gemessen). —
  **Ausgang:** bei der Closure einzutragen (Trigger des Verdikts: mehr als 10 000
  Zeilen in der Queue einer Quelle; der Slice nennt die Verdopplung dort als
  Konsequenz oder liest beide Zustände in einem Aufruf).
- **Zwei Regelstände, ein Fail-closed.** Transformations- und Routing-Regeln ändern
  sich unabhängig voneinander; die Prüfung vergleicht beide. — **Ausgang:** bei der
  Closure einzutragen (Testfall: nur der Routing-Stand wechselt).
- **Nichtanwendbarkeit im Run ohne Verdikt (V2).** Der Plan nimmt die Klasse `schema`
  als Erwartung; sie ist keine Entscheidung. — **Ausgang:** bei der Closure
  einzutragen (Verdikt-Verweis).
- **Bild-Parität.** Die Bedingung liest den Textwert wie das Row Image ihn trägt
  (`ADR-0115`); der Backfill-Pfad liest über das Text-Ergebnisformat des Snapshots,
  der WAL-Pfad über `pgoutput`-Text — gleiche Textform für den Typ-Satz ist
  *erwartet*, belegt für das Row Image, nicht für die Auswertung. — **Ausgang:** bei
  der Closure einzutragen (Typ-Satz-Test: dieselbe Zeile, beide Pfade, dasselbe Ziel).
- **Persistenz des Labels im Backfill-Insert.** Wird in `slice-routing-kern-label`
  gebaut; wenn der Writer das Feld dort nicht schreibt, fiele das Label hier
  unbemerkt weg. — **Ausgang:** bei der Closure einzutragen (Store-Test mit gesetztem
  Feld).
- **Coverage-Messgegenstand.** Der Backfill-Dienst liegt in der netzlos gemessenen
  Fläche, der Snapshot im DB-Gegenstand. — **Ausgang:** bei der Closure einzutragen
  (`make coverage-gate`, `make test-replication`).

## 7. Closure-Notiz

Wird bei der Closure gefüllt (vor dem `git mv` nach `done/`).

- **Was hat funktioniert:** —
- **Was ging anders als geplant:** —
- **Steering-Loop-Eintrag:** —
- **Beobachtungs-Register (`../observations/`):** —
- **Folge-Slices:** —
- **Risiken aus §6:** —

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit den
Pfaden `internal/application/usecase/backfill/`, `internal/adapters/driven/` und
`internal/bootstrap/` — eine Sub-Area.

**Vorgelagert — offene Beobachtungen sichten:** das Register wurde durchgegangen
(Zähler = Zahl der `evidence/`-Dateien, ausgezählt am 2026-10-01); die Treffer
stehen in §6 und in der Eröffnungs-Sichtung der Welle
[welle-routing](../welle-routing.md) §6:
`BEO-PGC/run-fehlerklasse-schema-im-transformations-backfill` (verkörpert, 1×),
`BEO-PGC/run-fehlertext-traegt-klasse-doppelt` (verkörpert, 1×),
`BEO-PGC/backfill-schema-version-hinter-snapshot-spalten` (verkörpert, 1×),
`BEO-PGC/db-gegenstand-enthaelt-netzlos-geprueften-code` (verkörpert, 3×),
`BEO-PGC/negativtest-ohne-bindung-an-seine-eingabe` (verkörpert, 21×),
`BEO-PGC/implementierung-weicht-von-adr-wortlaut-ab` (offen, 2× — V2 ist eine
benannte Frage statt einer stillen Abweichung). Keiner der offenen Einträge erreicht
mit diesem Slice 3×.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

