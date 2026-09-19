# Verifikationsbericht: slice-sdk-csharp-grpc-client-flaeche — 2026-09-19

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den
DoD-Vertrag
(`docs/plan/planning/in-progress/slice-sdk-csharp-grpc-client-flaeche.md`
§2) und die dort referenzierte ADR/Spec, in frischem Kontext. **Nicht**
gegen den Diff als solchen (Reviewer-Aufgabe, abgeschlossen mit
[`review-slice-sdk-csharp-grpc-client-flaeche.md`](review-slice-sdk-csharp-grpc-client-flaeche.md))
und **nicht** gegen realen Bedarf (Validator, hier nicht ausgelöst).

**Gegenstand:** `slice-sdk-csharp-grpc-client-flaeche`, Welle
`welle-sdk-csharp-lh-fa-sst-009`. Fünf Commits seit `87fa2f82`
(Abschluss von `slice-sdk-csharp-http-client-flaeche`): `09d8165f`
(open→next), `6eae87b1` (Verantwortlich gesetzt), `fc3426f1`
(next→in-progress), `06e78545` (Implementer-Commit: `PgChangeFeedGrpcClient`
+ Tests + `Directory.Packages.props`/`Dockerfile`/`.csproj`-Ergänzung +
Handbuch-Update), `d86f46b7` (Review-Report, 0 Findings), `899a3f80`
(rein mechanischer Docs-Check-Nachbessern des Review-Reports selbst, kein
Verdikt-Wechsel).

**Frischer Kontext:** Diese Sitzung hat Slice-Plan (vollständig, §1–§8),
Review-Report (vollständig), `ADR-0106` (vollständig, Festlegung 1–5,
§Kontext, §Konsequenzen) und `spec/pflichtenheft.md` `SPEC-020` selbst
gelesen, sowie jede in §2 des Plans genannte Quelldatei selbst geöffnet:
`PgChangeFeedGrpcClient.cs`, `PgChangeFeed.Client.csproj`,
`Directory.Packages.props`, `Dockerfile`, `ChangeMessageSchemaTests.cs`,
`PgChangeFeedGrpcClientTests.cs`, `benutzerhandbuch.md`-Diff,
`sdks/csharp/README.md`. Nichts aus DoD-Text, Commit-Message oder
Review-Report ungeprüft übernommen — eigener `make gates`-Lauf (ungepiped,
Exit-Code direkt geprüft), zwei eigene, unabhängige Docker-Build-Läufe (mit
und ohne `--build-context proto=proto`), eigene `git ls-files`/`grep`-Läufe
gegen die Import- und Proto-Stub-Grenze, eigener `git diff --stat` gegen
die in `ADR-0106` §5 als unberührt deklarierten Pfade.

---

## 1. DoD-Vertrag (§2) — jede Checkbox einzeln geprüft

### 1.1 Öffentliche Client-Klasse, `StreamChanges`-Methode, Bearer-Token im `authorization`-Metadata-Eintrag

Eigene Lektüre von
`sdks/csharp/PgChangeFeed.Client/Grpc/PgChangeFeedGrpcClient.cs`: die
öffentliche Klasse `PgChangeFeedGrpcClient` (`sealed`, `IDisposable`)
trägt `StreamChangesAsync(CancellationToken)` als
`IAsyncEnumerable<Change>`. Konstruktion (`Metadata { { "authorization",
"Bearer " + _options.ApiToken } }`, Zeile 97) legt den Bearer-Token exakt
in den `authorization`-Metadata-Eintrag — Wortform identisch mit
`SPEC-020`s Festlegung. Alle zehn `SPEC-020`-Felder (`change_id`,
`transaction_id`, `source_table_id`, `sequence`, `operation`,
`old_image`, `new_image`, `schema_version`, `schema`, `table`) werden
unverändert als generierte `Change`-Nachricht durchgereicht (kein
zusätzlicher DTO-Mapping-Schritt, der ein Feld verlieren könnte).
**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.2 `.csproj`/`Directory.Packages.props` — drei gRPC-Pakete zentral gepinnt, exakte Versionen

