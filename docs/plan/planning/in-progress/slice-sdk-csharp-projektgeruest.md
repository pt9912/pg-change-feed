# Slice sdk-csharp-projektgeruest: SDK-Projektgerüst `sdks/csharp/PgChangeFeed.Client/`

**Lifecycle:** Der Zustand dieses Slice ist das Verzeichnis, in dem diese
Datei liegt — eines von `open/`, `next/`, `in-progress/`, `done/`. Er
wechselt nur durch `git mv`.

**Welle:** [welle-sdk-csharp-lh-fa-sst-009](../welle-sdk-csharp-lh-fa-sst-009.md).

**Bezug:** [`LH-FA-SST-009`](../../../../spec/lastenheft.md) (Scope),
[`ADR-0106`](../../adr/0106-csharp-nuget-erstes-sdk-package.md)
Festlegung 1/2/3/5 (Umfang, Ort, Versionierung, was unberührt bleibt).

**Berührte Spec-Stellen:** [`LH-FA-SST-009.a`](../../../../spec/pflichtenheft.md)
(bleibt nach diesem Slice noch offen — der Träger-Nachzug folgt erst mit
`slice-sdk-csharp-pack-werkzeug`, wenn das Package real paketierbar ist).

**Verantwortlich:** Implementer-Agent, 2026-09-19.

**Autor:** Planner-Agent, direkt beauftragt (`ADR-0106` §Konsequenzen
Folgepflicht 1). **Datum:** 2026-09-19.

---

## 1. Ziel und Abgrenzung

**Ziel:** Ein neues, eigenständiges .NET-Projekt unter
`sdks/csharp/PgChangeFeed.Client/` anlegen — `.csproj` mit NuGet-Metadaten
(`PackageId=PgChangeFeed.Client`, `Version=0.1.0`, `Description`, `License`,
`RepositoryUrl`), ein eigenes, digest-gepinntes Docker-Bau-Setup (analog
`examples/csharp/Dockerfile`), ein englischsprachiges `README.md` und ein
leeres/minimales öffentliches API-Skelett — ohne jeden Import aus
`internal/**`/`cmd/**` dieses Repos.

**Ausdrücklich NICHT in diesem Slice** — je Punkt mit Begründung:

- **Öffentliche HTTP-/gRPC-API-Fläche** — `slice-sdk-csharp-http-client-flaeche`
  und `slice-sdk-csharp-grpc-client-flaeche` übernehmen das; dieser Slice
  liefert nur den Ort, in den sie schreiben (das API-Skelett bleibt leer
  oder trägt höchstens eine gemeinsame Konfigurations-/Auth-Grundklasse,
  falls sich beim Schreiben ein echter, unstrittiger gemeinsamer Nenner
  zeigt — kein Vorgriff auf Endpunkt-Methoden).
- **`make sdk-pack-csharp` (Pack-Werkzeug)** —
  `slice-sdk-csharp-pack-werkzeug` übernimmt das; dieser Slice liefert nur
  ein Docker-Bau-Setup, das `dotnet build`/`dotnet test` trägt (analog
  `examples/csharp/Dockerfile`s `build`-Stufe), kein `pack`-Aufruf und kein
  `make`-Ziel.
- **NuGet-Publish-Workflow** — `slice-sdk-csharp-publish-workflow`
  übernimmt das; dieser Slice pflegt nur die `<Version>` im `.csproj`, kein
  Tag-Trigger, kein Secret-Bezug.
- **Umbau von `examples/csharp/`** — `ADR-0106` Festlegung 2/5 verbietet
  das ausdrücklich; die Beispiele bleiben unverändertes Vorbild
  (`SPEC-023`).
- **Träger-Nachzug in `spec/pflichtenheft.md`** — `ADR-0106` §Konsequenzen
  Folgepflicht 2 bindet den Nachzug an den Zug, der das Package **real**
  existent macht (`dotnet pack` läuft real); das ist
  `slice-sdk-csharp-pack-werkzeug`, nicht dieser Slice, der nur ein
  Projekt-Gerüst ohne `pack`-Fähigkeit liefert.

