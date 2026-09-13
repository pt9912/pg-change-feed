# Architect-Verdikt: slice-044 — `cdc_admin`-`DELETE`-Grant gegen `welle-13` §6

**Rolle:** Architect (Modul 8 §Konflikt-Pfad als Rollen-Sequenz)
**Anlass:** Implementer-Fund während `slice-044` (Plan-Nachzug §3 Punkt 6,
Commit `824e001`): weder `cdc_capture` noch `cdc_admin` trugen ein
`DELETE`-Grant auf `cdc.transaction`/`cdc.change`, ohne das
`RunRetentionUseCase` real nichts löschen kann. `welle-13` §6 schließt eine
„Erweiterung der Least-Privilege-Rollen" ausdrücklich als Out-of-Scope aus
und benennt für einen tatsächlichen Bedarf „ein eigener Architect-Zug, kein
stiller Fortschritt dieser Welle" — dieser Zug wurde nicht durchlaufen; der
Implementer hat die Erweiterung stattdessen selbst vorgenommen, transparent
im Plan-Nachzug begründet und als neue Beobachtung geführt.
**Rolleninhaber:** pt9912 (Claude Sonnet 5, dieser Lauf)
**Datum:** 2026-09-13
**Bezug:** [`LH-FA-RET-002`](../../spec/lastenheft.md)…[`004`](../../spec/lastenheft.md),
[`LH-QA-SEC-001`](../../spec/lastenheft.md)…[`003`](../../spec/lastenheft.md),
[`ADR-0014`](../plan/adr/0014-retention-domain-policy.md),
[`ADR-0047`](../plan/adr/0047-rollenspezifische-dsn-verdrahtung.md),
[`ADR-0048`](../plan/adr/0048-heartbeat-grant-korrektur-select-ergaenzung.md)
(Präzedenzfall), [`ADR-0053`](../plan/adr/0053-retention-loeschausfuehrung-cdc-admin-delete-grant.md)
(diese Entscheidung), `docs/plan/planning/welle-13.md` §6,
`docs/plan/planning/in-progress/slice-044-retention-hintergrundjob.md` §3
Plan-Nachzug Punkt 6,
`docs/plan/planning/observations/BEO-PGC/architect-verdikt-rollen-scope-luecke`,
Commit `824e001`

---

## Verdikt

**Verdikt 2 — Folge-ADR mit `Supersedes ADR-0047` (teilweise).** Die
Grant-Erweiterung selbst ist inhaltlich richtig und bleibt bestehen — kein
Rückbau. Aber sie ist keine bloße DDL-Lückenschließung wie die Heartbeat-
oder Schema-Version-Nacharbeiten in `nacharbeit-roles.sql` (dort war die
Rollen-Zuordnung bereits durch eine `Accepted`-ADR vorentschieden, es fehlte
nur der exakte Grant-Text). Sie ist eine **neue Zuordnungsentscheidung**,
die in `ADR-0047`s Rollen-Tabelle nicht vorkam — `ADR-0047` wurde einen Tag
vor `RunRetentionUseCase`s Entstehung geschrieben und konnte diesen Aufrufer
gar nicht kennen. Diese Lücke wird jetzt mit `ADR-0053`
(`Supersedes ADR-0047`, dieselbe begrenzte Form wie `ADR-0048`) geschlossen.
`ADR-0053` ist Teil dieses Verdikts.

Verdikt 1 (ADR gilt, Plan hat falsch behauptet) und Verdikt 3 (Lockerung
legitim, aber nur nachträglich zu dokumentieren, ohne neue ADR) wurden
geprüft und **verworfen** — Begründung unten.

## Warum Verdikt 1 nicht trägt

`welle-13` §6 behauptet nicht fälschlich, dass keine Rollen-Erweiterung
nötig sei — sie formuliert die Bedingung ausdrücklich als offen: „träfe das
während der Umsetzung doch zu, wäre das ein eigener Architect-Zug". Der
ursprüngliche Architect-Verdikt
(`docs/reviews/architect-verdict-retention-loeschausfuehrung.md`, Frage 1)
prüfte real Domain-/Port-/ADR-Konsistenz (`ADR-0009`, `0011`, `0012`,
`0014`, `0029`) — und diese Prüfung war korrekt und vollständig für das, was
sie geprüft hat. Sie prüfte nur nicht die PostgreSQL-Rollen-/Grant-Ebene,
weil `ADR-0047`s Rollen-Tabelle zu diesem Zeitpunkt keinen
Retention-Aufrufer führte, den man hätte prüfen können. Es liegt also keine
falsche Behauptung vor, die zurückzuweisen wäre — die Lücke ist real und der
Implementer-Fund korrekt. Verdikt 1 passt nicht: Es gibt keine bestehende
Entscheidung, die „gilt" und die der Plan missverstanden hätte — die
Entscheidung fehlte schlicht.

## Warum Verdikt 3 (Bestätigung ohne neue ADR) nicht trägt

