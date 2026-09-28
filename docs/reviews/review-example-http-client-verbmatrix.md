# Review-Report: examples/http-client Verb-Matrix (10 HTTP-API-Fähigkeiten) — 2026-09-28

**Review-Art:** Code — geprüft gegen Plan (Implementer-Bericht der Erweiterung),
den realen Server-Kontrakt (`internal/adapters/driving/http/*.go`), `AGENTS.md`
§3 Hard Rules und `.harness/skills/reviewer.md`.

**Gegenstand:** `git diff be2d23e7..HEAD -- examples/http-client examples/README.md`
= Commit `b0948065` („feat(examples): Go http-client deckt alle zehn
HTTP-API-Fähigkeiten (LH-FA-SST-006)“, mit bereits eingefaltetem
Kommentar-`--amend`).

**Skill:** `.harness/skills/reviewer.md` @ `c5207cc1`
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-28

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `internal/adapters/driving/http/{server,registerconsumer,consumer,verwaltung,retention,readchanges,middleware}.go`
  (realer Server-Kontrakt: Routing, Rechtsklassen, JSON-Felder, Statuscodes)
- `internal/application/port/inbound/readchanges.go` (`ReadChangesQuery`-Vertrag:
  Start inklusiv, End exklusiv)
- `internal/bootstrap/wiring.go` (`administrationTableID`/
  `administrationSchemaVersionID`)
- `docs/user/benutzerhandbuch.md` §„Changes lesen“ (dokumentierte
  `from`/`to`-Semantik)
- `LH-FA-SST-006`, `LH-FA-CON-001/004/005/006`, `LH-FA-CFG-001/002/003`,
  `LH-FA-RET-002…004`, `LH-FA-REA-001 ff.`
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentarform), §3.9
  (Exit-Code-Disziplin)
- `.harness/skills/reviewer.md` (HIGH/MEDIUM-Klassifikation)

Verifikation lief real gegen `make example-demo-up` (Docker-only, eigener
Rundlauf, kein `docker exec` auf fremde Container), gegen den mit
`golang:1.27-alpine@sha256:cf6fca…` (repo-eigenes `TOOLCHAIN_IMAGE`) gebauten
Testlauf, sowie `make fmt-check`, `make kommentar-kennungen` und
`make docs-check` — alle unabhängig vom Implementer-Bericht neu ausgeführt.

---

## Findings

### F-1 — Flag-Hilfetext für `-from`/`-to` behauptet die genau umgekehrte Positionsgrenzen-Semantik

- `kategorie`: HIGH
- `quelle`: `.harness/skills/reviewer.md` HIGH „Kommentar trägt keine der
  Kommentar-Klassen“ (Unterpunkt Zusage); Server-Kontrakt
  `internal/application/port/inbound/readchanges.go:17` („Start **inklusiv**,
  End **exklusiv**“), `docs/user/benutzerhandbuch.md:1135`
  („eingegrenzt über `from` (inklusive) und `to` (exklusive)“)
- `pfad`: `examples/http-client/main.go:112-113`
- `befund`: Die Flag-Beschreibungen lauten `-from` „Untere Positions-Grenze,
  **ausschließlich**“ und `-to` „Obere Positions-Grenze, **einschließlich**“
  — genau umgekehrt zum real gemessenen und im Handbuch dokumentierten
  Server-Verhalten. Real reproduziert gegen `make example-demo-up` (vier
  Changes, Commit-Positionen `30206032`/`30560384`): `-from=30206032`
  liefert die Change an genau dieser Position mit (inklusiv), `-from=30206033`
  schließt sie aus; `-to=30560384` schließt die Changes an genau dieser
  Position aus (exklusiv), `-to=30560385` schließt sie ein. Ein Nutzer, der
  sich auf den Flag-Text verlässt, erhält bei einer Grenzabfrage eine um
  genau eine Transaktion verschobene, falsche Ergebnismenge, ohne dass ein
  Fehler sichtbar wird.
- `verifizierbar`: ja — reproduzierbar über `make example-demo-up` +
  `make example-run-go SURFACE=http ARGS="-verb=changes -source demo-source
  -from=<pos>"` bzw. `-to=<pos>` gegen bekannte `commit_position`-Werte (im
  Review real ausgeführt, vier Messungen, alle bestätigend).
- `klasse`: „Zusage widerspricht dem gemessenen Verhalten“

### F-2 — Kommentar zu `setIfNotEmpty` behauptet einen 400 für „leere“ Parameter, den der Server nicht liefert

- `kategorie`: HIGH
- `quelle`: `.harness/skills/reviewer.md` HIGH „Kommentar trägt keine der
  Kommentar-Klassen“ (Unterpunkt Zusage); Server-Kontrakt
  `internal/adapters/driving/http/readchanges.go` (`parseReadChangesQuery`,
  `readChangesPosition`, `readChangesLimit`: ein leerer Wert eines bekannten
  Parameters trägt „keine Grenze“, kein Fehler)
