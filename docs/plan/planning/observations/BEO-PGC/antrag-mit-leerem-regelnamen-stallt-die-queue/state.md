Zustand: **verkörpert** — Ausgang: **verkörpert** → Fix geliefert: `sqlexec.ReadPendingRequests`
(`internal/adapters/driven/postgresstorage/sqlexec/translate.go`) lehnt keine Zeile ab und reicht eine
vom Antrags-Konstruktor verworfene Zeile mit Kennung und Klartext durch; `processAdministrationRequests`
(`internal/bootstrap/wiring.go`) vermerkt sie `failed` mit dem Text in `error_message` und verarbeitet
die Zeilen dahinter in der Ordnung der Queue; `SPEC-019`, Absatz „Zeilen, die kein Antrag sind“, trägt
Ort der Prüfung (Verarbeitung, nicht Lesen), Klartext je Grund und die Grenze „Zeile ohne Kennung“ ·
seit slice-antragsqueue-lesefehler-failed (Beleg-Anker: `git grep -n 'Zeilen, die kein Antrag sind' --
spec/pflichtenheft.md` und `git grep -n 'func rejectionMessage' -- internal`). Die Regelfelder
(`rule_name`, `rule_spec`) schloss zuvor `slice-transformationen-antragsweg-usecase` (Prüfung in der
Verarbeitung, Fehlertext der Spec); der Fix erfasst jeden Grund, den der Konstruktor kennt und den eine
Kennung adressiert: leere Quelle, leeres Schema, leerer Tabellenname, unbekannte Antragsart, leere Spalte von
`exclude_column`/`include_column`; eine Zeile ohne Kennung ist die benannte Grenze (unten).

**Entscheidung: Option (a) in der allgemeinen Form.** Die Lesung der Queue lehnt keine einzelne Zeile ab;
die Prüfung liegt in der Verarbeitung (derselbe Ort wie bei den Regelfeldern). (b) scheidet aus: eine Prüfung
in den SQL-Funktionen berührt `ADR-0046` (keine Domänenlogik in SQL) und schützt nicht gegen einen
direkten `INSERT` von `cdc_admin`. (c) scheidet aus: eine Zeile einer vertrauten Rolle hielte den Betrieb
an, die Ursache stünde nur im Log. Der Fix ist Code, keine ADR.

**Belege:** Store-Test mit realen Zeilen `TestAdministrationRequestListPendingPassesRejectedRowsThrough`
(vier Gründe, je eine gültige Zeile davor und dahinter, Vermerke `failed`/`applied`), Whitebox-Tests der
Verarbeitung im Paket `internal/bootstrap`, Login-Test `TestAdministrationPathRunsUnderLeastPrivilegeLogins`
unter `cdc_admin`; Mutationen der Eingabeseite: Review 14 von 14 rot, Verifikation 21 von 23 rot (zwei
grün, äquivalent).

**Benannte Grenzen (Kenntnis):** `sqlexec.rejectionMessage` bildet den Klartext aus den Feldern der Zeile;
ein neuer Konstruktor-Grund trägt bis zu einem eigenen Fall den allgemeinen Klartext `Antrag ist
ungültig` (Adresse der Kopplung: der Kommentar an `rejectionMessage`). Die Gründe „Quelle ist leer“ und
„Antragsart ist unbekannt“ entstehen über die SQL-Funktionen nicht (Fremdschlüssel, `CHECK`) und sind nur
an Fake-Tests belegt. Eine Zeile ohne Kennung bleibt `pending` und erzeugt eine Warnung je Durchlauf.

Zähler (abgeleitet): 2× (evidence/slice-transformationen-antragsweg-schema.md,
evidence/slice-transformationen-antragsweg-usecase.md); der Slice des Fixes trägt keinen dritten Beleg,
Review und Verifikation fanden kein weiteres Auftreten der Klasse.
