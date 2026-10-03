# BEO-PGC/umbau-aendert-still-beobachtbares-verhalten

**Sub-Area:** `*` (`PGC`, Greenfield) — Umbauten, deren Plan oder ADR ein unverändertes
Verhalten zusagt.

Die Beobachtung: Ein Umbau mit der Zusage „Verhalten unverändert“ ändert still ein
beobachtbares Verhalten, das weder ein Test noch ein Plan-Satz noch die ADR nennt; der
Vergleich gegen den Parent war der einzige Sensor. Der erste Vorgang:
`slice-meldungscodes-registry-fehlerkopf` stellte die Fehlerklassifikation von einer
`errors.Is`-Liste auf eine code-getriebene Klasse um. Zehn Sentinels, die die Liste am Parent
nicht kannte und auf `internal` fallen ließ, bekamen dadurch eine andere Klasse; `error_class`
im Heartbeat und das Metrik-Label dieser Fälle hätten sich geändert, obwohl
[`ADR-0144`](../../../../adr/0144-meldungscodes-nutzerseitige-kennungen.md) Festlegung 2 „Klassen
unverändert“ sagt. Der Reviewer fand es durch das Lesen von `git show <Parent>:…` gegen den Kopf.

**Abgrenzung:** keine Bindungslücke eines Negativtests
(`negativtest-ohne-bindung-an-seine-eingabe`: dort hängt ein Test nicht an seiner Eingabe) und
keine Aussage, die breiter ist als ihre Messung (`adr-aussage-breiter-als-ihre-messung`: dort
trägt die Aussage keinen Beleg); hier ist die Zusage wahr gemeint und der Umbau verletzt sie
durch eine Nebenwirkung auf einem Pfad, den niemand als betroffen benannt hat.

**Warum das zählt:** Die Zusage der Gleichheit ist nur so stark wie der Vergleich, der sie
prüft; ohne benannten Vergleichsbefehl über alle betroffenen Pfade bleibt die Prüfung dem Zufall
des Lesers überlassen.

Deklaration: `slice-meldungscodes-registry-fehlerkopf`, Review F-1 (MEDIUM).
