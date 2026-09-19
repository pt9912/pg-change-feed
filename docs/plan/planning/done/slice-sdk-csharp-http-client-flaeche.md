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
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
      Fixrunde geprüft und Merge-Block aufgehoben, siehe
      `docs/reviews/review-slice-sdk-csharp-http-client-flaeche.md`
      §Fixrunden-Nachprüfung.
- [x] Doku-Update: `docs/user/benutzerhandbuch.md` bekommt einen
      SDK-Hinweis für die HTTP-Oberfläche (`ADR-0106` §Konsequenzen
      Folgepflicht 4) — getragen durch die bereits verkörperte
      Selbstprüf-Instruktion (`.claude/commands/implement-slice.md` Schritt
      17) und den Reviewer-HIGH-Punkt
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
  das ein Major-Bump bräuchte, um es zu ändern. **Ausgang: entschieden beim
  Schreiben.** Es liegt eine `PgChangeFeedException`-Hierarchie vor,
  konsistent über alle zehn Methoden (jede läuft durch denselben privaten
  `SendAsync`-Pfad; Review und Verifikation haben das unabhängig
  nachgezählt) — inklusive der in der Fixrunde ergänzten
  `PgChangeFeedMalformedResponseException` für einen malformten
  `2xx`-Erfolgsbody (F-2 des Reviews). `ADR-0106` Festlegung 3 bindet die
  SemVer-Major-Boundary an Draht-Änderungen, nicht an dieses interne
  Design — ein API-Redesign bleibt vor `1.0.0` folgenlos möglich, das Risiko
  ist damit für den aktuellen Stand geschlossen, nicht nur vertagt.
- Ein `HttpMessageHandler`-Fake für die Tests könnte reale
  Netzwerk-/Serialisierungs-Eigenheiten (z. B. Groß-/Kleinschreibung der
  JSON-Felder) verdecken, die erst gegen einen echten Server auffielen.
  **Ausgang: weiter offen** — unverändert zur Plan-Begründung. Ein realer
  Rundlauf-Beleg bleibt `make test-integration`s bestehendem
  `tools/harness/httpclient` vorbehalten (Wegwerf-Client, kein SDK-Import,
  `ADR-0068`); dieses SDK bekommt frühestens mit einem Folge-Slice einen
  eigenen Integrationsbeleg. Weder Review noch Verifikation haben einen
  realen Server gegen dieses SDK gefahren — die vier unabhängigen
  Docker-Testläufe (Implementer, Reviewer-Erstlauf, Reviewer-
  Fixrunden-Nachprüfung, Verifier) liefen alle netzlos gegen den
  `HttpMessageHandler`-Fake.

## 7. Closure-Notiz

- **Was hat funktioniert:** Vier unabhängige Docker-Builds über den ganzen
  Zyklus (Implementer, Reviewer-Erstlauf, Reviewer-Fixrunden-Nachprüfung,
  Verifier) bestätigten übereinstimmend dasselbe Ergebnis — 26/26 Tests grün
  vor der Fixrunde, 27/27 nach ihr (ein neuer Test deckt die in der
  Fixrunde behobene Lücke), 0 Warnings/0 Errors, keine Divergenz über die
  vier Läufe. Die zentrale `SendAsync`-Bündelung aller zehn Methoden über
  genau einen privaten Sende-Pfad erwies sich als tragfähig: Sowohl das
  Fehler-Mapping (`BuildException`) als auch der in der Fixrunde ergänzte
  Erfolgs-Deserialisierungs-Helfer (`DeserializeSuccessBody`) griffen ohne
  Methoden-Zweitpfad für alle zehn Fähigkeiten — Review und Verifikation
  konnten das je unabhängig durch Verfolgen einzelner Methodenketten
  bestätigen, statt es aus der Struktur zu vermuten. Das 3-Commit-
  Move-Muster (`git mv` next→in-progress · Inhalt · Fixrunde) hielt
  `AGENTS.md` §3.3 sauber.
