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

- `ADR-0127` Fitness-Function-Zeile 1 (Mutation „`clock_timestamp()` aus **einer** Funktion
  entfernen“ färbt den Go-Test rot; erprobt war der Text mit `now()` in **allen** sieben
  Funktionen; gemessen färbt sie bei `remove_transformation` und `backfill_table` keinen Test,
  weil die Farbe an der Stellung der Funktion in der Aufruffolge hängt): **akzeptiertes
  Negativ**, kein Supersede. Die Zeile hat keinen Verbraucher, das Produkt hält die Ordnung
  (Zusatztest der Verifikation), die Fixrunde 3 des Slice bindet alle sieben Funktionen in
  beiden Richtungen; die Berichtigung steht im Plan von `slice-transformationen-antragsweg-usecase`
  §3 (Mutationstabelle, Zeilen „Fixrunde 2“ und „Fixrunde 3 (V-2)“) und die nächste ADR zu
  `requested_at` oder zur Antrags-Queue trägt sie als Klausel (Trigger: die
  Re-Evaluierungs-Bedingungen von `ADR-0127`).

**Schärfung, Vorschlag an den Architect (Verkörperung 3b, Modul 8: Planner → Architect →
Planner; der Planner schärft `AGENTS.md` nicht selbst — Regel-Verkörperung ist eine
Entscheidung, keine Planung). Adresse: der nächste Architect-Zug zu einer ADR mit
Fitness-Function-Zeile — spätestens der Lese-Schritt der Closure von `welle-transformationen`.**
Zielort `AGENTS.md` §3.12, Absatz „Verfasser einer ADR“, ein Satz mehr, in der Fassung nach dem
achten Auftreten: „Eine in einer Fitness-Function-Zeile genannte Mutation nennt die Menge der
Stellen, an denen sie erprobt ist (eine · alle · welche), und die Instanz der Messung; jede
Verallgemeinerung darüber hinaus steht als hergeleitet, und ‚der Implementer fährt sie‘ ist eine
Erwartung, keine Erprobung.“ Die Lücke im heutigen Wortlaut: der Absatz nennt die Menge und den
Test, der eine Fitness-Function-Zeile trägt, nicht die Mutation, die den Test rot färben soll.
Das siebte Auftreten (`ADR-0126`) stand formal als Erwartung; das achte (`ADR-0127`) trug die
Grenze des Erprobten richtig (erprobt an allen sieben Funktionen, Go-Test nicht gefahren) und war
breiter in der Verallgemeinerung auf „eine Funktion“ — die erste Fassung des Satzes („erprobt
oder hergeleitet“) hätte die Zeile bestanden, die Ergänzung „Menge der Stellen“ trifft sie. Beide
Auftreten entstanden im Zug des Architects; der Vorschlag stand zur Zeit des zweiten im Register
und nicht in `AGENTS.md` §3.12 (gemessen: der Absatz „Verfasser einer ADR“ trägt den Satz nicht;
ob der Zug das Register las, ist nicht belegt). Ein Satz im Text, den jeder Zug liest, wirkt
sicherer als eine Adresse im Register — deshalb ist die Adresse der nächste Zug selbst (der
Orchestrator nennt den Vorschlag im Auftrag) und nicht erst die Closure. Kosten der Schärfung:
ein Satz; Kosten ihres Fehlens: eine falsche Angabe in einer `Accepted`-ADR je Auftreten, deren
Berichtigung nur über eine Folge-ADR geht (§3.5), und ein LOW-Finding samt Register-Zeile.
Gegenentscheidung („akzeptiertes Negativ, Schwere LOW, kein Verbraucher, der Implementer-Lauf fängt
die Zeile“) ist ein Verdikt des Architects; **ein Sensor ist ausgeschlossen** (Prosa über einen
Test, `AGENTS.md` §3.12 §Grenze).

Verwandt, nicht doppelt gezählt: `BEO-PGC/architect-verdikt-rollen-scope-luecke`.

Zähler: 8× (Dateien unter `evidence/`).
