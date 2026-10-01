# Slice routing-spec-nachzug: Spec-Nachzug — das Pflichtenheft trägt den beschlossenen Routing-Stand

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
Haupt-Bezug), [`LH-FA-CFG-008.a`](../../../../spec/pflichtenheft.md) (die offene,
ADR-pflichtige Frage „Routingform"),
[`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
Folgepflicht 1 (Spec-Nachzug — Träger dieses Slice),
[`LH-FA-CFG-007`](../../../../spec/lastenheft.md) (Transformationen — nur
abgegrenzt).

**Berührte Spec-Stellen:**
[`LH-FA-CFG-008.a`](../../../../spec/pflichtenheft.md),
[`SPEC-019`](../../../../spec/pflichtenheft.md) (Antrags-Datensatz — zwei weitere
Antragsarten), eine **neue** Kennung für die Routing-Regelform (neben
[`SPEC-030`](../../../../spec/pflichtenheft.md)),
[`SPEC-001`](../../../../spec/pflichtenheft.md) und
[`SPEC-002`](../../../../spec/pflichtenheft.md) (Spalte `route_target`),
[`SPEC-008`](../../../../spec/pflichtenheft.md) (Zeile `schema`),
[`SPEC-020`](../../../../spec/pflichtenheft.md),
[`SPEC-021`](../../../../spec/pflichtenheft.md),
[`SPEC-022`](../../../../spec/pflichtenheft.md) (Parameter `target`),
[`SPEC-024`](../../../../spec/pflichtenheft.md) (Zusatz-Subjekt), gegebenenfalls
`spec/architecture.md` (Antragsart-Tabelle, Capture-Sequenz). Gelesen, nicht
geändert: [`SPEC-017`](../../../../spec/pflichtenheft.md) (Wecksignal). Der
Verweis zeigt **aufwärts**: die Spec nennt diesen Slice nie.

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-routing](../welle-routing.md).
**Datum:** 2026-10-01.

---

## 1. Ziel und Abgrenzung

**Ziel:** Das Pflichtenheft (Rang 2) trägt den mit
[`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
beschlossenen Routing-Stand als Technik-Festlegung, **bevor** der erste umsetzende
Slice startet (die Modus-Deklaration in `harness/conventions.md` ist Greenfield:
die Doku führt). Umfang:

- (a) [`LH-FA-CFG-008.a`](../../../../spec/pflichtenheft.md): die Überschrift
  verliert das Wort „offen"; der Text legt Zielmodell (benannter Kanal je Quelle,
  Namensalphabet), Konfigurationsmechanismus (zwei Antragsarten), Ausdrucksform
  (`target`, `order`, optional `when` mit Gleichheit), Auflösung (erster Treffer in
  aufsteigender `order`, R1–R6 statisch ausgeschlossen), Wirkort (vor der
  Persistierung, Label an der Change) und Verhältnis zu Transformation,
  Spaltenausschluss und Backfill als Zusagen fest — in Zukunfts-Form, bis die
  Umsetzungs-Slices sie belegen;
- (b) [`SPEC-019`](../../../../spec/pflichtenheft.md): die Antragsarten
  `set_route`/`remove_route` (die `request_kind`-Menge wächst um zwei Werte), die
  Bedeutung von `rule_name`/`rule_spec` für sie (die Spalten bestehen), die
  zweite Bedeutung von `applied` (dauerhafter Regelstand, Ordnung `requested_at`,
  dann `administration_request_id`) und die `failed`-Fehlertexte je Verletzung
  (R1–R6, unbekannter Schlüssel, ungültiger Zielname);
- (c) die Routing-Regelform: `rule_spec` als JSON-Objekt (`target`, `order`,
  `when.column`/`when.equals`), die Validierung (Alphabet des Zielnamens, `order`
  als positive ganze Zahl, `equals` als Zeichenkette), die Bildbasis der Bedingung
  (Neu-Bild bei INSERT und UPDATE, Alt-Bild bei DELETE; Quellwert vor jeder
  Transformation) und die Randfälle, die der ADR-Text offen lässt — entweder in
  [`SPEC-019`](../../../../spec/pflichtenheft.md) oder unter der nächsten freien
  Kennung (die Entscheidung trifft der Slice; die Kennung ist die nächste freie am
  Parent-Stand, dort `SPEC-031` die letzte);
- (d) [`SPEC-001`](../../../../spec/pflichtenheft.md) und
  [`SPEC-002`](../../../../spec/pflichtenheft.md): `cdc.change` und `cdc.changes`
  tragen `route_target` (nullable, kein DEFAULT, kein CHECK, letzte Spalte der
  View; die Schlüsselmenge der Row Images bleibt unberührt, `NULL` heißt „kein
  Ziel");
- (e) [`SPEC-008`](../../../../spec/pflichtenheft.md): die Zeile `schema` nennt die
  Nichtanwendbarkeit einer Routing-Regel (Spalte fehlt in der Relation der Change)
  im Erfassungspfad samt Abhilfe-Weg (`cdc.remove_route`, Neustart — als Zusage,
  bis die Umsetzung sie belegt; Vorab-Bedingung V3 der Welle); der Run-Satz folgt
  erst mit dem Verdikt V2 (siehe NICHT unten);
- (f) [`SPEC-020`](../../../../spec/pflichtenheft.md),
  [`SPEC-021`](../../../../spec/pflichtenheft.md),
  [`SPEC-022`](../../../../spec/pflichtenheft.md): der optionale Parameter
  `target` (leer = kein Filter, Konjunktion mit `schema`/`table`; HTTP: unbekannte
  Parameter bleiben `400`); [`SPEC-024`](../../../../spec/pflichtenheft.md): das
  Zusatz-Subjekt `cdc.route.<source_id>.<ziel>` (Payload wie `cdc.stream…`, nur
  für Changes mit Ziel);
- (g) `spec/architecture.md`: die Antragsart-Tabelle (heute sieben Arten) trägt die
  zwei weiteren; die Capture-Sequenz sagt, dass eine Change ihr Ziel vor der
  Persistierung bestimmt — **ohne** ADR-, Slice- oder Wellen-Bezug
  ([`AGENTS.md`](../../../../AGENTS.md) §3.4). Ob die Sicht Change-Felder
  aufzählt, die `route_target` nachziehen müssten, klärt der Suchlauf in §3 (am
  Parent `30fd6cb5`: `origin`, das Schwesterfeld, kommt in `spec/architecture.md`
  0-mal vor, gemessen).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Lastenheft** — bleibt unverändert; das Pflichtenheft präzisiert, erweitert
  nie (Kopf von `spec/pflichtenheft.md`).
- **Zusagen ohne Verdikt** — die Vorab-Bedingungen V1 und V2 der Welle (Parameter
  am gRPC-RPC `ReadChanges`, Nichtanwendbarkeit im Backfill-Run) sind Lücken im Text
  der ADR. Liegt das Verdikt zum Start dieses Slice vor, gehören die Spec-Zeilen in
  (b), (e), (f); liegt es nicht vor, schreibt der Slice keine Zusage, die kein
  Verdikt trägt, und die Zeilen folgen im Slice, der das Verdikt braucht
  (`lesewege`, `backfill-pfad`; dort als Übergabe-Block in §2 geführt). Die
  Aufzählung in (f) nennt dann `SPEC-020` nur für den `StreamChangesRequest`.
- **Benutzerhandbuch** — beschreibt Betreiber-Oberfläche und darf keine Funktion
  nennen, die es noch nicht gibt; `slice-routing-betriebsdoku` trägt sie (Welle §4
  Abweichung 1).
- **Code, Schema, Skripte** — Umsetzung gehört den Folge-Slices.
- **Routing-Operatoren über Gleichheit hinaus, Standardziel, mehrere Ziele** — mit
  der ADR ausgeschlossen (Entscheidungen 3 und 5); die Spec sagt sie nicht zu.

## 2. Definition of Done

- [ ] [`LH-FA-CFG-008.a`](../../../../spec/pflichtenheft.md) ist beantwortet: die
      Überschrift trägt kein „offen" mehr, der Text nennt Zielmodell,
      Mechanismus, Ausdrucksform, Auflösung, Wirkort und das Verhältnis zu
      Transformation, Ausschluss und Backfill als Zusagen in Zukunfts-Form (was die
      Umsetzung liefern **muss**; bis zu den Umsetzungs-Slices unbelegt, nicht als
      geprüfte Aussage, [`AGENTS.md`](../../../../AGENTS.md) §3.12). Die Abhilfe
      nach einer nicht anwendbaren Regel steht als Zusage, nicht als Tatsache
      ([`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
      führt sie als Erwartung). *Zu belegen durch:* Lesen des Abschnitts und
      `make docs-check`.
- [ ] Die Datenstrukturen und Schnittstellen stehen:
      [`SPEC-019`](../../../../spec/pflichtenheft.md) (zwei Antragsarten, Menge,
      `applied`-Bedeutung, Fehlertexte R1–R6), die Routing-Regelform samt
      Randfällen und Bildbasis, [`SPEC-001`](../../../../spec/pflichtenheft.md) und
      [`SPEC-002`](../../../../spec/pflichtenheft.md) (`route_target`), die Zeile
      `schema` in [`SPEC-008`](../../../../spec/pflichtenheft.md), der Parameter
      `target` in [`SPEC-020`](../../../../spec/pflichtenheft.md),
      [`SPEC-021`](../../../../spec/pflichtenheft.md) und
      [`SPEC-022`](../../../../spec/pflichtenheft.md), das Zusatz-Subjekt in
      [`SPEC-024`](../../../../spec/pflichtenheft.md); eine neu vergebene Kennung
      ist die nächste freie (*zu belegen durch* `git grep -h -o 'SPEC-0[0-9][0-9]'
      <Stand> -- spec/pflichtenheft.md | sort -u | tail -1` am Parent-Stand und am
      Diff-Stand); §7 Historie trägt je Änderung eine Zeile ohne ADR-/Slice-Bezug.
- [ ] `spec/architecture.md` trägt die Antragsart-Tabelle mit den zwei weiteren
      Arten und die Aussage zum Ziel vor der Persistierung, ohne ADR-/Slice-/Wellen-Bezug.
      *Zu belegen durch:* `make docs-check` (`matrix`-Modul) und der Suchlauf in §3.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-routing-spec-nachzug.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: entfällt als eigener Punkt — der Slice **ist** das Doku-Update
      der Spec; das Benutzerhandbuch bleibt unberührt (§1).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
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
| `spec/pflichtenheft.md` §1 ([`LH-FA-CFG-008.a`](../../../../spec/pflichtenheft.md)) | update | Überschrift und Text auf den beantworteten Stand; Zusagen-Form. |
| `spec/pflichtenheft.md` §2 ([`SPEC-019`](../../../../spec/pflichtenheft.md)) | update | zwei Antragsarten, `request_kind`-Menge (Stand `30fd6cb5`: sieben Werte, nach dem Slice neun), Zählwörter der Spalten-Zeilen (`column_name`, `rule_name`, `rule_spec`), zweite Bedeutung von `applied`, Fehlertexte R1–R6 mit Prüfreihenfolge. |
| `spec/pflichtenheft.md` §2 (neue Kennung, nächste freie am Parent-Stand: `SPEC-032`) | neu | die Routing-Regelform als eigene Datenstruktur (Schlüssel, Bildbasis, Randfälle), Formvorbild `SPEC-030`; die Aufteilung gegenüber `SPEC-019` entscheidet der Slice (Festlegungen unten). |
| `spec/pflichtenheft.md` §2 ([`SPEC-001`](../../../../spec/pflichtenheft.md), [`SPEC-002`](../../../../spec/pflichtenheft.md)) | update | Spalte `route_target` an `cdc.change` und `cdc.changes` (letzte Spalte der View); Format wie die Zeile zu `origin`. |
| `spec/pflichtenheft.md` §2 ([`SPEC-020`](../../../../spec/pflichtenheft.md), [`SPEC-021`](../../../../spec/pflichtenheft.md), [`SPEC-022`](../../../../spec/pflichtenheft.md), [`SPEC-024`](../../../../spec/pflichtenheft.md)) | update | Parameter `target` bzw. Zusatz-Subjekt; die Zählwörter „zehn Felder" und „dreizehn Felder" bleiben (das Label ist nicht Teil der Nachrichten). |
| `spec/pflichtenheft.md` §4 ([`SPEC-008`](../../../../spec/pflichtenheft.md) Zeile `schema`) | update | Nichtanwendbarkeit einer Routing-Regel im Erfassungspfad; Absatzstruktur wie „Nicht anwendbare Regel (Klasse `schema`)" der Transformationen (Zellenlänge des `structure`-Moduls). |
| `spec/pflichtenheft.md` §7 Historie | update | je Änderung eine Zeile ohne ADR-/Slice-Bezug. |
| `spec/architecture.md` (Antragsart-Tabelle bei „sieben Arten", Capture-Sequenz) | update | die zwei weiteren Antragsarten; das Ziel einer Change steht vor der Persistierung. |

**Festzulegen im Slice** (Lücken im Text von
[`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md), die
die Spec schließen muss; Plan-Stand: offen, jede Festlegung bekommt im Slice einen
Bestands-Beleg wie in `slice-transformationen-spec-nachzug`):

- Vergleich von `when.column` (zeichengenau gegen den Katalog wie bei
  `SPEC-030`? ohne U+0000?) und von `when.equals` (Länge, Zeichen, leere Zeichenkette
  zulässig?).
- Obergrenze und Wertebereich von `order` („positive ganze Zahl": größter Wert?).
- Prüfreihenfolge der Fehlertexte (Form vor Konfliktfreiheit) und Wortlaut je
  Verletzung R1–R6; Verhalten von `remove_route` (nur Namens-Zeile und R6).
- Verhalten eines Antrags gegen eine Tabelle ohne laufende Bindung (Muster der
  Spalten- und Transformations-Antragsarten).
- Verhalten der Lesewege bei einem Wert von `target`, der das Alphabet verletzt oder
  nie vergeben wurde (`400`/`InvalidArgument` oder leere Antwort) — die Lesewege-
  Parameter von `SPEC-020`/`-021`/`-022` und das Subjekt von `SPEC-024` sagen es
  einheitlich.
- Aufteilung `SPEC-019` / neue Kennung; „führende Stelle je Sachverhalt" wie in
  `slice-transformationen-spec-nachzug` (Anwendbarkeit und Reihenfolge, Fehlertexte,
  Abhilfe — je an genau einer Stelle).

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „`LH-FA-CFG-008.a` ist
beantwortet; die Antragsarten-Menge trägt zwei weitere Werte; `cdc.changes` trägt
eine Spalte mehr; Stream-Anfragen tragen einen Parameter mehr"; Parent ist
`30fd6cb5`; der Implementer trägt Gefundenes und Nichtgefundenes je Zeile in die
Tabelle ein und ergänzt die `diff`-Zeilen im Block, Suchbefehle in Codeblöcken, nie
in Zellen):**

```suchlauf
30fd6cb5 5 -n -E 'Routingform|LH-FA-CFG-008\.a' -- spec docs/user harness README.md
30fd6cb5 6 -n -E 'sieben (Arten|Werte|Antragsarten)|fünf übrigen|sechs übrigen' -- spec docs/user harness
30fd6cb5 22 -n -E 'zehn Felder|dreizehn Felder' -- spec docs/user
30fd6cb5 11 -n -w origin -- spec
30fd6cb5 0 -n -w origin -- spec/architecture.md
30fd6cb5 0 -n -E 'Routing|Zustellziel' -- spec/architecture.md
```

| Träger | Messung am Parent (`30fd6cb5`, gemessen am 2026-10-01) | Behandlung und Befund am Diff |
|---|---|---|
| Anker und Überschrift „Routingform" | Zeile 1 des Blocks: 5 Zeilen, alle in `spec/pflichtenheft.md` (der Verweis in `LH-FA-CFG-007.a` `:347`, die Überschrift `:349`, drei Historie-Zeilen `:1265`, `:1287`, `:1298`); `docs/user`, `harness` und `README.md`: 0 | Verweis und Überschrift ziehen, Historie-Zeilen bleiben (Protokoll); Anker prüft `make docs-check` (Modul `anchors`). Die Zeile `make doc-trace` in `harness/README.md` nennt `LH-FA-CFG-008` ohne `.a` (nicht in diesem Suchraum) und bewegt sich erst mit `slice-routing-e2e` — Befund am Diff: einzutragen |
| Zählwörter zu Antragsarten und Werten | Zeile 2: 6 Zeilen (`spec/architecture.md:258` „sieben Arten", `spec/pflichtenheft.md:620`–`622`, Historie `:1288`, `harness/targets/schema-rollout.md:74` „die sieben SQL-Funktionen") | Spec-Zeilen ziehen; `harness/targets/schema-rollout.md` gemeldet an `slice-routing-antragsweg` (es beschreibt den Stand des Rollouts); die Historie-Zeile bleibt (Protokoll) — Befund am Diff: einzutragen |
| Nachrichtenschema-Zählwörter | Zeile 3: 22 Zeilen — erwartet unverändert am Diff (das Label ist nicht Teil der Nachrichten) | bei Abweichung ist der Zug falsch, nicht der Zähler — Befund am Diff: einzutragen |
| Schwesterfeld `origin` als Muster für Aufzählungen von Change-Feldern | Zeilen 4 und 5: 11 Zeilen in der Spec, davon 0 in der Architektur-Sicht | die Sicht zählt keine Change-Felder auf; ob sie eine andere Aufzählung trägt, die `route_target` brauchte, klärt das Lesen der Sequenz „Capture" am Diff — Befund: einzutragen |
| Routing in der Architektur-Sicht | Zeile 6: 0 Zeilen | die Sicht sagt zum Routing heute nichts; der Zug fügt höchstens einen Satz in der Capture-Sequenz und die zwei Tabellenzeilen ein |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn die Welle [welle-routing](../welle-routing.md)
eröffnet ist und kein anderer Slice in `in-progress/` liegt (WIP-Limit 1). Dieser
Slice kommt vor jedem übrigen der Welle: jeder Folge-Slice liest die hier
festgelegten Zusagen. Wünschenswert, nicht Bedingung: die Verdikte V1 und V2 (Welle
§5) liegen vor; dann zieht der Slice ihre Zeilen mit (§1, NICHT-Punkt „Zusagen ohne
Verdikt").

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht erwartet — ein
  reiner Doku-Zug über zwei Dateien; sprengte er den Umfang, wäre die
  Architektur-Sicht (Punkt g) der abtrennbare Teil.
- `in-progress` → `open` (blockiert): falls das Übertragen eine Aussage von
  [`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
  als widersprüchlich zum Bestand der Spec zeigt (etwa zum Backfill-Stand von
  `LH-FA-CAP-009.a` oder zur Zeile `schema`); die Frage geht an den Architect.

## 5. Closure-Trigger

DoD vollständig, Review ohne offenes HIGH/MEDIUM, `make gates` grün (Exit-Code
ungefiltert), Suchlauf-Block nachgemessen (`make suchlauf-nachmessen` Exit 0),
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Spec-Wiedergabe in eigenen Worten driftet von der ADR ab.** Die Zusagen stehen
  in Zukunfts-Form; ein Satz, der als Tatsache formuliert ist, bevor ein
  Umsetzungs-Slice ihn belegt, ist eine ungeprüfte Aussage
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B). — **Ausgang:** bei der
  Closure einzutragen (Leser: Reviewer, Verifier).
- **Zählwörter und Aufzählungen bleiben stehen** (`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`,
  verkörpert, 15×): „sieben Arten", „die fünf/sechs übrigen Antragsarten". —
  **Ausgang:** bei der Closure einzutragen (Suchlauf-Block in §3).
- **Doppelquelle Handbuch gegen Pflichtenheft**
  (`BEO-PGC/zwei-quellen-drift-handbuch-gegen-pflichtenheft`, offen, 3×): dieselben
  Fakten stehen später im Handbuch (`betriebsdoku`). — **Ausgang:** bei der Closure
  einzutragen; Adresse: `slice-routing-betriebsdoku` §1 (Handbuch verweist auf die
  führende Stelle der Spec je Sachverhalt, statt Fehlertexte zu wiederholen).
- **Doc-Gate-Regeln der Spec.** `matrix` verbietet Verweise Spec → ADR/Slice,
  `ids` verlangt Links auf nackte Kennungen, `structure` misst Zellenlängen (die
  Transformations-Fassung färbte `section-cell-oversized`). — **Ausgang:** bei der
  Closure einzutragen (`make docs-check` Exit 0 am Diff-Stand).
- **Zusagen ohne Verdikt** (V1, V2 der Welle). — **Ausgang:** bei der Closure
  einzutragen (Spec-Zeilen gezogen oder mit Adresse an `lesewege` bzw.
  `backfill-pfad` übergeben, dort als committeter Text in §2).

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
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit dem
Pfad `spec/` — eine Sub-Area, keine zu grobe Zusammenfassung.

**Vorgelagert — offene Beobachtungen sichten:** das Register wurde durchgegangen
(Zähler = Zahl der `evidence/`-Dateien, ausgezählt am 2026-10-01); die Treffer
dieses Slice stehen in §6 und in der Eröffnungs-Sichtung der Welle
[welle-routing](../welle-routing.md) §6:
`BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (verkörpert, 15×),
`BEO-PGC/zwei-quellen-drift-handbuch-gegen-pflichtenheft` (offen, 3×),
`BEO-PGC/zitat-nennt-die-falsche-stelle` (verkörpert, 9×). Kein offener Eintrag
erreicht mit diesem Slice 3×, den die Welle nicht schon führt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

