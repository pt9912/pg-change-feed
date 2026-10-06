**Vorgang:** slice-harness-baseline-v6-14-1 (Re-Review F-1, MEDIUM; Review F-5, LOW; dreizehnte Datei des Eintrags)

**Fund:** `ADR-0157` Entscheidung 1 Bedingung (b) bindet „Referent gemessen derselbe“ an einen `cmp` des Zielinhalts; Bedingung (a) lässt aber auch Anker und Zeilen-Lokator zu, und die ADR legt die verglichene Einheit nicht fest. Gemessen ist nur der `cmp` ganzer Dateien: Ein Anker-Wechsel in derselben Datei besteht ihn immer, ein verschobener Zeilen-Lokator fällt durch (Re-Review F-1). Im selben Vorgang sagt `ADR-0156` „der Referent ist je Zeile `cmp`-gleich“, gemessen war normalisiert; für `regelwerk/modul-13-quality-gates.md` ist der rohe `cmp` 1 (Review F-5).

**Form (Ausprägung):** bekannte Form — ADR-Prosa-Aussage über eine Messung, deren Einheit oder Instanz sie nicht nennt. F-5 zählt im selben Vorgang nicht ein zweites Mal. Schwere MEDIUM, daher eine Datei trotz Deckel; vor dem Merge gefunden. Adresse für F-1: Folge-Slice `slice-zitat-korrektur-vergleichseinheit` (Architect-Zug); F-5 bleibt Befund ohne Änderung (eine Ergänzung änderte die Aussage einer `Accepted`-ADR). Ausgang unverändert **verkörpert**.

Quelle: `docs/reviews/review-slice-harness-baseline-v6-14-1-fixrunde.md` (F-1) <!-- d-check:status-provenance -->
· `docs/reviews/review-slice-harness-baseline-v6-14-1.md` (F-5). <!-- d-check:status-provenance -->
