# BEO-PGC/verwaltung-keine-sql-administration

**Sub-Area:** CDC-Verwaltung/Administration (`LH-FA-CFG-*`, `LH-FA-ADM-*`;
Sub-Area-Kürzel `PGC` aus der Modus-Deklaration)

Die Beobachtung: [`LH-FA-ADM-001`](../../../../../../spec/lastenheft.md)
(Lastenheft, vertraglich abnahmebindend) verlangt in seinem Happy-Path-
Akzeptanzkriterium ausdrücklich, dass ein Administrator „Aktivierung,
Deaktivierung, Statusabfrage und Consumer-Verwaltung über SQL ausführen"
kann (Beispiele im Boundary-Kriterium: `cdc.enable_table(...)`,
`cdc.disable_table(...)`, `SELECT * FROM cdc.tables`, `SELECT * FROM
cdc.consumers`). Real existieren nur die Lese-Views
(`cdc.active_tables`, `cdc.consumer_status`) — keine SQL-Funktionen für
Aktivierung/Deaktivierung. Zusätzlich hat
[`LH-FA-CFG-002`](../../../../../../spec/lastenheft.md) (CDC-Deaktivierung)
im gesamten Repo **keinen** Live-Zugriffsweg: `disable.NewDisableTableService`
wird in `internal/bootstrap/wiring.go` nirgends konstruiert; der einzige
Aufrufer ist der weißbox-Unit-/Integrationstest selbst. Derselbe
Befund gilt für die CLI-Seite:
[`LH-FA-SST-003`](../../../../../../spec/lastenheft.md) verlangt, dass eine
CLI mindestens die Status-/Diagnoseabfragen aus
[`LH-FA-ADM-002`](../../../../../../spec/lastenheft.md)…`005` abdeckt — real
existieren nur die `register-consumer`/`acknowledge-consumer`-Befehle,
kein Status-/Diagnose-Befehl.

## Benannt, nicht gezählt

Gefunden bei einer Fork-Recherche zur Eröffnung von `welle-11`
(E2E-Abdeckung — Verwaltung & Observability), kein abgeschlossener
Vorgang trägt bisher einen Beleg.
