# Slice routing-spec-nachzug: Spec-Nachzug — das Pflichtenheft trägt den beschlossenen Routing-Stand

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-routing](welle-routing.md).

**Bezug:** [`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (Routing —
Haupt-Bezug), [`LH-FA-CFG-008.a`](../../../../spec/pflichtenheft.md) (die offene,
ADR-pflichtige Frage „Routingform"),
[`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
Folgepflicht 1 (Spec-Nachzug — Träger dieses Slice),
[`LH-FA-CFG-007`](../../../../spec/lastenheft.md) (Transformationen — nur
abgegrenzt),
[`ADR-0138`](../../adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)
Folgepflicht (Spec-Zeilen `SPEC-031`, `SPEC-008` — ebenfalls Träger dieses
Slice).

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
[`SPEC-024`](../../../../spec/pflichtenheft.md) (Zusatz-Subjekt),
[`SPEC-031`](../../../../spec/pflichtenheft.md) (Zeile `ReadChanges`: Request um
`target` ergänzt),
gegebenenfalls `spec/architecture.md` (Antragsart-Tabelle, Capture-Sequenz). Gelesen,
nicht geändert: [`SPEC-017`](../../../../spec/pflichtenheft.md) (Wecksignal) und
[`SPEC-029`](../../../../spec/pflichtenheft.md) (die Run-Fehlerklasse `schema` gilt
auch für eine Routing-Regel; die Festlegung steht in `SPEC-008` und
`LH-FA-CAP-009.a`, `SPEC-029` verweist generisch auf `SPEC-008`). Der
Verweis zeigt **aufwärts**: die Spec nennt diesen Slice nie.

**Verantwortlich:** — (gesetzt beim Übergang `open` → `next`).

**Autor:** Planner-Agent, Welle-Eröffnung [welle-routing](welle-routing.md).
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
  bis die Umsetzung sie belegt; Vorab-Bedingung V3 der Welle, offen laut
  `ADR-0138`); der Run-Satz (V2, entschieden mit `ADR-0138` Festlegung 2) steht in
  [`SPEC-008`](../../../../spec/pflichtenheft.md) (`SPEC-029` bleibt unberührt): eine im Run nicht anwendbare
  Routing-Regel (`when.column` fehlt in den Spalten des Snapshots) beendet den Run
  `failed` mit Klasse `schema`, einmal je Run vor der Schreibtransaktion, run-lokal,
  ohne Change (als Zusage in Zukunfts-Form), samt der Zeile zu `LH-FA-CAP-009.a`;
- (f) [`SPEC-020`](../../../../spec/pflichtenheft.md),
  [`SPEC-021`](../../../../spec/pflichtenheft.md),
  [`SPEC-022`](../../../../spec/pflichtenheft.md): der optionale Parameter
  `target` (leer = kein Filter, Konjunktion mit `schema`/`table`; HTTP: unbekannte
  Parameter bleiben `400`); [`SPEC-031`](../../../../spec/pflichtenheft.md): die
  Zeile `ReadChanges` trägt `target` als siebtes Request-Feld (`string target = 7`,
  leer = kein Filter, Konjunktion mit `schema`/`table`; `ChangeRecord` bleibt
  unverändert); [`SPEC-024`](../../../../spec/pflichtenheft.md): das
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
- **Zusagen ohne Entscheidung** — V1 und V2 der Welle sind mit
  [`ADR-0138`](../../adr/0138-routing-filter-target-readchanges-und-nichtanwendbarkeit-im-backfill-run.md)
  `Accepted`; ihre Spec-Zeilen (`SPEC-031` `ReadChanges`, `SPEC-029`/`SPEC-008` Run)
  trägt dieser Slice in (e) und (f). Es gibt keinen Übergabe-Block an `lesewege` und
  `backfill-pfad`. Offen bleibt V3 (Erreichbarkeit der Nichtanwendbarkeit): die
  Abhilfe-Aussage steht nur als Zusage (e).
- **Benutzerhandbuch** — beschreibt Betreiber-Oberfläche und darf keine Funktion
  nennen, die es noch nicht gibt; `slice-routing-betriebsdoku` trägt sie (Welle §4
  Abweichung 1).
- **Code, Schema, Skripte** — Umsetzung gehört den Folge-Slices.
- **Routing-Operatoren über Gleichheit hinaus, Standardziel, mehrere Ziele** — mit
  der ADR ausgeschlossen (Entscheidungen 3 und 5); die Spec sagt sie nicht zu.

## 2. Definition of Done

