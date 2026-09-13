# BEO-PGC/architect-verdikt-rollen-scope-luecke

**Sub-Area:** Bootstrap/Rollen-Verdrahtung (Sub-Area-Kürzel `PGC` aus der
Modus-Deklaration)

Die Beobachtung: Ein Architect-Verdikt, das eine neue Fähigkeit gegen
bestehende ADRs prüft (Domain-/Port-/Use-Case-Ebene), prüft nicht
durchgängig auch die Rollen-/Grant-Konsequenz eines neuen physischen
Schreib- oder Löschpfads auf PostgreSQL-Ebene (`ADR-0047`). Konkretes
Beispiel: `docs/reviews/architect-verdict-retention-loeschausfuehrung.md`
Frage 1 prüft `ADR-0009`/`0011`/`0012`/`0014`/`0029` und stellt fest,
dass die neue `ChangeStorePort`-Löschmethode keiner dieser Entscheidungen
widerspricht — die Frage, ob die Rolle, über die dieser physische
`DELETE` in der Verdrahtung real läuft (`cdc_admin`, `CDC_ADMIN_DSN`),
das dafür nötige `DELETE`-Grant auf `cdc.transaction`/`cdc.change`
überhaupt trägt, bleibt im Verdikt unbehandelt. `welle-13` §6 schließt
eine Rollen-Erweiterung explizit als Out-of-Scope aus, mit der
Bedingung, dass ein tatsächlicher Bedarf „ein eigener Architect-Zug,
kein stiller Fortschritt dieser Welle" wäre — real geprüft (`grep` gegen
`tools/schema/nacharbeit-roles.sql`) trug weder `cdc_capture` noch
`cdc_admin` vor `slice-044` ein `DELETE`-Grant auf diesen beiden
Tabellen.

Deklaration: `docs/reviews/architect-verdict-retention-loeschausfuehrung.md`
(Frage 1, keine Rollen-/Grant-Prüfung), `slice-044`-Plan-Nachzug §3
(Fund, Grant-Ergänzung ohne separaten Architect-Zug in eigenem Kontext —
Implementer-Session ohne getrennte Architect-Rolle verfügbar).
