# Welle welle-sdk-csharp-lh-fa-sst-009 — Closure-Notiz

**Welle:** welle-sdk-csharp-lh-fa-sst-009
**Abschluss:** 2026-09-19
**Verantwortlich:** pt9912

## Was wurde geliefert?

[`ADR-0106`](../adr/0106-csharp-nuget-erstes-sdk-package.md) (`Accepted`)
vollständig umgesetzt — das erste offizielle Client-Bibliothek-Package für
[`LH-FA-SST-009`](../../../spec/lastenheft.md), in fünf Slices:

- **`slice-sdk-csharp-projektgeruest`**: `sdks/csharp/PgChangeFeed.Client/`
  neu angelegt — `.csproj` mit NuGet-Metadaten (`PackageId=PgChangeFeed.Client`,
  `Version=0.1.0`, `TargetFramework net10.0`), eigenständiges, digest-gepinntes
  Docker-Bau-Setup, englischsprachiges `README.md`, ein realer gemeinsamer
  Konfigurations-Nenner (`PgChangeFeedClientOptions`) für die beiden
  Folge-Flächen.
- **`slice-sdk-csharp-http-client-flaeche`**: öffentliche .NET-API-Fläche für
  alle neun Port-gedeckten Fähigkeiten von `SPEC-018` plus das Changes-Lesen
  (`SPEC-022`) — Bearer-Token-Auth, eine `PgChangeFeedException`-Hierarchie
  konsistent über alle zehn Methoden, 27 eigene, netzlos laufende xUnit-Tests.
- **`slice-sdk-csharp-grpc-client-flaeche`**: öffentliche .NET-API-Fläche für
  den `StreamChanges`-RPC (`SPEC-020`) — `.proto`-Bezug über einen
  zusätzlichen, benannten Docker-Bau-Kontext (kein committeter Stub, Muster
  `examples/csharp/grpc-client`), Nachrichtenschema-Vollständigkeits- und
  Authn-Boundary-Tests gegen einen Fake-`CallInvoker`.
- **`slice-sdk-csharp-pack-werkzeug`**: `make sdk-pack-csharp` (Docker-only,
  kein Gate) baut, testet und paketiert beide Flächen zu einem realen
  `PgChangeFeed.Client.0.1.0.nupkg` — Träger-Nachzug in
  `spec/pflichtenheft.md` §1 (`LH-FA-SST-009.a`) und §6 (neue `SPEC-026`-Zeile)
  sowie `harness/README.md` §Werkzeuge.
- **`slice-sdk-csharp-publish-workflow`**: `.github/workflows/sdk-csharp-release.yml`
  (Trigger `push: tags: ['sdk-csharp-v*']`, eigener Tag-Namensraum getrennt
  von `release.yml`s `v*`), Tag-vs-`.csproj`-Versionsabgleich, `dotnet nuget
  push` gegen NuGet.org, referenziert `NUGET_API_KEY` (angelegt, nicht durch
  diesen Slice).

`make gates` grün auf dem Endstand (813 Dateien, 0 `docs-check`-Befunde,
Coverage 82,70 % ≥ 80 %, `a-check`/`generated-sync`/`commit-traceability` je
ohne Befund). Ein reales `.nupkg` existiert (`slice-sdk-csharp-pack-werkzeug`
§7, dreifach unabhängig bestätigt). Kein realer `sdk-csharp-v*`-Tag wurde
gepusht (`git tag -l "sdk-csharp-v*"` leer) — diese Welle liefert bewusst nur
den **Mechanismus**, seine erste Anwendung bleibt eine gesonderte,
irreversible Entscheidung beim Auftraggeber (`AGENTS.md` §3.10, Welle-Datei
§3 vierter Punkt).

## Was hat funktioniert?

Über den gesamten Wellen-Zyklus hinweg trug eine durchgehende, mehrfach
unabhängige Docker-Build-Verifikation — kein einziger Slice verließ sich auf
einen einzelnen Lauf:

- `slice-sdk-csharp-projektgeruest`: drei unabhängige Builds (Implementer,
  Reviewer, Verifier), 0 Warnings/0 Errors, 5/5 Tests grün, keine Divergenz.
