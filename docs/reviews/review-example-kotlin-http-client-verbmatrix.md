# Review-Report: examples/kotlin/http-client Verb-Matrix (10 HTTP-API-Fähigkeiten) — 2026-09-28

**Review-Art:** Code — geprüft gegen Plan (Implementer-Bericht der Erweiterung),
den realen Server-Kontrakt (`internal/adapters/driving/http/*.go`), das bereits
abgeschlossene Go-Vorbild (`examples/http-client`) und C#-Vorbild
(`examples/csharp/http-client`, Commits `9336350a`/`d1b4c5f2`), `AGENTS.md` §3
Hard Rules und `.harness/skills/reviewer.md`.

**Gegenstand:** `git diff d1b4c5f2..7a741896 -- examples/kotlin/http-client
examples/README.md` = Commit `7a741896` („feat(examples): Kotlin http-client
deckt alle zehn HTTP-API-Fähigkeiten (`LH-FA-SST-006`)“).

**Skill:** `.harness/skills/reviewer.md` @ `c5207cc1`
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-28

**Eingangs-Kontext** (die Verträge, gegen die geprüft wurde):

- `internal/adapters/driving/http/{server,registerconsumer,consumer,verwaltung,retention,readchanges}.go`
  (realer Server-Kontrakt: Routing, Rechtsklassen, JSON-Felder, Statuscodes)
- `internal/application/port/inbound/readchanges.go` (`ReadChangesQuery`-Vertrag:
  Start inklusiv, End exklusiv)
- `internal/bootstrap/wiring.go` (`administrationTableID`/
  `administrationSchemaVersionID`)
- `examples/http-client/main.go`, `examples/http-client/changes.go` (Form-Vorbild)
- `examples/kotlin/nats-stream-client/build.gradle.kts` (Gson-Pin-Vergleich)
- `docs/reviews/review-example-http-client-verbmatrix.md` (Go-Vorbild-Report,
  F-1…F-4)
- `docs/reviews/review-example-csharp-http-client-verbmatrix.md` (C#-Vorbild-
  Report, F-1 „Herkunft als mehrere Felder, Kette“) und die Folgekorrektur
  `d1b4c5f2`
- `LH-FA-SST-006`, `LH-FA-CON-001/004/005/006`, `LH-FA-CFG-001/002/003`,
  `LH-FA-RET-002…004`, `LH-FA-REA-001 ff.`
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentarform), §3.9
  (Exit-Code-Disziplin)
- `harness/sensors/kommentar-kennungen.md` (Werkzeug-Grenze: nur `.go`)
- `.harness/skills/reviewer.md` (HIGH/MEDIUM-Klassifikation)

Verifikation lief real gegen `make examples-kotlin` (Docker-Layer-Cache
bestätigt: unveränderter Inhalt seit letztem realen Bau, alle 46 Tests liefen
darin bereits grün — Zahl per `grep -rc "@Test"` über alle Testdateien
gegengezählt), `make example-demo-up`/`make example-run-kotlin` (eigener
Rundlauf: `register-consumer`/`table-status`/`remove-consumer` gegen die
bestehende Demo-Umgebung, `enable-table`/`changes`/`disable-table`/
`retention-run` gegen eine frische Wegwerf-Tabelle `public.review_probe`) und
`make kommentar-kennungen` — alle unabhängig vom Implementer-Bericht neu
ausgeführt. `make gates` wurde in diesem Review-Lauf **nicht** erneut
ausgeführt (kein Code außerhalb von `examples/` und `examples/README.md`
berührt, Vertrag prüft ausschließlich Reviewer-Skopus, nicht DoD — das bleibt
Verifier-Aufgabe).

---

## Findings

### F-1 — Datei-Kopfkommentar von `Main.kt` trägt zwei verschiedene Kennungen in einem Block — dieselbe Klasse, die im C#-Vorbild gefunden und dort per `d1b4c5f2` behoben wurde

- `kategorie`: MEDIUM
- `quelle`: `.harness/skills/reviewer.md` MEDIUM „Herkunft als mehrere Felder,
  Kette, „ff.“ oder Spec-Wiederholung“ (`AGENTS.md` §3.7 „Herkunft im
  Go-Kommentar“ — der Wortlaut ist Go-scoped, das Prinzip wird hier per
  Analogie auf den sprachfremden Fall angewendet, den `make
  kommentar-kennungen` strukturell nicht abdeckt, siehe unten)
