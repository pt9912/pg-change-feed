Zustand: **verkörpert** — Ausgang: **verkörpert** → `AGENTS.md` **§3.13** (neu:
*„Eine Arbeit, die eine beschriebene Eigenschaft bewegt, zieht ihre Träger nach“*)
· seit welle-20. Der **Lese-Schritt der `welle-20`-Closure** hat die Regel
**geschrieben**: kein Sensor (die Träger stehen nicht im Diff), aber ein
bestimmbarer **Leser** — die vier Fundstellen von `slice-094` fand der Implementer
auf einen `grep`-Auftrag hin, der aus diesem Eintrag stammte.

Zähler (abgeleitet): **4×** (evidence/slice-091.md, evidence/slice-093.md,
evidence/slice-094.md). Der dritte Vorgang trifft **vier** Stellen in **zwei**
Dateien — und die letzte ist eine **Korrektur**, die eine andere Stelle
derselben Klasse beheben sollte (zwei Herkünfte in einer Klammer).
`slice-093` ist der erste **angenommene** Fall: der Satz in
`harness/sensors/coverage-gate.md` war am Parent **wahr** und wurde durch die
Arbeit **falsch** (der neue netzlose Test fährt einen der zwei genannten Blöcke
deterministisch). Gefunden hat ihn der Implementer auf den `grep`-Auftrag hin,
der aus diesem Eintrag stammt — der Eintrag hat sich damit zum ersten Mal
**bezahlt**. Das Erstauftreten fiel im
Delta-Review zu `slice-091` auf: der Slice gab `driving/grpc/streamv1` eine
Testdatei und machte damit drei Sätze in `harness/sensors/coverage-gate.md`
falsch, die niemand im Diff hatte. Die **Reparatur** hat den Fall zunächst
verschärft (eine Zählung wurde tragend, die ihre Gruppe nicht erzeugt).

**Ein geprüfter Kandidat — und abgelehnt (nicht gezählt).** `slice-092` hat die
Frage aufgeworfen, ob eine **alternde Deixis** ohne Zahl-Drift ein Vorkommen
dieser Klasse ist: die Wendung „desselben, hier gegenständlichen
Produktionsstands" in `harness/sensors/coverage-gate.md` war nach einem
test-only-Slice irreführend, obwohl sich **keine** Zahl geändert hatte. Der
Delta-Review hat das **widerlegt**, aus drei Gründen: (1) hier **hilft**
`AGENTS.md` §3.12 — die Wendung hat eine Lauf-Größe als „der Ist-Stand"
vorgeführt, und das ist **bei Niederschrift** verboten, der Defekt war latent;
(2) „falsch **durch diese Arbeit**" trifft nicht zu — der Produktionsstand ist
gemessen derselbe geblieben; (3) die Form des Erstauftretens ist nicht
reproduziert (dort war der Satz bei Niederschrift **wahr**). Kein Beleg, kein
Zähler-Beitrag — die Grenze des Eintrags ist damit **geschärft**: was hierher
gehört, ist der Satz, den die Arbeit **umstößt**, nicht der, den sie nur
**sichtbar** macht.

**Nicht zu verwechseln** mit `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
(dort driftet ein Träger **schon**; hier driftet er **durch diese Arbeit**, und
die Arbeit war korrekt) und mit `BEO-PGC/beleg-befehl-traegt-seinen-satz-nicht`
(dort trägt ein **Beleg** seinen Satz nicht; hier trägt der **Satz** seinen Beleg
nicht mehr). §3.12 hilft **für den Fall dieses Eintrags** nicht: die Aussage trug ihren
Ursprung korrekt und war zum Zeitpunkt ihrer Niederschrift wahr — während sie
für den **abgelehnten Kandidaten** oben gerade **hilft** (dort war die Form
schon bei Niederschrift falsch gebunden).
