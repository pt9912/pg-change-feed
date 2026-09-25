# `make coverage-gate` — Go-Test-Coverage gegen eine Schwelle

## Vertrag

Wird dieses Target rot, unterschreitet die Gesamt-Coverage über der
**netzlos prüfbaren Fläche** — `./internal/...`+`./cmd/...`+`./gen/...`
(`./gen/...` seit `slice-097`: der Umzug der erzeugten Vertragsfläche,
`ADR-0076`, bewegt den Träger, nicht den Gegenstand — der Architect-Verdikt
zum Coverage-Messgegenstand von `slice-097`) **ohne**
die Pakete, deren Testlauf einen externen Dienst voraussetzt — die aktuell
gültige Schwelle (`THRESHOLD`). Vierte Docker-Multi-Stage-Stufe `coverage`
(nach `deps`, analog `d-check`s `Dockerfile`): `go test -coverpkg=<Pakete>
-coverprofile=… -covermode=atomic <Pakete>`, dann `go tool cover -func=…`,
dann `tools/coverage-gate.sh` gegen `THRESHOLD`
([`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md),
[`ADR-0054`](../../docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)).

Die Paketliste der Stufe kommt aus `go list`; der Filter nimmt die vier
Pakete aus, deren Testlauf einen externen Dienst voraussetzt
(`postgresstorage`, `postgresack`, `postgressnapshot`,
`replication/receive`). Das Muster ist END-verankert und trifft genau diese
Pakete, nicht ihre Unterpakete: `postgresstorage/mapper` und
`postgressnapshot/snapshotlogic` bleiben im Gegenstand. Die tragende Regel ist
die **Eigenschaft**, nicht die Liste: Ein Paket, dessen Testlauf einen externen
Dienst voraussetzt, ist nicht Gegenstand dieses Gates. Die DB-gestützte Ebene
dieser vier Pakete trägt ihre eigene, subjekt-qualifizierte Messung
(`ADR-0071` Punkt 3).

Vor dem Bau hält
[`tools/harness/db-package-lists-check.sh`](../../tools/harness/db-package-lists-check.sh)
die namentlichen Stellen der Liste gleich (Filter, `DB_COVERAGE_PKGS`, die
Paketlisten der beiden Messläufe; Vertrag:
[`db-adapter-coverage.md`](db-adapter-coverage.md) §Gegenstand); ein
Unterschied färbt `make coverage-gate` rot, bevor der Bau läuft (§Grenze
Nr. 8).

`test/integration/` bleibt außerhalb: eigenständige `integration_test`-
Paketwurzel unter `test/` mit Black-Box-Tests gegen einen laufenden
Compose-Container ([`ADR-0030`](../../docs/plan/adr/0030-testpyramide.md)
E2E-Tier) — kein Unit-Coverage-Kandidat.

## Kalibrierungs-Bindung (bootstrap-aware Gate)

| Stufe | Wert | Ereignis |
|---|---|---|
| Einstieg | **70 %** | real gemessener Ist-Stand auf der netzlos prüfbaren Fläche — die von der Stufe **gedruckte** Prozentzeile des **Kalibrierungs-Laufs** (71,3 %, Lauf `slice-081`) —, angehoben auf die nächste volle 5-%-Stufe ([`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md), Mechanik `ADR-0054` §(a), Anhebung durch den Subjekt-Transfer aus [`ADR-0077`](../../docs/plan/adr/0077-coverage-rampen-neu-bemessung-subjekt-transfer.md); die Größe dieser Zeile: §Zählbasis) |
| Endstufe | **80 %** | fest, Nutzer-Entscheidung ([`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md); `roadmap.md` §Nächste Wellen) |

**Geltende Stufe:** Der bewegliche Wert dieser Rampe steht ausschließlich in
[`harness/mk/coverage.mk`](../mk/coverage.mk) (`THRESHOLD`) — diese Sektion
beschreibt die Bindung (Rampe, Endstufe, Trigger) und führt ihn als **Träger**
nicht; sie zitiert ihn nur in den **Belegen** weiter unten, die einen konkreten
Lauf bezeugen und nicht wandern (`AGENTS.md` §3.7) · seit slice-076.

**Hochschalt-Trigger:** die nächste Coverage-Verbesserung schließt die
Lücke zur nächsten 5-%-Stufe (`THRESHOLD` in
[`harness/mk/coverage.mk`](../mk/coverage.mk) anheben), bis 80 % erreicht
ist — dieselbe Reifung, die d-check selbst durchlief (85 → 90 → 93,
`d-check`s `Makefile`). Das ist **keine** Schwellen-Senkung
(`AGENTS.md` §3.6 bleibt unverletzt): Die Endstufe steht fest, nur der
Einstiegspunkt hängt am real gemessenen Ist-Stand.

## Zählbasis der Zahlen dieser Datei

**Der Nenner der Stufe ist 2567** (Lauf `slice-backfill-slot-leerlauf-bestaetigung`,
Implementer-Lauf am Arbeitsbaum über `f4e32fba`: `make coverage-gate` baut die
Stufe `coverage`, das Profil `/out/coverage.out` des gebauten Images,
dedupliziert über die Block-Position mit Awk ausgezählt; gedeckt **2136 von
2567** = 83,21 %, gedruckt `total: (statements) 83.2%` und `coverage-gate: OK —
Coverage 83.20% erfüllt Schwelle 80%`). Der Nenner trägt gegenüber dem Nenner
2558 am Stand `8d8860f7` (unten) neun Statements mehr (**abgeleitet**); die Produktionsdateien
des Slice im Gegenstand sind `usecase/capture/service.go` (47 von 47 gedeckt),
`replication/mapper/mapper.go` (168 von 175) und `bootstrap/wiring.go` (237 von
604), aus demselben Profil (`git diff --stat f4e32fba` über `internal`, `cmd`,
`gen` ohne Testdateien zeigt außer diesen drei nur `replication/receive/receive.go`,
das im Gegenstand der DB-Adapter-Coverage liegt, und den Port
`port/inbound/idleconfirmation.go` ohne ausführbare Statements). Die
Lokatoren der zwei Blöcke von `wiring.go`, die der Abschnitt unten nennt, liegen
am Stand dieses Laufs bei `:1287.4,1288.1` (Kontext-Ende-Zweig von
`runAdministration`, 1 Statement; die Funktion trägt einen weiteren
`return`-Zweig bei `:1293.4,1294.1`) und `:1183.5,1184.13` (Fehlerzweig der
WAL-Rückstands-Messung in `runWALRetentionCheck`, 2 Statements), gemessen im
Profil dieses Laufs.

**Der Nenner der Stufe am Stand `8d8860f7` war 2558** (Lauf `slice-backfill-bench-richtgroesse`,
Closure am Stand `8d8860f7`: `make coverage-gate` baut die Stufe `coverage`, das
Profil `/out/coverage.out` des gebauten Images, dedupliziert über die
Block-Position mit Awk ausgezählt; gedeckt **2129 von 2558** = 83,23 %, gedruckt
`total: (statements) 83.2%` und `coverage-gate: OK — Coverage 83.20% erfüllt
Schwelle 80%`). Der Nenner setzt sich aus **2504** Statements der übrigen Pakete
(**abgeleitet**) und den **54** Statements des Unterpakets
`postgressnapshot/snapshotlogic` zusammen (54 von 54 gedeckt, aus demselben
Profil); das Unterpaket liegt im Gegenstand. Die Nenner-Größe ist an den
Code-Stand gebunden, die gedeckte Zahl an den Lauf: `make gates` im
Verifikations-Lauf am Stand `c97273b3` druckte `Coverage 83.20%`
(**übernommen** aus `verifikation-slice-backfill-bench-richtgroesse` §1, dieselbe
gedruckte Zeile wie der Lauf oben). Die Produktionsdateien des Slice liegen im
Gegenstand: `usecase/backfill/warn.go` trägt **9** Statements (9 gedeckt, aus
demselben Profil), `usecase/backfill/service.go` **214** (200 gedeckt, aus
demselben Profil); der Nenner-Unterschied von **17** Statements zum Nenner **2541**
des Laufs `slice-backfill-e2e` am Stand `cf7f2d02` (dort gedeckt 2112, gedruckt
`Coverage 83.10%`) ist **abgeleitet** und liegt in diesen zwei Dateien
(`git diff --stat cf7f2d02..HEAD` über `internal`, `cmd`, `gen` ohne
Testdateien zeigt außer einem Kommentar in `internal/domain/model/backfillrun.go`
nur sie).

**Der Nenner des Stands von `slice-097` war 1936, nicht 1903** (Lauf
`slice-097`, der
Architect-Verdikt zum Coverage-Messgegenstand von `slice-097` §2/§5) — die
Differenz trägt zwei Ursachen: `slice-096` bewegte Produktionscode in
`internal/bootstrap/{config_file,wiring}.go` (+33, abgeleitet), und `slice-097`
zog die erzeugte Vertragsfläche (`gen/cdc/stream/v1`, vormals
`internal/adapters/driving/grpc/streamv1`) an einen öffentlichen Pfad um und
nahm `./gen/...` wieder in die `go list`-Zeile auf — derselbe Gegenstand, kein
neuer (86 Statements, unverändert). Quote am `slice-097`-Stand: **83,3–83,4 %**
(Band 1 Statement = 0,05 pp). **Die 1903-Werte unten bleiben als datierte
Messung des `slice-085`-Laufs gültig** — sie sind kein Dauerwert und waren es
nie; wer den Ist-Stand braucht, liest diesen Absatz, nicht die Zahl darunter.

- **Statement-Zahlen** stammen aus dem Profil der `coverage`-Stufe über den
  Gegenstand (`/out/coverage.out`), **dedupliziert über die Block-Position**:
  jede Testbinary instrumentiert mit `-coverpkg` den ganzen Gegenstand, im
  gemergten Profil kommt dieselbe Block-Position darum mehrfach vor. „Gedeckt“
  heißt, dass **mindestens ein** Vorkommen `count > 0` trägt. Der **Nenner**
  ist die **Zustandsgröße** dieses Gegenstands: **1903 Statements** — sie hängt
  am Code-Stand, nicht am Lauf, ist darum **kein** Dauerwert und trägt den Lauf
  mit, in dem sie gemessen wurde (Lauf `slice-085`). Die **gedeckte** Zahl ist
  dagegen **kein** Zustand — sie trägt den Beleg **eines** Laufs: **1369 von
  1903** (**71,94 %**, gedruckt `71.90%`) und **1371** (**72,04 %**, gedruckt
  `72.00%`) sind die zwei beobachteten Enden des Stands der Läufe `slice-084`
  und `slice-085`. Der Träger der Schwankung liegt in
  `internal/bootstrap/wiring.go`; von den zwei Blöcken, die die
  Acht-Lauf-Messung unten benennt, trägt sie nach dem Stand von `slice-093`
  nur noch **einer**: der Kontext-Ende-Zweig von `runAdministration`
  (`:1287.4,1288.1`, 1 Statement). Den zweiten — den Fehlerzweig der
  WAL-Rückstands-Messung in `runWALRetentionCheck` (`:1183.5,1184.13`,
  2 Statements), den
  [`ADR-0082`](../../docs/plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
  §Kontext (2) als Takt-Zweig benennt — fährt der netzlose Test
  `TestRunWALRetentionCheckProtokolliertMessfehlerUndLaeuftWeiter`
  (`internal/bootstrap/wiring_rest_internal_test.go`) deterministisch: er
  liefert dem Messer einen Fehler, statt auf das Kontext-Ende zu warten, und
  der Block trug `count > 0` in jedem der vier Läufe von `slice-093`. Die
  gedeckte Zahl kann damit nur noch um das eine Statement des
  `runAdministration`-Zweigs wandern (**abgeleitet**, kein eigener Lauf).
  Über **acht** Läufe des Test-Stands von `slice-091`
  (`go test -count=1 -coverpkg=… -covermode=atomic`,
  Auswertung über die Block-Position, Lauf `slice-091`) lag die gedeckte Zahl
  zwischen **1468** und **1471**, die gedruckte Zeile zwischen `77.1%` und
  `77.3%`; der Takt-Zweig trug in **einem** dieser Läufe `count > 0`, der
  `runAdministration`-Zweig fiel in **einem** auf `count = 0`. Das ist die
  Messung jenes Test-Stands, nicht die des geltenden; Nenner und
  Abstand zur Schwelle sind von der Schwankung unberührt.

  Das hier gemessene Band ist **nicht** das oben genannte: der **Nenner** ist an
  beiden Ständen **1903** und im Messgegenstand hat sich zwischen ihnen **kein**
  Produktcode bewegt (`git diff` über `internal`/`cmd` ohne Testdateien und ohne
  das ausgenommene `replication/receive` → leer) — die zwei Bänder liegen auf
  zwei **Test**-Ständen **desselben** Produktionsstands. Ihre Differenz
  (**abgeleitet**, nicht gemessen: 1468/1471 gegen 1369/1371) beträgt je nach
  Paarung 97 bis 102 Statements.
- Die von der Stufe **gedruckte** Prozentzeile (`total: (statements) XX.X%`,
  z. B. `71.9%`) ruht auf **derselben** Basis: auch dort zählt ein Block als
  gedeckt, wenn er ein Vorkommen mit `count > 0` trägt — die Summierung über die
  Testbinaries ändert dieses Prädikat nicht, und damit auch nicht das Verhältnis
  gedeckter zu instrumentierten **Statements**. Sie ist deshalb **keine eigene
  Größe**, sondern dieselbe Messung in anderer Ausgabepräzision: `go tool cover`
  druckt eine Nachkommastelle (`71.9%`), `tools/coverage-gate.sh` liest genau
  diese gedruckte Zeile und gibt sie mit `%.2f` aus (`71.90%`; Lauf `slice-084`,
  im Lauf `slice-085` erneut gedruckt). Die zwei Nachkommastellen der
  **Rechnung** (`1369 von 1903` = `71,94 %`, erster Punkt oben) sind die
  **deduplizierte** Auswertung, keine gedruckte Zeile. Eine absolute
  Statement-Zahl trägt nur die deduplizierte Auswertung; die gedruckte Zeile
  trägt keine.
- Die Zahlen der **drei ausgenommenen Pakete** (§Grenze Punkt 4) stammen aus
  einer eigenen Messung mit `-coverpkg` über **alle** Pakete (die drei
  eingeschlossen, netzlos) — **30/32**, **31/472** und **112/187**, Lauf
  `slice-085`; die Stufe dieses Gates nimmt die drei Pakete aus und
  instrumentiert sie darum nicht. Das vierte ausgenommene Paket
  (`postgressnapshot`, 130 Statements laut Nenner der DB-Adapter-Coverage,
  Lauf `slice-backfill-e2e`)
  trägt netzlos **0** gedeckte Statements: jeder seiner Tests überspringt ohne
  `CDC_REPLICATION_TEST_DSN` (Lauf `slice-backfill-snapshot-reader`, Fixrunde,
  `go test -v` ohne Netz und ohne Variable: 21 Tests, 21 `SKIP`).

## Grenze — was das Grün nicht abdeckt

1. **Die DB-gestützte Fläche liegt außerhalb des Messgegenstands.** Die vier
   Pakete `internal/adapters/driven/postgresstorage` (ohne das Unterpaket
   `mapper`), `internal/adapters/driven/postgresack`,
   `internal/adapters/driven/postgressnapshot` und
   `internal/adapters/driving/replication/receive` setzen in ihren
   Testläufen einen externen Dienst voraus (PostgreSQL) und werden deshalb
   nicht in die Zahl dieses Gates gerechnet; ihre Netto-Abdeckung trägt die
   eigene, subjekt-qualifizierte Messung
   ([`db-adapter-coverage.md`](db-adapter-coverage.md)) aus
   [`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
   Punkt 3. Das Unterpaket `postgresstorage/mapper` bleibt im Gegenstand
   (20 Statements, 20 gedeckt, Lauf `slice-089`) — und ist damit **nicht** Teil der
   DB-Adapter-Coverage, deren Gegenstand `postgresstorage` ohne `mapper` führt;
   die zwei Zahlen überlappen nicht. Ebenso bleibt das Unterpaket
   `postgressnapshot/snapshotlogic` im Gegenstand (54 Statements, 54 gedeckt;
   Lauf `slice-backfill-e2e`, Closure): es trägt die netzlos
   prüfbare Logik des Snapshot-Adapters, die deshalb nicht im ausgenommenen
   Paket liegt.

   **Drei Pakete des Gegenstands führen keine Testdatei** — kein eigenes und
   kein externes Testpaket. Die Gruppierung ist mechanisch, nicht gezählt nach
   einem Einzelmaß: `go list -f '{{.ImportPath}} Test={{len .TestGoFiles}}
   XTest={{len .XTestGoFiles}}'` über den Gegenstand liefert für **genau drei**
   Pakete `Test=0 XTest=0` (Lauf `slice-094`) — `TestGoFiles` allein trifft
   **22** der 31 Pakete und ist darum **keine** Gruppierungsregel: die 22
   zerfallen in diese drei und **19**, die ausschließlich ein externes Testpaket
   führen (Lauf `slice-094`). Die drei tragen **keine ausführbaren Statements** —
   im Profil kommen sie nicht vor, der Lauf weist sie als `[no test files]` aus:
   `postgresstorage/queries` (SQL-Textkonstanten), `application/port/inbound`
   (Schnittstellen-Deklarationen) und `domain/errors`
   (Sentinel-Deklarationen).

   **`cmd/pg-change-feed` steht nicht in dieser Gruppe.** Das Paket führt ein
   eigenes Testpaket (`TestGoFiles` = **1**, `XTestGoFiles` = **0**; Lauf
   `slice-094`) — den Re-Exec-Harness des Argument-Dispatchs: der Test startet
   dasselbe Binary mit anderen Argumenten und wertet Exit-Code, stdout und
   stderr des Kindprozesses, ohne Naht im Produktionscode
   ([`ADR-0082`](../../docs/plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)).
   Das Paket trägt **49 von 49** Statements gedeckt (Lauf `slice-094`): auch die
   vier Sondermodi werden gefahren, mit vollständiger Umgebung und nicht
   erreichbarer Instanz, und enden am Verbindungsaufbau ihrer eigenen Rolle
   (`ADR-0047`, `ReaderDSN` bzw. `AdminDSN`) — die Aufrufe `main.go:44`
   (`Healthcheck`), `:62` (`RegisterConsumer`), `:86` (`AcknowledgeConsumer`)
   und `:101` (`Diagnose`) sind netzlos erreichbar. Dienstgebunden ist der
   **Rumpf** dieser vier Funktionen, nicht ihr Aufruf.

   **`gen/cdc/stream/v1` ist keins dieser drei Pakete.**
   Das Paket (die generierten `changestream*.pb.go`) führt ein **eigenes**,
   externes Testpaket (`changestream_test.go`, `package streamv1_test`;
   `XTestGoFiles` = **1**, `TestGoFiles` = **0**) und hat damit einen eigenen
   Testlauf — zugleich läuft seine **Deckung** weiter über die Paketgrenze: der
   eigene Lauf trägt **45 von 86** Statements (**52,3 %**, Lauf `slice-091`), im
   Gegenstand mit `-coverpkg` sind es **76 von 86** (**88,4 %**, abgeleitet aus
   76/86; Lauf `slice-091`). Es ist weder „ohne Testdatei" noch allein „über
   fremde Testpakete gedeckt".

   Die drei Pakete ohne ausführbare Statements weist der Lauf als
   `[no test files]` aus. In der Aufrufform der Stufe (`-coverpkg` über den
   ganzen Gegenstand) trägt **kein** Paket die Zeile
   `coverage: 0.0% of statements`, und kein
   Paket **mit** ausführbaren Statements trägt null gedeckte Statements (Lauf
   `slice-094`).
2. **Docker-Layer-Caching.** `--no-cache-filter coverage` erzwingt die
   Neu-Auswertung der Stage bei jedem `make coverage-gate`-Lauf — ohne
   diesen Flag könnte ein Cache-Hit einen veralteten Lauf überleben lassen.
3. **Keine Zeilen-Ausnahme.** Das Gate hat strukturell keinen
   Suppression-Pfad (`AGENTS.md` §3.2) — die Gesamt-Coverage besteht oder
   scheitert als Zahl.
4. **Die Rücknahme eines ausgenommenen Pakets ist durch die Prozent-Schwelle
   nur unvollständig gewächtert.** **Rückrechnung, datiert** — *kein eigener Lauf*, und die drei
   Paket-Zahlen stammen aus dem Stand von `slice-085`: wird `postgresack`
   wieder in `-coverpkg` genommen, blieb die Stufe **damals** (`THRESHOLD=70`)
   grün — `(1369 + 30) / (1903 + 32) = 72,30 %` ≥ 70; dasselbe galt für die
   Rücknahme von `replication/receive`
   (`(1369 + 112) / (1903 + 187) = 1481 / 2090 = 70,86 %` ≥ 70) — dessen
   netzlose Deckung von **112 von 187 (59,89 %)** macht die Rücknahme für die
   Prozent-Schwelle unsichtbar. Rot färbte sie nur die Rücknahme von
   `postgresstorage` (`(1369 + 31) / (1903 + 472) = 1400 / 2375 = 58,95 %`).
   **Auf dem heutigen Stand und der heutigen Stufe** (`Endstufe 80`,
   `slice-094`-Lauf: 1581 von 1903; die drei Paket-Zahlen weiter aus
   `slice-085` — **gemischter Stand, deshalb abgeleitet, nicht gemessen**)
   ergibt dieselbe Rechnung **83,26 %** (`postgresack`), **81,00 %**
   (`replication/receive`) und **67,87 %** (`postgresstorage`): die
   **Schlussfolgerung hält** — nur die Rücknahme von `postgresstorage` färbt
   die Stufe rot —, und die zwei grünen Fälle liegen jetzt **weiter** darüber
   als damals. **Die datierte Fassung und die heutige sind beide gültig; wer
   nur eine liest, liest die andere Stufe mit.**
   Alle drei Ausgänge liegen weiter als die Lauf-zu-Lauf-Schwankung

   (dem in §Zählbasis geführten Band — nach `slice-093` höchstens **1**
   Statement = 0,05 pp auf dem Gegenstands-Nenner, abgeleitet aus 1/1903,
   weil der zweite Block dort deterministisch gedeckt ist —) von der Schwelle
   entfernt: `postgresstorage` mit −11,05 Prozentpunkten darunter, die beiden
   grünen mit **+2,30** (`postgresack`) und **+0,86** Prozentpunkten
   (`replication/receive`) darüber — zum Kippen wären dort ≈45 bzw. ≈18
   Statements nötig. **Innerhalb der Stufe ist**
   die Prozent-Schwelle damit der einzige Wächter dieser Rücknahmen, und sie trägt die
   Gegenstands-Hälfte der Fitness Function aus
   [`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
   („keine Block-Position im Profil liegt in …“) nicht vollständig: einen
   eigenen Sensor hat innerhalb der Stufe kein Rücknahme-Fall, und die Schwelle
   fängt nur einen der drei — **zwei** bleiben dort grün. Der
   Re-Evaluierungs-Trigger (a) derselben ADR greift beim Kommen oder Gehen
   eines Pakets, nicht bei der Rücknahme eines bereits ausgenommenen.

   **Das vierte ausgenommene Paket, `postgressnapshot`, ist gemessen
   gewächtert.** Seine Tests überspringen netzlos alle, es trägt netzlos **0**
   gedeckte Statements; nimmt der Filter der Stufe es nicht mehr aus, färbt die
   Prozent-Schwelle rot: der Lauf `slice-backfill-snapshot-reader` (Fixrunde)
   mit dem Dockerfile-Filter ohne `postgressnapshot` druckt
   `coverage-gate: FAIL — Coverage 78.60% unter Schwelle 80%` (Exit ≠ 0 des
   `docker build`). Mit dem Eintrag druckt derselbe Stand
   `coverage-gate: OK — Coverage 83.20% erfüllt Schwelle 80%`. Von den vier
   Rücknahme-Fällen färben damit **zwei** die Stufe über die Prozent-Schwelle rot
   (`postgresstorage`, `postgressnapshot`) und **zwei** bleiben dort grün.
   Über die Prozent-Schwelle hinaus färbt `db-package-lists-check.sh`
   (`make coverage-gate` ruft es vor dem Bau) **jede** der vier Rücknahmen rot:
   der Filter weicht dann von `DB_COVERAGE_PKGS` ab (Nr. 8).
5. **Die Testpaket-Liste ist Disziplin, kein Sensor.** Ob die vier
   ausgenommenen Pakete in der Testpaket-Liste stehen oder nicht, ändert die
   Zahl nicht — ihre Testdateien überspringen netzlos ohnehin. **Der Wächter
   ist**: keiner; die Einhaltung ist eine Aussage des Rezepts, keine Messung.
6. **Was der Build-Kontext der Stufe nicht enthält, kann kein Test der Stufe
   lesen.** Ein Test, der ein **reales Artefakt** außerhalb von `cmd/`,
   `internal/`, `go.mod` und `go.sum` prüfen soll — die Form, die
   [`ADR-0082`](../../docs/plan/adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md)
   als Kern dieser Welle verlangt —, braucht eine benannte Ausnahme in
   `.dockerignore`. Deren **Klasse** (gelesen, nicht gebaut · genau eine Datei ·
   der Leser steht dabei · Image unberührt mit Beleg) führt
   [`ADR-0085`](../../docs/plan/adr/0085-build-kontext-ausnahme-test-only-zweck.md);
   **der Wächter ist**: keiner — die Klasse ist ein Zweck-Urteil, und ihr
   Entdecker ist der rote Bau (Lauf `slice-093`).
7. **Die Deckung von `cmd/pg-change-feed` hängt an einer Weitergabe, die keine
   Testzusage trägt.** Der Re-Exec-Harness des Argument-Dispatchs zählt im
   Profil nur, weil er `GOCOVERDIR` an die Kindprozesse durchreicht: `go test`
   setzt die Variable, die Kindprozesse schreiben ihre Zähler-Dateien dorthin,
   und `go test` mergt sie in dasselbe Profil. Fällt die Weitergabe weg, bleiben
   **alle** Tests dieses Pakets grün und seine Deckung fällt auf **0 von 49**
   zurück; die Gesamt-Coverage fällt mit ihr auf **1532 von 1903 = 80,50 %**,
   gedruckt `80.5%` (Lauf `slice-094`, netzlos, dedupliziert über die
   Block-Position). **Der Wächter ist**: keiner — die Bindung ist eine Zeile in
   `kindUmgebung` (`cmd/pg-change-feed/main_test.go`), und ihr Entdecker ist der
   Vergleich zweier Profile.

8. **Die Gleichheit der namentlichen Listen ist gewächtert, die Eigenschaft
   nicht.** `db-package-lists-check.sh` prüft, dass Dockerfile-Filter,
   `DB_COVERAGE_PKGS` und die Paketlisten der beiden Messläufe dieselben Pakete
   nennen: das Streichen von `postgressnapshot` aus je einer der vier Stellen, das
   Ersetzen von `postgresstorage` in der Store-Liste und ein unbekannter
   Filter-Eintrag enden je mit Exit 1 (Lauf `slice-backfill-snapshot-reader`,
   Fixrunde). Ob ein Paket die Eigenschaft „Testlauf setzt einen externen Dienst
   voraus" trägt, und ob jeder seiner Tests ohne Datenbank überspringt, prüft kein
   Sensor — Disziplin, benannt in
   [`db-adapter-coverage.md`](db-adapter-coverage.md) §Grenze Nr. 8.

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | Gesamt-Coverage ≥ `THRESHOLD` |
| 1 | Gesamt-Coverage < `THRESHOLD` (`coverage-gate: FAIL`) oder die Paketlisten weichen ab (`db-package-lists-check:`-Meldung, Bau nicht gestartet) |
| 2 | Coverage-Eingabe fehlt/leer, `total:`-Zeile fehlt, oder Prozentwert nicht parsbar |

Rot-/Grün-Beleg (real, [`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md);
die beiden Prozente sind die **gedruckten** Zeilen **eines** Laufs — Lauf `slice-081`,
beide Kommandos real gefahren im Verifikationsbericht zu `slice-081`,
Zählbasis §Zählbasis):
`THRESHOLD=75` (über dem Ist-Stand) lässt die Stage real scheitern
(`coverage-gate: FAIL — Coverage 71.30% unter Schwelle 75%`) — das Gate-Skript
endet Exit 1, `make` meldet für den gescheiterten Bauprozess Exit 2;
`THRESHOLD=70` besteht real (`coverage-gate: OK — Coverage 71.30% erfüllt
Schwelle 70%`, Exit 0).

## Bindung

[`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
· [`ADR-0054`](../../docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
· `tools/coverage-gate.sh` · `tools/harness/db-package-lists-check.sh` ·
`harness/mk/coverage.mk` · seit slice-049.