- `slice-sdk-csharp-http-client-flaeche`: vier unabhängige Builds über den
  ganzen Zyklus (Implementer, Reviewer-Erstlauf, Reviewer-Fixrunden-Nachprüfung,
  Verifier), 26/26 Tests grün vor der Fixrunde, 27/27 danach, keine Divergenz.
- `slice-sdk-csharp-grpc-client-flaeche`: fünf unabhängige Bestätigungen
  (Implementer, Reviewer je zweimal mit/ohne `--build-context proto=proto`,
  Verifier je zweimal), 32 Tests grün, glatter Durchlauf ohne Fixrunde.
- `slice-sdk-csharp-pack-werkzeug`: sechs eigenständige Bau-Läufe über den
  Zyklus inklusive einer vom Reviewer unabhängigen Mutations-Probe des
  Verifiers (andere Testdatei, andere Assertion, identisches Ergebnis: ein
  roter Test verhindert das `.nupkg` strukturell).
- `slice-sdk-csharp-publish-workflow`: Reviewer und Verifier maßen den
  Tag-Namensraum-Ausschluss gegenüber `ci.yml`/`e2e.yml`/`release.yml` je
  unabhängig selbst nach (`grep -n "tags"`), statt ihn zu übernehmen.

Keine Divergenz über keinen dieser Läufe — ein starker Beleg gegen
versteckten Nichtdeterminismus im neuen Sprach-Baum. Das etablierte
3-Commit-Move-Muster (`git mv` · Inhalt · `git mv`) hielt `AGENTS.md` §3.3
über alle fünf Slices sauber.

## Was ging anders als geplant?

Zwei der fünf Slices durchliefen eine Reviewer-Fixrunde, drei blieben beim
ersten Durchlauf unbeanstandet (0 HIGH/MEDIUM/LOW bzw. nur INFO):

- `slice-sdk-csharp-http-client-flaeche`: 1 HIGH (F-1, `sdks/csharp/README.md`
  behauptete weiterhin „follow-up release" für eine jetzt real gelieferte
  Fläche — Verstoß gegen `AGENTS.md` §3.13) + 1 MEDIUM (F-2, ein malformter
  `2xx`-Erfolgsbody führte zu einer rohen `JsonException` statt einer
  typisierten Exception) — beide in einer Fixrunde (`7bf7dece`) behoben, in
  der Fixrunden-Nachprüfung und unabhängig in der Verifikation ohne neuen
  Fund bestätigt.
- `slice-sdk-csharp-publish-workflow`: 1 MEDIUM (F-1, `docs/user/releasing.md`
  zog den neuen SDK-Release-Weg nicht mit) — behoben mit `c582eaaf`,
  Fixrunden-Nachprüfung ohne neuen Fund.
- `slice-sdk-csharp-projektgeruest`: 1 LOW ohne Fixrunde (DoD-Checkbox-Nachzug).
- `slice-sdk-csharp-grpc-client-flaeche`: 0 HIGH/MEDIUM/LOW, glatter
  Durchlauf.
- `slice-sdk-csharp-pack-werkzeug`: 0 HIGH/MEDIUM/LOW, 1 rein kosmetisches
  INFO ohne Fix.

Zusätzlich zwei durchgängige Plan-Nachzüge (kein Fehler, sondern die
vorgesehene Form): `spec/pflichtenheft.md`s `matrix`-Modul verbietet
mechanisch jede Referenz `spec → adr` (`slice-sdk-csharp-pack-werkzeug`) —
die neuen Zeilen bei `LH-FA-SST-009.a` und `SPEC-026` tragen den fachlichen
Inhalt ohne ADR-Bezug, dieselbe Struktur-Regel wie bei bestehenden
`SPEC-*`-Zeilen; und ein realer `Google.Protobuf`-Versionsunterschied
(`3.36.2` statt `3.36.1`) zwischen SDK und `examples/csharp/`, weil beide zu
unterschiedlichen Zeitpunkten real neu gemessen wurden (`AGENTS.md` §3.12).

## Trigger-Audit

Regeln dieser Sektion: Baseline-Regelwerk `modul-06-roadmap.md`
§Wellen-Closure-Prozedur, Schritt 2 — drei Artefaktklassen, je eine belegte
Feststellung.

