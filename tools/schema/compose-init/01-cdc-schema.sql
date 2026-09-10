-- Compose-Init der Test-Instanz: das Ziel-Schema der CDC-Tabellen und die
-- Schema-Auflösung der Rollout-Verbindung. Die Tabellen-DDL entsteht über
-- den d-migrate-Rollout (ADR-0043) und nicht hier — dieser Lauf trägt nur
-- die Umgebungs-Vorbedingung: das leere Schema `cdc` (SPEC-006) und den
-- search_path der Verbindungen als postgres in dieser Datenbank, damit
-- current_schema() der Rollout- und Test-Verbindungen das CDC-Schema liest.
CREATE SCHEMA IF NOT EXISTS cdc;
ALTER ROLE postgres IN DATABASE cdc SET search_path = cdc;
