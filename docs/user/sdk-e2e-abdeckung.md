# SDK-E2E-Abdeckung je Spec-Kennung

Erzeugt von `make test-sdk-csharp-integration` über
`tools/harness/run-sdk-csharp-integration-tests.sh`; die Sprach-Runner
der Folge-Slices (Kotlin, Python-HTTP) erweitern dieselbe Datei um ihre
marker-gegrenzten Abschnitte. Je Sprach-Abschnitt deklariert der
zuständige Runner seine Realserver-Phasen an Ort und Stelle. Diese
Datei ist eine **stabile Abdeckungs-Deklaration**, kein Lauf-Beleg: der
Runner schreibt sie nur bei inhaltlicher Abweichung. Sie trägt nur
Zeilen real existierender Runner-Phasen — ein Beleg steht hier nie,
bevor sein Lauf grün lief.

| Spec-Kennung | Kurzbeschreibung | Nachweis | Ort |
| --- | --- | --- | --- |
<!-- pgchangefeed-sdk-e2e:csharp-begin -->
| [`LH-FA-SST-008`](../../spec/lastenheft.md), [`LH-FA-SST-009`](../../spec/lastenheft.md) | ein C#-SDK-Client (`PgChangeFeedGrpcClient`) oeffnet real den gRPC-Server-Stream gegen den laufenden Feed-Container und empfaengt eine danach committete Aenderung; ein Oeffnungsversuch ohne gueltiges Token endet mit gRPC-Status `Unauthenticated` | `GrpcRealserverTests` | `tools/harness/run-sdk-csharp-integration-tests.sh` |
| [`LH-FA-SST-008`](../../spec/lastenheft.md), [`LH-FA-SST-009`](../../spec/lastenheft.md) | ein C#-SDK-Client (`PgChangeFeedSseClient`) oeffnet real `GET /changes/stream` und empfaengt eine danach committete Aenderung; ein Aufruf ohne gueltiges Token endet mit HTTP-Status 401 | `SseRealserverTests` | `tools/harness/run-sdk-csharp-integration-tests.sh` |
| [`LH-FA-SST-008`](../../spec/lastenheft.md), [`LH-FA-SST-009`](../../spec/lastenheft.md) | ein C#-SDK-Client (`PgChangeFeedNatsStreamClient`) verbindet sich real per NATS und empfaengt eine danach committete Aenderung als vollstaendiges JSON-Event; ein Verbindungsversuch mit falschem Token wird vom NATS-Server abgelehnt | `NatsRealserverTests` | `tools/harness/run-sdk-csharp-integration-tests.sh` |
| [`LH-FA-SST-006`](../../spec/lastenheft.md), [`LH-FA-CON-001`](../../spec/lastenheft.md), [`LH-FA-SST-009`](../../spec/lastenheft.md) | ein C#-SDK-Client (`PgChangeFeedHttpClient`) registriert real einen Consumer (admin-Token) und listet Tabellen (reader-Token); die Registrierung ist unabhaengig ueber `cdc.consumer` lesbar; ein Aufruf ohne gueltiges Token endet mit HTTP-Status 401 | `HttpRealserverTests` | `tools/harness/run-sdk-csharp-integration-tests.sh` |
<!-- pgchangefeed-sdk-e2e:csharp-end -->
<!-- pgchangefeed-sdk-e2e:kotlin-begin -->
| [`LH-FA-SST-008`](../../spec/lastenheft.md), [`LH-FA-SST-009`](../../spec/lastenheft.md) | ein Kotlin-SDK-Client (`PgChangeFeedGrpcClient`) oeffnet real den gRPC-Server-Stream gegen den laufenden Feed-Container und empfaengt eine danach committete Aenderung; ein Oeffnungsversuch ohne gueltiges Token endet mit gRPC-Status `Unauthenticated` | `GrpcRealserverTest` | `tools/harness/run-sdk-kotlin-integration-tests.sh` |
| [`LH-FA-SST-008`](../../spec/lastenheft.md), [`LH-FA-SST-009`](../../spec/lastenheft.md) | ein Kotlin-SDK-Client (`PgChangeFeedSseClient`) oeffnet real `GET /changes/stream` und empfaengt eine danach committete Aenderung; ein Aufruf ohne gueltiges Token endet mit HTTP-Status 401 | `SseRealserverTest` | `tools/harness/run-sdk-kotlin-integration-tests.sh` |
| [`LH-FA-SST-008`](../../spec/lastenheft.md), [`LH-FA-SST-009`](../../spec/lastenheft.md) | ein Kotlin-SDK-Client (`PgChangeFeedNatsStreamClient`) verbindet sich real per NATS und empfaengt eine danach committete Aenderung als vollstaendiges JSON-Event; ein Verbindungsversuch mit falschem Token wird vom NATS-Server abgelehnt | `NatsRealserverTest` | `tools/harness/run-sdk-kotlin-integration-tests.sh` |
| [`LH-FA-SST-006`](../../spec/lastenheft.md), [`LH-FA-CON-001`](../../spec/lastenheft.md), [`LH-FA-SST-009`](../../spec/lastenheft.md) | ein Kotlin-SDK-Client (`PgChangeFeedHttpClient`) registriert real einen Consumer (admin-Token) und listet Tabellen (reader-Token); die Registrierung ist unabhaengig ueber `cdc.consumer` lesbar; ein Aufruf ohne gueltiges Token endet mit HTTP-Status 401 | `HttpRealserverTest` | `tools/harness/run-sdk-kotlin-integration-tests.sh` |
<!-- pgchangefeed-sdk-e2e:kotlin-end -->
