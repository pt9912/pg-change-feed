# Review-Report: slice-sdk-csharp-nats-stream-client-flaeche — 2026-09-22

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-csharp-nats-stream-client-flaeche.md`),
`ADR-0106` (Accepted) und `AGENTS.md` Hard Rules (Modul 10 §Drei
Review-Arten). Kein DoD-Abgleich als solcher — das ist Verifier-Aufgabe
(Modul 11); wo eine DoD-Zeile eine **Zahl** behauptet, ist das Nachmessen
dieser Zahl trotzdem Reviewer-Aufgabe (`AGENTS.md` §3.12, Reviewer-Skill
§HIGH „Zahl im Träger … driftend").

**Gegenstand:** Diff-Range `4b93def4..d86d1965`, ein Commit (`d86d1965`,
„feat(sdk): C#-NATS-Vollinhalts-Client-Fläche für PgChangeFeed.Client"),
Slice `slice-sdk-csharp-nats-stream-client-flaeche`, **letzter**
Flächen-Slice der Welle `welle-sdk-csharp-vollabdeckung` (23 Dateien,
968 Insertions / 47 Deletions — bündelt die NATS-Fläche selbst mit einem
breiten Träger-Nachzug).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, seither um mehrere weitere HIGH-Klassen
ergänzt — u. a. `AGENTS.md` §3.13 Träger-Nachzug, „Slice-/Wellen-Chronik in
Produktionscode-Kommentar", „Zahl im Träger … driftend").
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-22.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-sdk-csharp-nats-stream-client-flaeche.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan inkl. Plan-Nachzug, §6 Risiken, §7
  Closure-Notiz, §8 Sub-Area/Modus)
- `ADR-0106` (Accepted) — Festlegung 1 (letzter Absatz, NATS-Vollinhalt als
  Folge-Package antizipiert), Festlegung 2/3 (räumliche Trennung,
  Versionierung)
- `spec/pflichtenheft.md` §2 `SPEC-024` (Subjekt-Schema, Nachrichtenform,
  Authentifizierung), §6 `SPEC-026`-Zeile
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentar-Disziplin), §3.12
  (Herkunft von Aussagen), §3.13 (Träger-Nachzug)
- `harness/conventions.md` (MR-000 ID-Schema, MR-002 Slice-Kennungen)
- `harness/README.md` §Sensors/§Minimal-Agent-Workflow
- `examples/csharp/nats-stream-client/` als Draht-Vorbild
- `sdks/csharp/PgChangeFeed.Client/{Http,Grpc,Sse}/**` als bestehende,
  bereits gereviewte Flächen
- `docs/reviews/review-slice-sdk-csharp-sse-client-flaeche.md` (F-2,
  DRY-Kopie-Finding — Implementer-Behauptung „hier bewusst durch eine
  eigenständige Exception-Klasse vermieden" gegengeprüft) und
  `docs/reviews/review-slice-sdk-csharp-grpc-client-flaeche.md` (etablierte
  Kommentar-Anker-Form „(slice-<name>, ADR-Festlegung …)" als
  Vergleichs-Präzedenz)
- `docs/reviews/architect-verdict-slice-chronik-in-code-kommentar-4x.md`
  (Ursprung der HIGH-Regel „Slice-/Wellen-Chronik in
  Produktionscode-Kommentar")

**Eigenständig durchgeführte Prüfungen (nicht nur Commit-Message/DoD-Text
übernommen):**

- Alle zehn `SPEC-024`-Felder aus `spec/pflichtenheft.md` §`SPEC-024`
  einzeln gegen `Nats/Models/Change.cs`s `JsonPropertyName`-Attribute
  gehalten: `change_id`, `transaction_id`, `source_table_id`, `sequence`,
  `operation`, `old_image`, `new_image`, `schema_version`, `schema`,
  `table` — alle zehn vorhanden, keins fehlt, keins zusätzlich; zusätzlich
  durch `ChangeMessageSchemaTests.ChangeHasExactlyTenProperties`
  (Reflection-Zähler) abgesichert.
- Design-Entscheidung „eigenständige `PgChangeFeedNatsMalformedMessageException`
  statt Wiederverwendung der HTTP-`PgChangeFeedException`-Hierarchie"
  gegen die HTTP-Hierarchie selbst gehalten (`Http/PgChangeFeedException.cs`):
  deren Basistyp trägt ein `StatusCode`-Property — ein echtes HTTP-Konzept,
  das ein NATS-Message-Payload strukturell nicht hat. Anders als beim
  SSE-Slice (F-2 dort: wörtliche Kopie derselben `int`-`switch`-Logik ohne
  strukturellen Unterschied) besteht hier ein echter struktureller
  Unterschied zwischen den beiden Transportarten — die Design-Abweichung
  ist begründet, keine unbegründete Inkonsistenz.
- Design-Entscheidung „keine eigene Auth-Exception, reale
  `NATS.Net`-Ausnahme propagiert unverändert" gegen das bestehende
  gRPC-Muster gehalten (`Grpc/PgChangeFeedGrpcClient.cs`s Klassenkommentar
  zu `RpcException`/`Unauthenticated`): dieselbe „nicht in eine zweite
  Hierarchie remappen"-Semantik, real durch
  `PgChangeFeedNatsStreamClientAuthBoundaryTests.RejectedConnection_PropagatesNatsServerExceptionUnwrapped`
  belegt (`Assert.Same(authError, ex)` — dieselbe Instanz, kein Remapping).
  Beide Design-Entscheidungen sind konsistent zueinander und zu den
  bestehenden Flächen, keine unbegründete Design-Inkonsistenz.
- `NATS.Net`-Versionsmessung real nachvollzogen:
  `curl https://api.nuget.org/v3-flatcontainer/nats.net/index.json`
  (2026-09-22) — jüngster Eintrag weiterhin `3.2.0`, kein neuerer Kandidat
  seit der Beispiel-Messung vom 2026-09-17. Die Implementer-Behauptung ist
  real bestätigt, nicht nur übernommen.
- `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj`s `<Version>`
  real gelesen: `0.2.0` (vorher `0.1.0`), additive Erweiterung, kein
  Breaking Change an den drei bestehenden Flächen.
- Eigener `bash tools/harness/sdk-pack-csharp.sh`-Lauf (gecacht) plus ein
  zusätzlicher `docker build --no-cache -f sdks/csharp/Dockerfile --target
  build --build-context proto=proto sdks/csharp`, um einen Cache-Treffer
  auszuschließen: **70/70 Tests grün**, reales Artefakt
  `sdks/csharp/dist/PgChangeFeed.Client.0.2.0.nupkg` (37451 Bytes,
  `lib/net10.0/PgChangeFeed.Client.dll` enthalten, per `unzip -l` geprüft)
  — Version- und Test-Behauptung der DoD sind real bestätigt.
- **Test-Anzahl real nachgemessen statt aus dem DoD-Text übernommen**
  (siehe F-2 unten): über `git worktree add --detach <scratch> 4b93def4`
  (nicht-destruktiv, kein Checkout im Haupt-Arbeitsbaum) ein isolierter
  Docker-Bau des **Elternstands** gefahren — `dotnet test` meldet dort
  `Passed: 49, Total: 49`. Der committete Stand (`d86d1965`) meldet real
  `Passed: 70, Total: 70`. Delta: **21** neue tatsächlich laufende
  Testfälle, nicht die im Plan behaupteten „sechzehn". Worktree danach per
  `git worktree remove --force` wieder entfernt, Hauptarbeitsbaum
  durchgehend unverändert (`git status --short` vor und nach der Prüfung
  leer).
- Mutation-Testing-Plausibilisierung (nicht automatisiert mutiert, aber an
  ihrer Eingabeseite nachvollzogen): `ChangeMessageSchemaTests.AllTenSpec024FieldsRoundTrip`/
  `ChangeHasExactlyTenProperties` binden an ein reales JSON-Literal bzw. an
  `typeof(Change).GetProperties().Length` — eine entfernte Eigenschaft aus
  `Change` würde beide Tests rot zeigen, keine Ausgabeseiten-Mutation eines
  Fakes. `PgChangeFeedNatsStreamClientAuthBoundaryTests.RejectedConnection_…`
  steuert über `FakeNatsClient.WithFailure` eine reale
  `Exception`-Instanz, die aus der Aufzählung selbst geworfen wird — eine
  entfernte oder abgefangene Ausnahme in `StreamChangesAsync` würde diesen
  Test rot zeigen, keine Bindung nur an einen Stub-Rückgabewert.
  `JsonNullPayload_ThrowsMalformedMessage`/`NonJsonPayload_ThrowsMalformedMessage`
  steuern über echte Byte-Payloads (`"null"`, `"not json at all"`) den
  realen `ParseChange`-Codepfad — eine entfernte Null-Prüfung ließe den
  ersten Test grün bleiben (`JsonSerializer.Deserialize` gibt `null`
  zurück, keine Exception), der Test würde also tatsächlich rot, wenn der
  `?? throw`-Pfad entfernt würde. Alle drei vom Implementer behaupteten
  Mutationsklassen sind an ihrer Eingabeseite gebunden, keine
  Ausgabeseiten-Mutation.
- Eigener, **bewusst breiterer** `AGENTS.md`-§3.13-Suchlauf über den vom
  Implementer gefahrenen hinaus: zusätzliche Formulierungsvarianten
  (`"zwei Flächen"`, `"two surfaces"`, `"HTTP and gRPC"`,
  `"only HTTP and gRPC"`, `"nur HTTP.*gRPC"`) über den gesamten Baum
  (`--include="*.md" --include="*.cs" --include="*.sh" --include="*.mk"
  --include="*.props" --include="*.csproj" --include="Dockerfile"`),
  zusätzlich gezielt `docs/user/releasing.md`, `sdks/kotlin/**/README.md`
  gelesen. **Kein zusätzlicher, vom Implementer übersehener Fund** — die
  einzigen weiteren Treffer sind entweder korrekt weiterhin gültig (Kotlin
  deckt real nur HTTP+gRPC, `docs/user/releasing.md`s `sdk-csharp-v0.1.0`
  ist ein historischer Release-Beleg, kein Ist-Stand-Satz) oder liegen in
  `ADR-*`/`done/`/`docs/reviews/**` (unberührbare bzw. historische
  Records) — deckt sich mit der eigenen „bewusst nicht geändert"-Liste des
  Beobachtungs-Belegs.
- `git diff 4b93def4..d86d1965 -- tools/harness/sdk-pack-csharp.sh
  harness/mk/sdk.mk` Zeile für Zeile gelesen: ausschließlich
  Kommentar-Text (`BEIDE Testflächen` → `alle vier Testflaechen`,
  `0.1.0.nupkg` → `0.2.0.nupkg`) und keine Skript-/Makefile-**Logik**
  geändert — Out-of-Scope-Zusage („kein Pack-Werkzeug-Verhalten geändert")
  real bestätigt.
- `grep -rn "internal/\|cmd/\|gen/" sdks/csharp/ --include="*.cs"
  --include="*.md" --include="*.csproj" --include="*.props"` selbst
  gefahren: einziger Treffer ein bereits im Vorgänger-Review geprüftes
  Doku-Kommentar-Zitat außerhalb dieses Diffs
  (`Grpc/ChangeMessageSchemaTests.cs:10`), kein Treffer in den neuen
  `Nats/**`-Dateien.
- `git diff 4b93def4..d86d1965 --stat -- spec/architecture.md
  .a-check.yml` — leer, wie im Plan §1 vorab benannt (kein Bezug).
- Backtick-Zeichen je geänderter Markdown-Datei real ausgezählt (neun
  Dateien: `README.de.md` 62, `README.md` 62,
  `slice-sdk-csharp-nats-stream-client-flaeche.md` 464,
  `evidence/slice-sdk-csharp-nats-stream-client-flaeche.md` 112,
  `arbeit-ueberholt-stehenden-traeger/state.md` 190,
  `docs/user/benutzerhandbuch.md` 1848, `harness/README.md` 1804,
  `sdks/csharp/README.md` 54, `spec/pflichtenheft.md` 1470) — alle neun
  gerade, alle paarig.
- Neue `.cs`-Dateien (`Nats/**`, `PgChangeFeed.Client.Tests/Nats/**`,
  `Sse/Models/Change.cs`) nach `slice-`/`welle-`-Nennungen durchsucht —
  kein Treffer. Die vier Produktionsdatei-Kommentare außerhalb von `.cs`
  (`.csproj`, `Directory.Packages.props`, `Dockerfile`, `sdk.mk`,
  `sdk-pack-csharp.sh`) einzeln gegen das etablierte, im gRPC-Review
  bereits als unbeanstandet bestätigte Anker-Muster „(slice-<name>,
  ADR-Festlegung …)" gehalten — vier von fünf folgen diesem Muster
  unverändert; einer bricht daraus aus (siehe F-1).
- `make gates` ungefiltert laufen lassen, Exit-Code direkt (kein Pipe/
  Wrapper dazwischen) geprüft: `0` — u. a. `generated-sync: OK`,
  `a-check: gesamt: 0 Befund(e)`.

---

## Findings

### F-1 — `PgChangeFeed.Client.csproj`s Versionskommentar erzählt Vorher/Nachher statt Zustand zu nennen

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.7 / Reviewer-Skill HIGH „Slice-/Wellen-Chronik in
  Produktionscode-Kommentar"
- `pfad`: `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj:12-20`
- `befund`: Der Kommentar über der `<PropertyGroup>` mit `<Version>` wurde
  von einer reinen Zustandsaussage (Elternstand: „Version startet bei
  0.1.0 …") zu einer Vorher/Nachher-Erzählung umgeschrieben: „Version
  **startete** bei 0.1.0 … **und wurde** auf 0.2.0 **gehoben**
  (`slice-sdk-csharp-nats-stream-client-flaeche`): additive,
  rückwärtskompatible Erweiterung …". Das ist exakt die Konstellation, die
  der Architect-Verdikt zur Slice-Chronik in Code-Kommentaren als eigenen
  HIGH-Punkt benennt: eine Aussage über einen **Produktionscode-Pfad**
  (hier: eine Datei, `.csproj`) wird mit impliziter Vorher/Nachher-Sprache
  **und** einer Slice-Nummer begründet, statt mit `ADR-*`/`LH-*` oder dem
  Herkunfts-Anker `· seit slice-<NNN>`. Zum Vergleich: Der unmittelbar
  benachbarte, ebenfalls neue Kommentar über der `NATS.Net`-`ItemGroup`
  (Zeilen 61-63: „NATS.Net trägt die NATS-Vollinhalts-Client-Fläche
  (`slice-sdk-csharp-nats-stream-client-flaeche`, `SPEC-024`, `ADR-0106`
  §Entscheidung Festlegung 1 letzter Absatz)") folgt dagegen exakt dem im
  gRPC-Review als unbeanstandet bestätigten Anker-Muster „Name +
  ADR-Festlegung", ohne Vorher/Nachher-Erzählung — nur der
  Versionskommentar bricht aus diesem Muster aus.
- `verifizierbar`: nein — kein Gate prüft Kommentar-Chronik in diesem
  Repo (Reviewer ist die tragende Instanz, 4/4 historische Trefferquote
  laut Architect-Verdikt).
- `klasse`: „Slice-/Wellen-Chronik in Produktionscode-Kommentar"

### F-2 — DoD behauptet „sechzehn neue Tests", real gemessen sind es 21

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.12 / Reviewer-Skill HIGH „Zahl im Träger ohne
  Ursprung — oder gegen die Messung driftend"
- `pfad`: `docs/plan/planning/in-progress/slice-sdk-csharp-nats-stream-client-flaeche.md:77`
  (DoD-Punkt 1, „`LH-FA-SST-009` erfüllt")
- `befund`: Der Plan behauptet „sechzehn neue Tests in vier Dateien
  (`Nats/`)". Ein isolierter, nicht-destruktiver Docker-Testlauf gegen den
  Elternstand (`git worktree`, `4b93def4`) liefert real `Passed: 49, Total:
  49`; derselbe Lauf gegen den committeten Stand (`d86d1965`) liefert real
  `Passed: 70, Total: 70` (durch einen zusätzlichen `--no-cache`-Bau als
  Cache-Treffer ausgeschlossen). Delta: **21**, nicht 16. Auch die
  Methoden-Anzahl (Attribut-Zählung `[Fact]`/`[Theory]` ohne
  Theory-Fallauffächerung) über die vier neuen Testdateien ergibt **15**,
  ebenfalls nicht 16 — die Behauptung trifft unter keiner der beiden
  naheliegenden Zähl-Interpretationen zu. Die begleitende „70/70
  grün"-Aussage selbst ist korrekt (real bestätigt), nur die Teil-Zahl
  „sechzehn" driftet gegen die eigene Messung.
- `verifizierbar`: ja — `dotnet test` im gepinnten Bau zeigt die reale
  Zahl unmittelbar (bereits in diesem Review-Lauf ausgeführt).
- `klasse`: „Zahl im Träger driftet gegen die Messung"

### F-3 — Beobachtungs-Beleg zählt neun statt real zehn betroffene Dateien

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.12 / Reviewer-Skill HIGH-Klasse „Zahl im Träger
  … driftend", hier mit reduzierter Tragweite (internes Bookkeeping-Feld
  des Beobachtungs-Registers, kein Consumer-/DoD-Vertrag)
- `pfad`: `docs/plan/planning/observations/BEO-PGC/arbeit-ueberholt-stehenden-traeger/state.md:32-42`,
  `evidence/slice-sdk-csharp-nats-stream-client-flaeche.md:55-58`
- `befund`: Beide Texte behaupten „acht real behobene Fundstellen … neun
  Dateien". Reales Nachzählen der in derselben Aufzählung genannten Pfade
  ergibt zehn distinkte Dateien: `spec/pflichtenheft.md`,
  `sdks/csharp/README.md`, `PgChangeFeed.Client.csproj`,
  `sdks/csharp/Dockerfile`, `harness/mk/sdk.mk`,
  `tools/harness/sdk-pack-csharp.sh`, `harness/README.md`, `README.md`,
  `README.de.md`, `Sse/Models/Change.cs` — bestätigt auch durch
  `git diff 4b93def4..d86d1965 --stat`, das exakt diese zehn (plus die
  neuen `Nats/**`- und `benutzerhandbuch.md`-Änderungen, die nicht Teil
  dieser Aufzählung sind) zeigt. Die niedrigere Kategorie gegenüber F-2:
  Diese Zahl steht in einem internen Register-Beleg, nicht in einer
  DoD-Zeile, die einen Nutzen-/Abnahme-Anspruch trägt.
- `verifizierbar`: ja — `git diff --stat` zeigt die reale Dateiliste
  unmittelbar (bereits in diesem Review-Lauf ausgeführt).
- `klasse`: „Zahl im Träger driftet gegen die Messung"

## Negativbefunde

- geprüft, ohne Befund: `SPEC-024`-Vollständigkeit — alle zehn
  Nachrichtenfelder Feld für Feld gegen `spec/pflichtenheft.md` §`SPEC-024`
  gehalten, `Nats/Models/Change.cs` deckt sie exakt, zusätzlich durch
  einen Feldanzahl-Drift-Wächter abgesichert.
- geprüft, ohne Befund: Design-Konsistenz der beiden benannten
  Abweichungen (eigenständige Exception-Klasse für Malformed-Payload,
  keine eigene Auth-Exception) — beide sind begründet und konsistent zu
  den bestehenden Flächen, keine unbegründete Design-Inkonsistenz; die
  eigenständige Exception-Klasse vermeidet real den im SSE-Review als F-2
  benannten DRY-Kopie-Effekt, weil hier (anders als bei SSE/HTTP) ein
  echter struktureller Unterschied besteht (kein `StatusCode`-Äquivalent
  bei NATS).
- geprüft, ohne Befund: `NATS.Net`-Versionsmessung — real gegen
  `api.nuget.org` nachverifiziert, `3.2.0` weiterhin jüngste stabile
  Version, keine Pin-Hebung nötig.
- geprüft, ohne Befund: Version-Bump — `PgChangeFeed.Client.csproj`s
  `<Version>` real `0.1.0` → `0.2.0`.
- geprüft, ohne Befund: reales `.nupkg` — eigener, nicht-gecachter Bau
  bestätigt `PgChangeFeed.Client.0.2.0.nupkg` mit allen vier
  Client-Flächen in einer einzigen Assembly.
- geprüft, ohne Befund: Mutation-Testing-Plausibilisierung — alle drei
  behaupteten Mutationsklassen (Nullpfad-Entfernung, geschluckte
  Auth-Exception, entferntes Feld) sind an ihrer jeweiligen Eingabeseite
  gebunden, keine Bindung ausschließlich an eine Ausgabeseite/einen
  Fake-Rückgabewert.
- geprüft, ohne Befund: Träger-Nachzug (`AGENTS.md` §3.13) — eigener,
  bewusst breiterer Suchlauf über zusätzliche Formulierungsvarianten
  („zwei Flächen", „two surfaces", „HTTP and gRPC", „only HTTP and gRPC")
  plus gezielte Prüfung von `docs/user/releasing.md` und
  `sdks/kotlin/**/README.md` fand **keine** vom Implementer übersehene
  Stelle; alle verbleibenden Treffer sind entweder weiterhin zutreffend
  (Kotlin) oder unberührbare/historische Records (ADRs, `done/`,
  `docs/reviews/**`, datierte Änderungshistorie-Zeilen).
- geprüft, ohne Befund: Import-Grenze (`internal/**`/`cmd/**`/`gen/**`) —
  eigener `grep`-Lauf, einziger Treffer außerhalb dieses Diffs (bereits im
  Vorgänger-Review geprüftes Doku-Zitat), kein Treffer in `Nats/**`.
- geprüft, ohne Befund: Out-of-Scope — kein realer Tag-Push, keine
  Verhaltensänderung an `tools/harness/sdk-pack-csharp.sh`/`harness/mk/sdk.mk`
  (Zeile-für-Zeile-Diff bestätigt: ausschließlich Kommentartext- und
  Versionsnummer-Updates, keine Logikänderung), kein Bezug zu
  `spec/architecture.md`/`.a-check.yml`.
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — Version 1.37 →
  1.38, Stand 2026-09-22, neue Änderungshistorie-Zeile vorhanden; beide
  neuen `**SDK:**`-Absätze (SSE, NATS-Vollinhalt) nennen „alle zehn
  Felder(n) der Tabelle oben" korrekt gegen die jeweils referenzierte
  Feldliste gehalten.
- geprüft, ohne Befund: Backtick-Parität — alle neun in diesem Diff
  geänderten Markdown-Dateien real nachgezählt, alle neun paarig.
- geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7) außerhalb
  von F-1 — kein Konjunktiv über eine verworfene Alternative, kein
  abwesender Text, keine Chronik in den neuen `.cs`-Produktionsdateien;
  die übrigen vier geänderten Nicht-`.cs`-Produktionskommentare
  (`Directory.Packages.props`, `Dockerfile`, `sdk.mk`,
  `sdk-pack-csharp.sh`) folgen dem etablierten Anker-Muster ohne
  Vorher/Nachher-Sprache.
- geprüft, ohne Befund: Traceability — Commit-Betreff `d86d1965` nennt
  `LH-FA-SST-009`/`ADR-0106`, kein `SPEC-*`/`ARC-*` im Betreff.
- geprüft, ohne Befund: `make gates` real gefahren, Exit-Code direkt
  (ungepiped) geprüft, `0`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 2 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** „Slice-/Wellen-Chronik in
Produktionscode-Kommentar" · „Zahl im Träger driftet gegen die Messung"
(2×, unterschiedliche Trägerklasse: DoD-Zeile vs. Beobachtungs-Beleg)

## Verdikt

**Merge-blockierend:** ja — zwei HIGH-Findings. Beide sind inhaltlich
schmal (eine Kommentar-Umformulierung, eine Zahlenkorrektur in einer
DoD-Zeile plus einer Register-Notiz) und berühren keinen Sicherheits- oder
Korrektheitspfad — die NATS-Fläche selbst (Nachrichtenschema,
Subjekt-Formatierung, Auth-Boundary, Fehlerpfad) ist vollständig,
konsistent designt und real durch 70/70 grüne Tests sowie ein reales
`.nupkg` belegt. Trotzdem: Nach dem Reviewer-Skill werden HIGH-Findings
nicht wegen Plausibilität oder geringer Tragweite herabgestuft — beide
Klassen haben eine etablierte, mehrfach bestätigte Präzedenz
(„Slice-/Wellen-Chronik in Produktionscode-Kommentar": 4/4 historische
Trefferquote vor Merge; „Zahl im Träger … driftend": 4× bereits vom
Reviewer durch eigenes Nachmessen gefunden).

**Übergabe:** Rückgabe-Pfeil an den Implementer für eine gezielte
Fixrunde — kein Rollen-Widerspruch (keine Implementer-Gegenrede zu
erwarten, beide Findings sind mechanisch nachvollziehbar), daher **keine**
Konflikt-Pfad-Sequenz über den Architect nötig (Modul 8: Konflikt-Pfad
greift bei Rollen-Widerspruch oder 3× gleichem Konflikttyp, hier liegt
keines von beiden vor). Erwarteter Umfang der Fixrunde:

1. F-1: `PgChangeFeed.Client.csproj`s Versionskommentar auf eine reine
   Zustandsaussage zurückführen (z. B. „Version `0.2.0` — additive
   Erweiterung um SSE-/NATS-Vollinhalts-Fläche, kein Breaking Change" mit
   `ADR-0106` Festlegung 3 und/oder dem Anker `· seit
   slice-sdk-csharp-nats-stream-client-flaeche`, ohne die
   „startete …/wurde … gehoben"-Erzählung) — Lösungsentscheidung liegt
   beim Implementer, nicht bei diesem Report.
2. F-2: DoD-Punkt 1 im Slice-Plan auf die real gemessene Zahl korrigieren
   (21 tatsächliche Testfälle bzw. 15 neue Testmethoden — Implementer
   wählt die zutreffende Zähleinheit und macht sie explizit, statt eine
   dritte, ungeprüfte Zahl zu setzen).
3. F-3: Beobachtungs-Beleg (`state.md` und die Evidenz-Datei) auf die
   real gezählten zehn Dateien korrigieren.

Die Finding-Klassen gehen in die Slice-Closure §7 und von dort in den
Steering-Loop-Zähler. Dieser Report ist ein Lauf-Beleg; er ersetzt keine
Verifikation gegen die volle DoD — das bleibt Verifier-Aufgabe (Modul 11).

**DoD-Checkbox-Nachzug:** entfällt — dieses Verdikt führt zu einer
Fixrunde, die DoD-Zeile „Review durchgeführt" wird regulär bei Schritt 21
des Implementer-Workflows nachgezogen (Reviewer-Skill §DoD-Checkbox-Nachzug
ohne Fixrunde, Grenzfall „nur ohne Fixrunde").
