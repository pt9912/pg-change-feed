# BEO-PGC/rollen-verdrahtung

**Sub-Area:** Bootstrap/Verdrahtung (Least-Privilege-Rollen-Adoption;
Sub-Area-Kürzel `PGC` aus der Modus-Deklaration)

Die Beobachtung: Die drei Least-Privilege-Rollen (`cdc_capture`/
`cdc_admin`/`cdc_reader`, [`LH-QA-SEC-001`](../../../../../../spec/lastenheft.md)…003)
existieren als DDL-Artefakt und sind über `SET ROLE` einzeln testbar,
aber `internal/bootstrap/wiring.go` verdrahtet weiterhin eine
gemeinsame Instanz-DSN für Store-, Aktivierungs- und Stream-Verbindung.
Die Rollen-Adoption in der tatsächlichen Laufzeit-Verdrahtung ist noch
nicht erfolgt.

Deklaration: `internal/bootstrap/wiring.go`, slice-011-Plan §1/§6.
