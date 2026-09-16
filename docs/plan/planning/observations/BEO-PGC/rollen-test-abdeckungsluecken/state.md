Zustand: offen — Ausgang: **weiter offen** → ein Test, der den
tatsächlichen `nacharbeit-roles.sql`-Inhalt liest (Punkt 1), sowie eine
Rollen-Vertauschungsprüfung, die real die Replication-Stream-/ACK-Adapter
(`receive.NewStream`/`postgresack.New`) mit rollenbeschränkten
Verbindungen aufruft (Punkt 2 — `slice-028` schloss nur die
PostgreSQL-Server-Ebene, nicht die Adapter-Ebene), sind eigenständiger
Aufwand; kein Slice dafür existiert.
Zähler (abgeleitet): 2× (evidence/slice-023.md, evidence/slice-028.md).

**Punkt 1 ist mit `slice-093` geschlossen — und der Zähler bleibt bei 2×.**
Der Slice hat die Lücke **behoben**, statt sie erneut zu treffen: der neue Test
`TestRolloutDateiTraegtDieRechteDerVerdrahtung` liest die **echte**
`tools/schema/nacharbeit-roles.sql` und bindet sieben Zusagen an ihren Inhalt;
eine Mutation in der Datei färbt **Test und Gate** (Bau-EC 1).
**Präzise, mit der Verifikation (F-4):** geschlossen ist die Lücke **im
Gate-Bündel** — nicht „nirgends"; der Integrationslauf liest den Heartbeat-Grant
real (er rollt die Datei aus und verlangt den Compose-Status `healthy`).
**Punkt 2 bleibt offen** (Rollen-Vertauschung für Replication-Stream/ACK auf
Adapter-Ebene; sie braucht rollenbeschränkte Logins in
`tools/harness/run-replication-tests.sh`, unberührt).
Ein **Beleg** wird dafür nicht angelegt: der Eintrag wird nicht erneut getroffen,
sondern aufgelöst — und ein Vorgang, der eine Lücke schließt, ist kein Vorkommen.
Sein Ausgang gehört dem **Lese-Schritt der `welle-20`-Closure** (Modul 6).
