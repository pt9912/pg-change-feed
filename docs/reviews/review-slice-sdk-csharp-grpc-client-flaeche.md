# Review-Report: slice-sdk-csharp-grpc-client-flaeche — 2026-09-19

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-csharp-grpc-client-flaeche.md`),
`ADR-0106` (Accepted) und `AGENTS.md` Hard Rules (Modul 10 §Drei
Review-Arten). Kein DoD-Abgleich — das ist Verifier-Aufgabe (Modul 11).

**Gegenstand:** Diff-Range `87fa2f82..HEAD` (Abschluss von
`slice-sdk-csharp-http-client-flaeche` bis Implementer-Commit dieses
Slice), Slice `slice-sdk-csharp-grpc-client-flaeche`, Welle
`welle-sdk-csharp-lh-fa-sst-009`. Vier Commits: `09d8165f` (open→next,
reiner Move, 0 Insertions/Deletions), `6eae87b1` (Verantwortlich gesetzt,
nur Slice-Datei), `fc3426f1` (next→in-progress, reiner Move, 0
Insertions/Deletions), `06e78545` (Inhalt: `PgChangeFeedGrpcClient` +
Tests + `Directory.Packages.props`/`Dockerfile`/`.csproj`-Ergänzung +
Handbuch-Update).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, seither um mehrere weitere
HIGH-Klassen ergänzt — u. a. `AGENTS.md` §3.13 Träger-Nachzug, „Arbeit
überholt stehenden Träger").
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-19.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-sdk-csharp-grpc-client-flaeche.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan inkl. Plan-Nachzug, §6 Risiken,
  §8 Sub-Area/Modus)
- `ADR-0106` (Accepted) — Festlegung 1/2/3/4/5, §Konsequenzen
  Folgepflicht 2/4
- `spec/pflichtenheft.md` §2 `SPEC-020` (Nachrichtenschema, RPC-Name,
  Stream-Semantik)
- `AGENTS.md` §3.1 (Docker-only), §3.3 (git mv + Inhalt = zwei Commits),
  §3.7 (Kommentar-Disziplin), §3.12 (Herkunft von Aussagen), §3.13
  (Träger-Nachzug)
- `harness/conventions.md` (MR-000 ID-Schema, MR-002 Slice-Kennungen)
- `docs/reviews/review-slice-sdk-csharp-http-client-flaeche.md` —
  Vorgeschichte: F-1 (HIGH, veralteter README-Statusabsatz) und F-2
  (MEDIUM, nicht-typisierter Fehlerpfad), beide dort in einer Fixrunde
  behoben — gezielt geprüft, ob dieselben zwei Klassen hier vermieden
  wurden
- `examples/csharp/grpc-client/Program.cs` als Draht-Vorbild (Bearer-Token-
  Metadata-Form, byte-für-byte verglichen)
- `docs/reviews/review-slice-sdk-csharp-http-client-flaeche.md` —
  Formvorbild für diesen Report

**Eigenständig durchgeführte Prüfungen (nicht nur Commit-Message/DoD-Text
übernommen):**

- `docker build --no-cache --build-context proto=proto -f
  sdks/csharp/Dockerfile sdks/csharp` real ausgeführt: Exit-Code direkt
  geprüft, `0`. `dotnet build` 0 Warnings/0 Errors, `dotnet test`:
  **32/32 Tests grün** — selbst nachgemessen, nicht übernommen.
  Test-Image danach mit `docker rmi` entfernt.
- `docker build --no-cache -f sdks/csharp/Dockerfile sdks/csharp` **ohne**
  den `proto`-Bau-Kontext real ausgeführt: Exit-Code direkt (in eine
  Log-Datei umgeleitet, danach separat gelesen) geprüft, `1`. Der Bau
  bricht sichtbar exakt an der `COPY --from=proto …`-Zeile ab
  („failed to resolve source metadata for docker.io/library/proto:latest:
  pull access denied …", Dockerfile-Zeile 57 im Fehler-Trace benannt) —
  kein stiller Fallback, kein übersprungener Schritt. Beleg für die
  DoD-Zusage „ohne ihn bricht der Bau an der COPY-Zeile ab" real erbracht,
  nicht nur den Code gelesen.
- `make docs-check` real ausgeführt, Exit-Code direkt geprüft: `0`
  (`d-check: 803 Datei(en) geprüft, 0 Befund(e)`).
- Bearer-Token-Metadata-Form Byte-für-Byte gegen
  `examples/csharp/grpc-client/Program.cs` verglichen:
  `AuthorizationMetadataKey = "authorization"`, `BearerPrefix = "Bearer "`,
  `new Metadata { { AuthorizationMetadataKey, BearerPrefix + token } }` —
  identisch in beiden Dateien, keine Abweichung.
- `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Grpc/ChangeMessageSchemaTests.cs`
  Feld für Feld gegen die zehn `SPEC-020`-Felder gezählt: `change_id`,
  `transaction_id`, `source_table_id`, `sequence`, `operation`,
  `old_image`, `new_image`, `schema_version`, `schema`, `table` — alle
  zehn vorhanden, keins fehlt, keins zusätzlich; zusätzlich abgesichert
  durch einen eigenen `Change.Descriptor.Fields.Count == 10`-Drift-Wächter.
- `grep -rniE "follow-up|added by follow|not yet|coming soon|to be
  added|will be added|future release|next release|remains? to be"
  sdks/csharp/` — zwei Treffer, beide **außerhalb** des Diffs
  (`PgChangeFeedClientOptions.cs:12`,
  `PgChangeFeedClientOptionsTests.cs:7`, bereits in der Fixrunden-
  Nachprüfung des Vorgänger-Reviews als Scope-Aussage über die
  `PgChangeFeedClientOptions`-Klasse selbst geprüft und nicht als
  Lieferstand-Behauptung eingestuft) — kein Treffer, der die
  gRPC-Fläche als „noch nicht geliefert" beschreibt.
  `sdks/csharp/README.md:9` selbst gelesen: nennt jetzt explizit „and the
  gRPC live-change-stream surface (`PgChangeFeedGrpcClient`,
  `StreamChangesAsync`)" als Teil des aktuellen Release-Standes — dieselbe
  Fehlerklasse wie F-1 im Vorgänger-Review ist hier vermieden.
- `grep -rn "internal/\|cmd/pg-change-feed\|gen/cdc" sdks/csharp/` — ein
  Treffer, `ChangeMessageSchemaTests.cs:10`, ein Doku-Kommentar-Zitat des
  Test-Vorbilds (`internal/adapters/driving/grpc/server_test.go`), kein
  Import. Zusätzlich geprüft: der generierte `Change`-Stub entsteht laut
  `.csproj` (`<Protobuf Include="Grpc/proto/changestream.proto" …>`) und
  Dockerfile-`COPY --from=proto`-Zeile im Bau aus der über den
  Zusatzkontext gelesenen `.proto` — kein Kopiervorgang aus `gen/cdc/**`.
- `git diff 87fa2f82..HEAD --stat -- examples/` — leer:
  `examples/csharp/grpc-client/**` real unverändert.
- `git show 09d8165f --stat`/`git show fc3426f1 --stat` geprüft: beide
  zeigen ausschließlich den Rename, 0 Insertions/Deletions — reine Moves,
  `AGENTS.md` §3.3 sauber eingehalten.
- `docs/user/benutzerhandbuch.md`-Diff gelesen: `Version: 1.33` → `1.34`
  korrekt hochgezogen, neue Zeile `| 1.34 | 2026-09-19 | … |` in
  `### Änderungshistorie` vorhanden; der neue `**SDK:**`-Absatz unter §4
  „Zugriff über den gRPC-Change-Stream" nennt „allen zehn Feldern der
  Tabelle oben" — gegen die unmittelbar vorangehende Feldtabelle (Zeile
  692–705) gehalten: exakt dieselben zehn Felder, kein Drift.
  Rest des Diffs (`sdks/csharp/**`) enthält weder `harness/mk/sdk.mk`
  noch `harness/mk/examples.mk`-Änderungen — deckt sich mit dem Plan-Nachzug
  „bleibt in diesem Slice unberührt".
- `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/PgChangeFeed.Client.Tests.csproj`
  gelesen — unverändert in diesem Diff; kein zusätzlicher `PackageReference`
  für `Grpc.Core` nötig, kommt transitiv über `ProjectReference` auf
  `PgChangeFeed.Client.csproj` (bestätigt durch den grünen Testlauf).
- Neue Dateien nach `slice-`/`welle-`-Nennungen durchsucht (`AGENTS.md`
  §3.7 Chronik-Klasse): Treffer ausschließlich in `Dockerfile`,
  `.csproj`, `Directory.Packages.props` — durchweg in der etablierten
  Herkunfts-Anker-Form „(slice-<name>, `ADR-0106` Festlegung …)", dieselbe
  Form, die bereits in `slice-sdk-csharp-projektgeruest`s Dockerfile
  unbeanstandet blieb; kein narrativer Vorher/Nachher-Satz, keine Chronik
  in `.cs`-Produktionscode-Dateien.

---

## Findings

Keine HIGH-, MEDIUM- oder LOW-Findings in diesem Lauf.

## Negativbefunde

- geprüft, ohne Befund: Bearer-Token-Metadata-Form byte-für-byte gegen
  `examples/csharp/grpc-client/Program.cs` — identisch (Schlüsselname,
  Präfix, Metadata-Konstruktion).
- geprüft, ohne Befund: Nachrichtenschema-Vollständigkeit — alle zehn
  `SPEC-020`-Felder in `ChangeMessageSchemaTests.cs` einzeln gezählt,
  keins fehlt, zusätzlich durch einen Feldanzahl-Drift-Wächter
  abgesichert.
- geprüft, ohne Befund: Authn-Boundary-Test (`FakeCallInvoker.WithStatus`,
  `StatusCode.Unauthenticated` surfaced als `RpcException` aus der
  Aufzählung selbst, nicht als stiller leerer Stream) — real gegen den
  Fake ausgeführt, grün.
- geprüft, ohne Befund: `sdks/csharp/Dockerfile` — der `proto`-Bau-Kontext
  ist real zwingend; ein Bau ohne `--build-context proto=proto` bricht
  real und sichtbar an der `COPY --from=proto`-Zeile ab (kein `ARG`-
  Fallback, kein `|| true`, kein stiller Leerlauf).
- geprüft, ohne Befund: `sdks/csharp/README.md` §Status — nennt die
  gRPC-Fläche jetzt korrekt als Teil des aktuellen Release-Standes, keine
  „follow-up release"-Behauptung mehr für gRPC (die im Vorgänger-Review
  gefundene Fehlerklasse F-1 ist hier vermieden); eigener Grep nach
  mehreren Formulierungsvarianten (`follow-up`, `not yet`, `coming soon`,
  `to be added`, `future release`, …) bestätigt das über den gesamten
  `sdks/csharp/`-Baum, die zwei verbleibenden Treffer liegen außerhalb
  dieses Diffs und behaupten keine veraltete Tatsache über die gRPC-Fläche.
- geprüft, ohne Befund: kein nicht-typisierter Fehlerpfad analog F-2 des
  Vorgänger-Reviews — der Authn-Boundary-Pfad wirft die bereits von
  `Grpc.Core` typisierte `RpcException`/`StatusCode.Unauthenticated`
  durch, kein zusätzlicher, selbst geschriebener Deserialisierungs-/
  Mapping-Schritt, an dem eine analoge Lücke entstehen könnte.
  `Console`-freier, reiner `IAsyncEnumerable`-Yield-Pfad ohne
  Zwischenschicht.
- geprüft, ohne Befund: Import-Grenze (`grep -rn
  "internal/\|cmd/pg-change-feed\|gen/cdc" sdks/csharp/` — ein Treffer,
  ein Doku-Kommentar-Zitat, kein Import); generierter Stub entsteht im
  Bau, kein committeter Stub, keine Kopie aus `gen/cdc/**`.
- geprüft, ohne Befund: `examples/csharp/grpc-client/**` unverändert in
  diesem Diff (`git diff 87fa2f82..HEAD --stat -- examples/` leer).
- geprüft, ohne Befund: `harness/mk/sdk.mk`/`harness/mk/examples.mk` —
  beide unberührt, wie im Plan-Nachzug §3 vorab benannt.
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — Version
  1.33→1.34 korrekt, Änderungshistorie-Zeile vorhanden und inhaltlich
  zutreffend (Feldzahl „zehn" stimmt gegen die unmittelbar vorangehende
  Tabelle).
- geprüft, ohne Befund: `git mv` + Inhaltsänderung als getrennte Commits
  (`AGENTS.md` §3.3) — beide Lifecycle-Moves (`09d8165f`, `fc3426f1`)
  tragen 0 Insertions/Deletions.
- geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7) — kein
  Konjunktiv über eine verworfene Alternative, kein abwesender Text,
  keine narrative Slice-/Wellen-Chronik in `.cs`-Produktionscode; die
  Slice-Nennungen in `Dockerfile`/`.csproj`/`Directory.Packages.props`
  folgen der etablierten Herkunfts-Anker-Form (Name + ADR-Festlegung),
  keine Vorher/Nachher-Erzählung.
- geprüft, ohne Befund: Traceability — Commit-Betreff `06e78545` nennt
  `LH-FA-SST-009` und `ADR-0106`, kein `SPEC-*`/`ARC-*` im Betreff.
- geprüft, ohne Befund: reale Docker-Build- und Testausführung (mit
  Zusatzkontext) — Exit-Code direkt geprüft, `0`; 32/32 Tests grün, Build
  ohne Warnungen (die einzige Docker-Meldung ist der plattform-bedingte
  `InvalidBaseImagePlatform`-Hinweis, kein neues Muster). Test-Image nach
  Prüfung entfernt.
- geprüft, ohne Befund: realer Bau-Abbruch ohne Zusatzkontext — Exit-Code
  `1`, Fehler exakt an der `COPY --from=proto`-Zeile verortet.
- geprüft, ohne Befund: `make docs-check` real gefahren, Exit-Code direkt
  geprüft, `0` (803 Dateien, 0 Befunde).
- geprüft, ohne Befund: `Google.Protobuf`-Versionsabweichung
  (`3.36.2` statt `examples/csharp`s `3.36.1`) — Plan-Nachzug benennt sie
  korrekt als reale Neu-Messung (`AGENTS.md` §3.12), kein unbegründeter
  Drift; `Grpc.Net.Client`/`Grpc.Tools` deckungsgleich.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 0 |
| MEDIUM | 0 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** keine — beide in der Vorgeschichte
benannten Fehlerklassen (veralteter README-Statusabsatz, nicht-typisierter
Fehlerpfad) sind in diesem Slice real vermieden, nicht nur behauptet.

## Verdikt

**Merge-blockierend:** nein — 0 HIGH, 0 MEDIUM, 0 LOW. Keine Fixrunde am
Implementer nötig.

**DoD-Checkbox-Nachzug ohne Fixrunde** (Skill-Regel, `.harness/skills/reviewer.md`
§DoD-Checkbox-Nachzug ohne Fixrunde): Da dieses Verdikt zu keiner Fixrunde
führt, wird die DoD-Zeile „Review durchgeführt, Report unter
`docs/reviews/` liegt vor" im Slice-Plan
(`docs/plan/planning/in-progress/slice-sdk-csharp-grpc-client-flaeche.md`)
im selben Commit, der diesen Report anlegt, auf `[x]` nachgezogen, mit
Verweis auf diesen Report-Pfad.

**Übergabe:** kein Rückgabe-Pfeil an den Implementer nötig. Dieser Report
ist ein Lauf-Beleg; er ersetzt keine Verifikation gegen die volle DoD —
das bleibt Verifier-Aufgabe (Modul 11).
