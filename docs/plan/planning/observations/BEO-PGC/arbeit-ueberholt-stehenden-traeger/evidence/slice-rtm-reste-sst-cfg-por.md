# Beleg: slice-rtm-reste-sst-cfg-por

Vorgang: `slice-rtm-reste-sst-cfg-por` — Review-Finding F-2 korrigierte einen
falschen Komponentennamen in einem Code-Kommentar
(`tools/harness/run-integration-tests.sh`: `cdc.enable_table` schreibt
`cdc.administration_request`, nicht das nicht-existente `cdc.active_table`).

Fund: Zwei Träger derselben Tatsache überholten die Korrektur nicht mit:

1. Der Slice-Plan selbst (`docs/plan/planning/in-progress/rtm-reste-sst-cfg-por.md`
   §1 und §6) wiederholte an zwei Stellen wörtlich dieselbe falsche
   Formulierung wie der ursprüngliche Code-Kommentar — beide vom
   Implementer zum selben Zeitpunkt geschrieben, nur der Code-Kommentar
   wurde in der Fixrunde korrigiert, die Plan-Prosa nicht mitgezogen.
   Anders als die bisherigen Fälle dieses Registers (ein Doku-Träger
   überholt eine bewegte *Zahl*) ist der bewegte Gegenstand hier eine
   *Tatsachenbehauptung in Prosa*, die an drei Stellen wortgleich
   dupliziert stand — der §3.13-Suchlauf (`grep` über Träger) hätte den
   exakten Fehltext (`cdc.active_table`) real gefunden, wurde aber vom
   Implementer nach der Fixrunde nicht erneut ausgeführt.
2. `docs/user/e2e-abdeckung.md` (generiert) war an zwei Zeilen-Lokatoren
   um exakt 3 Zeilen veraltet — die Fixrunde verlängerte den
   Code-Kommentar um 3 Zeilen, ohne die generierte Datei danach neu zu
   erzeugen (`make test-integration`). Dies ist die von `AGENTS.md` §3.13
   §Grenze explizit benannte Zahlen-Lücke: ein Zeilen-Lokator verschiebt
   sich mit jeder Nachbaränderung, ohne im Text eine wiederholbare Spur zu
   hinterlassen.

Beide Funde stammen nicht vom Reviewer (der die Fixrunde selbst prüfte,
aber nicht auf die Slice-Plan-Prosa bzw. die generierte Datei
zurückschaute), sondern vom unabhängigen Verifier — über ein reales
`git worktree`-Isolationsexperiment plus eigene Nachrechnung der
`awk`-Anker-Logik des Abdeckungsskripts, nicht über Diff-Lesen allein.

Quelle: `docs/reviews/verifikation-slice-rtm-reste-sst-cfg-por.md` <!-- d-check:status-provenance -->
· Fixrunden-Commit `89e9922f` zu
`docs/plan/planning/in-progress/rtm-reste-sst-cfg-por.md`.
