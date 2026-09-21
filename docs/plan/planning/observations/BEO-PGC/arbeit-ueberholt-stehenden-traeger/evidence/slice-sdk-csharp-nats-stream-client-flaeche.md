# Beleg: slice-sdk-csharp-nats-stream-client-flaeche

Vorgang: `slice-sdk-csharp-nats-stream-client-flaeche` liefert die
NATS-Vollinhalts-Client-Fläche (`PgChangeFeedNatsStreamClient`, `SPEC-024`)
und hebt `PgChangeFeed.Client` von `0.1.0` auf `0.2.0` — der vierte,
zuvor angekündigte Slice, der `PgChangeFeed.Client` von zwei auf vier
Client-Flächen bringt.

Fund: Der vorgeschriebene `AGENTS.md` §3.13-Suchlauf (`grep -rn "HTTP-API
und gRPC-Stream\|PgChangeFeed.Client.0.1.0" spec/ docs/ harness/`, plus
gezielte Prüfung von `docs/user/releasing.md`, `README.md`/`README.de.md`)
traf mehrere Stellen, die den Lieferstand von `PgChangeFeed.Client` als
„HTTP-API und gRPC-Stream" bzw. mit dem `0.1.0`-Artefaktnamen beschrieben
und durch diese Arbeit falsch wurden — real behoben, nicht nur gefunden:

- `spec/pflichtenheft.md` §1 (`LH-FA-SST-009.a`) und §6 (`SPEC-026`-Zeile)
  — die vom Slice-DoD selbst benannte Stelle.
- `sdks/csharp/README.md` §Status („NATS-vollinhalts delivery remains
  uncovered by this package") und die Einleitungszeile.
- `sdks/csharp/PgChangeFeed.Client/PgChangeFeed.Client.csproj`s
  `<Description>`-NuGet-Metadatum.
- `sdks/csharp/Dockerfile` (Kommentar-Erzeugnisname `0.1.0.nupkg`),
  `harness/mk/sdk.mk` (Kommentar-Erzeugnisname **und** „BEIDE Testflächen
  (HTTP + gRPC)"), `tools/harness/sdk-pack-csharp.sh` (dieselbe
  „BEIDE Testflächen"-Aussage — dieser dritte Fundort lag **außerhalb**
  des vorgeschriebenen `grep`-Musters: weder „HTTP-API und gRPC-Stream"
  noch „0.1.0" steht dort wörtlich, gefunden erst beim Lesen der Datei für
  den Docker-Bau-Kontext, nicht über den Suchlauf selbst).
- `harness/README.md` §Sensors, `make sdk-pack-csharp`-Zeile (dieselben
  zwei Aussagen).
- `README.md`/`README.de.md` §Distribution-Zeile („the HTTP API and gRPC
  stream are also available as an official C# client library").
- `sdks/csharp/PgChangeFeed.Client/Sse/Models/Change.cs`s Kommentar „three
  surfaces, three independent wire contracts that happen to share most
  field names" — dieser vierte Fundort lag ebenfalls **außerhalb** des
  `grep`-Musters (kein wörtlicher Treffer auf „HTTP-API und gRPC-Stream"
  oder „0.1.0"); aufgefallen beim Schreiben der neuen, wortgleich
  formulierten `Nats/Models/Change.cs`-Datei, nicht durch den Suchlauf
  selbst — derselbe Satz wäre sonst mit vier statt drei Wire-Contracts
  weiterhin bei „drei" stehen geblieben.

Bewusst **nicht** geändert (real geprüft, kein stiller Nachzug nötig):
`docs/plan/adr/0106-csharp-nuget-erstes-sdk-package.md`/`0109-…md` (ADRs,
`Accepted`, unberührbar — `AGENTS.md` §3.5), alle `docs/plan/planning/done/**`
und `docs/reviews/**`-Treffer (historische, abgeschlossene Records),
`spec/pflichtenheft.md`s datierte Änderungshistorie-Zeile vom 2026-09-19
(historischer Log-Eintrag, kein Ist-Stand-Satz), `docs/plan/planning/welle-sdk-csharp-vollabdeckung.md`
§2 Trigger-Begründung (beschreibt den Stand bei Welle-Eröffnung, nicht den
aktuellen), `docs/user/releasing.md`s `sdk-csharp-v0.1.0`-Erwähnungen
(Beleg eines real bereits erfolgten Releases, keine Ist-Stand-Behauptung),
und `sdks/kotlin/pgchangefeed-kotlin/README.md`/`spec/pflichtenheft.md`s
Kotlin-Absatz (Kotlin deckt real weiterhin nur HTTP-API und gRPC-Stream —
zutreffend, nicht stehen geblieben).

Einordnung: Diese Instanz reiht sich in die Erfolgsform („Implementer
findet über den vorgeschriebenen §3.13-Suchlauf selbst und behebt im
selben Commit", vgl. `slice-093`/`slice-094`/`slice-100`/`slice-101`/
`slice-102`) — mit einer Erweiterung: zwei der acht real behobenen
Fundstellen (`tools/harness/sdk-pack-csharp.sh`, `Sse/Models/Change.cs`)
lagen außerhalb des im DoD-Wortlaut vorgeschriebenen `grep`-Musters und
wurden nur durch aufmerksames Lesen benachbarter Dateien während der
Umsetzung gefunden, nicht durch den Suchlauf-Befehl selbst — derselbe
Grenzfall, den `AGENTS.md` §3.13 §Grenze bereits benennt
(„Zahlen"/Prosa-Umformulierungen trifft `grep` nicht zuverlässig), hier
erstmals für einen **Kommentar-Text** statt für eine Zahl oder einen
Symbolnamen belegt. Die acht Fundstellen betreffen zusammen **zehn**
Dateien, nicht neun — real ausgezählt per `git diff --stat` (korrigiert,
Fixrunde review-slice-sdk-csharp-nats-stream-client-flaeche F-3): die
Bündel-Fundstelle `sdks/csharp/Dockerfile`/`harness/mk/sdk.mk`/
`tools/harness/sdk-pack-csharp.sh` trägt drei Dateien statt
zwei, `README.md`/`README.de.md` trägt zwei statt einer.

Quelle: Implementer-Zug `slice-sdk-csharp-nats-stream-client-flaeche`, real
gegenseitig geprüfter Diff (`git diff --stat`) und die zitierten
Fundstellen-Grep-Läufe.
