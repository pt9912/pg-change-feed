# Slice sdk-kotlin-reale2e: Kotlin-SDK — vier Zustellweg-Flächen mit Realserver-Beleg

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-reale2e](welle-sdk-reale2e.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md)
(Client-Bibliotheken — die Fläche existiert, der Beleg-Stand wird stärker),
[`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Live-Streaming —
gRPC/SSE/NATS-Flächen, Boundary: kein Replay im Stream selbst),
[`LH-FA-SST-006`](../../../../spec/lastenheft.md) (HTTP-API),
[`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
§Entscheidung Festlegung 2/Folgepflicht 1 (die etablierte Mechanik-Klasse,
hier auf den Kotlin-Baum gespiegelt), [`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
(Ort `sdks/kotlin/`, Import-Grenze, unverändert gültig), [`ADR-0060`](../../adr/0060-grpc-streaming-mechanismus.md)/[`ADR-0061`](../../adr/0061-http-sse-zusaetzlich-zu-grpc.md)/[`ADR-0100`](../../adr/0100-nats-dritter-vollinhalts-zustellweg.md)
(Server-Verträge der drei Stream-Wege, werden vom SDK benutzt, nicht erweitert).

**Berührte Spec-Stellen:** [`SPEC-018`](../../../../spec/pflichtenheft.md),
[`SPEC-020`](../../../../spec/pflichtenheft.md),
[`SPEC-021`](../../../../spec/pflichtenheft.md),
[`SPEC-024`](../../../../spec/pflichtenheft.md) — gelesen als Draht-Vertrag,
nicht geändert.

**Verantwortlich:** Implementer-Agent, 2026-09-23.

**Autor:** Planner-Agent, Welle-Eröffnung
[welle-sdk-reale2e](welle-sdk-reale2e.md). **Datum:** 2026-09-23.

---

## 1. Ziel und Abgrenzung

**Ziel:** Der Realserver-Beleg für die **vier Kotlin-Zustellweg-Flächen**
des Packages `pgchangefeed-kotlin` (HTTP `SPEC-018`, gRPC `SPEC-020`, SSE
`SPEC-021`, NATS-Vollinhalt `SPEC-024`) — dieselbe Mechanik, die
`slice-sdk-csharp-reale2e` auf dem C#-Baum durchläuft (und beide im Muster
des Python-Vorbilds, real gebaut in
`slice-sdk-python-grpc-client-flaeche`, `done/`): (a) eine additive
Docker-Stufe `integration` in `sdks/kotlin/Dockerfile` (baut auf `build`
auf — derselbe installierte Jar-/Klassenpfad, kein zweiter Build-Pfad);
(b) ein Runner-Skript `tools/harness/run-sdk-kotlin-integration-tests.sh`
(Arbeitsname) — Bring-up im Python-Runner-Muster, Bau/Start der
`integration`-Stufe im selben Docker-Netz wie der Feed-Container, eine
Phase je Fläche; der Kotlin-Prüfling ist die **kompilierte
Client-Assembly** (der Integrationstest importiert
`io.github.pt9912.pgchangefeed` direkt, kein Wegwerf-Duplikat-Client);
(c) ein Make-Target `test-sdk-kotlin-integration` in `harness/mk/sdk.mk`;
(d) Nachzug der `harness/README.md`-Zeile im selben Zug.

**Je Fläche der Beleg (die vier Phasen des Runners)** — dieselbe
Formulierung wie im C#-Slice §1 (Phasen 1–4: gRPC-Stream-Empfang mit
SQL-Gegenprüfung und `Unauthenticated`-Ablehnung; SSE-Empfang mit
SQL-Gegenprüfung und 401; NATS-Vollinhalt-Empfang mit SQL-Gegenprüfung
und Token-Abweisung; HTTP-Rundlauf RegisterConsumer/ListTables über die
SDK-Methoden mit SQL-Gegenprüfung gegen `cdc.consumer` und 401-Ablehnung),
hier jeweils gegen die Kotlin-Client-Klassen (`PgChangeFeedGrpcClient`,
`PgChangeFeedSseClient`, `PgChangeFeedNatsStreamClient`,
`PgChangeFeedHttpClient` aus `io.github.pt9912.pgchangefeed.*`).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Die Python-HTTP-Fläche** — eigener Folge-Slice
  (`slice-sdk-python-http-reale2e`); er erweitert das bestehende
  Python-Runner-Skript, nicht dieses.
- **Der Abdeckungs-Träger-Erstzug** — die Datei
  `docs/user/sdk-e2e-abdeckung.md` existiert bereits real (erstes
  Erzeugnis des C#-Runners, Vorgänger-Slice); dieser Slice **erweitert**
  sie um den Kotlin-Abschnitt (idempotent, Marker-gegrenzt). Der
  `.d-check.yml`-Eintrag existiert bereits (Vorgänger-Slice) und wird
  hier nicht geändert.
- **Version-Bump oder Publish-Workflow-Änderung** — der Realserver-Test
  ist kein Artefakt-Vertrag; die Top-Level-`version` in
  `build.gradle.kts` und `.github/workflows/sdk-kotlin-release.yml`
  bleiben unverändert (Welle-Plan §6).
- **Ein `examples/kotlin/`-Bezug** — die Beispiele bleiben Doku mit
  Bau-Bindung (`SPEC-023`), unberührt.
- **`spec/pflichtenheft.md`-Nachzug** — kein falsch werdender Träger
  (Welle-Plan §6); `SPEC-027`s Deckungs-Aussage bleibt richtig.
- **Eine Aufnahme des Targets in `make gates`** — braucht DB-Zugang/
  Docker/Netz, dieselbe Klasse wie `make test-integration` (Welle-Plan §6).

## 2. Definition of Done

- [x] `LH-FA-SST-009`-Beleg-Stand stärker: alle vier Kotlin-Flächen
      tragen je einen realen Rundlauf gegen eine laufende Server-Instanz
      (§1, Phasen 1–4), jede Phase mit SQL-Gegenprüfung der empfangenen
      `change_id` gegen `cdc.changes` (gRPC/SSE/NATS) bzw. der
      Registrierung gegen `cdc.consumer` (HTTP) und je einem
      Ablehnungs-Beleg ohne gültiges Token (gRPC `Unauthenticated`,
      SSE/HTTP Status 401, NATS-Token-Abweisung durch den Server).
      *(Sensor-Beleg: `make test-sdk-kotlin-integration` EXIT=0 — der
      Lauf trug gRPC change_id=806-1, SSE 814-1, NATS 818-1 (je
      SQL-Gegenprüfung gegen `cdc.changes`) und die Consumer-Registrierung
      über `cdc.consumer`; die Ablehnungs-Belege trugen je Phase real
      gRPC `Unauthenticated`, HTTP `401` und die NATS-Verbindungsablehnung.
      Mutation real gefahren: gültiger Token im gRPC-Reject-Test — der
      Lauf färbte rot („der Ablehnungs-Beleg blieb aus“), Revert,
      Abschlusslauf grün.)*
- [x] **Mechanik** ([`ADR-0110`](../../adr/0110-python-sdk-umfang-erweitert-vollmatrix.md)
      §Entscheidung Festlegung 2, gespiegelt): additive Docker-Stufe
      `integration` in `sdks/kotlin/Dockerfile` (baut auf `build` auf,
      **dieselben Pins wiederverwendet** — `eclipse-temurin:21-jdk@sha256:085e…`
      aus der bestehenden `build`-Stufe (`sdks/kotlin/Dockerfile` Z. 43),
      **keine neuen Pins**), Runner-Skript
      `tools/harness/run-sdk-kotlin-integration-tests.sh` (Bring-up im
      Python-Runner-Muster, `set -euo pipefail`, Cleanup je Ausgang,
      Phase-Auswahl explizit — kein stiller Ausschluss,
      `BEO-PGC/test-runner-stiller-ausschluss`-Disziplin), Make-Target
      `test-sdk-kotlin-integration` in `harness/mk/sdk.mk`. Kein Gate.
      *(Sensor-Beleg: realer Lauf des Targets EXIT=0 — siehe oben; der
      Runner fuhr alle vier Phasen gegen den laufenden Feed-Container, und
      `make sdk-pack-kotlin` blieb grün nach dem Einzug der
      Integrations-Quellmenge — kein stiller Mitlauf in die andere
      Richtung, Risiko-§6-Ausgang 2.)*
- [x] **Träger-Erweiterung im selben Zug:**
      `docs/user/sdk-e2e-abdeckung.md` trägt den Kotlin-Abschnitt (aus
      derselben Messung, die ihn belegt, idempotent vom Runner
      geschrieben — die Datei trägt nur Zeilen real existierender
      Runner-Phasen).
- [x] Kein Import aus `internal/**`/`cmd/**`/`gen/**` dieses Repos in
      `sdks/kotlin/**` — der neue Integrationstest-Quelltext inklusive
      (Import-Zeilen-Prüfung, Muster der bestehenden SDK-Slices).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow ([`AGENTS.md`](../../../../AGENTS.md) §6),
      kein Self-Review (Modul 8). *(Report `review-slice-sdk-kotlin-reale2e`:
      0 HIGH · 1 MEDIUM · 2 LOW · 3 INFO; Fixrunde `c6523009` zog F-1,
      F-2, F-3, F-6 — F-4/F-5 bleiben INFO ohne erwartete Aktion
      (Review-Verdikt); die Verifikation urteilte „DoD erfüllt" mit drei
      nicht-blockierenden Beobachtungen (`e69a77eb`) — der Nachzug bei
      Schritt 21 des Minimal Agent Workflow geschieht in diesem
      Closure-Commit (`BEO-PGC/dod-checkbox-nachzug`-Form wie im
      C#-Vorgänger).)*
- [x] Doku-Update im selben Zug: `harness/README.md` §Werkzeuge bekommt
      die neue `make test-sdk-kotlin-integration`-Zeile, weil das Target
      hier real entsteht ([`AGENTS.md`](../../../../AGENTS.md) §4; die
      Zeile trägt ihren realen Lauf, `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag. *(§7 unten —
      Lerneintrag: eine Spiegelung überträgt nicht die Assertions-
      Semantik, sondern die Draht-Aussage; Anwendungs-Schärfung der
      verkörperten Klasse `BEO-PGC/arbeit-ueberholt-stehenden-traeger`,
      Anker [`AGENTS.md`](../../../../AGENTS.md) §3.13.)*
- [x] Reconciliation-Register — entfällt: keine Reconciliation-Datei in
      diesem Repo.
- [x] Beobachtungs-Register fortgeschrieben — kein Anfall ist ebenfalls
      eine Antwort und wird in §7 notiert. *(Beleg:
      `BEO-PGC/arbeit-ueberholt-stehenden-traeger`,
      `evidence/slice-sdk-kotlin-reale2e.md`; übrige Kandidaten geprüft,
      siehe §7.)*
- [x] Jedes Risiko aus §6 trägt einen Ausgang. *(fünf Ausgänge in §6, je
      entfallen mit Begründung am Ort; siehe §7.)*
- [x] Die drei Paarungen getragen — dieser Slice gehört zu
      [welle-sdk-reale2e](welle-sdk-reale2e.md). *(Anker: die Regel des
      Lerneintrags liegt als Anwendungs-Schärfung der verkörperten Klasse
      `BEO-PGC/arbeit-ueberholt-stehenden-traeger` vor (Anker
      [`AGENTS.md`](../../../../AGENTS.md) §3.13, am Ort existent);
      Folge-Slice: der genannte Plan existiert als Datei in `open/`;
      Register: der neue Beleg liegt in `evidence/` — §7.)*

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `tools/harness/run-sdk-kotlin-integration-tests.sh` | neu | Runner, wortwörtliche Spiegelung des C#-Runner-Musters (dasselbe Bring-up, dieselbe Phase-Form); Form-Vorbilder: `tools/harness/run-sdk-python-integration-tests.sh` (Origin) und der C#-Runner (Spiegel). |
| `sdks/kotlin/pgchangefeed-kotlin/src/integrationTest/kotlin/` (Arbeitsname) | neu | der Integrationstest importiert die Client-Klassen direkt; eigener Quellmenge-Ordner (Gradle-Quellset oder eigenes Modul — der umsetzende Zug entscheidet über die Gradle-Form), damit der Unit-Lauf (`./gradlew test`) unberührt bleibt und die `integration`-Stufe gezielt aufruft (kein stiller Ausschluss des Rests). |
| `sdks/kotlin/Dockerfile` | update | additive Stufe `integration` (baut auf `build` auf); Pins unverändert (`eclipse-temurin:21-jdk@sha256:085e…`, Wiederverwendung aus `build`, kein neuer Pin). |
| `harness/mk/sdk.mk` | update | neues Target `test-sdk-kotlin-integration`. |
| `docs/user/sdk-e2e-abdeckung.md` | update | Kotlin-Abschnitt als Erzeugnis des Kotlin-Runners (idempotent, Marker-gegrenzt). |
| `harness/README.md` §Werkzeuge | update | neue Zeile für `make test-sdk-kotlin-integration` (nach dem realen Lauf geschrieben). |

**Plan-Nachzug (im selben Lauf, vor dem Gate-Lauf):**

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts`: SourceSet `integrationTest` + Task | neu | die Gradle-Form (Plan §3 ließ die Form offen): eigener SourceSet + eigener `integrationTest`-Task, der **nicht** an `check` hängt — der netzlose Pack-Lauf (`make sdk-pack-kotlin`) bleibt unberührt (real gemessen, EXIT=0). Zwei eigene Konfigurationen: `integrationTestImplementation` (erbt `testImplementation`, trägt das main-Output als Prüfling) und `integrationTestRuntimeClasspath` (resolvable, erbt `testRuntimeOnly` **und** `runtimeOnly` — ohne `grpc-netty-shaded` endet der gRPC-Kanal real in der ProviderNotFoundException, im zweiten Lauf gesehen). Der Task-Classpath wird explizit aus diesen Konfigurationen gebaut; `showStandardStreams` reicht die Runner-Marker live durch (JVM-stdout-Pufferung, Risiko-§6-Ausgang 3). |
| `sdks/kotlin/pgchangefeed-kotlin/src/integrationTest/kotlin/` (5 Dateien) | neu | PhaseEnvironment (Env-Auslese + Marker-Druck mit Flush) + vier Testklassen (je Happy-Path + Reject-Beleg, `kotlin.test`/JUnit-Platform-Form wie der Unit-Baum). |
| gson-Semantik: `oldImage == null \|\| isJsonNull` | Erweiterung | gson trägt JSON-null in ein `JsonElement?`-Feld als `JsonNull.INSTANCE`, nicht als Kotlin-null — der Alt-Bild-Assert prüft die Semantik (kein Alt-Bild) über beide Formen; real im vierten Lauf als rote Assertion gesehen und in diesem Lauf gelöst. |

**§3.13-Suchlauf (committetes Feld — bewegte Eigenschaft: „die Kotlin-Flächen tragen reale Realserver-Belege"; beide Stände gemessen: Parent `b58cb173` und HEAD):**

| Träger | Befund | Behandlung |
|---|---|---|
| `harness/README.md` §Werkzeuge | Zeile fehlt (Target real, Zeile nicht) | in diesem Zug ergänzt (nach dem realen Lauf) |
| `docs/user/sdk-e2e-abdeckung.md` | Kotlin-Abschnitt fehlt | in diesem Zug ergänzt (C#-Abschnitt byte-identisch erhalten — Writer-Form-Erhalt auf beiden Seiten real gemessen) |
| `harness/README.md` §Werkzeuge (`make test-sdk-csharp-integration`-Zeile) | Endklause „erweitert sich Slice für Slice um die Kotlin- und Python-HTTP-Abschnitte" — mit diesem Slice zur Hälfte verbraucht | gemeldet, nicht gezogen — die Endklasse bleibt als Verlaufs-Aussage wahr (der Kotlin-Abschnitt ist mit diesem Slice real) |
| `sdks/kotlin/pgchangefeed-kotlin/README.md` §Status (die erste Feld-Fassung nannte fälschlich `sdks/kotlin/README.md` — Adresse ohne Artefakt, Review F-1) | geprüft — trägt keine Teststrategie-Aussage über Realserver-Läufe | kein Nachzug nötig |
| `docs/user/benutzerhandbuch.md` | geprüft — trägt die SDK-Hinweise, keine E2E-Beleg-Aussage | nichts zu ziehen |
| `spec/pflichtenheft.md` | geprüft — kein falsch werdender Träger | nichts zu ziehen |

**Ansatz:** Dieselbe `run_surface_phase`-Form wie im C#-Slice (begrenztes
Fire-and-Forget-Fenster, eindeutige Sentinel- und ID-Wertebereiche je
Phase, Reject-Marker je Protokoll). Der Runner schreibt den Träger-
Abschnitt idempotent aus derselben Messung, die ihn belegt
(`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`-Disziplin).

## 4. Trigger

**Start** (`next` → `in-progress`): wenn diese Welle eröffnet ist und kein
anderer Slice in `in-progress/` liegt (WIP-Limit 1); die Form ist durch
den C#-Vorgänger durchlaufen.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls die
  Gradle-Form der Integrations-Quellmenge (Quellmenge-Ordner vs. eigenes
  Modul vs. `--tests`-Filter) mehr Aufwand verlangt als erwartet — dann
  Rück zur Zerlegung mit der Gradle-Entscheidung als eigenem
  Mechanik-Kern.
- `in-progress` → `open` (blockiert — Carveout?): die `integration`-Stufe
  scheitert an einem Kotlin-/Gradle-spezifischen Verhalten, das der
  Vorbild-Mechanik widerspricht (unwahrscheinlich — die `build`-Stufe
  baut den Jar bereits real, `sdk-pack-kotlin` läuft grün).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + ein realer, grüner
`make test-sdk-kotlin-integration`-Lauf mit allen vier Phasen +
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- **Form-Vorbild-Kopie trägt deutsches Wortfragment weiter**
  (`BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`,
  offen, 2×): Runner-Skript und Dockerfile-Kommentar entstehen als Kopie
  zweier Vorbilder (Python-Origin, C#-Spiegel). *Zusage:* die
  übernommenen Form-Teile werden je separat auf Sprachreinheit geprüft;
  der Suchlauf im Bericht trägt das Ergebnis.
  **Ausgang:** entfallen — kein unübersetztes deutsches Wortfragment in
  einem englischen Klassen-Doc-Kommentar (Zähler der Klasse bleibt 2×);
  die je-teilige Sichtung hat den Treffer der Träger-Familie nicht
  ausgesiebt, gefunden hat ihn der Review (F-2: das Fragment „FlaecheN"
  im Runner-Kopf, wortgleich aus dem C#-Spiegel-Kopf — deutschsprachige
  Kommentar-Familie, ASCII-Transliteration), gezogen in der Fixrunde
  (`c6523009`), die beräumte Form vom Verifier bestätigt (Verifikation
  §3). Struktur-Abgrenzung im Register-Abschnitt §7.
- **Gradle-Quellmenge-Form:** `--tests`-Filter, eigener
  Quellmenge-Ordner (`src/integrationTest/` mit eigener Gradle-SourceSet-
  Registrierung) oder ein eigenes Gradle-Modul — die drei Formen haben
  unterschiedliche Unit-Lauf-Konsequenzen; die falsche Form schließt die
  Integrations-Tests aus dem Unit-Lauf aus (gut) oder lässt sie still
  mitlaufen (schlecht — netzloser `make sdk-pack-kotlin`-Lauf würde
  realserver-Tests verlangen und scheitern). *Erwartet, zu belegen
  durch:* ein realer `make sdk-pack-kotlin`-Lauf bleibt grün, nachdem
  die Integrations-Quellmenge existiert (kein stiller Ausschluss in die
  andere Richtung). **Ausgang:** entfallen — der reale
  `make sdk-pack-kotlin`-Lauf blieb über die Integrations-Quellmenge
  grün (Layer-Kette CACHED über den COPY der Quellmenge bis
  `pack-export`, EXIT=0 — Verifikation §1); die gebaute Form trägt der
  Plan-Nachzug (eigener SourceSet + Task, nicht an `check`; zwei eigene
  Konfigurationen mit Kommentar-Begründungszug).
- **JVM-stdout-Pufferung vs. `docker logs`-Marker-Polling:** der Runner
  liest die Marker über `docker logs`, während der Test läuft;
  JVM-stdout kann zeilenweise gepuffert sein. *Erwartet, zu belegen
  durch:* der erste reale Lauf zeigt die Marker fristnah; falls nicht,
  trägt der Fix die ungepufferte Ausgabe-Form (Gradle-Runner-Option bzw.
  explizite Flushes) und den Beleg.
  **Ausgang:** entfallen — die gebaute Form trägt beide Hälften
  (`System.out.flush()` je Marker in den vier Testklassen,
  `showStandardStreams` live in `build.gradle.kts`, Plan-Nachzug), und
  die Marker trafen in beiden Läufen (Implementer, Verifier — je EXIT=0
  mit vollständigen READY/RECEIVED/REJECTED-Zügen) fristnah.
- **zwei Token-Klassen in der HTTP-Phase** — dieselbe Ausgangslage wie im
  C#-Slice §6 (dritter Punkt); *Erwartet, zu belegen durch:* der reale
  Lauf. **Ausgang:** entfallen — die HTTP-Phase fährt beide Klassen real
  (admin-/reader-Paar je Phase-Env gebunden, Runner; Verifikation §4);
  der Ablehnungs-Beleg (401) lief am token-freien Aufruf, der
  Listen-Aufruf am reader-Token.
- **Sentinel-/ID-Kollisionen mit anderen Läufen** — der Kotlin-Runner
  wählt eigene Sentinel- und ID-Wertebereiche (Muster Python-Runner:
  300/310/320; C#-Runner wählt eigene). *Erwartet, zu belegen durch:* der
  reale Lauf. **Ausgang:** entfallen — eigene Bereiche je Phase
  (440/450/460/470 gegen C# 400/410/420/430 und Python 300/310/320,
  disjunkt; Sentinels sprachspezifisch disjunkt — Review-Negativbefund);
  Implementer- und Verifier-Lauf ohne Kollision, die IDs sind
  laufgebunden (Verifikation §6.2).

## 7. Closure-Notiz

- **Was hat funktioniert:** der Mechanik-Spiegel trug den Zug ohne zweite
  Infrastruktur — die additive `integration`-Stufe baut auf der bestehenden
  `build`-Stufe auf (kein neuer Pin, Verifikation §2 Zeile 2), und der
  Runner fuhr alle vier Phasen gegen dieselbe Compose-Umgebung:
  `make test-sdk-kotlin-integration` lief real grün (Implementer-Lauf
  EXIT=0: gRPC `change_id=806-1`, SSE `814-1`, NATS `818-1`, je
  SQL-Gegenprüfung gegen `cdc.changes`, die Registrierung über
  `cdc.consumer`; Verifier-Lauf EXIT=0 mit laufgebundenen Nachbar-IDs
  `805-1`/`810-1`/`817-1` — die IDs sind laufgebunden, je SQL-Gegenprüfung
  wirksam, Verifikation §1/§2). Die aus der C#-Kette übernommenen Lern-
  ketten trugen real: die Writer-Form-Randbedingung (C#-Review F-7,
  `done/slice-sdk-csharp-reale2e.md` §7 Folge-Slices) steckt im
  Kotlin-Runner als beidseitiger Erhalt (`awk` vor `kotlin-begin`/ab
  `kotlin-end`; der C#-Abschnitt blieb in allen Läufen byte-identisch),
  die SST-009-Zitierpflicht je Träger-Zeile (C#-Review F-2-Klasse) stand
  ab der ersten Kotlin-Zeile (Review-Negativbefund), und der netzlose
  Pack-Lauf blieb über die Integrations-Quellmenge grün (Layer-Ketten-
  Beleg, Verifikation §1). Die Rollen-Kette lief unabhängig: Haupt-Review
  F-1…F-6 (0 HIGH · 1 MEDIUM · 2 LOW · 3 INFO), Fixrunde `c6523009`
  (F-1, F-2, F-3, F-6), Verifikation „DoD erfüllt" mit eigenen
  Sensor-Läufen und drei nicht-blockierenden Beobachtungen ohne erwartete
  Aktion (Verifikation §6, `e69a77eb`).
- **Was ging anders als geplant:** die C#-Lernkette trug die Draht-
  Aussagen der Spiegelung — die neu aufgetretenen Lücken waren Gradle-/
  Bibliotheks-spezifische Formen, die keine Spiegelung mitträgt: die
  Config-Resolvability-Kette (eine reine Vererbung aus `runtimeOnly`/
  `integrationTestImplementation` ohne eigene resolvable Konfiguration
  endet real in der ProviderNotFoundException — `grpc-netty-shaded` fehlt
  im Lauf-Classpath; die gebaute Form trägt zwei eigene Konfigurationen,
  Plan-Nachzug) und die gson-Null-Semantik (gson trägt JSON-null als
  `JsonNull.INSTANCE`, nicht als Kotlin-null — die C#-`Assert.Null`-Form
  wäre auf dem gson-Baum falsch; die gebaute Form prüft die Aussage über
  beide Formen, Plan-Nachzug). Dazu das Suchlauf-Feld selbst: seine
  Prüf-Angaben trugen zwei Baum-Widerlegungen (F-1 — Adresse ohne
  Artefakt, `sdks/kotlin/README.md` existiert nicht, real ist es
  `sdks/kotlin/pgchangefeed-kotlin/README.md`; F-6 — Behandlungs-Angabe
  „gezogen" für eine Zeile, die an beiden Ständen byte-identisch bleibt,
  „gemeldet, nicht gezogen"), beide vom Reviewer gefunden, beide in der
  Fixrunde gezogen, beide vom Verifier gegen beide Stände bestätigt
  (Verifikation §3).
- **Steering-Loop-Eintrag (Lerneintrag):** geschärfte Regel
  (Anwendungs-Schärfung der bereits verkörperten Klasse
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger`, Anker
  [`AGENTS.md`](../../../../AGENTS.md) §3.13 · seit welle-20): **eine
  Spiegelung überträgt nicht die Assertions-Semantik, sondern die
  Draht-Aussage** — die C#-Lernkette lief als Draht-Aussage wortwörtlich
  weiter (Writer-Form beidseitig, SST-009 je Träger-Zeile), während die
  Assertions-Form des Vorbilds je Sprache neu zu binden ist:
  Bibliotheksspezifische Null-Formen (gson `JsonNull` vs. C# `null` vs.
  Python `None`) brauchen je Sprache eine eigene Bindungs-Form; die
  gebaute Form (`oldImage == null || oldImage!!.isJsonNull`, SSE/NATS-
  Test) trägt die Draht-Aussage (kein Alt-Bild am INSERT) über beide
  gson-Formen und bleibt rot-fähig bei real vorhandenem Alt-Bild. Kein
  neuer Sensor: die verfügbare Falsifikation bleibt die Messung an beiden
  Ständen; die Rot-Fähigkeit der Assertion trägt der Quelltext (Review-
  Negativbefund) — ein Wortmuster, das Bibliotheks-Semantik-Übernahmen
  träfe, bräuchte dieselbe Semantik-Entscheidung, die §3.13s Grenze für
  Zahlen/Prosa-Umformulierungen bereits benennt. Keine benannte
  Spec-Lücke.
- **Beobachtungs-Register (`../observations/`):**
  - **`BEO-PGC/arbeit-ueberholt-stehenden-traeger`** — neuer,
    fünfundzwanzigster Beleg: `evidence/slice-sdk-kotlin-reale2e.md`;
    F-1/F-6 tragen die Prüf-Angaben-Widerlegungen am Suchlauf-Feld, die
    Fixrunde zog sie; `state.md` trägt den abgeleiteten Zähler (25×,
    Datei-Anzahl unter `evidence/`, real ausgezählt).
  - **`BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter`** —
    geprüft: F-2 trägt die Klasse nicht, Zähler bleibt 2× — das Fragment
    „FlaecheN" sitzt in der deutschsprachigen Runner-Kommentar-Familie
    (ASCII-Transliteration, `tools/harness/`), nicht im englischen
    Klassen-Doc-Kommentar eines SDK-Sprachpakets (`sdks/*/`); dieselbe
    Struktur, die der C#-Vorgänger für seine F-3/F-4 als notierte Antwort
    führte. Der Befund (wortgleiche Übernahme aus dem C#-Spiegel-Kopf,
    Review F-2) ist in der Fixrunde gelöst und im Report konserviert; der
    C#-Spiegel-Kopf trägt dieselbe Form-Familie weiterhin — gemeldet,
    nicht gezogen: der C#-Baum gehört dem gelösten Vorgänger-Slice, hier
    kein Schreib-Ziel.
  - **F-3 (Writer-Fehlbestands-Pfad ohne Kopf)** — geprüft: keine
    Register-Klasse trägt ihn; degenerater Pfad, in der Fixrunde
    geschlossen (Kopf-Regeneration im Muster des C#-Runners), im Report
    konserviert.
  - **F-6 (Behandlungsspalte)** — trägt den Beleg von
    `BEO-PGC/arbeit-ueberholt-stehenden-traeger` mit (Prüf-Angabe des
    Suchlauf-Felds), kein eigener Eintrag.
  - **F-4/F-5 (Vorbild-Frist, Marker-Ordnung)** — INFO ohne erwartete
    Aktion (Review-Verdikt); keine Klasse trägt sie — keine Beobachtung
    angefallen.
- **Folge-Slices:** `slice-sdk-python-http-reale2e` (in `open/`,
  Welle-Plan §4) — der letzte Matrix-Slice erweitert den bestehenden
  Python-Runner um die HTTP-Phase und trägt die Matrix-Vollständigkeits-
  Prüfung; die Writer-Form-Randbedingung (C#-Review F-7-Kette) gilt für
  ihn unverändert und ist hier real gebaut und gemessen (beidseitiger
  Erhalt, Runner-Kommentar). Kein zusätzlicher Slice aus dieser Closure.
- **Risiken aus §6:** Risiko 1 (Form-Vorbild-Kopie) — **entfallen** mit
  Begründung am Ort (die je-teilige Sichtung hat den Treffer der
  Träger-Familie nicht ausgesiebt, gefunden hat ihn der Review F-2, gezogen
  in der Fixrunde `c6523009`; Zähler der Klasse bleibt 2×); Risiko 2
  (Gradle-Quellmenge-Form) — **entfallen**; Risiko 3 (Marker-Pufferung) —
  **entfallen**; Risiko 4 (zwei Token-Klassen) — **entfallen**; Risiko 5
  (ID-Kollisionen) — **entfallen** (alle Begründungen und Belege am Ort in
  §6). Kein Ausgang „weiter offen" ins Register — ein Registereintrag für
  ein im Slice gelöstes Muster trüge kein `evidence/` (Paarung (c)).
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-reale2e](welle-sdk-reale2e.md) (noch offen) — die Prüfung
  läuft regelkonform bei deren Closure. (a) Anker: die Regel des
  Lerneintrags liegt als Anwendungs-Schärfung der verkörperten Klasse
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` vor (Anker
  [`AGENTS.md`](../../../../AGENTS.md) §3.13, `seit welle-20`, am Ort
  existent); (b) Folge-Slice: der genannte Plan existiert als Datei in
  `open/`; (c) Register: jede genannte Kennung existiert als Verzeichnis,
  und das `evidence/` des genannten Eintrags trägt 25 Dateien (real
  ausgezählt).

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `sdks/kotlin/` — bereits
mit `slice-sdk-kotlin-projektgeruest` eröffnet (GF), keine erneute
Ausdifferenzierung nötig; `tools/harness/` ist Werkzeug-Area der
bestehenden Runner-Familie (GF, fortlaufend).

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/formvorbild-kopie-traegt-deutsches-wortfragment-weiter` (offen,
2×, einschlägig — Risiko §6; dieser Slice ist potenziell das dritte
Auftreten, das der Eintrag als Prüfpunkt benennt),
`BEO-PGC/test-runner-stiller-ausschluss` (offen, 2×, nicht einschlägig —
explizite Phase-Auswahl), `BEO-PGC/arbeit-ueberholt-stehenden-traeger`
(verkörpert, Suchlauf trägt README-Status-Abschnitt und
`harness/README.md`-Zeilen),
`BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` (offen, 3×,
Architect-Entscheidung — READMEs höchstens Suchlauf-Ziel, kein
Schreib-Ziel),
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` (verkörpert),
`BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht` (verkörpert).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Fortsetzung der
SDK-Bäume und der `tools/harness/`-Werkzeug-Familie; die Dockerfile-Stufe
ist eine additive Erweiterung eines GF-Artefakts).