# ADR-Index — PG Change Feed

Regeln dieser Datei: Neue ADRs ergänzen diesen Index (Baseline-Regelwerk
`modul-04-adrs.md`; `AGENTS.md` §5). Eine Zeile je ADR-Datei; die Kennung
ist `ADR-<NNNN>` (MR-000), der Datei-Name trägt dieselbe Nummer. Status:
`Accepted`-ADRs sind inhaltlich immutable — Korrekturen als Folge-ADR mit
`Supersedes` (`AGENTS.md` §3.5), eine **Zitat-Korrektur** am Zitat- und
Verweisgerüst bei unverändertem Referenten ausgenommen
([`ADR-0073`](0073-zitat-korrektur-an-immutablen-dokumenten.md)).

Quelle der Erst-Anlage: `architecture-decision-records.md` (Architektur-
Entwurf, 2026-09-09 in Einzel-ADRs überführt).

| ID | Titel | Status | Datum | Datei |
|---|---|---|---|---|
| ADR-0001 | Hexagonale Architektur | Accepted | 2026-09-09 | [0001-hexagonale-architektur.md](0001-hexagonale-architektur.md) |
| ADR-0002 | Abhängigkeitsrichtung | Accepted | 2026-09-09 | [0002-abhaengigkeitsrichtung.md](0002-abhaengigkeitsrichtung.md) |
| ADR-0003 | Physische Modulgrenzen | Accepted | 2026-09-09 | [0003-physische-modulgrenzen.md](0003-physische-modulgrenzen.md) |
| ADR-0004 | Eigenes CDC-Domänenmodell | Accepted | 2026-09-09 | [0004-cdc-domainmodell.md](0004-cdc-domainmodell.md) |
| ADR-0005 | SourcePosition abstrahiert LSN | Accepted | 2026-09-09 | [0005-sourceposition-abstrahiert-lsn.md](0005-sourceposition-abstrahiert-lsn.md) |
| ADR-0006 | Replication Stream als Driving Adapter | Accepted | 2026-09-09 | [0006-replication-stream-driving-adapter.md](0006-replication-stream-driving-adapter.md) |
| ADR-0007 | Source ACK als Outbound Port | Accepted | 2026-09-09 | [0007-source-ack-outbound-port.md](0007-source-ack-outbound-port.md) |
| ADR-0008 | pgoutput als Standard | Accepted | 2026-09-09 | [0008-pgoutput-standard.md](0008-pgoutput-standard.md) |
| ADR-0009 | Change Store als Outbound Port | Accepted | 2026-09-09 | [0009-change-store-outbound-port.md](0009-change-store-outbound-port.md) |
| ADR-0010 | PostgreSQL als CDC Store | Accepted | 2026-09-09 | [0010-postgresql-cdc-store.md](0010-postgresql-cdc-store.md) |
| ADR-0011 | Persist-before-ACK | Accepted | 2026-09-09 | [0011-persist-before-ack.md](0011-persist-before-ack.md) |
| ADR-0012 | At-Least-Once | Accepted | 2026-09-09 | [0012-at-least-once.md](0012-at-least-once.md) |
| ADR-0013 | Consumer als Domänenkonzept | Accepted | 2026-09-09 | [0013-consumer-domainkonzept.md](0013-consumer-domainkonzept.md) |
| ADR-0014 | Retention als Domain Policy | Accepted | 2026-09-09 | [0014-retention-domain-policy.md](0014-retention-domain-policy.md) |
| ADR-0015 | Schema Evolution | Accepted | 2026-09-09 | [0015-schema-evolution.md](0015-schema-evolution.md) |
| ADR-0016 | JSONB Row Images im MVP | Accepted | 2026-09-09 | [0016-jsonb-row-images-mvp.md](0016-jsonb-row-images-mvp.md) |
| ADR-0017 | Generische Change-Tabelle | Accepted | 2026-09-09 | [0017-generische-change-tabelle.md](0017-generische-change-tabelle.md) |
| ADR-0018 | SQL als Driving Adapter (→ ADR-0046) | Superseded | 2026-09-09 | [0018-sql-driving-adapter.md](0018-sql-driving-adapter.md) |
| ADR-0019 | CLI als Driving Adapter | Accepted | 2026-09-09 | [0019-cli-driving-adapter.md](0019-cli-driving-adapter.md) |
| ADR-0020 | HTTP/gRPC optional (→ ADR-0057) | Superseded | 2026-09-09 | [0020-http-grpc-optional.md](0020-http-grpc-optional.md) |
| ADR-0021 | Large Transaction Buffer | Proposed | 2026-09-09 | [0021-large-transaction-buffer.md](0021-large-transaction-buffer.md) |
| ADR-0022 | Filesystem Spool als Driven Adapter | Proposed | 2026-09-09 | [0022-filesystem-spool-driven-adapter.md](0022-filesystem-spool-driven-adapter.md) |
| ADR-0023 | Fehlerklassifikation | Accepted | 2026-09-09 | [0023-fehlerklassifikation.md](0023-fehlerklassifikation.md) |
| ADR-0024 | Observability außerhalb der Domain | Accepted | 2026-09-09 | [0024-observability-ausserhalb-der-domain.md](0024-observability-ausserhalb-der-domain.md) |
| ADR-0025 | Domain Events | Accepted | 2026-09-09 | [0025-domain-events.md](0025-domain-events.md) |
| ADR-0026 | Composition Root | Accepted | 2026-09-09 | [0026-composition-root.md](0026-composition-root.md) |
| ADR-0027 | Capture Application Service | Accepted | 2026-09-09 | [0027-capture-application-service.md](0027-capture-application-service.md) |
| ADR-0028 | Inbound Use Cases | Accepted | 2026-09-09 | [0028-inbound-use-cases.md](0028-inbound-use-cases.md) |
| ADR-0029 | Domain-Invarianten | Accepted | 2026-09-09 | [0029-domain-invarianten.md](0029-domain-invarianten.md) |
| ADR-0030 | Testpyramide | Accepted | 2026-09-09 | [0030-testpyramide.md](0030-testpyramide.md) |
| ADR-0031 | Paketgrenzen | Accepted | 2026-09-09 | [0031-paketgrenzen.md](0031-paketgrenzen.md) |
| ADR-0032 | PostgreSQL bleibt Adapterdetail | Accepted | 2026-09-09 | [0032-postgresql-adapterdetail.md](0032-postgresql-adapterdetail.md) |
| ADR-0033 | Kein Event-Sourcing-Framework | Accepted | 2026-09-09 | [0033-kein-event-sourcing-framework.md](0033-kein-event-sourcing-framework.md) |
| ADR-0034 | Ports nach Fähigkeiten | Accepted | 2026-09-09 | [0034-ports-nach-faehigkeiten.md](0034-ports-nach-faehigkeiten.md) |
| ADR-0035 | Gesamtarchitektur | Accepted | 2026-09-09 | [0035-gesamtarchitektur.md](0035-gesamtarchitektur.md) |
| ADR-0036 | Architekturprüfung in CI (→ ADR-0041) | Superseded | 2026-09-09 | [0036-architekturpruefung-ci.md](0036-architekturpruefung-ci.md) |
| ADR-0037 | Implementierungsreihenfolge | Accepted | 2026-09-09 | [0037-implementierungsreihenfolge.md](0037-implementierungsreihenfolge.md) |
| ADR-0038 | Implementierungssprache Go (→ ADR-0039) | Superseded | 2026-09-09 | [0038-implementierungssprache-go.md](0038-implementierungssprache-go.md) |
| ADR-0039 | Paketstruktur-Detaillierung (Go) (→ ADR-0042) | Superseded | 2026-09-09 | [0039-paketstruktur-detaillierung-go.md](0039-paketstruktur-detaillierung-go.md) |
| ADR-0040 | ClockPort (Zeit als Outbound Port) | Accepted | 2026-09-09 | [0040-clockport.md](0040-clockport.md) |
| ADR-0041 | a-check als Maschinenform der Architektur-Prüfung (→ ADR-0068, teilweise) | Accepted | 2026-09-09 | [0041-a-check-maschinenform-architekturpruefung.md](0041-a-check-maschinenform-architekturpruefung.md) |
| ADR-0042 | Transport-Typen am Port | Accepted | 2026-09-09 | [0042-transport-typen-am-port.md](0042-transport-typen-am-port.md) |
| ADR-0043 | Schemamigrationen mit d-migrate | Accepted | 2026-09-09 | [0043-schemamigrationen-mit-d-migrate.md](0043-schemamigrationen-mit-d-migrate.md) |
| ADR-0044 | Image-Beleg-Semantik (Digest ist lauf-gebunden) | Accepted | 2026-09-09 | [0044-image-beleg-semantik.md](0044-image-beleg-semantik.md) |
| ADR-0045 | Commit-Traceability als Standing-Gate (→ ADR-0062, teilweise) | Accepted | 2026-09-10 | [0045-commit-traceability-standing-gate.md](0045-commit-traceability-standing-gate.md) |
| ADR-0046 | SQL-Driving-Adapter: Lese-Views direkt, Schreiben über Ports | Accepted | 2026-09-10 | [0046-sql-driving-adapter-lese-schreib-trennung.md](0046-sql-driving-adapter-lese-schreib-trennung.md) |
| ADR-0047 | Rollen-spezifische DSN-Verdrahtung (→ ADR-0048, ADR-0053, teilweise) | Accepted | 2026-09-12 | [0047-rollenspezifische-dsn-verdrahtung.md](0047-rollenspezifische-dsn-verdrahtung.md) |
| ADR-0048 | Heartbeat-Grant-Korrektur (SELECT-Ergänzung) | Accepted | 2026-09-12 | [0048-heartbeat-grant-korrektur-select-ergaenzung.md](0048-heartbeat-grant-korrektur-select-ergaenzung.md) |
| ADR-0049 | Replication-Fehlerklassen-Trennung und WAL-Rückstand-Schwellen | Accepted | 2026-09-12 | [0049-replication-fehlerklassen-schwellen.md](0049-replication-fehlerklassen-schwellen.md) |
| ADR-0050 | Schreibende SQL-Administration über Antrags-Queue + Live-Reload | Accepted | 2026-09-13 | [0050-sql-administration-antragsqueue-und-live-reload.md](0050-sql-administration-antragsqueue-und-live-reload.md) |
| ADR-0051 | CI/CD-Pipeline über GitHub Actions | Accepted | 2026-09-13 | [0051-cicd-pipeline-github-actions.md](0051-cicd-pipeline-github-actions.md) |
| ADR-0052 | Optionale YAML-Konfigurationsdatei ergänzt Umgebungsvariablen | Accepted | 2026-09-13 | [0052-optionale-yaml-konfigurationsdatei.md](0052-optionale-yaml-konfigurationsdatei.md) |
| ADR-0053 | Retention-Löschausführung bindet an `cdc_admin` — `DELETE`-Grant-Ergänzung | Accepted | 2026-09-13 | [0053-retention-loeschausfuehrung-cdc-admin-delete-grant.md](0053-retention-loeschausfuehrung-cdc-admin-delete-grant.md) |
| ADR-0054 | Coverage-Gate und Benchmark-Infrastruktur (→ ADR-0071, teilweise) | Accepted | 2026-09-13 | [0054-coverage-gate-und-benchmark-infrastruktur.md](0054-coverage-gate-und-benchmark-infrastruktur.md) |
| ADR-0055 | NATS-Change-Notification als Wecksignal (→ ADR-0056, teilweise) | Accepted | 2026-09-13 | [0055-nats-change-notification-wecksignal.md](0055-nats-change-notification-wecksignal.md) |
| ADR-0056 | NATS-Wecksignal — tabellen-granulares Subjekt | Accepted | 2026-09-13 | [0056-nats-tabellen-granulares-subjekt.md](0056-nats-tabellen-granulares-subjekt.md) |
| ADR-0057 | HTTP/JSON-API mit Token-Authn (Supersedes ADR-0020) | Accepted | 2026-09-14 | [0057-http-grpc-api.md](0057-http-grpc-api.md) |
| ADR-0058 | Testansatz für fünf Lastenheft-Kennungen ohne Testbeleg (→ ADR-0063/0064) | Accepted | 2026-09-14 | [0058-testansatz-fuenf-luecken.md](0058-testansatz-fuenf-luecken.md) |
| ADR-0059 | Spaltenauswahl — Mechanismus, Granularität und Wirkort (→ ADR-0065, teilweise) | Accepted | 2026-09-14 | [0059-spaltenauswahl-mechanismus.md](0059-spaltenauswahl-mechanismus.md) |
| ADR-0060 | gRPC-Server-Streaming für Live-Change-Zustellung (→ ADR-0066, teilweise) | Accepted | 2026-09-14 | [0060-grpc-streaming-mechanismus.md](0060-grpc-streaming-mechanismus.md) |
| ADR-0061 | HTTP/SSE zusätzlich zu gRPC für Live-Change-Zustellung | Accepted | 2026-09-14 | [0061-http-sse-zusaetzlich-zu-grpc.md](0061-http-sse-zusaetzlich-zu-grpc.md) |
| ADR-0062 | Lokaler commit-msg-Hook (Supersedes ADR-0045, teilweise) (→ ADR-0069) | Accepted | 2026-09-14 | [0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md](0062-lokaler-commit-msg-hook-ergaenzt-standing-gate.md) |
| ADR-0063 | Testform-Korrektur „Entfernte Spalten“ (Supersedes ADR-0058, teilweise) | Accepted | 2026-09-14 | [0063-lh-fa-sch-003-testform-korrektur.md](0063-lh-fa-sch-003-testform-korrektur.md) |
| ADR-0064 | `LH-QA-OPS-005`-Testansatz-Korrektur (Supersedes ADR-0058, teilweise) | Accepted | 2026-09-14 | [0064-lh-qa-ops-005-testansatz-korrektur.md](0064-lh-qa-ops-005-testansatz-korrektur.md) |
| ADR-0065 | Spaltenausschluss — dauerhafter Träger (Supersedes ADR-0059, teilweise) | Accepted | 2026-09-14 | [0065-spaltenausschluss-dauerhafter-traeger.md](0065-spaltenausschluss-dauerhafter-traeger.md) |
| ADR-0066 | Broadcaster — begrenzte Empfangswarteschlange (Supersedes ADR-0060; → ADR-0067) | Accepted | 2026-09-14 | [0066-broadcaster-begrenzte-empfangswarteschlange.md](0066-broadcaster-begrenzte-empfangswarteschlange.md) |
| ADR-0067 | Publish-Einbindung — Fitness-Function-Zeile korrigiert (Supersedes ADR-0066) | Accepted | 2026-09-14 | [0067-capture-publish-einbindung-fitness-function-korrektur.md](0067-capture-publish-einbindung-fitness-function-korrektur.md) |
| ADR-0068 | Wegwerf-Harness-Clients (Supersedes ADR-0041, teilweise) | Accepted | 2026-09-14 | [0068-wegwerf-clients-begrenzte-import-berechtigung.md](0068-wegwerf-clients-begrenzte-import-berechtigung.md) |
| ADR-0069 | commit-msg-Hook — einseitige Zusage (Supersedes ADR-0062; → ADR-0070, teilw.) | Accepted | 2026-09-15 | [0069-commit-msg-hook-einseitige-zusage.md](0069-commit-msg-hook-einseitige-zusage.md) |
| ADR-0070 | Supersede-Reichweite und Klassengrenze (Supersedes ADR-0069, teilweise) | Accepted | 2026-09-15 | [0070-supersede-reichweite-und-klassengrenze.md](0070-supersede-reichweite-und-klassengrenze.md) |
| ADR-0071 | Coverage-Gate — Messgegenstand netzlos prüfbare Fläche (Supersedes ADR-0054) | Accepted | 2026-09-15 | [0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md](0071-coverage-gate-messgegenstand-netzlos-pruefbare-flaeche.md) |
| ADR-0072 | hostpaths aktiviert — kein Host-Pfad, ohne Ausnahme (→ ADR-0074/0075, teilw.) | Accepted | 2026-09-15 | [0072-hostpaths-modul-aktiviert-ohne-ausnahme.md](0072-hostpaths-modul-aktiviert-ohne-ausnahme.md) |
| ADR-0073 | Zitat-Korrektur an immutablen Dokumenten — die Klasse für §3.5 | Accepted | 2026-09-15 | [0073-zitat-korrektur-an-immutablen-dokumenten.md](0073-zitat-korrektur-an-immutablen-dokumenten.md) |
| ADR-0074 | Zitationsform Schwester-Repo — Hausform (Supersedes ADR-0072, teilweise) | Accepted | 2026-09-15 | [0074-zitationsform-schwester-repo-hausform.md](0074-zitationsform-schwester-repo-hausform.md) |
| ADR-0075 | `hostpaths`-Regel — Reichweite, §3.11-Entwurf, Lokator-Disposition | Accepted | 2026-09-15 | [0075-hostpaths-reichweite-und-wortlaut.md](0075-hostpaths-reichweite-und-wortlaut.md) |
| ADR-0076 | Beispiel-Clients `examples/` (Supers. ADR-0060/0068, teilw.; → ADR-0079/0098) | Accepted | 2026-09-15 | [0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md](0076-beispiel-clients-examples-oeffentlicher-draht-vertrag.md) |
| ADR-0077 | Coverage-Rampen — Neu-Bemessung bei Subjekt-Transfer (→ ADR-0078, teilweise) | Accepted | 2026-09-15 | [0077-coverage-rampen-neu-bemessung-subjekt-transfer.md](0077-coverage-rampen-neu-bemessung-subjekt-transfer.md) |
| ADR-0078 | Coverage-Rampen — Transfer-Nachweis statt Summen-Konstanz | Accepted | 2026-09-15 | [0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md](0078-coverage-rampen-transfer-nachweis-statt-summen-konstanz.md) |
| ADR-0079 | NATS-Beispielclient — vierter `examples/`-Client (Supersedes ADR-0076, teilw.) | Accepted | 2026-09-15 | [0079-nats-beispielclient-vierter-examples-client.md](0079-nats-beispielclient-vierter-examples-client.md) |
| ADR-0080 | Nähte der pgconn-Adapter — Treiber-Hülle, kein Subjekt-Transfer | Accepted | 2026-09-15 | [0080-nahtform-pgconn-adapter-treiberhuelle.md](0080-nahtform-pgconn-adapter-treiberhuelle.md) |
| ADR-0081 | Changes-Lesen über die HTTP-API (Supers. ADR-0057, teilw.) | Accepted | 2026-09-15 | [0081-changes-lesen-ueber-die-http-api.md](0081-changes-lesen-ueber-die-http-api.md) |
| ADR-0082 | Coverage 80 % — Schnittmaß; Composition Root netzlos nicht prüfbar | Accepted | 2026-09-16 | [0082-coverage-schnittmass-composition-root-nicht-netzlos.md](0082-coverage-schnittmass-composition-root-nicht-netzlos.md) |
| ADR-0083 | Herkunft von Aussagen in Trägern — Zahlenwert und Tatsachenbehauptung | Accepted | 2026-09-16 | [0083-herkunft-von-aussagen-in-traegern.md](0083-herkunft-von-aussagen-in-traegern.md) |
| ADR-0084 | Sync-Gate nur für das Erzeugnis mit einer netzlosen, deterministischen Quelle | Accepted | 2026-09-16 | [0084-sync-gate-fuer-generierte-artefakte.md](0084-sync-gate-fuer-generierte-artefakte.md) |
| ADR-0085 | Build-Kontext-Ausnahme — `test-only` auf den Zweck (Supers. ADR-0082, teilw.) | Accepted | 2026-09-16 | [0085-build-kontext-ausnahme-test-only-zweck.md](0085-build-kontext-ausnahme-test-only-zweck.md) |
| ADR-0086 | Herkunft von Aussagen — Schwere folgt der Konsequenz (Supers. ADR-0083 teilw.) | Accepted | 2026-09-17 | [0086-herkunft-aussagen-schwere-folgt-der-konsequenz.md](0086-herkunft-aussagen-schwere-folgt-der-konsequenz.md) |
| ADR-0087 | Beispiel-Clients C#/Kotlin — Werkzeugkette (Supers. ADR-0076; → ADR-0090/0093) | Accepted | 2026-09-17 | [0087-beispiel-clients-csharp-kotlin.md](0087-beispiel-clients-csharp-kotlin.md) |
| ADR-0088 | Konfigurationsdatei — Feldmenge, Zugangsdaten-Klasse (Supers. ADR-0052, teilw.) | Accepted | 2026-09-17 | [0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md](0088-konfigurationsdatei-feldmenge-zugangsdaten-klasse.md) |
| ADR-0089 | Feldmengen-Paarung — kein Sensor, Wächter Review (Supers. ADR-0088; → ADR-0092) | Accepted | 2026-09-17 | [0089-feldmengen-paarung-kein-sensor-review-waechter.md](0089-feldmengen-paarung-kein-sensor-review-waechter.md) |
| ADR-0090 | Beispiel-Clients — volle Matrix (Supers. ADR-0087, teilw.) | Accepted | 2026-09-17 | [0090-beispiel-clients-volle-matrix.md](0090-beispiel-clients-volle-matrix.md) |
| ADR-0091 | Zugangsdaten-Klasse — sechs Schlüssel (Supers. ADR-0088; → ADR-0101) | Accepted | 2026-09-17 | [0091-zugangsdaten-klasse-sechs-schluessel.md](0091-zugangsdaten-klasse-sechs-schluessel.md) |
| ADR-0092 | Feldmengen-Paarung — Reichweite der drei Träger (Supers. ADR-0089; → ADR-0101) | Accepted | 2026-09-17 | [0092-feldmengen-paarung-reichweite-der-drei-traeger.md](0092-feldmengen-paarung-reichweite-der-drei-traeger.md) |
| ADR-0093 | Digest-Korrektur — ADR-0087s Kotlin-Basis-Image-Zeile (Supers. ADR-0087, teilw.) | Accepted | 2026-09-17 | [0093-digest-korrektur-adr-0087-kotlin-basis-image.md](0093-digest-korrektur-adr-0087-kotlin-basis-image.md) |
| ADR-0094 | Review-Matrixklasse — Kennung statt Adresse (ergänzt ADR-0073) | Accepted | 2026-09-18 | [0094-review-matrixklasse-kennung-statt-adresse.md](0094-review-matrixklasse-kennung-statt-adresse.md) |
| ADR-0095 | Review-Klasse — Status-Ausnahme (ergänzt ADR-0094) | Accepted | 2026-09-18 | [0095-review-klasse-exempt-status-check.md](0095-review-klasse-exempt-status-check.md) |
| ADR-0096 | Altbestand-Schlüssel für wellenlosen Archiv-Bestand | Accepted | 2026-09-18 | [0096-altbestand-schluessel-fuer-wellenlosen-archiv-bestand.md](0096-altbestand-schluessel-fuer-wellenlosen-archiv-bestand.md) |
| ADR-0097 | `observation`-Matrixklasse — Review verboten (ergänzt ADR-0094) | Accepted | 2026-09-18 | [0097-observation-matrixklasse-review-verboten.md](0097-observation-matrixklasse-review-verboten.md) |
| ADR-0098 | Beispiel-Clients — Startform `make`/Dockerfile (Supers. ADR-0076, teilw.) | Accepted | 2026-09-18 | [0098-beispiel-clients-start-ueber-make-dockerfile.md](0098-beispiel-clients-start-ueber-make-dockerfile.md) |
| ADR-0099 | `slice`/`welle → review` zurückgenommen (Supers. ADR-0097, teilw.) | Accepted | 2026-09-18 | [0099-slice-welle-review-regel-zurueckgenommen.md](0099-slice-welle-review-regel-zurueckgenommen.md) |
| ADR-0100 | NATS — dritter Vollinhalts-Zustellweg für Live-Streaming | Accepted | 2026-09-18 | [0100-nats-dritter-vollinhalts-zustellweg.md](0100-nats-dritter-vollinhalts-zustellweg.md) |
| ADR-0101 | Zugangsdaten-Klasse — sieben Schlüssel (Supers. ADR-0088/0091/0092, teilw.) | Accepted | 2026-09-18 | [0101-zugangsdaten-klasse-sieben-schluessel.md](0101-zugangsdaten-klasse-sieben-schluessel.md) |
| ADR-0102 | Zugangsdaten-Klasse — Supersede-Liste vervollständigt (Supers. ADR-0101, teilw.) | Accepted | 2026-09-18 | [0102-zugangsdaten-klasse-supersede-liste-vervollstaendigt.md](0102-zugangsdaten-klasse-supersede-liste-vervollstaendigt.md) |
| ADR-0103 | `image-hash.txt` lokal statt committet (Supers. ADR-0044, teilweise) | Accepted | 2026-09-18 | [0103-image-hash-lokal-statt-committet.md](0103-image-hash-lokal-statt-committet.md) |
| ADR-0104 | Benchmark-Schwellen PER-001/002/003 (Supers. ADR-0054, teilweise) | Accepted | 2026-09-19 | [0104-benchmark-schwellen-per-001-002-003.md](0104-benchmark-schwellen-per-001-002-003.md) |
| ADR-0105 | CI-Matrix-RTM-Sichtbarkeit POR-001/002 | Accepted | 2026-09-19 | [0105-ci-matrix-rtm-sichtbarkeit-por-001-002.md](0105-ci-matrix-rtm-sichtbarkeit-por-001-002.md) |
| ADR-0106 | C#/NuGet als erstes SDK-Package für `LH-FA-SST-009` | Accepted | 2026-09-19 | [0106-csharp-nuget-erstes-sdk-package.md](0106-csharp-nuget-erstes-sdk-package.md) |
| ADR-0107 | Python/PyPI als zweites SDK-Package für `LH-FA-SST-009` (→ ADR-0108/0110) | Accepted | 2026-09-19 | [0107-python-pypi-zweites-sdk-package.md](0107-python-pypi-zweites-sdk-package.md) |
| ADR-0108 | Python-SDK — `uv` statt `build`+`twine` (Supers. ADR-0107, teilw.) | Accepted | 2026-09-19 | [0108-python-sdk-uv-statt-build-twine.md](0108-python-sdk-uv-statt-build-twine.md) |
| ADR-0109 | Kotlin/GitHub Packages als drittes SDK-Package für `LH-FA-SST-009` | Accepted | 2026-09-20 | [0109-kotlin-github-packages-drittes-sdk-package.md](0109-kotlin-github-packages-drittes-sdk-package.md) |
| ADR-0110 | Python-SDK-Umfang erweitert auf gRPC/SSE/NATS (Supers. ADR-0107, teilw.) | Accepted | 2026-09-21 | [0110-python-sdk-umfang-erweitert-vollmatrix.md](0110-python-sdk-umfang-erweitert-vollmatrix.md) |
| ADR-0111 | Backfill des Bestands: Bulk-Copy im Slot-Snapshot (`LH-FA-CAP-009`) | Accepted | 2026-09-23 | [0111-backfill-bestand-snapshot-bulk-copy.md](0111-backfill-bestand-snapshot-bulk-copy.md) |
| ADR-0112 | Transformationen: deklarative Regeln vor der Persistierung (`LH-FA-CFG-007`) | Accepted | 2026-09-23 | [0112-transformationsform-deklarative-regeln-vor-persistenz.md](0112-transformationsform-deklarative-regeln-vor-persistenz.md) |
| ADR-0113 | Backfill: Rollenschnitt, Aufnahme, Warn-Kriterium (Supers. ADR-0111, teilw.) | Accepted | 2026-09-24 | [0113-backfill-rollenschnitt-aufnahme-warnkriterium.md](0113-backfill-rollenschnitt-aufnahme-warnkriterium.md) |
| ADR-0114 | Schema-Rollout: Vorlauf für View-Signatur-Änderungen (ergänzt ADR-0043) | Accepted | 2026-09-24 | [0114-schema-rollout-vorlauf-view-signatur.md](0114-schema-rollout-vorlauf-view-signatur.md) |
| ADR-0115 | Backfill: Spaltenwerte im Text-Ergebnisformat (Supers. ADR-0111, teilw.) | Accepted | 2026-09-24 | [0115-backfill-spaltenwerte-text-ergebnisformat.md](0115-backfill-spaltenwerte-text-ergebnisformat.md) |
| ADR-0116 | Backfill: Reichweite der Schema-Version-Referenz (Supers. ADR-0111, teilw.) | Accepted | 2026-09-24 | [0116-backfill-schema-version-referenz-reichweite.md](0116-backfill-schema-version-referenz-reichweite.md) |
| ADR-0117 | Backfill: Fehlerklasse schema im Run (Supers. ADR-0111, teilw.) | Accepted | 2026-09-24 | [0117-backfill-run-fehlerklasse-schema.md](0117-backfill-run-fehlerklasse-schema.md) |
| ADR-0118 | Backfill: Umschreiben im Snapshot-Fenster erkannt (Supers. ADR-0111, teilw.) | Accepted | 2026-09-24 | [0118-backfill-umschreiben-im-snapshot-fenster.md](0118-backfill-umschreiben-im-snapshot-fenster.md) |
| ADR-0119 | Backfill: Wirkung der Lesesperre berichtigt (Supers. ADR-0118, teilw.) | Accepted | 2026-09-25 | [0119-backfill-wirkung-der-lesesperre-berichtigt.md](0119-backfill-wirkung-der-lesesperre-berichtigt.md) |
| ADR-0120 | Capture: Slot bestätigt WAL ohne Inhalt für die Publication (ergänzt ADR-0007) | Accepted | 2026-09-25 | [0120-capture-slot-leerlauf-bestaetigung.md](0120-capture-slot-leerlauf-bestaetigung.md) |
| ADR-0121 | Capture: Bindung der Leerlauf-Bedingung berichtigt (Supers. ADR-0120, teilw.) | Accepted | 2026-09-25 | [0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md](0121-capture-leerlauf-bedingung-store-bindung-berichtigt.md) |
| ADR-0122 | Backfill: Tier der Replay-Invariante (Supers. ADR-0111, teilw.) | Accepted | 2026-09-25 | [0122-backfill-replay-invariante-e2e-tier.md](0122-backfill-replay-invariante-e2e-tier.md) |
| ADR-0123 | Kotlin-SDK zusätzlich auf Cloudsmith (Supers. ADR-0109, teilw.) | Accepted | 2026-09-25 | [0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md](0123-kotlin-sdk-zusaetzlich-auf-cloudsmith.md) |
| ADR-0124 | Retention: Kandidaten seitenweise ohne Row Images (ergänzt ADR-0111) | Accepted | 2026-09-25 | [0124-retention-kandidaten-seitenweise-ohne-row-images.md](0124-retention-kandidaten-seitenweise-ohne-row-images.md) |
| ADR-0125 | Transformationen: Regelform als json-Parameter (Supers. ADR-0112, teilw.) | Accepted | 2026-09-26 | [0125-transformationen-parametertyp-regelform-json.md](0125-transformationen-parametertyp-regelform-json.md) |
| ADR-0126 | Transformationen: Annahmemenge von rule_spec (Supers. ADR-0125, teilw.) | Accepted | 2026-09-26 | [0126-transformationen-annahmemenge-rule-spec.md](0126-transformationen-annahmemenge-rule-spec.md) |
