Zustand: **verkörpert** — Ausgang: **verkörpert** → `AGENTS.md` §3.12 Instanz A
(*„die gedeckte Zahl hängt am Lauf … sie nennt ihn und nie ‚der Ist-Stand‘“*) und
`harness/sensors/coverage-gate.md` §Zählbasis + §Grenze 4 (das gemessene **Band**
dieses Gegenstands, sein verbleibender Block, die abgeleitete Schranke 1 Statement)
· seit slice-089 / seit slice-093. **Nicht gestrichen** — die Beobachtung kann noch
auttreten. **Verworfen** wurde der Kandidat „die Slice-Pläne als Risiko“: ein
Zeitdokument trägt keine Regel (`ADR-0083` §Verglichene Alternativen).
**Benannte Resthälfte:** der Erstauftreten-Fall `slice-057` hat keine Regel und
keinen Fix — er gehört eigenständig geführt, nicht in diesen Eintrag.

Die drei Belege treffen **denselben Gegenstand über verschiedene Träger**:
`slice-057` den `make test-integration`-Exit, `slice-090` die **Coverage-Zahl**
desselben Test-Objekts (`WAL-Retention`, `runWALRetentionCheck`) — dort brach ein
Lauf real ab, hier liefert `go tool cover -func` in 7 von 8 Läufen 87,5 % und in
einem 100,0 %, bei durchweg grünen Tests. Beide Male wechselt das Ergebnis bei
**unverändertem Stand**.