- `pfad`: `examples/kotlin/http-client/src/main/kotlin/cdcexamples/http/Main.kt:7-20`
- `befund`: Der Kopfkommentar von `Main.kt` trägt in einem einzigen
  `/** … */`-Block sowohl `LH-FA-SST-006` (Zeile 9) als auch `ADR-0087`
  Festlegung 3 (Zeile 12). Das C#-Vorbild trug an derselben Stelle
  (`Program.cs`) exakt dieselbe Zwei-Kennungen-Form, wurde im dortigen Review
  als F-1 gemeldet und in `d1b4c5f2` auf eine einzige Kennung (`LH-FA-SST-006`)
  gekürzt — der `ADR-0087`-Halbsatz wurde ersatzlos gestrichen. Die
  Kotlin-Erweiterung wurde nach dieser Korrektur committet (`7a741896` liegt
  nach `d1b4c5f2`), trägt die bereits identifizierte und bereits behobene Form
  aber unverändert.
- `verifizierbar`: nein für die Form direkt — `make kommentar-kennungen
  PATHS=examples/kotlin/http-client` meldet `0` Kandidaten, aber strukturell
  aus dem falschen Grund: Der Sensor liest laut eigenem Vertrag
  (`harness/sensors/kommentar-kennungen.md` §Grenze Punkt 2 „Nur
  Go-Kommentare“) ausschließlich `.go`-Dateien — `.kt`-Dateien liegen komplett
  außerhalb seines Suchraums, er hat den Block nie gesehen. Die zwei Kennungen
  im selben Block sind durch bloßes Lesen von `Main.kt:7-20` beobachtbar (`grep
  -no` bestätigt: Zeile 9 `LH-FA-SST-006`, Zeile 12 `ADR-0087`).
- `klasse`: „Herkunft als mehrere Felder, Kette“

## Negativbefunde

- geprüft, ohne Befund: **F-1-Vermeidung** (from/to-Semantik) — der
  `ChangesUrlBuilder`-Kommentar (`ChangesClient.kt:43-44`) benennt „from ist
  die untere Positions-Grenze einschließlich, to die obere Grenze
  ausschließlich“ korrekt, deckungsgleich mit
  `internal/application/port/inbound/readchanges.go:17` und
  `examples/http-client/changes.go`; kein Flag-Hilfetext, der etwas anderes
  behauptet
- geprüft, ohne Befund: **F-2-Vermeidung** — `Validator.kt` verlangt für
  `changes` ausschließlich `--source`; `schema`/`table`/`from`/`to`/`limit`
  bleiben unvalidiert optional, kein falscher 400-Kommentar; der
  `ChangesUrlBuilder`-Kommentar beschreibt korrekt, dass ein leerer Wert eines
  bekannten optionalen Parameters serverseitig unberücksichtigt bleibt und mit
  `200` endet (bestätigt gegen `parseReadChangesQuery`/`readChangesPosition`/
  `readChangesLimit` in `readchanges.go`)
- geprüft, ohne Befund: **F-3-Vermeidung** — alle neun neuen Aufruf-Funktionen
  (`registerConsumer`, `acknowledgeConsumer`, `consumerPosition`,
  `removeConsumer`, `enableTable`, `disableTable`, `tableStatus`,
  `readChanges`, `runRetention`) tragen einen eigenen `…FailsOnNon2xx`-Testfall
  mit eigenem Statuscode/Fehlertext (`ConsumerClientTest.kt` ×4,
  `TablesAdminClientTest.kt` ×3, `ChangesClientTest.kt`, `RetentionClientTest.kt`)
- geprüft, ohne Befund: **F-4-Vermeidung** — `HttpClient.newBuilder().build()`
  wird ausschließlich in `Main.kt` konstruiert und über `Dispatcher.dispatch`
  an alle zehn Verben durchgereicht (`grep -rn "HttpClient.new"` findet genau
  diese eine Stelle im gesamten Produktionscode)
- geprüft, ohne Befund: **F-5-Vermeidung (Bestand, Go-Kommentare-Sensor)** —
  `make kommentar-kennungen PATHS=examples/kotlin/http-client` meldet `0`
  Kandidaten; siehe F-1 oben zur Werkzeug-Grenze, die diesen Nullbefund
  strukturell aushöhlt, weil `.kt`-Dateien nicht gelesen werden
