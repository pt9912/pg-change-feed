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
[`SPEC-024`](../../../../spec/pflichtenheft.md) (Live-Nachrichten
`SPEC-020`/`-021`/`-024`: zehn Felder; `SPEC-022`: dreizehn, `origin`
inbegriffen — die Transformationen ändern keines),
[`LH-FA-CFG-008.a`](../../../../spec/pflichtenheft.md). Der Verweis zeigt
**aufwärts**: die Spec nennt diesen Slice nie.

**Verantwortlich:** Implementer-Agent, 2026-09-26.

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
  Zielname kollidiert mit einer Spalte der Relation) im Erfassungspfad **und**
  im Run (dort gegen die Spalten des Snapshots) samt Aktion und Abhilfe-Weg;
  [`LH-FA-CAP-009.a`](../../../../spec/pflichtenheft.md) sagt in einem Satz,
  dass eine im Run nicht anwendbare Regel den Run `failed` mit der Klasse
  `schema` beendet, ohne Change und ohne den Erfassungspfad zu berühren
  (Folgepflicht 1 von
  [`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md));
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

- [x] [`LH-FA-CFG-007.a`](../../../../spec/pflichtenheft.md) ist beantwortet:
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
- [x] Die Datenstrukturen stehen:
      [`SPEC-019`](../../../../spec/pflichtenheft.md) (zwei Antragsarten, zwei
      Spalten, Menge, `applied`-Bedeutung, Verhalten ohne laufende Bindung,
      Fehlertexte), die Regelform samt Randfällen,
      [`SPEC-002`](../../../../spec/pflichtenheft.md) (Schlüsselmenge folgt dem
      Regelstand), die Zeile `schema` in §4 (Erfassungspfad und Run) und der
      Satz zum Run in
      [`LH-FA-CAP-009.a`](../../../../spec/pflichtenheft.md); eine neu vergebene Kennung ist
      die nächste freie (*zu belegen durch* `grep -o 'SPEC-0[0-9][0-9]'
      spec/pflichtenheft.md | sort -u` am Parent-Stand); §7 Historie trägt je
      Änderung eine Zeile ohne ADR-/Slice-Bezug.
- [x] `spec/architecture.md` trägt die Antragsarten-Tabelle mit den zwei
      weiteren Arten und die Aussage zur Regelform vor der Persistierung, ohne
      ADR-/Slice-/Wellen-Bezug. *Zu belegen durch:* `make docs-check`
      (`matrix`-Modul) und der Suchlauf in §3.
- [x] `make gates` grün — Exit-Code des Laufs ungefiltert gesichert und
      gesondert ausgewertet ([`AGENTS.md`](../../../../AGENTS.md) §3.9).
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor, kein
      offenes HIGH/MEDIUM (`.harness/skills/reviewer.md`,
      `docs/reviews/review-slice-transformationen-spec-nachzug.md`: 0 HIGH,
      F-1 MEDIUM in der Fixrunde behoben) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6), kein
      Self-Review (Modul 8).
- [x] §3.13-Suchlauf: das committete Feld in §3 trägt Gefundenes **und**
      Nichtgefundenes je Träger, beide Stände gemessen (Parent und Diff)
      ([`AGENTS.md`](../../../../AGENTS.md) §3.13).
- [x] Doku-Update: entfällt als eigener Punkt — der Slice **ist** das
      Doku-Update der Spec; das Benutzerhandbuch bleibt unberührt (§1).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke).
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
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
| `spec/pflichtenheft.md` §4 (Zeile `schema`) und `LH-FA-CAP-009.a` (Absatz „Sichtbarkeit und Fehler des Runs“) | update | Bedingung, Aktion und Abhilfe-Weg der Fehlerklasse, im Erfassungspfad und im Run. |
| `spec/pflichtenheft.md` §7 Historie | update | je Änderung eine Zeile, ohne ADR-/Slice-Bezug. |
| `spec/architecture.md` (Antragsart-Tabelle, Capture-Sequenz) | update | die zwei weiteren Antragsarten; Row Images tragen den Regelstand vor der Persistierung. |
| `spec/pflichtenheft.md` §2 `SPEC-030` (neu; Kennung nächste freie am Parent-Stand `1718546b`: dort endet die Reihe bei `SPEC-029`) | neu (Plan-Nachzug) | die Regelform als eigene Datenstruktur statt als Breite in `SPEC-019` (Punkt c der Zielformulierung: die Entscheidung trifft der Slice). Die Randfälle sind hier festgelegt: leeres `values` und Zielname gleich Quellname enden `failed` (Letzteres folgt aus K3), die Abbildung eines Werts auf sich selbst ist zulässig. |
| `spec/pflichtenheft.md` §2 `SPEC-019` (Konfliktfreiheit K1–K4 und Fehlertext-Tabelle) | update (Plan-Nachzug) | die Fehlertexte stehen mit Klartext, Adresse und Prüfreihenfolge dort, nicht in der Regelform — sie sind Antragsverhalten. |
| `spec/pflichtenheft.md` §1 `LH-FA-CAP-009.a` (Absätze Markierung, Fail-closed vor dem Commit, Sichtbarkeit und Fehler des Runs) | update (Plan-Nachzug, über die Zielformulierung Punkt e hinaus) | neben dem Run-Satz der Klasse `schema` tragen zwei bestehende Aufzählungen Eigenschaften, die der Regelstand bewegt: „dieselbe Bild-Konstruktion … ausgeschlossene Spalten“ und „Ausschlussstand entspricht dem Stand der Blöcke“ nennen ihn jetzt mit (Abweichung des Regelstands im Lauf endet `configuration`). |
| `spec/pflichtenheft.md` §4 Absatz „Nicht anwendbare Regel (Klasse `schema`)“ | neu (Plan-Nachzug) | Aktion und Abhilfe-Weg je Pfad stehen unter der Tabelle statt in der Zeile (Zellenlänge des `structure`-Moduls). |
| `spec/architecture.md` §1 `ARC-001`-Zeile, Backfill-Sequenz (Fail-closed-Zeile, Absatz zur Ausführung), §5 Zeile „Dekodier- und Schemafehler“ | update (Plan-Nachzug, über Punkt f hinaus) | dieselben bewegten Eigenschaften in den Aufzählungen der Sicht: Domänenobjekte, Fail-closed-Prüfung, Fehlerklasse `schema`. `ARC-002`/`ARC-004`/`ARC-005` bleiben unverändert: ihre Zeilen nennen Fähigkeiten, keine Antragsarten oder Regeltypen (Suchlauf unten). |
| `spec/pflichtenheft.md` §2 `SPEC-030` (Absatz „Bezeichner", Tabelle, Abschnitt „Beispiele") | update (Fixrunde, Review-Befund F-1/F-2/F-5/F-9) | Bezeichner-Vergleich für `column`, `to` und den Regelnamen (Festlegungen unten); die „Position" der Schlüssel wird zu „Reihenfolge und Inhalt" (nicht zugesagte Schlüsselreihenfolge, zugesagter Inhalt); je Regeltyp ein Beispiel, dazu ein abgelehnter Antrag und die nicht anwendbare Regel; `SPEC-030` ist die führende Stelle für Anwendbarkeit und Reihenfolge. |
| `spec/pflichtenheft.md` §2 `SPEC-019` (Fehlertext-Tabelle, K3, Zeile `rule_name`) | update (Fixrunde, F-1/F-3) | drei neue Zeilen-Fälle (Regelname ungültig, `rule_spec` NULL/kein Objekt, Pflichtschlüssel getrennt vom unbekannten Regeltyp), Prüfreihenfolge der fünf Formzeilen, K3 zeichengenau. |
| `spec/pflichtenheft.md` §1 `LH-FA-CFG-007.a` (Auswertungsreihenfolge, Wirkort, Nicht anwendbare Regel, Abhilfe) | update (Fixrunde, F-2/F-9) | Reihenfolge und Anwendbarkeit verweisen auf `SPEC-030` statt sie zu wiederholen; die Abhilfe ist hier führend und trägt das Ersetzen der Regel; „byte-gleich" und „dieselbe Form" sind auf den Inhalt der Row Images (Schlüsselmenge und Werte) umformuliert. |
| `spec/pflichtenheft.md` §4 (Zeile `schema`, Absatz „Nicht anwendbare Regel") | update (Fixrunde, F-4/F-9) | die Abhilfe der Zeile gilt nur der nicht anwendbaren Regel; der Absatz verweist für Ursache, Erfassungspfad und Abhilfe und trägt nur noch den Run. |
| `spec/pflichtenheft.md` §7 Historie | update (Fixrunde) | eine weitere Zeile ohne ADR-/Slice-Bezug. |
| `spec/pflichtenheft.md:209`, `spec/architecture.md:402` | update (Fixrunde, F-6) | Zeilenumbruch nach der Teilersetzung, nur Form. |
| dieser Plan | update (Fixrunde) | §3 (diese Zeilen, Festlegungen, Suchlauf), §6, §8 (Register-Zahlen mit Messlauf), DoD-Haken. |

**Festlegungen des Slice** (Festlegungen der Spec über den Text von
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
hinaus; jede berührt keine ADR-Aussage, sie präzisiert
[`LH-FA-CFG-007.a`](../../../../spec/pflichtenheft.md); der Beleg des Bestands
steht je Zeile):

| Festlegung | Begründung und Bestands-Beleg |
|---|---|
| `column`: zeichengenau gegen den Katalog, keine Faltung, kein Quoting, nicht leer, ohne U+0000. | Der Bestand vergleicht Spaltennamen so: `ColumnExists` (`internal/adapters/driven/postgresstorage/tableactivation.go`) fragt `column_name = $3` gegen `information_schema.columns`, `BuildRowImage`/`containsName` (`internal/domain/model/rowimage.go`) vergleicht mit `==`, der Antrags-Konstruktor (`NewAdministrationRequest`) verlangt nur nichtleer, und das Adressformat `schema.table.column` steht in `excludecolumn`/`includecolumn`. U+0000 lehnt die Spec ab, weil der Katalog-Parameter mit dem Zeichen scheitert (Persistenzfehler statt `failed`-Antrag). |
| `to`: nicht leer, ohne U+0000, höchstens 63 Byte UTF-8, jedes weitere Zeichen zulässig; K3 zeichengenau. | Der Bestand prüft keinen Zielnamen. Kleinste sichere Form: der Name ist ein JSON-Schlüssel (jedes Zeichen außer U+0000 lässt sich in `jsonb` speichern) und bleibt im Längenraum der Quell-Bezeichner (`identifierShape` `{1,63}`); prüfbar durch Längenmessung, streng statt interpretierend. |
| `rule_name`: `[a-z0-9_]{1,63}`, Ungültiges wird abgelehnt, nicht gefaltet; Eindeutigkeit (K1) zeichengenau. | Das Alphabet ist das der Aktivierung für Schema und Tabelle (`identifierShape` in `tableactivation.go`); der Name geht in Fehleradressen und in eine Zerlegung `schema.table.name`, die dadurch eindeutig bleibt. |
| `column_name` bleibt bei beiden Transformations-Antragsarten NULL. | Die Spalte der Regel steht in `rule_spec`; `column_name` bleibt die Ziel-Spalte der beiden Spalten-Antragsarten (`SPEC-019`). Ein Antrag trägt so je Art genau eine Adressierungsform. |
| Prüfreihenfolge: Regelname, Form der `rule_spec`, Regeltyp, Schlüssel, Pflichtschlüssel und Werte, dann K1 bis K4; die Fehlertext-Wortlaute der Tabelle. | Der ADR-Text nennt die Verletzungen, keine Texte und keine Reihenfolge; die Form folgt dem Bestand (`ErrSourceColumnMissing`: Klartext, Doppelpunkt, Adresse). Die Reihenfolge trennt Form von Konfliktfreiheit, damit der erste Text eindeutig bestimmt ist. |
| Schlüsselreihenfolge im Image nicht zugesagt, zugesagt ist der Inhalt. | `jsonb` (`SPEC-002`) bewahrt die Reihenfolge nicht; über `cdc.changes` und `GET /changes` ist sie nicht beobachtbar. Der Assembler-Test darf sie zusätzlich prüfen, die Spec sagt sie nicht zu. |
| `remove_transformation` durchläuft nur die Regelnamen-Zeile und K4. | Es trägt keine `rule_spec` und keine Spalte. |

**Bekanntes Wortlaut-Versehen der ADR.**
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Folgepflicht 5 (b) nennt den Antragsstatus `requested`; der Bestand kennt nur
`pending`/`applied`/`failed`
(`tools/schema/nacharbeit-administration.sql`, `queries.go`,
[`ADR-0050`](../../adr/0050-sql-administration-antragsqueue-und-live-reload.md)).
Die Spec folgt dem Bestand (`pending`); die `Accepted` ADR bleibt unverändert
([`AGENTS.md`](../../../../AGENTS.md) §3.5), `slice-transformationen-e2e-abhilfe`
übersetzt das Wort bereits selbst. Eine Berichtigungs-ADR ist nicht nötig.

**Führende Stelle je Sachverhalt** (die andere Stelle verweist):
Anwendbarkeit, Auswertungsreihenfolge und Inhalt der Images → `SPEC-030`;
Abhilfe im gescheiterten Prozess → `LH-FA-CFG-007.a`; Fehlertexte,
Prüfreihenfolge und Konfliktfreiheit → `SPEC-019`. `SPEC-008` trägt die
Klassen-Zeile und die Run-Aussage.

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „`LH-FA-CFG-007.a`
ist beantwortet; die Antragsarten-Menge trägt zwei weitere Werte; Row Images
folgen dem Regelstand“; beide Stände gemessen: Parent und Diff — der
Implementer trägt Gefundenes und Nichtgefundenes je Zeile ein):**

| Träger | Suchbefehl | Befund | Behandlung |
|---|---|---|---|
| Anker auf die Überschrift „Transformationsform offen“ in anderen Dokumenten | `grep -rn 'transformationsform' --include=*.md .` (Parent: `git grep -n -i 'transformationsform' 1718546b -- '*.md'`; zusätzlich der Anker: `git grep -n 'lh-fa-cfg-007a\|transformationsform-offen' 1718546b`) | Parent: die Überschrift `spec/pflichtenheft.md:268` „Transformationsform offen“ und viele Treffer, die nur den Dateinamen der ADR (`…transformationsform-deklarative-regeln…`) tragen; **kein** Anker-Verweis auf `#lh-fa-cfg-007a…` in irgendeinem Dokument (zweiter Befehl: 0 Treffer). Diff: die Überschrift lautet „Transformationsform“ (Zeile 275), außerhalb der Datei-Namens-Treffer bleibt nur der Zitat-Text im Plan dieses Slice. | nichts mitzuziehen; Records nicht berührt (keine Zitat-Korrektur nötig) |
| Anker-Erhalt der Überschrift im Pflichtenheft | `make docs-check` (Modul `anchors`) | grün (1192 Dateien, 0 Befunde, Lauf am Diff-Stand) | — |
| Aufzählungen der Antragsarten | `grep -rn 'exclude_column' spec docs/user harness README.md` (Parent: `git grep -n 'exclude_column' 1718546b -- spec docs/user harness README.md`) | Parent: 10 Zeilen in 5 Dateien (`spec/architecture.md` 1, `spec/pflichtenheft.md` 5, `docs/user/benutzerhandbuch.md` 2, `harness/README.md` 1, `harness/targets/schema-rollout.md` 1). Diff: 11 Zeilen (die eine zusätzliche Zeile ist der Satz in `LH-FA-CFG-007.a`, den dieser Diff schreibt). | Spec-Zeilen gezogen (Antragsart-Tabelle der Architektur, `SPEC-019` Einleitung und `request_kind`-Menge). **Gemeldet, nicht mitgeändert:** `docs/user/benutzerhandbuch.md:265` (ein SQL-Beispiel des Ausschlusses, keine Aufzählung der Arten — Handbuch-Zug an `slice-transformationen-betriebsdoku`, Welle §4); `harness/targets/schema-rollout.md:74` („die fünf SQL-Funktionen“ beschreibt den Stand des Rollouts und bewegt sich mit `slice-transformationen-antragsweg-schema`); `harness/README.md:141` (Beschreibung des E2E-Laufs, bewegt sich mit den E2E-Slices). |
| Anzahl-Formulierungen zu Antragsarten und Werten der Menge | `grep -rn 'Antragsart' spec` (Parent: `git grep -n 'Antragsart' 1718546b -- spec`) und `grep -rn -i 'fünf Arten\|sieben Arten\|fünf übrigen\|drei übrigen\|sechs übrigen\|fünf Werte\|sieben Werte' spec` | Parent: 15 Zeilen mit „Antragsart“ (5 in der Architektur, 10 im Pflichtenheft); Zählwörter: „fünf Arten“ (`spec/architecture.md:248`), „die drei übrigen Antragsarten“ (`spec/pflichtenheft.md:535`), die `request_kind`-Menge mit fünf Werten (`:536`) und der Historien-Satz „(fünf Werte)“ (`:876`). Diff: 26 Zeilen; Zählwörter „sieben Arten“ (Architektur), „die fünf übrigen“/„die sechs übrigen“ (Spalten `column_name`/`rule_name`/`rule_spec`), die Menge mit sieben Werten; „fünf Werte“ steht nur noch in der Historie-Zeile vom 2026-09-24. | Zahlwörter und Aufzählungen auf die gemessene Menge gezogen (fünf Werte am Parent-Stand, `tools/schema/nacharbeit-administration.sql` Zeile 62 laut Welle §5; sieben nach diesem Diff); die Historie-Zeile bleibt (Protokoll, kein Zustand). Die Zählwörter „zehn Felder“ stehen nur in `SPEC-021` (`:621` am Parent) und `SPEC-024` (`:697` am Parent), unverändert; ein Zählwort „dreizehn“ steht im Pflichtenheft nicht (der Plan-Text dieses Slice nennt es für `SPEC-022`, gelesen, nicht nachgezählt, nicht geändert). |
| Fehlerklassen-Zeile `schema` in weiteren Trägern | `grep -rn 'nicht sicher interpretierbar' spec docs/user harness` (Parent: `git grep -n 'nicht sicher interpretierbar' 1718546b -- spec docs/user harness`) | Parent: 2 Zeilen (`spec/pflichtenheft.md:783`, `docs/user/benutzerhandbuch.md:1524`). Diff: 2 Zeilen, `spec/pflichtenheft.md:963` gezogen. | Spec-Zeile gezogen; Handbuch §6 (`docs/user/benutzerhandbuch.md`, Fehlerklassen-Tabelle) gemeldet an `slice-transformationen-betriebsdoku` |
| `LH-FA-CFG-007` in Trägern der Abdeckung | `grep -rn 'CFG-007' docs harness .d-check.yml` (Parent: `git grep -n 'CFG-007' 1718546b -- docs harness .d-check.yml`, ohne Planungs-/ADR-/Review-Records) | Parent: 1 Zeile außerhalb der Records, `harness/README.md:134` (Waisen-Aufzählung der `make doc-trace`-Zeile: „3 Waisen — `LH-FA-CAP-009`/`LH-FA-CFG-007`/`LH-FA-CFG-008`“, dort als Messung vom 2026-09-23 datiert). `docs/user/e2e-abdeckung.md`: kein Treffer. Diff: unverändert. Nachmessung `make doc-trace` am Diff-Stand: 80 Anforderungen, 2 Waisen (`LH-FA-CFG-007`, `LH-FA-CFG-008`). | `docs/user/e2e-abdeckung.md` unverändert bis `slice-transformationen-e2e-wirkung` (der Runner schreibt die Zeile). Die datierte Zahl „3 Waisen“ in `harness/README.md:134` ist gegen die heutige Messung („2 Waisen“) überholt, unabhängig von diesem Diff (`LH-FA-CAP-009` trägt inzwischen Belege): gemeldet an den Planner, nicht mitgeändert. |
| Architektur-Sicht: Aufzählungen in `ARC-002`/`ARC-004`/`ARC-005` | `git grep -n 'ARC-002\|ARC-004\|ARC-005' 1718546b -- spec/architecture.md` und Lesen der drei Zeilen in §1 | Parent: die Zeilen in §1 nennen „Konfiguration“ (`ARC-002`), „Fähigkeitsschnittstellen“ (`ARC-004`), „SQL-Funktionen/Views“ (`ARC-005`) — keine Aufzählung von Antragsarten oder Regeltypen. Diff: unverändert. | kein Nachzug nötig |
| Überholter Text im selben Dokument (Risiko §6) | Lesen von `spec/pflichtenheft.md` §1 `LH-FA-CFG-007.a`, §2 `SPEC-019` (Einleitung, Tabelle, Fließtext) und `spec/architecture.md` §4 von oben nach unten am Diff-Stand | Einleitungssatz von `SPEC-019` führte die fünf bestehenden SQL-Funktionen auf; Überschrift trug „offen“; die Architektur-Zeile „fünf Arten“; die Sequenz-Zeile „Bindung und Ausschlussstand erneut prüfen“; die `column_name`-Zelle „die drei übrigen Antragsarten“. Diff: alle fünf gezogen. | — |

**§3.13-Suchlauf der Fixrunde** (bewegte Aussagen: der Bezeichner-Vergleich,
die Abhilfe, die Fehlertext-Tabelle mit ihrer Prüfreihenfolge, „Position"/
„byte-gleich" der Row Images; Symbolname + Zählwort + Beschreibung; Parent ist
`a32a1931`, Diff ist der Arbeitsbaum der Fixrunde):

```text
git grep -n -E 'zeichengenau|Groß-/Kleinschreibung' a32a1931 -- spec
grep -rn -E 'zeichengenau|Groß-/Kleinschreibung' spec
git grep -n -E 'byte-gleich|Position seiner|Position der Quellspalte' a32a1931 -- spec docs/user
grep -rn -E 'byte-gleich|Position seiner|Position der Quellspalte' spec docs/user
git grep -n -E 'drei Formzeilen|fünf Formzeilen|Regelname ist ungültig' a32a1931 -- spec
grep -rn -E 'drei Formzeilen|fünf Formzeilen|Regelname ist ungültig' spec
git grep -n -E 'Abhilfe: Regelstand|Abhilfe nur bei' a32a1931 -- spec
grep -rn -E 'Abhilfe: Regelstand|Abhilfe nur bei' spec
git grep -n 'Abhilfe' a32a1931 -- spec docs/user harness
grep -rn 'Abhilfe' spec docs/user harness
```

| Bewegte Aussage | Befund Parent → Diff | Behandlung |
|---|---|---|
| Bezeichner-Vergleich (`zeichengenau`) | Parent: 1 Zeile (`spec/pflichtenheft.md:908`, nur `map_value`-Wert); Diff: 10 Zeilen, alle in `spec/pflichtenheft.md` (Bezeichner-Absatz, K3-Zusatz, Adressformat, Vergleich, Anwendbarkeit); `spec/architecture.md` und `docs/user`: kein Treffer an beiden Ständen | Architektur-Sicht nennt keine Bezeichner-Regeln, nichts nachzuziehen; das Handbuch nennt die SQL-Funktionen erst mit `slice-transformationen-betriebsdoku` |
| „Position"/„byte-gleich" der Row Images | Parent: 3 Zeilen (`spec/pflichtenheft.md:205`, `:299`, `:301`); Diff: 1 Zeile (`:205`) | `:299`/`:301` sind ersetzt. **Gemeldet, nicht geändert:** `:205` (`LH-FA-CAP-009.a` Absatz „Markierung": „das Row Image ist byte-gleich dem WAL-Image derselben Zeile") stammt aus dem Backfill-Stand und gilt dem Assembler-Ausgang, nicht dem `jsonb`-Lesepfad; ob die Zeile auf „Inhalt" umzuformulieren ist, entscheidet der Planner |
| Fehlertext-Tabelle und Prüfreihenfolge | Parent: „drei Formzeilen" (`:681`); Diff: „fünf Formzeilen" (`:683`), die neue Zeile `Regelname ist ungültig` (`:689`); ein weiteres Zählwort zu Formzeilen steht nirgends | gezogen; die Historie-Zeile vom 2026-09-26 zur ersten Fassung bleibt (Protokoll) |
| „Abhilfe: Regelstand ändern" in der `schema`-Zeile | Parent: 1 Zeile (`:963`); Diff: die Zeile trägt „Abhilfe nur bei nicht anwendbarer Regel" (`:1020`) | gezogen |
| Abhilfe an mehreren Stellen | Parent: `spec/pflichtenheft.md` 8 Zeilen (`:334`, `:337`, `:963`, `:970`, `:973`, `:976`, zwei Historie-Zeilen); Diff: 11 Zeilen, `LH-FA-CFG-007.a` führt (`:332`, `:335`), `SPEC-008` verweist (`:1020`, `:1027`, `:1030`); `docs/user` und `harness`: an beiden Ständen 3 Zeilen, keine gilt der Transformations-Abhilfe (WAL-Rückstand, Handbuch-Historie) | Führung deklariert (Absatz „Führende Stelle je Sachverhalt" oben); nichts außerhalb der Spec mitzuziehen |

**Belege des Laufs (Implementer):**

- Kennung: `git grep -h -o 'SPEC-0[0-9][0-9]' 1718546b -- spec/pflichtenheft.md | sort -u | tail -1` druckt `SPEC-029`; am Diff-Stand (`HEAD`) `SPEC-030` — die neu vergebene Kennung ist die nächste freie.
- Decken-Regel: `git diff 1718546b HEAD -- spec | grep '^+' | grep -c 'ADR-0\|slice-\|welle-'` druckt `0` (gemessen am Diff-Stand nach dem Spec-Commit).
- Handbuch unberührt: `git diff 1718546b HEAD -- docs/user/benutzerhandbuch.md | wc -l` druckt `0` — keine Versionshistorie-Pflicht.
- Sensoren: `make docs-check` Exit 0 (1192 Dateien, 0 Befunde); `make gates` Exit 0 am Stand nach dem Plan-Nachzug-Commit (`coverage-gate: OK — Coverage 83.40% erfüllt Schwelle 80%`, `commit-traceability: OK — 5 Commit(s)`, `generated-sync: OK`, `gesamt: 0 Befund(e)`); `make doc-trace`: `80 Anforderung(en), 2 Waise(n)` (`LH-FA-CFG-007`, `LH-FA-CFG-008`).
- Zusage · mutierte Eingabe · gesehenes Rot (die Spec-Zusage „Spec-Straten tragen keinen ADR-Bezug und keine nackte Kennung“ wird vom `matrix`- und `ids`-Modul gehalten): (1) `spec/pflichtenheft.md` bekommt einen Link auf `ADR-0112` → `make docs-check` Exit 2, `matrix-forbidden` „Referenz spec → adr ist nicht erlaubt“; (2) `spec/architecture.md` bekommt `LH-FA-CFG-007` ohne Link → Exit 2, `id-unlinked`; (3) zusätzlich färbte die erste Fassung des `ARC-001`-Zeilentexts das `structure`-Modul (`section-cell-oversized`, 247 von 220 Zeichen) — das Rot wurde vor dem Kürzen gesehen. Beide Mutationen wurden zurückgenommen (`cmp` gegen die Sicherungskopie: gleich). Die Inhalts-Zusagen der Spec selbst (Fehlertexte, Randfälle, Zusage-Form) hält kein Sensor; ihre Leser sind Reviewer und Verifier (Risiken §6).

- Fixrunde (Review `docs/reviews/review-slice-transformationen-spec-nachzug.md`, Parent `a32a1931`): Decken-Regel `git diff a32a1931 -- spec | grep '^+' | grep -c 'ADR-0\|slice-\|welle-'` druckt `0`; Handbuch `git diff a32a1931 -- docs/user | wc -l` druckt `0`; `make docs-check` Exit 0 (`d-check: 1193 Datei(en) geprüft, 0 Befund(e)`); `make gates` Exit 0 (`coverage-gate: OK — Coverage 83.40% erfüllt Schwelle 80%`, `generated-sync: OK`, `gesamt: 0 Befund(e)`); `make doc-trace`: `80 Anforderung(en), 2 Waise(n)`.
- Fixrunde, Zusage · mutierte Eingabe · gesehenes Rot: die Zusage „die neuen Absätze (Bezeichner, Beispiele, Fehlertext-Zeilen) tragen keine nackte Kennung" hält das `ids`-Modul: ein angehängter Satz mit `LH-FA-CFG-007` ohne Link am Ende von `spec/pflichtenheft.md` → `make docs-check` Exit 2, Befund `id-unlinked` in der angehängten Zeile; zurückgenommen (`cmp` gegen die Sicherungskopie: gleich). Die Inhalts-Zusagen (Bezeichner-Form, Fehlertext-Wortlaute, Beispiele) hält kein Sensor; ihre Leser sind Reviewer und Verifier, und `slice-transformationen-kern-rename` bindet sie an Tests (Risiko §6, Fehlertexte).

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
  (`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`, verkörpert, 12
  Evidence-Dateien, gemessen am 2026-09-26 mit `ls
  docs/plan/planning/observations/BEO-PGC/<eintrag>/evidence | wc -l`): die
  Überschrift „… offen“, der Einleitungssatz von
  [`SPEC-019`](../../../../spec/pflichtenheft.md) (nennt die
  bestehenden SQL-Funktionen) und der Satz zur Menge der Arten in `spec/architecture.md` können
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

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen; die
Zahlen sind die Zahl der Evidence-Dateien, gemessen am 2026-09-26 (HEAD
`a32a1931`) mit `ls docs/plan/planning/observations/BEO-PGC/<eintrag>/evidence
| wc -l` —
`BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (verkörpert, 12, einschlägig —
Risiko §6), `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert, 32,
Suchlauf §3), `BEO-PGC/dod-begruendung-unzutreffende-tatsachenbehauptung`
(verkörpert, 8, einschlägig — Risiko §6, zweiter Punkt),
`BEO-PGC/zitat-nennt-die-falsche-stelle` (verkörpert, 8, Review-Leser: jede
zitierte Stelle der Spec ist geprüft),
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert, 23, die Zahl
der Antragsarten und Werte trägt ihren Ursprung),
`BEO-PGC/adr-folgepflicht-ohne-traeger-slice` (offen, 1, dieser Slice ist der
Träger der Folgepflicht 1).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
