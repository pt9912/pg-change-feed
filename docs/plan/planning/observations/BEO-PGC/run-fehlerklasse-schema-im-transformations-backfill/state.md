Zustand: **verkörpert** — Ausgang: **verkörpert** →
[`ADR-0117`](../../../../adr/0117-backfill-run-fehlerklasse-schema.md) (Klasse `schema`
für eine im Run nicht anwendbare Regel, run-lokal; Folgepflicht 7 von `ADR-0112` für
den Run ausgefüllt) · seit welle-backfill-bestand
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt`
§4.3 und §5 (g)).

Der Start-Trigger von
[`slice-transformationen-backfill-pfad`](../../../open/slice-transformationen-backfill-pfad.md)
nennt `ADR-0117` und das Architect-Verdikt `architect-verdict-backfill-schema-klasse-rollen`.
Die Abbildung, an der die ADR ansetzt, steht in `classifyError` am Use Case des Runs.

Zähler: 1× (Datei unter `evidence/`).
