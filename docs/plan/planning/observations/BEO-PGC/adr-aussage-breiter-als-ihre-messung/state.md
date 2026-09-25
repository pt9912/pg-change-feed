Zustand: **verkörpert** — Ausgang: **verkörpert** → `AGENTS.md` §3.12, Absatz
„Verfasser einer ADR“ (der Architect nennt zu jeder Aussage über alle Werte einer
Menge die Menge, an der sie geprüft ist, und zu jeder Fitness-Function-Zeile den
Test, der sie trägt — erprobt oder als *hergeleitet* gekennzeichnet) und
`.claude/agents/architect.md` (Zeiger) · seit welle-backfill-bestand
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt`
§4.1, R5).

Die Belege und ihre Berichtigungen:

- `ADR-0114` Entscheidung 3 (Rechte-Verlust bei `DROP VIEW`) und `ADR-0113`
  Festlegung 1 (Recht auf `cdc.administration_request`): **akzeptiertes
  Negativ**, kein Supersede. `ADR-0114` Entscheidung 3 bleibt für jede Rolle wahr,
  die `tools/schema/nacharbeit-roles.sql` vergibt; die Grenze für Rechte außerhalb
  des Repos steht an drei Trägern (Handbuch §4, `harness/targets/schema-rollout.md`
  §Grenze Punkt 3, Plan von `slice-backfill-change-origin`). Das Recht von
  `cdc_admin` auf `cdc.administration_request` steht im Pflichtenheft (Rang 2, vor
  jeder ADR) und in `nacharbeit-roles.sql`.
- `ADR-0111` (byte-gleich): berichtigt mit `ADR-0115`.
- `ADR-0118` (Reichweite der Sperre, `RENAME COLUMN`): berichtigt mit `ADR-0119`.
- `ADR-0120` (Fitness-Function-Zeile): berichtigt mit `ADR-0121`.
- `ADR-0111` (Fitness-Function-Zeile der Replay-Invariante nennt
  `make test-replication`, der Beleg liegt in `make test-integration`): berichtigt
  mit `ADR-0122`.

Verwandt, nicht doppelt gezählt: `BEO-PGC/architect-verdikt-rollen-scope-luecke`.

Zähler: 5× (Dateien unter `evidence/`).
