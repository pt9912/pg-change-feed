# Beleg: slice-093

Vorgang: `slice-093` — Coverage Cluster D2 (Bootstrap-Rest und Telemetrie).

Fund: **Der erste angenommene Fall dieses Eintrags** — und der erste, in dem die
Arbeit einen Satz **umstößt**, statt ihn nur sichtbar zu machen.

`harness/sensors/coverage-gate.md` §Zählbasis führte: „Der Träger der Schwankung
liegt in `internal/bootstrap/wiring.go` und ist in **zwei** Blöcken gemessen: dem
Takt-Zweig von `runWALRetentionCheck` (`:991.5,992.13`, 2 Statements …) und dem
Kontext-Ende-Zweig von `runAdministration` (`:1091.4,1092.1`, 1 Statement)."

**Am Parent war der Satz wahr** — gemessen vom Verifier: `wiring.go:991.5,992.13`
trug `count = 0`. Der neue netzlose Test
(`TestRunWALRetentionCheckProtokolliertMessfehlerUndLaeuftWeiter`,
`internal/bootstrap/wiring_rest_internal_test.go`) liefert dem Messer einen
**Fehler**, statt auf das Kontext-Ende zu warten — und fährt den Block damit
**deterministisch**: `count > 0` in **6 von 6** Läufen der Verifikation. Der Satz
ist durch die Arbeit **falsch geworden**, obwohl die Datei **nicht im Diff** lag.

**Die Grenze, die der Eintrag selbst zieht — und die diesen Fall von `slice-092`
scheidet.** Dort hatte ich dieselbe Stelle als Kandidaten vermutet (eine alternde
**Deixis** ohne Zahl-Drift); der Delta-Review hat ihn **abgelehnt**, weil §3.12
dort half: die Formulierung war schon bei Niederschrift falsch gebunden. Hier ist
das anders — der Satz war bei Niederschrift **wahr**, wurde von einer korrekten
Änderung **umgestoßen**, und §3.12 hilft nicht (der Ursprung war korrekt und der
Zeitpunkt auch). Genau darum steht dieser Eintrag bei 1× und nicht schon bei 2×:
**ein Kandidat wird geprüft, nicht gezählt.**

Gefunden hat es der **Implementer selbst** — auf den `grep`-Auftrag hin, der aus
diesem Eintrag stammt (Erstauftreten `slice-091`). Damit hat die Beobachtung sich
zum ersten Mal **bezahlt**: sie hat einen Fehler gefunden, den weder Review noch
Verifikation gesucht hätten, weil beide auf die Stellen ihres Auftrags sahen.

Quelle: `docs/reviews/verify-slice-093.md` (Richtung 1 und 2) ·
`harness/sensors/coverage-gate.md` §Zählbasis (berichtigt, Herkunfts-Anker
`slice-093`) · `internal/bootstrap/wiring_rest_internal_test.go` ·
`docs/plan/planning/observations/BEO-PGC/arbeit-ueberholt-stehenden-traeger/state.md`
(der abgelehnte Kandidat aus `slice-092`).
