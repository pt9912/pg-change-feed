# Beleg: slice-090

Vorgang: `slice-090` — das Sync-Gate des generierten Protobuf-Codes.

Fund: In der Verifikation zu `slice-090` **V-1** lieferte derselbe Befehl am
**selben** Stand zwei verschiedene Coverage-Zahlen — `make gates` **74,80 %**,
`make coverage-gate` zweimal **74,70 %**. Ursache, lokalisiert: `go tool cover
-func` gibt für `runWALRetentionCheck` (`internal/bootstrap/wiring.go:981`) in
**7 von 8** Läufen **87,5 %**, in **einem** **100,0 %** — die Erreichung einer
Anweisung im `select`/Ticker-Pfad hängt von der Laufzeit ab. Der Test selbst ist
in **allen** Läufen grün.

**Derselbe Gegenstand, anderer Träger.** Der bestehende Beleg
(`evidence/slice-057.md`) dateit einen `make test-integration`-Lauf, der real
mit Exit 2 in der kombinierten Retention-Lebenszyklus-Timing-Zusicherung
(`LH-FA-RET-004`) abbrach und bei Wiederholung grün durchlief. Hier ist der
Träger die **Coverage-Zahl** desselben Test-Objekts (`WAL-Retention`) — die
Beobachtung ist dieselbe: *ein zeitabhängiger Test liefert bei unverändertem
Stand wechselnde Ergebnisse.* Die Zuordnung ist das Register-Urteil (Modul 6:
„Mensch urteilt, Maschine prüft Deckung"); ein Beleg-Vorgang liegt mit der
Closure dieses Slice vor.

**Warum das zählt, obwohl kein Gate gefährdet ist:** Solange `THRESHOLD = 70 %`
und die gemessene Zahl rund 4,8 pp darüber liegt, entscheidet der Flap nichts.
Bei der Rampe nach [`ADR-0077`](../../../../../adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md)
(Endstufe 80 %) rückt die gemessene Zahl an die Schwelle, und dann wird aus
0,1 pp eine **flappende Gate-Entscheidung**: dasselbe Gate wäre grün und rot,
ohne dass sich eine Zeile geändert hätte. Die Klasse ist für Tests unsichtbar
(alle grün) und für ein diff-skopiertes Review nur über eine
**Wiederholungsmessung** sichtbar.

Quelle: Verifikationsbericht zu `slice-090` (V-1) ·
`internal/bootstrap/walretention_internal_test.go:175`, `:219` ·
`internal/bootstrap/wiring.go:981` · `harness/mk/coverage.mk` (`THRESHOLD ?= 70`).
