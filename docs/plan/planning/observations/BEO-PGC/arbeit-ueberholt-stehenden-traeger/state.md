Zustand: **verkörpert** — Ausgang: **verkörpert** → `AGENTS.md` **§3.13** (neu:
*„Eine Arbeit, die eine beschriebene Eigenschaft bewegt, zieht ihre Träger nach“*)
· seit welle-20. Der **Lese-Schritt der `welle-20`-Closure** hat die Regel
**geschrieben**: kein Sensor (die Träger stehen nicht im Diff), aber ein
bestimmbarer **Leser** — die vier Fundstellen von `slice-094` fand der Implementer
auf einen `grep`-Auftrag hin, der aus diesem Eintrag stammte.

Zähler (Datei-Anzahl unter `evidence/`, real ausgezählt statt aus der
Ordinal-Erzählung unten übernommen): **21×** — die elf unten benannten
(`slice-091`, `slice-093`, `slice-094`, `slice-095`, `slice-096`,
`slice-097`, `slice-100`, `slice-101`, `slice-102`, `slice-103`,
`slice-105`) plus fünf zwischen `slice-105` und diesem Nachtrag ergänzte,
in dieser Ordinal-Erzählung bislang unbenannte Belege
(`evidence/slice-bench-schwellen-per-001-002-003.md`,
`evidence/slice-d-check-tracked-modul.md`,
`evidence/slice-nats-drittstream-core.md`,
`evidence/slice-rtm-letzte-zwoelf-tag-only.md`,
`evidence/slice-rtm-reste-sst-cfg-por.md` — jede für sich bereits ein
gültiger, in ihrer eigenen Datei begründeter Beleg dieser Klasse, hier nur
nachträglich in den Zähler aufgenommen, ohne die Ordinal-Erzählung für sie
rückwirkend zu schreiben) plus der achtzehnte Beleg,
`evidence/slice-sdk-python-http-client-flaeche.md`: Derselbe Vorgang wie
beim siebzehnten Beleg (Träger-Nachzug behauptet eine noch ausstehende
HTTP-Client-Fläche, real durch denselben Slice widerlegt), jetzt in der
Python-Sub-Area und mit einer geschärften Lehre — der Implementer-eigene
§3.13-Suchlauf war diesmal proaktiv **und** aus dem C#-Vorkommen abgeleitet,
deckte aber zwei unabhängige Lücken zugleich nicht ab: einen Datei-Glob, der
`pyproject.toml` ausschloss, und ein Wortmuster, das die deutsche
Formulierung „folgt erst" (statt „follow-up"/„added by") nicht traf.
Gefunden hat den dritten Treffer wieder der Reviewer (F-1, HIGH), nicht der
Implementer-Suchlauf. Details: `evidence/slice-sdk-python-http-client-flaeche.md`.
Der einundzwanzigste Beleg,
`evidence/slice-sdk-kotlin-nats-stream-client-flaeche.md`: anders als beim
zwanzigsten (C#-)Beleg fand hier **nicht** der Implementer-eigene
§3.13-Suchlauf die stehen gebliebene Stelle
(`sdks/kotlin/Dockerfile:77,104`, weiterhin `pgchangefeed-kotlin-0.1.0.jar`),
sondern der unabhängige Reviewer (F-2, MEDIUM) — dieselbe Unter-Klasse
„gefunden vom Reviewer, nicht vom Implementer-Suchlauf" wie bei `slice-095`/
`slice-097`/`slice-sdk-csharp-http-client-flaeche`. Zwei Sibling-Träger
(`harness/README.md`, `harness/mk/sdk.mk`), die denselben Sachverhalt
beschreiben, hatte der Implementer-Suchlauf bereits korrekt nachgezogen —
nur die beiden Dockerfile-Zeilen lagen außerhalb seines Suchraums. Fixrunde
real behoben, dritte/vierte unabhängige Bestätigung (Fixrunden-Reviewer,
Verifier). Ausgang bleibt **verkörpert**, kein neuer Regelschärfungs-Anlass.
Details: `evidence/slice-sdk-kotlin-nats-stream-client-flaeche.md`.
Der zwanzigste Beleg, `evidence/slice-sdk-csharp-nats-stream-client-flaeche.md`:
dieselbe Erfolgsform wie beim neunten/zehnten/elften Vorgang — der
Implementer-eigene §3.13-Suchlauf traf mehrere Fundstellen über **zehn**
Dateien (real ausgezählt per `git diff --stat`, korrigiert gegenüber der
zunächst genannten neun — Fixrunde review-slice-sdk-csharp-nats-stream-client-flaeche
F-3), die den
Lieferstand von `PgChangeFeed.Client` als „HTTP-API und gRPC-Stream" bzw.
über den `0.1.0`-Artefaktnamen beschrieben und durch die reale NATS-
Vollinhalts-Fläche (vierte Client-Fläche, Version-Hebung auf `0.2.0`)
falsch wurden, und behob sie alle im selben Zug — mit einer Erweiterung:
zwei dieser Fundstellen (`tools/harness/sdk-pack-csharp.sh`s Kommentar,
`Sse/Models/Change.cs`s „three surfaces"-Kommentar) lagen außerhalb des im
DoD-Wortlaut vorgeschriebenen `grep`-Musters und wurden nur durch
aufmerksames Lesen benachbarter Dateien gefunden — derselbe Grenzfall, den
`AGENTS.md` §3.13 §Grenze bereits für Zahlen/Prosa-Umformulierungen
benennt, hier erstmals für einen Kommentar-Text belegt. Anders als beim
Erstlauf trägt diese Zeile keine zusammenfassende „Fundstellen"-Zahl mehr
(`verifikation-slice-sdk-csharp-nats-stream-client-flaeche.md` §2: keine
im Text definierte, mechanisch eindeutige Zähleinheit für „Fundstelle" —
anders als für „Datei"); die mechanisch eindeutige Zahl bleibt **zehn
Dateien**. Details:
`evidence/slice-sdk-csharp-nats-stream-client-flaeche.md`.
Der neunzehnte Beleg, `evidence/slice-sdk-python-publish-workflow.md`: eine
**neue Form** innerhalb dieser Klasse — das überholende Ereignis liegt hier
zum ersten Mal vollständig **außerhalb** jedes Commits und jeder
Versionskontrolle (ein realer `sdk-csharp-v0.1.0`-Tag-Push samt `dotnet
nuget push` gegen NuGet.org, beides Betreiber-Handlungen nach der
`slice-sdk-csharp-publish-workflow`-Closure), das den Träger
`docs/user/releasing.md` überholte, ohne dass ein Sensor oder Diff dieses
Repos das je hätte zeigen können; gefunden hat es der Coordinator beim
Gegenlesen des Nachbarabschnitts, nicht der auf den eigenen
Slice-Gegenstand begrenzte §3.13-Suchlauf des Implementers — behoben in
einem eigenständigen Commit (`38a137f7`), dreifach unabhängig gegen
NuGet.org und `git tag -l` nachgeprüft (Reviewer, Verifier, Details in
`evidence/slice-sdk-python-publish-workflow.md`).
Der siebzehnte Beleg, `evidence/slice-sdk-csharp-http-client-flaeche.md`: Der
Vorgänger-Slice
(`slice-sdk-csharp-projektgeruest`) hatte `sdks/csharp/README.md` §Status
bewusst mit Verweis auf genau diesen Folge-Slice als HTTP-Fläche-Lieferant
geschrieben — bei Niederschrift wahr, durch die reale Auslieferung von
`PgChangeFeedHttpClient` (zehn Methoden) falsch geworden. Gefunden hat den
Fund nicht der Implementer-eigene §3.13-Suchlauf, sondern der **Reviewer**
(F-1, HIGH, merge-blockierend) — dieselbe Unter-Klasse „gefunden vom
Reviewer, nicht vom Implementer-Suchlauf" wie bei `slice-095`/`slice-097`.
In der Fixrunde behoben, Fixrunden-Nachprüfung und Verifikation bestätigen
unabhängig voneinander keinen weiteren stehen gebliebenen Satz (Details:
`evidence/slice-sdk-csharp-http-client-flaeche.md`). Der elfte Vorgang
(`slice-105`) trifft dieselbe Klasse über die
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
