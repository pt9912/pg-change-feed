# Slice transformationen-spec-nachzug: Spec-Nachzug — Pflichtenheft und Architektur-Sicht tragen den beschlossenen Transformations-Stand

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`, siehe
Baseline-Regelwerk `modul-05-planning-harness.md` §Lifecycle als State Machine.
Übernimmt ein anderer Slice den Gegenstand oder entfällt er, geht diese Datei
aus `open/` oder `next/` nach `done/` — §7 nennt in der Zeile `Gegenstand:`
Kennung oder Grund, die Liefer-Punkte der DoD bleiben leer
(§Ein Slice, dessen Gegenstand ein anderer übernimmt).

**Welle:** [welle-transformationen](../welle-transformationen.md).

**Bezug:** [`LH-FA-CFG-007`](../../../../spec/lastenheft.md) (Transformationen
— Haupt-Bezug), [`LH-FA-CFG-007.a`](../../../../spec/pflichtenheft.md) (die
offene, ADR-pflichtige Frage „Transformationsform“),
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Folgepflicht 1 (Spec-Nachzug — Träger dieses Slice),
[`LH-FA-CFG-008`](../../../../spec/lastenheft.md) (Routing — bleibt offen, nur
abgegrenzt).

**Berührte Spec-Stellen:**
[`LH-FA-CFG-007.a`](../../../../spec/pflichtenheft.md),
[`SPEC-019`](../../../../spec/pflichtenheft.md) (Antrags-Datensatz),
[`SPEC-008`](../../../../spec/pflichtenheft.md) (Fehlerklassen-Zeile `schema`,
§4), [`SPEC-002`](../../../../spec/pflichtenheft.md) (Row Images — die
Schlüsselmenge folgt dem Regelstand), gegebenenfalls eine **neue** Kennung für
die Regelform (`rule_spec`), [`ARC-002`](../../../../spec/architecture.md),
[`ARC-004`](../../../../spec/architecture.md),
[`ARC-005`](../../../../spec/architecture.md) und `spec/architecture.md` §4
(Sequenz). Gelesen, nicht geändert:
[`SPEC-020`](../../../../spec/pflichtenheft.md),
[`SPEC-021`](../../../../spec/pflichtenheft.md),
[`SPEC-022`](../../../../spec/pflichtenheft.md),
[`SPEC-024`](../../../../spec/pflichtenheft.md) (dieselben zehn Felder),
[`LH-FA-CFG-008.a`](../../../../spec/pflichtenheft.md). Der Verweis zeigt
**aufwärts**: die Spec nennt diesen Slice nie.

**Verantwortlich:** — (noch nicht priorisiert).

**Autor:** Planner-Agent, Welle-Eröffnung
[welle-transformationen](../welle-transformationen.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Die beiden führenden Spec-Straten — Pflichtenheft (Rang 2) und
Architektur-Sicht (Rang 3) — tragen den mit
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
beschlossenen Transformations-Stand als Technik-Festlegung und als Sicht,
**bevor** der erste umsetzende Slice startet (die Modus-Deklaration in
`harness/conventions.md` ist Greenfield: die Doku führt). Umfang:

- (a) [`LH-FA-CFG-007.a`](../../../../spec/pflichtenheft.md): die Überschrift
  verliert das Wort „offen“; der Text legt Konfigurationsmechanismus (zwei
  Antragsarten der Antrags-Queue), Ausdrucksform (geschlossener Satz:
  `rename_column`, `map_value`), Auswertungsreihenfolge (Ausschluss, dann
  Spaltenregeln in Relation-Spaltenreihenfolge), Auflösung der Mehrdeutigkeit
  (statisch ausgeschlossen durch vier Konfliktfreiheits-Invarianten K1–K4),
  Verhältnis zum Spaltenausschluss (der Ausschluss gilt zuerst) und Wirkort
  (vor der Persistierung) als Zusagen fest; die Routing-Frage
  [`LH-FA-CFG-008.a`](../../../../spec/pflichtenheft.md) bleibt unberührt
  offen;
- (b) [`SPEC-019`](../../../../spec/pflichtenheft.md): die Antragsarten
  `set_transformation`/`remove_transformation`, die Spalten
  `rule_name`/`rule_spec` (nullable; Pflicht für die beiden Arten, `rule_spec`
  nur bei `set_transformation`), die geschlossene `request_kind`-Menge (Menge
  des Parent-Stands plus zwei Werte), die zweite Bedeutung von `applied`
  (dauerhaft vermerkt: die `applied`-Zeilen sind die einzige Herkunft des
  Regelstands einer Tabelle, Ordnung `requested_at`, dann
  `administration_request_id`), das Verhalten eines Antrags gegen eine Tabelle
  ohne laufende Bindung (Muster der Spalten-Antragsarten) und die
  `failed`-Fehlertexte je Verletzung (K1–K4, unbekannter `kind`, unbekannter
  Schlüssel, nicht geführter Regelname);
- (c) die Regelform: `rule_spec` als JSON-Objekt mit Pflichtschlüssel `kind`,
  je Regeltyp die Schlüssel und die Wirkung auf `old_image` und `new_image`
  samt den Randfällen, die der ADR-Text offen lässt (leeres `values`-Objekt,
  Zielname gleich Quellname, Abbildung eines Werts auf sich selbst) — entweder
  in [`SPEC-019`](../../../../spec/pflichtenheft.md) oder unter der nächsten
  freien Kennung (die Entscheidung trifft der Slice, die Kennung ist die
  nächste freie am Parent-Stand);
- (d) [`SPEC-002`](../../../../spec/pflichtenheft.md): die Schlüsselmenge der
  Row Images folgt dem Regelstand zum Erfassungszeitpunkt; die übrigen Felder
  und die Tabellen-Identität bleiben Quell-Identität;
- (e) §4: die Zeile `schema` der Fehlerklassen-Tabelle nennt zusätzlich die
  Nichtanwendbarkeit einer Regel (Spalte fehlt in der Relation der Change;
  Zielname kollidiert mit einer Spalte der Relation) samt Aktion und
  Abhilfe-Weg;
- (f) `spec/architecture.md`: die Tabellen und Aufzählungen, die die
  Antragsarten nennen (Antragsart-Tabelle im Abschnitt zur SQL-Aktivierung),
  tragen die zwei weiteren Arten; die Capture-Sequenz sagt, dass Row Images die
  durch den Regelstand bestimmte Form tragen, bevor sie persistiert werden —
  **ohne** ADR-, Slice- oder Wellen-Bezug ([`AGENTS.md`](../../../../AGENTS.md)
  §3.4). Ob
  [`ARC-002`](../../../../spec/architecture.md)/[`ARC-004`](../../../../spec/architecture.md)
  Aufzählungen tragen, die nachzuziehen sind, klärt der Suchlauf in §3.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Lastenheft** — bleibt unverändert; das Pflichtenheft präzisiert, erweitert
  nie (Kopf von `spec/pflichtenheft.md`).
- **Routing** ([`LH-FA-CFG-008.a`](../../../../spec/pflichtenheft.md)) — bleibt
  offen; das Routing bekommt eine eigene Entscheidung
  ([`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Folgepflicht 8), die Leitplanken stehen nur in der ADR.
