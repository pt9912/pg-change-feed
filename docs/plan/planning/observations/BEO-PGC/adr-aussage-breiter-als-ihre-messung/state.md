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

- `ADR-0124` (Fitness-Function-Zeile 3: „die Spitze nach dem Run liegt im Bereich der
  Spitze im Run (8,9 bis 10,5 MiB) plus dem Bedarf einer Seite“; gemessen 14,9 bis
  17,6 MiB): **akzeptiertes Negativ**, kein Supersede. Die Zeile ist eine Messung ohne
  Pass/Fail, keine der sechs Festlegungen hängt an der Zahl; die Lesart (der Form nach
  erfüllt) steht in der Closure-Notiz von `slice-retention-lauf-speicher-begrenzung`
  und im Messbericht dieses Slice, Abschnitt 7.1.

- `ADR-0125` Festlegung 1 („alles, was gültiges JSON ist, wird angenommen“; gemessen:
  `{"a":"\u0000"}` wird abgelehnt): berichtigt mit `ADR-0126` (Teil-Supersedes), der
  `SPEC-019`-Absatz in place im selben Commit.
- `ADR-0126` Fitness-Function-Zeile (Mutation „Cast `::jsonb` entfernen färbt die
  `\u0000`-Zeilen rot“; gemessen: kein Rot, der Zuweisungs-Cast `json` → `jsonb` leistet
  dieselbe Umwandlung): **akzeptiertes Negativ**, kein Supersede. Die Zeile hat keinen
  Verbraucher, der Store-Test bindet die Annahmemenge mit der wirksamen Mutation
  `NULL::jsonb`, die Abweichung steht im Plan von
  `slice-transformationen-antragsweg-schema` §3; die Berichtigung trägt die nächste ADR
  zu `rule_spec` als Klausel (Trigger: die Re-Evaluierungs-Bedingungen von `ADR-0126`).

**Schärfung, Vorschlag an den Architect (Verkörperung 3b, Modul 8; Adresse: der
Lese-Schritt der Closure von `welle-transformationen`).** Zielort `AGENTS.md` §3.12, Absatz
„Verfasser einer ADR“, ein Satz mehr: „Eine in einer Fitness-Function-Zeile genannte
Mutation steht erprobt — an derselben Instanz wie die Messung — oder als hergeleitet
gekennzeichnet; ‚der Implementer fährt sie‘ ist eine Erwartung, keine Erprobung.“ Die
Lücke im heutigen Wortlaut: der Absatz nennt die Menge und den Test, der eine
Fitness-Function-Zeile trägt, nicht die Mutation, die den Test rot färben soll; das siebte
Auftreten entstand in der Rolle, die der Absatz adressiert, und stand formal als Erwartung
(Klausel „wird so sein“), die der Implementer erst am Bau erprobte. Kosten der Schärfung:
ein Satz; Kosten ihres Fehlens: eine falsche Angabe in einer `Accepted`-ADR, deren
Berichtigung nur über eine Folge-ADR geht (§3.5). Eine Gegenentscheidung („akzeptiertes
Negativ, Schwere LOW, kein Verbraucher“) ist ein Verdikt des Architects.

Verwandt, nicht doppelt gezählt: `BEO-PGC/architect-verdikt-rollen-scope-luecke`.

Zähler: 7× (Dateien unter `evidence/`).
