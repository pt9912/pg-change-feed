# BEO-PGC/wartegrenze-ohne-zeitgrenze-im-startpfad

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft den Startpfad der
Composition Root, keine eigene Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Der Prozessstart wartet vor dem Stream-Lauf auf einen
synchronen Vorlauf, der die offenen Anträge der Antrags-Queue verarbeitet
(`runStreamAfterAdministrationPass` in `internal/bootstrap/wiring.go`). Der Vorlauf
trägt keine eigene Frist: ihn beendet allein der Kontext des Aufrufs, in `Run` der
Kontext des Streams, den der Prozess-`ctx` und die WAL-Fehlerschwelle beenden. Ein
Antrag, der lange läuft, ohne zu hängen (etwa `ALTER PUBLICATION … ADD TABLE`, die
auf eine Tabellen-Sperre wartet — hergeleitet aus dem Code, nicht am laufenden
Prozess gemessen), oder eine lange Queue hält die Erfassung der Quelle an. Die
Lebenszeichen-Goroutine startet vor der Sequenz, und der `--healthcheck` liest nur
das Alter dieses Lebenszeichens: ein wartender Vorlauf erscheint als gesund
(hergeleitet, nicht gemessen). Gemessen ist das Ende des Wartens durch einen
endenden Kontext (`TestRunStreamAfterAdministrationPassHoldsTheStreamUntilThePassEndsAndContextCancelEndsIt`);
eine Zeitgrenze des Vorlaufs und eine Anzeige des Wartens gibt es nicht.

**Warum das zählt:** Die Abhilfe nach einer nicht anwendbaren Regel
(`LH-FA-CFG-007.a`) stellt sich mit dem Vorlauf einen neuen Warte-Zustand vor der
Erfassung; weder das Pflichtenheft noch `ADR-0112` nennen eine Frist oder eine
Sichtbarkeit dafür (`git grep -n -i -E 'Vorlauf|Zeitgrenze|eigene Frist|Wartegrenze' -- spec`
trifft keine Zeile, gemessen in der Closure), und `LH-FA-ADM-003` verlangt einen
erkennbaren Fehlerzustand, kein erkennbares Warten.

**Abgrenzung.** `BEO-PGC/lesesperre-ohne-zeitgrenze` beschreibt das Warten des
Backfill-Runs an einer fremden Sperre — dort war das Warten ein akzeptiertes Negativ
(das Handbuch nennt die Gegenrichtung, ein `lock_timeout` ist eine Designänderung). Hier hält
das Warten die **Erfassung der ganzen Quelle** an, und der Zustand ist im Heartbeat
nicht von „gesund“ zu unterscheiden.

Deklaration: `slice-transformationen-start-reihenfolge` (Review F-3, MEDIUM;
Verifikation V-4, LOW), Risiko §6 „Der Vorlauf verzögert den Stream-Start“,
Ausgang *weiter offen*.
