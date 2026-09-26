// Package errors trägt die Domänen-Fehler der Invarianten-Verletzungen
// (`ADR-0029`). Die Konstruktoren in `internal/domain/model` liefern sie
// als Fehlerursache; Adapter übersetzen sie an ihrer Grenze in die
// Fehlerklassen des Pflichtenhefts (§4, `SPEC-008`).
package errors

import (
	stderrors "errors"
)

// Jede Invariante aus `ADR-0029`, die ein Domänentyp erzwingt, hat hier
// einen Sentinel: die Konstruktoren wickeln ihn ein (`%w`), der Aufrufer
// klassifiziert über `errors.Is`.
var (
	// ErrEmptyIdentifier: Kennungen (Source, Tabelle, Change,
	// Transaktion, Consumer, Schema-Version) sind nicht leer.
	ErrEmptyIdentifier = stderrors.New("leere Kennung")

	// ErrInvalidOperation: Operation ist INSERT, UPDATE oder DELETE
	// (`SPEC-002`).
	ErrInvalidOperation = stderrors.New("unbekannter Operationstyp")

	// ErrNonPositiveSequence: die Sequenz innerhalb der Transaktion ist
	// mindestens 1 (`SPEC-002`, eindeutige Sequenz).
	ErrNonPositiveSequence = stderrors.New("Sequenz ist kleiner als 1")

	// ErrNonPositiveVersion: die Versionsnummer einer Schema-Version ist
	// mindestens 1 (`SPEC-004`).
	ErrNonPositiveVersion = stderrors.New("Versionsnummer ist kleiner als 1")

	// ErrInvalidPosition: eine Quellposition trägt einen Offset größer
	// als 0 (`SPEC-003`, sortierbare Position `LH-FA-DAT-004`).
	ErrInvalidPosition = stderrors.New("Position ohne Offset")

	// ErrPositionRegression: der Consumer-ACK verläuft regulär nur
	// vorwärts (`ADR-0029`, Regel 2); Wiederholung derselben Position ist
	// idempotent (`LH-FA-CON-004`).
	ErrPositionRegression = stderrors.New("Position rückt nicht vor")

	// ErrSourceMismatch: Positionen sind nur innerhalb ihrer Quelle
	// fachlich geordnet (`ADR-0005`); eine Position einer anderen Quelle
	// ist keine Bestätigungs- oder Commit-Größe für diese.
	ErrSourceMismatch = stderrors.New("Position einer anderen Quelle")

	// ErrTransactionMismatch: ein Change gehört genau einer Transaktion
	// an (`ADR-0029`, Regel 6).
	ErrTransactionMismatch = stderrors.New("Change gehört nicht zu dieser Transaktion")

	// ErrDuplicateSequence: die Sequenz ist innerhalb der Transaktion
	// eindeutig (`ADR-0029`, Regel 6).
	ErrDuplicateSequence = stderrors.New("Sequenz innerhalb der Transaktion bereits vergeben")

	// ErrTransactionAlreadyCommitted: ein Commit ist einmalig; einer
	// committed Transaktion fügt keine Changes mehr hinzu (`ADR-0029`,
	// Regel 3).
	ErrTransactionAlreadyCommitted = stderrors.New("Transaktion ist bereits committed")

	// ErrTransactionNotCommitted: offene Transaktionen sind nicht
	// konsumierbar (`ADR-0029`, Regel 3); ihre Changes sind erst nach dem
	// Commit lesbar (`LH-FA-CAP-006`).
	ErrTransactionNotCommitted = stderrors.New("Transaktion ist nicht committed")

	// ErrNegativeDuration: Zeiträume der Retention sind nicht negativ
	// (`LH-FA-RET-003`).
	ErrNegativeDuration = stderrors.New("negative Dauer")

	// ErrInvalidErrorClass: eine Fehlerklasse ist eine der sieben stabilen
	// Kategorien aus `ADR-0023` (`SPEC-008`) — eine leere oder unbekannte
	// Klasse verletzt die Invariante (`LH-FA-ADM-003`).
	ErrInvalidErrorClass = stderrors.New("unbekannte Fehlerklasse")

	// ErrEmptyColumns: ein TableSchema trägt mindestens eine Spalte
	// (`SPEC-004`, `ADR-0015` Folgepflicht) — eine Schema-Version ohne
	// Spaltenform trägt keine historisch stabile Interpretation.
	ErrEmptyColumns = stderrors.New("TableSchema ohne Spalten")

	// ErrInvalidAdministrationRequestKind: die Antragsart eines
	// Administrations-Antrags ist eine der geschlossenen Menge
	// `enable`/`disable`/`exclude_column`/`include_column`/`backfill`/
	// `set_transformation`/`remove_transformation`
	// (`chk_administration_request_kind`, `ADR-0050`, `LH-FA-CFG-005`,
	// `LH-FA-CAP-009`, `LH-FA-CFG-007`).
	ErrInvalidAdministrationRequestKind = stderrors.New("unbekannte Antragsart")

	// ErrInvalidChangeOrigin: die Herkunft eines Changes ist `wal` oder
	// `backfill` (`SPEC-002`, `LH-FA-CAP-009`); die Datenbankspalte
	// `cdc.change.origin` trägt keinen CHECK, die geschlossene Menge
	// erzwingt allein die Domäne.
	ErrInvalidChangeOrigin = stderrors.New("unbekannte Change-Herkunft")

	// ErrInvalidBackfillTransition: ein Backfill-Run wechselt nur entlang
	// `queued` → `running` → `completed`/`interrupted` und `queued`/`running`
	// → `failed` (`SPEC-029`, `LH-FA-CAP-009`); ein beendeter Run ist
	// endgültig, ein neuer Antrag legt einen neuen Run an.
	ErrInvalidBackfillTransition = stderrors.New("unzulässiger Statuswechsel des Backfill-Runs")

	// ErrBackfillProgressRegression: der Fortschrittszähler eines Runs
	// (`rows_copied`, `SPEC-029`) wächst nur.
	ErrBackfillProgressRegression = stderrors.New("Fortschritt des Backfill-Runs rückt nicht vor")

	// ErrNegativeRowCount: eine geschätzte Zeilenzahl ist mindestens 0;
	// „unbekannt" ist ein eigener Zustand, kein negativer Wert (`SPEC-029`).
	ErrNegativeRowCount = stderrors.New("Zeilenzahl ist kleiner als 0")

	// ErrBackfillBlockOverflow: die Blocknummer der synthetischen
	// Transaktions-Kennung ist achtstellig (`ADR-0111` Teilfrage 6); eine
	// größere Nummer verletzt die lexikographische Ordnung der Kennungen.
	ErrBackfillBlockOverflow = stderrors.New("Blocknummer außerhalb von 1 bis 99999999")

	// ErrTableNotActivated: die Tabelle trägt keine Bindung oder keine
	// Mitgliedschaft in der Publication der Quelle — die Vorbedingung eines
	// Backfills (`LH-FA-CAP-009`, `ADR-0111` Teilfrage 5) und der Fail-closed-
	// Prüfung vor dem Commit (`ADR-0111` Teilfrage 4). Fehlerklasse
	// `configuration` (`SPEC-008`).
	ErrTableNotActivated = stderrors.New("Tabelle nicht aktiviert oder nicht in der Publication")

	// ErrBackfillRunActive: für dieselbe Tabelle besteht ein Run im Zustand
	// `queued` oder `running` (`ADR-0111` Teilfrage 4, `ADR-0113`
	// Festlegung 1); der Annahme-Port meldet ihn, es entsteht keine zweite
	// Run-Zeile.
	ErrBackfillRunActive = stderrors.New("Für die Tabelle besteht bereits ein aktiver Backfill-Run")

	// ErrExclusionStateChanged: der Ausschlussstand der Tabelle
	// (`LH-FA-CFG-005`) weicht von dem ab, mit dem die Blöcke eines Runs
	// gebaut wurden (`ADR-0111` Teilfrage 4, `LH-QA-SEC-004`). Fehlerklasse
	// `configuration` (`SPEC-008`).
	ErrExclusionStateChanged = stderrors.New("Ausschlussstand während des Backfills geändert")

	// ErrTransformationStateChanged: der Regelstand der Tabelle
	// (`LH-FA-CFG-007`) weicht von dem ab, mit dem die Blöcke eines Runs
	// gebaut wurden (`ADR-0117` Festlegung 5, `ADR-0111` Teilfrage 4).
	// Fehlerklasse `configuration` (`SPEC-008`): der Zustand wechselt, keine
	// Regel ist auf eine Form nicht anwendbar.
	ErrTransformationStateChanged = stderrors.New("Regelstand während des Backfills geändert")

	// ErrInvalidTransformation: eine Transformationsregel verletzt die
	// Invarianten ihres Regeltyps — `column` oder `to` leer, mit dem
	// Zeichen U+0000, oder `to` länger als 63 Byte in UTF-8 (`SPEC-030`,
	// Bezeichner, `LH-FA-CFG-007`).
	ErrInvalidTransformation = stderrors.New("ungültige Transformationsregel")

	// ErrTransformationTargetIsColumn: der Zielname von `rename_column`
	// gleicht der Quellspalte; die Quellspalte ist ein Spaltenname der
	// Quelltabelle, der Fall gehört zu K3 (`SPEC-019`, `SPEC-030`
	// Randfälle), nicht zur Form der Regel.
	ErrTransformationTargetIsColumn = stderrors.New("Zielname gleicht der Quellspalte")

	// ErrTransformationColumnMissing: die Spalte einer Regel kommt in den
	// Spalten der Änderung nicht vor — die Regel ist nicht anwendbar
	// (`SPEC-030`, Anwendbarkeit).
	ErrTransformationColumnMissing = stderrors.New("Spalte der Regel fehlt in den Spalten der Änderung")

	// ErrTransformationTargetCollides: der Zielname einer Regel gleicht
	// einer Spalte der Änderung — die Regel ist nicht anwendbar, das Bild
	// trüge zwei gleichnamige Schlüssel (`SPEC-030`, Anwendbarkeit).
	ErrTransformationTargetCollides = stderrors.New("Zielname kollidiert mit einer Spalte der Änderung")

	// Die folgenden Sentinels tragen die Ablehnungsgründe eines
	// Regel-Antrags (`SPEC-019`, Fehlertext-Tabelle): ihr Text ist der
	// Klartext der Zeile; der Use Case hängt die Adresse an. Die fünf ersten
	// sind die Formzeilen, die vier letzten die Konfliktfreiheit K1 bis K4
	// (K4 trägt `inbound.ErrSourceColumnMissing` und `ErrRuleNotKept`).

	// ErrInvalidRuleName: der Regelname ist leer oder liegt außerhalb des
	// Alphabets `a`–`z`, `0`–`9`, `_` mit 1 bis 63 Zeichen (`SPEC-030`,
	// Bezeichner).
	ErrInvalidRuleName = stderrors.New("Regelname ist ungültig")

	// ErrInvalidRuleSpec: `rule_spec` ist kein JSON-Objekt mit einem
	// Zeichenketten-`kind`, oder ein Pflichtschlüssel des Regeltyps fehlt
	// oder hat den falschen Typ, oder `column`/`to` verletzt die
	// Bezeichner-Form (`SPEC-030`).
	ErrInvalidRuleSpec = stderrors.New("rule_spec ist ungültig")

	// ErrUnknownTransformationKind: `kind` nennt keinen Regeltyp aus
	// `TransformationKinds` (`SPEC-030`).
	ErrUnknownTransformationKind = stderrors.New("unbekannter Regeltyp")

	// ErrUnknownRuleSpecKey: `rule_spec` trägt einen Schlüssel, den der
	// Regeltyp nicht kennt (`SPEC-030`).
	ErrUnknownRuleSpecKey = stderrors.New("unbekannter Schlüssel in rule_spec")

	// ErrRuleNameTaken (K1): der Regelname ist je Tabelle vergeben
	// (`SPEC-019`).
	ErrRuleNameTaken = stderrors.New("Regelname bereits vergeben")

	// ErrColumnHasRule (K2): die Quellspalte trägt bereits eine
	// Spaltenregel (`SPEC-019`).
	ErrColumnHasRule = stderrors.New("Spalte trägt bereits eine Regel")

	// ErrTargetCollidesWithRule (K3): der Zielname gleicht dem Zielnamen
	// einer anderen Regel der Tabelle (`SPEC-019`).
	ErrTargetCollidesWithRule = stderrors.New("Zielname kollidiert mit einer anderen Regel")

	// ErrTargetCollidesWithColumn (K3): der Zielname gleicht einem
	// Spaltennamen der Quelltabelle, der Quellspalte der Regel
	// eingeschlossen (`SPEC-019`).
	ErrTargetCollidesWithColumn = stderrors.New("Zielname kollidiert mit einer Spalte der Tabelle")

	// ErrRuleNotKept (K4): ein `remove_transformation` nennt einen
	// Regelnamen, den die Tabelle nicht führt (`SPEC-019`).
	ErrRuleNotKept = stderrors.New("Regelname nicht geführt")
)
