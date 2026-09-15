# Beleg: slice-082

Vorgang: `slice-082` — der Fixture-Nachzug und die neue gemeinsame
Schema-Anwendung (`tools/schema/apply-rollout.sh`, die `make schema-rollout`
aufruft).

Fund: **Jeder** Lauf von `make schema-rollout` verändert die committete
`tools/schema/plan.yaml` — der Pflicht-Report trägt das Rollout-Ziel — und
hinterlässt damit einen Diff im Baum, den niemand bestellt hat und den kein
Sensor prüft. Der Implementer hat die Datei vor jedem Commit von Hand
zurückgenommen: die Bindung ist also wie bei den übrigen Trägern dieses Eintrags
**Disziplin, nicht Mechanik**. Sichtbar geworden ist der Nebeneffekt schon bei
`slice-080` (dort als Lauf-Artefakt im Baum aufgefallen und zurückgenommen); bei
`slice-082` trat er **planmäßig bei jedem Schema-Lauf** auf, weil der Nachzug den
Rollout in die Testläufe zieht.

Quelle: Implementer-Bericht `slice-082` · `tools/schema/apply-rollout.sh` ·
`Makefile` (`schema-rollout`, Pflicht-Report) · Sitzungsverlauf `slice-080`.
