# Beleg: slice-sdk-python-sse-client-flaeche

Vorgang: `slice-sdk-python-sse-client-flaeche` liefert die öffentliche
SSE-Stream-Client-Fläche (`PgChangeFeedSseClient`, `SPEC-021`) und
erweitert `make test-sdk-python-integration` um eine zweite Phasen-Fläche
(SSE neben gRPC, je Phase eigene Testdatei, eigener Sentinel und
eigener Reject-Marker). Derselbe Slice stellte in seiner Fixrunde die
Testdatei-Übergabe des Integration-Images von docker run-Argument auf
Umgebungsvariable (`PGCHANGEFEED_TEST_FILE`, Dockerfile-`CMD` mit
`:?`-Guard) um.

Fund: Der §3.13-Suchlauf des Implementers war je Slice und je **einer**
bewegten Eigenschaft gebunden („SSE bleibt außerhalb des Packages“) und
zog die drei SDK-Träger (`__init__.py`-Docstring, `options.py`-Docstring,
`sdks/python/README.md` §Status) korrekt nach — zwei Funde derselben
Lücken-Struktur entgingen ihm:

- **Haupt-Review F-2 (MEDIUM, vom Reviewer gefunden):**
  `harness/README.md`s `make test-sdk-python-integration`-Zeile trug die
  Ein-Flächen-Form weiter („`pgchangefeed.grpc_client`“, „gRPC-Status
  `Unauthenticated`“), während der Runner seit dem Implementation-Commit
  (`6bbe99d9`) je Fläche eine Phase fährt; der Plan-Nachzug listete die
  drei SDK-Träger, aber nicht diese Zeile. Fixrunde 1 (`3c941b0b`) zog
  sie auf die Zwei-Flächen-Form.
- **Re-Review FR-1 (MEDIUM, vom Fixrunden-Reviewer gefunden):** die
  Fixrunde 1 stellte die Testdatei-Übergabe auf die ENV um — und trug
  die neue Form an drei Stellen korrekt (Dockerfile-Kommentar,
  Planzeile, Runner-Code), während die Ist-Zustands-Träger derselben
  Eigenschaft auf der Argument-Form stehen blieben: der
  Runner-Skriptkopf („Stufe-ENTRYPOINT“, „docker run-Argument“) und die
  neue `harness/README.md`-Zeile. Der committete §3.13-Suchlauf deckte
  nur die sprachliche Eigenschaft, nicht die Mechanik-Umstellung —
  FR-1 ist der Befund, den ein zweiter Suchlauf gemeldet hätte.
  Fixrunde 2 (`beeddc3c`) zog die drei Phrase-Stellen auf die ENV-Form
  und erweiterte den Suchlauf um die zweite Eigenschaft.
- **Dritte Stelle der Kette — Verifikation V-1 (LOW):** die erweiterte
  Suchlauf-Tabelle deklarierte als Vorher-Stand `6bbe99d9`; für den
  `harness/README.md`-Träger ist der wahre Vorher-Stand `3c941b0b` (die
  Phrase kam erst mit Fixrunde 1 in den Träger). Gelöst im
  Verifikations-Commit (`25fad0e1`).

Einordnung: dieselbe Lücken-Struktur zweimal in einem Slice — die
Suchlauf-Pflicht klebt am Slice-Diff, nicht an der bewegten Eigenschaft;
der unabhängige Reviewer blieb beide Male die tragende Suchlinie
(dieselbe Unter-Klasse „gefunden vom Reviewer, nicht vom
Implementer-Suchlauf“ wie bei `slice-095`/`slice-097`/
`slice-sdk-csharp-http-client-flaeche`/
`slice-sdk-kotlin-nats-stream-client-flaeche`). Geschärfte Lehre
(Closure-Notiz §7 des Slices): der §3.13-Suchlauf wird je **bewegter
Eigenschaft** geführt, nicht je Slice — jede mechanische Umstellung im
Fix-Diff (hier: docker run-Argument → Umgebungsvariable) ist eine eigene
Eigenschaft mit eigenem Suchlauf über beide Stände. Ausgang bleibt
verkörpert (`AGENTS.md` §3.13), kein neuer Schwellen-Übertritt — die
Schärfung ist eine Anwendungs-Schärfung der verkörperten Regel; ob sie
einen eigenen Satz im Träger trägt, prüft der Lese-Schritt der
Welle-Closure.

Quelle: Haupt-Review `slice-sdk-python-sse-client-flaeche` (F-2),
Fixrunden-Report desselben Slices (FR-1 + Addendum),
Verifikationsbericht desselben Slices (V-1), Commits `6bbe99d9`,
`3c941b0b`, `beeddc3c`, `25fad0e1`.