- **Carveouts (Modul 7):** 0 offen — `docs/plan/carveouts/` enthält
  ausschließlich `.gitkeep`, kein Carveout referenziert `ADR-0106` oder eine
  dieser fünf Slices.
- **Bootstrap-aware Gates (Modul 13):** 0 betroffen — `coverage-gate` steht
  bereits auf der Endstufe (80 %, ausgeschöpft, `harness/mk/coverage.mk`) und
  wurde von keinem der fünf Slices berührt (reiner C#-Baum, außerhalb der
  netzlos prüfbaren Go-Fläche `./internal/...`+`./cmd/...`+`./gen/...`); die
  gemessene Coverage vor (82,80 %, `welle-release-pipeline-adr-0051-results.md`)
  und nach dieser Welle (82,70 %) bewegt sich innerhalb der für Go-Test-Rauschen
  erwartbaren Schwankung, kein struktureller Rückgang durch diese Welle.
  `.a-check.yml` `languages: go` liest `sdks/csharp/**` strukturell nicht
  (real bestätigt: `a-check` meldet 0 Befunde auf dem Endstand).
- **Entscheidung/ADR:** [`ADR-0106`](../adr/0106-csharp-nuget-erstes-sdk-package.md)s
  Re-Evaluierungs-Trigger sind durch diese Welle **nicht** ausgelöst:
  - Trigger 1 (zweite Sprache/zweiter Vertriebsweg verlangt) — nicht
    eingetreten; diese Welle liefert ausschließlich C#/NuGet, out-of-scope
    §6 der Welle-Datei.
  - Trigger 2 (real auf NuGet.org veröffentlicht + Nutzungsdaten legen eine
    andere Priorisierung nahe) — **nicht** eingetreten. Das Package ist real
    paketierbar (`.nupkg` existiert, `slice-sdk-csharp-pack-werkzeug`), aber
    nicht real veröffentlicht: `git tag -l "sdk-csharp-v*"` ist leer, kein
    Tag-Push erfolgte, `dotnet nuget push` lief nie real gegen NuGet.org.
    Ohne reale Veröffentlichung gibt es keine Download-Zahlen und keine
    Issue-Nachfrage — beide Trigger-Bestandteile sind unerfüllt, nicht nur
    einer. **Beobachtbare Auslösebedingung** (analog dem
    `ADR-0103`-Präzedenzfall in `welle-release-pipeline-adr-0051-results.md`):
    der erste reale `git tag sdk-csharp-v<SemVer> && git push --tags`, der
    einen grünen `sdk-csharp-release.yml`-Lauf erzeugt und danach messbare
    Nutzungssignale (NuGet.org-Downloadzahl, GitHub-Issues mit SDK-Bezug)
    zeigt — prüfbar über `git tag -l "sdk-csharp-v*"` (nicht mehr leer), einen
    grünen Actions-Lauf und einen NuGet.org-Paketseiten-Blick. Bis dahin
    bleibt `ADR-0106` Festlegung 1 (Umfang: HTTP + gRPC, kein SSE/NATS)
    unverändert in Kraft.
  - Trigger 3 (Abhängigkeits-Footprint von gRPC als echtes Consumer-Problem)
    — nicht prüfbar ohne reale Consumer, dieselbe Voraussetzung wie Trigger 2.
  - Trigger 4 (Server-SemVer erreicht `1.0.0`) — nicht eingetreten,
    `docs/user/version.md` bleibt unverändert bei `0.1.2` (diese Welle
    berührt die Datei nicht, `ADR-0106` Festlegung 3/5).

## Steering-Loop-Einträge

Kein neuer Eintrag erreicht in dieser Welle erstmals die 3×-Schwelle mit
Verkörperungs-Bedarf; die drei in den fünf Slices am häufigsten getroffenen
Beobachtungen waren bereits vor Wellen-Beginn verkörpert und bleiben es:

- `BEO-PGC/arbeit-ueberholt-stehenden-traeger` (verkörpert seit `AGENTS.md`
  §3.13): ein neuer Beleg in dieser Welle
  (`slice-sdk-csharp-http-client-flaeche`, F-1 — `sdks/csharp/README.md` §Status
  wurde durch die reale Auslieferung der HTTP-Fläche falsch), Zähler jetzt
  17× real ausgezählt (`state.md`). Bemerkenswert: gefunden hat den Fund der
  Reviewer, nicht der Implementer-eigene §3.13-Suchlauf — dieselbe bereits
  belegte Unter-Klasse wie bei `slice-095`/`slice-097`.
