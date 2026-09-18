# BEO-PGC/spec008-replication-luecke

**Sub-Area:** Fehlerbehandlung/Capture-Pfad (Sub-Area-Kürzel `PGC` aus der
Modus-Deklaration)

Die Beobachtung: `SPEC-008` legt für die Fehlerklasse `replication` als
Aktion „Überwachung über Schwellen (§5, WAL-Rückstand); kontrollierte
Fortsetzung" fest. Der tatsächliche Code (`internal/bootstrap/wiring.go`,
`Run`) kennt keinen Schwellen-Zähler und keinen Fortsetzungspfad — ein
`replication`-klassifizierter Fehler propagiert wie jede andere Klasse
direkt bis zum Prozessende (`os.Exit(1)`). `docs/user/benutzerhandbuch.md`
§6 beschreibt damit den heutigen Code korrekt („Sichtbarer Fehler"); die
Diskrepanz liegt zwischen `SPEC-008`s Ziel-Zustand und der tatsächlichen
Implementierung, nicht zwischen Handbuch und Code.

Deklaration: Review zu `slice-020`, F-1,
Verifikationsbericht zu `slice-020` (unabhängig bestätigt).
