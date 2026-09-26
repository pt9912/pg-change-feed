Zustand: **gestrichen** — Punkt 1 geschlossen, Punkt 2 gestrichen mit Begründung
· seit welle-backfill-bestand
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt`
§4.1 und §4.2).

- **Punkt 1** (Rollen-Rechte der Rollout-Datei): geschlossen. Der Test
  `TestRolloutDateiTraegtDieRechteDerVerdrahtung` liest die echte
  `tools/schema/nacharbeit-roles.sql` und bindet sieben Zusagen an ihren Inhalt; eine
  Mutation in der Datei färbt Test und Gate. Der Integrationslauf liest den
  Heartbeat-Grant real. Die Rollen-Tests der Backfill-Adapter und des
  Administrationswegs laufen unter Login-Identitäten der Rolle
  (`backfillroles_test.go`, `administration_roles_internal_test.go`).
- **Punkt 2** (Replikations- und ACK-Adapter gegen Rollen-Vertauschung): gestrichen.
  Beide Adapter sind an `cdc_capture` gebunden, dessen Login als einziger das
  Attribut `REPLICATION` trägt; der Compose-E2E fährt den Feed-Container mit den drei
  Rollen-DSNs und belegt an der PostgreSQL-Server-Ebene, dass ein Login ohne das
  Attribut an einer Replikationsverbindung scheitert (`harness/README.md`
  §Sensors, `make test-integration`). Vertauschte DSNs ließen den Container am Start
  scheitern — **hergeleitet**, nicht als Mutation gemessen.

- **Punkt 3** (Schreibweisen der Grants auf Funktionen in `tools/schema/nacharbeit-administration.sql`):
  geschlossen mit `slice-transformationen-antragsweg-schema`. `TestAdministrationDateiTraegtDieFunktionsRechte`
  zählt jede `GRANT`-Anweisung und jede `CREATE FUNCTION` der Datei und verlangt die engen Formen
  (`GRANT EXECUTE ON FUNCTION … TO …`, `CREATE FUNCTION cdc.<name>`); `GRANT ALL`,
  `GRANT … ON ALL FUNCTIONS IN SCHEMA` und eine Funktion außerhalb von `cdc.<name>` färben ihn rot.
  **Benannte Grenze:** ein dynamisch zusammengesetzter Grant (`EXECUTE format(…)`,
  `ALTER DEFAULT PRIVILEGES`) bleibt für den Textleser ungelesen; die Realinstanz-Belege
  (`TestAdministrationRequestTransformationFunctionsRequireCdcAdminMembership`, Lauf 5 des
  Guard-Skripts) tragen ihn für die zwei Transformations-Funktionen, eine spätere Funktion
  trägt ihn nicht von selbst.

Zähler: 4× (Dateien unter `evidence/`).
