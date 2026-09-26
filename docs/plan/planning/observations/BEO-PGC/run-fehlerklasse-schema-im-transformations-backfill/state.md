Zustand: **verkörpert** — Ausgang: **verkörpert** →
[`ADR-0117`](../../../../adr/0117-backfill-run-fehlerklasse-schema.md) (Klasse `schema`
für eine im Run nicht anwendbare Regel, run-lokal; Folgepflicht 7 von `ADR-0112` für
den Run ausgefüllt) · seit welle-backfill-bestand
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt`
§4.3 und §5 (g)).

Getragen: `classifyError` am Use Case des Runs bildet `ErrTransformationColumnMissing` und
`ErrTransformationTargetCollides` auf `schema` und `ErrTransformationStateChanged` auf
`configuration` ab; der Kommentar „vergibt `schema` nicht“ ist entfernt (`slice-transformationen-backfill-pfad`,
Start-Trigger: `ADR-0117` und das Architect-Verdikt `architect-verdict-backfill-schema-klasse-rollen`).
Beleg-Anker: `TestExecuteInapplicableRuleEndsRunAsSchema` und
`TestExecuteRuleStateChangeEndsRunAsConfiguration` (`make test`, je Abbildung eine Eingabeseiten-Mutation
rot) und die Phase „Backfill-Regelstand“ des Runners von `make test-integration`
(Run `failed`/`schema` ohne Change, Erfassungspfad läuft weiter, neuer Run nach dem Entfernen der Regel
`completed`; Zeile in `docs/user/e2e-abdeckung.md`). Offen bleibt der Betriebshinweis zur Abhilfe im Handbuch
(Folgepflicht 4 der ADR; Adresse: `slice-transformationen-betriebsdoku` §2, seit der Closure von
`slice-transformationen-backfill-pfad` mit dem Ablauf im Text).

Zähler: 1× (Datei unter `evidence/`).
