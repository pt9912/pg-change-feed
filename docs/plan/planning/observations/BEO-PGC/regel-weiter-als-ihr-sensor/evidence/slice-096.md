# Beleg: slice-096

Vorgang: `slice-096` — Konfigurationsdatei-Nachzug.

Fund: **`AGENTS.md` §3.13 ist weiter als das, was sie durchsetzt — an zwei
Stellen, und beide sind an ihrem ersten Fall gemessen.**

Die Regel (verkörpert seit der `welle-20`-Closure, `slice-089`/`-091`) verlangt
von einer Arbeit, die eine beschriebene Eigenschaft bewegt: den `grep`-Lauf nach
der **bewegten Eigenschaft** (nicht über den eigenen Diff), **beide Stände**
gemessen, und **das Gefundene und das Nichtgefundene** im Bericht.

**Grenze 1 — der Suchlauf greift Symbolnamen, nicht Zahlen.** Der Implementer hat
den Lauf korrekt geführt und den **Symbolnamen** gefunden. Die **zwei
Zeilen-Lokatoren** in derselben ADR fand er nicht, und das ist keine Nachlässigkeit:
der Suchlauf sucht **Namen**; Lokatoren sind **Zahlen**. Gefunden hat sie der
Review, den vierten Anker der Architect. Die Regel nennt ihre Suchform nicht —
und ihre Suchform deckt eine Teilmenge ihrer Klasse.

**Grenze 2 — das Ergebnis hat keinen Träger im Repo.** §3.13 sagt, das Ergebnis
stehe „im Bericht". Der Bericht eines Implementer-Laufs ist ein **Handoff** —
er wird nicht committet und ist über Läufe hinweg nicht nachlesbar. Damit ist der
Lauf selbst nicht belegbar: wer später prüfen will, ob der `grep` gelaufen ist,
findet nichts. Die Regel verlangt eine Handlung, deren **Spur** sie nicht
verlangt.

**Warum das zählt:** Eine Regel, deren Reichweite über ihre Durchsetzung
hinausgeht, ist die Klasse, die dieser Eintrag selbst beschreibt — hier trifft sie
ihren **eigenen** Träger. Die verfügbare Falsifikation ist die **Messung**: dieser
Beleg ist der erste Fall, in dem jemand den Lauf gegen seine Regel gehalten hat,
statt ihn zu glauben.

**Ausgang:** nicht in diesem Vorgang zu setzen — eine Hard Rule zu ändern ist ein
Architect-Zug (Modul 8). Der Eintrag steht damit auf der **Schwelle**; den
Ausgang weist der **Lese-Schritt der nächsten Wellen-Closure** zu (Modul 6).

Quelle: Verifikationsbericht zu `slice-096` (V-5) ·
`AGENTS.md` §3.13 · Review zu `slice-096` (F-2) ·
`BEO-PGC/arbeit-ueberholt-stehenden-traeger/evidence/slice-096.md`.
