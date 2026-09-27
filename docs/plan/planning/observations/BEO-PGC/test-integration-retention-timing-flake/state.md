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

**Vierter Beleg (`slice-capture-leerlauf-quellbelege`) — die Resthälfte hat ein
zweites Auftreten.** Der erste `e2e.yml`-Lauf nach dem Push (Lauf 36287009221, Leg
PostgreSQL 18, erster Versuch) endete rot in der Phase „Leerlauf-Bestätigung“ von
`make test-integration`: der Feed-Container endete mit dem WAL-Schwellen-Fehler
(Rückstand 15.238.216 B über der Fehlerschwelle 8.388.608 B), im Wiederholungsversuch
lief dieselbe Phase grün. Gemessen: 1 rotes erstes Ergebnis unter den 40 Läufen von
`e2e.yml`, deren Commit die Phase enthält (`evidence/slice-capture-leerlauf-quellbelege.md`).
Zähler: **4×**. Die Zählung teilt sich nach dem, was sie trägt: die **Zahl-Hälfte**
(Coverage, `slice-090`/`slice-091`) bleibt **verkörpert**; die **Resthälfte** — ein
`make test-integration`-Lauf endet rot bei unverändertem Stand und die Wiederholung
ist grün — steht bei **2×** (`slice-057`, `slice-capture-leerlauf-quellbelege`),
**ohne Ausgang**, unter der Schwelle von 3×: ein Slice ist nicht fällig, und der
Eintrag nimmt kein Ausgang vorweg, den nur die Ursache tragen kann.

**Offene Frage an den Architect** (Adresse: dieser Eintrag und
`evidence/slice-capture-leerlauf-quellbelege.md`; Zug: ein eigenes Verdikt, das die
Ursache entscheidet und, wenn nötig, einen Slice beauftragt). *Endete der Container
in der Phase „Leerlauf-Bestätigung“ aus einem Grund, den das Produkt trägt, oder aus
einem Grund des Testaufbaus?* Die gemessene Signatur: etwa 0,1 s nachdem ein einzelnes
`INSERT` über 60.000 Zeilen auf eine nicht aktivierte Tabelle zurückkehrte (Dauer etwa
20 s, abgeleitet aus den Log-Stempeln), meldete die Schwellen-Prüfung
(`runWALRetentionCheck`) einen Rückstand über der Fehlerschwelle, obwohl die
Leerlauf-Bestätigung (`ADR-0120`) den Rückstand einer Last ohne Inhalt für die
Publication absenken soll (harness/README.md, `make test-replication`: Form X1). Zwei
Lesarten, **keine geprüft**: (A) *Produkt* — die Schwellen-Prüfung kann bei einer
solchen Last auf einem langsamen Runner die Fehlerschwelle erreichen, bevor die
Bestätigung den Rückstand senkt; dann ist die Reichweite von `ADR-0120` gegen die
Fehlerschwelle die Frage (eine Norm-, keine Test-Frage). (B) *Testaufbau* — die Last
(rund das Doppelte der Fehlerschwelle in etwa 20 s) lässt dem Takt der Prüfung und dem
Keepalive-Takt der Quelle auf einem langsamen Runner zu wenig Reserve; dann braucht die
Phase eine andere Last-/Schwellen-Auslegung oder eine Wartebedingung. Die Entscheidung
liegt beim Architect; dieser Eintrag und der Plan des Quellbeleg-Slice entscheiden sie
nicht.

**Reihenfolge, falls der Architect einen Stabilisierungs-Slice beauftragt.** Er gehört
**vor** `slice-wal-fehlerschwelle-ausgangsklasse` und `slice-start-vorlauf-grenze`:
beide erweitern denselben Runner (`tools/harness/run-integration-tests.sh`) **hinter**
der Phase „Leerlauf-Bestätigung“, und ein Rot dort lässt jede Phase dahinter und die
Schritte des Legs dahinter ungelaufen (im Lauf 36287009221, Versuch 1, Leg
PostgreSQL 18: Coverage- und „Replication-Tier“-Schritte `skipped`) — ihre Belege
(die Mutation an der Phase „Fehlerschwelle beendet den Container“, die Phasen des
Startpfads) wären dann nicht widerlegt, sondern nicht gelaufen. Die Kanten:
`slice-capture-leerlauf-quellbelege` → *(Stabilisierung, falls beauftragt)* →
`slice-wal-fehlerschwelle-ausgangsklasse` → `slice-start-vorlauf-grenze` →
`slice-transformationen-e2e-abhilfe`. Der Slice ist nicht angelegt: die
Register-Regel verlangt ihn erst ab 3× ohne Ausgang, und ob er nötig ist, entscheidet
die Antwort auf die Frage oben.
