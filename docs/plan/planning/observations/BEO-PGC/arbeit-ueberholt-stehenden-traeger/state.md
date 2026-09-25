Deckel bei 32× (seit welle-backfill-bestand): weitere Auftreten, die vor dem Merge
vom Reviewer oder Verifier gefunden werden, Schwere ≤ LOW haben und einen bekannten
Träger-Typ treffen, bekommen keine `evidence/`-Datei, sondern stehen mit Finding-Kennung
in der Closure-Notiz des Slice (`../../README.md`, Deckel für verkörperte Einträge
ab 10×). Ausgang unverändert **verkörpert**, geschärft: `AGENTS.md` §3.13 §Suchform
(ganzer Baum, Symbolname · Zählwort · Beschreibung samt Hedge, Befehl im Codeblock);
der Suchlauf-Anteil bekommt ein Nachmess-Werkzeug in
[`slice-harness-suchlauf-nachmessen`](../../../open/slice-harness-suchlauf-nachmessen.md)
(Architect-Verdikt `architect-verdict-welle-backfill-bestand-lese-schritt`
§3.3 und §3.5).

Zustand: **verkörpert** — Ausgang: **verkörpert** → `AGENTS.md` **§3.13** (neu:
*„Eine Arbeit, die eine beschriebene Eigenschaft bewegt, zieht ihre Träger nach“*)
· seit welle-20. Der **Lese-Schritt der `welle-20`-Closure** hat die Regel
**geschrieben**: kein Sensor (die Träger stehen nicht im Diff), aber ein
bestimmbarer **Leser** — die vier Fundstellen von `slice-094` fand der Implementer
auf einen `grep`-Auftrag hin, der aus diesem Eintrag stammte.

