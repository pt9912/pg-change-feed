# Slice sdk-kotlin-sse-client-flaeche: Öffentliche SSE-Stream-Client-Fläche (`SPEC-021`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-kotlin-vollabdeckung](../welle-sdk-kotlin-vollabdeckung.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md),
[`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Live-Change-Stream —
Boundary: keine Zustellgarantie, kein Stream-internes Replay),
[`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
§Entscheidung Festlegung 1, letzter Absatz (SSE als Folge-Package
ausdrücklich antizipiert), Festlegung 3 (räumliche Trennung von
`examples/kotlin/`, hier unverändert gültig).

**Berührte Spec-Stellen:** [`SPEC-021`](../../../../spec/pflichtenheft.md)
(Endpunkt, Event-Form, Nachrichtenschema).

**Verantwortlich:** Implementer-Agent, 2026-09-22.

**Autor:** Planner-Agent, direkt beauftragt (Nutzerauftrag „alle drei SDKs
auf volle Vier-Wege-Parität"). **Datum:** 2026-09-21.

---

## 1. Ziel und Abgrenzung

**Ziel:** Eine öffentliche, stabile Kotlin-API-Fläche im bestehenden
Package `pgchangefeed-kotlin` für den SSE-Endpunkt von
[`SPEC-021`](../../../../spec/pflichtenheft.md) — `GET /changes/stream`,
Bearer-Token-Auth (`reader` oder `admin`), ein Frame-Parser (analog
`examples/kotlin/sse-client/SseStream.kt`), der die zehn Nachrichtenfelder
als `kotlinx.coroutines.flow.Flow<Change>` (oder gleichwertig idiomatisch,
dasselbe Muster wie die bestehende gRPC-Fläche) an den Consumer liefert.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **HTTP-API- oder gRPC-Fläche** — beide bereits geliefert
  (`slice-sdk-kotlin-http-client-flaeche`,
  `slice-sdk-kotlin-grpc-client-flaeche`).
- **NATS-Vollinhalts-Fläche** — eigener Slice
  (`slice-sdk-kotlin-nats-stream-client-flaeche`), getrennter Draht-Vertrag,
  unabhängige Fremdabhängigkeit (`io.nats:jnats`).
- **Stream-internes Replay oder `Last-Event-ID`-Auswertung** — `SPEC-021`
  trägt beides nicht.
- **Version-Bump, `spec/pflichtenheft.md`-Träger-Nachzug,
  `docs/user/benutzerhandbuch.md`-Nachzug** — gebündelt im Folge-Slice
  `slice-sdk-kotlin-nats-stream-client-flaeche` (Welle-Plan §4
  Reihenfolge).
- **CLI-Argument-Parsing wie `examples/kotlin/sse-client`** — ein
  SDK-Consumer ruft eine Methode auf, parst keine `argv`.
- **Ein Realserver-Integrationstest** — anders als bei der
  Python-Vollabdeckungs-Welle (`ADR-0110` Folgepflicht 1) verlangt weder
  `ADR-0109` noch diese Welle das für Kotlin (Welle-Plan §6, Begründung:
  Vorarbeits-Parität senkt das Wire-Risiko bereits).

## 2. Definition of Done

- [x] `LH-FA-SST-009` erfüllt: eine öffentliche, im Package sichtbare
      Client-Klasse öffnet `GET /changes/stream` (Bearer-Token in
      `Authorization`-Header), zerlegt SSE-Frames zu Events und liefert
      die zehn Nachrichtenfelder von
      [`SPEC-021`](../../../../spec/pflichtenheft.md) — Tests
      referenzieren `SPEC-021` (Frame-Parser netzlos, Authn-Boundary
      gegen eine gestubbte HTTP-Response ohne echten Server). Real
      umgesetzt: `PgChangeFeedSseClient.streamChanges()` (Package
      `io.github.pt9912.pgchangefeed.sse`), Frame-Parser `SseFrameParser`,
      Nachrichtenmodell `sse.model.Change` — 16 neue Tests (7
      `SseFrameParserTest` + 7 `PgChangeFeedSseClientAuthBoundaryTest` + 2
      `PgChangeFeedSseClientMessageSchemaTest`), real gezählt aus den
      JUnit-XML-Reports des Docker-Baus (`build/test-results/test/*.xml`);
      Gesamtsuite des Pakets 47 Tests, 0 Fehler.
- [x] Kein Import aus `internal/**`/`cmd/**`/`gen/**` dieses Repos — real
      geprüft (`grep -rn "internal/\|cmd/\|gen/" sdks/kotlin/`), Ausnahme
      nur Doku-Zitate.
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Report: `docs/reviews/review-slice-sdk-kotlin-sse-client-flaeche.md`
      (0 HIGH, 0 MEDIUM, zwei LOW, zwei INFO — keine Fixrunde nötig,
      DoD-Checkbox-Nachzug ohne Fixrunde per Skill-Regel).
- [x] Doku-Update (`docs/user/benutzerhandbuch.md`,
      `spec/pflichtenheft.md`): bewusst **nicht** in diesem Slice —
      gebündelt im Folge-Slice `slice-sdk-kotlin-nats-stream-client-flaeche`
      (§1 Abgrenzung). Träger-Nachzug-Suchlauf (`AGENTS.md` §3.13) gegen
      `sdks/kotlin/pgchangefeed-kotlin/README.md` durchgeführt: die dort
      stehengebliebene Aussage „SSE and NATS-vollinhalt delivery remain out
      of scope for this package's planned first full release" (Status-
      Absatz) wurde gefunden und korrigiert — SSE ist jetzt als gedeckt
      benannt, NATS bleibt offen benannt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`; keine Beobachtung angefallen
      ist ebenfalls eine Antwort und wird in §7 notiert.
- [x] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-kotlin-vollabdeckung](../welle-sdk-kotlin-vollabdeckung.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/.../sse/PgChangeFeedSseClient.kt` (Arbeitsname) | neu | öffentliche API-Fläche für den SSE-Stream (`SPEC-021`). |
| `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/.../sse/SseFrameParser.kt` (Arbeitsname) | neu | Frame-Parser (`event:`/`data:` → Event) — netzlos testbar, Muster `examples/kotlin/sse-client/SseStream.kt`. |
| `sdks/kotlin/pgchangefeed-kotlin/src/test/kotlin/.../sse/*Test.kt` (Arbeitsname, mehrere Dateien) | neu | Frame-Parser-Grenzfälle, Authn-Boundary, Nachrichtenschema-Vollständigkeit gegen `SPEC-021`. |

**Ansatz:** Referenzmaterial ist `examples/kotlin/sse-client/src/main/kotlin/cdcexamples/sse/SseStream.kt`
(Frame-Zerlegung, bereits real erprobt) und die bestehende HTTP-Fläche
dieses Packages (Auth-Header-Form, Fehlerklassen-Hierarchie — bereits real
mit `com.google.code.gson:gson` verdrahtet, dieselbe JSON-Bibliothek für
die Nachrichtenfelder).

**Bekannte Docker-Bau-Kopplung (`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut`,
2×, unter der Schwelle):** `sdks/kotlin/Dockerfile` trägt eine einzige
`build`-Stufe für alle Flächen — ein `docker build` dieses Slice braucht
deshalb weiterhin `--build-context proto=proto`, sonst bricht der Bau an
der bereits bestehenden `COPY --from=proto`-Zeile ab. `make sdk-pack-kotlin`
trägt den Flag bereits unconditional (`tools/harness/sdk-pack-kotlin.sh`) —
kein Änderungsbedarf, nur eine bewusst benannte, keine stille Kopplung.

## 4. Trigger

**Start** (`next` → `in-progress`): wenn diese Welle eröffnet ist und kein
anderer Slice in `in-progress/` liegt (WIP-Limit 1).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): entfällt aus
  heutiger Sicht — ein Endpunkt, ein Nachrichtenschema, bereits zweimal in
  diesem Package erprobtes Muster.