- **Benutzerhandbuch** — beschreibt Betreiber-Oberfläche und darf keine
  Funktion nennen, die es noch nicht gibt oder noch nicht wirkt;
  `slice-transformationen-betriebsdoku` trägt sie (Welle §4, Abweichung 3).
- **Code, Schema, Skripte** — Umsetzung gehört den Folge-Slices.
- **Die Backfill-Anteile von [`SPEC-019`](../../../../spec/pflichtenheft.md)**
  — bleiben Gegenstand von `slice-backfill-spec-nachzug`; dieser Slice
  erweitert deren Stand (Kollisions-Regel: §4 Start).

## 2. Definition of Done

- [ ] [`LH-FA-CFG-007.a`](../../../../spec/pflichtenheft.md) ist beantwortet:
      die Überschrift trägt kein „offen“ mehr, der Text nennt Mechanismus,
      Ausdrucksform, Auswertungsreihenfolge, Mehrdeutigkeits-Auflösung,
      Verhältnis zum Ausschluss und Wirkort als Zusagen in Zukunfts-Form (was
      die Umsetzung liefern **muss**; bis zu den Test-Slices unbelegt, nicht
      als geprüfte Aussage, [`AGENTS.md`](../../../../AGENTS.md) §3.12). Die
      Abhilfe nach einer nicht anwendbaren Regel steht als Zusage, nicht als
      Tatsache —
      [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
      führt sie als „erwartet, nicht am Code belegt“. *Zu belegen durch:* Lesen
      des Abschnitts und `make docs-check`.
- [ ] Die Datenstrukturen stehen:
      [`SPEC-019`](../../../../spec/pflichtenheft.md) (zwei Antragsarten, zwei
      Spalten, Menge, `applied`-Bedeutung, Verhalten ohne laufende Bindung,
      Fehlertexte), die Regelform samt Randfällen,
      [`SPEC-002`](../../../../spec/pflichtenheft.md) (Schlüsselmenge folgt dem
      Regelstand) und die Zeile `schema` in §4; eine neu vergebene Kennung ist
      die nächste freie (*zu belegen durch* `grep -o 'SPEC-0[0-9][0-9]'
      spec/pflichtenheft.md | sort -u` am Parent-Stand); §7 Historie trägt je
      Änderung eine Zeile ohne ADR-/Slice-Bezug.
- [ ] `spec/architecture.md` trägt die Antragsarten-Tabelle mit den zwei
      weiteren Arten und die Aussage zur Regelform vor der Persistierung, ohne
      ADR-/Slice-/Wellen-Bezug. *Zu belegen durch:* `make docs-check`
      (`matrix`-Modul) und der Suchlauf in §3.
- [ ] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [ ] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [ ] Doku-Update: entfällt als eigener Punkt — der Slice **ist** das
      Doku-Update der Spec; das Benutzerhandbuch bleibt unberührt (§1).
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
| `spec/pflichtenheft.md` §1 ([`LH-FA-CFG-007.a`](../../../../spec/pflichtenheft.md)) | update | Überschrift und Text auf den beantworteten Stand; Zusagen-Form. |
| `spec/pflichtenheft.md` §2 ([`SPEC-019`](../../../../spec/pflichtenheft.md), [`SPEC-002`](../../../../spec/pflichtenheft.md), gegebenenfalls neue Kennung) | update / neu | Datenstrukturen; die Regelform als Datenstruktur, nicht als Breiten-Regel (§3). |
| `spec/pflichtenheft.md` §4 (Zeile `schema`) | update | Bedingung, Aktion und Abhilfe-Weg der Fehlerklasse. |
| `spec/pflichtenheft.md` §7 Historie | update | je Änderung eine Zeile, ohne ADR-/Slice-Bezug. |
| `spec/architecture.md` (Antragsart-Tabelle, Capture-Sequenz) | update | die zwei weiteren Antragsarten; Row Images tragen den Regelstand vor der Persistierung. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „`LH-FA-CFG-007.a`
ist beantwortet; die Antragsarten-Menge trägt zwei weitere Werte; Row Images
folgen dem Regelstand“; beide Stände gemessen: Parent und Diff — der
Implementer trägt Gefundenes und Nichtgefundenes je Zeile ein):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Anker auf die Überschrift „Transformationsform offen“ in anderen Dokumenten | `grep -rn 'transformationsform' --include=*.md .` | *(Implementer trägt ein)* | Anker mitziehen; Records (`docs/reviews/**`, `done/**`) nur als Zitat-Korrektur nach [`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md) |
| Aufzählungen der Antragsarten | `grep -rn 'exclude_column' spec docs/user harness README.md` | *(Implementer trägt ein)* | Spec-Stellen dieses Slice ziehen; Handbuch-Stellen an `slice-transformationen-betriebsdoku` melden (Welle §4) |
| Anzahl-Formulierungen zu Antragsarten und Werten der Menge | `grep -rn 'Antragsart' spec` | *(Implementer trägt ein)* | Zahlwörter und Aufzählungen an die gemessene Menge ziehen ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz A) |
| Fehlerklassen-Zeile `schema` in weiteren Trägern | `grep -rn 'nicht sicher interpretierbar' spec docs/user harness` | *(Implementer trägt ein)* | Spec-Zeile dieses Slice ziehen; Handbuch §6 an `slice-transformationen-betriebsdoku` melden |
| `LH-FA-CFG-007` in Trägern der Abdeckung | `grep -rn 'CFG-007' docs harness .d-check.yml` | *(Implementer trägt ein)* | unverändert bis `slice-transformationen-e2e-wirkung` (der Runner schreibt die Abdeckungs-Zeile) |

## 4. Trigger

**Start** (`next` → `in-progress`): wenn die Welle
[welle-transformationen](../welle-transformationen.md) eröffnet ist,
`slice-backfill-spec-nachzug` in `done/` liegt und kein anderer Slice in
`in-progress/` liegt (WIP-Limit 1). Grund der zweiten Bedingung: beide Slices
ändern [`SPEC-019`](../../../../spec/pflichtenheft.md) (Antragsarten-Menge,
`applied`-Bedeutung) und die Kennungsvergabe im Pflichtenheft zählt fortlaufend
je Datei; dieser Slice liest den Stand des Backfill-Nachzugs und erweitert ihn,
statt ihn parallel zu überschreiben. Dieser Slice kommt vor jedem übrigen Slice
dieser Welle: jeder Folge-Slice liest die hier festgelegten Zusagen.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): nicht erwartet — ein
  reiner Doku-Zug über zwei Dateien; sprengte er den Umfang, wäre die
  Architektur-Sicht (Punkt f) der abtrennbare Teil.
- `in-progress` → `open` (blockiert): falls das Übertragen eine Aussage von
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  als widersprüchlich zum Bestand der Spec zeigt (etwa zum Backfill-Stand von
  [`SPEC-019`](../../../../spec/pflichtenheft.md)) — eine `Accepted` ADR wird
  nicht überschrieben ([`AGENTS.md`](../../../../AGENTS.md) §3.5), der
  Widerspruch ginge als Frage an den Architect und ggf. in eine Folge-ADR.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- **Die Spec-Straten dürfen keine ADR- und keine Slice-Kennung tragen**
  (`matrix`-Modul in `.d-check.yml`, [`AGENTS.md`](../../../../AGENTS.md) §3.4)
  — die Begründung der Entscheidung bleibt in
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md).
  *Erwartet, zu belegen durch:* `make docs-check` grün und ein `grep -n
  'ADR-0\|slice-\|welle-'` über beide Spec-Dateien am Diff. **Ausgang:** *(bei
  Closure)*
- **Zusagen als Tatsachen missverstanden**
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B): die Abhilfe im
  gescheiterten Prozess ist in
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  „erwartet“, die Konfliktfreiheit K1–K4 und die Regeltypen sind Festlegungen.
  Das Pflichtenheft darf nur zusagen, was die Umsetzung liefern muss.
  *Erwartet, zu belegen durch:* Review liest jeden Satz der neuen Abschnitte
  auf Zukunfts-Form. **Ausgang:** *(bei Closure)*
- **Überholter Text im selben Dokument**
  (`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`, offen, 2×): die
  Überschrift „… offen“, der Einleitungssatz von
  [`SPEC-019`](../../../../spec/pflichtenheft.md) (nennt die vier
  SQL-Funktionen) und der Satz „vier Arten“ in `spec/architecture.md` können
  die alte Aussage tragen. *Erwartet, zu belegen durch:* Lesen der Abschnitte
  von oben nach unten und Suchlauf §3, Zeilen 2–3. **Ausgang:** *(bei Closure)*
- **Kollision mit `slice-backfill-spec-nachzug`** in
  [`SPEC-019`](../../../../spec/pflichtenheft.md) und in der Kennungsvergabe.
  *Erwartet, zu belegen durch:* der Start-Trigger (der Backfill-Nachzug liegt
  in `done/`) und der Kennungs-Suchlauf am Parent-Stand. **Ausgang:** *(bei
  Closure)*
- **Fehlertexte stehen in der Spec, bevor der Code sie trägt.** Der Wortlaut
  ist bindend für `slice-transformationen-antragsweg-usecase`; eine Abweichung
  dort ist ein Plan-Nachzug hierher, kein stilles Umformulieren. *Erwartet, zu
  belegen durch:* der Review von `antragsweg-usecase` liest die Texte gegen
  [`SPEC-019`](../../../../spec/pflichtenheft.md). **Ausgang:** *(bei Closure)*
- **Die Randfälle der Regelform stehen nicht in
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)**
  (leeres `values`, Zielname gleich Quellname, Abbildung eines Werts auf sich
  selbst). Der Slice legt sie fest; berührt eine Festlegung eine ADR-Aussage,
  geht sie als Frage an den Architect. *Erwartet, zu belegen durch:* Review der
  Regelform gegen
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 2/3. **Ausgang:** *(bei Closure)*

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
`*`/`PGC` (Modus Greenfield, `harness/conventions.md` §Modus-Deklaration); die
Spec-Dateien bilden keine eigene Sub-Area — kein Anlass zur Ausdifferenzierung.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (offen, 2×, einschlägig —
Risiko §6), `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 26×,
Suchlauf §3), `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung`
(verkörpert, 6×, einschlägig — Risiko §6, zweiter Punkt),
`BEO-PGC/zitat-nennt-die-falsche-stelle` (verkörpert, 7×, Review-Leser: jede
zitierte Stelle der Spec ist geprüft),
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert, 13×, die Zahl
der Antragsarten und Werte trägt ihren Ursprung),
`BEO-PGC/adr-folgepflicht-ohne-traeger-slice` (offen, 1×, dieser Slice ist der
Träger der Folgepflicht 1).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
