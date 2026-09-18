# Welle welle-nats-drittstream — Closure-Notiz

**Welle:** welle-nats-drittstream
**Abschluss:** 2026-09-18
**Verantwortlich:** pt9912

## Was wurde geliefert?

- `ADR-0100` (`Accepted`) umgesetzt: NATS als dritter, paralleler
  Vollinhalts-Zustellweg für `LH-FA-SST-008`, neben gRPC und SSE — ein
  dritter Abonnent des bestehenden `Broadcaster`, ohne `ChangeStreamPort`,
  `CaptureService` oder das bestehende NATS-Wecksignal
  (`natsnotify`/`ADR-0055`/`ADR-0056`) anzufassen.
- `internal/adapters/driven/natsstream/` (`Publisher`): veröffentlicht
  jeden Change als vollständiges JSON-Event auf dem eigenen Subjekt-
  Namensraum `cdc.stream.<source_id>.<schema>.<table>`; Bootstrap-
  Verdrahtung über die Zwei-Bedingungen-Aktivierung (`CDC_NATS_URL` **und**
  `CDC_NATS_STREAM_TOKEN`); `compose.yaml`/`examples/compose.yaml` tragen
  die NATS-Server-Auth-Konfiguration für Test- und Demo-Umgebung.
- Volle Drei-Sprachen-Matrix für den neuen, vierten Zugriffsweg:
  `examples/nats-stream-client` (Go), `examples/csharp/nats-stream-client`,
  `examples/kotlin/nats-stream-client` — alle drei verbinden mit
  `CDC_NATS_STREAM_TOKEN`, abonnieren `cdc.stream.>` und geben jede
  empfangene Change lesbar aus. Mit dieser Welle ist die
  Beispiel-Client-Matrix zum ersten Mal tatsächlich vollständig: fünf
  Zugriffsarten × drei Sprachen, fünfzehn Programme
  (`examples/README.md`, `docs/user/benutzerhandbuch.md`,
  `spec/pflichtenheft.md` `SPEC-023` entsprechend nachgezogen).
- Realer `make test-integration`-Lauf am Endstand dieser Welle zeigt beide
  Belege aus `ADR-0100`s Fitness Function **in einem Durchlauf**: ein
  Wegwerf-Client (`tools/harness/natsstreamsub`) verbindet sich real mit
  gültigem Token und empfängt eine vollständige Change über
  `cdc.stream.<source_id>.<schema>.<table>` (`change_id=981-1`,
  unabhängig über `cdc.changes` gegengehalten); ein Verbindungsversuch
  ohne bzw. mit falschem Token wird vom NATS-Server abgelehnt
  (`REJECTED-NO-TOKEN`/`REJECTED-WRONG-TOKEN`); das bestehende Wecksignal
  (`natssub`) funktioniert mit demselben Test-Token unverändert weiter.
- `make gates`: grün auf dem Endstand (728 Dateien, 0 `docs-check`-Befunde,
  Coverage 82,80 % ≥ 80 %, `a-check`/`generated-sync`/
  `commit-traceability` je ohne Befund).

## Was hat funktioniert?

Die im `ADR-0100`-Slice-Schnitt vorgesehene Reihenfolge (Kernfähigkeit
zuerst, dann die beiden voneinander unabhängigen Beispiel-Client-Slices)
trug real: `slice-nats-drittstream-example-go` und
`slice-nats-drittstream-example-csharp-kotlin` brauchten sich zu keinem
Zeitpunkt gegenseitig, und beide bauten beim jeweils ersten realen
Docker-Lauf grün. Jeder der drei Slices belegte seine Behauptungen real
(reale Demo-Umgebungs-Rundläufe mit echtem Change-Empfang, nicht nur
Unit-Tests); der unabhängige Reviewer-Pass fuhr jede reale Verifikation
zusätzlich selbst nach (eigene Smoke-Tests mit eigenen Datensätzen, echte
`.dockerignore`-/Maven-Central-Nachmessungen) und fand über alle drei
Slices hinweg nur ein einziges HIGH-Finding und zwei LOW-Findings — beide
Klassen bereits etablierte Reviewer-Prüfpunkte
(„Beleg trägt seinen Satz nicht",
[`ADR-0073`](../../adr/0073-zitat-korrektur-an-immutablen-dokumenten.md)-Belegform),
keine neue Klasse.

## Was ging anders als geplant?

- Der Plan für `slice-nats-drittstream-example-go` ging von einer bereits
  vollständig offenen Handbuch-Sektion aus; real hatte
  `slice-nats-drittstream-core` die Sektion bereits mit einem
  Platzhalter-Absatz angelegt — der Implementer-Lauf ersetzte den
  Platzhalter, statt eine neue Sektion zu öffnen (Plan-Nachzug, keine
  Fehlplanung mit Konsequenz).
