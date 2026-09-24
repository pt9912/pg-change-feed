# BEO-PGC/backfill-schema-version-hinter-snapshot-spalten

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Schema-Version-Referenz der
Backfill-Changes, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Der Run liest die aktuelle Schema-Version der Tabelle (`CurrentVersion`) vor dem
Öffnen des Snapshots. Die Version wechselt allein mit der nächsten Relation-Nachricht des
WAL-Pfads. Sie kann deshalb hinter den Spalten des Snapshots liegen (kompatible
Spalten-Erweiterung ohne WAL-Änderung seit der Aktivierung) oder auf eine Version ohne
`TableSchema` verweisen (statische Erstaktivierung, Nachtrag erst mit der ersten
Relation-Nachricht). Das akzeptierte Negativ in `ADR-0111` nennt nur die Erweiterung „während
des Runs"; die Fälle davor stehen dort nicht, das Argument der ADR (das Bild ist
selbstbeschreibendes JSON) trägt sie ebenso. Die Unterscheidbarkeit der Versionen
(`LH-FA-SCH-005`) bleibt erfüllt. Die Grenze steht weder im Code noch in `ADR-0111`.

Deklaration: `slice-backfill-run-usecase`, Risiko §6 („An Architect gemeldet:
`schema_version` der Backfill-Changes"), Ausgang *weiter offen*; Review F-8 (INFO).