- **Was ging anders als geplant:** Zwei Plan-Nachzüge im
  Implementer-Zug (§3): Tests liegen in sechs Gruppendateien statt einer
  einzelnen `HttpClientTests.cs`, Modelle in vier Gruppendateien plus
  `ErrorResponse.cs` statt lose — reine Lesbarkeits-Entscheidung, kein
  fachlicher Umfangsunterschied (Fähigkeitsabdeckung deckungsgleich mit dem
  Plan, von Review und Verifikation je unabhängig nachgezählt). Zusätzlich
  war — anders als beim glatten Vorgänger-Slice — eine Fixrunde nötig: Der
  Review-Erstlauf fand 1 HIGH (F-1, `sdks/csharp/README.md` behauptete
  weiterhin „follow-up release" für eine jetzt real gelieferte Fläche —
  Verstoß gegen `AGENTS.md` §3.13) und 1 MEDIUM (F-2, ein malformter
  `2xx`-Erfolgsbody führte zu einer rohen `JsonException` statt einer
  typisierten Exception). Beide wurden in einer einzigen Fixrunde
  (`7bf7dece`) behoben und in der Fixrunden-Nachprüfung sowie unabhängig
  in der Verifikation als vollständig und ohne neuen Fund bestätigt.
- **Steering-Loop-Eintrag:** kein neuer Sensor — F-1 ist eine bereits
  verkörperte Hard Rule (`AGENTS.md` §3.13, kein Gate-fähiger Tatbestand,
  siehe Registerbegründung); F-2 ist eine lokale Code-Review-Fehlerklasse
  (unklare Fehlerbehandlung am Rand des Spec-Bereichs) ohne
  Schwellenwert-Charakter. Beide Findings sind stattdessen im
  Beobachtungs-Register festgehalten (F-1, siehe unten); F-2 bleibt ohne
  Registereintrag, da bislang kein zweites Auftreten dieser spezifischen
  Klasse (malformter Erfolgsbody statt Fehlerbody) im Bestand vorliegt.
- **Beobachtungs-Register (`../observations/`):** Beleg ergänzt — kein
  neues `BEO-PGC/<slug>/`. F-1 passt zur bereits verkörperten Beobachtung
  `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (Erstauftreten `slice-091`,
  bereits `AGENTS.md` §3.13): Der Vorgänger-Slice hatte
  `sdks/csharp/README.md` §Status bewusst mit Verweis auf genau diesen
  Folge-Slice geschrieben — bei Niederschrift wahr, durch die reale
  Auslieferung von `PgChangeFeedHttpClient` falsch geworden. Gefunden hat
  den Fund nicht der Implementer-eigene §3.13-Suchlauf, sondern der
  Reviewer — dieselbe bereits belegte Unter-Klasse „gefunden vom Reviewer,
  nicht vom Implementer-Suchlauf" wie bei `slice-095`/`slice-097`. Neue
  Evidenz-Datei:
  `../observations/BEO-PGC/arbeit-ueberholt-stehenden-traeger/evidence/slice-sdk-csharp-http-client-flaeche.md`;
  Zähler jetzt real ausgezählt 17× (siehe dortiges `state.md` — dabei fünf
  zwischenzeitlich ergänzte, in der Ordinal-Erzählung bislang unbenannte
  Belege nachgetragen, ohne deren Erzählung rückwirkend zu schreiben);
  Beobachtung bleibt bereits verkörpert, kein neuer Schwellen-Übertritt
  ausgelöst.
- **Folge-Slices:** keine aus diesem Slice selbst erwartet — Umfang bleibt
  innerhalb der Welle (gRPC-Fläche, Pack-Werkzeug, Publish-Workflow bleiben
  eigene, bereits geplante Slices).
- **Risiken aus §6:**
  - „Fehler-Antwortform-Design könnte spätere Consumer binden" —
    **Ausgang: entschieden während der Umsetzung** — Exception-Hierarchie,
    konsistent über alle zehn Methoden inkl. der neuen
    `PgChangeFeedMalformedResponseException` aus der Fixrunde; bleibt vor
    `1.0.0` folgenlos änderbar (`ADR-0106` Festlegung 3).
  - „`HttpMessageHandler`-Fake könnte reale Netzwerk-/
    Serialisierungs-Eigenheiten verdecken" — **Ausgang: weiter offen**,
    unverändert zur Plan-Begründung; ein realer Rundlauf-Beleg bleibt einem
    Folge-Slice vorbehalten.
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