- `BEO-PGC/report-nackte-id-ohne-link` (verkörpert seit `slice-063` in
  `AGENTS.md` §3.9): ein neuer Beleg
  (`slice-sdk-csharp-grpc-client-flaeche` — nackte `ADR-0106`-Erwähnung ohne
  Backticks im eigenen Review-Report), Zähler jetzt 7× real ausgezählt.
  Geprüft und **nicht** bestätigt: ein möglicher dritter Fund in
  `slice-sdk-csharp-pack-werkzeug` ließ sich aus den Artefakten nicht
  unabhängig belegen (kein separater Fix-Commit im `git log` dieses Slice) —
  Feststellung ohne neue Beleg-Datei, `AGENTS.md` §3.12.
- `BEO-PGC/workaround-uebersieht-etabliertes-muster-im-bestand` (Erstauftreten
  `slice-104`): ein neuer Beleg (`slice-sdk-csharp-projektgeruest` — ein
  fehlendes explizites `using Xunit;` bei `ImplicitUsings` löste einen ersten
  roten Build aus, obwohl 13/13 bestehende `examples/csharp`-Testdateien das
  Muster bereits durchgängig trugen), Zähler jetzt 2× — weiterhin unter der
  3×-Schwelle, keine Welle-Closure-Leseauslösung durch diese Welle.

**Neu angelegt, noch unter der Schwelle:**
`BEO-PGC/release-mechanismus-nicht-in-releasing-doku-nachgezogen`
(`slice-sdk-csharp-publish-workflow`, Reviewer-Finding F-1 — ein neuer
Release-Mechanismus wurde in `docs/user/releasing.md` nicht mitgezogen, weil
kein Sensor das erzwingt) — Zustand **offen**, Zähler real 1×
(`evidence/slice-sdk-csharp-publish-workflow.md`, `state.md`). Kein
Steering-Loop-Eintrag hier, da die 3×-Schwelle nicht erreicht ist; bewusst
als eigene, getrennte Klasse angelegt statt den bestehenden, enger gezogenen
`BEO-PGC/handbuch-nicht-nachgezogen-bei-neuer-betreiber-oberflaeche` zu
dehnen (der nennt `docs/user/benutzerhandbuch.md`, nicht `releasing.md`).

## Beobachtungs-Register (Zeiger)

Der Zähler steht in [`../observations/`](../observations/). Kein Eintrag
dieser Welle hat die 3×-Schwelle **erstmals** erreicht; die drei bestätigten
Wiederholungen oben bleiben bereits verkörpert, der neue Eintrag steht bei
1× (siehe oben).

## Nebenbefunde (außerhalb des Scopes dieser Welle, hier nur gemeldet)

- **`docs/user/releasing.md` §1 ist überholt.** Die Verifikation von
  `slice-sdk-csharp-publish-workflow` stellte real per `git tag -l`/
  `git ls-remote`/`gh release list` fest, dass bereits drei echte
  Server-Release-Tags (`v0.1.0`–`v0.1.2`) mit realen GitHub-Releases
  existieren, während `docs/user/releasing.md` §1 weiterhin behauptet: „Zum
  Zeitpunkt dieses Dokuments wurde noch kein realer Release-Tag gesetzt."
  Die Drift entstand **vor** dieser Welle (Fixrunde `c582eaaf` änderte laut
  `git diff --stat` ausschließlich §4, §1 blieb unberührt) und außerhalb
  ihres Scopes — kein Folge-Slice wird hier angelegt, die Entscheidung liegt
  beim Coordinator.
