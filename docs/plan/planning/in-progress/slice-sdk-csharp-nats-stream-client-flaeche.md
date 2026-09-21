# Slice sdk-csharp-nats-stream-client-flaeche: Öffentliche NATS-Vollinhalts-Client-Fläche (`SPEC-024`), Version-Hebung, Träger-Nachzug

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-csharp-vollabdeckung](../welle-sdk-csharp-vollabdeckung.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md),
[`LH-FA-SST-008`](../../../../spec/lastenheft.md) (Live-Change-Stream),
[`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md)
§Entscheidung Festlegung 1, letzter Absatz (NATS-Vollinhalt als
Folge-Package ausdrücklich antizipiert), [`ADR-0100`](../../adr/0100-nats-dritter-vollinhalts-zustellweg.md)
(NATS-Vollinhalts-Stream-Vertrag, wird vom SDK benutzt, nicht erweitert),
[`AGENTS.md`](../../../../AGENTS.md) §3.13 (Träger-Nachzug — dieser Slice
löst ihn für `spec/pflichtenheft.md` aus).

**Berührte Spec-Stellen:** [`SPEC-024`](../../../../spec/pflichtenheft.md)
(Subjekt-Schema, Nachrichtenform, Authentifizierung), `SPEC-026`
(`PgChangeFeed.Client`-Metadaten — Version-Hebung),
[`LH-FA-SST-009.a`](../../../../spec/pflichtenheft.md) (Träger-Nachzug:
„deckt HTTP-API und gRPC-Stream" wird durch diese Welle falsch).

**Verantwortlich:** Implementer-Agent, 2026-09-22.

**Autor:** Planner-Agent, direkt beauftragt (Nutzerauftrag „alle drei SDKs
auf volle Vier-Wege-Parität"). **Datum:** 2026-09-21.

---

## 1. Ziel und Abgrenzung

**Ziel:** Eine öffentliche, stabile .NET-API-Fläche im bestehenden Package
`PgChangeFeed.Client` für den NATS-Vollinhalts-Stream von
[`SPEC-024`](../../../../spec/pflichtenheft.md) — Verbindung über
`NATS.Net`, Subjekt-Abonnement `cdc.stream.<source_id>.<schema>.<table>`
(oder ein Platzhalter-Muster wie `cdc.stream.<source_id>.>`), Token-Auth
auf Verbindungsebene, Deserialisierung der zehn Nachrichtenfelder von
`SPEC-024` (dasselbe Schema wie SSE/gRPC). **Zusätzlich** bündelt dieser
Slice, weil er der letzte Flächen-Slice dieser Welle ist: Version-Hebung
von `PgChangeFeed.Client` und den Träger-Nachzug in
`spec/pflichtenheft.md`/`docs/user/benutzerhandbuch.md` für **beide** neu
gelieferten Flächen (SSE **und** NATS-Vollinhalt).

**Übernimmt:** `slice-sdk-csharp-sse-client-flaeche` — nicht den Code,
sondern dessen offen gelassenen Doku-/Träger-Nachzug (§1 jenes Slice
„Ausdrücklich NICHT in diesem Slice").

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **HTTP-API-, gRPC- oder SSE-Fläche** — HTTP/gRPC bereits geliefert; SSE
  ist der direkte Vorgänger-Slice dieser Welle.
- **Ein zweiter Vertriebsweg oder eine vierte Sprache** — unverändert
  `ADR-0106` §Re-Evaluierungs-Trigger 1, keine Entscheidung dieses
  Slice.
- **Ein realer `sdk-csharp-v<Version>`-Tag-Push** — Betreiber-Entscheidung
  nach `AGENTS.md` §3.10 (Welle-Plan §3), nicht Teil dieses Slice.
- **CLI-Argument-Parsing wie `examples/csharp/nats-stream-client`** —
  dasselbe Argument wie bei den übrigen Flächen: ein SDK-Consumer ruft
  eine Methode auf, parst keine `argv`.
- **Ein Realserver-Integrationstest** — dieselbe Begründung wie beim
  SSE-Slice (Welle-Plan §6): die Vorarbeits-Parität senkt das Wire-Risiko
  bereits, keine strengere Teststufe ohne fachlichen Grund.

## 2. Definition of Done

- [ ] `LH-FA-SST-009` erfüllt: eine öffentliche, im Package sichtbare
      Client-Klasse verbindet sich über `NATS.Net` mit Token-Auth,
      abonniert den Vollinhalts-Namensraum
      (`cdc.stream.<source_id>.<schema>.<table>`) und liefert die zehn
      Nachrichtenfelder von
      [`SPEC-024`](../../../../spec/pflichtenheft.md) — Tests
      referenzieren `SPEC-024` (Nachrichtenschema-Vollständigkeit,
      Subjekt-Formatierung netzlos, Verbindungs-Fehlerpfad ohne echten
      Server).
- [ ] `sdks/csharp/Directory.Packages.props` bekommt `NATS.Net` als
      gepinntes Paket, real zum Bau-Zeitpunkt neu gemessen (nicht blind
      aus `examples/csharp/Directory.Packages.props`s `3.2.0` übernommen —
      `AGENTS.md` §3.12).
- [ ] `PgChangeFeed.Client.csproj`s `<Version>` wird gehoben (von `0.1.0`
      auf `0.2.0` — additive, rückwärtskompatible Erweiterung, kein
      Breaking Change an HTTP-/gRPC-/SSE-Fläche).
- [ ] `spec/pflichtenheft.md` `LH-FA-SST-009.a`/`SPEC-026` nachgezogen: der
      Satz „deckt HTTP-API und gRPC-Stream" wird zu „deckt HTTP-API,
      gRPC-Stream, SSE und NATS-Vollinhalt" (`AGENTS.md` §3.13) —
      Suchlauf-Pflicht: `grep -rn "HTTP-API und gRPC-Stream\|PgChangeFeed.Client.0.1.0"`
      über `spec/`, `docs/`, `harness/` (Ergebnis wird in §7 berichtet,
      gefunden **und** nicht gefunden).
- [ ] Kein Import aus `internal/**`/`cmd/**`/`gen/**` dieses Repos.
- [ ] `make gates` grün.
- [ ] Ein real neu gebautes `.nupkg` (`make sdk-pack-csharp`) mit
      `<Version>0.2.0` als Smoke-Beleg — alle vier Client-Flächen im
      selben Artefakt.
- [ ] Review durchgeführt, Report unter `docs/reviews/` liegt vor
      (`.harness/skills/reviewer.md`) — Rollenwechsel nach Schritt 8 des
      Minimal Agent Workflow (`AGENTS.md` §6), kein Self-Review (Modul 8).
- [ ] Doku-Update: `docs/user/benutzerhandbuch.md` bekommt einen
      SDK-Hinweis für SSE **und** NATS-Vollinhalt (`ADR-0106`
      §Konsequenzen Folgepflicht 4-Muster) — getragen durch die bereits
      verkörperte Selbstprüf-Instruktion und den Reviewer-HIGH-Punkt
      (`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [ ] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo.
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis oder Beleg in `evidence/`; keine Beobachtung angefallen
      ist ebenfalls eine Antwort und wird in §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang.
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-csharp-vollabdeckung](../welle-sdk-csharp-vollabdeckung.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/csharp/PgChangeFeed.Client/Nats/PgChangeFeedNatsStreamClient.cs` (Arbeitsname) | neu | öffentliche API-Fläche für den NATS-Vollinhalts-Stream (`SPEC-024`). |
| `sdks/csharp/Directory.Packages.props` | update | `NATS.Net`-Paketversion ergänzen, real gemessen. |
| `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj` | update | `<Version>0.2.0</Version>`, `PackageReference` auf `NATS.Net`. |
| `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/Nats/*.cs` (Arbeitsname) | neu | Nachrichtenschema-Vollständigkeit, Subjekt-Formatierung, Verbindungs-Fehlerpfad. |
| `spec/pflichtenheft.md` §1 (`LH-FA-SST-009.a`), §6 (`SPEC-026`-Zeile) | update | Träger-Nachzug: volle Vier-Wege-Abdeckung für C#. |
| `docs/user/benutzerhandbuch.md` | update | SDK-Hinweis für SSE und NATS-Vollinhalt. |

**Ansatz:** Referenzmaterial ist `examples/csharp/nats-stream-client/`
(`Format.cs` für das Nachrichtenschema, `Program.cs`/`Cli.cs` für den
Verbindungs-/Subjekt-Aufbau, `NATS.Net`-Nutzung bereits real erprobt) und
die bestehende gRPC-Fläche dieses Packages (Fake-Verbindungs-Test-Muster,
`FakeCallInvoker.cs`-Analogie für einen gefakten NATS-Verbindungsfehler).

**Bekannte Docker-Bau-Kopplung:** wie beim SSE-Slice — `--build-context
proto=proto` bleibt für jeden Bau dieses Baums zwingend, unabhängig vom
NATS-Bezug (`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut`).

## 4. Trigger

**Start** (`next` → `in-progress`): wenn
`slice-sdk-csharp-sse-client-flaeche` in `done/` liegt (Welle-Plan §4
Reihenfolge — bewusst sequentiell, nicht parallel).

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): falls der
  Träger-Nachzug in `spec/pflichtenheft.md` mehr Stellen betrifft als
  erwartet (z. B. weitere `SPEC-026`-Referenzen außerhalb §1/§6) — dann
  Abspaltung eines reinen Doku-Nachzug-Slice.
- `in-progress` → `open` (blockiert — Carveout?): `slice-sdk-csharp-sse-client-flaeche`
  liegt noch nicht in `done/`.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + reales `.nupkg` mit `0.2.0` +
Closure-Notiz mit Lerneintrag geschrieben.

## 6. Risiken und offene Punkte

- Ein gefakter NATS-Verbindungsfehler-Test könnte das reale
  Token-Ablehnungsverhalten des NATS-Servers nicht exakt nachbilden.
  **Ausgang:** weiter offen — ein realer Rundlauf-Beleg bleibt
  `make test-integration`s bestehendem `tools/harness/natsstreamsub`
  vorbehalten (Wegwerf-Client, kein SDK-Import); dieselbe Teststrategie
  wie bei HTTP/gRPC/SSE.
- Der Träger-Nachzug in `spec/pflichtenheft.md` könnte eine Stelle
  übersehen, die denselben veralteten Satz an einer nicht per `grep`
  gefundenen Formulierung trägt (`AGENTS.md` §3.13 Grenze: Symbolnamen
  trifft `grep` zuverlässig, Prosa-Umformulierungen nicht). **Ausgang:**
  weiter offen — Review prüft den Diff unabhängig gegen denselben
  Suchraum.
- Version-Hebung `0.1.0` → `0.2.0` ohne SemVer-Minor-Konvention explizit
  in einer Datei dieses Repos festgeschrieben (nur `ADR-0106`
  §Entscheidung Festlegung 3 „ab `0.x.y`", kein Minor-Inkrement-Schema).
  **Ausgang:** entfallen — additive, rückwärtskompatible Erweiterung ist
  im SemVer-Vokabular unmissverständlich ein Minor-Bump, keine
  ADR-pflichtige Entscheidung.

## 7. Closure-Notiz

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <…>
- **Beobachtungs-Register (`../observations/`):** <…>
- **Folge-Slices:** <…>
- **Risiken aus §6:** <…>
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-csharp-vollabdeckung](../welle-sdk-csharp-vollabdeckung.md)
  (noch offen) — die Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Sub-Area `sdks/csharp/` — bereits
mit `slice-sdk-csharp-projektgeruest` eröffnet (GF), keine erneute
Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
`BEO-PGC/zusatzkontext-kopplung-breiter-als-dod-wortlaut` (§3 dieses Plans
trägt sie), `BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`
(verkörpert, DoD-Punkt „Doku-Update" trägt sie), `BEO-PGC/arbeit-ueberholt-stehenden-traeger`
(verkörpert, Suchlauf-Pflicht `AGENTS.md` §3.13 trägt den Träger-Nachzug
dieses Slice), `BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen`
(offen, 1×, nicht einschlägig — Publish-Mechanismus unverändert).

**Modus-Begründungsblock:** alle berührten Sub-Areas GF (Fortsetzung von
`slice-sdk-csharp-projektgeruest`).
