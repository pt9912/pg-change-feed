# ADR-Index — PG Change Feed

Regeln dieser Datei: Neue ADRs ergänzen diesen Index (Baseline-Regelwerk
`modul-04-adrs.md`; `AGENTS.md` §5). Eine Zeile je ADR-Datei; die Kennung
ist `ADR-<NNNN>` (MR-000), der Datei-Name trägt dieselbe Nummer. Status:
`Accepted`-ADRs sind inhaltlich immutable — Korrekturen als Folge-ADR mit
`Supersedes` (`AGENTS.md` §3.5).

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
| ADR-0020 | HTTP/gRPC optional | Accepted | 2026-09-09 | [0020-http-grpc-optional.md](0020-http-grpc-optional.md) |
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
| ADR-0041 | a-check als Maschinenform der Architektur-Prüfung | Accepted | 2026-09-09 | [0041-a-check-maschinenform-architekturpruefung.md](0041-a-check-maschinenform-architekturpruefung.md) |
| ADR-0042 | Transport-Typen am Port | Accepted | 2026-09-09 | [0042-transport-typen-am-port.md](0042-transport-typen-am-port.md) |
| ADR-0043 | Schemamigrationen mit d-migrate | Accepted | 2026-09-09 | [0043-schemamigrationen-mit-d-migrate.md](0043-schemamigrationen-mit-d-migrate.md) |
| ADR-0044 | Image-Beleg-Semantik (Digest ist lauf-gebunden) | Accepted | 2026-09-09 | [0044-image-beleg-semantik.md](0044-image-beleg-semantik.md) |
| ADR-0045 | Commit-Traceability als Standing-Gate | Accepted | 2026-09-10 | [0045-commit-traceability-standing-gate.md](0045-commit-traceability-standing-gate.md) |
| ADR-0046 | SQL-Driving-Adapter: Lese-Views direkt, Schreiben über Ports | Accepted | 2026-09-10 | [0046-sql-driving-adapter-lese-schreib-trennung.md](0046-sql-driving-adapter-lese-schreib-trennung.md) |
