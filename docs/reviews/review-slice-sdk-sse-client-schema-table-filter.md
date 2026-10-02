# Review-Report: slice-sdk-sse-client-schema-table-filter — 2026-10-02

**Review-Art:** Code — gegen Plan, Entscheidungen und Hard Rules (die DoD-Konformität prüft der Verifier).

**Gegenstand:** `git diff e97fd62d HEAD` (Commits `ed5db509` SDK, `d9990c1f` Beispiele, `07b10726` Doku und Plan); 27 Dateien, kein Server-, kein Produktivcode außerhalb von `sdks/` und `examples/` (`git diff --stat` gelesen).

**Skill:** `.harness/skills/reviewer.md` @ HEAD `07b10726`
**Modell:** Sonnet 5.5 (`claude-sonnet-5-5`) · **Datum:** 2026-10-02

**Eingangs-Kontext:**

- Slice-Plan `slice-sdk-sse-client-schema-table-filter` (DoD, §3 Umsetzung und Suchlauf-Feld)
- [`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) Teilfrage 4, [`ADR-0134`](../plan/adr/0134-sdk-public-doc-check-gate-make-gates.md), [`ADR-0090`](../plan/adr/0090-beispiel-clients-volle-matrix.md)
- [`LH-FA-SST-008`](../../spec/lastenheft.md), [`LH-FA-SST-009`](../../spec/lastenheft.md)
- Server-Auswertung selbst gelesen: `internal/adapters/driving/http/sse.go` (`parseStreamChangesFilter`, `streamChangesAllowedParams`) und `internal/domain/model/change.go` (`MatchesFilter`)
- [`AGENTS.md`](../../AGENTS.md) §3.1, §3.2, §3.7, §3.9, §3.12, §3.13, §3.15

---

## Server-Vertrag, gegen den die Clients gelesen wurden

- Die Reihenfolge der Query-Parameter ist dem Server gleich: `r.URL.Query()` entsteht ungeordnet, `values.Get` je Name.
- Prozent-Dekodierung durch `net/url` (`+` wird Leerzeichen, `%2B` bleibt `+`).
- Ein leerer Wert ist „kein Filter“ (`MatchesFilter`: `schema != ""`).
- Ein Parameter außerhalb von `schema`/`table`/`target` endet mit `400`.
- `table` ohne `schema` trifft jedes Schema. Die Docstrings der drei Clients sagen dasselbe.

## Mutationen (selbst gefahren, Kopie per `git archive` im Scratchpad, `sed … > Kopie` + `cp`, `git status --short` im Echtrepo leer)

| # | Mutation (Kopie) | Instanz | Gesehene Farbe |
|---|---|---|---|
| M1 | C# `if (table is not null)` zu `if (false)` | `docker build --target build` mit `dotnet test` | rot: 4 `PgChangeFeedSseClientTargetTests`-Fälle (`Table_…` ×3, `TableAndTarget`) |
| M2 | Python `("table", table)` aus dem Tupel entfernt | `docker build --target build` mit `pytest` | rot: `4 failed, 146 passed` (150 gesamt, deckt die berichtete Zahl) |
| M3 | Kotlin `"schema" to schema` entfernt | `docker build --target build` mit `./gradlew test` | rot: `116 tests completed, 3 failed` (Kotlin-Tests laufen) |
| M4 | Go `stream.go`: `"table": table` aus der Map entfernt | `go test` im Toolchain-Image | rot (`TestStreamURLCarriesSchemaAndTable`, Fall „table eu“) |
| M5 | Go `main.go`: `flag.StringVar(&cfg.table, …)` gelöscht | `go test ./examples/sse-client/` | **grün** |
| M6 | Go `main.go`: `StreamURL(…, cfg.table)` zu `StreamURL(…, "")` | `go test ./examples/sse-client/` | **grün** |

Nicht gefahren, **hergeleitet** aus dem Testschnitt (Parser-Tests und URL-Bau-Tests, keine Tests auf `Program.cs`/`Main.kt`): dieselbe Lücke wie M5/M6 bei den C#- und Kotlin-Beispielen.

Gelaufen: `make sdk-public-doc-check` Exit 0, `make fmt-check` Exit 0, `make kommentar-kennungen DIFF=e97fd62d` Exit 0, `make suchlauf-nachmessen PLAN=…` Exit 0 (10 Zeilen stimmen), `make docs-check` Exit 0, `make test` Exit 0. `make sdk-pack-*`, `make examples-*` nicht gefahren, ersetzt durch M1–M3 (der Test-Stufen-Bau ohne Cache druckt die Testzahl).

---

## Findings

### F-1 — Beispiel-Tests binden Flag → Anfrage nicht an der Verdrahtungsstelle

- `kategorie`: MEDIUM
- `quelle`: Plan-DoD (B) „ihre Tests belegen Flag → Anfrage“; Reviewer-Skill „Zusage ohne Bindung an ihre Eingabeseite“
- `pfad`: `examples/sse-client/main.go:52` und `:106-107`; analog `examples/csharp/sse-client/Program.cs:45`, `examples/kotlin/sse-client/src/main/kotlin/cdcexamples/sse/Main.kt:47` (analog hergeleitet, nicht gemutet)
- `befund`: Das Entfernen des `-table`-Flags (M5) und das Ersetzen von `cfg.table` durch `""` im `StreamURL`-Aufruf (M6) lassen `go test` grün; die Tests prüfen den URL-Bau (`StreamURL`) und, in C#/Kotlin, den Parser getrennt, nicht ihre Verbindung. Dieselbe Lücke besteht beim vorhandenen `-target` unverändert; der Diff erbt sie und trägt sie auf zwei Parameter weiter.
- `verifizierbar`: ja — M5/M6 (Mutation, kein Gate)
- `klasse`: Zusage ohne Bindung an ihre Eingabeseite (Verdrahtung Flag → Anfrage)

### F-2 — Eingabetabelle ohne `+`, Leerzeichen, Unicode, `%`; Draht-Bytes divergieren in einem Zeichen

- `kategorie`: LOW
- `quelle`: `BEO-PGC/drei-sprachen-kopie-divergiert-am-randfall`; Plan-DoD (A) (die feste Tabelle `null`/`""`/`eu`/`a&b=c` ist erfüllt)
- `pfad`: `sdks/*` SSE-Tests (C# `PgChangeFeedSseClientTargetTests.cs`, Kotlin `PgChangeFeedSseClientTargetTest.kt`, Python `test_sse_client.py`), Go-/C#-/Kotlin-Beispieltests
- `befund`: Keine Sprache testet `+`, Leerzeichen, Nicht-ASCII oder `%`. Aus den Kodierfunktionen hergeleitet (nicht gemessen): C# `Uri.EscapeDataString` und die Kotlin-Funktion `percentEncode` (UTF-8, unreserviert `A-Za-z0-9-_.~`) liefern `%20`/`%2B`/`%C3%A9`; Python (`httpx` `params`) und Go (`url.Values.Encode`) liefern für das Leerzeichen `+`. Der Server dekodiert beides zum selben Wert, die Bytes auf dem Draht sind unterschiedlich.
- `verifizierbar`: nein — kein Gate; Messung wäre ein Test je Sprache mit `"a b"`, `"a+b"`, `"ä"`
- `klasse`: Drei-Sprachen-Kopie am Randfall

### F-3 — Plan-DoD (A) „leer → byte-gleich zum Bestand“ und §3 „leer erscheint als `schema=`“ sprechen nebeneinander

- `kategorie`: LOW
- `quelle`: Maintainability (Plan-Konsistenz); [`ADR-0133`](../plan/adr/0133-tabellen-granulare-filterung-grpc-sse.md) Teilfrage 4
- `pfad`: Plan §1 (A) und §2 DoD (A) gegen §3 „Umsetzung (Implementer)“
- `befund`: Die Auslegung („`null` fehlt auf dem Draht, `""` erscheint als `schema=`“) ist die richtige gegen den Bestand: `target=""` erscheint schon am Parent als `?target=` (Parent-Tests Kotlin Zeile 25, Python Zeile 129), und der Server liest `schema=` als „kein Filter“. Die DoD-Worte nennen „leer“ ohne Unterscheidung `null`/`""` und verweisen in §3 nicht auf diese Festlegung; die Beispiele behandeln leer als „nicht gesetzt“, die Packages als „auf dem Draht“ (im Handbuch 1.87 benannt). Akzeptiert als Verhalten, offen als Text.
- `verifizierbar`: nein
- `klasse`: Nachzug widerspricht dem Nachbarn im selben Träger

### F-4 — API-Additivität: binär nicht kompatibel für bereits kompilierte Aufrufer

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `PgChangeFeedSseClient.cs` (`StreamChangesAsync`), `PgChangeFeedSseClient.kt` (`streamChanges`), `sse_client.py` (`stream_changes`)
- `befund`: C# und Kotlin sind quellkompatibel (neue Parameter hinten, Default `null`, bestehende Aufrufe und benannte Aufrufe `target:` bleiben gültig; positionelle Aufrufe `(ct, "eu")` bleiben gültig). Binär ändert sich die Signatur (C# optionale Parameter, Kotlin ohne `@JvmOverloads`, wie bei `target` schon): bereits gegen 0.2.x kompilierte Aufrufer bräuchten Neukompilierung; Python ist rein keyword-additiv. Kein Release in diesem Slice. Zusätzlich: die Parameter-Reihenfolge des SSE-Clients (`ct, target, schema, table`) unterscheidet sich von der des gRPC-Clients (`schema, table, ct, target`); die READMEs sagen „by name“.
- `verifizierbar`: nein
- `klasse`: API-Erweiterung, Binärkompatibilität

### F-5 — Bereich ohne Befund im Doku-Nachzug (Suchlauf selbst nachgefahren)

- `kategorie`: INFO
- `quelle`: [`AGENTS.md`](../../AGENTS.md) §3.13
- `pfad`: ganzer Baum außer `docs/reviews`, `done/`, `.harness/baseline`
- `befund`: Der eigene Suchlauf (`target only|only filter|einzige.*Filter|filtered by (delivery )?target|does not set|nur .* target`, mit `SSE`/`stream`-Zeilenfilter) findet keine stehen gebliebene Aussage „`target` ist der einzige Filter“ des SSE-Clients: die verbleibenden Treffer sind die Historienzeilen 1.85/1.86 des Handbuchs (bewusst), Plan-Prosa und eine Zeile in [`ADR-0139`](../plan/adr/0139-routing-run-regelstand-fail-closed-und-target-ausserhalb-alphabet.md) über die Wert-Prüfung (anderer Gegenstand). Die Zeile 682 des Handbuchs (Tabelle „Ziel lesen“) nennt für SSE nur `target`; sie beschreibt dort die Zielauswahl, keine Filtermenge, und ist wahr. Python-README Zeile 247 nennt für NATS „by delivery target only“ (wahr, anderer Zustellweg).
- `verifizierbar`: ja — `make suchlauf-nachmessen` (Zahlen), Suchraum und Muster sind Lese-Handlung
- `klasse`: —

---

## Negativbefunde

- geprüft, ohne Befund: `sdks/csharp/PgChangeFeed.Client/Sse/`, `sdks/kotlin/…/sse/`, `sdks/python/…/sse_client.py` — Query-Aufbau, Reihenfolge `schema`, `table`, `target` (dem Server gleich, für ihn ohne Belang), `null` fehlt, `""` erscheint wie bei `target`, `a&b=c` in allen drei Sprachen byte-gleich (`a%26b%3Dc`).
- geprüft, ohne Befund: unbekannter Parameter (`400`) und `table` ohne `schema` — die Clients senden nur die drei bekannten Namen, `table` ohne `schema` wird gesendet und vom Server wie dokumentiert gelesen.
- geprüft, ohne Befund: `examples/` Flag-Namen und „leer = nicht gesetzt“ in Go, C#, Kotlin; `examples/README.md`.
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — `Version: 1.87` und neue Historienzeile 1.87 im selben Diff (Skill-HIGH „Handbuch-Versionshistorie“ nicht ausgelöst), `make docs-check` Exit 0; die drei Package-READMEs englisch, ohne interne Kennung, `make sdk-public-doc-check` Exit 0.
- geprüft, ohne Befund: [`AGENTS.md`](../../AGENTS.md) §3.7 — `make kommentar-kennungen DIFF=e97fd62d` Exit 0, keine neue Kette; §3.2 kein `//nolint`/`# noqa`; §3.1 kein Host-Werkzeug im Diff; §3.15 keine Verweigerung im Bericht erkennbar (aus einem Diff nicht ablesbar).
- geprüft, ohne Befund: Server- und Produktivcode — nicht im Diff.
- geprüft, ohne Befund: Plan §3 Suchlauf-Feld — alle zehn Zeilen mit Exit 0 nachgemessen.
- nicht geprüft: Zählwort „drei optionale“ in der Handbuch-Prosa gegen den Realserver (Unit-Ebene, Realserver-Beleg ist ausdrücklich nicht in diesem Slice); `make sdk-pack-*` und `make examples-*` ungekürzt (durch M1–M3 ersetzt).

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 1 |
| LOW | 2 |
| INFO | 2 |

**Finding-Klassen dieses Laufs:** Zusage ohne Bindung an ihre Eingabeseite · Drei-Sprachen-Kopie am Randfall · Nachzug widerspricht dem Nachbarn im selben Träger · API-Erweiterung, Binärkompatibilität

## Verdikt

**Merge-blockierend:** ja — wegen F-1 (MEDIUM). Die DoD-Zeile „Review durchgeführt“ im Plan bleibt offen, bis F-1 behoben oder durch den Architect eingestuft ist. Die Einstufung MEDIUM statt HIGH begründet sich damit, dass URL-Bau und Parser je für sich gebunden sind (M1–M4 rot), die Lücke eine einzelne Verdrahtungszeile je Beispiel ist und das vorhandene `-target` sie unverändert trägt; die HIGH-Liste führt die Klasse „Zusage ohne Bindung“ dennoch, deshalb geht die Frage an den Architect.

**Übergabe:** F-1 bis F-3 gehen an den Implementer; die Finding-Klassen gehen in die Slice-Closure §7. Der Report ersetzt keine Verifikation.
