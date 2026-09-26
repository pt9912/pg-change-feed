Zustand: **verkörpert** — Ausgang: **verkörpert** → `AGENTS.md` §3.12, Absatz
„Verfasser einer ADR“ (der Architect nennt zu jeder Aussage über alle Werte einer
Menge die Menge, an der sie geprüft ist, und zu jeder Fitness-Function-Zeile den
Test, der sie trägt — erprobt oder als *hergeleitet* gekennzeichnet; eine in der Zeile
genannte Mutation nennt die Stellen, an denen sie erprobt ist, und die Instanz der
Messung), `.claude/agents/architect.md` (derselbe Satz an der Stelle, die der Zug beim
Schreiben liest) und `.harness/skills/reviewer.md` (Lese-Probe im Unterpunkt „Beleg trägt
seinen Satz nicht“) · seit welle-backfill-bestand, Satz zur Mutation seit
welle-transformationen.

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

**Schärfung zur Mutation (achtes Auftreten): verkörpert, ohne neue ADR.** Der Wortlaut des
Planners („erprobt oder hergeleitet“ plus Stellen und Instanz) trägt beide Fälle: `ADR-0126`
nannte eine Mutation, die weder erprobt noch als hergeleitet stand (der Satz verlangt Stellen
und Instanz, die es dort nicht gab); `ADR-0127` erprobte „in allen sieben Funktionen“ am
SQL-Text und schrieb „aus einer Funktion“ am Go-Test (der Satz verlangt die Menge „alle“ und
die Instanz „SQL-Text“, die Verallgemeinerung steht dann als hergeleitet). Gekürzt gegenüber
dem Vorschlag: „gesehene Farbe“ statt eines eigenen Satzes zur Erprobung, die zwei
Verallgemeinerungsrichtungen (alle → eine, Text → Test) als Beispiel im Satz. Keine neue ADR:
der Absatz „Verfasser einer ADR“ ist selbst ohne ADR in `AGENTS.md` §3.12 gewachsen;
`ADR-0083` (zwei Instanzen, Leser statt Sensor, Falsifikation ist die Messung) bleibt
unberührt, der Satz konkretisiert nur, was der Verfasser einer Mutationsangabe nennt. Ein
Sensor bleibt ausgeschlossen (Prosa über einen Test, `AGENTS.md` §3.12 §Grenze). Beleg-Anker:
`git grep -n "Stellen\*\*, an denen sie erprobt ist" -- AGENTS.md .claude/agents/architect.md`
liefert zwei Treffer.

Verwandt, nicht doppelt gezählt: `BEO-PGC/architect-verdikt-rollen-scope-luecke`.

Zähler: 8× (Dateien unter `evidence/`).
