# BEO-PGC/geschaetzter-wert-als-grenze

**Sub-Area:** `*`/`PGC` (Repo-weite Default-Sub-Area — betrifft die Form, in der
ein **geschätzter** Wert durch Weitergabe zu einer **Grenze** wird, keine eigene
Sub-Area im Sinn der Modus-Deklaration).

Die Beobachtung: Ein Wert, der als **Schätzung** notiert ist — mit Tilde, mit
„≈", als Spanne —, wird auf dem Weg durch **mehrere Träger** zu einer
**Tatsache**, und am Ende zu einer **Grenze**: einer Aussage darüber, was
**nicht** geht. Keine der Stationen hat ihn gemessen; jede hielt ihn für ein
Zitat aus der vorigen.

**Warum das schwer zu sehen ist:** Die Tilde fällt **stückweise** weg. Das
Nachbardokument schreibt „≈45"; der Plan schreibt „45 (nach der ADR)"; die
Umsetzung schreibt „nicht erreichbar". Jeder Schritt ist für sich plausibel, und
**keiner** ist eine Lüge — am Anfang steht eine ehrliche Schätzung, am Ende eine
falsche Grenze. Der Wert hat sich nie geändert; nur seine **Verbindlichkeit**.

**Der Unterschied zu den benachbarten Klassen.** `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`
beschreibt einen **genannten Beleg**, der seinen Satz nicht trägt; hier trägt
jeder Träger seinen Satz — er zitiert nur **nicht** als Zitat.
`BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung` beschreibt einen Wert, der
**gegen eine Messung** driftet; hier gibt es **keine Messung**, an der er
driften könnte — das ist der Kern. Und `BEO-PGC/arbeit-ueberholt-stehenden-traeger`
beschreibt einen Träger, den **die Arbeit** falsch macht; hier macht ihn der
**Weiterweg** falsch.

Belegt an einem abgeschlossenen Vorgang:

- **`slice-094`** (`review-slice-094` F-1, HIGH): [ADR-0082](../../../../adr/0082-coverage-schnittmass-composition-root-nicht-netzlos.md) §Kontext (4a) notiert
  **„≈45 netzlos erreichbar"** für `cmd/pg-change-feed` — eine Schätzung. Der
  **Planner** hat sie im Slice-Plan §1 als Grenze gelesen („die restlichen sind
  nicht erreichbar") und im §2 als LP3-Grund geführt; der **Implementer** hat
  sie übernommen und in `harness/sensors/coverage-gate.md` als **Grenze**
  geschrieben („liegen außerhalb des netzlosen Tiers"). Der **Reviewer** hat sie
  gemessen: ein Test, der dieselben vier Sondermodi mit **vollständigem** ENV
  fährt, deckt sie **alle** — `cmd` **49 von 49** statt 45, Gesamt **83,08 %**
  statt 82,87 %. Die Schätzung war um **vier Statements** zu niedrig, und die
  Differenz war als **Grenze** formuliert.

**Die Antwort ist eine Handlung an der letzten Station, kein Werkzeug.** Ein
Sensor müsste Schätzungen von Zitaten unterscheiden können — das kann er nicht.
Was hilft, ist die **Messung dort, wo der Wert zur Grenze wird**: wer „nicht
erreichbar" schreibt, hat eine Probe zu fahren. Der Reviewer hat genau das
getan, und es war eine Zeile Arbeit.
