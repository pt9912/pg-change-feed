Zustand: **verkörpert** — Ausgang: **verkörpert** → `AGENTS.md` §3.13 §Suchform
(Suchraum ganzer Baum, Muster aus drei Arten, Befehl im Codeblock, **Frist der
Meldung** in einer fremden Datei: die Closure des meldenden Slice) und
`.harness/skills/reviewer.md`, MEDIUM-Punkt „Nachzug widerspricht dem Nachbarn im
selben Träger“ (Probe: den Kontext um jede hinzugefügte Zeile lesen,
`git diff -U20`) · seit welle-backfill-bestand
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt`
§3.3, R1 und R2).

Begründung: die Belege sind dieselbe Fehlerklasse (ein Nachzug lässt eine
überholte Aussage stehen), von den Findern uneinheitlich als MEDIUM, LOW und INFO
eingeordnet; der Reviewer-Punkt ordnet sie einheitlich ein. Ein Sensor ist
ausgeschlossen: ob zwei Aussagen einander widersprechen, ist eine Lese-Handlung.

Zähler: 16× (Dateien unter `evidence/`; die sechzehnte,
`evidence/slice-routing-spec-nachzug.md`, trägt F-1 (MEDIUM): ein hinzugefügter Absatz in
`LH-FA-CAP-009.a` ließ den unberührten Nachbarabsatz „Fail-closed vor dem Commit“
unvollständig; vor dem Merge gefunden, in der Fixrunde behoben, Ausgang unverändert
**verkörpert**; die fünfzehnte,
`evidence/slice-transformationen-betriebsdoku.md`, trägt F-2 (MEDIUM) und eine **neue Form**:
keine widersprechende Nachbar-Aussage, sondern eine gegen eine im Plan explizit gezählte Liste
(drei benannte Beispiele einer Übergabe) unvollständige Übernahme in den Zielträger (Handbuch);
vor dem Merge gefunden, in der Fixrunde behoben, Ausgang unverändert **verkörpert**; die
vierzehnte,
`evidence/slice-capture-leerlauf-quellbelege.md`, trägt F-1 (MEDIUM) und V-1 (MEDIUM): der Nachzug
des Umbaus stand in §3, der Wortlaut von DoD 1 in §2 desselben Plans beschrieb den früheren
Aufbau; vor dem Merge gefunden, Ausgang unverändert **verkörpert**; die dreizehnte,
`evidence/slice-harness-suchlauf-nachmessen.md`, trägt F-3 (MEDIUM): der Plan hielt seinen alten
Wortlaut zum Ort der Rücknahme neben der Abweichungs-Zeile stehen; die zwölfte,
`evidence/slice-sdk-kotlin-cloudsmith.md`, trägt F-2 (MEDIUM) und F-3 (LOW): der Schlussabsatz
eines Betreiber-Abschnitts und die Nachbarzeile zum Werttyp eines Secrets im selben Träger; die zehnte,
`evidence/slice-backfill-speicher-untersuchung.md`, trägt F-3, F-5 und V-2: Nachbar-Aussagen
im selben Handbuch und im Trigger des Folge-Slice, ein Suchausdruck in fremder Datei, dessen
Meldung die Closure des meldenden Slice zieht; die elfte,
`evidence/slice-retention-lauf-speicher-begrenzung.md`, trägt F-3 und V-1, beide MEDIUM: die
Nachzüge eines Architect-Verdikts in Plan §6 und Handbuch).

Deckel bei 10× (seit slice-backfill-speicher-untersuchung): weitere Auftreten, die vor dem
Merge vom Reviewer oder Verifier gefunden werden, Schwere ≤ LOW haben und einen bekannten
Träger-Typ treffen, bekommen keine `evidence/`-Datei, sondern stehen mit Finding-Kennung in
der Closure-Notiz des Slice (`../../README.md`, Deckel für verkörperte Einträge ab 10×).
