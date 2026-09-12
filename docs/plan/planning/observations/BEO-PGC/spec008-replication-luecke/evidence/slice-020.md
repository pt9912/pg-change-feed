# Beleg: slice-020 (Handbuch-Fehlerklassen vollständig)

Vorgang: slice-020.

Fund: Beim Vervollständigen der Handbuch-Fehlerklassen-Tabelle fand der
Reviewer (F-1) und bestätigte der Verifier unabhängig: Die
`replication`-Zeile des Handbuchs beschreibt „Sichtbarer Fehler"
(Prozessabbruch), `SPEC-008` verlangt für dieselbe Klasse „kontrollierte
Fortsetzung" nach Schwellen-Überwachung. Der Verifier las
`internal/bootstrap/wiring.go`s `Run`-Funktion vollständig: kein
Schwellen-Zähler, kein Fortsetzungspfad existiert — jeder Fehler,
inklusive `replication`, propagiert bis `os.Exit(1)`. Das Handbuch ist
korrekt (Ist-Zustand); `SPEC-008` beschreibt an dieser Zeile einen noch
nicht gebauten Ziel-Zustand.

**Ausgang: weiter offen.** Schwellen-Überwachung mit kontrollierter
Fortsetzung für `replication`-Fehler ist ein eigenständiges
Umsetzungs-Stück (WAL-Rückstand-Zähler, Retry-/Backoff-Logik im
Capture-Pfad) oder erfordert eine Korrektur von `SPEC-008` selbst, falls
das Ziel nicht mehr verfolgt wird — beides eine bewusste Entscheidung,
kein Ein-Zeilen-Fix.

Quelle: `docs/reviews/review-slice-020.md` F-1, `docs/reviews/verify-slice-020.md`.