- geprüft, ohne Befund: JSON-Feldnamen (`@SerializedName`), Pflicht-/
  Optional-Status, HTTP-Methode, Pfad und erwarteter Erfolgs-Statuscode aller
  zehn Verben gegen die fünf Server-Handler-Dateien — vollständig
  deckungsgleich (inkl. `enableTableRequest`s sieben Felder in identischer
  Reihenfolge)
- geprüft, ohne Befund: Rechtsklassen-Zuordnung aller zehn Verben
  (`ConsumerClient`/`TablesAdminClient`/`RetentionClient`/`ChangesClient`
  nutzen `cfg.token`/`cfg.adminToken` exakt an den Stellen, an denen
  `server.go`s `withToken(...)`-Aufrufe `roleReader`/`roleAdmin` verlangen) und
  gegen `Validator.kt` (`requireReaderToken`/`requireAdminToken` je Verb)
- geprüft, ohne Befund: Gson-Versions-Behauptung — `gson:2.14.0` ist
  byte-identisch zum bereits gepinnten Wert in
  `examples/kotlin/nats-stream-client/build.gradle.kts:45`
- geprüft, ohne Befund: `@PublishedApi internal`-Sichtbarkeit von
  `RequestHelper.requestTimeout`/`.gson` — dies ist die von Kotlin selbst
  dokumentierte, engstmögliche Lösung für den Zugriff eines `public inline
  fun<reified T>` auf klassenlokale Felder über Dateigrenzen hinweg (die
  Alternative wäre, die Felder voll `public` zu machen oder `sendJson` selbst
  auf ein größeres API-Sichtbarkeitsproblem zu heben); die Felder sind über
  das Beispiel-Modul hinaus nicht sinnvoll nutzbar
- geprüft, ohne Befund: `enable-table`s Default-Ableitung von `tableId`/
  `schemaVersionId` (`Dispatcher.applyEnableTableDefaults`) gegen
  `internal/bootstrap/wiring.go`s `administrationTableID`/
  `administrationSchemaVersionID` — bytegleiche Form (`<schema>.<table>` /
  `<table_id>-v1`); real reproduziert gegen `make example-demo-up`:
  `enable-table --source demo-source --schema public --table review_probe
  --publication pub_demo` (ohne `--table-id`/`--schema-version-id`) liefert
  `"table_id":"public.review_probe"`
- geprüft, ohne Befund: reale End-to-End-Probe (eigene, vom Bericht
  unabhängige Wahl, da der Implementer-Bericht `register-consumer`/
  `table-status`/`remove-consumer`/`disable-table` nicht namentlich als
  einzeln verifiziert nennt): `register-consumer` → `table-status` (bestehende
  Tabelle `public.orders`) → `remove-consumer` (inkl. Idempotenz-Gegenprobe:
  zweiter Aufruf liefert `removed:false`) → `enable-table` (frische Tabelle
  `public.review_probe`, Default-Ableitung) → `disable-table` →
  erneutes `table-status` (`enabled:false`) → `retention-run` — alle sieben
  Antworten stimmen exakt mit dem jeweiligen Server-Response-Schema überein
- geprüft, ohne Befund: **Gson-Null-Weglassungs-Kosmetikbefund** — real
  reproduziert: `GET /changes` gegen die Demo-Zeile liefert im Server-JSON
  `"old_image":null` (`readChangesHandler`), die von `Main.kt` gedruckte
  Kotlin-Antwort lässt das Feld `old_image` vollständig weg (kein
  `"old_image":null`, kein `"old_image":` überhaupt). Ursache bestätigt:
  `TypeAdapters.JSON_ELEMENT` dekodiert JSON-`null` zu `JsonNull.INSTANCE`
  (kein Java-`null`), aber `JsonWriter.nullValue()` verwirft beim Schreiben
  den bereits gepufferten Feldnamen, wenn `serializeNulls` (Default: aus)
  nicht gesetzt ist — unabhängig davon, ob der Feldwert Java-`null` oder ein
  `JsonElement`, das `JsonNull` repräsentiert, ist. Die Einschätzung des
  Implementers (kein Effekt auf das **Parsen** der echten Server-Antwort) ist
  zutreffend — `ChangesClientTest.readChangesParsesSuccessResponseWithEmbeddedImages`
  bindet bereits `oldImage.isJsonNull` gegen eine Server-Antwort mit
  `"old_image":null` erfolgreich. Der Unterschied betrifft ausschließlich die
  **gedruckte** Konsolen-Ausgabe dieses Beispiels und ist damit stilistisch,
  nicht funktional — kein eigenes Finding (siehe unten zur fehlenden
  Dokumentation dieser Divergenz als separate Beobachtung)
