**Vorgang:** slice-pin-digests-aktualisieren-2026-10-b
**Fund:** Alle acht Digest-Pins außerhalb des Inventars P1 bis P9 wichen vom
Index-Digest ihres Tags ab, ohne dass ein Sensor es meldete: `nats:2-alpine`,
`aquasec/trivy` und die fünf Basis-Images der SDK-/Beispiel-Dockerfiles
(`dotnet/sdk`, `dotnet/runtime`, `eclipse-temurin` JDK und JRE,
`python:3.14-slim`) sowie `postgres:17-alpine` als Einzelplattform-Digest.
Gemessen je Referenz mit `docker buildx imagetools inspect <Tag> --format
'{{.Manifest.Digest}}'` vom Planner (Stand `39e27242`) und vom Implementer
(Stand `7bc4aadd`) unabhängig, vom Verifier am Endstand als gleich dem Pin
bestätigt. Ursache der Form: sieben lebende Träger nannten „`docker manifest
inspect` (amd64)“ als Gewinnungsweg. Der Beleg ist das Zweitauftreten der
Klasse und das Erstauftreten für `nats`, `trivy` und die fünf SDK-/Beispiel-Basis-Images
([`ADR-0146`](../../../../../adr/0146-pin-inventar-quantifizierte-regel-alle-digest-pins.md)).
Slice-Aufzeichnung: Closure-Notiz in
[`slice-pin-digests-aktualisieren-2026-10-b`](../../../../in-progress/slice-pin-digests-aktualisieren-2026-10-b.md)
§7.
