# Slice sdk-csharp-http-client-flaeche: Öffentliche HTTP-API-Client-Fläche (`SPEC-018`)

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-csharp-lh-fa-sst-009](../welle-sdk-csharp-lh-fa-sst-009.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md),
[`LH-FA-SST-006`](../../../../spec/lastenheft.md) (die neun Port-gedeckten
Fähigkeiten, die diese Fläche als Consumer anspricht),
[`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md)
Festlegung 1 (Umfang), [`ADR-0057`](../../adr/0057-http-grpc-api.md)
(HTTP/JSON-API-Vertrag, wird vom SDK benutzt, nicht erweitert).

**Berührte Spec-Stellen:** [`SPEC-018`](../../../../spec/pflichtenheft.md)
(Endpunkte, Token-Header-Form — das SDK benutzt diese Festlegungen,
verändert sie nicht).

**Verantwortlich:** Implementer-Agent, 2026-09-19.

**Autor:** Planner-Agent, direkt beauftragt (`ADR-0106` §Konsequenzen
Folgepflicht 1). **Datum:** 2026-09-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Eine öffentliche, stabile .NET-API-Fläche für alle neun
Port-gedeckten Fähigkeiten von [`SPEC-018`](../../../../spec/pflichtenheft.md)
(`RegisterConsumer`, `AcknowledgeConsumer`, `GetConsumerPosition`,
`RemoveConsumer`, `EnableTable`, `DisableTable`, `GetStatus`,
`ListTables`, `RunRetention`) sowie das Changes-Lesen
(`GET /changes`, [`SPEC-022`](../../../../spec/pflichtenheft.md),
[`ADR-0081`](../../adr/0081-changes-lesen-ueber-die-http-api.md)) — Bearer-
Token-Auth (`reader`/`admin`), eigene, von den Beispiel-Clients
unabhängige Tests.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **CLI-Argument-Parsing wie `examples/csharp/http-client`** —
  `ADR-0106` Festlegung 2 verlangt eine öffentliche, stabile API statt
  eines Kommandozeilen-Programms; ein SDK-Consumer ruft Methoden auf,
  parst keine `argv`.
- **gRPC-Stream** — `slice-sdk-csharp-grpc-client-flaeche` übernimmt das;
  getrennter Draht-Vertrag (`SPEC-020`), getrennte Fremdabhängigkeit
  (`ADR-0106` Festlegung 1 Abhängigkeits-Footprint-Begründung).
- **SSE- (`SPEC-021`) oder NATS-Vollinhalts-Stream (`SPEC-024`)** —
  `ADR-0106` Festlegung 1 grenzt v1 ausdrücklich auf HTTP-API und
  gRPC-Stream ein (Welle-Plan §6 Out-of-Scope).
- **Diagnose/Health-Endpunkte** — `SPEC-018` grenzt sie ausdrücklich aus
  den neun Port-gedeckten Fähigkeiten aus; kein SDK-Umfang.
- **Änderung von `SPEC-018` selbst** — das SDK benutzt den bestehenden
  Draht-Vertrag, verlangt keine Vertragsänderung (`ADR-0106` §Kontext
  Bindung 2, Analogie zu `SPEC-023`s „Verhältnis zum Draht").

## 2. Definition of Done

- [x] `sdks/csharp/PgChangeFeed.Client/Http/` (oder gleichwertiger
      Namensraum) trägt eine öffentliche Client-Klasse mit einer Methode je
      der neun Port-gedeckten Fähigkeiten von
      [`SPEC-018`](../../../../spec/pflichtenheft.md) plus dem Changes-Lesen
      (`SPEC-022`) — Signatur, Request-/Response-Form und Fehler-Antwortform
      (`400`/`401`/`403`/`404`/`500` → typisierte Exception oder
      Result-Form, konsistent über alle Methoden) spiegeln
      `SPEC-018`/`SPEC-022` exakt. Bearer-Token wird bei Konstruktion
      übergeben (kein globaler State).
- [x] Eigene Tests (xUnit, analog dem Pinnungs-Muster in
      `examples/csharp/Directory.Packages.props`) decken je Fähigkeit
      mindestens den Happy Path und die Auth-Boundary (`401` fehlendes/
      unbekanntes Token, `403` `reader`-Token gegen einen `admin`-Endpunkt)
      ab — netzlos prüfbar (kein realer Server nötig, `HttpMessageHandler`-
      Fake analog dem bestehenden Test-Muster der Beispiel-Clients).
- [x] Kein Import aus `internal/**`/`cmd/**` dieses Repos (`ADR-0106`
      §Kontext Bindung „Import-Grenze, hier ohne Ausnahme") — real geprüft:
      `grep -rn "internal/\|cmd/" sdks/csharp/` liefert keinen Treffer.
- [x] `make gates` grün.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [x] Doku-Update: `docs/user/benutzerhandbuch.md` bekommt einen
      SDK-Hinweis für die HTTP-Oberfläche (`ADR-0106` §Konsequenzen
      Folgepflicht 4) — getragen durch die bereits verkörperte
      Selbstprüf-Instruktion (`.claude/commands/implement-slice.md` Schritt
      17) und den Reviewer-HIGH-Punkt
      (`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`; keine Beobachtung angefallen
      ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-csharp-lh-fa-sst-009](../welle-sdk-csharp-lh-fa-sst-009.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/csharp/PgChangeFeed.Client/Http/PgChangeFeedHttpClient.cs` (Arbeitsname) | neu | öffentliche API-Fläche für die neun Port-gedeckten Fähigkeiten + Changes-Lesen. |
| `sdks/csharp/PgChangeFeed.Client/Http/Models/*.cs` | neu | typisierte Request-/Response-DTOs (spiegeln `SPEC-018`/`SPEC-022` JSON-Schemas). |
| `sdks/csharp/PgChangeFeed.Client/Http/PgChangeFeedException.cs` (Arbeitsname) | neu | typisierte Fehlerform für `400`/`401`/`403`/`404`/`500`. |
| `sdks/csharp/PgChangeFeed.Client.Tests/HttpClientTests.cs` (Arbeitsname) | neu | Happy Path je Fähigkeit (`SPEC-018`/`022` Referenz), Auth-Boundary `401`/`403`. |
| `docs/user/benutzerhandbuch.md` | update | SDK-Hinweis für die HTTP-Oberfläche, im selben Zug (`ADR-0106` Folgepflicht 4). |

**Ansatz:** Referenzmaterial ist `examples/csharp/http-client/TablesClient.cs`/
`TablesUrlBuilder.cs` (Draht-Kenntnis, kein `ProjectReference` — `ADR-0106`
Festlegung 2: „dieselbe Draht-Kenntnis, aber als eigenständiger,
paketierbarer Code neu geschrieben").

**Plan-Nachzug (Implementer-Zug):**

- Statt einer einzelnen `HttpClientTests.cs` (Arbeitsname) liegen die Tests
  in sechs Dateien unter
  `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Http/`
  (`PgChangeFeedHttpClientConstructionTests.cs`,
  `PgChangeFeedHttpClientConsumerTests.cs`,
  `PgChangeFeedHttpClientTableTests.cs`,
  `PgChangeFeedHttpClientRetentionAndChangesTests.cs`,
  `PgChangeFeedHttpClientAuthBoundaryTests.cs`, plus die gemeinsamen
  Test-Helfer `FakeHttpMessageHandler.cs`/`TestClientFactory.cs`) —
  gruppiert nach Fähigkeits-Gruppe (Consumer-/Tabellen-Verwaltung,
  Retention-und-Lesen, Auth-Boundary/Fehler-Mapping) statt einer einzigen,
  ca. 30 Tests tragenden Datei; reine Lesbarkeits-Entscheidung, kein
  fachlicher Umfangs-Unterschied.
- `sdks/csharp/PgChangeFeed.Client/Http/Models/*.cs` liegt in vier
  Gruppen-Dateien (`Consumers.cs`, `Tables.cs`, `Retention.cs`,
  `Changes.cs`) plus dem internen `ErrorResponse.cs` — bereits als
  Glob-Zielpfad im Plan vorgesehen, hier nur die konkrete Aufteilung
  benannt.
- Die neun Port-gedeckten Fähigkeiten von `SPEC-018` decken je genau einen
  Happy-Path-Test (`RegisterConsumer`, `AcknowledgeConsumer`,
  `GetConsumerPosition`, `RemoveConsumer`, `EnableTable`, `DisableTable`,
  `GetStatus`, `ListTables`, `RunRetention`) plus `ReadChanges`
  (`SPEC-022`); die Auth-Boundary (`401`/`403`) sowie `400`/`404`/`500`
  und ein Fallback für einen nicht dokumentierten Status sind zusätzlich
  je einmal generisch geprüft (nicht je Fähigkeit dupliziert) — das
  Fehler-Mapping selbst ist fähigkeits-unabhängiger Code
  (`PgChangeFeedHttpClient.BuildException`).

## 4. Trigger

**Start** (`next` → `in-progress`): wenn `slice-sdk-csharp-projektgeruest`
in `done/` liegt (siehe Welle-Plan §4 Reihenfolge).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls sich beim
  Schreiben zeigt, dass Fehler-Antwortform und DTO-Design für neun
  Fähigkeiten mehr als drei Liefer-Punkte brauchen — dann Aufteilung nach
  Fähigkeits-Gruppen (Consumer-Verwaltung / Tabellen-Verwaltung /
  Retention-und-Lesen).
- `in-progress` → `open` (blockiert — Carveout?): `slice-sdk-csharp-projektgeruest`
  liegt noch nicht in `done/`.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- Die Fehler-Antwortform (`{"error": "<Klartext>"}`) lässt sich auf
  unterschiedliche Arten in .NET abbilden (Exception-Hierarchie vs.
  Result-Typ) — eine falsche Wahl bindet spätere Consumer an ein API-Design,
  das ein Major-Bump bräuchte, um es zu ändern. **Ausgang:** weiter offen,
  entschieden beim Schreiben (`ADR-0106` Festlegung 3: die
  SemVer-Major-Boundary bindet ohnehin an Draht-Änderungen, nicht an
  dieses interne Design — ein API-Redesign vor `1.0.0` ist folgenlos
  möglich).
- Ein `HttpMessageHandler`-Fake für die Tests könnte reale
  Netzwerk-/Serialisierungs-Eigenheiten (z. B. Groß-/Kleinschreibung der
  JSON-Felder) verdecken, die erst gegen einen echten Server auffielen.
  **Ausgang:** weiter offen — ein realer Rundlauf-Beleg bleibt
  `make test-integration`s bestehendem `tools/harness/httpclient`
  vorbehalten (Wegwerf-Client, kein SDK-Import, `ADR-0068`); dieses SDK
  bekommt frühestens mit einem Folge-Slice einen eigenen Integrationsbeleg.

## 7. Closure-Notiz

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <Guide oder Sensor geschärft/ergänzt, oder
  „kein neuer Sensor" — je nach Lauf>.
- **Beobachtungs-Register (`../observations/`):** <neu angelegt | Beleg
  ergänzt | keine Beobachtung angefallen>.
- **Folge-Slices:** keine aus diesem Slice selbst erwartet — Umfang bleibt
  innerhalb der Welle.
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>.
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-csharp-lh-fa-sst-009](../welle-sdk-csharp-lh-fa-sst-009.md)
  (noch offen) — die Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `sdks/csharp/` — mit
`slice-sdk-csharp-projektgeruest` bereits eröffnet (GF, siehe dessen §8),
keine erneute Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
(bereits verkörpert) betrifft diesen Slice über den DoD-Punkt „Doku-Update"
oben; kein weiterer Treffer für diese Sub-Area (`grep`-Suche wie im
Welle-Plan §6 dokumentiert).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Fortsetzung von
`slice-sdk-csharp-projektgeruest`).
