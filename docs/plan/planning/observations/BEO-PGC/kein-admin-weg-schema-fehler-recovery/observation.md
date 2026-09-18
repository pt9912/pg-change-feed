# BEO-PGC/kein-admin-weg-schema-fehler-recovery

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft
Betriebsdokumentation/Administrationswege nach einem `schema`-Fehlerklasse-
Abbruch, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Eine nicht sicher interpretierbare Relation-Änderung
(`ErrIncompatibleSchemaChange`, Fehlerklasse `schema`) beendet den
Erfassungspfad dauerhaft (`os.Exit(1)`, `restart: "no"`) —
`LH-FA-SCH-004`s Negative-Kriterium fordert genau das. Es gibt aber
**keinen dokumentierten Administrationsweg**, mit dem ein Betreiber den
Prozess nach einem solchen Vorfall real wieder in Betrieb nimmt (neue
Schema-Version registrieren, Replication-Slot-Zustand bereinigen). Der
einzige real erprobte Weg ist der Testharness-interne Mechanismus aus
`slice-062` (direkte SQL-Eingriffe in `cdc.schema_version`/
`cdc.table_schema` plus Slot-Neuanlage, `tools/harness/run-integration-tests.sh`)
— nicht als Betriebsdokumentation für Endnutzer gedacht oder geeignet.

Deklaration: `slice-062` (Review-Finding F-1, Review zu `slice-062`;
Verifier bestätigte die Einschätzung unabhängig,
Verifikationsbericht zu `slice-062`).
