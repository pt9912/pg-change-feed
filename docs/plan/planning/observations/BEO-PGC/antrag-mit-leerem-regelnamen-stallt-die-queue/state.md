Zustand: **offen** (**2×**) — Hälfte geschlossen, Rest offen. Geschlossen mit
`slice-transformationen-antragsweg-usecase`: die Regelfelder (`rule_name`, `rule_spec`) werden
gelesen und verarbeitet statt beim Lesen abgelehnt; die Zeile endet `failed` mit dem Fehlertext
der Spec, die gültige Zeile dahinter wird `applied` (Store-Test mit realen Zeilen,
Whitebox-Test der Queue, Login-Test unter `cdc_admin`/`cdc_capture`); Pfad des Slice:
`docs/plan/planning/done/slice-transformationen-antragsweg-usecase.md` (§3 „Stelle der
Prüfung“).

**Rest (Review F-6, zweites Auftreten):** der Antrags-Konstruktor lehnt beim Lesen weiter ab,
und ein solcher Antrag hält die gesamte Queue an — `exclude_column`/`include_column` mit leerer
Spalte und jede Antragsart mit leerem Schema oder leerem Tabellennamen (`ErrEmptyIdentifier`;
Ist-Verhalten erprobt mit `TestReadPendingRequestsRejectsRowWithEmptySchemaOrTable`; ob die
SQL-Funktionen einen solchen Antrag schreiben, ist hergeleitet, nicht erprobt). Abhilfe: die
Zeile per `UPDATE` auf `failed` vermerken (`cdc_admin` trägt `UPDATE` auf der Antrags-Tabelle;
hergeleitet aus den Grants, nicht erprobt).

**Benannte Spec-Lücke:** `SPEC-019` nennt für ein leeres Schema, einen leeren Tabellennamen und
eine leere Spalte weder einen Fehlertext noch den Ort der Prüfung; die Zeile zur Spalte nennt den
Konstruktor (das Ist-Verhalten), sagt aber nicht, dass die Lesung dort endet.

**Adresse:** der Architect-Zug der Closure von `welle-transformationen` (Schritt 3b,
Planner → Architect → Planner) als Vorschlag des Planners, **unabhängig vom Zähler** (bei 2×
liest der Lese-Schritt den Eintrag nicht; eine angehaltene Queue trifft aber den Betrieb, und der
Architect ist dort ohnehin beauftragt): der Architect entscheidet zwischen (a) Verarbeiten mit
`failed`-Vermerk und Fehlertext je Feld samt Spec-Zeile (ein Folge-Slice), (b) einer Prüfung in
den SQL-Funktionen (Berührung von `ADR-0046`, keine Domänenlogik in SQL) und (c) einem
akzeptierten Negativ mit der Abhilfe im Handbuch (Träger dann ein
Übergabe-Text im Plan von `slice-transformationen-betriebsdoku`, den der Planner in diesem Fall
anlegt; heute trägt der Plan die Aussage nicht). Bis zur Entscheidung steht kein Ausgang.

Zähler (abgeleitet): 2× (evidence/slice-transformationen-antragsweg-schema.md,
evidence/slice-transformationen-antragsweg-usecase.md).
