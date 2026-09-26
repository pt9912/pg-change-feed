Zustand: **entschieden, Fix offen** (**2×**) — Hälfte geschlossen, Rest entschieden (Adresse `slice-antragsqueue-lesefehler-failed`). Geschlossen mit
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

**Entscheidung (Architect-Zug der Closure von `welle-transformationen`): Option (a), in der
allgemeinen Form.** Die Lesung der Queue lehnt keine einzelne Zeile ab: eine Zeile, deren
Antrags-Konstruktor sie verwirft (leeres Schema, leerer Tabellenname, leere Spalte, jeder
künftige Grund), wird als `failed` mit dem Fehlertext des Konstruktors vermerkt, und die Zeilen
dahinter werden verarbeitet. Das ist derselbe Mechanismus, der die Regelfelder schon schließt
(Stelle der Prüfung: Verarbeiten, nicht Lesen), und er entspricht der Spec-Aussage „`failed` mit
einem Fehlertext“. (b) scheidet aus: eine Prüfung in den SQL-Funktionen berührt `ADR-0046` (keine
Domänenlogik in SQL) und schützt nicht gegen einen direkten `INSERT` von `cdc_admin`, die Queue
bliebe für diese Zeile angehalten. (c) scheidet aus: eine Zeile einer vertrauten Rolle hielte den
Betrieb an, die Ursache stünde nur im Log, und die Abhilfe verlangte SQL-Zugriff auf die
Antrags-Tabelle. Der Fix braucht Code (Queue-Lesepfad und Use Case), keine neue ADR.

**Folgearbeit (Adresse `slice-antragsqueue-lesefehler-failed`, kleiner Fix-Slice):**
Lesepfad und Use Case so ändern, dass die verworfene Zeile mit Kennung und Fehlertext
durchgereicht und `failed` vermerkt wird; Store-Test mit realen Zeilen (leeres Schema, leerer
Tabellenname, `exclude_column`/`include_column` mit leerer Spalte, gültige Zeile dahinter wird
`applied`); den Test `TestReadPendingRequestsRejectsRowWithEmptySchemaOrTable` auf das neue
Verhalten umschreiben; `SPEC-019` im selben Commit nachziehen (Fehlertext je Feld, Ort der
Prüfung: Verarbeitung) — der heutige Wortlaut beschreibt das Ist-Verhalten und bleibt bis zum
Code stehen.

Zähler (abgeleitet): 2× (evidence/slice-transformationen-antragsweg-schema.md,
evidence/slice-transformationen-antragsweg-usecase.md).
