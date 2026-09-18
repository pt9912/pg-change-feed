Zustand: verkörpert → `AGENTS.md` §3.14 ("Ein Aufrufer von
`make schema-rollout` gegen ein möglicherweise bereits migriertes Ziel
trägt seine eigene Idempotenz-Wache").

Begründung: Der Root-Cause-Fix (die verbleibenden Fremdobjekte — vier
SQL-Funktionen — ins neutrale Modell überführen) ist an
`POST_EXECUTE_DRIFT` (Exit 5) gebunden, ein real isoliert reproduziertes
Upstream-Verhalten von d-migrate
(`docs/plan/planning/observations/BEO-PGC/d-migrate-nacharbeit/`), und
damit derzeit kein gangbarer Weg — kein Ein-Zeilen-Fix, kein
terminierbarer Zeitpunkt. Ein dokumentiertes
`--allow-destructive`-Handling im `schema-rollout`-Makefile-Target wäre
der falsche Gegenzug: es würde echte destruktive Änderungen (Spalten-
oder Tabellen-Abbau, nicht nur die sechs bekannten Fremdobjekte)
ununterscheidbar mit durchlassen — ein zu grobes Werkzeug für ein eng
umrissenes Problem. Der real dreifach unabhängig gefundene lokale
Workaround — ein Existenz-Check gegen ein bekanntes Kernobjekt vor dem
Rollout-Aufruf, wie zuletzt in `examples/bootstrap.sh` — ist dagegen
billig (zwei bis drei Zeilen) und trägt die nötige Information bereits
am richtigen Ort (dem Aufrufer, der weiß, ob sein Lauf wiederholbar ist).
Ein viertes Auftreten ist wahrscheinlich, weil der Objektumfang mit
jedem neuen `nacharbeit-*.sql`-Skript wächst und nichts andeutet, dass
das aufhört; die einzige noch fehlende Maßnahme war, das Muster einmal
zu benennen, statt es beim nächsten Aufrufer erneut still zu erfinden —
das leistet jetzt die Hard Rule.

Zähler (abgeleitet): 3× (evidence/slice-016.md, evidence/slice-063.md,
evidence/slice-beispiele-compose-bootstrap.md) — Schwelle erreicht und
mit diesem Ausgang aufgelöst.
