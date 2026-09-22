# Beleg: slice-sdk-kotlin-nats-stream-client-flaeche

Vorgang: `slice-sdk-kotlin-nats-stream-client-flaeche` liefert die
NATS-Vollinhalts-Client-Fläche (`PgChangeFeedNatsStreamClient`, `SPEC-024`)
und hebt `pgchangefeed-kotlin` von `0.1.0` auf `0.2.0` — der zweite und
letzte Flächen-Slice der Welle `welle-sdk-kotlin-vollabdeckung`, der das
Package von zwei auf vier Client-Flächen bringt.

Fund: Der vorgeschriebene `AGENTS.md` §3.13-Suchlauf des Implementers
(`grep -rn "HTTP-API und gRPC-Stream\|pgchangefeed-kotlin-0\.1\.0"` über
`spec/`, `docs/`, `harness/`) zog `harness/README.md`s
`make sdk-pack-kotlin`-Zeile und `harness/mk/sdk.mk`s gleichlautenden
Kommentarblock korrekt nach — **übersah aber zwei Zeilen außerhalb des
eigenen Suchraums**: `sdks/kotlin/Dockerfile:77` und `:104` trugen
weiterhin wörtlich „`build/libs/pgchangefeed-kotlin-0.1.0.jar`" als
Beispiel-Dateiname der `pack`-/`publish`-Stufen. Nicht der
Implementer-Suchlauf selbst, sondern der unabhängige **Reviewer** fand
diese Stelle (F-2, MEDIUM, siehe Review-Report zu diesem Slice) —
dieselbe Unter-Klasse „gefunden vom Reviewer, nicht vom
Implementer-Suchlauf" wie bei früheren Belegen dieser Beobachtung
(`slice-095`, `slice-097`, `slice-sdk-csharp-http-client-flaeche`). Die
Fixrunde zog beide Zeilen auf `pgchangefeed-kotlin-0.2.0.jar` nach; ein
frischer Fixrunden-Review und der Verifier bestätigten je eigenständig,
dass kein `0.1.0`-Treffer mehr in dieser Datei steht (dritte bzw. vierte
unabhängige Bestätigung).

Bewusst **nicht** geändert (real geprüft, kein stiller Nachzug nötig):
`docs/plan/adr/0109-…md` (`Accepted`, unberührbar, `AGENTS.md` §3.5),
`docs/plan/planning/done/**`/`docs/reviews/**`-Treffer (historische,
abgeschlossene Records), `spec/pflichtenheft.md`s datierte
Änderungshistorie-Zeilen, `docs/user/releasing.md`s illustratives
`sdk-kotlin-v0.1.0`-Tag-Beispiel (kein Ist-Stand-Anspruch),
`welle-sdk-kotlin-vollabdeckung.md`s §2-Trigger-Absatz (beschreibt den
Stand bei Welle-Eröffnung, nicht den aktuellen).

Einordnung: Zwanzigster+eins Vorgang derselben Klasse
(`AGENTS.md` §3.13, verkörpert seit `welle-20`) — die tragende
Verteidigungslinie bleibt der unabhängige Reviewer, nicht der
Implementer-eigene Suchlauf: Genau wie beim analogen C#-NATS-Sibling-Slice
(`evidence/slice-sdk-csharp-nats-stream-client-flaeche.md`, dort allerdings
vom Implementer-Suchlauf selbst über zehn Dateien gefunden) driftete auch
hier ein Träger außerhalb des im DoD-Wortlaut vorgeschriebenen
Such-Musters — hier aber **nicht** vom Implementer, sondern erst vom
Reviewer entdeckt. Kein neuer Handlungsbedarf, kein neuer
Regelschärfungs-Anlass — die bestehende, bereits mehrfach bestätigte
Diagnose („Reviewer als zweite, tragende Suchlinie neben dem
Implementer-eigenen §3.13-Suchlauf") trägt erneut.

Quelle: Review-Report `slice-sdk-kotlin-nats-stream-client-flaeche` (F-2),
Fixrunden-Commit `6b279311`, Fixrunden-Review-Report (Bestätigung),
Verifikationsbericht §2 F-2 (vierte Bestätigung).