- `pfad`: `examples/http-client/changes.go:75-77`
- `befund`: Der Kommentar sagt zu: „`GET /changes` trägt eine geschlossene
  Parameter-Menge und lässt einen unbekannten **oder leeren** Parameter mit
  `400` enden.“ Real gegen den laufenden Feed-Container gemessen (`curl` über
  das Docker-Netz `cdc-examples`, `Authorization: Bearer
  demo-reader-token`): `GET /changes?source=demo-source&schema=` → `200`,
  `&limit=` → `200`, `&from=` → `200` (Wert bleibt unberücksichtigt, keine
  Fehlermeldung). Nur ein unbekannter Parameter**name** (`&unbekannt=1` →
  `400`) und der spezifische, eigens geprüfte Pflichtfeld-Fall `source=`
  leer (`400`, „source ist Pflichtfeld“) enden sichtbar — das ist ein
  eigener, nicht verallgemeinerbarer Check auf genau einem Feld, keine
  generische „leerer Parameter“-Regel. Die Zusage ist für die vier übrigen
  optionalen Parameter (`schema`, `table`, `from`, `to`, `limit`) falsch;
  praktische Auswirkung ist gering, weil der Client selbst nie einen leeren
  Wert sendet (`setIfNotEmpty` lässt ihn weg) — der Kommentar beschreibt aber
  ein Server-Verhalten falsch, das ein Leser dieses Vorbild-Programms für
  bare Münze nimmt.
- `verifizierbar`: ja — reproduzierbar über `make example-demo-up` + direkte
  `curl`-Aufrufe gegen `GET /changes` mit leeren Query-Werten (im Review real
  ausgeführt, vier Messungen: zwei `200`, zwei `400` an den erwarteten
  Stellen).
- `klasse`: „Zusage widerspricht dem gemessenen Verhalten“

### F-3 — Fünf von neun neuen Aufruf-Funktionen ohne eigenen Nicht-2xx-Test

- `kategorie`: MEDIUM
- `quelle`: `.harness/skills/reviewer.md` MEDIUM „fehlende Negativtests bei
  neuem öffentlichem Vertrag“
- `pfad`: `examples/http-client/consumer_test.go`,
  `examples/http-client/tables_admin_test.go`
- `befund`: `acknowledgeConsumer`, `consumerPosition`, `removeConsumer`,
  `disableTable` und `tableStatus` haben je nur einen Erfolgstest; ein
  eigener Nicht-2xx-Test fehlt (anders als bei `readChanges`,
  `registerConsumer`, `enableTable`, `runRetention`, die je einen
  `TestXFailsOnNon2xx` tragen, zusätzlich zum generischen
  `TestDoRequestJSONFailsOnUnexpectedStatus` in `request_test.go`). Der
  Fehlerpfad selbst ist über den gemeinsamen `doRequestJSON`-Helfer bereits
  einmal generisch getestet — die fünf fehlenden Tests sind also aus reiner
  Code-Pfad-Sicht redundant, aber die Skill-Regel adressiert genau diese
  Lücke: ein neuer öffentlicher Vertrag (jede der zehn Aufruf-Funktionen ist
  eine eigenständige, exportierbare Fähigkeit dieses Beispiels) ohne eigenen
  Negativtest.
- `verifizierbar`: ja — `go test -v ./examples/http-client/...` (im Review
  real ausgeführt: 30 Tests, alle grün, keine der fünf genannten Funktionen
  hat einen `Non2xx`-benannten Testfall).
- `klasse`: „Fehlende Negativtests bei neuem öffentlichem Vertrag“

### F-4 — `listTables` reicht den in `dispatch` erzeugten `*http.Client` nicht durch

- `kategorie`: LOW
- `quelle`: Maintainability
- `pfad`: `examples/http-client/main.go` (`dispatch`, Fall `"tables"`;
  `listTables`)
- `befund`: `dispatch` erzeugt einen `*http.Client` und reicht ihn an alle
  neun neuen Aufruf-Funktionen durch; der Fall `"tables"` ruft stattdessen
  weiterhin `listTables(cfg)` auf, das intern einen zweiten,
  unabhängigen `*http.Client` mit demselben `requestTimeout` baut. Stilistisch
  inkonsistent zu den neun übrigen Verben, ohne funktionalen Unterschied.
- `verifizierbar`: nein — Lese-Handlung, kein Gate-Lauf zeigt das.
- `klasse`: „Stilistische Client-Inkonsistenz“

## Negativbefunde

