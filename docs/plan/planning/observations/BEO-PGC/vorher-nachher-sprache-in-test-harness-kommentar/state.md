Zustand: **verkörpert** — Ausgang: **verkörpert** → `.harness/skills/reviewer.md`,
HIGH-Punkt „Kommentar trägt keine der Kommentar-Klassen“, Klausel *Skopus*: der
Punkt gilt für Code, Konfiguration und Skripte einschließlich Tests und Runner
(`tools/harness/*.sh`); ein Kommentar mit Vorher/Nachher-Andeutung oder einer
verworfenen Alternative im Konjunktiv gehört hierher, nicht unter INFO. Der Punkt
„Slice-/Wellen-Chronik in Produktionscode-Kommentar“ bleibt auf Produktionscode
begrenzt · seit welle-backfill-bestand
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt`
§4.1, R6a).

Zähler: 5× (Dateien unter `evidence/`; die fünfte, `evidence/slice-antragsqueue-lesefehler-failed.md`,
trägt F-1 (HIGH, daher Datei): dieselbe Form im **Produktionscode** (der Kommentar von
`sqlexec.ReadPendingRequests` nannte die verworfene Alternative im Konjunktiv, F-8 (INFO) dieselbe Form im
selben Block); der Reviewer fand sie vor dem Merge, Ausgang unverändert **verkörpert**. Der Name des
Eintrags ist der des Erstauftretens; der Mechanismus gilt für jede Datei des Skopus. **Kandidat der
Schärfung (Adresse: der Lese-Schritt der Closure von `welle-transformationen`, Architect-Zug):**
`.claude/commands/implement-slice.md` Schritt 20 nennt den Konjunktiv als Frage der Probe, sein
Kandidatenlauf sucht aber nur Slice-/Wellen-Nummern; ein zweiter diff-skopierter Lauf auf
`wäre|würde|hielte|hätte` trifft am Diff dieses Slice genau die zwei Defekte (2 Zeilen vor der
Fixrunde, 0 danach; Messung in der Beleg-Datei), mit `statt|sonst` im Muster stiege er auf 13 Zeilen
mit 2 Defekten. Die vierte, `evidence/slice-harness-fmt-check.md`,
trägt F-1 (HIGH, daher Datei): der Kommentar eines neuen Skripts und der Vertrag
begründeten eine Regel mit der verworfenen Alternative — der erste Fund nach der
Verkörperung; der Reviewer ordnete ihn nach der Skopus-Klausel als HIGH ein und fand ihn
vor dem Merge, Ausgang unverändert **verkörpert**).
