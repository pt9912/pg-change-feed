Zustand: offen — Ausgang: **weiter offen**, adressiert. Adressen: `slice-backfill-run-store`
(DoD „Adapter-Pflichten": der Store-Test `Finish(failed)` auf einen `completed`-Run liefert nil
und die Zeile bleibt `completed`) und `slice-backfill-sql-administration` (DoD
„Verarbeitung": der Aufrufer liest `result.Run.Status` als Zustand des Aufrufs, nicht der
Zeile). Ob das Ergebnis die Zeile zurücklesen soll, ist offen und Sache dieser beiden
Slices; „entfallen" erst mit dem Store-Test.

Zähler (abgeleitet): 1× (evidence/slice-backfill-run-usecase.md).
