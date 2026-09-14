**Vorgang:** slice-062
**Fund:** Der Reviewer stellte beim Bauen des Recovery-Mechanismus für den
E2E-Testrundlauf fest, dass die einzige real erprobte Methode, einen nach
einem `schema`-Fehler dauerhaft beendeten Feed-Prozess wieder in Betrieb
zu nehmen (direkte SQL-Eingriffe in `cdc.schema_version`/
`cdc.table_schema` plus Replication-Slot-Neuanlage), ausschließlich im
Testharness existiert — kein dokumentierter oder für Endnutzer geeigneter
Administrationsweg. Der Verifier bestätigte die Einschätzung unabhängig.
Erstauftreten.