- geprüft, ohne Befund: Rechtsklassen-Zuordnung aller zehn Verben gegen
  `server.go`s `withToken`-Aufrufe (reader/admin je Endpunkt 1:1 deckungsgleich)
- geprüft, ohne Befund: JSON-Feldnamen, Pflicht-/Optional-Status, HTTP-Methode,
  Pfad und erwarteter Erfolgs-Statuscode aller zehn Verben gegen die fünf
  Server-Handler-Dateien (`registerconsumer.go`, `consumer.go`,
  `verwaltung.go`, `retention.go`, `readchanges.go`) — vollständig
  deckungsgleich
- geprüft, ohne Befund: `enable-table`s Default-Ableitung von
  `table_id`/`schema_version_id` gegen `internal/bootstrap/wiring.go`s
  `administrationTableID`/`administrationSchemaVersionID` — bytegleiche Form
  (`<schema>.<table>` / `<table_id>-v1`)
- geprüft, ohne Befund: unbekanntes `-verb` bricht in `validate()` vor jedem
  Netzwerkaufruf mit `os.Exit(2)` ab (`main.go:76-79`, `knownVerbs`-Prüfung
  als erster Schritt in `validate`); `dispatch` trägt denselben Schutz
  defensiv ein zweites Mal
- geprüft, ohne Befund: Rückwärtskompatibilität — `examples/http-client/tables.go`
  und `tables_test.go` sind im Diff `be2d23e7..HEAD` unverändert (leerer
  `git diff`), Default-Verb ist `tables`, bestehende Startform
  `ARGS="-source <quelle> -publication <publication>"` bleibt funktionsfähig
- geprüft, ohne Befund: Fehlerbehandlung bei Nicht-2xx — alle zehn Verben
  laufen durch denselben `doRequestJSON`- bzw. `listTables`-Pfad
  (Statuscode + Antworttext auf stderr, `os.Exit(1)`), konsistente Form
- geprüft, ohne Befund: `make fmt-check` (real ausgeführt: 271 Go-Dateien
  geprüft, alle formatiert)
- geprüft, ohne Befund: `make kommentar-kennungen PATHS=examples/http-client
  COUNT=1` (real ausgeführt: 0 Kandidaten — bestätigt, dass der
  Implementer-Fund im `main.go`-Kopfkommentar tatsächlich auf eine Kennung
  reduziert wurde)
- geprüft, ohne Befund: `make docs-check` (real ausgeführt: 1369 Dateien
  geprüft, 0 Befunde — die neue README-Zeile bricht keine Referenz)
- geprüft, ohne Befund: Docker-only (§3.1) — keine Host-Toolchain-Aufrufe im
  Diff, reine Go-Quellen und Tests
- geprüft, ohne Befund: Arbeitsbaum nach dem Commit sauber
  (`git status` clean), Commit `b0948065` trägt ausschließlich
  `examples/http-client/**` und `examples/README.md`
- geprüft, ohne Befund: Traceability — Commit-Betreff nennt `LH-FA-SST-006`
- geprüft, ohne Befund: Go-Unit-Tests des Pakets (`go test -v
  ./examples/http-client/...` gegen das repo-eigene `TOOLCHAIN_IMAGE`,
  real ausgeführt: 30/30 grün)
- geprüft, ohne Befund: reale E2E-Probe von `changes` mit gesetzten
  `-limit`/`-from`/`-to` gegen `make example-demo-up` — laut
  Implementer-Bericht die am dünnsten belegte Fläche; die Werte selbst
  wurden korrekt an den Server durchgereicht (siehe aber F-1/F-2 zur
  Hilfetext-/Kommentar-Beschreibung dieser Werte)

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 1 |
| LOW | 1 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Zusage widerspricht dem gemessenen Verhalten (2×) ·
Fehlende Negativtests bei neuem öffentlichem Vertrag · Stilistische Client-Inkonsistenz

## Verdikt

**Merge-blockierend:** ja — zwei HIGH-Findings (F-1, F-2) sind reale,
verifizierte Fehlbeschreibungen der Server-Semantik in einem als „Vorbild
(minimal, lesbar)“ deklarierten öffentlichen Beispielprogramm; beide sind
mechanisch günstig zu beheben (Flag-Hilfetext und Kommentar korrigieren),
aber sie gehen an den Implementer zurück, keine Direktkorrektur durch den
Reviewer.

**Übergabe:** F-1 bis F-4 gehen an den Implementer zurück (Fixrunde nötig).
Die Finding-Klassen gehen zusätzlich in die Slice-Closure §7 und von dort in
den Steering-Loop-Zähler des Reviewer-Skills. Dieser Report ist ein
Lauf-Beleg; er ersetzt keine Verifikation gegen die DoD (Verifier-Aufgabe,
Modul 11).
