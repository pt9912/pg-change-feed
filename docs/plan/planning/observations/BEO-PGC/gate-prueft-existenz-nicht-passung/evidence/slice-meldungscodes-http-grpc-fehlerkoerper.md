**Vorgang:** slice-meldungscodes-http-grpc-fehlerkoerper (Verifikation V-1; §5 (4))

**Fund (schwach):** Der Code `PCF-E2001` steht an zwei Stellen eines `503`/gRPC-`Internal`: im
Stream-Endpunkt (`sse.go`) und in `grpc/server.go`, wenn ein `ChangeStream` ohne verdrahteten
Broadcaster aufgerufen wird. Die Katalogzeile beschrieb den Code nur als fehlende oder falsche
Pflicht-Umgebungsvariable beim Start. `make meldungscodes-check` blieb Exit 0 (Code in Quelltext,
Tabelle und Katalog vorhanden). Der Verifier fand den Grenzfall beim Lesen der Emittenten je Code; der
Review hatte dieselbe Liste als passend geführt. Schwere LOW, Träger-Typ Fehlerkörper der API. Die
Katalogzeile im Handbuch nennt jetzt beide Fälle (Zustand und Ausgang regelkonform, kein eigener Code).
Gezählt als zweites, schwaches Vorkommen: die Zuordnung war teilweise tragend, nicht falsch.

Quelle: `docs/reviews/verifikation-slice-meldungscodes-http-grpc-fehlerkoerper.md` (V-1, §5 (4), §9). <!-- d-check:status-provenance -->
