**Vorgang:** slice-antragsqueue-lesefehler-failed (Review F-2, HIGH; Verifikation V-1, LOW)

**Fund:** Zwei Verweise im Slice-Plan nannten die falsche Stelle, beide vor dem Merge gefunden.

- **Review F-2 (HIGH):** die Träger-Tabelle stützte „keine Aussage zu einer Ablehnung oder einem Anhalten der Queue“ im Handbuch auf „Suchlauf Zeile 7 ohne Treffer“. Zeile 7 des `suchlauf`-Feldes zählt die Fundstellen von `NewAdministrationRequest` (`diff 7`, sieben Treffer); die Messung über `docs/user` und `harness` mit Ergebnis 0 stand in den Zeilen 11 und 12. Die Aussage war wahr, der genannte Beleg trug sie nicht; die Zeilennummern des Feldes verschoben sich, als Zeilen ergänzt wurden (übernommen aus dem Review). Die Fixrunde nannte die Zeilen 11 und 12 und die Verifikation las alle „Zeile N“-Verweise des Plans gegen den Block.
- **Verifikation V-1 (LOW):** der Fixrunden-Absatz des Plans zählte die Findings des Reviews falsch: „F-7 (Stand ‚Parent‘ ist `47a646b5` …)“ meint im Review F-9 (F-7 sind die benannten Grenzen der Verarbeitung), und „F-9, F-10, F-11 — keine Aktion“ verdeckte, dass F-9 einen Nachzug trug. Die Behebung selbst stand vollständig im Plan, nur die Zuordnung Finding zu Behandlung war verschoben; die Closure setzte die Nummern des Reviews ein.

**Form (Ausprägung):** dieselbe Klasse (der Verweis ist die Prüf-Form einer Aussage, er wird aufgeschlagen, nicht nachgemessen), in einer **vierten Form**: ein Verweis **innerhalb desselben Plans** — die Zeile eines nummerierten Feldes und die Nummer eines Findings des zugehörigen Reports —, nicht auf ein anderes Dokument. Die Nummer stammt in beiden Fällen aus einer Liste, die sich bewegt hat (Feld-Zeilen, Finding-Nummern des Reviews); der Ursprung der Verschiebung ist nicht ermittelt.

Quelle: `docs/reviews/review-slice-antragsqueue-lesefehler-failed.md` (F-2) <!-- d-check:status-provenance -->
· `docs/reviews/verifikation-slice-antragsqueue-lesefehler-failed.md` (§7 Zeile F-2, §8 V-1). <!-- d-check:status-provenance -->
