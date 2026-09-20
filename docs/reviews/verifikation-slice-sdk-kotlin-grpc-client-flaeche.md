# Verifikationsbericht: slice-sdk-kotlin-grpc-client-flaeche — 2026-09-20

**Rolle:** Verifier (Modul 11) — „Bauen wir es richtig?" gegen den DoD-Vertrag
(`docs/plan/planning/in-progress/slice-sdk-kotlin-grpc-client-flaeche.md` §2)
und die §6-Risiko-Ausgänge. **Nicht** gegen den Diff als solchen
(Reviewer-Aufgabe, mit
[`review-slice-sdk-kotlin-grpc-client-flaeche.md`](review-slice-sdk-kotlin-grpc-client-flaeche.md)
abgeschlossen, 0 HIGH/MEDIUM/LOW/INFO) und **nicht** gegen realen Bedarf
(Validator, hier nicht ausgelöst).

**Frischer Kontext:** Diese Sitzung hat `harness/README.md`, `AGENTS.md`,
`harness/conventions.md`, den vollständigen Slice-Plan (§1–§8), `ADR-0109`
(Accepted, vollständig), `ADR-0060` (referenziert), `spec/pflichtenheft.md`
`SPEC-020` und den Review-Report gelesen. Keine Implementer-/Reviewer-Behauptung
wurde unbesehen übernommen — jede DoD-Zeile wurde eigenständig gegen den
Code-/Doku-Stand nachgemessen, inklusive vier Maven-Central-/Gradle-Plugin-
Portal-`curl`-Abfragen (fünfte Messung insgesamt nach Implementer, Beispiel-
Vorbild und Reviewer), einem eigenen, ungecachten `docker build`-Lauf
(dritte unabhängige Bau-Bestätigung), einem eigenen `git blame`, einem
eigenen, ungepipten `make gates`-Lauf mit direkter Exit-Code-Prüfung
(`AGENTS.md` §3.9) und einer eigenen Backtick-Paritäts-Zählung.

**Gegenstand:** `slice-sdk-kotlin-grpc-client-flaeche`, zwei Commits auf
`main`:

- `473f3ee8` — Implementer-Commit (Kotlin-gRPC-Client-Fläche:
  `PgChangeFeedGrpcClient`, `GrpcStreamTransport`/`FakeGrpcStreamTransport`,
  zwei Tests, `build.gradle.kts`-/`Dockerfile`-Update, Handbuch-/README-Nachzug).