Eigene Lektüre von `PgChangeFeed.Client.csproj`: drei
`PackageReference`-Einträge ohne `Version`-Attribut (`Grpc.Net.Client`,
`Grpc.Tools` mit `PrivateAssets="All"`, `Google.Protobuf`) — die
Versionshoheit liegt zentral in `Directory.Packages.props`
(`ManagePackageVersionsCentrally=true`). Dort: `Version="2.83.0"`,
`Version="2.84.0"`, `Version="3.36.2"` — jeweils exakte Versionen, kein
`>=`, `*` oder Bereichs-Syntax. **Ergebnis: Checkbox berechtigt auf
`[x]`.**

### 1.3 `Dockerfile` — zusätzlicher `proto`-Bau-Kontext, kein committeter Stub

Eigene Lektüre von `sdks/csharp/Dockerfile` Zeile 57: `COPY --from=proto
cdc/stream/v1/changestream.proto
PgChangeFeed.Client/Grpc/proto/changestream.proto`. Eigener Befehl
`git ls-files sdks/csharp | grep -i proto` — **leer** (Exit 1, kein
Treffer): kein committeter Stub, keine committete Kopie der `.proto` im
SDK-Baum. **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.4 Tests — Nachrichtenschema-Vollständigkeit und Authn-Boundary

Eigene Lektüre von
`PgChangeFeed.Client.Tests/Grpc/ChangeMessageSchemaTests.cs`:
`AllTenSpec020FieldsRoundTrip` setzt und liest alle zehn Felder einzeln
gegen benannte Werte, zusätzlich `ChangeDescriptorHasExactlyTenFields`
(`Change.Descriptor.Fields.InFieldNumberOrder().Count == 10`) als
Drift-Wächter gegen ein künftig verändertes `.proto`. Eigene Lektüre von
`PgChangeFeedGrpcClientTests.cs`: drei Fakten-Tests gegen einen
`FakeCallInvoker` — Bearer-Token-Metadata-Form
(`StreamChangesAsync_SendsBearerTokenInAuthorizationMetadata`), Happy
Path mit vollem Inhalt
(`StreamChangesAsync_YieldsMessagesInOrderWithFullContent`), und die
Authn-Boundary
(`StreamChangesAsync_MissingOrInvalidToken_ThrowsUnauthenticated` — wirft
`RpcException` mit `StatusCode.Unauthenticated`, kein stiller leerer
Stream). Eigener Docker-Testlauf (Abschnitt 3 unten) bestätigt 32/32
grün. **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.5 Kein Import aus `internal/**`/`cmd/**`/`gen/**`