- `in-progress` → `open` (blockiert — Carveout?): der Frame-Parser lässt
  sich aus einem noch unbekannten Grund nicht ohne echten Server netzlos
  testen (unwahrscheinlich — `examples/kotlin/sse-client` belegt den
  Mechanismus bereits real).

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- Ein gestubbter HTTP-Response-Stream könnte das reale Chunked-Transfer-/
  Flush-Verhalten des Server-`http.Flusher` nicht exakt nachbilden.
  **Ausgang:** weiter offen — ein realer Rundlauf-Beleg bleibt
  `make test-integration`s bestehendem `tools/harness/sseclient`
  vorbehalten; dieselbe Teststrategie wie bei HTTP/gRPC dieses Packages.
- Die Docker-Bau-Kopplung (`--build-context proto=proto` auch ohne
  gRPC-Bezug) könnte bei einem künftigen, isolierten Bau-Versuch übersehen
  werden. **Ausgang:** weiter offen, aber transparent benannt (§3) — ein
  dritter Kotlin-Treffer würde `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut`
  auf die 3×-Schwelle heben.

## 7. Closure-Notiz

- **Was hat funktioniert:** Die bestehende HTTP-Fläche
  (`PgChangeFeedHttpClient`/`PgChangeFeedException`-Hierarchie) und der
  bereits real erprobte Frame-Parser aus
  `examples/kotlin/sse-client/SseStream.kt` trugen die Umsetzung fast
  vollständig — kein neuer Fremd-Abhängigkeitsbedarf, keine neue
  Fehlerklassen-Hierarchie: `PgChangeFeedSseClient` wirft dieselben sieben
  `PgChangeFeedException`-Subtypen wie `PgChangeFeedHttpClient`
  (`buildException`/`extractErrorMessage` dupliziert statt geteilt — dasselbe
  Muster wie der C#-Sibling `PgChangeFeed.Client.Sse.PgChangeFeedSseClient`).
  Die Docker-Bau-Kette (`bash tools/harness/sdk-pack-kotlin.sh`) lief beim
  ersten Versuch grün (Gradle-Wrapper-Layer bereits gecacht aus dem
  gRPC-Slice).
- **Was ging anders als geplant:** Die Plan-Tabelle (§3) nannte
  `kotlinx.coroutines.flow.Flow<Change>` als Zielform (analog der
  gRPC-Fläche); real umgesetzt wurde stattdessen ein `kotlin.sequences.Sequence<Change>`
  — kalt (nichts läuft vor der ersten Iteration) wie ein `Flow`, aber ohne
  die `flowOn(Dispatchers.IO)`-Umgehung, die ein `flow { }`-Builder um
  einen blockierenden `java.net.http.HttpClient.send()`-Aufruf bräuchte.
  Der Plan selbst öffnete diese Tür ausdrücklich ("oder gleichwertig
  idiomatisch"); die Design-Entscheidung samt Begründung steht in
  `PgChangeFeedSseClient`s Klassen-KDoc.
- **Steering-Loop-Eintrag:** Die aus der C#-SDK-Reihe übernommene Lektion
  zur `internal`-Sichtbarkeit (compile-time Kotlin-Grenze, keine
  JVM-Bytecode-Schranke) wurde von Anfang an korrekt in `SseTransport.kt`s
  KDoc dokumentiert, statt sie erst im Review nachzutragen — dieselbe
  Formulierung wie bei den beiden Vorgänger-Flächen
  (`HttpTransport`/`GrpcStreamTransport`), diesmal ohne Fixrunde.
- **Beobachtungs-Register (`../observations/`):** Keine neue Beobachtung
  angefallen. Der bereits geführte 2×-Fund
  `BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` wurde bewusst
  vermieden statt ein drittes Mal ausgelöst: jeder Docker-Bau dieses Slice
  trug `--build-context proto=proto` von Anfang an (§3 dieses Plans nennt
  den Grund) — der Zähler bleibt bei 2×.
- **Folge-Slices:** `slice-sdk-kotlin-nats-stream-client-flaeche` (nächster
  Slice der Welle, trägt zusätzlich den in diesem Slice bewusst
  ausgesparten Version-Bump und Doku-Träger-Nachzug für die SSE-Fläche).
- **Risiken aus §6:** Beide Risiken bleiben **weiter offen**, wie im Plan
  vorgesehen — keins wurde in diesem Slice aufgelöst: (1) der gestubbte
  HTTP-Response-Stream bildet nicht das reale Chunked-Transfer-/Flush-
  Verhalten nach, ein realer Rundlauf-Beleg bleibt
  `make test-integration`s `tools/harness/sseclient` vorbehalten; (2) die
  Docker-Bau-Kopplung wurde in diesem Slice erfolgreich vermieden (siehe
  Beobachtungs-Register oben), bleibt aber als latentes Risiko für künftige
  isolierte Bau-Versuche bestehen.
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-kotlin-vollabdeckung](../welle-sdk-kotlin-vollabdeckung.md)
  (noch offen) — die Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `sdks/kotlin/` — bereits
mit `slice-sdk-kotlin-projektgeruest` eröffnet (GF), keine erneute
Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` (2×, §3 dieses
Plans trägt sie), `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
(verkörpert), `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert),
`BEO-PGC/deutsches-fachwort-im-englischen-sdk-readme` (offen, 3×,
Architect-Entscheidung, kein Handlungsbedarf hier).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Fortsetzung von
`slice-sdk-kotlin-projektgeruest`).