Dieselbe Frage wurde für einen strukturell ähnlichen Fall bereits einmal
entschieden: `ADR-0048` korrigierte einen `ADR-0047`-Grant-Text (Heartbeat,
`SELECT`-Ergänzung) explizit **per Folge-ADR**, obwohl die Rollen-Zuordnung
selbst (`cdc_admin` für den Heartbeat-Schreibpfad) bereits in `ADR-0047`
stand und nur der exakte SQL-Text nachjustiert wurde — `ADR-0048`s eigener
Kontext-Abschnitt nennt das ausdrücklich „Verdikt 2 aus Modul 8
§Konflikt-Pfad: die ADR wird per Folge-ADR `supersedes`d, nicht der
Implementer-Fund als 'Lockerung' durchgewinkt." Der hier vorliegende Fall
ist **strenger** als der Heartbeat-Fall: Dort war die Rollen-*Zuordnung*
(„Heartbeat läuft über `cdc_admin`") bereits `Accepted` und nur der
Grant-Text unvollständig; hier fehlte die Rollen-*Zuordnung selbst*
vollständig, und die real geprüfte Alternative (eine eigene vierte Rolle,
§Frage unten) ist kein triviales Randdetail, sondern berührt `ADR-0047`s
eigenen Re-Evaluierungs-Trigger (a) — „eine vierte CDC-Rolle mit eigenem
Aufgabenschnitt". Ein bloßer Bestätigungs-Vermerk in diesem Report allein
würde diese Abwägung nicht mit derselben Verbindlichkeit tragen wie eine
ADR, die Implementer künftig als Constraint lesen (Modul 8 §Rollen-Regeln).
Wenn schon der leichtere Heartbeat-Fall eine ADR bekam, verdient dieser
Fall — reale Erweiterung der destruktiven Fähigkeit auf Kernstore-Daten,
nicht nur auf Steuerungsdaten — mindestens dieselbe Form.

## Die reale Abwägung: `cdc_admin` vs. eigene Rolle

Geprüft wurde, ob `DELETE` auf `cdc.transaction`/`cdc.change` kategorial zu
`cdc_admin`s bestehendem Verwaltungspfad passt oder eine neue Kategorie
wäre. Der volle Abwägungstext steht in `ADR-0053` §Verglichene Alternativen;
das Kernargument: `cdc_admin` ist bereits die Rolle für „alles außer dem
Erfassungs-`INSERT`-Pfad" (Heartbeat, Tabellen-Aktivierung,
Consumer-Registrierung — allesamt periodische oder administrative
Schreibzüge ohne Bezug zum Replication-Stream). Die Retention-Löschausführung
teilt exakt diesen Aufgabenschnitt (periodischer Hintergrundzug, physisch
ausführend für eine bereits im Domain Core getroffene Freigabe), nur die
Zieltabelle unterscheidet sich — und die Zieltabelle ist dabei die
Kernstore-Tabelle statt einer Steuerungstabelle, was den Blast-Radius eines
kompromittierten `cdc_admin`-Logins real vergrößert (unwiderruflicher
Verlust von Change-Historie statt korrigierbarer Steuerungsdaten-Korruption)
— das wird in `ADR-0053` als **Negativ**-Konsequenz benannt, nicht
verschwiegen. Eine eigene vierte Rolle würde diesen Blast-Radius
strukturell trennen, kostet aber ein viertes DSN/Login/Compose-Var/
Doku-Segment für einen einzelnen Hintergrundzug, während kein anderer
`cdc_admin`-Hintergrundzug eine eigene Rolle braucht. `ADR-0053` wählt
`cdc_admin` (Option A) und hält die vierte Rolle (Option B) als benannten,
künftig möglichen Schritt im Re-Evaluierungs-Trigger fest — keine
verschwiegene Alternative.

## Was das für Reviewer/Verifier bedeutet

Der Slice-Plan-Nachzug (`slice-044` §3 Punkt 6) selbst hat das Urteil, ob
seine eigene Grant-Erweiterung trägt, ausdrücklich offengelassen und an
Reviewer/Verifier delegiert. Diese Delegation ist mit diesem Verdikt
aufgelöst: Die Erweiterung **trägt inhaltlich** und bleibt bestehen; formal
fehlte die Architect-Entscheidung, die jetzt mit `ADR-0053` nachgetragen
ist. Reviewer/Verifier prüfen ab jetzt gegen `ADR-0053`, nicht mehr gegen
eine offene Frage.

## Folgeaktionen dieses Verdikts

1. [`ADR-0053`](../plan/adr/0053-retention-loeschausfuehrung-cdc-admin-delete-grant.md)
   geschrieben und in den ADR-Index aufgenommen (`Supersedes ADR-0047`,
   begrenzt auf die Rollen-Zuordnungstabelle).
2. `docs/plan/planning/welle-13.md` §6 erhält eine Klarstellung: die dort
   formulierte Bedingung ist real eingetreten und mit `ADR-0053` erfüllt —
   kein stiller Fortschritt, sondern ein nachgetragener, aber vollzogener
   Architect-Zug.
3. `BEO-PGC/architect-verdikt-rollen-scope-luecke` bleibt unverändert offen
   (1×, unter der 3×-Schwelle) — die dort benannte allgemeine Beobachtung
   (Architect-Verdikte prüfen nicht durchgängig die Rollen-/Grant-Konsequenz
   neuer physischer Schreib-/Löschpfade) ist mit diesem einen Vorfall nicht
   ausgeräumt und bleibt Gegenstand künftiger Sichtung.
4. Kein Rückbau von Commit `824e001` — die Grant-Erweiterung war inhaltlich
   richtig; nur die formale Architect-Bestätigung fehlte.

Weder `slice-044`s Code noch `nacharbeit-roles.sql` wurden im Rahmen dieses
Verdikts geändert.
