Zustand: **verkörpert** — Ausgang: **verkörpert** → `AGENTS.md` **§3.13** (neu:
*„Eine Arbeit, die eine beschriebene Eigenschaft bewegt, zieht ihre Träger nach“*)
· seit welle-20. Der **Lese-Schritt der `welle-20`-Closure** hat die Regel
**geschrieben**: kein Sensor (die Träger stehen nicht im Diff), aber ein
bestimmbarer **Leser** — die vier Fundstellen von `slice-094` fand der Implementer
auf einen `grep`-Auftrag hin, der aus diesem Eintrag stammte.

Zähler (abgeleitet): **11×** (evidence/slice-091.md, evidence/slice-093.md,
evidence/slice-094.md, evidence/slice-096.md, evidence/slice-097.md,
evidence/slice-095.md, evidence/slice-100.md, evidence/slice-101.md,
evidence/slice-102.md, evidence/slice-103.md, evidence/slice-105.md). Der
elfte Vorgang (`slice-105`) trifft dieselbe Klasse über die
Baseline-Versionierung statt über ein Runtime-Image: das Entfernen von
`.harness/baseline/v6.5.0/` überholte drei Links in `docs/reviews/**` — <!-- d-check:ignore (historische Vor-Migrations-Erwähnung, kein lebender Pin) -->
gefunden vom Implementer über denselben §3.13-Suchlauf, im selben Commit
per Zitat-Korrektur (`ADR-0073`) behoben, dieselbe Erfolgsform wie bei den
vorigen Vorgängen. Der zehnte Vorgang
(`slice-103`, letzter der sechs C#/Kotlin-Matrix-Slices) trifft dieselbe
Klasse ein zweites Mal in derselben Kotlin-Fortsetzung: (a) `harness/
README.md` „drei"→„vier" Images (`examples-kotlin`), direkt behoben in
`111f1cb`, dieselbe Form wie beim neunten Vorgang; (b) `slice-097`s §1, das
beim neunten Vorgang bereits als „für C# überholt, für Kotlin weiterhin
richtig" gemeldet wurde, ist jetzt — mit der Kotlin-Zelle geliefert — für
**beide** Sprachen überholt. Beide Funde entstammen demselben Vorgang und
zählen als **eine** Evidenzdatei (`evidence/slice-103.md`), kein neuer
Schwellen-Übertritt. Der neunte Vorgang (`slice-102`) trifft
`harness/README.md` §Sensors ein viertes Mal (Zahlenwort „drei"→„vier" bei
der `make examples-csharp`-Zeile, ausgelöst durch das vierte Runtime-Image)
— wieder gefunden vom Implementer über den vorgeschriebenen §3.13-Suchlauf
selbst und im selben Commit behoben, dieselbe Erfolgsform wie bei
`slice-093`/`slice-094`/`slice-100`/`slice-101`. Derselbe Vorgang trägt einen
**zweiten** Fund derselben Klasse, aber mit anderer Behandlungsform: `slice-
097`s §1 (`done/`, immutabel) wird durch `slice-102` teilweise überholt
(„ihre Bau-Kontexte erreichen sie heute nicht" gilt seither für C# nicht
mehr, für Kotlin weiterhin) — kein Zitat-Korrektur-Fall nach `ADR-0073`,
keine Editier-Gelegenheit, sondern der Regelfall von `AGENTS.md` §3.13:
gemeldet statt geändert. Beide Funde entstammen demselben Vorgang und
zählen als **eine** Evidenzdatei (`evidence/slice-102.md`), kein neuer
Schwellen-Übertritt. Der achte Vorgang (`slice-101`) trifft
`harness/README.md` §Sensors ein
drittes Mal (Zahlenwort „zwei"→„drei" bei den `make examples-csharp`/
`make examples-kotlin`-Zeilen, ausgelöst durch das dritte Runtime-Image je
Sprache) — wieder gefunden vom Implementer über den vorgeschriebenen
§3.13-Suchlauf selbst und im selben Commit behoben, dieselbe Erfolgsform
wie bei `slice-093`/`slice-094`/`slice-100`. Derselbe Slice trug einen
zweiten, davon unterschiedenen Fund (eine im Plan-Kopf/§6/§8 selbst
zitierte, bereits bei Niederschrift veraltete Zahl) — der zählt **nicht**
hier, sondern bei `BEO-PGC/zahl-in-traeger-driftet-gegen-die-messung`
(siehe dortiger Eintrag, Abgrenzung). Der siebte Vorgang
(`slice-095`, 2. Durchlauf) ist
die erste **selbstreferentielle** Variante: die Plan-Datei, die die
überholende Arbeit selbst trägt (§3-Begründung „keine Kante", durch die
eigene, spätere LP3-Lieferung desselben Slice falsch geworden), statt eines
externen Trägers — gefunden vom Reviewer (F-1), nicht vom §3.13-Suchlauf des
Implementers, der auf `grep`-Treffer über interne Code-Pfade beschränkt war
und diese Plan-Prosa-Zeile deshalb nicht erfasste. Der dritte Vorgang
trifft **vier** Stellen in **zwei** Dateien — und die letzte
ist eine **Korrektur**, die eine andere Stelle derselben Klasse beheben
sollte (zwei Herkünfte in einer Klammer). Der fünfte Vorgang
(`slice-097`) traf einen Träger **außerhalb** der vom Implementer-Suchlauf
durchsuchten Liste (`next/slice-095`, ein Planning-Dokument, kein
Code-/Harness-Träger) — gefunden erst durch einen zweiten, unabhängigen
Suchlauf (Reviewer).
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
