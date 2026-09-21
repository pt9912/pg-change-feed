# Review-Report: slice-sdk-csharp-sse-client-flaeche — 2026-09-21

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-csharp-sse-client-flaeche.md`),
`ADR-0106` (Accepted) und `AGENTS.md` Hard Rules (Modul 10 §Drei
Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Diff-Range `502283db..c0288550`, ein Commit (`c0288550`,
„feat(sdk): C#-SSE-Client-Fläche für PgChangeFeed.Client"), Slice
`slice-sdk-csharp-sse-client-flaeche`, Welle
`welle-sdk-csharp-vollabdeckung`.

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, seither um mehrere weitere HIGH-Klassen
ergänzt — u. a. `AGENTS.md` §3.13 Träger-Nachzug, „Beleg trägt seinen Satz
nicht").
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-21.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-sdk-csharp-sse-client-flaeche.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan inkl. Plan-Nachzug, §6 Risiken,
  §8 Sub-Area/Modus)
- `ADR-0106` (Accepted) — Festlegung 1 (letzter Absatz, SSE als
  Folge-Package antizipiert), Festlegung 2 (räumliche Trennung von
  `examples/csharp/`)
- `spec/pflichtenheft.md` §2 `SPEC-021` (Endpunkt, Event-Form,
  Nachrichtenschema)
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentar-Disziplin), §3.12
  (Herkunft von Aussagen), §3.13 (Träger-Nachzug)
- `harness/conventions.md` (MR-000 ID-Schema, MR-002 Slice-Kennungen)
- `harness/README.md` §Sensors/§Minimal-Agent-Workflow
- `examples/csharp/sse-client/SseStream.cs` als Draht-Vorbild (Frame-
  Zerlegung)
- `sdks/csharp/PgChangeFeed.Client/Http/PgChangeFeedHttpClient.cs`,
  `Grpc/PgChangeFeedGrpcClient.cs`, `Http/PgChangeFeedException.cs` als
  bestehende, bereits gereviewte Fläche (Auth-Header-Form,
  Exception-Hierarchie)
- `docs/reviews/review-slice-sdk-csharp-grpc-client-flaeche.md` —
  Formvorbild für diesen Report

**Eigenständig durchgeführte Prüfungen (nicht nur Commit-Message/DoD-Text
übernommen):**

- `docker build -f sdks/csharp/Dockerfile sdks/csharp` **ohne** den
  `--build-context proto=proto`-Zusatzkontext real ausgeführt: Bau bricht
  sichtbar exakt an der `COPY --from=proto`-Zeile ab
  („failed to resolve source metadata for docker.io/library/proto:latest:
  pull access denied …", Dockerfile-Zeile 57 im Fehler-Trace benannt) —
  die Plan-Behauptung „auch ohne gRPC-Bezug bricht der Bau ohne den Flag
  ab" ist real bestätigt, nicht nur aus dem Text übernommen.
- `grep -rn "internal/\|cmd/\|gen/" sdks/csharp/ --include="*.cs"
  --include="*.md" --include="*.csproj"` selbst gefahren: ein Treffer,
  `PgChangeFeed.Client.Tests/Grpc/ChangeMessageSchemaTests.cs:10` — ein
  Doku-Kommentar-Zitat des Go-Testvorbilds
  (`internal/adapters/driving/grpc/server_test.go`), kein Import, **und**
  eine Datei außerhalb dieses Diffs (aus `slice-sdk-csharp-grpc-client-
  flaeche`, dort bereits geprüft). Die neuen SSE-Dateien selbst (`Sse/**`,
  `PgChangeFeed.Client.Tests/Sse/**`) tragen keinen einzigen Treffer.
- Alle `using`-Zeilen der neuen Dateien einzeln gelesen
  (`Sse/PgChangeFeedSseClient.cs`, `Sse/SseFrameParser.cs`,
  `Sse/Models/Change.cs`, fünf Testdateien) — ausschließlich BCL-Typen und
  paketinterne `PgChangeFeed.Client.*`-Namespaces, keine Fremdabhängigkeit
  neu eingeführt (deckt sich mit der Plan-Zusage „kein Fremdmodul").
- `SseFrameParser.ReadFrameAsync`/`PgChangeFeedSseClient.StreamChangesAsync`
  Zeile für Zeile gegen `examples/csharp/sse-client/SseStream.cs` gehalten:
  identisches Grenzverhalten (Leerzeile beendet ein Frame, führende
  Leerzeilen werden übersprungen, ein am Quellende unvollständiges Frame
  liefert `null` statt eine Ausnahme) — real durch die vier
  `SseFrameParserTests`-Grenzfälle abgesichert und selbst nachvollzogen
  (nicht nur gelesen): Leerzeile-Trennung zweier Frames, unvollständiges
  letztes Frame, leere Quelle, führende Leerzeilen.
- Server-Gegenstück gelesen (`internal/adapters/driving/http/sse.go`):
  `fmt.Fprintf(w, "event: %s\ndata: %s\n\n", sseEventChange, payload)` mit
  `payload` aus `json.Marshal` (immer ein einzeiliges, kompaktes JSON-
  Objekt) — der Server schreibt nie mehr als eine `data:`-Zeile pro Event
  und nie eine Kommentarzeile (`:`-Präfix). Der Parser überschreibt bei
  mehreren `data:`-Zeilen statt sie zu verketten (siehe F-1 unten) — real
  gegen den Server-Code geprüft, dass dieser Pfad im aktuellen Vertrag nie
  ausgelöst wird.
- Alle zehn `SPEC-021`-Felder aus `spec/pflichtenheft.md` §`SPEC-021`
  einzeln gegen `Sse/Models/Change.cs`s `JsonPropertyName`-Attribute
  gehalten: `change_id`, `transaction_id`, `source_table_id`, `sequence`,
  `operation`, `old_image`, `new_image`, `schema_version`, `schema`,
  `table` — alle zehn vorhanden, keins fehlt, keins zusätzlich; zusätzlich
  durch `ChangeMessageSchemaTests.ChangeHasExactlyTenProperties`
  (Reflection-Zähler) abgesichert — real gelesen, nicht nur die
  Testnamen übernommen.
- Auth-Boundary-Mapping (`BuildException`) Byte-für-Byte gegen
  `Http/PgChangeFeedHttpClient.cs`s gleichnamige private Methode gehalten:
  identischer `switch`-Ausdruck (400/401/403/404/500 → je ein Exception-
  Typ, sonst `PgChangeFeedUnexpectedStatusException`), identisches
  `ExtractErrorMessage` — eine vollständige, wörtliche Kopie, keine
  abweichende Semantik (siehe F-2 unten zur Wartbarkeits-Bewertung dieser
  Kopie).
- `PgChangeFeedMalformedResponseException` ist keine neue, für diesen
  Slice erfundene Ausnahmeklasse: real in `Http/PgChangeFeedException.cs`
  als bereits bestehender Typ gefunden, den `Http/PgChangeFeedHttpClient.cs`
  bereits für eine 2xx-Antwort mit nicht parsebarem Body wirft — dieselbe
  Semantik wird hier für eine `data:`-Zeile wiederverwendet, keine
  unbegründete Design-Abweichung.
- `git diff 502283db..c0288550 -- sdks/csharp/README.md` gelesen: die
  §Status-Zeile „SSE and NATS-vollinhalts delivery remain uncovered by
  this package" ist ersetzt durch einen Satz, der die SSE-Fläche jetzt
  als Teil des Release-Standes nennt und ausschließlich NATS als offen
  belässt. Zusätzlich `git show 502283db:sdks/csharp/README.md | grep -n
  SSE` gegen den Elternstand gefahren: die einzige zweite `SSE`-Nennung
  (Zeile 21, „HTTP/gRPC/SSE/NATS delivery paths … see the project README")
  bestand bereits im Elternstand unverändert und behauptet keine
  Uncovered-Aussage über das Package — kein Rest der veralteten Aussage
  irgendwo im Dokument.
- `git diff 502283db..c0288550 --stat -- spec/ docs/user/ examples/
  .a-check.yml` — leer: keiner der vier Bereiche berührt, wie im Plan
  §1/§2 vorab benannt (Doku-Nachzug bewusst dem Folge-Slice vorbehalten).
- `grep -n Version sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj`
  — `<Version>0.1.0</Version>`, real unverändert gegenüber dem Elternstand.
- Backtick-Parität real nachgezählt (Backtick-Zeichen je Datei gezählt) für
  beide in diesem Diff geänderten Markdown-Dateien:
  `slice-sdk-csharp-sse-client-flaeche.md` → 252 (gerade),
  `sdks/csharp/README.md` → 48 (gerade) — beide paarig.
- Neue Dateien nach `slice-`/`welle-`-Nennungen und nach Konjunktiv-/
  Vorher-Nachher-Sprache durchsucht (`AGENTS.md` §3.7) — kein Treffer in
  `Sse/**`/`PgChangeFeed.Client.Tests/Sse/**`.
- `make gates` ungefiltert laufen lassen, Exit-Code direkt (kein Pipe/
  Wrapper dazwischen) geprüft: `0` — u. a. `generated-sync: OK`,
  `a-check: gesamt: 0 Befund(e)`.
- Testeigenständigkeit der Mutationsklassen plausibilisiert (nicht
  automatisiert mutiert, aber Bindung an die Eingabeseite nachvollzogen):
  Die Frame-Grenzfall-Tests rufen `SseFrameParser.ReadFrameAsync` direkt
  mit einem echten `StringReader` auf — kein Fake, keine Ausgabeseiten-
  Mutation möglich, jede Eingabe-Variante (fehlende Leerzeile, leere
  Quelle, führende Leerzeilen) ist eine reale Eingabe-Mutation. Die
  Auth-Boundary-Tests steuern über den Responder-Lambda den realen
  `HttpStatusCode` der gefakten Antwort; dieser Wert fließt unverändert
  als `(int)response.StatusCode` in den echten `switch`-Ausdruck von
  `BuildException` — ein geänderter Statuscode im Test änderte den
  Kontrollfluss der echten Mapping-Logik, keine Ausgabeseiten-Mutation
  eines Stubs.

---

## Findings

### F-1 — `SseFrameParser` überschreibt statt verkettet mehrere `data:`-Zeilen

- `kategorie`: LOW
- `quelle`: Maintainability (SSE-Spezifikation, RFC-Konvention „mehrere
  aufeinanderfolgende `data:`-Zeilen werden mit `\n` verkettet")
- `pfad`: `sdks/csharp/PgChangeFeed.Client/Sse/SseFrameParser.cs:44-46`
- `befund`: `data = line[DataPrefix.Length..];` überschreibt den bisher
  gesammelten Wert bei jeder weiteren `data:`-Zeile, statt die Zeilen zu
  verketten — bei mehreren `data:`-Zeilen in einem Frame ginge der
  vorherige Inhalt verloren. Der aktuelle Server (`internal/adapters/
  driving/http/sse.go`) schreibt nie mehr als eine `data:`-Zeile pro
  Event (ein kompaktes `json.Marshal`-Ergebnis), und das gelesene
  Draht-Vorbild (`examples/csharp/sse-client/SseStream.cs`) trägt exakt
  dasselbe Verhalten — kein neu eingeführter Fehler, sondern ein
  unverändert übernommenes, bislang folgenloses Verhalten.
- `verifizierbar`: ja — ein Test, der zwei `data:`-Zeilen in einem Frame
  sendet, würde den Verkettungs-Fall rot zeigen.
- `klasse`: „Vorbild-Verhalten unverändert übernommen, Spec-Lücke bleibt
  latent"

### F-2 — `BuildException`/`ExtractErrorMessage` wörtlich aus `PgChangeFeedHttpClient.cs` kopiert

- `kategorie`: LOW
- `quelle`: Maintainability (DRY)
- `pfad`: `sdks/csharp/PgChangeFeed.Client/Sse/PgChangeFeedSseClient.cs:127-153`
  vs. `sdks/csharp/PgChangeFeed.Client/Http/PgChangeFeedHttpClient.cs:224-249`
- `befund`: Beide private Methoden sind byte-identisch — eine künftige
  Änderung am `SPEC-018`-Fehler-Mapping (z. B. ein neuer Statuscode-Zweig)
  müsste an beiden Stellen nachgezogen werden, ohne dass ein Compiler-
  oder Testfehler das an der zweiten Stelle erzwingt. Der Plan begründet
  den Verzicht auf eine gemeinsame Extraktion explizit mit der bereits
  bestehenden Unabhängigkeit von HTTP- und gRPC-Fläche — dort ist die
  Unabhängigkeit aber strukturell (gRPC bildet `Grpc.Core.StatusCode`
  nativ ab, kein Analogon zu diesem `int`-`switch`), hier ist es eine
  wörtliche Kopie derselben Logik.
- `verifizierbar`: nein — kein Gate prüft Code-Duplikation in diesem Repo.
- `klasse`: „Wörtliche Kopie statt gemeinsamer Extraktion bei identischer
  Logik"

### F-3 — Kein Test für eine SSE-Kommentarzeile (`:`-Präfix)

- `kategorie`: LOW
- `quelle`: Maintainability (Testabdeckung Randfall)
- `pfad`: `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Sse/SseFrameParserTests.cs`
- `befund`: `SseFrameParser.ReadFrameAsync` ignoriert jede Zeile, die
  weder mit `event: ` noch mit `data: ` beginnt (also auch eine
  SSE-Kommentarzeile wie `: keep-alive`), korrekt per Konstruktion —
  aber kein Testfall bildet diesen Fall explizit ab. Der reale Server
  schreibt aktuell keine Kommentarzeilen (`internal/adapters/driving/
  http/sse.go` kennt kein Keep-Alive-Comment), das Verhalten ist also
  unbelegt statt falsch.
- `verifizierbar`: ja — ein zusätzlicher Testfall mit einer
  `:`-Präfix-Zeile würde die Erwartung explizit machen.
- `klasse`: „Korrektes Randfall-Verhalten ohne expliziten Beleg"

### F-4 — Frame-Name (`event:`-Wert) wird nie gegen `"change"` validiert

- `kategorie`: INFO
- `quelle`: Maintainability
- `pfad`: `sdks/csharp/PgChangeFeed.Client/Sse/PgChangeFeedSseClient.cs:88-99`
- `befund`: `StreamChangesAsync` verarbeitet jedes vollständige Frame als
  `Change`, unabhängig vom `Name`-Feld des `SseFrame` — `SPEC-021` kennt
  aktuell nur `event: change`, ein künftiger zweiter Event-Typ würde ohne
  Anpassung hier ebenfalls als `Change` interpretiert (oder an der
  JSON-Deserialisierung scheitern). Kein Verweis auf eine zuständige
  Rolle nötig — reiner Hinweis ohne erwartete Aktion, solange `SPEC-021`
  nur einen Event-Typ kennt.
- `verifizierbar`: nein — keine Gegenprobe ohne eine hypothetische
  zweite Event-Art.
- `klasse`: „Discriminator gelesen, aber nicht geprüft"

## Negativbefunde

- geprüft, ohne Befund: `SPEC-021`-Vollständigkeit — alle zehn
  Nachrichtenfelder Feld für Feld gegen `spec/pflichtenheft.md` §`SPEC-021`
  gehalten, `Change.cs` deckt sie exakt, zusätzlich durch einen
  Feldanzahl-Drift-Wächter (`ChangeHasExactlyTenProperties`) abgesichert.
- geprüft, ohne Befund: SSE-Frame-Parsing — Leerzeile-Grenze, führende
  Leerzeilen, unvollständiges letztes Frame am Stream-Ende, Kommentarzeilen
  (korrekt durch Konstruktion ignoriert, siehe F-3 zur Testabdeckung)
  — alle Fälle real nachvollzogen; die einzige echte Lücke ist die unter
  F-1 benannte Nicht-Verkettung mehrerer `data:`-Zeilen, aktuell folgenlos.
- geprüft, ohne Befund: Auth-Boundary — Bearer-Token wird identisch zur
  bestehenden HTTP-Fläche im `Authorization`-Header gesetzt
  (`Bearer <token>`); 401/403/500/503/unbekannt werden über die
  bestehende, geschlossene `PgChangeFeedException`-Hierarchie typisiert,
  keine neue, zweite Hierarchie.
- geprüft, ohne Befund: Design-Entscheidung „kein Exception bei
  regulärem Stream-Ende, `PgChangeFeedMalformedResponseException` bei
  Nicht-JSON-Payload" — beide Bausteine sind real bereits bestehende
  Muster (gRPC-Flächen-Fire-and-Forget-Semantik bzw. der bereits in der
  HTTP-Fläche existierende Exception-Typ für einen nicht parsebaren
  2xx-Body), keine unbegründete Neuerfindung.
- geprüft, ohne Befund: `--build-context proto=proto`-Zwang — real per
  eigenem `docker build` ohne das Flag reproduziert: Bau bricht exakt an
  der `COPY --from=proto`-Zeile ab, die Implementer-Behauptung ist
  bestätigt, nicht nur übernommen.
- geprüft, ohne Befund: Import-Grenze (`internal/**`/`cmd/**`/`gen/**`)
  — eigener `grep`-Lauf über `sdks/csharp/`, ein Treffer außerhalb dieses
  Diffs (`Grpc/ChangeMessageSchemaTests.cs:10`, Doku-Zitat, bereits im
  Vorgänger-Review geprüft); die neuen SSE-Dateien selbst tragen keinen
  Treffer.
- geprüft, ohne Befund: Plan-Nachzug (§3, drei benannte Abweichungen) —
  `Sse/Models/Change.cs` folgt demselben Trennungsmuster wie
  `Http/Models/Changes.cs`; fünf statt eine Testdatei folgt demselben
  Gruppierungsmuster wie der gRPC-Slice; keine Mapper-Extraktion aus
  `Http/` ist als bewusste Design-Entscheidung begründet (siehe F-2 zur
  eigenständigen Wartbarkeits-Bewertung dieser Entscheidung) — alle drei
  sind unkritische, begründete Ausgestaltung, kein Plan-Konflikt.
- geprüft, ohne Befund: Mutation-Testing-Plausibilisierung — Frame-
  Grenzlogik und Auth-Mapping sind an ihrer jeweiligen Eingabeseite
  gebunden (echter `StringReader`, echter `HttpStatusCode` durch den
  Responder), keine Bindung ausschließlich an eine Ausgabeseite/einen
  Fake-Rückgabewert.
- geprüft, ohne Befund: `sdks/csharp/README.md`-Korrektur — die
  stehengebliebene „SSE and NATS-vollinhalts delivery remain uncovered"-
  Aussage ist vollständig und korrekt behoben; kein Rest der veralteten
  Aussage an anderer Stelle des Dokuments.
- geprüft, ohne Befund: Version-Bump — `PgChangeFeed.Client.csproj`s
  `<Version>` real unverändert `0.1.0`.
- geprüft, ohne Befund: `spec/pflichtenheft.md`/`docs/user/
  benutzerhandbuch.md` — beide real unberührt in diesem Diff, korrekt dem
  Folge-Slice überlassen (Plan §1/§2 benennt das vorab).
- geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7) — kein
  Konjunktiv über eine verworfene Alternative, kein abwesender Text,
  keine narrative Slice-/Wellen-Chronik in den neuen `.cs`-Dateien.
- geprüft, ohne Befund: Backtick-Parität — beide in diesem Diff
  geänderten Markdown-Dateien real nachgezählt, beide paarig (252, 48).
- geprüft, ohne Befund: Out-of-Scope — kein Bezug zu
  `sdks/csharp/**Nats**`, keine Änderung an einem Pack-Werkzeug- oder
  Publish-Workflow-Skript, kein Bezug zu `spec/architecture.md` oder
  `.a-check.yml` (`git diff --stat` gegen den vollen Diff bestätigt: nur
  `Sse/**`, die fünf Testdateien und `README.md` sind berührt).
- geprüft, ohne Befund: Traceability — Commit-Betreff nennt
  `LH-FA-SST-009`/`ADR-0106`, kein `SPEC-*`/`ARC-*` im Betreff.
- geprüft, ohne Befund: `make gates` real gefahren, Exit-Code direkt
  (ungepiped) geprüft, `0`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 3 |
| INFO | 1 |

**Finding-Klassen dieses Laufs:** „Vorbild-Verhalten unverändert
übernommen, Spec-Lücke bleibt latent" · „Wörtliche Kopie statt
gemeinsamer Extraktion bei identischer Logik" · „Korrektes
Randfall-Verhalten ohne expliziten Beleg" · „Discriminator gelesen, aber
nicht geprüft"

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, drei LOW und ein INFO,
keines davon mit Sicherheits-, Korrektheits- oder ADR-Verstoß-Charakter.
Keine Fixrunde am Implementer nötig; die vier Findings sind Hinweise für
eine künftige Gelegenheit (z. B. den Folge-Slice, der ohnehin
`spec/pflichtenheft.md`/das Handbuch anfasst), keine Blocker dieses
Slice.

**DoD-Checkbox-Nachzug ohne Fixrunde** (Skill-Regel,
`.harness/skills/reviewer.md` §DoD-Checkbox-Nachzug ohne Fixrunde): Da
dieses Verdikt zu keiner Fixrunde führt, wird die DoD-Zeile „Review
durchgeführt, Report unter `docs/reviews/` liegt vor" im Slice-Plan
(`docs/plan/planning/in-progress/slice-sdk-csharp-sse-client-flaeche.md`)
im selben Commit, der diesen Report anlegt, auf `[x]` nachgezogen, mit
Verweis auf diesen Report-Pfad.

**Übergabe:** kein Rückgabe-Pfeil an den Implementer nötig. Die vier
Findings (drei LOW, ein INFO) gehen in die Slice-Closure §7 und von dort
in den Steering-Loop-Zähler. Dieser Report ist ein Lauf-Beleg; er ersetzt
keine Verifikation gegen die volle DoD — das bleibt Verifier-Aufgabe
(Modul 11).