Zähler (Datei-Anzahl unter `evidence/`, real ausgezählt statt aus der
Ordinal-Erzählung unten übernommen): **32×** — der Beleg
`evidence/slice-backfill-speicher-untersuchung.md` (der Höchstwert je Change bewegt, drei
Träger außerhalb des Diffs und der Suchwurzeln — `ADR-0124`, Architect-Verdikt, Plan des
Folge-Slice — tragen die frühere Zahl; Verifikation V-1 und V-2), die elf unten benannten
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
Der neunundzwanzigste Beleg,
`evidence/slice-backfill-snapshot-reader.md`: die bewegte Eigenschaft (Inhalt der
Phase `tier` — zwei Läufe statt eines) trug Beschreibungen in
`tools/harness/run-replication-tests.sh` und `.github/workflows/e2e.yml`, die
der Suchlauf nicht las, weil er Variablen- und Testnamen suchte, nicht die
Bezeichnung der Phase (Nachprüfungs-Review F-3, LOW); dazu zwei offene Pläne mit
drei statt vier ausgenommenen Paketen (F-16, LOW). Beide in den Fixrunden
gezogen, der Suchlauf trägt die gesuchte Eigenschaft seither mit beiden
Ständen — Anwendungs-Schärfung der verkörperten Regel `AGENTS.md` §3.13, kein
Schwellen-Übertritt.
Der einunddreißigste Beleg, `evidence/slice-backfill-bench-richtgroesse.md`: die
bewegte Eigenschaft (was `tools/bench-backfill.sh` misst — die Fixrunde ergänzte
die WAL-Messung) trug Beschreibungen im Vertrag `harness/targets/bench-backfill.md`
und in der `make bench`-Zeile von `harness/README.md`, die sie nicht nannten
(Review F-3, LOW); gefunden vom Reviewer, die Fixrunde zog beide Träger nach —
Anwendungs-Schärfung der verkörperten Regel `AGENTS.md` §3.13, kein
Schwellen-Übertritt.
Der dreißigste Beleg, `evidence/slice-backfill-e2e.md`: die bewegte Eigenschaft
(Ablauf des Snapshot-Imports: Sperre und Umschreib-Prüfung) trug eine
Sequenzdarstellung in `spec/architecture.md`, die das committete Suchlauf-Feld
nicht fand — sein Suchmuster nannte Sperre und Umschreiben, die Sicht führt den
Import als „Transaktion mit importiertem Snapshot“ (Review F-2, LOW). Gefunden
vom Reviewer, die Fixrunde zog die Sicht nach und das Feld trägt die Datei seither
mit beiden Ständen — Anwendungs-Schärfung der verkörperten Regel `AGENTS.md` §3.13
(das Suchmuster folgt der Bezeichnung, die der Träger führt, nicht der des
eigenen Diffs), kein Schwellen-Übertritt; der Lese-Schritt der Closure von
`welle-backfill-bestand` liest den Eintrag mit (32×).
Der achtundzwanzigste Beleg,
`evidence/slice-backfill-change-origin.md`: die bewegte Eigenschaft (Feldzahl
der `GET /changes`-Antwort) trug Träger in `*.kt` und `*.cs`
(SDK-Kommentare „twelve"/„eleven fields"), die das committete Suchlauf-Feld nicht
las — es suchte `docs/user`, `spec` und `*.go`/`*.proto`/`*.py`; das Kotlin-Wort
„twelve" stand am Zeilenende, „fields" in der Folgezeile. Gefunden vom Reviewer
(F-4, LOW), Suchraum in der Fixrunde auf den ganzen Baum erweitert, die Träger
als Meldung an den Folge-Slice `slice-backfill-sdk-origin` gegeben —
Anwendungs-Schärfung der verkörperten Regel `AGENTS.md` §3.13 (der Suchraum
folgt den Dateitypen, die die bewegte Eigenschaft tragen können, nicht denen des
eigenen Diffs), kein Schwellen-Übertritt.
Der siebenundzwanzigste Beleg,
`evidence/slice-backfill-spec-nachzug.md`: die Träger waren offene
Folge-Pläne (Plan-Sätze „leer, solange keine Auswertung sie setzt", „vier
Arten"), die der Spec-Nachzug durch die Festlegung der zwei Warn-Spalten und
der fünf Antragsarten überholte; das committete Suchlauf-Feld las `spec`,
`docs/user`, `harness` und die READMEs, nicht `docs/plan/planning` (Review F-2,
LOW), und das Muster der Fixrunde suchte die „leer"-Formen statt des Hedges
„Warn-Spalte(n)" (Verifikation V-1, INFO) — Anwendungs-Schärfung der
verkörperten Regel `AGENTS.md` §3.13, kein Schwellen-Übertritt: ein Slice, der
eine offene Zahl festlegt, sucht auch den Hedge des offenen Punkts über alle
Plan-Träger. Details: `evidence/slice-backfill-spec-nachzug.md`.
Der sechsundzwanzigste Beleg,
`evidence/slice-sdk-python-http-reale2e.md`: das §3.13-Suchlauf-Feld
nannte die eigene `harness/README.md`-Python-Zeile und verfehlte die zwei
Nachbar-Zeilen derselben Werkzeug-Tabelle, deren Endklassen die
Erweiterungs-Zusage („erweitert sich um den Python-HTTP-Abschnitt") als
ausstehend weitertrugen, obwohl beide Erweiterungen mit diesem Zug real
existierten — das vierte Auftreten der in `AGENTS.md` §3.13
dokumentierten Fund-Struktur (nach `slice-091`/`093`/`094`), diesmal mit
der Schreib-Verbreiterung als Wurzel: der Suchlauf folgt der
Schreib-Verbreiterung der Bewegung über alle Träger der Eigenschaft,
nicht dem deklarierten Eintäger-Mustersatz. Gefunden vom Reviewer (F-2,
MEDIUM), gezogen in der Fixrunde (`d668b9cc`), vom Verifier gegen beide
Stände bestätigt (Verifikation §3, Zeile 8). Ausgang bleibt **verkörpert**,
kein neuer Schwellen-Übertritt — die Schärfung ist eine Anwendungs-
Schärfung der verkörperten Regel `AGENTS.md` §3.13. Details:
`evidence/slice-sdk-python-http-reale2e.md`.
Der fünfundzwanzigste Beleg,
`evidence/slice-sdk-kotlin-reale2e.md`: die Prüf-Angaben des committeten
§3.13-Suchlauf-Felds trugen zwei Baum-Widerlegungen — F-1 (MEDIUM,
Reviewer): die Geprüft-Zeile nannte `sdks/kotlin/README.md` §Status, eine
Adresse ohne Artefakt (die reale Träger-Datei ist
`sdks/kotlin/pgchangefeed-kotlin/README.md`); F-6 (INFO): die
Behandlungs-Spalte zur C#-Zeile behauptete „gezogen", während die Zeile an
beiden Ständen byte-identisch bleibt — „gemeldet, nicht gezogen" ist die
wahre Behandlung, die Endklause bleibt als Verlaufs-Aussage wahr. Beide in
der Fixrunde (`c6523009`) gezogen, beide vom Verifier gegen beide Stände
bestätigt (Verifikation §3, alle sechs Feld-Zeilen); die bewegte Eigenschaft
selbst traf keinen falsch werdenden Träger — dieselbe Fundstruktur wie beim
dreiundzwanzigsten Beleg (der Suchlauf-Befund selbst wurde zum Träger,
dessen Prüf-Angaben der Review am Baum nachmisst), hier zusätzlich mit der
Adresse-Widerlegung. Ausgang bleibt **verkörpert**, kein neuer
Schwellen-Übertritt — die Schärfung ist eine Anwendungs-Schärfung der
verkörperten Regel `AGENTS.md` §3.13. Details:
`evidence/slice-sdk-kotlin-reale2e.md`.
Der vierundzwanzigste Beleg,
`evidence/slice-sdk-csharp-reale2e.md`: der §3.13-Suchlauf des
Implementers trug vier geplante Träger-Zeilen (alle bestätigt) und
verfehlte drei weitere — alle drei hängen an derselben Wurzel: der
`.d-check.yml`-`trace.coverage`-Eintrag des Diff verbreitert die Menge
der kuratierten Coverage-Dimensionen, und jede Prosa-Zeile, die diese
Menge aufzählt (`harness/README.md`s `make doc-trace`-Zeile,
`harness/sensors/docs-check.md` §Grenze, die Welle-Plan-§6), wird durch
exakt diese Verbreiterung überholt — der deklarierte Mustersatz
(geplante Träger, `Datei:Zeile`-Form) beschrieb den Suchraum, nicht die
Bewegung. Gefunden vom Reviewer (F-1, MEDIUM), gezogen in der Fixrunde
(`86892bdb`). Geschärfte Lehre (Closure-Notiz): der Suchraum eines
§3.13-Suchlaufs folgt der Schreib-Verbreiterung der Bewegung, nicht dem
deklarierten Mustersatz. Ausgang bleibt **verkörpert**, kein neuer
Schwellen-Übertritt — die Schärfung ist eine Anwendungs-Schärfung der
verkörperten Regel `AGENTS.md` §3.13. Details:
`evidence/slice-sdk-csharp-reale2e.md`.
Der dreiundzwanzigste Beleg,
`evidence/slice-sdk-python-nats-stream-client-flaeche.md`: die schmale
Plan-Mustersatz-Grep-Form des §3.13-Suchlaufs verfehlte alle fünf
Träger-Stellen in deren tatsächlicher Schreibform und ließ `harness/` und
die Wurzel-READMEs außerhalb des Pfadraums (F-4, gefunden vom Reviewer);
das committete Suchlauf-Feld behauptete danach zwei vollzogene Nachzüge,
die der Baum widerlegte (R-1/R-2) — der Suchlauf-Befund selbst wurde zum
Träger, dessen Behandlungs-Angaben das Re-Review am Baum nachmisst.
Geschärfte Lehre (Closure-Notiz): der Suchlauf folgt §3.13s Wortlaut
(über die Träger nach der bewegten Eigenschaft), nicht dem Plan-Mustersatz.
Ausgang bleibt **verkörpert**, kein neuer Schwellen-Übertritt. Details:
`evidence/slice-sdk-python-nats-stream-client-flaeche.md`.
Der zweiundzwanzigste Beleg,
`evidence/slice-sdk-python-sse-client-flaeche.md`: die Lücken-Struktur
wiederholt sich innerhalb eines Slice — der §3.13-Suchlauf war je Slice
und je **einer** bewegten Eigenschaft gebunden („SSE bleibt außerhalb
des Packages“), während der Fix-Zug desselben Slice eine **zweite**
Eigenschaft bewegte (Testdatei-Übergabe des Integration-Images:
docker run-Argument → Umgebungsvariable `PGCHANGEFEED_TEST_FILE`) und
ohne eigenen Suchlauf durchlief. Der Haupt-Review traf den Fall auf der
ersten Eigenschaft (`harness/README.md`s
`make test-sdk-python-integration`-Zeile trug die Ein-Flächen-Form,
F-2 MEDIUM — gefunden vom Reviewer, die Plan-Nachzug-Liste nannte die
Zeile nicht), der Re-Review auf der zweiten (FR-1 MEDIUM — drei
Phrase-Stellen: Runner-Skriptkopf „Stufe-ENTRYPOINT“/„docker
run-Argument“ und die neue `harness/README.md`-Zeile); Fixrunde 2 zog
die drei Stellen auf die ENV-Form und der Suchlauf wurde um die zweite
Eigenschaft erweitert — der Verifier trug die Stand-Deklaration der
Erweiterung nach (V-1: der wahre Vorher-Stand des
`harness/README.md`-Trägers ist `3c941b0b`, nicht `6bbe99d9`).
Geschärfte Lehre (Closure-Notiz): der Suchlauf wird je bewegter
Eigenschaft geführt, nicht je Slice. Ausgang bleibt **verkörpert**, kein
neuer Schwellen-Übertritt — die Schärfung ist eine Anwendungs-Schärfung
der verkörperten Regel `AGENTS.md` §3.13. Details:
`evidence/slice-sdk-python-sse-client-flaeche.md`.
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