- **`NUGET_API_KEY` existiert bereits real als Repository-Secret**
  (`gh secret list`: angelegt 2026-09-19T13:10:36Z, zeitlich neben
  `DOCKERHUB_TOKEN`/`DOCKERHUB_USERNAME`) — nicht durch einen Zug dieser
  Welle angelegt, vermutlich eine externe Betreiber-Handlung. Ein realer
  Tag-Push (`sdk-csharp-v0.1.0`) ist damit technisch möglich. Diese Welle löst
  ihn bewusst **nicht** aus (`AGENTS.md` §3.10: externe, irreversible Aktion,
  Entscheidung liegt beim Nutzer) — das Post-Push-Risiko bleibt strukturell
  weiter offen, unabhängig von der Secret-Existenz.

## Folge-Slices

Keine. `ADR-0106` ist mit dieser Welle für C#/NuGet vollständig umgesetzt.
Was ansteht, sind keine weiteren Slices dieser Welle, sondern eigenständige,
künftige Entscheidungen:

- Ein realer `sdk-csharp-v0.1.0`-Tag-Push — eine irreversible, extern
  sichtbare Aktion (öffentlicher NuGet.org-Push), die nur nach expliziter,
  gesonderter Rückfrage beim Auftraggeber läuft (Welle-Datei §3 vierter
  Punkt, `AGENTS.md` §3.10).
- Eine zweite SDK-Sprache oder ein zweiter Vertriebsweg — bleibt laut
  `ADR-0106` Re-Evaluierungs-Trigger 1 offen für eine künftige, separate ADR;
  kein Slice dieser Welle nimmt das vorweg.
- Der Nebenbefund `docs/user/releasing.md` §1 (siehe oben) — ein kleiner,
  eigenständiger Fix-Kandidat außerhalb dieser Welle.

## Verifikation

- `docs/reviews/review-slice-sdk-csharp-projektgeruest.md` (0 HIGH, 0 MEDIUM,
  1 LOW ohne Fixrunde).
- `docs/reviews/review-slice-sdk-csharp-http-client-flaeche.md` (1 HIGH + 1
  MEDIUM, Fixrunde `7bf7dece` real behoben, Fixrunden-Nachprüfung ohne neuen
  Fund) + `docs/reviews/verifikation-slice-sdk-csharp-http-client-flaeche.md`.
- `docs/reviews/review-slice-sdk-csharp-grpc-client-flaeche.md` (0
  HIGH/MEDIUM/LOW/INFO, kein Fixrunden-Bedarf) +
  `docs/reviews/verifikation-slice-sdk-csharp-grpc-client-flaeche.md`.
- `docs/reviews/review-slice-sdk-csharp-pack-werkzeug.md` (0 HIGH/MEDIUM/LOW,
  1 INFO ohne Fix) + `docs/reviews/verifikation-slice-sdk-csharp-pack-werkzeug.md`.
- `docs/reviews/review-slice-sdk-csharp-publish-workflow.md` (1 MEDIUM,
  Fixrunde `c582eaaf` real behoben, Fixrunden-Nachprüfung ohne neuen Fund) +
  `docs/reviews/verifikation-slice-sdk-csharp-publish-workflow.md`.
- `make gates`: grün auf dem Endstand (813 Dateien, 0 Befunde; Coverage
  82,70 % ≥ 80 %; `a-check`/`generated-sync`/`commit-traceability` je ohne
  Befund), ungepiped geprüft (`AGENTS.md` §3.9) nach jedem Commit dieser
  Welle sowie erneut vor dieser Closure.
- Reales `.nupkg`-Artefakt (`PgChangeFeed.Client.0.1.0.nupkg`) — dreifach
  unabhängig real erzeugt (Implementer, Reviewer, Verifier;
  `slice-sdk-csharp-pack-werkzeug` §7).
- `git tag -l "sdk-csharp-v*"`: leer — bestätigt, dass kein realer
  NuGet-Publish-Versuch stattfand (Welle-Datei §3 vierter Punkt eingehalten,
  Post-Push-Risiko bleibt strukturell offen nach `AGENTS.md` §3.10).
- `docs/plan/carveouts/`: nur `.gitkeep` — 0 offene Carveouts dieser Welle.

## Archivierung

Dieses Repo führt kein Archivierungs-Werkzeug für Wellen-Zeitdokumente (kein
`archiv`-Ziel in `Makefile`/`harness/mk/*.mk`, real geprüft per `grep`) — die
Bedingung für Schritt 4 der Closure-Prozedur ist nicht eingetreten, keine
Handarbeit als Ersatz.