Eigener Befehl `grep -rn "internal/\|cmd/pg-change-feed\|gen/cdc"
sdks/csharp/` — genau **ein** Treffer:
`ChangeMessageSchemaTests.cs:10`, ein Doku-Kommentar-Zitat
(„analogous to `internal/adapters/driving/grpc/server_test.go`'s …
field-completeness assertion") — kein `using`, kein `ProjectReference`,
kein Import. Der generierte `Change`-Stub entsteht laut `.csproj`
(`<Protobuf Include="Grpc/proto/changestream.proto" …>`) und
Dockerfile-`COPY --from=proto`-Zeile ausschließlich aus der über den
Zusatzkontext gelesenen `.proto`, keine Kopie aus `gen/cdc/**`.
**Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.6 `make gates` grün

Eigener Lauf, ungepiped, Exit-Code direkt geprüft — siehe Abschnitt 2
unten. **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.7 Review durchgeführt, Report liegt vor

Eigene Lektüre von
`docs/reviews/review-slice-sdk-csharp-grpc-client-flaeche.md`: 0 HIGH/0
MEDIUM/0 LOW/0 INFO, explizit dokumentierte eigenständige Docker-Build-
und -Test-Läufe (mit und ohne Zusatzkontext), `make docs-check` real
gefahren. Der Nachtrag-Commit `899a3f80` behebt ausschließlich einen
`id-unlinked`-Docs-Check-Befund im Report selbst (fehlende Backticks um
eine ADR-Erwähnung) — 1 Zeile geändert, kein inhaltlicher
Verdikt-Wechsel (eigene Lektüre des `git show`-Diffs bestätigt das).
Die DoD-Zeile im Slice-Plan trägt bereits `[x]` mit korrektem Verweis auf
den Report-Pfad. **Ergebnis: Checkbox berechtigt auf `[x]`.**

### 1.8 Doku-Update `docs/user/benutzerhandbuch.md`

Eigener `git diff 87fa2f82..HEAD -- docs/user/benutzerhandbuch.md`:
Versions-Kopf `1.33` → `1.34`, neue Änderungshistorie-Zeile `| 1.34 |
2026-09-19 | C#-SDK-Hinweis für den gRPC-Change-Stream ergänzt … |`, und
im Fließtext ein neuer `**SDK:**`-Absatz direkt nach dem
`**Beispiele:**`-Block in §4 „Zugriff über den gRPC-Change-Stream":
nennt `PgChangeFeedGrpcClient.StreamChangesAsync`, den
`ChangeStream/StreamChanges`-RPC, `IAsyncEnumerable<Change>` mit „allen
zehn Feldern der Tabelle oben" — gegen die unmittelbar vorangehende
Feldtabelle gehalten: deckungsgleich, kein Drift. **Ergebnis: Checkbox
berechtigt auf `[x]`.**

### 1.9 Vier verbleibende `[ ]`-Checkboxen (Closure-Notiz, Reconciliation, Beobachtungs-Register, Risiko-Ausgänge, drei Paarungen)

Der Slice liegt weiterhin unter `docs/plan/planning/in-progress/`. §7
„Closure-Notiz" trägt weiterhin ausschließlich Platzhalter (`<…>`). Für
einen Slice, der noch nicht nach `done/` gewandert ist, ist das korrekt —
diese Punkte sind Closure-Pflichten, keine Liefer-Punkte; ihr
`[ ]`-Zustand widerspricht sich nicht mit dem übrigen DoD-Bild.
Reconciliation ist im Plan bereits explizit als „entfällt" markiert
(keine Reconciliation-Datei in diesem Repo, `MR-000`). **Kein
DoD-Verstoß für den aktuellen `in-progress`-Stand.**

## 2. `make gates` real, ungepiped ausgeführt

```
$ git log -1 --format=%H
899a3f8073cfaeba511e64020a10e145295ee4b1
$ git status --short
(leer)
$ make gates > /tmp/verifier-gates.log 2>&1; echo $?
0
```

Einzelbelege aus demselben Lauf (Log vollständig eingesehen):

| Gate | Ergebnis |
|---|---|
| `baseline-verify` | `v6.9.0 OK — 54 Dateien (Integritaet + Vollstaendigkeit, netzlos)` |
| `docs-check` | `d-check: 804 Datei(en) geprüft, 0 Befund(e)` |
| `commit-traceability` | d-check-Modul `commits`: `804 Datei(en) geprüft, 0 Befund(e)`; `commit-traceability.sh`: `OK — 5 Commit(s) in "HEAD~5..HEAD", Betreffs ohne Struktur-ID` |
| `coverage-gate` | `coverage-gate: OK — Coverage 82.80% erfüllt Schwelle 80%` |
| `generated-sync` | `generated-sync: OK — das committete Erzeugnis ist byte-gleich der Ausgabe des gepinnten Generators (Stufe proto)` |
| `a-check` | `gesamt: 0 Befund(e)` |

**Ergebnis: Die DoD-Checkbox „`make gates` grün" ist berechtigt auf
`[x]` gesetzt** — real, ungepiped, `EXIT=0`.

## 3. Fünfte unabhängige Docker-Build-Bestätigung (Verifier-eigen)

**3.1 Mit Zusatzkontext (Happy Path):**

```
$ docker build --no-cache --build-context proto=proto -f sdks/csharp/Dockerfile sdks/csharp
...
Passed!  - Failed: 0, Passed: 32, Skipped: 0, Total: 32, Duration: 331 ms
...
$ echo $?
0
```

Deckungsgleich mit Implementer- und Reviewer-Beleg (32/32 Tests grün,
`dotnet build` 0 Warnungen/0 Fehler). Test-Image danach mit
`docker image prune -f` entfernt.

**3.2 Ohne Zusatzkontext (Negativ-Fall, eigenständig nachgefahren — nicht
nur den Reviewer-Beleg übernommen):**

```
$ docker build --no-cache -f sdks/csharp/Dockerfile sdks/csharp
...
ERROR: failed to build: failed to solve: proto: failed to resolve source
metadata for docker.io/library/proto:latest: pull access denied,
repository does not exist or may require authorization: server message:
insufficient_scope: authorization failed
Dockerfile:57
--------------------
  57 | >>> COPY --from=proto cdc/stream/v1/changestream.proto PgChangeFeed.Client/Grpc/proto/changestream.proto
--------------------
$ echo $?
1
```

Der Bau bricht real und sichtbar exakt an der `COPY --from=proto`-Zeile
(57) ab — kein stiller Fallback, kein übersprungener Schritt. **Beleg für
die DoD-Zusage „ohne ihn bricht der Bau … ab" real und unabhängig vom
Reviewer erbracht** — dies ist die fünfte unabhängige Bestätigung nach
Implementer- und Reviewer-Läufen (jeweils mit und ohne Kontext).

## 4. Spec-/ADR-Konformität

### 4.1 ADR-0106 Festlegung 1 — kein Vorgriff auf SSE/NATS

Eigene Lektüre der Datei-Liste unter `sdks/csharp/`: kein `Sse/`- oder
`Nats/`-Verzeichnis, keine `NATS.Net`-Paketreferenz. `sdks/csharp/README.md`
§Status nennt SSE und NATS-Vollinhalt ausdrücklich als „remain uncovered
by this package" mit Verweis auf den direkten Draht-Zugriffsweg — konform
zu Festlegung 1. **Kein Vorgriff im Diff.**

### 4.2 ADR-0106 §5 — was unberührt bleibt

```
$ git diff 87fa2f82..HEAD --stat -- .a-check.yml spec/architecture.md docs/user/version.md
(leer)
$ git diff 87fa2f82..HEAD --stat -- examples/
(leer)
```

Alle drei in §5 genannten Pfade sind seit dem letzten Closure-Commit
unverändert; `examples/csharp/grpc-client/**` ebenfalls unverändert.
**Konform.**

### 4.3 Beide Vorgeschichte-Fehlerklassen — selbst, unabhängig vom Reviewer geprüft

- **README-Statusbehauptung (F-1-Klasse des Vorgänger-Slice):** eigene
  Lektüre von `sdks/csharp/README.md` §Status — nennt die gRPC-Fläche
  explizit als Teil des „current release", nicht als „coming soon"/
  „follow-up release". Eigener `grep -rniE "follow-up|not yet|coming
  soon|to be added|future release" sdks/csharp/` — zwei Treffer, beide
  in `PgChangeFeedClientOptions.cs`/`-Tests.cs` außerhalb dieses Diffs
  und ohne Bezug auf die gRPC-Fläche. **Fehlerklasse nicht reproduziert.**
- **Nicht-typisierter Fehlerpfad (F-2-Klasse des Vorgänger-Slice):**
  eigene Lektüre von `StreamChangesAsync` — der Authn-Boundary-Pfad wirft
  die bereits von `Grpc.Core` typisierte `RpcException`/
  `StatusCode.Unauthenticated` unverändert durch; kein zusätzlicher,
  selbst geschriebener Deserialisierungs- oder Mapping-Schritt, an dem
  eine analoge Lücke entstehen könnte. **Fehlerklasse nicht
  reproduziert.**

---

## Verdikt

**DoD erfüllt** (für den aktuellen `in-progress`-Stand des Slice — die
Closure-Pflichten in §2 sind bewusst noch offen, siehe 1.9). Alle
beauftragten Prüfpunkte wurden real und unabhängig nachgemessen, nicht
aus Bericht oder Commit-Message übernommen:

1. Die öffentliche Client-Klasse, ihre Bearer-Token-Metadata-Form und die
   Feldvollständigkeit sind gegen den Quellcode selbst geprüft.
2. Die drei gRPC-Pakete sind zentral in `Directory.Packages.props` mit
   exakten Versionen gepinnt (keine Bereiche).
3. Kein committeter Protobuf-Stub im SDK-Baum (`git ls-files sdks/csharp
   | grep proto` leer); der Dockerfile-Zusatzkontext ist real zwingend.
4. Beide Testklassen (Schema-Vollständigkeit, Authn-Boundary) existieren,
   sind inhaltlich vollständig und laufen real grün (32/32, eigener
   Docker-Lauf).
5. Die Import-Grenze ist eingehalten — ein Treffer, ein Doku-Zitat, kein
   Import.
6. `make gates` lief eigenständig, ungepiped, mit `EXIT=0`; alle sechs
   Gates einzeln bestätigt grün (baseline-verify, docs-check,
   commit-traceability, coverage-gate 82.80 %≥80 %, generated-sync,
   a-check).
7. Der Docker-Build wurde fünftfach unabhängig bestätigt: mit
   Zusatzkontext Exit 0/32 Tests grün, ohne Zusatzkontext Exit 1 mit
   Abbruch exakt an der `COPY --from=proto`-Zeile — beide Fälle in dieser
   Sitzung selbst nachgefahren, nicht nur aus dem Reviewer-Bericht
   übernommen.
8. `docs/user/benutzerhandbuch.md` ist korrekt auf 1.34 hochgezogen, die
   Changelog-Zeile und der Fließtext-Absatz sind inhaltlich zutreffend.
9. `ADR-0106` Festlegung 1 (kein SSE/NATS-Vorgriff) und §5 (unberührte
   Pfade: `.a-check.yml`, `spec/architecture.md`,
   `docs/user/version.md`, `examples/**`) sind eingehalten.
10. Beide in der Vorgeschichte benannten Fehlerklassen (veralteter
    README-Statusabsatz, nicht-typisierter Fehlerpfad) treten in diesem
    Slice nicht auf — unabhängig vom Reviewer selbst geprüft.

**Nicht blockierende Beobachtung:** `spec/pflichtenheft.md`s
`LH-FA-SST-009.a` trägt weiterhin wörtlich „offen, ADR-pflichtig", obwohl
`ADR-0106` diese Frage für C#/NuGet bereits beantwortet hat
(`AGENTS.md` §3.13-Klasse Träger-Nachzug, `ADR-0106` §Konsequenzen
Folgepflicht 2 — dort ausdrücklich „Sache des umsetzenden Zuges", nicht
Teil der DoD dieses Einzel-Slice). Dem Planner zur Kenntnis für die
Welle-Closure, kein DoD-Blocker dieses Slice.

**Freigabe an den Planner:** Der Slice ist DoD-konform für seinen
aktuellen `in-progress`-Stand. Vor Closure (`git mv` nach `done/`)
bleiben die im Plan selbst bereits als offen geführten
Closure-Pflichten zu erfüllen (§7 Closure-Notiz mit Steering-Loop-
Lerneintrag, Beobachtungs-Register-Prüfung, formaler
Risiko-Ausgangs-Häkchen-Nachzug in §2 — beide §6-Risiken tragen
inhaltlich bereits einen zulässigen Ausgang laut Plantext —, die drei
Paarungen bei der nächsten Welle-Closure).
