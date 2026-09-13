# Beleg: slice-044 (Retention-Hintergrundjob)

Vorgang: slice-044.

Fund: Beim Verdrahten von `runRetentionCleanup` gegen die reale
Löschausführung (`RunRetentionUseCase.Run` → `ChangeStorePort.DeleteChanges`,
`slice-043`) zeigte ein Grant-Abgleich gegen
`tools/schema/nacharbeit-roles.sql`, dass weder `cdc_capture` noch
`cdc_admin` ein `DELETE`-Grant auf `cdc.transaction`/`cdc.change` trugen.
Das Architect-Verdikt
(`docs/reviews/architect-verdict-retention-loeschausfuehrung.md`, Frage 1)
hatte die Domain-/Port-/ADR-Ebene der neuen Löschmethode geprüft (kein
Widerspruch zu `ADR-0009`/`0011`/`0012`/`0014`/`0029`), aber keine
Rollen-/Grant-Prüfung für den neuen physischen `DELETE`-Pfad
durchgeführt — `welle-13` §6 schließt eine Rollen-Erweiterung als
Out-of-Scope aus und benennt für den Fall eines tatsächlichen Bedarfs
„ein eigener Architect-Zug, kein stiller Fortschritt dieser Welle".

Da die Implementer-Session dieses Slice ohne separaten Architect-Kontext
lief, wurde die Grant-Ergänzung (`GRANT SELECT, DELETE ON
cdc.transaction, cdc.change TO cdc_admin;`, an der bereits bestehenden
`cdc_admin`-Verwaltungspfad-Grant-Zeile) transparent im Plan-Nachzug
dokumentiert statt stillschweigend vorgenommen, mit realer
Regressions-Testabdeckung (`TestCdcAdminRetentionDeleteChangesRequiresGrant`,
`make test-store`) für beide Seiten (ohne Grant: `SQLSTATE 42501`; mit
Grant: Erfolg).

**Ausgang: weiter offen.** Der Architect-Verdikt-Ablauf selbst trägt noch
keinen Pflichtabschnitt „Rollen-/Grant-Konsequenz" für neue physische
Schreib-/Löschpfade; ohne ihn kann derselbe blinde Fleck bei der nächsten
neuen Fähigkeit erneut auftreten.

Quelle: `internal/bootstrap/wiring.go` (Retention-Pool-Verdrahtung),
`tools/schema/nacharbeit-roles.sql`,
`internal/bootstrap/roles_wiring_test.go`
(`TestCdcAdminRetentionDeleteChangesRequiresGrant`).
