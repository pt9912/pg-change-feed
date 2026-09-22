# Slice sdk-kotlin-nats-stream-client-flaeche: Öffentliche NATS-Vollinhalts-Client-Fläche (`SPEC-024`), Version-Hebung, Träger-Nachzug

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-kotlin-vollabdeckung](../welle-sdk-kotlin-vollabdeckung.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md),
[`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Live-Change-Stream),
[`ADR-0109`](../../adr/0109-kotlin-github-packages-drittes-sdk-package.md)
§Entscheidung Festlegung 1, letzter Absatz (NATS-Vollinhalt als
Folge-Package ausdrücklich antizipiert), [`ADR-0100`](../../adr/0100-nats-dritter-vollinhalts-zustellweg.md)
(NATS-Vollinhalts-Stream-Vertrag), [`AGENTS.md`](../../../../AGENTS.md)
§3.13 (Träger-Nachzug — dieser Slice löst ihn für
`spec/pflichtenheft.md` aus).

**Berührte Spec-Stellen:** [`SPEC-024`](../../../../spec/pflichtenheft.md)
(Subjekt-Schema, Nachrichtenform, Authentifizierung), `SPEC-028`
(`pgchangefeed-kotlin`-Metadaten — Version-Hebung),
[`LH-FA-SST-009.a`](../../../../spec/pflichtenheft.md) (Träger-Nachzug:
„deckt HTTP-API und gRPC-Stream" wird durch diese Welle falsch).

**Verantwortlich:** Implementer-Agent, 2026-09-22.

**Autor:** Planner-Agent, direkt beauftragt (Nutzerauftrag „alle drei SDKs
auf volle Vier-Wege-Parität"). **Datum:** 2026-09-21.

---

## 1. Ziel und Abgrenzung

**Ziel:** Eine öffentliche, stabile Kotlin-API-Fläche im bestehenden
Package `pgchangefeed-kotlin` für den NATS-Vollinhalts-Stream von
[`SPEC-024`](../../../../spec/pflichtenheft.md) — Verbindung über
`io.nats:jnats`, Subjekt-Abonnement
`cdc.stream.<source_id>.<schema>.<table>` (oder ein Platzhalter-Muster wie
`cdc.stream.<source_id>.>`), Token-Auth auf Verbindungsebene,
Deserialisierung der zehn Nachrichtenfelder von `SPEC-024`. **Zusätzlich**
bündelt dieser Slice, als letzter Flächen-Slice dieser Welle:
Version-Hebung von `pgchangefeed-kotlin` und den Träger-Nachzug in
`spec/pflichtenheft.md`/`docs/user/benutzerhandbuch.md` für **beide** neu
gelieferten Flächen (SSE **und** NATS-Vollinhalt).

**Übernimmt:** `slice-sdk-kotlin-sse-client-flaeche` — nicht den Code,
sondern dessen offen gelassenen Doku-/Träger-Nachzug (§1 jenes Slice).

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **HTTP-API-, gRPC- oder SSE-Fläche** — HTTP/gRPC bereits geliefert; SSE
  ist der direkte Vorgänger-Slice.
- **Ein zweiter Vertriebsweg oder eine vierte Sprache** — unverändert
  `ADR-0109` §Re-Evaluierungs-Trigger 1.
- **Ein Wechsel zu Maven Central als Vertriebsweg** — unverändert
  `ADR-0109` §Re-Evaluierungs-Trigger 5, nicht Teil dieses Slice.
- **Ein realer `sdk-kotlin-v<Version>`-Tag-Push** —
  Betreiber-Entscheidung nach `AGENTS.md` §3.10 (Welle-Plan §3).
- **CLI-Argument-Parsing wie `examples/kotlin/nats-stream-client`** — ein
  SDK-Consumer ruft eine Methode auf, parst keine `argv`.
- **Ein Realserver-Integrationstest** — dieselbe Begründung wie beim
  SSE-Slice (Welle-Plan §6).

## 2. Definition of Done

- [x] `LH-FA-SST-009` erfüllt: eine öffentliche, im Package sichtbare
      Client-Klasse verbindet sich über `io.nats:jnats` mit Token-Auth,
      abonniert den Vollinhalts-Namensraum und liefert die zehn
      Nachrichtenfelder von
      [`SPEC-024`](../../../../spec/pflichtenheft.md) — Tests
      referenzieren `SPEC-024` (Nachrichtenschema-Vollständigkeit,
      Subjekt-Formatierung netzlos, Verbindungs-Fehlerpfad ohne echten
      Server). Real gebaut/getestet (`docker build --target build`), 14 neue
      Tests grün (4 Auth-Boundary + 3 Nachrichtenschema + 7 Subjekt), zwei
      rot färbende Mutationen real geprüft (§7).
- [x] `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` bekommt
      `io.nats:jnats` als Abhängigkeit, real zum Bau-Zeitpunkt neu
      gemessen (nicht blind aus `examples/kotlin/nats-stream-client/build.gradle.kts`s
      `2.26.3` übernommen — `AGENTS.md` §3.12). Real gemessen 2026-09-22
      gegen Maven Central: weiterhin `2.26.3`, keine Drift (§7).
- [x] `build.gradle.kts`s `version` wird gehoben (von `0.1.0` auf `0.2.0`
      — additive, rückwärtskompatible Erweiterung). Top-Level UND
      `publishing`-Block-Koordinate beide gehoben (real geprüft, beide
      Stellen getrennt geführt).
- [x] `spec/pflichtenheft.md` `LH-FA-SST-009.a`/`SPEC-028` nachgezogen: der
      Satz „deckt HTTP-API und gRPC-Stream" wird zu „deckt HTTP-API,
      gRPC-Stream, SSE und NATS-Vollinhalt" (`AGENTS.md` §3.13) —
      Suchlauf-Pflicht: `grep -rn "HTTP-API und gRPC-Stream\|pgchangefeed-kotlin-0.1.0"`
      über `spec/`, `docs/`, `harness/` (Ergebnis in §7 berichtet).
- [x] Kein Import aus `internal/**`/`cmd/**`/`gen/**` dieses Repos (real
      geprüft: alle Importe in den neuen Dateien sind `io.github.pt9912.pgchangefeed.*`,
      `com.google.gson.*`, `io.nats.client.*`, `java.*`/`kotlin.*`).
- [x] `make gates` grün (Exit 0, ungepiped geprüft, `AGENTS.md` §3.9 — §7).
- [x] Ein real neu gebautes `.jar` (`make sdk-pack-kotlin`) mit
      `version=0.2.0` als Smoke-Beleg — alle vier Client-Flächen im
      selben Artefakt. `pgchangefeed-kotlin-0.2.0.jar`, 144598 Bytes, real
      erzeugt (§7).
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Bleibt für den Reviewer-Rollenwechsel offen — kein Self-Review.
- [x] Doku-Update: `docs/user/benutzerhandbuch.md` bekommt einen
      SDK-Hinweis für SSE **und** NATS-Vollinhalt.
- [x] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [x] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`; keine Beobachtung angefallen
      ist ebenfalls eine Antwort und wird in §7 notiert. Keine neue
      Beobachtung angefallen — der Suchlauf war vollständig, kein Fund über
      dieses Slice hinaus (§7).
- [x] Jedes Risiko aus §6 trägt einen Ausgang (§7).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-kotlin-vollabdeckung](../welle-sdk-kotlin-vollabdeckung.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.
      Bleibt bewusst offen bis zur Welle-Closure (eigener, nachfolgender
      Schritt, nicht Teil dieses Slices).

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/kotlin/pgchangefeed-kotlin/src/main/kotlin/.../nats/PgChangeFeedNatsStreamClient.kt` (Arbeitsname) | neu | öffentliche API-Fläche für den NATS-Vollinhalts-Stream (`SPEC-024`). |
| `sdks/kotlin/pgchangefeed-kotlin/build.gradle.kts` | update | `io.nats:jnats`-Abhängigkeit ergänzen, `version` auf `0.2.0` heben. |
| `sdks/kotlin/pgchangefeed-kotlin/src/test/kotlin/.../nats/*Test.kt` (Arbeitsname) | neu | Nachrichtenschema-Vollständigkeit, Subjekt-Formatierung, Verbindungs-Fehlerpfad. |
| `spec/pflichtenheft.md` §1 (`LH-FA-SST-009.a`), §6 (`SPEC-028`-Zeile) | update | Träger-Nachzug: volle Vier-Wege-Abdeckung für Kotlin. |
| `docs/user/benutzerhandbuch.md` | update | SDK-Hinweis für SSE und NATS-Vollinhalt. |

**Ansatz:** Referenzmaterial ist
`examples/kotlin/nats-stream-client/src/main/kotlin/cdcexamples/natsstream/{Format,Cli,Main}.kt`
(Nachrichtenschema, Verbindungs-/Subjekt-Aufbau, `io.nats:jnats` bereits
real erprobt) und die bestehende gRPC-Fläche dieses Packages
(Fake-Verbindungs-Test-Muster für einen gefakten NATS-Verbindungsfehler).

**Bekannte Docker-Bau-Kopplung:** wie beim SSE-Slice — `--build-context
proto=proto` bleibt für jeden Bau dieses Baums zwingend, unabhängig vom
NATS-Bezug.

## 4. Trigger

**Start** (`next` → `in-progress`): wenn
`slice-sdk-kotlin-sse-client-flaeche` in `done/` liegt (Welle-Plan §4
Reihenfolge — bewusst sequentiell).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls der
  Träger-Nachzug mehr Stellen betrifft als erwartet — dann Abspaltung
  eines reinen Doku-Nachzug-Slice.
- `in-progress` → `open` (blockiert — Carveout?): `slice-sdk-kotlin-sse-client-flaeche`
  liegt noch nicht in `done/`.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + reales `.jar` mit `0.2.0` +
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Ein gefakter NATS-Verbindungsfehler-Test könnte das reale
  Token-Ablehnungsverhalten des NATS-Servers nicht exakt nachbilden.
  **Ausgang:** weiter offen — ein realer Rundlauf-Beleg bleibt
  `make test-integration`s bestehendem `tools/harness/natsstreamsub`
  vorbehalten.
- Der Träger-Nachzug in `spec/pflichtenheft.md` könnte eine Stelle
  übersehen, die denselben veralteten Satz an einer nicht per `grep`
  gefundenen Formulierung trägt. **Ausgang:** weiter offen — Review prüft
  den Diff unabhängig gegen denselben Suchraum.
- Version-Hebung `0.1.0` → `0.2.0` ohne explizites Minor-Inkrement-Schema
  in einer Datei dieses Repos. **Ausgang:** entfallen — additive,
  rückwärtskompatible Erweiterung ist im SemVer-Vokabular unmissverständlich
  ein Minor-Bump.

## 7. Closure-Notiz

- **Was hat funktioniert:** Das bestehende `SseTransport`/`GrpcStreamTransport`-Seam-Muster
  hat sich unverändert auf NATS übertragen lassen (`NatsStreamTransport`,
  ein `() -> ByteArray?`-Next-Payload-Supplier statt SSEs `() -> String?`-Zeilen-Supplier)
  — netzlose Tests ohne echten NATS-Server, exakt wie bei den beiden
  Geschwister-Flächen. Die C#-Geschwister-Fläche
  (`sdks/csharp/PgChangeFeed.Client/Nats/PgChangeFeedNatsStreamClient.cs`)
  war ein direkt übertragbares Formvorbild für Konstruktoren
  (Convenience/Advanced/Test-only), Subjekt-Validierung
  (`buildSubject`/`buildSourceSubject`) und die Entscheidung, KEINE zweite
  Fehlerklassen-Hierarchie für Verbindungsfehler zu erfinden (nur für
  Payload-Malformed: `PgChangeFeedNatsMalformedMessageException`).
  Kotlin-Spezifisch: `Nats.connect()` verbindet synchron (anders als jnats'
  Java-Pendant der C#-Bibliothek), was einen kleinen Pair-basierten
  privaten Sekundär-Konstruktor brauchte, um eine
  Konstruktor-Arity-Kollision zwischen „owned Connection bauen" und
  „Connection injizieren" zu vermeiden (`connectOwned`-Companion-Helfer,
  im Quelltext begründet).
- **Was ging anders als geplant:** Der Träger-Nachzug war breiter als der
  im Slice-Plan §2 explizit genannte `grep`-Suchraum
  (`spec/`, `docs/`, `harness/`) vermuten ließ — der reale Lauf über
  genau diesen Suchraum fand zusätzlich zwei stehende Stellen, die keiner
  der beiden Kotlin-Vorgänger-Slices (HTTP, gRPC, SSE) angefasst hatte:
  `harness/README.md`s `make sdk-pack-kotlin`-Zeile („BEIDE Testflächen
  (HTTP + gRPC)", `pgchangefeed-kotlin-0.1.0.jar`) und
  `harness/mk/sdk.mk`s gleichlautender Kommentarblock — beide waren seit
  der SSE-Fläche bereits falsch (sie hätten „HTTP + gRPC + SSE" tragen
  müssen), wurden aber übersehen. Beide jetzt auf „vier Testflächen (HTTP +
  gRPC + SSE + NATS-Vollinhalt)"/`0.2.0.jar` nachgezogen — exakt die
  `AGENTS.md` §3.13-Disziplin, angewendet auf einen Fund, den kein
  vorheriger Slice gemacht hatte.
- **Steering-Loop-Eintrag:** Der im Slice-Auftrag genannte, vorab bekannte
  Fallstrick „kein Chronik-Kommentar" und „jede Zahl real messen" wurde
  durchgehalten: `jnats` real gegen Maven Central gemessen (2.26.3,
  unverändert), Testzahlen aus dem echten Gradle-Testlauf (61 Tests
  gesamt, 0 Fehler, `test-results/test/*.xml` extrahiert und aufsummiert —
  nicht geschätzt), Jar-Größe (144598 Bytes) aus einem realen `ls -la`.
  Zwei rot färbende Mutationen real gebaut und rot gesehen (nicht nur
  behauptet): (1) `yield(parseChange(payload))` um `.copy(schema = "")`
  ergänzt → `PgChangeFeedNatsStreamClientMessageSchemaTest` schlägt fehl;
  (2) `nextPayload()`-Aufruf in `try { … } catch (ex: Exception) { break }`
  gehüllt (verschluckt die Transport-Exception) →
  `PgChangeFeedNatsStreamClientAuthBoundaryTest` schlägt exakt am
  erwarteten Test fehl. Beide Mutationen zurückgenommen, finaler Bau grün
  (Docker-Layer-Cache-Treffer auf den bereits verifizierten guten Stand —
  deterministischer Beleg, kein erneuter blinder Lauf).
- **Beobachtungs-Register (`../observations/`):** Keine neue Beobachtung
  angefallen. Der vollständige Suchlauf (§2 DoD-Punkt) deckte alle
  stehenden, jetzt falschen Träger ab (`spec/pflichtenheft.md` §1/§6/§7,
  `docs/user/benutzerhandbuch.md` HTTP-/gRPC-Absätze +
  SSE-/NATS-`**SDK:**`-Absätze + Versionshistorie, `sdks/kotlin/pgchangefeed-kotlin/README.md`,
  `harness/README.md`, `harness/mk/sdk.mk`) — kein Fund, der über dieses
  Slice hinaus offen bliebe. `docs/user/releasing.md` wurde geprüft, trägt
  aber keine Faktenbehauptung zur aktuellen Versionsabdeckung (nur ein
  illustratives `z. B. sdk-kotlin-v0.1.0`-Tag-Beispiel) — kein Nachzug
  nötig.
- **Folge-Slices:** Keine neuen Folge-Slices ausgelöst. Eine vierte Sprache
  oder ein vierter Vertriebsweg für `LH-FA-SST-009.a` bleibt weiterhin eine
  eigene, künftige ADR-pflichtige Entscheidung (unverändert). Die
  Welle-Closure (`welle-sdk-kotlin-vollabdeckung` → `done/` mit
  `welle-sdk-kotlin-vollabdeckung-results.md`) ist der nächste Schritt
  nach diesem Slice, aber laut Welle-Plan/Slice-Plan bewusst ein eigener,
  nachfolgender Zug — nicht Teil dieses Slices.
- **Risiken aus §6:** (1) Gefakter NATS-Verbindungsfehler-Test bildet nicht
  exakt das reale Token-Ablehnungsverhalten nach — **weiter offen**,
  unverändert wie geplant; ein realer Rundlauf-Beleg bleibt `make
  test-integration`s bestehendem `tools/harness/natsstreamsub` vorbehalten,
  diese Fläche hat keinen eigenen Realserver-Test (Slice-Plan §1
  Out-of-Scope). (2) Träger-Nachzug könnte eine Stelle übersehen —
  **weiter offen**, Review prüft den Diff unabhängig; real ist der
  Suchlauf breiter ausgefallen als der im Plan genannte Minimal-Suchraum
  (siehe „Was ging anders als geplant"), was das Risiko in der Praxis
  gesenkt, aber nicht auf Null gebracht hat. (3) Version-Hebung ohne
  explizites Minor-Schema — **entfallen**, wie geplant (additiver
  SemVer-Minor-Bump ist unmissverständlich).
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-kotlin-vollabdeckung](../welle-sdk-kotlin-vollabdeckung.md)
  (noch offen) — die Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `sdks/kotlin/` — bereits
mit `slice-sdk-kotlin-projektgeruest` eröffnet (GF), keine erneute
Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` (§3 dieses Plans
trägt sie), `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
(verkörpert, DoD-Punkt „Doku-Update" trägt sie), `BEO-PGC/arbeit-ueberholt-stehenden-traeger`
(verkörpert, Suchlauf-Pflicht trägt den Träger-Nachzug dieses Slice),
`BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen` (offen,
1×, nicht einschlägig).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Fortsetzung von
`slice-sdk-kotlin-projektgeruest`).
