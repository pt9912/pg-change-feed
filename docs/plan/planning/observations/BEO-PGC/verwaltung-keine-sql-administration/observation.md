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

Zweiter Fund (Fork-Recherche zur Eröffnung der Feature-Welle
„Verwaltungsfunktionen — SQL-Administration & CLI-Diagnose",
2026-09-13): [`ADR-0046`](../../../../../plan/adr/README.md) (Accepted)
hat die Frage, wie eine schreibende SQL-Funktion
(`cdc.enable_table(...)`) einen Go-Inbound-Port erreichen soll, bereits
bewusst offen gelassen („dasselbe physikalische Problem wie bei den
Views … bleibt offen und ist nicht Gegenstand dieser Entscheidung").
Zusätzlich baut `Assembler.tables`
(`internal/adapters/driving/replication/mapper/mapper.go`, verdrahtet in
`internal/bootstrap/wiring.go`) seine Tabellen-Bindungen einmalig beim
Container-Start aus `CDC_TABLES` — ohne Reload-Mechanismus. Selbst eine
rein SQL-seitige Aktivierung (Bindungs-Zeile + `ALTER PUBLICATION`)
würde den laufenden Erfassungspfad nicht erreichen. Beides ist eine
Architect-Entscheidung wert, bevor eine Umsetzung geschnitten wird.