- geprüft, ohne Befund: `TablesClientTest.kt`/`CliTest.kt` nach der
  `Config`-Erweiterung — alle `Config(...)`-Aufrufe im gesamten Testbaum
  nutzen benannte Argumente, kein stiller Positionsversatz
- geprüft, ohne Befund: Docker-only (§3.1) — keine Host-Kotlin-/Gradle-Aufrufe
  im Diff außerhalb von `make examples-kotlin`/`make example-run-kotlin`, kein
  `sed -i`/keine Umleitung auf eine Repo-Datei
- geprüft, ohne Befund: Arbeitsbaum nach dem Commit sauber (`git status`
  clean), Commit `7a741896` trägt ausschließlich `examples/kotlin/http-client/**`
  und `examples/README.md`
- geprüft, ohne Befund: Traceability — Commit-Betreff nennt `LH-FA-SST-006`,
  Footer nennt `ADR-0090`
- geprüft, ohne Befund: `README.md`-Ergänzung nennt exakt die neun ergänzten
  Verben, deckungsgleich mit `Validator.knownVerbs` minus `tables`, und ist
  wörtlich parallel zur bereits akzeptierten C#-Zeile formuliert
- geprüft, ohne Befund: Rückwärtskompatibilität — Default-Verb ist `tables`,
  `TablesClient.kt`/`TablesUrlBuilder.kt` sind im Diff unverändert (kein
  Eintrag im `git diff --stat`)
- geprüft, ohne Befund: 46 Tests (`grep -rc "@Test"` über alle Testdateien
  summiert: 46) — deckungsgleich mit der Implementer-Angabe

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** Herkunft als mehrere Felder, Kette

## Verdikt

**Merge-blockierend:** nein — dieselbe Einschätzung wie beim C#-Vorbild-Review:
ein einzelnes MEDIUM, mechanisch günstig zu beheben (den `ADR-0087`-Halbsatz
aus dem Kopfkommentar entfernen, wie in `d1b4c5f2` für C# bereits vorgemacht),
kein Sicherheits-, Korrektheits- oder Vertragsproblem. Alle vier im
Go-Vorbild-Review gefundenen Fehlerklassen (F-1…F-4, davon zwei HIGH) sind in
dieser Kotlin-Erweiterung nachweislich vermieden, real gegenverifiziert statt
dem Implementer-Bericht geglaubt (inklusive einer eigenen, vom Bericht
unabhängigen End-to-End-Probe gegen `make example-demo-up`). Bemerkenswert:
Die bereits im C#-Review gefundene und dort per Folgecommit behobene
Kopfkommentar-Form (F-1 dort) wurde in der Kotlin-Erweiterung erneut
eingeführt — der Implementer hat offenbar den C#-Kopfkommentar **vor** der
Korrektur als Formvorbild kopiert, nicht die bereits bereinigte Fassung, und
`make kommentar-kennungen` kann das nicht auffangen, weil sein Suchraum
`.kt`-Dateien strukturell nicht trägt.

**Übergabe:** Das MEDIUM-Finding geht an den Implementer zur Fixrunde. Die
Finding-Klasse „Herkunft als mehrere Felder, Kette“ ist bereits als etablierte
Klasse im Reviewer-Skill geführt (`BEO-PGC/kommentar-herkunft-als-kette`, seit
slice-code-kommentare-kennungen); dies ist ihr drittes benanntes Vorkommen in
dieser Slice-Familie (Go-Review, C#-Review F-1, jetzt Kotlin-Review F-1) — bei
gleichem Konflikttyp zum dritten Mal greift der Konflikt-Pfad über den
Architect (Modul 8): Wird `make kommentar-kennungen` nicht auf `.kt`/`.cs`
erweitert, wiederholt sich dieser exakte Fund bei jeder künftigen
SDK-/Beispiel-Sprachwurzel, die einen bereits gefundenen und korrigierten
Kommentar unwissentlich vom **vorherigen** statt vom **korrigierten** Stand
des Formvorbilds kopiert. Dieser Report ersetzt keine Verifikation —
DoD-/Spec-Konformität prüft der Verifier separat (Modul 11).
