Zustand: **verkörpert** — Ausgang: **verkörpert** → `.harness/skills/reviewer.md`,
HIGH-Punkt „Kommentar trägt keine der Kommentar-Klassen“, Klausel *Skopus*: der
Punkt gilt für Code, Konfiguration und Skripte einschließlich Tests und Runner
(`tools/harness/*.sh`); ein Kommentar mit Vorher/Nachher-Andeutung oder einer
verworfenen Alternative im Konjunktiv gehört hierher, nicht unter INFO. Der Punkt
„Slice-/Wellen-Chronik in Produktionscode-Kommentar“ bleibt auf Produktionscode
begrenzt · seit welle-backfill-bestand
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt`
§4.1, R6a).

Zähler: 7× (Dateien unter `evidence/`; die siebte, `evidence/slice-transformationen-e2e-wirkung.md`,
trägt F-1 (HIGH, daher Datei): dieselbe Form im **Runner-Skript** (der Kommentar vor dem ersten
Neustart der Phase „Neustart und Ausschluss“ nannte den Zustand ohne die Zusage im Konjunktiv „trüge“); der
Reviewer fand sie vor dem Merge, Ausgang unverändert **verkörpert**. Der diff-skopierte
Kandidatenlauf von Schritt 20 druckte am Stand vor dem Review 0 Zeilen, weil „trüge“ nicht in seiner
Wortliste steht (gemessen bei der Closure, Befehl in der Beleg-Datei): die Wortliste ist eine
endliche Reihe von Konjunktiv-II-Formen, die Lücke ist benannt (Adresse: der Lese-Schritt der Closure
von `welle-transformationen`); die sechste, `evidence/slice-harness-guard-inplace-textwerkzeug.md`,
trägt zwei Ausprägungen (INFO/LOW, vor dem Merge von Verifier und Reviewer gefunden): Vorher-Nachher-Sätze in
Doku-Prosa von Plan und `MR-003` und ein Kandidatenlauf, der den Skopus der Skript-Kommentare nicht las; Ausgang
unverändert **verkörpert**; die fünfte, `evidence/slice-antragsqueue-lesefehler-failed.md`,
trägt F-1 (HIGH, daher Datei): dieselbe Form im **Produktionscode** (der Kommentar von
`sqlexec.ReadPendingRequests` nannte die verworfene Alternative im Konjunktiv, F-8 (INFO) dieselbe Form im
selben Block); der Reviewer fand sie vor dem Merge, Ausgang unverändert **verkörpert**. Der Name des
Eintrags ist der des Erstauftretens; der Mechanismus gilt für jede Datei des Skopus. **Die Schärfung des
Kandidatenlaufs ist verkörpert** → `.claude/commands/implement-slice.md` Schritt 20 (diff-skopierter
Konjunktiv-Kandidatenlauf über `'*.go' '*.sh' '*.awk'` mit `(//|#)` und den transliterierten Formen,
Beleg-Anker: `git grep -n 'wuerde' -- .claude/commands/implement-slice.md`) · der Lauf seit
`c625b572`, sein Skopus über Skripte und die transliterierten Formen seit
slice-harness-guard-inplace-textwerkzeug (Fixrunde, Review F-7). Die Messung des Kandidaten steht in der
Beleg-Datei der fünften: ein diff-skopierter Lauf auf `wäre|würde|hielte|hätte` traf am Diff dieses Slice
genau die zwei Defekte (2 Zeilen vor der Fixrunde, 0 danach), mit `statt|sonst` im Muster stiege er auf
13 Zeilen mit 2 Defekten. Die vierte, `evidence/slice-harness-fmt-check.md`,
trägt F-1 (HIGH, daher Datei): der Kommentar eines neuen Skripts und der Vertrag
begründeten eine Regel mit der verworfenen Alternative — der erste Fund nach der
Verkörperung; der Reviewer ordnete ihn nach der Skopus-Klausel als HIGH ein und fand ihn
vor dem Merge, Ausgang unverändert **verkörpert**).
