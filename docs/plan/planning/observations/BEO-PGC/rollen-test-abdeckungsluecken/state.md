Zustand: offen — Ausgang: **weiter offen** → ein Test, der den
tatsächlichen `nacharbeit-roles.sql`-Inhalt liest (Punkt 1), sowie eine
Rollen-Vertauschungsprüfung, die real die Replication-Stream-/ACK-Adapter
(`receive.NewStream`/`postgresack.New`) mit rollenbeschränkten
Verbindungen aufruft (Punkt 2 — `slice-028` schloss nur die
PostgreSQL-Server-Ebene, nicht die Adapter-Ebene), sind eigenständiger
Aufwand; kein Slice dafür existiert.
Zähler (abgeleitet): 2× (evidence/slice-023.md, evidence/slice-028.md).