- [x] [`LH-FA-CFG-008.a`](../../../../spec/pflichtenheft.md) ist beantwortet: die
      Überschrift trägt kein „offen" mehr, der Text nennt Zielmodell,
      Mechanismus, Ausdrucksform, Auflösung, Wirkort und das Verhältnis zu
      Transformation, Ausschluss und Backfill als Zusagen in Zukunfts-Form (was die
      Umsetzung liefern **muss**; bis zu den Umsetzungs-Slices unbelegt, nicht als
      geprüfte Aussage, [`AGENTS.md`](../../../../AGENTS.md) §3.12). Die Abhilfe
      nach einer nicht anwendbaren Regel steht als Zusage, nicht als Tatsache
      ([`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md)
      führt sie als Erwartung). *Zu belegen durch:* Lesen des Abschnitts und
      `make docs-check`. *Beleg (Verifier):* Verifikations-Report §2 Zeile 1
      (getragen), `make docs-check` Exit 0.
- [x] Die Datenstrukturen und Schnittstellen stehen:
      [`SPEC-019`](../../../../spec/pflichtenheft.md) (zwei Antragsarten, Menge,
      `applied`-Bedeutung, Fehlertexte R1–R6), die Routing-Regelform samt
      Randfällen und Bildbasis, [`SPEC-001`](../../../../spec/pflichtenheft.md) und
      [`SPEC-002`](../../../../spec/pflichtenheft.md) (`route_target`), die Zeile
      `schema` in [`SPEC-008`](../../../../spec/pflichtenheft.md) (Erfassungspfad
      und Backfill-Run), der Run-Satz in `LH-FA-CAP-009.a`, der Parameter
      `target` in [`SPEC-020`](../../../../spec/pflichtenheft.md),
      [`SPEC-021`](../../../../spec/pflichtenheft.md),
      [`SPEC-022`](../../../../spec/pflichtenheft.md) und als Feld `target` der Zeile
      `ReadChanges` in [`SPEC-031`](../../../../spec/pflichtenheft.md), das
      Zusatz-Subjekt in [`SPEC-024`](../../../../spec/pflichtenheft.md); eine neu vergebene Kennung
      ist die nächste freie (*zu belegen durch* `git grep -h -o 'SPEC-0[0-9][0-9]'
      <Stand> -- spec/pflichtenheft.md | sort -u | tail -1` am Parent-Stand und am
      Diff-Stand); §7 Historie trägt je Änderung eine Zeile ohne ADR-/Slice-Bezug.
      *Beleg (Verifier):* Verifikations-Report §2 Zeile 2 „getragen mit zwei
      LOW“; die beiden LOW (V-1: Leeraussage und U+0000 in der Zeile Query-Parameter
      von `SPEC-021`, V-2: „im Wertebereich“ im Fehlertext von `SPEC-019`) sind
      mit dem Nachzug-Commit `ab4a83c7` geschlossen (Diff gelesen, `make docs-check`
      Exit 0).
- [x] `spec/architecture.md` trägt die Antragsart-Tabelle mit den zwei weiteren
      Arten und die Aussage zum Ziel vor der Persistierung, ohne ADR-/Slice-/Wellen-Bezug.
      *Zu belegen durch:* `make docs-check` (`matrix`-Modul) und der Suchlauf in §3. *Beleg (Verifier):*
      Verifikations-Report §2 Zeile 3 und §1 (Bezug-Grep im Spec-Diff: 0 Treffer).
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein offenes
      HIGH/MEDIUM (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8). *Beleg:* Report
      `docs/reviews/review-slice-routing-spec-nachzug.md` (0 HIGH, 2 MEDIUM F-1/F-2,
      2 LOW, 2 INFO). Die Fixrunde ist **nicht** von einem Reviewer erneut gelesen
      worden: der Verifier hat F-1 bis F-4 am Spec-Text gegengeprüft (Verifikations-Report
      §4, „am Text geprüft, nicht am Fixrunden-Bericht“; V-3 in §7 als Prozessbefund
      geführt) und die Nachzüge V-1/V-2 im Commit `ab4a83c7` gelesen — V-3 ist damit
      durch diese Gegenprüfung geschlossen, ein separates Re-Review hat es nicht
      gegeben.
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff);
      `make suchlauf-nachmessen PLAN=docs/plan/planning/<Verzeichnis>/slice-routing-spec-nachzug.md`
      endet mit Exit 0 ([`AGENTS.md`](../../../../AGENTS.md) §3.13). *Beleg:*
      Exit 0, „16 Zeilen stimmen“ (Verifikations-Report §1; im Closure-Lauf am
      Stand nach `ab4a83c7` wiederholt, gleiches Ergebnis).
- [x] Doku-Update: entfällt als eigener Punkt — der Slice **ist** das Doku-Update
      der Spec; das Benutzerhandbuch bleibt unberührt (§1). *Beleg:* Verifikations-Report
      §2 Zeile 7.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke). *Beleg:* §7.
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls eine
      Antwort und wird in §7 notiert. *Beleg:* §7, zwei weitere `evidence/`-Dateien.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
      offen). *Beleg:* §6.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen — von der
      Closure der Welle [welle-routing](welle-routing.md) (die Roadmap führt sie
      unter *Abgeschlossene Wellen*, die Closure erfolgte am 2026-10-02). *Beleg:*
      `welle-routing-results.md`, Abschnitt „Drei Paarungen“ (Closure 2026-10-02).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `spec/pflichtenheft.md` §1 ([`LH-FA-CFG-008.a`](../../../../spec/pflichtenheft.md)) | update | Überschrift und Text auf den beantworteten Stand; Zusagen-Form. |
| `spec/pflichtenheft.md` §2 ([`SPEC-019`](../../../../spec/pflichtenheft.md)) | update | zwei Antragsarten, `request_kind`-Menge (Stand `30fd6cb5`: sieben Werte, nach dem Slice neun), Zählwörter der Spalten-Zeilen (`column_name`, `rule_name`, `rule_spec`), zweite Bedeutung von `applied`, Fehlertexte R1–R6 mit Prüfreihenfolge. |
| `spec/pflichtenheft.md` §2 (neue Kennung, nächste freie am Parent-Stand: `SPEC-032`) | neu | die Routing-Regelform als eigene Datenstruktur (Schlüssel, Bildbasis, Randfälle), Formvorbild `SPEC-030`; die Aufteilung gegenüber `SPEC-019` entscheidet der Slice (Festlegungen unten). |
| `spec/pflichtenheft.md` §2 ([`SPEC-001`](../../../../spec/pflichtenheft.md), [`SPEC-002`](../../../../spec/pflichtenheft.md)) | update | Spalte `route_target` an `cdc.change` und `cdc.changes` (letzte Spalte der View); Format wie die Zeile zu `origin`. |
| `spec/pflichtenheft.md` §2 ([`SPEC-020`](../../../../spec/pflichtenheft.md), [`SPEC-021`](../../../../spec/pflichtenheft.md), [`SPEC-022`](../../../../spec/pflichtenheft.md), [`SPEC-024`](../../../../spec/pflichtenheft.md)) | update | Parameter `target` bzw. Zusatz-Subjekt; die Zählwörter „zehn Felder" und „dreizehn Felder" bleiben (das Label ist nicht Teil der Nachrichten). |
| `spec/pflichtenheft.md` §2 ([`SPEC-031`](../../../../spec/pflichtenheft.md) Zeile `ReadChanges`, `LH-FA-CAP-009.a`) | update | Request-Feld `target` (siebtes, Nummer 7); Run-Klasse `schema` für eine Routing-Regel (Quelle der Entscheidungen: `ADR-0138`, in der Spec ohne ADR-Bezug formuliert). |
| `spec/pflichtenheft.md` §4 ([`SPEC-008`](../../../../spec/pflichtenheft.md) Zeile `schema`) | update | Nichtanwendbarkeit einer Routing-Regel im Erfassungspfad und im Backfill-Run; Absatzstruktur wie „Nicht anwendbare Regel (Klasse `schema`)" der Transformationen (Zellenlänge des `structure`-Moduls). |
| `spec/pflichtenheft.md` §1 (`LH-FA-CAP-009.a`, Absatz „Sichtbarkeit und Fehler des Runs" und neuer Punkt „Ziel der Backfill-Changes") | update | über den Plan hinaus: der Punkt „Ziel der Backfill-Changes" trägt die Aussage aus der Entscheidung, dass Backfill-Changes dieselbe Bestimmung durchlaufen; die Fail-closed-Prüfung des Routing-Regelstands vor dem Commit (Stand des Runs nach dem Öffnen des Snapshots, Klasse `configuration`) und die Definition von „Regelstand zum Run" stehen im Absatz „Fail-closed vor dem Commit" und im Punkt „Ziel der Backfill-Changes" (Quelle: [`ADR-0139`](../../adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md), Fixrunde nach dem Review). |
| `spec/pflichtenheft.md` §2 (`SPEC-020`, `SPEC-021`, `SPEC-022`, `SPEC-031` Zeile `ReadChanges`) | update (Fixrunde) | Der Satz „leer, kein Fehler" für ein `target` außerhalb des Alphabets nennt U+0000 und die Prüfung im gemeinsamen Use Case; der SQL-Zugriff `cdc.changes` ist ausgenommen (Quelle: `ADR-0139`). Das NUL-Verhalten von PostgreSQL steht nicht als erprobt, die Spec macht den Satz von ihm unabhängig. |
| `spec/pflichtenheft.md` §2 (`SPEC-032`, Zeile `order`) | update (Fixrunde) | Wertebereich auf den Wortlaut der Entscheidung („positive ganze Zahl") zurückgenommen; die Obergrenze 2147483647 und der Ausschluss von Bruchteil und Exponent hatten keinen Träger in den Entscheidungen und keinen genannten Grund. Eine Obergrenze ist Sache der Umsetzung, die der Folge-Slice `slice-routing-antragsweg` mit ihrer Begründung vorlegt. |
| `spec/architecture.md`, `spec/pflichtenheft.md` (zwei Zeilen) | update (Fixrunde) | Umbruch nach Teilersetzung nachgezogen; Historie-Zeile ergänzt. |
| `docs/plan/adr/README.md` | update (Fixrunde) | Zeile `ADR-0139`. |
| Pläne `slice-routing-backfill-pfad`, `slice-routing-lesewege`, `welle-routing` | update (Fixrunde) | Festlegungen aus `ADR-0139` nachgezogen (Fail-closed-Tests, Use-Case-Prüfung); A-1/A-2 geschlossen, A-3 und V3 offen. |
| `spec/pflichtenheft.md` §7 Historie | update | je Änderung eine Zeile ohne ADR-/Slice-Bezug. |
| `spec/architecture.md` (Komponentenliste `ARC-001`, Antragsart-Tabelle bei „sieben Arten", Capture-Sequenz, Backfill-Absatz, Fehlermodell-Zeile) | update | die zwei weiteren Antragsarten; das Ziel einer Change steht vor der Persistierung; über den Plan hinaus die Komponentenliste, der Backfill-Absatz und die Fehlermodell-Zeile, weil sie „Transformationsregel" allein nannten. |

**Festgelegt im Slice** (Lücken im Text von
[`ADR-0137`](../../adr/0137-routing-zustellziele-persistiertes-ziel-label.md), die
die Spec schließt; Ergebnis je Punkt unter der Liste, Beleg ist der Spec-Text selbst —
es ist Spec-Festlegung, kein Messergebnis, und als solche in `SPEC-032`/`SPEC-019`
als Zusage an die Umsetzung gekennzeichnet):

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
  einheitlich (`ADR-0138` führt für den RPC „leere Liste, kein Fehler" als
  Erwartung; die Spec legt es fest oder benennt den Gegenfall).
- Aufteilung `SPEC-019` / neue Kennung; „führende Stelle je Sachverhalt" wie in
  `slice-transformationen-spec-nachzug` (Anwendbarkeit und Reihenfolge, Fehlertexte,
  Abhilfe — je an genau einer Stelle).

Ergebnis (Plan-Stand, vor dem Review): `when.column` zeichengenau gegen den
Katalog, nicht leer, ohne U+0000; `when.equals` Zeichenkette (leer zulässig, keine
eigene Längengrenze), zeichengenauer Vergleich ohne Normalisierung; `order` JSON-
positive ganze Zahl (Wortlaut der Entscheidung, keine eigene Obergrenze der Spec);
Prüfreihenfolge Formzeilen, dann R1 bis R5, `remove_route`
nur Regelname und R6; R4 hat drei Fehlertexte (zweite Regel ohne `when`, Regel ohne
`when` nicht an höchster `order`, Regel mit `when` hinter der Regel ohne `when`);
Antrag ohne laufende Bindung endet `applied` (Muster der Transformationen);
Lesewege: ein `target`, das keine Change trägt oder das Alphabet verletzt (auch mit
U+0000), ist ein Filter ohne Treffer und kein `400`/`InvalidArgument` — in `SPEC-020`,
`SPEC-021`, `SPEC-022` und `SPEC-031` einheitlich, der SQL-Zugriff ausgenommen
(`ADR-0139`); der Routing-Regelstand im Run ist fail-closed (`ADR-0139`); Aufteilung: `SPEC-032` führt Regelform, Bildbasis,
Auswertung und Anwendbarkeit, `SPEC-019` die Antragsarten, R1 bis R6 und die
Fehlertexte, `LH-FA-CFG-008.a` die Abhilfe. Keiner dieser Texte ist am Code
gemessen.

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
30fd6cb5 1 -n -E 'sieben SQL-Funktionen' -- spec
30fd6cb5 0 -n -w route_target -- spec
diff 14 -n -E 'Routingform|LH-FA-CFG-008\.a' -- spec docs/user harness README.md
diff 3 -n -E 'sieben (Arten|Werte|Antragsarten)|fünf übrigen|sechs übrigen' -- spec docs/user harness
diff 22 -n -E 'zehn Felder|dreizehn Felder' -- spec docs/user
diff 12 -n -w origin -- spec
diff 0 -n -w origin -- spec/architecture.md
diff 9 -n -E 'Routing|Zustellziel' -- spec/architecture.md
diff 0 -n -E 'sieben SQL-Funktionen' -- spec
diff 22 -n -w route_target -- spec
```

| Träger | Messung am Parent (`30fd6cb5`, gemessen am 2026-10-01) | Behandlung und Befund am Diff |
|---|---|---|
| Anker und Überschrift „Routingform" | Zeile 1 des Blocks: 5 Zeilen, alle in `spec/pflichtenheft.md` (der Verweis in `LH-FA-CFG-007.a` `:347`, die Überschrift `:349`, drei Historie-Zeilen `:1265`, `:1287`, `:1298`); `docs/user`, `harness` und `README.md`: 0 | Verweis und Überschrift ziehen, Historie-Zeilen bleiben (Protokoll); Anker prüft `make docs-check` (Modul `anchors`). Die Zeile `make doc-trace` in `harness/README.md` nennt `LH-FA-CFG-008` ohne `.a` (nicht in diesem Suchraum) und bewegt sich erst mit `slice-routing-e2e`. **Befund am Diff:** 14 Zeilen, alle in `spec/pflichtenheft.md` — die Überschrift trägt kein „offen" mehr (Anker `#lh-fa-cfg-008a--routingform`, `make docs-check` Exit 0); neu hinzu kommen die Verweise in `LH-FA-CAP-009.a`, `SPEC-002`, `SPEC-032`, `SPEC-008`-Absatz und die Historie-Zeile. Nicht gefunden: ein weiterer Träger in `docs/user`, `harness` oder `README.md` (0). |
| Zählwörter zu Antragsarten und Werten | Zeile 2: 6 Zeilen (`spec/architecture.md:258` „sieben Arten", `spec/pflichtenheft.md:620`–`622`, Historie `:1288`, `harness/targets/schema-rollout.md:74` „die sieben SQL-Funktionen") | Spec-Zeilen ziehen; `harness/targets/schema-rollout.md` gemeldet an `slice-routing-antragsweg` (es beschreibt den Stand des Rollouts); die Historie-Zeile bleibt (Protokoll). **Befund am Diff:** Zeile 2 des Blocks 3 Zeilen — `spec/pflichtenheft.md` Spalte `rule_name` („die fünf übrigen Antragsarten": gezählt neu, `enable`/`disable`/`exclude_column`/`include_column`/`backfill` — richtig stehen gelassen), die Historie-Zeile vom 2026-09-26 und `harness/targets/schema-rollout.md:74` (gemeldet an `slice-routing-antragsweg`, dessen Suchlauf dieselbe Zeile führt). Die Zeile `Die sieben SQL-Funktionen` (`SPEC-019`, Ordnung der Verarbeitung) deckt das Muster der Parent-Zeile nicht; sie ist eigens gemessen (Parent 1, Diff 0) und auf „neun" gezogen. `spec/architecture.md` „sieben Arten" ist „neun Arten". Nicht gefunden: eine weitere Zählung der Antragsarten oder Werte in `docs/user` oder `harness` außer der gemeldeten (Muster der Zeile 2) |
| Nachrichtenschema-Zählwörter | Zeile 3: 22 Zeilen — erwartet unverändert am Diff (das Label ist nicht Teil der Nachrichten) | bei Abweichung ist der Zug falsch, nicht der Zähler. **Befund am Diff:** 22 Zeilen, unverändert; die Zeilen in `SPEC-021`, `SPEC-020` und `SPEC-031` ergänzen nur „ohne die Felder `origin` und `route_target`". Gefunden, nicht im Suchmuster: `docs/user/benutzerhandbuch.md` beschreibt den Filter `schema`/`table` der Wege mit „zwei optionale" (am Parent `:1301`, `:1378`, `:1403`) — bleibt bis zur Umsetzung wahr, das Handbuch zieht `slice-routing-betriebsdoku` (dessen Suchlauf führt diese Stellen); in dieser Spec-Datei trägt `SPEC-020`/`SPEC-021`/`SPEC-022` jetzt `target` neben `schema`/`table`. |
| Schwesterfeld `origin` als Muster für Aufzählungen von Change-Feldern | Zeilen 4 und 5: 11 Zeilen in der Spec, davon 0 in der Architektur-Sicht | die Sicht zählt keine Change-Felder auf; ob sie eine andere Aufzählung trägt, die `route_target` brauchte, klärt das Lesen der Sequenz „Capture" am Diff. **Befund am Diff:** Zeile 4 des Blocks 12 Zeilen (11 am Parent plus `route_target`-Absatz in `SPEC-002`), Zeile 5: 0 in der Architektur-Sicht; die Sequenz „Capture" und der Absatz „Form der Row Images" zählen keine Change-Felder auf, ein Ziel braucht dort keine Aufzählung. Nicht gefunden: eine Aufzählung der Change-Felder in `spec/architecture.md` (0 Treffer für `origin`). Der Sicht fehlte der Zustellziel-Satz; er steht jetzt als Absatz „Zustellziel der Change" |
| Routing in der Architektur-Sicht | Zeile 6: 0 Zeilen | die Sicht sagt zum Routing heute nichts; der Zug fügt höchstens einen Satz in der Capture-Sequenz und die zwei Tabellenzeilen ein. **Befund am Diff:** 9 Zeilen (Zeile 6 des Blocks, `diff`): Komponentenliste `ARC-001` (1), Absatz „Zustellziel der Change" (3 Zeilen), Absatz zu den Routing-Antragsarten (3 Zeilen), Backfill-Absatz (1), Fehlermodell-Zeile (1); die zwei Tabellenzeilen `set_route`/`remove_route` tragen die Wörter nicht. `make docs-check` (`matrix`) Exit 0: kein ADR-/Slice-/Wellen-Bezug in der Sicht |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn die Welle [welle-routing](welle-routing.md)
eröffnet ist und kein anderer Slice in `in-progress/` liegt (WIP-Limit 1). Dieser
Slice kommt vor jedem übrigen der Welle: jeder Folge-Slice liest die hier
festgelegten Zusagen. V1 und V2 (Welle §5) sind mit `ADR-0138` entschieden (erfüllt);
der Slice zieht ihre Zeilen mit (§1, NICHT-Punkt „Zusagen ohne Entscheidung").

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
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B). — **Ausgang: eingetreten,
  im Slice behoben.** Review F-2 (MEDIUM, `target` außerhalb des Alphabets als
  unbedingte Zusage auf allen Lesewegen) und F-3 (LOW, Grenzen von `order` ohne
  Herkunft) trafen genau diese Klasse; beide sind in der Fixrunde am Text
  geschlossen (Verifikations-Report §4). Gelesen von Reviewer und Verifier.
- **Zählwörter und Aufzählungen bleiben stehen** (`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`,
  verkörpert, 15×): „sieben Arten", „die fünf/sechs übrigen Antragsarten". —
  **Ausgang: entfallen für die Zählwörter, eingetreten für einen Nachbarabsatz.**
  Die Zählwörter sind am Diff gezogen (Suchlauf-Block in §3, Verifikations-Report
  §1: „neun“ konsistent); der Nachbarabsatz „Fail-closed vor dem Commit“ blieb
  zunächst stehen (Review F-1, MEDIUM) und ist in der Fixrunde nachgezogen.
- **Doppelquelle Handbuch gegen Pflichtenheft**
  (`BEO-PGC/zwei-quellen-drift-handbuch-gegen-pflichtenheft`, offen, 3×): dieselben
  Fakten stehen später im Handbuch (`betriebsdoku`). — **Ausgang: weiter offen,
  bewusst an die Adresse übergeben.** Dieser Slice berührt das Handbuch nicht
  (Verifikations-Report §2 Zeile 7); Adresse: `slice-routing-betriebsdoku` §1
  (Handbuch verweist auf die führende Stelle der Spec je Sachverhalt, statt
  Fehlertexte zu wiederholen).
- **Doc-Gate-Regeln der Spec.** `matrix` verbietet Verweise Spec → ADR/Slice,
  `ids` verlangt Links auf nackte Kennungen, `structure` misst Zellenlängen (die
  Transformations-Fassung färbte `section-cell-oversized`). — **Ausgang: entfallen.**
  `make docs-check` Exit 0 am Diff-Stand (Verifikations-Report §1; im Closure-Lauf
  vor dem Inhalts-Commit wiederholt).
- **V3 offen, V1/V2 entschieden.** `SPEC-031`, `SPEC-029` und `SPEC-008` tragen die
  Festlegungen aus `ADR-0138`; die Abhilfe-Aussage zur Nichtanwendbarkeit bleibt
  Zusage, bis V3 gemessen ist. — **Ausgang: V1/V2 entfallen** (Spec-Zeilen gezogen,
  Verifikations-Report §5), **V3 weiter offen** (Abhilfe-Aussage als Zusage
  gekennzeichnet, Verifikations-Report §6; Messadresse `slice-routing-kern-label`
  und `slice-routing-e2e`, `ADR-0138`).
- **A-3: Lesart „Zustellung“ im Lastenheft.** Die Spec trägt die Lesart
  „Zustellung = Abruf bzw. Abonnement“ als Zusage; die Klarstellung im Lastenheft
  liegt beim Auftraggeber ([`ADR-0139`](../../adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md)
  „Offen, nicht entschieden“, Review F-5/A-3). — **Ausgang: weiter offen.** Keine
  DoD-Zeile dieses Slice hängt von ihr ab (Verifikations-Report §6); Träger ist
  der Auftraggeber, die Adresse die Welle [welle-routing](welle-routing.md).

## 7. Closure-Notiz

- **Was hat funktioniert:** Das Pflichtenheft trägt den beschlossenen Routing-Stand
  (`LH-FA-CFG-008.a` beantwortet, `SPEC-032` als eigene Regelform, `SPEC-019` mit
  neun Antragsarten, `route_target` in `SPEC-001`/`SPEC-002`, `target` in
  `SPEC-020`/`-021`/`-022`/`-031`, Zusatz-Subjekt in `SPEC-024`, Nichtanwendbarkeit
  in `SPEC-008`), die Sicht trägt neun Arten und den Zustellziel-Satz ohne
  ADR-/Slice-/Wellen-Bezug (Bezug-Grep im Spec-Diff: 0, Verifikations-Report §1).
  Zwei Lücken, die der Text der Entscheidung ließ (Regelstand im Run, `target`
  außerhalb des Alphabets), gingen als Architect-Fragen A-1/A-2 an `ADR-0139`,
  bevor die Spec sie festlegte, statt vom Implementer geraten zu werden. Gemessen im
  Verifier-Lauf: `make gates` Exit 0, `make docs-check` Exit 0, Suchlauf 16 Zeilen
  (Verifikations-Report §1); im Closure-Lauf `make suchlauf-nachmessen` erneut Exit 0.
- **Was ging anders als geplant:** Mehrumfang gegenüber dem Plan, im Plan §3
  ausgewiesen: Komponentenliste, Backfill-Absatz und Fehlermodell-Zeile der Sicht,
  `LH-FA-CAP-009.a` (Ziel der Backfill-Changes, Fail-closed). Eine Fixrunde nach dem
  Review (F-1, F-2 MEDIUM; F-3, F-4 LOW) und ein Nachzug nach der Verifikation
  (`ab4a83c7`: V-1, V-2). `SPEC-029` blieb unberührt (die Festlegung steht in
  `SPEC-008` und `LH-FA-CAP-009.a`, `SPEC-029` verweist generisch); der Plan-Kopf nannte
  sie als zu ändernd und ist angeglichen (V-4). Die Fixrunde wurde nicht von einem
  Reviewer erneut gelesen, sondern vom Verifier am Text geprüft (DoD, Zeile Review).
- **Steering-Loop-Eintrag:** geschärfte Regel, kein neuer Sensor. Ein
  Fixrunden-Bericht (und eine Plan- oder Historie-Zeile), der „X ergänzt“ sagt, ist
  eine Trägerzusage und wird gegen den **Diff** geprüft, nicht gegen den Bericht:
  V-1 war genau das — Plan und Historie nannten `SPEC-021` als geändert, der Diff der
  Fixrunde berührte die Zeile nicht. Der Verifier fand es, indem er den Fixrunden-Diff
  las (Ursprung: **gemessen**, Verifikations-Report §3, §4, V-1). Zweiter Teil: eine
  Lücke im Text der Entscheidung wird vor dem Spec-Text zur ADR (`ADR-0139`), nicht
  vom Implementer als Setzung der Spec geraten (F-1, F-2). Träger der Regel ist die
  Lese-Handlung des Verifiers und des Planners (`AGENTS.md` §3.12 Instanz B); ein
  Sensor ist ausgeschlossen, weil kein Werkzeug „berichtet“ und „geändert“ vergleicht.
- **Beobachtungs-Register (`../observations/`):**
  - **`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`** (verkörpert) — F-1 (MEDIUM):
    neue `evidence/slice-routing-spec-nachzug.md`, Zähler **16×**. V-2 (LOW, Rest des
    zurückgenommenen Wertebereichs im Fehlertext) und V-4 (LOW, Plan-Kopf): vor dem
    Merge vom Verifier gefunden, bekannter Träger-Typ, nach der Deckel-Regel keine
    weitere Datei; sie stehen hier mit Finding-Kennung.
  - **`BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung`** (verkörpert) — F-2
    (MEDIUM, Setzung breiter als die Entscheidung): neue
    `evidence/slice-routing-spec-nachzug.md`, Zähler **12×**. V-1 (LOW, Bericht nennt eine
    Zeile, die der Diff nicht berührt) und F-3 (LOW, `order`-Grenzen ohne Herkunft): vor
    dem Merge gefunden, bekannter Träger-Typ, keine weitere Datei.
  - **`BEO-PGC/zwei-quellen-drift-handbuch-gegen-pflichtenheft`** (offen, 3×): kein neues
    Auftreten; das Handbuch ist nicht berührt, Adresse `slice-routing-betriebsdoku`.
  - F-4 (LOW, Umbruch nach Teilersetzung), F-5 (INFO, Lastenheft-Lesart: bleibt als A-3
    offen, §6) und F-6 (INFO, Use-Case-Namen in der Sicht, Bestandsmuster): **keine
    Beobachtung** — Einzelfälle ohne Träger-Klasse oder als Risiko geführt.
  - Kein Eintrag steht bei 3× oder mehr ohne Ausgang an, den dieser Slice neu erreichte.
- **Folge-Slices:** keine neuen. Die Folge-Slices der Welle
  [welle-routing](welle-routing.md) tragen die Festlegungen: `slice-routing-kern-label`
  (V3-Messung), `slice-routing-antragsweg` (`harness/targets/schema-rollout.md:74`,
  Zählwort „sieben SQL-Funktionen“), `slice-routing-backfill-pfad` und
  `slice-routing-lesewege` (Festlegungen aus `ADR-0139`), `slice-routing-betriebsdoku`
  (Handbuch-Stellen mit „zwei optionale“ Filtern).
- **Risiken aus §6:** je ein Ausgang am Ort: Zusagen als Tatsachen **eingetreten, im
  Slice behoben**; Zählwörter **entfallen** (Nachbarabsatz eingetreten, behoben);
  Doppelquelle Handbuch **weiter offen** (Adresse benannt); Doc-Gate-Regeln
  **entfallen**; V3 **weiter offen**, V1/V2 **entfallen**; A-3 **weiter offen**
  (Auftraggeber).
- **Drei Paarungen:** dieser Slice gehörte zu [welle-routing](welle-routing.md)
  — die Prüfung lief bei deren Closure (2026-10-02), die DoD-Zeile ist abgehakt.
  (a) Anker: der Lerneintrag verkörpert nichts neu; (b) Folge-Slice: keiner neu,
  die Übergaben stehen in den Plänen der Welle; (c) Register: die genannten Kennungen
  existieren als Verzeichnis, beide fortgeschriebenen tragen ihre Datei unter `evidence/`.
- **Validator (Modul 8):** entfällt — reiner Spec-Nachzug ohne End-Nutzer-Wert; der
  Nutzer-Bedarf ([`LH-FA-CFG-008`](../../../../spec/lastenheft.md)) wird erst durch die
  Umsetzung und den Wellen-Beleg validierbar.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** berührt ist die Default-Sub-Area `*`
(`harness/conventions.md`, Modus-Deklaration: Greenfield, Kürzel `PGC`) mit dem
Pfad `spec/` — eine Sub-Area, keine zu grobe Zusammenfassung.

**Vorgelagert — offene Beobachtungen sichten:** das Register wurde durchgegangen
(Zähler = Zahl der `evidence/`-Dateien, ausgezählt am 2026-10-01); die Treffer
dieses Slice stehen in §6 und in der Eröffnungs-Sichtung der Welle
[welle-routing](welle-routing.md) §6:
`BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (verkörpert, 15×),
`BEO-PGC/zwei-quellen-drift-handbuch-gegen-pflichtenheft` (offen, 3×),
`BEO-PGC/zitat-nennt-die-falsche-stelle` (verkörpert, 9×). Kein offener Eintrag
erreicht mit diesem Slice 3×, den die Welle nicht schon führt.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.

