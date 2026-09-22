# Review-Report: slice-sdk-kotlin-nats-stream-client-flaeche — 2026-09-22

**Review-Art:** Code — geprüft gegen Plan
(`docs/plan/planning/in-progress/slice-sdk-kotlin-nats-stream-client-flaeche.md`),
`ADR-0109` (Accepted) und `AGENTS.md` Hard Rules (Modul 10 §Drei
Review-Arten). Kein DoD-Abgleich als solcher — das ist Verifier-Aufgabe
(Modul 11); wo eine DoD-Zeile eine **Zahl** behauptet, ist das Nachmessen
dieser Zahl trotzdem Reviewer-Aufgabe (`AGENTS.md` §3.12, Reviewer-Skill
§HIGH „Zahl im Träger … driftend").

**Gegenstand:** Commit `c8c9e3ae` („feat(sdk-kotlin): NATS-Vollinhalts-
Client-Fläche + Version-Hebung 0.2.0"), Diff-Range `e8eddf79..c8c9e3ae`,
Slice `slice-sdk-kotlin-nats-stream-client-flaeche`, **letzter**
Flächen-Slice der Welle `welle-sdk-kotlin-vollabdeckung` (14 Dateien, 787
Insertions / 49 Deletions — bündelt die NATS-Fläche selbst mit einem
breiten Träger-Nachzug für zwei Flächen, SSE + NATS).

**Skill:** `.harness/skills/reviewer.md` @ Fassung „geschärft 2026-09-09"
(vier repo-spezifische HIGH-Regeln, seither um mehrere weitere HIGH-Klassen
ergänzt — u. a. `AGENTS.md` §3.13 Träger-Nachzug, „Slice-/Wellen-Chronik in
Produktionscode-Kommentar", „Zahl im Träger … driftend").
**Modell:** claude-sonnet-5 · **Datum:** 2026-09-22.

**Eingangs-Kontext:**

- `docs/plan/planning/in-progress/slice-sdk-kotlin-nats-stream-client-flaeche.md`
  (§1 Ziel/Abgrenzung, §2 DoD, §3 Plan, §6 Risiken, §7 Closure-Notiz, §8
  Sub-Area/Modus)
- `ADR-0109` (Accepted) — Festlegung 1 (letzter Absatz, NATS-Vollinhalt als
  Folge-Package antizipiert), Festlegung 2/3/5
- `spec/pflichtenheft.md` §2 `SPEC-024` (Subjekt-Schema, Nachrichtenform,
  Authentifizierung), §6 `SPEC-028`-Zeile
- `AGENTS.md` §3.1 (Docker-only), §3.7 (Kommentar-Disziplin), §3.12
  (Herkunft von Aussagen), §3.13 (Träger-Nachzug)
- `harness/conventions.md` (MR-000 ID-Schema, MR-002 Slice-Kennungen)
- `harness/README.md` §Sensors/§Minimal-Agent-Workflow
- `examples/kotlin/nats-stream-client/` als Draht-Vorbild
- `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/.../{http,grpc,sse}/**`
  als bestehende, bereits gereviewte Flächen
- `docs/reviews/review-slice-sdk-csharp-nats-stream-client-flaeche.md` als
  Kontrastfolie (analoger, letzter Flächen-Slice der C#-Welle, 2 HIGH +
  1 MEDIUM)

**Eigenständig durchgeführte Prüfungen (nicht nur Commit-Message/DoD-Text
übernommen):**

- Alle zehn `SPEC-024`-Felder aus `spec/pflichtenheft.md` §`SPEC-024`
  einzeln gegen `nats/model/Change.kt`s `@SerializedName`-Annotationen
  gehalten: `change_id`, `transaction_id`, `source_table_id`, `sequence`,
  `operation`, `old_image`, `new_image`, `schema_version`, `schema`,
  `table` — alle zehn vorhanden, keins fehlt, keins zusätzlich; zusätzlich
  Feld-für-Feld-Assertion in
  `PgChangeFeedNatsStreamClientMessageSchemaTest.streamChanges yields all
  SPEC-024 fields unchanged` gelesen.
- Subjekt-Schema real gegen `SPEC-024` gehalten: `buildSubject` erzeugt
  exakt `cdc.stream.<source_id>.<schema>.<table>`, `buildSourceSubject`
  exakt `cdc.stream.<source_id>.>`, `ALL_SOURCES_SUBJECT` exakt
  `cdc.stream.>` — alle drei Formen decken sich wörtlich mit der
  `SPEC-024`-Tabellenzeile „Subjekt-Schema".
- Authentifizierung gegen `SPEC-024` gehalten: Verbindungsebene über
  `Options.Builder().token(options.apiToken)`, kein Per-Message-Header —
  entspricht der Festlegung „Verbindungsebene … kein Verbindungsversuch
  ohne oder mit falschem Token wird vom Server abgelehnt".
- **Test-Anzahl real nachgemessen statt aus dem Commit-Text übernommen**
  (siehe Negativbefund unten): eigener, nicht-destruktiver Docker-Bau des
  **Elternstands** über `git worktree add --detach` (`e8eddf79`,
  danach `git worktree remove --force`, Hauptarbeitsbaum durchgehend
  unverändert) — `./gradlew test` liefert dort real (JUnit-XML-Summe über
  `build/test-results/test/*.xml`) **47 Tests, 0 Fehler**. Derselbe Lauf
  gegen den committeten Stand (`c8c9e3ae`) liefert real **61 Tests, 0
  Fehler**. Delta: **14**, exakt wie behauptet — aufgeschlüsselt real
  nachgezählt: `PgChangeFeedNatsStreamClientAuthBoundaryTest` 4,
  `PgChangeFeedNatsStreamClientMessageSchemaTest` 3,
  `PgChangeFeedNatsStreamClientSubjectTest` 7 → 4+3+7 = 14. Commit-Message
  und DoD-Zeile 1 sind real bestätigt, keine Drift.
- `io.nats:jnats`-Versionsmessung real nachvollzogen:
  `curl https://repo1.maven.org/maven2/io/nats/jnats/maven-metadata.xml`
  (2026-09-22) — `<latest>`/`<release>` weiterhin `2.26.3`. Die
  Implementer-Behauptung „keine Drift" ist real bestätigt, nicht nur
  übernommen.
- Eigener `bash tools/harness/sdk-pack-kotlin.sh`-Lauf: reales Artefakt
  `sdks/kotlin/dist/pgchangefeed-kotlin-0.2.0.jar` (144598 Bytes, per
  `unzip -l` geprüft — Klassen aus allen vier Paketen `http`/`grpc`/`sse`/
  `nats` im selben Jar).
- `build.gradle.kts`s Top-Level-`version`- **und**
  `publishing.publications.maven.version`-Zeile real gelesen: beide
  `0.2.0` (vorher `0.1.0`), additive Erweiterung.
- **`internal`-Transport-Abstraktion real gegen einen frischen Build
  geprüft** (dritte Anwendung dieser Lektion in dieser Welle, nach den
  gRPC-/SSE-Kotlin-Reviews): `javap -p` gegen die aus dem realen `.jar`
  entpackten `.class`-Dateien zeigt `NatsStreamTransport` als
  `public interface` und den `internal constructor(transport:
  NatsStreamTransport)` von `PgChangeFeedNatsStreamClient` als gewöhnlichen
  `public`-`<init>`-Bytecode-Eintrag — exakt wie die KDoc behauptet
  („compile-time Kotlin-Sichtbarkeitsgrenze, keine JVM-Bytecode-Grenze").
  Der separat deklarierte `private constructor(owned: Pair<…>)` erscheint
  dagegen real als `private` im Bytecode — ebenfalls wie dokumentiert.
- Design-Bewertung des als „ungewöhnlich" gemeldeten Pair-basierten
  privaten Sekundär-Konstruktors: `Nats.connect()` (anders als der C#
  Sibling, der laut Closure-Notiz asynchron verbindet) läuft synchron;
  der Convenience-Konstruktor `constructor(options: …) : this(connectOwned(options))`
  müsste ohne den Pair-Umweg entweder `connectOwned(options)` zweimal
  aufrufen (zwei reale `Nats.connect()`-Verbindungen für ein Objekt — ein
  echter Bug, nicht nur Stil) oder Kotlins Regel verletzen, dass ein
  Sekundärkonstruktor als **einzige** erste Anweisung ausschließlich einen
  `this(…)`/`super(…)`-Aufruf tragen darf (kein beliebiger Code davor).
  Der Pair ist die knappste Lösung für „eine seiteneffektbehaftete
  Berechnung liefert zwei Konstruktor-Parameter, aber nur ein
  Delegationsaufruf ist erlaubt" — ein bekanntes, wenn auch seltenes
  Kotlin-Idiom. Die Alternative (eine `companion`-Factory-Funktion statt
  eines öffentlichen Konstruktors) hätte `new PgChangeFeedNatsStreamClient(options)`
  für Java-Aufrufer verhindert. Bewertung: saubere, dokumentierte Lösung
  für ein reales Arity-Problem, kein Wartbarkeitsrisiko — die KDoc dieses
  Konstruktors benennt sowohl den Grund als auch die Abgrenzung zum
  gewöhnlichen `Connection`-Konstruktor explizit.
- `grep -rn "internal/\|cmd/\|gen/"` über alle neuen `.kt`-Dateien
  (`nats/**`, Test-Pendants) selbst gefahren: kein Treffer — alle Importe
  sind `io.github.pt9912.pgchangefeed.*`, `com.google.gson.*`,
  `io.nats.client.*`, `java.*`/`kotlin.*`.
- `git diff e8eddf79..c8c9e3ae -- harness/mk/sdk.mk` Zeile für Zeile
  gelesen: ausschließlich Kommentar-Text (`BEIDE Testflächen (HTTP + gRPC)`
  → `VIER Testflächen (HTTP + gRPC + SSE + NATS-Vollinhalt)`,
  `0.1.0.jar` → `0.2.0.jar`) und keine Skript-/Makefile-**Logik**
  geändert — Out-of-Scope-Zusage bestätigt.
- `git diff e8eddf79..c8c9e3ae --stat -- spec/architecture.md
  .a-check.yml .github/workflows/ tools/harness/sdk-pack-kotlin.sh` — leer,
  wie im Plan §1 vorab benannt (kein Bezug, keine
  Pack-/Publish-Verhaltensänderung).
- Backtick-Zeichen je geänderter Markdown-Datei real ausgezählt: fünf
  Dateien (`slice-sdk-kotlin-nats-stream-client-flaeche.md` 324,
  `docs/user/benutzerhandbuch.md` 1914, `harness/README.md` 1804,
  `sdks/kotlin/pgchangefeed-kotlin/README.md` 98, `spec/pflichtenheft.md`
  1486) — alle fünf gerade, alle paarig.
- **Eigener, bewusst breiterer `AGENTS.md`-§3.13-Suchlauf** über den vom
  Implementer gefahrenen hinaus (siehe F-2 unten für den einen real
  gefundenen Fehltreffer): zusätzliche Formulierungsvarianten
  (`HTTP.*gRPC`/`HTTP.*und.*gRPC`) über den gesamten Baum
  (`--include="*.md" --include="*.mk" --include="*.sh"
  --include="*.kts" --include="Dockerfile" --include="*.yml"`),
  zusätzlich gezielt `docs/user/releasing.md`,
  `.github/workflows/sdk-kotlin-release.yml`,
  `docs/plan/planning/welle-sdk-kotlin-vollabdeckung.md` gelesen. Die
  einzigen weiteren Treffer außer dem in F-2 benannten sind entweder
  bereits nachgezogen (die fünf im Diff geänderten Träger), datierte
  Historie-Einträge (`spec/pflichtenheft.md`s Änderungshistorie,
  `docs/user/benutzerhandbuch.md`s Versionshistorie) oder illustrative
  Beispiele ohne Ist-Stand-Anspruch (`docs/user/releasing.md`s
  „z. B. `sdk-kotlin-v0.1.0`", `welle-sdk-kotlin-vollabdeckung.md`s §1/§2
  Kontext-Absätze, explizit als Zustand „bei Welle-Eröffnung,
  2026-09-21" datiert).
- `jnats`-Deprecation-Warnung real geprüft statt übernommen: eigener
  `docker build --no-cache`-Lauf zeigt real
  `w: … 'fun token(p0: String!): Options.Builder!' is deprecated.
  Deprecated in Java.` an `PgChangeFeedNatsStreamClient.kt:206` (der
  `.token(options.apiToken)`-Aufruf in `connectOwned`). Real gegen
  `examples/kotlin/nats-stream-client/src/main/kotlin/cdcexamples/natsstream/Main.kt:52`
  gehalten: dieselbe `Options.Builder().server(…).token(…).build()`-Zeile
  existiert dort bereits unverändert — die Implementer-Einstufung „bereits
  im Beispiel vorhanden, kein neues Problem" ist real bestätigt.
- Mutation-Testing-Plausibilisierung nachgelesen (nicht selbst neu
  mutiert, aber Testcode und Assertion an ihrer Eingabeseite geprüft):
  `PgChangeFeedNatsStreamClientMessageSchemaTest`s Assertions binden an
  ein reales JSON-Literal je Feld — eine entfernte/geleerte Eigenschaft
  würde die jeweilige `assertEquals`-Zeile treffen, keine
  Ausgabeseiten-Mutation eines Fakes.
  `PgChangeFeedNatsStreamClientAuthBoundaryTest`s
  `a rejected connection propagates the underlying NATS exception
  unwrapped` steuert über `FakeNatsStreamTransport.withFailure` eine reale
  `IOException`-Instanz, die aus dem Next-Payload-Supplier selbst
  geworfen wird, und prüft `assertSame(authError, ex)` — ein verschluckter
  `catch`-Block im echten `streamChanges`-Code würde diesen Test treffen,
  keine Bindung nur an einen Stub-Rückgabewert. Beide vom Implementer
  benannten, real geprüften Mutationen (Feld-Leerung, verschluckte
  Exception) sind an ihrer Eingabeseite gebunden.
- `make gates` ungefiltert laufen lassen, Exit-Code direkt (kein Pipe/
  Wrapper dazwischen) geprüft: `0` — u. a. `generated-sync: OK`,
  `a-check: gesamt: 0 Befund(e)`, `coverage-gate: OK — 82.70%`.

---

## Findings

### F-1 — `build.gradle.kts`s Version-Hebungs-Kommentar erzählt Vorher/Nachher statt Zustand zu nennen

- `kategorie`: HIGH
- `quelle`: `AGENTS.md` §3.7 / Reviewer-Skill HIGH „Slice-/Wellen-Chronik in
  Produktionscode-Kommentar"
- `pfad`: `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts:97-102`
- `befund`: Der zweite Absatz des neuen Kommentarblocks über `plugins {}`
  lautet: „Diese Version-Hebung (`0.1.0` -> `0.2.0`) trägt außerdem den
  SSE-Client-Fläche (`slice-sdk-kotlin-sse-client-flaeche`, `SPEC-021`) —
  beide Flächen bündeln ihren Version-Bump gemeinsam in diesem letzten
  Flächen-Slice der Welle (`welle-sdk-kotlin-vollabdeckung` §1). Additiv,
  rückwärtskompatible Erweiterung — SemVer-Minor, keine ADR-pflichtige
  Ausnahme (Slice-Plan §6)." Das Arrow-Muster `0.1.0 -> 0.2.0` ist
  implizite Vorher/Nachher-Sprache über den Produktionscode-Pfad (die
  `version`-Zeile selbst), und die Aussage wird ausschließlich mit
  Slice-Kennungen (`slice-sdk-kotlin-sse-client-flaeche`) und einem Verweis
  auf den Slice-Plan (`Slice-Plan §6`) begründet — nicht mit `ADR-*`/`LH-*`
  oder dem Herkunfts-Anker `· seit slice-<NNN>`. Exakt die Konstellation,
  die der Architect-Verdikt zur Slice-Chronik in Code-Kommentaren als
  eigenen HIGH-Punkt benennt, und dieselbe Fehlerklasse wie F-1 des
  analogen C#-NATS-Reviews (dort: „Version startete bei 0.1.0 … und wurde
  auf 0.2.0 gehoben"). Zum Vergleich: Der unmittelbar vorangehende,
  ebenfalls neue Absatz über der `jnats`-`ItemGroup`-Zeile folgt dagegen
  dem etablierten, in diesem Repo bereits mehrfach unbeanstandet
  verwendeten Muster „Name + ADR-Festlegung + reale
  Messung-mit-Lauf-Datum" (`AGENTS.md` §3.12-konform) — nur der zweite,
  auf die Versionszahl selbst bezogene Absatz bricht daraus aus.
- `verifizierbar`: nein — kein Gate prüft Kommentar-Chronik in diesem
  Repo (Reviewer ist die tragende Instanz).
- `klasse`: „Slice-/Wellen-Chronik in Produktionscode-Kommentar"

### F-2 — `sdks/kotlin/Dockerfile` trägt zwei jetzt stale `pgchangefeed-kotlin-0.1.0.jar`-Kommentarzeilen, vom §3.13-Suchlauf übersehen

- `kategorie`: MEDIUM
- `quelle`: `AGENTS.md` §3.13 (Träger-Nachzug)
- `pfad`: `sdks/kotlin/Dockerfile:77`, `sdks/kotlin/Dockerfile:104`
- `befund`: Zwei Kommentarzeilen in `sdks/kotlin/Dockerfile` (Stufen
  `pack`/`publish`, außerhalb dieses Diffs) nennen wörtlich
  „`build/libs/pgchangefeed-kotlin-0.1.0.jar`" als Beispiel des vom
  Standard-`jar`-Task erzeugten Dateinamens — real gegen den
  gemessenen, aktuellen Jar-Namen (`pgchangefeed-kotlin-0.2.0.jar`, siehe
  Eingangsprüfung oben) driftend. Die beiden **sibling**-Träger, die
  denselben Sachverhalt (Jar-Name des Kotlin-SDK-Packages) beschreiben —
  `harness/README.md`s `make sdk-pack-kotlin`-Zeile und
  `harness/mk/sdk.mk`s Kommentarblock — **wurden** in genau diesem Diff
  korrekt von `0.1.0.jar` auf `0.2.0.jar` nachgezogen; der eigene, breitere
  Suchlauf dieses Reviews (`grep -rn "pgchangefeed-kotlin.*0\.1\.0"`, s. o.)
  fand als einzigen zusätzlichen, real übersehenen Treffer diese beiden
  Dockerfile-Zeilen. Die Closure-Notiz (§7 des Slice-Plans) behauptet „der
  vollständige Suchlauf … deckte alle stehenden, jetzt falschen Träger ab
  … kein Fund, der über dieses Slice hinaus offen bliebe" — diese Aussage
  trifft für die genannte Datei nicht zu. **Klassifikations-Hinweis:** Der
  Reviewer-Skill sieht für die verwandte HIGH-Klasse „Zahl im Träger …
  driftend" (`AGENTS.md` §3.12) einen expliziten Herabstufungs-Passus vor
  („Träger außerhalb des Diffs … bleiben INFO"); dieser Fund betrifft
  jedoch eine andere Hard Rule (`AGENTS.md` §3.13, Träger-Nachzug-Pflicht,
  keine eigene benannte Skill-Bullet) und liegt inhaltlich näher an einer
  fehlgeschlagenen, für dieses Slice ausdrücklich mandatierten
  Vollständigkeits-Suche als an einem beiläufig entdeckten, unverbundenen
  Fremd-Träger — deshalb MEDIUM statt INFO oder HIGH: real, konkret,
  kostengünstig zu beheben, aber ohne funktionale/Sicherheits-Auswirkung
  und in einer Datei außerhalb des Diffs.
- `verifizierbar`: ja — `grep -n "0\.1\.0" sdks/kotlin/Dockerfile` zeigt
  die beiden Treffer unmittelbar (bereits in diesem Review-Lauf
  ausgeführt).
- `klasse`: „Arbeit überholt stehenden Träger" (`BEO-PGC/arbeit-ueberholt-stehenden-traeger`)

## Negativbefunde

- geprüft, ohne Befund: `SPEC-024`-Vollständigkeit — alle zehn
  Nachrichtenfelder Feld für Feld gegen `spec/pflichtenheft.md` §`SPEC-024`
  gehalten, `nats/model/Change.kt` deckt sie exakt, zusätzlich durch
  `PgChangeFeedNatsStreamClientMessageSchemaTest` abgesichert.
- geprüft, ohne Befund: Subjekt-Schema und Authentifizierung — beide real
  gegen `SPEC-024` gehalten, keine Abweichung.
- geprüft, ohne Befund: Test-Anzahl — eigener, isolierter
  Elternstand-Bau (`git worktree`, `e8eddf79`) liefert real 47 Tests,
  committeter Stand (`c8c9e3ae`) real 61 Tests, Delta 14 = 4 (Auth) + 3
  (Schema) + 7 (Subjekt) — exakt wie im Commit/DoD behauptet, keine Drift.
- geprüft, ohne Befund: `jnats`-Versionsmessung — real gegen
  `repo1.maven.org` nachverifiziert, `2.26.3` weiterhin jüngste stabile
  Version, keine Pin-Hebung nötig.
- geprüft, ohne Befund: Version-Bump — Top-Level **und**
  `publishing`-Block-Koordinate beide real `0.1.0` → `0.2.0`.
- geprüft, ohne Befund: reales `.jar` — eigener `sdk-pack-kotlin.sh`-Lauf
  bestätigt `pgchangefeed-kotlin-0.2.0.jar` (144598 Bytes) mit allen vier
  Client-Flächen im selben Jar.
- geprüft, ohne Befund: `internal`-Transport-Abstraktion — `javap -p`
  gegen die realen `.class`-Dateien bestätigt die KDoc-Behauptung
  (compile-time Kotlin-Grenze, `public` im Bytecode) exakt.
- geprüft, ohne Befund: Pair-basierter privater Sekundärkonstruktor —
  begründete, dokumentierte Lösung für ein reales Arity-/Seiteneffekt-
  Problem (`Nats.connect()` synchron, Kotlin erlaubt keinen Code vor
  einem `this(…)`-Delegationsaufruf), kein Wartbarkeitsrisiko.
- geprüft, ohne Befund: Import-Grenze (`internal/**`/`cmd/**`/`gen/**`) —
  eigener `grep`-Lauf über alle neuen `.kt`-Dateien, kein Treffer.
- geprüft, ohne Befund: Out-of-Scope — kein realer Tag-Push, keine
  Verhaltensänderung an `tools/harness/sdk-pack-kotlin.sh`/
  `harness/mk/sdk.mk` (Zeile-für-Zeile-Diff bestätigt: ausschließlich
  Kommentartext- und Versionsnummer-Updates), kein Bezug zu
  `spec/architecture.md`/`.a-check.yml`/`.github/workflows/`, kein
  Maven-Central-Wechsel.
- geprüft, ohne Befund: `docs/user/benutzerhandbuch.md` — Version 1.38 →
  1.39, Stand 2026-09-22, neue Änderungshistorie-Zeile vorhanden; beide
  neuen `**SDK:**`-Absätze (SSE, NATS-Vollinhalt) korrekt gegen die
  jeweils referenzierte Zehn-Felder-Tabelle gehalten.
- geprüft, ohne Befund: `spec/pflichtenheft.md`/`SPEC-028`/
  `LH-FA-SST-009.a` — Träger-Nachzug korrekt („HTTP-API, gRPC-Stream, SSE
  und NATS-Vollinhalt" statt „HTTP-API und gRPC-Stream", Version
  `0.1.0` → `0.2.0`, neue Änderungshistorie-Zeile datiert 2026-09-22).
- geprüft, ohne Befund: Träger-Nachzug (`AGENTS.md` §3.13) darüber
  hinaus — eigener, bewusst breiterer Suchlauf über zusätzliche
  Formulierungsvarianten und weitere Dateien
  (`docs/user/releasing.md`, `.github/workflows/sdk-kotlin-release.yml`,
  `welle-sdk-kotlin-vollabdeckung.md`) fand **außer dem in F-2 benannten
  Fund** keine vom Implementer übersehene Stelle; alle übrigen Treffer sind
  entweder bereits nachgezogen, datierte Historie-Einträge oder
  illustrative Beispiele ohne Ist-Stand-Anspruch.
- geprüft, ohne Befund: `jnats`-Deprecation-Warnung — real per
  `docker build --no-cache` reproduziert, identisch bereits in
  `examples/kotlin/nats-stream-client/Main.kt:52` vorhanden — kein neues
  Problem dieses Slices.
- geprüft, ohne Befund: Mutation-Testing-Plausibilisierung — die beiden
  vom Implementer benannten, real geprüften Mutationen (Feld-Leerung,
  verschluckte Transport-Exception) sind an ihrer jeweiligen Eingabeseite
  gebunden, keine Bindung ausschließlich an eine Ausgabeseite/einen
  Fake-Rückgabewert.
- geprüft, ohne Befund: Backtick-Parität — alle fünf in diesem Diff
  geänderten Markdown-Dateien real nachgezählt, alle fünf paarig.
- geprüft, ohne Befund: Kommentar-Disziplin (`AGENTS.md` §3.7) außerhalb
  von F-1 — kein Konjunktiv über eine verworfene Alternative, kein
  abwesender Text, keine Chronik in den neuen `.kt`-Produktionsdateien
  selbst (`nats/**`); die Test-Kommentare, die Slice-Kennungen nennen
  („Rot färbende Mutation (real geprüft,
  slice-sdk-kotlin-nats-stream-client-flaeche): …"), sind zulässige
  Testfall-Provenienz — Satzsubjekt ist die Mutation/der Test, nicht der
  Produktionscode-Pfad.
- geprüft, ohne Befund: Traceability — Commit-Betreff `c8c9e3ae` nennt
  `LH-FA-SST-009`/`ADR-0109`, kein `SPEC-*`/`ARC-*` im Betreff.
- geprüft, ohne Befund: `make gates` real gefahren, Exit-Code direkt
  (ungepiped) geprüft, `0`.

## Summary

| Kategorie | Anzahl |
|---|---|
| HIGH | 1 |
| MEDIUM | 1 |
| LOW | 0 |
| INFO | 0 |

**Finding-Klassen dieses Laufs:** „Slice-/Wellen-Chronik in
Produktionscode-Kommentar" · „Arbeit überholt stehenden Träger"

## Verdikt

**Merge-blockierend:** ja — ein HIGH-Finding. Es ist inhaltlich schmal
(eine Kommentar-Umformulierung in `build.gradle.kts`) und berührt keinen
Sicherheits- oder Korrektheitspfad — die NATS-Fläche selbst
(Nachrichtenschema, Subjekt-Formatierung, Auth-Boundary, Fehlerpfad) ist
vollständig, konsistent designt und real durch 61/61 grüne Tests (14 neu,
real gegen den isolierten Elternstand nachgemessen) sowie ein reales
`.jar` belegt. Nach dem Reviewer-Skill wird ein HIGH-Finding nicht wegen
Plausibilität oder geringer Tragweite herabgestuft — die Klasse
„Slice-/Wellen-Chronik in Produktionscode-Kommentar" hat eine etablierte,
mehrfach bestätigte Präzedenz (4/4 historische Trefferquote vor Merge,
zuletzt identisch im analogen C#-NATS-Sibling-Review gefunden).

**Übergabe:** Rückgabe-Pfeil an den Implementer für eine gezielte
Fixrunde — kein Rollen-Widerspruch (keine Implementer-Gegenrede zu
erwarten, beide Findings sind mechanisch nachvollziehbar), daher **keine**
Konflikt-Pfad-Sequenz über den Architect nötig (Modul 8: Konflikt-Pfad
greift bei Rollen-Widerspruch oder 3× gleichem Konflikttyp, hier liegt
keines von beiden vor). Erwarteter Umfang der Fixrunde:

1. F-1: `build.gradle.kts`s zweiten Version-Kommentarabsatz auf eine reine
   Zustandsaussage zurückführen (z. B. „Version `0.2.0` — additive
   Erweiterung um SSE-/NATS-Vollinhalts-Fläche, kein Breaking Change" mit
   `ADR-0109` Festlegung 1 und/oder dem Anker `· seit
   slice-sdk-kotlin-nats-stream-client-flaeche`, ohne die
   „Version-Hebung (X -> Y) trägt …"-Erzählung) — Lösungsentscheidung
   liegt beim Implementer, nicht bei diesem Report.
2. F-2: `sdks/kotlin/Dockerfile:77,104` auf `pgchangefeed-kotlin-0.2.0.jar`
   nachziehen — derselbe Nachzug, der bereits für `harness/README.md`/
   `harness/mk/sdk.mk` in diesem Diff erfolgt ist.

Die Finding-Klassen gehen in die Slice-Closure §7 und von dort in den
Steering-Loop-Zähler. Dieser Report ist ein Lauf-Beleg; er ersetzt keine
Verifikation gegen die volle DoD — das bleibt Verifier-Aufgabe (Modul 11).

**DoD-Checkbox-Nachzug:** entfällt — dieses Verdikt führt zu einer
Fixrunde, die DoD-Zeile „Review durchgeführt" wird regulär bei Schritt 21
des Implementer-Workflows nachgezogen (Reviewer-Skill §DoD-Checkbox-Nachzug
ohne Fixrunde, Grenzfall „nur ohne Fixrunde").