- `c36a4169` — Reviewer-Commit (Review-Report, 0 Findings, DoD-Checkbox-Nachzug
  „Review durchgeführt" ohne Fixrunde).

Elternstand: `69c9e44b` (`slice-sdk-kotlin-http-client-flaeche` bereits in
`done/`). Der Slice liegt weiterhin in `in-progress/` — erwartungsgemäß, die
Closure (`git mv` nach `done/`) ist nicht Gegenstand dieser Verifikation.

---

## 1. DoD-Zeilen einzeln geprüft (§2 des Slice-Plans)

### 1.1 `StreamChanges`-RPC-Fläche, zehn `SPEC-020`-Felder

**Geprüft:** `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/io/github/pt9912/pgchangefeed/grpc/PgChangeFeedGrpcClient.kt`
gelesen. `PgChangeFeedGrpcClient.streamChanges()` öffnet `ChangeStream/StreamChanges`
über `GeneratedStubTransport`/`ChangeStreamGrpcKt.ChangeStreamCoroutineStub`,
liefert `Flow<Change>` unverändert, Bearer-Token im `authorization`-Metadata-
Eintrag (`AUTHORIZATION_METADATA_ENTRY`, `BEARER_PREFIX`).

`spec/pflichtenheft.md` `SPEC-020` (Zeile 440) real gelesen: die zehn Felder
sind `change_id`, `transaction_id`, `source_table_id`, `sequence`, `operation`,
`old_image`, `new_image`, `schema_version`, `schema`, `table` — exakt die zehn
Felder, die `PgChangeFeedGrpcClientMessageSchemaTest.kt` setzt und Feld für
Feld gegen `received.*` prüft. Keine Abweichung, keine Lücke.

**Ergebnis:** konform.

### 1.2 `build.gradle.kts`-Koordinaten gegen Maven Central/Gradle Plugin Portal

**Eigene, fünfte unabhängige Messung** (nach Beispiel-Vorbild, Implementer,
`ADR-0109`-Kontext-Messung, Reviewer) — real per `curl` gegen
`repo1.maven.org`/`plugins.gradle.org` ausgeführt:

| Koordinate | `build.gradle.kts` | eigene Messung (2026-09-20) |
|---|---|---|
| `org.jetbrains.kotlin:kotlin-gradle-plugin` | `2.4.20` | `latest`/`release` `2.4.20` ✓ |
| `io.grpc:grpc-bom` | `1.84.0` | `latest`/`release` `1.84.0` ✓ |
| `io.grpc:grpc-kotlin-stub`/`protoc-gen-grpc-kotlin` | `1.5.0` | `latest`/`release` zeigen real den Commit-Hash-Metadateneintrag `6f774052d1d6923f8af2e0023886d69949b695ee` — bestätigt exakt die im `build.gradle.kts`-Kommentar behauptete Nicht-Release-Anomalie; `1.5.0` bleibt zuletzt echtes Release |
| `com.google.protobuf:protoc`/`protobuf-java` | `4.36.2` | Maven-Metadaten-`<versions>`-Liste real abgerufen: letzte reguläre Version `4.36.2`, danach nur `21.0-rc-1` (Release-Candidate einer neuen Zählung) — bestätigt die im Kommentar behauptete Drift `4.36.1` → `4.36.2` und die RC-Anomalie |
| `org.jetbrains.kotlinx:kotlinx-coroutines-core` | `1.11.0` | `latest`/`release` `1.11.0` ✓ |
| `com.google.code.gson:gson` | `2.14.0` | `latest`/`release` `2.14.0` ✓ |
| `com.google.protobuf` Gradle-Plugin (Plugin Portal) | `0.10.0` | `latest`/`release` `0.10.0` ✓ |

Alle sieben eigenständig geprüften Koordinaten stimmen exakt mit dem
`build.gradle.kts`-Kommentar überein — keine Abweichung, kein unbelegter
Wert (`AGENTS.md` §3.12).

**Ergebnis:** konform.

### 1.3 `sdks/kotlin/Dockerfile` trägt den `proto`-Bau-Kontext korrekt

**Geprüft:** `sdks/kotlin/Dockerfile` gelesen — `WORKDIR /src` bis Zeile 36,
Wechsel auf `WORKDIR /src/pgchangefeed-kotlin` (Zeile 48), danach
`COPY --from=proto cdc/stream/v1/changestream.proto src/main/proto/changestream.proto`
(Zeile 60) — der Zielpfad ist relativ zum zu diesem Zeitpunkt bereits
gewechselten `WORKDIR`, kein doppelter Pfadanteil
(`pgchangefeed-kotlin/pgchangefeed-kotlin/…`, der im Implementer-Bericht §7 als
vor dem ersten Commit korrigierter Fehler beschrieben ist).

Real bestätigt über einen eigenen, ungecachten Docker-Bau (§2 unten):
`generateProto` lief real (nicht `NO-SOURCE`), `test`/`build` beide
`BUILD SUCCESSFUL`.

**Ergebnis:** konform.

### 1.4 Tests: Nachrichtenschema-Vollständigkeit + Authn-Boundary, netzlos

**Geprüft:** `PgChangeFeedGrpcClientMessageSchemaTest.kt` (Feld-für-Feld-
Assertion aller zehn `SPEC-020`-Felder gegen `FakeGrpcStreamTransport`) und
`PgChangeFeedGrpcClientAuthBoundaryTest.kt` (Bearer-Token-Metadata-Assertion +
`assertFailsWith<StatusException>` mit `Status.Code.UNAUTHENTICATED` bei
`FakeGrpcStreamTransport.withStatus(...)`) gelesen. Beide Tests laufen gegen
ein reines In-Memory-Fake (`FakeGrpcStreamTransport`, kein `io.grpc.Channel`,
kein Socket) — netzlos, `AGENTS.md` §3.1-konform.

Die im Implementer-Bericht behaupteten rot-färbenden Mutationen
(`.catch {}` verschluckt `StatusException`; `.map { clearSchema() }` entfernt
das `schema`-Feld) real anhand der Testlogik nachvollzogen: Beide Mutationen
würden je eine konkrete, benannte Assertion sichtbar rot färben — dieselbe
Prüftiefe wie der Reviewer, kein eigener Mutations-Docker-Lauf zusätzlich
gefahren (Prüfumfang bewusst gedeckt durch Testlogik-Nachvollzug statt
Doppel-Ausführung).

**Ergebnis:** konform.

### 1.5 Kein Import aus `internal/**`/`cmd/**`/`gen/**`

**Eigener `grep`-Lauf:**

```
grep -rn "internal/\|cmd/\|gen/" sdks/kotlin/
```

Liefert genau zwei Treffer:

1. `sdks/kotlin/pgchangefeed-kotlin/gradlew:60` — Upstream-Kommentarzeile
   (`org/gradle/api/internal/plugins/…`), kein Repo-Pfad, kein Import.
2. `PgChangeFeedGrpcClientMessageSchemaTest.kt:17` — KDoc-Zitat des
   Test-Vorbilds `internal/adapters/driving/grpc/server_test.go`, kein Import.

**Vorbestehend-Bewertung des `gradlew`-Treffers eigenständig geprüft:**
`git blame -L 58,62 473f3ee8 -- sdks/kotlin/pgchangefeed-kotlin/gradlew` zeigt
Commit `7500fefb8` (2026-09-20 05:01:19, „feat(sdk): Kotlin-SDK-Projektgerüst",
`slice-sdk-kotlin-projektgeruest`) als Ursprung der Zeile — real vor diesem
Slice (Implementer-Commit `473f3ee8`, 06:57:19) entstanden, unverändert im
Diff dieses Slice. Die dritte Bewertung (Implementer, Reviewer, jetzt Verifier)
kommt zum selben Ergebnis: harmlos, vorbestehend, kein neuer Import.

**Ergebnis:** konform.

### 1.6 `make gates` grün

Siehe §2 unten — eigener, ungepipter Lauf, Exit `0`.

**Ergebnis:** konform.

### 1.7 Review durchgeführt, kein offenes HIGH

**Geprüft:** `docs/reviews/review-slice-sdk-kotlin-grpc-client-flaeche.md`
vollständig gelesen. Summary-Tabelle: 0 HIGH / 0 MEDIUM / 0 LOW / 0 INFO.
Verdikt: „Merge-blockierend: nein". Der Report dokumentiert eigenständig
durchgeführte Prüfungen (nicht nur Commit-Message übernommen), inklusive
eines eigenen `docker build --no-cache` + `javap -p`-Laufs gegen die
Kotlin-`internal`-Semantik — dieselbe Lektion, die beim vorigen Slice als
HIGH F-1 gefunden wurde, wurde hier korrekt von Anfang an angewendet
(real durch den Reviewer per Bytecode-Inspektion verifiziert, nicht nur die
KDoc-Aussage übernommen).

**Ergebnis:** konform.

### 1.8 Doku-Update `docs/user/benutzerhandbuch.md`

**Geprüft:** `git show 473f3ee8 -- docs/user/benutzerhandbuch.md` real
gelesen. Bestätigt:

- Version `1.36` → `1.37` (Kopfzeile).
- Neuer `**SDK:**`-Absatz im gRPC-Abschnitt („Zugriff über den
  gRPC-Change-Stream"): nennt `PgChangeFeedGrpcClient.streamChanges()`,
  die Koordinate `io.github.pt9912:pgchangefeed-kotlin`, alle zehn Felder,
  den `authorization`-Metadata-Eintrag, `Unauthenticated`-Boundary — **und**
  den PAT-Hinweis (`read:packages`-Scope) für den Bezug über GitHub Packages
  (`ADR-0109` Festlegung 2).
- Die vorbestehende, veraltete Aussage im HTTP-`**SDK:**`-Absatz
  („gRPC-Change-Stream folgt in einem Folge-Release") ist korrigiert
  (Träger-Nachzug, `AGENTS.md` §3.13) — verweist jetzt auf den neuen
  gRPC-Absatz.
- Versionshistorie-Tabelle trägt eine neue `1.37`-Zeile mit Beleg-Ankern
  (`LH-FA-SST-009`, `ADR-0109`, `slice-sdk-kotlin-grpc-client-flaeche`).

Zusätzlich `sdks/kotlin/pgchangefeed-kotlin/README.md` (Diff, 2 Zeilen)
gegengelesen: dieselbe Korrektur des „follows in a subsequent release"-Status
im englischen Träger.

**Ergebnis:** konform.

---

## 2. Eigener Docker-Build (dritte unabhängige Bau-Bestätigung)

```
docker build --no-cache --build-context proto=proto -f sdks/kotlin/Dockerfile sdks/kotlin
```

Real ausgeführt (kein Cache-Hit — vollständiger Neubau). Ergebnis:

- `> Task :generateProto` (Zeile 73 des Bau-Logs, **nicht** `NO-SOURCE`) im
  `test`-Schritt — die `.proto` wird real über den Zusatzkontext gefunden und
  kompiliert.
- `compileKotlin`/`compileJava`/`test`: alle grün, `BUILD SUCCESSFUL in 49s`
  (Gradle-`test`-Ziel).
- `./gradlew --no-daemon build`: `generateProto UP-TO-DATE` (derselbe Task
  bereits im vorangehenden `test`-Lauf desselben Bau-Layers ausgeführt),
  `BUILD SUCCESSFUL in 11s`.
- Docker-Exit-Code `0`, direkt geprüft (kein Pipe dazwischen).

Damit ist der Doppelpfad-Fehler, den der Implementer-Bericht (§6/§7) als vor
dem ersten Commit korrigiert beschreibt, im committeten Zustand **nicht**
reproduzierbar — dritte unabhängige Bestätigung (Implementer, Reviewer,
jetzt Verifier), alle drei mit demselben Ergebnis. Test-Image nach Prüfung
entfernt (`docker rmi kotlin-grpc-verifier-test`).

---

## 3. `git diff` gegen Elternstand — Berührungs-Scope

```
git diff 69c9e44b..473f3ee8 --stat
```

Neun Dateien geändert, alle innerhalb des in §3 des Slice-Plans geplanten
Umfangs:

- `docs/plan/planning/in-progress/slice-sdk-kotlin-grpc-client-flaeche.md`
- `docs/user/benutzerhandbuch.md`
- `sdks/kotlin/Dockerfile`
- `sdks/kotlin/pgchangefeed-kotlin/README.md`
- `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts`
- `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/…/grpc/PgChangeFeedGrpcClient.kt`
- `sdks/kotlin/pgchangefeed-kotlin/src/test/kotlin/…/grpc/FakeGrpcStreamTransport.kt`
- `sdks/kotlin/pgchangefeed-kotlin/src/test/kotlin/…/grpc/PgChangeFeedGrpcClientAuthBoundaryTest.kt`
- `sdks/kotlin/pgchangefeed-kotlin/src/test/kotlin/…/grpc/PgChangeFeedGrpcClientMessageSchemaTest.kt`

**Bestätigt unberührt:** `examples/kotlin/**`, `.a-check.yml`,
`spec/architecture.md`, `docs/user/version.md`, `sdks/csharp/**`,
`sdks/python/**`, `gen/**`, `proto/**` — keiner dieser Pfade taucht im
`--stat`-Output auf.

**Ergebnis:** konform, keine unbeabsichtigte Berührung.

---

## 4. Backtick-Parität — eigenständig nachgezählt

Alle neun in diesem Slice geänderten/neuen Dateien plus den aktuellen
(nach dem Reviewer-Checkbox-Nachzug-Commit `c36a4169`) Stand des Slice-Plans:

| Datei | Anzahl | Parität |
|---|---|---|
| `docs/plan/planning/in-progress/slice-sdk-kotlin-grpc-client-flaeche.md` | 512 | gerade |
| `docs/user/benutzerhandbuch.md` | 1786 | gerade |
| `sdks/kotlin/pgchangefeed-kotlin/README.md` | 70 | gerade |
| `sdks/kotlin/Dockerfile` | 36 | gerade |
| `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` | 98 | gerade |
| `PgChangeFeedGrpcClient.kt` | 98 | gerade |
| `FakeGrpcStreamTransport.kt` | 16 | gerade |
| `PgChangeFeedGrpcClientAuthBoundaryTest.kt` | 20 | gerade |
| `PgChangeFeedGrpcClientMessageSchemaTest.kt` | 18 | gerade |

Alle neun gerade. (Der Slice-Plan zeigt 512 statt der vom Reviewer zum
Zeitpunkt seines Laufs — vor dem eigenen Checkbox-Nachzug-Commit — gezählten
510; die Differenz von zwei stammt aus dem Checkbox-Nachzug-Commit
`c36a4169` selbst, der Backtick-Inline-Code an einer DoD-Zeile ergänzt. Auch
dieser aktuellste Stand bleibt gerade.)

**Ergebnis:** konform.

---

## 5. §6-Risiken-Ausgänge — sachlich geprüft

- **Risiko 1** (kein aufrufendes `make`-Ziel bis Pack-Werkzeug-Slice).
  **Ausgang „eingetreten, akzeptiert"** — sachlich korrekt: `grep -rn
  "sdk-pack-kotlin\|sdk-kotlin" harness/mk/ Makefile` liefert keinen Treffer;
  `ADR-0109` §Konsequenzen Folgepflicht 1 sieht das Pack-Werkzeug ausdrücklich
  als eigenen, künftigen Slice vor. Der direkte `docker build
  --build-context proto=proto …`-Aufruf (§2 oben, real durchgeführt) trägt die
  DoD-Prüfung dieses Zwischenstands tatsächlich.
- **Risiko 2** (Copy-Zielpfad-Fehlerrisiko der Kotlin-/Gradle-Werkzeugkette).
  **Ausgang „eingetreten und aufgelöst"** — sachlich korrekt: der reale, real
  wiederholte Docker-Bau (Implementer, Reviewer, Verifier — dreifach, alle
  grün) bestätigt, dass `compileKotlin`/`compileJava` ohne `cannot find
  symbol` laufen und `generateProto` real läuft, nicht `NO-SOURCE`.
- **Risiko 3** (Fake/Stub bildet die reale `Unauthenticated`-Ablehnung nur in
  der Form nach, nicht am realen Server). **Ausgang „weiter offen"** —
  sachlich korrekt: `FakeGrpcStreamTransport.withStatus(...)` wirft einen
  `StatusException` direkt aus dem Flow, ohne einen echten
  `io.grpc.Channel`/Interceptor-Pfad zu durchlaufen; ein realer Rundlauf-Beleg
  existiert für dieses SDK bislang nicht (`make test-integration`s
  `tools/harness/grpcclient` ist ein Wegwerf-Client ohne SDK-Import — reale
  Rundlauf-Deckung des generierten Coroutine-Stubs selbst bleibt aus). Die
  Checkbox „Jedes Risiko aus §6 trägt einen Ausgang" ist zu Recht gesetzt —
  der Ausgang ist ehrlich als weiterhin offen deklariert, nicht stillschweigend
  als erledigt verbucht.

**Ergebnis:** alle drei Risiko-Ausgänge sachlich korrekt vorgetragen.

---

## 6. Was nicht Gegenstand dieser Verifikation war (Modul-11-Grenze)

- Kein Fix, keine Code-Änderung.
- Kein `git mv` nach `done/`.
- Keine Closure-Notiz-Fertigstellung (§7 des Slice-Plans bleibt Implementer-
  Entwurf, wie vom Auftrag vorgesehen).
- Reconciliation-Register/Beobachtungs-Register: die Slice-Plan-Aussagen
  „keine Reconciliation-Datei in diesem Repo" und „keine neue Beobachtung"
  wurden gegengeprüft (`ls docs/plan/planning/ | grep -i reconcil` → leer;
  `docs/plan/planning/observations/BEO-PGC` bereits bestehend, kein neues
  Verzeichnis nötig) — beide Aussagen sachlich korrekt.

---

## Verdikt

**DoD konform: ja.** Alle acht geprüften DoD-Zeilen aus §2 des Slice-Plans
sind gegen den realen Code-/Doku-Stand nachgemessen und bestätigt — keine
Behauptung wurde unbesehen übernommen. Der Review-Report ist sachlich
korrekt (0 HIGH/MEDIUM/LOW/INFO, Kotlin-`internal`-Lektion aus dem vorigen
Slice diesmal von Anfang an korrekt angewendet, real durch Bytecode-Inspektion
bestätigt). Alle drei §6-Risiko-Ausgänge sind ehrlich und sachlich korrekt
vorgetragen (zwei aufgelöst, einer bewusst weiter offen). Die Diff-Fläche
bleibt exakt im geplanten Umfang, keine unbeabsichtigte Berührung von
`examples/kotlin/**`, `.a-check.yml`, `spec/architecture.md`,
`docs/user/version.md`, `sdks/csharp/**`, `sdks/python/**`, `gen/**` oder
`proto/**`. Backtick-Parität ist in allen neun betroffenen Dateien gegeben.

Dieser Slice ist bereit für den Rollenwechsel zur Closure (`git mv` nach
`done/`, Planner-/Koordinator-Zug) — nicht Teil dieses Berichts.

**Real durchgeführte Prüfungen dieser Sitzung (Zusammenfassung):**

1. Vollständige Lektüre: `harness/README.md`, `AGENTS.md`,
   `harness/conventions.md`, Slice-Plan (vollständig), `ADR-0109`
   (vollständig), `spec/pflichtenheft.md` `SPEC-020`, Review-Report
   (vollständig).
2. `PgChangeFeedGrpcClient.kt`, `FakeGrpcStreamTransport.kt`,
   `PgChangeFeedGrpcClientMessageSchemaTest.kt`,
   `PgChangeFeedGrpcClientAuthBoundaryTest.kt`, `Dockerfile`,
   `build.gradle.kts` gelesen und gegen `SPEC-020`/DoD gehalten.
3. Sieben Maven-Central-/Gradle-Plugin-Portal-`curl`-Abfragen — alle sieben
   Koordinaten exakt bestätigt, inklusive der Commit-Hash- und
   Release-Candidate-Anomalien.
4. Eigener, ungecachter `docker build --no-cache --build-context
   proto=proto -f sdks/kotlin/Dockerfile sdks/kotlin` — Exit `0`,
   `generateProto` real gelaufen, `test`/`build` beide `BUILD SUCCESSFUL`.
5. `grep -rn "internal/\|cmd/\|gen/" sdks/kotlin/` — genau zwei
   Nicht-Import-Treffer, einer per `git blame` als vorbestehend (Commit
   `7500fefb8`) bestätigt.
6. `git diff 69c9e44b..473f3ee8 --stat` — neun Dateien, alle im geplanten
   Umfang, keine Fremdberührung.
7. Backtick-Parität in neun Dateien eigenständig nachgezählt — alle gerade.
8. `git show 473f3ee8 -- docs/user/benutzerhandbuch.md
   sdks/kotlin/pgchangefeed-kotlin/README.md` — PAT-Hinweis und
   Versionshistorie-Nachzug real bestätigt.
9. `git show c36a4169 --stat` — DoD-Checkbox-Nachzug-Commit real bestätigt.
10. `make gates` — eigener, ungepipter Lauf, Exit-Code direkt geprüft: `0`.

**Finaler `make gates`-Exit-Code (diese Sitzung):** `0`.