- Ein `make image`-Lauf während `slice-nats-drittstream-example-go`
  committete kurzzeitig einen Digest-Wechsel mit einer falschen
  Build-Kontext-Kausalbehauptung (Review-Fund F-1, HIGH) — real
  widerlegt gegen `.dockerignore` und zurückgenommen. Das Gespräch mit dem
  Auftraggeber im Anschluss deckte ein allgemeineres, wiederkehrendes
  Ärgernis auf (Merge-Rauschen durch einen nicht-deterministischen,
  inhaltslosen Digest) — siehe Steering-Loop-Einträge unten.
- `slice-nats-drittstream-example-csharp-kotlin`s Plan nahm eine rein
  mechanische Übertragung des `nats-client`-Musters an; real gebraucht
  wurde eine Vorbild-Mischung aus `nats-client` (Verbindung/Token) und
  `grpc-client` (Endlosschleife, strukturierte Feldausgabe), weil der
  Vollinhalts-Stream anders als das Wecksignal keinen zweiten
  HTTP-Holschritt braucht. Für Kotlin folgte daraus eine neue, im
  ursprünglichen Plan nicht vorgesehene Fremdbibliothek
  (`com.google.code.gson:gson`) — der Implementer-Lauf bewertete Version,
  Lizenz und transitive Abhängigkeiten mit derselben Rigorosität wie die
  bestehenden Pins und trug den Plan im selben Lauf nach.
- Eine seit `slice-nats-drittstream-core` mehrere Stunden laufende
  Demo-Umgebung hatte die dort neu eingeführte
  `CDC_NATS_STREAM_TOKEN`-Verdrahtung nicht übernommen (Container schon vor
  dem Env-/Image-Nachzug gestartet) — ein realer Closure-Rundlauf schlug
  dadurch zunächst fehl, bis die Umgebung neu hochgefahren wurde. Neue
  Beobachtung, siehe unten.

## Steering-Loop-Einträge

Kein Eintrag über der 3×-Schwelle in dieser Welle. Zwei neue Beobachtungen
unter der Schwelle (je 1×, keine Ausgangs-Entscheidung fällig):
`BEO-PGC/image-digest-nichtdeterminismus-erzeugt-merge-konflikt` (ein
Lösungs-Kandidat — Datei-Inhalt auf den sha256 des extrahierten Binaries
umstellen, `ADR-0044` Option B — ist mitgeführt, aber unentschieden,
Architect-Zuständigkeit) und
`BEO-PGC/lang-laufende-demo-umgebung-verpasst-neue-umgebungsvariable`
(kein Sensor deckt das bislang).

## Beobachtungs-Register (Zeiger)

Der Zähler steht in [`docs/plan/planning/observations/`](../observations/).
Kein Ausgang in dieser Welle fällig (beide neuen Einträge 1×, siehe oben).

## Folge-Slices

Keine — `ADR-0100`s Re-Evaluierungs-Trigger (echte Subjekt-scoped
NATS-Zugriffskontrolle, Stream-internes Replay/JetStream) ist nicht
eingetreten, bleibt `permanent` unbeobachtet fällig.

## Verifikation

- `docs/reviews/review-slice-nats-drittstream-core.md` (nach Fixrunde: 0
  offenes HIGH/MEDIUM; ursprünglich 4 HIGH, real behoben und
  re-verifiziert).
- `docs/reviews/review-slice-nats-drittstream-example-go.md` (nach
  Fixrunde: 0 HIGH, 0 MEDIUM, 1 LOW verbleibend nicht-blockierend;
  ursprünglich 1 HIGH, 1 MEDIUM, 2 LOW — alle vier adressiert).
- `docs/reviews/review-slice-nats-drittstream-example-csharp-kotlin.md`
  (0 HIGH, 0 MEDIUM, 1 LOW — LOW im selben Zug behoben).
- `make gates`: grün auf dem Endstand (728 Dateien, 0 Befunde;
  Coverage 82,80 %).
- `make test-integration`: grün auf dem Endstand, frischer Lauf — beide
  Belege aus `ADR-0100`s Fitness Function in einem Durchlauf (`RECEIVED
  change_id=981-1 table=feed_e2e_full operation=INSERT
  new_image={"id":"301","name":"NatsStreamE2ESentinel"}`,
  `REJECTED-NO-TOKEN`, `REJECTED-WRONG-TOKEN`, bestehendes Wecksignal
  unverändert funktionsfähig).
- Reale Ende-zu-Ende-Verifikation gegen die Demo-Umgebung, mehrfach
  unabhängig wiederholt (Implementer, Reviewer, Verifier — je mit
  eigenen Datensätzen, kein wiederholter Wert): Go-, C#- und
  Kotlin-Client empfangen je eine echte, unmittelbar zuvor eingefügte
  Change über `cdc.stream.>`.
