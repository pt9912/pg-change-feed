# Architect-Verdikt: Chronik-Sprache in Code-Kommentaren — 4. Auftreten (Reviewer als tragende Instanz bestätigt)

**Rolle:** Architect (Modul 8)
**Anlass:** `docs/plan/planning/observations/BEO-PGC/slice-chronik-in-code-kommentar`
— nach dem 3×-Verdikt (der Architect-Verdikt zur Slice-Chronik in Code-Kommentaren,
verkörpert als geschärfte Selbstprüf-Instruktion in
`.claude/commands/implement-slice.md` Schritt 20) tritt die Klasse im
*ersten* Slice nach der Schärfung erneut auf: `slice-052`
(Review zu `slice-052` F-1 HIGH,
`internal/adapters/driven/natsnotify/notify.go`). Die eigene
Eskalationsnotiz in `state.md` hatte dieses Szenario vorab benannt: „Tritt
die Klasse trotz der geschärften Instruktion ein viertes Mal auf, ist das
ein Signal, dass Enumeration allein nicht trägt — neue Beobachtung oder
Zähler-Fortschreibung, Urteil beim nächsten Lese-Schritt." Direkter
Architect-Zug aus einer Nutzerfrage, **vor** der formalen `slice-052`-Closure
(`slice-052` liegt zum Zeitpunkt dieses Verdikts noch in `in-progress/` —
der reguläre Lese-Schritt/Zähler-Eintrag durch den Planner steht noch
aus, siehe §Was dieses Verdikt NICHT tut).
**Rolleninhaber:** pt9912 (Claude Sonnet 5, dieser Lauf)
**Datum:** 2026-09-13
**Bezug:** [AGENTS.md](../../AGENTS.md) §3.7, Modul 8 §Kernidee („wer
geschrieben hat, reviewt nicht") und §Konflikt-Pfad als Rollen-Sequenz,
Modul 6 §Das Beobachtungs-Register, vorheriger
Architect-Verdikt zur Slice-Chronik in Code-Kommentaren,
das Review zu `slice-052` (F-1, Erstbefund),
der Review-Report zur Fixrunde von `slice-052` (Bestätigung + 4.-Beleg-Notiz),
`.harness/skills/reviewer.md`, `.claude/commands/implement-slice.md`
Schritt 20, `docs/plan/planning/observations/BEO-PGC/slice-chronik-in-code-kommentar/`.

---

## Frage

Trägt der 3×-Verdikt (geschärfte, diff-skopierte Selbstprüf-Instruktion,
kein mechanischer Sensor) noch, nachdem die Klasse trotzdem ein viertes
Mal auftrat — oder braucht es eine stärkere Mechanisierung (neuer Sensor,
Gate, oder eine grundsätzlich andere Zuordnung von Verantwortung
zwischen Implementer-Selbstprüfung und Reviewer)?

## Empirischer Befund: 4/4 vor Merge gefangen, 0/4 in `main`

Der entscheidende Unterschied zu den ersten drei Belegen: **keiner der
vier Fälle hat je gemergten Code erreicht.** Das Review zu `slice-041`,
der Review-Report zur Fixrunde von `slice-041` und das Review zu `slice-044` wurden allesamt
vor Merge korrigiert; Review zu `slice-052` F-1 ebenso — Fixrunde
`14e790e`, bestätigt im Review-Report zur Fixrunde von `slice-052`, danach
Verifikation (`461973e`). Die Hard Rule `AGENTS.md` §3.7 wurde in `main`
zu keinem Zeitpunkt verletzt. Das Zwei-Schichten-Design — Implementer-
Selbstprüfung als erste, Reviewer als unabhängige zweite Instanz — hat in
100 % der gezählten Fälle an der Stelle gehalten, die für die Zusage
zählt: dem Merge-Gate.

## Diagnose: kein weiterer Enumerations-Befund, sondern eine Struktur-Bestätigung

Der 3×-Verdikt hatte den dritten Fall als *Enumerations-Lücke* diagnostiziert
(eine von mehreren korrigierten Stellen im selben Commit wurde übersehen)
und die Schärfung entsprechend zielgerichtet: ein diff-skopierter
Kandidatenlauf, der **jede** geänderte Datei erfasst, nicht nur die vom
Implementer erinnerten.

Der vierte Fall ist **kein** Wiederauftreten derselben Diagnose. Das
Pattern aus Schritt 20 (`slice-[0-9]+|welle-[0-9]+|…`) hätte
„Folge-Slice `slice-053`" im geänderten `notify.go` zweifelsfrei
getroffen — es gab in diesem Commit auch keine zweite, vom Kandidatenlauf
unentdeckte Chronik-Stelle, die die Enumerationslogik hätte überfordern
können (F-2 aus dem Review zu `slice-052` ist eine andere Klasse: ein
ephemerer Artefakt-Verweis, keine Slice-Chronik). Der einzig plausible
Befund ist: **der Kandidatenlauf wurde in diesem Implementer-Durchlauf
nicht (oder nicht wirksam) ausgeführt** — eine Compliance-Lücke in der
Ausführung des Schritts, nicht eine Lücke in seiner Konstruktion. Eine
weitere Schärfung der Enumerationslogik (Option 3 der Anfrage — z. B. ein
zusätzliches Musterbeispiel für die Klammerform „Folge-Slice
`slice-NNN`") liefe gegen ein Muster, das das bestehende Pattern bereits
abdeckt, und würde diesen Fehlermodus nicht beheben.

Das ist strukturell erwartbar, nicht überraschend: Schritt 20 läuft **im
selben Kontext**, der den Kommentar geschrieben hat — exakt die
Konstellation, vor der Modul 8 §Kernidee warnt („wer geschrieben hat,
reviewt nicht … dieselbe Sicht denselben Fehler übersieht"). Ein
Selbstprüf-Schritt im eigenen Kontext hat eine strukturell begrenzte
Trefferquote, unabhängig davon, wie präzise seine Instruktion ist — der
tragende Schutz gegen genau diesen blinden Fleck ist im Harness-Design
bereits der unabhängige Reviewer, nicht die Selbstprüfung. Das
Registerdesign sah diese Rollenteilung immer schon vor (Modul 8, sechs
Rollen mit getrenntem Eingabe-Kontext); der 3×-Verdikt hatte lediglich
noch nicht ausdrücklich benannt, dass Schritt 20 *als solcher* nie mehr
als die erste, unvollkommene Linie sein kann.

## Verdikt

Der 3×-Verdikt **gilt weiter** und wird **nicht** per Folge-Verdikt
abgelöst — seine Kernanalyse (kein repo-weiter Textmuster-Sensor
praktikabel, wegen der weiterhin legitim wachsenden
Testfall-Provenienz-Konvention) bleibt richtig und ist vom vierten Fall
unberührt: F-1 aus `slice-052` ist strukturell derselbe Fall wie die drei
vorherigen (Godoc über Produktionscode, Slice-Nummer statt ADR-Bezug),
kein neuer Diskriminierungs-Fall, der die Sensor-Verwerfung in Frage
stellen würde.

Dies ist keines der drei kanonischen Konflikt-Pfad-Verdikte im engeren
Sinn (Modul 8) — es liegt kein Rollen-Widerspruch vor (Reviewer-Fund
sofort akzeptiert, keine Implementer-Gegenrede), keine ADR-Fehlbehauptung
im Plan, keine stillschweigende Lockerung. Die vierte, begründete Handlung
(Modul 8 lässt sie ausdrücklich zu, wo keines der drei passt):

**Status quo bestätigt, aber die bisher implizite Rollenteilung wird
explizit gemacht — zwei kleine, gezielte Verkörperungen statt eines
neuen Sensors:**

1. **`.harness/skills/reviewer.md`** bekommt einen eigenen, benannten
   HIGH-Unterpunkt „Slice-/Wellen-Chronik in Produktionscode-Kommentar"
   statt die Klasse nur implizit unter „Kommentar trägt keine der
   Kommentar-Klassen" zu führen. Begründung: Der Reviewer ist die
   tragende Instanz (4/4 Trefferquote); eine explizite, benannte Prüfregel
   macht diese Trefferquote weniger vom Zufall der jeweiligen
   Reviewer-Session abhängig und robuster gegen einen künftigen
   Reviewer-Lauf, der den Fund unter der generischen Klasse übersieht.
   Das ist dieselbe Logik, die die vier bereits benannten repo-spezifischen
   HIGH-Klassen tragen: eine wiederkehrende Klasse verdient einen eigenen
   Anker, kein Sammelbecken.
2. **`.claude/commands/implement-slice.md` Schritt 20** bekommt eine
   Grenz-Klarstellung: Der Schritt bleibt Pflicht (er senkt nachweislich
   die Zahl der Fixrunden — 3 von 4 Fällen wurden vor `slice-052` durch
   ihn selbst oder durch den unmittelbar folgenden Review-Zyklus
   behandelt), wird aber ausdrücklich als **erste, nicht tragende**
   Verteidigungslinie benannt. Ein Auftreten trotz gelaufenem Schritt 20
   ist kein Beleg für einen defekten Prozess, solange der Reviewer den
   Fall vor Merge fängt — das ist der Regelfall, für den das Design
   sorgt.

**Kein neuer Sensor, kein Gate, keine weitere Schärfung der
Enumerationslogik in Schritt 20** — aus dem oben genannten Grund: Das
Pattern deckt den Fall bereits ab; das Problem lag nicht am Muster.

## Was dieses Verdikt NICHT tut

- Es legt **keinen** neuen Eintrag in `evidence/` an und bumpt **keinen**
  Zähler. Der vierte Beleg (`slice-052`) wird regulär bei der
  `slice-052`-Closure durch den Planner als `evidence/slice-052.md`
  eingetragen (Modul 6 §Das Beobachtungs-Register, „Eingetragen wird bei
  der Slice-Closure") — dieser Zug läuft davor und unabhängig davon.
- Es überschreibt **nicht** den 3×-Verdikt (kein `supersedes`) — dessen
  Sensor-Analyse bleibt gültig.
- Für den **Lese-Schritt bei `slice-052`s Closure** ist die Antwort damit
  vorweggenommen: Ausgang **verkörpert** (erneut), Herkunfts-Anker
  `seit slice-052` auf die beiden Änderungen dieses Zugs (den neuen
  Reviewer-HIGH-Punkt und die Schritt-20-Grenzklärung) — der Planner muss
  bei der Closure keine neue Architect-Eskalation mehr auslösen, nur den
  Zähler und den Verweis nachtragen.

## Was das für Implementer und Reviewer bedeutet

Implementer: unverändertes Verhalten — Schritt 20 wie bisher ausführen,
jetzt mit der zusätzlichen Erwartungsklarheit, dass er nicht als
alleiniger Schutz zählt. Reviewer: unverändertes Verhalten in der Sache
(die Klasse wurde bereits zuverlässig gefunden), jetzt mit einem
eigenen, direkt zitierbaren HIGH-Anker statt einer Subsumtion unter die
generische Kommentar-Klassen-Regel.

Weder Produktionscode noch eine ADR-Datei wurden im Rahmen dieses
Verdikts geändert.
