# Slice sdk-kotlin-http-client-flaeche: Öffentliche HTTP-API-Client-Fläche (`SPEC-018`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-kotlin-lh-fa-sst-009](../welle-sdk-kotlin-lh-fa-sst-009.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md),
[`LH-FA-SST-006`](../../../../spec/lastenheft.md) (die neun Port-gedeckten
Fähigkeiten, die diese Fläche als Consumer anspricht),
[`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
Festlegung 1 (Umfang), [`ADR-0057`](../../adr/0057-http-grpc-api.md)
(HTTP/JSON-API-Vertrag, wird vom SDK benutzt, nicht erweitert).

**Berührte Spec-Stellen:** [`SPEC-018`](../../../../spec/pflichtenheft.md)
(Endpunkte, Token-Header-Form — das SDK benutzt diese Festlegungen,
verändert sie nicht).

**Verantwortlich:** Implementer-Agent, 2026-09-20.

**Autor:** Planner-Agent, direkt beauftragt (`ADR-0109` §Konsequenzen
Folgepflicht 1). **Datum:** 2026-09-20.

---

## 1. Ziel und Abgrenzung

**Ziel:** Eine öffentliche, stabile Kotlin-API-Fläche für alle neun
Port-gedeckten Fähigkeiten von [`SPEC-018`](../../../../spec/pflichtenheft.md)
(`RegisterConsumer`, `AcknowledgeConsumer`, `GetConsumerPosition`,
`RemoveConsumer`, `EnableTable`, `DisableTable`, `GetStatus`,
`ListTables`, `RunRetention`) sowie das Changes-Lesen
(`GET /changes`, [`SPEC-022`](../../../../spec/pflichtenheft.md),
[`ADR-0081`](../../adr/0081-changes-lesen-ueber-die-http-api.md)) — Bearer-
Token-Auth (`reader`/`admin`), eigene, von den Beispiel-Clients
unabhängige Tests.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **CLI-Argument-Parsing wie `examples/kotlin/http-client`** —
  `ADR-0109` Festlegung 3 verlangt eine öffentliche, stabile API statt
  eines Kommandozeilen-Programms; ein SDK-Consumer ruft Methoden auf,
  parst keine `argv`.
- **gRPC-Stream** — `slice-sdk-kotlin-grpc-client-flaeche` übernimmt das;
  getrennter Draht-Vertrag (`SPEC-020`), getrennte Fremdabhängigkeit
  (`ADR-0109` Festlegung 1 Abhängigkeits-Footprint-Begründung).
- **SSE- (`SPEC-021`) oder NATS-Vollinhalts-Stream (`SPEC-024`)** —
  `ADR-0109` Festlegung 1 grenzt v1 ausdrücklich auf HTTP-API und
  gRPC-Stream ein (Welle-Plan §6 Out-of-Scope), obwohl für Kotlin bereits
  funktionierende Beispiele für beide existieren.
- **Diagnose/Health-Endpunkte** — `SPEC-018` grenzt sie ausdrücklich aus
  den neun Port-gedeckten Fähigkeiten aus; kein SDK-Umfang.
- **Änderung von `SPEC-018` selbst** — das SDK benutzt den bestehenden
  Draht-Vertrag, verlangt keine Vertragsänderung (`ADR-0109` §Kontext
  Bindung, Analogie zu `SPEC-023`s „Verhältnis zum Draht").

## 2. Definition of Done

- [x] `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/…/http/` (oder
      gleichwertiger Namensraum) trägt eine öffentliche Client-Klasse mit
      einer Methode je der neun Port-gedeckten Fähigkeiten von
      [`SPEC-018`](../../../../spec/pflichtenheft.md) plus dem Changes-Lesen
      (`SPEC-022`) — Signatur, Request-/Response-Form und Fehler-Antwortform
      (`400`/`401`/`403`/`404`/`500` → typisierte Exception oder
      Result-Form, konsistent über alle Methoden) spiegeln
      `SPEC-018`/`SPEC-022` exakt. Bearer-Token wird bei Konstruktion
      übergeben (kein globaler State). Draht-Kenntnis kommt real aus
      `examples/kotlin/http-client/` **und** `spec/pflichtenheft.md` — hier
      existiert (anders als bei Python) ein reales Kotlin-Referenzprogramm
      als Vorbild, das genutzt wird, nicht nur als Struktur-Vorlage
      (`ADR-0109` §Kontext).
- [x] Eigene Tests decken je Fähigkeit mindestens den Happy Path und die
      Auth-Boundary (`401` fehlendes/unbekanntes Token, `403`
      `reader`-Token gegen einen `admin`-Endpunkt) ab — netzlos prüfbar
      (kein realer Server nötig; ein Fake/Mock des HTTP-Transports analog
      dem bestehenden Test-Muster der Beispiel-Clients bzw. der beiden
      vorigen SDKs, konkrete Form entscheidet der Implementer-Zug anhand
      des tatsächlich gewählten Kotlin-HTTP-Clients — `examples/kotlin/http-client`
      nutzt `java.net.http`, real nachzuprüfen).
- [x] Kein Import aus `internal/**`/`cmd/**` dieses Repos (`ADR-0109`
      §Kontext Bindung „Import-Grenze, hier ohne Ausnahme") — real geprüft:
      `grep -rn "internal/\|cmd/" sdks/kotlin/` liefert keinen Treffer.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Fixrunde geprüft und Merge-Block aufgehoben, falls nötig.
      `docs/reviews/review-slice-sdk-kotlin-http-client-flaeche.md` fand
      1 HIGH (F-1, Kotlin-`internal` fälschlich als JVM-Zugriffsschutz
      behauptet in `HttpTransport.kt`s KDoc und im Plan-Nachzug oben) — in
      dieser Fixrunde in beiden Trägern korrigiert, kein weiteres Vorkommen
      im Diff (`grep -rn "invisible outside\|no public API surface\|nur
      innerhalb des Gradle-Moduls" sdks/kotlin/` ohne Treffer); kein offenes
      HIGH mehr.
- [x] Doku-Update: `docs/user/benutzerhandbuch.md` bekommt einen
      SDK-Hinweis für die Kotlin-HTTP-Oberfläche (`ADR-0109` §Konsequenzen
      Folgepflicht 4, **inklusive** des expliziten PAT-Hinweises für den
      Bezug über GitHub Packages, Festlegung 2) — getragen durch die
      bereits verkörperte Selbstprüf-Instruktion
      (`.claude/commands/implement-slice.md` Schritt 17) und den
      Reviewer-HIGH-Punkt
      (`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`).
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`; keine Beobachtung angefallen
      ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [x] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-kotlin-lh-fa-sst-009](../welle-sdk-kotlin-lh-fa-sst-009.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/…/http/PgChangeFeedHttpClient.kt` (Arbeitsname) | neu | öffentliche API-Fläche für die neun Port-gedeckten Fähigkeiten + Changes-Lesen. |
| `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/…/http/model/*.kt` | neu | typisierte Request-/Response-Datenklassen (spiegeln `SPEC-018`/`SPEC-022` JSON-Schemas). |
| `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/…/http/PgChangeFeedException.kt` (Arbeitsname) | neu | typisierte Fehlerform für `400`/`401`/`403`/`404`/`500`. |
| `sdks/kotlin/pgchangefeed-kotlin/src/test/kotlin/…/http/*Test.kt` | neu | Happy Path je Fähigkeit (`SPEC-018`/`022` Referenz), Auth-Boundary `401`/`403`. |
| `docs/user/benutzerhandbuch.md` | update | SDK-Hinweis für die Kotlin-HTTP-Oberfläche inkl. PAT-Hinweis, im selben Zug (`ADR-0109` Folgepflicht 4). |

**Ansatz:** Referenzmaterial ist `examples/kotlin/http-client/` (Draht-
Kenntnis **und** reales Kotlin-Vorbild — anders als bei den C#-/Python-Wellen
existiert hier ein bereits funktionierendes Kotlin-Programm gegen dieselbe
API, `ADR-0109` §Kontext „Was das ändert"), kein `project(":…")`-Abhängigkeitspfad
darauf (`ADR-0109` Festlegung 3: „dieselbe Draht-Kenntnis, aber als
eigenständiger, paketierbarer Code neu geschrieben").

**Hinweise aus dem Beobachtungs-Register (proaktiv, nicht erst nach einem
Reviewer-Finding):**

- `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (Zähler real 19×,
  `AGENTS.md` §3.13): Vor Abschluss dieses Slice einen Träger-Nachzug-
  Suchlauf über `sdks/kotlin/pgchangefeed-kotlin/README.md`,
  `build.gradle.kts` und jede weitere Datei fahren, die eine Aussage über
  den Lieferstand dieser Fläche trägt (z. B. „HTTP folgt erst") — sowohl
  deutsche als auch englische Formulierungen (die Python-Welle fand real,
  dass ein aus einem vorigen Vorkommen übernommener Suchlauf beide Achsen
  unabhängig breit genug wählen muss, sonst bleibt eine Datei oder eine
  Formulierung unentdeckt).
- **Backtick-Paritäts-Check** vor jedem Commit dieses Slice.

**Plan-Nachzug (Implementer-Zug, `AGENTS.md`-Konvention „im selben Lauf
nachtragen"):**

- **Zusätzliche Dateien gegenüber der Tabelle oben** (Fähigkeits-/
  Verantwortungs-Zerlegung statt einer einzigen Modell-/Testdatei, kein
  Umfangs-Wachstum): `http/HttpTransport.kt` (Transport-Abstraktion, siehe
  §6 Risiko 3), `http/model/{Consumers,Tables,Retention,Changes,ErrorResponse}.kt`
  (fünf Dateien statt einer, analog dem bestehenden Datei-Schnitt der
  C#/Python-Geschwister), fünf Testdateien statt einer
  (`PgChangeFeedHttpClientConsumerTest.kt`,
  `PgChangeFeedHttpClientTableTest.kt`,
  `PgChangeFeedHttpClientRetentionAndChangesTest.kt`,
  `PgChangeFeedHttpClientAuthBoundaryTest.kt`) plus zwei Test-Helfer
  (`FakeHttpTransport.kt`, `TestClientFactory.kt`) — derselbe Datei-Schnitt
  wie `sdks/csharp/PgChangeFeed.Client.Tests/Http/`.
- **§6 Risiko 1 (Fehler-Antwortform) entschieden:** `PgChangeFeedException`
  ist eine Kotlin **sealed class** mit sieben konkreten Unterklassen (Namen
  identisch zur C#-Fassung: `PgChangeFeedBadRequestException` usw.) — Kotlin-
  idiomatischer als eine offene Exception-Hierarchie (exhaustives `when`
  möglich), bei identischem Feld-für-Feld-Verhalten zu C#/Python. Siehe
  KDoc in `PgChangeFeedException.kt` für die vollständige Begründung.
- **§6 Risiko 3 (HttpClient-Testbarkeit) real geprüft, nicht nur behauptet:**
  `java.net.http.HttpClient` hat — anders als C#s `HttpMessageHandler` oder
  Pythons `httpx.MockTransport` — keinen Pluggable-Handler; `HttpClient.send()`
  öffnet real einen Socket. Lösung: eine `internal fun interface HttpTransport`
  zwischen `PgChangeFeedHttpClient` und dem JDK-Client — der öffentliche
  Konstruktor nimmt weiterhin ein `java.net.http.HttpClient` (Parität zu
  C#/Python: Aufrufer besitzt/kontrolliert den Client), ein zweiter,
  **`internal`** Konstruktor nimmt direkt einen `HttpTransport` (das
  Kotlin-Gradle-Plugin bindet `test` per Default an `main`s
  `internal`-Sichtbarkeit) — jeder Test in diesem Slice ist dadurch **echt
  netzlos**, ohne Socket, ohne Loopback-Server. `internal` ist dabei eine
  **compile-time**-Sichtbarkeitsgrenze des Kotlin-Compilers gegenüber
  anderen Kotlin-Modulen, **keine** JVM-Bytecode-Zugriffsbeschränkung: Im
  kompilierten Jar sind sowohl dieser Konstruktor als auch `HttpTransport`
  selbst gewöhnliche `public`-Symbole (real mit `javap -p` gegen das gebaute
  Jar geprüft) — ein Java-Konsument oder Reflection kann beide trotzdem
  erreichen. Die Grenze schützt gegen versehentliche Nutzung aus anderen
  Kotlin/Gradle-Modulen, nicht gegen jeden JVM-Aufrufer.
- **Neue Fremdabhängigkeit `com.google.code.gson:gson:2.14.0`** — real
  gegen Maven Central nachgemessen (`maven-metadata.xml`, 2026-09-20,
  weiterhin aktuellste Version), dieselbe bereits im selben Repo bewertete
  Version wie `examples/kotlin/nats-stream-client/build.gradle.kts`
  (2026-09-18 dort gemessen, keine Drift) — keine neue, unabhängig zu
  bewertende Bibliothek.
- **`offset` ([`SPEC-018`](../../../../spec/pflichtenheft.md) `uint64`) als `Long` modelliert, nicht `ULong`:**
  Gsons reflektionsbasierter Codec unterstützt Kotlins `ULong`
  (Inline-/Value-Class) nicht korrekt — er würde das interne gewrappte
  `Long`-Feld (de-)serialisieren statt des Werts selbst und dabei die
  JSON-Form verfälschen. `Long` passt zu den übrigen 64-bit-Feldern dieses
  SDK und vermeidet dieses Gson-Verhalten, auf Kosten der oberen Hälfte des
  `uint64`-Wertebereichs — eine benannte, schmale Grenze für dieses
  Pre-1.0-Release (siehe KDoc in `model/Consumers.kt`).
- **Träger-Nachzug-Suchlauf real durchgeführt** (Hinweis oben): fand **eine**
  stehende Aussage in `sdks/kotlin/pgchangefeed-kotlin/README.md` §Status
  („a full HTTP API client surface … follow in subsequent releases") — durch
  diesen Slice falsch geworden, im selben Zug korrigiert (jetzt: „The
  current release provides … a full HTTP API client surface …"). Kein
  Treffer in `build.gradle.kts`, `Dockerfile` oder sonstwo unter `sdks/kotlin/`.
  Die bereits bekannte, separat verfolgte LOW-Klasse
  `BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` („vollinhalt", schon
  bei 3× Schwelle) bleibt bewusst unverändert — ihre Auflösung ist laut
  Closure-Notiz von `slice-sdk-kotlin-projektgeruest` eine Architect-Aufgabe
  (Pflege/Steering-Loop), kein Ad-hoc-Fix eines einzelnen Implementer-Zuges.
- **Mutation real gesehen (nicht nur behauptet, `AGENTS.md` §3.7/Schritt 19
  des Implementer-Workflows):** zwei Produktionscode-Mutationen real
  gefahren und wieder zurückgenommen — (1) `encode()` auf Identität gesetzt
  (keine Prozent-Kodierung mehr) macht
  `listTables percent-encodes reserved characters in query parameters`
  (`PgChangeFeedHttpClientTableTest.kt`) real rot; (2) die `403`-Zuordnung in
  `buildException` auf `PgChangeFeedBadRequestException` vertauscht macht
  `reader token against an admin endpoint throws Forbidden`
  (`PgChangeFeedHttpClientAuthBoundaryTest.kt`) real rot. Beide Male über
  `docker build --no-cache -f sdks/kotlin/Dockerfile sdks/kotlin` beobachtet
  (`BUILD FAILED`/„There were failing tests"), danach exakt zurückgesetzt
  (`diff` gegen die Vor-Mutation-Fassung bestätigt Identität) und erneut grün
  gebaut.

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-sdk-kotlin-projektgeruest`
in `done/` liegt (siehe Welle-Plan §4 Reihenfolge).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls sich beim
  Schreiben zeigt, dass Fehler-Antwortform und Datenklassen-Design für
  neun Fähigkeiten mehr als drei Liefer-Punkte brauchen — dann Aufteilung
  nach Fähigkeits-Gruppen (Consumer-Verwaltung / Tabellen-Verwaltung /
  Retention-und-Lesen).
- `in-progress` → `open` (blockiert — Carveout?): `slice-sdk-kotlin-projektgeruest`
  liegt noch nicht in `done/`.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- Die Fehler-Antwortform (`{"error": "<Klartext>"}`) lässt sich auf
  unterschiedliche Arten in Kotlin abbilden (Exception-Hierarchie vs.
  `Result`-Typ/sealed class) — eine falsche Wahl bindet spätere Consumer an
  ein API-Design, das ein Major-Bump bräuchte, um es zu ändern. **Ausgang:**
  entschieden beim Schreiben. `ADR-0109` Festlegung 4 bindet die
  SemVer-Major-Boundary an Draht-Änderungen, nicht an dieses interne
  Design — ein API-Redesign bleibt vor `1.0.0` folgenlos möglich.
  **Ausgang (Closure): entschieden, real bestätigt.** `PgChangeFeedException`
  ist eine sealed class mit sieben Unterklassen, konsistent mit der
  C#-Fassung geprüft (Review §Negativbefunde, Verifikation §1.1) — inklusive
  eines von Anfang an behandelten malformten `2xx`-Erfolgsbodys
  (`PgChangeFeedMalformedResponseException`), der bei C# erst eine
  Fixrunde brauchte. Risiko geschlossen.
- Ein netzloser Test-Fake für den gewählten Kotlin-HTTP-Client könnte
  reale Netzwerk-/Serialisierungs-Eigenheiten (z. B. Groß-/Kleinschreibung
  der JSON-Felder) verdecken, die erst gegen einen echten Server auffielen.
  **Ausgang:** weiter offen — ein realer Rundlauf-Beleg bleibt
  `make test-integration`s bestehendem `tools/harness/httpclient`
  vorbehalten (Wegwerf-Client, kein SDK-Import, `ADR-0068`); dieses SDK
  bekommt frühestens mit einem Folge-Slice einen eigenen
  Integrationsbeleg — dieselbe Grenze wie bei den beiden vorigen SDKs.
  **Ausgang (Closure): weiter offen, unverändert.** Weder C#- noch
  Python-Geschwister haben für diese strukturelle Grenze einen
  Beobachtungs-Register-Eintrag angelegt (`slice-sdk-csharp-http-client-flaeche`
  §7, `slice-sdk-python-http-client-flaeche` §7 — beide führen sie als
  Prosa-Risiko fort, kein Register); dieser Slice folgt demselben Muster,
  kein neuer Eintrag.
- `examples/kotlin/http-client` nutzt `java.net.http` ohne Fremdabhängigkeit
  (`ADR-0109` §Entscheidung Festlegung 1) — übernimmt das SDK dieselbe
  Wahl unreflektiert, ohne die für Tests nötige Ersetzbarkeit
  (`HttpClient`-Injektion) zu prüfen, könnte das netzlose Testen
  erschweren. **Ausgang:** weiter offen, entschieden beim Schreiben.
  **Ausgang (Closure): entschieden, real geprüft.** `java.net.http.HttpClient`
  bot keinen Pluggable-Handler; die Lösung (`internal fun interface
  HttpTransport` + zweiter `internal`-Konstruktor, Plan-Nachzug §3) macht
  alle Tests real netzlos (Review/Verifikation je unabhängig bestätigt: kein
  Socket, keine Loopback-Auflösung). Die Lösung selbst trug ein reales HIGH
  (F-1, siehe §7) — die Testbarkeits-**Wirkung** war korrekt, die
  begleitende KDoc-/Plan-Aussage über die Reichweite von `internal` war es
  nicht. Risiko geschlossen, mit dokumentiertem Lerneintrag.

## 7. Closure-Notiz

- **Was hat funktioniert:** Der Träger-Nachzug-Suchlauf (`AGENTS.md` §3.13,
  aus dem Beobachtungs-Register proaktiv in §3 aufgenommen) fing die
  stehende „follow in subsequent releases"-Aussage in
  `sdks/kotlin/pgchangefeed-kotlin/README.md` bereits im Implementer-Zug —
  anders als bei C# (dort erst dem Reviewer aufgefallen). Die sealed-class-
  Fehlerhierarchie vermied von Anfang an den bei C# erst per Fixrunde
  behobenen malformten-`2xx`-Body-Fehler (Review-Negativbefund, Verifikation
  §1.1). Drei unabhängige Docker-Builds (Implementer-Mutation-Beleg,
  Reviewer-Cache-Hit, Verifier `--no-cache`) bestätigten übereinstimmend
  denselben Grünzustand. `make gates` blieb über den gesamten Zyklus grün.
- **Was ging anders als geplant:** Eine Fixrunde war nötig — der Review fand
  1 HIGH (F-1): `HttpTransport.kt`s KDoc und der Plan-Nachzug (§3)
  behaupteten, der `internal`-Sekundärkonstruktor und `HttpTransport` selbst
  seien „außerhalb des Gradle-Moduls unsichtbar" bzw. fügten „keine
  öffentliche API-Fläche" hinzu. Real mit `javap -p` gegen das gebaute Jar
  geprüft, ist das falsch für den JVM-Bytecode: Kotlins `internal` ist eine
  **compile-time**-Grenze des Kotlin-Compiler-Frontends gegenüber anderen
  Kotlin-Modulen (Namens-Mangling), keine JVM-Zugriffsbeschränkung —
  Konstruktoren heißen im Bytecode immer `<init>` (nicht gemangelt) und ein
  `internal fun interface` kompiliert zu einem gewöhnlichen `public
  interface`. Ein Java-Konsument — den `ADR-0109` §Verglichene Alternativen
  A selbst als Zielgruppe benennt (Java-Binärkompatibilität) — kann beide
  technisch erreichen. Der Fix (`9181ae21`) korrigierte beide Träger auf die
  tatsächliche Grenze (Schutz gegen versehentliche Kotlin/Gradle-Modul-
  Nutzung, kein JVM-weiter Schutz); Verifikation bestätigte den neuen
  Wortlaut sachlich, nicht nur formal, als korrekt.
- **Steering-Loop-Eintrag (Lerneintrag):** Eine geschärfte Regel für künftige
  Kotlin-Arbeit in diesem Repo, ohne neuen Sensor (kein Gate-fähiger
  Tatbestand — dieselbe Klasse wie die bereits verkörperte
  `AGENTS.md` §3.7/§3.12-Disziplin: ein Kommentar/eine Aussage trägt seinen
  Beleg, hier speziell der JVM-Bytecode statt der Kotlin-Quellsemantik):
  **Kotlins Sichtbarkeitsmodifikatoren (`internal`, aber auch `private` auf
  Top-Level) sind Kotlin-Compiler-Grenzen, keine JVM-Bytecode-Grenzen** —
  eine Aussage über „unsichtbar außerhalb X" für kompilierten Kotlin-Code
  ist nur dann vollständig, wenn sie explizit auf den Kotlin/Gradle-
  Compile-Pfad eingeschränkt ist, nicht auf jeden JVM-Aufrufer (Java,
  Reflection). Diese Regel ist noch nicht in einem verkörperten Träger
  (Skill, ADR) verankert — sie steht hier als Vormerkung für den nächsten
  Kotlin-Slice.
- **Beobachtungs-Register (`../observations/`):** **Kein neues Verzeichnis
  angelegt** — bewusste Entscheidung, abweichend vom Vorschlag der
  Fixrunde (der Implementer-Zug schlug `BEO-PGC/kotlin-internal-faelschlich-als-jvm-zugriffsschutz`
  vor). Begründung: Das etablierte Muster dieses Repos öffnet ein neues
  `BEO-PGC/<slug>/`-Verzeichnis erst ab real gezählten **Mehrfach**-
  Vorkommen, nicht beim Erstfund — real belegt durch das nächstliegende
  Präzedens: `slice-sdk-csharp-http-client-flaeche` §7 hielt F-2
  (malformter `2xx`-Erfolgsbody) ausdrücklich **ohne** Registereintrag,
  „da bislang kein zweites Auftreten dieser spezifischen Klasse … im
  Bestand vorliegt". F-1 dieses Slice ist zwar über zwei **Träger**
  (Produktionscode-KDoc und Plan-Prosa) wiederholt, aber beide Träger
  tragen **dieselbe** ursprüngliche Fehleinschätzung desselben Autors
  innerhalb **desselben** Slice — kein zweites, unabhängiges Vorkommen in
  Zeit oder Kontext, sondern eine einzige Fundstelle, die sich in einen
  zweiten Text kopiert hat (dieselbe Struktur, die
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` bereits von einem
  „Zähler folgt den Vorkommen, nicht den Trägern" trennt). Der Reviewer
  selbst kategorisierte F-1 zudem unter die bereits **verkörperte**
  Skill-Klasse „Beleg trägt seinen Satz nicht" (`.harness/skills/reviewer.md`,
  aus `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`, Zähler bei 5× bereits
  geschlossen und in den Skill übernommen) — die allgemeine Regel „ein
  genannter Beleg muss die volle Aussage tragen" ist damit bereits
  gate-/skill-wirksam; ein zusätzliches, technisch enger gefasstes
  Register für genau die Kotlin-`internal`-Unterklasse wäre eine
  Vor-Verkörperung ohne zweites reales Vorkommen. **Stattdessen:**
  Vormerkung in diesem Absatz — sollte dieselbe Fehleinschätzung
  (Kotlin-Sichtbarkeitsmodifikator fälschlich als JVM-Grenze behauptet)
  in `slice-sdk-kotlin-grpc-client-flaeche` (plausibler Ort: eine
  ähnliche `internal`-Transport-Abstraktion dort ist im Plan bereits als
  wahrscheinlich benannt) oder einem späteren Kotlin-Slice ein zweites Mal
  real auftreten, öffnet der dortige Zug `BEO-PGC/kotlin-internal-faelschlich-als-jvm-zugriffsschutz`
  mit beiden Belegen (2× bereits bei Eröffnung) statt bei 1×.
- **Folge-Slices:** keine aus diesem Slice selbst erwartet — Umfang bleibt
  innerhalb der Welle (gRPC-Fläche, Pack-Werkzeug, Publish-Workflow bleiben
  eigene, bereits geplante Slices).
- **Risiken aus §6:**
  - „Fehler-Antwortform-Design könnte spätere Consumer binden" —
    **Ausgang: entschieden während der Umsetzung, geschlossen** —
    sealed-class-Hierarchie, sieben Unterklassen, konsistent mit C# geprüft,
    inklusive des dort erst per Fixrunde ergänzten
    Malformed-Response-Falls von Anfang an. Bleibt vor `1.0.0` folgenlos
    änderbar (`ADR-0109` Festlegung 4).
  - „Netzloser Test-Fake könnte reale Netzwerk-/Serialisierungs-Eigenheiten
    verdecken" — **Ausgang: weiter offen, unverändert** — dieselbe
    strukturelle Grenze wie bei den beiden vorigen SDKs, kein
    Registereintrag (weder C# noch Python legten einen an); ein realer
    Rundlauf-Beleg bleibt einem Folge-Slice bzw. `make test-integration`s
    bestehendem Wegwerf-Client vorbehalten.
  - „`java.net.http` ohne Pluggable-Handler könnte netzloses Testen
    erschweren" — **Ausgang: entschieden, geschlossen** —
    `internal fun interface HttpTransport` plus zweiter `internal`-
    Konstruktor löst es real (Review/Verifikation bestätigen Netzlosigkeit
    unabhängig); die Lösung selbst trug F-1 (siehe oben), inhaltlich aber
    korrekt und jetzt korrekt beschrieben.
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-kotlin-lh-fa-sst-009](../welle-sdk-kotlin-lh-fa-sst-009.md)
  (noch offen) — die Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `sdks/kotlin/` — mit
`slice-sdk-kotlin-projektgeruest` bereits eröffnet (GF, siehe dessen §8),
keine erneute Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
(bereits verkörpert) betrifft diesen Slice über den DoD-Punkt „Doku-Update"
oben; `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (19×) betrifft diesen
Slice über den in §3 aufgenommenen Suchlauf-Hinweis;
`BEO-PGC/unclosed-backtick-taeuscht-nackte-id-vor` (1×) über den
Backtick-Check in §3. Kein weiterer Treffer für diese Sub-Area.

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Fortsetzung von
`slice-sdk-kotlin-projektgeruest`).