## 2. Definition of Done

- [x] `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj` existiert:
      `TargetFramework net10.0` (Analogie zu `examples/csharp/*.csproj`,
      dieselbe .NET-Generation), `PackageId=PgChangeFeed.Client`,
      `Version=0.1.0` (`ADR-0106` Festlegung 3, Start bei `0.x.y`),
      `Description`, `Authors`, `License` (`MIT`, wie das Repo-Root-`LICENSE`),
      `PackageProjectUrl`/`RepositoryUrl` auf dieses Repo, `PackageReadmeFile`
      auf das SDK-eigene `README.md` verweisend. Kein
      `PackageReference`/`ProjectReference` auf einen privaten Baum dieses
      Repos (`ADR-0106` §Kontext Bindung „Import-Grenze, hier ohne
      Ausnahme").
- [x] `sdks/csharp/Dockerfile` (Bau-Kontext `sdks/csharp/`, eigenständig von
      `examples/csharp/Dockerfile` und der Wurzel-`Dockerfile`) mit
      digest-gepinnter `mcr.microsoft.com/dotnet/sdk`-Basis (real gemessener
      Digest zum Bau-Zeitpunkt, Kommentar-Pflicht analog
      `examples/csharp/Dockerfile`); `dotnet restore`/`build` laufen darin,
      kein `dotnet run`/Runtime-Stufe nötig (ein SDK ist keine
      startbare Anwendung).
- [x] `sdks/csharp/README.md` (Englisch, Vorgabe aus dem Auftrag) beschreibt
      Zweck, Installationsweg (`dotnet add package PgChangeFeed.Client`) und
      verweist auf das Repo-Root-`README.md` für den vollen Kontext; kein
      Duplikat der Draht-Doku (`SPEC-018`/`SPEC-020` bleiben die
      kanonische Quelle).
- [x] `make gates` grün.
- [x] Review durchgeführt, Report unter `docs/reviews/` liegt vor —
      [`review-slice-sdk-csharp-projektgeruest.md`](../../../reviews/review-slice-sdk-csharp-projektgeruest.md)
      (0 HIGH, 0 MEDIUM, 1 LOW ohne Fixrunde; Nachzug nach Skill-Regel
      „DoD-Checkbox-Nachzug ohne Fixrunde").
- [x] Doku-Update für `harness/README.md` entfällt in diesem Slice — kein
      neues `make`-Target entsteht hier (Pack-Werkzeug folgt in
      `slice-sdk-csharp-pack-werkzeug`).
- [ ] Closure-Notiz mit Steering-Loop-Lerneintrag.
- [x] Reconciliation-Register (`../reconciliation.md`) fortgeschrieben,
      **falls dieser Slice einen Inventur-Fund auflöst** — entfällt: keine
      Reconciliation-Datei in diesem Repo (kein Brownfield-Bootstrap).
- [ ] Beobachtungs-Register (`../observations/`) fortgeschrieben — neues
      Verzeichnis `BEO-PGC/<slug>/` oder eine weitere Datei in dessen
      `evidence/`; **kein Zähler wird gesetzt**, er folgt aus den Dateien.
      Keine Beobachtung angefallen ist ebenfalls eine Antwort und wird in
      §7 notiert.
- [ ] Jedes Risiko aus §6 trägt einen Ausgang (eingetreten / entfallen /
      weiter offen).
- [ ] Die drei Paarungen (Anker · Folge-Slice · Register) sind getragen —
      dieser Slice gehört zu
      [welle-sdk-csharp-lh-fa-sst-009](../welle-sdk-csharp-lh-fa-sst-009.md)
      (noch offen); die Prüfung läuft regelkonform bei deren Closure.

## 3. Plan (vor Code)

| Datei / Komponente | Änderungs-Art | Begründung |
|---|---|---|
| `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj` | neu | NuGet-Metadaten, `TargetFramework`, `Version=0.1.0` (`ADR-0106` Festlegung 2/3). |
| `sdks/csharp/PgChangeFeed.Client/` (leeres/minimales Skelett, z. B. `PgChangeFeedClientOptions.cs`) | neu | gemeinsamer Konfigurations-Nenner (Adresse, Token) für die beiden Folge-Slices, falls sich beim Schreiben ein echter zeigt — kein Vorgriff auf Endpunkt-Methoden. |
| `sdks/csharp/Dockerfile` | neu | digest-gepinnter Bau, analog `examples/csharp/Dockerfile`s `build`-Stufe, ohne Runtime-Stufe. |
| `sdks/csharp/README.md` | neu | Englisch, Installationsweg, Verweis auf Repo-Root-`README.md`. |
| `sdks/csharp/.gitignore` | neu (optional) | `bin/`/`obj/`, analog `examples/csharp/.gitignore`. |
| Testdatei (Platzhalter, falls das Skelett bereits eine echte Klasse trägt) | neu | Konstruktions-/Validierungstest der Optionsklasse (kein Draht-Verhalten, das kommt mit den Folge-Slices). |

**Plan-Nachzug (Implementer-Zug, `AGENTS.md`-Konvention „im selben Lauf
nachtragen"):**

- `sdks/csharp/Directory.Packages.props` — in der Tabelle oben nicht
  benannt, aber notwendig: Der Konstruktions-/Validierungstest der
  Optionsklasse braucht xUnit (`Microsoft.NET.Test.Sdk`, `xunit`,
  `xunit.runner.visualstudio`), zentral gepinnt analog
  `examples/csharp/Directory.Packages.props` (dieselben, dort am
  2026-09-17 gemessenen Versionen — keine neue, eigenständig zu bewertende
  Fremdabhängigkeit).
- Das Skelett zeigte tatsächlich einen echten gemeinsamen Nenner
  (`PgChangeFeedClientOptions`: `Address`/`ApiToken`) — die ADR nennt ihn
  bereits selbst (§Entscheidung Festlegung 1, „HTTP und gRPC teilen
  Auth-Header-Form … und Grundkonfiguration"); die Testdatei liegt unter
  `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.Tests/` (eigenes
  Testprojekt, analog `examples/csharp/http-client/HttpClient.Tests/`),
  nicht als loses „Platzhalter"-File.
- `sdks/csharp/README.md` referenziert das Repo-Root-`README.md`/`LICENSE`
  über absolute GitHub-Blob-URLs statt relativer Pfade — ein NuGet-Package
  wird außerhalb dieses Repository-Checkouts gelesen (NuGet.org-Paketseite),
  ein relativer Pfad wäre dort nicht auflösbar.

## 4. Trigger

**Start** (`next` → `in-progress`): sofort — `ADR-0106` ist `Accepted`,
keine weitere Vorbedingung.

**Rückführungen — vorab benennen, nicht erst im Nachhinein begründen:**

- `in-progress` → `next` (zu groß, zurück zur Zerlegung): entfällt aus
  heutiger Sicht — Umfang ist eine `.csproj`, ein Dockerfile, ein README.
- `in-progress` → `open` (blockiert — Carveout?): kein gepinntes
  `mcr.microsoft.com/dotnet/sdk`-Image mit `net10.0`-Unterstützung
  auffindbar — unwahrscheinlich, `examples/csharp/Dockerfile` nutzt es
  bereits real.

## 5. Closure-Trigger

DoD vollständig + `make gates` grün + Closure-Notiz mit Lerneintrag
geschrieben.

## 6. Risiken und offene Punkte

- Das leere API-Skelett könnte beim Schreiben der beiden Folge-Slices
  (`slice-sdk-csharp-http-client-flaeche`,
  `slice-sdk-csharp-grpc-client-flaeche`) doch keinen echten gemeinsamen
  Nenner zeigen (unterschiedliche Auth-Form trotz beide Bearer-Token) —
  dann bleibt das Skelett minimal (nur `.csproj`/Bau/README), kein
  erzwungener gemeinsamer Code. **Ausgang:** weiter offen, entschieden
  beim Schreiben der Folge-Slices.
- `net10.0` als `TargetFramework` schränkt den Konsumentenkreis auf sehr
  aktuelle .NET-Runtimes ein — ein SDK-Package hat potenziell breitere
  Zielgruppen als ein Docker-only-Beispielprogramm. **Ausgang:** weiter
  offen — Entscheidung bleibt bei diesem Slice (Analogie zu
  `examples/csharp`, keine Nutzungsdaten, die eine andere Wahl
  rechtfertigen); ein Multi-Targeting-Wechsel wäre eine spätere,
  eigenständige Entscheidung (`ADR-0106` §Re-Evaluierungs-Trigger 2,
  Nutzungsdaten-getrieben).

## 7. Closure-Notiz

- **Was hat funktioniert:** <…>
- **Was ging anders als geplant:** <…>
- **Steering-Loop-Eintrag:** <Guide oder Sensor geschärft/ergänzt, oder
  „kein neuer Sensor" — je nach Lauf>.
- **Beobachtungs-Register (`../observations/`):** <neu angelegt | Beleg
  ergänzt | keine Beobachtung angefallen>.
- **Folge-Slices:** `slice-sdk-csharp-http-client-flaeche`,
  `slice-sdk-csharp-grpc-client-flaeche` — beide bereits als Dateien in
  `open/` vorhanden.
- **Risiken aus §6:** <jedes mit genau einem Ausgang — siehe §6>.
- **Drei Paarungen:** dieser Slice gehört zu
  [welle-sdk-csharp-lh-fa-sst-009](../welle-sdk-csharp-lh-fa-sst-009.md)
  (noch offen) — die Prüfung läuft regelkonform bei deren Closure.

## 8. Sub-Area-Prüfungen und Modus-Begründung

**Vorgelagert — Sub-Area-Wahl prüfen:** Neue Sub-Area `sdks/csharp/` —
noch nie berührt, entsteht mit diesem Slice erstmalig. Analogie zur
bereits etablierten Sub-Area `examples/csharp/` (`ADR-0087`/`ADR-0090`,
Docker-only, digest-/paket-gepinnt) — dieselbe Konventionen-Dichte
übertragen, keine eigene Ausdifferenzierung nötig.

**Vorgelagert — offene Beobachtungen sichten:** Register durchgegangen —
kein Treffer für „Docker-only-Bau eines neuen Sprach-Ökosystems" oder
„NuGet-Metadaten" (`grep` über alle `observation.md`, siehe
Welle-Plan §6 Eröffnungs-Sichtung). Der einzige inhaltlich einschlägige
Eintrag (`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche`)
ist bereits verkörpert und betrifft die Folge-Slices (HTTP-/gRPC-Fläche),
nicht dieses Gerüst-Slice ohne Betreiber-Oberfläche.

**Modus-Begründungsblock:**

### Sub-Area: `sdks/csharp/`

- **Modus:** Greenfield — neuer Baum, Doku (`ADR-0106`) führt vor Code.
- **Konventionen-Dichte:** übernommen von `examples/csharp/`
  (`harness/conventions.md` §Modus-Deklaration, Default-Sub-Area `*`/`PGC`
  Greenfield) — Docker-only, digest-/paket-gepinnt, kein Host-`dotnet`.
- **Phase-Reife:** Phase 0 (Erstanlage) — Doc (`ADR-0106`) vollständig vor
  dem ersten Code dieser Sub-Area.
- **Evidenz-/Diskrepanz-Risiko:** niedrig (GF, keine Inventur nötig).
- **Reconciliation-Aufwand:** keiner — kein Brownfield-Bestand in diesem
  Baum.
