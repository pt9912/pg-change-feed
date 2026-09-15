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
| Einstieg | **65 %** | real gemessener Ist-Stand auf der netzlos prüfbaren Fläche — die von der Stufe **gedruckte** Prozentzeile des Kalibrierungs-Laufs (69,70 %) —, abgerundet auf die nächste volle 5-%-Stufe ([`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md), Mechanik `ADR-0054` §(a); die Größe dieser Zeile: §Zählbasis) |
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
  heißt, dass **mindestens ein** Vorkommen `count > 0` trägt (dedupliziert:
  1679 Statements, davon 1171 gedeckt = 69,74 %).
- Die von der Stufe **gedruckte** Prozentzeile (`total: (statements) XX.X%`,
  hier `69.7%`) ruht auf **derselben** Basis: auch dort zählt ein Block als
  gedeckt, wenn er ein Vorkommen mit `count > 0` trägt — die Summierung über die
  Testbinaries ändert dieses Prädikat nicht, und damit auch nicht das Verhältnis
  gedeckter zu instrumentierten **Statements**. Sie ist deshalb **keine eigene
  Größe**, sondern dieselbe Messung in anderer Ausgabepräzision: `go tool cover`
  druckt eine Nachkommastelle, `tools/coverage-gate.sh` formatiert `%.2f` —
  daraus werden `69.7%` und `69.70%`. Eine absolute Statement-Zahl trägt nur die
  deduplizierte Auswertung; die gedruckte Zeile trägt keine.
- Die Zahlen der **drei ausgenommenen Pakete** (§Grenze Punkt 4) stammen aus
  dem Profil des Gegenstands **vor** dem Schnitt — derselben Messung, die den
  Nenner `2467 → 1679` beziffert (vorher 1215 gedeckt, davon 44 in den drei
  Paketen).

## Grenze — was das Grün nicht abdeckt

1. **Die DB-gestützte Fläche liegt außerhalb des Messgegenstands.** Die drei
   Pakete `internal/adapters/driven/postgresstorage` (ohne das Unterpaket
   `mapper`), `internal/adapters/driven/postgresack` und
   `internal/adapters/driving/replication/receive` setzen in ihren
   Testläufen einen externen Dienst voraus (PostgreSQL) und werden deshalb
   nicht in die Zahl dieses Gates gerechnet; ihre Netto-Abdeckung trägt die
   eigene, subjekt-qualifizierte Messung aus
   [`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
   Punkt 3. Das Unterpaket `postgresstorage/mapper` bleibt im Gegenstand
   (15 Statements, 12 gedeckt).

   **Fünf Pakete des Gegenstands haben keine Testdatei** (`go list
   -f '{{len .TestGoFiles}}'` über den Gegenstand). Sie tragen drei
   verschiedene Rollen:

   - **ohne ausführbare Statements** — im Profil kommen sie nicht vor:
     `postgresstorage/queries` (SQL-Textkonstanten), `application/port/inbound`
     (Schnittstellen-Deklarationen) und `domain/errors`
     (Sentinel-Deklarationen);
   - **über fremde Testpakete gedeckt** — `internal/adapters/driving/grpc/streamv1`
     (die generierten `changestream*.pb.go`) trägt **86 Statements, 61 gedeckt
     (70,9 %)**; sie zählen, weil andere Testpakete mit `-coverpkg` über die
     Paketgrenze messen;
   - **vollständig ungedeckt** — `cmd/pg-change-feed` trägt **49 Statements,
     alle mit `count = 0`**. Es ist damit das **einzige Paket des Gegenstands
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
   (§Zählbasis: 1171 gedeckt von 1679 Statements im Gegenstand; 2/23, 31/610,
   11/155 in den drei Ausgenommenen) — kein eigener Lauf: wird `postgresack`
   wieder in `-coverpkg` genommen, bleibt die Stufe grün —
   `(1171 + 2) / (1679 + 23) = 68,92 %` ≥ 65; erst die Rücknahme von
   `postgresstorage` (`1202 / 2289 = 52,51 %`) oder `replication/receive`
   (`1182 / 1834 = 64,45 %`) färbt sie rot. Die drei Ausgänge liegen mit
   +3,9 / −12,5 / −0,6 Prozentpunkten weit genug von der Schwelle, dass die
   Lauf-zu-Lauf-Schwankung (wenige Statements) sie nicht umkehrt. **Der
   Wächter ist** damit allein die Prozent-Schwelle, und sie trägt die
   Gegenstands-Hälfte der Fitness Function aus
   [`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
   („keine Block-Position im Profil liegt in …“) nicht vollständig: für einen
   einzelnen Rücknahme-Fall gibt es keinen eigenen Sensor. Der
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
die beiden Prozente sind die **gedruckten** Zeilen der Läufe, §Zählbasis):
`THRESHOLD=75` (über dem Ist-Stand) lässt die Stage real scheitern
(`coverage-gate: FAIL — Coverage 69.70% unter Schwelle 75%`) — das Gate-Skript
endet Exit 1, `make` meldet für den gescheiterten Bauprozess Exit 2;
`THRESHOLD=65` besteht real (`coverage-gate: OK — Coverage 69.70% erfüllt
Schwelle 65%`, Exit 0).

## Bindung

[`ADR-0071`](../../docs/plan/adr/0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md)
· [`ADR-0054`](../../docs/plan/adr/0054-coverage-gate-und-benchmark-infrastruktur.md)
· `tools/coverage-gate.sh` · `harness/mk/coverage.mk` · seit slice-049.
