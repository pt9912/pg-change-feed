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
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag (geschärfte Regel · neuer
      Sensor · benannte Spec-Lücke): §7, benannte Spec-Lücke und neuer
      Register-Eintrag.
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in diesem
      Repo (Greenfield).
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder weitere `evidence/`-Datei; kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert: ein neues Verzeichnis, drei
      Auftreten unter dem Deckel in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen / weiter
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
| `spec/pflichtenheft.md` §1 `LH-FA-CAP-009.a` (Absätze Markierung, Fail-closed vor dem Commit, Sichtbarkeit und Fehler des Runs) | update (Plan-Nachzug, über die Zielformulierung Punkt e hinaus) | neben dem Run-Satz der Klasse `schema` tragen zwei bestehende Aufzählungen Eigenschaften, die der Regelstand bewegt: „dieselbe Bild-Konstruktion … ausgeschlossene Spalten“ und „Ausschlussstand entspricht dem Stand der Blöcke“ nennen ihn mit (Abweichung des Regelstands im Lauf endet `configuration`). |
| `spec/pflichtenheft.md` §4 Absatz „Nicht anwendbare Regel (Klasse `schema`)“ | neu (Plan-Nachzug) | Aktion und Abhilfe-Weg je Pfad stehen unter der Tabelle statt in der Zeile (Zellenlänge des `structure`-Moduls). |
| `spec/architecture.md` §1 `ARC-001`-Zeile, Backfill-Sequenz (Fail-closed-Zeile, Absatz zur Ausführung), §5 Zeile „Dekodier- und Schemafehler“ | update (Plan-Nachzug, über Punkt f hinaus) | dieselben bewegten Eigenschaften in den Aufzählungen der Sicht: Domänenobjekte, Fail-closed-Prüfung, Fehlerklasse `schema`. `ARC-002`/`ARC-004`/`ARC-005` bleiben unverändert: ihre Zeilen nennen Fähigkeiten, keine Antragsarten oder Regeltypen (Suchlauf unten). |
| `spec/pflichtenheft.md` §2 `SPEC-030` (Absatz „Bezeichner", Tabelle, Abschnitt „Beispiele") | update (Fixrunde, Review-Befund F-1/F-2/F-5/F-9) | Bezeichner-Vergleich für `column`, `to` und den Regelnamen (Festlegungen unten); die „Position" der Schlüssel wird zu „Reihenfolge und Inhalt" (nicht zugesagte Schlüsselreihenfolge, zugesagter Inhalt); je Regeltyp ein Beispiel, dazu ein abgelehnter Antrag und die nicht anwendbare Regel; `SPEC-030` ist die führende Stelle für Anwendbarkeit und Reihenfolge. |
| `spec/pflichtenheft.md` §2 `SPEC-019` (Fehlertext-Tabelle, K3, Zeile `rule_name`) | update (Fixrunde, F-1/F-3) | drei neue Zeilen-Fälle (Regelname ungültig, `rule_spec` NULL/kein Objekt, Pflichtschlüssel getrennt vom unbekannten Regeltyp), Prüfreihenfolge der fünf Formzeilen, K3 zeichengenau. |
| `spec/pflichtenheft.md` §1 `LH-FA-CFG-007.a` (Auswertungsreihenfolge, Wirkort, Nicht anwendbare Regel, Abhilfe) | update (Fixrunde, F-2/F-9) | Reihenfolge und Anwendbarkeit verweisen auf `SPEC-030` statt sie zu wiederholen; die Abhilfe ist hier führend und trägt das Ersetzen der Regel; „byte-gleich" und „dieselbe Form" sind auf den Inhalt der Row Images (Schlüsselmenge und Werte) umformuliert. |
| `spec/pflichtenheft.md` §4 (Zeile `schema`, Absatz „Nicht anwendbare Regel") | update (Fixrunde, F-4/F-9) | die Abhilfe der Zeile gilt nur der nicht anwendbaren Regel; der Absatz verweist für Ursache, Erfassungspfad und Abhilfe und trägt nur noch den Run. |
| `spec/pflichtenheft.md` §7 Historie | update (Fixrunde) | eine weitere Zeile ohne ADR-/Slice-Bezug. |
| `spec/pflichtenheft.md:209`, `spec/architecture.md:402` | update (Fixrunde, F-6) | Zeilenumbruch nach der Teilersetzung, nur Form. |
| dieser Plan | update (Fixrunde) | §3 (diese Zeilen, Festlegungen, Suchlauf), §6, §8 (Register-Zahlen mit Messlauf), DoD-Haken. |
| `spec/pflichtenheft.md` §1 `LH-FA-CAP-009.a` (Absatz „Markierung“), §4 `SPEC-008` (Absatz „Nicht anwendbare Regel“), §7 Historie | update (Closure, Verifikation V-4/V-5) | „inhaltsgleich (Schlüsselmenge und Werte)“ statt „byte-gleich“ für das Backfill-Row-Image am Lesepfad; ein Halbsatz, dass ein Prozessneustart für den neuen Run nicht Teil der Zusage ist; eine Historie-Zeile ohne ADR-/Slice-Bezug. |
| die offenen Pläne `slice-transformationen-kern-rename`, `-map-value`, `-backfill-pfad`, `-antragsweg-usecase`, `-antragsweg-schema` | update (Closure, Übergaben nach [`AGENTS.md`](../../../../AGENTS.md) §3.13 aus Verifikation V-3 und §7) | Position und Byte-Gleichheit sind Aussagen am Ausgang des `Assembler`, die Spec sagt am Lesepfad Schlüsselmenge und Werte zu; Form-Prüfungen von `column`/`to` in der Domäne, Sentinel für `to` gleich `column` auf den K3-Text, Zwischenzustand der Regeltyp-Menge bis `map-value`, Parameterlisten der SQL-Funktionen bei `antragsweg-schema`. Nur Pläne, keine Umsetzung. |
| dieser Plan | update (Closure) | Suchlauf-Berichtigungen (Träger `harness/README.md`, Zeilenzahl, Lokator), Vorspann der Festlegungen, Suchlauf der Closure, §6 Ausgänge, §7. |

**Festlegungen des Slice** (Festlegungen der Spec über den Text von
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
hinaus; sie präzisieren
[`LH-FA-CFG-007.a`](../../../../spec/pflichtenheft.md); der Beleg des Bestands
steht je Zeile). Eine Zeile — „Schlüsselreihenfolge im Image nicht zugesagt“ —
schwächt eine ADR-Aussage auf Spec-Ebene ab: die Festlegung von
[`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
Teilfrage 3, der umbenannte oder abgebildete Schlüssel behalte die Position
seiner Quellspalte, gilt am Ausgang des `Assembler` und bleibt dort prüfbar; die
Spec sagt am Lesepfad (`jsonb`) Schlüsselmenge und Werte zu. Das ist kein
Widerspruch zur ADR, weil beide Aussagen verschiedene Stellen betreffen:

| Festlegung | Begründung und Bestands-Beleg |
|---|---|
| `column`: zeichengenau gegen den Katalog, keine Faltung, kein Quoting, nicht leer, ohne U+0000. | Der Bestand vergleicht Spaltennamen so: `ColumnExists` (`internal/adapters/driven/postgresstorage/tableactivation.go`) fragt `column_name = $3` gegen `information_schema.columns`, `BuildRowImage`/`containsName` (`internal/domain/model/rowimage.go`) vergleicht mit `==`, der Antrags-Konstruktor (`NewAdministrationRequest`) verlangt nur nichtleer, und das Adressformat `schema.table.column` steht in `excludecolumn`/`includecolumn`. U+0000 lehnt die Spec ab, weil der Katalog-Parameter mit dem Zeichen scheitert (Persistenzfehler statt `failed`-Antrag). |
| `to`: nicht leer, ohne U+0000, höchstens 63 Byte UTF-8, jedes weitere Zeichen zulässig; K3 zeichengenau. | Der Bestand prüft keinen Zielnamen. Kleinste sichere Form: der Name ist ein JSON-Schlüssel (jedes Zeichen außer U+0000 lässt sich in `jsonb` speichern) und bleibt im Längenraum der Quell-Bezeichner (`identifierShape` `{1,63}`); prüfbar durch Längenmessung, streng statt interpretierend. |
| `rule_name`: `[a-z0-9_]{1,63}`, Ungültiges wird abgelehnt, nicht gefaltet; Eindeutigkeit (K1) zeichengenau. | Das Alphabet ist das der Aktivierung für Schema und Tabelle (`identifierShape` in `tableactivation.go`); der Name geht in Fehleradressen und in eine Zerlegung `schema.table.name`, die dadurch eindeutig bleibt. |
| `column_name` bleibt bei beiden Transformations-Antragsarten NULL. | Die Spalte der Regel steht in `rule_spec`; `column_name` bleibt die Ziel-Spalte der beiden Spalten-Antragsarten (`SPEC-019`). Ein Antrag trägt so je Art genau eine Adressierungsform. |
| Prüfreihenfolge: Regelname, Form der `rule_spec`, Regeltyp, Schlüssel, Pflichtschlüssel und Werte, dann K1 bis K4; die Fehlertext-Wortlaute der Tabelle. | Der ADR-Text nennt die Verletzungen, keine Texte und keine Reihenfolge; die Form folgt dem Bestand (`ErrSourceColumnMissing`: Klartext, Doppelpunkt, Adresse). Die Reihenfolge trennt Form von Konfliktfreiheit, damit der erste Text eindeutig bestimmt ist. |
| Schlüsselreihenfolge im Image nicht zugesagt, zugesagt ist der Inhalt (Schlüsselmenge und Werte). | `jsonb` (`SPEC-002`) bewahrt die Reihenfolge nicht; über `cdc.changes` und `GET /changes` ist sie nicht beobachtbar. Der Assembler-Test darf sie zusätzlich prüfen, die Spec sagt sie nicht zu. Die Position am Assembler-Ausgang bleibt Festlegung der ADR. |
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
| Aufzählungen der Antragsarten | `grep -rn 'exclude_column' spec docs/user harness README.md` (Parent: `git grep -n 'exclude_column' 1718546b -- spec docs/user harness README.md`) | Parent: 10 Zeilen in 5 Dateien (`spec/architecture.md` 1, `spec/pflichtenheft.md` 5, `docs/user/benutzerhandbuch.md` 2, `harness/README.md` 1, `harness/targets/schema-rollout.md` 1). Diff (Stand `527db536`): 12 Zeilen, `spec/pflichtenheft.md` 7 statt 5 (die zwei zusätzlichen Zeilen sind der Satz in `LH-FA-CFG-007.a` und der Satz zur Spaltenauswahl im Absatz „Bezeichner" von `SPEC-030`; `docs/user`, `harness` und `spec/architecture.md` tragen an beiden Ständen dieselbe Zeilenzahl). | Spec-Zeilen gezogen (Antragsart-Tabelle der Architektur, `SPEC-019` Einleitung und `request_kind`-Menge). **Gemeldet, nicht mitgeändert:** `docs/user/benutzerhandbuch.md:265` (ein SQL-Beispiel des Ausschlusses, keine Aufzählung der Arten — Handbuch-Zug an `slice-transformationen-betriebsdoku`, Welle §4); `harness/targets/schema-rollout.md:74` („die fünf SQL-Funktionen“ beschreibt den Stand des Rollouts und bewegt sich mit `slice-transformationen-antragsweg-schema`); `harness/README.md:141` (Beschreibung des E2E-Laufs, bewegt sich mit den E2E-Slices). |
| Anzahl-Formulierungen zu Antragsarten und Werten der Menge | `grep -rn 'Antragsart' spec` (Parent: `git grep -n 'Antragsart' 1718546b -- spec`) und `grep -rn -i 'fünf Arten\|sieben Arten\|fünf übrigen\|drei übrigen\|sechs übrigen\|fünf Werte\|sieben Werte' spec` | Parent: 15 Zeilen mit „Antragsart“ (5 in der Architektur, 10 im Pflichtenheft); Zählwörter: „fünf Arten“ (`spec/architecture.md:248`), „die drei übrigen Antragsarten“ (`spec/pflichtenheft.md:535`), die `request_kind`-Menge mit fünf Werten (`:536`) und der Historien-Satz „(fünf Werte)“ (`:876`). Diff: 26 Zeilen; Zählwörter „sieben Arten“ (Architektur), „die fünf übrigen“/„die sechs übrigen“ (Spalten `column_name`/`rule_name`/`rule_spec`), die Menge mit sieben Werten; „fünf Werte“ steht im Diff ausschließlich in der Historie-Zeile vom 2026-09-24. | Zahlwörter und Aufzählungen auf die gemessene Menge gezogen (fünf Werte am Parent-Stand, `tools/schema/nacharbeit-administration.sql` Zeile 62 laut Welle §5; sieben nach diesem Diff); die Historie-Zeile bleibt (Protokoll, kein Zustand). Die Zählwörter „zehn Felder“ stehen nur in `SPEC-021` (`:621` am Parent) und `SPEC-024` (`:697` am Parent), unverändert; ein Zählwort „dreizehn“ steht im Pflichtenheft nicht (der Plan-Text dieses Slice nennt es für `SPEC-022`, gelesen, nicht nachgezählt, nicht geändert). |
| Fehlerklassen-Zeile `schema` in weiteren Trägern | `grep -rn 'nicht sicher interpretierbar' spec docs/user harness` (Parent: `git grep -n 'nicht sicher interpretierbar' 1718546b -- spec docs/user harness`) | Parent: 2 Zeilen (`spec/pflichtenheft.md:783`, `docs/user/benutzerhandbuch.md:1524`). Diff (Stand `527db536`): 2 Zeilen, `spec/pflichtenheft.md:1020` gezogen. | Spec-Zeile gezogen; Handbuch §6 (`docs/user/benutzerhandbuch.md`, Fehlerklassen-Tabelle) gemeldet an `slice-transformationen-betriebsdoku` |
| `LH-FA-CFG-007` in Trägern der Abdeckung | `grep -rn 'CFG-007' docs harness .d-check.yml` (Parent: `git grep -n 'CFG-007' 1718546b -- docs harness .d-check.yml`, ohne Planungs-/ADR-/Review-Records) | Parent: 1 Zeile außerhalb der Records, `harness/README.md:134` (`make doc-trace`-Zeile; sie lautet an `1718546b` und an `527db536` „80 Anforderungen, **2 Waisen**“, gemessen `git grep -o -n '[0-9]* Anforderungen, \*\*[0-9]* Waisen\*\*' <Stand> -- harness/README.md`). `docs/user/e2e-abdeckung.md`: kein Treffer. Diff: unverändert. Nachmessung `make doc-trace` am Diff-Stand: 80 Anforderungen, 2 Waisen (`LH-FA-CFG-007`, `LH-FA-CFG-008`). | `docs/user/e2e-abdeckung.md` unverändert bis `slice-transformationen-e2e-wirkung` (der Runner schreibt die Zeile). Die Zeile in `harness/README.md:134` deckt sich mit der Messung (2 Waisen): keine Drift, keine Meldung. |
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
| „Position"/„byte-gleich" der Row Images | Parent: 3 Zeilen (`spec/pflichtenheft.md:205`, `:299`, `:301`); Diff: 1 Zeile (`:205`) | `:299`/`:301` sind ersetzt. `:205` (`LH-FA-CAP-009.a` Absatz „Markierung") sagt „inhaltsgleich (Schlüsselmenge und Werte)" statt „byte-gleich": der Backfill-Stand gilt dem Assembler-Ausgang, am `jsonb`-Lesepfad ist die Reihenfolge nicht beobachtbar. Suchlauf der Closure: nächste Tabelle |
| Fehlertext-Tabelle und Prüfreihenfolge | Parent: „drei Formzeilen" (`:681`); Diff: „fünf Formzeilen" (`:683`), die neue Zeile `Regelname ist ungültig` (`:689`); ein weiteres Zählwort zu Formzeilen steht nirgends | gezogen; die Historie-Zeile vom 2026-09-26 zur ersten Fassung bleibt (Protokoll) |
| „Abhilfe: Regelstand ändern" in der `schema`-Zeile | Parent: 1 Zeile (`:963`); Diff: die Zeile trägt „Abhilfe nur bei nicht anwendbarer Regel" (`:1020`) | gezogen |
| Abhilfe an mehreren Stellen | Parent: `spec/pflichtenheft.md` 8 Zeilen (`:334`, `:337`, `:963`, `:970`, `:973`, `:976`, zwei Historie-Zeilen); Diff: 11 Zeilen, `LH-FA-CFG-007.a` führt (`:332`, `:335`), `SPEC-008` verweist (`:1020`, `:1027`, `:1030`); `docs/user` und `harness`: an beiden Ständen 3 Zeilen, keine gilt der Transformations-Abhilfe (WAL-Rückstand, Handbuch-Historie) | Führung deklariert (Absatz „Führende Stelle je Sachverhalt" oben); nichts außerhalb der Spec mitzuziehen |

**§3.13-Suchlauf der Closure** (bewegte Aussagen: „Position“/„byte-gleich“ der
Row Images, die Abhilfe im Run; Suchraum die offenen Pläne der Welle, die Welle
selbst und die Roadmap, Parent ist `527db536`, Diff ist der Arbeitsbaum der
Closure):

```text
git grep -c -E 'byte-gleich|Position seiner|Position der Quellspalte|Schlüsselposition' 527db536 -- docs/plan/planning/open docs/plan/planning/welle-transformationen.md docs/plan/planning/in-progress/roadmap.md
grep -rc -E 'byte-gleich|Position seiner|Position der Quellspalte|Schlüsselposition' docs/plan/planning/open docs/plan/planning/welle-transformationen.md docs/plan/planning/in-progress/roadmap.md
git grep -n -i -E 'WAL-Image|byte-gleich' 527db536 -- docs/user spec
```

| Bewegte Aussage | Befund Parent → Diff | Behandlung |
|---|---|---|
| „Position“/„byte-gleich“ in den offenen Plänen | Parent und Diff: dieselben 3 Dateien und 7 Zeilen (`slice-transformationen-kern-rename` 4, `-map-value` 2, `-backfill-pfad` 1); `welle-transformationen.md` und die Roadmap: 0 an beiden Ständen | die sieben Zeilen tragen im Diff den Bezug „am Ausgang des `Assembler`“ bzw. der gemeinsamen Bild-Konstruktion und den Satz, dass die Spec am Lesepfad Schlüsselmenge und Werte zusagt; die Zeilen selbst bleiben (Festlegung der ADR, prüfbar am Assembler) |
| „byte-gleich“ dem WAL-Image in Spec und Handbuch | Parent: `spec/pflichtenheft.md:205` und `:206` (der Satz umbricht, 2 Zeilen); Diff: `:206` (WAL-Image, ohne „byte-gleich“) und die Historie-Zeile `:1139`; `docs/user`: 0 Treffer an beiden Ständen | Spec-Satz gezogen (`inhaltsgleich`); das Handbuch trägt keine Aussage, keine Version-/Historie-Pflicht |

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
  'ADR-0\|slice-\|welle-'` über beide Spec-Dateien am Diff. **Ausgang:**
  *entfallen* — `git diff 1718546b..HEAD -- spec | grep '^+' | grep -c
  'ADR-0\|slice-\|welle-'` druckt `0` (Verifikation §2, Zeile 2), `make
  docs-check` Exit 0, die Mutation mit einem ADR-Link in der Spec färbt das
  `matrix`-Modul rot (§3, Belege des Laufs).
- **Zusagen als Tatsachen missverstanden**
  ([`AGENTS.md`](../../../../AGENTS.md) §3.12 Instanz B): die Abhilfe im
  gescheiterten Prozess ist in
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  „erwartet“, die Konfliktfreiheit K1–K4 und die Regeltypen sind Festlegungen.
  Das Pflichtenheft darf nur zusagen, was die Umsetzung liefern muss.
  *Erwartet, zu belegen durch:* Review liest jeden Satz der neuen Abschnitte
  auf Zukunfts-Form. **Ausgang:** *entfallen* — Review und Verifikation lasen
  die Abschnitte auf Zukunfts-Form (Verifikation §2, Zeile 1 und §6, Konjunktiv-
  und Chronik-Prüfung); die Abhilfe schließt mit „erst mit dem Beleg am
  laufenden System eine Tatsache“.
- **Überholter Text im selben Dokument**
  (`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`, verkörpert, 12
  Evidence-Dateien, gemessen am 2026-09-26 mit `ls
  docs/plan/planning/observations/BEO-PGC/<eintrag>/evidence | wc -l`): die
  Überschrift „… offen“, der Einleitungssatz von
  [`SPEC-019`](../../../../spec/pflichtenheft.md) (nennt die
  bestehenden SQL-Funktionen) und der Satz zur Menge der Arten in `spec/architecture.md` können
  die alte Aussage tragen. *Erwartet, zu belegen durch:* Lesen der Abschnitte
  von oben nach unten und Suchlauf §3, Zeilen 2–3. **Ausgang:** *eingetreten,
  im Slice behoben* — die Zeile `schema` mit der unbedingten Abhilfe neben
  Ursachen ohne Regelbezug (Review F-4), „Position“/„byte-gleich“ (F-2) und
  `LH-FA-CAP-009.a` Absatz „Markierung“ (Verifikation V-4) sind in der
  Fixrunde und der Closure gezogen; die Folge-Pläne mit denselben Wörtern
  (Verifikation V-3) sind qualifiziert. Kein Folge-Slice nötig. Klasse im
  Register: `BEO-PGC/nachzug-laesst-ueberholten-text-stehen` (Deckel, §7).
- **Kollision mit `slice-backfill-spec-nachzug`** in
  [`SPEC-019`](../../../../spec/pflichtenheft.md) und in der Kennungsvergabe.
  *Erwartet, zu belegen durch:* der Start-Trigger (der Backfill-Nachzug liegt
  in `done/`) und der Kennungs-Suchlauf am Parent-Stand. **Ausgang:**
  *entfallen* — der Start-Trigger war erfüllt, `SPEC-030` ist die nächste freie
  Kennung (`git grep -h -o 'SPEC-0[0-9][0-9]' 1718546b -- spec/pflichtenheft.md
  | sort -u | tail -1` druckt `SPEC-029`; Verifikation §2, Zeile 2).
- **Fehlertexte stehen in der Spec, bevor der Code sie trägt.** Der Wortlaut
  ist bindend für `slice-transformationen-antragsweg-usecase`; eine Abweichung
  dort ist ein Plan-Nachzug hierher, kein stilles Umformulieren. *Erwartet, zu
  belegen durch:* der Review von `antragsweg-usecase` liest die Texte gegen
  [`SPEC-019`](../../../../spec/pflichtenheft.md). **Ausgang:** *entfallen für
  diesen Slice* — die Tabelle trägt fünf Formzeilen und die Prüfreihenfolge
  (Review F-3 behoben); den Rest der Bindung trägt `antragsweg-usecase`: sein §2
  Punkt 1 bindet den Fehlertext, die Prüfreihenfolge und die
  Sentinel-Abbildungen an die Spec, sein Review liest die Texte gegen
  `SPEC-019`.
- **Die Randfälle der Regelform stehen nicht in
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)**
  (leeres `values`, Zielname gleich Quellname, Abbildung eines Werts auf sich
  selbst). Der Slice legt sie fest; berührt eine Festlegung eine ADR-Aussage,
  geht sie als Frage an den Architect. *Erwartet, zu belegen durch:* Review der
  Regelform gegen
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 2/3. **Ausgang:** *eingetreten, im Slice behoben* — die Randfälle
  sind in `SPEC-030` festgelegt und mit [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  konform (Verifikation §3); der Bezeichner-Vergleich fehlte (Review F-1,
  MEDIUM) und ist mit Bestands-Belegen festgelegt; eine Festlegung
  (Schlüsselreihenfolge) schwächt eine ADR-Aussage auf Spec-Ebene ab und ist in
  §3 benannt. Klasse im Register:
  `BEO-PGC/spec-nachzug-laesst-festlegung-fuer-folge-slice-offen` (neu, 1×).

## 7. Closure-Notiz

- **Was hat funktioniert:** der Zug blieb ein reiner Doku-Zug über die zwei
  Spec-Dateien ohne Rückführung (§4): Lastenheft, Benutzerhandbuch, Code und
  Schema sind unverändert (`git diff --stat`, Verifikation §2 Zeile 7). Die
  Rollen-Kette lief unabhängig: Review (0 HIGH · 1 MEDIUM · 5 LOW · 5 INFO, aus
  dem Report übernommen), Fixrunde für F-1 bis F-11 (`8e977328`, `26d6fde7`),
  Verifikation „Bestätigt" (acht `[x]`-Zeilen je mit eigenem Beleg; die Spec
  gegen
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  und [`ADR-0117`](../../adr/0117-backfill-run-fehlerklasse-schema.md)
  konform). Die Festlegungs-Tabelle in §3 (Festlegung · Begründung und
  Bestands-Beleg je Zeile) machte die Festlegungen über den ADR-Text hinaus
  prüfbar: der Verifier fuhr jede Zeile am Bestand nach (`containsName`,
  `ColumnExists`, `identifierShape`) und bestätigte sie. Die Zusagen-Form trug:
  die Abhilfe steht als Zusage an die Umsetzung, nicht als geprüfte Tatsache.
- **Was ging anders als geplant:** der Plan wuchs im Lauf um `SPEC-030` als
  eigene Kennung (Regelform als Datenstruktur, Punkt c der Zielformulierung),
  um Nachzüge in `LH-FA-CAP-009.a` und in der Sicht (Aufzählungen, die der
  Regelstand bewegt) und um die Fixrunde (Bezeichner-Vergleich, Beispiele,
  Führung je Sachverhalt). Der Verifier fand zwei Drifts im Suchlauf-Feld des
  Plans (V-1 falsche Träger-Aussage, V-2 zwei überholte Angaben) und eine zu
  enge Suche (V-3); alle sind in der Closure gezogen.
- **Verifier-Beobachtungen (V-1 bis V-6):** *V-1* (LOW) gezogen: die Zeile
  `harness/README.md:134` trägt an `1718546b` und an `527db536` „2 Waisen“; der
  Plan hatte „3 Waisen“ gemeldet, ein Wert aus einer Nachbardatei ohne
  Nachmessen an beiden Ständen. *V-2* (INFO) gezogen: 12 statt 11 Zeilen,
  Lokator `:1020` statt `:963`, mit Stand `527db536`. *V-3* (LOW) gezogen: die
  Folge-Pläne `kern-rename`, `map-value`, `backfill-pfad` tragen den Bezug „am
  Ausgang des `Assembler`“; der Vorspann der Festlegungen benennt die
  Abschwächung von
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  Teilfrage 3 auf Spec-Ebene. *V-4* (LOW) gezogen: `LH-FA-CAP-009.a`
  „Markierung“ sagt „inhaltsgleich (Schlüsselmenge und Werte)“; das Handbuch
  trägt keine Aussage. *V-5* (INFO) gezogen: `SPEC-008` nennt den Neustart für
  den neuen Run nicht als Teil der Zusage. *V-6* (INFO) ohne Handlungsbedarf:
  `LH-FA-CFG-007` und `LH-FA-CFG-008` bleiben Waisen in `make doc-trace`, bis
  die umsetzenden Folge-Slices und `e2e-wirkung` die Kennung führen.
- **Übergaben nach [`AGENTS.md`](../../../../AGENTS.md) §3.13 (Verifikation §7,
  Anschlussfähigkeit für `slice-transformationen-kern-rename`):** vier Lücken,
  in den offenen Plänen festgelegt, keine Umsetzung. (1) Zuordnung der
  Form-Prüfungen: die Domäne trägt `column`/`to` nichtleer ohne U+0000 und `to`
  höchstens 63 Byte (`kern-rename` §1), der Use Case bildet sie auf `rule_spec
  ist ungültig` ab, das Alphabet des Regelnamens und `Regelname ist ungültig`
  legt `antragsweg-usecase` fest. (2) Der Domänen-Sentinel für `to` gleich
  `column` wird im Use Case auf den K3-Text abgebildet (`kern-rename` §1,
  `antragsweg-usecase` §2 Punkt 1). (3) Zwischenzustand der Regeltyp-Menge:
  bis `map-value` endet ein `map_value`-Antrag `failed` (`antragsweg-usecase`
  §6, Ausgang mit der Closure von `map-value`). (4) Die Parameterlisten der
  SQL-Funktionen legt `antragsweg-schema` fest (§1).
- **Steering-Loop-Eintrag (Lerneintrag):** *neuer Register-Eintrag, benannte
  Spec-Lücke:* die Festlegungen, die
  [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md)
  offen ließ (Bezeichner-Vergleich, Randfälle der Regelform, Fehlertexte),
  traf in diesem Slice die Spec, je Zeile mit Bestandsbeleg gestützt;
  unvollständig blieben sie an den Rändern, die der erste Folge-Slice braucht.
  Die Frage „kann ein Implementer den ersten Folge-Slice aus der Spec allein
  umsetzen, und welche Festlegung fehlt je Schicht?“ fand vor dem Start von
  `slice-transformationen-kern-rename` vier Lücken, die weder das Review (liest
  gegen ADR und Plan) noch ein Gate fand. Der Eintrag
  `BEO-PGC/spec-nachzug-laesst-festlegung-fuer-folge-slice-offen` (1×, offen)
  trägt sie; bei 3× ist zu prüfen, ob die Anschlussfähigkeits-Frage Prüfpunkt
  des Verifier-Auftrags für Spec-Nachzug-Slices wird. Kein neuer Sensor: ob
  eine Festlegung fehlt, ist eine Lese-Handlung.
- **Beobachtungs-Register (`../observations/`):**
  - **`BEO-PGC/spec-nachzug-laesst-festlegung-fuer-folge-slice-offen`** — neu,
    `evidence/slice-transformationen-spec-nachzug.md` (Review F-1, F-3,
    Verifikation §7); Zähler **1×** (real ausgezählt).
  - **`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`** (verkörpert,
    Deckel bei 23×) — V-1 und V-2 (LOW, vor dem Merge vom Verifier gefunden,
    bekannter Träger-Typ Suchlauf-Feld): keine weitere Datei, das Auftreten
    steht hier mit Finding-Kennung; Zähler bleibt **23×**.
  - **`BEO-PGC/nachzug-laesst-ueberholten-text-stehen`** (verkörpert, Deckel bei
    12×) — V-3 (LOW, vor dem Merge vom Verifier gefunden, bekannter
    Träger-Typ: die Folge-Pläne fehlten im Suchraum) und V-4 (LOW): keine
    weitere Datei; Zähler bleibt **12×**.
  - **`BEO-PGC/adr-folgepflicht-ohne-traeger-slice`** (offen, 1×) — dieser Slice
    ist der Träger der Folgepflicht 1 von
    [`ADR-0112`](../../adr/0112-transformationsform-deklarative-regeln-vor-persistenz.md);
    die Pflicht ist eingelöst, ein zweites Auftreten liegt nicht vor.
  - Kein Eintrag steht bei 3× oder mehr ohne Ausgang an: ein Vermerk im
    Lese-Schritt der Closure von `welle-transformationen` entfällt.
- **Validator (Modul 8):** entfällt ausdrücklich — der Slice ist ein reiner
  Spec-Nachzug ohne End-Nutzer-Wert; der Nutzer-Bedarf
  ([`LH-FA-CFG-007`](../../../../spec/lastenheft.md)) wird erst durch die
  Umsetzung und den Wellen-Beleg (`e2e-wirkung`) validierbar.
- **Closure-Notiz-Review (`.harness/skills/closure-note-reviewer.md`):** eine
  getrennte Rolle im frischen Kontext; der Skill prüft Slices in `done/` und
  greift nach dem `git mv` — hier nicht ausgeführt.
- **Folge-Slices:** keine neuen — die Folge-Slices der Welle
  [welle-transformationen](../welle-transformationen.md) liegen als Dateien in
  `open/`; die Handbuch-Meldungen (Zeile `schema` in §6 des Handbuchs, das
  SQL-Beispiel des Ausschlusses) gehen an `slice-transformationen-betriebsdoku`,
  `harness/targets/schema-rollout.md:74` an
  `slice-transformationen-antragsweg-schema`, `harness/README.md:141` an die
  E2E-Slices.
- **Risiken aus §6:** je ein Ausgang am Ort — Risiko 1 (Spec-Straten ohne
  Kennung) **entfallen**; Risiko 2 (Zusagen als Tatsachen) **entfallen**;
  Risiko 3 (überholter Text) **eingetreten**, im Slice behoben; Risiko 4
  (Kollision mit dem Backfill-Nachzug) **entfallen**; Risiko 5 (Fehlertexte vor
  dem Code) **entfallen für diesen Slice**, Träger des Rests ist
  `antragsweg-usecase`; Risiko 6 (Randfälle) **eingetreten**, im Slice behoben,
  Klasse im Register.
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-transformationen](../welle-transformationen.md) (offen) — die Prüfung
  läuft regelkonform bei deren Closure. (a) Anker: der Lerneintrag verkörpert
  nichts neu, er benennt eine Spec-Lücke und legt einen Register-Eintrag an;
  (b) Folge-Slice: keiner neu, die Übergaben stehen in den fünf offenen Plänen;
  (c) Register: die genannten Kennungen existieren als Verzeichnis, das neue
  trägt seine Datei unter `evidence/`.

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
Träger der Folgepflicht 1). Die Closure legt
`BEO-PGC/spec-nachzug-laesst-festlegung-fuer-folge-slice-offen` an (offen, 1,
gemessen mit demselben Befehl); die Zahlen der übrigen Einträge sind mit der
Closure am Stand `527db536` nachgemessen und unverändert.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF.
