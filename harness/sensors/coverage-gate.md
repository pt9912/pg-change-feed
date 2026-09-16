# `make coverage-gate` — Go-Test-Coverage gegen eine Schwelle

## Vertrag

Wird dieses Target rot, unterschreitet die Gesamt-Coverage über der
**netzlos prüfbaren Fläche** — `./internal/...`+`./cmd/...` **ohne** die
Pakete, deren Testlauf einen externen Dienst voraussetzt — die aktuell
gültige Schwelle (`THRESHOLD`). Vierte Docker-Multi-Stage-Stufe `coverage`
(nach `deps`, analog `d-check`s `Dockerfile`): `go test -coverpkg=<Pakete>
-coverprofile=… -covermode=atomic <Pakete>`, dann `go tool cover -func=…`,
dann `tools/coverage-gate.sh` gegen `THRESHOLD`
([`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md),
[`ADR-0054`](../../docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)).

Die Paketliste der Stufe kommt aus `go list`; der Filter nimmt die drei
Pakete aus, deren Testlauf einen externen Dienst voraussetzt
(`postgresstorage` ohne das Unterpaket `mapper`, `postgresack`,
`replication/receive`). Die tragende Regel ist die **Eigenschaft**, nicht die
Liste: Ein Paket, dessen Testlauf einen externen Dienst voraussetzt, ist
nicht Gegenstand dieses Gates. Die DB-gestützte Ebene dieser drei Pakete
trägt ihre eigene, subjekt-qualifizierte Messung (`ADR-0071` Punkt 3).

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

- **Statement-Zahlen** stammen aus dem Profil der `coverage`-Stufe über den
  Gegenstand (`/out/coverage.out`), **dedupliziert über die Block-Position**:
  jede Testbinary instrumentiert mit `-coverpkg` den ganzen Gegenstand, im
  gemergten Profil kommt dieselbe Block-Position darum mehrfach vor. „Gedeckt“
  heißt, dass **mindestens ein** Vorkommen `count > 0` trägt. Der **Nenner**
  ist die **Zustandsgröße** dieses Gegenstands: **1903 Statements** — sie hängt
  am Code-Stand, nicht am Lauf, ist darum **kein** Dauerwert und trägt den Lauf
  mit, in dem sie gemessen wurde (Lauf `slice-085`). Die **gedeckte** Zahl ist
  dagegen **kein** Zustand — sie trägt den Beleg **eines** Laufs und schwankt
  lauf-zu-lauf um ±2 Statements (ein `ctx`-abhängiger Pfad in
  `internal/adapters/driven/grpcstream/broadcaster.go` — `Publish` mit bereits
  beendetem `ctx`): **1369 von 1903** (**71,94 %**, gedruckt `71.90%`) und
  **1371** (**72,04 %**, gedruckt `72.00%`) sind die zwei beobachteten Enden
  desselben Stands (Lauf `slice-084`; im Lauf `slice-085` erneut beobachtet).
  Nenner und Abstand zur Schwelle sind von der Schwankung unberührt.
- Die von der Stufe **gedruckte** Prozentzeile (`total: (statements) XX.X%`,
  hier `71.9%`) ruht auf **derselben** Basis: auch dort zählt ein Block als
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
  instrumentiert sie darum nicht.

## Grenze — was das Grün nicht abdeckt

1. **Die DB-gestützte Fläche liegt außerhalb des Messgegenstands.** Die drei
   Pakete `internal/adapters/driven/postgresstorage` (ohne das Unterpaket
   `mapper`), `internal/adapters/driven/postgresack` und
   `internal/adapters/driving/replication/receive` setzen in ihren
   Testläufen einen externen Dienst voraus (PostgreSQL) und werden deshalb
   nicht in die Zahl dieses Gates gerechnet; ihre Netto-Abdeckung trägt die
   eigene, subjekt-qualifizierte Messung
   ([`db-adapter-coverage.md`](db-adapter-coverage.md)) aus
   [`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
   Punkt 3. Das Unterpaket `postgresstorage/mapper` bleibt im Gegenstand
   (20 Statements, 20 gedeckt, Lauf `slice-089`) — und ist damit **nicht** Teil der
   DB-Adapter-Coverage, deren Gegenstand `postgresstorage` ohne `mapper` führt;
   die zwei Zahlen überlappen nicht.

   **Fünf Pakete des Gegenstands haben keine Testdatei** (`go list
   -f '{{len .TestGoFiles}}'` über den Gegenstand). Sie tragen drei
   verschiedene Rollen:

   - **ohne ausführbare Statements** — im Profil kommen sie nicht vor:
     `postgresstorage/queries` (SQL-Textkonstanten), `application/port/inbound`
     (Schnittstellen-Deklarationen) und `domain/errors`
     (Sentinel-Deklarationen);
   - **über fremde Testpakete gedeckt** — `internal/adapters/driving/grpc/streamv1`
     (die generierten `changestream*.pb.go`) trägt **86 Statements, 61 gedeckt**
     (**70,9 %**, abgeleitet aus 61/86; Lauf `slice-089`); sie zählen, weil andere
     Testpakete mit `-coverpkg` über die
     Paketgrenze messen;
   - **vollständig ungedeckt** — `cmd/pg-change-feed` trägt **49 Statements,
     alle mit `count = 0`** (Lauf `slice-089`). Es ist damit das **einzige Paket des Gegenstands
     ohne ein einziges gedecktes Statement** und gehört zu der 80-%-Arbeit, die
     `welle-20` bündelt.

   Die beiden Pakete mit Statements weist der Lauf nicht als `[no test files]`
   aus, sondern als `coverage: 0.0% of statements`; das ist die Aufrufform mit
   `-coverpkg`.
2. **Docker-Layer-Caching.** `--no-cache-filter coverage` erzwingt die
   Neu-Auswertung der Stage bei jedem `make coverage-gate`-Lauf — ohne
   diesen Flag könnte ein Cache-Hit einen veralteten Lauf überleben lassen.
3. **Keine Zeilen-Ausnahme.** Das Gate hat strukturell keinen
   Suppression-Pfad (`AGENTS.md` §3.2) — die Gesamt-Coverage besteht oder
   scheitert als Zahl.
4. **Die Rücknahme eines ausgenommenen Pakets ist nur unvollständig
   gewächtert.** **Rückrechnung** aus den gemessenen Paket-Zahlen
   (§Zählbasis: 1369 gedeckt von 1903 Statements im Gegenstand; 30/32, 31/472,
   112/187 in den drei Ausgenommenen) — kein eigener Lauf: wird `postgresack`
   wieder in `-coverpkg` genommen, bleibt die Stufe grün —
   `(1369 + 30) / (1903 + 32) = 72,30 %` ≥ 70; dasselbe gilt für die Rücknahme
   von `replication/receive`
   (`(1369 + 112) / (1903 + 187) = 1481 / 2090 = 70,86 %` ≥ 70) — dessen
   netzlose Deckung von **112 von 187 (59,89 %)** macht die Rücknahme für die
   Prozent-Schwelle unsichtbar. Rot färbt sie nur die Rücknahme von
   `postgresstorage` (`(1369 + 31) / (1903 + 472) = 1400 / 2375 = 58,95 %`).
   Alle drei Ausgänge liegen weiter als die Lauf-zu-Lauf-Schwankung
   (±2 Statements = ±0,10 pp auf dem Gegenstands-Nenner) von der Schwelle
   entfernt: `postgresstorage` mit −11,05 Prozentpunkten darunter, die beiden
   grünen mit **+2,30** (`postgresack`) und **+0,86** Prozentpunkten
   (`replication/receive`) darüber — zum Kippen wären dort ≈45 bzw. ≈18
   Statements nötig. **Der
   Wächter ist** damit allein die Prozent-Schwelle, und sie trägt die
   Gegenstands-Hälfte der Fitness Function aus
   [`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
   („keine Block-Position im Profil liegt in …“) nicht vollständig: einen
   eigenen Sensor hat kein Rücknahme-Fall, und die Schwelle fängt nur einen der
   drei — **zwei** bleiben grün. Der
   Re-Evaluierungs-Trigger (a) derselben ADR greift beim Kommen oder Gehen
   eines Pakets, nicht bei der Rücknahme eines bereits ausgenommenen.
5. **Die Testpaket-Liste ist Disziplin, kein Sensor.** Ob die drei
   ausgenommenen Pakete in der Testpaket-Liste stehen oder nicht, ändert die
   Zahl nicht — ihre Testdateien überspringen netzlos ohnehin. **Der Wächter
   ist**: keiner; die Einhaltung ist eine Aussage des Rezepts, keine Messung.

## Ausgabe und Ausgänge

| Exit | Bedeutung |
|---|---|
| 0 | Gesamt-Coverage ≥ `THRESHOLD` |
| 1 | Gesamt-Coverage < `THRESHOLD` (`coverage-gate: FAIL`) |
| 2 | Coverage-Eingabe fehlt/leer, `total:`-Zeile fehlt, oder Prozentwert nicht parsbar |

Rot-/Grün-Beleg (real, [`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md);
die beiden Prozente sind die **gedruckten** Zeilen **eines** Laufs — Lauf `slice-081`,
beide Kommandos real gefahren in
[`verify-slice-081.md`](../../docs/reviews/verify-slice-081.md), Zählbasis §Zählbasis):
`THRESHOLD=75` (über dem Ist-Stand) lässt die Stage real scheitern
(`coverage-gate: FAIL — Coverage 71.30% unter Schwelle 75%`) — das Gate-Skript
endet Exit 1, `make` meldet für den gescheiterten Bauprozess Exit 2;
`THRESHOLD=70` besteht real (`coverage-gate: OK — Coverage 71.30% erfüllt
Schwelle 70%`, Exit 0).

## Bindung

[`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
· [`ADR-0054`](../../docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
· `tools/coverage-gate.sh` · `harness/mk/coverage.mk` · seit slice-049.
