# Review-Report: Tabellen-granulare Filterung für gRPC-Change-Stream und SSE ([ADR-0133](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md)) — 2026-09-28

**Review-Art:** Code — geprüft gegen Plan/ADR + Konventionen (Modul 10
§Drei Review-Arten). Die DoD-Frage bleibt beim Verifier.

**Gegenstand:** `git diff d987ca8f..ad7db2f6` — vier Commits: `d85f6ba3`
(`schema`/`table`-Filterpaar als Domänen-Prüfung `Change.MatchesFilter` +
Protobuf-Feld + erzeugter Code), `bb6d4f8f` (Filterprüfung im gRPC- und
SSE-Handler + Unit-Tests), `9530ca4a` (`SPEC-020`/`SPEC-021` +
Benutzerhandbuch), `ad7db2f6` (E2E-Filter-Rundlauf in
`run-integration-tests.sh`, `grpcclient`/`sseclient` um Filter-Argumente
erweitert).

**Skill:** `.harness/skills/reviewer.md` @ Stand 2026-09-09 (geschärft,
vier repo-spezifische HIGH-Regeln)
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-28

**Eingangs-Kontext:**

- [`docs/plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) (Accepted) — vollständig gelesen, alle sechs Teilfragen, Fitness-Function-Tabelle
- [`docs/plan/adr/0060-grpc-streaming-mechanismus.md`](../plan/adr/0060-grpc-streaming-mechanismus.md), [`docs/plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md`](../plan/adr/0061-http-sse-zusaetzlich-zu-grpc.md) (teilweise superseded — Vertagung der Filterung)
- [`docs/plan/adr/0081-changes-lesen-ueber-die-http-api.md`](../plan/adr/0081-changes-lesen-ueber-die-http-api.md) (Formvorbild `schema`/`table`-Filter)
- `AGENTS.md` §3 Hard Rules, insbesondere §3.1, §3.7, §3.9, §3.12, §3.13
- `.harness/skills/reviewer.md` vollständig

---

## Vorgehen (zusammengefasst)

- [ADR-0133](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) vollständig gelesen (sechs Teilfragen, Konsequenzen,
  Fitness-Function, Slice-Schnitt-Empfehlung, Re-Evaluierungs-Trigger).
- `proto/cdc/stream/v1/changestream.proto`-Diff gegen Teilfrage 2 gehalten:
  zwei additive Felder (`schema=1`, `table=2`) auf der zuvor leeren
  Message — wire-kompatibel.
- `internal/domain/model/change.go` (`MatchesFilter`) gelesen: reine,
  seiteneffektfreie Funktion, geteilt zwischen gRPC- und SSE-Handler
  (Teilfrage 3/4-Anforderung erfüllt).
- `internal/adapters/driving/grpc/server.go` und
  `internal/adapters/driving/http/sse.go` gegen Teilfrage 3 (Ansatzpunkt
  nach `Subscribe()`) und Teilfrage 4 (SSE-Query-Parameter, `400` vor
  Öffnen des Streams) gehalten.
- `git diff` gegen `internal/adapters/driven/grpcstream/` und
  `internal/application/port/` selbst ausgeführt: **leer** —
  `ChangeStreamPort`/`Broadcaster` tatsächlich unverändert, wie die ADR
  verlangt.
- Unit-Test-Diffs (`change_test.go`, `server_test.go`, `sse_test.go`)
  gegen die HIGH-Klasse „Zusage ohne Bindung an ihre Eingabeseite" geprüft:
  die neuen Fitness-Function-Tests mutieren den **zugestellten Change**
  (Schema/Tabelle), nicht nur die Filterwerte — Eingabeseiten-Bindung
  bestätigt.
- `spec/pflichtenheft.md`- und `docs/user/benutzerhandbuch.md`-Diffs
  gelesen: Korrektur sachlich zutreffend, Handbuch-Versionshistorie
  (1.73→1.74) im selben Commit nachgezogen, Beispiel-/SDK-Lücke explizit
  benannt.
- `tools/harness/grpcclient/main.go`/`sseclient/main.go`-Diffs gegen
  Regressionsfreiheit des bestehenden 2-/3-Argument-Aufrufs geprüft.
- `make image` (Cache-Treffer, `:dev` bereits aktuell), `make
  generated-sync` (OK, byte-gleich), `make kommentar-kennungen
  DIFF=d987ca8f` (0 Kandidaten), `make fmt-check` (285 Dateien, alle
  formatiert), `make a-check` (0 Befunde), `make docs-check` (0 Befunde),
  `make commit-traceability` (OK), `make baseline-verify` (OK) — alle
  selbst ausgeführt, Exit-Codes direkt geprüft.
- `make gates` zweimal real ausgeführt (Exit 0 beide Male, direkt
  geprüft): Coverage 80,60 % ≥ 80 % Schwelle.
- Die vier neuen Fitness-Function-Unit-Tests
  (`TestChangeMatchesFilter`, `TestStreamChangesFilterLaesstNurTreffer`,
  `TestStreamChangesOhneFilterLiefertAlle`,
  `TestStreamFilterLaesstNurTreffer`, `TestStreamOhneFilterLiefertAlle`,
  `TestStreamUnbekannterQueryParameterEndetMit400`) einzeln mit `-v`
  gefahren: alle PASS.
- `make test-integration` real ausgeführt (Hintergrundprozess, ~11
  Minuten bis zum Abbruch) — **Ergebnis: Exit 1, real reproduziert, siehe
  F-1**. Die beiden neuen Filter-Rundläufe selbst liefen dabei
  **erfolgreich** (beide `abdeckung_declare`-Anker-Selbstprüfungen und
  beide Beleg-Echos real gedruckt, inkl. echtem Negativbeleg — siehe
  Bestätigte Implementer-Behauptungen). Der Abbruch liegt **später** im
  selben Lauf, an einer von F-1 beschriebenen ID-Kollision. Compose-
  Teardown danach sauber (`docker ps -a` ohne `cdc-test-*`-Reste),
  `git status --short` nach dem Lauf leer (`docs/user/e2e-abdeckung.md`
  unangetastet, weil der Lauf vor `abdeckung_schreiben` abbrach).

---

## Findings

### F-1 — `make test-integration` scheitert real: die beiden neuen Filter-Negativbelege schreiben harte IDs in eine von einem späteren, unveränderten Go-Test exklusiv beanspruchte Tabelle — Primärschlüssel-Kollision

- `kategorie`: HIGH
- `quelle`: [ADR-0133](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) Fitness-Function-Zeile „`make test-integration`" (nicht erfüllt, weil der Gesamtlauf scheitert); `AGENTS.md` §3.13 („Eine Arbeit, die eine beschriebene Eigenschaft bewegt, zieht ihre Träger nach" — hier ist es kein Träger-Text, sondern ein realer Datenzustand derselben Tabelle, den der neue Testblock nicht gegen die späteren, bereits bestehenden Verbraucher dieser Tabelle geprüft hat)
- `pfad`: `tools/harness/run-integration-tests.sh:2502` (Kommentar „id=10 auf feed_e2e_schema — dort bereits 1/2 belegt"), `tools/harness/run-integration-tests.sh:2535` (`INSERT INTO public.$GRPC_FILTER_OTHER_TABLE (id, name) VALUES (10, …)`, `GRPC_FILTER_OTHER_TABLE=feed_e2e_schema`), `tools/harness/run-integration-tests.sh:2945` (Kommentar „id=11 auf feed_e2e_schema, getrennt vom id=10 des gRPC-Filter-Rundlaufs oben"), `tools/harness/run-integration-tests.sh:2979` (`INSERT INTO public.$SSE_FILTER_OTHER_TABLE (id, name) VALUES (11, …)`, `SSE_FILTER_OTHER_TABLE=feed_e2e_schema`) — gegen `test/integration/integration_test.go:1176` (`INSERT INTO … VALUES (10, 'Reject', …)`) und `test/integration/integration_test.go:1202` (`INSERT INTO … VALUES (11, 'TypeChanged', …)`), beide in `TestE2ESchemaChangeIncompatibleTypeChange` — eine bestehende, von diesem Diff **nicht berührte** Funktion.
- `befund`: Die beiden neuen Filter-Rundläufe wählen `feed_e2e_schema` als
  „eine andere, aktivierte Tabelle", gegen die der gefilterte Client
  **nichts** empfangen soll (Negativbeleg). Dafür fügen sie real Zeilen
  mit `id=10` (gRPC-Block) und `id=11` (SSE-Block) in `feed_e2e_schema`
  ein und lassen sie dort **stehen** — keine der beiden Sektionen räumt
  danach auf. Der Implementer-Kommentar an beiden Stellen benennt
  explizit einen Wertebereichs-Check, aber nur gegen die **früher**
  laufenden Tests derselben Tabelle („dort bereits 1/2 belegt", d. h.
  `TestE2ESchemaChangeAddColumn`s `id=1`/`id=2`) — nicht gegen den
  **später**, im selben Compose-Lauf noch ausstehenden
  `TestE2ESchemaChangeIncompatibleTypeChange`, der auf derselben Tabelle
  hart `id=10` (Fall 1: PostgreSQL soll die DDL selbst ablehnen) und
  `id=11` (Fall 2: die Zeile nach der zugelassenen Typänderung) einfügt.
  Beide IDs sind durch den neuen Diff bereits belegt, wenn dieser Test
  läuft — der `INSERT` von `id=10` scheitert real mit
  `duplicate key value violates unique constraint "feed_e2e_schema_pkey"
  (SQLSTATE 23505)`, der Test schlägt fehl, `test-integration` bricht mit
  Exit 1 ab, **bevor** die drei nachfolgenden `abdeckung_declare`-Phasen
  (u. a. „Transformationen-Nichtanwendbarkeit und Abhilfe") und der
  abschließende `abdeckung_schreiben`-Schritt je erreicht werden — deshalb
  bleibt `docs/user/e2e-abdeckung.md` von diesem Diff unangetastet, obwohl
  zwei neue Phasen real gelaufen sind.
- `verifizierbar`: ja — selbst reproduziert. `make image` (Cache-Treffer
  gegen den committeten Endstand `ad7db2f6`), danach `make
  test-integration`: der Lauf läuft ~11 Minuten real durch Kern-CDC-,
  Diagnose-, gRPC-/SSE-Stream- (inkl. beider neuer Filter-Rundläufe, siehe
  unten), NATS-, Backfill-, Leerlauf-Bestätigungs- und
  Transformations-Belege, dann:
  ```
  === RUN   TestE2ESchemaChangeIncompatibleTypeChange
      integration_test.go:1177: INSERT der nicht konvertierbaren Zeile: ERROR: duplicate key value violates unique constraint "feed_e2e_schema_pkey" (SQLSTATE 23505)
  --- FAIL: TestE2ESchemaChangeIncompatibleTypeChange (0.02s)
  FAIL
  FAIL	github.com/pt9912/pg-change-feed/test/integration	0.021s
  FAIL
  make: *** [Makefile:227: test-integration] Fehler 1
  ```
  `git diff d987ca8f..ad7db2f6 -- test/integration/integration_test.go`
  zeigt: diese Datei ist **nicht Teil des Diffs** — der Testcode selbst
  ist unverändert richtig, die Kollision entsteht ausschließlich durch die
  neuen `INSERT`-Zeilen des Diffs in derselben Tabelle.
- `klasse`: ID-Kollision zwischen neuem E2E-Testdatensatz und bestehendem Go-Test in derselben Tabelle (Wertebereichs-Prüfung nur rückwärts, nicht vorwärts im selben Lauf)

---

## Bestätigte Implementer-Behauptungen (eigenständig nachgeprüft)

- **Wire-Kompatibilität (Teilfrage 2):** `StreamChangesRequest` trägt
  genau zwei neue, additive `string`-Felder (`schema=1`, `table=2`) auf
  einer zuvor leeren Message; `make generated-sync` bestätigt
  Byte-Gleichheit von `gen/cdc/stream/v1/changestream.pb.go` gegen den
  gepinnten Generator. `changestream_grpc.pb.go` unverändert (kein neuer
  RPC, wie erwartet).
- **`ChangeStreamPort`/`Broadcaster` unverändert (Teilfrage 3):**
  `git diff d987ca8f..ad7db2f6 -- internal/adapters/driven/grpcstream/
  internal/application/port/` ist leer — die Filterung sitzt
  ausschließlich in den beiden Driving-Handlern, wie die ADR verlangt.
- **Geteilte Prüf-Funktion (Teilfrage 3/4, `LH-FA-SST-006`):**
  `model.Change.MatchesFilter` ist eine einzige, reine Funktion, die
  sowohl `grpc/server.go` als auch `http/sse.go` aufrufen — keine
  unabhängig driftende Doppelimplementierung.
- **Filter-Semantik (Teilfrage 1):** `TestChangeMatchesFilter` deckt alle
  acht Kombinationen des unabhängig optionalen `schema`/`table`-Paars
  (`change_test.go`); Regressionstests (`TestStreamChangesOhneFilterLiefertAlle`,
  `TestStreamOhneFilterLiefertAlle`) bestätigen, dass eine leere Request/
  Query weiterhin alle Changes liefert.
- **Fehlerform SSE (Teilfrage 4):** ein unbekannter Query-Parameter endet
  weiterhin mit `400`, geprüft **vor** `Subscribe()` und vor dem
  `200`-Header (`sse.go:121-125` vs. `142`ff.); Unit-Test
  `TestStreamUnbekannterQueryParameterEndetMit400` bindet das zusätzlich
  daran, dass der Handler sich in diesem Fall **nicht** am Broadcaster
  registriert.
- **Eingabeseiten-Bindung der Mutation (Commit-Message `bb6d4f8f`):**
  eigenständig nachgeprüft — sowohl `TestStreamChangesFilterLaesstNurTreffer`
  als auch `TestStreamFilterLaesstNurTreffer` mutieren den **zugestellten
  Change** (`nichtPassend.Schema`/`.Table` bzw. dieselben Felder), nicht
  die Filterwerte der Request/Query. Ein Fehler in `MatchesFilter` selbst
  (z. B. eine der beiden `if`-Bedingungen entfernt) würde diese Tests real
  rot färben — die Behauptung „Eingabeseite ist der zugestellte Change,
  nicht die Filterwerte" ist korrekt.
- **`spec/pflichtenheft.md`:** `SPEC-020`/`SPEC-021`-Korrektur sachlich
  zutreffend, `LH-FA-CFG-008.a` korrekt auf die verbleibende offene Frage
  (serverseitig konfigurierbares Routing-Zielmodell) verengt, kein
  ADR-Verweis im Fließtext (nur in der Changelog-Zeile, konsistent mit der
  etablierten Konvention).
- **`docs/user/benutzerhandbuch.md`:** Version 1.73→1.74 im Kopf **und**
  neue Zeile in der Änderungshistorie im selben Commit (`9530ca4a`); beide
  „Filterung nach Tabelle ist nicht Teil dieser Version"-Sätze ersetzt;
  beide neuen Absätze benennen die noch nicht angepassten Beispiel-Clients
  und SDK-Packages explizit als offenen Folge-Schritt (nicht
  stillschweigend ausgelassen) — konsistent mit der ADR-Folgepflicht.
- **`tools/harness/grpcclient/main.go`/`sseclient/main.go`:** die
  Argumentzahl-Prüfung akzeptiert weiterhin ausschließlich 3 oder 5
  Argumente — ein bestehender 3-Argument-Aufruf (alle bisherigen
  Rundläufe, die diese Clients nutzen) bleibt unverändert lauffähig.
- **`spec/architecture.md`/`.a-check.yml` unverändert:** `git diff` gegen
  beide Dateien ist leer, wie die ADR-Folgepflicht „braucht keine
  Änderung" behauptet; `make a-check` bestätigt 0 Befunde.
- **Gates/Sensoren — alle selbst gefahren, Exit-Codes direkt geprüft:**
  `make fmt-check` (285 Go-Dateien, alle formatiert), `make a-check`
  (0 Befunde), `make generated-sync` (OK, beide `.proto`-Quellen), `make
  docs-check` (1385 Dateien, 0 Befunde), `make commit-traceability` (OK,
  5 Commits), `make baseline-verify` (OK), `make gates` (Exit 0, zweimal
  geprüft; Coverage 80,60 % ≥ 80 % Schwelle).
- **Traceability:** alle vier Commit-Messages nennen `ADR-0133`, kein
  `SPEC-*`/`ARC-*` im Betreff.
- **`make kommentar-kennungen DIFF=d987ca8f`:** 0 Kandidaten — jeder neu
  geschriebene Kommentarblock trägt genau eine Kennung (`ADR-0133`).

---

## Negativbefunde

- geprüft, ohne Befund: `internal/domain/model/change.go`/`change_test.go`
  (`MatchesFilter`, reine Funktion, acht Kombinationsfälle, Kommentar
  trägt genau eine Kennung)
- geprüft, ohne Befund: `internal/adapters/driving/grpc/server.go`/
  `server_test.go` (Filterprüfung nach `Subscribe()`, Regressionstest,
  Eingabeseiten-Bindung)
- geprüft, mit Befund F-1 (Kollision, nicht in diesen Dateien selbst,
  sondern durch ihre Interaktion mit `test/integration/integration_test.go`):
  `internal/adapters/driving/http/sse.go`/`sse_test.go` selbst sind
  korrekt implementiert (Filterprüfung, `400`-Pfad vor Öffnen des Streams,
  Regressionstest, Eingabeseiten-Bindung)
- geprüft, ohne Befund: `proto/cdc/stream/v1/changestream.proto`,
  `gen/cdc/stream/v1/changestream.pb.go` (additiv, byte-gleich zur
  Generator-Ausgabe)
- geprüft, ohne Befund: `spec/pflichtenheft.md` (`SPEC-020`/`SPEC-021`,
  `LH-FA-CFG-008.a`)
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md`
  (Versionshistorie, beide Filter-Absätze, Beispiele-/SDK-Hinweis)
- geprüft, ohne Befund: `tools/harness/grpcclient/main.go`,
  `tools/harness/sseclient/main.go` (Regressionsfreiheit des bestehenden
  Aufrufs)
- geprüft, mit Befund F-1: `tools/harness/run-integration-tests.sh` (die
  beiden neuen Filter-Rundläufe selbst laufen korrekt und liefern die
  erwarteten Negativ-/Positiv-Belege; der Fund betrifft die von ihnen
  hinterlassenen, nicht aufgeräumten `id`-Werte in einer von einem
  späteren Test exklusiv benötigten Tabelle)
- geprüft, ohne Befund: `internal/adapters/driven/grpcstream/`,
  `internal/application/port/` (kein Diff, wie von der ADR verlangt)
- geprüft, ohne Befund: `spec/architecture.md`, `.a-check.yml` (kein Diff)
- geprüft, ohne Befund: Traceability — alle vier Commit-Messages nennen
  `ADR-0133`, kein `SPEC-*`/`ARC-*` im Betreff

---

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** ID-Kollision zwischen neuem
E2E-Testdatensatz und bestehendem Go-Test in derselben Tabelle
(Wertebereichs-Prüfung nur rückwärts, nicht vorwärts im selben Lauf)

## Verdikt

**Merge-blockierend: ja.**

F-1 ist unabhängig von jeder inhaltlichen Bewertung der Code-Substanz
bereits hinreichend: `make test-integration` — die von
[ADR-0133](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md)s
eigener Fitness-Function-Zeile verlangte Prüfung — **scheitert real und
reproduzierbar** am committeten Endstand (`ad7db2f6`). Bemerkenswert und
namentlich festzuhalten: Die beiden **neuen** Filter-Rundläufe, die dieser
Diff einführt, laufen dabei selbst **korrekt und erfolgreich** durch
(beide Anker-Selbstprüfungen und beide Beleg-Echos real gedruckt, inklusive
des echten Negativbelegs „Change auf der anderen Tabelle kommt nicht an").
Der Fehlschlag liegt nicht in der neuen Filter-Logik oder ihrer
E2E-Prüfung selbst, sondern in einer Nebenwirkung ihrer Testdaten: Die
beiden Negativbelege schreiben feste IDs (`10`, `11`) in eine gemeinsam
genutzte Tabelle (`feed_e2e_schema`), ohne sie danach zu entfernen, und
kollidieren damit mit einem später im selben Lauf ausgeführten,
unveränderten Bestandstest (`TestE2ESchemaChangeIncompatibleTypeChange`),
der genau diese beiden IDs für seine eigenen `INSERT`s braucht. Der
Implementer-Kommentar an beiden Einfügestellen zeigt, dass ein
Wertebereichs-Check versucht wurde — er prüfte nur die **vorangehenden**
IDs derselben Tabelle (`1`/`2`), nicht die **späteren**, bereits im
Bestand vorhandenen (`10`/`11`).

Die Code-Substanz der neuen Filterfähigkeit selbst (Domänen-Funktion,
Protobuf-Erweiterung, beide Driving-Handler, Unit-Tests mit korrekt
gebundener Eingabeseite, Spec-/Handbuch-Nachzug) ist nach eigenständiger
Nachprüfung strukturell
[ADR-0133](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md)-konform.
`make gates` ist grün (Coverage 80,60 % ≥ 80 %). Der einzige Mangel liegt
in der Testdaten-Isolation des neu hinzugefügten E2E-Abschnitts gegenüber
einem bestehenden, unveränderten Nachbartest in derselben Tabelle.

**Übergabe:** F-1 geht an den Implementer zurück (Fixrunde erforderlich —
z. B. andere, mit dem gesamten Skript kollisionsfreie IDs für die beiden
Filter-Negativbelege wählen, oder die eingefügten Zeilen am Ende der
jeweiligen Sektion wieder entfernen). Da eine Fixrunde ohnehin aussteht,
entfällt der DoD-Checkbox-Nachzug aus
`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde.
